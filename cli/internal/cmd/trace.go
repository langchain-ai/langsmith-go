package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/google/uuid"
	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

func newTraceCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "traces",
		Aliases: []string{"trace"},
		Short:   "Query and export traces (top-level agent runs and their full hierarchy)",
		Long: `Query and export traces (top-level agent runs and their full hierarchy).

A trace is a tree of runs representing one end-to-end invocation of your
application. The root run is the top-level entry point; child runs are
LLM calls, tool calls, retriever steps, etc.

Results are sorted newest-first by start time.

Examples:
  langsmith traces list --project my-app --limit 5
  langsmith traces list --project my-app --filter 'and(eq(input_key, "in_reply_to"), eq(input_value, "T-1001"))' --include-io
  langsmith traces get <trace-id> --full
  langsmith traces export ./traces --project my-app --limit 20 --full`,
	}

	cmd.AddCommand(newTraceListCmd())
	cmd.AddCommand(newTraceGetCmd())
	cmd.AddCommand(newTraceExportCmd())
	cmd.AddCommand(newTraceMessagesCmd())
	cmd.AddCommand(newTraceStatsCmd())
	cmd.AddCommand(newTraceSetupCmd())
	return cmd
}

func newTraceListCmd() *cobra.Command {
	var (
		ff              FilterFlags
		includeMetadata bool
		includeIO       bool
		includeFeedback bool
		includeFlagged  bool
		full            bool
		showHierarchy   bool
		outputFile      string
	)

	cmd := &cobra.Command{
		Use:   "list",
		Short: "List traces (root runs) matching filter criteria (default: 20, newest first)",
		Run: func(cmd *cobra.Command, args []string) {
			if full {
				includeMetadata = true
				includeIO = true
				includeFeedback = true
			}

			defaultLimit := 20
			if ff.Limit == 0 {
				ff.Limit = defaultLimit
			}

			c := MustGetClient()
			ctx := context.Background()
			sessionID, err := resolveSessionID(ctx, c, ff.Project, ff.ProjectID, "trace list")
			if err != nil {
				ExitErrorf("%v", err)
			}

			params := BuildRunQueryParams(&ff, true, ff.Limit)
			if sel := buildRunSelect(includeIO, includeFeedback); sel != nil {
				params.Select = langsmith.F(sel)
			}
			v2Select := buildRunSelectV2(includeIO, includeFeedback)
			since, before := feedbackWindow(&ff)
			fallback := false
			if !since.IsZero() || !before.IsZero() {
				ok, err := c.SupportsTraceFeedbackWindow(ctx)
				if err != nil {
					ExitErrorf("%v", err)
				}
				fallback = !ok
			}
			var (
				runs []langsmith.RunSchema
				next string
			)
			if fallback {
				// /v2/traces/query on this server ignores the feedback window, so
				// the roots come from /v2/runs/query and the trace totals are
				// added after.
				if ff.MinTokens > 0 {
					ExitErrorf("--min-tokens cannot be combined with --feedback-since or --feedback-before on this deployment (needs LangSmith 0.18.15 or later)")
				}
				runs, next, err = queryRunsAutoPage(ctx, c, params, v2Select, sessionID, ff.Cursor, ff.Limit, 0, feedbackWindowOption(ctx, c, &ff))
				if err == nil {
					err = hydrateTraceAggregates(ctx, c, sessionID, runs)
				}
			} else {
				runs, next, err = queryTracesAutoPage(ctx, c, params, v2Select, sessionID, ff.Cursor, ff.Limit, ff.MinTokens, traceFeedbackWindow(since, before))
			}
			if err != nil {
				ExitErrorf("%v", err)
			}

			var flaggedByTrace map[string]string
			if includeFlagged {
				flagged := FetchFlaggedTraces(ctx, c, sessionID, ff.LastNMinutes)
				flaggedByTrace = make(map[string]string, len(flagged))
				for _, ft := range flagged {
					flaggedByTrace[ft.TraceID] = ft.Comment
				}
			}

			fmt_ := GetFormat()

			if fmt_ == "pretty" {
				if showHierarchy {
					for _, run := range runs {
						allRuns, err := fetchTraceRuns(ctx, c, sessionID, run.TraceID, run.StartTime, time.Time{}, nil, v2Select)
						if err != nil {
							ExitErrorf("%v", err)
						}
						output.OutputTree(runsToTreeData(allRuns), "")
					}
				} else {
					data := extractRunsToMaps(runs, includeMetadata, includeIO, includeFeedback)
					output.PrintRunsTable(os.Stdout, data, includeMetadata, "Traces")
				}
				printNextCursorHint(next)
			} else {
				withCursors := func(items any) any {
					if !cmd.Flags().Changed("cursor") {
						return items
					}
					return map[string]any{"traces": items, "cursors": cursorsField(next)}
				}
				if showHierarchy {
					baseSelect := buildRunSelect(includeIO, includeFeedback)
					var result []map[string]any
					for _, run := range runs {
						allRuns, err := fetchTraceRuns(ctx, c, sessionID, run.TraceID, run.StartTime, time.Time{}, baseSelect, v2Select)
						if err != nil {
							ExitErrorf("%v", err)
						}
						result = append(result, map[string]any{
							"trace_id":  run.TraceID,
							"run_count": len(allRuns),
							"runs":      extractRunsToMaps(allRuns, includeMetadata, includeIO, includeFeedback),
						})
					}
					if err := output.OutputJSON(withCursors(result), outputFile); err != nil {
						ExitErrorf("%v", err)
					}
				} else {
					data := extractRunsToMaps(runs, includeMetadata, includeIO, includeFeedback)
					if flaggedByTrace != nil {
						annotateFlagged(data, flaggedByTrace)
					}
					if err := output.OutputJSON(withCursors(data), outputFile); err != nil {
						ExitErrorf("%v", err)
					}
				}
			}
		},
	}

	addCommonFilterFlags(cmd, &ff, false)
	addFeedbackWindowFlags(cmd, &ff)
	addCursorFlag(cmd, &ff.Cursor)
	cmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Add status, duration_ms, first_token_time, token_usage, costs, tags, custom_metadata (incl. revision_id)")
	cmd.Flags().BoolVar(&includeIO, "include-io", false, "Add inputs, outputs, error, and events fields")
	cmd.Flags().BoolVar(&includeFeedback, "include-feedback", false, "Add feedback_stats field")
	cmd.Flags().BoolVar(&includeFlagged, "include-flagged", false, "Add flagged_comment field populated from user-flagged trace feedback")
	cmd.Flags().BoolVar(&full, "full", false, "Shorthand for --include-metadata --include-io --include-feedback")
	cmd.Flags().BoolVar(&showHierarchy, "show-hierarchy", false, "Fetch the full run tree for each trace")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")

	return cmd
}

