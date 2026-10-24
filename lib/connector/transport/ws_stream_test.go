package transport

import (
	"bytes"
	"io"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestWebSocketStreamHalfCloseKeepsTheOtherDirectionOpen(t *testing.T) {
	clientConn, serverConn := newWebSocketPair(t)
	client := NewWebSocketStream(clientConn)
	server := NewWebSocketStream(serverConn)

	request := []byte("GET / HTTP/1.0\r\n\r\n")
	_, err := client.Write(request)
	require.NoError(t, err)
	require.NoError(t, client.CloseWrite())

	received, err := io.ReadAll(server)
	require.NoError(t, err)
	assert.Equal(t, request, received)

	response := []byte("HTTP/1.0 200 OK\r\n\r\nhello")
	_, err = server.Write(response)
	require.NoError(t, err)
	require.NoError(t, server.CloseWrite())

	received, err = io.ReadAll(client)
	require.NoError(t, err)
	assert.Equal(t, response, received)
}

func TestWebSocketStreamSplitsWritesLargerThanOneMessage(t *testing.T) {
	clientConn, serverConn := newWebSocketPair(t)
	client := NewWebSocketStream(clientConn)
	server := NewWebSocketStream(serverConn)

	payload := bytes.Repeat([]byte("x"), 2*MaxWebSocketMessageSizeBytes+17)
	done := make(chan error, 1)
	go func() {
		_, err := client.Write(payload)
		if err == nil {
			err = client.CloseWrite()
		}
		done <- err
	}()

	received, err := io.ReadAll(server)
	require.NoError(t, err)
	require.NoError(t, <-done)
	assert.Equal(t, len(payload), len(received))
}
