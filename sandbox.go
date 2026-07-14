// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"reflect"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/tidwall/gjson"
)

// SandboxService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxService] method instead.
type SandboxService struct {
	Options    []option.RequestOption
	Boxes      *SandboxBoxService
	Registries *SandboxRegistryService
	Snapshots  *SandboxSnapshotService
}

// NewSandboxService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewSandboxService(opts ...option.RequestOption) (r *SandboxService) {
	r = &SandboxService{}
	r.Options = opts
	r.Boxes = NewSandboxBoxService(opts...)
	r.Registries = NewSandboxRegistryService(opts...)
	r.Snapshots = NewSandboxSnapshotService(opts...)
	return
}

type SandboxListResponse struct {
	Offset    int64                   `json:"offset"`
	Sandboxes []SandboxResponse       `json:"sandboxes"`
	JSON      sandboxListResponseJSON `json:"-"`
}

// sandboxListResponseJSON contains the JSON metadata for the struct
// [SandboxListResponse]
type sandboxListResponseJSON struct {
	Offset      apijson.Field
	Sandboxes   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxListResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxResponse struct {
	ID                     string                     `json:"id"`
	CPUMillicores          int64                      `json:"cpu_millicores"`
	CreatedAt              string                     `json:"created_at"`
	CreatedBy              string                     `json:"created_by"`
	DataplaneURL           string                     `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                      `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                      `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                      `json:"idle_ttl_seconds"`
	MemBytes               int64                      `json:"mem_bytes"`
	MountConfig            SandboxResponseMountConfig `json:"mount_config"`
	Name                   string                     `json:"name"`
	ProxyConfig            SandboxResponseProxyConfig `json:"proxy_config"`
	SizeClass              string                     `json:"size_class"`
	SnapshotID             string                     `json:"snapshot_id"`
	Status                 string                     `json:"status"`
	StatusMessage          string                     `json:"status_message"`
	StoppedAt              string                     `json:"stopped_at"`
	UpdatedAt              string                     `json:"updated_at"`
	UpdatedBy              string                     `json:"updated_by"`
	Vcpus                  int64                      `json:"vcpus"`
	JSON                   sandboxResponseJSON        `json:"-"`
}

// sandboxResponseJSON contains the JSON metadata for the struct [SandboxResponse]
type sandboxResponseJSON struct {
	ID                     apijson.Field
	CPUMillicores          apijson.Field
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

func (r *SandboxResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfig struct {
	Auth   SandboxResponseMountConfigAuth    `json:"auth"`
	Mounts []SandboxResponseMountConfigMount `json:"mounts"`
	JSON   sandboxResponseMountConfigJSON    `json:"-"`
}

// sandboxResponseMountConfigJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfig]
type sandboxResponseMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuth struct {
	Aws  SandboxResponseMountConfigAuthAws  `json:"aws"`
	Gcp  SandboxResponseMountConfigAuthGcp  `json:"gcp"`
	JSON sandboxResponseMountConfigAuthJSON `json:"-"`
}

// sandboxResponseMountConfigAuthJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuth]
type sandboxResponseMountConfigAuthJSON struct {
	Aws         apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAws struct {
	AccessKeyID     SandboxResponseMountConfigAuthAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxResponseMountConfigAuthAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxResponseMountConfigAuthAwsJSON            `json:"-"`
}

// sandboxResponseMountConfigAuthAwsJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAws]
type sandboxResponseMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAwsAccessKeyID struct {
	Type  SandboxResponseMountConfigAuthAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                             `json:"is_set"`
	Value string                                           `json:"value"`
	JSON  sandboxResponseMountConfigAuthAwsAccessKeyIDJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAwsAccessKeyIDJSON contains the JSON metadata for
// the struct [SandboxResponseMountConfigAuthAwsAccessKeyID]
type sandboxResponseMountConfigAuthAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAwsAccessKeyIDType string

const (
	SandboxResponseMountConfigAuthAwsAccessKeyIDTypePlaintext       SandboxResponseMountConfigAuthAwsAccessKeyIDType = "plaintext"
	SandboxResponseMountConfigAuthAwsAccessKeyIDTypeOpaque          SandboxResponseMountConfigAuthAwsAccessKeyIDType = "opaque"
	SandboxResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret SandboxResponseMountConfigAuthAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAwsAccessKeyIDTypePlaintext, SandboxResponseMountConfigAuthAwsAccessKeyIDTypeOpaque, SandboxResponseMountConfigAuthAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigAuthAwsSecretAccessKey struct {
	Type  SandboxResponseMountConfigAuthAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                 `json:"is_set"`
	Value string                                               `json:"value"`
	JSON  sandboxResponseMountConfigAuthAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAwsSecretAccessKeyJSON contains the JSON metadata
// for the struct [SandboxResponseMountConfigAuthAwsSecretAccessKey]
type sandboxResponseMountConfigAuthAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAwsSecretAccessKeyType string

const (
	SandboxResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext       SandboxResponseMountConfigAuthAwsSecretAccessKeyType = "plaintext"
	SandboxResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque          SandboxResponseMountConfigAuthAwsSecretAccessKeyType = "opaque"
	SandboxResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret SandboxResponseMountConfigAuthAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAwsSecretAccessKeyTypePlaintext, SandboxResponseMountConfigAuthAwsSecretAccessKeyTypeOpaque, SandboxResponseMountConfigAuthAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigAuthGcp struct {
	ServiceAccountJson SandboxResponseMountConfigAuthGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxResponseMountConfigAuthGcpJSON               `json:"-"`
}

// sandboxResponseMountConfigAuthGcpJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthGcp]
type sandboxResponseMountConfigAuthGcpJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthGcpServiceAccountJson struct {
	Type  SandboxResponseMountConfigAuthGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                    `json:"is_set"`
	Value string                                                  `json:"value"`
	JSON  sandboxResponseMountConfigAuthGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxResponseMountConfigAuthGcpServiceAccountJsonJSON contains the JSON
// metadata for the struct [SandboxResponseMountConfigAuthGcpServiceAccountJson]
type sandboxResponseMountConfigAuthGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthGcpServiceAccountJsonType string

const (
	SandboxResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext       SandboxResponseMountConfigAuthGcpServiceAccountJsonType = "plaintext"
	SandboxResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque          SandboxResponseMountConfigAuthGcpServiceAccountJsonType = "opaque"
	SandboxResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret SandboxResponseMountConfigAuthGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthGcpServiceAccountJsonTypePlaintext, SandboxResponseMountConfigAuthGcpServiceAccountJsonTypeOpaque, SandboxResponseMountConfigAuthGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigMount struct {
	ID        string                               `json:"id" api:"required"`
	MountPath string                               `json:"mount_path" api:"required"`
	Type      SandboxResponseMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub].
	Contexthub interface{} `json:"contexthub"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3].
	S3    interface{}                         `json:"s3"`
	JSON  sandboxResponseMountConfigMountJSON `json:"-"`
	union SandboxResponseMountConfigMountsUnion
}

// sandboxResponseMountConfigMountJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigMount]
type sandboxResponseMountConfigMountJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r sandboxResponseMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseMountConfigMountsUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec].
func (r SandboxResponseMountConfigMount) AsUnion() SandboxResponseMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec] or
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec].
type SandboxResponseMountConfigMountsUnion interface {
	implementsSandboxResponseMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
	)
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID         string                                                                `json:"id" api:"required"`
	MountPath  string                                                                `json:"mount_path" api:"required"`
	S3         SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3         `json:"s3" api:"required"`
	Type       SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType       `json:"type" api:"required"`
	Cache      SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub `json:"contexthub"`
	Gcs        SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs        `json:"gcs"`
	Git        SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGit        `json:"git"`
	ReadOnly   bool                                                                  `json:"read_only"`
	JSON       sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON       `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
	ID          apijson.Field
	MountPath   apijson.Field
	S3          apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxResponseMountConfigMount() {
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                            `json:"bucket" api:"required"`
	Region      string                                                            `json:"region" api:"required"`
	EndpointURL string                                                            `json:"endpoint_url"`
	PathStyle   bool                                                              `json:"path_style"`
	Prefix      string                                                            `json:"prefix"`
	JSON        sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3         SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs        SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                `json:"max_size_bytes"`
	WritebackSeconds int64                                                                `json:"writeback_seconds"`
	JSON             sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                      `json:"initial_pull_only"`
	JSON            sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                             `json:"bucket" api:"required"`
	Prefix string                                                             `json:"prefix"`
	JSON   sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                             `json:"remote_url" api:"required"`
	Ref                    SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                              `json:"refresh_interval_seconds"`
	JSON                   sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                                `json:"name" api:"required"`
	Type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID         string                                                                 `json:"id" api:"required"`
	Gcs        SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs        `json:"gcs" api:"required"`
	MountPath  string                                                                 `json:"mount_path" api:"required"`
	Type       SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType       `json:"type" api:"required"`
	Cache      SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub `json:"contexthub"`
	Git        SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit        `json:"git"`
	ReadOnly   bool                                                                   `json:"read_only"`
	S3         SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3         `json:"s3"`
	JSON       sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON       `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
	ID          apijson.Field
	Gcs         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxResponseMountConfigMount() {
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                              `json:"bucket" api:"required"`
	Prefix string                                                              `json:"prefix"`
	JSON   sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3         SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs        SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                                 `json:"max_size_bytes"`
	WritebackSeconds int64                                                                 `json:"writeback_seconds"`
	JSON             sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                       `json:"initial_pull_only"`
	JSON            sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                              `json:"remote_url" api:"required"`
	Ref                    SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                               `json:"refresh_interval_seconds"`
	JSON                   sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                                 `json:"name" api:"required"`
	Type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                             `json:"bucket" api:"required"`
	Region      string                                                             `json:"region" api:"required"`
	EndpointURL string                                                             `json:"endpoint_url"`
	PathStyle   bool                                                               `json:"path_style"`
	Prefix      string                                                             `json:"prefix"`
	JSON        sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID         string                                                               `json:"id" api:"required"`
	Git        SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGit        `json:"git" api:"required"`
	MountPath  string                                                               `json:"mount_path" api:"required"`
	Type       SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType       `json:"type" api:"required"`
	Cache      SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub `json:"contexthub"`
	Gcs        SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs        `json:"gcs"`
	ReadOnly   bool                                                                 `json:"read_only"`
	S3         SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3         `json:"s3"`
	JSON       sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON       `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
	ID          apijson.Field
	Git         apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Gcs         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxResponseMountConfigMount() {
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                            `json:"remote_url" api:"required"`
	Ref                    SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                             `json:"refresh_interval_seconds"`
	JSON                   sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                               `json:"name" api:"required"`
	Type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3         SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs        SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                               `json:"max_size_bytes"`
	WritebackSeconds int64                                                               `json:"writeback_seconds"`
	JSON             sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                     `json:"initial_pull_only"`
	JSON            sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                            `json:"bucket" api:"required"`
	Prefix string                                                            `json:"prefix"`
	JSON   sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                           `json:"bucket" api:"required"`
	Region      string                                                           `json:"region" api:"required"`
	EndpointURL string                                                           `json:"endpoint_url"`
	PathStyle   bool                                                             `json:"path_style"`
	Prefix      string                                                           `json:"prefix"`
	JSON        sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec struct {
	ID         string                                                                      `json:"id" api:"required"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub `json:"contexthub" api:"required"`
	MountPath  string                                                                      `json:"mount_path" api:"required"`
	Type       SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType       `json:"type" api:"required"`
	Cache      SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache      `json:"cache"`
	Gcs        SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcs        `json:"gcs"`
	Git        SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGit        `json:"git"`
	ReadOnly   bool                                                                        `json:"read_only"`
	S3         SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3         `json:"s3"`
	JSON       sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecJSON       `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecJSON struct {
	ID          apijson.Field
	Contexthub  apijson.Field
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

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec) implementsSandboxResponseMountConfigMount() {
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                            `json:"initial_pull_only"`
	JSON            sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType string

const (
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3         SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "s3"
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs        SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "gcs"
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                                      `json:"max_size_bytes"`
	WritebackSeconds int64                                                                      `json:"writeback_seconds"`
	JSON             sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcs struct {
	Bucket string                                                                   `json:"bucket" api:"required"`
	Prefix string                                                                   `json:"prefix"`
	JSON   sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcs]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGit struct {
	RemoteURL              string                                                                   `json:"remote_url" api:"required"`
	Ref                    SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                    `json:"refresh_interval_seconds"`
	JSON                   sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGit]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef struct {
	Name string                                                                      `json:"name" api:"required"`
	Type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType string

const (
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeBranch SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType = "branch"
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeTag    SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType = "tag"
)

func (r SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeBranch, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3 struct {
	Bucket      string                                                                  `json:"bucket" api:"required"`
	Region      string                                                                  `json:"region" api:"required"`
	EndpointURL string                                                                  `json:"endpoint_url"`
	PathStyle   bool                                                                    `json:"path_style"`
	Prefix      string                                                                  `json:"prefix"`
	JSON        sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsType string

const (
	SandboxResponseMountConfigMountsTypeS3         SandboxResponseMountConfigMountsType = "s3"
	SandboxResponseMountConfigMountsTypeGcs        SandboxResponseMountConfigMountsType = "gcs"
	SandboxResponseMountConfigMountsTypeGit        SandboxResponseMountConfigMountsType = "git"
	SandboxResponseMountConfigMountsTypeContexthub SandboxResponseMountConfigMountsType = "contexthub"
)

func (r SandboxResponseMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsTypeS3, SandboxResponseMountConfigMountsTypeGcs, SandboxResponseMountConfigMountsTypeGit, SandboxResponseMountConfigMountsTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseProxyConfig struct {
	AccessControl SandboxResponseProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxResponseProxyConfigCallback    `json:"callbacks"`
	NoProxy       []string                                `json:"no_proxy"`
	Rules         []SandboxResponseProxyConfigRule        `json:"rules"`
	JSON          sandboxResponseProxyConfigJSON          `json:"-"`
}

// sandboxResponseProxyConfigJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfig]
type sandboxResponseProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxResponseProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigAccessControl struct {
	AllowList []string                                    `json:"allow_list"`
	DenyList  []string                                    `json:"deny_list"`
	JSON      sandboxResponseProxyConfigAccessControlJSON `json:"-"`
}

// sandboxResponseProxyConfigAccessControlJSON contains the JSON metadata for the
// struct [SandboxResponseProxyConfigAccessControl]
type sandboxResponseProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigCallback struct {
	MatchHosts     []string                                           `json:"match_hosts" api:"required"`
	TtlSeconds     int64                                              `json:"ttl_seconds" api:"required"`
	URL            string                                             `json:"url" api:"required"`
	FullRequest    bool                                               `json:"full_request"`
	RequestHeaders []SandboxResponseProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           sandboxResponseProxyConfigCallbackJSON             `json:"-"`
}

// sandboxResponseProxyConfigCallbackJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfigCallback]
type sandboxResponseProxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigCallbacksRequestHeader struct {
	Name  string                                                `json:"name" api:"required"`
	Type  SandboxResponseProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                                  `json:"is_set"`
	Value string                                                `json:"value"`
	JSON  sandboxResponseProxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// sandboxResponseProxyConfigCallbacksRequestHeaderJSON contains the JSON metadata
// for the struct [SandboxResponseProxyConfigCallbacksRequestHeader]
type sandboxResponseProxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigCallbacksRequestHeadersType string

const (
	SandboxResponseProxyConfigCallbacksRequestHeadersTypePlaintext       SandboxResponseProxyConfigCallbacksRequestHeadersType = "plaintext"
	SandboxResponseProxyConfigCallbacksRequestHeadersTypeOpaque          SandboxResponseProxyConfigCallbacksRequestHeadersType = "opaque"
	SandboxResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret SandboxResponseProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r SandboxResponseProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigCallbacksRequestHeadersTypePlaintext, SandboxResponseProxyConfigCallbacksRequestHeadersTypeOpaque, SandboxResponseProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRule struct {
	Name    string                                  `json:"name" api:"required"`
	Aws     SandboxResponseProxyConfigRulesAws      `json:"aws"`
	Enabled bool                                    `json:"enabled"`
	Gcp     SandboxResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxResponseProxyConfigRulesHeader `json:"headers"`
	// MatchHosts is only accepted for header injection rules. Provider auth rules use
	// built-in host matching.
	MatchHosts []string                           `json:"match_hosts"`
	MatchPaths []string                           `json:"match_paths"`
	Type       string                             `json:"type"`
	JSON       sandboxResponseProxyConfigRuleJSON `json:"-"`
}

// sandboxResponseProxyConfigRuleJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRule]
type sandboxResponseProxyConfigRuleJSON struct {
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

func (r *SandboxResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAws struct {
	AccessKeyID     SandboxResponseProxyConfigRulesAwsAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxResponseProxyConfigRulesAwsSecretAccessKey `json:"secret_access_key" api:"required"`
	JSON            sandboxResponseProxyConfigRulesAwsJSON            `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAws]
type sandboxResponseProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAwsAccessKeyID struct {
	Type  SandboxResponseProxyConfigRulesAwsAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                              `json:"is_set"`
	Value string                                            `json:"value"`
	JSON  sandboxResponseProxyConfigRulesAwsAccessKeyIDJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsAccessKeyIDJSON contains the JSON metadata for
// the struct [SandboxResponseProxyConfigRulesAwsAccessKeyID]
type sandboxResponseProxyConfigRulesAwsAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAwsAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAwsAccessKeyIDType string

const (
	SandboxResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext       SandboxResponseProxyConfigRulesAwsAccessKeyIDType = "plaintext"
	SandboxResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque          SandboxResponseProxyConfigRulesAwsAccessKeyIDType = "opaque"
	SandboxResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret SandboxResponseProxyConfigRulesAwsAccessKeyIDType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesAwsAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAwsAccessKeyIDTypePlaintext, SandboxResponseProxyConfigRulesAwsAccessKeyIDTypeOpaque, SandboxResponseProxyConfigRulesAwsAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesAwsSecretAccessKey struct {
	Type  SandboxResponseProxyConfigRulesAwsSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                  `json:"is_set"`
	Value string                                                `json:"value"`
	JSON  sandboxResponseProxyConfigRulesAwsSecretAccessKeyJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsSecretAccessKeyJSON contains the JSON metadata
// for the struct [SandboxResponseProxyConfigRulesAwsSecretAccessKey]
type sandboxResponseProxyConfigRulesAwsSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAwsSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAwsSecretAccessKeyType string

const (
	SandboxResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext       SandboxResponseProxyConfigRulesAwsSecretAccessKeyType = "plaintext"
	SandboxResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque          SandboxResponseProxyConfigRulesAwsSecretAccessKeyType = "opaque"
	SandboxResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret SandboxResponseProxyConfigRulesAwsSecretAccessKeyType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesAwsSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAwsSecretAccessKeyTypePlaintext, SandboxResponseProxyConfigRulesAwsSecretAccessKeyTypeOpaque, SandboxResponseProxyConfigRulesAwsSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesGcp struct {
	Scopes             []string                                             `json:"scopes" api:"required"`
	ServiceAccountJson SandboxResponseProxyConfigRulesGcpServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxResponseProxyConfigRulesGcpJSON               `json:"-"`
}

// sandboxResponseProxyConfigRulesGcpJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesGcp]
type sandboxResponseProxyConfigRulesGcpJSON struct {
	Scopes             apijson.Field
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesGcpServiceAccountJson struct {
	Type  SandboxResponseProxyConfigRulesGcpServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                     `json:"is_set"`
	Value string                                                   `json:"value"`
	JSON  sandboxResponseProxyConfigRulesGcpServiceAccountJsonJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesGcpServiceAccountJsonJSON contains the JSON
// metadata for the struct [SandboxResponseProxyConfigRulesGcpServiceAccountJson]
type sandboxResponseProxyConfigRulesGcpServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesGcpServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesGcpServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesGcpServiceAccountJsonType string

const (
	SandboxResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext       SandboxResponseProxyConfigRulesGcpServiceAccountJsonType = "plaintext"
	SandboxResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque          SandboxResponseProxyConfigRulesGcpServiceAccountJsonType = "opaque"
	SandboxResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret SandboxResponseProxyConfigRulesGcpServiceAccountJsonType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesGcpServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesGcpServiceAccountJsonTypePlaintext, SandboxResponseProxyConfigRulesGcpServiceAccountJsonTypeOpaque, SandboxResponseProxyConfigRulesGcpServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesHeader struct {
	Name  string                                     `json:"name" api:"required"`
	Type  SandboxResponseProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                       `json:"is_set"`
	Value string                                     `json:"value"`
	JSON  sandboxResponseProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxResponseProxyConfigRulesHeaderJSON contains the JSON metadata for the
// struct [SandboxResponseProxyConfigRulesHeader]
type sandboxResponseProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesHeadersType string

const (
	SandboxResponseProxyConfigRulesHeadersTypePlaintext       SandboxResponseProxyConfigRulesHeadersType = "plaintext"
	SandboxResponseProxyConfigRulesHeadersTypeOpaque          SandboxResponseProxyConfigRulesHeadersType = "opaque"
	SandboxResponseProxyConfigRulesHeadersTypeWorkspaceSecret SandboxResponseProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesHeadersTypePlaintext, SandboxResponseProxyConfigRulesHeadersTypeOpaque, SandboxResponseProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxStatusResponse struct {
	Status        string                    `json:"status"`
	StatusMessage string                    `json:"status_message"`
	JSON          sandboxStatusResponseJSON `json:"-"`
}

// sandboxStatusResponseJSON contains the JSON metadata for the struct
// [SandboxStatusResponse]
type sandboxStatusResponseJSON struct {
	Status        apijson.Field
	StatusMessage apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxStatusResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxStatusResponseJSON) RawJSON() string {
	return r.raw
}

type ServiceURLResponse struct {
	Token      string                 `json:"token"`
	BrowserURL string                 `json:"browser_url"`
	ExpiresAt  string                 `json:"expires_at"`
	ServiceURL string                 `json:"service_url"`
	JSON       serviceURLResponseJSON `json:"-"`
}

// serviceURLResponseJSON contains the JSON metadata for the struct
// [ServiceURLResponse]
type serviceURLResponseJSON struct {
	Token       apijson.Field
	BrowserURL  apijson.Field
	ExpiresAt   apijson.Field
	ServiceURL  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ServiceURLResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r serviceURLResponseJSON) RawJSON() string {
	return r.raw
}

type SnapshotListResponse struct {
	Offset    int64                    `json:"offset"`
	Snapshots []SnapshotResponse       `json:"snapshots"`
	JSON      snapshotListResponseJSON `json:"-"`
}

// snapshotListResponseJSON contains the JSON metadata for the struct
// [SnapshotListResponse]
type snapshotListResponseJSON struct {
	Offset      apijson.Field
	Snapshots   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SnapshotListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r snapshotListResponseJSON) RawJSON() string {
	return r.raw
}

type SnapshotResponse struct {
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
	MemorySnapshotSizeBytes int64                `json:"memory_snapshot_size_bytes"`
	Name                    string               `json:"name"`
	RegistryID              string               `json:"registry_id"`
	SourceSandboxID         string               `json:"source_sandbox_id"`
	Status                  string               `json:"status"`
	StatusMessage           string               `json:"status_message"`
	UpdatedAt               string               `json:"updated_at"`
	JSON                    snapshotResponseJSON `json:"-"`
}

// snapshotResponseJSON contains the JSON metadata for the struct
// [SnapshotResponse]
type snapshotResponseJSON struct {
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

func (r *SnapshotResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r snapshotResponseJSON) RawJSON() string {
	return r.raw
}