func newTraceGetCmd() *cobra.Command {
	var (
		project         string
		projectID       string
		since           string
		lastNMinutes    int
		includeMetadata bool
		includeIO       bool
		includeFeedback bool
		full            bool
		outputFile      string
	)

	cmd := &cobra.Command{
		Use:   "get TRACE_ID",
		Short: "Fetch every run in a single trace",
		Long: `Fetch every run in a single trace.

--project and --project-id are optional. Without them (or $LANGSMITH_PROJECT)
the trace is looked up by ID across the current workspace's projects and
experiments. That lookup only finds traces whose root run started in the last
7 days unless --since or --last-n-minutes is set, searches projects that were
active or created in that window (at most 1,000 of each), and skips projects
whose runs you can't read; pass the project to skip it.

LangSmith Cloud and self-hosted >= 0.16 then read the whole trace through
/v2/traces/{trace_id}/runs. Older self-hosted deployments read it through the
v1 runs query, which keeps the same 7-day default window.`,
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			traceID := args[0]
			// The v2 endpoint puts the ID in the URL path.
			if _, err := uuid.Parse(traceID); err != nil {
				ExitErrorf("invalid trace ID %q: must be a UUID", traceID)
			}

			if full {
				includeMetadata = true
				includeIO = true
				includeFeedback = true
			}

			c := MustGetClient()
			ctx := context.Background()
			runs, err := traceGetRuns(ctx, c, project, projectID, traceID, since, lastNMinutes,
				buildRunSelect(includeIO, includeFeedback), buildRunSelectV2(includeIO, includeFeedback))
			if err != nil {
				ExitErrorf("%v", err)
			}

			fmt_ := GetFormat()

			if fmt_ == "pretty" {
				output.OutputTree(runsToTreeData(runs), "")
			} else {
				data := map[string]any{
					"trace_id":  traceID,
					"run_count": len(runs),
					"runs":      extractRunsToMaps(runs, includeMetadata, includeIO, includeFeedback),
				}
				if err := output.OutputJSON(data, outputFile); err != nil {
					ExitErrorf("%v", err)
				}
			}
		},
	}

	addProjectFlags(cmd, &project, &projectID)
	cmd.Flags().StringVar(&since, "since", "", "Only include runs after this timestamp, e.g. 2024-01-15T00:00:00Z (overrides 7-day default)")
	cmd.Flags().IntVar(&lastNMinutes, "last-n-minutes", 0, "Only include runs from the last N minutes, e.g. 60 (overrides 7-day default)")
	cmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Add status, duration_ms, first_token_time, token_usage, costs, tags, custom_metadata (incl. revision_id)")
	cmd.Flags().BoolVar(&includeIO, "include-io", false, "Add inputs, outputs, error, and events fields")
	cmd.Flags().BoolVar(&includeFeedback, "include-feedback", false, "Add feedback_stats field")
	cmd.Flags().BoolVar(&full, "full", false, "Shorthand for --include-metadata --include-io --include-feedback")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")

	return cmd
}

