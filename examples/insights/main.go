package main

import (
	"context"
	"errors"
	"fmt"
	"os"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to create a one-time Insights report for a LangSmith project.
//
// Running this example creates a real report and may incur model usage costs.
// It does not configure a recurring schedule.
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT_ID: The UUID of the project to analyze
//
// Running:
//
//	go run ./examples/insights

const reportName = "Go SDK Insights example"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	projectID := os.Getenv("LANGSMITH_PROJECT_ID")
	if projectID == "" {
		return errors.New("LANGSMITH_PROJECT_ID is required")
	}

	client := langsmith.NewClient()
	config, err := client.Sessions.Insights.Configs.New(
		context.Background(),
		projectID,
		buildInsightsParams(),
	)
	if err != nil {
		return fmt.Errorf("creating Insights report: %w", err)
	}

	fmt.Printf("Created Insights report configuration %s\n", config.ID)
	return nil
}

func buildInsightsParams() langsmith.SessionInsightConfigNewParams {
	return langsmith.SessionInsightConfigNewParams{
		Name: langsmith.F(reportName),
		Config: langsmith.F(langsmith.CreateRunClusteringJobRequestParam{
			Name:          langsmith.F(reportName),
			Model:         langsmith.F(langsmith.CreateRunClusteringJobRequestModelOpenAI),
			LastNHours:    langsmith.Int(24),
			Filter:        langsmith.F("eq(is_root, true)"),
			SummaryPrompt: langsmith.F("Summarize the user's intent and any failures."),
		}),
	}
}
