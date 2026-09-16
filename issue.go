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

// IssueService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewIssueService] method instead.
type IssueService struct {
	Options []option.RequestOption
}

// NewIssueService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewIssueService(opts ...option.RequestOption) (r *IssueService) {
	r = &IssueService{}
	r.Options = opts
	return
}

// **Beta:** This endpoint is in active development and may change without notice.
//
// Returns one issue for the authenticated tenant.
func (r *IssueService) Get(ctx context.Context, id string, query IssueGetParams, opts ...option.RequestOption) (res *Issue, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/platform/issues/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// **Beta:** This endpoint is in active development and may change without notice.
//
// Returns issues for the authenticated tenant, optionally filtered by session,
// status, severity, tag, linked trace, or last modified time.
func (r *IssueService) List(ctx context.Context, query IssueListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationIssues[Issue], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/platform/issues"
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

// **Beta:** This endpoint is in active development and may change without notice.
//
// Returns issues for the authenticated tenant, optionally filtered by session,
// status, severity, tag, linked trace, or last modified time.
func (r *IssueService) ListAutoPaging(ctx context.Context, query IssueListParams, opts ...option.RequestOption) *pagination.OffsetPaginationIssuesAutoPager[Issue] {
	return pagination.NewOffsetPaginationIssuesAutoPager(r.List(ctx, query, opts...))
}

type Issue struct {
	ID                     string      `json:"id"`
	Actions                interface{} `json:"actions"`
	AutoResolutionEvidence interface{} `json:"auto_resolution_evidence"`
	// Nil unless eligible: "auto_close" or "prompt". Evidence carries the deciding
	// gate.
	AutoResolutionState  string               `json:"auto_resolution_state"`
	CreatedAt            string               `json:"created_at"`
	Description          string               `json:"description"`
	FirstSeenAt          string               `json:"first_seen_at"`
	FixBranch            string               `json:"fix_branch"`
	FixDispatchedAt      string               `json:"fix_dispatched_at"`
	FixPrNumber          int64                `json:"fix_pr_number"`
	FixPrompt            string               `json:"fix_prompt"`
	FixVerification      IssueFixVerification `json:"fix_verification"`
	LastSeenAt           string               `json:"last_seen_at"`
	LinearContext        IssueLinearContext   `json:"linear_context"`
	LinearSync           IssueLinearSync      `json:"linear_sync"`
	Name                 string               `json:"name"`
	ProposedContextFixes []interface{}        `json:"proposed_context_fixes"`
	ProposedExamples     []interface{}        `json:"proposed_examples"`
	ProposedFix          string               `json:"proposed_fix"`
	ProposedPromptFixes  []interface{}        `json:"proposed_prompt_fixes"`
	// RecurrencesSinceWatching counts linked traces whose run start_time is after
	// watching_since — i.e. recurrences observed during the current watch period.
	RecurrencesSinceWatching int64                 `json:"recurrences_since_watching"`
	SessionID                string                `json:"session_id"`
	Severity                 IssueSeverity         `json:"severity"`
	Status                   IssueStatus           `json:"status"`
	Tags                     []string              `json:"tags"`
	TenantID                 string                `json:"tenant_id"`
	Traces                   interface{}           `json:"traces"`
	UpdatedAt                string                `json:"updated_at"`
	ValidationResult         IssueValidationResult `json:"validation_result"`
	WatchingSince            string                `json:"watching_since"`
	JSON                     issueJSON             `json:"-"`
}

// issueJSON contains the JSON metadata for the struct [Issue]
type issueJSON struct {
	ID                       apijson.Field
	Actions                  apijson.Field
	AutoResolutionEvidence   apijson.Field
	AutoResolutionState      apijson.Field
	CreatedAt                apijson.Field
	Description              apijson.Field
	FirstSeenAt              apijson.Field
	FixBranch                apijson.Field
	FixDispatchedAt          apijson.Field
	FixPrNumber              apijson.Field
	FixPrompt                apijson.Field
	FixVerification          apijson.Field
	LastSeenAt               apijson.Field
	LinearContext            apijson.Field
	LinearSync               apijson.Field
	Name                     apijson.Field
	ProposedContextFixes     apijson.Field
	ProposedExamples         apijson.Field
	ProposedFix              apijson.Field
	ProposedPromptFixes      apijson.Field
	RecurrencesSinceWatching apijson.Field
	SessionID                apijson.Field
	Severity                 apijson.Field
	Status                   apijson.Field
	Tags                     apijson.Field
	TenantID                 apijson.Field
	Traces                   apijson.Field
	UpdatedAt                apijson.Field
	ValidationResult         apijson.Field
	WatchingSince            apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *Issue) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r issueJSON) RawJSON() string {
	return r.raw
}

type IssueFixVerification struct {
	Attempt             int64                      `json:"attempt"`
	ParentDeploymentID  string                     `json:"parent_deployment_id" format:"uuid"`
	PreviewDeploymentID string                     `json:"preview_deployment_id" format:"uuid"`
	Reason              string                     `json:"reason"`
	RootTraceIDs        []string                   `json:"root_trace_ids"`
	Status              IssueFixVerificationStatus `json:"status"`
	UpdatedAt           time.Time                  `json:"updated_at" format:"date-time"`
	JSON                issueFixVerificationJSON   `json:"-"`
}

// issueFixVerificationJSON contains the JSON metadata for the struct
// [IssueFixVerification]
type issueFixVerificationJSON struct {
	Attempt             apijson.Field
	ParentDeploymentID  apijson.Field
	PreviewDeploymentID apijson.Field
	Reason              apijson.Field
	RootTraceIDs        apijson.Field
	Status              apijson.Field
	UpdatedAt           apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *IssueFixVerification) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r issueFixVerificationJSON) RawJSON() string {
	return r.raw
}

