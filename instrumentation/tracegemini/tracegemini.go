// Package tracegemini provides OpenTelemetry tracing for the
// Google Gemini API client using LangSmith-compatible spans.
//
// Usage:
//
//	// Configure your Gemini client to use a traced HTTP client
//	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
//		APIKey:     apiKey,
//		HTTPClient: tracegemini.Client(),
//	})
//
//	// Or use a custom tracer provider:
//	tp := sdktrace.NewTracerProvider(...)
//	client, _ := genai.NewClient(ctx, &genai.ClientConfig{
//		APIKey:     apiKey,
//		HTTPClient: tracegemini.Client(tracegemini.WithTracerProvider(tp)),
//	})
//
//	// Your Gemini API calls will now be automatically traced with LangSmith attrs
//	// resp, err := client.Models.GenerateContent(ctx, model, contents, nil)
package tracegemini

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	semconv "go.opentelemetry.io/otel/semconv/v1.40.0"
	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go/internal/genaiattr"
	"github.com/langchain-ai/langsmith-go/internal/traceutil"
)

type contextKey struct{ name string }

var ctxKeyRunName = contextKey{"run_name"}

// WithRunNameContext sets the span (run) name for the next traced request made with ctx.
// The run name in LangSmith is the OTLP span name; there is no separate field.
// Use this so one client can emit runs with different names per call, e.g. in tests:
//
//	ctx = tracegemini.WithRunNameContext(ctx, "gemini_nonstreaming")
//	client.Models.GenerateContent(ctx, ...)
func WithRunNameContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, ctxKeyRunName, name)
}

// Option configures a traced HTTP client.
type Option func(*clientOptions)

type clientOptions struct {
	tracerProvider trace.TracerProvider
	runName        string
}

// WithTracerProvider returns an Option that sets the tracer provider.
// If not provided, the global tracer provider is used.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(opts *clientOptions) {
		opts.tracerProvider = tp
	}
}

// WithRunName sets the span (run) name to the given string when non-empty. Used by integration tests to identify runs in a shared project.
func WithRunName(name string) Option {
	return func(opts *clientOptions) {
		opts.runName = name
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
	client.Transport = newRoundTripper(transport, options.tracerProvider, options.runName)
	return client
}

type roundTripper struct {
	base           http.RoundTripper
	tracerProvider trace.TracerProvider
	runName        string
}

func newRoundTripper(base http.RoundTripper, tp trace.TracerProvider, runName string) http.RoundTripper {
	return &roundTripper{
		base:           base,
		tracerProvider: tp,
		runName:        runName,
	}
}

// RoundTrip intercepts requests/responses to add OpenTelemetry tracing.
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	if !isGeminiEndpoint(req.URL.Path) {
		return rt.base.RoundTrip(req)
	}

	ctx := req.Context()
	// Prefer run name from context (per-request); fall back to client option.
	if v := ctx.Value(ctxKeyRunName); v == nil && rt.runName != "" {
		ctx = context.WithValue(ctx, ctxKeyRunName, rt.runName)
		req = req.WithContext(ctx)
	}

	var tracer trace.Tracer
	if rt.tracerProvider != nil {
		tracer = rt.tracerProvider.Tracer("google.golang.org/genai")
	} else {
		tracer = otel.Tracer("google.golang.org/genai")
	}

	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))

	parentSpan := trace.SpanFromContext(ctx)

	model, action := parseModelAction(req.URL.Path)
	streaming := action == "streamGenerateContent"

	spanName := "gemini." + action
	if v := ctx.Value(ctxKeyRunName); v != nil {
		if s, ok := v.(string); ok && s != "" {
			spanName = s
		}
	}

	opName := "generate_content"
	if streaming {
		opName = "stream_generate_content"
	}

	// Detect Vertex AI vs Gemini API from the request host, matching the
	// official OTel semconv gen_ai.provider.name values.
	providerName := semconv.GenAIProviderNameGCPGemini
	if strings.Contains(req.URL.Host, "aiplatform.googleapis.com") {
		providerName = semconv.GenAIProviderNameGCPVertexAI
	}

	// Redact API key from URL if present as a query parameter.
	spanURL := req.URL.String()
	if req.URL.Query().Has("key") {
		redacted := *req.URL
		q := redacted.Query()
		q.Set("key", "REDACTED")
		redacted.RawQuery = q.Encode()
		spanURL = redacted.String()
	}

	ctx, span := tracer.Start(ctx, spanName,
		trace.WithAttributes(
			providerName,
			semconv.GenAIOperationNameKey.String(opName),
			genaiattr.HTTPMethodKey.String(req.Method),
			genaiattr.HTTPURLKey.String(spanURL),
		),
	)
	if model != "" {
		span.SetAttributes(semconv.GenAIRequestModel(model))
	}

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

	if len(requestBody) > 0 {
		extractRequestAttributes(span, requestBody)
	}

	otel.GetTextMapPropagator().Inject(ctx, propagation.HeaderCarrier(req.Header))
	req = req.WithContext(ctx)

	resp, err := rt.base.RoundTrip(req)
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
			if readErr != nil && readErr != io.EOF {
				span.RecordError(readErr)
				span.SetStatus(codes.Error, readErr.Error())
			}
			span.End()
			return
		}

		bodyText := string(data)
		if resp.StatusCode >= 400 {
			msg := bodyText
			if len(msg) > 500 {
				msg = msg[:500] + "..."
			}
			apiErr := fmt.Errorf("HTTP %d: %s", resp.StatusCode, msg)
			span.RecordError(apiErr)
			span.SetStatus(codes.Error, apiErr.Error())
		}
		// Gemini streams include finishReason in the final chunk; absence means early termination.
		incompleteStream := resp.StatusCode < 400 && streaming && !strings.Contains(bodyText, "\"finishReason\"")
		if incompleteStream {
			endErr := readErr
			if endErr == nil || endErr == io.EOF {
				endErr = fmt.Errorf("Cancelled")
			}
			span.RecordError(endErr)
			span.SetStatus(codes.Error, endErr.Error())
		}

		if streaming {
			extractStreamingResponseAttributes(span, data, parentSpan)
		} else {
			extractResponseAttributes(span, data, parentSpan)
		}
		if resp.StatusCode < 400 && !incompleteStream {
			span.SetStatus(codes.Ok, "")
		}
		span.End()
	})
	if streaming && resp.StatusCode < 400 {
		traceutil.OnFirstSSEMatch(br, isFirstContent, func() { span.AddEvent("new_token") })
	}
	resp.Body = br

	return resp, nil
}

