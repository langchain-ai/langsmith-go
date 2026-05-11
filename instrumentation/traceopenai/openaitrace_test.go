package traceopenai

import (
	"bytes"
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"

	"go.opentelemetry.io/otel/codes"
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

// doStreamingChat sends a chat-completions stream request through client and
// drains the response body so the SSE chunks flow through the wrapper.
func doStreamingChat(t *testing.T, client *http.Client) {
	t.Helper()
	body := `{"model":"gpt-4","messages":[{"role":"user","content":"hello"}],"stream":true}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/chat/completions", strings.NewReader(body))
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
	// The role chunk is preamble — it must NOT trip new_token. Only the
	// content chunk that follows should.
	sse := "data: {\"choices\":[{\"delta\":{\"role\":\"assistant\"}}]}\n\n" +
		"data: {\"choices\":[{\"delta\":{\"content\":\"hi\"}}]}\n\n" +
		"data: [DONE]\n\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreamingChat(t, client)

	if !hasEvent(exporter.GetSpans(), "new_token") {
		t.Fatalf("expected new_token event, got %v", exporter.GetSpans())
	}
}

func TestRoundTrip_StreamingEmitsNewTokenOnToolCallDelta(t *testing.T) {
	// A stream that only ever emits tool_calls (no text content) must still
	// register first-token time.
	sse := "data: {\"choices\":[{\"delta\":{\"role\":\"assistant\"}}]}\n\n" +
		"data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"c\",\"function\":{\"name\":\"x\"}}]}}]}\n\n" +
		"data: [DONE]\n\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreamingChat(t, client)

	if !hasEvent(exporter.GetSpans(), "new_token") {
		t.Fatalf("expected new_token event for tool_calls delta, got %v", exporter.GetSpans())
	}
}

func TestRoundTrip_StreamingPreambleOnlyDoesNotEmitNewToken(t *testing.T) {
	// Stream that contains only the role preamble and then ends — there's no
	// content, so new_token must NOT be emitted.
	sse := "data: {\"choices\":[{\"delta\":{\"role\":\"assistant\"}}]}\n\n" +
		"data: [DONE]\n\n"
	client, exporter := newTracedClient(t, []byte(sse))
	doStreamingChat(t, client)

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("preamble-only stream should not emit new_token, got %v", exporter.GetSpans())
	}
}

func TestRoundTrip_NonStreamingDoesNotEmitNewTokenEvent(t *testing.T) {
	resp := `{"choices":[{"message":{"role":"assistant","content":"hi"}}],"usage":{"prompt_tokens":1,"completion_tokens":1}}`
	client, exporter := newTracedClient(t, []byte(resp))

	body := `{"model":"gpt-4","messages":[{"role":"user","content":"hello"}]}`
	req, err := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/chat/completions", strings.NewReader(body))
	if err != nil {
		t.Fatal(err)
	}
	httpResp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(httpResp.Body); err != nil {
		t.Fatal(err)
	}
	httpResp.Body.Close()

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("non-streaming span should not record a new_token event")
	}
}

// errStatusTransport returns a non-streaming JSON error body with a 429.
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
		Transport: &errStatusTransport{body: []byte(`{"error":{"message":"rate limited"}}`)},
	}, WithTracerProvider(tp))

	body := `{"model":"gpt-4","messages":[{"role":"user","content":"hello"}],"stream":true}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/chat/completions", strings.NewReader(body))
	resp, err := client.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	if _, err := io.ReadAll(resp.Body); err != nil {
		t.Fatal(err)
	}
	resp.Body.Close()

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("streamed HTTP error should not emit new_token, got %v", exporter.GetSpans())
	}
}

func TestRoundTrip_ResponsesAPIStreamingFirstTextDelta(t *testing.T) {
	// Responses API: lifecycle envelopes (response.created, output_item.added)
	// arrive before the first text delta. Only the delta should register.
	sse := "data: {\"type\":\"response.created\",\"response\":{}}\n\n" +
		"data: {\"type\":\"response.in_progress\",\"response\":{}}\n\n" +
		"data: {\"type\":\"response.output_item.added\",\"item\":{}}\n\n" +
		"data: {\"type\":\"response.output_text.delta\",\"delta\":\"Hi\"}\n\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"Hi\"}]}]}}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))

	body := `{"model":"gpt-4o","input":"hi","stream":true}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/responses", strings.NewReader(body))
	resp, _ := client.Do(req)
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if !hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("expected new_token on response.output_text.delta, got %v", exporter.GetSpans())
	}
}

func TestRoundTrip_ResponsesAPIStreamingLifecycleOnlyDoesNotEmit(t *testing.T) {
	// Responses API stream that never produces a delta event must not emit
	// new_token even though many bytes flow through.
	sse := "data: {\"type\":\"response.created\",\"response\":{}}\n\n" +
		"data: {\"type\":\"response.in_progress\",\"response\":{}}\n\n" +
		"data: {\"type\":\"response.output_item.added\",\"item\":{}}\n\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"output\":[]}}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))

	body := `{"model":"gpt-4o","input":"hi","stream":true}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/responses", strings.NewReader(body))
	resp, _ := client.Do(req)
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if hasEvent(exporter.GetSpans(), "new_token") {
		t.Errorf("lifecycle-only Responses stream should not emit new_token")
	}
}

