package connector

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"errors"
	"fmt"
	"log/slog"
	"net"
	"sync"
	"sync/atomic"
	"time"

	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
)

const (
	defaultFlowRequestTimeout = 10 * time.Second
	overloadFrameTimeout      = 250 * time.Millisecond
	maxConcurrentRejections   = 16
	maxTargetDialTimeout      = 30 * time.Second
	flowResultWriteGrace      = 100 * time.Millisecond
)

// TargetDialer opens a configured local TCP target.
type TargetDialer interface {
	DialContext(ctx context.Context, network, address string) (net.Conn, error)
}

type flowHandler struct {
	targets         TargetMap
	dialer          TargetDialer
	resolver        Resolver
	targetTLSConfig *tls.Config
	slots           chan struct{}
	rejectionSlots  chan struct{}
	requestTimeout  time.Duration
	logger          *slog.Logger
}

func newFlowHandler(targets TargetMap, dialer TargetDialer, maxConcurrentFlows int) (*flowHandler, error) {
	if dialer == nil {
		return nil, fmt.Errorf("connector target dialer is required")
	}
	if maxConcurrentFlows <= 0 {
		return nil, fmt.Errorf("connector max concurrent flows must be positive")
	}
	return &flowHandler{
		targets:  targets,
		dialer:   dialer,
		resolver: net.DefaultResolver,
		targetTLSConfig: &tls.Config{
			MinVersion: tls.VersionTLS12,
			NextProtos: []string{"http/1.1"},
		},
		slots:          make(chan struct{}, maxConcurrentFlows),
		rejectionSlots: make(chan struct{}, maxConcurrentRejections),
		requestTimeout: defaultFlowRequestTimeout,
		logger:         slog.Default(),
	}, nil
}

func (h *flowHandler) handle(ctx context.Context, stream net.Conn) error {
	admitted := h.tryAcquire()
	if admitted {
		defer h.release()
	}
	return h.handleAdmitted(ctx, stream, admitted)
}

// dispatch keeps the accept loop nonblocking while bounding both active flows
// and the goroutines that send capacity errors. Once the rejection pool is
// full, additional streams are closed immediately.
func (h *flowHandler) dispatch(ctx context.Context, stream net.Conn, flows *sync.WaitGroup) {
	if !h.tryAcquire() {
		select {
		case h.rejectionSlots <- struct{}{}:
			flows.Add(1)
			go func() {
				defer flows.Done()
				defer func() { <-h.rejectionSlots }()
				_ = h.rejectAtCapacity(ctx, stream)
			}()
		default:
			_ = stream.Close()
		}
		return
	}
	flows.Add(1)
	go func() {
		defer flows.Done()
		defer h.release()
		_ = h.handleAdmitted(ctx, stream, true)
	}()
}

func (h *flowHandler) rejectAtCapacity(ctx context.Context, stream net.Conn) error {
	defer stream.Close()
	stopCancellation := context.AfterFunc(ctx, func() {
		_ = stream.SetDeadline(time.Now())
	})
	defer stopCancellation()
	if err := stream.SetDeadline(time.Now().Add(overloadFrameTimeout)); err != nil {
		return err
	}
	var open transport.FlowOpen
	if err := transport.ReadFrame(stream, transport.FrameTypeFlowOpen, &open); err != nil {
		return err
	}
	if err := transport.ValidateFlowOpen(open); err != nil {
		return err
	}
	return writeFlowResult(stream, open.FlowID, transport.FlowStatusCapacityExceeded)
}

func (h *flowHandler) tryAcquire() bool {
	select {
	case h.slots <- struct{}{}:
		return true
	default:
		return false
	}
}

func (h *flowHandler) release() {
	<-h.slots
}

