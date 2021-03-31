// Package main консольная команда для поиска и удаления дубликатов файлов
package main

import (
	"flag"
	"fmt"
	"os"
	"strings"

	"go.uber.org/zap"

	"github.com/phpCoder88/geekbrains-go2/duplicate-file-finder/duplicate"
)

var isRemove = flag.Bool("remove", false, "удалять дубликаты файлов")
var startDir = flag.String("path", ".", "Стартовая директория для поиска")
var maxDepth = flag.Int("maxdepth", 0, "максимальная глубина поиска по подкаталогам. --maxdepth <= 0 нет ограничений на вложенность")

func main() {
	flag.Parse()

	logger, _ := zap.NewProduction()
	defer func() {
		err := logger.Sync()
		if err != nil {
			fmt.Println(err)
		}
	}()

	logger = logger.With(zap.String("startSearchingDir", *startDir)).With(zap.Int("searchingDepth", *maxDepth)).With(zap.Bool("isRemove", *isRemove))

	finder := duplicate.NewDuplicateFinder(logger)
	logger.Info("Start searching...")
	files := finder.Seek(*startDir, *maxDepth)

	logger.Info("Printing searched results...")
	finder.PrintDuplicates(os.Stdout)

	if *isRemove && len(files) > 0 {
		var removeConfirm string
		fmt.Print("Удалить дубликаты(Y/n): ")
		_, err := fmt.Scanln(&removeConfirm)
		if err != nil {
			logger.Error("Can't scan removing confirm message")
			_, _ = fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}

		removeConfirm = strings.ToLower(strings.TrimSpace(removeConfirm))
		if removeConfirm != "y" && removeConfirm != "yes" {
			return
		}

		logger.Info("Removing files...")
		finder.RemoveAllDuplicates()
	}
}
