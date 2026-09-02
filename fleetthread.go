// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// FleetThreadService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewFleetThreadService] method instead.
type FleetThreadService struct {
	Options []option.RequestOption
}

// NewFleetThreadService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewFleetThreadService(opts ...option.RequestOption) (r *FleetThreadService) {
	r = &FleetThreadService{}
	r.Options = opts
	return
}

// Starts or resumes the sandbox referenced by the thread and returns when it is
// ready. The operation is idempotent. The thread must include
// sandbox.sandbox_slug.
func (r *FleetThreadService) ActivateSandbox(ctx context.Context, threadID string, opts ...option.RequestOption) (res *FleetThreadActivateSandboxResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if threadID == "" {
		err = errors.New("missing required thread_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v1/fleet/threads/%s/sandbox-activation", threadID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type FleetThreadActivateSandboxResponse struct {
	SandboxSlug string                                   `json:"sandbox_slug" api:"required"`
	Scope       FleetThreadActivateSandboxResponseScope  `json:"scope" api:"required"`
	Status      FleetThreadActivateSandboxResponseStatus `json:"status" api:"required"`
	JSON        fleetThreadActivateSandboxResponseJSON   `json:"-"`
}

// fleetThreadActivateSandboxResponseJSON contains the JSON metadata for the struct
// [FleetThreadActivateSandboxResponse]
type fleetThreadActivateSandboxResponseJSON struct {
	SandboxSlug apijson.Field
	Scope       apijson.Field
	Status      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *FleetThreadActivateSandboxResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r fleetThreadActivateSandboxResponseJSON) RawJSON() string {
	return r.raw
}

type FleetThreadActivateSandboxResponseScope string

const (
	FleetThreadActivateSandboxResponseScopeAgent  FleetThreadActivateSandboxResponseScope = "agent"
	FleetThreadActivateSandboxResponseScopeThread FleetThreadActivateSandboxResponseScope = "thread"
)

func (r FleetThreadActivateSandboxResponseScope) IsKnown() bool {
	switch r {
	case FleetThreadActivateSandboxResponseScopeAgent, FleetThreadActivateSandboxResponseScopeThread:
		return true
	}
	return false
}

type FleetThreadActivateSandboxResponseStatus string

const (
	FleetThreadActivateSandboxResponseStatusProvisioning FleetThreadActivateSandboxResponseStatus = "provisioning"
	FleetThreadActivateSandboxResponseStatusReady        FleetThreadActivateSandboxResponseStatus = "ready"
	FleetThreadActivateSandboxResponseStatusFailed       FleetThreadActivateSandboxResponseStatus = "failed"
	FleetThreadActivateSandboxResponseStatusStopped      FleetThreadActivateSandboxResponseStatus = "stopped"
	FleetThreadActivateSandboxResponseStatusDeleting     FleetThreadActivateSandboxResponseStatus = "deleting"
)

func (r FleetThreadActivateSandboxResponseStatus) IsKnown() bool {
	switch r {
	case FleetThreadActivateSandboxResponseStatusProvisioning, FleetThreadActivateSandboxResponseStatusReady, FleetThreadActivateSandboxResponseStatusFailed, FleetThreadActivateSandboxResponseStatusStopped, FleetThreadActivateSandboxResponseStatusDeleting:
		return true
	}
	return false
}
