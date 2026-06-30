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
	ID string `json:"id" api:"required" format:"uuid"`
	// Enum for custom chart types.
	ChartType     CustomChartsSectionChartsChartType     `json:"chart_type" api:"required"`
	Data          []CustomChartsSectionChartsData        `json:"data" api:"required"`
	Index         int64                                  `json:"index" api:"required"`
	Series        []CustomChartsSectionChartsSeries      `json:"series" api:"required"`
	Title         string                                 `json:"title" api:"required"`
	CommonFilters CustomChartsSectionChartsCommonFilters `json:"common_filters" api:"nullable"`
	Description   string                                 `json:"description" api:"nullable"`
	Metadata      map[string]interface{}                 `json:"metadata" api:"nullable"`
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
	CustomChartsSectionChartsChartTypeLine  CustomChartsSectionChartsChartType = "line"
	CustomChartsSectionChartsChartTypeBar   CustomChartsSectionChartsChartType = "bar"
	CustomChartsSectionChartsChartTypeTable CustomChartsSectionChartsChartType = "table"
	CustomChartsSectionChartsChartTypeKpi   CustomChartsSectionChartsChartType = "kpi"
	CustomChartsSectionChartsChartTypeTopK  CustomChartsSectionChartsChartType = "top-k"
	CustomChartsSectionChartsChartTypePie   CustomChartsSectionChartsChartType = "pie"
)

func (r CustomChartsSectionChartsChartType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsChartTypeLine, CustomChartsSectionChartsChartTypeBar, CustomChartsSectionChartsChartTypeTable, CustomChartsSectionChartsChartTypeKpi, CustomChartsSectionChartsChartTypeTopK, CustomChartsSectionChartsChartTypePie:
		return true
	}
	return false
}

type CustomChartsSectionChartsData struct {
	SeriesID  string                                  `json:"series_id" api:"required"`
	Timestamp time.Time                               `json:"timestamp" api:"required" format:"date-time"`
	Value     CustomChartsSectionChartsDataValueUnion `json:"value" api:"required,nullable"`
	Group     string                                  `json:"group" api:"nullable"`
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
	ID               string                                          `json:"id" api:"required" format:"uuid"`
	Name             string                                          `json:"name" api:"required"`
	FeedbackKey      string                                          `json:"feedback_key" api:"nullable"`
	FilterDefinition CustomChartsSectionChartsSeriesFilterDefinition `json:"filter_definition" api:"nullable"`
	Filters          CustomChartsSectionChartsSeriesFilters          `json:"filters" api:"nullable"`
	// Include additional information about where the group_by param was set.
	GroupBy            CustomChartsSectionChartsSeriesGroupBy             `json:"group_by" api:"nullable"`
	GroupByDefinitions []CustomChartsSectionChartsSeriesGroupByDefinition `json:"group_by_definitions" api:"nullable"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric           CustomChartsSectionChartsSeriesMetric           `json:"metric" api:"nullable"`
	MetricDefinition CustomChartsSectionChartsSeriesMetricDefinition `json:"metric_definition" api:"nullable"`
	// LGP Metrics you can chart.
	ProjectMetric CustomChartsSectionChartsSeriesProjectMetric `json:"project_metric" api:"nullable"`
	WorkspaceID   string                                       `json:"workspace_id" api:"nullable" format:"uuid"`
	JSON          customChartsSectionChartsSeriesJSON          `json:"-"`
}

// customChartsSectionChartsSeriesJSON contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeries]
type customChartsSectionChartsSeriesJSON struct {
	ID                 apijson.Field
	Name               apijson.Field
	FeedbackKey        apijson.Field
	FilterDefinition   apijson.Field
	Filters            apijson.Field
	GroupBy            apijson.Field
	GroupByDefinitions apijson.Field
	Metric             apijson.Field
	MetricDefinition   apijson.Field
	ProjectMetric      apijson.Field
	WorkspaceID        apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeries) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSeriesFilterDefinition struct {
	SourceType CustomChartsSectionChartsSeriesFilterDefinitionSourceType `json:"source_type" api:"required"`
	// This field can have the runtime type of [[]string].
	DatasetIDs interface{} `json:"dataset_ids"`
	// This field can have the runtime type of [[]string].
	ProjectIDs  interface{}                                         `json:"project_ids"`
	RunFilter   string                                              `json:"run_filter" api:"nullable"`
	TraceFilter string                                              `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                              `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSeriesFilterDefinitionJSON `json:"-"`
	union       CustomChartsSectionChartsSeriesFilterDefinitionUnion
}

// customChartsSectionChartsSeriesFilterDefinitionJSON contains the JSON metadata
// for the struct [CustomChartsSectionChartsSeriesFilterDefinition]
type customChartsSectionChartsSeriesFilterDefinitionJSON struct {
	SourceType  apijson.Field
	DatasetIDs  apijson.Field
	ProjectIDs  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSeriesFilterDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSeriesFilterDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSeriesFilterDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [CustomChartsSectionChartsSeriesFilterDefinitionUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject],
// [CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset].
func (r CustomChartsSectionChartsSeriesFilterDefinition) AsUnion() CustomChartsSectionChartsSeriesFilterDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject]
// or [CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset].
type CustomChartsSectionChartsSeriesFilterDefinitionUnion interface {
	implementsCustomChartsSectionChartsSeriesFilterDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSeriesFilterDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset{}),
		},
	)
}

type CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject struct {
	ProjectIDs  []string                                                                                   `json:"project_ids" api:"required" format:"uuid"`
	SourceType  CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType `json:"source_type" api:"required"`
	RunFilter   string                                                                                     `json:"run_filter" api:"nullable"`
	TraceFilter string                                                                                     `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                                                                     `json:"tree_filter" api:"nullable"`
	JSON        customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON       `json:"-"`
}

// customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject]
type customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON struct {
	ProjectIDs  apijson.Field
	SourceType  apijson.Field
	RunFilter   apijson.Field
	TraceFilter apijson.Field
	TreeFilter  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProject) implementsCustomChartsSectionChartsSeriesFilterDefinition() {
}

type CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType string

const (
	CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType = "tracing_project"
)

func (r CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset struct {
	DatasetIDs []string                                                                            `json:"dataset_ids" api:"required" format:"uuid"`
	SourceType CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType `json:"source_type" api:"required"`
	JSON       customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON       `json:"-"`
}

// customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset]
type customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON struct {
	DatasetIDs  apijson.Field
	SourceType  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDataset) implementsCustomChartsSectionChartsSeriesFilterDefinition() {
}

type CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType string

const (
	CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType = "dataset"
)

func (r CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesFilterDefinitionSourceType string

const (
	CustomChartsSectionChartsSeriesFilterDefinitionSourceTypeTracingProject CustomChartsSectionChartsSeriesFilterDefinitionSourceType = "tracing_project"
	CustomChartsSectionChartsSeriesFilterDefinitionSourceTypeDataset        CustomChartsSectionChartsSeriesFilterDefinitionSourceType = "dataset"
)

func (r CustomChartsSectionChartsSeriesFilterDefinitionSourceType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesFilterDefinitionSourceTypeTracingProject, CustomChartsSectionChartsSeriesFilterDefinitionSourceTypeDataset:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesFilters struct {
	Filter      string                                     `json:"filter" api:"nullable"`
	Session     []string                                   `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                     `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                     `json:"tree_filter" api:"nullable"`
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
	Attribute CustomChartsSectionChartsSeriesGroupByAttribute `json:"attribute" api:"required"`
	MaxGroups int64                                           `json:"max_groups"`
	Path      string                                          `json:"path" api:"nullable"`
	SetBy     CustomChartsSectionChartsSeriesGroupBySetBy     `json:"set_by" api:"nullable"`
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

type CustomChartsSectionChartsSeriesGroupByDefinition struct {
	Attribute CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute `json:"attribute" api:"required"`
	Path      string                                                     `json:"path"`
	JSON      customChartsSectionChartsSeriesGroupByDefinitionJSON       `json:"-"`
	union     CustomChartsSectionChartsSeriesGroupByDefinitionsUnion
}

// customChartsSectionChartsSeriesGroupByDefinitionJSON contains the JSON metadata
// for the struct [CustomChartsSectionChartsSeriesGroupByDefinition]
type customChartsSectionChartsSeriesGroupByDefinitionJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSeriesGroupByDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSeriesGroupByDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSeriesGroupByDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [CustomChartsSectionChartsSeriesGroupByDefinitionsUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain],
// [CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex].
func (r CustomChartsSectionChartsSeriesGroupByDefinition) AsUnion() CustomChartsSectionChartsSeriesGroupByDefinitionsUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain] or
// [CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex].
type CustomChartsSectionChartsSeriesGroupByDefinitionsUnion interface {
	implementsCustomChartsSectionChartsSeriesGroupByDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSeriesGroupByDefinitionsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex{}),
		},
	)
}

type CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain struct {
	Attribute CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute `json:"attribute" api:"required"`
	JSON      customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON      `json:"-"`
}

// customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain]
type customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON struct {
	Attribute   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlain) implementsCustomChartsSectionChartsSeriesGroupByDefinition() {
}

type CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute string

const (
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName    CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "name"
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "run_type"
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag     CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "tag"
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "project"
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus  CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "status"
)

func (r CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName, CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType, CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag, CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject, CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex struct {
	Attribute CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute `json:"attribute" api:"required"`
	Path      string                                                                              `json:"path" api:"required"`
	JSON      customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON      `json:"-"`
}

// customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex]
type customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON struct {
	Attribute   apijson.Field
	Path        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplex) implementsCustomChartsSectionChartsSeriesGroupByDefinition() {
}

type CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute string

const (
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata      CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "metadata"
	CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "feedback_label"
)

func (r CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata, CustomChartsSectionChartsSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute string

const (
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeName          CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "name"
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeRunType       CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "run_type"
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeTag           CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "tag"
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeProject       CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "project"
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeStatus        CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "status"
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeMetadata      CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "metadata"
	CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeFeedbackLabel CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute = "feedback_label"
)

func (r CustomChartsSectionChartsSeriesGroupByDefinitionsAttribute) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeName, CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeRunType, CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeTag, CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeProject, CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeStatus, CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeMetadata, CustomChartsSectionChartsSeriesGroupByDefinitionsAttributeFeedbackLabel:
		return true
	}
	return false
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

type CustomChartsSectionChartsSeriesMetricDefinition struct {
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator].
	Denominator interface{}                                          `json:"denominator"`
	Field       CustomChartsSectionChartsSeriesMetricDefinitionField `json:"field"`
	Filter      string                                               `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator].
	Numerator interface{} `json:"numerator"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParams].
	Params interface{}                                         `json:"params"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionType `json:"type"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionJSON `json:"-"`
	union  CustomChartsSectionChartsSeriesMetricDefinitionUnion
}

