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

// AnnotationQueueInfoService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAnnotationQueueInfoService] method instead.
type AnnotationQueueInfoService struct {
	Options []option.RequestOption
}

// NewAnnotationQueueInfoService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAnnotationQueueInfoService(opts ...option.RequestOption) (r *AnnotationQueueInfoService) {
	r = &AnnotationQueueInfoService{}
	r.Options = opts
	return
}

// Get information about the current deployment of LangSmith.
func (r *AnnotationQueueInfoService) List(ctx context.Context, opts ...option.RequestOption) (res *AnnotationQueueInfoListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/info"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// The LangSmith server info.
type AnnotationQueueInfoListResponse struct {
	Version string `json:"version,required"`
	// Batch ingest config.
	BatchIngestConfig AnnotationQueueInfoListResponseBatchIngestConfig `json:"batch_ingest_config"`
	// Customer info.
	CustomerInfo          AnnotationQueueInfoListResponseCustomerInfo `json:"customer_info,nullable"`
	InstanceFlags         map[string]interface{}                      `json:"instance_flags"`
	LicenseExpirationTime time.Time                                   `json:"license_expiration_time,nullable" format:"date-time"`
	JSON                  annotationQueueInfoListResponseJSON         `json:"-"`
}

// annotationQueueInfoListResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueInfoListResponse]
type annotationQueueInfoListResponseJSON struct {
	Version               apijson.Field
	BatchIngestConfig     apijson.Field
	CustomerInfo          apijson.Field
	InstanceFlags         apijson.Field
	LicenseExpirationTime apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *AnnotationQueueInfoListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueInfoListResponseJSON) RawJSON() string {
	return r.raw
}

// Batch ingest config.
type AnnotationQueueInfoListResponseBatchIngestConfig struct {
	ScaleDownNemptyTrigger int64                                                `json:"scale_down_nempty_trigger"`
	ScaleUpNthreadsLimit   int64                                                `json:"scale_up_nthreads_limit"`
	ScaleUpQsizeTrigger    int64                                                `json:"scale_up_qsize_trigger"`
	SizeLimit              int64                                                `json:"size_limit"`
	SizeLimitBytes         int64                                                `json:"size_limit_bytes"`
	UseMultipartEndpoint   bool                                                 `json:"use_multipart_endpoint"`
	JSON                   annotationQueueInfoListResponseBatchIngestConfigJSON `json:"-"`
}

// annotationQueueInfoListResponseBatchIngestConfigJSON contains the JSON metadata
// for the struct [AnnotationQueueInfoListResponseBatchIngestConfig]
type annotationQueueInfoListResponseBatchIngestConfigJSON struct {
	ScaleDownNemptyTrigger apijson.Field
	ScaleUpNthreadsLimit   apijson.Field
	ScaleUpQsizeTrigger    apijson.Field
	SizeLimit              apijson.Field
	SizeLimitBytes         apijson.Field
	UseMultipartEndpoint   apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *AnnotationQueueInfoListResponseBatchIngestConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueInfoListResponseBatchIngestConfigJSON) RawJSON() string {
	return r.raw
}

// Customer info.
type AnnotationQueueInfoListResponseCustomerInfo struct {
	CustomerID   string                                          `json:"customer_id,required"`
	CustomerName string                                          `json:"customer_name,required"`
	JSON         annotationQueueInfoListResponseCustomerInfoJSON `json:"-"`
}

// annotationQueueInfoListResponseCustomerInfoJSON contains the JSON metadata for
// the struct [AnnotationQueueInfoListResponseCustomerInfo]
type annotationQueueInfoListResponseCustomerInfoJSON struct {
	CustomerID   apijson.Field
	CustomerName apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *AnnotationQueueInfoListResponseCustomerInfo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueInfoListResponseCustomerInfoJSON) RawJSON() string {
	return r.raw
}
