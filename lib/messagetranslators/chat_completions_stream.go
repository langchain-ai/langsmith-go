package messagetranslators

import (
	"fmt"
	"sort"
)

type decodedChatCompletionsEvent struct {
	object map[string]any
	done   bool
	err    error
}

func decodeChatCompletionsEvent(e SSEEvent) decodedChatCompletionsEvent {
	if string(e.Data) == "[DONE]" {
		return decodedChatCompletionsEvent{done: true}
	}
	o, err := decodeObject(e.Data)
	return decodedChatCompletionsEvent{object: o, err: err}
}

func chatCompletionsChunkObject(e SSEEvent, decoded decodedChatCompletionsEvent) (map[string]any, bool, error) {
	if decoded.done {
		if e.Event != "" && e.Event != "message" {
			return nil, false, at("event", ErrInvalidWireData)
		}
		return nil, true, nil
	}
	if e.Event != "" && e.Event != "message" && e.Event != "error" {
		return nil, false, at("event "+e.Event, ErrUnsupported)
	}
	if decoded.err != nil {
		return nil, false, at("event", decoded.err)
	}
	if er, ok := obj(decoded.object["error"]); ok {
		return map[string]any{"_chat_completions_error": er}, false, nil
	}
	if decoded.object["object"] != "chat.completion.chunk" {
		return nil, false, at("event.object", ErrInvalidWireData)
	}
	return decoded.object, false, nil
}

type chatCompletionsStreamTool struct {
	chatCompletionsIndex, block int
	id, name                    string
	// args grows by append so a pre-event snapshot keeps its own length.
	args []byte
	open bool
}

// ChatCompletionsToAnthropicStream incrementally translates complete Chat
// Completions SSE events into complete Anthropic Messages SSE events.
type ChatCompletionsToAnthropicStream struct {
	cfg                                       config
	model                                     string
	started, role, finishSeen, terminal, text bool
	sourceID, id, sourceModel                 string
	created                                   int64
	nextBlock, textBlock                      int
	tools                                     map[int]*chatCompletionsStreamTool
	finish                                    string
	outputTokens                              int64
	inputTokens, cachedTokens                 int64
	startUsageKnown, terminalUsageKnown       bool
}

// NewChatCompletionsToAnthropicStream constructs a converter from Chat Completions to Anthropic Messages.
func NewChatCompletionsToAnthropicStream(model string, options ...Option) *ChatCompletionsToAnthropicStream {
	return &ChatCompletionsToAnthropicStream{model: model, cfg: newConfig(options), tools: map[int]*chatCompletionsStreamTool{}, textBlock: -1}
}

// Usage returns the token accounting observed so far, in Anthropic terms.
//
// A Chat Completions stream discloses input usage only in its optional terminal
// usage chunk, after the Anthropic message_start has been emitted. Reading it
// here recovers those counts without reparsing the source stream.
func (s *ChatCompletionsToAnthropicStream) Usage() Usage {
	return Usage{InputTokens: s.inputTokens, OutputTokens: s.outputTokens, CacheReadInputTokens: s.cachedTokens}
}

func (s *ChatCompletionsToAnthropicStream) applyUsage(v any, path string, terminal, compareOutput bool) (map[string]any, error) {
	if terminal && s.terminalUsageKnown {
		return nil, at(path, fmt.Errorf("%w: duplicate terminal usage", ErrInvalidSequence))
	}
	u, ok := obj(v)
	if !ok {
		return nil, at(path, ErrInvalidWireData)
	}
	usage, err := chatCompletionsUsageToAnthropic(u, path, s.cfg)
	if err != nil {
		return nil, err
	}
	input := usage["input_tokens"].(int64)
	cached := usage["cache_read_input_tokens"].(int64)
	output := usage["output_tokens"].(int64)
	if terminal {
		if s.startUsageKnown && (input != s.inputTokens || cached != s.cachedTokens || compareOutput && output != s.outputTokens) {
			return nil, at(path, fmt.Errorf("%w: usage changed", ErrInvalidSequence))
		}
		s.terminalUsageKnown = true
	} else {
		s.startUsageKnown = true
	}
	s.inputTokens, s.cachedTokens, s.outputTokens = input, cached, output
	return usage, nil
}

