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

// WorkspaceSecretService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkspaceSecretService] method instead.
type WorkspaceSecretService struct {
	Options []option.RequestOption
}

// NewWorkspaceSecretService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewWorkspaceSecretService(opts ...option.RequestOption) (r *WorkspaceSecretService) {
	r = &WorkspaceSecretService{}
	r.Options = opts
	return
}

// List the current workspace's secret keys. Values are never returned.
func (r *WorkspaceSecretService) List(ctx context.Context, query WorkspaceSecretListParams, opts ...option.RequestOption) (res *[]WorkspaceSecret, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces/current/secrets"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete a workspace secret and its allowed hosts.
//
// Requires workspace admin. Missing keys succeed.
func (r *WorkspaceSecretService) Delete(ctx context.Context, key string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if key == "" {
		err = errors.New("missing required key parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/secrets/%s", key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Create or update a workspace secret and the hosts it may be sent to.
//
// Requires workspace admin. Omit value to keep the current one; it is required to
// create a secret. configs replaces the secret's hosts; omitted or empty restores
// the built-in defaults for known provider keys. The value is never returned.
func (r *WorkspaceSecretService) Upsert(ctx context.Context, key string, body WorkspaceSecretUpsertParams, opts ...option.RequestOption) (res *WorkspaceSecret, err error) {
	opts = slices.Concat(r.Options, opts)
	if key == "" {
		err = errors.New("missing required key parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/workspaces/current/secrets/%s", key)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, &res, opts...)
	return res, err
}

// A workspace secret's key and the hosts it may be sent to. Never the value.
type WorkspaceSecret struct {
	// Secret name, e.g. OPENAI_API_KEY.
	Key string `json:"key" api:"required"`
	// Hosts the secret may be sent to, with built-in defaults applied. Listed secrets
	// include it only when include_configs is true.
	Configs []WorkspaceSecretConfig `json:"configs"`
	JSON    workspaceSecretJSON     `json:"-"`
}

// workspaceSecretJSON contains the JSON metadata for the struct [WorkspaceSecret]
type workspaceSecretJSON struct {
	Key         apijson.Field
	Configs     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *WorkspaceSecret) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r workspaceSecretJSON) RawJSON() string {
	return r.raw
}

type WorkspaceSecretConfig struct {
	// Exact host or a leading '\*.' wildcard, e.g. 'api.openai.com'.
	HostPattern string `json:"host_pattern" api:"required"`
	// DEFAULT for a built-in provider host, CUSTOM for one an admin set.
	Source WorkspaceSecretConfigsSource `json:"source" api:"required"`
	// Omitted for a built-in default.
	ID string `json:"id" api:"nullable" format:"uuid"`
	// Omitted for a built-in default.
	CreatedAt time.Time `json:"created_at" api:"nullable" format:"date-time"`
	// HTTP header the secret is sent in. Null leaves it to the consumer.
	HeaderName string `json:"header_name" api:"nullable"`
	// Header value with '{value}' replaced by the secret, e.g. 'Bearer {value}'. Null
	// sends the raw secret.
	HeaderValueTemplate string `json:"header_value_template" api:"nullable"`
	// Omitted for a built-in default.
	UpdatedAt time.Time                 `json:"updated_at" api:"nullable" format:"date-time"`
	JSON      workspaceSecretConfigJSON `json:"-"`
}

// workspaceSecretConfigJSON contains the JSON metadata for the struct
// [WorkspaceSecretConfig]
type workspaceSecretConfigJSON struct {
	HostPattern         apijson.Field
	Source              apijson.Field
	ID                  apijson.Field
	CreatedAt           apijson.Field
	HeaderName          apijson.Field
	HeaderValueTemplate apijson.Field
	UpdatedAt           apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *WorkspaceSecretConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r workspaceSecretConfigJSON) RawJSON() string {
	return r.raw
}

// DEFAULT for a built-in provider host, CUSTOM for one an admin set.
type WorkspaceSecretConfigsSource string

const (
	WorkspaceSecretConfigsSourceDefault WorkspaceSecretConfigsSource = "DEFAULT"
	WorkspaceSecretConfigsSourceCustom  WorkspaceSecretConfigsSource = "CUSTOM"
)

func (r WorkspaceSecretConfigsSource) IsKnown() bool {
	switch r {
	case WorkspaceSecretConfigsSourceDefault, WorkspaceSecretConfigsSourceCustom:
		return true
	}
	return false
}

type WorkspaceSecretListParams struct {
	// If true, include the hosts each secret may be sent to.
	IncludeConfigs param.Field[bool] `query:"include_configs"`
}

// URLQuery serializes [WorkspaceSecretListParams]'s query parameters as
// `url.Values`.
func (r WorkspaceSecretListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WorkspaceSecretUpsertParams struct {
	// Full replacement of the allowed hosts. Omitted or empty removes custom hosts so
	// the built-in defaults, if any, apply.
	Configs param.Field[[]WorkspaceSecretUpsertParamsConfig] `json:"configs"`
	// Secret value. Required to create a secret; omit to keep the current value.
	Value param.Field[string] `json:"value"`
}

func (r WorkspaceSecretUpsertParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// A host a workspace secret may be sent to, and how it is sent.
type WorkspaceSecretUpsertParamsConfig struct {
	// Exact host or a leading '\*.' wildcard, e.g. 'api.openai.com'.
	HostPattern param.Field[string] `json:"host_pattern" api:"required"`
	// HTTP header the secret is sent in. Null leaves it to the consumer.
	HeaderName param.Field[string] `json:"header_name"`
	// Header value with '{value}' replaced by the secret, e.g. 'Bearer {value}'. Null
	// sends the raw secret.
	HeaderValueTemplate param.Field[string] `json:"header_value_template"`
}

func (r WorkspaceSecretUpsertParamsConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
