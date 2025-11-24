// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// ExampleBulkService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewExampleBulkService] method instead.
type ExampleBulkService struct {
	Options []option.RequestOption
}

// NewExampleBulkService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewExampleBulkService(opts ...option.RequestOption) (r *ExampleBulkService) {
	r = &ExampleBulkService{}
	r.Options = opts
	return
}

// Create bulk examples.
func (r *ExampleBulkService) New(ctx context.Context, body ExampleBulkNewParams, opts ...option.RequestOption) (res *[]Example, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples/bulk"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Legacy update examples in bulk. For update involving attachments, use PATCH
// /v1/platform/datasets/{dataset_id}/examples instead.
func (r *ExampleBulkService) PatchAll(ctx context.Context, body ExampleBulkPatchAllParams, opts ...option.RequestOption) (res *ExampleBulkPatchAllResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples/bulk"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

type ExampleBulkPatchAllResponse = interface{}

type ExampleBulkNewParams struct {
	// Schema for a batch of examples to be created.
	Body []ExampleBulkNewParamsBody `json:"body,required"`
}

func (r ExampleBulkNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.Body)
}

// Example with optional created_at to prevent duplicate versions in bulk
// operations.
type ExampleBulkNewParamsBody struct {
	DatasetID   param.Field[string]                             `json:"dataset_id,required" format:"uuid"`
	ID          param.Field[string]                             `json:"id" format:"uuid"`
	CreatedAt   param.Field[string]                             `json:"created_at"`
	Inputs      param.Field[interface{}]                        `json:"inputs"`
	Metadata    param.Field[interface{}]                        `json:"metadata"`
	Outputs     param.Field[interface{}]                        `json:"outputs"`
	SourceRunID param.Field[string]                             `json:"source_run_id" format:"uuid"`
	Split       param.Field[ExampleBulkNewParamsBodySplitUnion] `json:"split"`
	// Use Legacy Message Format for LLM runs
	UseLegacyMessageFormat  param.Field[bool]     `json:"use_legacy_message_format"`
	UseSourceRunAttachments param.Field[[]string] `json:"use_source_run_attachments"`
	UseSourceRunIo          param.Field[bool]     `json:"use_source_run_io"`
}

func (r ExampleBulkNewParamsBody) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [ExampleBulkNewParamsBodySplitArray], [shared.UnionString].
type ExampleBulkNewParamsBodySplitUnion interface {
	ImplementsExampleBulkNewParamsBodySplitUnion()
}

type ExampleBulkNewParamsBodySplitArray []string

func (r ExampleBulkNewParamsBodySplitArray) ImplementsExampleBulkNewParamsBodySplitUnion() {}

type ExampleBulkPatchAllParams struct {
	Body []ExampleBulkPatchAllParamsBody `json:"body,required"`
}

func (r ExampleBulkPatchAllParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.Body)
}

// Bulk update class for Example (includes example id).
type ExampleBulkPatchAllParamsBody struct {
	ID                    param.Field[string]                                  `json:"id,required" format:"uuid"`
	AttachmentsOperations param.Field[AttachmentsOperationsParam]              `json:"attachments_operations"`
	DatasetID             param.Field[string]                                  `json:"dataset_id" format:"uuid"`
	Inputs                param.Field[interface{}]                             `json:"inputs"`
	Metadata              param.Field[interface{}]                             `json:"metadata"`
	Outputs               param.Field[interface{}]                             `json:"outputs"`
	Overwrite             param.Field[bool]                                    `json:"overwrite"`
	Split                 param.Field[ExampleBulkPatchAllParamsBodySplitUnion] `json:"split"`
}

func (r ExampleBulkPatchAllParamsBody) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [ExampleBulkPatchAllParamsBodySplitArray], [shared.UnionString].
type ExampleBulkPatchAllParamsBodySplitUnion interface {
	ImplementsExampleBulkPatchAllParamsBodySplitUnion()
}

type ExampleBulkPatchAllParamsBodySplitArray []string

func (r ExampleBulkPatchAllParamsBodySplitArray) ImplementsExampleBulkPatchAllParamsBodySplitUnion() {
}
