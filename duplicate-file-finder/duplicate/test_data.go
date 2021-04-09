package duplicate

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
