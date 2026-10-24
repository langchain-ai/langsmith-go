package transport

import (
	"context"
	"errors"
	"io"
	"net"
	"testing"
	"time"

	"github.com/hashicorp/yamux"
)

func TestBridgePreservesHalfCloseAcrossYamuxAndTCP(t *testing.T) {
	t.Parallel()

	bridgeStream, callerStream := newYamuxStreamPair(t)
	bridgeTCP, targetTCP := newTCPPair(t)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	bridgeResult := make(chan error, 1)
	go func() {
		bridgeResult <- Bridge(ctx, bridgeStream, bridgeTCP)
	}()

	request := []byte("request that requires an EOF")
	if _, err := callerStream.Write(request); err != nil {
		t.Fatalf("caller Write() error = %v", err)
	}
	if err := callerStream.Close(); err != nil {
		t.Fatalf("caller Close() error = %v", err)
	}

	targetResult := make(chan []byte, 1)
	go func() {
		got, err := io.ReadAll(targetTCP)
		if err != nil {
			targetResult <- nil
			return
		}
		targetResult <- got
		_, _ = targetTCP.Write([]byte("response after EOF"))
		_ = targetTCP.CloseWrite()
	}()

	response, err := io.ReadAll(callerStream)
	if err != nil {
		t.Fatalf("caller ReadAll() error = %v", err)
	}
	if string(response) != "response after EOF" {
		t.Fatalf("caller response = %q, want %q", response, "response after EOF")
	}
	if got := <-targetResult; string(got) != string(request) {
		t.Fatalf("target request = %q, want %q", got, request)
	}

	select {
	case err := <-bridgeResult:
		if err != nil {
			t.Fatalf("Bridge() error = %v", err)
		}
	case <-ctx.Done():
		t.Fatalf("Bridge() did not finish: %v", ctx.Err())
	}
}

func TestBridgeCancellationUnblocksCopies(t *testing.T) {
	t.Parallel()

	leftBridge, leftPeer := net.Pipe()
	rightBridge, rightPeer := net.Pipe()
	t.Cleanup(func() {
		_ = leftPeer.Close()
		_ = rightPeer.Close()
	})

	ctx, cancel := context.WithCancel(context.Background())
	result := make(chan error, 1)
	go func() {
		result <- Bridge(ctx, leftBridge, rightBridge)
	}()
	cancel()

	select {
	case err := <-result:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("Bridge() error = %v, want %v", err, context.Canceled)
		}
	case <-time.After(time.Second):
		t.Fatal("Bridge() did not stop after cancellation")
	}
}

func TestBridgePreservesCopyErrorIfContextIsCanceledDuringCleanup(t *testing.T) {
	t.Parallel()

	leftRaw, leftPeer := net.Pipe()
	rightRaw, rightPeer := net.Pipe()
	t.Cleanup(func() {
		_ = leftPeer.Close()
		_ = rightPeer.Close()
	})

	wantErr := errors.New("copy failed")
	left := &readErrorConn{Conn: leftRaw, err: wantErr}
	ctx, cancel := context.WithCancel(context.Background())
	right := &cancelOnCloseSlowReadConn{
		Conn:   rightRaw,
		cancel: cancel,
		delay:  50 * time.Millisecond,
	}

	err := Bridge(ctx, left, right)
	if !errors.Is(err, wantErr) {
		t.Fatalf("Bridge() error = %v, want %v", err, wantErr)
	}
}

func newYamuxStreamPair(t *testing.T) (*yamux.Stream, *yamux.Stream) {
	t.Helper()

	clientWire, serverWire := net.Pipe()
	clientSession, err := yamux.Client(clientWire, YamuxConfig())
	if err != nil {
		t.Fatalf("yamux.Client() error = %v", err)
	}
	serverSession, err := yamux.Server(serverWire, YamuxConfig())
	if err != nil {
		_ = clientSession.Close()
		t.Fatalf("yamux.Server() error = %v", err)
	}

	clientResult := make(chan struct {
		stream *yamux.Stream
		err    error
	}, 1)
	go func() {
		stream, openErr := clientSession.OpenStream()
		clientResult <- struct {
			stream *yamux.Stream
			err    error
		}{stream: stream, err: openErr}
	}()
	serverStream, err := serverSession.AcceptStream()
	if err != nil {
		_ = clientSession.Close()
		_ = serverSession.Close()
		t.Fatalf("AcceptStream() error = %v", err)
	}
	client := <-clientResult
	if client.err != nil {
		_ = serverStream.Close()
		_ = clientSession.Close()
		_ = serverSession.Close()
		t.Fatalf("OpenStream() error = %v", client.err)
	}

	t.Cleanup(func() {
		_ = client.stream.Close()
		_ = serverStream.Close()
		_ = clientSession.Close()
		_ = serverSession.Close()
		_ = clientWire.Close()
		_ = serverWire.Close()
	})
	return serverStream, client.stream
}

func newTCPPair(t *testing.T) (*net.TCPConn, *net.TCPConn) {
	t.Helper()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("127.0.0.1")})
	if err != nil {
		t.Fatalf("ListenTCP() error = %v", err)
	}
	accepted := make(chan struct {
		conn *net.TCPConn
		err  error
	}, 1)
	go func() {
		conn, acceptErr := listener.AcceptTCP()
		accepted <- struct {
			conn *net.TCPConn
			err  error
		}{conn: conn, err: acceptErr}
	}()

	client, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
	if err != nil {
		_ = listener.Close()
		t.Fatalf("DialTCP() error = %v", err)
	}
	server := <-accepted
	_ = listener.Close()
	if server.err != nil {
		_ = client.Close()
		t.Fatalf("AcceptTCP() error = %v", server.err)
	}

	t.Cleanup(func() {
		_ = client.Close()
		_ = server.conn.Close()
	})
	return client, server.conn
}

type readErrorConn struct {
	net.Conn
	err error
}

func (c *readErrorConn) Read([]byte) (int, error) {
	return 0, c.err
}

type cancelOnCloseSlowReadConn struct {
	net.Conn
	cancel context.CancelFunc
	delay  time.Duration
}

func (c *cancelOnCloseSlowReadConn) Read(buffer []byte) (int, error) {
	n, err := c.Conn.Read(buffer)
	if err != nil {
		time.Sleep(c.delay)
	}
	return n, err
}

func (c *cancelOnCloseSlowReadConn) Close() error {
	c.cancel()
	return c.Conn.Close()
}
