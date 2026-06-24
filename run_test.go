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
	)
	_, err := client.Runs.Update(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
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

func TestRunQueryV1WithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.QueryV1(context.TODO(), langsmith.RunQueryV1Params{
		ID:                    langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Cursor:                langsmith.F("cursor"),
		DataSourceType:        langsmith.F(langsmith.RunsFilterDataSourceTypeEnumCurrent),
		EndTime:               langsmith.F(time.Now()),
		Error:                 langsmith.F(true),
		ExecutionOrder:        langsmith.F(int64(1)),
		Filter:                langsmith.F("filter"),
		IsRoot:                langsmith.F(true),
		Limit:                 langsmith.F(int64(1)),
		Order:                 langsmith.F(langsmith.RunQueryV1ParamsOrderAsc),
		ParentRun:             langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Query:                 langsmith.F("query"),
		ReferenceExample:      langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		RunType:               langsmith.F(langsmith.RunTypeEnumTool),
		SearchFilter:          langsmith.F("search_filter"),
		Select:                langsmith.F([]langsmith.RunQueryV1ParamsSelect{langsmith.RunQueryV1ParamsSelectID}),
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

func TestRunQueryV2WithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.QueryV2(context.TODO(), langsmith.RunQueryV2Params{
		Cursor:             langsmith.F("eyJ2IjoxLCJhIjoicnVucy5xdWVyeSIsImsiOiJwYXNzIiwiYiI6InNkYiIsInQiOiJsdChjdXJzb3IsICcyMDI1LTEyLTEyIDE5OjAzOjI4LjQ4MTI1NTAxOWIxM2YyJykifQ"),
		Filter:             langsmith.F(`and(eq(run_type, "llm"), gt(latency, 5))`),
		HasError:           langsmith.F(false),
		IDs:                langsmith.F([]string{"018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9327", "f47ac10b-58cc-4372-a567-0e02b2c3d479"}),
		IsRoot:             langsmith.F(true),
		MaxStartTime:       langsmith.F(time.Now()),
		MinStartTime:       langsmith.F(time.Now()),
		PageSize:           langsmith.F(int64(100)),
		ProjectIDs:         langsmith.F([]string{"018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9327", "0190a1b2-c3d4-7ef0-a5b6-6ea3a82e9328"}),
		ReferenceDatasetID: langsmith.F("018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9327"),
		ReferenceExamples:  langsmith.F([]string{"b2c3d4e5-f6a7-4b5c-9d0e-1f2a3b4c5d6e", "c3d4e5f6-a7b8-4c5d-0e1f-2a3b4c5d6e7f"}),
		RunType:            langsmith.F(langsmith.RunQueryV2ParamsRunTypeLlm),
		Selects:            langsmith.F([]langsmith.RunQueryV2ParamsSelect{langsmith.RunQueryV2ParamsSelectID, langsmith.RunQueryV2ParamsSelectName, langsmith.RunQueryV2ParamsSelectProjectID, langsmith.RunQueryV2ParamsSelectStartTime, langsmith.RunQueryV2ParamsSelectRunType, langsmith.RunQueryV2ParamsSelectStatus}),
		TraceFilter:        langsmith.F(`eq(status, "success")`),
		TraceID:            langsmith.F("f47ac10b-58cc-4372-a567-0e02b2c3d479"),
		TreeFilter:         langsmith.F(`has(tags, "production")`),
		Accept:             langsmith.F("Accept"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRunGetV1WithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.GetV1(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.RunGetV1Params{
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

func TestRunGetV2WithOptionalParams(t *testing.T) {
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
	_, err := client.Runs.GetV2(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.RunGetV2Params{
			ProjectID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			StartTime: langsmith.F(time.Now()),
			Selects:   langsmith.F([]langsmith.RunGetV2ParamsSelect{langsmith.RunGetV2ParamsSelectID}),
			Accept:    langsmith.F("Accept"),
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
			IncludeDetails:        langsmith.F(true),
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
