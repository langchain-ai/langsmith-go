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

func TestThreadListTracesWithOptionalParams(t *testing.T) {
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
	_, err := client.Threads.ListTraces(
		context.TODO(),
		"thread_id",
		langsmith.ThreadListTracesParams{
			ProjectID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Cursor:    langsmith.F("cursor"),
			Filter:    langsmith.F("filter"),
			PageSize:  langsmith.F(int64(1)),
			Selects:   langsmith.F([]langsmith.ThreadListTracesParamsSelect{langsmith.ThreadListTracesParamsSelectThreadID}),
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

func TestThreadQueryWithOptionalParams(t *testing.T) {
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
	_, err := client.Threads.Query(context.TODO(), langsmith.ThreadQueryParams{
		Cursor:       langsmith.F("cursor"),
		Filter:       langsmith.F("filter"),
		MaxStartTime: langsmith.F(time.Now()),
		MinStartTime: langsmith.F(time.Now()),
		PageSize:     langsmith.F(int64(20)),
		ProjectID:    langsmith.F("0190a1b2-c3d4-7ef0-a5b6-6ea3a82e9328"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestThreadStatsWithOptionalParams(t *testing.T) {
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
	_, err := client.Threads.Stats(
		context.TODO(),
		"thread_id",
		langsmith.ThreadStatsParams{
			Selects:   langsmith.F([]langsmith.ThreadStatsParamsSelect{langsmith.ThreadStatsParamsSelectTurns}),
			SessionID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Filter:    langsmith.F("filter"),
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
