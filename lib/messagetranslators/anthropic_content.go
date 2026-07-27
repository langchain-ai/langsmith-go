package messagetranslators

import (
	"fmt"
	"strings"
)

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
