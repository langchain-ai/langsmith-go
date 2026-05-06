package langsmith

import (
	"testing"
	"time"
)

func TestSandboxErrorMessages(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want string
	}{
		{
			name: "dataplane not configured with sandbox name",
			err:  &SandboxDataplaneNotConfiguredError{SandboxName: "box-a"},
			want: `langsmith: sandbox "box-a" does not have a dataplane_url configured`,
		},
		{
			name: "sandbox not ready",
			err:  &SandboxNotReadyError{SandboxName: "box-a", Status: "starting"},
			want: `langsmith: sandbox "box-a" is not ready (status: starting)`,
		},
		{
			name: "connection error",
			err:  &SandboxConnectionError{Message: "dial failed"},
			want: "dial failed",
		},
		{
			name: "operation error includes type",
			err:  &SandboxOperationError{ErrorType: "CommandError", Message: "bad command"},
			want: "bad command [CommandError]",
		},
		{
			name: "command timeout",
			err:  &SandboxCommandTimeoutError{Message: "timed out"},
			want: "timed out",
		},
		{
			name: "resource timeout includes last status",
			err:  &SandboxResourceTimeoutError{ResourceType: "sandbox", ResourceID: "box-a", LastStatus: "starting", Timeout: 2 * time.Second},
			want: `langsmith: sandbox "box-a" not ready after 2s (last_status: starting)`,
		},
		{
			name: "resource creation includes message",
			err:  &SandboxResourceCreationError{ResourceType: "snapshot", ResourceID: "snap-a", Message: "builder failed"},
			want: `langsmith: snapshot "snap-a" failed: builder failed`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.err.Error(); got != tt.want {
				t.Fatalf("expected %q, got %q", tt.want, got)
			}
		})
	}
}
