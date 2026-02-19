// Package traceanthropic provides OpenTelemetry tracing for the
// Anthropic API client using LangSmith-compatible spans.
//
// Usage:
//
//	// Configure your Anthropic client to use a traced HTTP client
//	client := anthropic.NewClient(
//		anthropic.WithAPIKey(apiKey),
//		anthropic.WithHTTPClient(traceanthropic.Client()),
//	)
//
//	// Or use a custom tracer provider:
//	tp := sdktrace.NewTracerProvider(...)
//	client := anthropic.NewClient(
//		anthropic.WithAPIKey(apiKey),
//		anthropic.WithHTTPClient(traceanthropic.Client(traceanthropic.WithTracerProvider(tp))),
//	)
//
//	// Your Anthropic API calls will now be automatically traced with LangSmith attrs
//	// resp, err := client.Messages.Create(ctx, ...)
package traceanthropic

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go/internal/traceutil"
)

// Option configures a traced HTTP client.
type Option func(*clientOptions)

type clientOptions struct {
	tracerProvider trace.TracerProvider
}

// WithTracerProvider returns an Option that sets the tracer provider.
// If not provided, the global tracer provider is used.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(opts *clientOptions) {
		opts.tracerProvider = tp
	}
}

// Client returns a new http.Client configured with tracing middleware.
// Equivalent to WrapClient(nil, opts...), which wraps the default transport.
func Client(opts ...Option) *http.Client {
	return WrapClient(nil, opts...)
}

// WrapClient wraps an existing http.Client with tracing middleware.
// If client is nil, a new client with the default transport is created.
func WrapClient(client *http.Client, opts ...Option) *http.Client {
	options := &clientOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if client == nil {
		client = &http.Client{}
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = newRoundTripper(transport, options.tracerProvider)
	return client
}

type roundTripper struct {
	base           http.RoundTripper
	tracerProvider trace.TracerProvider
}

func newRoundTripper(base http.RoundTripper, tp trace.TracerProvider) http.RoundTripper {
	return &roundTripper{
		base:           base,
		tracerProvider: tp,
	}
}

// RoundTrip intercepts requests/responses to add OpenTelemetry tracing.
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only trace Anthropic API requests
	if !strings.Contains(req.URL.Host, "api.anthropic.com") {
		return rt.base.RoundTrip(req)
	}

	ctx := req.Context()
	var tracer trace.Tracer
	if rt.tracerProvider != nil {
		tracer = rt.tracerProvider.Tracer("github.com/anthropics/anthropic-sdk-go")
	} else {
		tracer = otel.Tracer("github.com/anthropics/anthropic-sdk-go")
	}

	// Extract span context from request headers
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))

	// Capture parent span before creating child span (for token propagation)
	parentSpan := trace.SpanFromContext(ctx)

	// Determine span name based on endpoint
	spanName := getSpanName(req.URL.Path)

	// Start span (child span)
	ctx, span := tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("gen_ai.system", "anthropic"),
			attribute.String("gen_ai.operation.name", getOperationName(req.URL.Path)),
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.URL.String()),
		),
	)

	// Read request body if present
	var requestBody []byte
	if req.Body != nil {
		var err error
		requestBody, err = io.ReadAll(req.Body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, fmt.Sprintf("failed to read request body: %v", err))
			span.End()
			return rt.base.RoundTrip(req)
		}
		req.Body = io.NopCloser(bytes.NewBuffer(requestBody))
	}

	// Extract request attributes from body
	if len(requestBody) > 0 {
		extractRequestAttributes(span, requestBody)
	}

	// Inject span context into request headers
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))

	// Make the actual request
	resp, err := rt.base.RoundTrip(req)
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
		if streaming {
			extractStreamingResponseAttributes(span, data, parentSpan)
		} else {
			extractResponseAttributes(span, data, parentSpan)
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
	if strings.Contains(path, "/v1/messages") {
		return "anthropic.messages"
	}
	return "anthropic.request"
}

// getOperationName returns the operation name for Gen AI semantic conventions.
func getOperationName(path string) string {
	if strings.Contains(path, "/v1/messages") {
		return "chat"
	}
	return "request"
}

