package langsmith_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"net/url"
	"slices"
	"sync/atomic"
	"testing"
	"time"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

var ruleListIDs = []string{
	"00000000-0000-4000-8000-000000000001",
	"00000000-0000-4000-8000-000000000002",
	"00000000-0000-4000-8000-000000000003",
}

var ruleListTags = []string{
	"00000000-0000-4000-8000-000000000011",
	"00000000-0000-4000-8000-000000000012",
}

const ruleListProject = "00000000-0000-4000-8000-000000000021"

func ruleListFixture(id string) map[string]any {
	return map[string]any{
		"id": id, "tenant_id": "00000000-0000-4000-8000-000000000031",
		"created_at": "2026-10-13T00:00:00Z", "updated_at": "2026-10-13T00:00:00Z",
		"display_name": "verification", "evaluator_version": 3, "sampling_rate": 1,
		"add_to_dataset_prefer_correction": false, "extend_only": false,
		"include_extended_stats": false, "is_enabled": true, "is_managed_evaluator": false,
		"is_tracing_disabled": false, "transient": false, "use_corrections_dataset": false,
		"webhooks":        nil,
		"evaluator_type":  "CODE",
		"code_evaluators": []map[string]any{{"code": "return 1", "require_attachments": false, "code_evaluator_input": "TRAJECTORY", "advanced_features_enabled": true}},
	}
}

func TestEvaluatorListV2AutoPagingCompatibility(t *testing.T) {
	t.Parallel()
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	wantIDs := slices.Clone(ruleListIDs)
	cursors := []string{"", "opaque+/cursor?one", "opaque+/cursor?two"}
	var requests atomic.Int32
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		request := int(requests.Add(1))
		if request > len(wantIDs) {
			http.Error(w, "request limit exceeded", http.StatusBadRequest)
			return
		}
		wantQuery := url.Values{
			"id": wantIDs, "tag_value_id": ruleListTags,
			"project_id": {ruleListProject}, "page_size": {"1"}, "type": {"PROJECT"},
			"name_contains": {"verification"}, "include_backfill_progress": {"true"},
		}
		if cursor := cursors[request-1]; cursor != "" {
			wantQuery.Set("cursor", cursor)
		}
		if r.Method != http.MethodGet || r.URL.Path != "/api/v2/runs/rules" {
			t.Errorf("request %d: got %s %s", request, r.Method, r.URL.Path)
		}
		if got := r.URL.Query().Encode(); got != wantQuery.Encode() {
			t.Errorf("request %d query: got %q, want %q", request, got, wantQuery.Encode())
		}
		var next any
		if request < len(wantIDs) {
			next = cursors[request]
		}
		w.Header().Set("Content-Type", "application/json")
		if err := json.NewEncoder(w).Encode(map[string]any{
			"items": []map[string]any{ruleListFixture(wantIDs[request-1])}, "next_cursor": next,
		}); err != nil {
			t.Errorf("encode page %d: %v", request, err)
		}
	}))
	defer server.Close()
	service := langsmith.NewEvaluatorService(
		option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic-test-key"), option.WithTenantID("synthetic-test-tenant"), option.WithMaxRetries(0),
	)
	pager := service.ListV2AutoPaging(ctx, langsmith.EvaluatorListV2Params{
		ID: langsmith.F(wantIDs), TagValueID: langsmith.F(ruleListTags),
		ProjectID: langsmith.F(ruleListProject), PageSize: langsmith.F(int64(1)),
		Type:         langsmith.F(langsmith.EvaluatorListV2ParamsTypeProject),
		NameContains: langsmith.F("verification"), IncludeBackfillProgress: langsmith.F(true),
	})
	var gotIDs []string
	for pager.Next() {
		if len(gotIDs) >= len(wantIDs) {
			t.Fatal("item limit exceeded")
		}
		gotIDs = append(gotIDs, pager.Current().ID)
		item := pager.Current()
		if item.EvaluatorType != "CODE" || len(item.CodeEvaluators) != 1 || item.CodeEvaluators[0].CodeEvaluatorInput != langsmith.RuleCodeEvaluatorsCodeEvaluatorInputTrajectory || !item.CodeEvaluators[0].AdvancedFeaturesEnabled {
			t.Fatalf("public code evaluator fields were not preserved: %+v", item)
		}
	}
	if err := pager.Err(); err != nil {
		t.Fatalf("auto-paging: %v", err)
	}
	if !slices.Equal(gotIDs, wantIDs) {
		t.Fatalf("yielded IDs: got %v, want %v", gotIDs, wantIDs)
	}
	if got := requests.Load(); got != int32(len(wantIDs)) {
		t.Fatalf("requests: got %d, want %d", got, len(wantIDs))
	}
	t.Logf("generated ListV2AutoPaging yielded %d items in exactly %d requests; final null cursor stopped iteration", len(gotIDs), requests.Load())
}