func (s *ChatCompletionsToAnthropicStream) start(o map[string]any) (SSEEvent, error) {
	id, err := requiredString(o, "id", "event")
	if err != nil {
		return SSEEvent{}, err
	}
	wireModel, err := requiredString(o, "model", "event")
	if err != nil {
		return SSEEvent{}, err
	}
	created, err := integer(o["created"], "event.created", false)
	if err != nil {
		return SSEEvent{}, err
	}
	s.sourceID, s.sourceModel, s.created = id, wireModel, created
	s.id = destinationID("msg", id, 0)
	dstModel := s.model
	if dstModel == "" {
		dstModel = wireModel
	}
	usage := map[string]any{"input_tokens": int64(0), "output_tokens": int64(0)}
	if o["usage"] != nil {
		usage, err = s.applyUsage(o["usage"], "event.usage", false, false)
		if err != nil {
			return SSEEvent{}, err
		}
	}
	s.started = true
	return event("message_start", map[string]any{"message": map[string]any{"id": s.id, "type": "message", "role": "assistant", "content": []any{}, "model": dstModel, "stop_reason": nil, "stop_sequence": nil, "usage": usage}}), nil
}

func (s *ChatCompletionsToAnthropicStream) validateIdentity(o map[string]any) error {
	id, err := requiredString(o, "id", "event")
	if err != nil {
		return err
	}
	m, err := requiredString(o, "model", "event")
	if err != nil {
		return err
	}
	created, err := integer(o["created"], "event.created", false)
	if err != nil {
		return err
	}
	if id != s.sourceID || m != s.sourceModel || created != s.created {
		return at("event", fmt.Errorf("%w: chunk identity changed", ErrInvalidSequence))
	}
	return nil
}

// Convert translates one complete Chat Completions SSE event. State is rolled
// back when validation fails, so callers never observe a partially accepted event.
func (s *ChatCompletionsToAnthropicStream) Convert(raw SSEEvent) ([]SSEEvent, error) {
	decoded := decodeChatCompletionsEvent(raw)
	if decoded.err == nil && !decoded.done {
		inspectChatCompletionsEvent(decoded.object, s.cfg)
	}
	before := *s
	before.tools = cloneMap(s.tools)
	for index, tool := range before.tools {
		copyTool := *tool
		before.tools[index] = &copyTool
	}
	out, err := s.convert(raw, decoded)
	if err != nil {
		*s = before
		return nil, err
	}
	return out, nil
}

