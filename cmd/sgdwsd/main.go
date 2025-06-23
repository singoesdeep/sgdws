package main

import (
	"log"
	"net"

	"github.com/singoesdeep/sgdws/internal/websocket"
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

func handleConn(nc net.Conn) {
	conn := websocket.NewConn(nc)
	defer conn.Close()

	err := websocket.PerformHandshake(conn.NetConn)
	if err != nil {
		log.Println("Handshake failed:", err)
		return
	}

	for {
		frame, err := conn.ReadFrame()
		if err != nil {
			log.Println("ReadFrame error:", err)
			break
		}

		switch frame.Opcode {
		case websocket.OpText:
			log.Printf("Text message: %s", string(frame.Payload))
			err := conn.WriteFrame(websocket.OpText, frame.Payload)
			if err != nil {
				log.Println("WriteFrame error:", err)
				return
			}

		case websocket.OpClose:
			log.Println("Close frame received")
			return

		case websocket.OpPing:
			log.Println("Ping frame received - ponging back")
			err := conn.WriteFrame(websocket.OpPong, nil)
			if err != nil {
				log.Println("WriteFrame pong error:", err)
				return
			}

		case websocket.OpPong:
			log.Println("Pong frame received")

		default:
			log.Printf("Unhandled opcode: %d", frame.Opcode)
		}
	}
}
