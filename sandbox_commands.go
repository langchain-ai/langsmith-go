package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"golang.org/x/net/websocket"
)

const (
	defaultSandboxCommandTimeoutSeconds = int64(60)
	defaultSandboxCommandShell          = "/bin/bash"
	defaultSandboxCommandIdleTimeout    = int64(300)
	defaultSandboxCommandTTLSeconds     = int64(600)
	sandboxCommandMaxAutoReconnects     = 5
	sandboxCommandReconnectBackoffBase  = 500 * time.Millisecond
	sandboxCommandReconnectBackoffMax   = 8 * time.Second
)

// SandboxExecutionResult is the result of executing a command in a sandbox.
type SandboxExecutionResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int64  `json:"exit_code"`
}

// Success reports whether the command exited with code 0.
func (r SandboxExecutionResult) Success() bool {
	return r.ExitCode == 0
}

// SandboxBoxRunParams configures a blocking sandbox command execution.
type SandboxBoxRunParams struct {
	Command param.Field[string]            `json:"command" api:"required"`
	Timeout param.Field[int64]             `json:"timeout"`
	Env     param.Field[map[string]string] `json:"env"`
	CWD     param.Field[string]            `json:"cwd"`
	Shell   param.Field[string]            `json:"shell"`
}

func (r SandboxBoxRunParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// SandboxCommandStartParams configures a streaming WebSocket command execution.
// Command is optional when Pty is true, in which case the sandbox starts the
// selected shell directly.
type SandboxCommandStartParams struct {
	Command            param.Field[string]            `json:"command" api:"required"`
	TimeoutSeconds     param.Field[int64]             `json:"timeout_seconds"`
	Env                param.Field[map[string]string] `json:"env"`
	CWD                param.Field[string]            `json:"cwd"`
	Shell              param.Field[string]            `json:"shell"`
	IdleTimeoutSeconds param.Field[int64]             `json:"idle_timeout_seconds"`
	KillOnDisconnect   param.Field[bool]              `json:"kill_on_disconnect"`
	TTLSeconds         param.Field[int64]             `json:"ttl_seconds"`
	Pty                param.Field[bool]              `json:"pty"`
	SSHAgentForward    param.Field[bool]              `json:"ssh_agent_forward"`
}

func (r SandboxCommandStartParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// SandboxCommandReconnectParams configures a command stream reconnection.
type SandboxCommandReconnectParams struct {
	StdoutOffset param.Field[int64] `json:"stdout_offset"`
	StderrOffset param.Field[int64] `json:"stderr_offset"`
}

func (r SandboxCommandReconnectParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// SandboxOutputChunk is a single chunk of streaming command output.
type SandboxOutputChunk struct {
	Stream string `json:"stream"`
	Data   string `json:"data"`
	Offset int64  `json:"offset"`
}

// SandboxCommandCallbacks are invoked as streaming command output arrives.
type SandboxCommandCallbacks struct {
	OnStdout        func(string)
	OnStderr        func(string)
	OnSSHAgentData  func(channelID string, data []byte)
	OnSSHAgentClose func(channelID string)
}

// Run executes a command in the named sandbox and waits for completion. The
// sandbox is fetched first so the current dataplane URL can be used.
func (r *SandboxBoxService) Run(ctx context.Context, name string, body SandboxBoxRunParams, opts ...option.RequestOption) (res *SandboxExecutionResult, err error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.RunWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// RunWithDataplaneURL executes a command against a sandbox dataplane URL and
// waits for completion.
func (r *SandboxBoxService) RunWithDataplaneURL(ctx context.Context, dataplaneURL string, body SandboxBoxRunParams, opts ...option.RequestOption) (res *SandboxExecutionResult, err error) {
	opts = slices.Concat(r.Options, opts)
	body, timeoutSeconds, err := normalizeSandboxRunParams(body)
	if err != nil {
		return nil, err
	}
	path, err := sandboxDataplaneURL(dataplaneURL, "execute")
	if err != nil {
		return nil, err
	}

	requestTimeout := time.Duration(timeoutSeconds+10) * time.Second
	opts = slices.Concat([]option.RequestOption{option.WithRequestTimeout(requestTimeout)}, opts)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// StartCommand starts a streaming command in the named sandbox. The returned
// handle can stream output, send stdin, kill the command, and reconnect.
func (r *SandboxBoxService) StartCommand(ctx context.Context, name string, body SandboxCommandStartParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.StartCommandWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// StartCommandWithDataplaneURL starts a streaming command against a sandbox
// dataplane URL.
func (r *SandboxBoxService) StartCommandWithDataplaneURL(ctx context.Context, dataplaneURL string, body SandboxCommandStartParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	return r.startCommandWithDataplaneURL(ctx, dataplaneURL, body, SandboxCommandCallbacks{}, opts...)
}

// RunWithCallbacks starts a command in the named sandbox over WebSocket, invokes
// callbacks for streamed output, and waits for completion.
func (r *SandboxBoxService) RunWithCallbacks(ctx context.Context, name string, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxExecutionResult, error) {
	handle, err := r.StartCommandWithCallbacks(ctx, name, body, callbacks, opts...)
	if err != nil {
		return nil, err
	}
	return handle.Result(ctx)
}

// RunWithDataplaneURLAndCallbacks starts a command over WebSocket against a
// dataplane URL, invokes callbacks for streamed output, and waits for completion.
func (r *SandboxBoxService) RunWithDataplaneURLAndCallbacks(ctx context.Context, dataplaneURL string, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxExecutionResult, error) {
	handle, err := r.StartCommandWithDataplaneURLAndCallbacks(ctx, dataplaneURL, body, callbacks, opts...)
	if err != nil {
		return nil, err
	}
	return handle.Result(ctx)
}

// StartCommandWithCallbacks starts a streaming command in the named sandbox and
// invokes callbacks for stdout/stderr chunks as they arrive.
func (r *SandboxBoxService) StartCommandWithCallbacks(ctx context.Context, name string, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.StartCommandWithDataplaneURLAndCallbacks(ctx, dataplaneURL, body, callbacks, opts...)
}

// StartCommandWithDataplaneURLAndCallbacks starts a streaming command against a
// dataplane URL and invokes callbacks for stdout/stderr chunks as they arrive.
func (r *SandboxBoxService) StartCommandWithDataplaneURLAndCallbacks(ctx context.Context, dataplaneURL string, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	return r.startCommandWithDataplaneURL(ctx, dataplaneURL, body, callbacks, opts...)
}

func (r *SandboxBoxService) startCommandWithDataplaneURL(ctx context.Context, dataplaneURL string, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	opts = slices.Concat(r.Options, opts)
	payload, err := normalizeSandboxCommandStartParams(body)
	if err != nil {
		return nil, err
	}

	ws, err := dialSandboxCommandWebSocket(ctx, dataplaneURL, opts...)
	if err != nil {
		return nil, err
	}
	if err := websocket.JSON.Send(ws, payload); err != nil {
		_ = ws.Close()
		return nil, &SandboxConnectionError{Message: fmt.Sprintf("langsmith: failed to send sandbox command request: %v", err)}
	}

	started, err := receiveSandboxWSMessage(ctx, ws)
	if err != nil {
		_ = ws.Close()
		return nil, err
	}
	if started.Type == "error" {
		_ = ws.Close()
		return nil, sandboxErrorFromWSMessage(started, "")
	}
	if started.Type != "started" {
		_ = ws.Close()
		return nil, &SandboxOperationError{
			Operation: "command",
			Message:   fmt.Sprintf("expected sandbox command started message, got %q", started.Type),
		}
	}

	handle := newSandboxCommandHandle(ws, dataplaneURL, opts, started.CommandID, started.PID, 0, 0)
	handle.callbacks = callbacks
	handle.start()
	return handle, nil
}

// ReconnectCommand reconnects to a running or recently-finished command in the
// named sandbox.
func (r *SandboxBoxService) ReconnectCommand(ctx context.Context, name string, commandID string, body SandboxCommandReconnectParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	if commandID == "" {
		return nil, errors.New("missing required commandID parameter")
	}
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.ReconnectCommandWithDataplaneURL(ctx, dataplaneURL, commandID, body, opts...)
}

// ReconnectCommandWithDataplaneURL reconnects to a command stream against a
// sandbox dataplane URL.
func (r *SandboxBoxService) ReconnectCommandWithDataplaneURL(ctx context.Context, dataplaneURL string, commandID string, body SandboxCommandReconnectParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	if commandID == "" {
		return nil, errors.New("missing required commandID parameter")
	}
	opts = slices.Concat(r.Options, opts)
	stdoutOffset := sandboxFieldValue(body.StdoutOffset, int64(0))
	stderrOffset := sandboxFieldValue(body.StderrOffset, int64(0))

	ws, err := dialSandboxCommandWebSocket(ctx, dataplaneURL, opts...)
	if err != nil {
		return nil, err
	}
	payload := sandboxCommandReconnectRequest{
		Type:         "reconnect",
		CommandID:    commandID,
		StdoutOffset: stdoutOffset,
		StderrOffset: stderrOffset,
	}
	if err := websocket.JSON.Send(ws, payload); err != nil {
		_ = ws.Close()
		return nil, &SandboxConnectionError{Message: fmt.Sprintf("langsmith: failed to send sandbox command reconnect request: %v", err)}
	}

	handle := newSandboxCommandHandle(ws, dataplaneURL, opts, commandID, 0, stdoutOffset, stderrOffset)
	handle.start()
	return handle, nil
}

type sandboxCommandStartRequest struct {
	Type               param.Field[string]            `json:"type" api:"required"`
	Command            param.Field[string]            `json:"command" api:"required"`
	TimeoutSeconds     param.Field[int64]             `json:"timeout_seconds"`
	Env                param.Field[map[string]string] `json:"env"`
	CWD                param.Field[string]            `json:"cwd"`
	Shell              param.Field[string]            `json:"shell"`
	IdleTimeoutSeconds param.Field[int64]             `json:"idle_timeout_seconds"`
	KillOnDisconnect   param.Field[bool]              `json:"kill_on_disconnect"`
	TTLSeconds         param.Field[int64]             `json:"ttl_seconds"`
	Pty                param.Field[bool]              `json:"pty"`
	SSHAgentForward    param.Field[bool]              `json:"ssh_agent_forward"`
}

func (r sandboxCommandStartRequest) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type sandboxCommandReconnectRequest struct {
	Type         string `json:"type"`
	CommandID    string `json:"command_id"`
	StdoutOffset int64  `json:"stdout_offset"`
	StderrOffset int64  `json:"stderr_offset"`
}

type sandboxWSMessage struct {
	Type      string `json:"type"`
	CommandID string `json:"command_id"`
	PID       int64  `json:"pid"`
	Stream    string `json:"stream"`
	Data      string `json:"data"`
	ChannelID string `json:"channel_id"`
	Offset    int64  `json:"offset"`
	ExitCode  int64  `json:"exit_code"`
	Error     string `json:"error"`
	ErrorType string `json:"error_type"`
}

func normalizeSandboxRunParams(body SandboxBoxRunParams) (SandboxBoxRunParams, int64, error) {
	command, ok := sandboxRequiredString(body.Command)
	if !ok {
		return SandboxBoxRunParams{}, 0, errors.New("missing required command parameter")
	}
	timeout := sandboxFieldValue(body.Timeout, defaultSandboxCommandTimeoutSeconds)
	shell := sandboxFieldValue(body.Shell, defaultSandboxCommandShell)
	body.Command = F(command)
	body.Timeout = F(timeout)
	body.Shell = F(shell)
	return body, timeout, nil
}

func normalizeSandboxCommandStartParams(body SandboxCommandStartParams) (sandboxCommandStartRequest, error) {
	command, ok := sandboxRequiredString(body.Command)
	pty := sandboxFieldValue(body.Pty, false)
	if !ok && !pty {
		return sandboxCommandStartRequest{}, errors.New("missing required command parameter")
	}
	out := sandboxCommandStartRequest{
		Type:               F("execute"),
		TimeoutSeconds:     F(sandboxFieldValue(body.TimeoutSeconds, defaultSandboxCommandTimeoutSeconds)),
		Env:                body.Env,
		CWD:                body.CWD,
		Shell:              F(sandboxFieldValue(body.Shell, defaultSandboxCommandShell)),
		IdleTimeoutSeconds: F(sandboxFieldValue(body.IdleTimeoutSeconds, defaultSandboxCommandIdleTimeout)),
		KillOnDisconnect:   F(sandboxFieldValue(body.KillOnDisconnect, false)),
		TTLSeconds:         F(sandboxFieldValue(body.TTLSeconds, defaultSandboxCommandTTLSeconds)),
		Pty:                body.Pty,
		SSHAgentForward:    body.SSHAgentForward,
	}
	if ok {
		out.Command = F(command)
	}
	return out, nil
}
