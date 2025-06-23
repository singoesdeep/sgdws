package websocket

import (
	"encoding/binary"
	"io"
	"net"
)

const (
	OpContinuation = 0x0
	OpText         = 0x1
	OpBinary       = 0x2
	OpClose        = 0x8
	OpPing         = 0x9
	OpPong         = 0xA
)

type Frame struct {
	Fin     bool
	Opcode  byte
	Masked  bool
	Payload []byte
}

// reads a single WebSocket frame from conn.
func ReadFrame(conn net.Conn) (*Frame, error) {
	header := make([]byte, 2)
	if _, err := io.ReadFull(conn, header); err != nil {
		return nil, err
	}

	fin := (header[0] & 0x80) != 0
	opcode := header[0] & 0x0F
	masked := (header[1] & 0x80) != 0
	payloadLen := int(header[1] & 0x7F)

	// Read extended payload length if needed
	switch payloadLen {
	case 126:
		ext := make([]byte, 2)
		if _, err := io.ReadFull(conn, ext); err != nil {
			return nil, err
		}
		payloadLen = int(binary.BigEndian.Uint16(ext))
	case 127:
		ext := make([]byte, 8)
		if _, err := io.ReadFull(conn, ext); err != nil {
			return nil, err
		}

		payloadLen = int(binary.BigEndian.Uint64(ext))
	}

	var maskingKey []byte
	if masked {
		maskingKey = make([]byte, 4)
		if _, err := io.ReadFull(conn, maskingKey); err != nil {
			return nil, err
		}
	}

	payload := make([]byte, payloadLen)
	if _, err := io.ReadFull(conn, payload); err != nil {
		return nil, err
	}

	if masked {
		for i := 0; i < payloadLen; i++ {
			payload[i] ^= maskingKey[i%4]
		}
	}

	return &Frame{
		Fin:     fin,
		Opcode:  opcode,
		Masked:  masked,
		Payload: payload,
	}, nil
}

func WriteFrame(conn net.Conn, opcode byte, payload []byte) error {
	finAndOpcode := byte(0x80) | (opcode & 0x0F)
	payloadLen := len(payload)

	var header []byte

	switch {
	case payloadLen <= 125:
		header = make([]byte, 2)
		header[0] = finAndOpcode
		header[1] = byte(payloadLen)
	case payloadLen <= 65535:
		header = make([]byte, 4)
		header[0] = finAndOpcode
		header[1] = 126
		binary.BigEndian.PutUint16(header[2:], uint16(payloadLen))
	default:
		header = make([]byte, 10)
		header[0] = finAndOpcode
		header[1] = 127
		binary.BigEndian.PutUint64(header[2:], uint64(payloadLen))
	}

	if _, err := conn.Write(header); err != nil {
		return err
	}
	if _, err := conn.Write(payload); err != nil {
		return err
	}

	return nil
}
