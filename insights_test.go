package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

// newInsightsTestClient creates a client pointed at the provided test server.
func newInsightsTestClient(t *testing.T, srv *httptest.Server) *langsmith.Client {
	t.Helper()
	return langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithTenantID("tenant-id"),
	)
}

func TestGenerateInsights(t *testing.T) {
	sessionID := "aaaaaaaa-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	jobID := "bbbbbbbb-bbbb-bbbb-bbbb-bbbbbbbbbbbb"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		// GET workspace secrets
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/workspaces/current/secrets":
			json.NewEncoder(w).Encode([]map[string]string{
				{"secret_name": "OPENAI_API_KEY"},
			})

		// POST create session
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":        sessionID,
				"tenant_id": "tenant-id",
				"name":      "test-insights",
				"start_time": time.Now().Format(time.RFC3339),
			})

		// POST ingest runs
		case r.Method == http.MethodPost && r.URL.Path == "/runs/batch":
			json.NewEncoder(w).Encode(map[string]interface{}{})

		// POST create insights job
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions/"+sessionID+"/insights":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":     jobID,
				"name":   "test-insights",
				"status": "pending",
				"error":  nil,
			})

		default:
			http.Error(w, "unexpected: "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	report, err := client.GenerateInsights(context.Background(), langsmith.GenerateInsightsParams{
		ChatHistories: [][]map[string]interface{}{
			{
				{"role": "user", "content": "hello"},
				{"role": "assistant", "content": "hi!"},
			},
		},
		Name:  "test-insights",
		Model: "openai",
	})
	if err != nil {
		t.Fatalf("GenerateInsights error: %v", err)
	}

	if report.ID != jobID {
		t.Errorf("want ID=%s, got %s", jobID, report.ID)
	}
	if report.ProjectID != sessionID {
		t.Errorf("want ProjectID=%s, got %s", sessionID, report.ProjectID)
	}
	if report.Status != "pending" {
		t.Errorf("want Status=pending, got %s", report.Status)
	}
}

func TestPollInsights(t *testing.T) {
	projectID := "cccccccc-cccc-cccc-cccc-cccccccccccc"
	jobID := "dddddddd-dddd-dddd-dddd-dddddddddddd"

	call := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		call++
		status := "pending"
		if call >= 2 {
			status = "success"
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       jobID,
			"name":     "my-report",
			"status":   status,
			"clusters": []interface{}{},
			"report":   nil,
			"error":    nil,
		})
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	report, err := client.PollInsights(context.Background(), langsmith.PollInsightsParams{
		ProjectID: projectID,
		ID:        jobID,
		Rate:      10 * time.Millisecond,
		Timeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatalf("PollInsights error: %v", err)
	}
	if report.Status != "success" {
		t.Errorf("want Status=success, got %s", report.Status)
	}
}

func TestPollInsights_Timeout(t *testing.T) {
	projectID := "eeeeeeee-eeee-eeee-eeee-eeeeeeeeeeee"
	jobID := "ffffffff-ffff-ffff-ffff-ffffffffffff"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id":       jobID,
			"name":     "slow-report",
			"status":   "pending",
			"clusters": []interface{}{},
			"error":    nil,
		})
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	_, err := client.PollInsights(context.Background(), langsmith.PollInsightsParams{
		ProjectID: projectID,
		ID:        jobID,
		Rate:      5 * time.Millisecond,
		Timeout:   20 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected timeout error, got nil")
	}
}

func TestGetInsightsReport(t *testing.T) {
	projectID := "11111111-1111-1111-1111-111111111111"
	jobID := "22222222-2222-2222-2222-222222222222"
	clusterID := "33333333-3333-3333-3333-333333333333"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/"+projectID+"/insights/"+jobID:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         jobID,
				"name":       "my-report",
				"status":     "success",
				"start_time": "2026-01-01T00:00:00Z",
				"end_time":   "2026-01-01T01:00:00Z",
				"shape":      map[string]int{"cluster-a": 3},
				"clusters": []map[string]interface{}{
					{
						"id":          clusterID,
						"name":        "cluster-a",
						"description": "Cluster A",
						"level":       0,
						"num_runs":    3,
						"stats":       nil,
					},
				},
				"report": map[string]interface{}{
					"key_points":         []string{"users ask about topic X"},
					"title":              "Usage Summary",
					"highlighted_traces": []interface{}{},
				},
				"error": nil,
			})

		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/"+projectID+"/insights/"+jobID+"/runs":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"runs":   []map[string]interface{}{{"id": "run-1"}, {"id": "run-2"}},
				"offset": 0,
			})

		default:
			http.Error(w, "unexpected: "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	result, err := client.GetInsightsReport(context.Background(), langsmith.GetInsightsReportParams{
		ProjectID:   projectID,
		ID:          jobID,
		IncludeRuns: true,
	})
	if err != nil {
		t.Fatalf("GetInsightsReport error: %v", err)
	}

	if result.Status != "success" {
		t.Errorf("want Status=success, got %s", result.Status)
	}
	if len(result.Clusters) != 1 {
		t.Errorf("want 1 cluster, got %d", len(result.Clusters))
	}
	if result.Clusters[0].Name != "cluster-a" {
		t.Errorf("want cluster name=cluster-a, got %s", result.Clusters[0].Name)
	}
	if result.Report == nil {
		t.Fatal("expected non-nil Report")
	}
	if result.Report.Title != "Usage Summary" {
		t.Errorf("want Report.Title=Usage Summary, got %s", result.Report.Title)
	}
	if len(result.Runs) != 2 {
		t.Errorf("want 2 runs, got %d", len(result.Runs))
	}
	if result.StartTime == nil {
		t.Error("expected non-nil StartTime")
	}
}

func TestGetInsightsReport_UsingReport(t *testing.T) {
	projectID := "44444444-4444-4444-4444-444444444444"
	jobID := "55555555-5555-5555-5555-555555555555"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/"+projectID+"/insights/"+jobID:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":       jobID,
				"name":     "report-via-struct",
				"status":   "success",
				"clusters": []interface{}{},
				"error":    nil,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/"+projectID+"/insights/"+jobID+"/runs":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"runs":   []interface{}{},
				"offset": 0,
			})
		default:
			http.Error(w, "unexpected: "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	report := &langsmith.InsightsReport{
		ID:        jobID,
		ProjectID: projectID,
		TenantID:  "tenant-id",
		Status:    "success",
	}

	result, err := client.GetInsightsReport(context.Background(), langsmith.GetInsightsReportParams{
		Report:      report,
		IncludeRuns: true,
	})
	if err != nil {
		t.Fatalf("GetInsightsReport error: %v", err)
	}
	if result.ID != jobID {
		t.Errorf("want ID=%s, got %s", jobID, result.ID)
	}
}
