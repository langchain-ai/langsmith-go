package langsmith

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
	"go.opentelemetry.io/otel/trace"
)

const (
	defaultEndpoint     = "api.smith.langchain.com"
	defaultURLPath      = "/otel/v1/traces"
	defaultBatchTimeout = 1 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

// TracerOption configures a Tracer.
type TracerOption func(*tracerConfig)

type tracerConfig struct {
	apiKey       string
	projectName  string
	serviceName  string
	endpoint     string
	batchTimeout time.Duration
}

// WithAPIKey sets the LangSmith API key.
func WithAPIKey(apiKey string) TracerOption {
	return func(c *tracerConfig) {
		c.apiKey = apiKey
	}
}

// WithProjectName sets the LangSmith project name.
func WithProjectName(projectName string) TracerOption {
	return func(c *tracerConfig) {
		c.projectName = projectName
	}
}

// WithServiceName sets the service name for the tracer.
func WithServiceName(serviceName string) TracerOption {
	return func(c *tracerConfig) {
		c.serviceName = serviceName
	}
}

// WithEndpoint sets the LangSmith endpoint URL.
func WithEndpoint(endpoint string) TracerOption {
	return func(c *tracerConfig) {
		c.endpoint = endpoint
	}
}

// WithBatchTimeout sets the batch timeout for trace exports.
func WithBatchTimeout(timeout time.Duration) TracerOption {
	return func(c *tracerConfig) {
		c.batchTimeout = timeout
	}
}

// Tracer manages an OpenTelemetry tracer provider for LangSmith.
type Tracer struct {
	tp *sdktrace.TracerProvider
}

// NewTracer creates a new Tracer with the given options.
func NewTracer(opts ...TracerOption) (*Tracer, error) {
	cfg := &tracerConfig{
		endpoint:     defaultEndpoint,
		batchTimeout: defaultBatchTimeout,
	}

	for _, opt := range opts {
		opt(cfg)
	}

	if cfg.apiKey == "" {
		return nil, fmt.Errorf("API key is required (use WithAPIKey)")
	}

	if cfg.projectName == "" {
		return nil, fmt.Errorf("project name is required (use WithProjectName)")
	}

	ctx := context.Background()

	// Create resource with service name if provided
	var res *resource.Resource
	var err error
	if cfg.serviceName != "" {
		res, err = resource.New(ctx,
			resource.WithAttributes(
				semconv.ServiceName(cfg.serviceName),
			),
		)
	} else {
		res, err = resource.New(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Create OTLP HTTP exporter with LangSmith endpoint and headers
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.endpoint),
		otlptracehttp.WithURLPath(defaultURLPath),
		otlptracehttp.WithHeaders(map[string]string{
			"x-api-key":         cfg.apiKey,
			"Langsmith-Project": cfg.projectName,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(cfg.batchTimeout)),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider (for backward compatibility)
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Tracer{tp: tp}, nil
}

// TracerProvider returns the underlying trace.TracerProvider.
func (t *Tracer) TracerProvider() trace.TracerProvider {
	return t.tp
}

// Shutdown gracefully shuts down the tracer provider.
func (t *Tracer) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, defaultShutdownTimeout)
	defer cancel()
	return t.tp.Shutdown(shutdownCtx)
}
