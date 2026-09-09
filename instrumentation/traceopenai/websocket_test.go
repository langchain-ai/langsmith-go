package traceopenai

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/gorilla/websocket"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/codes"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
	"go.opentelemetry.io/otel/trace"
)

func websocketPair(t *testing.T, ctx context.Context, opts ...Option) (*WebSocketConn, *websocket.Conn, *tracetest.InMemoryExporter) {
	t.Helper()
	connections := make(chan *websocket.Conn, 1)
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		conn, err := (&websocket.Upgrader{}).Upgrade(w, r, nil)
		if err == nil {
			connections <- conn
		}
	}))
	t.Cleanup(server.Close)
	raw, _, err := websocket.DefaultDialer.Dial("ws"+strings.TrimPrefix(server.URL, "http")+"/v1/responses", nil)
	require.NoError(t, err)
	peer := <-connections
	require.NoError(t, peer.SetReadDeadline(time.Now().Add(10*time.Second)))
	require.NoError(t, peer.SetWriteDeadline(time.Now().Add(10*time.Second)))
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
	t.Cleanup(func() { require.NoError(t, tp.Shutdown(context.Background())) })
	conn := WrapWebSocket(ctx, raw, append(opts, WithTracerProvider(tp))...)
	require.NoError(t, conn.SetReadDeadline(time.Now().Add(10*time.Second)))
	require.NoError(t, conn.SetWriteDeadline(time.Now().Add(10*time.Second)))
	t.Cleanup(func() { _ = conn.Close(); _ = peer.Close() })
	return conn, peer, exporter
}

func sendWebSocket(t *testing.T, conn *WebSocketConn, peer *websocket.Conn, data string) {
	t.Helper()
	require.NoError(t, conn.WriteMessage(websocket.TextMessage, []byte(data)))
	typ, received, err := peer.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, typ)
	require.Equal(t, data, string(received))
}

func receiveWebSocket(t *testing.T, conn *WebSocketConn, peer *websocket.Conn, data string) {
	t.Helper()
	require.NoError(t, peer.WriteMessage(websocket.TextMessage, []byte(data)))
	typ, received, err := conn.ReadMessage()
	require.NoError(t, err)
	require.Equal(t, websocket.TextMessage, typ)
	require.Equal(t, data, string(received))
}

func spanString(span tracetest.SpanStub, key string) string {
	for _, attr := range span.Attributes {
		if string(attr.Key) == key {
			return attr.Value.AsString()
		}
	}
	return ""
}

func TestWebSocketResponses(t *testing.T) {
	parent := trace.NewSpanContext(trace.SpanContextConfig{TraceID: trace.TraceID{1}, SpanID: trace.SpanID{2}, TraceFlags: trace.FlagsSampled})
	ctx := WithRunNameContext(trace.ContextWithSpanContext(context.Background(), parent), "astra")
	conn, peer, exporter := websocketPair(t, ctx, WithRunName("fallback"))
	sendWebSocket(t, conn, peer, `{"type":"response.create","model":"gpt-6-astra","input":[{"type":"configuration_update","reasoning":{"effort":"high"}},{"role":"user","content":"hello"}]}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"resp_1","model":"gpt-6-astra"}}`)
	require.Empty(t, exporter.GetSpans())
	receiveWebSocket(t, conn, peer, `{"type":"response.output_text.delta","response_id":"resp_1","delta":"hi"}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.output_text.delta","response_id":"resp_1","delta":"!"}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.completed","response":{"id":"resp_1","status":"completed","service_tier":"flex","output":[{"type":"message","content":[{"type":"output_text","text":"hi!"}]}],"usage":{"input_tokens":1000,"output_tokens":50,"input_tokens_details":{"cached_tokens":100,"cache_write_tokens":200}}}}`)
	spans := exporter.GetSpans()
	require.Len(t, spans, 1)
	span := spans[0]
	require.Equal(t, parent.SpanID(), span.Parent.SpanID())
	require.Equal(t, "astra", span.Name)
	require.Equal(t, codes.Ok, span.Status.Code)
	require.Equal(t, "gpt-6-astra", spanString(span, "gen_ai.request.model"))
	require.Contains(t, spanString(span, "gen_ai.prompt"), "configuration_update")
	require.Contains(t, spanString(span, "gen_ai.prompt"), "high")
	require.Contains(t, spanString(span, "gen_ai.completion"), "hi!")
	require.Len(t, span.Events, 1)
	require.Equal(t, "new_token", span.Events[0].Name)
	var usage map[string]any
	require.NoError(t, json.Unmarshal([]byte(spanString(span, string(usageMetadataKey))), &usage))
	require.Equal(t, map[string]any{"flex_cache_read": float64(100), "flex_cache_write": float64(200), "flex": float64(700)}, usage["input_token_details"])
	require.Empty(t, conn.active)
	require.Empty(t, conn.lanes)
	// Duplicate terminal events must not emit a second span.
	receiveWebSocket(t, conn, peer, `{"type":"response.completed","response":{"id":"resp_1"}}`)
	require.Len(t, exporter.GetSpans(), 1)
}

