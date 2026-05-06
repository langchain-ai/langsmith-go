package langsmith

import (
	"context"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
)

// SnapshotWaitParams configures polling for snapshot readiness.
type SnapshotWaitParams struct {
	Timeout      time.Duration
	PollInterval time.Duration
}

// Wait polls until a snapshot reaches ready or failed status.
func (r *SandboxSnapshotService) Wait(ctx context.Context, snapshotID string, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	timeout := params.Timeout
	if timeout == 0 {
		timeout = 300 * time.Second
	}
	pollInterval := params.PollInterval
	if pollInterval == 0 {
		pollInterval = 2 * time.Second
	}
	deadline := time.Now().Add(timeout)
	lastStatus := ""

	for {
		snapshot, err := r.Get(ctx, snapshotID, opts...)
		if err != nil {
			return nil, err
		}
		lastStatus = snapshot.Status
		switch snapshot.Status {
		case "ready":
			return snapshot, nil
		case "failed":
			return nil, &SandboxResourceCreationError{
				ResourceType: "snapshot",
				ResourceID:   snapshotID,
				Message:      snapshot.StatusMessage,
			}
		}

		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, &SandboxResourceTimeoutError{
				ResourceType: "snapshot",
				ResourceID:   snapshotID,
				LastStatus:   lastStatus,
				Timeout:      timeout,
			}
		}
		delay := minDuration(pollInterval, remaining)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
}

// NewAndWait creates a snapshot and waits until it is ready or failed.
func (r *SandboxSnapshotService) NewAndWait(ctx context.Context, body SandboxSnapshotNewParams, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	snapshot, err := r.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.Wait(ctx, snapshot.ID, params, opts...)
}

// CaptureSnapshotAndWait captures a snapshot from a sandbox and waits until it
// is ready or failed.
func (r *SandboxBoxService) CaptureSnapshotAndWait(ctx context.Context, name string, body SandboxBoxNewSnapshotParams, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	snapshot, err := r.NewSnapshot(ctx, name, body, opts...)
	if err != nil {
		return nil, err
	}
	return NewSandboxSnapshotService(r.Options...).Wait(ctx, snapshot.ID, params, opts...)
}

// CaptureSnapshot captures a snapshot from this sandbox.
func (s *Sandbox) CaptureSnapshot(ctx context.Context, body SandboxBoxNewSnapshotParams, opts ...option.RequestOption) (*SandboxBoxNewSnapshotResponse, error) {
	return s.boxes.NewSnapshot(ctx, s.Name, body, opts...)
}

// CaptureSnapshotAndWait captures a snapshot and waits until it is ready.
func (s *Sandbox) CaptureSnapshotAndWait(ctx context.Context, body SandboxBoxNewSnapshotParams, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	return s.boxes.CaptureSnapshotAndWait(ctx, s.Name, body, params, opts...)
}
