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

// Update Dataset Splits
func (r *DatasetSplitService) New(ctx context.Context, datasetID string, body DatasetSplitNewParams, opts ...option.RequestOption) (res *[]string, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/splits", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPut, path, body, &res, opts...)
	return
}

// Get Dataset Splits
func (r *DatasetSplitService) Get(ctx context.Context, datasetID string, query DatasetSplitGetParams, opts ...option.RequestOption) (res *[]string, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/splits", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

type DatasetSplitNewParams struct {
	Examples  param.Field[[]string] `json:"examples,required" format:"uuid"`
	SplitName param.Field[string]   `json:"split_name,required"`
	Remove    param.Field[bool]     `json:"remove"`
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
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
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