// traceGetRuns returns every run in traceID, oldest first, or an error when the
// trace has no runs in scope. Without a project it searches the current
// workspace. The v1 runs query filters by trace alone; every v2 runs and trace
// endpoint needs a project, so on v2 the trace's project is first found among
// the workspace's recently active projects (see traceProjectID) and the trace
// is then fetched from /v2/traces/{trace_id}/runs.
func traceGetRuns(ctx context.Context, c *client.Client, project, projectID, traceID, since string, lastNMinutes int, v1Select []langsmith.RunQueryParamsSelect, v2Selects []langsmith.RunSelectField) ([]langsmith.RunSchema, error) {
	useV2, err := c.UseV2API(ctx)
	if err != nil {
		return nil, err
	}
	v1Start := resolveStartTime(since, lastNMinutes)
	// v2 finds the trace without a window, so it only bounds the runs when a
	// window was asked for; v1 and the workspace lookup keep the 7-day default.
	var v2Since time.Time
	if since != "" || lastNMinutes > 0 {
		v2Since = v1Start
	}

	var sessionID string
	inProject := projectID != "" || ResolveProject(project) != ""
	if inProject {
		if sessionID, err = resolveSessionID(ctx, c, project, projectID, "trace get"); err != nil {
			return nil, err
		}
	} else if useV2 {
		if sessionID, err = traceProjectID(ctx, c, traceID, v1Start); err != nil {
			return nil, err
		}
		if sessionID == "" {
			return nil, traceNotFound(traceID, false, v1Start)
		}
	}

	runs, err := fetchTraceRuns(ctx, c, sessionID, traceID, v1Start, v2Since, v1Select, v2Selects)
	if err != nil {
		return nil, err
	}
	if len(runs) == 0 {
		window := v1Start
		if useV2 && inProject {
			window = v2Since
		}
		return nil, traceNotFound(traceID, inProject, window)
	}
	return runs, nil
}

// traceLookupPageSize is how many projects traceProjectID lists, and then
// searches for the trace, per request.
const traceLookupPageSize = 100

// traceLookupMaxPages bounds each pass of the workspace lookup to
// traceLookupMaxPages*traceLookupPageSize projects.
const traceLookupMaxPages = 10

// traceLookupActivityLag is how far a project's last_run_start_time may trail
// its newest run. The server updates it asynchronously, so a project that just
// received a run can still report an older time.
const traceLookupActivityLag = time.Hour

