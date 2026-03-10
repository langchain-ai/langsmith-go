package langsmithtracing_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/google/uuid"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing"
)

func zstdDecompress(t *testing.T, data []byte) []byte {
	t.Helper()
	dec, err := zstd.NewReader(nil)
	if err != nil {
		t.Fatalf("create zstd reader: %v", err)
	}
	defer dec.Close()
	out, err := dec.DecodeAll(data, nil)
	if err != nil {
		t.Fatalf("zstd decompress: %v", err)
	}
	return out
}

// This test sends real traces to LangSmith via the multipart ingestion endpoint.
//
// Run with:
//
//	export LANGSMITH_API_KEY="your-key"
//	export LANGSMITH_ENDPOINT="https://api.smith.langchain.com"   # optional, this is the default
//	go test -v -run TestMultipartTracing ./lib/langsmithtracing/ -count=1
func TestMultipartTracing(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set; skipping integration test")
	}

	ctx := context.Background()
	projectName := fmt.Sprintf("__go-multipart-test-%s", time.Now().UTC().Format("20060102-150405"))

	client := langsmithtracing.NewTracingClient(ctx,
		langsmithtracing.WithProject(projectName),
	)

	t.Logf("Sending traces to project %q", projectName)

	// --- Trace 1: simple chain with a child LLM span ---
	now := time.Now().UTC()
	chainID := uuid.New()
	chainStart := now.Add(-2 * time.Second)
	chainDotted := formatDottedOrder(chainStart, chainID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:      chainID,
		TraceID: chainID,
		Name:    "go-multipart-test",
		RunType: "chain",
		Inputs:  map[string]any{"question": "What is the capital of France?"},
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test"},
		},
		Tags:        []string{"test", "go-sdk"},
		StartTime:   chainStart,
		DottedOrder: chainDotted,
	}); err != nil {
		t.Fatalf("CreateRun (chain): %v", err)
	}

	llmID := uuid.New()
	llmStart := now.Add(-1500 * time.Millisecond)
	llmDotted := chainDotted + "." + formatDottedOrder(llmStart, llmID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:          llmID,
		TraceID:     chainID,
		ParentRunID: &chainID,
		Name:        "go-openai-chat",
		RunType:     "llm",
		Inputs: map[string]any{
			"messages": []map[string]any{
				{"role": "user", "content": "What is the capital of France?"},
			},
		},
		Extra: map[string]any{
			"invocation_params": map[string]any{
				"model":       "gpt-4",
				"temperature": 0.0,
			},
		},
		Tags:        []string{"test", "openai"},
		StartTime:   llmStart,
		DottedOrder: llmDotted,
	}); err != nil {
		t.Fatalf("CreateRun (llm): %v", err)
	}

	llmEnd := now.Add(-800 * time.Millisecond)
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:      llmID,
		TraceID: chainID,
		Outputs: map[string]any{
			"content": "The capital of France is Paris.",
			"usage": map[string]any{
				"prompt_tokens":     12,
				"completion_tokens": 8,
				"total_tokens":      20,
			},
		},
		EndTime:     llmEnd,
		DottedOrder: llmDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (llm): %v", err)
	}

	chainEnd := now.Add(-500 * time.Millisecond)
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:      chainID,
		TraceID: chainID,
		Outputs: map[string]any{
			"answer": "The capital of France is Paris.",
		},
		EndTime:     chainEnd,
		DottedOrder: chainDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (chain): %v", err)
	}

	// --- Trace 2: standalone tool call with attachments ---
	toolID := uuid.New()
	toolStart := now
	toolDotted := formatDottedOrder(toolStart, toolID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:      toolID,
		TraceID: toolID,
		Name:    "calculator",
		RunType: "tool",
		Inputs:  map[string]any{"expression": "2 + 2"},
		Extra: map[string]any{
			"metadata": map[string]any{"tool_name": "calculator", "source": "go-sdk-test"},
		},
		Tags:        []string{"test", "tool"},
		StartTime:   toolStart,
		DottedOrder: toolDotted,
		Attachments: map[string]langsmithtracing.Attachment{
			"test_payload": {
				ContentType: "application/json",
				Data:        []byte(`{"hello": "world"}`),
			},
			"readme": {
				ContentType: "text/plain",
				Data:        []byte("This is an attachment from the Go SDK integration test."),
			},
		},
	}); err != nil {
		t.Fatalf("CreateRun (tool): %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:          toolID,
		TraceID:     toolID,
		Outputs:     map[string]any{"result": 4},
		EndTime:     toolStart.Add(50 * time.Millisecond),
		DottedOrder: toolDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (tool): %v", err)
	}

	// --- Trace 3: errored run ---
	errID := uuid.New()
	errStart := now.Add(100 * time.Millisecond)
	errDotted := formatDottedOrder(errStart, errID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:      errID,
		TraceID: errID,
		Name:    "failing-chain",
		RunType: "chain",
		Inputs:  map[string]any{"query": "trigger error"},
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test", "expected_error": true},
		},
		Tags:        []string{"test", "error"},
		StartTime:   errStart,
		DottedOrder: errDotted,
	}); err != nil {
		t.Fatalf("CreateRun (error): %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:          errID,
		TraceID:     errID,
		Error:       "ValueError: something went wrong in the chain",
		EndTime:     errStart.Add(200 * time.Millisecond),
		DottedOrder: errDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (error): %v", err)
	}

	// --- Trace 4: attachment on update (proving update-time attachments work) ---
	updateAttachID := uuid.New()
	updateAttachStart := now.Add(300 * time.Millisecond)
	updateAttachDotted := formatDottedOrder(updateAttachStart, updateAttachID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:      updateAttachID,
		TraceID: updateAttachID,
		Name:    "attach-on-update",
		RunType: "chain",
		Inputs:  map[string]any{"prompt": "generate image"},
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test", "attachment_stage": "update"},
		},
		Tags:        []string{"test", "attachment"},
		StartTime:   updateAttachStart,
		DottedOrder: updateAttachDotted,
	}); err != nil {
		t.Fatalf("CreateRun (attach-on-update): %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:      updateAttachID,
		TraceID: updateAttachID,
		Outputs: map[string]any{"image_url": "data:image/png;base64,..."},
		EndTime: updateAttachStart.Add(100 * time.Millisecond),
		DottedOrder: updateAttachDotted,
		Attachments: map[string]langsmithtracing.Attachment{
			"generated_image": {
				ContentType: "image/png",
				Data:        []byte{0x89, 0x50, 0x4E, 0x47, 0x0D, 0x0A, 0x1A, 0x0A}, // PNG magic bytes
			},
		},
	}); err != nil {
		t.Fatalf("UpdateRun (attach-on-update): %v", err)
	}

	// --- Trace 5: empty-but-not-nil fields (inputs/outputs sent as {}) ---
	emptyID := uuid.New()
	emptyStart := now.Add(500 * time.Millisecond)
	emptyDotted := formatDottedOrder(emptyStart, emptyID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:          emptyID,
		TraceID:     emptyID,
		Name:        "empty-fields",
		RunType:     "chain",
		Inputs:      map[string]any{},
		Extra:       map[string]any{},
		StartTime:   emptyStart,
		DottedOrder: emptyDotted,
	}); err != nil {
		t.Fatalf("CreateRun (empty-fields): %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:          emptyID,
		TraceID:     emptyID,
		Outputs:     map[string]any{},
		EndTime:     emptyStart.Add(50 * time.Millisecond),
		DottedOrder: emptyDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (empty-fields): %v", err)
	}

	// --- Trace 6: runtime + metadata in extra (both auto-injected and user-provided) ---
	runtimeID := uuid.New()
	runtimeStart := now.Add(600 * time.Millisecond)
	runtimeDotted := formatDottedOrder(runtimeStart, runtimeID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:      runtimeID,
		TraceID: runtimeID,
		Name:    "runtime-and-metadata",
		RunType: "chain",
		Inputs:  map[string]any{"note": "verify extra.runtime and extra.metadata in the UI"},
		Extra: map[string]any{
			"runtime":  map[string]any{"sdk": "langsmith-go", "runtime": "go", "custom_key": "user-provided"},
			"metadata": map[string]any{"user_meta": "should_be_preserved"},
		},
		StartTime:   runtimeStart,
		DottedOrder: runtimeDotted,
	}); err != nil {
		t.Fatalf("CreateRun (runtime-and-metadata): %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:          runtimeID,
		TraceID:     runtimeID,
		Outputs:     map[string]any{"result": "ok"},
		EndTime:     runtimeStart.Add(50 * time.Millisecond),
		DottedOrder: runtimeDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (runtime-and-metadata): %v", err)
	}

	// --- Trace 7: LLM run with first-token event ---
	streamID := uuid.New()
	streamStart := now.Add(700 * time.Millisecond)
	streamDotted := formatDottedOrder(streamStart, streamID)
	firstTokenTime := streamStart.Add(120 * time.Millisecond)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:      streamID,
		TraceID: streamID,
		Name:    "streaming-llm",
		RunType: "llm",
		Inputs: map[string]any{
			"messages": []map[string]any{
				{"role": "user", "content": "Say hello"},
			},
		},
		Extra: map[string]any{
			"metadata":          map[string]any{"source": "go-sdk-test"},
			"invocation_params": map[string]any{"model": "gpt-4", "stream": true},
		},
		Tags:        []string{"test", "streaming"},
		StartTime:   streamStart,
		DottedOrder: streamDotted,
	}); err != nil {
		t.Fatalf("CreateRun (streaming-llm): %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:      streamID,
		TraceID: streamID,
		Outputs: map[string]any{
			"content": "Hello! How can I help you?",
		},
		Events: []map[string]any{
			{
				"name": "new_token",
				"time": firstTokenTime.UTC().Format(time.RFC3339Nano),
				"kwargs": map[string]any{
					"token": "Hello",
				},
			},
		},
		EndTime:     streamStart.Add(500 * time.Millisecond),
		DottedOrder: streamDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (streaming-llm): %v", err)
	}

	t.Log("All runs submitted, flushing...")
	client.Close()

	t.Logf("Done. Check LangSmith project %q for 7 traces:", projectName)
	t.Log("  1. go-multipart-test (chain) -> go-openai-chat (llm child)")
	t.Log("  2. calculator (tool) with 2 attachments: test_payload, readme")
	t.Log("  3. failing-chain (errored)")
	t.Log("  4. attach-on-update (chain) with attachment added at update time")
	t.Log("  5. empty-fields (chain) with empty {} inputs/outputs/extra")
	t.Log("  6. runtime-and-metadata — extra.runtime has sdk/runtime/custom_key, extra.metadata has user_meta + env vars")
	t.Log("  7. streaming-llm — LLM with new_token event (first token timing)")
}

