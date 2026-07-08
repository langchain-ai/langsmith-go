// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"os"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestManualPagination(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	} else {
		t.Skip("requires mock Prism server; set TEST_API_BASE_URL to run")
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	page, err := client.Runs.QueryV2(context.TODO(), langsmith.RunQueryV2Params{
		ProjectIDs: langsmith.F([]string{"00000000-0000-0000-0000-000000000000"}),
	})
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	for _, run := range page.Items {
		t.Logf("%+v\n", run.ID)
	}
	// The mock server isn't going to give us real pagination
	page, err = page.GetNextPage()
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if page != nil {
		for _, run := range page.Items {
			t.Logf("%+v\n", run.ID)
		}
	}
}
