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
