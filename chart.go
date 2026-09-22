// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/shared"
	"github.com/tidwall/gjson"
)

// ChartService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewChartService] method instead.
type ChartService struct {
	Options []option.RequestOption
}

// NewChartService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewChartService(opts ...option.RequestOption) (r *ChartService) {
	r = &ChartService{}
	r.Options = opts
	return
}

// Get a preview for a chart without actually creating it.
func (r *ChartService) Preview(ctx context.Context, body ChartPreviewParams, opts ...option.RequestOption) (res *ChartPreviewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/charts/preview"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type ChartPreviewResponse struct {
	Data []ChartPreviewResponseData `json:"data" api:"required"`
	JSON chartPreviewResponseJSON   `json:"-"`
}

// chartPreviewResponseJSON contains the JSON metadata for the struct
// [ChartPreviewResponse]
type chartPreviewResponseJSON struct {
	Data        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChartPreviewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chartPreviewResponseJSON) RawJSON() string {
	return r.raw
}

type ChartPreviewResponseData struct {
	SeriesID  string                             `json:"series_id" api:"required"`
	Timestamp time.Time                          `json:"timestamp" api:"required" format:"date-time"`
	Value     ChartPreviewResponseDataValueUnion `json:"value" api:"required,nullable"`
	Group     string                             `json:"group" api:"nullable"`
	JSON      chartPreviewResponseDataJSON       `json:"-"`
}

// chartPreviewResponseDataJSON contains the JSON metadata for the struct
// [ChartPreviewResponseData]
type chartPreviewResponseDataJSON struct {
	SeriesID    apijson.Field
	Timestamp   apijson.Field
	Value       apijson.Field
	Group       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ChartPreviewResponseData) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r chartPreviewResponseDataJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionFloat] or [ChartPreviewResponseDataValueMap].
type ChartPreviewResponseDataValueUnion interface {
	ImplementsChartPreviewResponseDataValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ChartPreviewResponseDataValueUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ChartPreviewResponseDataValueMap{}),
		},
	)
}

type ChartPreviewResponseDataValueMap map[string]interface{}

func (r ChartPreviewResponseDataValueMap) ImplementsChartPreviewResponseDataValueUnion() {}

type ChartPreviewParams struct {
	BucketInfo param.Field[ChartPreviewParamsBucketInfo] `json:"bucket_info" api:"required"`
	Chart      param.Field[ChartPreviewParamsChart]      `json:"chart" api:"required"`
}

