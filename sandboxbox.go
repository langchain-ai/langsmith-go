// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"reflect"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/tidwall/gjson"
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

// Create a new sandbox using server defaults. Provide at most one of
// `snapshot_id` or `snapshot_name` only when booting from a reusable snapshot.
func (r *SandboxBoxService) New(ctx context.Context, body SandboxBoxNewParams, opts ...option.RequestOption) (res *SandboxBoxNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/boxes"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Retrieve a sandbox by name. Stale provisioning sandboxes are auto-failed.
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

// Update a sandbox's display name. The name must be unique within the tenant.
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

// List sandboxes for the authenticated tenant, with optional filtering, sorting,
// and pagination.
func (r *SandboxBoxService) List(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) (res *SandboxBoxListResponse, err error) {
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

// Retrieve the lightweight status of a sandbox for polling.
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
	MountConfig            SandboxBoxNewResponseMountConfig `json:"mount_config"`
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
	MountConfig            apijson.Field
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

type SandboxBoxNewResponseMountConfig struct {
	Auth   SandboxBoxNewResponseMountConfigAuth    `json:"auth"`
	Mounts []SandboxBoxNewResponseMountConfigMount `json:"mounts"`
	JSON   sandboxBoxNewResponseMountConfigJSON    `json:"-"`
}

// sandboxBoxNewResponseMountConfigJSON contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfig]
type sandboxBoxNewResponseMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuth struct {
	Aws  SandboxBoxNewResponseMountConfigAuthAws  `json:"aws"`
	Gcp  SandboxBoxNewResponseMountConfigAuthGcp  `json:"gcp"`
	JSON sandboxBoxNewResponseMountConfigAuthJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigAuthJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseMountConfigAuth]
type sandboxBoxNewResponseMountConfigAuthJSON struct {
	Aws         apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuthAws struct {
	AccessKeyID     SandboxBoxNewResponseMountConfigAuthAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxNewResponseMountConfigAuthAwsJSON            `json:"-"`
}

// sandboxBoxNewResponseMountConfigAuthAwsJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseMountConfigAuthAws]
type sandboxBoxNewResponseMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuthAwsAccessKeyID struct {
	Type  SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                   `json:"is_set"`
	Value string                                                 `json:"value"`
	JSON  sandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxNewResponseMountConfigAuthAwsAccessKeyID]
type sandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigAuthAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxBoxNewResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKey struct {
	Type  SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                       `json:"is_set"`
	Value string                                                     `json:"value"`
	JSON  sandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct [SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKey]
type sandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxBoxNewResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigAuthGcp struct {
	ServiceAccountJson SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxNewResponseMountConfigAuthGcpJSON               `json:"-"`
}

// sandboxBoxNewResponseMountConfigAuthGcpJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseMountConfigAuthGcp]
type sandboxBoxNewResponseMountConfigAuthGcpJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJson struct {
	Type  SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                          `json:"is_set"`
	Value string                                                        `json:"value"`
	JSON  sandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJson]
type sandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxBoxNewResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMount struct {
	ID        string                                     `json:"id" api:"required"`
	MountPath string                                     `json:"mount_path" api:"required"`
	Type      SandboxBoxNewResponseMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                               `json:"s3"`
	JSON  sandboxBoxNewResponseMountConfigMountJSON `json:"-"`
	union SandboxBoxNewResponseMountConfigMountsUnion
}

// sandboxBoxNewResponseMountConfigMountJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseMountConfigMount]
type sandboxBoxNewResponseMountConfigMountJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sandboxBoxNewResponseMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxNewResponseMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxNewResponseMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxNewResponseMountConfigMountsUnion] interface which
// you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxNewResponseMountConfigMount) AsUnion() SandboxBoxNewResponseMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec].
type SandboxBoxNewResponseMountConfigMountsUnion interface {
	implementsSandboxBoxNewResponseMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxNewResponseMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                                 `json:"id" api:"required"`
	MountPath string                                                                 `json:"mount_path" api:"required"`
	S3        SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                   `json:"read_only"`
	JSON      sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	S3          apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxNewResponseMountConfigMount() {
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                                  `json:"bucket" api:"required"`
	EndpointURL string                                                                  `json:"endpoint_url" api:"required"`
	Region      string                                                                  `json:"region" api:"required"`
	PathStyle   bool                                                                    `json:"path_style"`
	Prefix      string                                                                  `json:"prefix"`
	JSON        sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                      `json:"max_size_bytes"`
	WritebackSeconds int64                                                                      `json:"writeback_seconds"`
	JSON             sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                                   `json:"bucket" api:"required"`
	Prefix string                                                                   `json:"prefix"`
	JSON   sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                                   `json:"remote_url" api:"required"`
	Ref                    SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                    `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                      `json:"name" api:"required"`
	Type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxNewResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                                  `json:"id" api:"required"`
	Gcs       SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                                  `json:"mount_path" api:"required"`
	Type      SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                    `json:"read_only"`
	S3        SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
	ID          apijson.Field
	Gcs         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxNewResponseMountConfigMount() {
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                                    `json:"bucket" api:"required"`
	Prefix string                                                                    `json:"prefix"`
	JSON   sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                       `json:"max_size_bytes"`
	WritebackSeconds int64                                                                       `json:"writeback_seconds"`
	JSON             sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                                    `json:"remote_url" api:"required"`
	Ref                    SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                     `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                       `json:"name" api:"required"`
	Type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                                   `json:"bucket" api:"required"`
	EndpointURL string                                                                   `json:"endpoint_url" api:"required"`
	Region      string                                                                   `json:"region" api:"required"`
	PathStyle   bool                                                                     `json:"path_style"`
	Prefix      string                                                                   `json:"prefix"`
	JSON        sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                                `json:"id" api:"required"`
	Git       SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                                `json:"mount_path" api:"required"`
	Type      SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                                  `json:"read_only"`
	S3        SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
	ID          apijson.Field
	Git         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxNewResponseMountConfigMount() {
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                                  `json:"remote_url" api:"required"`
	Ref                    SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                   `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                                     `json:"name" api:"required"`
	Type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                     `json:"max_size_bytes"`
	WritebackSeconds int64                                                                     `json:"writeback_seconds"`
	JSON             sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                                  `json:"bucket" api:"required"`
	Prefix string                                                                  `json:"prefix"`
	JSON   sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                                 `json:"bucket" api:"required"`
	EndpointURL string                                                                 `json:"endpoint_url" api:"required"`
	Region      string                                                                 `json:"region" api:"required"`
	PathStyle   bool                                                                   `json:"path_style"`
	Prefix      string                                                                 `json:"prefix"`
	JSON        sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountConfigMountsType string

const (
	SandboxBoxNewResponseMountConfigMountsTypeS3  SandboxBoxNewResponseMountConfigMountsType = "s3"
	SandboxBoxNewResponseMountConfigMountsTypeGcs SandboxBoxNewResponseMountConfigMountsType = "gcs"
	SandboxBoxNewResponseMountConfigMountsTypeGit SandboxBoxNewResponseMountConfigMountsType = "git"
)

func (r SandboxBoxNewResponseMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountConfigMountsTypeS3, SandboxBoxNewResponseMountConfigMountsTypeGcs, SandboxBoxNewResponseMountConfigMountsTypeGit:
		return true
	}
	return false
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
	Name    string                                        `json:"name" api:"required"`
	Aws     SandboxBoxNewResponseProxyConfigRulesAws      `json:"aws"`
	Enabled bool                                          `json:"enabled"`
	Gcp     SandboxBoxNewResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxBoxNewResponseProxyConfigRulesHeader `json:"headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts []string                                 `json:"match_hosts"`
	MatchPaths []string                                 `json:"match_paths"`
	Type       string                                   `json:"type"`
	JSON       sandboxBoxNewResponseProxyConfigRuleJSON `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseProxyConfigRule]
type sandboxBoxNewResponseProxyConfigRuleJSON struct {
	Name        apijson.Field
	Aws         apijson.Field
	Enabled     apijson.Field
	Gcp         apijson.Field
	Headers     apijson.Field
	MatchHosts  apijson.Field
	MatchPaths  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesAws struct {
	AccessKeyID     SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxNewResponseProxyConfigRulesAwsJSON            `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRulesAwsJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseProxyConfigRulesAws]
type sandboxBoxNewResponseProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyID struct {
	Type  SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                    `json:"is_set"`
	Value string                                                  `json:"value"`
	JSON  sandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyID]
type sandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxNewResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKey struct {
	Type  SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                        `json:"is_set"`
	Value string                                                      `json:"value"`
	JSON  sandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKey]
type sandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxNewResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxNewResponseProxyConfigRulesGcp struct {
	Scopes             []string                                                   `json:"scopes" api:"required"`
	ServiceAccountJson SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxNewResponseProxyConfigRulesGcpJSON               `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRulesGcpJSON contains the JSON metadata for the
// struct [SandboxBoxNewResponseProxyConfigRulesGcp]
type sandboxBoxNewResponseProxyConfigRulesGcpJSON struct {
	Scopes             apijson.Field
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJson struct {
	Type  SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                           `json:"is_set"`
	Value string                                                         `json:"value"`
	JSON  sandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJson]
type sandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxNewResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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
	MountConfig            SandboxBoxGetResponseMountConfig `json:"mount_config"`
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
	MountConfig            apijson.Field
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

type SandboxBoxGetResponseMountConfig struct {
	Auth   SandboxBoxGetResponseMountConfigAuth    `json:"auth"`
	Mounts []SandboxBoxGetResponseMountConfigMount `json:"mounts"`
	JSON   sandboxBoxGetResponseMountConfigJSON    `json:"-"`
}

// sandboxBoxGetResponseMountConfigJSON contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfig]
type sandboxBoxGetResponseMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuth struct {
	Aws  SandboxBoxGetResponseMountConfigAuthAws  `json:"aws"`
	Gcp  SandboxBoxGetResponseMountConfigAuthGcp  `json:"gcp"`
	JSON sandboxBoxGetResponseMountConfigAuthJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigAuthJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseMountConfigAuth]
type sandboxBoxGetResponseMountConfigAuthJSON struct {
	Aws         apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuthAws struct {
	AccessKeyID     SandboxBoxGetResponseMountConfigAuthAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxGetResponseMountConfigAuthAwsJSON            `json:"-"`
}

// sandboxBoxGetResponseMountConfigAuthAwsJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseMountConfigAuthAws]
type sandboxBoxGetResponseMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuthAwsAccessKeyID struct {
	Type  SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                   `json:"is_set"`
	Value string                                                 `json:"value"`
	JSON  sandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxGetResponseMountConfigAuthAwsAccessKeyID]
type sandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigAuthAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxBoxGetResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKey struct {
	Type  SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                       `json:"is_set"`
	Value string                                                     `json:"value"`
	JSON  sandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct [SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKey]
type sandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxBoxGetResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigAuthGcp struct {
	ServiceAccountJson SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxGetResponseMountConfigAuthGcpJSON               `json:"-"`
}

// sandboxBoxGetResponseMountConfigAuthGcpJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseMountConfigAuthGcp]
type sandboxBoxGetResponseMountConfigAuthGcpJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJson struct {
	Type  SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                          `json:"is_set"`
	Value string                                                        `json:"value"`
	JSON  sandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJson]
type sandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxBoxGetResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMount struct {
	ID        string                                     `json:"id" api:"required"`
	MountPath string                                     `json:"mount_path" api:"required"`
	Type      SandboxBoxGetResponseMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                               `json:"s3"`
	JSON  sandboxBoxGetResponseMountConfigMountJSON `json:"-"`
	union SandboxBoxGetResponseMountConfigMountsUnion
}

// sandboxBoxGetResponseMountConfigMountJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseMountConfigMount]
type sandboxBoxGetResponseMountConfigMountJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sandboxBoxGetResponseMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxGetResponseMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxGetResponseMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxGetResponseMountConfigMountsUnion] interface which
// you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxGetResponseMountConfigMount) AsUnion() SandboxBoxGetResponseMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec].
type SandboxBoxGetResponseMountConfigMountsUnion interface {
	implementsSandboxBoxGetResponseMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxGetResponseMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                                 `json:"id" api:"required"`
	MountPath string                                                                 `json:"mount_path" api:"required"`
	S3        SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                   `json:"read_only"`
	JSON      sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	S3          apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxGetResponseMountConfigMount() {
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                                  `json:"bucket" api:"required"`
	EndpointURL string                                                                  `json:"endpoint_url" api:"required"`
	Region      string                                                                  `json:"region" api:"required"`
	PathStyle   bool                                                                    `json:"path_style"`
	Prefix      string                                                                  `json:"prefix"`
	JSON        sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                      `json:"max_size_bytes"`
	WritebackSeconds int64                                                                      `json:"writeback_seconds"`
	JSON             sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                                   `json:"bucket" api:"required"`
	Prefix string                                                                   `json:"prefix"`
	JSON   sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                                   `json:"remote_url" api:"required"`
	Ref                    SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                    `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                      `json:"name" api:"required"`
	Type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxGetResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                                  `json:"id" api:"required"`
	Gcs       SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                                  `json:"mount_path" api:"required"`
	Type      SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                    `json:"read_only"`
	S3        SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
	ID          apijson.Field
	Gcs         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxGetResponseMountConfigMount() {
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                                    `json:"bucket" api:"required"`
	Prefix string                                                                    `json:"prefix"`
	JSON   sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                       `json:"max_size_bytes"`
	WritebackSeconds int64                                                                       `json:"writeback_seconds"`
	JSON             sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                                    `json:"remote_url" api:"required"`
	Ref                    SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                     `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                       `json:"name" api:"required"`
	Type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                                   `json:"bucket" api:"required"`
	EndpointURL string                                                                   `json:"endpoint_url" api:"required"`
	Region      string                                                                   `json:"region" api:"required"`
	PathStyle   bool                                                                     `json:"path_style"`
	Prefix      string                                                                   `json:"prefix"`
	JSON        sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                                `json:"id" api:"required"`
	Git       SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                                `json:"mount_path" api:"required"`
	Type      SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                                  `json:"read_only"`
	S3        SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
	ID          apijson.Field
	Git         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxGetResponseMountConfigMount() {
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                                  `json:"remote_url" api:"required"`
	Ref                    SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                   `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                                     `json:"name" api:"required"`
	Type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                     `json:"max_size_bytes"`
	WritebackSeconds int64                                                                     `json:"writeback_seconds"`
	JSON             sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                                  `json:"bucket" api:"required"`
	Prefix string                                                                  `json:"prefix"`
	JSON   sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                                 `json:"bucket" api:"required"`
	EndpointURL string                                                                 `json:"endpoint_url" api:"required"`
	Region      string                                                                 `json:"region" api:"required"`
	PathStyle   bool                                                                   `json:"path_style"`
	Prefix      string                                                                 `json:"prefix"`
	JSON        sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountConfigMountsType string

const (
	SandboxBoxGetResponseMountConfigMountsTypeS3  SandboxBoxGetResponseMountConfigMountsType = "s3"
	SandboxBoxGetResponseMountConfigMountsTypeGcs SandboxBoxGetResponseMountConfigMountsType = "gcs"
	SandboxBoxGetResponseMountConfigMountsTypeGit SandboxBoxGetResponseMountConfigMountsType = "git"
)

func (r SandboxBoxGetResponseMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountConfigMountsTypeS3, SandboxBoxGetResponseMountConfigMountsTypeGcs, SandboxBoxGetResponseMountConfigMountsTypeGit:
		return true
	}
	return false
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
	Name    string                                        `json:"name" api:"required"`
	Aws     SandboxBoxGetResponseProxyConfigRulesAws      `json:"aws"`
	Enabled bool                                          `json:"enabled"`
	Gcp     SandboxBoxGetResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxBoxGetResponseProxyConfigRulesHeader `json:"headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts []string                                 `json:"match_hosts"`
	MatchPaths []string                                 `json:"match_paths"`
	Type       string                                   `json:"type"`
	JSON       sandboxBoxGetResponseProxyConfigRuleJSON `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseProxyConfigRule]
type sandboxBoxGetResponseProxyConfigRuleJSON struct {
	Name        apijson.Field
	Aws         apijson.Field
	Enabled     apijson.Field
	Gcp         apijson.Field
	Headers     apijson.Field
	MatchHosts  apijson.Field
	MatchPaths  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesAws struct {
	AccessKeyID     SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxGetResponseProxyConfigRulesAwsJSON            `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRulesAwsJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseProxyConfigRulesAws]
type sandboxBoxGetResponseProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyID struct {
	Type  SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                    `json:"is_set"`
	Value string                                                  `json:"value"`
	JSON  sandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyID]
type sandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxGetResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKey struct {
	Type  SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                        `json:"is_set"`
	Value string                                                      `json:"value"`
	JSON  sandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKey]
type sandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxGetResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxGetResponseProxyConfigRulesGcp struct {
	Scopes             []string                                                   `json:"scopes" api:"required"`
	ServiceAccountJson SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxGetResponseProxyConfigRulesGcpJSON               `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRulesGcpJSON contains the JSON metadata for the
// struct [SandboxBoxGetResponseProxyConfigRulesGcp]
type sandboxBoxGetResponseProxyConfigRulesGcpJSON struct {
	Scopes             apijson.Field
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJson struct {
	Type  SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                           `json:"is_set"`
	Value string                                                         `json:"value"`
	JSON  sandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJson]
type sandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxGetResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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
	MountConfig            SandboxBoxUpdateResponseMountConfig `json:"mount_config"`
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
	MountConfig            apijson.Field
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

type SandboxBoxUpdateResponseMountConfig struct {
	Auth   SandboxBoxUpdateResponseMountConfigAuth    `json:"auth"`
	Mounts []SandboxBoxUpdateResponseMountConfigMount `json:"mounts"`
	JSON   sandboxBoxUpdateResponseMountConfigJSON    `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigJSON contains the JSON metadata for the
// struct [SandboxBoxUpdateResponseMountConfig]
type sandboxBoxUpdateResponseMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuth struct {
	Aws  SandboxBoxUpdateResponseMountConfigAuthAws  `json:"aws"`
	Gcp  SandboxBoxUpdateResponseMountConfigAuthGcp  `json:"gcp"`
	JSON sandboxBoxUpdateResponseMountConfigAuthJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigAuthJSON contains the JSON metadata for the
// struct [SandboxBoxUpdateResponseMountConfigAuth]
type sandboxBoxUpdateResponseMountConfigAuthJSON struct {
	Aws         apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuthAws struct {
	AccessKeyID     SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxUpdateResponseMountConfigAuthAwsJSON            `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigAuthAwsJSON contains the JSON metadata for
// the struct [SandboxBoxUpdateResponseMountConfigAuthAws]
type sandboxBoxUpdateResponseMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyID struct {
	Type  SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                      `json:"is_set"`
	Value string                                                    `json:"value"`
	JSON  sandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyID]
type sandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxBoxUpdateResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKey struct {
	Type  SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                          `json:"is_set"`
	Value string                                                        `json:"value"`
	JSON  sandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKey]
type sandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxBoxUpdateResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigAuthGcp struct {
	ServiceAccountJson SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxUpdateResponseMountConfigAuthGcpJSON               `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigAuthGcpJSON contains the JSON metadata for
// the struct [SandboxBoxUpdateResponseMountConfigAuthGcp]
type sandboxBoxUpdateResponseMountConfigAuthGcpJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJson struct {
	Type  SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                             `json:"is_set"`
	Value string                                                           `json:"value"`
	JSON  sandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJson]
type sandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxBoxUpdateResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMount struct {
	ID        string                                        `json:"id" api:"required"`
	MountPath string                                        `json:"mount_path" api:"required"`
	Type      SandboxBoxUpdateResponseMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                                  `json:"s3"`
	JSON  sandboxBoxUpdateResponseMountConfigMountJSON `json:"-"`
	union SandboxBoxUpdateResponseMountConfigMountsUnion
}

// sandboxBoxUpdateResponseMountConfigMountJSON contains the JSON metadata for the
// struct [SandboxBoxUpdateResponseMountConfigMount]
type sandboxBoxUpdateResponseMountConfigMountJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sandboxBoxUpdateResponseMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxUpdateResponseMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxUpdateResponseMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxUpdateResponseMountConfigMountsUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxUpdateResponseMountConfigMount) AsUnion() SandboxBoxUpdateResponseMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec].
type SandboxBoxUpdateResponseMountConfigMountsUnion interface {
	implementsSandboxBoxUpdateResponseMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxUpdateResponseMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                                    `json:"id" api:"required"`
	MountPath string                                                                    `json:"mount_path" api:"required"`
	S3        SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                      `json:"read_only"`
	JSON      sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	S3          apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxUpdateResponseMountConfigMount() {
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                                     `json:"bucket" api:"required"`
	EndpointURL string                                                                     `json:"endpoint_url" api:"required"`
	Region      string                                                                     `json:"region" api:"required"`
	PathStyle   bool                                                                       `json:"path_style"`
	Prefix      string                                                                     `json:"prefix"`
	JSON        sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                         `json:"max_size_bytes"`
	WritebackSeconds int64                                                                         `json:"writeback_seconds"`
	JSON             sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                                      `json:"bucket" api:"required"`
	Prefix string                                                                      `json:"prefix"`
	JSON   sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                                      `json:"remote_url" api:"required"`
	Ref                    SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                       `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                         `json:"name" api:"required"`
	Type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxUpdateResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                                     `json:"id" api:"required"`
	Gcs       SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                                     `json:"mount_path" api:"required"`
	Type      SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                       `json:"read_only"`
	S3        SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
	ID          apijson.Field
	Gcs         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxUpdateResponseMountConfigMount() {
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                                       `json:"bucket" api:"required"`
	Prefix string                                                                       `json:"prefix"`
	JSON   sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                          `json:"max_size_bytes"`
	WritebackSeconds int64                                                                          `json:"writeback_seconds"`
	JSON             sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                                       `json:"remote_url" api:"required"`
	Ref                    SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                        `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                          `json:"name" api:"required"`
	Type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                                      `json:"bucket" api:"required"`
	EndpointURL string                                                                      `json:"endpoint_url" api:"required"`
	Region      string                                                                      `json:"region" api:"required"`
	PathStyle   bool                                                                        `json:"path_style"`
	Prefix      string                                                                      `json:"prefix"`
	JSON        sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                                   `json:"id" api:"required"`
	Git       SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                                   `json:"mount_path" api:"required"`
	Type      SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                                     `json:"read_only"`
	S3        SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
	ID          apijson.Field
	Git         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxUpdateResponseMountConfigMount() {
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                                     `json:"remote_url" api:"required"`
	Ref                    SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                      `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                                        `json:"name" api:"required"`
	Type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                        `json:"max_size_bytes"`
	WritebackSeconds int64                                                                        `json:"writeback_seconds"`
	JSON             sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                                     `json:"bucket" api:"required"`
	Prefix string                                                                     `json:"prefix"`
	JSON   sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                                    `json:"bucket" api:"required"`
	EndpointURL string                                                                    `json:"endpoint_url" api:"required"`
	Region      string                                                                    `json:"region" api:"required"`
	PathStyle   bool                                                                      `json:"path_style"`
	Prefix      string                                                                    `json:"prefix"`
	JSON        sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountConfigMountsType string

const (
	SandboxBoxUpdateResponseMountConfigMountsTypeS3  SandboxBoxUpdateResponseMountConfigMountsType = "s3"
	SandboxBoxUpdateResponseMountConfigMountsTypeGcs SandboxBoxUpdateResponseMountConfigMountsType = "gcs"
	SandboxBoxUpdateResponseMountConfigMountsTypeGit SandboxBoxUpdateResponseMountConfigMountsType = "git"
)

func (r SandboxBoxUpdateResponseMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountConfigMountsTypeS3, SandboxBoxUpdateResponseMountConfigMountsTypeGcs, SandboxBoxUpdateResponseMountConfigMountsTypeGit:
		return true
	}
	return false
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
	Name    string                                           `json:"name" api:"required"`
	Aws     SandboxBoxUpdateResponseProxyConfigRulesAws      `json:"aws"`
	Enabled bool                                             `json:"enabled"`
	Gcp     SandboxBoxUpdateResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxBoxUpdateResponseProxyConfigRulesHeader `json:"headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts []string                                    `json:"match_hosts"`
	MatchPaths []string                                    `json:"match_paths"`
	Type       string                                      `json:"type"`
	JSON       sandboxBoxUpdateResponseProxyConfigRuleJSON `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxUpdateResponseProxyConfigRule]
type sandboxBoxUpdateResponseProxyConfigRuleJSON struct {
	Name        apijson.Field
	Aws         apijson.Field
	Enabled     apijson.Field
	Gcp         apijson.Field
	Headers     apijson.Field
	MatchHosts  apijson.Field
	MatchPaths  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesAws struct {
	AccessKeyID     SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxUpdateResponseProxyConfigRulesAwsJSON            `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRulesAwsJSON contains the JSON metadata for
// the struct [SandboxBoxUpdateResponseProxyConfigRulesAws]
type sandboxBoxUpdateResponseProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyID struct {
	Type  SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                       `json:"is_set"`
	Value string                                                     `json:"value"`
	JSON  sandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyID]
type sandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxUpdateResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKey struct {
	Type  SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                           `json:"is_set"`
	Value string                                                         `json:"value"`
	JSON  sandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKey]
type sandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxUpdateResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseProxyConfigRulesGcp struct {
	Scopes             []string                                                      `json:"scopes" api:"required"`
	ServiceAccountJson SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxUpdateResponseProxyConfigRulesGcpJSON               `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRulesGcpJSON contains the JSON metadata for
// the struct [SandboxBoxUpdateResponseProxyConfigRulesGcp]
type sandboxBoxUpdateResponseProxyConfigRulesGcpJSON struct {
	Scopes             apijson.Field
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJson struct {
	Type  SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                              `json:"is_set"`
	Value string                                                            `json:"value"`
	JSON  sandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJson]
type sandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxUpdateResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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
	MountConfig            SandboxBoxListResponseSandboxesMountConfig `json:"mount_config"`
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
	MountConfig            apijson.Field
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

type SandboxBoxListResponseSandboxesMountConfig struct {
	Auth   SandboxBoxListResponseSandboxesMountConfigAuth    `json:"auth"`
	Mounts []SandboxBoxListResponseSandboxesMountConfigMount `json:"mounts"`
	JSON   sandboxBoxListResponseSandboxesMountConfigJSON    `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigJSON contains the JSON metadata for
// the struct [SandboxBoxListResponseSandboxesMountConfig]
type sandboxBoxListResponseSandboxesMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuth struct {
	Aws  SandboxBoxListResponseSandboxesMountConfigAuthAws  `json:"aws"`
	Gcp  SandboxBoxListResponseSandboxesMountConfigAuthGcp  `json:"gcp"`
	JSON sandboxBoxListResponseSandboxesMountConfigAuthJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigAuthJSON contains the JSON metadata
// for the struct [SandboxBoxListResponseSandboxesMountConfigAuth]
type sandboxBoxListResponseSandboxesMountConfigAuthJSON struct {
	Aws         apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuthAws struct {
	AccessKeyID     SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxListResponseSandboxesMountConfigAuthAwsJSON            `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigAuthAwsJSON contains the JSON metadata
// for the struct [SandboxBoxListResponseSandboxesMountConfigAuthAws]
type sandboxBoxListResponseSandboxesMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyID struct {
	Type  SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                             `json:"is_set"`
	Value string                                                           `json:"value"`
	JSON  sandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDJSON contains the
// JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyID]
type sandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxBoxListResponseSandboxesMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKey struct {
	Type  SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                                 `json:"is_set"`
	Value string                                                               `json:"value"`
	JSON  sandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKey]
type sandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxBoxListResponseSandboxesMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigAuthGcp struct {
	ServiceAccountJson SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxListResponseSandboxesMountConfigAuthGcpJSON               `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigAuthGcpJSON contains the JSON metadata
// for the struct [SandboxBoxListResponseSandboxesMountConfigAuthGcp]
type sandboxBoxListResponseSandboxesMountConfigAuthGcpJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJson struct {
	Type  SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                                    `json:"is_set"`
	Value string                                                                  `json:"value"`
	JSON  sandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJson]
type sandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxBoxListResponseSandboxesMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMount struct {
	ID        string                                               `json:"id" api:"required"`
	MountPath string                                               `json:"mount_path" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                                         `json:"s3"`
	JSON  sandboxBoxListResponseSandboxesMountConfigMountJSON `json:"-"`
	union SandboxBoxListResponseSandboxesMountConfigMountsUnion
}

// sandboxBoxListResponseSandboxesMountConfigMountJSON contains the JSON metadata
// for the struct [SandboxBoxListResponseSandboxesMountConfigMount]
type sandboxBoxListResponseSandboxesMountConfigMountJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sandboxBoxListResponseSandboxesMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxListResponseSandboxesMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxListResponseSandboxesMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxListResponseSandboxesMountConfigMountsUnion]
// interface which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxListResponseSandboxesMountConfigMount) AsUnion() SandboxBoxListResponseSandboxesMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec]
// or [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec].
type SandboxBoxListResponseSandboxesMountConfigMountsUnion interface {
	implementsSandboxBoxListResponseSandboxesMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxListResponseSandboxesMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                                           `json:"id" api:"required"`
	MountPath string                                                                           `json:"mount_path" api:"required"`
	S3        SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                             `json:"read_only"`
	JSON      sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	S3          apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxListResponseSandboxesMountConfigMount() {
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                                            `json:"bucket" api:"required"`
	EndpointURL string                                                                            `json:"endpoint_url" api:"required"`
	Region      string                                                                            `json:"region" api:"required"`
	PathStyle   bool                                                                              `json:"path_style"`
	Prefix      string                                                                            `json:"prefix"`
	JSON        sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                                `json:"max_size_bytes"`
	WritebackSeconds int64                                                                                `json:"writeback_seconds"`
	JSON             sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                                             `json:"bucket" api:"required"`
	Prefix string                                                                             `json:"prefix"`
	JSON   sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                                             `json:"remote_url" api:"required"`
	Ref                    SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                              `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                                `json:"name" api:"required"`
	Type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                                            `json:"id" api:"required"`
	Gcs       SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                                            `json:"mount_path" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                              `json:"read_only"`
	S3        SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
	ID          apijson.Field
	Gcs         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxListResponseSandboxesMountConfigMount() {
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                                              `json:"bucket" api:"required"`
	Prefix string                                                                              `json:"prefix"`
	JSON   sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                                 `json:"max_size_bytes"`
	WritebackSeconds int64                                                                                 `json:"writeback_seconds"`
	JSON             sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                                              `json:"remote_url" api:"required"`
	Ref                    SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                               `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                                 `json:"name" api:"required"`
	Type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                                             `json:"bucket" api:"required"`
	EndpointURL string                                                                             `json:"endpoint_url" api:"required"`
	Region      string                                                                             `json:"region" api:"required"`
	PathStyle   bool                                                                               `json:"path_style"`
	Prefix      string                                                                             `json:"prefix"`
	JSON        sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                                          `json:"id" api:"required"`
	Git       SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                                          `json:"mount_path" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                                            `json:"read_only"`
	S3        SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
	ID          apijson.Field
	Git         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxListResponseSandboxesMountConfigMount() {
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                                            `json:"remote_url" api:"required"`
	Ref                    SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                             `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                                               `json:"name" api:"required"`
	Type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                               `json:"max_size_bytes"`
	WritebackSeconds int64                                                                               `json:"writeback_seconds"`
	JSON             sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                                            `json:"bucket" api:"required"`
	Prefix string                                                                            `json:"prefix"`
	JSON   sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                                           `json:"bucket" api:"required"`
	EndpointURL string                                                                           `json:"endpoint_url" api:"required"`
	Region      string                                                                           `json:"region" api:"required"`
	PathStyle   bool                                                                             `json:"path_style"`
	Prefix      string                                                                           `json:"prefix"`
	JSON        sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountConfigMountsType string

const (
	SandboxBoxListResponseSandboxesMountConfigMountsTypeS3  SandboxBoxListResponseSandboxesMountConfigMountsType = "s3"
	SandboxBoxListResponseSandboxesMountConfigMountsTypeGcs SandboxBoxListResponseSandboxesMountConfigMountsType = "gcs"
	SandboxBoxListResponseSandboxesMountConfigMountsTypeGit SandboxBoxListResponseSandboxesMountConfigMountsType = "git"
)

func (r SandboxBoxListResponseSandboxesMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountConfigMountsTypeS3, SandboxBoxListResponseSandboxesMountConfigMountsTypeGcs, SandboxBoxListResponseSandboxesMountConfigMountsTypeGit:
		return true
	}
	return false
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
	Name    string                                                  `json:"name" api:"required"`
	Aws     SandboxBoxListResponseSandboxesProxyConfigRulesAws      `json:"aws"`
	Enabled bool                                                    `json:"enabled"`
	Gcp     SandboxBoxListResponseSandboxesProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxBoxListResponseSandboxesProxyConfigRulesHeader `json:"headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts []string                                           `json:"match_hosts"`
	MatchPaths []string                                           `json:"match_paths"`
	Type       string                                             `json:"type"`
	JSON       sandboxBoxListResponseSandboxesProxyConfigRuleJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRuleJSON contains the JSON metadata
// for the struct [SandboxBoxListResponseSandboxesProxyConfigRule]
type sandboxBoxListResponseSandboxesProxyConfigRuleJSON struct {
	Name        apijson.Field
	Aws         apijson.Field
	Enabled     apijson.Field
	Gcp         apijson.Field
	Headers     apijson.Field
	MatchHosts  apijson.Field
	MatchPaths  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesAws struct {
	AccessKeyID     SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxListResponseSandboxesProxyConfigRulesAwsJSON            `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRulesAwsJSON contains the JSON
// metadata for the struct [SandboxBoxListResponseSandboxesProxyConfigRulesAws]
type sandboxBoxListResponseSandboxesProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyID struct {
	Type  SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                              `json:"is_set"`
	Value string                                                            `json:"value"`
	JSON  sandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDJSON contains the
// JSON metadata for the struct
// [SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyID]
type sandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxListResponseSandboxesProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKey struct {
	Type  SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                                  `json:"is_set"`
	Value string                                                                `json:"value"`
	JSON  sandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKey]
type sandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxListResponseSandboxesProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesProxyConfigRulesGcp struct {
	Scopes             []string                                                             `json:"scopes" api:"required"`
	ServiceAccountJson SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxListResponseSandboxesProxyConfigRulesGcpJSON               `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRulesGcpJSON contains the JSON
// metadata for the struct [SandboxBoxListResponseSandboxesProxyConfigRulesGcp]
type sandboxBoxListResponseSandboxesProxyConfigRulesGcpJSON struct {
	Scopes             apijson.Field
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJson struct {
	Type  SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                                     `json:"is_set"`
	Value string                                                                   `json:"value"`
	JSON  sandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJson]
type sandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxListResponseSandboxesProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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
	ID              string `json:"id"`
	CreatedAt       string `json:"created_at"`
	CreatedBy       string `json:"created_by"`
	DockerImage     string `json:"docker_image"`
	FsCapacityBytes int64  `json:"fs_capacity_bytes"`
	FsUsedBytes     int64  `json:"fs_used_bytes"`
	ImageDigest     string `json:"image_digest"`
	// MemorySnapshotSizeBytes is non-nil iff the snapshot was captured with VM memory
	// state. A non-nil value is the canonical signal that this snapshot can
	// warm-restore from memory; nil means rootfs only.
	MemorySnapshotSizeBytes int64                             `json:"memory_snapshot_size_bytes"`
	Name                    string                            `json:"name"`
	RegistryID              string                            `json:"registry_id"`
	SourceSandboxID         string                            `json:"source_sandbox_id"`
	Status                  string                            `json:"status"`
	StatusMessage           string                            `json:"status_message"`
	UpdatedAt               string                            `json:"updated_at"`
	JSON                    sandboxBoxNewSnapshotResponseJSON `json:"-"`
}

// sandboxBoxNewSnapshotResponseJSON contains the JSON metadata for the struct
// [SandboxBoxNewSnapshotResponse]
type sandboxBoxNewSnapshotResponseJSON struct {
	ID                      apijson.Field
	CreatedAt               apijson.Field
	CreatedBy               apijson.Field
	DockerImage             apijson.Field
	FsCapacityBytes         apijson.Field
	FsUsedBytes             apijson.Field
	ImageDigest             apijson.Field
	MemorySnapshotSizeBytes apijson.Field
	Name                    apijson.Field
	RegistryID              apijson.Field
	SourceSandboxID         apijson.Field
	Status                  apijson.Field
	StatusMessage           apijson.Field
	UpdatedAt               apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
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
	MountConfig            SandboxBoxStartResponseMountConfig `json:"mount_config"`
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
	MountConfig            apijson.Field
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

type SandboxBoxStartResponseMountConfig struct {
	Auth   SandboxBoxStartResponseMountConfigAuth    `json:"auth"`
	Mounts []SandboxBoxStartResponseMountConfigMount `json:"mounts"`
	JSON   sandboxBoxStartResponseMountConfigJSON    `json:"-"`
}

// sandboxBoxStartResponseMountConfigJSON contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfig]
type sandboxBoxStartResponseMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuth struct {
	Aws  SandboxBoxStartResponseMountConfigAuthAws  `json:"aws"`
	Gcp  SandboxBoxStartResponseMountConfigAuthGcp  `json:"gcp"`
	JSON sandboxBoxStartResponseMountConfigAuthJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigAuthJSON contains the JSON metadata for the
// struct [SandboxBoxStartResponseMountConfigAuth]
type sandboxBoxStartResponseMountConfigAuthJSON struct {
	Aws         apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuthAws struct {
	AccessKeyID     SandboxBoxStartResponseMountConfigAuthAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxStartResponseMountConfigAuthAwsJSON            `json:"-"`
}

// sandboxBoxStartResponseMountConfigAuthAwsJSON contains the JSON metadata for the
// struct [SandboxBoxStartResponseMountConfigAuthAws]
type sandboxBoxStartResponseMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuthAwsAccessKeyID struct {
	Type  SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                     `json:"is_set"`
	Value string                                                   `json:"value"`
	JSON  sandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxStartResponseMountConfigAuthAwsAccessKeyID]
type sandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigAuthAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxBoxStartResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKey struct {
	Type  SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                         `json:"is_set"`
	Value string                                                       `json:"value"`
	JSON  sandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKey]
type sandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxBoxStartResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigAuthGcp struct {
	ServiceAccountJson SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxStartResponseMountConfigAuthGcpJSON               `json:"-"`
}

// sandboxBoxStartResponseMountConfigAuthGcpJSON contains the JSON metadata for the
// struct [SandboxBoxStartResponseMountConfigAuthGcp]
type sandboxBoxStartResponseMountConfigAuthGcpJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJson struct {
	Type  SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                            `json:"is_set"`
	Value string                                                          `json:"value"`
	JSON  sandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJson]
type sandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxBoxStartResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMount struct {
	ID        string                                       `json:"id" api:"required"`
	MountPath string                                       `json:"mount_path" api:"required"`
	Type      SandboxBoxStartResponseMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                                 `json:"s3"`
	JSON  sandboxBoxStartResponseMountConfigMountJSON `json:"-"`
	union SandboxBoxStartResponseMountConfigMountsUnion
}

// sandboxBoxStartResponseMountConfigMountJSON contains the JSON metadata for the
// struct [SandboxBoxStartResponseMountConfigMount]
type sandboxBoxStartResponseMountConfigMountJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sandboxBoxStartResponseMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxStartResponseMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxStartResponseMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxStartResponseMountConfigMountsUnion] interface
// which you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxStartResponseMountConfigMount) AsUnion() SandboxBoxStartResponseMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec].
type SandboxBoxStartResponseMountConfigMountsUnion interface {
	implementsSandboxBoxStartResponseMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxStartResponseMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                                   `json:"id" api:"required"`
	MountPath string                                                                   `json:"mount_path" api:"required"`
	S3        SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                     `json:"read_only"`
	JSON      sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	S3          apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxStartResponseMountConfigMount() {
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                                    `json:"bucket" api:"required"`
	EndpointURL string                                                                    `json:"endpoint_url" api:"required"`
	Region      string                                                                    `json:"region" api:"required"`
	PathStyle   bool                                                                      `json:"path_style"`
	Prefix      string                                                                    `json:"prefix"`
	JSON        sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                        `json:"max_size_bytes"`
	WritebackSeconds int64                                                                        `json:"writeback_seconds"`
	JSON             sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                                     `json:"bucket" api:"required"`
	Prefix string                                                                     `json:"prefix"`
	JSON   sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                                     `json:"remote_url" api:"required"`
	Ref                    SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                      `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                        `json:"name" api:"required"`
	Type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxStartResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                                    `json:"id" api:"required"`
	Gcs       SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                                    `json:"mount_path" api:"required"`
	Type      SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                      `json:"read_only"`
	S3        SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
	ID          apijson.Field
	Gcs         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxStartResponseMountConfigMount() {
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                                      `json:"bucket" api:"required"`
	Prefix string                                                                      `json:"prefix"`
	JSON   sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                         `json:"max_size_bytes"`
	WritebackSeconds int64                                                                         `json:"writeback_seconds"`
	JSON             sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                                      `json:"remote_url" api:"required"`
	Ref                    SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                       `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                         `json:"name" api:"required"`
	Type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                                     `json:"bucket" api:"required"`
	EndpointURL string                                                                     `json:"endpoint_url" api:"required"`
	Region      string                                                                     `json:"region" api:"required"`
	PathStyle   bool                                                                       `json:"path_style"`
	Prefix      string                                                                     `json:"prefix"`
	JSON        sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                                  `json:"id" api:"required"`
	Git       SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                                  `json:"mount_path" api:"required"`
	Type      SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                                    `json:"read_only"`
	S3        SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
	ID          apijson.Field
	Git         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Gcs         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxStartResponseMountConfigMount() {
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                                    `json:"remote_url" api:"required"`
	Ref                    SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                     `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                                       `json:"name" api:"required"`
	Type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                       `json:"max_size_bytes"`
	WritebackSeconds int64                                                                       `json:"writeback_seconds"`
	JSON             sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                                    `json:"bucket" api:"required"`
	Prefix string                                                                    `json:"prefix"`
	JSON   sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                                   `json:"bucket" api:"required"`
	EndpointURL string                                                                   `json:"endpoint_url" api:"required"`
	Region      string                                                                   `json:"region" api:"required"`
	PathStyle   bool                                                                     `json:"path_style"`
	Prefix      string                                                                   `json:"prefix"`
	JSON        sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountConfigMountsType string

const (
	SandboxBoxStartResponseMountConfigMountsTypeS3  SandboxBoxStartResponseMountConfigMountsType = "s3"
	SandboxBoxStartResponseMountConfigMountsTypeGcs SandboxBoxStartResponseMountConfigMountsType = "gcs"
	SandboxBoxStartResponseMountConfigMountsTypeGit SandboxBoxStartResponseMountConfigMountsType = "git"
)

func (r SandboxBoxStartResponseMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountConfigMountsTypeS3, SandboxBoxStartResponseMountConfigMountsTypeGcs, SandboxBoxStartResponseMountConfigMountsTypeGit:
		return true
	}
	return false
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
	Name    string                                          `json:"name" api:"required"`
	Aws     SandboxBoxStartResponseProxyConfigRulesAws      `json:"aws"`
	Enabled bool                                            `json:"enabled"`
	Gcp     SandboxBoxStartResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxBoxStartResponseProxyConfigRulesHeader `json:"headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts []string                                   `json:"match_hosts"`
	MatchPaths []string                                   `json:"match_paths"`
	Type       string                                     `json:"type"`
	JSON       sandboxBoxStartResponseProxyConfigRuleJSON `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRuleJSON contains the JSON metadata for the
// struct [SandboxBoxStartResponseProxyConfigRule]
type sandboxBoxStartResponseProxyConfigRuleJSON struct {
	Name        apijson.Field
	Aws         apijson.Field
	Enabled     apijson.Field
	Gcp         apijson.Field
	Headers     apijson.Field
	MatchHosts  apijson.Field
	MatchPaths  apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesAws struct {
	AccessKeyID     SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxBoxStartResponseProxyConfigRulesAwsJSON            `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRulesAwsJSON contains the JSON metadata for
// the struct [SandboxBoxStartResponseProxyConfigRulesAws]
type sandboxBoxStartResponseProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyID struct {
	Type  SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                      `json:"is_set"`
	Value string                                                    `json:"value"`
	JSON  sandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDJSON `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDJSON contains the JSON
// metadata for the struct [SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyID]
type sandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxBoxStartResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKey struct {
	Type  SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                          `json:"is_set"`
	Value string                                                        `json:"value"`
	JSON  sandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKey]
type sandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxBoxStartResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxBoxStartResponseProxyConfigRulesGcp struct {
	Scopes             []string                                                     `json:"scopes" api:"required"`
	ServiceAccountJson SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxBoxStartResponseProxyConfigRulesGcpJSON               `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRulesGcpJSON contains the JSON metadata for
// the struct [SandboxBoxStartResponseProxyConfigRulesGcp]
type sandboxBoxStartResponseProxyConfigRulesGcpJSON struct {
	Scopes             apijson.Field
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJson struct {
	Type  SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                             `json:"is_set"`
	Value string                                                           `json:"value"`
	JSON  sandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJson]
type sandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxBoxStartResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
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
	ID        param.Field[string]                                   `json:"id" api:"required"`
	MountPath param.Field[string]                                   `json:"mount_path" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountConfigMountsType] `json:"type" api:"required"`
	Cache     param.Field[interface{}]                              `json:"cache"`
	Gcs       param.Field[interface{}]                              `json:"gcs"`
	Git       param.Field[interface{}]                              `json:"git"`
	ReadOnly  param.Field[bool]                                     `json:"read_only"`
	S3        param.Field[interface{}]                              `json:"s3"`
}

func (r SandboxBoxNewParamsMountConfigMount) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMount) implementsSandboxBoxNewParamsMountConfigMountUnion() {}

// Satisfied by [SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpec],
// [SandboxBoxNewParamsMountConfigMount].
type SandboxBoxNewParamsMountConfigMountUnion interface {
	implementsSandboxBoxNewParamsMountConfigMountUnion()
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID        param.Field[string]                                                               `json:"id" api:"required"`
	MountPath param.Field[string]                                                               `json:"mount_path" api:"required"`
	S3        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3]    `json:"s3" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType]  `json:"type" api:"required"`
	Cache     param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecCache] `json:"cache"`
	Gcs       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGcs]   `json:"gcs"`
	Git       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecGit]   `json:"git"`
	ReadOnly  param.Field[bool]                                                                 `json:"read_only"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxNewParamsMountConfigMountUnion() {
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiS3BucketMountSpecTypeGit:
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
	ID        param.Field[string]                                                                `json:"id" api:"required"`
	Gcs       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGcs]   `json:"gcs" api:"required"`
	MountPath param.Field[string]                                                                `json:"mount_path" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType]  `json:"type" api:"required"`
	Cache     param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecCache] `json:"cache"`
	Git       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecGit]   `json:"git"`
	ReadOnly  param.Field[bool]                                                                  `json:"read_only"`
	S3        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecS3]    `json:"s3"`
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
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit:
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
	EndpointURL param.Field[string] `json:"endpoint_url" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGcsBucketMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID        param.Field[string]                                                              `json:"id" api:"required"`
	Git       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGit]   `json:"git" api:"required"`
	MountPath param.Field[string]                                                              `json:"mount_path" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType]  `json:"type" api:"required"`
	Cache     param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecCache] `json:"cache"`
	Gcs       param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGcs]   `json:"gcs"`
	ReadOnly  param.Field[bool]                                                                `json:"read_only"`
	S3        param.Field[SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecS3]    `json:"s3"`
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
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecTypeGit:
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

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountConfigMountsSandboxapiGitRepoMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountConfigMountsType string

const (
	SandboxBoxNewParamsMountConfigMountsTypeS3  SandboxBoxNewParamsMountConfigMountsType = "s3"
	SandboxBoxNewParamsMountConfigMountsTypeGcs SandboxBoxNewParamsMountConfigMountsType = "gcs"
	SandboxBoxNewParamsMountConfigMountsTypeGit SandboxBoxNewParamsMountConfigMountsType = "git"
)

func (r SandboxBoxNewParamsMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountConfigMountsTypeS3, SandboxBoxNewParamsMountConfigMountsTypeGcs, SandboxBoxNewParamsMountConfigMountsTypeGit:
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
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
