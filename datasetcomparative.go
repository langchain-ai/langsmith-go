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

// DatasetComparativeService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetComparativeService] method instead.
type DatasetComparativeService struct {
	Options []option.RequestOption
}

// NewDatasetComparativeService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewDatasetComparativeService(opts ...option.RequestOption) (r *DatasetComparativeService) {
	r = &DatasetComparativeService{}
	r.Options = opts
	return
}

// Create a comparative experiment.
func (r *DatasetComparativeService) New(ctx context.Context, body DatasetComparativeNewParams, opts ...option.RequestOption) (res *DatasetComparativeNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/comparative"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Get all comparative experiments for a given dataset.
func (r *DatasetComparativeService) List(ctx context.Context, datasetID string, query DatasetComparativeListParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[DatasetComparativeListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/comparative", datasetID)
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

// Get all comparative experiments for a given dataset.
func (r *DatasetComparativeService) ListAutoPaging(ctx context.Context, datasetID string, query DatasetComparativeListParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[DatasetComparativeListResponse] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.List(ctx, datasetID, query, opts...))
}

// Delete a specific comparative experiment.
func (r *DatasetComparativeService) Delete(ctx context.Context, comparativeExperimentID string, opts ...option.RequestOption) (res *DatasetComparativeDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if comparativeExperimentID == "" {
		err = errors.New("missing required comparative_experiment_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/comparative/%s", comparativeExperimentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Simple experiment info schema for use with comparative experiments
type SimpleExperimentInfo struct {
	ID   string                   `json:"id,required" format:"uuid"`
	Name string                   `json:"name,required"`
	JSON simpleExperimentInfoJSON `json:"-"`
}

// simpleExperimentInfoJSON contains the JSON metadata for the struct
// [SimpleExperimentInfo]
type simpleExperimentInfoJSON struct {
	ID          apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SimpleExperimentInfo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r simpleExperimentInfoJSON) RawJSON() string {
	return r.raw
}

// Enum for available comparative experiment columns to sort by.
type SortByComparativeExperimentColumn string

const (
	SortByComparativeExperimentColumnName      SortByComparativeExperimentColumn = "name"
	SortByComparativeExperimentColumnCreatedAt SortByComparativeExperimentColumn = "created_at"
)

func (r SortByComparativeExperimentColumn) IsKnown() bool {
	switch r {
	case SortByComparativeExperimentColumnName, SortByComparativeExperimentColumnCreatedAt:
		return true
	}
	return false
}

// ComparativeExperiment schema.
type DatasetComparativeNewResponse struct {
	ID                 string                            `json:"id,required" format:"uuid"`
	CreatedAt          time.Time                         `json:"created_at,required" format:"date-time"`
	ModifiedAt         time.Time                         `json:"modified_at,required" format:"date-time"`
	ReferenceDatasetID string                            `json:"reference_dataset_id,required" format:"uuid"`
	TenantID           string                            `json:"tenant_id,required" format:"uuid"`
	Description        string                            `json:"description,nullable"`
	Extra              interface{}                       `json:"extra,nullable"`
	Name               string                            `json:"name,nullable"`
	JSON               datasetComparativeNewResponseJSON `json:"-"`
}

// datasetComparativeNewResponseJSON contains the JSON metadata for the struct
// [DatasetComparativeNewResponse]
type datasetComparativeNewResponseJSON struct {
	ID                 apijson.Field
	CreatedAt          apijson.Field
	ModifiedAt         apijson.Field
	ReferenceDatasetID apijson.Field
	TenantID           apijson.Field
	Description        apijson.Field
	Extra              apijson.Field
	Name               apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *DatasetComparativeNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetComparativeNewResponseJSON) RawJSON() string {
	return r.raw
}

// ComparativeExperiment schema.
type DatasetComparativeListResponse struct {
	ID                 string                             `json:"id,required" format:"uuid"`
	CreatedAt          time.Time                          `json:"created_at,required" format:"date-time"`
	ExperimentsInfo    []SimpleExperimentInfo             `json:"experiments_info,required"`
	ModifiedAt         time.Time                          `json:"modified_at,required" format:"date-time"`
	ReferenceDatasetID string                             `json:"reference_dataset_id,required" format:"uuid"`
	TenantID           string                             `json:"tenant_id,required" format:"uuid"`
	Description        string                             `json:"description,nullable"`
	Extra              interface{}                        `json:"extra,nullable"`
	FeedbackStats      interface{}                        `json:"feedback_stats,nullable"`
	Name               string                             `json:"name,nullable"`
	JSON               datasetComparativeListResponseJSON `json:"-"`
}

// datasetComparativeListResponseJSON contains the JSON metadata for the struct
// [DatasetComparativeListResponse]
type datasetComparativeListResponseJSON struct {
	ID                 apijson.Field
	CreatedAt          apijson.Field
	ExperimentsInfo    apijson.Field
	ModifiedAt         apijson.Field
	ReferenceDatasetID apijson.Field
	TenantID           apijson.Field
	Description        apijson.Field
	Extra              apijson.Field
	FeedbackStats      apijson.Field
	Name               apijson.Field
	raw                string
	ExtraFields        map[string]apijson.Field
}

func (r *DatasetComparativeListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetComparativeListResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetComparativeDeleteResponse = interface{}

type DatasetComparativeNewParams struct {
	ExperimentIDs      param.Field[[]string]    `json:"experiment_ids,required" format:"uuid"`
	ID                 param.Field[string]      `json:"id" format:"uuid"`
	CreatedAt          param.Field[time.Time]   `json:"created_at" format:"date-time"`
	Description        param.Field[string]      `json:"description"`
	Extra              param.Field[interface{}] `json:"extra"`
	ModifiedAt         param.Field[time.Time]   `json:"modified_at" format:"date-time"`
	Name               param.Field[string]      `json:"name"`
	ReferenceDatasetID param.Field[string]      `json:"reference_dataset_id" format:"uuid"`
}

func (r DatasetComparativeNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetComparativeListParams struct {
	ID           param.Field[[]string] `query:"id" format:"uuid"`
	Limit        param.Field[int64]    `query:"limit"`
	Name         param.Field[string]   `query:"name"`
	NameContains param.Field[string]   `query:"name_contains"`
	Offset       param.Field[int64]    `query:"offset"`
	// Enum for available comparative experiment columns to sort by.
	SortBy     param.Field[SortByComparativeExperimentColumn] `query:"sort_by"`
	SortByDesc param.Field[bool]                              `query:"sort_by_desc"`
}

// URLQuery serializes [DatasetComparativeListParams]'s query parameters as
// `url.Values`.
func (r DatasetComparativeListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
