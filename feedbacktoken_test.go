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
	"github.com/langchain-ai/langsmith-go/shared"
)

func TestFeedbackTokenNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.Tokens.New(context.TODO(), langsmith.FeedbackTokenNewParams{
		Body: langsmith.FeedbackIngestTokenCreateSchemaParam{
			FeedbackKey: langsmith.F("feedback_key"),
			RunID:       langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			ExpiresAt:   langsmith.F(time.Now()),
			ExpiresIn: langsmith.F(langsmith.TimedeltaInputParam{
				Days:    langsmith.F(int64(0)),
				Hours:   langsmith.F(int64(0)),
				Minutes: langsmith.F(int64(0)),
			}),
			FeedbackConfig: langsmith.F(langsmith.FeedbackIngestTokenCreateSchemaFeedbackConfigParam{
				Type: langsmith.F(langsmith.FeedbackIngestTokenCreateSchemaFeedbackConfigTypeContinuous),
				Categories: langsmith.F([]langsmith.FeedbackIngestTokenCreateSchemaFeedbackConfigCategoryParam{{
					Value: langsmith.F(0.000000),
					Label: langsmith.F("x"),
				}}),
				Max: langsmith.F(0.000000),
				Min: langsmith.F(0.000000),
			}),
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

func TestFeedbackTokenGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.Tokens.Get(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.FeedbackTokenGetParams{
			Comment:    langsmith.F("comment"),
			Correction: langsmith.F("correction"),
			Score:      langsmith.F[langsmith.FeedbackTokenGetParamsScoreUnion](shared.UnionFloat(0.000000)),
			Value:      langsmith.F[langsmith.FeedbackTokenGetParamsValueUnion](shared.UnionFloat(0.000000)),
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

func TestFeedbackTokenUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.Tokens.Update(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.FeedbackTokenUpdateParams{
			Comment: langsmith.F("comment"),
			Correction: langsmith.F[langsmith.FeedbackTokenUpdateParamsCorrectionUnion](langsmith.FeedbackTokenUpdateParamsCorrectionMap(map[string]interface{}{
				"foo": "bar",
			})),
			Metadata: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Score: langsmith.F[langsmith.FeedbackTokenUpdateParamsScoreUnion](shared.UnionFloat(0.000000)),
			Value: langsmith.F[langsmith.FeedbackTokenUpdateParamsValueUnion](shared.UnionFloat(0.000000)),
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

func TestFeedbackTokenList(t *testing.T) {
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
	_, err := client.Feedback.Tokens.List(context.TODO(), langsmith.FeedbackTokenListParams{
		RunID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
