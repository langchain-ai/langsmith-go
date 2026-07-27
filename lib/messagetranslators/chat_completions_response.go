package messagetranslators

import (
	"encoding/json"
	"fmt"
	"strings"
)

// ChatCompletionsResponseToAnthropic converts one completed Chat Completions response.
func ChatCompletionsResponseToAnthropic(body []byte, model string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectChatCompletionsObject(r, true, cfg, "$")
	return chatCompletionsResponseToAnthropic(r, model, cfg)
}

func chatCompletionsResponseToAnthropic(r map[string]any, model string, cfg config) ([]byte, error) {
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
	dst, err := resolveRequiredModel(r, model, "$")
	if err != nil {
		return nil, err
	}
	choices, ok := arr(r["choices"])
	if !ok {
		return nil, at("$.choices", fmt.Errorf("%w: exactly one choice required", ErrInvalidWireData))
	}
	if len(choices) != 1 {
		return nil, at("$.choices", fmt.Errorf("%w: exactly one choice required", ErrUnsupported))
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
	au, e := chatCompletionsUsageToAnthropic(u, "$.usage", cfg)
	if e != nil {
		return nil, e
	}
	cfg.reportUsage(anthropicUsage(au))
	return encode(map[string]any{"id": destinationID("msg", id, 0), "type": "message", "role": "assistant", "model": dst, "content": content, "stop_reason": stop, "stop_sequence": nil, "usage": au})
}

// AnthropicResponseToChatCompletions converts one completed Anthropic Messages response.
func AnthropicResponseToChatCompletions(body []byte, model string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectAnthropicObject(a, true, cfg, "$")
	return anthropicResponseToChatCompletions(a, model, cfg)
}

func anthropicResponseToChatCompletions(a map[string]any, model string, cfg config) ([]byte, error) {
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
	dst, err := resolveRequiredModel(a, model, "$")
	if err != nil {
		return nil, err
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
	cu, e := anthropicUsageToChatCompletions(u, "$.usage", cfg)
	if e != nil {
		return nil, e
	}
	cfg.reportUsage(anthropicUsage(map[string]any{
		"input_tokens":            cu["prompt_tokens"].(int64) - cu["prompt_tokens_details"].(map[string]any)["cached_tokens"].(int64),
		"output_tokens":           cu["completion_tokens"],
		"cache_read_input_tokens": cu["prompt_tokens_details"].(map[string]any)["cached_tokens"],
	}))
	now := cfg.now().Unix()
	return encode(map[string]any{"id": destinationID("chatcmpl", id, 0), "object": "chat.completion", "created": now, "model": dst, "choices": []any{map[string]any{"index": int64(0), "message": m, "finish_reason": finish, "logprobs": nil}}, "usage": cu})
}
