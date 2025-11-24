// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/tidwall/gjson"
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
func (r *DatasetRunService) New(ctx context.Context, datasetID string, params DatasetRunNewParams, opts ...option.RequestOption) (res *DatasetRunNewResponseUnion, err error) {
	var env apijson.UnionUnmarshaler[DatasetRunNewResponseUnion]
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/runs", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &env, opts...)
	if err != nil {
		return
	}
	res = &env.Value
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

// Example schema with list of runs.
type ExampleWithRuns struct {
	ID             string               `json:"id,required" format:"uuid"`
	DatasetID      string               `json:"dataset_id,required" format:"uuid"`
	Inputs         interface{}          `json:"inputs,required"`
	Name           string               `json:"name,required"`
	Runs           []ExampleWithRunsRun `json:"runs,required"`
	AttachmentURLs interface{}          `json:"attachment_urls,nullable"`
	CreatedAt      time.Time            `json:"created_at" format:"date-time"`
	Metadata       interface{}          `json:"metadata,nullable"`
	ModifiedAt     time.Time            `json:"modified_at,nullable" format:"date-time"`
	Outputs        interface{}          `json:"outputs,nullable"`
	SourceRunID    string               `json:"source_run_id,nullable" format:"uuid"`
	JSON           exampleWithRunsJSON  `json:"-"`
}

// exampleWithRunsJSON contains the JSON metadata for the struct [ExampleWithRuns]
type exampleWithRunsJSON struct {
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

func (r *ExampleWithRuns) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r exampleWithRunsJSON) RawJSON() string {
	return r.raw
}

// Run schema.
type ExampleWithRunsRun struct {
	ID          string `json:"id,required" format:"uuid"`
	AppPath     string `json:"app_path,required"`
	DottedOrder string `json:"dotted_order,required"`
	Name        string `json:"name,required"`
	// Enum for run types.
	RunType                ExampleWithRunsRunsRunType   `json:"run_type,required"`
	SessionID              string                       `json:"session_id,required" format:"uuid"`
	Status                 string                       `json:"status,required"`
	TraceID                string                       `json:"trace_id,required" format:"uuid"`
	ChildRunIDs            []string                     `json:"child_run_ids,nullable" format:"uuid"`
	CompletionCost         string                       `json:"completion_cost,nullable"`
	CompletionCostDetails  map[string]string            `json:"completion_cost_details,nullable"`
	CompletionTokenDetails map[string]int64             `json:"completion_token_details,nullable"`
	CompletionTokens       int64                        `json:"completion_tokens"`
	DirectChildRunIDs      []string                     `json:"direct_child_run_ids,nullable" format:"uuid"`
	EndTime                time.Time                    `json:"end_time,nullable" format:"date-time"`
	Error                  string                       `json:"error,nullable"`
	Events                 []interface{}                `json:"events,nullable"`
	ExecutionOrder         int64                        `json:"execution_order"`
	Extra                  interface{}                  `json:"extra,nullable"`
	FeedbackStats          map[string]interface{}       `json:"feedback_stats,nullable"`
	FirstTokenTime         time.Time                    `json:"first_token_time,nullable" format:"date-time"`
	InDataset              bool                         `json:"in_dataset,nullable"`
	Inputs                 interface{}                  `json:"inputs,nullable"`
	InputsPreview          string                       `json:"inputs_preview,nullable"`
	InputsS3URLs           interface{}                  `json:"inputs_s3_urls,nullable"`
	LastQueuedAt           time.Time                    `json:"last_queued_at,nullable" format:"date-time"`
	ManifestID             string                       `json:"manifest_id,nullable" format:"uuid"`
	ManifestS3ID           string                       `json:"manifest_s3_id,nullable" format:"uuid"`
	Outputs                interface{}                  `json:"outputs,nullable"`
	OutputsPreview         string                       `json:"outputs_preview,nullable"`
	OutputsS3URLs          interface{}                  `json:"outputs_s3_urls,nullable"`
	ParentRunID            string                       `json:"parent_run_id,nullable" format:"uuid"`
	ParentRunIDs           []string                     `json:"parent_run_ids,nullable" format:"uuid"`
	PriceModelID           string                       `json:"price_model_id,nullable" format:"uuid"`
	PromptCost             string                       `json:"prompt_cost,nullable"`
	PromptCostDetails      map[string]string            `json:"prompt_cost_details,nullable"`
	PromptTokenDetails     map[string]int64             `json:"prompt_token_details,nullable"`
	PromptTokens           int64                        `json:"prompt_tokens"`
	ReferenceDatasetID     string                       `json:"reference_dataset_id,nullable" format:"uuid"`
	ReferenceExampleID     string                       `json:"reference_example_id,nullable" format:"uuid"`
	S3URLs                 interface{}                  `json:"s3_urls,nullable"`
	Serialized             interface{}                  `json:"serialized,nullable"`
	ShareToken             string                       `json:"share_token,nullable" format:"uuid"`
	StartTime              time.Time                    `json:"start_time" format:"date-time"`
	Tags                   []string                     `json:"tags,nullable"`
	ThreadID               string                       `json:"thread_id,nullable"`
	TotalCost              string                       `json:"total_cost,nullable"`
	TotalTokens            int64                        `json:"total_tokens"`
	TraceFirstReceivedAt   time.Time                    `json:"trace_first_received_at,nullable" format:"date-time"`
	TraceMaxStartTime      time.Time                    `json:"trace_max_start_time,nullable" format:"date-time"`
	TraceMinStartTime      time.Time                    `json:"trace_min_start_time,nullable" format:"date-time"`
	TraceTier              ExampleWithRunsRunsTraceTier `json:"trace_tier,nullable"`
	TraceUpgrade           bool                         `json:"trace_upgrade"`
	TtlSeconds             int64                        `json:"ttl_seconds,nullable"`
	JSON                   exampleWithRunsRunJSON       `json:"-"`
}

// exampleWithRunsRunJSON contains the JSON metadata for the struct
// [ExampleWithRunsRun]
type exampleWithRunsRunJSON struct {
	ID                     apijson.Field
	AppPath                apijson.Field
	DottedOrder            apijson.Field
	Name                   apijson.Field
	RunType                apijson.Field
	SessionID              apijson.Field
	Status                 apijson.Field
	TraceID                apijson.Field
	ChildRunIDs            apijson.Field
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	DirectChildRunIDs      apijson.Field
	EndTime                apijson.Field
	Error                  apijson.Field
	Events                 apijson.Field
	ExecutionOrder         apijson.Field
	Extra                  apijson.Field
	FeedbackStats          apijson.Field
	FirstTokenTime         apijson.Field
	InDataset              apijson.Field
	Inputs                 apijson.Field
	InputsPreview          apijson.Field
	InputsS3URLs           apijson.Field
	LastQueuedAt           apijson.Field
	ManifestID             apijson.Field
	ManifestS3ID           apijson.Field
	Outputs                apijson.Field
	OutputsPreview         apijson.Field
	OutputsS3URLs          apijson.Field
	ParentRunID            apijson.Field
	ParentRunIDs           apijson.Field
	PriceModelID           apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	ReferenceDatasetID     apijson.Field
	ReferenceExampleID     apijson.Field
	S3URLs                 apijson.Field
	Serialized             apijson.Field
	ShareToken             apijson.Field
	StartTime              apijson.Field
	Tags                   apijson.Field
	ThreadID               apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	TraceFirstReceivedAt   apijson.Field
	TraceMaxStartTime      apijson.Field
	TraceMinStartTime      apijson.Field
	TraceTier              apijson.Field
	TraceUpgrade           apijson.Field
	TtlSeconds             apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *ExampleWithRunsRun) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r exampleWithRunsRunJSON) RawJSON() string {
	return r.raw
}

