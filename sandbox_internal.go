package langsmith

import (
	"bufio"
	"bytes"
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"golang.org/x/net/websocket"
)

func sandboxRequiredString(field param.Field[string]) (string, bool) {
	if !field.Present || field.Null {
		return "", false
	}
	if field.Value != "" {
		return field.Value, true
	}
	if raw, ok := field.Raw.(string); ok && raw != "" {
		return raw, true
	}
	return "", false
}

func sandboxFieldValue[T any](field param.Field[T], fallback T) T {
	if field.Present && !field.Null {
		return field.Value
	}
	return fallback
}

func requireSandboxDataplaneURL(name string, status string, dataplaneURL string) (string, error) {
	if status != "" && status != "ready" {
		return "", &SandboxNotReadyError{SandboxName: name, Status: status}
	}
	if dataplaneURL == "" {
		return "", &SandboxDataplaneNotConfiguredError{SandboxName: name}
	}
	return dataplaneURL, nil
}

func sandboxDataplaneURL(dataplaneURL string, path string) (string, error) {
	u, err := url.Parse(dataplaneURL)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid sandbox dataplane URL %q", dataplaneURL)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return u.String(), nil
}

func sandboxWebSocketURL(dataplaneURL string) (string, error) {
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
	u.Path = strings.TrimRight(u.Path, "/") + "/execute/ws"
	u.RawQuery = ""
	return u.String(), nil
}

func dialSandboxCommandWebSocket(ctx context.Context, dataplaneURL string, opts ...option.RequestOption) (*websocket.Conn, error) {
	wsURL, err := sandboxWebSocketURL(dataplaneURL)
	if err != nil {
		return nil, err
	}
	return dialSandboxWebSocketURL(ctx, wsURL, opts...)
}

func dialSandboxWebSocketURL(ctx context.Context, wsURL string, opts ...option.RequestOption) (*websocket.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	headers, err := sandboxHeaders(ctx, wsURL, opts...)
	if err != nil {
		return nil, err
	}
	origin, err := sandboxWebSocketOrigin(wsURL)
	if err != nil {
		return nil, err
	}
	config, err := websocket.NewConfig(wsURL, origin)
	if err != nil {
		return nil, err
	}
	config.Header = headers

	type dialResult struct {
		ws  *websocket.Conn
		err error
	}
	ch := make(chan dialResult, 1)
	go func() {
		ws, err := dialSandboxWebSocketConfig(ctx, config)
		ch <- dialResult{ws: ws, err: err}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			return nil, res.err
		}
		return res.ws, nil
	case <-ctx.Done():
		go func() {
			if res := <-ch; res.ws != nil {
				_ = res.ws.Close()
			}
		}()
		return nil, ctx.Err()
	}
}

// dialSandboxWebSocketConfig handshakes over a connection that records the
// server's reply, so a refused upgrade reports its HTTP status and body rather
// than websocket.ErrBadStatus, which drops both.
func dialSandboxWebSocketConfig(ctx context.Context, config *websocket.Config) (*websocket.Conn, error) {
	conn, err := dialSandboxWebSocketTransport(ctx, config)
	if err != nil {
		return nil, &SandboxConnectionError{Message: sandboxWebSocketDialErrorf("%v", err)}
	}
	recorder := &sandboxHandshakeRecorder{Conn: conn}
	ws, err := websocket.NewClient(config, recorder)
	if err != nil {
		_ = conn.Close()
		status, detail := recorder.handshakeResponse()
		if detail == "" {
			return nil, &SandboxConnectionError{Message: sandboxWebSocketDialErrorf("%v", err)}
		}
		return nil, &SandboxConnectionError{
			Message:    sandboxWebSocketDialErrorf("%s", detail),
			StatusCode: status,
		}
	}
	recorder.stop()
	return ws, nil
}

func sandboxWebSocketDialErrorf(format string, args ...any) string {
	return "langsmith: failed to connect to sandbox command WebSocket: " + fmt.Sprintf(format, args...)
}

func dialSandboxWebSocketTransport(ctx context.Context, config *websocket.Config) (net.Conn, error) {
	if config.Location == nil {
		return nil, errors.New("missing WebSocket URL")
	}
	dialer := config.Dialer
	if dialer == nil {
		dialer = &net.Dialer{}
	}
	host := config.Location.Host
	switch config.Location.Scheme {
	case "ws":
		if config.Location.Port() == "" {
			host = net.JoinHostPort(host, "80")
		}
		return dialer.DialContext(ctx, "tcp", host)
	case "wss":
		if config.Location.Port() == "" {
			host = net.JoinHostPort(host, "443")
		}
		return (&tls.Dialer{NetDialer: dialer, Config: config.TlsConfig}).DialContext(ctx, "tcp", host)
	default:
		return nil, fmt.Errorf("unsupported WebSocket URL scheme %q", config.Location.Scheme)
	}
}

const (
	sandboxHandshakeRecordLimitBytes = 8 << 10
	sandboxHandshakeDetailMaxBytes   = 512
)

type sandboxHandshakeRecorder struct {
	net.Conn

	mu      sync.Mutex
	buf     bytes.Buffer
	stopped bool
}

func (r *sandboxHandshakeRecorder) Read(p []byte) (int, error) {
	n, err := r.Conn.Read(p)
	if n > 0 {
		r.mu.Lock()
		if room := sandboxHandshakeRecordLimitBytes - r.buf.Len(); !r.stopped && room > 0 {
			r.buf.Write(p[:min(n, room)])
		}
		r.mu.Unlock()
	}
	return n, err
}

func (r *sandboxHandshakeRecorder) stop() {
	r.mu.Lock()
	defer r.mu.Unlock()
	r.stopped = true
	r.buf.Reset()
}

// handshakeResponse returns the status code and a "<status>: <body>" summary of
// the recorded reply, or (0, "") when it was not a readable HTTP response.
func (r *sandboxHandshakeRecorder) handshakeResponse() (int, string) {
	r.mu.Lock()
	raw := bytes.Clone(r.buf.Bytes())
	r.mu.Unlock()

	resp, err := http.ReadResponse(bufio.NewReader(bytes.NewReader(raw)), nil)
	if err != nil {
		return 0, ""
	}
	defer resp.Body.Close()

	body, _ := io.ReadAll(io.LimitReader(resp.Body, sandboxHandshakeDetailMaxBytes))
	detail := strings.Join(strings.Fields(string(body)), " ")
	if detail == "" {
		return resp.StatusCode, resp.Status
	}
	return resp.StatusCode, resp.Status + ": " + detail
}

func sandboxWebSocketOrigin(wsURL string) (string, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	u.Path = "/"
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func sandboxHeaders(ctx context.Context, requestURL string, opts ...option.RequestOption) (http.Header, error) {
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, requestURL, nil, nil, opts...)
	if err != nil {
		return nil, err
	}
	return cfg.Request.Header.Clone(), nil
}

func minDuration(a time.Duration, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
