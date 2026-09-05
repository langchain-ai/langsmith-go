package traceopenai

import (
	"context"
	"encoding/json"
	"errors"
	"net"
	"slices"
	"sync"
	"time"

	"github.com/gorilla/websocket"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// WebSocketConn traces Responses API messages on an existing WebSocket connection.
// Use its message or JSON methods for all application reads and writes. As with
// gorilla/websocket, callers may run one reader and one writer concurrently.
type WebSocketConn struct {
	conn    *websocket.Conn
	ctx     context.Context
	tracer  trace.Tracer
	name    string
	mu      sync.Mutex
	writeMu sync.Mutex
	closed  bool
	stop    func() bool
	lanes   map[string][]*webSocketResponse
	active  map[string]*webSocketResponse
	steers  map[string][]*webSocketSteer
}

type webSocketResponse struct {
	span       trace.Span
	id         string
	lane       string
	prompt     string
	firstToken bool
}

type webSocketSteer struct {
	id     string
	prompt string
}

// WrapWebSocket instruments a connection already dialed to the Responses API.
// The caller owns dialing, authentication, TLS, and reconnection. Each response
// gets a child span of ctx, including automatic steering continuations. Reading
// a terminal response ends its span; Close or cancellation of ctx closes the
// connection and ends outstanding spans. No background message reader is started.
func WrapWebSocket(ctx context.Context, conn *websocket.Conn, opts ...Option) *WebSocketConn {
	options := &clientOptions{}
	for _, opt := range opts {
		opt(options)
	}
	tp := options.tracerProvider
	if tp == nil {
		tp = otel.GetTracerProvider()
	}
	name := options.runName
	if v, ok := ctx.Value(ctxKeyRunName).(string); ok {
		name = v
	}
	if name == "" {
		name = "openai.responses"
	}
	c := &WebSocketConn{
		conn: conn, ctx: ctx, tracer: tp.Tracer("github.com/sashabaranov/go-openai"), name: name,
		lanes: make(map[string][]*webSocketResponse), active: make(map[string]*webSocketResponse),
		steers: make(map[string][]*webSocketSteer),
	}
	c.mu.Lock()
	c.stop = context.AfterFunc(ctx, func() {
		_ = conn.Close()
		c.finishAll(ctx.Err())
	})
	c.mu.Unlock()
	return c
}

// WriteMessage sends an unchanged WebSocket message and traces response requests.
func (c *WebSocketConn) WriteMessage(messageType int, data []byte) error {
	c.writeMu.Lock()
	defer c.writeMu.Unlock()
	c.mu.Lock()
	if !c.closed && messageType == websocket.TextMessage {
		var event map[string]any
		if json.Unmarshal(data, &event) == nil {
			switch stringField(event, "type") {
			case "response.create":
				c.startResponse(event, stringField(event, "stream_id"))
			case "response.steer":
				previous := stringField(event, "previous_response_id")
				prompt := parseRequestBody(data).inputMessages
				c.steers[previous] = append(c.steers[previous], &webSocketSteer{prompt: prompt})
				if r := c.active[previous]; r != nil {
					r.span.AddEvent("response.steer", trace.WithAttributes(attribute.String("gen_ai.prompt", prompt)))
				}
			}
		}
	}
	c.mu.Unlock()
	if err := c.conn.WriteMessage(messageType, data); err != nil {
		_ = c.conn.Close()
		c.finishAll(err)
		return err
	}
	return nil
}

// ReadMessage receives an unchanged message and records response lifecycle events.
func (c *WebSocketConn) ReadMessage() (int, []byte, error) {
	messageType, data, err := c.conn.ReadMessage()
	if err != nil {
		_ = c.conn.Close()
		c.finishAll(err)
		return messageType, data, err
	}
	if messageType == websocket.TextMessage {
		var event map[string]any
		if json.Unmarshal(data, &event) == nil {
			c.mu.Lock()
			if !c.closed {
				c.receive(event)
			}
			c.mu.Unlock()
		}
	}
	return messageType, data, nil
}

// WriteJSON marshals v and sends a traced text message.
func (c *WebSocketConn) WriteJSON(v any) error {
	data, err := json.Marshal(v)
	if err != nil {
		return err
	}
	return c.WriteMessage(websocket.TextMessage, data)
}

// ReadJSON receives a traced message and unmarshals it into v.
func (c *WebSocketConn) ReadJSON(v any) error {
	_, data, err := c.ReadMessage()
	if err != nil {
		return err
	}
	return json.Unmarshal(data, v)
}

// Close closes the connection and ends any unfinished response spans.
func (c *WebSocketConn) Close() error {
	err := c.conn.Close()
	c.finishAll(net.ErrClosed)
	return err
}

// SetReadDeadline sets the underlying connection's read deadline.
func (c *WebSocketConn) SetReadDeadline(t time.Time) error { return c.conn.SetReadDeadline(t) }

// SetWriteDeadline sets the underlying connection's write deadline.
func (c *WebSocketConn) SetWriteDeadline(t time.Time) error { return c.conn.SetWriteDeadline(t) }

// SetReadLimit sets the maximum incoming WebSocket message size in bytes.
func (c *WebSocketConn) SetReadLimit(limit int64) { c.conn.SetReadLimit(limit) }

func (c *WebSocketConn) startResponse(request map[string]any, lane string) *webSocketResponse {
	data, _ := json.Marshal(request)
	fields := parseRequestBody(data)
	attrs := []attribute.KeyValue{
		attribute.String("gen_ai.system", "openai"),
		attribute.String("gen_ai.operation.name", "responses"),
		attribute.String("network.protocol.name", "websocket"),
	}
	if fields.model != "" {
		attrs = append(attrs, attribute.String("gen_ai.request.model", fields.model))
	}
	if lane != "" {
		attrs = append(attrs, attribute.String("langsmith.metadata.stream_id", lane))
	}
	if previous := stringField(request, "previous_response_id"); previous != "" {
		attrs = append(attrs, attribute.String("langsmith.metadata.previous_response_id", previous))
	}
	bag := baggage.FromContext(c.ctx)
	for _, key := range []string{"session_id", "thread_id", "conversation_id"} {
		if member := bag.Member(key); member.Key() == key {
			attrs = append(attrs, attribute.String(key, member.Value()), attribute.String("langsmith.metadata."+key, member.Value()))
		}
	}
	_, span := c.tracer.Start(c.ctx, c.name, trace.WithAttributes(attrs...))
	r := &webSocketResponse{span: span, lane: lane, prompt: fields.inputMessages}
	if r.prompt != "" {
		span.SetAttributes(attribute.String("gen_ai.prompt", r.prompt))
	}
	c.lanes[lane] = append(c.lanes[lane], r)
	return r
}

func (c *WebSocketConn) receive(event map[string]any) {
	kind := stringField(event, "type")
	if kind == "response.steer.accepted" || kind == "response.steer.failed" || kind == "response.steer.pending" {
		c.receiveSteer(event, kind)
		return
	}
	response, _ := event["response"].(map[string]any)
	id := stringField(response, "id")
	if id == "" {
		id = stringField(event, "response_id")
	}
	lane := stringField(event, "stream_id")
	r := c.active[id]
	if r == nil && kind == "response.created" {
		if id == "" {
			return
		}
		for _, pending := range c.lanes[lane] {
			if pending.id == "" {
				r = pending
				break
			}
		}
		if r == nil {
			// Steering can create a response without a client response.create.
			r = c.startResponse(response, lane)
		}
		r.id = id
		c.active[id] = r
		r.span.SetAttributes(attribute.String("gen_ai.response.id", id))
		previous := stringField(response, "previous_response_id")
		var prompts []string
		var unacknowledged []*webSocketSteer
		for _, steer := range c.steers[previous] {
			if steer.id != "" {
				prompts = append(prompts, steer.prompt)
			} else {
				unacknowledged = append(unacknowledged, steer)
			}
		}
		if len(unacknowledged) == 0 {
			delete(c.steers, previous)
		} else {
			c.steers[previous] = unacknowledged
		}
		prompts = append(prompts, r.prompt)
		if prompt := mergeWebSocketPrompts(prompts); prompt != "" {
			r.span.SetAttributes(attribute.String("gen_ai.prompt", prompt))
		}
	}
	if r == nil && id == "" && len(c.lanes[lane]) > 0 {
		r = c.lanes[lane][0]
	}
	if kind == "error" {
		if r != nil {
			c.finishResponse(r, nil, errors.New("OpenAI WebSocket request failed"))
		}
		return
	}
	if r == nil {
		return
	}
	if model := stringField(response, "model"); model != "" {
		r.span.SetAttributes(attribute.String("gen_ai.response.model", model))
	}
	if !r.firstToken && isFirstContentResponses(event) {
		r.firstToken = true
		r.span.AddEvent("new_token")
	}
	switch kind {
	case "response.completed", "response.incomplete", "response.failed", "response.cancelled":
		if response == nil || id == "" {
			return
		}
		r.span.SetAttributes(attribute.String("langsmith.metadata.response_status", stringField(response, "status")))
		var err error
		if kind == "response.failed" || kind == "response.cancelled" {
			err = errors.New("OpenAI " + kind)
		}
		if details, ok := response["incomplete_details"].(map[string]any); ok {
			r.span.SetAttributes(attribute.String("langsmith.metadata.incomplete_reason", stringField(details, "reason")))
		}
		c.finishResponse(r, response, err)
	}
}

func (c *WebSocketConn) receiveSteer(event map[string]any, kind string) {
	steer, _ := event["steer"].(map[string]any)
	previous, id := stringField(steer, "previous_response_id"), stringField(steer, "id")
	if r := c.active[previous]; r != nil {
		r.span.AddEvent(kind, trace.WithAttributes(attribute.String("steer.id", id)))
	}
	for i, pending := range c.steers[previous] {
		if pending.id == id || (pending.id == "" && kind != "response.steer.pending") {
			if kind == "response.steer.failed" {
				queue := c.steers[previous]
				c.steers[previous] = slices.Delete(queue, i, i+1)
				if len(c.steers[previous]) == 0 {
					delete(c.steers, previous)
				}
			} else {
				pending.id = id
			}
			break
		}
	}
}

func (c *WebSocketConn) finishResponse(r *webSocketResponse, response map[string]any, err error) {
	if completion := extractResponsesOutput(response); completion != "" {
		r.span.SetAttributes(attribute.String("gen_ai.completion", completion))
	}
	setOpenAIUsageAttributes(r.span, trace.SpanFromContext(c.ctx), extractResponsesUsage(response))
	if err != nil {
		// Transport error strings can include caller-controlled URLs or credentials.
		r.span.SetStatus(codes.Error, "WebSocket response interrupted")
		r.span.RecordError(errors.New("WebSocket response interrupted"))
	} else {
		r.span.SetStatus(codes.Ok, "")
	}
	r.span.End()
	delete(c.active, r.id)
	queue := c.lanes[r.lane]
	for i, candidate := range queue {
		if candidate == r {
			queue = slices.Delete(queue, i, i+1)
			break
		}
	}
	if len(queue) == 0 {
		delete(c.lanes, r.lane)
	} else {
		c.lanes[r.lane] = queue
	}
}

func (c *WebSocketConn) finishAll(err error) {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.finishAllLocked(err)
}

func (c *WebSocketConn) finishAllLocked(err error) {
	if c.closed {
		return
	}
	c.closed = true
	if c.stop != nil {
		c.stop()
	}
	for _, queue := range c.lanes {
		for _, r := range append([]*webSocketResponse(nil), queue...) {
			c.finishResponse(r, nil, err)
		}
	}
	clear(c.steers)
}

func stringField(m map[string]any, key string) string {
	s, _ := m[key].(string)
	return s
}

func mergeWebSocketPrompts(prompts []string) string {
	var messages []any
	for _, prompt := range prompts {
		var body struct {
			Messages []any `json:"messages"`
		}
		if json.Unmarshal([]byte(prompt), &body) == nil {
			messages = append(messages, body.Messages...)
		}
	}
	if len(messages) == 0 {
		return ""
	}
	return marshalMessages(messages)
}
