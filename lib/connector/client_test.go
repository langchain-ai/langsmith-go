package connector

import (
	"context"
	"errors"
	"fmt"
	"io"
	"log/slog"
	"net"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"slices"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"github.com/hashicorp/yamux"
	"github.com/langchain-ai/langsmith-go/lib/connector/transport"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type carrierObservation struct {
	headers  http.Header
	register transport.Register
	payload  string
}

func TestClientCarrierRegistersAndAcceptsEdgeOpenedFlows(t *testing.T) {
	t.Parallel()

	connectionID := uuid.Must(uuid.NewV7())
	runtimeID := uuid.Must(uuid.NewV7())
	flowID := uuid.Must(uuid.NewV7()).String()
	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)

	observed := make(chan carrierObservation, 1)
	edgeErr := make(chan error, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin:  func(*http.Request) bool { return true },
		Subprotocols: []string{transport.WebSocketSubprotocol},
	}
	edgeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if err := validateCarrierRequest(r); err != nil {
			edgeErr <- err
			http.Error(w, "invalid carrier request", http.StatusBadRequest)
			return
		}
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			edgeErr <- fmt.Errorf("upgrade carrier: %w", err)
			return
		}
		defer ws.Close()
		if ws.Subprotocol() != transport.WebSocketSubprotocol {
			edgeErr <- fmt.Errorf("negotiated subprotocol %q", ws.Subprotocol())
			return
		}
		session, err := yamux.Server(transport.NewWebSocketConn(ws), transport.YamuxConfig())
		if err != nil {
			edgeErr <- fmt.Errorf("start edge yamux: %w", err)
			return
		}
		defer session.Close()

		control, err := session.AcceptStream()
		if err != nil {
			edgeErr <- fmt.Errorf("accept control stream: %w", err)
			return
		}
		defer control.Close()
		var register transport.Register
		if err := transport.ReadFrame(control, transport.FrameTypeRegister, &register); err != nil {
			edgeErr <- fmt.Errorf("read register: %w", err)
			return
		}
		if err := transport.ValidateRegister(register); err != nil {
			edgeErr <- fmt.Errorf("validate register: %w", err)
			return
		}
		if err := registerMatchesHeaders(register, r.Header); err != nil {
			edgeErr <- err
			return
		}
		if _, err := transport.WriteFrame(control, transport.FrameTypeRegistered, transport.Registered{
			ProtocolVersion:     transport.ProtocolVersion,
			ConnectionID:        register.ConnectionID,
			ConnectorInstanceID: register.ConnectorInstanceID,
			SessionID:           register.SessionID,
			RuntimeID:           runtimeID.String(),
		}); err != nil {
			edgeErr <- fmt.Errorf("write registered: %w", err)
			return
		}

		unexpectedControl := make(chan error, 1)
		go func() {
			stream, err := session.AcceptStream()
			if err == nil {
				_ = stream.Close()
				unexpectedControl <- errors.New("connector opened more than one control stream")
			}
		}()
		select {
		case err := <-unexpectedControl:
			edgeErr <- err
			return
		case <-time.After(25 * time.Millisecond):
		}

		flow, err := session.OpenStream()
		if err != nil {
			edgeErr <- fmt.Errorf("open application stream: %w", err)
			return
		}
		defer flow.Close()
		if _, err := transport.WriteFrame(flow, transport.FrameTypeFlowOpen, transport.FlowOpen{
			FlowID:                  flowID,
			TargetName:              "database",
			DialTimeoutMilliseconds: 1_000,
		}); err != nil {
			edgeErr <- fmt.Errorf("write flow open: %w", err)
			return
		}
		var result transport.FlowResult
		if err := transport.ReadFrame(flow, transport.FrameTypeFlowResult, &result); err != nil {
			edgeErr <- fmt.Errorf("read flow result: %w", err)
			return
		}
		if result != (transport.FlowResult{FlowID: flowID, Status: transport.FlowStatusReady}) {
			edgeErr <- fmt.Errorf("unexpected flow result: %#v", result)
			return
		}
		if _, err := flow.Write([]byte("ping")); err != nil {
			edgeErr <- fmt.Errorf("write opaque payload: %w", err)
			return
		}
		response := make([]byte, 4)
		if _, err := io.ReadFull(flow, response); err != nil {
			edgeErr <- fmt.Errorf("read opaque payload: %w", err)
			return
		}
		observed <- carrierObservation{
			headers:  r.Header.Clone(),
			register: register,
			payload:  string(response),
		}
		<-session.CloseChan()
	}))
	defer edgeServer.Close()

	apiKeyFile := filepath.Join(t.TempDir(), "api-key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte("test-api-key\n"), 0o600))
	config := Config{
		Endpoint:     "wss" + strings.TrimPrefix(edgeServer.URL, "https") + "/v1/connections/connect",
		ConnectionID: connectionID,
		Targets:      targets,
		apiKeyFile:   apiKeyFile,
	}
	var targetDials atomic.Int32
	targetDialer := targetDialerFunc(func(
		_ context.Context,
		network string,
		address string,
	) (net.Conn, error) {
		if network != "tcp" || address != "db.internal:5432" {
			return nil, fmt.Errorf("unexpected dial %q %q", network, address)
		}
		targetDials.Add(1)
		connectorSide, targetSide := net.Pipe()
		go func() {
			defer targetSide.Close()
			_, _ = io.Copy(targetSide, targetSide)
		}()
		return connectorSide, nil
	})
	tlsConfig := edgeServer.Client().Transport.(*http.Transport).TLSClientConfig.Clone()
	client, err := NewClient(
		config,
		WithWebSocketDialer(&websocket.Dialer{TLSClientConfig: tlsConfig}),
		WithTargetDialer(targetDialer),
		WithDrainGracePeriod(50*time.Millisecond),
	)
	require.NoError(t, err)

	ctx, cancel := context.WithCancel(context.Background())
	runDone := make(chan error, 1)
	go func() { runDone <- runSingleCarrier(ctx, client) }()
	var observation carrierObservation
	select {
	case observation = <-observed:
	case err := <-edgeErr:
		t.Fatalf("fake edge failed: %v", err)
	case err := <-runDone:
		t.Fatalf("connector stopped before flow: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for connector flow")
	}

	assert.Equal(t, "test-api-key", observation.headers.Get("X-Api-Key"))
	assert.Equal(t, connectionID.String(), observation.headers.Get(transport.HeaderConnectionID))
	assert.Equal(t, observation.headers.Get(transport.HeaderConnectorInstanceID), observation.register.ConnectorInstanceID)
	assert.Equal(t, observation.headers.Get(transport.HeaderSessionID), observation.register.SessionID)
	assert.Equal(t, connectionID.String(), observation.register.ConnectionID)
	assert.Equal(t, targets.RegistrationMetadata().Targets, observation.register.Targets)
	assert.Equal(t, targets.Hash(), observation.register.TargetMapHash)
	assert.Equal(t, "ping", observation.payload)
	assert.EqualValues(t, 1, targetDials.Load())
	assertUUIDv7(t, observation.register.ConnectorInstanceID)
	assertUUIDv7(t, observation.register.SessionID)

	cancel()
	select {
	case err := <-runDone:
		require.ErrorIs(t, err, context.Canceled)
	case err := <-edgeErr:
		t.Fatalf("fake edge failed during shutdown: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("connector did not stop after cancellation")
	}
}

func TestClientBoundsPendingEdgeOpenedStreams(t *testing.T) {
	t.Parallel()

	connectionID := uuid.Must(uuid.NewV7())
	runtimeID := uuid.Must(uuid.NewV7())
	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	excessClosed := make(chan error, 1)
	edgeErr := make(chan error, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin:  func(*http.Request) bool { return true },
		Subprotocols: []string{transport.WebSocketSubprotocol},
	}
	edgeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, upgradeErr := upgrader.Upgrade(w, r, nil)
		if upgradeErr != nil {
			edgeErr <- upgradeErr
			return
		}
		defer ws.Close()
		session, sessionErr := yamux.Server(transport.NewWebSocketConn(ws), transport.YamuxConfig())
		if sessionErr != nil {
			edgeErr <- sessionErr
			return
		}
		defer session.Close()
		control, acceptErr := session.AcceptStream()
		if acceptErr != nil {
			edgeErr <- acceptErr
			return
		}
		defer control.Close()
		var register transport.Register
		if readErr := transport.ReadFrame(control, transport.FrameTypeRegister, &register); readErr != nil {
			edgeErr <- readErr
			return
		}
		if _, writeErr := transport.WriteFrame(control, transport.FrameTypeRegistered, transport.Registered{
			ProtocolVersion:     transport.ProtocolVersion,
			ConnectionID:        register.ConnectionID,
			ConnectorInstanceID: register.ConnectorInstanceID,
			SessionID:           register.SessionID,
			RuntimeID:           runtimeID.String(),
		}); writeErr != nil {
			edgeErr <- writeErr
			return
		}

		streams := make([]*yamux.Stream, 0, 4)
		for range 4 {
			stream, openErr := session.OpenStream()
			if openErr != nil {
				edgeErr <- openErr
				return
			}
			streams = append(streams, stream)
			defer stream.Close()
		}
		flowID := uuid.Must(uuid.NewV7()).String()
		if _, writeErr := transport.WriteFrame(streams[2], transport.FrameTypeFlowOpen, transport.FlowOpen{
			FlowID:                  flowID,
			TargetName:              "database",
			DialTimeoutMilliseconds: 1_000,
		}); writeErr != nil {
			excessClosed <- fmt.Errorf("write overloaded FLOW_OPEN: %w", writeErr)
			return
		}
		var result transport.FlowResult
		if readErr := transport.ReadFrame(streams[2], transport.FrameTypeFlowResult, &result); readErr != nil {
			excessClosed <- fmt.Errorf("read overloaded FLOW_RESULT: %w", readErr)
			return
		}
		if result != (transport.FlowResult{FlowID: flowID, Status: transport.FlowStatusCapacityExceeded}) {
			excessClosed <- fmt.Errorf("overloaded result = %#v", result)
			return
		}

		_ = streams[3].SetReadDeadline(time.Now().Add(500 * time.Millisecond))
		var byteBuffer [1]byte
		_, readErr := streams[3].Read(byteBuffer[:])
		var netErr net.Error
		if errors.As(readErr, &netErr) && netErr.Timeout() {
			excessClosed <- errors.New("excess stream remained parked")
			return
		}
		if readErr == nil {
			excessClosed <- errors.New("excess stream unexpectedly carried data")
			return
		}
		excessClosed <- nil
		<-session.CloseChan()
	}))
	defer edgeServer.Close()

	apiKeyFile := filepath.Join(t.TempDir(), "api-key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte("test-api-key\n"), 0o600))
	config := Config{
		Endpoint:     "wss" + strings.TrimPrefix(edgeServer.URL, "https") + "/v1/connections/connect",
		ConnectionID: connectionID,
		Targets:      targets,
		apiKeyFile:   apiKeyFile,
	}
	var targetDials atomic.Int32
	client, err := NewClient(
		config,
		WithTLSConfig(edgeServer.Client().Transport.(*http.Transport).TLSClientConfig.Clone()),
		WithTargetDialer(targetDialerFunc(func(context.Context, string, string) (net.Conn, error) {
			targetDials.Add(1)
			return nil, errors.New("must not dial without FLOW_OPEN")
		})),
		WithMaxConcurrentFlows(2),
		WithDrainGracePeriod(50*time.Millisecond),
	)
	require.NoError(t, err)
	ctx, cancel := context.WithCancel(context.Background())
	runDone := make(chan error, 1)
	go func() { runDone <- runSingleCarrier(ctx, client) }()

	select {
	case err := <-excessClosed:
		require.NoError(t, err)
	case err := <-edgeErr:
		t.Fatalf("fake edge failed: %v", err)
	case err := <-runDone:
		t.Fatalf("connector stopped before enforcing capacity: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for excess stream closure")
	}
	assert.Zero(t, targetDials.Load())
	cancel()
	select {
	case err := <-runDone:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(3 * time.Second):
		t.Fatal("connector did not stop after cancellation")
	}
}

