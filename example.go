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

// ExampleService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewExampleService] method instead.
type ExampleService struct {
	Options  []option.RequestOption
	Bulk     *ExampleBulkService
	Validate *ExampleValidateService
}

// NewExampleService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewExampleService(opts ...option.RequestOption) (r *ExampleService) {
	r = &ExampleService{}
	r.Options = opts
	r.Bulk = NewExampleBulkService(opts...)
	r.Validate = NewExampleValidateService(opts...)
	return
}

// Create a new example.
func (r *ExampleService) New(ctx context.Context, body ExampleNewParams, opts ...option.RequestOption) (res *Example, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a specific example.
func (r *ExampleService) Get(ctx context.Context, exampleID string, query ExampleGetParams, opts ...option.RequestOption) (res *Example, err error) {
	opts = slices.Concat(r.Options, opts)
	if exampleID == "" {
		err = errors.New("missing required example_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/examples/%s", exampleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Update a specific example.
func (r *ExampleService) Update(ctx context.Context, exampleID string, body ExampleUpdateParams, opts ...option.RequestOption) (res *ExampleUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if exampleID == "" {
		err = errors.New("missing required example_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/examples/%s", exampleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// Get all examples by query params
func (r *ExampleService) List(ctx context.Context, query ExampleListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[Example], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/examples"
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

// Get all examples by query params
func (r *ExampleService) ListAutoPaging(ctx context.Context, query ExampleListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[Example] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, query, opts...))
}

// Soft delete an example. Only deletes the example in the 'latest' version of the
// dataset.
func (r *ExampleService) Delete(ctx context.Context, exampleID string, opts ...option.RequestOption) (res *ExampleDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if exampleID == "" {
		err = errors.New("missing required example_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/examples/%s", exampleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Soft delete examples. Only deletes the examples in the 'latest' version of the
// dataset.
func (r *ExampleService) DeleteAll(ctx context.Context, body ExampleDeleteAllParams, opts ...option.RequestOption) (res *ExampleDeleteAllResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, &res, opts...)
	return
}

// Count all examples by query params
func (r *ExampleService) GetCount(ctx context.Context, query ExampleGetCountParams, opts ...option.RequestOption) (res *int64, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples/count"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Upload examples from a CSV file.
//
// Note: For non-csv upload, please use the POST
// /v1/platform/datasets/{dataset_id}/examples endpoint which provides more
// efficient upload.
func (r *ExampleService) UploadFromCsv(ctx context.Context, datasetID string, body ExampleUploadFromCsvParams, opts ...option.RequestOption) (res *[]Example, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/examples/upload/%s", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type AttachmentsOperationsParam struct {
	// Mapping of old attachment names to new names
	Rename param.Field[map[string]string] `json:"rename"`
	// List of attachment names to keep
	Retain param.Field[[]string] `json:"retain"`
}

func (r AttachmentsOperationsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Example schema.
type Example struct {
	ID             string                 `json:"id,required" format:"uuid"`
	DatasetID      string                 `json:"dataset_id,required" format:"uuid"`
	Inputs         map[string]interface{} `json:"inputs,required"`
	Name           string                 `json:"name,required"`
	AttachmentURLs map[string]interface{} `json:"attachment_urls,nullable"`
	CreatedAt      time.Time              `json:"created_at" format:"date-time"`
	Metadata       map[string]interface{} `json:"metadata,nullable"`
	ModifiedAt     time.Time              `json:"modified_at,nullable" format:"date-time"`
	Outputs        map[string]interface{} `json:"outputs,nullable"`
	SourceRunID    string                 `json:"source_run_id,nullable" format:"uuid"`
	JSON           exampleJSON            `json:"-"`
}

// exampleJSON contains the JSON metadata for the struct [Example]
type exampleJSON struct {
	ID             apijson.Field
	DatasetID      apijson.Field
	Inputs         apijson.Field
	Name           apijson.Field
	AttachmentURLs apijson.Field
	CreatedAt      apijson.Field
	Metadata       apijson.Field
	ModifiedAt     apijson.Field
	Outputs        apijson.Field
	SourceRunID    apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *Example) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r exampleJSON) RawJSON() string {
	return r.raw
}

type ExampleSelect string

const (
	ExampleSelectID             ExampleSelect = "id"
	ExampleSelectCreatedAt      ExampleSelect = "created_at"
	ExampleSelectModifiedAt     ExampleSelect = "modified_at"
	ExampleSelectName           ExampleSelect = "name"
	ExampleSelectDatasetID      ExampleSelect = "dataset_id"
	ExampleSelectSourceRunID    ExampleSelect = "source_run_id"
	ExampleSelectMetadata       ExampleSelect = "metadata"
	ExampleSelectInputs         ExampleSelect = "inputs"
	ExampleSelectOutputs        ExampleSelect = "outputs"
	ExampleSelectAttachmentURLs ExampleSelect = "attachment_urls"
)

func (r ExampleSelect) IsKnown() bool {
	switch r {
	case ExampleSelectID, ExampleSelectCreatedAt, ExampleSelectModifiedAt, ExampleSelectName, ExampleSelectDatasetID, ExampleSelectSourceRunID, ExampleSelectMetadata, ExampleSelectInputs, ExampleSelectOutputs, ExampleSelectAttachmentURLs:
		return true
	}
	return false
}

type ExampleUpdateResponse = interface{}

type ExampleDeleteResponse = interface{}

type ExampleDeleteAllResponse = interface{}

type ExampleNewParams struct {
	DatasetID   param.Field[string]                     `json:"dataset_id,required" format:"uuid"`
	ID          param.Field[string]                     `json:"id" format:"uuid"`
	CreatedAt   param.Field[string]                     `json:"created_at"`
	Inputs      param.Field[map[string]interface{}]     `json:"inputs"`
	Metadata    param.Field[map[string]interface{}]     `json:"metadata"`
	Outputs     param.Field[map[string]interface{}]     `json:"outputs"`
	SourceRunID param.Field[string]                     `json:"source_run_id" format:"uuid"`
	Split       param.Field[ExampleNewParamsSplitUnion] `json:"split"`
	// Use Legacy Message Format for LLM runs
	UseLegacyMessageFormat  param.Field[bool]     `json:"use_legacy_message_format"`
	UseSourceRunAttachments param.Field[[]string] `json:"use_source_run_attachments"`
	UseSourceRunIo          param.Field[bool]     `json:"use_source_run_io"`
}

func (r ExampleNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [ExampleNewParamsSplitArray], [shared.UnionString].
type ExampleNewParamsSplitUnion interface {
	ImplementsExampleNewParamsSplitUnion()
}

type ExampleNewParamsSplitArray []string

func (r ExampleNewParamsSplitArray) ImplementsExampleNewParamsSplitUnion() {}

type ExampleGetParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[ExampleGetParamsAsOfUnion] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [ExampleGetParams]'s query parameters as `url.Values`.
func (r ExampleGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only modifications made on or before this time are included. If None, the latest
// version of the dataset is used.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type ExampleGetParamsAsOfUnion interface {
	ImplementsExampleGetParamsAsOfUnion()
}

type ExampleUpdateParams struct {
	AttachmentsOperations param.Field[AttachmentsOperationsParam]    `json:"attachments_operations"`
	DatasetID             param.Field[string]                        `json:"dataset_id" format:"uuid"`
	Inputs                param.Field[map[string]interface{}]        `json:"inputs"`
	Metadata              param.Field[map[string]interface{}]        `json:"metadata"`
	Outputs               param.Field[map[string]interface{}]        `json:"outputs"`
	Overwrite             param.Field[bool]                          `json:"overwrite"`
	Split                 param.Field[ExampleUpdateParamsSplitUnion] `json:"split"`
}

func (r ExampleUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [ExampleUpdateParamsSplitArray], [shared.UnionString].
type ExampleUpdateParamsSplitUnion interface {
	ImplementsExampleUpdateParamsSplitUnion()
}

type ExampleUpdateParamsSplitArray []string

func (r ExampleUpdateParamsSplitArray) ImplementsExampleUpdateParamsSplitUnion() {}

type ExampleListParams struct {
	ID param.Field[[]string] `query:"id" format:"uuid"`
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf             param.Field[ExampleListParamsAsOfUnion] `query:"as_of" format:"date-time"`
	Dataset          param.Field[string]                     `query:"dataset" format:"uuid"`
	Descending       param.Field[bool]                       `query:"descending"`
	Filter           param.Field[string]                     `query:"filter"`
	FullTextContains param.Field[[]string]                   `query:"full_text_contains"`
	Limit            param.Field[int64]                      `query:"limit"`
	Metadata         param.Field[string]                     `query:"metadata"`
	Offset           param.Field[int64]                      `query:"offset"`
	Order            param.Field[ExampleListParamsOrder]     `query:"order"`
	RandomSeed       param.Field[float64]                    `query:"random_seed"`
	Select           param.Field[[]ExampleSelect]            `query:"select"`
	Splits           param.Field[[]string]                   `query:"splits"`
}

// URLQuery serializes [ExampleListParams]'s query parameters as `url.Values`.
func (r ExampleListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only modifications made on or before this time are included. If None, the latest
// version of the dataset is used.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type ExampleListParamsAsOfUnion interface {
	ImplementsExampleListParamsAsOfUnion()
}

type ExampleListParamsOrder string

const (
	ExampleListParamsOrderRecent          ExampleListParamsOrder = "recent"
	ExampleListParamsOrderRandom          ExampleListParamsOrder = "random"
	ExampleListParamsOrderRecentlyCreated ExampleListParamsOrder = "recently_created"
	ExampleListParamsOrderID              ExampleListParamsOrder = "id"
)

func (r ExampleListParamsOrder) IsKnown() bool {
	switch r {
	case ExampleListParamsOrderRecent, ExampleListParamsOrderRandom, ExampleListParamsOrderRecentlyCreated, ExampleListParamsOrderID:
		return true
	}
	return false
}

type ExampleDeleteAllParams struct {
	ExampleIDs param.Field[[]string] `query:"example_ids,required" format:"uuid"`
}

// URLQuery serializes [ExampleDeleteAllParams]'s query parameters as `url.Values`.
func (r ExampleDeleteAllParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ExampleGetCountParams struct {
	ID param.Field[[]string] `query:"id" format:"uuid"`
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf             param.Field[ExampleGetCountParamsAsOfUnion] `query:"as_of" format:"date-time"`
	Dataset          param.Field[string]                         `query:"dataset" format:"uuid"`
	Filter           param.Field[string]                         `query:"filter"`
	FullTextContains param.Field[[]string]                       `query:"full_text_contains"`
	Metadata         param.Field[string]                         `query:"metadata"`
	Splits           param.Field[[]string]                       `query:"splits"`
}

// URLQuery serializes [ExampleGetCountParams]'s query parameters as `url.Values`.
func (r ExampleGetCountParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only modifications made on or before this time are included. If None, the latest
// version of the dataset is used.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type ExampleGetCountParamsAsOfUnion interface {
	ImplementsExampleGetCountParamsAsOfUnion()
}

type ExampleUploadFromCsvParams struct {
	File         param.Field[io.Reader] `json:"file,required" format:"binary"`
	InputKeys    param.Field[[]string]  `json:"input_keys,required"`
	MetadataKeys param.Field[[]string]  `json:"metadata_keys"`
	OutputKeys   param.Field[[]string]  `json:"output_keys"`
}

func (r ExampleUploadFromCsvParams) MarshalMultipart() (data []byte, contentType string, err error) {
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