// TestBatchFallbackOn404 verifies the full client → sink → exporter pipeline
// falls back to /runs/batch when /runs/multipart returns 404.
// This uses a local httptest server and does not require LANGSMITH_API_KEY.
func TestBatchFallbackOn404(t *testing.T) {
	var mu sync.Mutex
	var multipartCalls, batchCalls int
	var batchBodies [][]byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		switch r.URL.Path {
		case "/runs/multipart":
			multipartCalls++
			w.WriteHeader(http.StatusNotFound)
		case "/runs/batch":
			batchCalls++
			body, _ := io.ReadAll(r.Body)
			batchBodies = append(batchBodies, body)
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	ctx := context.Background()
	client := langsmithtracing.NewTracingClient(ctx,
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("fallback-test"),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID:          runID,
		TraceID:     runID,
		Name:        "fallback-run",
		RunType:     "chain",
		Inputs:      map[string]any{"q": "hello"},
		StartTime:   now,
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	endTime := now.Add(100 * time.Millisecond)
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID:          runID,
		TraceID:     runID,
		Outputs:     map[string]any{"a": "world"},
		EndTime:     endTime,
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()

	if multipartCalls < 1 {
		t.Errorf("expected at least 1 multipart attempt, got %d", multipartCalls)
	}
	if batchCalls < 1 {
		t.Fatalf("expected at least 1 batch fallback call, got %d", batchCalls)
	}

	var parsed struct {
		Post  []map[string]any `json:"post"`
		Patch []map[string]any `json:"patch"`
	}
	if err := json.Unmarshal(batchBodies[0], &parsed); err != nil {
		t.Fatalf("unmarshal batch body: %v", err)
	}

	totalRuns := len(parsed.Post) + len(parsed.Patch)
	if totalRuns == 0 {
		t.Fatal("batch body contained no runs")
	}

	t.Logf("multipart attempts: %d, batch calls: %d, runs in first batch: post=%d patch=%d",
		multipartCalls, batchCalls, len(parsed.Post), len(parsed.Patch))
}

// TestAutoScalingWorkers sends a burst of runs to LangSmith to exercise the
// auto-scaling worker pool. It uses a small batch size and low scale-up trigger
// so the sink spawns concurrent export goroutines under load.
//
// Run with:
//
//	export LANGSMITH_API_KEY="your-key"
//	go test -v -run TestAutoScalingWorkers ./lib/langsmithtracing/ -count=1
func TestAutoScalingWorkers(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set; skipping integration test")
	}

	projectName := fmt.Sprintf("__go-multipart-test-%s", time.Now().UTC().Format("20060102-150405"))

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.MaxBatchSize = 10
	cfg.DrainInterval = 50 * time.Millisecond
	cfg.ScaleUpQueueTrigger = 5
	cfg.MaxWorkers = 4
	cfg.ScaleDownEmptyTrigger = 2

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithProject(projectName),
		langsmithtracing.WithDrainConfig(cfg),
	)

	const numRuns = 200
	now := time.Now().UTC()
	t.Logf("Sending %d runs (create+update each) to project %q with MaxWorkers=%d", numRuns, projectName, cfg.MaxWorkers)

	type runInfo struct {
		id      uuid.UUID
		dotted  string
		startTs time.Time
	}
	runs := make([]runInfo, numRuns)

	start := time.Now()
	for i := 0; i < numRuns; i++ {
		id := uuid.New()
		ts := now.Add(time.Duration(i) * time.Millisecond)
		dotted := formatDottedOrder(ts, id)
		runs[i] = runInfo{id: id, dotted: dotted, startTs: ts}

		if err := client.CreateRun(&langsmithtracing.RunCreate{
			ID:      id,
			TraceID: id,
			Name:    fmt.Sprintf("scale-test-%d", i),
			RunType: "chain",
			Inputs:  map[string]any{"i": i},
			Extra: map[string]any{
				"metadata": map[string]any{"source": "go-sdk-test", "batch_index": i},
			},
			Tags:        []string{"test", "scaling"},
			StartTime:   ts,
			DottedOrder: dotted,
		}); err != nil {
			t.Fatalf("CreateRun %d: %v", i, err)
		}
	}

	for i, r := range runs {
		if err := client.UpdateRun(&langsmithtracing.RunUpdate{
			ID:          r.id,
			TraceID:     r.id,
			Outputs:     map[string]any{"result": fmt.Sprintf("done-%d", i)},
			EndTime:     r.startTs.Add(10 * time.Millisecond),
			DottedOrder: r.dotted,
		}); err != nil {
			t.Fatalf("UpdateRun %d: %v", i, err)
		}
	}

	client.Close()
	elapsed := time.Since(start)

	t.Logf("Done. %d runs flushed in %s. Check project %q in LangSmith.", numRuns, elapsed, projectName)
}

