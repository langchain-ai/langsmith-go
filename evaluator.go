// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
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

// EvaluatorService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewEvaluatorService] method instead.
type EvaluatorService struct {
	Options []option.RequestOption
}

// NewEvaluatorService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewEvaluatorService(opts ...option.RequestOption) (r *EvaluatorService) {
	r = &EvaluatorService{}
	r.Options = opts
	return
}

// List all run rules.
func (r *EvaluatorService) List(ctx context.Context, query EvaluatorListParams, opts ...option.RequestOption) (res *[]Evaluator, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/runs/rules"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type CodeEvaluatorTopLevel struct {
	Code     string                        `json:"code" api:"required"`
	Language CodeEvaluatorTopLevelLanguage `json:"language" api:"nullable"`
	JSON     codeEvaluatorTopLevelJSON     `json:"-"`
}

// codeEvaluatorTopLevelJSON contains the JSON metadata for the struct
// [CodeEvaluatorTopLevel]
type codeEvaluatorTopLevelJSON struct {
	Code        apijson.Field
	Language    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CodeEvaluatorTopLevel) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r codeEvaluatorTopLevelJSON) RawJSON() string {
	return r.raw
}

type CodeEvaluatorTopLevelLanguage string

const (
	CodeEvaluatorTopLevelLanguagePython     CodeEvaluatorTopLevelLanguage = "python"
	CodeEvaluatorTopLevelLanguageJavascript CodeEvaluatorTopLevelLanguage = "javascript"
)

func (r CodeEvaluatorTopLevelLanguage) IsKnown() bool {
	switch r {
	case CodeEvaluatorTopLevelLanguagePython, CodeEvaluatorTopLevelLanguageJavascript:
		return true
	}
	return false
}

// Run rules schema.
type Evaluator struct {
	ID                           string                    `json:"id" api:"required" format:"uuid"`
	CreatedAt                    time.Time                 `json:"created_at" api:"required" format:"date-time"`
	DisplayName                  string                    `json:"display_name" api:"required"`
	EvaluatorVersion             int64                     `json:"evaluator_version" api:"required"`
	SamplingRate                 float64                   `json:"sampling_rate" api:"required"`
	TenantID                     string                    `json:"tenant_id" api:"required" format:"uuid"`
	UpdatedAt                    time.Time                 `json:"updated_at" api:"required" format:"date-time"`
	Webhooks                     []EvaluatorWebhook        `json:"webhooks" api:"required,nullable"`
	AddToAnnotationQueueID       string                    `json:"add_to_annotation_queue_id" api:"nullable" format:"uuid"`
	AddToAnnotationQueueName     string                    `json:"add_to_annotation_queue_name" api:"nullable"`
	AddToDatasetID               string                    `json:"add_to_dataset_id" api:"nullable" format:"uuid"`
	AddToDatasetName             string                    `json:"add_to_dataset_name" api:"nullable"`
	AddToDatasetPreferCorrection bool                      `json:"add_to_dataset_prefer_correction"`
	Alerts                       []EvaluatorPagerdutyAlert `json:"alerts" api:"nullable"`
	AlignmentAnnotationQueueID   string                    `json:"alignment_annotation_queue_id" api:"nullable" format:"uuid"`
	BackfillCompletedAt          time.Time                 `json:"backfill_completed_at" api:"nullable" format:"date-time"`
	BackfillError                string                    `json:"backfill_error" api:"nullable"`
	BackfillFrom                 time.Time                 `json:"backfill_from" api:"nullable" format:"date-time"`
	BackfillID                   string                    `json:"backfill_id" api:"nullable" format:"uuid"`
	BackfillProgress             float64                   `json:"backfill_progress" api:"nullable"`
	BackfillStatus               string                    `json:"backfill_status" api:"nullable"`
	CodeEvaluators               []CodeEvaluatorTopLevel   `json:"code_evaluators" api:"nullable"`
	CorrectionsDatasetID         string                    `json:"corrections_dataset_id" api:"nullable" format:"uuid"`
	DatasetID                    string                    `json:"dataset_id" api:"nullable" format:"uuid"`
	DatasetName                  string                    `json:"dataset_name" api:"nullable"`
	EvaluatorID                  string                    `json:"evaluator_id" api:"nullable" format:"uuid"`
	Evaluators                   []EvaluatorTopLevel       `json:"evaluators" api:"nullable"`
	ExtendOnly                   bool                      `json:"extend_only"`
	Filter                       string                    `json:"filter" api:"nullable"`
	GroupBy                      EvaluatorGroupBy          `json:"group_by" api:"nullable"`
	IncludeExtendedStats         bool                      `json:"include_extended_stats"`
	IsEnabled                    bool                      `json:"is_enabled"`
	NumFewShotExamples           int64                     `json:"num_few_shot_examples" api:"nullable"`
	SessionID                    string                    `json:"session_id" api:"nullable" format:"uuid"`
	SessionName                  string                    `json:"session_name" api:"nullable"`
	TraceFilter                  string                    `json:"trace_filter" api:"nullable"`
	Transient                    bool                      `json:"transient"`
	TreeFilter                   string                    `json:"tree_filter" api:"nullable"`
	UseCorrectionsDataset        bool                      `json:"use_corrections_dataset"`
	JSON                         evaluatorJSON             `json:"-"`
}

