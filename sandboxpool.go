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

// SandboxPoolService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxPoolService] method instead.
type SandboxPoolService struct {
	Options []option.RequestOption
}

// NewSandboxPoolService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxPoolService(opts ...option.RequestOption) (r *SandboxPoolService) {
	r = &SandboxPoolService{}
	r.Options = opts
	return
}

// Create a new warm pool from a template with a specified replica count
func (r *SandboxPoolService) New(ctx context.Context, body SandboxPoolNewParams, opts ...option.RequestOption) (res *SandboxPoolNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/pools"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a warm pool by name
func (r *SandboxPoolService) Get(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxPoolGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/pools/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a warm pool's name or replica count
func (r *SandboxPoolService) Update(ctx context.Context, name string, body SandboxPoolUpdateParams, opts ...option.RequestOption) (res *SandboxPoolUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/pools/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List all warm pools for the current workspace
func (r *SandboxPoolService) List(ctx context.Context, query SandboxPoolListParams, opts ...option.RequestOption) (res *SandboxPoolListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/pools"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete a warm pool by name
func (r *SandboxPoolService) Delete(ctx context.Context, name string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("v2/sandboxes/pools/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type SandboxPoolNewResponse struct {
	ID           string                     `json:"id"`
	CreatedAt    string                     `json:"created_at"`
	Name         string                     `json:"name"`
	Replicas     int64                      `json:"replicas"`
	TemplateName string                     `json:"template_name"`
	UpdatedAt    string                     `json:"updated_at"`
	JSON         sandboxPoolNewResponseJSON `json:"-"`
}

// sandboxPoolNewResponseJSON contains the JSON metadata for the struct
// [SandboxPoolNewResponse]
type sandboxPoolNewResponseJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	Name         apijson.Field
	Replicas     apijson.Field
	TemplateName apijson.Field
	UpdatedAt    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxPoolNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxPoolNewResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxPoolGetResponse struct {
	ID           string                     `json:"id"`
	CreatedAt    string                     `json:"created_at"`
	Name         string                     `json:"name"`
	Replicas     int64                      `json:"replicas"`
	TemplateName string                     `json:"template_name"`
	UpdatedAt    string                     `json:"updated_at"`
	JSON         sandboxPoolGetResponseJSON `json:"-"`
}

// sandboxPoolGetResponseJSON contains the JSON metadata for the struct
// [SandboxPoolGetResponse]
type sandboxPoolGetResponseJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	Name         apijson.Field
	Replicas     apijson.Field
	TemplateName apijson.Field
	UpdatedAt    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxPoolGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxPoolGetResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxPoolUpdateResponse struct {
	ID           string                        `json:"id"`
	CreatedAt    string                        `json:"created_at"`
	Name         string                        `json:"name"`
	Replicas     int64                         `json:"replicas"`
	TemplateName string                        `json:"template_name"`
	UpdatedAt    string                        `json:"updated_at"`
	JSON         sandboxPoolUpdateResponseJSON `json:"-"`
}

// sandboxPoolUpdateResponseJSON contains the JSON metadata for the struct
// [SandboxPoolUpdateResponse]
type sandboxPoolUpdateResponseJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	Name         apijson.Field
	Replicas     apijson.Field
	TemplateName apijson.Field
	UpdatedAt    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxPoolUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxPoolUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxPoolListResponse struct {
	Offset int64                         `json:"offset"`
	Pools  []SandboxPoolListResponsePool `json:"pools"`
	JSON   sandboxPoolListResponseJSON   `json:"-"`
}

// sandboxPoolListResponseJSON contains the JSON metadata for the struct
// [SandboxPoolListResponse]
type sandboxPoolListResponseJSON struct {
	Offset      apijson.Field
	Pools       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxPoolListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxPoolListResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxPoolListResponsePool struct {
	ID           string                          `json:"id"`
	CreatedAt    string                          `json:"created_at"`
	Name         string                          `json:"name"`
	Replicas     int64                           `json:"replicas"`
	TemplateName string                          `json:"template_name"`
	UpdatedAt    string                          `json:"updated_at"`
	JSON         sandboxPoolListResponsePoolJSON `json:"-"`
}

// sandboxPoolListResponsePoolJSON contains the JSON metadata for the struct
// [SandboxPoolListResponsePool]
type sandboxPoolListResponsePoolJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	Name         apijson.Field
	Replicas     apijson.Field
	TemplateName apijson.Field
	UpdatedAt    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxPoolListResponsePool) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxPoolListResponsePoolJSON) RawJSON() string {
	return r.raw
}

type SandboxPoolNewParams struct {
	Name         param.Field[string] `json:"name" api:"required"`
	Replicas     param.Field[int64]  `json:"replicas" api:"required"`
	TemplateName param.Field[string] `json:"template_name" api:"required"`
	Timeout      param.Field[int64]  `json:"timeout"`
	WaitForReady param.Field[bool]   `json:"wait_for_ready"`
}

func (r SandboxPoolNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxPoolUpdateParams struct {
	Name     param.Field[string] `json:"name"`
	Replicas param.Field[int64]  `json:"replicas"`
}

func (r SandboxPoolUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxPoolListParams struct {
	// Max results
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Pagination offset
	Offset param.Field[int64] `query:"offset"`
}

// URLQuery serializes [SandboxPoolListParams]'s query parameters as `url.Values`.
func (r SandboxPoolListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
