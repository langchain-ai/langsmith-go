package messagetranslators

import (
	"encoding/base64"
	"fmt"
	"strings"
)

// anthropicImageURL validates an Anthropic image source and returns the URL
// representation shared by the Responses and Chat Completions converters.
func anthropicImageURL(o map[string]any, path string, incompleteSource error) (string, error) {
	s, ok := obj(o["source"])
	if !ok {
		return "", at(path+".source", ErrInvalidWireData)
	}
	switch s["type"] {
	case "url":
		u, ok := str(s["url"])
		if !ok || u == "" {
			return "", at(path+".source.url", ErrInvalidWireData)
		}
		return u, nil
	case "base64":
		media, mediaOK := str(s["media_type"])
		data, dataOK := str(s["data"])
		if !mediaOK || !dataOK || media == "" || data == "" {
			return "", at(path+".source", incompleteSource)
		}
		if _, err := base64.StdEncoding.DecodeString(data); err != nil {
			return "", at(path+".source.data", fmt.Errorf("%w: invalid base64", ErrInvalidWireData))
		}
		return "data:" + media + ";base64," + data, nil
	default:
		return "", at(path+".source.type", ErrUnsupported)
	}
}

// imageURLToAnthropic turns either an ordinary URL or a base64 data URL into
// an Anthropic image block. urlPath is kept explicit so callers retain their
// API-specific error paths.
func imageURLToAnthropic(u, urlPath string) (map[string]any, error) {
	if strings.HasPrefix(u, "data:") {
		parts := strings.SplitN(strings.TrimPrefix(u, "data:"), ",", 2)
		if len(parts) != 2 || !strings.HasSuffix(parts[0], ";base64") {
			return nil, at(urlPath, ErrUnsupported)
		}
		media := strings.TrimSuffix(parts[0], ";base64")
		if media == "" || parts[1] == "" {
			return nil, at(urlPath, ErrInvalidWireData)
		}
		if _, err := base64.StdEncoding.DecodeString(parts[1]); err != nil {
			return nil, at(urlPath, ErrInvalidWireData)
		}
		return map[string]any{"type": "image", "source": map[string]any{"type": "base64", "media_type": media, "data": parts[1]}}, nil
	}
	return map[string]any{"type": "image", "source": map[string]any{"type": "url", "url": u}}, nil
}

// chatToolCallsToAnthropic parses the function calls shared by request and
// completed Chat Completions messages.
func chatToolCallsToAnthropic(v any, path string, seen map[string]bool, duplicateErr error) ([]any, []string, error) {
	calls, ok := arr(v)
	if !ok || len(calls) == 0 {
		return nil, nil, at(path, ErrInvalidWireData)
	}
	out := make([]any, 0, len(calls))
	ids := make([]string, 0, len(calls))
	for i, x := range calls {
		p := fmt.Sprintf("%s[%d]", path, i)
		tc, ok := obj(x)
		if !ok || tc["type"] != "function" {
			return nil, nil, at(p+".type", ErrUnsupported)
		}
		id, err := requiredString(tc, "id", p)
		if err != nil {
			return nil, nil, err
		}
		if seen[id] {
			return nil, nil, at(p+".id", duplicateErr)
		}
		f, ok := obj(tc["function"])
		if !ok {
			return nil, nil, at(p+".function", ErrInvalidWireData)
		}
		name, err := requiredString(f, "name", p+".function")
		if err != nil {
			return nil, nil, err
		}
		args, err := parseArguments(f["arguments"], p+".function.arguments")
		if err != nil {
			return nil, nil, err
		}
		seen[id] = true
		ids = append(ids, id)
		out = append(out, map[string]any{"type": "tool_use", "id": id, "name": name, "input": args})
	}
	return out, ids, nil
}

func normalizeToolDefinition(o map[string]any, path, schemaKey string) (name string, description any, schema any, err error) {
	name, err = requiredString(o, "name", path)
	if err != nil {
		return "", nil, nil, err
	}
	var ok bool
	if description, ok = o["description"]; ok {
		if _, ok := str(description); !ok {
			return "", nil, nil, at(path+".description", ErrInvalidWireData)
		}
	}
	schema, ok = o[schemaKey]
	if !ok {
		schema = map[string]any{"type": "object"}
	} else if _, ok := obj(schema); !ok {
		return "", nil, nil, at(path+"."+schemaKey, ErrInvalidWireData)
	}
	return name, description, schema, nil
}

// flattenTextBlocks handles the string-or-text-block-array forms used by tool
// results in all three APIs. Callers opt into nil and metadata handling where
// their wire format permits it.
func flattenTextBlocks(v any, path string, allowNil bool, rejectedKeys ...string) (string, error) {
	if v == nil && allowNil {
		return "", nil
	}
	if s, ok := str(v); ok {
		return s, nil
	}
	a, ok := arr(v)
	if !ok {
		return "", at(path, ErrInvalidWireData)
	}
	var b strings.Builder
	for i, x := range a {
		p := fmt.Sprintf("%s[%d]", path, i)
		o, ok := obj(x)
		if !ok || o["type"] != "text" {
			return "", at(p, ErrUnsupported)
		}
		if err := rejectPresent(o, p, rejectedKeys...); err != nil {
			return "", err
		}
		t, ok := str(o["text"])
		if !ok {
			return "", at(p+".text", ErrInvalidWireData)
		}
		b.WriteString(t)
	}
	return b.String(), nil
}