// evaluatorJSON contains the JSON metadata for the struct [Evaluator]
type evaluatorJSON struct {
	ID                           apijson.Field
	CreatedAt                    apijson.Field
	DisplayName                  apijson.Field
	EvaluatorVersion             apijson.Field
	SamplingRate                 apijson.Field
	TenantID                     apijson.Field
	UpdatedAt                    apijson.Field
	Webhooks                     apijson.Field
	AddToAnnotationQueueID       apijson.Field
	AddToAnnotationQueueName     apijson.Field
	AddToDatasetID               apijson.Field
	AddToDatasetName             apijson.Field
	AddToDatasetPreferCorrection apijson.Field
	Alerts                       apijson.Field
	AlignmentAnnotationQueueID   apijson.Field
	BackfillCompletedAt          apijson.Field
	BackfillError                apijson.Field
	BackfillFrom                 apijson.Field
	BackfillID                   apijson.Field
	BackfillProgress             apijson.Field
	BackfillStatus               apijson.Field
	CodeEvaluators               apijson.Field
	CorrectionsDatasetID         apijson.Field
	DatasetID                    apijson.Field
	DatasetName                  apijson.Field
	EvaluatorID                  apijson.Field
	Evaluators                   apijson.Field
	ExtendOnly                   apijson.Field
	Filter                       apijson.Field
	GroupBy                      apijson.Field
	IncludeExtendedStats         apijson.Field
	IsEnabled                    apijson.Field
	NumFewShotExamples           apijson.Field
	SessionID                    apijson.Field
	SessionName                  apijson.Field
	TraceFilter                  apijson.Field
	Transient                    apijson.Field
	TreeFilter                   apijson.Field
	UseCorrectionsDataset        apijson.Field
	raw                          string
	ExtraFields                  map[string]apijson.Field
}

func (r *Evaluator) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluatorJSON) RawJSON() string {
	return r.raw
}

type EvaluatorGroupBy string

const (
	EvaluatorGroupByThreadID EvaluatorGroupBy = "thread_id"
)

func (r EvaluatorGroupBy) IsKnown() bool {
	switch r {
	case EvaluatorGroupByThreadID:
		return true
	}
	return false
}

type EvaluatorPagerdutyAlert struct {
	RoutingKey string `json:"routing_key" api:"required"`
	// Enum for severity.
	Severity EvaluatorPagerdutyAlertSeverity `json:"severity" api:"nullable"`
	Summary  string                          `json:"summary" api:"nullable"`
	// Enum for alert types.
	Type EvaluatorPagerdutyAlertType `json:"type" api:"nullable"`
	JSON evaluatorPagerdutyAlertJSON `json:"-"`
}

// evaluatorPagerdutyAlertJSON contains the JSON metadata for the struct
// [EvaluatorPagerdutyAlert]
type evaluatorPagerdutyAlertJSON struct {
	RoutingKey  apijson.Field
	Severity    apijson.Field
	Summary     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluatorPagerdutyAlert) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluatorPagerdutyAlertJSON) RawJSON() string {
	return r.raw
}

// Enum for severity.
type EvaluatorPagerdutyAlertSeverity string

const (
	EvaluatorPagerdutyAlertSeverityCritical EvaluatorPagerdutyAlertSeverity = "critical"
	EvaluatorPagerdutyAlertSeverityWarning  EvaluatorPagerdutyAlertSeverity = "warning"
	EvaluatorPagerdutyAlertSeverityError    EvaluatorPagerdutyAlertSeverity = "error"
	EvaluatorPagerdutyAlertSeverityInfo     EvaluatorPagerdutyAlertSeverity = "info"
)

