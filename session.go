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
	// [CustomChartsSectionChartsSingleCustomChartResponseCommonFilters].
	CommonFilters interface{} `json:"common_filters"`
	// This field can have the runtime type of
	// [[]CustomChartsSectionChartsSingleCustomChartResponseData].
	Data        interface{} `json:"data"`
	Description string      `json:"description" api:"nullable"`
	Markdown    string      `json:"markdown"`
	// This field can have the runtime type of [map[string]interface{}].
	Metadata interface{} `json:"metadata"`
	// This field can have the runtime type of
	// [[]CustomChartsSectionChartsSingleCustomChartResponseSeries].
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
// [CustomChartsSectionChartsSingleCustomChartResponse],
// [CustomChartsSectionChartsCustomTextBlock].
func (r CustomChartsSectionChart) AsUnion() CustomChartsSectionChartsUnion {
	return r.union
}

// Union satisfied by [CustomChartsSectionChartsSingleCustomChartResponse] or
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
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponse{}),
			DiscriminatorValue: "line",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponse{}),
			DiscriminatorValue: "bar",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponse{}),
			DiscriminatorValue: "table",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponse{}),
			DiscriminatorValue: "kpi",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponse{}),
			DiscriminatorValue: "top-k",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponse{}),
			DiscriminatorValue: "pie",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(CustomChartsSectionChartsCustomTextBlock{}),
			DiscriminatorValue: "text",
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponse struct {
	ID            string                                                          `json:"id" api:"required" format:"uuid"`
	ChartType     CustomChartsSectionChartsSingleCustomChartResponseChartType     `json:"chart_type" api:"required"`
	Data          []CustomChartsSectionChartsSingleCustomChartResponseData        `json:"data" api:"required"`
	Index         int64                                                           `json:"index" api:"required"`
	Series        []CustomChartsSectionChartsSingleCustomChartResponseSeries      `json:"series" api:"required"`
	Title         string                                                          `json:"title" api:"required"`
	CommonFilters CustomChartsSectionChartsSingleCustomChartResponseCommonFilters `json:"common_filters" api:"nullable"`
	Description   string                                                          `json:"description" api:"nullable"`
	Metadata      map[string]interface{}                                          `json:"metadata" api:"nullable"`
	JSON          customChartsSectionChartsSingleCustomChartResponseJSON          `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseJSON contains the JSON
// metadata for the struct [CustomChartsSectionChartsSingleCustomChartResponse]
type customChartsSectionChartsSingleCustomChartResponseJSON struct {
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

func (r *CustomChartsSectionChartsSingleCustomChartResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponse) implementsCustomChartsSectionChart() {}

type CustomChartsSectionChartsSingleCustomChartResponseChartType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseChartTypeLine  CustomChartsSectionChartsSingleCustomChartResponseChartType = "line"
	CustomChartsSectionChartsSingleCustomChartResponseChartTypeBar   CustomChartsSectionChartsSingleCustomChartResponseChartType = "bar"
	CustomChartsSectionChartsSingleCustomChartResponseChartTypeTable CustomChartsSectionChartsSingleCustomChartResponseChartType = "table"
	CustomChartsSectionChartsSingleCustomChartResponseChartTypeKpi   CustomChartsSectionChartsSingleCustomChartResponseChartType = "kpi"
	CustomChartsSectionChartsSingleCustomChartResponseChartTypeTopK  CustomChartsSectionChartsSingleCustomChartResponseChartType = "top-k"
	CustomChartsSectionChartsSingleCustomChartResponseChartTypePie   CustomChartsSectionChartsSingleCustomChartResponseChartType = "pie"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseChartTypeLine, CustomChartsSectionChartsSingleCustomChartResponseChartTypeBar, CustomChartsSectionChartsSingleCustomChartResponseChartTypeTable, CustomChartsSectionChartsSingleCustomChartResponseChartTypeKpi, CustomChartsSectionChartsSingleCustomChartResponseChartTypeTopK, CustomChartsSectionChartsSingleCustomChartResponseChartTypePie:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseData struct {
	SeriesID  string                                                           `json:"series_id" api:"required"`
	Timestamp time.Time                                                        `json:"timestamp" api:"required" format:"date-time"`
	Value     CustomChartsSectionChartsSingleCustomChartResponseDataValueUnion `json:"value" api:"required,nullable"`
	Group     string                                                           `json:"group" api:"nullable"`
	JSON      customChartsSectionChartsSingleCustomChartResponseDataJSON       `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseDataJSON contains the JSON
// metadata for the struct [CustomChartsSectionChartsSingleCustomChartResponseData]
type customChartsSectionChartsSingleCustomChartResponseDataJSON struct {
	SeriesID    apijson.Field
	Timestamp   apijson.Field
	Value       apijson.Field
	Group       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseDataJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionFloat] or
// [CustomChartsSectionChartsSingleCustomChartResponseDataValueMap].
type CustomChartsSectionChartsSingleCustomChartResponseDataValueUnion interface {
	ImplementsCustomChartsSectionChartsSingleCustomChartResponseDataValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseDataValueUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseDataValueMap{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseDataValueMap map[string]interface{}

func (r CustomChartsSectionChartsSingleCustomChartResponseDataValueMap) ImplementsCustomChartsSectionChartsSingleCustomChartResponseDataValueUnion() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeries struct {
	ID               string                                                                   `json:"id" api:"required" format:"uuid"`
	Name             string                                                                   `json:"name" api:"required"`
	FeedbackKey      string                                                                   `json:"feedback_key" api:"nullable"`
	FilterDefinition CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition `json:"filter_definition" api:"nullable"`
	Filters          CustomChartsSectionChartsSingleCustomChartResponseSeriesFilters          `json:"filters" api:"nullable"`
	// Include additional information about where the group_by param was set.
	GroupBy            CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBy             `json:"group_by" api:"nullable"`
	GroupByDefinitions []CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition `json:"group_by_definitions" api:"nullable"`
	Metadata           map[string]interface{}                                                      `json:"metadata" api:"nullable"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric           CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric           `json:"metric" api:"nullable"`
	MetricDefinition CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition `json:"metric_definition" api:"nullable"`
	// LGP Metrics you can chart.
	ProjectMetric CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric `json:"project_metric" api:"nullable"`
	WorkspaceID   string                                                                `json:"workspace_id" api:"nullable" format:"uuid"`
	JSON          customChartsSectionChartsSingleCustomChartResponseSeriesJSON          `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesJSON contains the JSON
// metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeries]
type customChartsSectionChartsSingleCustomChartResponseSeriesJSON struct {
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

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeries) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition struct {
	SourceType CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceType `json:"source_type" api:"required"`
	// This field can have the runtime type of [[]string].
	DatasetIDs interface{} `json:"dataset_ids"`
	// This field can have the runtime type of [[]string].
	ProjectIDs  interface{}                                                                  `json:"project_ids"`
	RunFilter   string                                                                       `json:"run_filter" api:"nullable"`
	TraceFilter string                                                                       `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                       `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionJSON `json:"-"`
	union       CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionUnion
}

// customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition]
type customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionJSON struct {
	SourceType  apijson.Field
	DatasetIDs  apijson.Field
	ProjectIDs  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset].
func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset].
type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject struct {
	ProjectIDs  []string                                                                                                            `json:"project_ids" api:"required" format:"uuid"`
	SourceType  CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType `json:"source_type" api:"required"`
	RunFilter   string                                                                                                              `json:"run_filter" api:"nullable"`
	TraceFilter string                                                                                                              `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                                                              `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON       `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject]
type customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON struct {
	ProjectIDs  apijson.Field
	SourceType  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProject) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType = "tracing_project"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset struct {
	DatasetIDs []string                                                                                                     `json:"dataset_ids" api:"required" format:"uuid"`
	SourceType CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetSourceType `json:"source_type" api:"required"`
	JSON       customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetJSON       `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset]
type customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetJSON struct {
	DatasetIDs  apijson.Field
	SourceType  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDataset) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetSourceType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetSourceType = "dataset"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceTypeTracingProject CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceType = "tracing_project"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceTypeDataset        CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceType = "dataset"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceTypeTracingProject, CustomChartsSectionChartsSingleCustomChartResponseSeriesFilterDefinitionSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesFilters struct {
	Filter      string                                                              `json:"filter" api:"nullable"`
	Session     []string                                                            `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                                              `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                              `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesFiltersJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesFiltersJSON contains the
// JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesFilters]
type customChartsSectionChartsSingleCustomChartResponseSeriesFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesFiltersJSON) RawJSON() string {
	return r.raw
}

// Include additional information about where the group_by param was set.
type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBy struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute `json:"attribute" api:"required"`
	MaxGroups int64                                                                    `json:"max_groups"`
	Path      string                                                                   `json:"path" api:"nullable"`
	SetBy     CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBy     `json:"set_by" api:"nullable"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSeriesGroupByJSON      `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesGroupByJSON contains the
// JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBy]
type customChartsSectionChartsSingleCustomChartResponseSeriesGroupByJSON struct {
	Attribute   apijson.Field
	MaxGroups   apijson.Field
	Path        apijson.Field
	SetBy       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBy) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesGroupByJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeName     CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute = "name"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeRunType  CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute = "run_type"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeTag      CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute = "tag"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeMetadata CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute = "metadata"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeName, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeRunType, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeTag, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByAttributeMetadata:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBy string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBySection CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBy = "section"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBySeries  CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBy = "series"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBy) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBySection, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupBySetBySeries:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute `json:"attribute" api:"required"`
	Path      string                                                                              `json:"path"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionJSON       `json:"-"`
	union     CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsUnion
}

// customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition]
type customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex].
func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex].
type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute `json:"attribute" api:"required"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainJSON      `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain]
type customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainJSON struct {
	Attribute   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlain) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName    CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "name"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "run_type"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag     CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "tag"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "project"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus  CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "status"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex struct {
	Attribute CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute `json:"attribute" api:"required"`
	Path      string                                                                                                       `json:"path" api:"required"`
	JSON      customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexJSON      `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex]
type customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplex) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata      CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "metadata"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "feedback_label"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeName          CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "name"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeRunType       CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "run_type"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeTag           CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "tag"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeProject       CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "project"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeStatus        CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "status"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeMetadata      CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "metadata"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeFeedbackLabel CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute = "feedback_label"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeName, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeRunType, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeTag, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeProject, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeStatus, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeMetadata, CustomChartsSectionChartsSingleCustomChartResponseSeriesGroupByDefinitionsAttributeFeedbackLabel:
		return true
	}
	return false
}

// Metrics you can chart. Feedback metrics are not available for
// organization-scoped charts.
type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricRunCount            CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "run_count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricLatencyP50          CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "latency_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricLatencyP99          CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "latency_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricLatencyAvg          CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "latency_avg"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFirstTokenP50       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "first_token_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFirstTokenP99       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "first_token_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricTotalTokens         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptTokens        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionTokens    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricMedianTokens        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "median_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionTokensP50 CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "completion_tokens_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptTokensP50     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "prompt_tokens_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricTokensP99           CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "tokens_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionTokensP99 CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "completion_tokens_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptTokensP99     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "prompt_tokens_p99"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFeedback            CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "feedback"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFeedbackScoreAvg    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "feedback_score_avg"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFeedbackValues      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "feedback_values"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricTotalCost           CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptCost          CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionCost      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricErrorRate           CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "error_rate"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricStreamingRate       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "streaming_rate"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCostP50             CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "cost_p50"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCostP99             CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric = "cost_p99"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricRunCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricLatencyP50, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricLatencyP99, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricLatencyAvg, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFirstTokenP50, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFirstTokenP99, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricMedianTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionTokensP50, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptTokensP50, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricTokensP99, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionTokensP99, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptTokensP99, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFeedback, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFeedbackScoreAvg, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricFeedbackValues, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricErrorRate, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricStreamingRate, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCostP50, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricCostP99:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition struct {
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator].
	Denominator interface{}                                                                    `json:"denominator"`
	Entity      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionEntity `json:"entity"`
	Field       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField  `json:"field"`
	Filter      string                                                                         `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator].
	Numerator interface{} `json:"numerator"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParams].
	Params interface{}                                                                  `json:"params"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionJSON `json:"-"`
	union  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionUnion
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionJSON struct {
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

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput].
func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput].
type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                       `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetric) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                           `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount struct {
	Filter string                                                                                             `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCount) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountTypeCount CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                             `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                 `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                               `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeSum CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                 `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                     `json:"feedback_key" api:"required"`
	P           float64                                                                                                                    `json:"p" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                    `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParams struct {
	P    float64                                                                                                       `json:"p" api:"required"`
	JSON customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput struct {
	Denominator CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator `json:"denominator" api:"required"`
	Numerator   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator   `json:"numerator" api:"required"`
	Type        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputType        `json:"type" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputJSON        `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputJSON struct {
	Denominator apijson.Field
	Numerator   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutput) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity `json:"entity"`
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField  `json:"field"`
	Filter string                                                                                                                `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams].
	Params interface{}                                                                                                         `json:"params"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON `json:"-"`
	union  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON struct {
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                                              `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetric) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                                                  `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricTypeCount CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount struct {
	Filter string                                                                                                                                    `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                                                    `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                                                        `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                                                      `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                                        `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                                                            `json:"feedback_key" api:"required"`
	P           float64                                                                                                                                                           `json:"p" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                           `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                                              `json:"p" api:"required"`
	JSON customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "feedback_score"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_cost"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFeedbackScore, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "avg"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity `json:"entity"`
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField  `json:"field"`
	Filter string                                                                                                              `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams],
	// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams].
	Params interface{}                                                                                                       `json:"params"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON `json:"-"`
	union  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON struct {
	Entity      apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) AsUnion() CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile]
