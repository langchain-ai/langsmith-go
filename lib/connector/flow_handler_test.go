package connector

import (
	"bufio"
	"context"
	"crypto/tls"
	"errors"
	"io"
	"log/slog"
	"math"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type targetDialerFunc func(context.Context, string, string) (net.Conn, error)

func (f targetDialerFunc) DialContext(
	ctx context.Context,
	network string,
	address string,
) (net.Conn, error) {
	return f(ctx, network, address)
}

type deadlineIgnoringConn struct {
	net.Conn
}

func (c deadlineIgnoringConn) SetDeadline(time.Time) error {
	return nil
}

type recordedLog struct {
	level   slog.Level
	message string
	attrs   map[string]string
}

// recordingHandler keeps every record so a test can assert what an operator
// would read in the connector's logs.
type recordingHandler struct {
	mu      sync.Mutex
	records []recordedLog
}

func (h *recordingHandler) Enabled(context.Context, slog.Level) bool { return true }

func (h *recordingHandler) Handle(_ context.Context, record slog.Record) error {
	attrs := make(map[string]string)
	record.Attrs(func(attr slog.Attr) bool {
		attrs[attr.Key] = attr.Value.String()
		return true
	})
	h.mu.Lock()
	defer h.mu.Unlock()
	h.records = append(h.records, recordedLog{level: record.Level, message: record.Message, attrs: attrs})
	return nil
}

func (h *recordingHandler) WithAttrs([]slog.Attr) slog.Handler { return h }

func (h *recordingHandler) WithGroup(string) slog.Handler { return h }

func (h *recordingHandler) find(message string) (recordedLog, bool) {
	h.mu.Lock()
	defer h.mu.Unlock()
	for _, record := range h.records {
		if record.message == message {
			return record, true
		}
	}
	return recordedLog{}, false
}

func TestFlowHandlerRejectsUnknownTargetWithoutDialing(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	var dialCalls atomic.Int32
	handler, err := newFlowHandler(targets, targetDialerFunc(func(
		context.Context,
		string,
		string,
	) (net.Conn, error) {
		dialCalls.Add(1)
		return nil, errors.New("must not dial")
	}), 1)
	require.NoError(t, err)

	edgeSide, connectorSide := net.Pipe()
	t.Cleanup(func() { _ = edgeSide.Close() })
	done := make(chan error, 1)
	go func() { done <- handler.handle(context.Background(), connectorSide) }()

	flowID := uuid.Must(uuid.NewV7()).String()
	_, err = transport.WriteFrame(edgeSide, transport.FrameTypeFlowOpen, transport.FlowOpen{
		FlowID:                  flowID,
		TargetName:              "unknown",
		DialTimeoutMilliseconds: 1_000,
	})
	require.NoError(t, err)
	var result transport.FlowResult
	require.NoError(t, transport.ReadFrame(edgeSide, transport.FrameTypeFlowResult, &result))
	assert.Equal(t, transport.FlowResult{
		FlowID: flowID,
		Status: transport.FlowStatusUnknownTarget,
	}, result)
	assert.Zero(t, dialCalls.Load())
	require.NoError(t, <-done)
}

func TestFlowHandlerRejectsInvalidOpenWithBoundedProtocolStatus(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	var dialCalls atomic.Int32
	handler, err := newFlowHandler(targets, targetDialerFunc(func(
		context.Context,
		string,
		string,
	) (net.Conn, error) {
		dialCalls.Add(1)
		return nil, errors.New("must not dial")
	}), 1)
	require.NoError(t, err)

	flowID := uuid.Must(uuid.NewV7()).String()
	result, handleErr := exchangeRejectedFlow(t, handler, transport.FlowOpen{
		FlowID:                  flowID,
		TargetName:              "database",
		DialTimeoutMilliseconds: 0,
	})
	assert.Equal(t, transport.FlowResult{
		FlowID: flowID,
		Status: transport.FlowStatusProtocolError,
	}, result)
	assert.Zero(t, dialCalls.Load())
	require.NoError(t, handleErr)
}

func TestFlowHandlerMapsLocalDialErrorsToBoundedStatuses(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name       string
		dial       targetDialerFunc
		wantStatus transport.FlowStatus
	}{
		{
			name: "timeout",
			dial: func(ctx context.Context, network, address string) (net.Conn, error) {
				assert.Equal(t, "tcp", network)
				assert.Equal(t, "db.internal:5432", address)
				<-ctx.Done()
				return nil, ctx.Err()
			},
			wantStatus: transport.FlowStatusDialTimeout,
		},
		{
			name: "failure",
			dial: func(context.Context, string, string) (net.Conn, error) {
				return nil, errors.New("sensitive local dial detail")
			},
			wantStatus: transport.FlowStatusDialFailed,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
			require.NoError(t, err)
			handler, err := newFlowHandler(targets, tt.dial, 1)
			require.NoError(t, err)
			flowID := uuid.Must(uuid.NewV7()).String()
			result, handleErr := exchangeRejectedFlow(t, handler, transport.FlowOpen{
				FlowID:                  flowID,
				TargetName:              "database",
				DialTimeoutMilliseconds: 5,
			})
			assert.Equal(t, transport.FlowResult{
				FlowID: flowID,
				Status: tt.wantStatus,
			}, result)
			require.NoError(t, handleErr)
		})
	}
}

