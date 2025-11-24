// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/stainless-sdks/langsmith-api-go/internal/apijson"
	"github.com/stainless-sdks/langsmith-api-go/internal/param"
	"github.com/stainless-sdks/langsmith-api-go/internal/requestconfig"
	"github.com/stainless-sdks/langsmith-api-go/option"
	"github.com/stainless-sdks/langsmith-api-go/shared"
	"github.com/tidwall/gjson"
)

// DatasetGroupService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetGroupService] method instead.
type DatasetGroupService struct {
	Options []option.RequestOption
}

// NewDatasetGroupService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetGroupService(opts ...option.RequestOption) (r *DatasetGroupService) {
	r = &DatasetGroupService{}
	r.Options = opts
	return
}

// Fetch examples for a dataset, and fetch the runs for each example if they are
// associated with the given session_ids.
func (r *DatasetGroupService) Runs(ctx context.Context, datasetID string, body DatasetGroupRunsParams, opts ...option.RequestOption) (res *DatasetGroupRunsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/group/runs", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Response for grouped comparison view of dataset examples.
//
// Returns dataset examples grouped by a run metadata value (e.g., model='gpt-4').
// Optional filters are applied to all runs before grouping.
//
// Shows:
//
// - Which examples were executed with each metadata value
// - Per-session aggregate statistics for runs on those examples
// - The actual example data with their associated runs
//
// Used for comparing how different sessions performed on the same set of examples.
type DatasetGroupRunsResponse struct {
	Groups []DatasetGroupRunsResponseGroup `json:"groups,required"`
	JSON   datasetGroupRunsResponseJSON    `json:"-"`
}

// datasetGroupRunsResponseJSON contains the JSON metadata for the struct
// [DatasetGroupRunsResponse]
type datasetGroupRunsResponseJSON struct {
	Groups      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetGroupRunsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetGroupRunsResponseJSON) RawJSON() string {
	return r.raw
}

// Group of examples with a specific metadata value across multiple sessions.
//
// Extends RunGroupBase with:
//
//   - group_key: metadata value that defines this group
//   - sessions: per-session stats for runs matching this metadata value
//   - examples: shared examples across all sessions (intersection logic) with flat
//     array of runs (each run has session_id field for frontend to determine column)
//   - example_count: unique example count (pagination-aware, same across all
//     sessions due to intersection)
//
// Inherited from RunGroupBase:
//
//   - filter: metadata filter for this group (e.g., "and(eq(is_root, true),
//     and(eq(metadata_key, 'model'), eq(metadata_value, 'gpt-4')))")
//   - count: total run count across all sessions (includes duplicate runs)
//   - total_tokens, total_cost: aggregate across sessions
//   - min_start_time, max_start_time: time range across sessions
//   - latency_p50, latency_p99: aggregate latency stats across sessions
//   - feedback_stats: weighted average feedback across sessions
//
// Additional aggregate stats (from ExampleWithRunsGroup):
//
// - prompt_tokens, completion_tokens: separate token counts
// - prompt_cost, completion_cost: separate costs
// - error_rate: average error rate
type DatasetGroupRunsResponseGroup struct {
	ExampleCount     int64                                       `json:"example_count,required"`
	Examples         []ExampleWithRunsCh                         `json:"examples,required"`
	Filter           string                                      `json:"filter,required"`
	GroupKey         DatasetGroupRunsResponseGroupsGroupKeyUnion `json:"group_key,required"`
	Sessions         []DatasetGroupRunsResponseGroupsSession     `json:"sessions,required"`
	CompletionCost   string                                      `json:"completion_cost,nullable"`
	CompletionTokens int64                                       `json:"completion_tokens,nullable"`
	Count            int64                                       `json:"count,nullable"`
	ErrorRate        float64                                     `json:"error_rate,nullable"`
	FeedbackStats    interface{}                                 `json:"feedback_stats,nullable"`
	LatencyP50       float64                                     `json:"latency_p50,nullable"`
	LatencyP99       float64                                     `json:"latency_p99,nullable"`
	MaxStartTime     time.Time                                   `json:"max_start_time,nullable" format:"date-time"`
	MinStartTime     time.Time                                   `json:"min_start_time,nullable" format:"date-time"`
	PromptCost       string                                      `json:"prompt_cost,nullable"`
	PromptTokens     int64                                       `json:"prompt_tokens,nullable"`
	TotalCost        string                                      `json:"total_cost,nullable"`
	TotalTokens      int64                                       `json:"total_tokens,nullable"`
	JSON             datasetGroupRunsResponseGroupJSON           `json:"-"`
}

// datasetGroupRunsResponseGroupJSON contains the JSON metadata for the struct
// [DatasetGroupRunsResponseGroup]
type datasetGroupRunsResponseGroupJSON struct {
	ExampleCount     apijson.Field
	Examples         apijson.Field
	Filter           apijson.Field
	GroupKey         apijson.Field
	Sessions         apijson.Field
	CompletionCost   apijson.Field
	CompletionTokens apijson.Field
	Count            apijson.Field
	ErrorRate        apijson.Field
	FeedbackStats    apijson.Field
	LatencyP50       apijson.Field
	LatencyP99       apijson.Field
	MaxStartTime     apijson.Field
	MinStartTime     apijson.Field
	PromptCost       apijson.Field
	PromptTokens     apijson.Field
	TotalCost        apijson.Field
	TotalTokens      apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *DatasetGroupRunsResponseGroup) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetGroupRunsResponseGroupJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionString] or [shared.UnionFloat].
