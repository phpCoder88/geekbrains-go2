package duplicate

import (
	"bytes"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type MemoryDuplicatesTestSuite struct {
	suite.Suite
	finder *Duplicates
}

func (s *MemoryDuplicatesTestSuite) SetupTest() {
	FileSystemTreeCopy := make(FileSystemStruct, len(FileSystemTree))

	for dir, dirValue := range FileSystemTree {
		FileSystemTreeCopy[dir] = make(map[string]FileInfoMock, len(dirValue))

		for file, fileInfo := range dirValue {
			FileSystemTreeCopy[dir][file] = fileInfo
		}
	}

	logger := zaptest.NewLogger(s.T())
	fs := NewFileSystemMock(FileSystemTreeCopy)
	s.finder = NewDuplicateFinder(fs, logger)
}

func (s *MemoryDuplicatesTestSuite) TestDuplicatesSeek() {
	for ind, tt := range FilesTestData {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.Name, func(t *testing.T) {
			dFiles := s.finder.Seek(tt.StartDir, tt.MaxDepth)
			assert.Equal(t, tt.WantResult, dFiles)
		})
	}
}

func (s *MemoryDuplicatesTestSuite) TestPrintDuplicates() {
	for ind, tt := range FilesTestData {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.Name, func(t *testing.T) {
			_ = s.finder.Seek(tt.StartDir, tt.MaxDepth)

			out := new(bytes.Buffer)
			s.finder.PrintDuplicates(out)
			result := out.String()
			assert.Equal(t, tt.WantPrinted, result)
		})
	}
}

func (s *MemoryDuplicatesTestSuite) TestRemoveAllDuplicates() {
	for ind, tt := range FilesTestData {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.Name, func(t *testing.T) {
			_ = s.finder.Seek(tt.StartDir, tt.MaxDepth)
			mock := s.finder.fs.(*FileSystemMock)

			s.finder.RemoveAllDuplicates()

			for _, filePath := range tt.WantDeletedFiles {
				dir := filepath.Dir(filePath)
				filename := filepath.Base(filePath)
				if _, fileExists := mock.fileSystem[dir][filename]; fileExists {
					s.T().Errorf("File %q exists", filePath)
				}
			}

			for _, filePath := range tt.WantPresentFiles {
				dir := filepath.Dir(filePath)
				filename := filepath.Base(filePath)
				if _, fileExists := mock.fileSystem[dir][filename]; !fileExists {
					s.T().Errorf("File %q does not exist", filePath)
				}
			}
		})
	}
}

func TestMemoryDuplicatesTestSuite(t *testing.T) {
	suite.Run(t, new(MemoryDuplicatesTestSuite))
}
