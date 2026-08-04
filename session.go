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

// Create a new project.
func (r *SessionService) New(ctx context.Context, params SessionNewParams, opts ...option.RequestOption) (res *TracerSessionWithoutVirtualFields, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/sessions"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Get a specific project.
func (r *SessionService) Get(ctx context.Context, sessionID string, params SessionGetParams, opts ...option.RequestOption) (res *TracerSession, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// Update a project.
func (r *SessionService) Update(ctx context.Context, sessionID string, body SessionUpdateParams, opts ...option.RequestOption) (res *TracerSessionWithoutVirtualFields, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List all projects.
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

// List all projects.
func (r *SessionService) ListAutoPaging(ctx context.Context, params SessionListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[TracerSession] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, params, opts...))
}

// Delete a specific project.
func (r *SessionService) Delete(ctx context.Context, sessionID string, opts ...option.RequestOption) (res *SessionDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Get a prebuilt dashboard for a tracing project.
func (r *SessionService) Dashboard(ctx context.Context, sessionID string, params SessionDashboardParams, opts ...option.RequestOption) (res *CustomChartsSection, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/dashboard", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type CustomChartsSection struct {
	ID          string                          `json:"id" api:"required" format:"uuid"`
	Charts      []CustomChartsSectionChart      `json:"charts" api:"required"`
	Title       string                          `json:"title" api:"required"`
	Description string                          `json:"description" api:"nullable"`
	Index       int64                           `json:"index" api:"nullable"`
	Layout      CustomChartsSectionLayout       `json:"layout" api:"nullable"`
	SessionID   string                          `json:"session_id" api:"nullable" format:"uuid"`
	SubSections []CustomChartsSectionSubSection `json:"sub_sections" api:"nullable"`
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
	Layout      apijson.Field
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
	ID        string                             `json:"id" api:"required" format:"uuid"`
	ChartType CustomChartsSectionChartsChartType `json:"chart_type" api:"required"`
	Index     int64                              `json:"index" api:"required"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedCommonFilters].
	CommonFilters interface{} `json:"common_filters"`
	// This field can have the runtime type of
	// [[]CustomChartsSectionChartsSingleCustomChartResponseSerializedData].
	Data        interface{} `json:"data"`
	Description string      `json:"description" api:"nullable"`
	Markdown    string      `json:"markdown"`
	// This field can have the runtime type of [map[string]interface{}].
	Metadata interface{} `json:"metadata"`
	// This field can have the runtime type of
	// [[]CustomChartsSectionChartsSingleCustomChartResponseSerializedSeries].
	Series interface{}                  `json:"series"`
	Title  string                       `json:"title"`
	JSON   customChartsSectionChartJSON `json:"-"`
	union  CustomChartsSectionChartsUnion
}

// customChartsSectionChartJSON contains the JSON metadata for the struct
// [CustomChartsSectionChart]
type customChartsSectionChartJSON struct {
	ID            apijson.Field
	ChartType     apijson.Field
	Index         apijson.Field
	CommonFilters apijson.Field
	Data          apijson.Field
	Description   apijson.Field
	Markdown      apijson.Field
	Metadata      apijson.Field
	Series        apijson.Field
	Title         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r customChartsSectionChartJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChart) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChart{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [CustomChartsSectionChartsUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSerialized],
// [CustomChartsSectionChartsCustomTextBlock].
func (r CustomChartsSectionChart) AsUnion() CustomChartsSectionChartsUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSerialized] or
// [CustomChartsSectionChartsCustomTextBlock].
type CustomChartsSectionChartsUnion interface {
	implementsCustomChartsSectionChart()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsUnion)(nil)).Elem(),
		"chart_type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerialized{}),
			DiscriminatorValue: "line",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerialized{}),
			DiscriminatorValue: "bar",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerialized{}),
			DiscriminatorValue: "table",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerialized{}),
			DiscriminatorValue: "kpi",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerialized{}),
			DiscriminatorValue: "top-k",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerialized{}),
			DiscriminatorValue: "pie",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsCustomTextBlock{}),
			DiscriminatorValue: "text",
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerialized struct {
	ID            string                                                                    `json:"id" api:"required" format:"uuid"`
	ChartType     CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType     `json:"chart_type" api:"required"`
	Data          []CustomChartsSectionChartsSingleCustomChartResponseSerializedData        `json:"data" api:"required"`
	Index         int64                                                                     `json:"index" api:"required"`
	Series        []CustomChartsSectionChartsSingleCustomChartResponseSerializedSeries      `json:"series" api:"required"`
	Title         string                                                                    `json:"title" api:"required"`
	CommonFilters CustomChartsSectionChartsSingleCustomChartResponseSerializedCommonFilters `json:"common_filters" api:"nullable"`
	Description   string                                                                    `json:"description" api:"nullable"`
	Metadata      map[string]interface{}                                                    `json:"metadata" api:"nullable"`
	JSON          customChartsSectionChartsSingleCustomChartResponseSerializedJSON          `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedJSON contains the
// JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerialized]
type customChartsSectionChartsSingleCustomChartResponseSerializedJSON struct {
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

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerialized) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerialized) implementsCustomChartsSectionChart() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeLine  CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType = "line"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeBar   CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType = "bar"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeTable CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType = "table"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeKpi   CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType = "kpi"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeTopK  CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType = "top-k"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypePie   CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType = "pie"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeLine, CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeBar, CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeTable, CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeKpi, CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypeTopK, CustomChartsSectionChartsSingleCustomChartResponseSerializedChartTypePie:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedData struct {
	SeriesID  string                                                                     `json:"series_id" api:"required"`
	Timestamp time.Time                                                                  `json:"timestamp" api:"required" format:"date-time"`
	Value     CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueUnion `json:"value" api:"required,nullable"`
	Group     string                                                                     `json:"group" api:"nullable"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSerializedDataJSON       `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedDataJSON contains
// the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedData]
type customChartsSectionChartsSingleCustomChartResponseSerializedDataJSON struct {
	SeriesID    apijson.Field
	Timestamp   apijson.Field
	Value       apijson.Field
	Group       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedDataJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionFloat] or
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueMap].
type CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueUnion interface {
	ImplementsCustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueMap{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueMap map[string]interface{}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueMap) ImplementsCustomChartsSectionChartsSingleCustomChartResponseSerializedDataValueUnion() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeries struct {
	ID               string                                                                             `json:"id" api:"required" format:"uuid"`
	Name             string                                                                             `json:"name" api:"required"`
	FeedbackKey      string                                                                             `json:"feedback_key" api:"nullable"`
	FilterDefinition CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition `json:"filter_definition" api:"nullable"`
	Filters          CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilters          `json:"filters" api:"nullable"`
	// Include additional information about where the group_by param was set.
	GroupBy            CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBy             `json:"group_by" api:"nullable"`
	GroupByDefinitions []CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition `json:"group_by_definitions" api:"nullable"`
	Metadata           map[string]interface{}                                                                `json:"metadata" api:"nullable"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric           CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric           `json:"metric" api:"nullable"`
	MetricDefinition CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition `json:"metric_definition" api:"nullable"`
	// LGP Metrics you can chart.
	ProjectMetric CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric `json:"project_metric" api:"nullable"`
	WorkspaceID   string                                                                          `json:"workspace_id" api:"nullable" format:"uuid"`
	JSON          customChartsSectionChartsSingleCustomChartResponseSerializedSeriesJSON          `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesJSON contains
// the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeries]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesJSON struct {
	ID                 apijson.Field
	Name               apijson.Field
	FeedbackKey        apijson.Field
	FilterDefinition   apijson.Field
	Filters            apijson.Field
	GroupBy            apijson.Field
	GroupByDefinitions apijson.Field
	Metadata           apijson.Field
	Metric             apijson.Field
	MetricDefinition   apijson.Field
	ProjectMetric      apijson.Field
	WorkspaceID        apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeries) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition struct {
	SourceType CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceType `json:"source_type" api:"required"`
	// This field can have the runtime type of [[]string].
	DatasetIDs interface{} `json:"dataset_ids"`
	// This field can have the runtime type of [[]string].
	ProjectIDs  interface{}                                                                            `json:"project_ids"`
	RunFilter   string                                                                                 `json:"run_filter" api:"nullable"`
	TraceFilter string                                                                                 `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                                 `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionJSON `json:"-"`
	union       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionUnion
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionJSON struct {
	SourceType  apijson.Field
	DatasetIDs  apijson.Field
	ProjectIDs  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset].
func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset].
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject struct {
	ProjectIDs  []string                                                                                                                      `json:"project_ids" api:"required" format:"uuid"`
	SourceType  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType `json:"source_type" api:"required"`
	RunFilter   string                                                                                                                        `json:"run_filter" api:"nullable"`
	TraceFilter string                                                                                                                        `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                                                                        `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON       `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON struct {
	ProjectIDs  apijson.Field
	SourceType  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProject) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType = "tracing_project"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset struct {
	DatasetIDs []string                                                                                                               `json:"dataset_ids" api:"required" format:"uuid"`
	SourceType CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetSourceType `json:"source_type" api:"required"`
	JSON       customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetJSON       `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetJSON struct {
	DatasetIDs  apijson.Field
	SourceType  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDataset) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetSourceType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetSourceType = "dataset"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceTypeTracingProject CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceType = "tracing_project"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceTypeDataset        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceType = "dataset"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceTypeTracingProject, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilterDefinitionSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilters struct {
	Filter      string                                                                        `json:"filter" api:"nullable"`
	Session     []string                                                                      `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                                                        `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                        `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFiltersJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFiltersJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilters]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesFiltersJSON) RawJSON() string {
	return r.raw
}

// Include additional information about where the group_by param was set.
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBy struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute `json:"attribute" api:"required"`
	MaxGroups int64                                                                              `json:"max_groups"`
	Path      string                                                                             `json:"path" api:"nullable"`
	SetBy     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBy     `json:"set_by" api:"nullable"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByJSON      `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBy]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByJSON struct {
	Attribute   apijson.Field
	MaxGroups   apijson.Field
	Path        apijson.Field
	SetBy       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeName     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute = "name"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeRunType  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute = "run_type"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeTag      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute = "tag"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeMetadata CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute = "metadata"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeName, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeRunType, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeTag, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByAttributeMetadata:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBy string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBySection CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBy = "section"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBySeries  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBy = "series"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBy) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBySection, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupBySetBySeries:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute `json:"attribute" api:"required"`
	Path      string                                                                                        `json:"path"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionJSON       `json:"-"`
	union     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsUnion
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex].
func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex].
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute `json:"attribute" api:"required"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainJSON      `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainJSON struct {
	Attribute   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlain) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "name"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "run_type"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "tag"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "project"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "status"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute `json:"attribute" api:"required"`
	Path      string                                                                                                                 `json:"path" api:"required"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexJSON      `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplex) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "metadata"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "feedback_label"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeName          CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "name"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeRunType       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "run_type"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeTag           CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "tag"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeProject       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "project"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeStatus        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "status"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeMetadata      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "metadata"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeFeedbackLabel CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute = "feedback_label"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeName, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeRunType, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeTag, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeProject, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeStatus, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeMetadata, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesGroupByDefinitionsAttributeFeedbackLabel:
		return true
	}
	return false
}

