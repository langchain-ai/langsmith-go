package connector

import (
	"context"
	"crypto/tls"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
)

const (
	defaultMaxConcurrentFlows = 128
	maxRejectionReasonBytes   = 512
	carrierHandshakeTimeout   = 15 * time.Second
	defaultDrainGracePeriod   = 30 * time.Second
	drainWriteTimeout         = 5 * time.Second
)

// ClientOption configures connector transport dependencies.
type ClientOption func(*clientOptions) error

type clientOptions struct {
	webSocketDialer    *websocket.Dialer
	tlsConfig          *tls.Config
	targetDialer       TargetDialer
	resolver           Resolver
	maxConcurrentFlows int
	drainGracePeriod   time.Duration
	logger             *slog.Logger
}

// WithWebSocketDialer replaces the outbound carrier dialer. The connector
// still restricts the WebSocket negotiation to its single wire subprotocol.
func WithWebSocketDialer(dialer *websocket.Dialer) ClientOption {
	return func(options *clientOptions) error {
		if dialer == nil {
			return fmt.Errorf("connector websocket dialer is required")
		}
		copy := *dialer
		options.webSocketDialer = &copy
		return nil
	}
}

// WithTLSConfig supplies TLS settings for the outbound WSS connection.
func WithTLSConfig(config *tls.Config) ClientOption {
	return func(options *clientOptions) error {
		if config == nil {
			return fmt.Errorf("connector TLS config is required")
		}
		options.tlsConfig = config.Clone()
		return nil
	}
}

// WithTargetDialer replaces the local TCP dialer.
func WithTargetDialer(dialer TargetDialer) ClientOption {
	return func(options *clientOptions) error {
		if dialer == nil {
			return fmt.Errorf("connector target dialer is required")
		}
		options.targetDialer = dialer
		return nil
	}
}

// WithResolver replaces how a relay looks up a destination's addresses. It
// defaults to net.DefaultResolver.
func WithResolver(resolver Resolver) ClientOption {
	return func(options *clientOptions) error {
		if resolver == nil {
			return fmt.Errorf("connector resolver is required")
		}
		options.resolver = resolver
		return nil
	}
}

// WithMaxConcurrentFlows bounds the number of active local TCP flows.
func WithMaxConcurrentFlows(limit int) ClientOption {
	return func(options *clientOptions) error {
		if limit <= 0 {
			return fmt.Errorf("connector max concurrent flows must be positive")
		}
		options.maxConcurrentFlows = limit
		return nil
	}
}

// WithDrainGracePeriod sets how long Run keeps accepting and serving flows
// after announcing a drain before the carrier is torn down.
func WithDrainGracePeriod(period time.Duration) ClientOption {
	return func(options *clientOptions) error {
		if period <= 0 {
			return fmt.Errorf("connector drain grace period must be positive")
		}
		options.drainGracePeriod = period
		return nil
	}
}

// WithLogger sets where the connector reports carrier and target failures.
// It defaults to slog.Default().
func WithLogger(logger *slog.Logger) ClientOption {
	return func(options *clientOptions) error {
		if logger == nil {
			return fmt.Errorf("connector logger is required")
		}
		options.logger = logger
		return nil
	}
}

// Client owns one connector instance identity and a shared target flow limit.
type Client struct {
	config              Config
	connectorInstanceID uuid.UUID
	webSocketDialer     *websocket.Dialer
	flowHandler         *flowHandler
	drainGracePeriod    time.Duration
	logger              *slog.Logger
}

// NewClient constructs a connector carrier client from validated config.
func NewClient(config Config, options ...ClientOption) (*Client, error) {
	if err := validateCarrierURL(config.Endpoint); err != nil {
		return nil, err
	}
	if config.ConnectionID == uuid.Nil {
		return nil, fmt.Errorf("connector connection ID is required")
	}
	if len(config.Targets.Names()) == 0 {
		return nil, fmt.Errorf("connector targets are required")
	}
	// The credential may come from a file or directly from the environment;
	// either way it must resolve before a carrier is attempted.
	if _, err := config.ReadAPIKey(); err != nil {
		return nil, err
	}

	dialer := &websocket.Dialer{
		Proxy:            http.ProxyFromEnvironment,
		HandshakeTimeout: carrierHandshakeTimeout,
	}
	clientOptions := clientOptions{
		webSocketDialer:    dialer,
		targetDialer:       newDefaultTargetDialer(),
		resolver:           net.DefaultResolver,
		maxConcurrentFlows: defaultMaxConcurrentFlows,
		drainGracePeriod:   defaultDrainGracePeriod,
		logger:             slog.Default(),
	}
	for _, option := range options {
		if option == nil {
			return nil, fmt.Errorf("connector client option is required")
		}
		if err := option(&clientOptions); err != nil {
			return nil, err
		}
	}
	if clientOptions.webSocketDialer == nil {
		return nil, fmt.Errorf("connector websocket dialer is required")
	}
	if clientOptions.targetDialer == nil {
		return nil, fmt.Errorf("connector target dialer is required")
	}
	if clientOptions.maxConcurrentFlows <= 0 {
		return nil, fmt.Errorf("connector max concurrent flows must be positive")
	}

	clientDialer := *clientOptions.webSocketDialer
	clientDialer.Subprotocols = []string{transport.WebSocketSubprotocol}
	if clientDialer.HandshakeTimeout <= 0 {
		clientDialer.HandshakeTimeout = carrierHandshakeTimeout
	}
	if clientOptions.tlsConfig != nil {
		clientDialer.TLSClientConfig = clientOptions.tlsConfig.Clone()
	} else if clientDialer.TLSClientConfig != nil {
		clientDialer.TLSClientConfig = clientDialer.TLSClientConfig.Clone()
	}

	flowHandler, err := newFlowHandler(
		config.Targets,
		clientOptions.targetDialer,
		clientOptions.maxConcurrentFlows,
	)
	if err != nil {
		return nil, err
	}
	flowHandler.logger = clientOptions.logger
	flowHandler.resolver = clientOptions.resolver
	connectorInstanceID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate connector instance ID: %w", err)
	}
	return &Client{
		config:              config,
		connectorInstanceID: connectorInstanceID,
		webSocketDialer:     &clientDialer,
		flowHandler:         flowHandler,
		drainGracePeriod:    clientOptions.drainGracePeriod,
		logger:              clientOptions.logger,
	}, nil
}

