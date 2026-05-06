package langsmith

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
)

// BridgeSandboxTunnel copies bytes bidirectionally between two read/write
// closers until either side closes or ctx is cancelled. Both sides are closed
// before it returns.
func BridgeSandboxTunnel(ctx context.Context, a io.ReadWriteCloser, b io.ReadWriteCloser) {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(2)
	go func() {
		defer wg.Done()
		defer cancel()
		copySandboxTunnelUntilDone(ctx, b, a, a)
	}()
	go func() {
		defer wg.Done()
		defer cancel()
		copySandboxTunnelUntilDone(ctx, a, b, b)
	}()
	wg.Wait()
	_ = a.Close()
	_ = b.Close()
}

// BridgeSandboxTunnelIO bridges input/output streams to a tunnel stream. It is
// intended for stdio-style integrations.
func BridgeSandboxTunnelIO(ctx context.Context, stream io.ReadWriteCloser, input io.Reader, output io.Writer) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	errs := make(chan error, 2)
	go func() {
		_, err := io.Copy(stream, input)
		errs <- err
		cancel()
		_ = stream.Close()
	}()
	go func() {
		_, err := io.Copy(output, stream)
		errs <- err
		cancel()
	}()

	select {
	case <-ctx.Done():
		_ = stream.Close()
		return nil
	case err := <-errs:
		_ = stream.Close()
		if errors.Is(err, io.EOF) {
			return nil
		}
		return err
	}
}

func (t *SandboxTunnel) start(ctx context.Context, localPort int) error {
	session, err := t.connect(ctx)
	if err != nil {
		return err
	}
	t.session = session

	listener, err := net.Listen("tcp", fmt.Sprintf("127.0.0.1:%d", localPort))
	if err != nil {
		session.close()
		return err
	}
	t.listener = listener
	t.LocalPort = listener.Addr().(*net.TCPAddr).Port

	go t.acceptLoop()
	return nil
}

func (t *SandboxTunnel) acceptLoop() {
	for {
		conn, err := t.listener.Accept()
		if err != nil {
			return
		}
		go t.handleConn(conn)
	}
}

func (t *SandboxTunnel) handleConn(conn net.Conn) {
	defer conn.Close()

	stream, err := t.Dial(context.Background())
	if err != nil {
		return
	}
	defer stream.Close()

	BridgeSandboxTunnel(context.Background(), stream, conn)
}

func (t *SandboxTunnel) ensureSession(ctx context.Context) (*yamuxSession, error) {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil, errors.New("sandbox tunnel is closed")
	}
	if t.session != nil && !t.session.isClosed() {
		session := t.session
		t.mu.Unlock()
		return session, nil
	}
	t.mu.Unlock()

	var lastErr error
	for attempt := 0; attempt < t.maxReconnects; attempt++ {
		session, err := t.connect(ctx)
		if err == nil {
			t.mu.Lock()
			if t.session != nil {
				t.session.close()
			}
			t.session = session
			t.mu.Unlock()
			return session, nil
		}
		lastErr = err
		delay := 500 * time.Millisecond << attempt
		if delay > 8*time.Second {
			delay = 8 * time.Second
		}
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
	return nil, fmt.Errorf("sandbox tunnel reconnect failed after %d attempts: %w", t.maxReconnects, lastErr)
}

func (t *SandboxTunnel) connect(ctx context.Context) (*yamuxSession, error) {
	return connectSandboxTunnelSession(ctx, t.dataplaneURL, t.opts)
}

func connectSandboxTunnelSession(ctx context.Context, dataplaneURL string, opts []option.RequestOption) (*yamuxSession, error) {
	wsURL, err := sandboxTunnelWebSocketURL(dataplaneURL)
	if err != nil {
		return nil, err
	}
	ws, err := dialSandboxWebSocketURL(ctx, wsURL, opts...)
	if err != nil {
		return nil, err
	}
	return newYamuxSession(newWebSocketByteStream(ws)), nil
}

func openSandboxTunnelStream(session *yamuxSession, remotePort int, closeSession bool) (*SandboxTunnelStream, error) {
	stream, err := session.openStream()
	if err != nil {
		return nil, err
	}
	header := []byte{SandboxTunnelProtocolVersion, byte(remotePort >> 8), byte(remotePort)}
	if _, err := stream.write(header); err != nil {
		stream.close()
		return nil, err
	}
	status := []byte{0}
	if _, err := io.ReadFull(stream, status); err != nil {
		stream.close()
		return nil, err
	}
	tunnelStatus := SandboxTunnelStatus(status[0])
	if tunnelStatus != SandboxTunnelStatusOK {
		stream.close()
		return nil, &SandboxTunnelStatusError{RemotePort: remotePort, Status: tunnelStatus}
	}
	return &SandboxTunnelStream{
		RemotePort:   remotePort,
		stream:       stream,
		session:      session,
		closeSession: closeSession,
	}, nil
}

func validateSandboxTunnelPort(port int, name string) error {
	if port < 1 || port > 65535 {
		return fmt.Errorf("%s must be between 1 and 65535, got %d", name, port)
	}
	return nil
}

func copySandboxTunnelUntilDone(ctx context.Context, dst io.Writer, src io.Reader, srcCloser io.Closer) {
	done := make(chan struct{})
	go func() {
		_, _ = io.Copy(dst, src)
		close(done)
	}()

	select {
	case <-done:
	case <-ctx.Done():
		_ = srcCloser.Close()
		<-done
	}
}
