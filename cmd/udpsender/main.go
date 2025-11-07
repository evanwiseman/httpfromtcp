package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"os"
)

const (
	addr = "127.0.0.1:42069"
)

func main() {
	// Resolve the address
	raddr, err := net.ResolveUDPAddr("udp", addr)
	if err != nil {
		log.Fatalf("[FATAL] Failed to resolve address (%s): %v", addr, err)
	}
	log.Println("[INFO] Resolved address", addr)

	// Dial the address to prepare the UDP connection
	conn, err := net.DialUDP("udp", nil, raddr)
	if err != nil {
		log.Fatalf("[FATAL] Unable to prepare UDP connection: %v", err)
	}
	defer conn.Close()
	log.Println("[INFO] Successfully prepared UDP connection", addr)

	// Start REPL
	reader := bufio.NewReader(os.Stdin)
	for {
		// Get line from user
		fmt.Print("> ")
		line, err := reader.ReadString('\n')
		if err != nil {
			log.Printf("[ERROR] Unable to parse line: %v", err)
			continue
		}

		// Send line from user
		conn.Write([]byte(line))
	}
}
