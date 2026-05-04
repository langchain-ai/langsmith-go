// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// AnnotationQueueRunService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAnnotationQueueRunService] method instead.
type AnnotationQueueRunService struct {
	Options []option.RequestOption
}

// NewAnnotationQueueRunService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAnnotationQueueRunService(opts ...option.RequestOption) (r *AnnotationQueueRunService) {
	r = &AnnotationQueueRunService{}
	r.Options = opts
	return
}

// Add Runs To Annotation Queue
func (r *AnnotationQueueRunService) New(ctx context.Context, queueID string, body AnnotationQueueRunNewParams, opts ...option.RequestOption) (res *[]AnnotationQueueRunNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/runs", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Update Run In Annotation Queue
func (r *AnnotationQueueRunService) Update(ctx context.Context, queueID string, queueRunID string, body AnnotationQueueRunUpdateParams, opts ...option.RequestOption) (res *AnnotationQueueRunUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	if queueRunID == "" {
		err = errors.New("missing required queue_run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/runs/%s", queueID, queueRunID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Get Runs From Annotation Queue
func (r *AnnotationQueueRunService) List(ctx context.Context, queueID string, query AnnotationQueueRunListParams, opts ...option.RequestOption) (res *[]RunSchemaWithAnnotationQueueInfo, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/runs", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Delete Runs From Annotation Queue
func (r *AnnotationQueueRunService) DeleteAll(ctx context.Context, queueID string, body AnnotationQueueRunDeleteAllParams, opts ...option.RequestOption) (res *AnnotationQueueRunDeleteAllResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/runs/delete", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Delete Run From Annotation Queue
func (r *AnnotationQueueRunService) DeleteQueue(ctx context.Context, queueID string, queueRunID string, opts ...option.RequestOption) (res *AnnotationQueueRunDeleteQueueResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	if queueRunID == "" {
		err = errors.New("missing required queue_run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/runs/%s", queueID, queueRunID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

type AnnotationQueueRunNewResponse struct {
	ID                      string                            `json:"id" api:"required" format:"uuid"`
	QueueID                 string                            `json:"queue_id" api:"required" format:"uuid"`
	RunID                   string                            `json:"run_id" api:"required" format:"uuid"`
	AddedAt                 time.Time                         `json:"added_at" format:"date-time"`
	LastReviewedTime        time.Time                         `json:"last_reviewed_time" api:"nullable" format:"date-time"`
	SourceProposedExampleID string                            `json:"source_proposed_example_id" api:"nullable" format:"uuid"`
	JSON                    annotationQueueRunNewResponseJSON `json:"-"`
}

// annotationQueueRunNewResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueRunNewResponse]
type annotationQueueRunNewResponseJSON struct {
	ID                      apijson.Field
	QueueID                 apijson.Field
	RunID                   apijson.Field
	AddedAt                 apijson.Field
	LastReviewedTime        apijson.Field
	SourceProposedExampleID apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *AnnotationQueueRunNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueRunNewResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueRunUpdateResponse = interface{}

type AnnotationQueueRunDeleteAllResponse = interface{}

type AnnotationQueueRunDeleteQueueResponse = interface{}

type AnnotationQueueRunNewParams struct {
	Body AnnotationQueueRunNewParamsBodyUnion `json:"body" api:"required" format:"uuid"`
}

func (r AnnotationQueueRunNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.Body)
}

// Satisfied by [AnnotationQueueRunNewParamsBodyRunsUuidArray],
// [AnnotationQueueRunNewParamsBodyRunsAnnotationQueueRunAddSchemaArray],
// [AnnotationQueueRunNewParamsBodyArray].
type AnnotationQueueRunNewParamsBodyUnion interface {
	implementsAnnotationQueueRunNewParamsBodyUnion()
}

type AnnotationQueueRunNewParamsBodyRunsUuidArray []string

func (r AnnotationQueueRunNewParamsBodyRunsUuidArray) implementsAnnotationQueueRunNewParamsBodyUnion() {
}

type AnnotationQueueRunNewParamsBodyRunsAnnotationQueueRunAddSchemaArray []AnnotationQueueRunNewParamsBodyRunsAnnotationQueueRunAddSchemaArrayItem

func (r AnnotationQueueRunNewParamsBodyRunsAnnotationQueueRunAddSchemaArray) implementsAnnotationQueueRunNewParamsBodyUnion() {
}

// Add a single run to AQ (CH path) with an optional back-pointer to the
// issues-agent proposal that seeded this add. Use when bulk-adding runs that come
// from different proposals — each row carries its own source_proposed_example_id.
// For unrelated bulk adds, prefer plain List[UUID] on the same endpoint.
type AnnotationQueueRunNewParamsBodyRunsAnnotationQueueRunAddSchemaArrayItem struct {
	RunID                   param.Field[string] `json:"run_id" api:"required" format:"uuid"`
	SourceProposedExampleID param.Field[string] `json:"source_proposed_example_id" format:"uuid"`
}

func (r AnnotationQueueRunNewParamsBodyRunsAnnotationQueueRunAddSchemaArrayItem) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueRunNewParamsBodyArray []AnnotationQueueRunNewParamsBodyArrayItem

func (r AnnotationQueueRunNewParamsBodyArray) implementsAnnotationQueueRunNewParamsBodyUnion() {}

// Deprecated: use plain UUID list or AddRunToQueueByKeyRequest instead.
type AnnotationQueueRunNewParamsBodyArrayItem struct {
	// Deprecated: deprecated
	RunID param.Field[string] `json:"run_id" api:"required" format:"uuid"`
	// Deprecated: deprecated
	ParentRunID param.Field[string] `json:"parent_run_id" format:"uuid"`
	// Deprecated: deprecated
	SessionID param.Field[string] `json:"session_id" format:"uuid"`
	// Deprecated: deprecated
	StartTime param.Field[time.Time] `json:"start_time" format:"date-time"`
	// Deprecated: deprecated
	TraceID param.Field[string] `json:"trace_id" format:"uuid"`
	// Deprecated: deprecated
	TraceTier param.Field[AnnotationQueueRunNewParamsBodyArrayTraceTier] `json:"trace_tier"`
}

func (r AnnotationQueueRunNewParamsBodyArrayItem) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueRunNewParamsBodyArrayTraceTier string

const (
	AnnotationQueueRunNewParamsBodyArrayTraceTierLonglived  AnnotationQueueRunNewParamsBodyArrayTraceTier = "longlived"
	AnnotationQueueRunNewParamsBodyArrayTraceTierShortlived AnnotationQueueRunNewParamsBodyArrayTraceTier = "shortlived"
)

func (r AnnotationQueueRunNewParamsBodyArrayTraceTier) IsKnown() bool {
	switch r {
	case AnnotationQueueRunNewParamsBodyArrayTraceTierLonglived, AnnotationQueueRunNewParamsBodyArrayTraceTierShortlived:
		return true
	}
	return false
}

type AnnotationQueueRunUpdateParams struct {
	AddedAt          param.Field[time.Time] `json:"added_at" format:"date-time"`
	LastReviewedTime param.Field[time.Time] `json:"last_reviewed_time" format:"date-time"`
}

func (r AnnotationQueueRunUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueRunListParams struct {
	Archived     param.Field[bool]                               `query:"archived"`
	IncludeStats param.Field[bool]                               `query:"include_stats"`
	Limit        param.Field[int64]                              `query:"limit"`
	Offset       param.Field[int64]                              `query:"offset"`
	Status       param.Field[AnnotationQueueRunListParamsStatus] `query:"status"`
}

// URLQuery serializes [AnnotationQueueRunListParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueRunListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AnnotationQueueRunListParamsStatus string

const (
	AnnotationQueueRunListParamsStatusNeedsMyReview     AnnotationQueueRunListParamsStatus = "needs_my_review"
	AnnotationQueueRunListParamsStatusNeedsOthersReview AnnotationQueueRunListParamsStatus = "needs_others_review"
	AnnotationQueueRunListParamsStatusCompleted         AnnotationQueueRunListParamsStatus = "completed"
)

func (r AnnotationQueueRunListParamsStatus) IsKnown() bool {
	switch r {
	case AnnotationQueueRunListParamsStatusNeedsMyReview, AnnotationQueueRunListParamsStatusNeedsOthersReview, AnnotationQueueRunListParamsStatusCompleted:
		return true
	}
	return false
}

type AnnotationQueueRunDeleteAllParams struct {
	DeleteAll     param.Field[bool]     `json:"delete_all"`
	ExcludeRunIDs param.Field[[]string] `json:"exclude_run_ids" format:"uuid"`
	RunIDs        param.Field[[]string] `json:"run_ids" format:"uuid"`
}

func (r AnnotationQueueRunDeleteAllParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
