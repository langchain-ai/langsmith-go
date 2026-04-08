package proxy_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/langchain-ai/langsmith-go/instrumentation/proxy"
	"github.com/langchain-ai/langsmith-go/internal/mockllm"
)

type testEnv struct {
	ProxyURL string
	Exporter *tracetest.InMemoryExporter
	TP       *sdktrace.TracerProvider
	cleanup  func()
}

func (e *testEnv) Close() { e.cleanup() }

func (e *testEnv) FlushAndGetSpans(t *testing.T) []sdktrace.ReadOnlySpan {
	t.Helper()
	e.TP.ForceFlush(context.Background())
	spans := e.Exporter.GetSpans()
	ro := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		ro[i] = spans[i].Snapshot()
	}
	return ro
}

func setup(t *testing.T) *testEnv {
	t.Helper()

	upstream := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.DefaultHandler()))

	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	p, err := proxy.New(proxy.Config{
		AnthropicUpstream: upstream.URL(),
		OpenAIUpstream:    upstream.URL(),
		TracerProvider:    tp,
	})
	if err != nil {
		upstream.Close()
		t.Fatalf("proxy.New: %v", err)
	}

	proxyServer := httptest.NewServer(p)

	return &testEnv{
		ProxyURL: proxyServer.URL,
		Exporter: exporter,
		TP:       tp,
		cleanup: func() {
			proxyServer.Close()
			upstream.Close()
			tp.Shutdown(context.Background())
		},
	}
}

func spanAttr(spans []sdktrace.ReadOnlySpan, key string) (string, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes() {
			if string(attr.Key) == key {
				return attr.Value.Emit(), true
			}
		}
	}
	return "", false
}

func spanAttrInt(spans []sdktrace.ReadOnlySpan, key string) (int64, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes() {
			if string(attr.Key) == key {
				return attr.Value.AsInt64(), true
			}
		}
	}
	return 0, false
}

// --- Anthropic ---

func TestProxy_Anthropic_NonStreaming(t *testing.T) {
	env := setup(t)
	defer env.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hello"}],"max_tokens":100}`
	resp, err := http.Post(env.ProxyURL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}

	spans := env.FlushAndGetSpans(t)
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	if v, ok := spanAttr(spans, "gen_ai.system"); !ok || v != "anthropic" {
		t.Errorf("gen_ai.system = %q, want anthropic", v)
	}
	if _, ok := spanAttr(spans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	}
	if _, ok := spanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	}
	if v, ok := spanAttrInt(spans, "gen_ai.usage.input_tokens"); !ok || v == 0 {
		t.Errorf("gen_ai.usage.input_tokens = %d", v)
	}
}

func TestProxy_Anthropic_Streaming(t *testing.T) {
	env := setup(t)
	defer env.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hello"}],"max_tokens":100,"stream":true}`
	resp, err := http.Post(env.ProxyURL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if !strings.Contains(string(data), "text_delta") {
		t.Error("expected SSE text_delta events in response")
	}

	spans := env.FlushAndGetSpans(t)
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	if v, ok := spanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v, "Hello") {
		t.Errorf("completion should contain 'Hello': %s", v)
	}
}

// --- OpenAI ---

func TestProxy_OpenAI_NonStreaming(t *testing.T) {
	env := setup(t)
	defer env.Close()

	body := `{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}]}`
	resp, err := http.Post(env.ProxyURL+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d", resp.StatusCode)
	}

	spans := env.FlushAndGetSpans(t)
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	if v, ok := spanAttr(spans, "gen_ai.system"); !ok || v != "openai" {
		t.Errorf("gen_ai.system = %q, want openai", v)
	}
	if _, ok := spanAttr(spans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	}
	if _, ok := spanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	}
}

