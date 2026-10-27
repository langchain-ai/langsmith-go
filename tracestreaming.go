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

// TraceStreamingService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewTraceStreamingService] method instead.
type TraceStreamingService struct {
	Options []option.RequestOption
}

// NewTraceStreamingService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewTraceStreamingService(opts ...option.RequestOption) (r *TraceStreamingService) {
	r = &TraceStreamingService{}
	r.Options = opts
	return
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
//
// When `ai_search` is set, `Accept: text/event-stream` is required; requests
// without it return 406. AI search is unavailable on deployments that route
// queries to the v1 backend and returns 501 there.
//
// Self-hosted deployments require LangSmith `v0.16` or later.
func (r *TraceStreamingService) QueryStreaming(ctx context.Context, params TraceStreamingQueryParams, opts ...option.RequestOption) (stream *ssestream.Stream[shared.QueryPatchEvent]) {
	var (
		raw *http.Response
		err error
	)
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	path := "api/v2/traces/query"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[shared.QueryPatchEvent](ssestream.NewDecoder(raw), err)
}

// QueryStreamingAutoPaging yields the items of the stream one at a time across
// pages, following the next page cursor.
func (r *TraceStreamingService) QueryStreamingAutoPaging(ctx context.Context, params TraceStreamingQueryParams, opts ...option.RequestOption) *ssestream.AutoPager[shared.QueryPatchEvent, Trace] {
	return ssestream.NewAutoPager[shared.QueryPatchEvent, Trace](
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
			return r.QueryStreaming(ctx, params, callOpts...)
		},
	)
}

type TraceStreamingQueryResponse struct {
	// `items` is the page of traces.
	Items []Trace `json:"items"`
	// `next_cursor` is the opaque cursor for the next page. Null on the final page.
	NextCursor string                          `json:"next_cursor"`
	JSON       traceStreamingQueryResponseJSON `json:"-"`
}

// traceStreamingQueryResponseJSON contains the JSON metadata for the struct
// [TraceStreamingQueryResponse]
type traceStreamingQueryResponseJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TraceStreamingQueryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r traceStreamingQueryResponseJSON) RawJSON() string {
	return r.raw
}

type TraceStreamingQueryParams struct {
	// `ai_search` is a plain-language criterion evaluated against the messages from
	// the agent trajectory scoped to the trace. AND-ed with the ordinary filters.
	// Requires semantic filtering enabled for the deployment. Must contain nonempty
	// text of at most 2000 UTF-8 bytes.
	AISearch param.Field[string] `json:"ai_search"`
	// `cursor` is the opaque string returned in a previous response's `next_cursor`.
	Cursor param.Field[string] `json:"cursor"`
	// `max_feedback_modified_at` is the exclusive upper bound (RFC3339) on feedback
	// `modified_at`, applied to the feedback conditions in `trace_filter` or
	// `tree_filter` (for example eq(feedback_key, "correctness")). Requires at least
	// one feedback condition. Root runs must still fall within the `start_time`
	// window.
	MaxFeedbackModifiedAt param.Field[time.Time] `json:"max_feedback_modified_at" format:"date-time"`
	// `max_start_time` is the exclusive upper bound for the root-run start time scan
	// (RFC3339). Defaults to the request time when omitted.
	MaxStartTime param.Field[time.Time] `json:"max_start_time" format:"date-time"`
	// `min_feedback_modified_at` is the inclusive lower bound (RFC3339) on feedback
	// `modified_at`, applied to the feedback conditions in `trace_filter` or
	// `tree_filter` (for example eq(feedback_key, "correctness")). Requires at least
	// one feedback condition. Root runs must still fall within the `start_time`
	// window.
	MinFeedbackModifiedAt param.Field[time.Time] `json:"min_feedback_modified_at" format:"date-time"`
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
	Accept     param.Field[string] `header:"Accept"`
}

func (r TraceStreamingQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
