package messagetranslators

import (
	"encoding/json"
	"fmt"
)

func anthropicImageToResponses(o map[string]any, path string) (map[string]any, error) {
	u, err := anthropicImageURL(o, path, fmt.Errorf("%w: non-empty media_type and data required", ErrInvalidWireData))
	if err != nil {
		return nil, err
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
		if err := rejectPresent(o, p, "phase"); err != nil {
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
			output, err := flattenTextBlocks(o["content"], p+".content", true)
			if err != nil {
				return nil, err
			}
			out = append(out, map[string]any{"type": "function_call_output", "call_id": id, "output": output})
		default:
			return nil, at(p+".type", fmt.Errorf("%w: Anthropic content block %q", ErrUnsupported, t))
		}
	}
	return out, nil
}

func responsesImageToAnthropic(o map[string]any, path string) (map[string]any, error) {
	u, ok := str(o["image_url"])
	if !ok || u == "" {
		return nil, at(path+".image_url", ErrInvalidWireData)
	}
	return imageURLToAnthropic(u, path+".image_url")
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
			if err := rejectPresent(o, p, "phase"); err != nil {
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
			if err := rejectPresent(o, p, "phase"); err != nil {
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
