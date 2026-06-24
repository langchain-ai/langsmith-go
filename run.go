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
	"github.com/langchain-ai/langsmith-go/packages/pagination"
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
func (r *RunService) QueryV1(ctx context.Context, body RunQueryV1Params, opts ...option.RequestOption) (res *pagination.CursorPagination[RunSchema], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/runs/query"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, &res, opts...)
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

// Query Runs
func (r *RunService) QueryV1AutoPaging(ctx context.Context, body RunQueryV1Params, opts ...option.RequestOption) *pagination.CursorPaginationAutoPager[RunSchema] {
	return pagination.NewCursorPaginationAutoPager(r.QueryV1(ctx, body, opts...))
}

// **Alpha:** The request and response contract may change; Returns a paginated
// list of runs for the given projects within min/max start_time. Supports filters,
// cursor pagination, and `selects` to select fields to return.
func (r *RunService) QueryV2(ctx context.Context, params RunQueryV2Params, opts ...option.RequestOption) (res *pagination.ItemsCursorPostPagination[QueryRunResponse], err error) {
	var raw *http.Response
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v2/runs/query"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, params, &res, opts...)
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

// **Alpha:** The request and response contract may change; Returns a paginated
// list of runs for the given projects within min/max start_time. Supports filters,
// cursor pagination, and `selects` to select fields to return.
func (r *RunService) QueryV2AutoPaging(ctx context.Context, params RunQueryV2Params, opts ...option.RequestOption) *pagination.ItemsCursorPostPaginationAutoPager[QueryRunResponse] {
	return pagination.NewItemsCursorPostPaginationAutoPager(r.QueryV2(ctx, params, opts...))
}

// Get a specific run.
func (r *RunService) GetV1(ctx context.Context, runID string, query RunGetV1Params, opts ...option.RequestOption) (res *RunSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/runs/%s", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// **Alpha:** The request and response contract may change; Returns one run by ID
// for the given session and start_time. Use the `selects` query parameter
// (repeatable) to select fields to return.
func (r *RunService) GetV2(ctx context.Context, runID string, params RunGetV2Params, opts ...option.RequestOption) (res *QueryRunResponse, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/runs/%s", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
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

// Get a specific run.
func (r *RunService) Get(ctx context.Context, runID string, body RunGetParams, opts ...option.RequestOption) (res *RunSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/runs/%s", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, body, &res, opts...)
	return res, err
}

// Query Runs
func (r *RunService) Query(ctx context.Context, body RunQueryParams, opts ...option.RequestOption) (res *pagination.CursorPagination[RunSchema], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/runs/query"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, &res, opts...)
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

// Query Runs
func (r *RunService) QueryAutoPaging(ctx context.Context, body RunQueryParams, opts ...option.RequestOption) *pagination.CursorPaginationAutoPager[RunSchema] {
	return pagination.NewCursorPaginationAutoPager(r.Query(ctx, body, opts...))
}

type QueryRunResponse struct {
	// `id` is this run's UUID.
	ID string `json:"id" format:"uuid"`
	// `app_path` identifies the application code location that produced this run, if
	// recorded.
	AppPath string `json:"app_path"`
	// `attachments` maps each attachment file name to a pre-signed HTTPS download URL.
	Attachments map[string]string `json:"attachments"`
	// `completion_cost` is estimated USD cost for the completion.
	CompletionCost float64 `json:"completion_cost"`
	// `completion_cost_details` is the per-category USD breakdown of
	// `completion_cost`. Categories mirror `completion_token_details`. Returned only
	// when the `COMPLETION_COST_DETAILS` field is requested.
	CompletionCostDetails QueryRunResponseCompletionCostDetails `json:"completion_cost_details"`
	// `completion_token_details` is the per-category breakdown of `completion_tokens`.
	// Category names are model-specific (for example `reasoning`, `audio`). Returned
	// only when the `COMPLETION_TOKEN_DETAILS` field is requested.
	CompletionTokenDetails QueryRunResponseCompletionTokenDetails `json:"completion_token_details"`
	// `completion_tokens` is the completion-side token count.
	CompletionTokens int64 `json:"completion_tokens"`
	// `dotted_order` is the hierarchical ordering key for trace trees.
	DottedOrder string `json:"dotted_order"`
	// `end_time` is when the run ended (RFC3339 date-time). JSON null if the run has
	// not finished yet.
	EndTime time.Time `json:"end_time" format:"date-time"`
	// `error` is the error message when `status` indicates failure.
	Error string `json:"error"`
	// `error_preview` is a truncated plain-text error snippet.
	ErrorPreview string `json:"error_preview"`
	// `events` is the ordered list of run events (for example streaming tokens).
	Events []QueryRunResponseEvent `json:"events"`
	// `extra` is additional runtime JSON attached to the run.
	Extra interface{} `json:"extra"`
	// `feedback_stats` aggregates feedback scores keyed by feedback key.
	FeedbackStats map[string]QueryRunResponseFeedbackStat `json:"feedback_stats"`
	// `first_token_time` is when the first output token was produced (RFC3339
	// date-time), when recorded for streamed runs.
	FirstTokenTime time.Time `json:"first_token_time" format:"date-time"`
	// `inputs` is the run input payload (arbitrary JSON object).
	Inputs interface{} `json:"inputs"`
	// `inputs_preview` is a truncated plain-text preview of inputs.
	InputsPreview string `json:"inputs_preview"`
	// `is_in_dataset` is true when this run is linked to a dataset example.
	IsInDataset bool `json:"is_in_dataset"`
	// `is_root` is true when this run has no parent (it is the trace root).
	IsRoot bool `json:"is_root"`
	// `latency_seconds` is wall-clock duration from start to end in seconds.
	LatencySeconds float64 `json:"latency_seconds"`
	// `manifest` is the serialized configuration of the traced component (for example
	// the model parameters, prompt template, or pipeline definition), when recorded.
	Manifest interface{} `json:"manifest"`
	// `metadata` is arbitrary user-defined JSON metadata.
	Metadata interface{} `json:"metadata"`
	// `name` is a human-readable label for the run (for example the model name,
	// function name, or step name chosen when the run was traced).
	Name string `json:"name"`
	// `outputs` is the run output payload (arbitrary JSON object).
	Outputs interface{} `json:"outputs"`
	// `outputs_preview` is a truncated plain-text preview of outputs.
	OutputsPreview string `json:"outputs_preview"`
	// `parent_run_ids` lists ancestor run UUIDs from the trace root down to the direct
	// parent.
	ParentRunIDs []string `json:"parent_run_ids" format:"uuid"`
	// `price_model_id` identifies the pricing model UUID used for cost estimates, when
	// recorded.
	PriceModelID string `json:"price_model_id" format:"uuid"`
	// `project_id` is the tracing project UUID this run was logged to.
	ProjectID string `json:"project_id" format:"uuid"`
	// `prompt_cost` is estimated USD cost for the prompt.
	PromptCost float64 `json:"prompt_cost"`
	// `prompt_cost_details` is the per-category USD breakdown of `prompt_cost`.
	// Categories mirror `prompt_token_details`. Returned only when the
	// `PROMPT_COST_DETAILS` field is requested.
	PromptCostDetails QueryRunResponsePromptCostDetails `json:"prompt_cost_details"`
	// `prompt_token_details` is the per-category breakdown of `prompt_tokens`.
	// Category names are model-specific (for example `cache_read`, `cache_write`).
	// Returned only when the `PROMPT_TOKEN_DETAILS` field is requested.
	PromptTokenDetails QueryRunResponsePromptTokenDetails `json:"prompt_token_details"`
	// `prompt_tokens` is the prompt-side token count.
	PromptTokens int64 `json:"prompt_tokens"`
	// `reference_dataset_id` is the dataset UUID for the reference example, if any.
	ReferenceDatasetID string `json:"reference_dataset_id" format:"uuid"`
	// `reference_example_id` is the dataset example UUID this run was compared
	// against, if any.
	ReferenceExampleID string `json:"reference_example_id" format:"uuid"`
	// `run_type` identifies what kind of operation this run represents (for example an
	// LLM call, a tool invocation, or a chain step). See the `RunType` enum for
	// allowed values.
	RunType QueryRunResponseRunType `json:"run_type"`
	// `share_url` is the fully-qualified URL of this run's public view, rooted at the
	// deployment's LangSmith app origin (for example
	// `https://smith.langchain.com/public/4f7a1b2c-8d9e-4a0b-9c1d-2e3f4a5b6c7d/r`). It
	// is returned only when `SHARE_URL` is included in `selects`, and only when the
	// run has been explicitly shared; the URL remains stable until the run is
	// unshared. Anyone with this URL can view the run anonymously, so treat it as a
	// secret and do not log it.
	ShareURL string `json:"share_url"`
	// `start_time` is when the run started (RFC3339 date-time).
	StartTime time.Time `json:"start_time" format:"date-time"`
	// `status` is the completion status of the run.
	Status QueryRunResponseStatus `json:"status"`
	// `tags` lists user-defined tags on this run.
	Tags []string `json:"tags"`
	// `thread_evaluation_time` is thread-level evaluation timing (RFC3339 date-time),
	// when recorded.
	ThreadEvaluationTime time.Time `json:"thread_evaluation_time" format:"date-time"`
	// `thread_id` is the conversation thread UUID this run belongs to, if any.
	ThreadID string `json:"thread_id" format:"uuid"`
	// `total_cost` is total estimated USD cost (prompt plus completion).
	TotalCost float64 `json:"total_cost"`
	// `total_tokens` is prompt plus completion tokens.
	TotalTokens int64 `json:"total_tokens"`
	// `trace_id` is the root trace UUID; for a root run it matches `id`.
	TraceID string               `json:"trace_id" format:"uuid"`
	JSON    queryRunResponseJSON `json:"-"`
}

// queryRunResponseJSON contains the JSON metadata for the struct
// [QueryRunResponse]
type queryRunResponseJSON struct {
	ID                     apijson.Field
	AppPath                apijson.Field
	Attachments            apijson.Field
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	DottedOrder            apijson.Field
	EndTime                apijson.Field
	Error                  apijson.Field
	ErrorPreview           apijson.Field
	Events                 apijson.Field
	Extra                  apijson.Field
	FeedbackStats          apijson.Field
	FirstTokenTime         apijson.Field
	Inputs                 apijson.Field
	InputsPreview          apijson.Field
	IsInDataset            apijson.Field
	IsRoot                 apijson.Field
	LatencySeconds         apijson.Field
	Manifest               apijson.Field
	Metadata               apijson.Field
	Name                   apijson.Field
	Outputs                apijson.Field
	OutputsPreview         apijson.Field
	ParentRunIDs           apijson.Field
	PriceModelID           apijson.Field
	ProjectID              apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	ReferenceDatasetID     apijson.Field
	ReferenceExampleID     apijson.Field
	RunType                apijson.Field
	ShareURL               apijson.Field
	StartTime              apijson.Field
	Status                 apijson.Field
	Tags                   apijson.Field
	ThreadEvaluationTime   apijson.Field
	ThreadID               apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	TraceID                apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *QueryRunResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponseJSON) RawJSON() string {
	return r.raw
}

// `completion_cost_details` is the per-category USD breakdown of
// `completion_cost`. Categories mirror `completion_token_details`. Returned only
// when the `COMPLETION_COST_DETAILS` field is requested.
type QueryRunResponseCompletionCostDetails struct {
	// `raw` maps each category name to its estimated USD cost.
	Raw  map[string]float64                        `json:"raw"`
	JSON queryRunResponseCompletionCostDetailsJSON `json:"-"`
}

// queryRunResponseCompletionCostDetailsJSON contains the JSON metadata for the
// struct [QueryRunResponseCompletionCostDetails]
type queryRunResponseCompletionCostDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *QueryRunResponseCompletionCostDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponseCompletionCostDetailsJSON) RawJSON() string {
	return r.raw
}

// `completion_token_details` is the per-category breakdown of `completion_tokens`.
// Category names are model-specific (for example `reasoning`, `audio`). Returned
// only when the `COMPLETION_TOKEN_DETAILS` field is requested.
type QueryRunResponseCompletionTokenDetails struct {
	// `raw` maps each category name to its completion-token count.
	Raw  map[string]int64                           `json:"raw"`
	JSON queryRunResponseCompletionTokenDetailsJSON `json:"-"`
}

// queryRunResponseCompletionTokenDetailsJSON contains the JSON metadata for the
// struct [QueryRunResponseCompletionTokenDetails]
type queryRunResponseCompletionTokenDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *QueryRunResponseCompletionTokenDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponseCompletionTokenDetailsJSON) RawJSON() string {
	return r.raw
}

type QueryRunResponseEvent struct {
	// `kwargs` is the event payload — an opaque JSON object whose shape depends on
	// `name` and on the emitting SDK. For example LangChain emits `{"token": {...}}`
	// for `new_token` events, tool-call start/end details for tool events, and
	// arbitrary user-defined payloads for custom events. Clients should treat `kwargs`
	// as untyped JSON: do not assume specific keys exist for a given `name`, and
	// tolerate additional unknown keys appearing over time.
	Kwargs interface{} `json:"kwargs"`
	// `name` is the event kind. Common values emitted by the LangChain/LangSmith
	// tracer SDKs include `"start"`, `"end"`, and `"new_token"`, but applications may
	// emit arbitrary strings for their own instrumentation.
	Name string `json:"name"`
	// `time` is when the event occurred (RFC3339 date-time with millisecond
	// precision).
	Time time.Time                 `json:"time" format:"date-time"`
	JSON queryRunResponseEventJSON `json:"-"`
}

// queryRunResponseEventJSON contains the JSON metadata for the struct
// [QueryRunResponseEvent]
type queryRunResponseEventJSON struct {
	Kwargs      apijson.Field
	Name        apijson.Field
	Time        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *QueryRunResponseEvent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponseEventJSON) RawJSON() string {
	return r.raw
}

type QueryRunResponseFeedbackStat struct {
	// `avg` is the arithmetic mean of numeric feedback scores for this key on the run,
	// or `null` when no numeric score has been recorded (for example purely
	// categorical feedback).
	Avg float64 `json:"avg"`
	// `comments` is a sample of human-readable comments attached to feedback points
	// for this key, in no particular order. May be empty; is not exhaustive when many
	// comments exist.
	Comments []string `json:"comments"`
	// `contains_thread_feedback` is true when at least one feedback point for this key
	// was submitted at the thread level (rather than at an individual run). Always
	// false on responses that already describe a single run in isolation.
	ContainsThreadFeedback bool `json:"contains_thread_feedback"`
	// `errors` is the number of feedback points recorded as errors rather than
	// successful scores (for example an automated evaluator that raised an exception).
	// Defaults to 0 when no errors occurred.
	Errors int64 `json:"errors"`
	// `max` is the largest numeric feedback score recorded for this key on the run, or
	// `null` when no numeric score has been recorded.
	Max float64 `json:"max"`
	// `min` is the smallest numeric feedback score recorded for this key on the run,
	// or `null` when no numeric score has been recorded.
	Min float64 `json:"min"`
	// `n` is the number of feedback points recorded for this key on the run. For
	// numeric feedback this is the sample size behind `avg`, `min`, `max`, and
	// `stdev`; for categorical feedback it is the sum of the `values` counts.
	N int64 `json:"n"`
	// `sources` is a sample of feedback sources for this key. Each entry is either a
	// plain string identifier (for example `"api"`, `"app"`, `"model"`) or a JSON
	// object describing a synthetic source (for example
	// `{"type": "__ls_composite_feedback"}` for a computed aggregate). Clients must
	// tolerate both shapes.
	Sources []interface{} `json:"sources"`
	// `stdev` is the sample standard deviation of numeric feedback scores for this key
	// on the run, or `null` when it cannot be computed (for example fewer than two
	// numeric scores, or purely categorical feedback).
	Stdev float64 `json:"stdev"`
	// `values` is the distribution of categorical feedback labels for this key,
	// mapping each label to its occurrence count. Empty (`{}`) for purely numeric
	// feedback.
	Values map[string]int64                 `json:"values"`
	JSON   queryRunResponseFeedbackStatJSON `json:"-"`
}

// queryRunResponseFeedbackStatJSON contains the JSON metadata for the struct
// [QueryRunResponseFeedbackStat]
type queryRunResponseFeedbackStatJSON struct {
	Avg                    apijson.Field
	Comments               apijson.Field
	ContainsThreadFeedback apijson.Field
	Errors                 apijson.Field
	Max                    apijson.Field
	Min                    apijson.Field
	N                      apijson.Field
	Sources                apijson.Field
	Stdev                  apijson.Field
	Values                 apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *QueryRunResponseFeedbackStat) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponseFeedbackStatJSON) RawJSON() string {
	return r.raw
}

// `prompt_cost_details` is the per-category USD breakdown of `prompt_cost`.
// Categories mirror `prompt_token_details`. Returned only when the
// `PROMPT_COST_DETAILS` field is requested.
type QueryRunResponsePromptCostDetails struct {
	// `raw` maps each category name to its estimated USD cost.
	Raw  map[string]float64                    `json:"raw"`
	JSON queryRunResponsePromptCostDetailsJSON `json:"-"`
}

// queryRunResponsePromptCostDetailsJSON contains the JSON metadata for the struct
// [QueryRunResponsePromptCostDetails]
type queryRunResponsePromptCostDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *QueryRunResponsePromptCostDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponsePromptCostDetailsJSON) RawJSON() string {
	return r.raw
}

// `prompt_token_details` is the per-category breakdown of `prompt_tokens`.
// Category names are model-specific (for example `cache_read`, `cache_write`).
// Returned only when the `PROMPT_TOKEN_DETAILS` field is requested.
type QueryRunResponsePromptTokenDetails struct {
	// `raw` maps each category name to its prompt-token count.
	Raw  map[string]int64                       `json:"raw"`
	JSON queryRunResponsePromptTokenDetailsJSON `json:"-"`
}

// queryRunResponsePromptTokenDetailsJSON contains the JSON metadata for the struct
// [QueryRunResponsePromptTokenDetails]
type queryRunResponsePromptTokenDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *QueryRunResponsePromptTokenDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r queryRunResponsePromptTokenDetailsJSON) RawJSON() string {
	return r.raw
}

