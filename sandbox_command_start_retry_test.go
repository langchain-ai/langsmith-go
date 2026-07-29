package langsmith_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/net/websocket"
)

func retryTestClient() *langsmith.Client {
	return langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
}

func TestSandboxCommandStartRetriesStreamEndedBeforeStarted(t *testing.T) {
	var attempts atomic.Int32
	commandIDs := make(chan string, 4)

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			return
		}
		commandID, _ := msg["command_id"].(string)
		commandIDs <- commandID
		// First attempt closes before "started"; the second one succeeds.
		if attempts.Add(1) == 1 {
			_ = ws.Close()
			return
		}
		_ = websocket.Message.Send(ws, `{"type":"started","command_id":"`+commandID+`","pid":7}`)
		_ = websocket.Message.Send(ws, `{"type":"exit","exit_code":0}`)
	}))
	defer srv.Close()

	handle, err := retryTestClient().Sandboxes.Boxes.StartCommandWithDataplaneURL(
		context.Background(),
		srv.URL,
		langsmith.SandboxCommandStartParams{Command: langsmith.String("echo hi")},
	)
	require.NoError(t, err)
	require.NotNil(t, handle)

	require.EqualValues(t, 2, attempts.Load())
	first := <-commandIDs
	second := <-commandIDs
	require.NotEmpty(t, first, "client must supply a command_id for idempotent re-issue")
	assert.Equal(t, first, second, "retry must reuse the command_id so the server dedupes")
	assert.Equal(t, first, handle.CommandID)
}

func TestSandboxCommandStartDoesNotRetryRejectedHandshake(t *testing.T) {
	var attempts atomic.Int32

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		attempts.Add(1)
		w.WriteHeader(http.StatusNotFound)
	}))
	defer srv.Close()

	_, err := retryTestClient().Sandboxes.Boxes.StartCommandWithDataplaneURL(
		context.Background(),
		srv.URL,
		langsmith.SandboxCommandStartParams{Command: langsmith.String("echo hi")},
	)
	require.Error(t, err)

	var permanent *langsmith.SandboxConnectionError
	require.ErrorAs(t, err, &permanent)
	var retryable *langsmith.SandboxConnectTimeoutError
	assert.NotErrorAs(t, err, &retryable, "a rejected upgrade must not be classified retryable")
	assert.EqualValues(t, 1, attempts.Load(), "a rejected upgrade must not be retried")
}

func TestSandboxCommandStartHonorsContextCancellation(t *testing.T) {
	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		var msg map[string]any
		_ = websocket.JSON.Receive(ws, &msg)
		_ = ws.Close() // always retryable, so the loop keeps backing off
	}))
	defer srv.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()

	_, err := retryTestClient().Sandboxes.Boxes.StartCommandWithDataplaneURL(
		ctx,
		srv.URL,
		langsmith.SandboxCommandStartParams{Command: langsmith.String("echo hi")},
	)
	// Deterministic timing for this path lives in the synctest suite.
	require.ErrorIs(t, err, context.DeadlineExceeded)
}
