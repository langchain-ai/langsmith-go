package messagetranslators

import (
	"fmt"
)

// AnthropicRequestToResponses converts an Anthropic Messages request JSON body
// to an OpenAI Responses request JSON body. modelOverride replaces the source model when non-empty.
func AnthropicRequestToResponses(body []byte, modelOverride string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectAnthropicObject(a, false, cfg, "$")
	return anthropicRequestToResponses(a, modelOverride, cfg)
}

func anthropicRequestToResponses(a map[string]any, modelOverride string, cfg config) ([]byte, error) {
	if err := rejectUnsupportedFields(a, "$", anthropicRequestUnsupported, anthropicToResponsesUnsupported); err != nil {
		return nil, err
	}
	if err := validateOptionalBool(a, "stream", "$"); err != nil {
		return nil, err
	}
	if v, ok := a["temperature"]; ok {
		if err := numberInRange(v, "$.temperature", 0, 1); err != nil {
			return nil, err
		}
	}
	if v, ok := a["top_p"]; ok {
		if err := numberInRange(v, "$.top_p", 0, 1); err != nil {
			return nil, err
		}
	}
	max, err := integer(a["max_tokens"], "$.max_tokens", true)
	if err != nil {
		return nil, err
	}
	dstModel, err := resolveModel(a, modelOverride, "$")
	if err != nil {
		return nil, err
	}
	ms, ok := arr(a["messages"])
	if !ok || len(ms) == 0 {
		return nil, at("$.messages", fmt.Errorf("%w: non-empty array required", ErrInvalidWireData))
	}
	r := map[string]any{"model": dstModel, "max_output_tokens": max}

	if s, err := anthropicSystem(a["system"]); err != nil {
		return nil, err
	} else if s != "" {
		r["instructions"] = s
	}
	copyIf(r, a, "temperature", "temperature")
	copyIf(r, a, "top_p", "top_p")
	copyIf(r, a, "stream", "stream")
	input := []any{}
	pendingCalls := map[string]bool{}
	for i, x := range ms {
		o, ok := obj(x)
		if !ok {
			return nil, at(fmt.Sprintf("$.messages[%d]", i), ErrInvalidWireData)
		}
		role, err := requiredString(o, "role", fmt.Sprintf("$.messages[%d]", i))
		if err != nil {
			return nil, err
		}
		if role != "user" && role != "assistant" {
			return nil, at(fmt.Sprintf("$.messages[%d].role", i), ErrUnsupported)
		}
		c, err := anthContent(o["content"], role, fmt.Sprintf("$.messages[%d].content", i))
		if err != nil {
			return nil, err
		}
		if len(c) == 0 {
			return nil, at(fmt.Sprintf("$.messages[%d].content", i), fmt.Errorf("%w: non-empty content required", ErrInvalidWireData))
		}
		var normal []any
		for _, z := range c {
			zo := z.(map[string]any)
			if zo["type"] == "function_call" {
				id := zo["call_id"].(string)
				if pendingCalls[id] {
					return nil, at(fmt.Sprintf("$.messages[%d].content", i), fmt.Errorf("%w: duplicate unresolved tool ID", ErrInvalidSequence))
				}
				pendingCalls[id] = true
			} else if zo["type"] == "function_call_output" {
				id := zo["call_id"].(string)
				if !pendingCalls[id] {
					return nil, at(fmt.Sprintf("$.messages[%d].content", i), fmt.Errorf("%w: tool_result has no preceding tool_use", ErrInvalidSequence))
				}
				delete(pendingCalls, id)
			}
			if zo["type"] == "function_call" || zo["type"] == "function_call_output" {
				if len(normal) > 0 {
					input = append(input, map[string]any{"type": "message", "role": role, "content": normal})
					normal = nil
				}
				input = append(input, zo)
			} else {
				normal = append(normal, zo)
			}
		}
		if len(normal) > 0 {
			input = append(input, map[string]any{"type": "message", "role": role, "content": normal})
		}
	}
	first := ms[0].(map[string]any)
	if first["role"] != "user" {
		return nil, at("$.messages[0].role", fmt.Errorf("%w: conversation must begin with user", ErrInvalidSequence))
	}
	if len(pendingCalls) != 0 {
		return nil, at("$.messages", fmt.Errorf("%w: tool_use requires a following tool_result", ErrInvalidSequence))
	}
	r["input"] = input
	if ts, ok := arr(a["tools"]); ok {
		tools := []any{}
		for i, x := range ts {
			o, ok := obj(x)
			if !ok {
				return nil, at(fmt.Sprintf("$.tools[%d]", i), ErrInvalidWireData)
			}
			name, err := requiredString(o, "name", fmt.Sprintf("$.tools[%d]", i))
			if err != nil {
				return nil, err
			}
			if _, ok := o["type"]; ok {
				return nil, at(fmt.Sprintf("$.tools[%d].type", i), ErrUnsupported)
			}
			_, description, schema, err := normalizeToolDefinition(o, fmt.Sprintf("$.tools[%d]", i), "input_schema")
			if err != nil {
				return nil, err
			}
			if _, ok := o["strict"]; ok {
				return nil, at(fmt.Sprintf("$.tools[%d].strict", i), ErrUnsupported)
			}
			t := map[string]any{"type": "function", "name": name, "parameters": schema}
			if description != nil {
				t["description"] = description
			}
			tools = append(tools, t)
		}
		r["tools"] = tools
	} else if a["tools"] != nil {
		return nil, at("$.tools", ErrInvalidWireData)
	}
	if tc, ok := obj(a["tool_choice"]); ok {
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
			name, err := requiredString(tc, "name", "$.tool_choice")
			if err != nil {
				return nil, err
			}
			r["tool_choice"] = map[string]any{"type": "function", "name": name}
		default:
			return nil, at("$.tool_choice.type", ErrUnsupported)
		}
	} else if a["tool_choice"] != nil {
		return nil, at("$.tool_choice", ErrInvalidWireData)
	}
	return encode(r)
}

