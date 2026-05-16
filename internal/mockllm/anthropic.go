// Package mockllm provides mock HTTP servers for LLM provider APIs.
// These servers return realistic canned responses and are intended for
// testing tracing, proxy, and gateway code without requiring real API keys.
//
// Behavior is controlled by a [Handler] callback that takes a provider-agnostic
// [Request] and returns a provider-agnostic [Response]. The same Handler works
// with both the Anthropic and OpenAI mock servers — each server handles the
// wire-format translation (JSON structure, SSE streaming, etc.).
package mockllm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"time"
)

// AnthropicServer is a configurable mock Anthropic API server.
type AnthropicServer struct {
	Server  *httptest.Server
	handler Handler

	mu       sync.Mutex
	requests []Request
}

// NewAnthropicServer creates and starts a mock Anthropic API server.
// If no handler is provided via WithHandler, DefaultHandler is used.
func NewAnthropicServer(opts ...ServerOption) *AnthropicServer {
	s := &AnthropicServer{}
	for _, opt := range opts {
		opt.applyAnthropic(s)
	}
	if s.handler == nil {
		s.handler = DefaultHandler()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/messages", s.handleMessages)
	s.Server = httptest.NewServer(mux)
	return s
}

// URL returns the base URL of the mock server.
func (s *AnthropicServer) URL() string {
	return s.Server.URL
}

// Close shuts down the server.
func (s *AnthropicServer) Close() {
	s.Server.Close()
}

// Requests returns all captured requests (in provider-agnostic form).
func (s *AnthropicServer) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]Request, len(s.requests))
	copy(cp, s.requests)
	return cp
}

// ServerOption configures a mock server. Options that apply to both providers
// implement both applyAnthropic and applyOpenAI.
type ServerOption interface {
	applyAnthropic(*AnthropicServer)
	applyOpenAI(*OpenAIServer)
	applyCombined(*CombinedServer)
}

type handlerOption struct{ h Handler }

func (o handlerOption) applyAnthropic(s *AnthropicServer) { s.handler = o.h }
func (o handlerOption) applyOpenAI(s *OpenAIServer)       { s.handler = o.h }

// WithHandler sets the Handler for the mock server.
func WithHandler(h Handler) ServerOption {
	return handlerOption{h}
}

