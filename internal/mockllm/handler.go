package mockllm

import (
	"math/rand"
	"time"
)

// Message represents a generic chat message.
type Message struct {
	Role    string
	Content string
}

// ToolDef represents a tool definition from the request.
type ToolDef struct {
	Name        string
	Description string
}

// Request is the provider-agnostic input to a Handler.
type Request struct {
	Model    string
	Messages []Message
	Tools    []ToolDef
}

// ToolCall represents a tool/function call in the response.
type ToolCall struct {
	ID        string
	Name      string
	Arguments string // JSON string
}

// Response is the provider-agnostic output from a Handler.
type Response struct {
	Content      string
	ToolCalls    []ToolCall
	InputTokens  int
	OutputTokens int
	StopReason   string // servers translate to provider-specific values

	// Error causes the mock server to return an HTTP error response.
	Error *ResponseError

	// NetworkError causes the mock server to simulate a network-level failure
	// by hijacking the connection and closing it abruptly. When set, Error
	// and all other fields are ignored.
	NetworkError bool

	// TruncateStream causes the mock server to close the connection mid-stream
	// (only relevant for streaming requests). The server writes some SSE events
	// then drops the connection without sending the final stop/done events.
	TruncateStream bool

	// StreamDelay, when non-nil, is called before each SSE text chunk to
	// introduce a pause. This simulates realistic token-by-token delivery.
	// Return 0 to skip the delay for a given chunk.
	StreamDelay func() time.Duration
}

// ResponseError causes the mock server to return an HTTP error.
type ResponseError struct {
	Status  int
	Message string
}

// Handler generates a Response from a Request.
// The same handler works with both the Anthropic and OpenAI mock servers.
type Handler func(Request) Response

// DefaultHandler returns a simple text response, or a tool call if tools are present.
func DefaultHandler() Handler {
	return func(req Request) Response {
		if len(req.Tools) > 0 {
			return Response{
				Content: "I'll check that for you.",
				ToolCalls: []ToolCall{
					{
						ID:        "call_mock_1",
						Name:      req.Tools[0].Name,
						Arguments: `{"location":"Paris"}`,
					},
				},
				InputTokens:  25,
				OutputTokens: 18,
				StopReason:   "tool_use",
			}
		}
		return Response{
			Content:      "Hello! How can I help you today?",
			InputTokens:  25,
			OutputTokens: 15,
			StopReason:   "end_turn",
		}
	}
}

// StaticHandler returns a handler that always produces the given text.
func StaticHandler(text string) Handler {
	return func(req Request) Response {
		return Response{
			Content:      text,
			InputTokens:  10,
			OutputTokens: len(text) / 4, // rough approximation
			StopReason:   "end_turn",
		}
	}
}

// ErrorHandler returns a handler that always produces an HTTP error response
// with the given status code and message.
func ErrorHandler(status int, message string) Handler {
	return func(req Request) Response {
		return Response{
			Error: &ResponseError{Status: status, Message: message},
		}
	}
}

// NetworkErrorHandler returns a handler that simulates a network-level failure
// by closing the connection before any response is written.
func NetworkErrorHandler() Handler {
	return func(req Request) Response {
		return Response{NetworkError: true}
	}
}

// JitteredDelay returns a StreamDelay function that sleeps for a random
// duration between min and max before each streamed text chunk.
func JitteredDelay(min, max time.Duration) func() time.Duration {
	return func() time.Duration {
		jitter := time.Duration(rand.Int63n(int64(max - min)))
		return min + jitter
	}
}

// NoDelay returns a nil StreamDelay (no artificial delay).
func NoDelay() func() time.Duration { return nil }

// DefaultStreamDelay returns the default Eliza-style jittered delay
// (30–120ms per chunk, simulating ~10-30 tokens/sec).
func DefaultStreamDelay() func() time.Duration {
	return JitteredDelay(30*time.Millisecond, 120*time.Millisecond)
}

// TruncatedStreamHandler returns a handler that writes a partial streaming
// response then drops the connection, simulating a mid-stream network failure.
// For non-streaming requests it behaves like NetworkErrorHandler.
func TruncatedStreamHandler(text string) Handler {
	return func(req Request) Response {
		return Response{
			Content:        text,
			InputTokens:    10,
			OutputTokens:   5,
			TruncateStream: true,
		}
	}
}
