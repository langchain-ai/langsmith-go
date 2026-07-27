package langsmith

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxCommandWrapperRunUsesCurrentDataplaneURL(t *testing.T) {
	var gotPayload map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/execute" {
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"stdout":    "ok",
			"stderr":    "",
			"exit_code": 0,
		})
	}))
	defer srv.Close()

	sandbox := &Sandbox{
		Name:         "box-a",
		Status:       "ready",
		DataplaneURL: srv.URL,
		boxes: NewSandboxBoxService(
			option.WithBaseURL("http://control-plane.test"),
			option.WithAPIKey("test-api-key"),
			option.WithMaxRetries(0),
		),
	}

	result, err := sandbox.Run(context.Background(), SandboxBoxRunParams{Command: String("echo ok")})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if result.Stdout != "ok" || !result.Success() {
		t.Fatalf("unexpected result: %#v", result)
	}
	if gotPayload["command"] != "echo ok" {
		t.Fatalf("unexpected command payload: %#v", gotPayload["command"])
	}
}

func TestSandboxCommandWrapperDoesNotGateOnStatus(t *testing.T) {
	// A non-ready box is NOT rejected client-side: the platform resumes a
	// stopped box when the dataplane request arrives. (The call still errors
	// here only because the test client has no configured base URL.)
	sandbox := &Sandbox{
		Name:         "box-a",
		Status:       "starting",
		DataplaneURL: "https://sandbox.example",
		boxes:        NewSandboxBoxService(),
	}
	_, err := sandbox.Run(context.Background(), SandboxBoxRunParams{Command: String("echo ok")})
	var notReady *SandboxNotReadyError
	if errors.As(err, &notReady) {
		t.Fatalf("client must not gate on status, but got SandboxNotReadyError: %v", err)
	}

	// A missing dataplane URL is still rejected client-side, regardless of status.
	noURL := &Sandbox{Name: "box-a", Status: "starting", boxes: NewSandboxBoxService()}
	_, err = noURL.Run(context.Background(), SandboxBoxRunParams{Command: String("echo ok")})
	var notConfigured *SandboxDataplaneNotConfiguredError
	if !errors.As(err, &notConfigured) {
		t.Fatalf("expected SandboxDataplaneNotConfiguredError, got %T: %v", err, err)
	}
}