// traceProjectID returns the project (session) ID of traceID, searching the
// current workspace; "" when no root run of traceID started at or after start.
//
// Every v2 runs and trace endpoint needs a project. Where the deployment lists
// projects through /v2/sessions, the lookup searches, page by page:
//  1. tracing projects, most recently active first, until their last run
//     started more than traceLookupActivityLag before start;
//  2. projects and experiments with no recorded activity yet (new projects,
//     and experiments, which have none), newest first, until one was created
//     more than traceLookupActivityLag before start.
//
// Older self-hosted deployments use the v1 runs query instead.
func traceProjectID(ctx context.Context, c *client.Client, traceID string, start time.Time) (string, error) {
	ok, err := c.SupportsSessionsListV2(ctx)
	if err != nil {
		return "", err
	}
	if !ok {
		return traceProjectIDV1(ctx, c, traceID, start)
	}

	cutoff := start.Add(-traceLookupActivityLag)
	searched := map[string]bool{}
	active := langsmith.SessionListV2Params{
		Type: langsmith.F([]langsmith.SessionListV2ParamsType{
			langsmith.SessionListV2ParamsTypeProject,
			langsmith.SessionListV2ParamsTypeAgent,
		}),
		SortBy: langsmith.F(langsmith.SessionListV2ParamsSortByLastRunStartTime),
	}
	id, activeCapped, err := traceSearchProjects(ctx, c, traceID, start, active, searched, func(s langsmith.Session) (search, done bool) {
		if !s.LastRunStartTime.IsZero() && s.LastRunStartTime.Before(cutoff) {
			return false, true
		}
		return true, false
	})
	if err != nil || id != "" {
		return id, err
	}
	unrecorded := langsmith.SessionListV2Params{
		Type: langsmith.F([]langsmith.SessionListV2ParamsType{
			langsmith.SessionListV2ParamsTypeProject,
			langsmith.SessionListV2ParamsTypeAgent,
			langsmith.SessionListV2ParamsTypeExperiment,
		}),
		SortBy: langsmith.F(langsmith.SessionListV2ParamsSortByStartTime),
	}
	id, newCapped, err := traceSearchProjects(ctx, c, traceID, start, unrecorded, searched, func(s langsmith.Session) (search, done bool) {
		if s.StartTime.Before(cutoff) {
			return false, true
		}
		return s.LastRunStartTime.IsZero(), false
	})
	if err != nil || id != "" {
		return id, err
	}
	if activeCapped || newCapped {
		return "", fmt.Errorf("trace %s not found in the %d most recently active projects; pass --project or --project-id",
			traceID, traceLookupMaxPages*traceLookupPageSize)
	}
	return "", nil
}

// traceSearchProjects pages through the projects params lists, newest first,
// and searches each page's projects that pick selects and searched doesn't
// hold yet. It stops at the first project pick marks done. capped reports that
// it stopped at traceLookupMaxPages with more projects left.
func traceSearchProjects(
	ctx context.Context,
	c *client.Client,
	traceID string,
	start time.Time,
	params langsmith.SessionListV2Params,
	searched map[string]bool,
	pick func(langsmith.Session) (search, done bool),
) (projectID string, capped bool, err error) {
	params.SortOrder = langsmith.F(langsmith.SessionListV2ParamsSortOrderDesc)
	params.PageSize = langsmith.F(int64(traceLookupPageSize))
	for page := 0; page < traceLookupMaxPages; page++ {
		resp, err := c.SDK.Sessions.ListV2(ctx, params)
		if err != nil {
			return "", false, fmt.Errorf("looking up trace %s in the workspace: listing projects: %w", traceID, err)
		}
		var ids []string
		done := false
		for _, s := range resp.Items {
			search, stop := pick(s)
			if stop {
				done = true
				break
			}
			if search && !searched[s.ID] {
				searched[s.ID] = true
				ids = append(ids, s.ID)
			}
		}
		if len(ids) > 0 {
			sessionID, err := traceRootProjectID(ctx, c, traceID, ids, start)
			if err != nil || sessionID != "" {
				return sessionID, false, err
			}
		}
		if done || resp.NextCursor == "" {
			return "", false, nil
		}
		params.Cursor = langsmith.F(resp.NextCursor)
	}
	return "", true, nil
}

