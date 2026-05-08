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

// SandboxBoxService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxBoxService] method instead.
type SandboxBoxService struct {
	Options []option.RequestOption
}

// NewSandboxBoxService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxBoxService(opts ...option.RequestOption) (r *SandboxBoxService) {
	r = &SandboxBoxService{}
	r.Options = opts
	return
}

// Create a new sandbox from a snapshot. The snapshot may be identified by
// `snapshot_id` (UUID) or by `snapshot_name` (tenant-scoped unique name); exactly
// one must be set.
func (r *SandboxBoxService) New(ctx context.Context, body SandboxBoxNewParams, opts ...option.RequestOption) (res *SandboxBoxNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/boxes"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a sandbox claim by name. Stale provisioning claims are auto-failed.
func (r *SandboxBoxService) Get(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxBoxGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a sandbox claim's display name. The name must be unique within the
// tenant.
func (r *SandboxBoxService) Update(ctx context.Context, name string, body SandboxBoxUpdateParams, opts ...option.RequestOption) (res *SandboxBoxUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List sandbox claims for the authenticated tenant, with optional filtering,
// sorting, and pagination.
func (r *SandboxBoxService) List(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) (res *SandboxBoxListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/boxes"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete a sandbox claim by name or UUID. Tears down the sandbox runtime and
// removes the DB record.
func (r *SandboxBoxService) Delete(ctx context.Context, name string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Create a snapshot by capturing the current state of a sandbox or promoting an
// existing checkpoint.
func (r *SandboxBoxService) NewSnapshot(ctx context.Context, name string, body SandboxBoxNewSnapshotParams, opts ...option.RequestOption) (res *SandboxBoxNewSnapshotResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/snapshot", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Create a short-lived JWT for accessing an HTTP service running on a specific
// port inside a sandbox. Returns a browser_url (sets auth cookie via redirect), a
// service_url (for use with the X-Langsmith-Sandbox-Service-Token header), the raw
// token, and its expiry.
func (r *SandboxBoxService) GenerateServiceURL(ctx context.Context, name string, body SandboxBoxGenerateServiceURLParams, opts ...option.RequestOption) (res *SandboxBoxGenerateServiceURLResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/service-url", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve the lightweight status of a sandbox claim for polling.
func (r *SandboxBoxService) GetStatus(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxBoxGetStatusResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/status", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Start a stopped or failed sandbox. This endpoint is not idempotent.
func (r *SandboxBoxService) Start(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxBoxStartResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/start", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Stop a ready sandbox. This endpoint is not idempotent; the filesystem is
// preserved for later restart.
func (r *SandboxBoxService) Stop(ctx context.Context, name string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/stop", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return err
}

type SandboxBoxNewResponse struct {
	ID                     string                           `json:"id"`
	CreatedAt              string                           `json:"created_at"`
	CreatedBy              string                           `json:"created_by"`
	DataplaneURL           string                           `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                            `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                            `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                            `json:"idle_ttl_seconds"`
	MemBytes               int64                            `json:"mem_bytes"`
	Name                   string                           `json:"name"`
	ProxyConfig            SandboxBoxNewResponseProxyConfig `json:"proxy_config"`
	SizeClass              string                           `json:"size_class"`
	SnapshotID             string                           `json:"snapshot_id"`
	Status                 string                           `json:"status"`
	StatusMessage          string                           `json:"status_message"`
	StoppedAt              string                           `json:"stopped_at"`
	UpdatedAt              string                           `json:"updated_at"`
	UpdatedBy              string                           `json:"updated_by"`
	Vcpus                  int64                            `json:"vcpus"`
	JSON                   sandboxBoxNewResponseJSON        `json:"-"`
}

// sandboxBoxNewResponseJSON contains the JSON metadata for the struct
// [SandboxBoxNewResponse]
type sandboxBoxNewResponseJSON struct {
	ID                     apijson.Field
	CreatedAt              apijson.Field
	CreatedBy              apijson.Field
	DataplaneURL           apijson.Field
	DeleteAfterStopSeconds apijson.Field
	FsCapacityBytes        apijson.Field
	IdleTtlSeconds         apijson.Field
	MemBytes               apijson.Field
	Name                   apijson.Field
	ProxyConfig            apijson.Field
	SizeClass              apijson.Field
	SnapshotID             apijson.Field
	Status                 apijson.Field
	StatusMessage          apijson.Field
	StoppedAt              apijson.Field
	UpdatedAt              apijson.Field
	UpdatedBy              apijson.Field
	Vcpus                  apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfig struct {
	AccessControl SandboxBoxNewResponseProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxBoxNewResponseProxyConfigCallback    `json:"callbacks"`
	NoProxy       []string                                      `json:"no_proxy"`
	Rules         []SandboxBoxNewResponseProxyConfigRule        `json:"rules"`
	JSON          sandboxBoxNewResponseProxyConfigJSON          `json:"-"`
}

// sandboxBoxNewResponseProxyConfigJSON contains the JSON metadata for the struct
// [SandboxBoxNewResponseProxyConfig]
type sandboxBoxNewResponseProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigAccessControl struct {
	AllowList []string                                          `json:"allow_list"`
	DenyList  []string                                          `json:"deny_list"`
	JSON      sandboxBoxNewResponseProxyConfigAccessControlJSON `json:"-"`
}

// sandboxBoxNewResponseProxyConfigAccessControlJSON contains the JSON metadata for
// the struct [SandboxBoxNewResponseProxyConfigAccessControl]
type sandboxBoxNewResponseProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigCallback struct {
	MatchHosts     []string                                                 `json:"match_hosts" api:"required"`
	TtlSeconds     int64                                                    `json:"ttl_seconds" api:"required"`
	URL            string                                                   `json:"url" api:"required"`
	FullRequest    bool                                                     `json:"full_request"`
	RequestHeaders []SandboxBoxNewResponseProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           sandboxBoxNewResponseProxyConfigCallbackJSON             `json:"-"`
}

// sandboxBoxNewResponseProxyConfigCallbackJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseProxyConfigCallback]
type sandboxBoxNewResponseProxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigCallbacksRequestHeader struct {
	Name  string                                                      `json:"name" api:"required"`
	Type  SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                                        `json:"is_set"`
	Value string                                                      `json:"value"`
	JSON  sandboxBoxNewResponseProxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// sandboxBoxNewResponseProxyConfigCallbacksRequestHeaderJSON contains the JSON
// metadata for the struct [SandboxBoxNewResponseProxyConfigCallbacksRequestHeader]
type sandboxBoxNewResponseProxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxNewResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewResponseProxyConfigRule struct {
	MatchHosts []string                                      `json:"match_hosts" api:"required"`
	Name       string                                        `json:"name" api:"required"`
	Enabled    bool                                          `json:"enabled"`
	Headers    []SandboxBoxNewResponseProxyConfigRulesHeader `json:"headers"`
	MatchPaths []string                                      `json:"match_paths"`
	JSON       sandboxBoxNewResponseProxyConfigRuleJSON      `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseProxyConfigRule]
type sandboxBoxNewResponseProxyConfigRuleJSON struct {
	MatchHosts  apijson.Field
	Name        apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	MatchPaths  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesHeader struct {
	Name  string                                           `json:"name" api:"required"`
	Type  SandboxBoxNewResponseProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                             `json:"is_set"`
	Value string                                           `json:"value"`
	JSON  sandboxBoxNewResponseProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRulesHeaderJSON contains the JSON metadata for
// the struct [SandboxBoxNewResponseProxyConfigRulesHeader]
type sandboxBoxNewResponseProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesHeadersType string

const (
	SandboxBoxNewResponseProxyConfigRulesHeadersTypePlaintext       SandboxBoxNewResponseProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxNewResponseProxyConfigRulesHeadersTypeOpaque          SandboxBoxNewResponseProxyConfigRulesHeadersType = "opaque"
	SandboxBoxNewResponseProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxNewResponseProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxNewResponseProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseProxyConfigRulesHeadersTypePlaintext, SandboxBoxNewResponseProxyConfigRulesHeadersTypeOpaque, SandboxBoxNewResponseProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponse struct {
	ID                     string                           `json:"id"`
	CreatedAt              string                           `json:"created_at"`
	CreatedBy              string                           `json:"created_by"`
	DataplaneURL           string                           `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                            `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                            `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                            `json:"idle_ttl_seconds"`
	MemBytes               int64                            `json:"mem_bytes"`
	Name                   string                           `json:"name"`
	ProxyConfig            SandboxBoxGetResponseProxyConfig `json:"proxy_config"`
	SizeClass              string                           `json:"size_class"`
	SnapshotID             string                           `json:"snapshot_id"`
	Status                 string                           `json:"status"`
	StatusMessage          string                           `json:"status_message"`
	StoppedAt              string                           `json:"stopped_at"`
	UpdatedAt              string                           `json:"updated_at"`
	UpdatedBy              string                           `json:"updated_by"`
	Vcpus                  int64                            `json:"vcpus"`
	JSON                   sandboxBoxGetResponseJSON        `json:"-"`
}

// sandboxBoxGetResponseJSON contains the JSON metadata for the struct
// [SandboxBoxGetResponse]
type sandboxBoxGetResponseJSON struct {
	ID                     apijson.Field
	CreatedAt              apijson.Field
	CreatedBy              apijson.Field
	DataplaneURL           apijson.Field
	DeleteAfterStopSeconds apijson.Field
	FsCapacityBytes        apijson.Field
	IdleTtlSeconds         apijson.Field
	MemBytes               apijson.Field
	Name                   apijson.Field
	ProxyConfig            apijson.Field
	SizeClass              apijson.Field
	SnapshotID             apijson.Field
	Status                 apijson.Field
	StatusMessage          apijson.Field
	StoppedAt              apijson.Field
	UpdatedAt              apijson.Field
	UpdatedBy              apijson.Field
	Vcpus                  apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfig struct {
	AccessControl SandboxBoxGetResponseProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxBoxGetResponseProxyConfigCallback    `json:"callbacks"`
	NoProxy       []string                                      `json:"no_proxy"`
	Rules         []SandboxBoxGetResponseProxyConfigRule        `json:"rules"`
	JSON          sandboxBoxGetResponseProxyConfigJSON          `json:"-"`
}

// sandboxBoxGetResponseProxyConfigJSON contains the JSON metadata for the struct
// [SandboxBoxGetResponseProxyConfig]
type sandboxBoxGetResponseProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigAccessControl struct {
	AllowList []string                                          `json:"allow_list"`
	DenyList  []string                                          `json:"deny_list"`
	JSON      sandboxBoxGetResponseProxyConfigAccessControlJSON `json:"-"`
}

// sandboxBoxGetResponseProxyConfigAccessControlJSON contains the JSON metadata for
// the struct [SandboxBoxGetResponseProxyConfigAccessControl]
type sandboxBoxGetResponseProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigCallback struct {
	MatchHosts     []string                                                 `json:"match_hosts" api:"required"`
	TtlSeconds     int64                                                    `json:"ttl_seconds" api:"required"`
	URL            string                                                   `json:"url" api:"required"`
	FullRequest    bool                                                     `json:"full_request"`
	RequestHeaders []SandboxBoxGetResponseProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           sandboxBoxGetResponseProxyConfigCallbackJSON             `json:"-"`
}

// sandboxBoxGetResponseProxyConfigCallbackJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseProxyConfigCallback]
type sandboxBoxGetResponseProxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigCallbacksRequestHeader struct {
	Name  string                                                      `json:"name" api:"required"`
	Type  SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                                        `json:"is_set"`
	Value string                                                      `json:"value"`
	JSON  sandboxBoxGetResponseProxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// sandboxBoxGetResponseProxyConfigCallbacksRequestHeaderJSON contains the JSON
// metadata for the struct [SandboxBoxGetResponseProxyConfigCallbacksRequestHeader]
type sandboxBoxGetResponseProxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxGetResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponseProxyConfigRule struct {
	MatchHosts []string                                      `json:"match_hosts" api:"required"`
	Name       string                                        `json:"name" api:"required"`
	Enabled    bool                                          `json:"enabled"`
	Headers    []SandboxBoxGetResponseProxyConfigRulesHeader `json:"headers"`
	MatchPaths []string                                      `json:"match_paths"`
	JSON       sandboxBoxGetResponseProxyConfigRuleJSON      `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseProxyConfigRule]
type sandboxBoxGetResponseProxyConfigRuleJSON struct {
	MatchHosts  apijson.Field
	Name        apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	MatchPaths  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesHeader struct {
	Name  string                                           `json:"name" api:"required"`
	Type  SandboxBoxGetResponseProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                             `json:"is_set"`
	Value string                                           `json:"value"`
	JSON  sandboxBoxGetResponseProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRulesHeaderJSON contains the JSON metadata for
// the struct [SandboxBoxGetResponseProxyConfigRulesHeader]
type sandboxBoxGetResponseProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesHeadersType string

const (
	SandboxBoxGetResponseProxyConfigRulesHeadersTypePlaintext       SandboxBoxGetResponseProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxGetResponseProxyConfigRulesHeadersTypeOpaque          SandboxBoxGetResponseProxyConfigRulesHeadersType = "opaque"
	SandboxBoxGetResponseProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxGetResponseProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxGetResponseProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseProxyConfigRulesHeadersTypePlaintext, SandboxBoxGetResponseProxyConfigRulesHeadersTypeOpaque, SandboxBoxGetResponseProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponse struct {
	ID                     string                              `json:"id"`
	CreatedAt              string                              `json:"created_at"`
	CreatedBy              string                              `json:"created_by"`
	DataplaneURL           string                              `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                               `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                               `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                               `json:"idle_ttl_seconds"`
	MemBytes               int64                               `json:"mem_bytes"`
	Name                   string                              `json:"name"`
	ProxyConfig            SandboxBoxUpdateResponseProxyConfig `json:"proxy_config"`
	SizeClass              string                              `json:"size_class"`
	SnapshotID             string                              `json:"snapshot_id"`
	Status                 string                              `json:"status"`
	StatusMessage          string                              `json:"status_message"`
	StoppedAt              string                              `json:"stopped_at"`
	UpdatedAt              string                              `json:"updated_at"`
	UpdatedBy              string                              `json:"updated_by"`
	Vcpus                  int64                               `json:"vcpus"`
	JSON                   sandboxBoxUpdateResponseJSON        `json:"-"`
}

// sandboxBoxUpdateResponseJSON contains the JSON metadata for the struct
// [SandboxBoxUpdateResponse]
type sandboxBoxUpdateResponseJSON struct {
	ID                     apijson.Field
	CreatedAt              apijson.Field
	CreatedBy              apijson.Field
	DataplaneURL           apijson.Field
	DeleteAfterStopSeconds apijson.Field
	FsCapacityBytes        apijson.Field
	IdleTtlSeconds         apijson.Field
	MemBytes               apijson.Field
	Name                   apijson.Field
	ProxyConfig            apijson.Field
	SizeClass              apijson.Field
	SnapshotID             apijson.Field
	Status                 apijson.Field
	StatusMessage          apijson.Field
	StoppedAt              apijson.Field
	UpdatedAt              apijson.Field
	UpdatedBy              apijson.Field
	Vcpus                  apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfig struct {
	AccessControl SandboxBoxUpdateResponseProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxBoxUpdateResponseProxyConfigCallback    `json:"callbacks"`
	NoProxy       []string                                         `json:"no_proxy"`
	Rules         []SandboxBoxUpdateResponseProxyConfigRule        `json:"rules"`
	JSON          sandboxBoxUpdateResponseProxyConfigJSON          `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigJSON contains the JSON metadata for the
// struct [SandboxBoxUpdateResponseProxyConfig]
type sandboxBoxUpdateResponseProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigAccessControl struct {
	AllowList []string                                             `json:"allow_list"`
	DenyList  []string                                             `json:"deny_list"`
	JSON      sandboxBoxUpdateResponseProxyConfigAccessControlJSON `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigAccessControlJSON contains the JSON metadata
// for the struct [SandboxBoxUpdateResponseProxyConfigAccessControl]
type sandboxBoxUpdateResponseProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigCallback struct {
	MatchHosts     []string                                                    `json:"match_hosts" api:"required"`
	TtlSeconds     int64                                                       `json:"ttl_seconds" api:"required"`
	URL            string                                                      `json:"url" api:"required"`
	FullRequest    bool                                                        `json:"full_request"`
	RequestHeaders []SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           sandboxBoxUpdateResponseProxyConfigCallbackJSON             `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigCallbackJSON contains the JSON metadata for
// the struct [SandboxBoxUpdateResponseProxyConfigCallback]
type sandboxBoxUpdateResponseProxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeader struct {
	Name  string                                                         `json:"name" api:"required"`
	Type  SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                                           `json:"is_set"`
	Value string                                                         `json:"value"`
	JSON  sandboxBoxUpdateResponseProxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigCallbacksRequestHeaderJSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeader]
type sandboxBoxUpdateResponseProxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxUpdateResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseProxyConfigRule struct {
	MatchHosts []string                                         `json:"match_hosts" api:"required"`
	Name       string                                           `json:"name" api:"required"`
	Enabled    bool                                             `json:"enabled"`
	Headers    []SandboxBoxUpdateResponseProxyConfigRulesHeader `json:"headers"`
	MatchPaths []string                                         `json:"match_paths"`
	JSON       sandboxBoxUpdateResponseProxyConfigRuleJSON      `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxUpdateResponseProxyConfigRule]
type sandboxBoxUpdateResponseProxyConfigRuleJSON struct {
	MatchHosts  apijson.Field
	Name        apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	MatchPaths  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesHeader struct {
	Name  string                                              `json:"name" api:"required"`
	Type  SandboxBoxUpdateResponseProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                                `json:"is_set"`
	Value string                                              `json:"value"`
	JSON  sandboxBoxUpdateResponseProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRulesHeaderJSON contains the JSON metadata
// for the struct [SandboxBoxUpdateResponseProxyConfigRulesHeader]
type sandboxBoxUpdateResponseProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesHeadersType string

const (
	SandboxBoxUpdateResponseProxyConfigRulesHeadersTypePlaintext       SandboxBoxUpdateResponseProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxUpdateResponseProxyConfigRulesHeadersTypeOpaque          SandboxBoxUpdateResponseProxyConfigRulesHeadersType = "opaque"
	SandboxBoxUpdateResponseProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxUpdateResponseProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseProxyConfigRulesHeadersTypePlaintext, SandboxBoxUpdateResponseProxyConfigRulesHeadersTypeOpaque, SandboxBoxUpdateResponseProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponse struct {
	Offset    int64                           `json:"offset"`
	Sandboxes []SandboxBoxListResponseSandbox `json:"sandboxes"`
	JSON      sandboxBoxListResponseJSON      `json:"-"`
}

// sandboxBoxListResponseJSON contains the JSON metadata for the struct
// [SandboxBoxListResponse]
type sandboxBoxListResponseJSON struct {
	Offset      apijson.Field
	Sandboxes   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandbox struct {
	ID                     string                                     `json:"id"`
	CreatedAt              string                                     `json:"created_at"`
	CreatedBy              string                                     `json:"created_by"`
	DataplaneURL           string                                     `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                                      `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                                      `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                                      `json:"idle_ttl_seconds"`
	MemBytes               int64                                      `json:"mem_bytes"`
	Name                   string                                     `json:"name"`
	ProxyConfig            SandboxBoxListResponseSandboxesProxyConfig `json:"proxy_config"`
	SizeClass              string                                     `json:"size_class"`
	SnapshotID             string                                     `json:"snapshot_id"`
	Status                 string                                     `json:"status"`
	StatusMessage          string                                     `json:"status_message"`
	StoppedAt              string                                     `json:"stopped_at"`
	UpdatedAt              string                                     `json:"updated_at"`
	UpdatedBy              string                                     `json:"updated_by"`
	Vcpus                  int64                                      `json:"vcpus"`
	JSON                   sandboxBoxListResponseSandboxJSON          `json:"-"`
}

// sandboxBoxListResponseSandboxJSON contains the JSON metadata for the struct
// [SandboxBoxListResponseSandbox]
type sandboxBoxListResponseSandboxJSON struct {
	ID                     apijson.Field
	CreatedAt              apijson.Field
	CreatedBy              apijson.Field
	DataplaneURL           apijson.Field
	DeleteAfterStopSeconds apijson.Field
	FsCapacityBytes        apijson.Field
	IdleTtlSeconds         apijson.Field
	MemBytes               apijson.Field
	Name                   apijson.Field
	ProxyConfig            apijson.Field
	SizeClass              apijson.Field
	SnapshotID             apijson.Field
	Status                 apijson.Field
	StatusMessage          apijson.Field
	StoppedAt              apijson.Field
	UpdatedAt              apijson.Field
	UpdatedBy              apijson.Field
	Vcpus                  apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandbox) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfig struct {
	AccessControl SandboxBoxListResponseSandboxesProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxBoxListResponseSandboxesProxyConfigCallback    `json:"callbacks"`
	NoProxy       []string                                                `json:"no_proxy"`
	Rules         []SandboxBoxListResponseSandboxesProxyConfigRule        `json:"rules"`
	JSON          sandboxBoxListResponseSandboxesProxyConfigJSON          `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigJSON contains the JSON metadata for
// the struct [SandboxBoxListResponseSandboxesProxyConfig]
type sandboxBoxListResponseSandboxesProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigAccessControl struct {
	AllowList []string                                                    `json:"allow_list"`
	DenyList  []string                                                    `json:"deny_list"`
	JSON      sandboxBoxListResponseSandboxesProxyConfigAccessControlJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigAccessControlJSON contains the JSON
// metadata for the struct
// [SandboxBoxListResponseSandboxesProxyConfigAccessControl]
type sandboxBoxListResponseSandboxesProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigCallback struct {
	MatchHosts     []string                                                           `json:"match_hosts" api:"required"`
	TtlSeconds     int64                                                              `json:"ttl_seconds" api:"required"`
	URL            string                                                             `json:"url" api:"required"`
	FullRequest    bool                                                               `json:"full_request"`
	RequestHeaders []SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           sandboxBoxListResponseSandboxesProxyConfigCallbackJSON             `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigCallbackJSON contains the JSON
// metadata for the struct [SandboxBoxListResponseSandboxesProxyConfigCallback]
type sandboxBoxListResponseSandboxesProxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeader struct {
	Name  string                                                                `json:"name" api:"required"`
	Type  SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                                                  `json:"is_set"`
	Value string                                                                `json:"value"`
	JSON  sandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeaderJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeader]
type sandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxListResponseSandboxesProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesProxyConfigRule struct {
	MatchHosts []string                                                `json:"match_hosts" api:"required"`
	Name       string                                                  `json:"name" api:"required"`
	Enabled    bool                                                    `json:"enabled"`
	Headers    []SandboxBoxListResponseSandboxesProxyConfigRulesHeader `json:"headers"`
	MatchPaths []string                                                `json:"match_paths"`
	JSON       sandboxBoxListResponseSandboxesProxyConfigRuleJSON      `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRuleJSON contains the JSON metadata
// for the struct [SandboxBoxListResponseSandboxesProxyConfigRule]
type sandboxBoxListResponseSandboxesProxyConfigRuleJSON struct {
	MatchHosts  apijson.Field
	Name        apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	MatchPaths  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesHeader struct {
	Name  string                                                     `json:"name" api:"required"`
	Type  SandboxBoxListResponseSandboxesProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                                       `json:"is_set"`
	Value string                                                     `json:"value"`
	JSON  sandboxBoxListResponseSandboxesProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRulesHeaderJSON contains the JSON
// metadata for the struct [SandboxBoxListResponseSandboxesProxyConfigRulesHeader]
type sandboxBoxListResponseSandboxesProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesHeadersType string

const (
	SandboxBoxListResponseSandboxesProxyConfigRulesHeadersTypePlaintext       SandboxBoxListResponseSandboxesProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxListResponseSandboxesProxyConfigRulesHeadersTypeOpaque          SandboxBoxListResponseSandboxesProxyConfigRulesHeadersType = "opaque"
	SandboxBoxListResponseSandboxesProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxListResponseSandboxesProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesProxyConfigRulesHeadersTypePlaintext, SandboxBoxListResponseSandboxesProxyConfigRulesHeadersTypeOpaque, SandboxBoxListResponseSandboxesProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewSnapshotResponse struct {
	ID              string                            `json:"id"`
	CreatedAt       string                            `json:"created_at"`
	CreatedBy       string                            `json:"created_by"`
	DockerImage     string                            `json:"docker_image"`
	FsCapacityBytes int64                             `json:"fs_capacity_bytes"`
	FsUsedBytes     int64                             `json:"fs_used_bytes"`
	ImageDigest     string                            `json:"image_digest"`
	Name            string                            `json:"name"`
	RegistryID      string                            `json:"registry_id"`
	SourceSandboxID string                            `json:"source_sandbox_id"`
	Status          string                            `json:"status"`
	StatusMessage   string                            `json:"status_message"`
	UpdatedAt       string                            `json:"updated_at"`
	JSON            sandboxBoxNewSnapshotResponseJSON `json:"-"`
}

// sandboxBoxNewSnapshotResponseJSON contains the JSON metadata for the struct
// [SandboxBoxNewSnapshotResponse]
type sandboxBoxNewSnapshotResponseJSON struct {
	ID              apijson.Field
	CreatedAt       apijson.Field
	CreatedBy       apijson.Field
	DockerImage     apijson.Field
	FsCapacityBytes apijson.Field
	FsUsedBytes     apijson.Field
	ImageDigest     apijson.Field
	Name            apijson.Field
	RegistryID      apijson.Field
	SourceSandboxID apijson.Field
	Status          apijson.Field
	StatusMessage   apijson.Field
	UpdatedAt       apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxNewSnapshotResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewSnapshotResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGenerateServiceURLResponse struct {
	Token      string                                   `json:"token"`
	BrowserURL string                                   `json:"browser_url"`
	ExpiresAt  string                                   `json:"expires_at"`
	ServiceURL string                                   `json:"service_url"`
	JSON       sandboxBoxGenerateServiceURLResponseJSON `json:"-"`
}

// sandboxBoxGenerateServiceURLResponseJSON contains the JSON metadata for the
// struct [SandboxBoxGenerateServiceURLResponse]
type sandboxBoxGenerateServiceURLResponseJSON struct {
	Token       apijson.Field
	BrowserURL  apijson.Field
	ExpiresAt   apijson.Field
	ServiceURL  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGenerateServiceURLResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGenerateServiceURLResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetStatusResponse struct {
	Status        string                          `json:"status"`
	StatusMessage string                          `json:"status_message"`
	JSON          sandboxBoxGetStatusResponseJSON `json:"-"`
}

// sandboxBoxGetStatusResponseJSON contains the JSON metadata for the struct
// [SandboxBoxGetStatusResponse]
type sandboxBoxGetStatusResponseJSON struct {
	Status        apijson.Field
	StatusMessage apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxBoxGetStatusResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetStatusResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponse struct {
	ID                     string                             `json:"id"`
	CreatedAt              string                             `json:"created_at"`
	CreatedBy              string                             `json:"created_by"`
	DataplaneURL           string                             `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                              `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                              `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                              `json:"idle_ttl_seconds"`
	MemBytes               int64                              `json:"mem_bytes"`
	Name                   string                             `json:"name"`
	ProxyConfig            SandboxBoxStartResponseProxyConfig `json:"proxy_config"`
	SizeClass              string                             `json:"size_class"`
	SnapshotID             string                             `json:"snapshot_id"`
	Status                 string                             `json:"status"`
	StatusMessage          string                             `json:"status_message"`
	StoppedAt              string                             `json:"stopped_at"`
	UpdatedAt              string                             `json:"updated_at"`
	UpdatedBy              string                             `json:"updated_by"`
	Vcpus                  int64                              `json:"vcpus"`
	JSON                   sandboxBoxStartResponseJSON        `json:"-"`
}

// sandboxBoxStartResponseJSON contains the JSON metadata for the struct
// [SandboxBoxStartResponse]
type sandboxBoxStartResponseJSON struct {
	ID                     apijson.Field
	CreatedAt              apijson.Field
	CreatedBy              apijson.Field
	DataplaneURL           apijson.Field
	DeleteAfterStopSeconds apijson.Field
	FsCapacityBytes        apijson.Field
	IdleTtlSeconds         apijson.Field
	MemBytes               apijson.Field
	Name                   apijson.Field
	ProxyConfig            apijson.Field
	SizeClass              apijson.Field
	SnapshotID             apijson.Field
	Status                 apijson.Field
	StatusMessage          apijson.Field
	StoppedAt              apijson.Field
	UpdatedAt              apijson.Field
	UpdatedBy              apijson.Field
	Vcpus                  apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfig struct {
	AccessControl SandboxBoxStartResponseProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxBoxStartResponseProxyConfigCallback    `json:"callbacks"`
	NoProxy       []string                                        `json:"no_proxy"`
	Rules         []SandboxBoxStartResponseProxyConfigRule        `json:"rules"`
	JSON          sandboxBoxStartResponseProxyConfigJSON          `json:"-"`
}

// sandboxBoxStartResponseProxyConfigJSON contains the JSON metadata for the struct
// [SandboxBoxStartResponseProxyConfig]
type sandboxBoxStartResponseProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigAccessControl struct {
	AllowList []string                                            `json:"allow_list"`
	DenyList  []string                                            `json:"deny_list"`
	JSON      sandboxBoxStartResponseProxyConfigAccessControlJSON `json:"-"`
}

// sandboxBoxStartResponseProxyConfigAccessControlJSON contains the JSON metadata
// for the struct [SandboxBoxStartResponseProxyConfigAccessControl]
type sandboxBoxStartResponseProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigCallback struct {
	MatchHosts     []string                                                   `json:"match_hosts" api:"required"`
	TtlSeconds     int64                                                      `json:"ttl_seconds" api:"required"`
	URL            string                                                     `json:"url" api:"required"`
	FullRequest    bool                                                       `json:"full_request"`
	RequestHeaders []SandboxBoxStartResponseProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           sandboxBoxStartResponseProxyConfigCallbackJSON             `json:"-"`
}

// sandboxBoxStartResponseProxyConfigCallbackJSON contains the JSON metadata for
// the struct [SandboxBoxStartResponseProxyConfigCallback]
type sandboxBoxStartResponseProxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigCallbacksRequestHeader struct {
	Name  string                                                        `json:"name" api:"required"`
	Type  SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                                          `json:"is_set"`
	Value string                                                        `json:"value"`
	JSON  sandboxBoxStartResponseProxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// sandboxBoxStartResponseProxyConfigCallbacksRequestHeaderJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseProxyConfigCallbacksRequestHeader]
type sandboxBoxStartResponseProxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxStartResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxStartResponseProxyConfigRule struct {
	MatchHosts []string                                        `json:"match_hosts" api:"required"`
	Name       string                                          `json:"name" api:"required"`
	Enabled    bool                                            `json:"enabled"`
	Headers    []SandboxBoxStartResponseProxyConfigRulesHeader `json:"headers"`
	MatchPaths []string                                        `json:"match_paths"`
	JSON       sandboxBoxStartResponseProxyConfigRuleJSON      `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxStartResponseProxyConfigRule]
type sandboxBoxStartResponseProxyConfigRuleJSON struct {
	MatchHosts  apijson.Field
	Name        apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	MatchPaths  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesHeader struct {
	Name  string                                             `json:"name" api:"required"`
	Type  SandboxBoxStartResponseProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                               `json:"is_set"`
	Value string                                             `json:"value"`
	JSON  sandboxBoxStartResponseProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRulesHeaderJSON contains the JSON metadata for
// the struct [SandboxBoxStartResponseProxyConfigRulesHeader]
type sandboxBoxStartResponseProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesHeadersType string

const (
	SandboxBoxStartResponseProxyConfigRulesHeadersTypePlaintext       SandboxBoxStartResponseProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxStartResponseProxyConfigRulesHeadersTypeOpaque          SandboxBoxStartResponseProxyConfigRulesHeadersType = "opaque"
	SandboxBoxStartResponseProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxStartResponseProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxStartResponseProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseProxyConfigRulesHeadersTypePlaintext, SandboxBoxStartResponseProxyConfigRulesHeadersTypeOpaque, SandboxBoxStartResponseProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParams struct {
	DeleteAfterStopSeconds param.Field[int64]                          `json:"delete_after_stop_seconds"`
	FsCapacityBytes        param.Field[int64]                          `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64]                          `json:"idle_ttl_seconds"`
	MemBytes               param.Field[int64]                          `json:"mem_bytes"`
	Name                   param.Field[string]                         `json:"name"`
	ProxyConfig            param.Field[SandboxBoxNewParamsProxyConfig] `json:"proxy_config"`
	SnapshotID             param.Field[string]                         `json:"snapshot_id"`
	SnapshotName           param.Field[string]                         `json:"snapshot_name"`
	Vcpus                  param.Field[int64]                          `json:"vcpus"`
}

func (r SandboxBoxNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfig struct {
	AccessControl param.Field[SandboxBoxNewParamsProxyConfigAccessControl] `json:"access_control"`
	Callbacks     param.Field[[]SandboxBoxNewParamsProxyConfigCallback]    `json:"callbacks"`
	NoProxy       param.Field[[]string]                                    `json:"no_proxy"`
	Rules         param.Field[[]SandboxBoxNewParamsProxyConfigRule]        `json:"rules"`
}

func (r SandboxBoxNewParamsProxyConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigAccessControl struct {
	AllowList param.Field[[]string] `json:"allow_list"`
	DenyList  param.Field[[]string] `json:"deny_list"`
}

func (r SandboxBoxNewParamsProxyConfigAccessControl) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigCallback struct {
	MatchHosts     param.Field[[]string]                                               `json:"match_hosts" api:"required"`
	TtlSeconds     param.Field[int64]                                                  `json:"ttl_seconds" api:"required"`
	URL            param.Field[string]                                                 `json:"url" api:"required"`
	FullRequest    param.Field[bool]                                                   `json:"full_request"`
	RequestHeaders param.Field[[]SandboxBoxNewParamsProxyConfigCallbacksRequestHeader] `json:"request_headers"`
}

func (r SandboxBoxNewParamsProxyConfigCallback) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigCallbacksRequestHeader struct {
	Name  param.Field[string]                                                    `json:"name" api:"required"`
	Type  param.Field[SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                      `json:"is_set"`
	Value param.Field[string]                                                    `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigCallbacksRequestHeader) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRule struct {
	MatchHosts param.Field[[]string]                                    `json:"match_hosts" api:"required"`
	Name       param.Field[string]                                      `json:"name" api:"required"`
	Enabled    param.Field[bool]                                        `json:"enabled"`
	Headers    param.Field[[]SandboxBoxNewParamsProxyConfigRulesHeader] `json:"headers"`
	MatchPaths param.Field[[]string]                                    `json:"match_paths"`
}

func (r SandboxBoxNewParamsProxyConfigRule) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesHeader struct {
	Name  param.Field[string]                                         `json:"name" api:"required"`
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesHeadersType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                           `json:"is_set"`
	Value param.Field[string]                                         `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesHeader) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesHeadersType string

const (
	SandboxBoxNewParamsProxyConfigRulesHeadersTypePlaintext       SandboxBoxNewParamsProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesHeadersTypeOpaque          SandboxBoxNewParamsProxyConfigRulesHeadersType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesHeadersTypePlaintext, SandboxBoxNewParamsProxyConfigRulesHeadersTypeOpaque, SandboxBoxNewParamsProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParams struct {
	DeleteAfterStopSeconds param.Field[int64]                             `json:"delete_after_stop_seconds"`
	FsCapacityBytes        param.Field[int64]                             `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64]                             `json:"idle_ttl_seconds"`
	MemBytes               param.Field[int64]                             `json:"mem_bytes"`
	Name                   param.Field[string]                            `json:"name"`
	ProxyConfig            param.Field[SandboxBoxUpdateParamsProxyConfig] `json:"proxy_config"`
	Vcpus                  param.Field[int64]                             `json:"vcpus"`
}

func (r SandboxBoxUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfig struct {
	AccessControl param.Field[SandboxBoxUpdateParamsProxyConfigAccessControl] `json:"access_control"`
	Callbacks     param.Field[[]SandboxBoxUpdateParamsProxyConfigCallback]    `json:"callbacks"`
	NoProxy       param.Field[[]string]                                       `json:"no_proxy"`
	Rules         param.Field[[]SandboxBoxUpdateParamsProxyConfigRule]        `json:"rules"`
}

func (r SandboxBoxUpdateParamsProxyConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigAccessControl struct {
	AllowList param.Field[[]string] `json:"allow_list"`
	DenyList  param.Field[[]string] `json:"deny_list"`
}

func (r SandboxBoxUpdateParamsProxyConfigAccessControl) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigCallback struct {
	MatchHosts     param.Field[[]string]                                                  `json:"match_hosts" api:"required"`
	TtlSeconds     param.Field[int64]                                                     `json:"ttl_seconds" api:"required"`
	URL            param.Field[string]                                                    `json:"url" api:"required"`
	FullRequest    param.Field[bool]                                                      `json:"full_request"`
	RequestHeaders param.Field[[]SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeader] `json:"request_headers"`
}

func (r SandboxBoxUpdateParamsProxyConfigCallback) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeader struct {
	Name  param.Field[string]                                                       `json:"name" api:"required"`
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                         `json:"is_set"`
	Value param.Field[string]                                                       `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeader) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersType string

const (
	SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRule struct {
	MatchHosts param.Field[[]string]                                       `json:"match_hosts" api:"required"`
	Name       param.Field[string]                                         `json:"name" api:"required"`
	Enabled    param.Field[bool]                                           `json:"enabled"`
	Headers    param.Field[[]SandboxBoxUpdateParamsProxyConfigRulesHeader] `json:"headers"`
	MatchPaths param.Field[[]string]                                       `json:"match_paths"`
}

func (r SandboxBoxUpdateParamsProxyConfigRule) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesHeader struct {
	Name  param.Field[string]                                            `json:"name" api:"required"`
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesHeadersType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                              `json:"is_set"`
	Value param.Field[string]                                            `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesHeader) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesHeadersType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesHeadersTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesHeadersType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesHeadersTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesHeadersType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesHeadersTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesHeadersTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesHeadersTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListParams struct {
	// Maximum number of results
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Pagination offset
	Offset param.Field[int64] `query:"offset"`
	// Sort column (name, status, created_at)
	SortBy param.Field[string] `query:"sort_by"`
	// Sort direction (asc, desc)
	SortDirection param.Field[string] `query:"sort_direction"`
	// Filter by status (provisioning, ready, failed, stopped, deleting)
	Status param.Field[string] `query:"status"`
}

// URLQuery serializes [SandboxBoxListParams]'s query parameters as `url.Values`.
func (r SandboxBoxListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SandboxBoxNewSnapshotParams struct {
	Name param.Field[string] `json:"name" api:"required"`
	// if omitted, creates a fresh checkpoint from the running VM
	Checkpoint param.Field[string] `json:"checkpoint"`
}

func (r SandboxBoxNewSnapshotParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxGenerateServiceURLParams struct {
	ExpiresInSeconds param.Field[int64] `json:"expires_in_seconds"`
	Port             param.Field[int64] `json:"port"`
}

func (r SandboxBoxGenerateServiceURLParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