func (s *AnthropicServer) handleMessages(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAnthropicError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	body, err := io.ReadAll(r.Body)
	if err != nil {
		writeAnthropicError(w, http.StatusBadRequest, "invalid_request_error", "failed to read body")
		return
	}

	var raw struct {
		Model    string           `json:"model"`
		Messages []map[string]any `json:"messages"`
		System   any              `json:"system,omitempty"`
		Stream   bool             `json:"stream"`
		Tools    []map[string]any `json:"tools,omitempty"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		writeAnthropicError(w, http.StatusBadRequest, "invalid_request_error", "invalid JSON")
		return
	}

	// Check for invalid API key simulation
	apiKey := r.Header.Get("x-api-key")
	if strings.HasPrefix(apiKey, "sk-ant-invalid") {
		writeAnthropicError(w, http.StatusUnauthorized, "authentication_error", "invalid x-api-key")
		return
	}

	// Build provider-agnostic request
	req := Request{Model: raw.Model}
	// Prepend system message if present
	if sys, ok := raw.System.(string); ok && sys != "" {
		req.Messages = append(req.Messages, Message{Role: "system", Content: sys})
	}
	for _, m := range raw.Messages {
		role, _ := m["role"].(string)
		content, _ := m["content"].(string)
		req.Messages = append(req.Messages, Message{Role: role, Content: content})
	}
	for _, t := range raw.Tools {
		name, _ := t["name"].(string)
		desc, _ := t["description"].(string)
		req.Tools = append(req.Tools, ToolDef{Name: name, Description: desc})
	}

	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()

	// Call the handler
	resp := s.handler(req)

	if resp.NetworkError {
		hijackAndClose(w)
		return
	}

	if resp.Error != nil {
		writeAnthropicError(w, resp.Error.Status, "api_error", resp.Error.Message)
		return
	}

	if raw.Stream {
		writeAnthropicStreaming(w, raw.Model, resp)
	} else {
		if resp.TruncateStream {
			hijackAndClose(w)
			return
		}
		writeAnthropicJSON(w, raw.Model, resp)
	}
}

func writeAnthropicError(w http.ResponseWriter, status int, errType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"type": "error",
		"error": map[string]any{
			"type":    errType,
			"message": message,
		},
	})
}

func writeAnthropicJSON(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "application/json")

	var content []map[string]any
	if resp.Content != "" {
		content = append(content, map[string]any{"type": "text", "text": resp.Content})
	}
	for _, tc := range resp.ToolCalls {
		var input any
		if err := json.Unmarshal([]byte(tc.Arguments), &input); err != nil {
			input = tc.Arguments
		}
		content = append(content, map[string]any{
			"type":  "tool_use",
			"id":    tc.ID,
			"name":  tc.Name,
			"input": input,
		})
	}

	stopReason := resp.StopReason
	if stopReason == "" {
		if len(resp.ToolCalls) > 0 {
			stopReason = "tool_use"
		} else {
			stopReason = "end_turn"
		}
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":      "msg_mock_001",
		"type":    "message",
		"role":    "assistant",
		"model":   model,
		"content": content,
		"usage": map[string]any{
			"input_tokens":  resp.InputTokens,
			"output_tokens": resp.OutputTokens,
		},
		"stop_reason": stopReason,
	})
}

func writeAnthropicStreaming(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	writeSSE := func(eventType, data string) {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
		flusher.Flush()
	}

	writeSSE("message_start", fmt.Sprintf(`{"type":"message_start","message":{"id":"msg_mock_stream_001","type":"message","role":"assistant","content":[],"model":"%s","usage":{"input_tokens":%d}}}`, model, resp.InputTokens))

	idx := 0
	// Text content
	if resp.Content != "" {
		writeSSE("content_block_start", fmt.Sprintf(`{"type":"content_block_start","index":%d,"content_block":{"type":"text","text":""}}`, idx))
		// Split text into chunks for realistic streaming
		chunks := splitText(resp.Content, 3)
		for _, chunk := range chunks {
			if resp.StreamDelay != nil {
				if d := resp.StreamDelay(); d > 0 {
					time.Sleep(d)
				}
			}
			escaped, _ := json.Marshal(chunk)
			writeSSE("content_block_delta", fmt.Sprintf(`{"type":"content_block_delta","index":%d,"delta":{"type":"text_delta","text":%s}}`, idx, escaped))
		}
		writeSSE("content_block_stop", fmt.Sprintf(`{"type":"content_block_stop","index":%d}`, idx))
		idx++
	}

	// Tool calls
	for _, tc := range resp.ToolCalls {
		writeSSE("content_block_start", fmt.Sprintf(`{"type":"content_block_start","index":%d,"content_block":{"type":"tool_use","id":"%s","name":"%s"}}`, idx, tc.ID, tc.Name))
		// Stream arguments in two chunks
		args := tc.Arguments
		mid := len(args) / 2
		if mid > 0 {
			escaped1, _ := json.Marshal(args[:mid])
			escaped2, _ := json.Marshal(args[mid:])
			writeSSE("content_block_delta", fmt.Sprintf(`{"type":"content_block_delta","index":%d,"delta":{"type":"input_json_delta","partial_json":%s}}`, idx, escaped1))
			writeSSE("content_block_delta", fmt.Sprintf(`{"type":"content_block_delta","index":%d,"delta":{"type":"input_json_delta","partial_json":%s}}`, idx, escaped2))
		}
		writeSSE("content_block_stop", fmt.Sprintf(`{"type":"content_block_stop","index":%d}`, idx))
		idx++
	}

	stopReason := resp.StopReason
	if stopReason == "" {
		if len(resp.ToolCalls) > 0 {
			stopReason = "tool_use"
		} else {
			stopReason = "end_turn"
		}
	}

	if resp.TruncateStream {
		// Simulate mid-stream network failure: close connection before stop events
		hijackAndClose(w)
		return
	}

	writeSSE("message_delta", fmt.Sprintf(`{"type":"message_delta","delta":{"stop_reason":"%s"},"usage":{"output_tokens":%d}}`, stopReason, resp.OutputTokens))
	writeSSE("message_stop", `{"type":"message_stop"}`)
}

// splitText splits s into n roughly equal chunks.
func splitText(s string, n int) []string {
	if n <= 1 || len(s) <= n {
		return []string{s}
	}
	chunkSize := (len(s) + n - 1) / n
	var chunks []string
	for i := 0; i < len(s); i += chunkSize {
		end := i + chunkSize
		if end > len(s) {
			end = len(s)
		}
		chunks = append(chunks, s[i:end])
	}
	return chunks
}
