// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// AISearchService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAISearchService] method instead.
type AISearchService struct {
	Options []option.RequestOption
}

// NewAISearchService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAISearchService(opts ...option.RequestOption) (r *AISearchService) {
	r = &AISearchService{}
	r.Options = opts
	return
}

// Translates a natural-language search into the scope, filters, semantic criterion
// and time bounds for the runs, traces or threads query API, without running the
// search. No resources are created. Send the returned fields to the query API
// named by scope, with the same projects. When scope is omitted, it is inferred
// from the search and defaults to traces. A time period in the search sets
// min_start_time or max_start_time; bounds the search does not mention are omitted
// so you can keep your own time range. Requires AI search to be enabled for the
// organization and read access to every project. Returns a single JSON response.
func (r *AISearchService) Parse(ctx context.Context, body AISearchParseParams, opts ...option.RequestOption) (res *AISearchParseResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v2/ai-search/parse"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type AISearchParseResponse struct {
	// Interpretations to choose from. Returned only when allow_clarification is true
	// and the search is ambiguous; no other fields are set.
	Clarification AISearchParseResponseClarification `json:"clarification"`
	// Filter expression for the query API's filter parameter. Omitted when there is
	// none.
	Filter string `json:"filter"`
	// Upper bound on start time (RFC 3339, UTC) for the query API's max_start_time
	// parameter. Omitted when the search names no end; keep your own time range.
	MaxStartTime time.Time `json:"max_start_time" format:"date-time"`
	// Lower bound on start time (RFC 3339, UTC) for the query API's min_start_time
	// parameter. Omitted when the search names no start; keep your own time range.
	MinStartTime time.Time `json:"min_start_time" format:"date-time"`
	// Query API to send the plan to: runs, traces or threads.
	Scope AISearchParseResponseScope `json:"scope"`
	// Plain-language criterion for the query API's sem_filter parameter, matched
	// against conversation messages. Null when the search has no conversational
	// criterion.
	SemFilter string `json:"sem_filter" api:"nullable"`
	// Filter expression for the threads query API's thread_filter parameter. Omitted
	// when there is none.
	ThreadFilter string `json:"thread_filter"`
	// Filter expression for the query API's trace_filter parameter. Omitted when there
	// is none.
	TraceFilter string `json:"trace_filter"`
	// Filter expression for the query API's tree_filter parameter. Omitted when there
	// is none.
	TreeFilter string                    `json:"tree_filter"`
	JSON       aiSearchParseResponseJSON `json:"-"`
}

// aiSearchParseResponseJSON contains the JSON metadata for the struct
// [AISearchParseResponse]
type aiSearchParseResponseJSON struct {
	Clarification apijson.Field
	Filter        apijson.Field
	MaxStartTime  apijson.Field
	MinStartTime  apijson.Field
	Scope         apijson.Field
	SemFilter     apijson.Field
	ThreadFilter  apijson.Field
	TraceFilter   apijson.Field
	TreeFilter    apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *AISearchParseResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r aiSearchParseResponseJSON) RawJSON() string {
	return r.raw
}

// Interpretations to choose from. Returned only when allow_clarification is true
// and the search is ambiguous; no other fields are set.
type AISearchParseResponseClarification struct {
	// Two or three interpretations of the search.
	Options []AISearchParseResponseClarificationOption `json:"options"`
	JSON    aiSearchParseResponseClarificationJSON     `json:"-"`
}

// aiSearchParseResponseClarificationJSON contains the JSON metadata for the struct
// [AISearchParseResponseClarification]
type aiSearchParseResponseClarificationJSON struct {
	Options     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AISearchParseResponseClarification) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r aiSearchParseResponseClarificationJSON) RawJSON() string {
	return r.raw
}

