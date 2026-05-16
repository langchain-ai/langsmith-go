package mockllm

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestOpenAIServer_ChatCompletion_NonStreaming(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}],"stream":false}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	if result["model"] != "gpt-4o-mini" {
		t.Errorf("model = %v", result["model"])
	}
	choices := result["choices"].([]any)
	if len(choices) == 0 {
		t.Fatal("expected at least one choice")
	}
	choice := choices[0].(map[string]any)
	msg := choice["message"].(map[string]any)
	if msg["content"] == nil || msg["content"] == "" {
		t.Error("expected non-empty content")
	}

	reqs := srv.Requests()
	if len(reqs) != 1 {
		t.Fatalf("captured %d requests, want 1", len(reqs))
	}
}

func TestOpenAIServer_ChatCompletion_ToolCalls(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"weather?"}],"stream":false,"tools":[{"type":"function","function":{"name":"get_weather"}}]}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	choices := result["choices"].([]any)
	choice := choices[0].(map[string]any)
	if choice["finish_reason"] != "tool_calls" {
		t.Errorf("finish_reason = %v, want tool_calls", choice["finish_reason"])
	}
	msg := choice["message"].(map[string]any)
	toolCalls, ok := msg["tool_calls"].([]any)
	if !ok || len(toolCalls) == 0 {
		t.Fatal("expected tool_calls in message")
	}
}

func TestOpenAIServer_ChatCompletion_Streaming(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}],"stream":true}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Errorf("content-type = %q, want text/event-stream", ct)
	}

	data, _ := io.ReadAll(resp.Body)
	sseBody := string(data)

	if !strings.Contains(sseBody, "Hello") {
		t.Error("expected 'Hello' in stream")
	}
	if !strings.Contains(sseBody, "[DONE]") {
		t.Error("expected [DONE] sentinel")
	}
}

func TestOpenAIServer_ChatCompletion_Streaming_ToolCalls(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"weather?"}],"stream":true,"tools":[{"type":"function","function":{"name":"get_weather"}}]}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	sseBody := string(data)

	if !strings.Contains(sseBody, "get_weather") {
		t.Error("expected tool name in stream")
	}
	if !strings.Contains(sseBody, "call_mock_1") {
		t.Error("expected tool call ID in stream")
	}
}

func TestOpenAIServer_Responses_NonStreaming(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o","input":"What is Go?","stream":false}`
	resp, err := http.Post(srv.URL()+"/v1/responses", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	if result["object"] != "response" {
		t.Errorf("object = %v, want response", result["object"])
	}
	output := result["output"].([]any)
	if len(output) == 0 {
		t.Fatal("expected non-empty output")
	}
}

func TestOpenAIServer_Responses_Streaming(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o","input":"What is Go?","stream":true}`
	resp, err := http.Post(srv.URL()+"/v1/responses", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	sseBody := string(data)

	if !strings.Contains(sseBody, "response.completed") {
		t.Error("expected response.completed event")
	}
	if !strings.Contains(sseBody, "[DONE]") {
		t.Error("expected [DONE]")
	}
}

func TestOpenAIServer_InvalidAPIKey(t *testing.T) {
	srv := NewOpenAIServer()
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/chat/completions", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer sk-invalid-test")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestOpenAIServer_CustomHandler(t *testing.T) {
	srv := NewOpenAIServer(WithHandler(ErrorHandler(429, "rate limited")))
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("status = %d, want 429", resp.StatusCode)
	}
}

func TestOpenAIServer_NetworkError(t *testing.T) {
	srv := NewOpenAIServer(WithHandler(NetworkErrorHandler()))
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		// Connection reset / EOF is the expected network error
		return
	}
	defer resp.Body.Close()
	_, err = io.ReadAll(resp.Body)
	if err == nil {
		t.Error("expected network error (connection reset or EOF)")
	}
}

func TestOpenAIServer_TruncatedStream(t *testing.T) {
	srv := NewOpenAIServer(WithHandler(TruncatedStreamHandler("partial")))
	defer srv.Close()

	body := `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}],"stream":true}`
	resp, err := http.Post(srv.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	sseBody := string(data)

	if strings.Contains(sseBody, "[DONE]") {
		t.Error("truncated stream should NOT contain [DONE]")
	}
}

func TestOpenAIServer_SameHandlerBothProviders(t *testing.T) {
	// Verify the same Handler produces correct responses from both servers
	h := StaticHandler("shared response")

	anthropic := NewAnthropicServer(WithHandler(h))
	defer anthropic.Close()
	openai := NewOpenAIServer(WithHandler(h))
	defer openai.Close()

	// Anthropic
	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	resp, err := http.Post(anthropic.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("anthropic POST: %v", err)
	}
	var aResult map[string]any
	json.NewDecoder(resp.Body).Decode(&aResult)
	resp.Body.Close()

	aContent := aResult["content"].([]any)
	aBlock := aContent[0].(map[string]any)
	if aBlock["text"] != "shared response" {
		t.Errorf("anthropic text = %v", aBlock["text"])
	}

	// OpenAI
	body = `{"model":"gpt-4o-mini","messages":[{"role":"user","content":"hi"}]}`
	resp, err = http.Post(openai.URL()+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("openai POST: %v", err)
	}
	var oResult map[string]any
	json.NewDecoder(resp.Body).Decode(&oResult)
	resp.Body.Close()

	choices := oResult["choices"].([]any)
	msg := choices[0].(map[string]any)["message"].(map[string]any)
	if msg["content"] != "shared response" {
		t.Errorf("openai content = %v", msg["content"])
	}
}
