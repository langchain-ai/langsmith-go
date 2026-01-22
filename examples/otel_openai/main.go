package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to make real OpenAI API calls with OpenTelemetry tracing to LangSmith.
//
// This example demonstrates:
//   - Configuring OpenTelemetry to send traces to LangSmith
//   - Making actual API calls to OpenAI with tool definitions
//   - Manual tool call span creation
//   - Multi-turn conversations with tool execution
//   - Viewing rich traces in the LangSmith dashboard
//
// Prerequisites:
//   - OPENAI_API_KEY: Your OpenAI API key
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT: Your LangSmith project name (defaults to "default")
//
// Running:
//
//	go run ./examples/otel_openai

const (
	defaultProjectName     = "default"
	serviceName            = "langsmith-go-openai-example"
	separator              = "============================================================"
	traceFlushWaitDuration = 2 * time.Second
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("=== OpenAI + LangSmith OpenTelemetry Example ===")
	fmt.Println()

	// Load configuration
	cfg, err := loadConfig()
	if err != nil {
		return err
	}

	printConfig(cfg)

	// Initialize OpenTelemetry tracer
	ls, err := langsmith.NewTracer(
		langsmith.WithAPIKey(cfg.langsmithKey),
		langsmith.WithProjectName(cfg.projectName),
		langsmith.WithServiceName(serviceName),
	)
	if err != nil {
		return fmt.Errorf("initializing tracer: %w", err)
	}
	defer ls.Shutdown(context.Background())

	fmt.Println("✓ OpenTelemetry configured for LangSmith")
	fmt.Println()

	// Create OpenAI client
	client := openai.NewClient(cfg.openaiKey)

	// Run the agent workflow
	ctx := context.Background()
	tracer := otel.Tracer(serviceName)

	ctx, workflowSpan := tracer.Start(ctx, "openai_agent_workflow",
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "agent_workflow"),
			attribute.String("langsmith.span.kind", "chain"),
			attribute.String("langsmith.trace.name", "OpenAI Agent with Tools"),
		),
	)
	defer workflowSpan.End()

	fmt.Println(separator)
	fmt.Println("Agent Workflow: Chat with Tool Calls")
	fmt.Println(separator)

	// Set initial input on root span
	initialPrompt := "What is the capital of France and what's the current weather there?"
	workflowSpan.SetAttributes(attribute.String("gen_ai.prompt", initialPrompt))

	finalContent, err := runAgentWorkflow(ctx, client, tracer)
	if err != nil {
		workflowSpan.RecordError(err)
		workflowSpan.SetStatus(codes.Error, err.Error())
		return err
	}

	// Display results
	printResults(finalContent, cfg.projectName)

	// Set output on root span
	workflowSpan.SetAttributes(
		attribute.String("gen_ai.completion", finalContent),
		attribute.String("response.content", finalContent),
	)
	workflowSpan.SetStatus(codes.Ok, "")

	// Flush traces
	flushTraces(cfg.projectName)

	return nil
}

// config holds the application configuration.
type config struct {
	openaiKey    string
	langsmithKey string
	projectName  string
}

// loadConfig loads configuration from environment variables.
func loadConfig() (*config, error) {
	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return nil, fmt.Errorf("OPENAI_API_KEY environment variable is required\nGet your API key from: https://platform.openai.com/api-keys")
	}

	langsmithKey := os.Getenv("LANGSMITH_API_KEY")
	if langsmithKey == "" {
		return nil, fmt.Errorf("LANGSMITH_API_KEY environment variable is required\nGet your API key from: https://smith.langchain.com/settings")
	}

	projectName := os.Getenv("LANGSMITH_PROJECT")
	if projectName == "" {
		projectName = defaultProjectName
	}

	return &config{
		openaiKey:    openaiKey,
		langsmithKey: langsmithKey,
		projectName:  projectName,
	}, nil
}

// printConfig prints the configuration.
func printConfig(cfg *config) {
	fmt.Println("Configuration:")
	fmt.Printf("  LangSmith Project: %s\n", cfg.projectName)
	fmt.Printf("  Service Name: %s\n", serviceName)
	fmt.Println()
}

