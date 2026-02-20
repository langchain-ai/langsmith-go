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
		option.WithOrganizationID("My Organization ID"),
	)
	page, err := client.Datasets.List(context.TODO(), langsmith.DatasetListParams{
		Limit: langsmith.F(int64(100)),
	})
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	for _, dataset := range page.Items {
		t.Logf("%+v\n", dataset.ID)
	}
	// The mock server isn't going to give us real pagination
	page, err = page.GetNextPage()
	if err != nil {
		t.Fatalf("err should be nil: %s", err.Error())
	}
	if page != nil {
		for _, dataset := range page.Items {
			t.Logf("%+v\n", dataset.ID)
		}
	}
}
