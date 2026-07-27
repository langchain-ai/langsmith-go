package messagetranslators

import (
	"bytes"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"strings"
	"time"
)

func decodeObject(body []byte) (map[string]any, error) {
	dec := json.NewDecoder(bytes.NewReader(body))
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWireData, err)
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return nil, fmt.Errorf("%w: trailing or malformed JSON", ErrInvalidWireData)
	}
	o, ok := v.(map[string]any)
	if !ok {
		return nil, fmt.Errorf("%w: top level must be an object", ErrInvalidWireData)
	}
	return o, nil
}

func encode(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWireData, err)
	}
	return b, nil
}

func str(v any) (string, bool)         { s, ok := v.(string); return s, ok }
func arr(v any) ([]any, bool)          { a, ok := v.([]any); return a, ok }
func obj(v any) (map[string]any, bool) { o, ok := v.(map[string]any); return o, ok }

func requiredString(o map[string]any, key, path string) (string, error) {
	s, ok := str(o[key])
	if !ok || s == "" {
		return "", at(path+"."+key, fmt.Errorf("%w: non-empty string required", ErrInvalidWireData))
	}
	return s, nil
}

func copyIf(dst, src map[string]any, to, from string) {
	if v, ok := src[from]; ok {
		dst[to] = v
	}
}

func integer(v any, path string, positive bool) (int64, error) {
	n, ok := v.(json.Number)
	if !ok {
		return 0, at(path, fmt.Errorf("%w: integer required", ErrInvalidWireData))
	}
	i, err := n.Int64()
	if err != nil || i < 0 || (positive && i == 0) {
		return 0, at(path, fmt.Errorf("%w: %s integer required", ErrInvalidWireData, map[bool]string{true: "positive", false: "nonnegative"}[positive]))
	}
	return i, nil
}

func numberInRange(v any, path string, min, max float64) error {
	n, ok := v.(json.Number)
	if !ok {
		return at(path, fmt.Errorf("%w: number required", ErrInvalidWireData))
	}
	f, err := n.Float64()
	if err != nil || f < min || f > max {
		return at(path, fmt.Errorf("%w: number must be in [%g,%g]", ErrInvalidWireData, min, max))
	}
	return nil
}

func validateOptionalBool(o map[string]any, key, path string) error {
	if v, ok := o[key]; ok {
		if _, ok := v.(bool); !ok {
			return at(path+"."+key, fmt.Errorf("%w: boolean required", ErrInvalidWireData))
		}
	}
	return nil
}

func token(v any, path string) (int64, error) { return integer(v, path, false) }

// rejectNonEmptyArray implements the v0 policy for metadata such as
// citations and annotations: an explicitly empty SDK-default array is harmless,
// but populated semantic metadata must never be silently discarded.
func rejectNonEmptyArray(o map[string]any, key, path string) error {
	v, exists := o[key]
	if !exists {
		return nil
	}
	a, ok := arr(v)
	if !ok {
		return at(path+"."+key, ErrInvalidWireData)
	}
	if len(a) != 0 {
		return at(path+"."+key, ErrUnsupported)
	}
	return nil
}

func rejectPhase(o map[string]any, path string) error {
	if _, ok := o["phase"]; ok {
		return at(path+".phase", ErrUnsupported)
	}
	return nil
}

