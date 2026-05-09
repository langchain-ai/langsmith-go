// Package openaitrace provides OpenTelemetry middleware for OpenAI API requests.
package traceopenai

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go/internal/traceutil"
)

// MiddlewareNext is a function which is called by the middleware to pass an HTTP request
// to the next stage in the middleware chain (the actual HTTP transport).
type MiddlewareNext func(*http.Request) (*http.Response, error)

// Middleware adds OpenTelemetry tracing to OpenAI API requests.
// It intercepts the request, creates spans, extracts/records attributes,
// and then calls next to make the actual HTTP request.
// Uses the global tracer provider.
func Middleware(req *http.Request, next MiddlewareNext) (*http.Response, error) {
	return MiddlewareWithTracerProvider(req, next, nil)
}

// MiddlewareWithTracerProvider adds OpenTelemetry tracing to OpenAI API requests.
// It intercepts the request, creates spans, extracts/records attributes,
// and then calls next to make the actual HTTP request.
// If tp is nil, uses the global tracer provider.
func MiddlewareWithTracerProvider(req *http.Request, next MiddlewareNext, tp trace.TracerProvider) (*http.Response, error) {
	// Only trace known OpenAI-compatible API endpoints (path-based, works with
	// OpenRouter, Azure, local proxies, etc.)
	if !isOpenAIEndpoint(req.URL.Path) {
		return next(req)
	}

	ctx := req.Context()
	var tracer trace.Tracer
	if tp != nil {
		tracer = tp.Tracer("github.com/sashabaranov/go-openai")
	} else {
		tracer = otel.Tracer("github.com/sashabaranov/go-openai")
	}

	// Extract span context from request headers
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))

	// Capture parent span before creating child span (for token propagation)
	parentSpan := trace.SpanFromContext(ctx)

	spanName := getSpanName(req.URL.Path)
	if v := req.Context().Value(ctxKeyRunName); v != nil {
		if s, ok := v.(string); ok && s != "" {
			spanName = s
		}
	}

	// Build attributes for child span
	spanAttrs := []attribute.KeyValue{
		attribute.String("gen_ai.system", "openai"),
		attribute.String("gen_ai.operation.name", getOperationName(req.URL.Path)),
		attribute.String("http.method", req.Method),
		attribute.String("http.url", req.URL.String()),
	}

	// Propagate thread metadata from baggage
	// This ensures all child spans are part of the same thread
	// According to LangSmith docs: thread metadata must be on ALL spans including children
	// LangSmith looks for: session_id, thread_id, or conversation_id in span attributes/metadata
	// Set in multiple formats: standard (session_id), LangSmith metadata format (langsmith.metadata.session_id), and compatibility (session.id)
	bag := baggage.FromContext(ctx)
	// Check for thread metadata keys in baggage and propagate to span attributes
	if member := bag.Member("session_id"); member.Key() == "session_id" {
		value := member.Value()
		spanAttrs = append(spanAttrs,
			attribute.String("session_id", value),                    // Standard format per docs
			attribute.String("langsmith.metadata.session_id", value), // LangSmith metadata format
			attribute.String("session.id", value),                    // Compatibility format
		)
	}
	if member := bag.Member("thread_id"); member.Key() == "thread_id" {
		value := member.Value()
		spanAttrs = append(spanAttrs,
			attribute.String("thread_id", value),
			attribute.String("langsmith.metadata.thread_id", value),
			attribute.String("thread.id", value),
		)
	}
	if member := bag.Member("conversation_id"); member.Key() == "conversation_id" {
		value := member.Value()
		spanAttrs = append(spanAttrs,
			attribute.String("conversation_id", value),
			attribute.String("langsmith.metadata.conversation_id", value),
			attribute.String("conversation.id", value),
		)
	}

	// Start span (child span)
	ctx, span := tracer.Start(ctx, spanName, trace.WithAttributes(spanAttrs...))

	// Read request body if present
	var requestBody []byte
	if req.Body != nil {
		var err error
		requestBody, err = io.ReadAll(req.Body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, fmt.Sprintf("failed to read request body: %v", err))
			span.End()
			return next(req)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// Extract request attributes
	var streaming bool
	if len(requestBody) > 0 {
		reqFields := parseRequestBody(requestBody)
		if reqFields.inputMessages != "" {
			span.SetAttributes(attribute.String("gen_ai.prompt", reqFields.inputMessages))
		}
		if reqFields.model != "" {
			span.SetAttributes(attribute.String("gen_ai.request.model", reqFields.model))
		}
		streaming = reqFields.streaming
	}

	responsesAPI := strings.HasSuffix(req.URL.Path, "/responses")

	// Inject span context into request headers and update request context
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	req = req.WithContext(ctx)

	// Make the actual request via next middleware/transport
	resp, err := next(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return resp, err
	}

	br := traceutil.NewBufferedReader(resp.Body, func(r io.Reader, readErr error) {
		data, err := io.ReadAll(r)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			span.End()
			return
		}
		if len(data) == 0 {
			if resp.StatusCode >= 400 {
				apiErr := fmt.Errorf("HTTP %d", resp.StatusCode)
				span.RecordError(apiErr)
				span.SetStatus(codes.Error, apiErr.Error())
			}
			// readErr is the error that ended the read (e.g. context.Canceled); record it so run has real error
			if readErr != nil && readErr != io.EOF {
				span.RecordError(readErr)
				span.SetStatus(codes.Error, readErr.Error())
			}
			span.End()
			return
		}

		bodyText := string(data)
		if resp.StatusCode >= 400 {
			// Record an error so backends (e.g. LangSmith) show the trace as failed and populate run.error
			msg := bodyText
			if len(msg) > 500 {
				msg = msg[:500] + "..."
			}
			apiErr := fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
			span.RecordError(apiErr)
			span.SetStatus(codes.Error, apiErr.Error())
		}
		// Early stream termination: use real error from read (e.g. context.Canceled) or synthetic "Cancelled".
		// Chat Completions streams end with "data: [DONE]"; Responses API streams
		// end with a "response.completed" event, not [DONE].
		// See https://developers.openai.com/api/docs/guides/streaming-responses#read-the-responses
		var incompleteStream bool
		if resp.StatusCode < 400 && streaming {
			if responsesAPI {
				incompleteStream = !strings.Contains(bodyText, `"response.completed"`)
			} else {
				incompleteStream = !strings.Contains(bodyText, "[DONE]")
			}
		}
		if incompleteStream {
			endErr := readErr
			if endErr == nil || endErr == io.EOF {
				endErr = fmt.Errorf("Cancelled")
			}
			span.RecordError(endErr)
			span.SetStatus(codes.Error, endErr.Error())
		}

		var completion string
		var usage usageInfo
		if responsesAPI {
			if streaming {
				completion, usage = extractStreamingResponsesCompletion(data)
			} else {
				completion, usage = extractResponsesCompletion(data)
			}
		} else if streaming {
			completion, usage = extractStreamingCompletion(data)
		} else {
			completion, usage = extractCompletionFromResponse(data)
		}

		if completion != "" {
			span.SetAttributes(attribute.String("gen_ai.completion", completion))
		}
		if usage.InputTokens > 0 {
			inputTokens := int64(usage.InputTokens)
			span.SetAttributes(attribute.Int64("gen_ai.usage.input_tokens", inputTokens))
			if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
				parentSpan.SetAttributes(attribute.Int64("gen_ai.usage.input_tokens", inputTokens))
			}
		}
		if usage.OutputTokens > 0 {
			outputTokens := int64(usage.OutputTokens)
			span.SetAttributes(attribute.Int64("gen_ai.usage.output_tokens", outputTokens))
			if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
				parentSpan.SetAttributes(attribute.Int64("gen_ai.usage.output_tokens", outputTokens))
			}
		}
		if resp.StatusCode < 400 && !incompleteStream {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	})
	// LangSmith ingest reads new_token to derive first_token_time; skip on
	// HTTP errors so an error body doesn't inflate it.
	if streaming && resp.StatusCode < 400 {
		isMatch := isFirstContentChat
		if responsesAPI {
			isMatch = isFirstContentResponses
		}
		traceutil.OnFirstSSEMatch(br, isMatch, func() { span.AddEvent("new_token") })
	}
	resp.Body = br

	return resp, nil
}

