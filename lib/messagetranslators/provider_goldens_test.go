package messagetranslators

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

type anthropicGolden struct {
	SchemaVersion int             `json:"schema_version"`
	Provider      string          `json:"provider"`
	API           string          `json:"api"`
	APIVersion    string          `json:"api_version"`
	Mode          string          `json:"mode"`
	Scenario      string          `json:"scenario"`
	Request       json.RawMessage `json:"request"`
	Response      json.RawMessage `json:"response"`
	Events        []struct {
		Event string          `json:"event"`
		Data  json.RawMessage `json:"data"`
	} `json:"events"`
}

func loadAnthropicGolden(t *testing.T, name string) anthropicGolden {
	t.Helper()
	body, err := os.ReadFile(filepath.Join("testdata", "provider_goldens", name))
	if err != nil {
		t.Fatal(err)
	}
	var fixture anthropicGolden
	if err := json.Unmarshal(body, &fixture); err != nil {
		t.Fatal(err)
	}
	if fixture.SchemaVersion != 1 || fixture.Provider != "anthropic" || fixture.API != "messages" || fixture.APIVersion != "2023-06-01" {
		t.Fatalf("unexpected fixture metadata: %#v", fixture)
	}
	return fixture
}

func TestAnthropicProviderGoldens(t *testing.T) {
	for _, name := range []string{
		"anthropic_messages_completed_text.json",
		"anthropic_messages_completed_tool.json",
		"anthropic_messages_stream_text.json",
		"anthropic_messages_stream_tool.json",
	} {
		t.Run(name, func(t *testing.T) {
			fixture := loadAnthropicGolden(t, name)
			warnings := &WarningCollector{}
			option := WithWarningHandler(warnings.HandleWarning)

			if _, err := AnthropicRequestToChatCompletions(fixture.Request, "", option); err != nil {
				t.Fatalf("request to Chat Completions: %v", err)
			}
			if _, err := AnthropicRequestToResponses(fixture.Request, "", option); err != nil {
				t.Fatalf("request to Responses: %v", err)
			}

			switch fixture.Mode {
			case "completed":
				if _, err := AnthropicResponseToChatCompletions(fixture.Response, "", option); err != nil {
					t.Fatalf("response to Chat Completions: %v", err)
				}
				if _, err := AnthropicResponseToResponses(fixture.Response, "", option); err != nil {
					t.Fatalf("response to Responses: %v", err)
				}
			case "stream":
				chat := NewAnthropicToChatCompletionsStream("", option)
				responses := NewAnthropicToResponsesStream("", option)
				for i, event := range fixture.Events {
					raw := SSEEvent{Event: event.Event, Data: event.Data}
					if _, err := chat.Convert(raw); err != nil {
						t.Fatalf("event %d to Chat Completions: %v", i, err)
					}
					if _, err := responses.Convert(raw); err != nil {
						t.Fatalf("event %d to Responses: %v", i, err)
					}
				}
				if err := chat.Finish(); err != nil {
					t.Fatalf("finish Chat Completions stream: %v", err)
				}
				if err := responses.Finish(); err != nil {
					t.Fatalf("finish Responses stream: %v", err)
				}
			default:
				t.Fatalf("unknown mode %q", fixture.Mode)
			}

			if got := warnings.Warnings(); len(got) != 0 {
				t.Fatalf("known provider fields produced warnings: %#v", got)
			}
		})
	}
}

func warningCount(warnings []Warning, code WarningCode) int {
	count := 0
	for _, warning := range warnings {
		if warning.Code == code {
			count++
		}
	}
	return count
}

