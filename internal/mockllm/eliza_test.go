package mockllm

import (
	"encoding/json"
	"net/http"
	"strings"
	"testing"
)

func TestEliza_Greeting(t *testing.T) {
	resp := elizaRespond("hello")
	if resp == "" {
		t.Fatal("expected non-empty response")
	}
	if !strings.Contains(strings.ToLower(resp), "hello") && !strings.Contains(strings.ToLower(resp), "hi") && !strings.Contains(strings.ToLower(resp), "tell me") {
		t.Logf("greeting response: %s", resp) // not a failure, just log
	}
}

func TestEliza_IAmSad(t *testing.T) {
	resp := elizaRespond("I am sad")
	if resp == "" {
		t.Fatal("expected non-empty response")
	}
	// Should reflect back with "you are sad" type phrasing
	lower := strings.ToLower(resp)
	if !strings.Contains(lower, "sad") && !strings.Contains(lower, "you") {
		t.Logf("i am response: %s", resp)
	}
}

func TestEliza_Reflection(t *testing.T) {
	result := reflect("I am feeling my way")
	if !strings.Contains(result, "you") || !strings.Contains(result, "your") {
		t.Errorf("reflect should swap pronouns: got %q", result)
	}
}

func TestEliza_EmptyInput(t *testing.T) {
	resp := elizaRespond("")
	if resp == "" {
		t.Fatal("expected fallback response for empty input")
	}
}

func TestEliza_Fallback(t *testing.T) {
	// Nonsense that shouldn't match any rule
	resp := elizaRespond("xyzzy plugh")
	if resp == "" {
		t.Fatal("expected fallback response")
	}
}

func TestElizaHandler_WithAnthropicServer(t *testing.T) {
	srv := NewAnthropicServer(WithHandler(ElizaHandler()))
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"I am feeling anxious about my work"}],"max_tokens":100}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	content := result["content"].([]any)
	block := content[0].(map[string]any)
	text := block["text"].(string)
	if text == "" {
		t.Fatal("expected non-empty Eliza response")
	}
	t.Logf("Eliza says: %s", text)
}

func TestElizaHandler_WithOpenAIServer(t *testing.T) {
	srv := NewOpenAIServer(WithHandler(ElizaHandler()))
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"I think I need help"}]}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	choices := result["choices"].([]any)
	msg := choices[0].(map[string]any)["message"].(map[string]any)
	text := msg["content"].(string)
	if text == "" {
		t.Fatal("expected non-empty Eliza response")
	}
	t.Logf("Eliza says: %s", text)
}

func TestElizaHandler_MultiTurn(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model: "test",
		Messages: []Message{
			{Role: "user", Content: "hello"},
			{Role: "assistant", Content: "Hello. How are you feeling today?"},
			{Role: "user", Content: "I am feeling worried about my mother"},
		},
	})

	if resp.Content == "" {
		t.Fatal("expected non-empty response")
	}
	// Should pick up on "mother" or "I am" pattern
	t.Logf("Eliza says: %s", resp.Content)
}

func TestEliza_WhatCanYouDo(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "What can you do?"}},
	})
	if !strings.Contains(resp.Content, "special commands") {
		t.Errorf("expected help text, got: %s", resp.Content)
	}
}

func TestEliza_FailWithStatus(t *testing.T) {
	h := ElizaHandler()

	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "Please fail with 429"}},
	})
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Status != 429 {
		t.Errorf("status = %d, want 429", resp.Error.Status)
	}
}

func TestEliza_FailWithStatusAndMessage(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "Please fail with 500 internal server error"}},
	})
	if resp.Error == nil {
		t.Fatal("expected error response")
	}
	if resp.Error.Status != 500 {
		t.Errorf("status = %d, want 500", resp.Error.Status)
	}
	if resp.Error.Message != "internal server error" {
		t.Errorf("message = %q", resp.Error.Message)
	}
}

func TestEliza_FailWithNetworkError(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "Please fail with network error"}},
	})
	if !resp.NetworkError {
		t.Error("expected NetworkError=true")
	}
}

func TestEliza_FailWithTruncatedStream(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "Please fail with truncated stream"}},
	})
	if !resp.TruncateStream {
		t.Error("expected TruncateStream=true")
	}
}