func TestClientRunDrainsFlowRacingShutdown(t *testing.T) {
	t.Parallel()

	const (
		drainGracePeriod = 500 * time.Millisecond
		lateFlowDelay    = 200 * time.Millisecond
	)
	carrierReady := make(chan struct{})
	drainObserved := make(chan struct{})
	racedFlowServed := make(chan struct{})
	edgeErr := make(chan error, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin:  func(*http.Request) bool { return true },
		Subprotocols: []string{transport.WebSocketSubprotocol},
	}
	edgeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			edgeErr <- fmt.Errorf("upgrade carrier: %w", err)
			return
		}
		defer ws.Close()
		session, err := yamux.Server(transport.NewWebSocketConn(ws), transport.YamuxConfig())
		if err != nil {
			edgeErr <- fmt.Errorf("start edge yamux: %w", err)
			return
		}
		defer session.Close()
		control, err := session.AcceptStream()
		if err != nil {
			edgeErr <- fmt.Errorf("accept control stream: %w", err)
			return
		}
		defer control.Close()

		var register transport.Register
		if err := transport.ReadFrame(control, transport.FrameTypeRegister, &register); err != nil {
			edgeErr <- fmt.Errorf("read register: %w", err)
			return
		}
		if _, err := transport.WriteFrame(control, transport.FrameTypeRegistered, transport.Registered{
			ProtocolVersion:     transport.ProtocolVersion,
			ConnectionID:        register.ConnectionID,
			ConnectorInstanceID: register.ConnectorInstanceID,
			SessionID:           register.SessionID,
			RuntimeID:           uuid.Must(uuid.NewV7()).String(),
		}); err != nil {
			edgeErr <- fmt.Errorf("write registered: %w", err)
			return
		}
		flow, err := session.OpenStream()
		if err != nil {
			edgeErr <- fmt.Errorf("open readiness flow: %w", err)
			return
		}
		flowID := uuid.Must(uuid.NewV7()).String()
		if _, err := transport.WriteFrame(flow, transport.FrameTypeFlowOpen, transport.FlowOpen{
			FlowID:                  flowID,
			TargetName:              "unknown",
			DialTimeoutMilliseconds: 1_000,
		}); err != nil {
			edgeErr <- fmt.Errorf("write readiness flow: %w", err)
			return
		}
		var result transport.FlowResult
		if err := transport.ReadFrame(flow, transport.FrameTypeFlowResult, &result); err != nil {
			edgeErr <- fmt.Errorf("read readiness flow result: %w", err)
			return
		}
		if result != (transport.FlowResult{FlowID: flowID, Status: transport.FlowStatusUnknownTarget}) {
			edgeErr <- fmt.Errorf("unexpected readiness flow result: %#v", result)
			return
		}
		_ = flow.Close()
		close(carrierReady)

		var drain transport.Drain
		if err := transport.ReadFrame(control, transport.FrameTypeDrain, &drain); err != nil {
			edgeErr <- fmt.Errorf("read drain: %w", err)
			return
		}
		if err := transport.ValidateDrain(drain); err != nil {
			edgeErr <- fmt.Errorf("validate drain: %w", err)
			return
		}
		close(drainObserved)
		time.Sleep(lateFlowDelay)

		flow, err = session.OpenStream()
		if err != nil {
			edgeErr <- fmt.Errorf("open raced drain flow: %w", err)
			return
		}
		flowID = uuid.Must(uuid.NewV7()).String()
		if _, err := transport.WriteFrame(flow, transport.FrameTypeFlowOpen, transport.FlowOpen{
			FlowID:                  flowID,
			TargetName:              "unknown",
			DialTimeoutMilliseconds: 1_000,
		}); err != nil {
			edgeErr <- fmt.Errorf("write raced drain flow: %w", err)
			return
		}
		result = transport.FlowResult{}
		if err := transport.ReadFrame(flow, transport.FrameTypeFlowResult, &result); err != nil {
			edgeErr <- fmt.Errorf("read raced drain flow result: %w", err)
			return
		}
		if result != (transport.FlowResult{FlowID: flowID, Status: transport.FlowStatusUnknownTarget}) {
			edgeErr <- fmt.Errorf("unexpected raced drain flow result: %#v", result)
			return
		}
		_ = flow.Close()
		close(racedFlowServed)
		<-session.CloseChan()
	}))
	defer edgeServer.Close()

	client := newTestClient(t, edgeServer, "test-api-key")
	client.drainGracePeriod = drainGracePeriod
	ctx, cancel := context.WithCancel(context.Background())
	runDone := make(chan error, 1)
	go func() { runDone <- runSingleCarrier(ctx, client) }()

	select {
	case <-carrierReady:
	case err := <-edgeErr:
		t.Fatalf("fake edge failed during registration: %v", err)
	case <-time.After(3 * time.Second):
		t.Fatal("timed out waiting for carrier readiness")
	}
	require.Eventually(t, func() bool {
		return len(client.flowHandler.slots) == 0
	}, time.Second, time.Millisecond)
	cancelledAt := time.Now()
	cancel()
	select {
	case <-drainObserved:
	case err := <-edgeErr:
		t.Fatalf("fake edge failed during drain: %v", err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for DRAIN")
	}
	select {
	case <-racedFlowServed:
	case err := <-edgeErr:
		t.Fatalf("fake edge failed while racing drain: %v", err)
	case err := <-runDone:
		t.Fatalf("connector stopped before serving a flow already racing DRAIN: %v", err)
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for a flow already racing DRAIN")
	}
	select {
	case err := <-runDone:
		require.ErrorIs(t, err, context.Canceled)
		assert.Less(t, time.Since(cancelledAt), time.Second)
	case err := <-edgeErr:
		t.Fatalf("fake edge failed during shutdown: %v", err)
	case <-time.After(2 * time.Second):
		t.Fatal("connector did not stop after the drain grace period")
	}
}

