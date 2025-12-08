package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to record experiments in LangSmith using datasets, sessions, and runs.
//
// This example is inspired by the Python implementation at:
// https://github.com/jacoblee93/langsmith-rest-experiments
//
// This example shows:
//   - Defining experiment results with inputs, reference outputs, and actual outputs
//   - Getting or creating a dataset by name
//   - Generating deterministic example IDs using UUID v5
//   - Creating or updating examples in the dataset with idempotent IDs
//   - Creating an experiment session to group experiment runs
//   - Logging mocked LLM runs with inputs, outputs, and metadata
//   - Using LLM-as-judge evaluators to compare actual vs reference outputs
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//
// Running:
//
//	go run ./examples/record_experiment

// ExperimentResult represents a single experiment result with input, reference output, and actual output.
type ExperimentResult struct {
	Input           map[string]interface{} `json:"input"`
	ReferenceOutput map[string]interface{} `json:"reference_output"`
	ActualOutput    map[string]interface{} `json:"actual_output"`
	StartTime       time.Time              `json:"start_time"`
	EndTime         time.Time              `json:"end_time"`
	Metadata        map[string]interface{} `json:"metadata"`
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := langsmith.NewClient()
	ctx := context.Background()

	datasetName := "Experiment Dataset - Go Example"
	experimentName := fmt.Sprintf("My First Experiment - %s", time.Now().Format(time.RFC3339))

	printHeader()

	fmt.Println("Defining experiment results...")
	experimentResults := defineExperimentResults()
	fmt.Printf("   ✓ Defined %d experiment results\n\n", len(experimentResults))

	fmt.Printf("1. Getting or creating dataset '%s'...\n", datasetName)
	dataset, err := getOrCreateDataset(ctx, client, datasetName)
	if err != nil {
		return fmt.Errorf("getting or creating dataset: %w", err)
	}
	fmt.Printf("   → View dataset: %s\n\n", buildDatasetURL(dataset.TenantID, dataset.ID))

	fmt.Println("2. Creating dataset examples from experiment results...")
	createdExamples, err := createExamplesFromResults(ctx, client, dataset.ID, experimentResults)
	if err != nil {
		return fmt.Errorf("creating examples: %w", err)
	}
	fmt.Printf("   ✓ Created/updated %d examples with deterministic IDs\n\n", len(createdExamples))

	fmt.Println("3. Creating experiment session...")
	session, err := createExperimentSession(ctx, client, experimentName, dataset.ID)
	if err != nil {
		return fmt.Errorf("creating session: %w", err)
	}
	fmt.Printf("   ✓ Created experiment session: %s\n", session.Name)
	fmt.Printf("   ✓ Session ID: %s\n", session.ID)
	fmt.Printf("   → View experiment: %s\n\n", buildSessionURL(session.TenantID, session.ID))

	fmt.Println("4. Logging experiment runs to LangSmith...")
	if err := createAndIngestRuns(ctx, client, experimentResults, createdExamples, experimentName, session.ID); err != nil {
		return fmt.Errorf("ingesting runs: %w", err)
	}
	fmt.Printf("   ✓ Successfully logged %d runs to experiment session\n\n", len(experimentResults))

	printSummary(summaryParams{
		datasetName:    datasetName,
		experimentName: experimentName,
		exampleCount:   len(createdExamples),
		runCount:       len(experimentResults),
		datasetID:      dataset.ID,
		session:        session,
	})
	return nil
}

// summaryParams holds parameters for printing the experiment summary.
type summaryParams struct {
	datasetName    string
	experimentName string
	exampleCount   int
	runCount       int
	datasetID      string
	session        *langsmith.TracerSessionWithoutVirtualFields
}

// printHeader prints the example header.
func printHeader() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("LangSmith Experiment Recording Example")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
}

