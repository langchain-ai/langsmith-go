package messagetranslators

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

// SSEEvent is one complete server-sent event. Data contains only the event's
// data field; callers remain responsible for HTTP/SSE framing.
type SSEEvent struct {
	Event string
	Data  []byte
}

type decodedStreamEvent struct {
	object map[string]any
	typ    string
	err    error
	done   bool
}

func decodeStreamEvent(e SSEEvent) decodedStreamEvent {
	if string(e.Data) == "[DONE]" {
		return decodedStreamEvent{done: true}
	}
	o, err := decodeObject(e.Data)
	typ := e.Event
	if err == nil && typ == "" {
		typ, _ = str(o["type"])
	}
	return decodedStreamEvent{object: o, typ: typ, err: err}
}

func streamObject(e SSEEvent, decoded decodedStreamEvent) (map[string]any, string, error) {
	if decoded.done {
		return nil, "", at("event", fmt.Errorf("%w: untyped done sentinel", ErrUnsupported))
	}
	if decoded.err != nil {
		return nil, "", at("event "+e.Event, decoded.err)
	}
	o := decoded.object
	dataType, hasDataType := str(o["type"])
	if e.Event != "" && hasDataType && dataType != e.Event {
		return nil, "", at("event.type", fmt.Errorf("%w: Event header %q does not match data type %q", ErrInvalidWireData, e.Event, dataType))
	}
	if decoded.typ == "" {
		return nil, "", at("event", fmt.Errorf("%w: missing event type", ErrInvalidWireData))
	}
	return o, decoded.typ, nil
}
func event(typ string, v map[string]any) (SSEEvent, error) {
	v["type"] = typ
	b, err := json.Marshal(v)
	return SSEEvent{Event: typ, Data: b}, err
}
func one(typ string, v map[string]any) ([]SSEEvent, error) {
	e, err := event(typ, v)
	if err != nil {
		return nil, err
	}
	return []SSEEvent{e}, nil
}

// ResponsesToAnthropicStream incrementally translates complete Responses API
// SSE events into complete Anthropic Messages SSE events.
type responseItemKind uint8

const (
	responseMessageItem responseItemKind = iota + 1
	responseToolItem
)

type responseItemState struct {
	id   string
	kind responseItemKind
	done bool
}

type responseBlockKind uint8

const (
	responseTextBlock responseBlockKind = iota + 1
	responseToolBlock
)

type responseBlockKey struct {
	outputIndex  int
	contentIndex int
	kind         responseBlockKind
}

type responseBlockState struct {
	index        int
	open         bool
	arguments    string
	argumentDone bool
	textDone     bool
}

type ResponsesToAnthropicStream struct {
	options           ConversionOptions
	model             string
	started, terminal bool
	id, sourceID      string
	next              int
	blocks            map[responseBlockKey]responseBlockState
	items             map[int]responseItemState
	inputUsageKnown   bool
	inputTokens       int64
	cacheReadTokens   int64
}

// NewResponsesToAnthropicStream constructs a Responses-to-Anthropic stream converter.
func NewResponsesToAnthropicStream(model string) *ResponsesToAnthropicStream {
	return NewResponsesToAnthropicStreamWithOptions(model, ConversionOptions{})
}

// NewResponsesToAnthropicStreamWithOptions constructs a converter with per-stream warning options.
func NewResponsesToAnthropicStreamWithOptions(model string, options ConversionOptions) *ResponsesToAnthropicStream {
	return &ResponsesToAnthropicStream{
		model: model, options: options,
		blocks: map[responseBlockKey]responseBlockState{},
		items:  map[int]responseItemState{},
	}
}

func textBlockKey(outputIndex, contentIndex int) responseBlockKey {
	return responseBlockKey{outputIndex: outputIndex, contentIndex: contentIndex, kind: responseTextBlock}
}

func toolBlockKey(outputIndex int) responseBlockKey {
	return responseBlockKey{outputIndex: outputIndex, kind: responseToolBlock}
}

func (s *ResponsesToAnthropicStream) allocate(key responseBlockKey) responseBlockState {
	if block, ok := s.blocks[key]; ok {
		return block
	}
	block := responseBlockState{index: s.next}
	s.next++
	s.blocks[key] = block
	return block
}

func (s *ResponsesToAnthropicStream) close(key responseBlockKey) ([]SSEEvent, error) {
	block, ok := s.blocks[key]
	if !ok || !block.open {
		return nil, fmt.Errorf("%w: duplicate stop", ErrInvalidSequence)
	}
	block.open = false
	s.blocks[key] = block
	return one("content_block_stop", map[string]any{"index": block.index})
}

// Convert translates one complete event and may return zero or more events.
func (s *ResponsesToAnthropicStream) Convert(e SSEEvent) ([]SSEEvent, error) {
	decoded := decodeStreamEvent(e)
	if decoded.err == nil && !decoded.done {
		inspectResponsesEvent(decoded.object, decoded.typ, s.options)
	}
	before := *s
	before.blocks = cloneMap(s.blocks)
	before.items = cloneMap(s.items)
	out, err := s.convert(e, decoded)
	if err != nil {
		*s = before
		return nil, err
	}
	return out, nil
}

