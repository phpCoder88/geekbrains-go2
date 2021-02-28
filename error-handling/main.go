package main

import (
	"bufio"
	"fmt"
	"os"

	"github.com/phpCoder88/geekbrains-go2/error-handling/myerrors"
)

func main() {
	var reader *bufio.Reader
	fmt.Println(PrintData(reader))

	err := FileCreate("error-handling/test.txt", "Hello there\n")
	if err != nil {
		fmt.Println(err)
		return
	}
}

func PrintData(reader *bufio.Reader) (err error) {
	defer func() {
		if errFromPanic := recover(); errFromPanic != nil {
			_, _ = fmt.Fprintln(os.Stderr, errFromPanic)
			err = myerrors.NewErrorWithTime("nil pointer exception")
			return
		}
	}()

	data, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	fmt.Println(data)

	return nil
}

func FileCreate(path string, content string) error {
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	return err
}
