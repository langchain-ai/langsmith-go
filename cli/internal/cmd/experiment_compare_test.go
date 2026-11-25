package cmd

import (
	"encoding/json"
	"io"
	"net/http"
	"strings"
	"testing"
)

const (
	compareTenantID  = "019a0000-0000-7000-8000-0000000000ff"
	compareDatasetID = "019a0000-0000-7000-8000-0000000000d1"
	compareExpA      = "019a0000-0000-7000-8000-0000000000a1"
	compareExpB      = "019a0000-0000-7000-8000-0000000000b1"
)

func compareSessionJSON(id, name, datasetID, stats string) string {
	return `{"id":"` + id + `","tenant_id":"` + compareTenantID + `","name":"` + name + `","start_time":"2026-10-26T00:00:00Z",` +
		`"reference_dataset_id":"` + datasetID + `",` + stats + `}`
}

// compareRun renders one experiment run with the given feedback scores and error.
func compareRun(id, project, errMsg string, scores map[string]float64) string {
	stats := map[string]any{}
	for k, v := range scores {
		stats[k] = map[string]any{"avg": v, "n": 1}
	}
	run := map[string]any{"id": id, "project_id": project, "feedback_stats": stats}
	if errMsg != "" {
		run["error"] = errMsg
	}
	b, _ := json.Marshal(run)
	return string(b)
}

func compareExample(id string, runs ...string) string {
	return `{"id":"` + id + `","dataset_id":"` + compareDatasetID + `","runs":[` + strings.Join(runs, ",") + `]}`
}

// newCompareServer serves two experiments on compareDatasetID (bDatasetID for
// B) and two pages of experiment runs, recording each experiment-runs request body.
func newCompareServer(t *testing.T, bDatasetID string, bodies *[]map[string]any) string {
	t.Helper()
	pages := []string{
		`{"items":[` +
			compareExample("ex-pass-fail",
				compareRun("r1", compareExpA, "", map[string]float64{"correctness": 1}),
				compareRun("r2", compareExpB, "", map[string]float64{"correctness": 0})) + "," +
			compareExample("ex-same",
				compareRun("r3", compareExpA, "", map[string]float64{"correctness": 1}),
				compareRun("r4", compareExpB, "", map[string]float64{"correctness": 1})) +
			`],"next_cursor":"page2"}`,
		`{"items":[` +
			compareExample("ex-drop-and-error",
				compareRun("r5", compareExpA, "", map[string]float64{"relevance": 0.9}),
				compareRun("r6", compareExpB, "ValueError: boom", map[string]float64{"relevance": 0.4})) + "," +
			compareExample("ex-improved",
				compareRun("r7", compareExpA, "", map[string]float64{"correctness": 0}),
				compareRun("r8", compareExpB, "", map[string]float64{"correctness": 1})) + "," +
			compareExample("ex-only-a",
				compareRun("r9", compareExpA, "", map[string]float64{"correctness": 1})) +
			`],"next_cursor":null}`,
	}
	srv := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sessions/"+compareExpA:
			_, _ = io.WriteString(w, compareSessionJSON(compareExpA, "baseline", compareDatasetID,
				`"feedback_stats":{"correctness":{"n":4,"avg":0.75,"type":"NUMERIC"},"relevance":{"n":4,"avg":0.9,"type":"NUMERIC"}},`+
					`"latency_p50":1.5,"latency_p99":4,"total_tokens":1000,"total_cost":0.5,"run_count":5`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sessions/"+compareExpB:
			_, _ = io.WriteString(w, compareSessionJSON(compareExpB, "candidate", bDatasetID,
				`"feedback_stats":{"correctness":{"n":4,"avg":0.5,"type":"NUMERIC"},"helpfulness":{"n":4,"avg":0.8,"type":"NUMERIC"}},`+
					`"latency_p50":1,"latency_p99":null,"total_tokens":null,"error_rate":0.25,"total_cost":0.25,"run_count":4`))
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/datasets/"+compareDatasetID+"/experiment-runs":
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decoding experiment-runs body: %v", err)
			}
			*bodies = append(*bodies, body)
			page := 0
			if body["cursor"] == "page2" {
				page = 1
			}
			_, _ = io.WriteString(w, pages[page])
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.NotFound(w, r)
		}
	})
	return srv.URL
}

func runCompare(t *testing.T, format string, args ...string) (string, error) {
	t.Helper()
	flagOutputFormat = format
	var err error
	out := captureStdout(t, func() {
		cmd := newExperimentCmd()
		cmd.SetArgs(append([]string{"compare"}, args...))
		err = cmd.Execute()
	})
	return out, err
}

