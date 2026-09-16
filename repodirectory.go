// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/tidwall/gjson"
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
	path := fmt.Sprintf("api/v1/platform/hub/repos/%s/%s/directories", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Deletes an agent or skill repository and its owned child file repositories.
func (r *RepoDirectoryService) Delete(ctx context.Context, owner string, repo string, body RepoDirectoryDeleteParams, opts ...option.RequestOption) (err error) {
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
	path := fmt.Sprintf("api/v1/platform/hub/repos/%s/%s/directories", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return err
}

// Creates a new directory commit for an agent or skill repository by applying
// file/link create, update, and delete operations. Linked directories default to
// the LATEST selector; use COMMIT to pin one commit. The legacy commit_id write
// field is deprecated and resolves as LATEST.
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
	path := fmt.Sprintf("api/v1/platform/hub/repos/%s/%s/directories/commits", owner, repo)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type DirectorySelector struct {
	Type     DirectorySelectorType `json:"type" api:"required"`
	CommitID string                `json:"commit_id" format:"uuid"`
	JSON     directorySelectorJSON `json:"-"`
	union    DirectorySelectorUnion
}

// directorySelectorJSON contains the JSON metadata for the struct
// [DirectorySelector]
type directorySelectorJSON struct {
	Type        apijson.Field
	CommitID    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r directorySelectorJSON) RawJSON() string {
	return r.raw
}

func (r *DirectorySelector) UnmarshalJSON(data []byte) (err error) {
	*r = DirectorySelector{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [DirectorySelectorUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are
// [DirectorySelectorDirectoryLatestSelector],
// [DirectorySelectorDirectoryCommitSelector].
func (r DirectorySelector) AsUnion() DirectorySelectorUnion {
	return r.union
}

// Union satisfied by [DirectorySelectorDirectoryLatestSelector] or
// [DirectorySelectorDirectoryCommitSelector].
type DirectorySelectorUnion interface {
	implementsDirectorySelector()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*DirectorySelectorUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(DirectorySelectorDirectoryLatestSelector{}),
			DiscriminatorValue: "LATEST",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(DirectorySelectorDirectoryCommitSelector{}),
			DiscriminatorValue: "COMMIT",
		},
	)
}

type DirectorySelectorDirectoryLatestSelector struct {
	Type DirectorySelectorDirectoryLatestSelectorType `json:"type" api:"required"`
	JSON directorySelectorDirectoryLatestSelectorJSON `json:"-"`
}

// directorySelectorDirectoryLatestSelectorJSON contains the JSON metadata for the
// struct [DirectorySelectorDirectoryLatestSelector]
type directorySelectorDirectoryLatestSelectorJSON struct {
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DirectorySelectorDirectoryLatestSelector) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r directorySelectorDirectoryLatestSelectorJSON) RawJSON() string {
	return r.raw
}

func (r DirectorySelectorDirectoryLatestSelector) implementsDirectorySelector() {}

type DirectorySelectorDirectoryLatestSelectorType string

const (
	DirectorySelectorDirectoryLatestSelectorTypeLatest DirectorySelectorDirectoryLatestSelectorType = "LATEST"
)

func (r DirectorySelectorDirectoryLatestSelectorType) IsKnown() bool {
	switch r {
	case DirectorySelectorDirectoryLatestSelectorTypeLatest:
		return true
	}
	return false
}

type DirectorySelectorDirectoryCommitSelector struct {
	CommitID string                                       `json:"commit_id" api:"required" format:"uuid"`
	Type     DirectorySelectorDirectoryCommitSelectorType `json:"type" api:"required"`
	JSON     directorySelectorDirectoryCommitSelectorJSON `json:"-"`
}

// directorySelectorDirectoryCommitSelectorJSON contains the JSON metadata for the
// struct [DirectorySelectorDirectoryCommitSelector]
type directorySelectorDirectoryCommitSelectorJSON struct {
	CommitID    apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DirectorySelectorDirectoryCommitSelector) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r directorySelectorDirectoryCommitSelectorJSON) RawJSON() string {
	return r.raw
}

func (r DirectorySelectorDirectoryCommitSelector) implementsDirectorySelector() {}

type DirectorySelectorDirectoryCommitSelectorType string

const (
	DirectorySelectorDirectoryCommitSelectorTypeCommit DirectorySelectorDirectoryCommitSelectorType = "COMMIT"
)

func (r DirectorySelectorDirectoryCommitSelectorType) IsKnown() bool {
	switch r {
	case DirectorySelectorDirectoryCommitSelectorTypeCommit:
		return true
	}
	return false
}

type DirectorySelectorType string

const (
	DirectorySelectorTypeLatest DirectorySelectorType = "LATEST"
	DirectorySelectorTypeCommit DirectorySelectorType = "COMMIT"
)

func (r DirectorySelectorType) IsKnown() bool {
	switch r {
	case DirectorySelectorTypeLatest, DirectorySelectorTypeCommit:
		return true
	}
	return false
}

type DirectorySelectorParam struct {
	Type     param.Field[DirectorySelectorType] `json:"type" api:"required"`
	CommitID param.Field[string]                `json:"commit_id" format:"uuid"`
}

func (r DirectorySelectorParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r DirectorySelectorParam) implementsDirectorySelectorUnionParam() {}

// Satisfied by [DirectorySelectorDirectoryLatestSelectorParam],
// [DirectorySelectorDirectoryCommitSelectorParam], [DirectorySelectorParam].
type DirectorySelectorUnionParam interface {
	implementsDirectorySelectorUnionParam()
}

type DirectorySelectorDirectoryLatestSelectorParam struct {
	Type param.Field[DirectorySelectorDirectoryLatestSelectorType] `json:"type" api:"required"`
}

func (r DirectorySelectorDirectoryLatestSelectorParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r DirectorySelectorDirectoryLatestSelectorParam) implementsDirectorySelectorUnionParam() {}

type DirectorySelectorDirectoryCommitSelectorParam struct {
	CommitID param.Field[string]                                       `json:"commit_id" api:"required" format:"uuid"`
	Type     param.Field[DirectorySelectorDirectoryCommitSelectorType] `json:"type" api:"required"`
}

func (r DirectorySelectorDirectoryCommitSelectorParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r DirectorySelectorDirectoryCommitSelectorParam) implementsDirectorySelectorUnionParam() {}

type FileEntry struct {
	Content string        `json:"content" api:"required"`
	Type    FileEntryType `json:"type" api:"required"`
	JSON    fileEntryJSON `json:"-"`
}

// fileEntryJSON contains the JSON metadata for the struct [FileEntry]
type fileEntryJSON struct {
	Content     apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FileEntry) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fileEntryJSON) RawJSON() string {
	return r.raw
}

