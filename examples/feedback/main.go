package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/shared"
)

// Demonstrates how to manage feedback in LangSmith.
//
// This example shows:
//   - Creating a run to provide feedback on
//   - Submitting feedback (score, comment, correction) for a run
//   - Updating existing feedback
//   - Listing feedback with filters
//   - Deleting feedback
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//
// Running:
//
//	go run ./examples/feedback
func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	client := langsmith.NewClient()
	ctx := context.Background()

	printHeader()

	// 1. Create a dummy run to provide feedback on
	// In a real scenario, you would have a run ID from an actual LLM call.
	runID := uuid.New().String()
	fmt.Printf("1. Creating a dummy run (ID: %s)...\n", runID)
	
	// We use IngestBatch to create a run quickly for this example
	_, err := client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Post: langsmith.F([]langsmith.RunParam{
			{
				ID:        langsmith.F(runID),
				Name:      langsmith.F("Feedback Example Run"),
				RunType:   langsmith.F(langsmith.RunRunTypeLlm),
				StartTime: langsmith.F(time.Now().Format(time.RFC3339)),
				Inputs:    langsmith.F(map[string]interface{}{"question": "What is 2+2?"}),
				Outputs:   langsmith.F(map[string]interface{}{"answer": "4"}),
			},
		}),
	})
	if err != nil {
		return fmt.Errorf("creating dummy run: %w", err)
	}
	fmt.Println("   ✓ Run created")

	// 2. Submit feedback for the run
	fmt.Println("\n2. Submitting feedback for the run...")
	feedback, err := client.Feedback.New(ctx, langsmith.FeedbackNewParams{
		FeedbackCreateSchema: langsmith.FeedbackCreateSchemaParam{
			RunID:   langsmith.F(runID),
			Key:     langsmith.F("user_score"),
			Score:   langsmith.F[langsmith.FeedbackCreateSchemaScoreUnionParam](shared.UnionFloat(1.0)),
			Comment: langsmith.F("Excellent answer!"),
			Value:   langsmith.F[langsmith.FeedbackCreateSchemaValueUnionParam](shared.UnionString("correct")),
		},
	})
	if err != nil {
		return fmt.Errorf("submitting feedback: %w", err)
	}
	fmt.Printf("   ✓ Feedback submitted (ID: %s)\n", feedback.ID)

	// 3. Update the feedback
	fmt.Println("\n3. Updating the feedback comment...")
	updatedFeedback, err := client.Feedback.Update(ctx, feedback.ID, langsmith.FeedbackUpdateParams{
		Comment: langsmith.F("Excellent answer! Very concise."),
	})
	if err != nil {
		return fmt.Errorf("updating feedback: %w", err)
	}
	fmt.Printf("   ✓ Feedback updated. New comment: %s\n", updatedFeedback.Comment)

	// 4. List feedback for the run
	fmt.Println("\n4. Listing feedback for the run...")
	listRes, err := client.Feedback.List(ctx, langsmith.FeedbackListParams{
		Run: langsmith.F[langsmith.FeedbackListParamsRunUnion](langsmith.FeedbackListParamsRunArray([]string{runID})),
	})
	if err != nil {
		return fmt.Errorf("listing feedback: %w", err)
	}
	fmt.Printf("   ✓ Found %d feedback entries for run %s\n", len(listRes.Items), runID)
	for _, f := range listRes.Items {
		fmt.Printf("     - Key: %s, Score: %v, Comment: %s\n", f.Key, f.Score, f.Comment)
	}

	// 5. Delete the feedback
	fmt.Println("\n5. Deleting the feedback...")
	_, err = client.Feedback.Delete(ctx, feedback.ID)
	if err != nil {
		return fmt.Errorf("deleting feedback: %w", err)
	}
	fmt.Println("   ✓ Feedback deleted")

	printSummary()
	return nil
}

func printHeader() {
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("LangSmith Feedback Management Example")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println()
}

func printSummary() {
	fmt.Println()
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("Example Complete!")
	fmt.Println(strings.Repeat("=", 60))
	fmt.Println("\nIn this example, we demonstrated how to:")
	fmt.Println("  1. Create a run in LangSmith")
	fmt.Println("  2. Submit programmatic feedback (scores and comments)")
	fmt.Println("  3. Update existing feedback entries")
	fmt.Println("  4. Query feedback by run ID")
	fmt.Println("  5. Clean up by deleting feedback")
	fmt.Println("\nFeedback is essential for improving LLM applications by:")
	fmt.Println("  - Collecting human-in-the-loop ratings")
	fmt.Println("  - Recording automated evaluator results")
	fmt.Println("  - Tracking performance over time in your LangSmith projects")
}
