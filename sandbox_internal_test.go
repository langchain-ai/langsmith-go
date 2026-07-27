package langsmith

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxInternalFieldHelpers(t *testing.T) {
	if got, ok := sandboxRequiredString(String("echo ok")); !ok || got != "echo ok" {
		t.Fatalf("expected required string from value, got %q ok=%v", got, ok)
	}
	if got, ok := sandboxRequiredString(Raw[string]("raw command")); !ok || got != "raw command" {
		t.Fatalf("expected required string from raw value, got %q ok=%v", got, ok)
	}
	if got, ok := sandboxRequiredString(Null[string]()); ok || got != "" {
		t.Fatalf("expected null string to be missing, got %q ok=%v", got, ok)
	}
	if got := sandboxFieldValue(Int(7), int64(3)); got != 7 {
		t.Fatalf("expected explicit field value, got %d", got)
	}
	if got := sandboxFieldValue(Null[int64](), int64(3)); got != 3 {
		t.Fatalf("expected fallback value, got %d", got)
	}
}

func TestSandboxInternalURLHelpers(t *testing.T) {
	dataplaneURL, err := sandboxDataplaneURL("https://sandbox.example/base/", "/execute")
	if err != nil {
		t.Fatalf("sandboxDataplaneURL returned error: %v", err)
	}
	if dataplaneURL != "https://sandbox.example/base/execute" {
		t.Fatalf("unexpected dataplane URL: %q", dataplaneURL)
	}

	wsURL, err := sandboxWebSocketURL("https://sandbox.example/base/")
	if err != nil {
		t.Fatalf("sandboxWebSocketURL returned error: %v", err)
	}
	if wsURL != "wss://sandbox.example/base/execute/ws" {
		t.Fatalf("unexpected websocket URL: %q", wsURL)
	}

	origin, err := sandboxWebSocketOrigin("wss://sandbox.example/base/execute/ws")
	if err != nil {
		t.Fatalf("sandboxWebSocketOrigin returned error: %v", err)
	}
	if origin != "https://sandbox.example/" {
		t.Fatalf("unexpected websocket origin: %q", origin)
	}

	if _, err := sandboxDataplaneURL("sandbox.example", "execute"); err == nil {
		t.Fatal("expected invalid dataplane URL error")
	}
	if _, err := sandboxWebSocketURL("ftp://sandbox.example"); err == nil {
		t.Fatal("expected unsupported websocket scheme error")
	}
}

func TestRequireSandboxDataplaneURL(t *testing.T) {
	got, err := requireSandboxDataplaneURL("box-a", "ready", "https://sandbox.example")
	if err != nil {
		t.Fatalf("requireSandboxDataplaneURL returned error: %v", err)
	}
	if got != "https://sandbox.example" {
		t.Fatalf("unexpected dataplane URL: %q", got)
	}

	var notReady *SandboxNotReadyError
	if _, err := requireSandboxDataplaneURL("box-a", "starting", "https://sandbox.example"); !errors.As(err, &notReady) {
		t.Fatalf("expected SandboxNotReadyError, got %T: %v", err, err)
	}

	var notConfigured *SandboxDataplaneNotConfiguredError
	if _, err := requireSandboxDataplaneURL("box-a", "ready", ""); !errors.As(err, &notConfigured) {
		t.Fatalf("expected SandboxDataplaneNotConfiguredError, got %T: %v", err, err)
	}
}

func TestSandboxHeadersAndMinDuration(t *testing.T) {
	headers, err := sandboxHeaders(context.Background(), "https://sandbox.example/execute", option.WithAPIKey("test-api-key"))
	if err != nil {
		t.Fatalf("sandboxHeaders returned error: %v", err)
	}
	if got := headers.Get("X-API-Key"); got != "test-api-key" {
		t.Fatalf("expected API key header, got %q", got)
	}
	if got := minDuration(10*time.Millisecond, time.Second); got != 10*time.Millisecond {
		t.Fatalf("unexpected min duration: %s", got)
	}
}