func (r ChartPreviewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsBucketInfo struct {
	EndTime   param.Field[time.Time] `json:"end_time" format:"date-time"`
	OmitData  param.Field[bool]      `json:"omit_data"`
	StartTime param.Field[time.Time] `json:"start_time" format:"date-time"`
	// Timedelta input.
	Stride   param.Field[TimedeltaInputParam] `json:"stride"`
	Timezone param.Field[string]              `json:"timezone"`
}

func (r ChartPreviewParamsBucketInfo) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChart struct {
	Series        param.Field[[]ChartPreviewParamsChartSeries]      `json:"series" api:"required"`
	CommonFilters param.Field[ChartPreviewParamsChartCommonFilters] `json:"common_filters"`
}

func (r ChartPreviewParamsChart) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeries struct {
	ID               param.Field[string]                                             `json:"id" api:"required" format:"uuid"`
	Name             param.Field[string]                                             `json:"name" api:"required"`
	FeedbackKey      param.Field[string]                                             `json:"feedback_key"`
	FilterDefinition param.Field[ChartPreviewParamsChartSeriesFilterDefinitionUnion] `json:"filter_definition"`
	Filters          param.Field[ChartPreviewParamsChartSeriesFilters]               `json:"filters"`
	// Include additional information about where the group_by param was set.
	GroupBy            param.Field[ChartPreviewParamsChartSeriesGroupBy]                  `json:"group_by"`
	GroupByDefinitions param.Field[[]ChartPreviewParamsChartSeriesGroupByDefinitionUnion] `json:"group_by_definitions"`
	Metadata           param.Field[map[string]interface{}]                                `json:"metadata"`
	// Metrics you can chart. Feedback metrics are not available for
	// organization-scoped charts.
	Metric           param.Field[ChartPreviewParamsChartSeriesMetric]                `json:"metric"`
	MetricDefinition param.Field[ChartPreviewParamsChartSeriesMetricDefinitionUnion] `json:"metric_definition"`
	// LGP Metrics you can chart.
	ProjectMetric param.Field[ChartPreviewParamsChartSeriesProjectMetric] `json:"project_metric"`
	WorkspaceID   param.Field[string]                                     `json:"workspace_id" format:"uuid"`
}

func (r ChartPreviewParamsChartSeries) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesFilterDefinition struct {
	SourceType  param.Field[ChartPreviewParamsChartSeriesFilterDefinitionSourceType] `json:"source_type" api:"required"`
	DatasetIDs  param.Field[interface{}]                                             `json:"dataset_ids"`
	ProjectIDs  param.Field[interface{}]                                             `json:"project_ids"`
	RunFilter   param.Field[string]                                                  `json:"run_filter"`
	TraceFilter param.Field[string]                                                  `json:"trace_filter"`
	TreeFilter  param.Field[string]                                                  `json:"tree_filter"`
}

func (r ChartPreviewParamsChartSeriesFilterDefinition) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesFilterDefinition) implementsChartPreviewParamsChartSeriesFilterDefinitionUnion() {
}

// Satisfied by
// [ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProject],
// [ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDataset],
// [ChartPreviewParamsChartSeriesFilterDefinition].
type ChartPreviewParamsChartSeriesFilterDefinitionUnion interface {
	implementsChartPreviewParamsChartSeriesFilterDefinitionUnion()
}

type ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProject struct {
	ProjectIDs  param.Field[[]string]                                                                                 `json:"project_ids" api:"required" format:"uuid"`
	SourceType  param.Field[ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType] `json:"source_type" api:"required"`
	RunFilter   param.Field[string]                                                                                   `json:"run_filter"`
	TraceFilter param.Field[string]                                                                                   `json:"trace_filter"`
	TreeFilter  param.Field[string]                                                                                   `json:"tree_filter"`
}

func (r ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProject) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProject) implementsChartPreviewParamsChartSeriesFilterDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType string

const (
	ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType = "tracing_project"
)

func (r ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByTracingProjectSourceTypeTracingProject:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDataset struct {
	DatasetIDs param.Field[[]string]                                                                          `json:"dataset_ids" api:"required" format:"uuid"`
	SourceType param.Field[ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDatasetSourceType] `json:"source_type" api:"required"`
}

func (r ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDataset) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDataset) implementsChartPreviewParamsChartSeriesFilterDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDatasetSourceType string

const (
	ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDatasetSourceType = "dataset"
)

