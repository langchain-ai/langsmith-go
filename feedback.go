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

// FeedbackService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFeedbackService] method instead.
type FeedbackService struct {
	Options []option.RequestOption
	Tokens  *FeedbackTokenService
}

// NewFeedbackService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFeedbackService(opts ...option.RequestOption) (r *FeedbackService) {
	r = &FeedbackService{}
	r.Options = opts
	r.Tokens = NewFeedbackTokenService(opts...)
	return
}

// Create a new feedback.
func (r *FeedbackService) New(ctx context.Context, body FeedbackNewParams, opts ...option.RequestOption) (res *FeedbackSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/feedback"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a specific feedback.
func (r *FeedbackService) Get(ctx context.Context, feedbackID string, query FeedbackGetParams, opts ...option.RequestOption) (res *FeedbackSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if feedbackID == "" {
		err = errors.New("missing required feedback_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/feedback/%s", feedbackID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Replace an existing feedback entry with a new, modified entry.
func (r *FeedbackService) Update(ctx context.Context, feedbackID string, body FeedbackUpdateParams, opts ...option.RequestOption) (res *FeedbackSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if feedbackID == "" {
		err = errors.New("missing required feedback_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/feedback/%s", feedbackID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// List all Feedback by query params.
func (r *FeedbackService) List(ctx context.Context, query FeedbackListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[FeedbackSchema], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/feedback"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
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

// List all Feedback by query params.
func (r *FeedbackService) ListAutoPaging(ctx context.Context, query FeedbackListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[FeedbackSchema] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, query, opts...))
}

// Delete a feedback.
func (r *FeedbackService) Delete(ctx context.Context, feedbackID string, opts ...option.RequestOption) (res *FeedbackDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if feedbackID == "" {
		err = errors.New("missing required feedback_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/feedback/%s", feedbackID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// API feedback source.
type APIFeedbackSourceParam struct {
	Metadata param.Field[interface{}]           `json:"metadata"`
	Type     param.Field[APIFeedbackSourceType] `json:"type"`
}

func (r APIFeedbackSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r APIFeedbackSourceParam) implementsFeedbackCreateSchemaFeedbackSourceUnionParam() {}

type APIFeedbackSourceType string

const (
	APIFeedbackSourceTypeAPI APIFeedbackSourceType = "api"
)

func (r APIFeedbackSourceType) IsKnown() bool {
	switch r {
	case APIFeedbackSourceTypeAPI:
		return true
	}
	return false
}

// Feedback from the LangChainPlus App.
type AppFeedbackSourceParam struct {
	Metadata param.Field[interface{}]           `json:"metadata"`
	Type     param.Field[AppFeedbackSourceType] `json:"type"`
}

func (r AppFeedbackSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r AppFeedbackSourceParam) implementsFeedbackCreateSchemaFeedbackSourceUnionParam() {}

type AppFeedbackSourceType string

const (
	AppFeedbackSourceTypeApp AppFeedbackSourceType = "app"
)

func (r AppFeedbackSourceType) IsKnown() bool {
	switch r {
	case AppFeedbackSourceTypeApp:
		return true
	}
	return false
}

// Auto eval feedback source.
type AutoEvalFeedbackSourceParam struct {
	Metadata param.Field[interface{}]                `json:"metadata"`
	Type     param.Field[AutoEvalFeedbackSourceType] `json:"type"`
}

func (r AutoEvalFeedbackSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r AutoEvalFeedbackSourceParam) implementsFeedbackCreateSchemaFeedbackSourceUnionParam() {}

type AutoEvalFeedbackSourceType string

const (
	AutoEvalFeedbackSourceTypeAutoEval AutoEvalFeedbackSourceType = "auto_eval"
)

func (r AutoEvalFeedbackSourceType) IsKnown() bool {
	switch r {
	case AutoEvalFeedbackSourceTypeAutoEval:
		return true
	}
	return false
}

// Schema used for creating feedback.
type FeedbackCreateSchemaParam struct {
	Key                     param.Field[string]                                  `json:"key,required"`
	ID                      param.Field[string]                                  `json:"id" format:"uuid"`
	Comment                 param.Field[string]                                  `json:"comment"`
	ComparativeExperimentID param.Field[string]                                  `json:"comparative_experiment_id" format:"uuid"`
	Correction              param.Field[interface{}]                             `json:"correction"`
	CreatedAt               param.Field[time.Time]                               `json:"created_at" format:"date-time"`
	Error                   param.Field[bool]                                    `json:"error"`
	FeedbackConfig          param.Field[FeedbackCreateSchemaFeedbackConfigParam] `json:"feedback_config"`
	FeedbackGroupID         param.Field[string]                                  `json:"feedback_group_id" format:"uuid"`
	// Feedback from the LangChainPlus App.
	FeedbackSource param.Field[FeedbackCreateSchemaFeedbackSourceUnionParam] `json:"feedback_source"`
	ModifiedAt     param.Field[time.Time]                                    `json:"modified_at" format:"date-time"`
	RunID          param.Field[string]                                       `json:"run_id" format:"uuid"`
	Score          param.Field[FeedbackCreateSchemaScoreUnionParam]          `json:"score"`
	SessionID      param.Field[string]                                       `json:"session_id" format:"uuid"`
	TraceID        param.Field[string]                                       `json:"trace_id" format:"uuid"`
	Value          param.Field[interface{}]                                  `json:"value"`
}

func (r FeedbackCreateSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type FeedbackCreateSchemaFeedbackConfigParam struct {
	// Enum for feedback types.
	Type       param.Field[FeedbackCreateSchemaFeedbackConfigType]            `json:"type,required"`
	Categories param.Field[[]FeedbackCreateSchemaFeedbackConfigCategoryParam] `json:"categories"`
	Max        param.Field[float64]                                           `json:"max"`
	Min        param.Field[float64]                                           `json:"min"`
}

func (r FeedbackCreateSchemaFeedbackConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Enum for feedback types.
type FeedbackCreateSchemaFeedbackConfigType string

const (
	FeedbackCreateSchemaFeedbackConfigTypeContinuous  FeedbackCreateSchemaFeedbackConfigType = "continuous"
	FeedbackCreateSchemaFeedbackConfigTypeCategorical FeedbackCreateSchemaFeedbackConfigType = "categorical"
	FeedbackCreateSchemaFeedbackConfigTypeFreeform    FeedbackCreateSchemaFeedbackConfigType = "freeform"
)

func (r FeedbackCreateSchemaFeedbackConfigType) IsKnown() bool {
	switch r {
	case FeedbackCreateSchemaFeedbackConfigTypeContinuous, FeedbackCreateSchemaFeedbackConfigTypeCategorical, FeedbackCreateSchemaFeedbackConfigTypeFreeform:
		return true
	}
	return false
}

// Specific value and label pair for feedback
type FeedbackCreateSchemaFeedbackConfigCategoryParam struct {
	Value param.Field[float64] `json:"value,required"`
	Label param.Field[string]  `json:"label"`
}

func (r FeedbackCreateSchemaFeedbackConfigCategoryParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Feedback from the LangChainPlus App.
type FeedbackCreateSchemaFeedbackSourceParam struct {
	Metadata param.Field[interface{}]                            `json:"metadata"`
	Type     param.Field[FeedbackCreateSchemaFeedbackSourceType] `json:"type"`
}

func (r FeedbackCreateSchemaFeedbackSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FeedbackCreateSchemaFeedbackSourceParam) implementsFeedbackCreateSchemaFeedbackSourceUnionParam() {
}

// Feedback from the LangChainPlus App.
//
// Satisfied by [AppFeedbackSourceParam], [APIFeedbackSourceParam],
// [ModelFeedbackSourceParam], [AutoEvalFeedbackSourceParam],
// [FeedbackCreateSchemaFeedbackSourceParam].
type FeedbackCreateSchemaFeedbackSourceUnionParam interface {
	implementsFeedbackCreateSchemaFeedbackSourceUnionParam()
}

type FeedbackCreateSchemaFeedbackSourceType string

const (
	FeedbackCreateSchemaFeedbackSourceTypeApp      FeedbackCreateSchemaFeedbackSourceType = "app"
	FeedbackCreateSchemaFeedbackSourceTypeAPI      FeedbackCreateSchemaFeedbackSourceType = "api"
	FeedbackCreateSchemaFeedbackSourceTypeModel    FeedbackCreateSchemaFeedbackSourceType = "model"
	FeedbackCreateSchemaFeedbackSourceTypeAutoEval FeedbackCreateSchemaFeedbackSourceType = "auto_eval"
)

func (r FeedbackCreateSchemaFeedbackSourceType) IsKnown() bool {
	switch r {
	case FeedbackCreateSchemaFeedbackSourceTypeApp, FeedbackCreateSchemaFeedbackSourceTypeAPI, FeedbackCreateSchemaFeedbackSourceTypeModel, FeedbackCreateSchemaFeedbackSourceTypeAutoEval:
		return true
	}
	return false
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool].
type FeedbackCreateSchemaScoreUnionParam interface {
	ImplementsFeedbackCreateSchemaScoreUnionParam()
}

// Enum for feedback levels.
type FeedbackLevel string

const (
	FeedbackLevelRun     FeedbackLevel = "run"
	FeedbackLevelSession FeedbackLevel = "session"
)

func (r FeedbackLevel) IsKnown() bool {
	switch r {
	case FeedbackLevelRun, FeedbackLevelSession:
		return true
	}
	return false
}

// Schema for getting feedback.
type FeedbackSchema struct {
	ID                      string      `json:"id,required" format:"uuid"`
	Key                     string      `json:"key,required"`
	Comment                 string      `json:"comment,nullable"`
	ComparativeExperimentID string      `json:"comparative_experiment_id,nullable" format:"uuid"`
	Correction              interface{} `json:"correction,nullable"`
	CreatedAt               time.Time   `json:"created_at" format:"date-time"`
	Extra                   interface{} `json:"extra,nullable"`
	FeedbackGroupID         string      `json:"feedback_group_id,nullable" format:"uuid"`
	// The feedback source loaded from the database.
	FeedbackSource   FeedbackSchemaFeedbackSource `json:"feedback_source,nullable"`
	FeedbackThreadID string                       `json:"feedback_thread_id,nullable"`
	ModifiedAt       time.Time                    `json:"modified_at" format:"date-time"`
	RunID            string                       `json:"run_id,nullable" format:"uuid"`
	Score            FeedbackSchemaScoreUnion     `json:"score,nullable"`
	SessionID        string                       `json:"session_id,nullable" format:"uuid"`
	StartTime        time.Time                    `json:"start_time,nullable" format:"date-time"`
	TraceID          string                       `json:"trace_id,nullable" format:"uuid"`
	Value            interface{}                  `json:"value,nullable"`
	JSON             feedbackSchemaJSON           `json:"-"`
}

// feedbackSchemaJSON contains the JSON metadata for the struct [FeedbackSchema]
type feedbackSchemaJSON struct {
	ID                      apijson.Field
	Key                     apijson.Field
	Comment                 apijson.Field
	ComparativeExperimentID apijson.Field
	Correction              apijson.Field
	CreatedAt               apijson.Field
	Extra                   apijson.Field
	FeedbackGroupID         apijson.Field
	FeedbackSource          apijson.Field
	FeedbackThreadID        apijson.Field
	ModifiedAt              apijson.Field
	RunID                   apijson.Field
	Score                   apijson.Field
	SessionID               apijson.Field
	StartTime               apijson.Field
	TraceID                 apijson.Field
	Value                   apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *FeedbackSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r feedbackSchemaJSON) RawJSON() string {
	return r.raw
}

// The feedback source loaded from the database.
type FeedbackSchemaFeedbackSource struct {
	LsUserID string                           `json:"ls_user_id,nullable" format:"uuid"`
	Metadata interface{}                      `json:"metadata,nullable"`
	Type     string                           `json:"type,nullable"`
	UserID   string                           `json:"user_id,nullable" format:"uuid"`
	UserName string                           `json:"user_name,nullable"`
	JSON     feedbackSchemaFeedbackSourceJSON `json:"-"`
}

// feedbackSchemaFeedbackSourceJSON contains the JSON metadata for the struct
// [FeedbackSchemaFeedbackSource]
type feedbackSchemaFeedbackSourceJSON struct {
	LsUserID    apijson.Field
	Metadata    apijson.Field
	Type        apijson.Field
	UserID      apijson.Field
	UserName    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FeedbackSchemaFeedbackSource) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r feedbackSchemaFeedbackSourceJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [shared.UnionFloat] or [shared.UnionBool].
type FeedbackSchemaScoreUnion interface {
	ImplementsFeedbackSchemaScoreUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FeedbackSchemaScoreUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.Number,
			Type:       reflect.TypeOf(shared.UnionFloat(0)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.True,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.False,
			Type:       reflect.TypeOf(shared.UnionBool(false)),
		},
	)
}

// Model feedback source.
type ModelFeedbackSourceParam struct {
	Metadata param.Field[interface{}]             `json:"metadata"`
	Type     param.Field[ModelFeedbackSourceType] `json:"type"`
}

func (r ModelFeedbackSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ModelFeedbackSourceParam) implementsFeedbackCreateSchemaFeedbackSourceUnionParam() {}

type ModelFeedbackSourceType string

const (
	ModelFeedbackSourceTypeModel ModelFeedbackSourceType = "model"
)

func (r ModelFeedbackSourceType) IsKnown() bool {
	switch r {
	case ModelFeedbackSourceTypeModel:
		return true
	}
	return false
}

// Enum for feedback source types.
type SourceType string

const (
	SourceTypeAPI      SourceType = "api"
	SourceTypeModel    SourceType = "model"
	SourceTypeApp      SourceType = "app"
	SourceTypeAutoEval SourceType = "auto_eval"
)

func (r SourceType) IsKnown() bool {
	switch r {
	case SourceTypeAPI, SourceTypeModel, SourceTypeApp, SourceTypeAutoEval:
		return true
	}
	return false
}

type FeedbackDeleteResponse = interface{}

type FeedbackNewParams struct {
	// Schema used for creating feedback.
	FeedbackCreateSchema FeedbackCreateSchemaParam `json:"feedback_create_schema,required"`
}

func (r FeedbackNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.FeedbackCreateSchema)
}

type FeedbackGetParams struct {
	IncludeUserNames param.Field[bool] `query:"include_user_names"`
}

// URLQuery serializes [FeedbackGetParams]'s query parameters as `url.Values`.
func (r FeedbackGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FeedbackUpdateParams struct {
	Comment        param.Field[string]                             `json:"comment"`
	Correction     param.Field[interface{}]                        `json:"correction"`
	FeedbackConfig param.Field[FeedbackUpdateParamsFeedbackConfig] `json:"feedback_config"`
	Score          param.Field[FeedbackUpdateParamsScoreUnion]     `json:"score"`
	Value          param.Field[interface{}]                        `json:"value"`
}

func (r FeedbackUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type FeedbackUpdateParamsFeedbackConfig struct {
	// Enum for feedback types.
	Type       param.Field[FeedbackUpdateParamsFeedbackConfigType]       `json:"type,required"`
	Categories param.Field[[]FeedbackUpdateParamsFeedbackConfigCategory] `json:"categories"`
	Max        param.Field[float64]                                      `json:"max"`
	Min        param.Field[float64]                                      `json:"min"`
}

func (r FeedbackUpdateParamsFeedbackConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Enum for feedback types.
type FeedbackUpdateParamsFeedbackConfigType string

const (
	FeedbackUpdateParamsFeedbackConfigTypeContinuous  FeedbackUpdateParamsFeedbackConfigType = "continuous"
	FeedbackUpdateParamsFeedbackConfigTypeCategorical FeedbackUpdateParamsFeedbackConfigType = "categorical"
	FeedbackUpdateParamsFeedbackConfigTypeFreeform    FeedbackUpdateParamsFeedbackConfigType = "freeform"
)

func (r FeedbackUpdateParamsFeedbackConfigType) IsKnown() bool {
	switch r {
	case FeedbackUpdateParamsFeedbackConfigTypeContinuous, FeedbackUpdateParamsFeedbackConfigTypeCategorical, FeedbackUpdateParamsFeedbackConfigTypeFreeform:
		return true
	}
	return false
}

// Specific value and label pair for feedback
type FeedbackUpdateParamsFeedbackConfigCategory struct {
	Value param.Field[float64] `json:"value,required"`
	Label param.Field[string]  `json:"label"`
}

func (r FeedbackUpdateParamsFeedbackConfigCategory) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool].
type FeedbackUpdateParamsScoreUnion interface {
	ImplementsFeedbackUpdateParamsScoreUnion()
}

type FeedbackListParams struct {
	ComparativeExperimentID param.Field[string]   `query:"comparative_experiment_id" format:"uuid"`
	HasComment              param.Field[bool]     `query:"has_comment"`
	HasScore                param.Field[bool]     `query:"has_score"`
	IncludeUserNames        param.Field[bool]     `query:"include_user_names"`
	Key                     param.Field[[]string] `query:"key"`
	// Enum for feedback levels.
	Level        param.Field[FeedbackLevel] `query:"level"`
	Limit        param.Field[int64]         `query:"limit"`
	MaxCreatedAt param.Field[time.Time]     `query:"max_created_at" format:"date-time"`
	MinCreatedAt param.Field[time.Time]     `query:"min_created_at" format:"date-time"`
	Offset       param.Field[int64]         `query:"offset"`
	Run          param.Field[[]string]      `query:"run" format:"uuid"`
	Session      param.Field[[]string]      `query:"session" format:"uuid"`
	Source       param.Field[[]SourceType]  `query:"source"`
	User         param.Field[[]string]      `query:"user" format:"uuid"`
}

// URLQuery serializes [FeedbackListParams]'s query parameters as `url.Values`.
func (r FeedbackListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
