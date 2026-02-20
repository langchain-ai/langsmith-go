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

func TestDatasetRunNewWithOptionalParams(t *testing.T) {
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Datasets.Runs.New(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetRunNewParams{
			SessionIDs:              langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
			Format:                  langsmith.F(langsmith.DatasetRunNewParamsFormatCsv),
			ComparativeExperimentID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			ExampleIDs:              langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
			Filters: langsmith.F(map[string][]string{
				"foo": {"string"},
			}),
			Limit:   langsmith.F(int64(1)),
			Offset:  langsmith.F(int64(0)),
			Preview: langsmith.F(true),
			SortParams: langsmith.F(langsmith.SortParamsForRunsComparisonView{
				SortBy:    langsmith.F("sort_by"),
				SortOrder: langsmith.F(langsmith.SortParamsForRunsComparisonViewSortOrderAsc),
			}),
			Stream: langsmith.F(true),
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

func TestDatasetRunDeltaWithOptionalParams(t *testing.T) {
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Datasets.Runs.Delta(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetRunDeltaParams{
			QueryFeedbackDelta: langsmith.QueryFeedbackDeltaParam{
				BaselineSessionID:       langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
				ComparisonSessionIDs:    langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
				FeedbackKey:             langsmith.F("feedback_key"),
				ComparativeExperimentID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
				Filters: langsmith.F(map[string][]string{
					"foo": {"string"},
				}),
				Limit:  langsmith.F(int64(1)),
				Offset: langsmith.F(int64(0)),
			},
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
