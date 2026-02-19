package langsmith

import (
	"context"
	"fmt"
	"os"
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

// Tracer manages a LangSmith span processor registered on an OpenTelemetry tracer provider.
type Tracer struct {
	tp        *sdktrace.TracerProvider
	processor sdktrace.SpanProcessor
	ownsTP    bool
}

// New registers a LangSmith exporter on the provided TracerProvider.
//
// Example:
//
//	tp := sdktrace.NewTracerProvider()
//	defer tp.Shutdown(context.Background())
//	otel.SetTracerProvider(tp)
//
//	ls, err := langsmith.New(tp,
//		langsmith.WithAPIKey("your-api-key"),
//		langsmith.WithProjectName("my-project"),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer ls.Shutdown(context.Background())
func New(tp *sdktrace.TracerProvider, opts ...TracerOption) (*Tracer, error) {
	cfg := resolveConfig(opts)

	if cfg.apiKey == "" {
		return nil, fmt.Errorf("API key is required (use WithAPIKey or set LANGSMITH_API_KEY environment variable)")
	}

	processor, err := createProcessor(cfg)
	if err != nil {
		return nil, err
	}

	tp.RegisterSpanProcessor(processor)

	return &Tracer{
		tp:        tp,
		processor: processor,
		ownsTP:    false,
	}, nil
}

// NewTracer creates a new Tracer that owns its own TracerProvider.
// For sharing a TracerProvider with other libraries, use New instead.
func NewTracer(opts ...TracerOption) (*Tracer, error) {
	cfg := resolveConfig(opts)

	if cfg.apiKey == "" {
		return nil, fmt.Errorf("API key is required (use WithAPIKey or set LANGSMITH_API_KEY environment variable)")
	}

	ctx := context.Background()

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

	tp := sdktrace.NewTracerProvider(sdktrace.WithResource(res))

	processor, err := createProcessor(cfg)
	if err != nil {
		return nil, err
	}

	tp.RegisterSpanProcessor(processor)

	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	return &Tracer{
		tp:        tp,
		processor: processor,
		ownsTP:    true,
	}, nil
}

func resolveConfig(opts []TracerOption) *tracerConfig {
	cfg := &tracerConfig{
		endpoint:     defaultEndpoint,
		batchTimeout: defaultBatchTimeout,
	}
	for _, opt := range opts {
		opt(cfg)
	}
	if cfg.apiKey == "" {
		cfg.apiKey = os.Getenv("LANGSMITH_API_KEY")
	}
	if cfg.projectName == "" {
		cfg.projectName = os.Getenv("LANGSMITH_PROJECT")
	}
	if cfg.projectName == "" {
		cfg.projectName = "default"
	}
	return cfg
}

func createProcessor(cfg *tracerConfig) (sdktrace.SpanProcessor, error) {
	ctx := context.Background()

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

	return sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithBatchTimeout(cfg.batchTimeout),
	), nil
}

// TracerProvider returns the underlying trace.TracerProvider.
func (t *Tracer) TracerProvider() trace.TracerProvider {
	return t.tp
}

// Tracer returns a trace.Tracer with the given name.
// This is a convenience method that wraps TracerProvider().Tracer(name).
func (t *Tracer) Tracer(name string) trace.Tracer {
	return t.tp.Tracer(name)
}

// Shutdown gracefully shuts down the tracer.
// If the Tracer was created with New, only the LangSmith processor is shut down.
// If the Tracer was created with NewTracer, the entire TracerProvider is shut down.
func (t *Tracer) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, defaultShutdownTimeout)
	defer cancel()
	if t.ownsTP {
		return t.tp.Shutdown(shutdownCtx)
	}
	return t.processor.Shutdown(shutdownCtx)
}
