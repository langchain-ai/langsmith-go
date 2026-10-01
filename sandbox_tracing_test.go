package langsmith

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"sync/atomic"
	"testing"

	"github.com/langchain-ai/langsmith-go/option"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
	"golang.org/x/net/websocket"
)

const tracingSandboxID = "7ca3b03e-2a49-4dd6-aa52-f4f578a4ea9c"

func sandboxTraceTestContext(t *testing.T) (context.Context, trace.Span, *tracetest.SpanRecorder) {
	t.Helper()
	recorder := tracetest.NewSpanRecorder()
	provider := sdktrace.NewTracerProvider(sdktrace.WithSpanProcessor(recorder))
	t.Cleanup(func() { require.NoError(t, provider.Shutdown(context.Background())) })
	ctx, span := provider.Tracer("sandbox-test").Start(context.Background(), "caller")
	span.SetAttributes(attribute.String("existing", "preserved"))
	return ctx, span, recorder
}

func requireSandboxTrace(t *testing.T, span trace.Span, recorder *tracetest.SpanRecorder, id string) {
	t.Helper()
	span.End()
	spans := recorder.Ended()
	require.Len(t, spans, 1)
	attrs := attribute.NewSet(spans[0].Attributes()...)
	existing, ok := attrs.Value("existing")
	require.True(t, ok)
	require.Equal(t, "preserved", existing.AsString())
	sandboxID, ok := attrs.Value("langsmith.metadata.sandbox_id")
	if id == "" {
		require.False(t, ok)
	} else {
		require.True(t, ok)
		require.Equal(t, id, sandboxID.AsString())
	}
}

func TestSandboxTracingMetadata(t *testing.T) {
	ctx, span, recorder := sandboxTraceTestContext(t)
	traceSandbox(ctx, "first")
	traceSandbox(ctx, "")
	traceSandbox(ctx, tracingSandboxID)
	requireSandboxTrace(t, span, recorder, tracingSandboxID)
	require.NotPanics(t, func() { traceSandbox(context.Background(), tracingSandboxID) })
}

func TestSandboxTracingReferences(t *testing.T) {
	for _, ref := range []string{"display-name", tracingSandboxID} {
		t.Run(ref, func(t *testing.T) {
			ctx, span, recorder := sandboxTraceTestContext(t)
			traceSandboxReference(ctx, ref)
			want := ""
			if ref == tracingSandboxID {
				want = ref
			}
			requireSandboxTrace(t, span, recorder, want)
		})
	}
}

func TestSandboxTracingConvenienceCalls(t *testing.T) {
	var requests atomic.Int32
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		requests.Add(1)
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"id": tracingSandboxID, "name": "display-name", "status": "ready",
			"dataplane_url": srv.URL, "source_sandbox_id": tracingSandboxID,
			"service_url": srv.URL, "browser_url": srv.URL,
			"items": []map[string]string{{"id": tracingSandboxID, "name": "display-name"}},
		})
	}))
	defer srv.Close()
	boxes := NewSandboxBoxService(option.WithBaseURL(srv.URL), option.WithAPIKey("test-key"), option.WithMaxRetries(0))
	sandbox := &Sandbox{ID: tracingSandboxID, Name: "display-name", DataplaneURL: srv.URL, boxes: boxes}
	tests := []struct {
		name     string
		call     func(context.Context) error
		requests int32
		want     string
	}{
		{"create", func(ctx context.Context) error { _, err := boxes.NewSandbox(ctx, SandboxBoxNewParams{}); return err }, 1, tracingSandboxID},
		{"get", func(ctx context.Context) error { _, err := boxes.GetSandbox(ctx, sandbox.Name); return err }, 1, tracingSandboxID},
		{"list", func(ctx context.Context) error {
			_, err := boxes.ListSandboxes(ctx, SandboxBoxListParams{})
			return err
		}, 1, ""},
		{"run", func(ctx context.Context) error {
			_, err := sandbox.Run(ctx, SandboxBoxRunParams{Command: String("pwd")})
			return err
		}, 1, tracingSandboxID},
		{"run by name", func(ctx context.Context) error {
			_, err := boxes.Run(ctx, sandbox.Name, SandboxBoxRunParams{Command: String("pwd")})
			return err
		}, 2, tracingSandboxID},
		{"read", func(ctx context.Context) error { _, err := sandbox.ReadFile(ctx, "/file"); return err }, 1, tracingSandboxID},
		{"read by name", func(ctx context.Context) error { _, err := boxes.ReadFile(ctx, sandbox.Name, "/file"); return err }, 2, tracingSandboxID},
		{"write", func(ctx context.Context) error { return sandbox.WriteFile(ctx, "/file", []byte("content")) }, 1, tracingSandboxID},
		{"write by name", func(ctx context.Context) error { return boxes.WriteFile(ctx, sandbox.Name, "/file", []byte("content")) }, 2, tracingSandboxID},
		{"glob", func(ctx context.Context) error { _, err := sandbox.Glob(ctx, SandboxGlobParams{}); return err }, 1, tracingSandboxID},
		{"glob by name", func(ctx context.Context) error {
			_, err := boxes.Glob(ctx, sandbox.Name, SandboxGlobParams{})
			return err
		}, 2, tracingSandboxID},
		{"ls", func(ctx context.Context) error { _, err := sandbox.Ls(ctx, "/"); return err }, 1, tracingSandboxID},
		{"grep", func(ctx context.Context) error { _, err := sandbox.Grep(ctx, SandboxGrepParams{}); return err }, 1, tracingSandboxID},
		{"stat", func(ctx context.Context) error { _, err := sandbox.StatFile(ctx, "/file"); return err }, 1, tracingSandboxID},
		{"range", func(ctx context.Context) error {
			_, err := sandbox.ReadFileRange(ctx, "/file", SandboxReadRangeParams{Start: Int(0)})
			return err
		}, 1, tracingSandboxID},
		{"refresh", func(ctx context.Context) error { return sandbox.Refresh(ctx) }, 1, tracingSandboxID},
		{"update", func(ctx context.Context) error { return sandbox.Update(ctx, SandboxBoxUpdateParams{}) }, 1, tracingSandboxID},
		{"start", func(ctx context.Context) error { return sandbox.Start(ctx, SandboxWaitParams{}) }, 3, tracingSandboxID},
		{"wait", func(ctx context.Context) error {
			_, err := boxes.WaitSandbox(ctx, sandbox.Name, SandboxWaitParams{})
			return err
		}, 2, tracingSandboxID},
		{"stop", func(ctx context.Context) error { return sandbox.Stop(ctx) }, 1, tracingSandboxID},
		{"delete", func(ctx context.Context) error { return sandbox.Delete(ctx) }, 1, tracingSandboxID},
		{"snapshot", func(ctx context.Context) error {
			_, err := sandbox.CaptureSnapshot(ctx, SandboxBoxNewSnapshotParams{})
			return err
		}, 1, tracingSandboxID},
		{"snapshot wait", func(ctx context.Context) error {
			_, err := sandbox.CaptureSnapshotAndWait(ctx, SandboxBoxNewSnapshotParams{}, SnapshotWaitParams{})
			return err
		}, 2, tracingSandboxID},
		{"service", func(ctx context.Context) error {
			_, err := sandbox.Service(ctx, SandboxBoxGenerateServiceURLParams{Port: Int(8080)})
			return err
		}, 1, tracingSandboxID},
		{"unknown service name", func(ctx context.Context) error {
			_, err := boxes.Service(ctx, sandbox.Name, SandboxBoxGenerateServiceURLParams{Port: Int(8080)})
			return err
		}, 1, ""},
		{"bare dataplane", func(ctx context.Context) error {
			_, err := boxes.ReadFileWithDataplaneURL(ctx, srv.URL, "/file")
			return err
		}, 1, ""},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			requests.Store(0)
			ctx, span, recorder := sandboxTraceTestContext(t)
			require.NoError(t, test.call(ctx))
			require.Equal(t, test.requests, requests.Load())
			requireSandboxTrace(t, span, recorder, test.want)
		})
	}
	service, err := sandbox.Service(context.Background(), SandboxBoxGenerateServiceURLParams{Port: Int(8080)})
	require.NoError(t, err)
	ctx, span, recorder := sandboxTraceTestContext(t)
	resp, err := service.Get(ctx, "/", nil)
	require.NoError(t, err)
	require.NoError(t, resp.Body.Close())
	requireSandboxTrace(t, span, recorder, tracingSandboxID)
}

