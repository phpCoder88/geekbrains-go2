package main

import (
	"fmt"
	"sync"
)

func main() {
	waits()
}

func waits() {
	const numStreams = 100
	var wg sync.WaitGroup

	for i := 0; i < numStreams; i++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()
			// Do something
			fmt.Printf("Hello from goroutine %d\n", counter)
		}(i)
	}

	wg.Wait()
}
