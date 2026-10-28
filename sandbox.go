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

type ProxyConfig struct {
	AccessControl ProxyConfigAccessControl `json:"access_control"`
	Callbacks     []ProxyConfigCallback    `json:"callbacks"`
	// Description says what this configuration as a whole lets the sandbox reach,
	// complementing the per-rule descriptions. At most 1024 characters.
	Description string            `json:"description"`
	NoProxy     []string          `json:"no_proxy"`
	Rules       []ProxyConfigRule `json:"rules"`
	JSON        proxyConfigJSON   `json:"-"`
}

// proxyConfigJSON contains the JSON metadata for the struct [ProxyConfig]
type proxyConfigJSON struct {
	AccessControl apijson.Field
	Callbacks     apijson.Field
	Description   apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *ProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigAccessControl struct {
	AllowList []string                     `json:"allow_list"`
	DenyList  []string                     `json:"deny_list"`
	JSON      proxyConfigAccessControlJSON `json:"-"`
}

// proxyConfigAccessControlJSON contains the JSON metadata for the struct
// [ProxyConfigAccessControl]
type proxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigCallback struct {
	MatchHosts     []string                            `json:"match_hosts" api:"required"`
	TtlSeconds     int64                               `json:"ttl_seconds" api:"required"`
	URL            string                              `json:"url" api:"required"`
	FullRequest    bool                                `json:"full_request"`
	RequestHeaders []ProxyConfigCallbacksRequestHeader `json:"request_headers"`
	JSON           proxyConfigCallbackJSON             `json:"-"`
}

// proxyConfigCallbackJSON contains the JSON metadata for the struct
// [ProxyConfigCallback]
type proxyConfigCallbackJSON struct {
	MatchHosts     apijson.Field
	TtlSeconds     apijson.Field
	URL            apijson.Field
	FullRequest    apijson.Field
	RequestHeaders apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *ProxyConfigCallback) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigCallbackJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigCallbacksRequestHeader struct {
	Name  string                                 `json:"name" api:"required"`
	Type  ProxyConfigCallbacksRequestHeadersType `json:"type" api:"required"`
	IsSet bool                                   `json:"is_set"`
	Value string                                 `json:"value"`
	JSON  proxyConfigCallbacksRequestHeaderJSON  `json:"-"`
}

// proxyConfigCallbacksRequestHeaderJSON contains the JSON metadata for the struct
// [ProxyConfigCallbacksRequestHeader]
type proxyConfigCallbacksRequestHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigCallbacksRequestHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigCallbacksRequestHeaderJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigCallbacksRequestHeadersType string

const (
	ProxyConfigCallbacksRequestHeadersTypePlaintext       ProxyConfigCallbacksRequestHeadersType = "plaintext"
	ProxyConfigCallbacksRequestHeadersTypeOpaque          ProxyConfigCallbacksRequestHeadersType = "opaque"
	ProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret ProxyConfigCallbacksRequestHeadersType = "workspace_secret"
)

func (r ProxyConfigCallbacksRequestHeadersType) IsKnown() bool {
	switch r {
	case ProxyConfigCallbacksRequestHeadersTypePlaintext, ProxyConfigCallbacksRequestHeadersTypeOpaque, ProxyConfigCallbacksRequestHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type ProxyConfigRule struct {
	Name  string                `json:"name" api:"required"`
	Aws   ProxyConfigRulesAws   `json:"aws"`
	Azure ProxyConfigRulesAzure `json:"azure"`
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
	EnvVars map[string]string        `json:"env_vars"`
	Gcp     ProxyConfigRulesGcp      `json:"gcp"`
	Headers []ProxyConfigRulesHeader `json:"headers"`
	// MatchHeaders restricts a header injection rule to requests carrying every listed
	// header, each written "name: value" with a lowercase name and an exact value.
	// Pair with headers of the same name to swap a placeholder the sandbox sends (e.g.
	// "authorization: Bearer account-b") for a real credential, so one host can serve
	// several accounts. Rules are evaluated in order and the first match wins.
	MatchHeaders []string `json:"match_headers"`
	// MatchHosts is only accepted for header injection and relay rules. Provider auth
	// rules use built-in host matching. A relay rule may also list IP addresses and
	// CIDR blocks (10.0.0.0/8), which match connections the sandbox makes to an IP
	// address rather than a name.
	MatchHosts []string `json:"match_hosts"`
	MatchPaths []string `json:"match_paths"`
	// RelayID is the ID of the LangSmith relay (a connection of kind relay) that
	// carries connections to match_hosts when type is relay. Those connections leave
	// from the relay's connector, and its network and IP, instead of the sandbox host.
	// Access control still applies first, the first enabled matching relay rule wins,
	// and a matched host never falls back to a direct connection.
	RelayID string              `json:"relay_id" format:"uuid"`
	Type    string              `json:"type"`
	JSON    proxyConfigRuleJSON `json:"-"`
}

// proxyConfigRuleJSON contains the JSON metadata for the struct [ProxyConfigRule]
type proxyConfigRuleJSON struct {
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
	RelayID      apijson.Field
	Type         apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesAws struct {
	// This field can have the runtime type of
	// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID].
	AccessKeyID interface{} `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn string `json:"role_arn"`
	// This field can have the runtime type of
	// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey].
	SecretAccessKey interface{} `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                  `json:"subject_scope"`
	JSON         proxyConfigRulesAwsJSON `json:"-"`
	union        ProxyConfigRulesAwsUnion
}

// proxyConfigRulesAwsJSON contains the JSON metadata for the struct
// [ProxyConfigRulesAws]
type proxyConfigRulesAwsJSON struct {
	AccessKeyID     apijson.Field
	RoleArn         apijson.Field
	SecretAccessKey apijson.Field
	SubjectScope    apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r proxyConfigRulesAwsJSON) RawJSON() string {
	return r.raw
}

func (r *ProxyConfigRulesAws) UnmarshalJSON(data []byte) (err error) {
	*r = ProxyConfigRulesAws{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ProxyConfigRulesAwsUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are
// [ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig],
// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig].
func (r ProxyConfigRulesAws) AsUnion() ProxyConfigRulesAwsUnion {
	return r.union
}

// Union satisfied by [ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig] or
// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig].
type ProxyConfigRulesAwsUnion interface {
	implementsProxyConfigRulesAws()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ProxyConfigRulesAwsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig{}),
		},
	)
}

type ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn string `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                                             `json:"subject_scope"`
	JSON         proxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON `json:"-"`
}

// proxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON contains the JSON metadata
// for the struct [ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig]
type proxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON struct {
	RoleArn      apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAwsSandboxesProxyAwsRoleConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsRoleConfig) implementsProxyConfigRulesAws() {}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig struct {
	AccessKeyID     ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn `json:"role_arn"`
	JSON    proxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON    `json:"-"`
}

// proxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON contains the JSON metadata
// for the struct [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig]
type proxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	RoleArn         apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAwsSandboxesProxyAwsStaticConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfig) implementsProxyConfigRulesAws() {}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID struct {
	Type  ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                            `json:"is_set"`
	Value string                                                          `json:"value"`
	JSON  proxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON `json:"-"`
}

// proxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON contains the
// JSON metadata for the struct
// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID]
type proxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType string

const (
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext       ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "plaintext"
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque          ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "opaque"
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType = "workspace_secret"
)

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypePlaintext, ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeOpaque, ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey struct {
	Type  ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                                `json:"is_set"`
	Value string                                                              `json:"value"`
	JSON  proxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON `json:"-"`
}

// proxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON contains the
// JSON metadata for the struct
// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey]
type proxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType string

const (
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext       ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "plaintext"
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque          ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "opaque"
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType = "workspace_secret"
)

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypePlaintext, ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeOpaque, ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace. The role's trust policy must trust
// LangSmith's issuer for the token subject. Credentials renew automatically, and
// access follows the role's effective AWS permissions. Mutually exclusive with
// static keys, and cannot be changed after creation.
type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn string

const (
	ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn = ""
)

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn) IsKnown() bool {
	switch r {
	case ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArnEmpty:
		return true
	}
	return false
}

type ProxyConfigRulesAzure struct {
	// This field can have the runtime type of
	// [ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey].
	AccountKey interface{} `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID string `json:"client_id"`
	// This field can have the runtime type of
	// [[]ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource].
	Resources interface{} `json:"resources"`
	// This field can have the runtime type of
	// [ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken].
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
	TenantID string                    `json:"tenant_id"`
	JSON     proxyConfigRulesAzureJSON `json:"-"`
	union    ProxyConfigRulesAzureUnion
}

// proxyConfigRulesAzureJSON contains the JSON metadata for the struct
// [ProxyConfigRulesAzure]
type proxyConfigRulesAzureJSON struct {
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

func (r proxyConfigRulesAzureJSON) RawJSON() string {
	return r.raw
}

func (r *ProxyConfigRulesAzure) UnmarshalJSON(data []byte) (err error) {
	*r = ProxyConfigRulesAzure{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ProxyConfigRulesAzureUnion] interface which you can cast to
// the specific types for more type safety.
//
// Possible runtime types of the union are
// [ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig],
// [ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig],
// [ProxyConfigRulesAzureSandboxesProxyAzureSasConfig].
func (r ProxyConfigRulesAzure) AsUnion() ProxyConfigRulesAzureUnion {
	return r.union
}

// Union satisfied by
// [ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig],
// [ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig] or
// [ProxyConfigRulesAzureSandboxesProxyAzureSasConfig].
type ProxyConfigRulesAzureUnion interface {
	implementsProxyConfigRulesAzure()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ProxyConfigRulesAzureUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesAzureSandboxesProxyAzureSasConfig{}),
		},
	)
}

type ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig struct {
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
	Resources []ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource `json:"resources"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope string                                                             `json:"subject_scope"`
	JSON         proxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON `json:"-"`
}

// proxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON contains the
// JSON metadata for the struct
// [ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig]
type proxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON struct {
	ClientID     apijson.Field
	TenantID     apijson.Field
	Resources    apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfig) implementsProxyConfigRulesAzure() {
}

type ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource string

const (
	ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom           ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://storage.azure.com"
	ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet             ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://vault.azure.net"
	ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom        ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://management.azure.com"
	ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource = "https://cognitiveservices.azure.com"
)

func (r ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource) IsKnown() bool {
	switch r {
	case ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom, ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSVaultAzureNet, ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSManagementAzureCom, ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSCognitiveservicesAzureCom:
		return true
	}
	return false
}

type ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig struct {
	AccountKey ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey `json:"account_key" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount string                                                       `json:"storage_account" api:"required"`
	JSON           proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON `json:"-"`
}

// proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON contains the JSON
// metadata for the struct
// [ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig]
type proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON struct {
	AccountKey     apijson.Field
	StorageAccount apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfig) implementsProxyConfigRulesAzure() {}

type ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey struct {
	Type  ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType `json:"type" api:"required"`
	IsSet bool                                                                   `json:"is_set"`
	Value string                                                                 `json:"value"`
	JSON  proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON `json:"-"`
}

// proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON contains
// the JSON metadata for the struct
// [ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey]
type proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType string

const (
	ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext       ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "plaintext"
	ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque          ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "opaque"
	ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType = "workspace_secret"
)

func (r ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType) IsKnown() bool {
	switch r {
	case ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypePlaintext, ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeOpaque, ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type ProxyConfigRulesAzureSandboxesProxyAzureSasConfig struct {
	SasToken ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken `json:"sas_token" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount string                                                `json:"storage_account" api:"required"`
	JSON           proxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON `json:"-"`
}

// proxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON contains the JSON metadata
// for the struct [ProxyConfigRulesAzureSandboxesProxyAzureSasConfig]
type proxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON struct {
	SasToken       apijson.Field
	StorageAccount apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *ProxyConfigRulesAzureSandboxesProxyAzureSasConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAzureSandboxesProxyAzureSasConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureSasConfig) implementsProxyConfigRulesAzure() {}

type ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken struct {
	Type  ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType `json:"type" api:"required"`
	IsSet bool                                                          `json:"is_set"`
	Value string                                                        `json:"value"`
	JSON  proxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON `json:"-"`
}

// proxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON contains the JSON
// metadata for the struct
// [ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken]
type proxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasToken) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType string

const (
	ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext       ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "plaintext"
	ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque          ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "opaque"
	ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType = "workspace_secret"
)

func (r ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType) IsKnown() bool {
	switch r {
	case ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypePlaintext, ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeOpaque, ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type ProxyConfigRulesGcp struct {
	// This field can have the runtime type of [[]string].
	Scopes interface{} `json:"scopes"`
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly.
	ServiceAccountEmail string `json:"service_account_email"`
	// This field can have the runtime type of
	// [ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson].
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
	WorkloadIdentityProvider string                  `json:"workload_identity_provider"`
	JSON                     proxyConfigRulesGcpJSON `json:"-"`
	union                    ProxyConfigRulesGcpUnion
}

// proxyConfigRulesGcpJSON contains the JSON metadata for the struct
// [ProxyConfigRulesGcp]
type proxyConfigRulesGcpJSON struct {
	Scopes                   apijson.Field
	ServiceAccountEmail      apijson.Field
	ServiceAccountJson       apijson.Field
	SubjectScope             apijson.Field
	WorkloadIdentityProvider apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r proxyConfigRulesGcpJSON) RawJSON() string {
	return r.raw
}

func (r *ProxyConfigRulesGcp) UnmarshalJSON(data []byte) (err error) {
	*r = ProxyConfigRulesGcp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [ProxyConfigRulesGcpUnion] interface which you can cast to the
// specific types for more type safety.
//
// Possible runtime types of the union are
// [ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig],
// [ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig].
func (r ProxyConfigRulesGcp) AsUnion() ProxyConfigRulesGcpUnion {
	return r.union
}

// Union satisfied by [ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig]
// or [ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig].
type ProxyConfigRulesGcpUnion interface {
	implementsProxyConfigRulesGcp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ProxyConfigRulesGcpUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig{}),
		},
	)
}

type ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig struct {
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
	SubjectScope string                                                         `json:"subject_scope"`
	JSON         proxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON `json:"-"`
}

// proxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON contains the JSON
// metadata for the struct
// [ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig]
type proxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON struct {
	WorkloadIdentityProvider apijson.Field
	Scopes                   apijson.Field
	ServiceAccountEmail      apijson.Field
	SubjectScope             apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfig) implementsProxyConfigRulesGcp() {}

type ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig struct {
	ServiceAccountJson ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson `json:"service_account_json" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes []string                                                     `json:"scopes"`
	JSON   proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON `json:"-"`
}

// proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON contains the JSON
// metadata for the struct
// [ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig]
type proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON struct {
	ServiceAccountJson apijson.Field
	Scopes             apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigJSON) RawJSON() string {
	return r.raw
}

func (r ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfig) implementsProxyConfigRulesGcp() {}

type ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson struct {
	Type  ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                                           `json:"is_set"`
	Value string                                                                         `json:"value"`
	JSON  proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON `json:"-"`
}

// proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON
// contains the JSON metadata for the struct
// [ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson]
type proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType string

const (
	ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext       ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "plaintext"
	ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque          ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "opaque"
	ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType = "workspace_secret"
)

func (r ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypePlaintext, ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeOpaque, ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type ProxyConfigRulesHeader struct {
	Name  string                      `json:"name" api:"required"`
	Type  ProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                        `json:"is_set"`
	Value string                      `json:"value"`
	JSON  proxyConfigRulesHeaderJSON  `json:"-"`
}

// proxyConfigRulesHeaderJSON contains the JSON metadata for the struct
// [ProxyConfigRulesHeader]
type proxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r proxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type ProxyConfigRulesHeadersType string

const (
	ProxyConfigRulesHeadersTypePlaintext       ProxyConfigRulesHeadersType = "plaintext"
	ProxyConfigRulesHeadersTypeOpaque          ProxyConfigRulesHeadersType = "opaque"
	ProxyConfigRulesHeadersTypeWorkspaceSecret ProxyConfigRulesHeadersType = "workspace_secret"
)

func (r ProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case ProxyConfigRulesHeadersTypePlaintext, ProxyConfigRulesHeadersTypeOpaque, ProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type ProxyConfigParam struct {
	AccessControl param.Field[ProxyConfigAccessControlParam] `json:"access_control"`
	Callbacks     param.Field[[]ProxyConfigCallbackParam]    `json:"callbacks"`
	// Description says what this configuration as a whole lets the sandbox reach,
	// complementing the per-rule descriptions. At most 1024 characters.
	Description param.Field[string]                 `json:"description"`
	NoProxy     param.Field[[]string]               `json:"no_proxy"`
	Rules       param.Field[[]ProxyConfigRuleParam] `json:"rules"`
}

func (r ProxyConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigAccessControlParam struct {
	AllowList param.Field[[]string] `json:"allow_list"`
	DenyList  param.Field[[]string] `json:"deny_list"`
}

func (r ProxyConfigAccessControlParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigCallbackParam struct {
	MatchHosts     param.Field[[]string]                                 `json:"match_hosts" api:"required"`
	TtlSeconds     param.Field[int64]                                    `json:"ttl_seconds" api:"required"`
	URL            param.Field[string]                                   `json:"url" api:"required"`
	FullRequest    param.Field[bool]                                     `json:"full_request"`
	RequestHeaders param.Field[[]ProxyConfigCallbacksRequestHeaderParam] `json:"request_headers"`
}

func (r ProxyConfigCallbackParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigCallbacksRequestHeaderParam struct {
	Name  param.Field[string]                                 `json:"name" api:"required"`
	Type  param.Field[ProxyConfigCallbacksRequestHeadersType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                   `json:"is_set"`
	Value param.Field[string]                                 `json:"value"`
}

func (r ProxyConfigCallbacksRequestHeaderParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRuleParam struct {
	Name  param.Field[string]                          `json:"name" api:"required"`
	Aws   param.Field[ProxyConfigRulesAwsUnionParam]   `json:"aws"`
	Azure param.Field[ProxyConfigRulesAzureUnionParam] `json:"azure"`
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
	EnvVars param.Field[map[string]string]             `json:"env_vars"`
	Gcp     param.Field[ProxyConfigRulesGcpUnionParam] `json:"gcp"`
	Headers param.Field[[]ProxyConfigRulesHeaderParam] `json:"headers"`
	// MatchHeaders restricts a header injection rule to requests carrying every listed
	// header, each written "name: value" with a lowercase name and an exact value.
	// Pair with headers of the same name to swap a placeholder the sandbox sends (e.g.
	// "authorization: Bearer account-b") for a real credential, so one host can serve
	// several accounts. Rules are evaluated in order and the first match wins.
	MatchHeaders param.Field[[]string] `json:"match_headers"`
	// MatchHosts is only accepted for header injection and relay rules. Provider auth
	// rules use built-in host matching. A relay rule may also list IP addresses and
	// CIDR blocks (10.0.0.0/8), which match connections the sandbox makes to an IP
	// address rather than a name.
	MatchHosts param.Field[[]string] `json:"match_hosts"`
	MatchPaths param.Field[[]string] `json:"match_paths"`
	// RelayID is the ID of the LangSmith relay (a connection of kind relay) that
	// carries connections to match_hosts when type is relay. Those connections leave
	// from the relay's connector, and its network and IP, instead of the sandbox host.
	// Access control still applies first, the first enabled matching relay rule wins,
	// and a matched host never falls back to a direct connection.
	RelayID param.Field[string] `json:"relay_id" format:"uuid"`
	Type    param.Field[string] `json:"type"`
}

func (r ProxyConfigRuleParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRulesAwsParam struct {
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

func (r ProxyConfigRulesAwsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAwsParam) implementsProxyConfigRulesAwsUnionParam() {}

// Satisfied by [ProxyConfigRulesAwsSandboxesProxyAwsRoleConfigParam],
// [ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigParam],
// [ProxyConfigRulesAwsParam].
type ProxyConfigRulesAwsUnionParam interface {
	implementsProxyConfigRulesAwsUnionParam()
}

type ProxyConfigRulesAwsSandboxesProxyAwsRoleConfigParam struct {
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

func (r ProxyConfigRulesAwsSandboxesProxyAwsRoleConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsRoleConfigParam) implementsProxyConfigRulesAwsUnionParam() {
}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigParam struct {
	AccessKeyID     param.Field[ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDParam]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyParam] `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace. The role's trust policy must trust
	// LangSmith's issuer for the token subject. Credentials renew automatically, and
	// access follows the role's effective AWS permissions. Mutually exclusive with
	// static keys, and cannot be changed after creation.
	RoleArn param.Field[ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigRoleArn] `json:"role_arn"`
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigParam) implementsProxyConfigRulesAwsUnionParam() {
}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDParam struct {
	Type  param.Field[ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                            `json:"is_set"`
	Value param.Field[string]                                                          `json:"value"`
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigAccessKeyIDParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyParam struct {
	Type  param.Field[ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                `json:"is_set"`
	Value param.Field[string]                                                              `json:"value"`
}

func (r ProxyConfigRulesAwsSandboxesProxyAwsStaticConfigSecretAccessKeyParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRulesAzureParam struct {
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

func (r ProxyConfigRulesAzureParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAzureParam) implementsProxyConfigRulesAzureUnionParam() {}

// Satisfied by
// [ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigParam],
// [ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigParam],
// [ProxyConfigRulesAzureSandboxesProxyAzureSasConfigParam],
// [ProxyConfigRulesAzureParam].
type ProxyConfigRulesAzureUnionParam interface {
	implementsProxyConfigRulesAzureUnionParam()
}

type ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigParam struct {
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
	Resources param.Field[[]ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource] `json:"resources"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope param.Field[string] `json:"subject_scope"`
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigParam) implementsProxyConfigRulesAzureUnionParam() {
}

type ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigParam struct {
	AccountKey param.Field[ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyParam] `json:"account_key" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigParam) implementsProxyConfigRulesAzureUnionParam() {
}

type ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyParam struct {
	Type  param.Field[ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                   `json:"is_set"`
	Value param.Field[string]                                                                 `json:"value"`
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureAccountKeyConfigAccountKeyParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRulesAzureSandboxesProxyAzureSasConfigParam struct {
	SasToken param.Field[ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenParam] `json:"sas_token" api:"required"`
	// StorageAccount is the storage account static credentials authenticate to:
	// requests to its blob, dfs, queue, table and file endpoints. Requires exactly one
	// of account_key (requests are signed with SharedKey) or sas_token (appended to
	// each request).
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureSasConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureSasConfigParam) implementsProxyConfigRulesAzureUnionParam() {
}

type ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenParam struct {
	Type  param.Field[ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                          `json:"is_set"`
	Value param.Field[string]                                                        `json:"value"`
}

func (r ProxyConfigRulesAzureSandboxesProxyAzureSasConfigSasTokenParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRulesGcpParam struct {
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

func (r ProxyConfigRulesGcpParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesGcpParam) implementsProxyConfigRulesGcpUnionParam() {}

// Satisfied by [ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigParam],
// [ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigParam],
// [ProxyConfigRulesGcpParam].
type ProxyConfigRulesGcpUnionParam interface {
	implementsProxyConfigRulesGcpUnionParam()
}

type ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigParam struct {
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

func (r ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigParam) implementsProxyConfigRulesGcpUnionParam() {
}

type ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigParam struct {
	ServiceAccountJson param.Field[ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonParam] `json:"service_account_json" api:"required"`
	// Scopes are the OAuth scopes requested for the access token. Omit them to use
	// https://www.googleapis.com/auth/cloud-platform.
	Scopes param.Field[[]string] `json:"scopes"`
}

func (r ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigParam) implementsProxyConfigRulesGcpUnionParam() {
}

type ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonParam struct {
	Type  param.Field[ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                           `json:"is_set"`
	Value param.Field[string]                                                                         `json:"value"`
}

func (r ProxyConfigRulesGcpSandboxesProxyGcpServiceAccountConfigServiceAccountJsonParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProxyConfigRulesHeaderParam struct {
	Name  param.Field[string]                      `json:"name" api:"required"`
	Type  param.Field[ProxyConfigRulesHeadersType] `json:"type" api:"required"`
	IsSet param.Field[bool]                        `json:"is_set"`
	Value param.Field[string]                      `json:"value"`
}

func (r ProxyConfigRulesHeaderParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
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

type SandboxMountConfig struct {
	Auth   SandboxMountConfigAuth    `json:"auth"`
	Mounts []SandboxMountConfigMount `json:"mounts"`
	JSON   sandboxMountConfigJSON    `json:"-"`
}

// sandboxMountConfigJSON contains the JSON metadata for the struct
// [SandboxMountConfig]
type sandboxMountConfigJSON struct {
	Auth        apijson.Field
	Mounts      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuth struct {
	Aws   SandboxMountConfigAuthAws   `json:"aws"`
	Azure SandboxMountConfigAuthAzure `json:"azure"`
	Gcp   SandboxMountConfigAuthGcp   `json:"gcp"`
	JSON  sandboxMountConfigAuthJSON  `json:"-"`
}

// sandboxMountConfigAuthJSON contains the JSON metadata for the struct
// [SandboxMountConfigAuth]
type sandboxMountConfigAuthJSON struct {
	Aws         apijson.Field
	Azure       apijson.Field
	Gcp         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuth) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuthAws struct {
	// This field can have the runtime type of
	// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID].
	AccessKeyID interface{} `json:"access_key_id"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn string `json:"role_arn"`
	// This field can have the runtime type of
	// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey].
	SecretAccessKey interface{} `json:"secret_access_key"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                        `json:"subject_scope"`
	JSON         sandboxMountConfigAuthAwsJSON `json:"-"`
	union        SandboxMountConfigAuthAwsUnion
}

// sandboxMountConfigAuthAwsJSON contains the JSON metadata for the struct
// [SandboxMountConfigAuthAws]
type sandboxMountConfigAuthAwsJSON struct {
	AccessKeyID     apijson.Field
	RoleArn         apijson.Field
	SecretAccessKey apijson.Field
	SubjectScope    apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r sandboxMountConfigAuthAwsJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxMountConfigAuthAws) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxMountConfigAuthAws{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxMountConfigAuthAwsUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig],
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig].
func (r SandboxMountConfigAuthAws) AsUnion() SandboxMountConfigAuthAwsUnion {
	return r.union
}

// Union satisfied by
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig] or
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig].
type SandboxMountConfigAuthAwsUnion interface {
	implementsSandboxMountConfigAuthAws()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxMountConfigAuthAwsUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig{}),
		},
	)
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig struct {
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn string `json:"role_arn" api:"required"`
	// SubjectScope picks the token subject the role trusts: "workspace" (the default)
	// or "tag:<tag value id>", which must be one of the sandbox's tag_value_ids. Only
	// valid with role_arn.
	SubjectScope string                                                              `json:"subject_scope"`
	JSON         sandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON `json:"-"`
}

// sandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig]
type sandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON struct {
	RoleArn      apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfig) implementsSandboxMountConfigAuthAws() {
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig struct {
	AccessKeyID     SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID     `json:"access_key_id" api:"required"`
	SecretAccessKey SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn `json:"role_arn"`
	JSON    sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON    `json:"-"`
}

// sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig]
type sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON struct {
	AccessKeyID     apijson.Field
	SecretAccessKey apijson.Field
	RoleArn         apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfig) implementsSandboxMountConfigAuthAws() {
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID struct {
	Type  SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType `json:"type" api:"required"`
	IsSet bool                                                                             `json:"is_set"`
	Value string                                                                           `json:"value"`
	JSON  sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON `json:"-"`
}

// sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID]
type sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyID) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType string

const (
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypePlaintext       SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "plaintext"
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeOpaque          SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "opaque"
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeWorkspaceSecret SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType = "workspace_secret"
)

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType) IsKnown() bool {
	switch r {
	case SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypePlaintext, SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeOpaque, SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey struct {
	Type  SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType `json:"type" api:"required"`
	IsSet bool                                                                                 `json:"is_set"`
	Value string                                                                               `json:"value"`
	JSON  sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON `json:"-"`
}

// sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey]
type sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType string

const (
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypePlaintext       SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "plaintext"
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeOpaque          SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "opaque"
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeWorkspaceSecret SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType = "workspace_secret"
)

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType) IsKnown() bool {
	switch r {
	case SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypePlaintext, SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeOpaque, SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
// token LangSmith signs for this workspace; the role's trust policy must trust
// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
// be changed after creation.
type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn string

const (
	SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArnEmpty SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn = ""
)

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn) IsKnown() bool {
	switch r {
	case SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArnEmpty:
		return true
	}
	return false
}

type SandboxMountConfigAuthAzure struct {
	// This field can have the runtime type of
	// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey].
	AccountKey interface{} `json:"account_key"`
	// ClientID is the application (client) ID, a GUID, of the app registration or
	// user-assigned managed identity whose federated identity credential trusts
	// LangSmith's issuer. Required with tenant_id.
	ClientID string `json:"client_id"`
	// This field can have the runtime type of
	// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken].
	SasToken interface{} `json:"sas_token"`
	// SubjectScope picks the token subject the federated credential trusts:
	// "workspace" (the default) or "tag:<tag value id>", which must be one of the
	// sandbox's tag_value_ids. Only valid with tenant_id.
	SubjectScope string `json:"subject_scope"`
	// TenantID is the Entra tenant (directory) ID, a GUID. Setting it selects workload
	// identity: LangSmith signs an OIDC token for this workspace and exchanges it at
	// Entra ID for storage access tokens of client_id. Cannot be changed after
	// creation.
	TenantID string                          `json:"tenant_id"`
	JSON     sandboxMountConfigAuthAzureJSON `json:"-"`
	union    SandboxMountConfigAuthAzureUnion
}

// sandboxMountConfigAuthAzureJSON contains the JSON metadata for the struct
// [SandboxMountConfigAuthAzure]
type sandboxMountConfigAuthAzureJSON struct {
	AccountKey   apijson.Field
	ClientID     apijson.Field
	SasToken     apijson.Field
	SubjectScope apijson.Field
	TenantID     apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r sandboxMountConfigAuthAzureJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxMountConfigAuthAzure) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxMountConfigAuthAzure{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxMountConfigAuthAzureUnion] interface which you can
// cast to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig],
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig],
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig].
func (r SandboxMountConfigAuthAzure) AsUnion() SandboxMountConfigAuthAzureUnion {
	return r.union
}

// Union satisfied by
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig],
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig] or
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig].
type SandboxMountConfigAuthAzureUnion interface {
	implementsSandboxMountConfigAuthAzure()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxMountConfigAuthAzureUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig{}),
		},
	)
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig struct {
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
	SubjectScope string                                                                              `json:"subject_scope"`
	JSON         sandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON `json:"-"`
}

// sandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig]
type sandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON struct {
	ClientID     apijson.Field
	TenantID     apijson.Field
	SubjectScope apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfig) implementsSandboxMountConfigAuthAzure() {
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig struct {
	AccountKey SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey `json:"account_key" api:"required"`
	JSON       sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON       `json:"-"`
}

// sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig]
type sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON struct {
	AccountKey  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfig) implementsSandboxMountConfigAuthAzure() {
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey struct {
	Type  SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType `json:"type" api:"required"`
	IsSet bool                                                                                    `json:"is_set"`
	Value string                                                                                  `json:"value"`
	JSON  sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON `json:"-"`
}

// sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey]
type sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKey) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType string

const (
	SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypePlaintext       SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "plaintext"
	SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeOpaque          SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "opaque"
	SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeWorkspaceSecret SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType = "workspace_secret"
)

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType) IsKnown() bool {
	switch r {
	case SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypePlaintext, SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeOpaque, SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig struct {
	SasToken SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken `json:"sas_token" api:"required"`
	JSON     sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON     `json:"-"`
}

// sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig]
type sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON struct {
	SasToken    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfig) implementsSandboxMountConfigAuthAzure() {
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken struct {
	Type  SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType `json:"type" api:"required"`
	IsSet bool                                                                           `json:"is_set"`
	Value string                                                                         `json:"value"`
	JSON  sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON `json:"-"`
}

// sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken]
type sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasToken) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType string

const (
	SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypePlaintext       SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "plaintext"
	SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeOpaque          SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "opaque"
	SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeWorkspaceSecret SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType = "workspace_secret"
)

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType) IsKnown() bool {
	switch r {
	case SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypePlaintext, SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeOpaque, SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxMountConfigAuthGcp struct {
	// ServiceAccountEmail is the service account to impersonate with the federated
	// token. Omit it to use the federated principal directly. Only valid with
	// workload_identity_provider.
	ServiceAccountEmail string `json:"service_account_email"`
	// This field can have the runtime type of
	// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson].
	ServiceAccountJson interface{} `json:"service_account_json"`
	// SubjectScope picks the token subject the provider trusts: "workspace" (the
	// default) or "tag:<tag value id>", which must be one of the sandbox's
	// tag_value_ids. Only valid with workload_identity_provider.
	SubjectScope string `json:"subject_scope"`
	// WorkloadIdentityProvider federates an OIDC token LangSmith signs for this
	// workspace through Workload Identity Federation instead of a service account key:
	// //iam.googleapis.com/projects/<number>/locations/global/workloadIdentityPools/<pool>/providers/<provider>.
	// Mutually exclusive with service_account_json. Cannot be changed after creation.
	WorkloadIdentityProvider string                        `json:"workload_identity_provider"`
	JSON                     sandboxMountConfigAuthGcpJSON `json:"-"`
	union                    SandboxMountConfigAuthGcpUnion
}

// sandboxMountConfigAuthGcpJSON contains the JSON metadata for the struct
// [SandboxMountConfigAuthGcp]
type sandboxMountConfigAuthGcpJSON struct {
	ServiceAccountEmail      apijson.Field
	ServiceAccountJson       apijson.Field
	SubjectScope             apijson.Field
	WorkloadIdentityProvider apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r sandboxMountConfigAuthGcpJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxMountConfigAuthGcp) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxMountConfigAuthGcp{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxMountConfigAuthGcpUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig],
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig].
func (r SandboxMountConfigAuthGcp) AsUnion() SandboxMountConfigAuthGcpUnion {
	return r.union
}

// Union satisfied by
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig] or
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig].
type SandboxMountConfigAuthGcpUnion interface {
	implementsSandboxMountConfigAuthGcp()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxMountConfigAuthGcpUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig{}),
		},
	)
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig struct {
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
	SubjectScope string                                                                          `json:"subject_scope"`
	JSON         sandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON `json:"-"`
}

// sandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig]
type sandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON struct {
	WorkloadIdentityProvider apijson.Field
	ServiceAccountEmail      apijson.Field
	SubjectScope             apijson.Field
	raw                      string
	ExtraFields              map[string]apijson.Field
}

func (r *SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfig) implementsSandboxMountConfigAuthGcp() {
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig struct {
	ServiceAccountJson SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson `json:"service_account_json" api:"required"`
	JSON               sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON               `json:"-"`
}

// sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig]
type sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON struct {
	ServiceAccountJson apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfig) implementsSandboxMountConfigAuthGcp() {
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson struct {
	Type  SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType `json:"type" api:"required"`
	IsSet bool                                                                                            `json:"is_set"`
	Value string                                                                                          `json:"value"`
	JSON  sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON `json:"-"`
}

// sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson]
type sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON struct {
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJson) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType string

const (
	SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypePlaintext       SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "plaintext"
	SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeOpaque          SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "opaque"
	SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeWorkspaceSecret SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType = "workspace_secret"
)

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType) IsKnown() bool {
	switch r {
	case SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypePlaintext, SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeOpaque, SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxMountConfigMount struct {
	ID        string                       `json:"id" api:"required"`
	MountPath string                       `json:"mount_path" api:"required"`
	Type      SandboxMountConfigMountsType `json:"type" api:"required"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob].
	AzureBlob interface{} `json:"azure_blob"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecCache],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCache],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecCache],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCache],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache].
	Cache interface{} `json:"cache"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthub],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthub],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthub],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub].
	Contexthub interface{} `json:"contexthub"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryption],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryption],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryption],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption].
	Encryption interface{} `json:"encryption"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcs],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcs],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcs],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcs],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs].
	Gcs interface{} `json:"gcs"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecGit],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGit],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecGit],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGit],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit].
	Git      interface{} `json:"git"`
	ReadOnly bool        `json:"read_only"`
	// This field can have the runtime type of
	// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3],
	// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3],
	// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3],
	// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3],
	// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3].
	S3    interface{}                 `json:"s3"`
	JSON  sandboxMountConfigMountJSON `json:"-"`
	union SandboxMountConfigMountsUnion
}

// sandboxMountConfigMountJSON contains the JSON metadata for the struct
// [SandboxMountConfigMount]
type sandboxMountConfigMountJSON struct {
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

func (r sandboxMountConfigMountJSON) RawJSON() string {
	return r.raw
}

func (r *SandboxMountConfigMount) UnmarshalJSON(data []byte) (err error) {
	*r = SandboxMountConfigMount{}
	err = apijson.UnmarshalRoot(data, &r.union)
	if err != nil {
		return err
	}
	return apijson.Port(r.union, &r)
}

// AsUnion returns a [SandboxMountConfigMountsUnion] interface which you can cast
// to the specific types for more type safety.
//
// Possible runtime types of the union are
// [SandboxMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpec],
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec],
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec].
func (r SandboxMountConfigMount) AsUnion() SandboxMountConfigMountsUnion {
	return r.union
}

// Union satisfied by [SandboxMountConfigMountsSandboxapiS3BucketMountSpec],
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpec],
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpec],
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec] or
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec].
type SandboxMountConfigMountsUnion interface {
	implementsSandboxMountConfigMount()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*SandboxMountConfigMountsUnion)(nil)).Elem(),
		"type",
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "azure_blob",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiS3BucketMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "azure_blob",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGcsBucketMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "azure_blob",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiGitRepoMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "azure_blob",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "s3",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "gcs",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "azure_blob",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "git",
		},
		apijson.UnionVariant{
			TypeFilter:         gjson.JSON,
			Type:               reflect.TypeOf(SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec{}),
			DiscriminatorValue: "contexthub",
		},
	)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpec struct {
	ID         string                                                        `json:"id" api:"required"`
	MountPath  string                                                        `json:"mount_path" api:"required"`
	S3         SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3         `json:"s3" api:"required"`
	Type       SandboxMountConfigMountsSandboxapiS3BucketMountSpecType       `json:"type" api:"required"`
	AzureBlob  SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxMountConfigMountsSandboxapiS3BucketMountSpecCache      `json:"cache"`
	Contexthub SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthub `json:"contexthub"`
	Encryption SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryption `json:"encryption"`
	Gcs        SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcs        `json:"gcs"`
	Git        SandboxMountConfigMountsSandboxapiS3BucketMountSpecGit        `json:"git"`
	ReadOnly   bool                                                          `json:"read_only"`
	JSON       sandboxMountConfigMountsSandboxapiS3BucketMountSpecJSON       `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiS3BucketMountSpec]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecJSON struct {
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

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpec) implementsSandboxMountConfigMount() {}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3 struct {
	Bucket      string                                                    `json:"bucket" api:"required"`
	Region      string                                                    `json:"region" api:"required"`
	EndpointURL string                                                    `json:"endpoint_url"`
	PathStyle   bool                                                      `json:"path_style"`
	Prefix      string                                                    `json:"prefix"`
	JSON        sandboxMountConfigMountsSandboxapiS3BucketMountSpecS3JSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecS3JSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecType string

const (
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeS3         SandboxMountConfigMountsSandboxapiS3BucketMountSpecType = "s3"
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs        SandboxMountConfigMountsSandboxapiS3BucketMountSpecType = "gcs"
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeAzureBlob  SandboxMountConfigMountsSandboxapiS3BucketMountSpecType = "azure_blob"
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeGit        SandboxMountConfigMountsSandboxapiS3BucketMountSpecType = "git"
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub SandboxMountConfigMountsSandboxapiS3BucketMountSpecType = "contexthub"
)

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeS3, SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeGcs, SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeAzureBlob, SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeGit, SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                           `json:"storage_account" api:"required"`
	Prefix         string                                                           `json:"prefix"`
	JSON           sandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecCache struct {
	MaxSizeBytes     int64                                                        `json:"max_size_bytes"`
	WritebackSeconds int64                                                        `json:"writeback_seconds"`
	JSON             sandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecCache]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                              `json:"initial_pull_only"`
	JSON            sandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthub]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                            `json:"kms_key" api:"required"`
	JSON   sandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryption]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcs struct {
	Bucket string                                                     `json:"bucket" api:"required"`
	Prefix string                                                     `json:"prefix"`
	JSON   sandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcs]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGit struct {
	RemoteURL              string                                                     `json:"remote_url" api:"required"`
	Ref                    SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                      `json:"refresh_interval_seconds"`
	JSON                   sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiS3BucketMountSpecGit]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRef struct {
	Name string                                                        `json:"name" api:"required"`
	Type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRef]
type sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefType string

const (
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "branch"
	SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag    SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefType = "tag"
)

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch, SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpec struct {
	ID         string                                                         `json:"id" api:"required"`
	Gcs        SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcs        `json:"gcs" api:"required"`
	MountPath  string                                                         `json:"mount_path" api:"required"`
	Type       SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType       `json:"type" api:"required"`
	AzureBlob  SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCache      `json:"cache"`
	Contexthub SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthub `json:"contexthub"`
	Encryption SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryption `json:"encryption"`
	Git        SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGit        `json:"git"`
	ReadOnly   bool                                                           `json:"read_only"`
	S3         SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3         `json:"s3"`
	JSON       sandboxMountConfigMountsSandboxapiGcsBucketMountSpecJSON       `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiGcsBucketMountSpec]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecJSON struct {
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

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpec) implementsSandboxMountConfigMount() {}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcs struct {
	Bucket string                                                      `json:"bucket" api:"required"`
	Prefix string                                                      `json:"prefix"`
	JSON   sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcs]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType string

const (
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3         SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType = "s3"
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs        SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType = "gcs"
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeAzureBlob  SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType = "azure_blob"
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit        SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType = "git"
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType = "contexthub"
)

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeS3, SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeGcs, SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeAzureBlob, SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeGit, SandboxMountConfigMountsSandboxapiGcsBucketMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                            `json:"storage_account" api:"required"`
	Prefix         string                                                            `json:"prefix"`
	JSON           sandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCache struct {
	MaxSizeBytes     int64                                                         `json:"max_size_bytes"`
	WritebackSeconds int64                                                         `json:"writeback_seconds"`
	JSON             sandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCache]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                               `json:"initial_pull_only"`
	JSON            sandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthub]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                             `json:"kms_key" api:"required"`
	JSON   sandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryption]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGit struct {
	RemoteURL              string                                                      `json:"remote_url" api:"required"`
	Ref                    SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                       `json:"refresh_interval_seconds"`
	JSON                   sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGit]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRef struct {
	Name string                                                         `json:"name" api:"required"`
	Type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRef]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType string

const (
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "branch"
	SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag    SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType = "tag"
)

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeBranch, SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3 struct {
	Bucket      string                                                     `json:"bucket" api:"required"`
	Region      string                                                     `json:"region" api:"required"`
	EndpointURL string                                                     `json:"endpoint_url"`
	PathStyle   bool                                                       `json:"path_style"`
	Prefix      string                                                     `json:"prefix"`
	JSON        sandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3]
type sandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpec struct {
	ID         string                                                       `json:"id" api:"required"`
	Git        SandboxMountConfigMountsSandboxapiGitRepoMountSpecGit        `json:"git" api:"required"`
	MountPath  string                                                       `json:"mount_path" api:"required"`
	Type       SandboxMountConfigMountsSandboxapiGitRepoMountSpecType       `json:"type" api:"required"`
	AzureBlob  SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxMountConfigMountsSandboxapiGitRepoMountSpecCache      `json:"cache"`
	Contexthub SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthub `json:"contexthub"`
	Encryption SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryption `json:"encryption"`
	Gcs        SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcs        `json:"gcs"`
	ReadOnly   bool                                                         `json:"read_only"`
	S3         SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3         `json:"s3"`
	JSON       sandboxMountConfigMountsSandboxapiGitRepoMountSpecJSON       `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiGitRepoMountSpec]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecJSON struct {
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

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpec) implementsSandboxMountConfigMount() {}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGit struct {
	RemoteURL              string                                                    `json:"remote_url" api:"required"`
	Ref                    SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                     `json:"refresh_interval_seconds"`
	JSON                   sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiGitRepoMountSpecGit]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRef struct {
	Name string                                                       `json:"name" api:"required"`
	Type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRef]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefType string

const (
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "branch"
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag    SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefType = "tag"
)

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeBranch, SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecType string

const (
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeS3         SandboxMountConfigMountsSandboxapiGitRepoMountSpecType = "s3"
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs        SandboxMountConfigMountsSandboxapiGitRepoMountSpecType = "gcs"
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeAzureBlob  SandboxMountConfigMountsSandboxapiGitRepoMountSpecType = "azure_blob"
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeGit        SandboxMountConfigMountsSandboxapiGitRepoMountSpecType = "git"
	SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub SandboxMountConfigMountsSandboxapiGitRepoMountSpecType = "contexthub"
)

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeS3, SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeGcs, SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeAzureBlob, SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeGit, SandboxMountConfigMountsSandboxapiGitRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                          `json:"storage_account" api:"required"`
	Prefix         string                                                          `json:"prefix"`
	JSON           sandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                       `json:"max_size_bytes"`
	WritebackSeconds int64                                                       `json:"writeback_seconds"`
	JSON             sandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecCache]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                             `json:"initial_pull_only"`
	JSON            sandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthub]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                           `json:"kms_key" api:"required"`
	JSON   sandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryption]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcs struct {
	Bucket string                                                    `json:"bucket" api:"required"`
	Prefix string                                                    `json:"prefix"`
	JSON   sandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcs]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3 struct {
	Bucket      string                                                   `json:"bucket" api:"required"`
	Region      string                                                   `json:"region" api:"required"`
	EndpointURL string                                                   `json:"endpoint_url"`
	PathStyle   bool                                                     `json:"path_style"`
	Prefix      string                                                   `json:"prefix"`
	JSON        sandboxMountConfigMountsSandboxapiGitRepoMountSpecS3JSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiGitRepoMountSpecS3JSON contains the JSON
// metadata for the struct [SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3]
type sandboxMountConfigMountsSandboxapiGitRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiGitRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec struct {
	ID         string                                                              `json:"id" api:"required"`
	Contexthub SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub `json:"contexthub" api:"required"`
	MountPath  string                                                              `json:"mount_path" api:"required"`
	Type       SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType       `json:"type" api:"required"`
	AzureBlob  SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob  `json:"azure_blob"`
	Cache      SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCache      `json:"cache"`
	Encryption SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption `json:"encryption"`
	Gcs        SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcs        `json:"gcs"`
	Git        SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGit        `json:"git"`
	ReadOnly   bool                                                                `json:"read_only"`
	S3         SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3         `json:"s3"`
	JSON       sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecJSON       `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecJSON contains the JSON
// metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecJSON struct {
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

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpec) implementsSandboxMountConfigMount() {
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                    `json:"initial_pull_only"`
	JSON            sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType string

const (
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3         SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType = "s3"
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs        SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType = "gcs"
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeAzureBlob  SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType = "azure_blob"
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit        SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType = "git"
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType = "contexthub"
)

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeS3, SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGcs, SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeAzureBlob, SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeGit, SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                 `json:"storage_account" api:"required"`
	Prefix         string                                                                 `json:"prefix"`
	JSON           sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCache struct {
	MaxSizeBytes     int64                                                              `json:"max_size_bytes"`
	WritebackSeconds int64                                                              `json:"writeback_seconds"`
	JSON             sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCache]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                  `json:"kms_key" api:"required"`
	JSON   sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcs struct {
	Bucket string                                                           `json:"bucket" api:"required"`
	Prefix string                                                           `json:"prefix"`
	JSON   sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcs]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGit struct {
	RemoteURL              string                                                           `json:"remote_url" api:"required"`
	Ref                    SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                            `json:"refresh_interval_seconds"`
	JSON                   sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGit]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef struct {
	Name string                                                              `json:"name" api:"required"`
	Type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType string

const (
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeBranch SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType = "branch"
	SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeTag    SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType = "tag"
)

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeBranch, SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3 struct {
	Bucket      string                                                          `json:"bucket" api:"required"`
	Region      string                                                          `json:"region" api:"required"`
	EndpointURL string                                                          `json:"endpoint_url"`
	PathStyle   bool                                                            `json:"path_style"`
	Prefix      string                                                          `json:"prefix"`
	JSON        sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3]
type sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec struct {
	ID         string                                                                  `json:"id" api:"required"`
	AzureBlob  SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob  `json:"azure_blob" api:"required"`
	MountPath  string                                                                  `json:"mount_path" api:"required"`
	Type       SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType       `json:"type" api:"required"`
	Cache      SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache      `json:"cache"`
	Contexthub SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub `json:"contexthub"`
	Encryption SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption `json:"encryption"`
	Gcs        SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs        `json:"gcs"`
	Git        SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit        `json:"git"`
	ReadOnly   bool                                                                    `json:"read_only"`
	S3         SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3         `json:"s3"`
	JSON       sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON       `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON struct {
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

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecJSON) RawJSON() string {
	return r.raw
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpec) implementsSandboxMountConfigMount() {
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob struct {
	Container string `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount string                                                                     `json:"storage_account" api:"required"`
	Prefix         string                                                                     `json:"prefix"`
	JSON           sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON struct {
	Container      apijson.Field
	StorageAccount apijson.Field
	Prefix         apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlob) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType string

const (
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeS3         SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "s3"
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGcs        SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "gcs"
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeAzureBlob  SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "azure_blob"
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGit        SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "git"
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeContexthub SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType = "contexthub"
)

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeS3, SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGcs, SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeAzureBlob, SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeGit, SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecTypeContexthub:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache struct {
	MaxSizeBytes     int64                                                                  `json:"max_size_bytes"`
	WritebackSeconds int64                                                                  `json:"writeback_seconds"`
	JSON             sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON struct {
	MaxSizeBytes     apijson.Field
	WritebackSeconds apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCache) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo string `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly bool                                                                        `json:"initial_pull_only"`
	JSON            sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON struct {
	Repo            apijson.Field
	InitialPullOnly apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthub) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey string                                                                      `json:"kms_key" api:"required"`
	JSON   sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON
// contains the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON struct {
	KmsKey      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryption) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs struct {
	Bucket string                                                               `json:"bucket" api:"required"`
	Prefix string                                                               `json:"prefix"`
	JSON   sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON struct {
	Bucket      apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcs) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit struct {
	RemoteURL              string                                                               `json:"remote_url" api:"required"`
	Ref                    SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef  `json:"ref"`
	RefreshIntervalSeconds int64                                                                `json:"refresh_interval_seconds"`
	JSON                   sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON struct {
	RemoteURL              apijson.Field
	Ref                    apijson.Field
	RefreshIntervalSeconds apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGit) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef struct {
	Name string                                                                  `json:"name" api:"required"`
	Type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType `json:"type" api:"required"`
	JSON sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON contains
// the JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRef) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefJSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType string

const (
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeBranch SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType = "branch"
	SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeTag    SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType = "tag"
)

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeBranch, SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefTypeTag:
		return true
	}
	return false
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3 struct {
	Bucket      string                                                              `json:"bucket" api:"required"`
	Region      string                                                              `json:"region" api:"required"`
	EndpointURL string                                                              `json:"endpoint_url"`
	PathStyle   bool                                                                `json:"path_style"`
	Prefix      string                                                              `json:"prefix"`
	JSON        sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON `json:"-"`
}

// sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON contains the
// JSON metadata for the struct
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3]
type sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON struct {
	Bucket      apijson.Field
	Region      apijson.Field
	EndpointURL apijson.Field
	PathStyle   apijson.Field
	Prefix      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3JSON) RawJSON() string {
	return r.raw
}

type SandboxMountConfigMountsType string

const (
	SandboxMountConfigMountsTypeS3         SandboxMountConfigMountsType = "s3"
	SandboxMountConfigMountsTypeGcs        SandboxMountConfigMountsType = "gcs"
	SandboxMountConfigMountsTypeAzureBlob  SandboxMountConfigMountsType = "azure_blob"
	SandboxMountConfigMountsTypeGit        SandboxMountConfigMountsType = "git"
	SandboxMountConfigMountsTypeContexthub SandboxMountConfigMountsType = "contexthub"
)

func (r SandboxMountConfigMountsType) IsKnown() bool {
	switch r {
	case SandboxMountConfigMountsTypeS3, SandboxMountConfigMountsTypeGcs, SandboxMountConfigMountsTypeAzureBlob, SandboxMountConfigMountsTypeGit, SandboxMountConfigMountsTypeContexthub:
		return true
	}
	return false
}

type SandboxMountConfigParam struct {
	Auth   param.Field[SandboxMountConfigAuthParam]          `json:"auth"`
	Mounts param.Field[[]SandboxMountConfigMountsUnionParam] `json:"mounts"`
}

func (r SandboxMountConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigAuthParam struct {
	Aws   param.Field[SandboxMountConfigAuthAwsUnionParam]   `json:"aws"`
	Azure param.Field[SandboxMountConfigAuthAzureUnionParam] `json:"azure"`
	Gcp   param.Field[SandboxMountConfigAuthGcpUnionParam]   `json:"gcp"`
}

func (r SandboxMountConfigAuthParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigAuthAwsParam struct {
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

func (r SandboxMountConfigAuthAwsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAwsParam) implementsSandboxMountConfigAuthAwsUnionParam() {}

// Satisfied by
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigParam],
// [SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigParam],
// [SandboxMountConfigAuthAwsParam].
type SandboxMountConfigAuthAwsUnionParam interface {
	implementsSandboxMountConfigAuthAwsUnionParam()
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigParam struct {
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

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigParam) implementsSandboxMountConfigAuthAwsUnionParam() {
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigParam struct {
	AccessKeyID     param.Field[SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDParam]     `json:"access_key_id" api:"required"`
	SecretAccessKey param.Field[SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyParam] `json:"secret_access_key" api:"required"`
	// IAM role to assume with AssumeRoleWithWebIdentity, using a short-lived OIDC
	// token LangSmith signs for this workspace; the role's trust policy must trust
	// LangSmith's issuer for the token subject. LangSmith scopes the credentials to
	// the configured S3 mounts. Mutually exclusive with static credentials, and cannot
	// be changed after creation.
	RoleArn param.Field[SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigRoleArn] `json:"role_arn"`
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigParam) implementsSandboxMountConfigAuthAwsUnionParam() {
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDParam struct {
	Type  param.Field[SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                             `json:"is_set"`
	Value param.Field[string]                                                                           `json:"value"`
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigAccessKeyIDParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyParam struct {
	Type  param.Field[SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                 `json:"is_set"`
	Value param.Field[string]                                                                               `json:"value"`
}

func (r SandboxMountConfigAuthAwsSandboxesSandboxAwsMountStaticAuthConfigSecretAccessKeyParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigAuthAzureParam struct {
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

func (r SandboxMountConfigAuthAzureParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAzureParam) implementsSandboxMountConfigAuthAzureUnionParam() {}

// Satisfied by
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigParam],
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigParam],
// [SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigParam],
// [SandboxMountConfigAuthAzureParam].
type SandboxMountConfigAuthAzureUnionParam interface {
	implementsSandboxMountConfigAuthAzureUnionParam()
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigParam struct {
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

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigParam) implementsSandboxMountConfigAuthAzureUnionParam() {
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigParam struct {
	AccountKey param.Field[SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyParam] `json:"account_key" api:"required"`
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigParam) implementsSandboxMountConfigAuthAzureUnionParam() {
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyParam struct {
	Type  param.Field[SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                    `json:"is_set"`
	Value param.Field[string]                                                                                  `json:"value"`
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountAccountKeyAuthConfigAccountKeyParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigParam struct {
	SasToken param.Field[SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenParam] `json:"sas_token" api:"required"`
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigParam) implementsSandboxMountConfigAuthAzureUnionParam() {
}

type SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenParam struct {
	Type  param.Field[SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                           `json:"is_set"`
	Value param.Field[string]                                                                         `json:"value"`
}

func (r SandboxMountConfigAuthAzureSandboxesSandboxAzureMountSasAuthConfigSasTokenParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigAuthGcpParam struct {
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

func (r SandboxMountConfigAuthGcpParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthGcpParam) implementsSandboxMountConfigAuthGcpUnionParam() {}

// Satisfied by
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigParam],
// [SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigParam],
// [SandboxMountConfigAuthGcpParam].
type SandboxMountConfigAuthGcpUnionParam interface {
	implementsSandboxMountConfigAuthGcpUnionParam()
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigParam struct {
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

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigParam) implementsSandboxMountConfigAuthGcpUnionParam() {
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigParam struct {
	ServiceAccountJson param.Field[SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonParam] `json:"service_account_json" api:"required"`
}

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigParam) implementsSandboxMountConfigAuthGcpUnionParam() {
}

type SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonParam struct {
	Type  param.Field[SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonType] `json:"type" api:"required"`
	IsSet param.Field[bool]                                                                                            `json:"is_set"`
	Value param.Field[string]                                                                                          `json:"value"`
}

func (r SandboxMountConfigAuthGcpSandboxesSandboxGcpMountServiceAccountAuthConfigServiceAccountJsonParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountParam struct {
	ID         param.Field[string]                       `json:"id" api:"required"`
	MountPath  param.Field[string]                       `json:"mount_path" api:"required"`
	Type       param.Field[SandboxMountConfigMountsType] `json:"type" api:"required"`
	AzureBlob  param.Field[interface{}]                  `json:"azure_blob"`
	Cache      param.Field[interface{}]                  `json:"cache"`
	Contexthub param.Field[interface{}]                  `json:"contexthub"`
	Encryption param.Field[interface{}]                  `json:"encryption"`
	Gcs        param.Field[interface{}]                  `json:"gcs"`
	Git        param.Field[interface{}]                  `json:"git"`
	ReadOnly   param.Field[bool]                         `json:"read_only"`
	S3         param.Field[interface{}]                  `json:"s3"`
}

func (r SandboxMountConfigMountParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigMountParam) implementsSandboxMountConfigMountsUnionParam() {}

// Satisfied by [SandboxMountConfigMountsSandboxapiS3BucketMountSpecParam],
// [SandboxMountConfigMountsSandboxapiGcsBucketMountSpecParam],
// [SandboxMountConfigMountsSandboxapiGitRepoMountSpecParam],
// [SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecParam],
// [SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecParam],
// [SandboxMountConfigMountParam].
type SandboxMountConfigMountsUnionParam interface {
	implementsSandboxMountConfigMountsUnionParam()
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecParam struct {
	ID         param.Field[string]                                                             `json:"id" api:"required"`
	MountPath  param.Field[string]                                                             `json:"mount_path" api:"required"`
	S3         param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3Param]         `json:"s3" api:"required"`
	Type       param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecType]            `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobParam]  `json:"azure_blob"`
	Cache      param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheParam]      `json:"cache"`
	Contexthub param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubParam] `json:"contexthub"`
	Encryption param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionParam] `json:"encryption"`
	Gcs        param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsParam]        `json:"gcs"`
	Git        param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitParam]        `json:"git"`
	ReadOnly   param.Field[bool]                                                               `json:"read_only"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecParam) implementsSandboxMountConfigMountsUnionParam() {
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3Param struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobParam struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheParam struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubParam struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionParam struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsParam struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitParam struct {
	RemoteURL              param.Field[string]                                                         `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefParam] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                          `json:"refresh_interval_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefParam struct {
	Name param.Field[string]                                                        `json:"name" api:"required"`
	Type param.Field[SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecParam struct {
	ID         param.Field[string]                                                              `json:"id" api:"required"`
	Gcs        param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsParam]        `json:"gcs" api:"required"`
	MountPath  param.Field[string]                                                              `json:"mount_path" api:"required"`
	Type       param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecType]            `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobParam]  `json:"azure_blob"`
	Cache      param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheParam]      `json:"cache"`
	Contexthub param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubParam] `json:"contexthub"`
	Encryption param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionParam] `json:"encryption"`
	Git        param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitParam]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                `json:"read_only"`
	S3         param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3Param]         `json:"s3"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecParam) implementsSandboxMountConfigMountsUnionParam() {
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsParam struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGcsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobParam struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecAzureBlobParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheParam struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecCacheParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubParam struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecContexthubParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionParam struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecEncryptionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitParam struct {
	RemoteURL              param.Field[string]                                                          `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefParam] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                           `json:"refresh_interval_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefParam struct {
	Name param.Field[string]                                                         `json:"name" api:"required"`
	Type param.Field[SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecGitRefParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3Param struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiGcsBucketMountSpecS3Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecParam struct {
	ID         param.Field[string]                                                            `json:"id" api:"required"`
	Git        param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitParam]        `json:"git" api:"required"`
	MountPath  param.Field[string]                                                            `json:"mount_path" api:"required"`
	Type       param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecType]            `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobParam]  `json:"azure_blob"`
	Cache      param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheParam]      `json:"cache"`
	Contexthub param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubParam] `json:"contexthub"`
	Encryption param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionParam] `json:"encryption"`
	Gcs        param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsParam]        `json:"gcs"`
	ReadOnly   param.Field[bool]                                                              `json:"read_only"`
	S3         param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3Param]         `json:"s3"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecParam) implementsSandboxMountConfigMountsUnionParam() {
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitParam struct {
	RemoteURL              param.Field[string]                                                        `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefParam] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                         `json:"refresh_interval_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefParam struct {
	Name param.Field[string]                                                       `json:"name" api:"required"`
	Type param.Field[SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecGitRefParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobParam struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecAzureBlobParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheParam struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecCacheParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubParam struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecContexthubParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionParam struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecEncryptionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsParam struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecGcsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3Param struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiGitRepoMountSpecS3Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecParam struct {
	ID         param.Field[string]                                                                   `json:"id" api:"required"`
	Contexthub param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubParam] `json:"contexthub" api:"required"`
	MountPath  param.Field[string]                                                                   `json:"mount_path" api:"required"`
	Type       param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecType]            `json:"type" api:"required"`
	AzureBlob  param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobParam]  `json:"azure_blob"`
	Cache      param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheParam]      `json:"cache"`
	Encryption param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionParam] `json:"encryption"`
	Gcs        param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsParam]        `json:"gcs"`
	Git        param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitParam]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                     `json:"read_only"`
	S3         param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3Param]         `json:"s3"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecParam) implementsSandboxMountConfigMountsUnionParam() {
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubParam struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecContexthubParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobParam struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecAzureBlobParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheParam struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecCacheParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionParam struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecEncryptionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsParam struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGcsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitParam struct {
	RemoteURL              param.Field[string]                                                               `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefParam] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                `json:"refresh_interval_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefParam struct {
	Name param.Field[string]                                                              `json:"name" api:"required"`
	Type param.Field[SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecGitRefParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3Param struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiContextHubRepoMountSpecS3Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecParam struct {
	ID         param.Field[string]                                                                       `json:"id" api:"required"`
	AzureBlob  param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobParam]  `json:"azure_blob" api:"required"`
	MountPath  param.Field[string]                                                                       `json:"mount_path" api:"required"`
	Type       param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecType]            `json:"type" api:"required"`
	Cache      param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheParam]      `json:"cache"`
	Contexthub param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubParam] `json:"contexthub"`
	Encryption param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionParam] `json:"encryption"`
	Gcs        param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsParam]        `json:"gcs"`
	Git        param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitParam]        `json:"git"`
	ReadOnly   param.Field[bool]                                                                         `json:"read_only"`
	S3         param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3Param]         `json:"s3"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecParam) implementsSandboxMountConfigMountsUnionParam() {
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobParam struct {
	Container param.Field[string] `json:"container" api:"required"`
	// StorageAccount is the storage account name: 3 to 24 lowercase letters and
	// digits.
	StorageAccount param.Field[string] `json:"storage_account" api:"required"`
	Prefix         param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecAzureBlobParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheParam struct {
	MaxSizeBytes     param.Field[int64] `json:"max_size_bytes"`
	WritebackSeconds param.Field[int64] `json:"writeback_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecCacheParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubParam struct {
	// Repo is the Context Hub repository to sync, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace). The repo's latest commit tree is mirrored
	// into the mount path.
	Repo param.Field[string] `json:"repo" api:"required"`
	// InitialPullOnly syncs the repo once at startup instead of polling for updates
	// for the sandbox's lifetime.
	InitialPullOnly param.Field[bool] `json:"initial_pull_only"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecContexthubParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionParam struct {
	// KMSKey is the key to encrypt uploads with. S3 mounts take a KMS key ARN in the
	// mount's region.
	KmsKey param.Field[string] `json:"kms_key" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecEncryptionParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsParam struct {
	Bucket param.Field[string] `json:"bucket" api:"required"`
	Prefix param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGcsParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitParam struct {
	RemoteURL              param.Field[string]                                                                   `json:"remote_url" api:"required"`
	Ref                    param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefParam] `json:"ref"`
	RefreshIntervalSeconds param.Field[int64]                                                                    `json:"refresh_interval_seconds"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefParam struct {
	Name param.Field[string]                                                                  `json:"name" api:"required"`
	Type param.Field[SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefType] `json:"type" api:"required"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecGitRefParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3Param struct {
	Bucket      param.Field[string] `json:"bucket" api:"required"`
	Region      param.Field[string] `json:"region" api:"required"`
	EndpointURL param.Field[string] `json:"endpoint_url"`
	PathStyle   param.Field[bool]   `json:"path_style"`
	Prefix      param.Field[string] `json:"prefix"`
}

func (r SandboxMountConfigMountsSandboxapiAzureBlobContainerMountSpecS3Param) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
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
	MountConfig            SandboxMountConfig              `json:"mount_config"`
	Name                   string                          `json:"name"`
	// Persistent reports whether the sandbox keeps its filesystem across stops. When
	// false, its writes are discarded whenever it stops and it is not charged for
	// storage.
	Persistent           bool        `json:"persistent"`
	PreserveMemoryOnStop bool        `json:"preserve_memory_on_stop"`
	ProxyConfig          ProxyConfig `json:"proxy_config"`
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
	// Volumes (beta) are the volumes attached to the sandbox, in slot order. Absent
	// when the workspace does not have sandbox volumes enabled.
	Volumes []SandboxResponseVolume `json:"volumes"`
	JSON    sandboxResponseJSON     `json:"-"`
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
	Persistent             apijson.Field
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
	Volumes                apijson.Field
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

type SandboxResponseVolume struct {
	MountPath  string                    `json:"mount_path"`
	ReadOnly   bool                      `json:"read_only"`
	VolumeID   string                    `json:"volume_id"`
	VolumeName string                    `json:"volume_name"`
	JSON       sandboxResponseVolumeJSON `json:"-"`
}

// sandboxResponseVolumeJSON contains the JSON metadata for the struct
// [SandboxResponseVolume]
type sandboxResponseVolumeJSON struct {
	MountPath   apijson.Field
	ReadOnly    apijson.Field
	VolumeID    apijson.Field
	VolumeName  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxResponseVolume) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxResponseVolumeJSON) RawJSON() string {
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
	// EmbedOrigins echoes the origins that may embed the URL, normalized; omitted when
	// there are none.
	EmbedOrigins []string `json:"embed_origins"`
	// EmbedURL opens the service in an iframe: the token rides the URL fragment, which
	// no server, log or Referer header sees, and the page it loads trades it for the
	// cookie. Append &redirect=/path to land elsewhere than /. Omitted in LangSmith
	// login mode.
	EmbedURL   string                 `json:"embed_url"`
	ExpiresAt  string                 `json:"expires_at"`
	ServiceURL string                 `json:"service_url"`
	JSON       serviceURLResponseJSON `json:"-"`
}

// serviceURLResponseJSON contains the JSON metadata for the struct
// [ServiceURLResponse]
type serviceURLResponseJSON struct {
	Token        apijson.Field
	Access       apijson.Field
	BrowserURL   apijson.Field
	EmbedOrigins apijson.Field
	EmbedURL     apijson.Field
	ExpiresAt    apijson.Field
	ServiceURL   apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
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
	Description     string `json:"description"`
	DockerImage     string `json:"docker_image"`
	FsCapacityBytes int64  `json:"fs_capacity_bytes"`
	FsUsedBytes     int64  `json:"fs_used_bytes"`
	ImageDigest     string `json:"image_digest"`
	// The settings this snapshot carries, in include terms: sandboxes created from it
	// start with each one their create request leaves unset. Env var values are never
	// returned.
	Included []SnapshotResponseIncluded `json:"included"`
	Labels   map[string]string          `json:"labels"`
	// MemorySnapshotSizeBytes is non-nil iff the snapshot was captured with VM memory
	// state. A non-nil value is the canonical signal that this snapshot can
	// warm-restore from memory; nil means rootfs only.
	MemorySnapshotSizeBytes int64 `json:"memory_snapshot_size_bytes"`
	// The mount_config sandboxes created from this snapshot start with when their
	// create request has none, with secret values redacted. Present only on snapshots
	// captured with include "mounts" or "all" from a sandbox that had mounts.
	MountConfig SandboxMountConfig `json:"mount_config"`
	Name        string             `json:"name"`
	// The proxy_config sandboxes created from this snapshot start with when their
	// create request has none, with secret values redacted. Present only on snapshots
	// captured with include "proxy_config" or "all".
	ProxyConfig ProxyConfig `json:"proxy_config"`
	RegistryID  string      `json:"registry_id"`
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
	Included                apijson.Field
	Labels                  apijson.Field
	MemorySnapshotSizeBytes apijson.Field
	MountConfig             apijson.Field
	Name                    apijson.Field
	ProxyConfig             apijson.Field
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

type SnapshotResponseIncluded string

const (
	SnapshotResponseIncludedMounts      SnapshotResponseIncluded = "mounts"
	SnapshotResponseIncludedProxyConfig SnapshotResponseIncluded = "proxy_config"
	SnapshotResponseIncludedEnvVars     SnapshotResponseIncluded = "env_vars"
	SnapshotResponseIncludedResources   SnapshotResponseIncluded = "resources"
	SnapshotResponseIncludedLifecycle   SnapshotResponseIncluded = "lifecycle"
	SnapshotResponseIncludedLabels      SnapshotResponseIncluded = "labels"
)

func (r SnapshotResponseIncluded) IsKnown() bool {
	switch r {
	case SnapshotResponseIncludedMounts, SnapshotResponseIncludedProxyConfig, SnapshotResponseIncludedEnvVars, SnapshotResponseIncludedResources, SnapshotResponseIncludedLifecycle, SnapshotResponseIncludedLabels:
		return true
	}
	return false
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
	// Recorded compute usage in Standard Units, as a decimal string, separate from
	// LCU. Null means complete Standard compute usage is unavailable; "0" means known
	// zero usage.
	ComputeStandardUnits string `json:"compute_standard_units" api:"nullable" format:"decimal"`
	// Memory usage in GiB-hours, as a decimal string. Available with RESOURCE
	// granularity; null means complete usage is unavailable.
	MemoryGibHours string `json:"memory_gib_hours" api:"nullable" format:"decimal"`
	// Storage usage in GiB-hours from hourly attributed-byte samples, including
	// checkpoints. Available with RESOURCE granularity; null means complete usage is
	// unavailable.
	StorageGibHours string `json:"storage_gib_hours" api:"nullable" format:"decimal"`
	// Allocated storage usage in Standard Units, as a decimal string, separate from
	// LSU. Null means complete Standard storage usage is unavailable; "0" means known
	// zero usage.
	StorageStandardUnits string `json:"storage_standard_units" api:"nullable" format:"decimal"`
	// CPU usage in vCPU-hours, as a decimal string. Available with RESOURCE
	// granularity; null means complete usage is unavailable.
	VcpuHours string                            `json:"vcpu_hours" api:"nullable" format:"decimal"`
	JSON      sandboxListUsageCostsResponseJSON `json:"-"`
}

// sandboxListUsageCostsResponseJSON contains the JSON metadata for the struct
// [SandboxListUsageCostsResponse]
type sandboxListUsageCostsResponseJSON struct {
	Lcu                  apijson.Field
	Lsu                  apijson.Field
	PeriodStart          apijson.Field
	ResourceID           apijson.Field
	ResourceType         apijson.Field
	ComputeStandardUnits apijson.Field
	MemoryGibHours       apijson.Field
	StorageGibHours      apijson.Field
	StorageStandardUnits apijson.Field
	VcpuHours            apijson.Field
	raw                  string
	ExtraFields          map[string]apijson.Field
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
	// Include bounded daily physical usage for the top 10 resources per metric and
	// Other. Requires RESOURCE and include_totals. Omitted when unrestricted access or
	// service support is unavailable. CLI: use --format raw to show this summary.
	IncludeDaily param.Field[bool] `query:"include_daily"`
	// Include usage totals across the filtered interval (RESOURCE only). Omitted when
	// unrestricted resource access cannot be verified or the service does not support
	// totals. CLI: use --format raw to show this summary.
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