// `run_type` identifies what kind of operation this run represents (for example an
// LLM call, a tool invocation, or a chain step). See the `RunType` enum for
// allowed values.
type QueryRunResponseRunType string

const (
	QueryRunResponseRunTypeTool      QueryRunResponseRunType = "TOOL"
	QueryRunResponseRunTypeChain     QueryRunResponseRunType = "CHAIN"
	QueryRunResponseRunTypeLlm       QueryRunResponseRunType = "LLM"
	QueryRunResponseRunTypeRetriever QueryRunResponseRunType = "RETRIEVER"
	QueryRunResponseRunTypeEmbedding QueryRunResponseRunType = "EMBEDDING"
	QueryRunResponseRunTypePrompt    QueryRunResponseRunType = "PROMPT"
	QueryRunResponseRunTypeParser    QueryRunResponseRunType = "PARSER"
)

func (r QueryRunResponseRunType) IsKnown() bool {
	switch r {
	case QueryRunResponseRunTypeTool, QueryRunResponseRunTypeChain, QueryRunResponseRunTypeLlm, QueryRunResponseRunTypeRetriever, QueryRunResponseRunTypeEmbedding, QueryRunResponseRunTypePrompt, QueryRunResponseRunTypeParser:
		return true
	}
	return false
}

// `status` is the completion status of the run.
type QueryRunResponseStatus string

