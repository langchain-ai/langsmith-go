// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetRunService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetRunService] method instead.
type DatasetRunService struct {
	Options []option.RequestOption
}

// NewDatasetRunService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetRunService(opts ...option.RequestOption) (r *DatasetRunService) {
	r = &DatasetRunService{}
	r.Options = opts
	return
}

// Fetch examples for a dataset, and fetch the runs for each example if they are
// associated with the given session_ids.
func (r *DatasetRunService) New(ctx context.Context, datasetID string, params DatasetRunNewParams, opts ...option.RequestOption) (res *[]ExampleWithRunsCh, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/runs", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

// Fetch the number of regressions/improvements for each example in a dataset,
// between sessions[0] and sessions[1].
func (r *DatasetRunService) Delta(ctx context.Context, datasetID string, body DatasetRunDeltaParams, opts ...option.RequestOption) (res *SessionFeedbackDelta, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/runs/delta", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Example schema with list of runs from ClickHouse.
//
// For non-grouped endpoint (/datasets/{dataset_id}/runs): runs from single
// session. For grouped endpoint (/datasets/{dataset_id}/group/runs): flat array of
// runs from all sessions, where each run has a session_id field for frontend to
// determine column placement.
type ExampleWithRunsCh struct {
	ID             string                 `json:"id" api:"required" format:"uuid"`
	DatasetID      string                 `json:"dataset_id" api:"required" format:"uuid"`
	Inputs         map[string]interface{} `json:"inputs" api:"required"`
	Name           string                 `json:"name" api:"required"`
	Runs           []ExampleWithRunsChRun `json:"runs" api:"required"`
	AttachmentURLs map[string]interface{} `json:"attachment_urls" api:"nullable"`
	CreatedAt      time.Time              `json:"created_at" format:"date-time"`
	Metadata       map[string]interface{} `json:"metadata" api:"nullable"`
	ModifiedAt     time.Time              `json:"modified_at" api:"nullable" format:"date-time"`
	Outputs        map[string]interface{} `json:"outputs" api:"nullable"`
	SourceRunID    string                 `json:"source_run_id" api:"nullable" format:"uuid"`
	JSON           exampleWithRunsChJSON  `json:"-"`
}

// exampleWithRunsChJSON contains the JSON metadata for the struct
// [ExampleWithRunsCh]
type exampleWithRunsChJSON struct {
	ID             apijson.Field
	DatasetID      apijson.Field
	Inputs         apijson.Field
	Name           apijson.Field
	Runs           apijson.Field
	AttachmentURLs apijson.Field
	CreatedAt      apijson.Field
	Metadata       apijson.Field
	ModifiedAt     apijson.Field
	Outputs        apijson.Field
	SourceRunID    apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *ExampleWithRunsCh) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r exampleWithRunsChJSON) RawJSON() string {
	return r.raw
}

// Run schema for comparison view.
type ExampleWithRunsChRun struct {
	ID   string `json:"id" api:"required" format:"uuid"`
	Name string `json:"name" api:"required"`
	// Enum for run types.
	RunType            RunTypeEnum                       `json:"run_type" api:"required"`
	SessionID          string                            `json:"session_id" api:"required" format:"uuid"`
	Status             string                            `json:"status" api:"required"`
	TraceID            string                            `json:"trace_id" api:"required" format:"uuid"`
	AppPath            string                            `json:"app_path" api:"nullable"`
	CompletionCost     string                            `json:"completion_cost" api:"nullable"`
	CompletionTokens   int64                             `json:"completion_tokens" api:"nullable"`
	DottedOrder        string                            `json:"dotted_order" api:"nullable"`
	EndTime            time.Time                         `json:"end_time" api:"nullable" format:"date-time"`
	Error              string                            `json:"error" api:"nullable"`
	Events             []map[string]interface{}          `json:"events" api:"nullable"`
	ExecutionOrder     int64                             `json:"execution_order"`
	Extra              map[string]interface{}            `json:"extra" api:"nullable"`
	FeedbackStats      map[string]map[string]interface{} `json:"feedback_stats" api:"nullable"`
	Feedbacks          []FeedbackSchema                  `json:"feedbacks"`
	Inputs             map[string]interface{}            `json:"inputs" api:"nullable"`
	InputsPreview      string                            `json:"inputs_preview" api:"nullable"`
	InputsS3URLs       map[string]interface{}            `json:"inputs_s3_urls" api:"nullable"`
	ManifestID         string                            `json:"manifest_id" api:"nullable" format:"uuid"`
	ManifestS3ID       string                            `json:"manifest_s3_id" api:"nullable" format:"uuid"`
	Outputs            map[string]interface{}            `json:"outputs" api:"nullable"`
	OutputsPreview     string                            `json:"outputs_preview" api:"nullable"`
	OutputsS3URLs      map[string]interface{}            `json:"outputs_s3_urls" api:"nullable"`
	ParentRunID        string                            `json:"parent_run_id" api:"nullable" format:"uuid"`
	PromptCost         string                            `json:"prompt_cost" api:"nullable"`
	PromptTokens       int64                             `json:"prompt_tokens" api:"nullable"`
	ReferenceExampleID string                            `json:"reference_example_id" api:"nullable" format:"uuid"`
	S3URLs             map[string]interface{}            `json:"s3_urls" api:"nullable"`
	Serialized         map[string]interface{}            `json:"serialized" api:"nullable"`
	StartTime          time.Time                         `json:"start_time" format:"date-time"`
	Tags               []string                          `json:"tags" api:"nullable"`
	TotalCost          string                            `json:"total_cost" api:"nullable"`
	TotalTokens        int64                             `json:"total_tokens" api:"nullable"`
	TraceMaxStartTime  time.Time                         `json:"trace_max_start_time" api:"nullable" format:"date-time"`
	TraceMinStartTime  time.Time                         `json:"trace_min_start_time" api:"nullable" format:"date-time"`
	JSON               exampleWithRunsChRunJSON          `json:"-"`
}

// exampleWithRunsChRunJSON contains the JSON metadata for the struct
// [ExampleWithRunsChRun]
type exampleWithRunsChRunJSON struct {
	ID                 apijson.Field
	Name               apijson.Field
	RunType            apijson.Field
	SessionID          apijson.Field
	Status             apijson.Field
	TraceID            apijson.Field
	AppPath            apijson.Field
	CompletionCost     apijson.Field
	CompletionTokens   apijson.Field
	DottedOrder        apijson.Field
	EndTime            apijson.Field
	Error              apijson.Field
	Events             apijson.Field
	ExecutionOrder     apijson.Field
	Extra              apijson.Field
	FeedbackStats      apijson.Field
	Feedbacks          apijson.Field
	Inputs             apijson.Field
	InputsPreview      apijson.Field
	InputsS3URLs       apijson.Field
	ManifestID         apijson.Field
	ManifestS3ID       apijson.Field
	Outputs            apijson.Field
	OutputsPreview     apijson.Field
	OutputsS3URLs      apijson.Field
	ParentRunID        apijson.Field
	PromptCost         apijson.Field
	PromptTokens       apijson.Field
	ReferenceExampleID apijson.Field
	S3URLs             apijson.Field
	Serialized         apijson.Field
	StartTime          apijson.Field
	Tags               apijson.Field
	TotalCost          apijson.Field
	TotalTokens        apijson.Field
	TraceMaxStartTime  apijson.Field
	TraceMinStartTime  apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *ExampleWithRunsChRun) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r exampleWithRunsChRunJSON) RawJSON() string {
	return r.raw
}