func (s *ResponsesToAnthropicStream) convert(e SSEEvent, decoded decodedStreamEvent) ([]SSEEvent, error) {
	if s.terminal {
		return nil, fmt.Errorf("%w: event after terminal event", ErrInvalidSequence)
	}
	o, typ, err := streamObject(e, decoded)
	if err != nil {
		return nil, err
	}
	for _, key := range []string{"output_index", "content_index"} {
		if v, ok := o[key]; ok {
			if _, err := integer(v, "event."+key, false); err != nil {
				return nil, err
			}
		}
	}
	switch typ {
	case "response.output_item.added", "response.output_item.done", "response.function_call_arguments.delta", "response.function_call_arguments.done":
		if _, err := integer(o["output_index"], "event.output_index", false); err != nil {
			return nil, err
		}
	case "response.content_part.added", "response.content_part.done", "response.output_text.delta", "response.output_text.done":
		if _, err := integer(o["output_index"], "event.output_index", false); err != nil {
			return nil, err
		}
		if _, err := integer(o["content_index"], "event.content_index", false); err != nil {
			return nil, err
		}
	}
	if itemID, present := o["item_id"]; present {
		id, ok := str(itemID)
		if !ok || id == "" {
			return nil, at("event.item_id", ErrInvalidWireData)
		}
		if v, ok := o["output_index"]; ok {
			out, _ := integer(v, "event.output_index", false)
			if known := s.items[int(out)].id; known != "" && known != id {
				return nil, at("event.item_id", fmt.Errorf("%w: item ID mismatch", ErrInvalidSequence))
			}
		}
	}
	switch typ {
	case "response.in_progress":
		if !s.started {
			return nil, fmt.Errorf("%w: in_progress before response.created", ErrInvalidSequence)
		}
		return nil, nil
	case "response.created":
		if s.started {
			return nil, fmt.Errorf("%w: duplicate response start", ErrInvalidSequence)
		}
		r, _ := obj(o["response"])
		if r == nil {
			r = o
		}
		s.sourceID, _ = str(r["id"])
		if s.sourceID == "" {
			return nil, at("event.response.id", ErrInvalidWireData)
		}
		s.id = destinationID("msg", s.sourceID, 0)
		m := s.model
		if m == "" {
			m, _ = str(r["model"])
		}
		if m == "" {
			return nil, at("event.response.model", ErrInvalidWireData)
		}
		usage := map[string]any{"input_tokens": int64(0), "output_tokens": int64(0)}
		if u, ok := obj(r["usage"]); ok {
			au, err := openAIUsageToAnthropic(u, "event.response.usage")
			if err != nil {
				return nil, err
			}
			usage = au
			s.inputUsageKnown = true
			s.inputTokens = au["input_tokens"].(int64)
			s.cacheReadTokens = au["cache_read_input_tokens"].(int64)
		} else if r["usage"] != nil {
			return nil, at("event.response.usage", ErrInvalidWireData)
		}
		s.started = true
		return one("message_start", map[string]any{"message": map[string]any{"id": s.id, "type": "message", "role": "assistant", "content": []any{}, "model": m, "stop_reason": nil, "stop_sequence": nil, "usage": usage}})
	case "response.output_item.added":
		if !s.started {
			return nil, fmt.Errorf("%w: output before response start", ErrInvalidSequence)
		}
		item, ok := obj(o["item"])
		if !ok {
			return nil, at("event.item", ErrInvalidWireData)
		}
		out64, _ := integer(o["output_index"], "event.output_index", false)
		outIndex := int(out64)
		if int64(outIndex) != out64 {
			return nil, at("event.output_index", ErrInvalidWireData)
		}
		if _, exists := s.items[outIndex]; exists {
			return nil, fmt.Errorf("%w: duplicate output item", ErrInvalidSequence)
		}
		if err := rejectPresent(item, "event.item", "phase"); err != nil {
			return nil, err
		}
		if item["type"] == "message" {
			itemID, _ := str(item["id"])
			if itemID == "" {
				return nil, at("event.item.id", ErrInvalidWireData)
			}
			if v, exists := item["role"]; exists && v != "assistant" {
				return nil, at("event.item.role", ErrInvalidWireData)
			}
			s.items[outIndex] = responseItemState{id: itemID, kind: responseMessageItem}
			return nil, nil
		}
		if item["type"] != "function_call" {
			return nil, at("event.item.type", ErrUnsupported)
		}
		key := toolBlockKey(outIndex)
		block := s.allocate(key)
		itemID, _ := str(item["id"])
		if itemID == "" {
			return nil, at("event.item.id", ErrInvalidWireData)
		}
		s.items[outIndex] = responseItemState{id: itemID, kind: responseToolItem}
		id, _ := str(item["call_id"])
		name, _ := str(item["name"])
		if id == "" || name == "" {
			return nil, at("event.item", ErrInvalidWireData)
		}
		block.open = true
		s.blocks[key] = block
		return one("content_block_start", map[string]any{"index": block.index, "content_block": map[string]any{"type": "tool_use", "id": id, "name": name, "input": map[string]any{}}})
	case "response.content_part.added":
		if !s.started {
			return nil, fmt.Errorf("%w: content before start", ErrInvalidSequence)
		}
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		item, exists := s.items[out]
		if int64(out) != out64 || !exists {
			return nil, fmt.Errorf("%w: content for unknown output item", ErrInvalidSequence)
		}
		if item.done {
			return nil, fmt.Errorf("%w: content after output item done", ErrInvalidSequence)
		}
		if item.kind != responseMessageItem {
			return nil, fmt.Errorf("%w: content part on tool output item", ErrInvalidSequence)
		}
		part, ok := obj(o["part"])
		if !ok {
			return nil, at("event.part", ErrInvalidWireData)
		}
		if part["type"] != "output_text" {
			return nil, at("event.part.type", ErrUnsupported)
		}
		if err := rejectNonEmptyArray(part, "annotations", "event.part"); err != nil {
			return nil, err
		}
		if err := rejectPresent(part, "event.part", "phase"); err != nil {
			return nil, err
		}
		initial, ok := str(part["text"])
		if !ok {
			return nil, at("event.part.text", ErrInvalidWireData)
		}
		content64, _ := integer(o["content_index"], "event.content_index", false)
		content := int(content64)
		if int64(content) != content64 {
			return nil, at("event.content_index", ErrInvalidWireData)
		}
		key := textBlockKey(out, content)
		block := s.allocate(key)
		if block.open {
			return nil, fmt.Errorf("%w: duplicate content part", ErrInvalidSequence)
		}
		block.open = true
		s.blocks[key] = block
		return one("content_block_start", map[string]any{"index": block.index, "content_block": map[string]any{"type": "text", "text": initial}})
	case "response.output_text.delta":
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		item, ok := s.items[out]
		if int64(out) != out64 || !ok || item.kind != responseMessageItem || item.done {
			return nil, fmt.Errorf("%w: text delta on non-message or done item", ErrInvalidSequence)
		}
		content64, _ := integer(o["content_index"], "event.content_index", false)
		content := int(content64)
		block, ok := s.blocks[textBlockKey(out, content)]
		if int64(content) != content64 || !ok || !block.open || block.textDone {
			return nil, fmt.Errorf("%w: text delta before content start", ErrInvalidSequence)
		}
		d, ok := str(o["delta"])
		if !ok {
			return nil, at("event.delta", ErrInvalidWireData)
		}
		return one("content_block_delta", map[string]any{"index": block.index, "delta": map[string]any{"type": "text_delta", "text": d}})
	case "response.function_call_arguments.delta":
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		item, ok := s.items[out]
		if int64(out) != out64 || !ok || item.kind != responseToolItem || item.done {
			return nil, fmt.Errorf("%w: arguments delta on non-tool or done item", ErrInvalidSequence)
		}
		key := toolBlockKey(out)
		block, ok := s.blocks[key]
		if !ok || !block.open || block.argumentDone {
			return nil, fmt.Errorf("%w: arguments delta before tool start", ErrInvalidSequence)
		}
		d, ok := str(o["delta"])
		if !ok {
			return nil, at("event.delta", ErrInvalidWireData)
		}
		block.arguments += d
		s.blocks[key] = block
		return one("content_block_delta", map[string]any{"index": block.index, "delta": map[string]any{"type": "input_json_delta", "partial_json": d}})
	case "response.content_part.done":
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		item, ok := s.items[out]
		if int64(out) != out64 || !ok || item.kind != responseMessageItem || item.done {
			return nil, fmt.Errorf("%w: content done on non-message or done item", ErrInvalidSequence)
		}
		if part, exists := o["part"]; exists {
			q, ok := obj(part)
			if !ok || q["type"] != "output_text" {
				return nil, at("event.part", ErrUnsupported)
			}
			if err := rejectNonEmptyArray(q, "annotations", "event.part"); err != nil {
				return nil, err
			}
		}
		content64, _ := integer(o["content_index"], "event.content_index", false)
		content := int(content64)
		key := textBlockKey(out, content)
		if int64(content) != content64 {
			return nil, fmt.Errorf("%w: unknown content part", ErrInvalidSequence)
		}
		if _, ok := s.blocks[key]; !ok {
			return nil, fmt.Errorf("%w: unknown content part", ErrInvalidSequence)
		}
		return s.close(key)
	case "response.output_item.done":
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		if int64(out) != out64 {
			return nil, at("event.output_index", ErrInvalidWireData)
		}
		itemState, exists := s.items[out]
		if !exists {
			return nil, fmt.Errorf("%w: done for unknown output item", ErrInvalidSequence)
		}
		if itemState.done {
			return nil, fmt.Errorf("%w: duplicate output item done", ErrInvalidSequence)
		}
		if item, exists := o["item"]; exists {
			q, ok := obj(item)
			if !ok {
				return nil, at("event.item", ErrInvalidWireData)
			}
			wantType := map[responseItemKind]string{responseMessageItem: "message", responseToolItem: "function_call"}[itemState.kind]
			if q["type"] != nil && q["type"] != wantType {
				return nil, at("event.item.type", fmt.Errorf("%w: output item kind mismatch", ErrInvalidSequence))
			}
			if id, ok := str(q["id"]); ok && id != itemState.id {
				return nil, at("event.item.id", fmt.Errorf("%w: output item ID mismatch", ErrInvalidSequence))
			}
		}
		if itemState.kind == responseMessageItem {
			for key, block := range s.blocks {
				if key.kind == responseTextBlock && key.outputIndex == out && block.open {
					return nil, fmt.Errorf("%w: output item done with open content part", ErrInvalidSequence)
				}
			}
		}
		itemState.done = true
		s.items[out] = itemState
		key := toolBlockKey(out)
		block, isTool := s.blocks[key]
		if !isTool {
			return nil, nil
		}
		final := block.arguments
		if item, ok := obj(o["item"]); ok {
			if x, ok := str(item["arguments"]); ok {
				final = x
			}
		}
		if _, err := parseArguments(final, "event.item.arguments"); err != nil {
			return nil, err
		}
		return s.close(key)
	case "response.function_call_arguments.done":
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		item, ok := s.items[out]
		if int64(out) != out64 || !ok || item.kind != responseToolItem || item.done {
			return nil, fmt.Errorf("%w: arguments done on non-tool or done item", ErrInvalidSequence)
		}
		key := toolBlockKey(out)
		block, ok := s.blocks[key]
		if !ok || !block.open {
			return nil, fmt.Errorf("%w: arguments done for unknown or stopped tool", ErrInvalidSequence)
		}
		if block.argumentDone {
			return nil, fmt.Errorf("%w: duplicate arguments done", ErrInvalidSequence)
		}
		block.argumentDone = true
		final := block.arguments
		if x, ok := str(o["arguments"]); ok {
			final = x
		}
		if _, err := parseArguments(final, "event.arguments"); err != nil {
			return nil, err
		}
		block.arguments = final
		s.blocks[key] = block
		return nil, nil
	case "response.output_text.done":
		out64, _ := integer(o["output_index"], "event.output_index", false)
		out := int(out64)
		item, ok := s.items[out]
		if int64(out) != out64 || !ok || item.kind != responseMessageItem || item.done {
			return nil, fmt.Errorf("%w: text done on non-message or done item", ErrInvalidSequence)
		}
		if err := rejectNonEmptyArray(o, "annotations", "event"); err != nil {
			return nil, err
		}
		content64, _ := integer(o["content_index"], "event.content_index", false)
		content := int(content64)
		key := textBlockKey(out, content)
		block, ok := s.blocks[key]
		if int64(content) != content64 || !ok || !block.open {
			return nil, fmt.Errorf("%w: text done for unknown or stopped text", ErrInvalidSequence)
		}
		if block.textDone {
			return nil, fmt.Errorf("%w: duplicate text done", ErrInvalidSequence)
		}
		if _, ok := str(o["text"]); !ok {
			return nil, at("event.text", ErrInvalidWireData)
		}
		block.textDone = true
		s.blocks[key] = block
		return nil, nil
	case "response.completed", "response.incomplete":
		if !s.started {
			return nil, fmt.Errorf("%w: terminal event before start", ErrInvalidSequence)
		}
		r, _ := obj(o["response"])
		if r == nil {
			r = o
		}
		terminalID, ok := str(r["id"])
		if !ok || terminalID == "" {
			return nil, at("event.response.id", ErrInvalidWireData)
		}
		if terminalID != s.sourceID {
			return nil, at("event.response.id", fmt.Errorf("%w: response ID mismatch", ErrInvalidSequence))
		}
		status, ok := str(r["status"])
		expected := "completed"
		if typ == "response.incomplete" {
			expected = "incomplete"
		}
		if !ok || status != expected {
			return nil, at("event.response.status", fmt.Errorf("%w: terminal event/status mismatch", ErrInvalidWireData))
		}
		hasTool := false
		for key, block := range s.blocks {
			if key.kind == responseToolBlock {
				hasTool = true
			}
			if !block.open {
				continue
			}
			if status == "incomplete" && key.kind == responseToolBlock {
				return nil, at("event.response.output", fmt.Errorf("%w: incomplete response contains a truncated function call", ErrUnsupported))
			}
			return nil, fmt.Errorf("%w: terminal event with open output", ErrInvalidSequence)
		}
		reason := "end_turn"
		if hasTool {
			reason = "tool_use"
		}
		if status == "incomplete" {
			details, ok := obj(r["incomplete_details"])
			if !ok {
				return nil, at("event.response.incomplete_details", fmt.Errorf("%w: incomplete terminal requires details", ErrInvalidWireData))
			}
			why, ok := str(details["reason"])
			if !ok || why == "" {
				return nil, at("event.response.incomplete_details.reason", ErrInvalidWireData)
			}
			if why != "max_output_tokens" {
				return nil, at("event.response.incomplete_details.reason", ErrUnsupported)
			}
			reason = "max_tokens"
		} else if details, exists := r["incomplete_details"]; exists && details != nil {
			return nil, at("event.response.incomplete_details", fmt.Errorf("%w: completed terminal cannot have incomplete details", ErrInvalidWireData))
		}
		usage := map[string]any{}
		if u, ok := obj(r["usage"]); ok {
			au, err := openAIUsageToAnthropic(u, "event.response.usage")
			if err != nil {
				return nil, err
			}
			if s.inputUsageKnown && (au["input_tokens"].(int64) != s.inputTokens || au["cache_read_input_tokens"].(int64) != s.cacheReadTokens) {
				return nil, at("event.response.usage", fmt.Errorf("%w: terminal input usage does not match response.created", ErrInvalidSequence))
			}
			usage["output_tokens"] = au["output_tokens"]
		} else {
			return nil, at("event.response.usage", fmt.Errorf("%w: terminal usage object required", ErrInvalidWireData))
		}
		d1, _ := event("message_delta", map[string]any{"delta": map[string]any{"stop_reason": reason, "stop_sequence": nil}, "usage": usage})
		d2, _ := event("message_stop", map[string]any{})
		s.terminal = true
		return []SSEEvent{d1, d2}, nil
	case "error", "response.failed":
		var source map[string]any
		if typ == "response.failed" {
			r, ok := obj(o["response"])
			if !ok {
				return nil, at("event.response", ErrInvalidWireData)
			}
			source, ok = obj(r["error"])
			if !ok {
				return nil, at("event.response.error", ErrInvalidWireData)
			}
		} else if nested, ok := obj(o["error"]); ok {
			source = nested // tolerate the older nested OpenAI error representation.
		} else {
			source = o
		}
		message, err := requiredString(source, "message", "event.error")
		if err != nil {
			return nil, err
		}
		code, _ := str(source["code"])
		if code == "" {
			code = "openai_error"
		}
		// Anthropic requires every started content block to be stopped. Preserve
		// that lifecycle even when Responses fails midway through an item.
		keys := make([]responseBlockKey, 0, len(s.blocks))
		for key, block := range s.blocks {
			if block.open {
				keys = append(keys, key)
			}
		}
		sort.Slice(keys, func(i, j int) bool {
			return s.blocks[keys[i]].index < s.blocks[keys[j]].index
		})
		out := make([]SSEEvent, 0, len(keys)+1)
		for _, key := range keys {
			stops, closeErr := s.close(key)
			if closeErr != nil {
				return nil, closeErr
			}
			out = append(out, stops...)
		}
		errEvent, marshalErr := event("error", map[string]any{"error": map[string]any{"type": code, "message": message}})
		if marshalErr != nil {
			return nil, marshalErr
		}
		out = append(out, errEvent)
		s.terminal = true
		return out, nil
	case "response.refusal.delta", "response.refusal.done", "response.output_text.annotation.added", "response.output_text.annotation.delta", "response.output_text.annotation.done", "response.reasoning_summary_part.added", "response.reasoning_summary_part.done", "response.reasoning_summary_text.delta", "response.reasoning_summary_text.done":
		return nil, at("event "+typ, ErrUnsupported)
	default:
		return nil, at("event "+typ, ErrUnsupported)
	}
}

