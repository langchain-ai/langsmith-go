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

func TestExampleBulkNew(t *testing.T) {
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
	_, err := client.Examples.Bulk.New(context.TODO(), langsmith.ExampleBulkNewParams{
		Body: []langsmith.ExampleBulkNewParamsBody{{
			DatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			ID:        langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			CreatedAt: langsmith.F("created_at"),
			Inputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Metadata: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Outputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			SourceRunID:             langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Split:                   langsmith.F[langsmith.ExampleBulkNewParamsBodySplitUnion](langsmith.ExampleBulkNewParamsBodySplitArray([]string{"string"})),
			UseLegacyMessageFormat:  langsmith.F(true),
			UseSourceRunAttachments: langsmith.F([]string{"string"}),
			UseSourceRunIo:          langsmith.F(true),
		}},
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestExampleBulkPatchAll(t *testing.T) {
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
	_, err := client.Examples.Bulk.PatchAll(context.TODO(), langsmith.ExampleBulkPatchAllParams{
		Body: []langsmith.ExampleBulkPatchAllParamsBody{{
			ID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			AttachmentsOperations: langsmith.F(langsmith.AttachmentsOperationsParam{
				Rename: langsmith.F(map[string]string{
					"foo": "string",
				}),
				Retain: langsmith.F([]string{"string"}),
			}),
			DatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
			Inputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Metadata: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Outputs: langsmith.F(map[string]interface{}{
				"foo": "bar",
			}),
			Overwrite: langsmith.F(true),
			Split:     langsmith.F[langsmith.ExampleBulkPatchAllParamsBodySplitUnion](langsmith.ExampleBulkPatchAllParamsBodySplitArray([]string{"string"})),
		}},
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
