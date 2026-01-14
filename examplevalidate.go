// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"reflect"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/shared"
	"github.com/tidwall/gjson"
)

// ExampleValidateService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewExampleValidateService] method instead.
type ExampleValidateService struct {
	Options []option.RequestOption
}

// NewExampleValidateService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewExampleValidateService(opts ...option.RequestOption) (r *ExampleValidateService) {
	r = &ExampleValidateService{}
	r.Options = opts
	return
}

// Validate an example.
func (r *ExampleValidateService) New(ctx context.Context, opts ...option.RequestOption) (res *ExampleValidationResult, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples/validate"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Validate examples in bulk.
func (r *ExampleValidateService) Bulk(ctx context.Context, opts ...option.RequestOption) (res *[]ExampleValidationResult, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/examples/validate/bulk"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return
}

// Validation result for Example, combining fields from Create/Base/Update schemas.
type ExampleValidationResult struct {
	ID             string                            `json:"id,nullable" format:"uuid"`
	CreatedAt      time.Time                         `json:"created_at,nullable" format:"date-time"`
	DatasetID      string                            `json:"dataset_id,nullable" format:"uuid"`
	Inputs         map[string]interface{}            `json:"inputs,nullable"`
	Metadata       map[string]interface{}            `json:"metadata,nullable"`
	Outputs        map[string]interface{}            `json:"outputs,nullable"`
	Overwrite      bool                              `json:"overwrite"`
	SourceRunID    string                            `json:"source_run_id,nullable" format:"uuid"`
	Split          ExampleValidationResultSplitUnion `json:"split,nullable"`
	UseSourceRunIo bool                              `json:"use_source_run_io"`
	JSON           exampleValidationResultJSON       `json:"-"`
}

// exampleValidationResultJSON contains the JSON metadata for the struct
// [ExampleValidationResult]
type exampleValidationResultJSON struct {
	ID             apijson.Field
	CreatedAt      apijson.Field
	DatasetID      apijson.Field
	Inputs         apijson.Field
	Metadata       apijson.Field
	Outputs        apijson.Field
	Overwrite      apijson.Field
	SourceRunID    apijson.Field
	Split          apijson.Field
	UseSourceRunIo apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *ExampleValidationResult) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r exampleValidationResultJSON) RawJSON() string {
	return r.raw
}

// Union satisfied by [ExampleValidationResultSplitArray] or [shared.UnionString].
type ExampleValidationResultSplitUnion interface {
	ImplementsExampleValidationResultSplitUnion()
}

func init() {
	apijson.RegisterUnion(
		reflect.TypeOf((*ExampleValidationResultSplitUnion)(nil)).Elem(),
		"",
		apijson.UnionVariant{
			TypeFilter: gjson.JSON,
			Type:       reflect.TypeOf(ExampleValidationResultSplitArray{}),
		},
		apijson.UnionVariant{
			TypeFilter: gjson.String,
			Type:       reflect.TypeOf(shared.UnionString("")),
		},
	)
}

type ExampleValidationResultSplitArray []string

func (r ExampleValidationResultSplitArray) ImplementsExampleValidationResultSplitUnion() {}
