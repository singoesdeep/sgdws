package websocket

import (
	"bufio"
	"crypto/sha1"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"net"
	"strings"
)

const websocketGUID = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"

func PerformHandshake(conn net.Conn) error {
	reader := bufio.NewReader(conn)

	requestLine, err := reader.ReadString('\n')
	if err != nil {
		return err
	}
	if !strings.HasPrefix(requestLine, "GET") {
		return errors.New("not a GET request")
	}

	headers := make(map[string]string)
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			return err
		}
		line = strings.TrimSpace(line)
		if line == "" {
			break
		}
		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}
		key := strings.ToLower(strings.TrimSpace(parts[0]))
		value := strings.TrimSpace(parts[1])
		headers[key] = value
	}

	if strings.ToLower(headers["connection"]) != "upgrade" {
		return errors.New("missing or invalid Connection header")
	}
	if strings.ToLower(headers["upgrade"]) != "websocket" {
		return errors.New("missing or invalid Upgrade header")
	}
	secKey, ok := headers["sec-websocket-key"]
	if !ok || secKey == "" {
		return errors.New("missing Sec-WebSocket-Key header")
	}
	secVersion, ok := headers["sec-websocket-version"]
	if !ok || secVersion != "13" {
		return errors.New("unsupported Sec-WebSocket-Version, expected 13")
	}

	h := sha1.New()
	h.Write([]byte(secKey + websocketGUID))
	acceptKey := base64.StdEncoding.EncodeToString(h.Sum(nil))

	response := fmt.Sprintf(
		"HTTP/1.1 101 Switching Protocols\r\n"+
			"Upgrade: websocket\r\n"+
			"Connection: Upgrade\r\n"+
			"Sec-WebSocket-Accept: %s\r\n\r\n",
		acceptKey)

	_, err = io.WriteString(conn, response)
	return err
}
