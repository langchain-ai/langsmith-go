package transport

import (
	"errors"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/gorilla/websocket"
)

const (
	MaxWebSocketMessageSizeBytes = 256 * 1024
	webSocketCloseTimeout        = time.Second
)

var (
	ErrUnexpectedWebSocketMessageType = errors.New("unexpected websocket message type")
	ErrWebSocketMessageTooLarge       = errors.New("websocket message too large")
)

// WebSocketConn adapts Gorilla's message transport to the byte stream yamux
// expects. Yamux has one reader but may write from multiple goroutines.
type WebSocketConn struct {
	conn *websocket.Conn

	reader  io.Reader
	writeMu sync.Mutex

	closeOnce sync.Once
	closeErr  error
}

var _ io.ReadWriteCloser = (*WebSocketConn)(nil)

func NewWebSocketConn(conn *websocket.Conn) *WebSocketConn {
	conn.SetReadLimit(MaxWebSocketMessageSizeBytes)
	return &WebSocketConn{conn: conn}
}

func (c *WebSocketConn) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}

	for {
		if c.reader == nil {
			messageType, reader, err := c.conn.NextReader()
			if err != nil {
				return 0, err
			}
			if messageType != websocket.BinaryMessage {
				return 0, fmt.Errorf("%w: %d", ErrUnexpectedWebSocketMessageType, messageType)
			}
			c.reader = reader
		}

		n, err := c.reader.Read(p)
		if errors.Is(err, io.EOF) {
			c.reader = nil
			if n > 0 {
				return n, nil
			}
			continue
		}
		return n, err
	}
}

func (c *WebSocketConn) Write(p []byte) (int, error) {
	if len(p) > MaxWebSocketMessageSizeBytes {
		return 0, fmt.Errorf("%w: %d bytes", ErrWebSocketMessageTooLarge, len(p))
	}

	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	if err := c.conn.WriteMessage(websocket.BinaryMessage, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (c *WebSocketConn) Close() error {
	c.closeOnce.Do(func() {
		_ = c.conn.WriteControl(
			websocket.CloseMessage,
			websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
			time.Now().Add(webSocketCloseTimeout),
		)
		c.closeErr = c.conn.Close()
	})
	return c.closeErr
}
