package langsmith_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxWaitFailureAndTimeout(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/failed-box/status":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"status":         "failed",
				"status_message": "boot failed",
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/slow-box/status":
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "starting"})
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)

	_, err := client.Sandboxes.Boxes.Wait(context.Background(), "failed-box", langsmith.SandboxWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	var creationErr *langsmith.SandboxResourceCreationError
	if !errors.As(err, &creationErr) {
		t.Fatalf("expected SandboxResourceCreationError, got %T: %v", err, err)
	}
	if creationErr.Message != "boot failed" {
		t.Fatalf("unexpected creation error: %#v", creationErr)
	}

	_, err = client.Sandboxes.Boxes.Wait(context.Background(), "slow-box", langsmith.SandboxWaitParams{
		Timeout:      time.Millisecond,
		PollInterval: time.Millisecond,
	})
	var timeoutErr *langsmith.SandboxResourceTimeoutError
	if !errors.As(err, &timeoutErr) {
		t.Fatalf("expected SandboxResourceTimeoutError, got %T: %v", err, err)
	}
	if timeoutErr.LastStatus != "starting" {
		t.Fatalf("unexpected timeout error: %#v", timeoutErr)
	}
}

func TestSandboxStartAndWait(t *testing.T) {
	var statusCalls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/boxes/box-a/start":
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "starting"})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/box-a/status":
			statusCalls++
			status := "starting"
			if statusCalls > 1 {
				status = "ready"
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"status": status})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/box-a":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"name":          "box-a",
				"status":        "ready",
				"dataplane_url": "https://sandbox.example",
			})
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	sandbox, err := client.Sandboxes.Boxes.StartSandbox(context.Background(), "box-a", langsmith.SandboxWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("StartSandbox returned error: %v", err)
	}
	if sandbox.Name != "box-a" || sandbox.Status != "ready" {
		t.Fatalf("unexpected sandbox: %#v", sandbox)
	}
}
