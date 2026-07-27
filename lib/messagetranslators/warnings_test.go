package messagetranslators

import (
	"errors"
	"sync"
	"testing"
)

func warningOptions() (ConversionOptions, *WarningCollector) {
	collector := &WarningCollector{}
	return ConversionOptions{WarningHandler: collector.HandleWarning}, collector
}

func hasWarning(warnings []Warning, path, field string) bool {
	for _, warning := range warnings {
		if warning.Code == WarningUnknownField && warning.Path == path && warning.Field == field {
			return true
		}
	}
	return false
}

func TestJSONWarningsAcrossProviderFamiliesAndDirections(t *testing.T) {
	anthropicRequest := []byte(`{"model":"a","max_tokens":2,"future_envelope":true,"messages":[{"role":"user","content":[{"type":"text","text":"hi","future_part":1}]}]}`)
	responsesRequest := []byte(`{"model":"r","max_output_tokens":2,"future_envelope":true,"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hi","future_part":1}]}]}`)
	chatRequest := []byte(`{"model":"g","max_tokens":2,"future_envelope":true,"messages":[{"role":"user","content":[{"type":"text","text":"hi","future_part":1}]}]}`)
	anthropicResponse := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"ok","future_part":1}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1},"future_envelope":true}`)
	responsesResponse := []byte(`{"id":"r","object":"response","status":"completed","model":"g","output":[{"id":"m","type":"message","role":"assistant","status":"completed","future_item":1,"content":[{"type":"output_text","text":"ok","annotations":[],"future_part":1}]}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2},"future_envelope":true}`)
	chatResponse := []byte(`{"id":"c","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"future_choice":1,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2},"future_envelope":true}`)

	tests := []struct {
		name       string
		body       []byte
		convert    func([]byte, string, ConversionOptions) ([]byte, error)
		nestedPath string
		nested     string
	}{
		{"anthropic-request-responses", anthropicRequest, AnthropicRequestToResponsesWithOptions, "$.messages[0].content[0].future_part", "future_part"},
		{"anthropic-request-chat", anthropicRequest, AnthropicRequestToChatCompletionsWithOptions, "$.messages[0].content[0].future_part", "future_part"},
		{"responses-request", responsesRequest, ResponsesRequestToAnthropicWithOptions, "$.input[0].content[0].future_part", "future_part"},
		{"chat-request", chatRequest, ChatCompletionsRequestToAnthropicWithOptions, "$.messages[0].content[0].future_part", "future_part"},
		{"responses-response", responsesResponse, ResponsesResponseToAnthropicWithOptions, "$.output[0].content[0].future_part", "future_part"},
		{"chat-response", chatResponse, ChatCompletionsResponseToAnthropicWithOptions, "$.choices[0].future_choice", "future_choice"},
		{"anthropic-response-responses", anthropicResponse, AnthropicResponseToResponsesWithOptions, "$.content[0].future_part", "future_part"},
		{"anthropic-response-chat", anthropicResponse, AnthropicResponseToChatCompletionsWithOptions, "$.content[0].future_part", "future_part"},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			options, collector := warningOptions()
			if _, err := test.convert(test.body, "", options); err != nil {
				t.Fatal(err)
			}
			warnings := collector.Warnings()
			if !hasWarning(warnings, "$.future_envelope", "future_envelope") {
				t.Fatalf("missing top-level warning: %#v", warnings)
			}
			if !hasWarning(warnings, test.nestedPath, test.nested) {
				t.Fatalf("missing nested warning: %#v", warnings)
			}
		})
	}
}

func TestWarningsDoNotDowngradeErrors(t *testing.T) {
	options, _ := warningOptions()
	if _, err := AnthropicRequestToResponsesWithOptions([]byte(`{"model":"a","max_tokens":1,"thinking":{},"messages":[{"role":"user","content":"x"}]}`), "", options); !errors.Is(err, ErrUnsupported) {
		t.Fatalf("known unsupported field: %v", err)
	}
	if _, err := ChatCompletionsRequestToAnthropicWithOptions([]byte(`{"model":"g","max_tokens":1}`), "", options); !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("missing required field: %v", err)
	}
}

