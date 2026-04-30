// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// InfoService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewInfoService] method instead.
type InfoService struct {
	Options []option.RequestOption
}

// NewInfoService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewInfoService(opts ...option.RequestOption) (r *InfoService) {
	r = &InfoService{}
	r.Options = opts
	return
}

// Get information about the current deployment of LangSmith.
func (r *InfoService) List(ctx context.Context, opts ...option.RequestOption) (res *InfoListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/info"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// The LangSmith server info.
type InfoListResponse struct {
	Version string `json:"version" api:"required"`
	// Batch ingest config.
	BatchIngestConfig InfoListResponseBatchIngestConfig `json:"batch_ingest_config"`
	// Customer info.
	CustomerInfo          InfoListResponseCustomerInfo `json:"customer_info" api:"nullable"`
	GitSha                string                       `json:"git_sha" api:"nullable"`
	InstanceFlags         map[string]interface{}       `json:"instance_flags"`
	LicenseExpirationTime time.Time                    `json:"license_expiration_time" api:"nullable" format:"date-time"`
	JSON                  infoListResponseJSON         `json:"-"`
}

// infoListResponseJSON contains the JSON metadata for the struct
// [InfoListResponse]
type infoListResponseJSON struct {
	Version               apijson.Field
	BatchIngestConfig     apijson.Field
	CustomerInfo          apijson.Field
	GitSha                apijson.Field
	InstanceFlags         apijson.Field
	LicenseExpirationTime apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *InfoListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r infoListResponseJSON) RawJSON() string {
	return r.raw
}

// Batch ingest config.
type InfoListResponseBatchIngestConfig struct {
	ScaleDownNemptyTrigger int64                                 `json:"scale_down_nempty_trigger"`
	ScaleUpNthreadsLimit   int64                                 `json:"scale_up_nthreads_limit"`
	ScaleUpQsizeTrigger    int64                                 `json:"scale_up_qsize_trigger"`
	SizeLimit              int64                                 `json:"size_limit"`
	SizeLimitBytes         int64                                 `json:"size_limit_bytes"`
	UseMultipartEndpoint   bool                                  `json:"use_multipart_endpoint"`
	JSON                   infoListResponseBatchIngestConfigJSON `json:"-"`
}

// infoListResponseBatchIngestConfigJSON contains the JSON metadata for the struct
// [InfoListResponseBatchIngestConfig]
type infoListResponseBatchIngestConfigJSON struct {
	ScaleDownNemptyTrigger apijson.Field
	ScaleUpNthreadsLimit   apijson.Field
	ScaleUpQsizeTrigger    apijson.Field
	SizeLimit              apijson.Field
	SizeLimitBytes         apijson.Field
	UseMultipartEndpoint   apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *InfoListResponseBatchIngestConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r infoListResponseBatchIngestConfigJSON) RawJSON() string {
	return r.raw
}

// Customer info.
type InfoListResponseCustomerInfo struct {
	CustomerID   string                           `json:"customer_id" api:"required"`
	CustomerName string                           `json:"customer_name" api:"required"`
	JSON         infoListResponseCustomerInfoJSON `json:"-"`
}

// infoListResponseCustomerInfoJSON contains the JSON metadata for the struct
// [InfoListResponseCustomerInfo]
type infoListResponseCustomerInfoJSON struct {
	CustomerID   apijson.Field
	CustomerName apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *InfoListResponseCustomerInfo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r infoListResponseCustomerInfoJSON) RawJSON() string {
	return r.raw
}
