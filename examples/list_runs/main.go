package main

import (
	"context"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to query and list runs from a LangSmith project.
//
// This example shows:
//   - Resolving a project by name with an upsert create
//   - Configuring V2 run query parameters with a project and a time window
//   - Selecting which run fields to return
//   - Paging through runs with Runs.QueryV2AutoPaging
//   - Iterating through and displaying run information
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_PROJECT: The project name to query runs from (defaults to "default")
//
// Running:
//
//	go run ./examples/list_runs
const (
	// defaultProjectName is the project queried when LANGSMITH_PROJECT is unset.
	defaultProjectName = "default"
	// defaultPageSize is how many runs to ask for in a single page. The auto-pager
	// requests another page whenever the current one is exhausted.
	defaultPageSize = 10
	// maxRuns bounds how many runs this example prints, so a busy project does not
	// page indefinitely. It must stay above defaultPageSize for paging to happen.
	maxRuns = 25
	// defaultLookback bounds the query to recently started runs. QueryV2 defaults
	// to the last day when min_start_time is omitted.
	defaultLookback = 7 * 24 * time.Hour
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

	fmt.Println("=== LangSmith List Runs Example ===")
	fmt.Println()

	// Resolve the project to query. Creating with upsert=true returns the existing
	// project when one already has this name, so the example always ends up with a
	// project ID without the reader having to look one up in the UI first.
	projectName := os.Getenv("LANGSMITH_PROJECT")
	if projectName == "" {
		projectName = defaultProjectName
	}
	project, err := client.Sessions.New(ctx, buildProjectParams(projectName))
	if err != nil {
		return fmt.Errorf("resolving project %q: %w", projectName, err)
	}
	fmt.Printf("Project: %s (%s)\n\n", project.Name, project.ID)

	// Query runs with the V2 endpoint (POST /api/v2/runs/query). A query must be
	// scoped to either project_ids or reference_dataset_id. The auto-pager follows
	// next_cursor for us, so we iterate runs rather than pages.
	pager := client.Runs.QueryV2AutoPaging(ctx, buildQueryParams(project.ID, time.Now()))

	count := 0
	for pager.Next() {
		count++
		fmt.Print(formatRun(count, pager.Current()))
		if count == maxRuns {
			fmt.Printf("Stopping at %d run(s); raise maxRuns to page through more.\n\n", maxRuns)
			break
		}
	}
	// Paging errors surface here rather than from Next.
	if err := pager.Err(); err != nil {
		return fmt.Errorf("querying runs: %w", err)
	}

	if count == 0 {
		fmt.Println("No runs found for this project in the last 7 days.")
		return nil
	}

	fmt.Printf("Listed %d run(s).\n", count)
	return nil
}

// buildProjectParams builds the create call used to resolve a project by name.
// upsert=true makes the create idempotent: the API returns the existing project
// instead of failing when the name is already taken.
func buildProjectParams(projectName string) langsmith.SessionNewParams {
	return langsmith.SessionNewParams{
		Upsert: langsmith.F(true),
		Name:   langsmith.F(projectName),
	}
}

// runSelectFields lists the run properties this example prints. QueryV2 returns
// only `id` on each run unless the fields are requested explicitly.
func runSelectFields() []langsmith.RunSelectField {
	return []langsmith.RunSelectField{
		langsmith.RunSelectFieldName,
		langsmith.RunSelectFieldRunType,
		langsmith.RunSelectFieldStatus,
		langsmith.RunSelectFieldStartTime,
		langsmith.RunSelectFieldLatencySeconds,
	}
}

// buildQueryParams builds the V2 query for a project's recent runs.
func buildQueryParams(projectID string, now time.Time) langsmith.RunQueryV2Params {
	return langsmith.RunQueryV2Params{
		ProjectIDs:   langsmith.F([]string{projectID}),
		MinStartTime: langsmith.F(now.Add(-defaultLookback)),
		PageSize:     langsmith.F(int64(defaultPageSize)),
		Selects:      langsmith.F(runSelectFields()),
	}
}

// formatRun renders one run as the display block printed by this example.
func formatRun(position int, run langsmith.Run) string {
	var b strings.Builder
	fmt.Fprintf(&b, "%d. Run ID: %s\n", position, run.ID)
	fmt.Fprintf(&b, "   Name: %s\n", run.Name)
	if run.RunType != "" {
		fmt.Fprintf(&b, "   Type: %s\n", run.RunType)
	}
	if run.Status != "" {
		fmt.Fprintf(&b, "   Status: %s\n", run.Status)
	}
	if !run.StartTime.IsZero() {
		fmt.Fprintf(&b, "   Started: %s\n", run.StartTime.Format(time.RFC3339))
	}
	if run.LatencySeconds > 0 {
		fmt.Fprintf(&b, "   Latency: %.2fs\n", run.LatencySeconds)
	}
	b.WriteString("\n")
	return b.String()
}
