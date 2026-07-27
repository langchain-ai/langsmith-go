package messagetranslators

import (
	"encoding/json"
	"fmt"
	"sort"
	"time"
)

func chatCompletionsChunkObject(e SSEEvent) (map[string]any, bool, error) {
	if string(e.Data) == "[DONE]" {
		if e.Event != "" && e.Event != "message" {
			return nil, false, at("event", ErrInvalidWireData)
		}
		return nil, true, nil
	}
	if e.Event != "" && e.Event != "message" && e.Event != "error" {
		return nil, false, at("event "+e.Event, ErrUnsupported)
	}
	o, err := decodeObject(e.Data)
	if err != nil {
		return nil, false, at("event", err)
	}
	if er, ok := obj(o["error"]); ok {
		return map[string]any{"_chat_completions_error": er}, false, nil
	}
	if o["object"] != "chat.completion.chunk" {
		return nil, false, at("event.object", ErrInvalidWireData)
	}
	return o, false, nil
}

type chatCompletionsStreamTool struct {
	chatCompletionsIndex, block int
	id, name, args              string
	open                        bool
}

// ChatCompletionsToAnthropicStream incrementally translates complete Chat
// Completions SSE events into complete Anthropic Messages SSE events.
type ChatCompletionsToAnthropicStream struct {
	options                                   ConversionOptions
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
func NewChatCompletionsToAnthropicStream(model string) *ChatCompletionsToAnthropicStream {
	return NewChatCompletionsToAnthropicStreamWithOptions(model, ConversionOptions{})
}

// NewChatCompletionsToAnthropicStreamWithOptions constructs a converter with per-stream warning options.
func NewChatCompletionsToAnthropicStreamWithOptions(model string, options ConversionOptions) *ChatCompletionsToAnthropicStream {
	return &ChatCompletionsToAnthropicStream{model: model, options: options, tools: map[int]*chatCompletionsStreamTool{}, textBlock: -1}
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
	if u, ok := obj(o["usage"]); ok {
		au, e := chatCompletionsUsageToAnthropic(u, "event.usage")
		if e != nil {
			return SSEEvent{}, e
		}
		usage = au
		s.inputTokens = au["input_tokens"].(int64)
		s.cachedTokens = au["cache_read_input_tokens"].(int64)
		s.outputTokens = au["output_tokens"].(int64)
		s.startUsageKnown = true
	} else if o["usage"] != nil {
		return SSEEvent{}, at("event.usage", ErrInvalidWireData)
	}
	e, err := event("message_start", map[string]any{"message": map[string]any{"id": s.id, "type": "message", "role": "assistant", "content": []any{}, "model": dstModel, "stop_reason": nil, "stop_sequence": nil, "usage": usage}})
	if err == nil {
		s.started = true
	}
	return e, err
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
func (s *ChatCompletionsToAnthropicStream) Convert(e SSEEvent) ([]SSEEvent, error) {
	inspectChatCompletionsEvent(e, s.options)
	before := *s
	before.tools = make(map[int]*chatCompletionsStreamTool, len(s.tools))
	for index, tool := range s.tools {
		copyTool := *tool
		before.tools[index] = &copyTool
	}
	out, err := s.convert(e)
	if err != nil {
		*s = before
		return nil, err
	}
	return out, nil
}

func (s *ChatCompletionsToAnthropicStream) convert(e SSEEvent) ([]SSEEvent, error) {
	if s.terminal {
		return nil, fmt.Errorf("%w: event after terminal event", ErrInvalidSequence)
	}
	o, done, err := chatCompletionsChunkObject(e)
	if err != nil {
		return nil, err
	}
	if done {
		if !s.started || !s.finishSeen {
			return nil, fmt.Errorf("%w: [DONE] before finish_reason", ErrInvalidSequence)
		}
		usage := map[string]any{"output_tokens": s.outputTokens}
		d1, _ := event("message_delta", map[string]any{"delta": map[string]any{"stop_reason": s.finish, "stop_sequence": nil}, "usage": usage})
		d2, _ := event("message_stop", map[string]any{})
		s.terminal = true
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
				x, _ := event("content_block_stop", map[string]any{"index": block})
				out = append(out, x)
			}
		}
		x, _ := event("error", map[string]any{"error": map[string]any{"type": code, "message": message}})
		out = append(out, x)
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
		if s.terminalUsageKnown {
			return nil, at("event.usage", fmt.Errorf("%w: duplicate terminal usage", ErrInvalidSequence))
		}
		u, ok := obj(o["usage"])
		if !ok {
			return nil, at("event.usage", ErrInvalidWireData)
		}
		au, e := chatCompletionsUsageToAnthropic(u, "event.usage")
		if e != nil {
			return nil, e
		}
		if s.startUsageKnown && (au["input_tokens"] != s.inputTokens || au["cache_read_input_tokens"] != s.cachedTokens) {
			return nil, at("event.usage", fmt.Errorf("%w: usage changed", ErrInvalidSequence))
		}
		s.inputTokens = au["input_tokens"].(int64)
		s.cachedTokens = au["cache_read_input_tokens"].(int64)
		s.outputTokens = au["output_tokens"].(int64)
		s.terminalUsageKnown = true
		return out, nil
	}
	if len(choices) != 1 {
		return nil, at("event.choices", ErrUnsupported)
	}
	c, ok := obj(choices[0])
	if !ok {
		return nil, at("event.choices[0]", ErrInvalidWireData)
	}
	idx, e2 := integer(c["index"], "event.choices[0].index", false)
	if e2 != nil || idx != 0 {
		if e2 != nil {
			return nil, e2
		}
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
		return nil, fmt.Errorf("%w: content before assistant role delta", ErrInvalidSequence)
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
				x, _ := event("content_block_start", map[string]any{"index": s.textBlock, "content_block": map[string]any{"type": "text", "text": ""}})
				out = append(out, x)
			}
			x, _ := event("content_block_delta", map[string]any{"index": s.textBlock, "delta": map[string]any{"type": "text_delta", "text": t}})
			out = append(out, x)
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
				x, _ := event("content_block_start", map[string]any{"index": tool.block, "content_block": map[string]any{"type": "tool_use", "id": id, "name": name, "input": map[string]any{}}})
				out = append(out, x)
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
					tool.args += frag
					x, _ := event("content_block_delta", map[string]any{"index": tool.block, "delta": map[string]any{"type": "input_json_delta", "partial_json": frag}})
					out = append(out, x)
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
		hasTools := len(s.tools) > 0
		switch f {
		case "stop":
			if hasTools {
				return nil, at("event.choices[0].finish_reason", ErrInvalidSequence)
			}
			s.finish = "end_turn"
		case "tool_calls":
			if !hasTools {
				return nil, at("event.choices[0].finish_reason", ErrInvalidSequence)
			}
			s.finish = "tool_use"
		case "length":
			s.finish = "max_tokens"
		case "content_filter", "function_call":
			return nil, at("event.choices[0].finish_reason", ErrUnsupported)
		default:
			return nil, at("event.choices[0].finish_reason", ErrInvalidWireData)
		}
		// Validate every accumulated argument before emitting any stops.
		for _, tool := range s.tools {
			if _, e := parseArguments(tool.args, "event.tool_calls.arguments"); e != nil {
				return nil, e
			}
		}
		if s.text {
			x, _ := event("content_block_stop", map[string]any{"index": s.textBlock})
			out = append(out, x)
		}
		blocks := make([]int, 0, len(s.tools))
		byBlock := map[int]*chatCompletionsStreamTool{}
		for _, tool := range s.tools {
			blocks = append(blocks, tool.block)
			byBlock[tool.block] = tool
		}
		sort.Ints(blocks)
		for _, b := range blocks {
			x, _ := event("content_block_stop", map[string]any{"index": b})
			out = append(out, x)
			byBlock[b].open = false
		}
		s.finishSeen = true
	}
	if u, ok := obj(o["usage"]); ok {
		if !s.finishSeen {
			return nil, at("event.usage", fmt.Errorf("%w: usage before finish_reason", ErrInvalidSequence))
		}
		if s.terminalUsageKnown {
			return nil, at("event.usage", fmt.Errorf("%w: duplicate terminal usage", ErrInvalidSequence))
		}
		au, e := chatCompletionsUsageToAnthropic(u, "event.usage")
		if e != nil {
			return nil, e
		}
		if s.startUsageKnown && (au["input_tokens"] != s.inputTokens || au["cache_read_input_tokens"] != s.cachedTokens || au["output_tokens"] != s.outputTokens) {
			return nil, at("event.usage", fmt.Errorf("%w: usage changed", ErrInvalidSequence))
		}
		s.inputTokens = au["input_tokens"].(int64)
		s.cachedTokens = au["cache_read_input_tokens"].(int64)
		s.outputTokens = au["output_tokens"].(int64)
		s.terminalUsageKnown = true
	} else if o["usage"] != nil {
		return nil, at("event.usage", ErrInvalidWireData)
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
	options                                  ConversionOptions
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
func NewAnthropicToChatCompletionsStream(model string) *AnthropicToChatCompletionsStream {
	return NewAnthropicToChatCompletionsStreamWithOptions(model, ConversionOptions{})
}

// NewAnthropicToChatCompletionsStreamWithOptions constructs a converter with per-stream warning options.
func NewAnthropicToChatCompletionsStreamWithOptions(model string, options ConversionOptions) *AnthropicToChatCompletionsStream {
	return &AnthropicToChatCompletionsStream{model: model, options: options, blocks: map[int]*anthropicBlock{}, toolIDs: map[string]bool{}}
}
func (s *AnthropicToChatCompletionsStream) chunk(choices []any, usage any) (SSEEvent, error) {
	return chatCompletionsEvent(map[string]any{"id": s.id, "object": "chat.completion.chunk", "created": s.created, "model": s.model, "choices": choices, "usage": usage})
}
func chatCompletionsEvent(v map[string]any) (SSEEvent, error) {
	b, e := json.Marshal(v)
	return SSEEvent{Data: b}, e
}
func (s *AnthropicToChatCompletionsStream) choice(delta map[string]any, finish any) (SSEEvent, error) {
	return s.chunk([]any{map[string]any{"index": int64(0), "delta": delta, "finish_reason": finish, "logprobs": nil}}, nil)
}

// Convert translates one complete Anthropic Messages SSE event. Failed events
// do not advance the converter state.
func (s *AnthropicToChatCompletionsStream) Convert(e SSEEvent) ([]SSEEvent, error) {
	inspectAnthropicEvent(e, s.options)
	before := *s
	before.blocks = make(map[int]*anthropicBlock, len(s.blocks))
	for index, block := range s.blocks {
		copyBlock := *block
		before.blocks[index] = &copyBlock
	}
	before.toolIDs = make(map[string]bool, len(s.toolIDs))
	for id := range s.toolIDs {
		before.toolIDs[id] = true
	}
	out, err := s.convert(e)
	if err != nil {
		*s = before
		return nil, err
	}
	return out, nil
}

func (s *AnthropicToChatCompletionsStream) convert(e SSEEvent) ([]SSEEvent, error) {
	if s.terminal {
		return nil, fmt.Errorf("%w: event after terminal event", ErrInvalidSequence)
	}
	o, typ, err := streamObject(e)
	if err != nil {
		return nil, err
	}
	if s.messageDelta && typ != "message_stop" && typ != "error" && typ != "ping" {
		return nil, fmt.Errorf("%w: event after message_delta", ErrInvalidSequence)
	}
	switch typ {
	case "message_start":
		if s.started {
			return nil, fmt.Errorf("%w: duplicate message_start", ErrInvalidSequence)
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
		s.created = time.Now().Unix()
		if u, ok := obj(m["usage"]); ok {
			cu, e := anthropicUsageToChatCompletions(u, "event.message.usage")
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
		x, _ := s.choice(map[string]any{"role": "assistant"}, nil)
		return []SSEEvent{x}, nil
	case "content_block_start":
		if !s.started {
			return nil, fmt.Errorf("%w: block before message_start", ErrInvalidSequence)
		}
		idx, e := eventIndex(o)
		if e != nil {
			return nil, e
		}
		if _, ok := s.blocks[idx]; ok {
			return nil, fmt.Errorf("%w: duplicate block", ErrInvalidSequence)
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
			b.kind = "text"
			b.text = t
			s.blocks[idx] = b
			if t == "" {
				return nil, nil
			}
			x, _ := s.choice(map[string]any{"content": t}, nil)
			return []SSEEvent{x}, nil
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
			x, _ := s.choice(map[string]any{"tool_calls": []any{map[string]any{"index": b.output, "id": id, "type": "function", "function": map[string]any{"name": name, "arguments": ""}}}}, nil)
			return []SSEEvent{x}, nil
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
			return nil, fmt.Errorf("%w: delta for unknown block", ErrInvalidSequence)
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
			b.text += t
			x, _ := s.choice(map[string]any{"content": t}, nil)
			return []SSEEvent{x}, nil
		}
		if d["type"] != "input_json_delta" {
			return nil, at("event.delta.type", ErrInvalidWireData)
		}
		frag, ok := str(d["partial_json"])
		if !ok {
			return nil, at("event.delta.partial_json", ErrInvalidWireData)
		}
		b.args += frag
		x, _ := s.choice(map[string]any{"tool_calls": []any{map[string]any{"index": b.output, "function": map[string]any{"arguments": frag}}}}, nil)
		return []SSEEvent{x}, nil
	case "content_block_stop":
		idx, e := eventIndex(o)
		if e != nil {
			return nil, e
		}
		b, ok := s.blocks[idx]
		if !ok {
			return nil, fmt.Errorf("%w: stop for unknown block", ErrInvalidSequence)
		}
		if b.kind == "tool" {
			if _, e := parseArguments(b.args, "event.arguments"); e != nil {
				return nil, e
			}
		}
		delete(s.blocks, idx)
		return nil, nil
	case "message_delta":
		if !s.started || s.messageDelta {
			return nil, fmt.Errorf("%w: invalid message_delta", ErrInvalidSequence)
		}
		if len(s.blocks) > 0 {
			return nil, fmt.Errorf("%w: message_delta with open blocks", ErrInvalidSequence)
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
		switch stop {
		case "end_turn", "stop_sequence":
			if s.sawTool {
				return nil, at("event.delta.stop_reason", ErrInvalidSequence)
			}
			s.finish = "stop"
		case "tool_use":
			if !s.sawTool {
				return nil, at("event.delta.stop_reason", ErrInvalidSequence)
			}
			s.finish = "tool_calls"
		case "max_tokens":
			s.finish = "length"
		case "pause_turn", "refusal":
			return nil, at("event.delta.stop_reason", ErrUnsupported)
		default:
			return nil, at("event.delta.stop_reason", ErrInvalidWireData)
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
		x, _ := s.choice(map[string]any{}, s.finish)
		return []SSEEvent{x}, nil
	case "message_stop":
		if !s.started || !s.messageDelta {
			return nil, fmt.Errorf("%w: message_stop before message_delta", ErrInvalidSequence)
		}
		usage := map[string]any{"prompt_tokens": s.inputTokens, "completion_tokens": s.outputTokens, "total_tokens": s.inputTokens + s.outputTokens, "prompt_tokens_details": map[string]any{"cached_tokens": s.cachedTokens}}
		x, _ := s.chunk([]any{}, usage)
		s.terminal = true
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
		x, _ := chatCompletionsEvent(map[string]any{"error": map[string]any{"message": message, "type": code, "code": code}})
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
