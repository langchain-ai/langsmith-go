// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"net/url"
	"reflect"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/pagination"
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
	Browsers   *SandboxBrowserService
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
	r.Browsers = NewSandboxBrowserService(opts...)
	return
}

// Returns priced usage per sandbox or snapshot and UTC hour in the half-open
// requested interval. LCU uses the recorded compute amount for sandboxes;
// snapshots have zero LCU. LSU allocates the recorded workspace storage amount
// proportionally to attributed bytes, including checkpoints on their sandbox and
// snapshots as separate resources. Resource filters preserve each resource's
// share. Rate changes do not reprice recorded amounts. An access-filtered page can
// have no items and a non-null next_cursor; continue until next_cursor is null.
func (r *SandboxService) ListUsageCosts(ctx context.Context, query SandboxListUsageCostsParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[SandboxListUsageCostsResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v2/sandboxes/usage/costs"
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

// Returns priced usage per sandbox or snapshot and UTC hour in the half-open
// requested interval. LCU uses the recorded compute amount for sandboxes;
// snapshots have zero LCU. LSU allocates the recorded workspace storage amount
// proportionally to attributed bytes, including checkpoints on their sandbox and
// snapshots as separate resources. Resource filters preserve each resource's
// share. Rate changes do not reprice recorded amounts. An access-filtered page can
// have no items and a non-null next_cursor; continue until next_cursor is null.
func (r *SandboxService) ListUsageCostsAutoPaging(ctx context.Context, query SandboxListUsageCostsParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[SandboxListUsageCostsResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.ListUsageCosts(ctx, query, opts...))
}

type DownloadURLResponse struct {
	Token       string `json:"token" api:"required"`
	DownloadURL string `json:"download_url" api:"required"`
	// ExpiresAt is null for a link that never expires.
	ExpiresAt string                  `json:"expires_at" api:"nullable"`
	JSON      downloadURLResponseJSON `json:"-"`
}

// downloadURLResponseJSON contains the JSON metadata for the struct
// [DownloadURLResponse]
type downloadURLResponseJSON struct {
	Token       apijson.Field
	DownloadURL apijson.Field
	ExpiresAt   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DownloadURLResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r downloadURLResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxListResponse struct {
	// This page of sandboxes.
	Items []SandboxResponse `json:"items"`
	// Cursor for the next page, or null on the last page. A non-null value is the only
	// signal that more pages exist. Treat it as opaque.
	NextCursor string `json:"next_cursor"`
	// Deprecated: use next_cursor. Offset to request for the next page, or 0 when no
	// pages remain.
	Offset int64 `json:"offset"`
	// Deprecated: use items. Duplicates items.
	Sandboxes []SandboxResponse       `json:"sandboxes"`
	JSON      sandboxListResponseJSON `json:"-"`
}

// sandboxListResponseJSON contains the JSON metadata for the struct
// [SandboxListResponse]
type sandboxListResponseJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
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
	ID string `json:"id"`
	// AccessDelegation is the LangSmith access this sandbox was granted, absent when
	// it has none. It is reported exactly as it was requested at create.
	AccessDelegation       SandboxResponseAccessDelegation `json:"access_delegation"`
	CPUMillicores          int64                           `json:"cpu_millicores"`
	CreatedAt              string                          `json:"created_at"`
	CreatedBy              string                          `json:"created_by"`
	DataplaneURL           string                          `json:"dataplane_url"`
	DeleteAfterStopSeconds int64                           `json:"delete_after_stop_seconds"`
	FsCapacityBytes        int64                           `json:"fs_capacity_bytes"`
	IdleTtlSeconds         int64                           `json:"idle_ttl_seconds"`
	Labels                 map[string]string               `json:"labels"`
	MemBytes               int64                           `json:"mem_bytes"`
	MountConfig            SandboxResponseMountConfig      `json:"mount_config"`
	Name                   string                          `json:"name"`
	PreserveMemoryOnStop   bool                            `json:"preserve_memory_on_stop"`
	ProxyConfig            SandboxResponseProxyConfig      `json:"proxy_config"`
	// RunConfig is what the sandbox's commands run with: the user, working directory
	// and base env beneath env_vars.
	RunConfig     SandboxResponseRunConfig `json:"run_config"`
	SizeClass     string                   `json:"size_class"`
	SnapshotID    string                   `json:"snapshot_id"`
	Status        string                   `json:"status"`
	StatusMessage string                   `json:"status_message"`
	StoppedAt     string                   `json:"stopped_at"`
	UpdatedAt     string                   `json:"updated_at"`
	UpdatedBy     string                   `json:"updated_by"`
	Vcpus         int64                    `json:"vcpus"`
	JSON          sandboxResponseJSON      `json:"-"`
}

// sandboxResponseJSON contains the JSON metadata for the struct [SandboxResponse]
type sandboxResponseJSON struct {
	ID                     apijson.Field
	AccessDelegation       apijson.Field
	CPUMillicores          apijson.Field
	CreatedAt              apijson.Field
	CreatedBy              apijson.Field
	DataplaneURL           apijson.Field
	DeleteAfterStopSeconds apijson.Field
	FsCapacityBytes        apijson.Field
	IdleTtlSeconds         apijson.Field
	Labels                 apijson.Field
	MemBytes               apijson.Field
	MountConfig            apijson.Field
	Name                   apijson.Field
	PreserveMemoryOnStop   apijson.Field
	ProxyConfig            apijson.Field
	RunConfig              apijson.Field
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

// AccessDelegation is the LangSmith access this sandbox was granted, absent when
// it has none. It is reported exactly as it was requested at create.
type SandboxResponseAccessDelegation struct {
	Mode        SandboxResponseAccessDelegationMode `json:"mode" api:"required"`
	Permissions []string                            `json:"permissions"`
	JSON        sandboxResponseAccessDelegationJSON `json:"-"`
}

// sandboxResponseAccessDelegationJSON contains the JSON metadata for the struct
// [SandboxResponseAccessDelegation]
type sandboxResponseAccessDelegationJSON struct {
	Mode        apijson.Field
	Permissions apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseAccessDelegation) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseAccessDelegationJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseAccessDelegationMode string

const (
	SandboxResponseAccessDelegationModeInherit  SandboxResponseAccessDelegationMode = "INHERIT"
	SandboxResponseAccessDelegationModeExplicit SandboxResponseAccessDelegationMode = "EXPLICIT"
)

func (r SandboxResponseAccessDelegationMode) IsKnown() bool {
	switch r {
	case SandboxResponseAccessDelegationModeInherit, SandboxResponseAccessDelegationModeExplicit:
		return true
	}
	return false
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
	Aws   SandboxResponseMountConfigAuthAws   `json:"aws"`
	Azure SandboxResponseMountConfigAuthAzure `json:"azure"`
	Gcp   SandboxResponseMountConfigAuthGcp   `json:"gcp"`
	JSON  sandboxResponseMountConfigAuthJSON  `json:"-"`
}

// sandboxResponseMountConfigAuthJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuth]
type sandboxResponseMountConfigAuthJSON struct {
	Aws         apijson.Field
	Azure       apijson.Field
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
	// This field can have the runtime type of
	// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID].
	AccessKeyID interface{} `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn string `json:"role_arn"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey].
	SecretAccessKey interface{} `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                                `json:"subject_scope"`
	JSON         sandboxResponseMountConfigAuthAwsJSON `json:"-"`
	union        SandboxResponseMountConfigAuthAwsUnion
}

// sandboxResponseMountConfigAuthAwsJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAws]
type sandboxResponseMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	RoleArn         apijson.Field
	SecretAccessKey apijson.Field
	SubjectScope    apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r sandboxResponseMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseMountConfigAuthAws{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseMountConfigAuthAwsUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig],
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig].
func (r SandboxResponseMountConfigAuthAws) AsUnion() SandboxResponseMountConfigAuthAwsUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig] or
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig].
type SandboxResponseMountConfigAuthAwsUnion interface {
	implementsSandboxResponseMountConfigAuthAws()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseMountConfigAuthAwsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig{}),
		},
	)
}

type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn string `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                                                                      `json:"subject_scope"`
	JSON         sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig]
type sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON struct {
	RoleArn      apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig) implementsSandboxResponseMountConfigAuthAws() {
}

type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig struct {
	AccessKeyID     SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn `json:"role_arn"`
	JSON    sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON    `json:"-"`
}

// sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig]
type sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	RoleArn         apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig) implementsSandboxResponseMountConfigAuthAws() {
}

type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID struct {
	Type  SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                                                     `json:"is_set"`
	Value string                                                                                   `json:"value"`
	JSON  sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID]
type sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType string

const (
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypePlaintext       SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "plaintext"
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeOpaque          SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "opaque"
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeWorkspaceSecret SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypePlaintext, SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeOpaque, SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey struct {
	Type  SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                                                         `json:"is_set"`
	Value string                                                                                       `json:"value"`
	JSON  sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey]
type sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType string

const (
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypePlaintext       SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "plaintext"
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeOpaque          SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "opaque"
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeWorkspaceSecret SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypePlaintext, SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeOpaque, SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace; the role's trust policy must trust
// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
// be changed after creation.
type SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn string

const (
	SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArnEmpty SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn = ""
)

func (r SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArnEmpty:
		return true
	}
	return false
}

type SandboxResponseMountConfigAuthAzure struct {
	// This field can have the runtime type of
	// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey].
	AccountKey interface{} `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID string `json:"client_id"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken].
	SasToken interface{} `json:"sas_token"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope string `json:"subject_scope"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for storage access tokens of client_id. Cannot be changed after
	// creation.
	TenantID string                                  `json:"tenant_id"`
	JSON     sandboxResponseMountConfigAuthAzureJSON `json:"-"`
	union    SandboxResponseMountConfigAuthAzureUnion
}

// sandboxResponseMountConfigAuthAzureJSON contains the JSON metadata for the
// struct [SandboxResponseMountConfigAuthAzure]
type sandboxResponseMountConfigAuthAzureJSON struct {
	AccountKey   apijson.Field
	ClientID     apijson.Field
	SasToken     apijson.Field
	SubjectScope apijson.Field
	TenantID     apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r sandboxResponseMountConfigAuthAzureJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseMountConfigAuthAzure) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseMountConfigAuthAzure{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseMountConfigAuthAzureUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig],
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig],
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig].
func (r SandboxResponseMountConfigAuthAzure) AsUnion() SandboxResponseMountConfigAuthAzureUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig],
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig]
// or [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig].
type SandboxResponseMountConfigAuthAzureUnion interface {
	implementsSandboxResponseMountConfigAuthAzure()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseMountConfigAuthAzureUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig{}),
		},
	)
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig struct {
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID string `json:"client_id" api:"required"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for storage access tokens of client_id. Cannot be changed after
	// creation.
	TenantID string `json:"tenant_id" api:"required"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope string                                                                                      `json:"subject_scope"`
	JSON         sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig]
type sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON struct {
	ClientID     apijson.Field
	TenantID     apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig) implementsSandboxResponseMountConfigAuthAzure() {
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig struct {
	AccountKey SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey `json:"account_key" api:"required"`
	JSON       sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON       `json:"-"`
}

// sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig]
type sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON struct {
	AccountKey  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig) implementsSandboxResponseMountConfigAuthAzure() {
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey struct {
	Type  SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType `json:"type" api:"required"`
	IsSet bool                                                                                            `json:"is_set"`
	Value string                                                                                          `json:"value"`
	JSON  sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey]
type sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType string

const (
	SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypePlaintext       SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "plaintext"
	SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeOpaque          SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "opaque"
	SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeWorkspaceSecret SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypePlaintext, SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeOpaque, SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig struct {
	SasToken SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken `json:"sas_token" api:"required"`
	JSON     sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON     `json:"-"`
}

// sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig]
type sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON struct {
	SasToken    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig) implementsSandboxResponseMountConfigAuthAzure() {
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken struct {
	Type  SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType `json:"type" api:"required"`
	IsSet bool                                                                                   `json:"is_set"`
	Value string                                                                                 `json:"value"`
	JSON  sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON `json:"-"`
}

// sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken]
type sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType string

const (
	SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypePlaintext       SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "plaintext"
	SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeOpaque          SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "opaque"
	SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeWorkspaceSecret SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypePlaintext, SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeOpaque, SandboxResponseMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigAuthGcp struct {
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly. Only valid with
	// workload_identity_provider.
	ServiceAccountEmail string `json:"service_account_email"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson].
	ServiceAccountJson interface{} `json:"service_account_json"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope string `json:"subject_scope"`
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Mutually exclusive with service_account_json. Cannot be changed after creation.
	WorkloadIdentityProvider string                                `json:"workload_identity_provider"`
	JSON                     sandboxResponseMountConfigAuthGcpJSON `json:"-"`
	union                    SandboxResponseMountConfigAuthGcpUnion
}

// sandboxResponseMountConfigAuthGcpJSON contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthGcp]
type sandboxResponseMountConfigAuthGcpJSON struct {
	ServiceAccountEmail      apijson.Field
	ServiceAccountJson       apijson.Field
	SubjectScope             apijson.Field
	WorkloadIdentityProvider apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r sandboxResponseMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseMountConfigAuthGcp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseMountConfigAuthGcpUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig],
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig].
func (r SandboxResponseMountConfigAuthGcp) AsUnion() SandboxResponseMountConfigAuthGcpUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig]
// or
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig].
type SandboxResponseMountConfigAuthGcpUnion interface {
	implementsSandboxResponseMountConfigAuthGcp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseMountConfigAuthGcpUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig{}),
		},
	)
}

type SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig struct {
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Mutually exclusive with service_account_json. Cannot be changed after creation.
	WorkloadIdentityProvider string `json:"workload_identity_provider" api:"required"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly. Only valid with
	// workload_identity_provider.
	ServiceAccountEmail string `json:"service_account_email"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope string                                                                                  `json:"subject_scope"`
	JSON         sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON `json:"-"`
}

// sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig]
type sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON struct {
	WorkloadIdentityProvider apijson.Field
	ServiceAccountEmail      apijson.Field
	SubjectScope             apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig) implementsSandboxResponseMountConfigAuthGcp() {
}

type SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig struct {
	ServiceAccountJson SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON               `json:"-"`
}

// sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig]
type sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig) implementsSandboxResponseMountConfigAuthGcp() {
}

type SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson struct {
	Type  SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                                                                    `json:"is_set"`
	Value string                                                                                                  `json:"value"`
	JSON  sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON `json:"-"`
}

// sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson]
type sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType string

const (
	SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypePlaintext       SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "plaintext"
	SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeOpaque          SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "opaque"
	SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeWorkspaceSecret SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "workspace_secret"
)

func (r SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypePlaintext, SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeOpaque, SandboxResponseMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseMountConfigMount struct {
	ID        string                               `json:"id" api:"required"`
	MountPath string                               `json:"mount_path" api:"required"`
	Type      SandboxResponseMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob].
	AzureBlob interface{} `json:"azure_blob"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub].
	Contexthub interface{} `json:"contexthub"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryption],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryption],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryption],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption].
	Encryption interface{} `json:"encryption"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGcs],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecGit],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecS3],
	// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3].
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
	AzureBlob   apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Encryption  apijson.Field
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
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec].
func (r SandboxResponseMountConfigMount) AsUnion() SandboxResponseMountConfigMountsUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpec],
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpec] or
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec].
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
			DiscriminatorValue: "azure_blob",
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
			DiscriminatorValue: "azure_blob",
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
			DiscriminatorValue: "azure_blob",
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
			DiscriminatorValue: "azure_blob",
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
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "azure_blob",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
	)
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID         string                                                                `json:"id" api:"required"`
	MountPath  string                                                                `json:"mount_path" api:"required"`
	S3         SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecS3         `json:"s3" api:"required"`
	Type       SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType       `json:"type" api:"required"`
	AzureBlob  SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecContexthub `json:"contexthub"`
	Encryption SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryption `json:"encryption"`
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
	AzureBlob   apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Encryption  apijson.Field
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
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeAzureBlob  SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "azure_blob"
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeAzureBlob, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                   `json:"storage_account" api:"required"`
	Prefix         string                                                                   `json:"prefix"`
	JSON           sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
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

type SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                    `json:"kms_key" api:"required"`
	JSON   sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryption]
type sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON) RawJSON() string {
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
	AzureBlob  SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecContexthub `json:"contexthub"`
	Encryption SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryption `json:"encryption"`
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
	AzureBlob   apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Encryption  apijson.Field
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
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeAzureBlob  SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "azure_blob"
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeAzureBlob, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                    `json:"storage_account" api:"required"`
	Prefix         string                                                                    `json:"prefix"`
	JSON           sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
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

type SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                     `json:"kms_key" api:"required"`
	JSON   sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryption]
type sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON) RawJSON() string {
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
	AzureBlob  SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecContexthub `json:"contexthub"`
	Encryption SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryption `json:"encryption"`
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
	AzureBlob   apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Encryption  apijson.Field
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
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeAzureBlob  SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "azure_blob"
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeAzureBlob, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                  `json:"storage_account" api:"required"`
	Prefix         string                                                                  `json:"prefix"`
	JSON           sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON contains
// the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
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

type SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                   `json:"kms_key" api:"required"`
	JSON   sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryption]
type sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON) RawJSON() string {
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
	AzureBlob  SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecCache      `json:"cache"`
	Encryption SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption `json:"encryption"`
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
	AzureBlob   apijson.Field
	Cache       apijson.Field
	Encryption  apijson.Field
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
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeAzureBlob  SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "azure_blob"
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeAzureBlob, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                         `json:"storage_account" api:"required"`
	Prefix         string                                                                         `json:"prefix"`
	JSON           sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
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

type SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                          `json:"kms_key" api:"required"`
	JSON   sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption]
type sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON) RawJSON() string {
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

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec struct {
	ID         string                                                                          `json:"id" api:"required"`
	AzureBlob  SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob  `json:"azure_blob" api:"required"`
	MountPath  string                                                                          `json:"mount_path" api:"required"`
	Type       SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType       `json:"type" api:"required"`
	Cache      SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache      `json:"cache"`
	Contexthub SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub `json:"contexthub"`
	Encryption SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption `json:"encryption"`
	Gcs        SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs        `json:"gcs"`
	Git        SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit        `json:"git"`
	ReadOnly   bool                                                                            `json:"read_only"`
	S3         SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3         `json:"s3"`
	JSON       sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON       `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON struct {
	ID          apijson.Field
	AzureBlob   apijson.Field
	MountPath   apijson.Field
	Type        apijson.Field
	Cache       apijson.Field
	Contexthub  apijson.Field
	Encryption  apijson.Field
	Gcs         apijson.Field
	Git         apijson.Field
	ReadOnly    apijson.Field
	S3          apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpec) implementsSandboxResponseMountConfigMount() {
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                             `json:"storage_account" api:"required"`
	Prefix         string                                                                             `json:"prefix"`
	JSON           sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType string

const (
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeS3         SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "s3"
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGcs        SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "gcs"
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeAzureBlob  SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "azure_blob"
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGit        SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "git"
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeContexthub SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "contexthub"
)

func (r SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeS3, SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGcs, SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeAzureBlob, SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGit, SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache struct {
	MaxSizeBytes     int64                                                                          `json:"max_size_bytes"`
	WritebackSeconds int64                                                                          `json:"writeback_seconds"`
	JSON             sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                                `json:"initial_pull_only"`
	JSON            sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                              `json:"kms_key" api:"required"`
	JSON   sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs struct {
	Bucket string                                                                       `json:"bucket" api:"required"`
	Prefix string                                                                       `json:"prefix"`
	JSON   sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit struct {
	RemoteURL              string                                                                       `json:"remote_url" api:"required"`
	Ref                    SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                        `json:"refresh_interval_seconds"`
	JSON                   sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef struct {
	Name string                                                                          `json:"name" api:"required"`
	Type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType string

const (
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeBranch SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType = "branch"
	SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeTag    SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType = "tag"
)

func (r SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeBranch, SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3 struct {
	Bucket      string                                                                      `json:"bucket" api:"required"`
	Region      string                                                                      `json:"region" api:"required"`
	EndpointURL string                                                                      `json:"endpoint_url"`
	PathStyle   bool                                                                        `json:"path_style"`
	Prefix      string                                                                      `json:"prefix"`
	JSON        sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON `json:"-"`
}

// sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON
// contains the JSON metadata for the struct
// [SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3]
type sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxResponseMountConfigMountsType string

const (
	SandboxResponseMountConfigMountsTypeS3         SandboxResponseMountConfigMountsType = "s3"
	SandboxResponseMountConfigMountsTypeGcs        SandboxResponseMountConfigMountsType = "gcs"
	SandboxResponseMountConfigMountsTypeAzureBlob  SandboxResponseMountConfigMountsType = "azure_blob"
	SandboxResponseMountConfigMountsTypeGit        SandboxResponseMountConfigMountsType = "git"
	SandboxResponseMountConfigMountsTypeContexthub SandboxResponseMountConfigMountsType = "contexthub"
)

func (r SandboxResponseMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxResponseMountConfigMountsTypeS3, SandboxResponseMountConfigMountsTypeGcs, SandboxResponseMountConfigMountsTypeAzureBlob, SandboxResponseMountConfigMountsTypeGit, SandboxResponseMountConfigMountsTypeContexthub:
		return true
	}
	return false
}

type SandboxResponseProxyConfig struct {
	AccessControl SandboxResponseProxyConfigAccessControl `json:"access_control"`
	Callbacks     []SandboxResponseProxyConfigCallback    `json:"callbacks"`
	// Description says what this configuration as a whole lets the sandbox reach,
	// complementing the per-rule descriptions. At most 1024 characters.
	Description string                           `json:"description"`
	NoProxy     []string                         `json:"no_proxy"`
	Rules       []SandboxResponseProxyConfigRule `json:"rules"`
	JSON        sandboxResponseProxyConfigJSON   `json:"-"`
}

// sandboxResponseProxyConfigJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfig]
type sandboxResponseProxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	Description   apijson.Field
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
	Name  string                               `json:"name" api:"required"`
	Aws   SandboxResponseProxyConfigRulesAws   `json:"aws"`
	Azure SandboxResponseProxyConfigRulesAzure `json:"azure"`
	// Description says what this rule lets the sandbox reach, so an agent driving the
	// sandbox can be told its capabilities. At most 1024 characters.
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
	// EnvVars are plaintext env vars set for every command in the sandbox while this
	// rule is enabled. Use them for tools that refuse to run unless a credential env
	// var is present (e.g. gh needs GH_TOKEN) even though this rule injects the real
	// credential on the wire — set a dummy value here so the command starts. Explicit
	// per-sandbox env_vars win over these, and provider-managed (AWS/GCP) vars win
	// over both.
	EnvVars map[string]string                       `json:"env_vars"`
	Gcp     SandboxResponseProxyConfigRulesGcp      `json:"gcp"`
	Headers []SandboxResponseProxyConfigRulesHeader `json:"headers"`
	// MatchHeaders restricts a header injection rule to requests carrying every listed
	// header, each written "name: value" with a lowercase name and an exact value.
	// Pair with headers of the same name to swap a placeholder the sandbox sends (e.g.
	// "authorization: Bearer account-b") for a real credential, so one host can serve
	// several accounts. Rules are evaluated in order and the first match wins.
	MatchHeaders []string `json:"match_headers"`
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
	Name         apijson.Field
	Aws          apijson.Field
	Azure        apijson.Field
	Description  apijson.Field
	Enabled      apijson.Field
	EnvVars      apijson.Field
	Gcp          apijson.Field
	Headers      apijson.Field
	MatchHeaders apijson.Field
	MatchHosts   apijson.Field
	MatchPaths   apijson.Field
	Type         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAws struct {
	// This field can have the runtime type of
	// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID].
	AccessKeyID interface{} `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn string `json:"role_arn"`
	// This field can have the runtime type of
	// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey].
	SecretAccessKey interface{} `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                                 `json:"subject_scope"`
	JSON         sandboxResponseProxyConfigRulesAwsJSON `json:"-"`
	union        SandboxResponseProxyConfigRulesAwsUnion
}

// sandboxResponseProxyConfigRulesAwsJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAws]
type sandboxResponseProxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	RoleArn         apijson.Field
	SecretAccessKey apijson.Field
	SubjectScope    apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r sandboxResponseProxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseProxyConfigRulesAws{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseProxyConfigRulesAwsUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig],
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig].
func (r SandboxResponseProxyConfigRulesAws) AsUnion() SandboxResponseProxyConfigRulesAwsUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig] or
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig].
type SandboxResponseProxyConfigRulesAwsUnion interface {
	implementsSandboxResponseProxyConfigRulesAws()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseProxyConfigRulesAwsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig{}),
		},
	)
}

