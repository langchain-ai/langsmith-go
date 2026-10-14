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
	"github.com/langchain-ai/langsmith-go/packages/pagination"
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

// Create a new sandbox from a snapshot. Provide at most one of `snapshot_id` or
// `snapshot_name`; if neither is provided, the server uses the default snapshot.
// `snapshot_name` accepts a Docker-style `name` or `name:tag` reference (a bare
// name resolves to `name:latest`).
func (r *SandboxBoxService) New(ctx context.Context, body SandboxBoxNewParams, opts ...option.RequestOption) (res *SandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v2/sandboxes/boxes"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a sandbox by name. Stale provisioning sandboxes are auto-failed.
func (r *SandboxBoxService) Get(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a sandbox's display name, retention, resources, tags, or proxy
// configuration. The name must be unique within the tenant. Proxy configuration
// sent to a sandbox that is not running is stored and applied when it next starts.
func (r *SandboxBoxService) Update(ctx context.Context, name string, body SandboxBoxUpdateParams, opts ...option.RequestOption) (res *SandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List sandboxes for the authenticated tenant, with optional filtering, sorting,
// and pagination. Page with page_size and cursor: replay the response's
// next_cursor until it comes back null, which is the only signal that no pages
// remain. Cursors are opaque and only valid on this endpoint; do not parse or
// construct one.
func (r *SandboxBoxService) List(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[SandboxResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v2/sandboxes/boxes"
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List sandboxes for the authenticated tenant, with optional filtering, sorting,
// and pagination. Page with page_size and cursor: replay the response's
// next_cursor until it comes back null, which is the only signal that no pages
// remain. Cursors are opaque and only valid on this endpoint; do not parse or
// construct one.
func (r *SandboxBoxService) ListAutoPaging(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[SandboxResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Delete a sandbox by name or UUID. Tears down the sandbox runtime and removes the
// DB record.
func (r *SandboxBoxService) Delete(ctx context.Context, name string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Create a snapshot by capturing the current state of a sandbox or promoting an
// existing checkpoint.
func (r *SandboxBoxService) NewSnapshot(ctx context.Context, name string, body SandboxBoxNewSnapshotParams, opts ...option.RequestOption) (res *SnapshotResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/snapshot", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Removes the sharing grant for one port, or for every port when port is omitted.
// A LangSmith login URL stops working immediately. A previously minted service
// token is not revoked and stays valid until it expires, but no new one can be
// issued from the removed grant.
func (r *SandboxBoxService) DeleteServiceURL(ctx context.Context, name string, body SandboxBoxDeleteServiceURLParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/service-urls", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return err
}

// Generate a tokenized link that downloads a single file from a sandbox with no
// further authentication. This mints a token rather than creating an addressable
// resource, so it returns 200 with no Location header. The token pins the sandbox,
// the file path, the response content type and disposition, and the sandbox flags,
// so a link cannot be repointed at another file or served under a weaker policy.
// The file is always served with a Content-Security-Policy: a sandbox directive,
// plus a default-src holding every fetch to the file's own download host and a set
// of pre-approved third-party origins. csp_sandbox_flags may loosen the sandbox
// with allow-downloads, allow-forms, allow-modals, allow-orientation-lock,
// allow-pointer-lock, allow-popups, allow-presentation, allow-same-origin,
// allow-scripts, or allow-top-navigation-by-user-activation. Every file is served
// from its own host, derived from the sandbox and the path, so allow-same-origin
// gives a page localStorage and IndexedDB that no other file can read, and
// re-minting a link for the same file keeps them. csp_sandbox set to false drops
// the sandbox directive altogether, and csp_sandbox_flags must then be omitted.
// csp_source_bundles selects the third-party origins: cdnjs, google-fonts,
// jsdelivr, and unpkg are all allowed when the field is omitted, 'none' holds the
// file to its own host, and 'any' sends no default-src at all. Links never expire
// unless expires_in_seconds is set. The link is served from the sandbox service
// domain, not the API host.
func (r *SandboxBoxService) GenerateDownloadURL(ctx context.Context, name string, body SandboxBoxGenerateDownloadURLParams, opts ...option.RequestOption) (res *DownloadURLResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/download-url", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Create a short-lived JWT for accessing an HTTP service running on a specific
// port inside a sandbox. Returns a browser_url (sets auth cookie via redirect), a
// service_url (for use with the X-Langsmith-Sandbox-Service-Token header), the raw
// token, and its expiry. Set access=restricted|workspace to instead enable durable
// LangSmith login (no token; users authenticate with their normal LangSmith
// session), or access=off to disable it. LangSmith login and token access are
// mutually exclusive per service URL.
func (r *SandboxBoxService) GenerateServiceURL(ctx context.Context, name string, body SandboxBoxGenerateServiceURLParams, opts ...option.RequestOption) (res *ServiceURLResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/service-url", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve the lightweight status of a sandbox for polling.
func (r *SandboxBoxService) GetStatus(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxStatusResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/status", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Returns one entry per port the sandbox is currently reachable on, so a caller
// can see what is shared before turning it off. Expired token grants are omitted.
// Cursors are opaque and only valid on this endpoint; do not parse or construct
// one.
func (r *SandboxBoxService) ListServiceURLs(ctx context.Context, name string, query SandboxBoxListServiceURLsParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[SandboxBoxListServiceURLsResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/service-urls", name)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// Returns one entry per port the sandbox is currently reachable on, so a caller
// can see what is shared before turning it off. Expired token grants are omitted.
// Cursors are opaque and only valid on this endpoint; do not parse or construct
// one.
func (r *SandboxBoxService) ListServiceURLsAutoPaging(ctx context.Context, name string, query SandboxBoxListServiceURLsParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[SandboxBoxListServiceURLsResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.ListServiceURLs(ctx, name, query, opts...))
}

// Start a stopped or failed sandbox. This endpoint is not idempotent.
func (r *SandboxBoxService) Start(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/start", name)
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
	path := fmt.Sprintf("api/v2/sandboxes/boxes/%s/stop", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, nil, opts...)
	return err
}

type SandboxBoxListServiceURLsResponse struct {
	// How the port is shared: "token" for a minted service token, or
	// "restricted"/"workspace" for LangSmith login.
	Access    SandboxBoxListServiceURLsResponseAccess `json:"access" api:"required"`
	CreatedAt time.Time                               `json:"created_at" api:"required" format:"date-time"`
	Port      int64                                   `json:"port" api:"required"`
	// The LangSmith user who first shared this port, when known.
	CreatedBy string `json:"created_by" format:"uuid"`
	// When the share expires. Set only for "token"; a login grant does not expire.
	ExpiresAt time.Time                             `json:"expires_at" format:"date-time"`
	JSON      sandboxBoxListServiceURLsResponseJSON `json:"-"`
}

// sandboxBoxListServiceURLsResponseJSON contains the JSON metadata for the struct
// [SandboxBoxListServiceURLsResponse]
type sandboxBoxListServiceURLsResponseJSON struct {
	Access      apijson.Field
	CreatedAt   apijson.Field
	Port        apijson.Field
	CreatedBy   apijson.Field
	ExpiresAt   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListServiceURLsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListServiceURLsResponseJSON) RawJSON() string {
	return r.raw
}

// How the port is shared: "token" for a minted service token, or
// "restricted"/"workspace" for LangSmith login.
type SandboxBoxListServiceURLsResponseAccess string

const (
	SandboxBoxListServiceURLsResponseAccessToken      SandboxBoxListServiceURLsResponseAccess = "token"
	SandboxBoxListServiceURLsResponseAccessRestricted SandboxBoxListServiceURLsResponseAccess = "restricted"
	SandboxBoxListServiceURLsResponseAccessWorkspace  SandboxBoxListServiceURLsResponseAccess = "workspace"
)

func (r SandboxBoxListServiceURLsResponseAccess) IsKnown() bool {
	switch r {
	case SandboxBoxListServiceURLsResponseAccessToken, SandboxBoxListServiceURLsResponseAccessRestricted, SandboxBoxListServiceURLsResponseAccessWorkspace:
		return true
	}
	return false
}

type SandboxBoxNewParams struct {
	// AccessDelegation lets code inside the sandbox call the LangSmith API as you,
	// with at most the permissions granted here. Omit for no access.
	AccessDelegation param.Field[SandboxBoxNewParamsAccessDelegation] `json:"access_delegation"`
	// CPUMillicores optionally requests CPU at millicore granularity (e.g. 500 = 0.5
	// vCPU); takes precedence over VCPUs. Fractional (sub-vCPU) values are not
	// available for every sandbox.
	CPUMillicores          param.Field[int64]             `json:"cpu_millicores"`
	DeleteAfterStopSeconds param.Field[int64]             `json:"delete_after_stop_seconds"`
	EnvVars                param.Field[map[string]string] `json:"env_vars"`
	FsCapacityBytes        param.Field[int64]             `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64]             `json:"idle_ttl_seconds"`
	// Labels are free-form key/value metadata persisted with the sandbox and returned
	// on reads. Labels from the source snapshot are inherited unless overridden here.
	Labels param.Field[map[string]string] `json:"labels"`
	// Memory for the sandbox, in bytes. Memory is tied to CPU at 4 GiB per vCPU: omit
	// it and it follows that ratio; set it and it must stay within 50% of the ratio
	// for the requested CPU, so a 1 vCPU sandbox accepts 2-6 GiB. Setting memory
	// without CPU derives the CPU from the same ratio. Maximum 64 GiB.
	MemBytes    param.Field[int64]                          `json:"mem_bytes"`
	MountConfig param.Field[SandboxBoxNewParamsMountConfig] `json:"mount_config"`
	Name        param.Field[string]                         `json:"name"`
	// PreserveMemoryOnStop, when true, suspends the sandbox's memory on a voluntary
	// stop (idle timeout or explicit stop) so the next start resumes from where it
	// left off. Default false discards memory and keeps only the filesystem, so the
	// next start is a cold boot. Restarts triggered by infrastructure maintenance
	// always preserve memory regardless of this setting.
	PreserveMemoryOnStop param.Field[bool]                           `json:"preserve_memory_on_stop"`
	ProxyConfig          param.Field[SandboxBoxNewParamsProxyConfig] `json:"proxy_config"`
	// RestoreMemory selects how the sandbox handles a snapshot's captured memory:
	//
	// nil → if-present: resume from memory when the snapshot has it, else cold-boot
	// (default). true → always: resume from memory; rejected if the snapshot has none.
	// false → never: always cold-boot.
	//
	// Applies to this request only.
	RestoreMemory param.Field[bool] `json:"restore_memory"`
	// RunConfig overrides the snapshot's run config for this sandbox: user and
	// work_dir replace the snapshot's, env_vars merge over it. The result is what the
	// sandbox boots with, and what a snapshot captured from it carries.
	RunConfig param.Field[SandboxBoxNewParamsRunConfig] `json:"run_config"`
	// Snapshot is a Docker-style name or name:tag reference to boot from. A bare name
	// resolves to name:latest.
	Snapshot   param.Field[string] `json:"snapshot"`
	SnapshotID param.Field[string] `json:"snapshot_id"`
	// SnapshotName is a synonym for Snapshot, accepted for compatibility with clients
	// that predate it. Set one or the other.
	SnapshotName param.Field[string]   `json:"snapshot_name"`
	TagValueIDs  param.Field[[]string] `json:"tag_value_ids"`
	Vcpus        param.Field[int64]    `json:"vcpus"`
}

func (r SandboxBoxNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// AccessDelegation lets code inside the sandbox call the LangSmith API as you,
// with at most the permissions granted here. Omit for no access.
type SandboxBoxNewParamsAccessDelegation struct {
	Mode        param.Field[SandboxBoxNewParamsAccessDelegationMode] `json:"mode" api:"required"`
	Permissions param.Field[[]string]                                `json:"permissions"`
}

func (r SandboxBoxNewParamsAccessDelegation) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsAccessDelegationMode string

const (
	SandboxBoxNewParamsAccessDelegationModeInherit  SandboxBoxNewParamsAccessDelegationMode = "INHERIT"
	SandboxBoxNewParamsAccessDelegationModeExplicit SandboxBoxNewParamsAccessDelegationMode = "EXPLICIT"
)

func (r SandboxBoxNewParamsAccessDelegationMode) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsAccessDelegationModeInherit, SandboxBoxNewParamsAccessDelegationModeExplicit:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfig struct {
	Auth   param.Field[SandboxBoxNewParamsMountConfigAuth]         `json:"auth"`
	Mounts param.Field[[]SandboxBoxNewParamsMountConfigMountUnion] `json:"mounts"`
}

func (r SandboxBoxNewParamsMountConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuth struct {
	Aws   param.Field[SandboxBoxNewParamsMountConfigAuthAwsUnion]   `json:"aws"`
	Azure param.Field[SandboxBoxNewParamsMountConfigAuthAzureUnion] `json:"azure"`
	Gcp   param.Field[SandboxBoxNewParamsMountConfigAuthGcpUnion]   `json:"gcp"`
}

func (r SandboxBoxNewParamsMountConfigAuth) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAws struct {
	AccessKeyID param.Field[interface{}] `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn         param.Field[string]      `json:"role_arn"`
	SecretAccessKey param.Field[interface{}] `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsMountConfigAuthAws) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAws) implementsSandboxBoxNewParamsMountConfigAuthAwsUnion() {
}

// Satisfied by
// [SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthAws].
type SandboxBoxNewParamsMountConfigAuthAwsUnion interface {
	implementsSandboxBoxNewParamsMountConfigAuthAwsUnion()
}

type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn param.Field[string] `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthAwsUnion() {
}

type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig struct {
	AccessKeyID     param.Field[SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey] `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn param.Field[SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn] `json:"role_arn"`
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthAwsUnion() {
}

type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                         `json:"is_set"`
	Value param.Field[string]                                                                                       `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType string

const (
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypePlaintext       SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeOpaque          SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "opaque"
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypePlaintext, SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeOpaque, SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                             `json:"is_set"`
	Value param.Field[string]                                                                                           `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType string

const (
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypePlaintext       SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeOpaque          SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "opaque"
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypePlaintext, SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeOpaque, SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace; the role's trust policy must trust
// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
// be changed after creation.
type SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn string

const (
	SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArnEmpty SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn = ""
)

func (r SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArnEmpty:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigAuthAzure struct {
	AccountKey param.Field[interface{}] `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID param.Field[string]      `json:"client_id"`
	SasToken param.Field[interface{}] `json:"sas_token"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for storage access tokens of client_id. Cannot be changed after
	// creation.
	TenantID param.Field[string] `json:"tenant_id"`
}

func (r SandboxBoxNewParamsMountConfigAuthAzure) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAzure) implementsSandboxBoxNewParamsMountConfigAuthAzureUnion() {
}

// Satisfied by
// [SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthAzure].
type SandboxBoxNewParamsMountConfigAuthAzureUnion interface {
	implementsSandboxBoxNewParamsMountConfigAuthAzureUnion()
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig struct {
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID param.Field[string] `json:"client_id" api:"required"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for storage access tokens of client_id. Cannot be changed after
	// creation.
	TenantID param.Field[string] `json:"tenant_id" api:"required"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthAzureUnion() {
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig struct {
	AccountKey param.Field[SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey] `json:"account_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthAzureUnion() {
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                                `json:"is_set"`
	Value param.Field[string]                                                                                              `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType string

const (
	SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypePlaintext       SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeOpaque          SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "opaque"
	SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypePlaintext, SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeOpaque, SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig struct {
	SasToken param.Field[SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken] `json:"sas_token" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthAzureUnion() {
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                       `json:"is_set"`
	Value param.Field[string]                                                                                     `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType string

const (
	SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypePlaintext       SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeOpaque          SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "opaque"
	SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypePlaintext, SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeOpaque, SandboxBoxNewParamsMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigAuthGcp struct {
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly. Only valid with
	// workload_identity_provider.
	ServiceAccountEmail param.Field[string]      `json:"service_account_email"`
	ServiceAccountJson  param.Field[interface{}] `json:"service_account_json"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope param.Field[string] `json:"subject_scope"`
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Mutually exclusive with service_account_json. Cannot be changed after creation.
	WorkloadIdentityProvider param.Field[string] `json:"workload_identity_provider"`
}

func (r SandboxBoxNewParamsMountConfigAuthGcp) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthGcp) implementsSandboxBoxNewParamsMountConfigAuthGcpUnion() {
}

// Satisfied by
// [SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig],
// [SandboxBoxNewParamsMountConfigAuthGcp].
type SandboxBoxNewParamsMountConfigAuthGcpUnion interface {
	implementsSandboxBoxNewParamsMountConfigAuthGcpUnion()
}

type SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig struct {
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Mutually exclusive with service_account_json. Cannot be changed after creation.
	WorkloadIdentityProvider param.Field[string] `json:"workload_identity_provider" api:"required"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly. Only valid with
	// workload_identity_provider.
	ServiceAccountEmail param.Field[string] `json:"service_account_email"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthGcpUnion() {
}

type SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig struct {
	ServiceAccountJson param.Field[SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson] `json:"service_account_json" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig) implementsSandboxBoxNewParamsMountConfigAuthGcpUnion() {
}

type SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                                        `json:"is_set"`
	Value param.Field[string]                                                                                                      `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType string

const (
	SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypePlaintext       SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeOpaque          SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "opaque"
	SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypePlaintext, SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeOpaque, SandboxBoxNewParamsMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMount struct {
	ID         param.Field[string]                                   `json:"id" api:"required"`
	MountPath  param.Field[string]                                   `json:"mount_path" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsType] `json:"type" api:"required"`
	AzureBlob  param.Field[interface{}]                              `json:"azure_blob"`
	Cache      param.Field[interface{}]                              `json:"cache"`
	Contexthub param.Field[interface{}]                              `json:"contexthub"`
	Encryption param.Field[interface{}]                              `json:"encryption"`
	Gcs        param.Field[interface{}]                              `json:"gcs"`
	Git        param.Field[interface{}]                              `json:"git"`
	ReadOnly   param.Field[bool]                                     `json:"read_only"`
	S3         param.Field[interface{}]                              `json:"s3"`
}

func (r SandboxBoxNewParamsMountConfigMount) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMount) implementsSandboxBoxNewParamsMountConfigMountUnion() {}

// Satisfied by [SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpec],
// [SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpec],
// [SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpec],
// [SandboxBoxNewParamsMountConfigMount].
type SandboxBoxNewParamsMountConfigMountUnion interface {
	implementsSandboxBoxNewParamsMountConfigMountUnion()
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID         param.Field[string]                                                                    `json:"id" api:"required"`
	MountPath  param.Field[string]                                                                    `json:"mount_path" api:"required"`
	S3         param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3]         `json:"s3" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType]       `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob]  `json:"azure_blob"`
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecContexthub] `json:"contexthub"`
	Encryption param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecEncryption] `json:"encryption"`
	Gcs        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGcs]        `json:"gcs"`
	Git        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGit]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                      `json:"read_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxNewParamsMountConfigMountUnion() {
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeS3         SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs        SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeAzureBlob  SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "azure_blob"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeAzureBlob, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecContexthub) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecEncryption) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              param.Field[string]                                                                `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                 `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name param.Field[string]                                                                    `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID         param.Field[string]                                                                     `json:"id" api:"required"`
	Gcs        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGcs]        `json:"gcs" api:"required"`
	MountPath  param.Field[string]                                                                     `json:"mount_path" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType]       `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob]  `json:"azure_blob"`
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecContexthub] `json:"contexthub"`
	Encryption param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecEncryption] `json:"encryption"`
	Git        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGit]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                       `json:"read_only"`
	S3         param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecS3]         `json:"s3"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxNewParamsMountConfigMountUnion() {
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3         SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs        SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeAzureBlob  SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "azure_blob"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeAzureBlob, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecContexthub) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecEncryption) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              param.Field[string]                                                                 `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                  `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name param.Field[string]                                                                     `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID         param.Field[string]                                                                   `json:"id" api:"required"`
	Git        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGit]        `json:"git" api:"required"`
	MountPath  param.Field[string]                                                                   `json:"mount_path" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType]       `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob]  `json:"azure_blob"`
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecContexthub] `json:"contexthub"`
	Encryption param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecEncryption] `json:"encryption"`
	Gcs        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGcs]        `json:"gcs"`
	ReadOnly   param.Field[bool]                                                                     `json:"read_only"`
	S3         param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecS3]         `json:"s3"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxNewParamsMountConfigMountUnion() {
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              param.Field[string]                                                               `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name param.Field[string]                                                                   `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeS3         SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs        SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeAzureBlob  SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "azure_blob"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeAzureBlob, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecContexthub) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecEncryption) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpec struct {
	ID         param.Field[string]                                                                          `json:"id" api:"required"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub] `json:"contexthub" api:"required"`
	MountPath  param.Field[string]                                                                          `json:"mount_path" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType]       `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob]  `json:"azure_blob"`
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecCache]      `json:"cache"`
	Encryption param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption] `json:"encryption"`
	Gcs        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGcs]        `json:"gcs"`
	Git        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGit]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                            `json:"read_only"`
	S3         param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecS3]         `json:"s3"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpec) implementsSandboxBoxNewParamsMountConfigMountUnion() {
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3         SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs        SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeAzureBlob  SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "azure_blob"
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeAzureBlob, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGit struct {
	RemoteURL              param.Field[string]                                                                      `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                       `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef struct {
	Name param.Field[string]                                                                          `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeBranch SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeTag    SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpec struct {
	ID         param.Field[string]                                                                              `json:"id" api:"required"`
	AzureBlob  param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob]  `json:"azure_blob" api:"required"`
	MountPath  param.Field[string]                                                                              `json:"mount_path" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType]       `json:"type" api:"required"`
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub] `json:"contexthub"`
	Encryption param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption] `json:"encryption"`
	Gcs        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs]        `json:"gcs"`
	Git        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                                `json:"read_only"`
	S3         param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3]         `json:"s3"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpec) implementsSandboxBoxNewParamsMountConfigMountUnion() {
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeS3         SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGcs        SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeAzureBlob  SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "azure_blob"
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeAzureBlob, SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit struct {
	RemoteURL              param.Field[string]                                                                          `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                           `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef struct {
	Name param.Field[string]                                                                              `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeBranch SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeTag    SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsType string

const (
	SandboxBoxNewParamsMountConfigMountsTypeS3         SandboxBoxNewParamsMountConfigMountsType = "s3"
	SandboxBoxNewParamsMountConfigMountsTypeGcs        SandboxBoxNewParamsMountConfigMountsType = "gcs"
	SandboxBoxNewParamsMountConfigMountsTypeAzureBlob  SandboxBoxNewParamsMountConfigMountsType = "azure_blob"
	SandboxBoxNewParamsMountConfigMountsTypeGit        SandboxBoxNewParamsMountConfigMountsType = "git"
	SandboxBoxNewParamsMountConfigMountsTypeContexthub SandboxBoxNewParamsMountConfigMountsType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsTypeS3, SandboxBoxNewParamsMountConfigMountsTypeGcs, SandboxBoxNewParamsMountConfigMountsTypeAzureBlob, SandboxBoxNewParamsMountConfigMountsTypeGit, SandboxBoxNewParamsMountConfigMountsTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfig struct {
	AccessControl param.Field[SandboxBoxNewParamsProxyConfigAccessControl] `json:"access_control"`
	Callbacks     param.Field[[]SandboxBoxNewParamsProxyConfigCallback]    `json:"callbacks"`
	// Description says what this configuration as a whole lets the sandbox reach,
	// complementing the per-rule descriptions. At most 1024 characters.
	Description param.Field[string]                               `json:"description"`
	NoProxy     param.Field[[]string]                             `json:"no_proxy"`
	Rules       param.Field[[]SandboxBoxNewParamsProxyConfigRule] `json:"rules"`
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
	Name  param.Field[string]                                        `json:"name" api:"required"`
	Aws   param.Field[SandboxBoxNewParamsProxyConfigRulesAwsUnion]   `json:"aws"`
	Azure param.Field[SandboxBoxNewParamsProxyConfigRulesAzureUnion] `json:"azure"`
	// Description says what this rule lets the sandbox reach, so an agent driving the
	// sandbox can be told its capabilities. At most 1024 characters.
	Description param.Field[string] `json:"description"`
	Enabled     param.Field[bool]   `json:"enabled"`
	// EnvVars are plaintext env vars set for every command in the sandbox while this
	// rule is enabled. Use them for tools that refuse to run unless a credential env
	// var is present (e.g. gh needs GH_TOKEN) even though this rule injects the real
	// credential on the wire — set a dummy value here so the command starts. Explicit
	// per-sandbox env_vars win over these, and provider-managed (AWS/GCP) vars win
	// over both.
	EnvVars param.Field[map[string]string]                           `json:"env_vars"`
	Gcp     param.Field[SandboxBoxNewParamsProxyConfigRulesGcpUnion] `json:"gcp"`
	Headers param.Field[[]SandboxBoxNewParamsProxyConfigRulesHeader] `json:"headers"`
	// MatchHeaders restricts a header injection rule to requests carrying every listed
	// header, each written "name: value" with a lowercase name and an exact value.
	// Pair with headers of the same name to swap a placeholder the sandbox sends (e.g.
	// "authorization: Bearer account-b") for a real credential, so one host can serve
	// several accounts. Rules are evaluated in order and the first match wins.
	MatchHeaders param.Field[[]string] `json:"match_headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts param.Field[[]string] `json:"match_hosts"`
	MatchPaths param.Field[[]string] `json:"match_paths"`
	Type       param.Field[string]   `json:"type"`
}

func (r SandboxBoxNewParamsProxyConfigRule) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAws struct {
	AccessKeyID param.Field[interface{}] `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn         param.Field[string]      `json:"role_arn"`
	SecretAccessKey param.Field[interface{}] `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAws) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAws) implementsSandboxBoxNewParamsProxyConfigRulesAwsUnion() {
}

// Satisfied by
// [SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig],
// [SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig],
// [SandboxBoxNewParamsProxyConfigRulesAws].
type SandboxBoxNewParamsProxyConfigRulesAwsUnion interface {
	implementsSandboxBoxNewParamsProxyConfigRulesAwsUnion()
}

type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn param.Field[string] `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) implementsSandboxBoxNewParamsProxyConfigRulesAwsUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig struct {
	AccessKeyID     param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey] `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn] `json:"role_arn"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) implementsSandboxBoxNewParamsProxyConfigRulesAwsUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                               `json:"is_set"`
	Value param.Field[string]                                                                             `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType string

const (
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext       SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque          SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext, SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque, SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                   `json:"is_set"`
	Value param.Field[string]                                                                                 `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType string

const (
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext       SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque          SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext, SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque, SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace. The role's trust policy must trust
// LangSmith's issuer for the token subject. Credentials renew automatically, and
// access follows the role's effective AWS permissions. Mutually exclusive with
// static keys, and cannot be changed after creation.
type SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn string

const (
	SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn = ""
)

func (r SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesAzure struct {
	AccountKey param.Field[interface{}] `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID  param.Field[string]      `json:"client_id"`
	Resources param.Field[interface{}] `json:"resources"`
	SasToken  param.Field[interface{}] `json:"sas_token"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for access tokens of client_id, so no secret is stored. Cannot be
	// changed after creation.
	TenantID param.Field[string] `json:"tenant_id"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAzure) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAzure) implementsSandboxBoxNewParamsProxyConfigRulesAzureUnion() {
}

// Satisfied by
// [SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig],
// [SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig],
// [SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig],
// [SandboxBoxNewParamsProxyConfigRulesAzure].
type SandboxBoxNewParamsProxyConfigRulesAzureUnion interface {
	implementsSandboxBoxNewParamsProxyConfigRulesAzureUnion()
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig struct {
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID param.Field[string] `json:"client_id" api:"required"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for access tokens of client_id, so no secret is stored. Cannot be
	// changed after creation.
	TenantID param.Field[string] `json:"tenant_id" api:"required"`
	// Resources narrows the Entra resources access tokens are obtained for, from
	// https://storage.azure.com, https://vault.azure.net, https://management.azure.com
	// and https://cognitiveservices.azure.com. Defaults to all of them. Only valid
	// with tenant_id.
	Resources param.Field[[]SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource] `json:"resources"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) implementsSandboxBoxNewParamsProxyConfigRulesAzureUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource string

const (
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom           SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://storage.azure.com"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet             SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://vault.azure.net"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom        SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://management.azure.com"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://cognitiveservices.azure.com"
)

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig struct {
	AccountKey param.Field[SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey] `json:"account_key" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) implementsSandboxBoxNewParamsProxyConfigRulesAzureUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                      `json:"is_set"`
	Value param.Field[string]                                                                                    `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType string

const (
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext       SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque          SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig struct {
	SasToken param.Field[SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken] `json:"sas_token" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig) implementsSandboxBoxNewParamsProxyConfigRulesAzureUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                             `json:"is_set"`
	Value param.Field[string]                                                                           `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType string

const (
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext       SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque          SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque, SandboxBoxNewParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesGcp struct {
	Scopes param.Field[interface{}] `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail param.Field[string]      `json:"service_account_email"`
	ServiceAccountJson  param.Field[interface{}] `json:"service_account_json"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope param.Field[string] `json:"subject_scope"`
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key.
	// Use the full resource name:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Cannot be changed after creation.
	WorkloadIdentityProvider param.Field[string] `json:"workload_identity_provider"`
}

func (r SandboxBoxNewParamsProxyConfigRulesGcp) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesGcp) implementsSandboxBoxNewParamsProxyConfigRulesGcpUnion() {
}

// Satisfied by
// [SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig],
// [SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig],
// [SandboxBoxNewParamsProxyConfigRulesGcp].
type SandboxBoxNewParamsProxyConfigRulesGcpUnion interface {
	implementsSandboxBoxNewParamsProxyConfigRulesGcpUnion()
}

type SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig struct {
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key.
	// Use the full resource name:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Cannot be changed after creation.
	WorkloadIdentityProvider param.Field[string] `json:"workload_identity_provider" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes param.Field[[]string] `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail param.Field[string] `json:"service_account_email"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) implementsSandboxBoxNewParamsProxyConfigRulesGcpUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig struct {
	ServiceAccountJson param.Field[SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson] `json:"service_account_json" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes param.Field[[]string] `json:"scopes"`
}

func (r SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) implementsSandboxBoxNewParamsProxyConfigRulesGcpUnion() {
}

type SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                              `json:"is_set"`
	Value param.Field[string]                                                                                            `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType string

const (
	SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext       SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque          SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext, SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque, SandboxBoxNewParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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

// RunConfig overrides the snapshot's run config for this sandbox: user and
// work_dir replace the snapshot's, env_vars merge over it. The result is what the
// sandbox boots with, and what a snapshot captured from it carries.
type SandboxBoxNewParamsRunConfig struct {
	EnvVars param.Field[map[string]string] `json:"env_vars"`
	User    param.Field[string]            `json:"user"`
	WorkDir param.Field[string]            `json:"work_dir"`
}

func (r SandboxBoxNewParamsRunConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParams struct {
	CPUMillicores          param.Field[int64] `json:"cpu_millicores"`
	DeleteAfterStopSeconds param.Field[int64] `json:"delete_after_stop_seconds"`
	FsCapacityBytes        param.Field[int64] `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64] `json:"idle_ttl_seconds"`
	// New memory for the sandbox, in bytes. The 4 GiB per vCPU ratio applies when the
	// sandbox is created; a resize enforces only the maximum of 64 GiB.
	MemBytes    param.Field[int64]                             `json:"mem_bytes"`
	Name        param.Field[string]                            `json:"name"`
	ProxyConfig param.Field[SandboxBoxUpdateParamsProxyConfig] `json:"proxy_config"`
	// RunConfig changes what subsequent commands run with: user and work_dir replace
	// the current values, env_vars merge over them. Commands already running are
	// unaffected.
	RunConfig   param.Field[SandboxBoxUpdateParamsRunConfig] `json:"run_config"`
	TagValueIDs param.Field[[]string]                        `json:"tag_value_ids"`
	Vcpus       param.Field[int64]                           `json:"vcpus"`
}

func (r SandboxBoxUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfig struct {
	AccessControl param.Field[SandboxBoxUpdateParamsProxyConfigAccessControl] `json:"access_control"`
	Callbacks     param.Field[[]SandboxBoxUpdateParamsProxyConfigCallback]    `json:"callbacks"`
	// Description says what this configuration as a whole lets the sandbox reach,
	// complementing the per-rule descriptions. At most 1024 characters.
	Description param.Field[string]                                  `json:"description"`
	NoProxy     param.Field[[]string]                                `json:"no_proxy"`
	Rules       param.Field[[]SandboxBoxUpdateParamsProxyConfigRule] `json:"rules"`
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
	Name  param.Field[string]                                           `json:"name" api:"required"`
	Aws   param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsUnion]   `json:"aws"`
	Azure param.Field[SandboxBoxUpdateParamsProxyConfigRulesAzureUnion] `json:"azure"`
	// Description says what this rule lets the sandbox reach, so an agent driving the
	// sandbox can be told its capabilities. At most 1024 characters.
	Description param.Field[string] `json:"description"`
	Enabled     param.Field[bool]   `json:"enabled"`
	// EnvVars are plaintext env vars set for every command in the sandbox while this
	// rule is enabled. Use them for tools that refuse to run unless a credential env
	// var is present (e.g. gh needs GH_TOKEN) even though this rule injects the real
	// credential on the wire — set a dummy value here so the command starts. Explicit
	// per-sandbox env_vars win over these, and provider-managed (AWS/GCP) vars win
	// over both.
	EnvVars param.Field[map[string]string]                              `json:"env_vars"`
	Gcp     param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcpUnion] `json:"gcp"`
	Headers param.Field[[]SandboxBoxUpdateParamsProxyConfigRulesHeader] `json:"headers"`
	// MatchHeaders restricts a header injection rule to requests carrying every listed
	// header, each written "name: value" with a lowercase name and an exact value.
	// Pair with headers of the same name to swap a placeholder the sandbox sends (e.g.
	// "authorization: Bearer account-b") for a real credential, so one host can serve
	// several accounts. Rules are evaluated in order and the first match wins.
	MatchHeaders param.Field[[]string] `json:"match_headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts param.Field[[]string] `json:"match_hosts"`
	MatchPaths param.Field[[]string] `json:"match_paths"`
	Type       param.Field[string]   `json:"type"`
}

func (r SandboxBoxUpdateParamsProxyConfigRule) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAws struct {
	AccessKeyID param.Field[interface{}] `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn         param.Field[string]      `json:"role_arn"`
	SecretAccessKey param.Field[interface{}] `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAws) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAws) implementsSandboxBoxUpdateParamsProxyConfigRulesAwsUnion() {
}

// Satisfied by
// [SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesAws].
type SandboxBoxUpdateParamsProxyConfigRulesAwsUnion interface {
	implementsSandboxBoxUpdateParamsProxyConfigRulesAwsUnion()
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn param.Field[string] `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesAwsUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig struct {
	AccessKeyID     param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey] `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn] `json:"role_arn"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesAwsUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                  `json:"is_set"`
	Value param.Field[string]                                                                                `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                      `json:"is_set"`
	Value param.Field[string]                                                                                    `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace. The role's trust policy must trust
// LangSmith's issuer for the token subject. Credentials renew automatically, and
// access follows the role's effective AWS permissions. Mutually exclusive with
// static keys, and cannot be changed after creation.
type SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn = ""
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesAzure struct {
	AccountKey param.Field[interface{}] `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID  param.Field[string]      `json:"client_id"`
	Resources param.Field[interface{}] `json:"resources"`
	SasToken  param.Field[interface{}] `json:"sas_token"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for access tokens of client_id, so no secret is stored. Cannot be
	// changed after creation.
	TenantID param.Field[string] `json:"tenant_id"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzure) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzure) implementsSandboxBoxUpdateParamsProxyConfigRulesAzureUnion() {
}

// Satisfied by
// [SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesAzure].
type SandboxBoxUpdateParamsProxyConfigRulesAzureUnion interface {
	implementsSandboxBoxUpdateParamsProxyConfigRulesAzureUnion()
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig struct {
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID param.Field[string] `json:"client_id" api:"required"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for access tokens of client_id, so no secret is stored. Cannot be
	// changed after creation.
	TenantID param.Field[string] `json:"tenant_id" api:"required"`
	// Resources narrows the Entra resources access tokens are obtained for, from
	// https://storage.azure.com, https://vault.azure.net, https://management.azure.com
	// and https://cognitiveservices.azure.com. Defaults to all of them. Only valid
	// with tenant_id.
	Resources param.Field[[]SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource] `json:"resources"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesAzureUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom           SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://storage.azure.com"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet             SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://vault.azure.net"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom        SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://management.azure.com"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://cognitiveservices.azure.com"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig struct {
	AccountKey param.Field[SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey] `json:"account_key" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesAzureUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                         `json:"is_set"`
	Value param.Field[string]                                                                                       `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig struct {
	SasToken param.Field[SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken] `json:"sas_token" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesAzureUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                `json:"is_set"`
	Value param.Field[string]                                                                              `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesGcp struct {
	Scopes param.Field[interface{}] `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail param.Field[string]      `json:"service_account_email"`
	ServiceAccountJson  param.Field[interface{}] `json:"service_account_json"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope param.Field[string] `json:"subject_scope"`
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key.
	// Use the full resource name:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Cannot be changed after creation.
	WorkloadIdentityProvider param.Field[string] `json:"workload_identity_provider"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcp) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcp) implementsSandboxBoxUpdateParamsProxyConfigRulesGcpUnion() {
}

// Satisfied by
// [SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig],
// [SandboxBoxUpdateParamsProxyConfigRulesGcp].
type SandboxBoxUpdateParamsProxyConfigRulesGcpUnion interface {
	implementsSandboxBoxUpdateParamsProxyConfigRulesGcpUnion()
}

type SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig struct {
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key.
	// Use the full resource name:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Cannot be changed after creation.
	WorkloadIdentityProvider param.Field[string] `json:"workload_identity_provider" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes param.Field[[]string] `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail param.Field[string] `json:"service_account_email"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesGcpUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig struct {
	ServiceAccountJson param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson] `json:"service_account_json" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes param.Field[[]string] `json:"scopes"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) implementsSandboxBoxUpdateParamsProxyConfigRulesGcpUnion() {
}

type SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                                 `json:"is_set"`
	Value param.Field[string]                                                                                               `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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

// RunConfig changes what subsequent commands run with: user and work_dir replace
// the current values, env_vars merge over them. Commands already running are
// unaffected.
type SandboxBoxUpdateParamsRunConfig struct {
	EnvVars param.Field[map[string]string] `json:"env_vars"`
	User    param.Field[string]            `json:"user"`
	WorkDir param.Field[string]            `json:"work_dir"`
}

func (r SandboxBoxUpdateParamsRunConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxListParams struct {
	// Filter by creator identity. Only 'me' is supported.
	CreatedBy param.Field[string] `query:"created_by"`
	// Opaque pagination cursor from a prior response's next_cursor
	Cursor param.Field[string] `query:"cursor"`
	// Filter by label. Repeatable; all must match. Use 'key' to match on key presence
	// or 'key=value' for equality.
	Label param.Field[[]string] `query:"label"`
	// Deprecated: use page_size. Maximum number of results
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Deprecated: use cursor. Pagination offset
	Offset param.Field[int64] `query:"offset"`
	// Number of results per page
	PageSize param.Field[int64] `query:"page_size"`
	// Sort column (name, status, created_at, stopped_at, idle_ttl_seconds,
	// delete_after_stop_seconds)
	SortBy param.Field[string] `query:"sort_by"`
	// Deprecated: use sort_order. Sort direction (asc, desc)
	SortDirection param.Field[string] `query:"sort_direction"`
	// Sort direction (asc, desc)
	SortOrder param.Field[string] `query:"sort_order"`
	// Filter by status (provisioning, ready, failed, stopped, deleting)
	Status param.Field[string] `query:"status"`
	// Filter by workspace resource tag value IDs; all must match
	TagValueID param.Field[[]string] `query:"tag_value_id"`
}

// URLQuery serializes [SandboxBoxListParams]'s query parameters as `url.Values`.
func (r SandboxBoxListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SandboxBoxNewSnapshotParams struct {
	Name param.Field[string] `json:"name" api:"required"`
	// if omitted, creates a fresh checkpoint from the running VM
	Checkpoint param.Field[string] `json:"checkpoint"`
	// Description says what this snapshot's image can do, so a caller can hand it to
	// an agent as a capability summary. At most 1024 characters.
	Description param.Field[string] `json:"description"`
	// sandbox-local Docker image to export
	DockerImage param.Field[string] `json:"docker_image"`
	// required for Docker image export unless the sandbox has a capacity
	FsCapacityBytes param.Field[int64] `json:"fs_capacity_bytes"`
	// IncludeMemory, when true, captures a full VM memory snapshot alongside the
	// filesystem clone. Only honored when the sandbox is running AND Checkpoint is
	// omitted (i.e. a fresh in-VM checkpoint is requested). Defaults to false to keep
	// snapshots small unless memory restore is explicitly desired.
	IncludeMemory param.Field[bool] `json:"include_memory"`
	// Labels seed the captured snapshot's labels.
	Labels param.Field[map[string]string] `json:"labels"`
	// RunConfig overrides the runtime configuration the snapshot carries: for a
	// docker_image export, the image's USER, WORKDIR and ENV; for a capture of the
	// running VM, the sandbox's own. user and work_dir replace, env_vars merge.
	RunConfig param.Field[SandboxBoxNewSnapshotParamsRunConfig] `json:"run_config"`
	// mutable Docker-style tag; defaults to "latest"
	Tag param.Field[string] `json:"tag"`
}

func (r SandboxBoxNewSnapshotParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// RunConfig overrides the runtime configuration the snapshot carries: for a
// docker_image export, the image's USER, WORKDIR and ENV; for a capture of the
// running VM, the sandbox's own. user and work_dir replace, env_vars merge.
type SandboxBoxNewSnapshotParamsRunConfig struct {
	EnvVars param.Field[map[string]string] `json:"env_vars"`
	User    param.Field[string]            `json:"user"`
	WorkDir param.Field[string]            `json:"work_dir"`
}

func (r SandboxBoxNewSnapshotParamsRunConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxDeleteServiceURLParams struct {
	// Port to stop sharing. Omit to stop sharing every port.
	Port param.Field[int64] `query:"port"`
}

// URLQuery serializes [SandboxBoxDeleteServiceURLParams]'s query parameters as
// `url.Values`.
func (r SandboxBoxDeleteServiceURLParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SandboxBoxGenerateDownloadURLParams struct {
	Path               param.Field[string] `json:"path" api:"required"`
	ContentDisposition param.Field[string] `json:"content_disposition"`
	ContentType        param.Field[string] `json:"content_type"`
	// CSPSandbox false serves the file with no CSP sandbox directive; omit to keep it.
	CspSandbox param.Field[bool] `json:"csp_sandbox"`
	// CSPSandboxFlags loosen the CSP sandbox the file is served under; omit for the
	// most restrictive policy.
	CspSandboxFlags param.Field[[]SandboxBoxGenerateDownloadURLParamsCspSandboxFlag] `json:"csp_sandbox_flags"`
	// CSPSourceBundles allow the served file to fetch from named third-party origins;
	// omit to send no fetch directive.
	CspSourceBundles param.Field[[]SandboxBoxGenerateDownloadURLParamsCspSourceBundle] `json:"csp_source_bundles"`
	// ExpiresInSeconds is optional; a link with no expiry never expires.
	ExpiresInSeconds param.Field[int64] `json:"expires_in_seconds"`
}

func (r SandboxBoxGenerateDownloadURLParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxGenerateDownloadURLParamsCspSandboxFlag string

const (
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowDownloads                     SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-downloads"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowForms                         SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-forms"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowModals                        SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-modals"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowOrientationLock               SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-orientation-lock"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowPointerLock                   SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-pointer-lock"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowPopups                        SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-popups"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowPresentation                  SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-presentation"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowScripts                       SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-scripts"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowTopNavigationByUserActivation SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-top-navigation-by-user-activation"
	SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowSameOrigin                    SandboxBoxGenerateDownloadURLParamsCspSandboxFlag = "allow-same-origin"
)

func (r SandboxBoxGenerateDownloadURLParamsCspSandboxFlag) IsKnown() bool {
	switch r {
	case SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowDownloads, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowForms, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowModals, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowOrientationLock, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowPointerLock, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowPopups, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowPresentation, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowScripts, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowTopNavigationByUserActivation, SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowSameOrigin:
		return true
	}
	return false
}

type SandboxBoxGenerateDownloadURLParamsCspSourceBundle string

const (
	SandboxBoxGenerateDownloadURLParamsCspSourceBundleCdnjs       SandboxBoxGenerateDownloadURLParamsCspSourceBundle = "cdnjs"
	SandboxBoxGenerateDownloadURLParamsCspSourceBundleGoogleFonts SandboxBoxGenerateDownloadURLParamsCspSourceBundle = "google-fonts"
	SandboxBoxGenerateDownloadURLParamsCspSourceBundleJsdelivr    SandboxBoxGenerateDownloadURLParamsCspSourceBundle = "jsdelivr"
	SandboxBoxGenerateDownloadURLParamsCspSourceBundleUnpkg       SandboxBoxGenerateDownloadURLParamsCspSourceBundle = "unpkg"
	SandboxBoxGenerateDownloadURLParamsCspSourceBundleNone        SandboxBoxGenerateDownloadURLParamsCspSourceBundle = "none"
	SandboxBoxGenerateDownloadURLParamsCspSourceBundleAny         SandboxBoxGenerateDownloadURLParamsCspSourceBundle = "any"
)

func (r SandboxBoxGenerateDownloadURLParamsCspSourceBundle) IsKnown() bool {
	switch r {
	case SandboxBoxGenerateDownloadURLParamsCspSourceBundleCdnjs, SandboxBoxGenerateDownloadURLParamsCspSourceBundleGoogleFonts, SandboxBoxGenerateDownloadURLParamsCspSourceBundleJsdelivr, SandboxBoxGenerateDownloadURLParamsCspSourceBundleUnpkg, SandboxBoxGenerateDownloadURLParamsCspSourceBundleNone, SandboxBoxGenerateDownloadURLParamsCspSourceBundleAny:
		return true
	}
	return false
}

type SandboxBoxGenerateServiceURLParams struct {
	// Access selects the login mode, mutually exclusive with the minted token. Omit
	// the field for token mode: mint a short-lived service token (default).
	// "restricted" — LangSmith login: the sandbox's creator, or any user with
	// SandboxesExec on it (admins by default). "workspace" — LangSmith login: any
	// member of the owning workspace. "off" — remove an existing LangSmith login grant
	// and mint a token. A LangSmith login grant is durable; token mode is refused
	// (409) while one exists.
	Access           param.Field[SandboxBoxGenerateServiceURLParamsAccess] `json:"access"`
	ExpiresInSeconds param.Field[int64]                                    `json:"expires_in_seconds"`
	Port             param.Field[int64]                                    `json:"port"`
}

func (r SandboxBoxGenerateServiceURLParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Access selects the login mode, mutually exclusive with the minted token. Omit
// the field for token mode: mint a short-lived service token (default).
// "restricted" — LangSmith login: the sandbox's creator, or any user with
// SandboxesExec on it (admins by default). "workspace" — LangSmith login: any
// member of the owning workspace. "off" — remove an existing LangSmith login grant
// and mint a token. A LangSmith login grant is durable; token mode is refused
// (409) while one exists.
type SandboxBoxGenerateServiceURLParamsAccess string

const (
	SandboxBoxGenerateServiceURLParamsAccessRestricted SandboxBoxGenerateServiceURLParamsAccess = "restricted"
	SandboxBoxGenerateServiceURLParamsAccessWorkspace  SandboxBoxGenerateServiceURLParamsAccess = "workspace"
	SandboxBoxGenerateServiceURLParamsAccessOff        SandboxBoxGenerateServiceURLParamsAccess = "off"
)

func (r SandboxBoxGenerateServiceURLParamsAccess) IsKnown() bool {
	switch r {
	case SandboxBoxGenerateServiceURLParamsAccessRestricted, SandboxBoxGenerateServiceURLParamsAccessWorkspace, SandboxBoxGenerateServiceURLParamsAccessOff:
		return true
	}
	return false
}

type SandboxBoxListServiceURLsParams struct {
	// Opaque pagination cursor from a prior response's next_cursor
	Cursor param.Field[string] `query:"cursor"`
	// Number of results per page
	PageSize param.Field[int64] `query:"page_size"`
}

// URLQuery serializes [SandboxBoxListServiceURLsParams]'s query parameters as
// `url.Values`.
func (r SandboxBoxListServiceURLsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