func TestProxy_OpenAI_Streaming(t *testing.T) {
	env := setup(t)
	defer env.Close()

	body := `{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}],"stream":true}`
	resp, err := http.Post(env.ProxyURL+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if !strings.Contains(string(data), "[DONE]") {
		t.Error("expected [DONE] in SSE stream")
	}

	spans := env.FlushAndGetSpans(t)
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	if v, ok := spanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v, "Hello") {
		t.Errorf("completion should contain 'Hello': %s", v)
	}
}

func TestProxy_OpenAI_ToolCalls(t *testing.T) {
	env := setup(t)
	defer env.Close()

	body := `{"model":"gpt-4o","messages":[{"role":"user","content":"weather?"}],"tools":[{"type":"function","function":{"name":"get_weather"}}]}`
	resp, err := http.Post(env.ProxyURL+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	spans := env.FlushAndGetSpans(t)
	if v, ok := spanAttr(spans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v, "get_weather") {
		t.Errorf("completion should contain tool call: %s", v)
	}
}

// --- Passthrough ---

func TestProxy_Models_Passthrough(t *testing.T) {
	env := setup(t)
	defer env.Close()

	resp, err := http.Get(env.ProxyURL + "/v1/models")
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	data, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	var result map[string]any
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	if result["object"] != "list" {
		t.Errorf("expected models list, got %v", result["object"])
	}

	// No tracing spans for /v1/models
	spans := env.FlushAndGetSpans(t)
	for _, s := range spans {
		if strings.Contains(s.Name(), "model") {
			t.Error("unexpected tracing span for /v1/models")
		}
	}
}

// --- Auth headers ---

func TestProxy_AuthHeaders_Forwarded(t *testing.T) {
	// Use a custom upstream that captures headers
	var capturedHeaders http.Header
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		capturedHeaders = r.Header.Clone()
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]any{
			"id": "msg_001", "type": "message", "role": "assistant", "model": "test",
			"content":     []map[string]any{{"type": "text", "text": "ok"}},
			"usage":       map[string]any{"input_tokens": 1, "output_tokens": 1},
			"stop_reason": "end_turn",
		})
	}))
	defer upstream.Close()

	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	defer tp.Shutdown(context.Background())

	p, err := proxy.New(proxy.Config{
		AnthropicUpstream: upstream.URL,
		OpenAIUpstream:    upstream.URL,
		TracerProvider:    tp,
	})
	if err != nil {
		t.Fatalf("proxy.New: %v", err)
	}
	proxyServer := httptest.NewServer(p)
	defer proxyServer.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":10}`
	req, _ := http.NewRequest("POST", proxyServer.URL+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "sk-ant-test-key-123")
	req.Header.Set("Authorization", "Bearer sk-test-key-456")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if capturedHeaders.Get("X-Api-Key") != "sk-ant-test-key-123" {
		t.Errorf("x-api-key not forwarded: %q", capturedHeaders.Get("X-Api-Key"))
	}
	if !strings.Contains(capturedHeaders.Get("Authorization"), "sk-test-key-456") {
		t.Errorf("authorization not forwarded: %q", capturedHeaders.Get("Authorization"))
	}
}

// --- Error tracing ---

func TestProxy_Error_Traced(t *testing.T) {
	// Upstream returns 401
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(map[string]any{
			"type": "error",
			"error": map[string]any{
				"type":    "authentication_error",
				"message": "invalid key",
			},
		})
	}))
	defer upstream.Close()

	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	defer tp.Shutdown(context.Background())

	p, _ := proxy.New(proxy.Config{
		AnthropicUpstream: upstream.URL,
		TracerProvider:    tp,
	})
	proxyServer := httptest.NewServer(p)
	defer proxyServer.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":10}`
	resp, err := http.Post(proxyServer.URL+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("expected 401, got %d", resp.StatusCode)
	}

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	foundError := false
	for _, s := range spans {
		snap := s.Snapshot()
		if snap.Status().Code == 2 { // codes.Error
			foundError = true
			break
		}
		for _, ev := range snap.Events() {
			if ev.Name == "exception" {
				foundError = true
				break
			}
		}
	}
	if !foundError {
		t.Error("expected error span for 401 response")
	}
}