// or
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion interface {
	implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric struct {
	Entity CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity `json:"entity" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams `json:"params" api:"required"`
	Filter string                                                                                                                                            `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType   `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON struct {
	Entity      apijson.Field
	Params      apijson.Field
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetric) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey string                                                                                                                                                `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricTypeCount CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount struct {
	Filter string                                                                                                                                  `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType   `json:"type" api:"required"`
	Filter string                                                                                                                                                  `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey string                                                                                                                                                      `json:"feedback_key" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON struct {
	FeedbackKey apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                                                    `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                                      `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey string                                                                                                                                                          `json:"feedback_key" api:"required"`
	P           float64                                                                                                                                                         `json:"p" api:"required"`
	JSON        customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON struct {
	FeedbackKey apijson.Field
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                                         `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) implementsCustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "feedback_score"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                                            `json:"p" api:"required"`
	JSON customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams]
type customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "feedback_score"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_cost"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFeedbackScore, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "avg"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "percentile"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputType = "ratio"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionEntity string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionEntityFeedback CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionEntity = "feedback"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionEntity) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionEntityFeedback:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldFeedbackScore     CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "feedback_score"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldLatencySeconds    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "latency_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldFirstTokenSeconds CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "first_token_seconds"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldTotalTokens       CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "total_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldPromptTokens      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "prompt_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldCompletionTokens  CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "completion_tokens"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldTotalCost         CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "total_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldPromptCost        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "prompt_cost"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldCompletionCost    CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField = "completion_cost"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldFeedbackScore, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldLatencySeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldFirstTokenSeconds, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldTotalTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldPromptTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldCompletionTokens, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldTotalCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldPromptCost, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeCount      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeMax        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "max"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeMin        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "min"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeAvg        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "avg"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeSum        CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "sum"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypePercentile CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "percentile"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeRatio      CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType = "ratio"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeMax, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeMin, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeAvg, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeSum, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypePercentile, CustomChartsSectionChartsSingleCustomChartResponseSeriesMetricDefinitionTypeRatio:
		return true
	}
	return false
}

