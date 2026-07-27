package messagetranslators

import (
	"errors"
	"strings"
	"testing"
)

// These tests pin the "Known incompatibilities" table in README.md. A change in
// either direction is a documentation change: a throw becoming best effort means
// something is now being silently dropped, and best effort becoming a throw
// means traffic that used to convert now fails.

type incompatCase struct {
	name    string
	body    string
	convert func([]byte, string, ...Option) ([]byte, error)
	// dropped, when set, must not appear anywhere in a successful conversion.
	dropped string
}

func runThrows(t *testing.T, cases []incompatCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.convert([]byte(tc.body), "")
			if out != nil {
				t.Fatalf("documented as throwing, but converted: %s", out)
			}
			if !errors.Is(err, ErrUnsupported) {
				t.Fatalf("documented as ErrUnsupported, got %v", err)
			}
		})
	}
}

func runBestEffort(t *testing.T, cases []incompatCase) {
	t.Helper()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			out, err := tc.convert([]byte(tc.body), "")
			if err != nil {
				t.Fatalf("documented as best effort, but threw: %v", err)
			}
			if tc.dropped != "" && strings.Contains(string(out), tc.dropped) {
				t.Fatalf("documented as dropped, but %q survived: %s", tc.dropped, out)
			}
		})
	}
}

func anthropicToolExchange(result string) string {
	return `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},` +
		`{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f","input":{}}]},` +
		`{"role":"user","content":[` + result + `]}]}`
}

func TestDocumentedThrowsAnthropicSource(t *testing.T) {
	req := func(content string) string {
		return `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":` + content + `}]}`
	}
	resp := func(content, stop string) string {
		return `{"id":"m","type":"message","role":"assistant","model":"a","content":` + content +
			`,"stop_reason":"` + stop + `","usage":{"input_tokens":1,"output_tokens":1}}`
	}
	runThrows(t, []incompatCase{
		{"thinking-block", req(`[{"type":"thinking","thinking":"x","signature":"s"}]`), AnthropicRequestToResponses, ""},
		{"redacted-thinking-block", req(`[{"type":"redacted_thinking","data":"x"}]`), AnthropicRequestToChatCompletions, ""},
		{"document-block", req(`[{"type":"document","source":{"type":"url","url":"https://x/y.pdf"}}]`), AnthropicRequestToResponses, ""},
		{"image-file-source", req(`[{"type":"image","source":{"type":"file","file_id":"f"}}]`), AnthropicRequestToResponses, ""},
		{"tool-result-is-error", anthropicToolExchange(`{"type":"tool_result","tool_use_id":"c","is_error":true,"content":"boom"}`), AnthropicRequestToResponses, ""},
		{"tool-result-image", anthropicToolExchange(`{"type":"tool_result","tool_use_id":"c","content":[{"type":"image","source":{"type":"url","url":"https://x"}}]}`), AnthropicRequestToResponses, ""},
		{"stop-sequences-to-responses", `{"model":"a","max_tokens":1,"stop_sequences":["a"],"messages":[{"role":"user","content":"hi"}]}`, AnthropicRequestToResponses, ""},
		{"over-four-stop-sequences-to-chat", `{"model":"a","max_tokens":1,"stop_sequences":["a","b","c","d","e"],"messages":[{"role":"user","content":"hi"}]}`, AnthropicRequestToChatCompletions, ""},
		{"text-after-tool-use-to-chat", resp(`[{"type":"tool_use","id":"c","name":"f","input":{}},{"type":"text","text":"after"}]`, "tool_use"), AnthropicResponseToChatCompletions, ""},
		{"pause-turn", resp(`[{"type":"text","text":"x"}]`, "pause_turn"), AnthropicResponseToChatCompletions, ""},
		{"refusal-to-responses", resp(`[{"type":"text","text":"x"}]`, "refusal"), AnthropicResponseToResponses, ""},
	})
}

