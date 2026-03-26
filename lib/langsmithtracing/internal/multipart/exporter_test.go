package multipart

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/klauspost/compress/zstd"
	"github.com/google/uuid"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/models"
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

func makeOp(kind models.OpKind) *models.SerializedOp {
	id := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	info, _ := json.Marshal(map[string]any{"id": id.String(), "name": "test"})
	inputs, _ := json.Marshal(map[string]any{"q": "hello"})
	return &models.SerializedOp{
		Kind:    kind,
		ID:      id,
		TraceID: id,
		RunInfo: info,
		Inputs:  inputs,
	}
}

func TestExporter_AttachmentParts(t *testing.T) {
	var captured []byte
	var contentType string
	var contentEncoding string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		contentEncoding = r.Header.Get("Content-Encoding")
		captured, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)

	runID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	traceID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	ops := []*models.SerializedOp{
		{
			Kind:    models.OpKindPost,
			ID:      runID,
			TraceID: traceID,
			RunInfo: []byte(`{"id":"11111111-1111-1111-1111-111111111111","name":"test"}`),
			Inputs:  []byte(`{"question":"hello"}`),
			Attachments: map[string]models.Attachment{
				"my_image": {
					ContentType: "image/png",
					Data:        []byte("fake-png-data"),
				},
				"report": {
					ContentType: "application/pdf",
					Data:        []byte("fake-pdf-data"),
				},
			},
		},
	}

	endpoint := models.WriteEndpoint{
		URL:     srv.URL,
		Key:     "test-key",
		Project: "test-project",
	}

	if err := exp.Export(context.Background(), endpoint, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	if contentEncoding != "zstd" {
		t.Fatalf("expected Content-Encoding: zstd, got %q", contentEncoding)
	}

	decompressed := zstdDecompress(t, captured)

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("parse content-type: %v", err)
	}
	boundary := params["boundary"]
	if boundary == "" {
		t.Fatal("missing boundary")
	}

	reader := multipart.NewReader(bytes.NewReader(decompressed), boundary)

	partsByName := map[string]partInfo{}
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("next part: %v", err)
		}
		data, _ := io.ReadAll(p)
		name := p.FormName()
		partsByName[name] = partInfo{
			contentType: p.Header.Get("Content-Type"),
			data:        data,
		}
	}

	runIDStr := runID.String()

	if _, ok := partsByName["post."+runIDStr]; !ok {
		t.Error("missing post.<runID> part")
	}
	if _, ok := partsByName["post."+runIDStr+".inputs"]; !ok {
		t.Error("missing post.<runID>.inputs part")
	}

	imgKey := "attachment." + runIDStr + ".my_image"
	if p, ok := partsByName[imgKey]; !ok {
		t.Errorf("missing attachment part %q", imgKey)
	} else {
		if !strings.HasPrefix(p.contentType, "image/png") {
			t.Errorf("image part content-type = %q, want prefix image/png", p.contentType)
		}
		if string(p.data) != "fake-png-data" {
			t.Errorf("image data = %q", string(p.data))
		}
	}

	pdfKey := "attachment." + runIDStr + ".report"
	if p, ok := partsByName[pdfKey]; !ok {
		t.Errorf("missing attachment part %q", pdfKey)
	} else {
		if !strings.HasPrefix(p.contentType, "application/pdf") {
			t.Errorf("pdf part content-type = %q, want prefix application/pdf", p.contentType)
		}
		if string(p.data) != "fake-pdf-data" {
			t.Errorf("pdf data = %q", string(p.data))
		}
	}

	t.Logf("Total parts received: %d", len(partsByName))
	for name := range partsByName {
		t.Logf("  part: %s", name)
	}
}