// isFirstContent returns true when the SSE chunk contains generated content
// (text or function call). Used for first-token-time tracking.
func isFirstContent(chunk map[string]any) bool {
	candidates, ok := chunk["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return false
	}
	candidate, ok := candidates[0].(map[string]any)
	if !ok {
		return false
	}
	content, ok := candidate["content"].(map[string]any)
	if !ok {
		return false
	}
	parts, ok := content["parts"].([]any)
	if !ok || len(parts) == 0 {
		return false
	}
	for _, p := range parts {
		part, ok := p.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := part["text"].(string); ok && text != "" {
			return true
		}
		if _, ok := part["functionCall"]; ok {
			return true
		}
	}
	return false
}

// isGeminiEndpoint returns true if the path matches a Gemini generateContent
// or streamGenerateContent endpoint.
func isGeminiEndpoint(path string) bool {
	if !strings.Contains(path, "/models/") {
		return false
	}
	_, action := parseModelAction(path)
	return action == "generateContent" || action == "streamGenerateContent"
}

// parseModelAction extracts the model name and action from a Gemini API path.
// e.g. /v1beta/models/gemini-2.0-flash:generateContent → ("gemini-2.0-flash", "generateContent")
func parseModelAction(path string) (model, action string) {
	idx := strings.Index(path, "/models/")
	if idx < 0 {
		return "", ""
	}
	rest := path[idx+len("/models/"):]
	if m, a, ok := strings.Cut(rest, ":"); ok {
		return m, a
	}
	return rest, ""
}

// extractRequestAttributes extracts attributes from the Gemini request body
// and sets them on the span. Handles contents, systemInstruction, and generationConfig.
func extractRequestAttributes(span trace.Span, body []byte) {
	var req map[string]any
	if err := json.Unmarshal(body, &req); err != nil {
		return
	}

	if cfg, ok := req["generationConfig"].(map[string]any); ok {
		if temp, ok := cfg["temperature"].(float64); ok {
			span.SetAttributes(semconv.GenAIRequestTemperature(temp))
		}
		if maxTokens, ok := cfg["maxOutputTokens"].(float64); ok {
			span.SetAttributes(semconv.GenAIRequestMaxTokens(int(maxTokens)))
		}
		if topP, ok := cfg["topP"].(float64); ok {
			span.SetAttributes(semconv.GenAIRequestTopP(topP))
		}
	}

	var messages []any

	// systemInstruction → system message (prepended)
	if sysInstr, ok := req["systemInstruction"].(map[string]any); ok {
		if text := partsText(sysInstr); text != "" {
			messages = append(messages, map[string]any{"role": "system", "content": text})
		}
	}

	if contents, ok := req["contents"].([]any); ok {
		for _, c := range contents {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			messages = append(messages, contentToMessage(cm))
		}
	}

	if len(messages) > 0 {
		if out, err := json.Marshal(map[string]any{"messages": messages}); err == nil {
			span.SetAttributes(genaiattr.PromptKey.String(string(out)))
		}
	}
}