// traceRootProjectID returns the project of traceID's root run among
// projectIDs, or "" when none of them has a root run of the trace that started
// at or after start. The runs query rejects the whole request when the caller
// can't read runs in any one project, so a rejected batch is split in half and
// projects whose runs can't be read are skipped.
func traceRootProjectID(ctx context.Context, c *client.Client, traceID string, projectIDs []string, start time.Time) (string, error) {
	page, err := c.SDK.Runs.QueryV2(client.ExpectForbidden(ctx), langsmith.RunQueryV2Params{
		ProjectIDs:   langsmith.F(projectIDs),
		TraceID:      langsmith.F(traceID),
		IsRoot:       langsmith.F(true),
		MinStartTime: langsmith.F(start),
		PageSize:     langsmith.F(int64(1)),
		Selects:      langsmith.F([]langsmith.RunSelectField{langsmith.RunSelectFieldID, langsmith.RunSelectFieldProjectID}),
	})
	if isHTTPStatus(err, http.StatusForbidden) {
		if len(projectIDs) == 1 {
			return "", nil
		}
		mid := len(projectIDs) / 2
		if id, err := traceRootProjectID(ctx, c, traceID, projectIDs[:mid], start); err != nil || id != "" {
			return id, err
		}
		return traceRootProjectID(ctx, c, traceID, projectIDs[mid:], start)
	}
	if err != nil {
		return "", fmt.Errorf("looking up trace %s in the workspace: %w", traceID, err)
	}
	if len(page.Items) == 0 {
		return "", nil
	}
	return page.Items[0].ProjectID, nil
}

// traceProjectIDV1 is traceProjectID for self-hosted deployments that predate
// the v2 project list. It reads the root run's project from the v1 runs query,
// which filters by trace alone. The server ignores limit when filtering by
// trace, so is_root keeps the response to one row. Deployments without
// ClickHouse querying (SmithDB only) reject the v1 query with 501; the error
// then asks for the project.
func traceProjectIDV1(ctx context.Context, c *client.Client, traceID string, start time.Time) (string, error) {
	params := langsmith.RunQueryParams{
		Trace:     langsmith.F(traceID),
		IsRoot:    langsmith.F(true),
		StartTime: langsmith.F(start),
		Limit:     langsmith.F(int64(1)),
		Select:    langsmith.F([]langsmith.RunQueryParamsSelect{langsmith.RunQueryParamsSelectSessionID}),
	}
	runs, err := queryRuns(ctx, c, params, "", 1, 0)
	if isHTTPStatus(err, http.StatusNotImplemented) {
		return "", fmt.Errorf("this deployment cannot look up trace %s without a project; pass --project or --project-id", traceID)
	}
	if err != nil {
		return "", fmt.Errorf("looking up trace %s in the workspace: %w", traceID, err)
	}
	if len(runs) == 0 {
		return "", nil
	}
	return runs[0].SessionID, nil
}

// traceNotFound explains where trace get searched. since is zero when the
// search had no start-time bound.
func traceNotFound(traceID string, inProject bool, since time.Time) error {
	scope := "the current workspace"
	if inProject {
		scope = "the given project"
	}
	if since.IsZero() {
		return fmt.Errorf("trace %s not found in %s", traceID, scope)
	}
	return fmt.Errorf("trace %s not found in %s among runs since %s; pass --since with an earlier time if the trace is older",
		traceID, scope, since.UTC().Format(time.RFC3339))
}

