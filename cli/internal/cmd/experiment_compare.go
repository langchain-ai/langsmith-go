package cmd

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

// experimentRunsPageSize is the largest page the experiment-runs endpoint accepts.
const experimentRunsPageSize = 100

func newExperimentCompareCmd() *cobra.Command {
	var (
		limit      int
		outputFile string
	)

	cmd := &cobra.Command{
		Use:   "compare BASELINE CANDIDATE",
		Short: "Compare two experiments on the same dataset",
		Long: `Compare a candidate experiment with a baseline on the same dataset.

Prints, for baseline (a) and candidate (b), and the change (b - a):
  - each feedback key's average score
  - total cost, latency p50 and p99 (seconds), total tokens, and error rate
  - the examples that got worse: a feedback score that went from 1 to 0
    (pass_to_fail) or dropped (score_drop), or a run that errored in the
    candidate but not the baseline (new_error)

Scores are treated as higher-is-better. Each experiment is given by name or
ID; they must be two different experiments on the same dataset. --limit caps the regressed examples
listed; regression_count is always the full count.

Examples:
  langsmith experiments compare baseline-exp candidate-exp
  langsmith experiments compare baseline-exp candidate-exp --format json --limit 50
  langsmith experiments compare baseline-exp candidate-exp --jq '.feedback.correctness.delta'`,
		Args: cobra.ExactArgs(2),
		RunE: func(cmd *cobra.Command, args []string) error {
			c, err := getClient()
			if err != nil {
				return err
			}
			result, err := compareExperiments(cmd.Context(), c, args[0], args[1], limit)
			if err != nil {
				return err
			}
			if GetFormat() == "pretty" && outputFile == "" {
				printExperimentComparison(result)
				return nil
			}
			return output.OutputJSON(result, outputFile)
		},
	}

	cmd.Flags().IntVarP(&limit, "limit", "n", 20, "Maximum number of regressed examples to list (0 for all)")
	cmd.Flags().StringVarP(&outputFile, "output", "o", "", "Write JSON output to a file")
	return cmd
}

type experimentRef struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

// metricChange is one metric in both experiments. A, B and Delta are nil
// when a side has no value.
type metricChange struct {
	A     *float64 `json:"a"`
	B     *float64 `json:"b"`
	Delta *float64 `json:"delta"`
}

type exampleChange struct {
	Key  string `json:"key"`
	Kind string `json:"kind"` // pass_to_fail, score_drop, or new_error
	A    any    `json:"a"`
	B    any    `json:"b"`
}

type exampleRegression struct {
	ExampleID string          `json:"example_id"`
	Changes   []exampleChange `json:"changes"`
}

type experimentComparison struct {
	DatasetID        string                  `json:"dataset_id"`
	A                experimentRef           `json:"a"`
	B                experimentRef           `json:"b"`
	Feedback         map[string]metricChange `json:"feedback"`
	RunStats         map[string]metricChange `json:"run_stats"`
	ExamplesCompared int                     `json:"examples_compared"`
	RegressionCount  int                     `json:"regression_count"`
	Regressions      []exampleRegression     `json:"regressions"`
}

// runStatOrder is the order run stats are printed in.
var runStatOrder = []string{"total_cost", "latency_p50", "latency_p99", "total_tokens", "error_rate"}