// defineExperimentResults returns sample experiment results for demonstration.
// In a real scenario, these would come from running your LLM application.
func defineExperimentResults() []ExperimentResult {
	now := time.Now()
	return []ExperimentResult{
		{
			Input: map[string]interface{}{
				"question": "What is the capital of France?",
			},
			ReferenceOutput: map[string]interface{}{
				"answer": "Paris",
			},
			ActualOutput: map[string]interface{}{
				"answer": "The capital of France is Marseille.",
			},
			StartTime: now.Add(-10 * time.Second),
			EndTime:   now.Add(-8 * time.Second),
			Metadata: map[string]interface{}{
				"model":       "gpt-4",
				"temperature": 0.7,
				"tokens":      25,
				"cost":        0.0015,
			},
		},
		{
			Input: map[string]interface{}{
				"question": "What is 2 + 2?",
			},
			ReferenceOutput: map[string]interface{}{
				"answer": "4",
			},
			ActualOutput: map[string]interface{}{
				"answer": "2 + 2 equals 4",
			},
			StartTime: now.Add(-8 * time.Second),
			EndTime:   now.Add(-7 * time.Second),
			Metadata: map[string]interface{}{
				"model":       "gpt-4",
				"temperature": 0.7,
				"tokens":      18,
				"cost":        0.001,
			},
		},
		{
			Input: map[string]interface{}{
				"question": "Who wrote Romeo and Juliet?",
			},
			ReferenceOutput: map[string]interface{}{
				"answer": "William Shakespeare",
			},
			ActualOutput: map[string]interface{}{
				"answer": "Romeo and Juliet was written by William Shakespeare.",
			},
			StartTime: now.Add(-7 * time.Second),
			EndTime:   now.Add(-5 * time.Second),
			Metadata: map[string]interface{}{
				"model":       "gpt-4",
				"temperature": 0.7,
				"tokens":      22,
				"cost":        0.0013,
			},
		},
		{
			Input: map[string]interface{}{
				"question": "What is the largest planet in our solar system?",
			},
			ReferenceOutput: map[string]interface{}{
				"answer": "Jupiter",
			},
			ActualOutput: map[string]interface{}{
				"answer": "Jupiter is the largest planet in our solar system.",
			},
			StartTime: now.Add(-5 * time.Second),
			EndTime:   now.Add(-3 * time.Second),
			Metadata: map[string]interface{}{
				"model":       "gpt-4",
				"temperature": 0.7,
				"tokens":      20,
				"cost":        0.0012,
			},
		},
	}
}

// getOrCreateDataset gets an existing dataset by name or creates a new one.
func getOrCreateDataset(ctx context.Context, client *langsmith.Client, datasetName string) (*langsmith.Dataset, error) {
	datasets, err := client.Datasets.List(ctx, langsmith.DatasetListParams{
		Name: langsmith.F(datasetName),
	})
	if err != nil {
		return nil, fmt.Errorf("listing datasets: %w", err)
	}

	if len(datasets.Items) > 0 {
		existing := &datasets.Items[0]
		fmt.Printf("   ✓ Found existing dataset (ID: %s)\n", existing.ID)
		return existing, nil
	}

	createParams := langsmith.DatasetNewParams{
		Name:        langsmith.F(datasetName),
		Description: langsmith.F("Dataset for recording experiment results. Inspired by langsmith-rest-experiments."),
	}
	created, err := client.Datasets.New(ctx, createParams)
	if err != nil {
		return nil, fmt.Errorf("creating dataset: %w", err)
	}
	fmt.Printf("   ✓ Created new dataset (ID: %s)\n", created.ID)
	return created, nil
}

// generateExampleID generates a deterministic UUID v5 based on dataset ID, input, and reference output.
// This ensures idempotent example creation. UUID v5 uses SHA-1 hashing with a namespace UUID.
func generateExampleID(datasetID string, input, referenceOutput map[string]interface{}) (string, error) {
	inputJSON, err := json.Marshal(input)
	if err != nil {
		return "", fmt.Errorf("marshaling input: %w", err)
	}

	outputJSON, err := json.Marshal(referenceOutput)
	if err != nil {
		return "", fmt.Errorf("marshaling reference output: %w", err)
	}

	namespace, err := uuid.Parse(datasetID)
	if err != nil {
		return "", fmt.Errorf("parsing dataset ID as UUID: %w", err)
	}

	name := fmt.Sprintf("%s|%s", string(inputJSON), string(outputJSON))
	id := uuid.NewSHA1(namespace, []byte(name))

	return id.String(), nil
}


