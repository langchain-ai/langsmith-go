// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"github.com/langchain-ai/langsmith-go/option"
)

// SandboxPoolService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxPoolService] method instead.
type SandboxPoolService struct {
	Options []option.RequestOption
}

// NewSandboxPoolService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxPoolService(opts ...option.RequestOption) (r *SandboxPoolService) {
	r = &SandboxPoolService{}
	r.Options = opts
	return
}