func TestFlowHandlerCompletesHTTPSHandshakeBeforeReady(t *testing.T) {
	t.Parallel()

	origin := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "/private", r.URL.Path)
		_, _ = w.Write([]byte("through-tunnel"))
	}))
	t.Cleanup(origin.Close)

	targets, err := NewTargetMap(map[string]string{"api": origin.URL})
	require.NoError(t, err)
	handler, err := newFlowHandler(targets, &net.Dialer{}, 1)
	require.NoError(t, err)
	handler.targetTLSConfig = origin.Client().Transport.(*http.Transport).TLSClientConfig.Clone()

	edgeSide, connectorSide := net.Pipe()
	t.Cleanup(func() { _ = edgeSide.Close() })
	done := make(chan error, 1)
	go func() { done <- handler.handle(context.Background(), connectorSide) }()

	flowID := uuid.Must(uuid.NewV7()).String()
	_, err = transport.WriteFrame(edgeSide, transport.FrameTypeFlowOpen, transport.FlowOpen{
		FlowID:                  flowID,
		TargetName:              "api",
		DialTimeoutMilliseconds: 1_000,
	})
	require.NoError(t, err)
	var result transport.FlowResult
	require.NoError(t, transport.ReadFrame(edgeSide, transport.FrameTypeFlowResult, &result))
	require.Equal(t, transport.FlowStatusReady, result.Status)

	request, err := http.NewRequest(http.MethodGet, "http://tunnel.invalid/private", nil)
	require.NoError(t, err)
	require.NoError(t, request.Write(edgeSide))
	response, err := http.ReadResponse(bufio.NewReader(edgeSide), request)
	require.NoError(t, err)
	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)
	require.NoError(t, err)
	require.Equal(t, "through-tunnel", string(body))
	require.NoError(t, edgeSide.Close())
	require.NoError(t, <-done)
}

func TestFlowHandlerRejectsUntrustedHTTPSOriginBeforeReady(t *testing.T) {
	t.Parallel()

	origin := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	t.Cleanup(origin.Close)
	targets, err := NewTargetMap(map[string]string{"api": origin.URL})
	require.NoError(t, err)
	handler, err := newFlowHandler(targets, &net.Dialer{}, 1)
	require.NoError(t, err)
	handler.targetTLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}

	flowID := uuid.Must(uuid.NewV7()).String()
	result, handleErr := exchangeRejectedFlow(t, handler, transport.FlowOpen{
		FlowID:                  flowID,
		TargetName:              "api",
		DialTimeoutMilliseconds: 1_000,
	})
	require.NoError(t, handleErr)
	require.Equal(t, transport.FlowStatusDialFailed, result.Status)
}