// Metrics you can chart. Feedback metrics are not available for
// organization-scoped charts.
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricRunCount            CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "run_count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricLatencyP50          CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "latency_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricLatencyP99          CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "latency_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricLatencyAvg          CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "latency_avg"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFirstTokenP50       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "first_token_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFirstTokenP99       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "first_token_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricTotalTokens         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptTokens        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionTokens    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricMedianTokens        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "median_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionTokensP50 CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "completion_tokens_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptTokensP50     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "prompt_tokens_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricTokensP99           CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "tokens_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionTokensP99 CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "completion_tokens_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptTokensP99     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "prompt_tokens_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFeedback            CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "feedback"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFeedbackScoreAvg    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "feedback_score_avg"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFeedbackValues      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "feedback_values"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricTotalCost           CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptCost          CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionCost      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricErrorRate           CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "error_rate"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricStreamingRate       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "streaming_rate"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCostP50             CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "cost_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCostP99             CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric = "cost_p99"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricRunCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricLatencyP50, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricLatencyP99, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricLatencyAvg, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFirstTokenP50, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFirstTokenP99, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricMedianTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionTokensP50, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptTokensP50, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricTokensP99, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionTokensP99, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptTokensP99, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFeedback, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFeedbackScoreAvg, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricFeedbackValues, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricErrorRate, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricStreamingRate, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCostP50, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricCostP99:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition struct {
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator].
	Denominator interface{}                                                                              `json:"denominator"`
	Entity      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionEntity `json:"entity"`
	Field       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField  `json:"field"`
	Filter      string                                                                                   `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator].
	Numerator interface{} `json:"numerator"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParams].
	Params interface{}                                                                            `json:"params"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionJSON `json:"-"`
	union  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionUnion
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionJSON struct {
	Denominator apijson.Field
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Numerator   apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput].
func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput].
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                 `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetric) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                     `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount struct {
	Filter string                                                                                                       `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCount) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountTypeCount CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                       `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                           `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                         `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeSum CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                           `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                               `json:"feedback_key" api:"required"`
	P           float64                                                                                                                              `json:"p" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                              `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                 `json:"p" api:"required"`
	JSON customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput struct {
	Denominator CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator `json:"denominator" api:"required"`
	Numerator   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator   `json:"numerator" api:"required"`
	Type        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputType        `json:"type" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputJSON        `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputJSON struct {
	Denominator apijson.Field
	Numerator   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutput) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity `json:"entity"`
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField  `json:"field"`
	Filter string                                                                                                                          `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams].
	Params interface{}                                                                                                                   `json:"params"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON `json:"-"`
	union  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON struct {
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                                                        `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                                                            `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricTypeCount CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount struct {
	Filter string                                                                                                                                              `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                                                              `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                                                                  `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                                                                `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                                                  `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                                                                      `json:"feedback_key" api:"required"`
	P           float64                                                                                                                                                                     `json:"p" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                                     `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                                                        `json:"p" api:"required"`
	JSON customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "feedback_score"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_cost"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFeedbackScore, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "avg"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity `json:"entity"`
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField  `json:"field"`
	Filter string                                                                                                                        `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams].
	Params interface{}                                                                                                                 `json:"params"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON `json:"-"`
	union  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON struct {
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                                                      `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                                                          `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricTypeCount CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount struct {
	Filter string                                                                                                                                            `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                                                            `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                                                                `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                                                              `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                                                `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                                                                    `json:"feedback_key" api:"required"`
	P           float64                                                                                                                                                                   `json:"p" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                                   `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                                                      `json:"p" api:"required"`
	JSON customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "feedback_score"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_cost"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFeedbackScore, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "avg"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputType = "ratio"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "feedback_score"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField = "completion_cost"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldFeedbackScore, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeCount      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeMax        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeMin        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeAvg        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "avg"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeSum        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "percentile"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeRatio      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType = "ratio"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeAvg, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypePercentile, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesMetricDefinitionTypeRatio:
		return true
	}
	return false
}