// extractRequestAttributes extracts attributes from Anthropic request body.
func extractRequestAttributes(span trace.Span, body []byte) {
	var req map[string]interface{}
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}

	// Extract model
	if model, ok := req["model"].(string); ok {
		span.SetAttributes(attribute.String("gen_ai.request.model", model))
	}

	// Extract max_tokens
	if maxTokens, ok := req["max_tokens"].(float64); ok {
		span.SetAttributes(attribute.Int64("gen_ai.request.max_tokens", int64(maxTokens)))
	}

	// Extract temperature
	if temp, ok := req["temperature"].(float64); ok {
		span.SetAttributes(attribute.Float64("gen_ai.request.temperature", temp))
	}

	// Extract system message if present
	var messages []interface{}
	if sys, ok := req["system"]; ok {
		messages = append(messages, map[string]interface{}{
			"role":    "system",
			"content": sys,
		})
	}

	// Extract user messages
	if msgs, ok := req["messages"].([]interface{}); ok {
		messages = append(messages, msgs...)
	}

	// Build prompt from messages
	if len(messages) > 0 {
		prompt := buildPromptFromMessages(messages)
		if prompt != "" {
			span.SetAttributes(attribute.String("gen_ai.prompt", prompt))
		}
	}
}

// buildPromptFromMessages builds a prompt string from Anthropic messages.
func buildPromptFromMessages(messages []interface{}) string {
	var parts []string
	for _, msg := range messages {
		if msgMap, ok := msg.(map[string]interface{}); ok {
			if role, _ := msgMap["role"].(string); role == "user" || role == "system" {
				if content, ok := msgMap["content"].(string); ok {
					parts = append(parts, content)
				} else if contentArray, ok := msgMap["content"].([]interface{}); ok {
					// Handle content blocks
					for _, block := range contentArray {
						if blockMap, ok := block.(map[string]interface{}); ok {
							if blockType, _ := blockMap["type"].(string); blockType == "text" {
								if text, _ := blockMap["text"].(string); text != "" {
									parts = append(parts, text)
								}
							}
						}
					}
				}
			}
		}
	}
	return strings.Join(parts, "\n")
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

// extractStreamingResponseAttributes parses an SSE response body and sets
// the same span attributes as the non-streaming path: gen_ai.completion,
// usage tokens, and cache breakdown.
//
// Anthropic SSE events of interest:
//   - message_start  — contains message.usage (input_tokens, cache tokens)
//   - content_block_delta — text_delta carries incremental text
//   - message_delta  — contains usage.output_tokens
func extractStreamingResponseAttributes(span trace.Span, data []byte, parentSpan trace.Span) {
	chunks, err := traceutil.ParseSSEChunks(bytes.NewReader(data))
	if err != nil || len(chunks) == 0 {
		return
	}

	usage := make(map[string]interface{})
	var content strings.Builder

	for _, chunk := range chunks {
		eventType, _ := chunk["type"].(string)

		switch eventType {
		case "message_start":
			// message_start contains message.usage with input/cache tokens
			if message, ok := chunk["message"].(map[string]any); ok {
				if curUsage, ok := message["usage"].(map[string]any); ok {
					for k, v := range curUsage {
						usage[k] = v
					}
				}
				// Extract model from message_start
				if model, ok := message["model"].(string); ok {
					span.SetAttributes(attribute.String("gen_ai.response.model", model))
				}
			}

		case "content_block_delta":
			// Aggregate text from text_delta events
			if delta, ok := chunk["delta"].(map[string]any); ok {
				if deltaType, _ := delta["type"].(string); deltaType == "text_delta" {
					if text, ok := delta["text"].(string); ok {
						content.WriteString(text)
					}
				}
			}

		case "message_delta":
			// message_delta contains usage.output_tokens and stop_reason
			if curUsage, ok := chunk["usage"].(map[string]any); ok {
				for k, v := range curUsage {
					usage[k] = v
				}
			}
			if delta, ok := chunk["delta"].(map[string]any); ok {
				if stopReason, ok := delta["stop_reason"].(string); ok && stopReason != "" {
					span.SetAttributes(attribute.String("langsmith.metadata.stop_reason", stopReason))
				}
			}
		}
	}

	if text := content.String(); text != "" {
		span.SetAttributes(attribute.String("gen_ai.completion", text))
	}

	if len(usage) > 0 {
		setUsageAttributes(span, usage, parentSpan)
	}
}

// extractResponseAttributes extracts attributes from Anthropic response body.
func extractResponseAttributes(span trace.Span, body []byte, parentSpan trace.Span) {
	var resp map[string]interface{}
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}

	// Extract model
	if model, ok := resp["model"].(string); ok {
		span.SetAttributes(attribute.String("gen_ai.response.model", model))
	}

	// Extract usage
	if usage, ok := resp["usage"].(map[string]interface{}); ok {
		setUsageAttributes(span, usage, parentSpan)
	}

	// Extract stop_reason
	if stopReason, ok := resp["stop_reason"].(string); ok && stopReason != "" {
		span.SetAttributes(attribute.String("langsmith.metadata.stop_reason", stopReason))
	}

	// Extract completion from content
	if content, ok := resp["content"].([]interface{}); ok && len(content) > 0 {
		completion := extractCompletionFromContent(content)
		if completion != "" {
			span.SetAttributes(attribute.String("gen_ai.completion", completion))
		}
	}
}

