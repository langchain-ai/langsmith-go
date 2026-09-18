// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
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

// Create an Insights job configuration for a project.
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

// Update an Insights job configuration for a project.
func (r *SessionInsightConfigService) Update(ctx context.Context, sessionID string, configID string, body SessionInsightConfigUpdateParams, opts ...option.RequestOption) (res *SessionInsightConfigUpdateResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if configID == "" {
		err = errors.New("missing required config_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/configs/%s", sessionID, configID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
	return res, err
}

// List Insights job configurations for a project.
func (r *SessionInsightConfigService) List(ctx context.Context, sessionID string, query SessionInsightConfigListParams, opts ...option.RequestOption) (res *[]SessionInsightConfigListResponse, err error) {
	var env SessionInsightConfigListResponseEnvelope
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/configs", sessionID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, query, &env, opts...)
	if err != nil {
		return nil, err
	}
	res = &env.Configs
	return res, nil
}

// Delete an Insights job configuration for a project.
func (r *SessionInsightConfigService) Delete(ctx context.Context, sessionID string, configID string, opts ...option.RequestOption) (res *SessionInsightConfigDeleteResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if sessionID == "" {
		err = errors.New("missing required session_id parameter")
		return nil, err
	}
	if configID == "" {
		err = errors.New("missing required config_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v1/sessions/%s/insights/configs/%s", sessionID, configID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, &res, opts...)
	return res, err
}

// An Insights job configuration.
type SessionInsightConfigNewResponse struct {
	ID string `json:"id" api:"required" format:"uuid"`
	// Saved configuration for an Insights job.
	Config       SessionInsightConfigNewResponseConfig `json:"config" api:"required"`
	Description  string                                `json:"description" api:"required,nullable"`
	Name         string                                `json:"name" api:"required"`
	ScheduleCron string                                `json:"schedule_cron" api:"nullable"`
	JSON         sessionInsightConfigNewResponseJSON   `json:"-"`
}

// sessionInsightConfigNewResponseJSON contains the JSON metadata for the struct
// [SessionInsightConfigNewResponse]
type sessionInsightConfigNewResponseJSON struct {
	ID           apijson.Field
	Config       apijson.Field
	Description  apijson.Field
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

// Saved configuration for an Insights job.
type SessionInsightConfigNewResponseConfig struct {
	AttributeSchemas map[string]interface{}                     `json:"attribute_schemas" api:"required,nullable"`
	Filter           string                                     `json:"filter" api:"required,nullable"`
	Hierarchy        []int64                                    `json:"hierarchy" api:"required,nullable"`
	Model            SessionInsightConfigNewResponseConfigModel `json:"model" api:"required"`
	Name             string                                     `json:"name" api:"required,nullable"`
	Partitions       map[string]string                          `json:"partitions" api:"required,nullable"`
	Sample           float64                                    `json:"sample" api:"required,nullable"`
	SummaryPrompt    string                                     `json:"summary_prompt" api:"required,nullable"`
	ClusterModel     string                                     `json:"cluster_model" api:"nullable"`
	EndTime          string                                     `json:"end_time" api:"nullable"`
	LastNHours       int64                                      `json:"last_n_hours" api:"nullable"`
	StartTime        string                                     `json:"start_time" api:"nullable"`
	SummaryModel     string                                     `json:"summary_model" api:"nullable"`
	UserContext      map[string]string                          `json:"user_context" api:"nullable"`
	JSON             sessionInsightConfigNewResponseConfigJSON  `json:"-"`
}

// sessionInsightConfigNewResponseConfigJSON contains the JSON metadata for the
// struct [SessionInsightConfigNewResponseConfig]
type sessionInsightConfigNewResponseConfigJSON struct {
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

func (r *SessionInsightConfigNewResponseConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigNewResponseConfigJSON) RawJSON() string {
	return r.raw
}

type SessionInsightConfigNewResponseConfigModel string

const (
	SessionInsightConfigNewResponseConfigModelOpenAI    SessionInsightConfigNewResponseConfigModel = "openai"
	SessionInsightConfigNewResponseConfigModelAnthropic SessionInsightConfigNewResponseConfigModel = "anthropic"
)

func (r SessionInsightConfigNewResponseConfigModel) IsKnown() bool {
	switch r {
	case SessionInsightConfigNewResponseConfigModelOpenAI, SessionInsightConfigNewResponseConfigModelAnthropic:
		return true
	}
	return false
}

// An Insights job configuration.
type SessionInsightConfigUpdateResponse struct {
	ID string `json:"id" api:"required" format:"uuid"`
	// Saved configuration for an Insights job.
	Config       SessionInsightConfigUpdateResponseConfig `json:"config" api:"required"`
	Description  string                                   `json:"description" api:"required,nullable"`
	Name         string                                   `json:"name" api:"required"`
	ScheduleCron string                                   `json:"schedule_cron" api:"nullable"`
	JSON         sessionInsightConfigUpdateResponseJSON   `json:"-"`
}

// sessionInsightConfigUpdateResponseJSON contains the JSON metadata for the struct
// [SessionInsightConfigUpdateResponse]
type sessionInsightConfigUpdateResponseJSON struct {
	ID           apijson.Field
	Config       apijson.Field
	Description  apijson.Field
	Name         apijson.Field
	ScheduleCron apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SessionInsightConfigUpdateResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigUpdateResponseJSON) RawJSON() string {
	return r.raw
}

// Saved configuration for an Insights job.
type SessionInsightConfigUpdateResponseConfig struct {
	AttributeSchemas map[string]interface{}                        `json:"attribute_schemas" api:"required,nullable"`
	Filter           string                                        `json:"filter" api:"required,nullable"`
	Hierarchy        []int64                                       `json:"hierarchy" api:"required,nullable"`
	Model            SessionInsightConfigUpdateResponseConfigModel `json:"model" api:"required"`
	Name             string                                        `json:"name" api:"required,nullable"`
	Partitions       map[string]string                             `json:"partitions" api:"required,nullable"`
	Sample           float64                                       `json:"sample" api:"required,nullable"`
	SummaryPrompt    string                                        `json:"summary_prompt" api:"required,nullable"`
	ClusterModel     string                                        `json:"cluster_model" api:"nullable"`
	EndTime          string                                        `json:"end_time" api:"nullable"`
	LastNHours       int64                                         `json:"last_n_hours" api:"nullable"`
	StartTime        string                                        `json:"start_time" api:"nullable"`
	SummaryModel     string                                        `json:"summary_model" api:"nullable"`
	UserContext      map[string]string                             `json:"user_context" api:"nullable"`
	JSON             sessionInsightConfigUpdateResponseConfigJSON  `json:"-"`
}

// sessionInsightConfigUpdateResponseConfigJSON contains the JSON metadata for the
// struct [SessionInsightConfigUpdateResponseConfig]
type sessionInsightConfigUpdateResponseConfigJSON struct {
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

func (r *SessionInsightConfigUpdateResponseConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigUpdateResponseConfigJSON) RawJSON() string {
	return r.raw
}

type SessionInsightConfigUpdateResponseConfigModel string

const (
	SessionInsightConfigUpdateResponseConfigModelOpenAI    SessionInsightConfigUpdateResponseConfigModel = "openai"
	SessionInsightConfigUpdateResponseConfigModelAnthropic SessionInsightConfigUpdateResponseConfigModel = "anthropic"
)

func (r SessionInsightConfigUpdateResponseConfigModel) IsKnown() bool {
	switch r {
	case SessionInsightConfigUpdateResponseConfigModelOpenAI, SessionInsightConfigUpdateResponseConfigModelAnthropic:
		return true
	}
	return false
}

// An Insights job configuration.
type SessionInsightConfigListResponse struct {
	ID string `json:"id" api:"required" format:"uuid"`
	// Saved configuration for an Insights job.
	Config       SessionInsightConfigListResponseConfig `json:"config" api:"required"`
	Name         string                                 `json:"name" api:"required"`
	Prebuilt     bool                                   `json:"prebuilt" api:"required"`
	Description  string                                 `json:"description" api:"nullable"`
	ScheduleCron string                                 `json:"schedule_cron" api:"nullable"`
	JSON         sessionInsightConfigListResponseJSON   `json:"-"`
}

// sessionInsightConfigListResponseJSON contains the JSON metadata for the struct
// [SessionInsightConfigListResponse]
type sessionInsightConfigListResponseJSON struct {
	ID           apijson.Field
	Config       apijson.Field
	Name         apijson.Field
	Prebuilt     apijson.Field
	Description  apijson.Field
	ScheduleCron apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *SessionInsightConfigListResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigListResponseJSON) RawJSON() string {
	return r.raw
}

// Saved configuration for an Insights job.
type SessionInsightConfigListResponseConfig struct {
	AttributeSchemas map[string]interface{}                      `json:"attribute_schemas" api:"required,nullable"`
	Filter           string                                      `json:"filter" api:"required,nullable"`
	Hierarchy        []int64                                     `json:"hierarchy" api:"required,nullable"`
	Model            SessionInsightConfigListResponseConfigModel `json:"model" api:"required"`
	Name             string                                      `json:"name" api:"required,nullable"`
	Partitions       map[string]string                           `json:"partitions" api:"required,nullable"`
	Sample           float64                                     `json:"sample" api:"required,nullable"`
	SummaryPrompt    string                                      `json:"summary_prompt" api:"required,nullable"`
	ClusterModel     string                                      `json:"cluster_model" api:"nullable"`
	EndTime          string                                      `json:"end_time" api:"nullable"`
	LastNHours       int64                                       `json:"last_n_hours" api:"nullable"`
	StartTime        string                                      `json:"start_time" api:"nullable"`
	SummaryModel     string                                      `json:"summary_model" api:"nullable"`
	UserContext      map[string]string                           `json:"user_context" api:"nullable"`
	JSON             sessionInsightConfigListResponseConfigJSON  `json:"-"`
}

// sessionInsightConfigListResponseConfigJSON contains the JSON metadata for the
// struct [SessionInsightConfigListResponseConfig]
type sessionInsightConfigListResponseConfigJSON struct {
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

func (r *SessionInsightConfigListResponseConfig) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigListResponseConfigJSON) RawJSON() string {
	return r.raw
}

type SessionInsightConfigListResponseConfigModel string

const (
	SessionInsightConfigListResponseConfigModelOpenAI    SessionInsightConfigListResponseConfigModel = "openai"
	SessionInsightConfigListResponseConfigModelAnthropic SessionInsightConfigListResponseConfigModel = "anthropic"
)

func (r SessionInsightConfigListResponseConfigModel) IsKnown() bool {
	switch r {
	case SessionInsightConfigListResponseConfigModelOpenAI, SessionInsightConfigListResponseConfigModelAnthropic:
		return true
	}
	return false
}

// Confirmation that an Insights job configuration was deleted.
type SessionInsightConfigDeleteResponse struct {
	ID      string                                 `json:"id" api:"required" format:"uuid"`
	Message string                                 `json:"message" api:"required"`
	JSON    sessionInsightConfigDeleteResponseJSON `json:"-"`
}

// sessionInsightConfigDeleteResponseJSON contains the JSON metadata for the struct
// [SessionInsightConfigDeleteResponse]
type sessionInsightConfigDeleteResponseJSON struct {
	ID          apijson.Field
	Message     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightConfigDeleteResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigDeleteResponseJSON) RawJSON() string {
	return r.raw
}

type SessionInsightConfigNewParams struct {
	// Configuration for an Insights job.
	Config       param.Field[CreateRunClusteringJobRequestParam] `json:"config" api:"required"`
	Name         param.Field[string]                             `json:"name" api:"required"`
	Description  param.Field[string]                             `json:"description"`
	ScheduleCron param.Field[string]                             `json:"schedule_cron"`
}

func (r SessionInsightConfigNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionInsightConfigUpdateParams struct {
	// Configuration for an Insights job.
	Config       param.Field[CreateRunClusteringJobRequestParam] `json:"config"`
	Description  param.Field[string]                             `json:"description"`
	Name         param.Field[string]                             `json:"name"`
	ScheduleCron param.Field[string]                             `json:"schedule_cron"`
}

func (r SessionInsightConfigUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SessionInsightConfigListParams struct {
	IncludePrebuilts param.Field[bool] `query:"include_prebuilts"`
}

// URLQuery serializes [SessionInsightConfigListParams]'s query parameters as
// `url.Values`.
func (r SessionInsightConfigListParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}

// A collection of Insights job configurations.
type SessionInsightConfigListResponseEnvelope struct {
	Configs []SessionInsightConfigListResponse           `json:"configs" api:"required"`
	JSON    sessionInsightConfigListResponseEnvelopeJSON `json:"-"`
}

// sessionInsightConfigListResponseEnvelopeJSON contains the JSON metadata for the
// struct [SessionInsightConfigListResponseEnvelope]
type sessionInsightConfigListResponseEnvelopeJSON struct {
	Configs     apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SessionInsightConfigListResponseEnvelope) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sessionInsightConfigListResponseEnvelopeJSON) RawJSON() string {
	return r.raw
}