// LGP Metrics you can chart.
type CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricMemoryUsage             CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "memory_usage"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricCPUUsage                CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "cpu_usage"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricDiskUsage               CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "disk_usage"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricRestartCount            CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "restart_count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricReplicaCount            CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "replica_count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricWorkerCount             CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "worker_count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricLgRunCount              CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "lg_run_count"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricResponsesPerSecond      CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "responses_per_second"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricErrorResponsesPerSecond CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "error_responses_per_second"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricP95Latency              CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "p95_latency"
	CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricRunQueueWaitTime        CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric = "run_queue_wait_time"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricMemoryUsage, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricCPUUsage, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricDiskUsage, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricRestartCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricReplicaCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricWorkerCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricLgRunCount, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricResponsesPerSecond, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricErrorResponsesPerSecond, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricP95Latency, CustomChartsSectionChartsSingleCustomChartResponseSerializedSeriesProjectMetricRunQueueWaitTime:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSerializedCommonFilters struct {
	Filter      string                                                                        `json:"filter" api:"nullable"`
	Session     []string                                                                      `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                                                        `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                        `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSerializedCommonFiltersJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSerializedCommonFiltersJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSerializedCommonFilters]
type customChartsSectionChartsSingleCustomChartResponseSerializedCommonFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSerializedCommonFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSerializedCommonFiltersJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsCustomTextBlock struct {
	ID        string                                            `json:"id" api:"required" format:"uuid"`
	ChartType CustomChartsSectionChartsCustomTextBlockChartType `json:"chart_type" api:"required"`
	Index     int64                                             `json:"index" api:"required"`
	Markdown  string                                            `json:"markdown" api:"required"`
	Metadata  map[string]interface{}                            `json:"metadata" api:"nullable"`
	JSON      customChartsSectionChartsCustomTextBlockJSON      `json:"-"`
}

// customChartsSectionChartsCustomTextBlockJSON contains the JSON metadata for the
// struct [CustomChartsSectionChartsCustomTextBlock]
type customChartsSectionChartsCustomTextBlockJSON struct {
	ID          apijson.Field
	ChartType   apijson.Field
	Index       apijson.Field
	Markdown    apijson.Field
	Metadata    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsCustomTextBlock) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsCustomTextBlockJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsCustomTextBlock) implementsCustomChartsSectionChart() {}

type CustomChartsSectionChartsCustomTextBlockChartType string

const (
	CustomChartsSectionChartsCustomTextBlockChartTypeText CustomChartsSectionChartsCustomTextBlockChartType = "text"
)

func (r CustomChartsSectionChartsCustomTextBlockChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsCustomTextBlockChartTypeText:
		return true
	}
	return false
}

type CustomChartsSectionChartsChartType string

const (
	CustomChartsSectionChartsChartTypeLine  CustomChartsSectionChartsChartType = "line"
	CustomChartsSectionChartsChartTypeBar   CustomChartsSectionChartsChartType = "bar"
	CustomChartsSectionChartsChartTypeTable CustomChartsSectionChartsChartType = "table"
	CustomChartsSectionChartsChartTypeKpi   CustomChartsSectionChartsChartType = "kpi"
	CustomChartsSectionChartsChartTypeTopK  CustomChartsSectionChartsChartType = "top-k"
	CustomChartsSectionChartsChartTypePie   CustomChartsSectionChartsChartType = "pie"
	CustomChartsSectionChartsChartTypeText  CustomChartsSectionChartsChartType = "text"
)

func (r CustomChartsSectionChartsChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsChartTypeLine, CustomChartsSectionChartsChartTypeBar, CustomChartsSectionChartsChartTypeTable, CustomChartsSectionChartsChartTypeKpi, CustomChartsSectionChartsChartTypeTopK, CustomChartsSectionChartsChartTypePie, CustomChartsSectionChartsChartTypeText:
		return true
	}
	return false
}

type CustomChartsSectionLayout struct {
	Breakpoints CustomChartsSectionLayoutBreakpoints `json:"breakpoints" api:"required"`
	Version     CustomChartsSectionLayoutVersion     `json:"version" api:"required"`
	JSON        customChartsSectionLayoutJSON        `json:"-"`
}

// customChartsSectionLayoutJSON contains the JSON metadata for the struct
// [CustomChartsSectionLayout]
type customChartsSectionLayoutJSON struct {
	Breakpoints apijson.Field
	Version     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayout) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpoints struct {
	Md   CustomChartsSectionLayoutBreakpointsMd   `json:"md" api:"required"`
	Sm   CustomChartsSectionLayoutBreakpointsSm   `json:"sm" api:"required"`
	JSON customChartsSectionLayoutBreakpointsJSON `json:"-"`
}

// customChartsSectionLayoutBreakpointsJSON contains the JSON metadata for the
// struct [CustomChartsSectionLayoutBreakpoints]
type customChartsSectionLayoutBreakpointsJSON struct {
	Md          apijson.Field
	Sm          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpoints) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpointsMd struct {
	Rows []CustomChartsSectionLayoutBreakpointsMdRow `json:"rows" api:"required"`
	JSON customChartsSectionLayoutBreakpointsMdJSON  `json:"-"`
}

// customChartsSectionLayoutBreakpointsMdJSON contains the JSON metadata for the
// struct [CustomChartsSectionLayoutBreakpointsMd]
type customChartsSectionLayoutBreakpointsMdJSON struct {
	Rows        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpointsMd) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsMdJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpointsMdRow struct {
	HeightUnits int64                                            `json:"height_units" api:"required"`
	Items       []CustomChartsSectionLayoutBreakpointsMdRowsItem `json:"items" api:"required"`
	JSON        customChartsSectionLayoutBreakpointsMdRowJSON    `json:"-"`
}

// customChartsSectionLayoutBreakpointsMdRowJSON contains the JSON metadata for the
// struct [CustomChartsSectionLayoutBreakpointsMdRow]
type customChartsSectionLayoutBreakpointsMdRowJSON struct {
	HeightUnits apijson.Field
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpointsMdRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsMdRowJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpointsMdRowsItem struct {
	ChartID    string                                             `json:"chart_id" api:"required" format:"uuid"`
	WidthUnits int64                                              `json:"width_units" api:"required"`
	JSON       customChartsSectionLayoutBreakpointsMdRowsItemJSON `json:"-"`
}

// customChartsSectionLayoutBreakpointsMdRowsItemJSON contains the JSON metadata
// for the struct [CustomChartsSectionLayoutBreakpointsMdRowsItem]
type customChartsSectionLayoutBreakpointsMdRowsItemJSON struct {
	ChartID     apijson.Field
	WidthUnits  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpointsMdRowsItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsMdRowsItemJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpointsSm struct {
	Rows []CustomChartsSectionLayoutBreakpointsSmRow `json:"rows" api:"required"`
	JSON customChartsSectionLayoutBreakpointsSmJSON  `json:"-"`
}

// customChartsSectionLayoutBreakpointsSmJSON contains the JSON metadata for the
// struct [CustomChartsSectionLayoutBreakpointsSm]
type customChartsSectionLayoutBreakpointsSmJSON struct {
	Rows        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpointsSm) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsSmJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpointsSmRow struct {
	HeightUnits int64                                            `json:"height_units" api:"required"`
	Items       []CustomChartsSectionLayoutBreakpointsSmRowsItem `json:"items" api:"required"`
	JSON        customChartsSectionLayoutBreakpointsSmRowJSON    `json:"-"`
}

// customChartsSectionLayoutBreakpointsSmRowJSON contains the JSON metadata for the
// struct [CustomChartsSectionLayoutBreakpointsSmRow]
type customChartsSectionLayoutBreakpointsSmRowJSON struct {
	HeightUnits apijson.Field
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpointsSmRow) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsSmRowJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutBreakpointsSmRowsItem struct {
	ChartID    string                                             `json:"chart_id" api:"required" format:"uuid"`
	WidthUnits int64                                              `json:"width_units" api:"required"`
	JSON       customChartsSectionLayoutBreakpointsSmRowsItemJSON `json:"-"`
}

// customChartsSectionLayoutBreakpointsSmRowsItemJSON contains the JSON metadata
// for the struct [CustomChartsSectionLayoutBreakpointsSmRowsItem]
type customChartsSectionLayoutBreakpointsSmRowsItemJSON struct {
	ChartID     apijson.Field
	WidthUnits  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionLayoutBreakpointsSmRowsItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionLayoutBreakpointsSmRowsItemJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionLayoutVersion int64

const (
	CustomChartsSectionLayoutVersion1 CustomChartsSectionLayoutVersion = 1
)

func (r CustomChartsSectionLayoutVersion) IsKnown() bool {
	switch r {
	case CustomChartsSectionLayoutVersion1:
		return true
	}
	return false
}

type CustomChartsSectionSubSection struct {
	ID          string                                `json:"id" api:"required" format:"uuid"`
	Charts      []CustomChartsSectionSubSectionsChart `json:"charts" api:"required"`
	Index       int64                                 `json:"index" api:"required"`
	Title       string                                `json:"title" api:"required"`
	Description string                                `json:"description" api:"nullable"`
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
	ID            string                                            `json:"id" api:"required" format:"uuid"`
	ChartType     CustomChartsSectionSubSectionsChartsChartType     `json:"chart_type" api:"required"`
	Data          []CustomChartsSectionSubSectionsChartsData        `json:"data" api:"required"`
	Index         int64                                             `json:"index" api:"required"`
	Series        []CustomChartsSectionSubSectionsChartsSeries      `json:"series" api:"required"`
	Title         string                                            `json:"title" api:"required"`
	CommonFilters CustomChartsSectionSubSectionsChartsCommonFilters `json:"common_filters" api:"nullable"`
	Description   string                                            `json:"description" api:"nullable"`
	Metadata      map[string]interface{}                            `json:"metadata" api:"nullable"`
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

type CustomChartsSectionSubSectionsChartsChartType string

const (
	CustomChartsSectionSubSectionsChartsChartTypeLine  CustomChartsSectionSubSectionsChartsChartType = "line"
	CustomChartsSectionSubSectionsChartsChartTypeBar   CustomChartsSectionSubSectionsChartsChartType = "bar"
	CustomChartsSectionSubSectionsChartsChartTypeTable CustomChartsSectionSubSectionsChartsChartType = "table"
	CustomChartsSectionSubSectionsChartsChartTypeKpi   CustomChartsSectionSubSectionsChartsChartType = "kpi"
	CustomChartsSectionSubSectionsChartsChartTypeTopK  CustomChartsSectionSubSectionsChartsChartType = "top-k"
	CustomChartsSectionSubSectionsChartsChartTypePie   CustomChartsSectionSubSectionsChartsChartType = "pie"
)

func (r CustomChartsSectionSubSectionsChartsChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsChartTypeLine, CustomChartsSectionSubSectionsChartsChartTypeBar, CustomChartsSectionSubSectionsChartsChartTypeTable, CustomChartsSectionSubSectionsChartsChartTypeKpi, CustomChartsSectionSubSectionsChartsChartTypeTopK, CustomChartsSectionSubSectionsChartsChartTypePie:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsData struct {
	SeriesID  string                                             `json:"series_id" api:"required"`
	Timestamp time.Time                                          `json:"timestamp" api:"required" format:"date-time"`
	Value     CustomChartsSectionSubSectionsChartsDataValueUnion `json:"value" api:"required,nullable"`
	Group     string                                             `json:"group" api:"nullable"`
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
	ID               string                                                     `json:"id" api:"required" format:"uuid"`
	Name             string                                                     `json:"name" api:"required"`
	FeedbackKey      string                                                     `json:"feedback_key" api:"nullable"`
	FilterDefinition CustomChartsSectionSubSectionsChartsSeriesFilterDefinition `json:"filter_definition" api:"nullable"`
	Filters          CustomChartsSectionSubSectionsChartsSeriesFilters          `json:"filters" api:"nullable"`
	// Include additional information about where the group_by param was set.
	GroupBy            CustomChartsSectionSubSectionsChartsSeriesGroupBy             `json:"group_by" api:"nullable"`
	GroupByDefinitions []CustomChartsSectionSubSectionsChartsSeriesGroupByDefinition `json:"group_by_definitions" api:"nullable"`
	Metadata           map[string]interface{}                                        `json:"metadata" api:"nullable"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric           CustomChartsSectionSubSectionsChartsSeriesMetric           `json:"metric" api:"nullable"`
	MetricDefinition CustomChartsSectionSubSectionsChartsSeriesMetricDefinition `json:"metric_definition" api:"nullable"`
	// LGP Metrics you can chart.
	ProjectMetric CustomChartsSectionSubSectionsChartsSeriesProjectMetric `json:"project_metric" api:"nullable"`
	WorkspaceID   string                                                  `json:"workspace_id" api:"nullable" format:"uuid"`
	JSON          customChartsSectionSubSectionsChartsSeriesJSON          `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesJSON contains the JSON metadata for
// the struct [CustomChartsSectionSubSectionsChartsSeries]
type customChartsSectionSubSectionsChartsSeriesJSON struct {
	ID                 apijson.Field
	Name               apijson.Field
	FeedbackKey        apijson.Field
	FilterDefinition   apijson.Field
	Filters            apijson.Field
	GroupBy            apijson.Field
	GroupByDefinitions apijson.Field
	Metadata           apijson.Field
	Metric             apijson.Field
	MetricDefinition   apijson.Field
	ProjectMetric      apijson.Field
	WorkspaceID        apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeries) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesFilterDefinition struct {
	SourceType CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceType `json:"source_type" api:"required"`
	// This field can have the runtime type of [[]string].
	DatasetIDs interface{} `json:"dataset_ids"`
	// This field can have the runtime type of [[]string].
	ProjectIDs  interface{}                                                    `json:"project_ids"`
	RunFilter   string                                                         `json:"run_filter" api:"nullable"`
	TraceFilter string                                                         `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                         `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionSubSectionsChartsSeriesFilterDefinitionJSON `json:"-"`
	union       CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionUnion
}

// customChartsSectionSubSectionsChartsSeriesFilterDefinitionJSON contains the JSON
// metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinition]
type customChartsSectionSubSectionsChartsSeriesFilterDefinitionJSON struct {
	SourceType  apijson.Field
	DatasetIDs  apijson.Field
	ProjectIDs  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionSubSectionsChartsSeriesFilterDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionSubSectionsChartsSeriesFilterDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionSubSectionsChartsSeriesFilterDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject],
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset].
func (r CustomChartsSectionSubSectionsChartsSeriesFilterDefinition) AsUnion() CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject]
// or
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset].
type CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionUnion interface {
	implementsCustomChartsSectionSubSectionsChartsSeriesFilterDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset{}),
		},
	)
}

type CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject struct {
	ProjectIDs  []string                                                                                              `json:"project_ids" api:"required" format:"uuid"`
	SourceType  CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType `json:"source_type" api:"required"`
	RunFilter   string                                                                                                `json:"run_filter" api:"nullable"`
	TraceFilter string                                                                                                `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                                                `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON       `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject]
type customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON struct {
	ProjectIDs  apijson.Field
	SourceType  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProject) implementsCustomChartsSectionSubSectionsChartsSeriesFilterDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType string

const (
	CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType = "tracing_project"
)

func (r CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset struct {
	DatasetIDs []string                                                                                       `json:"dataset_ids" api:"required" format:"uuid"`
	SourceType CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType `json:"source_type" api:"required"`
	JSON       customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON       `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset]
type customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON struct {
	DatasetIDs  apijson.Field
	SourceType  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDataset) implementsCustomChartsSectionSubSectionsChartsSeriesFilterDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType string

const (
	CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType = "dataset"
)

func (r CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceType string

const (
	CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceTypeTracingProject CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceType = "tracing_project"
	CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceTypeDataset        CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceType = "dataset"
)

func (r CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceTypeTracingProject, CustomChartsSectionSubSectionsChartsSeriesFilterDefinitionSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesFilters struct {
	Filter      string                                                `json:"filter" api:"nullable"`
	Session     []string                                              `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                                `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                `json:"tree_filter" api:"nullable"`
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
	Attribute CustomChartsSectionSubSectionsChartsSeriesGroupByAttribute `json:"attribute" api:"required"`
	MaxGroups int64                                                      `json:"max_groups"`
	Path      string                                                     `json:"path" api:"nullable"`
	SetBy     CustomChartsSectionSubSectionsChartsSeriesGroupBySetBy     `json:"set_by" api:"nullable"`
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

type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinition struct {
	Attribute CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute `json:"attribute" api:"required"`
	Path      string                                                                `json:"path"`
	JSON      customChartsSectionSubSectionsChartsSeriesGroupByDefinitionJSON       `json:"-"`
	union     CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsUnion
}

// customChartsSectionSubSectionsChartsSeriesGroupByDefinitionJSON contains the
// JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinition]
type customChartsSectionSubSectionsChartsSeriesGroupByDefinitionJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionSubSectionsChartsSeriesGroupByDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionSubSectionsChartsSeriesGroupByDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionSubSectionsChartsSeriesGroupByDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain],
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex].
func (r CustomChartsSectionSubSectionsChartsSeriesGroupByDefinition) AsUnion() CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain]
// or
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex].
type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsUnion interface {
	implementsCustomChartsSectionSubSectionsChartsSeriesGroupByDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex{}),
		},
	)
}

type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain struct {
	Attribute CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute `json:"attribute" api:"required"`
	JSON      customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON      `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain]
type customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON struct {
	Attribute   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlain) implementsCustomChartsSectionSubSectionsChartsSeriesGroupByDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute string

const (
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName    CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "name"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "run_type"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag     CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "tag"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "project"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus  CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "status"
)

func (r CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex struct {
	Attribute CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute `json:"attribute" api:"required"`
	Path      string                                                                                         `json:"path" api:"required"`
	JSON      customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON      `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex]
type customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplex) implementsCustomChartsSectionSubSectionsChartsSeriesGroupByDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute string

const (
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata      CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "metadata"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "feedback_label"
)

func (r CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute string

const (
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeName          CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "name"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeRunType       CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "run_type"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeTag           CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "tag"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeProject       CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "project"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeStatus        CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "status"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeMetadata      CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "metadata"
	CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeFeedbackLabel CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute = "feedback_label"
)

func (r CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeName, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeRunType, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeTag, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeProject, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeStatus, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeMetadata, CustomChartsSectionSubSectionsChartsSeriesGroupByDefinitionsAttributeFeedbackLabel:
		return true
	}
	return false
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

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinition struct {
	// This field can have the runtime type of
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator].
	Denominator interface{}                                                      `json:"denominator"`
	Entity      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionEntity `json:"entity"`
	Field       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField  `json:"field"`
	Filter      string                                                           `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator].
	Numerator interface{} `json:"numerator"`
	// This field can have the runtime type of
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParams].
	Params interface{}                                                    `json:"params"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionJSON `json:"-"`
	union  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionUnion
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionJSON contains the JSON
// metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinition]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionJSON struct {
	Denominator apijson.Field
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Numerator   apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionSubSectionsChartsSeriesMetricDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput].
func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinition) AsUnion() CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile]
// or
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput].
type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionUnion interface {
	implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput{}),
		},
	)
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                         `json:"filter" api:"nullable"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetric) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                             `json:"feedback_key" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount struct {
	Filter string                                                                               `json:"filter" api:"nullable"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountTypeCount CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                               `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                   `json:"feedback_key" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                 `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeSum CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMax CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMin CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeAvg CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                   `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                       `json:"feedback_key" api:"required"`
	P           float64                                                                                                      `json:"p" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                      `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParams struct {
	P    float64                                                                                         `json:"p" api:"required"`
	JSON customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput struct {
	Denominator CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator `json:"denominator" api:"required"`
	Numerator   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator   `json:"numerator" api:"required"`
	Type        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType        `json:"type" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON        `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON struct {
	Denominator apijson.Field
	Numerator   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinition() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator struct {
	Entity CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity `json:"entity"`
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField  `json:"field"`
	Filter string                                                                                                  `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams].
	Params interface{}                                                                                           `json:"params"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON `json:"-"`
	union  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON struct {
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) AsUnion() CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile]
// or
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion interface {
	implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                                `json:"filter" api:"nullable"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                                    `json:"feedback_key" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricTypeCount CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount struct {
	Filter string                                                                                                                      `json:"filter" api:"nullable"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                                      `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                                          `json:"feedback_key" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                                        `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                          `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                                              `json:"feedback_key" api:"required"`
	P           float64                                                                                                                                             `json:"p" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                             `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                                `json:"p" api:"required"`
	JSON customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntityFeedback CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity = "feedback"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "feedback_score"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_cost"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFeedbackScore, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "count"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator struct {
	Entity CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity `json:"entity"`
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField  `json:"field"`
	Filter string                                                                                                `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams].
	Params interface{}                                                                                         `json:"params"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON `json:"-"`
	union  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON struct {
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) AsUnion() CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile]
// or
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion interface {
	implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                              `json:"filter" api:"nullable"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                                  `json:"feedback_key" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricTypeCount CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount struct {
	Filter string                                                                                                                    `json:"filter" api:"nullable"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                                    `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                                        `json:"feedback_key" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                                      `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                        `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                                            `json:"feedback_key" api:"required"`
	P           float64                                                                                                                                           `json:"p" api:"required"`
	JSON        customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                           `json:"filter" api:"nullable"`
	JSON   customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) implementsCustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                              `json:"p" api:"required"`
	JSON customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams]
type customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntityFeedback CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity = "feedback"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "feedback_score"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_cost"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFeedbackScore, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "count"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType = "ratio"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionEntity string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionEntityFeedback CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionEntity = "feedback"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldFeedbackScore     CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "feedback_score"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldLatencySeconds    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "latency_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldFirstTokenSeconds CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "first_token_seconds"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldTotalTokens       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "total_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldPromptTokens      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "prompt_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldCompletionTokens  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "completion_tokens"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldTotalCost         CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "total_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldPromptCost        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "prompt_cost"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldCompletionCost    CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField = "completion_cost"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldFeedbackScore, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeCount      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "count"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMax        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMin        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeAvg        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeSum        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "percentile"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeRatio      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "ratio"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeCount, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeAvg, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypePercentile, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeRatio:
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
	CustomChartsSectionSubSectionsChartsSeriesProjectMetricRunQueueWaitTime        CustomChartsSectionSubSectionsChartsSeriesProjectMetric = "run_queue_wait_time"
)

func (r CustomChartsSectionSubSectionsChartsSeriesProjectMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesProjectMetricMemoryUsage, CustomChartsSectionSubSectionsChartsSeriesProjectMetricCPUUsage, CustomChartsSectionSubSectionsChartsSeriesProjectMetricDiskUsage, CustomChartsSectionSubSectionsChartsSeriesProjectMetricRestartCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricReplicaCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricWorkerCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricLgRunCount, CustomChartsSectionSubSectionsChartsSeriesProjectMetricResponsesPerSecond, CustomChartsSectionSubSectionsChartsSeriesProjectMetricErrorResponsesPerSecond, CustomChartsSectionSubSectionsChartsSeriesProjectMetricP95Latency, CustomChartsSectionSubSectionsChartsSeriesProjectMetricRunQueueWaitTime:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsCommonFilters struct {
	Filter      string                                                `json:"filter" api:"nullable"`
	Session     []string                                              `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                                `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                `json:"tree_filter" api:"nullable"`
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
	Attribute param.Field[RunStatsGroupByAttribute] `json:"attribute" api:"required"`
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
)

func (r SessionSortableColumns) IsKnown() bool {
	switch r {
	case SessionSortableColumnsName, SessionSortableColumnsStartTime, SessionSortableColumnsLastRunStartTime, SessionSortableColumnsLatencyP50, SessionSortableColumnsLatencyP99, SessionSortableColumnsErrorRate, SessionSortableColumnsFeedback:
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
	ID                   string                          `json:"id" api:"required" format:"uuid"`
	TenantID             string                          `json:"tenant_id" api:"required" format:"uuid"`
	CompletionCost       string                          `json:"completion_cost" api:"nullable"`
	CompletionTokens     int64                           `json:"completion_tokens" api:"nullable"`
	DefaultDatasetID     string                          `json:"default_dataset_id" api:"nullable" format:"uuid"`
	Description          string                          `json:"description" api:"nullable"`
	EndTime              time.Time                       `json:"end_time" api:"nullable" format:"date-time"`
	ErrorRate            float64                         `json:"error_rate" api:"nullable"`
	ExperimentProgress   TracerSessionExperimentProgress `json:"experiment_progress" api:"nullable"`
	Extra                map[string]interface{}          `json:"extra" api:"nullable"`
	FeedbackStats        map[string]interface{}          `json:"feedback_stats" api:"nullable"`
	FirstTokenP50        float64                         `json:"first_token_p50" api:"nullable"`
	FirstTokenP99        float64                         `json:"first_token_p99" api:"nullable"`
	LastRunStartTime     time.Time                       `json:"last_run_start_time" api:"nullable" format:"date-time"`
	LastRunStartTimeLive time.Time                       `json:"last_run_start_time_live" api:"nullable" format:"date-time"`
	LatencyP50           float64                         `json:"latency_p50" api:"nullable"`
	LatencyP99           float64                         `json:"latency_p99" api:"nullable"`
	Name                 string                          `json:"name"`
	PromptCost           string                          `json:"prompt_cost" api:"nullable"`
	PromptTokens         int64                           `json:"prompt_tokens" api:"nullable"`
	ReferenceDatasetID   string                          `json:"reference_dataset_id" api:"nullable" format:"uuid"`
	RunCount             int64                           `json:"run_count" api:"nullable"`
	RunFacets            []map[string]interface{}        `json:"run_facets" api:"nullable"`
	SessionFeedbackStats map[string]interface{}          `json:"session_feedback_stats" api:"nullable"`
	StartTime            time.Time                       `json:"start_time" format:"date-time"`
	StreamingRate        float64                         `json:"streaming_rate" api:"nullable"`
	TestRunNumber        int64                           `json:"test_run_number" api:"nullable"`
	TotalCost            string                          `json:"total_cost" api:"nullable"`
	TotalTokens          int64                           `json:"total_tokens" api:"nullable"`
	TraceTier            TracerSessionTraceTier          `json:"trace_tier" api:"nullable"`
	JSON                 tracerSessionJSON               `json:"-"`
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
	ExperimentProgress   apijson.Field
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

type TracerSessionExperimentProgress struct {
	EvaluatorProgress map[string]float64                  `json:"evaluator_progress" api:"required"`
	ExpectedRunCount  int64                               `json:"expected_run_count" api:"required"`
	RunProgress       float64                             `json:"run_progress" api:"required"`
	JSON              tracerSessionExperimentProgressJSON `json:"-"`
}

// tracerSessionExperimentProgressJSON contains the JSON metadata for the struct
// [TracerSessionExperimentProgress]
type tracerSessionExperimentProgressJSON struct {
	EvaluatorProgress apijson.Field
	ExpectedRunCount  apijson.Field
	RunProgress       apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *TracerSessionExperimentProgress) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r tracerSessionExperimentProgressJSON) RawJSON() string {
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
	ID                   string                                     `json:"id" api:"required" format:"uuid"`
	TenantID             string                                     `json:"tenant_id" api:"required" format:"uuid"`
	DefaultDatasetID     string                                     `json:"default_dataset_id" api:"nullable" format:"uuid"`
	Description          string                                     `json:"description" api:"nullable"`
	EndTime              time.Time                                  `json:"end_time" api:"nullable" format:"date-time"`
	Extra                map[string]interface{}                     `json:"extra" api:"nullable"`
	LastRunStartTimeLive time.Time                                  `json:"last_run_start_time_live" api:"nullable" format:"date-time"`
	Name                 string                                     `json:"name"`
	ReferenceDatasetID   string                                     `json:"reference_dataset_id" api:"nullable" format:"uuid"`
	StartTime            time.Time                                  `json:"start_time" format:"date-time"`
	TraceTier            TracerSessionWithoutVirtualFieldsTraceTier `json:"trace_tier" api:"nullable"`
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
	EvaluatorKeys      param.Field[[]string]                  `json:"evaluator_keys"`
	Extra              param.Field[map[string]interface{}]    `json:"extra"`
	KickedOffBy        param.Field[string]                    `json:"kicked_off_by"`
	Name               param.Field[string]                    `json:"name"`
	NumExamples        param.Field[int64]                     `json:"num_examples"`
	NumRepetitions     param.Field[int64]                     `json:"num_repetitions"`
	ReferenceDatasetID param.Field[string]                    `json:"reference_dataset_id" format:"uuid"`
	StartTime          param.Field[time.Time]                 `json:"start_time" format:"date-time"`
	TagValueIDs        param.Field[[]string]                  `json:"tag_value_ids" format:"uuid"`
	TraceTier          param.Field[SessionNewParamsTraceTier] `json:"trace_tier"`
}

func (r SessionNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [SessionNewParams]'s query parameters as `url.Values`.
func (r SessionNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
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
	ID                   param.Field[[]string]                              `query:"id" format:"uuid"`
	DatasetVersion       param.Field[string]                                `query:"dataset_version"`
	Facets               param.Field[bool]                                  `query:"facets"`
	Filter               param.Field[string]                                `query:"filter"`
	IncludeStats         param.Field[bool]                                  `query:"include_stats"`
	Limit                param.Field[int64]                                 `query:"limit"`
	Metadata             param.Field[string]                                `query:"metadata"`
	Name                 param.Field[string]                                `query:"name"`
	NameContains         param.Field[string]                                `query:"name_contains"`
	Offset               param.Field[int64]                                 `query:"offset"`
	ReferenceDataset     param.Field[[]string]                              `query:"reference_dataset" format:"uuid"`
	ReferenceFree        param.Field[bool]                                  `query:"reference_free"`
	SortBy               param.Field[SessionSortableColumns]                `query:"sort_by"`
	SortByDesc           param.Field[bool]                                  `query:"sort_by_desc"`
	SortByFeedbackKey    param.Field[string]                                `query:"sort_by_feedback_key"`
	SortByFeedbackSource param.Field[SessionListParamsSortByFeedbackSource] `query:"sort_by_feedback_source"`
	StatsFilter          param.Field[string]                                `query:"stats_filter"`
	StatsSelect          param.Field[[]string]                              `query:"stats_select"`
	StatsStartTime       param.Field[time.Time]                             `query:"stats_start_time" format:"date-time"`
	TagValueID           param.Field[[]string]                              `query:"tag_value_id" format:"uuid"`
	UseApproxStats       param.Field[bool]                                  `query:"use_approx_stats"`
	Accept               param.Field[string]                                `header:"accept"`
}

// URLQuery serializes [SessionListParams]'s query parameters as `url.Values`.
func (r SessionListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionListParamsSortByFeedbackSource string

const (
	SessionListParamsSortByFeedbackSourceSession SessionListParamsSortByFeedbackSource = "session"
	SessionListParamsSortByFeedbackSourceRun     SessionListParamsSortByFeedbackSource = "run"
)

func (r SessionListParamsSortByFeedbackSource) IsKnown() bool {
	switch r {
	case SessionListParamsSortByFeedbackSourceSession, SessionListParamsSortByFeedbackSourceRun:
		return true
	}
	return false
}

type SessionDashboardParams struct {
	CustomChartsSectionRequest CustomChartsSectionRequestParam `json:"custom_charts_section_request" api:"required"`
	Accept                     param.Field[string]             `header:"accept"`
}

func (r SessionDashboardParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.CustomChartsSectionRequest)
}
