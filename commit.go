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

// CommitService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewCommitService] method instead.
type CommitService struct {
	Options []option.RequestOption
}

// NewCommitService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewCommitService(opts ...option.RequestOption) (r *CommitService) {
	r = &CommitService{}
	r.Options = opts
	return
}

// Download a repo.
func (r *CommitService) Get(ctx context.Context, owner string, repo string, commit string, query CommitGetParams, opts ...option.RequestOption) (res *CommitManifestResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return
	}
	if commit == "" {
		err = errors.New("missing required commit parameter")
		return
	}
	path := fmt.Sprintf("api/v1/commits/%s/%s/%s", owner, repo, commit)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Upload a repo.
func (r *CommitService) Update(ctx context.Context, owner string, repo string, body CommitUpdateParams, opts ...option.RequestOption) (res *CommitUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return
	}
	path := fmt.Sprintf("api/v1/commits/%s/%s", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get all commits.
func (r *CommitService) List(ctx context.Context, owner string, repo string, query CommitListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationCommits[CommitWithLookups], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return
	}
	path := fmt.Sprintf("api/v1/commits/%s/%s", owner, repo)
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

// Get all commits.
func (r *CommitService) ListAutoPaging(ctx context.Context, owner string, repo string, query CommitListParams, opts ...option.RequestOption) *pagination.OffsetPaginationCommitsAutoPager[CommitWithLookups] {
	return pagination.NewOffsetPaginationCommitsAutoPager(r.List(ctx, owner, repo, query, opts...))
}

// Response model for get_commit_manifest.
type CommitManifestResponse struct {
	CommitHash string                          `json:"commit_hash,required"`
	Manifest   interface{}                     `json:"manifest,required"`
	Examples   []CommitManifestResponseExample `json:"examples,nullable"`
	JSON       commitManifestResponseJSON      `json:"-"`
}

// commitManifestResponseJSON contains the JSON metadata for the struct
// [CommitManifestResponse]
type commitManifestResponseJSON struct {
	CommitHash  apijson.Field
	Manifest    apijson.Field
	Examples    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitManifestResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitManifestResponseJSON) RawJSON() string {
	return r.raw
}

// Response model for example runs
type CommitManifestResponseExample struct {
	ID        string                            `json:"id,required" format:"uuid"`
	SessionID string                            `json:"session_id,required" format:"uuid"`
	Inputs    interface{}                       `json:"inputs,nullable"`
	Outputs   interface{}                       `json:"outputs,nullable"`
	StartTime time.Time                         `json:"start_time,nullable" format:"date-time"`
	JSON      commitManifestResponseExampleJSON `json:"-"`
}

// commitManifestResponseExampleJSON contains the JSON metadata for the struct
// [CommitManifestResponseExample]
type commitManifestResponseExampleJSON struct {
	ID          apijson.Field
	SessionID   apijson.Field
	Inputs      apijson.Field
	Outputs     apijson.Field
	StartTime   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitManifestResponseExample) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitManifestResponseExampleJSON) RawJSON() string {
	return r.raw
}

// All database fields for commits, plus helpful computed fields and user info for
// private prompts.
type CommitWithLookups struct {
	ID               string                `json:"id,required" format:"uuid"`
	CommitHash       string                `json:"commit_hash,required"`
	CreatedAt        time.Time             `json:"created_at,required" format:"date-time"`
	ExampleRunIDs    []string              `json:"example_run_ids,required" format:"uuid"`
	Manifest         interface{}           `json:"manifest,required"`
	NumDownloads     int64                 `json:"num_downloads,required"`
	NumViews         int64                 `json:"num_views,required"`
	RepoID           string                `json:"repo_id,required" format:"uuid"`
	UpdatedAt        time.Time             `json:"updated_at,required" format:"date-time"`
	FullName         string                `json:"full_name,nullable"`
	ParentCommitHash string                `json:"parent_commit_hash,nullable"`
	ParentID         string                `json:"parent_id,nullable" format:"uuid"`
	JSON             commitWithLookupsJSON `json:"-"`
}

// commitWithLookupsJSON contains the JSON metadata for the struct
// [CommitWithLookups]
type commitWithLookupsJSON struct {
	ID               apijson.Field
	CommitHash       apijson.Field
	CreatedAt        apijson.Field
	ExampleRunIDs    apijson.Field
	Manifest         apijson.Field
	NumDownloads     apijson.Field
	NumViews         apijson.Field
	RepoID           apijson.Field
	UpdatedAt        apijson.Field
	FullName         apijson.Field
	ParentCommitHash apijson.Field
	ParentID         apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *CommitWithLookups) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitWithLookupsJSON) RawJSON() string {
	return r.raw
}

type CommitUpdateResponse struct {
	// All database fields for commits, plus helpful computed fields and user info for
	// private prompts.
	Commit CommitWithLookups        `json:"commit,required"`
	JSON   commitUpdateResponseJSON `json:"-"`
}

// commitUpdateResponseJSON contains the JSON metadata for the struct
// [CommitUpdateResponse]
type commitUpdateResponseJSON struct {
	Commit      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type CommitGetParams struct {
	GetExamples  param.Field[bool] `query:"get_examples"`
	IncludeModel param.Field[bool] `query:"include_model"`
	IsView       param.Field[bool] `query:"is_view"`
}

// URLQuery serializes [CommitGetParams]'s query parameters as `url.Values`.
func (r CommitGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type CommitUpdateParams struct {
	Manifest      param.Field[interface{}]                         `json:"manifest,required"`
	ExampleRunIDs param.Field[[]string]                            `json:"example_run_ids" format:"uuid"`
	ParentCommit  param.Field[string]                              `json:"parent_commit"`
	SkipWebhooks  param.Field[CommitUpdateParamsSkipWebhooksUnion] `json:"skip_webhooks" format:"uuid"`
}

func (r CommitUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [shared.UnionBool], [CommitUpdateParamsSkipWebhooksArray].
type CommitUpdateParamsSkipWebhooksUnion interface {
	ImplementsCommitUpdateParamsSkipWebhooksUnion()
}

type CommitUpdateParamsSkipWebhooksArray []string

func (r CommitUpdateParamsSkipWebhooksArray) ImplementsCommitUpdateParamsSkipWebhooksUnion() {}

type CommitListParams struct {
	Limit  param.Field[int64] `query:"limit"`
	Offset param.Field[int64] `query:"offset"`
}

// URLQuery serializes [CommitListParams]'s query parameters as `url.Values`.
func (r CommitListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
