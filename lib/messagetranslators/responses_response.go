package messagetranslators

import (
	"encoding/json"
	"errors"
	"fmt"
)

// ResponsesResponseToAnthropic converts a completed OpenAI Responses response
// JSON body to an Anthropic Messages response JSON body.
func ResponsesResponseToAnthropic(body []byte, modelOverride string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	r, err := decodeObject(body)
	if err != nil {
		return nil, err
	}
	inspectResponsesObject(r, true, cfg, "$")
	return responsesResponseToAnthropic(r, modelOverride, cfg)
}

func responsesResponseToAnthropic(r map[string]any, modelOverride string, cfg config) ([]byte, error) {
	status, ok := str(r["status"])
	if !ok || (status != "completed" && status != "incomplete") {
		return nil, at("$.status", fmt.Errorf("%w: completed or incomplete status required", ErrInvalidWireData))
	}
	if v, exists := r["object"]; exists && v != "response" {
		return nil, at("$.object", ErrInvalidWireData)
	}
	if err := rejectPresent(r, "$", "phase"); err != nil {
		return nil, err
	}
	sourceID, err := requiredString(r, "id", "$")
	if err != nil {
		return nil, err
	}
	a := map[string]any{"type": "message", "role": "assistant", "id": destinationID("msg", sourceID, 0)}
	if a["model"], err = resolveModel(r, modelOverride, "$"); err != nil {
		return nil, err
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
		if err := rejectPresent(o, p, "phase"); err != nil {
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
				if err := rejectPresent(q, cp, "phase"); err != nil {
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
	u, ok := obj(r["usage"])
	if !ok {
		return nil, at("$.usage", fmt.Errorf("%w: completed Responses response requires usage", ErrInvalidWireData))
	}
	au, err := openAIUsageToAnthropic(u, "$.usage", cfg)
	if err != nil {
		return nil, err
	}
	a["usage"] = au
	cfg.reportUsage(anthropicUsage(au))
	return encode(a)
}

// AnthropicResponseToResponses converts a completed Anthropic Messages response
// JSON body to an OpenAI Responses response JSON body.
func AnthropicResponseToResponses(body []byte, modelOverride string, options ...Option) ([]byte, error) {
	cfg := newConfig(options)
	a, err := decodeAnthropicPayload(body, true, cfg)
	if err != nil {
		return nil, err
	}
	return anthropicResponseToResponses(a, modelOverride, cfg)
}

func anthropicResponseToResponses(a map[string]any, modelOverride string, cfg config) ([]byte, error) {
	var err error
	if a["type"] != "message" {
		return nil, at("$.type", fmt.Errorf("%w: completed Anthropic response type must be message", ErrInvalidWireData))
	}
	if a["role"] != "assistant" {
		return nil, at("$.role", fmt.Errorf("%w: completed Anthropic response role must be assistant", ErrInvalidWireData))
	}
	if err := rejectPresent(a, "$", "phase"); err != nil {
		return nil, err
	}
	sourceID, _ := str(a["id"])
	if sourceID == "" {
		return nil, at("$.id", fmt.Errorf("%w: non-empty string required", ErrInvalidWireData))
	}
	createdAt := cfg.now().Unix()
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
	if r["model"], err = resolveModel(a, modelOverride, "$"); err != nil {
		return nil, err
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
		if err := rejectPresent(o, p, "phase"); err != nil {
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
	stop, ok := str(a["stop_reason"])
	if !ok {
		return nil, at("$.stop_reason", fmt.Errorf("%w: terminal stop_reason required", ErrInvalidWireData))
	}
	switch stop {
	case "max_tokens":
		r["status"] = "incomplete"
		r["incomplete_details"] = map[string]any{"reason": "max_output_tokens"}
	case "tool_use", "end_turn", "stop_sequence":
		// Responses encodes no tool-specific terminal state, so a source whose
		// stop_reason disagrees with its own content changes nothing here. See
		// reason_mapping.go for the same policy where the distinction is carried.
		r["status"] = "completed"
	case "pause_turn", "refusal":
		return nil, at("$.stop_reason", ErrUnsupported)
	default:
		return nil, at("$.stop_reason", fmt.Errorf("%w: unknown stop reason %q", ErrInvalidWireData, stop))
	}
	u, ok := obj(a["usage"])
	if !ok {
		return nil, at("$.usage", fmt.Errorf("%w: usage object required", ErrInvalidWireData))
	}
	ru, err := anthropicUsageToOpenAI(u, "$.usage", cfg)
	if err != nil {
		return nil, err
	}
	r["usage"] = ru
	cfg.reportUsage(openAIUsage(ru))
	return encode(r)
}