func (r ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDatasetSourceType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesFilterDefinitionCustomChartFilterByDatasetSourceTypeDataset:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesFilterDefinitionSourceType string

const (
	ChartPreviewParamsChartSeriesFilterDefinitionSourceTypeTracingProject ChartPreviewParamsChartSeriesFilterDefinitionSourceType = "tracing_project"
	ChartPreviewParamsChartSeriesFilterDefinitionSourceTypeDataset        ChartPreviewParamsChartSeriesFilterDefinitionSourceType = "dataset"
)

func (r ChartPreviewParamsChartSeriesFilterDefinitionSourceType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesFilterDefinitionSourceTypeTracingProject, ChartPreviewParamsChartSeriesFilterDefinitionSourceTypeDataset:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesFilters struct {
	Filter      param.Field[string]   `json:"filter"`
	Session     param.Field[[]string] `json:"session" format:"uuid"`
	TraceFilter param.Field[string]   `json:"trace_filter"`
	TreeFilter  param.Field[string]   `json:"tree_filter"`
}

func (r ChartPreviewParamsChartSeriesFilters) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Include additional information about where the group_by param was set.
type ChartPreviewParamsChartSeriesGroupBy struct {
	Attribute param.Field[ChartPreviewParamsChartSeriesGroupByAttribute] `json:"attribute" api:"required"`
	MaxGroups param.Field[int64]                                         `json:"max_groups"`
	Path      param.Field[string]                                        `json:"path"`
	SetBy     param.Field[ChartPreviewParamsChartSeriesGroupBySetBy]     `json:"set_by"`
}

func (r ChartPreviewParamsChartSeriesGroupBy) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesGroupByAttribute string

const (
	ChartPreviewParamsChartSeriesGroupByAttributeName     ChartPreviewParamsChartSeriesGroupByAttribute = "name"
	ChartPreviewParamsChartSeriesGroupByAttributeRunType  ChartPreviewParamsChartSeriesGroupByAttribute = "run_type"
	ChartPreviewParamsChartSeriesGroupByAttributeTag      ChartPreviewParamsChartSeriesGroupByAttribute = "tag"
	ChartPreviewParamsChartSeriesGroupByAttributeMetadata ChartPreviewParamsChartSeriesGroupByAttribute = "metadata"
)

func (r ChartPreviewParamsChartSeriesGroupByAttribute) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesGroupByAttributeName, ChartPreviewParamsChartSeriesGroupByAttributeRunType, ChartPreviewParamsChartSeriesGroupByAttributeTag, ChartPreviewParamsChartSeriesGroupByAttributeMetadata:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesGroupBySetBy string

const (
	ChartPreviewParamsChartSeriesGroupBySetBySection ChartPreviewParamsChartSeriesGroupBySetBy = "section"
	ChartPreviewParamsChartSeriesGroupBySetBySeries  ChartPreviewParamsChartSeriesGroupBySetBy = "series"
)

func (r ChartPreviewParamsChartSeriesGroupBySetBy) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesGroupBySetBySection, ChartPreviewParamsChartSeriesGroupBySetBySeries:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesGroupByDefinition struct {
	Attribute param.Field[ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute] `json:"attribute" api:"required"`
	Path      param.Field[string]                                                   `json:"path"`
}

func (r ChartPreviewParamsChartSeriesGroupByDefinition) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesGroupByDefinition) implementsChartPreviewParamsChartSeriesGroupByDefinitionUnion() {
}

// Satisfied by
// [ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlain],
// [ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplex],
// [ChartPreviewParamsChartSeriesGroupByDefinition].
type ChartPreviewParamsChartSeriesGroupByDefinitionUnion interface {
	implementsChartPreviewParamsChartSeriesGroupByDefinitionUnion()
}

type ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlain struct {
	Attribute param.Field[ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute] `json:"attribute" api:"required"`
}

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlain) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlain) implementsChartPreviewParamsChartSeriesGroupByDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute string

const (
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName    ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "name"
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "run_type"
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag     ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "tag"
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "project"
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus  ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute = "status"
)

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttribute) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeName, ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeRunType, ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeTag, ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeProject, ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByPlainAttributeStatus:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplex struct {
	Attribute param.Field[ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute] `json:"attribute" api:"required"`
	Path      param.Field[string]                                                                            `json:"path" api:"required"`
}

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplex) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplex) implementsChartPreviewParamsChartSeriesGroupByDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute string

const (
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata      ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "metadata"
	ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute = "feedback_label"
)

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttribute) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeMetadata, ChartPreviewParamsChartSeriesGroupByDefinitionsCustomChartGroupByComplexAttributeFeedbackLabel:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute string

const (
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeName          ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "name"
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeRunType       ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "run_type"
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeTag           ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "tag"
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeProject       ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "project"
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeStatus        ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "status"
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeMetadata      ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "metadata"
	ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeFeedbackLabel ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute = "feedback_label"
)