// contentToMessage converts a Gemini Content object to a LangSmith-style message,
// matching the format produced by langsmith-python's _process_gemini_inputs:
//   - text-only parts → content is a simple "\n"-joined string
//   - mixed parts (text + function calls/responses) → content is a typed parts list
func contentToMessage(content map[string]any) map[string]any {
	role := mapGeminiRole(content["role"])
	msg := map[string]any{"role": role}

	parts, _ := content["parts"].([]any)
	var textParts []string
	var structured []any
	var hasNonText bool

	for _, p := range parts {
		part, ok := p.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := part["text"].(string); ok && text != "" {
			textParts = append(textParts, text)
			structured = append(structured, map[string]any{"type": "text", "text": text})
		}
		if fc, ok := part["functionCall"].(map[string]any); ok {
			hasNonText = true
			formatted := map[string]any{"name": fc["name"]}
			if args, ok := fc["args"]; ok {
				formatted["arguments"] = args
			}
			structured = append(structured, map[string]any{
				"type":          "function_call",
				"function_call": formatted,
			})
		}
		if fr, ok := part["functionResponse"].(map[string]any); ok {
			hasNonText = true
			structured = append(structured, map[string]any{
				"type": "function_response",
				"function_response": map[string]any{
					"name":     fr["name"],
					"response": fr["response"],
				},
			})
		}
	}

	if !hasNonText && len(textParts) > 0 {
		msg["content"] = strings.Join(textParts, "\n")
	} else if len(structured) > 0 {
		msg["content"] = structured
	}
	return msg
}

// partsText concatenates text from all parts in a Content-like object.
func partsText(content map[string]any) string {
	parts, ok := content["parts"].([]any)
	if !ok {
		return ""
	}
	var texts []string
	for _, p := range parts {
		part, ok := p.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := part["text"].(string); ok {
			texts = append(texts, text)
		}
	}
	return strings.Join(texts, "\n")
}

func mapGeminiRole(role any) string {
	s, _ := role.(string)
	switch s {
	case "model":
		return "assistant"
	case "user":
		return "user"
	case "function":
		return "tool"
	case "":
		return "user"
	default:
		return s
	}
}

// extractResponseAttributes extracts attributes from a non-streaming Gemini response body.
func extractResponseAttributes(span trace.Span, body []byte, parentSpan trace.Span) {
	var resp map[string]any
	if err := json.Unmarshal(body, &resp); err != nil {
		return
	}
	processResponse(span, resp, parentSpan)
}

// extractStreamingResponseAttributes merges SSE chunks into a single
// synthetic response and processes it through the same path as non-streaming.
func extractStreamingResponseAttributes(span trace.Span, data []byte, parentSpan trace.Span) {
	chunks, err := traceutil.ParseSSEChunks(bytes.NewReader(data))
	if err != nil || len(chunks) == 0 {
		return
	}
	processResponse(span, mergeStreamingChunks(chunks), parentSpan)
}

// mergeStreamingChunks folds SSE chunks (each a GenerateContentResponse)
// into one synthetic response by accumulating content parts and keeping the
// last-seen metadata/usage.
func mergeStreamingChunks(chunks []map[string]any) map[string]any {
	resp := map[string]any{}
	var allParts []any
	var finishReason string

	for _, chunk := range chunks {
		if mv, ok := chunk["modelVersion"].(string); ok {
			resp["modelVersion"] = mv
		}
		if rid, ok := chunk["responseId"].(string); ok {
			resp["responseId"] = rid
		}
		if u, ok := chunk["usageMetadata"].(map[string]any); ok {
			resp["usageMetadata"] = u
		}
		candidates, ok := chunk["candidates"].([]any)
		if !ok || len(candidates) == 0 {
			continue
		}
		candidate, ok := candidates[0].(map[string]any)
		if !ok {
			continue
		}
		if r, ok := candidate["finishReason"].(string); ok && r != "" {
			finishReason = r
		}
		content, ok := candidate["content"].(map[string]any)
		if !ok {
			continue
		}
		if parts, ok := content["parts"].([]any); ok {
			allParts = append(allParts, parts...)
		}
	}

	candidate := map[string]any{
		"content": map[string]any{"parts": allParts},
	}
	if finishReason != "" {
		candidate["finishReason"] = finishReason
	}
	resp["candidates"] = []any{candidate}
	return resp
}

