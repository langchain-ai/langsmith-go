// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestRunNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.New(context.TODO(), langsmith.RunNewParams{
		Run: langsmith.RunParam{
			ID:          langsmith.F("id"),
			DottedOrder: langsmith.F("dotted_order"),
			EndTime:     langsmith.F("end_time"),
			Error:       langsmith.F("error"),
			Events: langsmith.F([]map[string]interface{}{{
				"foo": "bar",
			}}),
			Extra: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			InputAttachments: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Inputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Name: langsmith.F("name"),
			OutputAttachments: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Outputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			ParentRunID:        langsmith.F("parent_run_id"),
			ReferenceExampleID: langsmith.F("reference_example_id"),
			RunType:            langsmith.F(langsmith.RunRunTypeTool),
			Serialized: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			SessionID:   langsmith.F("session_id"),
			SessionName: langsmith.F("session_name"),
			StartTime:   langsmith.F("start_time"),
			Status:      langsmith.F("status"),
			Tags:        langsmith.F([]string{"string"}),
			TraceID:     langsmith.F("trace_id"),
		},
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRunGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.Get(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.RunGetParams{
			ExcludeS3StoredAttributes: langsmith.F(true),
			ExcludeSerialized:         langsmith.F(true),
			IncludeMessages:           langsmith.F(true),
			SessionID:                 langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			StartTime:                 langsmith.F(time.Now()),
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

func TestRunUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.Update(
		context.TODO(),
		"run_id",
		langsmith.RunUpdateParams{
			Run: langsmith.RunParam{
				ID:          langsmith.F("id"),
				DottedOrder: langsmith.F("dotted_order"),
				EndTime:     langsmith.F("end_time"),
				Error:       langsmith.F("error"),
				Events: langsmith.F([]map[string]interface{}{{
					"foo": "bar",
				}}),
				Extra: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				InputAttachments: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				Inputs: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				Name: langsmith.F("name"),
				OutputAttachments: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				Outputs: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				ParentRunID:        langsmith.F("parent_run_id"),
				ReferenceExampleID: langsmith.F("reference_example_id"),
				RunType:            langsmith.F(langsmith.RunRunTypeTool),
				Serialized: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				SessionID:   langsmith.F("session_id"),
				SessionName: langsmith.F("session_name"),
				StartTime:   langsmith.F("start_time"),
				Status:      langsmith.F("status"),
				Tags:        langsmith.F([]string{"string"}),
				TraceID:     langsmith.F("trace_id"),
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

func TestRunIngestBatchWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.IngestBatch(context.TODO(), langsmith.RunIngestBatchParams{
		Patch: langsmith.F([]langsmith.RunParam{{
			ID:          langsmith.F("id"),
			DottedOrder: langsmith.F("dotted_order"),
			EndTime:     langsmith.F("end_time"),
			Error:       langsmith.F("error"),
			Events: langsmith.F([]map[string]interface{}{{
				"foo": "bar",
			}}),
			Extra: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			InputAttachments: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Inputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Name: langsmith.F("name"),
			OutputAttachments: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Outputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			ParentRunID:        langsmith.F("parent_run_id"),
			ReferenceExampleID: langsmith.F("reference_example_id"),
			RunType:            langsmith.F(langsmith.RunRunTypeTool),
			Serialized: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			SessionID:   langsmith.F("session_id"),
			SessionName: langsmith.F("session_name"),
			StartTime:   langsmith.F("start_time"),
			Status:      langsmith.F("status"),
			Tags:        langsmith.F([]string{"string"}),
			TraceID:     langsmith.F("trace_id"),
		}}),
		Post: langsmith.F([]langsmith.RunParam{{
			ID:          langsmith.F("id"),
			DottedOrder: langsmith.F("dotted_order"),
			EndTime:     langsmith.F("end_time"),
			Error:       langsmith.F("error"),
			Events: langsmith.F([]map[string]interface{}{{
				"foo": "bar",
			}}),
			Extra: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			InputAttachments: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Inputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Name: langsmith.F("name"),
			OutputAttachments: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Outputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			ParentRunID:        langsmith.F("parent_run_id"),
			ReferenceExampleID: langsmith.F("reference_example_id"),
			RunType:            langsmith.F(langsmith.RunRunTypeTool),
			Serialized: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			SessionID:   langsmith.F("session_id"),
			SessionName: langsmith.F("session_name"),
			StartTime:   langsmith.F("start_time"),
			Status:      langsmith.F("status"),
			Tags:        langsmith.F([]string{"string"}),
			TraceID:     langsmith.F("trace_id"),
		}}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRunIngestMultipartWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.IngestMultipart(context.TODO(), langsmith.RunIngestMultipartParams{
		AttachmentRunIDFilename: langsmith.F(io.Reader(bytes.NewBuffer([]byte("Example data")))),
		FeedbackRunID:           langsmith.F(io.Reader(bytes.NewBuffer([]byte("Example data")))),
		PatchRunID:              langsmith.F(io.Reader(bytes.NewBuffer([]byte("Example data")))),
		PatchRunIDOutputs:       langsmith.F(io.Reader(bytes.NewBuffer([]byte("Example data")))),
		PostRunID:               langsmith.F(io.Reader(bytes.NewBuffer([]byte("Example data")))),
		PostRunIDInputs:         langsmith.F(io.Reader(bytes.NewBuffer([]byte("Example data")))),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRunQueryWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.Query(context.TODO(), langsmith.RunQueryParams{
		ID:                    langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Cursor:                langsmith.F("cursor"),
		DataSourceType:        langsmith.F(langsmith.RunsFilterDataSourceTypeEnumCurrent),
		EndTime:               langsmith.F(time.Now()),
		Error:                 langsmith.F(true),
		ExecutionOrder:        langsmith.F(int64(1)),
		Filter:                langsmith.F("filter"),
		IsRoot:                langsmith.F(true),
		Limit:                 langsmith.F(int64(1)),
		Order:                 langsmith.F(langsmith.RunQueryParamsOrderAsc),
		ParentRun:             langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Query:                 langsmith.F("query"),
		ReferenceExample:      langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		RunType:               langsmith.F(langsmith.RunTypeEnumTool),
		SearchFilter:          langsmith.F("search_filter"),
		Select:                langsmith.F([]langsmith.RunQueryParamsSelect{langsmith.RunQueryParamsSelectID}),
		Session:               langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		SkipPagination:        langsmith.F(true),
		SkipPrevCursor:        langsmith.F(true),
		StartTime:             langsmith.F(time.Now()),
		Trace:                 langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		TraceFilter:           langsmith.F("trace_filter"),
		TreeFilter:            langsmith.F("tree_filter"),
		UseExperimentalSearch: langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRunStatsWithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.Stats(context.TODO(), langsmith.RunStatsParams{
		RunStatsQueryParams: langsmith.RunStatsQueryParams{
			ID:             langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
			DataSourceType: langsmith.F(langsmith.RunsFilterDataSourceTypeEnumCurrent),
			EndTime:        langsmith.F(time.Now()),
			Error:          langsmith.F(true),
			ExecutionOrder: langsmith.F(int64(1)),
			Filter:         langsmith.F("filter"),
			GroupBy: langsmith.F(langsmith.RunStatsGroupByParam{
				Attribute: langsmith.F(langsmith.RunStatsGroupByAttributeName),
				MaxGroups: langsmith.F(int64(0)),
				Path:      langsmith.F("path"),
			}),
			Groups:                langsmith.F([]string{"string"}),
			IsRoot:                langsmith.F(true),
			ParentRun:             langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Query:                 langsmith.F("query"),
			ReferenceExample:      langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
			RunType:               langsmith.F(langsmith.RunTypeEnumTool),
			SearchFilter:          langsmith.F("search_filter"),
			Select:                langsmith.F([]langsmith.RunStatsQueryParamsSelect{langsmith.RunStatsQueryParamsSelectRunCount}),
			Session:               langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
			SkipPagination:        langsmith.F(true),
			StartTime:             langsmith.F(time.Now()),
			Trace:                 langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			TraceFilter:           langsmith.F("trace_filter"),
			TreeFilter:            langsmith.F("tree_filter"),
			UseExperimentalSearch: langsmith.F(true),
		},
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRunUpdate2(t *testing.T) {
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
	_, err := client.Runs.Update2(context.TODO(), "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