// Chat streams open with a delta.role-only preamble; with n>1, content can
// land first on any choice index.
func isFirstContentChat(chunk map[string]any) bool {
	choices, ok := chunk["choices"].([]any)
	if !ok {
		return false
	}
	for _, raw := range choices {
		choice, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		// Legacy /v1/completions: text on the choice itself.
		if text, ok := choice["text"].(string); ok && text != "" {
			return true
		}
		delta, ok := choice["delta"].(map[string]any)
		if !ok {
			continue
		}
		if text, ok := delta["content"].(string); ok && text != "" {
			return true
		}
		if tcs, ok := delta["tool_calls"].([]any); ok && len(tcs) > 0 {
			return true
		}
	}
	return false
}

// Responses API uses *.delta event types for tokens; lifecycle envelopes
// (response.created, .in_progress, .completed, …) don't.
func isFirstContentResponses(chunk map[string]any) bool {
	t, _ := chunk["type"].(string)
	return strings.HasSuffix(t, ".delta")
}

// isOpenAIEndpoint returns true if the path matches a known OpenAI API endpoint.
func isOpenAIEndpoint(path string) bool {
	return strings.HasSuffix(path, "/chat/completions") ||
		strings.HasSuffix(path, "/completions") ||
		strings.HasSuffix(path, "/embeddings") ||
		strings.HasSuffix(path, "/responses")
}

