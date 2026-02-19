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
	// Only trace OpenAI API requests
	if !strings.Contains(req.URL.Host, "api.openai.com") {
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

	// Determine span name based on endpoint
	spanName := getSpanName(req.URL.Path)

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
			attribute.String("session_id", value), // Standard format per docs
			attribute.String("langsmith.metadata.session_id", value), // LangSmith metadata format
			attribute.String("session.id", value), // Compatibility format
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

	// Extract prompt from request body
	if len(requestBody) > 0 {
		prompt := extractPromptFromRequest(requestBody)
		if prompt != "" {
			span.SetAttributes(attribute.String("gen_ai.prompt", prompt))
		}

		// Extract model if present
		model := extractModelFromRequest(requestBody)
		if model != "" {
			span.SetAttributes(attribute.String("gen_ai.request.model", model))
		}
	}

	// Inject span context into request headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Make the actual request via next middleware/transport
	resp, err := next(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		span.End()
		return resp, err
	}

	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
	}

	streaming := isStreaming(requestBody)

	resp.Body = traceutil.NewBufferedReader(resp.Body, func(r io.Reader) {
		data, err := io.ReadAll(r)
		if err != nil || len(data) == 0 {
			span.End()
			return
		}

		var completion string
		var usage usageInfo
		if streaming {
			completion, usage = extractStreamingCompletion(data)
		} else {
			completion, usage = extractCompletionFromResponse(data)
		}

		if completion != "" {
			span.SetAttributes(attribute.String("gen_ai.completion", completion))
		}
		if usage.InputTokens > 0 {
			inputTokens := int64(usage.InputTokens)
			outputTokens := int64(usage.OutputTokens)
			span.SetAttributes(
				attribute.Int64("gen_ai.usage.input_tokens", inputTokens),
				attribute.Int64("gen_ai.usage.output_tokens", outputTokens),
			)
			if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
				parentSpan.SetAttributes(
					attribute.Int64("gen_ai.usage.input_tokens", inputTokens),
					attribute.Int64("gen_ai.usage.output_tokens", outputTokens),
				)
			}
		}
		if resp.StatusCode < 400 {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	})

	return resp, nil
}

// getSpanName returns an appropriate span name based on the API endpoint.
func getSpanName(path string) string {
	if strings.Contains(path, "/chat/completions") {
		return "openai.chat.completion"
	}
	if strings.Contains(path, "/completions") {
		return "openai.completion"
	}
	if strings.Contains(path, "/embeddings") {
		return "openai.embedding"
	}
	return "openai.request"
}

// getOperationName returns the operation name for Gen AI semantic conventions.
func getOperationName(path string) string {
	if strings.Contains(path, "/chat/completions") {
		return "chat"
	}
	if strings.Contains(path, "/completions") {
		return "completion"
	}
	if strings.Contains(path, "/embeddings") {
		return "embedding"
	}
	return "request"
}

// extractPromptFromRequest extracts the prompt text from OpenAI request body.
func extractPromptFromRequest(body []byte) string {
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		return ""
	}

	// For chat completions, extract messages
	if messages, ok := req["messages"].([]interface{}); ok {
		var promptParts []string
		for _, msg := range messages {
			if msgMap, ok := msg.(map[string]interface{}); ok {
				if role, _ := msgMap["role"].(string); role == "user" || role == "system" {
					if content, _ := msgMap["content"].(string); content != "" {
						promptParts = append(promptParts, content)
					}
				}
			}
		}
		if len(promptParts) > 0 {
			return strings.Join(promptParts, "\n")
		}
	}

	// For completions, extract prompt
	if prompt, ok := req["prompt"].(string); ok {
		return prompt
	}

	return ""
}

// extractModelFromRequest extracts the model name from OpenAI request body.
func extractModelFromRequest(body []byte) string {
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		return ""
	}

	if model, ok := req["model"].(string); ok {
		return model
	}

	return ""
}

// isStreaming checks the request JSON for "stream":true.
func isStreaming(requestBody []byte) bool {
	var req map[string]any
	if err := json.Unmarshal(requestBody, &req); err != nil {
		return false
	}
	v, ok := req["stream"].(bool)
	return ok && v
}

// extractStreamingCompletion parses an SSE response body, aggregates
// delta.content across chunks, and extracts usage from the chunk that
// contains it (last chunk when stream_options.include_usage is set).
func extractStreamingCompletion(data []byte) (string, usageInfo) {
	chunks, err := traceutil.ParseSSEChunks(bytes.NewReader(data))
	if err != nil || len(chunks) == 0 {
		return "", usageInfo{}
	}

	var content strings.Builder
	var usage usageInfo

	for _, chunk := range chunks {
		// Aggregate choices[0].delta.content
		if choices, ok := chunk["choices"].([]any); ok && len(choices) > 0 {
			if choice, ok := choices[0].(map[string]any); ok {
				if delta, ok := choice["delta"].(map[string]any); ok {
					if text, ok := delta["content"].(string); ok {
						content.WriteString(text)
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

	return content.String(), usage
}

// usageInfo holds token usage information.
type usageInfo struct {
	InputTokens  int
	OutputTokens int
}

// extractCompletionFromResponse extracts completion text and usage from OpenAI response.
func extractCompletionFromResponse(body []byte) (string, usageInfo) {
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return "", usageInfo{}
	}

	var completion string
	var usage usageInfo

	// Extract from choices array (chat completions)
	if choices, ok := resp["choices"].([]interface{}); ok && len(choices) > 0 {
		if choice, ok := choices[0].(map[string]interface{}); ok {
			if message, ok := choice["message"].(map[string]interface{}); ok {
				if content, ok := message["content"].(string); ok {
					completion = content
				}
			}
			// For completions endpoint
			if text, ok := choice["text"].(string); ok {
				completion = text
			}
		}
	}

	// Extract usage
	if usageMap, ok := resp["usage"].(map[string]interface{}); ok {
		if promptTokens, ok := usageMap["prompt_tokens"].(float64); ok {
			usage.InputTokens = int(promptTokens)
		}
		if completionTokens, ok := usageMap["completion_tokens"].(float64); ok {
			usage.OutputTokens = int(completionTokens)
		}
	}

	return completion, usage
}
