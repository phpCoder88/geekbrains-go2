// Package duplicate ищет и удаляет дубликаты файлов(одинаковые имя и размер файлов)
package duplicate

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path"
	"path/filepath"
	"sort"
	"sync"
	"text/tabwriter"

	"go.uber.org/zap"
)

// FSReader описывает чтение директории
type FSReader interface {
	ReadDir(path string) ([]os.FileInfo, error)
}

// FSDeleter описывает удаление файла
type FSDeleter interface {
	Remove(name string) error
}

// FSReadDeleter описывает чтение директории и удаление файла
type FSReadDeleter interface {
	FSReader
	FSDeleter
}

// FileSystem представляет работу с файловой системой
type FileSystem struct{}

// ReadDir читает содержимое указанной директории
func (dr FileSystem) ReadDir(dirPath string) ([]os.FileInfo, error) {
	return ioutil.ReadDir(dirPath)
}

// Remove удаляет файл по указанному пути
func (dr FileSystem) Remove(name string) error {
	return os.Remove(name)
}

// File описывает единичный файл в поиске
type File struct {
	Name string
	Path string
	Size int64
}

// Files описывает все найденные файлы, сгруппированные по копиям
type Files map[string][]File

// Duplicates управляет поиском дубликатов
type Duplicates struct {
	fs FSReadDeleter
	sync.Mutex
	files Files
	sync.WaitGroup
	logger *zap.Logger
}

// NewDuplicateFinder инициализирует поиск
func NewDuplicateFinder(fs FSReadDeleter, logger *zap.Logger) *Duplicates {
	return &Duplicates{
		fs:     fs,
		files:  make(Files),
		logger: logger,
	}
}

// Seek ищет дубликаты файлов
func (d *Duplicates) Seek(startPath string, maxDepth int) Files {
	d.Add(1)
	go d.scanDir(path.Clean(startPath), maxDepth, 1)

	d.Wait()

	d.filterFiles()

	return d.files
}

// scanDir рекурсивно сканирует директории в поисках дубликатов
func (d *Duplicates) scanDir(dirPath string, maxDepth, level int) {
	defer d.Done()

	d.logger.Info("Start scanning dir " + dirPath)
	list, err := d.fs.ReadDir(dirPath)
	if err != nil {
		d.logger.Error("Can't read dir " + dirPath)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, val := range list {
		currPath := filepath.Join(dirPath, val.Name())

		if val.IsDir() {
			if maxDepth <= 0 || level < maxDepth {
				d.Add(1)
				go d.scanDir(currPath, maxDepth, level+1)
			}
			continue
		}

		fileToken := fmt.Sprintf("%s_%d", val.Name(), val.Size())
		d.Lock()
		d.files[fileToken] = append(d.files[fileToken], File{
			Name: val.Name(),
			Path: currPath,
			Size: val.Size(),
		})
		d.Unlock()
	}
}

// RemoveAllDuplicates удаляет все дубликаты файлов
func (d *Duplicates) RemoveAllDuplicates() {
	for fileSetKey := range d.files {
		d.Add(1)
		go d.removeFileDuplicates(fileSetKey)
	}

	d.Wait()
}

// removeFileDuplicates удаляет дубликаты одного файла
func (d *Duplicates) removeFileDuplicates(fileSetKey string) {
	defer d.Done()

	files, ok := d.files[fileSetKey]
	if !ok {
		return
	}

	for fileInd, file := range files {
		if fileInd != 0 {
			d.logger.Info("Removing file " + file.Path)
			err := d.fs.Remove(file.Path)
			if err != nil {
				d.logger.Error("Removing file " + file.Path)
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// filterFiles фильтрует найденные файлы и сортирует дубликаты
func (d *Duplicates) filterFiles() {
	minFileFilter := 2
	for ind, dFiles := range d.files {
		if len(dFiles) < minFileFilter {
			delete(d.files, ind)
			continue
		}

		sort.Sort(byFilePath(dFiles))
	}
}

// PrintDuplicates Вывод найденных дубликатов
func (d *Duplicates) PrintDuplicates(out io.Writer) {
	if len(d.files) == 0 {
		return
	}

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n", "File Name", "File Path", "File Size")

	keys := make([]string, 0, len(d.files))
	for k := range d.files {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, file := range d.files[key] {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t\n", file.Name, file.Path, file.Size)
		}
	}
	_ = w.Flush()
}

// byFilePath Сортирует слайс имен файлов
type byFilePath []File

func (f byFilePath) Len() int {
	return len(f)
}

func (f byFilePath) Swap(i, j int) {
	f[i], f[j] = f[j], f[i]
}

func (f byFilePath) Less(i, j int) bool {
	return len(f[i].Path) < len(f[j].Path)
}