func (h *flowHandler) handleAdmitted(ctx context.Context, stream net.Conn, admitted bool) error {
	defer stream.Close()
	stopCancellation := context.AfterFunc(ctx, func() {
		_ = stream.SetDeadline(time.Now())
	})
	defer stopCancellation()
	if err := stream.SetDeadline(time.Now().Add(h.requestTimeout)); err != nil {
		return fmt.Errorf("set flow request deadline: %w", err)
	}

	var open transport.FlowOpen
	if err := transport.ReadFrame(stream, transport.FrameTypeFlowOpen, &open); err != nil {
		return fmt.Errorf("read flow request: %w", err)
	}
	if err := transport.ValidateFlowOpen(open); err != nil {
		return writeFlowResult(stream, open.FlowID, transport.FlowStatusProtocolError)
	}
	targetConfig, ok := h.targets.Resolve(open.TargetName)
	if !ok {
		h.logger.WarnContext(ctx, "edge requested a target this connector does not serve",
			"target", open.TargetName)
		return writeFlowResult(stream, open.FlowID, transport.FlowStatusUnknownTarget)
	}
	// Only a proxy flow names its destination; any other target's address is
	// fixed by this connector's configuration.
	if (targetConfig.Proxy != nil) != (open.Destination != "") {
		return writeFlowResult(stream, open.FlowID, transport.FlowStatusProtocolError)
	}
	address := targetConfig.Address
	if targetConfig.Proxy != nil {
		address = open.Destination
	}

	if !admitted {
		return writeFlowResult(stream, open.FlowID, transport.FlowStatusCapacityExceeded)
	}

	dialTimeout := min(
		time.Duration(open.DialTimeoutMilliseconds)*time.Millisecond,
		maxTargetDialTimeout,
	)
	if err := stream.SetDeadline(time.Now().Add(dialTimeout + flowResultWriteGrace)); err != nil {
		return fmt.Errorf("set flow response deadline: %w", err)
	}
	dialCtx, cancelDial := context.WithTimeout(ctx, dialTimeout)
	target, err := h.dialTarget(dialCtx, targetConfig, open.Destination)
	dialContextErr := dialCtx.Err()
	cancelDial()
	if errors.Is(err, errDestinationDenied) {
		h.logger.WarnContext(ctx, "relay policy denied destination",
			"flow_id", open.FlowID, "destination", open.Destination)
		return writeFlowResult(stream, open.FlowID, transport.FlowStatusDestinationDenied)
	}
	if err != nil {
		if ctx.Err() != nil {
			return ctx.Err()
		}
		status := transport.FlowStatusDialFailed
		var netErr net.Error
		if errors.Is(err, context.DeadlineExceeded) ||
			errors.Is(dialContextErr, context.DeadlineExceeded) ||
			errors.As(err, &netErr) && netErr.Timeout() {
			status = transport.FlowStatusDialTimeout
		}
		// The edge only learns the bounded status, so this line is where an
		// operator finds out why their target could not be reached.
		attrs := []any{
			"flow_id", open.FlowID,
			"target", open.TargetName,
			"address", address,
			"status", status,
			"error", err,
		}
		if hint := dialFailureHint(err); hint != "" {
			attrs = append(attrs, "hint", hint)
		}
		h.logger.WarnContext(ctx, "connector could not reach target", attrs...)
		return writeFlowResult(stream, open.FlowID, status)
	}
	defer target.Close()
	if targetConfig.Proxy != nil {
		h.logger.InfoContext(ctx, "connector opened relay flow",
			"flow_id", open.FlowID, "destination", address)
	} else {
		h.logger.DebugContext(ctx, "connector opened flow to target",
			"target", open.TargetName, "address", address)
	}

	if err := writeFlowResult(stream, open.FlowID, transport.FlowStatusReady); err != nil {
		return err
	}
	if err := stream.SetDeadline(time.Time{}); err != nil {
		return fmt.Errorf("clear flow stream deadline: %w", err)
	}
	return h.bridgeTarget(ctx, stream, target, open, targetConfig, address)
}

// bridgeTarget carries a flow's bytes and then reports how it ended. An HTTP
// target that never answered is a warning: the edge only tells the caller it
// gave up, so this line is where an operator learns the origin was too slow
// or closed without responding.
func (h *flowHandler) bridgeTarget(
	ctx context.Context,
	stream net.Conn,
	target net.Conn,
	open transport.FlowOpen,
	targetConfig Target,
	address string,
) error {
	observed := &observedTarget{
		Conn: target,
		http: targetConfig.Descriptor.Kind == transport.TargetKindHTTP,
	}
	started := time.Now()
	err := transport.Bridge(ctx, stream, observed)
	attrs := []any{
		"flow_id", open.FlowID,
		"target", open.TargetName,
		"address", address,
		"duration", time.Since(started).Round(time.Millisecond).String(),
		"bytes_to_target", observed.written.Load(),
		"bytes_from_target", observed.read.Load(),
	}
	if err != nil {
		attrs = append(attrs, "error", err)
	}
	unanswered := observed.http && observed.written.Load() > 0 && !observed.answered.Load()
	if unanswered && ctx.Err() == nil {
		h.logger.WarnContext(ctx, "target did not respond before the flow closed", attrs...)
	} else {
		h.logger.DebugContext(ctx, "connector closed flow to target", attrs...)
	}
	return err
}

