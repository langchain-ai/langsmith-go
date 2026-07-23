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
	path := "v2/sandboxes/snapshots"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a sandbox snapshot by ID.
func (r *SandboxSnapshotService) Get(ctx context.Context, snapshotID string, opts ...option.RequestOption) (res *SnapshotResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if snapshotID == "" {
		err = errors.New("missing required snapshot_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/snapshots/%s", snapshotID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// List sandbox snapshots for the authenticated tenant, with optional filtering,
// sorting, and pagination.
func (r *SandboxSnapshotService) List(ctx context.Context, query SandboxSnapshotListParams, opts ...option.RequestOption) (res *SnapshotListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/snapshots"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete a snapshot by ID. The underlying storage is reclaimed asynchronously.
func (r *SandboxSnapshotService) Delete(ctx context.Context, snapshotID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if snapshotID == "" {
		err = errors.New("missing required snapshot_id parameter")
		return err
	}
	path := fmt.Sprintf("v2/sandboxes/snapshots/%s", snapshotID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type SandboxSnapshotNewParams struct {
	DockerImage     param.Field[string] `json:"docker_image" api:"required"`
	FsCapacityBytes param.Field[int64]  `json:"fs_capacity_bytes" api:"required"`
	Name            param.Field[string] `json:"name" api:"required"`
	// Labels seed the snapshot's labels, overriding any label of the same key derived
	// from the Docker image.
	Labels     param.Field[map[string]string] `json:"labels"`
	RegistryID param.Field[string]            `json:"registry_id"`
}

func (r SandboxSnapshotNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxSnapshotListParams struct {
	// Filter by creator identity. Only 'me' is supported.
	CreatedBy param.Field[string] `query:"created_by"`
	// Filter by label. Repeatable; all must match. Use 'key' to match on key presence
	// or 'key=value' for equality.
	Label param.Field[[]string] `query:"label"`
	// Maximum number of results
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Pagination offset
	Offset param.Field[int64] `query:"offset"`
	// Sort column (name, status, created_at)
	SortBy param.Field[string] `query:"sort_by"`
	// Sort direction (asc, desc)
	SortDirection param.Field[string] `query:"sort_direction"`
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
