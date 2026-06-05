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
	"go.opentelemetry.io/otel/attribute"
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
		incompleteStream := resp.StatusCode < 400 && streaming && readErr != io.EOF
		if incompleteStream {
			endErr := readErr
			if endErr == nil {
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
		rawParts, _ := sysInstr["parts"].([]any)
		if p := parseParts(rawParts); len(p.texts) > 0 {
			messages = append(messages, map[string]any{
				"role":    "system",
				"content": strings.Join(p.texts, "\n"),
			})
		}
	}

	if contents, ok := req["contents"].([]any); ok {
		for _, c := range contents {
			cm, ok := c.(map[string]any)
			if !ok {
				continue
			}
			messages = append(messages, contentToMessages(cm)...)
		}
	}

	if len(messages) > 0 {
		if out, err := json.Marshal(map[string]any{"messages": messages}); err == nil {
			span.SetAttributes(genaiattr.PromptKey.String(string(out)))
		}
	}
}

// parsedParts holds the result of scanning Gemini content parts.
type parsedParts struct {
	texts         []string
	toolCalls     []any
	toolResponses []any
}

// parseParts scans Gemini content parts and returns text, tool calls
// (OpenAI-compatible format), and tool response messages.
func parseParts(parts []any) parsedParts {
	var p parsedParts
	for _, raw := range parts {
		part, ok := raw.(map[string]any)
		if !ok {
			continue
		}
		if text, ok := part["text"].(string); ok && text != "" {
			p.texts = append(p.texts, text)
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
				id = fmt.Sprintf("call_%d", len(p.toolCalls))
			}
			p.toolCalls = append(p.toolCalls, map[string]any{
				"id":    id,
				"type":  "function",
				"index": len(p.toolCalls),
				"function": map[string]any{
					"name":      name,
					"arguments": argsJSON,
				},
			})
		}
		if fr, ok := part["functionResponse"].(map[string]any); ok {
			name, _ := fr["name"].(string)
			respContent := "{}"
			if r, ok := fr["response"]; ok {
				if rb, err := json.Marshal(r); err == nil {
					respContent = string(rb)
				}
			}
			p.toolResponses = append(p.toolResponses, map[string]any{
				"role":    "tool",
				"name":    name,
				"content": respContent,
			})
		}
	}
	return p
}

// contentToMessages converts a Gemini Content object to one or more
// OpenAI-compatible messages. A single Gemini content turn may expand to
// multiple messages:
//   - text-only → one message with string content
//   - function calls → one assistant message with tool_calls array
//   - function responses → one "tool" message per response
func contentToMessages(content map[string]any) []any {
	role := mapGeminiRole(content["role"])
	rawParts, _ := content["parts"].([]any)
	p := parseParts(rawParts)

	if len(p.toolCalls) > 0 {
		msg := map[string]any{"role": "assistant"}
		if len(p.texts) > 0 {
			msg["content"] = strings.Join(p.texts, "\n")
		} else {
			msg["content"] = nil
		}
		msg["tool_calls"] = p.toolCalls
		return append([]any{msg}, p.toolResponses...)
	}

	if len(p.toolResponses) > 0 {
		return p.toolResponses
	}

	msg := map[string]any{"role": role}
	if len(p.texts) > 0 {
		msg["content"] = strings.Join(p.texts, "\n")
	}
	return []any{msg}
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
// using OpenAI-compatible format for tool calls.
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

	p := parseParts(parts)
	if len(p.texts) == 0 && len(p.toolCalls) == 0 {
		return "", finishReason
	}

	msg := map[string]any{"role": "assistant"}
	if finishReason != "" {
		msg["finish_reason"] = finishReason
	}
	if len(p.toolCalls) > 0 {
		if len(p.texts) > 0 {
			msg["content"] = strings.Join(p.texts, "")
		} else {
			msg["content"] = nil
		}
		msg["tool_calls"] = p.toolCalls
	} else {
		msg["content"] = strings.Join(p.texts, "")
	}

	out, err := json.Marshal(map[string]any{"messages": []any{msg}})
	if err != nil {
		return "", finishReason
	}
	return string(out), finishReason
}