// extractCompletionFromContent extracts text completion from Anthropic content blocks.
func extractCompletionFromContent(content []interface{}) string {
	var parts []string
	for _, block := range content {
		if blockMap, ok := block.(map[string]interface{}); ok {
			if blockType, _ := blockMap["type"].(string); blockType == "text" {
				if text, _ := blockMap["text"].(string); text != "" {
					parts = append(parts, text)
				}
			}
		}
	}
	return strings.Join(parts, "\n")
}

// setUsageAttributes sets usage-related attributes on the span.
func setUsageAttributes(span trace.Span, usage map[string]interface{}, parentSpan trace.Span) {
	var inputTokens, outputTokens, totalTokens int64

	if v, ok := usage["input_tokens"].(float64); ok {
		inputTokens = int64(v)
	}
	if v, ok := usage["output_tokens"].(float64); ok {
		outputTokens = int64(v)
	}
	if v, ok := usage["total_tokens"].(float64); ok {
		totalTokens = int64(v)
	}

	// Handle cache tokens if present
	var cacheCreate, cacheRead int64
	if v, ok := usage["cache_creation_input_tokens"].(float64); ok {
		cacheCreate = int64(v)
	}
	if v, ok := usage["cache_read_input_tokens"].(float64); ok {
		cacheRead = int64(v)
	}

	// Total prompt tokens includes cache
	totalPrompt := inputTokens + cacheCreate + cacheRead

	if totalPrompt > 0 {
		span.SetAttributes(attribute.Int64("gen_ai.usage.input_tokens", totalPrompt))
		// Propagate usage to parent span if it exists and is valid
		// This ensures token counts appear in Thread list view (which aggregates from root spans)
		if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
			parentSpan.SetAttributes(attribute.Int64("gen_ai.usage.input_tokens", totalPrompt))
		}
	}
	if outputTokens > 0 {
		span.SetAttributes(attribute.Int64("gen_ai.usage.output_tokens", outputTokens))
		// Propagate usage to parent span if it exists and is valid
		if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
			parentSpan.SetAttributes(attribute.Int64("gen_ai.usage.output_tokens", outputTokens))
		}
	}
	if totalTokens > 0 {
		span.SetAttributes(attribute.Int64("gen_ai.usage.total_tokens", totalTokens))
	}

	// Cache breakdown in metadata
	if cacheCreate > 0 {
		span.SetAttributes(attribute.Int64("langsmith.metadata.usage_metadata.input_token_details.cache_creation", cacheCreate))
	}
	if cacheRead > 0 {
		span.SetAttributes(attribute.Int64("langsmith.metadata.usage_metadata.input_token_details.cache_read", cacheRead))
	}

	// Reasoning tokens if present (for Claude Sonnet 4.5+)
	if v, ok := usage["reasoning_tokens"].(float64); ok {
		span.SetAttributes(attribute.Int64("gen_ai.usage.details.reasoning_tokens", int64(v)))
	}
}
