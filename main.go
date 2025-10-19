package main

import (
	"fmt"
	"io"
	"os"
)

const (
	filename = "messages.txt"
)

func main() {
	file, err := os.Open(filename)
	if err != nil {
		os.Exit(1)
	}
	for {
		// Byte slice of size 8
		buffer := make([]byte, 8)

		// Read file into buffer
		_, err := file.Read(buffer)

		if err == io.EOF {
			break
		}
		if err != nil {
			os.Exit(2)
		}

		fmt.Printf("read: %s\n", string(buffer))
	}
}