// TestAutoScalingConcurrency is a local-only test (no API key needed) that
// verifies the sink actually achieves concurrent exports via a mock server.
func TestAutoScalingConcurrency(t *testing.T) {
	var (
		mu           sync.Mutex
		inflight     int
		peakInflight int
		totalExports int
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		inflight++
		if inflight > peakInflight {
			peakInflight = inflight
		}
		totalExports++
		mu.Unlock()

		time.Sleep(100 * time.Millisecond)

		mu.Lock()
		inflight--
		mu.Unlock()

		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.MaxBatchSize = 10
	cfg.DrainInterval = 50 * time.Millisecond
	cfg.ScaleUpQueueTrigger = 5
	cfg.MaxWorkers = 4
	cfg.ScaleDownEmptyTrigger = 2

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("concurrency-test"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()
	for i := 0; i < 200; i++ {
		id := uuid.New()
		ts := now.Add(time.Duration(i) * time.Millisecond)
		dotted := formatDottedOrder(ts, id)
		if err := client.CreateRun(&langsmithtracing.RunCreate{
			ID:      id,
			TraceID: id,
			Name:    fmt.Sprintf("concurrency-test-%d", i),
			RunType: "chain",
			Inputs:  map[string]any{"i": i},
			Extra: map[string]any{
				"metadata": map[string]any{"source": "go-sdk-test", "batch_index": i},
			},
			Tags:        []string{"test", "scaling"},
			StartTime:   ts,
			DottedOrder: dotted,
		}); err != nil {
			t.Fatalf("CreateRun %d: %v", i, err)
		}
	}

	client.Close()

	mu.Lock()
	peak := peakInflight
	exports := totalExports
	mu.Unlock()

	t.Logf("peak concurrent exports: %d, total export calls: %d", peak, exports)
	if peak <= 1 {
		t.Errorf("expected peak concurrency > 1 (auto-scaling should spawn workers), got %d", peak)
	}
}

// TestSamplingRateZero verifies that with sample rate 0, no runs are exported.
func TestSamplingRateZero(t *testing.T) {
	var mu sync.Mutex
	var totalRequests int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		totalRequests++
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("sampling-zero"),
		langsmithtracing.WithSampleRate(0),
	)

	now := time.Now().UTC()
	for i := 0; i < 20; i++ {
		id := uuid.New()
		ts := now.Add(time.Duration(i) * time.Millisecond)
		dotted := formatDottedOrder(ts, id)
		_ = client.CreateRun(&langsmithtracing.RunCreate{
			ID: id, TraceID: id, Name: fmt.Sprintf("sampled-out-%d", i),
			RunType: "chain", Inputs: map[string]any{"i": i},
			StartTime: ts, DottedOrder: dotted,
		})
		_ = client.UpdateRun(&langsmithtracing.RunUpdate{
			ID: id, TraceID: id,
			Outputs: map[string]any{"r": i}, EndTime: ts.Add(time.Millisecond),
			DottedOrder: dotted,
		})
	}

	client.Close()

	mu.Lock()
	reqs := totalRequests
	mu.Unlock()

	if reqs != 0 {
		t.Errorf("sample rate 0: expected 0 requests, got %d", reqs)
	}
}

// TestSamplingRateOne verifies that with sample rate 1, all runs are exported.
func TestSamplingRateOne(t *testing.T) {
	var mu sync.Mutex
	var totalRequests int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		totalRequests++
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("sampling-one"),
		langsmithtracing.WithSampleRate(1),
	)

	now := time.Now().UTC()
	for i := 0; i < 5; i++ {
		id := uuid.New()
		ts := now.Add(time.Duration(i) * time.Millisecond)
		dotted := formatDottedOrder(ts, id)
		_ = client.CreateRun(&langsmithtracing.RunCreate{
			ID: id, TraceID: id, Name: fmt.Sprintf("sampled-in-%d", i),
			RunType: "chain", Inputs: map[string]any{"i": i},
			StartTime: ts, DottedOrder: dotted,
		})
		_ = client.UpdateRun(&langsmithtracing.RunUpdate{
			ID: id, TraceID: id,
			Outputs: map[string]any{"r": i}, EndTime: ts.Add(time.Millisecond),
			DottedOrder: dotted,
		})
	}

	client.Close()

	mu.Lock()
	reqs := totalRequests
	mu.Unlock()

	if reqs == 0 {
		t.Error("sample rate 1: expected at least 1 request, got 0")
	}
	t.Logf("sample rate 1: %d requests for 5 runs", reqs)
}

