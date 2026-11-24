package cmd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"slices"
	"strings"
	"testing"
)

const (
	fromRunDatasetID = "11111111-1111-4111-8111-111111111111"
	fromRunSession   = "22222222-2222-4222-8222-222222222222"
	runA             = "aaaaaaaa-aaaa-4aaa-8aaa-aaaaaaaaaaaa"
	runB             = "bbbbbbbb-bbbb-4bbb-8bbb-bbbbbbbbbbbb"
	runC             = "cccccccc-cccc-4ccc-8ccc-cccccccccccc"
	runD             = "dddddddd-dddd-4ddd-8ddd-dddddddddddd"
)

// fromRunServer fakes the dataset, example and run endpoints that
// `examples create --from-run` uses and records the create requests.
type fromRunServer struct {
	runs     map[string]map[string]any
	existing []map[string]any
	single   []map[string]any
	bulk     [][]map[string]any
	// rawBodies holds every create request body as sent.
	rawBodies []string
	// dataType is the dataset's data_type ("" for kv).
	dataType string
	// v1Disabled answers v1 single-run reads with 501, as a deployment with
	// ClickHouse queries disabled does; runs are then read from v2 in
	// fromRunSession.
	v1Disabled bool
	// converted maps a source run ID to the inputs /examples/validate/bulk
	// returns for it.
	converted map[string]any
}

func (f *fromRunServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		path := r.URL.Path
		switch {
		case r.Method == http.MethodGet && path == "/api/v1/datasets/"+fromRunDatasetID:
			ds := map[string]any{"id": fromRunDatasetID, "name": "regressions"}
			if f.dataType != "" {
				ds["data_type"] = f.dataType
			}
			_ = json.NewEncoder(w).Encode(ds)
		case r.Method == http.MethodGet && path == "/api/v1/examples":
			if r.URL.Query().Get("dataset") != fromRunDatasetID {
				t.Errorf("examples listed for dataset %q", r.URL.Query().Get("dataset"))
			}
			if off := r.URL.Query().Get("offset"); off != "" && off != "0" {
				_, _ = w.Write([]byte("[]"))
				return
			}
			_ = json.NewEncoder(w).Encode(f.existing)
		case r.Method == http.MethodGet && strings.HasPrefix(path, "/api/v1/runs/") && f.v1Disabled:
			w.WriteHeader(http.StatusNotImplemented)
			_, _ = w.Write([]byte(`{"detail":"Single run read requires ClickHouse query support, which is disabled."}`))
		case r.Method == http.MethodGet && strings.HasPrefix(path, "/api/v2/runs/"):
			if got := r.URL.Query().Get("project_id"); got != fromRunSession {
				t.Errorf("v2 run read for project %q, want %s", got, fromRunSession)
			}
			run, ok := f.runs[strings.TrimPrefix(path, "/api/v2/runs/")]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			v2 := map[string]any{"project_id": run["session_id"]}
			for _, k := range []string{"id", "trace_id", "start_time", "inputs", "outputs"} {
				v2[k] = run[k]
			}
			if slices.Contains(r.URL.Query()["selects"], "ATTACHMENTS") {
				v2["attachments"] = run["attachments"]
			}
			_ = json.NewEncoder(w).Encode(v2)
		case r.Method == http.MethodPost && path == "/api/v1/examples/validate/bulk":
			var body []map[string]any
			decodeBody(t, r.Body, &body)
			resp := make([]map[string]any, len(body))
			for i, b := range body {
				inputs, ok := f.converted[b["source_run_id"].(string)]
				if !ok {
					t.Errorf("run %v sent for conversion", b["source_run_id"])
				}
				resp[i] = map[string]any{"id": b["id"], "inputs": inputs}
			}
			_ = json.NewEncoder(w).Encode(resp)
		case r.Method == http.MethodGet && strings.HasPrefix(path, "/api/v1/runs/"):
			run, ok := f.runs[strings.TrimPrefix(path, "/api/v1/runs/")]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				_, _ = w.Write([]byte(`{"detail":"Run not found"}`))
				return
			}
			_ = json.NewEncoder(w).Encode(run)
		case r.Method == http.MethodPost && path == "/api/v1/examples":
			var body map[string]any
			decodeBody(t, f.record(r.Body), &body)
			f.single = append(f.single, body)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": body["id"], "dataset_id": body["dataset_id"]})
		case r.Method == http.MethodPost && path == "/api/v1/examples/bulk":
			var body []map[string]any
			decodeBody(t, f.record(r.Body), &body)
			f.bulk = append(f.bulk, body)
			resp := make([]map[string]any, len(body))
			for i, b := range body {
				resp[i] = map[string]any{"id": b["id"], "dataset_id": b["dataset_id"]}
			}
			_ = json.NewEncoder(w).Encode(resp)
		default:
			t.Errorf("unexpected request %s %s", r.Method, path)
			w.WriteHeader(http.StatusNotFound)
		}
	}
}

