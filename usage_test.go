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

func TestUsage(t *testing.T) {
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
	customChartsSection, err := client.Sessions.Dashboard(
		context.TODO(),
		"1ffaeba7-541e-469f-bae7-df3208ea3d45",
		langsmith.SessionDashboardParams{
			CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{},
		},
	)
	if err != nil {
		t.Error(err)
		return
	}
	t.Logf("%+v\n", customChartsSection.ID)
}
