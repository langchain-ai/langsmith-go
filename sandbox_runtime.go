package langsmith

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
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
	OnStdout func(string)
	OnStderr func(string)
}

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

// SandboxCommandHandle is a handle to a running sandbox command.
type SandboxCommandHandle struct {
	CommandID string
	PID       int64

	dataplaneURL string
	opts         []option.RequestOption

	ws        *websocket.Conn
	sendMu    sync.Mutex
	stateMu   sync.Mutex
	chunks    chan SandboxOutputChunk
	done      chan struct{}
	result    *SandboxExecutionResult
	err       error
	killed    bool
	closed    bool
	callbacks SandboxCommandCallbacks
	stdout    strings.Builder
	stderr    strings.Builder
	stdoutOf  int64
	stderrOf  int64
}

func newSandboxCommandHandle(ws *websocket.Conn, dataplaneURL string, opts []option.RequestOption, commandID string, pid int64, stdoutOffset int64, stderrOffset int64) *SandboxCommandHandle {
	return &SandboxCommandHandle{
		CommandID:    commandID,
		PID:          pid,
		dataplaneURL: dataplaneURL,
		opts:         append([]option.RequestOption{}, opts...),
		ws:           ws,
		chunks:       make(chan SandboxOutputChunk, 64),
		done:         make(chan struct{}),
		stdoutOf:     stdoutOffset,
		stderrOf:     stderrOffset,
	}
}

func (h *SandboxCommandHandle) start() {
	go h.readLoop()
}

// Next returns the next stdout/stderr chunk. If ok is false, the command stream
// has ended and Result returns the final command result.
func (h *SandboxCommandHandle) Next(ctx context.Context) (chunk SandboxOutputChunk, ok bool, err error) {
	select {
	case chunk, ok = <-h.chunks:
		if !ok {
			return SandboxOutputChunk{}, false, h.Err()
		}
		return chunk, true, nil
	case <-ctx.Done():
		return SandboxOutputChunk{}, false, ctx.Err()
	}
}

// Result waits for the command to exit and returns its final result. If output
// chunks have not been consumed, Result drains them before returning.
func (h *SandboxCommandHandle) Result(ctx context.Context) (*SandboxExecutionResult, error) {
	for {
		select {
		case _, ok := <-h.chunks:
			if !ok {
				h.stateMu.Lock()
				defer h.stateMu.Unlock()
				if h.err != nil {
					return nil, h.err
				}
				if h.result == nil {
					return nil, &SandboxOperationError{Operation: "command", Message: "sandbox command stream ended without exit message"}
				}
				return h.result, nil
			}
		case <-ctx.Done():
			return nil, ctx.Err()
		}
	}
}

// Err returns the terminal stream error, if any.
func (h *SandboxCommandHandle) Err() error {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	return h.err
}

// Done is closed when the command stream exits or fails.
func (h *SandboxCommandHandle) Done() <-chan struct{} {
	return h.done
}

// Close closes the command stream connection. The command may continue running
// on the sandbox unless it was started with KillOnDisconnect.
func (h *SandboxCommandHandle) Close() error {
	h.stateMu.Lock()
	h.closed = true
	h.stateMu.Unlock()
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	return h.ws.Close()
}

// LastStdoutOffset returns the last known stdout byte offset.
func (h *SandboxCommandHandle) LastStdoutOffset() int64 {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	return h.stdoutOf
}

// LastStderrOffset returns the last known stderr byte offset.
func (h *SandboxCommandHandle) LastStderrOffset() int64 {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	return h.stderrOf
}

// SendInput writes data to the running command's stdin.
func (h *SandboxCommandHandle) SendInput(data string) error {
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	return websocket.JSON.Send(h.ws, map[string]string{"type": "input", "data": data})
}

// Kill sends a kill request for the running command.
func (h *SandboxCommandHandle) Kill() error {
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	h.stateMu.Lock()
	h.killed = true
	h.stateMu.Unlock()
	return websocket.JSON.Send(h.ws, map[string]string{"type": "kill"})
}

