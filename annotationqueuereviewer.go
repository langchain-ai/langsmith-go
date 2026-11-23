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
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// AnnotationQueueReviewerService contains methods and other services that help
// with interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAnnotationQueueReviewerService] method instead.
type AnnotationQueueReviewerService struct {
	Options []option.RequestOption
}

// NewAnnotationQueueReviewerService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAnnotationQueueReviewerService(opts ...option.RequestOption) (r *AnnotationQueueReviewerService) {
	r = &AnnotationQueueReviewerService{}
	r.Options = opts
	return
}

// List a queue's reviewers with how many reviews each has waiting, how many they
// have completed, and when they last completed one. Includes everyone on the
// queue's roster and anyone else who has reviewed on the queue, such as people
// removed from the roster; is_on_roster tells them apart. name and email are left
// out when a reviewer cannot be looked up. Only single queues are supported, and
// the endpoint returns 501 unless multi-reviewer is on for the workspace.
func (r *AnnotationQueueReviewerService) List(ctx context.Context, queueID string, query AnnotationQueueReviewerListParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[AnnotationQueueReviewerListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/platform/annotation-queues/%s/reviewers", queueID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodGet, path, query, &res, opts...)
	if err != nil {
		return nil, err
	}
	err = cfg.Execute()
	if err != nil {
		return nil, err
	}
	res.SetPageConfig(cfg, raw)
	return res, nil
}

// List a queue's reviewers with how many reviews each has waiting, how many they
// have completed, and when they last completed one. Includes everyone on the
// queue's roster and anyone else who has reviewed on the queue, such as people
// removed from the roster; is_on_roster tells them apart. name and email are left
// out when a reviewer cannot be looked up. Only single queues are supported, and
// the endpoint returns 501 unless multi-reviewer is on for the workspace.
func (r *AnnotationQueueReviewerService) ListAutoPaging(ctx context.Context, queueID string, query AnnotationQueueReviewerListParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[AnnotationQueueReviewerListResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.List(ctx, queueID, query, opts...))
}

type AnnotationQueueReviewerListResponse struct {
	CompletedCount int64                                   `json:"completed_count"`
	Email          string                                  `json:"email"`
	IdentityID     string                                  `json:"identity_id"`
	IsOnRoster     bool                                    `json:"is_on_roster"`
	LastActiveAt   time.Time                               `json:"last_active_at" format:"date-time"`
	Name           string                                  `json:"name"`
	PendingCount   int64                                   `json:"pending_count"`
	JSON           annotationQueueReviewerListResponseJSON `json:"-"`
}

// annotationQueueReviewerListResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueReviewerListResponse]
type annotationQueueReviewerListResponseJSON struct {
	CompletedCount apijson.Field
	Email          apijson.Field
	IdentityID     apijson.Field
	IsOnRoster     apijson.Field
	LastActiveAt   apijson.Field
	Name           apijson.Field
	PendingCount   apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *AnnotationQueueReviewerListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueReviewerListResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueReviewerListParams struct {
	// Opaque pagination cursor
	Cursor param.Field[string] `query:"cursor"`
	// Page size (max 100)
	PageSize param.Field[int64] `query:"page_size"`
}

// URLQuery serializes [AnnotationQueueReviewerListParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueReviewerListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