// LGP Metrics you can chart.
type CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric string

const (
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricMemoryUsage             CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "memory_usage"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricCPUUsage                CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "cpu_usage"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricDiskUsage               CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "disk_usage"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricRestartCount            CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "restart_count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricReplicaCount            CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "replica_count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricWorkerCount             CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "worker_count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricLgRunCount              CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "lg_run_count"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricResponsesPerSecond      CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "responses_per_second"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricErrorResponsesPerSecond CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "error_responses_per_second"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricP95Latency              CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "p95_latency"
	CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricRunQueueWaitTime        CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric = "run_queue_wait_time"
)

func (r CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetric) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricMemoryUsage, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricCPUUsage, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricDiskUsage, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricRestartCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricReplicaCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricWorkerCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricLgRunCount, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricResponsesPerSecond, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricErrorResponsesPerSecond, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricP95Latency, CustomChartsSectionChartsSingleCustomChartResponseSeriesProjectMetricRunQueueWaitTime:
		return true
	}
	return false
}

type CustomChartsSectionChartsSingleCustomChartResponseCommonFilters struct {
	Filter      string                                                              `json:"filter" api:"nullable"`
	Session     []string                                                            `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                                              `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                              `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSingleCustomChartResponseCommonFiltersJSON `json:"-"`
}

// customChartsSectionChartsSingleCustomChartResponseCommonFiltersJSON contains the
// JSON metadata for the struct
// [CustomChartsSectionChartsSingleCustomChartResponseCommonFilters]
type customChartsSectionChartsSingleCustomChartResponseCommonFiltersJSON struct {
	Filter      apijson.Field
	Session     apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSingleCustomChartResponseCommonFilters) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSingleCustomChartResponseCommonFiltersJSON) RawJSON() string {
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
