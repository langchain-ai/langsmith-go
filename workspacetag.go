// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// WorkspaceTagService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewWorkspaceTagService] method instead.
type WorkspaceTagService struct {
	Options []option.RequestOption
	Keys    *WorkspaceTagKeyService
}

// NewWorkspaceTagService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewWorkspaceTagService(opts ...option.RequestOption) (r *WorkspaceTagService) {
	r = &WorkspaceTagService{}
	r.Options = opts
	r.Keys = NewWorkspaceTagKeyService(opts...)
	return
}

// List the workspace's tag keys, each with its values.
//
// With resource_type, only values applied to at least one resource of that type
// are returned, and keys with no such values are left out.
func (r *WorkspaceTagService) List(ctx context.Context, query WorkspaceTagListParams, opts ...option.RequestOption) (res *[]TagKeyWithValues, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/workspaces/current/tags"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type TagKeyWithValues struct {
	ID          string               `json:"id" api:"required" format:"uuid"`
	CreatedAt   time.Time            `json:"created_at" api:"required" format:"date-time"`
	Key         string               `json:"key" api:"required"`
	UpdatedAt   time.Time            `json:"updated_at" api:"required" format:"date-time"`
	Values      []TagValue           `json:"values" api:"required"`
	Description string               `json:"description" api:"nullable"`
	JSON        tagKeyWithValuesJSON `json:"-"`
}

// tagKeyWithValuesJSON contains the JSON metadata for the struct
// [TagKeyWithValues]
type tagKeyWithValuesJSON struct {
	ID          apijson.Field
	CreatedAt   apijson.Field
	Key         apijson.Field
	UpdatedAt   apijson.Field
	Values      apijson.Field
	Description apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *TagKeyWithValues) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r tagKeyWithValuesJSON) RawJSON() string {
	return r.raw
}

type WorkspaceTagListParams struct {
	ResourceType param.Field[WorkspaceTagListParamsResourceType] `query:"resource_type"`
}

// URLQuery serializes [WorkspaceTagListParams]'s query parameters as `url.Values`.
func (r WorkspaceTagListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type WorkspaceTagListParamsResourceType string

const (
	WorkspaceTagListParamsResourceTypeAgent            WorkspaceTagListParamsResourceType = "agent"
	WorkspaceTagListParamsResourceTypeCustomApp        WorkspaceTagListParamsResourceType = "custom_app"
	WorkspaceTagListParamsResourceTypeDashboard        WorkspaceTagListParamsResourceType = "dashboard"
	WorkspaceTagListParamsResourceTypeDataset          WorkspaceTagListParamsResourceType = "dataset"
	WorkspaceTagListParamsResourceTypeDeployment       WorkspaceTagListParamsResourceType = "deployment"
	WorkspaceTagListParamsResourceTypeEvaluator        WorkspaceTagListParamsResourceType = "evaluator"
	WorkspaceTagListParamsResourceTypeExperiment       WorkspaceTagListParamsResourceType = "experiment"
	WorkspaceTagListParamsResourceTypeFleetIntegration WorkspaceTagListParamsResourceType = "fleet_integration"
	WorkspaceTagListParamsResourceTypeMcpServer        WorkspaceTagListParamsResourceType = "mcp_server"
	WorkspaceTagListParamsResourceTypeProject          WorkspaceTagListParamsResourceType = "project"
	WorkspaceTagListParamsResourceTypePrompt           WorkspaceTagListParamsResourceType = "prompt"
	WorkspaceTagListParamsResourceTypeQueue            WorkspaceTagListParamsResourceType = "queue"
	WorkspaceTagListParamsResourceTypeSandbox          WorkspaceTagListParamsResourceType = "sandbox"
	WorkspaceTagListParamsResourceTypeSkill            WorkspaceTagListParamsResourceType = "skill"
	WorkspaceTagListParamsResourceTypeSnapshotName     WorkspaceTagListParamsResourceType = "snapshot_name"
)

func (r WorkspaceTagListParamsResourceType) IsKnown() bool {
	switch r {
	case WorkspaceTagListParamsResourceTypeAgent, WorkspaceTagListParamsResourceTypeCustomApp, WorkspaceTagListParamsResourceTypeDashboard, WorkspaceTagListParamsResourceTypeDataset, WorkspaceTagListParamsResourceTypeDeployment, WorkspaceTagListParamsResourceTypeEvaluator, WorkspaceTagListParamsResourceTypeExperiment, WorkspaceTagListParamsResourceTypeFleetIntegration, WorkspaceTagListParamsResourceTypeMcpServer, WorkspaceTagListParamsResourceTypeProject, WorkspaceTagListParamsResourceTypePrompt, WorkspaceTagListParamsResourceTypeQueue, WorkspaceTagListParamsResourceTypeSandbox, WorkspaceTagListParamsResourceTypeSkill, WorkspaceTagListParamsResourceTypeSnapshotName:
		return true
	}
	return false
}
