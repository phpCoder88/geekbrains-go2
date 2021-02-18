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
}

func PrintData(reader *bufio.Reader) (err error) {
	defer func() {
		if errFromPanic := recover(); errFromPanic != nil {
			fmt.Fprintln(os.Stderr, errFromPanic)
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
