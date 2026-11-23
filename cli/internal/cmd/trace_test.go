package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"testing"
	"time"

	langsmith "github.com/langchain-ai/langsmith-go"
)

func TestTraceExportPathUsesTraceID(t *testing.T) {
	first := langsmith.RunSchema{TraceID: "trace-1", Name: "same"}
	second := langsmith.RunSchema{TraceID: "trace-2", Name: "same"}
	if got := traceExportPath("out", "{trace_id}.jsonl", first); got == traceExportPath("out", "{trace_id}.jsonl", second) {
		t.Fatal("trace_id pattern should produce unique paths")
	}
}

func TestTraceExportPathCollidesForConstantAndNamePatterns(t *testing.T) {
	first := langsmith.RunSchema{TraceID: "trace-1", Name: "same"}
	second := langsmith.RunSchema{TraceID: "trace-2", Name: "same"}
	for _, pattern := range []string{"trace.jsonl", "{name}.jsonl"} {
		if traceExportPath("out", pattern, first) != traceExportPath("out", pattern, second) {
			t.Errorf("pattern %q should collide for equal names", pattern)
		}
	}
}

// ==================== Command structure ====================

func TestTraceCmd_Subcommands(t *testing.T) {
	cmd := newTraceCmd()
	expected := map[string]bool{"list": false, "get": false, "export": false, "messages": false}
	for _, sub := range cmd.Commands() {
		if _, ok := expected[sub.Name()]; ok {
			expected[sub.Name()] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("trace missing subcommand %q", name)
		}
	}
}

func TestTraceCmd_UseField(t *testing.T) {
	cmd := newTraceCmd()
	if cmd.Use != "traces" {
		t.Errorf("expected Use=traces, got %q", cmd.Use)
	}
}

// ==================== trace list flags ====================

func TestTraceListCmd_Flags(t *testing.T) {
	cmd := newTraceListCmd()
	tests := []struct {
		name   string
		defVal string
		short  string
	}{
		{"include-metadata", "false", ""},
		{"include-io", "false", ""},
		{"full", "false", ""},
		{"show-hierarchy", "false", ""},
		{"output", "", "o"},
	}
	for _, tc := range tests {
		f := cmd.Flags().Lookup(tc.name)
		if f == nil {
			t.Errorf("flag --%s not found", tc.name)
			continue
		}
		if f.DefValue != tc.defVal {
			t.Errorf("flag --%s: expected default %q, got %q", tc.name, tc.defVal, f.DefValue)
		}
		if tc.short != "" && f.Shorthand != tc.short {
			t.Errorf("flag --%s: expected shorthand %q, got %q", tc.name, tc.short, f.Shorthand)
		}
	}
}

func TestTraceListCmd_HasCommonFilterFlags(t *testing.T) {
	cmd := newTraceListCmd()
	common := []string{"trace-ids", "limit", "project", "last-n-minutes", "since",
		"error", "no-error", "name", "min-latency", "max-latency", "min-tokens", "tags", "filter"}
	for _, name := range common {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("trace list missing common filter flag --%s", name)
		}
	}
}

func TestTraceListCmd_NoRunTypeFlag(t *testing.T) {
	cmd := newTraceListCmd()
	if cmd.Flags().Lookup("run-type") != nil {
		t.Error("trace list should not have --run-type flag")
	}
}

// ==================== trace get flags ====================

func TestTraceGetCmd_Flags(t *testing.T) {
	cmd := newTraceGetCmd()
	tests := []struct {
		name   string
		defVal string
	}{
		{"project", ""},
		{"include-metadata", "false"},
		{"include-io", "false"},
		{"full", "false"},
		{"output", ""},
	}
	for _, tc := range tests {
		f := cmd.Flags().Lookup(tc.name)
		if f == nil {
			t.Errorf("flag --%s not found", tc.name)
			continue
		}
		if f.DefValue != tc.defVal {
			t.Errorf("flag --%s: expected default %q, got %q", tc.name, tc.defVal, f.DefValue)
		}
	}
}

