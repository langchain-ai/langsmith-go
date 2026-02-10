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

func TestFeedbackNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.New(context.TODO(), langsmith.FeedbackNewParams{
		FeedbackCreateSchema: langsmith.FeedbackCreateSchemaParam{
			Key:                     langsmith.F("key"),
			ID:                      langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Comment:                 langsmith.F("comment"),
			ComparativeExperimentID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Correction: langsmith.F[langsmith.FeedbackCreateSchemaCorrectionUnionParam](langsmith.FeedbackCreateSchemaCorrectionMapParam(map[string]interface{}{
				"foo": "bar",
			})),
			CreatedAt: langsmith.F(time.Now()),
			Error:     langsmith.F(true),
			FeedbackConfig: langsmith.F(langsmith.FeedbackCreateSchemaFeedbackConfigParam{
				Type: langsmith.F(langsmith.FeedbackCreateSchemaFeedbackConfigTypeContinuous),
				Categories: langsmith.F([]langsmith.FeedbackCreateSchemaFeedbackConfigCategoryParam{{
					Value: langsmith.F(0.000000),
					Label: langsmith.F("x"),
				}}),
				Max: langsmith.F(0.000000),
				Min: langsmith.F(0.000000),
			}),
			FeedbackGroupID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			FeedbackSource: langsmith.F[langsmith.FeedbackCreateSchemaFeedbackSourceUnionParam](langsmith.AppFeedbackSourceParam{
				Metadata: langsmith.F(map[string]interface{}{
					"foo": "bar",
				}),
				Type: langsmith.F(langsmith.AppFeedbackSourceTypeApp),
			}),
			ModifiedAt: langsmith.F(time.Now()),
			RunID:      langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Score:      langsmith.F[langsmith.FeedbackCreateSchemaScoreUnionParam](shared.UnionFloat(0.000000)),
			SessionID:  langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			StartTime:  langsmith.F(time.Now()),
			TraceID:    langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Value:      langsmith.F[langsmith.FeedbackCreateSchemaValueUnionParam](shared.UnionFloat(0.000000)),
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

func TestFeedbackGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.Get(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.FeedbackGetParams{
			IncludeUserNames: langsmith.F(true),
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

func TestFeedbackUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.Update(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.FeedbackUpdateParams{
			Comment: langsmith.F("comment"),
			Correction: langsmith.F[langsmith.FeedbackUpdateParamsCorrectionUnion](langsmith.FeedbackUpdateParamsCorrectionMap(map[string]interface{}{
				"foo": "bar",
			})),
			FeedbackConfig: langsmith.F(langsmith.FeedbackUpdateParamsFeedbackConfig{
				Type: langsmith.F(langsmith.FeedbackUpdateParamsFeedbackConfigTypeContinuous),
				Categories: langsmith.F([]langsmith.FeedbackUpdateParamsFeedbackConfigCategory{{
					Value: langsmith.F(0.000000),
					Label: langsmith.F("x"),
				}}),
				Max: langsmith.F(0.000000),
				Min: langsmith.F(0.000000),
			}),
			Score: langsmith.F[langsmith.FeedbackUpdateParamsScoreUnion](shared.UnionFloat(0.000000)),
			Value: langsmith.F[langsmith.FeedbackUpdateParamsValueUnion](shared.UnionFloat(0.000000)),
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

func TestFeedbackListWithOptionalParams(t *testing.T) {
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
	_, err := client.Feedback.List(context.TODO(), langsmith.FeedbackListParams{
		ComparativeExperimentID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		HasComment:              langsmith.F(true),
		HasScore:                langsmith.F(true),
		IncludeUserNames:        langsmith.F(true),
		Key:                     langsmith.F([]string{"string"}),
		Level:                   langsmith.F(langsmith.FeedbackLevelRun),
		Limit:                   langsmith.F(int64(1)),
		MaxCreatedAt:            langsmith.F(time.Now()),
		MinCreatedAt:            langsmith.F(time.Now()),
		Offset:                  langsmith.F(int64(0)),
		Run:                     langsmith.F[langsmith.FeedbackListParamsRunUnion](langsmith.FeedbackListParamsRunArray([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"})),
		Session:                 langsmith.F[langsmith.FeedbackListParamsSessionUnion](langsmith.FeedbackListParamsSessionArray([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"})),
		Source:                  langsmith.F([]langsmith.SourceType{langsmith.SourceTypeAPI}),
		User:                    langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestFeedbackDelete(t *testing.T) {
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
	_, err := client.Feedback.Delete(context.TODO(), "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
