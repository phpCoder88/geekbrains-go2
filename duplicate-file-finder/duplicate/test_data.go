package duplicate

var FilesTestData = []struct {
	Name             string
	StartDir         string
	MaxDepth         int
	WantResult       Files
	WantDeletedFiles []string
	WantPresentFiles []string
	WantPrinted      string
}{
	{
		Name:     "Max Depth 0",
		StartDir: "./tmp",
		MaxDepth: 0,
		WantResult: Files{
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
		WantDeletedFiles: []string{
			"tmp/A/copy1.txt",
			"tmp/A/AA/copy1.txt",
			"tmp/B/copy2.txt",
		},
		WantPresentFiles: []string{
			"tmp/unique.txt",
			"tmp/A/AB/copy1.txt",
		},
		WantPrinted: `   File Name|            File Path|   File Size|
   copy1.txt|        tmp/copy1.txt|          28|
   copy1.txt|      tmp/A/copy1.txt|          28|
   copy1.txt|   tmp/A/AA/copy1.txt|          28|
   copy2.txt|        tmp/copy2.txt|          28|
   copy2.txt|      tmp/B/copy2.txt|          28|
`,
	},

	{
		Name:     "Max Depth 2",
		StartDir: "./tmp",
		MaxDepth: 2,
		WantResult: Files{
			"copy1.txt_28": []File{
				{Name: "copy1.txt", Path: "tmp/copy1.txt", Size: 28},
				{Name: "copy1.txt", Path: "tmp/A/copy1.txt", Size: 28},
			},
			"copy2.txt_28": []File{
				{Name: "copy2.txt", Path: "tmp/copy2.txt", Size: 28},
				{Name: "copy2.txt", Path: "tmp/B/copy2.txt", Size: 28},
			},
		},
		WantDeletedFiles: []string{
			"tmp/A/copy1.txt",
			"tmp/B/copy2.txt",
		},
		WantPresentFiles: []string{
			"tmp/unique.txt",
			"tmp/A/AA/copy1.txt",
			"tmp/A/AB/copy1.txt",
		},
		WantPrinted: `   File Name|         File Path|   File Size|
   copy1.txt|     tmp/copy1.txt|          28|
   copy1.txt|   tmp/A/copy1.txt|          28|
   copy2.txt|     tmp/copy2.txt|          28|
   copy2.txt|   tmp/B/copy2.txt|          28|
`,
	},
}
