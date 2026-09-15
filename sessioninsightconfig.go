// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SessionInsightConfigService contains methods and other services that help with
// interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSessionInsightConfigService] method instead.
type SessionInsightConfigService struct {
	Options []option.RequestOption
}

// NewSessionInsightConfigService generates a new service that applies the given
// options to each request. These options are applied after the parent client's
// options (if there is one), and before any request-specific options.
func NewSessionInsightConfigService(opts ...option.RequestOption) (r *SessionInsightConfigService) {
	r = &SessionInsightConfigService{}
	r.Options = opts
	return
}

// Save an insights job config.
func (r *SessionInsightConfigService) New(ctx context.Context, sessionID string, body SessionInsightConfigNewParams, opts ...option.RequestOption) (res *SessionInsightConfigNewResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/configs", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Response to create a clustering job config.
type SessionInsightConfigNewResponse struct {
	Config       SavedRunClusteringJobRequest        `json:"config" api:"required"`
	Description  string                              `json:"description" api:"required,nullable"`
	ID           string                              `json:"id" api:"required" format:"uuid"`
	Name         string                              `json:"name" api:"required"`
	ScheduleCron string                              `json:"schedule_cron" api:"nullable"`
	JSON         sessionInsightConfigNewResponseJSON `json:"-"`
}

// sessionInsightConfigNewResponseJSON contains the JSON metadata for the struct
// [SessionInsightConfigNewResponse].
type sessionInsightConfigNewResponseJSON struct {
	Config       apijson.Field
	Description  apijson.Field
	ID           apijson.Field
	Name         apijson.Field
	ScheduleCron apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SessionInsightConfigNewResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigNewResponseJSON) RawJSON() string {
	return r.raw
}

// Request to create a run clustering job.
type SavedRunClusteringJobRequest struct {
	AttributeSchemas map[string]interface{}             `json:"attribute_schemas" api:"required,nullable"`
	Filter           string                             `json:"filter" api:"required,nullable"`
	Hierarchy        []int64                            `json:"hierarchy" api:"required,nullable"`
	Model            CreateRunClusteringJobRequestModel `json:"model" api:"required"`
	Name             string                             `json:"name" api:"required,nullable"`
	Partitions       map[string]string                  `json:"partitions" api:"required,nullable"`
	Sample           float64                            `json:"sample" api:"required,nullable"`
	SummaryPrompt    string                             `json:"summary_prompt" api:"required,nullable"`
	ClusterModel     string                             `json:"cluster_model" api:"nullable"`
	EndTime          string                             `json:"end_time" api:"nullable"`
	LastNHours       int64                              `json:"last_n_hours" api:"nullable"`
	StartTime        string                             `json:"start_time" api:"nullable"`
	SummaryModel     string                             `json:"summary_model" api:"nullable"`
	UserContext      map[string]string                  `json:"user_context" api:"nullable"`
	JSON             savedRunClusteringJobRequestJSON   `json:"-"`
}

// savedRunClusteringJobRequestJSON contains the JSON metadata for the struct
// [SavedRunClusteringJobRequest].
type savedRunClusteringJobRequestJSON struct {
	AttributeSchemas apijson.Field
	Filter           apijson.Field
	Hierarchy        apijson.Field
	Model            apijson.Field
	Name             apijson.Field
	Partitions       apijson.Field
	Sample           apijson.Field
	SummaryPrompt    apijson.Field
	ClusterModel     apijson.Field
	EndTime          apijson.Field
	LastNHours       apijson.Field
	StartTime        apijson.Field
	SummaryModel     apijson.Field
	UserContext      apijson.Field
	raw              string
	ExtraFields      map[string]apijson.Field
}

func (r *SavedRunClusteringJobRequest) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r savedRunClusteringJobRequestJSON) RawJSON() string {
	return r.raw
}

type SessionInsightConfigNewParams struct {
	Config       param.Field[CreateRunClusteringJobRequestParam] `json:"config" api:"required"`
	Name         param.Field[string]                             `json:"name" api:"required"`
	Description  param.Field[string]                             `json:"description"`
	ScheduleCron param.Field[string]                             `json:"schedule_cron"`
}

func (r SessionInsightConfigNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
