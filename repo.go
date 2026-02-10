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

// RepoService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewRepoService] method instead.
type RepoService struct {
	Options []option.RequestOption
}

// NewRepoService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewRepoService(opts ...option.RequestOption) (r *RepoService) {
	r = &RepoService{}
	r.Options = opts
	return
}

// Create a repo.
func (r *RepoService) New(ctx context.Context, body RepoNewParams, opts ...option.RequestOption) (res *CreateRepoResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/repos"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get a repo.
func (r *RepoService) Get(ctx context.Context, owner string, repo string, opts ...option.RequestOption) (res *GetRepoResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return
	}
	path := fmt.Sprintf("api/v1/repos/%s/%s", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update a repo.
func (r *RepoService) Update(ctx context.Context, owner string, repo string, body RepoUpdateParams, opts ...option.RequestOption) (res *CreateRepoResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return
	}
	path := fmt.Sprintf("api/v1/repos/%s/%s", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// Get all repos.
func (r *RepoService) List(ctx context.Context, query RepoListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationRepos[RepoWithLookups], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/repos"
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

// Get all repos.
func (r *RepoService) ListAutoPaging(ctx context.Context, query RepoListParams, opts ...option.RequestOption) *pagination.OffsetPaginationReposAutoPager[RepoWithLookups] {
	return pagination.NewOffsetPaginationReposAutoPager(r.List(ctx, query, opts...))
}

// Delete a repo.
func (r *RepoService) Delete(ctx context.Context, owner string, repo string, opts ...option.RequestOption) (res *RepoDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return
	}
	path := fmt.Sprintf("api/v1/repos/%s/%s", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type CreateRepoResponse struct {
	// All database fields for repos, plus helpful computed fields.
	Repo RepoWithLookups        `json:"repo,required"`
	JSON createRepoResponseJSON `json:"-"`
}

// createRepoResponseJSON contains the JSON metadata for the struct
// [CreateRepoResponse]
type createRepoResponseJSON struct {
	Repo        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *CreateRepoResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r createRepoResponseJSON) RawJSON() string {
	return r.raw
}

type GetRepoResponse struct {
	// All database fields for repos, plus helpful computed fields.
	Repo RepoWithLookups     `json:"repo,required"`
	JSON getRepoResponseJSON `json:"-"`
}

// getRepoResponseJSON contains the JSON metadata for the struct [GetRepoResponse]
type getRepoResponseJSON struct {
	Repo        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *GetRepoResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r getRepoResponseJSON) RawJSON() string {
	return r.raw
}

// All database fields for repos, plus helpful computed fields.
type RepoWithLookups struct {
	ID             string    `json:"id,required" format:"uuid"`
	CreatedAt      time.Time `json:"created_at,required" format:"date-time"`
	FullName       string    `json:"full_name,required"`
	IsArchived     bool      `json:"is_archived,required"`
	IsPublic       bool      `json:"is_public,required"`
	NumCommits     int64     `json:"num_commits,required"`
	NumDownloads   int64     `json:"num_downloads,required"`
	NumLikes       int64     `json:"num_likes,required"`
	NumViews       int64     `json:"num_views,required"`
	Owner          string    `json:"owner,required,nullable"`
	RepoHandle     string    `json:"repo_handle,required"`
	Tags           []string  `json:"tags,required"`
	TenantID       string    `json:"tenant_id,required" format:"uuid"`
	UpdatedAt      time.Time `json:"updated_at,required" format:"date-time"`
	CommitTags     []string  `json:"commit_tags"`
	CreatedBy      string    `json:"created_by,nullable"`
	Description    string    `json:"description,nullable"`
	LastCommitHash string    `json:"last_commit_hash,nullable"`
	// Response model for get_commit_manifest.
	LatestCommitManifest CommitManifestResponse `json:"latest_commit_manifest,nullable"`
	LikedByAuthUser      bool                   `json:"liked_by_auth_user,nullable"`
	OriginalRepoFullName string                 `json:"original_repo_full_name,nullable"`
	OriginalRepoID       string                 `json:"original_repo_id,nullable" format:"uuid"`
	Readme               string                 `json:"readme,nullable"`
	UpstreamRepoFullName string                 `json:"upstream_repo_full_name,nullable"`
	UpstreamRepoID       string                 `json:"upstream_repo_id,nullable" format:"uuid"`
	JSON                 repoWithLookupsJSON    `json:"-"`
}

// repoWithLookupsJSON contains the JSON metadata for the struct [RepoWithLookups]
type repoWithLookupsJSON struct {
	ID                   apijson.Field
	CreatedAt            apijson.Field
	FullName             apijson.Field
	IsArchived           apijson.Field
	IsPublic             apijson.Field
	NumCommits           apijson.Field
	NumDownloads         apijson.Field
	NumLikes             apijson.Field
	NumViews             apijson.Field
	Owner                apijson.Field
	RepoHandle           apijson.Field
	Tags                 apijson.Field
	TenantID             apijson.Field
	UpdatedAt            apijson.Field
	CommitTags           apijson.Field
	CreatedBy            apijson.Field
	Description          apijson.Field
	LastCommitHash       apijson.Field
	LatestCommitManifest apijson.Field
	LikedByAuthUser      apijson.Field
	OriginalRepoFullName apijson.Field
	OriginalRepoID       apijson.Field
	Readme               apijson.Field
	UpstreamRepoFullName apijson.Field
	UpstreamRepoID       apijson.Field
	raw                  string
	ExtraFields          map[string]apijson.Field
}

func (r *RepoWithLookups) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r repoWithLookupsJSON) RawJSON() string {
	return r.raw
}

type RepoDeleteResponse = interface{}

type RepoNewParams struct {
	IsPublic    param.Field[bool]     `json:"is_public,required"`
	RepoHandle  param.Field[string]   `json:"repo_handle,required"`
	Description param.Field[string]   `json:"description"`
	Readme      param.Field[string]   `json:"readme"`
	Tags        param.Field[[]string] `json:"tags"`
}

func (r RepoNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RepoUpdateParams struct {
	Description param.Field[string]   `json:"description"`
	IsArchived  param.Field[bool]     `json:"is_archived"`
	IsPublic    param.Field[bool]     `json:"is_public"`
	Readme      param.Field[string]   `json:"readme"`
	Tags        param.Field[[]string] `json:"tags"`
}

func (r RepoUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RepoListParams struct {
	HasCommits         param.Field[bool]                        `query:"has_commits"`
	IsArchived         param.Field[RepoListParamsIsArchived]    `query:"is_archived"`
	IsPublic           param.Field[RepoListParamsIsPublic]      `query:"is_public"`
	Limit              param.Field[int64]                       `query:"limit"`
	Offset             param.Field[int64]                       `query:"offset"`
	Query              param.Field[string]                      `query:"query"`
	SortDirection      param.Field[RepoListParamsSortDirection] `query:"sort_direction"`
	SortField          param.Field[RepoListParamsSortField]     `query:"sort_field"`
	TagValueID         param.Field[[]string]                    `query:"tag_value_id" format:"uuid"`
	Tags               param.Field[[]string]                    `query:"tags"`
	TenantHandle       param.Field[string]                      `query:"tenant_handle"`
	TenantID           param.Field[string]                      `query:"tenant_id" format:"uuid"`
	UpstreamRepoHandle param.Field[string]                      `query:"upstream_repo_handle"`
	UpstreamRepoOwner  param.Field[string]                      `query:"upstream_repo_owner"`
	WithLatestManifest param.Field[bool]                        `query:"with_latest_manifest"`
}

// URLQuery serializes [RepoListParams]'s query parameters as `url.Values`.
func (r RepoListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type RepoListParamsIsArchived string

const (
	RepoListParamsIsArchivedTrue  RepoListParamsIsArchived = "true"
	RepoListParamsIsArchivedAllow RepoListParamsIsArchived = "allow"
	RepoListParamsIsArchivedFalse RepoListParamsIsArchived = "false"
)

func (r RepoListParamsIsArchived) IsKnown() bool {
	switch r {
	case RepoListParamsIsArchivedTrue, RepoListParamsIsArchivedAllow, RepoListParamsIsArchivedFalse:
		return true
	}
	return false
}

type RepoListParamsIsPublic string

const (
	RepoListParamsIsPublicTrue  RepoListParamsIsPublic = "true"
	RepoListParamsIsPublicFalse RepoListParamsIsPublic = "false"
)

func (r RepoListParamsIsPublic) IsKnown() bool {
	switch r {
	case RepoListParamsIsPublicTrue, RepoListParamsIsPublicFalse:
		return true
	}
	return false
}

type RepoListParamsSortDirection string

const (
	RepoListParamsSortDirectionAsc  RepoListParamsSortDirection = "asc"
	RepoListParamsSortDirectionDesc RepoListParamsSortDirection = "desc"
)

func (r RepoListParamsSortDirection) IsKnown() bool {
	switch r {
	case RepoListParamsSortDirectionAsc, RepoListParamsSortDirectionDesc:
		return true
	}
	return false
}

type RepoListParamsSortField string

const (
	RepoListParamsSortFieldNumLikes     RepoListParamsSortField = "num_likes"
	RepoListParamsSortFieldNumDownloads RepoListParamsSortField = "num_downloads"
	RepoListParamsSortFieldNumViews     RepoListParamsSortField = "num_views"
	RepoListParamsSortFieldUpdatedAt    RepoListParamsSortField = "updated_at"
	RepoListParamsSortFieldRelevance    RepoListParamsSortField = "relevance"
)

func (r RepoListParamsSortField) IsKnown() bool {
	switch r {
	case RepoListParamsSortFieldNumLikes, RepoListParamsSortFieldNumDownloads, RepoListParamsSortFieldNumViews, RepoListParamsSortFieldUpdatedAt, RepoListParamsSortFieldRelevance:
		return true
	}
	return false
}
