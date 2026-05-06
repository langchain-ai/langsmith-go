package langsmith_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxFileReadWriteWithDataplaneURL(t *testing.T) {
	var uploadedPath string
	var uploadedContent string
	var uploadAPIKey string
	var readAPIKey string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/upload":
			uploadAPIKey = r.Header.Get("X-API-Key")
			uploadedPath = r.URL.Query().Get("path")
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("read upload file: %v", err)
			}
			defer file.Close()
			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("read upload content: %v", err)
			}
			uploadedContent = string(data)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/download":
			readAPIKey = r.Header.Get("X-API-Key")
			if r.URL.Query().Get("path") != "/app/test.txt" {
				t.Fatalf("unexpected read path: %q", r.URL.Query().Get("path"))
			}
			_, _ = w.Write([]byte("file contents"))
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	if err := client.Sandboxes.Boxes.WriteFileWithDataplaneURL(context.Background(), srv.URL, "/app/test.txt", []byte("hello world")); err != nil {
		t.Fatalf("WriteFileWithDataplaneURL returned error: %v", err)
	}
	if uploadedPath != "/app/test.txt" {
		t.Fatalf("unexpected upload path: %q", uploadedPath)
	}
	if uploadedContent != "hello world" {
		t.Fatalf("unexpected upload content: %q", uploadedContent)
	}
	if uploadAPIKey != "test-api-key" {
		t.Fatalf("expected upload API key header, got %q", uploadAPIKey)
	}

	content, err := client.Sandboxes.Boxes.ReadFileWithDataplaneURL(context.Background(), srv.URL, "/app/test.txt")
	if err != nil {
		t.Fatalf("ReadFileWithDataplaneURL returned error: %v", err)
	}
	if string(content) != "file contents" {
		t.Fatalf("unexpected read content: %q", string(content))
	}
	if readAPIKey != "test-api-key" {
		t.Fatalf("expected read API key header, got %q", readAPIKey)
	}
}

func TestSandboxWaitHelpers(t *testing.T) {
	var sandboxStatusCalls atomic.Int64
	var snapshotGetCalls atomic.Int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/test-box/status":
			if sandboxStatusCalls.Add(1) == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{"status": "provisioning"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "ready"})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/test-box":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":            "box-id",
				"name":          "test-box",
				"status":        "ready",
				"dataplane_url": "http://dataplane.test",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/snapshots":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "snap-1",
				"name":              "snap",
				"status":            "building",
				"fs_capacity_bytes": 1024,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/snapshots/snap-1":
			if snapshotGetCalls.Add(1) == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id":                "snap-1",
					"name":              "snap",
					"status":            "building",
					"fs_capacity_bytes": 1024,
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "snap-1",
				"name":              "snap",
				"status":            "ready",
				"fs_capacity_bytes": 1024,
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

	sandbox, err := client.Sandboxes.Boxes.WaitSandbox(context.Background(), "test-box", langsmith.SandboxWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("WaitSandbox returned error: %v", err)
	}
	if sandbox.Name != "test-box" || sandbox.DataplaneURL != "http://dataplane.test" {
		t.Fatalf("unexpected sandbox: %#v", sandbox)
	}
	if sandboxStatusCalls.Load() != 2 {
		t.Fatalf("expected 2 sandbox status calls, got %d", sandboxStatusCalls.Load())
	}

	snapshot, err := client.Sandboxes.Snapshots.NewAndWait(context.Background(), langsmith.SandboxSnapshotNewParams{
		Name:            langsmith.String("snap"),
		DockerImage:     langsmith.String("python:3.12-slim"),
		FsCapacityBytes: langsmith.Int(1024),
	}, langsmith.SnapshotWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewAndWait returned error: %v", err)
	}
	if snapshot.ID != "snap-1" || snapshot.Status != "ready" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	if snapshotGetCalls.Load() != 2 {
		t.Fatalf("expected 2 snapshot get calls, got %d", snapshotGetCalls.Load())
	}
}

func TestSandboxServiceURLRefreshAndRequest(t *testing.T) {
	var serviceURLCalls atomic.Int64
	var firstPayload map[string]any
	var serviceAuthHeader string
	var serviceCustomHeader string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/boxes/test-box/service-url":
			call := serviceURLCalls.Add(1)
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode service URL request: %v", err)
			}
			if call == 1 {
				firstPayload = payload
				_ = json.NewEncoder(w).Encode(map[string]any{
					"browser_url": srvURL(r) + "/browser",
					"service_url": srvURL(r) + "/service",
					"token":       "token-1",
					"expires_at":  time.Now().Add(time.Second).UTC().Format(time.RFC3339),
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"browser_url": srvURL(r) + "/browser",
				"service_url": srvURL(r) + "/service",
				"token":       "token-2",
				"expires_at":  time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			})
		case r.Method == http.MethodGet && r.URL.Path == "/service/path":
			serviceAuthHeader = r.Header.Get("X-Langsmith-Sandbox-Service-Token")
			serviceCustomHeader = r.Header.Get("X-Test")
			_, _ = w.Write([]byte(`{"ok":true}`))
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
	service, err := client.Sandboxes.Boxes.Service(context.Background(), "test-box", langsmith.SandboxBoxGenerateServiceURLParams{
		Port: langsmith.Int(3000),
	})
	if err != nil {
		t.Fatalf("Service returned error: %v", err)
	}
	if firstPayload["port"] != float64(3000) {
		t.Fatalf("unexpected service port payload: %#v", firstPayload["port"])
	}
	if firstPayload["expires_in_seconds"] != float64(600) {
		t.Fatalf("expected default service URL TTL, got %#v", firstPayload["expires_in_seconds"])
	}

	token, err := service.Token(context.Background())
	if err != nil {
		t.Fatalf("Token returned error: %v", err)
	}
	if token != "token-2" {
		t.Fatalf("expected refreshed token, got %q", token)
	}
	if serviceURLCalls.Load() != 2 {
		t.Fatalf("expected 2 service URL calls, got %d", serviceURLCalls.Load())
	}

	resp, err := service.Get(context.Background(), "/path", http.Header{"X-Test": []string{"yes"}})
	if err != nil {
		t.Fatalf("service Get returned error: %v", err)
	}
	_ = resp.Body.Close()
	if serviceAuthHeader != "token-2" {
		t.Fatalf("expected service auth token header, got %q", serviceAuthHeader)
	}
	if serviceCustomHeader != "yes" {
		t.Fatalf("expected custom service header, got %q", serviceCustomHeader)
	}
}
