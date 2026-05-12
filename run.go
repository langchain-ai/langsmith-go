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

// RunService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewRunService] method instead.
type RunService struct {
	Options []option.RequestOption
	Rules   *RunRuleService
}

// NewRunService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewRunService(opts ...option.RequestOption) (r *RunService) {
	r = &RunService{}
	r.Options = opts
	r.Rules = NewRunRuleService(opts...)
	return
}

// Queues a single run for ingestion. The request body must be a JSON-encoded run
// object that follows the Run schema.
func (r *RunService) New(ctx context.Context, body RunNewParams, opts ...option.RequestOption) (res *RunNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "runs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a specific run.
func (r *RunService) Get(ctx context.Context, runID string, query RunGetParams, opts ...option.RequestOption) (res *RunSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/runs/%s", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Updates a run identified by its ID. The body should contain only the fields to
// be changed; unknown fields are ignored.
func (r *RunService) Update(ctx context.Context, runID string, body RunUpdateParams, opts ...option.RequestOption) (res *RunUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("runs/%s", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Ingests a batch of runs in a single JSON payload. The payload must have `post`
// and/or `patch` arrays containing run objects. Prefer this endpoint over
// single‑run ingestion when submitting hundreds of runs, but `/runs/multipart`
// offers better handling for very large fields and attachments.
func (r *RunService) IngestBatch(ctx context.Context, body RunIngestBatchParams, opts ...option.RequestOption) (res *RunIngestBatchResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "runs/batch"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Query Runs
func (r *RunService) Query(ctx context.Context, body RunQueryParams, opts ...option.RequestOption) (res *RunQueryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/runs/query"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get all runs by query in body payload.
func (r *RunService) Stats(ctx context.Context, body RunStatsParams, opts ...option.RequestOption) (res *RunStatsResponseUnion, err error) {
	var env apijson.UnionUnmarshaler[RunStatsResponseUnion]
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/runs/stats"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &env, opts...)
	if err != nil {
		return nil, err
	}
	res = &env.Value
	return res, nil
}

// Update a run.
func (r *RunService) Update2(ctx context.Context, runID string, opts ...option.RequestOption) (res *RunUpdate2Response, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/runs/%s", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, nil, &res, opts...)
	return res, err
}

type RunParam struct {
	ID                 param.Field[string]                   `json:"id"`
	DottedOrder        param.Field[string]                   `json:"dotted_order"`
	EndTime            param.Field[string]                   `json:"end_time"`
	Error              param.Field[string]                   `json:"error"`
	Events             param.Field[[]map[string]interface{}] `json:"events"`
	Extra              param.Field[map[string]interface{}]   `json:"extra"`
	InputAttachments   param.Field[map[string]interface{}]   `json:"input_attachments"`
	Inputs             param.Field[map[string]interface{}]   `json:"inputs"`
	Name               param.Field[string]                   `json:"name"`
	OutputAttachments  param.Field[map[string]interface{}]   `json:"output_attachments"`
	Outputs            param.Field[map[string]interface{}]   `json:"outputs"`
	ParentRunID        param.Field[string]                   `json:"parent_run_id"`
	ReferenceExampleID param.Field[string]                   `json:"reference_example_id"`
	RunType            param.Field[RunRunType]               `json:"run_type"`
	Serialized         param.Field[map[string]interface{}]   `json:"serialized"`
	SessionID          param.Field[string]                   `json:"session_id"`
	SessionName        param.Field[string]                   `json:"session_name"`
	StartTime          param.Field[string]                   `json:"start_time"`
	Status             param.Field[string]                   `json:"status"`
	Tags               param.Field[[]string]                 `json:"tags"`
	TraceID            param.Field[string]                   `json:"trace_id"`
}

func (r RunParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RunRunType string

const (
	RunRunTypeTool      RunRunType = "tool"
	RunRunTypeChain     RunRunType = "chain"
	RunRunTypeLlm       RunRunType = "llm"
	RunRunTypeRetriever RunRunType = "retriever"
	RunRunTypeEmbedding RunRunType = "embedding"
	RunRunTypePrompt    RunRunType = "prompt"
	RunRunTypeParser    RunRunType = "parser"
)

func (r RunRunType) IsKnown() bool {
	switch r {
	case RunRunTypeTool, RunRunTypeChain, RunRunTypeLlm, RunRunTypeRetriever, RunRunTypeEmbedding, RunRunTypePrompt, RunRunTypeParser:
		return true
	}
	return false
}

// Run schema.
type RunSchema struct {
	ID          string `json:"id" api:"required" format:"uuid"`
	AppPath     string `json:"app_path" api:"required"`
	DottedOrder string `json:"dotted_order" api:"required"`
	Name        string `json:"name" api:"required"`
	// Enum for run types.
	RunType                RunTypeEnum                       `json:"run_type" api:"required"`
	SessionID              string                            `json:"session_id" api:"required" format:"uuid"`
	Status                 string                            `json:"status" api:"required"`
	TraceID                string                            `json:"trace_id" api:"required" format:"uuid"`
	ChildRunIDs            []string                          `json:"child_run_ids" api:"nullable" format:"uuid"`
	CompletionCost         float64                           `json:"completion_cost" api:"nullable"`
	CompletionCostDetails  map[string]string                 `json:"completion_cost_details" api:"nullable"`
	CompletionTokenDetails map[string]int64                  `json:"completion_token_details" api:"nullable"`
	CompletionTokens       int64                             `json:"completion_tokens"`
	DirectChildRunIDs      []string                          `json:"direct_child_run_ids" api:"nullable" format:"uuid"`
	EndTime                time.Time                         `json:"end_time" api:"nullable" format:"date-time"`
	Error                  string                            `json:"error" api:"nullable"`
	Events                 []map[string]interface{}          `json:"events" api:"nullable"`
	ExecutionOrder         int64                             `json:"execution_order"`
	Extra                  map[string]interface{}            `json:"extra" api:"nullable"`
	FeedbackStats          map[string]map[string]interface{} `json:"feedback_stats" api:"nullable"`
	FirstTokenTime         time.Time                         `json:"first_token_time" api:"nullable" format:"date-time"`
	InDataset              bool                              `json:"in_dataset" api:"nullable"`
	Inputs                 map[string]interface{}            `json:"inputs" api:"nullable"`
	InputsPreview          string                            `json:"inputs_preview" api:"nullable"`
	InputsS3URLs           map[string]interface{}            `json:"inputs_s3_urls" api:"nullable"`
	LastQueuedAt           time.Time                         `json:"last_queued_at" api:"nullable" format:"date-time"`
	ManifestID             string                            `json:"manifest_id" api:"nullable" format:"uuid"`
	ManifestS3ID           string                            `json:"manifest_s3_id" api:"nullable" format:"uuid"`
	Messages               []map[string]interface{}          `json:"messages" api:"nullable"`
	Outputs                map[string]interface{}            `json:"outputs" api:"nullable"`
	OutputsPreview         string                            `json:"outputs_preview" api:"nullable"`
	OutputsS3URLs          map[string]interface{}            `json:"outputs_s3_urls" api:"nullable"`
	ParentRunID            string                            `json:"parent_run_id" api:"nullable" format:"uuid"`
	ParentRunIDs           []string                          `json:"parent_run_ids" api:"nullable" format:"uuid"`
	PriceModelID           string                            `json:"price_model_id" api:"nullable" format:"uuid"`
	PromptCost             float64                           `json:"prompt_cost" api:"nullable"`
	PromptCostDetails      map[string]string                 `json:"prompt_cost_details" api:"nullable"`
	PromptTokenDetails     map[string]int64                  `json:"prompt_token_details" api:"nullable"`
	PromptTokens           int64                             `json:"prompt_tokens"`
	ReferenceDatasetID     string                            `json:"reference_dataset_id" api:"nullable" format:"uuid"`
	ReferenceExampleID     string                            `json:"reference_example_id" api:"nullable" format:"uuid"`
	S3URLs                 map[string]interface{}            `json:"s3_urls" api:"nullable"`
	Serialized             map[string]interface{}            `json:"serialized" api:"nullable"`
	ShareToken             string                            `json:"share_token" api:"nullable" format:"uuid"`
	StartTime              time.Time                         `json:"start_time" format:"date-time"`
	Tags                   []string                          `json:"tags" api:"nullable"`
	ThreadID               string                            `json:"thread_id" api:"nullable"`
	TotalCost              float64                           `json:"total_cost" api:"nullable"`
	TotalTokens            int64                             `json:"total_tokens"`
	TraceFirstReceivedAt   time.Time                         `json:"trace_first_received_at" api:"nullable" format:"date-time"`
	TraceMaxStartTime      time.Time                         `json:"trace_max_start_time" api:"nullable" format:"date-time"`
	TraceMinStartTime      time.Time                         `json:"trace_min_start_time" api:"nullable" format:"date-time"`
	TraceTier              RunSchemaTraceTier                `json:"trace_tier" api:"nullable"`
	TraceUpgrade           bool                              `json:"trace_upgrade"`
	TtlSeconds             int64                             `json:"ttl_seconds" api:"nullable"`
	JSON                   runSchemaJSON                     `json:"-"`
}

// runSchemaJSON contains the JSON metadata for the struct [RunSchema]
type runSchemaJSON struct {
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
	Messages               apijson.Field
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

func (r *RunSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runSchemaJSON) RawJSON() string {
	return r.raw
}

type RunSchemaTraceTier string

const (
	RunSchemaTraceTierLonglived  RunSchemaTraceTier = "longlived"
	RunSchemaTraceTierShortlived RunSchemaTraceTier = "shortlived"
)

func (r RunSchemaTraceTier) IsKnown() bool {
	switch r {
	case RunSchemaTraceTierLonglived, RunSchemaTraceTierShortlived:
		return true
	}
	return false
}

// Query params for run stats.
type RunStatsQueryParams struct {
	ID param.Field[[]string] `json:"id" format:"uuid"`
	// Enum for run data source types.
	DataSourceType param.Field[RunsFilterDataSourceTypeEnum] `json:"data_source_type"`
	EndTime        param.Field[time.Time]                    `json:"end_time" format:"date-time"`
	Error          param.Field[bool]                         `json:"error"`
	ExecutionOrder param.Field[int64]                        `json:"execution_order"`
	Filter         param.Field[string]                       `json:"filter"`
	// Group by param for run stats.
	GroupBy          param.Field[RunStatsGroupByParam] `json:"group_by"`
	Groups           param.Field[[]string]             `json:"groups"`
	IsRoot           param.Field[bool]                 `json:"is_root"`
	ParentRun        param.Field[string]               `json:"parent_run" format:"uuid"`
	Query            param.Field[string]               `json:"query"`
	ReferenceExample param.Field[[]string]             `json:"reference_example" format:"uuid"`
	// Enum for run types.
	RunType               param.Field[RunTypeEnum]                 `json:"run_type"`
	SearchFilter          param.Field[string]                      `json:"search_filter"`
	Select                param.Field[[]RunStatsQueryParamsSelect] `json:"select"`
	Session               param.Field[[]string]                    `json:"session" format:"uuid"`
	SkipPagination        param.Field[bool]                        `json:"skip_pagination"`
	StartTime             param.Field[time.Time]                   `json:"start_time" format:"date-time"`
	Trace                 param.Field[string]                      `json:"trace" format:"uuid"`
	TraceFilter           param.Field[string]                      `json:"trace_filter"`
	TreeFilter            param.Field[string]                      `json:"tree_filter"`
	UseExperimentalSearch param.Field[bool]                        `json:"use_experimental_search"`
}

func (r RunStatsQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Metrics you can select from run stats endpoint.
type RunStatsQueryParamsSelect string

const (
	RunStatsQueryParamsSelectRunCount               RunStatsQueryParamsSelect = "run_count"
	RunStatsQueryParamsSelectLatencyP50             RunStatsQueryParamsSelect = "latency_p50"
	RunStatsQueryParamsSelectLatencyP99             RunStatsQueryParamsSelect = "latency_p99"
	RunStatsQueryParamsSelectLatencyAvg             RunStatsQueryParamsSelect = "latency_avg"
	RunStatsQueryParamsSelectFirstTokenP50          RunStatsQueryParamsSelect = "first_token_p50"
	RunStatsQueryParamsSelectFirstTokenP99          RunStatsQueryParamsSelect = "first_token_p99"
	RunStatsQueryParamsSelectTotalTokens            RunStatsQueryParamsSelect = "total_tokens"
	RunStatsQueryParamsSelectPromptTokens           RunStatsQueryParamsSelect = "prompt_tokens"
	RunStatsQueryParamsSelectCompletionTokens       RunStatsQueryParamsSelect = "completion_tokens"
	RunStatsQueryParamsSelectMedianTokens           RunStatsQueryParamsSelect = "median_tokens"
	RunStatsQueryParamsSelectCompletionTokensP50    RunStatsQueryParamsSelect = "completion_tokens_p50"
	RunStatsQueryParamsSelectPromptTokensP50        RunStatsQueryParamsSelect = "prompt_tokens_p50"
	RunStatsQueryParamsSelectTokensP99              RunStatsQueryParamsSelect = "tokens_p99"
	RunStatsQueryParamsSelectCompletionTokensP99    RunStatsQueryParamsSelect = "completion_tokens_p99"
	RunStatsQueryParamsSelectPromptTokensP99        RunStatsQueryParamsSelect = "prompt_tokens_p99"
	RunStatsQueryParamsSelectLastRunStartTime       RunStatsQueryParamsSelect = "last_run_start_time"
	RunStatsQueryParamsSelectFeedbackStats          RunStatsQueryParamsSelect = "feedback_stats"
	RunStatsQueryParamsSelectThreadFeedbackStats    RunStatsQueryParamsSelect = "thread_feedback_stats"
	RunStatsQueryParamsSelectRunFacets              RunStatsQueryParamsSelect = "run_facets"
	RunStatsQueryParamsSelectErrorRate              RunStatsQueryParamsSelect = "error_rate"
	RunStatsQueryParamsSelectStreamingRate          RunStatsQueryParamsSelect = "streaming_rate"
	RunStatsQueryParamsSelectTotalCost              RunStatsQueryParamsSelect = "total_cost"
	RunStatsQueryParamsSelectPromptCost             RunStatsQueryParamsSelect = "prompt_cost"
	RunStatsQueryParamsSelectCompletionCost         RunStatsQueryParamsSelect = "completion_cost"
	RunStatsQueryParamsSelectCostP50                RunStatsQueryParamsSelect = "cost_p50"
	RunStatsQueryParamsSelectCostP99                RunStatsQueryParamsSelect = "cost_p99"
	RunStatsQueryParamsSelectSessionFeedbackStats   RunStatsQueryParamsSelect = "session_feedback_stats"
	RunStatsQueryParamsSelectAllRunStats            RunStatsQueryParamsSelect = "all_run_stats"
	RunStatsQueryParamsSelectAllTokenStats          RunStatsQueryParamsSelect = "all_token_stats"
	RunStatsQueryParamsSelectPromptTokenDetails     RunStatsQueryParamsSelect = "prompt_token_details"
	RunStatsQueryParamsSelectCompletionTokenDetails RunStatsQueryParamsSelect = "completion_token_details"
	RunStatsQueryParamsSelectPromptCostDetails      RunStatsQueryParamsSelect = "prompt_cost_details"
	RunStatsQueryParamsSelectCompletionCostDetails  RunStatsQueryParamsSelect = "completion_cost_details"
)

func (r RunStatsQueryParamsSelect) IsKnown() bool {
	switch r {
	case RunStatsQueryParamsSelectRunCount, RunStatsQueryParamsSelectLatencyP50, RunStatsQueryParamsSelectLatencyP99, RunStatsQueryParamsSelectLatencyAvg, RunStatsQueryParamsSelectFirstTokenP50, RunStatsQueryParamsSelectFirstTokenP99, RunStatsQueryParamsSelectTotalTokens, RunStatsQueryParamsSelectPromptTokens, RunStatsQueryParamsSelectCompletionTokens, RunStatsQueryParamsSelectMedianTokens, RunStatsQueryParamsSelectCompletionTokensP50, RunStatsQueryParamsSelectPromptTokensP50, RunStatsQueryParamsSelectTokensP99, RunStatsQueryParamsSelectCompletionTokensP99, RunStatsQueryParamsSelectPromptTokensP99, RunStatsQueryParamsSelectLastRunStartTime, RunStatsQueryParamsSelectFeedbackStats, RunStatsQueryParamsSelectThreadFeedbackStats, RunStatsQueryParamsSelectRunFacets, RunStatsQueryParamsSelectErrorRate, RunStatsQueryParamsSelectStreamingRate, RunStatsQueryParamsSelectTotalCost, RunStatsQueryParamsSelectPromptCost, RunStatsQueryParamsSelectCompletionCost, RunStatsQueryParamsSelectCostP50, RunStatsQueryParamsSelectCostP99, RunStatsQueryParamsSelectSessionFeedbackStats, RunStatsQueryParamsSelectAllRunStats, RunStatsQueryParamsSelectAllTokenStats, RunStatsQueryParamsSelectPromptTokenDetails, RunStatsQueryParamsSelectCompletionTokenDetails, RunStatsQueryParamsSelectPromptCostDetails, RunStatsQueryParamsSelectCompletionCostDetails:
		return true
	}
	return false
}

// Enum for run types.
type RunTypeEnum string

const (
	RunTypeEnumTool      RunTypeEnum = "tool"
	RunTypeEnumChain     RunTypeEnum = "chain"
	RunTypeEnumLlm       RunTypeEnum = "llm"
	RunTypeEnumRetriever RunTypeEnum = "retriever"
	RunTypeEnumEmbedding RunTypeEnum = "embedding"
	RunTypeEnumPrompt    RunTypeEnum = "prompt"
	RunTypeEnumParser    RunTypeEnum = "parser"
)

func (r RunTypeEnum) IsKnown() bool {
	switch r {
	case RunTypeEnumTool, RunTypeEnumChain, RunTypeEnumLlm, RunTypeEnumRetriever, RunTypeEnumEmbedding, RunTypeEnumPrompt, RunTypeEnumParser:
		return true
	}
	return false
}

// Enum for run data source types.
type RunsFilterDataSourceTypeEnum string

const (
	RunsFilterDataSourceTypeEnumCurrent              RunsFilterDataSourceTypeEnum = "current"
	RunsFilterDataSourceTypeEnumHistorical           RunsFilterDataSourceTypeEnum = "historical"
	RunsFilterDataSourceTypeEnumLite                 RunsFilterDataSourceTypeEnum = "lite"
	RunsFilterDataSourceTypeEnumRootLite             RunsFilterDataSourceTypeEnum = "root_lite"
	RunsFilterDataSourceTypeEnumRunsFeedbacksRmtWide RunsFilterDataSourceTypeEnum = "runs_feedbacks_rmt_wide"
)

func (r RunsFilterDataSourceTypeEnum) IsKnown() bool {
	switch r {
	case RunsFilterDataSourceTypeEnumCurrent, RunsFilterDataSourceTypeEnumHistorical, RunsFilterDataSourceTypeEnumLite, RunsFilterDataSourceTypeEnumRootLite, RunsFilterDataSourceTypeEnumRunsFeedbacksRmtWide:
		return true
	}
	return false
}

type RunNewResponse map[string]RunNewResponseItem

type RunNewResponseItem struct {
	JSON runNewResponseItemJSON `json:"-"`
}

// runNewResponseItemJSON contains the JSON metadata for the struct
// [RunNewResponseItem]
type runNewResponseItemJSON struct {
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunNewResponseItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runNewResponseItemJSON) RawJSON() string {
	return r.raw
}

type RunUpdateResponse map[string]RunUpdateResponseItem

type RunUpdateResponseItem struct {
	JSON runUpdateResponseItemJSON `json:"-"`
}

// runUpdateResponseItemJSON contains the JSON metadata for the struct
// [RunUpdateResponseItem]
type runUpdateResponseItemJSON struct {
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunUpdateResponseItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runUpdateResponseItemJSON) RawJSON() string {
	return r.raw
}

type RunIngestBatchResponse map[string]RunIngestBatchResponseItem

type RunIngestBatchResponseItem struct {
	JSON runIngestBatchResponseItemJSON `json:"-"`
}

// runIngestBatchResponseItemJSON contains the JSON metadata for the struct
// [RunIngestBatchResponseItem]
type runIngestBatchResponseItemJSON struct {
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunIngestBatchResponseItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runIngestBatchResponseItemJSON) RawJSON() string {
	return r.raw
}

type RunQueryResponse struct {
	Cursors       map[string]string      `json:"cursors" api:"required"`
	Runs          []RunSchema            `json:"runs" api:"required"`
	ParsedQuery   string                 `json:"parsed_query" api:"nullable"`
	SearchCursors map[string]interface{} `json:"search_cursors" api:"nullable"`
	JSON          runQueryResponseJSON   `json:"-"`
}

// runQueryResponseJSON contains the JSON metadata for the struct
// [RunQueryResponse]
type runQueryResponseJSON struct {
	Cursors       apijson.Field
	Runs          apijson.Field
	ParsedQuery   apijson.Field
	SearchCursors apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *RunQueryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runQueryResponseJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [RunStatsResponseRunStats] or [RunStatsResponseMap].
type RunStatsResponseUnion interface {
	implementsRunStatsResponseUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*RunStatsResponseUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(RunStatsResponseRunStats{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(RunStatsResponseMap{}),
		},
	)
}

type RunStatsResponseRunStats struct {
	CompletionCost         float64                      `json:"completion_cost" api:"nullable"`
	CompletionCostDetails  map[string]interface{}       `json:"completion_cost_details" api:"nullable"`
	CompletionTokenDetails map[string]interface{}       `json:"completion_token_details" api:"nullable"`
	CompletionTokens       int64                        `json:"completion_tokens" api:"nullable"`
	CompletionTokensP50    int64                        `json:"completion_tokens_p50" api:"nullable"`
	CompletionTokensP99    int64                        `json:"completion_tokens_p99" api:"nullable"`
	CostP50                float64                      `json:"cost_p50" api:"nullable"`
	CostP99                float64                      `json:"cost_p99" api:"nullable"`
	ErrorRate              float64                      `json:"error_rate" api:"nullable"`
	FeedbackStats          map[string]interface{}       `json:"feedback_stats" api:"nullable"`
	FirstTokenP50          float64                      `json:"first_token_p50" api:"nullable"`
	FirstTokenP99          float64                      `json:"first_token_p99" api:"nullable"`
	LastRunStartTime       time.Time                    `json:"last_run_start_time" api:"nullable" format:"date-time"`
	LatencyP50             float64                      `json:"latency_p50" api:"nullable"`
	LatencyP99             float64                      `json:"latency_p99" api:"nullable"`
	MedianTokens           int64                        `json:"median_tokens" api:"nullable"`
	PromptCost             float64                      `json:"prompt_cost" api:"nullable"`
	PromptCostDetails      map[string]interface{}       `json:"prompt_cost_details" api:"nullable"`
	PromptTokenDetails     map[string]interface{}       `json:"prompt_token_details" api:"nullable"`
	PromptTokens           int64                        `json:"prompt_tokens" api:"nullable"`
	PromptTokensP50        int64                        `json:"prompt_tokens_p50" api:"nullable"`
	PromptTokensP99        int64                        `json:"prompt_tokens_p99" api:"nullable"`
	RunCount               int64                        `json:"run_count" api:"nullable"`
	RunFacets              []map[string]interface{}     `json:"run_facets" api:"nullable"`
	StreamingRate          float64                      `json:"streaming_rate" api:"nullable"`
	TokensP99              int64                        `json:"tokens_p99" api:"nullable"`
	TotalCost              float64                      `json:"total_cost" api:"nullable"`
	TotalTokens            int64                        `json:"total_tokens" api:"nullable"`
	JSON                   runStatsResponseRunStatsJSON `json:"-"`
}

// runStatsResponseRunStatsJSON contains the JSON metadata for the struct
// [RunStatsResponseRunStats]
type runStatsResponseRunStatsJSON struct {
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	CompletionTokensP50    apijson.Field
	CompletionTokensP99    apijson.Field
	CostP50                apijson.Field
	CostP99                apijson.Field
	ErrorRate              apijson.Field
	FeedbackStats          apijson.Field
	FirstTokenP50          apijson.Field
	FirstTokenP99          apijson.Field
	LastRunStartTime       apijson.Field
	LatencyP50             apijson.Field
	LatencyP99             apijson.Field
	MedianTokens           apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	PromptTokensP50        apijson.Field
	PromptTokensP99        apijson.Field
	RunCount               apijson.Field
	RunFacets              apijson.Field
	StreamingRate          apijson.Field
	TokensP99              apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *RunStatsResponseRunStats) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStatsResponseRunStatsJSON) RawJSON() string {
	return r.raw
}

func (r RunStatsResponseRunStats) implementsRunStatsResponseUnion() {}

type RunStatsResponseMap map[string]RunStatsResponseMapItem

func (r RunStatsResponseMap) implementsRunStatsResponseUnion() {}

type RunStatsResponseMapItem struct {
	CompletionCost         float64                     `json:"completion_cost" api:"nullable"`
	CompletionCostDetails  map[string]interface{}      `json:"completion_cost_details" api:"nullable"`
	CompletionTokenDetails map[string]interface{}      `json:"completion_token_details" api:"nullable"`
	CompletionTokens       int64                       `json:"completion_tokens" api:"nullable"`
	CompletionTokensP50    int64                       `json:"completion_tokens_p50" api:"nullable"`
	CompletionTokensP99    int64                       `json:"completion_tokens_p99" api:"nullable"`
	CostP50                float64                     `json:"cost_p50" api:"nullable"`
	CostP99                float64                     `json:"cost_p99" api:"nullable"`
	ErrorRate              float64                     `json:"error_rate" api:"nullable"`
	FeedbackStats          map[string]interface{}      `json:"feedback_stats" api:"nullable"`
	FirstTokenP50          float64                     `json:"first_token_p50" api:"nullable"`
	FirstTokenP99          float64                     `json:"first_token_p99" api:"nullable"`
	LastRunStartTime       time.Time                   `json:"last_run_start_time" api:"nullable" format:"date-time"`
	LatencyP50             float64                     `json:"latency_p50" api:"nullable"`
	LatencyP99             float64                     `json:"latency_p99" api:"nullable"`
	MedianTokens           int64                       `json:"median_tokens" api:"nullable"`
	PromptCost             float64                     `json:"prompt_cost" api:"nullable"`
	PromptCostDetails      map[string]interface{}      `json:"prompt_cost_details" api:"nullable"`
	PromptTokenDetails     map[string]interface{}      `json:"prompt_token_details" api:"nullable"`
	PromptTokens           int64                       `json:"prompt_tokens" api:"nullable"`
	PromptTokensP50        int64                       `json:"prompt_tokens_p50" api:"nullable"`
	PromptTokensP99        int64                       `json:"prompt_tokens_p99" api:"nullable"`
	RunCount               int64                       `json:"run_count" api:"nullable"`
	RunFacets              []map[string]interface{}    `json:"run_facets" api:"nullable"`
	StreamingRate          float64                     `json:"streaming_rate" api:"nullable"`
	TokensP99              int64                       `json:"tokens_p99" api:"nullable"`
	TotalCost              float64                     `json:"total_cost" api:"nullable"`
	TotalTokens            int64                       `json:"total_tokens" api:"nullable"`
	JSON                   runStatsResponseMapItemJSON `json:"-"`
}

// runStatsResponseMapItemJSON contains the JSON metadata for the struct
// [RunStatsResponseMapItem]
type runStatsResponseMapItemJSON struct {
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	CompletionTokensP50    apijson.Field
	CompletionTokensP99    apijson.Field
	CostP50                apijson.Field
	CostP99                apijson.Field
	ErrorRate              apijson.Field
	FeedbackStats          apijson.Field
	FirstTokenP50          apijson.Field
	FirstTokenP99          apijson.Field
	LastRunStartTime       apijson.Field
	LatencyP50             apijson.Field
	LatencyP99             apijson.Field
	MedianTokens           apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	PromptTokensP50        apijson.Field
	PromptTokensP99        apijson.Field
	RunCount               apijson.Field
	RunFacets              apijson.Field
	StreamingRate          apijson.Field
	TokensP99              apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *RunStatsResponseMapItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStatsResponseMapItemJSON) RawJSON() string {
	return r.raw
}

type RunUpdate2Response = interface{}

type RunNewParams struct {
	Run RunParam `json:"run" api:"required"`
}

func (r RunNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.Run)
}

type RunGetParams struct {
	ExcludeS3StoredAttributes param.Field[bool]      `query:"exclude_s3_stored_attributes"`
	ExcludeSerialized         param.Field[bool]      `query:"exclude_serialized"`
	IncludeMessages           param.Field[bool]      `query:"include_messages"`
	SessionID                 param.Field[string]    `query:"session_id" format:"uuid"`
	StartTime                 param.Field[time.Time] `query:"start_time" format:"date-time"`
}

// URLQuery serializes [RunGetParams]'s query parameters as `url.Values`.
func (r RunGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type RunUpdateParams struct {
	Run RunParam `json:"run" api:"required"`
}

func (r RunUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.Run)
}

type RunIngestBatchParams struct {
	Patch param.Field[[]RunParam] `json:"patch"`
	Post  param.Field[[]RunParam] `json:"post"`
}

func (r RunIngestBatchParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RunQueryParams struct {
	ID     param.Field[[]string] `json:"id" format:"uuid"`
	Cursor param.Field[string]   `json:"cursor"`
	// Enum for run data source types.
	DataSourceType param.Field[RunsFilterDataSourceTypeEnum] `json:"data_source_type"`
	EndTime        param.Field[time.Time]                    `json:"end_time" format:"date-time"`
	Error          param.Field[bool]                         `json:"error"`
	ExecutionOrder param.Field[int64]                        `json:"execution_order"`
	Filter         param.Field[string]                       `json:"filter"`
	IsRoot         param.Field[bool]                         `json:"is_root"`
	Limit          param.Field[int64]                        `json:"limit"`
	// Enum for run start date order.
	Order            param.Field[RunQueryParamsOrder] `json:"order"`
	ParentRun        param.Field[string]              `json:"parent_run" format:"uuid"`
	Query            param.Field[string]              `json:"query"`
	ReferenceExample param.Field[[]string]            `json:"reference_example" format:"uuid"`
	// Enum for run types.
	RunType               param.Field[RunTypeEnum]            `json:"run_type"`
	SearchFilter          param.Field[string]                 `json:"search_filter"`
	Select                param.Field[[]RunQueryParamsSelect] `json:"select"`
	Session               param.Field[[]string]               `json:"session" format:"uuid"`
	SkipPagination        param.Field[bool]                   `json:"skip_pagination"`
	SkipPrevCursor        param.Field[bool]                   `json:"skip_prev_cursor"`
	StartTime             param.Field[time.Time]              `json:"start_time" format:"date-time"`
	Trace                 param.Field[string]                 `json:"trace" format:"uuid"`
	TraceFilter           param.Field[string]                 `json:"trace_filter"`
	TreeFilter            param.Field[string]                 `json:"tree_filter"`
	UseExperimentalSearch param.Field[bool]                   `json:"use_experimental_search"`
}

func (r RunQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Enum for run start date order.
type RunQueryParamsOrder string

const (
	RunQueryParamsOrderAsc  RunQueryParamsOrder = "asc"
	RunQueryParamsOrderDesc RunQueryParamsOrder = "desc"
)

func (r RunQueryParamsOrder) IsKnown() bool {
	switch r {
	case RunQueryParamsOrderAsc, RunQueryParamsOrderDesc:
		return true
	}
	return false
}

// Enum for available run columns.
type RunQueryParamsSelect string

const (
	RunQueryParamsSelectID                     RunQueryParamsSelect = "id"
	RunQueryParamsSelectName                   RunQueryParamsSelect = "name"
	RunQueryParamsSelectRunType                RunQueryParamsSelect = "run_type"
	RunQueryParamsSelectStartTime              RunQueryParamsSelect = "start_time"
	RunQueryParamsSelectEndTime                RunQueryParamsSelect = "end_time"
	RunQueryParamsSelectStatus                 RunQueryParamsSelect = "status"
	RunQueryParamsSelectError                  RunQueryParamsSelect = "error"
	RunQueryParamsSelectExtra                  RunQueryParamsSelect = "extra"
	RunQueryParamsSelectEvents                 RunQueryParamsSelect = "events"
	RunQueryParamsSelectInputs                 RunQueryParamsSelect = "inputs"
	RunQueryParamsSelectInputsPreview          RunQueryParamsSelect = "inputs_preview"
	RunQueryParamsSelectInputsS3URLs           RunQueryParamsSelect = "inputs_s3_urls"
	RunQueryParamsSelectInputsOrSignedURL      RunQueryParamsSelect = "inputs_or_signed_url"
	RunQueryParamsSelectOutputs                RunQueryParamsSelect = "outputs"
	RunQueryParamsSelectOutputsPreview         RunQueryParamsSelect = "outputs_preview"
	RunQueryParamsSelectOutputsS3URLs          RunQueryParamsSelect = "outputs_s3_urls"
	RunQueryParamsSelectOutputsOrSignedURL     RunQueryParamsSelect = "outputs_or_signed_url"
	RunQueryParamsSelectS3URLs                 RunQueryParamsSelect = "s3_urls"
	RunQueryParamsSelectErrorOrSignedURL       RunQueryParamsSelect = "error_or_signed_url"
	RunQueryParamsSelectEventsOrSignedURL      RunQueryParamsSelect = "events_or_signed_url"
	RunQueryParamsSelectExtraOrSignedURL       RunQueryParamsSelect = "extra_or_signed_url"
	RunQueryParamsSelectSerializedOrSignedURL  RunQueryParamsSelect = "serialized_or_signed_url"
	RunQueryParamsSelectParentRunID            RunQueryParamsSelect = "parent_run_id"
	RunQueryParamsSelectManifestID             RunQueryParamsSelect = "manifest_id"
	RunQueryParamsSelectManifestS3ID           RunQueryParamsSelect = "manifest_s3_id"
	RunQueryParamsSelectManifest               RunQueryParamsSelect = "manifest"
	RunQueryParamsSelectSessionID              RunQueryParamsSelect = "session_id"
	RunQueryParamsSelectSerialized             RunQueryParamsSelect = "serialized"
	RunQueryParamsSelectReferenceExampleID     RunQueryParamsSelect = "reference_example_id"
	RunQueryParamsSelectReferenceDatasetID     RunQueryParamsSelect = "reference_dataset_id"
	RunQueryParamsSelectTotalTokens            RunQueryParamsSelect = "total_tokens"
	RunQueryParamsSelectPromptTokens           RunQueryParamsSelect = "prompt_tokens"
	RunQueryParamsSelectPromptTokenDetails     RunQueryParamsSelect = "prompt_token_details"
	RunQueryParamsSelectCompletionTokens       RunQueryParamsSelect = "completion_tokens"
	RunQueryParamsSelectCompletionTokenDetails RunQueryParamsSelect = "completion_token_details"
	RunQueryParamsSelectTotalCost              RunQueryParamsSelect = "total_cost"
	RunQueryParamsSelectPromptCost             RunQueryParamsSelect = "prompt_cost"
	RunQueryParamsSelectPromptCostDetails      RunQueryParamsSelect = "prompt_cost_details"
	RunQueryParamsSelectCompletionCost         RunQueryParamsSelect = "completion_cost"
	RunQueryParamsSelectCompletionCostDetails  RunQueryParamsSelect = "completion_cost_details"
	RunQueryParamsSelectPriceModelID           RunQueryParamsSelect = "price_model_id"
	RunQueryParamsSelectFirstTokenTime         RunQueryParamsSelect = "first_token_time"
	RunQueryParamsSelectTraceID                RunQueryParamsSelect = "trace_id"
	RunQueryParamsSelectDottedOrder            RunQueryParamsSelect = "dotted_order"
	RunQueryParamsSelectLastQueuedAt           RunQueryParamsSelect = "last_queued_at"
	RunQueryParamsSelectFeedbackStats          RunQueryParamsSelect = "feedback_stats"
	RunQueryParamsSelectChildRunIDs            RunQueryParamsSelect = "child_run_ids"
	RunQueryParamsSelectParentRunIDs           RunQueryParamsSelect = "parent_run_ids"
	RunQueryParamsSelectTags                   RunQueryParamsSelect = "tags"
	RunQueryParamsSelectInDataset              RunQueryParamsSelect = "in_dataset"
	RunQueryParamsSelectAppPath                RunQueryParamsSelect = "app_path"
	RunQueryParamsSelectShareToken             RunQueryParamsSelect = "share_token"
	RunQueryParamsSelectTraceTier              RunQueryParamsSelect = "trace_tier"
	RunQueryParamsSelectTraceFirstReceivedAt   RunQueryParamsSelect = "trace_first_received_at"
	RunQueryParamsSelectTtlSeconds             RunQueryParamsSelect = "ttl_seconds"
	RunQueryParamsSelectTraceUpgrade           RunQueryParamsSelect = "trace_upgrade"
	RunQueryParamsSelectThreadID               RunQueryParamsSelect = "thread_id"
	RunQueryParamsSelectTraceMinMaxStartTime   RunQueryParamsSelect = "trace_min_max_start_time"
	RunQueryParamsSelectMessages               RunQueryParamsSelect = "messages"
	RunQueryParamsSelectInsertedAt             RunQueryParamsSelect = "inserted_at"
)

func (r RunQueryParamsSelect) IsKnown() bool {
	switch r {
	case RunQueryParamsSelectID, RunQueryParamsSelectName, RunQueryParamsSelectRunType, RunQueryParamsSelectStartTime, RunQueryParamsSelectEndTime, RunQueryParamsSelectStatus, RunQueryParamsSelectError, RunQueryParamsSelectExtra, RunQueryParamsSelectEvents, RunQueryParamsSelectInputs, RunQueryParamsSelectInputsPreview, RunQueryParamsSelectInputsS3URLs, RunQueryParamsSelectInputsOrSignedURL, RunQueryParamsSelectOutputs, RunQueryParamsSelectOutputsPreview, RunQueryParamsSelectOutputsS3URLs, RunQueryParamsSelectOutputsOrSignedURL, RunQueryParamsSelectS3URLs, RunQueryParamsSelectErrorOrSignedURL, RunQueryParamsSelectEventsOrSignedURL, RunQueryParamsSelectExtraOrSignedURL, RunQueryParamsSelectSerializedOrSignedURL, RunQueryParamsSelectParentRunID, RunQueryParamsSelectManifestID, RunQueryParamsSelectManifestS3ID, RunQueryParamsSelectManifest, RunQueryParamsSelectSessionID, RunQueryParamsSelectSerialized, RunQueryParamsSelectReferenceExampleID, RunQueryParamsSelectReferenceDatasetID, RunQueryParamsSelectTotalTokens, RunQueryParamsSelectPromptTokens, RunQueryParamsSelectPromptTokenDetails, RunQueryParamsSelectCompletionTokens, RunQueryParamsSelectCompletionTokenDetails, RunQueryParamsSelectTotalCost, RunQueryParamsSelectPromptCost, RunQueryParamsSelectPromptCostDetails, RunQueryParamsSelectCompletionCost, RunQueryParamsSelectCompletionCostDetails, RunQueryParamsSelectPriceModelID, RunQueryParamsSelectFirstTokenTime, RunQueryParamsSelectTraceID, RunQueryParamsSelectDottedOrder, RunQueryParamsSelectLastQueuedAt, RunQueryParamsSelectFeedbackStats, RunQueryParamsSelectChildRunIDs, RunQueryParamsSelectParentRunIDs, RunQueryParamsSelectTags, RunQueryParamsSelectInDataset, RunQueryParamsSelectAppPath, RunQueryParamsSelectShareToken, RunQueryParamsSelectTraceTier, RunQueryParamsSelectTraceFirstReceivedAt, RunQueryParamsSelectTtlSeconds, RunQueryParamsSelectTraceUpgrade, RunQueryParamsSelectThreadID, RunQueryParamsSelectTraceMinMaxStartTime, RunQueryParamsSelectMessages, RunQueryParamsSelectInsertedAt:
		return true
	}
	return false
}

type RunStatsParams struct {
	// Query params for run stats.
	RunStatsQueryParams RunStatsQueryParams `json:"run_stats_query_params" api:"required"`
}

func (r RunStatsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.RunStatsQueryParams)
}
