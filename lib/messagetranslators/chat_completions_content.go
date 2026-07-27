package messagetranslators

import (
	"fmt"
)

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

func chatCompletionsToolText(v any, path string) (string, error) {
	return flattenTextBlocks(v, path, false)
}