func TestClientRejectsRegisteredIdentityMismatch(t *testing.T) {
	t.Parallel()

	edgeErr := make(chan error, 1)
	upgrader := websocket.Upgrader{
		CheckOrigin:  func(*http.Request) bool { return true },
		Subprotocols: []string{transport.WebSocketSubprotocol},
	}
	edgeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			edgeErr <- err
			return
		}
		defer ws.Close()
		session, err := yamux.Server(transport.NewWebSocketConn(ws), transport.YamuxConfig())
		if err != nil {
			edgeErr <- err
			return
		}
		defer session.Close()
		control, err := session.AcceptStream()
		if err != nil {
			edgeErr <- err
			return
		}
		defer control.Close()
		var register transport.Register
		if err := transport.ReadFrame(control, transport.FrameTypeRegister, &register); err != nil {
			edgeErr <- err
			return
		}
		_, err = transport.WriteFrame(control, transport.FrameTypeRegistered, transport.Registered{
			ProtocolVersion:     transport.ProtocolVersion,
			ConnectionID:        register.ConnectionID,
			ConnectorInstanceID: register.ConnectorInstanceID,
			SessionID:           uuid.Must(uuid.NewV7()).String(),
			RuntimeID:           uuid.Must(uuid.NewV7()).String(),
		})
		if err != nil {
			edgeErr <- err
		}
	}))
	defer edgeServer.Close()

	client := newTestClient(t, edgeServer, "test-api-key")
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	err := runSingleCarrier(ctx, client)
	require.ErrorContains(t, err, "identity mismatch")
	select {
	case err := <-edgeErr:
		t.Fatalf("fake edge failed: %v", err)
	default:
	}
}

