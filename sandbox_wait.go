package langsmith

import (
	"context"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
)

// SandboxWaitParams configures polling for sandbox readiness.
type SandboxWaitParams struct {
	Timeout      time.Duration
	PollInterval time.Duration
}

// Wait polls the generated status endpoint until the sandbox is ready or failed.
func (r *SandboxBoxService) Wait(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*SandboxBoxGetResponse, error) {
	timeout := params.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	pollInterval := params.PollInterval
	if pollInterval == 0 {
		pollInterval = time.Second
	}
	deadline := time.Now().Add(timeout)
	lastStatus := ""

	for {
		status, err := r.GetStatus(ctx, name, opts...)
		if err != nil {
			return nil, err
		}
		lastStatus = status.Status
		switch status.Status {
		case "ready":
			return r.Get(ctx, name, opts...)
		case "failed":
			return nil, &SandboxResourceCreationError{
				ResourceType: "sandbox",
				ResourceID:   name,
				Message:      status.StatusMessage,
			}
		}

		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, &SandboxResourceTimeoutError{
				ResourceType: "sandbox",
				ResourceID:   name,
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

// WaitSandbox polls until ready and returns the convenience wrapper.
func (r *SandboxBoxService) WaitSandbox(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.Wait(ctx, name, params, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromGetResponse(res, r), nil
}

// StartAndWait starts a stopped sandbox and waits until it is ready.
func (r *SandboxBoxService) StartAndWait(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*SandboxBoxGetResponse, error) {
	if _, err := r.Start(ctx, name, opts...); err != nil {
		return nil, err
	}
	return r.Wait(ctx, name, params, opts...)
}

// StartSandbox starts a stopped sandbox and returns the convenience wrapper once
// it is ready.
func (r *SandboxBoxService) StartSandbox(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.StartAndWait(ctx, name, params, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromGetResponse(res, r), nil
}