func TestExporter_NoAttachments(t *testing.T) {
	var captured []byte
	var contentType string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		captured, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)

	runID := uuid.New()
	ops := []*models.SerializedOp{
		{
			Kind:    models.OpKindPost,
			ID:      runID,
			TraceID: runID,
			RunInfo: []byte(`{"id":"test","name":"test"}`),
		},
	}

	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	if err := exp.Export(context.Background(), endpoint, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	decompressed := zstdDecompress(t, captured)

	_, params, _ := mime.ParseMediaType(contentType)
	reader := multipart.NewReader(bytes.NewReader(decompressed), params["boundary"])

	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("next part: %v", err)
		}
		name := p.FormName()
		if strings.HasPrefix(name, "attachment.") {
			t.Errorf("unexpected attachment part: %s", name)
		}
	}
}

func TestExporter_FallbackToBatchOn404(t *testing.T) {
	var multipartCalls, batchCalls int
	var batchBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/runs/multipart":
			multipartCalls++
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("not found"))
		case "/runs/batch":
			batchCalls++
			batchBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusBadRequest)
		}
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}

	runID := uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa")
	ops := []*models.SerializedOp{
		{
			Kind:    models.OpKindPost,
			ID:      runID,
			TraceID: runID,
			RunInfo: []byte(`{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","name":"test","run_type":"chain"}`),
			Inputs:  []byte(`{"question":"hello"}`),
		},
	}

	if err := exp.Export(context.Background(), endpoint, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	if multipartCalls != 1 {
		t.Errorf("expected 1 multipart call, got %d", multipartCalls)
	}
	if batchCalls != 1 {
		t.Errorf("expected 1 batch call, got %d", batchCalls)
	}

	var parsed struct {
		Post  []map[string]any `json:"post"`
		Patch []map[string]any `json:"patch"`
	}
	if err := json.Unmarshal(batchBody, &parsed); err != nil {
		t.Fatalf("unmarshal batch body: %v", err)
	}
	if len(parsed.Post) != 1 {
		t.Fatalf("expected 1 post, got %d", len(parsed.Post))
	}
	if parsed.Post[0]["name"] != "test" {
		t.Errorf("post name = %v, want %q", parsed.Post[0]["name"], "test")
	}
	if inputs, ok := parsed.Post[0]["inputs"].(map[string]any); !ok || inputs["question"] != "hello" {
		t.Errorf("post inputs = %v, want {question: hello}", parsed.Post[0]["inputs"])
	}
}

func TestExporter_MultipartDisabledSkipsMultipart(t *testing.T) {
	var multipartCalls, batchCalls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/runs/multipart":
			multipartCalls++
			w.WriteHeader(http.StatusNotFound)
		case "/runs/batch":
			batchCalls++
			w.WriteHeader(http.StatusAccepted)
		}
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}

	runID := uuid.New()
	mkOps := func() []*models.SerializedOp {
		return []*models.SerializedOp{{
			Kind:    models.OpKindPost,
			ID:      runID,
			TraceID: runID,
			RunInfo: []byte(`{"id":"test","name":"test"}`),
		}}
	}

	// First call: multipart → 404 → batch fallback.
	if err := exp.Export(context.Background(), endpoint, mkOps()); err != nil {
		t.Fatalf("first Export: %v", err)
	}
	if multipartCalls != 1 {
		t.Errorf("after first export: multipart calls = %d, want 1", multipartCalls)
	}
	if batchCalls != 1 {
		t.Errorf("after first export: batch calls = %d, want 1", batchCalls)
	}

	// Second call: should go directly to batch (no multipart attempt).
	if err := exp.Export(context.Background(), endpoint, mkOps()); err != nil {
		t.Fatalf("second Export: %v", err)
	}
	if multipartCalls != 1 {
		t.Errorf("after second export: multipart calls = %d, want 1 (should not retry)", multipartCalls)
	}
	if batchCalls != 2 {
		t.Errorf("after second export: batch calls = %d, want 2", batchCalls)
	}
}