// getSpanName returns an appropriate span name based on the API endpoint.
func getSpanName(path string) string {
	if strings.HasSuffix(path, "/chat/completions") {
		return "openai.chat.completion"
	}
	if strings.HasSuffix(path, "/completions") {
		return "openai.completion"
	}
	if strings.HasSuffix(path, "/embeddings") {
		return "openai.embedding"
	}
	if strings.HasSuffix(path, "/responses") {
		return "openai.responses"
	}
	return "openai.request"
}

// getOperationName returns the operation name for Gen AI semantic conventions.
func getOperationName(path string) string {
	if strings.HasSuffix(path, "/chat/completions") {
		return "chat"
	}
	if strings.HasSuffix(path, "/completions") {
		return "completion"
	}
	if strings.HasSuffix(path, "/embeddings") {
		return "embedding"
	}
	if strings.HasSuffix(path, "/responses") {
		return "responses"
	}
	return "request"
}

// requestFields holds fields extracted from the request body.
type requestFields struct {
	inputMessages string
	model         string
	streaming     bool
}

// parseRequestBody extracts input messages, model, and streaming flag from
// the request body.
func parseRequestBody(body []byte) requestFields {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return requestFields{}
	}

	var fields requestFields

	// Model
	fields.model, _ = req["model"].(string)

	// Streaming
	fields.streaming, _ = req["stream"].(bool)

	// Input messages — chat completions
	if messages, ok := req["messages"].([]any); ok && len(messages) > 0 {
		fields.inputMessages = marshalMessages(messages)
		return fields
	}

	// Input — Responses API (string or array)
	if input, ok := req["input"]; ok {
		var msgs []any
		if instr, ok := req["instructions"].(string); ok && instr != "" {
			msgs = append(msgs, map[string]any{"role": "system", "content": instr})
		}
		switch v := input.(type) {
		case string:
			msgs = append(msgs, map[string]any{"role": "user", "content": v})
		case []any:
			msgs = append(msgs, normalizeResponsesInput(v)...)
		}
		if len(msgs) > 0 {
			fields.inputMessages = marshalMessages(msgs)
		}
		return fields
	}

	// Input — legacy completions
	if prompt, ok := req["prompt"].(string); ok {
		fields.inputMessages = marshalMessages([]any{map[string]any{"role": "user", "content": prompt}})
	}

	return fields
}

