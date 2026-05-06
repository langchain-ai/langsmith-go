package langsmith

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
	"golang.org/x/net/websocket"
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

func sandboxTunnelWebSocketURL(dataplaneURL string) (string, error) {
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
	u.Path = strings.TrimRight(u.Path, "/") + "/tunnel"
	u.RawQuery = ""
	return u.String(), nil
}

type webSocketByteStream struct {
	ws      *websocket.Conn
	readMu  sync.Mutex
	writeMu sync.Mutex
	buf     []byte
}

func newWebSocketByteStream(ws *websocket.Conn) *webSocketByteStream {
	return &webSocketByteStream{ws: ws}
}

func (s *webSocketByteStream) Read(p []byte) (int, error) {
	s.readMu.Lock()
	defer s.readMu.Unlock()
	for len(s.buf) == 0 {
		var msg []byte
		if err := websocket.Message.Receive(s.ws, &msg); err != nil {
			return 0, err
		}
		s.buf = msg
	}
	n := copy(p, s.buf)
	s.buf = s.buf[n:]
	return n, nil
}

func (s *webSocketByteStream) Write(p []byte) (int, error) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if err := websocket.Message.Send(s.ws, p); err != nil {
		return 0, err
	}
	return len(p), nil
}

func (s *webSocketByteStream) Close() error {
	return s.ws.Close()
}

type yamuxSession struct {
	conn      io.ReadWriteCloser
	streams   map[uint32]*yamuxStream
	nextID    uint32
	mu        sync.Mutex
	writeMu   sync.Mutex
	closed    bool
	closeOnce sync.Once
}

func newYamuxSession(conn io.ReadWriteCloser) *yamuxSession {
	s := &yamuxSession{
		conn:    conn,
		streams: make(map[uint32]*yamuxStream),
		nextID:  1,
	}
	go s.readLoop()
	go s.keepAliveLoop()
	return s
}

func (s *yamuxSession) isClosed() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.closed
}

func (s *yamuxSession) openStream() (*yamuxStream, error) {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return nil, errors.New("yamux session is closed")
	}
	id := s.nextID
	s.nextID += 2
	stream := newYamuxStream(id, s)
	s.streams[id] = stream
	s.mu.Unlock()

	if err := s.sendFrame(yamuxTypeWindowUpdate, yamuxFlagSYN, id, 0); err != nil {
		return nil, err
	}
	return stream, nil
}

func (s *yamuxSession) close() {
	s.closeOnce.Do(func() {
		s.mu.Lock()
		s.closed = true
		for _, stream := range s.streams {
			stream.receiveRST()
		}
		s.mu.Unlock()
		_ = s.sendFrame(yamuxTypeGoAway, 0, 0, 0)
		_ = s.conn.Close()
	})
}

func (s *yamuxSession) readLoop() {
	header := make([]byte, yamuxHeaderSize)
	for {
		if _, err := io.ReadFull(s.conn, header); err != nil {
			s.close()
			return
		}
		msgType := header[1]
		flags := binary.BigEndian.Uint16(header[2:4])
		streamID := binary.BigEndian.Uint32(header[4:8])
		length := binary.BigEndian.Uint32(header[8:12])

		switch msgType {
		case yamuxTypeData:
			payload := make([]byte, length)
			if length > 0 {
				if _, err := io.ReadFull(s.conn, payload); err != nil {
					s.close()
					return
				}
			}
			s.handleData(flags, streamID, payload)
		case yamuxTypeWindowUpdate:
			s.handleWindowUpdate(flags, streamID, length)
		case yamuxTypePing:
			if flags&yamuxFlagSYN != 0 {
				_ = s.sendFrame(yamuxTypePing, yamuxFlagACK, 0, length)
			}
		case yamuxTypeGoAway:
			s.close()
			return
		}
	}
}

func (s *yamuxSession) handleData(flags uint16, streamID uint32, payload []byte) {
	stream := s.getStream(streamID)
	if stream == nil {
		return
	}
	if len(payload) > 0 {
		stream.receiveData(payload)
	}
	if flags&yamuxFlagFIN != 0 {
		stream.receiveFIN()
	}
	if flags&yamuxFlagRST != 0 {
		stream.receiveRST()
	}
}

func (s *yamuxSession) handleWindowUpdate(flags uint16, streamID uint32, delta uint32) {
	stream := s.getStream(streamID)
	if stream == nil {
		return
	}
	if delta > 0 {
		stream.updateSendWindow(delta)
	}
	if flags&yamuxFlagFIN != 0 {
		stream.receiveFIN()
	}
	if flags&yamuxFlagRST != 0 {
		stream.receiveRST()
	}
}

func (s *yamuxSession) keepAliveLoop() {
	ticker := time.NewTicker(30 * time.Second)
	defer ticker.Stop()
	var ping uint32
	for range ticker.C {
		if s.isClosed() {
			return
		}
		ping++
		if err := s.sendFrame(yamuxTypePing, yamuxFlagSYN, 0, ping); err != nil {
			s.close()
			return
		}
	}
}

