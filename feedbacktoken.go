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
	"github.com/tidwall/gjson"
)

// FeedbackTokenService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFeedbackTokenService] method instead.
type FeedbackTokenService struct {
	Options []option.RequestOption
}

// NewFeedbackTokenService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFeedbackTokenService(opts ...option.RequestOption) (r *FeedbackTokenService) {
	r = &FeedbackTokenService{}
	r.Options = opts
	return
}

// Create a new feedback ingest token.
func (r *FeedbackTokenService) New(ctx context.Context, body FeedbackTokenNewParams, opts ...option.RequestOption) (res *FeedbackTokenNewResponseUnion, err error) {
	var env apijson.UnionUnmarshaler[FeedbackTokenNewResponseUnion]
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/feedback/tokens"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &env, opts...)
	if err != nil {
		return
	}
	res = &env.Value
	return
}

// Create a new feedback with a token.
func (r *FeedbackTokenService) Get(ctx context.Context, token string, query FeedbackTokenGetParams, opts ...option.RequestOption) (res *FeedbackTokenGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if token == "" {
		err = errors.New("missing required token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/feedback/tokens/%s", token)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Create a new feedback with a token.
func (r *FeedbackTokenService) Update(ctx context.Context, token string, body FeedbackTokenUpdateParams, opts ...option.RequestOption) (res *FeedbackTokenUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if token == "" {
		err = errors.New("missing required token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/feedback/tokens/%s", token)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// List all feedback ingest tokens for a run.
func (r *FeedbackTokenService) List(ctx context.Context, query FeedbackTokenListParams, opts ...option.RequestOption) (res *[]FeedbackIngestTokenSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/feedback/tokens"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Feedback ingest token create schema.
type FeedbackIngestTokenCreateSchemaParam struct {
	FeedbackKey param.Field[string]    `json:"feedback_key,required"`
	RunID       param.Field[string]    `json:"run_id,required" format:"uuid"`
	ExpiresAt   param.Field[time.Time] `json:"expires_at" format:"date-time"`
	// Timedelta input.
	ExpiresIn      param.Field[TimedeltaInputParam]                                `json:"expires_in"`
	FeedbackConfig param.Field[FeedbackIngestTokenCreateSchemaFeedbackConfigParam] `json:"feedback_config"`
}

func (r FeedbackIngestTokenCreateSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FeedbackIngestTokenCreateSchemaParam) implementsFeedbackTokenNewParamsBodyUnion() {}

type FeedbackIngestTokenCreateSchemaFeedbackConfigParam struct {
	// Enum for feedback types.
	Type       param.Field[FeedbackIngestTokenCreateSchemaFeedbackConfigType]            `json:"type,required"`
	Categories param.Field[[]FeedbackIngestTokenCreateSchemaFeedbackConfigCategoryParam] `json:"categories"`
	Max        param.Field[float64]                                                      `json:"max"`
	Min        param.Field[float64]                                                      `json:"min"`
}

func (r FeedbackIngestTokenCreateSchemaFeedbackConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Enum for feedback types.
type FeedbackIngestTokenCreateSchemaFeedbackConfigType string

const (
	FeedbackIngestTokenCreateSchemaFeedbackConfigTypeContinuous  FeedbackIngestTokenCreateSchemaFeedbackConfigType = "continuous"
	FeedbackIngestTokenCreateSchemaFeedbackConfigTypeCategorical FeedbackIngestTokenCreateSchemaFeedbackConfigType = "categorical"
	FeedbackIngestTokenCreateSchemaFeedbackConfigTypeFreeform    FeedbackIngestTokenCreateSchemaFeedbackConfigType = "freeform"
)

func (r FeedbackIngestTokenCreateSchemaFeedbackConfigType) IsKnown() bool {
	switch r {
	case FeedbackIngestTokenCreateSchemaFeedbackConfigTypeContinuous, FeedbackIngestTokenCreateSchemaFeedbackConfigTypeCategorical, FeedbackIngestTokenCreateSchemaFeedbackConfigTypeFreeform:
		return true
	}
	return false
}

// Specific value and label pair for feedback
type FeedbackIngestTokenCreateSchemaFeedbackConfigCategoryParam struct {
	Value param.Field[float64] `json:"value,required"`
	Label param.Field[string]  `json:"label"`
}

func (r FeedbackIngestTokenCreateSchemaFeedbackConfigCategoryParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Feedback ingest token schema.
type FeedbackIngestTokenSchema struct {
	ID          string                        `json:"id,required" format:"uuid"`
	ExpiresAt   time.Time                     `json:"expires_at,required" format:"date-time"`
	FeedbackKey string                        `json:"feedback_key,required"`
	URL         string                        `json:"url,required"`
	JSON        feedbackIngestTokenSchemaJSON `json:"-"`
}

// feedbackIngestTokenSchemaJSON contains the JSON metadata for the struct
// [FeedbackIngestTokenSchema]
type feedbackIngestTokenSchemaJSON struct {
	ID          apijson.Field
	ExpiresAt   apijson.Field
	FeedbackKey apijson.Field
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FeedbackIngestTokenSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r feedbackIngestTokenSchemaJSON) RawJSON() string {
	return r.raw
}

func (r FeedbackIngestTokenSchema) implementsFeedbackTokenNewResponseUnion() {}

// Feedback ingest token schema.
//
// Union satisfied by [FeedbackIngestTokenSchema] or
// [FeedbackTokenNewResponseArray].
type FeedbackTokenNewResponseUnion interface {
	implementsFeedbackTokenNewResponseUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*FeedbackTokenNewResponseUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(FeedbackIngestTokenSchema{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(FeedbackTokenNewResponseArray{}),
		},
	)
}

type FeedbackTokenNewResponseArray []FeedbackIngestTokenSchema

func (r FeedbackTokenNewResponseArray) implementsFeedbackTokenNewResponseUnion() {}

type FeedbackTokenGetResponse = interface{}

type FeedbackTokenUpdateResponse = interface{}

type FeedbackTokenNewParams struct {
	// Feedback ingest token create schema.
	Body FeedbackTokenNewParamsBodyUnion `json:"body,required"`
}

func (r FeedbackTokenNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.Body)
}

// Feedback ingest token create schema.
//
// Satisfied by [FeedbackIngestTokenCreateSchemaParam],
// [FeedbackTokenNewParamsBodyArray].
type FeedbackTokenNewParamsBodyUnion interface {
	implementsFeedbackTokenNewParamsBodyUnion()
}

type FeedbackTokenNewParamsBodyArray []FeedbackIngestTokenCreateSchemaParam

func (r FeedbackTokenNewParamsBodyArray) implementsFeedbackTokenNewParamsBodyUnion() {}

type FeedbackTokenGetParams struct {
	Comment    param.Field[string]                           `query:"comment"`
	Correction param.Field[string]                           `query:"correction"`
	Score      param.Field[FeedbackTokenGetParamsScoreUnion] `query:"score"`
	Value      param.Field[FeedbackTokenGetParamsValueUnion] `query:"value"`
}

// URLQuery serializes [FeedbackTokenGetParams]'s query parameters as `url.Values`.
func (r FeedbackTokenGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool].
type FeedbackTokenGetParamsScoreUnion interface {
	ImplementsFeedbackTokenGetParamsScoreUnion()
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool], [shared.UnionString].
type FeedbackTokenGetParamsValueUnion interface {
	ImplementsFeedbackTokenGetParamsValueUnion()
}

type FeedbackTokenUpdateParams struct {
	Comment    param.Field[string]                                   `json:"comment"`
	Correction param.Field[FeedbackTokenUpdateParamsCorrectionUnion] `json:"correction"`
	Metadata   param.Field[map[string]interface{}]                   `json:"metadata"`
	Score      param.Field[FeedbackTokenUpdateParamsScoreUnion]      `json:"score"`
	Value      param.Field[FeedbackTokenUpdateParamsValueUnion]      `json:"value"`
}

func (r FeedbackTokenUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [FeedbackTokenUpdateParamsCorrectionMap], [shared.UnionString].
type FeedbackTokenUpdateParamsCorrectionUnion interface {
	ImplementsFeedbackTokenUpdateParamsCorrectionUnion()
}

type FeedbackTokenUpdateParamsCorrectionMap map[string]interface{}

func (r FeedbackTokenUpdateParamsCorrectionMap) ImplementsFeedbackTokenUpdateParamsCorrectionUnion() {
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool].
type FeedbackTokenUpdateParamsScoreUnion interface {
	ImplementsFeedbackTokenUpdateParamsScoreUnion()
}

// Satisfied by [shared.UnionFloat], [shared.UnionBool], [shared.UnionString].
type FeedbackTokenUpdateParamsValueUnion interface {
	ImplementsFeedbackTokenUpdateParamsValueUnion()
}

type FeedbackTokenListParams struct {
	RunID param.Field[string] `query:"run_id,required" format:"uuid"`
}

// URLQuery serializes [FeedbackTokenListParams]'s query parameters as
// `url.Values`.
func (r FeedbackTokenListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
