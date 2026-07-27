package messagetranslators

import (
	"encoding/json"
	"fmt"
	"strings"
)

// AnthropicRequestToChatCompletions converts an Anthropic Messages request to
// an OpenAI Chat Completions request. model overrides the wire model when non-empty.
func AnthropicRequestToChatCompletions(body []byte, model string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectAnthropicObject(a, false, cfg, "$")
	return anthropicRequestToChatCompletions(a, model, cfg)
}

func anthropicRequestToChatCompletions(a map[string]any, model string, cfg config) ([]byte, error) {
	if err := rejectUnsupportedFields(a, "$", anthropicRequestUnsupported); err != nil {
		return nil, err
	}
	if v, ok := a["stop_sequences"]; ok {
		ss, ok := arr(v)
		if !ok {
			return nil, at("$.stop_sequences", ErrInvalidWireData)
		}
		if len(ss) > 4 {
			// Not malformed Anthropic input: Chat Completions simply caps stop at 4.
			return nil, at("$.stop_sequences", fmt.Errorf("%w: Chat Completions accepts at most 4 stop sequences", ErrUnsupported))
		}
		for i, x := range ss {
			if s, ok := str(x); !ok || s == "" {
				return nil, at(fmt.Sprintf("$.stop_sequences[%d]", i), ErrInvalidWireData)
			}
		}
	}
	if err := validateOptionalBool(a, "stream", "$"); err != nil {
		return nil, err
	}
	for _, key := range []string{"temperature", "top_p"} {
		if v, ok := a[key]; ok {
			if err := numberInRange(v, "$."+key, 0, 1); err != nil {
				return nil, err
			}
		}
	}
	max, err := integer(a["max_tokens"], "$.max_tokens", true)
	if err != nil {
		return nil, err
	}
	dstModel, err := resolveModel(a, model, "$")
	if err != nil {
		return nil, err
	}
	ms, ok := arr(a["messages"])
	if !ok || len(ms) == 0 {
		return nil, at("$.messages", fmt.Errorf("%w: non-empty array required", ErrInvalidWireData))
	}
	r := map[string]any{"model": dstModel, "max_tokens": max}
	copyIf(r, a, "temperature", "temperature")
	copyIf(r, a, "top_p", "top_p")
	copyIf(r, a, "stream", "stream")
	if streaming, _ := a["stream"].(bool); streaming {
		r["stream_options"] = map[string]any{"include_usage": true}
	}
	copyIf(r, a, "stop", "stop_sequences")
	out := []any{}
	if s, e := anthropicSystem(a["system"]); e != nil {
		return nil, e
	} else if s != "" {
		out = append(out, map[string]any{"role": "system", "content": s})
	}
	pending, seen := map[string]bool{}, map[string]bool{}
	appendMessage := func(role string, content any) {
		out = append(out, map[string]any{"role": role, "content": content})
	}
	for i, x := range ms {
		p := fmt.Sprintf("$.messages[%d]", i)
		m, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if _, ok := m["name"]; ok {
			return nil, at(p+".name", ErrUnsupported)
		}
		role, err := requiredString(m, "role", p)
		if err != nil {
			return nil, err
		}
		if role != "user" && role != "assistant" {
			return nil, at(p+".role", ErrUnsupported)
		}
		var blocks []any
		if s, ok := str(m["content"]); ok {
			blocks = []any{map[string]any{"type": "text", "text": s}}
		} else if blocks, ok = arr(m["content"]); !ok || len(blocks) == 0 {
			return nil, at(p+".content", fmt.Errorf("%w: non-empty string or array required", ErrInvalidWireData))
		}
		if role == "assistant" {
			if len(pending) != 0 {
				return nil, at(p+".role", fmt.Errorf("%w: assistant message before all tool results", ErrInvalidSequence))
			}
			var text strings.Builder
			calls := []any{}
			sawTool := false
			for j, z := range blocks {
				q := fmt.Sprintf("%s.content[%d]", p, j)
				b, ok := obj(z)
				if !ok {
					return nil, at(q, ErrInvalidWireData)
				}
				if err := rejectPresent(b, q, "phase"); err != nil {
					return nil, err
				}
				switch b["type"] {
				case "text":
					if sawTool {
						return nil, at(q, fmt.Errorf("%w: Chat Completions orders assistant text before tool calls", ErrUnsupported))
					}
					if err := rejectNonEmptyArray(b, "citations", q); err != nil {
						return nil, err
					}
					t, ok := str(b["text"])
					if !ok {
						return nil, at(q+".text", ErrInvalidWireData)
					}
					text.WriteString(t)
				case "tool_use":
					sawTool = true
					id, e := requiredString(b, "id", q)
					if e != nil {
						return nil, e
					}
					if seen[id] {
						return nil, at(q+".id", fmt.Errorf("%w: duplicate tool ID", ErrInvalidSequence))
					}
					name, e := requiredString(b, "name", q)
					if e != nil {
						return nil, e
					}
					input := b["input"]
					if input == nil {
						input = map[string]any{}
					}
					if _, ok := obj(input); !ok {
						return nil, at(q+".input", ErrInvalidWireData)
					}
					args, _ := json.Marshal(input)
					calls = append(calls, map[string]any{"id": id, "type": "function", "function": map[string]any{"name": name, "arguments": string(args)}})
					seen[id], pending[id] = true, true
				default:
					return nil, at(q+".type", ErrUnsupported)
				}
			}
			cm := map[string]any{"role": "assistant"}
			if text.Len() > 0 {
				cm["content"] = text.String()
			} else {
				cm["content"] = nil
			}
			if len(calls) > 0 {
				cm["tool_calls"] = calls
			}
			out = append(out, cm)
			continue
		}
		var normal []any
		flush := func() {
			if len(normal) > 0 {
				appendMessage("user", normal)
				normal = nil
			}
		}
		for j, z := range blocks {
			q := fmt.Sprintf("%s.content[%d]", p, j)
			b, ok := obj(z)
			if !ok {
				return nil, at(q, ErrInvalidWireData)
			}
			if err := rejectPresent(b, q, "phase"); err != nil {
				return nil, err
			}
			switch b["type"] {
			case "text":
				if len(pending) != 0 {
					return nil, at(q, fmt.Errorf("%w: Chat Completions requires tool results before user content", ErrUnsupported))
				}
				if err := rejectNonEmptyArray(b, "citations", q); err != nil {
					return nil, err
				}
				t, ok := str(b["text"])
				if !ok {
					return nil, at(q+".text", ErrInvalidWireData)
				}
				normal = append(normal, map[string]any{"type": "text", "text": t})
			case "image":
				if len(pending) != 0 {
					return nil, at(q, fmt.Errorf("%w: Chat Completions requires tool results before user content", ErrUnsupported))
				}
				im, e := anthropicImageToChatCompletions(b, q)
				if e != nil {
					return nil, e
				}
				normal = append(normal, im)
			case "tool_result":
				flush()
				if v, ok := b["is_error"]; ok {
					bv, valid := v.(bool)
					if !valid {
						return nil, at(q+".is_error", ErrInvalidWireData)
					}
					if bv {
						return nil, at(q+".is_error", ErrUnsupported)
					}
				}
				id, e := requiredString(b, "tool_use_id", q)
				if e != nil {
					return nil, e
				}
				if !pending[id] {
					return nil, at(q+".tool_use_id", fmt.Errorf("%w: orphan or duplicate tool result", ErrInvalidSequence))
				}
				text, e := anthropicToolResultText(b["content"], q+".content")
				if e != nil {
					return nil, e
				}
				out = append(out, map[string]any{"role": "tool", "tool_call_id": id, "content": text})
				delete(pending, id)
			default:
				return nil, at(q+".type", ErrUnsupported)
			}
		}
		flush()
	}
	if len(pending) != 0 {
		return nil, at("$.messages", fmt.Errorf("%w: unresolved tool calls", ErrInvalidSequence))
	}
	r["messages"] = out
	if v, ok := a["tools"]; ok {
		r["tools"], err = anthropicToolsToChatCompletions(v, "$.tools")
		if err != nil {
			return nil, err
		}
	}
	if v, ok := a["tool_choice"]; ok {
		tc, ok := obj(v)
		if !ok {
			return nil, at("$.tool_choice", ErrInvalidWireData)
		}
		if _, ok := tc["disable_parallel_tool_use"]; ok {
			return nil, at("$.tool_choice.disable_parallel_tool_use", ErrUnsupported)
		}
		switch tc["type"] {
		case "auto":
			r["tool_choice"] = "auto"
		case "any":
			r["tool_choice"] = "required"
		case "none":
			r["tool_choice"] = "none"
		case "tool":
			name, e := requiredString(tc, "name", "$.tool_choice")
			if e != nil {
				return nil, e
			}
			r["tool_choice"] = map[string]any{"type": "function", "function": map[string]any{"name": name}}
		default:
			return nil, at("$.tool_choice.type", ErrUnsupported)
		}
	}
	return encode(r)
}

