package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/langchain-ai/langsmith-go/shared"
	"github.com/spf13/cobra"
)

func newExampleCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "examples",
		Aliases: []string{"example"},
		Short:   "Manage individual examples within datasets",
		Long: `Manage individual examples within datasets.

Examples are the individual input/output pairs stored in a dataset.
Use these commands to list, add, or remove examples.

List results are paginated and return at most 20 examples by default
(use --limit to change). Use --offset to paginate through results.

Examples:
  langsmith examples list --dataset my-dataset
  langsmith examples create --dataset my-dataset --inputs '{"question": "What is LangSmith?"}'
  langsmith examples delete <example-id> --yes`,
	}

	cmd.AddCommand(newExampleListCmd())
	cmd.AddCommand(newExampleCreateCmd())
	cmd.AddCommand(newExampleDeleteCmd())
	return cmd
}

func newExampleListCmd() *cobra.Command {
	var (
		datasetName string
		limit       int
		offset      int
		split       string
		outputFile  string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List examples in a dataset (default: 20)",
		Run: func(cmd *cobra.Command, args []string) {
			c := MustGetClient()
			ctx := context.Background()

			ds, err := resolveDataset(ctx, c, datasetName)
			if err != nil {
				ExitErrorf("%v", err)
			}

			pageSize := int64(20)
			if limit > 0 && int64(limit) < pageSize {
				pageSize = int64(limit)
			}
			params := exampleListParams(ds.ID, pageSize, offset, split)
			var examples []langsmith.Example
			pager := c.SDK.Examples.ListAutoPaging(ctx, params)
			for pager.Next() {
				examples = append(examples, pager.Current())
				if limit > 0 && len(examples) >= limit {
					break
				}
			}
			if err := pager.Err(); err != nil {
				ExitErrorf("listing examples: %v", err)
			}
			fmt_ := GetFormat()

			if fmt_ == "pretty" {
				columns := []string{"ID", "Split", "Created", "Inputs Preview"}
				var rows [][]string
				for _, ex := range examples {
					id := ex.ID
					splitVal := exampleSplitDisplay(ex.Metadata)
					created := "N/A"
					if !ex.CreatedAt.IsZero() {
						created = ex.CreatedAt.Format("2006-01-02")
					}
					inputsPreview := "N/A"
					if ex.Inputs != nil {
						b, _ := json.Marshal(ex.Inputs)
						inputsPreview = string(b)
						if len(inputsPreview) > 60 {
							inputsPreview = inputsPreview[:60] + "..."
						}
					}
					rows = append(rows, []string{id, splitVal, created, inputsPreview})
				}
				output.OutputTable(columns, rows, fmt.Sprintf("Examples in %s", ds.Name))
			} else {
				var data []map[string]any
				for _, ex := range examples {
					entry := map[string]any{
						"id":         ex.ID,
						"inputs":     ex.Inputs,
						"outputs":    ex.Outputs,
						"metadata":   ex.Metadata,
						"created_at": formatTimeISO(ex.CreatedAt),
					}
					data = append(data, entry)
				}
				if err := output.OutputJSON(data, outputFile); err != nil {
					ExitErrorf("%v", err)
				}
			}
		},
	}

	cmd.Flags().StringVar(&datasetName, "dataset", "", "Dataset name or UUID (required)")
	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum number of examples to return")
	cmd.Flags().IntVar(&offset, "offset", 0, "Number of examples to skip (pagination)")
	cmd.Flags().StringVar(&split, "split", "", "Filter by split (train, test, validation)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")
	_ = cmd.MarkFlagRequired("dataset")

	return cmd
}

func exampleSplitDisplay(metadata map[string]any) string {
	if metadata == nil {
		return "N/A"
	}
	switch splits := metadata["dataset_split"].(type) {
	case []any:
		values := make([]string, 0, len(splits))
		for _, split := range splits {
			if value, ok := split.(string); ok {
				values = append(values, value)
			}
		}
		if len(values) > 0 {
			return strings.Join(values, ", ")
		}
	case []string:
		if len(splits) > 0 {
			return strings.Join(splits, ", ")
		}
	case string:
		if splits != "" {
			return splits
		}
	}
	return "N/A"
}

