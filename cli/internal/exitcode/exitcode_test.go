package exitcode

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
)

func TestFromHTTPStatus(t *testing.T) {
	tests := map[int]int{
		200: OK,
		400: Input,
		401: Auth,
		402: Conflict,
		403: Auth,
		404: NotFound,
		406: Input,
		408: Unavailable,
		409: Conflict,
		410: NotFound,
		422: Input,
		429: RateLimited,
		500: Unavailable,
		503: Unavailable,
	}
	for status, want := range tests {
		if got := FromHTTPStatus(status); got != want {
			t.Errorf("FromHTTPStatus(%d) = %d, want %d", status, got, want)
		}
	}
}

func TestFor(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want int
	}{
		{"nil", nil, OK},
		{"plain error", errors.New("--limit must be positive"), Input},
		{"SDK API error", fmt.Errorf("listing projects: %w", &langsmith.Error{StatusCode: 404}), NotFound},
		{"SDK rate limit", &langsmith.Error{StatusCode: 429}, RateLimited},
		{"SDK server error", &langsmith.Error{StatusCode: 502}, Unavailable},
		{"lookup with no match", fmt.Errorf("project %w: demo", client.ErrNotFound), NotFound},
		{"network failure", &url.Error{Op: "Get", URL: "https://example.invalid", Err: errors.New("connection refused")}, Unavailable},
		{"canceled request", &url.Error{Op: "Get", URL: "https://example.invalid", Err: context.Canceled}, Input},
		{"invalid URL", &url.Error{Op: "parse", URL: "://bad", Err: errors.New("missing protocol scheme")}, Input},
		{"request timeout", &url.Error{Op: "Get", URL: "https://example.invalid", Err: context.DeadlineExceeded}, Unavailable},
		{"explicit code", WithCode(errors.New("timed out waiting"), Timeout), Timeout},
		{"explicit code wins over the status", WithCode(&langsmith.Error{StatusCode: 404}, Auth), Auth},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := For(tt.err); got != tt.want {
				t.Errorf("For(%v) = %d, want %d", tt.err, got, tt.want)
			}
		})
	}
}

func TestWithCodeKeepsMessageAndChain(t *testing.T) {
	err := WithCode(fmt.Errorf("wrapped: %w", client.ErrNotFound), Timeout)
	if err.Error() != "wrapped: not found" {
		t.Errorf("message = %q", err.Error())
	}
	if !errors.Is(err, client.ErrNotFound) {
		t.Error("WithCode hid the wrapped error")
	}
	if WithCode(nil, Timeout) != nil {
		t.Error("WithCode(nil) is not nil")
	}
}

func TestForRawClientErrors(t *testing.T) {
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusConflict)
		_, _ = w.Write([]byte(`{"detail":"already exists"}`))
	}))
	defer server.Close()
	c := client.NewWithOptions(client.Options{APIURL: server.URL, APIKey: "test-key"})
	err := c.RawGet(context.Background(), "/api/v1/datasets", nil)
	if got := For(err); got != Conflict {
		t.Errorf("For(%v) = %d, want %d", err, got, Conflict)
	}
}