// Enum for run types.
type ExampleWithRunsRunsRunType string

const (
	ExampleWithRunsRunsRunTypeTool      ExampleWithRunsRunsRunType = "tool"
	ExampleWithRunsRunsRunTypeChain     ExampleWithRunsRunsRunType = "chain"
	ExampleWithRunsRunsRunTypeLlm       ExampleWithRunsRunsRunType = "llm"
	ExampleWithRunsRunsRunTypeRetriever ExampleWithRunsRunsRunType = "retriever"
	ExampleWithRunsRunsRunTypeEmbedding ExampleWithRunsRunsRunType = "embedding"
	ExampleWithRunsRunsRunTypePrompt    ExampleWithRunsRunsRunType = "prompt"
	ExampleWithRunsRunsRunTypeParser    ExampleWithRunsRunsRunType = "parser"
)

func (r ExampleWithRunsRunsRunType) IsKnown() bool {
	switch r {
	case ExampleWithRunsRunsRunTypeTool, ExampleWithRunsRunsRunTypeChain, ExampleWithRunsRunsRunTypeLlm, ExampleWithRunsRunsRunTypeRetriever, ExampleWithRunsRunsRunTypeEmbedding, ExampleWithRunsRunsRunTypePrompt, ExampleWithRunsRunsRunTypeParser:
		return true
	}
	return false
}

