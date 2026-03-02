// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetIndexService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetIndexService] method instead.
type DatasetIndexService struct {
	Options []option.RequestOption
}

// NewDatasetIndexService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewDatasetIndexService(opts ...option.RequestOption) (r *DatasetIndexService) {
	r = &DatasetIndexService{}
	r.Options = opts
	return
}
