package messagetranslators

import (
	"encoding/json"
	"strings"
	"testing"
	"time"
)

const (
	anthropicStartEvent    = `{"type":"message_start","message":{"id":"m1","type":"message","role":"assistant","content":[],"model":"claude-x","stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":5,"output_tokens":1}}}`
	anthropicCompletedBody = `{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"ok"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`
)

// A Responses response.created carries no output yet, but the field is a
// required list: emitting null there fails strict client-side validation.
func TestResponsesCreatedCarriesEmptyOutputArray(t *testing.T) {
	s := NewAnthropicToResponsesStream("")
	out, err := s.Convert(SSEEvent{Data: []byte(anthropicStartEvent)})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 {
		t.Fatalf("got %d events", len(out))
	}
	for _, e := range out {
		var decoded struct {
			Response struct {
				Output *[]any `json:"output"`
			} `json:"response"`
		}
		if err := json.Unmarshal(e.Data, &decoded); err != nil {
			t.Fatal(err)
		}
		if decoded.Response.Output == nil {
			t.Fatalf("%s emitted a null output: %s", e.Event, e.Data)
		}
	}
}

func TestWithClockMakesGeneratedTimestampsDeterministic(t *testing.T) {
	fixed := time.Unix(1700000000, 0)
	clock := WithClock(func() time.Time { return fixed })

	out, err := AnthropicResponseToResponses([]byte(anthropicCompletedBody), "", clock)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"created_at":1700000000`) {
		t.Fatalf("clock was not applied: %s", out)
	}

	chat, err := AnthropicResponseToChatCompletions([]byte(anthropicCompletedBody), "", clock)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(chat), `"created":1700000000`) {
		t.Fatalf("clock was not applied: %s", chat)
	}

	s := NewAnthropicToChatCompletionsStream("m", clock)
	events, err := s.Convert(SSEEvent{Data: []byte(anthropicStartEvent)})
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(events[0].Data), `"created":1700000000`) {
		t.Fatalf("clock was not applied to the stream: %s", events[0].Data)
	}
}

// Input usage arrives only in a Chat Completions terminal usage chunk, after
// message_start has shipped. It must still be recoverable.
func TestStreamUsageIsRecoverableAfterMessageStart(t *testing.T) {
	var reported []Usage
	s := NewChatCompletionsToAnthropicStream("", WithUsageHandler(func(u Usage) { reported = append(reported, u) }))
	for _, chunk := range []string{
		`{"id":"c1","object":"chat.completion.chunk","created":1,"model":"m","choices":[{"index":0,"delta":{"role":"assistant","content":""},"finish_reason":null}]}`,
		`{"id":"c1","object":"chat.completion.chunk","created":1,"model":"m","choices":[{"index":0,"delta":{"content":"hi"},"finish_reason":null}]}`,
		`{"id":"c1","object":"chat.completion.chunk","created":1,"model":"m","choices":[{"index":0,"delta":{},"finish_reason":"stop"}]}`,
		`{"id":"c1","object":"chat.completion.chunk","created":1,"model":"m","choices":[],"usage":{"prompt_tokens":100,"completion_tokens":7,"total_tokens":107,"prompt_tokens_details":{"cached_tokens":40}}}`,
		`[DONE]`,
	} {
		if _, err := s.Convert(SSEEvent{Data: []byte(chunk)}); err != nil {
			t.Fatal(err)
		}
	}
	want := Usage{InputTokens: 60, OutputTokens: 7, CacheReadInputTokens: 40}
	if got := s.Usage(); got != want {
		t.Fatalf("Usage() = %+v, want %+v", got, want)
	}
	if len(reported) != 1 || reported[0] != want {
		t.Fatalf("usage handler got %+v, want exactly one %+v", reported, want)
	}
}

func TestResponsesStreamReportsTerminalOnlyInputUsage(t *testing.T) {
	s := NewResponsesToAnthropicStream("")
	for _, e := range []SSEEvent{
		{Event: "response.created", Data: []byte(`{"type":"response.created","response":{"id":"r1","model":"gpt","status":"in_progress"}}`)},
		{Event: "response.completed", Data: []byte(`{"type":"response.completed","response":{"id":"r1","model":"gpt","status":"completed","output":[],"usage":{"input_tokens":11,"output_tokens":3,"total_tokens":14,"input_tokens_details":{"cached_tokens":4}}}}`)},
	} {
		if _, err := s.Convert(e); err != nil {
			t.Fatal(err)
		}
	}
	want := Usage{InputTokens: 7, OutputTokens: 3, CacheReadInputTokens: 4}
	if got := s.Usage(); got != want {
		t.Fatalf("Usage() = %+v, want %+v", got, want)
	}
}

// SDKs send documented defaults explicitly; those must not be rejected, while
// the same field carrying a real request still must be.
func TestDocumentedDefaultsAreAcceptedAndRealRequestsRejected(t *testing.T) {
	chat := func(extra string) []byte {
		return []byte(`{"model":"g","max_tokens":1,` + extra + `,"messages":[{"role":"user","content":"hi"}]}`)
	}
	responses := func(extra string) []byte {
		return []byte(`{"model":"r","max_output_tokens":1,` + extra + `,"input":"hi"}`)
	}
	for _, tc := range []struct {
		name    string
		body    []byte
		convert func([]byte, string, ...Option) ([]byte, error)
		wantErr bool
	}{
		{"cc-parallel-default", chat(`"parallel_tool_calls":true`), ChatCompletionsRequestToAnthropic, false},
		{"cc-parallel-disabled", chat(`"parallel_tool_calls":false`), ChatCompletionsRequestToAnthropic, true},
		{"cc-format-text", chat(`"response_format":{"type":"text"}`), ChatCompletionsRequestToAnthropic, false},
		{"cc-format-json", chat(`"response_format":{"type":"json_object"}`), ChatCompletionsRequestToAnthropic, true},
		{"cc-zero-penalties", chat(`"frequency_penalty":0,"presence_penalty":0`), ChatCompletionsRequestToAnthropic, false},
		{"cc-real-penalty", chat(`"frequency_penalty":0.5`), ChatCompletionsRequestToAnthropic, true},
		{"cc-logprobs-off", chat(`"logprobs":false`), ChatCompletionsRequestToAnthropic, false},
		{"cc-logprobs-on", chat(`"logprobs":true`), ChatCompletionsRequestToAnthropic, true},
		{"cc-null-is-unset", chat(`"seed":null,"user":null`), ChatCompletionsRequestToAnthropic, false},
		{"resp-truncation-default", responses(`"truncation":"disabled"`), ResponsesRequestToAnthropic, false},
		{"resp-truncation-auto", responses(`"truncation":"auto"`), ResponsesRequestToAnthropic, true},
		{"resp-text-default", responses(`"text":{"format":{"type":"text"}}`), ResponsesRequestToAnthropic, false},
		{"resp-text-schema", responses(`"text":{"format":{"type":"json_schema","name":"x"}}`), ResponsesRequestToAnthropic, true},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.convert(tc.body, "")
			if tc.wantErr {
				requireErrorIs(t, err, ErrUnsupported)
				return
			}
			if err != nil {
				t.Fatalf("documented default was rejected: %v", err)
			}
		})
	}
}

// Fields that change what the provider does must never be dropped in silence,
// in either direction.
func TestSemanticallySignificantFieldsAreRejectedNotDropped(t *testing.T) {
	for _, tc := range []struct {
		name    string
		body    []byte
		convert func([]byte, string, ...Option) ([]byte, error)
	}{
		{"responses-background", []byte(`{"model":"r","max_output_tokens":1,"background":true,"input":"hi"}`), ResponsesRequestToAnthropic},
		{"responses-conversation", []byte(`{"model":"r","max_output_tokens":1,"conversation":"conv_1","input":"hi"}`), ResponsesRequestToAnthropic},
		{"responses-include", []byte(`{"model":"r","max_output_tokens":1,"include":["reasoning.encrypted_content"],"input":"hi"}`), ResponsesRequestToAnthropic},
		{"anthropic-mcp-to-responses", []byte(`{"model":"a","max_tokens":1,"mcp_servers":[{"name":"x"}],"messages":[{"role":"user","content":"hi"}]}`), AnthropicRequestToResponses},
		{"anthropic-context-to-responses", []byte(`{"model":"a","max_tokens":1,"context_management":{"edits":[]},"messages":[{"role":"user","content":"hi"}]}`), AnthropicRequestToResponses},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.convert(tc.body, "")
			if out != nil {
				t.Fatalf("field was silently dropped: %s", out)
			}
			requireErrorIs(t, err, ErrUnsupported)
		})
	}
}

// The same Anthropic body must get the same verdict regardless of destination.
func TestAnthropicRequestPolicyIsConsistentAcrossDestinations(t *testing.T) {
	for _, field := range []string{"mcp_servers", "context_management", "thinking", "output_config", "container"} {
		body := []byte(`{"model":"a","max_tokens":1,"` + field + `":{"x":1},"messages":[{"role":"user","content":"hi"}]}`)
		_, toResponses := AnthropicRequestToResponses(body, "")
		_, toChat := AnthropicRequestToChatCompletions(body, "")
		if (toResponses == nil) != (toChat == nil) {
			t.Fatalf("%s: Responses err=%v but Chat Completions err=%v", field, toResponses, toChat)
		}
	}
}

func TestLossyConversionWarningsAreEmitted(t *testing.T) {
	collector := &WarningCollector{}
	body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"ok","cache_control":{"type":"ephemeral"}}],"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1,"cache_creation_input_tokens":9}}`)
	if _, err := AnthropicResponseToResponses(body, "", WithWarningHandler(collector.HandleWarning)); err != nil {
		t.Fatal(err)
	}
	var lossy []string
	for _, w := range collector.Warnings() {
		if w.Code == WarningLossyConversion {
			lossy = append(lossy, w.Field)
		}
	}
	if len(lossy) != 2 {
		t.Fatalf("want cache_control and cache_creation_input_tokens warnings, got %v", lossy)
	}
}
