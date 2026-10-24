package transport

import (
	"errors"
	"io"
	"net"
	"time"

	"github.com/gorilla/websocket"
)

// StreamSubprotocol identifies a WebSocket that carries one opaque byte
// stream to a connector target, as opposed to a connector carrier.
const StreamSubprotocol = "langsmith-stream.v1"

// WebSocketStream presents one WebSocket as a net.Conn carrying a single byte
// stream. A close frame is the stream's half-close: CloseWrite sends one and a
// normal close from the peer reads as io.EOF, while data still flows the
// other way until both sides have closed.
type WebSocketStream struct {
	*WebSocketConn
}

var _ net.Conn = (*WebSocketStream)(nil)

// NewWebSocketStream wraps conn as a byte stream.
func NewWebSocketStream(conn *websocket.Conn) *WebSocketStream {
	// Gorilla answers a close frame with its own by default, which would end
	// this side's sending as soon as the peer half-closed.
	conn.SetCloseHandler(func(int, string) error { return nil })
	return &WebSocketStream{WebSocketConn: NewWebSocketConn(conn)}
}

func (s *WebSocketStream) Read(p []byte) (int, error) {
	n, err := s.WebSocketConn.Read(p)
	if websocket.IsCloseError(err, websocket.CloseNormalClosure) {
		return n, io.EOF
	}
	return n, err
}

// Write splits p into messages no larger than the peer's read limit.
func (s *WebSocketStream) Write(p []byte) (int, error) {
	written := 0
	for written < len(p) {
		end := min(len(p), written+MaxWebSocketMessageSizeBytes)
		n, err := s.WebSocketConn.Write(p[written:end])
		written += n
		if err != nil {
			return written, err
		}
	}
	return written, nil
}

// CloseWrite tells the peer this side has finished sending.
func (s *WebSocketStream) CloseWrite() error {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	err := s.conn.WriteControl(
		websocket.CloseMessage,
		websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""),
		time.Now().Add(webSocketCloseTimeout),
	)
	if errors.Is(err, websocket.ErrCloseSent) {
		return nil
	}
	return err
}

func (s *WebSocketStream) LocalAddr() net.Addr  { return s.conn.LocalAddr() }
func (s *WebSocketStream) RemoteAddr() net.Addr { return s.conn.RemoteAddr() }

func (s *WebSocketStream) SetDeadline(t time.Time) error {
	return errors.Join(s.conn.SetReadDeadline(t), s.conn.SetWriteDeadline(t))
}

func (s *WebSocketStream) SetReadDeadline(t time.Time) error  { return s.conn.SetReadDeadline(t) }
func (s *WebSocketStream) SetWriteDeadline(t time.Time) error { return s.conn.SetWriteDeadline(t) }
