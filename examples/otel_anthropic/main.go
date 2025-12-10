package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/anthropics/anthropic-sdk-go"
	"github.com/anthropics/anthropic-sdk-go/option"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/langchain-ai/langsmith-go/examples/otel_anthropic/traceanthropic"
)

// Demonstrates automatic OpenTelemetry tracing for Anthropic API calls.
//
// This example shows:
//   - Using traceanthropic.Client() to automatically trace Anthropic HTTP requests
//   - No manual span creation needed - tracing happens automatically
//   - All Anthropic API calls are automatically instrumented with LangSmith attributes
//   - Viewing rich traces in the LangSmith dashboard
//
// Prerequisites:
//   - ANTHROPIC_API_KEY: Your Anthropic API key
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT: Your LangSmith project name (defaults to "default")
//
// Running:
//
//	go run ./examples/otel_anthropic

const (
	defaultProjectName    = "default"
	serviceName           = "langsmith-go-anthropic-auto"
	traceFlushWait        = 2 * time.Second
	tracerShutdownTimeout = 10 * time.Second
	batchTimeout          = 1 * time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("=== Anthropic Automatic Tracing Example ===")
	fmt.Println()

	// Load configuration
	langsmithKey := os.Getenv("LANGSMITH_API_KEY")
	if langsmithKey == "" {
		return fmt.Errorf("LANGSMITH_API_KEY environment variable is required")
	}

	anthropicKey := os.Getenv("ANTHROPIC_API_KEY")
	if anthropicKey == "" {
		return fmt.Errorf("ANTHROPIC_API_KEY environment variable is required")
	}

	projectName := getProjectName()

	fmt.Printf("Configuration:\n")
	fmt.Printf("  LangSmith Project: %s\n", projectName)
	fmt.Printf("  Service Name: %s\n", serviceName)
	fmt.Println()

	// Initialize OpenTelemetry tracer
	shutdown, err := initTracer(langsmithKey, projectName)
	if err != nil {
		return fmt.Errorf("initializing tracer: %w", err)
	}
	defer shutdown()

	fmt.Println("✓ OpenTelemetry configured for LangSmith")
	fmt.Println()

	// Create Anthropic client with automatic tracing
	// The traceanthropic.Client() wraps the HTTP client to automatically trace all requests
	client := anthropic.NewClient(
		option.WithAPIKey(anthropicKey),
		option.WithHTTPClient(traceanthropic.Client()),
	)

	fmt.Println("✓ Anthropic client configured with automatic tracing")
	fmt.Println("  All API calls will be automatically traced!")
	fmt.Println()

	// Make Anthropic API calls - tracing happens automatically!
	ctx := context.Background()

	fmt.Println("Making Anthropic API calls...")
	fmt.Println()

	// Example 1: Simple message
	fmt.Println("1. Simple message:")
	msg1, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.Model("claude-3-opus-20240229"),
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.ContentBlockParamUnion{
				OfText: &anthropic.TextBlockParam{
					Text: "What is the capital of France?",
				},
			}),
		},
		MaxTokens: 1024,
	})
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}

	if len(msg1.Content) > 0 && msg1.Content[0].Type == "text" && msg1.Content[0].Text != "" {
		fmt.Printf("   Response: %s\n", msg1.Content[0].Text)
	}
	totalTokens := msg1.Usage.InputTokens + msg1.Usage.OutputTokens
	if totalTokens > 0 {
		fmt.Printf("   Tokens used: %d\n", totalTokens)
	}
	fmt.Println()

	// Example 2: Message with system prompt
	fmt.Println("2. Message with system prompt:")
	msg2, err := client.Messages.New(ctx, anthropic.MessageNewParams{
		Model: anthropic.Model("claude-3-opus-20240229"),
		System: []anthropic.TextBlockParam{
			{Text: "You are a helpful assistant that provides concise answers."},
		},
		Messages: []anthropic.MessageParam{
			anthropic.NewUserMessage(anthropic.ContentBlockParamUnion{
				OfText: &anthropic.TextBlockParam{
					Text: "Explain quantum computing in one sentence.",
				},
			}),
		},
		MaxTokens: 1024,
	})
	if err != nil {
		return fmt.Errorf("creating message: %w", err)
	}

	if len(msg2.Content) > 0 && msg2.Content[0].Type == "text" && msg2.Content[0].Text != "" {
		fmt.Printf("   Response: %s\n", msg2.Content[0].Text)
	}
	totalTokens2 := msg2.Usage.InputTokens + msg2.Usage.OutputTokens
	if totalTokens2 > 0 {
		fmt.Printf("   Tokens used: %d\n", totalTokens2)
	}
	fmt.Println()

	// Flush traces
	fmt.Println("Flushing traces to LangSmith...")
	time.Sleep(traceFlushWait)
	fmt.Println("✓ All traces flushed successfully")
	fmt.Println()

	fmt.Println("=== Summary ===")
	fmt.Println("✓ Made 2 Anthropic API calls")
	fmt.Println("✓ All calls were automatically traced with OpenTelemetry")
	fmt.Println("✓ Traces sent to LangSmith with proper Gen AI attributes")
	fmt.Println()
	fmt.Println("View your traces in LangSmith:")
	fmt.Printf("  https://smith.langchain.com/o/{tenant_id}/projects/{project_name}\n")
	fmt.Println()
	fmt.Println("Benefits of automatic tracing:")
	fmt.Println("  • No manual span creation needed")
	fmt.Println("  • All Anthropic requests/responses automatically captured")
	fmt.Println("  • Proper Gen AI semantic conventions applied")
	fmt.Println("  • Token usage and latency metrics included")
	fmt.Println("  • Cache token breakdown (if using cached prompts)")

	return nil
}

// initTracer initializes the OpenTelemetry tracer with LangSmith OTLP exporter.
func initTracer(apiKey, projectName string) (func(), error) {
	ctx := context.Background()

	// Create resource with service name
	res, err := resource.New(ctx,
		resource.WithAttributes(
			semconv.ServiceName(serviceName),
		),
	)
	if err != nil {
		return nil, fmt.Errorf("creating resource: %w", err)
	}

	// Create OTLP HTTP exporter with LangSmith endpoint and headers
	exporter, err := otlptracehttp.New(ctx,
		otlptracehttp.WithEndpoint("api.smith.langchain.com"),
		otlptracehttp.WithURLPath("/otel/v1/traces"),
		otlptracehttp.WithHeaders(map[string]string{
			"x-api-key":         apiKey,
			"Langsmith-Project": projectName,
		}),
	)
	if err != nil {
		return nil, fmt.Errorf("creating OTLP exporter: %w", err)
	}

	// Create tracer provider
	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exporter, sdktrace.WithBatchTimeout(batchTimeout)),
		sdktrace.WithResource(res),
	)

	// Set global tracer provider
	otel.SetTracerProvider(tp)
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// Return shutdown function
	return func() {
		shutdownCtx, cancel := context.WithTimeout(ctx, tracerShutdownTimeout)
		defer cancel()
		if err := tp.Shutdown(shutdownCtx); err != nil {
			fmt.Fprintf(os.Stderr, "Error shutting down tracer: %v\n", err)
		}
	}, nil
}

// getProjectName returns the project name from environment or default.
func getProjectName() string {
	if projectName := os.Getenv("LANGSMITH_PROJECT"); projectName != "" {
		return projectName
	}
	return defaultProjectName
}
