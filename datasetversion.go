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
	"github.com/stainless-sdks/langsmith-api-go/packages/pagination"
)

// DatasetVersionService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetVersionService] method instead.
type DatasetVersionService struct {
	Options []option.RequestOption
}

// NewDatasetVersionService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetVersionService(opts ...option.RequestOption) (r *DatasetVersionService) {
	r = &DatasetVersionService{}
	r.Options = opts
	return
}

// Get dataset versions.
func (r *DatasetVersionService) List(ctx context.Context, datasetID string, query DatasetVersionListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[DatasetVersion], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/versions", datasetID)
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

// Get dataset versions.
func (r *DatasetVersionService) ListAutoPaging(ctx context.Context, datasetID string, query DatasetVersionListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[DatasetVersion] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, datasetID, query, opts...))
}

// Get diff between two dataset versions.
func (r *DatasetVersionService) GetDiff(ctx context.Context, datasetID string, query DatasetVersionGetDiffParams, opts ...option.RequestOption) (res *DatasetVersionGetDiffResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/versions/diff", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Dataset diff schema.
type DatasetVersionGetDiffResponse struct {
	ExamplesAdded    []string                          `json:"examples_added,required" format:"uuid"`
	ExamplesModified []string                          `json:"examples_modified,required" format:"uuid"`
	ExamplesRemoved  []string                          `json:"examples_removed,required" format:"uuid"`
	JSON             datasetVersionGetDiffResponseJSON `json:"-"`
}

// datasetVersionGetDiffResponseJSON contains the JSON metadata for the struct
// [DatasetVersionGetDiffResponse]
type datasetVersionGetDiffResponseJSON struct {
	ExamplesAdded    apijson.Field
	ExamplesModified apijson.Field
	ExamplesRemoved  apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *DatasetVersionGetDiffResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetVersionGetDiffResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetVersionListParams struct {
	Example param.Field[string] `query:"example" format:"uuid"`
	Limit   param.Field[int64]  `query:"limit"`
	Offset  param.Field[int64]  `query:"offset"`
	Search  param.Field[string] `query:"search"`
}

// URLQuery serializes [DatasetVersionListParams]'s query parameters as
// `url.Values`.
func (r DatasetVersionListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type DatasetVersionGetDiffParams struct {
	FromVersion param.Field[DatasetVersionGetDiffParamsFromVersionUnion] `query:"from_version,required" format:"date-time"`
	ToVersion   param.Field[DatasetVersionGetDiffParamsToVersionUnion]   `query:"to_version,required" format:"date-time"`
}

// URLQuery serializes [DatasetVersionGetDiffParams]'s query parameters as
// `url.Values`.
func (r DatasetVersionGetDiffParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Satisfied by [shared.UnionTime], [shared.UnionString].
type DatasetVersionGetDiffParamsFromVersionUnion interface {
	ImplementsDatasetVersionGetDiffParamsFromVersionUnion()
}

// Satisfied by [shared.UnionTime], [shared.UnionString].
type DatasetVersionGetDiffParamsToVersionUnion interface {
	ImplementsDatasetVersionGetDiffParamsToVersionUnion()
}