func TestFlowHandlerLogsWhyATargetCouldNotBeReached(t *testing.T) {
	t.Parallel()

	untrustedOrigin := httptest.NewTLSServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) {}))
	t.Cleanup(untrustedOrigin.Close)

	tests := []struct {
		name      string
		target    string
		dialer    TargetDialer
		wantError string
		wantHint  string
	}{
		{
			name:      "origin certificate from an untrusted CA",
			target:    untrustedOrigin.URL,
			dialer:    &net.Dialer{},
			wantError: "certificate signed by unknown authority",
			wantHint:  "SSL_CERT_FILE",
		},
		{
			name:   "origin refuses the connection",
			target: "db.internal:5432",
			dialer: targetDialerFunc(func(context.Context, string, string) (net.Conn, error) {
				return nil, errors.New("connect: connection refused")
			}),
			wantError: "connection refused",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			targets, err := NewTargetMap(map[string]string{"api": tt.target})
			require.NoError(t, err)
			handler, err := newFlowHandler(targets, tt.dialer, 1)
			require.NoError(t, err)
			handler.targetTLSConfig = &tls.Config{MinVersion: tls.VersionTLS12}
			logs := &recordingHandler{}
			handler.logger = slog.New(logs)

			result, handleErr := exchangeRejectedFlow(t, handler, transport.FlowOpen{
				FlowID:                  uuid.Must(uuid.NewV7()).String(),
				TargetName:              "api",
				DialTimeoutMilliseconds: 1_000,
			})
			require.NoError(t, handleErr)
			require.Equal(t, transport.FlowStatusDialFailed, result.Status)

			record, ok := logs.find("connector could not reach target")
			require.True(t, ok, "dial failure was not logged")
			assert.Equal(t, slog.LevelWarn, record.level)
			assert.Equal(t, "api", record.attrs["target"])
			assert.Equal(t, string(transport.FlowStatusDialFailed), record.attrs["status"])
			assert.Contains(t, record.attrs["error"], tt.wantError)
			if tt.wantHint == "" {
				assert.NotContains(t, record.attrs, "hint")
			} else {
				assert.Contains(t, record.attrs["hint"], tt.wantHint)
			}
		})
	}
}

func TestFlowHandlerBoundsWaitForFlowOpen(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	handler, err := newFlowHandler(targets, &net.Dialer{}, 1)
	require.NoError(t, err)
	handler.requestTimeout = 20 * time.Millisecond
	edgeSide, connectorSide := net.Pipe()
	t.Cleanup(func() { _ = edgeSide.Close() })
	done := make(chan error, 1)
	go func() { done <- handler.handle(context.Background(), connectorSide) }()

	select {
	case err := <-done:
		require.Error(t, err)
	case <-time.After(time.Second):
		t.Fatal("flow handler remained parked without FLOW_OPEN")
	}
}

func TestFlowHandlerCapsEdgeRequestedDialTimeout(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	observedBudget := make(chan time.Duration, 1)
	handler, err := newFlowHandler(targets, targetDialerFunc(func(
		ctx context.Context,
		_ string,
		_ string,
	) (net.Conn, error) {
		deadline, ok := ctx.Deadline()
		if !ok {
			observedBudget <- 0
		} else {
			observedBudget <- time.Until(deadline)
		}
		return nil, errors.New("dial failed")
	}), 1)
	require.NoError(t, err)
	flowID := uuid.Must(uuid.NewV7()).String()
	result, handleErr := exchangeRejectedFlow(t, handler, transport.FlowOpen{
		FlowID:                  flowID,
		TargetName:              "database",
		DialTimeoutMilliseconds: math.MaxUint32,
	})
	require.NoError(t, handleErr)
	assert.Equal(t, transport.FlowStatusDialFailed, result.Status)
	budget := <-observedBudget
	assert.Positive(t, budget)
	assert.LessOrEqual(t, budget, maxTargetDialTimeout)
}