func TestClientReadsAPIKeyFileForEveryCarrierAttempt(t *testing.T) {
	t.Parallel()

	apiKeys := make(chan string, 2)
	edgeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKeys <- r.Header.Get("X-Api-Key")
		http.Error(w, "carrier unavailable", http.StatusServiceUnavailable)
	}))
	defer edgeServer.Close()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	apiKeyFile := filepath.Join(t.TempDir(), "api-key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte("first-key\n"), 0o600))
	config := Config{
		Endpoint:     "wss" + strings.TrimPrefix(edgeServer.URL, "https") + "/v1/connections/connect",
		ConnectionID: uuid.Must(uuid.NewV7()),
		Targets:      targets,
		apiKeyFile:   apiKeyFile,
	}
	client, err := NewClient(
		config,
		WithTLSConfig(edgeServer.Client().Transport.(*http.Transport).TLSClientConfig),
	)
	require.NoError(t, err)

	require.Error(t, runSingleCarrier(context.Background(), client))
	require.NoError(t, os.WriteFile(apiKeyFile, []byte("second-key\n"), 0o600))
	require.Error(t, runSingleCarrier(context.Background(), client))
	assert.Equal(t, "first-key", <-apiKeys)
	assert.Equal(t, "second-key", <-apiKeys)
}

