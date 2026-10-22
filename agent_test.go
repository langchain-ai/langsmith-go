// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestAgentNewSourceTypeDefault(t *testing.T) {
	tests := []struct {
		name       string
		sourceType param.Field[langsmith.AgentNewParamsSourceType]
		want       any
	}{
		{name: "omitted", want: "EXTERNAL"},
		{name: "external", sourceType: langsmith.F(langsmith.AgentNewParamsSourceTypeExternal), want: "EXTERNAL"},
		{name: "lsd", sourceType: langsmith.F(langsmith.AgentNewParamsSourceTypeLsd), want: "LSD"},
		{name: "mda", sourceType: langsmith.F(langsmith.AgentNewParamsSourceTypeMda), want: "MDA"},
		{name: "fleet", sourceType: langsmith.F(langsmith.AgentNewParamsSourceTypeFleet), want: "FLEET"},
		{name: "empty", sourceType: langsmith.F(langsmith.AgentNewParamsSourceType("")), want: ""},
		{name: "null", sourceType: langsmith.Null[langsmith.AgentNewParamsSourceType](), want: nil},
		{name: "raw", sourceType: langsmith.Raw[langsmith.AgentNewParamsSourceType]("future-type"), want: "future-type"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodPost || r.URL.Path != "/api/v1/agents" {
					t.Errorf("unexpected request: %s %s", r.Method, r.URL.Path)
				}
				var body map[string]any
				if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
					t.Errorf("decode request: %v", err)
				}
				want := map[string]any{"id": "example-agent", "name": "Example agent", "source_type": tt.want}
				if !reflect.DeepEqual(body, want) {
					t.Errorf("body = %#v, want %#v", body, want)
				}
				for key, value := range map[string]string{"x-api-key": "test-api-key", "x-tenant-id": "test-workspace", "x-test-header": "preserved"} {
					if r.Header.Get(key) != value {
						t.Errorf("header %s = %q, want %q", key, r.Header.Get(key), value)
					}
				}
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"example-agent"}`))
			}))
			defer server.Close()
			client := langsmith.NewClient(option.WithBaseURL(server.URL), option.WithAPIKey("test-api-key"), option.WithTenantID("test-workspace"))
			body := langsmith.AgentNewParams{ID: langsmith.F("example-agent"), Name: langsmith.F("Example agent"), SourceType: tt.sourceType}
			if _, err := client.Agents.New(context.Background(), body, option.WithHeader("x-test-header", "preserved")); err != nil {
				t.Fatal(err)
			}
			if !reflect.DeepEqual(body.SourceType, tt.sourceType) {
				t.Fatal("caller source type was mutated")
			}
		})
	}
}

func TestAgentNewWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Agents.New(context.TODO(), langsmith.AgentNewParams{
		ID:             langsmith.F("id"),
		SourceType:     langsmith.F(langsmith.AgentNewParamsSourceTypeLsd),
		AgentType:      langsmith.F(langsmith.AgentNewParamsAgentTypeCustom),
		Name:           langsmith.F("name"),
		ProductionOnly: langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAgentGetWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Agents.Get(
		context.TODO(),
		"agent_id",
		langsmith.AgentGetParams{
			Include: langsmith.F(langsmith.AgentGetParamsIncludeRuntimeBindings),
		},
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAgentUpdateWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Agents.Update(
		context.TODO(),
		"agent_id",
		langsmith.AgentUpdateParams{
			AgentType:  langsmith.F(langsmith.AgentUpdateParamsAgentTypeCustom),
			Name:       langsmith.F("name"),
			OwnerID:    langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			SourceType: langsmith.F(langsmith.AgentUpdateParamsSourceTypeLsd),
		},
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAgentListWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Agents.List(context.TODO(), langsmith.AgentListParams{
		ID:         langsmith.F([]string{"string"}),
		Cursor:     langsmith.F("cursor"),
		Name:       langsmith.F("name"),
		PageSize:   langsmith.F(int64(1)),
		SortBy:     langsmith.F(langsmith.AgentListParamsSortByLastRunStartTime),
		SortOrder:  langsmith.F(langsmith.AgentListParamsSortOrderAsc),
		SourceType: langsmith.F(langsmith.AgentListParamsSourceTypeLsd),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestAgentDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	err := client.Agents.Delete(context.TODO(), "agent_id")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
