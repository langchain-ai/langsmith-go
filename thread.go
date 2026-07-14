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
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// ThreadService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewThreadService] method instead.
type ThreadService struct {
	Options []option.RequestOption
}

// NewThreadService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewThreadService(opts ...option.RequestOption) (r *ThreadService) {
	r = &ThreadService{}
	r.Options = opts
	return
}

// **Alpha:** The request and response contract may change; Retrieve all traces
// belonging to a specific thread within a project.
func (r *ThreadService) ListTraces(ctx context.Context, threadID string, query ThreadListTracesParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[ThreadTrace], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/threads/%s/traces", threadID)
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

// **Alpha:** The request and response contract may change; Retrieve all traces
// belonging to a specific thread within a project.
func (r *ThreadService) ListTracesAutoPaging(ctx context.Context, threadID string, query ThreadListTracesParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[ThreadTrace] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.ListTraces(ctx, threadID, query, opts...))
}

// **Alpha:** The request and response contract may change; Query threads within a
// project (session), with cursor-based pagination. Returns threads matching the
// given time range and optional filter.
func (r *ThreadService) Query(ctx context.Context, body ThreadQueryParams, opts ...option.RequestOption) (res *pagination.ItemsCursorPostPagination[Thread], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v2/threads/query"
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

// **Alpha:** The request and response contract may change; Query threads within a
// project (session), with cursor-based pagination. Returns threads matching the
// given time range and optional filter.
func (r *ThreadService) QueryAutoPaging(ctx context.Context, body ThreadQueryParams, opts ...option.RequestOption) *pagination.ItemsCursorPostPaginationAutoPager[Thread] {
	return pagination.NewItemsCursorPostPaginationAutoPager(r.Query(ctx, body, opts...))
}

// **Alpha:** The request and response contract may change; Compute aggregate stats
// for a single thread (turn count, latency percentiles, token/cost sums, and
// detail breakdowns) within a project.
func (r *ThreadService) Stats(ctx context.Context, threadID string, query ThreadStatsParams, opts ...option.RequestOption) (res *ThreadStats, err error) {
	opts = slices.Concat(r.Options, opts)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/threads/%s/stats", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type Thread struct {
	// `count` is how many root traces (conversation turns) fall in this thread for the
	// query time range.
	Count int64 `json:"count"`
	// `feedback_stats` is the aggregated feedback across traces in the thread, keyed
	// by feedback key; shape matches `feedback_stats` on a single run.
	FeedbackStats map[string]ThreadFeedbackStat `json:"feedback_stats"`
	// `first_inputs` is a truncated preview of inputs from the earliest trace in the
	// thread for the query window.
	FirstInputs string `json:"first_inputs"`
	// `first_trace_id` is the root trace UUID for the chronologically first trace in
	// the query time window.
	FirstTraceID string `json:"first_trace_id" format:"uuid"`
	// `last_error` is a short error summary from the most recent failing trace in the
	// thread. Absent when there is no error in the window.
	LastError string `json:"last_error"`
	// `last_outputs` is a truncated preview of outputs from the latest trace in the
	// thread for the query window.
	LastOutputs string `json:"last_outputs"`
	// `last_trace_id` is the root trace UUID for the chronologically last trace in the
	// query time window.
	LastTraceID string `json:"last_trace_id" format:"uuid"`
	// `latency_p50` is the approximate median end-to-end latency of traces in the
	// thread, in seconds.
	LatencyP50 float64 `json:"latency_p50"`
	// `latency_p99` is the approximate 99th percentile end-to-end latency of traces in
	// the thread, in seconds.
	LatencyP99 float64 `json:"latency_p99"`
	// `max_start_time` is the latest trace start time in the thread (RFC3339
	// date-time).
	MaxStartTime time.Time `json:"max_start_time" format:"date-time"`
	// `min_start_time` is the earliest trace start time in the thread (RFC3339
	// date-time).
	MinStartTime time.Time `json:"min_start_time" format:"date-time"`
	// `num_errored_turns` is the count of root traces in the thread (within the query
	// window) whose status was an error.
	NumErroredTurns int64 `json:"num_errored_turns"`
	// `start_time` is a reference start time for this row (RFC3339 date-time), such as
	// for sorting.
	StartTime time.Time `json:"start_time" format:"date-time"`
	// `thread_id` identifies this conversation thread within the project from the
	// request body `project_id`.
	ThreadID string `json:"thread_id" format:"uuid"`
	// `total_cost` is the sum of estimated USD cost across those traces.
	TotalCost float64 `json:"total_cost"`
	// `total_cost_details` sums per-category estimated USD cost across traces in the
	// thread. Keys mirror `total_token_details`.
	//
	// Example: `{"cache_read": 0.012, "reasoning": 0.008}`.
	TotalCostDetails map[string]float64 `json:"total_cost_details"`
	// `total_token_details` sums per-category token counts across traces in the
	// thread. Keys are model-specific category names (for example `cache_read`,
	// `cache_write`, `reasoning`, `audio`).
	//
	// Example: `{"cache_read": 400, "reasoning": 120}`.
	TotalTokenDetails map[string]int64 `json:"total_token_details"`
	// `total_tokens` is the sum of token usage across those traces.
	TotalTokens int64 `json:"total_tokens"`
	// `trace_id` is a representative root trace UUID when the summary includes one,
	// for example for deep links.
	TraceID string     `json:"trace_id" format:"uuid"`
	JSON    threadJSON `json:"-"`
}

// threadJSON contains the JSON metadata for the struct [Thread]
type threadJSON struct {
	Count             apijson.Field
	FeedbackStats     apijson.Field
	FirstInputs       apijson.Field
	FirstTraceID      apijson.Field
	LastError         apijson.Field
	LastOutputs       apijson.Field
	LastTraceID       apijson.Field
	LatencyP50        apijson.Field
	LatencyP99        apijson.Field
	MaxStartTime      apijson.Field
	MinStartTime      apijson.Field
	NumErroredTurns   apijson.Field
	StartTime         apijson.Field
	ThreadID          apijson.Field
	TotalCost         apijson.Field
	TotalCostDetails  apijson.Field
	TotalTokenDetails apijson.Field
	TotalTokens       apijson.Field
	TraceID           apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *Thread) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadJSON) RawJSON() string {
	return r.raw
}

