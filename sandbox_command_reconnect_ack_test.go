package langsmith

import (
	"context"
	"fmt"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"testing/synctest"

	"github.com/langchain-ai/langsmith-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

// ackingServer answers each attachment with a "started" frame naming commandID(n)
// (empty means send nothing), then hands the connection to onAttached.
//
// The server must be started outside any synctest bubble: its Accept loop blocks
// on real I/O, which never counts as durably blocked, so inside a bubble it would
// stall the fake clock and the test would hang until its deadline.
func ackingServer(t *testing.T, commandID func(n int64) string, onAttached func(ws *websocket.Conn, n int64)) (string, *atomic.Int64) {
	t.Helper()
	var attachments atomic.Int64
	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		n := attachments.Add(1)
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			return
		}
		if id := commandID(n); id != "" {
			if err := websocket.Message.Send(ws, fmt.Sprintf(`{"type":"started","command_id":%q,"pid":42}`, id)); err != nil {
				return
			}
		}
		if onAttached != nil {
			onAttached(ws, n)
		}
	}))
	t.Cleanup(srv.Close)
	return srv.URL, &attachments
}

func startSilentCommand(t *testing.T, dataplaneURL string) *SandboxCommandHandle {
	t.Helper()
	client := NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	handle, err := client.Sandboxes.Boxes.StartCommandWithDataplaneURL(context.Background(), dataplaneURL, SandboxCommandStartParams{
		Command:            String("sleep 60"),
		IdleTimeoutSeconds: Int(-1),
	})
	require.NoError(t, err)
	return handle
}

// A command and its WebSocket are separate things: the command keeps running on
// the server and the socket is only this client's attachment to it. The server
// acknowledges every successful reattachment with a "started" frame, which for a
// command producing no output is the only evidence the reattachment landed.
//
// The reconnect budget is meant to bound *consecutive failed* reattachments.
// Reset it only on stdout/stderr and it instead counts socket closures since the
// last output, so a healthy, attached, quiet command dies one loss at a time.
func TestSandboxCommandSilentCommandSurvivesMoreLossesThanBudget(t *testing.T) {
	losses := int64(sandboxCommandMaxAutoReconnects + 3)
	url, attachments := ackingServer(t,
		func(int64) string { return "cmd-123" },
		func(ws *websocket.Conn, n int64) {
			if n > losses {
				_ = websocket.Message.Send(ws, `{"type":"exit","exit_code":7}`)
				return
			}
			// Otherwise drop the socket under the client, with no exit frame.
		})

	// Bubbled so the reconnect backoff — 500ms doubling to an 8s cap, ~40s over
	// this many losses — is virtual rather than slept through.
	synctest.Test(t, func(t *testing.T) {
		handle := startSilentCommand(t, url)
		result, err := handle.Result(context.Background())
		require.NoError(t, err, "every reattachment was acknowledged, so the budget must not run out")
		assert.Equal(t, int64(7), result.ExitCode)
		assert.Empty(t, result.Stdout)
		assert.Equal(t, losses+1, attachments.Load())
	})
}

// A reattachment that is never acknowledged still has to hit the limit.
func TestSandboxCommandUnacknowledgedReattachmentsExhaustBudget(t *testing.T) {
	url, attachments := ackingServer(t,
		// Only the first attachment is acknowledged; reattachments die silently.
		func(n int64) string {
			if n == 1 {
				return "cmd-123"
			}
			return ""
		}, nil)

	synctest.Test(t, func(t *testing.T) {
		handle := startSilentCommand(t, url)
		_, err := handle.Result(context.Background())
		var connErr *SandboxConnectionError
		require.ErrorAs(t, err, &connErr)
		// One initial attachment plus the budget's worth of reattachments.
		assert.Equal(t, int64(sandboxCommandMaxAutoReconnects+1), attachments.Load())
	})
}

// An acknowledgement naming a different command is not proof of anything.
func TestSandboxCommandForeignAcknowledgementDoesNotResetBudget(t *testing.T) {
	url, attachments := ackingServer(t,
		func(n int64) string {
			if n == 1 {
				return "cmd-123"
			}
			return "someone-elses-cmd"
		}, nil)

	synctest.Test(t, func(t *testing.T) {
		handle := startSilentCommand(t, url)
		_, err := handle.Result(context.Background())
		var connErr *SandboxConnectionError
		require.ErrorAs(t, err, &connErr)
		assert.Equal(t, int64(sandboxCommandMaxAutoReconnects+1), attachments.Load())
	})
}
