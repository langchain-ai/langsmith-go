package langsmith

import (
	"context"

	"github.com/langchain-ai/langsmith-go/option"
)

// Run executes a command and waits for completion.
func (s *Sandbox) Run(ctx context.Context, body SandboxBoxRunParams, opts ...option.RequestOption) (*SandboxExecutionResult, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.RunWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// StartCommand starts a streaming command in this sandbox.
func (s *Sandbox) StartCommand(ctx context.Context, body SandboxCommandStartParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.StartCommandWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// RunWithCallbacks starts a WebSocket command, invokes callbacks for output,
// and waits for completion.
func (s *Sandbox) RunWithCallbacks(ctx context.Context, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxExecutionResult, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.RunWithDataplaneURLAndCallbacks(ctx, dataplaneURL, body, callbacks, opts...)
}

// ReconnectCommand reconnects to a command stream.
func (s *Sandbox) ReconnectCommand(ctx context.Context, commandID string, body SandboxCommandReconnectParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.ReconnectCommandWithDataplaneURL(ctx, dataplaneURL, commandID, body, opts...)
}