func TestRepresentativePayloadsHaveNoWarningsAndLegacyWrappersWork(t *testing.T) {
	anthropicRequest := []byte(`{"model":"a","max_tokens":2,"messages":[{"role":"user","content":[{"type":"text","text":"hi"}]}]}`)
	responsesRequest := []byte(`{"model":"r","max_output_tokens":2,"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"hi"}]}]}`)
	chatRequest := []byte(`{"model":"g","max_tokens":2,"messages":[{"role":"user","content":"hi"}]}`)
	anthropicResponse := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"ok"}],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":1,"output_tokens":1}}`)
	responsesResponse := []byte(`{"id":"r","object":"response","status":"completed","model":"g","output":[{"id":"m","type":"message","role":"assistant","status":"completed","content":[{"type":"output_text","text":"ok","annotations":[]}]}],"usage":{"input_tokens":1,"output_tokens":1,"total_tokens":2}}`)
	chatResponse := []byte(`{"id":"c","object":"chat.completion","created":1,"model":"g","choices":[{"index":0,"message":{"role":"assistant","content":"ok"},"finish_reason":"stop"}],"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}}`)
	tests := []struct {
		body    []byte
		with    func([]byte, string, ConversionOptions) ([]byte, error)
		without func([]byte, string) ([]byte, error)
	}{
		{anthropicRequest, AnthropicRequestToResponsesWithOptions, AnthropicRequestToResponses},
		{anthropicRequest, AnthropicRequestToChatCompletionsWithOptions, AnthropicRequestToChatCompletions},
		{responsesRequest, ResponsesRequestToAnthropicWithOptions, ResponsesRequestToAnthropic},
		{chatRequest, ChatCompletionsRequestToAnthropicWithOptions, ChatCompletionsRequestToAnthropic},
		{responsesResponse, ResponsesResponseToAnthropicWithOptions, ResponsesResponseToAnthropic},
		{chatResponse, ChatCompletionsResponseToAnthropicWithOptions, ChatCompletionsResponseToAnthropic},
		{anthropicResponse, AnthropicResponseToResponsesWithOptions, AnthropicResponseToResponses},
		{anthropicResponse, AnthropicResponseToChatCompletionsWithOptions, AnthropicResponseToChatCompletions},
	}
	for n, test := range tests {
		options, collector := warningOptions()
		if _, err := test.with(test.body, "", options); err != nil {
			t.Fatalf("options conversion %d: %v", n, err)
		}
		if warnings := collector.Warnings(); len(warnings) != 0 {
			t.Fatalf("conversion %d warnings: %#v", n, warnings)
		}
		if _, err := test.without(test.body, ""); err != nil {
			t.Fatalf("legacy conversion %d: %v", n, err)
		}
	}
}

func TestStreamWarningAndFailedEventRollback(t *testing.T) {
	options, collector := warningOptions()
	s := NewChatCompletionsToAnthropicStreamWithOptions("claude", options)
	first := SSEEvent{Data: []byte(`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","future_envelope":1,"choices":[{"index":0,"delta":{"role":"assistant","future_delta":1},"finish_reason":null}]}`)}
	if _, err := s.Convert(first); err != nil {
		t.Fatal(err)
	}
	bad := SSEEvent{Data: []byte(`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{"content":7,"observed_before_failure":true},"finish_reason":null}]}`)}
	if _, err := s.Convert(bad); !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("bad event: %v", err)
	}
	good := SSEEvent{Data: []byte(`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{"content":"still usable"},"finish_reason":null}]}`)}
	if _, err := s.Convert(good); err != nil {
		t.Fatalf("state was not usable after failed event: %v", err)
	}
	warnings := collector.Warnings()
	for _, expected := range []struct{ path, field string }{
		{"event.future_envelope", "future_envelope"},
		{"event.choices[0].delta.future_delta", "future_delta"},
		{"event.choices[0].delta.observed_before_failure", "observed_before_failure"},
	} {
		if !hasWarning(warnings, expected.path, expected.field) {
			t.Fatalf("missing %s warning: %#v", expected.path, warnings)
		}
	}
}

func TestAllStreamConstructorsWithOptions(t *testing.T) {
	options := ConversionOptions{WarningHandler: func(Warning) {}}
	if NewResponsesToAnthropicStream("") == nil || NewResponsesToAnthropicStreamWithOptions("", options) == nil ||
		NewAnthropicToResponsesStream("") == nil || NewAnthropicToResponsesStreamWithOptions("", options) == nil ||
		NewChatCompletionsToAnthropicStream("") == nil || NewChatCompletionsToAnthropicStreamWithOptions("", options) == nil ||
		NewAnthropicToChatCompletionsStream("") == nil || NewAnthropicToChatCompletionsStreamWithOptions("", options) == nil {
		t.Fatal("constructor returned nil")
	}
}

func TestWarningHandlerIsSynchronousAndPanicPropagates(t *testing.T) {
	called := false
	options := ConversionOptions{WarningHandler: func(Warning) {
		called = true
		panic("handler panic")
	}}
	defer func() {
		if recover() == nil {
			t.Fatal("warning-handler panic was unexpectedly recovered")
		}
		if !called {
			t.Fatal("warning handler was not called synchronously")
		}
	}()
	_, _ = AnthropicRequestToResponsesWithOptions([]byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"new_field":true}`), "", options)
}

func TestWarningCollectorConcurrent(t *testing.T) {
	collector := &WarningCollector{}
	const goroutines, perGoroutine = 16, 100
	var wg sync.WaitGroup
	for n := 0; n < goroutines; n++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < perGoroutine; j++ {
				collector.HandleWarning(Warning{Code: WarningUnknownField, Field: "x"})
			}
		}()
	}
	wg.Wait()
	got := collector.Warnings()
	if len(got) != goroutines*perGoroutine {
		t.Fatalf("got %d warnings", len(got))
	}
	got[0].Field = "mutated"
	if collector.Warnings()[0].Field != "x" {
		t.Fatal("Warnings did not return a copy")
	}
}
