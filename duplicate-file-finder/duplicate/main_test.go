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
	for ind, tt := range tests {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.name, func(t *testing.T) {
			dFiles := s.finder.Seek(tt.startDir, tt.maxDepth)
			assert.Equal(t, tt.wantResult, dFiles)
		})
	}
}

func (s *MemoryDuplicatesTestSuite) TestPrintDuplicates() {
	for ind, tt := range tests {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.name, func(t *testing.T) {
			_ = s.finder.Seek(tt.startDir, tt.maxDepth)

			out := new(bytes.Buffer)
			s.finder.PrintDuplicates(out)
			result := out.String()
			assert.Equal(t, tt.wantPrinted, result)
		})
	}
}

func (s *MemoryDuplicatesTestSuite) TestRemoveAllDuplicates() {
	for ind, tt := range tests {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.name, func(t *testing.T) {
			_ = s.finder.Seek(tt.startDir, tt.maxDepth)
			mock := s.finder.fs.(*FileSystemMock)

			s.finder.RemoveAllDuplicates()

			for _, filePath := range tt.wantDeletedFiles {
				dir := filepath.Dir(filePath)
				filename := filepath.Base(filePath)
				if _, fileExists := mock.fileSystem[dir][filename]; fileExists {
					s.T().Errorf("File %q exists", filePath)
				}
			}

			for _, filePath := range tt.wantPresentFiles {
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