func (s *yamuxSession) getStream(streamID uint32) *yamuxStream {
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.streams[streamID]
}

func (s *yamuxSession) sendFrame(msgType byte, flags uint16, streamID uint32, length uint32) error {
	var header [yamuxHeaderSize]byte
	header[0] = yamuxVersion
	header[1] = msgType
	binary.BigEndian.PutUint16(header[2:4], flags)
	binary.BigEndian.PutUint32(header[4:8], streamID)
	binary.BigEndian.PutUint32(header[8:12], length)
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	_, err := s.conn.Write(header[:])
	return err
}

func (s *yamuxSession) sendData(streamID uint32, data []byte) error {
	var header [yamuxHeaderSize]byte
	header[0] = yamuxVersion
	header[1] = yamuxTypeData
	binary.BigEndian.PutUint32(header[4:8], streamID)
	binary.BigEndian.PutUint32(header[8:12], uint32(len(data)))
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	if _, err := s.conn.Write(header[:]); err != nil {
		return err
	}
	_, err := s.conn.Write(data)
	return err
}

func (s *yamuxSession) sendWindowUpdate(streamID uint32, delta uint32) error {
	return s.sendFrame(yamuxTypeWindowUpdate, 0, streamID, delta)
}

type yamuxStream struct {
	id      uint32
	session *yamuxSession

	recvMu     sync.Mutex
	recvCond   *sync.Cond
	recvBuf    []byte
	recvClosed bool
	recvErr    bool
	recvWindow uint32

	sendMu     sync.Mutex
	sendCond   *sync.Cond
	sendClosed bool
	sendWindow uint32
}

func newYamuxStream(id uint32, session *yamuxSession) *yamuxStream {
	stream := &yamuxStream{
		id:         id,
		session:    session,
		recvWindow: yamuxInitialWindowSize,
		sendWindow: yamuxInitialWindowSize,
	}
	stream.recvCond = sync.NewCond(&stream.recvMu)
	stream.sendCond = sync.NewCond(&stream.sendMu)
	return stream
}

func (s *yamuxStream) Read(p []byte) (int, error) {
	s.recvMu.Lock()
	for len(s.recvBuf) == 0 && !s.recvClosed && !s.recvErr {
		s.recvCond.Wait()
	}
	if s.recvErr && len(s.recvBuf) == 0 {
		s.recvMu.Unlock()
		return 0, io.ErrUnexpectedEOF
	}
	if len(s.recvBuf) == 0 {
		s.recvMu.Unlock()
		return 0, io.EOF
	}
	n := copy(p, s.recvBuf)
	s.recvBuf = s.recvBuf[n:]
	consumed := yamuxInitialWindowSize - s.recvWindow
	if consumed >= yamuxInitialWindowSize/2 {
		s.recvWindow += consumed
		s.recvMu.Unlock()
		_ = s.session.sendWindowUpdate(s.id, consumed)
		return n, nil
	}
	s.recvMu.Unlock()
	return n, nil
}

func (s *yamuxStream) Write(p []byte) (int, error) {
	return s.write(p)
}

func (s *yamuxStream) write(p []byte) (int, error) {
	offset := 0
	for offset < len(p) {
		s.sendMu.Lock()
		for s.sendWindow == 0 && !s.sendClosed {
			s.sendCond.Wait()
		}
		if s.sendClosed {
			s.sendMu.Unlock()
			return offset, io.ErrClosedPipe
		}
		chunk := min(len(p)-offset, int(s.sendWindow))
		s.sendWindow -= uint32(chunk)
		s.sendMu.Unlock()

		if err := s.session.sendData(s.id, p[offset:offset+chunk]); err != nil {
			return offset, err
		}
		offset += chunk
	}
	return len(p), nil
}

func (s *yamuxStream) Close() error {
	s.close()
	return nil
}

func (s *yamuxStream) close() {
	s.sendMu.Lock()
	if !s.sendClosed {
		s.sendClosed = true
		_ = s.session.sendFrame(yamuxTypeData, yamuxFlagFIN, s.id, 0)
	}
	s.sendMu.Unlock()
	s.recvMu.Lock()
	s.recvClosed = true
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
}

func (s *yamuxStream) receiveData(data []byte) {
	s.recvMu.Lock()
	s.recvBuf = append(s.recvBuf, data...)
	s.recvWindow -= uint32(len(data))
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
}

func (s *yamuxStream) receiveFIN() {
	s.recvMu.Lock()
	s.recvClosed = true
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
}

func (s *yamuxStream) receiveRST() {
	s.recvMu.Lock()
	s.recvErr = true
	s.recvCond.Broadcast()
	s.recvMu.Unlock()
	s.sendMu.Lock()
	s.sendClosed = true
	s.sendCond.Broadcast()
	s.sendMu.Unlock()
}

func (s *yamuxStream) updateSendWindow(delta uint32) {
	s.sendMu.Lock()
	s.sendWindow += delta
	s.sendCond.Broadcast()
	s.sendMu.Unlock()
}