type IssueFixVerificationStatus string

const (
	IssueFixVerificationStatusAwaitingPreview IssueFixVerificationStatus = "awaiting_preview"
	IssueFixVerificationStatusVerifying       IssueFixVerificationStatus = "verifying"
	IssueFixVerificationStatusPassed          IssueFixVerificationStatus = "passed"
	IssueFixVerificationStatusFailed          IssueFixVerificationStatus = "failed"
	IssueFixVerificationStatusInconclusive    IssueFixVerificationStatus = "inconclusive"
	IssueFixVerificationStatusTimeout         IssueFixVerificationStatus = "timeout"
	IssueFixVerificationStatusError           IssueFixVerificationStatus = "error"
)

func (r IssueFixVerificationStatus) IsKnown() bool {
	switch r {
	case IssueFixVerificationStatusAwaitingPreview, IssueFixVerificationStatusVerifying, IssueFixVerificationStatusPassed, IssueFixVerificationStatusFailed, IssueFixVerificationStatusInconclusive, IssueFixVerificationStatusTimeout, IssueFixVerificationStatusError:
		return true
	}
	return false
}

type IssueLinearContext struct {
	GitHubPrURLs  []string               `json:"github_pr_urls"`
	WorkflowState string                 `json:"workflow_state"`
	JSON          issueLinearContextJSON `json:"-"`
}

// issueLinearContextJSON contains the JSON metadata for the struct
// [IssueLinearContext]
type issueLinearContextJSON struct {
	GitHubPrURLs  apijson.Field
	WorkflowState apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *IssueLinearContext) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r issueLinearContextJSON) RawJSON() string {
	return r.raw
}

type IssueLinearSync struct {
	Identifier      string               `json:"identifier"`
	IssueID         string               `json:"issue_id"`
	LastAttemptedAt time.Time            `json:"last_attempted_at" format:"date-time"`
	LastError       string               `json:"last_error"`
	LastSyncedAt    time.Time            `json:"last_synced_at" format:"date-time"`
	LinearIssueID   string               `json:"linear_issue_id"`
	State           IssueLinearSyncState `json:"state"`
	URL             string               `json:"url"`
	JSON            issueLinearSyncJSON  `json:"-"`
}

// issueLinearSyncJSON contains the JSON metadata for the struct [IssueLinearSync]
type issueLinearSyncJSON struct {
	Identifier      apijson.Field
	IssueID         apijson.Field
	LastAttemptedAt apijson.Field
	LastError       apijson.Field
	LastSyncedAt    apijson.Field
	LinearIssueID   apijson.Field
	State           apijson.Field
	URL             apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *IssueLinearSync) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r issueLinearSyncJSON) RawJSON() string {
	return r.raw
}

type IssueLinearSyncState string

