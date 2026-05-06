package langsmith

import (
	"context"
	"fmt"
	"net"
	"slices"
	"sync"

	"github.com/langchain-ai/langsmith-go/option"
)

const (
	// SandboxTunnelProtocolVersion is the per-stream connect header protocol
	// version layered on top of yamux streams.
	SandboxTunnelProtocolVersion = byte(0x01)
)

// SandboxTunnelStatus is the daemon's one-byte response to a tunnel stream
// connect request.
type SandboxTunnelStatus byte

const (
	SandboxTunnelStatusOK                 SandboxTunnelStatus = 0x00
	SandboxTunnelStatusPortNotAllowed     SandboxTunnelStatus = 0x01
	SandboxTunnelStatusDialFailed         SandboxTunnelStatus = 0x02
	SandboxTunnelStatusUnsupportedVersion SandboxTunnelStatus = 0x03

	yamuxVersion           = byte(0)
	yamuxTypeData          = byte(0)
	yamuxTypeWindowUpdate  = byte(1)
	yamuxTypePing          = byte(2)
	yamuxTypeGoAway        = byte(3)
	yamuxFlagSYN           = uint16(0x0001)
	yamuxFlagACK           = uint16(0x0002)
	yamuxFlagFIN           = uint16(0x0004)
	yamuxFlagRST           = uint16(0x0008)
	yamuxHeaderSize        = 12
	yamuxInitialWindowSize = uint32(256 * 1024)
)

// Reason returns a stable human-readable status reason.
func (s SandboxTunnelStatus) Reason() string {
	switch s {
	case SandboxTunnelStatusOK:
		return "ok"
	case SandboxTunnelStatusPortNotAllowed:
		return "port not allowed"
	case SandboxTunnelStatusDialFailed:
		return "dial failed"
	case SandboxTunnelStatusUnsupportedVersion:
		return "unsupported protocol version"
	default:
		return "unknown"
	}
}

// SandboxTunnelParams configures a TCP tunnel to a sandbox port.
type SandboxTunnelParams struct {
	LocalPort     int
	MaxReconnects int
}

// SandboxTunnelStatusError is returned when the sandbox daemon rejects a tunnel
// stream connection.
type SandboxTunnelStatusError struct {
	RemotePort int
	Status     SandboxTunnelStatus
}

func (e *SandboxTunnelStatusError) Error() string {
	return fmt.Sprintf("sandbox tunnel rejected port %d: %s", e.RemotePort, e.Status.Reason())
}

// SandboxTunnelStream is a single connected stream to a TCP port inside a
// sandbox.
type SandboxTunnelStream struct {
	RemotePort int

	stream       *yamuxStream
	session      *yamuxSession
	closeSession bool
}

// SandboxTunnel forwards local TCP connections to a port inside a sandbox.
type SandboxTunnel struct {
	LocalPort  int
	RemotePort int

	dataplaneURL  string
	opts          []option.RequestOption
	maxReconnects int

	listener net.Listener
	session  *yamuxSession
	mu       sync.Mutex
	closed   bool
}

// Tunnel opens a TCP tunnel to a port inside the named sandbox.
func (r *SandboxBoxService) Tunnel(ctx context.Context, name string, remotePort int, params SandboxTunnelParams, opts ...option.RequestOption) (*SandboxTunnel, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.TunnelWithDataplaneURL(ctx, dataplaneURL, remotePort, params, opts...)
}

// OpenTunnelStream opens a single TCP stream to a port inside the named
// sandbox. This is useful for stdio bridges such as SSH ProxyCommand.
func (r *SandboxBoxService) OpenTunnelStream(ctx context.Context, name string, remotePort int, opts ...option.RequestOption) (*SandboxTunnelStream, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.OpenTunnelStreamWithDataplaneURL(ctx, dataplaneURL, remotePort, opts...)
}

// OpenTunnelStreamWithDataplaneURL opens a single TCP stream directly against a
// sandbox dataplane URL.
func (r *SandboxBoxService) OpenTunnelStreamWithDataplaneURL(ctx context.Context, dataplaneURL string, remotePort int, opts ...option.RequestOption) (*SandboxTunnelStream, error) {
	if err := validateSandboxTunnelPort(remotePort, "remotePort"); err != nil {
		return nil, err
	}
	session, err := connectSandboxTunnelSession(ctx, dataplaneURL, slices.Concat(r.Options, opts))
	if err != nil {
		return nil, err
	}
	stream, err := openSandboxTunnelStream(session, remotePort, true)
	if err != nil {
		session.close()
		return nil, err
	}
	return stream, nil
}

// TunnelWithDataplaneURL opens a TCP tunnel directly against a sandbox dataplane
// URL.
func (r *SandboxBoxService) TunnelWithDataplaneURL(ctx context.Context, dataplaneURL string, remotePort int, params SandboxTunnelParams, opts ...option.RequestOption) (*SandboxTunnel, error) {
	if err := validateSandboxTunnelPort(remotePort, "remotePort"); err != nil {
		return nil, err
	}
	if params.LocalPort < 0 || params.LocalPort > 65535 {
		return nil, fmt.Errorf("localPort must be between 0 and 65535, got %d", params.LocalPort)
	}
	maxReconnects := params.MaxReconnects
	if maxReconnects == 0 {
		maxReconnects = 3
	}
	t := &SandboxTunnel{
		RemotePort:    remotePort,
		dataplaneURL:  dataplaneURL,
		opts:          slices.Concat(r.Options, opts),
		maxReconnects: maxReconnects,
	}
	if err := t.start(ctx, params.LocalPort); err != nil {
		return nil, err
	}
	return t, nil
}

// OpenTunnelStream opens a single TCP stream to a port inside this sandbox.
func (s *Sandbox) OpenTunnelStream(ctx context.Context, remotePort int, opts ...option.RequestOption) (*SandboxTunnelStream, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.OpenTunnelStreamWithDataplaneURL(ctx, dataplaneURL, remotePort, opts...)
}

// Tunnel opens a TCP tunnel from localhost to a port inside this sandbox.
func (s *Sandbox) Tunnel(ctx context.Context, remotePort int, params SandboxTunnelParams, opts ...option.RequestOption) (*SandboxTunnel, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.TunnelWithDataplaneURL(ctx, dataplaneURL, remotePort, params, opts...)
}

// Close shuts down the tunnel.
func (t *SandboxTunnel) Close() error {
	t.mu.Lock()
	if t.closed {
		t.mu.Unlock()
		return nil
	}
	t.closed = true
	listener := t.listener
	session := t.session
	t.mu.Unlock()

	var err error
	if listener != nil {
		err = listener.Close()
	}
	if session != nil {
		session.close()
	}
	return err
}

// Dial opens a single stream over this tunnel's managed session.
func (t *SandboxTunnel) Dial(ctx context.Context) (*SandboxTunnelStream, error) {
	session, err := t.ensureSession(ctx)
	if err != nil {
		return nil, err
	}
	return openSandboxTunnelStream(session, t.RemotePort, false)
}

// Read implements io.Reader.
func (s *SandboxTunnelStream) Read(p []byte) (int, error) {
	return s.stream.Read(p)
}

// Write implements io.Writer.
func (s *SandboxTunnelStream) Write(p []byte) (int, error) {
	return s.stream.Write(p)
}

// Close closes the stream. Streams opened with OpenTunnelStream also close
// their underlying tunnel session.
func (s *SandboxTunnelStream) Close() error {
	err := s.stream.Close()
	if s.closeSession && s.session != nil {
		s.session.close()
	}
	return err
}