func TestAnthropicAdditiveMetadataCompleted(t *testing.T) {
	base := `{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"tool_use","id":"c","name":"f","input":{"x":1},"caller":%s}],"stop_details":%s,"stop_reason":"tool_use","usage":{"input_tokens":0,"output_tokens":0}}`
	tests := []struct {
		name, caller, details string
		lossy, unknown        int
	}{
		{"direct-default", `{"type":"direct"}`, `null`, 0, 0},
		{"null-caller", `null`, `null`, 0, 0},
		{"populated-stop-details", `{"type":"direct"}`, `{"type":"refusal"}`, 1, 0},
		{"malformed-stop-details", `{"type":"direct"}`, `"bad"`, 1, 0},
		{"malformed-caller-scalar", `"direct"`, `null`, 1, 0},
		{"malformed-caller-missing-type", `{"tool_id":"x"}`, `null`, 1, 0},
		{"malformed-caller-type", `{"type":1}`, `null`, 1, 0},
		{"non-direct-caller", `{"type":"code_execution","tool_id":"x"}`, `null`, 1, 0},
		{"direct-caller-with-tool-id", `{"type":"direct","tool_id":"x"}`, `null`, 1, 0},
		{"direct-caller-with-unknown-metadata", `{"type":"direct","future":true}`, `null`, 1, 1},
	}
	converters := map[string]func([]byte, string, ...Option) ([]byte, error){
		"chat": AnthropicResponseToChatCompletions, "responses": AnthropicResponseToResponses,
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			body := []byte(fmt.Sprintf(base, tc.caller, tc.details))
			for name, convert := range converters {
				if _, err := convert(body, ""); err != nil {
					t.Errorf("%s without warning handler: %v", name, err)
				}
				warnings := &WarningCollector{}
				out, err := convert(body, "", WithWarningHandler(warnings.HandleWarning))
				if err != nil {
					t.Errorf("%s with warning handler: %v", name, err)
					continue
				}
				if !strings.Contains(string(out), `"name":"f"`) || !strings.Contains(string(out), `x`) {
					t.Errorf("%s lost core tool call fields: %s", name, out)
				}
				got := warnings.Warnings()
				if n := warningCount(got, WarningLossyConversion); n != tc.lossy {
					t.Errorf("%s lossy warnings = %d, want %d: %#v", name, n, tc.lossy, got)
				}
				if n := warningCount(got, WarningUnknownField); n != tc.unknown {
					t.Errorf("%s unknown warnings = %d, want %d: %#v", name, n, tc.unknown, got)
				}
			}
		})
	}
}

func TestAnthropicUnknownExtraResponseFieldIsIgnored(t *testing.T) {
	body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"ok"}],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":1},"future_metadata":{"anything":true}}`)
	for name, convert := range map[string]func([]byte, string, ...Option) ([]byte, error){
		"chat": AnthropicResponseToChatCompletions, "responses": AnthropicResponseToResponses,
	} {
		if _, err := convert(body, ""); err != nil {
			t.Errorf("%s without warning handler: %v", name, err)
		}
		warnings := &WarningCollector{}
		if _, err := convert(body, "", WithWarningHandler(warnings.HandleWarning)); err != nil {
			t.Errorf("%s with warning handler: %v", name, err)
		}
		got := warnings.Warnings()
		if warningCount(got, WarningUnknownField) != 1 || warningCount(got, WarningLossyConversion) != 0 {
			t.Errorf("%s warnings = %#v", name, got)
		}
	}
}

func TestAnthropicAdditiveCallerInRequestHistory(t *testing.T) {
	request := func(caller string) []byte {
		return []byte(`{"model":"a","max_tokens":1,"messages":[` +
			`{"role":"user","content":"call it"},` +
			`{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f","input":{"x":1},"caller":` + caller + `}]},` +
			`{"role":"user","content":[{"type":"tool_result","tool_use_id":"c","content":"ok"}]}` +
			`]}`)
	}
	for _, tc := range []struct {
		name, caller string
		lossy        int
	}{
		{"direct-default", `{"type":"direct"}`, 0},
		{"null", `null`, 0},
		{"malformed", `false`, 1},
		{"non-direct", `{"type":"code_execution","tool_id":"srv"}`, 1},
		{"direct-with-tool-id", `{"type":"direct","tool_id":"srv"}`, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for name, convert := range map[string]func([]byte, string, ...Option) ([]byte, error){
				"chat": AnthropicRequestToChatCompletions, "responses": AnthropicRequestToResponses,
			} {
				if _, err := convert(request(tc.caller), ""); err != nil {
					t.Errorf("%s without warning handler: %v", name, err)
				}
				warnings := &WarningCollector{}
				out, err := convert(request(tc.caller), "", WithWarningHandler(warnings.HandleWarning))
				if err != nil {
					t.Errorf("%s: %v", name, err)
					continue
				}
				if !strings.Contains(string(out), `"name":"f"`) || !strings.Contains(string(out), `x`) {
					t.Errorf("%s lost core tool call fields: %s", name, out)
				}
				if got := warningCount(warnings.Warnings(), WarningLossyConversion); got != tc.lossy {
					t.Errorf("%s lossy warnings = %d, want %d", name, got, tc.lossy)
				}
			}
		})
	}
}

type anthropicStreamConverter interface {
	Convert(SSEEvent) ([]SSEEvent, error)
}

func anthropicStreamConstructors() map[string]func(...Option) anthropicStreamConverter {
	return map[string]func(...Option) anthropicStreamConverter{
		"chat": func(options ...Option) anthropicStreamConverter {
			return NewAnthropicToChatCompletionsStream("", options...)
		},
		"responses": func(options ...Option) anthropicStreamConverter {
			return NewAnthropicToResponsesStream("", options...)
		},
	}
}

