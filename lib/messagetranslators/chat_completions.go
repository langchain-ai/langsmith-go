package messagetranslators

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

// chatCompletionsUsageToAnthropic converts Chat Completions usage. cached_tokens is a
// subset of prompt_tokens, as it is in the Responses API.
func chatCompletionsUsageToAnthropic(u map[string]any, path string) (map[string]any, error) {
	in, err := token(u["prompt_tokens"], path+".prompt_tokens")
	if err != nil {
		return nil, err
	}
	out, err := token(u["completion_tokens"], path+".completion_tokens")
	if err != nil {
		return nil, err
	}
	var cached int64
	if v, ok := u["prompt_tokens_details"]; ok {
		d, ok := obj(v)
		if !ok {
			return nil, at(path+".prompt_tokens_details", ErrInvalidWireData)
		}
		if v, ok := d["cached_tokens"]; ok {
			cached, err = token(v, path+".prompt_tokens_details.cached_tokens")
			if err != nil {
				return nil, err
			}
		}
		if v, ok := d["audio_tokens"]; ok {
			n, e := token(v, path+".prompt_tokens_details.audio_tokens")
			if e != nil {
				return nil, e
			}
			if n != 0 {
				return nil, at(path+".prompt_tokens_details.audio_tokens", ErrUnsupported)
			}
		}
	}
	if cached > in {
		return nil, at(path+".prompt_tokens_details.cached_tokens", fmt.Errorf("%w: cached_tokens exceeds prompt_tokens", ErrInvalidWireData))
	}
	if v, ok := u["completion_tokens_details"]; ok {
		d, ok := obj(v)
		if !ok {
			return nil, at(path+".completion_tokens_details", ErrInvalidWireData)
		}
		// Keep parity with Responses usage: Anthropic has no reasoning-token
		// category, so validate reasoning_tokens and intentionally drop it. Audio
		// and speculative-prediction counts describe unsupported output modes and
		// may only be present as zero-valued SDK defaults.
		if v, ok := d["reasoning_tokens"]; ok {
			if _, e := token(v, path+".completion_tokens_details.reasoning_tokens"); e != nil {
				return nil, e
			}
		}
		for _, key := range []string{"audio_tokens", "accepted_prediction_tokens", "rejected_prediction_tokens"} {
			if v, ok := d[key]; ok {
				n, e := token(v, path+".completion_tokens_details."+key)
				if e != nil {
					return nil, e
				}
				if n != 0 {
					return nil, at(path+".completion_tokens_details."+key, ErrUnsupported)
				}
			}
		}
	}
	if v, ok := u["total_tokens"]; ok {
		total, e := token(v, path+".total_tokens")
		if e != nil {
			return nil, e
		}
		if in > int64(^uint64(0)>>1)-out || total != in+out {
			return nil, at(path+".total_tokens", fmt.Errorf("%w: total_tokens must equal prompt_tokens + completion_tokens", ErrInvalidWireData))
		}
	}
	return map[string]any{"input_tokens": in - cached, "output_tokens": out, "cache_read_input_tokens": cached}, nil
}

func anthropicUsageToChatCompletions(u map[string]any, path string) (map[string]any, error) {
	ru, err := anthropicUsageToOpenAI(u, path)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"prompt_tokens":     ru["input_tokens"],
		"completion_tokens": ru["output_tokens"],
		"total_tokens":      ru["total_tokens"],
		"prompt_tokens_details": map[string]any{
			"cached_tokens": ru["input_tokens_details"].(map[string]any)["cached_tokens"],
		},
	}, nil
}

func anthropicImageToChatCompletions(o map[string]any, path string) (map[string]any, error) {
	u, err := anthropicImageURL(o, path, ErrInvalidWireData)
	if err != nil {
		return nil, err
	}
	return map[string]any{"type": "image_url", "image_url": map[string]any{"url": u}}, nil
}

func chatCompletionsImageToAnthropic(o map[string]any, path string) (map[string]any, error) {
	i, ok := obj(o["image_url"])
	if !ok {
		return nil, at(path+".image_url", ErrInvalidWireData)
	}
	if _, ok := i["detail"]; ok {
		return nil, at(path+".image_url.detail", ErrUnsupported)
	}
	u, ok := str(i["url"])
	if !ok || u == "" {
		return nil, at(path+".image_url.url", ErrInvalidWireData)
	}
	return imageURLToAnthropic(u, path+".image_url.url")
}

