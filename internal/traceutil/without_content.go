package traceutil

import "context"

type withoutContentKey struct{}

// WithoutContent returns a context whose traced requests omit request and
// response content (gen_ai.prompt, gen_ai.completion).
func WithoutContent(ctx context.Context) context.Context {
	return context.WithValue(ctx, withoutContentKey{}, true)
}

// IsWithoutContent reports whether requests made with ctx omit content.
func IsWithoutContent(ctx context.Context) bool {
	return ctx.Value(withoutContentKey{}) != nil
}