func (c *Client) runCarrier(ctx context.Context) (*carrierSession, error) {
	apiKey, err := c.config.ReadAPIKey()
	if err != nil {
		return nil, err
	}
	sessionID, err := uuid.NewV7()
	if err != nil {
		return nil, fmt.Errorf("generate connector session ID: %w", err)
	}
	session, carrierTransport, err := c.dialCarrier(ctx, apiKey, sessionID)
	if err != nil {
		return nil, err
	}
	control, runtimeID, err := c.registerCarrier(ctx, session, sessionID)
	if err != nil {
		_ = session.Close()
		_ = carrierTransport.Close()
		return nil, err
	}
	return &carrierSession{
		sessionID: sessionID,
		runtimeID: runtimeID,
		yamux:     session,
		control:   control,
		transport: carrierTransport,
	}, nil
}

func (c *Client) serveCarrier(ctx context.Context, carrier *carrierSession) error {
	// Flows run under a context independent of ctx so a graceful shutdown can
	// let in-flight flows finish before the session closes.
	flowsCtx, cancelFlows := context.WithCancel(context.Background())
	var flows sync.WaitGroup
	defer func() {
		cancelFlows()
		carrier.close()
		flows.Wait()
	}()

	// Serve edge-opened flows until the context is cancelled.
	if err := c.acceptFlows(ctx, flowsCtx, carrier.yamux, &flows); err != nil {
		return err
	}

	// The context was cancelled (e.g. SIGTERM): make-before-break. Tell the edge
	// to stop routing NEW flows to this carrier first, then keep accepting and
	// serving for a bounded grace window so a flow the edge routed before it
	// processed DRAIN is still served here instead of dropped. The deferred
	// cleanup then cancels any still-running flows and closes the session. The
	// protocol has no drain acknowledgement, so an idle interval cannot prove
	// that a flow routed over a high-latency link before DRAIN has arrived.
	c.sendDrain(carrier.control, carrier.sessionID)
	graceCtx, cancelGrace := context.WithTimeout(context.Background(), c.drainGracePeriod)
	acceptDone := make(chan struct{})
	go func() {
		_ = c.acceptFlows(graceCtx, flowsCtx, carrier.yamux, &flows)
		close(acceptDone)
	}()
	select {
	case <-graceCtx.Done():
	case <-acceptDone:
	}
	cancelGrace()
	<-acceptDone
	return ctx.Err()
}

func (c *Client) dialCarrier(
	ctx context.Context,
	apiKey string,
	sessionID uuid.UUID,
) (*yamux.Session, io.ReadWriteCloser, error) {
	headers := http.Header{}
	headers.Set("X-Api-Key", apiKey)
	headers.Set(transport.HeaderConnectionID, c.config.ConnectionID.String())
	headers.Set(transport.HeaderConnectorInstanceID, c.connectorInstanceID.String())
	headers.Set(transport.HeaderSessionID, sessionID.String())
	headers.Set(transport.HeaderConnectorKind, c.config.Targets.Kind())
	ws, response, err := c.webSocketDialer.DialContext(ctx, c.config.Endpoint, headers)
	headers.Del("X-Api-Key")
	if err != nil {
		if response == nil {
			return nil, nil, fmt.Errorf("connect carrier websocket: %w", err)
		}
		reason := ""
		if response.Body != nil {
			// The edge explains rejections it expects the operator to fix, such
			// as a connector running the wrong kind for its connection.
			body, _ := io.ReadAll(io.LimitReader(response.Body, maxRejectionReasonBytes))
			_ = response.Body.Close()
			reason = strings.TrimSpace(string(body))
		}
		if reason != "" {
			return nil, nil, fmt.Errorf("connect carrier websocket (HTTP %d: %s): %w", response.StatusCode, reason, err)
		}
		return nil, nil, fmt.Errorf("connect carrier websocket (HTTP %d): %w", response.StatusCode, err)
	}
	if response != nil && response.Body != nil {
		_ = response.Body.Close()
	}
	if ws.Subprotocol() != transport.WebSocketSubprotocol {
		_ = ws.Close()
		return nil, nil, fmt.Errorf("connect carrier websocket: unexpected subprotocol")
	}

	streamConn := transport.NewWebSocketConn(ws)
	session, err := yamux.Client(streamConn, transport.YamuxConfig())
	if err != nil {
		_ = streamConn.Close()
		return nil, nil, fmt.Errorf("start connector yamux client: %w", err)
	}
	return session, streamConn, nil
}

