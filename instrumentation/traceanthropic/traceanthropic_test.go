package traceanthropic

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

// fakeTransport returns a canned response, capturing the request for inspection.
type fakeTransport struct {
	body []byte
}

func (t *fakeTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusOK,
		Body:       io.NopCloser(bytes.NewReader(t.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

// hasEvent reports whether any exported span contains an event with the given name.
func hasEvent(spans tracetest.SpanStubs, name string) bool {
	for _, s := range spans {
		for _, e := range s.Events {
			if e.Name == name {
				return true
			}
		}
	}
	return false
}

func newTracedClient(t *testing.T, body []byte) (*http.Client, *tracetest.InMemoryExporter) {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	client := WrapClient(&http.Client{Transport: &fakeTransport{body: body}}, WithTracerProvider(tp))
	return client, exporter
}

func doStreamingMessages(t *testing.T, client *http.Client) {
	t.Helper()
	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":16,"stream":true}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.anthropic.com/v1/messages", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()
}

func TestRoundTrip_StreamingEmitsNewTokenOnFirstContentDelta(t *testing.T) {
	// message_start and content_block_start arrive before any text — only the
	// first content_block_delta should trip new_token.
	sse := "data: {\"type\":\"message_start\",\"message\":{\"model\":\"claude-sonnet-4-20250514\",\"usage\":{\"input_tokens\":3}}}\n" +
		"data: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\"}}\n" +
		"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"hi\"}}\n" +
		"data: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":1}}\n" +
		"data: {\"type\":\"message_stop\"}\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreamingMessages(t, client)

	if !hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("expected new_token event, got events=%v", exporter.GetSpans())
	}
}

func TestRoundTrip_StreamingPreambleOnlyDoesNotEmitNewToken(t *testing.T) {
	// Stream cancelled after content_block_start but before any content_block_delta.
	sse := "data: {\"type\":\"message_start\",\"message\":{\"model\":\"claude-sonnet-4-20250514\",\"usage\":{\"input_tokens\":3}}}\n" +
		"data: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\"}}\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreamingMessages(t, client)

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("preamble-only stream should not emit new_token, got %v", exporter.GetSpans())
	}
}

// errStatusTransport returns an HTTP error body with no SSE frames.
type errStatusTransport struct{ body []byte }

func (e *errStatusTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	return &http.Response{
		StatusCode: http.StatusTooManyRequests,
		Body:       io.NopCloser(bytes.NewReader(e.body)),
		Header:     make(http.Header),
		Request:    req,
	}, nil
}

func TestRoundTrip_StreamingHTTPErrorDoesNotEmitNewToken(t *testing.T) {
	// stream=true but the API returns an error JSON body — first-token time is
	// undefined and we must not record one.
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	client := WrapClient(&http.Client{
		Transport: &errStatusTransport{body: []byte(`{"type":"error","error":{"message":"overloaded"}}`)},
	}, WithTracerProvider(tp))
	doStreamingMessages(t, client)

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("streamed HTTP error should not emit new_token, got %v", exporter.GetSpans())
	}
}

