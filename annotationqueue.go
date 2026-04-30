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

// AnnotationQueueService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewAnnotationQueueService] method instead.
type AnnotationQueueService struct {
	Options    []option.RequestOption
	Runs       *AnnotationQueueRunService
	Info       *AnnotationQueueInfoService
	Workspaces *AnnotationQueueWorkspaceService
}

// NewAnnotationQueueService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAnnotationQueueService(opts ...option.RequestOption) (r *AnnotationQueueService) {
	r = &AnnotationQueueService{}
	r.Options = opts
	r.Runs = NewAnnotationQueueRunService(opts...)
	r.Info = NewAnnotationQueueInfoService(opts...)
	r.Workspaces = NewAnnotationQueueWorkspaceService(opts...)
	return
}

// Get Annotation Queue
func (r *AnnotationQueueService) Get(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Update Annotation Queue
func (r *AnnotationQueueService) Update(ctx context.Context, queueID string, body AnnotationQueueUpdateParams, opts ...option.RequestOption) (res *AnnotationQueueUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Delete Annotation Queue
func (r *AnnotationQueueService) Delete(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// Create Annotation Queue
func (r *AnnotationQueueService) AnnotationQueues(ctx context.Context, body AnnotationQueueAnnotationQueuesParams, opts ...option.RequestOption) (res *AnnotationQueueSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/annotation-queues"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Create Identity Annotation Queue Run Status
func (r *AnnotationQueueService) NewRunStatus(ctx context.Context, annotationQueueRunID string, body AnnotationQueueNewRunStatusParams, opts ...option.RequestOption) (res *AnnotationQueueNewRunStatusResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if annotationQueueRunID == "" {
		err = errors.New("missing required annotation_queue_run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/status/%s", annotationQueueRunID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Export Annotation Queue Archived Runs
func (r *AnnotationQueueService) Export(ctx context.Context, queueID string, body AnnotationQueueExportParams, opts ...option.RequestOption) (res *AnnotationQueueExportResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/export", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Populate annotation queue with runs from an experiment.
func (r *AnnotationQueueService) Populate(ctx context.Context, body AnnotationQueuePopulateParams, opts ...option.RequestOption) (res *AnnotationQueuePopulateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/annotation-queues/populate"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Get Annotation Queues
func (r *AnnotationQueueService) GetAnnotationQueues(ctx context.Context, query AnnotationQueueGetAnnotationQueuesParams, opts ...option.RequestOption) (res *pagination.OffsetPaginationTopLevelArray[AnnotationQueueGetAnnotationQueuesResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	path := "api/v1/annotation-queues"
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

// Get Annotation Queues
func (r *AnnotationQueueService) GetAnnotationQueuesAutoPaging(ctx context.Context, query AnnotationQueueGetAnnotationQueuesParams, opts ...option.RequestOption) *pagination.OffsetPaginationTopLevelArrayAutoPager[AnnotationQueueGetAnnotationQueuesResponse] {
	return pagination.NewOffsetPaginationTopLevelArrayAutoPager(r.GetAnnotationQueues(ctx, query, opts...))
}

// Get Annotation Queues For Run
func (r *AnnotationQueueService) GetQueues(ctx context.Context, runID string, opts ...option.RequestOption) (res *[]AnnotationQueueSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if runID == "" {
		err = errors.New("missing required run_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/queues", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get a run from an annotation queue
func (r *AnnotationQueueService) GetRun(ctx context.Context, queueID string, index int64, query AnnotationQueueGetRunParams, opts ...option.RequestOption) (res *RunSchemaWithAnnotationQueueInfo, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/run/%v", queueID, index)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get Size From Annotation Queue
func (r *AnnotationQueueService) GetSize(ctx context.Context, queueID string, query AnnotationQueueGetSizeParams, opts ...option.RequestOption) (res *AnnotationQueueSizeSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/size", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get Total Archived From Annotation Queue
func (r *AnnotationQueueService) GetTotalArchived(ctx context.Context, queueID string, query AnnotationQueueGetTotalArchivedParams, opts ...option.RequestOption) (res *AnnotationQueueSizeSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/total_archived", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Get Total Size From Annotation Queue
func (r *AnnotationQueueService) GetTotalSize(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueSizeSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/total_size", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

type AnnotationQueueRubricItemSchema struct {
	FeedbackKey       string                              `json:"feedback_key" api:"required"`
	Description       string                              `json:"description" api:"nullable"`
	IsAssertion       bool                                `json:"is_assertion" api:"nullable"`
	IsRequired        bool                                `json:"is_required" api:"nullable"`
	ScoreDescriptions map[string]string                   `json:"score_descriptions" api:"nullable"`
	ValueDescriptions map[string]string                   `json:"value_descriptions" api:"nullable"`
	JSON              annotationQueueRubricItemSchemaJSON `json:"-"`
}

// annotationQueueRubricItemSchemaJSON contains the JSON metadata for the struct
// [AnnotationQueueRubricItemSchema]
type annotationQueueRubricItemSchemaJSON struct {
	FeedbackKey       apijson.Field
	Description       apijson.Field
	IsAssertion       apijson.Field
	IsRequired        apijson.Field
	ScoreDescriptions apijson.Field
	ValueDescriptions apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *AnnotationQueueRubricItemSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueRubricItemSchemaJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueRubricItemSchemaParam struct {
	FeedbackKey       param.Field[string]            `json:"feedback_key" api:"required"`
	Description       param.Field[string]            `json:"description"`
	IsAssertion       param.Field[bool]              `json:"is_assertion"`
	IsRequired        param.Field[bool]              `json:"is_required"`
	ScoreDescriptions param.Field[map[string]string] `json:"score_descriptions"`
	ValueDescriptions param.Field[map[string]string] `json:"value_descriptions"`
}

func (r AnnotationQueueRubricItemSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// AnnotationQueue schema.
type AnnotationQueueSchema struct {
	ID                  string                                  `json:"id" api:"required" format:"uuid"`
	Name                string                                  `json:"name" api:"required"`
	QueueType           AnnotationQueueSchemaQueueType          `json:"queue_type" api:"required"`
	TenantID            string                                  `json:"tenant_id" api:"required" format:"uuid"`
	AssignedReviewers   []AnnotationQueueSchemaAssignedReviewer `json:"assigned_reviewers"`
	CreatedAt           time.Time                               `json:"created_at" format:"date-time"`
	DefaultDataset      string                                  `json:"default_dataset" api:"nullable" format:"uuid"`
	Description         string                                  `json:"description" api:"nullable"`
	EnableReservations  bool                                    `json:"enable_reservations" api:"nullable"`
	Metadata            map[string]interface{}                  `json:"metadata" api:"nullable"`
	NumReviewersPerItem int64                                   `json:"num_reviewers_per_item" api:"nullable"`
	ReservationMinutes  int64                                   `json:"reservation_minutes" api:"nullable"`
	ReviewerAccessMode  string                                  `json:"reviewer_access_mode"`
	RunRuleID           string                                  `json:"run_rule_id" api:"nullable" format:"uuid"`
	SourceRuleID        string                                  `json:"source_rule_id" api:"nullable" format:"uuid"`
	UpdatedAt           time.Time                               `json:"updated_at" format:"date-time"`
	JSON                annotationQueueSchemaJSON               `json:"-"`
}

// annotationQueueSchemaJSON contains the JSON metadata for the struct
// [AnnotationQueueSchema]
type annotationQueueSchemaJSON struct {
	ID                  apijson.Field
	Name                apijson.Field
	QueueType           apijson.Field
	TenantID            apijson.Field
	AssignedReviewers   apijson.Field
	CreatedAt           apijson.Field
	DefaultDataset      apijson.Field
	Description         apijson.Field
	EnableReservations  apijson.Field
	Metadata            apijson.Field
	NumReviewersPerItem apijson.Field
	ReservationMinutes  apijson.Field
	ReviewerAccessMode  apijson.Field
	RunRuleID           apijson.Field
	SourceRuleID        apijson.Field
	UpdatedAt           apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *AnnotationQueueSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueSchemaJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueSchemaQueueType string

const (
	AnnotationQueueSchemaQueueTypeSingle   AnnotationQueueSchemaQueueType = "single"
	AnnotationQueueSchemaQueueTypePairwise AnnotationQueueSchemaQueueType = "pairwise"
)

func (r AnnotationQueueSchemaQueueType) IsKnown() bool {
	switch r {
	case AnnotationQueueSchemaQueueTypeSingle, AnnotationQueueSchemaQueueTypePairwise:
		return true
	}
	return false
}

// Identity info for an assigned reviewer on an annotation queue.
type AnnotationQueueSchemaAssignedReviewer struct {
	ID    string                                    `json:"id" api:"required" format:"uuid"`
	Email string                                    `json:"email" api:"nullable"`
	Name  string                                    `json:"name" api:"nullable"`
	JSON  annotationQueueSchemaAssignedReviewerJSON `json:"-"`
}

// annotationQueueSchemaAssignedReviewerJSON contains the JSON metadata for the
// struct [AnnotationQueueSchemaAssignedReviewer]
type annotationQueueSchemaAssignedReviewerJSON struct {
	ID          apijson.Field
	Email       apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueSchemaAssignedReviewer) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueSchemaAssignedReviewerJSON) RawJSON() string {
	return r.raw
}

// Size of an Annotation Queue
type AnnotationQueueSizeSchema struct {
	Size int64                         `json:"size" api:"required"`
	JSON annotationQueueSizeSchemaJSON `json:"-"`
}

// annotationQueueSizeSchemaJSON contains the JSON metadata for the struct
// [AnnotationQueueSizeSchema]
type annotationQueueSizeSchemaJSON struct {
	Size        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueSizeSchema) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueSizeSchemaJSON) RawJSON() string {
	return r.raw
}

// Run schema with annotation queue info.
type RunSchemaWithAnnotationQueueInfo struct {
	ID          string `json:"id" api:"required" format:"uuid"`
	AppPath     string `json:"app_path" api:"required"`
	DottedOrder string `json:"dotted_order" api:"required"`
	Name        string `json:"name" api:"required"`
	QueueRunID  string `json:"queue_run_id" api:"required" format:"uuid"`
	// Enum for run types.
	RunType                RunTypeEnum                               `json:"run_type" api:"required"`
	SessionID              string                                    `json:"session_id" api:"required" format:"uuid"`
	Status                 string                                    `json:"status" api:"required"`
	TraceID                string                                    `json:"trace_id" api:"required" format:"uuid"`
	AddedAt                time.Time                                 `json:"added_at" api:"nullable" format:"date-time"`
	ChildRunIDs            []string                                  `json:"child_run_ids" api:"nullable" format:"uuid"`
	CompletedBy            []string                                  `json:"completed_by" format:"uuid"`
	CompletionCost         string                                    `json:"completion_cost" api:"nullable"`
	CompletionCostDetails  map[string]string                         `json:"completion_cost_details" api:"nullable"`
	CompletionTokenDetails map[string]int64                          `json:"completion_token_details" api:"nullable"`
	CompletionTokens       int64                                     `json:"completion_tokens"`
	DirectChildRunIDs      []string                                  `json:"direct_child_run_ids" api:"nullable" format:"uuid"`
	EffectiveAddedAt       time.Time                                 `json:"effective_added_at" api:"nullable" format:"date-time"`
	EndTime                time.Time                                 `json:"end_time" api:"nullable" format:"date-time"`
	Error                  string                                    `json:"error" api:"nullable"`
	Events                 []map[string]interface{}                  `json:"events" api:"nullable"`
	ExecutionOrder         int64                                     `json:"execution_order"`
	Extra                  map[string]interface{}                    `json:"extra" api:"nullable"`
	FeedbackStats          map[string]map[string]interface{}         `json:"feedback_stats" api:"nullable"`
	FirstTokenTime         time.Time                                 `json:"first_token_time" api:"nullable" format:"date-time"`
	InDataset              bool                                      `json:"in_dataset" api:"nullable"`
	Inputs                 map[string]interface{}                    `json:"inputs" api:"nullable"`
	InputsPreview          string                                    `json:"inputs_preview" api:"nullable"`
	InputsS3URLs           map[string]interface{}                    `json:"inputs_s3_urls" api:"nullable"`
	LastQueuedAt           time.Time                                 `json:"last_queued_at" api:"nullable" format:"date-time"`
	LastReviewedTime       time.Time                                 `json:"last_reviewed_time" api:"nullable" format:"date-time"`
	ManifestID             string                                    `json:"manifest_id" api:"nullable" format:"uuid"`
	ManifestS3ID           string                                    `json:"manifest_s3_id" api:"nullable" format:"uuid"`
	Messages               []map[string]interface{}                  `json:"messages" api:"nullable"`
	Outputs                map[string]interface{}                    `json:"outputs" api:"nullable"`
	OutputsPreview         string                                    `json:"outputs_preview" api:"nullable"`
	OutputsS3URLs          map[string]interface{}                    `json:"outputs_s3_urls" api:"nullable"`
	ParentRunID            string                                    `json:"parent_run_id" api:"nullable" format:"uuid"`
	ParentRunIDs           []string                                  `json:"parent_run_ids" api:"nullable" format:"uuid"`
	PriceModelID           string                                    `json:"price_model_id" api:"nullable" format:"uuid"`
	PromptCost             string                                    `json:"prompt_cost" api:"nullable"`
	PromptCostDetails      map[string]string                         `json:"prompt_cost_details" api:"nullable"`
	PromptTokenDetails     map[string]int64                          `json:"prompt_token_details" api:"nullable"`
	PromptTokens           int64                                     `json:"prompt_tokens"`
	ReferenceDatasetID     string                                    `json:"reference_dataset_id" api:"nullable" format:"uuid"`
	ReferenceExampleID     string                                    `json:"reference_example_id" api:"nullable" format:"uuid"`
	ReservedBy             []string                                  `json:"reserved_by" format:"uuid"`
	S3URLs                 map[string]interface{}                    `json:"s3_urls" api:"nullable"`
	Serialized             map[string]interface{}                    `json:"serialized" api:"nullable"`
	ShareToken             string                                    `json:"share_token" api:"nullable" format:"uuid"`
	StartTime              time.Time                                 `json:"start_time" format:"date-time"`
	Tags                   []string                                  `json:"tags" api:"nullable"`
	ThreadID               string                                    `json:"thread_id" api:"nullable"`
	TotalCost              string                                    `json:"total_cost" api:"nullable"`
	TotalTokens            int64                                     `json:"total_tokens"`
	TraceFirstReceivedAt   time.Time                                 `json:"trace_first_received_at" api:"nullable" format:"date-time"`
	TraceMaxStartTime      time.Time                                 `json:"trace_max_start_time" api:"nullable" format:"date-time"`
	TraceMinStartTime      time.Time                                 `json:"trace_min_start_time" api:"nullable" format:"date-time"`
	TraceTier              RunSchemaWithAnnotationQueueInfoTraceTier `json:"trace_tier" api:"nullable"`
	TraceUpgrade           bool                                      `json:"trace_upgrade"`
	TtlSeconds             int64                                     `json:"ttl_seconds" api:"nullable"`
	JSON                   runSchemaWithAnnotationQueueInfoJSON      `json:"-"`
}

// runSchemaWithAnnotationQueueInfoJSON contains the JSON metadata for the struct
// [RunSchemaWithAnnotationQueueInfo]
type runSchemaWithAnnotationQueueInfoJSON struct {
	ID                     apijson.Field
	AppPath                apijson.Field
	DottedOrder            apijson.Field
	Name                   apijson.Field
	QueueRunID             apijson.Field
	RunType                apijson.Field
	SessionID              apijson.Field
	Status                 apijson.Field
	TraceID                apijson.Field
	AddedAt                apijson.Field
	ChildRunIDs            apijson.Field
	CompletedBy            apijson.Field
	CompletionCost         apijson.Field
	CompletionCostDetails  apijson.Field
	CompletionTokenDetails apijson.Field
	CompletionTokens       apijson.Field
	DirectChildRunIDs      apijson.Field
	EffectiveAddedAt       apijson.Field
	EndTime                apijson.Field
	Error                  apijson.Field
	Events                 apijson.Field
	ExecutionOrder         apijson.Field
	Extra                  apijson.Field
	FeedbackStats          apijson.Field
	FirstTokenTime         apijson.Field
	InDataset              apijson.Field
	Inputs                 apijson.Field
	InputsPreview          apijson.Field
	InputsS3URLs           apijson.Field
	LastQueuedAt           apijson.Field
	LastReviewedTime       apijson.Field
	ManifestID             apijson.Field
	ManifestS3ID           apijson.Field
	Messages               apijson.Field
	Outputs                apijson.Field
	OutputsPreview         apijson.Field
	OutputsS3URLs          apijson.Field
	ParentRunID            apijson.Field
	ParentRunIDs           apijson.Field
	PriceModelID           apijson.Field
	PromptCost             apijson.Field
	PromptCostDetails      apijson.Field
	PromptTokenDetails     apijson.Field
	PromptTokens           apijson.Field
	ReferenceDatasetID     apijson.Field
	ReferenceExampleID     apijson.Field
	ReservedBy             apijson.Field
	S3URLs                 apijson.Field
	Serialized             apijson.Field
	ShareToken             apijson.Field
	StartTime              apijson.Field
	Tags                   apijson.Field
	ThreadID               apijson.Field
	TotalCost              apijson.Field
	TotalTokens            apijson.Field
	TraceFirstReceivedAt   apijson.Field
	TraceMaxStartTime      apijson.Field
	TraceMinStartTime      apijson.Field
	TraceTier              apijson.Field
	TraceUpgrade           apijson.Field
	TtlSeconds             apijson.Field
	raw                    string
	ExtraFields            map[string]apijson.Field
}

func (r *RunSchemaWithAnnotationQueueInfo) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r runSchemaWithAnnotationQueueInfoJSON) RawJSON() string {
	return r.raw
}

type RunSchemaWithAnnotationQueueInfoTraceTier string

const (
	RunSchemaWithAnnotationQueueInfoTraceTierLonglived  RunSchemaWithAnnotationQueueInfoTraceTier = "longlived"
	RunSchemaWithAnnotationQueueInfoTraceTierShortlived RunSchemaWithAnnotationQueueInfoTraceTier = "shortlived"
)

func (r RunSchemaWithAnnotationQueueInfoTraceTier) IsKnown() bool {
	switch r {
	case RunSchemaWithAnnotationQueueInfoTraceTierLonglived, RunSchemaWithAnnotationQueueInfoTraceTierShortlived:
		return true
	}
	return false
}

// AnnotationQueue schema with rubric.
type AnnotationQueueGetResponse struct {
	ID                  string                                       `json:"id" api:"required" format:"uuid"`
	Name                string                                       `json:"name" api:"required"`
	QueueType           AnnotationQueueGetResponseQueueType          `json:"queue_type" api:"required"`
	TenantID            string                                       `json:"tenant_id" api:"required" format:"uuid"`
	AssignedReviewers   []AnnotationQueueGetResponseAssignedReviewer `json:"assigned_reviewers"`
	CreatedAt           time.Time                                    `json:"created_at" format:"date-time"`
	DefaultDataset      string                                       `json:"default_dataset" api:"nullable" format:"uuid"`
	Description         string                                       `json:"description" api:"nullable"`
	EnableReservations  bool                                         `json:"enable_reservations" api:"nullable"`
	Metadata            map[string]interface{}                       `json:"metadata" api:"nullable"`
	NumReviewersPerItem int64                                        `json:"num_reviewers_per_item" api:"nullable"`
	ReservationMinutes  int64                                        `json:"reservation_minutes" api:"nullable"`
	ReviewerAccessMode  string                                       `json:"reviewer_access_mode"`
	RubricInstructions  string                                       `json:"rubric_instructions" api:"nullable"`
	RubricItems         []AnnotationQueueRubricItemSchema            `json:"rubric_items" api:"nullable"`
	RunRuleID           string                                       `json:"run_rule_id" api:"nullable" format:"uuid"`
	SourceRuleID        string                                       `json:"source_rule_id" api:"nullable" format:"uuid"`
	UpdatedAt           time.Time                                    `json:"updated_at" format:"date-time"`
	JSON                annotationQueueGetResponseJSON               `json:"-"`
}

// annotationQueueGetResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueGetResponse]
type annotationQueueGetResponseJSON struct {
	ID                  apijson.Field
	Name                apijson.Field
	QueueType           apijson.Field
	TenantID            apijson.Field
	AssignedReviewers   apijson.Field
	CreatedAt           apijson.Field
	DefaultDataset      apijson.Field
	Description         apijson.Field
	EnableReservations  apijson.Field
	Metadata            apijson.Field
	NumReviewersPerItem apijson.Field
	ReservationMinutes  apijson.Field
	ReviewerAccessMode  apijson.Field
	RubricInstructions  apijson.Field
	RubricItems         apijson.Field
	RunRuleID           apijson.Field
	SourceRuleID        apijson.Field
	UpdatedAt           apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *AnnotationQueueGetResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueGetResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueGetResponseQueueType string

const (
	AnnotationQueueGetResponseQueueTypeSingle   AnnotationQueueGetResponseQueueType = "single"
	AnnotationQueueGetResponseQueueTypePairwise AnnotationQueueGetResponseQueueType = "pairwise"
)

func (r AnnotationQueueGetResponseQueueType) IsKnown() bool {
	switch r {
	case AnnotationQueueGetResponseQueueTypeSingle, AnnotationQueueGetResponseQueueTypePairwise:
		return true
	}
	return false
}

// Identity info for an assigned reviewer on an annotation queue.
type AnnotationQueueGetResponseAssignedReviewer struct {
	ID    string                                         `json:"id" api:"required" format:"uuid"`
	Email string                                         `json:"email" api:"nullable"`
	Name  string                                         `json:"name" api:"nullable"`
	JSON  annotationQueueGetResponseAssignedReviewerJSON `json:"-"`
}

// annotationQueueGetResponseAssignedReviewerJSON contains the JSON metadata for
// the struct [AnnotationQueueGetResponseAssignedReviewer]
type annotationQueueGetResponseAssignedReviewerJSON struct {
	ID          apijson.Field
	Email       apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueGetResponseAssignedReviewer) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueGetResponseAssignedReviewerJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueUpdateResponse = interface{}

type AnnotationQueueDeleteResponse = interface{}

type AnnotationQueueNewRunStatusResponse = interface{}

type AnnotationQueueExportResponse = interface{}

type AnnotationQueuePopulateResponse = interface{}

// AnnotationQueue schema with size.
type AnnotationQueueGetAnnotationQueuesResponse struct {
	ID                  string                                                       `json:"id" api:"required" format:"uuid"`
	Name                string                                                       `json:"name" api:"required"`
	QueueType           AnnotationQueueGetAnnotationQueuesResponseQueueType          `json:"queue_type" api:"required"`
	TenantID            string                                                       `json:"tenant_id" api:"required" format:"uuid"`
	TotalRuns           int64                                                        `json:"total_runs" api:"required"`
	AssignedReviewers   []AnnotationQueueGetAnnotationQueuesResponseAssignedReviewer `json:"assigned_reviewers"`
	CreatedAt           time.Time                                                    `json:"created_at" format:"date-time"`
	DefaultDataset      string                                                       `json:"default_dataset" api:"nullable" format:"uuid"`
	Description         string                                                       `json:"description" api:"nullable"`
	EnableReservations  bool                                                         `json:"enable_reservations" api:"nullable"`
	Metadata            map[string]interface{}                                       `json:"metadata" api:"nullable"`
	NumReviewersPerItem int64                                                        `json:"num_reviewers_per_item" api:"nullable"`
	ReservationMinutes  int64                                                        `json:"reservation_minutes" api:"nullable"`
	ReviewerAccessMode  string                                                       `json:"reviewer_access_mode"`
	RunRuleID           string                                                       `json:"run_rule_id" api:"nullable" format:"uuid"`
	SourceRuleID        string                                                       `json:"source_rule_id" api:"nullable" format:"uuid"`
	UpdatedAt           time.Time                                                    `json:"updated_at" format:"date-time"`
	JSON                annotationQueueGetAnnotationQueuesResponseJSON               `json:"-"`
}

// annotationQueueGetAnnotationQueuesResponseJSON contains the JSON metadata for
// the struct [AnnotationQueueGetAnnotationQueuesResponse]
type annotationQueueGetAnnotationQueuesResponseJSON struct {
	ID                  apijson.Field
	Name                apijson.Field
	QueueType           apijson.Field
	TenantID            apijson.Field
	TotalRuns           apijson.Field
	AssignedReviewers   apijson.Field
	CreatedAt           apijson.Field
	DefaultDataset      apijson.Field
	Description         apijson.Field
	EnableReservations  apijson.Field
	Metadata            apijson.Field
	NumReviewersPerItem apijson.Field
	ReservationMinutes  apijson.Field
	ReviewerAccessMode  apijson.Field
	RunRuleID           apijson.Field
	SourceRuleID        apijson.Field
	UpdatedAt           apijson.Field
	raw                 string
	ExtraFields         map[string]apijson.Field
}

func (r *AnnotationQueueGetAnnotationQueuesResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueGetAnnotationQueuesResponseJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueGetAnnotationQueuesResponseQueueType string

const (
	AnnotationQueueGetAnnotationQueuesResponseQueueTypeSingle   AnnotationQueueGetAnnotationQueuesResponseQueueType = "single"
	AnnotationQueueGetAnnotationQueuesResponseQueueTypePairwise AnnotationQueueGetAnnotationQueuesResponseQueueType = "pairwise"
)

func (r AnnotationQueueGetAnnotationQueuesResponseQueueType) IsKnown() bool {
	switch r {
	case AnnotationQueueGetAnnotationQueuesResponseQueueTypeSingle, AnnotationQueueGetAnnotationQueuesResponseQueueTypePairwise:
		return true
	}
	return false
}

// Identity info for an assigned reviewer on an annotation queue.
type AnnotationQueueGetAnnotationQueuesResponseAssignedReviewer struct {
	ID    string                                                         `json:"id" api:"required" format:"uuid"`
	Email string                                                         `json:"email" api:"nullable"`
	Name  string                                                         `json:"name" api:"nullable"`
	JSON  annotationQueueGetAnnotationQueuesResponseAssignedReviewerJSON `json:"-"`
}

// annotationQueueGetAnnotationQueuesResponseAssignedReviewerJSON contains the JSON
// metadata for the struct
// [AnnotationQueueGetAnnotationQueuesResponseAssignedReviewer]
type annotationQueueGetAnnotationQueuesResponseAssignedReviewerJSON struct {
	ID          apijson.Field
	Email       apijson.Field
	Name        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *AnnotationQueueGetAnnotationQueuesResponseAssignedReviewer) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r annotationQueueGetAnnotationQueuesResponseAssignedReviewerJSON) RawJSON() string {
	return r.raw
}

type AnnotationQueueUpdateParams struct {
	DefaultDataset      param.Field[string]                                              `json:"default_dataset" format:"uuid"`
	Description         param.Field[string]                                              `json:"description"`
	EnableReservations  param.Field[bool]                                                `json:"enable_reservations"`
	Metadata            param.Field[AnnotationQueueUpdateParamsMetadataUnion]            `json:"metadata"`
	Name                param.Field[string]                                              `json:"name"`
	NumReviewersPerItem param.Field[AnnotationQueueUpdateParamsNumReviewersPerItemUnion] `json:"num_reviewers_per_item"`
	ReservationMinutes  param.Field[int64]                                               `json:"reservation_minutes"`
	ReviewerAccessMode  param.Field[AnnotationQueueUpdateParamsReviewerAccessMode]       `json:"reviewer_access_mode"`
	RubricInstructions  param.Field[string]                                              `json:"rubric_instructions"`
	RubricItems         param.Field[[]AnnotationQueueRubricItemSchemaParam]              `json:"rubric_items"`
}

func (r AnnotationQueueUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [AnnotationQueueUpdateParamsMetadataMap], [MissingParam].
type AnnotationQueueUpdateParamsMetadataUnion interface {
	implementsAnnotationQueueUpdateParamsMetadataUnion()
}

type AnnotationQueueUpdateParamsMetadataMap map[string]interface{}

func (r AnnotationQueueUpdateParamsMetadataMap) implementsAnnotationQueueUpdateParamsMetadataUnion() {
}

// Satisfied by [shared.UnionInt], [MissingParam].
type AnnotationQueueUpdateParamsNumReviewersPerItemUnion interface {
	ImplementsAnnotationQueueUpdateParamsNumReviewersPerItemUnion()
}

type AnnotationQueueUpdateParamsReviewerAccessMode string

const (
	AnnotationQueueUpdateParamsReviewerAccessModeAny      AnnotationQueueUpdateParamsReviewerAccessMode = "any"
	AnnotationQueueUpdateParamsReviewerAccessModeAssigned AnnotationQueueUpdateParamsReviewerAccessMode = "assigned"
)

func (r AnnotationQueueUpdateParamsReviewerAccessMode) IsKnown() bool {
	switch r {
	case AnnotationQueueUpdateParamsReviewerAccessModeAny, AnnotationQueueUpdateParamsReviewerAccessModeAssigned:
		return true
	}
	return false
}

type AnnotationQueueAnnotationQueuesParams struct {
	Name                param.Field[string]                                 `json:"name" api:"required"`
	ID                  param.Field[string]                                 `json:"id" format:"uuid"`
	CreatedAt           param.Field[time.Time]                              `json:"created_at" format:"date-time"`
	DefaultDataset      param.Field[string]                                 `json:"default_dataset" format:"uuid"`
	Description         param.Field[string]                                 `json:"description"`
	EnableReservations  param.Field[bool]                                   `json:"enable_reservations"`
	Metadata            param.Field[map[string]interface{}]                 `json:"metadata"`
	NumReviewersPerItem param.Field[int64]                                  `json:"num_reviewers_per_item"`
	ReservationMinutes  param.Field[int64]                                  `json:"reservation_minutes"`
	ReviewerAccessMode  param.Field[string]                                 `json:"reviewer_access_mode"`
	RubricInstructions  param.Field[string]                                 `json:"rubric_instructions"`
	RubricItems         param.Field[[]AnnotationQueueRubricItemSchemaParam] `json:"rubric_items"`
	SessionIDs          param.Field[[]string]                               `json:"session_ids" format:"uuid"`
	UpdatedAt           param.Field[time.Time]                              `json:"updated_at" format:"date-time"`
}

func (r AnnotationQueueAnnotationQueuesParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueNewRunStatusParams struct {
	OverrideAddedAt param.Field[time.Time] `json:"override_added_at" format:"date-time"`
	Status          param.Field[string]    `json:"status"`
}

func (r AnnotationQueueNewRunStatusParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueExportParams struct {
	EndTime                param.Field[time.Time] `json:"end_time" format:"date-time"`
	IncludeAnnotatorDetail param.Field[bool]      `json:"include_annotator_detail"`
	StartTime              param.Field[time.Time] `json:"start_time" format:"date-time"`
}

func (r AnnotationQueueExportParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueuePopulateParams struct {
	QueueID    param.Field[string]   `json:"queue_id" api:"required" format:"uuid"`
	SessionIDs param.Field[[]string] `json:"session_ids" api:"required" format:"uuid"`
}

func (r AnnotationQueuePopulateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueGetAnnotationQueuesParams struct {
	AssignedToMe param.Field[bool]                                              `query:"assigned_to_me"`
	DatasetID    param.Field[string]                                            `query:"dataset_id" format:"uuid"`
	IDs          param.Field[[]string]                                          `query:"ids" format:"uuid"`
	Limit        param.Field[int64]                                             `query:"limit"`
	Name         param.Field[string]                                            `query:"name"`
	NameContains param.Field[string]                                            `query:"name_contains"`
	Offset       param.Field[int64]                                             `query:"offset"`
	QueueType    param.Field[AnnotationQueueGetAnnotationQueuesParamsQueueType] `query:"queue_type"`
	SortBy       param.Field[string]                                            `query:"sort_by"`
	SortByDesc   param.Field[bool]                                              `query:"sort_by_desc"`
	TagValueID   param.Field[[]string]                                          `query:"tag_value_id" format:"uuid"`
}

// URLQuery serializes [AnnotationQueueGetAnnotationQueuesParams]'s query
// parameters as `url.Values`.
func (r AnnotationQueueGetAnnotationQueuesParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AnnotationQueueGetAnnotationQueuesParamsQueueType string

const (
	AnnotationQueueGetAnnotationQueuesParamsQueueTypeSingle   AnnotationQueueGetAnnotationQueuesParamsQueueType = "single"
	AnnotationQueueGetAnnotationQueuesParamsQueueTypePairwise AnnotationQueueGetAnnotationQueuesParamsQueueType = "pairwise"
)

func (r AnnotationQueueGetAnnotationQueuesParamsQueueType) IsKnown() bool {
	switch r {
	case AnnotationQueueGetAnnotationQueuesParamsQueueTypeSingle, AnnotationQueueGetAnnotationQueuesParamsQueueTypePairwise:
		return true
	}
	return false
}

type AnnotationQueueGetRunParams struct {
	IncludeExtra param.Field[bool] `query:"include_extra"`
}

// URLQuery serializes [AnnotationQueueGetRunParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueGetRunParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AnnotationQueueGetSizeParams struct {
	Status param.Field[AnnotationQueueGetSizeParamsStatus] `query:"status"`
}

// URLQuery serializes [AnnotationQueueGetSizeParams]'s query parameters as
// `url.Values`.
func (r AnnotationQueueGetSizeParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type AnnotationQueueGetSizeParamsStatus string

const (
	AnnotationQueueGetSizeParamsStatusNeedsMyReview     AnnotationQueueGetSizeParamsStatus = "needs_my_review"
	AnnotationQueueGetSizeParamsStatusNeedsOthersReview AnnotationQueueGetSizeParamsStatus = "needs_others_review"
	AnnotationQueueGetSizeParamsStatusCompleted         AnnotationQueueGetSizeParamsStatus = "completed"
)

func (r AnnotationQueueGetSizeParamsStatus) IsKnown() bool {
	switch r {
	case AnnotationQueueGetSizeParamsStatusNeedsMyReview, AnnotationQueueGetSizeParamsStatusNeedsOthersReview, AnnotationQueueGetSizeParamsStatusCompleted:
		return true
	}
	return false
}

type AnnotationQueueGetTotalArchivedParams struct {
	EndTime   param.Field[time.Time] `query:"end_time" format:"date-time"`
	StartTime param.Field[time.Time] `query:"start_time" format:"date-time"`
}

// URLQuery serializes [AnnotationQueueGetTotalArchivedParams]'s query parameters
// as `url.Values`.
func (r AnnotationQueueGetTotalArchivedParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