type QueryFeedbackDeltaParam struct {
	BaselineSessionID       param.Field[string]              `json:"baseline_session_id" api:"required" format:"uuid"`
	ComparisonSessionIDs    param.Field[[]string]            `json:"comparison_session_ids" api:"required" format:"uuid"`
	FeedbackKey             param.Field[string]              `json:"feedback_key" api:"required"`
	ComparativeExperimentID param.Field[string]              `json:"comparative_experiment_id" format:"uuid"`
	Filters                 param.Field[map[string][]string] `json:"filters"`
	Limit                   param.Field[int64]               `json:"limit"`
	Offset                  param.Field[int64]               `json:"offset"`
}

func (r QueryFeedbackDeltaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// List of feedback keys with number of improvements and regressions for each.
type SessionFeedbackDelta struct {
	FeedbackDeltas map[string]SessionFeedbackDeltaFeedbackDelta `json:"feedback_deltas" api:"required"`
	JSON           sessionFeedbackDeltaJSON                     `json:"-"`
}

// sessionFeedbackDeltaJSON contains the JSON metadata for the struct
// [SessionFeedbackDelta]
type sessionFeedbackDeltaJSON struct {
	FeedbackDeltas apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SessionFeedbackDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionFeedbackDeltaJSON) RawJSON() string {
	return r.raw
}

// Feedback key with number of improvements and regressions.
type SessionFeedbackDeltaFeedbackDelta struct {
	ImprovedExamples  []string                              `json:"improved_examples" api:"required" format:"uuid"`
	RegressedExamples []string                              `json:"regressed_examples" api:"required" format:"uuid"`
	JSON              sessionFeedbackDeltaFeedbackDeltaJSON `json:"-"`
}

// sessionFeedbackDeltaFeedbackDeltaJSON contains the JSON metadata for the struct
// [SessionFeedbackDeltaFeedbackDelta]
type sessionFeedbackDeltaFeedbackDeltaJSON struct {
	ImprovedExamples  apijson.Field
	RegressedExamples apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *SessionFeedbackDeltaFeedbackDelta) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionFeedbackDeltaFeedbackDeltaJSON) RawJSON() string {
	return r.raw
}

