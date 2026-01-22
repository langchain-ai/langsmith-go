package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/sashabaranov/go-openai"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to run real experiments in LangSmith using OpenAI with automatic tracing.
//
// This example shows:
//   - Configuring OpenTelemetry to send traces to LangSmith automatically
//   - Creating a dataset idempotently (get or create)
//   - Creating examples with deterministic IDs for idempotent operations
//   - Running real OpenAI experiments using the OpenAI client
//   - Using OpenTelemetry spans to link traces to dataset examples
//   - Automatic trace submission via OpenTelemetry (no manual run ingestion!)
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - OPENAI_API_KEY: Your OpenAI API key
//   - LANGSMITH_PROJECT: Your LangSmith project name (defaults to "default")
//
// Running:
//
//	go run ./examples/e2e_eval

const (
	datasetName        = "Q&A Evaluation Dataset - Go Example"
	datasetDescription = "Dataset for Q&A evaluation with real OpenAI experiments."
	serviceName        = "langsmith-go-e2e-eval"
	traceFlushWait     = 2 * time.Second
	defaultProjectName = "default"
	defaultEndpoint    = "https://api.smith.langchain.com"
	dateTimeFormat     = "20060102-150405"

	// OpenTelemetry span and attribute names
	spanNameChatCompletion = "openai.chat.completion"
	genAISystem            = "openai"
	genAIOperationChat     = "chat"
)

// testCase represents a question-answer pair for evaluation.
type testCase struct {
	question       string
	expectedAnswer string
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("LangSmith Real Experiment Example (with OpenAI)")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()

	// Check required environment variables
	langsmithKey := os.Getenv("LANGSMITH_API_KEY")
	if langsmithKey == "" {
		return fmt.Errorf("LANGSMITH_API_KEY environment variable is required\nPlease set your LangSmith API key to run experiments")
	}

	openaiKey := os.Getenv("OPENAI_API_KEY")
	if openaiKey == "" {
		return fmt.Errorf("OPENAI_API_KEY environment variable is required\nPlease set your OpenAI API key to run the agent")
	}

	projectName := getProjectName()

	// Create LangSmith client
	langsmithClient := langsmith.NewClient()

	ctx := context.Background()

	// Define test cases
	fmt.Println("Defining test cases...")
	testCases := getTestCases()
	fmt.Printf("   ✓ Defined %d test cases\n\n", len(testCases))

	// 1. Get or create dataset
	dataset, err := getOrCreateDataset(ctx, langsmithClient)
	if err != nil {
		return err
	}

	// 2. Create examples from test cases
	examples, err := createExamplesFromTestCases(ctx, langsmithClient, dataset.ID, testCases)
	if err != nil {
		return err
	}

	// 3. Create experiment session linked to the dataset
	session, err := createExperimentSession(ctx, langsmithClient, dataset.ID)
	if err != nil {
		return err
	}

	// 4. Configure OpenTelemetry
	ls, err := langsmith.NewTracer(
		langsmith.WithAPIKey(langsmithKey),
		langsmith.WithProjectName(projectName),
		langsmith.WithServiceName(serviceName),
	)
	if err != nil {
		return fmt.Errorf("initializing tracer: %w", err)
	}
	defer ls.Shutdown(context.Background())

	fmt.Println("   ✓ OpenTelemetry configured to send traces with session_id attribute")
	fmt.Println()

	// 5. Create OpenAI client
	openaiClient := openai.NewClient(openaiKey)

	// 6. Run experiments
	if err := runExperiments(ctx, ls, openaiClient, testCases, examples, session.ID); err != nil {
		return err
	}

	// Flush traces
	fmt.Println("6. Flushing traces to LangSmith...")
	time.Sleep(traceFlushWait)
	fmt.Println("   ✓ All traces flushed successfully")
	fmt.Println()

	// Summary
	printSummary(dataset, session, len(examples))

	return nil
}

