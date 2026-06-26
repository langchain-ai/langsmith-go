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
func (r *SandboxRegistryService) New(ctx context.Context, body SandboxRegistryNewParams, opts ...option.RequestOption) (res *RegistryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/registries"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a sandbox registry by name.
func (r *SandboxRegistryService) Get(ctx context.Context, name string, opts ...option.RequestOption) (res *RegistryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/registries/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a sandbox registry's name and/or credentials.
func (r *SandboxRegistryService) Update(ctx context.Context, name string, body SandboxRegistryUpdateParams, opts ...option.RequestOption) (res *RegistryResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/registries/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List sandbox registries for pulling private images.
func (r *SandboxRegistryService) List(ctx context.Context, query SandboxRegistryListParams, opts ...option.RequestOption) (res *RegistryListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/registries"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete a sandbox registry by name.
func (r *SandboxRegistryService) Delete(ctx context.Context, name string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("v2/sandboxes/registries/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type RegistryListResponse struct {
	Offset     int64                    `json:"offset"`
	Registries []RegistryResponse       `json:"registries"`
	JSON       registryListResponseJSON `json:"-"`
}

// registryListResponseJSON contains the JSON metadata for the struct
// [RegistryListResponse]
type registryListResponseJSON struct {
	Offset      apijson.Field
	Registries  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RegistryListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r registryListResponseJSON) RawJSON() string {
	return r.raw
}

type RegistryResponse struct {
	ID        string               `json:"id"`
	CreatedAt string               `json:"created_at"`
	CreatedBy string               `json:"created_by"`
	Name      string               `json:"name"`
	UpdatedAt string               `json:"updated_at"`
	UpdatedBy string               `json:"updated_by"`
	URL       string               `json:"url"`
	JSON      registryResponseJSON `json:"-"`
}

// registryResponseJSON contains the JSON metadata for the struct
// [RegistryResponse]
type registryResponseJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	CreatedBy   apijson.Field
	Name        apijson.Field
	UpdatedAt   apijson.Field
	UpdatedBy   apijson.Field
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RegistryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r registryResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxRegistryNewParams struct {
	Name     param.Field[string] `json:"name" api:"required"`
	Password param.Field[string] `json:"password" api:"required"`
	URL      param.Field[string] `json:"url" api:"required"`
	Username param.Field[string] `json:"username" api:"required"`
}

func (r SandboxRegistryNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxRegistryUpdateParams struct {
	Name     param.Field[string] `json:"name"`
	Password param.Field[string] `json:"password"`
	URL      param.Field[string] `json:"url"`
	Username param.Field[string] `json:"username"`
}

func (r SandboxRegistryUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxRegistryListParams struct {
	// Maximum number of registries to return
	Limit param.Field[int64] `query:"limit"`
	// Filter to registries whose name contains this substring
	NameContains param.Field[string] `query:"name_contains"`
	// Number of registries to skip
	Offset param.Field[int64] `query:"offset"`
}

// URLQuery serializes [SandboxRegistryListParams]'s query parameters as
// `url.Values`.
func (r SandboxRegistryListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
