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

// SessionInsightService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionInsightService] method instead.
type SessionInsightService struct {
	Options []option.RequestOption
}

// NewSessionInsightService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSessionInsightService(opts ...option.RequestOption) (r *SessionInsightService) {
	r = &SessionInsightService{}
	r.Options = opts
	return
}

// Create an insights job.
func (r *SessionInsightService) New(ctx context.Context, sessionID string, body SessionInsightNewParams, opts ...option.RequestOption) (res *SessionInsightNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return
}

// Update a session cluster job.
func (r *SessionInsightService) Update(ctx context.Context, sessionID string, jobID string, body SessionInsightUpdateParams, opts ...option.RequestOption) (res *SessionInsightUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return
}

// Delete a session cluster job.
func (r *SessionInsightService) Delete(ctx context.Context, sessionID string, jobID string, opts ...option.RequestOption) (res *SessionInsightDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return
}

// Get a specific cluster job for a session.
func (r *SessionInsightService) GetJob(ctx context.Context, sessionID string, jobID string, opts ...option.RequestOption) (res *SessionInsightGetJobResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return
}

// Get all runs for a cluster job, optionally filtered by cluster.
func (r *SessionInsightService) GetRuns(ctx context.Context, sessionID string, jobID string, query SessionInsightGetRunsParams, opts ...option.RequestOption) (res *SessionInsightGetRunsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s/runs", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return
}

// Request to create a run clustering job.
type CreateRunClusteringJobRequestParam struct {
	AttributeSchemas     param.Field[interface{}]                        `json:"attribute_schemas"`
	EndTime              param.Field[time.Time]                          `json:"end_time" format:"date-time"`
	Filter               param.Field[string]                             `json:"filter"`
	Hierarchy            param.Field[[]int64]                            `json:"hierarchy"`
	LastNHours           param.Field[int64]                              `json:"last_n_hours"`
	Model                param.Field[CreateRunClusteringJobRequestModel] `json:"model"`
	Name                 param.Field[string]                             `json:"name"`
	Partitions           param.Field[map[string]string]                  `json:"partitions"`
	Sample               param.Field[float64]                            `json:"sample"`
	StartTime            param.Field[time.Time]                          `json:"start_time" format:"date-time"`
	SummaryPrompt        param.Field[string]                             `json:"summary_prompt"`
	UserContext          param.Field[map[string]string]                  `json:"user_context"`
	ValidateModelSecrets param.Field[bool]                               `json:"validate_model_secrets"`
}

func (r CreateRunClusteringJobRequestParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type CreateRunClusteringJobRequestModel string

const (
	CreateRunClusteringJobRequestModelOpenAI    CreateRunClusteringJobRequestModel = "openai"
	CreateRunClusteringJobRequestModelAnthropic CreateRunClusteringJobRequestModel = "anthropic"
)

func (r CreateRunClusteringJobRequestModel) IsKnown() bool {
	switch r {
	case CreateRunClusteringJobRequestModelOpenAI, CreateRunClusteringJobRequestModelAnthropic:
		return true
	}
	return false
}

// Response to creating a run clustering job.
type SessionInsightNewResponse struct {
	ID     string                        `json:"id,required" format:"uuid"`
	Name   string                        `json:"name,required"`
	Status string                        `json:"status,required"`
	Error  string                        `json:"error,nullable"`
	JSON   sessionInsightNewResponseJSON `json:"-"`
}

// sessionInsightNewResponseJSON contains the JSON metadata for the struct
// [SessionInsightNewResponse]
type sessionInsightNewResponseJSON struct {
	ID          apijson.Field
	Name        apijson.Field
	Status      apijson.Field
	Error       apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightNewResponseJSON) RawJSON() string {
	return r.raw
}

// Response to update a session cluster job.
type SessionInsightUpdateResponse struct {
	Name   string                           `json:"name,required"`
	Status string                           `json:"status,required"`
	JSON   sessionInsightUpdateResponseJSON `json:"-"`
}

// sessionInsightUpdateResponseJSON contains the JSON metadata for the struct
// [SessionInsightUpdateResponse]
type sessionInsightUpdateResponseJSON struct {
	Name        apijson.Field
	Status      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightUpdateResponseJSON) RawJSON() string {
	return r.raw
}

// Response to delete a session cluster job.
type SessionInsightDeleteResponse struct {
	ID      string                           `json:"id,required" format:"uuid"`
	Message string                           `json:"message,required"`
	JSON    sessionInsightDeleteResponseJSON `json:"-"`
}

// sessionInsightDeleteResponseJSON contains the JSON metadata for the struct
// [SessionInsightDeleteResponse]
type sessionInsightDeleteResponseJSON struct {
	ID          apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightDeleteResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightDeleteResponseJSON) RawJSON() string {
	return r.raw
}

// Response to get a specific cluster job for a session.
type SessionInsightGetJobResponse struct {
	ID        string                                `json:"id,required" format:"uuid"`
	Clusters  []SessionInsightGetJobResponseCluster `json:"clusters,required"`
	Name      string                                `json:"name,required"`
	Status    string                                `json:"status,required"`
	EndTime   time.Time                             `json:"end_time,nullable" format:"date-time"`
	Error     string                                `json:"error,nullable"`
	Metadata  interface{}                           `json:"metadata,nullable"`
	Shape     map[string]int64                      `json:"shape,nullable"`
	StartTime time.Time                             `json:"start_time,nullable" format:"date-time"`
	JSON      sessionInsightGetJobResponseJSON      `json:"-"`
}

// sessionInsightGetJobResponseJSON contains the JSON metadata for the struct
// [SessionInsightGetJobResponse]
type sessionInsightGetJobResponseJSON struct {
	ID          apijson.Field
	Clusters    apijson.Field
	Name        apijson.Field
	Status      apijson.Field
	EndTime     apijson.Field
	Error       apijson.Field
	Metadata    apijson.Field
	Shape       apijson.Field
	StartTime   apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightGetJobResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightGetJobResponseJSON) RawJSON() string {
	return r.raw
}

// A single cluster of runs.
type SessionInsightGetJobResponseCluster struct {
	ID          string                                  `json:"id,required" format:"uuid"`
	Description string                                  `json:"description,required"`
	Level       int64                                   `json:"level,required"`
	Name        string                                  `json:"name,required"`
	NumRuns     int64                                   `json:"num_runs,required"`
	Stats       interface{}                             `json:"stats,required,nullable"`
	ParentID    string                                  `json:"parent_id,nullable" format:"uuid"`
	ParentName  string                                  `json:"parent_name,nullable"`
	JSON        sessionInsightGetJobResponseClusterJSON `json:"-"`
}

// sessionInsightGetJobResponseClusterJSON contains the JSON metadata for the
// struct [SessionInsightGetJobResponseCluster]
type sessionInsightGetJobResponseClusterJSON struct {
	ID          apijson.Field
	Description apijson.Field
	Level       apijson.Field
	Name        apijson.Field
	NumRuns     apijson.Field
	Stats       apijson.Field
	ParentID    apijson.Field
	ParentName  apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightGetJobResponseCluster) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightGetJobResponseClusterJSON) RawJSON() string {
	return r.raw
}

type SessionInsightGetRunsResponse struct {
	Offset int64                             `json:"offset,required,nullable"`
	Runs   []interface{}                     `json:"runs,required"`
	JSON   sessionInsightGetRunsResponseJSON `json:"-"`
}

// sessionInsightGetRunsResponseJSON contains the JSON metadata for the struct
// [SessionInsightGetRunsResponse]
type sessionInsightGetRunsResponseJSON struct {
	Offset      apijson.Field
	Runs        apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightGetRunsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightGetRunsResponseJSON) RawJSON() string {
	return r.raw
}

type SessionInsightNewParams struct {
	// Request to create a run clustering job.
	CreateRunClusteringJobRequest CreateRunClusteringJobRequestParam `json:"create_run_clustering_job_request,required"`
}

func (r SessionInsightNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.CreateRunClusteringJobRequest)
}

type SessionInsightUpdateParams struct {
	Name param.Field[string] `json:"name,required"`
}

func (r SessionInsightUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionInsightGetRunsParams struct {
	AttributeSortKey   param.Field[string]                                        `query:"attribute_sort_key"`
	AttributeSortOrder param.Field[SessionInsightGetRunsParamsAttributeSortOrder] `query:"attribute_sort_order"`
	ClusterID          param.Field[string]                                        `query:"cluster_id" format:"uuid"`
	Limit              param.Field[int64]                                         `query:"limit"`
	Offset             param.Field[int64]                                         `query:"offset"`
}

// URLQuery serializes [SessionInsightGetRunsParams]'s query parameters as
// `url.Values`.
func (r SessionInsightGetRunsParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatComma,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

type SessionInsightGetRunsParamsAttributeSortOrder string

const (
	SessionInsightGetRunsParamsAttributeSortOrderAsc  SessionInsightGetRunsParamsAttributeSortOrder = "asc"
	SessionInsightGetRunsParamsAttributeSortOrderDesc SessionInsightGetRunsParamsAttributeSortOrder = "desc"
)

func (r SessionInsightGetRunsParamsAttributeSortOrder) IsKnown() bool {
	switch r {
	case SessionInsightGetRunsParamsAttributeSortOrderAsc, SessionInsightGetRunsParamsAttributeSortOrderDesc:
		return true
	}
	return false
}