type ThreadFeedbackStat struct {
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
	Values map[string]int64       `json:"values"`
	JSON   threadFeedbackStatJSON `json:"-"`
}

// threadFeedbackStatJSON contains the JSON metadata for the struct
// [ThreadFeedbackStat]
type threadFeedbackStatJSON struct {
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

func (r *ThreadFeedbackStat) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadFeedbackStatJSON) RawJSON() string {
	return r.raw
}

type ThreadStats struct {
	// `completion_cost` is the sum of per-trace completion costs across the thread, in
	// USD. Populated when `COMPLETION_COST` is selected.
	CompletionCost float64 `json:"completion_cost"`
	// `completion_cost_details` is the per-sub-category sum of completion cost details
	// across the thread. Populated when `COMPLETION_COST_DETAILS` is selected.
	CompletionCostDetails ThreadStatsCompletionCostDetails `json:"completion_cost_details"`
	// `completion_token_details` is the per-sub-category sum of completion token
	// details across the thread. Populated when `COMPLETION_TOKEN_DETAILS` is
	// selected.
	CompletionTokenDetails ThreadStatsCompletionTokenDetails `json:"completion_token_details"`
	// `completion_tokens` is the sum of per-trace completion token counts across the
	// thread. Populated when `COMPLETION_TOKENS` is selected.
	CompletionTokens int64 `json:"completion_tokens"`
	// `feedback_stats` aggregates run-level feedback across the thread's traces, keyed
	// by feedback key. Populated when `FEEDBACK_STATS` is selected.
	FeedbackStats map[string]ThreadStatsFeedbackStat `json:"feedback_stats"`
	// `first_start_time` is the earliest trace start time in the thread (RFC3339).
	// Populated when `FIRST_START_TIME` is selected.
	FirstStartTime time.Time `json:"first_start_time" format:"date-time"`
	// `last_end_time` is the latest trace end time in the thread (RFC3339). Populated
	// when `LAST_END_TIME` is selected.
	LastEndTime time.Time `json:"last_end_time" format:"date-time"`
	// `last_start_time` is the latest trace start time in the thread (RFC3339).
	// Populated when `LAST_START_TIME` is selected.
	LastStartTime time.Time `json:"last_start_time" format:"date-time"`
	// `latency_p50_seconds` is the approximate p50 of trace latency across the thread,
	// in seconds. Populated when `LATENCY_P50` is selected.
	LatencyP50Seconds float64 `json:"latency_p50_seconds"`
	// `latency_p99_seconds` is the approximate p99 of trace latency across the thread,
	// in seconds. Populated when `LATENCY_P99` is selected.
	LatencyP99Seconds float64 `json:"latency_p99_seconds"`
	// `prompt_cost` is the sum of per-trace prompt costs across the thread, in USD.
	// Populated when `PROMPT_COST` is selected.
	PromptCost float64 `json:"prompt_cost"`
	// `prompt_cost_details` is the per-sub-category sum of prompt cost details across
	// the thread. Populated when `PROMPT_COST_DETAILS` is selected.
	PromptCostDetails ThreadStatsPromptCostDetails `json:"prompt_cost_details"`
	// `prompt_token_details` is the per-sub-category sum of prompt token details
	// across the thread. Populated when `PROMPT_TOKEN_DETAILS` is selected.
	PromptTokenDetails ThreadStatsPromptTokenDetails `json:"prompt_token_details"`
	// `prompt_tokens` is the sum of per-trace prompt token counts across the thread.
	// Populated when `PROMPT_TOKENS` is selected.
	PromptTokens int64 `json:"prompt_tokens"`
	// `total_cost` is the sum of per-trace total costs across the thread, in USD.
	// Populated when `TOTAL_COST` is selected.
	TotalCost float64 `json:"total_cost"`
	// `total_tokens` is the sum of per-trace total token counts across the thread.
	// Populated when `TOTAL_TOKENS` is selected.
	TotalTokens int64 `json:"total_tokens"`
	// `turns` is the number of distinct traces (turns) in the thread. Populated when
	// `TURNS` is selected.
	Turns int64           `json:"turns"`
	JSON  threadStatsJSON `json:"-"`
}

