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

func TestOnlineEvaluatorNewWithOptionalParams(t *testing.T) {
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
	_, err := client.OnlineEvaluators.New(context.TODO(), langsmith.OnlineEvaluatorNewParams{
		CreateOnlineEvaluatorRequest: langsmith.CreateOnlineEvaluatorRequestParam{
			CodeEvaluator: langsmith.F(langsmith.CreateOnlineCodeEvaluatorRequestParam{
				Code:     langsmith.F("code"),
				Language: langsmith.F("language"),
			}),
			LlmEvaluator: langsmith.F(langsmith.CreateOnlineLlmEvaluatorRequestParam{
				CommitHashOrTag:  langsmith.F("commit_hash_or_tag"),
				PromptRepoHandle: langsmith.F("prompt_repo_handle"),
				VariableMapping:  langsmith.F[any](map[string]interface{}{}),
			}),
			Name: langsmith.F("name"),
			Type: langsmith.F(langsmith.OnlineEvaluatorTypeLlm),
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

func TestOnlineEvaluatorGet(t *testing.T) {
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
	_, err := client.OnlineEvaluators.Get(context.TODO(), "evaluator_id")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOnlineEvaluatorUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.OnlineEvaluators.Update(
		context.TODO(),
		"evaluator_id",
		langsmith.OnlineEvaluatorUpdateParams{
			UpdateOnlineEvaluatorRequest: langsmith.UpdateOnlineEvaluatorRequestParam{
				CodeEvaluator: langsmith.F(langsmith.UpdateOnlineCodeEvaluatorRequestParam{
					Code:     langsmith.F("code"),
					Language: langsmith.F("language"),
				}),
				LlmEvaluator: langsmith.F(langsmith.UpdateOnlineLlmEvaluatorRequestParam{
					CommitHashOrTag:       langsmith.F("commit_hash_or_tag"),
					NumFewShotExamples:    langsmith.F(int64(0)),
					PromptRepoHandle:      langsmith.F("prompt_repo_handle"),
					UseCorrectionsDataset: langsmith.F(true),
					VariableMapping:       langsmith.F[any](map[string]interface{}{}),
				}),
				Name: langsmith.F("name"),
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

func TestOnlineEvaluatorListWithOptionalParams(t *testing.T) {
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
	_, err := client.OnlineEvaluators.List(context.TODO(), langsmith.OnlineEvaluatorListParams{
		FeedbackKey:  langsmith.F("feedback_key"),
		Limit:        langsmith.F(int64(0)),
		NameContains: langsmith.F("name_contains"),
		Offset:       langsmith.F(int64(0)),
		ResourceID:   langsmith.F([]string{"string"}),
		SortBy:       langsmith.F("sort_by"),
		SortByDesc:   langsmith.F(true),
		TagValueID:   langsmith.F([]string{"string"}),
		Type:         langsmith.F("type"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOnlineEvaluatorDeleteWithOptionalParams(t *testing.T) {
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
	err := client.OnlineEvaluators.Delete(
		context.TODO(),
		"evaluator_id",
		langsmith.OnlineEvaluatorDeleteParams{
			DeleteRunRules: langsmith.F(true),
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

func TestOnlineEvaluatorBulkDeleteWithOptionalParams(t *testing.T) {
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
	_, err := client.OnlineEvaluators.BulkDelete(context.TODO(), langsmith.OnlineEvaluatorBulkDeleteParams{
		EvaluatorIDs:   langsmith.F([]string{"string"}),
		DeleteRunRules: langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestOnlineEvaluatorSpendWithOptionalParams(t *testing.T) {
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
	_, err := client.OnlineEvaluators.Spend(context.TODO(), langsmith.OnlineEvaluatorSpendParams{
		PeriodStart: langsmith.F("period_start"),
		DatasetID:   langsmith.F("dataset_id"),
		EvaluatorID: langsmith.F("evaluator_id"),
		FeedbackKey: langsmith.F("feedback_key"),
		GroupBy:     langsmith.F("group_by"),
		ResourceID:  langsmith.F([]string{"string"}),
		SessionID:   langsmith.F("session_id"),
		Type:        langsmith.F("type"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
