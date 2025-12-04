package main

import (
	"bufio"
	"context"
	"fmt"
	"os"
	"strings"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to create a dataset and add examples to it programmatically.
//
// This example shows:
//   - Checking if a dataset already exists by name
//   - Creating a new dataset
//   - Adding individual examples with inputs and outputs
//   - Adding examples in bulk
//   - Retrieving the dataset and showing example count
//   - Deleting the dataset
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//
// Running:
//
//	go run ./examples/dataset

const (
	datasetName        = "Sample Dataset created in Go"
	datasetDescription = "A sample dataset in LangSmith."
	defaultEndpoint    = "https://api.smith.langchain.com"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := langsmith.NewClient()
	ctx := context.Background()

	// 1. Check if dataset exists and delete if it does
	if err := checkAndDeleteExistingDataset(ctx, client); err != nil {
		return err
	}

	// 2. Create a new dataset
	dataset, err := createDataset(ctx, client)
	if err != nil {
		return err
	}

	// 3. Add individual examples
	if err := createIndividualExamples(ctx, client, dataset.ID); err != nil {
		return err
	}

	// 4. Add examples in bulk
	if err := createBulkExamples(ctx, client, dataset.ID); err != nil {
		return err
	}

	// 5. Retrieve and display dataset info
	if err := retrieveAndDisplayDataset(ctx, client, dataset.ID); err != nil {
		return err
	}

	// 6. Delete the dataset
	fmt.Println("\n6. Deleting the dataset...")
	return deleteDatasetWithConfirmation(ctx, client, dataset)
}

// exampleData represents a question-answer pair for dataset examples.
type exampleData struct {
	question string
	answer   string
}

// getIndividualExamples returns the examples to create individually.
func getIndividualExamples() []exampleData {
	return []exampleData{
		{"Which country is Mount Kilimanjaro located in?", "Mount Kilimanjaro is located in Tanzania."},
		{"What is Earth's lowest point?", "Earth's lowest point is The Dead Sea."},
	}
}

// getBulkExamples returns the examples to create in bulk.
func getBulkExamples() []exampleData {
	return []exampleData{
		{"What is the capital of France?", "The capital of France is Paris."},
		{"Which ocean is the largest?", "The Pacific Ocean is the largest."},
	}
}

// checkAndDeleteExistingDataset checks if a dataset exists and deletes it if found.
func checkAndDeleteExistingDataset(ctx context.Context, client *langsmith.Client) error {
	fmt.Println("1. Listing datasets with name filter using client.Datasets.List()...")
	listParams := langsmith.DatasetListParams{
		Name: langsmith.F(datasetName),
	}

	existingDatasets, err := client.Datasets.List(ctx, listParams)
	if err != nil {
		return fmt.Errorf("listing datasets: %w", err)
	}

	if len(existingDatasets.Items) > 0 {
		existing := existingDatasets.Items[0]
		fmt.Fprintf(os.Stderr, "Error: Dataset with name '%s' already exists\n", datasetName)
		fmt.Fprintf(os.Stderr, "View existing dataset: %s\n", buildDatasetURL(&existing))
		return deleteDatasetWithConfirmation(ctx, client, &existing)
	}

	fmt.Printf("   ✓ No existing dataset found with name '%s'\n\n", datasetName)
	return nil
}

// createDataset creates a new dataset.
func createDataset(ctx context.Context, client *langsmith.Client) (*langsmith.Dataset, error) {
	fmt.Println("2. Creating dataset using client.Datasets.New()...")
	datasetParams := langsmith.DatasetNewParams{
		Name:        langsmith.F(datasetName),
		Description: langsmith.F(datasetDescription),
	}

	dataset, err := client.Datasets.New(ctx, datasetParams)
	if err != nil {
		return nil, fmt.Errorf("creating dataset: %w", err)
	}
	fmt.Printf("   ✓ Created dataset: %s (ID: %s)\n\n", dataset.Name, dataset.ID)
	return dataset, nil
}

// createIndividualExamples creates individual examples.
func createIndividualExamples(ctx context.Context, client *langsmith.Client, datasetID string) error {
	fmt.Println("3. Adding examples using client.Examples.New()...")

	examples := getIndividualExamples()
	for i, ex := range examples {
		params := buildExampleParams(datasetID, ex)
		example, err := client.Examples.New(ctx, params)
		if err != nil {
			return fmt.Errorf("creating example %d: %w", i+1, err)
		}
		fmt.Printf("   ✓ Created example %d: %s\n", i+1, example.ID)
	}
	fmt.Println()
	return nil
}

// createBulkExamples creates examples in bulk.
func createBulkExamples(ctx context.Context, client *langsmith.Client, datasetID string) error {
	fmt.Println("4. Adding examples using client.Examples.Bulk.New()...")

	examples := getBulkExamples()
	bulkParams := langsmith.ExampleBulkNewParams{
		Body: make([]langsmith.ExampleBulkNewParamsBody, 0, len(examples)),
	}

	for _, ex := range examples {
		bulkParams.Body = append(bulkParams.Body, buildBulkExampleBody(datasetID, ex))
	}

	results, err := client.Examples.Bulk.New(ctx, bulkParams)
	if err != nil {
		return fmt.Errorf("creating bulk examples: %w", err)
	}
	fmt.Printf("   ✓ Created %d examples in bulk\n\n", len(*results))
	return nil
}

// buildExampleParams builds ExampleNewParams from example data.
func buildExampleParams(datasetID string, ex exampleData) langsmith.ExampleNewParams {
	return langsmith.ExampleNewParams{
		DatasetID: langsmith.F(datasetID),
		Inputs: langsmith.F(map[string]interface{}{
			"question": ex.question,
		}),
		Outputs: langsmith.F(map[string]interface{}{
			"answer": ex.answer,
		}),
	}
}

// buildBulkExampleBody builds ExampleBulkNewParamsBody from example data.
func buildBulkExampleBody(datasetID string, ex exampleData) langsmith.ExampleBulkNewParamsBody {
	return langsmith.ExampleBulkNewParamsBody{
		DatasetID: langsmith.F(datasetID),
		Inputs: langsmith.F(map[string]interface{}{
			"question": ex.question,
		}),
		Outputs: langsmith.F(map[string]interface{}{
			"answer": ex.answer,
		}),
	}
}

// retrieveAndDisplayDataset retrieves and displays dataset information.
func retrieveAndDisplayDataset(ctx context.Context, client *langsmith.Client, datasetID string) error {
	fmt.Println("5. Retrieving dataset using client.Datasets.Get()...")
	dataset, err := client.Datasets.Get(ctx, datasetID)
	if err != nil {
		return fmt.Errorf("retrieving dataset: %w", err)
	}
	fmt.Printf("   ✓ Retrieved dataset: %s\n", dataset.Name)
	fmt.Printf("   ✓ Number of examples: %d\n", dataset.ExampleCount)
	fmt.Printf("   ✓ View dataset: %s\n", buildDatasetURL(dataset))
	return nil
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

// deleteDatasetWithConfirmation deletes a dataset after user confirmation.
func deleteDatasetWithConfirmation(ctx context.Context, client *langsmith.Client, dataset *langsmith.Dataset) error {
	fmt.Println("Press Enter to delete the dataset using client.Datasets.Delete()...")
	reader := bufio.NewReader(os.Stdin)
	_, err := reader.ReadString('\n')
	if err != nil {
		// If stdin is not available (e.g., running in CI), skip confirmation
		fmt.Println("   (Skipping confirmation - stdin not available)")
	}

	_, err = client.Datasets.Delete(ctx, dataset.ID)
	if err != nil {
		return fmt.Errorf("deleting dataset: %w", err)
	}
	fmt.Printf("   ✓ Deleted dataset: %s (ID: %s)\n", dataset.Name, dataset.ID)
	return nil
}