func TestRoundTrip_ResponsesAPIStreamingCompleteIsNotError(t *testing.T) {
	// A Responses API stream that ends with response.completed must NOT be
	// marked as an incomplete/cancelled stream — it completed normally.
	sse := "data: {\"type\":\"response.created\",\"response\":{}}\n\n" +
		"data: {\"type\":\"response.output_text.delta\",\"delta\":\"Hi\"}\n\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"Hi\"}]}],\"usage\":{\"input_tokens\":2,\"output_tokens\":1}}}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))

	body := `{"model":"gpt-4o","input":"hi","stream":true}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/responses", strings.NewReader(body))
	resp, _ := client.Do(req)
	io.ReadAll(resp.Body)
	resp.Body.Close()

	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	for _, s := range spans {
		if s.Status.Code == codes.Error {
			t.Errorf("span %q should not be errored, got status %v", s.Name, s.Status)
		}
	}
}

func TestRoundTrip_ResponsesAPIStreamingIncompleteIsError(t *testing.T) {
	// A Responses API stream that ends WITHOUT response.completed should be
	// flagged as incomplete (e.g. cancelled or network failure).
	sse := "data: {\"type\":\"response.created\",\"response\":{}}\n\n" +
		"data: {\"type\":\"response.output_text.delta\",\"delta\":\"Hi\"}\n\n"
	client, exporter := newTracedClient(t, []byte(sse))

	body := `{"model":"gpt-4o","input":"hi","stream":true}`
	req, _ := http.NewRequestWithContext(context.Background(), http.MethodPost,
		"https://api.openai.com/v1/responses", strings.NewReader(body))
	resp, _ := client.Do(req)
	io.ReadAll(resp.Body)
	resp.Body.Close()

	spans := exporter.GetSpans()
	if len(spans) == 0 {
		t.Fatal("expected at least one span")
	}
	found := false
	for _, s := range spans {
		if s.Status.Code == codes.Error {
			found = true
		}
	}
	if !found {
		t.Error("incomplete Responses stream should have error status")
	}
}

func TestIsFirstContentChat(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{"role preamble", `{"choices":[{"delta":{"role":"assistant"}}]}`, false},
		{"empty content", `{"choices":[{"delta":{"content":""}}]}`, false},
		{"content delta", `{"choices":[{"delta":{"content":"hi"}}]}`, true},
		{"tool_calls delta", `{"choices":[{"delta":{"tool_calls":[{"id":"x"}]}}]}`, true},
		{"legacy text", `{"choices":[{"text":"hi"}]}`, true},
		{"no choices", `{"id":"x"}`, false},
		{"empty choices", `{"choices":[]}`, false},
		{"n>1 content on second choice", `{"choices":[{"index":0,"delta":{"role":"assistant"}},{"index":1,"delta":{"content":"hi"}}]}`, true},
		{"n>1 all preamble", `{"choices":[{"index":0,"delta":{"role":"assistant"}},{"index":1,"delta":{"role":"assistant"}}]}`, false},
	}
	for _, tt := range tests {
		var c map[string]any
		_ = json.Unmarshal([]byte(tt.body), &c)
		if got := isFirstContentChat(c); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsFirstContentResponses(t *testing.T) {
	tests := []struct {
		name string
		body string
		want bool
	}{
		{"created lifecycle", `{"type":"response.created"}`, false},
		{"in_progress", `{"type":"response.in_progress"}`, false},
		{"output_item.added", `{"type":"response.output_item.added"}`, false},
		{"output_text.delta", `{"type":"response.output_text.delta","delta":"hi"}`, true},
		{"function_call_arguments.delta", `{"type":"response.function_call_arguments.delta"}`, true},
		{"reasoning_summary_text.delta", `{"type":"response.reasoning_summary_text.delta"}`, true},
		{"completed lifecycle", `{"type":"response.completed"}`, false},
	}
	for _, tt := range tests {
		var c map[string]any
		_ = json.Unmarshal([]byte(tt.body), &c)
		if got := isFirstContentResponses(c); got != tt.want {
			t.Errorf("%s: got %v, want %v", tt.name, got, tt.want)
		}
	}
}

func TestIsOpenAIEndpoint(t *testing.T) {
	tests := []struct {
		path string
		want bool
	}{
		{"/v1/chat/completions", true},
		{"/v1/completions", true},
		{"/v1/embeddings", true},
		{"/v1/responses", true},
		{"/openai/deployments/gpt-4/chat/completions", true},
		{"/v1/models", false},
		{"/v1/files", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := isOpenAIEndpoint(tt.path); got != tt.want {
			t.Errorf("isOpenAIEndpoint(%q) = %v, want %v", tt.path, got, tt.want)
		}
	}
}

func TestGetSpanName(t *testing.T) {
	tests := []struct {
		path string
		want string
	}{
		{"/v1/chat/completions", "openai.chat.completion"},
		{"/v1/completions", "openai.completion"},
		{"/v1/embeddings", "openai.embedding"},
		{"/v1/responses", "openai.responses"},
		{"/v1/models", "openai.request"},
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
		{"/v1/chat/completions", "chat"},
		{"/v1/completions", "completion"},
		{"/v1/embeddings", "embedding"},
		{"/v1/responses", "responses"},
		{"/v1/models", "request"},
	}
	for _, tt := range tests {
		if got := getOperationName(tt.path); got != tt.want {
			t.Errorf("getOperationName(%q) = %q, want %q", tt.path, got, tt.want)
		}
	}
}

func TestParseRequestBody_ChatCompletion(t *testing.T) {
	body := `{"model":"gpt-4","messages":[{"role":"user","content":"hello"}],"stream":true}`
	fields := parseRequestBody([]byte(body))

	if fields.model != "gpt-4" {
		t.Errorf("model = %q, want gpt-4", fields.model)
	}
	if !fields.streaming {
		t.Error("expected streaming=true")
	}
	if fields.inputMessages == "" {
		t.Fatal("expected non-empty inputMessages")
	}
	if !strings.Contains(fields.inputMessages, "hello") {
		t.Errorf("inputMessages should contain 'hello': %s", fields.inputMessages)
	}
}

func TestParseRequestBody_ResponsesAPI(t *testing.T) {
	body := `{"model":"gpt-4","input":"What is Go?"}`
	fields := parseRequestBody([]byte(body))

	if fields.model != "gpt-4" {
		t.Errorf("model = %q, want gpt-4", fields.model)
	}
	if fields.inputMessages == "" {
		t.Fatal("expected non-empty inputMessages")
	}
	if !strings.Contains(fields.inputMessages, "What is Go?") {
		t.Errorf("inputMessages should contain input text: %s", fields.inputMessages)
	}
}

func TestParseRequestBody_LegacyCompletion(t *testing.T) {
	body := `{"model":"text-davinci-003","prompt":"Once upon a time"}`
	fields := parseRequestBody([]byte(body))

	if fields.model != "text-davinci-003" {
		t.Errorf("model = %q, want text-davinci-003", fields.model)
	}
	if !strings.Contains(fields.inputMessages, "Once upon a time") {
		t.Errorf("inputMessages should contain prompt: %s", fields.inputMessages)
	}
}

func TestParseRequestBody_InvalidJSON(t *testing.T) {
	fields := parseRequestBody([]byte("not json"))
	if fields.model != "" || fields.streaming || fields.inputMessages != "" {
		t.Errorf("expected empty fields for invalid JSON, got %+v", fields)
	}
}

func TestMarshalMessages(t *testing.T) {
	msgs := []any{map[string]any{"role": "user", "content": "hi"}}
	result := marshalMessages(msgs)
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	var parsed map[string]any
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if _, ok := parsed["messages"]; !ok {
		t.Error("expected 'messages' key in output")
	}
}

func TestExtractCompletionFromResponse_ChatCompletion(t *testing.T) {
	body := `{
		"choices": [{"message": {"role": "assistant", "content": "Hello!"}}],
		"usage": {"prompt_tokens": 10, "completion_tokens": 5}
	}`
	completion, usage := extractCompletionFromResponse([]byte(body))

	if !strings.Contains(completion, "Hello!") {
		t.Errorf("completion should contain 'Hello!': %s", completion)
	}
	if usage.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want 10", usage.InputTokens)
	}
	if usage.OutputTokens != 5 {
		t.Errorf("OutputTokens = %d, want 5", usage.OutputTokens)
	}
}

func TestExtractCompletionFromResponse_WithToolCalls(t *testing.T) {
	body := `{
		"choices": [{"message": {"role": "assistant", "content": null, "tool_calls": [{"id": "call_1", "type": "function", "function": {"name": "get_weather", "arguments": "{\"city\":\"SF\"}"}}]}}],
		"usage": {"prompt_tokens": 15, "completion_tokens": 8}
	}`
	completion, usage := extractCompletionFromResponse([]byte(body))

	if !strings.Contains(completion, "get_weather") {
		t.Errorf("completion should contain tool call: %s", completion)
	}
	if usage.InputTokens != 15 {
		t.Errorf("InputTokens = %d, want 15", usage.InputTokens)
	}
}

func TestExtractCompletionFromResponse_LegacyCompletion(t *testing.T) {
	body := `{
		"choices": [{"text": "Once upon a time"}],
		"usage": {"prompt_tokens": 5, "completion_tokens": 4}
	}`
	completion, _ := extractCompletionFromResponse([]byte(body))

	if !strings.Contains(completion, "Once upon a time") {
		t.Errorf("completion should contain text: %s", completion)
	}
}

func TestExtractCompletionFromResponse_InvalidJSON(t *testing.T) {
	completion, usage := extractCompletionFromResponse([]byte("bad"))
	if completion != "" {
		t.Errorf("expected empty completion, got %q", completion)
	}
	if usage.InputTokens != 0 || usage.OutputTokens != 0 {
		t.Errorf("expected zero usage, got %+v", usage)
	}
}

func TestExtractStreamingCompletion_TextOnly(t *testing.T) {
	sse := "data: {\"choices\":[{\"delta\":{\"content\":\"Hello\"}}]}\n" +
		"data: {\"choices\":[{\"delta\":{\"content\":\" world\"}}]}\n" +
		"data: {\"usage\":{\"prompt_tokens\":5,\"completion_tokens\":2}}\n" +
		"data: [DONE]\n"
	completion, usage := extractStreamingCompletion([]byte(sse))

	if !strings.Contains(completion, "Hello world") {
		t.Errorf("completion should contain 'Hello world': %s", completion)
	}
	if usage.InputTokens != 5 {
		t.Errorf("InputTokens = %d, want 5", usage.InputTokens)
	}
	if usage.OutputTokens != 2 {
		t.Errorf("OutputTokens = %d, want 2", usage.OutputTokens)
	}
}

func TestExtractStreamingCompletion_WithToolCalls(t *testing.T) {
	sse := "data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"id\":\"call_1\",\"type\":\"function\",\"function\":{\"name\":\"get_weather\",\"arguments\":\"\"}}]}}]}\n" +
		"data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\"{\\\"city\\\"\"}}]}}]}\n" +
		"data: {\"choices\":[{\"delta\":{\"tool_calls\":[{\"index\":0,\"function\":{\"arguments\":\": \\\"SF\\\"}\"}}]}}]}\n" +
		"data: [DONE]\n"
	completion, _ := extractStreamingCompletion([]byte(sse))

	if !strings.Contains(completion, "get_weather") {
		t.Errorf("completion should contain function name: %s", completion)
	}
	if !strings.Contains(completion, "call_1") {
		t.Errorf("completion should contain call id: %s", completion)
	}
}

func TestExtractStreamingCompletion_Empty(t *testing.T) {
	completion, usage := extractStreamingCompletion([]byte(""))
	if completion != "" {
		t.Errorf("expected empty completion, got %q", completion)
	}
	if usage.InputTokens != 0 || usage.OutputTokens != 0 {
		t.Errorf("expected zero usage, got %+v", usage)
	}
}

func TestExtractResponsesCompletion(t *testing.T) {
	body := `{
		"output": [{"type": "message", "content": [{"type": "output_text", "text": "Hi there"}]}],
		"usage": {"input_tokens": 3, "output_tokens": 2}
	}`
	completion, usage := extractResponsesCompletion([]byte(body))

	if !strings.Contains(completion, "Hi there") {
		t.Errorf("completion should contain 'Hi there': %s", completion)
	}
	if usage.InputTokens != 3 {
		t.Errorf("InputTokens = %d, want 3", usage.InputTokens)
	}
	if usage.OutputTokens != 2 {
		t.Errorf("OutputTokens = %d, want 2", usage.OutputTokens)
	}
}

func TestExtractResponsesCompletion_FunctionCall(t *testing.T) {
	body := `{
		"output": [{"type": "function_call", "name": "search", "arguments": "{\"q\":\"Go\"}", "call_id": "fc_1"}],
		"usage": {"input_tokens": 1, "output_tokens": 1}
	}`
	completion, _ := extractResponsesCompletion([]byte(body))

	if !strings.Contains(completion, "search") {
		t.Errorf("completion should contain function name: %s", completion)
	}
}

func TestExtractStreamingResponsesCompletion(t *testing.T) {
	// Responses API streams do not send [DONE]; the stream ends after the
	// response.completed event and the connection closes.
	sse := "data: {\"type\":\"response.created\"}\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"done\"}]}],\"usage\":{\"input_tokens\":7,\"output_tokens\":3}}}\n"
	completion, usage := extractStreamingResponsesCompletion([]byte(sse))

	if !strings.Contains(completion, "done") {
		t.Errorf("completion should contain 'done': %s", completion)
	}
	if usage.InputTokens != 7 {
		t.Errorf("InputTokens = %d, want 7", usage.InputTokens)
	}
}

func TestExtractStreamingResponsesCompletion_NoCompletedEvent(t *testing.T) {
	sse := "data: {\"type\":\"response.created\"}\n"
	completion, usage := extractStreamingResponsesCompletion([]byte(sse))
	if completion != "" {
		t.Errorf("expected empty completion, got %q", completion)
	}
	if usage.InputTokens != 0 || usage.OutputTokens != 0 {
		t.Errorf("expected zero usage, got %+v", usage)
	}
}

func TestExtractResponsesUsage(t *testing.T) {
	resp := map[string]any{
		"usage": map[string]any{
			"input_tokens":  float64(10),
			"output_tokens": float64(20),
		},
	}
	usage := extractResponsesUsage(resp)
	if usage.InputTokens != 10 {
		t.Errorf("InputTokens = %d, want 10", usage.InputTokens)
	}
	if usage.OutputTokens != 20 {
		t.Errorf("OutputTokens = %d, want 20", usage.OutputTokens)
	}
}

func TestExtractResponsesUsage_NoUsage(t *testing.T) {
	usage := extractResponsesUsage(map[string]any{})
	if usage.InputTokens != 0 || usage.OutputTokens != 0 {
		t.Errorf("expected zero usage, got %+v", usage)
	}
}

func TestExtractResponsesOutput_TextAndFunctionCalls(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type": "message",
				"content": []any{
					map[string]any{"type": "output_text", "text": "Here's the result"},
				},
			},
			map[string]any{
				"type":      "function_call",
				"name":      "search",
				"arguments": `{"q":"test"}`,
				"call_id":   "fc_2",
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "Here's the result") {
		t.Errorf("output should contain text: %s", output)
	}
	if !strings.Contains(output, "search") {
		t.Errorf("output should contain function call: %s", output)
	}
}

func TestExtractResponsesOutput_Empty(t *testing.T) {
	output := extractResponsesOutput(map[string]any{})
	if output != "" {
		t.Errorf("expected empty output, got %q", output)
	}
}

func TestExtractResponsesOutput_Refusal(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type": "message",
				"content": []any{
					map[string]any{"type": "refusal", "refusal": "I cannot help with that"},
				},
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "I cannot help with that") {
		t.Errorf("output should contain refusal text: %s", output)
	}
}

func TestExtractResponsesOutput_Reasoning(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type": "reasoning",
				"summary": []any{
					map[string]any{"type": "summary_text", "text": "Let me think about this"},
				},
			},
			map[string]any{
				"type": "message",
				"content": []any{
					map[string]any{"type": "output_text", "text": "The answer is 42"},
				},
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "Let me think about this") {
		t.Errorf("output should contain reasoning summary: %s", output)
	}
	if !strings.Contains(output, "The answer is 42") {
		t.Errorf("output should contain message text: %s", output)
	}
}

func TestExtractResponsesOutput_WebSearchCall(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":   "web_search_call",
				"id":     "ws_1",
				"status": "completed",
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "web_search") {
		t.Errorf("output should contain web_search tool call: %s", output)
	}
	if !strings.Contains(output, `"type":"function"`) {
		t.Errorf("output should use chat-completions tool_calls format: %s", output)
	}
}

func TestExtractResponsesOutput_FileSearchCall(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":    "file_search_call",
				"id":      "fs_1",
				"queries": []any{"search query"},
				"results": []any{map[string]any{"file_id": "f1", "text": "result"}},
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "file_search") {
		t.Errorf("output should contain file_search tool call: %s", output)
	}
	if !strings.Contains(output, "search query") {
		t.Errorf("output should contain query: %s", output)
	}
}

func TestExtractResponsesOutput_CodeInterpreterCall(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":    "code_interpreter_call",
				"id":      "ci_1",
				"code":    "print('hello')",
				"results": []any{map[string]any{"type": "logs", "logs": "hello"}},
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "code_interpreter") {
		t.Errorf("output should contain code_interpreter tool call: %s", output)
	}
	if !strings.Contains(output, "print") {
		t.Errorf("output should contain code: %s", output)
	}
}

func TestExtractResponsesOutput_ComputerCall(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":    "computer_call",
				"call_id": "comp_1",
				"action":  map[string]any{"type": "click", "x": float64(100), "y": float64(200)},
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "computer") {
		t.Errorf("output should contain computer tool call: %s", output)
	}
	if !strings.Contains(output, "click") {
		t.Errorf("output should contain action: %s", output)
	}
}

func TestExtractResponsesOutput_McpCall(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":         "mcp_call",
				"id":           "mcp_1",
				"server_label": "my_server",
				"name":         "my_tool",
				"arguments":    `{"key":"value"}`,
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "my_server:my_tool") {
		t.Errorf("output should contain server_label:name: %s", output)
	}
}

func TestExtractResponsesOutput_McpListTools(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":         "mcp_list_tools",
				"id":           "mlt_1",
				"server_label": "my_server",
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "mcp_list_tools") {
		t.Errorf("output should contain mcp_list_tools: %s", output)
	}
}

func TestExtractResponsesOutput_ImageGenerationCall(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":   "image_generation_call",
				"id":     "ig_1",
				"status": "completed",
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, "image_generation") {
		t.Errorf("output should contain image_generation tool call: %s", output)
	}
}

func TestExtractResponsesOutput_AgentToolCalls(t *testing.T) {
	tests := []struct {
		name     string
		item     map[string]any
		contains []string
	}{
		{
			name: "local_shell_call",
			item: map[string]any{
				"type": "local_shell_call", "id": "ls_1", "call_id": "lsc_1",
				"action": map[string]any{"type": "exec", "command": []any{"ls", "-la"}},
			},
			contains: []string{"local_shell", "lsc_1", `"type":"function"`},
		},
		{
			name: "shell_call",
			item: map[string]any{
				"type": "shell_call", "call_id": "sc_1",
				"action": map[string]any{"commands": []any{"echo hello"}},
			},
			contains: []string{`"shell"`, "sc_1", `"type":"function"`},
		},
		{
			name: "apply_patch_call",
			item: map[string]any{
				"type": "apply_patch_call", "call_id": "ap_1",
				"operation": map[string]any{"type": "create", "path": "foo.txt", "diff": "+hello"},
			},
			contains: []string{"apply_patch", `"type":"function"`},
		},
		{
			name: "tool_search_call",
			item: map[string]any{
				"type": "tool_search_call", "call_id": "ts_1",
				"arguments": map[string]any{"query": "find files"},
			},
			contains: []string{"tool_search", `"type":"function"`},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resp := map[string]any{"output": []any{tt.item}}
			result := extractResponsesOutput(resp)
			for _, s := range tt.contains {
				if !strings.Contains(result, s) {
					t.Errorf("expected %q in output: %s", s, result)
				}
			}
		})
	}
}

func TestExtractResponsesOutput_CodeInterpreterUsesOutputs(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type": "code_interpreter_call", "id": "ci_1",
				"code": "print('hi')", "container_id": "ctr_1",
				"outputs": []any{map[string]any{"type": "logs", "logs": "hi"}},
				"status":  "completed",
			},
		},
	}
	result := extractResponsesOutput(resp)
	if !strings.Contains(result, "outputs") {
		t.Errorf("expected 'outputs' key (not 'results'), got %s", result)
	}
}

func TestExtractResponsesOutput_FunctionCallFormat(t *testing.T) {
	resp := map[string]any{
		"output": []any{
			map[string]any{
				"type":      "function_call",
				"name":      "get_weather",
				"arguments": `{"city":"SF"}`,
				"call_id":   "fc_1",
			},
		},
	}
	output := extractResponsesOutput(resp)
	if !strings.Contains(output, `"type":"function"`) {
		t.Errorf("function_call should use chat-completions format: %s", output)
	}
	if !strings.Contains(output, `"name":"get_weather"`) {
		t.Errorf("should contain function name: %s", output)
	}
	if !strings.Contains(output, `"id":"fc_1"`) {
		t.Errorf("should contain call id: %s", output)
	}
}

func TestParseRequestBody_ResponsesAPI_ArrayInput(t *testing.T) {
	body := []byte(`{
		"model": "gpt-4o",
		"instructions": "You are a helpful assistant.",
		"input": [
			{"role": "user", "content": [{"type": "input_text", "text": "Hello"}]},
			{"role": "assistant", "content": [{"type": "output_text", "text": "Hi there"}]},
			{"role": "user", "content": "What is 2+2?"}
		]
	}`)
	fields := parseRequestBody(body)
	if fields.model != "gpt-4o" {
		t.Errorf("model = %q, want gpt-4o", fields.model)
	}
	if !strings.Contains(fields.inputMessages, "You are a helpful assistant.") {
		t.Errorf("should contain instructions as system message: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, "Hello") {
		t.Errorf("should contain flattened input_text: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, "Hi there") {
		t.Errorf("should contain flattened output_text: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, "What is 2+2?") {
		t.Errorf("should contain plain string content: %s", fields.inputMessages)
	}
}

func TestParseRequestBody_ResponsesAPI_MultiTurnWithToolCalls(t *testing.T) {
	body := []byte(`{
		"model": "gpt-4o",
		"input": [
			{"role": "user", "content": "What is the weather in SF?"},
			{"type": "function_call", "name": "get_weather", "arguments": "{\"city\":\"SF\"}", "call_id": "fc_1"},
			{"type": "function_call_output", "call_id": "fc_1", "output": "72°F and sunny"},
			{"role": "assistant", "content": "The weather in SF is 72°F and sunny."},
			{"role": "user", "content": "Thanks!"}
		]
	}`)
	fields := parseRequestBody(body)

	if !strings.Contains(fields.inputMessages, "What is the weather in SF?") {
		t.Errorf("should contain user message: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, "get_weather") {
		t.Errorf("should contain function_call as tool_calls: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, `"role":"tool"`) {
		t.Errorf("function_call_output should become role:tool message: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, "72°F and sunny") {
		t.Errorf("should contain tool output: %s", fields.inputMessages)
	}
	if !strings.Contains(fields.inputMessages, "Thanks!") {
		t.Errorf("should contain final user message: %s", fields.inputMessages)
	}
}

func TestNormalizeResponsesInput_WebSearchInInput(t *testing.T) {
	items := []any{
		map[string]any{
			"type":   "web_search_call",
			"id":     "ws_1",
			"status": "completed",
		},
	}
	result := normalizeResponsesInput(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	msg, ok := result[0].(map[string]any)
	if !ok {
		t.Fatal("expected map")
	}
	if msg["role"] != "assistant" {
		t.Errorf("role = %v, want assistant", msg["role"])
	}
	tcs, ok := msg["tool_calls"].([]any)
	if !ok || len(tcs) == 0 {
		t.Errorf("expected tool_calls: %+v", msg)
	}
}

func TestNormalizeResponsesInput_ItemReferenceSkipped(t *testing.T) {
	items := []any{
		map[string]any{
			"type": "item_reference",
			"id":   "ref_123",
		},
		map[string]any{
			"role":    "user",
			"content": "Hello",
		},
	}
	result := normalizeResponsesInput(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 item (item_reference skipped), got %d", len(result))
	}
	msg := result[0].(map[string]any)
	if msg["content"] != "Hello" {
		t.Errorf("content = %v, want Hello", msg["content"])
	}
}

func TestNormalizeResponsesInput_ComputerCallOutput(t *testing.T) {
	items := []any{
		map[string]any{
			"type":    "computer_call",
			"call_id": "comp_1",
			"action":  map[string]any{"type": "click", "x": float64(50), "y": float64(100)},
		},
		map[string]any{
			"type":    "computer_call_output",
			"call_id": "comp_1",
			"output":  "screenshot captured",
		},
	}
	result := normalizeResponsesInput(items)
	if len(result) != 2 {
		t.Fatalf("expected 2 items, got %d", len(result))
	}
	call := result[0].(map[string]any)
	if call["role"] != "assistant" {
		t.Errorf("computer_call role = %v, want assistant", call["role"])
	}
	output := result[1].(map[string]any)
	if output["role"] != "tool" {
		t.Errorf("computer_call_output role = %v, want tool", output["role"])
	}
	if output["content"] != "screenshot captured" {
		t.Errorf("content = %v, want 'screenshot captured'", output["content"])
	}
}

func TestNormalizeResponsesInput_ReasoningItem(t *testing.T) {
	items := []any{
		map[string]any{
			"type": "reasoning",
			"id":   "rs_1",
			"summary": []any{
				map[string]any{"type": "summary_text", "text": "Let me think about this"},
			},
		},
	}
	result := normalizeResponsesInput(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	msg := result[0].(map[string]any)
	if msg["role"] != "assistant" {
		t.Errorf("role = %v, want assistant", msg["role"])
	}
	content, _ := msg["content"].(string)
	if !strings.Contains(content, "Let me think about this") {
		t.Errorf("content should contain reasoning summary: %s", content)
	}
}

func TestNormalizeResponsesInput_ReasoningEmptySummarySkipped(t *testing.T) {
	items := []any{
		map[string]any{
			"type":    "reasoning",
			"id":      "rs_1",
			"summary": []any{},
		},
		map[string]any{
			"role":    "user",
			"content": "Hello",
		},
	}
	result := normalizeResponsesInput(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 item (empty reasoning skipped), got %d", len(result))
	}
}

func TestNormalizeResponsesInput_UnknownTypeGetsRole(t *testing.T) {
	items := []any{
		map[string]any{
			"type": "some_future_type",
			"id":   "x_1",
		},
	}
	result := normalizeResponsesInput(items)
	if len(result) != 1 {
		t.Fatalf("expected 1 item, got %d", len(result))
	}
	msg := result[0].(map[string]any)
	if msg["role"] != "assistant" {
		t.Errorf("unknown type should get role=assistant, got %v", msg["role"])
	}
	content, _ := msg["content"].(string)
	if !strings.Contains(content, "some_future_type") {
		t.Errorf("content should contain the type name: %s", content)
	}
}

func TestNormalizeResponsesInput_AgentToolCallsBecomeAssistant(t *testing.T) {
	toolCallItems := []map[string]any{
		{"type": "local_shell_call", "call_id": "ls_1", "action": map[string]any{"type": "exec"}},
		{"type": "shell_call", "call_id": "sc_1", "action": map[string]any{"commands": []any{"ls"}}},
		{"type": "apply_patch_call", "call_id": "ap_1", "operation": map[string]any{"type": "create"}},
		{"type": "tool_search_call", "call_id": "ts_1"},
	}
	for _, item := range toolCallItems {
		t.Run(item["type"].(string), func(t *testing.T) {
			result := normalizeResponsesInput([]any{item})
			if len(result) != 1 {
				t.Fatalf("expected 1 item, got %d", len(result))
			}
			msg := result[0].(map[string]any)
			if msg["role"] != "assistant" {
				t.Errorf("role = %v, want assistant", msg["role"])
			}
			if _, ok := msg["tool_calls"].([]any); !ok {
				t.Errorf("expected tool_calls: %+v", msg)
			}
		})
	}
}

func TestNormalizeResponsesInput_AgentToolOutputs(t *testing.T) {
	tests := []struct {
		name       string
		item       map[string]any
		wantCallID string
		wantIn     string
	}{
		{
			name:       "local_shell_call_output",
			item:       map[string]any{"type": "local_shell_call_output", "id": "lso_1", "output": "command output here"},
			wantCallID: "lso_1",
			wantIn:     "command output here",
		},
		{
			name:       "apply_patch_call_output",
			item:       map[string]any{"type": "apply_patch_call_output", "call_id": "ap_1", "output": "patch applied"},
			wantCallID: "ap_1",
			wantIn:     "patch applied",
		},
		{
			name:       "custom_tool_call_output",
			item:       map[string]any{"type": "custom_tool_call_output", "call_id": "ct_1", "output": "custom result"},
			wantCallID: "ct_1",
			wantIn:     "custom result",
		},
		{
			name:       "mcp_approval_response",
			item:       map[string]any{"type": "mcp_approval_response", "approval_request_id": "ar_1", "approve": true},
			wantCallID: "ar_1",
			wantIn:     "approved: true",
		},
		{
			name: "tool_search_output",
			item: map[string]any{
				"type":  "tool_search_output",
				"tools": []any{map[string]any{"name": "grep", "description": "search files"}},
			},
			wantIn: "grep",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := normalizeResponsesInput([]any{tt.item})
			if len(result) != 1 {
				t.Fatalf("expected 1 item, got %d", len(result))
			}
			msg := result[0].(map[string]any)
			if msg["role"] != "tool" {
				t.Errorf("role = %v, want tool", msg["role"])
			}
			if tt.wantCallID != "" && msg["tool_call_id"] != tt.wantCallID {
				t.Errorf("tool_call_id = %v, want %s", msg["tool_call_id"], tt.wantCallID)
			}
			content, _ := msg["content"].(string)
			if !strings.Contains(content, tt.wantIn) {
				t.Errorf("content %q should contain %q", content, tt.wantIn)
			}
		})
	}
}

func TestNormalizeResponsesInput_ShellCallOutputExtractsStdout(t *testing.T) {
	items := []any{
		map[string]any{
			"type":    "shell_call_output",
			"call_id": "sc_1",
			"output": []any{
				map[string]any{
					"stdout":  "hello world",
					"stderr":  "",
					"outcome": map[string]any{"type": "exit", "exit_code": float64(0)},
				},
			},
		},
	}
	result := normalizeResponsesInput(items)
	msg := result[0].(map[string]any)
	if msg["role"] != "tool" {
		t.Errorf("role = %v, want tool", msg["role"])
	}
	if msg["tool_call_id"] != "sc_1" {
		t.Errorf("tool_call_id = %v, want sc_1", msg["tool_call_id"])
	}
	if !strings.Contains(msg["content"].(string), "hello world") {
		t.Errorf("content should contain stdout: %v", msg["content"])
	}
}

func TestNormalizeResponsesInput_ShellCallOutputConcatsStderr(t *testing.T) {
	items := []any{
		map[string]any{
			"type":    "shell_call_output",
			"call_id": "sc_2",
			"output": []any{
				map[string]any{
					"stdout":  "some output",
					"stderr":  "warning: something",
					"outcome": map[string]any{"type": "exit", "exit_code": float64(1)},
				},
			},
		},
	}
	result := normalizeResponsesInput(items)
	msg := result[0].(map[string]any)
	content := msg["content"].(string)
	if !strings.Contains(content, "some output") {
		t.Errorf("content should contain stdout: %s", content)
	}
	if !strings.Contains(content, "warning: something") {
		t.Errorf("content should contain stderr: %s", content)
	}
}

// TestResponsesAPI_RoundTrip exercises parseRequestBody and extractResponsesOutput
// with a multi-turn conversation modeled after a Codex coding agent session.
// Field shapes match captured OpenAI wire format (content arrays with input_text,
// reasoning with encrypted_content, web_search_call with action.queries, etc.).
func TestResponsesAPI_RoundTrip(t *testing.T) {
	reqBody := []byte(`{
		"model": "gpt-5.5",
		"instructions": "You are a coding agent.",
		"input": [
			{"type": "message", "role": "developer", "content": [
				{"type": "input_text", "text": "Sandbox mode is workspace-write."}
			]},
			{"type": "message", "role": "user", "content": [
				{"type": "input_text", "text": "List the files in my home directory"}
			]},

			{"type": "reasoning", "summary": [
				{"type": "summary_text", "text": "I'll run ls to list the directory contents."}
			], "content": null, "encrypted_content": "enc_abc"},
			{"type": "function_call", "name": "exec_command", "call_id": "call_ls1",
			 "arguments": "{\"cmd\":\"ls\",\"workdir\":\"/Users/dev\"}"},
			{"type": "function_call_output", "call_id": "call_ls1",
			 "output": "Desktop\nDocuments\nDownloads\nsrc\n"},
			{"type": "message", "role": "assistant", "content": [
				{"type": "output_text", "text": "Here are the files in your home directory."}
			]},

			{"type": "message", "role": "user", "content": [
				{"type": "input_text", "text": "Now search for recent Go news"}
			]},

			{"type": "reasoning", "summary": [], "content": null, "encrypted_content": "enc_def"},
			{"type": "web_search_call", "status": "completed", "action": {
				"type": "search",
				"query": "Go programming language news 2026",
				"queries": ["Go programming language news 2026", "Go 1.25 release"]
			}},
			{"type": "message", "role": "assistant", "content": [
				{"type": "output_text", "text": "Go 1.25 was released with improved generics."}
			]},

			{"type": "message", "role": "user", "content": [
				{"type": "input_text", "text": "Create a hello.go file and run it"}
			]},

			{"type": "shell_call", "call_id": "sc_1", "action": {
				"type": "exec", "commands": ["echo 'package main' > hello.go"]
			}},
			{"type": "shell_call_output", "call_id": "sc_1", "output": [
				{"stdout": "", "stderr": "", "outcome": {"type": "exit", "exit_code": 0}}
			]},
			{"type": "apply_patch_call", "call_id": "ap_1", "operation": {
				"type": "create", "path": "hello.go",
				"diff": "+package main\n+\n+func main() {\n+\tprintln(\"hello\")\n+}"
			}},
			{"type": "apply_patch_call_output", "call_id": "ap_1", "output": "patch applied"},

			{"type": "message", "role": "user", "content": [
				{"type": "input_text", "text": "Run it"}
			]}
		]
	}`)

	fields := parseRequestBody(reqBody)
	if fields.model != "gpt-5.5" {
		t.Errorf("model = %q, want gpt-5.5", fields.model)
	}

	in := fields.inputMessages
	checks := []struct {
		desc string
		want string
	}{
		{"instructions as system message", "coding agent"},
		{"developer message", "Sandbox mode"},
		{"first user message", "List the files"},
		{"reasoning summary", "run ls"},
		{"function_call tool_calls", "exec_command"},
		{"function_call_output as tool role", `"role":"tool"`},
		{"function output content", "Desktop"},
		{"assistant reply", "files in your home"},
		{"second user turn", "search for recent Go"},
		{"web_search_call", "web_search"},
		{"third user turn", "Create a hello.go"},
		{"shell_call", `"shell"`},
		{"shell_call_output as tool", "sc_1"},
		{"apply_patch_call", "apply_patch"},
		{"apply_patch_call_output as tool", "patch applied"},
		{"final user turn", "Run it"},
	}
	for _, c := range checks {
		if !strings.Contains(in, c.want) {
			t.Errorf("%s: %q not found in input", c.desc, c.want)
		}
	}

	respBody := map[string]any{
		"id": "resp_abc", "object": "response", "model": "gpt-5.5-2026-04-23",
		"output": []any{
			map[string]any{
				"id": "rs_1", "type": "reasoning",
				"encrypted_content": "enc_xyz",
				"summary":           []any{},
			},
			map[string]any{
				"id": "ws_1", "type": "web_search_call",
				"status": "completed",
				"action": map[string]any{
					"type":    "search",
					"query":   "how to run Go programs",
					"queries": []any{"how to run Go programs", "go run command"},
				},
			},
			map[string]any{
				"id": "sc_2", "type": "shell_call", "call_id": "sc_2",
				"action": map[string]any{"commands": []any{"go run hello.go"}},
			},
			map[string]any{
				"id": "fc_1", "type": "function_call", "call_id": "call_run1",
				"name": "exec_command", "arguments": `{"cmd":"go run hello.go"}`,
				"status": "completed",
			},
			map[string]any{
				"id": "msg_1", "type": "message", "role": "assistant",
				"status": "completed",
				"content": []any{
					map[string]any{
						"type":        "output_text",
						"text":        "The program printed hello.",
						"logprobs":    []any{},
						"annotations": []any{},
					},
				},
			},
		},
		"usage": map[string]any{
			"input_tokens":  float64(500),
			"output_tokens": float64(50),
			"total_tokens":  float64(550),
			"output_tokens_details": map[string]any{
				"reasoning_tokens": float64(22),
			},
		},
	}

	out := extractResponsesOutput(respBody)
	if out == "" {
		t.Fatal("extractResponsesOutput returned empty")
	}

	outChecks := []struct {
		desc string
		want string
	}{
		{"web_search as tool call", "web_search"},
		{"web_search uses function format", `"type":"function"`},
		{"shell_call name", `"shell"`},
		{"function_call name", "exec_command"},
		{"function_call args", "go run hello.go"},
		{"message text", "printed hello"},
	}
	for _, c := range outChecks {
		if !strings.Contains(out, c.want) {
			t.Errorf("%s: %q not found in output", c.desc, c.want)
		}
	}
}
