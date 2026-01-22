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
	"net/http"

	"go.opentelemetry.io/otel/trace"
)

// Option configures a traced HTTP client.
type Option func(*clientOptions)

type clientOptions struct {
	tracerProvider trace.TracerProvider
}

// WithTracerProvider returns an Option that sets the tracer provider.
// If not provided, the global tracer provider is used.
func WithTracerProvider(tp trace.TracerProvider) Option {
	return func(opts *clientOptions) {
		opts.tracerProvider = tp
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
	client.Transport = newRoundTripper(transport, options.tracerProvider)
	return client
}

type roundTripper struct {
	base           http.RoundTripper
	tracerProvider trace.TracerProvider
}

func newRoundTripper(base http.RoundTripper, tp trace.TracerProvider) http.RoundTripper {
	return &roundTripper{
		base:           base,
		tracerProvider: tp,
	}
}

// RoundTrip intercepts requests/responses to add tracing via the OpenAI middleware.
// This follows the middleware pattern:
//  1. roundTripper.RoundTrip() intercepts here
//  2. Calls Middleware() to add tracing
//  3. Middleware() calls base.RoundTrip() for actual HTTP request
func (rt *roundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	// Use the OpenAI middleware to add tracing
	// The middleware will call next (base.RoundTrip) to make the actual HTTP request
	next := func(r *http.Request) (*http.Response, error) {
		return rt.base.RoundTrip(r)
	}

	return MiddlewareWithTracerProvider(req, next, rt.tracerProvider)
}
