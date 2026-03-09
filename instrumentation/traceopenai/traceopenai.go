// Package traceopenai provides OpenTelemetry tracing for the
// github.com/sashabaranov/go-openai client using LangSmith-compatible spans.
//
// Usage:
//
//	// Configure your sashabaranov client to use a traced HTTP client
//	cfg := openai.DefaultConfig(apiKey)
//	cfg.HTTPClient = traceopenai.Client()
//	client := openai.NewClientWithConfig(cfg)
//
//	// Or use a custom tracer provider:
//	tp := sdktrace.NewTracerProvider(...)
//	cfg.HTTPClient = traceopenai.Client(traceopenai.WithTracerProvider(tp))
//	client := openai.NewClientWithConfig(cfg)
//
//	// Your OpenAI API calls will now be automatically traced with LangSmith attrs
//	// resp, err := client.CreateChatCompletion(ctx, ...)
package traceopenai

import (
	"context"
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

type contextKey struct{ name string }

var ctxKeyRunName = contextKey{"run_name"}

// WithRunNameContext sets the span (run) name for the next traced request made with ctx.
// The run name in LangSmith is the OTLP span name; there is no separate field.
// Use this so one client can emit runs with different names per call, e.g. in tests:
//
//	ctx = traceopenai.WithRunNameContext(ctx, "openai_nonstreaming")
//	client.CreateChatCompletion(ctx, ...)
func WithRunNameContext(ctx context.Context, name string) context.Context {
	return context.WithValue(ctx, ctxKeyRunName, name)
}

// Option configures a traced HTTP client.
type Option func(*clientOptions)

type clientOptions struct {
	tracerProvider trace.TracerProvider
	runName        string
}

// WithTracerProvider returns an Option that sets the tracer provider.
// If not provided, the global tracer provider is used.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(opts *clientOptions) {
		opts.tracerProvider = tp
	}
}

// WithRunName sets the span (run) name to the given string when non-empty. Used by integration tests to identify runs in a shared project.
func WithRunName(name string) Option {
	return func(opts *clientOptions) {
		opts.runName = name
	}
}

// Client returns a new http.Client configured with tracing middleware.
// Equivalent to WrapClient(nil, opts...), which wraps the default transport.
func Client(opts ...Option) *http.Client {
	return WrapClient(nil, opts...)
}

// WrapClient wraps an existing http.Client with tracing middleware.
// If client is nil, a new client with the default transport is created.
func WrapClient(client *http.Client, opts ...Option) *http.Client {
	options := &clientOptions{}
	for _, opt := range opts {
		opt(options)
	}

	if client == nil {
		client = &http.Client{}
	}
	transport := client.Transport
	if transport == nil {
		transport = http.DefaultTransport
	}
	client.Transport = newRoundTripper(transport, options.tracerProvider, options.runName)
	return client
}

type roundTripper struct {
	base           http.RoundTripper
	tracerProvider trace.TracerProvider
	runName        string
}

func newRoundTripper(base http.RoundTripper, tp trace.TracerProvider, runName string) http.RoundTripper {
	return &roundTripper{
		base:           base,
		tracerProvider: tp,
		runName:        runName,
	}
}

// RoundTrip intercepts requests/responses to add tracing via the OpenAI middleware.
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	ctx := req.Context()
	// Prefer run name from context (per-request); fall back to client option.
	if v := ctx.Value(ctxKeyRunName); v == nil && rt.runName != "" {
		ctx = context.WithValue(ctx, ctxKeyRunName, rt.runName)
		req = req.WithContext(ctx)
	}
	next := func(r *http.Request) (*http.Response, error) {
		return rt.base.RoundTrip(r)
	}
	return MiddlewareWithTracerProvider(req, next, rt.tracerProvider)
}
