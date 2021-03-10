package main

import (
	"context"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"syscall"
	"time"
)

type TimeTicker struct {
	ticker   *time.Ticker
	stopChan chan struct{}
}

func main() {
	interrupt := make(chan os.Signal, 1)
	// c Golang 1.16 можно использовать NotifyContext из os/signal вместо Notify
	signal.Notify(interrupt, syscall.SIGTERM, syscall.SIGINT)

	service := NewTimeTicker()
	go service.Start()

	<-interrupt

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	service.Shutdown(ctx)
}

func NewTimeTicker() *TimeTicker {
	return &TimeTicker{
		stopChan: make(chan struct{}, 1),
	}
}

func (tt *TimeTicker) Start() {
	tt.ticker = time.NewTicker(time.Second)
	defer tt.ticker.Stop()

	for {
		select {
		case <-tt.ticker.C:
			fmt.Println(time.Now().Format("02.01.2006 15:04:05"))

		case <-tt.stopChan:
			// Эмуляция остановки сервиса
			r := rand.New(rand.NewSource(time.Now().UnixNano()))
			stopTime := time.Duration(r.Intn(2000))
			fmt.Println(stopTime)
			time.Sleep(stopTime * time.Millisecond)
			tt.stopChan <- struct{}{}
			return
		}
	}

}
func (tt *TimeTicker) Shutdown(ctx context.Context) {
	tt.stopChan <- struct{}{}

	select {
	case <-tt.stopChan:
		fmt.Println("Graceful shutdown!!!")
		return

	case <-ctx.Done():
		fmt.Println("Time is out. Shutdown!")
		return
	}
}