// Reconnect opens a new stream for this command from the last known offsets.
func (h *SandboxCommandHandle) Reconnect(ctx context.Context) (*SandboxCommandHandle, error) {
	if h.CommandID == "" {
		return nil, &SandboxOperationError{Operation: "reconnect", Message: "cannot reconnect: command ID is not available"}
	}
	ws, err := dialSandboxCommandWebSocket(ctx, h.dataplaneURL, h.opts...)
	if err != nil {
		return nil, err
	}
	payload := sandboxCommandReconnectRequest{
		Type:         "reconnect",
		CommandID:    h.CommandID,
		StdoutOffset: h.LastStdoutOffset(),
		StderrOffset: h.LastStderrOffset(),
	}
	if err := websocket.JSON.Send(ws, payload); err != nil {
		_ = ws.Close()
		return nil, &SandboxConnectionError{Message: fmt.Sprintf("langsmith: failed to send sandbox command reconnect request: %v", err)}
	}
	reconnected := newSandboxCommandHandle(ws, h.dataplaneURL, h.opts, h.CommandID, h.PID, payload.StdoutOffset, payload.StderrOffset)
	reconnected.callbacks = h.callbacks
	reconnected.start()
	return reconnected, nil
}

func (h *SandboxCommandHandle) readLoop() {
	defer close(h.done)
	defer close(h.chunks)
	defer h.closeCurrentWebSocket()

	reconnectAttempts := 0
	for {
		msg, err := receiveSandboxWSMessage(context.Background(), h.ws)
		if err != nil {
			if !h.shouldReconnect(err, reconnectAttempts) {
				h.setErr(err)
				return
			}
			reconnectAttempts++
			if !h.reconnectForReadLoop(reconnectAttempts) {
				h.setErr(&SandboxConnectionError{Message: fmt.Sprintf("langsmith: lost sandbox command connection %d times in succession, giving up", reconnectAttempts)})
				return
			}
			continue
		}

		switch msg.Type {
		case "stdout", "stderr":
			reconnectAttempts = 0
			chunk := SandboxOutputChunk{Stream: msg.Type, Data: msg.Data, Offset: msg.Offset}
			h.appendChunk(chunk)
			h.invokeCallback(chunk)
			h.chunks <- chunk
		case "exit":
			h.setResult(msg.ExitCode)
			return
		case "error":
			h.setErr(sandboxErrorFromWSMessage(msg, h.CommandID))
			return
		case "started":
			if h.CommandID == "" {
				h.CommandID = msg.CommandID
			}
			if h.PID == 0 {
				h.PID = msg.PID
			}
		default:
			h.setErr(&SandboxOperationError{Operation: "command", Message: fmt.Sprintf("unknown sandbox command message type %q", msg.Type)})
			return
		}
	}
}

func (h *SandboxCommandHandle) shouldReconnect(err error, attempts int) bool {
	var connErr *SandboxConnectionError
	if !errors.As(err, &connErr) {
		return false
	}
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	if h.killed || h.closed || h.CommandID == "" {
		return false
	}
	return attempts < sandboxCommandMaxAutoReconnects
}

func (h *SandboxCommandHandle) reconnectForReadLoop(attempt int) bool {
	delay := sandboxCommandReconnectBackoffBase << max(0, attempt-1)
	if delay > sandboxCommandReconnectBackoffMax {
		delay = sandboxCommandReconnectBackoffMax
	}
	time.Sleep(delay)
	if h.isClosedOrKilled() {
		return false
	}

	ws, err := dialSandboxCommandWebSocket(context.Background(), h.dataplaneURL, h.opts...)
	if err != nil {
		return false
	}
	payload := sandboxCommandReconnectRequest{
		Type:         "reconnect",
		CommandID:    h.CommandID,
		StdoutOffset: h.LastStdoutOffset(),
		StderrOffset: h.LastStderrOffset(),
	}
	if err := websocket.JSON.Send(ws, payload); err != nil {
		_ = ws.Close()
		return false
	}

	h.sendMu.Lock()
	old := h.ws
	h.ws = ws
	h.sendMu.Unlock()
	_ = old.Close()
	return true
}