func TestDocumentedThrowsResponsesSource(t *testing.T) {
	req := func(extra, input string) string {
		return `{"model":"r","max_output_tokens":1` + extra + `,"input":` + input + `}`
	}
	runThrows(t, []incompatCase{
		{"reasoning-item", req("", `[{"type":"reasoning","id":"rs","summary":[]}]`), ResponsesRequestToAnthropic, ""},
		{"web-search-call-item", req("", `[{"type":"web_search_call","id":"ws","status":"completed"}]`), ResponsesRequestToAnthropic, ""},
		{"input-file-part", req("", `[{"type":"message","role":"user","content":[{"type":"input_file","file_id":"f"}]}]`), ResponsesRequestToAnthropic, ""},
		{"image-detail", req("", `[{"type":"message","role":"user","content":[{"type":"input_image","image_url":"https://x","detail":"high"}]}]`), ResponsesRequestToAnthropic, ""},
		{"system-role-message", req("", `[{"type":"message","role":"system","content":"be nice"},{"type":"message","role":"user","content":"hi"}]`), ResponsesRequestToAnthropic, ""},
		{"strict-tool", req(`,"tools":[{"type":"function","name":"f","parameters":{},"strict":true}]`, `"hi"`), ResponsesRequestToAnthropic, ""},
		{"structured-function-output", req("", `[{"type":"message","role":"user","content":"u"},{"type":"function_call","call_id":"c","name":"f","arguments":"{}"},{"type":"function_call_output","call_id":"c","output":{"structured":true}}]`), ResponsesRequestToAnthropic, ""},
		{"incomplete-content-filter", `{"id":"r","object":"response","status":"incomplete","model":"g","incomplete_details":{"reason":"content_filter"},"output":[],"usage":{"input_tokens":1,"output_tokens":1}}`, ResponsesResponseToAnthropic, ""},
		{"truncated-function-call", `{"id":"r","object":"response","status":"incomplete","model":"g","incomplete_details":{"reason":"max_output_tokens"},"output":[{"id":"fc","type":"function_call","status":"incomplete","call_id":"c","name":"f","arguments":"{\"a\":"}],"usage":{"input_tokens":1,"output_tokens":1}}`, ResponsesResponseToAnthropic, ""},
	})
}

func TestDocumentedThrowsChatCompletionsSource(t *testing.T) {
	req := func(extra, msgs string) string {
		return `{"model":"g","max_tokens":1` + extra + `,"messages":` + msgs + `}`
	}
	completed := func(extra, message, finish string) string {
		return `{"id":"c","object":"chat.completion","created":1,"model":"g"` + extra +
			`,"choices":[{"index":0,"message":` + message + `,"finish_reason":"` + finish +
			`"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`
	}
	runThrows(t, []incompatCase{
		{"input-audio-part", req("", `[{"role":"user","content":[{"type":"input_audio","input_audio":{"data":"x","format":"wav"}}]}]`), ChatCompletionsRequestToAnthropic, ""},
		{"file-part", req("", `[{"role":"user","content":[{"type":"file","file":{"file_id":"f"}}]}]`), ChatCompletionsRequestToAnthropic, ""},
		{"image-detail", req("", `[{"role":"user","content":[{"type":"image_url","image_url":{"url":"https://x","detail":"high"}}]}]`), ChatCompletionsRequestToAnthropic, ""},
		{"message-name", req("", `[{"role":"user","content":"hi","name":"bob"}]`), ChatCompletionsRequestToAnthropic, ""},
		{"late-system-message", req("", `[{"role":"user","content":"hi"},{"role":"system","content":"late"}]`), ChatCompletionsRequestToAnthropic, ""},
		{"strict-tool", req(`,"tools":[{"type":"function","function":{"name":"f","parameters":{},"strict":true}}]`, `[{"role":"user","content":"hi"}]`), ChatCompletionsRequestToAnthropic, ""},
		{"n-two", req(`,"n":2`, `[{"role":"user","content":"hi"}]`), ChatCompletionsRequestToAnthropic, ""},
		{"two-choices", `{"id":"c","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"message":{"role":"assistant","content":"a"},"finish_reason":"stop"},{"index":1,"message":{"role":"assistant","content":"b"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`, ChatCompletionsResponseToAnthropic, ""},
		{"choice-logprobs", `{"id":"c","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"message":{"role":"assistant","content":"a"},"finish_reason":"stop","logprobs":{"content":[]}}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`, ChatCompletionsResponseToAnthropic, ""},
		{"message-refusal", completed("", `{"role":"assistant","content":null,"refusal":"no"}`, "stop"), ChatCompletionsResponseToAnthropic, ""},
		{"response-service-tier", completed(`,"service_tier":"default"`, `{"role":"assistant","content":"a"}`, "stop"), ChatCompletionsResponseToAnthropic, ""},
		{"audio-tokens", `{"id":"c","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"message":{"role":"assistant","content":"a"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2,"prompt_tokens_details":{"audio_tokens":3}}}`, ChatCompletionsResponseToAnthropic, ""},
	})
}

