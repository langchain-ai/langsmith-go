package mockllm

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"

	"compress/gzip"

	"github.com/klauspost/compress/zstd"
)

// CombinedServer serves both Anthropic and OpenAI APIs on a single port,
// backed by a shared Handler. It also serves GET /v1/models for client
// discovery (e.g. Codex).
type CombinedServer struct {
	Server  *httptest.Server
	handler Handler

	mu       sync.Mutex
	requests []Request
}

// NewCombinedServer creates and starts a combined mock server.
// If no handler is provided, ElizaHandler is used.
func NewCombinedServer(opts ...ServerOption) *CombinedServer {
	s := &CombinedServer{}
	for _, opt := range opts {
		opt.applyCombined(s)
	}
	if s.handler == nil {
		s.handler = ElizaHandler()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/messages", s.handleAnthropic)
	mux.HandleFunc("/v1/chat/completions", s.handleChatCompletions)
	mux.HandleFunc("/v1/responses", s.handleResponses)
	mux.HandleFunc("/v1/models", s.handleModels)
	s.Server = httptest.NewServer(mux)
	return s
}

// URL returns the base URL of the server.
func (s *CombinedServer) URL() string { return s.Server.URL }

// Close shuts down the server.
func (s *CombinedServer) Close() { s.Server.Close() }

// Requests returns all captured requests.
func (s *CombinedServer) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]Request, len(s.requests))
	copy(cp, s.requests)
	return cp
}

func (s *CombinedServer) record(req Request) {
	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()
}

// --- Anthropic /v1/messages ---

func (s *CombinedServer) handleAnthropic(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeAnthropicError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}

	body, err := decompressRequestBody(r)
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

	apiKey := r.Header.Get("x-api-key")
	if strings.HasPrefix(apiKey, "sk-ant-invalid") {
		writeAnthropicError(w, http.StatusUnauthorized, "authentication_error", "invalid x-api-key")
		return
	}

	req := Request{Model: raw.Model}
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

	s.record(req)
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
		writeAnthropicJSON(w, raw.Model, resp)
	}
}

// --- OpenAI /v1/chat/completions ---

func (s *CombinedServer) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeOpenAIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	if !checkOpenAIAuth(w, r) {
		return
	}

	req, streaming, err := parseOpenAIBody(r)
	if err != nil {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "invalid JSON")
		return
	}

	s.record(req)
	resp := s.handler(req)

	if resp.NetworkError {
		hijackAndClose(w)
		return
	}
	if resp.Error != nil {
		writeOpenAIError(w, resp.Error.Status, "api_error", resp.Error.Message)
		return
	}

	if streaming {
		writeChatStreamingResponse(w, req.Model, resp)
	} else {
		writeChatJSONResponse(w, req.Model, resp)
	}
}

// --- OpenAI /v1/responses ---

func (s *CombinedServer) handleResponses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeOpenAIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	if !checkOpenAIAuth(w, r) {
		return
	}

	req, streaming, err := parseOpenAIBody(r)
	if err != nil {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "invalid JSON")
		return
	}

	s.record(req)
	resp := s.handler(req)

	if resp.NetworkError {
		hijackAndClose(w)
		return
	}
	if resp.Error != nil {
		writeOpenAIError(w, resp.Error.Status, "api_error", resp.Error.Message)
		return
	}

	if streaming {
		writeResponsesStreamingCodex(w, req.Model, resp)
	} else {
		writeResponsesJSON(w, req.Model, resp)
	}
}

// --- Models ---

func (s *CombinedServer) handleModels(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(map[string]any{
		"object": "list",
		"data": []map[string]any{
			{"id": "eliza", "object": "model", "created": 1700000000, "owned_by": "eliza"},
		},
	})
}

// --- Shared helpers ---

func checkOpenAIAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if strings.Contains(auth, "invalid") {
		writeOpenAIError(w, http.StatusUnauthorized, "invalid_api_key", "Incorrect API key provided")
		return false
	}
	return true
}