func (c *Client) registerCarrier(
	ctx context.Context,
	session *yamux.Session,
	sessionID uuid.UUID,
) (net.Conn, uuid.UUID, error) {
	stopHandshakeCancel := context.AfterFunc(ctx, func() { _ = session.Close() })
	defer stopHandshakeCancel()
	handshakeErr := func(err error) error {
		if ctxErr := ctx.Err(); ctxErr != nil {
			return ctxErr
		}
		return err
	}

	control, err := session.OpenStream()
	if err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("open connector control stream: %w", err))
	}
	keepControl := false
	defer func() {
		if !keepControl {
			_ = control.Close()
		}
	}()
	if err := control.SetDeadline(carrierDeadline(ctx, carrierHandshakeTimeout)); err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("set connector registration deadline: %w", err))
	}
	metadata := c.config.Targets.RegistrationMetadata()
	register := transport.Register{
		ProtocolVersion:     transport.ProtocolVersion,
		ConnectionID:        c.config.ConnectionID.String(),
		ConnectorInstanceID: c.connectorInstanceID.String(),
		SessionID:           sessionID.String(),
		Targets:             metadata.Targets,
		TargetMapHash:       metadata.TargetMapHash,
	}
	if err := transport.ValidateRegister(register); err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("validate connector registration: %w", err))
	}
	if _, err := transport.WriteFrame(control, transport.FrameTypeRegister, register); err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("write connector registration: %w", err))
	}
	var registered transport.Registered
	if err := transport.ReadFrame(control, transport.FrameTypeRegistered, &registered); err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("read edge registration: %w", err))
	}
	if err := validateRegisteredResponse(register, registered); err != nil {
		return nil, uuid.Nil, handshakeErr(err)
	}
	runtimeID, err := uuid.Parse(registered.RuntimeID)
	if err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("parse edge runtime ID: %w", err))
	}
	if err := control.SetDeadline(time.Time{}); err != nil {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("clear connector registration deadline: %w", err))
	}
	if !stopHandshakeCancel() {
		return nil, uuid.Nil, handshakeErr(fmt.Errorf("connector registration canceled"))
	}
	keepControl = true
	return control, runtimeID, nil
}

// acceptFlows accepts and dispatches edge-opened flow streams until acceptCtx is
// cancelled (returns nil) or the session fails (returns the error). Dispatched
// flows run under flowsCtx, so they can outlive acceptCtx during a drain.
func (c *Client) acceptFlows(
	acceptCtx context.Context,
	flowsCtx context.Context,
	session *yamux.Session,
	flows *sync.WaitGroup,
) error {
	for {
		stream, err := session.AcceptStreamWithContext(acceptCtx)
		if err != nil {
			if acceptCtx.Err() != nil {
				return nil
			}
			return fmt.Errorf("accept edge flow stream: %w", err)
		}
		c.flowHandler.dispatch(flowsCtx, stream, flows)
	}
}

// sendDrain tells the edge to stop routing new flows to this carrier so a
// replacement carrier takes them (make-before-break). It is best-effort: if the
// notice cannot be written, the drain grace window and teardown still proceed.
func (c *Client) sendDrain(control net.Conn, sessionID uuid.UUID) {
	_ = control.SetDeadline(time.Now().Add(drainWriteTimeout))
	_, _ = transport.WriteFrame(control, transport.FrameTypeDrain, transport.Drain{
		ProtocolVersion:     transport.ProtocolVersion,
		ConnectionID:        c.config.ConnectionID.String(),
		ConnectorInstanceID: c.connectorInstanceID.String(),
		SessionID:           sessionID.String(),
	})
}

func validateRegisteredResponse(register transport.Register, registered transport.Registered) error {
	if err := transport.ValidateRegistered(registered); err != nil {
		return fmt.Errorf("validate edge registration: %w", err)
	}
	if registered.ProtocolVersion != register.ProtocolVersion ||
		registered.ConnectionID != register.ConnectionID ||
		registered.ConnectorInstanceID != register.ConnectorInstanceID ||
		registered.SessionID != register.SessionID {
		return fmt.Errorf("validate edge registration: identity mismatch")
	}
	return nil
}

func carrierDeadline(ctx context.Context, timeout time.Duration) time.Time {
	deadline := time.Now().Add(timeout)
	if contextDeadline, ok := ctx.Deadline(); ok && contextDeadline.Before(deadline) {
		return contextDeadline
	}
	return deadline
}
