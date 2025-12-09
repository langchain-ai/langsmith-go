// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// FeedbackConfigService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFeedbackConfigService] method instead.
type FeedbackConfigService struct {
	Options []option.RequestOption
}

// NewFeedbackConfigService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFeedbackConfigService(opts ...option.RequestOption) (r *FeedbackConfigService) {
	r = &FeedbackConfigService{}
	r.Options = opts
	return
}

// Soft delete a feedback config by marking it as deleted.
//
// The config can be recreated later with the same key (simple reuse pattern).
// Existing feedback records with this key will remain unchanged.
func (r *FeedbackConfigService) Delete(ctx context.Context, body FeedbackConfigDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	path := "api/v1/feedback-configs"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return
}

type FeedbackConfigDeleteParams struct {
	FeedbackKey param.Field[string] `query:"feedback_key,required"`
}

// URLQuery serializes [FeedbackConfigDeleteParams]'s query parameters as
// `url.Values`.
func (r FeedbackConfigDeleteParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
