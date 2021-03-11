package gen

import (
	"fmt"
	"time"
)

func Counter() {
	go func() {
		fmt.Println("go 1")
	}()

	if true {
		go func() {
			fmt.Println("go 2")
		}()
	}

	go func() {
		fmt.Println("go 3")
	}()

	time.Sleep(time.Second)
}

func Counter2() {
	go func() {
		fmt.Println("test")
	}()

	time.Sleep(time.Second)
}