func compareExperiments(ctx context.Context, c *client.Client, nameA, nameB string, limit int) (*experimentComparison, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	a, err := fetchExperiment(ctx, c, nameA)
	if err != nil {
		return nil, err
	}
	b, err := fetchExperiment(ctx, c, nameB)
	if err != nil {
		return nil, err
	}
	if a.ID == b.ID {
		return nil, fmt.Errorf("%q and %q are the same experiment (%s); pass two different experiments", nameA, nameB, a.ID)
	}
	if a.DatasetID == "" || b.DatasetID == "" {
		return nil, fmt.Errorf("both experiments must be evaluations on a dataset: %q has dataset %q, %q has dataset %q",
			a.Name, a.DatasetID, b.Name, b.DatasetID)
	}
	if a.DatasetID != b.DatasetID {
		return nil, fmt.Errorf("experiments use different datasets: %q is on dataset %s, %q is on dataset %s",
			a.Name, a.DatasetID, b.Name, b.DatasetID)
	}

	result := &experimentComparison{
		DatasetID:   a.DatasetID,
		A:           experimentRef{ID: a.ID, Name: a.Name},
		B:           experimentRef{ID: b.ID, Name: b.Name},
		Feedback:    map[string]metricChange{},
		RunStats:    map[string]metricChange{},
		Regressions: []exampleRegression{},
	}

	aScores, bScores := feedbackAverages(a.FeedbackStats), feedbackAverages(b.FeedbackStats)
	for key := range aScores {
		result.Feedback[key] = newMetricChange(aScores[key], bScores[key])
	}
	for key := range bScores {
		if _, ok := result.Feedback[key]; !ok {
			result.Feedback[key] = newMetricChange(nil, bScores[key])
		}
	}

	aStats, bStats := experimentRunStats(a), experimentRunStats(b)
	for _, key := range runStatOrder {
		result.RunStats[key] = newMetricChange(aStats[key], bStats[key])
	}

	compared, regressions, err := exampleRegressions(ctx, c, a.DatasetID, a.ID, b.ID)
	if err != nil {
		return nil, err
	}
	result.ExamplesCompared = compared
	result.RegressionCount = len(regressions)
	if limit > 0 && len(regressions) > limit {
		regressions = regressions[:limit]
	}
	result.Regressions = append(result.Regressions, regressions...)
	return result, nil
}

func newMetricChange(a, b *float64) metricChange {
	m := metricChange{A: a, B: b}
	if a != nil && b != nil {
		d := *b - *a
		m.Delta = &d
	}
	return m
}

// feedbackAverages returns the numeric average score per feedback key.
func feedbackAverages(stats map[string]any) map[string]*float64 {
	out := map[string]*float64{}
	for key, v := range stats {
		stat, ok := v.(map[string]any)
		if !ok {
			continue
		}
		if avg, ok := stat["avg"].(float64); ok {
			out[key] = &avg
		} else {
			out[key] = nil
		}
	}
	return out
}

func experimentRunStats(e experimentSummary) map[string]*float64 {
	stats := map[string]*float64{}
	if cost, ok := e.TotalCost.(float64); ok {
		stats["total_cost"] = &cost
	}
	if e.HasLatencyP50 {
		p50 := e.LatencyP50
		stats["latency_p50"] = &p50
	}
	if p99, ok := e.LatencyP99.(float64); ok {
		stats["latency_p99"] = &p99
	}
	if e.HasTokens {
		tokens := float64(e.TotalTokens)
		stats["total_tokens"] = &tokens
	}
	if e.HasErrorRate {
		errorRate := e.ErrorRate
		stats["error_rate"] = &errorRate
	}
	return stats
}

// exampleOutcome is one experiment's outcome on one dataset example, averaged
// over the experiment's runs (repetitions) of that example.
type exampleOutcome struct {
	scores map[string]float64
	err    string
}

// exampleRegressions reads both experiments' runs per dataset example and
// returns how many examples both experiments ran, and the examples where the
// candidate did worse than the baseline.
func exampleRegressions(ctx context.Context, c *client.Client, datasetID, aID, bID string) (int, []exampleRegression, error) {
	pager := c.SDK.Datasets.ExperimentRuns.QueryAutoPaging(ctx, datasetID, langsmith.DatasetExperimentRunQueryParams{
		ExperimentIDs: langsmith.F([]string{aID, bID}),
		PageSize:      langsmith.F(int64(experimentRunsPageSize)),
		Selects: langsmith.F([]langsmith.RunSelectField{
			langsmith.RunSelectFieldID,
			langsmith.RunSelectFieldProjectID,
			langsmith.RunSelectFieldError,
			langsmith.RunSelectFieldFeedbackStats,
		}),
	})
	compared := 0
	var regressions []exampleRegression
	for pager.Next() {
		example := pager.Current()
		var runsA, runsB []langsmith.Run
		for _, run := range example.Runs {
			switch run.ProjectID {
			case aID:
				runsA = append(runsA, run)
			case bID:
				runsB = append(runsB, run)
			}
		}
		if len(runsA) == 0 || len(runsB) == 0 {
			continue
		}
		compared++
		if changes := compareExampleOutcomes(summarizeRuns(runsA), summarizeRuns(runsB)); len(changes) > 0 {
			regressions = append(regressions, exampleRegression{ExampleID: example.ID, Changes: changes})
		}
	}
	if err := pager.Err(); err != nil {
		return 0, nil, fmt.Errorf("reading experiment runs: %w", err)
	}
	return compared, regressions, nil
}