func (r FileEntry) implementsRepoDirectoryListResponseFile() {}

type FileEntryType string

const (
	FileEntryTypeFile FileEntryType = "file"
)

func (r FileEntryType) IsKnown() bool {
	switch r {
	case FileEntryTypeFile:
		return true
	}
	return false
}

type FileEntryParam struct {
	Content param.Field[string]        `json:"content" api:"required"`
	Type    param.Field[FileEntryType] `json:"type" api:"required"`
}

func (r FileEntryParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r FileEntryParam) implementsRepoDirectoryCommitParamsFilesUnion() {}

type RepoDirectoryListResponse struct {
	CommitHash string                                   `json:"commit_hash" api:"required"`
	CommitID   string                                   `json:"commit_id" api:"required" format:"uuid"`
	Files      map[string]RepoDirectoryListResponseFile `json:"files" api:"required"`
	JSON       repoDirectoryListResponseJSON            `json:"-"`
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

type RepoDirectoryListResponseFile struct {
	Type       RepoDirectoryListResponseFilesType `json:"type" api:"required"`
	CommitHash string                             `json:"commit_hash"`
	Content    string                             `json:"content"`
	Owner      string                             `json:"owner"`
	RepoHandle string                             `json:"repo_handle"`
	// The authored selection policy for this linked directory.
	Selector DirectorySelector                 `json:"selector"`
	JSON     repoDirectoryListResponseFileJSON `json:"-"`
	union    RepoDirectoryListResponseFilesUnion
}

// repoDirectoryListResponseFileJSON contains the JSON metadata for the struct
// [RepoDirectoryListResponseFile]
type repoDirectoryListResponseFileJSON struct {
	Type        apijson.Field
	CommitHash  apijson.Field
	Content     apijson.Field
	Owner       apijson.Field
	RepoHandle  apijson.Field
	Selector    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r repoDirectoryListResponseFileJSON) RawJSON() string {
	return r.raw
}

func (r *RepoDirectoryListResponseFile) UnmarshalJSON(data []byte) (err error) {
	*r = RepoDirectoryListResponseFile{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [RepoDirectoryListResponseFilesUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are [FileEntry],
// [RepoDirectoryListResponseFilesDirectoryAgentEntryOutput],
// [RepoDirectoryListResponseFilesDirectorySkillEntryOutput].
func (r RepoDirectoryListResponseFile) AsUnion() RepoDirectoryListResponseFilesUnion {
	return r.union
}

// Union satisfied by [FileEntry],
// [RepoDirectoryListResponseFilesDirectoryAgentEntryOutput] or
// [RepoDirectoryListResponseFilesDirectorySkillEntryOutput].
type RepoDirectoryListResponseFilesUnion interface {
	implementsRepoDirectoryListResponseFile()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*RepoDirectoryListResponseFilesUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(FileEntry{}),
			DiscriminatorValue: "file",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RepoDirectoryListResponseFilesDirectoryAgentEntryOutput{}),
			DiscriminatorValue: "agent",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(RepoDirectoryListResponseFilesDirectorySkillEntryOutput{}),
			DiscriminatorValue: "skill",
		},
	)
}

type RepoDirectoryListResponseFilesDirectoryAgentEntryOutput struct {
	CommitHash string                                                      `json:"commit_hash" api:"required"`
	Owner      string                                                      `json:"owner" api:"required"`
	RepoHandle string                                                      `json:"repo_handle" api:"required"`
	Type       RepoDirectoryListResponseFilesDirectoryAgentEntryOutputType `json:"type" api:"required"`
	// The authored selection policy for this linked directory.
	Selector DirectorySelector                                           `json:"selector"`
	JSON     repoDirectoryListResponseFilesDirectoryAgentEntryOutputJSON `json:"-"`
}

// repoDirectoryListResponseFilesDirectoryAgentEntryOutputJSON contains the JSON
// metadata for the struct
// [RepoDirectoryListResponseFilesDirectoryAgentEntryOutput]
type repoDirectoryListResponseFilesDirectoryAgentEntryOutputJSON struct {
	CommitHash  apijson.Field
	Owner       apijson.Field
	RepoHandle  apijson.Field
	Type        apijson.Field
	Selector    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RepoDirectoryListResponseFilesDirectoryAgentEntryOutput) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r repoDirectoryListResponseFilesDirectoryAgentEntryOutputJSON) RawJSON() string {
	return r.raw
}