func TestExporter_BatchPostAndPatch(t *testing.T) {
	var batchBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/runs/multipart":
			w.WriteHeader(http.StatusNotFound)
		case "/runs/batch":
			batchBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusAccepted)
		}
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}

	postID := uuid.MustParse("11111111-1111-1111-1111-111111111111")
	patchID := uuid.MustParse("22222222-2222-2222-2222-222222222222")

	ops := []*models.SerializedOp{
		{
			Kind:    models.OpKindPost,
			ID:      postID,
			TraceID: postID,
			RunInfo: []byte(`{"id":"11111111-1111-1111-1111-111111111111","name":"create-run"}`),
			Inputs:  []byte(`{"x":1}`),
		},
		{
			Kind:    models.OpKindPatch,
			ID:      patchID,
			TraceID: patchID,
			RunInfo: []byte(`{"id":"22222222-2222-2222-2222-222222222222"}`),
			Outputs: []byte(`{"y":2}`),
			Error:   []byte(`"something broke"`),
		},
	}

	if err := exp.Export(context.Background(), endpoint, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	var parsed struct {
		Post  []map[string]any `json:"post"`
		Patch []map[string]any `json:"patch"`
	}
	if err := json.Unmarshal(batchBody, &parsed); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if len(parsed.Post) != 1 {
		t.Fatalf("post count = %d, want 1", len(parsed.Post))
	}
	if len(parsed.Patch) != 1 {
		t.Fatalf("patch count = %d, want 1", len(parsed.Patch))
	}
	if parsed.Post[0]["name"] != "create-run" {
		t.Errorf("post name = %v", parsed.Post[0]["name"])
	}
	if parsed.Patch[0]["error"] != "something broke" {
		t.Errorf("patch error = %v", parsed.Patch[0]["error"])
	}
	if outputs, ok := parsed.Patch[0]["outputs"].(map[string]any); !ok || outputs["y"] != float64(2) {
		t.Errorf("patch outputs = %v", parsed.Patch[0]["outputs"])
	}
}

func TestExporter_NonNotFoundErrorDoesNotFallback(t *testing.T) {
	var batchCalls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/runs/multipart":
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("internal error"))
		case "/runs/batch":
			batchCalls++
			w.WriteHeader(http.StatusAccepted)
		}
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}

	ops := []*models.SerializedOp{{
		Kind:    models.OpKindPost,
		ID:      uuid.New(),
		TraceID: uuid.New(),
		RunInfo: []byte(`{"id":"test","name":"test"}`),
	}}

	err := exp.Export(context.Background(), endpoint, ops)
	if err == nil {
		t.Fatal("expected error from 500, got nil")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T", err)
	}
	if apiErr.StatusCode != 500 {
		t.Errorf("status = %d, want 500", apiErr.StatusCode)
	}
	if batchCalls != 0 {
		t.Errorf("batch should not be called for non-404 errors, got %d calls", batchCalls)
	}
}

func TestExporter_EmptyPartsAreSent(t *testing.T) {
	partsByName := exportAndParseParts(t, []*models.SerializedOp{
		{
			Kind:    models.OpKindPost,
			ID:      uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			TraceID: uuid.MustParse("aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"),
			RunInfo: []byte(`{"id":"aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa","name":"test"}`),
			Inputs:  []byte(`{}`),    // empty map, explicitly provided
			Outputs: []byte(`{}`),    // empty map, explicitly provided
			Events:  nil,             // not provided → should be skipped
			Extra:   []byte(`{}`),    // empty map, explicitly provided
			Error:   nil,             // not provided → should be skipped
		},
	})

	id := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"

	// Empty-but-provided fields must produce parts.
	for _, suffix := range []string{".inputs", ".outputs", ".extra"} {
		key := "post." + id + suffix
		p, ok := partsByName[key]
		if !ok {
			t.Errorf("expected part %q to be sent (empty-but-not-nil), but it was missing", key)
			continue
		}
		if string(p.data) != "{}" {
			t.Errorf("part %q data = %q, want %q", key, string(p.data), "{}")
		}
	}

	// Nil fields must NOT produce parts.
	for _, suffix := range []string{".events", ".error"} {
		key := "post." + id + suffix
		if _, ok := partsByName[key]; ok {
			t.Errorf("expected part %q to be skipped (nil), but it was present", key)
		}
	}
}

