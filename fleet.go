// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"github.com/langchain-ai/langsmith-go/option"
)

// FleetService contains methods and other services that help with interacting with
// the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFleetService] method instead.
type FleetService struct {
	Options []option.RequestOption
	Threads *FleetThreadService
}

// NewFleetService generates a new service that applies the given options to each
// request. These options are applied after the parent client's options (if there
// is one), and before any request-specific options.
func NewFleetService(opts ...option.RequestOption) (r *FleetService) {
	r = &FleetService{}
	r.Options = opts
	r.Threads = NewFleetThreadService(opts...)
	return
}