// threadStatsJSON contains the JSON metadata for the struct [ThreadStats]
type threadStatsJSON struct {
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	FeedbackStats          apijson.Field
	FirstStartTime         apijson.Field
	LastEndTime            apijson.Field
	LastStartTime          apijson.Field
	LatencyP50Seconds      apijson.Field
	LatencyP99Seconds      apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	Turns                  apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *ThreadStats) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStatsJSON) RawJSON() string {
	return r.raw
}

// `completion_cost_details` is the per-sub-category sum of completion cost details
// across the thread. Populated when `COMPLETION_COST_DETAILS` is selected.
type ThreadStatsCompletionCostDetails struct {
	// `raw` maps each category name to its estimated USD cost.
	Raw  map[string]float64                   `json:"raw"`
	JSON threadStatsCompletionCostDetailsJSON `json:"-"`
}

// threadStatsCompletionCostDetailsJSON contains the JSON metadata for the struct
// [ThreadStatsCompletionCostDetails]
type threadStatsCompletionCostDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadStatsCompletionCostDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStatsCompletionCostDetailsJSON) RawJSON() string {
	return r.raw
}

// `completion_token_details` is the per-sub-category sum of completion token
// details across the thread. Populated when `COMPLETION_TOKEN_DETAILS` is
// selected.
type ThreadStatsCompletionTokenDetails struct {
	// `raw` maps each category name to its completion-token count.
	Raw  map[string]int64                      `json:"raw"`
	JSON threadStatsCompletionTokenDetailsJSON `json:"-"`
}

// threadStatsCompletionTokenDetailsJSON contains the JSON metadata for the struct
// [ThreadStatsCompletionTokenDetails]
type threadStatsCompletionTokenDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadStatsCompletionTokenDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStatsCompletionTokenDetailsJSON) RawJSON() string {
	return r.raw
}

