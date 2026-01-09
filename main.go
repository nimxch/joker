package main

import (
	"fmt"
	"io"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/nimxch/joker/learn"
)

const (
	MAX_NODE_SIZE = 1024 * 8 // 8KB fits into CPU L1/L2 Cache
)

func main() {
	fmt.Println("Hello Joker")
	q := &learn.Queue{}
	q.Enqueue([]byte("hello\n"))
	node, ok := q.Peek()
	if !ok {
		fmt.Println("NOT OK")
	}
	fmt.Print(string(node))
	q.Enqueue([]byte("World!\n"))
	node, _ = q.Peek()
	fmt.Print(string(node))
	q.Enqueue([]byte("Nimai!\n"))
	q.Enqueue([]byte("Charan Nimai!\n"))
	node, _ = q.Peek()
	q.Dequeue()
	deq, _ := q.Dequeue()
	fmt.Print(string(node))
	fmt.Print("Dequeue: ", string(deq))
}

func InitServerSync() {
	listener, err := net.Listen("tcp", ":8080")
	if err != nil {
		fmt.Println("Error:", err)
		return
	}

	// Handle Ctrl+C
	go func() {
		var sig os.Signal
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, os.Interrupt, syscall.SIGTERM)
		sig = <-sigCh
		fmt.Println("\nReceived signal:", sig)
		listener.Close()
	}()

	for {
		conn, err := listener.Accept()
		if err != nil {
			fmt.Println("Server shutting down")
			return
		}

		fmt.Println("Handling connection")
		go handleConnection(conn)
	}
}

func handleConnection(conn net.Conn) {
	defer conn.Close()

	buffer := make([]byte, 1024)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			if err != io.EOF {
				fmt.Println("Error while reading data: ", err)
			}
			fmt.Println("Error: ", err)
			break
		}

		fmt.Println(string(buffer[:n]))
	}
}
