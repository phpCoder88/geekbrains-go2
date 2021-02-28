package main

import (
	"fmt"
	"sync"
	"sync/atomic"
)

func main() {
	fmt.Println(channelCounter(50, 1000))
	fmt.Println(atomicCounter(50, 1000))
}

func channelCounter(workerCount int, jobCount int) int {
	var counter int

	jobs := make(chan struct{})
	inc := make(chan struct{}, jobCount)
	wg := &sync.WaitGroup{}
	worker := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for range jobs {
			inc <- struct{}{}
		}
	}
	jobMaker := func() {
		for j := 1; j <= jobCount; j++ {
			jobs <- struct{}{}
		}
		close(jobs)
	}

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(wg)
	}
	go jobMaker()

	wg.Wait()
	close(inc)

	for range inc {
		counter++
	}

	return counter
}

func atomicCounter(workerCount int, jobCount int) int {
	var counter int32
	wg := &sync.WaitGroup{}
	jobs := make(chan struct{})
	worker := func(wg *sync.WaitGroup) {
		defer wg.Done()
		for range jobs {
			atomic.AddInt32(&counter, 1)
		}
	}
	jobMaker := func() {
		for j := 1; j <= jobCount; j++ {
			jobs <- struct{}{}
		}
		close(jobs)
	}

	for w := 1; w <= workerCount; w++ {
		wg.Add(1)
		go worker(wg)
	}

	go jobMaker()
	wg.Wait()

	return int(counter)
}
