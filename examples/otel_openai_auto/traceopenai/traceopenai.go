// Package traceopenai provides OpenTelemetry tracing for the
// github.com/sashabaranov/go-openai client using LangSmith-compatible spans.
//
// Usage:
//
//	// Configure your sashabaranov client to use a traced HTTP client
//	cfg := openai.DefaultConfig(apiKey)
//	cfg.HTTPClient = traceopenai.Client()
//	client := openai.NewClientWithConfig(cfg)
//
//	// Your OpenAI API calls will now be automatically traced with LangSmith attrs
//	// resp, err := client.CreateChatCompletion(ctx, ...)
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
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/trace"
)

// Client returns a new http.Client configured with tracing middleware.
// Equivalent to WrapClient(nil), which wraps the default transport.
func Client() *http.Client {
	return WrapClient(nil)
}

// WrapClient wraps an existing http.Client with tracing middleware.
// If client is nil, a new client with the default transport is created.
func WrapClient(client *http.Client) *http.Client {
	if client == nil {
		client = &http.Client{}
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = newRoundTripper(transport)
	return client
}

type roundTripper struct {
	base http.RoundTripper
}

func newRoundTripper(base http.RoundTripper) http.RoundTripper {
	return &roundTripper{base: base}
}

// RoundTrip intercepts requests/responses to add OpenTelemetry tracing.
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Only trace OpenAI API requests
	if !strings.Contains(req.URL.Host, "api.openai.com") {
		return rt.base.RoundTrip(req)
	}

	ctx := req.Context()
	tracer := otel.Tracer("github.com/sashabaranov/go-openai")

	// Extract span context from request headers
	ctx = otel.GetTextMapPropagator().Extract(ctx, propagation.HeaderCarrier(req.Header))

	// Determine span name based on endpoint
	spanName := getSpanName(req.URL.Path)

	// Start span
	ctx, span := tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("gen_ai.system", "openai"),
			attribute.String("gen_ai.operation.name", getOperationName(req.URL.Path)),
			attribute.String("http.method", req.Method),
			attribute.String("http.url", req.URL.String()),
		),
	)
	defer span.End()

	// Read request body if present
	var requestBody []byte
	if req.Body != nil {
		var err error
		requestBody, err = io.ReadAll(req.Body)
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, fmt.Sprintf("failed to read request body: %v", err))
			return rt.base.RoundTrip(req)
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

	// Make the actual request
	resp, err := rt.base.RoundTrip(req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return resp, err
	}

	// Read response body
	var responseBody []byte
	if resp.Body != nil {
		var readErr error
		responseBody, readErr = io.ReadAll(resp.Body)
		if readErr != nil {
			span.RecordError(readErr)
			span.SetStatus(codes.Error, fmt.Sprintf("failed to read response body: %v", readErr))
			// Continue with empty body rather than failing the request
		} else {
			resp.Body = io.NopCloser(bytes.NewBuffer(responseBody))
		}
	}

	// Extract completion and usage from response
	if len(responseBody) > 0 {
		completion, usage := extractCompletionFromResponse(responseBody)
		if completion != "" {
			span.SetAttributes(attribute.String("gen_ai.completion", completion))
		}
		if usage.InputTokens > 0 {
			span.SetAttributes(
				attribute.Int64("gen_ai.usage.input_tokens", int64(usage.InputTokens)),
				attribute.Int64("gen_ai.usage.output_tokens", int64(usage.OutputTokens)),
			)
		}
	}

	// Set status based on HTTP status code
	if resp.StatusCode >= 400 {
		span.SetStatus(codes.Error, fmt.Sprintf("HTTP %d", resp.StatusCode))
	} else {
		span.SetStatus(codes.Ok, "")
	}

	return resp, err
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
