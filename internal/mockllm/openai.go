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

// OpenAIServer is a configurable mock OpenAI API server.
type OpenAIServer struct {
	Server  *httptest.Server
	handler Handler

	mu       sync.Mutex
	requests []Request
}

// NewOpenAIServer creates and starts a mock OpenAI API server.
// If no handler is provided via WithHandler, DefaultHandler is used.
func NewOpenAIServer(opts ...ServerOption) *OpenAIServer {
	s := &OpenAIServer{}
	for _, opt := range opts {
		opt.applyOpenAI(s)
	}
	if s.handler == nil {
		s.handler = DefaultHandler()
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/v1/chat/completions", s.handleChatCompletions)
	mux.HandleFunc("/v1/responses", s.handleResponses)
	s.Server = httptest.NewServer(mux)
	return s
}

// URL returns the base URL of the mock server.
func (s *OpenAIServer) URL() string {
	return s.Server.URL
}

// Close shuts down the server.
func (s *OpenAIServer) Close() {
	s.Server.Close()
}

// Requests returns all captured requests (in provider-agnostic form).
func (s *OpenAIServer) Requests() []Request {
	s.mu.Lock()
	defer s.mu.Unlock()
	cp := make([]Request, len(s.requests))
	copy(cp, s.requests)
	return cp
}

func (s *OpenAIServer) checkAuth(w http.ResponseWriter, r *http.Request) bool {
	auth := r.Header.Get("Authorization")
	if strings.Contains(auth, "invalid") {
		writeOpenAIError(w, http.StatusUnauthorized, "invalid_api_key", "Incorrect API key provided")
		return false
	}
	return true
}

func (s *OpenAIServer) parseRequest(r *http.Request) (Request, bool, error) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		return Request{}, false, err
	}

	var raw struct {
		Model    string           `json:"model"`
		Messages []map[string]any `json:"messages,omitempty"`
		Input    any              `json:"input,omitempty"` // Responses API
		Stream   bool             `json:"stream"`
		Tools    []map[string]any `json:"tools,omitempty"`
	}
	if err := json.Unmarshal(body, &raw); err != nil {
		return Request{}, false, err
	}

	req := Request{Model: raw.Model}

	// Chat completions messages
	for _, m := range raw.Messages {
		role, _ := m["role"].(string)
		content, _ := m["content"].(string)
		req.Messages = append(req.Messages, Message{Role: role, Content: content})
	}

	// Responses API input
	if raw.Input != nil && len(req.Messages) == 0 {
		switch v := raw.Input.(type) {
		case string:
			req.Messages = append(req.Messages, Message{Role: "user", Content: v})
		case []any:
			for _, item := range v {
				if m, ok := item.(map[string]any); ok {
					role, _ := m["role"].(string)
					content, _ := m["content"].(string)
					req.Messages = append(req.Messages, Message{Role: role, Content: content})
				}
			}
		}
	}

	// Tools
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

// --- Chat Completions ---

func (s *OpenAIServer) handleChatCompletions(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeOpenAIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	if !s.checkAuth(w, r) {
		return
	}

	req, streaming, err := s.parseRequest(r)
	if err != nil {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "invalid JSON")
		return
	}

	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()

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
		if resp.TruncateStream {
			hijackAndClose(w)
			return
		}
		writeChatJSONResponse(w, req.Model, resp)
	}
}

func writeOpenAIError(w http.ResponseWriter, status int, code, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(map[string]any{
		"error": map[string]any{
			"message": message,
			"type":    "invalid_request_error",
			"code":    code,
		},
	})
}

func writeChatJSONResponse(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "application/json")

	message := map[string]any{"role": "assistant"}
	finishReason := "stop"

	if len(resp.ToolCalls) > 0 {
		message["content"] = nil
		tcs := make([]map[string]any, len(resp.ToolCalls))
		for i, tc := range resp.ToolCalls {
			tcs[i] = map[string]any{
				"id":   tc.ID,
				"type": "function",
				"function": map[string]any{
					"name":      tc.Name,
					"arguments": tc.Arguments,
				},
			}
		}
		message["tool_calls"] = tcs
		finishReason = "tool_calls"
	} else {
		message["content"] = resp.Content
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":      "chatcmpl-mock-001",
		"object":  "chat.completion",
		"created": 1700000000,
		"model":   model,
		"choices": []map[string]any{
			{
				"index":         0,
				"message":       message,
				"finish_reason": finishReason,
			},
		},
		"usage": map[string]any{
			"prompt_tokens":     resp.InputTokens,
			"completion_tokens": resp.OutputTokens,
			"total_tokens":      resp.InputTokens + resp.OutputTokens,
		},
	})
}