func (s *ChatCompletionsToAnthropicStream) convert(raw SSEEvent, decoded decodedChatCompletionsEvent) ([]SSEEvent, error) {
	if s.terminal {
		return nil, at("event", fmt.Errorf("%w: event after terminal event", ErrInvalidSequence))
	}
	o, done, err := chatCompletionsChunkObject(raw, decoded)
	if err != nil {
		return nil, err
	}
	if done {
		if !s.started || !s.finishSeen {
			return nil, at("event", fmt.Errorf("%w: [DONE] before finish_reason", ErrInvalidSequence))
		}
		// Output usage only, matching the documented v0 contract. Input counts
		// learned from a terminal usage chunk cannot amend the already-emitted
		// message_start, so callers that need them read Usage instead.
		usage := map[string]any{"output_tokens": s.outputTokens}
		d1 := event("message_delta", map[string]any{"delta": map[string]any{"stop_reason": s.finish, "stop_sequence": nil}, "usage": usage})
		d2 := event("message_stop", map[string]any{})
		s.terminal = true
		s.cfg.reportUsage(s.Usage())
		return []SSEEvent{d1, d2}, nil
	}
	if er, ok := obj(o["_chat_completions_error"]); ok {
		message, e := requiredString(er, "message", "event.error")
		if e != nil {
			return nil, e
		}
		code, _ := str(er["code"])
		if code == "" {
			code = "openai_error"
		}
		out := []SSEEvent{}
		if s.started && !s.finishSeen {
			blocks := []int{}
			if s.text {
				blocks = append(blocks, s.textBlock)
			}
			for _, tool := range s.tools {
				if tool.open {
					blocks = append(blocks, tool.block)
				}
			}
			sort.Ints(blocks)
			for _, block := range blocks {
				out = append(out, event("content_block_stop", map[string]any{"index": block}))
			}
		}
		out = append(out, event("error", map[string]any{"error": map[string]any{"type": code, "message": message}}))
		s.terminal = true
		return out, nil
	}
	var start *SSEEvent
	if !s.started {
		x, e := s.start(o)
		if e != nil {
			return nil, e
		}
		start = &x
	} else if err := s.validateIdentity(o); err != nil {
		return nil, err
	}
	if err := rejectPresent(o, "event", "service_tier", "system_fingerprint"); err != nil {
		return nil, err
	}
	choices, ok := arr(o["choices"])
	if !ok {
		return nil, at("event.choices", ErrInvalidWireData)
	}
	out := []SSEEvent{}
	if start != nil {
		out = append(out, *start)
	}
	if s.finishSeen && len(choices) != 0 {
		return nil, at("event.choices", fmt.Errorf("%w: choice after finish_reason", ErrInvalidSequence))
	}
	if len(choices) == 0 {
		if !s.finishSeen {
			return nil, at("event.choices", fmt.Errorf("%w: usage-only chunk before finish", ErrInvalidSequence))
		}
		if _, err := s.applyUsage(o["usage"], "event.usage", true, false); err != nil {
			return nil, err
		}
		return out, nil
	}
	if len(choices) != 1 {
		return nil, at("event.choices", ErrUnsupported)
	}
	c, ok := obj(choices[0])
	if !ok {
		return nil, at("event.choices[0]", ErrInvalidWireData)
	}
	idx, err := integer(c["index"], "event.choices[0].index", false)
	if err != nil {
		return nil, err
	}
	if idx != 0 {
		return nil, at("event.choices[0].index", ErrUnsupported)
	}
	if v, ok := c["logprobs"]; ok && v != nil {
		return nil, at("event.choices[0].logprobs", ErrUnsupported)
	}
	d, ok := obj(c["delta"])
	if !ok {
		return nil, at("event.choices[0].delta", ErrInvalidWireData)
	}
	if err := rejectPresent(d, "event.choices[0].delta", "audio", "function_call", "refusal"); err != nil {
		return nil, err
	}
	if v, ok := d["role"]; ok {
		if v != "assistant" || s.role {
			return nil, at("event.choices[0].delta.role", fmt.Errorf("%w: invalid or duplicate role delta", ErrInvalidSequence))
		}
		s.role = true
	}
	if !s.role && (d["content"] != nil || d["tool_calls"] != nil) {
		return nil, at("event", fmt.Errorf("%w: content before assistant role delta", ErrInvalidSequence))
	}
	if v, ok := d["content"]; ok && v != nil {
		t, ok := str(v)
		if !ok {
			return nil, at("event.choices[0].delta.content", ErrInvalidWireData)
		}
		if t != "" {
			// Providers commonly put an empty content field beside the role delta.
			// It carries no content and must not manufacture an Anthropic block.
			if !s.text {
				s.text = true
				s.textBlock = s.nextBlock
				s.nextBlock++
				out = append(out, event("content_block_start", map[string]any{"index": s.textBlock, "content_block": map[string]any{"type": "text", "text": ""}}))
			}
			out = append(out, event("content_block_delta", map[string]any{"index": s.textBlock, "delta": map[string]any{"type": "text_delta", "text": t}}))
		}
	}
	if v, ok := d["tool_calls"]; ok {
		calls, ok := arr(v)
		if !ok || len(calls) == 0 {
			return nil, at("event.choices[0].delta.tool_calls", ErrInvalidWireData)
		}
		inEvent := map[int]bool{}
		for i, x := range calls {
			p := fmt.Sprintf("event.choices[0].delta.tool_calls[%d]", i)
			tc, ok := obj(x)
			if !ok {
				return nil, at(p, ErrInvalidWireData)
			}
			n, e := integer(tc["index"], p+".index", false)
			if e != nil {
				return nil, e
			}
			ci := int(n)
			if int64(ci) != n {
				return nil, at(p+".index", ErrInvalidWireData)
			}
			if inEvent[ci] {
				return nil, at(p+".index", ErrInvalidSequence)
			}
			inEvent[ci] = true
			tool := s.tools[ci]
			f, hasFunction := obj(tc["function"])
			if tool == nil {
				if ci != len(s.tools) {
					return nil, at(p+".index", fmt.Errorf("%w: new tool indexes must be contiguous", ErrInvalidSequence))
				}
				if tc["type"] != "function" {
					return nil, at(p+".type", ErrInvalidWireData)
				}
				id, e := requiredString(tc, "id", p)
				if e != nil {
					return nil, e
				}
				for _, existing := range s.tools {
					if existing.id == id {
						return nil, at(p+".id", fmt.Errorf("%w: duplicate tool call ID", ErrInvalidSequence))
					}
				}
				if !hasFunction {
					return nil, at(p+".function", ErrInvalidWireData)
				}
				name, e := requiredString(f, "name", p+".function")
				if e != nil {
					return nil, e
				}
				tool = &chatCompletionsStreamTool{chatCompletionsIndex: ci, block: s.nextBlock, id: id, name: name, open: true}
				s.nextBlock++
				s.tools[ci] = tool
				out = append(out, event("content_block_start", map[string]any{"index": tool.block, "content_block": map[string]any{"type": "tool_use", "id": id, "name": name, "input": map[string]any{}}}))
			} else {
				if id, ok := str(tc["id"]); ok && id != tool.id {
					return nil, at(p+".id", ErrInvalidSequence)
				}
				if typ, ok := str(tc["type"]); ok && typ != "function" {
					return nil, at(p+".type", ErrInvalidSequence)
				}
				if hasFunction {
					if name, ok := str(f["name"]); ok && name != "" && name != tool.name {
						return nil, at(p+".function.name", ErrInvalidSequence)
					}
				}
			}
			if hasFunction {
				if av, ok := f["arguments"]; ok {
					frag, ok := str(av)
					if !ok {
						return nil, at(p+".function.arguments", ErrInvalidWireData)
					}
					tool.args = append(tool.args, frag...)
					out = append(out, event("content_block_delta", map[string]any{"index": tool.block, "delta": map[string]any{"type": "input_json_delta", "partial_json": frag}}))
				}
			}
		}
	}
	if fr, exists := c["finish_reason"]; exists && fr != nil {
		if s.finishSeen {
			return nil, at("event.choices[0].finish_reason", ErrInvalidSequence)
		}
		f, ok := str(fr)
		if !ok {
			return nil, at("event.choices[0].finish_reason", ErrInvalidWireData)
		}
		s.finish, err = chatCompletionsFinishReasonToAnthropic(f, len(s.tools) > 0, "event.choices[0].finish_reason")
		if err != nil {
			return nil, err
		}
		// Validate every accumulated argument before emitting any stops.
		for _, tool := range s.tools {
			if _, e := parseArgumentsRaw(tool.args, "event.tool_calls.arguments"); e != nil {
				return nil, e
			}
		}
		if s.text {
			out = append(out, event("content_block_stop", map[string]any{"index": s.textBlock}))
		}
		blocks := make([]int, 0, len(s.tools))
		byBlock := map[int]*chatCompletionsStreamTool{}
		for _, tool := range s.tools {
			blocks = append(blocks, tool.block)
			byBlock[tool.block] = tool
		}
		sort.Ints(blocks)
		for _, b := range blocks {
			out = append(out, event("content_block_stop", map[string]any{"index": b}))
			byBlock[b].open = false
		}
		s.finishSeen = true
	}
	if o["usage"] != nil {
		if _, ok := obj(o["usage"]); !ok {
			return nil, at("event.usage", ErrInvalidWireData)
		}
		if !s.finishSeen {
			return nil, at("event.usage", fmt.Errorf("%w: usage before finish_reason", ErrInvalidSequence))
		}
		if _, err := s.applyUsage(o["usage"], "event.usage", true, true); err != nil {
			return nil, err
		}
	}
	return out, nil
}

