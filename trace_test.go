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

func TestTraceListRunsWithOptionalParams(t *testing.T) {
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
	_, err := client.Traces.ListRuns(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.TraceListRunsParams{
			ProjectID:    langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Filter:       langsmith.F("filter"),
			MaxStartTime: langsmith.F(time.Now()),
			MinStartTime: langsmith.F(time.Now()),
			Selects:      langsmith.F([]langsmith.TraceListRunsParamsSelect{langsmith.TraceListRunsParamsSelectID}),
			Accept:       langsmith.F("Accept"),
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

func TestTraceQueryWithOptionalParams(t *testing.T) {
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
	_, err := client.Traces.Query(context.TODO(), langsmith.TraceQueryParams{
		Cursor:       langsmith.F("cursor"),
		MaxStartTime: langsmith.F(time.Now()),
		MinStartTime: langsmith.F(time.Now()),
		PageSize:     langsmith.F(int64(20)),
		ProjectID:    langsmith.F("018e4c7e-a9fb-7ef0-a5b6-6ea3a82e9327"),
		Selects:      langsmith.F([]langsmith.TraceQueryParamsSelect{langsmith.TraceQueryParamsSelectID, langsmith.TraceQueryParamsSelectName, langsmith.TraceQueryParamsSelectStartTime, langsmith.TraceQueryParamsSelectStatus, langsmith.TraceQueryParamsSelectTotalTokens, langsmith.TraceQueryParamsSelectTotalCost, langsmith.TraceQueryParamsSelectFirstTokenTime}),
		TraceFilter:  langsmith.F(`eq(status, "error")`),
		TraceIDs:     langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		TreeFilter:   langsmith.F(`has(tags, "production")`),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
