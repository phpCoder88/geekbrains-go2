package main

import (
	"fmt"
	"os"
	"runtime/trace"
	"sync"
)

func main() {
	task1()
}

func task1() {
	_ = trace.Start(os.Stderr)
	defer trace.Stop()

	var mu sync.Mutex
	var wg sync.WaitGroup
	var counter = 0

	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()

			mu.Lock()
			counter++
			mu.Unlock()
		}()
	}

	wg.Wait()
	fmt.Println(counter)
}