type ThreadStatsFeedbackStat struct {
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
	Values map[string]int64            `json:"values"`
	JSON   threadStatsFeedbackStatJSON `json:"-"`
}

// threadStatsFeedbackStatJSON contains the JSON metadata for the struct
// [ThreadStatsFeedbackStat]
type threadStatsFeedbackStatJSON struct {
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

func (r *ThreadStatsFeedbackStat) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStatsFeedbackStatJSON) RawJSON() string {
	return r.raw
}

// `prompt_cost_details` is the per-sub-category sum of prompt cost details across
// the thread. Populated when `PROMPT_COST_DETAILS` is selected.
type ThreadStatsPromptCostDetails struct {
	// `raw` maps each category name to its estimated USD cost.
	Raw  map[string]float64               `json:"raw"`
	JSON threadStatsPromptCostDetailsJSON `json:"-"`
}

// threadStatsPromptCostDetailsJSON contains the JSON metadata for the struct
// [ThreadStatsPromptCostDetails]
type threadStatsPromptCostDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadStatsPromptCostDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStatsPromptCostDetailsJSON) RawJSON() string {
	return r.raw
}

// `prompt_token_details` is the per-sub-category sum of prompt token details
// across the thread. Populated when `PROMPT_TOKEN_DETAILS` is selected.
type ThreadStatsPromptTokenDetails struct {
	// `raw` maps each category name to its prompt-token count.
	Raw  map[string]int64                  `json:"raw"`
	JSON threadStatsPromptTokenDetailsJSON `json:"-"`
}

// threadStatsPromptTokenDetailsJSON contains the JSON metadata for the struct
// [ThreadStatsPromptTokenDetails]
type threadStatsPromptTokenDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadStatsPromptTokenDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStatsPromptTokenDetailsJSON) RawJSON() string {
	return r.raw
}

type ThreadTrace struct {
	// `completion_cost` is the estimated USD cost for the completion. Omitted unless
	// included in `selects`.
	CompletionCost float64 `json:"completion_cost"`
	// `completion_cost_details` is the USD cost breakdown for completion-side
	// categories; per-category values are under `raw`. Omitted unless included in
	// `selects`.
	CompletionCostDetails ThreadTraceCompletionCostDetails `json:"completion_cost_details"`
	// `completion_token_details` is the completion-side token breakdown by category;
	// per-category counts are under `raw`. Omitted unless included in `selects`.
	CompletionTokenDetails ThreadTraceCompletionTokenDetails `json:"completion_token_details"`
	// `completion_tokens` is the completion-side token count. Omitted unless included
	// in `selects`.
	CompletionTokens int64 `json:"completion_tokens"`
	// `end_time` is when the root run ended (RFC3339 date-time). JSON null if the run
	// is still in progress. Omitted unless included in `selects`.
	EndTime time.Time `json:"end_time" format:"date-time"`
	// `error` is the full root run error message when the run failed. Omitted unless
	// included in `selects`.
	Error string `json:"error"`
	// `error_preview` is a short error summary when the run failed. Omitted unless
	// included in `selects`.
	ErrorPreview string `json:"error_preview"`
	// `first_token_time` is when the first output token was produced (RFC3339
	// date-time), for streamed runs when that metadata exists. Omitted unless included
	// in `selects`.
	FirstTokenTime time.Time `json:"first_token_time" format:"date-time"`
	// `inputs` is the full root run input payload. Omitted unless included in
	// `selects`.
	Inputs interface{} `json:"inputs"`
	// `inputs_preview` is a truncated text preview of inputs. Omitted unless included
	// in `selects`.
	InputsPreview string `json:"inputs_preview"`
	// `latency` is wall-clock duration from start to end in seconds. Omitted unless
	// included in `selects`.
	Latency float64 `json:"latency"`
	// `name` is a human-readable label for the root run (for example the model name,
	// function name, or step name chosen when the run was traced). Omitted unless
	// included in `selects`.
	Name string `json:"name"`
	// `op` is a numeric code identifying the root run's `run_type` (for example LLM
	// vs. tool vs. chain). Encoded as a number for compatibility with legacy clients;
	// prefer the string `run_type` on `RunResponse` when available. Omitted unless
	// included in `selects`.
	Op float64 `json:"op"`
	// `outputs` is the full root run output payload. Omitted unless included in
	// `selects`.
	Outputs interface{} `json:"outputs"`
	// `outputs_preview` is a truncated text preview of outputs. Omitted unless
	// included in `selects`.
	OutputsPreview string `json:"outputs_preview"`
	// `prompt_cost` is the estimated USD cost for the prompt. Omitted unless included
	// in `selects`.
	PromptCost float64 `json:"prompt_cost"`
	// `prompt_cost_details` is the USD cost breakdown for prompt-side categories;
	// per-category values are under `raw`. Omitted unless included in `selects`.
	PromptCostDetails ThreadTracePromptCostDetails `json:"prompt_cost_details"`
	// `prompt_token_details` is the prompt-side token breakdown by category;
	// per-category counts are under nested `raw`. Omitted unless included in
	// `selects`.
	PromptTokenDetails ThreadTracePromptTokenDetails `json:"prompt_token_details"`
	// `prompt_tokens` is the prompt-side token count. Omitted unless included in
	// `selects`.
	PromptTokens int64 `json:"prompt_tokens"`
	// `start_time` is when the trace started (RFC3339 date-time). Omitted unless
	// included in `selects`.
	StartTime time.Time `json:"start_time" format:"date-time"`
	// `thread_id` is the conversation thread UUID that contains this trace. Matches
	// the `thread_id` path parameter of the request. Omitted unless included in
	// `selects`.
	ThreadID string `json:"thread_id" format:"uuid"`
	// `total_cost` is the estimated total USD cost for the root run. Omitted unless
	// included in `selects`.
	TotalCost float64 `json:"total_cost"`
	// `total_tokens` is the total token count (prompt plus completion). Omitted unless
	// included in `selects`.
	TotalTokens int64 `json:"total_tokens"`
	// `trace_id` is the UUID of this trace (the root run). Always present.
	TraceID string          `json:"trace_id" format:"uuid"`
	JSON    threadTraceJSON `json:"-"`
}

