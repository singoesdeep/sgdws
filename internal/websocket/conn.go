package websocket

import (
	"net"
	"sync"
)

type Conn struct {
	NetConn net.Conn
	writeMu sync.Mutex
}

func NewConn(nc net.Conn) *Conn {
	return &Conn{
		NetConn: nc,
	}
}

func (c *Conn) ReadFrame() (*Frame, error) {
	return ReadFrame(c.NetConn)
}

func (c *Conn) WriteFrame(opcode byte, payload []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()

	return WriteFrame(c.NetConn, opcode, payload)
}

func (c *Conn) Close() error {
	return c.NetConn.Close()
}
