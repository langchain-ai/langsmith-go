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
)

// RepoDirectoryService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewRepoDirectoryService] method instead.
type RepoDirectoryService struct {
	Options []option.RequestOption
}

// NewRepoDirectoryService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewRepoDirectoryService(opts ...option.RequestOption) (r *RepoDirectoryService) {
	r = &RepoDirectoryService{}
	r.Options = opts
	return
}

// Resolves the flattened file tree for an agent or skill repository at a specific
// commit, tag, or latest.
func (r *RepoDirectoryService) List(ctx context.Context, owner string, repo string, query RepoDirectoryListParams, opts ...option.RequestOption) (res *RepoDirectoryListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return nil, err
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/platform/hub/repos/%s/%s/directories", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Deletes an agent or skill repository and its owned child file repositories.
func (r *RepoDirectoryService) Delete(ctx context.Context, owner string, repo string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return err
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return err
	}
	path := fmt.Sprintf("v1/platform/hub/repos/%s/%s/directories", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Creates a new directory commit for an agent or skill repository by applying
// file/link create, update, and delete operations.
func (r *RepoDirectoryService) Commit(ctx context.Context, owner string, repo string, body RepoDirectoryCommitParams, opts ...option.RequestOption) (res *RepoDirectoryCommitResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if owner == "" {
		err = errors.New("missing required owner parameter")
		return nil, err
	}
	if repo == "" {
		err = errors.New("missing required repo parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/platform/hub/repos/%s/%s/directories/commits", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type RepoDirectoryListResponse struct {
	CommitHash string                        `json:"commit_hash"`
	CommitID   string                        `json:"commit_id"`
	Files      map[string]interface{}        `json:"files"`
	JSON       repoDirectoryListResponseJSON `json:"-"`
}

// repoDirectoryListResponseJSON contains the JSON metadata for the struct
// [RepoDirectoryListResponse]
type repoDirectoryListResponseJSON struct {
	CommitHash  apijson.Field
	CommitID    apijson.Field
	Files       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RepoDirectoryListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r repoDirectoryListResponseJSON) RawJSON() string {
	return r.raw
}

type RepoDirectoryCommitResponse struct {
	Commit RepoDirectoryCommitResponseCommit `json:"commit"`
	JSON   repoDirectoryCommitResponseJSON   `json:"-"`
}

// repoDirectoryCommitResponseJSON contains the JSON metadata for the struct
// [RepoDirectoryCommitResponse]
type repoDirectoryCommitResponseJSON struct {
	Commit      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RepoDirectoryCommitResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r repoDirectoryCommitResponseJSON) RawJSON() string {
	return r.raw
}

type RepoDirectoryCommitResponseCommit struct {
	ID         string                                `json:"id"`
	CommitHash string                                `json:"commit_hash"`
	CreatedAt  string                                `json:"created_at"`
	JSON       repoDirectoryCommitResponseCommitJSON `json:"-"`
}

// repoDirectoryCommitResponseCommitJSON contains the JSON metadata for the struct
// [RepoDirectoryCommitResponseCommit]
type repoDirectoryCommitResponseCommitJSON struct {
	ID          apijson.Field
	CommitHash  apijson.Field
	CreatedAt   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RepoDirectoryCommitResponseCommit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r repoDirectoryCommitResponseCommitJSON) RawJSON() string {
	return r.raw
}

type RepoDirectoryListParams struct {
	// Commit hash/tag to resolve (defaults to latest)
	Commit param.Field[string] `query:"commit"`
}

// URLQuery serializes [RepoDirectoryListParams]'s query parameters as
// `url.Values`.
func (r RepoDirectoryListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type RepoDirectoryCommitParams struct {
	// Files maps path to an Entry (object = create/update/link, null = delete/unlink).
	Files        param.Field[map[string]interface{}] `json:"files"`
	ParentCommit param.Field[string]                 `json:"parent_commit"`
}

func (r RepoDirectoryCommitParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
