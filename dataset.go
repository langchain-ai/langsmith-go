// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apiform"
	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// DatasetService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetService] method instead.
type DatasetService struct {
	Options              []option.RequestOption
	Versions             *DatasetVersionService
	Runs                 *DatasetRunService
	Group                *DatasetGroupService
	Experiments          *DatasetExperimentService
	Share                *DatasetShareService
	Comparative          *DatasetComparativeService
	Splits               *DatasetSplitService
	Index                *DatasetIndexService
	PlaygroundExperiment *DatasetPlaygroundExperimentService
}

// NewDatasetService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewDatasetService(opts ...option.RequestOption) (r *DatasetService) {
	r = &DatasetService{}
	r.Options = opts
	r.Versions = NewDatasetVersionService(opts...)
	r.Runs = NewDatasetRunService(opts...)
	r.Group = NewDatasetGroupService(opts...)
	r.Experiments = NewDatasetExperimentService(opts...)
	r.Share = NewDatasetShareService(opts...)
	r.Comparative = NewDatasetComparativeService(opts...)
	r.Splits = NewDatasetSplitService(opts...)
	r.Index = NewDatasetIndexService(opts...)
	r.PlaygroundExperiment = NewDatasetPlaygroundExperimentService(opts...)
	return
}

// Create a new dataset.
func (r *DatasetService) New(ctx context.Context, body DatasetNewParams, opts ...option.RequestOption) (res *Dataset, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a specific dataset.
func (r *DatasetService) Get(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *Dataset, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update a specific dataset.
func (r *DatasetService) Update(ctx context.Context, datasetID string, body DatasetUpdateParams, opts ...option.RequestOption) (res *DatasetUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// Get all datasets by query params and owner.
func (r *DatasetService) List(ctx context.Context, query DatasetListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[Dataset], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/datasets"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Get all datasets by query params and owner.
func (r *DatasetService) ListAutoPaging(ctx context.Context, query DatasetListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[Dataset] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, query, opts...))
}

// Delete a specific dataset.
func (r *DatasetService) Delete(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *DatasetDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Clone a dataset.
func (r *DatasetService) Clone(ctx context.Context, body DatasetCloneParams, opts ...option.RequestOption) (res *[]DatasetCloneResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/clone"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Download a dataset as CSV format.
func (r *DatasetService) GetCsv(ctx context.Context, datasetID string, query DatasetGetCsvParams, opts ...option.RequestOption) (res *DatasetGetCsvResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/csv", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Download a dataset as CSV format.
func (r *DatasetService) GetJSONL(ctx context.Context, datasetID string, query DatasetGetJSONLParams, opts ...option.RequestOption) (res *DatasetGetJSONLResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/jsonl", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Download a dataset as OpenAI Evals Jsonl format.
func (r *DatasetService) GetOpenAI(ctx context.Context, datasetID string, query DatasetGetOpenAIParams, opts ...option.RequestOption) (res *DatasetGetOpenAIResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/openai", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Download a dataset as OpenAI Jsonl format.
func (r *DatasetService) GetOpenAIFt(ctx context.Context, datasetID string, query DatasetGetOpenAIFtParams, opts ...option.RequestOption) (res *DatasetGetOpenAIFtResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/openai_ft", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Get dataset version by as_of or exact tag.
func (r *DatasetService) GetVersion(ctx context.Context, datasetID string, query DatasetGetVersionParams, opts ...option.RequestOption) (res *DatasetVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/version", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Set a tag on a dataset version.
func (r *DatasetService) UpdateTags(ctx context.Context, datasetID string, body DatasetUpdateTagsParams, opts ...option.RequestOption) (res *DatasetVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/tags", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, &res, opts...)
	return
}

// Create a new dataset from a CSV file.
func (r *DatasetService) Upload(ctx context.Context, body DatasetUploadParams, opts ...option.RequestOption) (res *Dataset, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/upload"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Enum for dataset data types.
type DataType string

const (
	DataTypeKv   DataType = "kv"
	DataTypeLlm  DataType = "llm"
	DataTypeChat DataType = "chat"
)

func (r DataType) IsKnown() bool {
	switch r {
	case DataTypeKv, DataTypeLlm, DataTypeChat:
		return true
	}
	return false
}

func (r DataType) implementsDatasetListParamsDataTypeUnion() {}

// Dataset schema.
type Dataset struct {
	ID           string    `json:"id,required" format:"uuid"`
	ModifiedAt   time.Time `json:"modified_at,required" format:"date-time"`
	Name         string    `json:"name,required"`
	SessionCount int64     `json:"session_count,required"`
	TenantID     string    `json:"tenant_id,required" format:"uuid"`
	CreatedAt    time.Time `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                DataType                `json:"data_type,nullable"`
	Description             string                  `json:"description,nullable"`
	ExampleCount            int64                   `json:"example_count,nullable"`
	ExternallyManaged       bool                    `json:"externally_managed,nullable"`
	InputsSchemaDefinition  map[string]interface{}  `json:"inputs_schema_definition,nullable"`
	LastSessionStartTime    time.Time               `json:"last_session_start_time,nullable" format:"date-time"`
	Metadata                map[string]interface{}  `json:"metadata,nullable"`
	OutputsSchemaDefinition map[string]interface{}  `json:"outputs_schema_definition,nullable"`
	Transformations         []DatasetTransformation `json:"transformations,nullable"`
	JSON                    datasetJSON             `json:"-"`
}

// datasetJSON contains the JSON metadata for the struct [Dataset]
type datasetJSON struct {
	ID                      apijson.Field
	ModifiedAt              apijson.Field
	Name                    apijson.Field
	SessionCount            apijson.Field
	TenantID                apijson.Field
	CreatedAt               apijson.Field
	DataType                apijson.Field
	Description             apijson.Field
	ExampleCount            apijson.Field
	ExternallyManaged       apijson.Field
	InputsSchemaDefinition  apijson.Field
	LastSessionStartTime    apijson.Field
	Metadata                apijson.Field
	OutputsSchemaDefinition apijson.Field
	Transformations         apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *Dataset) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetJSON) RawJSON() string {
	return r.raw
}

type DatasetTransformation struct {
	Path []string `json:"path,required"`
	// Enum for dataset transformation types. Ordering determines the order in which
	// transformations are applied if there are multiple transformations on the same
	// path.
	TransformationType DatasetTransformationTransformationType `json:"transformation_type,required"`
	JSON               datasetTransformationJSON               `json:"-"`
}

// datasetTransformationJSON contains the JSON metadata for the struct
// [DatasetTransformation]
type datasetTransformationJSON struct {
	Path               apijson.Field
	TransformationType apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *DatasetTransformation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetTransformationJSON) RawJSON() string {
	return r.raw
}

// Enum for dataset transformation types. Ordering determines the order in which
// transformations are applied if there are multiple transformations on the same
// path.
type DatasetTransformationTransformationType string

const (
	DatasetTransformationTransformationTypeConvertToOpenAIMessage DatasetTransformationTransformationType = "convert_to_openai_message"
	DatasetTransformationTransformationTypeConvertToOpenAITool    DatasetTransformationTransformationType = "convert_to_openai_tool"
	DatasetTransformationTransformationTypeRemoveSystemMessages   DatasetTransformationTransformationType = "remove_system_messages"
	DatasetTransformationTransformationTypeRemoveExtraFields      DatasetTransformationTransformationType = "remove_extra_fields"
	DatasetTransformationTransformationTypeExtractToolsFromRun    DatasetTransformationTransformationType = "extract_tools_from_run"
)

func (r DatasetTransformationTransformationType) IsKnown() bool {
	switch r {
	case DatasetTransformationTransformationTypeConvertToOpenAIMessage, DatasetTransformationTransformationTypeConvertToOpenAITool, DatasetTransformationTransformationTypeRemoveSystemMessages, DatasetTransformationTransformationTypeRemoveExtraFields, DatasetTransformationTransformationTypeExtractToolsFromRun:
		return true
	}
	return false
}

type DatasetTransformationParam struct {
	Path param.Field[[]string] `json:"path,required"`
	// Enum for dataset transformation types. Ordering determines the order in which
	// transformations are applied if there are multiple transformations on the same
	// path.
	TransformationType param.Field[DatasetTransformationTransformationType] `json:"transformation_type,required"`
}

func (r DatasetTransformationParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Dataset version schema.
type DatasetVersion struct {
	AsOf time.Time          `json:"as_of,required" format:"date-time"`
	Tags []string           `json:"tags,nullable"`
	JSON datasetVersionJSON `json:"-"`
}

// datasetVersionJSON contains the JSON metadata for the struct [DatasetVersion]
type datasetVersionJSON struct {
	AsOf        apijson.Field
	Tags        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetVersion) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetVersionJSON) RawJSON() string {
	return r.raw
}

type MissingParam struct {
	Missing param.Field[Missing_Missing] `json:"__missing__,required"`
}

func (r MissingParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r MissingParam) ImplementsDatasetUpdateParamsDescriptionUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsInputsSchemaDefinitionUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsMetadataUnion() {}

func (r MissingParam) ImplementsDatasetUpdateParamsNameUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsOutputsSchemaDefinitionUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsTransformationsUnion() {}

func (r MissingParam) implementsAnnotationQueueUpdateParamsMetadataUnion() {}

func (r MissingParam) ImplementsAnnotationQueueUpdateParamsNumReviewersPerItemUnion() {}

type Missing_Missing string

const (
	Missing_Missing_Missing Missing_Missing = "__missing__"
)

func (r Missing_Missing) IsKnown() bool {
	switch r {
	case Missing_Missing_Missing:
		return true
	}
	return false
}

// Enum for available dataset columns to sort by.
type SortByDatasetColumn string

const (
	SortByDatasetColumnName                 SortByDatasetColumn = "name"
	SortByDatasetColumnCreatedAt            SortByDatasetColumn = "created_at"
	SortByDatasetColumnLastSessionStartTime SortByDatasetColumn = "last_session_start_time"
	SortByDatasetColumnExampleCount         SortByDatasetColumn = "example_count"
	SortByDatasetColumnSessionCount         SortByDatasetColumn = "session_count"
	SortByDatasetColumnModifiedAt           SortByDatasetColumn = "modified_at"
)

func (r SortByDatasetColumn) IsKnown() bool {
	switch r {
	case SortByDatasetColumnName, SortByDatasetColumnCreatedAt, SortByDatasetColumnLastSessionStartTime, SortByDatasetColumnExampleCount, SortByDatasetColumnSessionCount, SortByDatasetColumnModifiedAt:
		return true
	}
	return false
}

type DatasetUpdateResponse struct {
	ID        string    `json:"id,required" format:"uuid"`
	Name      string    `json:"name,required"`
	TenantID  string    `json:"tenant_id,required" format:"uuid"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                DataType                  `json:"data_type,nullable"`
	Description             string                    `json:"description,nullable"`
	ExternallyManaged       bool                      `json:"externally_managed,nullable"`
	InputsSchemaDefinition  map[string]interface{}    `json:"inputs_schema_definition,nullable"`
	OutputsSchemaDefinition map[string]interface{}    `json:"outputs_schema_definition,nullable"`
	Transformations         []DatasetTransformation   `json:"transformations,nullable"`
	JSON                    datasetUpdateResponseJSON `json:"-"`
}

// datasetUpdateResponseJSON contains the JSON metadata for the struct
// [DatasetUpdateResponse]
type datasetUpdateResponseJSON struct {
	ID                      apijson.Field
	Name                    apijson.Field
	TenantID                apijson.Field
	CreatedAt               apijson.Field
	DataType                apijson.Field
	Description             apijson.Field
	ExternallyManaged       apijson.Field
	InputsSchemaDefinition  apijson.Field
	OutputsSchemaDefinition apijson.Field
	Transformations         apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *DatasetUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetDeleteResponse = interface{}

type DatasetCloneResponse map[string]interface{}

type DatasetGetCsvResponse = interface{}

type DatasetGetJSONLResponse = interface{}

type DatasetGetOpenAIResponse = interface{}

type DatasetGetOpenAIFtResponse = interface{}

type DatasetNewParams struct {
	Name      param.Field[string]    `json:"name,required"`
	ID        param.Field[string]    `json:"id" format:"uuid"`
	CreatedAt param.Field[time.Time] `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                param.Field[DataType]                     `json:"data_type"`
	Description             param.Field[string]                       `json:"description"`
	ExternallyManaged       param.Field[bool]                         `json:"externally_managed"`
	Extra                   param.Field[map[string]interface{}]       `json:"extra"`
	InputsSchemaDefinition  param.Field[map[string]interface{}]       `json:"inputs_schema_definition"`
	OutputsSchemaDefinition param.Field[map[string]interface{}]       `json:"outputs_schema_definition"`
	Transformations         param.Field[[]DatasetTransformationParam] `json:"transformations"`
}

func (r DatasetNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetUpdateParams struct {
	Description             param.Field[DatasetUpdateParamsDescriptionUnion]             `json:"description"`
	InputsSchemaDefinition  param.Field[DatasetUpdateParamsInputsSchemaDefinitionUnion]  `json:"inputs_schema_definition"`
	Metadata                param.Field[DatasetUpdateParamsMetadataUnion]                `json:"metadata"`
	Name                    param.Field[DatasetUpdateParamsNameUnion]                    `json:"name"`
	OutputsSchemaDefinition param.Field[DatasetUpdateParamsOutputsSchemaDefinitionUnion] `json:"outputs_schema_definition"`
	PatchExamples           param.Field[map[string]DatasetUpdateParamsPatchExamples]     `json:"patch_examples"`
	Transformations         param.Field[DatasetUpdateParamsTransformationsUnion]         `json:"transformations"`
}

func (r DatasetUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [shared.UnionString], [MissingParam].
type DatasetUpdateParamsDescriptionUnion interface {
	ImplementsDatasetUpdateParamsDescriptionUnion()
}

// Satisfied by [DatasetUpdateParamsInputsSchemaDefinitionMap], [MissingParam].
type DatasetUpdateParamsInputsSchemaDefinitionUnion interface {
	implementsDatasetUpdateParamsInputsSchemaDefinitionUnion()
}

type DatasetUpdateParamsInputsSchemaDefinitionMap map[string]interface{}

func (r DatasetUpdateParamsInputsSchemaDefinitionMap) implementsDatasetUpdateParamsInputsSchemaDefinitionUnion() {
}

// Satisfied by [DatasetUpdateParamsMetadataMap], [MissingParam].
type DatasetUpdateParamsMetadataUnion interface {
	implementsDatasetUpdateParamsMetadataUnion()
}

type DatasetUpdateParamsMetadataMap map[string]interface{}

func (r DatasetUpdateParamsMetadataMap) implementsDatasetUpdateParamsMetadataUnion() {}

// Satisfied by [shared.UnionString], [MissingParam].
type DatasetUpdateParamsNameUnion interface {
	ImplementsDatasetUpdateParamsNameUnion()
}

// Satisfied by [DatasetUpdateParamsOutputsSchemaDefinitionMap], [MissingParam].
type DatasetUpdateParamsOutputsSchemaDefinitionUnion interface {
	implementsDatasetUpdateParamsOutputsSchemaDefinitionUnion()
}

type DatasetUpdateParamsOutputsSchemaDefinitionMap map[string]interface{}

func (r DatasetUpdateParamsOutputsSchemaDefinitionMap) implementsDatasetUpdateParamsOutputsSchemaDefinitionUnion() {
}

// Update class for Example.
type DatasetUpdateParamsPatchExamples struct {
	AttachmentsOperations param.Field[AttachmentsOperationsParam]                 `json:"attachments_operations"`
	DatasetID             param.Field[string]                                     `json:"dataset_id" format:"uuid"`
	Inputs                param.Field[map[string]interface{}]                     `json:"inputs"`
	Metadata              param.Field[map[string]interface{}]                     `json:"metadata"`
	Outputs               param.Field[map[string]interface{}]                     `json:"outputs"`
	Overwrite             param.Field[bool]                                       `json:"overwrite"`
	Split                 param.Field[DatasetUpdateParamsPatchExamplesSplitUnion] `json:"split"`
}

func (r DatasetUpdateParamsPatchExamples) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [DatasetUpdateParamsPatchExamplesSplitArray], [shared.UnionString].
type DatasetUpdateParamsPatchExamplesSplitUnion interface {
	ImplementsDatasetUpdateParamsPatchExamplesSplitUnion()
}

type DatasetUpdateParamsPatchExamplesSplitArray []string

func (r DatasetUpdateParamsPatchExamplesSplitArray) ImplementsDatasetUpdateParamsPatchExamplesSplitUnion() {
}

// Satisfied by [DatasetUpdateParamsTransformationsArray], [MissingParam].
type DatasetUpdateParamsTransformationsUnion interface {
	implementsDatasetUpdateParamsTransformationsUnion()
}

type DatasetUpdateParamsTransformationsArray []DatasetTransformationParam

func (r DatasetUpdateParamsTransformationsArray) implementsDatasetUpdateParamsTransformationsUnion() {
}

type DatasetListParams struct {
	ID param.Field[[]string] `query:"id" format:"uuid"`
	// Enum for dataset data types.
	Datatype                   param.Field[DatasetListParamsDataTypeUnion] `query:"data_type"`
	Exclude                    param.Field[[]DatasetListParamsExclude]     `query:"exclude"`
	ExcludeCorrectionsDatasets param.Field[bool]                           `query:"exclude_corrections_datasets"`
	Limit                      param.Field[int64]                          `query:"limit"`
	Metadata                   param.Field[string]                         `query:"metadata"`
	Name                       param.Field[string]                         `query:"name"`
	NameContains               param.Field[string]                         `query:"name_contains"`
	Offset                     param.Field[int64]                          `query:"offset"`
	// Enum for available dataset columns to sort by.
	SortBy     param.Field[SortByDatasetColumn] `query:"sort_by"`
	SortByDesc param.Field[bool]                `query:"sort_by_desc"`
	TagValueID param.Field[[]string]            `query:"tag_value_id" format:"uuid"`
}

// URLQuery serializes [DatasetListParams]'s query parameters as `url.Values`.
func (r DatasetListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Enum for dataset data types.
//
// Satisfied by [DatasetListParamsDataTypeArray], [DataType].
type DatasetListParamsDataTypeUnion interface {
	implementsDatasetListParamsDataTypeUnion()
}

type DatasetListParamsDataTypeArray []DataType

func (r DatasetListParamsDataTypeArray) implementsDatasetListParamsDataTypeUnion() {}

type DatasetListParamsExclude string

const (
	DatasetListParamsExcludeExampleCount DatasetListParamsExclude = "example_count"
)

func (r DatasetListParamsExclude) IsKnown() bool {
	switch r {
	case DatasetListParamsExcludeExampleCount:
		return true
	}
	return false
}

type DatasetCloneParams struct {
	SourceDatasetID param.Field[string] `json:"source_dataset_id,required" format:"uuid"`
	TargetDatasetID param.Field[string] `json:"target_dataset_id,required" format:"uuid"`
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf     param.Field[DatasetCloneParamsAsOfUnion] `json:"as_of" format:"date-time"`
	Examples param.Field[[]string]                    `json:"examples" format:"uuid"`
}

func (r DatasetCloneParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Only modifications made on or before this time are included. If None, the latest
// version of the dataset is used.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type DatasetCloneParamsAsOfUnion interface {
	ImplementsDatasetCloneParamsAsOfUnion()
}

type DatasetGetCsvParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[time.Time] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [DatasetGetCsvParams]'s query parameters as `url.Values`.
func (r DatasetGetCsvParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetGetJSONLParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[time.Time] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [DatasetGetJSONLParams]'s query parameters as `url.Values`.
func (r DatasetGetJSONLParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetGetOpenAIParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[time.Time] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [DatasetGetOpenAIParams]'s query parameters as `url.Values`.
func (r DatasetGetOpenAIParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetGetOpenAIFtParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[time.Time] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [DatasetGetOpenAIFtParams]'s query parameters as
// `url.Values`.
func (r DatasetGetOpenAIFtParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetGetVersionParams struct {
	AsOf param.Field[time.Time] `query:"as_of" format:"date-time"`
	Tag  param.Field[string]    `query:"tag"`
}

// URLQuery serializes [DatasetGetVersionParams]'s query parameters as
// `url.Values`.
func (r DatasetGetVersionParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetUpdateTagsParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[DatasetUpdateTagsParamsAsOfUnion] `json:"as_of,required" format:"date-time"`
	Tag  param.Field[string]                           `json:"tag,required"`
}

func (r DatasetUpdateTagsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Only modifications made on or before this time are included. If None, the latest
// version of the dataset is used.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type DatasetUpdateTagsParamsAsOfUnion interface {
	ImplementsDatasetUpdateTagsParamsAsOfUnion()
}

type DatasetUploadParams struct {
	File      param.Field[io.Reader] `json:"file,required" format:"binary"`
	InputKeys param.Field[[]string]  `json:"input_keys,required"`
	// Enum for dataset data types.
	DataType                param.Field[DataType] `json:"data_type"`
	Description             param.Field[string]   `json:"description"`
	InputKeyMappings        param.Field[string]   `json:"input_key_mappings"`
	InputsSchemaDefinition  param.Field[string]   `json:"inputs_schema_definition"`
	MetadataKeyMappings     param.Field[string]   `json:"metadata_key_mappings"`
	MetadataKeys            param.Field[[]string] `json:"metadata_keys"`
	Name                    param.Field[string]   `json:"name"`
	OutputKeyMappings       param.Field[string]   `json:"output_key_mappings"`
	OutputKeys              param.Field[[]string] `json:"output_keys"`
	OutputsSchemaDefinition param.Field[string]   `json:"outputs_schema_definition"`
	Transformations         param.Field[string]   `json:"transformations"`
}

func (r DatasetUploadParams) MarshalMultipart() (data []byte, contentType string, err error) {
	buf := bytes.NewBuffer(nil)
	writer := multipart.NewWriter(buf)
	err = apiform.MarshalRoot(r, writer)
	if err != nil {
		writer.Close()
		return nil, "", err
	}
	err = writer.Close()
	if err != nil {
		return nil, "", err
	}
	return buf.Bytes(), writer.FormDataContentType(), nil
}