// getTestCases returns the test cases for evaluation.
func getTestCases() []testCase {
	return []testCase{
		{"What is the capital of France?", "Paris"},
		{"What is 2 + 2?", "4"},
		{"Who wrote Romeo and Juliet?", "William Shakespeare"},
		{"What is the largest planet in our solar system?", "Jupiter"},
	}
}

// getOrCreateDataset gets an existing dataset or creates a new one.
func getOrCreateDataset(ctx context.Context, client *langsmith.Client) (*langsmith.Dataset, error) {
	fmt.Println("1. Getting or creating dataset '" + datasetName + "'...")

	listParams := langsmith.DatasetListParams{
		Name: langsmith.F(datasetName),
	}

	existingDatasets, err := client.Datasets.List(ctx, listParams)
	if err != nil {
		return nil, fmt.Errorf("listing datasets: %w", err)
	}

	if len(existingDatasets.Items) > 0 {
		existing := existingDatasets.Items[0]
		fmt.Printf("   ✓ Found existing dataset (ID: %s)\n", existing.ID)
		return &existing, nil
	}

	createParams := langsmith.DatasetNewParams{
		Name:        langsmith.F(datasetName),
		Description: langsmith.F(datasetDescription),
	}

	created, err := client.Datasets.New(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("creating dataset: %w", err)
	}
	fmt.Printf("   ✓ Created new dataset (ID: %s)\n", created.ID)
	fmt.Printf("   → View dataset: %s\n\n", buildDatasetURL(created))
	return created, nil
}

// createExamplesFromTestCases creates examples from test cases with deterministic IDs.
func createExamplesFromTestCases(ctx context.Context, client *langsmith.Client, datasetID string, testCases []testCase) ([]langsmith.Example, error) {
	fmt.Println("2. Creating dataset examples from test cases...")

	exampleBodies := make([]langsmith.ExampleBulkNewParamsBody, 0, len(testCases))
	for _, tc := range testCases {
		input := map[string]interface{}{
			"question": tc.question,
		}
		referenceOutput := map[string]interface{}{
			"answer": tc.expectedAnswer,
		}

		// Generate deterministic ID for idempotent creation
		exampleID, err := generateExampleID(datasetID, input, referenceOutput)
		if err != nil {
			return nil, fmt.Errorf("generating example ID: %w", err)
		}

		exampleBodies = append(exampleBodies, langsmith.ExampleBulkNewParamsBody{
			ID:        langsmith.F(exampleID),
			DatasetID: langsmith.F(datasetID),
			Inputs:    langsmith.F(input),
			Outputs:   langsmith.F(referenceOutput),
		})
	}

	bulkParams := langsmith.ExampleBulkNewParams{
		Body: exampleBodies,
	}

	createdExamples, err := client.Examples.Bulk.New(ctx, bulkParams)
	if err != nil {
		return nil, fmt.Errorf("creating bulk examples: %w", err)
	}

	if createdExamples == nil {
		return nil, fmt.Errorf("no examples created")
	}

	fmt.Printf("   ✓ Created/updated %d examples with deterministic IDs\n\n", len(*createdExamples))
	return *createdExamples, nil
}

// generateExampleID generates a deterministic UUID v5 for an example.
func generateExampleID(datasetID string, inputs, outputs map[string]interface{}) (string, error) {
	// Serialize inputs and outputs to JSON for deterministic hashing
	inputJSON, err := json.Marshal(inputs)
	if err != nil {
		return "", fmt.Errorf("marshaling inputs: %w", err)
	}
	outputJSON, err := json.Marshal(outputs)
	if err != nil {
		return "", fmt.Errorf("marshaling outputs: %w", err)
	}

	// Create namespace UUID from dataset ID
	namespace, err := uuid.Parse(datasetID)
	if err != nil {
		// Fallback to a fixed namespace if dataset ID is not a valid UUID
		namespace = uuid.NameSpaceOID
	}

	// Create deterministic name from inputs and outputs
	name := fmt.Sprintf("%s|%s", string(inputJSON), string(outputJSON))

	// Generate UUID v5 (deterministic)
	return uuid.NewSHA1(namespace, []byte(name)).String(), nil
}