// createExamplesFromResults creates examples from experiment results with deterministic IDs.
func createExamplesFromResults(ctx context.Context, client *langsmith.Client, datasetID string, results []ExperimentResult) ([]langsmith.Example, error) {
	exampleBodies := make([]langsmith.ExampleBulkNewParamsBody, 0, len(results))
	for _, result := range results {
		exampleID, err := generateExampleID(datasetID, result.Input, result.ReferenceOutput)
		if err != nil {
			return nil, fmt.Errorf("generating example ID: %w", err)
		}

		exampleBodies = append(exampleBodies, langsmith.ExampleBulkNewParamsBody{
			ID:        langsmith.F(exampleID),
			DatasetID: langsmith.F(datasetID),
			Inputs:    langsmith.F(result.Input),
			Outputs:   langsmith.F(result.ReferenceOutput),
		})
	}

	bulkParams := langsmith.ExampleBulkNewParams{
		Body: exampleBodies,
	}

	createdExamples, err := client.Examples.Bulk.New(ctx, bulkParams)
	if err != nil {
		return nil, fmt.Errorf("bulk creating examples: %w", err)
	}

	if createdExamples == nil {
		return nil, fmt.Errorf("bulk creating examples: received nil response")
	}

	return *createdExamples, nil
}

// createExperimentSession creates a new experiment session linked to a reference dataset.
func createExperimentSession(ctx context.Context, client *langsmith.Client, experimentName, datasetID string) (*langsmith.TracerSessionWithoutVirtualFields, error) {
	sessionParams := langsmith.SessionNewParams{
		Name:               langsmith.F(experimentName),
		Description:        langsmith.F("Experiment session for testing LLM responses. Use LLM-as-judge evaluators to compare actual_output vs reference_output."),
		StartTime:          langsmith.F(time.Now()),
		ReferenceDatasetID: langsmith.F(datasetID),
	}

	session, err := client.Sessions.New(ctx, sessionParams)
	if err != nil {
		return nil, fmt.Errorf("creating session: %w", err)
	}

	return session, nil
}

// generateUUID generates a random UUID v4 using the standard library.
func generateUUID() (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", fmt.Errorf("generating UUID: %w", err)
	}
	return id.String(), nil
}

// createAndIngestRuns creates runs from experiment results and ingests them in a batch.
func createAndIngestRuns(ctx context.Context, client *langsmith.Client, results []ExperimentResult, examples []langsmith.Example, experimentName, sessionID string) error {
	if len(results) != len(examples) {
		return fmt.Errorf("mismatched results and examples: %d results, %d examples", len(results), len(examples))
	}

	runs := make([]langsmith.RunParam, 0, len(results))
	for i, result := range results {
		runID, err := generateUUID()
		if err != nil {
			return fmt.Errorf("generating run ID: %w", err)
		}

		dottedOrder := formatDottedOrder(result.StartTime, runID)

		run := langsmith.RunParam{
			ID:                 langsmith.F(runID),
			Name:               langsmith.F(fmt.Sprintf("Experiment Run - Question %d", i+1)),
			RunType:            langsmith.F(langsmith.RunRunTypeLlm),
			Inputs:             langsmith.F(result.Input),
			Outputs:            langsmith.F(result.ActualOutput),
			ReferenceExampleID: langsmith.F(examples[i].ID),
			SessionName:        langsmith.F(experimentName),
			SessionID:          langsmith.F(sessionID),
			StartTime:          langsmith.F(result.StartTime.Format(time.RFC3339)),
			EndTime:            langsmith.F(result.EndTime.Format(time.RFC3339)),
			TraceID:            langsmith.F(runID),
			DottedOrder:        langsmith.F(dottedOrder),
			Extra:              langsmith.F(result.Metadata),
		}
		runs = append(runs, run)
	}

	ingestParams := langsmith.RunIngestBatchParams{
		Post: langsmith.F(runs),
	}

	if _, err := client.Runs.IngestBatch(ctx, ingestParams); err != nil {
		return fmt.Errorf("ingesting batch: %w", err)
	}

	return nil
}

