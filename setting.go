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

// Get settings.
func (r *SettingService) List(ctx context.Context, opts ...option.RequestOption) (res *AppHubCrudTenantsTenant, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/settings"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

type AppHubCrudTenantsTenant struct {
	ID           string                      `json:"id,required" format:"uuid"`
	CreatedAt    time.Time                   `json:"created_at,required" format:"date-time"`
	DisplayName  string                      `json:"display_name,required"`
	TenantHandle string                      `json:"tenant_handle,nullable"`
	JSON         appHubCrudTenantsTenantJSON `json:"-"`
}

// appHubCrudTenantsTenantJSON contains the JSON metadata for the struct
// [AppHubCrudTenantsTenant]
type appHubCrudTenantsTenantJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	DisplayName  apijson.Field
	TenantHandle apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *AppHubCrudTenantsTenant) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r appHubCrudTenantsTenantJSON) RawJSON() string {
	return r.raw
}
