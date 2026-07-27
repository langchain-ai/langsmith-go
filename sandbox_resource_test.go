package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxResourceWrappers(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/boxes":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "new-id",
				"name":              "new-box",
				"status":            "starting",
				"dataplane_url":     "https://sandbox.example/new",
				"ttl_seconds":       600,
				"fs_capacity_bytes": 1024,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/new-box":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "get-id",
				"name":              "new-box",
				"status":            "ready",
				"dataplane_url":     "https://sandbox.example/get",
				"ttl_seconds":       700,
				"fs_capacity_bytes": 2048,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"offset": 0,
				"sandboxes": []map[string]any{
					{
						"id":                "list-id",
						"name":              "listed-box",
						"status":            "ready",
						"dataplane_url":     "https://sandbox.example/list",
						"ttl_seconds":       800,
						"fs_capacity_bytes": 4096,
					},
				},
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

	created, err := client.Sandboxes.Boxes.NewSandbox(context.Background(), langsmith.SandboxBoxNewParams{
		SnapshotID: langsmith.String("snapshot-id"),
		Name:       langsmith.String("new-box"),
	})
	if err != nil {
		t.Fatalf("NewSandbox returned error: %v", err)
	}
	if created.Name != "new-box" || created.DataplaneURL != "https://sandbox.example/new" || created.TTLSeconds != 600 {
		t.Fatalf("unexpected created sandbox: %#v", created)
	}

	got, err := client.Sandboxes.Boxes.GetSandbox(context.Background(), "new-box")
	if err != nil {
		t.Fatalf("GetSandbox returned error: %v", err)
	}
	if got.Name != "new-box" || got.Status != "ready" || got.FsCapacityBytes != 2048 {
		t.Fatalf("unexpected fetched sandbox: %#v", got)
	}

	listed, err := client.Sandboxes.Boxes.ListSandboxes(context.Background(), langsmith.SandboxBoxListParams{})
	if err != nil {
		t.Fatalf("ListSandboxes returned error: %v", err)
	}
	if len(listed) != 1 || listed[0].Name != "listed-box" || listed[0].DataplaneURL != "https://sandbox.example/list" {
		t.Fatalf("unexpected listed sandboxes: %#v", listed)
	}
}

func TestSandboxRefreshUpdateAndStop(t *testing.T) {
	var stopped bool

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/box-a":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":            "box-id",
				"name":          "box-a",
				"status":        "ready",
				"dataplane_url": "https://sandbox.example/refreshed",
			})
		case r.Method == http.MethodPatch && r.URL.Path == "/v2/sandboxes/boxes/box-a":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":            "box-id",
				"name":          "box-renamed",
				"status":        "ready",
				"dataplane_url": "https://sandbox.example/updated",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/boxes/box-renamed/stop":
			stopped = true
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
	sandbox, err := client.Sandboxes.Boxes.GetSandbox(context.Background(), "box-a")
	if err != nil {
		t.Fatalf("GetSandbox returned error: %v", err)
	}
	if sandbox.DataplaneURL != "https://sandbox.example/refreshed" {
		t.Fatalf("unexpected refreshed dataplane URL: %q", sandbox.DataplaneURL)
	}
	if err := sandbox.Update(context.Background(), langsmith.SandboxBoxUpdateParams{Name: langsmith.String("box-renamed")}); err != nil {
		t.Fatalf("Update returned error: %v", err)
	}
	if sandbox.Name != "box-renamed" || sandbox.DataplaneURL != "https://sandbox.example/updated" {
		t.Fatalf("unexpected updated sandbox: %#v", sandbox)
	}
	if err := sandbox.Stop(context.Background()); err != nil {
		t.Fatalf("Stop returned error: %v", err)
	}
	if !stopped || sandbox.Status != "stopped" || sandbox.DataplaneURL != "" {
		t.Fatalf("unexpected stopped state: stopped=%v sandbox=%#v", stopped, sandbox)
	}
}
