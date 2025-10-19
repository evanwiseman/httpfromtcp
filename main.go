package main

import (
	"fmt"
	"io"
	"log"
	"net"
	"strings"
)

const (
	addr = ":42069"
)

func getLinesChannel(f io.ReadCloser) <-chan string {
	ch := make(chan string)
	log.Printf("[INFO] Opened read channel")

	go func() {
		defer func() {
			close(ch)
			log.Printf("[INFO] Closed read channel")
		}()
		var line string
		for {
			buffer := make([]byte, 8)

			// Read file into buffer until EOF
			n, err := f.Read(buffer)
			if err == io.EOF {
				break
			}
			if err != nil {
				log.Printf("[ERROR] Error reading from stream: %v", err)
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
	listener, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("[FATAL] Failed to open listener: %v", err)
	}
	defer listener.Close()
	log.Printf("[INFO] Server started on: %s", addr)

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Fatalf("[ERROR] Unable to accept connection: %v", err)
		}
		log.Printf("[INFO] Accepted new connection from: %s", conn.RemoteAddr().String())

		// Read all lines from the channel
		linesCh := getLinesChannel(conn)
		for line := range linesCh {
			fmt.Printf("%s\n", line)
		}
	}

}
