package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSqlQueryContract(t *testing.T) {
	projectIDs := []string{"11111111-1111-4111-8111-111111111111", "22222222-2222-4222-8222-222222222222"}
	query := "SELECT 1 AS value"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/v2/sql/query" {
			t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		var body struct {
			ProjectIDs []string `json:"project_ids"`
			Query      string   `json:"query"`
		}
		if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
			t.Errorf("decode request: %v", err)
		}
		if !reflect.DeepEqual(body.ProjectIDs, projectIDs) || body.Query != query {
			t.Errorf("unexpected body: %#v", body)
		}
		if r.Header.Get("X-Api-Key") != "test-api-key" || r.Header.Get("X-Tenant-Id") != "test-workspace" {
			t.Error("missing authentication or workspace header")
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"columns":[{"name":"value","data_type":"integer"}],"rows":[[1,"text",true,null,[1],{"key":"value"}]],"row_count":1,"is_truncated":false,"truncation_reason":null,"execution_time_ms":1.5}`))
	}))
	defer server.Close()
	client := langsmith.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-api-key"), option.WithTenantID("test-workspace"))
	result, err := client.Sql.Query(context.Background(), langsmith.SqlQueryParams{ProjectIDs: langsmith.F(projectIDs), Query: langsmith.F(query)})
	if err != nil {
		t.Fatal(err)
	}
	wantRows := [][]any{{float64(1), "text", true, nil, []any{float64(1)}, map[string]any{"key": "value"}}}
	if !reflect.DeepEqual(result.Rows, wantRows) || result.RowCount != 1 || result.IsTruncated || !result.JSON.TruncationReason.IsNull() || result.ExecutionTimeMs != 1.5 {
		t.Errorf("unexpected result: %#v", result)
	}
}
