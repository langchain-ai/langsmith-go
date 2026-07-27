package messagetranslators

import (
	"fmt"
	"strings"
)

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