func summarizeRuns(runs []langsmith.Run) exampleOutcome {
	sums := map[string]float64{}
	counts := map[string]int{}
	var result exampleOutcome
	for _, run := range runs {
		if run.Error != "" && result.err == "" {
			result.err = run.Error
		}
		for key, stat := range run.FeedbackStats {
			if stat.JSON.Avg.IsNull() || stat.JSON.Avg.IsMissing() {
				continue
			}
			sums[key] += stat.Avg
			counts[key]++
		}
	}
	result.scores = map[string]float64{}
	for key, sum := range sums {
		result.scores[key] = sum / float64(counts[key])
	}
	return result
}

// compareExampleOutcomes lists how b is worse than a on one example: a new
// error, then each feedback key whose score dropped, in key order.
func compareExampleOutcomes(a, b exampleOutcome) []exampleChange {
	var changes []exampleChange
	if a.err == "" && b.err != "" {
		changes = append(changes, exampleChange{Key: "error", Kind: "new_error", A: nil, B: b.err})
	}
	keys := make([]string, 0, len(a.scores))
	for key := range a.scores {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	for _, key := range keys {
		scoreA := a.scores[key]
		scoreB, ok := b.scores[key]
		if !ok || scoreB >= scoreA {
			continue
		}
		kind := "score_drop"
		if scoreA == 1 && scoreB == 0 {
			kind = "pass_to_fail"
		}
		changes = append(changes, exampleChange{Key: key, Kind: kind, A: scoreA, B: scoreB})
	}
	return changes
}

func printExperimentComparison(r *experimentComparison) {
	fmt.Printf("Dataset: %s\n", r.DatasetID)
	fmt.Printf("A (baseline):  %s (%s)\n", r.A.Name, r.A.ID)
	fmt.Printf("B (candidate): %s (%s)\n\n", r.B.Name, r.B.ID)

	columns := []string{"Metric", "A", "B", "Change"}
	keys := make([]string, 0, len(r.Feedback))
	for key := range r.Feedback {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	var rows [][]string
	for _, key := range keys {
		rows = append(rows, metricRow(key, r.Feedback[key]))
	}
	if len(rows) > 0 {
		output.OutputTable(columns, rows, "Feedback")
		fmt.Println()
	}

	rows = nil
	for _, key := range runStatOrder {
		rows = append(rows, metricRow(key, r.RunStats[key]))
	}
	output.OutputTable(columns, rows, "Run stats")
	fmt.Println()

	title := fmt.Sprintf("Regressions: %d of %d examples", r.RegressionCount, r.ExamplesCompared)
	if len(r.Regressions) < r.RegressionCount {
		title += fmt.Sprintf(" (showing %d; use --limit)", len(r.Regressions))
	}
	if r.RegressionCount == 0 {
		fmt.Println(title)
		return
	}
	rows = nil
	for _, reg := range r.Regressions {
		for _, ch := range reg.Changes {
			rows = append(rows, []string{reg.ExampleID, ch.Key, ch.Kind, formatChangeValue(ch.A), formatChangeValue(ch.B)})
		}
	}
	output.OutputTable([]string{"Example ID", "Key", "Kind", "A", "B"}, rows, title)
}

func metricRow(name string, m metricChange) []string {
	change := "N/A"
	if m.Delta != nil {
		change = formatMetric(*m.Delta)
		if *m.Delta > 0 {
			change = "+" + change
		}
	}
	return []string{name, formatOptionalMetric(m.A), formatOptionalMetric(m.B), change}
}

func formatOptionalMetric(v *float64) string {
	if v == nil {
		return "N/A"
	}
	return formatMetric(*v)
}

func formatMetric(v float64) string {
	return strconv.FormatFloat(v, 'g', 6, 64)
}

func formatChangeValue(v any) string {
	switch x := v.(type) {
	case nil:
		return "-"
	case float64:
		return formatMetric(x)
	case string:
		return truncate(x, 60)
	default:
		return fmt.Sprint(x)
	}
}
