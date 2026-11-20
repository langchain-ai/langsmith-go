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

// WorkspaceTagKeyValueService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkspaceTagKeyValueService] method instead.
type WorkspaceTagKeyValueService struct {
	Options []option.RequestOption
}

// NewWorkspaceTagKeyValueService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewWorkspaceTagKeyValueService(opts ...option.RequestOption) (r *WorkspaceTagKeyValueService) {
	r = &WorkspaceTagKeyValueService{}
	r.Options = opts
	return
}

// Create Tag Value
func (r *WorkspaceTagKeyValueService) New(ctx context.Context, tagKeyID string, body WorkspaceTagKeyValueNewParams, opts ...option.RequestOption) (res *TagValue, err error) {
	opts = slices.Concat(r.Options, opts)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s/tag-values", tagKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get Tag Value
func (r *WorkspaceTagKeyValueService) Get(ctx context.Context, tagKeyID string, tagValueID string, opts ...option.RequestOption) (res *TagValue, err error) {
	opts = slices.Concat(r.Options, opts)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return nil, err
	}
	if tagValueID == "" {
		err = errors.New("missing required tag_value_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s/tag-values/%s", tagKeyID, tagValueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Tag Value
func (r *WorkspaceTagKeyValueService) Update(ctx context.Context, tagKeyID string, tagValueID string, body WorkspaceTagKeyValueUpdateParams, opts ...option.RequestOption) (res *TagValue, err error) {
	opts = slices.Concat(r.Options, opts)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return nil, err
	}
	if tagValueID == "" {
		err = errors.New("missing required tag_value_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s/tag-values/%s", tagKeyID, tagValueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List Tag Values
func (r *WorkspaceTagKeyValueService) List(ctx context.Context, tagKeyID string, opts ...option.RequestOption) (res *[]TagValue, err error) {
	opts = slices.Concat(r.Options, opts)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s/tag-values", tagKeyID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Delete Tag Value
func (r *WorkspaceTagKeyValueService) Delete(ctx context.Context, tagKeyID string, tagValueID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if tagKeyID == "" {
		err = errors.New("missing required tag_key_id parameter")
		return err
	}
	if tagValueID == "" {
		err = errors.New("missing required tag_value_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/tag-keys/%s/tag-values/%s", tagKeyID, tagValueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type TagValue struct {
	ID          string       `json:"id" api:"required" format:"uuid"`
	CreatedAt   time.Time    `json:"created_at" api:"required" format:"date-time"`
	TagKeyID    string       `json:"tag_key_id" api:"required" format:"uuid"`
	UpdatedAt   time.Time    `json:"updated_at" api:"required" format:"date-time"`
	Value       string       `json:"value" api:"required"`
	Description string       `json:"description" api:"nullable"`
	JSON        tagValueJSON `json:"-"`
}

// tagValueJSON contains the JSON metadata for the struct [TagValue]
type tagValueJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	TagKeyID    apijson.Field
	UpdatedAt   apijson.Field
	Value       apijson.Field
	Description apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TagValue) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r tagValueJSON) RawJSON() string {
	return r.raw
}

type WorkspaceTagKeyValueNewParams struct {
	Value       param.Field[string] `json:"value" api:"required"`
	Description param.Field[string] `json:"description"`
}

func (r WorkspaceTagKeyValueNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type WorkspaceTagKeyValueUpdateParams struct {
	Description param.Field[string] `json:"description"`
	Value       param.Field[string] `json:"value"`
}

func (r WorkspaceTagKeyValueUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
