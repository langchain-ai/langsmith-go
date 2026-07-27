package messagetranslators

import (
	"errors"
	"testing"
)

func TestChatCompletionsFinishReasonToAnthropic(t *testing.T) {
	tests := []struct {
		reason   string
		hasTools bool
		want     string
		wantErr  error
	}{
		{reason: "stop", want: "end_turn"},
		{reason: "tool_calls", hasTools: true, want: "tool_use"},
		{reason: "length", want: "max_tokens"},
		{reason: "content_filter", want: "refusal"},
		{reason: "function_call", wantErr: ErrUnsupported},
		{reason: "unknown", wantErr: ErrInvalidWireData},
		// Content wins over an advisory finish_reason that disagrees with it, so
		// that OpenAI-compatible servers reporting "stop" beside a populated
		// tool_calls array stay translatable.
		{reason: "stop", hasTools: true, want: "tool_use"},
		{reason: "tool_calls", want: "end_turn"},
	}
	for _, test := range tests {
		t.Run(test.reason, func(t *testing.T) {
			got, err := chatCompletionsFinishReasonToAnthropic(test.reason, test.hasTools, "reason")
			if got != test.want || !errors.Is(err, test.wantErr) {
				t.Fatalf("got (%q, %v), want (%q, %v)", got, err, test.want, test.wantErr)
			}
		})
	}
}

func TestAnthropicStopReasonToChatCompletions(t *testing.T) {
	tests := []struct {
		reason   string
		hasTools bool
		want     string
		wantErr  error
	}{
		{reason: "end_turn", want: "stop"},
		{reason: "stop_sequence", want: "stop"},
		{reason: "tool_use", hasTools: true, want: "tool_calls"},
		{reason: "max_tokens", want: "length"},
		{reason: "refusal", want: "content_filter"},
		{reason: "pause_turn", wantErr: ErrUnsupported},
		{reason: "unknown", wantErr: ErrInvalidWireData},
		// Content wins over a disagreeing stop_reason, as above.
		{reason: "end_turn", hasTools: true, want: "tool_calls"},
		{reason: "tool_use", want: "stop"},
	}
	for _, test := range tests {
		t.Run(test.reason, func(t *testing.T) {
			got, err := anthropicStopReasonToChatCompletions(test.reason, test.hasTools, "reason")
			if got != test.want || !errors.Is(err, test.wantErr) {
				t.Fatalf("got (%q, %v), want (%q, %v)", got, err, test.want, test.wantErr)
			}
		})
	}
}
