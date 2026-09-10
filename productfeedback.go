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

// ProductFeedbackService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewProductFeedbackService] method instead.
type ProductFeedbackService struct {
	Options []option.RequestOption
}

// NewProductFeedbackService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewProductFeedbackService(opts ...option.RequestOption) (r *ProductFeedbackService) {
	r = &ProductFeedbackService{}
	r.Options = opts
	return
}

// **Alpha:** This endpoint is in active development and may change without notice.
//
// Submits concise product feedback with optional non-sensitive client details.
func (r *ProductFeedbackService) New(ctx context.Context, params ProductFeedbackNewParams, opts ...option.RequestOption) (res *ProductFeedbackNewResponse, err error) {
	if params.IdempotencyKey.Present {
		opts = append(opts, option.WithHeader("Idempotency-Key", fmt.Sprintf("%v", params.IdempotencyKey)))
	}
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/platform/product-feedbacks"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// **Alpha:** This endpoint is in active development and may change without notice.
func (r *ProductFeedbackService) Get(ctx context.Context, id string, opts ...option.RequestOption) (res *ProductFeedbackGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if id == "" {
		err = errors.New("missing required id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/platform/product-feedbacks/%s", id)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type ProductFeedbackNewResponse struct {
	ID        string                             `json:"id" api:"required" format:"uuid"`
	Category  ProductFeedbackNewResponseCategory `json:"category" api:"required"`
	CreatedAt time.Time                          `json:"created_at" api:"required" format:"date-time"`
	Message   string                             `json:"message" api:"required"`
	Source    ProductFeedbackNewResponseSource   `json:"source" api:"required"`
	Client    ProductFeedbackNewResponseClient   `json:"client"`
	JSON      productFeedbackNewResponseJSON     `json:"-"`
}

// productFeedbackNewResponseJSON contains the JSON metadata for the struct
// [ProductFeedbackNewResponse]
type productFeedbackNewResponseJSON struct {
	ID          apijson.Field
	Category    apijson.Field
	CreatedAt   apijson.Field
	Message     apijson.Field
	Source      apijson.Field
	Client      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProductFeedbackNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r productFeedbackNewResponseJSON) RawJSON() string {
	return r.raw
}

type ProductFeedbackNewResponseCategory string

const (
	ProductFeedbackNewResponseCategoryBug            ProductFeedbackNewResponseCategory = "BUG"
	ProductFeedbackNewResponseCategoryFeatureRequest ProductFeedbackNewResponseCategory = "FEATURE_REQUEST"
	ProductFeedbackNewResponseCategoryUsability      ProductFeedbackNewResponseCategory = "USABILITY"
	ProductFeedbackNewResponseCategoryDocumentation  ProductFeedbackNewResponseCategory = "DOCUMENTATION"
	ProductFeedbackNewResponseCategoryOther          ProductFeedbackNewResponseCategory = "OTHER"
)

func (r ProductFeedbackNewResponseCategory) IsKnown() bool {
	switch r {
	case ProductFeedbackNewResponseCategoryBug, ProductFeedbackNewResponseCategoryFeatureRequest, ProductFeedbackNewResponseCategoryUsability, ProductFeedbackNewResponseCategoryDocumentation, ProductFeedbackNewResponseCategoryOther:
		return true
	}
	return false
}

type ProductFeedbackNewResponseSource string

const (
	ProductFeedbackNewResponseSourceLangsmithCli ProductFeedbackNewResponseSource = "LANGSMITH_CLI"
)

func (r ProductFeedbackNewResponseSource) IsKnown() bool {
	switch r {
	case ProductFeedbackNewResponseSourceLangsmithCli:
		return true
	}
	return false
}

type ProductFeedbackNewResponseClient struct {
	Architecture string                               `json:"architecture"`
	Os           string                               `json:"os"`
	Version      string                               `json:"version"`
	JSON         productFeedbackNewResponseClientJSON `json:"-"`
}

// productFeedbackNewResponseClientJSON contains the JSON metadata for the struct
// [ProductFeedbackNewResponseClient]
type productFeedbackNewResponseClientJSON struct {
	Architecture apijson.Field
	Os           apijson.Field
	Version      apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ProductFeedbackNewResponseClient) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r productFeedbackNewResponseClientJSON) RawJSON() string {
	return r.raw
}

type ProductFeedbackGetResponse struct {
	ID        string                             `json:"id" api:"required" format:"uuid"`
	Category  ProductFeedbackGetResponseCategory `json:"category" api:"required"`
	CreatedAt time.Time                          `json:"created_at" api:"required" format:"date-time"`
	Message   string                             `json:"message" api:"required"`
	Source    ProductFeedbackGetResponseSource   `json:"source" api:"required"`
	Client    ProductFeedbackGetResponseClient   `json:"client"`
	JSON      productFeedbackGetResponseJSON     `json:"-"`
}

// productFeedbackGetResponseJSON contains the JSON metadata for the struct
// [ProductFeedbackGetResponse]
type productFeedbackGetResponseJSON struct {
	ID          apijson.Field
	Category    apijson.Field
	CreatedAt   apijson.Field
	Message     apijson.Field
	Source      apijson.Field
	Client      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *ProductFeedbackGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r productFeedbackGetResponseJSON) RawJSON() string {
	return r.raw
}

type ProductFeedbackGetResponseCategory string

const (
	ProductFeedbackGetResponseCategoryBug            ProductFeedbackGetResponseCategory = "BUG"
	ProductFeedbackGetResponseCategoryFeatureRequest ProductFeedbackGetResponseCategory = "FEATURE_REQUEST"
	ProductFeedbackGetResponseCategoryUsability      ProductFeedbackGetResponseCategory = "USABILITY"
	ProductFeedbackGetResponseCategoryDocumentation  ProductFeedbackGetResponseCategory = "DOCUMENTATION"
	ProductFeedbackGetResponseCategoryOther          ProductFeedbackGetResponseCategory = "OTHER"
)

func (r ProductFeedbackGetResponseCategory) IsKnown() bool {
	switch r {
	case ProductFeedbackGetResponseCategoryBug, ProductFeedbackGetResponseCategoryFeatureRequest, ProductFeedbackGetResponseCategoryUsability, ProductFeedbackGetResponseCategoryDocumentation, ProductFeedbackGetResponseCategoryOther:
		return true
	}
	return false
}

type ProductFeedbackGetResponseSource string

const (
	ProductFeedbackGetResponseSourceLangsmithCli ProductFeedbackGetResponseSource = "LANGSMITH_CLI"
)

func (r ProductFeedbackGetResponseSource) IsKnown() bool {
	switch r {
	case ProductFeedbackGetResponseSourceLangsmithCli:
		return true
	}
	return false
}

type ProductFeedbackGetResponseClient struct {
	Architecture string                               `json:"architecture"`
	Os           string                               `json:"os"`
	Version      string                               `json:"version"`
	JSON         productFeedbackGetResponseClientJSON `json:"-"`
}

// productFeedbackGetResponseClientJSON contains the JSON metadata for the struct
// [ProductFeedbackGetResponseClient]
type productFeedbackGetResponseClientJSON struct {
	Architecture apijson.Field
	Os           apijson.Field
	Version      apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *ProductFeedbackGetResponseClient) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r productFeedbackGetResponseClientJSON) RawJSON() string {
	return r.raw
}

type ProductFeedbackNewParams struct {
	Category       param.Field[ProductFeedbackNewParamsCategory] `json:"category" api:"required"`
	Message        param.Field[string]                           `json:"message" api:"required"`
	Source         param.Field[ProductFeedbackNewParamsSource]   `json:"source" api:"required"`
	Client         param.Field[ProductFeedbackNewParamsClient]   `json:"client"`
	IdempotencyKey param.Field[string]                           `header:"Idempotency-Key"`
}

func (r ProductFeedbackNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type ProductFeedbackNewParamsCategory string

const (
	ProductFeedbackNewParamsCategoryBug            ProductFeedbackNewParamsCategory = "BUG"
	ProductFeedbackNewParamsCategoryFeatureRequest ProductFeedbackNewParamsCategory = "FEATURE_REQUEST"
	ProductFeedbackNewParamsCategoryUsability      ProductFeedbackNewParamsCategory = "USABILITY"
	ProductFeedbackNewParamsCategoryDocumentation  ProductFeedbackNewParamsCategory = "DOCUMENTATION"
	ProductFeedbackNewParamsCategoryOther          ProductFeedbackNewParamsCategory = "OTHER"
)

func (r ProductFeedbackNewParamsCategory) IsKnown() bool {
	switch r {
	case ProductFeedbackNewParamsCategoryBug, ProductFeedbackNewParamsCategoryFeatureRequest, ProductFeedbackNewParamsCategoryUsability, ProductFeedbackNewParamsCategoryDocumentation, ProductFeedbackNewParamsCategoryOther:
		return true
	}
	return false
}

type ProductFeedbackNewParamsSource string

const (
	ProductFeedbackNewParamsSourceLangsmithCli ProductFeedbackNewParamsSource = "LANGSMITH_CLI"
)

func (r ProductFeedbackNewParamsSource) IsKnown() bool {
	switch r {
	case ProductFeedbackNewParamsSourceLangsmithCli:
		return true
	}
	return false
}

type ProductFeedbackNewParamsClient struct {
	Architecture param.Field[string] `json:"architecture"`
	Os           param.Field[string] `json:"os"`
	Version      param.Field[string] `json:"version"`
}

func (r ProductFeedbackNewParamsClient) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
