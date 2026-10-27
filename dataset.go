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
	Options        []option.RequestOption
	Examples       *DatasetExampleService
	Versions       *DatasetVersionService
	Runs           *DatasetRunService
	ExperimentRuns *DatasetExperimentRunService
	Share          *DatasetShareService
	Comparative    *DatasetComparativeService
	Splits         *DatasetSplitService
}

// NewDatasetService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewDatasetService(opts ...option.RequestOption) (r *DatasetService) {
	r = &DatasetService{}
	r.Options = opts
	r.Examples = NewDatasetExampleService(opts...)
	r.Versions = NewDatasetVersionService(opts...)
	r.Runs = NewDatasetRunService(opts...)
	r.ExperimentRuns = NewDatasetExperimentRunService(opts...)
	r.Share = NewDatasetShareService(opts...)
	r.Comparative = NewDatasetComparativeService(opts...)
	r.Splits = NewDatasetSplitService(opts...)
	return
}

// Create a new dataset.
func (r *DatasetService) New(ctx context.Context, body DatasetNewParams, opts ...option.RequestOption) (res *Dataset, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a specific dataset.
func (r *DatasetService) Get(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *Dataset, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a dataset.
//
// Only the fields in the request change; omitted fields keep their value. Setting
// a schema or transformations checks every existing example against the result.
// Returns 400, leaving the dataset unchanged, if an example does not match, the
// dataset is not a key-value dataset, or a transformation points into inputs or
// outputs that have no schema. The X-Updated-Examples-Count response header holds
// the number of examples changed through patch_examples. Returns 404 if the
// dataset does not exist and 409 if the new name is already used.
func (r *DatasetService) Update(ctx context.Context, datasetID string, body DatasetUpdateParams, opts ...option.RequestOption) (res *DatasetUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
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
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Copy examples from a source dataset into an existing target dataset.
//
// Each copied example gets a new ID and keeps the inputs, outputs, metadata, and
// source run ID of its source example. Attachments are copied when blob storage is
// enabled. Copying the same examples twice adds them twice. Requires the
// datasets:clone permission on both datasets. Returns the ID, dataset ID, creation
// and modification times, and source run ID of each new example. Without example
// IDs, a source or target that is not found can return an empty list instead of an
// error.
func (r *DatasetService) Clone(ctx context.Context, body DatasetCloneParams, opts ...option.RequestOption) (res *[]DatasetCloneResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/clone"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Download a dataset as CSV format.
func (r *DatasetService) GetCsv(ctx context.Context, datasetID string, query DatasetGetCsvParams, opts ...option.RequestOption) (res *DatasetGetCsvResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/csv", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Download a dataset as CSV format.
func (r *DatasetService) GetJSONL(ctx context.Context, datasetID string, query DatasetGetJSONLParams, opts ...option.RequestOption) (res *DatasetGetJSONLResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/jsonl", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Download a dataset as OpenAI Evals Jsonl format.
func (r *DatasetService) GetOpenAI(ctx context.Context, datasetID string, query DatasetGetOpenAIParams, opts ...option.RequestOption) (res *DatasetGetOpenAIResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/openai", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Download a dataset as OpenAI Jsonl format.
func (r *DatasetService) GetOpenAIFt(ctx context.Context, datasetID string, query DatasetGetOpenAIFtParams, opts ...option.RequestOption) (res *DatasetGetOpenAIFtResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/openai_ft", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get dataset version by as_of or exact tag.
func (r *DatasetService) GetVersion(ctx context.Context, datasetID string, query DatasetGetVersionParams, opts ...option.RequestOption) (res *DatasetVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/version", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Tag a dataset version.
//
// Points the tag at the version given by as_of, creating the tag or moving it if
// the dataset already has it. Other tags are unchanged. Requires permission to
// update datasets. Returns the version's time and all of its tags. Returns 400 if
// the tag is latest, and 404 if the dataset or the as_of tag does not exist.
func (r *DatasetService) UpdateTags(ctx context.Context, datasetID string, body DatasetUpdateTagsParams, opts ...option.RequestOption) (res *DatasetVersion, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/tags", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, &res, opts...)
	return res, err
}

// Create a new dataset from a CSV or JSONL file.
func (r *DatasetService) Upload(ctx context.Context, body DatasetUploadParams, opts ...option.RequestOption) (res *Dataset, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/upload"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
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
	ID                   string    `json:"id" api:"required" format:"uuid"`
	ModifiedAt           time.Time `json:"modified_at" api:"required" format:"date-time"`
	Name                 string    `json:"name" api:"required"`
	SessionCount         int64     `json:"session_count" api:"required"`
	TenantID             string    `json:"tenant_id" api:"required" format:"uuid"`
	BaselineExperimentID string    `json:"baseline_experiment_id" api:"nullable" format:"uuid"`
	CreatedAt            time.Time `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                DataType                `json:"data_type" api:"nullable"`
	Description             string                  `json:"description" api:"nullable"`
	ExampleCount            int64                   `json:"example_count" api:"nullable"`
	ExternallyManaged       bool                    `json:"externally_managed" api:"nullable"`
	InputsSchemaDefinition  map[string]interface{}  `json:"inputs_schema_definition" api:"nullable"`
	LastSessionStartTime    time.Time               `json:"last_session_start_time" api:"nullable" format:"date-time"`
	Metadata                map[string]interface{}  `json:"metadata" api:"nullable"`
	OutputsSchemaDefinition map[string]interface{}  `json:"outputs_schema_definition" api:"nullable"`
	Transformations         []DatasetTransformation `json:"transformations" api:"nullable"`
	JSON                    datasetJSON             `json:"-"`
}

// datasetJSON contains the JSON metadata for the struct [Dataset]
type datasetJSON struct {
	ID                      apijson.Field
	ModifiedAt              apijson.Field
	Name                    apijson.Field
	SessionCount            apijson.Field
	TenantID                apijson.Field
	BaselineExperimentID    apijson.Field
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
	Path []string `json:"path" api:"required"`
	// Enum for dataset transformation types. Ordering determines the order in which
	// transformations are applied if there are multiple transformations on the same
	// path.
	TransformationType DatasetTransformationTransformationType `json:"transformation_type" api:"required"`
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
	Path param.Field[[]string] `json:"path" api:"required"`
	// Enum for dataset transformation types. Ordering determines the order in which
	// transformations are applied if there are multiple transformations on the same
	// path.
	TransformationType param.Field[DatasetTransformationTransformationType] `json:"transformation_type" api:"required"`
}

func (r DatasetTransformationParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Dataset version schema.
type DatasetVersion struct {
	AsOf time.Time          `json:"as_of" api:"required" format:"date-time"`
	Tags []string           `json:"tags" api:"nullable"`
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

type Missing struct {
	Missing Missing_Missing `json:"__missing__" api:"required"`
	JSON    missingJSON     `json:"-"`
}

// missingJSON contains the JSON metadata for the struct [Missing]
type missingJSON struct {
	Missing     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *Missing) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r missingJSON) RawJSON() string {
	return r.raw
}

func (r Missing) ImplementsAnnotationQueueRubricItemSchemaRegexValidatorUnion() {}

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

type MissingParam struct {
	Missing param.Field[Missing_Missing] `json:"__missing__" api:"required"`
}

func (r MissingParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r MissingParam) ImplementsDatasetUpdateParamsBaselineExperimentIDUnion() {}

func (r MissingParam) ImplementsDatasetUpdateParamsDescriptionUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsInputsSchemaDefinitionUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsMetadataUnion() {}

func (r MissingParam) ImplementsDatasetUpdateParamsNameUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsOutputsSchemaDefinitionUnion() {}

func (r MissingParam) implementsDatasetUpdateParamsTransformationsUnion() {}

func (r MissingParam) ImplementsAnnotationQueueRubricItemSchemaRegexValidatorUnionParam() {}

func (r MissingParam) ImplementsAnnotationQueueUpdateParamsEnableReservationsUnion() {}

func (r MissingParam) implementsAnnotationQueueUpdateParamsMetadataUnion() {}

func (r MissingParam) ImplementsAnnotationQueueUpdateParamsNumReviewersPerItemUnion() {}

func (r MissingParam) ImplementsAnnotationQueueUpdateParamsReservationMinutesUnion() {}

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
	ID        string    `json:"id" api:"required" format:"uuid"`
	Name      string    `json:"name" api:"required"`
	TenantID  string    `json:"tenant_id" api:"required" format:"uuid"`
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                DataType                  `json:"data_type" api:"nullable"`
	Description             string                    `json:"description" api:"nullable"`
	ExternallyManaged       bool                      `json:"externally_managed" api:"nullable"`
	InputsSchemaDefinition  map[string]interface{}    `json:"inputs_schema_definition" api:"nullable"`
	OutputsSchemaDefinition map[string]interface{}    `json:"outputs_schema_definition" api:"nullable"`
	Transformations         []DatasetTransformation   `json:"transformations" api:"nullable"`
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
	Name      param.Field[string]    `json:"name" api:"required"`
	ID        param.Field[string]    `json:"id" format:"uuid"`
	CreatedAt param.Field[time.Time] `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                param.Field[DataType]                     `json:"data_type"`
	Description             param.Field[string]                       `json:"description"`
	ExternallyManaged       param.Field[bool]                         `json:"externally_managed"`
	Extra                   param.Field[map[string]interface{}]       `json:"extra"`
	InputsSchemaDefinition  param.Field[map[string]interface{}]       `json:"inputs_schema_definition"`
	OutputsSchemaDefinition param.Field[map[string]interface{}]       `json:"outputs_schema_definition"`
	TagValueIDs             param.Field[[]string]                     `json:"tag_value_ids" format:"uuid"`
	Transformations         param.Field[[]DatasetTransformationParam] `json:"transformations"`
}

func (r DatasetNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetUpdateParams struct {
	// ID of the experiment used as the baseline when comparing experiments on this
	// dataset. Must be an experiment run on this dataset, otherwise returns 400. Null
	// clears it. Omit to keep the current baseline.
	BaselineExperimentID param.Field[DatasetUpdateParamsBaselineExperimentIDUnion] `json:"baseline_experiment_id" format:"uuid"`
	// New dataset description. Null clears it. Omit to keep the current description.
	Description param.Field[DatasetUpdateParamsDescriptionUnion] `json:"description"`
	// JSON Schema that example inputs must match. Replaces the current inputs schema
	// as a whole; null removes it, which returns 400 while transformations still point
	// into inputs. Only key-value datasets support schemas. Every existing example
	// must match the new schema, otherwise the request returns 400 with the IDs of the
	// examples that do not, and nothing is changed. Omit to keep the current schema.
	InputsSchemaDefinition param.Field[DatasetUpdateParamsInputsSchemaDefinitionUnion] `json:"inputs_schema_definition"`
	// Dataset metadata as a JSON object. Replaces the current metadata as a whole;
	// keys are not merged. Send an empty object to clear it. Omit to keep the current
	// metadata.
	Metadata param.Field[DatasetUpdateParamsMetadataUnion] `json:"metadata"`
	// New dataset name. Must be unique in the workspace; a name already used by
	// another dataset returns 409. Null returns 400. Omit to keep the current name.
	Name param.Field[DatasetUpdateParamsNameUnion] `json:"name"`
	// JSON Schema that example outputs must match. Replaces the current outputs schema
	// as a whole; null removes it, which returns 400 while transformations still point
	// into outputs. Only key-value datasets support schemas. Every existing example
	// must match the new schema, otherwise the request returns 400 with the IDs of the
	// examples that do not, and nothing is changed. Omit to keep the current schema.
	OutputsSchemaDefinition param.Field[DatasetUpdateParamsOutputsSchemaDefinitionUnion] `json:"outputs_schema_definition"`
	// Example updates keyed by example ID, applied in the same transaction so existing
	// examples can be fixed to match a new schema. Only used when the request also
	// sets inputs_schema_definition, outputs_schema_definition, or transformations and
	// the dataset ends up with at least one schema; otherwise it is ignored. Each
	// update replaces the inputs, outputs, metadata, or split it sets and keeps the
	// rest; attachment operations are ignored. An ID that is not an example in this
	// dataset returns 404.
	PatchExamples param.Field[map[string]DatasetUpdateParamsPatchExamples] `json:"patch_examples"`
	// Transformations applied to new examples, each a path into inputs or outputs and
	// a transformation type. Replaces the whole current list; null or an empty list
	// removes all of them. Existing examples are not transformed. Requires a schema on
	// the inputs or outputs the path points into, otherwise returns 400. Omit to keep
	// the current list.
	Transformations param.Field[DatasetUpdateParamsTransformationsUnion] `json:"transformations"`
}

func (r DatasetUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// ID of the experiment used as the baseline when comparing experiments on this
// dataset. Must be an experiment run on this dataset, otherwise returns 400. Null
// clears it. Omit to keep the current baseline.
//
// Satisfied by [shared.UnionString], [MissingParam].
type DatasetUpdateParamsBaselineExperimentIDUnion interface {
	ImplementsDatasetUpdateParamsBaselineExperimentIDUnion()
}

// New dataset description. Null clears it. Omit to keep the current description.
//
// Satisfied by [shared.UnionString], [MissingParam].
type DatasetUpdateParamsDescriptionUnion interface {
	ImplementsDatasetUpdateParamsDescriptionUnion()
}

// JSON Schema that example inputs must match. Replaces the current inputs schema
// as a whole; null removes it, which returns 400 while transformations still point
// into inputs. Only key-value datasets support schemas. Every existing example
// must match the new schema, otherwise the request returns 400 with the IDs of the
// examples that do not, and nothing is changed. Omit to keep the current schema.
//
// Satisfied by [DatasetUpdateParamsInputsSchemaDefinitionMap], [MissingParam].
type DatasetUpdateParamsInputsSchemaDefinitionUnion interface {
	implementsDatasetUpdateParamsInputsSchemaDefinitionUnion()
}

type DatasetUpdateParamsInputsSchemaDefinitionMap map[string]interface{}

func (r DatasetUpdateParamsInputsSchemaDefinitionMap) implementsDatasetUpdateParamsInputsSchemaDefinitionUnion() {
}

// Dataset metadata as a JSON object. Replaces the current metadata as a whole;
// keys are not merged. Send an empty object to clear it. Omit to keep the current
// metadata.
//
// Satisfied by [DatasetUpdateParamsMetadataMap], [MissingParam].
type DatasetUpdateParamsMetadataUnion interface {
	implementsDatasetUpdateParamsMetadataUnion()
}

type DatasetUpdateParamsMetadataMap map[string]interface{}

func (r DatasetUpdateParamsMetadataMap) implementsDatasetUpdateParamsMetadataUnion() {}

// New dataset name. Must be unique in the workspace; a name already used by
// another dataset returns 409. Null returns 400. Omit to keep the current name.
//
// Satisfied by [shared.UnionString], [MissingParam].
type DatasetUpdateParamsNameUnion interface {
	ImplementsDatasetUpdateParamsNameUnion()
}

// JSON Schema that example outputs must match. Replaces the current outputs schema
// as a whole; null removes it, which returns 400 while transformations still point
// into outputs. Only key-value datasets support schemas. Every existing example
// must match the new schema, otherwise the request returns 400 with the IDs of the
// examples that do not, and nothing is changed. Omit to keep the current schema.
//
// Satisfied by [DatasetUpdateParamsOutputsSchemaDefinitionMap], [MissingParam].
type DatasetUpdateParamsOutputsSchemaDefinitionUnion interface {
	implementsDatasetUpdateParamsOutputsSchemaDefinitionUnion()
}

type DatasetUpdateParamsOutputsSchemaDefinitionMap map[string]interface{}

func (r DatasetUpdateParamsOutputsSchemaDefinitionMap) implementsDatasetUpdateParamsOutputsSchemaDefinitionUnion() {
}

// Update class for Example.
type DatasetUpdateParamsPatchExamples struct {
	// Attachments to rename or keep.
	AttachmentsOperations param.Field[AttachmentsOperationsParam] `json:"attachments_operations"`
	// The example's current dataset. An example can't move to another dataset.
	DatasetID param.Field[string] `json:"dataset_id" format:"uuid"`
	// New inputs. Replaces the existing inputs; omit to keep them.
	Inputs param.Field[map[string]interface{}] `json:"inputs"`
	// New metadata. Replaces the existing metadata; omit to keep it.
	Metadata param.Field[map[string]interface{}] `json:"metadata"`
	// New outputs. Replaces the existing outputs; omit to keep them.
	Outputs param.Field[map[string]interface{}] `json:"outputs"`
	// Not used.
	Overwrite param.Field[bool] `json:"overwrite"`
	// Split or splits to put the example in. Replaces its current splits.
	Split param.Field[DatasetUpdateParamsPatchExamplesSplitUnion] `json:"split"`
}

func (r DatasetUpdateParamsPatchExamples) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Split or splits to put the example in. Replaces its current splits.
//
// Satisfied by [DatasetUpdateParamsPatchExamplesSplitArray], [shared.UnionString].
type DatasetUpdateParamsPatchExamplesSplitUnion interface {
	ImplementsDatasetUpdateParamsPatchExamplesSplitUnion()
}

type DatasetUpdateParamsPatchExamplesSplitArray []string

func (r DatasetUpdateParamsPatchExamplesSplitArray) ImplementsDatasetUpdateParamsPatchExamplesSplitUnion() {
}

// Transformations applied to new examples, each a path into inputs or outputs and
// a transformation type. Replaces the whole current list; null or an empty list
// removes all of them. Existing examples are not transformed. Requires a schema on
// the inputs or outputs the path points into, otherwise returns 400. Omit to keep
// the current list.
//
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
	// ID of the dataset in this workspace to copy examples from. It is not modified.
	SourceDatasetID param.Field[string] `json:"source_dataset_id" api:"required" format:"uuid"`
	// ID of the dataset in this workspace that receives the copied examples. Its
	// current examples are kept.
	TargetDatasetID param.Field[string] `json:"target_dataset_id" api:"required" format:"uuid"`
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[DatasetCloneParamsAsOfUnion] `json:"as_of" format:"date-time"`
	// IDs of source examples to copy. Defaults to every example in the source version.
	// When IDs are given and nothing is copied, because none match or either dataset
	// is not found, the request fails with 404.
	Examples param.Field[[]string] `json:"examples" format:"uuid"`
	// Split name, or list of split names, assigned to every copied example in place of
	// its source splits. Defaults to keeping the source splits.
	Split param.Field[DatasetCloneParamsSplitUnion] `json:"split"`
	// List of resource tag value IDs to attach to the target dataset, up to 100.
	// Requires the datasets:tag_on_create permission.
	TagValueIDs param.Field[[]string] `json:"tag_value_ids" format:"uuid"`
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

// Split name, or list of split names, assigned to every copied example in place of
// its source splits. Defaults to keeping the source splits.
//
// Satisfied by [shared.UnionString], [DatasetCloneParamsSplitArray].
type DatasetCloneParamsSplitUnion interface {
	ImplementsDatasetCloneParamsSplitUnion()
}

type DatasetCloneParamsSplitArray []string

func (r DatasetCloneParamsSplitArray) ImplementsDatasetCloneParamsSplitUnion() {}

type DatasetGetCsvParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[time.Time] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [DatasetGetCsvParams]'s query parameters as `url.Values`.
func (r DatasetGetCsvParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetUpdateTagsParams struct {
	// The version to tag: an ISO 8601 timestamp, or the name of an existing tag on
	// this dataset (latest resolves to the most recent example change, and returns 404
	// for a dataset with no examples). A value that does not parse as a timestamp is
	// treated as a tag name, and an unknown tag name returns 404. Timestamps are not
	// checked against the dataset history: the tag covers every example change made on
	// or before that time, so a future time keeps picking up new changes until it
	// passes. Use the exact as_of of an existing version, since looking up a version
	// by tag returns 404 when the tag's time does not match a version.
	AsOf param.Field[DatasetUpdateTagsParamsAsOfUnion] `json:"as_of" api:"required" format:"date-time"`
	// The tag to set. If the dataset already has this tag, it is moved to the new
	// version. The tag latest is reserved and returns 400.
	Tag param.Field[string] `json:"tag" api:"required"`
}

func (r DatasetUpdateTagsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The version to tag: an ISO 8601 timestamp, or the name of an existing tag on
// this dataset (latest resolves to the most recent example change, and returns 404
// for a dataset with no examples). A value that does not parse as a timestamp is
// treated as a tag name, and an unknown tag name returns 404. Timestamps are not
// checked against the dataset history: the tag covers every example change made on
// or before that time, so a future time keeps picking up new changes until it
// passes. Use the exact as_of of an existing version, since looking up a version
// by tag returns 404 when the tag's time does not match a version.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type DatasetUpdateTagsParamsAsOfUnion interface {
	ImplementsDatasetUpdateTagsParamsAsOfUnion()
}

type DatasetUploadParams struct {
	File      param.Field[io.Reader] `json:"file" api:"required" format:"binary"`
	InputKeys param.Field[[]string]  `json:"input_keys" api:"required"`
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
	TagValueIDs             param.Field[string]   `json:"tag_value_ids"`
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
