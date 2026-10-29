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

func TestEvaluatorListWithOptionalParams(t *testing.T) {
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
	_, err := client.Evaluators.List(context.TODO(), langsmith.EvaluatorListParams{
		ID:                      langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e", "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		DatasetID:               langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		EvaluatorID:             langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		IncludeBackfillProgress: langsmith.F(true),
		NameContains:            langsmith.F("name_contains"),
		SessionID:               langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		TagValueID:              langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e", "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Type:                    langsmith.F(langsmith.EvaluatorListParamsTypeSession),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestEvaluatorNewDecisionModelWithOptionalParams(t *testing.T) {
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
	_, err := client.Evaluators.NewDecisionModel(context.TODO(), langsmith.EvaluatorNewDecisionModelParams{
		DisplayName: langsmith.F("display_name"),
		Prompt:      langsmith.F([][]interface{}{{map[string]interface{}{}, map[string]interface{}{}}}),
		Questions: langsmith.F(map[string]langsmith.EvaluatorNewDecisionModelParamsQuestions{
			"foo": {
				Instructions: langsmith.F("instructions"),
				Type:         langsmith.F(langsmith.EvaluatorNewDecisionModelParamsQuestionsTypeNoul),
				Levels:       langsmith.F([]string{"string"}),
				No:           langsmith.F("no"),
				Options:      langsmith.F[langsmith.EvaluatorNewDecisionModelParamsQuestionsOptionsUnion](langsmith.EvaluatorNewDecisionModelParamsQuestionsOptionsArray([]string{"string"})),
				Yes:          langsmith.F("yes"),
			},
		}),
		BaseURL:    langsmith.F("base_url"),
		DatasetID:  langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Filter:     langsmith.F("filter"),
		MaxRetries: langsmith.F(int64(0)),
		Model:      langsmith.F("model"),
		ModelHeaders: langsmith.F(map[string]string{
			"foo": "string",
		}),
		PlaygroundSettingsID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Provider:             langsmith.F(langsmith.EvaluatorNewDecisionModelParamsProviderTypesafe),
		SamplingRate:         langsmith.F(0.000000),
		SessionID:            langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		TemplateFormat:       langsmith.F(langsmith.EvaluatorNewDecisionModelParamsTemplateFormatMustache),
		TimeoutSeconds:       langsmith.F(int64(0)),
		VariableMapping: langsmith.F(map[string]string{
			"foo": "string",
		}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
