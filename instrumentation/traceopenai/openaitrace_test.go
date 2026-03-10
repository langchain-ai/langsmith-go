package traceopenai

import (
	"encoding/json"
	"strings"
	"testing"
)

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
	sse := "data: {\"type\":\"response.created\"}\n" +
		"data: {\"type\":\"response.completed\",\"response\":{\"output\":[{\"type\":\"message\",\"content\":[{\"type\":\"output_text\",\"text\":\"done\"}]}],\"usage\":{\"input_tokens\":7,\"output_tokens\":3}}}\n" +
		"data: [DONE]\n"
	completion, usage := extractStreamingResponsesCompletion([]byte(sse))

	if !strings.Contains(completion, "done") {
		t.Errorf("completion should contain 'done': %s", completion)
	}
	if usage.InputTokens != 7 {
		t.Errorf("InputTokens = %d, want 7", usage.InputTokens)
	}
}

func TestExtractStreamingResponsesCompletion_NoCompletedEvent(t *testing.T) {
	sse := "data: {\"type\":\"response.created\"}\ndata: [DONE]\n"
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
