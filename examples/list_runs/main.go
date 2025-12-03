package main

import (
	"context"
	"fmt"
	"os"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to query and list runs from a LangSmith project.
//
// This example shows:
//   - Configuring run query parameters with filters
//   - Querying runs by project/session ID
//   - Iterating through and displaying run information
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT_ID: The project ID to query runs from (find in LangSmith UI → Settings)
//
// Running:
//
//	go run ./examples/list_runs

const defaultLimit = 10

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := langsmith.NewClient()
	ctx := context.Background()

	projectID := os.Getenv("LANGSMITH_PROJECT_ID")
	if projectID == "" {
		return fmt.Errorf("LANGSMITH_PROJECT_ID environment variable is required.\nYou can find your project ID in the LangSmith UI under Settings")
	}

	fmt.Println("=== LangSmith List Runs Example ===")
	fmt.Println()

	// Create query parameters
	// The API requires at least one filter (session, id, parent_run, trace, or reference_example)
	// This example queries runs for a specific project/session
	params := langsmith.RunQueryParams{
		Session: langsmith.F([]string{projectID}),
		Limit:   langsmith.F(int64(defaultLimit)),
	}

	// Query runs
	response, err := client.Runs.Query(ctx, params)
	if err != nil {
		return fmt.Errorf("querying runs: %w", err)
	}

	// Print runs
	fmt.Printf("Found %d run(s):\n", len(response.Runs))
	if len(response.Runs) == 0 {
		fmt.Println("No runs found for this project.")
		return nil
	}

	fmt.Println()
	for i, run := range response.Runs {
		fmt.Printf("%d. Run ID: %s\n", i+1, run.ID)
		fmt.Printf("   Name: %s\n", run.Name)
		if run.RunType != "" {
			fmt.Printf("   Type: %s\n", run.RunType)
		}
		if run.Status != "" {
			fmt.Printf("   Status: %s\n", run.Status)
		}
		fmt.Println()
	}

	return nil
}