// processResponse sets span attributes from a (possibly merged) Gemini response.
func processResponse(span trace.Span, resp map[string]any, parentSpan trace.Span) {
	if mv, ok := resp["modelVersion"].(string); ok {
		span.SetAttributes(semconv.GenAIResponseModel(mv))
	}
	if rid, ok := resp["responseId"].(string); ok {
		span.SetAttributes(semconv.GenAIResponseID(rid))
	}

	if completion, finishReason := buildCompletion(resp); completion != "" {
		span.SetAttributes(genaiattr.CompletionKey.String(completion))
		if finishReason != "" {
			span.SetAttributes(genaiattr.StopReasonKey.String(finishReason))
		}
	}

	if usage, ok := resp["usageMetadata"].(map[string]any); ok {
		setUsageAttributes(span, usage, parentSpan)
	}
}

// buildCompletion builds a {"messages":[...]} JSON string from candidates,
// matching langsmith-python's _process_generate_content_response:
//   - tool calls use OpenAI-compatible format (id, type, index, function)
//   - finish_reason is included in the message
//   - content is null (not omitted) when only tool calls are present
func buildCompletion(resp map[string]any) (completion, finishReason string) {
	candidates, ok := resp["candidates"].([]any)
	if !ok || len(candidates) == 0 {
		return "", ""
	}
	candidate, ok := candidates[0].(map[string]any)
	if !ok {
		return "", ""
	}
	finishReason, _ = candidate["finishReason"].(string)

	content, ok := candidate["content"].(map[string]any)
	if !ok {
		return "", finishReason
	}
	parts, ok := content["parts"].([]any)
	if !ok || len(parts) == 0 {
		return "", finishReason
	}

	var texts []string
	var toolCalls []any
	for _, p := range parts {
		part, ok := p.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := part["text"].(string); ok && text != "" {
			texts = append(texts, text)
		}
		if fc, ok := part["functionCall"].(map[string]any); ok {
			name, _ := fc["name"].(string)
			argsJSON := "{}"
			if args, ok := fc["args"]; ok {
				if ab, err := json.Marshal(args); err == nil {
					argsJSON = string(ab)
				}
			}
			id, _ := fc["id"].(string)
			if id == "" {
				id = fmt.Sprintf("call_%d", len(toolCalls))
			}
			toolCalls = append(toolCalls, map[string]any{
				"id":    id,
				"type":  "function",
				"index": len(toolCalls),
				"function": map[string]any{
					"name":      name,
					"arguments": argsJSON,
				},
			})
		}
	}
	if len(texts) == 0 && len(toolCalls) == 0 {
		return "", finishReason
	}

	msg := map[string]any{"role": "assistant"}
	if finishReason != "" {
		msg["finish_reason"] = finishReason
	}
	if len(toolCalls) > 0 {
		if len(texts) > 0 {
			msg["content"] = strings.Join(texts, "")
		} else {
			msg["content"] = nil
		}
		msg["tool_calls"] = toolCalls
	} else {
		msg["content"] = strings.Join(texts, "")
	}

	out, err := json.Marshal(map[string]any{"messages": []any{msg}})
	if err != nil {
		return "", finishReason
	}
	return string(out), finishReason
}

// setUsageAttributes sets usage-related attributes on the span.
// Gemini uses promptTokenCount/candidatesTokenCount instead of input_tokens/output_tokens.
func setUsageAttributes(span trace.Span, usage map[string]any, parentSpan trace.Span) {
	var inputTokens, outputTokens int64

	if v, ok := usage["promptTokenCount"].(float64); ok {
		inputTokens = int64(v)
	}
	if v, ok := usage["candidatesTokenCount"].(float64); ok {
		outputTokens = int64(v)
	}

	if inputTokens > 0 {
		span.SetAttributes(semconv.GenAIUsageInputTokens(int(inputTokens)))
		if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
			parentSpan.SetAttributes(semconv.GenAIUsageInputTokens(int(inputTokens)))
		}
	}
	if outputTokens > 0 {
		span.SetAttributes(semconv.GenAIUsageOutputTokens(int(outputTokens)))
		if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
			parentSpan.SetAttributes(semconv.GenAIUsageOutputTokens(int(outputTokens)))
		}
	}

	if v, ok := usage["cachedContentTokenCount"].(float64); ok && v > 0 {
		span.SetAttributes(
			genaiattr.CacheReadInputTokensKey.Int64(int64(v)),
			semconv.GenAIUsageCacheReadInputTokens(int(v)),
		)
	}
	if v, ok := usage["thoughtsTokenCount"].(float64); ok && v > 0 {
		span.SetAttributes(genaiattr.UsageReasoningTokensKey.Int64(int64(v)))
	}
}
