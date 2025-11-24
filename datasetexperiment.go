// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/stainless-sdks/langsmith-api-go/internal/apijson"
	"github.com/stainless-sdks/langsmith-api-go/internal/param"
	"github.com/stainless-sdks/langsmith-api-go/internal/requestconfig"
	"github.com/stainless-sdks/langsmith-api-go/option"
)

// DatasetExperimentService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetExperimentService] method instead.
type DatasetExperimentService struct {
	Options []option.RequestOption
}

// NewDatasetExperimentService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewDatasetExperimentService(opts ...option.RequestOption) (r *DatasetExperimentService) {
	r = &DatasetExperimentService{}
	r.Options = opts
	return
}

// Stream grouped and aggregated experiments.
func (r *DatasetExperimentService) Grouped(ctx context.Context, datasetID string, body DatasetExperimentGroupedParams, opts ...option.RequestOption) (res *DatasetExperimentGroupedResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/datasets/%s/experiments/grouped", datasetID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

type DatasetExperimentGroupedResponse = interface{}

type DatasetExperimentGroupedParams struct {
	MetadataKeys    param.Field[[]string]  `json:"metadata_keys,required"`
	DatasetVersion  param.Field[string]    `json:"dataset_version"`
	ExperimentLimit param.Field[int64]     `json:"experiment_limit"`
	Filter          param.Field[string]    `json:"filter"`
	NameContains    param.Field[string]    `json:"name_contains"`
	StatsStartTime  param.Field[time.Time] `json:"stats_start_time" format:"date-time"`
	TagValueID      param.Field[[]string]  `json:"tag_value_id" format:"uuid"`
	UseApproxStats  param.Field[bool]      `json:"use_approx_stats"`
}

func (r DatasetExperimentGroupedParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