// exportAndParseParts is a test helper that exports ops through a local server,
// decompresses the zstd stream, and returns the parsed multipart parts by name.
func exportAndParseParts(t *testing.T, ops []*models.SerializedOp) map[string]partInfo {
	t.Helper()

	var captured []byte
	var contentType string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		contentType = r.Header.Get("Content-Type")
		captured, _ = io.ReadAll(r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	endpoint := models.WriteEndpoint{URL: srv.URL, Key: "k", Project: "p"}
	if err := exp.Export(context.Background(), endpoint, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	decompressed := zstdDecompress(t, captured)

	_, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		t.Fatalf("parse content-type: %v", err)
	}

	reader := multipart.NewReader(bytes.NewReader(decompressed), params["boundary"])
	result := map[string]partInfo{}
	for {
		p, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("next part: %v", err)
		}
		data, _ := io.ReadAll(p)
		result[p.FormName()] = partInfo{
			contentType: p.Header.Get("Content-Type"),
			data:        data,
		}
	}
	return result
}

type partInfo struct {
	contentType string
	data        []byte
}

func TestExporter_MultipartRetriesOn500(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after retries, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
}

func TestExporter_MultipartRetriesExhausted(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		mu.Unlock()
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("always failing"))
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	err := exp.Export(context.Background(), ep, ops)
	if err == nil {
		t.Fatal("expected error after exhausting retries")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 500 {
		t.Fatalf("expected status 500, got %d", apiErr.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 3 {
		t.Fatalf("expected 3 attempts, got %d", calls)
	}
}

func TestExporter_NoRetryOn400(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		mu.Unlock()
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("bad request"))
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	err := exp.Export(context.Background(), ep, ops)
	if err == nil {
		t.Fatal("expected error for 400")
	}
	apiErr, ok := err.(*APIError)
	if !ok {
		t.Fatalf("expected *APIError, got %T: %v", err, err)
	}
	if apiErr.StatusCode != 400 {
		t.Fatalf("expected status 400, got %d", apiErr.StatusCode)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected 1 attempt (no retry on 400), got %d", calls)
	}
}

func TestExporter_RetriesOn408(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(http.StatusRequestTimeout)
			w.Write([]byte("timeout"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 attempts, got %d", calls)
	}
}

func TestExporter_BatchRetriesOn500(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.ReadAll(r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("server error"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)
	exp.multipartDisabled.Store(true)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after batch retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 batch attempts, got %d", calls)
	}
}

func TestExporter_RetriesOnConnectionError(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(&http.Client{Timeout: 120 * time.Second}, rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}

	// First call to unreachable host triggers connection errors that get retried,
	// then exhausts all attempts.
	ep := models.WriteEndpoint{URL: "http://127.0.0.1:1", Key: "k"}
	err := exp.Export(context.Background(), ep, ops)
	if err == nil {
		t.Fatal("expected error for unreachable host")
	}
	// Should not be an *APIError since the connection never succeeded.
	if _, ok := err.(*APIError); ok {
		t.Fatal("expected non-APIError for connection failure")
	}
}

func TestExporter_BatchSplitsOversizedPayload(t *testing.T) {
	var mu sync.Mutex
	var receivedBodies []json.RawMessage

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		body, _ := io.ReadAll(r.Body)
		r.Body.Close()
		mu.Lock()
		receivedBodies = append(receivedBodies, json.RawMessage(body))
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 1, BackoffBase: time.Millisecond, BackoffMax: time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)
	exp.multipartDisabled.Store(true)
	exp.batchSizeLimitBytes = 500 // very small limit to force splitting

	// Create 5 ops with ~150 bytes each in their JSON form — total ~750 bytes,
	// which exceeds the 500-byte limit and should produce multiple chunks.
	var ops []*models.SerializedOp
	for i := 0; i < 5; i++ {
		id := uuid.New()
		info, _ := json.Marshal(map[string]any{
			"id":   id.String(),
			"name": "run-" + id.String()[:8],
		})
		inputs, _ := json.Marshal(map[string]any{
			"data": strings.Repeat("x", 50),
		})
		ops = append(ops, &models.SerializedOp{
			Kind:    models.OpKindPost,
			ID:      id,
			TraceID: id,
			RunInfo: info,
			Inputs:  inputs,
		})
	}

	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}
	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if len(receivedBodies) < 2 {
		t.Fatalf("expected at least 2 batch requests due to size splitting, got %d", len(receivedBodies))
	}

	totalPost := 0
	for i, body := range receivedBodies {
		var parsed struct {
			Post  []json.RawMessage `json:"post"`
			Patch []json.RawMessage `json:"patch"`
		}
		if err := json.Unmarshal(body, &parsed); err != nil {
			t.Fatalf("unmarshal chunk %d: %v", i, err)
		}
		totalPost += len(parsed.Post)
		t.Logf("chunk %d: %d post, %d patch, %d bytes", i+1, len(parsed.Post), len(parsed.Patch), len(body))
	}

	if totalPost != 5 {
		t.Fatalf("expected 5 total post ops across all chunks, got %d", totalPost)
	}
}

