package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello Joker")
	ch := make(chan int, 1)

	go func() {
		time.Sleep(5 * time.Second)
		ch <- 42
	}()

	for {
		select {
		case x, ok := <-ch:
			if !ok {
				fmt.Println("Channel closed")
				return
			}
			fmt.Println("Value available:", x)
			return
		default:
			fmt.Println("No value available, doing other work...")
			time.Sleep(1 * time.Second)
		}
	}
}
