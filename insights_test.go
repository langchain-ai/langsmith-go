package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

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
		// GET workspace secrets — returns key field (not secret_name)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/workspaces/current/secrets":
			json.NewEncoder(w).Encode([]map[string]string{
				{"key": "OPENAI_API_KEY"},
			})

		// POST create session
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id":         sessionID,
				"tenant_id":  "tenant-id",
				"name":       "test-insights",
				"start_time": time.Now().Format(time.RFC3339),
			})

		// POST ingest runs — verify run name is "trace"
		case r.Method == http.MethodPost && r.URL.Path == "/runs/batch":
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			posts, _ := body["post"].([]interface{})
			if len(posts) == 0 {
				http.Error(w, "no runs posted", http.StatusBadRequest)
				return
			}
			run, _ := posts[0].(map[string]interface{})
			if run["name"] != "trace" {
				http.Error(w, "expected run name 'trace', got: "+run["name"].(string), http.StatusBadRequest)
				return
			}
			json.NewEncoder(w).Encode(map[string]interface{}{})

		// POST create insights job — verify user_context and last_n_hours are sent
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions/"+sessionID+"/insights":
			var body map[string]interface{}
			json.NewDecoder(r.Body).Decode(&body)
			if body["last_n_hours"] == nil {
				http.Error(w, "expected last_n_hours in body", http.StatusBadRequest)
				return
			}
			if body["user_context"] == nil {
				http.Error(w, "expected user_context in body", http.StatusBadRequest)
				return
			}
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
	// Link should be populated
	if !strings.Contains(report.Link(), report.ID) {
		t.Errorf("Link() should contain job ID, got: %s", report.Link())
	}
}

func TestGenerateInsights_StoresAPIKey(t *testing.T) {
	sessionID := "cccccccc-cccc-cccc-cccc-cccccccccccc"
	jobID := "dddddddd-dddd-dddd-dddd-dddddddddddd"
	storedKey := ""

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		// No existing secrets
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/workspaces/current/secrets":
			json.NewEncoder(w).Encode([]interface{}{})

		// POST secrets — verify array format with key/value fields
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/workspaces/current/secrets":
			var body []map[string]string
			json.NewDecoder(r.Body).Decode(&body)
			if len(body) == 0 {
				http.Error(w, "expected array body", http.StatusBadRequest)
				return
			}
			storedKey = body[0]["key"]
			json.NewEncoder(w).Encode(map[string]interface{}{})

		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": sessionID, "tenant_id": "t", "name": "x",
				"start_time": time.Now().Format(time.RFC3339),
			})

		case r.Method == http.MethodPost && r.URL.Path == "/runs/batch":
			json.NewEncoder(w).Encode(map[string]interface{}{})

		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/sessions/"+sessionID+"/insights":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": jobID, "name": "x", "status": "pending", "error": nil,
			})

		default:
			http.Error(w, "unexpected: "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	_, err := client.GenerateInsights(context.Background(), langsmith.GenerateInsightsParams{
		ChatHistories:   [][]map[string]interface{}{},
		OpenAIAPIKey:    "sk-test",
	})
	if err != nil {
		t.Fatalf("GenerateInsights error: %v", err)
	}
	if storedKey != "OPENAI_API_KEY" {
		t.Errorf("expected OPENAI_API_KEY to be stored, got %q", storedKey)
	}
}

func TestPollInsights_Success(t *testing.T) {
	projectID := "11111111-1111-1111-1111-111111111111"
	jobID := "22222222-2222-2222-2222-222222222222"

	call := 0
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		call++
		status := "pending"
		if call >= 2 {
			status = "success"
		}
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": jobID, "name": "report", "status": status,
			"clusters": []interface{}{}, "error": nil,
		})
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	report, err := client.PollInsights(context.Background(), langsmith.PollInsightsParams{
		ProjectID: projectID,
		ID:        jobID,
		Rate:      5 * time.Millisecond,
		Timeout:   5 * time.Second,
	})
	if err != nil {
		t.Fatalf("PollInsights error: %v", err)
	}
	if report.Status != "success" {
		t.Errorf("want Status=success, got %s", report.Status)
	}
}

func TestPollInsights_ErrorStatus(t *testing.T) {
	projectID := "33333333-3333-3333-3333-333333333333"
	jobID := "44444444-4444-4444-4444-444444444444"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": jobID, "name": "report", "status": "error",
			"error": "something went wrong", "clusters": []interface{}{},
		})
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	_, err := client.PollInsights(context.Background(), langsmith.PollInsightsParams{
		ProjectID: projectID,
		ID:        jobID,
		Rate:      5 * time.Millisecond,
	})
	if err == nil {
		t.Fatal("expected error when status is 'error', got nil")
	}
	if !strings.Contains(err.Error(), "something went wrong") {
		t.Errorf("expected error to contain job error message, got: %v", err)
	}
}