// TestSamplingChildFollowsParent verifies that when a root run is sampled out,
// its child runs and updates are also dropped.
func TestSamplingChildFollowsParent(t *testing.T) {
	var mu sync.Mutex
	var totalRequests int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		totalRequests++
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("sampling-child"),
		langsmithtracing.WithSampleRate(0),
	)

	now := time.Now().UTC()
	parentID := uuid.New()
	parentDotted := formatDottedOrder(now, parentID)

	_ = client.CreateRun(&langsmithtracing.RunCreate{
		ID: parentID, TraceID: parentID, Name: "parent",
		RunType: "chain", Inputs: map[string]any{"q": "hello"},
		StartTime: now, DottedOrder: parentDotted,
	})

	childID := uuid.New()
	childStart := now.Add(time.Millisecond)
	childDotted := parentDotted + "." + formatDottedOrder(childStart, childID)

	_ = client.CreateRun(&langsmithtracing.RunCreate{
		ID: childID, TraceID: parentID, ParentRunID: &parentID,
		Name: "child", RunType: "llm",
		Inputs:    map[string]any{"msg": "hi"},
		StartTime: childStart, DottedOrder: childDotted,
	})

	_ = client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: childID, TraceID: parentID,
		Outputs: map[string]any{"resp": "hey"}, EndTime: childStart.Add(time.Millisecond),
		DottedOrder: childDotted,
	})

	_ = client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: parentID, TraceID: parentID,
		Outputs: map[string]any{"a": "world"}, EndTime: now.Add(2 * time.Millisecond),
		DottedOrder: parentDotted,
	})

	client.Close()

	mu.Lock()
	reqs := totalRequests
	mu.Unlock()

	if reqs != 0 {
		t.Errorf("sampling parent+child with rate=0: expected 0 requests, got %d", reqs)
	}
}

// TestSerializedPartLLMvsChain verifies that the .serialized multipart part is
// emitted for "llm" and "prompt" run types (with the "graph" key stripped) and
// is NOT emitted for other run types like "chain".
func TestSerializedPartLLMvsChain(t *testing.T) {
	var mu sync.Mutex
	var capturedBodies [][]byte
	var capturedContentTypes []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		capturedBodies = append(capturedBodies, body)
		capturedContentTypes = append(capturedContentTypes, r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("serialized-test"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()

	// Run 1: LLM with Serialized (should be kept, graph stripped)
	llmID := uuid.New()
	llmDotted := formatDottedOrder(now, llmID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: llmID, TraceID: llmID, Name: "llm-with-serialized",
		RunType: "llm",
		Inputs:  map[string]any{"msg": "hello"},
		Serialized: map[string]any{
			"name":  "ChatOpenAI",
			"type":  "llm",
			"graph": map[string]any{"nodes": []any{"a", "b"}},
		},
		StartTime: now, DottedOrder: llmDotted,
	}); err != nil {
		t.Fatalf("CreateRun (llm): %v", err)
	}
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: llmID, TraceID: llmID,
		Outputs: map[string]any{"content": "hi"},
		EndTime: now.Add(time.Millisecond), DottedOrder: llmDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (llm): %v", err)
	}

	// Run 2: chain with Serialized (should be dropped entirely)
	chainID := uuid.New()
	chainDotted := formatDottedOrder(now.Add(10*time.Millisecond), chainID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: chainID, TraceID: chainID, Name: "chain-with-serialized",
		RunType: "chain",
		Inputs:  map[string]any{"q": "test"},
		Serialized: map[string]any{
			"name": "MyChain",
			"type": "chain",
		},
		StartTime: now.Add(10 * time.Millisecond), DottedOrder: chainDotted,
	}); err != nil {
		t.Fatalf("CreateRun (chain): %v", err)
	}
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: chainID, TraceID: chainID,
		Outputs: map[string]any{"a": "b"},
		EndTime: now.Add(11 * time.Millisecond), DottedOrder: chainDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (chain): %v", err)
	}

	// Run 3: prompt with Serialized (should be kept, graph stripped)
	promptID := uuid.New()
	promptDotted := formatDottedOrder(now.Add(20*time.Millisecond), promptID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: promptID, TraceID: promptID, Name: "prompt-with-serialized",
		RunType: "prompt",
		Inputs:  map[string]any{"template": "hello {name}"},
		Serialized: map[string]any{
			"name":  "ChatPromptTemplate",
			"type":  "prompt",
			"graph": map[string]any{"edges": []any{1, 2}},
		},
		StartTime: now.Add(20 * time.Millisecond), DottedOrder: promptDotted,
	}); err != nil {
		t.Fatalf("CreateRun (prompt): %v", err)
	}
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: promptID, TraceID: promptID,
		Outputs: map[string]any{"formatted": "hello world"},
		EndTime: now.Add(21 * time.Millisecond), DottedOrder: promptDotted,
	}); err != nil {
		t.Fatalf("UpdateRun (prompt): %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()

	// Parse all captured multipart bodies and collect parts by name.
	allParts := make(map[string][]byte)
	for i, body := range capturedBodies {
		decompressed := zstdDecompress(t, body)
		_, params, err := mime.ParseMediaType(capturedContentTypes[i])
		if err != nil {
			t.Fatalf("parse content-type %d: %v", i, err)
		}
		boundary := params["boundary"]
		if boundary == "" {
			t.Fatalf("missing boundary in body %d", i)
		}
		reader := multipart.NewReader(bytes.NewReader(decompressed), boundary)
		for {
			p, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next part body %d: %v", i, err)
			}
			data, _ := io.ReadAll(p)
			allParts[p.FormName()] = data
		}
	}

	// LLM: .serialized part should exist with graph stripped.
	llmSerKey := "post." + llmID.String() + ".serialized"
	if raw, ok := allParts[llmSerKey]; !ok {
		t.Errorf("LLM run: expected part %q to be present", llmSerKey)
	} else {
		var parsed map[string]any
		if err := json.Unmarshal(raw, &parsed); err != nil {
			t.Fatalf("unmarshal LLM serialized: %v", err)
		}
		if parsed["name"] != "ChatOpenAI" {
			t.Errorf("LLM serialized.name = %v, want ChatOpenAI", parsed["name"])
		}
		if _, hasGraph := parsed["graph"]; hasGraph {
			t.Error("LLM serialized should NOT contain 'graph' (it should be stripped)")
		}
	}

	// Chain: .serialized part should NOT exist.
	chainSerKey := "post." + chainID.String() + ".serialized"
	if _, ok := allParts[chainSerKey]; ok {
		t.Error("chain run: .serialized part should NOT be emitted for run_type=chain")
	}

	// Prompt: .serialized part should exist with graph stripped.
	promptSerKey := "post." + promptID.String() + ".serialized"
	if raw, ok := allParts[promptSerKey]; !ok {
		t.Errorf("prompt run: expected part %q to be present", promptSerKey)
	} else {
		var parsed map[string]any
		if err := json.Unmarshal(raw, &parsed); err != nil {
			t.Fatalf("unmarshal prompt serialized: %v", err)
		}
		if parsed["name"] != "ChatPromptTemplate" {
			t.Errorf("prompt serialized.name = %v, want ChatPromptTemplate", parsed["name"])
		}
		if _, hasGraph := parsed["graph"]; hasGraph {
			t.Error("prompt serialized should NOT contain 'graph' (it should be stripped)")
		}
	}

	// Log all parts for debugging.
	for name := range allParts {
		if strings.Contains(name, ".serialized") {
			t.Logf("  serialized part: %s (%d bytes)", name, len(allParts[name]))
		}
	}
}

