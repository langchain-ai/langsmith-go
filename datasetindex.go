// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetIndexService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetIndexService] method instead.
type DatasetIndexService struct {
	Options []option.RequestOption
}

// NewDatasetIndexService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetIndexService(opts ...option.RequestOption) (r *DatasetIndexService) {
	r = &DatasetIndexService{}
	r.Options = opts
	return
}

// Index a dataset.
func (r *DatasetIndexService) New(ctx context.Context, datasetID string, body DatasetIndexNewParams, opts ...option.RequestOption) (res *DatasetIndexNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/index", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get index info.
func (r *DatasetIndexService) Get(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *DatasetIndexGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/index", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Remove an index for a dataset.
func (r *DatasetIndexService) DeleteAll(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *DatasetIndexDeleteAllResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/index", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Sync an index for a dataset.
func (r *DatasetIndexService) Sync(ctx context.Context, datasetID string, opts ...option.RequestOption) (res *DatasetIndexSyncResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/index/sync", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

type DatasetIndexNewResponse = interface{}

// Dataset schema for serving.
type DatasetIndexGetResponse struct {
	DatasetID          string                      `json:"dataset_id,required" format:"uuid"`
	LastUpdatedVersion time.Time                   `json:"last_updated_version,nullable" format:"date-time"`
	Tag                string                      `json:"tag,nullable"`
	JSON               datasetIndexGetResponseJSON `json:"-"`
}

// datasetIndexGetResponseJSON contains the JSON metadata for the struct
// [DatasetIndexGetResponse]
type datasetIndexGetResponseJSON struct {
	DatasetID          apijson.Field
	LastUpdatedVersion apijson.Field
	Tag                apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *DatasetIndexGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetIndexGetResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetIndexDeleteAllResponse = interface{}

type DatasetIndexSyncResponse = interface{}

type DatasetIndexNewParams struct {
	Tag param.Field[string] `json:"tag"`
}

func (r DatasetIndexNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