const (
	IssueLinearSyncStatePending      IssueLinearSyncState = "pending"
	IssueLinearSyncStateSynced       IssueLinearSyncState = "synced"
	IssueLinearSyncStateFailed       IssueLinearSyncState = "failed"
	IssueLinearSyncStateAuthRequired IssueLinearSyncState = "auth_required"
	IssueLinearSyncStatePaused       IssueLinearSyncState = "paused"
)

func (r IssueLinearSyncState) IsKnown() bool {
	switch r {
	case IssueLinearSyncStatePending, IssueLinearSyncStateSynced, IssueLinearSyncStateFailed, IssueLinearSyncStateAuthRequired, IssueLinearSyncStatePaused:
		return true
	}
	return false
}

type IssueSeverity int64

const (
	IssueSeverity0 IssueSeverity = 0
	IssueSeverity1 IssueSeverity = 1
	IssueSeverity2 IssueSeverity = 2
	IssueSeverity3 IssueSeverity = 3
)

func (r IssueSeverity) IsKnown() bool {
	switch r {
	case IssueSeverity0, IssueSeverity1, IssueSeverity2, IssueSeverity3:
		return true
	}
	return false
}

type IssueStatus string

const (
	IssueStatusOpen      IssueStatus = "open"
	IssueStatusFixing    IssueStatus = "fixing"
	IssueStatusWatching  IssueStatus = "watching"
	IssueStatusCompleted IssueStatus = "completed"
	IssueStatusIgnored   IssueStatus = "ignored"
)

func (r IssueStatus) IsKnown() bool {
	switch r {
	case IssueStatusOpen, IssueStatusFixing, IssueStatusWatching, IssueStatusCompleted, IssueStatusIgnored:
		return true
	}
	return false
}

type IssueValidationResult struct {
	ActiveRevisionID string                       `json:"active_revision_id" format:"uuid"`
	CompletedAt      time.Time                    `json:"completed_at" format:"date-time"`
	DeploymentID     string                       `json:"deployment_id" format:"uuid"`
	Outcome          IssueValidationResultOutcome `json:"outcome"`
	Reason           string                       `json:"reason"`
	RootTraceIDs     []string                     `json:"root_trace_ids"`
	JSON             issueValidationResultJSON    `json:"-"`
}

// issueValidationResultJSON contains the JSON metadata for the struct
// [IssueValidationResult]
type issueValidationResultJSON struct {
	ActiveRevisionID apijson.Field
	CompletedAt      apijson.Field
	DeploymentID     apijson.Field
	Outcome          apijson.Field
	Reason           apijson.Field
	RootTraceIDs     apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *IssueValidationResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r issueValidationResultJSON) RawJSON() string {
	return r.raw
}

type IssueValidationResultOutcome string

const (
	IssueValidationResultOutcomeReproduced    IssueValidationResultOutcome = "reproduced"
	IssueValidationResultOutcomeNotReproduced IssueValidationResultOutcome = "not_reproduced"
	IssueValidationResultOutcomeInconclusive  IssueValidationResultOutcome = "inconclusive"
	IssueValidationResultOutcomeError         IssueValidationResultOutcome = "error"
)

func (r IssueValidationResultOutcome) IsKnown() bool {
	switch r {
	case IssueValidationResultOutcomeReproduced, IssueValidationResultOutcomeNotReproduced, IssueValidationResultOutcomeInconclusive, IssueValidationResultOutcomeError:
		return true
	}
	return false
}

type IssueGetParams struct {
	// Include current Linear workflow state and validated linked GitHub pull request
	// URLs
	IncludeLinearContext param.Field[bool] `query:"include_linear_context"`
}

