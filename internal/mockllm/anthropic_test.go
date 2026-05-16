package mockllm

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

func TestAnthropicServer_NonStreaming(t *testing.T) {
	srv := NewAnthropicServer()
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100,"stream":false}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("status = %d, want 200", resp.StatusCode)
	}

	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatalf("decode: %v", err)
	}

	if result["model"] != "claude-sonnet-4-20250514" {
		t.Errorf("model = %v", result["model"])
	}
	if result["stop_reason"] != "end_turn" {
		t.Errorf("stop_reason = %v", result["stop_reason"])
	}
	content, ok := result["content"].([]any)
	if !ok || len(content) == 0 {
		t.Fatal("expected non-empty content")
	}

	reqs := srv.Requests()
	if len(reqs) != 1 {
		t.Fatalf("captured %d requests, want 1", len(reqs))
	}
	if reqs[0].Model != "claude-sonnet-4-20250514" {
		t.Errorf("captured model = %q", reqs[0].Model)
	}
}

func TestAnthropicServer_NonStreaming_ToolUse(t *testing.T) {
	srv := NewAnthropicServer()
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"weather?"}],"max_tokens":100,"stream":false,"tools":[{"name":"get_weather"}]}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	if result["stop_reason"] != "tool_use" {
		t.Errorf("stop_reason = %v, want tool_use", result["stop_reason"])
	}

	content := result["content"].([]any)
	foundToolUse := false
	for _, c := range content {
		block := c.(map[string]any)
		if block["type"] == "tool_use" {
			foundToolUse = true
			if block["name"] != "get_weather" {
				t.Errorf("tool name = %v", block["name"])
			}
		}
	}
	if !foundToolUse {
		t.Error("expected tool_use content block")
	}
}

func TestAnthropicServer_Streaming(t *testing.T) {
	srv := NewAnthropicServer()
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100,"stream":true}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if ct := resp.Header.Get("Content-Type"); !strings.Contains(ct, "text/event-stream") {
		t.Errorf("content-type = %q, want text/event-stream", ct)
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	sseBody := string(data)
	if !strings.Contains(sseBody, "message_start") {
		t.Error("expected message_start event")
	}
	if !strings.Contains(sseBody, "text_delta") {
		t.Error("expected text_delta events")
	}
	if !strings.Contains(sseBody, "message_stop") {
		t.Error("expected message_stop event")
	}
	if !strings.Contains(sseBody, "Hello") {
		t.Error("expected 'Hello' in streamed text")
	}
}

func TestAnthropicServer_Streaming_ToolUse(t *testing.T) {
	srv := NewAnthropicServer()
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"weather?"}],"max_tokens":100,"stream":true,"tools":[{"name":"get_weather"}]}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	sseBody := string(data)

	if !strings.Contains(sseBody, "input_json_delta") {
		t.Error("expected input_json_delta for tool use streaming")
	}
	if !strings.Contains(sseBody, "get_weather") {
		t.Error("expected tool name in stream")
	}
}

func TestAnthropicServer_InvalidAPIKey(t *testing.T) {
	srv := NewAnthropicServer()
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	req, _ := http.NewRequest("POST", srv.URL()+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "sk-ant-invalid-test")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status = %d, want 401", resp.StatusCode)
	}
}

func TestAnthropicServer_CustomHandler(t *testing.T) {
	srv := NewAnthropicServer(WithHandler(ErrorHandler(429, "rate limited")))
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusTooManyRequests {
		t.Errorf("status = %d, want 429", resp.StatusCode)
	}
}

func TestAnthropicServer_StaticHandler(t *testing.T) {
	srv := NewAnthropicServer(WithHandler(StaticHandler("custom response")))
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("POST: %v", err)
	}
	defer resp.Body.Close()

	var result map[string]any
	json.NewDecoder(resp.Body).Decode(&result)

	content := result["content"].([]any)
	block := content[0].(map[string]any)
	if block["text"] != "custom response" {
		t.Errorf("text = %v, want 'custom response'", block["text"])
	}
}

func TestAnthropicServer_NetworkError(t *testing.T) {
	srv := NewAnthropicServer(WithHandler(NetworkErrorHandler()))
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		// Connection reset / EOF is the expected network error
		return
	}
	defer resp.Body.Close()
	// If we got a response, reading the body should fail
	_, err = io.ReadAll(resp.Body)
	if err == nil {
		t.Error("expected network error (connection reset or EOF)")
	}
}

func TestAnthropicServer_TruncatedStream(t *testing.T) {
	srv := NewAnthropicServer(WithHandler(TruncatedStreamHandler("partial")))
	defer srv.Close()

	body := `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hi"}],"max_tokens":100,"stream":true}`
	resp, err := http.Post(srv.URL()+"/v1/messages", "application/json", strings.NewReader(body))
	if err != nil {
		// Connection error is acceptable
		return
	}
	defer resp.Body.Close()

	data, _ := io.ReadAll(resp.Body)
	sseBody := string(data)

	// Should have some content but NOT the final message_stop
	if strings.Contains(sseBody, "message_stop") {
		t.Error("truncated stream should NOT contain message_stop")
	}
}
