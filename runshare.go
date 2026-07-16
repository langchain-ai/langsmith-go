// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// RunShareService contains methods and other services that help with interacting
// with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewRunShareService] method instead.
type RunShareService struct {
	Options []option.RequestOption
}

// NewRunShareService generates a new service that applies the given options to
// each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewRunShareService(opts ...option.RequestOption) (r *RunShareService) {
	r = &RunShareService{}
	r.Options = opts
	return
}

// Creates or returns a share token for a run. Child runs share their trace root.
func (r *RunShareService) New(ctx context.Context, runID string, body RunShareNewParams, opts ...option.RequestOption) (res *RunShareNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/runs/%s/share", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Deletes the share token for the trace identified by trace_id and session_id.
// Idempotent: returns 204 whether or not a share token existed.
func (r *RunShareService) Delete(ctx context.Context, traceID string, body RunShareDeleteParams, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if traceID == "" {
		err = errors.New("missing required trace_id parameter")
		return err
	}
	path := fmt.Sprintf("v2/runs/%s/share", traceID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, body, nil, opts...)
	return err
}

type RunShareNewResponse struct {
	ShareToken string                  `json:"share_token" format:"uuid"`
	JSON       runShareNewResponseJSON `json:"-"`
}

// runShareNewResponseJSON contains the JSON metadata for the struct
// [RunShareNewResponse]
type runShareNewResponseJSON struct {
	ShareToken  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *RunShareNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runShareNewResponseJSON) RawJSON() string {
	return r.raw
}

type RunShareNewParams struct {
	// session_id is the tracing project UUID containing the trace.
	SessionID param.Field[string] `json:"session_id" format:"uuid"`
	// trace_id is the root trace UUID to share.
	TraceID param.Field[string] `json:"trace_id" format:"uuid"`
}

func (r RunShareNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RunShareDeleteParams struct {
	// session_id is the tracing project UUID containing the trace.
	SessionID param.Field[string] `json:"session_id" format:"uuid"`
}

func (r RunShareDeleteParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
