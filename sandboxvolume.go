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

// SandboxVolumeService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxVolumeService] method instead.
type SandboxVolumeService struct {
	Options []option.RequestOption
}

// NewSandboxVolumeService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxVolumeService(opts ...option.RequestOption) (r *SandboxVolumeService) {
	r = &SandboxVolumeService{}
	r.Options = opts
	return
}

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. Create an empty
// volume that can be attached to sandboxes at a mount path. The size is rounded up
// to a whole GiB.
func (r *SandboxVolumeService) New(ctx context.Context, body SandboxVolumeNewParams, opts ...option.RequestOption) (res *VolumeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v2/sandboxes/volumes"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. Get a sandbox volume
// by ID or name, with its current attachments.
func (r *SandboxVolumeService) Get(ctx context.Context, volume string, opts ...option.RequestOption) (res *VolumeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if volume == "" {
		err = errors.New("missing required volume parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/volumes/%s", volume)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. Update a volume's
// description or labels. Each field sent replaces the stored value; omitted fields
// are unchanged. Null, an empty description, or an empty labels object clears it.
// Only the volume's creator, or a role with sandboxes:exec, may update it.
func (r *SandboxVolumeService) Update(ctx context.Context, volume string, body SandboxVolumeUpdateParams, opts ...option.RequestOption) (res *VolumeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if volume == "" {
		err = errors.New("missing required volume parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/volumes/%s", volume)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. List the workspace's
// sandbox volumes with their current attachments.
func (r *SandboxVolumeService) List(ctx context.Context, query SandboxVolumeListParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[VolumeResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v2/sandboxes/volumes"
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

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. List the workspace's
// sandbox volumes with their current attachments.
func (r *SandboxVolumeService) ListAutoPaging(ctx context.Context, query SandboxVolumeListParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[VolumeResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.List(ctx, query, opts...))
}

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. Delete a volume and
// its data. A volume attached to any sandbox cannot be deleted. Storage is
// reclaimed asynchronously. Only the volume's creator, or a role with
// sandboxes:exec, may delete it.
func (r *SandboxVolumeService) Delete(ctx context.Context, volume string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if volume == "" {
		err = errors.New("missing required volume parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/sandboxes/volumes/%s", volume)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// **Beta:** This API is in active development and may change without notice.
// Requires sandbox volumes to be enabled for your workspace. Create a new volume
// holding a copy-on-write copy of an existing one. A volume attached read-write to
// a sandbox that is not stopped cannot be cloned. Only the volume's creator, or a
// role with sandboxes:exec, may clone it.
func (r *SandboxVolumeService) Clone(ctx context.Context, volume string, body SandboxVolumeCloneParams, opts ...option.RequestOption) (res *VolumeResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if volume == "" {
		err = errors.New("missing required volume parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/volumes/%s/clone", volume)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type VolumeResponse struct {
	ID             string                     `json:"id"`
	Attachments    []VolumeResponseAttachment `json:"attachments"`
	CreatedAt      time.Time                  `json:"created_at" format:"date-time"`
	CreatedBy      string                     `json:"created_by"`
	Description    string                     `json:"description"`
	Labels         map[string]string          `json:"labels"`
	Name           string                     `json:"name"`
	SizeBytes      int64                      `json:"size_bytes"`
	SourceVolumeID string                     `json:"source_volume_id"`
	Status         string                     `json:"status"`
	StatusMessage  string                     `json:"status_message"`
	UpdatedAt      time.Time                  `json:"updated_at" format:"date-time"`
	UpdatedBy      string                     `json:"updated_by"`
	JSON           volumeResponseJSON         `json:"-"`
}

// volumeResponseJSON contains the JSON metadata for the struct [VolumeResponse]
type volumeResponseJSON struct {
	ID             apijson.Field
	Attachments    apijson.Field
	CreatedAt      apijson.Field
	CreatedBy      apijson.Field
	Description    apijson.Field
	Labels         apijson.Field
	Name           apijson.Field
	SizeBytes      apijson.Field
	SourceVolumeID apijson.Field
	Status         apijson.Field
	StatusMessage  apijson.Field
	UpdatedAt      apijson.Field
	UpdatedBy      apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *VolumeResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r volumeResponseJSON) RawJSON() string {
	return r.raw
}

type VolumeResponseAttachment struct {
	MountPath   string                       `json:"mount_path"`
	ReadOnly    bool                         `json:"read_only"`
	SandboxID   string                       `json:"sandbox_id"`
	SandboxName string                       `json:"sandbox_name"`
	JSON        volumeResponseAttachmentJSON `json:"-"`
}

// volumeResponseAttachmentJSON contains the JSON metadata for the struct
// [VolumeResponseAttachment]
type volumeResponseAttachmentJSON struct {
	MountPath   apijson.Field
	ReadOnly    apijson.Field
	SandboxID   apijson.Field
	SandboxName apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *VolumeResponseAttachment) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r volumeResponseAttachmentJSON) RawJSON() string {
	return r.raw
}

type SandboxVolumeNewParams struct {
	// Name of the volume, unique within the workspace.
	Name param.Field[string] `json:"name" api:"required"`
	// Capacity in bytes, rounded up to a whole GiB.
	SizeBytes param.Field[int64] `json:"size_bytes" api:"required"`
	// Free-text description, at most 1024 characters.
	Description param.Field[string] `json:"description"`
	// Labels to filter and identify the volume by. Keys are lowercase letters, digits,
	// dots and hyphens.
	Labels param.Field[map[string]string] `json:"labels"`
}

func (r SandboxVolumeNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxVolumeUpdateParams struct {
	// Description replaces the description; null or an empty string clears it.
	Description param.Field[string] `json:"description"`
	// Labels replaces every label; null or an empty object clears them.
	Labels param.Field[map[string]string] `json:"labels"`
}

func (r SandboxVolumeUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxVolumeListParams struct {
	// Opaque pagination cursor from a prior response's next_cursor
	Cursor param.Field[string] `query:"cursor"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Number of results per page
	PageSize param.Field[int64] `query:"page_size"`
}

// URLQuery serializes [SandboxVolumeListParams]'s query parameters as
// `url.Values`.
func (r SandboxVolumeListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SandboxVolumeCloneParams struct {
	// Name of the new volume, unique within the workspace.
	Name param.Field[string] `json:"name" api:"required"`
	// Free-text description of the new volume, at most 1024 characters.
	Description param.Field[string] `json:"description"`
	// Labels for the new volume. Keys are lowercase letters, digits, dots and hyphens.
	Labels param.Field[map[string]string] `json:"labels"`
}

func (r SandboxVolumeCloneParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
