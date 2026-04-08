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

// TestE2E_ElizaUpstream tests the full proxy stack: real HTTP client →
// proxy → Eliza mock upstream, verifying the response passes through
// correctly for all supported endpoints.
func TestE2E_ElizaUpstream(t *testing.T) {
	upstream := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.ElizaHandler()))
	defer upstream.Close()

	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	defer tp.Shutdown(context.Background())

	p, err := proxy.New(proxy.Config{
		AnthropicUpstream: upstream.URL(),
		OpenAIUpstream:    upstream.URL(),
		TracerProvider:    tp,
	})
	if err != nil {
		t.Fatalf("proxy.New: %v", err)
	}
	proxyServer := httptest.NewServer(p)
	defer proxyServer.Close()

	t.Run("OpenAI_ChatCompletions", func(t *testing.T) {
		resp, err := http.Post(proxyServer.URL+"/v1/chat/completions", "application/json",
			strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}]}`))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Fatalf("status %d: %s", resp.StatusCode, data)
		}

		var result map[string]any
		json.Unmarshal(data, &result)
		choices, _ := result["choices"].([]any)
		if len(choices) == 0 {
			t.Fatal("expected choices")
		}
		msg := choices[0].(map[string]any)["message"].(map[string]any)
		content, _ := msg["content"].(string)
		if content == "" {
			t.Error("expected non-empty Eliza response")
		}
		t.Logf("Eliza says: %s", content)
	})

	t.Run("Anthropic_Messages", func(t *testing.T) {
		req, _ := http.NewRequest("POST", proxyServer.URL+"/v1/messages",
			strings.NewReader(`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"I am sad"}],"max_tokens":100}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "fake")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		if resp.StatusCode != 200 {
			t.Fatalf("status %d: %s", resp.StatusCode, data)
		}

		var result map[string]any
		json.Unmarshal(data, &result)
		blocks, _ := result["content"].([]any)
		if len(blocks) == 0 {
			t.Fatal("expected content blocks")
		}
		text, _ := blocks[0].(map[string]any)["text"].(string)
		if text == "" {
			t.Error("expected non-empty Eliza response")
		}
		t.Logf("Eliza says: %s", text)
	})

	t.Run("Anthropic_Streaming", func(t *testing.T) {
		req, _ := http.NewRequest("POST", proxyServer.URL+"/v1/messages",
			strings.NewReader(`{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"hello"}],"max_tokens":100,"stream":true}`))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("x-api-key", "fake")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		if !strings.Contains(string(data), "text_delta") {
			t.Error("expected text_delta events")
		}
		if !strings.Contains(string(data), "message_stop") {
			t.Error("expected message_stop event")
		}
	})

	t.Run("OpenAI_Streaming", func(t *testing.T) {
		resp, err := http.Post(proxyServer.URL+"/v1/chat/completions", "application/json",
			strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"hello"}],"stream":true}`))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		if !strings.Contains(string(data), "[DONE]") {
			t.Error("expected [DONE]")
		}
	})

	t.Run("Models_Passthrough", func(t *testing.T) {
		resp, err := http.Get(proxyServer.URL + "/v1/models")
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		var result map[string]any
		json.Unmarshal(data, &result)
		if result["object"] != "list" {
			t.Errorf("expected models list, got: %s", data)
		}
	})

	t.Run("Eliza_SpecialCommand", func(t *testing.T) {
		resp, err := http.Post(proxyServer.URL+"/v1/chat/completions", "application/json",
			strings.NewReader(`{"model":"gpt-4o","messages":[{"role":"user","content":"What can you do?"}]}`))
		if err != nil {
			t.Fatalf("request: %v", err)
		}
		defer resp.Body.Close()
		data, _ := io.ReadAll(resp.Body)

		if !strings.Contains(string(data), "special commands") {
			t.Error("expected Eliza help text through proxy")
		}
	})
}
