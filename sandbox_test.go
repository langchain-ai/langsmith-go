// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"errors"
	"os"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxListUsageCostsWithOptionalParams(t *testing.T) {
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
	_, err := client.Sandboxes.ListUsageCosts(context.TODO(), langsmith.SandboxListUsageCostsParams{
		EndTime:      langsmith.F(time.Now()),
		StartTime:    langsmith.F(time.Now()),
		Cursor:       langsmith.F("cursor"),
		Granularity:  langsmith.F(langsmith.SandboxListUsageCostsParamsGranularityHour),
		PageSize:     langsmith.F(int64(1)),
		ResourceIDs:  langsmith.F([]string{"string"}),
		ResourceType: langsmith.F(langsmith.SandboxListUsageCostsParamsResourceTypeSandbox),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