// Finish validates that the Chat Completions stream ended with [DONE].
func (s *ChatCompletionsToAnthropicStream) Finish() error {
	if !s.terminal {
		return ErrTruncatedStream
	}
	return nil
}

// AnthropicToChatCompletionsStream incrementally translates complete Anthropic
// Messages SSE events into complete Chat Completions SSE events.
type AnthropicToChatCompletionsStream struct {
	cfg                                      config
	model, id, sourceID                      string
	created                                  int64
	started, messageDelta, terminal, sawTool bool
	blocks                                   map[int]*anthropicBlock
	toolIDs                                  map[string]bool
	nextTool                                 int
	inputTokens, cachedTokens, outputTokens  int64
	finish                                   string
}

// NewAnthropicToChatCompletionsStream constructs a converter from Anthropic Messages to Chat Completions.
func NewAnthropicToChatCompletionsStream(model string, options ...Option) *AnthropicToChatCompletionsStream {
	return &AnthropicToChatCompletionsStream{model: model, cfg: newConfig(options), blocks: map[int]*anthropicBlock{}, toolIDs: map[string]bool{}}
}

// Usage returns the token accounting observed so far, in Anthropic terms.
func (s *AnthropicToChatCompletionsStream) Usage() Usage {
	return Usage{InputTokens: s.inputTokens - s.cachedTokens, OutputTokens: s.outputTokens, CacheReadInputTokens: s.cachedTokens}
}
func (s *AnthropicToChatCompletionsStream) chunk(choices []any, usage any) SSEEvent {
	return chatCompletionsEvent(map[string]any{"id": s.id, "object": "chat.completion.chunk", "created": s.created, "model": s.model, "choices": choices, "usage": usage})
}

