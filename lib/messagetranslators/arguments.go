package messagetranslators

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

func parseArguments(v any, path string) (any, error) {
	s, ok := str(v)
	if !ok {
		return nil, at(path, ErrInvalidWireData)
	}
	return parseArgumentsRaw([]byte(s), path)
}

// parseArgumentsRaw validates accumulated tool-call arguments. Stream
// converters accumulate into a byte slice, so they avoid the string conversion
// parseArguments needs.
func parseArgumentsRaw(raw []byte, path string) (any, error) {
	dec := json.NewDecoder(bytes.NewReader(raw))
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
