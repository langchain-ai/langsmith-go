package runsv2_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go/lib/runsv2"
)

func TestQuery_BuildsRequestAndDecodesResponse(t *testing.T) {
	var (
		gotMethod   string
		gotPath     string
		gotAPIKey   string
		gotBody     map[string]any
		nextCursor  = "next-cursor-abc"
		hasMore     = true
		respPayload = map[string]any{
			"items": []map[string]any{
				{
					"id":        "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9327",
					"name":      "ChatOpenAI",
					"run_type":  "llm",
					"status":    "SUCCESS",
					"is_root":   true,
					"trace_id":  "018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9327",
					"project_id": "0190a1b2-c3d4-7ef0-a5b6-6ea3a82e9328",
				},
			},
			"next_cursor": nextCursor,
			"has_more":    hasMore,
		}
	)

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotMethod = r.Method
		gotPath = r.URL.Path
		gotAPIKey = r.Header.Get("X-API-Key")
		body, _ := io.ReadAll(r.Body)
		_ = json.Unmarshal(body, &gotBody)

		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(respPayload)
	}))
	defer srv.Close()

	c := runsv2.NewClient(srv.URL, "secret-key")
	projectID := "0190a1b2-c3d4-7ef0-a5b6-6ea3a82e9328"
	minStart := "2024-01-01T00:00:00Z"

	resp, err := c.Query(context.Background(), runsv2.QueryRequest{
		ProjectIDs:   []string{projectID},
		MinStartTime: &minStart,
		IsRoot:       runsv2.Ptr(true),
		PageSize:     runsv2.Ptr(uint64(50)),
		Selects:      []runsv2.SelectField{runsv2.SelectID, runsv2.SelectName, runsv2.SelectRunType},
	})
	if err != nil {
		t.Fatalf("Query returned error: %v", err)
	}

	if gotMethod != http.MethodPost {
		t.Errorf("method = %q, want POST", gotMethod)
	}
	if gotPath != "/v2/runs/query" {
		t.Errorf("path = %q, want /v2/runs/query", gotPath)
	}
	if gotAPIKey != "secret-key" {
		t.Errorf("X-API-Key = %q, want secret-key", gotAPIKey)
	}
	if ids, _ := gotBody["project_ids"].([]any); len(ids) != 1 || ids[0] != projectID {
		t.Errorf("project_ids = %v, want [%q]", gotBody["project_ids"], projectID)
	}
	if gotBody["min_start_time"] != minStart {
		t.Errorf("min_start_time = %v, want %q", gotBody["min_start_time"], minStart)
	}
	if gotBody["is_root"] != true {
		t.Errorf("is_root = %v, want true", gotBody["is_root"])
	}
	if _, ok := gotBody["sort_order"]; ok {
		t.Errorf("sort_order should be omitted when nil")
	}

	if len(resp.Items) != 1 {
		t.Fatalf("len(items) = %d, want 1", len(resp.Items))
	}
	if resp.Items[0].Name == nil || *resp.Items[0].Name != "ChatOpenAI" {
		t.Errorf("items[0].Name = %v, want ChatOpenAI", resp.Items[0].Name)
	}
	if resp.NextCursor == nil || *resp.NextCursor != nextCursor {
		t.Errorf("NextCursor = %v, want %q", resp.NextCursor, nextCursor)
	}
	if !resp.HasMore {
		t.Errorf("HasMore = false, want true")
	}
}

func TestQuery_StripsTrailingAPIV1FromBaseURL(t *testing.T) {
	var gotPath string
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotPath = r.URL.Path
		_, _ = w.Write([]byte(`{"items":[],"has_more":false}`))
	}))
	defer srv.Close()

	c := runsv2.NewClient(srv.URL+"/api/v1", "k")
	if _, err := c.Query(context.Background(), runsv2.QueryRequest{}); err != nil {
		t.Fatalf("Query: %v", err)
	}
	if gotPath != "/v2/runs/query" {
		t.Errorf("path = %q, want /v2/runs/query (api/v1 suffix should be stripped)", gotPath)
	}
}

func TestQuery_ReturnsHTTPErrorOnNon2xx(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnprocessableEntity)
		_, _ = w.Write([]byte(`{"detail":"missing project_ids"}`))
	}))
	defer srv.Close()

	c := runsv2.NewClient(srv.URL, "k")
	_, err := c.Query(context.Background(), runsv2.QueryRequest{})
	if err == nil {
		t.Fatal("expected error, got nil")
	}
	httpErr, ok := err.(*runsv2.HTTPError)
	if !ok {
		t.Fatalf("error type = %T, want *HTTPError", err)
	}
	if httpErr.StatusCode != http.StatusUnprocessableEntity {
		t.Errorf("StatusCode = %d, want 422", httpErr.StatusCode)
	}
	if !strings.Contains(httpErr.Body, "missing project_ids") {
		t.Errorf("Body = %q, want it to contain server detail", httpErr.Body)
	}
}
