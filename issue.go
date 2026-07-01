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
func (r *IssueService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *Issue, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/platform/issues/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// **Beta:** This endpoint is in active development and may change without notice.
//
// Returns issues for the authenticated tenant, optionally filtered by session,
// status, severity, tag, or last modified time.
func (r *IssueService) List(ctx context.Context, query IssueListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationIssues[Issue], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "v1/platform/issues"
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
// status, severity, tag, or last modified time.
func (r *IssueService) ListAutoPaging(ctx context.Context, query IssueListParams, opts ...option.RequestOption) *pagination.OffsetPaginationIssuesAutoPager[Issue] {
	return pagination.NewOffsetPaginationIssuesAutoPager(r.List(ctx, query, opts...))
}

type Issue struct {
	ID                   string        `json:"id"`
	Actions              interface{}   `json:"actions"`
	CreatedAt            string        `json:"created_at"`
	Description          string        `json:"description"`
	FirstSeenAt          string        `json:"first_seen_at"`
	FixBranch            string        `json:"fix_branch"`
	FixDispatchedAt      string        `json:"fix_dispatched_at"`
	FixPrNumber          int64         `json:"fix_pr_number"`
	FixPrompt            string        `json:"fix_prompt"`
	FixVerification      interface{}   `json:"fix_verification"`
	LastSeenAt           string        `json:"last_seen_at"`
	Name                 string        `json:"name"`
	ProposedContextFixes []interface{} `json:"proposed_context_fixes"`
	ProposedExamples     []interface{} `json:"proposed_examples"`
	ProposedFix          string        `json:"proposed_fix"`
	ProposedPromptFixes  []interface{} `json:"proposed_prompt_fixes"`
	SessionID            string        `json:"session_id"`
	Severity             IssueSeverity `json:"severity"`
	Status               IssueStatus   `json:"status"`
	Tags                 []string      `json:"tags"`
	TenantID             string        `json:"tenant_id"`
	Traces               interface{}   `json:"traces"`
	UpdatedAt            string        `json:"updated_at"`
	JSON                 issueJSON     `json:"-"`
}

// issueJSON contains the JSON metadata for the struct [Issue]
type issueJSON struct {
	ID                   apijson.Field
	Actions              apijson.Field
	CreatedAt            apijson.Field
	Description          apijson.Field
	FirstSeenAt          apijson.Field
	FixBranch            apijson.Field
	FixDispatchedAt      apijson.Field
	FixPrNumber          apijson.Field
	FixPrompt            apijson.Field
	FixVerification      apijson.Field
	LastSeenAt           apijson.Field
	Name                 apijson.Field
	ProposedContextFixes apijson.Field
	ProposedExamples     apijson.Field
	ProposedFix          apijson.Field
	ProposedPromptFixes  apijson.Field
	SessionID            apijson.Field
	Severity             apijson.Field
	Status               apijson.Field
	Tags                 apijson.Field
	TenantID             apijson.Field
	Traces               apijson.Field
	UpdatedAt            apijson.Field
	raw                  string
	ExtraFields          map[string]apijson.Field
}

func (r *Issue) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r issueJSON) RawJSON() string {
	return r.raw
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
	IssueStatusCompleted IssueStatus = "completed"
	IssueStatusIgnored   IssueStatus = "ignored"
)

func (r IssueStatus) IsKnown() bool {
	switch r {
	case IssueStatusOpen, IssueStatusCompleted, IssueStatusIgnored:
		return true
	}
	return false
}

type IssueListParams struct {
	// Page size (positive integer; defaults to 50, capped at 500)
	Limit param.Field[int64] `query:"limit"`
	// Page offset (non-negative integer)
	Offset param.Field[int64] `query:"offset"`
	// Filter by session ID (UUID)
	SessionID param.Field[string] `query:"session_id"`
	// Filter by session name (exact match)
	SessionName param.Field[string] `query:"session_name"`
	// Filter by severity
	Severity param.Field[IssueListParamsSeverity] `query:"severity"`
	// Sort field
	SortBy param.Field[IssueListParamsSortBy] `query:"sort_by"`
	// Filter by status
	Status param.Field[IssueListParamsStatus] `query:"status"`
	// Filter by tag (exact match)
	Tag param.Field[string] `query:"tag"`
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

// Sort field
type IssueListParamsSortBy string

const (
	IssueListParamsSortByCreatedAt IssueListParamsSortBy = "created_at"
	IssueListParamsSortByUpdatedAt IssueListParamsSortBy = "updated_at"
	IssueListParamsSortBySeverity  IssueListParamsSortBy = "severity"
)

func (r IssueListParamsSortBy) IsKnown() bool {
	switch r {
	case IssueListParamsSortByCreatedAt, IssueListParamsSortByUpdatedAt, IssueListParamsSortBySeverity:
		return true
	}
	return false
}

// Filter by status
type IssueListParamsStatus string

const (
	IssueListParamsStatusOpen      IssueListParamsStatus = "open"
	IssueListParamsStatusCompleted IssueListParamsStatus = "completed"
	IssueListParamsStatusIgnored   IssueListParamsStatus = "ignored"
)

func (r IssueListParamsStatus) IsKnown() bool {
	switch r {
	case IssueListParamsStatusOpen, IssueListParamsStatusCompleted, IssueListParamsStatusIgnored:
		return true
	}
	return false
}