func writeChatStreamingResponse(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	writeSSE := func(data string) {
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	id := "chatcmpl-mock-stream"

	if len(resp.ToolCalls) > 0 {
		// First chunk: role + first tool call header
		for i, tc := range resp.ToolCalls {
			writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{"role":"assistant","tool_calls":[{"index":%d,"id":"%s","type":"function","function":{"name":"%s","arguments":""}}]}}]}`, id, model, i, tc.ID, tc.Name))
			// Stream arguments in two chunks
			args := tc.Arguments
			mid := len(args) / 2
			if mid > 0 {
				escaped1, _ := json.Marshal(args[:mid])
				escaped2, _ := json.Marshal(args[mid:])
				writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{"tool_calls":[{"index":%d,"function":{"arguments":%s}}]}}]}`, id, model, i, escaped1))
				writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{"tool_calls":[{"index":%d,"function":{"arguments":%s}}]}}]}`, id, model, i, escaped2))
			}
		}
		if resp.TruncateStream {
			hijackAndClose(w)
			return
		}
		writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":%d,"completion_tokens":%d}}`, id, model, resp.InputTokens, resp.OutputTokens))
	} else {
		// Role chunk
		writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{"role":"assistant","content":""}}]}`, id, model))
		// Text chunks
		chunks := splitText(resp.Content, 3)
		for _, chunk := range chunks {
			if resp.StreamDelay != nil {
				if d := resp.StreamDelay(); d > 0 {
					time.Sleep(d)
				}
			}
			escaped, _ := json.Marshal(chunk)
			writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{"content":%s}}]}`, id, model, escaped))
		}
		if resp.TruncateStream {
			hijackAndClose(w)
			return
		}
		// Final chunk with usage
		writeSSE(fmt.Sprintf(`{"id":"%s","object":"chat.completion.chunk","model":"%s","choices":[{"index":0,"delta":{},"finish_reason":"stop"}],"usage":{"prompt_tokens":%d,"completion_tokens":%d}}`, id, model, resp.InputTokens, resp.OutputTokens))
	}

	writeSSE("[DONE]")
}

// --- Responses API ---

func (s *OpenAIServer) handleResponses(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		writeOpenAIError(w, http.StatusMethodNotAllowed, "method_not_allowed", "method not allowed")
		return
	}
	if !s.checkAuth(w, r) {
		return
	}

	req, streaming, err := s.parseRequest(r)
	if err != nil {
		writeOpenAIError(w, http.StatusBadRequest, "invalid_request_error", "invalid JSON")
		return
	}

	s.mu.Lock()
	s.requests = append(s.requests, req)
	s.mu.Unlock()

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
		writeResponsesStreaming(w, req.Model, resp)
	} else {
		if resp.TruncateStream {
			hijackAndClose(w)
			return
		}
		writeResponsesJSON(w, req.Model, resp)
	}
}

func writeResponsesJSON(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "application/json")

	var output []map[string]any
	if resp.Content != "" {
		output = append(output, map[string]any{
			"type": "message",
			"content": []map[string]any{
				{"type": "output_text", "text": resp.Content},
			},
		})
	}
	for _, tc := range resp.ToolCalls {
		output = append(output, map[string]any{
			"type":      "function_call",
			"name":      tc.Name,
			"arguments": tc.Arguments,
			"call_id":   tc.ID,
		})
	}

	json.NewEncoder(w).Encode(map[string]any{
		"id":     "resp_mock_001",
		"object": "response",
		"model":  model,
		"output": output,
		"usage": map[string]any{
			"input_tokens":  resp.InputTokens,
			"output_tokens": resp.OutputTokens,
		},
	})
}

func writeResponsesStreaming(w http.ResponseWriter, model string, resp Response) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	w.Header().Set("Connection", "keep-alive")

	flusher, ok := w.(http.Flusher)
	if !ok {
		http.Error(w, "streaming not supported", http.StatusInternalServerError)
		return
	}

	writeSSE := func(data string) {
		fmt.Fprintf(w, "data: %s\n\n", data)
		flusher.Flush()
	}

	// Build the full response object for response.completed
	var output []map[string]any
	if resp.Content != "" {
		output = append(output, map[string]any{
			"type": "message",
			"content": []map[string]any{
				{"type": "output_text", "text": resp.Content},
			},
		})
	}

	completedResp := map[string]any{
		"id":     "resp_mock_stream_001",
		"model":  model,
		"output": output,
		"usage": map[string]any{
			"input_tokens":  resp.InputTokens,
			"output_tokens": resp.OutputTokens,
		},
	}

	writeSSE(fmt.Sprintf(`{"type":"response.created","response":{"id":"resp_mock_stream_001","model":"%s"}}`, model))
	completedJSON, _ := json.Marshal(map[string]any{
		"type":     "response.completed",
		"response": completedResp,
	})
	writeSSE(string(completedJSON))
	writeSSE("[DONE]")
}