// marshalMessages wraps a messages slice in {"messages":[...]} and marshals to JSON.
func marshalMessages(messages []any) string {
	out, err := json.Marshal(map[string]any{"messages": messages})
	if err != nil {
		return ""
	}
	return string(out)
}

// normalizeResponsesInput converts Responses API input items into
// chat-completions messages that LangSmith can render.
func normalizeResponsesInput(items []any) []any {
	out := make([]any, 0, len(items))
	for _, item := range items {
		m, ok := item.(map[string]any)
		if !ok {
			out = append(out, item)
			continue
		}

		itemType, _ := m["type"].(string)

		// Items with a role are messages.
		if role, ok := m["role"].(string); ok {
			msg := map[string]any{"role": role}
			switch content := m["content"].(type) {
			case string:
				msg["content"] = content
			case []any:
				if text := flattenContentParts(content); text != "" {
					msg["content"] = text
				} else {
					b, _ := json.Marshal(content)
					msg["content"] = string(b)
				}
			}
			out = append(out, msg)
			continue
		}

		// Non-message items: convert to chat-completions equivalents.
		switch itemType {
		case "function_call_output":
			out = append(out, toolOutputFromItem(m, "call_id"))
		case "computer_call_output":
			out = append(out, toolOutputFromItem(m, "call_id"))
		case "mcp_call_output":
			out = append(out, toolOutputFromItem(m, "id"))
		case "mcp_approval_response":
			msg := map[string]any{"role": "tool"}
			if reqID, ok := m["approval_request_id"].(string); ok {
				msg["tool_call_id"] = reqID
			}
			approved, _ := m["approve"].(bool)
			msg["content"] = fmt.Sprintf("approved: %v", approved)
			out = append(out, msg)
		case "reasoning":
			if text := extractSummaryText(m); text != "" {
				out = append(out, map[string]any{
					"role":    "assistant",
					"content": "[reasoning] " + text,
				})
			}
		case "item_reference", "compaction":
			continue
		default:
			if tc, ok := responsesItemToToolCall(itemType, m); ok {
				out = append(out, map[string]any{
					"role":       "assistant",
					"tool_calls": []any{tc},
				})
			} else if itemType != "" {
				out = append(out, map[string]any{
					"role":    "assistant",
					"content": "[" + itemType + "]",
				})
			}
		}
	}
	return out
}

// toolOutputFromItem builds a tool-output message from a Responses API output item.
func toolOutputFromItem(m map[string]any, idKey string) map[string]any {
	msg := map[string]any{"role": "tool"}
	if callID, ok := m[idKey].(string); ok {
		msg["tool_call_id"] = callID
	}
	switch v := m["output"].(type) {
	case string:
		msg["content"] = v
	case []any:
		if text := flattenContentParts(v); text != "" {
			msg["content"] = text
		} else {
			b, _ := json.Marshal(v)
			msg["content"] = string(b)
		}
	}
	return msg
}

