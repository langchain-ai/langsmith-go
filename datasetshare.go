// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/stainless-sdks/langsmith-api-go/internal/apijson"
	"github.com/stainless-sdks/langsmith-api-go/internal/apiquery"
	"github.com/stainless-sdks/langsmith-api-go/internal/param"
	"github.com/stainless-sdks/langsmith-api-go/internal/requestconfig"
	"github.com/stainless-sdks/langsmith-api-go/option"
)

// DatasetShareService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetShareService] method instead.
type DatasetShareService struct {
	Options []option.RequestOption
}

// NewDatasetShareService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetShareService(opts ...option.RequestOption) (r *DatasetShareService) {
	r = &DatasetShareService{}
	r.Options = opts
	return
}

// Share a dataset.
func (r *DatasetShareService) New(ctx context.Context, datasetID string, body DatasetShareNewParams, opts ...option.RequestOption) (res *DatasetShareSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/share", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, &res, opts...)
	return
}

// Get the state of sharing a dataset
func (r *DatasetShareService) Get(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *DatasetShareSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/share", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Unshare a dataset.
func (r *DatasetShareService) DeleteAll(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *DatasetShareDeleteAllResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/share", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

type DatasetShareSchema struct {
	DatasetID  string                 `json:"dataset_id,required" format:"uuid"`
	ShareToken string                 `json:"share_token,required" format:"uuid"`
	JSON       datasetShareSchemaJSON `json:"-"`
}

// datasetShareSchemaJSON contains the JSON metadata for the struct
// [DatasetShareSchema]
type datasetShareSchemaJSON struct {
	DatasetID   apijson.Field
	ShareToken  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetShareSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetShareSchemaJSON) RawJSON() string {
	return r.raw
}

type DatasetShareDeleteAllResponse = interface{}

type DatasetShareNewParams struct {
	ShareProjects param.Field[bool] `query:"share_projects"`
}

// URLQuery serializes [DatasetShareNewParams]'s query parameters as `url.Values`.
func (r DatasetShareNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
