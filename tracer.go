package langsmith

import (
	"context"
	"errors"
	"fmt"
	"os"
	"strings"
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
	defaultEndpoint        = "api.smith.langchain.com"
	defaultURLPath         = "/otel/v1/traces"
	defaultBatchTimeout    = 1 * time.Second
	defaultShutdownTimeout = 10 * time.Second
)

// OTelTracerOption configures an [OTelTracer].
type OTelTracerOption func(*tracerConfig)

type tracerConfig struct {
	apiKey       string
	projectName  string
	address      AgentAddress
	serviceName  string
	endpoint     string
	batchTimeout time.Duration
}

// WithAPIKey sets the LangSmith API key.
func WithAPIKey(apiKey string) OTelTracerOption {
	return func(c *tracerConfig) {
		c.apiKey = apiKey
	}
}

// WithProjectName sets the LangSmith project name.
func WithProjectName(projectName string) OTelTracerOption {
	return func(c *tracerConfig) {
		c.projectName = projectName
	}
}

// WithAddress (beta) sends traces to an agent environment instead of a project. It cannot be
// combined with [WithProjectName]. Without either option, LANGSMITH_AGENT_ID and
// LANGSMITH_AGENT_ENVIRONMENT address the agent when both are set. Agent addressing is enabled
// per workspace.
func WithAddress(a AgentAddress) OTelTracerOption {
	return func(c *tracerConfig) {
		c.address = a
	}
}

// WithServiceName sets the service name for the tracer.
func WithServiceName(serviceName string) OTelTracerOption {
	return func(c *tracerConfig) {
		c.serviceName = serviceName
	}
}

// WithEndpoint sets the LangSmith endpoint URL.
func WithEndpoint(endpoint string) OTelTracerOption {
	return func(c *tracerConfig) {
		c.endpoint = endpoint
	}
}

// WithBatchTimeout sets the batch timeout for trace exports.
func WithBatchTimeout(timeout time.Duration) OTelTracerOption {
	return func(c *tracerConfig) {
		c.batchTimeout = timeout
	}
}

// OTelTracer manages a LangSmith span processor registered on an OpenTelemetry tracer provider.
type OTelTracer struct {
	tp        *sdktrace.TracerProvider
	processor sdktrace.SpanProcessor
	ownsTP    bool
}

// NewOTel registers a LangSmith exporter on the provided TracerProvider.
//
// Example:
//
//	tp := sdktrace.NewTracerProvider()
//	defer tp.Shutdown(context.Background())
//	otel.SetTracerProvider(tp)
//
//	ls, err := langsmith.NewOTel(tp,
//		langsmith.WithAPIKey("your-api-key"),
//		langsmith.WithProjectName("my-project"),
//	)
//	if err != nil {
//		log.Fatal(err)
//	}
//	defer ls.Shutdown(context.Background())
func NewOTel(tp *sdktrace.TracerProvider, opts ...OTelTracerOption) (*OTelTracer, error) {
	if tp == nil {
		return nil, fmt.Errorf("TracerProvider must not be nil (use NewOTelTracer to create one automatically)")
	}

	cfg, err := resolveConfig(opts)
	if err != nil {
		return nil, err
	}

	if cfg.apiKey == "" {
		return nil, fmt.Errorf("API key is required (use WithAPIKey or set LANGSMITH_API_KEY environment variable)")
	}

	processor, err := createProcessor(cfg)
	if err != nil {
		return nil, err
	}

	tp.RegisterSpanProcessor(processor)

	return &OTelTracer{
		tp:        tp,
		processor: processor,
		ownsTP:    false,
	}, nil
}

// NewOTelTracer creates a new [OTelTracer] that owns its own TracerProvider.
// For sharing a TracerProvider with other libraries, use [NewOTel] instead.
func NewOTelTracer(opts ...OTelTracerOption) (*OTelTracer, error) {
	cfg, err := resolveConfig(opts)
	if err != nil {
		return nil, err
	}

	if cfg.apiKey == "" {
		return nil, fmt.Errorf("API key is required (use WithAPIKey or set LANGSMITH_API_KEY environment variable)")
	}

	ctx := context.Background()

	var res *resource.Resource
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

	return &OTelTracer{
		tp:        tp,
		processor: processor,
		ownsTP:    true,
	}, nil
}

func resolveConfig(opts []OTelTracerOption) (*tracerConfig, error) {
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
	if !cfg.address.IsZero() && cfg.projectName != "" {
		return nil, errors.New("langsmith: WithAddress and WithProjectName cannot be combined")
	}
	if cfg.address.IsZero() && cfg.projectName == "" {
		// The environment decides when neither option names a destination.
		envAddress, err := AgentAddressFromEnv()
		if err != nil {
			return nil, err
		}
		envProject := os.Getenv("LANGSMITH_PROJECT")
		if !envAddress.IsZero() && envProject != "" {
			return nil, fmt.Errorf("%w: LANGSMITH_PROJECT and the LANGSMITH_AGENT_* variables are both set", ErrEnvAddress)
		}
		cfg.address, cfg.projectName = envAddress, envProject
	}
	if cfg.address.IsZero() && cfg.projectName == "" {
		cfg.projectName = "default"
	}
	return cfg, nil
}

// headers returns the OTLP request headers. An address replaces the project header, because
// ingestion rejects a request that sets both.
func (c *tracerConfig) headers() map[string]string {
	headers := map[string]string{"x-api-key": c.apiKey}
	if c.address.IsZero() {
		headers["Langsmith-Project"] = c.projectName
		return headers
	}
	headers["Langsmith-Agent-Id"] = c.address.ID()
	headers["Langsmith-Agent-Environment"] = strings.ToLower(c.address.Environment())
	return headers
}

func createProcessor(cfg *tracerConfig) (sdktrace.SpanProcessor, error) {
	ctx := context.Background()

	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint(cfg.endpoint),
		otlptracehttp.WithURLPath(defaultURLPath),
		otlptracehttp.WithHeaders(cfg.headers()),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	return sdktrace.NewBatchSpanProcessor(exporter,
		sdktrace.WithBatchTimeout(cfg.batchTimeout),
	), nil
}

// TracerProvider returns the underlying trace.TracerProvider.
func (t *OTelTracer) TracerProvider() trace.TracerProvider {
	return t.tp
}

// Tracer returns a trace.Tracer with the given name.
func (t *OTelTracer) Tracer(name string) trace.Tracer {
	return t.tp.Tracer(name)
}

// Shutdown gracefully shuts down the tracer.
// If the OTelTracer was created with [NewOTel], only the LangSmith processor is shut down.
// If it was created with [NewOTelTracer], the entire TracerProvider is shut down.
func (t *OTelTracer) Shutdown(ctx context.Context) error {
	shutdownCtx, cancel := context.WithTimeout(ctx, defaultShutdownTimeout)
	defer cancel()
	if t.ownsTP {
		return t.tp.Shutdown(shutdownCtx)
	}
	return t.processor.Shutdown(shutdownCtx)
}

// Deprecated: Use [OTelTracerOption] instead.
type TracerOption = OTelTracerOption

// Deprecated: Use [OTelTracer] instead.
type Tracer = OTelTracer

// Deprecated: Use [NewOTel] instead.
var New = NewOTel

// Deprecated: Use [NewOTelTracer] instead.
var NewTracer = NewOTelTracer
