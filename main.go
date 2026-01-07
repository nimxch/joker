package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	fmt.Println("Hello Joker")
	shutdown := make(chan struct{})
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)

	go func() {
		<-sigCh
		close(shutdown)
	}()

	for {
		select {
		case <-shutdown:
			fmt.Println("Interrupted....")
			return
		default:
			fmt.Println("Normal...")
			time.Sleep(time.Second)
		}
	}
}
