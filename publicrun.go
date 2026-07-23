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
)

// PublicRunService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewPublicRunService] method instead.
type PublicRunService struct {
	Options []option.RequestOption
}

// NewPublicRunService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewPublicRunService(opts ...option.RequestOption) (r *PublicRunService) {
	r = &PublicRunService{}
	r.Options = opts
	return
}

// **Alpha:** The request and response contract may change; Returns one run within
// the trace identified by the share token. The request supplies only the run ID
// and that run's exact start_time coordinate.
func (r *PublicRunService) Get(ctx context.Context, shareToken string, runID string, params PublicRunGetParams, opts ...option.RequestOption) (res *Run, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return nil, err
	}
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/public/%s/run/%s", shareToken, runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, params, &res, opts...)
	return res, err
}

// **Alpha:** The request and response contract may change; Returns all runs within
// the trace identified by the share token. The share token supplies the tenant,
// project, and trace scope.
func (r *PublicRunService) Query(ctx context.Context, shareToken string, params PublicRunQueryParams, opts ...option.RequestOption) (res *PublicRunQueryResponse, err error) {
	if params.Accept.Present {
		opts = append(opts, option.WithHeader("Accept", fmt.Sprintf("%v", params.Accept)))
	}
	opts = slices.Concat(r.Options, opts)
	if shareToken == "" {
		err = errors.New("missing required share_token parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/public/%s/runs/v2/query", shareToken)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

type PublicRunQueryResponse struct {
	// `items` lists runs in the trace for the requested time window, in `start_time`
	// order.
	Items []Run                      `json:"items"`
	JSON  publicRunQueryResponseJSON `json:"-"`
}

// publicRunQueryResponseJSON contains the JSON metadata for the struct
// [PublicRunQueryResponse]
type publicRunQueryResponseJSON struct {
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *PublicRunQueryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r publicRunQueryResponseJSON) RawJSON() string {
	return r.raw
}

type PublicRunGetParams struct {
	// repeatable public run fields to include
	Selects param.Field[[]string] `query:"selects" api:"required"`
	// Run start_time coordinate (RFC3339)
	StartTime param.Field[time.Time] `query:"start_time" api:"required" format:"date-time"`
	Accept    param.Field[string]    `header:"Accept"`
}

// URLQuery serializes [PublicRunGetParams]'s query parameters as `url.Values`.
func (r PublicRunGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type PublicRunQueryParams struct {
	// `selects` lists which public run properties to include on each returned run.
	Selects param.Field[[]PublicRunQueryParamsSelect] `json:"selects"`
	Accept  param.Field[string]                       `header:"Accept"`
}

func (r PublicRunQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type PublicRunQueryParamsSelect string

const (
	PublicRunQueryParamsSelectID                     PublicRunQueryParamsSelect = "ID"
	PublicRunQueryParamsSelectName                   PublicRunQueryParamsSelect = "NAME"
	PublicRunQueryParamsSelectRunType                PublicRunQueryParamsSelect = "RUN_TYPE"
	PublicRunQueryParamsSelectStatus                 PublicRunQueryParamsSelect = "STATUS"
	PublicRunQueryParamsSelectStartTime              PublicRunQueryParamsSelect = "START_TIME"
	PublicRunQueryParamsSelectEndTime                PublicRunQueryParamsSelect = "END_TIME"
	PublicRunQueryParamsSelectLatencySeconds         PublicRunQueryParamsSelect = "LATENCY_SECONDS"
	PublicRunQueryParamsSelectFirstTokenTime         PublicRunQueryParamsSelect = "FIRST_TOKEN_TIME"
	PublicRunQueryParamsSelectError                  PublicRunQueryParamsSelect = "ERROR"
	PublicRunQueryParamsSelectErrorPreview           PublicRunQueryParamsSelect = "ERROR_PREVIEW"
	PublicRunQueryParamsSelectExtra                  PublicRunQueryParamsSelect = "EXTRA"
	PublicRunQueryParamsSelectMetadata               PublicRunQueryParamsSelect = "METADATA"
	PublicRunQueryParamsSelectInputsPreview          PublicRunQueryParamsSelect = "INPUTS_PREVIEW"
	PublicRunQueryParamsSelectOutputsPreview         PublicRunQueryParamsSelect = "OUTPUTS_PREVIEW"
	PublicRunQueryParamsSelectParentRunID            PublicRunQueryParamsSelect = "PARENT_RUN_ID"
	PublicRunQueryParamsSelectParentRunIDs           PublicRunQueryParamsSelect = "PARENT_RUN_IDS"
	PublicRunQueryParamsSelectProjectID              PublicRunQueryParamsSelect = "PROJECT_ID"
	PublicRunQueryParamsSelectTraceID                PublicRunQueryParamsSelect = "TRACE_ID"
	PublicRunQueryParamsSelectThreadID               PublicRunQueryParamsSelect = "THREAD_ID"
	PublicRunQueryParamsSelectDottedOrder            PublicRunQueryParamsSelect = "DOTTED_ORDER"
	PublicRunQueryParamsSelectIsRoot                 PublicRunQueryParamsSelect = "IS_ROOT"
	PublicRunQueryParamsSelectReferenceDatasetID     PublicRunQueryParamsSelect = "REFERENCE_DATASET_ID"
	PublicRunQueryParamsSelectTotalTokens            PublicRunQueryParamsSelect = "TOTAL_TOKENS"
	PublicRunQueryParamsSelectPromptTokens           PublicRunQueryParamsSelect = "PROMPT_TOKENS"
	PublicRunQueryParamsSelectCompletionTokens       PublicRunQueryParamsSelect = "COMPLETION_TOKENS"
	PublicRunQueryParamsSelectTotalCost              PublicRunQueryParamsSelect = "TOTAL_COST"
	PublicRunQueryParamsSelectPromptCost             PublicRunQueryParamsSelect = "PROMPT_COST"
	PublicRunQueryParamsSelectCompletionCost         PublicRunQueryParamsSelect = "COMPLETION_COST"
	PublicRunQueryParamsSelectPromptTokenDetails     PublicRunQueryParamsSelect = "PROMPT_TOKEN_DETAILS"
	PublicRunQueryParamsSelectCompletionTokenDetails PublicRunQueryParamsSelect = "COMPLETION_TOKEN_DETAILS"
	PublicRunQueryParamsSelectPromptCostDetails      PublicRunQueryParamsSelect = "PROMPT_COST_DETAILS"
	PublicRunQueryParamsSelectCompletionCostDetails  PublicRunQueryParamsSelect = "COMPLETION_COST_DETAILS"
	PublicRunQueryParamsSelectPriceModelID           PublicRunQueryParamsSelect = "PRICE_MODEL_ID"
	PublicRunQueryParamsSelectTags                   PublicRunQueryParamsSelect = "TAGS"
	PublicRunQueryParamsSelectThreadEvaluationTime   PublicRunQueryParamsSelect = "THREAD_EVALUATION_TIME"
	PublicRunQueryParamsSelectFeedbackStats          PublicRunQueryParamsSelect = "FEEDBACK_STATS"
)

func (r PublicRunQueryParamsSelect) IsKnown() bool {
	switch r {
	case PublicRunQueryParamsSelectID, PublicRunQueryParamsSelectName, PublicRunQueryParamsSelectRunType, PublicRunQueryParamsSelectStatus, PublicRunQueryParamsSelectStartTime, PublicRunQueryParamsSelectEndTime, PublicRunQueryParamsSelectLatencySeconds, PublicRunQueryParamsSelectFirstTokenTime, PublicRunQueryParamsSelectError, PublicRunQueryParamsSelectErrorPreview, PublicRunQueryParamsSelectExtra, PublicRunQueryParamsSelectMetadata, PublicRunQueryParamsSelectInputsPreview, PublicRunQueryParamsSelectOutputsPreview, PublicRunQueryParamsSelectParentRunID, PublicRunQueryParamsSelectParentRunIDs, PublicRunQueryParamsSelectProjectID, PublicRunQueryParamsSelectTraceID, PublicRunQueryParamsSelectThreadID, PublicRunQueryParamsSelectDottedOrder, PublicRunQueryParamsSelectIsRoot, PublicRunQueryParamsSelectReferenceDatasetID, PublicRunQueryParamsSelectTotalTokens, PublicRunQueryParamsSelectPromptTokens, PublicRunQueryParamsSelectCompletionTokens, PublicRunQueryParamsSelectTotalCost, PublicRunQueryParamsSelectPromptCost, PublicRunQueryParamsSelectCompletionCost, PublicRunQueryParamsSelectPromptTokenDetails, PublicRunQueryParamsSelectCompletionTokenDetails, PublicRunQueryParamsSelectPromptCostDetails, PublicRunQueryParamsSelectCompletionCostDetails, PublicRunQueryParamsSelectPriceModelID, PublicRunQueryParamsSelectTags, PublicRunQueryParamsSelectThreadEvaluationTime, PublicRunQueryParamsSelectFeedbackStats:
		return true
	}
	return false
}
