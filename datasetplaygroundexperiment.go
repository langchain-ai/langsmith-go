// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// DatasetPlaygroundExperimentService contains methods and other services that help
// with interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewDatasetPlaygroundExperimentService] method instead.
type DatasetPlaygroundExperimentService struct {
	Options []option.RequestOption
}

// NewDatasetPlaygroundExperimentService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewDatasetPlaygroundExperimentService(opts ...option.RequestOption) (r *DatasetPlaygroundExperimentService) {
	r = &DatasetPlaygroundExperimentService{}
	r.Options = opts
	return
}

// Dataset Handler
func (r *DatasetPlaygroundExperimentService) Batch(ctx context.Context, body DatasetPlaygroundExperimentBatchParams, opts ...option.RequestOption) (res *DatasetPlaygroundExperimentBatchResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/playground_experiment/batch"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Stream Dataset Handler
func (r *DatasetPlaygroundExperimentService) Stream(ctx context.Context, body DatasetPlaygroundExperimentStreamParams, opts ...option.RequestOption) (res *DatasetPlaygroundExperimentStreamResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v1/datasets/playground_experiment/stream"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Configuration for a Runnable.
type RunnableConfigParam struct {
	Callbacks      param.Field[[]interface{}]          `json:"callbacks"`
	Configurable   param.Field[map[string]interface{}] `json:"configurable"`
	MaxConcurrency param.Field[int64]                  `json:"max_concurrency"`
	Metadata       param.Field[map[string]interface{}] `json:"metadata"`
	RecursionLimit param.Field[int64]                  `json:"recursion_limit"`
	RunID          param.Field[string]                 `json:"run_id" format:"uuid"`
	RunName        param.Field[string]                 `json:"run_name"`
	Tags           param.Field[[]string]               `json:"tags"`
}

func (r RunnableConfigParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type RunnerContextEnum string

const (
	RunnerContextEnumLangsmithUi         RunnerContextEnum = "langsmith_ui"
	RunnerContextEnumLangsmithAlignEvals RunnerContextEnum = "langsmith_align_evals"
)

func (r RunnerContextEnum) IsKnown() bool {
	switch r {
	case RunnerContextEnumLangsmithUi, RunnerContextEnumLangsmithAlignEvals:
		return true
	}
	return false
}

type DatasetPlaygroundExperimentBatchResponse = interface{}

type DatasetPlaygroundExperimentStreamResponse = interface{}

type DatasetPlaygroundExperimentBatchParams struct {
	DatasetID param.Field[string]      `json:"dataset_id,required" format:"uuid"`
	Manifest  param.Field[interface{}] `json:"manifest,required"`
	// Configuration for a Runnable.
	Options                         param.Field[RunnableConfigParam]    `json:"options,required"`
	ProjectName                     param.Field[string]                 `json:"project_name,required"`
	Secrets                         param.Field[map[string]string]      `json:"secrets,required"`
	BatchSize                       param.Field[int64]                  `json:"batch_size"`
	Commit                          param.Field[string]                 `json:"commit"`
	DatasetSplits                   param.Field[[]string]               `json:"dataset_splits"`
	EvaluatorRules                  param.Field[[]string]               `json:"evaluator_rules" format:"uuid"`
	Metadata                        param.Field[map[string]interface{}] `json:"metadata"`
	Owner                           param.Field[string]                 `json:"owner"`
	ParallelToolCalls               param.Field[bool]                   `json:"parallel_tool_calls"`
	Repetitions                     param.Field[int64]                  `json:"repetitions"`
	RepoHandle                      param.Field[string]                 `json:"repo_handle"`
	RepoID                          param.Field[string]                 `json:"repo_id"`
	RequestsPerSecond               param.Field[int64]                  `json:"requests_per_second"`
	RunID                           param.Field[string]                 `json:"run_id"`
	RunnerContext                   param.Field[RunnerContextEnum]      `json:"runner_context"`
	ToolChoice                      param.Field[string]                 `json:"tool_choice"`
	Tools                           param.Field[[]interface{}]          `json:"tools"`
	UseOrFallbackToWorkspaceSecrets param.Field[bool]                   `json:"use_or_fallback_to_workspace_secrets"`
}

func (r DatasetPlaygroundExperimentBatchParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type DatasetPlaygroundExperimentStreamParams struct {
	DatasetID param.Field[string]      `json:"dataset_id,required" format:"uuid"`
	Manifest  param.Field[interface{}] `json:"manifest,required"`
	// Configuration for a Runnable.
	Options                         param.Field[RunnableConfigParam]    `json:"options,required"`
	ProjectName                     param.Field[string]                 `json:"project_name,required"`
	Secrets                         param.Field[map[string]string]      `json:"secrets,required"`
	Commit                          param.Field[string]                 `json:"commit"`
	DatasetSplits                   param.Field[[]string]               `json:"dataset_splits"`
	EvaluatorRules                  param.Field[[]string]               `json:"evaluator_rules" format:"uuid"`
	Metadata                        param.Field[map[string]interface{}] `json:"metadata"`
	Owner                           param.Field[string]                 `json:"owner"`
	ParallelToolCalls               param.Field[bool]                   `json:"parallel_tool_calls"`
	Repetitions                     param.Field[int64]                  `json:"repetitions"`
	RepoHandle                      param.Field[string]                 `json:"repo_handle"`
	RepoID                          param.Field[string]                 `json:"repo_id"`
	RequestsPerSecond               param.Field[int64]                  `json:"requests_per_second"`
	RunID                           param.Field[string]                 `json:"run_id"`
	RunnerContext                   param.Field[RunnerContextEnum]      `json:"runner_context"`
	ToolChoice                      param.Field[string]                 `json:"tool_choice"`
	Tools                           param.Field[[]interface{}]          `json:"tools"`
	UseOrFallbackToWorkspaceSecrets param.Field[bool]                   `json:"use_or_fallback_to_workspace_secrets"`
}

func (r DatasetPlaygroundExperimentStreamParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
