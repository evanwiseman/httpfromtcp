package main

import (
	"fmt"
	"io"
	"log"
	"net"
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
	listener, err := net.Listen("tcp", "localhost:42069")
	if err != nil {
		log.Fatalf("Failed to open listener: %v", err)
	}
	defer listener.Close()
	log.Println("Successfully opened listener")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("Unable to accept connection: %v", err)
		}
		log.Println("Successfully accepted connection")

		// Read all lines from the channel
		linesCh := getLinesChannel(conn)
		for line := range linesCh {
			fmt.Printf("%s\n", line)
		}
		<-linesCh
		log.Println("Successfully closed channel")
	}

}