func newTraceExportCmd() *cobra.Command {
	var (
		ff              FilterFlags
		includeMetadata bool
		includeIO       bool
		includeFeedback bool
		full            bool
		filenamePattern string
	)

	cmd := &cobra.Command{
		Use:   "export OUTPUT_DIR",
		Short: "Export traces to a directory as JSONL files (one file per trace)",
		Args:  cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			outputDir := args[0]

			if full {
				includeMetadata = true
				includeIO = true
				includeFeedback = true
			}

			if ff.Limit == 0 {
				ff.Limit = 10
			}

			if err := os.MkdirAll(outputDir, 0755); err != nil {
				ExitErrorf("creating output directory: %v", err)
			}

			c := MustGetClient()
			ctx := context.Background()
			sessionID, err := resolveSessionID(ctx, c, ff.Project, ff.ProjectID, "trace export")
			if err != nil {
				ExitErrorf("%v", err)
			}

			params := BuildRunQueryParams(&ff, true, ff.Limit)
			sel := buildRunSelect(includeIO, includeFeedback)
			if sel != nil {
				params.Select = langsmith.F(sel)
			}
			v2Select := buildRunSelectV2(includeIO, includeFeedback)
			rootRuns, err := queryTracesAuto(ctx, c, params, v2Select, sessionID, ff.Limit, ff.MinTokens)
			if err != nil {
				ExitErrorf("%v", err)
			}

			// Resolve every destination before creating any files. A pattern such as
			// "trace.jsonl" or {name} can map multiple traces to one path; allowing
			// os.Create to handle that would silently truncate an earlier export.
			paths := make(map[string]string, len(rootRuns))
			for _, root := range rootRuns {
				fpath := traceExportPath(outputDir, filenamePattern, root)
				if previous, ok := paths[fpath]; ok {
					ExitErrorf("filename pattern resolves multiple traces to %s (trace IDs %s and %s); include {trace_id} or another unique placeholder", fpath, previous, root.TraceID)
				}
				paths[fpath] = root.TraceID
			}

			exported := 0
			for _, root := range rootRuns {
				tid := root.TraceID

				allRuns, err := fetchTraceRuns(ctx, c, sessionID, tid, root.StartTime, time.Time{}, sel, v2Select)
				if err != nil {
					ExitErrorf("%v", err)
				}

				fpath := traceExportPath(outputDir, filenamePattern, root)

				f, err := os.Create(fpath)
				if err != nil {
					ExitErrorf("creating file %s: %v", fpath, err)
				}

				for _, run := range allRuns {
					data := extractRunsToMaps([]langsmith.RunSchema{run}, includeMetadata, includeIO, includeFeedback)
					line, _ := json.Marshal(data[0])
					_, _ = f.Write(line)
					_, _ = f.WriteString("\n")
				}
				f.Close()
				exported++
			}
			if err := output.OutputJSON(map[string]any{
				"status":     "exported",
				"count":      exported,
				"output_dir": outputDir,
			}, ""); err != nil {
				ExitErrorf("%v", err)
			}
		},
	}

	addCommonFilterFlags(cmd, &ff, false)
	cmd.Flags().BoolVar(&includeMetadata, "include-metadata", false, "Add status, duration_ms, first_token_time, token_usage, costs, tags, custom_metadata (incl. revision_id)")
	cmd.Flags().BoolVar(&includeIO, "include-io", false, "Add inputs, outputs, error, and events fields")
	cmd.Flags().BoolVar(&includeFeedback, "include-feedback", false, "Add feedback_stats field")
	cmd.Flags().BoolVar(&full, "full", false, "Shorthand for --include-metadata --include-io --include-feedback")
	cmd.Flags().StringVar(&filenamePattern, "filename-pattern", "{trace_id}.jsonl",
		"Filename pattern. Supports {trace_id} and {name} placeholders.")

	return cmd
}

func traceExportPath(outputDir, pattern string, root langsmith.RunSchema) string {
	name := root.Name
	if name == "" {
		name = "unknown"
	}
	filename := strings.ReplaceAll(pattern, "{trace_id}", root.TraceID)
	filename = strings.ReplaceAll(filename, "{name}", name)
	return filepath.Join(outputDir, filepath.Base(filename))
}

// annotateFlagged mutates each run map to add a "flagged_comment" field
// (string or null) based on whether the trace_id has a user-flagged feedback
// entry. Runs without a match get "flagged_comment": null so downstream jq
// filters can distinguish "missing" from "not flagged".
func annotateFlagged(data []map[string]any, flaggedByTrace map[string]string) {
	for _, row := range data {
		tid, _ := row["trace_id"].(string)
		if comment, ok := flaggedByTrace[tid]; ok {
			if comment == "" {
				comment = "User flagged for review"
			}
			row["flagged_comment"] = comment
		} else {
			row["flagged_comment"] = nil
		}
	}
}