func TestExporter_BatchNoSplitWhenUnderLimit(t *testing.T) {
	var mu sync.Mutex
	var requestCount int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.ReadAll(r.Body)
		r.Body.Close()
		mu.Lock()
		requestCount++
		mu.Unlock()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 1, BackoffBase: time.Millisecond, BackoffMax: time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)
	exp.multipartDisabled.Store(true)
	// default 20MiB limit — a small batch should fit in one request

	ops := []*models.SerializedOp{makeOp(models.OpKindPost), makeOp(models.OpKindPatch)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if requestCount != 1 {
		t.Fatalf("expected 1 batch request (all ops fit), got %d", requestCount)
	}
}

func TestExporter_CompressionDisabledEnvVar(t *testing.T) {
	t.Setenv("LANGSMITH_DISABLE_RUN_COMPRESSION", "true")

	var mu sync.Mutex
	var receivedEncoding string
	var receivedBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedEncoding = r.Header.Get("Content-Encoding")
		receivedBody, _ = io.ReadAll(r.Body)
		mu.Unlock()
		r.Body.Close()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, true, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if receivedEncoding != "" {
		t.Errorf("expected no Content-Encoding header, got %q", receivedEncoding)
	}

	// Body should be raw multipart (not zstd-compressed), so it should start
	// with the boundary delimiter "--".
	if !strings.HasPrefix(string(receivedBody), "--") {
		t.Errorf("expected raw multipart body starting with '--', got prefix %q", string(receivedBody[:min(20, len(receivedBody))]))
	}

	// Extract the boundary from the first line of the body.
	idx := bytes.Index(receivedBody, []byte("\r\n"))
	if idx < 0 {
		t.Fatal("no CRLF found in body; cannot extract boundary")
	}
	boundary := strings.TrimPrefix(string(receivedBody[:idx]), "--")
	reader := multipart.NewReader(bytes.NewReader(receivedBody), boundary)
	partCount := 0
	for {
		_, err := reader.NextPart()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("failed to parse multipart: %v", err)
		}
		partCount++
	}
	if partCount == 0 {
		t.Error("expected at least 1 multipart part in uncompressed body")
	}
	t.Logf("Uncompressed body has %d parts", partCount)
}

func TestExporter_CompressionEnabledByDefault(t *testing.T) {
	var mu sync.Mutex
	var receivedEncoding string
	var receivedBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		receivedEncoding = r.Header.Get("Content-Encoding")
		receivedBody, _ = io.ReadAll(r.Body)
		mu.Unlock()
		r.Body.Close()
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()

	if receivedEncoding != "zstd" {
		t.Errorf("expected Content-Encoding: zstd, got %q", receivedEncoding)
	}

	// Body should be zstd-compressed; verify by decompressing.
	decompressed := zstdDecompress(t, receivedBody)
	if !strings.HasPrefix(string(decompressed), "--") {
		t.Errorf("decompressed body should be multipart starting with '--', got prefix %q",
			string(decompressed[:min(20, len(decompressed))]))
	}
}