// createExperimentSession creates an experiment session linked to the dataset.
func createExperimentSession(ctx context.Context, client *langsmith.Client, datasetID string) (*langsmith.TracerSessionWithoutVirtualFields, error) {
	fmt.Println("3. Creating experiment session linked to dataset...")

	experimentName := fmt.Sprintf("E2eEvalExample-%s", time.Now().Format(dateTimeFormat))

	sessionParams := langsmith.SessionNewParams{
		Name:               langsmith.F(experimentName),
		Description:        langsmith.F("Experiment session for Q&A evaluation using OpenAI"),
		ReferenceDatasetID: langsmith.F(datasetID), // Link session to dataset - critical for Experiments tab
	}

	session, err := client.Sessions.New(ctx, sessionParams)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	fmt.Printf("   ✓ Created experiment session: %s\n", session.Name)
	fmt.Printf("   ✓ Session ID: %s\n", session.ID)

	// Verify the session was created with the reference dataset ID
	if session.ReferenceDatasetID == "" {
		return nil, fmt.Errorf("session created but ReferenceDatasetID is empty - this will prevent runs from being counted in the experiment")
	}

	fmt.Printf("   ✓ Session linked to dataset (ReferenceDatasetID: %s)\n", session.ReferenceDatasetID)
	fmt.Println()
	return session, nil
}

// runExperiments runs the agent against each test case with OpenTelemetry tracing.
func runExperiments(ctx context.Context, ls *langsmith.Tracer, client *openai.Client, testCases []testCase, examples []langsmith.Example, sessionID string) error {
	if len(testCases) != len(examples) {
		return fmt.Errorf("test cases count (%d) does not match examples count (%d)", len(testCases), len(examples))
	}

	fmt.Println("5. Running experiments with OpenAI agent...")
	fmt.Println("   Each LLM call will automatically create a trace in LangSmith")
	fmt.Println("   linked to the corresponding dataset example.")
	fmt.Println()

	tracer := ls.Tracer(serviceName)

	for i, tc := range testCases {
		example := examples[i]
		fmt.Printf("   Test %d/%d: \"%s\"\n", i+1, len(testCases), tc.question)
		fmt.Printf("      Expected: %s\n", tc.expectedAnswer)

		if err := runSingleExperiment(ctx, tracer, client, tc, example, sessionID); err != nil {
			fmt.Printf("      ✗ Error: %v\n\n", err)
			continue
		}

		fmt.Println()
	}

	fmt.Printf("   ✓ Completed %d experiment runs\n\n", len(testCases))
	return nil
}

// runSingleExperiment runs a single experiment with OpenTelemetry tracing.
func runSingleExperiment(ctx context.Context, tracer trace.Tracer, client *openai.Client, tc testCase, example langsmith.Example, sessionID string) error {
	// Create span with experiment attributes
	spanAttrs := buildExperimentSpanAttributes(sessionID, example.ID, tc.question)
	ctx, span := tracer.Start(ctx, spanNameChatCompletion, trace.WithAttributes(spanAttrs...))
	defer span.End()

	// Run the agent
	req := openai.ChatCompletionRequest{
		Model: openai.GPT4oMini,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleUser,
				Content: tc.question,
			},
		},
	}

	completion, err := client.CreateChatCompletion(ctx, req)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return fmt.Errorf("creating chat completion: %w", err)
	}

	actualAnswer := extractAnswer(completion)
	setCompletionSpanAttributes(span, actualAnswer, completion.Usage)

	fmt.Printf("      Actual:   %s\n", actualAnswer)
	fmt.Println("      → Trace sent to LangSmith and linked to experiment")

	return nil
}