// record keeps a create request body as sent and returns a reader over it.
func (f *fromRunServer) record(r io.Reader) io.Reader {
	b, _ := io.ReadAll(r)
	f.rawBodies = append(f.rawBodies, string(b))
	return bytes.NewReader(b)
}

// decodeBody decodes numbers as json.Number so large integers survive.
func decodeBody(t *testing.T, r io.Reader, v any) {
	t.Helper()
	dec := json.NewDecoder(r)
	dec.UseNumber()
	if err := dec.Decode(v); err != nil {
		t.Errorf("decoding request body: %v", err)
	}
}

func fakeRun(id string, inputs, outputs map[string]any) map[string]any {
	return map[string]any{
		"id": id, "name": "agent", "run_type": "chain", "trace_id": id,
		"session_id": fromRunSession, "start_time": "2026-11-01T10:00:00Z",
		"inputs": inputs, "outputs": outputs,
	}
}

// runFromRunCreate runs `examples create` with JSON output and returns stdout.
func runFromRunCreate(t *testing.T, srv *fromRunServer, stdin string, set map[string][]string) (string, error) {
	t.Helper()
	ts := newTestServer(t, srv.handler(t))
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	cmd := newExampleCreateCmd()
	cmd.SetIn(strings.NewReader(stdin))
	for name, values := range set {
		for _, v := range values {
			if err := cmd.Flags().Set(name, v); err != nil {
				t.Fatalf("setting --%s: %v", name, err)
			}
		}
	}
	var runErr error
	out := captureStdout(t, func() { runErr = runTestCommand(t, cmd, nil) })
	return out, runErr
}

