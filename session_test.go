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

func TestSessionNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Sessions.New(context.TODO(), langsmith.SessionNewParams{
		Upsert:           langsmith.F(true),
		ID:               langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		DefaultDatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Description:      langsmith.F("description"),
		EndTime:          langsmith.F(time.Now()),
		Extra: langsmith.F(map[string]interface{}{
			"foo": "bar",
		}),
		Name:               langsmith.F("name"),
		ReferenceDatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		StartTime:          langsmith.F(time.Now()),
		TraceTier:          langsmith.F(langsmith.SessionNewParamsTraceTierLonglived),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSessionGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Sessions.Get(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.SessionGetParams{
			IncludeStats:   langsmith.F(true),
			StatsStartTime: langsmith.F(time.Now()),
			Accept:         langsmith.F("accept"),
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

func TestSessionUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Sessions.Update(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.SessionUpdateParams{
			DefaultDatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Description:      langsmith.F("description"),
			EndTime:          langsmith.F(time.Now()),
			Extra: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Name:      langsmith.F("name"),
			TraceTier: langsmith.F(langsmith.SessionUpdateParamsTraceTierLonglived),
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

func TestSessionListWithOptionalParams(t *testing.T) {
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
	_, err := client.Sessions.List(context.TODO(), langsmith.SessionListParams{
		ID:                langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		DatasetVersion:    langsmith.F("dataset_version"),
		Facets:            langsmith.F(true),
		Filter:            langsmith.F("filter"),
		IncludeStats:      langsmith.F(true),
		Limit:             langsmith.F(int64(1)),
		Metadata:          langsmith.F("metadata"),
		Name:              langsmith.F("name"),
		NameContains:      langsmith.F("name_contains"),
		Offset:            langsmith.F(int64(0)),
		ReferenceDataset:  langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		ReferenceFree:     langsmith.F(true),
		SortBy:            langsmith.F(langsmith.SessionSortableColumnsName),
		SortByDesc:        langsmith.F(true),
		SortByFeedbackKey: langsmith.F("sort_by_feedback_key"),
		StatsStartTime:    langsmith.F(time.Now()),
		TagValueID:        langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		UseApproxStats:    langsmith.F(true),
		Accept:            langsmith.F("accept"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSessionDelete(t *testing.T) {
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
	_, err := client.Sessions.Delete(context.TODO(), "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSessionDashboardWithOptionalParams(t *testing.T) {
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
	_, err := client.Sessions.Dashboard(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.SessionDashboardParams{
			CustomChartsSectionRequest: langsmith.CustomChartsSectionRequestParam{
				EndTime: langsmith.F(time.Now()),
				GroupBy: langsmith.F(langsmith.RunStatsGroupByParam{
					Attribute: langsmith.F(langsmith.RunStatsGroupByAttributeName),
					MaxGroups: langsmith.F(int64(0)),
					Path:      langsmith.F("path"),
				}),
				OmitData:  langsmith.F(true),
				StartTime: langsmith.F(time.Now()),
				Stride: langsmith.F(langsmith.TimedeltaInputParam{
					Days:    langsmith.F(int64(0)),
					Hours:   langsmith.F(int64(0)),
					Minutes: langsmith.F(int64(0)),
				}),
				Timezone: langsmith.F("timezone"),
			},
			Accept: langsmith.F("accept"),
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
