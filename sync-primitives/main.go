package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

func main() {
	waits()
}

func waits() {
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	var numStreams = r.Intn(200)
	var wg sync.WaitGroup
	var streamCounter int
	var mu sync.Mutex

	for i := 0; i < numStreams; i++ {
		wg.Add(1)
		go func(counter int) {
			defer wg.Done()

			mu.Lock()
			defer mu.Unlock()
			streamCounter++

			fmt.Printf("Hello from goroutine %d\n", counter)
		}(i)
	}

	wg.Wait()

	fmt.Println(streamCounter)
}
