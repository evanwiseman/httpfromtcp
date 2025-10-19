package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

const (
	filename = "messages.txt"
)

func main() {
	file, err := os.Open(filename)
	if err != nil {
		os.Exit(1)
	}

	// Store the current line to output
	var currentLine string
	for {
		// Byte slice of size 8
		buffer := make([]byte, 8)

		// Read file into buffer
		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}
		if err != nil {
			os.Exit(2)
		}

		// Split buffer into parts and add first part to current line
		parts := strings.Split(string(buffer[:n]), "\n")

		// Append
		currentLine += parts[0]

		// iterate through parts
		for i, part := range parts {
			if i == 0 {
				continue
			}
			fmt.Printf("read: %s\n", currentLine) // print current line
			currentLine = ""                      // reset current line
			currentLine += part                   // add to current line
		}
	}

	if currentLine != "" { // ensure we've read everything
		fmt.Printf("read: %s\n", currentLine)
	}
}