// TestUpdateRunFields verifies that all optional fields on RunUpdate
// (Inputs, Name, RunType, StartTime, Tags, Extra) are sent as multipart parts.
// Uses a local mock server — no API key needed.
func TestUpdateRunFields(t *testing.T) {
	var mu sync.Mutex
	var capturedBodies [][]byte
	var capturedContentTypes []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		capturedBodies = append(capturedBodies, body)
		capturedContentTypes = append(capturedContentTypes, r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("update-fields-test"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name:    "original-name",
		RunType: "chain",
		Inputs:  map[string]any{"original": true},
		Extra:   map[string]any{"metadata": map[string]any{"version": "v1"}},
		Tags:    []string{"initial"},
		StartTime: now, DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	adjustedStart := now.Add(-500 * time.Millisecond)
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: runID, TraceID: runID,
		Name:    "renamed-run",
		RunType: "llm",
		Inputs:  map[string]any{"replaced": true, "prompt": "hello"},
		Outputs: map[string]any{"answer": "world"},
		Extra:   map[string]any{"metadata": map[string]any{"version": "v2", "updated": true}},
		Tags:    []string{"updated", "v2"},
		StartTime: adjustedStart,
		EndTime:   now.Add(100 * time.Millisecond),
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()

	allParts := make(map[string][]byte)
	for i, body := range capturedBodies {
		decompressed := zstdDecompress(t, body)
		_, params, err := mime.ParseMediaType(capturedContentTypes[i])
		if err != nil {
			t.Fatalf("parse content-type %d: %v", i, err)
		}
		reader := multipart.NewReader(bytes.NewReader(decompressed), params["boundary"])
		for {
			p, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next part: %v", err)
			}
			data, _ := io.ReadAll(p)
			allParts[p.FormName()] = data
		}
	}

	idStr := runID.String()

	// Because create+update for the same run are coalesced into a single "post" op,
	// the update fields should be merged into the post. Check the coalesced post.
	postKey := "post." + idStr
	if raw, ok := allParts[postKey]; ok {
		var info map[string]any
		if err := json.Unmarshal(raw, &info); err != nil {
			t.Fatalf("unmarshal post run info: %v", err)
		}
		if info["name"] != "renamed-run" {
			t.Errorf("name = %v, want renamed-run", info["name"])
		}
		if info["run_type"] != "llm" {
			t.Errorf("run_type = %v, want llm", info["run_type"])
		}
		if _, ok := info["start_time"]; !ok {
			t.Error("start_time should be present in coalesced run info")
		}
		if tags, ok := info["tags"].([]any); !ok || len(tags) == 0 {
			t.Errorf("tags = %v, want non-empty", info["tags"])
		}
		t.Logf("coalesced post run info: %s", string(raw))
	} else {
		// If they arrived in separate batches, check the patch.
		patchKey := "patch." + idStr
		raw, ok := allParts[patchKey]
		if !ok {
			t.Fatal("neither post nor patch part found for run")
		}
		var info map[string]any
		if err := json.Unmarshal(raw, &info); err != nil {
			t.Fatalf("unmarshal patch run info: %v", err)
		}
		if info["name"] != "renamed-run" {
			t.Errorf("name = %v, want renamed-run", info["name"])
		}
		if info["run_type"] != "llm" {
			t.Errorf("run_type = %v, want llm", info["run_type"])
		}
		t.Logf("patch run info: %s", string(raw))
	}

	// Check inputs part (should contain the updated inputs).
	inputsKey := "post." + idStr + ".inputs"
	if _, ok := allParts[inputsKey]; !ok {
		inputsKey = "patch." + idStr + ".inputs"
	}
	if raw, ok := allParts[inputsKey]; ok {
		var inputs map[string]any
		if err := json.Unmarshal(raw, &inputs); err != nil {
			t.Fatalf("unmarshal inputs: %v", err)
		}
		if inputs["replaced"] != true {
			t.Errorf("inputs.replaced = %v, want true", inputs["replaced"])
		}
		t.Logf("inputs: %s", string(raw))
	} else {
		t.Error("inputs part not found")
	}

	// Check extra part.
	extraKey := "post." + idStr + ".extra"
	if _, ok := allParts[extraKey]; !ok {
		extraKey = "patch." + idStr + ".extra"
	}
	if raw, ok := allParts[extraKey]; ok {
		t.Logf("extra: %s", string(raw))
	} else {
		t.Error("extra part not found")
	}

	// Check outputs part.
	outputsKey := "post." + idStr + ".outputs"
	if _, ok := allParts[outputsKey]; !ok {
		outputsKey = "patch." + idStr + ".outputs"
	}
	if raw, ok := allParts[outputsKey]; ok {
		var outputs map[string]any
		if err := json.Unmarshal(raw, &outputs); err != nil {
			t.Fatalf("unmarshal outputs: %v", err)
		}
		if outputs["answer"] != "world" {
			t.Errorf("outputs.answer = %v, want world", outputs["answer"])
		}
	} else {
		t.Error("outputs part not found")
	}
}

// TestUpdateRunFieldsLive sends a trace to LangSmith that exercises all the
// optional update fields (name, run_type, inputs, extra, tags, start_time).
// Requires LANGSMITH_API_KEY.
func TestUpdateRunFieldsLive(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set; skipping integration test")
	}

	projectName := fmt.Sprintf("__go-multipart-test-%s", time.Now().UTC().Format("20060102-150405"))
	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithProject(projectName),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name:    "original-name",
		RunType: "chain",
		Inputs:  map[string]any{"question": "original"},
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test", "stage": "create"},
		},
		Tags:      []string{"create-tag"},
		StartTime: now, DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	adjustedStart := now.Add(-1 * time.Second)
	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: runID, TraceID: runID,
		Name:    "renamed-on-update",
		RunType: "llm",
		Inputs:  map[string]any{"question": "replaced on update", "model": "gpt-4"},
		Outputs: map[string]any{"answer": "updated answer"},
		Extra: map[string]any{
			"metadata":          map[string]any{"source": "go-sdk-test", "stage": "update"},
			"invocation_params": map[string]any{"model": "gpt-4", "temperature": 0.5},
		},
		Tags:      []string{"update-tag", "v2"},
		StartTime: adjustedStart,
		EndTime:   now.Add(200 * time.Millisecond),
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	t.Logf("Done. Check project %q for run %s:", projectName, runID)
	t.Log("  - name should be 'renamed-on-update'")
	t.Log("  - run_type should be 'llm'")
	t.Log("  - inputs should have 'question: replaced on update' and 'model: gpt-4'")
	t.Log("  - tags should include 'update-tag' and 'v2'")
	t.Log("  - start_time should be ~1s before the original")
	t.Log("  - extra.metadata.stage should be 'update'")
}

