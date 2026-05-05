package tracegemini

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

func doNonStreaming(t *testing.T, client *http.Client) {
	t.Helper()
	body := `{"contents":[{"role":"user","parts":[{"text":"hello"}]}]}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=test-key",
		strings.NewReader(body))
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

func doStreaming(t *testing.T, client *http.Client) {
	t.Helper()
	body := `{"contents":[{"role":"user","parts":[{"text":"hello"}]}]}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:streamGenerateContent?alt=sse&key=test-key",
		strings.NewReader(body))
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

func getSpanAttr(spans tracetest.SpanStubs, key string) (string, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes {
			if string(attr.Key) == key {
				return attr.Value.Emit(), true
			}
		}
	}
	return "", false
}

func getSpanAttrInt(spans tracetest.SpanStubs, key string) (int64, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes {
			if string(attr.Key) == key {
				return attr.Value.AsInt64(), true
			}
		}
	}
	return 0, false
}

// --- Endpoint detection ---

func TestIsGeminiEndpoint(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/v1beta/models/gemini-2.0-flash:generateContent", true},
		{"/v1beta/models/gemini-2.0-flash:streamGenerateContent", true},
		{"/v1beta/models/gemini-1.5-pro:generateContent", true},
		{"/v1beta/models/gemini-2.0-flash:countTokens", false},
		{"/v1beta/models/gemini-2.0-flash", false},
		{"/v1/chat/completions", false},
		// Vertex AI paths
		{"/v1beta1/publishers/google/models/gemini-2.0-flash:generateContent", true},
		{"/v1beta1/publishers/google/models/gemini-2.0-flash:streamGenerateContent", true},
		{"/v1/projects/my-proj/locations/us-central1/publishers/google/models/gemini-2.0-flash:generateContent", true},
		{"", false},
	}
	for _, tt := range tests {
		if got := isGeminiEndpoint(tt.path); got != tt.want {
			t.Errorf("isGeminiEndpoint(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestParseModelAction(t *testing.T) {
	tests := []struct {
		path       string
		wantModel  string
		wantAction string
	}{
		{"/v1beta/models/gemini-2.0-flash:generateContent", "gemini-2.0-flash", "generateContent"},
		{"/v1beta/models/gemini-1.5-pro:streamGenerateContent", "gemini-1.5-pro", "streamGenerateContent"},
		{"/v1beta/models/gemini-2.0-flash:countTokens", "gemini-2.0-flash", "countTokens"},
		{"/v1beta/models/gemini-2.0-flash", "gemini-2.0-flash", ""},
		{"/v1/models/gemini-2.0-flash:generateContent", "gemini-2.0-flash", "generateContent"},
		// Vertex AI paths
		{"/v1beta1/publishers/google/models/gemini-2.0-flash:generateContent", "gemini-2.0-flash", "generateContent"},
		{"/v1/projects/my-proj/locations/us-central1/publishers/google/models/gemini-2.0-flash:streamGenerateContent", "gemini-2.0-flash", "streamGenerateContent"},
		{"/v1beta/models/", "", ""},
		{"/no-models-here", "", ""},
	}
	for _, tt := range tests {
		model, action := parseModelAction(tt.path)
		if model != tt.wantModel || action != tt.wantAction {
			t.Errorf("parseModelAction(%q) = (%q, %q), want (%q, %q)",
				tt.path, model, action, tt.wantModel, tt.wantAction)
		}
	}
}

// --- Non-streaming ---

func TestRoundTrip_NonStreaming(t *testing.T) {
	respBody := `{
		"candidates": [{"content": {"parts": [{"text": "Hello!"}], "role": "model"}, "finishReason": "STOP"}],
		"modelVersion": "gemini-2.0-flash-001",
		"responseId": "resp-abc",
		"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 5}
	}`
	client, exporter := newTracedClient(t, []byte(respBody))
	doNonStreaming(t, client)

	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}

	if v, ok := getSpanAttr(spans, "gen_ai.provider.name"); !ok || v != "gcp.gemini" {
		t.Errorf("gen_ai.provider.name = %q, want 'gcp.gemini'", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.operation.name"); !ok || v != "generate_content" {
		t.Errorf("gen_ai.operation.name = %q, want 'generate_content'", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.request.model"); !ok || v != "gemini-2.0-flash" {
		t.Errorf("gen_ai.request.model = %q", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.response.model"); !ok || v != "gemini-2.0-flash-001" {
		t.Errorf("gen_ai.response.model = %q", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.response.id"); !ok || v != "resp-abc" {
		t.Errorf("gen_ai.response.id = %q", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.completion"); !ok || !strings.Contains(v, "Hello!") {
		t.Errorf("gen_ai.completion should contain 'Hello!': %q", v)
	} else if !strings.Contains(v, `"finish_reason":"STOP"`) {
		t.Errorf("gen_ai.completion should contain finish_reason: %q", v)
	}
	if v, ok := getSpanAttr(spans, "langsmith.metadata.stop_reason"); !ok || v != "STOP" {
		t.Errorf("stop_reason = %q, want 'STOP'", v)
	}
	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.input_tokens"); !ok || v != 10 {
		t.Errorf("input_tokens = %d, want 10", v)
	}
	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.output_tokens"); !ok || v != 5 {
		t.Errorf("output_tokens = %d, want 5", v)
	}
}

func TestRoundTrip_NonStreamingDoesNotEmitNewToken(t *testing.T) {
	respBody := `{
		"candidates": [{"content": {"parts": [{"text": "hi"}]}, "finishReason": "STOP"}],
		"usageMetadata": {"promptTokenCount": 1, "candidatesTokenCount": 1}
	}`
	client, exporter := newTracedClient(t, []byte(respBody))
	doNonStreaming(t, client)

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("non-streaming span should not record a new_token event")
	}
}

func TestRoundTrip_NonStreamingURLRedactsKey(t *testing.T) {
	respBody := `{"candidates": [{"content": {"parts": [{"text": "hi"}]}, "finishReason": "STOP"}]}`
	client, exporter := newTracedClient(t, []byte(respBody))
	doNonStreaming(t, client)

	spans := exporter.GetSpans()
	if v, ok := getSpanAttr(spans, "http.url"); !ok {
		t.Fatal("expected http.url attribute")
	} else {
		if strings.Contains(v, "test-key") {
			t.Errorf("http.url should not contain API key, got %q", v)
		}
		if !strings.Contains(v, "REDACTED") {
			t.Errorf("http.url should contain REDACTED, got %q", v)
		}
	}
}

// --- Streaming ---

func TestRoundTrip_StreamingEmitsNewTokenOnFirstContent(t *testing.T) {
	sse := "data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"Hi\"}],\"role\":\"model\"}}],\"modelVersion\":\"gemini-2.0-flash-001\"}\n\n" +
		"data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"!\"}],\"role\":\"model\"},\"finishReason\":\"STOP\"}],\"usageMetadata\":{\"promptTokenCount\":3,\"candidatesTokenCount\":2},\"modelVersion\":\"gemini-2.0-flash-001\"}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreaming(t, client)

	if !hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("expected new_token event on first content chunk")
	}
}

func TestRoundTrip_StreamingAccumulatesText(t *testing.T) {
	sse := "data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"Hello\"}],\"role\":\"model\"}}]}\n\n" +
		"data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\" world\"}],\"role\":\"model\"},\"finishReason\":\"STOP\"}],\"usageMetadata\":{\"promptTokenCount\":5,\"candidatesTokenCount\":2},\"modelVersion\":\"gemini-2.0-flash-001\"}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreaming(t, client)

	spans := exporter.GetSpans()
	if v, ok := getSpanAttr(spans, "gen_ai.completion"); !ok || !strings.Contains(v, "Hello world") {
		t.Errorf("completion should contain accumulated text 'Hello world': %q", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.operation.name"); !ok || v != "stream_generate_content" {
		t.Errorf("operation name = %q, want 'stream_generate_content'", v)
	}
	if v, ok := getSpanAttr(spans, "gen_ai.response.model"); !ok || v != "gemini-2.0-flash-001" {
		t.Errorf("response model = %q", v)
	}
	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.input_tokens"); !ok || v != 5 {
		t.Errorf("input_tokens = %d, want 5", v)
	}
	if v, ok := getSpanAttrInt(spans, "gen_ai.usage.output_tokens"); !ok || v != 2 {
		t.Errorf("output_tokens = %d, want 2", v)
	}
}

func TestRoundTrip_StreamingEmptyChunksNoNewToken(t *testing.T) {
	// Chunks with no candidates shouldn't trigger new_token
	sse := "data: {\"modelVersion\":\"gemini-2.0-flash-001\"}\n\n" +
		"data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"\"}],\"role\":\"model\"}}]}\n\n" +
		"data: {\"candidates\":[{\"finishReason\":\"STOP\"}],\"usageMetadata\":{\"promptTokenCount\":1,\"candidatesTokenCount\":0}}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreaming(t, client)

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("empty text chunks should not emit new_token")
	}
}

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
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	client := WrapClient(&http.Client{
		Transport: &errStatusTransport{body: []byte(`{"error":{"message":"rate limited"}}`)},
	}, WithTracerProvider(tp))
	doStreaming(t, client)

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("streamed HTTP error should not emit new_token")
	}
}

