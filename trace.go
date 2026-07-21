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

// TraceService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTraceService] method instead.
type TraceService struct {
	Options []option.RequestOption
}

// NewTraceService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewTraceService(opts ...option.RequestOption) (r *TraceService) {
	r = &TraceService{}
	r.Options = opts
	return
}

// **Alpha:** The request and response contract may change; Returns runs for a
// trace ID within min/max start time. Optional `filter`; repeatable `selects` to
// select fields to return.
func (r *TraceService) ListRuns(ctx context.Context, traceID string, params TraceListRunsParams, opts ...option.RequestOption) (res *TraceListRunsResponse, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if traceID == "" {
		err = errors.New("missing required trace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/traces/%s/runs", traceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Returns a paginated list of traces (root runs) for a single tracing project.
// Each item carries the trace's root run plus optional trace-wide aggregates
// (`total_tokens`, `total_cost`, `first_token_time`) under `trace_aggregates`, so
// clients never have to merge by `trace_id`.
//
// Traces are scanned within a `start_time` window: `min_start_time` defaults to 24
// hours before the request, `max_start_time` defaults to the request time. Set
// either explicitly to widen or narrow the window.
//
// Supports filters (`trace_filter`, `tree_filter`), cursor pagination (`cursor`),
// and field projection (`selects`).
func (r *TraceService) Query(ctx context.Context, body TraceQueryParams, opts ...option.RequestOption) (res *pagination.ItemsCursorPostPagination[Trace], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v2/traces/query"
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

// Returns a paginated list of traces (root runs) for a single tracing project.
// Each item carries the trace's root run plus optional trace-wide aggregates
// (`total_tokens`, `total_cost`, `first_token_time`) under `trace_aggregates`, so
// clients never have to merge by `trace_id`.
//
// Traces are scanned within a `start_time` window: `min_start_time` defaults to 24
// hours before the request, `max_start_time` defaults to the request time. Set
// either explicitly to widen or narrow the window.
//
// Supports filters (`trace_filter`, `tree_filter`), cursor pagination (`cursor`),
// and field projection (`selects`).
func (r *TraceService) QueryAutoPaging(ctx context.Context, body TraceQueryParams, opts ...option.RequestOption) *pagination.ItemsCursorPostPaginationAutoPager[Trace] {
	return pagination.NewItemsCursorPostPaginationAutoPager(r.Query(ctx, body, opts...))
}

type Trace struct {
	// `root_run` is the trace's root run. Which properties are populated is controlled
	// by `selects` in the request.
	RootRun Run `json:"root_run"`
	// `trace_aggregates` carries trace-wide aggregate metrics. Omitted when no
	// aggregate field was selected, or `null` (then later filled) on the streaming
	// wire while the aggregate values are still being computed.
	TraceAggregates TraceAggregates `json:"trace_aggregates"`
	JSON            traceJSON       `json:"-"`
}

// traceJSON contains the JSON metadata for the struct [Trace]
type traceJSON struct {
	RootRun         apijson.Field
	TraceAggregates apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *Trace) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r traceJSON) RawJSON() string {
	return r.raw
}

type TraceAggregates struct {
	// `first_token_time` is when the first output token was produced anywhere in the
	// trace (RFC3339), when recorded.
	FirstTokenTime time.Time `json:"first_token_time" format:"date-time"`
	// `total_cost` is total estimated USD cost across every run in the trace.
	TotalCost float64 `json:"total_cost"`
	// `total_tokens` is prompt plus completion tokens summed across every run in the
	// trace.
	TotalTokens int64               `json:"total_tokens"`
	JSON        traceAggregatesJSON `json:"-"`
}

// traceAggregatesJSON contains the JSON metadata for the struct [TraceAggregates]
type traceAggregatesJSON struct {
	FirstTokenTime apijson.Field
	TotalCost      apijson.Field
	TotalTokens    apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *TraceAggregates) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r traceAggregatesJSON) RawJSON() string {
	return r.raw
}

type TraceListRunsResponse struct {
	// `items` lists runs in the trace for the requested time window, in `start_time`
	// order.
	Items []Run                     `json:"items"`
	JSON  traceListRunsResponseJSON `json:"-"`
}

// traceListRunsResponseJSON contains the JSON metadata for the struct
// [TraceListRunsResponse]
type traceListRunsResponseJSON struct {
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TraceListRunsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r traceListRunsResponseJSON) RawJSON() string {
	return r.raw
}

type TraceListRunsParams struct {
	// `project_id` is the UUID of the tracing project that owns the trace.
	ProjectID param.Field[string] `query:"project_id" api:"required" format:"uuid"`
	// `filter` narrows which runs within this trace are returned, using a LangSmith
	// filter expression evaluated against each run. For example: `eq(run_type, "llm")`
	// for LLM runs only, or `eq(status, "error")` for failed runs. See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	Filter param.Field[string] `query:"filter"`
	// `max_start_time` is the optional inclusive upper bound for run `start_time`
	// (RFC3339 date-time). Required together with `min_start_time`.
	MaxStartTime param.Field[time.Time] `query:"max_start_time" format:"date-time"`
	// `min_start_time` is the optional inclusive lower bound for run `start_time`
	// (RFC3339 date-time). Required together with `max_start_time`.
	MinStartTime param.Field[time.Time] `query:"min_start_time" format:"date-time"`
	// `selects` lists which properties to include on each returned run (repeatable
	// query parameter). Accepts any value of the `RunSelectField` enum. If omitted,
	// only `id` is returned.
	Selects param.Field[[]TraceListRunsParamsSelect] `query:"selects"`
	Accept  param.Field[string]                      `header:"Accept"`
}

// URLQuery serializes [TraceListRunsParams]'s query parameters as `url.Values`.
func (r TraceListRunsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type TraceListRunsParamsSelect string

const (
	TraceListRunsParamsSelectID                     TraceListRunsParamsSelect = "ID"
	TraceListRunsParamsSelectName                   TraceListRunsParamsSelect = "NAME"
	TraceListRunsParamsSelectRunType                TraceListRunsParamsSelect = "RUN_TYPE"
	TraceListRunsParamsSelectStatus                 TraceListRunsParamsSelect = "STATUS"
	TraceListRunsParamsSelectStartTime              TraceListRunsParamsSelect = "START_TIME"
	TraceListRunsParamsSelectEndTime                TraceListRunsParamsSelect = "END_TIME"
	TraceListRunsParamsSelectLatencySeconds         TraceListRunsParamsSelect = "LATENCY_SECONDS"
	TraceListRunsParamsSelectFirstTokenTime         TraceListRunsParamsSelect = "FIRST_TOKEN_TIME"
	TraceListRunsParamsSelectError                  TraceListRunsParamsSelect = "ERROR"
	TraceListRunsParamsSelectErrorPreview           TraceListRunsParamsSelect = "ERROR_PREVIEW"
	TraceListRunsParamsSelectExtra                  TraceListRunsParamsSelect = "EXTRA"
	TraceListRunsParamsSelectMetadata               TraceListRunsParamsSelect = "METADATA"
	TraceListRunsParamsSelectEvents                 TraceListRunsParamsSelect = "EVENTS"
	TraceListRunsParamsSelectInputs                 TraceListRunsParamsSelect = "INPUTS"
	TraceListRunsParamsSelectInputsPreview          TraceListRunsParamsSelect = "INPUTS_PREVIEW"
	TraceListRunsParamsSelectOutputs                TraceListRunsParamsSelect = "OUTPUTS"
	TraceListRunsParamsSelectOutputsPreview         TraceListRunsParamsSelect = "OUTPUTS_PREVIEW"
	TraceListRunsParamsSelectManifest               TraceListRunsParamsSelect = "MANIFEST"
	TraceListRunsParamsSelectParentRunIDs           TraceListRunsParamsSelect = "PARENT_RUN_IDS"
	TraceListRunsParamsSelectProjectID              TraceListRunsParamsSelect = "PROJECT_ID"
	TraceListRunsParamsSelectTraceID                TraceListRunsParamsSelect = "TRACE_ID"
	TraceListRunsParamsSelectThreadID               TraceListRunsParamsSelect = "THREAD_ID"
	TraceListRunsParamsSelectDottedOrder            TraceListRunsParamsSelect = "DOTTED_ORDER"
	TraceListRunsParamsSelectIsRoot                 TraceListRunsParamsSelect = "IS_ROOT"
	TraceListRunsParamsSelectReferenceExampleID     TraceListRunsParamsSelect = "REFERENCE_EXAMPLE_ID"
	TraceListRunsParamsSelectReferenceDatasetID     TraceListRunsParamsSelect = "REFERENCE_DATASET_ID"
	TraceListRunsParamsSelectTotalTokens            TraceListRunsParamsSelect = "TOTAL_TOKENS"
	TraceListRunsParamsSelectPromptTokens           TraceListRunsParamsSelect = "PROMPT_TOKENS"
	TraceListRunsParamsSelectCompletionTokens       TraceListRunsParamsSelect = "COMPLETION_TOKENS"
	TraceListRunsParamsSelectTotalCost              TraceListRunsParamsSelect = "TOTAL_COST"
	TraceListRunsParamsSelectPromptCost             TraceListRunsParamsSelect = "PROMPT_COST"
	TraceListRunsParamsSelectCompletionCost         TraceListRunsParamsSelect = "COMPLETION_COST"
	TraceListRunsParamsSelectPromptTokenDetails     TraceListRunsParamsSelect = "PROMPT_TOKEN_DETAILS"
	TraceListRunsParamsSelectCompletionTokenDetails TraceListRunsParamsSelect = "COMPLETION_TOKEN_DETAILS"
	TraceListRunsParamsSelectPromptCostDetails      TraceListRunsParamsSelect = "PROMPT_COST_DETAILS"
	TraceListRunsParamsSelectCompletionCostDetails  TraceListRunsParamsSelect = "COMPLETION_COST_DETAILS"
	TraceListRunsParamsSelectPriceModelID           TraceListRunsParamsSelect = "PRICE_MODEL_ID"
	TraceListRunsParamsSelectTags                   TraceListRunsParamsSelect = "TAGS"
	TraceListRunsParamsSelectAppPath                TraceListRunsParamsSelect = "APP_PATH"
	TraceListRunsParamsSelectAttachments            TraceListRunsParamsSelect = "ATTACHMENTS"
	TraceListRunsParamsSelectThreadEvaluationTime   TraceListRunsParamsSelect = "THREAD_EVALUATION_TIME"
	TraceListRunsParamsSelectIsInDataset            TraceListRunsParamsSelect = "IS_IN_DATASET"
	TraceListRunsParamsSelectLastQueuedAt           TraceListRunsParamsSelect = "LAST_QUEUED_AT"
	TraceListRunsParamsSelectShareURL               TraceListRunsParamsSelect = "SHARE_URL"
	TraceListRunsParamsSelectFeedbackStats          TraceListRunsParamsSelect = "FEEDBACK_STATS"
)

func (r TraceListRunsParamsSelect) IsKnown() bool {
	switch r {
	case TraceListRunsParamsSelectID, TraceListRunsParamsSelectName, TraceListRunsParamsSelectRunType, TraceListRunsParamsSelectStatus, TraceListRunsParamsSelectStartTime, TraceListRunsParamsSelectEndTime, TraceListRunsParamsSelectLatencySeconds, TraceListRunsParamsSelectFirstTokenTime, TraceListRunsParamsSelectError, TraceListRunsParamsSelectErrorPreview, TraceListRunsParamsSelectExtra, TraceListRunsParamsSelectMetadata, TraceListRunsParamsSelectEvents, TraceListRunsParamsSelectInputs, TraceListRunsParamsSelectInputsPreview, TraceListRunsParamsSelectOutputs, TraceListRunsParamsSelectOutputsPreview, TraceListRunsParamsSelectManifest, TraceListRunsParamsSelectParentRunIDs, TraceListRunsParamsSelectProjectID, TraceListRunsParamsSelectTraceID, TraceListRunsParamsSelectThreadID, TraceListRunsParamsSelectDottedOrder, TraceListRunsParamsSelectIsRoot, TraceListRunsParamsSelectReferenceExampleID, TraceListRunsParamsSelectReferenceDatasetID, TraceListRunsParamsSelectTotalTokens, TraceListRunsParamsSelectPromptTokens, TraceListRunsParamsSelectCompletionTokens, TraceListRunsParamsSelectTotalCost, TraceListRunsParamsSelectPromptCost, TraceListRunsParamsSelectCompletionCost, TraceListRunsParamsSelectPromptTokenDetails, TraceListRunsParamsSelectCompletionTokenDetails, TraceListRunsParamsSelectPromptCostDetails, TraceListRunsParamsSelectCompletionCostDetails, TraceListRunsParamsSelectPriceModelID, TraceListRunsParamsSelectTags, TraceListRunsParamsSelectAppPath, TraceListRunsParamsSelectAttachments, TraceListRunsParamsSelectThreadEvaluationTime, TraceListRunsParamsSelectIsInDataset, TraceListRunsParamsSelectLastQueuedAt, TraceListRunsParamsSelectShareURL, TraceListRunsParamsSelectFeedbackStats:
		return true
	}
	return false
}

type TraceQueryParams struct {
	// `cursor` is the opaque string returned in a previous response's `next_cursor`.
	Cursor param.Field[string] `json:"cursor"`
	// `max_start_time` is the exclusive upper bound for the root-run start time scan
	// (RFC3339). Defaults to the request time when omitted.
	MaxStartTime param.Field[time.Time] `json:"max_start_time" format:"date-time"`
	// `min_start_time` is the inclusive lower bound for the root-run start time scan
	// (RFC3339). Defaults to 24 hours before the request when omitted.
	MinStartTime param.Field[time.Time] `json:"min_start_time" format:"date-time"`
	// `page_size` is the maximum number of traces to return per page. Defaults to 20;
	// must be between 1 and 100 when set.
	PageSize param.Field[int64] `json:"page_size"`
	// `project_id` is the UUID of the tracing project that owns the traces. Required.
	ProjectID param.Field[string] `json:"project_id" format:"uuid"`
	// `selects` lists which properties to include on each returned trace. Properties
	// listed here are routed to the appropriate sub-object on each item:
	// `total_tokens`, `total_cost`, and `first_token_time` appear under
	// `trace_aggregates`; everything else appears under `root_run`. If omitted, only
	// `id` is returned on `root_run`.
	Selects param.Field[[]RunSelectField] `json:"selects"`
	// `trace_filter` narrows results to traces whose root run matches this LangSmith
	// filter expression. This filter targets root runs only — `is_root = true` is
	// implied. See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	TraceFilter param.Field[string] `json:"trace_filter"`
	// `trace_ids` is an optional fast-path restriction to a known set of trace UUIDs.
	// Equivalent in result to including each UUID in a `trace_filter`, but more
	// efficient at scale.
	TraceIDs param.Field[[]string] `json:"trace_ids" format:"uuid"`
	// `tree_filter` narrows results to traces containing at least one run anywhere in
	// the run tree (root or descendant) that matches this LangSmith filter expression.
	TreeFilter param.Field[string] `json:"tree_filter"`
}

func (r TraceQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