// threadTraceJSON contains the JSON metadata for the struct [ThreadTrace]
type threadTraceJSON struct {
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	EndTime                apijson.Field
	Error                  apijson.Field
	ErrorPreview           apijson.Field
	FirstTokenTime         apijson.Field
	Inputs                 apijson.Field
	InputsPreview          apijson.Field
	Latency                apijson.Field
	Name                   apijson.Field
	Op                     apijson.Field
	Outputs                apijson.Field
	OutputsPreview         apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	StartTime              apijson.Field
	ThreadID               apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	TraceID                apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *ThreadTrace) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadTraceJSON) RawJSON() string {
	return r.raw
}

// `completion_cost_details` is the USD cost breakdown for completion-side
// categories; per-category values are under `raw`. Omitted unless included in
// `selects`.
type ThreadTraceCompletionCostDetails struct {
	// `raw` maps each category name to its estimated USD cost.
	Raw  map[string]float64                   `json:"raw"`
	JSON threadTraceCompletionCostDetailsJSON `json:"-"`
}

// threadTraceCompletionCostDetailsJSON contains the JSON metadata for the struct
// [ThreadTraceCompletionCostDetails]
type threadTraceCompletionCostDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadTraceCompletionCostDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadTraceCompletionCostDetailsJSON) RawJSON() string {
	return r.raw
}

// `completion_token_details` is the completion-side token breakdown by category;
// per-category counts are under `raw`. Omitted unless included in `selects`.
type ThreadTraceCompletionTokenDetails struct {
	// `raw` maps each category name to its completion-token count.
	Raw  map[string]int64                      `json:"raw"`
	JSON threadTraceCompletionTokenDetailsJSON `json:"-"`
}

// threadTraceCompletionTokenDetailsJSON contains the JSON metadata for the struct
// [ThreadTraceCompletionTokenDetails]
type threadTraceCompletionTokenDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadTraceCompletionTokenDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadTraceCompletionTokenDetailsJSON) RawJSON() string {
	return r.raw
}