const (
	QueryRunResponseStatusSuccess QueryRunResponseStatus = "SUCCESS"
	QueryRunResponseStatusError   QueryRunResponseStatus = "ERROR"
	QueryRunResponseStatusPending QueryRunResponseStatus = "PENDING"
)

func (r QueryRunResponseStatus) IsKnown() bool {
	switch r {
	case QueryRunResponseStatusSuccess, QueryRunResponseStatusError, QueryRunResponseStatusPending:
		return true
	}
	return false
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
	IncludeDetails   param.Field[bool]                 `json:"include_details"`
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
	RunStatsQueryParamsSelectGroupCount             RunStatsQueryParamsSelect = "group_count"
	RunStatsQueryParamsSelectPromptTokenDetails     RunStatsQueryParamsSelect = "prompt_token_details"
	RunStatsQueryParamsSelectCompletionTokenDetails RunStatsQueryParamsSelect = "completion_token_details"
	RunStatsQueryParamsSelectPromptCostDetails      RunStatsQueryParamsSelect = "prompt_cost_details"
	RunStatsQueryParamsSelectCompletionCostDetails  RunStatsQueryParamsSelect = "completion_cost_details"
)

func (r RunStatsQueryParamsSelect) IsKnown() bool {
	switch r {
	case RunStatsQueryParamsSelectRunCount, RunStatsQueryParamsSelectLatencyP50, RunStatsQueryParamsSelectLatencyP99, RunStatsQueryParamsSelectLatencyAvg, RunStatsQueryParamsSelectFirstTokenP50, RunStatsQueryParamsSelectFirstTokenP99, RunStatsQueryParamsSelectTotalTokens, RunStatsQueryParamsSelectPromptTokens, RunStatsQueryParamsSelectCompletionTokens, RunStatsQueryParamsSelectMedianTokens, RunStatsQueryParamsSelectCompletionTokensP50, RunStatsQueryParamsSelectPromptTokensP50, RunStatsQueryParamsSelectTokensP99, RunStatsQueryParamsSelectCompletionTokensP99, RunStatsQueryParamsSelectPromptTokensP99, RunStatsQueryParamsSelectLastRunStartTime, RunStatsQueryParamsSelectFeedbackStats, RunStatsQueryParamsSelectThreadFeedbackStats, RunStatsQueryParamsSelectRunFacets, RunStatsQueryParamsSelectErrorRate, RunStatsQueryParamsSelectStreamingRate, RunStatsQueryParamsSelectTotalCost, RunStatsQueryParamsSelectPromptCost, RunStatsQueryParamsSelectCompletionCost, RunStatsQueryParamsSelectCostP50, RunStatsQueryParamsSelectCostP99, RunStatsQueryParamsSelectSessionFeedbackStats, RunStatsQueryParamsSelectAllRunStats, RunStatsQueryParamsSelectAllTokenStats, RunStatsQueryParamsSelectGroupCount, RunStatsQueryParamsSelectPromptTokenDetails, RunStatsQueryParamsSelectCompletionTokenDetails, RunStatsQueryParamsSelectPromptCostDetails, RunStatsQueryParamsSelectCompletionCostDetails:
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

type RunQueryV1Params struct {
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
	Order            param.Field[RunQueryV1ParamsOrder] `json:"order"`
	ParentRun        param.Field[string]                `json:"parent_run" format:"uuid"`
	Query            param.Field[string]                `json:"query"`
	ReferenceExample param.Field[[]string]              `json:"reference_example" format:"uuid"`
	// Enum for run types.
	RunType               param.Field[RunTypeEnum]              `json:"run_type"`
	SearchFilter          param.Field[string]                   `json:"search_filter"`
	Select                param.Field[[]RunQueryV1ParamsSelect] `json:"select"`
	Session               param.Field[[]string]                 `json:"session" format:"uuid"`
	SkipPagination        param.Field[bool]                     `json:"skip_pagination"`
	SkipPrevCursor        param.Field[bool]                     `json:"skip_prev_cursor"`
	StartTime             param.Field[time.Time]                `json:"start_time" format:"date-time"`
	Trace                 param.Field[string]                   `json:"trace" format:"uuid"`
	TraceFilter           param.Field[string]                   `json:"trace_filter"`
	TreeFilter            param.Field[string]                   `json:"tree_filter"`
	UseExperimentalSearch param.Field[bool]                     `json:"use_experimental_search"`
}

func (r RunQueryV1Params) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Enum for run start date order.
type RunQueryV1ParamsOrder string

const (
	RunQueryV1ParamsOrderAsc  RunQueryV1ParamsOrder = "asc"
	RunQueryV1ParamsOrderDesc RunQueryV1ParamsOrder = "desc"
)

func (r RunQueryV1ParamsOrder) IsKnown() bool {
	switch r {
	case RunQueryV1ParamsOrderAsc, RunQueryV1ParamsOrderDesc:
		return true
	}
	return false
}

// Enum for available run columns.
type RunQueryV1ParamsSelect string

const (
	RunQueryV1ParamsSelectID                     RunQueryV1ParamsSelect = "id"
	RunQueryV1ParamsSelectName                   RunQueryV1ParamsSelect = "name"
	RunQueryV1ParamsSelectRunType                RunQueryV1ParamsSelect = "run_type"
	RunQueryV1ParamsSelectStartTime              RunQueryV1ParamsSelect = "start_time"
	RunQueryV1ParamsSelectEndTime                RunQueryV1ParamsSelect = "end_time"
	RunQueryV1ParamsSelectStatus                 RunQueryV1ParamsSelect = "status"
	RunQueryV1ParamsSelectError                  RunQueryV1ParamsSelect = "error"
	RunQueryV1ParamsSelectExtra                  RunQueryV1ParamsSelect = "extra"
	RunQueryV1ParamsSelectEvents                 RunQueryV1ParamsSelect = "events"
	RunQueryV1ParamsSelectInputs                 RunQueryV1ParamsSelect = "inputs"
	RunQueryV1ParamsSelectInputsPreview          RunQueryV1ParamsSelect = "inputs_preview"
	RunQueryV1ParamsSelectInputsS3URLs           RunQueryV1ParamsSelect = "inputs_s3_urls"
	RunQueryV1ParamsSelectInputsOrSignedURL      RunQueryV1ParamsSelect = "inputs_or_signed_url"
	RunQueryV1ParamsSelectOutputs                RunQueryV1ParamsSelect = "outputs"
	RunQueryV1ParamsSelectOutputsPreview         RunQueryV1ParamsSelect = "outputs_preview"
	RunQueryV1ParamsSelectOutputsS3URLs          RunQueryV1ParamsSelect = "outputs_s3_urls"
	RunQueryV1ParamsSelectOutputsOrSignedURL     RunQueryV1ParamsSelect = "outputs_or_signed_url"
	RunQueryV1ParamsSelectS3URLs                 RunQueryV1ParamsSelect = "s3_urls"
	RunQueryV1ParamsSelectErrorOrSignedURL       RunQueryV1ParamsSelect = "error_or_signed_url"
	RunQueryV1ParamsSelectEventsOrSignedURL      RunQueryV1ParamsSelect = "events_or_signed_url"
	RunQueryV1ParamsSelectExtraOrSignedURL       RunQueryV1ParamsSelect = "extra_or_signed_url"
	RunQueryV1ParamsSelectSerializedOrSignedURL  RunQueryV1ParamsSelect = "serialized_or_signed_url"
	RunQueryV1ParamsSelectParentRunID            RunQueryV1ParamsSelect = "parent_run_id"
	RunQueryV1ParamsSelectManifestID             RunQueryV1ParamsSelect = "manifest_id"
	RunQueryV1ParamsSelectManifestS3ID           RunQueryV1ParamsSelect = "manifest_s3_id"
	RunQueryV1ParamsSelectManifest               RunQueryV1ParamsSelect = "manifest"
	RunQueryV1ParamsSelectSessionID              RunQueryV1ParamsSelect = "session_id"
	RunQueryV1ParamsSelectSerialized             RunQueryV1ParamsSelect = "serialized"
	RunQueryV1ParamsSelectReferenceExampleID     RunQueryV1ParamsSelect = "reference_example_id"
	RunQueryV1ParamsSelectReferenceDatasetID     RunQueryV1ParamsSelect = "reference_dataset_id"
	RunQueryV1ParamsSelectTotalTokens            RunQueryV1ParamsSelect = "total_tokens"
	RunQueryV1ParamsSelectPromptTokens           RunQueryV1ParamsSelect = "prompt_tokens"
	RunQueryV1ParamsSelectPromptTokenDetails     RunQueryV1ParamsSelect = "prompt_token_details"
	RunQueryV1ParamsSelectCompletionTokens       RunQueryV1ParamsSelect = "completion_tokens"
	RunQueryV1ParamsSelectCompletionTokenDetails RunQueryV1ParamsSelect = "completion_token_details"
	RunQueryV1ParamsSelectTotalCost              RunQueryV1ParamsSelect = "total_cost"
	RunQueryV1ParamsSelectPromptCost             RunQueryV1ParamsSelect = "prompt_cost"
	RunQueryV1ParamsSelectPromptCostDetails      RunQueryV1ParamsSelect = "prompt_cost_details"
	RunQueryV1ParamsSelectCompletionCost         RunQueryV1ParamsSelect = "completion_cost"
	RunQueryV1ParamsSelectCompletionCostDetails  RunQueryV1ParamsSelect = "completion_cost_details"
	RunQueryV1ParamsSelectPriceModelID           RunQueryV1ParamsSelect = "price_model_id"
	RunQueryV1ParamsSelectFirstTokenTime         RunQueryV1ParamsSelect = "first_token_time"
	RunQueryV1ParamsSelectTraceID                RunQueryV1ParamsSelect = "trace_id"
	RunQueryV1ParamsSelectDottedOrder            RunQueryV1ParamsSelect = "dotted_order"
	RunQueryV1ParamsSelectLastQueuedAt           RunQueryV1ParamsSelect = "last_queued_at"
	RunQueryV1ParamsSelectFeedbackStats          RunQueryV1ParamsSelect = "feedback_stats"
	RunQueryV1ParamsSelectChildRunIDs            RunQueryV1ParamsSelect = "child_run_ids"
	RunQueryV1ParamsSelectParentRunIDs           RunQueryV1ParamsSelect = "parent_run_ids"
	RunQueryV1ParamsSelectTags                   RunQueryV1ParamsSelect = "tags"
	RunQueryV1ParamsSelectInDataset              RunQueryV1ParamsSelect = "in_dataset"
	RunQueryV1ParamsSelectAppPath                RunQueryV1ParamsSelect = "app_path"
	RunQueryV1ParamsSelectShareToken             RunQueryV1ParamsSelect = "share_token"
	RunQueryV1ParamsSelectTraceTier              RunQueryV1ParamsSelect = "trace_tier"
	RunQueryV1ParamsSelectTraceFirstReceivedAt   RunQueryV1ParamsSelect = "trace_first_received_at"
	RunQueryV1ParamsSelectTtlSeconds             RunQueryV1ParamsSelect = "ttl_seconds"
	RunQueryV1ParamsSelectTraceUpgrade           RunQueryV1ParamsSelect = "trace_upgrade"
	RunQueryV1ParamsSelectThreadID               RunQueryV1ParamsSelect = "thread_id"
	RunQueryV1ParamsSelectTraceMinMaxStartTime   RunQueryV1ParamsSelect = "trace_min_max_start_time"
	RunQueryV1ParamsSelectMessages               RunQueryV1ParamsSelect = "messages"
	RunQueryV1ParamsSelectInsertedAt             RunQueryV1ParamsSelect = "inserted_at"
)

func (r RunQueryV1ParamsSelect) IsKnown() bool {
	switch r {
	case RunQueryV1ParamsSelectID, RunQueryV1ParamsSelectName, RunQueryV1ParamsSelectRunType, RunQueryV1ParamsSelectStartTime, RunQueryV1ParamsSelectEndTime, RunQueryV1ParamsSelectStatus, RunQueryV1ParamsSelectError, RunQueryV1ParamsSelectExtra, RunQueryV1ParamsSelectEvents, RunQueryV1ParamsSelectInputs, RunQueryV1ParamsSelectInputsPreview, RunQueryV1ParamsSelectInputsS3URLs, RunQueryV1ParamsSelectInputsOrSignedURL, RunQueryV1ParamsSelectOutputs, RunQueryV1ParamsSelectOutputsPreview, RunQueryV1ParamsSelectOutputsS3URLs, RunQueryV1ParamsSelectOutputsOrSignedURL, RunQueryV1ParamsSelectS3URLs, RunQueryV1ParamsSelectErrorOrSignedURL, RunQueryV1ParamsSelectEventsOrSignedURL, RunQueryV1ParamsSelectExtraOrSignedURL, RunQueryV1ParamsSelectSerializedOrSignedURL, RunQueryV1ParamsSelectParentRunID, RunQueryV1ParamsSelectManifestID, RunQueryV1ParamsSelectManifestS3ID, RunQueryV1ParamsSelectManifest, RunQueryV1ParamsSelectSessionID, RunQueryV1ParamsSelectSerialized, RunQueryV1ParamsSelectReferenceExampleID, RunQueryV1ParamsSelectReferenceDatasetID, RunQueryV1ParamsSelectTotalTokens, RunQueryV1ParamsSelectPromptTokens, RunQueryV1ParamsSelectPromptTokenDetails, RunQueryV1ParamsSelectCompletionTokens, RunQueryV1ParamsSelectCompletionTokenDetails, RunQueryV1ParamsSelectTotalCost, RunQueryV1ParamsSelectPromptCost, RunQueryV1ParamsSelectPromptCostDetails, RunQueryV1ParamsSelectCompletionCost, RunQueryV1ParamsSelectCompletionCostDetails, RunQueryV1ParamsSelectPriceModelID, RunQueryV1ParamsSelectFirstTokenTime, RunQueryV1ParamsSelectTraceID, RunQueryV1ParamsSelectDottedOrder, RunQueryV1ParamsSelectLastQueuedAt, RunQueryV1ParamsSelectFeedbackStats, RunQueryV1ParamsSelectChildRunIDs, RunQueryV1ParamsSelectParentRunIDs, RunQueryV1ParamsSelectTags, RunQueryV1ParamsSelectInDataset, RunQueryV1ParamsSelectAppPath, RunQueryV1ParamsSelectShareToken, RunQueryV1ParamsSelectTraceTier, RunQueryV1ParamsSelectTraceFirstReceivedAt, RunQueryV1ParamsSelectTtlSeconds, RunQueryV1ParamsSelectTraceUpgrade, RunQueryV1ParamsSelectThreadID, RunQueryV1ParamsSelectTraceMinMaxStartTime, RunQueryV1ParamsSelectMessages, RunQueryV1ParamsSelectInsertedAt:
		return true
	}
	return false
}

type RunQueryV2Params struct {
	// `cursor` is the opaque string from a previous response's `next_cursor`. Treat it
	// as opaque and pass it back unmodified.
	Cursor param.Field[string] `json:"cursor"`
	// `filter` narrows results to runs matching this LangSmith filter expression,
	// evaluated against each individual run. For example: and(eq(run_type, "llm"),
	// gt(latency, 5)) or eq(status, "error"). See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	Filter param.Field[string] `json:"filter"`
	// `has_error` filters to runs that errored (true) or completed without error
	// (false).
	HasError param.Field[bool] `json:"has_error"`
	// `ids` optionally limits the request to these run UUIDs.
	IDs param.Field[[]string] `json:"ids" format:"uuid"`
	// `is_root` returns only root runs (true) or only non-root runs (false).
	IsRoot param.Field[bool] `json:"is_root"`
	// `max_start_time` is the upper bound for run `start_time` (RFC3339). Defaults to
	// now.
	MaxStartTime param.Field[time.Time] `json:"max_start_time" format:"date-time"`
	// `min_start_time` is the lower bound for run `start_time` (RFC3339). Defaults to
	// 1 day ago.
	MinStartTime param.Field[time.Time] `json:"min_start_time" format:"date-time"`
	// `page_size` is the maximum number of runs to return in this response. Defaults
	// to 100 when omitted; must be between 1 and 1000 inclusive when set.
	PageSize param.Field[int64] `json:"page_size"`
	// `project_ids` lists tracing project UUIDs to query. Required unless
	// `reference_dataset_id` is set. Mutually exclusive with `reference_dataset_id` —
	// set exactly one of them.
	ProjectIDs param.Field[[]string] `json:"project_ids" format:"uuid"`
	// `reference_dataset_id` resolves session IDs server-side from the dataset.
	// Required unless `project_ids` is set. Mutually exclusive with `project_ids` —
	// set exactly one of them. When provided and `min_start_time` is omitted, the
	// server derives it from the earliest session creation date.
	ReferenceDatasetID param.Field[string] `json:"reference_dataset_id" format:"uuid"`
	// `reference_examples` optionally limits to runs linked to these dataset example
	// UUIDs.
	ReferenceExamples param.Field[[]string] `json:"reference_examples" format:"uuid"`
	// `run_type`, when set, restricts results to runs whose `run_type` equals this
	// value.
	RunType param.Field[RunQueryV2ParamsRunType] `json:"run_type"`
	// `selects` lists which properties to include on each returned run. If omitted,
	// only `id` is returned. Properties not listed are omitted from each run object.
	Selects param.Field[[]RunQueryV2ParamsSelect] `json:"selects"`
	// `trace_filter` narrows results to runs whose root trace matches this LangSmith
	// filter expression. Use this to filter by properties of the trace's root run —
	// for example eq(status, "success") to include only traces that completed without
	// error. See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	TraceFilter param.Field[string] `json:"trace_filter"`
	// `trace_id` optionally limits results to runs belonging to this trace UUID.
	TraceID param.Field[string] `json:"trace_id" format:"uuid"`
	// `tree_filter` narrows results to runs that belong to a trace containing at least
	// one run matching this LangSmith filter expression anywhere in the run tree (not
	// just the root). Use this to find runs inside traces that involved a specific
	// tool, tag, or model — for example has(tags, "production") or eq(name,
	// "my_tool"). See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	TreeFilter param.Field[string] `json:"tree_filter"`
	Accept     param.Field[string] `header:"Accept"`
}

func (r RunQueryV2Params) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// `run_type`, when set, restricts results to runs whose `run_type` equals this
// value.
type RunQueryV2ParamsRunType string

const (
	RunQueryV2ParamsRunTypeTool      RunQueryV2ParamsRunType = "TOOL"
	RunQueryV2ParamsRunTypeChain     RunQueryV2ParamsRunType = "CHAIN"
	RunQueryV2ParamsRunTypeLlm       RunQueryV2ParamsRunType = "LLM"
	RunQueryV2ParamsRunTypeRetriever RunQueryV2ParamsRunType = "RETRIEVER"
	RunQueryV2ParamsRunTypeEmbedding RunQueryV2ParamsRunType = "EMBEDDING"
	RunQueryV2ParamsRunTypePrompt    RunQueryV2ParamsRunType = "PROMPT"
	RunQueryV2ParamsRunTypeParser    RunQueryV2ParamsRunType = "PARSER"
)

func (r RunQueryV2ParamsRunType) IsKnown() bool {
	switch r {
	case RunQueryV2ParamsRunTypeTool, RunQueryV2ParamsRunTypeChain, RunQueryV2ParamsRunTypeLlm, RunQueryV2ParamsRunTypeRetriever, RunQueryV2ParamsRunTypeEmbedding, RunQueryV2ParamsRunTypePrompt, RunQueryV2ParamsRunTypeParser:
		return true
	}
	return false
}

type RunQueryV2ParamsSelect string

const (
	RunQueryV2ParamsSelectID                     RunQueryV2ParamsSelect = "ID"
	RunQueryV2ParamsSelectName                   RunQueryV2ParamsSelect = "NAME"
	RunQueryV2ParamsSelectRunType                RunQueryV2ParamsSelect = "RUN_TYPE"
	RunQueryV2ParamsSelectStatus                 RunQueryV2ParamsSelect = "STATUS"
	RunQueryV2ParamsSelectStartTime              RunQueryV2ParamsSelect = "START_TIME"
	RunQueryV2ParamsSelectEndTime                RunQueryV2ParamsSelect = "END_TIME"
	RunQueryV2ParamsSelectLatencySeconds         RunQueryV2ParamsSelect = "LATENCY_SECONDS"
	RunQueryV2ParamsSelectFirstTokenTime         RunQueryV2ParamsSelect = "FIRST_TOKEN_TIME"
	RunQueryV2ParamsSelectError                  RunQueryV2ParamsSelect = "ERROR"
	RunQueryV2ParamsSelectErrorPreview           RunQueryV2ParamsSelect = "ERROR_PREVIEW"
	RunQueryV2ParamsSelectExtra                  RunQueryV2ParamsSelect = "EXTRA"
	RunQueryV2ParamsSelectMetadata               RunQueryV2ParamsSelect = "METADATA"
	RunQueryV2ParamsSelectEvents                 RunQueryV2ParamsSelect = "EVENTS"
	RunQueryV2ParamsSelectInputs                 RunQueryV2ParamsSelect = "INPUTS"
	RunQueryV2ParamsSelectInputsPreview          RunQueryV2ParamsSelect = "INPUTS_PREVIEW"
	RunQueryV2ParamsSelectOutputs                RunQueryV2ParamsSelect = "OUTPUTS"
	RunQueryV2ParamsSelectOutputsPreview         RunQueryV2ParamsSelect = "OUTPUTS_PREVIEW"
	RunQueryV2ParamsSelectManifest               RunQueryV2ParamsSelect = "MANIFEST"
	RunQueryV2ParamsSelectParentRunIDs           RunQueryV2ParamsSelect = "PARENT_RUN_IDS"
	RunQueryV2ParamsSelectProjectID              RunQueryV2ParamsSelect = "PROJECT_ID"
	RunQueryV2ParamsSelectTraceID                RunQueryV2ParamsSelect = "TRACE_ID"
	RunQueryV2ParamsSelectThreadID               RunQueryV2ParamsSelect = "THREAD_ID"
	RunQueryV2ParamsSelectDottedOrder            RunQueryV2ParamsSelect = "DOTTED_ORDER"
	RunQueryV2ParamsSelectIsRoot                 RunQueryV2ParamsSelect = "IS_ROOT"
	RunQueryV2ParamsSelectReferenceExampleID     RunQueryV2ParamsSelect = "REFERENCE_EXAMPLE_ID"
	RunQueryV2ParamsSelectReferenceDatasetID     RunQueryV2ParamsSelect = "REFERENCE_DATASET_ID"
	RunQueryV2ParamsSelectTotalTokens            RunQueryV2ParamsSelect = "TOTAL_TOKENS"
	RunQueryV2ParamsSelectPromptTokens           RunQueryV2ParamsSelect = "PROMPT_TOKENS"
	RunQueryV2ParamsSelectCompletionTokens       RunQueryV2ParamsSelect = "COMPLETION_TOKENS"
	RunQueryV2ParamsSelectTotalCost              RunQueryV2ParamsSelect = "TOTAL_COST"
	RunQueryV2ParamsSelectPromptCost             RunQueryV2ParamsSelect = "PROMPT_COST"
	RunQueryV2ParamsSelectCompletionCost         RunQueryV2ParamsSelect = "COMPLETION_COST"
	RunQueryV2ParamsSelectPromptTokenDetails     RunQueryV2ParamsSelect = "PROMPT_TOKEN_DETAILS"
	RunQueryV2ParamsSelectCompletionTokenDetails RunQueryV2ParamsSelect = "COMPLETION_TOKEN_DETAILS"
	RunQueryV2ParamsSelectPromptCostDetails      RunQueryV2ParamsSelect = "PROMPT_COST_DETAILS"
	RunQueryV2ParamsSelectCompletionCostDetails  RunQueryV2ParamsSelect = "COMPLETION_COST_DETAILS"
	RunQueryV2ParamsSelectPriceModelID           RunQueryV2ParamsSelect = "PRICE_MODEL_ID"
	RunQueryV2ParamsSelectTags                   RunQueryV2ParamsSelect = "TAGS"
	RunQueryV2ParamsSelectAppPath                RunQueryV2ParamsSelect = "APP_PATH"
	RunQueryV2ParamsSelectAttachments            RunQueryV2ParamsSelect = "ATTACHMENTS"
	RunQueryV2ParamsSelectThreadEvaluationTime   RunQueryV2ParamsSelect = "THREAD_EVALUATION_TIME"
	RunQueryV2ParamsSelectIsInDataset            RunQueryV2ParamsSelect = "IS_IN_DATASET"
	RunQueryV2ParamsSelectShareURL               RunQueryV2ParamsSelect = "SHARE_URL"
	RunQueryV2ParamsSelectFeedbackStats          RunQueryV2ParamsSelect = "FEEDBACK_STATS"
)

func (r RunQueryV2ParamsSelect) IsKnown() bool {
	switch r {
	case RunQueryV2ParamsSelectID, RunQueryV2ParamsSelectName, RunQueryV2ParamsSelectRunType, RunQueryV2ParamsSelectStatus, RunQueryV2ParamsSelectStartTime, RunQueryV2ParamsSelectEndTime, RunQueryV2ParamsSelectLatencySeconds, RunQueryV2ParamsSelectFirstTokenTime, RunQueryV2ParamsSelectError, RunQueryV2ParamsSelectErrorPreview, RunQueryV2ParamsSelectExtra, RunQueryV2ParamsSelectMetadata, RunQueryV2ParamsSelectEvents, RunQueryV2ParamsSelectInputs, RunQueryV2ParamsSelectInputsPreview, RunQueryV2ParamsSelectOutputs, RunQueryV2ParamsSelectOutputsPreview, RunQueryV2ParamsSelectManifest, RunQueryV2ParamsSelectParentRunIDs, RunQueryV2ParamsSelectProjectID, RunQueryV2ParamsSelectTraceID, RunQueryV2ParamsSelectThreadID, RunQueryV2ParamsSelectDottedOrder, RunQueryV2ParamsSelectIsRoot, RunQueryV2ParamsSelectReferenceExampleID, RunQueryV2ParamsSelectReferenceDatasetID, RunQueryV2ParamsSelectTotalTokens, RunQueryV2ParamsSelectPromptTokens, RunQueryV2ParamsSelectCompletionTokens, RunQueryV2ParamsSelectTotalCost, RunQueryV2ParamsSelectPromptCost, RunQueryV2ParamsSelectCompletionCost, RunQueryV2ParamsSelectPromptTokenDetails, RunQueryV2ParamsSelectCompletionTokenDetails, RunQueryV2ParamsSelectPromptCostDetails, RunQueryV2ParamsSelectCompletionCostDetails, RunQueryV2ParamsSelectPriceModelID, RunQueryV2ParamsSelectTags, RunQueryV2ParamsSelectAppPath, RunQueryV2ParamsSelectAttachments, RunQueryV2ParamsSelectThreadEvaluationTime, RunQueryV2ParamsSelectIsInDataset, RunQueryV2ParamsSelectShareURL, RunQueryV2ParamsSelectFeedbackStats:
		return true
	}
	return false
}

type RunGetV1Params struct {
	ExcludeS3StoredAttributes param.Field[bool]      `query:"exclude_s3_stored_attributes"`
	ExcludeSerialized         param.Field[bool]      `query:"exclude_serialized"`
	IncludeMessages           param.Field[bool]      `query:"include_messages"`
	SessionID                 param.Field[string]    `query:"session_id" format:"uuid"`
	StartTime                 param.Field[time.Time] `query:"start_time" format:"date-time"`
}

// URLQuery serializes [RunGetV1Params]'s query parameters as `url.Values`.
func (r RunGetV1Params) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type RunGetV2Params struct {
	// `project_id` is the UUID of the tracing project that owns the run.
	ProjectID param.Field[string] `query:"project_id" api:"required" format:"uuid"`
	// `start_time` is the run's `start_time` (RFC3339 date-time), used together with
	// `project_id` to locate the run.
	StartTime param.Field[time.Time] `query:"start_time" api:"required" format:"date-time"`
	// `selects` lists which properties to include on the returned run (repeatable
	// query parameter). Accepts any value of the `RunSelectField` enum. If omitted,
	// only `id` is returned.
	Selects param.Field[[]RunGetV2ParamsSelect] `query:"selects"`
	Accept  param.Field[string]                 `header:"Accept"`
}

// URLQuery serializes [RunGetV2Params]'s query parameters as `url.Values`.
func (r RunGetV2Params) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type RunGetV2ParamsSelect string

const (
	RunGetV2ParamsSelectID                     RunGetV2ParamsSelect = "ID"
	RunGetV2ParamsSelectName                   RunGetV2ParamsSelect = "NAME"
	RunGetV2ParamsSelectRunType                RunGetV2ParamsSelect = "RUN_TYPE"
	RunGetV2ParamsSelectStatus                 RunGetV2ParamsSelect = "STATUS"
	RunGetV2ParamsSelectStartTime              RunGetV2ParamsSelect = "START_TIME"
	RunGetV2ParamsSelectEndTime                RunGetV2ParamsSelect = "END_TIME"
	RunGetV2ParamsSelectLatencySeconds         RunGetV2ParamsSelect = "LATENCY_SECONDS"
	RunGetV2ParamsSelectFirstTokenTime         RunGetV2ParamsSelect = "FIRST_TOKEN_TIME"
	RunGetV2ParamsSelectError                  RunGetV2ParamsSelect = "ERROR"
	RunGetV2ParamsSelectErrorPreview           RunGetV2ParamsSelect = "ERROR_PREVIEW"
	RunGetV2ParamsSelectExtra                  RunGetV2ParamsSelect = "EXTRA"
	RunGetV2ParamsSelectMetadata               RunGetV2ParamsSelect = "METADATA"
	RunGetV2ParamsSelectEvents                 RunGetV2ParamsSelect = "EVENTS"
	RunGetV2ParamsSelectInputs                 RunGetV2ParamsSelect = "INPUTS"
	RunGetV2ParamsSelectInputsPreview          RunGetV2ParamsSelect = "INPUTS_PREVIEW"
	RunGetV2ParamsSelectOutputs                RunGetV2ParamsSelect = "OUTPUTS"
	RunGetV2ParamsSelectOutputsPreview         RunGetV2ParamsSelect = "OUTPUTS_PREVIEW"
	RunGetV2ParamsSelectManifest               RunGetV2ParamsSelect = "MANIFEST"
	RunGetV2ParamsSelectParentRunIDs           RunGetV2ParamsSelect = "PARENT_RUN_IDS"
	RunGetV2ParamsSelectProjectID              RunGetV2ParamsSelect = "PROJECT_ID"
	RunGetV2ParamsSelectTraceID                RunGetV2ParamsSelect = "TRACE_ID"
	RunGetV2ParamsSelectThreadID               RunGetV2ParamsSelect = "THREAD_ID"
	RunGetV2ParamsSelectDottedOrder            RunGetV2ParamsSelect = "DOTTED_ORDER"
	RunGetV2ParamsSelectIsRoot                 RunGetV2ParamsSelect = "IS_ROOT"
	RunGetV2ParamsSelectReferenceExampleID     RunGetV2ParamsSelect = "REFERENCE_EXAMPLE_ID"
	RunGetV2ParamsSelectReferenceDatasetID     RunGetV2ParamsSelect = "REFERENCE_DATASET_ID"
	RunGetV2ParamsSelectTotalTokens            RunGetV2ParamsSelect = "TOTAL_TOKENS"
	RunGetV2ParamsSelectPromptTokens           RunGetV2ParamsSelect = "PROMPT_TOKENS"
	RunGetV2ParamsSelectCompletionTokens       RunGetV2ParamsSelect = "COMPLETION_TOKENS"
	RunGetV2ParamsSelectTotalCost              RunGetV2ParamsSelect = "TOTAL_COST"
	RunGetV2ParamsSelectPromptCost             RunGetV2ParamsSelect = "PROMPT_COST"
	RunGetV2ParamsSelectCompletionCost         RunGetV2ParamsSelect = "COMPLETION_COST"
	RunGetV2ParamsSelectPromptTokenDetails     RunGetV2ParamsSelect = "PROMPT_TOKEN_DETAILS"
	RunGetV2ParamsSelectCompletionTokenDetails RunGetV2ParamsSelect = "COMPLETION_TOKEN_DETAILS"
	RunGetV2ParamsSelectPromptCostDetails      RunGetV2ParamsSelect = "PROMPT_COST_DETAILS"
	RunGetV2ParamsSelectCompletionCostDetails  RunGetV2ParamsSelect = "COMPLETION_COST_DETAILS"
	RunGetV2ParamsSelectPriceModelID           RunGetV2ParamsSelect = "PRICE_MODEL_ID"
	RunGetV2ParamsSelectTags                   RunGetV2ParamsSelect = "TAGS"
	RunGetV2ParamsSelectAppPath                RunGetV2ParamsSelect = "APP_PATH"
	RunGetV2ParamsSelectAttachments            RunGetV2ParamsSelect = "ATTACHMENTS"
	RunGetV2ParamsSelectThreadEvaluationTime   RunGetV2ParamsSelect = "THREAD_EVALUATION_TIME"
	RunGetV2ParamsSelectIsInDataset            RunGetV2ParamsSelect = "IS_IN_DATASET"
	RunGetV2ParamsSelectShareURL               RunGetV2ParamsSelect = "SHARE_URL"
	RunGetV2ParamsSelectFeedbackStats          RunGetV2ParamsSelect = "FEEDBACK_STATS"
)

func (r RunGetV2ParamsSelect) IsKnown() bool {
	switch r {
	case RunGetV2ParamsSelectID, RunGetV2ParamsSelectName, RunGetV2ParamsSelectRunType, RunGetV2ParamsSelectStatus, RunGetV2ParamsSelectStartTime, RunGetV2ParamsSelectEndTime, RunGetV2ParamsSelectLatencySeconds, RunGetV2ParamsSelectFirstTokenTime, RunGetV2ParamsSelectError, RunGetV2ParamsSelectErrorPreview, RunGetV2ParamsSelectExtra, RunGetV2ParamsSelectMetadata, RunGetV2ParamsSelectEvents, RunGetV2ParamsSelectInputs, RunGetV2ParamsSelectInputsPreview, RunGetV2ParamsSelectOutputs, RunGetV2ParamsSelectOutputsPreview, RunGetV2ParamsSelectManifest, RunGetV2ParamsSelectParentRunIDs, RunGetV2ParamsSelectProjectID, RunGetV2ParamsSelectTraceID, RunGetV2ParamsSelectThreadID, RunGetV2ParamsSelectDottedOrder, RunGetV2ParamsSelectIsRoot, RunGetV2ParamsSelectReferenceExampleID, RunGetV2ParamsSelectReferenceDatasetID, RunGetV2ParamsSelectTotalTokens, RunGetV2ParamsSelectPromptTokens, RunGetV2ParamsSelectCompletionTokens, RunGetV2ParamsSelectTotalCost, RunGetV2ParamsSelectPromptCost, RunGetV2ParamsSelectCompletionCost, RunGetV2ParamsSelectPromptTokenDetails, RunGetV2ParamsSelectCompletionTokenDetails, RunGetV2ParamsSelectPromptCostDetails, RunGetV2ParamsSelectCompletionCostDetails, RunGetV2ParamsSelectPriceModelID, RunGetV2ParamsSelectTags, RunGetV2ParamsSelectAppPath, RunGetV2ParamsSelectAttachments, RunGetV2ParamsSelectThreadEvaluationTime, RunGetV2ParamsSelectIsInDataset, RunGetV2ParamsSelectShareURL, RunGetV2ParamsSelectFeedbackStats:
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
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
