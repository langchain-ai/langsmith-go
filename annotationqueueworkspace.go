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

// AnnotationQueueWorkspaceService contains methods and other services that help
// with interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAnnotationQueueWorkspaceService] method instead.
type AnnotationQueueWorkspaceService struct {
	Options []option.RequestOption
}

// NewAnnotationQueueWorkspaceService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewAnnotationQueueWorkspaceService(opts ...option.RequestOption) (r *AnnotationQueueWorkspaceService) {
	r = &AnnotationQueueWorkspaceService{}
	r.Options = opts
	return
}

// Create a new workspace.
func (r *AnnotationQueueWorkspaceService) New(ctx context.Context, body AnnotationQueueWorkspaceNewParams, opts ...option.RequestOption) (res *AnnotationQueueWorkspaceNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Update a workspace.
func (r *AnnotationQueueWorkspaceService) Update(ctx context.Context, workspaceID string, body AnnotationQueueWorkspaceUpdateParams, opts ...option.RequestOption) (res *AnnotationQueueWorkspaceUpdateResponse, err error) {
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
func (r *AnnotationQueueWorkspaceService) List(ctx context.Context, query AnnotationQueueWorkspaceListParams, opts ...option.RequestOption) (res *[]AnnotationQueueWorkspaceListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete Workspace
func (r *AnnotationQueueWorkspaceService) Delete(ctx context.Context, workspaceID string, opts ...option.RequestOption) (res *AnnotationQueueWorkspaceDeleteResponse, err error) {
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
type AnnotationQueueWorkspaceNewResponse struct {
	ID             string                                  `json:"id" api:"required" format:"uuid"`
	CreatedAt      time.Time                               `json:"created_at" api:"required" format:"date-time"`
	DisplayName    string                                  `json:"display_name" api:"required"`
	IsDeleted      bool                                    `json:"is_deleted" api:"required"`
	IsPersonal     bool                                    `json:"is_personal" api:"required"`
	DataPlaneURL   string                                  `json:"data_plane_url" api:"nullable"`
	OrganizationID string                                  `json:"organization_id" api:"nullable" format:"uuid"`
	TenantHandle   string                                  `json:"tenant_handle" api:"nullable"`
	JSON           annotationQueueWorkspaceNewResponseJSON `json:"-"`
}

// annotationQueueWorkspaceNewResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueWorkspaceNewResponse]
type annotationQueueWorkspaceNewResponseJSON struct {
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

func (r *AnnotationQueueWorkspaceNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueWorkspaceNewResponseJSON) RawJSON() string {
	return r.raw
}

// Tenant schema.
type AnnotationQueueWorkspaceUpdateResponse struct {
	ID             string                                     `json:"id" api:"required" format:"uuid"`
	CreatedAt      time.Time                                  `json:"created_at" api:"required" format:"date-time"`
	DisplayName    string                                     `json:"display_name" api:"required"`
	IsDeleted      bool                                       `json:"is_deleted" api:"required"`
	IsPersonal     bool                                       `json:"is_personal" api:"required"`
	DataPlaneURL   string                                     `json:"data_plane_url" api:"nullable"`
	OrganizationID string                                     `json:"organization_id" api:"nullable" format:"uuid"`
	TenantHandle   string                                     `json:"tenant_handle" api:"nullable"`
	JSON           annotationQueueWorkspaceUpdateResponseJSON `json:"-"`
}

// annotationQueueWorkspaceUpdateResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueWorkspaceUpdateResponse]
type annotationQueueWorkspaceUpdateResponseJSON struct {
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

func (r *AnnotationQueueWorkspaceUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueWorkspaceUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueWorkspaceListResponse struct {
	ID             string    `json:"id" api:"required" format:"uuid"`
	CreatedAt      time.Time `json:"created_at" api:"required" format:"date-time"`
	DisplayName    string    `json:"display_name" api:"required"`
	IsDeleted      bool      `json:"is_deleted" api:"required"`
	IsPersonal     bool      `json:"is_personal" api:"required"`
	DataPlaneURL   string    `json:"data_plane_url" api:"nullable"`
	OrganizationID string    `json:"organization_id" api:"nullable" format:"uuid"`
	Permissions    []string  `json:"permissions" api:"nullable"`
	// Deprecated: deprecated
	ReadOnly     bool                                     `json:"read_only"`
	RoleID       string                                   `json:"role_id" api:"nullable" format:"uuid"`
	RoleName     string                                   `json:"role_name" api:"nullable"`
	TenantHandle string                                   `json:"tenant_handle" api:"nullable"`
	JSON         annotationQueueWorkspaceListResponseJSON `json:"-"`
}

// annotationQueueWorkspaceListResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueWorkspaceListResponse]
type annotationQueueWorkspaceListResponseJSON struct {
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

func (r *AnnotationQueueWorkspaceListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueWorkspaceListResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueWorkspaceDeleteResponse = interface{}

type AnnotationQueueWorkspaceNewParams struct {
	DisplayName  param.Field[string] `json:"display_name" api:"required"`
	ID           param.Field[string] `json:"id" format:"uuid"`
	TenantHandle param.Field[string] `json:"tenant_handle"`
}

func (r AnnotationQueueWorkspaceNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueWorkspaceUpdateParams struct {
	DisplayName param.Field[string] `json:"display_name" api:"required"`
}

func (r AnnotationQueueWorkspaceUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueWorkspaceListParams struct {
	IncludeDeleted param.Field[bool] `query:"include_deleted"`
}

// URLQuery serializes [AnnotationQueueWorkspaceListParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueWorkspaceListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