// responsesItemToToolCall converts a Responses API tool item into a
// chat-completions tool_call. Returns false for unknown types.
func responsesItemToToolCall(itemType string, m map[string]any) (map[string]any, bool) {
	switch itemType {
	case "function_call":
		name, _ := m["name"].(string)
		args, _ := m["arguments"].(string)
		return chatToolCall(m["call_id"], name, args), true
	case "web_search_call":
		args := make(map[string]any)
		if status, ok := m["status"].(string); ok {
			args["status"] = status
		}
		if action, ok := m["action"].(map[string]any); ok {
			if q, ok := action["query"].(string); ok {
				args["query"] = q
			}
		}
		b, _ := json.Marshal(args)
		return chatToolCall(m["id"], "web_search", string(b)), true
	case "file_search_call":
		args := make(map[string]any)
		if q, ok := m["queries"].([]any); ok {
			args["queries"] = q
		}
		if r, ok := m["results"].([]any); ok {
			args["results"] = r
		}
		b, _ := json.Marshal(args)
		return chatToolCall(m["id"], "file_search", string(b)), true
	case "code_interpreter_call":
		args := make(map[string]any)
		if code, ok := m["code"].(string); ok {
			args["code"] = code
		}
		if r, ok := m["results"].([]any); ok {
			args["results"] = r
		}
		b, _ := json.Marshal(args)
		return chatToolCall(m["id"], "code_interpreter", string(b)), true
	case "computer_call":
		args := make(map[string]any)
		if action, ok := m["action"].(map[string]any); ok {
			args["action"] = action
		}
		b, _ := json.Marshal(args)
		return chatToolCall(m["call_id"], "computer", string(b)), true
	case "mcp_call":
		name := "mcp"
		if sl, ok := m["server_label"].(string); ok {
			if tn, ok := m["name"].(string); ok {
				name = sl + ":" + tn
			}
		}
		args, _ := m["arguments"].(string)
		if args == "" {
			args = "{}"
		}
		return chatToolCall(m["id"], name, args), true
	case "mcp_list_tools":
		return chatToolCall(m["id"], "mcp_list_tools", marshalToolArgs(m, "server_label")), true
	case "mcp_approval_request":
		name := "mcp_approval_request"
		if n, ok := m["name"].(string); ok {
			name = n
		}
		args, _ := m["arguments"].(string)
		if args == "" {
			args = "{}"
		}
		return chatToolCall(m["id"], name, args), true
	case "image_generation_call":
		return chatToolCall(m["id"], "image_generation", marshalToolArgs(m, "status")), true
	default:
		return nil, false
	}
}

// extractSummaryText extracts text from a reasoning item's summary.
func extractSummaryText(m map[string]any) string {
	summary, ok := m["summary"].([]any)
	if !ok {
		return ""
	}
	var b strings.Builder
	for _, s := range summary {
		sMap, ok := s.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := sMap["text"].(string); ok {
			if b.Len() > 0 {
				b.WriteByte('\n')
			}
			b.WriteString(text)
		}
	}
	return b.String()
}

// flattenContentParts extracts text from Responses API content part arrays.
func flattenContentParts(parts []any) string {
	var b strings.Builder
	for _, part := range parts {
		pm, ok := part.(map[string]any)
		if !ok {
			continue
		}
		switch pm["type"] {
		case "input_text", "output_text", "text":
			if text, ok := pm["text"].(string); ok {
				if b.Len() > 0 {
					b.WriteByte('\n')
				}
				b.WriteString(text)
			}
		case "refusal":
			if text, ok := pm["refusal"].(string); ok {
				if b.Len() > 0 {
					b.WriteByte('\n')
				}
				b.WriteString(text)
			}
		}
	}
	return b.String()
}

