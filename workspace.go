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
)

// WorkspaceService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkspaceService] method instead.
type WorkspaceService struct {
	Options []option.RequestOption
}

// NewWorkspaceService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewWorkspaceService(opts ...option.RequestOption) (r *WorkspaceService) {
	r = &WorkspaceService{}
	r.Options = opts
	return
}

// Create a new workspace.
func (r *WorkspaceService) New(ctx context.Context, body WorkspaceNewParams, opts ...option.RequestOption) (res *WorkspaceNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Update a workspace.
func (r *WorkspaceService) Update(ctx context.Context, workspaceID string, body WorkspaceUpdateParams, opts ...option.RequestOption) (res *WorkspaceUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/%s", workspaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Get all workspaces visible to this auth in the current org. Does not create a
// new workspace/org.
func (r *WorkspaceService) List(ctx context.Context, query WorkspaceListParams, opts ...option.RequestOption) (res *[]WorkspaceListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete Workspace
func (r *WorkspaceService) Delete(ctx context.Context, workspaceID string, opts ...option.RequestOption) (res *WorkspaceDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if workspaceID == "" {
		err = errors.New("missing required workspace_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/%s", workspaceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Tenant schema.
type WorkspaceNewResponse struct {
	ID             string                   `json:"id" api:"required" format:"uuid"`
	CreatedAt      time.Time                `json:"created_at" api:"required" format:"date-time"`
	DisplayName    string                   `json:"display_name" api:"required"`
	IsDeleted      bool                     `json:"is_deleted" api:"required"`
	IsPersonal     bool                     `json:"is_personal" api:"required"`
	DataPlaneURL   string                   `json:"data_plane_url" api:"nullable"`
	OrganizationID string                   `json:"organization_id" api:"nullable" format:"uuid"`
	TenantHandle   string                   `json:"tenant_handle" api:"nullable"`
	JSON           workspaceNewResponseJSON `json:"-"`
}

// workspaceNewResponseJSON contains the JSON metadata for the struct
// [WorkspaceNewResponse]
type workspaceNewResponseJSON struct {
	ID             apijson.Field
	CreatedAt      apijson.Field
	DisplayName    apijson.Field
	IsDeleted      apijson.Field
	IsPersonal     apijson.Field
	DataPlaneURL   apijson.Field
	OrganizationID apijson.Field
	TenantHandle   apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *WorkspaceNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r workspaceNewResponseJSON) RawJSON() string {
	return r.raw
}

// Tenant schema.
type WorkspaceUpdateResponse struct {
	ID             string                      `json:"id" api:"required" format:"uuid"`
	CreatedAt      time.Time                   `json:"created_at" api:"required" format:"date-time"`
	DisplayName    string                      `json:"display_name" api:"required"`
	IsDeleted      bool                        `json:"is_deleted" api:"required"`
	IsPersonal     bool                        `json:"is_personal" api:"required"`
	DataPlaneURL   string                      `json:"data_plane_url" api:"nullable"`
	OrganizationID string                      `json:"organization_id" api:"nullable" format:"uuid"`
	TenantHandle   string                      `json:"tenant_handle" api:"nullable"`
	JSON           workspaceUpdateResponseJSON `json:"-"`
}

// workspaceUpdateResponseJSON contains the JSON metadata for the struct
// [WorkspaceUpdateResponse]
type workspaceUpdateResponseJSON struct {
	ID             apijson.Field
	CreatedAt      apijson.Field
	DisplayName    apijson.Field
	IsDeleted      apijson.Field
	IsPersonal     apijson.Field
	DataPlaneURL   apijson.Field
	OrganizationID apijson.Field
	TenantHandle   apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *WorkspaceUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r workspaceUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type WorkspaceListResponse struct {
	ID             string    `json:"id" api:"required" format:"uuid"`
	CreatedAt      time.Time `json:"created_at" api:"required" format:"date-time"`
	DisplayName    string    `json:"display_name" api:"required"`
	IsDeleted      bool      `json:"is_deleted" api:"required"`
	IsPersonal     bool      `json:"is_personal" api:"required"`
	DataPlaneURL   string    `json:"data_plane_url" api:"nullable"`
	OrganizationID string    `json:"organization_id" api:"nullable" format:"uuid"`
	Permissions    []string  `json:"permissions" api:"nullable"`
	// Deprecated: deprecated
	ReadOnly     bool                      `json:"read_only"`
	RoleID       string                    `json:"role_id" api:"nullable" format:"uuid"`
	RoleName     string                    `json:"role_name" api:"nullable"`
	TenantHandle string                    `json:"tenant_handle" api:"nullable"`
	JSON         workspaceListResponseJSON `json:"-"`
}

// workspaceListResponseJSON contains the JSON metadata for the struct
// [WorkspaceListResponse]
type workspaceListResponseJSON struct {
	ID             apijson.Field
	CreatedAt      apijson.Field
	DisplayName    apijson.Field
	IsDeleted      apijson.Field
	IsPersonal     apijson.Field
	DataPlaneURL   apijson.Field
	OrganizationID apijson.Field
	Permissions    apijson.Field
	ReadOnly       apijson.Field
	RoleID         apijson.Field
	RoleName       apijson.Field
	TenantHandle   apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *WorkspaceListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r workspaceListResponseJSON) RawJSON() string {
	return r.raw
}

type WorkspaceDeleteResponse = interface{}

type WorkspaceNewParams struct {
	DisplayName  param.Field[string] `json:"display_name" api:"required"`
	ID           param.Field[string] `json:"id" format:"uuid"`
	TenantHandle param.Field[string] `json:"tenant_handle"`
}

func (r WorkspaceNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type WorkspaceUpdateParams struct {
	DisplayName param.Field[string] `json:"display_name" api:"required"`
}

func (r WorkspaceUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type WorkspaceListParams struct {
	IncludeDeleted param.Field[bool] `query:"include_deleted"`
}

// URLQuery serializes [WorkspaceListParams]'s query parameters as `url.Values`.
func (r WorkspaceListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