type AISearchParseResponseClarificationOption struct {
	// Filter expression for the query API's filter parameter. Omitted when there is
	// none.
	Filter string `json:"filter"`
	// Short description of this interpretation to show the user.
	Label string `json:"label"`
	// Upper bound on start time (RFC 3339, UTC) for the query API's max_start_time
	// parameter. Omitted when the search names no end; keep your own time range.
	MaxStartTime time.Time `json:"max_start_time" format:"date-time"`
	// Lower bound on start time (RFC 3339, UTC) for the query API's min_start_time
	// parameter. Omitted when the search names no start; keep your own time range.
	MinStartTime time.Time `json:"min_start_time" format:"date-time"`
	// Query API this interpretation is planned for.
	Scope AISearchParseResponseClarificationOptionsScope `json:"scope"`
	// Plain-language criterion for the query API's sem_filter parameter, matched
	// against conversation messages. Null when the search has no conversational
	// criterion.
	SemFilter string `json:"sem_filter" api:"nullable"`
	// Filter expression for the threads query API's thread_filter parameter. Omitted
	// when there is none.
	ThreadFilter string `json:"thread_filter"`
	// Filter expression for the query API's trace_filter parameter. Omitted when there
	// is none.
	TraceFilter string `json:"trace_filter"`
	// Filter expression for the query API's tree_filter parameter. Omitted when there
	// is none.
	TreeFilter string                                       `json:"tree_filter"`
	JSON       aiSearchParseResponseClarificationOptionJSON `json:"-"`
}

// aiSearchParseResponseClarificationOptionJSON contains the JSON metadata for the
// struct [AISearchParseResponseClarificationOption]
type aiSearchParseResponseClarificationOptionJSON struct {
	Filter       apijson.Field
	Label        apijson.Field
	MaxStartTime apijson.Field
	MinStartTime apijson.Field
	Scope        apijson.Field
	SemFilter    apijson.Field
	ThreadFilter apijson.Field
	TraceFilter  apijson.Field
	TreeFilter   apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *AISearchParseResponseClarificationOption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r aiSearchParseResponseClarificationOptionJSON) RawJSON() string {
	return r.raw
}

// Query API this interpretation is planned for.
type AISearchParseResponseClarificationOptionsScope string

const (
	AISearchParseResponseClarificationOptionsScopeRuns    AISearchParseResponseClarificationOptionsScope = "runs"
	AISearchParseResponseClarificationOptionsScopeTraces  AISearchParseResponseClarificationOptionsScope = "traces"
	AISearchParseResponseClarificationOptionsScopeThreads AISearchParseResponseClarificationOptionsScope = "threads"
)

func (r AISearchParseResponseClarificationOptionsScope) IsKnown() bool {
	switch r {
	case AISearchParseResponseClarificationOptionsScopeRuns, AISearchParseResponseClarificationOptionsScopeTraces, AISearchParseResponseClarificationOptionsScopeThreads:
		return true
	}
	return false
}

// Query API to send the plan to: runs, traces or threads.
type AISearchParseResponseScope string

const (
	AISearchParseResponseScopeRuns    AISearchParseResponseScope = "runs"
	AISearchParseResponseScopeTraces  AISearchParseResponseScope = "traces"
	AISearchParseResponseScopeThreads AISearchParseResponseScope = "threads"
)

func (r AISearchParseResponseScope) IsKnown() bool {
	switch r {
	case AISearchParseResponseScopeRuns, AISearchParseResponseScopeTraces, AISearchParseResponseScopeThreads:
		return true
	}
	return false
}

type AISearchParseParams struct {
	// Natural-language description of the runs, traces or threads to find, at most
	// 2000 UTF-8 bytes.
	AISearch param.Field[string] `json:"ai_search" api:"required"`
	// Tracing project UUIDs used to interpret the search. Send the same projects to
	// the query API.
	ProjectIDs param.Field[[]string] `json:"project_ids" api:"required" format:"uuid"`
	// When true, an ambiguous search returns clarification options instead of a single
	// plan. Defaults to false.
	AllowClarification param.Field[bool] `json:"allow_clarification"`
	// Query API to plan for: runs, traces or threads. When omitted, it is inferred
	// from the search and defaults to traces.
	Scope param.Field[AISearchParseParamsScope] `json:"scope"`
}

func (r AISearchParseParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Query API to plan for: runs, traces or threads. When omitted, it is inferred
// from the search and defaults to traces.
type AISearchParseParamsScope string

const (
	AISearchParseParamsScopeRuns    AISearchParseParamsScope = "runs"
	AISearchParseParamsScopeTraces  AISearchParseParamsScope = "traces"
	AISearchParseParamsScopeThreads AISearchParseParamsScope = "threads"
)

func (r AISearchParseParamsScope) IsKnown() bool {
	switch r {
	case AISearchParseParamsScopeRuns, AISearchParseParamsScopeTraces, AISearchParseParamsScopeThreads:
		return true
	}
	return false
}