// runAgentWorkflow executes the agent workflow with tool calls.
func runAgentWorkflow(ctx context.Context, client *openai.Client, tracer trace.Tracer) (string, error) {
	toolDefinition := createWeatherToolDefinition()
	initialMessages := []openai.ChatCompletionMessage{
		{
			Role:    openai.ChatMessageRoleUser,
			Content: "What is the capital of France and what's the current weather there?",
		},
	}

	fmt.Println("\n1. Making initial API call with tool definitions...")

	// Make initial API call
	ctx, llmSpan1 := createLLMSpan(ctx, tracer, "openai.chat.completion")
	promptText := buildPromptText(initialMessages)
	llmSpan1.SetAttributes(attribute.String("gen_ai.prompt", promptText))

	req := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: initialMessages,
		Tools: []openai.Tool{
			{
				Type:     openai.ToolTypeFunction,
				Function: &toolDefinition,
			},
		},
	}

	completion, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		recordSpanError(llmSpan1, err)
		return "", fmt.Errorf("creating chat completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		llmSpan1.End()
		return "", fmt.Errorf("no choices in completion response")
	}

	message := completion.Choices[0].Message
	var finalContent string

	// Set completion on initial span - this is the assistant's response (may include tool calls)
	firstSpanCompletion := message.Content
	if firstSpanCompletion == "" && len(message.ToolCalls) > 0 {
		firstSpanCompletion = "Tool calls requested"
	}
	setSpanCompletion(llmSpan1, firstSpanCompletion, completion.Usage)
	llmSpan1.End()

	// Handle tool calls if present
	if len(message.ToolCalls) > 0 {
		finalContent, err = handleToolCalls(ctx, client, tracer, initialMessages, message, completion)
		if err != nil {
			return "", err
		}
	} else {
		finalContent = message.Content
	}

	return finalContent, nil
}

// createWeatherToolDefinition creates the weather tool definition.
func createWeatherToolDefinition() openai.FunctionDefinition {
	return openai.FunctionDefinition{
		Name:        "get_weather",
		Description: "Get the current weather for a given location",
		Parameters: map[string]interface{}{
			"type": "object",
			"properties": map[string]interface{}{
				"location": map[string]interface{}{
					"type":        "string",
					"description": "The city and state, e.g., San Francisco, CA",
				},
			},
			"required": []string{"location"},
		},
	}
}

// createLLMSpan creates a span for an LLM call.
func createLLMSpan(ctx context.Context, tracer trace.Tracer, spanName string) (context.Context, trace.Span) {
	return tracer.Start(ctx, spanName,
		trace.WithAttributes(
			attribute.String("gen_ai.operation.name", "chat"),
			attribute.String("gen_ai.system", "openai"),
			attribute.String("gen_ai.request.model", "gpt-4o-mini"),
			attribute.String("service.name", serviceName),
		),
	)
}

// handleToolCalls processes tool calls and makes a follow-up request.
func handleToolCalls(ctx context.Context, client *openai.Client, tracer trace.Tracer, initialMessages []openai.ChatCompletionMessage, message openai.ChatCompletionMessage, initialCompletion openai.ChatCompletionResponse) (string, error) {
	fmt.Println("   ✓ Tool calls detected in response!")

	// Build messages list for follow-up request
	messages := make([]openai.ChatCompletionMessage, 0, len(initialMessages)+len(message.ToolCalls)+1)
	messages = append(messages, initialMessages[0]) // Original user message

	// Add assistant message with tool calls
	messages = append(messages, openai.ChatCompletionMessage{
		Role:       openai.ChatMessageRoleAssistant,
		Content:    message.Content,
		ToolCalls:  message.ToolCalls,
		ToolCallID: "",
	})

	// Execute each tool call
	fmt.Println("\n2. Executing tool calls...")
	for _, toolCall := range message.ToolCalls {
		if toolCall.Type != openai.ToolTypeFunction {
			continue
		}

		toolName := toolCall.Function.Name
		toolArguments := toolCall.Function.Arguments

		fmt.Printf("   - Tool: %s | Args: %s\n", toolName, toolArguments)

		// Create and execute tool span
		_, toolSpan := tracer.Start(ctx, toolName,
			trace.WithAttributes(
				attribute.String("gen_ai.operation.name", "tool"),
				attribute.String("tool.name", toolName),
				attribute.String("service.name", serviceName),
				attribute.String("gen_ai.prompt", toolArguments),
			),
		)

		toolResult := executeTool(toolName, toolArguments)
		fmt.Printf("   - Result: %s\n", toolResult)

		toolSpan.SetAttributes(attribute.String("gen_ai.completion", toolResult))
		toolSpan.End()

		// Add tool result message
		messages = append(messages, openai.ChatCompletionMessage{
			Role:       openai.ChatMessageRoleTool,
			Content:    toolResult,
			ToolCallID: toolCall.ID,
		})
	}

	// Send follow-up request with tool results
	fmt.Println("\n3. Sending follow-up request with tool results...")

	ctx, llmSpan2 := createLLMSpan(ctx, tracer, "openai.chat.completion")
	followUpPromptText := buildPromptText(messages)
	llmSpan2.SetAttributes(attribute.String("gen_ai.prompt", followUpPromptText))

	followUpReq := openai.ChatCompletionRequest{
		Model:    openai.GPT4oMini,
		Messages: messages,
	}

	completion, err := client.CreateChatCompletion(ctx, followUpReq)
	if err != nil {
		recordSpanError(llmSpan2, err)
		return "", fmt.Errorf("creating follow-up chat completion: %w", err)
	}

	if len(completion.Choices) == 0 {
		llmSpan2.End()
		return "", fmt.Errorf("no choices in follow-up completion response")
	}

	finalContent := completion.Choices[0].Message.Content
	setSpanCompletion(llmSpan2, finalContent, completion.Usage)
	llmSpan2.End()

	return finalContent, nil
}