func TestDocumentedBestEffortDrops(t *testing.T) {
	runBestEffort(t, []incompatCase{
		{"anthropic-cache-control", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[{"type":"text","text":"hi","cache_control":{"type":"ephemeral"}}]}]}`, AnthropicRequestToResponses, "cache_control"},
		{"anthropic-stop-sequence-value", `{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"x"}],"stop_reason":"stop_sequence","stop_sequence":"END","usage":{"input_tokens":1,"output_tokens":1}}`, AnthropicResponseToResponses, "END"},
		{"anthropic-container-on-response", `{"id":"m","type":"message","role":"assistant","model":"a","container":{"id":"ctr_1"},"content":[{"type":"text","text":"x"}],"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`, AnthropicResponseToResponses, "ctr_1"},
		{"responses-reasoning-tokens", `{"id":"r","object":"response","status":"completed","model":"g","output":[{"id":"m","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok","annotations":[]}]}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2,"output_tokens_details":{"reasoning_tokens":50}}}`, ResponsesResponseToAnthropic, "reasoning"},
		{"chat-system-fingerprint", `{"id":"c","object":"chat.completion","created":1,"model":"g","system_fingerprint":"fp_1","choices":[{"index":0,"message":{"role":"assistant","content":"a"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`, ChatCompletionsResponseToAnthropic, "fp_1"},
		{"chat-reasoning-tokens", `{"id":"c","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"message":{"role":"assistant","content":"a"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2,"completion_tokens_details":{"reasoning_tokens":9}}}`, ChatCompletionsResponseToAnthropic, "reasoning"},
		// Anthropic text after a tool_use converts cleanly to Responses; only the
		// Chat Completions destination cannot express the ordering.
		{"text-after-tool-use-to-responses", `{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"tool_use","id":"c","name":"f","input":{}},{"type":"text","text":"after"}],"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":1}}`, AnthropicResponseToResponses, ""},
	})
}

// Multiple text blocks are joined without a separator in both directions.
func TestDocumentedConcatenationHasNoSeparator(t *testing.T) {
	out, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"system":[{"type":"text","text":"A"},{"type":"text","text":"B"}],"messages":[{"role":"user","content":"hi"}]}`), "")
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(out), `"instructions":"AB"`) {
		t.Fatalf("system blocks are documented as joined without a separator: %s", out)
	}
}

// Generated IDs are derived from the source, never passed through.
func TestDocumentedIDsAreNotPassedThrough(t *testing.T) {
	out, err := ChatCompletionsResponseToAnthropic([]byte(`{"id":"chatcmpl-upstream-123","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"message":{"role":"assistant","content":"a"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`), "")
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(out), "chatcmpl-upstream-123") {
		t.Fatalf("upstream ID is documented as replaced, not forwarded: %s", out)
	}
}

// A Chat Completions stream always ends with a usage chunk and [DONE], whether
// or not the client asked for usage.
func TestDocumentedTerminalUsageChunkIsAlwaysEmitted(t *testing.T) {
	s := NewAnthropicToChatCompletionsStream("m")
	for _, e := range []string{
		anthropicStartEvent,
		`{"type":"content_block_start","index":0,"content_block":{"type":"text","text":"hi"}}`,
		`{"type":"content_block_stop","index":0}`,
		`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":2}}`,
	} {
		if _, err := s.Convert(SSEEvent{Data: []byte(e)}); err != nil {
			t.Fatal(err)
		}
	}
	out, err := s.Convert(SSEEvent{Data: []byte(`{"type":"message_stop"}`)})
	if err != nil {
		t.Fatal(err)
	}
	if len(out) != 2 || string(out[1].Data) != "[DONE]" {
		t.Fatalf("want a usage chunk followed by [DONE], got %d events", len(out))
	}
	if !strings.Contains(string(out[0].Data), `"usage":{`) || !strings.Contains(string(out[0].Data), `"choices":[]`) {
		t.Fatalf("terminal chunk should carry usage and no choices: %s", out[0].Data)
	}
}