func TestExporter_RetriesOn429WithRetryAfter(t *testing.T) {
	var mu sync.Mutex
	var calls int
	var firstCallTime, secondCallTime time.Time

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		if n == 1 {
			firstCallTime = time.Now()
		} else if n == 2 {
			secondCallTime = time.Now()
		}
		mu.Unlock()
		if n == 1 {
			w.Header().Set("Retry-After", "1")
			w.WriteHeader(http.StatusTooManyRequests)
			w.Write([]byte("rate limited"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after 429 retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 attempts, got %d", calls)
	}

	gap := secondCallTime.Sub(firstCallTime)
	if gap < 900*time.Millisecond {
		t.Errorf("expected ~1s Retry-After delay, got %v", gap)
	}
}

func TestExporter_429DefaultRetryAfterWithoutHeader(t *testing.T) {
	// Without a Retry-After header, 429 should use the default 10s delay.
	// Verify by checking the APIError.RetryAfter field directly.
	resp := &http.Response{
		StatusCode: 429,
		Header:     http.Header{},
	}
	d := parseRetryAfter(resp)
	if d != defaultRetryAfter429 {
		t.Errorf("expected default %v for 429 without header, got %v", defaultRetryAfter429, d)
	}

	// Non-429 without Retry-After should return 0.
	resp500 := &http.Response{
		StatusCode: 500,
		Header:     http.Header{},
	}
	if d := parseRetryAfter(resp500); d != 0 {
		t.Errorf("expected 0 for 500 without Retry-After, got %v", d)
	}
}

func TestExporter_RetriesOn502(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(http.StatusBadGateway)
			w.Write([]byte("bad gateway"))
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after 502 retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 attempts (502 then 202), got %d", calls)
	}
}

func TestExporter_RetriesOn503(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(http.StatusServiceUnavailable)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after 503 retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 attempts (503 then 202), got %d", calls)
	}
}

func TestExporter_NoRetryOn422(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		mu.Unlock()
		w.WriteHeader(http.StatusUnprocessableEntity)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	err := exp.Export(context.Background(), ep, ops)
	if err == nil {
		t.Fatal("expected error for 422")
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 1 {
		t.Fatalf("expected 1 attempt (no retry on 422), got %d", calls)
	}
}

func TestExporter_EmptyOpsNoOp(t *testing.T) {
	var calls int
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		calls++
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, nil); err != nil {
		t.Fatalf("Export(nil): %v", err)
	}
	if err := exp.Export(context.Background(), ep, []*models.SerializedOp{}); err != nil {
		t.Fatalf("Export(empty): %v", err)
	}
	if calls != 0 {
		t.Fatalf("expected 0 HTTP calls for empty ops, got %d", calls)
	}
}

func TestExporter_MissingRunInfoReturnsError(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	ops := []*models.SerializedOp{{
		Kind:    models.OpKindPost,
		ID:      uuid.New(),
		TraceID: uuid.New(),
		RunInfo: nil,
	}}

	err := exp.Export(context.Background(), ep, ops)
	if err == nil {
		t.Fatal("expected error for missing RunInfo, got nil")
	}
	if !strings.Contains(err.Error(), "missing run info") {
		t.Errorf("error should mention 'missing run info', got: %v", err)
	}
}

func TestExporter_EmptyContentTypeFallback(t *testing.T) {
	runID := uuid.MustParse("bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb")
	parts := exportAndParseParts(t, []*models.SerializedOp{{
		Kind:    models.OpKindPost,
		ID:      runID,
		TraceID: runID,
		RunInfo: []byte(`{"id":"bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb","name":"test"}`),
		Attachments: map[string]models.Attachment{
			"blob": {ContentType: "", Data: []byte("binary-data")},
		},
	}})

	key := "attachment." + runID.String() + ".blob"
	p, ok := parts[key]
	if !ok {
		t.Fatalf("missing attachment part %q", key)
	}
	if !strings.HasPrefix(p.contentType, "application/octet-stream") {
		t.Errorf("expected fallback content type application/octet-stream, got %q", p.contentType)
	}
	if string(p.data) != "binary-data" {
		t.Errorf("attachment data = %q, want %q", string(p.data), "binary-data")
	}
}