func TestSandboxTracingFailure(t *testing.T) {
	for _, id := range []string{tracingSandboxID, ""} {
		ctx, span, recorder := sandboxTraceTestContext(t)
		sandbox := &Sandbox{ID: id, Name: "display-name"}
		_, err := sandbox.Run(ctx, SandboxBoxRunParams{})
		require.Error(t, err)
		requireSandboxTrace(t, span, recorder, id)
	}
}

func TestSandboxTracingTunnel(t *testing.T) {
	sandbox := &Sandbox{ID: tracingSandboxID, Name: "display-name", DataplaneURL: "https://sandbox.example", boxes: NewSandboxBoxService()}
	for _, call := range []func(context.Context) error{
		func(ctx context.Context) error { _, err := sandbox.Tunnel(ctx, 0, SandboxTunnelParams{}); return err },
		func(ctx context.Context) error { _, err := sandbox.OpenTunnelStream(ctx, 0); return err },
		func(ctx context.Context) error {
			_, err := (&SandboxTunnel{sandboxID: tracingSandboxID, closed: true}).Dial(ctx)
			return err
		},
	} {
		ctx, span, recorder := sandboxTraceTestContext(t)
		require.Error(t, call(ctx))
		requireSandboxTrace(t, span, recorder, tracingSandboxID)
	}
}

func TestSandboxTracingCommandHandle(t *testing.T) {
	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		defer ws.Close()
		var payload map[string]any
		if websocket.JSON.Receive(ws, &payload) != nil {
			return
		}
		_ = websocket.JSON.Send(ws, map[string]any{"type": "started", "command_id": "command", "pid": 1})
		_ = websocket.JSON.Send(ws, map[string]any{"type": "exit", "exit_code": 0})
	}))
	defer srv.Close()
	sandbox := &Sandbox{ID: tracingSandboxID, Name: "display-name", DataplaneURL: srv.URL, boxes: NewSandboxBoxService()}
	ctx, span, recorder := sandboxTraceTestContext(t)
	handle, err := sandbox.StartCommand(ctx, SandboxCommandStartParams{Command: String("pwd")})
	require.NoError(t, err)
	requireSandboxTrace(t, span, recorder, tracingSandboxID)
	ctx, span, recorder = sandboxTraceTestContext(t)
	_, err = handle.Result(ctx)
	require.NoError(t, err)
	requireSandboxTrace(t, span, recorder, tracingSandboxID)
	ctx, span, recorder = sandboxTraceTestContext(t)
	reconnected, err := handle.Reconnect(ctx)
	require.NoError(t, err)
	requireSandboxTrace(t, span, recorder, tracingSandboxID)
	ctx, span, recorder = sandboxTraceTestContext(t)
	_, _, err = reconnected.Next(ctx)
	require.NoError(t, err)
	requireSandboxTrace(t, span, recorder, tracingSandboxID)
}