func anthropicToolsToChatCompletions(v any, path string) ([]any, error) {
	ts, ok := arr(v)
	if !ok {
		return nil, at(path, ErrInvalidWireData)
	}
	out := make([]any, 0, len(ts))
	for i, x := range ts {
		p := fmt.Sprintf("%s[%d]", path, i)
		o, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if err := rejectPresent(o, p, "type", "strict"); err != nil {
			return nil, err
		}
		name, description, schema, err := normalizeToolDefinition(o, p, "input_schema")
		if err != nil {
			return nil, err
		}
		f := map[string]any{"name": name, "parameters": schema}
		if description != nil {
			f["description"] = description
		}
		out = append(out, map[string]any{"type": "function", "function": f})
	}
	return out, nil
}

func chatCompletionsToolsToAnthropic(v any, path string) ([]any, error) {
	ts, ok := arr(v)
	if !ok {
		return nil, at(path, ErrInvalidWireData)
	}
	out := make([]any, 0, len(ts))
	for i, x := range ts {
		p := fmt.Sprintf("%s[%d]", path, i)
		o, ok := obj(x)
		if !ok || o["type"] != "function" {
			return nil, at(p+".type", ErrUnsupported)
		}
		if err := rejectPresent(o, p, "strict"); err != nil {
			return nil, err
		}
		f, ok := obj(o["function"])
		if !ok {
			return nil, at(p+".function", ErrInvalidWireData)
		}
		if err := rejectPresent(f, p+".function", "strict"); err != nil {
			return nil, err
		}
		name, description, schema, err := normalizeToolDefinition(f, p+".function", "parameters")
		if err != nil {
			return nil, err
		}
		t := map[string]any{"name": name, "input_schema": schema}
		if description != nil {
			t["description"] = description
		}
		out = append(out, t)
	}
	return out, nil
}

// AnthropicRequestToChatCompletions converts an Anthropic Messages request to
// an OpenAI Chat Completions request. model overrides the wire model when non-empty.
func AnthropicRequestToChatCompletions(body []byte, model string) ([]byte, error) {
	return AnthropicRequestToChatCompletionsWithOptions(body, model, ConversionOptions{})
}

// AnthropicRequestToChatCompletionsWithOptions is AnthropicRequestToChatCompletions with per-conversion warning options.
func AnthropicRequestToChatCompletionsWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	if options.WarningHandler != nil {
		inspectAnthropicObject(a, false, options, "$")
	}
	return anthropicRequestToChatCompletions(a, model)
}