func TestFlowHandlerRejectsFlowWhenAtCapacity(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	dialStarted := make(chan struct{})
	releaseDial := make(chan struct{})
	var once sync.Once
	handler, err := newFlowHandler(targets, targetDialerFunc(func(
		context.Context,
		string,
		string,
	) (net.Conn, error) {
		once.Do(func() { close(dialStarted) })
		<-releaseDial
		return nil, errors.New("released")
	}), 1)
	require.NoError(t, err)

	firstEdge, firstConnector := net.Pipe()
	firstDone := make(chan error, 1)
	go func() { firstDone <- handler.handle(context.Background(), firstConnector) }()
	firstFlowID := uuid.Must(uuid.NewV7()).String()
	_, err = transport.WriteFrame(firstEdge, transport.FrameTypeFlowOpen, transport.FlowOpen{
		FlowID:                  firstFlowID,
		TargetName:              "database",
		DialTimeoutMilliseconds: 10_000,
	})
	require.NoError(t, err)
	select {
	case <-dialStarted:
	case <-time.After(time.Second):
		t.Fatal("first flow did not start dialing")
	}

	secondFlowID := uuid.Must(uuid.NewV7()).String()
	secondResult, secondErr := exchangeRejectedFlow(t, handler, transport.FlowOpen{
		FlowID:                  secondFlowID,
		TargetName:              "database",
		DialTimeoutMilliseconds: 1_000,
	})
	assert.Equal(t, transport.FlowResult{
		FlowID: secondFlowID,
		Status: transport.FlowStatusCapacityExceeded,
	}, secondResult)
	require.NoError(t, secondErr)

	close(releaseDial)
	var firstResult transport.FlowResult
	require.NoError(t, transport.ReadFrame(firstEdge, transport.FrameTypeFlowResult, &firstResult))
	assert.Equal(t, transport.FlowStatusDialFailed, firstResult.Status)
	require.NoError(t, firstEdge.Close())
	require.NoError(t, <-firstDone)
}

func TestFlowHandlerDoesNotBlockLaterCapacityRejections(t *testing.T) {
	t.Parallel()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	handler, err := newFlowHandler(targets, &net.Dialer{}, 1)
	require.NoError(t, err)
	require.True(t, handler.tryAcquire())
	defer handler.release()

	ctx, cancel := context.WithCancel(context.Background())
	var flows sync.WaitGroup
	stalledEdge, stalledConnector := net.Pipe()
	laterEdge, laterConnector := net.Pipe()
	t.Cleanup(func() {
		cancel()
		_ = stalledEdge.Close()
		_ = stalledConnector.Close()
		_ = laterEdge.Close()
		_ = laterConnector.Close()
	})
	dispatchDone := make(chan struct{})
	go func() {
		handler.dispatch(ctx, deadlineIgnoringConn{Conn: stalledConnector}, &flows)
		handler.dispatch(ctx, laterConnector, &flows)
		close(dispatchDone)
	}()

	flowID := uuid.Must(uuid.NewV7()).String()
	writeDone := make(chan error, 1)
	go func() {
		_, writeErr := transport.WriteFrame(laterEdge, transport.FrameTypeFlowOpen, transport.FlowOpen{
			FlowID:                  flowID,
			TargetName:              "database",
			DialTimeoutMilliseconds: 1_000,
		})
		writeDone <- writeErr
	}()
	require.NoError(t, laterEdge.SetReadDeadline(time.Now().Add(time.Second)))
	var result transport.FlowResult
	resultErr := transport.ReadFrame(laterEdge, transport.FrameTypeFlowResult, &result)

	cancel()
	require.NoError(t, stalledEdge.Close())
	require.NoError(t, laterEdge.Close())
	select {
	case <-dispatchDone:
	case <-time.After(time.Second):
		t.Fatal("flow dispatch remained blocked during cleanup")
	}
	flows.Wait()

	require.NoError(t, resultErr)
	require.NoError(t, <-writeDone)
	assert.Equal(t, transport.FlowResult{
		FlowID: flowID,
		Status: transport.FlowStatusCapacityExceeded,
	}, result)
}