func TestExampleCreateFromRuns_BulkSkipsDuplicateInputs(t *testing.T) {
	srv := &fromRunServer{
		runs: map[string]map[string]any{
			runA: fakeRun(runA, map[string]any{"q": "new", "opts": map[string]any{"b": 2, "a": 1}}, map[string]any{"answer": "x"}),
			// Same inputs as an existing example, with keys in another order.
			runB: fakeRun(runB, map[string]any{"opts": map[string]any{"a": 1, "b": 2}, "q": "old"}, nil),
			runC: fakeRun(runC, map[string]any{"q": "other"}, map[string]any{"answer": "y"}),
			// Same inputs as runA, earlier in this call.
			runD: fakeRun(runD, map[string]any{"opts": map[string]any{"a": 1, "b": 2}, "q": "new"}, nil),
		},
		existing: []map[string]any{{"id": "e1", "inputs": map[string]any{"q": "old", "opts": map[string]any{"a": 1, "b": 2}}}},
	}

	out, err := runFromRunCreate(t, srv, runC+"\n\n"+runD+"\n", map[string][]string{
		"dataset":         {fromRunDatasetID},
		"from-run":        {runA, runB, "-"},
		"skip-duplicates": {"true"},
		"split":           {"regression"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(srv.single) != 0 || len(srv.bulk) != 1 {
		t.Fatalf("want one bulk request, got single=%d bulk=%d", len(srv.single), len(srv.bulk))
	}
	sent := srv.bulk[0]
	if len(sent) != 2 || sent[0]["source_run_id"] != runA || sent[1]["source_run_id"] != runC {
		t.Fatalf("bulk body = %v, want examples from %s and %s", sent, runA, runC)
	}
	first := sent[0]
	if first["dataset_id"] != fromRunDatasetID || first["source_session_id"] != fromRunSession || first["split"] != "regression" {
		t.Errorf("example body = %v", first)
	}
	if got, _ := json.Marshal(first["inputs"]); string(got) != `{"opts":{"a":1,"b":2},"q":"new"}` {
		t.Errorf("inputs = %s", got)
	}
	if got, _ := json.Marshal(first["outputs"]); string(got) != `{"answer":"x"}` {
		t.Errorf("outputs = %s", got)
	}

	var result struct {
		Created       int      `json:"created"`
		Skipped       int      `json:"skipped"`
		ExampleIDs    []string `json:"example_ids"`
		SkippedRunIDs []string `json:"skipped_run_ids"`
	}
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("output %q: %v", out, err)
	}
	if result.Created != 2 || result.Skipped != 2 {
		t.Errorf("created=%d skipped=%d, want 2 and 2", result.Created, result.Skipped)
	}
	if len(result.ExampleIDs) != 2 || result.ExampleIDs[0] != sent[0]["id"] || result.ExampleIDs[1] != sent[1]["id"] {
		t.Errorf("example_ids = %v, want the IDs sent %v, %v", result.ExampleIDs, sent[0]["id"], sent[1]["id"])
	}
	if strings.Join(result.SkippedRunIDs, ",") != runB+","+runD {
		t.Errorf("skipped_run_ids = %v, want [%s %s]", result.SkippedRunIDs, runB, runD)
	}
}

func TestExampleCreateFromRuns_SingleRunUsesCreateEndpoint(t *testing.T) {
	srv := &fromRunServer{runs: map[string]map[string]any{
		runA: fakeRun(runA, map[string]any{"q": "hi"}, map[string]any{"a": "hello"}),
	}}
	out, err := runFromRunCreate(t, srv, "", map[string][]string{
		"dataset":  {fromRunDatasetID},
		"from-run": {runA},
		"metadata": {`{"incident":"INC-1"}`},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(srv.single) != 1 || len(srv.bulk) != 0 {
		t.Fatalf("want one single create, got single=%d bulk=%d", len(srv.single), len(srv.bulk))
	}
	body := srv.single[0]
	if body["source_run_id"] != runA || body["source_run_start_time"] == nil {
		t.Errorf("create body = %v", body)
	}
	if md, _ := body["metadata"].(map[string]any); md["incident"] != "INC-1" {
		t.Errorf("metadata = %v", body["metadata"])
	}
	if !strings.Contains(out, `"created": 1`) || !strings.Contains(out, `"skipped_run_ids": []`) {
		t.Errorf("output = %s", out)
	}
}

func TestExampleCreateFromRuns_MissingRunCreatesNothing(t *testing.T) {
	srv := &fromRunServer{runs: map[string]map[string]any{
		runA: fakeRun(runA, map[string]any{"q": "hi"}, nil),
	}}
	_, err := runFromRunCreate(t, srv, "", map[string][]string{
		"dataset":  {fromRunDatasetID},
		"from-run": {runA, runB},
	})
	if err == nil || !strings.Contains(err.Error(), runB) {
		t.Fatalf("error = %v, want one naming %s", err, runB)
	}
	if len(srv.single)+len(srv.bulk) != 0 {
		t.Errorf("examples were created despite the missing run")
	}
}

func TestExampleCreateFromRuns_ManyRunsKeepInputOrder(t *testing.T) {
	srv := &fromRunServer{runs: map[string]map[string]any{}}
	var ids []string
	for i := range 3 * sourceRunFetchConcurrency {
		id := fmt.Sprintf("eeeeeeee-eeee-4eee-8eee-%012d", i)
		srv.runs[id] = fakeRun(id, map[string]any{"i": i}, nil)
		ids = append(ids, id)
	}
	_, err := runFromRunCreate(t, srv, strings.Join(ids, "\n"), map[string][]string{
		"dataset":  {fromRunDatasetID},
		"from-run": {"-"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(srv.bulk) != 1 || len(srv.bulk[0]) != len(ids) {
		t.Fatalf("want one bulk request with %d examples, got %v", len(ids), srv.bulk)
	}
	for i, ex := range srv.bulk[0] {
		if ex["source_run_id"] != ids[i] {
			t.Fatalf("example %d is from run %v, want %s", i, ex["source_run_id"], ids[i])
		}
	}
}

func TestExampleCreate_FromRunFlagErrors(t *testing.T) {
	tests := []struct {
		name  string
		set   map[string][]string
		stdin string
		want  string
	}{
		{"with inputs", map[string][]string{"from-run": {runA}, "inputs": {`{}`}}, "", "can't be combined with --inputs or --outputs"},
		{"with outputs", map[string][]string{"from-run": {runA}, "outputs": {`{}`}}, "", "can't be combined with --inputs or --outputs"},
		{"neither", nil, "", "one of --inputs or --from-run is required"},
		{"skip without from-run", map[string][]string{"inputs": {`{}`}, "skip-duplicates": {"true"}}, "", "--skip-duplicates requires --from-run"},
		{"bad run id", map[string][]string{"from-run": {"not-a-uuid"}}, "", `invalid run ID "not-a-uuid"`},
		{"empty stdin", map[string][]string{"from-run": {"-"}}, "\n", "--from-run got no run IDs"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cleanup := setupTestEnv(t, "http://127.0.0.1:1")
			defer cleanup()
			cmd := newExampleCreateCmd()
			cmd.SetIn(strings.NewReader(tt.stdin))
			_ = cmd.Flags().Set("dataset", "anything")
			for name, values := range tt.set {
				for _, v := range values {
					_ = cmd.Flags().Set(name, v)
				}
			}
			err := runTestCommand(t, cmd, nil)
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("error = %v, want substring %q", err, tt.want)
			}
		})
	}
}

func TestExampleCreateFromRuns_KeepsExactTimestampsAndIntegers(t *testing.T) {
	run := fakeRun(runA, map[string]any{"record_id": int64(9007199254740993)}, map[string]any{"total": int64(9007199254740995)})
	run["start_time"] = "2026-11-01T10:00:00.123456Z"
	srv := &fromRunServer{
		runs: map[string]map[string]any{runA: run},
		// Equal to the run's inputs only after rounding through float64.
		existing: []map[string]any{{"id": "e1", "inputs": map[string]any{"record_id": int64(9007199254740992)}}},
	}
	if _, err := runFromRunCreate(t, srv, "", map[string][]string{
		"dataset":         {fromRunDatasetID},
		"from-run":        {runA},
		"skip-duplicates": {"true"},
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(srv.rawBodies) != 1 {
		t.Fatalf("want one create request, got %d", len(srv.rawBodies))
	}
	sent := srv.rawBodies[0]
	for _, want := range []string{
		`"source_run_start_time":"2026-11-01T10:00:00.123456Z"`,
		`"inputs":{"record_id":9007199254740993}`,
		`"outputs":{"total":9007199254740995}`,
	} {
		if !strings.Contains(sent, want) {
			t.Errorf("create body %s\nmissing %s", sent, want)
		}
	}
}

func TestExampleCreateFromRuns_ReadsRunsFromV2WithProject(t *testing.T) {
	run := fakeRun(runA, map[string]any{"q": "hi"}, map[string]any{"a": "hello"})
	run["attachments"] = map[string]any{"photo.png": "https://example.com/photo.png"}
	srv := &fromRunServer{
		v1Disabled: true,
		runs:       map[string]map[string]any{runA: run},
	}
	if _, err := runFromRunCreate(t, srv, "", map[string][]string{
		"dataset":    {fromRunDatasetID},
		"from-run":   {runA},
		"project-id": {fromRunSession},
	}); err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(srv.single) != 1 {
		t.Fatalf("want one single create, got %d", len(srv.single))
	}
	body := srv.single[0]
	if body["source_run_id"] != runA || body["source_session_id"] != fromRunSession || body["source_run_start_time"] != "2026-11-01T10:00:00Z" {
		t.Errorf("create body = %v", body)
	}
	if got, _ := json.Marshal(body["inputs"]); string(got) != `{"q":"hi"}` {
		t.Errorf("inputs = %s", got)
	}
	if got, _ := json.Marshal(body["use_source_run_attachments"]); string(got) != `["photo.png"]` {
		t.Errorf("use_source_run_attachments = %s, want the run's attachment", got)
	}

	// Without a project there is no v2 read to fall back to.
	srv = &fromRunServer{v1Disabled: true, runs: srv.runs}
	_, err := runFromRunCreate(t, srv, "", map[string][]string{
		"dataset":  {fromRunDatasetID},
		"from-run": {runA},
	})
	if err == nil || !strings.Contains(err.Error(), "pass --project or --project-id") {
		t.Fatalf("error = %v, want a hint to pass --project", err)
	}
}

func TestExampleCreateFromRuns_ChatDatasetConvertsOnlyModelRuns(t *testing.T) {
	chatInputs := func(text string) map[string]any {
		return map[string]any{"messages": []any{[]any{map[string]any{"type": "human", "content": text}}}}
	}
	stored := func(text string) map[string]any {
		return map[string]any{"input": []any{map[string]any{"type": "human", "data": map[string]any{"content": text}}}}
	}
	formatted := fakeRun(runC, stored("formatted"), map[string]any{"output": "ok"})
	formatted["s3_urls"] = map[string]any{
		"attachment.document.pdf": map[string]any{"s3_url": "s3://bucket/document.pdf"},
		"inputs":                  map[string]any{"s3_url": "s3://bucket/inputs"},
	}
	srv := &fromRunServer{
		dataType: "chat",
		runs: map[string]map[string]any{
			runA: fakeRun(runA, chatInputs("new"), map[string]any{"generations": []any{}}),
			runB: fakeRun(runB, chatInputs("old"), nil),
			// Already in the dataset's format, so copied as is, keeping its attachment.
			runC: formatted,
		},
		converted: map[string]any{runA: stored("new"), runB: stored("old")},
		existing:  []map[string]any{{"id": "e1", "inputs": stored("old")}},
	}
	out, err := runFromRunCreate(t, srv, "", map[string][]string{
		"dataset":         {fromRunDatasetID},
		"from-run":        {runA, runB, runC},
		"skip-duplicates": {"true"},
	})
	if err != nil {
		t.Fatalf("create: %v", err)
	}
	if len(srv.bulk) != 1 || len(srv.bulk[0]) != 2 {
		t.Fatalf("want one bulk create of 2, got single=%d bulk=%v", len(srv.single), srv.bulk)
	}
	converted, copied := srv.bulk[0][0], srv.bulk[0][1]
	if converted["source_run_id"] != runA || converted["use_source_run_io"] != true {
		t.Errorf("create body = %v, want %s with use_source_run_io", converted, runA)
	}
	if _, ok := converted["inputs"]; ok {
		t.Errorf("create body sends raw model-run inputs to a chat dataset: %v", converted["inputs"])
	}
	if copied["source_run_id"] != runC || copied["use_source_run_io"] != nil || copied["inputs"] == nil {
		t.Errorf("create body = %v, want %s copied without conversion", copied, runC)
	}
	if got, _ := json.Marshal(copied["use_source_run_attachments"]); string(got) != `["document.pdf"]` {
		t.Errorf("use_source_run_attachments = %s, want the run's attachment", got)
	}
	if !strings.Contains(out, runB) {
		t.Errorf("output = %s, want %s skipped", out, runB)
	}
}