func newExampleCreateCmd() *cobra.Command {
	var (
		datasetName    string
		inputs         string
		outputs        string
		metadata       string
		split          string
		fromRuns       []string
		skipDuplicates bool
		project        string
		projectID      string
	)

	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new example in a dataset, or one example per run with --from-run",
		Long: `Create a new example in a dataset.

Pass --inputs (and optionally --outputs) to create one example from JSON, or
--from-run to copy each run's inputs, outputs and attachments into a new
example that keeps the run as its source. --from-run is repeatable;
--from-run - reads run IDs from stdin, one per line. --skip-duplicates skips runs whose inputs already
match an example in the dataset (or an earlier run in the same call). For
chat and LLM datasets the server converts each run to the dataset's format.
Pass --project or --project-id to read the runs through the v2 runs API,
which deployments without v1 run reads require.

Examples:
  langsmith examples create --dataset my-dataset --inputs '{"question": "What is LangSmith?"}'
  langsmith examples create --dataset regressions --from-run <run-id> --skip-duplicates
  langsmith runs list --project my-app --error --format json --jq '.[].run_id' | langsmith examples create --dataset regressions --from-run - --skip-duplicates`,
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(fromRuns) > 0 {
				if inputs != "" || outputs != "" {
					return fmt.Errorf("--from-run copies inputs and outputs from the run; it can't be combined with --inputs or --outputs")
				}
			} else {
				if inputs == "" {
					return fmt.Errorf("one of --inputs or --from-run is required")
				}
				if skipDuplicates {
					return fmt.Errorf("--skip-duplicates requires --from-run")
				}
				if project != "" || projectID != "" {
					return fmt.Errorf("--project and --project-id require --from-run")
				}
			}

			var parsedMetadata map[string]any
			if metadata != "" {
				if err := json.Unmarshal([]byte(metadata), &parsedMetadata); err != nil {
					return fmt.Errorf("Invalid JSON for --metadata: %v", err)
				}
			}

			if len(fromRuns) > 0 {
				runIDs, err := collectRunIDs(fromRuns, cmd.InOrStdin())
				if err != nil {
					return err
				}
				c := MustGetClient()
				ctx := cmd.Context()
				if ctx == nil {
					ctx = context.Background()
				}
				var sessionID string
				if project != "" || projectID != "" {
					if sessionID, err = resolveSessionID(ctx, c, project, projectID, "examples create"); err != nil {
						return err
					}
				}
				return createExamplesFromRuns(ctx, c, datasetName, runIDs, sessionID, parsedMetadata, split, skipDuplicates)
			}

			// Parse JSON inputs
			var parsedInputs map[string]any
			if err := json.Unmarshal([]byte(inputs), &parsedInputs); err != nil {
				return fmt.Errorf("Invalid JSON for --inputs: %v", err)
			}

			var parsedOutputs map[string]any
			if outputs != "" {
				if err := json.Unmarshal([]byte(outputs), &parsedOutputs); err != nil {
					return fmt.Errorf("Invalid JSON for --outputs: %v", err)
				}
			}

			c := MustGetClient()
			ctx := context.Background()

			ds, err := resolveDataset(ctx, c, datasetName)
			if err != nil {
				ExitErrorf("%v", err)
			}

			params := exampleCreateParams(ds.ID, parsedInputs, parsedOutputs, parsedMetadata, split)

			ex, err := c.SDK.Examples.New(ctx, params)
			if err != nil {
				ExitErrorf("creating example: %v", err)
			}
			return output.OutputJSON(map[string]any{
				"status":     "created",
				"id":         ex.ID,
				"dataset_id": ex.DatasetID,
				"inputs":     ex.Inputs,
				"outputs":    ex.Outputs,
			}, "")
		},
	}

	cmd.Flags().StringVar(&datasetName, "dataset", "", "Dataset name or UUID (required)")
	cmd.Flags().StringVar(&inputs, "inputs", "", "JSON string of input fields (required unless --from-run is set)")
	cmd.Flags().StringVar(&outputs, "outputs", "", "JSON string of output fields")
	cmd.Flags().StringVar(&metadata, "metadata", "", "JSON string of metadata")
	cmd.Flags().StringVar(&split, "split", "", "Assign to a split (train, test, validation)")
	cmd.Flags().StringArrayVar(&fromRuns, "from-run", nil, "Run ID to copy into a new example (repeatable; - reads IDs from stdin, one per line)")
	cmd.Flags().BoolVar(&skipDuplicates, "skip-duplicates", false, "With --from-run, skip runs whose inputs already match an example in the dataset")
	cmd.Flags().StringVar(&project, "project", "", "With --from-run, name of the project the runs belong to; reads them through the v2 runs API")
	cmd.Flags().StringVar(&projectID, "project-id", "", "With --from-run, UUID of the project the runs belong to; reads them through the v2 runs API")
	cmd.MarkFlagsMutuallyExclusive("project", "project-id")
	_ = cmd.MarkFlagRequired("dataset")

	return cmd
}

func exampleListParams(datasetID string, pageSize int64, offset int, split string) langsmith.ExampleListParams {
	params := langsmith.ExampleListParams{Dataset: langsmith.F(datasetID), Limit: langsmith.F(pageSize)}
	if offset > 0 {
		params.Offset = langsmith.F(int64(offset))
	}
	if split != "" {
		params.Splits = langsmith.F([]string{split})
	}
	return params
}

func exampleCreateParams(datasetID string, inputs, outputs, metadata map[string]any, split string) langsmith.ExampleNewParams {
	params := langsmith.ExampleNewParams{DatasetID: langsmith.F(datasetID), Inputs: langsmith.F(inputs)}
	if outputs != nil {
		params.Outputs = langsmith.F(outputs) //nolint:staticcheck // older servers reject reference_outputs
	}
	if metadata != nil {
		params.Metadata = langsmith.F(metadata)
	}
	if split != "" {
		params.Split = langsmith.F[langsmith.ExampleNewParamsSplitUnion](shared.UnionString(split))
	}
	return params
}

func newExampleDeleteCmd() *cobra.Command {
	var yes bool

	cmd := &cobra.Command{
		Use:   "delete EXAMPLE_ID",
		Short: "Delete an example by its UUID",
		Args:  cobra.ExactArgs(1),
		RunE: errorsToStderr(func(cmd *cobra.Command, args []string) error {
			exampleID := args[0]

			if !yes {
				if err := confirmDelete(cmd, deleteConfirmation{
					acceptsUnterminated: true,
					target:              "the example",
					identity:            "Example: " + exampleID,
				}); err != nil {
					return err
				}
			}

			c := MustGetClient()
			ctx := context.Background()

			if _, err := c.SDK.Examples.Delete(ctx, exampleID); err != nil {
				return fmt.Errorf("deleting example: %w", err)
			}
			return output.OutputJSON(map[string]any{
				"status": "deleted",
				"id":     exampleID,
			}, "")
		}),
	}

	addYesFlag(cmd, &yes, "Skip confirmation prompt")
	return cmd
}
