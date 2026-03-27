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
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Update a session cluster job.
func (r *SessionInsightService) Update(ctx context.Context, sessionID string, jobID string, body SessionInsightUpdateParams, opts ...option.RequestOption) (res *SessionInsightUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// Delete a session cluster job.
func (r *SessionInsightService) Delete(ctx context.Context, sessionID string, jobID string, opts ...option.RequestOption) (res *SessionInsightDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// List all insights jobs for a session.
func (r *SessionInsightService) List(ctx context.Context, sessionID string, opts ...option.RequestOption) (res []SessionInsightListResponseItem, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get a specific cluster job for a session.
func (r *SessionInsightService) GetJob(ctx context.Context, sessionID string, jobID string, opts ...option.RequestOption) (res *SessionInsightGetJobResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Get all runs for a cluster job, optionally filtered by cluster.
func (r *SessionInsightService) GetRuns(ctx context.Context, sessionID string, jobID string, query SessionInsightGetRunsParams, opts ...option.RequestOption) (res *SessionInsightGetRunsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if jobID == "" {
		err = errors.New("missing required job_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/%s/runs", sessionID, jobID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &res, opts...)
	return res, err
}

// Request to create a run clustering job.
type CreateRunClusteringJobRequestParam struct {
	AttributeSchemas     param.Field[map[string]interface{}]             `json:"attribute_schemas"`
	ClusterModel         param.Field[string]                             `json:"cluster_model"`
	ConfigID             param.Field[string]                             `json:"config_id" format:"uuid"`
	EndTime              param.Field[time.Time]                          `json:"end_time" format:"date-time"`
	Filter               param.Field[string]                             `json:"filter"`
	Hierarchy            param.Field[[]int64]                            `json:"hierarchy"`
	IsScheduled          param.Field[bool]                               `json:"is_scheduled"`
	LastNHours           param.Field[int64]                              `json:"last_n_hours"`
	Model                param.Field[CreateRunClusteringJobRequestModel] `json:"model"`
	Name                 param.Field[string]                             `json:"name"`
	Partitions           param.Field[map[string]string]                  `json:"partitions"`
	Sample               param.Field[float64]                            `json:"sample"`
	StartTime            param.Field[time.Time]                          `json:"start_time" format:"date-time"`
	SummaryModel         param.Field[string]                             `json:"summary_model"`
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
	ID     string                        `json:"id" api:"required" format:"uuid"`
	Name   string                        `json:"name" api:"required"`
	Status string                        `json:"status" api:"required"`
	Error  string                        `json:"error" api:"nullable"`
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
	Name   string                           `json:"name" api:"required"`
	Status string                           `json:"status" api:"required"`
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
	ID      string                           `json:"id" api:"required" format:"uuid"`
	Message string                           `json:"message" api:"required"`
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
	ID        string                                `json:"id" api:"required" format:"uuid"`
	Clusters  []SessionInsightGetJobResponseCluster `json:"clusters" api:"required"`
	CreatedAt time.Time                             `json:"created_at" api:"required" format:"date-time"`
	Name      string                                `json:"name" api:"required"`
	Status    string                                `json:"status" api:"required"`
	ConfigID  string                                `json:"config_id" api:"nullable" format:"uuid"`
	EndTime   time.Time                             `json:"end_time" api:"nullable" format:"date-time"`
	Error     string                                `json:"error" api:"nullable"`
	Metadata  map[string]interface{}                `json:"metadata" api:"nullable"`
	// High level summary of an insights job that pulls out patterns and specific
	// traces.
	Report    SessionInsightGetJobResponseReport `json:"report" api:"nullable"`
	Shape     map[string]int64                   `json:"shape" api:"nullable"`
	StartTime time.Time                          `json:"start_time" api:"nullable" format:"date-time"`
	JSON      sessionInsightGetJobResponseJSON   `json:"-"`
}

// sessionInsightGetJobResponseJSON contains the JSON metadata for the struct
// [SessionInsightGetJobResponse]
type sessionInsightGetJobResponseJSON struct {
	ID          apijson.Field
	Clusters    apijson.Field
	CreatedAt   apijson.Field
	Name        apijson.Field
	Status      apijson.Field
	ConfigID    apijson.Field
	EndTime     apijson.Field
	Error       apijson.Field
	Metadata    apijson.Field
	Report      apijson.Field
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
	ID          string                                  `json:"id" api:"required" format:"uuid"`
	Description string                                  `json:"description" api:"required"`
	Level       int64                                   `json:"level" api:"required"`
	Name        string                                  `json:"name" api:"required"`
	NumRuns     int64                                   `json:"num_runs" api:"required"`
	Stats       map[string]interface{}                  `json:"stats" api:"required,nullable"`
	ParentID    string                                  `json:"parent_id" api:"nullable" format:"uuid"`
	ParentName  string                                  `json:"parent_name" api:"nullable"`
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

// High level summary of an insights job that pulls out patterns and specific
// traces.
type SessionInsightGetJobResponseReport struct {
	CreatedAt         time.Time                                            `json:"created_at" api:"nullable" format:"date-time"`
	HighlightedTraces []SessionInsightGetJobResponseReportHighlightedTrace `json:"highlighted_traces"`
	KeyPoints         []string                                             `json:"key_points"`
	Title             string                                               `json:"title" api:"nullable"`
	JSON              sessionInsightGetJobResponseReportJSON               `json:"-"`
}

// sessionInsightGetJobResponseReportJSON contains the JSON metadata for the struct
// [SessionInsightGetJobResponseReport]
type sessionInsightGetJobResponseReportJSON struct {
	CreatedAt         apijson.Field
	HighlightedTraces apijson.Field
	KeyPoints         apijson.Field
	Title             apijson.Field
	raw               string
	ExtraFields       map[string]apijson.Field
}

func (r *SessionInsightGetJobResponseReport) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightGetJobResponseReportJSON) RawJSON() string {
	return r.raw
}

// A trace highlighted in an insights report summary. Up to 10 per insights job.
type SessionInsightGetJobResponseReportHighlightedTrace struct {
	HighlightReason string                                                 `json:"highlight_reason" api:"required"`
	Rank            int64                                                  `json:"rank" api:"required"`
	RunID           string                                                 `json:"run_id" api:"required" format:"uuid"`
	ClusterID       string                                                 `json:"cluster_id" api:"nullable" format:"uuid"`
	ClusterName     string                                                 `json:"cluster_name" api:"nullable"`
	Summary         string                                                 `json:"summary" api:"nullable"`
	JSON            sessionInsightGetJobResponseReportHighlightedTraceJSON `json:"-"`
}

// sessionInsightGetJobResponseReportHighlightedTraceJSON contains the JSON
// metadata for the struct [SessionInsightGetJobResponseReportHighlightedTrace]
type sessionInsightGetJobResponseReportHighlightedTraceJSON struct {
	HighlightReason apijson.Field
	Rank            apijson.Field
	RunID           apijson.Field
	ClusterID       apijson.Field
	ClusterName     apijson.Field
	Summary         apijson.Field
	raw             string
	ExtraFields     map[string]apijson.Field
}

func (r *SessionInsightGetJobResponseReportHighlightedTrace) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightGetJobResponseReportHighlightedTraceJSON) RawJSON() string {
	return r.raw
}

type SessionInsightListResponseItem struct {
	ID        string                                `json:"id" api:"required" format:"uuid"`
	Name      string                                `json:"name" api:"required"`
	Status    string                                `json:"status" api:"required"`
	CreatedAt time.Time                             `json:"created_at" api:"required" format:"date-time"`
	ConfigID  string                                `json:"config_id" api:"nullable" format:"uuid"`
	EndTime   time.Time                             `json:"end_time" api:"nullable" format:"date-time"`
	Error     string                                `json:"error" api:"nullable"`
	Metadata  map[string]interface{}                `json:"metadata" api:"nullable"`
	Shape     map[string]int64                      `json:"shape" api:"nullable"`
	StartTime time.Time                             `json:"start_time" api:"nullable" format:"date-time"`
	Clusters  []SessionInsightGetJobResponseCluster `json:"clusters"`
	JSON      sessionInsightListResponseItemJSON    `json:"-"`
}

type sessionInsightListResponseItemJSON struct {
	ID          apijson.Field
	Name        apijson.Field
	Status      apijson.Field
	CreatedAt   apijson.Field
	ConfigID    apijson.Field
	EndTime     apijson.Field
	Error       apijson.Field
	Metadata    apijson.Field
	Shape       apijson.Field
	StartTime   apijson.Field
	Clusters    apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightListResponseItem) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightListResponseItemJSON) RawJSON() string {
	return r.raw
}

type SessionInsightGetRunsResponse struct {
	Offset int64                             `json:"offset" api:"required,nullable"`
	Runs   []map[string]interface{}          `json:"runs" api:"required"`
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
	CreateRunClusteringJobRequest CreateRunClusteringJobRequestParam `json:"create_run_clustering_job_request" api:"required"`
}

func (r SessionInsightNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r.CreateRunClusteringJobRequest)
}

type SessionInsightUpdateParams struct {
	Name param.Field[string] `json:"name" api:"required"`
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
