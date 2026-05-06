package langsmith

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"strings"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
	"golang.org/x/net/websocket"
)

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

// Resize updates the command PTY size.
func (h *SandboxCommandHandle) Resize(cols int, rows int) error {
	if cols <= 0 || rows <= 0 {
		return errors.New("cols and rows must be greater than zero")
	}
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	return websocket.JSON.Send(h.ws, map[string]any{"type": "resize", "cols": cols, "rows": rows})
}

// SendSSHAgentData sends SSH agent response bytes for a forwarded channel.
func (h *SandboxCommandHandle) SendSSHAgentData(channelID string, data []byte) error {
	if channelID == "" {
		return errors.New("missing SSH agent channel ID")
	}
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	return websocket.JSON.Send(h.ws, map[string]string{
		"type":       "ssh_agent_data",
		"channel_id": channelID,
		"data":       base64.StdEncoding.EncodeToString(data),
	})
}

// CloseSSHAgentChannel tells the sandbox that a forwarded SSH agent channel is closed.
func (h *SandboxCommandHandle) CloseSSHAgentChannel(channelID string) error {
	if channelID == "" {
		return errors.New("missing SSH agent channel ID")
	}
	h.sendMu.Lock()
	defer h.sendMu.Unlock()
	return websocket.JSON.Send(h.ws, map[string]string{"type": "ssh_agent_close", "channel_id": channelID})
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
		case "ssh_agent_data":
			h.invokeSSHAgentData(msg)
		case "ssh_agent_close":
			h.invokeSSHAgentClose(msg.ChannelID)
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

func (h *SandboxCommandHandle) invokeSSHAgentData(msg sandboxWSMessage) {
	if h.callbacks.OnSSHAgentData == nil || msg.ChannelID == "" || msg.Data == "" {
		return
	}
	data, err := base64.StdEncoding.DecodeString(msg.Data)
	if err != nil {
		return
	}
	h.callbacks.OnSSHAgentData(msg.ChannelID, data)
}

func (h *SandboxCommandHandle) invokeSSHAgentClose(channelID string) {
	if h.callbacks.OnSSHAgentClose != nil && channelID != "" {
		h.callbacks.OnSSHAgentClose(channelID)
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
