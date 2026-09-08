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

// SandboxSnapshotService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxSnapshotService] method instead.
type SandboxSnapshotService struct {
	Options []option.RequestOption
}

// NewSandboxSnapshotService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxSnapshotService(opts ...option.RequestOption) (r *SandboxSnapshotService) {
	r = &SandboxSnapshotService{}
	r.Options = opts
	return
}

// Create a snapshot from a Docker image (async build).
func (r *SandboxSnapshotService) New(ctx context.Context, body SandboxSnapshotNewParams, opts ...option.RequestOption) (res *SnapshotResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v2/sandboxes/snapshots"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a sandbox snapshot by ID or by a Docker-style reference. A bare name means
// name:latest, falling back to the newest ready untagged snapshot of that name. To
// list the tags under a name, use /api/v2/sandboxes/snapshots-by-name/{name}.
func (r *SandboxSnapshotService) Get(ctx context.Context, snapshotID string, opts ...option.RequestOption) (res *SnapshotResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if snapshotID == "" {
		err = errors.New("missing required snapshot_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/snapshots/%s", snapshotID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List sandbox snapshots for the authenticated tenant, with optional filtering,
// sorting, and pagination. Page with page_size and cursor: replay the response's
// next_cursor until it comes back null, which is the only signal that no pages
// remain. Cursors are opaque and only valid on this endpoint; do not parse or
// construct one.
func (r *SandboxSnapshotService) List(ctx context.Context, query SandboxSnapshotListParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[SnapshotResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v2/sandboxes/snapshots"
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

// List sandbox snapshots for the authenticated tenant, with optional filtering,
// sorting, and pagination. Page with page_size and cursor: replay the response's
// next_cursor until it comes back null, which is the only signal that no pages
// remain. Cursors are opaque and only valid on this endpoint; do not parse or
// construct one.
func (r *SandboxSnapshotService) ListAutoPaging(ctx context.Context, query SandboxSnapshotListParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[SnapshotResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a snapshot by ID or by a Docker-style name[:tag] reference. The
// underlying storage is reclaimed asynchronously.
func (r *SandboxSnapshotService) Delete(ctx context.Context, snapshotID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if snapshotID == "" {
		err = errors.New("missing required snapshot_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/sandboxes/snapshots/%s", snapshotID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Get a snapshot name and every tag under it, with the snapshot each tag resolves
// to. To fetch one snapshot, use /api/v2/sandboxes/snapshots/{snapshot_id}.
func (r *SandboxSnapshotService) GetByName(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxSnapshotGetByNameResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/snapshots-by-name/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type SandboxSnapshotGetByNameResponse struct {
	Name string                                `json:"name"`
	Tags []SandboxSnapshotGetByNameResponseTag `json:"tags"`
	JSON sandboxSnapshotGetByNameResponseJSON  `json:"-"`
}

// sandboxSnapshotGetByNameResponseJSON contains the JSON metadata for the struct
// [SandboxSnapshotGetByNameResponse]
type sandboxSnapshotGetByNameResponseJSON struct {
	Name        apijson.Field
	Tags        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxSnapshotGetByNameResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxSnapshotGetByNameResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxSnapshotGetByNameResponseTag struct {
	SnapshotID string                                  `json:"snapshot_id"`
	Tag        string                                  `json:"tag"`
	JSON       sandboxSnapshotGetByNameResponseTagJSON `json:"-"`
}

// sandboxSnapshotGetByNameResponseTagJSON contains the JSON metadata for the
// struct [SandboxSnapshotGetByNameResponseTag]
type sandboxSnapshotGetByNameResponseTagJSON struct {
	SnapshotID  apijson.Field
	Tag         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxSnapshotGetByNameResponseTag) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxSnapshotGetByNameResponseTagJSON) RawJSON() string {
	return r.raw
}

type SandboxSnapshotNewParams struct {
	DockerImage     param.Field[string] `json:"docker_image" api:"required"`
	FsCapacityBytes param.Field[int64]  `json:"fs_capacity_bytes" api:"required"`
	Name            param.Field[string] `json:"name" api:"required"`
	// Description says what this snapshot's image can do, so a caller can hand it to
	// an agent as a capability summary. At most 1024 characters.
	Description param.Field[string] `json:"description"`
	// Labels seed the snapshot's labels, overriding any label of the same key derived
	// from the Docker image.
	Labels     param.Field[map[string]string] `json:"labels"`
	RegistryID param.Field[string]            `json:"registry_id"`
	// RunConfig overrides the runtime configuration taken from the Docker image. Every
	// sandbox created from the snapshot runs as the image's USER, in its WORKDIR, with
	// its ENV beneath the sandbox's own env_vars; user and work_dir given here replace
	// the image's, and env_vars merge over it.
	RunConfig param.Field[SandboxSnapshotNewParamsRunConfig] `json:"run_config"`
	// mutable Docker-style tag; defaults to "latest"
	Tag param.Field[string] `json:"tag"`
}

func (r SandboxSnapshotNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// RunConfig overrides the runtime configuration taken from the Docker image. Every
// sandbox created from the snapshot runs as the image's USER, in its WORKDIR, with
// its ENV beneath the sandbox's own env_vars; user and work_dir given here replace
// the image's, and env_vars merge over it.
type SandboxSnapshotNewParamsRunConfig struct {
	EnvVars param.Field[map[string]string] `json:"env_vars"`
	User    param.Field[string]            `json:"user"`
	WorkDir param.Field[string]            `json:"work_dir"`
}

func (r SandboxSnapshotNewParamsRunConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxSnapshotListParams struct {
	// Filter by creator identity. Only 'me' is supported.
	CreatedBy param.Field[string] `query:"created_by"`
	// Opaque pagination cursor from a prior response's next_cursor
	Cursor param.Field[string] `query:"cursor"`
	// Filter by label. Repeatable; all must match. Use 'key' to match on key presence
	// or 'key=value' for equality.
	Label param.Field[[]string] `query:"label"`
	// Deprecated: use page_size. Maximum number of results
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Deprecated: use cursor. Pagination offset
	Offset param.Field[int64] `query:"offset"`
	// Number of results per page
	PageSize param.Field[int64] `query:"page_size"`
	// Sort column (name, status, created_at)
	SortBy param.Field[string] `query:"sort_by"`
	// Deprecated: use sort_order. Sort direction (asc, desc)
	SortDirection param.Field[string] `query:"sort_direction"`
	// Sort direction (asc, desc)
	SortOrder param.Field[string] `query:"sort_order"`
	// Filter by status (building, ready, failed, deleting)
	Status param.Field[string] `query:"status"`
}

// URLQuery serializes [SandboxSnapshotListParams]'s query parameters as
// `url.Values`.
func (r SandboxSnapshotListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