// recordSpanError records an error on a span.
func recordSpanError(span trace.Span, err error) {
	span.RecordError(err)
	span.SetStatus(codes.Error, err.Error())
	span.End()
}

// setSpanCompletion sets completion and usage attributes on a span.
func setSpanCompletion(span trace.Span, content string, usage openai.Usage) {
	span.SetAttributes(attribute.String("gen_ai.completion", content))
	if usage.TotalTokens > 0 {
		span.SetAttributes(
			attribute.Int64("gen_ai.usage.input_tokens", int64(usage.PromptTokens)),
			attribute.Int64("gen_ai.usage.output_tokens", int64(usage.CompletionTokens)),
		)
	}
}

// printResults prints the final results and token usage.
func printResults(finalContent, projectName string) {
	fmt.Println("\n" + separator)
	fmt.Println("Final Response:")
	fmt.Println(finalContent)
	fmt.Println(separator)
}

// flushTraces flushes traces and prints completion message.
func flushTraces(projectName string) {
	fmt.Println("\n" + separator)
	fmt.Println("Flushing traces to LangSmith...")
	time.Sleep(traceFlushWaitDuration)

	fmt.Println("✓ Traces sent successfully!")
	fmt.Printf("\nView your traces at:\n  https://smith.langchain.com/projects/%s\n", projectName)
	fmt.Println(separator)
	fmt.Println("\nNote: Check the trace waterfall in LangSmith UI to see:")
	fmt.Println("  - Parent workflow span (chain)")
	fmt.Println("  - Child LLM spans (manually created)")
	fmt.Println("  - Tool call spans (manually created)")
}

// buildPromptText builds a text representation of messages for span attributes.
func buildPromptText(messages []openai.ChatCompletionMessage) string {
	var prompt strings.Builder
	for i, msg := range messages {
		if i > 0 {
			prompt.WriteString("\n")
		}
		prompt.WriteString(fmt.Sprintf("%s: %s", msg.Role, msg.Content))
		if len(msg.ToolCalls) > 0 {
			for _, tc := range msg.ToolCalls {
				if tc.Type == openai.ToolTypeFunction {
					prompt.WriteString(fmt.Sprintf("\n[Tool Call: %s(%s)]", tc.Function.Name, tc.Function.Arguments))
				}
			}
		}
	}
	return prompt.String()
}

// executeTool simulates executing a tool based on its name and arguments.
func executeTool(toolName, arguments string) string {
	var args map[string]interface{}
	if err := json.Unmarshal([]byte(arguments), &args); err != nil {
		return fmt.Sprintf(`{"error": "Failed to parse arguments: %s"}`, err.Error())
	}

	switch toolName {
	case "get_weather":
		return executeWeatherTool(args)
	default:
		errorMap := map[string]string{
			"error": fmt.Sprintf("Unknown tool: %s", toolName),
		}
		errorJSON, _ := json.Marshal(errorMap)
		return string(errorJSON)
	}
}

// executeWeatherTool simulates a weather API call.
func executeWeatherTool(args map[string]interface{}) string {
	location := "unknown"
	if loc, ok := args["location"].(string); ok {
		location = loc
	}

	result := map[string]interface{}{
		"location":    location,
		"temperature": "18°C",
		"condition":    "Partly Cloudy",
		"humidity":     "65%",
		"wind":         "15 km/h",
	}

	resultJSON, err := json.Marshal(result)
	if err != nil {
		return fmt.Sprintf(`{"error": "Failed to serialize result: %s"}`, err.Error())
	}
	return string(resultJSON)
}


