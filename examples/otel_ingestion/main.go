package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to send OpenTelemetry traces to LangSmith UI.
//
// This is a mock/demo example that simulates LLM calls without requiring API keys.
// It demonstrates the tracing structure and waterfall visualization.
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT: Your LangSmith project name (defaults to "default")
//
// Running:
//
//	go run ./examples/otel_ingestion

const (
	defaultProjectName     = "default"
	otelEndpoint           = "https://api.smith.langchain.com/otel/v1/traces"
	serviceName            = "langsmith-go"
	tracerName             = "langsmith.go.example"
	traceFlushWait         = 2 * time.Second
	llmSpan1Duration       = 500 * time.Millisecond
	toolSpanDuration       = 300 * time.Millisecond
	retrieverSpanDuration  = 200 * time.Millisecond
	llmSpan2Duration       = 400 * time.Millisecond
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("=== LangSmith OpenTelemetry Example ===")
	fmt.Println()

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	// Initialize OpenTelemetry tracer
	ls, err := langsmith.NewTracer(
		langsmith.WithAPIKey(cfg.apiKey),
		langsmith.WithProjectName(cfg.projectName),
		langsmith.WithServiceName(serviceName),
	)
	if err != nil {
		return fmt.Errorf("initializing tracer: %w", err)
	}
	defer ls.Shutdown(context.Background())

	// Create spans
	ctx := context.Background()
	tracer := ls.Tracer(tracerName)
	sessionID := uuid.New().String()

	printWaterfallStructure()

	ctx, rootSpan := createRootSpan(ctx, tracer, sessionID)
	defer func() {
		// Set input and output on root span
		rootSpan.SetAttributes(
			attribute.String("gen_ai.prompt", "What's the weather in San Francisco?"),
			attribute.String("gen_ai.completion", "The weather in San Francisco is sunny with a temperature of 72°F."),
		)
		rootSpan.End()
	}()

	createChildSpans(ctx, tracer, sessionID)

	flushTraces()

	return nil
}

// config holds the application configuration.
type config struct {
	apiKey     string
	projectName string
}

// loadConfig loads configuration from environment variables.
func loadConfig() (*config, error) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		return nil, fmt.Errorf("LANGSMITH_API_KEY environment variable is required")
	}

	projectName := os.Getenv("LANGSMITH_PROJECT")
	if projectName == "" {
		projectName = defaultProjectName
	}

	return &config{
		apiKey:      apiKey,
		projectName: projectName,
	}, nil
}

// printConfig prints the configuration.
func printConfig(cfg *config) {
	fmt.Println("Configuration:")
	fmt.Printf("  Endpoint: %s\n", otelEndpoint)
	fmt.Printf("  Project: %s\n", cfg.projectName)
	fmt.Printf("  Service name: %s\n", serviceName)
	fmt.Println()
}

// printWaterfallStructure prints the expected waterfall structure.
func printWaterfallStructure() {
	fmt.Println("Creating waterfall with 5 spans:")
	fmt.Println("  1. agent.chain (root, 2s)")
	fmt.Println("     ├─ 2. openai.llm (500ms)")
	fmt.Println("     ├─ 3. weather.tool (300ms)")
	fmt.Println("     └─ 4. openai.llm (600ms)")
	fmt.Println("        └─ 5. database.retriever (200ms)")
	fmt.Println()
}

// createRootSpan creates the root agent chain span.
func createRootSpan(ctx context.Context, tracer trace.Tracer, sessionID string) (context.Context, trace.Span) {
	return tracer.Start(ctx, "agent.chain",
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "chat"),
			attribute.String("service.name", serviceName),
			attribute.String("session.id", sessionID),
		),
	)
}

// createChildSpans creates all child spans in the waterfall.
func createChildSpans(ctx context.Context, tracer trace.Tracer, sessionID string) {
	// CHILD 1: First LLM call
	ctx, llmSpan1 := createLLMSpan(ctx, tracer, sessionID, "openai.llm.call",
		"What's the weather in San Francisco?")
	time.Sleep(llmSpan1Duration)
	setLLMSpanCompletion(llmSpan1, "Let me check the weather for you.", 15, 12)
	llmSpan1.End()

	// CHILD 2: Tool call
	ctx, toolSpan := createToolSpan(ctx, tracer, sessionID, "weather.tool", "get_weather")
	time.Sleep(toolSpanDuration)
	toolSpan.SetAttributes(
		attribute.String("gen_ai.completion", `{"location":"San Francisco","temperature":"72°F","condition":"Sunny"}`),
	)
	toolSpan.End()

	// CHILD 3: Second LLM call with nested retriever
	ctx, llmSpan2 := createLLMSpan(ctx, tracer, sessionID, "openai.llm.final",
		"Based on the weather data, provide a summary.")

	// NESTED CHILD: Retriever call inside LLM
	_, retrieverSpan := createRetrieverSpan(ctx, tracer, sessionID)
	time.Sleep(retrieverSpanDuration)
	retrieverSpan.SetAttributes(
		attribute.String("gen_ai.completion", "Temperature: 72F, Sunny"),
	)
	retrieverSpan.End()

	time.Sleep(llmSpan2Duration)
	setLLMSpanCompletion(llmSpan2, "The weather in San Francisco is sunny with a temperature of 72°F.", 25, 18)
	llmSpan2.End()
}

// createLLMSpan creates a span for an LLM call.
func createLLMSpan(ctx context.Context, tracer trace.Tracer, sessionID, spanName, prompt string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "chat"),
			attribute.String("gen_ai.system", "openai"),
			attribute.String("gen_ai.request.model", "gpt-4"),
			attribute.String("service.name", serviceName),
			attribute.String("session.id", sessionID),
			attribute.String("gen_ai.prompt", prompt),
		),
	)
}

// createToolSpan creates a span for a tool call.
func createToolSpan(ctx context.Context, tracer trace.Tracer, sessionID, spanName, toolName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "tool"),
			attribute.String("tool.name", toolName),
			attribute.String("service.name", serviceName),
			attribute.String("session.id", sessionID),
			attribute.String("gen_ai.prompt", `{"location":"San Francisco"}`),
		),
	)
}

// createRetrieverSpan creates a span for a retriever call.
func createRetrieverSpan(ctx context.Context, tracer trace.Tracer, sessionID string) (context.Context, trace.Span) {
	return tracer.Start(ctx, "database.retriever",
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "retrieval"),
			attribute.String("service.name", serviceName),
			attribute.String("session.id", sessionID),
			attribute.String("gen_ai.prompt", "weather forecast data"),
		),
	)
}

// setLLMSpanCompletion sets completion and usage attributes on an LLM span.
func setLLMSpanCompletion(span trace.Span, completion string, inputTokens, outputTokens int64) {
	span.SetAttributes(
		attribute.String("gen_ai.completion", completion),
		attribute.Int64("gen_ai.usage.input_tokens", inputTokens),
		attribute.Int64("gen_ai.usage.output_tokens", outputTokens),
	)
}

// flushTraces flushes traces and prints completion message.
func flushTraces() {
	fmt.Println("\nAll spans ended. Flushing to LangSmith...")
	time.Sleep(traceFlushWait)
}