// customChartsSectionChartsSeriesMetricDefinitionJSON contains the JSON metadata
// for the struct [CustomChartsSectionChartsSeriesMetricDefinition]
type customChartsSectionChartsSeriesMetricDefinitionJSON struct {
	Denominator apijson.Field
	Field       apijson.Field
	Filter      apijson.Field
	Numerator   apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSeriesMetricDefinitionJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSeriesMetricDefinition) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSeriesMetricDefinition{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [CustomChartsSectionChartsSeriesMetricDefinitionUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput].
func (r CustomChartsSectionChartsSeriesMetricDefinition) AsUnion() CustomChartsSectionChartsSeriesMetricDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile] or
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput].
type CustomChartsSectionChartsSeriesMetricDefinitionUnion interface {
	implementsCustomChartsSectionChartsSeriesMetricDefinition()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSeriesMetricDefinitionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput{}),
		},
	)
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount struct {
	Filter string                                                                    `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCount) implementsCustomChartsSectionChartsSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountTypeCount CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                      `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalar) implementsCustomChartsSectionChartsSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeSum CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMax CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMin CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                           `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentile) implementsCustomChartsSectionChartsSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParams struct {
	P    float64                                                                              `json:"p" api:"required"`
	JSON customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParams]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput struct {
	Denominator CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator `json:"denominator" api:"required"`
	Numerator   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator   `json:"numerator" api:"required"`
	Type        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType        `json:"type" api:"required"`
	JSON        customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON        `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON struct {
	Denominator apijson.Field
	Numerator   apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutput) implementsCustomChartsSectionChartsSeriesMetricDefinition() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField `json:"field"`
	Filter string                                                                                      `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams].
	Params interface{}                                                                                `json:"params"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType `json:"type"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON `json:"-"`
	union  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON struct {
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) AsUnion() CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar]
// or
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion interface {
	implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount struct {
	Filter string                                                                                                           `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount) implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                             `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar) implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                  `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile) implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                     `json:"p" api:"required"`
	JSON customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "count"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "sum"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "max"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "min"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "avg"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "percentile"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField `json:"field"`
	Filter string                                                                                    `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams].
	Params interface{}                                                                              `json:"params"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType `json:"type"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON `json:"-"`
	union  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON struct {
	Field       apijson.Field
	Filter      apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorJSON) RawJSON() string {
	return r.raw
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) UnmarshalJSON(data []byte) (err error) {
	*r = CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) AsUnion() CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar]
// or
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion interface {
	implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile{}),
		},
	)
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount struct {
	Filter string                                                                                                         `json:"filter" api:"nullable"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType `json:"type"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON struct {
	Filter      apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount) implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType = "count"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField `json:"field" api:"required"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType  `json:"type" api:"required"`
	Filter string                                                                                                           `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON  `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON struct {
	Field       apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar) implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "sum"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "max"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "min"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType = "avg"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeSum, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMax, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeMin, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile struct {
	Field  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField  `json:"field" api:"required"`
	Params CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams `json:"params" api:"required"`
	Type   CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType   `json:"type" api:"required"`
	Filter string                                                                                                                `json:"filter" api:"nullable"`
	JSON   customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON   `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON struct {
	Field       apijson.Field
	Params      apijson.Field
	Type        apijson.Field
	Filter      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileJSON) RawJSON() string {
	return r.raw
}

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile) implementsCustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator() {
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams struct {
	P    float64                                                                                                                   `json:"p" api:"required"`
	JSON customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON `json:"-"`
}

// customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON
// contains the JSON metadata for the struct
// [CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams]
type customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON struct {
	P           apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParams) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r customChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileParamsJSON) RawJSON() string {
	return r.raw
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType = "percentile"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount      CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "count"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "sum"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "max"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "min"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg        CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "avg"
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "percentile"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg, CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType = "ratio"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionCustomChartMetricRatioOutputTypeRatio:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionField string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionFieldLatencySeconds    CustomChartsSectionChartsSeriesMetricDefinitionField = "latency_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldFirstTokenSeconds CustomChartsSectionChartsSeriesMetricDefinitionField = "first_token_seconds"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldTotalTokens       CustomChartsSectionChartsSeriesMetricDefinitionField = "total_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldPromptTokens      CustomChartsSectionChartsSeriesMetricDefinitionField = "prompt_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldCompletionTokens  CustomChartsSectionChartsSeriesMetricDefinitionField = "completion_tokens"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldTotalCost         CustomChartsSectionChartsSeriesMetricDefinitionField = "total_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldPromptCost        CustomChartsSectionChartsSeriesMetricDefinitionField = "prompt_cost"
	CustomChartsSectionChartsSeriesMetricDefinitionFieldCompletionCost    CustomChartsSectionChartsSeriesMetricDefinitionField = "completion_cost"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionField) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionFieldLatencySeconds, CustomChartsSectionChartsSeriesMetricDefinitionFieldFirstTokenSeconds, CustomChartsSectionChartsSeriesMetricDefinitionFieldTotalTokens, CustomChartsSectionChartsSeriesMetricDefinitionFieldPromptTokens, CustomChartsSectionChartsSeriesMetricDefinitionFieldCompletionTokens, CustomChartsSectionChartsSeriesMetricDefinitionFieldTotalCost, CustomChartsSectionChartsSeriesMetricDefinitionFieldPromptCost, CustomChartsSectionChartsSeriesMetricDefinitionFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionChartsSeriesMetricDefinitionType string

const (
	CustomChartsSectionChartsSeriesMetricDefinitionTypeCount      CustomChartsSectionChartsSeriesMetricDefinitionType = "count"
	CustomChartsSectionChartsSeriesMetricDefinitionTypeSum        CustomChartsSectionChartsSeriesMetricDefinitionType = "sum"
	CustomChartsSectionChartsSeriesMetricDefinitionTypeMax        CustomChartsSectionChartsSeriesMetricDefinitionType = "max"
	CustomChartsSectionChartsSeriesMetricDefinitionTypeMin        CustomChartsSectionChartsSeriesMetricDefinitionType = "min"
	CustomChartsSectionChartsSeriesMetricDefinitionTypeAvg        CustomChartsSectionChartsSeriesMetricDefinitionType = "avg"
	CustomChartsSectionChartsSeriesMetricDefinitionTypePercentile CustomChartsSectionChartsSeriesMetricDefinitionType = "percentile"
	CustomChartsSectionChartsSeriesMetricDefinitionTypeRatio      CustomChartsSectionChartsSeriesMetricDefinitionType = "ratio"
)

func (r CustomChartsSectionChartsSeriesMetricDefinitionType) IsKnown() bool {
	switch r {
	case CustomChartsSectionChartsSeriesMetricDefinitionTypeCount, CustomChartsSectionChartsSeriesMetricDefinitionTypeSum, CustomChartsSectionChartsSeriesMetricDefinitionTypeMax, CustomChartsSectionChartsSeriesMetricDefinitionTypeMin, CustomChartsSectionChartsSeriesMetricDefinitionTypeAvg, CustomChartsSectionChartsSeriesMetricDefinitionTypePercentile, CustomChartsSectionChartsSeriesMetricDefinitionTypeRatio:
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
	Filter      string                                     `json:"filter" api:"nullable"`
	Session     []string                                   `json:"session" api:"nullable" format:"uuid"`
	TraceFilter string                                     `json:"trace_filter" api:"nullable"`
	TreeFilter  string                                     `json:"tree_filter" api:"nullable"`
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
	ID string `json:"id" api:"required" format:"uuid"`
	// Enum for custom chart types.
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

// Enum for custom chart types.
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
	Denominator interface{}                                                     `json:"denominator"`
	Field       CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField `json:"field"`
	Filter      string                                                          `json:"filter" api:"nullable"`
	// This field can have the runtime type of
	// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator].
	Numerator interface{} `json:"numerator"`
	// This field can have the runtime type of
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
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentile],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutput].
func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinition) AsUnion() CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar],
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
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalar{}),
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
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost:
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
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost:
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
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField `json:"field"`
	Filter string                                                                                                 `json:"filter" api:"nullable"`
	// This field can have the runtime type of
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
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile].
func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominator) AsUnion() CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar]
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
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentile{}),
		},
	)
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
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricScalarFieldCompletionCost:
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
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorCustomChartMetricPercentileFieldCompletionCost:
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

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorField string

const (
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
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "count"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeCount, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypeAvg, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputDenominatorTypePercentile:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator struct {
	Field  CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField `json:"field"`
	Filter string                                                                                               `json:"filter" api:"nullable"`
	// This field can have the runtime type of
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
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile].
func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumerator) AsUnion() CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorUnion {
	return r.union
}

// Union satisfied by
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount],
// [CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar]
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
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricCount{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalar{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentile{}),
		},
	)
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
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricScalarFieldCompletionCost:
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
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorCustomChartMetricPercentileFieldCompletionCost:
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

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorField string

const (
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
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "count"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType = "percentile"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeCount, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypeAvg, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionCustomChartMetricRatioOutputNumeratorTypePercentile:
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

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionField string

const (
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
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldLatencySeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldFirstTokenSeconds, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldTotalTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldPromptTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldCompletionTokens, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldTotalCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldPromptCost, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionFieldCompletionCost:
		return true
	}
	return false
}

type CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType string

const (
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeCount      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "count"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeSum        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "sum"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMax        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "max"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMin        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "min"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeAvg        CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "avg"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypePercentile CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "percentile"
	CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeRatio      CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType = "ratio"
)

func (r CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionType) IsKnown() bool {
	switch r {
	case CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeCount, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeSum, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMax, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeMin, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeAvg, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypePercentile, CustomChartsSectionSubSectionsChartsSeriesMetricDefinitionTypeRatio:
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