func (h *SandboxCommandHandle) isClosedOrKilled() bool {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	return h.closed || h.killed
}

func (h *SandboxCommandHandle) closeCurrentWebSocket() {
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	_ = h.ws.Close()
}

func (h *SandboxCommandHandle) invokeCallback(chunk SandboxOutputChunk) {
	if chunk.Stream == "stdout" && h.callbacks.OnStdout != nil {
		h.callbacks.OnStdout(chunk.Data)
	}
	if chunk.Stream == "stderr" && h.callbacks.OnStderr != nil {
		h.callbacks.OnStderr(chunk.Data)
	}
}

func (h *SandboxCommandHandle) appendChunk(chunk SandboxOutputChunk) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	if chunk.Stream == "stdout" {
		h.stdout.WriteString(chunk.Data)
		h.stdoutOf = chunk.Offset + int64(len([]byte(chunk.Data)))
		return
	}
	h.stderr.WriteString(chunk.Data)
	h.stderrOf = chunk.Offset + int64(len([]byte(chunk.Data)))
}

func (h *SandboxCommandHandle) setResult(exitCode int64) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	h.result = &SandboxExecutionResult{
		Stdout:   h.stdout.String(),
		Stderr:   h.stderr.String(),
		ExitCode: exitCode,
	}
}

func (h *SandboxCommandHandle) setErr(err error) {
	h.stateMu.Lock()
	defer h.stateMu.Unlock()
	h.err = err
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
	if !ok {
		return sandboxCommandStartRequest{}, errors.New("missing required command parameter")
	}
	return sandboxCommandStartRequest{
		Type:               F("execute"),
		Command:            F(command),
		TimeoutSeconds:     F(sandboxFieldValue(body.TimeoutSeconds, defaultSandboxCommandTimeoutSeconds)),
		Env:                body.Env,
		CWD:                body.CWD,
		Shell:              F(sandboxFieldValue(body.Shell, defaultSandboxCommandShell)),
		IdleTimeoutSeconds: F(sandboxFieldValue(body.IdleTimeoutSeconds, defaultSandboxCommandIdleTimeout)),
		KillOnDisconnect:   F(sandboxFieldValue(body.KillOnDisconnect, false)),
		TTLSeconds:         F(sandboxFieldValue(body.TTLSeconds, defaultSandboxCommandTTLSeconds)),
		Pty:                body.Pty,
	}, nil
}

func sandboxRequiredString(field param.Field[string]) (string, bool) {
	if !field.Present || field.Null {
		return "", false
	}
	if field.Value != "" {
		return field.Value, true
	}
	if raw, ok := field.Raw.(string); ok && raw != "" {
		return raw, true
	}
	return "", false
}

func sandboxFieldValue[T any](field param.Field[T], fallback T) T {
	if field.Present && !field.Null {
		return field.Value
	}
	return fallback
}

func requireSandboxDataplaneURL(name string, status string, dataplaneURL string) (string, error) {
	if status != "" && status != "ready" {
		return "", &SandboxNotReadyError{SandboxName: name, Status: status}
	}
	if dataplaneURL == "" {
		return "", &SandboxDataplaneNotConfiguredError{SandboxName: name}
	}
	return dataplaneURL, nil
}

func sandboxDataplaneURL(dataplaneURL string, path string) (string, error) {
	u, err := url.Parse(dataplaneURL)
	if err != nil {
		return "", err
	}
	if u.Scheme == "" || u.Host == "" {
		return "", fmt.Errorf("invalid sandbox dataplane URL %q", dataplaneURL)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/" + strings.TrimLeft(path, "/")
	return u.String(), nil
}

func sandboxWebSocketURL(dataplaneURL string) (string, error) {
	u, err := url.Parse(dataplaneURL)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "http":
		u.Scheme = "ws"
	case "https":
		u.Scheme = "wss"
	case "ws", "wss":
	default:
		return "", fmt.Errorf("unsupported sandbox dataplane URL scheme %q", u.Scheme)
	}
	u.Path = strings.TrimRight(u.Path, "/") + "/execute/ws"
	u.RawQuery = ""
	return u.String(), nil
}

