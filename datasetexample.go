// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetExampleService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetExampleService] method instead.
type DatasetExampleService struct {
	Options []option.RequestOption
}

// NewDatasetExampleService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetExampleService(opts ...option.RequestOption) (r *DatasetExampleService) {
	r = &DatasetExampleService{}
	r.Options = opts
	return
}

// Soft-delete an example, preserving prior versions and their attachments. If the
// latest version is already deleted, the request succeeds without creating another
// version. Deletion is recorded at the current time or just after the latest
// version, whichever is later. For future-dated versions, latest reads reflect
// deletion immediately; timestamp reads reflect deletion only at or after the
// recorded deletion timestamp.
func (r *DatasetExampleService) Delete(ctx context.Context, datasetID string, exampleID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return err
	}
	if exampleID == "" {
		err = errors.New("missing required example_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v1/platform/datasets/%s/examples/%s", datasetID, exampleID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}