// modalityTokens sums tokenCount across a Gemini *TokensDetails array
// (e.g. promptTokensDetails) for entries matching the given modality
// (e.g. "AUDIO", "IMAGE", "TEXT").
func modalityTokens(usage map[string]any, field, modality string) int64 {
	arr, ok := usage[field].([]any)
	if !ok {
		return 0
	}
	var sum int64
	for _, e := range arr {
		m, ok := e.(map[string]any)
		if !ok {
			continue
		}
		if mod, _ := m["modality"].(string); mod == modality {
			if tc, ok := m["tokenCount"].(float64); ok {
				sum += int64(tc)
			}
		}
	}
	return sum
}

// buildUsageMetadata maps a Gemini usageMetadata object onto the LangSmith
// usage_metadata schema. It is pure (no span side-effects) and returns the
// assembled usage_metadata map along with the input, output, and total token
// totals (used for flat gen_ai.usage.* aggregation).
//
// The full Gemini usageMetadata schema (see the official GenerateContentResponse
// reference at https://ai.google.dev/api/generate-content#UsageMetadata) is:
//
//	"usageMetadata": {
//	  "promptTokenCount": 2095,            // [used] inclusive input total (cache is a subset)
//	  "cachedContentTokenCount": 1800,     // [used] prompt tokens served from cache
//	  "candidatesTokenCount": 503,         // [used] visible output, EXCLUDING thinking
//	  "thoughtsTokenCount": 64,            // [used] thinking/reasoning tokens
//	  "toolUsePromptTokenCount": 0,        // [ignored] tokens for tool-use prompts
//	  "totalTokenCount": 2662,             // [used] grand total (incl. tool-use)
//	  // The *Details arrays break each count down by modality (TEXT, IMAGE,
//	  // AUDIO, VIDEO, DOCUMENT). We pull AUDIO out of promptTokensDetails (it
//	  // has a separate input rate); the rest are not mapped — generated-image
//	  // output ("image") is priced per-image rather than per-token and has no
//	  // client convention yet, and no other modality is separately priced.
//	  "promptTokensDetails":      [{ "modality": "TEXT",  "tokenCount": 2095 }],
//	  "cacheTokensDetails":       [{ "modality": "TEXT",  "tokenCount": 1800 }],
//	  "candidatesTokensDetails":  [{ "modality": "TEXT",  "tokenCount": 503  }],
//	  "toolUsePromptTokensDetails": []
//	}
//
// candidatesTokenCount excludes thinking, so thoughts are folded into
// output_tokens and cachedContentTokenCount becomes a cache_read subset of
// input_tokens, mirroring langchain-google-genai (_response_to_result).
//
// Long-context tier: above 200k prompt tokens Google bills the whole request
// (input and output) at the higher rate, not just the overage
// (https://ai.google.dev/gemini-api/docs/pricing). We route everything into the
// over_200k / cache_read_over_200k buckets so every token is charged there.
func buildUsageMetadata(usage map[string]any) (usageMetadata map[string]any, totalInput, outputTokens, totalTokens int64) {
	getInt := func(key string) int64 {
		if v, ok := usage[key].(float64); ok {
			return int64(v)
		}
		return 0
	}

	totalInput = getInt("promptTokenCount")
	thoughts := getInt("thoughtsTokenCount")
	outputTokens = getInt("candidatesTokenCount") + thoughts
	cacheRead := getInt("cachedContentTokenCount")

	// Prefer the API's totalTokenCount (it also accounts for tool-use prompt
	// tokens); fall back to input+output when absent.
	totalTokens = getInt("totalTokenCount")
	if totalTokens == 0 {
		totalTokens = totalInput + outputTokens
	}

	// The long-context tier is a cliff keyed on the input prompt size; once
	// crossed it applies to output too. Threshold and high-tier rates come from
	// Google's rate card: https://ai.google.dev/gemini-api/docs/pricing
	// (Vertex AI: https://cloud.google.com/vertex-ai/generative-ai/pricing).
	const longContextThreshold = 200_000
	tier200k := totalInput > longContextThreshold

	inputTokenDetails := map[string]any{}
	switch {
	case tier200k:
		// Whole prompt at the high tier: cached portion at cache_read_over_200k,
		// the rest at over_200k. Together they sum to input_tokens (no leftover
		// charged at the base rate).
		if cacheRead > 0 {
			inputTokenDetails["cache_read_over_200k"] = cacheRead
		}
		if nonCached := totalInput - cacheRead; nonCached > 0 {
			inputTokenDetails["over_200k"] = nonCached
		}
	default:
		// cache_read and audio are independent subsets of input_tokens; the rest
		// (plain text) is charged at the base rate. (Audio isn't broken out in
		// the over_200k path: there's no audio_over_200k rate.)
		if cacheRead > 0 {
			inputTokenDetails["cache_read"] = cacheRead
		}
		// Audio input bills at a separate (higher) rate than text on some models.
		if audioInput := modalityTokens(usage, "promptTokensDetails", "AUDIO"); audioInput > 0 {
			inputTokenDetails["audio"] = audioInput
		}
	}

	outputTokenDetails := map[string]any{}
	switch {
	case tier200k:
		// Above the threshold all output (thinking included) bills at the
		// output over_200k rate. Gemini has no separate reasoning price, so we
		// put the whole output_tokens into over_200k rather than splitting out
		// reasoning, which would leave the buckets not summing to output_tokens.
		if outputTokens > 0 {
			outputTokenDetails["over_200k"] = outputTokens
		}
	case thoughts > 0:
		outputTokenDetails["reasoning"] = thoughts
	}

	usageMetadata = map[string]any{}
	if totalInput > 0 {
		usageMetadata["input_tokens"] = totalInput
	}
	if outputTokens > 0 {
		usageMetadata["output_tokens"] = outputTokens
	}
	if totalInput > 0 || outputTokens > 0 {
		usageMetadata["total_tokens"] = totalTokens
	}
	if len(inputTokenDetails) > 0 {
		usageMetadata["input_token_details"] = inputTokenDetails
	}
	if len(outputTokenDetails) > 0 {
		usageMetadata["output_token_details"] = outputTokenDetails
	}
	return usageMetadata, totalInput, outputTokens, totalTokens
}