func (r RepoDirectoryListResponseFilesDirectoryAgentEntryOutput) implementsRepoDirectoryListResponseFile() {
}

type RepoDirectoryListResponseFilesDirectoryAgentEntryOutputType string

const (
	RepoDirectoryListResponseFilesDirectoryAgentEntryOutputTypeAgent RepoDirectoryListResponseFilesDirectoryAgentEntryOutputType = "agent"
)

func (r RepoDirectoryListResponseFilesDirectoryAgentEntryOutputType) IsKnown() bool {
	switch r {
	case RepoDirectoryListResponseFilesDirectoryAgentEntryOutputTypeAgent:
		return true
	}
	return false
}

type RepoDirectoryListResponseFilesDirectorySkillEntryOutput struct {
	CommitHash string                                                      `json:"commit_hash" api:"required"`
	Owner      string                                                      `json:"owner" api:"required"`
	RepoHandle string                                                      `json:"repo_handle" api:"required"`
	Type       RepoDirectoryListResponseFilesDirectorySkillEntryOutputType `json:"type" api:"required"`
	// The authored selection policy for this linked directory.
	Selector DirectorySelector                                           `json:"selector"`
	JSON     repoDirectoryListResponseFilesDirectorySkillEntryOutputJSON `json:"-"`
}

