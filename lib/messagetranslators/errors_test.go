package messagetranslators

import (
	"errors"
	"fmt"
	"testing"
)

func TestSentinelErrorStrings(t *testing.T) {
	tests := []struct {
		err  error
		want string
	}{
		{ErrInvalidWireData, "invalid wire data"},
		{ErrUnsupported, "unsupported feature"},
		{ErrInvalidSequence, "invalid stream sequence"},
		{ErrTruncatedStream, "truncated stream"},
	}

	for _, test := range tests {
		if got := test.err.Error(); got != test.want {
			t.Errorf("%T.Error() = %q, want %q", test.err, got, test.want)
		}
	}
}

func TestConversionErrorFormattingAndUnwrap(t *testing.T) {
	tests := []struct {
		name   string
		err    error
		target error
		want   string
	}{
		{
			name:   "path",
			err:    at("$.messages[0].content", ErrUnsupported),
			target: ErrUnsupported,
			want:   "messagetranslators: $.messages[0].content: unsupported feature",
		},
		{
			name:   "empty path",
			err:    at("", ErrInvalidWireData),
			target: ErrInvalidWireData,
			want:   "messagetranslators: invalid wire data",
		},
		{
			name:   "wrapped detail",
			err:    at("$.stop_reason", fmt.Errorf("%w: pause_turn", ErrUnsupported)),
			target: ErrUnsupported,
			want:   "messagetranslators: $.stop_reason: unsupported feature: pause_turn",
		},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := test.err.Error(); got != test.want {
				t.Fatalf("Error() = %q, want %q", got, test.want)
			}
			if !errors.Is(test.err, test.target) {
				t.Fatalf("errors.Is(%v, %v) = false, want true", test.err, test.target)
			}
		})
	}
}