func TestFlowHandlerWritesReadyBeforeOpaqueBytesAndPreservesHalfClose(t *testing.T) {
	t.Parallel()

	targetListener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	t.Cleanup(func() { _ = targetListener.Close() })
	targetDone := make(chan error, 1)
	go func() {
		conn, acceptErr := targetListener.AcceptTCP()
		if acceptErr != nil {
			targetDone <- acceptErr
			return
		}
		defer conn.Close()
		request, readErr := io.ReadAll(conn)
		if readErr != nil {
			targetDone <- readErr
			return
		}
		if string(request) != "request" {
			targetDone <- errors.New("target received unexpected request")
			return
		}
		_, writeErr := conn.Write([]byte("response-after-eof"))
		targetDone <- writeErr
	}()

	targets, err := NewTargetMap(map[string]string{"database": targetListener.Addr().String()})
	require.NoError(t, err)
	handler, err := newFlowHandler(targets, &net.Dialer{}, 1)
	require.NoError(t, err)
	edgeSide, connectorSide := tcpConnectionPair(t)
	done := make(chan error, 1)
	go func() { done <- handler.handle(context.Background(), connectorSide) }()

	flowID := uuid.Must(uuid.NewV7()).String()
	_, err = transport.WriteFrame(edgeSide, transport.FrameTypeFlowOpen, transport.FlowOpen{
		FlowID:                  flowID,
		TargetName:              "database",
		DialTimeoutMilliseconds: 100,
	})
	require.NoError(t, err)
	var result transport.FlowResult
	require.NoError(t, transport.ReadFrame(edgeSide, transport.FrameTypeFlowResult, &result))
	assert.Equal(t, transport.FlowResult{FlowID: flowID, Status: transport.FlowStatusReady}, result)

	require.NoError(t, edgeSide.SetDeadline(time.Now().Add(time.Second)))
	_, err = edgeSide.Write([]byte("request"))
	require.NoError(t, err)
	require.NoError(t, edgeSide.CloseWrite())
	response, err := io.ReadAll(edgeSide)
	require.NoError(t, err)
	assert.Equal(t, "response-after-eof", string(response))
	require.NoError(t, <-targetDone)
	require.NoError(t, <-done)
}

