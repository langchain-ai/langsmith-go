package transport

import (
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/gorilla/websocket"
)

func TestWebSocketConnReadsAcrossBinaryMessages(t *testing.T) {
	t.Parallel()

	client, server := newWebSocketPair(t)
	conn := NewWebSocketConn(client)
	t.Cleanup(func() { _ = conn.Close() })

	for _, message := range [][]byte{nil, []byte("abc"), []byte("defg")} {
		if err := server.WriteMessage(websocket.BinaryMessage, message); err != nil {
			t.Fatalf("WriteMessage() error = %v", err)
		}
	}

	got := make([]byte, 7)
	if _, err := io.ReadFull(conn, got); err != nil {
		t.Fatalf("ReadFull() error = %v", err)
	}
	if string(got) != "abcdefg" {
		t.Fatalf("ReadFull() = %q, want %q", got, "abcdefg")
	}
}

func TestWebSocketConnRejectsTextMessages(t *testing.T) {
	t.Parallel()

	client, server := newWebSocketPair(t)
	conn := NewWebSocketConn(client)
	t.Cleanup(func() { _ = conn.Close() })

	if err := server.WriteMessage(websocket.TextMessage, []byte("not carrier bytes")); err != nil {
		t.Fatalf("WriteMessage() error = %v", err)
	}
	buffer := make([]byte, 1)
	_, err := conn.Read(buffer)
	if !errors.Is(err, ErrUnexpectedWebSocketMessageType) {
		t.Fatalf("Read() error = %v, want %v", err, ErrUnexpectedWebSocketMessageType)
	}
}

func TestWebSocketConnSerializesConcurrentWrites(t *testing.T) {
	t.Parallel()

	client, server := newWebSocketPair(t)
	conn := NewWebSocketConn(client)
	t.Cleanup(func() { _ = conn.Close() })

	const messageCount = 64
	errorsCh := make(chan error, messageCount)
	var writers sync.WaitGroup
	writers.Add(messageCount)
	for i := range messageCount {
		go func() {
			defer writers.Done()
			payload := []byte(fmt.Sprintf("message-%02d", i))
			written, err := conn.Write(payload)
			if err == nil && written != len(payload) {
				err = fmt.Errorf("Write() wrote %d bytes, want %d", written, len(payload))
			}
			errorsCh <- err
		}()
	}

	received := make(map[string]bool, messageCount)
	for range messageCount {
		messageType, payload, err := server.ReadMessage()
		if err != nil {
			t.Fatalf("ReadMessage() error = %v", err)
		}
		if messageType != websocket.BinaryMessage {
			t.Fatalf("ReadMessage() type = %d, want binary", messageType)
		}
		key := string(payload)
		if received[key] {
			t.Fatalf("received duplicate payload %q", key)
		}
		received[key] = true
	}

	writers.Wait()
	close(errorsCh)
	for err := range errorsCh {
		if err != nil {
			t.Fatalf("Write() error = %v", err)
		}
	}
	if len(received) != messageCount {
		t.Fatalf("received %d messages, want %d", len(received), messageCount)
	}
}

func TestWebSocketConnEnforcesMessageSizeBounds(t *testing.T) {
	t.Parallel()

	t.Run("outbound", func(t *testing.T) {
		t.Parallel()

		client, _ := newWebSocketPair(t)
		conn := NewWebSocketConn(client)
		t.Cleanup(func() { _ = conn.Close() })

		written, err := conn.Write(make([]byte, MaxWebSocketMessageSizeBytes+1))
		if !errors.Is(err, ErrWebSocketMessageTooLarge) {
			t.Fatalf("Write() error = %v, want %v", err, ErrWebSocketMessageTooLarge)
		}
		if written != 0 {
			t.Fatalf("Write() wrote %d bytes, want 0", written)
		}
	})

	t.Run("inbound", func(t *testing.T) {
		t.Parallel()

		client, server := newWebSocketPair(t)
		conn := NewWebSocketConn(client)
		t.Cleanup(func() { _ = conn.Close() })

		if err := server.WriteMessage(websocket.BinaryMessage, make([]byte, MaxWebSocketMessageSizeBytes+1)); err != nil {
			t.Fatalf("WriteMessage() error = %v", err)
		}
		_, err := io.ReadAll(conn)
		if !errors.Is(err, websocket.ErrReadLimit) {
			t.Fatalf("ReadAll() error = %v, want %v", err, websocket.ErrReadLimit)
		}
	})
}

func TestWebSocketConnCloseUnblocksRead(t *testing.T) {
	t.Parallel()

	client, _ := newWebSocketPair(t)
	conn := NewWebSocketConn(client)

	started := make(chan struct{})
	result := make(chan error, 1)
	go func() {
		close(started)
		buffer := make([]byte, 1)
		_, err := conn.Read(buffer)
		result <- err
	}()
	<-started

	if err := conn.Close(); err != nil {
		t.Fatalf("Close() error = %v", err)
	}
	select {
	case err := <-result:
		if err == nil {
			t.Fatal("Read() returned nil error after Close()")
		}
	case <-time.After(time.Second):
		t.Fatal("Close() did not unblock Read()")
	}

	if err := conn.Close(); err != nil {
		t.Fatalf("second Close() error = %v", err)
	}
}

func newWebSocketPair(t *testing.T) (*websocket.Conn, *websocket.Conn) {
	t.Helper()

	type upgradeResult struct {
		conn *websocket.Conn
		err  error
	}
	resultCh := make(chan upgradeResult, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin: func(*http.Request) bool { return true },
	}
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		resultCh <- upgradeResult{conn: conn, err: err}
	}))

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http")
	client, response, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
	if err != nil {
		server.Close()
		t.Fatalf("Dial() error = %v", err)
	}
	result := <-resultCh
	if result.err != nil {
		_ = client.Close()
		server.Close()
		t.Fatalf("Upgrade() error = %v", result.err)
	}

	t.Cleanup(func() {
		_ = client.Close()
		_ = result.conn.Close()
		server.Close()
	})
	return client, result.conn
}