// Finish validates that a terminal event was observed. It never fabricates completion.
func (s *ResponsesToAnthropicStream) Finish() error {
	if !s.terminal {
		return ErrTruncatedStream
	}
	return nil
}

type anthropicBlock struct {
	kind, id, itemID, name string
	index, output          int
	text, args             string
}

// AnthropicToResponsesStream incrementally translates complete Anthropic
// Messages SSE events into complete Responses API SSE events.
type AnthropicToResponsesStream struct {
	options                                    ConversionOptions
	model, id                                  string
	started, terminal                          bool
	blocks                                     map[int]*anthropicBlock
	outputs                                    []any
	inputTokens, cacheReadTokens, outputTokens int64
	stop                                       string
	nextOutput                                 int
	sequence                                   int64
	createdAt                                  int64
	messageDelta                               bool
}

// NewAnthropicToResponsesStream constructs an Anthropic-to-Responses stream converter.
func NewAnthropicToResponsesStream(model string) *AnthropicToResponsesStream {
	return NewAnthropicToResponsesStreamWithOptions(model, ConversionOptions{})
}

// NewAnthropicToResponsesStreamWithOptions constructs a converter with per-stream warning options.
func NewAnthropicToResponsesStreamWithOptions(model string, options ConversionOptions) *AnthropicToResponsesStream {
	return &AnthropicToResponsesStream{model: model, options: options, blocks: map[int]*anthropicBlock{}}
}
func (s *AnthropicToResponsesStream) response(status string, terminal bool) map[string]any {
	r := map[string]any{
		"id": s.id, "object": "response", "created_at": s.createdAt,
		"completed_at": nil, "status": status, "error": nil,
		"incomplete_details": nil, "instructions": nil,
		"max_output_tokens": nil, "max_tool_calls": nil, "metadata": map[string]any{},
		"model": s.model, "output": s.outputs, "parallel_tool_calls": true,
		"previous_response_id": nil, "prompt_cache_key": nil, "prompt_cache_retention": nil,
		"reasoning": nil, "safety_identifier": nil, "service_tier": "default",
		"store": true, "temperature": 1.0,
		"text":        map[string]any{"format": map[string]any{"type": "text"}},
		"tool_choice": "auto", "tools": []any{}, "top_logprobs": 0,
		"top_p": 1.0, "truncation": "disabled", "user": nil,
	}
	if terminal {
		r["completed_at"] = time.Now().Unix()
		r["usage"] = map[string]any{"input_tokens": s.inputTokens, "output_tokens": s.outputTokens, "total_tokens": s.inputTokens + s.outputTokens, "input_tokens_details": map[string]any{"cached_tokens": s.cacheReadTokens}, "output_tokens_details": map[string]any{"reasoning_tokens": int64(0)}}
	} else {
		r["usage"] = nil
	}
	if status == "incomplete" {
		r["incomplete_details"] = map[string]any{"reason": "max_output_tokens"}
	}
	return r
}

