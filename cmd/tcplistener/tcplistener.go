package main

import (
	"log"
	"net"

	"github.com/evanwiseman/httpfromtcp/internal/request"
)

const (
	addr = ":42069"
)

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

		req, err := request.RequestFromReader(conn)
		if err != nil {
			log.Fatalf("[ERROR] Unable to read data: %v", err)
		}
		request.PrintRequest(req)
	}

}
