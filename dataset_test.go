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

func TestDatasetNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.New(context.TODO(), langsmith.DatasetNewParams{
		Name:              langsmith.F("name"),
		ID:                langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		CreatedAt:         langsmith.F(time.Now()),
		DataType:          langsmith.F(langsmith.DataTypeKv),
		Description:       langsmith.F("description"),
		ExternallyManaged: langsmith.F(true),
		Extra: langsmith.F(map[string]interface{}{
			"foo": "bar",
		}),
		InputsSchemaDefinition: langsmith.F(map[string]interface{}{
			"foo": "bar",
		}),
		OutputsSchemaDefinition: langsmith.F(map[string]interface{}{
			"foo": "bar",
		}),
		Transformations: langsmith.F([]langsmith.DatasetTransformationParam{{
			Path:               langsmith.F([]string{"string"}),
			TransformationType: langsmith.F(langsmith.DatasetTransformationTransformationTypeConvertToOpenAIMessage),
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

func TestDatasetGet(t *testing.T) {
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
	_, err := client.Datasets.Get(context.TODO(), "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestDatasetUpdateWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.Update(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetUpdateParams{
			BaselineExperimentID: langsmith.F[langsmith.DatasetUpdateParamsBaselineExperimentIDUnion](shared.UnionString("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")),
			Description:          langsmith.F[langsmith.DatasetUpdateParamsDescriptionUnion](shared.UnionString("string")),
			InputsSchemaDefinition: langsmith.F[langsmith.DatasetUpdateParamsInputsSchemaDefinitionUnion](langsmith.DatasetUpdateParamsInputsSchemaDefinitionMap(map[string]interface{}{
				"foo": "bar",
			})),
			Metadata: langsmith.F[langsmith.DatasetUpdateParamsMetadataUnion](langsmith.DatasetUpdateParamsMetadataMap(map[string]interface{}{
				"foo": "bar",
			})),
			Name: langsmith.F[langsmith.DatasetUpdateParamsNameUnion](shared.UnionString("string")),
			OutputsSchemaDefinition: langsmith.F[langsmith.DatasetUpdateParamsOutputsSchemaDefinitionUnion](langsmith.DatasetUpdateParamsOutputsSchemaDefinitionMap(map[string]interface{}{
				"foo": "bar",
			})),
			PatchExamples: langsmith.F(map[string]langsmith.DatasetUpdateParamsPatchExamples{
				"foo": {
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
					Split:     langsmith.F[langsmith.DatasetUpdateParamsPatchExamplesSplitUnion](langsmith.DatasetUpdateParamsPatchExamplesSplitArray([]string{"string"})),
				},
			}),
			Transformations: langsmith.F[langsmith.DatasetUpdateParamsTransformationsUnion](langsmith.DatasetUpdateParamsTransformationsArray([]langsmith.DatasetTransformationParam{{
				Path:               langsmith.F([]string{"string"}),
				TransformationType: langsmith.F(langsmith.DatasetTransformationTransformationTypeConvertToOpenAIMessage),
			}})),
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

func TestDatasetListWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.List(context.TODO(), langsmith.DatasetListParams{
		ID:                         langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Datatype:                   langsmith.F[langsmith.DatasetListParamsDataTypeUnion](langsmith.DatasetListParamsDataTypeArray([]langsmith.DataType{langsmith.DataTypeKv})),
		Exclude:                    langsmith.F([]langsmith.DatasetListParamsExclude{langsmith.DatasetListParamsExcludeExampleCount}),
		ExcludeCorrectionsDatasets: langsmith.F(true),
		Limit:                      langsmith.F(int64(1)),
		Metadata:                   langsmith.F("metadata"),
		Name:                       langsmith.F("name"),
		NameContains:               langsmith.F("name_contains"),
		Offset:                     langsmith.F(int64(0)),
		SortBy:                     langsmith.F(langsmith.SortByDatasetColumnName),
		SortByDesc:                 langsmith.F(true),
		TagValueID:                 langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestDatasetDelete(t *testing.T) {
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
	_, err := client.Datasets.Delete(context.TODO(), "182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestDatasetCloneWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.Clone(context.TODO(), langsmith.DatasetCloneParams{
		SourceDatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		TargetDatasetID: langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		AsOf:            langsmith.F[langsmith.DatasetCloneParamsAsOfUnion](shared.UnionTime(time.Now())),
		Examples:        langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Split:           langsmith.F[langsmith.DatasetCloneParamsSplitUnion](shared.UnionString("string")),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestDatasetGetCsvWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.GetCsv(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetGetCsvParams{
			AsOf: langsmith.F(time.Now()),
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

func TestDatasetGetJSONLWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.GetJSONL(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetGetJSONLParams{
			AsOf: langsmith.F(time.Now()),
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

func TestDatasetGetOpenAIWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.GetOpenAI(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetGetOpenAIParams{
			AsOf: langsmith.F(time.Now()),
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

func TestDatasetGetOpenAIFtWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.GetOpenAIFt(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetGetOpenAIFtParams{
			AsOf: langsmith.F(time.Now()),
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

func TestDatasetGetVersionWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.GetVersion(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetGetVersionParams{
			AsOf: langsmith.F(time.Now()),
			Tag:  langsmith.F("tag"),
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

func TestDatasetUpdateTags(t *testing.T) {
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
	_, err := client.Datasets.UpdateTags(
		context.TODO(),
		"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e",
		langsmith.DatasetUpdateTagsParams{
			AsOf: langsmith.F[langsmith.DatasetUpdateTagsParamsAsOfUnion](shared.UnionTime(time.Now())),
			Tag:  langsmith.F("tag"),
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

func TestDatasetUploadWithOptionalParams(t *testing.T) {
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
	_, err := client.Datasets.Upload(context.TODO(), langsmith.DatasetUploadParams{
		File:                    langsmith.F(io.Reader(bytes.NewBuffer([]byte("some file contents")))),
		InputKeys:               langsmith.F([]string{"string"}),
		DataType:                langsmith.F(langsmith.DataTypeKv),
		Description:             langsmith.F("description"),
		InputKeyMappings:        langsmith.F("input_key_mappings"),
		InputsSchemaDefinition:  langsmith.F("inputs_schema_definition"),
		MetadataKeyMappings:     langsmith.F("metadata_key_mappings"),
		MetadataKeys:            langsmith.F([]string{"string"}),
		Name:                    langsmith.F("name"),
		OutputKeyMappings:       langsmith.F("output_key_mappings"),
		OutputKeys:              langsmith.F([]string{"string"}),
		OutputsSchemaDefinition: langsmith.F("outputs_schema_definition"),
		Transformations:         langsmith.F("transformations"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
