package mockllm_test

import (
	"bytes"
	"context"
	"io"
	"net/http"
	"strings"
	"testing"

	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"

	"github.com/langchain-ai/langsmith-go/instrumentation/traceanthropic"
	"github.com/langchain-ai/langsmith-go/instrumentation/traceopenai"
	"github.com/langchain-ai/langsmith-go/internal/mockllm"
)

func newTestTP() (*sdktrace.TracerProvider, *tracetest.InMemoryExporter) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	return tp, exporter
}

func getSpanAttr(spans []sdktrace.ReadOnlySpan, key string) (string, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes() {
			if string(attr.Key) == key {
				return attr.Value.Emit(), true
			}
		}
	}
	return "", false
}

func getSpanAttrInt(spans []sdktrace.ReadOnlySpan, key string) (int64, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes() {
			if string(attr.Key) == key {
				return attr.Value.AsInt64(), true
			}
		}
	}
	return 0, false
}

// --- Anthropic through traceanthropic ---

func TestTracing_Anthropic_NonStreaming(t *testing.T) {
	srv := mockllm.NewAnthropicServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceanthropic.Client(
		traceanthropic.WithTracerProvider(tp),
		traceanthropic.WithTraceAllHosts(),
	)

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if len(roSpans) == 0 {
		t.Fatal("expected at least one span")
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.system"); !ok || v != "anthropic" {
		t.Errorf("gen_ai.system = %q, want anthropic", v)
	}
	if v, ok := getSpanAttr(roSpans, "gen_ai.request.model"); !ok || v != "claude-sonnet-4-20250514" {
		t.Errorf("gen_ai.request.model = %q", v)
	}
	if _, ok := getSpanAttr(roSpans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	}
	if _, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	}
	if v, ok := getSpanAttrInt(roSpans, "gen_ai.usage.input_tokens"); !ok || v == 0 {
		t.Errorf("gen_ai.usage.input_tokens = %d", v)
	}
	if v, ok := getSpanAttrInt(roSpans, "gen_ai.usage.output_tokens"); !ok || v == 0 {
		t.Errorf("gen_ai.usage.output_tokens = %d", v)
	}
}

func TestTracing_Anthropic_Streaming(t *testing.T) {
	srv := mockllm.NewAnthropicServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceanthropic.Client(
		traceanthropic.WithTracerProvider(tp),
		traceanthropic.WithTraceAllHosts(),
	)

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100,"stream":true}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if len(roSpans) == 0 {
		t.Fatal("expected at least one span")
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion from streaming")
	} else if !strings.Contains(v, "Hello") {
		t.Errorf("completion should contain 'Hello': %s", v)
	}
	if v, ok := getSpanAttr(roSpans, "gen_ai.response.model"); !ok || !strings.Contains(v, "claude") {
		t.Errorf("gen_ai.response.model = %q", v)
	}
}

func TestTracing_Anthropic_ToolUse(t *testing.T) {
	srv := mockllm.NewAnthropicServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceanthropic.Client(
		traceanthropic.WithTracerProvider(tp),
		traceanthropic.WithTraceAllHosts(),
	)

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"weather?"}],"max_tokens":100,"tools":[{"name":"get_weather"}]}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v, "tool_use") {
		t.Errorf("completion should contain tool_use: %s", v)
	}
}

func TestTracing_Anthropic_Error(t *testing.T) {
	srv := mockllm.NewAnthropicServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceanthropic.Client(
		traceanthropic.WithTracerProvider(tp),
		traceanthropic.WithTraceAllHosts(),
	)

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "sk-ant-invalid-test")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if len(roSpans) == 0 {
		t.Fatal("expected span even on error")
	}

	// Check span recorded an error event or has error status
	foundError := false
	for _, s := range roSpans {
		if s.Status().Code == 2 { // codes.Error
			foundError = true
			break
		}
		for _, ev := range s.Events() {
			if ev.Name == "exception" {
				foundError = true
				break
			}
		}
	}
	if !foundError {
		for _, s := range roSpans {
			t.Logf("span %q: status=%v events=%v", s.Name(), s.Status(), s.Events())
		}
		t.Error("expected at least one span with error status or error event")
	}
}

// --- OpenAI through traceopenai ---

func TestTracing_OpenAI_NonStreaming(t *testing.T) {
	srv := mockllm.NewOpenAIServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceopenai.Client(traceopenai.WithTracerProvider(tp))

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if len(roSpans) == 0 {
		t.Fatal("expected at least one span")
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.system"); !ok || v != "openai" {
		t.Errorf("gen_ai.system = %q, want openai", v)
	}
	if v, ok := getSpanAttr(roSpans, "gen_ai.request.model"); !ok || v != "gpt-4o-mini" {
		t.Errorf("gen_ai.request.model = %q", v)
	}
	if _, ok := getSpanAttr(roSpans, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	}
	if _, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	}
	if v, ok := getSpanAttrInt(roSpans, "gen_ai.usage.input_tokens"); !ok || v == 0 {
		t.Errorf("gen_ai.usage.input_tokens = %d", v)
	}
}

func TestTracing_OpenAI_Streaming(t *testing.T) {
	srv := mockllm.NewOpenAIServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceopenai.Client(traceopenai.WithTracerProvider(tp))

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}],"stream":true}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/chat/completions", bytes.NewReader([]byte(body)))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if len(roSpans) == 0 {
		t.Fatal("expected at least one span")
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion from streaming")
	} else if !strings.Contains(v, "Hello") {
		t.Errorf("completion should contain 'Hello': %s", v)
	}
}

func TestTracing_OpenAI_ToolCalls(t *testing.T) {
	srv := mockllm.NewOpenAIServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceopenai.Client(traceopenai.WithTracerProvider(tp))

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"weather?"}],"tools":[{"type":"function","function":{"name":"get_weather"}}]}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v, "get_weather") {
		t.Errorf("completion should contain get_weather: %s", v)
	}
}

func TestTracing_OpenAI_Responses(t *testing.T) {
	srv := mockllm.NewOpenAIServer()
	defer srv.Close()

	tp, exporter := newTestTP()
	client := traceopenai.Client(traceopenai.WithTracerProvider(tp))

	body := `{"model":"gpt-4o","input":"What is Go?"}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/responses", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")

	resp, err := client.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	tp.ForceFlush(context.Background())
	spans := exporter.GetSpans()

	roSpans := make([]sdktrace.ReadOnlySpan, len(spans))
	for i := range spans {
		roSpans[i] = spans[i].Snapshot()
	}

	if len(roSpans) == 0 {
		t.Fatal("expected at least one span")
	}

	if v, ok := getSpanAttr(roSpans, "gen_ai.operation.name"); !ok || v != "responses" {
		t.Errorf("gen_ai.operation.name = %q, want responses", v)
	}
	if _, ok := getSpanAttr(roSpans, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion from responses API")
	}
}