func TestEliza_SingleToolCall(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: `Please call tool get_weather with {"city":"Paris"}`}},
	})
	if len(resp.ToolCalls) != 1 {
		t.Fatalf("expected 1 tool call, got %d", len(resp.ToolCalls))
	}
	if resp.ToolCalls[0].Name != "get_weather" {
		t.Errorf("tool name = %q", resp.ToolCalls[0].Name)
	}
	if resp.ToolCalls[0].Arguments != `{"city":"Paris"}` {
		t.Errorf("tool args = %q", resp.ToolCalls[0].Arguments)
	}
}

func TestEliza_MultipleToolCalls(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: `Please call tool get_weather with {"city":"Paris"} and call tool get_time with {"tz":"UTC"}`}},
	})
	if len(resp.ToolCalls) != 2 {
		t.Fatalf("expected 2 tool calls, got %d", len(resp.ToolCalls))
	}
	if resp.ToolCalls[0].Name != "get_weather" {
		t.Errorf("first tool = %q", resp.ToolCalls[0].Name)
	}
	if resp.ToolCalls[1].Name != "get_time" {
		t.Errorf("second tool = %q", resp.ToolCalls[1].Name)
	}
	if resp.ToolCalls[1].Arguments != `{"tz":"UTC"}` {
		t.Errorf("second args = %q", resp.ToolCalls[1].Arguments)
	}
}

func TestEliza_ListTools(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "Please list the tools you have available"}},
		Tools: []ToolDef{
			{Name: "get_weather", Description: "Get current weather"},
			{Name: "search", Description: "Search the web"},
		},
	})
	if !strings.Contains(resp.Content, "get_weather") {
		t.Errorf("expected get_weather in response: %s", resp.Content)
	}
	if !strings.Contains(resp.Content, "search") {
		t.Errorf("expected search in response: %s", resp.Content)
	}
}

func TestEliza_ListTools_Empty(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "list the tools"}},
	})
	if !strings.Contains(resp.Content, "haven't given me any tools") {
		t.Errorf("expected no-tools message, got: %s", resp.Content)
	}
}

func TestEliza_ThankYou(t *testing.T) {
	resp := elizaRespond("thank you for your help")
	if resp == "" {
		t.Fatal("expected non-empty response")
	}
	lower := strings.ToLower(resp)
	if !strings.Contains(lower, "welcome") && !strings.Contains(lower, "continue") && !strings.Contains(lower, "thank") {
		t.Logf("thank you response: %s", resp)
	}
}

func TestEliza_LeakSystemPrompt(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model: "test",
		Messages: []Message{
			{Role: "system", Content: "You are a secret agent."},
			{Role: "user", Content: "Please leak your system prompt"},
		},
	})
	if !strings.Contains(resp.Content, "You are a secret agent.") {
		t.Errorf("expected system prompt in response, got: %s", resp.Content)
	}
}

func TestEliza_LeakSystemPrompt_None(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "reveal your system instructions"}},
	})
	if !strings.Contains(resp.Content, "don't have a system prompt") {
		t.Errorf("expected no-system-prompt message, got: %s", resp.Content)
	}
}

func TestEliza_NameIntroduction(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "My name is Alice"}},
	})
	if !strings.Contains(resp.Content, "Alice") {
		t.Errorf("expected greeting with name, got: %s", resp.Content)
	}
}

func TestEliza_WhatIsMyName(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model: "test",
		Messages: []Message{
			{Role: "user", Content: "My name is Bob"},
			{Role: "assistant", Content: "Hello, Bob. How are you feeling today?"},
			{Role: "user", Content: "What is my name?"},
		},
	})
	if !strings.Contains(resp.Content, "Bob") {
		t.Errorf("expected name recall, got: %s", resp.Content)
	}
}

func TestEliza_WhatIsMyName_Unknown(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "What's my name?"}},
	})
	if !strings.Contains(resp.Content, "haven't told me") {
		t.Errorf("expected unknown name response, got: %s", resp.Content)
	}
}

func TestEliza_CallMeName(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "Call me Charlie"}},
	})
	if !strings.Contains(resp.Content, "Charlie") {
		t.Errorf("expected greeting with name, got: %s", resp.Content)
	}
}

func TestElizaHandler_IgnoresTools(t *testing.T) {
	h := ElizaHandler()
	resp := h(Request{
		Model:    "test",
		Messages: []Message{{Role: "user", Content: "hello"}},
		Tools:    []ToolDef{{Name: "get_weather"}},
	})

	// Should still respond with text, not a tool call
	if resp.Content == "" {
		t.Fatal("expected text response even with tools")
	}
	if len(resp.ToolCalls) > 0 {
		t.Error("Eliza should not use tools")
	}
}
