// Package instrumentation holds settings shared by the langsmith-go provider
// instrumentations (traceanthropic, tracegemini, traceopenai).
package instrumentation

import (
	"context"

	"github.com/langchain-ai/langsmith-go/internal/traceutil"
)

// WithoutContent returns a context whose traced requests omit request and
// response content (gen_ai.prompt, gen_ai.completion). It applies to every
// langsmith-go instrumentation. Model, usage, timing, and status are still
// recorded.
func WithoutContent(ctx context.Context) context.Context {
	return traceutil.WithoutContent(ctx)
}
