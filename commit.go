// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
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

// Creates a new commit in a repository. Requires authentication and write access
// to the repository.
func (r *CommitService) New(ctx context.Context, owner interface{}, repo interface{}, body CommitNewParams, opts ...option.RequestOption) (res *CommitNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("commits/%v/%v", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Retrieves a specific commit by hash, tag, or "latest" for a repository. This
// endpoint supports both authenticated and unauthenticated access. Authenticated
// users can access private repos, while unauthenticated users can only access
// public repos. Commit resolution logic:
//
// - "latest" or empty: Get the most recent commit
// - Less than 8 characters: Only check for tags
// - 8 or more characters: Prioritize commit hash over tag, check both
func (r *CommitService) Get(ctx context.Context, owner interface{}, repo interface{}, commit interface{}, query CommitGetParams, opts ...option.RequestOption) (res *CommitGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := fmt.Sprintf("commits/%v/%v/%v", owner, repo, commit)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Lists all commits for a repository with pagination support. This endpoint
// supports both authenticated and unauthenticated access. Authenticated users can
// access private repos, while unauthenticated users can only access public repos.
func (r *CommitService) List(ctx context.Context, owner interface{}, repo interface{}, query CommitListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationCommits[CommitWithLookups], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := fmt.Sprintf("commits/%v/%v", owner, repo)
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

// Lists all commits for a repository with pagination support. This endpoint
// supports both authenticated and unauthenticated access. Authenticated users can
// access private repos, while unauthenticated users can only access public repos.
func (r *CommitService) ListAutoPaging(ctx context.Context, owner interface{}, repo interface{}, query CommitListParams, opts ...option.RequestOption) *pagination.OffsetPaginationCommitsAutoPager[CommitWithLookups] {
	return pagination.NewOffsetPaginationCommitsAutoPager(r.List(ctx, owner, repo, query, opts...))
}

// Response model for get_commit_manifest.
type CommitManifestResponse struct {
	CommitHash string                          `json:"commit_hash,required"`
	Manifest   map[string]interface{}          `json:"manifest,required"`
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
	Inputs    map[string]interface{}            `json:"inputs,nullable"`
	Outputs   map[string]interface{}            `json:"outputs,nullable"`
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

type CommitWithLookups struct {
	// The commit ID
	ID string `json:"id" format:"uuid"`
	// The hash of the commit
	CommitHash string `json:"commit_hash"`
	// When the commit was created
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// Example run IDs associated with the commit
	ExampleRunIDs []string `json:"example_run_ids" format:"uuid"`
	// Author's full name
	FullName string `json:"full_name"`
	// The manifest of the commit
	Manifest interface{} `json:"manifest"`
	// The SHA of the manifest
	ManifestSha []int64 `json:"manifest_sha"`
	// Number of API downloads
	NumDownloads int64 `json:"num_downloads"`
	// Number of web views
	NumViews int64 `json:"num_views"`
	// The hash of the parent commit
	ParentCommitHash string `json:"parent_commit_hash"`
	// The ID of the parent commit
	ParentID string `json:"parent_id" format:"uuid"`
	// Repository ID
	RepoID string `json:"repo_id" format:"uuid"`
	// When the commit was last updated
	UpdatedAt time.Time             `json:"updated_at" format:"date-time"`
	JSON      commitWithLookupsJSON `json:"-"`
}

// commitWithLookupsJSON contains the JSON metadata for the struct
// [CommitWithLookups]
type commitWithLookupsJSON struct {
	ID               apijson.Field
	CommitHash       apijson.Field
	CreatedAt        apijson.Field
	ExampleRunIDs    apijson.Field
	FullName         apijson.Field
	Manifest         apijson.Field
	ManifestSha      apijson.Field
	NumDownloads     apijson.Field
	NumViews         apijson.Field
	ParentCommitHash apijson.Field
	ParentID         apijson.Field
	RepoID           apijson.Field
	UpdatedAt        apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *CommitWithLookups) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitWithLookupsJSON) RawJSON() string {
	return r.raw
}

type CommitNewResponse struct {
	Commit CommitWithLookups     `json:"commit"`
	JSON   commitNewResponseJSON `json:"-"`
}

// commitNewResponseJSON contains the JSON metadata for the struct
// [CommitNewResponse]
type commitNewResponseJSON struct {
	Commit      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitNewResponseJSON) RawJSON() string {
	return r.raw
}

type CommitGetResponse struct {
	CommitHash string                     `json:"commit_hash"`
	Examples   []CommitGetResponseExample `json:"examples"`
	Manifest   interface{}                `json:"manifest"`
	JSON       commitGetResponseJSON      `json:"-"`
}

// commitGetResponseJSON contains the JSON metadata for the struct
// [CommitGetResponse]
type commitGetResponseJSON struct {
	CommitHash  apijson.Field
	Examples    apijson.Field
	Manifest    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitGetResponseJSON) RawJSON() string {
	return r.raw
}

type CommitGetResponseExample struct {
	ID        string                       `json:"id" format:"uuid"`
	Inputs    interface{}                  `json:"inputs"`
	Outputs   interface{}                  `json:"outputs"`
	SessionID string                       `json:"session_id" format:"uuid"`
	StartTime string                       `json:"start_time"`
	JSON      commitGetResponseExampleJSON `json:"-"`
}

// commitGetResponseExampleJSON contains the JSON metadata for the struct
// [CommitGetResponseExample]
type commitGetResponseExampleJSON struct {
	ID          apijson.Field
	Inputs      apijson.Field
	Outputs     apijson.Field
	SessionID   apijson.Field
	StartTime   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CommitGetResponseExample) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r commitGetResponseExampleJSON) RawJSON() string {
	return r.raw
}

type CommitNewParams struct {
	Manifest     param.Field[interface{}] `json:"manifest"`
	ParentCommit param.Field[string]      `json:"parent_commit"`
	// SkipWebhooks allows skipping webhook notifications. Can be true (boolean) to
	// skip all, or an array of webhook UUIDs to skip specific ones.
	SkipWebhooks param.Field[interface{}] `json:"skip_webhooks"`
}

func (r CommitNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
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

type CommitListParams struct {
	// Limit is the pagination limit
	Limit param.Field[int64] `query:"limit"`
	// Offset is the pagination offset
	Offset param.Field[int64] `query:"offset"`
}

// URLQuery serializes [CommitListParams]'s query parameters as `url.Values`.
func (r CommitListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