// --- Non-Gemini requests pass through ---

func TestRoundTrip_NonGeminiPassesThrough(t *testing.T) {
	client, exporter := newTracedClient(t, []byte(`{"ok":true}`))

	req, _ := http.NewRequestWithContext(context.Background(), http.MethodGet,
		"https://api.openai.com/v1/models", nil)
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if len(exporter.GetSpans()) > 0 {
		t.Errorf("non-Gemini request should not produce spans")
	}
}

// --- Run name ---

func TestRoundTrip_RunNameFromContext(t *testing.T) {
	respBody := `{"candidates": [{"content": {"parts": [{"text": "hi"}]}, "finishReason": "STOP"}]}`
	client, exporter := newTracedClient(t, []byte(respBody))

	body := `{"contents":[{"role":"user","parts":[{"text":"hello"}]}]}`
	ctx := WithRunNameContext(context.Background(), "my_custom_name")
	req, _ := http.NewRequestWithContext(ctx, http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent",
		strings.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected span")
	}
	if spans[0].Name != "my_custom_name" {
		t.Errorf("span name = %q, want 'my_custom_name'", spans[0].Name)
	}
}

func TestRoundTrip_RunNameFromOption(t *testing.T) {
	respBody := `{"candidates": [{"content": {"parts": [{"text": "hi"}]}, "finishReason": "STOP"}]}`
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	client := WrapClient(
		&http.Client{Transport: &fakeTransport{body: []byte(respBody)}},
		WithTracerProvider(tp),
		WithRunName("option_name"),
	)

	body := `{"contents":[{"role":"user","parts":[{"text":"hello"}]}]}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent",
		strings.NewReader(body))
	resp, _ := client.Do(req)
	io.ReadAll(resp.Body)
	resp.Body.Close()

	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected span")
	}
	if spans[0].Name != "option_name" {
		t.Errorf("span name = %q, want 'option_name'", spans[0].Name)
	}
}

// --- isFirstContent ---

func TestIsFirstContent(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{"text content", `{"candidates":[{"content":{"parts":[{"text":"hi"}]}}]}`, true},
		{"empty text", `{"candidates":[{"content":{"parts":[{"text":""}]}}]}`, false},
		{"function call", `{"candidates":[{"content":{"parts":[{"functionCall":{"name":"f","args":{}}}]}}]}`, true},
		{"no candidates", `{"modelVersion":"v1"}`, false},
		{"empty candidates", `{"candidates":[]}`, false},
		{"no content", `{"candidates":[{"finishReason":"STOP"}]}`, false},
		{"no parts", `{"candidates":[{"content":{}}]}`, false},
		{"empty parts", `{"candidates":[{"content":{"parts":[]}}]}`, false},
	}
	for _, tt := range tests {
		var c map[string]any
		_ = json.Unmarshal([]byte(tt.body), &c)
		if got := isFirstContent(c); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

// --- Request attributes ---

func TestExtractRequestAttributes_BasicContents(t *testing.T) {
	span, _ := startTestSpan(t)
	body := `{
		"contents": [{"role": "user", "parts": [{"text": "hello"}]}],
		"generationConfig": {"temperature": 0.7, "maxOutputTokens": 1024, "topP": 0.9}
	}`
	extractRequestAttributes(span, []byte(body))

	if v, ok := getAttr(span, "gen_ai.request.temperature"); !ok || v.AsFloat64() != 0.7 {
		t.Errorf("temperature: got %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.request.max_tokens"); !ok || v.AsInt64() != 1024 {
		t.Errorf("max_tokens: got %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.request.top_p"); !ok || v.AsFloat64() != 0.9 {
		t.Errorf("top_p: got %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.prompt"); !ok {
		t.Error("expected gen_ai.prompt")
	} else if !strings.Contains(v.AsString(), "hello") {
		t.Errorf("prompt should contain 'hello': %s", v.AsString())
	}
}

func TestExtractRequestAttributes_SystemInstruction(t *testing.T) {
	span, _ := startTestSpan(t)
	body := `{
		"systemInstruction": {"parts": [{"text": "You are helpful"}]},
		"contents": [{"role": "user", "parts": [{"text": "hi"}]}]
	}`
	extractRequestAttributes(span, []byte(body))

	v, ok := getAttr(span, "gen_ai.prompt")
	if !ok {
		t.Fatal("expected gen_ai.prompt")
	}
	prompt := v.AsString()
	if !strings.Contains(prompt, "system") {
		t.Errorf("prompt should include system role: %s", prompt)
	}
	if !strings.Contains(prompt, "You are helpful") {
		t.Errorf("prompt should include system content: %s", prompt)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(prompt), &parsed); err != nil {
		t.Fatalf("prompt is not valid JSON: %v", err)
	}
	msgs, _ := parsed["messages"].([]any)
	if len(msgs) != 2 {
		t.Errorf("expected 2 messages (system + user), got %d", len(msgs))
	}
	if first, ok := msgs[0].(map[string]any); ok {
		if first["role"] != "system" {
			t.Errorf("first message role = %v, want 'system'", first["role"])
		}
	}
}

func TestExtractRequestAttributes_MultiTurnConversation(t *testing.T) {
	span, _ := startTestSpan(t)
	body := `{
		"contents": [
			{"role": "user", "parts": [{"text": "What is 2+2?"}]},
			{"role": "model", "parts": [{"text": "4"}]},
			{"role": "user", "parts": [{"text": "And times 3?"}]}
		]
	}`
	extractRequestAttributes(span, []byte(body))

	v, ok := getAttr(span, "gen_ai.prompt")
	if !ok {
		t.Fatal("expected gen_ai.prompt")
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(v.AsString()), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	msgs, _ := parsed["messages"].([]any)
	if len(msgs) != 3 {
		t.Errorf("expected 3 messages, got %d", len(msgs))
	}
	// "model" role should be mapped to "assistant"
	if second, ok := msgs[1].(map[string]any); ok {
		if second["role"] != "assistant" {
			t.Errorf("model role should map to 'assistant', got %v", second["role"])
		}
	}
}

func TestExtractRequestAttributes_FunctionCall(t *testing.T) {
	span, _ := startTestSpan(t)
	body := `{
		"contents": [
			{"role": "user", "parts": [{"text": "What is the weather?"}]},
			{"role": "model", "parts": [{"functionCall": {"name": "get_weather", "args": {"city": "Paris"}}}]},
			{"role": "function", "parts": [{"functionResponse": {"name": "get_weather", "response": {"temp": 20}}}]}
		]
	}`
	extractRequestAttributes(span, []byte(body))

	v, ok := getAttr(span, "gen_ai.prompt")
	if !ok {
		t.Fatal("expected gen_ai.prompt")
	}
	prompt := v.AsString()

	if !strings.Contains(prompt, "get_weather") {
		t.Errorf("prompt should contain function name: %s", prompt)
	}
	// Assistant message should have OpenAI-style tool_calls
	if !strings.Contains(prompt, "tool_calls") {
		t.Errorf("prompt should contain tool_calls: %s", prompt)
	}
	// Tool response should be a separate "tool" role message
	if !strings.Contains(prompt, `"role":"tool"`) {
		t.Errorf("prompt should contain tool role message: %s", prompt)
	}

	var parsed map[string]any
	if err := json.Unmarshal([]byte(prompt), &parsed); err != nil {
		t.Fatalf("invalid JSON: %v", err)
	}
	msgs, _ := parsed["messages"].([]any)
	if len(msgs) != 3 {
		t.Fatalf("expected 3 messages, got %d", len(msgs))
	}
	// First message (user, text-only) should have string content
	if first, ok := msgs[0].(map[string]any); ok {
		if _, isStr := first["content"].(string); !isStr {
			t.Error("text-only user message should have string content")
		}
	}
	// Second message (assistant with tool_calls) should have null content and tool_calls
	if second, ok := msgs[1].(map[string]any); ok {
		if second["role"] != "assistant" {
			t.Errorf("expected assistant role, got %v", second["role"])
		}
		if _, hasTC := second["tool_calls"].([]any); !hasTC {
			t.Error("assistant message should have tool_calls array")
		}
	}
	// Third message (tool response)
	if third, ok := msgs[2].(map[string]any); ok {
		if third["role"] != "tool" {
			t.Errorf("expected tool role, got %v", third["role"])
		}
		if third["name"] != "get_weather" {
			t.Errorf("expected tool name get_weather, got %v", third["name"])
		}
	}
}

func TestExtractRequestAttributes_InvalidJSON(t *testing.T) {
	span, _ := startTestSpan(t)
	extractRequestAttributes(span, []byte("not json"))
	if _, ok := getAttr(span, "gen_ai.prompt"); ok {
		t.Error("should not set prompt for invalid JSON")
	}
}

// --- Non-streaming response attributes ---

func TestExtractResponseAttributes(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	body := `{
		"candidates": [{"content": {"parts": [{"text": "Hello!"}], "role": "model"}, "finishReason": "STOP"}],
		"modelVersion": "gemini-2.0-flash-001",
		"responseId": "resp-123",
		"usageMetadata": {"promptTokenCount": 10, "candidatesTokenCount": 5}
	}`
	extractResponseAttributes(span, []byte(body), parentSpan)

	if v, ok := getAttr(span, "gen_ai.response.model"); !ok || v.AsString() != "gemini-2.0-flash-001" {
		t.Errorf("model: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.response.id"); !ok || v.AsString() != "resp-123" {
		t.Errorf("response id: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.completion"); !ok || !strings.Contains(v.AsString(), "Hello!") {
		t.Errorf("completion should contain 'Hello!': %v", v)
	}
	if v, ok := getAttr(span, "langsmith.metadata.stop_reason"); !ok || v.AsString() != "STOP" {
		t.Errorf("stop_reason: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 10 {
		t.Errorf("input_tokens: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 5 {
		t.Errorf("output_tokens: %v", v)
	}
}

func TestExtractResponseAttributes_FunctionCall(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	body := `{
		"candidates": [{
			"content": {
				"parts": [{"functionCall": {"name": "get_weather", "args": {"city": "SF"}}}],
				"role": "model"
			},
			"finishReason": "STOP"
		}],
		"usageMetadata": {"promptTokenCount": 15, "candidatesTokenCount": 8}
	}`
	extractResponseAttributes(span, []byte(body), parentSpan)

	v, ok := getAttr(span, "gen_ai.completion")
	if !ok {
		t.Fatal("expected gen_ai.completion")
	}
	completion := v.AsString()

	// OpenAI-compatible tool_calls format
	if !strings.Contains(completion, "tool_calls") {
		t.Errorf("completion should contain tool_calls: %s", completion)
	}
	if !strings.Contains(completion, "get_weather") {
		t.Errorf("completion should contain function name: %s", completion)
	}
	if !strings.Contains(completion, `"type":"function"`) {
		t.Errorf("tool_calls should have type 'function': %s", completion)
	}
	if !strings.Contains(completion, `"id":"call_0"`) {
		t.Errorf("tool_calls should have synthetic id: %s", completion)
	}
	// content should be null when only tool calls
	if !strings.Contains(completion, `"content":null`) {
		t.Errorf("content should be null with tool-only response: %s", completion)
	}
	// finish_reason should be in the message
	if !strings.Contains(completion, `"finish_reason":"STOP"`) {
		t.Errorf("completion should include finish_reason: %s", completion)
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

func TestExtractResponseAttributes_MultipleTextParts(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	body := `{
		"candidates": [{
			"content": {"parts": [{"text": "Hello!"}, {"text": " World"}], "role": "model"},
			"finishReason": "STOP"
		}]
	}`
	extractResponseAttributes(span, []byte(body), parentSpan)

	v, ok := getAttr(span, "gen_ai.completion")
	if !ok {
		t.Fatal("expected gen_ai.completion")
	}
	if !strings.Contains(v.AsString(), "Hello! World") {
		t.Errorf("completion should concatenate text parts: %s", v.AsString())
	}
}

// --- Streaming response attributes ---

func TestExtractStreamingResponseAttributes_TextContent(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	sse := "data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\"Hello\"}],\"role\":\"model\"}}],\"modelVersion\":\"gemini-2.0-flash-001\"}\n" +
		"data: {\"candidates\":[{\"content\":{\"parts\":[{\"text\":\" world\"}],\"role\":\"model\"},\"finishReason\":\"STOP\"}],\"usageMetadata\":{\"promptTokenCount\":10,\"candidatesTokenCount\":5},\"modelVersion\":\"gemini-2.0-flash-001\",\"responseId\":\"resp-xyz\"}\n"

	extractStreamingResponseAttributes(span, []byte(sse), parentSpan)

	if v, ok := getAttr(span, "gen_ai.response.model"); !ok || v.AsString() != "gemini-2.0-flash-001" {
		t.Errorf("model: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.response.id"); !ok || v.AsString() != "resp-xyz" {
		t.Errorf("response id: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.completion"); !ok || !strings.Contains(v.AsString(), "Hello world") {
		t.Errorf("completion should contain 'Hello world': %v", v)
	}
	if v, ok := getAttr(span, "langsmith.metadata.stop_reason"); !ok || v.AsString() != "STOP" {
		t.Errorf("stop_reason: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 10 {
		t.Errorf("input_tokens: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 5 {
		t.Errorf("output_tokens: %v", v)
	}
}

func TestExtractStreamingResponseAttributes_FunctionCall(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	sse := "data: {\"candidates\":[{\"content\":{\"parts\":[{\"functionCall\":{\"name\":\"get_weather\",\"args\":{\"city\":\"Paris\"}}}],\"role\":\"model\"},\"finishReason\":\"STOP\"}],\"usageMetadata\":{\"promptTokenCount\":20,\"candidatesTokenCount\":12}}\n"

	extractStreamingResponseAttributes(span, []byte(sse), parentSpan)

	v, ok := getAttr(span, "gen_ai.completion")
	if !ok {
		t.Fatal("expected gen_ai.completion")
	}
	completion := v.AsString()
	if !strings.Contains(completion, "get_weather") {
		t.Errorf("completion should contain function name: %s", completion)
	}
	if !strings.Contains(completion, `"type":"function"`) {
		t.Errorf("streaming tool_calls should have type 'function': %s", completion)
	}
	if !strings.Contains(completion, `"id":"call_0"`) {
		t.Errorf("streaming tool_calls should have synthetic id: %s", completion)
	}
	if !strings.Contains(completion, `"content":null`) {
		t.Errorf("content should be null with tool-only response: %s", completion)
	}
	if !strings.Contains(completion, `"finish_reason":"STOP"`) {
		t.Errorf("streaming completion should include finish_reason: %s", completion)
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

// --- Usage attributes ---

func TestSetUsageAttributes(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)

	usage := map[string]any{
		"promptTokenCount":        float64(100),
		"candidatesTokenCount":    float64(50),
		"cachedContentTokenCount": float64(10),
		"thoughtsTokenCount":      float64(5),
	}
	setUsageAttributes(span, usage, parentSpan)

	if v, ok := getAttr(span, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 100 {
		t.Errorf("input_tokens: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 50 {
		t.Errorf("output_tokens: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.cache_read_input_tokens"); !ok || v.AsInt64() != 10 {
		t.Errorf("cache_read: %v", v)
	}
	if v, ok := getAttr(span, "gen_ai.usage.details.reasoning_tokens"); !ok || v.AsInt64() != 5 {
		t.Errorf("reasoning_tokens: %v", v)
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

	usage := map[string]any{
		"promptTokenCount":     float64(10),
		"candidatesTokenCount": float64(5),
	}
	setUsageAttributes(rwChild, usage, rwParent)

	if v, ok := getAttr(rwParent, "gen_ai.usage.input_tokens"); !ok || v.AsInt64() != 10 {
		t.Errorf("parent input_tokens: %v", v)
	}
	if v, ok := getAttr(rwParent, "gen_ai.usage.output_tokens"); !ok || v.AsInt64() != 5 {
		t.Errorf("parent output_tokens: %v", v)
	}
}

func TestSetUsageAttributes_EmptyUsage(t *testing.T) {
	span, _ := startTestSpan(t)
	parentSpan, _ := startTestSpan(t)
	setUsageAttributes(span, map[string]any{}, parentSpan)

	if _, ok := getAttr(span, "gen_ai.usage.input_tokens"); ok {
		t.Error("should not set input_tokens for empty usage")
	}
}

// --- mapGeminiRole ---

func TestMapGeminiRole(t *testing.T) {
	tests := []struct {
		input any
		want  string
	}{
		{"model", "assistant"},
		{"user", "user"},
		{"function", "tool"},
		{"", "user"},
		{nil, "user"},
		{"custom", "custom"},
	}
	for _, tt := range tests {
		if got := mapGeminiRole(tt.input); got != tt.want {
			t.Errorf("mapGeminiRole(%v) = %q, want %q", tt.input, got, tt.want)
		}
	}
}
