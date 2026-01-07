package main

import (
	"fmt"
	"time"

	"github.com/nimxch/joker/learn"
)

func main() {
	fmt.Println("Hello Joker")
	i := 0
	for {
		learn.TestFsync(i)
		i++
		time.Sleep(2 * time.Second)
	}
}