// formatDottedOrder formats a timestamp and UUID into a dotted order string for root runs.
// Format: YYYYMMDDTHHMMSSmmmmmmZ{run_id} where mmmmmm is microseconds (6 digits).
// For root runs: single part (no dots): {timestamp}Z{run_id}
// For child runs: {parent_dotted_order}.{timestamp}Z{run_id}
// Rules:
//   - Root runs: single part (no dots)
//   - First part must match trace_id (root run ID)
//   - Last part must match the current run_id
//   - Each part: {timestamp}Z{run_id} where timestamp is YYYYMMDDTHHMMSSmmmmmm
func formatDottedOrder(t time.Time, runID string) string {
	// Format: YYYYMMDDTHHMMSS (date and time)
	//         mmmmmm (microseconds, 6 digits)
	//         Z (literal Z)
	//         {run_id} (UUID)
	timestamp := t.Format("20060102T150405")
	microseconds := t.Nanosecond() / 1000 // Convert nanoseconds to microseconds
	return fmt.Sprintf("%s%06dZ%s", timestamp, microseconds, runID)
}

// buildDatasetURL builds a URL to view a dataset in LangSmith.
func buildDatasetURL(tenantID, datasetID string) string {
	return fmt.Sprintf("https://smith.langchain.com/o/%s/datasets/%s", tenantID, datasetID)
}

// buildSessionURL builds a URL to view a session/experiment in LangSmith.
// The session ID is used as the project ID in the URL.
func buildSessionURL(tenantID, sessionID string) string {
	return fmt.Sprintf("https://smith.langchain.com/o/%s/projects/p/%s", tenantID, sessionID)
}

// printSummary prints a summary of the experiment and next steps.
func printSummary(p summaryParams) {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Experiment Complete!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
	fmt.Println("What was logged:")
	fmt.Printf("  • Dataset: %s with %d examples\n", p.datasetName, p.exampleCount)
	fmt.Printf("  • Experiment: %s with %d runs\n", p.experimentName, p.runCount)
	fmt.Println("  • Each run has:")
	fmt.Println("    - Input (the question)")
	fmt.Println("    - Reference output (ground truth answer)")
	fmt.Println("    - Actual output (simulated LLM response)")
	fmt.Println("    - Metadata (model, tokens, cost, etc.)")
	fmt.Println()
	fmt.Println("Next steps:")
	fmt.Println("  1. View your experiment results:")
	fmt.Printf("     %s\n", buildSessionURL(p.session.TenantID, p.session.ID))
	fmt.Println()
	fmt.Println("  2. In the LangSmith UI, you can:")
	fmt.Println("     - Compare actual outputs vs reference outputs side-by-side")
	fmt.Println("     - Add LLM-as-judge evaluators to automatically score quality")
	fmt.Println("     - Create annotation queues for human feedback")
	fmt.Println("     - Run additional experiments and compare results")
	fmt.Println()
	fmt.Println("  3. For your real application:")
	fmt.Println("     - Replace experiment_results with actual LLM application outputs")
	fmt.Println("     - Log runs as they complete using client.Runs.IngestBatch()")
	fmt.Println("     - Use the same session to group related runs")
	fmt.Println()
	fmt.Printf("Dataset URL:    %s\n", buildDatasetURL(p.session.TenantID, p.datasetID))
	fmt.Printf("Experiment URL: %s\n", buildSessionURL(p.session.TenantID, p.session.ID))
	fmt.Println()
	fmt.Println("Reference: https://github.com/jacoblee93/langsmith-rest-experiments")
}
