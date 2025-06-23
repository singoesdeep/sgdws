package main

import (
	"log"
	"net"
)

func main() {
	ln, err := net.Listen("tcp", ":8080")
	if err != nil {
		log.Fatalf("Failed to start listener: %v", err)
	}
	log.Println("Server listening on port 8080...")

	for {
		conn, err := ln.Accept()
		if err != nil {
			log.Printf("Connection error: %v", err)
			continue
		}
		go handleConn(conn)
	}
}

func handleConn(conn net.Conn) {
	defer conn.Close()

	buf := make([]byte, 1024)
	n, err := conn.Read(buf)
	if err != nil {
		log.Println("Read error:", err)
		return
	}
	log.Printf("🔹 Received data: %s", string(buf[:n]))
}
