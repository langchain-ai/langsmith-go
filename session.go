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
	"github.com/langchain-ai/langsmith-go/shared"
	"github.com/tidwall/gjson"
)

// SessionService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionService] method instead.
type SessionService struct {
	Options  []option.RequestOption
	Insights *SessionInsightService
}

// NewSessionService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSessionService(opts ...option.RequestOption) (r *SessionService) {
	r = &SessionService{}
	r.Options = opts
	r.Insights = NewSessionInsightService(opts...)
	return
}

// Create a new session.
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *TracerSessionWithoutVirtualFields, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

// Get a specific session.
func (r *SessionService) Get(ctx context.Context, sessionID string, params SessionGetParams, opts ...option.RequestOption) (res *TracerSession, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return
}

// Create a new session.
func (r *SessionService) Update(ctx context.Context, sessionID string, body SessionUpdateParams, opts ...option.RequestOption) (res *TracerSessionWithoutVirtualFields, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// Get all sessions.
func (r *SessionService) List(ctx context.Context, params SessionListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[TracerSession], err error) {
	var raw *http.Response
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/sessions"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, params, &res, opts...)
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

// Get all sessions.
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[TracerSession] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, params, opts...))
}

// Delete a specific session.
func (r *SessionService) Delete(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *SessionDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Get a prebuilt dashboard for a tracing project.
func (r *SessionService) Dashboard(ctx context.Context, sessionID string, params SessionDashboardParams, opts ...option.RequestOption) (res *CustomChartsSection, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s/dashboard", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return
}

type CustomChartsSection struct {
	ID          string                          `json:"id,required" format:"uuid"`
	Charts      []CustomChartsSectionChart      `json:"charts,required"`
	Title       string                          `json:"title,required"`
	Description string                          `json:"description,nullable"`
	Index       int64                           `json:"index,nullable"`
	SessionID   string                          `json:"session_id,nullable" format:"uuid"`
	SubSections []CustomChartsSectionSubSection `json:"sub_sections,nullable"`
	JSON        customChartsSectionJSON         `json:"-"`
}

// customChartsSectionJSON contains the JSON metadata for the struct
// [CustomChartsSection]
type customChartsSectionJSON struct {
	ID          apijson.Field
	Charts      apijson.Field
	Title       apijson.Field
	Description apijson.Field
	Index       apijson.Field
	SessionID   apijson.Field
	SubSections apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSection) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChart struct {
	ID string `json:"id,required" format:"uuid"`
	// Enum for custom chart types.
	ChartType     CustomChartsSectionChartsChartType     `json:"chart_type,required"`
	Data          []CustomChartsSectionChartsData        `json:"data,required"`
	Index         int64                                  `json:"index,required"`
	Series        []CustomChartsSectionChartsSeries      `json:"series,required"`
	Title         string                                 `json:"title,required"`
	CommonFilters CustomChartsSectionChartsCommonFilters `json:"common_filters,nullable"`
	Description   string                                 `json:"description,nullable"`
	Metadata      map[string]interface{}                 `json:"metadata,nullable"`
	JSON          customChartsSectionChartJSON           `json:"-"`
}

// customChartsSectionChartJSON contains the JSON metadata for the struct
// [CustomChartsSectionChart]
type customChartsSectionChartJSON struct {
	ID            apijson.Field
	ChartType     apijson.Field
	Data          apijson.Field
	Index         apijson.Field
	Series        apijson.Field
	Title         apijson.Field
	CommonFilters apijson.Field
	Description   apijson.Field
	Metadata      apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *CustomChartsSectionChart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartJSON) RawJSON() string {
	return r.raw
}

// Enum for custom chart types.
type CustomChartsSectionChartsChartType string

const (
	CustomChartsSectionChartsChartTypeLine CustomChartsSectionChartsChartType = "line"
	CustomChartsSectionChartsChartTypeBar  CustomChartsSectionChartsChartType = "bar"
)

func (r CustomChartsSectionChartsChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsChartTypeLine, CustomChartsSectionChartsChartTypeBar:
		return true
	}
	return false
}

type CustomChartsSectionChartsData struct {
	SeriesID  string                                  `json:"series_id,required"`
	Timestamp time.Time                               `json:"timestamp,required" format:"date-time"`
	Value     CustomChartsSectionChartsDataValueUnion `json:"value,required,nullable"`
	Group     string                                  `json:"group,nullable"`
	JSON      customChartsSectionChartsDataJSON       `json:"-"`
}

// customChartsSectionChartsDataJSON contains the JSON metadata for the struct
// [CustomChartsSectionChartsData]
type customChartsSectionChartsDataJSON struct {
	SeriesID    apijson.Field
	Timestamp   apijson.Field
	Value       apijson.Field
	Group       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsDataJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionFloat] or
// [CustomChartsSectionChartsDataValueMap].
type CustomChartsSectionChartsDataValueUnion interface {
	ImplementsCustomChartsSectionChartsDataValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsDataValueUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsDataValueMap{}),
		},
	)
}

type CustomChartsSectionChartsDataValueMap map[string]interface{}

func (r CustomChartsSectionChartsDataValueMap) ImplementsCustomChartsSectionChartsDataValueUnion() {}