// `prompt_cost_details` is the USD cost breakdown for prompt-side categories;
// per-category values are under `raw`. Omitted unless included in `selects`.
type ThreadTracePromptCostDetails struct {
	// `raw` maps each category name to its estimated USD cost.
	Raw  map[string]float64               `json:"raw"`
	JSON threadTracePromptCostDetailsJSON `json:"-"`
}

// threadTracePromptCostDetailsJSON contains the JSON metadata for the struct
// [ThreadTracePromptCostDetails]
type threadTracePromptCostDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadTracePromptCostDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadTracePromptCostDetailsJSON) RawJSON() string {
	return r.raw
}

// `prompt_token_details` is the prompt-side token breakdown by category;
// per-category counts are under nested `raw`. Omitted unless included in
// `selects`.
type ThreadTracePromptTokenDetails struct {
	// `raw` maps each category name to its prompt-token count.
	Raw  map[string]int64                  `json:"raw"`
	JSON threadTracePromptTokenDetailsJSON `json:"-"`
}

// threadTracePromptTokenDetailsJSON contains the JSON metadata for the struct
// [ThreadTracePromptTokenDetails]
type threadTracePromptTokenDetailsJSON struct {
	Raw         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadTracePromptTokenDetails) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadTracePromptTokenDetailsJSON) RawJSON() string {
	return r.raw
}

type ThreadListTracesParams struct {
	// `project_id` is the tracing project UUID (required).
	ProjectID param.Field[string] `query:"project_id" api:"required" format:"uuid"`
	// `cursor` is the opaque string from a previous response's `next_cursor`. Omit on
	// the first request; pass the returned cursor to fetch the next page.
	Cursor param.Field[string] `query:"cursor"`
	// `filter` narrows which traces are returned for this thread, using a LangSmith
	// filter expression evaluated against each root trace run. For example: eq(status,
	// "success") or has(tags, "production"). See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	Filter param.Field[string] `query:"filter"`
	// `page_size` is the maximum number of traces to return in this response. Defaults
	// to 20 when omitted; must be between 1 and 100 inclusive when set.
	PageSize param.Field[int64] `query:"page_size"`
	// `selects` lists which properties to include on each returned trace (repeatable
	// query parameter). Accepts any value of the `ThreadTraceSelectField` enum.
	// Properties not listed are omitted from each trace object; `trace_id` is always
	// returned.
	Selects param.Field[[]ThreadListTracesParamsSelect] `query:"selects"`
}