type ExampleWithRunsRunsTraceTier string

const (
	ExampleWithRunsRunsTraceTierLonglived  ExampleWithRunsRunsTraceTier = "longlived"
	ExampleWithRunsRunsTraceTierShortlived ExampleWithRunsRunsTraceTier = "shortlived"
)

func (r ExampleWithRunsRunsTraceTier) IsKnown() bool {
	switch r {
	case ExampleWithRunsRunsTraceTierLonglived, ExampleWithRunsRunsTraceTierShortlived:
		return true
	}
	return false
}

// Example schema with list of runs from ClickHouse.
//
// For non-grouped endpoint (/datasets/{dataset_id}/runs): runs from single
// session. For grouped endpoint (/datasets/{dataset_id}/group/runs): flat array of
// runs from all sessions, where each run has a session_id field for frontend to
// determine column placement.
type ExampleWithRunsCh struct {
	ID             string                 `json:"id,required" format:"uuid"`
	DatasetID      string                 `json:"dataset_id,required" format:"uuid"`
	Inputs         interface{}            `json:"inputs,required"`
	Name           string                 `json:"name,required"`
	Runs           []ExampleWithRunsChRun `json:"runs,required"`
	AttachmentURLs interface{}            `json:"attachment_urls,nullable"`
	CreatedAt      time.Time              `json:"created_at" format:"date-time"`
	Metadata       interface{}            `json:"metadata,nullable"`
	ModifiedAt     time.Time              `json:"modified_at,nullable" format:"date-time"`
	Outputs        interface{}            `json:"outputs,nullable"`
	SourceRunID    string                 `json:"source_run_id,nullable" format:"uuid"`
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
	ID   string `json:"id,required" format:"uuid"`
	Name string `json:"name,required"`
	// Enum for run types.
	RunType            ExampleWithRunsChRunsRunType `json:"run_type,required"`
	SessionID          string                       `json:"session_id,required" format:"uuid"`
	Status             string                       `json:"status,required"`
	TraceID            string                       `json:"trace_id,required" format:"uuid"`
	AppPath            string                       `json:"app_path,nullable"`
	CompletionCost     string                       `json:"completion_cost,nullable"`
	CompletionTokens   int64                        `json:"completion_tokens,nullable"`
	DottedOrder        string                       `json:"dotted_order,nullable"`
	EndTime            time.Time                    `json:"end_time,nullable" format:"date-time"`
	Error              string                       `json:"error,nullable"`
	Events             []interface{}                `json:"events,nullable"`
	ExecutionOrder     int64                        `json:"execution_order"`
	Extra              interface{}                  `json:"extra,nullable"`
	FeedbackStats      map[string]interface{}       `json:"feedback_stats,nullable"`
	Inputs             interface{}                  `json:"inputs,nullable"`
	InputsPreview      string                       `json:"inputs_preview,nullable"`
	InputsS3URLs       interface{}                  `json:"inputs_s3_urls,nullable"`
	ManifestID         string                       `json:"manifest_id,nullable" format:"uuid"`
	ManifestS3ID       string                       `json:"manifest_s3_id,nullable" format:"uuid"`
	Outputs            interface{}                  `json:"outputs,nullable"`
	OutputsPreview     string                       `json:"outputs_preview,nullable"`
	OutputsS3URLs      interface{}                  `json:"outputs_s3_urls,nullable"`
	ParentRunID        string                       `json:"parent_run_id,nullable" format:"uuid"`
	PromptCost         string                       `json:"prompt_cost,nullable"`
	PromptTokens       int64                        `json:"prompt_tokens,nullable"`
	ReferenceExampleID string                       `json:"reference_example_id,nullable" format:"uuid"`
	S3URLs             interface{}                  `json:"s3_urls,nullable"`
	Serialized         interface{}                  `json:"serialized,nullable"`
	StartTime          time.Time                    `json:"start_time" format:"date-time"`
	Tags               []string                     `json:"tags,nullable"`
	TotalCost          string                       `json:"total_cost,nullable"`
	TotalTokens        int64                        `json:"total_tokens,nullable"`
	TraceMaxStartTime  time.Time                    `json:"trace_max_start_time,nullable" format:"date-time"`
	TraceMinStartTime  time.Time                    `json:"trace_min_start_time,nullable" format:"date-time"`
	JSON               exampleWithRunsChRunJSON     `json:"-"`
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

// Enum for run types.
type ExampleWithRunsChRunsRunType string

const (
	ExampleWithRunsChRunsRunTypeTool      ExampleWithRunsChRunsRunType = "tool"
	ExampleWithRunsChRunsRunTypeChain     ExampleWithRunsChRunsRunType = "chain"
	ExampleWithRunsChRunsRunTypeLlm       ExampleWithRunsChRunsRunType = "llm"
	ExampleWithRunsChRunsRunTypeRetriever ExampleWithRunsChRunsRunType = "retriever"
	ExampleWithRunsChRunsRunTypeEmbedding ExampleWithRunsChRunsRunType = "embedding"
	ExampleWithRunsChRunsRunTypePrompt    ExampleWithRunsChRunsRunType = "prompt"
	ExampleWithRunsChRunsRunTypeParser    ExampleWithRunsChRunsRunType = "parser"
)

func (r ExampleWithRunsChRunsRunType) IsKnown() bool {
	switch r {
	case ExampleWithRunsChRunsRunTypeTool, ExampleWithRunsChRunsRunTypeChain, ExampleWithRunsChRunsRunTypeLlm, ExampleWithRunsChRunsRunTypeRetriever, ExampleWithRunsChRunsRunTypeEmbedding, ExampleWithRunsChRunsRunTypePrompt, ExampleWithRunsChRunsRunTypeParser:
		return true
	}
	return false
}

type QueryExampleSchemaWithRunsParam struct {
	SessionIDs              param.Field[[]string]                        `json:"session_ids,required" format:"uuid"`
	ComparativeExperimentID param.Field[string]                          `json:"comparative_experiment_id" format:"uuid"`
	Filters                 param.Field[map[string][]string]             `json:"filters"`
	Limit                   param.Field[int64]                           `json:"limit"`
	Offset                  param.Field[int64]                           `json:"offset"`
	Preview                 param.Field[bool]                            `json:"preview"`
	SortParams              param.Field[SortParamsForRunsComparisonView] `json:"sort_params"`
}

func (r QueryExampleSchemaWithRunsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type QueryFeedbackDeltaParam struct {
	BaselineSessionID       param.Field[string]              `json:"baseline_session_id,required" format:"uuid"`
	ComparisonSessionIDs    param.Field[[]string]            `json:"comparison_session_ids,required" format:"uuid"`
	FeedbackKey             param.Field[string]              `json:"feedback_key,required"`
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
	FeedbackDeltas map[string]SessionFeedbackDeltaFeedbackDelta `json:"feedback_deltas,required"`
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
	ImprovedExamples  []string                              `json:"improved_examples,required" format:"uuid"`
	RegressedExamples []string                              `json:"regressed_examples,required" format:"uuid"`
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
	SortBy    param.Field[string]                                   `json:"sort_by,required"`
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

// Union satisfied by [DatasetRunNewResponseArray] or [DatasetRunNewResponseArray].
type DatasetRunNewResponseUnion interface {
	implementsDatasetRunNewResponseUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*DatasetRunNewResponseUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(DatasetRunNewResponseArray{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(DatasetRunNewResponseArray{}),
		},
	)
}

type DatasetRunNewResponseArray []ExampleWithRuns

func (r DatasetRunNewResponseArray) implementsDatasetRunNewResponseUnion() {}

type DatasetRunNewParams struct {
	QueryExampleSchemaWithRuns QueryExampleSchemaWithRunsParam `json:"query_example_schema_with_runs,required"`
	// Response format, e.g., 'csv'
	Format param.Field[string] `query:"format"`
}

func (r DatasetRunNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.QueryExampleSchemaWithRuns)
}

// URLQuery serializes [DatasetRunNewParams]'s query parameters as `url.Values`.
func (r DatasetRunNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetRunDeltaParams struct {
	QueryFeedbackDelta QueryFeedbackDeltaParam `json:"query_feedback_delta,required"`
}

func (r DatasetRunDeltaParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.QueryFeedbackDelta)
}