func (r ChartPreviewParamsChartSeriesGroupByDefinitionsAttribute) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeName, ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeRunType, ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeTag, ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeProject, ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeStatus, ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeMetadata, ChartPreviewParamsChartSeriesGroupByDefinitionsAttributeFeedbackLabel:
		return true
	}
	return false
}

// Metrics you can chart. Feedback metrics are not available for
// organization-scoped charts.
type ChartPreviewParamsChartSeriesMetric string

const (
	ChartPreviewParamsChartSeriesMetricRunCount            ChartPreviewParamsChartSeriesMetric = "run_count"
	ChartPreviewParamsChartSeriesMetricLatencyP50          ChartPreviewParamsChartSeriesMetric = "latency_p50"
	ChartPreviewParamsChartSeriesMetricLatencyP99          ChartPreviewParamsChartSeriesMetric = "latency_p99"
	ChartPreviewParamsChartSeriesMetricLatencyAvg          ChartPreviewParamsChartSeriesMetric = "latency_avg"
	ChartPreviewParamsChartSeriesMetricFirstTokenP50       ChartPreviewParamsChartSeriesMetric = "first_token_p50"
	ChartPreviewParamsChartSeriesMetricFirstTokenP99       ChartPreviewParamsChartSeriesMetric = "first_token_p99"
	ChartPreviewParamsChartSeriesMetricTotalTokens         ChartPreviewParamsChartSeriesMetric = "total_tokens"
	ChartPreviewParamsChartSeriesMetricPromptTokens        ChartPreviewParamsChartSeriesMetric = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricCompletionTokens    ChartPreviewParamsChartSeriesMetric = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricMedianTokens        ChartPreviewParamsChartSeriesMetric = "median_tokens"
	ChartPreviewParamsChartSeriesMetricCompletionTokensP50 ChartPreviewParamsChartSeriesMetric = "completion_tokens_p50"
	ChartPreviewParamsChartSeriesMetricPromptTokensP50     ChartPreviewParamsChartSeriesMetric = "prompt_tokens_p50"
	ChartPreviewParamsChartSeriesMetricTokensP99           ChartPreviewParamsChartSeriesMetric = "tokens_p99"
	ChartPreviewParamsChartSeriesMetricCompletionTokensP99 ChartPreviewParamsChartSeriesMetric = "completion_tokens_p99"
	ChartPreviewParamsChartSeriesMetricPromptTokensP99     ChartPreviewParamsChartSeriesMetric = "prompt_tokens_p99"
	ChartPreviewParamsChartSeriesMetricFeedback            ChartPreviewParamsChartSeriesMetric = "feedback"
	ChartPreviewParamsChartSeriesMetricFeedbackScoreAvg    ChartPreviewParamsChartSeriesMetric = "feedback_score_avg"
	ChartPreviewParamsChartSeriesMetricFeedbackValues      ChartPreviewParamsChartSeriesMetric = "feedback_values"
	ChartPreviewParamsChartSeriesMetricTotalCost           ChartPreviewParamsChartSeriesMetric = "total_cost"
	ChartPreviewParamsChartSeriesMetricPromptCost          ChartPreviewParamsChartSeriesMetric = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricCompletionCost      ChartPreviewParamsChartSeriesMetric = "completion_cost"
	ChartPreviewParamsChartSeriesMetricErrorRate           ChartPreviewParamsChartSeriesMetric = "error_rate"
	ChartPreviewParamsChartSeriesMetricStreamingRate       ChartPreviewParamsChartSeriesMetric = "streaming_rate"
	ChartPreviewParamsChartSeriesMetricCostP50             ChartPreviewParamsChartSeriesMetric = "cost_p50"
	ChartPreviewParamsChartSeriesMetricCostP99             ChartPreviewParamsChartSeriesMetric = "cost_p99"
)

