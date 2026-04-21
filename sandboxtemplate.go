// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SandboxTemplateService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxTemplateService] method instead.
type SandboxTemplateService struct {
	Options []option.RequestOption
}

// NewSandboxTemplateService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxTemplateService(opts ...option.RequestOption) (r *SandboxTemplateService) {
	r = &SandboxTemplateService{}
	r.Options = opts
	return
}

// List sandbox templates for the authenticated tenant, with optional filtering,
// sorting, and pagination.
func (r *SandboxTemplateService) List(ctx context.Context, query SandboxTemplateListParams, opts ...option.RequestOption) (res *SandboxTemplateListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "v2/sandboxes/templates"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type SandboxTemplateListResponse struct {
	Offset    int64                                 `json:"offset"`
	Templates []SandboxTemplateListResponseTemplate `json:"templates"`
	JSON      sandboxTemplateListResponseJSON       `json:"-"`
}

// sandboxTemplateListResponseJSON contains the JSON metadata for the struct
// [SandboxTemplateListResponse]
type sandboxTemplateListResponseJSON struct {
	Offset      apijson.Field
	Templates   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxTemplateListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListResponseTemplate struct {
	ID        string `json:"id"`
	CreatedAt string `json:"created_at"`
	Image     string `json:"image"`
	Name      string `json:"name"`
	// Deprecated: template-level proxy config is legacy; configure proxy rules per box
	// instead.
	ProxyConfig  SandboxTemplateListResponseTemplatesProxyConfig   `json:"proxy_config"`
	RegistryID   string                                            `json:"registry_id"`
	RegistryName string                                            `json:"registry_name"`
	Resources    SandboxTemplateListResponseTemplatesResources     `json:"resources"`
	UpdatedAt    string                                            `json:"updated_at"`
	VolumeMounts []SandboxTemplateListResponseTemplatesVolumeMount `json:"volume_mounts"`
	JSON         sandboxTemplateListResponseTemplateJSON           `json:"-"`
}

// sandboxTemplateListResponseTemplateJSON contains the JSON metadata for the
// struct [SandboxTemplateListResponseTemplate]
type sandboxTemplateListResponseTemplateJSON struct {
	ID           apijson.Field
	CreatedAt    apijson.Field
	Image        apijson.Field
	Name         apijson.Field
	ProxyConfig  apijson.Field
	RegistryID   apijson.Field
	RegistryName apijson.Field
	Resources    apijson.Field
	UpdatedAt    apijson.Field
	VolumeMounts apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplate) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplateJSON) RawJSON() string {
	return r.raw
}

// Deprecated: template-level proxy config is legacy; configure proxy rules per box
// instead.
type SandboxTemplateListResponseTemplatesProxyConfig struct {
	AccessControl SandboxTemplateListResponseTemplatesProxyConfigAccessControl `json:"access_control"`
	NoProxy       []string                                                     `json:"no_proxy"`
	Rules         []SandboxTemplateListResponseTemplatesProxyConfigRule        `json:"rules"`
	JSON          sandboxTemplateListResponseTemplatesProxyConfigJSON          `json:"-"`
}

// sandboxTemplateListResponseTemplatesProxyConfigJSON contains the JSON metadata
// for the struct [SandboxTemplateListResponseTemplatesProxyConfig]
type sandboxTemplateListResponseTemplatesProxyConfigJSON struct {
	AccessControl apijson.Field
	NoProxy       apijson.Field
	Rules         apijson.Field
	raw           string
	ExtraFields   map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplatesProxyConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplatesProxyConfigJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListResponseTemplatesProxyConfigAccessControl struct {
	AllowList []string                                                         `json:"allow_list"`
	DenyList  []string                                                         `json:"deny_list"`
	JSON      sandboxTemplateListResponseTemplatesProxyConfigAccessControlJSON `json:"-"`
}

// sandboxTemplateListResponseTemplatesProxyConfigAccessControlJSON contains the
// JSON metadata for the struct
// [SandboxTemplateListResponseTemplatesProxyConfigAccessControl]
type sandboxTemplateListResponseTemplatesProxyConfigAccessControlJSON struct {
	AllowList   apijson.Field
	DenyList    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplatesProxyConfigAccessControl) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplatesProxyConfigAccessControlJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListResponseTemplatesProxyConfigRule struct {
	MatchHosts []string                                                     `json:"match_hosts" api:"required"`
	Name       string                                                       `json:"name" api:"required"`
	Enabled    bool                                                         `json:"enabled"`
	Headers    []SandboxTemplateListResponseTemplatesProxyConfigRulesHeader `json:"headers"`
	MatchPaths []string                                                     `json:"match_paths"`
	JSON       sandboxTemplateListResponseTemplatesProxyConfigRuleJSON      `json:"-"`
}

// sandboxTemplateListResponseTemplatesProxyConfigRuleJSON contains the JSON
// metadata for the struct [SandboxTemplateListResponseTemplatesProxyConfigRule]
type sandboxTemplateListResponseTemplatesProxyConfigRuleJSON struct {
	MatchHosts  apijson.Field
	Name        apijson.Field
	Enabled     apijson.Field
	Headers     apijson.Field
	MatchPaths  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplatesProxyConfigRule) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplatesProxyConfigRuleJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListResponseTemplatesProxyConfigRulesHeader struct {
	Name  string                                                          `json:"name" api:"required"`
	Type  SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersType `json:"type" api:"required"`
	IsSet bool                                                            `json:"is_set"`
	Value string                                                          `json:"value"`
	JSON  sandboxTemplateListResponseTemplatesProxyConfigRulesHeaderJSON  `json:"-"`
}

// sandboxTemplateListResponseTemplatesProxyConfigRulesHeaderJSON contains the JSON
// metadata for the struct
// [SandboxTemplateListResponseTemplatesProxyConfigRulesHeader]
type sandboxTemplateListResponseTemplatesProxyConfigRulesHeaderJSON struct {
	Name        apijson.Field
	Type        apijson.Field
	IsSet       apijson.Field
	Value       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplatesProxyConfigRulesHeader) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplatesProxyConfigRulesHeaderJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersType string

const (
	SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersTypePlaintext       SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersType = "plaintext"
	SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersTypeOpaque          SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersType = "opaque"
	SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersTypeWorkspaceSecret SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersType = "workspace_secret"
)

func (r SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersType) IsKnown() bool {
	switch r {
	case SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersTypePlaintext, SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersTypeOpaque, SandboxTemplateListResponseTemplatesProxyConfigRulesHeadersTypeWorkspaceSecret:
		return true
	}
	return false
}

type SandboxTemplateListResponseTemplatesResources struct {
	CPU     string                                            `json:"cpu"`
	Memory  string                                            `json:"memory"`
	Storage string                                            `json:"storage"`
	JSON    sandboxTemplateListResponseTemplatesResourcesJSON `json:"-"`
}

// sandboxTemplateListResponseTemplatesResourcesJSON contains the JSON metadata for
// the struct [SandboxTemplateListResponseTemplatesResources]
type sandboxTemplateListResponseTemplatesResourcesJSON struct {
	CPU         apijson.Field
	Memory      apijson.Field
	Storage     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplatesResources) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplatesResourcesJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListResponseTemplatesVolumeMount struct {
	MountPath  string                                              `json:"mount_path" api:"required"`
	VolumeName string                                              `json:"volume_name" api:"required"`
	JSON       sandboxTemplateListResponseTemplatesVolumeMountJSON `json:"-"`
}

// sandboxTemplateListResponseTemplatesVolumeMountJSON contains the JSON metadata
// for the struct [SandboxTemplateListResponseTemplatesVolumeMount]
type sandboxTemplateListResponseTemplatesVolumeMountJSON struct {
	MountPath   apijson.Field
	VolumeName  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxTemplateListResponseTemplatesVolumeMount) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxTemplateListResponseTemplatesVolumeMountJSON) RawJSON() string {
	return r.raw
}

type SandboxTemplateListParams struct {
	// Maximum number of results
	Limit param.Field[int64] `query:"limit"`
	// Filter by name substring
	NameContains param.Field[string] `query:"name_contains"`
	// Pagination offset
	Offset param.Field[int64] `query:"offset"`
	// Sort column (name, image, created_at)
	SortBy param.Field[string] `query:"sort_by"`
	// Sort direction (asc, desc)
	SortDirection param.Field[string] `query:"sort_direction"`
}

// URLQuery serializes [SandboxTemplateListParams]'s query parameters as
// `url.Values`.
func (r SandboxTemplateListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