// setUsageAttributes records token usage on the span: the cost-driving
// langsmith.usage_metadata JSON, the flat gen_ai.usage.* token counts, and the
// legacy underscore-format cache/reasoning attributes. See the inline notes for
// which consumer each group actually serves.
func setUsageAttributes(span trace.Span, usage map[string]any, parentSpan trace.Span) {
	usageMetadata, totalInput, outputTokens, totalTokens := buildUsageMetadata(usage)

	if len(usageMetadata) > 0 {
		if out, err := json.Marshal(usageMetadata); err == nil {
			span.SetAttributes(genaiattr.UsageMetadataKey.String(string(out)))
		}
	}

	// Flat gen_ai.usage.* attributes. On this span they are redundant for
	// LangSmith (the converter prefers the usage_metadata JSON above and ignores
	// these), but are kept for OpenTelemetry-standard consumers. Propagating
	// input/output to the parent IS load-bearing: the root span has no
	// usage_metadata of its own, so the converter's fallback path turns these
	// propagated attrs into the root run's token totals, which Thread-list stats
	// aggregate.
	setSelfAndParent := func(kv attribute.KeyValue) {
		span.SetAttributes(kv)
		if parentSpan.SpanContext().IsValid() && parentSpan.IsRecording() {
			parentSpan.SetAttributes(kv)
		}
	}
	if totalInput > 0 {
		setSelfAndParent(genaiattr.UsageInputTokensKey.Int64(totalInput))
	}
	if outputTokens > 0 {
		setSelfAndParent(genaiattr.UsageOutputTokensKey.Int64(outputTokens))
	}
	if totalTokens > 0 {
		span.SetAttributes(genaiattr.UsageTotalTokensKey.Int64(totalTokens))
	}

	// Legacy flat detail attributes (subsets of the totals above) for the
	// converter's underscore-format path.
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