// TestRunInfoFields verifies that session_name, session_id, reference_example_id,
// input_attachments, and output_attachments appear in the multipart run JSON.
func TestRunInfoFields(t *testing.T) {
	var mu sync.Mutex
	var capturedBodies [][]byte
	var capturedContentTypes []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		capturedBodies = append(capturedBodies, body)
		capturedContentTypes = append(capturedContentTypes, r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("default-project"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)
	sessionID := uuid.New()
	refExampleID := uuid.New()

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name: "run-info-test", RunType: "chain",
		Inputs:    map[string]any{"q": "hello"},
		StartTime: now, DottedOrder: dotted,
		SessionName:        "per-run-project",
		SessionID:          &sessionID,
		ReferenceExampleID: &refExampleID,
		InputAttachments: map[string]any{
			"my_doc": map[string]any{"presigned_url": "https://example.com/doc.pdf"},
		},
		OutputAttachments: map[string]any{
			"result_img": map[string]any{"presigned_url": "https://example.com/img.png"},
		},
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: runID, TraceID: runID,
		Outputs: map[string]any{"a": "world"},
		EndTime: now.Add(time.Millisecond), DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()

	allParts := make(map[string][]byte)
	for i, body := range capturedBodies {
		decompressed := zstdDecompress(t, body)
		_, params, err := mime.ParseMediaType(capturedContentTypes[i])
		if err != nil {
			t.Fatalf("parse content-type %d: %v", i, err)
		}
		reader := multipart.NewReader(bytes.NewReader(decompressed), params["boundary"])
		for {
			p, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next part: %v", err)
			}
			data, _ := io.ReadAll(p)
			allParts[p.FormName()] = data
		}
	}

	postKey := "post." + runID.String()
	raw, ok := allParts[postKey]
	if !ok {
		t.Fatalf("missing part %q", postKey)
	}

	var info map[string]any
	if err := json.Unmarshal(raw, &info); err != nil {
		t.Fatalf("unmarshal run info: %v", err)
	}

	if info["session_name"] != "per-run-project" {
		t.Errorf("session_name = %v, want per-run-project", info["session_name"])
	}
	if info["session_id"] != sessionID.String() {
		t.Errorf("session_id = %v, want %s", info["session_id"], sessionID)
	}
	if info["reference_example_id"] != refExampleID.String() {
		t.Errorf("reference_example_id = %v, want %s", info["reference_example_id"], refExampleID)
	}
	if info["input_attachments"] == nil {
		t.Error("input_attachments should be present")
	}
	if info["output_attachments"] == nil {
		t.Error("output_attachments should be present")
	}
	t.Logf("run info: %s", string(raw))
}