// ResponsesRequestToAnthropic converts an OpenAI Responses request JSON body
// to an Anthropic Messages request JSON body. modelOverride replaces the source model when non-empty.
func ResponsesRequestToAnthropic(body []byte, modelOverride string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectResponsesObject(r, false, cfg, "$")
	return responsesRequestToAnthropic(r, modelOverride, cfg)
}

func responsesRequestToAnthropic(r map[string]any, modelOverride string, cfg config) ([]byte, error) {
	if err := rejectUnsupportedFields(r, "$", responsesRequestUnsupported); err != nil {
		return nil, err
	}
	if err := validateOptionalBool(r, "stream", "$"); err != nil {
		return nil, err
	}
	if v, ok := r["temperature"]; ok {
		if err := numberInRange(v, "$.temperature", 0, 1); err != nil {
			return nil, err
		}
	}
	if v, ok := r["top_p"]; ok {
		if err := numberInRange(v, "$.top_p", 0, 1); err != nil {
			return nil, err
		}
	}
	max, err := integer(r["max_output_tokens"], "$.max_output_tokens", true)
	if err != nil {
		return nil, err
	}
	dstModel, err := resolveModel(r, modelOverride, "$")
	if err != nil {
		return nil, err
	}
	a := map[string]any{"model": dstModel, "max_tokens": max}
	copyIf(a, r, "temperature", "temperature")
	copyIf(a, r, "top_p", "top_p")
	copyIf(a, r, "stream", "stream")
	if v, ok := r["instructions"]; ok {
		if _, ok := str(v); !ok {
			return nil, at("$.instructions", ErrUnsupported)
		}
		a["system"] = v
	}
	messages := []any{}
	appendMsg := func(role string, content []any) {
		if len(messages) > 0 {
			last := messages[len(messages)-1].(map[string]any)
			if last["role"] == role {
				last["content"] = append(last["content"].([]any), content...)
				return
			}
		}
		messages = append(messages, map[string]any{"role": role, "content": content})
	}
	pendingCalls := map[string]bool{}
	switch input := r["input"].(type) {
	case string:
		appendMsg("user", []any{map[string]any{"type": "text", "text": input}})
	case []any:
		for i, x := range input {
			p := fmt.Sprintf("$.input[%d]", i)
			o, ok := obj(x)
			if !ok {
				return nil, at(p, ErrInvalidWireData)
			}
			typ, _ := str(o["type"])
			if typ == "" || typ == "message" {
				role, err := requiredString(o, "role", p)
				if err != nil {
					return nil, err
				}
				if role != "user" && role != "assistant" && role != "system" && role != "developer" {
					return nil, at(p+".role", ErrUnsupported)
				}
				if role == "system" || role == "developer" {
					return nil, at(p+".role", fmt.Errorf("%w: instructions must carry system messages", ErrUnsupported))
				}
				c, err := responseMessageContent(o["content"], role, p+".content")
				if err != nil {
					return nil, err
				}
				if len(c) == 0 {
					return nil, at(p+".content", fmt.Errorf("%w: non-empty content required", ErrInvalidWireData))
				}
				appendMsg(role, c)
			} else if typ == "function_call" {
				id, err := requiredString(o, "call_id", p)
				if err != nil {
					return nil, err
				}
				name, err := requiredString(o, "name", p)
				if err != nil {
					return nil, err
				}
				args, err := parseArguments(o["arguments"], p+".arguments")
				if err != nil {
					return nil, err
				}
				if pendingCalls[id] {
					return nil, at(p+".call_id", fmt.Errorf("%w: duplicate unresolved call_id", ErrInvalidSequence))
				}
				pendingCalls[id] = true
				appendMsg("assistant", []any{map[string]any{"type": "tool_use", "id": id, "name": name, "input": args}})
			} else if typ == "function_call_output" {
				id, err := requiredString(o, "call_id", p)
				if err != nil {
					return nil, err
				}
				if !pendingCalls[id] {
					return nil, at(p+".call_id", fmt.Errorf("%w: function output has no preceding call", ErrInvalidSequence))
				}
				delete(pendingCalls, id)
				output, ok := str(o["output"])
				if !ok {
					return nil, at(p+".output", ErrUnsupported)
				}
				appendMsg("user", []any{map[string]any{"type": "tool_result", "tool_use_id": id, "content": output}})
			} else {
				return nil, at(p+".type", ErrUnsupported)
			}
		}
	case nil:
	default:
		return nil, at("$.input", ErrInvalidWireData)
	}
	if len(messages) == 0 {
		return nil, at("$.input", fmt.Errorf("%w: input must produce non-empty messages", ErrInvalidWireData))
	}
	if messages[0].(map[string]any)["role"] != "user" {
		return nil, at("$.input", fmt.Errorf("%w: Anthropic conversation must begin with a user message", ErrInvalidSequence))
	}
	if len(pendingCalls) != 0 {
		return nil, at("$.input", fmt.Errorf("%w: function calls require following outputs", ErrInvalidSequence))
	}
	a["messages"] = messages
	if ts, ok := arr(r["tools"]); ok {
		tools := []any{}
		for i, x := range ts {
			p := fmt.Sprintf("$.tools[%d]", i)
			o, ok := obj(x)
			if !ok || o["type"] != "function" {
				return nil, at(p, ErrUnsupported)
			}
			if _, ok := o["strict"]; ok {
				return nil, at(p+".strict", ErrUnsupported)
			}
			name, description, schema, err := normalizeToolDefinition(o, p, "parameters")
			if err != nil {
				return nil, err
			}
			t := map[string]any{"name": name, "input_schema": schema}
			if description != nil {
				t["description"] = description
			}
			tools = append(tools, t)
		}
		a["tools"] = tools
	} else if r["tools"] != nil {
		return nil, at("$.tools", ErrInvalidWireData)
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
			name, err := requiredString(x, "name", "$.tool_choice")
			if err != nil {
				return nil, err
			}
			a["tool_choice"] = map[string]any{"type": "tool", "name": name}
		default:
			return nil, at("$.tool_choice", ErrInvalidWireData)
		}
	}
	return encode(a)
}
