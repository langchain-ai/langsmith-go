package langsmith

import (
	"fmt"
	"net/url"
	"strings"
	"sync"

	"golang.org/x/net/websocket"
)

func sandboxTunnelWebSocketURL(dataplaneURL string) (string, error) {
	u, err := url.Parse(dataplaneURL)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported sandbox dataplane URL scheme %q", u.Scheme)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/tunnel"
	u.RawQuery = ""
	return u.String(), nil
}

type webSocketByteStream struct {
	ws      *websocket.Conn
	readMu  sync.Mutex
	writeMu sync.Mutex
	buf     []byte
}

func newWebSocketByteStream(ws *websocket.Conn) *webSocketByteStream {
	return &webSocketByteStream{ws: ws}
}

func (s *webSocketByteStream) Read(p []byte) (int, error) {
	s.readMu.Lock()
	defer s.readMu.Unlock()
	for len(s.buf) == 0 {
		var msg []byte
		if err := websocket.Message.Receive(s.ws, &msg); err != nil {
			return 0, err
		}
		s.buf = msg
	}
	n := copy(p, s.buf)
	s.buf = s.buf[n:]
	return n, nil
}

func (s *webSocketByteStream) Write(p []byte) (int, error) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if err := websocket.Message.Send(s.ws, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *webSocketByteStream) Close() error {
	return s.ws.Close()
}