// TestRunInfoFieldsLive sends a trace with session_name override and
// reference_example_id to LangSmith. Requires LANGSMITH_API_KEY.
func TestRunInfoFieldsLive(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set; skipping integration test")
	}

	projectName := fmt.Sprintf("__go-multipart-test-%s", time.Now().UTC().Format("20060102-150405"))
	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithProject(projectName),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)
	refExampleID := uuid.New()

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name: "run-info-fields-live", RunType: "chain",
		Inputs: map[string]any{"q": "test run info fields"},
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test"},
		},
		StartTime: now, DottedOrder: dotted,
		ReferenceExampleID: &refExampleID,
		InputAttachments: map[string]any{
			"doc": map[string]any{"presigned_url": "https://example.com/doc.pdf"},
		},
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: runID, TraceID: runID,
		Outputs: map[string]any{"result": "ok"},
		EndTime: now.Add(100 * time.Millisecond), DottedOrder: dotted,
		OutputAttachments: map[string]any{
			"report": map[string]any{"presigned_url": "https://example.com/report.pdf"},
		},
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	t.Logf("Done. Check project %q for run %s:", projectName, runID)
	t.Log("  - reference_example_id should be set")
	t.Log("  - input_attachments should have 'doc' entry")
	t.Log("  - output_attachments should have 'report' entry")
}

// TestSingleShotCreateRun verifies that Outputs, EndTime, and Error on RunCreate
// produce a complete run in a single call (no UpdateRun needed).
func TestSingleShotCreateRun(t *testing.T) {
	var mu sync.Mutex
	var capturedBodies [][]byte
	var capturedContentTypes []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		capturedBodies = append(capturedBodies, body)
		capturedContentTypes = append(capturedContentTypes, r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("single-shot-test"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()

	// Run 1: successful single-shot run with outputs + end_time.
	okID := uuid.New()
	okDotted := formatDottedOrder(now, okID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: okID, TraceID: okID,
		Name: "single-shot-ok", RunType: "chain",
		Inputs:    map[string]any{"q": "hello"},
		Outputs:   map[string]any{"a": "world"},
		StartTime: now, EndTime: now.Add(50 * time.Millisecond),
		DottedOrder: okDotted,
	}); err != nil {
		t.Fatalf("CreateRun (ok): %v", err)
	}

	// Run 2: failed single-shot run with error + end_time.
	errID := uuid.New()
	errDotted := formatDottedOrder(now.Add(100*time.Millisecond), errID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: errID, TraceID: errID,
		Name: "single-shot-error", RunType: "chain",
		Inputs:    map[string]any{"q": "fail"},
		Error:     "ValueError: something broke",
		StartTime: now.Add(100 * time.Millisecond),
		EndTime:   now.Add(150 * time.Millisecond),
		DottedOrder: errDotted,
	}); err != nil {
		t.Fatalf("CreateRun (error): %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()

	allParts := make(map[string][]byte)
	for i, body := range capturedBodies {
		decompressed := zstdDecompress(t, body)
		_, params, err := mime.ParseMediaType(capturedContentTypes[i])
		if err != nil {
			t.Fatalf("parse content-type %d: %v", i, err)
		}
		reader := multipart.NewReader(bytes.NewReader(decompressed), params["boundary"])
		for {
			p, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next part: %v", err)
			}
			data, _ := io.ReadAll(p)
			allParts[p.FormName()] = data
		}
	}

	// Verify run 1: should have outputs and end_time, no error.
	okInfoKey := "post." + okID.String()
	if raw, ok := allParts[okInfoKey]; ok {
		var info map[string]any
		if err := json.Unmarshal(raw, &info); err != nil {
			t.Fatalf("unmarshal ok run info: %v", err)
		}
		if info["end_time"] == nil {
			t.Error("ok run: end_time should be set")
		}
		if info["status"] != nil {
			t.Errorf("ok run: status should not be set, got %v", info["status"])
		}
	} else {
		t.Fatalf("missing part %q", okInfoKey)
	}

	okOutputsKey := "post." + okID.String() + ".outputs"
	if raw, ok := allParts[okOutputsKey]; ok {
		var outputs map[string]any
		if err := json.Unmarshal(raw, &outputs); err != nil {
			t.Fatalf("unmarshal ok outputs: %v", err)
		}
		if outputs["a"] != "world" {
			t.Errorf("ok outputs.a = %v, want world", outputs["a"])
		}
	} else {
		t.Error("ok run: outputs part missing")
	}

	// Verify run 2: should have error and status=error.
	errInfoKey := "post." + errID.String()
	if raw, ok := allParts[errInfoKey]; ok {
		var info map[string]any
		if err := json.Unmarshal(raw, &info); err != nil {
			t.Fatalf("unmarshal error run info: %v", err)
		}
		if info["end_time"] == nil {
			t.Error("error run: end_time should be set")
		}
		if info["status"] != "error" {
			t.Errorf("error run: status = %v, want error", info["status"])
		}
	} else {
		t.Fatalf("missing part %q", errInfoKey)
	}

	errErrorKey := "post." + errID.String() + ".error"
	if raw, ok := allParts[errErrorKey]; ok {
		var errStr string
		if err := json.Unmarshal(raw, &errStr); err != nil {
			t.Fatalf("unmarshal error part: %v", err)
		}
		if errStr != "ValueError: something broke" {
			t.Errorf("error = %q, want ValueError: something broke", errStr)
		}
	} else {
		t.Error("error run: .error part missing")
	}
}

// TestSingleShotCreateRunLive sends single-shot complete runs to LangSmith.
// Requires LANGSMITH_API_KEY.
func TestSingleShotCreateRunLive(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set; skipping integration test")
	}

	projectName := fmt.Sprintf("__go-multipart-test-%s", time.Now().UTC().Format("20060102-150405"))
	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithProject(projectName),
	)

	now := time.Now().UTC()

	// Successful single-shot run.
	okID := uuid.New()
	okDotted := formatDottedOrder(now, okID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: okID, TraceID: okID,
		Name: "single-shot-success", RunType: "chain",
		Inputs:  map[string]any{"question": "What is 2+2?"},
		Outputs: map[string]any{"answer": "4"},
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test", "single_shot": true},
		},
		Tags:      []string{"test", "single-shot"},
		StartTime: now, EndTime: now.Add(50 * time.Millisecond),
		DottedOrder: okDotted,
	}); err != nil {
		t.Fatalf("CreateRun (ok): %v", err)
	}

	// Failed single-shot run.
	errID := uuid.New()
	errStart := now.Add(100 * time.Millisecond)
	errDotted := formatDottedOrder(errStart, errID)
	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: errID, TraceID: errID,
		Name: "single-shot-error", RunType: "chain",
		Inputs: map[string]any{"question": "trigger failure"},
		Error:  "RuntimeError: single-shot run failed",
		Extra: map[string]any{
			"metadata": map[string]any{"source": "go-sdk-test", "single_shot": true},
		},
		Tags:      []string{"test", "single-shot", "error"},
		StartTime: errStart, EndTime: errStart.Add(30 * time.Millisecond),
		DottedOrder: errDotted,
	}); err != nil {
		t.Fatalf("CreateRun (error): %v", err)
	}

	client.Close()

	t.Logf("Done. Check project %q for 2 single-shot traces:", projectName)
	t.Logf("  1. single-shot-success (%s) — complete with outputs, no UpdateRun needed", okID)
	t.Logf("  2. single-shot-error (%s) — failed with error, no UpdateRun needed", errID)
}