// ChatCompletionsRequestToAnthropic converts an OpenAI Chat Completions request
// to an Anthropic Messages request. Legacy functions/function_call are rejected.
func ChatCompletionsRequestToAnthropic(body []byte, model string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectChatCompletionsObject(r, false, cfg, "$")
	return chatCompletionsRequestToAnthropic(r, model, cfg)
}

func chatCompletionsRequestToAnthropic(r map[string]any, model string, cfg config) ([]byte, error) {
	if err := rejectUnsupportedFields(r, "$", chatCompletionsRequestUnsupported); err != nil {
		return nil, err
	}
	if v, ok := r["n"]; ok {
		n, e := integer(v, "$.n", true)
		if e != nil {
			return nil, e
		}
		if n != 1 {
			return nil, at("$.n", ErrUnsupported)
		}
	}
	if err := validateOptionalBool(r, "stream", "$"); err != nil {
		return nil, err
	}
	if v, ok := r["stream_options"]; ok {
		o, valid := obj(v)
		if !valid {
			return nil, at("$.stream_options", ErrInvalidWireData)
		}
		include, exists := o["include_usage"]
		if !exists {
			return nil, at("$.stream_options.include_usage", ErrInvalidWireData)
		}
		if _, valid := include.(bool); !valid {
			return nil, at("$.stream_options.include_usage", ErrInvalidWireData)
		}
		if len(o) != 1 {
			return nil, at("$.stream_options", ErrUnsupported)
		}
		streaming, _ := r["stream"].(bool)
		if !streaming {
			return nil, at("$.stream_options", fmt.Errorf("%w: stream_options requires stream=true", ErrInvalidWireData))
		}
	}
	for _, key := range []string{"temperature", "top_p"} {
		if v, ok := r[key]; ok {
			if e := numberInRange(v, "$."+key, 0, 1); e != nil {
				return nil, e
			}
		}
	}
	var max int64
	var err error
	_, hasMax := r["max_tokens"]
	_, hasCompletion := r["max_completion_tokens"]
	if hasMax && hasCompletion {
		return nil, at("$.max_completion_tokens", fmt.Errorf("%w: specify only one token limit", ErrInvalidWireData))
	}
	if hasCompletion {
		max, err = integer(r["max_completion_tokens"], "$.max_completion_tokens", true)
	} else {
		max, err = integer(r["max_tokens"], "$.max_tokens", true)
	}
	if err != nil {
		return nil, err
	}
	dstModel, err := resolveModel(r, model, "$")
	if err != nil {
		return nil, err
	}
	a := map[string]any{"model": dstModel, "max_tokens": max}
	copyIf(a, r, "temperature", "temperature")
	copyIf(a, r, "top_p", "top_p")
	copyIf(a, r, "stream", "stream")
	if v, ok := r["stop"]; ok {
		switch x := v.(type) {
		case string:
			if x == "" {
				return nil, at("$.stop", ErrInvalidWireData)
			}
			a["stop_sequences"] = []any{x}
		case []any:
			if len(x) > 4 {
				return nil, at("$.stop", ErrInvalidWireData)
			}
			for i, z := range x {
				if s, ok := str(z); !ok || s == "" {
					return nil, at(fmt.Sprintf("$.stop[%d]", i), ErrInvalidWireData)
				}
			}
			a["stop_sequences"] = x
		default:
			return nil, at("$.stop", ErrInvalidWireData)
		}
	}
	ms, ok := arr(r["messages"])
	if !ok || len(ms) == 0 {
		return nil, at("$.messages", fmt.Errorf("%w: non-empty array required", ErrInvalidWireData))
	}
	system := []any{}
	messages := []any{}
	appendMsg := func(role string, c []any) {
		if len(messages) > 0 {
			last := messages[len(messages)-1].(map[string]any)
			if last["role"] == role {
				last["content"] = append(last["content"].([]any), c...)
				return
			}
		}
		messages = append(messages, map[string]any{"role": role, "content": c})
	}
	pending, seen := map[string]bool{}, map[string]bool{}
	conversation := false
	for i, x := range ms {
		p := fmt.Sprintf("$.messages[%d]", i)
		m, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if _, ok := m["name"]; ok {
			return nil, at(p+".name", ErrUnsupported)
		}
		if err := rejectPresent(m, p, "function_call"); err != nil {
			return nil, err
		}
		if v, ok := m["audio"]; ok && v != nil {
			return nil, at(p+".audio", ErrUnsupported)
		}
		if v, ok := m["refusal"]; ok && v != nil && v != "" {
			return nil, at(p+".refusal", ErrUnsupported)
		}
		if err := rejectNonEmptyArray(m, "annotations", p); err != nil {
			return nil, err
		}
		role, e := requiredString(m, "role", p)
		if e != nil {
			return nil, e
		}
		if role == "system" || role == "developer" {
			if conversation {
				return nil, at(p+".role", fmt.Errorf("%w: system/developer messages must precede the conversation", ErrUnsupported))
			}
			c, e := chatCompletionsTextContent(m["content"], role, p+".content")
			if e != nil {
				return nil, e
			}
			for _, z := range c {
				q := z.(map[string]any)
				if q["type"] != "text" {
					return nil, at(p+".content", ErrUnsupported)
				}
				system = append(system, q)
			}
			continue
		}
		conversation = true
		if len(pending) != 0 && role != "tool" {
			return nil, at(p+".role", fmt.Errorf("%w: tool calls must be followed by role=tool messages", ErrInvalidSequence))
		}
		switch role {
		case "user":
			c, e := chatCompletionsTextContent(m["content"], role, p+".content")
			if e != nil {
				return nil, e
			}
			if len(c) == 0 {
				return nil, at(p+".content", ErrInvalidWireData)
			}
			appendMsg("user", c)
		case "assistant":
			content := []any{}
			if m["content"] != nil {
				content, e = chatCompletionsTextContent(m["content"], role, p+".content")
				if e != nil {
					return nil, e
				}
			}
			if v, ok := m["tool_calls"]; ok {
				toolUses, ids, e := chatToolCallsToAnthropic(v, p+".tool_calls", seen, fmt.Errorf("%w: duplicate tool call ID", ErrInvalidSequence))
				if e != nil {
					return nil, e
				}
				content = append(content, toolUses...)
				for _, id := range ids {
					pending[id] = true
				}
			}
			if len(content) == 0 {
				return nil, at(p+".content", fmt.Errorf("%w: assistant message needs text or tool_calls", ErrInvalidWireData))
			}
			appendMsg("assistant", content)
		case "tool":
			id, e := requiredString(m, "tool_call_id", p)
			if e != nil {
				return nil, e
			}
			if !pending[id] {
				return nil, at(p+".tool_call_id", fmt.Errorf("%w: orphan or duplicate tool message", ErrInvalidSequence))
			}
			text, e := chatCompletionsToolText(m["content"], p+".content")
			if e != nil {
				return nil, e
			}
			appendMsg("user", []any{map[string]any{"type": "tool_result", "tool_use_id": id, "content": text}})
			delete(pending, id)
		default:
			return nil, at(p+".role", ErrUnsupported)
		}
	}
	if len(messages) == 0 {
		return nil, at("$.messages", ErrInvalidWireData)
	}
	if messages[0].(map[string]any)["role"] != "user" {
		return nil, at("$.messages", fmt.Errorf("%w: Anthropic conversation must begin with user", ErrInvalidSequence))
	}
	if len(pending) != 0 {
		return nil, at("$.messages", fmt.Errorf("%w: unresolved tool calls", ErrInvalidSequence))
	}
	a["messages"] = messages
	if len(system) > 0 {
		a["system"] = system
	}
	if v, ok := r["tools"]; ok {
		a["tools"], err = chatCompletionsToolsToAnthropic(v, "$.tools")
		if err != nil {
			return nil, err
		}
	}
	if v, ok := r["tool_choice"]; ok {
		switch x := v.(type) {
		case string:
			switch x {
			case "auto":
				a["tool_choice"] = map[string]any{"type": "auto"}
			case "required":
				a["tool_choice"] = map[string]any{"type": "any"}
			case "none":
				a["tool_choice"] = map[string]any{"type": "none"}
			default:
				return nil, at("$.tool_choice", ErrUnsupported)
			}
		case map[string]any:
			if x["type"] != "function" {
				return nil, at("$.tool_choice.type", ErrUnsupported)
			}
			f, ok := obj(x["function"])
			if !ok {
				return nil, at("$.tool_choice.function", ErrInvalidWireData)
			}
			name, e := requiredString(f, "name", "$.tool_choice.function")
			if e != nil {
				return nil, e
			}
			a["tool_choice"] = map[string]any{"type": "tool", "name": name}
		default:
			return nil, at("$.tool_choice", ErrInvalidWireData)
		}
	}
	return encode(a)
}