// openAIUsageToAnthropic accounts for cached_tokens being a subset of input_tokens.
func openAIUsageToAnthropic(u map[string]any, path string) (map[string]any, error) {
	in, err := token(u["input_tokens"], path+".input_tokens")
	if err != nil {
		return nil, err
	}
	out, err := token(u["output_tokens"], path+".output_tokens")
	if err != nil {
		return nil, err
	}
	var cached int64
	if details, exists := u["input_tokens_details"]; exists {
		d, ok := obj(details)
		if !ok {
			return nil, at(path+".input_tokens_details", ErrInvalidWireData)
		}
		if v, exists := d["cached_tokens"]; exists {
			cached, err = token(v, path+".input_tokens_details.cached_tokens")
			if err != nil {
				return nil, err
			}
		}
	}
	if cached > in {
		return nil, at(path+".input_tokens_details.cached_tokens", fmt.Errorf("%w: cached_tokens exceeds input_tokens", ErrInvalidWireData))
	}
	// Anthropic has no reasoning-token usage category. Validate the wire value
	// before intentionally dropping it (Python/JS parity), rather than accepting
	// malformed details merely because the destination cannot represent them.
	if details, exists := u["output_tokens_details"]; exists {
		d, ok := obj(details)
		if !ok {
			return nil, at(path+".output_tokens_details", ErrInvalidWireData)
		}
		if v, exists := d["reasoning_tokens"]; exists {
			if _, err := token(v, path+".output_tokens_details.reasoning_tokens"); err != nil {
				return nil, err
			}
		}
	}
	if v, ok := u["total_tokens"]; ok {
		total, err := token(v, path+".total_tokens")
		if err != nil {
			return nil, err
		}
		if in > int64(^uint64(0)>>1)-out || total != in+out {
			return nil, at(path+".total_tokens", fmt.Errorf("%w: total_tokens must equal input_tokens + output_tokens", ErrInvalidWireData))
		}
	}
	return map[string]any{"input_tokens": in - cached, "output_tokens": out, "cache_read_input_tokens": cached}, nil
}

// anthropicUsageToOpenAI sums Anthropic's disjoint input token categories.
func anthropicUsageToOpenAI(u map[string]any, path string) (map[string]any, error) {
	in, err := token(u["input_tokens"], path+".input_tokens")
	if err != nil {
		return nil, err
	}
	out, err := token(u["output_tokens"], path+".output_tokens")
	if err != nil {
		return nil, err
	}
	var read, create int64
	if v, ok := u["cache_read_input_tokens"]; ok {
		read, err = token(v, path+".cache_read_input_tokens")
		if err != nil {
			return nil, err
		}
	}
	if v, ok := u["cache_creation_input_tokens"]; ok {
		create, err = token(v, path+".cache_creation_input_tokens")
		if err != nil {
			return nil, err
		}
	}
	const maxInt64 = int64(^uint64(0) >> 1)
	if read > maxInt64-in || create > maxInt64-in-read || out > maxInt64-in-read-create {
		return nil, at(path, fmt.Errorf("%w: token total overflows int64", ErrInvalidWireData))
	}
	totalIn := in + read + create
	return map[string]any{"input_tokens": totalIn, "output_tokens": out, "total_tokens": totalIn + out, "input_tokens_details": map[string]any{"cached_tokens": read}, "output_tokens_details": map[string]any{"reasoning_tokens": int64(0)}}, nil
}

func destinationID(prefix, source string, index int) string {
	h := sha256.Sum256([]byte(fmt.Sprintf("%s\x00%s\x00%d", prefix, source, index)))
	return prefix + "_" + hex.EncodeToString(h[:12])
}

func anthropicSystem(v any) (string, error) {
	if v == nil {
		return "", nil
	}
	if s, ok := str(v); ok {
		return s, nil
	}
	a, ok := arr(v)
	if !ok {
		return "", at("$.system", fmt.Errorf("%w: expected string or array", ErrInvalidWireData))
	}
	var b strings.Builder
	for i, x := range a {
		o, ok := obj(x)
		if !ok {
			return "", at(fmt.Sprintf("$.system[%d]", i), ErrInvalidWireData)
		}
		if o["type"] != "text" {
			return "", at(fmt.Sprintf("$.system[%d].type", i), ErrUnsupported)
		}
		t, ok := str(o["text"])
		if !ok {
			return "", at(fmt.Sprintf("$.system[%d].text", i), ErrInvalidWireData)
		}
		b.WriteString(t)
	}
	return b.String(), nil
}

