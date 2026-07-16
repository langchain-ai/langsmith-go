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

func TestDatasetExperimentRunQueryWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.ExperimentRuns.Query(
		context.TODO(),
		"dataset_id",
		langsmith.DatasetExperimentRunQueryParams{
			ComparativeExperimentID: langsmith.F("comparative_experiment_id"),
			Cursor:                  langsmith.F("cursor"),
			ExampleIDs:              langsmith.F([]string{"string"}),
			ExperimentIDs:           langsmith.F([]string{"string"}),
			Filters: langsmith.F(map[string][]string{
				"foo": {"string"},
			}),
			PageSize: langsmith.F(int64(0)),
			Selects:  langsmith.F([]langsmith.RunSelectField{langsmith.RunSelectFieldID}),
			Sort: langsmith.F(langsmith.DatasetExperimentRunQueryParamsSort{
				By:    langsmith.F("by"),
				Order: langsmith.F("order"),
			}),
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