func (r ChartPreviewParamsChartSeriesMetric) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricRunCount, ChartPreviewParamsChartSeriesMetricLatencyP50, ChartPreviewParamsChartSeriesMetricLatencyP99, ChartPreviewParamsChartSeriesMetricLatencyAvg, ChartPreviewParamsChartSeriesMetricFirstTokenP50, ChartPreviewParamsChartSeriesMetricFirstTokenP99, ChartPreviewParamsChartSeriesMetricTotalTokens, ChartPreviewParamsChartSeriesMetricPromptTokens, ChartPreviewParamsChartSeriesMetricCompletionTokens, ChartPreviewParamsChartSeriesMetricMedianTokens, ChartPreviewParamsChartSeriesMetricCompletionTokensP50, ChartPreviewParamsChartSeriesMetricPromptTokensP50, ChartPreviewParamsChartSeriesMetricTokensP99, ChartPreviewParamsChartSeriesMetricCompletionTokensP99, ChartPreviewParamsChartSeriesMetricPromptTokensP99, ChartPreviewParamsChartSeriesMetricFeedback, ChartPreviewParamsChartSeriesMetricFeedbackScoreAvg, ChartPreviewParamsChartSeriesMetricFeedbackValues, ChartPreviewParamsChartSeriesMetricTotalCost, ChartPreviewParamsChartSeriesMetricPromptCost, ChartPreviewParamsChartSeriesMetricCompletionCost, ChartPreviewParamsChartSeriesMetricErrorRate, ChartPreviewParamsChartSeriesMetricStreamingRate, ChartPreviewParamsChartSeriesMetricCostP50, ChartPreviewParamsChartSeriesMetricCostP99:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinition struct {
	Denominator param.Field[interface{}]                                         `json:"denominator"`
	Entity      param.Field[ChartPreviewParamsChartSeriesMetricDefinitionEntity] `json:"entity"`
	Field       param.Field[ChartPreviewParamsChartSeriesMetricDefinitionField]  `json:"field"`
	Filter      param.Field[string]                                              `json:"filter"`
	Numerator   param.Field[interface{}]                                         `json:"numerator"`
	Params      param.Field[interface{}]                                         `json:"params"`
	Type        param.Field[ChartPreviewParamsChartSeriesMetricDefinitionType]   `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinition) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinition) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

// Satisfied by
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetric],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCount],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalar],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentile],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInput],
// [ChartPreviewParamsChartSeriesMetricDefinition].
type ChartPreviewParamsChartSeriesMetricDefinitionUnion interface {
	implementsChartPreviewParamsChartSeriesMetricDefinitionUnion()
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetric struct {
	Entity param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity] `json:"entity" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricParams] `json:"params" api:"required"`
	Filter param.Field[string]                                                                            `json:"filter"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricType]   `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetric) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetric) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricParams struct {
	FeedbackKey param.Field[string] `json:"feedback_key" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricType = "count"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCount struct {
	Filter param.Field[string]                                                                  `json:"filter"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCountType] `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCount) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCount) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCountType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCountTypeCount ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCountType = "count"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                  `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalar) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey param.Field[string] `json:"feedback_key" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeSum ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalar struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField] `json:"field" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType]  `json:"type" api:"required"`
	Filter param.Field[string]                                                                    `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalar) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalar) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "completion_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldCompletionCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeSum ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeMax ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeMin ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeAvg ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType = "avg"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                      `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentile) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey param.Field[string]  `json:"feedback_key" api:"required"`
	P           param.Field[float64] `json:"p" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentile struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                         `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentile) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentile) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "completion_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldCompletionCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileParams struct {
	P param.Field[float64] `json:"p" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInput struct {
	Denominator param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion] `json:"denominator" api:"required"`
	Numerator   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion]   `json:"numerator" api:"required"`
	Type        param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputType]             `json:"type" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInput) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInput) implementsChartPreviewParamsChartSeriesMetricDefinitionUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominator struct {
	Entity param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorEntity] `json:"entity"`
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField]  `json:"field"`
	Filter param.Field[string]                                                                                    `json:"filter"`
	Params param.Field[interface{}]                                                                               `json:"params"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType]   `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominator) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominator) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