type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn string `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                                                            `json:"subject_scope"`
	JSON         sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON contains the
// JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig]
type sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON struct {
	RoleArn      apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) implementsSandboxResponseProxyConfigRulesAws() {
}

type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig struct {
	AccessKeyID     SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn `json:"role_arn"`
	JSON    sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON    `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON contains the
// JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig]
type sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	RoleArn         apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) implementsSandboxResponseProxyConfigRulesAws() {
}

type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID struct {
	Type  SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                                           `json:"is_set"`
	Value string                                                                         `json:"value"`
	JSON  sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID]
type sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType string

const (
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext       SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "plaintext"
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque          SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "opaque"
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext, SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque, SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey struct {
	Type  SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                                               `json:"is_set"`
	Value string                                                                             `json:"value"`
	JSON  sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey]
type sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType string

const (
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext       SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "plaintext"
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque          SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "opaque"
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext, SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque, SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace. The role's trust policy must trust
// LangSmith's issuer for the token subject. Credentials renew automatically, and
// access follows the role's effective AWS permissions. Mutually exclusive with
// static keys, and cannot be changed after creation.
type SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn string

const (
	SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn = ""
)

func (r SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesAzure struct {
	// This field can have the runtime type of
	// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey].
	AccountKey interface{} `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID string `json:"client_id"`
	// This field can have the runtime type of
	// [[]SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource].
	Resources interface{} `json:"resources"`
	// This field can have the runtime type of
	// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken].
	SasToken interface{} `json:"sas_token"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount string `json:"storage_account"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope string `json:"subject_scope"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for access tokens of client_id, so no secret is stored. Cannot be
	// changed after creation.
	TenantID string                                   `json:"tenant_id"`
	JSON     sandboxResponseProxyConfigRulesAzureJSON `json:"-"`
	union    SandboxResponseProxyConfigRulesAzureUnion
}

// sandboxResponseProxyConfigRulesAzureJSON contains the JSON metadata for the
// struct [SandboxResponseProxyConfigRulesAzure]
type sandboxResponseProxyConfigRulesAzureJSON struct {
	AccountKey     apijson.Field
	ClientID       apijson.Field
	Resources      apijson.Field
	SasToken       apijson.Field
	StorageAccount apijson.Field
	SubjectScope   apijson.Field
	TenantID       apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r sandboxResponseProxyConfigRulesAzureJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseProxyConfigRulesAzure) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseProxyConfigRulesAzure{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseProxyConfigRulesAzureUnion] interface which
// you can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig],
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig],
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig].
func (r SandboxResponseProxyConfigRulesAzure) AsUnion() SandboxResponseProxyConfigRulesAzureUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig],
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig] or
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig].
type SandboxResponseProxyConfigRulesAzureUnion interface {
	implementsSandboxResponseProxyConfigRulesAzure()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseProxyConfigRulesAzureUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig{}),
		},
	)
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig struct {
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID string `json:"client_id" api:"required"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for access tokens of client_id, so no secret is stored. Cannot be
	// changed after creation.
	TenantID string `json:"tenant_id" api:"required"`
	// Resources narrows the Entra resources access tokens are obtained for, from
	// https://storage.azure.com, https://vault.azure.net, https://management.azure.com
	// and https://cognitiveservices.azure.com. Defaults to all of them. Only valid
	// with tenant_id.
	Resources []SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource `json:"resources"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope string                                                                            `json:"subject_scope"`
	JSON         sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig]
type sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON struct {
	ClientID     apijson.Field
	TenantID     apijson.Field
	Resources    apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) implementsSandboxResponseProxyConfigRulesAzure() {
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource string

const (
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom           SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://storage.azure.com"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet             SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://vault.azure.net"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom        SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://management.azure.com"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://cognitiveservices.azure.com"
)

func (r SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig struct {
	AccountKey SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey `json:"account_key" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount string                                                                      `json:"storage_account" api:"required"`
	JSON           sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig]
type sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON struct {
	AccountKey     apijson.Field
	StorageAccount apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) implementsSandboxResponseProxyConfigRulesAzure() {
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey struct {
	Type  SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType `json:"type" api:"required"`
	IsSet bool                                                                                  `json:"is_set"`
	Value string                                                                                `json:"value"`
	JSON  sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey]
type sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType string

const (
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext       SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "plaintext"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque          SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "opaque"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig struct {
	SasToken SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken `json:"sas_token" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount string                                                               `json:"storage_account" api:"required"`
	JSON           sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON contains
// the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig]
type sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON struct {
	SasToken       apijson.Field
	StorageAccount apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfig) implementsSandboxResponseProxyConfigRulesAzure() {
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken struct {
	Type  SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType `json:"type" api:"required"`
	IsSet bool                                                                         `json:"is_set"`
	Value string                                                                       `json:"value"`
	JSON  sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken]
type sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType string

const (
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext       SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "plaintext"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque          SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "opaque"
	SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque, SandboxResponseProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxResponseProxyConfigRulesGcp struct {
	// This field can have the runtime type of [[]string].
	Scopes interface{} `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail string `json:"service_account_email"`
	// This field can have the runtime type of
	// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson].
	ServiceAccountJson interface{} `json:"service_account_json"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope string `json:"subject_scope"`
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key.
	// Use the full resource name:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Cannot be changed after creation.
	WorkloadIdentityProvider string                                 `json:"workload_identity_provider"`
	JSON                     sandboxResponseProxyConfigRulesGcpJSON `json:"-"`
	union                    SandboxResponseProxyConfigRulesGcpUnion
}

// sandboxResponseProxyConfigRulesGcpJSON contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesGcp]
type sandboxResponseProxyConfigRulesGcpJSON struct {
	Scopes                   apijson.Field
	ServiceAccountEmail      apijson.Field
	ServiceAccountJson       apijson.Field
	SubjectScope             apijson.Field
	WorkloadIdentityProvider apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r sandboxResponseProxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxResponseProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxResponseProxyConfigRulesGcp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxResponseProxyConfigRulesGcpUnion] interface which you
// can cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig],
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig].
func (r SandboxResponseProxyConfigRulesGcp) AsUnion() SandboxResponseProxyConfigRulesGcpUnion {
	return r.union
}

// Union satisfied by
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig] or
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig].
type SandboxResponseProxyConfigRulesGcpUnion interface {
	implementsSandboxResponseProxyConfigRulesGcp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxResponseProxyConfigRulesGcpUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig{}),
		},
	)
}

type SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig struct {
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key.
	// Use the full resource name:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Cannot be changed after creation.
	WorkloadIdentityProvider string `json:"workload_identity_provider" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes []string `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail string `json:"service_account_email"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope string                                                                        `json:"subject_scope"`
	JSON         sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig]
type sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON struct {
	WorkloadIdentityProvider apijson.Field
	Scopes                   apijson.Field
	ServiceAccountEmail      apijson.Field
	SubjectScope             apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) implementsSandboxResponseProxyConfigRulesGcp() {
}

type SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig struct {
	ServiceAccountJson SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson `json:"service_account_json" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes []string                                                                    `json:"scopes"`
	JSON   sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig]
type sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON struct {
	ServiceAccountJson apijson.Field
	Scopes             apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) implementsSandboxResponseProxyConfigRulesGcp() {
}

type SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson struct {
	Type  SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                                                          `json:"is_set"`
	Value string                                                                                        `json:"value"`
	JSON  sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON `json:"-"`
}

// sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON
// contains the JSON metadata for the struct
// [SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson]
type sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType string

const (
	SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext       SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "plaintext"
	SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque          SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "opaque"
	SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "workspace_secret"
)

func (r SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext, SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque, SandboxResponseProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret:
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

// RunConfig is what the sandbox's commands run with: the user, working directory
// and base env beneath env_vars.
type SandboxResponseRunConfig struct {
	EnvVars map[string]string            `json:"env_vars"`
	User    string                       `json:"user"`
	WorkDir string                       `json:"work_dir"`
	JSON    sandboxResponseRunConfigJSON `json:"-"`
}

// sandboxResponseRunConfigJSON contains the JSON metadata for the struct
// [SandboxResponseRunConfig]
type sandboxResponseRunConfigJSON struct {
	EnvVars     apijson.Field
	User        apijson.Field
	WorkDir     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseRunConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseRunConfigJSON) RawJSON() string {
	return r.raw
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
	// Token and ExpiresAt are empty in LangSmith login mode (no token is minted).
	Token string `json:"token"`
	// Access echoes the enabled LangSmith login level ("restricted"/"workspace");
	// omitted in token mode.
	Access     ServiceURLResponseAccess `json:"access"`
	BrowserURL string                   `json:"browser_url"`
	ExpiresAt  string                   `json:"expires_at"`
	ServiceURL string                   `json:"service_url"`
	JSON       serviceURLResponseJSON   `json:"-"`
}

// serviceURLResponseJSON contains the JSON metadata for the struct
// [ServiceURLResponse]
type serviceURLResponseJSON struct {
	Token       apijson.Field
	Access      apijson.Field
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

// Access echoes the enabled LangSmith login level ("restricted"/"workspace");
// omitted in token mode.
type ServiceURLResponseAccess string

const (
	ServiceURLResponseAccessRestricted ServiceURLResponseAccess = "restricted"
	ServiceURLResponseAccessWorkspace  ServiceURLResponseAccess = "workspace"
)

func (r ServiceURLResponseAccess) IsKnown() bool {
	switch r {
	case ServiceURLResponseAccessRestricted, ServiceURLResponseAccessWorkspace:
		return true
	}
	return false
}

type SnapshotListResponse struct {
	// This page of snapshots.
	Items []SnapshotResponse `json:"items"`
	// Cursor for the next page, or null on the last page. A non-null value is the only
	// signal that more pages exist. Treat it as opaque.
	NextCursor string `json:"next_cursor"`
	// Deprecated: use next_cursor. Offset to request for the next page, or 0 when no
	// pages remain.
	Offset int64 `json:"offset"`
	// Deprecated: use items. Duplicates items.
	Snapshots []SnapshotResponse       `json:"snapshots"`
	JSON      snapshotListResponseJSON `json:"-"`
}

// snapshotListResponseJSON contains the JSON metadata for the struct
// [SnapshotListResponse]
type snapshotListResponseJSON struct {
	Items       apijson.Field
	NextCursor  apijson.Field
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
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	CreatedBy string `json:"created_by"`
	// Description says what this snapshot's image can do, so a caller can hand it to
	// an agent as a capability summary.
	Description     string            `json:"description"`
	DockerImage     string            `json:"docker_image"`
	FsCapacityBytes int64             `json:"fs_capacity_bytes"`
	FsUsedBytes     int64             `json:"fs_used_bytes"`
	ImageDigest     string            `json:"image_digest"`
	Labels          map[string]string `json:"labels"`
	// MemorySnapshotSizeBytes is non-nil iff the snapshot was captured with VM memory
	// state. A non-nil value is the canonical signal that this snapshot can
	// warm-restore from memory; nil means rootfs only.
	MemorySnapshotSizeBytes int64  `json:"memory_snapshot_size_bytes"`
	Name                    string `json:"name"`
	RegistryID              string `json:"registry_id"`
	// RunConfig is what sandboxes from this snapshot boot with. Absent on snapshots
	// built before it was recorded, which run as root with their own env.
	RunConfig       SnapshotResponseRunConfig `json:"run_config"`
	SourceSandboxID string                    `json:"source_sandbox_id"`
	Status          string                    `json:"status"`
	StatusMessage   string                    `json:"status_message"`
	// Tags currently resolving to this snapshot, under Name. A snapshot with no tags
	// is dangling — addressable only by id.
	Tags      []string             `json:"tags"`
	UpdatedAt string               `json:"updated_at"`
	JSON      snapshotResponseJSON `json:"-"`
}

// snapshotResponseJSON contains the JSON metadata for the struct
// [SnapshotResponse]
type snapshotResponseJSON struct {
	ID                      apijson.Field
	CreatedAt               apijson.Field
	CreatedBy               apijson.Field
	Description             apijson.Field
	DockerImage             apijson.Field
	FsCapacityBytes         apijson.Field
	FsUsedBytes             apijson.Field
	ImageDigest             apijson.Field
	Labels                  apijson.Field
	MemorySnapshotSizeBytes apijson.Field
	Name                    apijson.Field
	RegistryID              apijson.Field
	RunConfig               apijson.Field
	SourceSandboxID         apijson.Field
	Status                  apijson.Field
	StatusMessage           apijson.Field
	Tags                    apijson.Field
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

// RunConfig is what sandboxes from this snapshot boot with. Absent on snapshots
// built before it was recorded, which run as root with their own env.
type SnapshotResponseRunConfig struct {
	EnvVars map[string]string             `json:"env_vars"`
	User    string                        `json:"user"`
	WorkDir string                        `json:"work_dir"`
	JSON    snapshotResponseRunConfigJSON `json:"-"`
}

// snapshotResponseRunConfigJSON contains the JSON metadata for the struct
// [SnapshotResponseRunConfig]
type snapshotResponseRunConfigJSON struct {
	EnvVars     apijson.Field
	User        apijson.Field
	WorkDir     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SnapshotResponseRunConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r snapshotResponseRunConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxListUsageCostsResponse struct {
	// Recorded compute usage in LangSmith Compute Units (LCU), as a decimal string
	// with up to six fractional digits and trailing zeros omitted. Snapshots return
	// "0".
	Lcu string `json:"lcu" api:"required" format:"decimal"`
	// Allocated storage usage in LangSmith Storage Units (LSU), as a decimal string
	// with up to six fractional digits and trailing zeros omitted.
	Lsu          string                                    `json:"lsu" api:"required" format:"decimal"`
	PeriodStart  time.Time                                 `json:"period_start" api:"required" format:"date-time"`
	ResourceID   string                                    `json:"resource_id" api:"required" format:"uuid"`
	ResourceType SandboxListUsageCostsResponseResourceType `json:"resource_type" api:"required"`
	JSON         sandboxListUsageCostsResponseJSON         `json:"-"`
}

// sandboxListUsageCostsResponseJSON contains the JSON metadata for the struct
// [SandboxListUsageCostsResponse]
type sandboxListUsageCostsResponseJSON struct {
	Lcu          apijson.Field
	Lsu          apijson.Field
	PeriodStart  apijson.Field
	ResourceID   apijson.Field
	ResourceType apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxListUsageCostsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxListUsageCostsResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxListUsageCostsResponseResourceType string

const (
	SandboxListUsageCostsResponseResourceTypeSandbox  SandboxListUsageCostsResponseResourceType = "SANDBOX"
	SandboxListUsageCostsResponseResourceTypeSnapshot SandboxListUsageCostsResponseResourceType = "SNAPSHOT"
)

func (r SandboxListUsageCostsResponseResourceType) IsKnown() bool {
	switch r {
	case SandboxListUsageCostsResponseResourceTypeSandbox, SandboxListUsageCostsResponseResourceTypeSnapshot:
		return true
	}
	return false
}

type SandboxListUsageCostsParams struct {
	// Exclusive RFC3339 end time; the range must not exceed 31 days
	EndTime param.Field[time.Time] `query:"end_time" api:"required" format:"date-time"`
	// Inclusive RFC3339 start time
	StartTime param.Field[time.Time] `query:"start_time" api:"required" format:"date-time"`
	// Opaque pagination cursor
	Cursor param.Field[string] `query:"cursor"`
	// HOUR returns hourly buckets. RESOURCE sums each resource over the requested
	// interval and sets period_start to start_time.
	Granularity param.Field[SandboxListUsageCostsParamsGranularity] `query:"granularity"`
	// Include usage totals across the filtered interval (RESOURCE only). Omitted when
	// unrestricted resource access cannot be verified or the service does not support
	// totals.
	IncludeTotals param.Field[bool] `query:"include_totals"`
	// Maximum rows to return
	PageSize param.Field[int64] `query:"page_size"`
	// Resource UUID filter; repeat this parameter up to 100 times
	ResourceIDs param.Field[[]string] `query:"resource_ids"`
	// Resource type filter
	ResourceType param.Field[SandboxListUsageCostsParamsResourceType] `query:"resource_type"`
}

// URLQuery serializes [SandboxListUsageCostsParams]'s query parameters as
// `url.Values`.
func (r SandboxListUsageCostsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// HOUR returns hourly buckets. RESOURCE sums each resource over the requested
// interval and sets period_start to start_time.
type SandboxListUsageCostsParamsGranularity string

const (
	SandboxListUsageCostsParamsGranularityHour     SandboxListUsageCostsParamsGranularity = "HOUR"
	SandboxListUsageCostsParamsGranularityResource SandboxListUsageCostsParamsGranularity = "RESOURCE"
)

func (r SandboxListUsageCostsParamsGranularity) IsKnown() bool {
	switch r {
	case SandboxListUsageCostsParamsGranularityHour, SandboxListUsageCostsParamsGranularityResource:
		return true
	}
	return false
}

// Resource type filter
type SandboxListUsageCostsParamsResourceType string

const (
	SandboxListUsageCostsParamsResourceTypeSandbox  SandboxListUsageCostsParamsResourceType = "SANDBOX"
	SandboxListUsageCostsParamsResourceTypeSnapshot SandboxListUsageCostsParamsResourceType = "SNAPSHOT"
)

func (r SandboxListUsageCostsParamsResourceType) IsKnown() bool {
	switch r {
	case SandboxListUsageCostsParamsResourceTypeSandbox, SandboxListUsageCostsParamsResourceTypeSnapshot:
		return true
	}
	return false
}