// URLQuery serializes [IssueGetParams]'s query parameters as `url.Values`.
func (r IssueGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type IssueListParams struct {
	// Filter by Engine activity (repeatable; OR semantics)
	Activity param.Field[[]IssueListParamsActivity] `query:"activity"`
	// Page size (positive integer; defaults to 50, capped at 500)
	Limit param.Field[int64] `query:"limit"`
	// Page offset (non-negative integer; at most 100000)
	Offset param.Field[int64] `query:"offset"`
	// Filter by session ID (UUID)
	SessionID param.Field[string] `query:"session_id"`
	// Filter by session name (exact match)
	SessionName param.Field[string] `query:"session_name"`
	// Filter by severity
	Severity param.Field[IssueListParamsSeverity] `query:"severity"`
	// Filter by exact severity (repeatable; OR semantics)
	SeverityExact param.Field[[]IssueListParamsSeverityExact] `query:"severity_exact"`
	// Sort field
	SortBy param.Field[IssueListParamsSortBy] `query:"sort_by"`
	// Filter by status
	Status param.Field[IssueListParamsStatus] `query:"status"`
	// Group results by issue lifecycle status before applying sort_by
	StatusFirst param.Field[bool] `query:"status_first"`
	// Filter by tag (exact match)
	Tag param.Field[string] `query:"tag"`
	// Return only issues with a linked run in this trace
	TraceID param.Field[string] `query:"trace_id" format:"uuid"`
	// Return only issues updated at or after this RFC3339 timestamp
	UpdatedAt param.Field[string] `query:"updated_at"`
}

// URLQuery serializes [IssueListParams]'s query parameters as `url.Values`.
func (r IssueListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type IssueListParamsActivity string

const (
	IssueListParamsActivityFixing   IssueListParamsActivity = "fixing"
	IssueListParamsActivityWatching IssueListParamsActivity = "watching"
	IssueListParamsActivityRecurred IssueListParamsActivity = "recurred"
)

func (r IssueListParamsActivity) IsKnown() bool {
	switch r {
	case IssueListParamsActivityFixing, IssueListParamsActivityWatching, IssueListParamsActivityRecurred:
		return true
	}
	return false
}

// Filter by severity
type IssueListParamsSeverity int64

const (
	IssueListParamsSeverity0 IssueListParamsSeverity = 0
	IssueListParamsSeverity1 IssueListParamsSeverity = 1
	IssueListParamsSeverity2 IssueListParamsSeverity = 2
	IssueListParamsSeverity3 IssueListParamsSeverity = 3
)

func (r IssueListParamsSeverity) IsKnown() bool {
	switch r {
	case IssueListParamsSeverity0, IssueListParamsSeverity1, IssueListParamsSeverity2, IssueListParamsSeverity3:
		return true
	}
	return false
}

type IssueListParamsSeverityExact int64

const (
	IssueListParamsSeverityExact0 IssueListParamsSeverityExact = 0
	IssueListParamsSeverityExact1 IssueListParamsSeverityExact = 1
	IssueListParamsSeverityExact2 IssueListParamsSeverityExact = 2
	IssueListParamsSeverityExact3 IssueListParamsSeverityExact = 3
)

func (r IssueListParamsSeverityExact) IsKnown() bool {
	switch r {
	case IssueListParamsSeverityExact0, IssueListParamsSeverityExact1, IssueListParamsSeverityExact2, IssueListParamsSeverityExact3:
		return true
	}
	return false
}

// Sort field
type IssueListParamsSortBy string

const (
	IssueListParamsSortByDefault     IssueListParamsSortBy = "default"
	IssueListParamsSortByCreatedAt   IssueListParamsSortBy = "created_at"
	IssueListParamsSortByUpdatedAt   IssueListParamsSortBy = "updated_at"
	IssueListParamsSortByLastSeen    IssueListParamsSortBy = "last_seen"
	IssueListParamsSortByLastUpdated IssueListParamsSortBy = "last_updated"
	IssueListParamsSortByTraceCount  IssueListParamsSortBy = "trace_count"
	IssueListParamsSortBySeverity    IssueListParamsSortBy = "severity"
)

func (r IssueListParamsSortBy) IsKnown() bool {
	switch r {
	case IssueListParamsSortByDefault, IssueListParamsSortByCreatedAt, IssueListParamsSortByUpdatedAt, IssueListParamsSortByLastSeen, IssueListParamsSortByLastUpdated, IssueListParamsSortByTraceCount, IssueListParamsSortBySeverity:
		return true
	}
	return false
}

// Filter by status
type IssueListParamsStatus string

const (
	IssueListParamsStatusOpen      IssueListParamsStatus = "open"
	IssueListParamsStatusFixing    IssueListParamsStatus = "fixing"
	IssueListParamsStatusWatching  IssueListParamsStatus = "watching"
	IssueListParamsStatusCompleted IssueListParamsStatus = "completed"
	IssueListParamsStatusIgnored   IssueListParamsStatus = "ignored"
)

func (r IssueListParamsStatus) IsKnown() bool {
	switch r {
	case IssueListParamsStatusOpen, IssueListParamsStatusFixing, IssueListParamsStatusWatching, IssueListParamsStatusCompleted, IssueListParamsStatusIgnored:
		return true
	}
	return false
}
