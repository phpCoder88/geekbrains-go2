// +build integration

package duplicate

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"go.uber.org/zap/zaptest"
)

type RealDuplicateFilesTestSuite struct {
	suite.Suite
	finder *Duplicates
}

func (s *RealDuplicateFilesTestSuite) SetupTest() {
	var files = []struct {
		path    string
		content string
	}{
		{path: "./tmp/unique.txt", content: "Unique content for ./unique.txt"},
		{path: "./tmp/copy1.txt", content: "Some content for ./copy1.txt"},
		{path: "./tmp/copy2.txt", content: "Some content for ./copy2.txt"},
		{path: "./tmp/A/copy1.txt", content: "Some content for ./copy1.txt"},
		{path: "./tmp/B/copy2.txt", content: "Some content for ./copy2.txt"},
		{path: "./tmp/A/AA/copy1.txt", content: "Some content for ./copy1.txt"},
		{path: "./tmp/A/AB/copy1.txt", content: "Some other content for ./copy1.txt"},
	}

	for _, fileItem := range files {
		err := os.MkdirAll(path.Dir(fileItem.path), 0755)
		if err != nil {
			s.T().Fatal(err)
		}

		file, err := os.Create(fileItem.path)
		if err != nil {
			s.T().Fatal(err)
		}

		_, _ = file.WriteString(fileItem.content)
		_ = file.Close()
	}
	logger := zaptest.NewLogger(s.T())
	fs := FileSystem{}
	s.finder = NewDuplicateFinder(fs, logger)
}

func (s *RealDuplicateFilesTestSuite) TearDownTest() {
	err := os.RemoveAll("./tmp")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *RealDuplicateFilesTestSuite) TestDuplicatesSeek() {
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

func (s *RealDuplicateFilesTestSuite) TestPrintDuplicates() {
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

func (s *RealDuplicateFilesTestSuite) TestRemoveAllDuplicates() {
	for ind, tt := range tests {
		if ind != 0 {
			s.SetupTest()
		}

		s.T().Run(tt.name, func(t *testing.T) {
			_ = s.finder.Seek(tt.startDir, tt.maxDepth)
			s.finder.RemoveAllDuplicates()

			for _, filePath := range tt.wantDeletedFiles {
				assert.NoFileExists(t, filePath)
			}

			for _, filePath := range tt.wantPresentFiles {
				assert.FileExists(t, filePath)
			}
		})
	}
}

func TestRealDuplicateFilesTestSuite(t *testing.T) {
	suite.Run(t, new(RealDuplicateFilesTestSuite))
}