func dialSandboxCommandWebSocket(ctx context.Context, dataplaneURL string, opts ...option.RequestOption) (*websocket.Conn, error) {
	wsURL, err := sandboxWebSocketURL(dataplaneURL)
	if err != nil {
		return nil, err
	}
	return dialSandboxWebSocketURL(ctx, wsURL, opts...)
}

func dialSandboxWebSocketURL(ctx context.Context, wsURL string, opts ...option.RequestOption) (*websocket.Conn, error) {
	if err := ctx.Err(); err != nil {
		return nil, err
	}
	headers, err := sandboxHeaders(ctx, wsURL, opts...)
	if err != nil {
		return nil, err
	}
	origin, err := sandboxWebSocketOrigin(wsURL)
	if err != nil {
		return nil, err
	}
	config, err := websocket.NewConfig(wsURL, origin)
	if err != nil {
		return nil, err
	}
	config.Header = headers

	type dialResult struct {
		ws  *websocket.Conn
		err error
	}
	ch := make(chan dialResult, 1)
	go func() {
		ws, err := websocket.DialConfig(config)
		ch <- dialResult{ws: ws, err: err}
	}()

	select {
	case res := <-ch:
		if res.err != nil {
			return nil, &SandboxConnectionError{Message: fmt.Sprintf("langsmith: failed to connect to sandbox command WebSocket: %v", res.err)}
		}
		return res.ws, nil
	case <-ctx.Done():
		return nil, ctx.Err()
	}
}

func sandboxWebSocketOrigin(wsURL string) (string, error) {
	u, err := url.Parse(wsURL)
	if err != nil {
		return "", err
	}
	switch u.Scheme {
	case "ws":
		u.Scheme = "http"
	case "wss":
		u.Scheme = "https"
	}
	u.Path = "/"
	u.RawQuery = ""
	u.Fragment = ""
	return u.String(), nil
}

func sandboxHeaders(ctx context.Context, requestURL string, opts ...option.RequestOption) (http.Header, error) {
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, requestURL, nil, nil, opts...)
	if err != nil {
		return nil, err
	}
	return cfg.Request.Header.Clone(), nil
}

func receiveSandboxWSMessage(ctx context.Context, ws *websocket.Conn) (sandboxWSMessage, error) {
	type receiveResult struct {
		msg sandboxWSMessage
		err error
	}
	ch := make(chan receiveResult, 1)
	go func() {
		var raw string
		if err := websocket.Message.Receive(ws, &raw); err != nil {
			ch <- receiveResult{err: &SandboxConnectionError{Message: fmt.Sprintf("langsmith: sandbox command WebSocket closed unexpectedly: %v", err)}}
			return
		}
		var msg sandboxWSMessage
		if err := json.Unmarshal([]byte(raw), &msg); err != nil {
			ch <- receiveResult{err: &SandboxOperationError{Operation: "command", Message: fmt.Sprintf("failed to parse sandbox command message: %v", err)}}
			return
		}
		ch <- receiveResult{msg: msg}
	}()

	select {
	case res := <-ch:
		return res.msg, res.err
	case <-ctx.Done():
		_ = ws.Close()
		return sandboxWSMessage{}, ctx.Err()
	}
}

func sandboxErrorFromWSMessage(msg sandboxWSMessage, commandID string) error {
	errorType := msg.ErrorType
	if errorType == "" {
		errorType = "CommandError"
	}
	message := msg.Error
	if message == "" {
		message = "unknown sandbox command error"
	}
	if errorType == "CommandTimeout" {
		return &SandboxCommandTimeoutError{Message: message}
	}
	if errorType == "CommandNotFound" && commandID != "" {
		message = fmt.Sprintf("command not found: %s", commandID)
	}
	if errorType == "SessionExpired" && commandID != "" {
		message = fmt.Sprintf("session expired: %s", commandID)
	}
	operation := "command"
	if commandID != "" {
		operation = "reconnect"
	}
	return &SandboxOperationError{
		Operation: operation,
		ErrorType: errorType,
		Message:   message,
	}
}