func (s *AnthropicToResponsesStream) emit(typ string, v map[string]any) (SSEEvent, error) {
	v["sequence_number"] = s.sequence
	s.sequence++
	return event(typ, v)
}

func (s *AnthropicToResponsesStream) one(typ string, v map[string]any) ([]SSEEvent, error) {
	e, err := s.emit(typ, v)
	if err != nil {
		return nil, err
	}
	return []SSEEvent{e}, nil
}

// Convert translates one complete event and may return zero or more events.
func (s *AnthropicToResponsesStream) Convert(e SSEEvent) ([]SSEEvent, error) {
	decoded := decodeStreamEvent(e)
	if decoded.err == nil && !decoded.done {
		inspectDecodedAnthropicEvent(decoded.object, decoded.typ, s.options)
	}
	before := *s
	before.blocks = cloneBlockMap(s.blocks)
	before.outputs = append([]any(nil), s.outputs...)
	out, err := s.convert(e, decoded)
	if err != nil {
		*s = before
		return nil, err
	}
	return out, nil
}

func (s *AnthropicToResponsesStream) convert(e SSEEvent, decoded decodedStreamEvent) ([]SSEEvent, error) {
	if s.terminal {
		return nil, fmt.Errorf("%w: event after terminal event", ErrInvalidSequence)
	}
	o, typ, err := streamObject(e, decoded)
	if err != nil {
		return nil, err
	}
	if s.messageDelta && typ != "message_stop" && typ != "error" && typ != "ping" {
		return nil, fmt.Errorf("%w: event after message_delta", ErrInvalidSequence)
	}
	switch typ {
	case "message_start":
		if s.started {
			return nil, fmt.Errorf("%w: duplicate message start", ErrInvalidSequence)
		}
		m, ok := obj(o["message"])
		if !ok {
			return nil, at("event.message", ErrInvalidWireData)
		}
		sourceID, _ := str(m["id"])
		if sourceID == "" {
			return nil, at("event.message.id", ErrInvalidWireData)
		}
		s.id = destinationID("resp", sourceID, 0)
		s.createdAt = time.Now().Unix()
		if s.model == "" {
			s.model, _ = str(m["model"])
		}
		if s.model == "" {
			return nil, at("event.message.model", ErrInvalidWireData)
		}
		if u, ok := obj(m["usage"]); ok {
			ru, err := anthropicUsageToOpenAI(u, "event.message.usage")
			if err != nil {
				return nil, err
			}
			s.inputTokens = ru["input_tokens"].(int64)
			s.outputTokens = ru["output_tokens"].(int64)
			s.cacheReadTokens = ru["input_tokens_details"].(map[string]any)["cached_tokens"].(int64)
		} else if m["usage"] != nil {
			return nil, at("event.message.usage", ErrInvalidWireData)
		}
		s.started = true
		e1, _ := s.emit("response.created", map[string]any{"response": s.response("in_progress", false)})
		e2, _ := s.emit("response.in_progress", map[string]any{"response": s.response("in_progress", false)})
		return []SSEEvent{e1, e2}, nil
	case "content_block_start":
		if !s.started {
			return nil, fmt.Errorf("%w: block before message start", ErrInvalidSequence)
		}
		idx, err := eventIndex(o)
		if err != nil {
			return nil, err
		}
		if _, ok := s.blocks[idx]; ok {
			return nil, fmt.Errorf("%w: duplicate block index", ErrInvalidSequence)
		}
		c, ok := obj(o["content_block"])
		if !ok {
			return nil, at("event.content_block", ErrInvalidWireData)
		}
		if err := rejectPresent(c, "event.content_block", "phase"); err != nil {
			return nil, err
		}
		// Validate the complete start event before reserving an output/index. This
		// prevents a rejected event from corrupting the converter's FSM state.
		kind, _ := str(c["type"])
		b := &anthropicBlock{index: idx}
		switch kind {
		case "text":
			text, ok := str(c["text"])
			if !ok {
				return nil, at("event.content_block.text", ErrInvalidWireData)
			}
			if err := rejectNonEmptyArray(c, "citations", "event.content_block"); err != nil {
				return nil, err
			}
			b.kind, b.text = "text", text
		case "tool_use":
			if input, ok := obj(c["input"]); !ok || len(input) != 0 {
				return nil, at("event.content_block.input", fmt.Errorf("%w: streaming tool block must start with empty input", ErrInvalidWireData))
			}
			b.kind = "tool"
			b.id, _ = str(c["id"])
			b.name, _ = str(c["name"])
			if b.id == "" || b.name == "" {
				return nil, at("event.content_block", ErrInvalidWireData)
			}
		default:
			return nil, at("event.content_block.type", ErrUnsupported)
		}

		out := s.nextOutput
		s.nextOutput++
		b.output = out
		s.blocks[idx] = b
		if b.kind == "text" {
			b.id = destinationID("msg", s.id, out)
			item := map[string]any{"id": b.id, "type": "message", "status": "in_progress", "role": "assistant", "content": []any{}}
			e1, _ := s.emit("response.output_item.added", map[string]any{"output_index": out, "item": item})
			// Non-empty initial text is legal in Anthropic streams. Carry it in the
			// added part and final item; do not invent an empty delta.
			e2, _ := s.emit("response.content_part.added", map[string]any{"item_id": b.id, "output_index": out, "content_index": 0, "part": map[string]any{"type": "output_text", "text": b.text, "annotations": []any{}, "logprobs": []any{}}})
			return []SSEEvent{e1, e2}, nil
		}
		b.itemID = destinationID("fc", s.id, out)
		item := map[string]any{"id": b.itemID, "call_id": b.id, "type": "function_call", "name": b.name, "arguments": "", "status": "in_progress"}
		return s.one("response.output_item.added", map[string]any{"output_index": out, "item": item})
	case "content_block_delta":
		idx, err := eventIndex(o)
		if err != nil {
			return nil, err
		}
		b, ok := s.blocks[idx]
		if !ok {
			return nil, fmt.Errorf("%w: delta for unknown block", ErrInvalidSequence)
		}
		d, ok := obj(o["delta"])
		if !ok {
			return nil, at("event.delta", ErrInvalidWireData)
		}
		if b.kind == "text" {
			if d["type"] == "citations_delta" || d["type"] == "citation_delta" {
				return nil, at("event.delta.type", ErrUnsupported)
			}
			if d["type"] != "text_delta" {
				return nil, at("event.delta.type", ErrInvalidWireData)
			}
			if err := rejectNonEmptyArray(d, "citations", "event.delta"); err != nil {
				return nil, err
			}
			x, ok := str(d["text"])
			if !ok {
				return nil, at("event.delta.text", ErrInvalidWireData)
			}
			b.text += x
			return s.one("response.output_text.delta", map[string]any{"item_id": b.id, "output_index": b.output, "content_index": 0, "delta": x})
		}
		if d["type"] != "input_json_delta" {
			return nil, at("event.delta.type", ErrInvalidWireData)
		}
		x, ok := str(d["partial_json"])
		if !ok {
			return nil, at("event.delta.partial_json", ErrInvalidWireData)
		}
		b.args += x
		return s.one("response.function_call_arguments.delta", map[string]any{"item_id": b.itemID, "output_index": b.output, "delta": x})
	case "content_block_stop":
		idx, err := eventIndex(o)
		if err != nil {
			return nil, err
		}
		b, ok := s.blocks[idx]
		if !ok {
			return nil, fmt.Errorf("%w: stop for unknown block", ErrInvalidSequence)
		}
		delete(s.blocks, idx)
		if b.kind == "text" {
			part := map[string]any{"type": "output_text", "text": b.text, "annotations": []any{}, "logprobs": []any{}}
			item := map[string]any{"id": b.id, "type": "message", "status": "completed", "role": "assistant", "content": []any{part}}
			for len(s.outputs) <= b.output {
				s.outputs = append(s.outputs, nil)
			}
			s.outputs[b.output] = item
			e1, _ := s.emit("response.output_text.done", map[string]any{"item_id": b.id, "output_index": b.output, "content_index": 0, "text": b.text})
			e2, _ := s.emit("response.content_part.done", map[string]any{"item_id": b.id, "output_index": b.output, "content_index": 0, "part": part})
			e3, _ := s.emit("response.output_item.done", map[string]any{"output_index": b.output, "item": item})
			return []SSEEvent{e1, e2, e3}, nil
		}
		if _, err := parseArguments(b.args, "event.arguments"); err != nil {
			return nil, err
		}
		item := map[string]any{"id": b.itemID, "call_id": b.id, "type": "function_call", "name": b.name, "arguments": b.args, "status": "completed"}
		for len(s.outputs) <= b.output {
			s.outputs = append(s.outputs, nil)
		}
		s.outputs[b.output] = item
		e1, _ := s.emit("response.function_call_arguments.done", map[string]any{"item_id": item["id"], "output_index": b.output, "arguments": b.args})
		e2, _ := s.emit("response.output_item.done", map[string]any{"output_index": b.output, "item": item})
		return []SSEEvent{e1, e2}, nil
	case "message_delta":
		if !s.started {
			return nil, fmt.Errorf("%w: message_delta before message_start", ErrInvalidSequence)
		}
		if s.messageDelta {
			return nil, fmt.Errorf("%w: duplicate message_delta", ErrInvalidSequence)
		}
		d, ok := obj(o["delta"])
		if !ok {
			return nil, at("event.delta", ErrInvalidWireData)
		}
		s.stop, _ = str(d["stop_reason"])
		switch s.stop {
		case "end_turn", "tool_use", "stop_sequence", "max_tokens":
		case "pause_turn", "refusal":
			return nil, at("event.delta.stop_reason", ErrUnsupported)
		default:
			return nil, at("event.delta.stop_reason", fmt.Errorf("%w: unknown stop reason", ErrInvalidWireData))
		}
		if u, ok := obj(o["usage"]); ok {
			v, exists := u["output_tokens"]
			if !exists {
				return nil, at("event.usage.output_tokens", ErrInvalidWireData)
			}
			n, err := token(v, "event.usage.output_tokens")
			if err != nil {
				return nil, err
			}
			if n > int64(^uint64(0)>>1)-s.inputTokens {
				return nil, at("event.usage.output_tokens", fmt.Errorf("%w: token total overflows int64", ErrInvalidWireData))
			}
			s.outputTokens = n
		} else {
			return nil, at("event.usage", ErrInvalidWireData)
		}
		s.messageDelta = true
		return nil, nil
	case "message_stop":
		if !s.started || !s.messageDelta {
			return nil, fmt.Errorf("%w: message_stop requires message_start and message_delta", ErrInvalidSequence)
		}
		if len(s.blocks) > 0 {
			return nil, fmt.Errorf("%w: message stop with open block", ErrInvalidSequence)
		}
		status := "completed"
		typ := "response.completed"
		if s.stop == "max_tokens" {
			status = "incomplete"
			typ = "response.incomplete"
		}
		hasTool := false
		for _, output := range s.outputs {
			if output == nil {
				return nil, fmt.Errorf("%w: output blocks completed out of order with a missing output", ErrInvalidSequence)
			}
			if output.(map[string]any)["type"] == "function_call" {
				hasTool = true
			}
		}
		if s.stop == "tool_use" && !hasTool {
			return nil, fmt.Errorf("%w: tool_use stop without tool output", ErrInvalidSequence)
		}
		if (s.stop == "end_turn" || s.stop == "stop_sequence") && hasTool {
			return nil, fmt.Errorf("%w: tool output without tool_use stop", ErrInvalidSequence)
		}
		s.terminal = true
		return s.one(typ, map[string]any{"response": s.response(status, true)})
	case "error":
		eo, ok := obj(o["error"])
		if !ok {
			return nil, at("event.error", ErrInvalidWireData)
		}
		message, err := requiredString(eo, "message", "event.error")
		if err != nil {
			return nil, err
		}
		code, _ := str(eo["type"])
		if code == "" {
			code, _ = str(eo["code"])
		}
		if code == "" {
			code = "anthropic_error"
		}
		s.terminal = true
		return s.one("error", map[string]any{"code": code, "message": message, "param": nil})
	case "ping":
		return nil, nil
	default:
		return nil, at("event "+typ, ErrUnsupported)
	}
}
func eventIndex(o map[string]any) (int, error) {
	i, err := integer(o["index"], "event.index", false)
	if err != nil {
		return 0, err
	}
	idx := int(i)
	if int64(idx) != i {
		return 0, at("event.index", fmt.Errorf("%w: index out of range", ErrInvalidWireData))
	}
	return idx, nil
}

// Finish validates that a terminal event was observed. It never fabricates completion.
func (s *AnthropicToResponsesStream) Finish() error {
	if !s.terminal {
		return ErrTruncatedStream
	}
	return nil
}
