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

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)

	go func() {
		defer close(ch)
		var line string
		for {
			buffer := make([]byte, 8)

			// Read file into buffer until EOF
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				return
			}

			// Send to ch until last new line
			parts := strings.Split(string(buffer[:n]), "\n")
			line += parts[0]
			for _, part := range parts[1:] {
				ch <- line
				line = part
			}
		}
		// Ensure everything is read
		if line != "" {
			ch <- line
		}
	}()

	return ch
}

func main() {
	// Open the file
	f, err := os.Open(filename)
	if err != nil {
		os.Exit(1)
	}
	defer f.Close()

	// Read all lines from the channel
	linesCh := getLinesChannel(f)
	for line := range linesCh {
		fmt.Printf("read: %s\n", line)
	}
}
