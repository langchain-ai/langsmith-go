package langsmith

import (
	"fmt"
	"time"
)

// SandboxDataplaneNotConfiguredError is returned when runtime operations are
// requested before the sandbox has a dataplane URL.
type SandboxDataplaneNotConfiguredError struct {
	SandboxName string
}

func (e *SandboxDataplaneNotConfiguredError) Error() string {
	if e.SandboxName == "" {
		return "langsmith: sandbox does not have a dataplane_url configured"
	}
	return fmt.Sprintf("langsmith: sandbox %q does not have a dataplane_url configured", e.SandboxName)
}

// SandboxNotReadyError is returned when runtime operations are requested for a
// sandbox that is not ready.
type SandboxNotReadyError struct {
	SandboxName string
	Status      string
}

func (e *SandboxNotReadyError) Error() string {
	if e.SandboxName == "" {
		return fmt.Sprintf("langsmith: sandbox is not ready (status: %s)", e.Status)
	}
	return fmt.Sprintf("langsmith: sandbox %q is not ready (status: %s)", e.SandboxName, e.Status)
}

// SandboxConnectionError is returned when a sandbox WebSocket command stream
// cannot be established or is interrupted unexpectedly.
type SandboxConnectionError struct {
	Message string
}

func (e *SandboxConnectionError) Error() string {
	return e.Message
}

// SandboxOperationError is returned when the sandbox dataplane reports a
// command operation error.
type SandboxOperationError struct {
	Operation string
	ErrorType string
	Message   string
}

func (e *SandboxOperationError) Error() string {
	if e.ErrorType != "" {
		return fmt.Sprintf("%s [%s]", e.Message, e.ErrorType)
	}
	return e.Message
}

// SandboxCommandTimeoutError is returned when the sandbox reports a command
// timeout.
type SandboxCommandTimeoutError struct {
	Message string
}

func (e *SandboxCommandTimeoutError) Error() string {
	return e.Message
}

// SandboxResourceTimeoutError is returned when waiting for a sandbox resource
// exceeds the configured timeout.
type SandboxResourceTimeoutError struct {
	ResourceType string
	ResourceID   string
	LastStatus   string
	Timeout      time.Duration
}

func (e *SandboxResourceTimeoutError) Error() string {
	if e.LastStatus != "" {
		return fmt.Sprintf("langsmith: %s %q not ready after %s (last_status: %s)", e.ResourceType, e.ResourceID, e.Timeout, e.LastStatus)
	}
	return fmt.Sprintf("langsmith: %s %q not ready after %s", e.ResourceType, e.ResourceID, e.Timeout)
}

// SandboxResourceCreationError is returned when a sandbox resource reaches a
// failed provisioning state.
type SandboxResourceCreationError struct {
	ResourceType string
	ResourceID   string
	Message      string
}

func (e *SandboxResourceCreationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("langsmith: %s %q failed: %s", e.ResourceType, e.ResourceID, e.Message)
	}
	return fmt.Sprintf("langsmith: %s %q failed", e.ResourceType, e.ResourceID)
}
