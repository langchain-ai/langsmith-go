// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stainless-sdks/langsmith-api-go"
	"github.com/stainless-sdks/langsmith-api-go/internal/testutil"
	"github.com/stainless-sdks/langsmith-api-go/option"
)

func TestDatasetPlaygroundExperimentBatchWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.PlaygroundExperiment.Batch(context.TODO(), langsmith.DatasetPlaygroundExperimentBatchParams{
		DatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Manifest:  langsmith.F[any](map[string]interface{}{}),
		Options: langsmith.F(langsmith.RunnableConfigParam{
			Callbacks:      langsmith.F([]interface{}{map[string]interface{}{}}),
			Configurable:   langsmith.F[any](map[string]interface{}{}),
			MaxConcurrency: langsmith.F(int64(0)),
			Metadata:       langsmith.F[any](map[string]interface{}{}),
			RecursionLimit: langsmith.F(int64(0)),
			RunID:          langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			RunName:        langsmith.F("run_name"),
			Tags:           langsmith.F([]string{"string"}),
		}),
		ProjectName: langsmith.F("project_name"),
		Secrets: langsmith.F(map[string]string{
			"foo": "string",
		}),
		BatchSize:                       langsmith.F(int64(1)),
		Commit:                          langsmith.F("commit"),
		DatasetSplits:                   langsmith.F([]string{"string"}),
		EvaluatorRules:                  langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Metadata:                        langsmith.F[any](map[string]interface{}{}),
		Owner:                           langsmith.F("owner"),
		ParallelToolCalls:               langsmith.F(true),
		Repetitions:                     langsmith.F(int64(1)),
		RepoHandle:                      langsmith.F("repo_handle"),
		RepoID:                          langsmith.F("repo_id"),
		RequestsPerSecond:               langsmith.F(int64(0)),
		RunID:                           langsmith.F("run_id"),
		RunnerContext:                   langsmith.F(langsmith.RunnerContextEnumLangsmithUi),
		ToolChoice:                      langsmith.F("tool_choice"),
		Tools:                           langsmith.F([]interface{}{map[string]interface{}{}}),
		UseOrFallbackToWorkspaceSecrets: langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestDatasetPlaygroundExperimentStreamWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.PlaygroundExperiment.Stream(context.TODO(), langsmith.DatasetPlaygroundExperimentStreamParams{
		DatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Manifest:  langsmith.F[any](map[string]interface{}{}),
		Options: langsmith.F(langsmith.RunnableConfigParam{
			Callbacks:      langsmith.F([]interface{}{map[string]interface{}{}}),
			Configurable:   langsmith.F[any](map[string]interface{}{}),
			MaxConcurrency: langsmith.F(int64(0)),
			Metadata:       langsmith.F[any](map[string]interface{}{}),
			RecursionLimit: langsmith.F(int64(0)),
			RunID:          langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			RunName:        langsmith.F("run_name"),
			Tags:           langsmith.F([]string{"string"}),
		}),
		ProjectName: langsmith.F("project_name"),
		Secrets: langsmith.F(map[string]string{
			"foo": "string",
		}),
		Commit:                          langsmith.F("commit"),
		DatasetSplits:                   langsmith.F([]string{"string"}),
		EvaluatorRules:                  langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Metadata:                        langsmith.F[any](map[string]interface{}{}),
		Owner:                           langsmith.F("owner"),
		ParallelToolCalls:               langsmith.F(true),
		Repetitions:                     langsmith.F(int64(1)),
		RepoHandle:                      langsmith.F("repo_handle"),
		RepoID:                          langsmith.F("repo_id"),
		RequestsPerSecond:               langsmith.F(int64(0)),
		RunID:                           langsmith.F("run_id"),
		RunnerContext:                   langsmith.F(langsmith.RunnerContextEnumLangsmithUi),
		ToolChoice:                      langsmith.F("tool_choice"),
		Tools:                           langsmith.F([]interface{}{map[string]interface{}{}}),
		UseOrFallbackToWorkspaceSecrets: langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
