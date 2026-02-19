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

	// Extract request attributes and streaming flag
	var streaming bool
	if len(requestBody) > 0 {
		streaming = extractRequestAttributes(span, requestBody)
	}

	// Inject span context into request headers and update request context
	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	req = req.WithContext(ctx)

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

// extractRequestAttributes extracts attributes from the Anthropic request body
// and returns whether the request is streaming.
func extractRequestAttributes(span trace.Span, body []byte) (streaming bool) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return false
	}

	// Model
	if model, ok := req["model"].(string); ok {
		span.SetAttributes(attribute.String("gen_ai.request.model", model))
	}

	// Max tokens
	if maxTokens, ok := req["max_tokens"].(float64); ok {
		span.SetAttributes(attribute.Int64("gen_ai.request.max_tokens", int64(maxTokens)))
	}

	// Temperature
	if temp, ok := req["temperature"].(float64); ok {
		span.SetAttributes(attribute.Float64("gen_ai.request.temperature", temp))
	}

	// Streaming flag
	streaming, _ = req["stream"].(bool)

	// Build input messages — system prepended, all roles preserved
	var messages []any
	if sys, ok := req["system"]; ok {
		messages = append(messages, map[string]any{
			"role":    "system",
			"content": sys,
		})
	}
	if msgs, ok := req["messages"].([]any); ok {
		messages = append(messages, msgs...)
	}

	if len(messages) > 0 {
		if out, err := json.Marshal(map[string]any{"messages": messages}); err == nil {
			span.SetAttributes(attribute.String("gen_ai.prompt", string(out)))
		}
	}

	return streaming
}

// extractStreamingResponseAttributes parses an SSE response body and sets
// the same span attributes as the non-streaming path: gen_ai.completion,
// usage tokens, and cache breakdown.
//
// Anthropic SSE events of interest:
//   - message_start       — contains message.usage (input_tokens, cache tokens)
//   - content_block_start — initialises a content block (text or tool_use)
//   - content_block_delta — text_delta or input_json_delta carries incremental data
//   - message_delta       — contains usage.output_tokens and stop_reason
func extractStreamingResponseAttributes(span trace.Span, data []byte, parentSpan trace.Span) {
	chunks, err := traceutil.ParseSSEChunks(bytes.NewReader(data))
	if err != nil || len(chunks) == 0 {
		return
	}

	usage := make(map[string]interface{})

	// Track content blocks by index for proper multi-block reconstruction
	type contentBlock struct {
		blockType string // "text" or "tool_use"
		id        string // tool_use id
		name      string // tool_use function name
		buf       strings.Builder
	}
	var blocks []*contentBlock

	for _, chunk := range chunks {
		eventType, _ := chunk["type"].(string)

		switch eventType {
		case "message_start":
			if message, ok := chunk["message"].(map[string]any); ok {
				if curUsage, ok := message["usage"].(map[string]any); ok {
					for k, v := range curUsage {
						usage[k] = v
					}
				}
				if model, ok := message["model"].(string); ok {
					span.SetAttributes(attribute.String("gen_ai.response.model", model))
				}
			}

		case "content_block_start":
			idxF, ok := chunk["index"].(float64)
			if !ok {
				continue
			}
			idx := int(idxF)
			for len(blocks) <= idx {
				blocks = append(blocks, nil)
			}
			block := &contentBlock{}
			if cb, ok := chunk["content_block"].(map[string]any); ok {
				block.blockType, _ = cb["type"].(string)
				block.id, _ = cb["id"].(string)
				block.name, _ = cb["name"].(string)
			}
			blocks[idx] = block

		case "content_block_delta":
			idxF, ok := chunk["index"].(float64)
			if !ok {
				continue
			}
			idx := int(idxF)
			for len(blocks) <= idx {
				blocks = append(blocks, nil)
			}
			if blocks[idx] == nil {
				blocks[idx] = &contentBlock{}
			}
			if delta, ok := chunk["delta"].(map[string]any); ok {
				switch deltaType, _ := delta["type"].(string); deltaType {
				case "text_delta":
					if text, ok := delta["text"].(string); ok {
						blocks[idx].buf.WriteString(text)
						if blocks[idx].blockType == "" {
							blocks[idx].blockType = "text"
						}
					}
				case "input_json_delta":
					if partialJSON, ok := delta["partial_json"].(string); ok {
						blocks[idx].buf.WriteString(partialJSON)
						if blocks[idx].blockType == "" {
							blocks[idx].blockType = "tool_use"
						}
					}
				}
			}

		case "message_delta":
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

	// Reconstruct content blocks into an assistant message
	var contentBlocks []map[string]any
	for _, b := range blocks {
		if b == nil {
			continue
		}
		switch b.blockType {
		case "text":
			contentBlocks = append(contentBlocks, map[string]any{
				"type": "text",
				"text": b.buf.String(),
			})
		case "tool_use":
			block := map[string]any{
				"type": "tool_use",
				"id":   b.id,
				"name": b.name,
			}
			var input any
			if err := json.Unmarshal([]byte(b.buf.String()), &input); err == nil {
				block["input"] = input
			} else {
				block["input"] = b.buf.String()
			}
			contentBlocks = append(contentBlocks, block)
		}
	}

	if len(contentBlocks) > 0 {
		msg := map[string]any{"messages": []any{
			map[string]any{"role": "assistant", "content": contentBlocks},
		}}
		if out, err := json.Marshal(msg); err == nil {
			span.SetAttributes(attribute.String("gen_ai.completion", string(out)))
		}
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

	// Extract output — wrap full content array in an assistant message
	if content, ok := resp["content"].([]any); ok && len(content) > 0 {
		msg := map[string]any{"messages": []any{
			map[string]any{"role": "assistant", "content": content},
		}}
		if out, err := json.Marshal(msg); err == nil {
			span.SetAttributes(attribute.String("gen_ai.completion", string(out)))
		}
	}
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