func TestFlowHandlerLogsHowAnHTTPFlowEnded(t *testing.T) {
	t.Parallel()

	const request = "POST /v1/responses HTTP/1.1\r\nHost: api\r\nContent-Length: 2\r\n\r\nhi"
	const response = "HTTP/1.1 200 OK\r\nContent-Length: 2\r\nConnection: close\r\n\r\nok"
	const partialHeader = "HTTP/1.1 200 OK\r\nContent-Length:"
	answer := func(conn net.Conn) {
		if _, err := http.ReadRequest(bufio.NewReader(conn)); err == nil {
			_, _ = io.WriteString(conn, response)
		}
	}
	tests := []struct {
		name string
		// origin serves the one request; returning closes its connection.
		origin func(net.Conn)
		// edgeReads is the response the edge reads before it closes the flow;
		// when empty the edge gives up right after sending.
		edgeReads      string
		wantUnanswered bool
	}{
		{
			name:      "target answers while the edge waits",
			origin:    answer,
			edgeReads: response,
		},
		{
			name: "target sends part of its header and stalls",
			origin: func(conn net.Conn) {
				if _, err := http.ReadRequest(bufio.NewReader(conn)); err == nil {
					_, _ = io.WriteString(conn, partialHeader)
					_, _ = io.Copy(io.Discard, conn)
				}
			},
			edgeReads:      partialHeader,
			wantUnanswered: true,
		},
		{
			name:           "edge gives up before the target answers",
			origin:         func(conn net.Conn) { _, _ = io.Copy(io.Discard, conn) },
			wantUnanswered: true,
		},
		{
			name: "target answers only after the edge gave up",
			origin: func(conn net.Conn) {
				_, _ = io.Copy(io.Discard, conn)
				_, _ = io.WriteString(conn, response)
			},
			wantUnanswered: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			originListener, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			t.Cleanup(func() { _ = originListener.Close() })
			go func() {
				conn, acceptErr := originListener.Accept()
				if acceptErr != nil {
					return
				}
				defer conn.Close()
				tt.origin(conn)
			}()

			targets, err := NewTargetMap(map[string]string{"api": "http://" + originListener.Addr().String()})
			require.NoError(t, err)
			handler, err := newFlowHandler(targets, &net.Dialer{}, 1)
			require.NoError(t, err)
			logs := &recordingHandler{}
			handler.logger = slog.New(logs)
			edgeSide, connectorSide := tcpConnectionPair(t)
			done := make(chan error, 1)
			go func() { done <- handler.handle(context.Background(), connectorSide) }()

			_, err = transport.WriteFrame(edgeSide, transport.FrameTypeFlowOpen, transport.FlowOpen{
				FlowID:                  uuid.Must(uuid.NewV7()).String(),
				TargetName:              "api",
				DialTimeoutMilliseconds: 1_000,
			})
			require.NoError(t, err)
			var result transport.FlowResult
			require.NoError(t, transport.ReadFrame(edgeSide, transport.FrameTypeFlowResult, &result))
			require.Equal(t, transport.FlowStatusReady, result.Status)
			require.NoError(t, edgeSide.SetDeadline(time.Now().Add(5*time.Second)))
			_, err = io.WriteString(edgeSide, request)
			require.NoError(t, err)
			if tt.edgeReads != "" {
				received := make([]byte, len(tt.edgeReads))
				_, readErr := io.ReadFull(edgeSide, received)
				require.NoError(t, readErr)
				require.Equal(t, tt.edgeReads, string(received))
			}
			require.NoError(t, edgeSide.Close())
			select {
			case <-done:
			case <-time.After(5 * time.Second):
				t.Fatal("flow did not end after the edge closed it")
			}

			unanswered, warned := logs.find("target did not respond before the flow closed")
			if !tt.wantUnanswered {
				assert.False(t, warned, "an answered flow was reported as unanswered")
				closed, ok := logs.find("connector closed flow to target")
				require.True(t, ok, "flow end was not logged")
				assert.Equal(t, slog.LevelDebug, closed.level)
				assert.Equal(t, strconv.Itoa(len(response)), closed.attrs["bytes_from_target"])
				return
			}
			require.True(t, warned, "unanswered flow was not logged")
			assert.Equal(t, slog.LevelWarn, unanswered.level)
			assert.Equal(t, "api", unanswered.attrs["target"])
			assert.Equal(t, strconv.Itoa(len(request)), unanswered.attrs["bytes_to_target"])
			assert.NotEmpty(t, unanswered.attrs["duration"])
		})
	}
}

func exchangeRejectedFlow(
	t *testing.T,
	handler *flowHandler,
	open transport.FlowOpen,
) (transport.FlowResult, error) {
	t.Helper()

	edgeSide, connectorSide := net.Pipe()
	done := make(chan error, 1)
	go func() { done <- handler.handle(context.Background(), connectorSide) }()
	_, err := transport.WriteFrame(edgeSide, transport.FrameTypeFlowOpen, open)
	require.NoError(t, err)
	var result transport.FlowResult
	require.NoError(t, transport.ReadFrame(edgeSide, transport.FrameTypeFlowResult, &result))
	require.NoError(t, edgeSide.Close())
	return result, <-done
}

func tcpConnectionPair(t *testing.T) (*net.TCPConn, *net.TCPConn) {
	t.Helper()

	listener, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.IPv4(127, 0, 0, 1)})
	require.NoError(t, err)
	accepted := make(chan *net.TCPConn, 1)
	acceptErr := make(chan error, 1)
	go func() {
		conn, err := listener.AcceptTCP()
		if err != nil {
			acceptErr <- err
			return
		}
		accepted <- conn
	}()
	edgeSide, err := net.DialTCP("tcp", nil, listener.Addr().(*net.TCPAddr))
	require.NoError(t, err)
	select {
	case connectorSide := <-accepted:
		require.NoError(t, listener.Close())
		t.Cleanup(func() { _ = edgeSide.Close() })
		return edgeSide, connectorSide
	case err := <-acceptErr:
		require.NoError(t, err)
	case <-time.After(time.Second):
		t.Fatal("timed out accepting loopback TCP connection")
	}
	return nil, nil
}