func TestPollInsights_Timeout(t *testing.T) {
	projectID := "55555555-5555-5555-5555-555555555555"
	jobID := "66666666-6666-6666-6666-666666666666"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(map[string]interface{}{
			"id": jobID, "name": "report", "status": "pending",
			"clusters": []interface{}{}, "error": nil,
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

func TestPollInsights_MutuallyExclusive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	_, err := client.PollInsights(context.Background(), langsmith.PollInsightsParams{
		Report:    &langsmith.InsightsReport{ID: "x", ProjectID: "y"},
		ID:        "x",
		ProjectID: "y",
	})
	if err == nil {
		t.Fatal("expected error when both Report and IDs are provided")
	}
}

func TestGetInsightsReport(t *testing.T) {
	projectID := "77777777-7777-7777-7777-777777777777"
	jobID := "88888888-8888-8888-8888-888888888888"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/"+projectID+"/insights/"+jobID:
			json.NewEncoder(w).Encode(map[string]interface{}{
				"id": jobID, "name": "report", "status": "success",
				"start_time": "2026-01-01T00:00:00Z",
				"end_time":   "2026-01-01T01:00:00Z",
				"shape":      map[string]int{"cluster-a": 3},
				"clusters": []map[string]interface{}{
					{"id": "cid", "name": "cluster-a", "description": "Cluster A",
						"level": 0, "num_runs": 3, "stats": nil},
				},
				"report": map[string]interface{}{
					"key_points":         []string{"topic X"},
					"title":              "Usage Summary",
					"highlighted_traces": []interface{}{},
				},
				"error": nil,
			})

		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/sessions/"+projectID+"/insights/"+jobID+"/runs":
			json.NewEncoder(w).Encode(map[string]interface{}{
				"runs":   []map[string]interface{}{{"id": "r1"}, {"id": "r2"}},
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
	if len(result.Clusters) != 1 || result.Clusters[0].Name != "cluster-a" {
		t.Errorf("unexpected clusters: %+v", result.Clusters)
	}
	if result.Report == nil || result.Report.Title != "Usage Summary" {
		t.Errorf("unexpected report: %+v", result.Report)
	}
	if len(result.Runs) != 2 {
		t.Errorf("want 2 runs, got %d", len(result.Runs))
	}
	if result.StartTime == nil {
		t.Error("expected non-nil StartTime")
	}
}

func TestListInsightsJobs(t *testing.T) {
	sessionID := "aaaabbbb-aaaa-aaaa-aaaa-aaaaaaaaaaaa"
	job1ID := "11112222-1111-1111-1111-111111111111"
	job2ID := "22223333-2222-2222-2222-222222222222"

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method != http.MethodGet || r.URL.Path != "/api/v1/sessions/"+sessionID+"/insights" {
			http.Error(w, "unexpected: "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
			return
		}
		json.NewEncoder(w).Encode([]map[string]interface{}{
			{"id": job1ID, "name": "job-one", "status": "success", "created_at": "2026-01-01T00:00:00Z", "clusters": []interface{}{}},
			{"id": job2ID, "name": "job-two", "status": "pending", "created_at": "2026-01-02T00:00:00Z", "clusters": []interface{}{}},
		})
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	jobs, err := client.ListInsightsJobs(context.Background(), sessionID)
	if err != nil {
		t.Fatalf("ListInsightsJobs error: %v", err)
	}
	if len(jobs) != 2 {
		t.Fatalf("want 2 jobs, got %d", len(jobs))
	}
	if jobs[0].ID != job1ID || jobs[0].Name != "job-one" || jobs[0].Status != "success" {
		t.Errorf("unexpected first job: %+v", jobs[0])
	}
	if jobs[1].ID != job2ID || jobs[1].Status != "pending" {
		t.Errorf("unexpected second job: %+v", jobs[1])
	}
}

func TestGetInsightsReport_MutuallyExclusive(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer srv.Close()

	client := newInsightsTestClient(t, srv)

	_, err := client.GetInsightsReport(context.Background(), langsmith.GetInsightsReportParams{
		Report:    &langsmith.InsightsReport{ID: "x", ProjectID: "y"},
		ID:        "x",
		ProjectID: "y",
	})
	if err == nil {
		t.Fatal("expected error when both Report and IDs are provided")
	}
}
