package main

import (
	"fmt"
	"os"
	"runtime"
	"runtime/trace"
	"sync"
)

const (
	iterationsNum = 7
	goroutinesNum = 5
)

func main() {
	_ = trace.Start(os.Stderr)
	defer trace.Stop()

	var wg sync.WaitGroup

	for i := 0; i < goroutinesNum; i++ {
		wg.Add(1)
		go startWorker(i, &wg)
	}

	wg.Wait()
}

func startWorker(in int, wg *sync.WaitGroup) {
	defer wg.Done()

	for j := 0; j < iterationsNum; j++ {
		fmt.Printf("thread %d iteration %d\n", in, j)
		runtime.Gosched()
	}
}
