# sgdws

`sgdws` is a minimal, from-scratch implementation of the WebSocket protocol written in pure Go with `net` package.  
It does not rely on `net/http` or any WebSocket libraries — all handshake, framing, and message handling is implemented manually.


## Goals

- Implement RFC 6455-compliant WebSocket handshake
- Build frame encoder/decoder (text, binary, ping, pong, close)
- Support raw TCP WebSocket connections (no Gorilla, no stdlib HTTP)
- Write clean, testable and idiomatic Go code
- Learn low-level networking and protocol design

## Why?

This project is a learning-driven exploration of how WebSocket really works under the hood.  
Instead of using high-level libraries, `sgdws` gets its hands dirty with:

- Manual header parsing
- Base64 & SHA-1 handling
- Byte-level frame parsing & masking
- Connection lifecycle and concurrency

## Features

- WebSocket handshake over raw TCP
- RFC-compliant frame parsing (FIN, opcode, payload, mask)
- Frame writing (echo server works!)
- Ping / Pong frame support
- Close frame support
- Multi-client management (planned)
- Broadcast & hub (planned)

## Usage

Run the example echo server:

```bash
go run ./cmd/sgdws
```

Then connect using any WebSocket client:

```bash
websocat ws://localhost:8080
```

Send a message and get it echoed back!

Project Structure
```cs
sgdws/
├── cmd/sgdws/            # Example server
│   └── main.go
├── internal/websocket/    # Core protocol implementation
│   ├── handshake.go
│   ├── frame.go
│   └── conn.go
└── README.md
```

### Inspired by
RFC 6455
Wireshark
gorilla/websocket (but we don't use it)
Curiosity 😄

