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

// AgentService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAgentService] method instead.
type AgentService struct {
	Options []option.RequestOption
}

// NewAgentService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewAgentService(opts ...option.RequestOption) (r *AgentService) {
	r = &AgentService{}
	r.Options = opts
	return
}

// Creates an Agent with a caller-assigned ID. Returns 409 if the ID already
// exists. Creates all four environments unless production_only is set for a FLEET
// agent. source_type defaults to EXTERNAL in this SDK. agent_type defaults to
// CUSTOM. The authenticated user becomes the owner
// when available. Requires permission to create projects. The response does not
// include runtime bindings.
func (r *AgentService) New(ctx context.Context, body AgentNewParams, opts ...option.RequestOption) (res *Agent, err error) {
	if !body.SourceType.Present {
		body.SourceType = F(AgentNewParamsSourceTypeExternal)
	}
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/agents"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns an Agent with the environments visible to the caller. The path value is
// the user-assigned Agent ID. Pass include=runtime_bindings to embed each
// environment's runtime bindings.
func (r *AgentService) Get(ctx context.Context, agentID string, query AgentGetParams, opts ...option.RequestOption) (res *Agent, err error) {
	opts = slices.Concat(r.Options, opts)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/agents/%s", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Updates an Agent's name, source_type, agent_type, or owner_id. Omitted fields
// keep their values. Changing or clearing owner_id additionally requires
// workspaces:manage; assigning an owner also requires workspaces:read and active
// workspace membership. Fleet ownership cannot be changed through this endpoint.
// Changing agent_type requires workspaces:manage in addition to project update
// access; repeating the current type is a no-op. Source type can change only for
// an EXTERNAL Agent without runtime bindings. The id is immutable and is rejected
// when present. Requires update permission on every environment project.
func (r *AgentService) Update(ctx context.Context, agentID string, body AgentUpdateParams, opts ...option.RequestOption) (res *Agent, err error) {
	opts = slices.Concat(r.Options, opts)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/agents/%s", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Returns Agents visible to the caller. Each item is the same Agent shape as GET
// /api/v1/agents/{agent_id}, including the environments the caller may read. The
// list never embeds runtime bindings. An unknown query parameter, including limit,
// returns 400; use page_size to set the page size.
func (r *AgentService) List(ctx context.Context, query AgentListParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[Agent], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/agents"
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

// Returns Agents visible to the caller. Each item is the same Agent shape as GET
// /api/v1/agents/{agent_id}, including the environments the caller may read. The
// list never embeds runtime bindings. An unknown query parameter, including limit,
// returns 400; use page_size to set the page size.
func (r *AgentService) ListAutoPaging(ctx context.Context, query AgentListParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[Agent] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.List(ctx, query, opts...))
}

// Permanently deletes an Agent, along with its environments, runtime bindings,
// runtimes, and tracing projects. This cannot be undone. Deleting an Agent that
// does not exist returns 204.
func (r *AgentService) Delete(ctx context.Context, agentID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if agentID == "" {
		err = errors.New("missing required agent_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/agents/%s", agentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type Agent struct {
	ID           string             `json:"id" api:"required"`
	AgentType    AgentAgentType     `json:"agent_type" api:"required"`
	CreatedAt    time.Time          `json:"created_at" api:"required" format:"date-time"`
	Environments []AgentEnvironment `json:"environments" api:"required"`
	Name         string             `json:"name" api:"required"`
	SourceType   AgentSourceType    `json:"source_type" api:"required"`
	// AvatarURL is a short-lived URL for the agent avatar, when present.
	AvatarURL string `json:"avatar_url" api:"nullable"`
	// Owner is the owner's full name, falling back to email when available.
	Owner      string    `json:"owner"`
	OwnerID    string    `json:"owner_id" api:"nullable" format:"uuid"`
	TagKeyID   string    `json:"tag_key_id" api:"nullable" format:"uuid"`
	TagValueID string    `json:"tag_value_id" api:"nullable" format:"uuid"`
	JSON       agentJSON `json:"-"`
}

// agentJSON contains the JSON metadata for the struct [Agent]
type agentJSON struct {
	ID           apijson.Field
	AgentType    apijson.Field
	CreatedAt    apijson.Field
	Environments apijson.Field
	Name         apijson.Field
	SourceType   apijson.Field
	AvatarURL    apijson.Field
	Owner        apijson.Field
	OwnerID      apijson.Field
	TagKeyID     apijson.Field
	TagValueID   apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *Agent) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentJSON) RawJSON() string {
	return r.raw
}

type AgentAgentType string

const (
	AgentAgentTypeCustom AgentAgentType = "CUSTOM"
	AgentAgentTypeCoding AgentAgentType = "CODING"
)

func (r AgentAgentType) IsKnown() bool {
	switch r {
	case AgentAgentTypeCustom, AgentAgentTypeCoding:
		return true
	}
	return false
}

type AgentEnvironment struct {
	ID                string                            `json:"id" api:"required" format:"uuid"`
	AgentID           string                            `json:"agent_id" api:"required"`
	Environment       AgentEnvironmentsEnvironment      `json:"environment" api:"required"`
	TracerSessionID   string                            `json:"tracer_session_id" api:"required" format:"uuid"`
	TracerSessionName string                            `json:"tracer_session_name" api:"required"`
	RuntimeBindings   []AgentEnvironmentsRuntimeBinding `json:"runtime_bindings"`
	JSON              agentEnvironmentJSON              `json:"-"`
}

// agentEnvironmentJSON contains the JSON metadata for the struct
// [AgentEnvironment]
type agentEnvironmentJSON struct {
	ID                apijson.Field
	AgentID           apijson.Field
	Environment       apijson.Field
	TracerSessionID   apijson.Field
	TracerSessionName apijson.Field
	RuntimeBindings   apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *AgentEnvironment) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentEnvironmentJSON) RawJSON() string {
	return r.raw
}

type AgentEnvironmentsEnvironment string

const (
	AgentEnvironmentsEnvironmentLocal       AgentEnvironmentsEnvironment = "LOCAL"
	AgentEnvironmentsEnvironmentDevelopment AgentEnvironmentsEnvironment = "DEVELOPMENT"
	AgentEnvironmentsEnvironmentStaging     AgentEnvironmentsEnvironment = "STAGING"
	AgentEnvironmentsEnvironmentProduction  AgentEnvironmentsEnvironment = "PRODUCTION"
)

func (r AgentEnvironmentsEnvironment) IsKnown() bool {
	switch r {
	case AgentEnvironmentsEnvironmentLocal, AgentEnvironmentsEnvironmentDevelopment, AgentEnvironmentsEnvironmentStaging, AgentEnvironmentsEnvironmentProduction:
		return true
	}
	return false
}

type AgentEnvironmentsRuntimeBinding struct {
	ID                 string                              `json:"id" api:"required" format:"uuid"`
	AgentEnvironmentID string                              `json:"agent_environment_id" api:"required" format:"uuid"`
	SourceID           string                              `json:"source_id" api:"required"`
	JSON               agentEnvironmentsRuntimeBindingJSON `json:"-"`
}

// agentEnvironmentsRuntimeBindingJSON contains the JSON metadata for the struct
// [AgentEnvironmentsRuntimeBinding]
type agentEnvironmentsRuntimeBindingJSON struct {
	ID                 apijson.Field
	AgentEnvironmentID apijson.Field
	SourceID           apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *AgentEnvironmentsRuntimeBinding) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r agentEnvironmentsRuntimeBindingJSON) RawJSON() string {
	return r.raw
}

type AgentSourceType string

const (
	AgentSourceTypeLsd      AgentSourceType = "LSD"
	AgentSourceTypeMda      AgentSourceType = "MDA"
	AgentSourceTypeExternal AgentSourceType = "EXTERNAL"
	AgentSourceTypeFleet    AgentSourceType = "FLEET"
)

func (r AgentSourceType) IsKnown() bool {
	switch r {
	case AgentSourceTypeLsd, AgentSourceTypeMda, AgentSourceTypeExternal, AgentSourceTypeFleet:
		return true
	}
	return false
}

type AgentNewParams struct {
	// Unique agent identifier: 1-63 lowercase letters, digits, or hyphens, starting
	// with a letter and ending with a letter or digit.
	ID param.Field[string] `json:"id" api:"required"`
	// Agent runtime kind. Use EXTERNAL for an application hosted outside LangSmith.
	// Defaults to EXTERNAL when omitted.
	SourceType param.Field[AgentNewParamsSourceType] `json:"source_type"`
	// AgentType defaults to CUSTOM when omitted.
	AgentType param.Field[AgentNewParamsAgentType] `json:"agent_type"`
	// Display name. Defaults to the agent ID.
	Name param.Field[string] `json:"name"`
	// ProductionOnly creates only the PRODUCTION environment. It is accepted only when
	// SourceType is FLEET; other sources always get all four environments.
	ProductionOnly param.Field[bool] `json:"production_only"`
}

func (r AgentNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Agent runtime kind. Use EXTERNAL for an application hosted outside LangSmith.
type AgentNewParamsSourceType string

const (
	AgentNewParamsSourceTypeLsd      AgentNewParamsSourceType = "LSD"
	AgentNewParamsSourceTypeMda      AgentNewParamsSourceType = "MDA"
	AgentNewParamsSourceTypeExternal AgentNewParamsSourceType = "EXTERNAL"
	AgentNewParamsSourceTypeFleet    AgentNewParamsSourceType = "FLEET"
)

func (r AgentNewParamsSourceType) IsKnown() bool {
	switch r {
	case AgentNewParamsSourceTypeLsd, AgentNewParamsSourceTypeMda, AgentNewParamsSourceTypeExternal, AgentNewParamsSourceTypeFleet:
		return true
	}
	return false
}

// AgentType defaults to CUSTOM when omitted.
type AgentNewParamsAgentType string

const (
	AgentNewParamsAgentTypeCustom AgentNewParamsAgentType = "CUSTOM"
	AgentNewParamsAgentTypeCoding AgentNewParamsAgentType = "CODING"
)

func (r AgentNewParamsAgentType) IsKnown() bool {
	switch r {
	case AgentNewParamsAgentTypeCustom, AgentNewParamsAgentTypeCoding:
		return true
	}
	return false
}

type AgentGetParams struct {
	// Relation to embed. runtime_bindings adds each environment's runtime bindings; an
	// empty array means no runtime is linked.
	Include param.Field[AgentGetParamsInclude] `query:"include"`
}

// URLQuery serializes [AgentGetParams]'s query parameters as `url.Values`.
func (r AgentGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Relation to embed. runtime_bindings adds each environment's runtime bindings; an
// empty array means no runtime is linked.
type AgentGetParamsInclude string

const (
	AgentGetParamsIncludeRuntimeBindings AgentGetParamsInclude = "runtime_bindings"
)

func (r AgentGetParamsInclude) IsKnown() bool {
	switch r {
	case AgentGetParamsIncludeRuntimeBindings:
		return true
	}
	return false
}

type AgentUpdateParams struct {
	// New classification. Omit to keep the current classification.
	AgentType param.Field[AgentUpdateParamsAgentType] `json:"agent_type"`
	// New display name. Omit to keep the current name.
	Name param.Field[string] `json:"name"`
	// Omit to leave the owner unchanged; null clears it.
	OwnerID param.Field[string] `json:"owner_id" format:"uuid"`
	// New runtime kind. Only EXTERNAL agents without runtime bindings can change kind.
	SourceType param.Field[AgentUpdateParamsSourceType] `json:"source_type"`
}

func (r AgentUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// New classification. Omit to keep the current classification.
type AgentUpdateParamsAgentType string

const (
	AgentUpdateParamsAgentTypeCustom AgentUpdateParamsAgentType = "CUSTOM"
	AgentUpdateParamsAgentTypeCoding AgentUpdateParamsAgentType = "CODING"
)

func (r AgentUpdateParamsAgentType) IsKnown() bool {
	switch r {
	case AgentUpdateParamsAgentTypeCustom, AgentUpdateParamsAgentTypeCoding:
		return true
	}
	return false
}

// New runtime kind. Only EXTERNAL agents without runtime bindings can change kind.
type AgentUpdateParamsSourceType string

const (
	AgentUpdateParamsSourceTypeLsd      AgentUpdateParamsSourceType = "LSD"
	AgentUpdateParamsSourceTypeMda      AgentUpdateParamsSourceType = "MDA"
	AgentUpdateParamsSourceTypeExternal AgentUpdateParamsSourceType = "EXTERNAL"
	AgentUpdateParamsSourceTypeFleet    AgentUpdateParamsSourceType = "FLEET"
)

func (r AgentUpdateParamsSourceType) IsKnown() bool {
	switch r {
	case AgentUpdateParamsSourceTypeLsd, AgentUpdateParamsSourceTypeMda, AgentUpdateParamsSourceTypeExternal, AgentUpdateParamsSourceTypeFleet:
		return true
	}
	return false
}

type AgentListParams struct {
	// Exact Agent ids (repeat the parameter; maximum 100)
	ID param.Field[[]string] `query:"id"`
	// Opaque pagination cursor
	Cursor param.Field[string] `query:"cursor"`
	// Case-insensitive partial Agent name
	Name param.Field[string] `query:"name"`
	// Maximum number of Agents to return
	PageSize param.Field[int64] `query:"page_size"`
	// Field to sort by. When name is set, exact name matches still come first.
	SortBy param.Field[AgentListParamsSortBy] `query:"sort_by"`
	// Sort direction
	SortOrder param.Field[AgentListParamsSortOrder] `query:"sort_order"`
	// Agent source type
	SourceType param.Field[AgentListParamsSourceType] `query:"source_type"`
}

// URLQuery serializes [AgentListParams]'s query parameters as `url.Values`.
func (r AgentListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Field to sort by. When name is set, exact name matches still come first.
type AgentListParamsSortBy string

const (
	AgentListParamsSortByLastRunStartTime AgentListParamsSortBy = "last_run_start_time"
	AgentListParamsSortByName             AgentListParamsSortBy = "name"
	AgentListParamsSortByCreatedAt        AgentListParamsSortBy = "created_at"
)

func (r AgentListParamsSortBy) IsKnown() bool {
	switch r {
	case AgentListParamsSortByLastRunStartTime, AgentListParamsSortByName, AgentListParamsSortByCreatedAt:
		return true
	}
	return false
}

// Sort direction
type AgentListParamsSortOrder string

const (
	AgentListParamsSortOrderAsc  AgentListParamsSortOrder = "asc"
	AgentListParamsSortOrderDesc AgentListParamsSortOrder = "desc"
)

func (r AgentListParamsSortOrder) IsKnown() bool {
	switch r {
	case AgentListParamsSortOrderAsc, AgentListParamsSortOrderDesc:
		return true
	}
	return false
}

// Agent source type
type AgentListParamsSourceType string

const (
	AgentListParamsSourceTypeLsd      AgentListParamsSourceType = "LSD"
	AgentListParamsSourceTypeMda      AgentListParamsSourceType = "MDA"
	AgentListParamsSourceTypeExternal AgentListParamsSourceType = "EXTERNAL"
	AgentListParamsSourceTypeFleet    AgentListParamsSourceType = "FLEET"
)

func (r AgentListParamsSourceType) IsKnown() bool {
	switch r {
	case AgentListParamsSourceTypeLsd, AgentListParamsSourceTypeMda, AgentListParamsSourceTypeExternal, AgentListParamsSourceTypeFleet:
		return true
	}
	return false
}