// repoDirectoryListResponseFilesDirectorySkillEntryOutputJSON contains the JSON
// metadata for the struct
// [RepoDirectoryListResponseFilesDirectorySkillEntryOutput]
type repoDirectoryListResponseFilesDirectorySkillEntryOutputJSON struct {
	CommitHash  apijson.Field
	Owner       apijson.Field
	RepoHandle  apijson.Field
	Type        apijson.Field
	Selector    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RepoDirectoryListResponseFilesDirectorySkillEntryOutput) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r repoDirectoryListResponseFilesDirectorySkillEntryOutputJSON) RawJSON() string {
	return r.raw
}

func (r RepoDirectoryListResponseFilesDirectorySkillEntryOutput) implementsRepoDirectoryListResponseFile() {
}

type RepoDirectoryListResponseFilesDirectorySkillEntryOutputType string

const (
	RepoDirectoryListResponseFilesDirectorySkillEntryOutputTypeSkill RepoDirectoryListResponseFilesDirectorySkillEntryOutputType = "skill"
)

func (r RepoDirectoryListResponseFilesDirectorySkillEntryOutputType) IsKnown() bool {
	switch r {
	case RepoDirectoryListResponseFilesDirectorySkillEntryOutputTypeSkill:
		return true
	}
	return false
}

type RepoDirectoryListResponseFilesType string

const (
	RepoDirectoryListResponseFilesTypeFile  RepoDirectoryListResponseFilesType = "file"
	RepoDirectoryListResponseFilesTypeAgent RepoDirectoryListResponseFilesType = "agent"
	RepoDirectoryListResponseFilesTypeSkill RepoDirectoryListResponseFilesType = "skill"
)

