package duplicate

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

// MockNotFileExistErr
var MockNotFileExistErr = errors.New("no such file or directory")

// FileSystemStruct описывает структуру мока файловой системы
type FileSystemStruct map[string]map[string]FileInfoMock

// FileSystemTree описывает содержимое мока файловой системы
var FileSystemTree = FileSystemStruct{
	"tmp": {
		"unique.txt": {name: "unique.txt", content: "Unique content for ./unique.txt"},
		"copy1.txt":  {name: "copy1.txt", content: "Some content for ./copy1.txt"},
		"copy2.txt":  {name: "copy2.txt", content: "Some content for ./copy2.txt"},
		"A":          {name: "A", isDir: true},
		"B":          {name: "B", isDir: true},
	},
	"tmp/A": {
		"copy1.txt": {name: "copy1.txt", content: "Some content for ./copy1.txt"},
		"AA":        {name: "AA", isDir: true},
		"AB":        {name: "AB", isDir: true},
	},
	"tmp/B": {
		"copy2.txt": {name: "copy2.txt", content: "Some content for ./copy2.txt"},
	},
	"tmp/A/AA": {
		"copy1.txt": {name: "copy1.txt", content: "Some content for ./copy1.txt"},
	},
	"tmp/A/AB": {
		"copy1.txt": {name: "copy1.txt", content: "Some other content for ./copy1.txt"},
	},
}

// FileSystemMock описывает мок файловой системы
type FileSystemMock struct {
	fileSystem FileSystemStruct
}

// NewFileSystemMock создает мок файловой системы
func NewFileSystemMock(fileSystem FileSystemStruct) *FileSystemMock {
	return &FileSystemMock{
		fileSystem: fileSystem,
	}
}

// ReadDir читает содержимое директории в FileSystemMock
func (dr *FileSystemMock) ReadDir(path string) ([]os.FileInfo, error) {
	if fileInfos, ok := dr.fileSystem[path]; ok {
		files := make([]os.FileInfo, len(fileInfos))

		var i int
		for _, file := range fileInfos {
			files[i] = file
			i++
		}

		return files, nil
	}

	return nil, fmt.Errorf("open %s: %w", path, MockNotFileExistErr)
}

// Remove удаляет файл из FileSystemMock
func (dr *FileSystemMock) Remove(path string) error {
	dir := filepath.Dir(path)
	if _, ok := dr.fileSystem[dir]; !ok {
		return fmt.Errorf("open %s: %w", path, MockNotFileExistErr)
	}

	filename := filepath.Base(path)
	if _, ok := dr.fileSystem[dir][filename]; !ok {
		return fmt.Errorf("stat %s: %w", path, MockNotFileExistErr)
	}

	delete(dr.fileSystem[dir], filename)
	return nil
}

// FileInfoMock описывает мок файла
type FileInfoMock struct {
	name     string
	mode     os.FileMode
	modeTime time.Time
	isDir    bool
	content  string
}

// Name возвращает имя файла
func (f FileInfoMock) Name() string {
	return f.name
}

// Size возвращает размер файла
func (f FileInfoMock) Size() int64 {
	return int64(len(f.content))
}

// Mode возвращает размер файла
func (f FileInfoMock) Mode() os.FileMode {
	return f.mode
}

// ModTime возвращает время изменения файла
func (f FileInfoMock) ModTime() time.Time {
	return f.modeTime
}

// IsDir возвращает true если файл - директория, и false в противном случае
func (f FileInfoMock) IsDir() bool {
	return f.isDir
}

// Sys underlying data source (can return nil)
func (f FileInfoMock) Sys() interface{} {
	return nil
}