// Satisfied by
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetric],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCount],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalar],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalar],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentile],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentile],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominator].
type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion interface {
	implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion()
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetric struct {
	Entity param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricEntity] `json:"entity" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricParams] `json:"params" api:"required"`
	Filter param.Field[string]                                                                                                                  `json:"filter"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricType]   `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetric) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetric) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricEntity string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricEntityFeedback ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey param.Field[string] `json:"feedback_key" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricTypeCount ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricType = "count"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCount struct {
	Filter param.Field[string]                                                                                                        `json:"filter"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCountType] `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCount) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCount) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCountType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCountTypeCount ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCountType = "count"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalar struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                                        `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalar) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalar) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey param.Field[string] `json:"feedback_key" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeSum ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalar struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField] `json:"field" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType]  `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                          `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalar) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalar) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "completion_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldCompletionCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeSum ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeMax ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeMin ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeAvg ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType = "avg"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentile struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                                            `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentile) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentile) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey param.Field[string]  `json:"feedback_key" api:"required"`
	P           param.Field[float64] `json:"p" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentile struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                               `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentile) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentile) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "completion_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldCompletionCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileParams struct {
	P param.Field[float64] `json:"p" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorEntity string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorEntityFeedback ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorEntity = "feedback"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorEntity) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorEntityFeedback:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "feedback_score"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField = "completion_cost"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldFeedbackScore, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorFieldCompletionCost:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeCount      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType = "count"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeSum        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeMax        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeMin        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeAvg        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType = "avg"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeCount, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypeAvg, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputDenominatorTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumerator struct {
	Entity param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorEntity] `json:"entity"`
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField]  `json:"field"`
	Filter param.Field[string]                                                                                  `json:"filter"`
	Params param.Field[interface{}]                                                                             `json:"params"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType]   `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumerator) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumerator) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

// Satisfied by
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetric],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCount],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalar],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalar],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentile],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentile],
// [ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumerator].
type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion interface {
	implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion()
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetric struct {
	Entity param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricEntity] `json:"entity" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricParams] `json:"params" api:"required"`
	Filter param.Field[string]                                                                                                                `json:"filter"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricType]   `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetric) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetric) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricEntity string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricEntityFeedback ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricEntity = "feedback"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricEntity) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricEntityFeedback:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricParams struct {
	FeedbackKey param.Field[string] `json:"feedback_key" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricTypeCount ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricType = "count"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackCountMetricTypeCount:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCount struct {
	Filter param.Field[string]                                                                                                      `json:"filter"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCountType] `json:"type"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCount) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCount) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCountType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCountTypeCount ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCountType = "count"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCountType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricCountTypeCount:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalar struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                                      `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalar) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalar) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarParams struct {
	FeedbackKey param.Field[string] `json:"feedback_key" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeSum ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType = "avg"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricScalarTypeAvg:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalar struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField] `json:"field" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType]  `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                        `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalar) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalar) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "completion_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldCompletionCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeSum ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeMax ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeMin ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeAvg ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType = "avg"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricScalarTypeAvg:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentile struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                                          `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentile) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentile) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileParams struct {
	FeedbackKey param.Field[string]  `json:"feedback_key" api:"required"`
	P           param.Field[float64] `json:"p" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartFeedbackScoreMetricPercentileTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentile struct {
	Field  param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField]  `json:"field" api:"required"`
	Params param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileParams] `json:"params" api:"required"`
	Type   param.Field[ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileType]   `json:"type" api:"required"`
	Filter param.Field[string]                                                                                                             `json:"filter"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentile) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentile) implementsChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorUnion() {
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "completion_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField = "feedback_score"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldCompletionCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileFieldFeedbackScore:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileParams struct {
	P param.Field[float64] `json:"p" api:"required"`
}

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorCustomChartMetricPercentileTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorEntity string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorEntityFeedback ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorEntity = "feedback"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorEntity) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorEntityFeedback:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "feedback_score"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField = "completion_cost"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldFeedbackScore, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorFieldCompletionCost:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeCount      ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType = "count"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeSum        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeMax        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeMin        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeAvg        ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType = "avg"
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType = "percentile"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeCount, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypeAvg, ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputNumeratorTypePercentile:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputTypeRatio ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputType = "ratio"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionCustomChartMetricRatioInputTypeRatio:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionEntity string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionEntityFeedback ChartPreviewParamsChartSeriesMetricDefinitionEntity = "feedback"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionEntity) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionEntityFeedback:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionField string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionFieldFeedbackScore     ChartPreviewParamsChartSeriesMetricDefinitionField = "feedback_score"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldLatencySeconds    ChartPreviewParamsChartSeriesMetricDefinitionField = "latency_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldFirstTokenSeconds ChartPreviewParamsChartSeriesMetricDefinitionField = "first_token_seconds"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldTotalTokens       ChartPreviewParamsChartSeriesMetricDefinitionField = "total_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldPromptTokens      ChartPreviewParamsChartSeriesMetricDefinitionField = "prompt_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldCompletionTokens  ChartPreviewParamsChartSeriesMetricDefinitionField = "completion_tokens"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldTotalCost         ChartPreviewParamsChartSeriesMetricDefinitionField = "total_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldPromptCost        ChartPreviewParamsChartSeriesMetricDefinitionField = "prompt_cost"
	ChartPreviewParamsChartSeriesMetricDefinitionFieldCompletionCost    ChartPreviewParamsChartSeriesMetricDefinitionField = "completion_cost"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionField) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionFieldFeedbackScore, ChartPreviewParamsChartSeriesMetricDefinitionFieldLatencySeconds, ChartPreviewParamsChartSeriesMetricDefinitionFieldFirstTokenSeconds, ChartPreviewParamsChartSeriesMetricDefinitionFieldTotalTokens, ChartPreviewParamsChartSeriesMetricDefinitionFieldPromptTokens, ChartPreviewParamsChartSeriesMetricDefinitionFieldCompletionTokens, ChartPreviewParamsChartSeriesMetricDefinitionFieldTotalCost, ChartPreviewParamsChartSeriesMetricDefinitionFieldPromptCost, ChartPreviewParamsChartSeriesMetricDefinitionFieldCompletionCost:
		return true
	}
	return false
}

