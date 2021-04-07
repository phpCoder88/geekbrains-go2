package duplicate

import (
	"bytes"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/zap/zaptest"

	"github.com/stretchr/testify/suite"
)

var tests = []struct {
	name             string
	startDir         string
	maxDepth         int
	wantResult       Files
	wantDeletedFiles []string
	wantPresentFiles []string
	wantPrinted      string
}{
	{
		name:     "Max Depth 0",
		startDir: "./tmp",
		maxDepth: 0,
		wantResult: Files{
			"copy1.txt_28": []File{
				{Name: "copy1.txt", Path: "tmp/copy1.txt", Size: 28},
				{Name: "copy1.txt", Path: "tmp/A/copy1.txt", Size: 28},
				{Name: "copy1.txt", Path: "tmp/A/AA/copy1.txt", Size: 28},
			},
			"copy2.txt_28": []File{
				{Name: "copy2.txt", Path: "tmp/copy2.txt", Size: 28},
				{Name: "copy2.txt", Path: "tmp/B/copy2.txt", Size: 28},
			},
		},
		wantDeletedFiles: []string{
			"tmp/A/copy1.txt",
			"tmp/A/AA/copy1.txt",
			"tmp/B/copy2.txt",
		},
		wantPresentFiles: []string{
			"tmp/unique.txt",
			"tmp/A/AB/copy1.txt",
		},
		wantPrinted: `   File Name|            File Path|   File Size|
   copy1.txt|        tmp/copy1.txt|          28|
   copy1.txt|      tmp/A/copy1.txt|          28|
   copy1.txt|   tmp/A/AA/copy1.txt|          28|
   copy2.txt|        tmp/copy2.txt|          28|
   copy2.txt|      tmp/B/copy2.txt|          28|
`,
	},

	{
		name:     "Max Depth 2",
		startDir: "./tmp",
		maxDepth: 2,
		wantResult: Files{
			"copy1.txt_28": []File{
				{Name: "copy1.txt", Path: "tmp/copy1.txt", Size: 28},
				{Name: "copy1.txt", Path: "tmp/A/copy1.txt", Size: 28},
			},
			"copy2.txt_28": []File{
				{Name: "copy2.txt", Path: "tmp/copy2.txt", Size: 28},
				{Name: "copy2.txt", Path: "tmp/B/copy2.txt", Size: 28},
			},
		},
		wantDeletedFiles: []string{
			"tmp/A/copy1.txt",
			"tmp/B/copy2.txt",
		},
		wantPresentFiles: []string{
			"tmp/unique.txt",
			"tmp/A/AA/copy1.txt",
			"tmp/A/AB/copy1.txt",
		},
		wantPrinted: `   File Name|         File Path|   File Size|
   copy1.txt|     tmp/copy1.txt|          28|
   copy1.txt|   tmp/A/copy1.txt|          28|
   copy2.txt|     tmp/copy2.txt|          28|
   copy2.txt|   tmp/B/copy2.txt|          28|
`,
	},
}

type DuplicatesTestSuite struct {
	suite.Suite
	finder *Duplicates
}

func (s *DuplicatesTestSuite) SetupTest() {
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
	s.finder = NewDuplicateFinder(logger)
}

func (s *DuplicatesTestSuite) TearDownTest() {
	err := os.RemoveAll("./tmp")
	if err != nil {
		s.T().Fatal(err)
	}
}

func (s *DuplicatesTestSuite) TestDuplicatesSeek() {
	tt := tests[0]
	s.T().Run(tt.name, func(t *testing.T) {
		dFiles := s.finder.Seek(tt.startDir, tt.maxDepth)
		assert.Equal(t, tt.wantResult, dFiles)
	})
}

func (s *DuplicatesTestSuite) TestDuplicatesSeek_Depth2() {
	tt := tests[1]
	s.T().Run(tt.name, func(t *testing.T) {
		dFiles := s.finder.Seek(tt.startDir, tt.maxDepth)
		assert.Equal(t, tt.wantResult, dFiles)
	})
}

func (s *DuplicatesTestSuite) TestPrintDuplicates() {
	tt := tests[0]
	s.T().Run(tt.name, func(t *testing.T) {
		_ = s.finder.Seek(tt.startDir, tt.maxDepth)

		out := new(bytes.Buffer)
		s.finder.PrintDuplicates(out)
		result := out.String()
		assert.Equal(t, tt.wantPrinted, result)
	})
}

func (s *DuplicatesTestSuite) TestPrintDuplicates_Depth2() {
	tt := tests[1]
	s.T().Run(tt.name, func(t *testing.T) {
		_ = s.finder.Seek(tt.startDir, tt.maxDepth)

		out := new(bytes.Buffer)
		s.finder.PrintDuplicates(out)
		result := out.String()

		assert.Equal(t, tt.wantPrinted, result)
	})
}

func (s *DuplicatesTestSuite) TestRemoveAllDuplicates() {
	tt := tests[0]
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

func (s *DuplicatesTestSuite) TestRemoveAllDuplicates_Depth2() {
	tt := tests[1]
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

func TestExampleTestSuite(t *testing.T) {
	suite.Run(t, new(DuplicatesTestSuite))
}
