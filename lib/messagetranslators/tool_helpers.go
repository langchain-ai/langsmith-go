package messagetranslators

import "fmt"

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