func (r EvaluatorPagerdutyAlertSeverity) IsKnown() bool {
	switch r {
	case EvaluatorPagerdutyAlertSeverityCritical, EvaluatorPagerdutyAlertSeverityWarning, EvaluatorPagerdutyAlertSeverityError, EvaluatorPagerdutyAlertSeverityInfo:
		return true
	}
	return false
}

// Enum for alert types.
type EvaluatorPagerdutyAlertType string

const (
	EvaluatorPagerdutyAlertTypePagerduty EvaluatorPagerdutyAlertType = "pagerduty"
)

func (r EvaluatorPagerdutyAlertType) IsKnown() bool {
	switch r {
	case EvaluatorPagerdutyAlertTypePagerduty:
		return true
	}
	return false
}

type EvaluatorTopLevel struct {
	// Evaluator structured output schema.
	Structured EvaluatorTopLevelStructured `json:"structured" api:"required"`
	JSON       evaluatorTopLevelJSON       `json:"-"`
}

// evaluatorTopLevelJSON contains the JSON metadata for the struct
// [EvaluatorTopLevel]
type evaluatorTopLevelJSON struct {
	Structured  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluatorTopLevel) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluatorTopLevelJSON) RawJSON() string {
	return r.raw
}

// Evaluator structured output schema.
type EvaluatorTopLevelStructured struct {
	HubRef          string                          `json:"hub_ref" api:"nullable"`
	Model           map[string]interface{}          `json:"model" api:"nullable"`
	Prompt          [][]interface{}                 `json:"prompt" api:"nullable"`
	Schema          map[string]interface{}          `json:"schema" api:"nullable"`
	TemplateFormat  string                          `json:"template_format" api:"nullable"`
	VariableMapping map[string]string               `json:"variable_mapping" api:"nullable"`
	JSON            evaluatorTopLevelStructuredJSON `json:"-"`
}

// evaluatorTopLevelStructuredJSON contains the JSON metadata for the struct
// [EvaluatorTopLevelStructured]
type evaluatorTopLevelStructuredJSON struct {
	HubRef          apijson.Field
	Model           apijson.Field
	Prompt          apijson.Field
	Schema          apijson.Field
	TemplateFormat  apijson.Field
	VariableMapping apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *EvaluatorTopLevelStructured) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluatorTopLevelStructuredJSON) RawJSON() string {
	return r.raw
}

type EvaluatorWebhook struct {
	URL     string               `json:"url" api:"required"`
	Headers map[string]string    `json:"headers" api:"nullable"`
	JSON    evaluatorWebhookJSON `json:"-"`
}

// evaluatorWebhookJSON contains the JSON metadata for the struct
// [EvaluatorWebhook]
type evaluatorWebhookJSON struct {
	URL         apijson.Field
	Headers     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *EvaluatorWebhook) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r evaluatorWebhookJSON) RawJSON() string {
	return r.raw
}

type EvaluatorListParams struct {
	ID                      param.Field[[]string]                `query:"id" format:"uuid"`
	DatasetID               param.Field[string]                  `query:"dataset_id" format:"uuid"`
	EvaluatorID             param.Field[string]                  `query:"evaluator_id" format:"uuid"`
	IncludeBackfillProgress param.Field[bool]                    `query:"include_backfill_progress"`
	NameContains            param.Field[string]                  `query:"name_contains"`
	SessionID               param.Field[string]                  `query:"session_id" format:"uuid"`
	TagValueID              param.Field[[]string]                `query:"tag_value_id" format:"uuid"`
	Type                    param.Field[EvaluatorListParamsType] `query:"type"`
}

// URLQuery serializes [EvaluatorListParams]'s query parameters as `url.Values`.
func (r EvaluatorListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type EvaluatorListParamsType string

const (
	EvaluatorListParamsTypeSession EvaluatorListParamsType = "session"
	EvaluatorListParamsTypeDataset EvaluatorListParamsType = "dataset"
)

func (r EvaluatorListParamsType) IsKnown() bool {
	switch r {
	case EvaluatorListParamsTypeSession, EvaluatorListParamsTypeDataset:
		return true
	}
	return false
}