type CustomChartsSectionChartsSeries struct {
	ID string `json:"id,required" format:"uuid"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric      CustomChartsSectionChartsSeriesMetric  `json:"metric,required"`
	Name        string                                 `json:"name,required"`
	FeedbackKey string                                 `json:"feedback_key,nullable"`
	Filters     CustomChartsSectionChartsSeriesFilters `json:"filters,nullable"`
	// Include additional information about where the group_by param was set.
	GroupBy CustomChartsSectionChartsSeriesGroupBy `json:"group_by,nullable"`
	// LGP Metrics you can chart.
	ProjectMetric CustomChartsSectionChartsSeriesProjectMetric `json:"project_metric,nullable"`
	WorkspaceID   string                                       `json:"workspace_id,nullable" format:"uuid"`
	JSON          customChartsSectionChartsSeriesJSON          `json:"-"`
}

// customChartsSectionChartsSeriesJSON contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeries]
type customChartsSectionChartsSeriesJSON struct {
	ID            apijson.Field
	Metric        apijson.Field
	Name          apijson.Field
	FeedbackKey   apijson.Field
	Filters       apijson.Field
	GroupBy       apijson.Field
	ProjectMetric apijson.Field
	WorkspaceID   apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeries) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesJSON) RawJSON() string {
	return r.raw
}

// Metrics you can chart. Feedback metrics are not available for
// organization-scoped charts.
type CustomChartsSectionChartsSeriesMetric string

const (
	CustomChartsSectionChartsSeriesMetricRunCount            CustomChartsSectionChartsSeriesMetric = "run_count"
	CustomChartsSectionChartsSeriesMetricLatencyP50          CustomChartsSectionChartsSeriesMetric = "latency_p50"
	CustomChartsSectionChartsSeriesMetricLatencyP99          CustomChartsSectionChartsSeriesMetric = "latency_p99"
	CustomChartsSectionChartsSeriesMetricLatencyAvg          CustomChartsSectionChartsSeriesMetric = "latency_avg"
	CustomChartsSectionChartsSeriesMetricFirstTokenP50       CustomChartsSectionChartsSeriesMetric = "first_token_p50"
	CustomChartsSectionChartsSeriesMetricFirstTokenP99       CustomChartsSectionChartsSeriesMetric = "first_token_p99"
	CustomChartsSectionChartsSeriesMetricTotalTokens         CustomChartsSectionChartsSeriesMetric = "total_tokens"
	CustomChartsSectionChartsSeriesMetricPromptTokens        CustomChartsSectionChartsSeriesMetric = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricCompletionTokens    CustomChartsSectionChartsSeriesMetric = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricMedianTokens        CustomChartsSectionChartsSeriesMetric = "median_tokens"
	CustomChartsSectionChartsSeriesMetricCompletionTokensP50 CustomChartsSectionChartsSeriesMetric = "completion_tokens_p50"
	CustomChartsSectionChartsSeriesMetricPromptTokensP50     CustomChartsSectionChartsSeriesMetric = "prompt_tokens_p50"
	CustomChartsSectionChartsSeriesMetricTokensP99           CustomChartsSectionChartsSeriesMetric = "tokens_p99"
	CustomChartsSectionChartsSeriesMetricCompletionTokensP99 CustomChartsSectionChartsSeriesMetric = "completion_tokens_p99"
	CustomChartsSectionChartsSeriesMetricPromptTokensP99     CustomChartsSectionChartsSeriesMetric = "prompt_tokens_p99"
	CustomChartsSectionChartsSeriesMetricFeedback            CustomChartsSectionChartsSeriesMetric = "feedback"
	CustomChartsSectionChartsSeriesMetricFeedbackScoreAvg    CustomChartsSectionChartsSeriesMetric = "feedback_score_avg"
	CustomChartsSectionChartsSeriesMetricFeedbackValues      CustomChartsSectionChartsSeriesMetric = "feedback_values"
	CustomChartsSectionChartsSeriesMetricTotalCost           CustomChartsSectionChartsSeriesMetric = "total_cost"
	CustomChartsSectionChartsSeriesMetricPromptCost          CustomChartsSectionChartsSeriesMetric = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricCompletionCost      CustomChartsSectionChartsSeriesMetric = "completion_cost"
	CustomChartsSectionChartsSeriesMetricErrorRate           CustomChartsSectionChartsSeriesMetric = "error_rate"
	CustomChartsSectionChartsSeriesMetricStreamingRate       CustomChartsSectionChartsSeriesMetric = "streaming_rate"
	CustomChartsSectionChartsSeriesMetricCostP50             CustomChartsSectionChartsSeriesMetric = "cost_p50"
	CustomChartsSectionChartsSeriesMetricCostP99             CustomChartsSectionChartsSeriesMetric = "cost_p99"
)

func (r CustomChartsSectionChartsSeriesMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricRunCount, CustomChartsSectionChartsSeriesMetricLatencyP50, CustomChartsSectionChartsSeriesMetricLatencyP99, CustomChartsSectionChartsSeriesMetricLatencyAvg, CustomChartsSectionChartsSeriesMetricFirstTokenP50, CustomChartsSectionChartsSeriesMetricFirstTokenP99, CustomChartsSectionChartsSeriesMetricTotalTokens, CustomChartsSectionChartsSeriesMetricPromptTokens, CustomChartsSectionChartsSeriesMetricCompletionTokens, CustomChartsSectionChartsSeriesMetricMedianTokens, CustomChartsSectionChartsSeriesMetricCompletionTokensP50, CustomChartsSectionChartsSeriesMetricPromptTokensP50, CustomChartsSectionChartsSeriesMetricTokensP99, CustomChartsSectionChartsSeriesMetricCompletionTokensP99, CustomChartsSectionChartsSeriesMetricPromptTokensP99, CustomChartsSectionChartsSeriesMetricFeedback, CustomChartsSectionChartsSeriesMetricFeedbackScoreAvg, CustomChartsSectionChartsSeriesMetricFeedbackValues, CustomChartsSectionChartsSeriesMetricTotalCost, CustomChartsSectionChartsSeriesMetricPromptCost, CustomChartsSectionChartsSeriesMetricCompletionCost, CustomChartsSectionChartsSeriesMetricErrorRate, CustomChartsSectionChartsSeriesMetricStreamingRate, CustomChartsSectionChartsSeriesMetricCostP50, CustomChartsSectionChartsSeriesMetricCostP99:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesFilters struct {
	Filter      string                                     `json:"filter,nullable"`
	Session     []string                                   `json:"session,nullable" format:"uuid"`
	TraceFilter string                                     `json:"trace_filter,nullable"`
	TreeFilter  string                                     `json:"tree_filter,nullable"`
	JSON        customChartsSectionChartsSeriesFiltersJSON `json:"-"`
}

// customChartsSectionChartsSeriesFiltersJSON contains the JSON metadata for the
// struct [CustomChartsSectionChartsSeriesFilters]
type customChartsSectionChartsSeriesFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesFiltersJSON) RawJSON() string {
	return r.raw
}

// Include additional information about where the group_by param was set.
type CustomChartsSectionChartsSeriesGroupBy struct {
	Attribute CustomChartsSectionChartsSeriesGroupByAttribute `json:"attribute,required"`
	MaxGroups int64                                           `json:"max_groups"`
	Path      string                                          `json:"path,nullable"`
	SetBy     CustomChartsSectionChartsSeriesGroupBySetBy     `json:"set_by,nullable"`
	JSON      customChartsSectionChartsSeriesGroupByJSON      `json:"-"`
}

// customChartsSectionChartsSeriesGroupByJSON contains the JSON metadata for the
// struct [CustomChartsSectionChartsSeriesGroupBy]
type customChartsSectionChartsSeriesGroupByJSON struct {
	Attribute   apijson.Field
	MaxGroups   apijson.Field
	Path        apijson.Field
	SetBy       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesGroupBy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesGroupByJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSeriesGroupByAttribute string

const (
	CustomChartsSectionChartsSeriesGroupByAttributeName     CustomChartsSectionChartsSeriesGroupByAttribute = "name"
	CustomChartsSectionChartsSeriesGroupByAttributeRunType  CustomChartsSectionChartsSeriesGroupByAttribute = "run_type"
	CustomChartsSectionChartsSeriesGroupByAttributeTag      CustomChartsSectionChartsSeriesGroupByAttribute = "tag"
	CustomChartsSectionChartsSeriesGroupByAttributeMetadata CustomChartsSectionChartsSeriesGroupByAttribute = "metadata"
)

func (r CustomChartsSectionChartsSeriesGroupByAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesGroupByAttributeName, CustomChartsSectionChartsSeriesGroupByAttributeRunType, CustomChartsSectionChartsSeriesGroupByAttributeTag, CustomChartsSectionChartsSeriesGroupByAttributeMetadata:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesGroupBySetBy string

const (
	CustomChartsSectionChartsSeriesGroupBySetBySection CustomChartsSectionChartsSeriesGroupBySetBy = "section"
	CustomChartsSectionChartsSeriesGroupBySetBySeries  CustomChartsSectionChartsSeriesGroupBySetBy = "series"
)

func (r CustomChartsSectionChartsSeriesGroupBySetBy) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesGroupBySetBySection, CustomChartsSectionChartsSeriesGroupBySetBySeries:
		return true
	}
	return false
}

// LGP Metrics you can chart.
type CustomChartsSectionChartsSeriesProjectMetric string

const (
	CustomChartsSectionChartsSeriesProjectMetricMemoryUsage             CustomChartsSectionChartsSeriesProjectMetric = "memory_usage"
	CustomChartsSectionChartsSeriesProjectMetricCPUUsage                CustomChartsSectionChartsSeriesProjectMetric = "cpu_usage"
	CustomChartsSectionChartsSeriesProjectMetricDiskUsage               CustomChartsSectionChartsSeriesProjectMetric = "disk_usage"
	CustomChartsSectionChartsSeriesProjectMetricRestartCount            CustomChartsSectionChartsSeriesProjectMetric = "restart_count"
	CustomChartsSectionChartsSeriesProjectMetricReplicaCount            CustomChartsSectionChartsSeriesProjectMetric = "replica_count"
	CustomChartsSectionChartsSeriesProjectMetricWorkerCount             CustomChartsSectionChartsSeriesProjectMetric = "worker_count"
	CustomChartsSectionChartsSeriesProjectMetricLgRunCount              CustomChartsSectionChartsSeriesProjectMetric = "lg_run_count"
	CustomChartsSectionChartsSeriesProjectMetricResponsesPerSecond      CustomChartsSectionChartsSeriesProjectMetric = "responses_per_second"
	CustomChartsSectionChartsSeriesProjectMetricErrorResponsesPerSecond CustomChartsSectionChartsSeriesProjectMetric = "error_responses_per_second"
	CustomChartsSectionChartsSeriesProjectMetricP95Latency              CustomChartsSectionChartsSeriesProjectMetric = "p95_latency"
)

func (r CustomChartsSectionChartsSeriesProjectMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesProjectMetricMemoryUsage, CustomChartsSectionChartsSeriesProjectMetricCPUUsage, CustomChartsSectionChartsSeriesProjectMetricDiskUsage, CustomChartsSectionChartsSeriesProjectMetricRestartCount, CustomChartsSectionChartsSeriesProjectMetricReplicaCount, CustomChartsSectionChartsSeriesProjectMetricWorkerCount, CustomChartsSectionChartsSeriesProjectMetricLgRunCount, CustomChartsSectionChartsSeriesProjectMetricResponsesPerSecond, CustomChartsSectionChartsSeriesProjectMetricErrorResponsesPerSecond, CustomChartsSectionChartsSeriesProjectMetricP95Latency:
		return true
	}
	return false
}

type CustomChartsSectionChartsCommonFilters struct {
	Filter      string                                     `json:"filter,nullable"`
	Session     []string                                   `json:"session,nullable" format:"uuid"`
	TraceFilter string                                     `json:"trace_filter,nullable"`
	TreeFilter  string                                     `json:"tree_filter,nullable"`
	JSON        customChartsSectionChartsCommonFiltersJSON `json:"-"`
}

// customChartsSectionChartsCommonFiltersJSON contains the JSON metadata for the
// struct [CustomChartsSectionChartsCommonFilters]
type customChartsSectionChartsCommonFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsCommonFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsCommonFiltersJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSection struct {
	ID          string                                `json:"id,required" format:"uuid"`
	Charts      []CustomChartsSectionSubSectionsChart `json:"charts,required"`
	Index       int64                                 `json:"index,required"`
	Title       string                                `json:"title,required"`
	Description string                                `json:"description,nullable"`
	JSON        customChartsSectionSubSectionJSON     `json:"-"`
}

// customChartsSectionSubSectionJSON contains the JSON metadata for the struct
// [CustomChartsSectionSubSection]
type customChartsSectionSubSectionJSON struct {
	ID          apijson.Field
	Charts      apijson.Field
	Index       apijson.Field
	Title       apijson.Field
	Description apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSection) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChart struct {
	ID string `json:"id,required" format:"uuid"`
	// Enum for custom chart types.
	ChartType     CustomChartsSectionSubSectionsChartsChartType     `json:"chart_type,required"`
	Data          []CustomChartsSectionSubSectionsChartsData        `json:"data,required"`
	Index         int64                                             `json:"index,required"`
	Series        []CustomChartsSectionSubSectionsChartsSeries      `json:"series,required"`
	Title         string                                            `json:"title,required"`
	CommonFilters CustomChartsSectionSubSectionsChartsCommonFilters `json:"common_filters,nullable"`
	Description   string                                            `json:"description,nullable"`
	Metadata      map[string]interface{}                            `json:"metadata,nullable"`
	JSON          customChartsSectionSubSectionsChartJSON           `json:"-"`
}

// customChartsSectionSubSectionsChartJSON contains the JSON metadata for the
// struct [CustomChartsSectionSubSectionsChart]
type customChartsSectionSubSectionsChartJSON struct {
	ID            apijson.Field
	ChartType     apijson.Field
	Data          apijson.Field
	Index         apijson.Field
	Series        apijson.Field
	Title         apijson.Field
	CommonFilters apijson.Field
	Description   apijson.Field
	Metadata      apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChart) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartJSON) RawJSON() string {
	return r.raw
}

// Enum for custom chart types.
type CustomChartsSectionSubSectionsChartsChartType string

const (
	CustomChartsSectionSubSectionsChartsChartTypeLine CustomChartsSectionSubSectionsChartsChartType = "line"
	CustomChartsSectionSubSectionsChartsChartTypeBar  CustomChartsSectionSubSectionsChartsChartType = "bar"
)

func (r CustomChartsSectionSubSectionsChartsChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsChartTypeLine, CustomChartsSectionSubSectionsChartsChartTypeBar:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsData struct {
	SeriesID  string                                             `json:"series_id,required"`
	Timestamp time.Time                                          `json:"timestamp,required" format:"date-time"`
	Value     CustomChartsSectionSubSectionsChartsDataValueUnion `json:"value,required,nullable"`
	Group     string                                             `json:"group,nullable"`
	JSON      customChartsSectionSubSectionsChartsDataJSON       `json:"-"`
}

// customChartsSectionSubSectionsChartsDataJSON contains the JSON metadata for the
// struct [CustomChartsSectionSubSectionsChartsData]
type customChartsSectionSubSectionsChartsDataJSON struct {
	SeriesID    apijson.Field
	Timestamp   apijson.Field
	Value       apijson.Field
	Group       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsDataJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionFloat] or
// [CustomChartsSectionSubSectionsChartsDataValueMap].
type CustomChartsSectionSubSectionsChartsDataValueUnion interface {
	ImplementsCustomChartsSectionSubSectionsChartsDataValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionSubSectionsChartsDataValueUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsDataValueMap{}),
		},
	)
}

type CustomChartsSectionSubSectionsChartsDataValueMap map[string]interface{}

func (r CustomChartsSectionSubSectionsChartsDataValueMap) ImplementsCustomChartsSectionSubSectionsChartsDataValueUnion() {
}

type CustomChartsSectionSubSectionsChartsSeries struct {
	ID string `json:"id,required" format:"uuid"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric      CustomChartsSectionSubSectionsChartsSeriesMetric  `json:"metric,required"`
	Name        string                                            `json:"name,required"`
	FeedbackKey string                                            `json:"feedback_key,nullable"`
	Filters     CustomChartsSectionSubSectionsChartsSeriesFilters `json:"filters,nullable"`
	// Include additional information about where the group_by param was set.
	GroupBy CustomChartsSectionSubSectionsChartsSeriesGroupBy `json:"group_by,nullable"`
	// LGP Metrics you can chart.
	ProjectMetric CustomChartsSectionSubSectionsChartsSeriesProjectMetric `json:"project_metric,nullable"`
	WorkspaceID   string                                                  `json:"workspace_id,nullable" format:"uuid"`
	JSON          customChartsSectionSubSectionsChartsSeriesJSON          `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesJSON contains the JSON metadata for
// the struct [CustomChartsSectionSubSectionsChartsSeries]
type customChartsSectionSubSectionsChartsSeriesJSON struct {
	ID            apijson.Field
	Metric        apijson.Field
	Name          apijson.Field
	FeedbackKey   apijson.Field
	Filters       apijson.Field
	GroupBy       apijson.Field
	ProjectMetric apijson.Field
	WorkspaceID   apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeries) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesJSON) RawJSON() string {
	return r.raw
}

// Metrics you can chart. Feedback metrics are not available for
// organization-scoped charts.
type CustomChartsSectionSubSectionsChartsSeriesMetric string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricRunCount            CustomChartsSectionSubSectionsChartsSeriesMetric = "run_count"
	CustomChartsSectionSubSectionsChartsSeriesMetricLatencyP50          CustomChartsSectionSubSectionsChartsSeriesMetric = "latency_p50"
	CustomChartsSectionSubSectionsChartsSeriesMetricLatencyP99          CustomChartsSectionSubSectionsChartsSeriesMetric = "latency_p99"
	CustomChartsSectionSubSectionsChartsSeriesMetricLatencyAvg          CustomChartsSectionSubSectionsChartsSeriesMetric = "latency_avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricFirstTokenP50       CustomChartsSectionSubSectionsChartsSeriesMetric = "first_token_p50"
	CustomChartsSectionSubSectionsChartsSeriesMetricFirstTokenP99       CustomChartsSectionSubSectionsChartsSeriesMetric = "first_token_p99"
	CustomChartsSectionSubSectionsChartsSeriesMetricTotalTokens         CustomChartsSectionSubSectionsChartsSeriesMetric = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricPromptTokens        CustomChartsSectionSubSectionsChartsSeriesMetric = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricCompletionTokens    CustomChartsSectionSubSectionsChartsSeriesMetric = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricMedianTokens        CustomChartsSectionSubSectionsChartsSeriesMetric = "median_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricCompletionTokensP50 CustomChartsSectionSubSectionsChartsSeriesMetric = "completion_tokens_p50"
	CustomChartsSectionSubSectionsChartsSeriesMetricPromptTokensP50     CustomChartsSectionSubSectionsChartsSeriesMetric = "prompt_tokens_p50"
	CustomChartsSectionSubSectionsChartsSeriesMetricTokensP99           CustomChartsSectionSubSectionsChartsSeriesMetric = "tokens_p99"
	CustomChartsSectionSubSectionsChartsSeriesMetricCompletionTokensP99 CustomChartsSectionSubSectionsChartsSeriesMetric = "completion_tokens_p99"
	CustomChartsSectionSubSectionsChartsSeriesMetricPromptTokensP99     CustomChartsSectionSubSectionsChartsSeriesMetric = "prompt_tokens_p99"
	CustomChartsSectionSubSectionsChartsSeriesMetricFeedback            CustomChartsSectionSubSectionsChartsSeriesMetric = "feedback"
	CustomChartsSectionSubSectionsChartsSeriesMetricFeedbackScoreAvg    CustomChartsSectionSubSectionsChartsSeriesMetric = "feedback_score_avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricFeedbackValues      CustomChartsSectionSubSectionsChartsSeriesMetric = "feedback_values"
	CustomChartsSectionSubSectionsChartsSeriesMetricTotalCost           CustomChartsSectionSubSectionsChartsSeriesMetric = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricPromptCost          CustomChartsSectionSubSectionsChartsSeriesMetric = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricCompletionCost      CustomChartsSectionSubSectionsChartsSeriesMetric = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricErrorRate           CustomChartsSectionSubSectionsChartsSeriesMetric = "error_rate"
	CustomChartsSectionSubSectionsChartsSeriesMetricStreamingRate       CustomChartsSectionSubSectionsChartsSeriesMetric = "streaming_rate"
	CustomChartsSectionSubSectionsChartsSeriesMetricCostP50             CustomChartsSectionSubSectionsChartsSeriesMetric = "cost_p50"
	CustomChartsSectionSubSectionsChartsSeriesMetricCostP99             CustomChartsSectionSubSectionsChartsSeriesMetric = "cost_p99"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricRunCount, CustomChartsSectionSubSectionsChartsSeriesMetricLatencyP50, CustomChartsSectionSubSectionsChartsSeriesMetricLatencyP99, CustomChartsSectionSubSectionsChartsSeriesMetricLatencyAvg, CustomChartsSectionSubSectionsChartsSeriesMetricFirstTokenP50, CustomChartsSectionSubSectionsChartsSeriesMetricFirstTokenP99, CustomChartsSectionSubSectionsChartsSeriesMetricTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricMedianTokens, CustomChartsSectionSubSectionsChartsSeriesMetricCompletionTokensP50, CustomChartsSectionSubSectionsChartsSeriesMetricPromptTokensP50, CustomChartsSectionSubSectionsChartsSeriesMetricTokensP99, CustomChartsSectionSubSectionsChartsSeriesMetricCompletionTokensP99, CustomChartsSectionSubSectionsChartsSeriesMetricPromptTokensP99, CustomChartsSectionSubSectionsChartsSeriesMetricFeedback, CustomChartsSectionSubSectionsChartsSeriesMetricFeedbackScoreAvg, CustomChartsSectionSubSectionsChartsSeriesMetricFeedbackValues, CustomChartsSectionSubSectionsChartsSeriesMetricTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricErrorRate, CustomChartsSectionSubSectionsChartsSeriesMetricStreamingRate, CustomChartsSectionSubSectionsChartsSeriesMetricCostP50, CustomChartsSectionSubSectionsChartsSeriesMetricCostP99:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesFilters struct {
	Filter      string                                                `json:"filter,nullable"`
	Session     []string                                              `json:"session,nullable" format:"uuid"`
	TraceFilter string                                                `json:"trace_filter,nullable"`
	TreeFilter  string                                                `json:"tree_filter,nullable"`
	JSON        customChartsSectionSubSectionsChartsSeriesFiltersJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesFiltersJSON contains the JSON metadata
// for the struct [CustomChartsSectionSubSectionsChartsSeriesFilters]
type customChartsSectionSubSectionsChartsSeriesFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesFiltersJSON) RawJSON() string {
	return r.raw
}

// Include additional information about where the group_by param was set.
type CustomChartsSectionSubSectionsChartsSeriesGroupBy struct {
	Attribute CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute `json:"attribute,required"`
	MaxGroups int64                                                      `json:"max_groups"`
	Path      string                                                     `json:"path,nullable"`
	SetBy     CustomChartsSectionSubSectionsChartsSeriesGroupBySetBy     `json:"set_by,nullable"`
	JSON      customChartsSectionSubSectionsChartsSeriesGroupByJSON      `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesGroupByJSON contains the JSON metadata
// for the struct [CustomChartsSectionSubSectionsChartsSeriesGroupBy]
type customChartsSectionSubSectionsChartsSeriesGroupByJSON struct {
	Attribute   apijson.Field
	MaxGroups   apijson.Field
	Path        apijson.Field
	SetBy       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesGroupBy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesGroupByJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute string

const (
	CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeName     CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute = "name"
	CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeRunType  CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute = "run_type"
	CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeTag      CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute = "tag"
	CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeMetadata CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute = "metadata"
)

func (r CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeName, CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeRunType, CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeTag, CustomChartsSectionSubSectionsChartsSeriesGroupByAttributeMetadata:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesGroupBySetBy string

const (
	CustomChartsSectionSubSectionsChartsSeriesGroupBySetBySection CustomChartsSectionSubSectionsChartsSeriesGroupBySetBy = "section"
	CustomChartsSectionSubSectionsChartsSeriesGroupBySetBySeries  CustomChartsSectionSubSectionsChartsSeriesGroupBySetBy = "series"
)

func (r CustomChartsSectionSubSectionsChartsSeriesGroupBySetBy) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesGroupBySetBySection, CustomChartsSectionSubSectionsChartsSeriesGroupBySetBySeries:
		return true
	}
	return false
}

// LGP Metrics you can chart.
type CustomChartsSectionSubSectionsChartsSeriesProjectMetric string

const (
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricMemoryUsage             CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "memory_usage"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricCPUUsage                CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "cpu_usage"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricDiskUsage               CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "disk_usage"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricRestartCount            CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "restart_count"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricReplicaCount            CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "replica_count"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricWorkerCount             CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "worker_count"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricLgRunCount              CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "lg_run_count"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricResponsesPerSecond      CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "responses_per_second"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricErrorResponsesPerSecond CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "error_responses_per_second"
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricP95Latency              CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "p95_latency"
)

func (r CustomChartsSectionSubSectionsChartsSeriesProjectMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesProjectMetricMemoryUsage, CustomChartsSectionSubSectionsChartsSeriesProjectMetricCPUUsage, CustomChartsSectionSubSectionsChartsSeriesProjectMetricDiskUsage, CustomChartsSectionSubSectionsChartsSeriesProjectMetricRestartCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricReplicaCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricWorkerCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricLgRunCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricResponsesPerSecond, CustomChartsSectionSubSectionsChartsSeriesProjectMetricErrorResponsesPerSecond, CustomChartsSectionSubSectionsChartsSeriesProjectMetricP95Latency:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsCommonFilters struct {
	Filter      string                                                `json:"filter,nullable"`
	Session     []string                                              `json:"session,nullable" format:"uuid"`
	TraceFilter string                                                `json:"trace_filter,nullable"`
	TreeFilter  string                                                `json:"tree_filter,nullable"`
	JSON        customChartsSectionSubSectionsChartsCommonFiltersJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsCommonFiltersJSON contains the JSON metadata
// for the struct [CustomChartsSectionSubSectionsChartsCommonFilters]
type customChartsSectionSubSectionsChartsCommonFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsCommonFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsCommonFiltersJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionRequestParam struct {
	EndTime param.Field[time.Time] `json:"end_time" format:"date-time"`
	// Group by param for run stats.
	GroupBy   param.Field[RunStatsGroupByParam] `json:"group_by"`
	OmitData  param.Field[bool]                 `json:"omit_data"`
	StartTime param.Field[time.Time]            `json:"start_time" format:"date-time"`
	// Timedelta input.
	Stride   param.Field[TimedeltaInputParam] `json:"stride"`
	Timezone param.Field[string]              `json:"timezone"`
}

func (r CustomChartsSectionRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Group by param for run stats.
type RunStatsGroupByParam struct {
	Attribute param.Field[RunStatsGroupByAttribute] `json:"attribute,required"`
	MaxGroups param.Field[int64]                    `json:"max_groups"`
	Path      param.Field[string]                   `json:"path"`
}

func (r RunStatsGroupByParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RunStatsGroupByAttribute string

const (
	RunStatsGroupByAttributeName     RunStatsGroupByAttribute = "name"
	RunStatsGroupByAttributeRunType  RunStatsGroupByAttribute = "run_type"
	RunStatsGroupByAttributeTag      RunStatsGroupByAttribute = "tag"
	RunStatsGroupByAttributeMetadata RunStatsGroupByAttribute = "metadata"
)

func (r RunStatsGroupByAttribute) IsKnown() bool {
	switch r {
	case RunStatsGroupByAttributeName, RunStatsGroupByAttributeRunType, RunStatsGroupByAttributeTag, RunStatsGroupByAttributeMetadata:
		return true
	}
	return false
}

type SessionSortableColumns string

const (
	SessionSortableColumnsName             SessionSortableColumns = "name"
	SessionSortableColumnsStartTime        SessionSortableColumns = "start_time"
	SessionSortableColumnsLastRunStartTime SessionSortableColumns = "last_run_start_time"
	SessionSortableColumnsLatencyP50       SessionSortableColumns = "latency_p50"
	SessionSortableColumnsLatencyP99       SessionSortableColumns = "latency_p99"
	SessionSortableColumnsErrorRate        SessionSortableColumns = "error_rate"
	SessionSortableColumnsFeedback         SessionSortableColumns = "feedback"
	SessionSortableColumnsRunsCount        SessionSortableColumns = "runs_count"
)

func (r SessionSortableColumns) IsKnown() bool {
	switch r {
	case SessionSortableColumnsName, SessionSortableColumnsStartTime, SessionSortableColumnsLastRunStartTime, SessionSortableColumnsLatencyP50, SessionSortableColumnsLatencyP99, SessionSortableColumnsErrorRate, SessionSortableColumnsFeedback, SessionSortableColumnsRunsCount:
		return true
	}
	return false
}

// Timedelta input.
type TimedeltaInputParam struct {
	Days    param.Field[int64] `json:"days"`
	Hours   param.Field[int64] `json:"hours"`
	Minutes param.Field[int64] `json:"minutes"`
}

func (r TimedeltaInputParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// TracerSession schema.
type TracerSession struct {
	ID                   string                   `json:"id,required" format:"uuid"`
	TenantID             string                   `json:"tenant_id,required" format:"uuid"`
	CompletionCost       string                   `json:"completion_cost,nullable"`
	CompletionTokens     int64                    `json:"completion_tokens,nullable"`
	DefaultDatasetID     string                   `json:"default_dataset_id,nullable" format:"uuid"`
	Description          string                   `json:"description,nullable"`
	EndTime              time.Time                `json:"end_time,nullable" format:"date-time"`
	ErrorRate            float64                  `json:"error_rate,nullable"`
	Extra                map[string]interface{}   `json:"extra,nullable"`
	FeedbackStats        map[string]interface{}   `json:"feedback_stats,nullable"`
	FirstTokenP50        float64                  `json:"first_token_p50,nullable"`
	FirstTokenP99        float64                  `json:"first_token_p99,nullable"`
	LastRunStartTime     time.Time                `json:"last_run_start_time,nullable" format:"date-time"`
	LastRunStartTimeLive time.Time                `json:"last_run_start_time_live,nullable" format:"date-time"`
	LatencyP50           float64                  `json:"latency_p50,nullable"`
	LatencyP99           float64                  `json:"latency_p99,nullable"`
	Name                 string                   `json:"name"`
	PromptCost           string                   `json:"prompt_cost,nullable"`
	PromptTokens         int64                    `json:"prompt_tokens,nullable"`
	ReferenceDatasetID   string                   `json:"reference_dataset_id,nullable" format:"uuid"`
	RunCount             int64                    `json:"run_count,nullable"`
	RunFacets            []map[string]interface{} `json:"run_facets,nullable"`
	SessionFeedbackStats map[string]interface{}   `json:"session_feedback_stats,nullable"`
	StartTime            time.Time                `json:"start_time" format:"date-time"`
	StreamingRate        float64                  `json:"streaming_rate,nullable"`
	TestRunNumber        int64                    `json:"test_run_number,nullable"`
	TotalCost            string                   `json:"total_cost,nullable"`
	TotalTokens          int64                    `json:"total_tokens,nullable"`
	TraceTier            TracerSessionTraceTier   `json:"trace_tier,nullable"`
	JSON                 tracerSessionJSON        `json:"-"`
}

// tracerSessionJSON contains the JSON metadata for the struct [TracerSession]
type tracerSessionJSON struct {
	ID                   apijson.Field
	TenantID             apijson.Field
	CompletionCost       apijson.Field
	CompletionTokens     apijson.Field
	DefaultDatasetID     apijson.Field
	Description          apijson.Field
	EndTime              apijson.Field
	ErrorRate            apijson.Field
	Extra                apijson.Field
	FeedbackStats        apijson.Field
	FirstTokenP50        apijson.Field
	FirstTokenP99        apijson.Field
	LastRunStartTime     apijson.Field
	LastRunStartTimeLive apijson.Field
	LatencyP50           apijson.Field
	LatencyP99           apijson.Field
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

func (r *TracerSession) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r tracerSessionJSON) RawJSON() string {
	return r.raw
}

type TracerSessionTraceTier string

const (
	TracerSessionTraceTierLonglived  TracerSessionTraceTier = "longlived"
	TracerSessionTraceTierShortlived TracerSessionTraceTier = "shortlived"
)

func (r TracerSessionTraceTier) IsKnown() bool {
	switch r {
	case TracerSessionTraceTierLonglived, TracerSessionTraceTierShortlived:
		return true
	}
	return false
}

// TracerSession schema.
type TracerSessionWithoutVirtualFields struct {
	ID                   string                                     `json:"id,required" format:"uuid"`
	TenantID             string                                     `json:"tenant_id,required" format:"uuid"`
	DefaultDatasetID     string                                     `json:"default_dataset_id,nullable" format:"uuid"`
	Description          string                                     `json:"description,nullable"`
	EndTime              time.Time                                  `json:"end_time,nullable" format:"date-time"`
	Extra                map[string]interface{}                     `json:"extra,nullable"`
	LastRunStartTimeLive time.Time                                  `json:"last_run_start_time_live,nullable" format:"date-time"`
	Name                 string                                     `json:"name"`
	ReferenceDatasetID   string                                     `json:"reference_dataset_id,nullable" format:"uuid"`
	StartTime            time.Time                                  `json:"start_time" format:"date-time"`
	TraceTier            TracerSessionWithoutVirtualFieldsTraceTier `json:"trace_tier,nullable"`
	JSON                 tracerSessionWithoutVirtualFieldsJSON      `json:"-"`
}

// tracerSessionWithoutVirtualFieldsJSON contains the JSON metadata for the struct
// [TracerSessionWithoutVirtualFields]
type tracerSessionWithoutVirtualFieldsJSON struct {
	ID                   apijson.Field
	TenantID             apijson.Field
	DefaultDatasetID     apijson.Field
	Description          apijson.Field
	EndTime              apijson.Field
	Extra                apijson.Field
	LastRunStartTimeLive apijson.Field
	Name                 apijson.Field
	ReferenceDatasetID   apijson.Field
	StartTime            apijson.Field
	TraceTier            apijson.Field
	raw                  string
	ExtraFields          map[string]apijson.Field
}

func (r *TracerSessionWithoutVirtualFields) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r tracerSessionWithoutVirtualFieldsJSON) RawJSON() string {
	return r.raw
}

type TracerSessionWithoutVirtualFieldsTraceTier string

const (
	TracerSessionWithoutVirtualFieldsTraceTierLonglived  TracerSessionWithoutVirtualFieldsTraceTier = "longlived"
	TracerSessionWithoutVirtualFieldsTraceTierShortlived TracerSessionWithoutVirtualFieldsTraceTier = "shortlived"
)

func (r TracerSessionWithoutVirtualFieldsTraceTier) IsKnown() bool {
	switch r {
	case TracerSessionWithoutVirtualFieldsTraceTierLonglived, TracerSessionWithoutVirtualFieldsTraceTierShortlived:
		return true
	}
	return false
}

type SessionDeleteResponse = interface{}

type SessionNewParams struct {
	Upsert             param.Field[bool]                      `query:"upsert"`
	ID                 param.Field[string]                    `json:"id" format:"uuid"`
	DefaultDatasetID   param.Field[string]                    `json:"default_dataset_id" format:"uuid"`
	Description        param.Field[string]                    `json:"description"`
	EndTime            param.Field[time.Time]                 `json:"end_time" format:"date-time"`
	Extra              param.Field[map[string]interface{}]    `json:"extra"`
	Name               param.Field[string]                    `json:"name"`
	ReferenceDatasetID param.Field[string]                    `json:"reference_dataset_id" format:"uuid"`
	StartTime          param.Field[time.Time]                 `json:"start_time" format:"date-time"`
	TraceTier          param.Field[SessionNewParamsTraceTier] `json:"trace_tier"`
}

func (r SessionNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionNewParams]'s query parameters as `url.Values`.
func (r SessionNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionNewParamsTraceTier string

const (
	SessionNewParamsTraceTierLonglived  SessionNewParamsTraceTier = "longlived"
	SessionNewParamsTraceTierShortlived SessionNewParamsTraceTier = "shortlived"
)

func (r SessionNewParamsTraceTier) IsKnown() bool {
	switch r {
	case SessionNewParamsTraceTierLonglived, SessionNewParamsTraceTierShortlived:
		return true
	}
	return false
}

type SessionGetParams struct {
	IncludeStats   param.Field[bool]      `query:"include_stats"`
	StatsStartTime param.Field[time.Time] `query:"stats_start_time" format:"date-time"`
	Accept         param.Field[string]    `header:"accept"`
}

// URLQuery serializes [SessionGetParams]'s query parameters as `url.Values`.
func (r SessionGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionUpdateParams struct {
	DefaultDatasetID param.Field[string]                       `json:"default_dataset_id" format:"uuid"`
	Description      param.Field[string]                       `json:"description"`
	EndTime          param.Field[time.Time]                    `json:"end_time" format:"date-time"`
	Extra            param.Field[map[string]interface{}]       `json:"extra"`
	Name             param.Field[string]                       `json:"name"`
	TraceTier        param.Field[SessionUpdateParamsTraceTier] `json:"trace_tier"`
}

func (r SessionUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionUpdateParamsTraceTier string

const (
	SessionUpdateParamsTraceTierLonglived  SessionUpdateParamsTraceTier = "longlived"
	SessionUpdateParamsTraceTierShortlived SessionUpdateParamsTraceTier = "shortlived"
)

func (r SessionUpdateParamsTraceTier) IsKnown() bool {
	switch r {
	case SessionUpdateParamsTraceTierLonglived, SessionUpdateParamsTraceTierShortlived:
		return true
	}
	return false
}

type SessionListParams struct {
	ID                param.Field[[]string]               `query:"id" format:"uuid"`
	DatasetVersion    param.Field[string]                 `query:"dataset_version"`
	Facets            param.Field[bool]                   `query:"facets"`
	Filter            param.Field[string]                 `query:"filter"`
	IncludeStats      param.Field[bool]                   `query:"include_stats"`
	Limit             param.Field[int64]                  `query:"limit"`
	Metadata          param.Field[string]                 `query:"metadata"`
	Name              param.Field[string]                 `query:"name"`
	NameContains      param.Field[string]                 `query:"name_contains"`
	Offset            param.Field[int64]                  `query:"offset"`
	ReferenceDataset  param.Field[[]string]               `query:"reference_dataset" format:"uuid"`
	ReferenceFree     param.Field[bool]                   `query:"reference_free"`
	SortBy            param.Field[SessionSortableColumns] `query:"sort_by"`
	SortByDesc        param.Field[bool]                   `query:"sort_by_desc"`
	SortByFeedbackKey param.Field[string]                 `query:"sort_by_feedback_key"`
	StatsStartTime    param.Field[time.Time]              `query:"stats_start_time" format:"date-time"`
	TagValueID        param.Field[[]string]               `query:"tag_value_id" format:"uuid"`
	UseApproxStats    param.Field[bool]                   `query:"use_approx_stats"`
	Accept            param.Field[string]                 `header:"accept"`
}

// URLQuery serializes [SessionListParams]'s query parameters as `url.Values`.
func (r SessionListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionDashboardParams struct {
	CustomChartsSectionRequest CustomChartsSectionRequestParam `json:"custom_charts_section_request,required"`
	Accept                     param.Field[string]             `header:"accept"`
}

func (r SessionDashboardParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.CustomChartsSectionRequest)
}