type DatasetGroupRunsResponseGroupsGroupKeyUnion interface {
	ImplementsDatasetGroupRunsResponseGroupsGroupKeyUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*DatasetGroupRunsResponseGroupsGroupKeyUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
	)
}

// TracerSession stats filtered to runs matching a specific metadata value.
//
// Extends TracerSession with:
//
//   - example_count: unique examples (vs run_count = total runs including
//     duplicates)
//   - filter: ClickHouse filter for fetching runs in this session/group
//   - min/max_start_time: time range for runs in this session/group
type DatasetGroupRunsResponseGroupsSession struct {
	ID                   string                                          `json:"id,required" format:"uuid"`
	Filter               string                                          `json:"filter,required"`
	TenantID             string                                          `json:"tenant_id,required" format:"uuid"`
	CompletionCost       string                                          `json:"completion_cost,nullable"`
	CompletionTokens     int64                                           `json:"completion_tokens,nullable"`
	DefaultDatasetID     string                                          `json:"default_dataset_id,nullable" format:"uuid"`
	Description          string                                          `json:"description,nullable"`
	EndTime              time.Time                                       `json:"end_time,nullable" format:"date-time"`
	ErrorRate            float64                                         `json:"error_rate,nullable"`
	ExampleCount         int64                                           `json:"example_count,nullable"`
	Extra                interface{}                                     `json:"extra,nullable"`
	FeedbackStats        interface{}                                     `json:"feedback_stats,nullable"`
	FirstTokenP50        float64                                         `json:"first_token_p50,nullable"`
	FirstTokenP99        float64                                         `json:"first_token_p99,nullable"`
	LastRunStartTime     time.Time                                       `json:"last_run_start_time,nullable" format:"date-time"`
	LastRunStartTimeLive time.Time                                       `json:"last_run_start_time_live,nullable" format:"date-time"`
	LatencyP50           float64                                         `json:"latency_p50,nullable"`
	LatencyP99           float64                                         `json:"latency_p99,nullable"`
	MaxStartTime         time.Time                                       `json:"max_start_time,nullable" format:"date-time"`
	MinStartTime         time.Time                                       `json:"min_start_time,nullable" format:"date-time"`
	Name                 string                                          `json:"name"`
	PromptCost           string                                          `json:"prompt_cost,nullable"`
	PromptTokens         int64                                           `json:"prompt_tokens,nullable"`
	ReferenceDatasetID   string                                          `json:"reference_dataset_id,nullable" format:"uuid"`
	RunCount             int64                                           `json:"run_count,nullable"`
	RunFacets            []interface{}                                   `json:"run_facets,nullable"`
	SessionFeedbackStats interface{}                                     `json:"session_feedback_stats,nullable"`
	StartTime            time.Time                                       `json:"start_time" format:"date-time"`
	StreamingRate        float64                                         `json:"streaming_rate,nullable"`
	TestRunNumber        int64                                           `json:"test_run_number,nullable"`
	TotalCost            string                                          `json:"total_cost,nullable"`
	TotalTokens          int64                                           `json:"total_tokens,nullable"`
	TraceTier            DatasetGroupRunsResponseGroupsSessionsTraceTier `json:"trace_tier,nullable"`
	JSON                 datasetGroupRunsResponseGroupsSessionJSON       `json:"-"`
}