// URLQuery serializes [ThreadListTracesParams]'s query parameters as `url.Values`.
func (r ThreadListTracesParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ThreadListTracesParamsSelect string

const (
	ThreadListTracesParamsSelectThreadID               ThreadListTracesParamsSelect = "THREAD_ID"
	ThreadListTracesParamsSelectTraceID                ThreadListTracesParamsSelect = "TRACE_ID"
	ThreadListTracesParamsSelectOp                     ThreadListTracesParamsSelect = "OP"
	ThreadListTracesParamsSelectPromptTokens           ThreadListTracesParamsSelect = "PROMPT_TOKENS"
	ThreadListTracesParamsSelectCompletionTokens       ThreadListTracesParamsSelect = "COMPLETION_TOKENS"
	ThreadListTracesParamsSelectTotalTokens            ThreadListTracesParamsSelect = "TOTAL_TOKENS"
	ThreadListTracesParamsSelectStartTime              ThreadListTracesParamsSelect = "START_TIME"
	ThreadListTracesParamsSelectEndTime                ThreadListTracesParamsSelect = "END_TIME"
	ThreadListTracesParamsSelectLatency                ThreadListTracesParamsSelect = "LATENCY"
	ThreadListTracesParamsSelectFirstTokenTime         ThreadListTracesParamsSelect = "FIRST_TOKEN_TIME"
	ThreadListTracesParamsSelectInputsPreview          ThreadListTracesParamsSelect = "INPUTS_PREVIEW"
	ThreadListTracesParamsSelectOutputsPreview         ThreadListTracesParamsSelect = "OUTPUTS_PREVIEW"
	ThreadListTracesParamsSelectInputs                 ThreadListTracesParamsSelect = "INPUTS"
	ThreadListTracesParamsSelectOutputs                ThreadListTracesParamsSelect = "OUTPUTS"
	ThreadListTracesParamsSelectError                  ThreadListTracesParamsSelect = "ERROR"
	ThreadListTracesParamsSelectPromptCost             ThreadListTracesParamsSelect = "PROMPT_COST"
	ThreadListTracesParamsSelectCompletionCost         ThreadListTracesParamsSelect = "COMPLETION_COST"
	ThreadListTracesParamsSelectTotalCost              ThreadListTracesParamsSelect = "TOTAL_COST"
	ThreadListTracesParamsSelectPromptTokenDetails     ThreadListTracesParamsSelect = "PROMPT_TOKEN_DETAILS"
	ThreadListTracesParamsSelectCompletionTokenDetails ThreadListTracesParamsSelect = "COMPLETION_TOKEN_DETAILS"
	ThreadListTracesParamsSelectPromptCostDetails      ThreadListTracesParamsSelect = "PROMPT_COST_DETAILS"
	ThreadListTracesParamsSelectCompletionCostDetails  ThreadListTracesParamsSelect = "COMPLETION_COST_DETAILS"
	ThreadListTracesParamsSelectName                   ThreadListTracesParamsSelect = "NAME"
	ThreadListTracesParamsSelectErrorPreview           ThreadListTracesParamsSelect = "ERROR_PREVIEW"
)

func (r ThreadListTracesParamsSelect) IsKnown() bool {
	switch r {
	case ThreadListTracesParamsSelectThreadID, ThreadListTracesParamsSelectTraceID, ThreadListTracesParamsSelectOp, ThreadListTracesParamsSelectPromptTokens, ThreadListTracesParamsSelectCompletionTokens, ThreadListTracesParamsSelectTotalTokens, ThreadListTracesParamsSelectStartTime, ThreadListTracesParamsSelectEndTime, ThreadListTracesParamsSelectLatency, ThreadListTracesParamsSelectFirstTokenTime, ThreadListTracesParamsSelectInputsPreview, ThreadListTracesParamsSelectOutputsPreview, ThreadListTracesParamsSelectInputs, ThreadListTracesParamsSelectOutputs, ThreadListTracesParamsSelectError, ThreadListTracesParamsSelectPromptCost, ThreadListTracesParamsSelectCompletionCost, ThreadListTracesParamsSelectTotalCost, ThreadListTracesParamsSelectPromptTokenDetails, ThreadListTracesParamsSelectCompletionTokenDetails, ThreadListTracesParamsSelectPromptCostDetails, ThreadListTracesParamsSelectCompletionCostDetails, ThreadListTracesParamsSelectName, ThreadListTracesParamsSelectErrorPreview:
		return true
	}
	return false
}

type ThreadQueryParams struct {
	// `cursor` is the opaque string from a previous response's `next_cursor`. Omit on
	// the first request; pass the returned cursor to fetch the next page.
	Cursor param.Field[string] `json:"cursor"`
	// `filter` narrows which threads are returned, using a LangSmith filter expression
	// evaluated against each thread's root run. For example: has(tags, "production")
	// or eq(status, "error"). See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	Filter param.Field[string] `json:"filter"`
	// `max_start_time` is the exclusive upper bound on thread activity (RFC3339
	// date-time). Defaults to now (UTC) when omitted.
	MaxStartTime param.Field[time.Time] `json:"max_start_time" format:"date-time"`
	// `min_start_time` is the inclusive lower bound on thread activity (RFC3339
	// date-time). Defaults to 1 day before now (UTC) when omitted.
	MinStartTime param.Field[time.Time] `json:"min_start_time" format:"date-time"`
	// `page_size` is the maximum number of threads to return in this response.
	// Defaults to 20 when omitted; must be between 1 and 100 inclusive when set. The
	// response may contain fewer threads than `page_size` even when `next_cursor` is
	// non-null.
	PageSize param.Field[int64] `json:"page_size"`
	// `project_id` is the tracing project UUID.
	ProjectID param.Field[string] `json:"project_id" format:"uuid"`
}

func (r ThreadQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ThreadStatsParams struct {
	// `selects` lists which aggregate stats to compute and return (repeatable query
	// parameter). At least one value is required. Accepts any value of
	// `SingleThreadStatsSelectField`.
	Selects param.Field[[]ThreadStatsParamsSelect] `query:"selects" api:"required"`
	// `session_id` is the tracing project (session) UUID (required).
	SessionID param.Field[string] `query:"session_id" api:"required" format:"uuid"`
	// `filter` narrows which of the thread's traces are aggregated, using a LangSmith
	// filter expression. For example: lt(start_time, "2025-01-01T00:00:00Z") or
	// eq(trace_id, "0190a1b2-c3d4-7ef0-a5b6-6ea3a82e9328"). See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	Filter param.Field[string] `query:"filter"`
}

// URLQuery serializes [ThreadStatsParams]'s query parameters as `url.Values`.
func (r ThreadStatsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ThreadStatsParamsSelect string

const (
	ThreadStatsParamsSelectTurns                  ThreadStatsParamsSelect = "TURNS"
	ThreadStatsParamsSelectFirstStartTime         ThreadStatsParamsSelect = "FIRST_START_TIME"
	ThreadStatsParamsSelectLastStartTime          ThreadStatsParamsSelect = "LAST_START_TIME"
	ThreadStatsParamsSelectLastEndTime            ThreadStatsParamsSelect = "LAST_END_TIME"
	ThreadStatsParamsSelectLatencyP50             ThreadStatsParamsSelect = "LATENCY_P50"
	ThreadStatsParamsSelectLatencyP99             ThreadStatsParamsSelect = "LATENCY_P99"
	ThreadStatsParamsSelectPromptTokens           ThreadStatsParamsSelect = "PROMPT_TOKENS"
	ThreadStatsParamsSelectPromptCost             ThreadStatsParamsSelect = "PROMPT_COST"
	ThreadStatsParamsSelectCompletionTokens       ThreadStatsParamsSelect = "COMPLETION_TOKENS"
	ThreadStatsParamsSelectCompletionCost         ThreadStatsParamsSelect = "COMPLETION_COST"
	ThreadStatsParamsSelectTotalTokens            ThreadStatsParamsSelect = "TOTAL_TOKENS"
	ThreadStatsParamsSelectTotalCost              ThreadStatsParamsSelect = "TOTAL_COST"
	ThreadStatsParamsSelectPromptTokenDetails     ThreadStatsParamsSelect = "PROMPT_TOKEN_DETAILS"
	ThreadStatsParamsSelectCompletionTokenDetails ThreadStatsParamsSelect = "COMPLETION_TOKEN_DETAILS"
	ThreadStatsParamsSelectPromptCostDetails      ThreadStatsParamsSelect = "PROMPT_COST_DETAILS"
	ThreadStatsParamsSelectCompletionCostDetails  ThreadStatsParamsSelect = "COMPLETION_COST_DETAILS"
	ThreadStatsParamsSelectFeedbackStats          ThreadStatsParamsSelect = "FEEDBACK_STATS"
)

func (r ThreadStatsParamsSelect) IsKnown() bool {
	switch r {
	case ThreadStatsParamsSelectTurns, ThreadStatsParamsSelectFirstStartTime, ThreadStatsParamsSelectLastStartTime, ThreadStatsParamsSelectLastEndTime, ThreadStatsParamsSelectLatencyP50, ThreadStatsParamsSelectLatencyP99, ThreadStatsParamsSelectPromptTokens, ThreadStatsParamsSelectPromptCost, ThreadStatsParamsSelectCompletionTokens, ThreadStatsParamsSelectCompletionCost, ThreadStatsParamsSelectTotalTokens, ThreadStatsParamsSelectTotalCost, ThreadStatsParamsSelectPromptTokenDetails, ThreadStatsParamsSelectCompletionTokenDetails, ThreadStatsParamsSelectPromptCostDetails, ThreadStatsParamsSelectCompletionCostDetails, ThreadStatsParamsSelectFeedbackStats:
		return true
	}
	return false
}