func TestIsFirstContent(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{"message_start", `{"type":"message_start","message":{}}`, false},
		{"content_block_start", `{"type":"content_block_start","index":0}`, false},
		{"content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"}}`, true},
		{"message_delta", `{"type":"message_delta"}`, false},
		{"message_stop", `{"type":"message_stop"}`, false},
	}
	for _, tt := range tests {
		var c map[string]any
		_ = json.Unmarshal([]byte(tt.body), &c)
		if got := isFirstContent(c); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestRoundTrip_NonStreamingDoesNotEmitNewTokenEvent(t *testing.T) {
	respBody := `{"model":"claude-sonnet-4-20250514","content":[{"type":"text","text":"hi"}],"usage":{"input_tokens":3,"output_tokens":1},"stop_reason":"end_turn"}`
	client, exporter := newTracedClient(t, []byte(respBody))

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":16}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.anthropic.com/v1/messages", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("non-streaming span should not record a new_token event")
	}
}

func TestGetSpanName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/v1/messages", "anthropic.messages"},
		{"/v1/messages?beta=true", "anthropic.messages"},
		{"/v1/complete", "anthropic.request"},
		{"", "anthropic.request"},
	}
	for _, tt := range tests {
		if got := getSpanName(tt.path); got != tt.want {
			t.Errorf("getSpanName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestGetOperationName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/v1/messages", "chat"},
		{"/v1/complete", "request"},
	}
	for _, tt := range tests {
		if got := getOperationName(tt.path); got != tt.want {
			t.Errorf("getOperationName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func startTestSpan(t *testing.T) (sdktrace.ReadWriteSpan, *tracetest.InMemoryExporter) {
	t.Helper()
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	_, span := tp.Tracer("test").Start(context.Background(), "test")
	rwSpan, ok := span.(sdktrace.ReadWriteSpan)
	if !ok {
		t.Fatal("span does not implement ReadWriteSpan")
	}
	return rwSpan, exporter
}

func getAttr(span sdktrace.ReadWriteSpan, key string) (attribute.Value, bool) {
	for _, attr := range span.Attributes() {
		if string(attr.Key) == key {
			return attr.Value, true
		}
	}
	return attribute.Value{}, false
}

func TestExtractRequestAttributes_BasicMessage(t *testing.T) {
	span, _ := startTestSpan(t)
	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":1024,"temperature":0.7,"stream":false}`

	streaming := extractRequestAttributes(span, []byte(body))
	if streaming {
		t.Error("expected streaming=false")
	}

	if v, ok := getAttr(span, "gen_ai.request.model"); !ok || v.AsString() != "claude-sonnet-4-20250514" {
		t.Errorf("model attr: got %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.request.max_tokens"); !ok || v.AsInt64() != 1024 {
		t.Errorf("max_tokens attr: got %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.request.temperature"); !ok || v.AsFloat64() != 0.7 {
		t.Errorf("temperature attr: got %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt attribute")
	} else if !strings.Contains(v.AsString(), "hi") {
		t.Errorf("prompt should contain 'hi': %s", v.AsString())
	}
}

func TestExtractRequestAttributes_StreamingWithSystem(t *testing.T) {
	span, _ := startTestSpan(t)
	body := `{"model":"claude-3-haiku","system":"You are helpful","messages":[{"role":"user","content":"hello"}],"stream":true}`

	streaming := extractRequestAttributes(span, []byte(body))
	if !streaming {
		t.Error("expected streaming=true")
	}

	v, ok := getAttr(span, "gen_ai.prompt")
	if !ok {
		t.Fatal("expected gen_ai.prompt attribute")
	}
	prompt := v.AsString()
	if !strings.Contains(prompt, "system") {
		t.Errorf("prompt should include system message: %s", prompt)
	}
	if !strings.Contains(prompt, "You are helpful") {
		t.Errorf("prompt should include system content: %s", prompt)
	}
}

func TestExtractRequestAttributes_InvalidJSON(t *testing.T) {
	span, _ := startTestSpan(t)
	streaming := extractRequestAttributes(span, []byte("not json"))
	if streaming {
		t.Error("expected streaming=false for invalid JSON")
	}
}

func TestExtractResponseAttributes(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	body := `{
		"model": "claude-sonnet-4-20250514",
		"content": [{"type": "text", "text": "Hello!"}],
		"usage": {"input_tokens": 10, "output_tokens": 5},
		"stop_reason": "end_turn"
	}`

	extractResponseAttributes(span, []byte(body), parentSpan)

	if v, ok := getAttr(span, "gen_ai.response.model"); !ok || v.AsString() != "claude-sonnet-4-20250514" {
		t.Errorf("model attr: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(span, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion attribute")
	} else if !strings.Contains(v.AsString(), "Hello!") {
		t.Errorf("completion should contain 'Hello!': %s", v.AsString())
	}
	if v, ok := getAttr(span, "langsmith.metadata.stop_reason"); !ok || v.AsString() != "end_turn" {
		t.Errorf("stop_reason attr: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(span, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 10 {
		t.Errorf("input_tokens attr: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(span, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 5 {
		t.Errorf("output_tokens attr: got %v, ok=%v", v, ok)
	}
}

func TestExtractResponseAttributes_InvalidJSON(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)
	extractResponseAttributes(span, []byte("nope"), parentSpan)
	if _, ok := getAttr(span, "gen_ai.response.model"); ok {
		t.Error("should not set attributes for invalid JSON")
	}
}

func TestExtractStreamingResponseAttributes_TextContent(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	sse := "data: {\"type\":\"message_start\",\"message\":{\"model\":\"claude-sonnet-4-20250514\",\"usage\":{\"input_tokens\":15}}}\n" +
		"data: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"text\"}}\n" +
		"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\"Hello\"}}\n" +
		"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"text_delta\",\"text\":\" world\"}}\n" +
		"data: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"end_turn\"},\"usage\":{\"output_tokens\":8}}\n"

	extractStreamingResponseAttributes(span, []byte(sse), parentSpan)

	if v, ok := getAttr(span, "gen_ai.response.model"); !ok || v.AsString() != "claude-sonnet-4-20250514" {
		t.Errorf("model: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(span, "gen_ai.completion"); !ok {
		t.Error("expected gen_ai.completion")
	} else if !strings.Contains(v.AsString(), "Hello world") {
		t.Errorf("completion should contain 'Hello world': %s", v.AsString())
	}
	if v, ok := getAttr(span, "langsmith.metadata.stop_reason"); !ok || v.AsString() != "end_turn" {
		t.Errorf("stop_reason: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(span, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 15 {
		t.Errorf("input_tokens: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(span, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 8 {
		t.Errorf("output_tokens: got %v, ok=%v", v, ok)
	}
}

func TestExtractStreamingResponseAttributes_ToolUse(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	sse := "data: {\"type\":\"message_start\",\"message\":{\"model\":\"claude-sonnet-4-20250514\",\"usage\":{\"input_tokens\":20}}}\n" +
		"data: {\"type\":\"content_block_start\",\"index\":0,\"content_block\":{\"type\":\"tool_use\",\"id\":\"toolu_1\",\"name\":\"get_weather\"}}\n" +
		"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"input_json_delta\",\"partial_json\":\"{\\\"city\\\"\"}}\n" +
		"data: {\"type\":\"content_block_delta\",\"index\":0,\"delta\":{\"type\":\"input_json_delta\",\"partial_json\":\": \\\"SF\\\"}\"}}\n" +
		"data: {\"type\":\"message_delta\",\"delta\":{\"stop_reason\":\"tool_use\"},\"usage\":{\"output_tokens\":12}}\n"

	extractStreamingResponseAttributes(span, []byte(sse), parentSpan)

	v, ok := getAttr(span, "gen_ai.completion")
	if !ok {
		t.Fatal("expected gen_ai.completion")
	}
	completion := v.AsString()
	if !strings.Contains(completion, "tool_use") {
		t.Errorf("completion should contain tool_use: %s", completion)
	}
	if !strings.Contains(completion, "get_weather") {
		t.Errorf("completion should contain function name: %s", completion)
	}
}

func TestExtractStreamingResponseAttributes_Empty(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)
	extractStreamingResponseAttributes(span, []byte(""), parentSpan)
	if _, ok := getAttr(span, "gen_ai.completion"); ok {
		t.Error("should not set completion for empty data")
	}
}

func TestSetUsageAttributes(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	usage := map[string]interface{}{
		"input_tokens":                float64(100),
		"output_tokens":               float64(50),
		"cache_creation_input_tokens": float64(10),
		"cache_read_input_tokens":     float64(5),
	}
	setUsageAttributes(span, usage, parentSpan)

	if v, ok := getAttr(span, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 115 {
		t.Errorf("input_tokens: got %v (expected 100+10+5=115)", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 50 {
		t.Errorf("output_tokens: got %v", v)
	}
	if v, ok := getAttr(span, "langsmith.metadata.usage_metadata.input_token_details.cache_creation"); !ok || v.AsInt64() != 10 {
		t.Errorf("cache_creation: got %v", v)
	}
	if v, ok := getAttr(span, "langsmith.metadata.usage_metadata.input_token_details.cache_read"); !ok || v.AsInt64() != 5 {
		t.Errorf("cache_read: got %v", v)
	}
}

func TestSetUsageAttributes_PropagatesTokensToParent(t *testing.T) {
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	ctx, parentSpan := tp.Tracer("test").Start(context.Background(), "parent")
	_, childSpan := tp.Tracer("test").Start(ctx, "child")

	rwChild, ok := childSpan.(sdktrace.ReadWriteSpan)
	if !ok {
		t.Fatal("child span is not ReadWriteSpan")
	}
	rwParent, ok := parentSpan.(sdktrace.ReadWriteSpan)
	if !ok {
		t.Fatal("parent span is not ReadWriteSpan")
	}

	usage := map[string]interface{}{
		"input_tokens":  float64(10),
		"output_tokens": float64(5),
	}
	setUsageAttributes(rwChild, usage, rwParent)

	if v, ok := getAttr(rwParent, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 10 {
		t.Errorf("parent input_tokens: got %v, ok=%v", v, ok)
	}
	if v, ok := getAttr(rwParent, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 5 {
		t.Errorf("parent output_tokens: got %v, ok=%v", v, ok)
	}
}

func TestSetUsageAttributes_EmptyUsage(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)
	setUsageAttributes(span, map[string]interface{}{}, parentSpan)

	if _, ok := getAttr(span, "gen_ai.usage.input_tokens"); ok {
		t.Error("should not set input_tokens for empty usage")
	}
}

func TestExtractResponseAttributes_FullRoundTrip(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	resp := map[string]any{
		"model": "claude-sonnet-4-20250514",
		"content": []any{
			map[string]any{"type": "text", "text": "Result"},
			map[string]any{
				"type":  "tool_use",
				"id":    "toolu_1",
				"name":  "calculator",
				"input": map[string]any{"expression": "2+2"},
			},
		},
		"usage":       map[string]any{"input_tokens": float64(20), "output_tokens": float64(15)},
		"stop_reason": "tool_use",
	}
	body, _ := json.Marshal(resp)
	extractResponseAttributes(span, body, parentSpan)

	v, ok := getAttr(span, "gen_ai.completion")
	if !ok {
		t.Fatal("expected gen_ai.completion")
	}
	if !strings.Contains(v.AsString(), "Result") {
		t.Errorf("completion should contain 'Result': %s", v.AsString())
	}
	if !strings.Contains(v.AsString(), "calculator") {
		t.Errorf("completion should contain tool_use: %s", v.AsString())
	}
}
