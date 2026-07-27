package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
)

// DecodeObject decodes a single JSON object while preserving JSON numbers.
func DecodeObject(body []byte) (map[string]any, error) {
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

// Encode marshals v as JSON and classifies marshal failures as invalid wire data.
func Encode(v any) ([]byte, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return nil, fmt.Errorf("%w: %v", ErrInvalidWireData, err)
	}
	return b, nil
}

func String(v any) (string, bool)         { s, ok := v.(string); return s, ok }
func Array(v any) ([]any, bool)           { a, ok := v.([]any); return a, ok }
func Object(v any) (map[string]any, bool) { o, ok := v.(map[string]any); return o, ok }

func RequiredString(o map[string]any, key, path string) (string, error) {
	s, ok := String(o[key])
	if !ok || s == "" {
		return "", At(path+"."+key, fmt.Errorf("%w: non-empty string required", ErrInvalidWireData))
	}
	return s, nil
}

func CopyIf(dst, src map[string]any, to, from string) {
	if v, ok := src[from]; ok {
		dst[to] = v
	}
}

func Integer(v any, path string, positive bool) (int64, error) {
	n, ok := v.(json.Number)
	if !ok {
		return 0, At(path, fmt.Errorf("%w: integer required", ErrInvalidWireData))
	}
	i, err := n.Int64()
	if err != nil || i < 0 || (positive && i == 0) {
		kind := "nonnegative"
		if positive {
			kind = "positive"
		}
		return 0, At(path, fmt.Errorf("%w: %s integer required", ErrInvalidWireData, kind))
	}
	return i, nil
}

func NumberInRange(v any, path string, min, max float64) error {
	n, ok := v.(json.Number)
	if !ok {
		return At(path, fmt.Errorf("%w: number required", ErrInvalidWireData))
	}
	f, err := n.Float64()
	if err != nil || f < min || f > max {
		return At(path, fmt.Errorf("%w: number must be in [%g,%g]", ErrInvalidWireData, min, max))
	}
	return nil
}

func ValidateOptionalBool(o map[string]any, key, path string) error {
	if v, ok := o[key]; ok {
		if _, ok := v.(bool); !ok {
			return At(path+"."+key, fmt.Errorf("%w: boolean required", ErrInvalidWireData))
		}
	}
	return nil
}
