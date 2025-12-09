// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// PublicService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewPublicService] method instead.
type PublicService struct {
	Options  []option.RequestOption
	Datasets *PublicDatasetService
}

// NewPublicService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewPublicService(opts ...option.RequestOption) (r *PublicService) {
	r = &PublicService{}
	r.Options = opts
	r.Datasets = NewPublicDatasetService(opts...)
	return
}

// Read Shared Feedbacks
func (r *PublicService) GetFeedbacks(ctx context.Context, shareToken string, query PublicGetFeedbacksParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[FeedbackSchema], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return
	}
	path := fmt.Sprintf("api/v1/public/%s/feedbacks", shareToken)
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

// Read Shared Feedbacks
func (r *PublicService) GetFeedbacksAutoPaging(ctx context.Context, shareToken string, query PublicGetFeedbacksParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[FeedbackSchema] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.GetFeedbacks(ctx, shareToken, query, opts...))
}

type PublicGetFeedbacksParams struct {
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

// URLQuery serializes [PublicGetFeedbacksParams]'s query parameters as
// `url.Values`.
func (r PublicGetFeedbacksParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