type ChartPreviewParamsChartSeriesMetricDefinitionType string

const (
	ChartPreviewParamsChartSeriesMetricDefinitionTypeCount      ChartPreviewParamsChartSeriesMetricDefinitionType = "count"
	ChartPreviewParamsChartSeriesMetricDefinitionTypeSum        ChartPreviewParamsChartSeriesMetricDefinitionType = "sum"
	ChartPreviewParamsChartSeriesMetricDefinitionTypeMax        ChartPreviewParamsChartSeriesMetricDefinitionType = "max"
	ChartPreviewParamsChartSeriesMetricDefinitionTypeMin        ChartPreviewParamsChartSeriesMetricDefinitionType = "min"
	ChartPreviewParamsChartSeriesMetricDefinitionTypeAvg        ChartPreviewParamsChartSeriesMetricDefinitionType = "avg"
	ChartPreviewParamsChartSeriesMetricDefinitionTypePercentile ChartPreviewParamsChartSeriesMetricDefinitionType = "percentile"
	ChartPreviewParamsChartSeriesMetricDefinitionTypeRatio      ChartPreviewParamsChartSeriesMetricDefinitionType = "ratio"
)

func (r ChartPreviewParamsChartSeriesMetricDefinitionType) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesMetricDefinitionTypeCount, ChartPreviewParamsChartSeriesMetricDefinitionTypeSum, ChartPreviewParamsChartSeriesMetricDefinitionTypeMax, ChartPreviewParamsChartSeriesMetricDefinitionTypeMin, ChartPreviewParamsChartSeriesMetricDefinitionTypeAvg, ChartPreviewParamsChartSeriesMetricDefinitionTypePercentile, ChartPreviewParamsChartSeriesMetricDefinitionTypeRatio:
		return true
	}
	return false
}