// Chat Completions streams are unnamed data-only events, so no Event field is set.
func chatCompletionsEvent(v map[string]any) SSEEvent { return SSEEvent{Data: mustMarshal(v)} }

func (s *AnthropicToChatCompletionsStream) choice(delta map[string]any, finish any) SSEEvent {
	return s.chunk([]any{map[string]any{"index": int64(0), "delta": delta, "finish_reason": finish, "logprobs": nil}}, nil)
}

// Convert translates one complete Anthropic Messages SSE event. Failed events
// do not advance the converter state.
func (s *AnthropicToChatCompletionsStream) Convert(raw SSEEvent) ([]SSEEvent, error) {
	decoded := decodeStreamEvent(raw)
	if decoded.err == nil && !decoded.done {
		inspectDecodedAnthropicEvent(decoded.object, decoded.typ, s.cfg)
	}
	before := *s
	before.blocks = cloneBlockMap(s.blocks)
	before.toolIDs = cloneMap(s.toolIDs)
	out, err := s.convert(raw, decoded)
	if err != nil {
		*s = before
		return nil, err
	}
	return out, nil
}

func (s *AnthropicToChatCompletionsStream) convert(raw SSEEvent, decoded decodedStreamEvent) ([]SSEEvent, error) {
	if s.terminal {
		return nil, at("event", fmt.Errorf("%w: event after terminal event", ErrInvalidSequence))
	}
	o, typ, err := streamObject(raw, decoded)
	if err != nil {
		return nil, err
	}
	if s.messageDelta && typ != "message_stop" && typ != "error" && typ != "ping" {
		return nil, at("event "+typ, fmt.Errorf("%w: event after message_delta", ErrInvalidSequence))
	}
	switch typ {
	case "message_start":
		if s.started {
			return nil, at("event "+typ, fmt.Errorf("%w: duplicate message_start", ErrInvalidSequence))
		}
		m, ok := obj(o["message"])
		if !ok {
			return nil, at("event.message", ErrInvalidWireData)
		}
		if m["type"] != "message" || m["role"] != "assistant" {
			return nil, at("event.message", ErrInvalidWireData)
		}
		content, ok := arr(m["content"])
		if !ok || len(content) != 0 {
			return nil, at("event.message.content", fmt.Errorf("%w: message_start content must be empty", ErrInvalidWireData))
		}
		if v, ok := m["stop_reason"]; ok && v != nil {
			return nil, at("event.message.stop_reason", ErrInvalidWireData)
		}
		if v, ok := m["stop_sequence"]; ok && v != nil {
			return nil, at("event.message.stop_sequence", ErrInvalidWireData)
		}
		sid, err := requiredString(m, "id", "event.message")
		if err != nil {
			return nil, err
		}
		wire, err := requiredString(m, "model", "event.message")
		if err != nil {
			return nil, err
		}
		s.sourceID = sid
		s.id = destinationID("chatcmpl", sid, 0)
		if s.model == "" {
			s.model = wire
		}
		s.created = s.cfg.now().Unix()
		if u, ok := obj(m["usage"]); ok {
			cu, e := anthropicUsageToChatCompletions(u, "event.message.usage", s.cfg)
			if e != nil {
				return nil, e
			}
			s.inputTokens = cu["prompt_tokens"].(int64)
			s.cachedTokens = cu["prompt_tokens_details"].(map[string]any)["cached_tokens"].(int64)
			s.outputTokens = cu["completion_tokens"].(int64)
		} else if m["usage"] != nil {
			return nil, at("event.message.usage", ErrInvalidWireData)
		}
		s.started = true
		return []SSEEvent{s.choice(map[string]any{"role": "assistant"}, nil)}, nil
	case "content_block_start":
		if !s.started {
			return nil, at("event "+typ, fmt.Errorf("%w: block before message_start", ErrInvalidSequence))
		}
		idx, e := eventIndex(o)
		if e != nil {
			return nil, e
		}
		if _, ok := s.blocks[idx]; ok {
			return nil, at("event "+typ, fmt.Errorf("%w: duplicate block", ErrInvalidSequence))
		}
		c, ok := obj(o["content_block"])
		if !ok {
			return nil, at("event.content_block", ErrInvalidWireData)
		}
		if err := rejectPresent(c, "event.content_block", "phase"); err != nil {
			return nil, err
		}
		b := &anthropicBlock{index: idx}
		switch c["type"] {
		case "text":
			if s.sawTool {
				return nil, at("event.content_block.type", fmt.Errorf("%w: text after tool use", ErrUnsupported))
			}
			t, ok := str(c["text"])
			if !ok {
				return nil, at("event.content_block.text", ErrInvalidWireData)
			}
			if err := rejectNonEmptyArray(c, "citations", "event.content_block"); err != nil {
				return nil, err
			}
			// Unlike the Responses converter, this one forwards each fragment as it
			// arrives and never needs the assembled text, so none is accumulated.
			b.kind = "text"
			s.blocks[idx] = b
			if t == "" {
				return nil, nil
			}
			return []SSEEvent{s.choice(map[string]any{"content": t}, nil)}, nil
		case "tool_use":
			input, ok := obj(c["input"])
			if !ok || len(input) != 0 {
				return nil, at("event.content_block.input", ErrInvalidWireData)
			}
			id, e := requiredString(c, "id", "event.content_block")
			if e != nil {
				return nil, e
			}
			if s.toolIDs[id] {
				return nil, at("event.content_block.id", fmt.Errorf("%w: duplicate tool call ID", ErrInvalidSequence))
			}
			name, e := requiredString(c, "name", "event.content_block")
			if e != nil {
				return nil, e
			}
			b.kind = "tool"
			b.id = id
			b.name = name
			b.output = s.nextTool
			s.nextTool++
			s.sawTool = true
			s.toolIDs[id] = true
			s.blocks[idx] = b
			return []SSEEvent{s.choice(map[string]any{"tool_calls": []any{map[string]any{"index": b.output, "id": id, "type": "function", "function": map[string]any{"name": name, "arguments": ""}}}}, nil)}, nil
		default:
			return nil, at("event.content_block.type", ErrUnsupported)
		}
	case "content_block_delta":
		idx, e := eventIndex(o)
		if e != nil {
			return nil, e
		}
		b, ok := s.blocks[idx]
		if !ok {
			return nil, at("event "+typ, fmt.Errorf("%w: delta for unknown block", ErrInvalidSequence))
		}
		d, ok := obj(o["delta"])
		if !ok {
			return nil, at("event.delta", ErrInvalidWireData)
		}
		if b.kind == "text" {
			if d["type"] != "text_delta" {
				if d["type"] == "citations_delta" {
					return nil, at("event.delta.type", ErrUnsupported)
				}
				return nil, at("event.delta.type", ErrInvalidWireData)
			}
			t, ok := str(d["text"])
			if !ok {
				return nil, at("event.delta.text", ErrInvalidWireData)
			}
			return []SSEEvent{s.choice(map[string]any{"content": t}, nil)}, nil
		}
		if d["type"] != "input_json_delta" {
			return nil, at("event.delta.type", ErrInvalidWireData)
		}
		frag, ok := str(d["partial_json"])
		if !ok {
			return nil, at("event.delta.partial_json", ErrInvalidWireData)
		}
		b.args = append(b.args, frag...)
		return []SSEEvent{s.choice(map[string]any{"tool_calls": []any{map[string]any{"index": b.output, "function": map[string]any{"arguments": frag}}}}, nil)}, nil
	case "content_block_stop":
		idx, e := eventIndex(o)
		if e != nil {
			return nil, e
		}
		b, ok := s.blocks[idx]
		if !ok {
			return nil, at("event "+typ, fmt.Errorf("%w: stop for unknown block", ErrInvalidSequence))
		}
		if b.kind == "tool" {
			if _, e := parseArgumentsRaw(b.args, "event.arguments"); e != nil {
				return nil, e
			}
		}
		delete(s.blocks, idx)
		return nil, nil
	case "message_delta":
		if !s.started || s.messageDelta {
			return nil, at("event "+typ, fmt.Errorf("%w: invalid message_delta", ErrInvalidSequence))
		}
		if len(s.blocks) > 0 {
			return nil, at("event "+typ, fmt.Errorf("%w: message_delta with open blocks", ErrInvalidSequence))
		}
		d, ok := obj(o["delta"])
		if !ok {
			return nil, at("event.delta", ErrInvalidWireData)
		}
		stop, ok := str(d["stop_reason"])
		if !ok {
			return nil, at("event.delta.stop_reason", ErrInvalidWireData)
		}
		if sequence, exists := d["stop_sequence"]; exists {
			if stop == "stop_sequence" {
				if value, ok := str(sequence); !ok || value == "" {
					return nil, at("event.delta.stop_sequence", ErrInvalidWireData)
				}
			} else if sequence != nil {
				return nil, at("event.delta.stop_sequence", ErrInvalidWireData)
			}
		}
		s.finish, err = anthropicStopReasonToChatCompletions(stop, s.sawTool, "event.delta.stop_reason")
		if err != nil {
			return nil, err
		}
		u, ok := obj(o["usage"])
		if !ok {
			return nil, at("event.usage", ErrInvalidWireData)
		}
		n, e := token(u["output_tokens"], "event.usage.output_tokens")
		if e != nil {
			return nil, e
		}
		s.outputTokens = n
		s.messageDelta = true
		return []SSEEvent{s.choice(map[string]any{}, s.finish)}, nil
	case "message_stop":
		if !s.started || !s.messageDelta {
			return nil, at("event "+typ, fmt.Errorf("%w: message_stop before message_delta", ErrInvalidSequence))
		}
		usage := map[string]any{"prompt_tokens": s.inputTokens, "completion_tokens": s.outputTokens, "total_tokens": s.inputTokens + s.outputTokens, "prompt_tokens_details": map[string]any{"cached_tokens": s.cachedTokens}}
		x := s.chunk([]any{}, usage)
		s.terminal = true
		s.cfg.reportUsage(s.Usage())
		return []SSEEvent{x, {Data: []byte("[DONE]")}}, nil
	case "error":
		er, ok := obj(o["error"])
		if !ok {
			return nil, at("event.error", ErrInvalidWireData)
		}
		message, e := requiredString(er, "message", "event.error")
		if e != nil {
			return nil, e
		}
		code, _ := str(er["type"])
		if code == "" {
			code = "anthropic_error"
		}
		x := chatCompletionsEvent(map[string]any{"error": map[string]any{"message": message, "type": code, "code": code}})
		s.terminal = true
		return []SSEEvent{x}, nil
	case "ping":
		return nil, nil
	default:
		return nil, at("event "+typ, ErrUnsupported)
	}
}

// Finish validates that an Anthropic terminal event was observed.
func (s *AnthropicToChatCompletionsStream) Finish() error {
	if !s.terminal {
		return ErrTruncatedStream
	}
	return nil
}
