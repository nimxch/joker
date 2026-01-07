package main

import (
	"fmt"
	"time"
)

func main() {
	fmt.Println("Hello Joker")
	i := 0
	for {
		TestFsync(i)
		i++
		time.Sleep(2 * time.Second)
	}
}
