// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// OnlineEvaluatorService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewOnlineEvaluatorService] method instead.
type OnlineEvaluatorService struct {
	Options []option.RequestOption
}

// NewOnlineEvaluatorService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewOnlineEvaluatorService(opts ...option.RequestOption) (r *OnlineEvaluatorService) {
	r = &OnlineEvaluatorService{}
	r.Options = opts
	return
}

// Create a new LLM or code evaluator for the current workspace.
func (r *OnlineEvaluatorService) New(ctx context.Context, body OnlineEvaluatorNewParams, opts ...option.RequestOption) (res *CreateOnlineEvaluatorResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/platform/evaluators"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a single evaluator by its ID.
func (r *OnlineEvaluatorService) Get(ctx context.Context, evaluatorID string, opts ...option.RequestOption) (res *OnlineEvaluator, err error) {
	opts = slices.Concat(r.Options, opts)
	if evaluatorID == "" {
		err = errors.New("missing required evaluator_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/platform/evaluators/%s", evaluatorID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update an existing evaluator's name, LLM configuration, or code configuration.
func (r *OnlineEvaluatorService) Update(ctx context.Context, evaluatorID string, body OnlineEvaluatorUpdateParams, opts ...option.RequestOption) (res *UpdateOnlineEvaluatorResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if evaluatorID == "" {
		err = errors.New("missing required evaluator_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/platform/evaluators/%s", evaluatorID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List evaluators for the current workspace, with optional filtering by type,
// name, tag, feedback key, or resource ID.
func (r *OnlineEvaluatorService) List(ctx context.Context, query OnlineEvaluatorListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationOnlineEvaluators[OnlineEvaluator], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/platform/evaluators"
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

// List evaluators for the current workspace, with optional filtering by type,
// name, tag, feedback key, or resource ID.
func (r *OnlineEvaluatorService) ListAutoPaging(ctx context.Context, query OnlineEvaluatorListParams, opts ...option.RequestOption) *pagination.OffsetPaginationOnlineEvaluatorsAutoPager[OnlineEvaluator] {
	return pagination.NewOffsetPaginationOnlineEvaluatorsAutoPager(r.List(ctx, query, opts...))
}

// Delete an evaluator. When delete_run_rules is true, all run rules referencing
// this evaluator are deleted first (same tenant). Associated llm_evaluators and
// code_evaluators rows are removed by foreign-key cascade when the evaluator row
// is deleted.
func (r *OnlineEvaluatorService) Delete(ctx context.Context, evaluatorID string, body OnlineEvaluatorDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if evaluatorID == "" {
		err = errors.New("missing required evaluator_id parameter")
		return err
	}
	path := fmt.Sprintf("v1/platform/evaluators/%s", evaluatorID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return err
}

// Delete multiple evaluators by their IDs. Returns per-item success/failure.
func (r *OnlineEvaluatorService) BulkDelete(ctx context.Context, body OnlineEvaluatorBulkDeleteParams, opts ...option.RequestOption) (res *BulkDeleteEvaluatorsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/platform/evaluators"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, &res, opts...)
	return res, err
}

// Returns per-day LLM evaluator spend for the requested 7-day period, grouped by
// evaluator, resource, or run rule. Exactly one of group_by, evaluator_id,
// session_id, or dataset_id is required. resource_id, type, and feedback_key may
// be supplied with group_by to narrow listing aggregations.
func (r *OnlineEvaluatorService) Spend(ctx context.Context, query OnlineEvaluatorSpendParams, opts ...option.RequestOption) (res *GetOnlineEvaluatorSpendResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v1/platform/evaluators/spend"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type BulkDeleteEvaluatorFailedItem struct {
	ID    string                            `json:"id"`
	Error string                            `json:"error"`
	JSON  bulkDeleteEvaluatorFailedItemJSON `json:"-"`
}

// bulkDeleteEvaluatorFailedItemJSON contains the JSON metadata for the struct
// [BulkDeleteEvaluatorFailedItem]
type bulkDeleteEvaluatorFailedItemJSON struct {
	ID          apijson.Field
	Error       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BulkDeleteEvaluatorFailedItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r bulkDeleteEvaluatorFailedItemJSON) RawJSON() string {
	return r.raw
}

type BulkDeleteEvaluatorsResponse struct {
	Failed    []BulkDeleteEvaluatorFailedItem  `json:"failed"`
	Succeeded []string                         `json:"succeeded"`
	JSON      bulkDeleteEvaluatorsResponseJSON `json:"-"`
}

// bulkDeleteEvaluatorsResponseJSON contains the JSON metadata for the struct
// [BulkDeleteEvaluatorsResponse]
type bulkDeleteEvaluatorsResponseJSON struct {
	Failed      apijson.Field
	Succeeded   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BulkDeleteEvaluatorsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r bulkDeleteEvaluatorsResponseJSON) RawJSON() string {
	return r.raw
}

type CreateOnlineCodeEvaluatorRequestParam struct {
	Code param.Field[string] `json:"code"`
	// Default: "python"
	Language param.Field[string] `json:"language"`
}

func (r CreateOnlineCodeEvaluatorRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type CreateOnlineEvaluatorRequestParam struct {
	CodeEvaluator param.Field[CreateOnlineCodeEvaluatorRequestParam] `json:"code_evaluator"`
	LlmEvaluator  param.Field[CreateOnlineLlmEvaluatorRequestParam]  `json:"llm_evaluator"`
	Name          param.Field[string]                                `json:"name"`
	Type          param.Field[OnlineEvaluatorType]                   `json:"type"`
}

func (r CreateOnlineEvaluatorRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type CreateOnlineEvaluatorResponse struct {
	Evaluator OnlineEvaluator                   `json:"evaluator"`
	JSON      createOnlineEvaluatorResponseJSON `json:"-"`
}

// createOnlineEvaluatorResponseJSON contains the JSON metadata for the struct
// [CreateOnlineEvaluatorResponse]
type createOnlineEvaluatorResponseJSON struct {
	Evaluator   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CreateOnlineEvaluatorResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r createOnlineEvaluatorResponseJSON) RawJSON() string {
	return r.raw
}

type CreateOnlineLlmEvaluatorRequestParam struct {
	CommitHashOrTag  param.Field[string]      `json:"commit_hash_or_tag"`
	PromptRepoHandle param.Field[string]      `json:"prompt_repo_handle"`
	VariableMapping  param.Field[interface{}] `json:"variable_mapping"`
}

func (r CreateOnlineLlmEvaluatorRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type GetOnlineEvaluatorSpendResponse struct {
	Groups      []OnlineEvaluatorSpendGroup         `json:"groups"`
	PeriodEnd   string                              `json:"period_end"`
	PeriodStart string                              `json:"period_start"`
	JSON        getOnlineEvaluatorSpendResponseJSON `json:"-"`
}

// getOnlineEvaluatorSpendResponseJSON contains the JSON metadata for the struct
// [GetOnlineEvaluatorSpendResponse]
type getOnlineEvaluatorSpendResponseJSON struct {
	Groups      apijson.Field
	PeriodEnd   apijson.Field
	PeriodStart apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *GetOnlineEvaluatorSpendResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r getOnlineEvaluatorSpendResponseJSON) RawJSON() string {
	return r.raw
}

type OnlineCodeEvaluator struct {
	Code        string `json:"code"`
	EvaluatorID string `json:"evaluator_id"`
	// Default: "python"
	Language string                  `json:"language"`
	JSON     onlineCodeEvaluatorJSON `json:"-"`
}

// onlineCodeEvaluatorJSON contains the JSON metadata for the struct
// [OnlineCodeEvaluator]
type onlineCodeEvaluatorJSON struct {
	Code        apijson.Field
	EvaluatorID apijson.Field
	Language    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OnlineCodeEvaluator) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineCodeEvaluatorJSON) RawJSON() string {
	return r.raw
}

type OnlineEvaluator struct {
	ID            string              `json:"id"`
	CodeEvaluator OnlineCodeEvaluator `json:"code_evaluator"`
	CreatedAt     string              `json:"created_at"`
	CreatedBy     string              `json:"created_by"`
	FeedbackKeys  []string            `json:"feedback_keys"`
	// Embedded child evaluator (populated based on type)
	LlmEvaluator OnlineLlmEvaluator       `json:"llm_evaluator"`
	Name         string                   `json:"name"`
	RunRules     []OnlineEvaluatorRunRule `json:"run_rules"`
	TenantID     string                   `json:"tenant_id"`
	Type         OnlineEvaluatorType      `json:"type"`
	UpdatedAt    string                   `json:"updated_at"`
	JSON         onlineEvaluatorJSON      `json:"-"`
}

// onlineEvaluatorJSON contains the JSON metadata for the struct [OnlineEvaluator]
type onlineEvaluatorJSON struct {
	ID            apijson.Field
	CodeEvaluator apijson.Field
	CreatedAt     apijson.Field
	CreatedBy     apijson.Field
	FeedbackKeys  apijson.Field
	LlmEvaluator  apijson.Field
	Name          apijson.Field
	RunRules      apijson.Field
	TenantID      apijson.Field
	Type          apijson.Field
	UpdatedAt     apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *OnlineEvaluator) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineEvaluatorJSON) RawJSON() string {
	return r.raw
}

type OnlineEvaluatorRunRule struct {
	ID                   string `json:"id"`
	CorrectionsDatasetID string `json:"corrections_dataset_id"`
	DatasetID            string `json:"dataset_id"`
	DatasetName          string `json:"dataset_name"`
	GroupBy              string `json:"group_by"`
	NumFewShotExamples   int64  `json:"num_few_shot_examples"`
	SessionID            string `json:"session_id"`
	SessionName          string `json:"session_name"`
	// SpendLimit is the effective spend-cap limit for this rule (nil when
	// unconfigured).
	SpendLimit OnlineSpendLimit `json:"spend_limit"`
	// Per-rule spend for the current ISO week (omitted when feature is disabled).
	// LLM-evaluator rules are initialized to 0; code-evaluator rules remain nil.
	SpendUsd              float64                    `json:"spend_usd"`
	TraceCount            int64                      `json:"trace_count"`
	UseCorrectionsDataset bool                       `json:"use_corrections_dataset"`
	JSON                  onlineEvaluatorRunRuleJSON `json:"-"`
}

// onlineEvaluatorRunRuleJSON contains the JSON metadata for the struct
// [OnlineEvaluatorRunRule]
type onlineEvaluatorRunRuleJSON struct {
	ID                    apijson.Field
	CorrectionsDatasetID  apijson.Field
	DatasetID             apijson.Field
	DatasetName           apijson.Field
	GroupBy               apijson.Field
	NumFewShotExamples    apijson.Field
	SessionID             apijson.Field
	SessionName           apijson.Field
	SpendLimit            apijson.Field
	SpendUsd              apijson.Field
	TraceCount            apijson.Field
	UseCorrectionsDataset apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *OnlineEvaluatorRunRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineEvaluatorRunRuleJSON) RawJSON() string {
	return r.raw
}

type OnlineEvaluatorSpendDay struct {
	Date       string                      `json:"date"`
	SpendUsd   float64                     `json:"spend_usd"`
	TraceCount int64                       `json:"trace_count"`
	JSON       onlineEvaluatorSpendDayJSON `json:"-"`
}

// onlineEvaluatorSpendDayJSON contains the JSON metadata for the struct
// [OnlineEvaluatorSpendDay]
type onlineEvaluatorSpendDayJSON struct {
	Date        apijson.Field
	SpendUsd    apijson.Field
	TraceCount  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *OnlineEvaluatorSpendDay) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineEvaluatorSpendDayJSON) RawJSON() string {
	return r.raw
}

type OnlineEvaluatorSpendGroup struct {
	DatasetID           string                        `json:"dataset_id"`
	DatasetName         string                        `json:"dataset_name"`
	Days                []OnlineEvaluatorSpendDay     `json:"days"`
	EvaluatorID         string                        `json:"evaluator_id"`
	EvaluatorName       string                        `json:"evaluator_name"`
	PrevTotalSpendUsd   float64                       `json:"prev_total_spend_usd"`
	PrevTotalTraceCount int64                         `json:"prev_total_trace_count"`
	RunRuleID           string                        `json:"run_rule_id"`
	RunRuleName         string                        `json:"run_rule_name"`
	SessionID           string                        `json:"session_id"`
	SessionName         string                        `json:"session_name"`
	SpendLimit          OnlineSpendLimit              `json:"spend_limit"`
	TotalSpendUsd       float64                       `json:"total_spend_usd"`
	TotalTraceCount     int64                         `json:"total_trace_count"`
	JSON                onlineEvaluatorSpendGroupJSON `json:"-"`
}

// onlineEvaluatorSpendGroupJSON contains the JSON metadata for the struct
// [OnlineEvaluatorSpendGroup]
type onlineEvaluatorSpendGroupJSON struct {
	DatasetID           apijson.Field
	DatasetName         apijson.Field
	Days                apijson.Field
	EvaluatorID         apijson.Field
	EvaluatorName       apijson.Field
	PrevTotalSpendUsd   apijson.Field
	PrevTotalTraceCount apijson.Field
	RunRuleID           apijson.Field
	RunRuleName         apijson.Field
	SessionID           apijson.Field
	SessionName         apijson.Field
	SpendLimit          apijson.Field
	TotalSpendUsd       apijson.Field
	TotalTraceCount     apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *OnlineEvaluatorSpendGroup) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineEvaluatorSpendGroupJSON) RawJSON() string {
	return r.raw
}

type OnlineEvaluatorType string

const (
	OnlineEvaluatorTypeLlm  OnlineEvaluatorType = "llm"
	OnlineEvaluatorTypeCode OnlineEvaluatorType = "code"
)

func (r OnlineEvaluatorType) IsKnown() bool {
	switch r {
	case OnlineEvaluatorTypeLlm, OnlineEvaluatorTypeCode:
		return true
	}
	return false
}

type OnlineLlmEvaluator struct {
	AnnotationQueueID    string `json:"annotation_queue_id"`
	CommitHashOrTag      string `json:"commit_hash_or_tag"`
	CorrectionsDatasetID string `json:"corrections_dataset_id"`
	EvaluatorID          string `json:"evaluator_id"`
	NumFewShotExamples   int64  `json:"num_few_shot_examples"`
	PromptID             string `json:"prompt_id"`
	PromptRepoHandle     string `json:"prompt_repo_handle"`
	// Derived from the evaluator's run rules — shared across all rules on this
	// evaluator. Nil when the evaluator has no run rules.
	UseCorrectionsDataset bool `json:"use_corrections_dataset"`
	// JSONB
	VariableMapping interface{}            `json:"variable_mapping"`
	JSON            onlineLlmEvaluatorJSON `json:"-"`
}

// onlineLlmEvaluatorJSON contains the JSON metadata for the struct
// [OnlineLlmEvaluator]
type onlineLlmEvaluatorJSON struct {
	AnnotationQueueID     apijson.Field
	CommitHashOrTag       apijson.Field
	CorrectionsDatasetID  apijson.Field
	EvaluatorID           apijson.Field
	NumFewShotExamples    apijson.Field
	PromptID              apijson.Field
	PromptRepoHandle      apijson.Field
	UseCorrectionsDataset apijson.Field
	VariableMapping       apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *OnlineLlmEvaluator) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineLlmEvaluatorJSON) RawJSON() string {
	return r.raw
}

type OnlineSpendLimit struct {
	LimitUsd       float64              `json:"limit_usd"`
	UtilizationPct float64              `json:"utilization_pct"`
	Window         string               `json:"window"`
	JSON           onlineSpendLimitJSON `json:"-"`
}

// onlineSpendLimitJSON contains the JSON metadata for the struct
// [OnlineSpendLimit]
type onlineSpendLimitJSON struct {
	LimitUsd       apijson.Field
	UtilizationPct apijson.Field
	Window         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *OnlineSpendLimit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r onlineSpendLimitJSON) RawJSON() string {
	return r.raw
}

type UpdateOnlineCodeEvaluatorRequestParam struct {
	Code     param.Field[string] `json:"code"`
	Language param.Field[string] `json:"language"`
}

func (r UpdateOnlineCodeEvaluatorRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type UpdateOnlineEvaluatorRequestParam struct {
	CodeEvaluator param.Field[UpdateOnlineCodeEvaluatorRequestParam] `json:"code_evaluator"`
	LlmEvaluator  param.Field[UpdateOnlineLlmEvaluatorRequestParam]  `json:"llm_evaluator"`
	Name          param.Field[string]                                `json:"name"`
}

func (r UpdateOnlineEvaluatorRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type UpdateOnlineEvaluatorResponse struct {
	Evaluator OnlineEvaluator                   `json:"evaluator"`
	JSON      updateOnlineEvaluatorResponseJSON `json:"-"`
}

// updateOnlineEvaluatorResponseJSON contains the JSON metadata for the struct
// [UpdateOnlineEvaluatorResponse]
type updateOnlineEvaluatorResponseJSON struct {
	Evaluator   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *UpdateOnlineEvaluatorResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r updateOnlineEvaluatorResponseJSON) RawJSON() string {
	return r.raw
}

type UpdateOnlineLlmEvaluatorRequestParam struct {
	CommitHashOrTag       param.Field[string]      `json:"commit_hash_or_tag"`
	NumFewShotExamples    param.Field[int64]       `json:"num_few_shot_examples"`
	PromptRepoHandle      param.Field[string]      `json:"prompt_repo_handle"`
	UseCorrectionsDataset param.Field[bool]        `json:"use_corrections_dataset"`
	VariableMapping       param.Field[interface{}] `json:"variable_mapping"`
}

func (r UpdateOnlineLlmEvaluatorRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type OnlineEvaluatorNewParams struct {
	CreateOnlineEvaluatorRequest CreateOnlineEvaluatorRequestParam `json:"create_online_evaluator_request" api:"required"`
}

func (r OnlineEvaluatorNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.CreateOnlineEvaluatorRequest)
}

type OnlineEvaluatorUpdateParams struct {
	UpdateOnlineEvaluatorRequest UpdateOnlineEvaluatorRequestParam `json:"update_online_evaluator_request" api:"required"`
}

func (r OnlineEvaluatorUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.UpdateOnlineEvaluatorRequest)
}

type OnlineEvaluatorListParams struct {
	// Filter by feedback key
	FeedbackKey param.Field[string] `query:"feedback_key"`
	// Maximum number of results (1-100)
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring (also searches creator names)
	NameContains param.Field[string] `query:"name_contains"`
	// Offset for pagination
	Offset param.Field[int64] `query:"offset"`
	// Filter by resource IDs
	ResourceID param.Field[[]string] `query:"resource_id"`
	// Field to sort by
	SortBy param.Field[string] `query:"sort_by"`
	// Sort in descending order
	SortByDesc param.Field[bool] `query:"sort_by_desc"`
	// Filter by tag value IDs
	TagValueID param.Field[[]string] `query:"tag_value_id"`
	// Filter by evaluator type
	Type param.Field[string] `query:"type"`
}

// URLQuery serializes [OnlineEvaluatorListParams]'s query parameters as
// `url.Values`.
func (r OnlineEvaluatorListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OnlineEvaluatorDeleteParams struct {
	// When true, delete all run rules for this evaluator before deleting the evaluator
	DeleteRunRules param.Field[bool] `query:"delete_run_rules"`
}

// URLQuery serializes [OnlineEvaluatorDeleteParams]'s query parameters as
// `url.Values`.
func (r OnlineEvaluatorDeleteParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OnlineEvaluatorBulkDeleteParams struct {
	// Evaluator IDs to delete
	EvaluatorIDs param.Field[[]string] `query:"evaluator_ids" api:"required"`
	// When true, delete all run rules for this evaluator before deleting the evaluator
	DeleteRunRules param.Field[bool] `query:"delete_run_rules"`
}

// URLQuery serializes [OnlineEvaluatorBulkDeleteParams]'s query parameters as
// `url.Values`.
func (r OnlineEvaluatorBulkDeleteParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type OnlineEvaluatorSpendParams struct {
	// Start of the 7-day window (YYYY-MM-DD).
	PeriodStart param.Field[string] `query:"period_start" api:"required"`
	// Filter to a specific dataset (UUID). Mutually exclusive with group_by.
	DatasetID param.Field[string] `query:"dataset_id"`
	// Filter to a specific evaluator (UUID). Mutually exclusive with group_by.
	EvaluatorID param.Field[string] `query:"evaluator_id"`
	// Filter grouped results by evaluator feedback key. Only valid with group_by.
	FeedbackKey param.Field[string] `query:"feedback_key"`
	// Aggregation mode: 'evaluator', 'resource', or 'run_rule'. Mutually exclusive
	// with entity filters.
	GroupBy param.Field[string] `query:"group_by"`
	// Filter grouped results to evaluators attached to all supplied project or dataset
	// IDs. Only valid with group_by.
	ResourceID param.Field[[]string] `query:"resource_id"`
	// Filter to a specific project (UUID). Mutually exclusive with group_by.
	SessionID param.Field[string] `query:"session_id"`
	// Filter grouped results by evaluator type: 'llm' or 'code'. Only valid with
	// group_by.
	Type param.Field[string] `query:"type"`
}

// URLQuery serializes [OnlineEvaluatorSpendParams]'s query parameters as
// `url.Values`.
func (r OnlineEvaluatorSpendParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