func TestExperimentCompare_JSON(t *testing.T) {
	var bodies []map[string]any
	defer setupTestEnv(t, newCompareServer(t, compareDatasetID, &bodies))()

	out, err := runCompare(t, "json", compareExpA, compareExpB)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}

	var got experimentComparison
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not a comparison object: %v\n%s", err, out)
	}
	if got.DatasetID != compareDatasetID || got.A.Name != "baseline" || got.B.Name != "candidate" {
		t.Errorf("header = %s %+v %+v", got.DatasetID, got.A, got.B)
	}

	assertMetric := func(section, key string, m metricChange, a, b, delta *float64) {
		t.Helper()
		eq := func(x, y *float64) bool {
			return (x == nil) == (y == nil) && (x == nil || *x-*y < 1e-9 && *y-*x < 1e-9)
		}
		if !eq(m.A, a) || !eq(m.B, b) || !eq(m.Delta, delta) {
			t.Errorf("%s.%s = %s, want a=%v b=%v delta=%v", section, key, mustJSON(t, m), deref(a), deref(b), deref(delta))
		}
	}
	f := func(v float64) *float64 { return &v }
	assertMetric("feedback", "correctness", got.Feedback["correctness"], f(0.75), f(0.5), f(-0.25))
	assertMetric("feedback", "relevance", got.Feedback["relevance"], f(0.9), nil, nil)
	assertMetric("feedback", "helpfulness", got.Feedback["helpfulness"], nil, f(0.8), nil)
	assertMetric("run_stats", "total_cost", got.RunStats["total_cost"], f(0.5), f(0.25), f(-0.25))
	assertMetric("run_stats", "latency_p50", got.RunStats["latency_p50"], f(1.5), f(1), f(-0.5))
	assertMetric("run_stats", "latency_p99", got.RunStats["latency_p99"], f(4), nil, nil)
	assertMetric("run_stats", "total_tokens", got.RunStats["total_tokens"], f(1000), nil, nil)
	assertMetric("run_stats", "error_rate", got.RunStats["error_rate"], nil, f(0.25), nil)

	if got.ExamplesCompared != 4 || got.RegressionCount != 2 {
		t.Errorf("examples_compared=%d regression_count=%d, want 4 and 2", got.ExamplesCompared, got.RegressionCount)
	}
	wantRegressions := `[{"example_id":"ex-pass-fail","changes":[{"key":"correctness","kind":"pass_to_fail","a":1,"b":0}]},` +
		`{"example_id":"ex-drop-and-error","changes":[{"key":"error","kind":"new_error","a":null,"b":"ValueError: boom"},` +
		`{"key":"relevance","kind":"score_drop","a":0.9,"b":0.4}]}]`
	if gotRegressions := mustJSON(t, got.Regressions); gotRegressions != wantRegressions {
		t.Errorf("regressions =\n%s\nwant\n%s", gotRegressions, wantRegressions)
	}

	if len(bodies) != 2 {
		t.Fatalf("experiment-runs requests = %d, want 2 pages", len(bodies))
	}
	if ids := mustJSON(t, bodies[0]["experiment_ids"]); ids != `["`+compareExpA+`","`+compareExpB+`"]` {
		t.Errorf("experiment_ids = %s", ids)
	}
	selects := mustJSON(t, bodies[0]["selects"])
	for _, field := range []string{"PROJECT_ID", "ERROR", "FEEDBACK_STATS"} {
		if !strings.Contains(selects, field) {
			t.Errorf("selects %s missing %s", selects, field)
		}
	}
}

func TestExperimentCompare_LimitCapsListedRegressions(t *testing.T) {
	var bodies []map[string]any
	defer setupTestEnv(t, newCompareServer(t, compareDatasetID, &bodies))()

	out, err := runCompare(t, "json", compareExpA, compareExpB, "--limit", "1")
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	var got experimentComparison
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("output is not a comparison object: %v\n%s", err, out)
	}
	if got.RegressionCount != 2 || len(got.Regressions) != 1 || got.Regressions[0].ExampleID != "ex-pass-fail" {
		t.Errorf("regression_count=%d regressions=%s, want count 2 and only ex-pass-fail listed", got.RegressionCount, mustJSON(t, got.Regressions))
	}
}

func TestExperimentCompare_Pretty(t *testing.T) {
	var bodies []map[string]any
	defer setupTestEnv(t, newCompareServer(t, compareDatasetID, &bodies))()

	out, err := runCompare(t, "pretty", compareExpA, compareExpB)
	if err != nil {
		t.Fatalf("execute: %v", err)
	}
	for _, want := range []string{"baseline", "candidate", "correctness", "-0.25", "latency_p99", "Regressions: 2 of 4 examples", "ex-pass-fail", "pass_to_fail", "new_error"} {
		if !strings.Contains(out, want) {
			t.Errorf("pretty output missing %q:\n%s", want, out)
		}
	}
}

func TestExperimentCompare_DifferentDatasets(t *testing.T) {
	var bodies []map[string]any
	defer setupTestEnv(t, newCompareServer(t, "019a0000-0000-7000-8000-0000000000d2", &bodies))()

	_, err := runCompare(t, "json", compareExpA, compareExpB)
	if err == nil || !strings.Contains(err.Error(), "different datasets") {
		t.Fatalf("error = %v, want a different-datasets error", err)
	}
	if len(bodies) != 0 {
		t.Errorf("read experiment runs %d times despite the dataset mismatch", len(bodies))
	}
}

func TestExperimentCompare_SameExperiment(t *testing.T) {
	var bodies []map[string]any
	defer setupTestEnv(t, newCompareServer(t, compareDatasetID, &bodies))()

	_, err := runCompare(t, "json", compareExpA, compareExpA)
	if err == nil || !strings.Contains(err.Error(), "same experiment") {
		t.Fatalf("error = %v, want a same-experiment error", err)
	}
	if len(bodies) != 0 {
		t.Errorf("read experiment runs %d times for a self-comparison", len(bodies))
	}
}

func deref(v *float64) any {
	if v == nil {
		return nil
	}
	return *v
}