func TestWebSocketMultiplexAndQueuedTurns(t *testing.T) {
	conn, peer, exporter := websocketPair(t, context.Background())
	for _, item := range []struct{ lane, prompt string }{{"a", "first"}, {"b", "parallel"}, {"a", "next"}} {
		sendWebSocket(t, conn, peer, fmt.Sprintf(`{"type":"response.create","stream_id":%q,"model":"gpt-6-astra","input":%q}`, item.lane, item.prompt))
	}
	for _, event := range []string{
		`{"type":"response.created","stream_id":"a","response":{"id":"a1"}}`,
		`{"type":"response.created","stream_id":"b","response":{"id":"b1"}}`,
		`{"type":"response.output_text.delta","stream_id":"b","delta":"b"}`,
		`{"type":"response.completed","stream_id":"b","response":{"id":"b1"}}`,
		`{"type":"response.completed","stream_id":"a","response":{"id":"a1"}}`,
		`{"type":"response.created","stream_id":"a","response":{"id":"a2"}}`,
		`{"type":"response.completed","stream_id":"a","response":{"id":"a2"}}`,
	} {
		receiveWebSocket(t, conn, peer, event)
	}
	spans := exporter.GetSpans()
	require.Len(t, spans, 3)
	for i, expected := range []struct{ id, prompt string }{{"b1", "parallel"}, {"a1", "first"}, {"a2", "next"}} {
		require.Equal(t, expected.id, spanString(spans[i], "gen_ai.response.id"))
		require.Contains(t, spanString(spans[i], "gen_ai.prompt"), expected.prompt)
	}
	require.Len(t, spans[0].Events, 1)
	require.Empty(t, spans[1].Events)
	require.Empty(t, spans[2].Events)
	require.Empty(t, conn.lanes)
}

func TestWebSocketSteering(t *testing.T) {
	for _, explicit := range []bool{false, true} {
		t.Run(fmt.Sprintf("tool_result=%v", explicit), func(t *testing.T) {
			conn, peer, exporter := websocketPair(t, context.Background())
			sendWebSocket(t, conn, peer, `{"type":"response.create","model":"gpt-6-astra","input":"draft a plan"}`)
			receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"original"}}`)
			sendWebSocket(t, conn, peer, `{"type":"response.steer","previous_response_id":"original","input":"keep it small"}`)
			receiveWebSocket(t, conn, peer, `{"type":"response.steer.accepted","steer":{"id":"steer_1","previous_response_id":"original"}}`)
			receiveWebSocket(t, conn, peer, `{"type":"response.incomplete","response":{"id":"original","status":"incomplete","incomplete_details":{"reason":"steered"}}}`)
			if explicit {
				receiveWebSocket(t, conn, peer, `{"type":"response.steer.pending","steer":{"id":"steer_1","previous_response_id":"original"}}`)
				sendWebSocket(t, conn, peer, `{"type":"response.create","previous_response_id":"original","model":"gpt-6-astra","input":[{"type":"function_call_output","call_id":"tool_1","output":"tool result"}]}`)
			}
			receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"successor","previous_response_id":"original","model":"gpt-6-astra"}}`)
			receiveWebSocket(t, conn, peer, `{"type":"response.completed","response":{"id":"successor","status":"completed"}}`)
			spans := exporter.GetSpans()
			require.Len(t, spans, 2)
			require.Equal(t, codes.Ok, spans[0].Status.Code)
			require.Equal(t, "steered", spanString(spans[0], "langsmith.metadata.incomplete_reason"))
			require.Contains(t, spanString(spans[1], "gen_ai.prompt"), "keep it small")
			require.Equal(t, "original", spanString(spans[1], "langsmith.metadata.previous_response_id"))
			if explicit {
				require.Contains(t, spanString(spans[1], "gen_ai.prompt"), "tool result")
			}
			require.Empty(t, conn.steers)
		})
	}
}

