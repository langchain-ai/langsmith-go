package langsmith_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxFileMethodsResolveNamedSandbox(t *testing.T) {
	var uploaded bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/box-a":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"name":          "box-a",
				"status":        "ready",
				"dataplane_url": srvURL(r),
			})
		case r.Method == http.MethodGet && r.URL.Path == "/download":
			if got := r.URL.Query().Get("path"); got != "/tmp/read.txt" {
				t.Fatalf("unexpected download path: %q", got)
			}
			_, _ = w.Write([]byte("downloaded"))
		case r.Method == http.MethodPost && r.URL.Path == "/upload":
			if got := r.URL.Query().Get("path"); got != "/tmp/write.txt" {
				t.Fatalf("unexpected upload path: %q", got)
			}
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("read upload file: %v", err)
			}
			defer file.Close()
			content, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("read upload content: %v", err)
			}
			if string(content) != "uploaded" {
				t.Fatalf("unexpected upload content: %q", string(content))
			}
			uploaded = true
			w.WriteHeader(http.StatusNoContent)
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

	content, err := client.Sandboxes.Boxes.ReadFile(context.Background(), "box-a", "/tmp/read.txt")
	if err != nil {
		t.Fatalf("ReadFile returned error: %v", err)
	}
	if string(content) != "downloaded" {
		t.Fatalf("unexpected read content: %q", string(content))
	}

	if err := client.Sandboxes.Boxes.WriteFile(context.Background(), "box-a", "/tmp/write.txt", []byte("uploaded")); err != nil {
		t.Fatalf("WriteFile returned error: %v", err)
	}
	if !uploaded {
		t.Fatal("expected upload request")
	}
}
