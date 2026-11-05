// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetExampleService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetExampleService] method instead.
type DatasetExampleService struct {
	Options []option.RequestOption
}

// NewDatasetExampleService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetExampleService(opts ...option.RequestOption) (r *DatasetExampleService) {
	r = &DatasetExampleService{}
	r.Options = opts
	return
}

// Soft-delete an example, preserving prior versions and their attachments. If the
// latest version is already deleted, the request succeeds without creating another
// version. Deletion is recorded at the current time or just after the latest
// version, whichever is later. For future-dated versions, latest reads reflect
// deletion immediately; timestamp reads reflect deletion only at or after the
// recorded deletion timestamp.
func (r *DatasetExampleService) Delete(ctx context.Context, datasetID string, exampleID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return err
	}
	if exampleID == "" {
		err = errors.New("missing required example_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/platform/datasets/%s/examples/%s", datasetID, exampleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Applies the same metadata changes, in order, to every listed example in a
// dataset. Examples that are not in the dataset or were deleted are skipped; the
// rest are updated together, so if the update fails none of them change. Answers
// 207 with one status per submitted example ID: UPDATED, or NOT_FOUND for a
// skipped example. Metadata keys the operations do not name keep their values. To
// add examples to a split or remove them from one, use ARRAY_ADD or ARRAY_REMOVE
// on the dataset_split key. Requires permission to update datasets.
func (r *DatasetExampleService) UpdateMetadata(ctx context.Context, datasetID string, body DatasetExampleUpdateMetadataParams, opts ...option.RequestOption) (res *DatasetExampleUpdateMetadataResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/platform/datasets/%s/examples/metadata", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

type DatasetExampleUpdateMetadataResponse struct {
	// One entry per submitted example ID, in request order.
	Statuses []DatasetExampleUpdateMetadataResponseStatus `json:"statuses" api:"required"`
	JSON     datasetExampleUpdateMetadataResponseJSON     `json:"-"`
}

// datasetExampleUpdateMetadataResponseJSON contains the JSON metadata for the
// struct [DatasetExampleUpdateMetadataResponse]
type datasetExampleUpdateMetadataResponseJSON struct {
	Statuses    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetExampleUpdateMetadataResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetExampleUpdateMetadataResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetExampleUpdateMetadataResponseStatus struct {
	// The example ID, as submitted.
	ID string `json:"id" api:"required" format:"uuid"`
	// UPDATED if the example's metadata was changed; NOT_FOUND if the example is not
	// in the dataset or was deleted.
	Status DatasetExampleUpdateMetadataResponseStatusesStatus `json:"status" api:"required"`
	// Why the example was not updated. Omitted when it was.
	Description string                                         `json:"description"`
	JSON        datasetExampleUpdateMetadataResponseStatusJSON `json:"-"`
}

// datasetExampleUpdateMetadataResponseStatusJSON contains the JSON metadata for
// the struct [DatasetExampleUpdateMetadataResponseStatus]
type datasetExampleUpdateMetadataResponseStatusJSON struct {
	ID          apijson.Field
	Status      apijson.Field
	Description apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *DatasetExampleUpdateMetadataResponseStatus) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetExampleUpdateMetadataResponseStatusJSON) RawJSON() string {
	return r.raw
}

// UPDATED if the example's metadata was changed; NOT_FOUND if the example is not
// in the dataset or was deleted.
type DatasetExampleUpdateMetadataResponseStatusesStatus string

const (
	DatasetExampleUpdateMetadataResponseStatusesStatusUpdated  DatasetExampleUpdateMetadataResponseStatusesStatus = "UPDATED"
	DatasetExampleUpdateMetadataResponseStatusesStatusNotFound DatasetExampleUpdateMetadataResponseStatusesStatus = "NOT_FOUND"
)

func (r DatasetExampleUpdateMetadataResponseStatusesStatus) IsKnown() bool {
	switch r {
	case DatasetExampleUpdateMetadataResponseStatusesStatusUpdated, DatasetExampleUpdateMetadataResponseStatusesStatusNotFound:
		return true
	}
	return false
}

type DatasetExampleUpdateMetadataParams struct {
	// The IDs of the examples to change; at most 10000. IDs of examples that are not
	// in the dataset, or were deleted, are skipped and reported as NOT_FOUND.
	ExampleIDs param.Field[[]string] `json:"example_ids" api:"required" format:"uuid"`
	// The metadata changes to apply, in order; at most 50, of which at most 10 are
	// array operations.
	Operations param.Field[[]DatasetExampleUpdateMetadataParamsOperation] `json:"operations" api:"required"`
}

func (r DatasetExampleUpdateMetadataParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetExampleUpdateMetadataParamsOperation struct {
	// The top-level metadata key to change, at most 255 bytes. Use dataset_split with
	// ARRAY_ADD or ARRAY_REMOVE to change which splits the examples belong to.
	Key param.Field[string] `json:"key" api:"required"`
	// The change to make: SET replaces the value at key; ARRAY_ADD adds values to the
	// list at key; ARRAY_REMOVE removes values from it. For array operations, a
	// missing or non-list value starts as an empty list, except for dataset_split,
	// where it starts as the base split.
	Operation param.Field[DatasetExampleUpdateMetadataParamsOperationsOperation] `json:"operation" api:"required"`
	// The new value for a SET operation; any JSON value, at most 64 KB. Ignored by
	// array operations.
	Value param.Field[interface{}] `json:"value"`
	// The values to add or remove for an array operation; at most 1000 values of at
	// most 255 bytes each. Ignored by SET.
	Values param.Field[[]string] `json:"values"`
}

func (r DatasetExampleUpdateMetadataParamsOperation) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// The change to make: SET replaces the value at key; ARRAY_ADD adds values to the
// list at key; ARRAY_REMOVE removes values from it. For array operations, a
// missing or non-list value starts as an empty list, except for dataset_split,
// where it starts as the base split.
type DatasetExampleUpdateMetadataParamsOperationsOperation string

const (
	DatasetExampleUpdateMetadataParamsOperationsOperationSet         DatasetExampleUpdateMetadataParamsOperationsOperation = "SET"
	DatasetExampleUpdateMetadataParamsOperationsOperationArrayAdd    DatasetExampleUpdateMetadataParamsOperationsOperation = "ARRAY_ADD"
	DatasetExampleUpdateMetadataParamsOperationsOperationArrayRemove DatasetExampleUpdateMetadataParamsOperationsOperation = "ARRAY_REMOVE"
)

func (r DatasetExampleUpdateMetadataParamsOperationsOperation) IsKnown() bool {
	switch r {
	case DatasetExampleUpdateMetadataParamsOperationsOperationSet, DatasetExampleUpdateMetadataParamsOperationsOperationArrayAdd, DatasetExampleUpdateMetadataParamsOperationsOperationArrayRemove:
		return true
	}
	return false
}
