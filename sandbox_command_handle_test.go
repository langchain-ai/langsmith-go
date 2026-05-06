package langsmith

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestSandboxCommandHandleNextContextCancellation(t *testing.T) {
	handle := &SandboxCommandHandle{
		chunks: make(chan SandboxOutputChunk),
	}
	ctx, cancel := context.WithTimeout(context.Background(), time.Millisecond)
	defer cancel()

	_, ok, err := handle.Next(ctx)
	if ok {
		t.Fatal("expected no chunk after context cancellation")
	}
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected context deadline error, got %T: %v", err, err)
	}
}

func TestSandboxCommandHandleErrAndDone(t *testing.T) {
	errDone := errors.New("stream failed")
	handle := &SandboxCommandHandle{
		done: make(chan struct{}),
	}
	handle.setErr(errDone)
	close(handle.done)

	if got := handle.Err(); !errors.Is(got, errDone) {
		t.Fatalf("expected stored error, got %v", got)
	}
	select {
	case <-handle.Done():
	case <-time.After(time.Second):
		t.Fatal("expected Done channel to be closed")
	}
}