func TestExporter_BatchDropsAttachments(t *testing.T) {
	var batchBody []byte

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/runs/batch":
			batchBody, _ = io.ReadAll(r.Body)
			w.WriteHeader(http.StatusAccepted)
		default:
			w.WriteHeader(http.StatusNotFound)
		}
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	exp.multipartDisabled.Store(true)
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	runID := uuid.New()
	ops := []*models.SerializedOp{{
		Kind:    models.OpKindPost,
		ID:      runID,
		TraceID: runID,
		RunInfo: []byte(`{"id":"test","name":"with-attachment"}`),
		Attachments: map[string]models.Attachment{
			"file": {ContentType: "text/plain", Data: []byte("file-data")},
		},
	}}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("Export: %v", err)
	}

	if !strings.Contains(string(batchBody), "with-attachment") {
		t.Fatal("batch body should contain the run data")
	}
	if strings.Contains(string(batchBody), "file-data") {
		t.Error("batch body should NOT contain attachment data")
	}
}

func TestExporter_ContextCancellationAbortsExport(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(5 * time.Second)
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	exp := NewExporter(srv.Client(), RetryConfig{MaxAttempts: 1}, false, nil)
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}
	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}

	ctx, cancel := context.WithCancel(context.Background())
	cancel()

	err := exp.Export(ctx, ep, ops)
	if err == nil {
		t.Fatal("expected error from cancelled context")
	}
}

func TestExporter_RetriesOn425(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(425) // Too Early
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after 425 retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 attempts (425 then 202), got %d", calls)
	}
}

func TestExporter_RetriesOn504(t *testing.T) {
	var mu sync.Mutex
	var calls int

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		io.Copy(io.Discard, r.Body)
		r.Body.Close()
		mu.Lock()
		calls++
		n := calls
		mu.Unlock()
		if n == 1 {
			w.WriteHeader(http.StatusGatewayTimeout)
			return
		}
		w.WriteHeader(http.StatusAccepted)
	}))
	defer srv.Close()

	rc := RetryConfig{MaxAttempts: 3, BackoffBase: 10 * time.Millisecond, BackoffMax: 50 * time.Millisecond}
	exp := NewExporter(srv.Client(), rc, false, nil)

	ops := []*models.SerializedOp{makeOp(models.OpKindPost)}
	ep := models.WriteEndpoint{URL: srv.URL, Key: "k"}

	if err := exp.Export(context.Background(), ep, ops); err != nil {
		t.Fatalf("expected success after 504 retry, got: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if calls != 2 {
		t.Fatalf("expected 2 attempts (504 then 202), got %d", calls)
	}
}

func TestRetryConfig_BackoffJitterBounds(t *testing.T) {
	rc := RetryConfig{
		MaxAttempts: 5,
		BackoffBase: 100 * time.Millisecond,
		BackoffMax:  1 * time.Second,
	}

	for attempt := 0; attempt < 4; attempt++ {
		expectedMax := float64(rc.BackoffBase) * float64(int(1)<<attempt)
		if expectedMax > float64(rc.BackoffMax) {
			expectedMax = float64(rc.BackoffMax)
		}

		for i := 0; i < 100; i++ {
			d := rc.backoff(attempt)
			if d < 0 {
				t.Fatalf("attempt %d: backoff returned negative duration %v", attempt, d)
			}
			if d > time.Duration(expectedMax) {
				t.Fatalf("attempt %d: backoff %v exceeds cap %v", attempt, d, time.Duration(expectedMax))
			}
		}
	}
}

func TestRetryConfig_RetryDelayUsesRetryAfter(t *testing.T) {
	rc := RetryConfig{
		MaxAttempts: 3,
		BackoffBase: 100 * time.Millisecond,
		BackoffMax:  1 * time.Second,
	}

	apiErr := &APIError{StatusCode: 429, RetryAfter: 5 * time.Second}
	d := rc.retryDelay(apiErr, 0)
	if d != 5*time.Second {
		t.Errorf("expected 5s from Retry-After, got %v", d)
	}

	d = rc.retryDelay(nil, 0)
	if d < 0 || d > rc.BackoffBase {
		t.Errorf("expected backoff in [0, %v], got %v", rc.BackoffBase, d)
	}
}
