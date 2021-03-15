package duplicate

import (
	"bytes"
	"fmt"
	"os"
	"path"
	"reflect"
	"sort"
	"testing"
	"text/tabwriter"
)

var tests = []struct {
	name             string
	maxDepth         int
	wantResult       Files
	wantDeletedFiles []string
	wantPresentFiles []string
}{
	{
		name:     "Max Depth 0",
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
	},

	{
		name:     "Max Depth 2",
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
	},
}

func TestDuplicates_Seek(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, files, tearDown := newTestFiles(t, tt.maxDepth)
			defer tearDown()

			if !reflect.DeepEqual(files, tt.wantResult) {
				t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", files, tt.wantResult)
			}
		})
	}
}

func TestDuplicates_PrintDuplicates(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder, _, tearDown := newTestFiles(t, tt.maxDepth)
			defer tearDown()

			out := new(bytes.Buffer)
			finder.PrintDuplicates(out)
			result := out.String()

			wantResult := printedResult(tt.wantResult)
			if result != wantResult {
				t.Errorf("test Failed - results not match\nGot:\n%v\nExpected:\n%v", result, wantResult)
			}
		})
	}
}

func TestDuplicates_RemoveAllDuplicates(t *testing.T) {
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			finder, _, tearDown := newTestFiles(t, tt.maxDepth)
			defer tearDown()

			finder.RemoveAllDuplicates()

			for _, filePath := range tt.wantDeletedFiles {
				if _, err := os.Stat(filePath); err == nil {
					t.Fatalf("test Failed - File %s was not deleted\n", filePath)
				}
			}

			for _, filePath := range tt.wantPresentFiles {
				if _, err := os.Stat(filePath); os.IsNotExist(err) {
					t.Fatalf("test Failed - File %s was deleted\n", filePath)
				}
			}
		})
	}

}

func newTestFiles(t *testing.T, maxDepth int) (*Duplicates, Files, func()) {
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
			t.Fatal(err)
		}

		file, err := os.Create(fileItem.path)
		if err != nil {
			t.Fatal(err)
		}

		_, _ = file.WriteString(fileItem.content)
		_ = file.Close()
	}

	finder := NewDuplicateFinder()
	dFiles := finder.Seek("./tmp", maxDepth)

	return finder, dFiles, func() {
		err := os.RemoveAll("./tmp")
		if err != nil {
			t.Fatal(err)
		}
	}
}

func printedResult(wantResult Files) string {
	out := new(bytes.Buffer)

	w := tabwriter.NewWriter(out, 0, 0, 3, ' ', tabwriter.AlignRight|tabwriter.Debug)
	_, _ = fmt.Fprintf(w, "%s\t%s\t%s\t\n", "File Name", "File Path", "File Size")

	keys := make([]string, 0, len(wantResult))
	for k := range wantResult {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	for _, key := range keys {
		for _, file := range wantResult[key] {
			_, _ = fmt.Fprintf(w, "%s\t%s\t%d\t\n", file.Name, file.Path, file.Size)
		}
	}
	_ = w.Flush()

	return out.String()
}
