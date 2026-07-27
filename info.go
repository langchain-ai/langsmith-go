// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"

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

// Returns information about the current LangSmith deployment: version, instance
// feature flags, batch-ingest limits, and max SDK versions. Unauthenticated by
// default; set FF_INFO_ENDPOINT_AUTH_REQUIRED=true to require auth.
func (r *InfoService) List(ctx context.Context, opts ...option.RequestOption) (res *InfoListResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/info"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type InfoListResponse struct {
	BatchIngestConfig     InfoListResponseBatchIngestConfig `json:"batch_ingest_config"`
	CustomerInfo          InfoListResponseCustomerInfo      `json:"customer_info"`
	GitSha                string                            `json:"git_sha"`
	InstanceFlags         map[string]interface{}            `json:"instance_flags"`
	LicenseExpirationTime string                            `json:"license_expiration_time"`
	SDKVersions           InfoListResponseSDKVersions       `json:"sdk_versions"`
	Version               string                            `json:"version"`
	JSON                  infoListResponseJSON              `json:"-"`
}

// infoListResponseJSON contains the JSON metadata for the struct
// [InfoListResponse]
type infoListResponseJSON struct {
	BatchIngestConfig     apijson.Field
	CustomerInfo          apijson.Field
	GitSha                apijson.Field
	InstanceFlags         apijson.Field
	LicenseExpirationTime apijson.Field
	SDKVersions           apijson.Field
	Version               apijson.Field
	raw                   string
	ExtraFields           map[string]apijson.Field
}

func (r *InfoListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r infoListResponseJSON) RawJSON() string {
	return r.raw
}

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

type InfoListResponseCustomerInfo struct {
	CustomerID   string                           `json:"customer_id"`
	CustomerName string                           `json:"customer_name"`
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

type InfoListResponseSDKVersions struct {
	MaxGoSDKVersion     string                          `json:"max_go_sdk_version"`
	MaxJavaSDKVersion   string                          `json:"max_java_sdk_version"`
	MaxJsSDKVersion     string                          `json:"max_js_sdk_version"`
	MaxPythonSDKVersion string                          `json:"max_python_sdk_version"`
	JSON                infoListResponseSDKVersionsJSON `json:"-"`
}

// infoListResponseSDKVersionsJSON contains the JSON metadata for the struct
// [InfoListResponseSDKVersions]
type infoListResponseSDKVersionsJSON struct {
	MaxGoSDKVersion     apijson.Field
	MaxJavaSDKVersion   apijson.Field
	MaxJsSDKVersion     apijson.Field
	MaxPythonSDKVersion apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *InfoListResponseSDKVersions) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r infoListResponseSDKVersionsJSON) RawJSON() string {
	return r.raw
}
