// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestIssueGet(t *testing.T) {
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
	_, err := client.Issues.Get(context.TODO(), "id")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestIssueListWithOptionalParams(t *testing.T) {
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
	_, err := client.Issues.List(context.TODO(), langsmith.IssueListParams{
		Activity:      langsmith.F([]langsmith.IssueListParamsActivity{langsmith.IssueListParamsActivityFixing}),
		Limit:         langsmith.F(int64(0)),
		Offset:        langsmith.F(int64(0)),
		SessionID:     langsmith.F("session_id"),
		SessionName:   langsmith.F("session_name"),
		Severity:      langsmith.F(langsmith.IssueListParamsSeverity0),
		SeverityExact: langsmith.F([]langsmith.IssueListParamsSeverityExact{langsmith.IssueListParamsSeverityExact0}),
		SortBy:        langsmith.F(langsmith.IssueListParamsSortByDefault),
		Status:        langsmith.F(langsmith.IssueListParamsStatusOpen),
		StatusFirst:   langsmith.F(true),
		Tag:           langsmith.F("tag"),
		TraceID:       langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		UpdatedAt:     langsmith.F("updated_at"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
