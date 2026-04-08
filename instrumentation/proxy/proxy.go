// Package proxy provides a reverse proxy that transparently adds LangSmith
// tracing to Anthropic and OpenAI API calls. It routes requests by path to
// the correct upstream and uses the existing traceanthropic/traceopenai
// instrumentation to create OpenTelemetry spans.
//
// Usage:
//
//	p, err := proxy.New(proxy.Config{
//	    LangSmithAPIKey:  os.Getenv("LANGSMITH_API_KEY"),
//	    LangSmithProject: "my-project",
//	})
//	defer p.Shutdown(ctx)
//	http.ListenAndServe(":8090", p)
//
// Then point any SDK at the proxy:
//
//	export ANTHROPIC_BASE_URL=http://localhost:8090
//	export OPENAI_BASE_URL=http://localhost:8090/v1
package proxy

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"strings"

	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/instrumentation/traceanthropic"
	"github.com/langchain-ai/langsmith-go/instrumentation/traceopenai"
)

// Config configures a tracing proxy.
type Config struct {
	// Upstream URLs. Empty strings default to the real provider APIs.
	AnthropicUpstream string // default: https://api.anthropic.com
	OpenAIUpstream    string // default: https://api.openai.com

	// LangSmith credentials. Empty strings are read from LANGSMITH_API_KEY,
	// LANGSMITH_PROJECT, and LANGSMITH_ENDPOINT environment variables.
	LangSmithAPIKey   string
	LangSmithProject  string
	LangSmithEndpoint string

	// TracerProvider overrides the internally created provider. When set,
	// the proxy does not create or shut down its own provider. Use this
	// for testing with tracetest.InMemoryExporter.
	TracerProvider trace.TracerProvider
}

func (c *Config) resolve() {
	if c.AnthropicUpstream == "" {
		c.AnthropicUpstream = "https://api.anthropic.com"
	}
	if c.OpenAIUpstream == "" {
		c.OpenAIUpstream = "https://api.openai.com"
	}
	if c.LangSmithAPIKey == "" {
		c.LangSmithAPIKey = os.Getenv("LANGSMITH_API_KEY")
	}
	if c.LangSmithProject == "" {
		c.LangSmithProject = os.Getenv("LANGSMITH_PROJECT")
	}
	if c.LangSmithEndpoint == "" {
		c.LangSmithEndpoint = os.Getenv("LANGSMITH_ENDPOINT")
	}
}

// Proxy is an http.Handler that forwards LLM API requests to upstream
// providers while adding LangSmith tracing.
type Proxy struct {
	handler http.Handler
	tracer  *langsmith.OTelTracer // nil when an external TracerProvider was provided
}

// New creates a tracing proxy from the given configuration.
func New(cfg Config) (*Proxy, error) {
	cfg.resolve()

	anthropicURL, err := url.Parse(cfg.AnthropicUpstream)
	if err != nil {
		return nil, fmt.Errorf("parsing anthropic upstream: %w", err)
	}
	openaiURL, err := url.Parse(cfg.OpenAIUpstream)
	if err != nil {
		return nil, fmt.Errorf("parsing openai upstream: %w", err)
	}

	// Obtain a TracerProvider.
	var tp trace.TracerProvider
	var tracer *langsmith.OTelTracer

	if cfg.TracerProvider != nil {
		tp = cfg.TracerProvider
	} else {
		var opts []langsmith.OTelTracerOption
		if cfg.LangSmithAPIKey != "" {
			opts = append(opts, langsmith.WithAPIKey(cfg.LangSmithAPIKey))
		}
		if cfg.LangSmithProject != "" {
			opts = append(opts, langsmith.WithProjectName(cfg.LangSmithProject))
		}
		if cfg.LangSmithEndpoint != "" {
			// WithEndpoint expects a host, but LANGSMITH_ENDPOINT is often
			// a full URL. Strip the scheme to avoid double-wrapping.
			ep := strings.TrimPrefix(cfg.LangSmithEndpoint, "https://")
			ep = strings.TrimPrefix(ep, "http://")
			opts = append(opts, langsmith.WithEndpoint(ep))
		}
		opts = append(opts, langsmith.WithServiceName("langsmith-proxy"))

		t, tErr := langsmith.NewOTelTracer(opts...)
		if tErr != nil {
			return nil, fmt.Errorf("creating langsmith tracer: %w", tErr)
		}
		tracer = t
		tp = t.TracerProvider()
	}

	// Build traced transports for each provider.
	anthropicTransport := traceanthropic.WrapClient(nil,
		traceanthropic.WithTracerProvider(tp),
		traceanthropic.WithTraceAllHosts(),
	).Transport

	openaiTransport := traceopenai.WrapClient(nil,
		traceopenai.WithTracerProvider(tp),
	).Transport

	plainTransport := http.DefaultTransport

	rt := &routingTransport{
		anthropic:    anthropicTransport,
		openai:       openaiTransport,
		passthrough:  plainTransport,
		anthropicURL: anthropicURL,
		openaiURL:    openaiURL,
	}

	rp := &httputil.ReverseProxy{
		Director: func(req *http.Request) {
			// Route to the correct upstream based on request path.
			if isAnthropicPath(req.URL.Path) {
				req.URL.Scheme = anthropicURL.Scheme
				req.URL.Host = anthropicURL.Host
				req.Host = anthropicURL.Host
			} else {
				req.URL.Scheme = openaiURL.Scheme
				req.URL.Host = openaiURL.Host
				req.Host = openaiURL.Host
			}
		},
		Transport:     rt,
		FlushInterval: -1, // immediate flush for SSE streaming
	}

	return &Proxy{handler: rp, tracer: tracer}, nil
}

// ServeHTTP implements http.Handler.
func (p *Proxy) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	p.handler.ServeHTTP(w, r)
}

// Shutdown flushes pending spans and releases resources. If the proxy was
// created with an external TracerProvider, Shutdown is a no-op.
func (p *Proxy) Shutdown(ctx context.Context) error {
	if p.tracer != nil {
		return p.tracer.Shutdown(ctx)
	}
	return nil
}

// routingTransport dispatches requests to the correct traced transport
// based on the URL path.
type routingTransport struct {
	anthropic   http.RoundTripper
	openai      http.RoundTripper
	passthrough http.RoundTripper

	anthropicURL *url.URL
	openaiURL    *url.URL
}

func (rt *routingTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	path := req.URL.Path

	if isAnthropicPath(path) {
		return rt.anthropic.RoundTrip(req)
	}
	if isOpenAIPath(path) {
		return rt.openai.RoundTrip(req)
	}
	// Passthrough for /v1/models, health checks, etc.
	return rt.passthrough.RoundTrip(req)
}

func isAnthropicPath(path string) bool {
	return strings.Contains(path, "/v1/messages")
}

func isOpenAIPath(path string) bool {
	return strings.HasSuffix(path, "/chat/completions") ||
		strings.HasSuffix(path, "/completions") ||
		strings.HasSuffix(path, "/embeddings") ||
		strings.HasSuffix(path, "/responses")
}
