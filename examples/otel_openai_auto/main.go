package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"

	"github.com/langchain-ai/langsmith-go/examples/otel_openai_auto/traceopenai"
)

// Demonstrates automatic OpenTelemetry tracing for OpenAI API calls.
//
// This example shows:
//   - Using traceopenai.Client() to automatically trace OpenAI HTTP requests
//   - No manual span creation needed - tracing happens automatically
//   - All OpenAI API calls are automatically instrumented with LangSmith attributes
//   - Viewing rich traces in the LangSmith dashboard
//
// Prerequisites:
//   - OPENAI_API_KEY: Your OpenAI API key
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT: Your LangSmith project name (defaults to "default")
//
// Running:
//
//	go run ./examples/otel_openai_auto

const (
	defaultProjectName    = "default"
	serviceName           = "langsmith-go-openai-auto"
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
	fmt.Println("=== OpenAI Automatic Tracing Example ===")
	fmt.Println()

	// Load configuration
	langsmithKey := os.Getenv("LANGSMITH_API_KEY")
	if langsmithKey == "" {
		return fmt.Errorf("LANGSMITH_API_KEY environment variable is required")
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is required")
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

	// Create OpenAI client with automatic tracing
	// The traceopenai.Client() wraps the HTTP client to automatically trace all requests
	cfg := openai.DefaultConfig(openaiKey)
	cfg.HTTPClient = traceopenai.Client()
	client := openai.NewClientWithConfig(cfg)

	fmt.Println("✓ OpenAI client configured with automatic tracing")
	fmt.Println("  All API calls will be automatically traced!")
	fmt.Println()

	// Make OpenAI API calls - tracing happens automatically!
	ctx := context.Background()

	fmt.Println("Making OpenAI API calls...")
	fmt.Println()

	// Example 1: Simple chat completion
	fmt.Println("1. Simple chat completion:")
	req1 := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "What is the capital of France?",
			},
		},
	}

	completion1, err := client.CreateChatCompletion(ctx, req1)
	if err != nil {
		return fmt.Errorf("creating chat completion: %w", err)
	}

	if len(completion1.Choices) > 0 && completion1.Choices[0].Message.Content != "" {
		fmt.Printf("   Response: %s\n", completion1.Choices[0].Message.Content)
	}
	if completion1.Usage.TotalTokens > 0 {
		fmt.Printf("   Tokens used: %d\n", completion1.Usage.TotalTokens)
	}
	fmt.Println()

	// Example 2: Chat completion with system message
	fmt.Println("2. Chat completion with system message:")
	req2 := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: "You are a helpful assistant that provides concise answers.",
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Explain quantum computing in one sentence.",
			},
		},
	}

	completion2, err := client.CreateChatCompletion(ctx, req2)
	if err != nil {
		return fmt.Errorf("creating chat completion: %w", err)
	}

	if len(completion2.Choices) > 0 && completion2.Choices[0].Message.Content != "" {
		fmt.Printf("   Response: %s\n", completion2.Choices[0].Message.Content)
	}
	if completion2.Usage.TotalTokens > 0 {
		fmt.Printf("   Tokens used: %d\n", completion2.Usage.TotalTokens)
	}
	fmt.Println()

	// Flush traces
	fmt.Println("Flushing traces to LangSmith...")
	time.Sleep(traceFlushWait)
	fmt.Println("✓ All traces flushed successfully")
	fmt.Println()

	fmt.Println("=== Summary ===")
	fmt.Println("✓ Made 2 OpenAI API calls")
	fmt.Println("✓ All calls were automatically traced with OpenTelemetry")
	fmt.Println("✓ Traces sent to LangSmith with proper Gen AI attributes")
	fmt.Println()
	fmt.Println("View your traces in LangSmith:")
	fmt.Printf("  https://smith.langchain.com/o/{tenant_id}/projects/{project_name}\n")
	fmt.Println()
	fmt.Println("Benefits of automatic tracing:")
	fmt.Println("  • No manual span creation needed")
	fmt.Println("  • All OpenAI requests/responses automatically captured")
	fmt.Println("  • Proper Gen AI semantic conventions applied")
	fmt.Println("  • Token usage and latency metrics included")

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