// observedTarget counts the bytes a flow exchanges with its target and, for an
// HTTP target, whether its final response header arrived while the edge was
// still waiting. The edge gives up on a response whose header it has not
// finished reading, so a partial header does not count as an answer. It
// forwards the target's half-close so a finished request still reaches the
// origin as EOF.
type observedTarget struct {
	net.Conn
	http       bool
	read       atomic.Int64
	written    atomic.Int64
	answered   atomic.Bool
	edgeClosed atomic.Bool
	// Only the goroutine copying the target's response calls Read, so the
	// scanner needs no lock.
	header responseHeaderScanner
}

func (t *observedTarget) Read(p []byte) (int, error) {
	n, err := t.Conn.Read(p)
	if n > 0 {
		t.read.Add(int64(n))
		if t.http && !t.answered.Load() && !t.edgeClosed.Load() && t.header.scan(p[:n]) {
			t.answered.Store(true)
		}
	}
	return n, err
}

func (t *observedTarget) Write(p []byte) (int, error) {
	n, err := t.Conn.Write(p)
	t.written.Add(int64(n))
	return n, err
}

// CloseWrite runs when the edge stops sending, which for HTTP means it has
// either read the whole response or given up on it.
func (t *observedTarget) CloseWrite() error {
	t.edgeClosed.Store(true)
	if conn, ok := t.Conn.(interface{ CloseWrite() error }); ok {
		return conn.CloseWrite()
	}
	return nil
}

func (t *observedTarget) Close() error {
	t.edgeClosed.Store(true)
	return t.Conn.Close()
}

func (h *flowHandler) dialTarget(ctx context.Context, target Target, destination string) (net.Conn, error) {
	if target.Proxy != nil {
		return target.Proxy.dial(ctx, destination, h.resolver, h.dialer)
	}
	conn, err := h.dialer.DialContext(ctx, "tcp", target.Address)
	if err != nil {
		return nil, err
	}
	if target.Descriptor.Kind != transport.TargetKindHTTP ||
		target.Descriptor.OriginScheme != "https" {
		return conn, nil
	}

	host, _, err := net.SplitHostPort(target.Address)
	if err != nil {
		_ = conn.Close()
		return nil, err
	}
	tlsConfig := h.targetTLSConfig.Clone()
	tlsConfig.ServerName = host
	if len(tlsConfig.NextProtos) == 0 {
		tlsConfig.NextProtos = []string{"http/1.1"}
	}
	tlsConn := tls.Client(conn, tlsConfig)
	if err := tlsConn.HandshakeContext(ctx); err != nil {
		_ = tlsConn.Close()
		return nil, err
	}
	return tlsConn, nil
}

// dialFailureHint names the fix for TLS failures an operator resolves in the
// connector's own environment, or returns "" when the error speaks for itself.
func dialFailureHint(err error) string {
	var unknownAuthority x509.UnknownAuthorityError
	if errors.As(err, &unknownAuthority) {
		return "the target's certificate is signed by a CA the connector does not trust; " +
			"mount that CA and point SSL_CERT_FILE at it"
	}
	var hostname x509.HostnameError
	if errors.As(err, &hostname) {
		return "the target's certificate does not list the hostname in the target address"
	}
	return ""
}

func writeFlowResult(stream net.Conn, flowID string, status transport.FlowStatus) error {
	result := transport.FlowResult{FlowID: flowID, Status: status}
	if _, err := transport.WriteFrame(stream, transport.FrameTypeFlowResult, result); err != nil {
		return fmt.Errorf("write flow result: %w", err)
	}
	return nil
}
