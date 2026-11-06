// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SettingService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSettingService] method instead.
type SettingService struct {
	Options []option.RequestOption
}

// NewSettingService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSettingService(opts ...option.RequestOption) (r *SettingService) {
	r = &SettingService{}
	r.Options = opts
	return
}

// Get the current workspace.
func (r *SettingService) List(ctx context.Context, opts ...option.RequestOption) (res *SettingListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/settings"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type SettingListResponse struct {
	ID              string    `json:"id" api:"required" format:"uuid"`
	CreatedAt       time.Time `json:"created_at" api:"required" format:"date-time"`
	DisplayName     string    `json:"display_name" api:"required"`
	IsDeleted       bool      `json:"is_deleted" api:"required"`
	IsPersonal      bool      `json:"is_personal" api:"required"`
	DataPlaneName   string    `json:"data_plane_name" api:"nullable"`
	DataPlaneRegion string    `json:"data_plane_region" api:"nullable"`
	DataPlaneURL    string    `json:"data_plane_url" api:"nullable"`
	OrganizationID  string    `json:"organization_id" api:"nullable" format:"uuid"`
	Permissions     []string  `json:"permissions" api:"nullable"`
	// Deprecated: deprecated
	ReadOnly     bool                    `json:"read_only"`
	RoleID       string                  `json:"role_id" api:"nullable" format:"uuid"`
	RoleName     string                  `json:"role_name" api:"nullable"`
	TenantHandle string                  `json:"tenant_handle" api:"nullable"`
	JSON         settingListResponseJSON `json:"-"`
}

// settingListResponseJSON contains the JSON metadata for the struct
// [SettingListResponse]
type settingListResponseJSON struct {
	ID              apijson.Field
	CreatedAt       apijson.Field
	DisplayName     apijson.Field
	IsDeleted       apijson.Field
	IsPersonal      apijson.Field
	DataPlaneName   apijson.Field
	DataPlaneRegion apijson.Field
	DataPlaneURL    apijson.Field
	OrganizationID  apijson.Field
	Permissions     apijson.Field
	ReadOnly        apijson.Field
	RoleID          apijson.Field
	RoleName        apijson.Field
	TenantHandle    apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SettingListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r settingListResponseJSON) RawJSON() string {
	return r.raw
}