func TestWebSocketSteeringFailure(t *testing.T) {
	conn, peer, exporter := websocketPair(t, context.Background())
	sendWebSocket(t, conn, peer, `{"type":"response.create","model":"gpt-6-astra","input":"hello"}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"original"}}`)
	sendWebSocket(t, conn, peer, `{"type":"response.steer","previous_response_id":"original","input":"rejected"}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.steer.failed","steer":{"id":"rejected","previous_response_id":"original"},"error":{"code":"invalid_input"}}`)
	require.Empty(t, conn.steers)
	require.Empty(t, exporter.GetSpans())
	receiveWebSocket(t, conn, peer, `{"type":"response.completed","response":{"id":"original"}}`)
	require.Equal(t, codes.Ok, exporter.GetSpans()[0].Status.Code)
}

func TestWebSocketTerminalStates(t *testing.T) {
	for _, tc := range []struct {
		kind string
		code codes.Code
	}{
		{"response.failed", codes.Error}, {"response.cancelled", codes.Error},
		{"response.incomplete", codes.Ok}, {"error", codes.Error},
	} {
		t.Run(tc.kind, func(t *testing.T) {
			conn, peer, exporter := websocketPair(t, context.Background())
			sendWebSocket(t, conn, peer, `{"type":"response.create","model":"gpt-6-astra","input":"hello"}`)
			receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"r"}}`)
			receiveWebSocket(t, conn, peer, fmt.Sprintf(`{"type":%q,"response":{"id":"r","usage":{"input_tokens":10,"output_tokens":5}}}`, tc.kind))
			spans := exporter.GetSpans()
			require.Len(t, spans, 1)
			require.Equal(t, tc.code, spans[0].Status.Code)
			require.Empty(t, conn.active)
		})
	}
}

func TestWebSocketCleanup(t *testing.T) {
	for _, mode := range []string{"close", "disconnect", "cancel", "write_failure"} {
		t.Run(mode, func(t *testing.T) {
			ctx, cancel := context.WithCancel(context.Background())
			defer cancel()
			conn, peer, exporter := websocketPair(t, ctx)
			sendWebSocket(t, conn, peer, `{"type":"response.create","model":"gpt-6-astra","input":"hello"}`)
			switch mode {
			case "close":
				require.NoError(t, conn.Close())
			case "disconnect":
				require.NoError(t, peer.Close())
				_, _, err := conn.ReadMessage()
				require.Error(t, err)
			case "cancel":
				cancel()
				_, _, err := conn.ReadMessage()
				require.Error(t, err)
			case "write_failure":
				require.NoError(t, conn.conn.Close())
				require.Error(t, conn.WriteJSON(map[string]any{"type": "response.create", "input": "second"}))
			}
			spans := exporter.GetSpans()
			if mode == "write_failure" {
				require.Len(t, spans, 2)
			} else {
				require.Len(t, spans, 1)
			}
			for _, span := range spans {
				require.Equal(t, codes.Error, span.Status.Code)
			}
			_ = conn.Close()
			require.Len(t, exporter.GetSpans(), len(spans))
		})
	}
}

func TestWebSocketConcurrentReadWrite(t *testing.T) {
	conn, peer, exporter := websocketPair(t, context.Background())
	serverErrors := make(chan error, 1)
	go func() {
		for range 20 {
			if _, _, err := peer.ReadMessage(); err != nil {
				serverErrors <- err
				return
			}
			if err := peer.WriteJSON(map[string]any{"type": "response.created", "response": map[string]any{"id": "r"}}); err != nil {
				serverErrors <- err
				return
			}
			if err := peer.WriteJSON(map[string]any{"type": "response.completed", "response": map[string]any{"id": "r"}}); err != nil {
				serverErrors <- err
				return
			}
		}
		serverErrors <- nil
	}()
	readErrors := make(chan error, 1)
	go func() {
		for range 40 {
			var event map[string]any
			if err := conn.ReadJSON(&event); err != nil {
				readErrors <- err
				return
			}
		}
		readErrors <- nil
	}()
	for range 20 {
		require.NoError(t, conn.WriteJSON(map[string]any{"type": "response.create", "model": "gpt-6-astra", "input": "hi"}))
	}
	require.NoError(t, <-readErrors)
	require.NoError(t, <-serverErrors)
	require.Len(t, exporter.GetSpans(), 20)
}

func TestWebSocketIgnoresNonResponseMessages(t *testing.T) {
	conn, peer, exporter := websocketPair(t, context.Background())
	for _, data := range []string{"not json", `{"type":"unknown"}`} {
		sendWebSocket(t, conn, peer, data)
		receiveWebSocket(t, conn, peer, data)
	}
	require.Error(t, conn.WriteJSON(make(chan int)))
	require.Empty(t, exporter.GetSpans())
}

func TestWebSocketErrorDoesNotEndOtherLanes(t *testing.T) {
	conn, peer, exporter := websocketPair(t, context.Background())
	for _, lane := range []string{"a", "b"} {
		sendWebSocket(t, conn, peer, fmt.Sprintf(`{"type":"response.create","stream_id":%q,"input":"hello"}`, lane))
	}
	// A pre-response error ends only its own request; raw error text is not exported.
	receiveWebSocket(t, conn, peer, `{"type":"error","stream_id":"a","error":{"message":"private request details"}}`)
	require.Len(t, exporter.GetSpans(), 1)
	require.Equal(t, "a", spanString(exporter.GetSpans()[0], "langsmith.metadata.stream_id"))
	require.NotContains(t, fmt.Sprint(exporter.GetSpans()), "private request details")
	// An error in the default lane must not end a named lane's response.
	receiveWebSocket(t, conn, peer, `{"type":"error","error":{"message":"default lane failed"}}`)
	require.Len(t, exporter.GetSpans(), 1)
	receiveWebSocket(t, conn, peer, `{"type":"response.created","stream_id":"b","response":{"id":"b"}}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.completed","stream_id":"b","response":{"id":"b"}}`)
	require.Len(t, exporter.GetSpans(), 2)
	require.Equal(t, codes.Ok, exporter.GetSpans()[1].Status.Code)
}

func TestWebSocketMultipleSteers(t *testing.T) {
	conn, peer, exporter := websocketPair(t, context.Background())
	sendWebSocket(t, conn, peer, `{"type":"response.create","model":"gpt-6-astra","input":"hello"}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"original"}}`)
	for _, input := range []string{"first", "rejected", "last"} {
		sendWebSocket(t, conn, peer, fmt.Sprintf(`{"type":"response.steer","previous_response_id":"original","input":%q}`, input))
		receiveWebSocket(t, conn, peer, fmt.Sprintf(`{"type":"response.steer.accepted","steer":{"id":%q,"previous_response_id":"original"}}`, input))
	}
	receiveWebSocket(t, conn, peer, `{"type":"response.completed","response":{"id":"original"}}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.steer.failed","steer":{"id":"rejected","previous_response_id":"original"}}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.created","response":{"id":"next","previous_response_id":"original","model":"gpt-6-astra"}}`)
	receiveWebSocket(t, conn, peer, `{"type":"response.completed","response":{"id":"next"}}`)
	spans := exporter.GetSpans()
	require.Len(t, spans, 2)
	var prompt struct {
		Messages []map[string]any `json:"messages"`
	}
	require.NoError(t, json.Unmarshal([]byte(spanString(spans[1], "gen_ai.prompt")), &prompt))
	require.Equal(t, []map[string]any{{"role": "user", "content": "first"}, {"role": "user", "content": "last"}}, prompt.Messages)
	require.Empty(t, conn.steers)
}