func TestTraceGetCmd_OutputShorthand(t *testing.T) {
	cmd := newTraceGetCmd()
	f := cmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("--output flag not found")
	}
	if f.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got %q", f.Shorthand)
	}
}

func TestTraceGetCmd_TimeFlags(t *testing.T) {
	cmd := newTraceGetCmd()
	for _, name := range []string{"since", "last-n-minutes"} {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("trace get missing flag --%s", name)
		}
	}
}

func TestTraceGet_V1WithoutProjectQueriesByTraceOnly(t *testing.T) {
	withDeploymentVersion(t, "0.15.0")
	t.Setenv("LANGSMITH_PROJECT", "")
	const traceID = "01a10aa3-0e17-7eb6-b255-d8a470732c91"
	var paths []string
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		_ = json.NewDecoder(r.Body).Decode(&body)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"runs":[{"id":"r1","trace_id":"` + traceID + `"}],"cursors":{}}`))
	})

	got := runTraceCmdJSON(t, ts.URL, "get", traceID).(map[string]any)

	if len(paths) != 1 || paths[0] != "/api/v1/runs/query" {
		t.Fatalf("requests = %v, want one /api/v1/runs/query", paths)
	}
	if body["trace"] != traceID {
		t.Errorf("trace = %v, want %s", body["trace"], traceID)
	}
	if _, ok := body["session"]; ok {
		t.Errorf("sent a session filter without --project: %v", body["session"])
	}
	if _, ok := body["start_time"]; !ok {
		t.Error("start_time missing: the v1 lookup should keep its default window")
	}
	if got["run_count"] != float64(1) {
		t.Errorf("run_count = %v, want 1", got["run_count"])
	}
}

func TestTraceGet_V2WithoutProjectFindsProjectByTrace(t *testing.T) {
	t.Setenv("LANGSMITH_PROJECT", "")
	const traceID = "01a10aa3-0e17-7eb6-b255-d8a470732c91"
	const staleProjectID = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9328"
	recent := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	stale := time.Now().Add(-30 * 24 * time.Hour).UTC().Format(time.RFC3339)
	var paths []string
	var sessionsQuery url.Values
	var lookup map[string]any
	var traceQuery map[string][]string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		paths = append(paths, r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/sessions":
			sessionsQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"items":[` +
				`{"id":"` + testProjectID + `","last_run_start_time":"` + recent + `"},` +
				`{"id":"` + staleProjectID + `","last_run_start_time":"` + stale + `"}` +
				`],"next_cursor":"more"}`))
		case "/api/v2/runs/query":
			_ = json.NewDecoder(r.Body).Decode(&lookup)
			_, _ = w.Write([]byte(`{"items":[{"id":"r1","project_id":"` + testProjectID + `"}]}`))
		case "/api/v2/traces/" + traceID + "/runs":
			traceQuery = r.URL.Query()
			_, _ = w.Write([]byte(`{"items":[{"id":"r1","start_time":"2026-01-01T00:00:01Z"}]}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
		}
	})

	got := runTraceCmdJSON(t, ts.URL, "get", traceID).(map[string]any)

	want := []string{"/api/v2/sessions", "/api/v2/runs/query", "/api/v2/traces/" + traceID + "/runs"}
	if strings.Join(paths, " ") != strings.Join(want, " ") {
		t.Fatalf("requests = %v, want %v", paths, want)
	}
	if sessionsQuery.Get("sort_by") != "LAST_RUN_START_TIME" || sessionsQuery.Get("sort_order") != "desc" {
		t.Errorf("sessions sort = %v %v, want most recently active first", sessionsQuery["sort_by"], sessionsQuery["sort_order"])
	}
	if types := sessionsQuery["type"]; strings.Join(types, ",") != "PROJECT,AGENT" {
		t.Errorf("sessions type = %v, want tracing projects only", types)
	}
	// Only projects with a run in the window are searched; the stale one ends
	// the listing, so the next page is never requested.
	if ids, _ := lookup["project_ids"].([]any); len(ids) != 1 || ids[0] != testProjectID {
		t.Errorf("lookup project_ids = %v, want only the recently active project", lookup["project_ids"])
	}
	if lookup["trace_id"] != traceID || lookup["is_root"] != true {
		t.Errorf("lookup trace_id = %v, is_root = %v, want the trace's root run", lookup["trace_id"], lookup["is_root"])
	}
	if _, ok := lookup["min_start_time"]; !ok {
		t.Error("lookup min_start_time missing: it should keep the default window")
	}
	if traceQuery["project_id"][0] != testProjectID {
		t.Errorf("project_id = %v, want the project from the lookup", traceQuery["project_id"])
	}
	if got["run_count"] != float64(1) {
		t.Errorf("run_count = %v, want 1", got["run_count"])
	}
}

func TestTraceGet_V2LookupPagesThroughActiveProjects(t *testing.T) {
	t.Setenv("LANGSMITH_PROJECT", "")
	const traceID = "01a10aa3-0e17-7eb6-b255-d8a470732c91"
	const firstPageProjectID = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9328"
	recent := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	var cursors []string
	var searched [][]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/sessions":
			cursor := r.URL.Query().Get("cursor")
			cursors = append(cursors, cursor)
			if cursor == "" {
				_, _ = w.Write([]byte(`{"items":[{"id":"` + firstPageProjectID + `","last_run_start_time":"` + recent + `"}],"next_cursor":"page-2"}`))
				return
			}
			_, _ = w.Write([]byte(`{"items":[{"id":"` + testProjectID + `","last_run_start_time":"` + recent + `"}]}`))
		case "/api/v2/runs/query":
			var body map[string]any
			_ = json.NewDecoder(r.Body).Decode(&body)
			ids, _ := body["project_ids"].([]any)
			searched = append(searched, ids)
			if len(ids) == 1 && ids[0] == testProjectID {
				_, _ = w.Write([]byte(`{"items":[{"id":"r1","project_id":"` + testProjectID + `"}]}`))
				return
			}
			_, _ = w.Write([]byte(`{"items":[]}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	got, err := traceProjectID(context.Background(), MustGetClient(), traceID, time.Now().Add(-7*24*time.Hour))
	if err != nil {
		t.Fatal(err)
	}
	if got != testProjectID {
		t.Errorf("project = %q, want %q from the second page", got, testProjectID)
	}
	if strings.Join(cursors, ",") != ",page-2" {
		t.Errorf("sessions cursors = %q, want the first page then page-2", cursors)
	}
	if len(searched) != 2 {
		t.Errorf("runs queries = %v, want one per page of projects", searched)
	}
}

// The project activity time lags new runs and is empty for new projects and
// experiments, and the runs query rejects a batch holding any project whose
// runs the caller can't read. None of that may hide an accessible trace.
func TestTraceProjectID_LaggingActivityAndDeniedProjects(t *testing.T) {
	const (
		traceID    = "01a10aa3-0e17-7eb6-b255-d8a470732c91"
		laggingID  = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9301"
		deniedID   = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9302"
		staleID    = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9303"
		newID      = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9304"
		experiment = "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9305"
	)
	start := time.Now().Add(-7 * 24 * time.Hour)
	at := func(t time.Time) string { return t.UTC().Format(time.RFC3339) }
	var types []string
	var searched []string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/sessions":
			q := r.URL.Query()
			types = append(types, strings.Join(q["type"], ","))
			if q.Get("sort_by") == "LAST_RUN_START_TIME" {
				_, _ = w.Write([]byte(`{"items":[` +
					`{"id":"` + laggingID + `","last_run_start_time":"` + at(start.Add(-30*time.Minute)) + `"},` +
					`{"id":"` + deniedID + `","last_run_start_time":"` + at(start.Add(time.Hour)) + `"},` +
					`{"id":"` + staleID + `","last_run_start_time":"` + at(start.Add(-48*time.Hour)) + `"}` +
					`],"next_cursor":"more"}`))
				return
			}
			_, _ = w.Write([]byte(`{"items":[` +
				`{"id":"` + newID + `","start_time":"` + at(time.Now()) + `"},` +
				`{"id":"` + laggingID + `","start_time":"` + at(start) + `","last_run_start_time":"` + at(start) + `"},` +
				`{"id":"` + experiment + `","start_time":"` + at(start.Add(time.Hour)) + `"},` +
				`{"id":"` + staleID + `","start_time":"` + at(start.Add(-48*time.Hour)) + `"}` +
				`],"next_cursor":"more"}`))
		case "/api/v2/runs/query":
			var body struct {
				ProjectIDs []string `json:"project_ids"`
			}
			_ = json.NewDecoder(r.Body).Decode(&body)
			searched = append(searched, strings.Join(body.ProjectIDs, "+"))
			if slices.Contains(body.ProjectIDs, deniedID) {
				http.Error(w, `{"detail":"forbidden"}`, http.StatusForbidden)
				return
			}
			if slices.Contains(body.ProjectIDs, experiment) {
				_, _ = w.Write([]byte(`{"items":[{"id":"r1","project_id":"` + experiment + `"}]}`))
				return
			}
			_, _ = w.Write([]byte(`{"items":[]}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	got, err := traceProjectID(context.Background(), MustGetClient(), traceID, start)
	if err != nil {
		t.Fatal(err)
	}
	if got != experiment {
		t.Errorf("project = %q, want the experiment %q", got, experiment)
	}
	if want := []string{"PROJECT,AGENT", "PROJECT,AGENT,EXPERIMENT"}; !slices.Equal(types, want) {
		t.Errorf("listed types = %v, want %v", types, want)
	}
	want := []string{
		laggingID + "+" + deniedID, // within the lag margin; rejected for deniedID
		laggingID, deniedID,        // split; deniedID skipped
		newID + "+" + experiment, // no recorded activity yet
	}
	if !slices.Equal(searched, want) {
		t.Errorf("runs queries = %v, want %v", searched, want)
	}
}

func TestTraceGetRuns_NotFound(t *testing.T) {
	t.Setenv("LANGSMITH_PROJECT", "")
	const traceID = "01a10aa3-0e17-7eb6-b255-d8a470732c91"
	cases := []struct {
		name, version, projectID string
		want, lookup             []string
	}{
		{"v2 workspace", "dev", "", []string{"not found in the current workspace", "--since"},
			[]string{"/api/v2/sessions", "/api/v2/sessions"}},
		// v2 runs API without the v2 project list: the lookup stays on v1.
		{"v2 workspace before sessions list", "0.18.22", "", []string{"not found in the current workspace", "--since"},
			[]string{"/api/v1/runs/query"}},
		{"v1 workspace", "0.15.0", "", []string{"not found in the current workspace", "--since"},
			[]string{"/api/v1/runs/query"}},
		{"v2 project", "dev", testProjectID, []string{"not found in the given project"}, nil},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			withDeploymentVersion(t, tc.version)
			var paths []string
			ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				paths = append(paths, r.URL.Path)
				w.Header().Set("Content-Type", "application/json")
				if strings.HasPrefix(r.URL.Path, "/api/v2/") {
					_, _ = w.Write([]byte(`{"items":[]}`))
					return
				}
				_, _ = w.Write([]byte(`{"runs":[],"cursors":{}}`))
			})
			cleanup := setupTestEnv(t, ts.URL)
			defer cleanup()

			_, err := traceGetRuns(context.Background(), MustGetClient(), "", tc.projectID, traceID, "", 0, nil, nil)
			if err == nil {
				t.Fatal("want a not-found error")
			}
			for _, w := range tc.want {
				if !strings.Contains(err.Error(), w) {
					t.Errorf("err = %q, want it to mention %q", err, w)
				}
			}
			if tc.projectID == "" && strings.Join(paths, ",") != strings.Join(tc.lookup, ",") {
				t.Errorf("requests = %v, want only the workspace lookup %v", paths, tc.lookup)
			}
		})
	}
}

func TestTraceGetCmd_ExactArgs(t *testing.T) {
	cmd := newTraceGetCmd()
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error for 0 args")
	}
	if err := cmd.Args(cmd, []string{"trace-123"}); err != nil {
		t.Errorf("expected no error for 1 arg, got %v", err)
	}
	if err := cmd.Args(cmd, []string{"a", "b"}); err == nil {
		t.Error("expected error for 2 args")
	}
}

// ==================== trace export flags ====================

func TestTraceExportCmd_Flags(t *testing.T) {
	cmd := newTraceExportCmd()
	tests := []struct {
		name   string
		defVal string
	}{
		{"include-metadata", "false"},
		{"include-io", "false"},
		{"full", "false"},
		{"filename-pattern", "{trace_id}.jsonl"},
	}
	for _, tc := range tests {
		f := cmd.Flags().Lookup(tc.name)
		if f == nil {
			t.Errorf("flag --%s not found", tc.name)
			continue
		}
		if f.DefValue != tc.defVal {
			t.Errorf("flag --%s: expected default %q, got %q", tc.name, tc.defVal, f.DefValue)
		}
	}
}

func TestTraceExportCmd_HasCommonFilterFlags(t *testing.T) {
	cmd := newTraceExportCmd()
	common := []string{"trace-ids", "limit", "project", "last-n-minutes", "since",
		"error", "no-error", "name", "min-latency", "max-latency", "min-tokens", "tags", "filter"}
	for _, name := range common {
		if cmd.Flags().Lookup(name) == nil {
			t.Errorf("trace export missing common filter flag --%s", name)
		}
	}
}

func TestTraceExportCmd_NoRunTypeFlag(t *testing.T) {
	cmd := newTraceExportCmd()
	if cmd.Flags().Lookup("run-type") != nil {
		t.Error("trace export should not have --run-type flag")
	}
}

func TestTraceExportCmd_ExactArgs(t *testing.T) {
	cmd := newTraceExportCmd()
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error for 0 args")
	}
	if err := cmd.Args(cmd, []string{"./output"}); err != nil {
		t.Errorf("expected no error for 1 arg, got %v", err)
	}
	if err := cmd.Args(cmd, []string{"a", "b"}); err == nil {
		t.Error("expected error for 2 args")
	}
}

func TestTraceGetRuns_V2LookupBoundedAsksForProject(t *testing.T) {
	t.Setenv("LANGSMITH_PROJECT", "")
	withDeploymentVersion(t, "dev")
	const traceID = "01a10aa3-0e17-7eb6-b255-d8a470732c91"
	recent := time.Now().Add(-time.Hour).UTC().Format(time.RFC3339)
	var lists, queries int
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/api/v2/sessions":
			lists++
			id := fmt.Sprintf("018e4c7e-a9fb-7ef0-a5b6-%012d", lists)
			_, _ = w.Write([]byte(`{"items":[{"id":"` + id + `","last_run_start_time":"` + recent + `"}],"next_cursor":"next"}`))
		case "/api/v2/runs/query":
			queries++
			_, _ = w.Write([]byte(`{"items":[]}`))
		default:
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	_, err := traceGetRuns(context.Background(), MustGetClient(), "", "", traceID, "", 0, nil, nil)
	if err == nil || !strings.Contains(err.Error(), "--project-id") {
		t.Errorf("err = %v, want it to ask for --project or --project-id", err)
	}
	// The active-project pass stops at its page cap; the pass over projects
	// without recorded activity stops at the first one created before the window.
	if lists != traceLookupMaxPages+1 || queries != traceLookupMaxPages {
		t.Errorf("lists = %d, queries = %d, want the lookup to stop after %d pages", lists, queries, traceLookupMaxPages)
	}
}