func parseOpenAIBody(r *http.Request) (Request, bool, error) {
	body, err := decompressRequestBody(r)
	if err != nil {
		return Request{}, false, err
	}

	var raw struct {
		Model        string           `json:"model"`
		Messages     []map[string]any `json:"messages,omitempty"`
		Input        any              `json:"input,omitempty"`
		Instructions string           `json:"instructions,omitempty"`
		Stream       bool             `json:"stream"`
		Tools        []map[string]any `json:"tools,omitempty"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Request{}, false, err
	}

	req := Request{Model: raw.Model}

	if raw.Instructions != "" {
		req.Messages = append(req.Messages, Message{Role: "system", Content: raw.Instructions})
	}

	for _, m := range raw.Messages {
		role, _ := m["role"].(string)
		content, _ := m["content"].(string)
		req.Messages = append(req.Messages, Message{Role: role, Content: content})
	}

	if raw.Input != nil && len(raw.Messages) == 0 {
		switch v := raw.Input.(type) {
		case string:
			req.Messages = append(req.Messages, Message{Role: "user", Content: v})
		case []any:
			for _, item := range v {
				m, ok := item.(map[string]any)
				if !ok {
					continue
				}
				role, _ := m["role"].(string)
				switch itemType, _ := m["type"].(string); itemType {
				case "message":
					if contentArr, ok := m["content"].([]any); ok {
						for _, c := range contentArr {
							if cm, ok := c.(map[string]any); ok {
								if t, _ := cm["type"].(string); t == "input_text" {
									text, _ := cm["text"].(string)
									req.Messages = append(req.Messages, Message{Role: role, Content: text})
								}
							}
						}
					}
					if content, ok := m["content"].(string); ok {
						req.Messages = append(req.Messages, Message{Role: role, Content: content})
					}
				case "function_call_output":
					output, _ := m["output"].(string)
					req.Messages = append(req.Messages, Message{Role: "tool", Content: output})
				default:
					if role != "" {
						if content, ok := m["content"].(string); ok {
							req.Messages = append(req.Messages, Message{Role: role, Content: content})
						}
					}
				}
			}
		}
	}

	for _, t := range raw.Tools {
		td := ToolDef{}
		if fn, ok := t["function"].(map[string]any); ok {
			td.Name, _ = fn["name"].(string)
			td.Description, _ = fn["description"].(string)
		} else {
			td.Name, _ = t["name"].(string)
			td.Description, _ = t["description"].(string)
		}
		req.Tools = append(req.Tools, td)
	}

	return req, raw.Stream, nil
}

func decompressRequestBody(r *http.Request) ([]byte, error) {
	var reader io.Reader = r.Body
	switch r.Header.Get("Content-Encoding") {
	case "zstd":
		dec, err := zstd.NewReader(r.Body)
		if err != nil {
			return nil, fmt.Errorf("zstd init: %w", err)
		}
		defer dec.Close()
		reader = dec
	case "gzip":
		dec, err := gzip.NewReader(r.Body)
		if err != nil {
			return nil, fmt.Errorf("gzip init: %w", err)
		}
		defer dec.Close()
		reader = dec
	}
	return io.ReadAll(reader)
}

// writeResponsesStreamingCodex writes Responses API SSE with proper event:
// prefix lines as expected by Codex and the OpenAI SDK.
func writeResponsesStreamingCodex(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	sseEvent := func(eventType, data string) {
		fmt.Fprintf(w, "event: %s\ndata: %s\n\n", eventType, data)
		flusher.Flush()
	}

	respID := "resp_mock_stream"

	created, _ := json.Marshal(map[string]any{
		"type":     "response.created",
		"response": map[string]any{"id": respID},
	})
	sseEvent("response.created", string(created))

	outputItem, _ := json.Marshal(map[string]any{
		"type": "response.output_item.done",
		"item": map[string]any{
			"type": "message",
			"role": "assistant",
			"id":   "msg_mock_001",
			"content": []map[string]any{
				{"type": "output_text", "text": resp.Content},
			},
		},
	})
	sseEvent("response.output_item.done", string(outputItem))

	completed, _ := json.Marshal(map[string]any{
		"type": "response.completed",
		"response": map[string]any{
			"id": respID,
			"usage": map[string]any{
				"input_tokens":  resp.InputTokens,
				"output_tokens": resp.OutputTokens,
				"total_tokens":  resp.InputTokens + resp.OutputTokens,
			},
		},
	})
	sseEvent("response.completed", string(completed))
}

// Add combined support to ServerOption.
func (o handlerOption) applyCombined(s *CombinedServer) { s.handler = o.h }
