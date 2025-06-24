package websocket

import (
	"sync"
)

type Hub struct {
	mu      sync.RWMutex
	clients map[*Conn]bool
}

func NewHub() *Hub {
	return &Hub{
		clients: make(map[*Conn]bool),
	}
}

func (h *Hub) AddClient(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	h.clients[c] = true
}

func (h *Hub) RemoveClient(c *Conn) {
	h.mu.Lock()
	defer h.mu.Unlock()
	delete(h.clients, c)
}

func (h *Hub) Broadcast(from *Conn, opcode byte, payload []byte) {
	h.mu.RLock()
	defer h.mu.RUnlock()

	for client := range h.clients {
		if client != from {
			client.WriteFrame(opcode, payload)
		}
	}
}
