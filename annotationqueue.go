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
	Options []option.RequestOption
	Runs    *AnnotationQueueRunService
}

// NewAnnotationQueueService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewAnnotationQueueService(opts ...option.RequestOption) (r *AnnotationQueueService) {
	r = &AnnotationQueueService{}
	r.Options = opts
	r.Runs = NewAnnotationQueueRunService(opts...)
	return
}

// Get Annotation Queue
func (r *AnnotationQueueService) Get(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueGetResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Update Annotation Queue
func (r *AnnotationQueueService) Update(ctx context.Context, queueID string, body AnnotationQueueUpdateParams, opts ...option.RequestOption) (res *AnnotationQueueUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// Delete Annotation Queue
func (r *AnnotationQueueService) Delete(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Create Annotation Queue
func (r *AnnotationQueueService) AnnotationQueues(ctx context.Context, body AnnotationQueueAnnotationQueuesParams, opts ...option.RequestOption) (res *AnnotationQueueSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/annotation-queues"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Create Identity Annotation Queue Run Status
func (r *AnnotationQueueService) NewRunStatus(ctx context.Context, annotationQueueRunID string, body AnnotationQueueNewRunStatusParams, opts ...option.RequestOption) (res *AnnotationQueueNewRunStatusResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if annotationQueueRunID == "" {
		err = errors.New("missing required annotation_queue_run_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/status/%s", annotationQueueRunID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Export Annotation Queue Archived Runs
func (r *AnnotationQueueService) Export(ctx context.Context, queueID string, body AnnotationQueueExportParams, opts ...option.RequestOption) (res *AnnotationQueueExportResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/export", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Populate annotation queue with runs from an experiment.
func (r *AnnotationQueueService) Populate(ctx context.Context, body AnnotationQueuePopulateParams, opts ...option.RequestOption) (res *AnnotationQueuePopulateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/annotation-queues/populate"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
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
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/queues", runID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Get a run from an annotation queue
func (r *AnnotationQueueService) GetRun(ctx context.Context, queueID string, index int64, query AnnotationQueueGetRunParams, opts ...option.RequestOption) (res *RunSchemaWithAnnotationQueueInfo, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/run/%v", queueID, index)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Get Size From Annotation Queue
func (r *AnnotationQueueService) GetSize(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueSizeSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/size", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Get Total Archived From Annotation Queue
func (r *AnnotationQueueService) GetTotalArchived(ctx context.Context, queueID string, query AnnotationQueueGetTotalArchivedParams, opts ...option.RequestOption) (res *AnnotationQueueSizeSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/total_archived", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Get Total Size From Annotation Queue
func (r *AnnotationQueueService) GetTotalSize(ctx context.Context, queueID string, opts ...option.RequestOption) (res *AnnotationQueueSizeSchema, err error) {
	opts = slices.Concat(r.Options, opts)
	if queueID == "" {
		err = errors.New("missing required queue_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/annotation-queues/%s/total_size", queueID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

type AnnotationQueueRubricItemSchema struct {
	FeedbackKey       string                              `json:"feedback_key,required"`
	Description       string                              `json:"description,nullable"`
	ScoreDescriptions map[string]string                   `json:"score_descriptions,nullable"`
	ValueDescriptions map[string]string                   `json:"value_descriptions,nullable"`
	JSON              annotationQueueRubricItemSchemaJSON `json:"-"`
}

// annotationQueueRubricItemSchemaJSON contains the JSON metadata for the struct
// [AnnotationQueueRubricItemSchema]
type annotationQueueRubricItemSchemaJSON struct {
	FeedbackKey       apijson.Field
	Description       apijson.Field
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
	FeedbackKey       param.Field[string]            `json:"feedback_key,required"`
	Description       param.Field[string]            `json:"description"`
	ScoreDescriptions param.Field[map[string]string] `json:"score_descriptions"`
	ValueDescriptions param.Field[map[string]string] `json:"value_descriptions"`
}

func (r AnnotationQueueRubricItemSchemaParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// AnnotationQueue schema.
type AnnotationQueueSchema struct {
	ID                  string                         `json:"id,required" format:"uuid"`
	Name                string                         `json:"name,required"`
	QueueType           AnnotationQueueSchemaQueueType `json:"queue_type,required"`
	TenantID            string                         `json:"tenant_id,required" format:"uuid"`
	CreatedAt           time.Time                      `json:"created_at" format:"date-time"`
	DefaultDataset      string                         `json:"default_dataset,nullable" format:"uuid"`
	Description         string                         `json:"description,nullable"`
	EnableReservations  bool                           `json:"enable_reservations,nullable"`
	Metadata            interface{}                    `json:"metadata,nullable"`
	NumReviewersPerItem int64                          `json:"num_reviewers_per_item,nullable"`
	ReservationMinutes  int64                          `json:"reservation_minutes,nullable"`
	RunRuleID           string                         `json:"run_rule_id,nullable" format:"uuid"`
	SourceRuleID        string                         `json:"source_rule_id,nullable" format:"uuid"`
	UpdatedAt           time.Time                      `json:"updated_at" format:"date-time"`
	JSON                annotationQueueSchemaJSON      `json:"-"`
}

// annotationQueueSchemaJSON contains the JSON metadata for the struct
// [AnnotationQueueSchema]
type annotationQueueSchemaJSON struct {
	ID                  apijson.Field
	Name                apijson.Field
	QueueType           apijson.Field
	TenantID            apijson.Field
	CreatedAt           apijson.Field
	DefaultDataset      apijson.Field
	Description         apijson.Field
	EnableReservations  apijson.Field
	Metadata            apijson.Field
	NumReviewersPerItem apijson.Field
	ReservationMinutes  apijson.Field
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

// Size of an Annotation Queue
type AnnotationQueueSizeSchema struct {
	Size int64                         `json:"size,required"`
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
	ID          string `json:"id,required" format:"uuid"`
	AppPath     string `json:"app_path,required"`
	DottedOrder string `json:"dotted_order,required"`
	Name        string `json:"name,required"`
	QueueRunID  string `json:"queue_run_id,required" format:"uuid"`
	// Enum for run types.
	RunType                RunSchemaWithAnnotationQueueInfoRunType   `json:"run_type,required"`
	SessionID              string                                    `json:"session_id,required" format:"uuid"`
	Status                 string                                    `json:"status,required"`
	TraceID                string                                    `json:"trace_id,required" format:"uuid"`
	AddedAt                time.Time                                 `json:"added_at,nullable" format:"date-time"`
	ChildRunIDs            []string                                  `json:"child_run_ids,nullable" format:"uuid"`
	CompletionCost         string                                    `json:"completion_cost,nullable"`
	CompletionCostDetails  map[string]string                         `json:"completion_cost_details,nullable"`
	CompletionTokenDetails map[string]int64                          `json:"completion_token_details,nullable"`
	CompletionTokens       int64                                     `json:"completion_tokens"`
	DirectChildRunIDs      []string                                  `json:"direct_child_run_ids,nullable" format:"uuid"`
	EffectiveAddedAt       time.Time                                 `json:"effective_added_at,nullable" format:"date-time"`
	EndTime                time.Time                                 `json:"end_time,nullable" format:"date-time"`
	Error                  string                                    `json:"error,nullable"`
	Events                 []interface{}                             `json:"events,nullable"`
	ExecutionOrder         int64                                     `json:"execution_order"`
	Extra                  interface{}                               `json:"extra,nullable"`
	FeedbackStats          map[string]interface{}                    `json:"feedback_stats,nullable"`
	FirstTokenTime         time.Time                                 `json:"first_token_time,nullable" format:"date-time"`
	InDataset              bool                                      `json:"in_dataset,nullable"`
	Inputs                 interface{}                               `json:"inputs,nullable"`
	InputsPreview          string                                    `json:"inputs_preview,nullable"`
	InputsS3URLs           interface{}                               `json:"inputs_s3_urls,nullable"`
	LastQueuedAt           time.Time                                 `json:"last_queued_at,nullable" format:"date-time"`
	LastReviewedTime       time.Time                                 `json:"last_reviewed_time,nullable" format:"date-time"`
	ManifestID             string                                    `json:"manifest_id,nullable" format:"uuid"`
	ManifestS3ID           string                                    `json:"manifest_s3_id,nullable" format:"uuid"`
	Outputs                interface{}                               `json:"outputs,nullable"`
	OutputsPreview         string                                    `json:"outputs_preview,nullable"`
	OutputsS3URLs          interface{}                               `json:"outputs_s3_urls,nullable"`
	ParentRunID            string                                    `json:"parent_run_id,nullable" format:"uuid"`
	ParentRunIDs           []string                                  `json:"parent_run_ids,nullable" format:"uuid"`
	PriceModelID           string                                    `json:"price_model_id,nullable" format:"uuid"`
	PromptCost             string                                    `json:"prompt_cost,nullable"`
	PromptCostDetails      map[string]string                         `json:"prompt_cost_details,nullable"`
	PromptTokenDetails     map[string]int64                          `json:"prompt_token_details,nullable"`
	PromptTokens           int64                                     `json:"prompt_tokens"`
	ReferenceDatasetID     string                                    `json:"reference_dataset_id,nullable" format:"uuid"`
	ReferenceExampleID     string                                    `json:"reference_example_id,nullable" format:"uuid"`
	S3URLs                 interface{}                               `json:"s3_urls,nullable"`
	Serialized             interface{}                               `json:"serialized,nullable"`
	ShareToken             string                                    `json:"share_token,nullable" format:"uuid"`
	StartTime              time.Time                                 `json:"start_time" format:"date-time"`
	Tags                   []string                                  `json:"tags,nullable"`
	ThreadID               string                                    `json:"thread_id,nullable"`
	TotalCost              string                                    `json:"total_cost,nullable"`
	TotalTokens            int64                                     `json:"total_tokens"`
	TraceFirstReceivedAt   time.Time                                 `json:"trace_first_received_at,nullable" format:"date-time"`
	TraceMaxStartTime      time.Time                                 `json:"trace_max_start_time,nullable" format:"date-time"`
	TraceMinStartTime      time.Time                                 `json:"trace_min_start_time,nullable" format:"date-time"`
	TraceTier              RunSchemaWithAnnotationQueueInfoTraceTier `json:"trace_tier,nullable"`
	TraceUpgrade           bool                                      `json:"trace_upgrade"`
	TtlSeconds             int64                                     `json:"ttl_seconds,nullable"`
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

// Enum for run types.
type RunSchemaWithAnnotationQueueInfoRunType string

const (
	RunSchemaWithAnnotationQueueInfoRunTypeTool      RunSchemaWithAnnotationQueueInfoRunType = "tool"
	RunSchemaWithAnnotationQueueInfoRunTypeChain     RunSchemaWithAnnotationQueueInfoRunType = "chain"
	RunSchemaWithAnnotationQueueInfoRunTypeLlm       RunSchemaWithAnnotationQueueInfoRunType = "llm"
	RunSchemaWithAnnotationQueueInfoRunTypeRetriever RunSchemaWithAnnotationQueueInfoRunType = "retriever"
	RunSchemaWithAnnotationQueueInfoRunTypeEmbedding RunSchemaWithAnnotationQueueInfoRunType = "embedding"
	RunSchemaWithAnnotationQueueInfoRunTypePrompt    RunSchemaWithAnnotationQueueInfoRunType = "prompt"
	RunSchemaWithAnnotationQueueInfoRunTypeParser    RunSchemaWithAnnotationQueueInfoRunType = "parser"
)

func (r RunSchemaWithAnnotationQueueInfoRunType) IsKnown() bool {
	switch r {
	case RunSchemaWithAnnotationQueueInfoRunTypeTool, RunSchemaWithAnnotationQueueInfoRunTypeChain, RunSchemaWithAnnotationQueueInfoRunTypeLlm, RunSchemaWithAnnotationQueueInfoRunTypeRetriever, RunSchemaWithAnnotationQueueInfoRunTypeEmbedding, RunSchemaWithAnnotationQueueInfoRunTypePrompt, RunSchemaWithAnnotationQueueInfoRunTypeParser:
		return true
	}
	return false
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
	ID                  string                              `json:"id,required" format:"uuid"`
	Name                string                              `json:"name,required"`
	QueueType           AnnotationQueueGetResponseQueueType `json:"queue_type,required"`
	TenantID            string                              `json:"tenant_id,required" format:"uuid"`
	CreatedAt           time.Time                           `json:"created_at" format:"date-time"`
	DefaultDataset      string                              `json:"default_dataset,nullable" format:"uuid"`
	Description         string                              `json:"description,nullable"`
	EnableReservations  bool                                `json:"enable_reservations,nullable"`
	Metadata            interface{}                         `json:"metadata,nullable"`
	NumReviewersPerItem int64                               `json:"num_reviewers_per_item,nullable"`
	ReservationMinutes  int64                               `json:"reservation_minutes,nullable"`
	RubricInstructions  string                              `json:"rubric_instructions,nullable"`
	RubricItems         []AnnotationQueueRubricItemSchema   `json:"rubric_items,nullable"`
	RunRuleID           string                              `json:"run_rule_id,nullable" format:"uuid"`
	SourceRuleID        string                              `json:"source_rule_id,nullable" format:"uuid"`
	UpdatedAt           time.Time                           `json:"updated_at" format:"date-time"`
	JSON                annotationQueueGetResponseJSON      `json:"-"`
}

// annotationQueueGetResponseJSON contains the JSON metadata for the struct
// [AnnotationQueueGetResponse]
type annotationQueueGetResponseJSON struct {
	ID                  apijson.Field
	Name                apijson.Field
	QueueType           apijson.Field
	TenantID            apijson.Field
	CreatedAt           apijson.Field
	DefaultDataset      apijson.Field
	Description         apijson.Field
	EnableReservations  apijson.Field
	Metadata            apijson.Field
	NumReviewersPerItem apijson.Field
	ReservationMinutes  apijson.Field
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

type AnnotationQueueUpdateResponse = interface{}

type AnnotationQueueDeleteResponse = interface{}

type AnnotationQueueNewRunStatusResponse = interface{}

type AnnotationQueueExportResponse = interface{}

type AnnotationQueuePopulateResponse = interface{}

// AnnotationQueue schema with size.
type AnnotationQueueGetAnnotationQueuesResponse struct {
	ID                  string                                              `json:"id,required" format:"uuid"`
	Name                string                                              `json:"name,required"`
	QueueType           AnnotationQueueGetAnnotationQueuesResponseQueueType `json:"queue_type,required"`
	TenantID            string                                              `json:"tenant_id,required" format:"uuid"`
	TotalRuns           int64                                               `json:"total_runs,required"`
	CreatedAt           time.Time                                           `json:"created_at" format:"date-time"`
	DefaultDataset      string                                              `json:"default_dataset,nullable" format:"uuid"`
	Description         string                                              `json:"description,nullable"`
	EnableReservations  bool                                                `json:"enable_reservations,nullable"`
	Metadata            interface{}                                         `json:"metadata,nullable"`
	NumReviewersPerItem int64                                               `json:"num_reviewers_per_item,nullable"`
	ReservationMinutes  int64                                               `json:"reservation_minutes,nullable"`
	RunRuleID           string                                              `json:"run_rule_id,nullable" format:"uuid"`
	SourceRuleID        string                                              `json:"source_rule_id,nullable" format:"uuid"`
	UpdatedAt           time.Time                                           `json:"updated_at" format:"date-time"`
	JSON                annotationQueueGetAnnotationQueuesResponseJSON      `json:"-"`
}

// annotationQueueGetAnnotationQueuesResponseJSON contains the JSON metadata for
// the struct [AnnotationQueueGetAnnotationQueuesResponse]
type annotationQueueGetAnnotationQueuesResponseJSON struct {
	ID                  apijson.Field
	Name                apijson.Field
	QueueType           apijson.Field
	TenantID            apijson.Field
	TotalRuns           apijson.Field
	CreatedAt           apijson.Field
	DefaultDataset      apijson.Field
	Description         apijson.Field
	EnableReservations  apijson.Field
	Metadata            apijson.Field
	NumReviewersPerItem apijson.Field
	ReservationMinutes  apijson.Field
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

type AnnotationQueueUpdateParams struct {
	DefaultDataset      param.Field[string]                                              `json:"default_dataset" format:"uuid"`
	Description         param.Field[string]                                              `json:"description"`
	EnableReservations  param.Field[bool]                                                `json:"enable_reservations"`
	Metadata            param.Field[MissingParam]                                        `json:"metadata"`
	Name                param.Field[string]                                              `json:"name"`
	NumReviewersPerItem param.Field[AnnotationQueueUpdateParamsNumReviewersPerItemUnion] `json:"num_reviewers_per_item"`
	ReservationMinutes  param.Field[int64]                                               `json:"reservation_minutes"`
	RubricInstructions  param.Field[string]                                              `json:"rubric_instructions"`
	RubricItems         param.Field[[]AnnotationQueueRubricItemSchemaParam]              `json:"rubric_items"`
}

func (r AnnotationQueueUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Satisfied by [shared.UnionInt], [MissingParam].
type AnnotationQueueUpdateParamsNumReviewersPerItemUnion interface {
	ImplementsAnnotationQueueUpdateParamsNumReviewersPerItemUnion()
}

type AnnotationQueueAnnotationQueuesParams struct {
	Name                param.Field[string]                                 `json:"name,required"`
	ID                  param.Field[string]                                 `json:"id" format:"uuid"`
	CreatedAt           param.Field[time.Time]                              `json:"created_at" format:"date-time"`
	DefaultDataset      param.Field[string]                                 `json:"default_dataset" format:"uuid"`
	Description         param.Field[string]                                 `json:"description"`
	EnableReservations  param.Field[bool]                                   `json:"enable_reservations"`
	Metadata            param.Field[interface{}]                            `json:"metadata"`
	NumReviewersPerItem param.Field[int64]                                  `json:"num_reviewers_per_item"`
	ReservationMinutes  param.Field[int64]                                  `json:"reservation_minutes"`
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
	EndTime   param.Field[time.Time] `json:"end_time" format:"date-time"`
	StartTime param.Field[time.Time] `json:"start_time" format:"date-time"`
}

func (r AnnotationQueueExportParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueuePopulateParams struct {
	QueueID    param.Field[string]   `json:"queue_id,required" format:"uuid"`
	SessionIDs param.Field[[]string] `json:"session_ids,required" format:"uuid"`
}

func (r AnnotationQueuePopulateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type AnnotationQueueGetAnnotationQueuesParams struct {
	DatasetID    param.Field[string]                                            `query:"dataset_id" format:"uuid"`
	IDs          param.Field[[]string]                                          `query:"ids" format:"uuid"`
	Limit        param.Field[int64]                                             `query:"limit"`
	Name         param.Field[string]                                            `query:"name"`
	NameContains param.Field[string]                                            `query:"name_contains"`
	Offset       param.Field[int64]                                             `query:"offset"`
	QueueType    param.Field[AnnotationQueueGetAnnotationQueuesParamsQueueType] `query:"queue_type"`
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
