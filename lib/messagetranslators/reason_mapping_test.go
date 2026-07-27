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
		{reason: "stop", hasTools: true, wantErr: ErrInvalidSequence},
		{reason: "tool_calls", wantErr: ErrInvalidSequence},
		{reason: "content_filter", wantErr: ErrUnsupported},
		{reason: "function_call", wantErr: ErrUnsupported},
		{reason: "unknown", wantErr: ErrInvalidWireData},
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
		{reason: "end_turn", hasTools: true, wantErr: ErrInvalidSequence},
		{reason: "tool_use", wantErr: ErrInvalidSequence},
		{reason: "pause_turn", wantErr: ErrUnsupported},
		{reason: "refusal", wantErr: ErrUnsupported},
		{reason: "unknown", wantErr: ErrInvalidWireData},
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
