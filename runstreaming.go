// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/ssestream"
	"github.com/langchain-ai/langsmith-go/shared"
)

// RunStreamingService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewRunStreamingService] method instead.
type RunStreamingService struct {
	Options []option.RequestOption
}

// NewRunStreamingService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewRunStreamingService(opts ...option.RequestOption) (r *RunStreamingService) {
	r = &RunStreamingService{}
	r.Options = opts
	return
}

// Returns a paginated list of runs for the given projects within min/max
// start_time. Supports filters, cursor pagination, and `selects` to select fields
// to return.
//
// When `ai_search` is set, `Accept: text/event-stream` is required; requests
// without it return 406. AI search is unavailable on deployments that route
// queries to the v1 backend and returns 501 there.
//
// Self-hosted deployments require LangSmith `v0.16` or later.
func (r *RunStreamingService) QueryV2Streaming(ctx context.Context, params RunStreamingQueryV2Params, opts ...option.RequestOption) (stream *ssestream.Stream[shared.QueryPatchEvent]) {
	var (
		raw *http.Response
		err error
	)
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	path := "api/v2/runs/query"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[shared.QueryPatchEvent](ssestream.NewDecoder(raw), err)
}

// QueryV2StreamingAutoPaging yields the items of the stream one at a time across
// pages, following the next page cursor.
func (r *RunStreamingService) QueryV2StreamingAutoPaging(ctx context.Context, params RunStreamingQueryV2Params, opts ...option.RequestOption) *ssestream.AutoPager[shared.QueryPatchEvent, Run] {
	return ssestream.NewAutoPager[shared.QueryPatchEvent, Run](
		ssestream.AutoPagerConfig[shared.QueryPatchEvent]{
			Patches: func(event shared.QueryPatchEvent) []ssestream.Patch {
				patches := make([]ssestream.Patch, 0, len(event))
				for _, op := range event {
					patches = append(patches, ssestream.Patch{Op: string(op.Op), Path: op.Path, Value: []byte(op.JSON.Value.Raw())})
				}
				return patches
			},
			ItemsPath:      "/items",
			NextCursorPath: "/next_cursor",
		},
		func(cursor string) *ssestream.Stream[shared.QueryPatchEvent] {
			callOpts := opts
			if cursor != "" {
				callOpts = append(append([]option.RequestOption{}, opts...), option.WithJSONSet("cursor", cursor))
			}
			return r.QueryV2Streaming(ctx, params, callOpts...)
		},
	)
}

type RunStreamingQueryV2Response struct {
	// `items` is the page of runs, sorted by `start_time` descending.
	Items []Run `json:"items"`
	// `next_cursor` is the opaque cursor to pass as `cursor` on the next request. Null
	// on the final page.
	NextCursor string                          `json:"next_cursor"`
	JSON       runStreamingQueryV2ResponseJSON `json:"-"`
}

// runStreamingQueryV2ResponseJSON contains the JSON metadata for the struct
// [RunStreamingQueryV2Response]
type runStreamingQueryV2ResponseJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunStreamingQueryV2Response) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runStreamingQueryV2ResponseJSON) RawJSON() string {
	return r.raw
}

type RunStreamingQueryV2Params struct {
	// `ai_search` is a plain-language criterion evaluated against the messages from
	// the agent trajectory scoped to the run. AND-ed with the ordinary filters.
	// Requires semantic filtering enabled for the deployment. Must contain nonempty
	// text of at most 2000 UTF-8 bytes. Not supported on public dataset queries.
	AISearch param.Field[string] `json:"ai_search"`
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
	// `include_subtree_usage` enriches selected token and cost fields on root and
	// intermediate runs with subtree totals, and selected `first_token_time` with the
	// earliest observed subtree timestamp. Descendants contribute regardless of this
	// request's filters or time bounds. Runs without observed descendants retain their
	// original values. Defaults to false.
	IncludeSubtreeUsage param.Field[bool] `json:"include_subtree_usage"`
	// `is_root` returns only root runs (true) or only non-root runs (false).
	IsRoot param.Field[bool] `json:"is_root"`
	// `max_feedback_modified_at` is the exclusive upper bound (RFC3339) on feedback
	// `modified_at`, applied to the feedback conditions in `filter`, `trace_filter`,
	// or `tree_filter` (for example eq(feedback_key, "correctness")). Requires at
	// least one feedback condition. Runs must still fall within the `start_time`
	// window.
	MaxFeedbackModifiedAt param.Field[time.Time] `json:"max_feedback_modified_at" format:"date-time"`
	// `max_start_time` is the upper bound for run `start_time` (RFC3339). Defaults to
	// now.
	MaxStartTime param.Field[time.Time] `json:"max_start_time" format:"date-time"`
	// `min_feedback_modified_at` is the inclusive lower bound (RFC3339) on feedback
	// `modified_at`, applied to the feedback conditions in `filter`, `trace_filter`,
	// or `tree_filter` (for example eq(feedback_key, "correctness")). Requires at
	// least one feedback condition. Runs must still fall within the `start_time`
	// window.
	MinFeedbackModifiedAt param.Field[time.Time] `json:"min_feedback_modified_at" format:"date-time"`
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
	RunType param.Field[RunType] `json:"run_type"`
	// `selects` lists which properties to include on each returned run. If omitted,
	// only `id` is returned. Properties not listed are omitted from each run object.
	Selects param.Field[[]RunSelectField] `json:"selects"`
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

func (r RunStreamingQueryV2Params) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
