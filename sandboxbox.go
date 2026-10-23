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
// session), or access=off to disable it. With LangSmith login, embed_origins names
// the pages that may show the URL in an iframe; a viewer signs in from inside the
// frame through a LangSmith popup. LangSmith login and token access are mutually
// exclusive per service URL.
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
	// Origins that may embed a LangSmith login URL in an iframe; omitted when there
	// are none.
	EmbedOrigins []string `json:"embed_origins"`
	// When the share expires. Set only for "token"; a login grant does not expire.
	ExpiresAt time.Time                             `json:"expires_at" format:"date-time"`
	JSON      sandboxBoxListServiceURLsResponseJSON `json:"-"`
}

// sandboxBoxListServiceURLsResponseJSON contains the JSON metadata for the struct
// [SandboxBoxListServiceURLsResponse]
type sandboxBoxListServiceURLsResponseJSON struct {
	Access       apijson.Field
	CreatedAt    apijson.Field
	Port         apijson.Field
	CreatedBy    apijson.Field
	EmbedOrigins apijson.Field
	ExpiresAt    apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
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
	MemBytes    param.Field[int64]                   `json:"mem_bytes"`
	MountConfig param.Field[SandboxMountConfigParam] `json:"mount_config"`
	Name        param.Field[string]                  `json:"name"`
	// PreserveMemoryOnStop, when true, suspends the sandbox's memory on a voluntary
	// stop (idle timeout or explicit stop) so the next start resumes from where it
	// left off. Default false discards memory and keeps only the filesystem, so the
	// next start is a cold boot. Restarts triggered by infrastructure maintenance
	// always preserve memory regardless of this setting.
	PreserveMemoryOnStop param.Field[bool]             `json:"preserve_memory_on_stop"`
	ProxyConfig          param.Field[ProxyConfigParam] `json:"proxy_config"`
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
	// Volumes (beta) attaches existing volumes at mount paths inside the sandbox. At
	// most 4; requires sandbox volumes to be enabled for your workspace. Attaching a
	// volume requires being its creator or holding sandboxes:exec.
	Volumes param.Field[[]SandboxBoxNewParamsVolume] `json:"volumes"`
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

type SandboxBoxNewParamsVolume struct {
	// MountPath is the absolute directory inside the sandbox where the volume appears.
	MountPath param.Field[string] `json:"mount_path" api:"required"`
	// Volume is the volume's ID or name.
	Volume param.Field[string] `json:"volume" api:"required"`
	// ReadOnly mounts the volume read-only; many sandboxes may do so at once.
	ReadOnly param.Field[bool] `json:"read_only"`
}

func (r SandboxBoxNewParamsVolume) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParams struct {
	CPUMillicores          param.Field[int64] `json:"cpu_millicores"`
	DeleteAfterStopSeconds param.Field[int64] `json:"delete_after_stop_seconds"`
	FsCapacityBytes        param.Field[int64] `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64] `json:"idle_ttl_seconds"`
	// New memory for the sandbox, in bytes. The 4 GiB per vCPU ratio applies when the
	// sandbox is created; a resize enforces only the maximum of 64 GiB.
	MemBytes    param.Field[int64]            `json:"mem_bytes"`
	Name        param.Field[string]           `json:"name"`
	ProxyConfig param.Field[ProxyConfigParam] `json:"proxy_config"`
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
	// Must be true when the included settings hold literal secret values (opaque
	// values rather than workspace secret references, and any env var value),
	// acknowledging that the snapshot stores them and every sandbox created from it
	// can use them.
	ConfirmAllowSnapshotSecretValues param.Field[bool] `json:"confirm_allow_snapshot_secret_values"`
	// Description says what this snapshot's image can do, so a caller can hand it to
	// an agent as a capability summary. At most 1024 characters.
	Description param.Field[string] `json:"description"`
	// sandbox-local Docker image to export
	DockerImage param.Field[string] `json:"docker_image"`
	// required for Docker image export unless the sandbox has a capacity
	FsCapacityBytes param.Field[int64] `json:"fs_capacity_bytes"`
	// ImageConfig sets the command and ports the snapshot records, as a Docker image's
	// ENTRYPOINT, CMD and EXPOSE do. By default a capture keeps those of the snapshot
	// the sandbox was created from; each field given here replaces that one, and an
	// empty list clears it. A function serving the snapshot runs this command and
	// listens on its one exposed port.
	ImageConfig param.Field[SandboxBoxNewSnapshotParamsImageConfig] `json:"image_config"`
	// Sandbox settings to carry into the snapshot. Sandboxes created from it start
	// with each carried setting their create request leaves unset: "mounts" (the
	// mount_config, with its auth), "proxy_config", "env_vars", "resources" (CPU and
	// memory), "lifecycle" (idle_ttl_seconds, delete_after_stop_seconds,
	// preserve_memory_on_stop), and "labels" (merged under this request's labels;
	// reserved com.langchain.\* labels are left out). "all" names every one of these.
	// VM memory is captured with include_memory, not here. Warning: mounts,
	// proxy_config and env_vars carry the sandbox's credentials, and every sandbox
	// created from the snapshot can use them.
	Include param.Field[[]SandboxBoxNewSnapshotParamsInclude] `json:"include"`
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

// ImageConfig sets the command and ports the snapshot records, as a Docker image's
// ENTRYPOINT, CMD and EXPOSE do. By default a capture keeps those of the snapshot
// the sandbox was created from; each field given here replaces that one, and an
// empty list clears it. A function serving the snapshot runs this command and
// listens on its one exposed port.
type SandboxBoxNewSnapshotParamsImageConfig struct {
	// Cmd follows Entrypoint, as CMD: its arguments, or the whole command when
	// Entrypoint is empty.
	Cmd param.Field[[]string] `json:"cmd"`
	// Entrypoint is the program the command starts with, as ENTRYPOINT.
	Entrypoint param.Field[[]string] `json:"entrypoint"`
	// ExposedPorts are the TCP ports the command listens on, as EXPOSE.
	ExposedPorts param.Field[[]int64] `json:"exposed_ports"`
}

func (r SandboxBoxNewSnapshotParamsImageConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewSnapshotParamsInclude string

const (
	SandboxBoxNewSnapshotParamsIncludeAll         SandboxBoxNewSnapshotParamsInclude = "all"
	SandboxBoxNewSnapshotParamsIncludeMounts      SandboxBoxNewSnapshotParamsInclude = "mounts"
	SandboxBoxNewSnapshotParamsIncludeProxyConfig SandboxBoxNewSnapshotParamsInclude = "proxy_config"
	SandboxBoxNewSnapshotParamsIncludeEnvVars     SandboxBoxNewSnapshotParamsInclude = "env_vars"
	SandboxBoxNewSnapshotParamsIncludeResources   SandboxBoxNewSnapshotParamsInclude = "resources"
	SandboxBoxNewSnapshotParamsIncludeLifecycle   SandboxBoxNewSnapshotParamsInclude = "lifecycle"
	SandboxBoxNewSnapshotParamsIncludeLabels      SandboxBoxNewSnapshotParamsInclude = "labels"
)

func (r SandboxBoxNewSnapshotParamsInclude) IsKnown() bool {
	switch r {
	case SandboxBoxNewSnapshotParamsIncludeAll, SandboxBoxNewSnapshotParamsIncludeMounts, SandboxBoxNewSnapshotParamsIncludeProxyConfig, SandboxBoxNewSnapshotParamsIncludeEnvVars, SandboxBoxNewSnapshotParamsIncludeResources, SandboxBoxNewSnapshotParamsIncludeLifecycle, SandboxBoxNewSnapshotParamsIncludeLabels:
		return true
	}
	return false
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
	Access param.Field[SandboxBoxGenerateServiceURLParamsAccess] `json:"access"`
	// EmbedOrigins are the pages allowed to show a LangSmith login service URL in an
	// iframe, besides LangSmith itself: at most 5 origins such as
	// https://app.example.com (http only for localhost). Valid only with access
	// "restricted" or "workspace". Each call sets the URL's whole list, so omitting
	// the field removes any origins set before.
	EmbedOrigins     param.Field[[]string] `json:"embed_origins"`
	ExpiresInSeconds param.Field[int64]    `json:"expires_in_seconds"`
	Port             param.Field[int64]    `json:"port"`
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
