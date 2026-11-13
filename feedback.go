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
	Configs *FeedbackConfigService
}

// NewFeedbackService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFeedbackService(opts ...option.RequestOption) (r *FeedbackService) {
	r = &FeedbackService{}
	r.Options = opts
	r.Tokens = NewFeedbackTokenService(opts...)
	r.Configs = NewFeedbackConfigService(opts...)
	return
}

// Create a new feedback.
//
// `session_id` identifies the tracing project the feedback belongs to. It is
// required unless the feedback is addressed by `address`, which names that project
// through an Agent environment that already exists.
func (r *FeedbackService) New(ctx context.Context, body FeedbackNewParams, opts ...option.RequestOption) (res *FeedbackSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/feedback"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a specific feedback.
func (r *FeedbackService) Get(ctx context.Context, feedbackID string, query FeedbackGetParams, opts ...option.RequestOption) (res *FeedbackSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if feedbackID == "" {
		err = errors.New("missing required feedback_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/feedback/%s", feedbackID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Replace an existing feedback entry with a new, modified entry.
func (r *FeedbackService) Update(ctx context.Context, feedbackID string, body FeedbackUpdateParams, opts ...option.RequestOption) (res *FeedbackSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if feedbackID == "" {
		err = errors.New("missing required feedback_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/feedback/%s", feedbackID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List feedback, newest first.
//
// At least one of run, key, or session is required.
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

// List feedback, newest first.
//
// At least one of run, key, or session is required.
func (r *FeedbackService) ListAutoPaging(ctx context.Context, query FeedbackListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[FeedbackSchema] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, query, opts...))
}

// Delete a feedback.
func (r *FeedbackService) Delete(ctx context.Context, feedbackID string, opts ...option.RequestOption) (res *FeedbackDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if feedbackID == "" {
		err = errors.New("missing required feedback_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/feedback/%s", feedbackID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// API feedback source.
type APIFeedbackSourceParam struct {
	Metadata param.Field[map[string]interface{}] `json:"metadata"`
	Type     param.Field[APIFeedbackSourceType]  `json:"type"`
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
	Metadata param.Field[map[string]interface{}] `json:"metadata"`
	Type     param.Field[AppFeedbackSourceType]  `json:"type"`
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
	Metadata param.Field[map[string]interface{}]     `json:"metadata"`
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
	// The feedback key: the metric, tag, or aspect being scored, for example
	// correctness.
	Key param.Field[string] `json:"key" api:"required"`
	// The ID to give the feedback. Defaults to a new random UUID.
	ID param.Field[string] `json:"id" format:"uuid"`
	// Addresses the tracing project through an Agent environment, as
	// lrn:agents/{id}/environments/{environment}, instead of session_id. The id is 1
	// to 63 lowercase ASCII letters, digits, or hyphens, starting with a letter and
	// ending with a letter or digit. The environment is local, development, staging,
	// or production, matched case-insensitively. Never combined with session_id. The
	// Agent and the environment must already exist; sending feedback does not create
	// them. Self-hosted deployments need LangSmith v0.18 or later.
	Address param.Field[string] `json:"address"`
	// A free-text comment explaining the feedback. Truncated at the feedback size
	// limit.
	Comment param.Field[string] `json:"comment"`
	// The ID of the comparative (pairwise) experiment this feedback was written for.
	ComparativeExperimentID param.Field[string] `json:"comparative_experiment_id" format:"uuid"`
	// The corrected output for the run, as a string or a JSON object. Dropped if it
	// exceeds the feedback size limit.
	Correction param.Field[FeedbackCreateSchemaCorrectionUnionParam] `json:"correction"`
	// When the feedback was created, as an ISO 8601 timestamp. Defaults to now.
	CreatedAt param.Field[time.Time] `json:"created_at" format:"date-time"`
	// Deprecated. Use extra.error instead. If both values are provided, error takes
	// precedence.
	//
	// Deprecated: deprecated
	Error param.Field[bool] `json:"error"`
	// Whether writing this feedback may upgrade the trace to extended retention.
	// Defaults to true.
	ExtendTraceRetention param.Field[bool] `json:"extend_trace_retention"`
	// Additional metadata to store with the feedback, as a JSON object.
	Extra param.Field[map[string]interface{}] `json:"extra"`
	// The scale the score is checked against. If the feedback key already has a
	// feedback config, this must match it exactly or the request is rejected; if it
	// has none, this becomes the key's feedback config.
	FeedbackConfig param.Field[FeedbackCreateSchemaFeedbackConfigParam] `json:"feedback_config"`
	// The ID of the feedback config this feedback was written against. It must be a
	// feedback config for this feedback key in the same workspace, or the feedback is
	// rejected. It doesn't change how the score is checked; feedback_config does that.
	// If omitted, the server may record the key's enforced feedback config, or the one
	// created for the key by this request.
	FeedbackConfigID param.Field[string] `json:"feedback_config_id" format:"uuid"`
	// Groups feedback written together, such as one comparison in a pairwise
	// experiment.
	FeedbackGroupID param.Field[string] `json:"feedback_group_id" format:"uuid"`
	// Where the feedback came from, as an object whose type is api, app, model, or
	// auto_eval.
	FeedbackSource param.Field[FeedbackCreateSchemaFeedbackSourceUnionParam] `json:"feedback_source"`
	// The ID of the thread the feedback is about. Requires session_id (or address) and
	// cannot be combined with run_id.
	FeedbackThreadID param.Field[string] `json:"feedback_thread_id"`
	// When the feedback was last modified, as an ISO 8601 timestamp. Defaults to now.
	ModifiedAt param.Field[time.Time] `json:"modified_at" format:"date-time"`
	// The ID of the run the feedback is about. Omit it for feedback on a whole tracing
	// project or a thread.
	RunID param.Field[string] `json:"run_id" format:"uuid"`
	// The numeric or boolean score. Numbers must be between -99999.9999 and 99999.9999
	// with at most 4 decimal places.
	Score param.Field[FeedbackCreateSchemaScoreUnionParam] `json:"score"`
	// Required unless the feedback is addressed by address. The ID of the tracing
	// project (session) the feedback belongs to.
	SessionID param.Field[string] `json:"session_id" format:"uuid"`
	// The start time of run_id, as an ISO 8601 timestamp. With session_id and
	// trace_id, it lets the server write run feedback without looking the run up.
	StartTime param.Field[time.Time] `json:"start_time" format:"date-time"`
	// The ID of the trace that contains run_id. With session_id and start_time, it
	// lets the server write run feedback without looking the run up.
	TraceID param.Field[string] `json:"trace_id" format:"uuid"`
	// A non-numeric value, such as a category label or a JSON object. Dropped if it
	// exceeds the feedback size limit.
	Value param.Field[FeedbackCreateSchemaValueUnionParam] `json:"value"`
}

func (r FeedbackCreateSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The corrected output for the run, as a string or a JSON object. Dropped if it
// exceeds the feedback size limit.
//
// Satisfied by [FeedbackCreateSchemaCorrectionMapParam], [shared.UnionString].
type FeedbackCreateSchemaCorrectionUnionParam interface {
	ImplementsFeedbackCreateSchemaCorrectionUnionParam()
}

type FeedbackCreateSchemaCorrectionMapParam map[string]interface{}

func (r FeedbackCreateSchemaCorrectionMapParam) ImplementsFeedbackCreateSchemaCorrectionUnionParam() {
}

// The scale the score is checked against. If the feedback key already has a
// feedback config, this must match it exactly or the request is rejected; if it
// has none, this becomes the key's feedback config.
type FeedbackCreateSchemaFeedbackConfigParam struct {
	// Enum for feedback types.
	Type param.Field[FeedbackCreateSchemaFeedbackConfigType] `json:"type" api:"required"`
	// Whether a reviewer can pick several categories. Categorical feedback only.
	AllowsMultipleCategories param.Field[bool]                                              `json:"allows_multiple_categories"`
	Categories               param.Field[[]FeedbackCreateSchemaFeedbackConfigCategoryParam] `json:"categories"`
	Max                      param.Field[float64]                                           `json:"max"`
	Min                      param.Field[float64]                                           `json:"min"`
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
	Value param.Field[float64] `json:"value" api:"required"`
	Label param.Field[string]  `json:"label"`
}

func (r FeedbackCreateSchemaFeedbackConfigCategoryParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Where the feedback came from, as an object whose type is api, app, model, or
// auto_eval.
type FeedbackCreateSchemaFeedbackSourceParam struct {
	Metadata param.Field[interface{}]                            `json:"metadata"`
	Type     param.Field[FeedbackCreateSchemaFeedbackSourceType] `json:"type"`
}

func (r FeedbackCreateSchemaFeedbackSourceParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FeedbackCreateSchemaFeedbackSourceParam) implementsFeedbackCreateSchemaFeedbackSourceUnionParam() {
}

// Where the feedback came from, as an object whose type is api, app, model, or
// auto_eval.
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

// The numeric or boolean score. Numbers must be between -99999.9999 and 99999.9999
// with at most 4 decimal places.
//
// Satisfied by [shared.UnionFloat], [shared.UnionBool].
type FeedbackCreateSchemaScoreUnionParam interface {
	ImplementsFeedbackCreateSchemaScoreUnionParam()
}

// A non-numeric value, such as a category label or a JSON object. Dropped if it
// exceeds the feedback size limit.
//
// Satisfied by [shared.UnionFloat], [shared.UnionBool], [shared.UnionString],
// [FeedbackCreateSchemaValueMapParam].
type FeedbackCreateSchemaValueUnionParam interface {
	ImplementsFeedbackCreateSchemaValueUnionParam()
}

type FeedbackCreateSchemaValueMapParam map[string]interface{}

func (r FeedbackCreateSchemaValueMapParam) ImplementsFeedbackCreateSchemaValueUnionParam() {}

// Feedback level: run for feedback on a run, session for feedback on a whole
// tracing project (session level requires a session filter).
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
	ID string `json:"id" api:"required" format:"uuid"`
	// The feedback key: the metric, tag, or aspect being scored, for example
	// correctness.
	Key string `json:"key" api:"required"`
	// A free-text comment explaining the feedback. Truncated at the feedback size
	// limit.
	Comment string `json:"comment" api:"nullable"`
	// The ID of the comparative (pairwise) experiment this feedback was written for.
	ComparativeExperimentID string `json:"comparative_experiment_id" api:"nullable" format:"uuid"`
	// The corrected output for the run, as a string or a JSON object. Dropped if it
	// exceeds the feedback size limit.
	Correction FeedbackSchemaCorrectionUnion `json:"correction" api:"nullable"`
	// When the feedback was created, as an ISO 8601 timestamp. Defaults to now.
	CreatedAt time.Time              `json:"created_at" format:"date-time"`
	Extra     map[string]interface{} `json:"extra" api:"nullable"`
	// Groups feedback written together, such as one comparison in a pairwise
	// experiment.
	FeedbackGroupID string `json:"feedback_group_id" api:"nullable" format:"uuid"`
	// The feedback source loaded from the database.
	FeedbackSource   FeedbackSchemaFeedbackSource `json:"feedback_source" api:"nullable"`
	FeedbackThreadID string                       `json:"feedback_thread_id" api:"nullable"`
	IsRoot           bool                         `json:"is_root"`
	// When the feedback was last modified, as an ISO 8601 timestamp. Defaults to now.
	ModifiedAt time.Time `json:"modified_at" format:"date-time"`
	// The ID of the run the feedback is about. Omit it for feedback on a whole tracing
	// project or a thread.
	RunID string `json:"run_id" api:"nullable" format:"uuid"`
	// The numeric or boolean score. Numbers must be between -99999.9999 and 99999.9999
	// with at most 4 decimal places.
	Score FeedbackSchemaScoreUnion `json:"score" api:"nullable"`
	// The ID of the tracing project (session) the feedback belongs to.
	SessionID string    `json:"session_id" api:"nullable" format:"uuid"`
	StartTime time.Time `json:"start_time" api:"nullable" format:"date-time"`
	TraceID   string    `json:"trace_id" api:"nullable" format:"uuid"`
	// A non-numeric value, such as a category label or a JSON object. Dropped if it
	// exceeds the feedback size limit.
	Value FeedbackSchemaValueUnion `json:"value" api:"nullable"`
	JSON  feedbackSchemaJSON       `json:"-"`
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
	IsRoot                  apijson.Field
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

// The corrected output for the run, as a string or a JSON object. Dropped if it
// exceeds the feedback size limit.
//
// Union satisfied by [FeedbackSchemaCorrectionMap] or [shared.UnionString].
type FeedbackSchemaCorrectionUnion interface {
	ImplementsFeedbackSchemaCorrectionUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FeedbackSchemaCorrectionUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(FeedbackSchemaCorrectionMap{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
	)
}

type FeedbackSchemaCorrectionMap map[string]interface{}

func (r FeedbackSchemaCorrectionMap) ImplementsFeedbackSchemaCorrectionUnion() {}

// The feedback source loaded from the database.
type FeedbackSchemaFeedbackSource struct {
	LsUserID string                           `json:"ls_user_id" api:"nullable" format:"uuid"`
	Metadata map[string]interface{}           `json:"metadata" api:"nullable"`
	Type     string                           `json:"type" api:"nullable"`
	UserID   string                           `json:"user_id" api:"nullable" format:"uuid"`
	UserName string                           `json:"user_name" api:"nullable"`
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

// The numeric or boolean score. Numbers must be between -99999.9999 and 99999.9999
// with at most 4 decimal places.
//
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

// A non-numeric value, such as a category label or a JSON object. Dropped if it
// exceeds the feedback size limit.
//
// Union satisfied by [shared.UnionFloat], [shared.UnionBool], [shared.UnionString]
// or [FeedbackSchemaValueMap].
type FeedbackSchemaValueUnion interface {
	ImplementsFeedbackSchemaValueUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FeedbackSchemaValueUnion)(nil)).Elem(),
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
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(FeedbackSchemaValueMap{}),
		},
	)
}

type FeedbackSchemaValueMap map[string]interface{}

func (r FeedbackSchemaValueMap) ImplementsFeedbackSchemaValueUnion() {}

// Model feedback source.
type ModelFeedbackSourceParam struct {
	Metadata param.Field[map[string]interface{}]  `json:"metadata"`
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
	FeedbackCreateSchema FeedbackCreateSchemaParam `json:"feedback_create_schema" api:"required"`
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type FeedbackUpdateParams struct {
	Comment        param.Field[string]                              `json:"comment"`
	Correction     param.Field[FeedbackUpdateParamsCorrectionUnion] `json:"correction"`
	FeedbackConfig param.Field[FeedbackUpdateParamsFeedbackConfig]  `json:"feedback_config"`
	Score          param.Field[FeedbackUpdateParamsScoreUnion]      `json:"score"`
	Value          param.Field[FeedbackUpdateParamsValueUnion]      `json:"value"`
}

func (r FeedbackUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [FeedbackUpdateParamsCorrectionMap], [shared.UnionString].
type FeedbackUpdateParamsCorrectionUnion interface {
	ImplementsFeedbackUpdateParamsCorrectionUnion()
}

type FeedbackUpdateParamsCorrectionMap map[string]interface{}

func (r FeedbackUpdateParamsCorrectionMap) ImplementsFeedbackUpdateParamsCorrectionUnion() {}

type FeedbackUpdateParamsFeedbackConfig struct {
	// Enum for feedback types.
	Type param.Field[FeedbackUpdateParamsFeedbackConfigType] `json:"type" api:"required"`
	// Whether a reviewer can pick several categories. Categorical feedback only.
	AllowsMultipleCategories param.Field[bool]                                         `json:"allows_multiple_categories"`
	Categories               param.Field[[]FeedbackUpdateParamsFeedbackConfigCategory] `json:"categories"`
	Max                      param.Field[float64]                                      `json:"max"`
	Min                      param.Field[float64]                                      `json:"min"`
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
	Value param.Field[float64] `json:"value" api:"required"`
	Label param.Field[string]  `json:"label"`
}

func (r FeedbackUpdateParamsFeedbackConfigCategory) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool].
type FeedbackUpdateParamsScoreUnion interface {
	ImplementsFeedbackUpdateParamsScoreUnion()
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool], [shared.UnionString],
// [FeedbackUpdateParamsValueMap].
type FeedbackUpdateParamsValueUnion interface {
	ImplementsFeedbackUpdateParamsValueUnion()
}

type FeedbackUpdateParamsValueMap map[string]interface{}

func (r FeedbackUpdateParamsValueMap) ImplementsFeedbackUpdateParamsValueUnion() {}

type FeedbackListParams struct {
	// Only feedback written for this comparative (pairwise) experiment.
	ComparativeExperimentID param.Field[string] `query:"comparative_experiment_id" format:"uuid"`
	// Only feedback on this thread. Requires exactly one session and cannot be
	// combined with run.
	FeedbackThreadID param.Field[string] `query:"feedback_thread_id"`
	// true returns only feedback with a non-empty comment; false only feedback without
	// one.
	HasComment param.Field[bool] `query:"has_comment"`
	// true returns only feedback with a score; false only feedback without one.
	HasScore param.Field[bool] `query:"has_score"`
	// Whether to add the author's name to each feedback's source.
	IncludeUserNames param.Field[bool] `query:"include_user_names"`
	// Only feedback with these keys.
	Key param.Field[[]string] `query:"key"`
	// run returns only feedback on a run; session returns only feedback on a whole
	// tracing project, and requires session.
	Level param.Field[FeedbackLevel] `query:"level"`
	// The maximum number of feedback to return, 1 to 100.
	Limit param.Field[int64] `query:"limit"`
	// Only feedback created at or before this ISO 8601 time.
	MaxCreatedAt param.Field[time.Time] `query:"max_created_at" format:"date-time"`
	// Only feedback created at or after this ISO 8601 time.
	MinCreatedAt param.Field[time.Time] `query:"min_created_at" format:"date-time"`
	// The number of feedback to skip, for paging.
	Offset param.Field[int64] `query:"offset"`
	// Only feedback on these run IDs.
	Run param.Field[FeedbackListParamsRunUnion] `query:"run" format:"uuid"`
	// Only feedback in these tracing project (session) IDs.
	Session param.Field[FeedbackListParamsSessionUnion] `query:"session" format:"uuid"`
	// Only feedback from these source types: api, model, app, or auto_eval.
	Source param.Field[[]SourceType] `query:"source"`
	// Only feedback written by these user IDs.
	User param.Field[[]string] `query:"user" format:"uuid"`
}

// URLQuery serializes [FeedbackListParams]'s query parameters as `url.Values`.
func (r FeedbackListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only feedback on these run IDs.
//
// Satisfied by [FeedbackListParamsRunArray], [shared.UnionString].
type FeedbackListParamsRunUnion interface {
	ImplementsFeedbackListParamsRunUnion()
}

type FeedbackListParamsRunArray []string

func (r FeedbackListParamsRunArray) ImplementsFeedbackListParamsRunUnion() {}

// Only feedback in these tracing project (session) IDs.
//
// Satisfied by [FeedbackListParamsSessionArray], [shared.UnionString].
type FeedbackListParamsSessionUnion interface {
	ImplementsFeedbackListParamsSessionUnion()
}

type FeedbackListParamsSessionArray []string

func (r FeedbackListParamsSessionArray) ImplementsFeedbackListParamsSessionUnion() {}
