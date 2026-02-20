// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"bytes"
	"context"
	"errors"
	"io"
	"os"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/shared"
)

func TestExampleNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Examples.New(context.TODO(), langsmith.ExampleNewParams{
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
		Split:                   langsmith.F[langsmith.ExampleNewParamsSplitUnion](langsmith.ExampleNewParamsSplitArray([]string{"string"})),
		UseLegacyMessageFormat:  langsmith.F(true),
		UseSourceRunAttachments: langsmith.F([]string{"string"}),
		UseSourceRunIo:          langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestExampleGetWithOptionalParams(t *testing.T) {
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
	_, err := client.Examples.Get(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.ExampleGetParams{
			AsOf:    langsmith.F[langsmith.ExampleGetParamsAsOfUnion](shared.UnionTime(time.Now())),
			Dataset: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
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

func TestExampleUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Examples.Update(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.ExampleUpdateParams{
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
			Split:     langsmith.F[langsmith.ExampleUpdateParamsSplitUnion](langsmith.ExampleUpdateParamsSplitArray([]string{"string"})),
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

func TestExampleListWithOptionalParams(t *testing.T) {
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
	_, err := client.Examples.List(context.TODO(), langsmith.ExampleListParams{
		ID:               langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		AsOf:             langsmith.F[langsmith.ExampleListParamsAsOfUnion](shared.UnionTime(time.Now())),
		Dataset:          langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Descending:       langsmith.F(true),
		Filter:           langsmith.F("filter"),
		FullTextContains: langsmith.F([]string{"string"}),
		Limit:            langsmith.F(int64(1)),
		Metadata:         langsmith.F("metadata"),
		Offset:           langsmith.F(int64(0)),
		Order:            langsmith.F(langsmith.ExampleListParamsOrderRecent),
		RandomSeed:       langsmith.F(0.000000),
		Select:           langsmith.F([]langsmith.ExampleSelect{langsmith.ExampleSelectID}),
		Splits:           langsmith.F([]string{"string"}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestExampleDelete(t *testing.T) {
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
	_, err := client.Examples.Delete(context.TODO(), "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestExampleDeleteAll(t *testing.T) {
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
	_, err := client.Examples.DeleteAll(context.TODO(), langsmith.ExampleDeleteAllParams{
		ExampleIDs: langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestExampleGetCountWithOptionalParams(t *testing.T) {
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
	_, err := client.Examples.GetCount(context.TODO(), langsmith.ExampleGetCountParams{
		ID:               langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		AsOf:             langsmith.F[langsmith.ExampleGetCountParamsAsOfUnion](shared.UnionTime(time.Now())),
		Dataset:          langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		Filter:           langsmith.F("filter"),
		FullTextContains: langsmith.F([]string{"string"}),
		Metadata:         langsmith.F("metadata"),
		Splits:           langsmith.F([]string{"string"}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestExampleUploadFromCsvWithOptionalParams(t *testing.T) {
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
	_, err := client.Examples.UploadFromCsv(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.ExampleUploadFromCsvParams{
			File:         langsmith.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
			InputKeys:    langsmith.F([]string{"string"}),
			MetadataKeys: langsmith.F([]string{"string"}),
			OutputKeys:   langsmith.F([]string{"string"}),
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
