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
func (r *DatasetExperimentRunService) Query(ctx context.Context, datasetID string, body DatasetExperimentRunQueryParams, opts ...option.RequestOption) (res *pagination.ItemsCursorPostPagination[DatasetExperimentRunQueryResponse], err error) {
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
func (r *DatasetExperimentRunService) QueryAutoPaging(ctx context.Context, datasetID string, body DatasetExperimentRunQueryParams, opts ...option.RequestOption) *pagination.ItemsCursorPostPaginationAutoPager[DatasetExperimentRunQueryResponse] {
	return pagination.NewItemsCursorPostPaginationAutoPager(r.Query(ctx, datasetID, body, opts...))
}

type DatasetExperimentRunQueryResponse struct {
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
	Runs []Run `json:"runs"`
	// `source_run_id` is the run UUID the example was created from, if any.
	SourceRunID string                                `json:"source_run_id" format:"uuid"`
	JSON        datasetExperimentRunQueryResponseJSON `json:"-"`
}

// datasetExperimentRunQueryResponseJSON contains the JSON metadata for the struct
// [DatasetExperimentRunQueryResponse]
type datasetExperimentRunQueryResponseJSON struct {
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

func (r *DatasetExperimentRunQueryResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r datasetExperimentRunQueryResponseJSON) RawJSON() string {
	return r.raw
}

type DatasetExperimentRunQueryParams struct {
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
	Selects param.Field[[]DatasetExperimentRunQueryParamsSelect] `json:"selects"`
	// `sort` controls feedback-score sorting (single project only).
	Sort param.Field[DatasetExperimentRunQueryParamsSort] `json:"sort"`
}

func (r DatasetExperimentRunQueryParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetExperimentRunQueryParamsSelect string

const (
	DatasetExperimentRunQueryParamsSelectID                     DatasetExperimentRunQueryParamsSelect = "ID"
	DatasetExperimentRunQueryParamsSelectName                   DatasetExperimentRunQueryParamsSelect = "NAME"
	DatasetExperimentRunQueryParamsSelectRunType                DatasetExperimentRunQueryParamsSelect = "RUN_TYPE"
	DatasetExperimentRunQueryParamsSelectStatus                 DatasetExperimentRunQueryParamsSelect = "STATUS"
	DatasetExperimentRunQueryParamsSelectStartTime              DatasetExperimentRunQueryParamsSelect = "START_TIME"
	DatasetExperimentRunQueryParamsSelectEndTime                DatasetExperimentRunQueryParamsSelect = "END_TIME"
	DatasetExperimentRunQueryParamsSelectLatencySeconds         DatasetExperimentRunQueryParamsSelect = "LATENCY_SECONDS"
	DatasetExperimentRunQueryParamsSelectFirstTokenTime         DatasetExperimentRunQueryParamsSelect = "FIRST_TOKEN_TIME"
	DatasetExperimentRunQueryParamsSelectError                  DatasetExperimentRunQueryParamsSelect = "ERROR"
	DatasetExperimentRunQueryParamsSelectErrorPreview           DatasetExperimentRunQueryParamsSelect = "ERROR_PREVIEW"
	DatasetExperimentRunQueryParamsSelectExtra                  DatasetExperimentRunQueryParamsSelect = "EXTRA"
	DatasetExperimentRunQueryParamsSelectMetadata               DatasetExperimentRunQueryParamsSelect = "METADATA"
	DatasetExperimentRunQueryParamsSelectEvents                 DatasetExperimentRunQueryParamsSelect = "EVENTS"
	DatasetExperimentRunQueryParamsSelectInputs                 DatasetExperimentRunQueryParamsSelect = "INPUTS"
	DatasetExperimentRunQueryParamsSelectInputsPreview          DatasetExperimentRunQueryParamsSelect = "INPUTS_PREVIEW"
	DatasetExperimentRunQueryParamsSelectOutputs                DatasetExperimentRunQueryParamsSelect = "OUTPUTS"
	DatasetExperimentRunQueryParamsSelectOutputsPreview         DatasetExperimentRunQueryParamsSelect = "OUTPUTS_PREVIEW"
	DatasetExperimentRunQueryParamsSelectManifest               DatasetExperimentRunQueryParamsSelect = "MANIFEST"
	DatasetExperimentRunQueryParamsSelectParentRunIDs           DatasetExperimentRunQueryParamsSelect = "PARENT_RUN_IDS"
	DatasetExperimentRunQueryParamsSelectProjectID              DatasetExperimentRunQueryParamsSelect = "PROJECT_ID"
	DatasetExperimentRunQueryParamsSelectTraceID                DatasetExperimentRunQueryParamsSelect = "TRACE_ID"
	DatasetExperimentRunQueryParamsSelectThreadID               DatasetExperimentRunQueryParamsSelect = "THREAD_ID"
	DatasetExperimentRunQueryParamsSelectDottedOrder            DatasetExperimentRunQueryParamsSelect = "DOTTED_ORDER"
	DatasetExperimentRunQueryParamsSelectIsRoot                 DatasetExperimentRunQueryParamsSelect = "IS_ROOT"
	DatasetExperimentRunQueryParamsSelectReferenceExampleID     DatasetExperimentRunQueryParamsSelect = "REFERENCE_EXAMPLE_ID"
	DatasetExperimentRunQueryParamsSelectReferenceDatasetID     DatasetExperimentRunQueryParamsSelect = "REFERENCE_DATASET_ID"
	DatasetExperimentRunQueryParamsSelectTotalTokens            DatasetExperimentRunQueryParamsSelect = "TOTAL_TOKENS"
	DatasetExperimentRunQueryParamsSelectPromptTokens           DatasetExperimentRunQueryParamsSelect = "PROMPT_TOKENS"
	DatasetExperimentRunQueryParamsSelectCompletionTokens       DatasetExperimentRunQueryParamsSelect = "COMPLETION_TOKENS"
	DatasetExperimentRunQueryParamsSelectTotalCost              DatasetExperimentRunQueryParamsSelect = "TOTAL_COST"
	DatasetExperimentRunQueryParamsSelectPromptCost             DatasetExperimentRunQueryParamsSelect = "PROMPT_COST"
	DatasetExperimentRunQueryParamsSelectCompletionCost         DatasetExperimentRunQueryParamsSelect = "COMPLETION_COST"
	DatasetExperimentRunQueryParamsSelectPromptTokenDetails     DatasetExperimentRunQueryParamsSelect = "PROMPT_TOKEN_DETAILS"
	DatasetExperimentRunQueryParamsSelectCompletionTokenDetails DatasetExperimentRunQueryParamsSelect = "COMPLETION_TOKEN_DETAILS"
	DatasetExperimentRunQueryParamsSelectPromptCostDetails      DatasetExperimentRunQueryParamsSelect = "PROMPT_COST_DETAILS"
	DatasetExperimentRunQueryParamsSelectCompletionCostDetails  DatasetExperimentRunQueryParamsSelect = "COMPLETION_COST_DETAILS"
	DatasetExperimentRunQueryParamsSelectPriceModelID           DatasetExperimentRunQueryParamsSelect = "PRICE_MODEL_ID"
	DatasetExperimentRunQueryParamsSelectTags                   DatasetExperimentRunQueryParamsSelect = "TAGS"
	DatasetExperimentRunQueryParamsSelectAppPath                DatasetExperimentRunQueryParamsSelect = "APP_PATH"
	DatasetExperimentRunQueryParamsSelectAttachments            DatasetExperimentRunQueryParamsSelect = "ATTACHMENTS"
	DatasetExperimentRunQueryParamsSelectThreadEvaluationTime   DatasetExperimentRunQueryParamsSelect = "THREAD_EVALUATION_TIME"
	DatasetExperimentRunQueryParamsSelectIsInDataset            DatasetExperimentRunQueryParamsSelect = "IS_IN_DATASET"
	DatasetExperimentRunQueryParamsSelectShareURL               DatasetExperimentRunQueryParamsSelect = "SHARE_URL"
	DatasetExperimentRunQueryParamsSelectFeedbackStats          DatasetExperimentRunQueryParamsSelect = "FEEDBACK_STATS"
)

func (r DatasetExperimentRunQueryParamsSelect) IsKnown() bool {
	switch r {
	case DatasetExperimentRunQueryParamsSelectID, DatasetExperimentRunQueryParamsSelectName, DatasetExperimentRunQueryParamsSelectRunType, DatasetExperimentRunQueryParamsSelectStatus, DatasetExperimentRunQueryParamsSelectStartTime, DatasetExperimentRunQueryParamsSelectEndTime, DatasetExperimentRunQueryParamsSelectLatencySeconds, DatasetExperimentRunQueryParamsSelectFirstTokenTime, DatasetExperimentRunQueryParamsSelectError, DatasetExperimentRunQueryParamsSelectErrorPreview, DatasetExperimentRunQueryParamsSelectExtra, DatasetExperimentRunQueryParamsSelectMetadata, DatasetExperimentRunQueryParamsSelectEvents, DatasetExperimentRunQueryParamsSelectInputs, DatasetExperimentRunQueryParamsSelectInputsPreview, DatasetExperimentRunQueryParamsSelectOutputs, DatasetExperimentRunQueryParamsSelectOutputsPreview, DatasetExperimentRunQueryParamsSelectManifest, DatasetExperimentRunQueryParamsSelectParentRunIDs, DatasetExperimentRunQueryParamsSelectProjectID, DatasetExperimentRunQueryParamsSelectTraceID, DatasetExperimentRunQueryParamsSelectThreadID, DatasetExperimentRunQueryParamsSelectDottedOrder, DatasetExperimentRunQueryParamsSelectIsRoot, DatasetExperimentRunQueryParamsSelectReferenceExampleID, DatasetExperimentRunQueryParamsSelectReferenceDatasetID, DatasetExperimentRunQueryParamsSelectTotalTokens, DatasetExperimentRunQueryParamsSelectPromptTokens, DatasetExperimentRunQueryParamsSelectCompletionTokens, DatasetExperimentRunQueryParamsSelectTotalCost, DatasetExperimentRunQueryParamsSelectPromptCost, DatasetExperimentRunQueryParamsSelectCompletionCost, DatasetExperimentRunQueryParamsSelectPromptTokenDetails, DatasetExperimentRunQueryParamsSelectCompletionTokenDetails, DatasetExperimentRunQueryParamsSelectPromptCostDetails, DatasetExperimentRunQueryParamsSelectCompletionCostDetails, DatasetExperimentRunQueryParamsSelectPriceModelID, DatasetExperimentRunQueryParamsSelectTags, DatasetExperimentRunQueryParamsSelectAppPath, DatasetExperimentRunQueryParamsSelectAttachments, DatasetExperimentRunQueryParamsSelectThreadEvaluationTime, DatasetExperimentRunQueryParamsSelectIsInDataset, DatasetExperimentRunQueryParamsSelectShareURL, DatasetExperimentRunQueryParamsSelectFeedbackStats:
		return true
	}
	return false
}

// `sort` controls feedback-score sorting (single project only).
type DatasetExperimentRunQueryParamsSort struct {
	// `by` is the feedback selector, e.g. `feedback.correctness` (the `feedback.`
	// prefix is optional).
	By param.Field[string] `json:"by"`
	// `order` is `ASC` or `DESC` (defaults to `DESC`).
	Order param.Field[string] `json:"order"`
}

func (r DatasetExperimentRunQueryParamsSort) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