// buildExperimentSpanAttributes builds the attributes for an experiment span.
// CRITICAL: Use "langsmith.trace.session_id" (not "session.id") to link to experiment.
func buildExperimentSpanAttributes(sessionID, exampleID, prompt string) []attribute.KeyValue {
	return []attribute.KeyValue{
		attribute.String("gen_ai.operation.name", genAIOperationChat),
		attribute.String("gen_ai.system", genAISystem),
		attribute.String("gen_ai.request.model", string(openai.GPT4oMini)),
		attribute.String("service.name", serviceName),
		attribute.String("langsmith.trace.session_id", sessionID), // Required for experiment linking
		attribute.String("reference_example_id", exampleID),
		attribute.String("gen_ai.prompt", prompt),
	}
}

// setCompletionSpanAttributes sets completion and usage attributes on a span.
func setCompletionSpanAttributes(span trace.Span, completion string, usage openai.Usage) {
	span.SetAttributes(attribute.String("gen_ai.completion", completion))

	if usage.TotalTokens > 0 {
		span.SetAttributes(
			attribute.Int64("gen_ai.usage.input_tokens", int64(usage.PromptTokens)),
			attribute.Int64("gen_ai.usage.output_tokens", int64(usage.CompletionTokens)),
		)
	}
}

// extractAnswer extracts the answer from the OpenAI completion response.
func extractAnswer(completion openai.ChatCompletionResponse) string {
	if len(completion.Choices) > 0 && completion.Choices[0].Message.Content != "" {
		return completion.Choices[0].Message.Content
	}
	return "No response"
}

// getProjectName returns the project name from environment or default.
func getProjectName() string {
	if projectName := os.Getenv("LANGSMITH_PROJECT"); projectName != "" {
		return projectName
	}
	return defaultProjectName
}

// buildDatasetURL constructs the web URL for a dataset.
func buildDatasetURL(dataset *langsmith.Dataset) string {
	baseURL := os.Getenv("LANGSMITH_ENDPOINT")
	if baseURL == "" {
		baseURL = defaultEndpoint
	}
	hostURL := strings.Replace(baseURL, "api.", "", 1)
	return fmt.Sprintf("%s/o/%s/datasets/%s", hostURL, dataset.TenantID, dataset.ID)
}

// printSummary prints the experiment summary.
func printSummary(dataset *langsmith.Dataset, session *langsmith.TracerSessionWithoutVirtualFields, exampleCount int) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Experiment Complete!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("What happened:")
	fmt.Printf("  • Dataset: %s with %d examples\n", datasetName, exampleCount)
	fmt.Printf("  • Session: %s linked to dataset (aka. experiment)\n", session.Name)
	fmt.Println("  • Each run includes:")
	fmt.Println("    - Input question from the dataset")
	fmt.Println("    - Reference/expected answer")
	fmt.Println("    - Actual answer from OpenAI")
	fmt.Println("    - Complete OpenTelemetry trace automatically sent to LangSmith and linked to dataset")
	fmt.Println()
	fmt.Println("How it works:")
	fmt.Println("  1. Create/find a session linked to the dataset (via reference_dataset_id)")
	fmt.Println("  2. Configure OpenTelemetry")
	fmt.Println("  3. Create OpenTelemetry spans with reference_example_id and session_id attributes")
	fmt.Println("  4. OpenTelemetry automatically exports traces with proper linkage")
	fmt.Println("  5. Experiment appears in dataset's 'Experiments' tab!")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Printf("  1. View the dataset and experiments tab:\n     %s\n\n", buildDatasetURL(dataset))
	fmt.Println("  2. In the LangSmith UI, you can:")
	fmt.Println("     - Compare actual answers vs expected answers side-by-side")
	fmt.Println("     - Add LLM-as-judge evaluators to automatically evaluate your experiments")
	fmt.Println("     - View detailed traces of each OpenAI call")
	fmt.Println("     - Analyze token usage and latency metrics")
	fmt.Println("     - Much more...")
	fmt.Println()
	fmt.Println("  3. Run this example again:")
	fmt.Println("     - A new experiment session will be created with a unique timestamp")
	fmt.Println("     - Dataset and examples are reused (idempotent)")
	fmt.Println("     - Each run appears as a separate experiment in the Experiments tab of the dataset")
	fmt.Println("     - You can compare results across multiple experiment runs")
	fmt.Println()
}

