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

func TestRunIngestBatchWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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

func TestRunQueryWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
		DataSourceType:        langsmith.F(langsmith.RunQueryParamsDataSourceTypeCurrent),
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
		RunType:               langsmith.F(langsmith.RunQueryParamsRunTypeTool),
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
