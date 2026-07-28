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

// AnnotationQueueItemService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAnnotationQueueItemService] method instead.
type AnnotationQueueItemService struct {
	Options []option.RequestOption
}

// NewAnnotationQueueItemService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewAnnotationQueueItemService(opts ...option.RequestOption) (r *AnnotationQueueItemService) {
	r = &AnnotationQueueItemService{}
	r.Options = opts
	return
}

// Add RUN or THREAD items to a single annotation queue. RUN items require run_id
// unless they are created from a suggested example. THREAD items require thread_id
// and session_id.
func (r *AnnotationQueueItemService) New(ctx context.Context, queueID string, params AnnotationQueueItemNewParams, opts ...option.RequestOption) (res *AnnotationQueueItemNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/items", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, params, &res, opts...)
	return res, err
}

// Partially update mutable timestamps (added_at, last_reviewed_time) for a RUN or
// THREAD annotation queue item. Omit a field, or pass JSON null, to leave it
// unchanged.
func (r *AnnotationQueueItemService) Update(ctx context.Context, queueID string, itemID string, body AnnotationQueueItemUpdateParams, opts ...option.RequestOption) (res *AnnotationQueueItemUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	if itemID == "" {
		err = errors.New("missing required item_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/items/%s", queueID, itemID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List RUN and THREAD items in a single annotation queue for one review status
// section, with opaque cursor pagination. Optional item_type=RUN|THREAD filters
// the page. direction=backward returns items before the supplied cursor. The
// response contains item metadata only, not expanded run or thread payloads.
// status=archived returns items whose queue review requirements have been
// satisfied, not merely items the caller personally marked completed.
func (r *AnnotationQueueItemService) List(ctx context.Context, queueID string, query AnnotationQueueItemListParams, opts ...option.RequestOption) (res *pagination.ItemsCursorGetPagination[AnnotationQueueItemListResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/items", queueID)
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

// List RUN and THREAD items in a single annotation queue for one review status
// section, with opaque cursor pagination. Optional item_type=RUN|THREAD filters
// the page. direction=backward returns items before the supplied cursor. The
// response contains item metadata only, not expanded run or thread payloads.
// status=archived returns items whose queue review requirements have been
// satisfied, not merely items the caller personally marked completed.
func (r *AnnotationQueueItemService) ListAutoPaging(ctx context.Context, queueID string, query AnnotationQueueItemListParams, opts ...option.RequestOption) *pagination.ItemsCursorGetPaginationAutoPager[AnnotationQueueItemListResponse] {
	return pagination.NewItemsCursorGetPaginationAutoPager(r.List(ctx, queueID, query, opts...))
}

// Log the caller's reviewer status for a RUN or THREAD annotation queue item. A
// null status re-shows the item for this reviewer.
func (r *AnnotationQueueItemService) NewStatus(ctx context.Context, queueItemID string, body AnnotationQueueItemNewStatusParams, opts ...option.RequestOption) (res *AnnotationQueueItemNewStatusResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueItemID == "" {
		err = errors.New("missing required queue_item_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/items/%s/status", queueItemID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Remove RUN or THREAD items from a single annotation queue by item ID.
func (r *AnnotationQueueItemService) DeleteAll(ctx context.Context, queueID string, body AnnotationQueueItemDeleteAllParams, opts ...option.RequestOption) (res *AnnotationQueueItemDeleteAllResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/items/delete", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Returns the number of annotation queue items for the requested reviewer-specific
// or archived bucket.
func (r *AnnotationQueueItemService) GetCount(ctx context.Context, queueID string, query AnnotationQueueItemGetCountParams, opts ...option.RequestOption) (res *AnnotationQueueItemGetCountResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/items/count", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Resolve a RUN or THREAD item to its current review section and zero-based
// position for deep linking.
func (r *AnnotationQueueItemService) GetPlacement(ctx context.Context, queueID string, itemID string, opts ...option.RequestOption) (res *AnnotationQueueItemGetPlacementResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	if itemID == "" {
		err = errors.New("missing required item_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/items/%s/placement", queueID, itemID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type AnnotationQueueItemNewResponse struct {
	Items []AnnotationQueueItemNewResponseItem `json:"items"`
	JSON  annotationQueueItemNewResponseJSON   `json:"-"`
}

// annotationQueueItemNewResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueItemNewResponse]
type annotationQueueItemNewResponseJSON struct {
	Items       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueItemNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemNewResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemNewResponseItem struct {
	ID       string                                      `json:"id"`
	AddedAt  string                                      `json:"added_at"`
	ItemType AnnotationQueueItemNewResponseItemsItemType `json:"item_type"`
	// LastReviewedTime is always present on the wire (null until reviewed).
	LastReviewedTime        string                                 `json:"last_reviewed_time"`
	QueueID                 string                                 `json:"queue_id"`
	RunID                   string                                 `json:"run_id"`
	SessionID               string                                 `json:"session_id"`
	SourceProposedExampleID string                                 `json:"source_proposed_example_id"`
	StartTime               string                                 `json:"start_time"`
	ThreadID                string                                 `json:"thread_id"`
	JSON                    annotationQueueItemNewResponseItemJSON `json:"-"`
}

// annotationQueueItemNewResponseItemJSON contains the JSON metadata for the struct
// [AnnotationQueueItemNewResponseItem]
type annotationQueueItemNewResponseItemJSON struct {
	ID                      apijson.Field
	AddedAt                 apijson.Field
	ItemType                apijson.Field
	LastReviewedTime        apijson.Field
	QueueID                 apijson.Field
	RunID                   apijson.Field
	SessionID               apijson.Field
	SourceProposedExampleID apijson.Field
	StartTime               apijson.Field
	ThreadID                apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *AnnotationQueueItemNewResponseItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemNewResponseItemJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemNewResponseItemsItemType string

const (
	AnnotationQueueItemNewResponseItemsItemTypeRun    AnnotationQueueItemNewResponseItemsItemType = "RUN"
	AnnotationQueueItemNewResponseItemsItemTypeThread AnnotationQueueItemNewResponseItemsItemType = "THREAD"
)

func (r AnnotationQueueItemNewResponseItemsItemType) IsKnown() bool {
	switch r {
	case AnnotationQueueItemNewResponseItemsItemTypeRun, AnnotationQueueItemNewResponseItemsItemTypeThread:
		return true
	}
	return false
}

type AnnotationQueueItemUpdateResponse struct {
	ID       string                                    `json:"id"`
	AddedAt  string                                    `json:"added_at"`
	ItemType AnnotationQueueItemUpdateResponseItemType `json:"item_type"`
	// LastReviewedTime is always present on the wire (null until reviewed).
	LastReviewedTime        string                                `json:"last_reviewed_time"`
	QueueID                 string                                `json:"queue_id"`
	RunID                   string                                `json:"run_id"`
	SessionID               string                                `json:"session_id"`
	SourceProposedExampleID string                                `json:"source_proposed_example_id"`
	StartTime               string                                `json:"start_time"`
	ThreadID                string                                `json:"thread_id"`
	JSON                    annotationQueueItemUpdateResponseJSON `json:"-"`
}

// annotationQueueItemUpdateResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueItemUpdateResponse]
type annotationQueueItemUpdateResponseJSON struct {
	ID                      apijson.Field
	AddedAt                 apijson.Field
	ItemType                apijson.Field
	LastReviewedTime        apijson.Field
	QueueID                 apijson.Field
	RunID                   apijson.Field
	SessionID               apijson.Field
	SourceProposedExampleID apijson.Field
	StartTime               apijson.Field
	ThreadID                apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *AnnotationQueueItemUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemUpdateResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemUpdateResponseItemType string

const (
	AnnotationQueueItemUpdateResponseItemTypeRun    AnnotationQueueItemUpdateResponseItemType = "RUN"
	AnnotationQueueItemUpdateResponseItemTypeThread AnnotationQueueItemUpdateResponseItemType = "THREAD"
)

func (r AnnotationQueueItemUpdateResponseItemType) IsKnown() bool {
	switch r {
	case AnnotationQueueItemUpdateResponseItemTypeRun, AnnotationQueueItemUpdateResponseItemTypeThread:
		return true
	}
	return false
}

type AnnotationQueueItemListResponse struct {
	ID               string                                  `json:"id"`
	AddedAt          string                                  `json:"added_at"`
	CompletedBy      []string                                `json:"completed_by"`
	EffectiveAddedAt string                                  `json:"effective_added_at"`
	ItemType         AnnotationQueueItemListResponseItemType `json:"item_type"`
	// LastReviewedTime is always present on the wire (null until reviewed).
	LastReviewedTime        string                              `json:"last_reviewed_time"`
	QueueID                 string                              `json:"queue_id"`
	ReservedBy              []string                            `json:"reserved_by"`
	RunID                   string                              `json:"run_id"`
	SessionID               string                              `json:"session_id"`
	SourceProposedExampleID string                              `json:"source_proposed_example_id"`
	StartTime               string                              `json:"start_time"`
	ThreadID                string                              `json:"thread_id"`
	JSON                    annotationQueueItemListResponseJSON `json:"-"`
}

// annotationQueueItemListResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueItemListResponse]
type annotationQueueItemListResponseJSON struct {
	ID                      apijson.Field
	AddedAt                 apijson.Field
	CompletedBy             apijson.Field
	EffectiveAddedAt        apijson.Field
	ItemType                apijson.Field
	LastReviewedTime        apijson.Field
	QueueID                 apijson.Field
	ReservedBy              apijson.Field
	RunID                   apijson.Field
	SessionID               apijson.Field
	SourceProposedExampleID apijson.Field
	StartTime               apijson.Field
	ThreadID                apijson.Field
	raw                     string
	ExtraFields             map[string]apijson.Field
}

func (r *AnnotationQueueItemListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemListResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemListResponseItemType string

const (
	AnnotationQueueItemListResponseItemTypeRun    AnnotationQueueItemListResponseItemType = "RUN"
	AnnotationQueueItemListResponseItemTypeThread AnnotationQueueItemListResponseItemType = "THREAD"
)

func (r AnnotationQueueItemListResponseItemType) IsKnown() bool {
	switch r {
	case AnnotationQueueItemListResponseItemTypeRun, AnnotationQueueItemListResponseItemTypeThread:
		return true
	}
	return false
}

type AnnotationQueueItemNewStatusResponse struct {
	IsArchived      bool                                       `json:"is_archived"`
	OverrideAddedAt string                                     `json:"override_added_at"`
	QueueItemID     string                                     `json:"queue_item_id"`
	Status          AnnotationQueueItemNewStatusResponseStatus `json:"status"`
	JSON            annotationQueueItemNewStatusResponseJSON   `json:"-"`
}

// annotationQueueItemNewStatusResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueItemNewStatusResponse]
type annotationQueueItemNewStatusResponseJSON struct {
	IsArchived      apijson.Field
	OverrideAddedAt apijson.Field
	QueueItemID     apijson.Field
	Status          apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *AnnotationQueueItemNewStatusResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemNewStatusResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemNewStatusResponseStatus string

const (
	AnnotationQueueItemNewStatusResponseStatusViewed    AnnotationQueueItemNewStatusResponseStatus = "viewed"
	AnnotationQueueItemNewStatusResponseStatusCompleted AnnotationQueueItemNewStatusResponseStatus = "completed"
)

func (r AnnotationQueueItemNewStatusResponseStatus) IsKnown() bool {
	switch r {
	case AnnotationQueueItemNewStatusResponseStatusViewed, AnnotationQueueItemNewStatusResponseStatusCompleted:
		return true
	}
	return false
}

type AnnotationQueueItemDeleteAllResponse map[string]string

type AnnotationQueueItemGetCountResponse struct {
	Count int64                                   `json:"count"`
	JSON  annotationQueueItemGetCountResponseJSON `json:"-"`
}

// annotationQueueItemGetCountResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueItemGetCountResponse]
type annotationQueueItemGetCountResponseJSON struct {
	Count       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueItemGetCountResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemGetCountResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemGetPlacementResponse struct {
	Cursor   string                                          `json:"cursor"`
	ItemType AnnotationQueueItemGetPlacementResponseItemType `json:"item_type"`
	Position int64                                           `json:"position"`
	Section  AnnotationQueueItemGetPlacementResponseSection  `json:"section"`
	JSON     annotationQueueItemGetPlacementResponseJSON     `json:"-"`
}

// annotationQueueItemGetPlacementResponseJSON contains the JSON metadata for the
// struct [AnnotationQueueItemGetPlacementResponse]
type annotationQueueItemGetPlacementResponseJSON struct {
	Cursor      apijson.Field
	ItemType    apijson.Field
	Position    apijson.Field
	Section     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueItemGetPlacementResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueItemGetPlacementResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueItemGetPlacementResponseItemType string

const (
	AnnotationQueueItemGetPlacementResponseItemTypeRun    AnnotationQueueItemGetPlacementResponseItemType = "RUN"
	AnnotationQueueItemGetPlacementResponseItemTypeThread AnnotationQueueItemGetPlacementResponseItemType = "THREAD"
)

func (r AnnotationQueueItemGetPlacementResponseItemType) IsKnown() bool {
	switch r {
	case AnnotationQueueItemGetPlacementResponseItemTypeRun, AnnotationQueueItemGetPlacementResponseItemTypeThread:
		return true
	}
	return false
}

type AnnotationQueueItemGetPlacementResponseSection string

const (
	AnnotationQueueItemGetPlacementResponseSectionNeedsMyReview     AnnotationQueueItemGetPlacementResponseSection = "needs_my_review"
	AnnotationQueueItemGetPlacementResponseSectionNeedsOthersReview AnnotationQueueItemGetPlacementResponseSection = "needs_others_review"
	AnnotationQueueItemGetPlacementResponseSectionArchived          AnnotationQueueItemGetPlacementResponseSection = "archived"
)

func (r AnnotationQueueItemGetPlacementResponseSection) IsKnown() bool {
	switch r {
	case AnnotationQueueItemGetPlacementResponseSectionNeedsMyReview, AnnotationQueueItemGetPlacementResponseSectionNeedsOthersReview, AnnotationQueueItemGetPlacementResponseSectionArchived:
		return true
	}
	return false
}

type AnnotationQueueItemNewParams struct {
	// Extend trace retention for added run items
	ExtendTraceRetention param.Field[bool]                               `query:"extend_trace_retention"`
	Items                param.Field[[]AnnotationQueueItemNewParamsItem] `json:"items"`
}

func (r AnnotationQueueItemNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// URLQuery serializes [AnnotationQueueItemNewParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueItemNewParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AnnotationQueueItemNewParamsItem struct {
	ItemType param.Field[AnnotationQueueItemNewParamsItemsItemType] `json:"item_type"`
	// RUN fields
	RunID param.Field[string] `json:"run_id"`
	// SessionID is the ID of the tracing project that contains the run or thread.
	SessionID param.Field[string] `json:"session_id"`
	// SourceProposedExampleID links the queue item to the suggested example it was
	// created from, when applicable.
	SourceProposedExampleID param.Field[string] `json:"source_proposed_example_id"`
	// StartTime is the start time of the run being added, used to identify it.
	StartTime param.Field[time.Time] `json:"start_time" format:"date-time"`
	// ThreadID is the ID of the thread being added.
	ThreadID param.Field[string] `json:"thread_id"`
}

func (r AnnotationQueueItemNewParamsItem) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueItemNewParamsItemsItemType string

const (
	AnnotationQueueItemNewParamsItemsItemTypeRun    AnnotationQueueItemNewParamsItemsItemType = "RUN"
	AnnotationQueueItemNewParamsItemsItemTypeThread AnnotationQueueItemNewParamsItemsItemType = "THREAD"
)

func (r AnnotationQueueItemNewParamsItemsItemType) IsKnown() bool {
	switch r {
	case AnnotationQueueItemNewParamsItemsItemTypeRun, AnnotationQueueItemNewParamsItemsItemTypeThread:
		return true
	}
	return false
}

type AnnotationQueueItemUpdateParams struct {
	AddedAt          param.Field[time.Time] `json:"added_at" format:"date-time"`
	LastReviewedTime param.Field[time.Time] `json:"last_reviewed_time" format:"date-time"`
}

func (r AnnotationQueueItemUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueItemListParams struct {
	// Review section: needs_my_review, needs_others_review, or archived
	Status param.Field[AnnotationQueueItemListParamsStatus] `query:"status" api:"required"`
	// Opaque pagination cursor
	Cursor param.Field[string] `query:"cursor"`
	// Pagination direction. backward requires cursor
	Direction param.Field[AnnotationQueueItemListParamsDirection] `query:"direction"`
	// Filter to RUN or THREAD
	ItemType param.Field[AnnotationQueueItemListParamsItemType] `query:"item_type"`
	// Page size (max 100)
	PageSize param.Field[int64] `query:"page_size"`
}

// URLQuery serializes [AnnotationQueueItemListParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueItemListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// Review section: needs_my_review, needs_others_review, or archived
type AnnotationQueueItemListParamsStatus string

const (
	AnnotationQueueItemListParamsStatusNeedsMyReview     AnnotationQueueItemListParamsStatus = "needs_my_review"
	AnnotationQueueItemListParamsStatusNeedsOthersReview AnnotationQueueItemListParamsStatus = "needs_others_review"
	AnnotationQueueItemListParamsStatusArchived          AnnotationQueueItemListParamsStatus = "archived"
)

func (r AnnotationQueueItemListParamsStatus) IsKnown() bool {
	switch r {
	case AnnotationQueueItemListParamsStatusNeedsMyReview, AnnotationQueueItemListParamsStatusNeedsOthersReview, AnnotationQueueItemListParamsStatusArchived:
		return true
	}
	return false
}

// Pagination direction. backward requires cursor
type AnnotationQueueItemListParamsDirection string

const (
	AnnotationQueueItemListParamsDirectionForward  AnnotationQueueItemListParamsDirection = "forward"
	AnnotationQueueItemListParamsDirectionBackward AnnotationQueueItemListParamsDirection = "backward"
)

func (r AnnotationQueueItemListParamsDirection) IsKnown() bool {
	switch r {
	case AnnotationQueueItemListParamsDirectionForward, AnnotationQueueItemListParamsDirectionBackward:
		return true
	}
	return false
}

// Filter to RUN or THREAD
type AnnotationQueueItemListParamsItemType string

const (
	AnnotationQueueItemListParamsItemTypeRun    AnnotationQueueItemListParamsItemType = "RUN"
	AnnotationQueueItemListParamsItemTypeThread AnnotationQueueItemListParamsItemType = "THREAD"
)

func (r AnnotationQueueItemListParamsItemType) IsKnown() bool {
	switch r {
	case AnnotationQueueItemListParamsItemTypeRun, AnnotationQueueItemListParamsItemTypeThread:
		return true
	}
	return false
}

type AnnotationQueueItemNewStatusParams struct {
	OverrideAddedAt param.Field[string]                                   `json:"override_added_at"`
	Status          param.Field[AnnotationQueueItemNewStatusParamsStatus] `json:"status"`
}

func (r AnnotationQueueItemNewStatusParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueItemNewStatusParamsStatus string

const (
	AnnotationQueueItemNewStatusParamsStatusViewed    AnnotationQueueItemNewStatusParamsStatus = "viewed"
	AnnotationQueueItemNewStatusParamsStatusCompleted AnnotationQueueItemNewStatusParamsStatus = "completed"
)

func (r AnnotationQueueItemNewStatusParamsStatus) IsKnown() bool {
	switch r {
	case AnnotationQueueItemNewStatusParamsStatusViewed, AnnotationQueueItemNewStatusParamsStatusCompleted:
		return true
	}
	return false
}

type AnnotationQueueItemDeleteAllParams struct {
	ItemIDs param.Field[[]string] `json:"item_ids"`
}

func (r AnnotationQueueItemDeleteAllParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueItemGetCountParams struct {
	// Count bucket: all, needs_my_review, needs_others_review, or archived.
	Status param.Field[string] `query:"status" api:"required"`
	// Exclusive upper bound for archived item timestamp
	EndTime param.Field[string] `query:"end_time"`
	// Exclusive lower bound for archived item timestamp
	StartTime param.Field[string] `query:"start_time"`
}

// URLQuery serializes [AnnotationQueueItemGetCountParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueItemGetCountParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
