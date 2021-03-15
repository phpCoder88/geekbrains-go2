// Package main консольная команда для поиска и удаления дубликатов файлов
package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/phpCoder88/geekbrains-go2/duplicate-file-finder/duplicate"
)

var isRemove = flag.Bool("remove", false, "удалять дубликаты файлов")
var startDir = flag.String("path", ".", "Стартовая директория для поиска")
var maxDepth = flag.Int("maxdepth", 0, "максимальная глубина поиска по подкаталогам. --maxdepth <= 0 нет ограничений на вложенность")

func main() {
	flag.Parse()

	finder := duplicate.NewDuplicateFinder()

	files := finder.Seek(*startDir, *maxDepth)

	finder.PrintDuplicates(os.Stdout)

	if *isRemove && len(files) > 0 {
		var removeConfirm string
		fmt.Print("Удалить дубликаты(Y/n): ")
		_, err := fmt.Scanln(&removeConfirm)
		if err != nil {
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		removeConfirm = strings.ToLower(strings.TrimSpace(removeConfirm))
		if removeConfirm != "y" && removeConfirm != "yes" {
			return
		}

		log.Println("Removing...")
		finder.RemoveAllDuplicates()
	}
}
