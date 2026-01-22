package main

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/baggage"
	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/examples/otel_go_client_openai/traceopenai"
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
//	go run ./examples/otel_go_client_openai

const (
	defaultProjectName = "default"
	serviceName        = "langsmith-go-openai-auto"
	traceFlushWait     = 2 * time.Second
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
	ls, err := langsmith.NewTracer(
		langsmith.WithAPIKey(langsmithKey),
		langsmith.WithProjectName(projectName),
		langsmith.WithServiceName(serviceName),
	)
	if err != nil {
		return fmt.Errorf("initializing tracer: %w", err)
	}
	defer ls.Shutdown(context.Background())

	fmt.Println("✓ OpenTelemetry configured for LangSmith")
	fmt.Println()

	// Create OpenAI client with automatic tracing
	// The traceopenai.Client() wraps the HTTP client to automatically trace all requests
	cfg := openai.DefaultConfig(openaiKey)
	cfg.HTTPClient = traceopenai.Client(traceopenai.WithTracerProvider(ls.TracerProvider()))
	client := openai.NewClientWithConfig(cfg)

	fmt.Println("✓ OpenAI client configured with automatic tracing")
	fmt.Println("  All API calls will be automatically traced!")
	fmt.Println()

	// Generate a thread ID to group all traces together
	// IMPORTANT: For threads, each API call should be its own trace (not one big trace)
	// All traces share the same session_id to group them into a thread
	threadID := uuid.New().String()
	fmt.Printf("Thread ID: %s\n", threadID)
	fmt.Println("  Each API call will be its own trace, all grouped into this thread")
	fmt.Println()

	tracer := ls.Tracer(serviceName)

	fmt.Println("Creating separate traces for each API call (grouped in thread)...")
	fmt.Println()
	fmt.Println("Each trace will have:")
	fmt.Println("  • agent.chain (root span)")
	fmt.Println("    └─ agent.step (child span)")
	fmt.Println("        └─ openai.chat.completion (auto-traced child)")
	fmt.Println()
	fmt.Println("All traces share session_id: " + threadID)
	fmt.Println()

	// Trace 1: First API call
	fmt.Println("Trace 1: Initial question")
	ctx1 := context.Background()
	
	// Propagate thread ID through baggage for this trace
	member1, _ := baggage.NewMember("session_id", threadID)
	bag1, _ := baggage.New(member1)
	ctx1 = baggage.ContextWithBaggage(ctx1, bag1)
	
	ctx1, trace1Root := tracer.Start(ctx1, "trace.1.initial_question",
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "chat"),
			attribute.String("service.name", serviceName),
			// Set thread metadata in multiple formats for compatibility
			attribute.String("session_id", threadID),
			attribute.String("langsmith.metadata.session_id", threadID),
			attribute.String("session.id", threadID),
		),
	)
	
	ctx1, trace1Step := tracer.Start(ctx1, "agent.step",
		trace.WithAttributes(
			attribute.String("step.type", "initial_query"),
			attribute.String("step.number", "1"),
			// Set thread metadata in multiple formats for compatibility
			attribute.String("session_id", threadID),
			attribute.String("langsmith.metadata.session_id", threadID),
			attribute.String("session.id", threadID),
		),
	)
	
	time.Sleep(50 * time.Millisecond)

	req1 := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "What is the capital of France?",
			},
		},
	}

	// Set input (prompt) on root span for LangSmith UI
	trace1Root.SetAttributes(attribute.String("gen_ai.prompt", "What is the capital of France?"))

	completion1, err := client.CreateChatCompletion(ctx1, req1)
	if err != nil {
		return fmt.Errorf("creating chat completion: %w", err)
	}

	// Extract completion text
	var completion1Text string
	if len(completion1.Choices) > 0 && completion1.Choices[0].Message.Content != "" {
		completion1Text = completion1.Choices[0].Message.Content
		fmt.Printf("   Response: %s\n", completion1Text)
	}
	
	// Set output on root span for LangSmith UI
	if completion1Text != "" {
		trace1Root.SetAttributes(attribute.String("gen_ai.completion", completion1Text))
	}
	
	if completion1.Usage.TotalTokens > 0 {
		fmt.Printf("   Tokens used: %d\n", completion1.Usage.TotalTokens)
	}
	trace1Step.End()
	trace1Root.End()
	time.Sleep(200 * time.Millisecond) // Delay between traces
	fmt.Println()

	// Trace 2: Second API call (same thread)
	fmt.Println("Trace 2: Another question")
	ctx2 := context.Background()
	
	member2, _ := baggage.NewMember("session_id", threadID)
	bag2, _ := baggage.New(member2)
	ctx2 = baggage.ContextWithBaggage(ctx2, bag2)
	
	ctx2, trace2Root := tracer.Start(ctx2, "trace.2.simple_math",
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "chat"),
			attribute.String("service.name", serviceName),
			attribute.String("session_id", threadID),
			attribute.String("langsmith.metadata.session_id", threadID),
			attribute.String("session.id", threadID),
		),
	)
	
	ctx2, trace2Step := tracer.Start(ctx2, "agent.step",
		trace.WithAttributes(
			attribute.String("step.type", "additional_query"),
			attribute.String("step.number", "2"),
			attribute.String("session_id", threadID),
			attribute.String("langsmith.metadata.session_id", threadID),
			attribute.String("session.id", threadID),
		),
	)
	
	time.Sleep(50 * time.Millisecond)

	req2 := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "What is 2+2?",
			},
		},
	}

	// Set input on root span for LangSmith UI
	trace2Root.SetAttributes(attribute.String("gen_ai.prompt", "What is 2+2?"))

	completion2, err := client.CreateChatCompletion(ctx2, req2)
	if err != nil {
		return fmt.Errorf("creating chat completion: %w", err)
	}

	// Extract completion text
	var completion2Text string
	if len(completion2.Choices) > 0 && completion2.Choices[0].Message.Content != "" {
		completion2Text = completion2.Choices[0].Message.Content
		fmt.Printf("   Response: %s\n", completion2Text)
	}
	
	// Set output on root span for LangSmith UI
	if completion2Text != "" {
		trace2Root.SetAttributes(attribute.String("gen_ai.completion", completion2Text))
	}
	
	if completion2.Usage.TotalTokens > 0 {
		fmt.Printf("   Tokens used: %d\n", completion2.Usage.TotalTokens)
	}
	trace2Step.End()
	trace2Root.End()
	time.Sleep(200 * time.Millisecond)
	fmt.Println()

	// Trace 3: Third API call (same thread)
	fmt.Println("Trace 3: Final question")
	ctx3 := context.Background()
	
	member3, _ := baggage.NewMember("session_id", threadID)
	bag3, _ := baggage.New(member3)
	ctx3 = baggage.ContextWithBaggage(ctx3, bag3)
	
	ctx3, trace3Root := tracer.Start(ctx3, "trace.3.famous_scientist",
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "chat"),
			attribute.String("service.name", serviceName),
			attribute.String("session_id", threadID),
			attribute.String("langsmith.metadata.session_id", threadID),
			attribute.String("session.id", threadID),
		),
	)
	
	ctx3, trace3Step := tracer.Start(ctx3, "agent.step",
		trace.WithAttributes(
			attribute.String("step.type", "final_query"),
			attribute.String("step.number", "3"),
			attribute.String("session_id", threadID),
			attribute.String("langsmith.metadata.session_id", threadID),
			attribute.String("session.id", threadID),
		),
	)
	
	time.Sleep(50 * time.Millisecond)

	req3 := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: "Name a famous scientist.",
			},
		},
	}

	// Set input on root span for LangSmith UI
	trace3Root.SetAttributes(attribute.String("gen_ai.prompt", "Name a famous scientist."))

	completion3, err := client.CreateChatCompletion(ctx3, req3)
	if err != nil {
		return fmt.Errorf("creating chat completion: %w", err)
	}

	// Extract completion text
	var completion3Text string
	if len(completion3.Choices) > 0 && completion3.Choices[0].Message.Content != "" {
		completion3Text = completion3.Choices[0].Message.Content
		fmt.Printf("   Response: %s\n", completion3Text)
	}
	
	// Set output on root span for LangSmith UI
	if completion3Text != "" {
		trace3Root.SetAttributes(attribute.String("gen_ai.completion", completion3Text))
	}
	
	if completion3.Usage.TotalTokens > 0 {
		fmt.Printf("   Tokens used: %d\n", completion3.Usage.TotalTokens)
	}
	trace3Step.End()
	trace3Root.End()
	time.Sleep(200 * time.Millisecond)
	fmt.Println()

	// Flush traces
	fmt.Println("Flushing traces to LangSmith...")
	time.Sleep(traceFlushWait)
	fmt.Println("✓ All traces flushed successfully")
	return nil
}


// getProjectName returns the project name from environment or default.
func getProjectName() string {
	if projectName := os.Getenv("LANGSMITH_PROJECT"); projectName != "" {
		return projectName
	}
	return defaultProjectName
}