// LGP Metrics you can chart.
type ChartPreviewParamsChartSeriesProjectMetric string

const (
	ChartPreviewParamsChartSeriesProjectMetricMemoryUsage             ChartPreviewParamsChartSeriesProjectMetric = "memory_usage"
	ChartPreviewParamsChartSeriesProjectMetricCPUUsage                ChartPreviewParamsChartSeriesProjectMetric = "cpu_usage"
	ChartPreviewParamsChartSeriesProjectMetricDiskUsage               ChartPreviewParamsChartSeriesProjectMetric = "disk_usage"
	ChartPreviewParamsChartSeriesProjectMetricRestartCount            ChartPreviewParamsChartSeriesProjectMetric = "restart_count"
	ChartPreviewParamsChartSeriesProjectMetricReplicaCount            ChartPreviewParamsChartSeriesProjectMetric = "replica_count"
	ChartPreviewParamsChartSeriesProjectMetricWorkerCount             ChartPreviewParamsChartSeriesProjectMetric = "worker_count"
	ChartPreviewParamsChartSeriesProjectMetricLgRunCount              ChartPreviewParamsChartSeriesProjectMetric = "lg_run_count"
	ChartPreviewParamsChartSeriesProjectMetricResponsesPerSecond      ChartPreviewParamsChartSeriesProjectMetric = "responses_per_second"
	ChartPreviewParamsChartSeriesProjectMetricErrorResponsesPerSecond ChartPreviewParamsChartSeriesProjectMetric = "error_responses_per_second"
	ChartPreviewParamsChartSeriesProjectMetricP95Latency              ChartPreviewParamsChartSeriesProjectMetric = "p95_latency"
	ChartPreviewParamsChartSeriesProjectMetricRunQueueWaitTime        ChartPreviewParamsChartSeriesProjectMetric = "run_queue_wait_time"
)

func (r ChartPreviewParamsChartSeriesProjectMetric) IsKnown() bool {
	switch r {
	case ChartPreviewParamsChartSeriesProjectMetricMemoryUsage, ChartPreviewParamsChartSeriesProjectMetricCPUUsage, ChartPreviewParamsChartSeriesProjectMetricDiskUsage, ChartPreviewParamsChartSeriesProjectMetricRestartCount, ChartPreviewParamsChartSeriesProjectMetricReplicaCount, ChartPreviewParamsChartSeriesProjectMetricWorkerCount, ChartPreviewParamsChartSeriesProjectMetricLgRunCount, ChartPreviewParamsChartSeriesProjectMetricResponsesPerSecond, ChartPreviewParamsChartSeriesProjectMetricErrorResponsesPerSecond, ChartPreviewParamsChartSeriesProjectMetricP95Latency, ChartPreviewParamsChartSeriesProjectMetricRunQueueWaitTime:
		return true
	}
	return false
}

type ChartPreviewParamsChartCommonFilters struct {
	Filter      param.Field[string]   `json:"filter"`
	Session     param.Field[[]string] `json:"session" format:"uuid"`
	TraceFilter param.Field[string]   `json:"trace_filter"`
	TreeFilter  param.Field[string]   `json:"tree_filter"`
}

func (r ChartPreviewParamsChartCommonFilters) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