// TestRetryOnServerError verifies that the client retries on 500 and eventually
// succeeds when the server recovers, exercising the full pipeline (client → sink → exporter).
func TestRetryOnServerError(t *testing.T) {
	var mu sync.Mutex
	var attempts int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		attempts++
		n := attempts
		mu.Unlock()
		if n <= 2 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("temporary failure"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("retry-test"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name: "retry-run", RunType: "chain",
		Inputs:      map[string]any{"q": "test"},
		StartTime:   now,
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: runID, TraceID: runID,
		Outputs:     map[string]any{"a": "ok"},
		EndTime:     now.Add(50 * time.Millisecond),
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()
	if attempts < 3 {
		t.Errorf("expected at least 3 attempts (2 failures + 1 success), got %d", attempts)
	}
	t.Logf("Server received %d request(s) (2 were 500, final was 202)", attempts)
}

// TestNoRetryOnClientError verifies that 4xx errors (except 408) are not retried.
func TestNoRetryOnClientError(t *testing.T) {
	var mu sync.Mutex
	var attempts int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		attempts++
		mu.Unlock()
		w.WriteHeader(http.StatusUnprocessableEntity)
		w.Write([]byte("validation error"))
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("no-retry-test"),
		langsmithtracing.WithDrainConfig(cfg),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name: "no-retry-run", RunType: "chain",
		Inputs:      map[string]any{"q": "test"},
		StartTime:   now,
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()
	if attempts != 1 {
		t.Errorf("expected exactly 1 attempt (no retry on 422), got %d", attempts)
	}
}

// TestRunTransformHook verifies that WithRunTransform receives batched ops
// and can modify run data before it's sent to the server.
func TestRunTransformHook(t *testing.T) {
	var mu sync.Mutex
	var capturedBodies [][]byte
	var capturedContentTypes []string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		body, _ := io.ReadAll(r.Body)
		capturedBodies = append(capturedBodies, body)
		capturedContentTypes = append(capturedContentTypes, r.Header.Get("Content-Type"))
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	cfg := langsmithtracing.DefaultDrainConfig()
	cfg.DrainInterval = 50 * time.Millisecond

	// Transform that adds a "transformed" key to every run's extra.metadata.
	transform := func(ops []langsmithtracing.RunOp) []langsmithtracing.RunOp {
		for i := range ops {
			extra, _ := ops[i].Data["extra"].(map[string]any)
			if extra == nil {
				extra = make(map[string]any)
			}
			meta, _ := extra["metadata"].(map[string]any)
			if meta == nil {
				meta = make(map[string]any)
			}
			meta["transformed"] = true
			meta["hook_version"] = "1.0"
			extra["metadata"] = meta
			ops[i].Data["extra"] = extra
		}
		return ops
	}

	client := langsmithtracing.NewTracingClient(context.Background(),
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithProject("transform-test"),
		langsmithtracing.WithDrainConfig(cfg),
		langsmithtracing.WithRunTransform(transform),
	)

	now := time.Now().UTC()
	runID := uuid.New()
	dotted := formatDottedOrder(now, runID)

	if err := client.CreateRun(&langsmithtracing.RunCreate{
		ID: runID, TraceID: runID,
		Name: "transformed-run", RunType: "chain",
		Inputs:      map[string]any{"q": "hello"},
		Extra:       map[string]any{"metadata": map[string]any{"source": "test"}},
		StartTime:   now,
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("CreateRun: %v", err)
	}

	if err := client.UpdateRun(&langsmithtracing.RunUpdate{
		ID: runID, TraceID: runID,
		Outputs:     map[string]any{"a": "world"},
		EndTime:     now.Add(50 * time.Millisecond),
		DottedOrder: dotted,
	}); err != nil {
		t.Fatalf("UpdateRun: %v", err)
	}

	client.Close()

	mu.Lock()
	defer mu.Unlock()

	// Parse all multipart bodies and collect extra parts.
	allParts := make(map[string][]byte)
	for i, body := range capturedBodies {
		decompressed := zstdDecompress(t, body)
		_, params, err := mime.ParseMediaType(capturedContentTypes[i])
		if err != nil {
			t.Fatalf("parse content-type %d: %v", i, err)
		}
		reader := multipart.NewReader(bytes.NewReader(decompressed), params["boundary"])
		for {
			p, err := reader.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				t.Fatalf("next part: %v", err)
			}
			data, _ := io.ReadAll(p)
			allParts[p.FormName()] = data
		}
	}

	// Find the extra part for the post (which should be coalesced).
	extraKey := "post." + runID.String() + ".extra"
	raw, ok := allParts[extraKey]
	if !ok {
		t.Fatalf("missing part %q", extraKey)
	}

	var extra map[string]any
	if err := json.Unmarshal(raw, &extra); err != nil {
		t.Fatalf("unmarshal extra: %v", err)
	}

	meta, _ := extra["metadata"].(map[string]any)
	if meta == nil {
		t.Fatal("extra.metadata is nil")
	}
	if meta["transformed"] != true {
		t.Errorf("expected extra.metadata.transformed=true, got %v", meta["transformed"])
	}
	if meta["hook_version"] != "1.0" {
		t.Errorf("expected extra.metadata.hook_version=1.0, got %v", meta["hook_version"])
	}
	t.Logf("Transform hook injected: transformed=%v, hook_version=%v", meta["transformed"], meta["hook_version"])
}

func formatDottedOrder(t time.Time, id uuid.UUID) string {
	return fmt.Sprintf("%s%06dZ%s",
		t.UTC().Format("20060102T150405"),
		t.UTC().Nanosecond()/1000,
		id.String(),
	)
}