func newTestClient(
	t *testing.T,
	edgeServer *httptest.Server,
	apiKey string,
	options ...ClientOption,
) *Client {
	t.Helper()

	targets, err := NewTargetMap(map[string]string{"database": "db.internal:5432"})
	require.NoError(t, err)
	apiKeyFile := filepath.Join(t.TempDir(), "api-key")
	require.NoError(t, os.WriteFile(apiKeyFile, []byte(apiKey), 0o600))
	client, err := NewClient(
		Config{
			Endpoint:     "wss" + strings.TrimPrefix(edgeServer.URL, "https") + "/v1/connections/connect",
			ConnectionID: uuid.Must(uuid.NewV7()),
			Targets:      targets,
			apiKeyFile:   apiKeyFile,
		},
		append([]ClientOption{
			WithTLSConfig(edgeServer.Client().Transport.(*http.Transport).TLSClientConfig),
		}, options...)...,
	)
	require.NoError(t, err)
	return client
}

func TestClientRunLogsWhyItCannotRegister(t *testing.T) {
	t.Parallel()

	edgeServer := httptest.NewTLSServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		http.Error(w, "missing connections:connect", http.StatusForbidden)
	}))
	t.Cleanup(edgeServer.Close)
	logs := &recordingHandler{}
	client := newTestClient(t, edgeServer, "lsv2_test_key", WithLogger(slog.New(logs)))

	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() { done <- client.Run(ctx) }()
	var record recordedLog
	require.Eventually(t, func() bool {
		var ok bool
		record, ok = logs.find("connector could not register with LangSmith; retrying")
		return ok
	}, 5*time.Second, 10*time.Millisecond)
	cancel()
	require.ErrorIs(t, <-done, context.Canceled)

	assert.Equal(t, slog.LevelWarn, record.level)
	assert.Contains(t, record.attrs["error"], "HTTP 403")
	assert.NotContains(t, record.attrs["error"], "lsv2_test_key")
}

