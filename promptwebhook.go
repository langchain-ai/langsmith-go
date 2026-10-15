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

// PromptWebhookService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewPromptWebhookService] method instead.
type PromptWebhookService struct {
	Options []option.RequestOption
}

// NewPromptWebhookService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewPromptWebhookService(opts ...option.RequestOption) (r *PromptWebhookService) {
	r = &PromptWebhookService{}
	r.Options = opts
	return
}

// Create a new prompt webhook.
//
// LangSmith sends a POST request to the webhook URL for each event listed in
// triggers on the prompts it matches. A workspace can have one prompt webhook
// unless its limit was raised; creating one past the limit returns 400.
func (r *PromptWebhookService) New(ctx context.Context, body PromptWebhookNewParams, opts ...option.RequestOption) (res *PromptWebhook, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/prompt-webhooks"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get a specific prompt webhook.
func (r *PromptWebhookService) Get(ctx context.Context, webhookID string, opts ...option.RequestOption) (res *PromptWebhook, err error) {
	opts = slices.Concat(r.Options, opts)
	if webhookID == "" {
		err = errors.New("missing required webhook_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/prompt-webhooks/%s", webhookID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a specific prompt webhook.
//
// Omitted fields keep their current values; list fields and headers that are set
// replace the stored values.
func (r *PromptWebhookService) Update(ctx context.Context, webhookID string, body PromptWebhookUpdateParams, opts ...option.RequestOption) (res *PromptWebhook, err error) {
	opts = slices.Concat(r.Options, opts)
	if webhookID == "" {
		err = errors.New("missing required webhook_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/prompt-webhooks/%s", webhookID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List all prompt webhooks for the current tenant.
func (r *PromptWebhookService) List(ctx context.Context, opts ...option.RequestOption) (res *[]PromptWebhook, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/prompt-webhooks"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Delete a specific prompt webhook.
func (r *PromptWebhookService) Delete(ctx context.Context, webhookID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if webhookID == "" {
		err = errors.New("missing required webhook_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/prompt-webhooks/%s", webhookID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Test a specific prompt webhook.
//
// Sends the sample payload to the webhook configuration in the request body
// without saving it, but only if the webhook's triggers include the payload's
// event and its prompt filters match the payload's prompt. Otherwise nothing is
// sent. Fails if the destination does not accept the payload.
func (r *PromptWebhookService) Test(ctx context.Context, body PromptWebhookTestParams, opts ...option.RequestOption) (res *PromptWebhookTestResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/prompt-webhooks/test"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Schema for a prompt webhook.
type PromptWebhook struct {
	ID        string    `json:"id" api:"required" format:"uuid"`
	CreatedAt time.Time `json:"created_at" api:"required" format:"date-time"`
	TenantID  string    `json:"tenant_id" api:"required" format:"uuid"`
	UpdatedAt time.Time `json:"updated_at" api:"required" format:"date-time"`
	// URL that receives a POST request with the prompt event payload.
	URL string `json:"url" api:"required" format:"uri"`
	// IDs of prompts whose events never fire the webhook.
	ExcludePrompts []string `json:"exclude_prompts" api:"nullable" format:"uuid"`
	// HTTP headers sent with each request, for example an Authorization header. Values
	// are stored encrypted and are masked in responses to callers who cannot update
	// prompts.
	Headers map[string]interface{} `json:"headers" api:"nullable"`
	// IDs of the prompts whose events fire the webhook. When empty, events from every
	// prompt in the workspace fire it.
	IncludePrompts []string `json:"include_prompts" api:"nullable" format:"uuid"`
	// Events that fire the webhook: commit (a new prompt commit), tag:create, and
	// tag:update. A webhook with no triggers never fires.
	Triggers []PromptWebhookTrigger `json:"triggers"`
	JSON     promptWebhookJSON      `json:"-"`
}

// promptWebhookJSON contains the JSON metadata for the struct [PromptWebhook]
type promptWebhookJSON struct {
	ID             apijson.Field
	CreatedAt      apijson.Field
	TenantID       apijson.Field
	UpdatedAt      apijson.Field
	URL            apijson.Field
	ExcludePrompts apijson.Field
	Headers        apijson.Field
	IncludePrompts apijson.Field
	Triggers       apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *PromptWebhook) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r promptWebhookJSON) RawJSON() string {
	return r.raw
}

// Valid trigger types for prompt webhooks.
type PromptWebhookTrigger string

const (
	PromptWebhookTriggerCommit    PromptWebhookTrigger = "commit"
	PromptWebhookTriggerTagCreate PromptWebhookTrigger = "tag:create"
	PromptWebhookTriggerTagUpdate PromptWebhookTrigger = "tag:update"
)

func (r PromptWebhookTrigger) IsKnown() bool {
	switch r {
	case PromptWebhookTriggerCommit, PromptWebhookTriggerTagCreate, PromptWebhookTriggerTagUpdate:
		return true
	}
	return false
}

type PromptWebhookTestResponse map[string]string

type PromptWebhookNewParams struct {
	// URL that receives a POST request with the prompt event payload.
	URL param.Field[string] `json:"url" api:"required" format:"uri"`
	// ID to assign to the webhook. Generated when omitted.
	ID param.Field[string] `json:"id" format:"uuid"`
	// IDs of prompts whose events never fire the webhook.
	ExcludePrompts param.Field[[]string] `json:"exclude_prompts" format:"uuid"`
	// HTTP headers sent with each request, for example an Authorization header. Values
	// are stored encrypted and are masked in responses to callers who cannot update
	// prompts.
	Headers param.Field[map[string]interface{}] `json:"headers"`
	// IDs of the prompts whose events fire the webhook. When empty, events from every
	// prompt in the workspace fire it.
	IncludePrompts param.Field[[]string] `json:"include_prompts" format:"uuid"`
	// Events that fire the webhook: commit (a new prompt commit), tag:create, and
	// tag:update. A webhook with no triggers never fires.
	Triggers param.Field[[]PromptWebhookNewParamsTrigger] `json:"triggers"`
}

func (r PromptWebhookNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Valid trigger types for prompt webhooks.
type PromptWebhookNewParamsTrigger string

const (
	PromptWebhookNewParamsTriggerCommit    PromptWebhookNewParamsTrigger = "commit"
	PromptWebhookNewParamsTriggerTagCreate PromptWebhookNewParamsTrigger = "tag:create"
	PromptWebhookNewParamsTriggerTagUpdate PromptWebhookNewParamsTrigger = "tag:update"
)

func (r PromptWebhookNewParamsTrigger) IsKnown() bool {
	switch r {
	case PromptWebhookNewParamsTriggerCommit, PromptWebhookNewParamsTriggerTagCreate, PromptWebhookNewParamsTriggerTagUpdate:
		return true
	}
	return false
}

type PromptWebhookUpdateParams struct {
	// IDs of prompts whose events never fire the webhook, replacing the current list.
	// Omit to keep it.
	ExcludePrompts param.Field[[]string] `json:"exclude_prompts" format:"uuid"`
	// HTTP headers sent with each request, replacing all current headers. Omit to keep
	// them. Masked values (**\*\*\*\***) from a read are rejected.
	Headers param.Field[map[string]interface{}] `json:"headers"`
	// IDs of the prompts whose events fire the webhook, replacing the current list.
	// Omit to keep it.
	IncludePrompts param.Field[[]string] `json:"include_prompts" format:"uuid"`
	// Events that fire the webhook (commit, tag:create, tag:update), replacing the
	// current list. Omit to keep it.
	Triggers param.Field[[]PromptWebhookUpdateParamsTrigger] `json:"triggers"`
	// URL that receives the POST requests. Omit to keep the current URL.
	URL param.Field[string] `json:"url" format:"uri"`
}

func (r PromptWebhookUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Valid trigger types for prompt webhooks.
type PromptWebhookUpdateParamsTrigger string

const (
	PromptWebhookUpdateParamsTriggerCommit    PromptWebhookUpdateParamsTrigger = "commit"
	PromptWebhookUpdateParamsTriggerTagCreate PromptWebhookUpdateParamsTrigger = "tag:create"
	PromptWebhookUpdateParamsTriggerTagUpdate PromptWebhookUpdateParamsTrigger = "tag:update"
)

func (r PromptWebhookUpdateParamsTrigger) IsKnown() bool {
	switch r {
	case PromptWebhookUpdateParamsTriggerCommit, PromptWebhookUpdateParamsTriggerTagCreate, PromptWebhookUpdateParamsTriggerTagUpdate:
		return true
	}
	return false
}

type PromptWebhookTestParams struct {
	Payload param.Field[PromptWebhookTestParamsPayload] `json:"payload" api:"required"`
	// Base schema for prompt webhooks.
	Webhook param.Field[PromptWebhookTestParamsWebhook] `json:"webhook" api:"required"`
}

func (r PromptWebhookTestParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type PromptWebhookTestParamsPayload struct {
	CommitHash param.Field[string] `json:"commit_hash" api:"required"`
	CreatedAt  param.Field[string] `json:"created_at" api:"required"`
	CreatedBy  param.Field[string] `json:"created_by" api:"required"`
	// Valid trigger types for prompt webhooks.
	Event      param.Field[PromptWebhookTestParamsPayloadEvent] `json:"event" api:"required"`
	Manifest   param.Field[map[string]interface{}]              `json:"manifest" api:"required"`
	PromptID   param.Field[string]                              `json:"prompt_id" api:"required"`
	PromptName param.Field[string]                              `json:"prompt_name" api:"required"`
	TagName    param.Field[string]                              `json:"tag_name"`
}

func (r PromptWebhookTestParamsPayload) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Valid trigger types for prompt webhooks.
type PromptWebhookTestParamsPayloadEvent string

const (
	PromptWebhookTestParamsPayloadEventCommit    PromptWebhookTestParamsPayloadEvent = "commit"
	PromptWebhookTestParamsPayloadEventTagCreate PromptWebhookTestParamsPayloadEvent = "tag:create"
	PromptWebhookTestParamsPayloadEventTagUpdate PromptWebhookTestParamsPayloadEvent = "tag:update"
)

func (r PromptWebhookTestParamsPayloadEvent) IsKnown() bool {
	switch r {
	case PromptWebhookTestParamsPayloadEventCommit, PromptWebhookTestParamsPayloadEventTagCreate, PromptWebhookTestParamsPayloadEventTagUpdate:
		return true
	}
	return false
}

// Base schema for prompt webhooks.
type PromptWebhookTestParamsWebhook struct {
	// URL that receives a POST request with the prompt event payload.
	URL param.Field[string] `json:"url" api:"required" format:"uri"`
	// IDs of prompts whose events never fire the webhook.
	ExcludePrompts param.Field[[]string] `json:"exclude_prompts" format:"uuid"`
	// HTTP headers sent with each request, for example an Authorization header. Values
	// are stored encrypted and are masked in responses to callers who cannot update
	// prompts.
	Headers param.Field[map[string]interface{}] `json:"headers"`
	// IDs of the prompts whose events fire the webhook. When empty, events from every
	// prompt in the workspace fire it.
	IncludePrompts param.Field[[]string] `json:"include_prompts" format:"uuid"`
	// Events that fire the webhook: commit (a new prompt commit), tag:create, and
	// tag:update. A webhook with no triggers never fires.
	Triggers param.Field[[]PromptWebhookTestParamsWebhookTrigger] `json:"triggers"`
}

func (r PromptWebhookTestParamsWebhook) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Valid trigger types for prompt webhooks.
type PromptWebhookTestParamsWebhookTrigger string

const (
	PromptWebhookTestParamsWebhookTriggerCommit    PromptWebhookTestParamsWebhookTrigger = "commit"
	PromptWebhookTestParamsWebhookTriggerTagCreate PromptWebhookTestParamsWebhookTrigger = "tag:create"
	PromptWebhookTestParamsWebhookTriggerTagUpdate PromptWebhookTestParamsWebhookTrigger = "tag:update"
)

func (r PromptWebhookTestParamsWebhookTrigger) IsKnown() bool {
	switch r {
	case PromptWebhookTestParamsWebhookTriggerCommit, PromptWebhookTestParamsWebhookTriggerTagCreate, PromptWebhookTestParamsWebhookTriggerTagUpdate:
		return true
	}
	return false
}
