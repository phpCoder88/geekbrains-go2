// Package duplicate ищет и удаляет дубликаты файлов(одинаковые имя и размер файлов)
package duplicate

import (
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"path/filepath"
	"sort"
	"sync"
	"text/tabwriter"

	"go.uber.org/zap"
)

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
	sync.Mutex
	files Files
	sync.WaitGroup
	logger *zap.Logger
}

// NewDuplicateFinder инициализирует поиск
func NewDuplicateFinder(logger *zap.Logger) *Duplicates {
	return &Duplicates{
		files:  make(Files),
		logger: logger,
	}
}

// Seek ищет дубликаты файлов
func (f *Duplicates) Seek(path string, maxDepth int) Files {
	f.Add(1)
	go f.scanDir(path, maxDepth, 1)

	f.Wait()

	f.filterFiles()

	return f.files
}

// scanDir рекурсивно сканирует директории в поисках дубликатов
func (f *Duplicates) scanDir(path string, maxDepth int, level int) {
	defer f.Done()

	f.logger.Info("Start scanning dir " + path)
	list, err := ioutil.ReadDir(path)
	if err != nil {
		f.logger.Error("Can't read dir " + path)
		_, _ = fmt.Fprintln(os.Stderr, err)
		return
	}

	for _, val := range list {
		currPath := filepath.Join(path, val.Name())

		if val.IsDir() {
			if maxDepth <= 0 || level < maxDepth {
				f.Add(1)
				go f.scanDir(currPath, maxDepth, level+1)
			}
			continue
		}

		fileToken := fmt.Sprintf("%s_%d", val.Name(), val.Size())
		f.Lock()
		f.files[fileToken] = append(f.files[fileToken], File{
			Name: val.Name(),
			Path: currPath,
			Size: val.Size(),
		})
		f.Unlock()
	}
}

// RemoveAllDuplicates удаляет все дубликаты файлов
func (f *Duplicates) RemoveAllDuplicates() {
	for fileSetKey := range f.files {
		f.Add(1)
		go f.removeFileDuplicates(fileSetKey)
	}

	f.Wait()
}

// removeFileDuplicates удаляет дубликаты одного файла
func (f *Duplicates) removeFileDuplicates(fileSetKey string) {
	defer f.Done()

	files, ok := f.files[fileSetKey]
	if !ok {
		return
	}

	for fileInd, file := range files {
		if fileInd != 0 {
			f.logger.Info("Removing file " + file.Path)
			err := os.Remove(file.Path)
			if err != nil {
				f.logger.Error("Removing file " + file.Path)
				_, _ = fmt.Fprintln(os.Stderr, err)
			}
		}
	}
}

// filterFiles фильтрует найденные файлы и сортирует дубликаты
func (f *Duplicates) filterFiles() {
	for ind, dFiles := range f.files {
		if len(dFiles) < 2 {
			delete(f.files, ind)
			continue
		}

		sort.Sort(byFilePath(dFiles))
	}
}

// PrintDuplicates Вывод найденных дубликатов
func (f *Duplicates) PrintDuplicates(out io.Writer) {
	if len(f.files) == 0 {
		return
	}

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n", "File Name", "File Path", "File Size")

	keys := make([]string, 0, len(f.files))
	for k := range f.files {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, file := range f.files[key] {
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
