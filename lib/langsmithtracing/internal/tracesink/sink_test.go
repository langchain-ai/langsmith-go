package tracesink

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/models"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/multipart"
)

func makeOp() *models.SerializedOp {
	return &models.SerializedOp{
		Kind:    models.OpKindPost,
		ID:      uuid.New(),
		TraceID: uuid.New(),
		RunInfo: []byte(`{"name":"test"}`),
	}
}

func makeOpWithSize(n int) *models.SerializedOp {
	padding := strings.Repeat("x", n-11)
	return &models.SerializedOp{
		Kind:    models.OpKindPost,
		ID:      uuid.New(),
		TraceID: uuid.New(),
		RunInfo: []byte(fmt.Sprintf(`{"name":"%s"}`, padding)),
	}
}

func testServer(t *testing.T) (*httptest.Server, *atomic.Int64) {
	t.Helper()
	var reqCount atomic.Int64
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		reqCount.Add(1)
		w.WriteHeader(202)
	}))
	t.Cleanup(srv.Close)
	return srv, &reqCount
}

func testDrainConfig(maxQueueSize int) DrainConfig {
	return DrainConfig{
		MaxBatchSize:  100,
		MaxBatchBytes: 20 * 1024 * 1024,
		MaxQueueSize:  maxQueueSize,
		DrainInterval: 50 * time.Millisecond,
		MaxWorkers:    1,
	}
}

func TestQueueFullDrop(t *testing.T) {
	srv, _ := testServer(t)
	exp := multipart.NewExporter(srv.Client(), multipart.RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	cfg := testDrainConfig(3)
	cfg.DrainInterval = 10 * time.Second // don't drain during test
	sink := NewTraceSink(context.Background(), exp, cfg, endpoint, nil, nil)

	for i := 0; i < 3; i++ {
		if err := sink.Submit(makeOp()); err != nil {
			t.Fatalf("submit %d: %v", i, err)
		}
	}

	if got := len(sink.queue); got != 3 {
		t.Fatalf("queue length after 3 submits: got %d, want 3", got)
	}

	if err := sink.Submit(makeOp()); err != nil {
		t.Fatalf("4th submit returned error: %v", err)
	}

	if got := len(sink.queue); got != 3 {
		t.Fatalf("queue length after 4th submit: got %d, want 3", got)
	}

	sink.Close()
}

func TestSubmitAfterClose(t *testing.T) {
	srv, reqCount := testServer(t)
	exp := multipart.NewExporter(srv.Client(), multipart.RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	sink := NewTraceSink(context.Background(), exp, testDrainConfig(100), endpoint, nil, nil)

	sink.Close()

	if err := sink.Submit(makeOp()); err != nil {
		t.Fatalf("submit after close returned error: %v", err)
	}

	if got := reqCount.Load(); got != 0 {
		t.Fatalf("server received %d requests, want 0", got)
	}
}

func TestDoubleClose(t *testing.T) {
	srv, _ := testServer(t)
	exp := multipart.NewExporter(srv.Client(), multipart.RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	sink := NewTraceSink(context.Background(), exp, testDrainConfig(100), endpoint, nil, nil)

	done := make(chan struct{})
	go func() {
		sink.Close()
		sink.Close()
		close(done)
	}()

	select {
	case <-done:
	case <-time.After(5 * time.Second):
		t.Fatal("double Close deadlocked")
	}
}

func TestCollectBatchRespectsMaxBatchBytes(t *testing.T) {
	cfg := testDrainConfig(100)
	cfg.MaxBatchSize = 100
	cfg.MaxBatchBytes = 100
	cfg.DrainInterval = 10 * time.Second

	srv, _ := testServer(t)
	exp := multipart.NewExporter(srv.Client(), multipart.RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	sink := NewTraceSink(context.Background(), exp, cfg, endpoint, nil, nil)
	defer sink.Close()

	for i := 0; i < 5; i++ {
		sink.queue <- makeOpWithSize(50)
	}

	batch, leftover := sink.collectBatch(nil)
	if len(batch) != 2 {
		t.Fatalf("collectBatch returned %d ops, want 2", len(batch))
	}
	if leftover == nil {
		t.Fatal("expected leftover op, got nil")
	}
	if got := len(sink.queue); got != 2 {
		t.Fatalf("remaining queue has %d ops, want 2", got)
	}
}

func TestAllItemsDrained(t *testing.T) {
	srv, reqCount := testServer(t)
	exp := multipart.NewExporter(srv.Client(), multipart.RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	cfg := testDrainConfig(100)
	cfg.MaxBatchSize = 5
	sink := NewTraceSink(context.Background(), exp, cfg, endpoint, nil, nil)

	for i := 0; i < 20; i++ {
		if err := sink.Submit(makeOp()); err != nil {
			t.Fatalf("submit %d: %v", i, err)
		}
	}

	sink.Close()

	if got := reqCount.Load(); got == 0 {
		t.Fatal("server received 0 requests, expected at least 1")
	}
}

func TestDefaultDrainConfigValues(t *testing.T) {
	cfg := DefaultDrainConfig()

	wantMaxWorkers := runtime.GOMAXPROCS(0)
	if wantMaxWorkers > 32 {
		wantMaxWorkers = 32
	}

	checks := []struct {
		name string
		got  any
		want any
	}{
		{"MaxBatchSize", cfg.MaxBatchSize, 100},
		{"MaxBatchBytes", cfg.MaxBatchBytes, 20 * 1024 * 1024},
		{"MaxQueueSize", cfg.MaxQueueSize, 10_000},
		{"DrainInterval", cfg.DrainInterval, 250 * time.Millisecond},
		{"MaxWorkers", cfg.MaxWorkers, wantMaxWorkers},
	}

	for _, c := range checks {
		if c.got != c.want {
			t.Errorf("%s: got %v, want %v", c.name, c.got, c.want)
		}
	}
}