func runSingleCarrier(ctx context.Context, client *Client) error {
	carrier, err := client.runCarrier(ctx)
	if err != nil {
		return err
	}
	return client.serveCarrier(ctx, carrier)
}

func validateCarrierRequest(r *http.Request) error {
	if r.Method != http.MethodGet {
		return fmt.Errorf("carrier method = %s", r.Method)
	}
	if r.URL.Path != "/v1/connections/connect" || r.URL.RawQuery != "" {
		return fmt.Errorf("carrier URL = %s", r.URL.String())
	}
	if r.Header.Get("X-Api-Key") != "test-api-key" {
		return errors.New("missing API key")
	}
	if !slices.Equal(websocket.Subprotocols(r), []string{transport.WebSocketSubprotocol}) {
		return fmt.Errorf("offered subprotocols = %#v", websocket.Subprotocols(r))
	}
	for _, name := range []string{
		transport.HeaderConnectionID,
		transport.HeaderConnectorInstanceID,
		transport.HeaderSessionID,
	} {
		if len(r.Header.Values(name)) != 1 {
			return fmt.Errorf("header %s values = %#v", name, r.Header.Values(name))
		}
	}
	return nil
}

func registerMatchesHeaders(register transport.Register, headers http.Header) error {
	for _, pair := range []struct {
		name  string
		value string
	}{
		{name: transport.HeaderConnectionID, value: register.ConnectionID},
		{name: transport.HeaderConnectorInstanceID, value: register.ConnectorInstanceID},
		{name: transport.HeaderSessionID, value: register.SessionID},
	} {
		if headers.Get(pair.name) != pair.value {
			return fmt.Errorf("header %s does not match REGISTER", pair.name)
		}
	}
	return nil
}

func assertUUIDv7(t *testing.T, value string) {
	t.Helper()
	id, err := uuid.Parse(value)
	require.NoError(t, err)
	assert.Equal(t, uuid.Version(7), id.Version())
}
