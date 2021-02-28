package main

import (
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var shutdownSuccess = make(chan struct{}, 1)

func main() {
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	stopTicker := make(chan struct{})

	go printDatetime(stopTicker)

	<-interrupt
	stopTicker <- struct{}{}

	timer := time.NewTimer(time.Second)
	select {
	case <-timer.C:
		fmt.Println("Time is out. Shutdown!")
		os.Exit(1)
	case <-shutdownSuccess:
		fmt.Println("Graceful shutdown!!!")
		return
	}
}

func printDatetime(stop <-chan struct{}) {
	ticker := time.NewTicker(time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			fmt.Println(time.Now().Format("02.01.2006 15:04:05"))
		case <-stop:
			// Эмуляция остановки сервиса
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			stopTime := time.Duration(r.Intn(2000))
			fmt.Println(stopTime)
			time.Sleep(stopTime * time.Millisecond)
			shutdownSuccess <- struct{}{}
			return
		}
	}
}
