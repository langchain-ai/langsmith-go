// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// WorkspaceTagKeyService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkspaceTagKeyService] method instead.
type WorkspaceTagKeyService struct {
	Options []option.RequestOption
	Values  *WorkspaceTagKeyValueService
}

// NewWorkspaceTagKeyService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewWorkspaceTagKeyService(opts ...option.RequestOption) (r *WorkspaceTagKeyService) {
	r = &WorkspaceTagKeyService{}
	r.Options = opts
	r.Values = NewWorkspaceTagKeyValueService(opts...)
	return
}

// Create Tag Key
func (r *WorkspaceTagKeyService) New(ctx context.Context, body WorkspaceTagKeyNewParams, opts ...option.RequestOption) (res *TagKey, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces/current/tag-keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get Tag Key
func (r *WorkspaceTagKeyService) Get(ctx context.Context, tagKeyID string, opts ...option.RequestOption) (res *TagKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s", tagKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Tag Key
func (r *WorkspaceTagKeyService) Update(ctx context.Context, tagKeyID string, body WorkspaceTagKeyUpdateParams, opts ...option.RequestOption) (res *TagKey, err error) {
	opts = slices.Concat(r.Options, opts)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s", tagKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List Tag Keys
func (r *WorkspaceTagKeyService) List(ctx context.Context, opts ...option.RequestOption) (res *[]TagKey, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces/current/tag-keys"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Delete Tag Key
func (r *WorkspaceTagKeyService) Delete(ctx context.Context, tagKeyID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s", tagKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type TagKey struct {
	ID          string     `json:"id" api:"required" format:"uuid"`
	CreatedAt   time.Time  `json:"created_at" api:"required" format:"date-time"`
	Key         string     `json:"key" api:"required"`
	UpdatedAt   time.Time  `json:"updated_at" api:"required" format:"date-time"`
	Description string     `json:"description" api:"nullable"`
	JSON        tagKeyJSON `json:"-"`
}

// tagKeyJSON contains the JSON metadata for the struct [TagKey]
type tagKeyJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	Key         apijson.Field
	UpdatedAt   apijson.Field
	Description apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TagKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r tagKeyJSON) RawJSON() string {
	return r.raw
}

type WorkspaceTagKeyNewParams struct {
	Key         param.Field[string] `json:"key" api:"required"`
	Description param.Field[string] `json:"description"`
}

func (r WorkspaceTagKeyNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type WorkspaceTagKeyUpdateParams struct {
	Description param.Field[string] `json:"description"`
	Key         param.Field[string] `json:"key"`
}

func (r WorkspaceTagKeyUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