func (r RepoDirectoryListResponseFilesType) IsKnown() bool {
	switch r {
	case RepoDirectoryListResponseFilesTypeFile, RepoDirectoryListResponseFilesTypeAgent, RepoDirectoryListResponseFilesTypeSkill:
		return true
	}
	return false
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type RepoDirectoryDeleteParams struct {
	// Repository type to delete; a different type is treated as not found
	RepoType param.Field[RepoDirectoryDeleteParamsRepoType] `query:"repo_type"`
}

// URLQuery serializes [RepoDirectoryDeleteParams]'s query parameters as
// `url.Values`.
func (r RepoDirectoryDeleteParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Repository type to delete; a different type is treated as not found
type RepoDirectoryDeleteParamsRepoType string

const (
	RepoDirectoryDeleteParamsRepoTypeAgent RepoDirectoryDeleteParamsRepoType = "agent"
	RepoDirectoryDeleteParamsRepoTypeSkill RepoDirectoryDeleteParamsRepoType = "skill"
)

func (r RepoDirectoryDeleteParamsRepoType) IsKnown() bool {
	switch r {
	case RepoDirectoryDeleteParamsRepoTypeAgent, RepoDirectoryDeleteParamsRepoTypeSkill:
		return true
	}
	return false
}

type RepoDirectoryCommitParams struct {
	// Paths to create, update, link, delete, or unlink. Use null to delete or unlink
	// an existing path.
	Files        param.Field[map[string]RepoDirectoryCommitParamsFilesUnion] `json:"files"`
	ParentCommit param.Field[string]                                         `json:"parent_commit"`
	// SkipWebhooks suppresses Context Hub commit webhooks for this commit.
	SkipWebhooks param.Field[bool] `json:"skip_webhooks"`
}

func (r RepoDirectoryCommitParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RepoDirectoryCommitParamsFiles struct {
	Type param.Field[RepoDirectoryCommitParamsFilesType] `json:"type" api:"required"`
	// Deprecated write input. It is accepted for compatibility but ignored for
	// selection, so the link resolves as LATEST. Omit it for LATEST or replace it with
	// selector {"type": "COMMIT", "commit_id": "<uuid>"} to pin a commit. commit_id
	// and selector are mutually exclusive.
	//
	// Deprecated: deprecated
	CommitID   param.Field[string] `json:"commit_id" format:"uuid"`
	Content    param.Field[string] `json:"content"`
	RepoHandle param.Field[string] `json:"repo_handle"`
	// How to select the linked commit. Omit this field to use LATEST.
	Selector param.Field[DirectorySelectorUnionParam] `json:"selector"`
}

func (r RepoDirectoryCommitParamsFiles) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r RepoDirectoryCommitParamsFiles) implementsRepoDirectoryCommitParamsFilesUnion() {}

// Satisfied by [FileEntryParam],
// [RepoDirectoryCommitParamsFilesDirectoryAgentEntryInput],
// [RepoDirectoryCommitParamsFilesDirectorySkillEntryInput],
// [RepoDirectoryCommitParamsFiles].
type RepoDirectoryCommitParamsFilesUnion interface {
	implementsRepoDirectoryCommitParamsFilesUnion()
}

type RepoDirectoryCommitParamsFilesDirectoryAgentEntryInput struct {
	RepoHandle param.Field[string]                                                     `json:"repo_handle" api:"required"`
	Type       param.Field[RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputType] `json:"type" api:"required"`
	// Deprecated write input. It is accepted for compatibility but ignored for
	// selection, so the link resolves as LATEST. Omit it for LATEST or replace it with
	// selector {"type": "COMMIT", "commit_id": "<uuid>"} to pin a commit. commit_id
	// and selector are mutually exclusive.
	//
	// Deprecated: deprecated
	CommitID param.Field[string] `json:"commit_id" format:"uuid"`
	// How to select the linked commit. Omit this field to use LATEST.
	Selector param.Field[DirectorySelectorUnionParam] `json:"selector"`
}

func (r RepoDirectoryCommitParamsFilesDirectoryAgentEntryInput) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r RepoDirectoryCommitParamsFilesDirectoryAgentEntryInput) implementsRepoDirectoryCommitParamsFilesUnion() {
}

type RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputType string

const (
	RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputTypeAgent RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputType = "agent"
)

func (r RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputType) IsKnown() bool {
	switch r {
	case RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputTypeAgent:
		return true
	}
	return false
}

type RepoDirectoryCommitParamsFilesDirectorySkillEntryInput struct {
	RepoHandle param.Field[string]                                                     `json:"repo_handle" api:"required"`
	Type       param.Field[RepoDirectoryCommitParamsFilesDirectorySkillEntryInputType] `json:"type" api:"required"`
	// Deprecated write input. It is accepted for compatibility but ignored for
	// selection, so the link resolves as LATEST. Omit it for LATEST or replace it with
	// selector {"type": "COMMIT", "commit_id": "<uuid>"} to pin a commit. commit_id
	// and selector are mutually exclusive.
	//
	// Deprecated: deprecated
	CommitID param.Field[string] `json:"commit_id" format:"uuid"`
	// How to select the linked commit. Omit this field to use LATEST.
	Selector param.Field[DirectorySelectorUnionParam] `json:"selector"`
}

func (r RepoDirectoryCommitParamsFilesDirectorySkillEntryInput) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r RepoDirectoryCommitParamsFilesDirectorySkillEntryInput) implementsRepoDirectoryCommitParamsFilesUnion() {
}

type RepoDirectoryCommitParamsFilesDirectorySkillEntryInputType string

const (
	RepoDirectoryCommitParamsFilesDirectorySkillEntryInputTypeSkill RepoDirectoryCommitParamsFilesDirectorySkillEntryInputType = "skill"
)

func (r RepoDirectoryCommitParamsFilesDirectorySkillEntryInputType) IsKnown() bool {
	switch r {
	case RepoDirectoryCommitParamsFilesDirectorySkillEntryInputTypeSkill:
		return true
	}
	return false
}

type RepoDirectoryCommitParamsFilesType string

const (
	RepoDirectoryCommitParamsFilesTypeFile  RepoDirectoryCommitParamsFilesType = "file"
	RepoDirectoryCommitParamsFilesTypeAgent RepoDirectoryCommitParamsFilesType = "agent"
	RepoDirectoryCommitParamsFilesTypeSkill RepoDirectoryCommitParamsFilesType = "skill"
)

func (r RepoDirectoryCommitParamsFilesType) IsKnown() bool {
	switch r {
	case RepoDirectoryCommitParamsFilesTypeFile, RepoDirectoryCommitParamsFilesTypeAgent, RepoDirectoryCommitParamsFilesTypeSkill:
		return true
	}
	return false
}
