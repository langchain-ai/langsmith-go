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

// ThreadStreamingService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewThreadStreamingService] method instead.
type ThreadStreamingService struct {
	Options []option.RequestOption
}

// NewThreadStreamingService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewThreadStreamingService(opts ...option.RequestOption) (r *ThreadStreamingService) {
	r = &ThreadStreamingService{}
	r.Options = opts
	return
}

// GET with body payload — no resources created. Returns the threads in a tracing
// project matching the query passed as the request payload. Results are filtered
// by the given time range and optional filters, with cursor-based pagination.
//
// When `ai_search` is set, `Accept: text/event-stream` is required; requests
// without it return 406. AI search is unavailable on deployments that route
// queries to the v1 backend and returns 501 there.
//
// Self-hosted deployments require LangSmith `v0.16` or later.
func (r *ThreadStreamingService) QueryStreaming(ctx context.Context, params ThreadStreamingQueryParams, opts ...option.RequestOption) (stream *ssestream.Stream[shared.QueryPatchEvent]) {
	var (
		raw *http.Response
		err error
	)
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "text/event-stream")}, opts...)
	path := "api/v2/threads/query"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &raw, opts...)
	return ssestream.NewStream[shared.QueryPatchEvent](ssestream.NewDecoder(raw), err)
}

// QueryStreamingAutoPaging yields the items of the stream one at a time across
// pages, following the next page cursor.
func (r *ThreadStreamingService) QueryStreamingAutoPaging(ctx context.Context, params ThreadStreamingQueryParams, opts ...option.RequestOption) *ssestream.AutoPager[shared.QueryPatchEvent, Thread] {
	return ssestream.NewAutoPager[shared.QueryPatchEvent, Thread](
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

type ThreadStreamingQueryResponse struct {
	// `items` is the page of thread summaries, sorted by the thread's most recent
	// activity.
	Items []Thread `json:"items"`
	// `next_cursor` is the opaque cursor to pass as `cursor` on the next request. Null
	// on the final page.
	NextCursor string                           `json:"next_cursor"`
	JSON       threadStreamingQueryResponseJSON `json:"-"`
}

// threadStreamingQueryResponseJSON contains the JSON metadata for the struct
// [ThreadStreamingQueryResponse]
type threadStreamingQueryResponseJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadStreamingQueryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadStreamingQueryResponseJSON) RawJSON() string {
	return r.raw
}

type ThreadStreamingQueryParams struct {
	// `ai_search` is a plain-language criterion evaluated against the messages from
	// the agent trajectory scoped to the thread. AND-ed with the ordinary filters.
	// Requires semantic filtering enabled for the deployment. Must contain nonempty
	// text of at most 2000 UTF-8 bytes.
	AISearch param.Field[string] `json:"ai_search"`
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
	// `thread_filter` narrows results using a LangSmith filter expression evaluated
	// against each complete thread summary. Self-hosted deployments require LangSmith
	// v0.17 or later; unsupported deployments return 501. See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	ThreadFilter param.Field[string] `json:"thread_filter"`
	// `trace_filter` narrows results to threads containing at least one trace whose
	// root run matches this LangSmith filter expression. Trace-level aggregate fields
	// are evaluated using the complete trace summary. Self-hosted deployments require
	// LangSmith v0.17 or later; unsupported deployments return 501. See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	TraceFilter param.Field[string] `json:"trace_filter"`
	// `tree_filter` narrows results to threads containing at least one trace with a
	// matching run anywhere in its run tree. Self-hosted deployments require LangSmith
	// v0.17 or later; unsupported deployments return 501. See
	// https://docs.langchain.com/langsmith/trace-query-syntax#filter-query-language
	// for syntax.
	TreeFilter param.Field[string] `json:"tree_filter"`
	Accept     param.Field[string] `header:"Accept"`
}

func (r ThreadStreamingQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
