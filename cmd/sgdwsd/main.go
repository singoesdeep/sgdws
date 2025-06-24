package main

import (
	"log"
	"net"

	"github.com/singoesdeep/sgdws/internal/websocket"
)

var hub = websocket.NewHub()

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
	hub.AddClient(conn)
	defer hub.RemoveClient(conn)
	defer conn.Close()

	if err := websocket.PerformHandshake(conn.NetConn); err != nil {
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
			hub.Broadcast(conn, websocket.OpText, frame.Payload)

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
