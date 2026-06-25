// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SandboxRegistryService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxRegistryService] method instead.
type SandboxRegistryService struct {
	Options []option.RequestOption
}

// NewSandboxRegistryService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxRegistryService(opts ...option.RequestOption) (r *SandboxRegistryService) {
	r = &SandboxRegistryService{}
	r.Options = opts
	return
}

// Create a sandbox registry for pulling private images.
func (r *SandboxRegistryService) New(ctx context.Context, body SandboxRegistryNewParams, opts ...option.RequestOption) (res *SandboxRegistryNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/registries"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type SandboxRegistryNewResponse struct {
	ID        string                         `json:"id"`
	Name      string                         `json:"name"`
	URL       string                         `json:"url"`
	CreatedAt string                         `json:"created_at"`
	CreatedBy string                         `json:"created_by"`
	UpdatedAt string                         `json:"updated_at"`
	UpdatedBy string                         `json:"updated_by"`
	JSON      sandboxRegistryNewResponseJSON `json:"-"`
}

// sandboxRegistryNewResponseJSON contains the JSON metadata for the struct
// [SandboxRegistryNewResponse]
type sandboxRegistryNewResponseJSON struct {
	ID          apijson.Field
	Name        apijson.Field
	URL         apijson.Field
	CreatedAt   apijson.Field
	CreatedBy   apijson.Field
	UpdatedAt   apijson.Field
	UpdatedBy   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxRegistryNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxRegistryNewResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxRegistryNewParams struct {
	Name     param.Field[string] `json:"name" api:"required"`
	URL      param.Field[string] `json:"url" api:"required"`
	Username param.Field[string] `json:"username" api:"required"`
	Password param.Field[string] `json:"password" api:"required"`
}

func (r SandboxRegistryNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
