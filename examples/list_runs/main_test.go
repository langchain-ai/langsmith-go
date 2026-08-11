package main

import (
	"encoding/json"
	"slices"
	"strings"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
)

func TestBuildProjectParams(t *testing.T) {
	params := buildProjectParams("my-project")

	if got := params.Name.Value; got != "my-project" {
		t.Errorf("expected project name my-project, got %q", got)
	}
	// Without upsert the create fails once the project exists, and the example
	// would never get a project ID back for an already-created project.
	if got := params.Upsert.Value; !got {
		t.Error("expected upsert to be true")
	}
}

// TestBuildProjectParamsQuery guards that upsert travels as a query parameter
// rather than in the request body.
func TestBuildProjectParamsQuery(t *testing.T) {
	query := buildProjectParams("my-project").URLQuery()

	if got := query.Get("upsert"); got != "true" {
		t.Errorf("expected upsert=true in the query string, got %q", query.Encode())
	}
}

func TestRunSelectFields(t *testing.T) {
	fields := runSelectFields()

	// QueryV2 returns only `id` unless the other fields are requested, so every
	// property formatRun prints must appear here.
	for _, want := range []langsmith.RunSelectField{
		langsmith.RunSelectFieldName,
		langsmith.RunSelectFieldRunType,
		langsmith.RunSelectFieldStatus,
		langsmith.RunSelectFieldStartTime,
		langsmith.RunSelectFieldLatencySeconds,
	} {
		if !slices.Contains(fields, want) {
			t.Errorf("expected %q in selected fields, got %v", want, fields)
		}
	}
}

// TestPagingBounds guards the invariant that makes the auto-pager request more
// than one page: a page smaller than the print budget.
func TestPagingBounds(t *testing.T) {
	if defaultPageSize >= maxRuns {
		t.Errorf("expected defaultPageSize (%d) below maxRuns (%d) so the auto-pager pages", defaultPageSize, maxRuns)
	}
}

func TestBuildQueryParams(t *testing.T) {
	now := time.Date(2026, time.August, 7, 12, 0, 0, 0, time.UTC)

	params := buildQueryParams("ds-1", now)

	if got := params.ProjectIDs.Value; !slices.Equal(got, []string{"ds-1"}) {
		t.Errorf("expected project IDs [ds-1], got %v", got)
	}
	if got := params.PageSize.Value; got != defaultPageSize {
		t.Errorf("expected page size %d, got %d", defaultPageSize, got)
	}
	if want := now.Add(-defaultLookback); !params.MinStartTime.Value.Equal(want) {
		t.Errorf("expected min start time %s, got %s", want, params.MinStartTime.Value)
	}
	if got := params.Selects.Value; !slices.Equal(got, runSelectFields()) {
		t.Errorf("expected selects %v, got %v", runSelectFields(), got)
	}
}

// TestBuildQueryParamsBody guards the V2 request shape: the legacy v1 query took
// `session`, while /api/v2/runs/query takes `project_ids` and `page_size`.
func TestBuildQueryParamsBody(t *testing.T) {
	now := time.Date(2026, time.August, 7, 12, 0, 0, 0, time.UTC)

	data, err := json.Marshal(buildQueryParams("ds-1", now))
	if err != nil {
		t.Fatalf("marshaling params: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("unmarshaling params: %v", err)
	}

	for _, key := range []string{"project_ids", "page_size", "min_start_time", "selects"} {
		if _, ok := body[key]; !ok {
			t.Errorf("expected %q in request body, got %v", key, body)
		}
	}
	if _, ok := body["session"]; ok {
		t.Errorf("unexpected legacy v1 %q key in request body: %v", "session", body)
	}
}

func TestFormatRun(t *testing.T) {
	started := time.Date(2026, time.August, 7, 9, 30, 0, 0, time.UTC)
	out := formatRun(2, langsmith.Run{
		ID:             "0f8fad5b-d9cb-469f-a165-70867728950e",
		Name:           "my_chain",
		RunType:        langsmith.RunTypeChain,
		Status:         langsmith.RunStatusSuccess,
		StartTime:      started,
		LatencySeconds: 1.234,
	})

	for _, want := range []string{
		"2. Run ID: 0f8fad5b-d9cb-469f-a165-70867728950e",
		"Name: my_chain",
		"Type: CHAIN",
		"Status: SUCCESS",
		"Started: " + started.Format(time.RFC3339),
		"Latency: 1.23s",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("expected %q in output, got:\n%s", want, out)
		}
	}
}

func TestFormatRunOmitsUnsetFields(t *testing.T) {
	out := formatRun(1, langsmith.Run{
		ID:   "0f8fad5b-d9cb-469f-a165-70867728950e",
		Name: "my_chain",
	})

	if !strings.Contains(out, "1. Run ID: 0f8fad5b-d9cb-469f-a165-70867728950e") {
		t.Errorf("expected run ID in output, got:\n%s", out)
	}
	for _, unwanted := range []string{"Type:", "Status:", "Started:", "Latency:"} {
		if strings.Contains(out, unwanted) {
			t.Errorf("expected no %q line for an unset field, got:\n%s", unwanted, out)
		}
	}
}