var anthropicTestStreamStart = SSEEvent{Event: "message_start", Data: []byte(`{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"a","content":[],"stop_details":null,"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":0,"output_tokens":0}}}`)}

func TestAnthropicAdditiveCallerStreaming(t *testing.T) {
	for _, tc := range []struct {
		name, caller string
		lossy        int
	}{
		{"direct-default", `{"type":"direct"}`, 0},
		{"null", `null`, 0},
		{"malformed-scalar", `false`, 1},
		{"malformed-object", `{}`, 1},
		{"non-direct", `{"type":"code_execution"}`, 1},
		{"direct-with-tool-id", `{"type":"direct","tool_id":"x"}`, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			data := []byte(`{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{},"caller":` + tc.caller + `}}`)
			for name, construct := range anthropicStreamConstructors() {
				plain := construct()
				if _, err := plain.Convert(anthropicTestStreamStart); err != nil {
					t.Fatalf("%s plain start: %v", name, err)
				}
				if _, err := plain.Convert(SSEEvent{Event: "content_block_start", Data: data}); err != nil {
					t.Errorf("%s without warning handler: %v", name, err)
				}

				warnings := &WarningCollector{}
				converter := construct(WithWarningHandler(warnings.HandleWarning))
				if _, err := converter.Convert(anthropicTestStreamStart); err != nil {
					t.Fatalf("%s start: %v", name, err)
				}
				out, err := converter.Convert(SSEEvent{Event: "content_block_start", Data: data})
				if err != nil {
					t.Errorf("%s with warning handler: %v", name, err)
					continue
				}
				encoded, _ := json.Marshal(out)
				if !strings.Contains(string(encoded), `f`) || !strings.Contains(string(encoded), `c`) {
					t.Errorf("%s lost core tool fields: %s", name, encoded)
				}
				if got := warningCount(warnings.Warnings(), WarningLossyConversion); got != tc.lossy {
					t.Errorf("%s lossy warnings = %d, want %d: %#v", name, got, tc.lossy, warnings.Warnings())
				}
			}
		})
	}
}

func TestAnthropicAdditiveStopDetailsStreaming(t *testing.T) {
	for _, tc := range []struct {
		name, details string
		lossy         int
	}{
		{"null", `null`, 0},
		{"populated", `{"type":"refusal"}`, 1},
		{"malformed", `false`, 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			for name, construct := range anthropicStreamConstructors() {
				data := []byte(`{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_details":` + tc.details + `},"usage":{"output_tokens":0}}`)
				plain := construct()
				if _, err := plain.Convert(anthropicTestStreamStart); err != nil {
					t.Fatalf("%s plain start: %v", name, err)
				}
				if _, err := plain.Convert(SSEEvent{Event: "message_delta", Data: data}); err != nil {
					t.Errorf("%s without warning handler: %v", name, err)
				}

				warnings := &WarningCollector{}
				converter := construct(WithWarningHandler(warnings.HandleWarning))
				if _, err := converter.Convert(anthropicTestStreamStart); err != nil {
					t.Fatalf("%s start: %v", name, err)
				}
				if _, err := converter.Convert(SSEEvent{Event: "message_delta", Data: data}); err != nil {
					t.Errorf("%s with warning handler: %v", name, err)
				}
				if got := warningCount(warnings.Warnings(), WarningLossyConversion); got != tc.lossy {
					t.Errorf("%s lossy warnings = %d, want %d: %#v", name, got, tc.lossy, warnings.Warnings())
				}
			}
		})
	}
}

func TestAnthropicAdditiveStopDetailsAtMessageStart(t *testing.T) {
	for _, details := range []string{`{}`, `"malformed"`} {
		t.Run(details, func(t *testing.T) {
			data := []byte(`{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"a","content":[],"stop_details":` + details + `,"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":0,"output_tokens":0}}}`)
			for name, construct := range anthropicStreamConstructors() {
				if _, err := construct().Convert(SSEEvent{Event: "message_start", Data: data}); err != nil {
					t.Errorf("%s without warning handler: %v", name, err)
				}
				warnings := &WarningCollector{}
				if _, err := construct(WithWarningHandler(warnings.HandleWarning)).Convert(SSEEvent{Event: "message_start", Data: data}); err != nil {
					t.Errorf("%s with warning handler: %v", name, err)
				}
				if got := warningCount(warnings.Warnings(), WarningLossyConversion); got != 1 {
					t.Errorf("%s lossy warnings = %d, want 1: %#v", name, got, warnings.Warnings())
				}
			}
		})
	}
}
