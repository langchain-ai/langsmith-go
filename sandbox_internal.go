package langsmith

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"strings"
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
		ws, err := websocket.DialConfig(config)
		ch <- dialResult{ws: ws, err: err}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			return nil, &SandboxConnectionError{Message: fmt.Sprintf("langsmith: failed to connect to sandbox command WebSocket: %v", res.err)}
		}
		return res.ws, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
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