// datasetGroupRunsResponseGroupsSessionJSON contains the JSON metadata for the
// struct [DatasetGroupRunsResponseGroupsSession]
type datasetGroupRunsResponseGroupsSessionJSON struct {
	ID                   apijson.Field
	Filter               apijson.Field
	TenantID             apijson.Field
	CompletionCost       apijson.Field
	CompletionTokens     apijson.Field
	DefaultDatasetID     apijson.Field
	Description          apijson.Field
	EndTime              apijson.Field
	ErrorRate            apijson.Field
	ExampleCount         apijson.Field
	Extra                apijson.Field
	FeedbackStats        apijson.Field
	FirstTokenP50        apijson.Field
	FirstTokenP99        apijson.Field
	LastRunStartTime     apijson.Field
	LastRunStartTimeLive apijson.Field
	LatencyP50           apijson.Field
	LatencyP99           apijson.Field
	MaxStartTime         apijson.Field
	MinStartTime         apijson.Field
	Name                 apijson.Field
	PromptCost           apijson.Field
	PromptTokens         apijson.Field
	ReferenceDatasetID   apijson.Field
	RunCount             apijson.Field
	RunFacets            apijson.Field
	SessionFeedbackStats apijson.Field
	StartTime            apijson.Field
	StreamingRate        apijson.Field
	TestRunNumber        apijson.Field
	TotalCost            apijson.Field
	TotalTokens          apijson.Field
	TraceTier            apijson.Field
	raw                  string
	ExtraFields          map[string]apijson.Field
}

func (r *DatasetGroupRunsResponseGroupsSession) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetGroupRunsResponseGroupsSessionJSON) RawJSON() string {
	return r.raw
}

type DatasetGroupRunsResponseGroupsSessionsTraceTier string

const (
	DatasetGroupRunsResponseGroupsSessionsTraceTierLonglived  DatasetGroupRunsResponseGroupsSessionsTraceTier = "longlived"
	DatasetGroupRunsResponseGroupsSessionsTraceTierShortlived DatasetGroupRunsResponseGroupsSessionsTraceTier = "shortlived"
)

func (r DatasetGroupRunsResponseGroupsSessionsTraceTier) IsKnown() bool {
	switch r {
	case DatasetGroupRunsResponseGroupsSessionsTraceTierLonglived, DatasetGroupRunsResponseGroupsSessionsTraceTierShortlived:
		return true
	}
	return false
}

type DatasetGroupRunsParams struct {
	GroupBy       param.Field[DatasetGroupRunsParamsGroupBy] `json:"group_by,required"`
	MetadataKey   param.Field[string]                        `json:"metadata_key,required"`
	SessionIDs    param.Field[[]string]                      `json:"session_ids,required" format:"uuid"`
	Filters       param.Field[map[string][]string]           `json:"filters"`
	Limit         param.Field[int64]                         `json:"limit"`
	Offset        param.Field[int64]                         `json:"offset"`
	PerGroupLimit param.Field[int64]                         `json:"per_group_limit"`
	Preview       param.Field[bool]                          `json:"preview"`
}

func (r DatasetGroupRunsParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetGroupRunsParamsGroupBy string

const (
	DatasetGroupRunsParamsGroupByRunMetadata     DatasetGroupRunsParamsGroupBy = "run_metadata"
	DatasetGroupRunsParamsGroupByExampleMetadata DatasetGroupRunsParamsGroupBy = "example_metadata"
)

func (r DatasetGroupRunsParamsGroupBy) IsKnown() bool {
	switch r {
	case DatasetGroupRunsParamsGroupByRunMetadata, DatasetGroupRunsParamsGroupByExampleMetadata:
		return true
	}
	return false
}
