package main

import (
	"fmt"
	"os"
)

func TestFsync(iter int) {
	file_name := "wal.log"

	file, err := os.OpenFile(file_name, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println("Error: ", err)
		return
	}
	defer file.Close()
	text := fmt.Sprintf("This is a content {%d}\n", iter)
	_, err = file.Write([]byte(text))
	if err != nil {
		fmt.Println("Error: ", err)
	}
}
