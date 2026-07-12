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

// Create a new sandbox from a snapshot. Provide at most one of `snapshot_id` or
// `snapshot_name`; if neither is provided, the server uses the default snapshot.
func (r *SandboxBoxService) New(ctx context.Context, body SandboxBoxNewParams, opts ...option.RequestOption) (res *SandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/boxes"
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
	path := fmt.Sprintf("v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update a sandbox's display name. The name must be unique within the tenant.
func (r *SandboxBoxService) Update(ctx context.Context, name string, body SandboxBoxUpdateParams, opts ...option.RequestOption) (res *SandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List sandboxes for the authenticated tenant, with optional filtering, sorting,
// and pagination.
func (r *SandboxBoxService) List(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) (res *SandboxListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/boxes"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
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
	path := fmt.Sprintf("v2/sandboxes/boxes/%s", name)
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
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/snapshot", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Create a short-lived JWT for accessing an HTTP service running on a specific
// port inside a sandbox. Returns a browser_url (sets auth cookie via redirect), a
// service_url (for use with the X-Langsmith-Sandbox-Service-Token header), the raw
// token, and its expiry.
func (r *SandboxBoxService) GenerateServiceURL(ctx context.Context, name string, body SandboxBoxGenerateServiceURLParams, opts ...option.RequestOption) (res *ServiceURLResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/service-url", name)
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
	path := fmt.Sprintf("v2/sandboxes/boxes/%s/status", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Start a stopped or failed sandbox. This endpoint is not idempotent.
func (r *SandboxBoxService) Start(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxResponse, err error) {
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

type SandboxBoxNewParams struct {
	// CPUMillicores optionally requests CPU at millicore granularity (e.g. 500 = 0.5
	// vCPU); takes precedence over VCPUs. Fractional (sub-vCPU) values are not
	// available for every sandbox.
	CPUMillicores          param.Field[int64]                          `json:"cpu_millicores"`
	DeleteAfterStopSeconds param.Field[int64]                          `json:"delete_after_stop_seconds"`
	EnvVars                param.Field[map[string]string]              `json:"env_vars"`
	FsCapacityBytes        param.Field[int64]                          `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64]                          `json:"idle_ttl_seconds"`
	MemBytes               param.Field[int64]                          `json:"mem_bytes"`
	MountConfig            param.Field[SandboxBoxNewParamsMountConfig] `json:"mount_config"`
	Name                   param.Field[string]                         `json:"name"`
	ProxyConfig            param.Field[SandboxBoxNewParamsProxyConfig] `json:"proxy_config"`
	// RestoreMemory selects how the sandbox handles a snapshot's captured memory:
	//
	// nil → if-present: resume from memory when the snapshot has it, else cold-boot
	// (default). true → always: resume from memory; rejected if the snapshot has none.
	// false → never: always cold-boot.
	//
	// Applies to this request only.
	RestoreMemory param.Field[bool]     `json:"restore_memory"`
	SnapshotID    param.Field[string]   `json:"snapshot_id"`
	SnapshotName  param.Field[string]   `json:"snapshot_name"`
	TagValueIDs   param.Field[[]string] `json:"tag_value_ids"`
	Vcpus         param.Field[int64]    `json:"vcpus"`
}

func (r SandboxBoxNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfig struct {
	Auth   param.Field[SandboxBoxNewParamsMountConfigAuth]         `json:"auth"`
	Mounts param.Field[[]SandboxBoxNewParamsMountConfigMountUnion] `json:"mounts"`
}

func (r SandboxBoxNewParamsMountConfig) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuth struct {
	Aws param.Field[SandboxBoxNewParamsMountConfigAuthAws] `json:"aws"`
	Gcp param.Field[SandboxBoxNewParamsMountConfigAuthGcp] `json:"gcp"`
}

func (r SandboxBoxNewParamsMountConfigAuth) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAws struct {
	AccessKeyID     param.Field[SandboxBoxNewParamsMountConfigAuthAwsAccessKeyID]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKey] `json:"secret_access_key" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigAuthAws) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAwsAccessKeyID struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                 `json:"is_set"`
	Value param.Field[string]                                               `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthAwsAccessKeyID) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxBoxNewParamsMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKey struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                     `json:"is_set"`
	Value param.Field[string]                                                   `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxBoxNewParamsMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigAuthGcp struct {
	ServiceAccountJson param.Field[SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJson] `json:"service_account_json" api:"required"`
}

func (r SandboxBoxNewParamsMountConfigAuthGcp) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJson struct {
	Type  param.Field[SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                        `json:"is_set"`
	Value param.Field[string]                                                      `json:"value"`
}

func (r SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJson) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxBoxNewParamsMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMount struct {
	ID         param.Field[string]                                   `json:"id" api:"required"`
	MountPath  param.Field[string]                                   `json:"mount_path" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsType] `json:"type" api:"required"`
	Cache      param.Field[interface{}]                              `json:"cache"`
	Contexthub param.Field[interface{}]                              `json:"contexthub"`
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
// [SandboxBoxNewParamsMountConfigMount].
type SandboxBoxNewParamsMountConfigMountUnion interface {
	implementsSandboxBoxNewParamsMountConfigMountUnion()
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID         param.Field[string]                                                                    `json:"id" api:"required"`
	MountPath  param.Field[string]                                                                    `json:"mount_path" api:"required"`
	S3         param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3]         `json:"s3" api:"required"`
	Type       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType]       `json:"type" api:"required"`
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecContexthub] `json:"contexthub"`
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
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub:
		return true
	}
	return false
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
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecContexthub] `json:"contexthub"`
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
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub:
		return true
	}
	return false
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
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecCache]      `json:"cache"`
	Contexthub param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecContexthub] `json:"contexthub"`
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
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub:
		return true
	}
	return false
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
	Cache      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecCache]      `json:"cache"`
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
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit        SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "git"
	SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit, SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecCache) MarshalJSON() (data []byte, err error) {
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

type SandboxBoxNewParamsMountConfigMountsType string

const (
	SandboxBoxNewParamsMountConfigMountsTypeS3         SandboxBoxNewParamsMountConfigMountsType = "s3"
	SandboxBoxNewParamsMountConfigMountsTypeGcs        SandboxBoxNewParamsMountConfigMountsType = "gcs"
	SandboxBoxNewParamsMountConfigMountsTypeGit        SandboxBoxNewParamsMountConfigMountsType = "git"
	SandboxBoxNewParamsMountConfigMountsTypeContexthub SandboxBoxNewParamsMountConfigMountsType = "contexthub"
)

func (r SandboxBoxNewParamsMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsTypeS3, SandboxBoxNewParamsMountConfigMountsTypeGcs, SandboxBoxNewParamsMountConfigMountsTypeGit, SandboxBoxNewParamsMountConfigMountsTypeContexthub:
		return true
	}
	return false
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
	Name    param.Field[string]                                      `json:"name" api:"required"`
	Aws     param.Field[SandboxBoxNewParamsProxyConfigRulesAws]      `json:"aws"`
	Enabled param.Field[bool]                                        `json:"enabled"`
	Gcp     param.Field[SandboxBoxNewParamsProxyConfigRulesGcp]      `json:"gcp"`
	Headers param.Field[[]SandboxBoxNewParamsProxyConfigRulesHeader] `json:"headers"`
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
	AccessKeyID     param.Field[SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyID]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKey] `json:"secret_access_key" api:"required"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAws) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyID struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                  `json:"is_set"`
	Value param.Field[string]                                                `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyID) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKey struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                      `json:"is_set"`
	Value param.Field[string]                                                    `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewParamsProxyConfigRulesGcp struct {
	Scopes             param.Field[[]string]                                                 `json:"scopes" api:"required"`
	ServiceAccountJson param.Field[SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJson] `json:"service_account_json" api:"required"`
}

func (r SandboxBoxNewParamsProxyConfigRulesGcp) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJson struct {
	Type  param.Field[SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                         `json:"is_set"`
	Value param.Field[string]                                                       `json:"value"`
}

func (r SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJson) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
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

type SandboxBoxUpdateParams struct {
	CPUMillicores          param.Field[int64]                             `json:"cpu_millicores"`
	DeleteAfterStopSeconds param.Field[int64]                             `json:"delete_after_stop_seconds"`
	FsCapacityBytes        param.Field[int64]                             `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64]                             `json:"idle_ttl_seconds"`
	MemBytes               param.Field[int64]                             `json:"mem_bytes"`
	Name                   param.Field[string]                            `json:"name"`
	ProxyConfig            param.Field[SandboxBoxUpdateParamsProxyConfig] `json:"proxy_config"`
	TagValueIDs            param.Field[[]string]                          `json:"tag_value_ids"`
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
	Name    param.Field[string]                                         `json:"name" api:"required"`
	Aws     param.Field[SandboxBoxUpdateParamsProxyConfigRulesAws]      `json:"aws"`
	Enabled param.Field[bool]                                           `json:"enabled"`
	Gcp     param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcp]      `json:"gcp"`
	Headers param.Field[[]SandboxBoxUpdateParamsProxyConfigRulesHeader] `json:"headers"`
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
	AccessKeyID     param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyID]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKey] `json:"secret_access_key" api:"required"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAws) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyID struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                     `json:"is_set"`
	Value param.Field[string]                                                   `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyID) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKey struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                         `json:"is_set"`
	Value param.Field[string]                                                       `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKey) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateParamsProxyConfigRulesGcp struct {
	Scopes             param.Field[[]string]                                                    `json:"scopes" api:"required"`
	ServiceAccountJson param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJson] `json:"service_account_json" api:"required"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcp) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJson struct {
	Type  param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                            `json:"is_set"`
	Value param.Field[string]                                                          `json:"value"`
}

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJson) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
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

type SandboxBoxListParams struct {
	// Filter by creator identity. Only 'me' is supported.
	CreatedBy param.Field[string] `query:"created_by"`
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
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SandboxBoxNewSnapshotParams struct {
	Name param.Field[string] `json:"name" api:"required"`
	// if omitted, creates a fresh checkpoint from the running VM
	Checkpoint param.Field[string] `json:"checkpoint"`
	// sandbox-local Docker image to export
	DockerImage param.Field[string] `json:"docker_image"`
	// required for Docker image export unless the sandbox has a capacity
	FsCapacityBytes param.Field[int64] `json:"fs_capacity_bytes"`
	// IncludeMemory, when true, captures a full VM memory snapshot alongside the
	// filesystem clone. Only honored when the sandbox is running AND Checkpoint is
	// omitted (i.e. a fresh in-VM checkpoint is requested). Defaults to false to keep
	// snapshots small unless memory restore is explicitly desired.
	IncludeMemory param.Field[bool] `json:"include_memory"`
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
