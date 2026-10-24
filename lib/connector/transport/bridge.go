package transport

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"time"

	"github.com/hashicorp/yamux"
)

type closeWriter interface {
	CloseWrite() error
}

type bridgeCopyResult struct {
	direction string
	err       error
}

// Bridge copies opaque bytes in both directions. A clean EOF half-closes the
// destination while leaving the response direction alive.
func Bridge(ctx context.Context, left, right net.Conn) error {
	defer closeBridgeConnections(left, right)

	results := make(chan bridgeCopyResult, 2)
	go bridgeCopy("left to right", right, left, results)
	go bridgeCopy("right to left", left, right, results)

	remaining := 2
	ctxDone := ctx.Done()
	var bridgeErr error
	for remaining > 0 {
		select {
		case <-ctxDone:
			if bridgeErr == nil {
				bridgeErr = ctx.Err()
			}
			interruptBridgeConnections(left, right)
			ctxDone = nil
		case result := <-results:
			remaining--
			if err := materialBridgeError(result.err); err != nil && bridgeErr == nil {
				bridgeErr = fmt.Errorf("bridge %s: %w", result.direction, err)
				interruptBridgeConnections(left, right)
			}
		}
	}
	return bridgeErr
}

func bridgeCopy(
	direction string,
	destination net.Conn,
	source net.Conn,
	results chan<- bridgeCopyResult,
) {
	_, err := io.Copy(destination, source)
	if err == nil {
		err = halfCloseWrite(destination)
	}
	results <- bridgeCopyResult{direction: direction, err: err}
}

func halfCloseWrite(conn net.Conn) error {
	if conn, ok := conn.(closeWriter); ok {
		return conn.CloseWrite()
	}
	if stream, ok := conn.(*yamux.Stream); ok {
		return stream.Close()
	}
	return nil
}

func interruptBridgeConnections(connections ...net.Conn) {
	deadline := time.Now()
	for _, conn := range connections {
		_ = conn.SetDeadline(deadline)
		_ = conn.Close()
	}
}

func closeBridgeConnections(connections ...net.Conn) {
	for _, conn := range connections {
		_ = conn.Close()
	}
}

func materialBridgeError(err error) error {
	if err == nil ||
		errors.Is(err, net.ErrClosed) ||
		errors.Is(err, io.ErrClosedPipe) {
		return nil
	}
	return err
}