func TestEvaluatorListV2TerminalPages(t *testing.T) {
	t.Parallel()
	for _, scenario := range []string{"empty", "failure"} {
		t.Run(scenario, func(t *testing.T) {
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				count := requests.Add(1)
				w.Header().Set("Content-Type", "application/json")
				if count > 2 {
					t.Error("unexpected extra page request")
					w.WriteHeader(http.StatusBadRequest)
					return
				}
				if scenario == "empty" {
					_ = json.NewEncoder(w).Encode(map[string]any{"items": []any{}, "next_cursor": nil})
					return
				}
				if count == 1 {
					_ = json.NewEncoder(w).Encode(map[string]any{"items": []any{ruleListFixture(ruleListIDs[0])}, "next_cursor": "opaque-cursor"})
					return
				}
				if r.URL.Query().Get("cursor") != "opaque-cursor" {
					t.Error("continuation cursor was not preserved")
				}
				w.WriteHeader(http.StatusServiceUnavailable)
				_ = json.NewEncoder(w).Encode(map[string]any{"type": "about:blank", "title": "Service Unavailable", "status": 503, "detail": "Please retry later"})
			}))
			defer server.Close()
			client := langsmith.NewClient(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic-test-key"), option.WithTenantID("synthetic-test-tenant"), option.WithMaxRetries(0))
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			pager := client.Evaluators.ListV2AutoPaging(ctx, langsmith.EvaluatorListV2Params{})
			count := 0
			for pager.Next() {
				count++
			}
			if scenario == "empty" {
				if count != 0 || requests.Load() != 1 || pager.Err() != nil {
					t.Fatalf("empty page: count=%d requests=%d error=%v", count, requests.Load(), pager.Err())
				}
				return
			}
			if count != 1 || requests.Load() != 2 || pager.Err() == nil {
				t.Fatalf("failed continuation: count=%d requests=%d error=%v", count, requests.Load(), pager.Err())
			}
			var apiError *langsmith.Error
			if !errors.As(pager.Err(), &apiError) || apiError.StatusCode != http.StatusServiceUnavailable {
				t.Fatalf("continuation error did not preserve status 503: %v", pager.Err())
			}
		})
	}
}

func TestEvaluatorListV1Compatibility(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name string
		call func(context.Context, *langsmith.EvaluatorService) (*[]langsmith.Evaluator, error)
	}{
		{
			name: "legacy List signature and route",
			call: func(ctx context.Context, service *langsmith.EvaluatorService) (*[]langsmith.Evaluator, error) {
				var originalSignature func(context.Context, langsmith.EvaluatorListParams, ...option.RequestOption) (*[]langsmith.Evaluator, error) = service.List
				return originalSignature(ctx, langsmith.EvaluatorListParams{
					ID: langsmith.F(ruleListIDs), TagValueID: langsmith.F(ruleListTags), SessionID: langsmith.F(ruleListProject),
				})
			},
		},
		{
			name: "explicit ListV1 signature and route",
			call: func(ctx context.Context, service *langsmith.EvaluatorService) (*[]langsmith.Evaluator, error) {
				return service.ListV1(ctx, langsmith.EvaluatorListV1Params{
					ID: langsmith.F(ruleListIDs), TagValueID: langsmith.F(ruleListTags), SessionID: langsmith.F(ruleListProject),
				})
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancel()
			var requests atomic.Int32
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if requests.Add(1) > 1 {
					http.Error(w, "request limit exceeded", http.StatusBadRequest)
					return
				}
				wantQuery := url.Values{"id": ruleListIDs, "tag_value_id": ruleListTags, "session_id": {ruleListProject}}
				if r.Method != http.MethodGet || r.URL.Path != "/api/v1/runs/rules" || r.URL.Query().Encode() != wantQuery.Encode() {
					t.Errorf("V1 request: got %s %s?%s, want GET /api/v1/runs/rules?%s", r.Method, r.URL.Path, r.URL.RawQuery, wantQuery.Encode())
				}
				w.Header().Set("Content-Type", "application/json")
				legacy := ruleListFixture(ruleListIDs[0])
				// V1 returns enum values in snake_case.
				legacy["code_evaluators"].([]map[string]any)[0]["code_evaluator_input"] = "trajectory"
				legacy["evaluator_type"] = "code"
				if err := json.NewEncoder(w).Encode([]map[string]any{legacy}); err != nil {
					t.Errorf("encode V1 response: %v", err)
				}
			}))
			defer server.Close()
			service := langsmith.NewEvaluatorService(option.WithBaseURL(server.URL), option.WithHTTPClient(server.Client()), option.WithAPIKey("synthetic-test-key"), option.WithTenantID("synthetic-test-tenant"), option.WithMaxRetries(0))
			items, err := test.call(ctx, service)
			if err != nil {
				t.Fatalf("V1 list: %v", err)
			}
			if items == nil || len(*items) != 1 || (*items)[0].ID != ruleListIDs[0] || requests.Load() != 1 {
				t.Fatalf("V1 array response or request count mismatch: items=%v, requests=%d", items, requests.Load())
			}
			t.Log("generated V1 method retained its array result, repeated filters and /api/v1/runs/rules route")
		})
	}
}
