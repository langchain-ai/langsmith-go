package langsmith

import (
	"context"

	"github.com/google/uuid"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

type sandboxIDContextKey struct{}

func traceSandbox(ctx context.Context, id string) context.Context {
	if id == "" {
		return ctx
	}
	trace.SpanFromContext(ctx).SetAttributes(attribute.String("langsmith.metadata.sandbox_id", id))
	return context.WithValue(ctx, sandboxIDContextKey{}, id)
}

func traceSandboxReference(ctx context.Context, name string) context.Context {
	if sandboxIDFromContext(ctx) != "" {
		return ctx
	}
	if id, err := uuid.Parse(name); err == nil {
		return traceSandbox(ctx, id.String())
	}
	return ctx
}

func sandboxIDFromContext(ctx context.Context) string {
	id, _ := ctx.Value(sandboxIDContextKey{}).(string)
	return id
}
