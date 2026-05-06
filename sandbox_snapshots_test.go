package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxCaptureSnapshotAndWait(t *testing.T) {
	var snapshotGets atomic.Int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/boxes/box-a/snapshot":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "snap-1",
				"name":              "snap",
				"status":            "building",
				"fs_capacity_bytes": 1024,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/snapshots/snap-1":
			status := "building"
			if snapshotGets.Add(1) > 1 {
				status = "ready"
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "snap-1",
				"name":              "snap",
				"status":            status,
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
	snapshot, err := client.Sandboxes.Boxes.CaptureSnapshotAndWait(context.Background(), "box-a", langsmith.SandboxBoxNewSnapshotParams{
		Name: langsmith.String("snap"),
	}, langsmith.SnapshotWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("CaptureSnapshotAndWait returned error: %v", err)
	}
	if snapshot.ID != "snap-1" || snapshot.Status != "ready" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	if snapshotGets.Load() != 2 {
		t.Fatalf("expected 2 snapshot get calls, got %d", snapshotGets.Load())
	}
}