type SortParamsForRunsComparisonView struct {
	SortBy    param.Field[string]                                   `json:"sort_by" api:"required"`
	SortOrder param.Field[SortParamsForRunsComparisonViewSortOrder] `json:"sort_order"`
}

func (r SortParamsForRunsComparisonView) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SortParamsForRunsComparisonViewSortOrder string

const (
	SortParamsForRunsComparisonViewSortOrderAsc  SortParamsForRunsComparisonViewSortOrder = "ASC"
	SortParamsForRunsComparisonViewSortOrderDesc SortParamsForRunsComparisonViewSortOrder = "DESC"
)

func (r SortParamsForRunsComparisonViewSortOrder) IsKnown() bool {
	switch r {
	case SortParamsForRunsComparisonViewSortOrderAsc, SortParamsForRunsComparisonViewSortOrderDesc:
		return true
	}
	return false
}

type DatasetRunNewParams struct {
	SessionIDs param.Field[[]string] `json:"session_ids" api:"required" format:"uuid"`
	// Response format, e.g., 'csv'
	Format                  param.Field[DatasetRunNewParamsFormat]       `query:"format"`
	ComparativeExperimentID param.Field[string]                          `json:"comparative_experiment_id" format:"uuid"`
	ExampleIDs              param.Field[[]string]                        `json:"example_ids" format:"uuid"`
	Filters                 param.Field[map[string][]string]             `json:"filters"`
	Limit                   param.Field[int64]                           `json:"limit"`
	Offset                  param.Field[int64]                           `json:"offset"`
	Preview                 param.Field[bool]                            `json:"preview"`
	SortParams              param.Field[SortParamsForRunsComparisonView] `json:"sort_params"`
	Stream                  param.Field[bool]                            `json:"stream"`
}

func (r DatasetRunNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [DatasetRunNewParams]'s query parameters as `url.Values`.
func (r DatasetRunNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Response format, e.g., 'csv'
type DatasetRunNewParamsFormat string

const (
	DatasetRunNewParamsFormatCsv DatasetRunNewParamsFormat = "csv"
)

func (r DatasetRunNewParamsFormat) IsKnown() bool {
	switch r {
	case DatasetRunNewParamsFormatCsv:
		return true
	}
	return false
}

type DatasetRunDeltaParams struct {
	QueryFeedbackDelta QueryFeedbackDeltaParam `json:"query_feedback_delta" api:"required"`
}

func (r DatasetRunDeltaParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.QueryFeedbackDelta)
}