func anthropicImageToResponses(o map[string]any, path string) (map[string]any, error) {
	s, ok := obj(o["source"])
	if !ok {
		return nil, at(path+".source", ErrInvalidWireData)
	}
	t, _ := str(s["type"])
	var u string
	switch t {
	case "url":
		u, ok = str(s["url"])
		if !ok || u == "" {
			return nil, at(path+".source.url", ErrInvalidWireData)
		}
	case "base64":
		media, mok := str(s["media_type"])
		data, dok := str(s["data"])
		if !mok || !dok || media == "" || data == "" {
			return nil, at(path+".source", fmt.Errorf("%w: non-empty media_type and data required", ErrInvalidWireData))
		}
		if _, err := base64.StdEncoding.DecodeString(data); err != nil {
			return nil, at(path+".source.data", fmt.Errorf("%w: invalid base64", ErrInvalidWireData))
		}
		u = "data:" + media + ";base64," + data
	default:
		return nil, at(path+".source.type", ErrUnsupported)
	}
	return map[string]any{"type": "input_image", "image_url": u}, nil
}

func anthContent(v any, role, path string) ([]any, error) {
	if s, ok := str(v); ok {
		textType := "input_text"
		if role == "assistant" {
			textType = "output_text"
		}
		return []any{map[string]any{"type": textType, "text": s}}, nil
	}
	a, ok := arr(v)
	if !ok {
		return nil, at(path, fmt.Errorf("%w: expected string or array", ErrInvalidWireData))
	}
	out := make([]any, 0, len(a))
	for i, x := range a {
		p := fmt.Sprintf("%s[%d]", path, i)
		o, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if err := rejectPhase(o, p); err != nil {
			return nil, err
		}
		t, _ := str(o["type"])
		switch t {
		case "text":
			txt, ok := str(o["text"])
			if !ok {
				return nil, at(p+".text", ErrInvalidWireData)
			}
			if err := rejectNonEmptyArray(o, "citations", p); err != nil {
				return nil, err
			}
			textType := "input_text"
			if role == "assistant" {
				textType = "output_text"
			}
			out = append(out, map[string]any{"type": textType, "text": txt})
		case "image":
			if role != "user" {
				return nil, at(p, fmt.Errorf("%w: image is only valid in user messages", ErrInvalidWireData))
			}
			im, err := anthropicImageToResponses(o, p)
			if err != nil {
				return nil, err
			}
			out = append(out, im)
		case "tool_use":
			if role != "assistant" {
				return nil, at(p, fmt.Errorf("%w: tool_use is only valid in assistant messages", ErrInvalidWireData))
			}
			id, err := requiredString(o, "id", p)
			if err != nil {
				return nil, err
			}
			name, err := requiredString(o, "name", p)
			if err != nil {
				return nil, err
			}
			input, ok := o["input"]
			if !ok {
				input = map[string]any{}
			}
			if _, ok := input.(map[string]any); !ok {
				return nil, at(p+".input", fmt.Errorf("%w: tool input must be an object", ErrInvalidWireData))
			}
			args, err := json.Marshal(input)
			if err != nil {
				return nil, at(p+".input", ErrInvalidWireData)
			}
			out = append(out, map[string]any{"type": "function_call", "call_id": id, "name": name, "arguments": string(args)})
		case "tool_result":
			if role != "user" {
				return nil, at(p, fmt.Errorf("%w: tool_result is only valid in user messages", ErrInvalidWireData))
			}
			if v, exists := o["is_error"]; exists {
				isError, ok := v.(bool)
				if !ok {
					return nil, at(p+".is_error", ErrInvalidWireData)
				}
				if isError {
					return nil, at(p+".is_error", ErrUnsupported)
				}
			}
			id, err := requiredString(o, "tool_use_id", p)
			if err != nil {
				return nil, err
			}
			var output string
			switch c := o["content"].(type) {
			case string:
				output = c
			case nil:
				output = ""
			case []any:
				var b strings.Builder
				for j, z := range c {
					q, ok := obj(z)
					if !ok || q["type"] != "text" {
						return nil, at(fmt.Sprintf("%s.content[%d]", p, j), ErrUnsupported)
					}
					tx, ok := str(q["text"])
					if !ok {
						return nil, at(fmt.Sprintf("%s.content[%d].text", p, j), ErrInvalidWireData)
					}
					b.WriteString(tx)
				}
				output = b.String()
			default:
				return nil, at(p+".content", ErrInvalidWireData)
			}
			out = append(out, map[string]any{"type": "function_call_output", "call_id": id, "output": output})
		default:
			return nil, at(p+".type", fmt.Errorf("%w: Anthropic content block %q", ErrUnsupported, t))
		}
	}
	return out, nil
}