// extractStreamingCompletion parses an SSE response body, aggregates
// delta.content and delta.tool_calls across chunks, and extracts usage from
// the chunk that contains it (last chunk when stream_options.include_usage is set).
func extractStreamingCompletion(data []byte) (string, usageInfo) {
	chunks, err := traceutil.ParseSSEChunks(bytes.NewReader(data))
	if err != nil || len(chunks) == 0 {
		return "", usageInfo{}
	}

	var content strings.Builder
	var usage usageInfo

	// Track tool calls by index. Each entry holds id, type, function name,
	// and a builder that accumulates the streamed function arguments.
	type toolCallAcc struct {
		ID   string
		Type string
		Name string
		Args strings.Builder
	}
	var toolCalls []*toolCallAcc

	for _, chunk := range chunks {
		if choices, ok := chunk["choices"].([]any); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]any); ok {
				if delta, ok := choice["delta"].(map[string]any); ok {
					// Aggregate text content
					if text, ok := delta["content"].(string); ok {
						content.WriteString(text)
					}
					// Aggregate tool call deltas
					if tcs, ok := delta["tool_calls"].([]any); ok {
						for _, tc := range tcs {
							tcMap, ok := tc.(map[string]any)
							if !ok {
								continue
							}
							idx := 0
							if idxF, ok := tcMap["index"].(float64); ok {
								idx = int(idxF)
							}
							for len(toolCalls) <= idx {
								toolCalls = append(toolCalls, &toolCallAcc{})
							}
							if id, ok := tcMap["id"].(string); ok {
								toolCalls[idx].ID = id
							}
							if typ, ok := tcMap["type"].(string); ok {
								toolCalls[idx].Type = typ
							}
							if fn, ok := tcMap["function"].(map[string]any); ok {
								if name, ok := fn["name"].(string); ok {
									toolCalls[idx].Name = name
								}
								if args, ok := fn["arguments"].(string); ok {
									toolCalls[idx].Args.WriteString(args)
								}
							}
						}
					}
				}
			}
		}

		// Extract usage (present in the last chunk when include_usage is set)
		if usageMap, ok := chunk["usage"].(map[string]any); ok {
			if v, ok := usageMap["prompt_tokens"].(float64); ok {
				usage.InputTokens = int(v)
			}
			if v, ok := usageMap["completion_tokens"].(float64); ok {
				usage.OutputTokens = int(v)
			}
		}
	}

	text := content.String()

	// Build assistant message
	msg := map[string]any{"role": "assistant"}
	if text != "" {
		msg["content"] = text
	}
	if len(toolCalls) > 0 {
		tcOut := make([]map[string]any, len(toolCalls))
		for i, tc := range toolCalls {
			tcOut[i] = map[string]any{
				"id":   tc.ID,
				"type": tc.Type,
				"function": map[string]any{
					"name":      tc.Name,
					"arguments": tc.Args.String(),
				},
			}
		}
		msg["tool_calls"] = tcOut
	}

	if len(msg) == 1 {
		// Only "role" — no content or tool calls
		return "", usage
	}
	return marshalMessages([]any{msg}), usage
}

// usageInfo holds token usage information.
type usageInfo struct {
	InputTokens  int
	OutputTokens int
}

// extractCompletionFromResponse extracts the assistant message and usage from
// an OpenAI response. Returns a {"messages":[message]} JSON string.
func extractCompletionFromResponse(body []byte) (string, usageInfo) {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", usageInfo{}
	}

	var usage usageInfo
	if usageMap, ok := resp["usage"].(map[string]any); ok {
		if v, ok := usageMap["prompt_tokens"].(float64); ok {
			usage.InputTokens = int(v)
		}
		if v, ok := usageMap["completion_tokens"].(float64); ok {
			usage.OutputTokens = int(v)
		}
	}

	if choices, ok := resp["choices"].([]any); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]any); ok {
			// Chat completions — full message object (has role, content, tool_calls, etc.)
			if message, ok := choice["message"].(map[string]any); ok {
				return marshalMessages([]any{message}), usage
			}
			// Legacy completions — wrap text in an assistant message
			if text, ok := choice["text"].(string); ok {
				return marshalMessages([]any{map[string]any{"role": "assistant", "content": text}}), usage
			}
		}
	}

	return "", usage
}

// --- Responses API (/v1/responses) ---

// extractResponsesCompletion extracts completion and usage from a non-streaming
// Responses API response.
func extractResponsesCompletion(body []byte) (string, usageInfo) {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", usageInfo{}
	}
	return extractResponsesOutput(resp), extractResponsesUsage(resp)
}

