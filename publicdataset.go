// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// PublicDatasetService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewPublicDatasetService] method instead.
type PublicDatasetService struct {
	Options []option.RequestOption
}

// NewPublicDatasetService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewPublicDatasetService(opts ...option.RequestOption) (r *PublicDatasetService) {
	r = &PublicDatasetService{}
	r.Options = opts
	return
}

// Get dataset by ids or the shared dataset if not specifed.
func (r *PublicDatasetService) List(ctx context.Context, shareToken string, query PublicDatasetListParams, opts ...option.RequestOption) (res *PublicDatasetListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/public/%s/datasets", shareToken)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Get all comparative experiments for a given dataset.
func (r *PublicDatasetService) ListComparative(ctx context.Context, shareToken string, query PublicDatasetListComparativeParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[PublicDatasetListComparativeResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/public/%s/datasets/comparative", shareToken)
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

// Get all comparative experiments for a given dataset.
func (r *PublicDatasetService) ListComparativeAutoPaging(ctx context.Context, shareToken string, query PublicDatasetListComparativeParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[PublicDatasetListComparativeResponse] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.ListComparative(ctx, shareToken, query, opts...))
}

// Get feedback for runs in projects run over a dataset that has been shared.
func (r *PublicDatasetService) ListFeedback(ctx context.Context, shareToken string, query PublicDatasetListFeedbackParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[FeedbackSchema], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/public/%s/datasets/feedback", shareToken)
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

// Get feedback for runs in projects run over a dataset that has been shared.
func (r *PublicDatasetService) ListFeedbackAutoPaging(ctx context.Context, shareToken string, query PublicDatasetListFeedbackParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[FeedbackSchema] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.ListFeedback(ctx, shareToken, query, opts...))
}

// Get projects run on a dataset that has been shared.
func (r *PublicDatasetService) ListSessions(ctx context.Context, shareToken string, params PublicDatasetListSessionsParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[TracerSession], err error) {
	var raw *http.Response
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("accept", fmt.Sprintf("%s", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/public/%s/datasets/sessions", shareToken)
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

// Get projects run on a dataset that has been shared.
func (r *PublicDatasetService) ListSessionsAutoPaging(ctx context.Context, shareToken string, params PublicDatasetListSessionsParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[TracerSession] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.ListSessions(ctx, shareToken, params, opts...))
}

// Get sessions from multiple datasets using share tokens.
func (r *PublicDatasetService) GetSessionsBulk(ctx context.Context, query PublicDatasetGetSessionsBulkParams, opts ...option.RequestOption) (res *[]TracerSession, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/public/datasets/sessions-bulk"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Public schema for datasets.
//
// Doesn't currently include session counts/stats since public test project sharing
// is not yet shipped
type PublicDatasetListResponse struct {
	ID           string    `json:"id,required" format:"uuid"`
	ExampleCount int64     `json:"example_count,required"`
	Name         string    `json:"name,required"`
	CreatedAt    time.Time `json:"created_at" format:"date-time"`
	// Enum for dataset data types.
	DataType                DataType                      `json:"data_type,nullable"`
	Description             string                        `json:"description,nullable"`
	ExternallyManaged       bool                          `json:"externally_managed,nullable"`
	InputsSchemaDefinition  interface{}                   `json:"inputs_schema_definition,nullable"`
	OutputsSchemaDefinition interface{}                   `json:"outputs_schema_definition,nullable"`
	Transformations         []DatasetTransformation       `json:"transformations,nullable"`
	JSON                    publicDatasetListResponseJSON `json:"-"`
}

// publicDatasetListResponseJSON contains the JSON metadata for the struct
// [PublicDatasetListResponse]
type publicDatasetListResponseJSON struct {
	ID                      apijson.Field
	ExampleCount            apijson.Field
	Name                    apijson.Field
	CreatedAt               apijson.Field
	DataType                apijson.Field
	Description             apijson.Field
	ExternallyManaged       apijson.Field
	InputsSchemaDefinition  apijson.Field
	OutputsSchemaDefinition apijson.Field
	Transformations         apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *PublicDatasetListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r publicDatasetListResponseJSON) RawJSON() string {
	return r.raw
}

// Publicly-shared ComparativeExperiment schema.
type PublicDatasetListComparativeResponse struct {
	ID              string                                   `json:"id,required" format:"uuid"`
	CreatedAt       time.Time                                `json:"created_at,required" format:"date-time"`
	ExperimentsInfo []SimpleExperimentInfo                   `json:"experiments_info,required"`
	ModifiedAt      time.Time                                `json:"modified_at,required" format:"date-time"`
	Description     string                                   `json:"description,nullable"`
	Extra           interface{}                              `json:"extra,nullable"`
	FeedbackStats   interface{}                              `json:"feedback_stats,nullable"`
	Name            string                                   `json:"name,nullable"`
	JSON            publicDatasetListComparativeResponseJSON `json:"-"`
}

// publicDatasetListComparativeResponseJSON contains the JSON metadata for the
// struct [PublicDatasetListComparativeResponse]
type publicDatasetListComparativeResponseJSON struct {
	ID              apijson.Field
	CreatedAt       apijson.Field
	ExperimentsInfo apijson.Field
	ModifiedAt      apijson.Field
	Description     apijson.Field
	Extra           apijson.Field
	FeedbackStats   apijson.Field
	Name            apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *PublicDatasetListComparativeResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r publicDatasetListComparativeResponseJSON) RawJSON() string {
	return r.raw
}

type PublicDatasetListParams struct {
	Limit  param.Field[int64] `query:"limit"`
	Offset param.Field[int64] `query:"offset"`
	// Enum for available dataset columns to sort by.
	SortBy     param.Field[SortByDatasetColumn] `query:"sort_by"`
	SortByDesc param.Field[bool]                `query:"sort_by_desc"`
}

// URLQuery serializes [PublicDatasetListParams]'s query parameters as
// `url.Values`.
func (r PublicDatasetListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PublicDatasetListComparativeParams struct {
	Limit        param.Field[int64]  `query:"limit"`
	Name         param.Field[string] `query:"name"`
	NameContains param.Field[string] `query:"name_contains"`
	Offset       param.Field[int64]  `query:"offset"`
	// Enum for available comparative experiment columns to sort by.
	SortBy     param.Field[SortByComparativeExperimentColumn] `query:"sort_by"`
	SortByDesc param.Field[bool]                              `query:"sort_by_desc"`
}

// URLQuery serializes [PublicDatasetListComparativeParams]'s query parameters as
// `url.Values`.
func (r PublicDatasetListComparativeParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PublicDatasetListFeedbackParams struct {
	HasComment param.Field[bool]     `query:"has_comment"`
	HasScore   param.Field[bool]     `query:"has_score"`
	Key        param.Field[[]string] `query:"key"`
	// Enum for feedback levels.
	Level   param.Field[FeedbackLevel] `query:"level"`
	Limit   param.Field[int64]         `query:"limit"`
	Offset  param.Field[int64]         `query:"offset"`
	Run     param.Field[[]string]      `query:"run" format:"uuid"`
	Session param.Field[[]string]      `query:"session" format:"uuid"`
	Source  param.Field[[]SourceType]  `query:"source"`
	User    param.Field[[]string]      `query:"user" format:"uuid"`
}

// URLQuery serializes [PublicDatasetListFeedbackParams]'s query parameters as
// `url.Values`.
func (r PublicDatasetListFeedbackParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PublicDatasetListSessionsParams struct {
	ID                param.Field[[]string]               `query:"id" format:"uuid"`
	DatasetVersion    param.Field[string]                 `query:"dataset_version"`
	Facets            param.Field[bool]                   `query:"facets"`
	Limit             param.Field[int64]                  `query:"limit"`
	Name              param.Field[string]                 `query:"name"`
	NameContains      param.Field[string]                 `query:"name_contains"`
	Offset            param.Field[int64]                  `query:"offset"`
	SortBy            param.Field[SessionSortableColumns] `query:"sort_by"`
	SortByDesc        param.Field[bool]                   `query:"sort_by_desc"`
	SortByFeedbackKey param.Field[string]                 `query:"sort_by_feedback_key"`
	Accept            param.Field[string]                 `header:"accept"`
}

// URLQuery serializes [PublicDatasetListSessionsParams]'s query parameters as
// `url.Values`.
func (r PublicDatasetListSessionsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PublicDatasetGetSessionsBulkParams struct {
	ShareTokens param.Field[[]string] `query:"share_tokens,required"`
}

// URLQuery serializes [PublicDatasetGetSessionsBulkParams]'s query parameters as
// `url.Values`.
func (r PublicDatasetGetSessionsBulkParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
