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

// ThreadShareService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewThreadShareService] method instead.
type ThreadShareService struct {
	Options []option.RequestOption
}

// NewThreadShareService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewThreadShareService(opts ...option.RequestOption) (r *ThreadShareService) {
	r = &ThreadShareService{}
	r.Options = opts
	return
}

// Mints a public share token for a thread. Idempotent: sharing an already-shared
// thread returns the existing token.
func (r *ThreadShareService) New(ctx context.Context, threadID string, body ThreadShareNewParams, opts ...option.RequestOption) (res *ThreadShareNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/threads/%s/share", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns the share token for a thread, or 404 when it is not shared. Gated on
// runs:share so the control's state matches the control's permission.
func (r *ThreadShareService) Get(ctx context.Context, threadID string, query ThreadShareGetParams, opts ...option.RequestOption) (res *ThreadShareGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/threads/%s/share", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Deletes the share token for a thread. Idempotent: returns 204 whether or not a
// share token existed. Deliberately does not verify the thread still exists.
func (r *ThreadShareService) Delete(ctx context.Context, threadID string, body ThreadShareDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/threads/%s/share", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return err
}

type ThreadShareNewResponse struct {
	ShareToken string                     `json:"share_token" api:"required" format:"uuid"`
	JSON       threadShareNewResponseJSON `json:"-"`
}

// threadShareNewResponseJSON contains the JSON metadata for the struct
// [ThreadShareNewResponse]
type threadShareNewResponseJSON struct {
	ShareToken  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadShareNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadShareNewResponseJSON) RawJSON() string {
	return r.raw
}

type ThreadShareGetResponse struct {
	ShareToken string                     `json:"share_token" api:"required" format:"uuid"`
	JSON       threadShareGetResponseJSON `json:"-"`
}

// threadShareGetResponseJSON contains the JSON metadata for the struct
// [ThreadShareGetResponse]
type threadShareGetResponseJSON struct {
	ShareToken  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ThreadShareGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r threadShareGetResponseJSON) RawJSON() string {
	return r.raw
}

type ThreadShareNewParams struct {
	// project_id is the tracing project UUID containing the thread.
	ProjectID param.Field[string] `json:"project_id" api:"required" format:"uuid"`
}

func (r ThreadShareNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ThreadShareGetParams struct {
	// Project UUID
	ProjectID param.Field[string] `query:"project_id" api:"required" format:"uuid"`
}

// URLQuery serializes [ThreadShareGetParams]'s query parameters as `url.Values`.
func (r ThreadShareGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type ThreadShareDeleteParams struct {
	// Project UUID
	ProjectID param.Field[string] `query:"project_id" api:"required" format:"uuid"`
}

// URLQuery serializes [ThreadShareDeleteParams]'s query parameters as
// `url.Values`.
func (r ThreadShareDeleteParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