func anthropicRequestToChatCompletions(a map[string]any, model string) ([]byte, error) {
	var err error
	if err := rejectPresent(a, "$", "thinking", "top_k", "service_tier", "output_config", "metadata", "mcp_servers", "context_management"); err != nil {
		return nil, err
	}
	if v, ok := a["stop_sequences"]; ok {
		ss, ok := arr(v)
		if !ok || len(ss) > 4 {
			return nil, at("$.stop_sequences", ErrInvalidWireData)
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
	if _, ok := a["model"]; ok {
		if _, err := requiredString(a, "model", "$"); err != nil {
			return nil, err
		}
	}
	dstModel := model
	if dstModel == "" {
		dstModel, err = requiredString(a, "model", "$")
		if err != nil {
			return nil, err
		}
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

func anthropicToolResultText(v any, path string) (string, error) {
	return flattenTextBlocks(v, path, true, "citations")
}

func chatCompletionsTextContent(v any, role, path string) ([]any, error) {
	if s, ok := str(v); ok {
		return []any{map[string]any{"type": "text", "text": s}}, nil
	}
	a, ok := arr(v)
	if !ok {
		return nil, at(path, ErrInvalidWireData)
	}
	out := []any{}
	for i, x := range a {
		p := fmt.Sprintf("%s[%d]", path, i)
		o, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		switch o["type"] {
		case "text":
			t, ok := str(o["text"])
			if !ok {
				return nil, at(p+".text", ErrInvalidWireData)
			}
			out = append(out, map[string]any{"type": "text", "text": t})
		case "image_url":
			if role != "user" {
				return nil, at(p+".type", ErrUnsupported)
			}
			im, err := chatCompletionsImageToAnthropic(o, p)
			if err != nil {
				return nil, err
			}
			out = append(out, im)
		case "input_audio", "file", "refusal":
			return nil, at(p+".type", ErrUnsupported)
		default:
			return nil, at(p+".type", ErrUnsupported)
		}
	}
	return out, nil
}

// ChatCompletionsRequestToAnthropic converts an OpenAI Chat Completions request
// to an Anthropic Messages request. Legacy functions/function_call are rejected.
func ChatCompletionsRequestToAnthropic(body []byte, model string) ([]byte, error) {
	return ChatCompletionsRequestToAnthropicWithOptions(body, model, ConversionOptions{})
}

// ChatCompletionsRequestToAnthropicWithOptions is ChatCompletionsRequestToAnthropic with per-conversion warning options.
func ChatCompletionsRequestToAnthropicWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	if options.WarningHandler != nil {
		inspectChatCompletionsObject(r, false, options, "$")
	}
	return chatCompletionsRequestToAnthropic(r, model)
}

func chatCompletionsRequestToAnthropic(r map[string]any, model string) ([]byte, error) {
	var err error
	if err := rejectPresent(r, "$", "functions", "function_call", "response_format", "parallel_tool_calls", "logprobs", "top_logprobs", "prediction", "modalities", "audio", "service_tier", "store", "metadata", "seed", "user", "reasoning_effort", "web_search_options", "frequency_penalty", "presence_penalty", "logit_bias", "verbosity", "safety_identifier"); err != nil {
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
	if _, ok := r["model"]; ok {
		if _, e := requiredString(r, "model", "$"); e != nil {
			return nil, e
		}
	}
	dstModel := model
	if dstModel == "" {
		dstModel, err = requiredString(r, "model", "$")
		if err != nil {
			return nil, err
		}
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

func chatCompletionsToolText(v any, path string) (string, error) {
	return flattenTextBlocks(v, path, false)
}

// ChatCompletionsResponseToAnthropic converts one completed Chat Completions response.
func ChatCompletionsResponseToAnthropic(body []byte, model string) ([]byte, error) {
	return ChatCompletionsResponseToAnthropicWithOptions(body, model, ConversionOptions{})
}

// ChatCompletionsResponseToAnthropicWithOptions is ChatCompletionsResponseToAnthropic with per-conversion warning options.
func ChatCompletionsResponseToAnthropicWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	if options.WarningHandler != nil {
		inspectChatCompletionsObject(r, true, options, "$")
	}
	return chatCompletionsResponseToAnthropic(r, model)
}

func chatCompletionsResponseToAnthropic(r map[string]any, model string) ([]byte, error) {
	if r["object"] != "chat.completion" {
		return nil, at("$.object", ErrInvalidWireData)
	}
	if err := rejectPresent(r, "$", "service_tier"); err != nil {
		return nil, err
	}
	if _, err := integer(r["created"], "$.created", false); err != nil {
		return nil, err
	}
	id, err := requiredString(r, "id", "$")
	if err != nil {
		return nil, err
	}
	wireModel, err := requiredString(r, "model", "$")
	if err != nil {
		return nil, err
	}
	choices, ok := arr(r["choices"])
	if !ok || len(choices) != 1 {
		return nil, at("$.choices", fmt.Errorf("%w: exactly one choice required", map[bool]error{true: ErrUnsupported, false: ErrInvalidWireData}[ok]))
	}
	c, ok := obj(choices[0])
	if !ok {
		return nil, at("$.choices[0]", ErrInvalidWireData)
	}
	idx, e := integer(c["index"], "$.choices[0].index", false)
	if e != nil || idx != 0 {
		if e != nil {
			return nil, e
		}
		return nil, at("$.choices[0].index", ErrInvalidWireData)
	}
	if v, ok := c["logprobs"]; ok && v != nil {
		return nil, at("$.choices[0].logprobs", ErrUnsupported)
	}
	m, ok := obj(c["message"])
	if !ok {
		return nil, at("$.choices[0].message", ErrInvalidWireData)
	}
	if m["role"] != "assistant" {
		return nil, at("$.choices[0].message.role", ErrInvalidWireData)
	}
	if err := rejectPresent(m, "$.choices[0].message", "function_call"); err != nil {
		return nil, err
	}
	if v, ok := m["audio"]; ok && v != nil {
		return nil, at("$.choices[0].message.audio", ErrUnsupported)
	}
	if err := rejectNonEmptyArray(m, "annotations", "$.choices[0].message"); err != nil {
		return nil, err
	}
	if v, ok := m["refusal"]; ok && v != nil && v != "" {
		return nil, at("$.choices[0].message.refusal", ErrUnsupported)
	}
	content := []any{}
	contentPresent := false
	if v, ok := m["content"]; ok && v != nil {
		contentPresent = true
		parts, e := chatCompletionsTextContent(v, "assistant", "$.choices[0].message.content")
		if e != nil {
			return nil, e
		}
		content = append(content, parts...)
	}
	if v, ok := m["tool_calls"]; ok {
		toolUses, _, e := chatToolCallsToAnthropic(v, "$.choices[0].message.tool_calls", map[string]bool{}, ErrInvalidSequence)
		if e != nil {
			return nil, e
		}
		content = append(content, toolUses...)
	}
	if !contentPresent && len(content) == 0 {
		return nil, at("$.choices[0].message.content", fmt.Errorf("%w: content is required when no tool calls are present", ErrInvalidWireData))
	}
	finish, ok := str(c["finish_reason"])
	if !ok {
		return nil, at("$.choices[0].finish_reason", ErrInvalidWireData)
	}
	hasTool := false
	for _, x := range content {
		if x.(map[string]any)["type"] == "tool_use" {
			hasTool = true
		}
	}
	stop, err := chatCompletionsFinishReasonToAnthropic(finish, hasTool, "$.choices[0].finish_reason")
	if err != nil {
		return nil, err
	}
	u, ok := obj(r["usage"])
	if !ok {
		return nil, at("$.usage", fmt.Errorf("%w: completed response requires usage", ErrInvalidWireData))
	}
	au, e := chatCompletionsUsageToAnthropic(u, "$.usage")
	if e != nil {
		return nil, e
	}
	dst := model
	if dst == "" {
		dst = wireModel
	}
	return encode(map[string]any{"id": destinationID("msg", id, 0), "type": "message", "role": "assistant", "model": dst, "content": content, "stop_reason": stop, "stop_sequence": nil, "usage": au})
}

// AnthropicResponseToChatCompletions converts one completed Anthropic Messages response.
func AnthropicResponseToChatCompletions(body []byte, model string) ([]byte, error) {
	return AnthropicResponseToChatCompletionsWithOptions(body, model, ConversionOptions{})
}

// AnthropicResponseToChatCompletionsWithOptions is AnthropicResponseToChatCompletions with per-conversion warning options.
func AnthropicResponseToChatCompletionsWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	if options.WarningHandler != nil {
		inspectAnthropicObject(a, true, options, "$")
	}
	return anthropicResponseToChatCompletions(a, model)
}

func anthropicResponseToChatCompletions(a map[string]any, model string) ([]byte, error) {
	if a["type"] != "message" {
		return nil, at("$.type", ErrInvalidWireData)
	}
	if a["role"] != "assistant" {
		return nil, at("$.role", ErrInvalidWireData)
	}
	if err := rejectPresent(a, "$", "phase"); err != nil {
		return nil, err
	}
	id, err := requiredString(a, "id", "$")
	if err != nil {
		return nil, err
	}
	wireModel, err := requiredString(a, "model", "$")
	if err != nil {
		return nil, err
	}
	dst := model
	if dst == "" {
		dst = wireModel
	}
	cs, ok := arr(a["content"])
	if !ok {
		return nil, at("$.content", ErrInvalidWireData)
	}
	var text strings.Builder
	calls := []any{}
	callIDs := map[string]bool{}
	sawTool := false
	for i, x := range cs {
		p := fmt.Sprintf("$.content[%d]", i)
		o, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if err := rejectPresent(o, p, "phase"); err != nil {
			return nil, err
		}
		switch o["type"] {
		case "text":
			if sawTool {
				return nil, at(p, fmt.Errorf("%w: Chat Completions orders text before tool calls", ErrUnsupported))
			}
			if err := rejectNonEmptyArray(o, "citations", p); err != nil {
				return nil, err
			}
			t, ok := str(o["text"])
			if !ok {
				return nil, at(p+".text", ErrInvalidWireData)
			}
			text.WriteString(t)
		case "tool_use":
			sawTool = true
			cid, e := requiredString(o, "id", p)
			if e != nil {
				return nil, e
			}
			if callIDs[cid] {
				return nil, at(p+".id", fmt.Errorf("%w: duplicate tool call ID", ErrInvalidSequence))
			}
			callIDs[cid] = true
			name, e := requiredString(o, "name", p)
			if e != nil {
				return nil, e
			}
			input := o["input"]
			if input == nil {
				input = map[string]any{}
			}
			if _, ok := obj(input); !ok {
				return nil, at(p+".input", ErrInvalidWireData)
			}
			args, _ := json.Marshal(input)
			calls = append(calls, map[string]any{"id": cid, "type": "function", "function": map[string]any{"name": name, "arguments": string(args)}})
		default:
			return nil, at(p+".type", ErrUnsupported)
		}
	}
	stop, ok := str(a["stop_reason"])
	if !ok {
		return nil, at("$.stop_reason", ErrInvalidWireData)
	}
	finish, err := anthropicStopReasonToChatCompletions(stop, len(calls) > 0, "$.stop_reason")
	if err != nil {
		return nil, err
	}
	m := map[string]any{"role": "assistant", "content": text.String(), "refusal": nil}
	if text.Len() == 0 && len(calls) > 0 {
		m["content"] = nil
	}
	if len(calls) > 0 {
		m["tool_calls"] = calls
	}
	u, ok := obj(a["usage"])
	if !ok {
		return nil, at("$.usage", ErrInvalidWireData)
	}
	cu, e := anthropicUsageToChatCompletions(u, "$.usage")
	if e != nil {
		return nil, e
	}
	now := time.Now().Unix()
	return encode(map[string]any{"id": destinationID("chatcmpl", id, 0), "object": "chat.completion", "created": now, "model": dst, "choices": []any{map[string]any{"index": int64(0), "message": m, "finish_reason": finish, "logprobs": nil}}, "usage": cu})
}