// extractStreamingResponsesCompletion extracts completion and usage from a
// streaming Responses API response. The response.completed event contains the
// full response object.
func extractStreamingResponsesCompletion(data []byte) (string, usageInfo) {
	chunks, err := traceutil.ParseSSEChunks(bytes.NewReader(data))
	if err != nil || len(chunks) == 0 {
		return "", usageInfo{}
	}

	for _, chunk := range chunks {
		if msgType, _ := chunk["type"].(string); msgType == "response.completed" {
			if response, ok := chunk["response"].(map[string]any); ok {
				return extractResponsesOutput(response), extractResponsesUsage(response)
			}
		}
	}
	return "", usageInfo{}
}

// extractResponsesUsage extracts usage from a Responses API response object.
// The Responses API uses input_tokens/output_tokens (not prompt_tokens/completion_tokens).
func extractResponsesUsage(resp map[string]any) usageInfo {
	var usage usageInfo
	if usageMap, ok := resp["usage"].(map[string]any); ok {
		if v, ok := usageMap["input_tokens"].(float64); ok {
			usage.InputTokens = int(v)
		}
		if v, ok := usageMap["output_tokens"].(float64); ok {
			usage.OutputTokens = int(v)
		}
	}
	return usage
}

// extractResponsesOutput builds a chat-completions output from a Responses API response.
func extractResponsesOutput(resp map[string]any) string {
	output, ok := resp["output"].([]any)
	if !ok || len(output) == 0 {
		return ""
	}

	var textParts []string
	var refusalParts []string
	var reasoningParts []string
	var toolCalls []map[string]any

	for _, item := range output {
		itemMap, ok := item.(map[string]any)
		if !ok {
			continue
		}
		itemType, _ := itemMap["type"].(string)
		switch itemType {
		case "message":
			if content, ok := itemMap["content"].([]any); ok {
				for _, block := range content {
					blockMap, ok := block.(map[string]any)
					if !ok {
						continue
					}
					switch blockType, _ := blockMap["type"].(string); blockType {
					case "output_text":
						if text, ok := blockMap["text"].(string); ok {
							textParts = append(textParts, text)
						}
					case "refusal":
						if text, ok := blockMap["refusal"].(string); ok {
							refusalParts = append(refusalParts, text)
						}
					}
				}
			}
		case "reasoning":
			if text := extractSummaryText(itemMap); text != "" {
				reasoningParts = append(reasoningParts, text)
			}
		default:
			if tc, ok := responsesItemToToolCall(itemType, itemMap); ok {
				toolCalls = append(toolCalls, tc)
			}
		}
	}

	msg := map[string]any{"role": "assistant"}
	if len(textParts) > 0 {
		msg["content"] = strings.Join(textParts, "\n")
	}
	if len(refusalParts) > 0 {
		msg["refusal"] = strings.Join(refusalParts, "\n")
	}
	if len(reasoningParts) > 0 {
		msg["reasoning"] = strings.Join(reasoningParts, "\n")
	}
	if len(toolCalls) > 0 {
		msg["tool_calls"] = toolCalls
	}
	if len(msg) == 1 {
		return ""
	}
	return marshalMessages([]any{msg})
}

// chatToolCall builds a tool_call in the chat-completions format.
func chatToolCall(id any, name, arguments string) map[string]any {
	tc := map[string]any{
		"type": "function",
		"function": map[string]any{
			"name":      name,
			"arguments": arguments,
		},
	}
	if s, ok := id.(string); ok {
		tc["id"] = s
	}
	return tc
}

// marshalToolArgs extracts the given keys from src into a JSON object string.
func marshalToolArgs(src map[string]any, keys ...string) string {
	args := make(map[string]any, len(keys))
	for _, k := range keys {
		if v, ok := src[k].(string); ok {
			args[k] = v
		}
	}
	b, _ := json.Marshal(args)
	return string(b)
}
