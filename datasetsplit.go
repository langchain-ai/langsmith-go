// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetSplitService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetSplitService] method instead.
type DatasetSplitService struct {
	Options []option.RequestOption
}

// NewDatasetSplitService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetSplitService(opts ...option.RequestOption) (r *DatasetSplitService) {
	r = &DatasetSplitService{}
	r.Options = opts
	return
}

// Add examples to a split, or remove them from it.
//
// Returns the IDs of the given examples found in the dataset. Each of them gets a
// new version, even when its splits do not change. An example's other splits are
// kept, and an example without a split counts as being in the base split. Removing
// an example's last split leaves it in no split.
func (r *DatasetSplitService) New(ctx context.Context, datasetID string, body DatasetSplitNewParams, opts ...option.RequestOption) (res *[]string, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/splits", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, &res, opts...)
	return res, err
}

// List the split names a dataset's examples belong to.
//
// An example whose split was never set counts as the base split.
func (r *DatasetSplitService) Get(ctx context.Context, datasetID string, query DatasetSplitGetParams, opts ...option.RequestOption) (res *[]string, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/%s/splits", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

type DatasetSplitNewParams struct {
	// IDs of the examples to change. IDs of examples in other datasets are ignored.
	Examples param.Field[[]string] `json:"examples" api:"required" format:"uuid"`
	// Split to add the examples to, or remove them from.
	SplitName param.Field[string] `json:"split_name" api:"required"`
	// Remove the examples from the split instead of adding them.
	Remove param.Field[bool] `json:"remove"`
}

func (r DatasetSplitNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetSplitGetParams struct {
	// Only modifications made on or before this time are included. If None, the latest
	// version of the dataset is used.
	AsOf param.Field[DatasetSplitGetParamsAsOfUnion] `query:"as_of" format:"date-time"`
}

// URLQuery serializes [DatasetSplitGetParams]'s query parameters as `url.Values`.
func (r DatasetSplitGetParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Only modifications made on or before this time are included. If None, the latest
// version of the dataset is used.
//
// Satisfied by [shared.UnionTime], [shared.UnionString].
type DatasetSplitGetParamsAsOfUnion interface {
	ImplementsDatasetSplitGetParamsAsOfUnion()
}