// AnthropicRequestToResponses converts an Anthropic Messages request JSON body
// to an OpenAI Responses request JSON body. model overrides the wire model when non-empty.
func AnthropicRequestToResponses(body []byte, model string) ([]byte, error) {
	return AnthropicRequestToResponsesWithOptions(body, model, ConversionOptions{})
}

// AnthropicRequestToResponsesWithOptions is AnthropicRequestToResponses with per-conversion warning options.
func AnthropicRequestToResponsesWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	if options.WarningHandler != nil {
		o, err := decodeObject(body)
		if err != nil {
			return nil, err
		}
		inspectAnthropicObject(o, false, options, "$")
	}
	return anthropicRequestToResponses(body, model)
}

func anthropicRequestToResponses(body []byte, model string) ([]byte, error) {
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	for _, key := range []string{"thinking", "stop_sequences", "top_k", "service_tier", "output_config"} {
		if _, ok := a[key]; ok {
			return nil, at("$."+key, ErrUnsupported)
		}
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
	if _, exists := a["model"]; exists {
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
	r := map[string]any{"model": dstModel, "max_output_tokens": max}

	if s, err := anthropicSystem(a["system"]); err != nil {
		return nil, err
	} else if s != "" {
		r["instructions"] = s
	}
	copyIf(r, a, "temperature", "temperature")
	copyIf(r, a, "top_p", "top_p")
	copyIf(r, a, "stream", "stream")
	if ms, ok := arr(a["messages"]); ok {
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
	} else if a["messages"] != nil {
		return nil, at("$.messages", ErrInvalidWireData)
	}
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
			if v, ok := o["description"]; ok {
				if _, ok := str(v); !ok {
					return nil, at(fmt.Sprintf("$.tools[%d].description", i), ErrInvalidWireData)
				}
			}
			if v, ok := o["input_schema"]; ok {
				if _, ok := obj(v); !ok {
					return nil, at(fmt.Sprintf("$.tools[%d].input_schema", i), ErrInvalidWireData)
				}
			}
			if _, ok := o["strict"]; ok {
				return nil, at(fmt.Sprintf("$.tools[%d].strict", i), ErrUnsupported)
			}
			t := map[string]any{"type": "function", "name": name}
			copyIf(t, o, "description", "description")
			if p, ok := o["input_schema"]; ok {
				t["parameters"] = p
			} else {
				t["parameters"] = map[string]any{"type": "object"}
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

func responsesImageToAnthropic(o map[string]any, path string) (map[string]any, error) {
	u, ok := str(o["image_url"])
	if !ok || u == "" {
		return nil, at(path+".image_url", ErrInvalidWireData)
	}
	if strings.HasPrefix(u, "data:") {
		parts := strings.SplitN(strings.TrimPrefix(u, "data:"), ",", 2)
		if len(parts) != 2 || !strings.HasSuffix(parts[0], ";base64") {
			return nil, at(path+".image_url", ErrUnsupported)
		}
		media := strings.TrimSuffix(parts[0], ";base64")
		if media == "" || parts[1] == "" {
			return nil, at(path+".image_url", ErrInvalidWireData)
		}
		if _, e := base64.StdEncoding.DecodeString(parts[1]); e != nil {
			return nil, at(path+".image_url", ErrInvalidWireData)
		}
		return map[string]any{"type": "image", "source": map[string]any{"type": "base64", "media_type": media, "data": parts[1]}}, nil
	}
	return map[string]any{"type": "image", "source": map[string]any{"type": "url", "url": u}}, nil
}

func responseMessageContent(v any, role, path string) ([]any, error) {
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
		case "input_text", "output_text":
			t, ok := str(o["text"])
			if !ok {
				return nil, at(p+".text", ErrInvalidWireData)
			}
			if o["type"] == "output_text" {
				if err := rejectNonEmptyArray(o, "annotations", p); err != nil {
					return nil, err
				}
			}
			if err := rejectPhase(o, p); err != nil {
				return nil, err
			}
			out = append(out, map[string]any{"type": "text", "text": t})
		case "input_image":
			if role != "user" {
				return nil, at(p, fmt.Errorf("%w: image is only valid in user messages", ErrInvalidWireData))
			}
			if _, exists := o["detail"]; exists {
				return nil, at(p+".detail", ErrUnsupported)
			}
			if err := rejectPhase(o, p); err != nil {
				return nil, err
			}
			im, e := responsesImageToAnthropic(o, p)
			if e != nil {
				return nil, e
			}
			out = append(out, im)
		default:
			return nil, at(p+".type", ErrUnsupported)
		}
	}
	return out, nil
}

func parseArguments(v any, path string) (any, error) {
	s, ok := str(v)
	if !ok {
		return nil, at(path, ErrInvalidWireData)
	}
	dec := json.NewDecoder(strings.NewReader(s))
	dec.UseNumber()
	var x any
	if err := dec.Decode(&x); err != nil {
		return nil, at(path, fmt.Errorf("%w: malformed function arguments: %v", ErrInvalidWireData, err))
	}
	var trailing any
	if err := dec.Decode(&trailing); err != io.EOF {
		return nil, at(path, fmt.Errorf("%w: trailing or malformed function arguments JSON", ErrInvalidWireData))
	}
	if _, ok := x.(map[string]any); !ok {
		return nil, at(path, fmt.Errorf("%w: function arguments must be an object", ErrInvalidWireData))
	}
	return x, nil
}

// ResponsesRequestToAnthropic converts an OpenAI Responses request JSON body
// to an Anthropic Messages request JSON body. model overrides the wire model when non-empty.
func ResponsesRequestToAnthropic(body []byte, model string) ([]byte, error) {
	return ResponsesRequestToAnthropicWithOptions(body, model, ConversionOptions{})
}

// ResponsesRequestToAnthropicWithOptions is ResponsesRequestToAnthropic with per-conversion warning options.
func ResponsesRequestToAnthropicWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	if options.WarningHandler != nil {
		o, err := decodeObject(body)
		if err != nil {
			return nil, err
		}
		inspectResponsesObject(o, false, options, "$")
	}
	return responsesRequestToAnthropic(body, model)
}

func responsesRequestToAnthropic(body []byte, model string) ([]byte, error) {
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	for _, key := range []string{"previous_response_id", "reasoning", "truncation", "max_tool_calls", "parallel_tool_calls", "text", "phase", "prompt_cache_key", "prompt_cache_retention"} {
		if _, ok := r[key]; ok {
			return nil, at("$."+key, ErrUnsupported)
		}
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
	if _, exists := r["model"]; exists {
		if _, err := requiredString(r, "model", "$"); err != nil {
			return nil, err
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
			if v, ok := o["description"]; ok {
				if _, ok := str(v); !ok {
					return nil, at(p+".description", ErrInvalidWireData)
				}
			}
			if v, ok := o["parameters"]; ok {
				if _, ok := obj(v); !ok {
					return nil, at(p+".parameters", ErrInvalidWireData)
				}
			}
			name, err := requiredString(o, "name", p)
			if err != nil {
				return nil, err
			}
			t := map[string]any{"name": name}
			copyIf(t, o, "description", "description")
			if v, ok := o["parameters"]; ok {
				t["input_schema"] = v
			} else {
				t["input_schema"] = map[string]any{"type": "object"}
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

// ResponsesResponseToAnthropic converts a completed OpenAI Responses response
// JSON body to an Anthropic Messages response JSON body.
func ResponsesResponseToAnthropic(body []byte, model string) ([]byte, error) {
	return ResponsesResponseToAnthropicWithOptions(body, model, ConversionOptions{})
}

// ResponsesResponseToAnthropicWithOptions is ResponsesResponseToAnthropic with per-conversion warning options.
func ResponsesResponseToAnthropicWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	if options.WarningHandler != nil {
		o, err := decodeObject(body)
		if err != nil {
			return nil, err
		}
		inspectResponsesObject(o, true, options, "$")
	}
	return responsesResponseToAnthropic(body, model)
}

func responsesResponseToAnthropic(body []byte, model string) ([]byte, error) {
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	status, ok := str(r["status"])
	if !ok || (status != "completed" && status != "incomplete") {
		return nil, at("$.status", fmt.Errorf("%w: completed or incomplete status required", ErrInvalidWireData))
	}
	if v, exists := r["object"]; exists && v != "response" {
		return nil, at("$.object", ErrInvalidWireData)
	}
	if err := rejectPhase(r, "$"); err != nil {
		return nil, err
	}
	sourceID, err := requiredString(r, "id", "$")
	if err != nil {
		return nil, err
	}
	a := map[string]any{"type": "message", "role": "assistant", "id": destinationID("msg", sourceID, 0)}
	if _, exists := r["model"]; exists {
		if _, err := requiredString(r, "model", "$"); err != nil {
			return nil, err
		}
	}
	if model != "" {
		a["model"] = model
	} else {
		a["model"], err = requiredString(r, "model", "$")
		if err != nil {
			return nil, err
		}
	}
	if status == "incomplete" {
		d, ok := obj(r["incomplete_details"])
		if !ok {
			return nil, at("$.incomplete_details", fmt.Errorf("%w: incomplete response requires details", ErrInvalidWireData))
		}
		if reason, ok := str(d["reason"]); !ok || reason == "" {
			return nil, at("$.incomplete_details.reason", ErrInvalidWireData)
		}
	} else if details, exists := r["incomplete_details"]; exists && details != nil {
		if _, isObject := obj(details); !isObject {
			return nil, at("$.incomplete_details", ErrInvalidWireData)
		}
		return nil, at("$.incomplete_details", fmt.Errorf("%w: completed response cannot have incomplete details", ErrInvalidWireData))
	}
	content := []any{}
	outs, ok := arr(r["output"])
	if !ok {
		return nil, at("$.output", fmt.Errorf("%w: output array required", ErrInvalidWireData))
	}
	for i, x := range outs {
		p := fmt.Sprintf("$.output[%d]", i)
		o, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if err := rejectPhase(o, p); err != nil {
			return nil, err
		}
		if o["type"] == nil || o["type"] == "" {
			return nil, at(p+".type", fmt.Errorf("%w: output item type required", ErrInvalidWireData))
		}
		switch o["type"] {
		case "message":
			if _, err := requiredString(o, "id", p); err != nil {
				return nil, err
			}
			if o["role"] != "assistant" {
				return nil, at(p+".role", fmt.Errorf("%w: assistant role required", ErrInvalidWireData))
			}
			itemStatus, ok := str(o["status"])
			if !ok || (itemStatus != "completed" && itemStatus != "incomplete") {
				return nil, at(p+".status", ErrInvalidWireData)
			}
			cs, ok := arr(o["content"])
			if !ok {
				return nil, at(p+".content", ErrInvalidWireData)
			}
			for j, z := range cs {
				cp := fmt.Sprintf("%s.content[%d]", p, j)
				q, ok := obj(z)
				if !ok {
					return nil, at(cp, ErrInvalidWireData)
				}
				if q["type"] == nil || q["type"] == "" {
					return nil, at(cp+".type", fmt.Errorf("%w: content part type required", ErrInvalidWireData))
				}
				if q["type"] != "output_text" {
					return nil, at(cp+".type", ErrUnsupported)
				}
				if err := rejectNonEmptyArray(q, "annotations", cp); err != nil {
					return nil, err
				}
				if err := rejectPhase(q, cp); err != nil {
					return nil, err
				}
				t, ok := str(q["text"])
				if !ok {
					return nil, at(cp+".text", ErrInvalidWireData)
				}
				content = append(content, map[string]any{"type": "text", "text": t})
			}
		case "function_call":
			if _, err := requiredString(o, "id", p); err != nil {
				return nil, err
			}
			itemStatus, ok := str(o["status"])
			if !ok || (itemStatus != "completed" && itemStatus != "incomplete") {
				return nil, at(p+".status", ErrInvalidWireData)
			}
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
				if status == "incomplete" && errors.Is(err, ErrInvalidWireData) {
					return nil, at(p+".arguments", fmt.Errorf("%w: truncated function call cannot be represented", ErrUnsupported))
				}
				return nil, err
			}
			content = append(content, map[string]any{"type": "tool_use", "id": id, "name": name, "input": args})
		default:
			return nil, at(p+".type", ErrUnsupported)
		}
	}
	a["content"] = content
	hasTool := false
	for _, x := range content {
		if x.(map[string]any)["type"] == "tool_use" {
			hasTool = true
		}
	}
	switch status {
	case "incomplete":
		reason := "max_tokens"
		if d, ok := obj(r["incomplete_details"]); ok {
			if x, _ := str(d["reason"]); x != "" && x != "max_output_tokens" {
				return nil, at("$.incomplete_details.reason", ErrUnsupported)
			}
		}
		a["stop_reason"] = reason
	case "completed":
		if hasTool {
			a["stop_reason"] = "tool_use"
		} else {
			a["stop_reason"] = "end_turn"
		}
	}
	a["stop_sequence"] = nil
	if u, ok := obj(r["usage"]); ok {
		au, err := openAIUsageToAnthropic(u, "$.usage")
		if err != nil {
			return nil, err
		}
		a["usage"] = au
	} else {
		return nil, at("$.usage", fmt.Errorf("%w: completed Responses response requires usage", ErrInvalidWireData))
	}
	return encode(a)
}

// AnthropicResponseToResponses converts a completed Anthropic Messages response
// JSON body to an OpenAI Responses response JSON body.
func AnthropicResponseToResponses(body []byte, model string) ([]byte, error) {
	return AnthropicResponseToResponsesWithOptions(body, model, ConversionOptions{})
}

// AnthropicResponseToResponsesWithOptions is AnthropicResponseToResponses with per-conversion warning options.
func AnthropicResponseToResponsesWithOptions(body []byte, model string, options ConversionOptions) ([]byte, error) {
	if options.WarningHandler != nil {
		o, err := decodeObject(body)
		if err != nil {
			return nil, err
		}
		inspectAnthropicObject(o, true, options, "$")
	}
	return anthropicResponseToResponses(body, model)
}

func anthropicResponseToResponses(body []byte, model string) ([]byte, error) {
	a, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	if a["type"] != "message" {
		return nil, at("$.type", fmt.Errorf("%w: completed Anthropic response type must be message", ErrInvalidWireData))
	}
	if a["role"] != "assistant" {
		return nil, at("$.role", fmt.Errorf("%w: completed Anthropic response role must be assistant", ErrInvalidWireData))
	}
	if err := rejectPhase(a, "$"); err != nil {
		return nil, err
	}
	sourceID, _ := str(a["id"])
	if sourceID == "" {
		return nil, at("$.id", fmt.Errorf("%w: non-empty string required", ErrInvalidWireData))
	}
	createdAt := time.Now().Unix()
	r := map[string]any{
		"id":                     destinationID("resp", sourceID, 0),
		"object":                 "response",
		"created_at":             createdAt,
		"completed_at":           createdAt,
		"error":                  nil,
		"incomplete_details":     nil,
		"instructions":           nil,
		"max_output_tokens":      nil,
		"max_tool_calls":         nil,
		"metadata":               map[string]any{},
		"parallel_tool_calls":    true,
		"previous_response_id":   nil,
		"prompt_cache_key":       nil,
		"prompt_cache_retention": nil,
		"reasoning":              nil,
		"safety_identifier":      nil,
		"service_tier":           "default",
		"store":                  true,
		"temperature":            1.0,
		"text":                   map[string]any{"format": map[string]any{"type": "text"}},
		"tool_choice":            "auto",
		"tools":                  []any{},
		"top_logprobs":           0,
		"top_p":                  1.0,
		"truncation":             "disabled",
		"user":                   nil,
	}
	if _, exists := a["model"]; exists {
		if _, err := requiredString(a, "model", "$"); err != nil {
			return nil, err
		}
	}
	if model != "" {
		r["model"] = model
	} else {
		r["model"], err = requiredString(a, "model", "$")
		if err != nil {
			return nil, err
		}
	}
	outs := []any{}
	cs, ok := arr(a["content"])
	if !ok {
		return nil, at("$.content", fmt.Errorf("%w: content array required", ErrInvalidWireData))
	}
	var textParts []any
	textRun := 0
	flushText := func() {
		if len(textParts) == 0 {
			return
		}
		outs = append(outs, map[string]any{"type": "message", "id": destinationID("msg", sourceID, textRun), "status": "completed", "role": "assistant", "content": textParts})
		textParts = nil
	}
	for i, x := range cs {
		p := fmt.Sprintf("$.content[%d]", i)
		o, ok := obj(x)
		if !ok {
			return nil, at(p, ErrInvalidWireData)
		}
		if err := rejectPhase(o, p); err != nil {
			return nil, err
		}
		if o["type"] == nil || o["type"] == "" {
			return nil, at(p+".type", fmt.Errorf("%w: content block type required", ErrInvalidWireData))
		}
		switch o["type"] {
		case "text":
			t, ok := str(o["text"])
			if !ok {
				return nil, at(p+".text", ErrInvalidWireData)
			}
			if err := rejectNonEmptyArray(o, "citations", p); err != nil {
				return nil, err
			}
			if len(textParts) == 0 {
				textRun = i
			}
			textParts = append(textParts, map[string]any{"type": "output_text", "text": t, "annotations": []any{}, "logprobs": []any{}})
		case "tool_use":
			flushText()
			id, err := requiredString(o, "id", p)
			if err != nil {
				return nil, err
			}
			name, err := requiredString(o, "name", p)
			if err != nil {
				return nil, err
			}
			input := o["input"]
			if input == nil {
				input = map[string]any{}
			}
			if _, ok := input.(map[string]any); !ok {
				return nil, at(p+".input", fmt.Errorf("%w: tool input must be an object", ErrInvalidWireData))
			}
			args, err := json.Marshal(input)
			if err != nil {
				return nil, at(p+".input", ErrInvalidWireData)
			}
			outs = append(outs, map[string]any{"type": "function_call", "id": destinationID("fc", sourceID, i), "call_id": id, "name": name, "arguments": string(args), "status": "completed"})
		default:
			return nil, at(p+".type", ErrUnsupported)
		}
	}
	flushText()
	r["output"] = outs
	hasTool := false
	for _, out := range outs {
		if out.(map[string]any)["type"] == "function_call" {
			hasTool = true
		}
	}
	stop, ok := str(a["stop_reason"])
	if !ok {
		return nil, at("$.stop_reason", fmt.Errorf("%w: terminal stop_reason required", ErrInvalidWireData))
	}
	switch stop {
	case "max_tokens":
		r["status"] = "incomplete"
		r["incomplete_details"] = map[string]any{"reason": "max_output_tokens"}
	case "tool_use":
		if !hasTool {
			return nil, at("$.stop_reason", fmt.Errorf("%w: tool_use stop requires a tool_use block", ErrInvalidSequence))
		}
		r["status"] = "completed"
	case "end_turn", "stop_sequence":
		if hasTool {
			return nil, at("$.stop_reason", fmt.Errorf("%w: tool output requires tool_use stop", ErrInvalidSequence))
		}
		r["status"] = "completed"
	case "pause_turn", "refusal":
		return nil, at("$.stop_reason", ErrUnsupported)
	default:
		return nil, at("$.stop_reason", fmt.Errorf("%w: unknown stop reason %q", ErrInvalidWireData, stop))
	}
	if u, ok := obj(a["usage"]); ok {
		ru, err := anthropicUsageToOpenAI(u, "$.usage")
		if err != nil {
			return nil, err
		}
		r["usage"] = ru
	} else {
		return nil, at("$.usage", fmt.Errorf("%w: usage object required", ErrInvalidWireData))
	}
	return encode(r)
}
