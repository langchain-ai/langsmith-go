// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/langchain-ai/langsmith-go/packages/pagination"
)

// DatasetExperimentRunService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetExperimentRunService] method instead.
type DatasetExperimentRunService struct {
	Options []option.RequestOption
}

// NewDatasetExperimentRunService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewDatasetExperimentRunService(opts ...option.RequestOption) (r *DatasetExperimentRunService) {
	r = &DatasetExperimentRunService{}
	r.Options = opts
	return
}

// Returns a paginated page of dataset examples with runs from the requested
// experiments. Response uses the canonical `{items, next_cursor}` envelope.
func (r *DatasetExperimentRunService) New(ctx context.Context, datasetID string, body DatasetExperimentRunNewParams, opts ...option.RequestOption) (res *pagination.ItemsCursorPostPagination[DatasetExperimentRunNewResponse], err error) {
	var raw *http.Response
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithResponseInto(&raw)}, opts...)
	if datasetID == "" {
		err = errors.New("missing required dataset_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("v2/datasets/%s/experiment-runs", datasetID)
	cfg, err := requestconfig.NewRequestConfig(ctx, http.MethodPost, path, body, &res, opts...)
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

// Returns a paginated page of dataset examples with runs from the requested
// experiments. Response uses the canonical `{items, next_cursor}` envelope.
func (r *DatasetExperimentRunService) NewAutoPaging(ctx context.Context, datasetID string, body DatasetExperimentRunNewParams, opts ...option.RequestOption) *pagination.ItemsCursorPostPaginationAutoPager[DatasetExperimentRunNewResponse] {
	return pagination.NewItemsCursorPostPaginationAutoPager(r.New(ctx, datasetID, body, opts...))
}

type DatasetExperimentRunNewResponse struct {
	// `id` is the dataset example UUID.
	ID string `json:"id" format:"uuid"`
	// `attachment_urls` maps each attachment name to a pre-signed download URL.
	AttachmentURLs interface{} `json:"attachment_urls"`
	// `created_at` is when the example was created (RFC3339 date-time).
	CreatedAt time.Time `json:"created_at" format:"date-time"`
	// `dataset_id` is the parent dataset UUID.
	DatasetID string `json:"dataset_id" format:"uuid"`
	// `inputs` is the example input payload (arbitrary JSON object).
	Inputs interface{} `json:"inputs"`
	// `metadata` is arbitrary user-defined JSON metadata on the example.
	Metadata interface{} `json:"metadata"`
	// `modified_at` is when the example was last modified (RFC3339 date-time).
	ModifiedAt time.Time `json:"modified_at" format:"date-time"`
	// `name` is the example's optional name.
	Name string `json:"name"`
	// `outputs` is the example reference-output payload (arbitrary JSON object).
	Outputs interface{} `json:"outputs"`
	// `runs` is the list of experiment runs produced for this example.
	Runs []QueryRunResponse `json:"runs"`
	// `source_run_id` is the run UUID the example was created from, if any.
	SourceRunID string                              `json:"source_run_id" format:"uuid"`
	JSON        datasetExperimentRunNewResponseJSON `json:"-"`
}

// datasetExperimentRunNewResponseJSON contains the JSON metadata for the struct
// [DatasetExperimentRunNewResponse]
type datasetExperimentRunNewResponseJSON struct {
	ID             apijson.Field
	AttachmentURLs apijson.Field
	CreatedAt      apijson.Field
	DatasetID      apijson.Field
	Inputs         apijson.Field
	Metadata       apijson.Field
	ModifiedAt     apijson.Field
	Name           apijson.Field
	Outputs        apijson.Field
	Runs           apijson.Field
	SourceRunID    apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *DatasetExperimentRunNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetExperimentRunNewResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetExperimentRunNewParams struct {
	// `comparative_experiment_id` scopes pairwise-annotation feedback (optional).
	ComparativeExperimentID param.Field[string] `json:"comparative_experiment_id"`
	// `cursor` is the opaque string from a previous response's `next_cursor`. Absent
	// for the first page.
	Cursor param.Field[string] `json:"cursor"`
	// `example_ids` optionally restricts the page to these dataset example UUIDs (max
	// 1000).
	ExampleIDs param.Field[[]string] `json:"example_ids"`
	// `experiment_ids` lists the experiment (tracing session) UUIDs to query.
	// Required, non-empty.
	ExperimentIDs param.Field[[]string] `json:"experiment_ids"`
	// `filters` maps a project (session) UUID string to a list of filter expressions
	// (optional).
	Filters param.Field[map[string][]string] `json:"filters"`
	// `page_size` is the maximum number of examples to return. Defaults to 20,
	// max 100.
	PageSize param.Field[int64] `json:"page_size"`
	// `selects` lists which run properties to include. Omitted => only `id`. Tokens
	// mirror /v2/runs/query.
	Selects param.Field[[]DatasetExperimentRunNewParamsSelect] `json:"selects"`
	// `sort` controls feedback-score sorting (single project only).
	Sort param.Field[DatasetExperimentRunNewParamsSort] `json:"sort"`
}

func (r DatasetExperimentRunNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetExperimentRunNewParamsSelect string

const (
	DatasetExperimentRunNewParamsSelectID                     DatasetExperimentRunNewParamsSelect = "ID"
	DatasetExperimentRunNewParamsSelectName                   DatasetExperimentRunNewParamsSelect = "NAME"
	DatasetExperimentRunNewParamsSelectRunType                DatasetExperimentRunNewParamsSelect = "RUN_TYPE"
	DatasetExperimentRunNewParamsSelectStatus                 DatasetExperimentRunNewParamsSelect = "STATUS"
	DatasetExperimentRunNewParamsSelectStartTime              DatasetExperimentRunNewParamsSelect = "START_TIME"
	DatasetExperimentRunNewParamsSelectEndTime                DatasetExperimentRunNewParamsSelect = "END_TIME"
	DatasetExperimentRunNewParamsSelectLatencySeconds         DatasetExperimentRunNewParamsSelect = "LATENCY_SECONDS"
	DatasetExperimentRunNewParamsSelectFirstTokenTime         DatasetExperimentRunNewParamsSelect = "FIRST_TOKEN_TIME"
	DatasetExperimentRunNewParamsSelectError                  DatasetExperimentRunNewParamsSelect = "ERROR"
	DatasetExperimentRunNewParamsSelectErrorPreview           DatasetExperimentRunNewParamsSelect = "ERROR_PREVIEW"
	DatasetExperimentRunNewParamsSelectExtra                  DatasetExperimentRunNewParamsSelect = "EXTRA"
	DatasetExperimentRunNewParamsSelectMetadata               DatasetExperimentRunNewParamsSelect = "METADATA"
	DatasetExperimentRunNewParamsSelectEvents                 DatasetExperimentRunNewParamsSelect = "EVENTS"
	DatasetExperimentRunNewParamsSelectInputs                 DatasetExperimentRunNewParamsSelect = "INPUTS"
	DatasetExperimentRunNewParamsSelectInputsPreview          DatasetExperimentRunNewParamsSelect = "INPUTS_PREVIEW"
	DatasetExperimentRunNewParamsSelectOutputs                DatasetExperimentRunNewParamsSelect = "OUTPUTS"
	DatasetExperimentRunNewParamsSelectOutputsPreview         DatasetExperimentRunNewParamsSelect = "OUTPUTS_PREVIEW"
	DatasetExperimentRunNewParamsSelectManifest               DatasetExperimentRunNewParamsSelect = "MANIFEST"
	DatasetExperimentRunNewParamsSelectParentRunIDs           DatasetExperimentRunNewParamsSelect = "PARENT_RUN_IDS"
	DatasetExperimentRunNewParamsSelectProjectID              DatasetExperimentRunNewParamsSelect = "PROJECT_ID"
	DatasetExperimentRunNewParamsSelectTraceID                DatasetExperimentRunNewParamsSelect = "TRACE_ID"
	DatasetExperimentRunNewParamsSelectThreadID               DatasetExperimentRunNewParamsSelect = "THREAD_ID"
	DatasetExperimentRunNewParamsSelectDottedOrder            DatasetExperimentRunNewParamsSelect = "DOTTED_ORDER"
	DatasetExperimentRunNewParamsSelectIsRoot                 DatasetExperimentRunNewParamsSelect = "IS_ROOT"
	DatasetExperimentRunNewParamsSelectReferenceExampleID     DatasetExperimentRunNewParamsSelect = "REFERENCE_EXAMPLE_ID"
	DatasetExperimentRunNewParamsSelectReferenceDatasetID     DatasetExperimentRunNewParamsSelect = "REFERENCE_DATASET_ID"
	DatasetExperimentRunNewParamsSelectTotalTokens            DatasetExperimentRunNewParamsSelect = "TOTAL_TOKENS"
	DatasetExperimentRunNewParamsSelectPromptTokens           DatasetExperimentRunNewParamsSelect = "PROMPT_TOKENS"
	DatasetExperimentRunNewParamsSelectCompletionTokens       DatasetExperimentRunNewParamsSelect = "COMPLETION_TOKENS"
	DatasetExperimentRunNewParamsSelectTotalCost              DatasetExperimentRunNewParamsSelect = "TOTAL_COST"
	DatasetExperimentRunNewParamsSelectPromptCost             DatasetExperimentRunNewParamsSelect = "PROMPT_COST"
	DatasetExperimentRunNewParamsSelectCompletionCost         DatasetExperimentRunNewParamsSelect = "COMPLETION_COST"
	DatasetExperimentRunNewParamsSelectPromptTokenDetails     DatasetExperimentRunNewParamsSelect = "PROMPT_TOKEN_DETAILS"
	DatasetExperimentRunNewParamsSelectCompletionTokenDetails DatasetExperimentRunNewParamsSelect = "COMPLETION_TOKEN_DETAILS"
	DatasetExperimentRunNewParamsSelectPromptCostDetails      DatasetExperimentRunNewParamsSelect = "PROMPT_COST_DETAILS"
	DatasetExperimentRunNewParamsSelectCompletionCostDetails  DatasetExperimentRunNewParamsSelect = "COMPLETION_COST_DETAILS"
	DatasetExperimentRunNewParamsSelectPriceModelID           DatasetExperimentRunNewParamsSelect = "PRICE_MODEL_ID"
	DatasetExperimentRunNewParamsSelectTags                   DatasetExperimentRunNewParamsSelect = "TAGS"
	DatasetExperimentRunNewParamsSelectAppPath                DatasetExperimentRunNewParamsSelect = "APP_PATH"
	DatasetExperimentRunNewParamsSelectAttachments            DatasetExperimentRunNewParamsSelect = "ATTACHMENTS"
	DatasetExperimentRunNewParamsSelectThreadEvaluationTime   DatasetExperimentRunNewParamsSelect = "THREAD_EVALUATION_TIME"
	DatasetExperimentRunNewParamsSelectIsInDataset            DatasetExperimentRunNewParamsSelect = "IS_IN_DATASET"
	DatasetExperimentRunNewParamsSelectShareURL               DatasetExperimentRunNewParamsSelect = "SHARE_URL"
	DatasetExperimentRunNewParamsSelectFeedbackStats          DatasetExperimentRunNewParamsSelect = "FEEDBACK_STATS"
)

func (r DatasetExperimentRunNewParamsSelect) IsKnown() bool {
	switch r {
	case DatasetExperimentRunNewParamsSelectID, DatasetExperimentRunNewParamsSelectName, DatasetExperimentRunNewParamsSelectRunType, DatasetExperimentRunNewParamsSelectStatus, DatasetExperimentRunNewParamsSelectStartTime, DatasetExperimentRunNewParamsSelectEndTime, DatasetExperimentRunNewParamsSelectLatencySeconds, DatasetExperimentRunNewParamsSelectFirstTokenTime, DatasetExperimentRunNewParamsSelectError, DatasetExperimentRunNewParamsSelectErrorPreview, DatasetExperimentRunNewParamsSelectExtra, DatasetExperimentRunNewParamsSelectMetadata, DatasetExperimentRunNewParamsSelectEvents, DatasetExperimentRunNewParamsSelectInputs, DatasetExperimentRunNewParamsSelectInputsPreview, DatasetExperimentRunNewParamsSelectOutputs, DatasetExperimentRunNewParamsSelectOutputsPreview, DatasetExperimentRunNewParamsSelectManifest, DatasetExperimentRunNewParamsSelectParentRunIDs, DatasetExperimentRunNewParamsSelectProjectID, DatasetExperimentRunNewParamsSelectTraceID, DatasetExperimentRunNewParamsSelectThreadID, DatasetExperimentRunNewParamsSelectDottedOrder, DatasetExperimentRunNewParamsSelectIsRoot, DatasetExperimentRunNewParamsSelectReferenceExampleID, DatasetExperimentRunNewParamsSelectReferenceDatasetID, DatasetExperimentRunNewParamsSelectTotalTokens, DatasetExperimentRunNewParamsSelectPromptTokens, DatasetExperimentRunNewParamsSelectCompletionTokens, DatasetExperimentRunNewParamsSelectTotalCost, DatasetExperimentRunNewParamsSelectPromptCost, DatasetExperimentRunNewParamsSelectCompletionCost, DatasetExperimentRunNewParamsSelectPromptTokenDetails, DatasetExperimentRunNewParamsSelectCompletionTokenDetails, DatasetExperimentRunNewParamsSelectPromptCostDetails, DatasetExperimentRunNewParamsSelectCompletionCostDetails, DatasetExperimentRunNewParamsSelectPriceModelID, DatasetExperimentRunNewParamsSelectTags, DatasetExperimentRunNewParamsSelectAppPath, DatasetExperimentRunNewParamsSelectAttachments, DatasetExperimentRunNewParamsSelectThreadEvaluationTime, DatasetExperimentRunNewParamsSelectIsInDataset, DatasetExperimentRunNewParamsSelectShareURL, DatasetExperimentRunNewParamsSelectFeedbackStats:
		return true
	}
	return false
}

// `sort` controls feedback-score sorting (single project only).
type DatasetExperimentRunNewParamsSort struct {
	// `by` is the feedback selector, e.g. `feedback.correctness` (the `feedback.`
	// prefix is optional).
	By param.Field[string] `json:"by"`
	// `order` is `ASC` or `DESC` (defaults to `DESC`).
	Order param.Field[string] `json:"order"`
}

func (r DatasetExperimentRunNewParamsSort) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
