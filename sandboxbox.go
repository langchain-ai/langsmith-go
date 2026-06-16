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
	Mounts                 []SandboxBoxNewResponseMount     `json:"mounts"`
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
	Mounts                 apijson.Field
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

type SandboxBoxNewResponseMount struct {
	ID        string                          `json:"id" api:"required"`
	MountPath string                          `json:"mount_path" api:"required"`
	Type      SandboxBoxNewResponseMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                    `json:"s3"`
	JSON  sandboxBoxNewResponseMountJSON `json:"-"`
	union SandboxBoxNewResponseMountsUnion
}

// sandboxBoxNewResponseMountJSON contains the JSON metadata for the struct
// [SandboxBoxNewResponseMount]
type sandboxBoxNewResponseMountJSON struct {
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

func (r sandboxBoxNewResponseMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxNewResponseMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxNewResponseMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxNewResponseMountsUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxNewResponseMount) AsUnion() SandboxBoxNewResponseMountsUnion {
	return r.union
}

// Union satisfied by [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec].
type SandboxBoxNewResponseMountsUnion interface {
	implementsSandboxBoxNewResponseMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxNewResponseMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                      `json:"id" api:"required"`
	MountPath string                                                      `json:"mount_path" api:"required"`
	S3        SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                        `json:"read_only"`
	JSON      sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecJSON contains the JSON
// metadata for the struct [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec]
type sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecJSON struct {
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

func (r *SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxNewResponseMount() {
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                       `json:"bucket" api:"required"`
	EndpointURL string                                                       `json:"endpoint_url" api:"required"`
	Region      string                                                       `json:"region" api:"required"`
	PathStyle   bool                                                         `json:"path_style"`
	Prefix      string                                                       `json:"prefix"`
	JSON        sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                           `json:"max_size_bytes"`
	WritebackSeconds int64                                                           `json:"writeback_seconds"`
	JSON             sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                        `json:"bucket" api:"required"`
	Prefix string                                                        `json:"prefix"`
	JSON   sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                        `json:"remote_url" api:"required"`
	Ref                    SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                         `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                           `json:"name" api:"required"`
	Type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxNewResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                       `json:"id" api:"required"`
	Gcs       SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                       `json:"mount_path" api:"required"`
	Type      SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                         `json:"read_only"`
	S3        SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecJSON struct {
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

func (r *SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxNewResponseMount() {
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                         `json:"bucket" api:"required"`
	Prefix string                                                         `json:"prefix"`
	JSON   sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                            `json:"max_size_bytes"`
	WritebackSeconds int64                                                            `json:"writeback_seconds"`
	JSON             sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                         `json:"remote_url" api:"required"`
	Ref                    SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                          `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                            `json:"name" api:"required"`
	Type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                        `json:"bucket" api:"required"`
	EndpointURL string                                                        `json:"endpoint_url" api:"required"`
	Region      string                                                        `json:"region" api:"required"`
	PathStyle   bool                                                          `json:"path_style"`
	Prefix      string                                                        `json:"prefix"`
	JSON        sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                     `json:"id" api:"required"`
	Git       SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                     `json:"mount_path" api:"required"`
	Type      SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                       `json:"read_only"`
	S3        SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecJSON contains the JSON
// metadata for the struct [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec]
type sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecJSON struct {
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

func (r *SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxNewResponseMount() {
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                       `json:"remote_url" api:"required"`
	Ref                    SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                        `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                          `json:"name" api:"required"`
	Type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                          `json:"max_size_bytes"`
	WritebackSeconds int64                                                          `json:"writeback_seconds"`
	JSON             sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCacheJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                       `json:"bucket" api:"required"`
	Prefix string                                                       `json:"prefix"`
	JSON   sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                      `json:"bucket" api:"required"`
	EndpointURL string                                                      `json:"endpoint_url" api:"required"`
	Region      string                                                      `json:"region" api:"required"`
	PathStyle   bool                                                        `json:"path_style"`
	Prefix      string                                                      `json:"prefix"`
	JSON        sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxNewResponseMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxNewResponseMountsType string

const (
	SandboxBoxNewResponseMountsTypeS3  SandboxBoxNewResponseMountsType = "s3"
	SandboxBoxNewResponseMountsTypeGcs SandboxBoxNewResponseMountsType = "gcs"
	SandboxBoxNewResponseMountsTypeGit SandboxBoxNewResponseMountsType = "git"
)

func (r SandboxBoxNewResponseMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxNewResponseMountsTypeS3, SandboxBoxNewResponseMountsTypeGcs, SandboxBoxNewResponseMountsTypeGit:
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
	Name       string                                        `json:"name" api:"required"`
	Aws        SandboxBoxNewResponseProxyConfigRulesAws      `json:"aws"`
	Enabled    bool                                          `json:"enabled"`
	Gcp        SandboxBoxNewResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers    []SandboxBoxNewResponseProxyConfigRulesHeader `json:"headers"`
	MatchHosts []string                                      `json:"match_hosts"`
	MatchPaths []string                                      `json:"match_paths"`
	Type       string                                        `json:"type"`
	JSON       sandboxBoxNewResponseProxyConfigRuleJSON      `json:"-"`
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
	Mounts                 []SandboxBoxGetResponseMount     `json:"mounts"`
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
	Mounts                 apijson.Field
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

type SandboxBoxGetResponseMount struct {
	ID        string                          `json:"id" api:"required"`
	MountPath string                          `json:"mount_path" api:"required"`
	Type      SandboxBoxGetResponseMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                    `json:"s3"`
	JSON  sandboxBoxGetResponseMountJSON `json:"-"`
	union SandboxBoxGetResponseMountsUnion
}

// sandboxBoxGetResponseMountJSON contains the JSON metadata for the struct
// [SandboxBoxGetResponseMount]
type sandboxBoxGetResponseMountJSON struct {
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

func (r sandboxBoxGetResponseMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxGetResponseMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxGetResponseMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxGetResponseMountsUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxGetResponseMount) AsUnion() SandboxBoxGetResponseMountsUnion {
	return r.union
}

// Union satisfied by [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec].
type SandboxBoxGetResponseMountsUnion interface {
	implementsSandboxBoxGetResponseMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxGetResponseMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                      `json:"id" api:"required"`
	MountPath string                                                      `json:"mount_path" api:"required"`
	S3        SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                        `json:"read_only"`
	JSON      sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecJSON contains the JSON
// metadata for the struct [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec]
type sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecJSON struct {
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

func (r *SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxGetResponseMount() {
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                       `json:"bucket" api:"required"`
	EndpointURL string                                                       `json:"endpoint_url" api:"required"`
	Region      string                                                       `json:"region" api:"required"`
	PathStyle   bool                                                         `json:"path_style"`
	Prefix      string                                                       `json:"prefix"`
	JSON        sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                           `json:"max_size_bytes"`
	WritebackSeconds int64                                                           `json:"writeback_seconds"`
	JSON             sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                        `json:"bucket" api:"required"`
	Prefix string                                                        `json:"prefix"`
	JSON   sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                        `json:"remote_url" api:"required"`
	Ref                    SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                         `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                           `json:"name" api:"required"`
	Type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxGetResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                       `json:"id" api:"required"`
	Gcs       SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                       `json:"mount_path" api:"required"`
	Type      SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                         `json:"read_only"`
	S3        SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecJSON struct {
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

func (r *SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxGetResponseMount() {
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                         `json:"bucket" api:"required"`
	Prefix string                                                         `json:"prefix"`
	JSON   sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                            `json:"max_size_bytes"`
	WritebackSeconds int64                                                            `json:"writeback_seconds"`
	JSON             sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                         `json:"remote_url" api:"required"`
	Ref                    SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                          `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                            `json:"name" api:"required"`
	Type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                        `json:"bucket" api:"required"`
	EndpointURL string                                                        `json:"endpoint_url" api:"required"`
	Region      string                                                        `json:"region" api:"required"`
	PathStyle   bool                                                          `json:"path_style"`
	Prefix      string                                                        `json:"prefix"`
	JSON        sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                     `json:"id" api:"required"`
	Git       SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                     `json:"mount_path" api:"required"`
	Type      SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                       `json:"read_only"`
	S3        SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecJSON contains the JSON
// metadata for the struct [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec]
type sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecJSON struct {
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

func (r *SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxGetResponseMount() {
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                       `json:"remote_url" api:"required"`
	Ref                    SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                        `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                          `json:"name" api:"required"`
	Type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                          `json:"max_size_bytes"`
	WritebackSeconds int64                                                          `json:"writeback_seconds"`
	JSON             sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCacheJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                       `json:"bucket" api:"required"`
	Prefix string                                                       `json:"prefix"`
	JSON   sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                      `json:"bucket" api:"required"`
	EndpointURL string                                                      `json:"endpoint_url" api:"required"`
	Region      string                                                      `json:"region" api:"required"`
	PathStyle   bool                                                        `json:"path_style"`
	Prefix      string                                                      `json:"prefix"`
	JSON        sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxGetResponseMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxGetResponseMountsType string

const (
	SandboxBoxGetResponseMountsTypeS3  SandboxBoxGetResponseMountsType = "s3"
	SandboxBoxGetResponseMountsTypeGcs SandboxBoxGetResponseMountsType = "gcs"
	SandboxBoxGetResponseMountsTypeGit SandboxBoxGetResponseMountsType = "git"
)

func (r SandboxBoxGetResponseMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxGetResponseMountsTypeS3, SandboxBoxGetResponseMountsTypeGcs, SandboxBoxGetResponseMountsTypeGit:
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
	Name       string                                        `json:"name" api:"required"`
	Aws        SandboxBoxGetResponseProxyConfigRulesAws      `json:"aws"`
	Enabled    bool                                          `json:"enabled"`
	Gcp        SandboxBoxGetResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers    []SandboxBoxGetResponseProxyConfigRulesHeader `json:"headers"`
	MatchHosts []string                                      `json:"match_hosts"`
	MatchPaths []string                                      `json:"match_paths"`
	Type       string                                        `json:"type"`
	JSON       sandboxBoxGetResponseProxyConfigRuleJSON      `json:"-"`
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
	Mounts                 []SandboxBoxUpdateResponseMount     `json:"mounts"`
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
	Mounts                 apijson.Field
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

type SandboxBoxUpdateResponseMount struct {
	ID        string                             `json:"id" api:"required"`
	MountPath string                             `json:"mount_path" api:"required"`
	Type      SandboxBoxUpdateResponseMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                       `json:"s3"`
	JSON  sandboxBoxUpdateResponseMountJSON `json:"-"`
	union SandboxBoxUpdateResponseMountsUnion
}

// sandboxBoxUpdateResponseMountJSON contains the JSON metadata for the struct
// [SandboxBoxUpdateResponseMount]
type sandboxBoxUpdateResponseMountJSON struct {
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

func (r sandboxBoxUpdateResponseMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxUpdateResponseMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxUpdateResponseMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxUpdateResponseMountsUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxUpdateResponseMount) AsUnion() SandboxBoxUpdateResponseMountsUnion {
	return r.union
}

// Union satisfied by [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec].
type SandboxBoxUpdateResponseMountsUnion interface {
	implementsSandboxBoxUpdateResponseMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxUpdateResponseMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                         `json:"id" api:"required"`
	MountPath string                                                         `json:"mount_path" api:"required"`
	S3        SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                           `json:"read_only"`
	JSON      sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec]
type sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecJSON struct {
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

func (r *SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxUpdateResponseMount() {
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                          `json:"bucket" api:"required"`
	EndpointURL string                                                          `json:"endpoint_url" api:"required"`
	Region      string                                                          `json:"region" api:"required"`
	PathStyle   bool                                                            `json:"path_style"`
	Prefix      string                                                          `json:"prefix"`
	JSON        sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                              `json:"max_size_bytes"`
	WritebackSeconds int64                                                              `json:"writeback_seconds"`
	JSON             sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                           `json:"bucket" api:"required"`
	Prefix string                                                           `json:"prefix"`
	JSON   sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                           `json:"remote_url" api:"required"`
	Ref                    SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                            `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                              `json:"name" api:"required"`
	Type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxUpdateResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                          `json:"id" api:"required"`
	Gcs       SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                          `json:"mount_path" api:"required"`
	Type      SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                            `json:"read_only"`
	S3        SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecJSON struct {
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

func (r *SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxUpdateResponseMount() {
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                            `json:"bucket" api:"required"`
	Prefix string                                                            `json:"prefix"`
	JSON   sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                               `json:"max_size_bytes"`
	WritebackSeconds int64                                                               `json:"writeback_seconds"`
	JSON             sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                            `json:"remote_url" api:"required"`
	Ref                    SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                             `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                               `json:"name" api:"required"`
	Type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON contains
// the JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                           `json:"bucket" api:"required"`
	EndpointURL string                                                           `json:"endpoint_url" api:"required"`
	Region      string                                                           `json:"region" api:"required"`
	PathStyle   bool                                                             `json:"path_style"`
	Prefix      string                                                           `json:"prefix"`
	JSON        sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                        `json:"id" api:"required"`
	Git       SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                        `json:"mount_path" api:"required"`
	Type      SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                          `json:"read_only"`
	S3        SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec]
type sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecJSON struct {
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

func (r *SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxUpdateResponseMount() {
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                          `json:"remote_url" api:"required"`
	Ref                    SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                           `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                             `json:"name" api:"required"`
	Type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                             `json:"max_size_bytes"`
	WritebackSeconds int64                                                             `json:"writeback_seconds"`
	JSON             sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                          `json:"bucket" api:"required"`
	Prefix string                                                          `json:"prefix"`
	JSON   sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                         `json:"bucket" api:"required"`
	EndpointURL string                                                         `json:"endpoint_url" api:"required"`
	Region      string                                                         `json:"region" api:"required"`
	PathStyle   bool                                                           `json:"path_style"`
	Prefix      string                                                         `json:"prefix"`
	JSON        sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxUpdateResponseMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxUpdateResponseMountsType string

const (
	SandboxBoxUpdateResponseMountsTypeS3  SandboxBoxUpdateResponseMountsType = "s3"
	SandboxBoxUpdateResponseMountsTypeGcs SandboxBoxUpdateResponseMountsType = "gcs"
	SandboxBoxUpdateResponseMountsTypeGit SandboxBoxUpdateResponseMountsType = "git"
)

func (r SandboxBoxUpdateResponseMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxUpdateResponseMountsTypeS3, SandboxBoxUpdateResponseMountsTypeGcs, SandboxBoxUpdateResponseMountsTypeGit:
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
	Name       string                                           `json:"name" api:"required"`
	Aws        SandboxBoxUpdateResponseProxyConfigRulesAws      `json:"aws"`
	Enabled    bool                                             `json:"enabled"`
	Gcp        SandboxBoxUpdateResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers    []SandboxBoxUpdateResponseProxyConfigRulesHeader `json:"headers"`
	MatchHosts []string                                         `json:"match_hosts"`
	MatchPaths []string                                         `json:"match_paths"`
	Type       string                                           `json:"type"`
	JSON       sandboxBoxUpdateResponseProxyConfigRuleJSON      `json:"-"`
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
	Mounts                 []SandboxBoxListResponseSandboxesMount     `json:"mounts"`
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
	Mounts                 apijson.Field
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

type SandboxBoxListResponseSandboxesMount struct {
	ID        string                                    `json:"id" api:"required"`
	MountPath string                                    `json:"mount_path" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                              `json:"s3"`
	JSON  sandboxBoxListResponseSandboxesMountJSON `json:"-"`
	union SandboxBoxListResponseSandboxesMountsUnion
}

// sandboxBoxListResponseSandboxesMountJSON contains the JSON metadata for the
// struct [SandboxBoxListResponseSandboxesMount]
type sandboxBoxListResponseSandboxesMountJSON struct {
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

func (r sandboxBoxListResponseSandboxesMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxListResponseSandboxesMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxListResponseSandboxesMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxListResponseSandboxesMountsUnion] interface which
// you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxListResponseSandboxesMount) AsUnion() SandboxBoxListResponseSandboxesMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec].
type SandboxBoxListResponseSandboxesMountsUnion interface {
	implementsSandboxBoxListResponseSandboxesMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxListResponseSandboxesMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                                `json:"id" api:"required"`
	MountPath string                                                                `json:"mount_path" api:"required"`
	S3        SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                  `json:"read_only"`
	JSON      sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec]
type sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecJSON struct {
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

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxListResponseSandboxesMount() {
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                                 `json:"bucket" api:"required"`
	EndpointURL string                                                                 `json:"endpoint_url" api:"required"`
	Region      string                                                                 `json:"region" api:"required"`
	PathStyle   bool                                                                   `json:"path_style"`
	Prefix      string                                                                 `json:"prefix"`
	JSON        sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                     `json:"max_size_bytes"`
	WritebackSeconds int64                                                                     `json:"writeback_seconds"`
	JSON             sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                                  `json:"bucket" api:"required"`
	Prefix string                                                                  `json:"prefix"`
	JSON   sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcsJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                                  `json:"remote_url" api:"required"`
	Ref                    SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                   `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                     `json:"name" api:"required"`
	Type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxListResponseSandboxesMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                                 `json:"id" api:"required"`
	Gcs       SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                                 `json:"mount_path" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                                   `json:"read_only"`
	S3        SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecJSON struct {
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

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxListResponseSandboxesMount() {
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                                   `json:"bucket" api:"required"`
	Prefix string                                                                   `json:"prefix"`
	JSON   sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                      `json:"max_size_bytes"`
	WritebackSeconds int64                                                                      `json:"writeback_seconds"`
	JSON             sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                                   `json:"remote_url" api:"required"`
	Ref                    SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                    `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                      `json:"name" api:"required"`
	Type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                                  `json:"bucket" api:"required"`
	EndpointURL string                                                                  `json:"endpoint_url" api:"required"`
	Region      string                                                                  `json:"region" api:"required"`
	PathStyle   bool                                                                    `json:"path_style"`
	Prefix      string                                                                  `json:"prefix"`
	JSON        sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                               `json:"id" api:"required"`
	Git       SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                               `json:"mount_path" api:"required"`
	Type      SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                                 `json:"read_only"`
	S3        SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecJSON contains the
// JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec]
type sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecJSON struct {
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

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxListResponseSandboxesMount() {
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                                 `json:"remote_url" api:"required"`
	Ref                    SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                  `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                                    `json:"name" api:"required"`
	Type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                    `json:"max_size_bytes"`
	WritebackSeconds int64                                                                    `json:"writeback_seconds"`
	JSON             sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                                 `json:"bucket" api:"required"`
	Prefix string                                                                 `json:"prefix"`
	JSON   sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcsJSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                                `json:"bucket" api:"required"`
	EndpointURL string                                                                `json:"endpoint_url" api:"required"`
	Region      string                                                                `json:"region" api:"required"`
	PathStyle   bool                                                                  `json:"path_style"`
	Prefix      string                                                                `json:"prefix"`
	JSON        sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxListResponseSandboxesMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxListResponseSandboxesMountsType string

const (
	SandboxBoxListResponseSandboxesMountsTypeS3  SandboxBoxListResponseSandboxesMountsType = "s3"
	SandboxBoxListResponseSandboxesMountsTypeGcs SandboxBoxListResponseSandboxesMountsType = "gcs"
	SandboxBoxListResponseSandboxesMountsTypeGit SandboxBoxListResponseSandboxesMountsType = "git"
)

func (r SandboxBoxListResponseSandboxesMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxListResponseSandboxesMountsTypeS3, SandboxBoxListResponseSandboxesMountsTypeGcs, SandboxBoxListResponseSandboxesMountsTypeGit:
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
	Name       string                                                  `json:"name" api:"required"`
	Aws        SandboxBoxListResponseSandboxesProxyConfigRulesAws      `json:"aws"`
	Enabled    bool                                                    `json:"enabled"`
	Gcp        SandboxBoxListResponseSandboxesProxyConfigRulesGcp      `json:"gcp"`
	Headers    []SandboxBoxListResponseSandboxesProxyConfigRulesHeader `json:"headers"`
	MatchHosts []string                                                `json:"match_hosts"`
	MatchPaths []string                                                `json:"match_paths"`
	Type       string                                                  `json:"type"`
	JSON       sandboxBoxListResponseSandboxesProxyConfigRuleJSON      `json:"-"`
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
	Mounts                 []SandboxBoxStartResponseMount     `json:"mounts"`
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
	Mounts                 apijson.Field
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

type SandboxBoxStartResponseMount struct {
	ID        string                            `json:"id" api:"required"`
	MountPath string                            `json:"mount_path" api:"required"`
	Type      SandboxBoxStartResponseMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3].
	S3    interface{}                      `json:"s3"`
	JSON  sandboxBoxStartResponseMountJSON `json:"-"`
	union SandboxBoxStartResponseMountsUnion
}

// sandboxBoxStartResponseMountJSON contains the JSON metadata for the struct
// [SandboxBoxStartResponseMount]
type sandboxBoxStartResponseMountJSON struct {
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

func (r sandboxBoxStartResponseMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxBoxStartResponseMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxBoxStartResponseMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxBoxStartResponseMountsUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec].
func (r SandboxBoxStartResponseMount) AsUnion() SandboxBoxStartResponseMountsUnion {
	return r.union
}

// Union satisfied by [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec] or
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec].
type SandboxBoxStartResponseMountsUnion interface {
	implementsSandboxBoxStartResponseMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxBoxStartResponseMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
	)
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec struct {
	ID        string                                                        `json:"id" api:"required"`
	MountPath string                                                        `json:"mount_path" api:"required"`
	S3        SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3    `json:"s3" api:"required"`
	Type      SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCache `json:"cache"`
	Gcs       SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcs   `json:"gcs"`
	Git       SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                          `json:"read_only"`
	JSON      sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecJSON  `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec]
type sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecJSON struct {
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

func (r *SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxStartResponseMount() {
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                         `json:"bucket" api:"required"`
	EndpointURL string                                                         `json:"endpoint_url" api:"required"`
	Region      string                                                         `json:"region" api:"required"`
	PathStyle   bool                                                           `json:"path_style"`
	Prefix      string                                                         `json:"prefix"`
	JSON        sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3]
type sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                             `json:"max_size_bytes"`
	WritebackSeconds int64                                                             `json:"writeback_seconds"`
	JSON             sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCache]
type sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                          `json:"bucket" api:"required"`
	Prefix string                                                          `json:"prefix"`
	JSON   sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcs]
type sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                          `json:"remote_url" api:"required"`
	Ref                    SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                           `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGit]
type sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                             `json:"name" api:"required"`
	Type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxStartResponseMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec struct {
	ID        string                                                         `json:"id" api:"required"`
	Gcs       SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcs   `json:"gcs" api:"required"`
	MountPath string                                                         `json:"mount_path" api:"required"`
	Type      SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCache `json:"cache"`
	Git       SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGit   `json:"git"`
	ReadOnly  bool                                                           `json:"read_only"`
	S3        SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3    `json:"s3"`
	JSON      sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecJSON  `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec]
type sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecJSON struct {
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

func (r *SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxStartResponseMount() {
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                           `json:"bucket" api:"required"`
	Prefix string                                                           `json:"prefix"`
	JSON   sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                              `json:"max_size_bytes"`
	WritebackSeconds int64                                                              `json:"writeback_seconds"`
	JSON             sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCache]
type sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                           `json:"remote_url" api:"required"`
	Ref                    SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                            `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGit]
type sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                              `json:"name" api:"required"`
	Type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                          `json:"bucket" api:"required"`
	EndpointURL string                                                          `json:"endpoint_url" api:"required"`
	Region      string                                                          `json:"region" api:"required"`
	PathStyle   bool                                                            `json:"path_style"`
	Prefix      string                                                          `json:"prefix"`
	JSON        sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3]
type sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec struct {
	ID        string                                                       `json:"id" api:"required"`
	Git       SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGit   `json:"git" api:"required"`
	MountPath string                                                       `json:"mount_path" api:"required"`
	Type      SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecType  `json:"type" api:"required"`
	Cache     SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCache `json:"cache"`
	Gcs       SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcs   `json:"gcs"`
	ReadOnly  bool                                                         `json:"read_only"`
	S3        SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3    `json:"s3"`
	JSON      sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecJSON  `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec]
type sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecJSON struct {
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

func (r *SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxStartResponseMount() {
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                         `json:"remote_url" api:"required"`
	Ref                    SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                          `json:"refresh_interval_seconds"`
	JSON                   sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGit]
type sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                            `json:"name" api:"required"`
	Type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                            `json:"max_size_bytes"`
	WritebackSeconds int64                                                            `json:"writeback_seconds"`
	JSON             sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCache]
type sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                         `json:"bucket" api:"required"`
	Prefix string                                                         `json:"prefix"`
	JSON   sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcs]
type sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                        `json:"bucket" api:"required"`
	EndpointURL string                                                        `json:"endpoint_url" api:"required"`
	Region      string                                                        `json:"region" api:"required"`
	PathStyle   bool                                                          `json:"path_style"`
	Prefix      string                                                        `json:"prefix"`
	JSON        sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3JSON contains the JSON
// metadata for the struct
// [SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3]
type sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	EndpointURL apijson.Field
	Region      apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBoxStartResponseMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxBoxStartResponseMountsType string

const (
	SandboxBoxStartResponseMountsTypeS3  SandboxBoxStartResponseMountsType = "s3"
	SandboxBoxStartResponseMountsTypeGcs SandboxBoxStartResponseMountsType = "gcs"
	SandboxBoxStartResponseMountsTypeGit SandboxBoxStartResponseMountsType = "git"
)

func (r SandboxBoxStartResponseMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxStartResponseMountsTypeS3, SandboxBoxStartResponseMountsTypeGcs, SandboxBoxStartResponseMountsTypeGit:
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
	Name       string                                          `json:"name" api:"required"`
	Aws        SandboxBoxStartResponseProxyConfigRulesAws      `json:"aws"`
	Enabled    bool                                            `json:"enabled"`
	Gcp        SandboxBoxStartResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers    []SandboxBoxStartResponseProxyConfigRulesHeader `json:"headers"`
	MatchHosts []string                                        `json:"match_hosts"`
	MatchPaths []string                                        `json:"match_paths"`
	Type       string                                          `json:"type"`
	JSON       sandboxBoxStartResponseProxyConfigRuleJSON      `json:"-"`
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
	DeleteAfterStopSeconds param.Field[int64]                           `json:"delete_after_stop_seconds"`
	EnvVars                param.Field[map[string]string]               `json:"env_vars"`
	FsCapacityBytes        param.Field[int64]                           `json:"fs_capacity_bytes"`
	IdleTtlSeconds         param.Field[int64]                           `json:"idle_ttl_seconds"`
	MemBytes               param.Field[int64]                           `json:"mem_bytes"`
	Mounts                 param.Field[[]SandboxBoxNewParamsMountUnion] `json:"mounts"`
	Name                   param.Field[string]                          `json:"name"`
	ProxyConfig            param.Field[SandboxBoxNewParamsProxyConfig]  `json:"proxy_config"`
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

type SandboxBoxNewParamsMount struct {
	ID        param.Field[string]                        `json:"id" api:"required"`
	MountPath param.Field[string]                        `json:"mount_path" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountsType] `json:"type" api:"required"`
	Cache     param.Field[interface{}]                   `json:"cache"`
	Gcs       param.Field[interface{}]                   `json:"gcs"`
	Git       param.Field[interface{}]                   `json:"git"`
	ReadOnly  param.Field[bool]                          `json:"read_only"`
	S3        param.Field[interface{}]                   `json:"s3"`
}

func (r SandboxBoxNewParamsMount) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMount) implementsSandboxBoxNewParamsMountUnion() {}

// Satisfied by [SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpec],
// [SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpec],
// [SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpec],
// [SandboxBoxNewParamsMount].
type SandboxBoxNewParamsMountUnion interface {
	implementsSandboxBoxNewParamsMountUnion()
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpec struct {
	ID        param.Field[string]                                                    `json:"id" api:"required"`
	MountPath param.Field[string]                                                    `json:"mount_path" api:"required"`
	S3        param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecS3]    `json:"s3" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecType]  `json:"type" api:"required"`
	Cache     param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecCache] `json:"cache"`
	Gcs       param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGcs]   `json:"gcs"`
	Git       param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGit]   `json:"git"`
	ReadOnly  param.Field[bool]                                                      `json:"read_only"`
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpec) implementsSandboxBoxNewParamsMountUnion() {
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeS3  SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeGcs SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeGit SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecType = "git"
)

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeS3, SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              param.Field[string]                                                     `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                      `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name param.Field[string]                                                         `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpec struct {
	ID        param.Field[string]                                                     `json:"id" api:"required"`
	Gcs       param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGcs]   `json:"gcs" api:"required"`
	MountPath param.Field[string]                                                     `json:"mount_path" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecType]  `json:"type" api:"required"`
	Cache     param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecCache] `json:"cache"`
	Git       param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGit]   `json:"git"`
	ReadOnly  param.Field[bool]                                                       `json:"read_only"`
	S3        param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecS3]    `json:"s3"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpec) implementsSandboxBoxNewParamsMountUnion() {
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecTypeS3  SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecTypeGcs SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecTypeGit SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecType = "git"
)

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              param.Field[string]                                                      `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                       `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name param.Field[string]                                                          `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGcsBucketMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpec struct {
	ID        param.Field[string]                                                   `json:"id" api:"required"`
	Git       param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGit]   `json:"git" api:"required"`
	MountPath param.Field[string]                                                   `json:"mount_path" api:"required"`
	Type      param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecType]  `json:"type" api:"required"`
	Cache     param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecCache] `json:"cache"`
	Gcs       param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGcs]   `json:"gcs"`
	ReadOnly  param.Field[bool]                                                     `json:"read_only"`
	S3        param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecS3]    `json:"s3"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpec) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpec) implementsSandboxBoxNewParamsMountUnion() {
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              param.Field[string]                                                    `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRef] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                     `json:"refresh_interval_seconds"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGit) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name param.Field[string]                                                        `json:"name" api:"required"`
	Type param.Field[SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRef) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecTypeS3  SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecTypeGcs SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecTypeGit SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecType = "git"
)

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecTypeS3, SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecTypeGit:
		return true
	}
	return false
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecCache) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecGcs) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxBoxNewParamsMountsSandboxapiGitRepoMountSpecS3) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBoxNewParamsMountsType string

const (
	SandboxBoxNewParamsMountsTypeS3  SandboxBoxNewParamsMountsType = "s3"
	SandboxBoxNewParamsMountsTypeGcs SandboxBoxNewParamsMountsType = "gcs"
	SandboxBoxNewParamsMountsTypeGit SandboxBoxNewParamsMountsType = "git"
)

func (r SandboxBoxNewParamsMountsType) IsKnown() bool {
	switch r {
	case SandboxBoxNewParamsMountsTypeS3, SandboxBoxNewParamsMountsTypeGcs, SandboxBoxNewParamsMountsTypeGit:
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
	Name       param.Field[string]                                      `json:"name" api:"required"`
	Aws        param.Field[SandboxBoxNewParamsProxyConfigRulesAws]      `json:"aws"`
	Enabled    param.Field[bool]                                        `json:"enabled"`
	Gcp        param.Field[SandboxBoxNewParamsProxyConfigRulesGcp]      `json:"gcp"`
	Headers    param.Field[[]SandboxBoxNewParamsProxyConfigRulesHeader] `json:"headers"`
	MatchHosts param.Field[[]string]                                    `json:"match_hosts"`
	MatchPaths param.Field[[]string]                                    `json:"match_paths"`
	Type       param.Field[string]                                      `json:"type"`
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
	Name       param.Field[string]                                         `json:"name" api:"required"`
	Aws        param.Field[SandboxBoxUpdateParamsProxyConfigRulesAws]      `json:"aws"`
	Enabled    param.Field[bool]                                           `json:"enabled"`
	Gcp        param.Field[SandboxBoxUpdateParamsProxyConfigRulesGcp]      `json:"gcp"`
	Headers    param.Field[[]SandboxBoxUpdateParamsProxyConfigRulesHeader] `json:"headers"`
	MatchHosts param.Field[[]string]                                       `json:"match_hosts"`
	MatchPaths param.Field[[]string]                                       `json:"match_paths"`
	Type       param.Field[string]                                         `json:"type"`
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
