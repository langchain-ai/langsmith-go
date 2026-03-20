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
	return res, err
}

// Delete a specific comparative experiment.
func (r *DatasetComparativeService) Delete(ctx context.Context, comparativeExperimentID string, opts ...option.RequestOption) (res *DatasetComparativeDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if comparativeExperimentID == "" {
		err = errors.New("missing required comparative_experiment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/datasets/comparative/%s", comparativeExperimentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Simple experiment info schema for use with comparative experiments
type SimpleExperimentInfo struct {
	ID   string                   `json:"id" api:"required" format:"uuid"`
	Name string                   `json:"name" api:"required"`
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
	ID                 string                            `json:"id" api:"required" format:"uuid"`
	CreatedAt          time.Time                         `json:"created_at" api:"required" format:"date-time"`
	ModifiedAt         time.Time                         `json:"modified_at" api:"required" format:"date-time"`
	ReferenceDatasetID string                            `json:"reference_dataset_id" api:"required" format:"uuid"`
	TenantID           string                            `json:"tenant_id" api:"required" format:"uuid"`
	Description        string                            `json:"description" api:"nullable"`
	Extra              map[string]interface{}            `json:"extra" api:"nullable"`
	Name               string                            `json:"name" api:"nullable"`
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

type DatasetComparativeDeleteResponse = interface{}

type DatasetComparativeNewParams struct {
	ExperimentIDs      param.Field[[]string]               `json:"experiment_ids" api:"required" format:"uuid"`
	ID                 param.Field[string]                 `json:"id" format:"uuid"`
	CreatedAt          param.Field[time.Time]              `json:"created_at" format:"date-time"`
	Description        param.Field[string]                 `json:"description"`
	Extra              param.Field[map[string]interface{}] `json:"extra"`
	ModifiedAt         param.Field[time.Time]              `json:"modified_at" format:"date-time"`
	Name               param.Field[string]                 `json:"name"`
	ReferenceDatasetID param.Field[string]                 `json:"reference_dataset_id" format:"uuid"`
}

func (r DatasetComparativeNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
