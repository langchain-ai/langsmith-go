package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/option"
)

const (
	defaultInsightsInstructions = "How are people using my agent? What are they asking about?"
	maxInsightsChatHistories    = 1000
)

// InsightsReport is returned after creating an insights job via [Client.GenerateInsights].
// It contains the job ID and project info needed to poll for results.
type InsightsReport struct {
	// ID is the unique identifier of the insights job.
	ID string
	// Name is the display name of the insights job.
	Name string
	// Status is the current status of the insights job (e.g. "pending", "success", "error").
	Status string
	// Error contains an error message when Status is "error".
	Error string
	// ProjectID is the ID of the tracing project that holds the ingested chat histories.
	ProjectID string
	// TenantID is the ID of the tenant/workspace.
	TenantID string
}

// InsightsReportResult contains the full results of a completed insights job,
// including clusters, summary report, and optionally individual runs.
type InsightsReportResult struct {
	// ID is the unique identifier of the insights job.
	ID string
	// Name is the display name of the insights job.
	Name string
	// Status is the terminal status of the job ("success" or "error").
	Status string
	// StartTime is when the job started processing.
	StartTime *time.Time
	// EndTime is when the job finished processing.
	EndTime *time.Time
	// ConfigID is the ID of the insights configuration used.
	ConfigID string
	// Metadata contains arbitrary job metadata.
	Metadata map[string]interface{}
	// Shape maps each cluster name to its run count.
	Shape map[string]int64
	// Error contains an error message when Status is "error".
	Error string
	// Clusters is the list of run clusters produced by the insights agent.
	Clusters []SessionInsightGetJobResponseCluster
	// Report is the high-level summary produced by the insights agent.
	Report *SessionInsightGetJobResponseReport
	// Runs contains all runs associated with this job.
	// Populated only when GetInsightsReportParams.IncludeRuns is true.
	Runs []map[string]interface{}
}

// GenerateInsightsParams contains parameters for [Client.GenerateInsights].
type GenerateInsightsParams struct {
	// ChatHistories is a list of conversation histories to analyse.
	// Each item is a slice of messages; each message is a map with
	// at least "role" and "content" keys (OpenAI message format).
	// Truncated to 1000 items automatically.
	ChatHistories [][]map[string]interface{}

	// Instructions for the insights agent.
	// Defaults to "How are people using my agent? What are they asking about?"
	Instructions string

	// Name for the insights report. A random name is used when empty.
	Name string

	// Model provider: "openai" or "anthropic".
	// Auto-detected from available workspace secrets when empty.
	Model string

	// OpenAIAPIKey is an optional OpenAI API key.
	// Stored as a workspace secret when not already present.
	OpenAIAPIKey string

	// AnthropicAPIKey is an optional Anthropic API key.
	// Stored as a workspace secret when not already present.
	AnthropicAPIKey string
}

// PollInsightsParams contains parameters for [Client.PollInsights].
type PollInsightsParams struct {
	// Report is the InsightsReport returned by GenerateInsights.
	// Mutually exclusive with ID + ProjectID.
	Report *InsightsReport

	// ID is the insights job ID. Required when Report is nil.
	ID string

	// ProjectID is the tracing project ID. Required when Report is nil.
	ProjectID string

	// Rate controls how often to poll. Defaults to 30 seconds.
	Rate time.Duration

	// Timeout is the maximum time to wait before returning an error.
	// Defaults to 30 minutes.
	Timeout time.Duration

	// Verbose prints polling timestamps to stdout when true.
	Verbose bool
}

// GetInsightsReportParams contains parameters for [Client.GetInsightsReport].
type GetInsightsReportParams struct {
	// Report is the InsightsReport returned by GenerateInsights or PollInsights.
	// Mutually exclusive with ID + ProjectID.
	Report *InsightsReport

	// ID is the insights job ID. Required when Report is nil.
	ID string

	// ProjectID is the tracing project ID. Required when Report is nil.
	ProjectID string

	// IncludeRuns fetches all runs associated with the job when true.
	// Defaults to true.
	IncludeRuns bool
}

// GenerateInsights generates an insights report from a list of chat histories.
//
// The method:
//  1. Ensures an OpenAI or Anthropic API key is available as a workspace secret.
//  2. Creates a tracing project and ingests the chat histories as runs.
//  3. Submits an insights clustering job against that project.
//
// NOTE: Requires a Plus or higher tier LangSmith account.
// Report generation can take up to 30 minutes; use [Client.PollInsights] to wait
// for completion and [Client.GetInsightsReport] to retrieve full results.
func (r *Client) GenerateInsights(ctx context.Context, params GenerateInsightsParams, opts ...option.RequestOption) (*InsightsReport, error) {
	model, err := r.ensureInsightsAPIKey(ctx, params.OpenAIAPIKey, params.AnthropicAPIKey, params.Model, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: insights API key: %w", err)
	}

	project, err := r.ingestInsightsRuns(ctx, params.ChatHistories, params.Name, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: ingest insights runs: %w", err)
	}

	jobParams := SessionInsightNewParams{
		CreateRunClusteringJobRequest: CreateRunClusteringJobRequestParam{
			Name:  F(params.Name),
			Model: F(CreateRunClusteringJobRequestModel(model)),
		},
	}
	instructions := params.Instructions
	if instructions == "" {
		instructions = defaultInsightsInstructions
	}
	jobParams.CreateRunClusteringJobRequest.SummaryPrompt = F(instructions)

	resp, err := r.Sessions.Insights.New(ctx, project.ID, jobParams, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: create insights job: %w", err)
	}

	return &InsightsReport{
		ID:        resp.ID,
		Name:      resp.Name,
		Status:    resp.Status,
		Error:     resp.Error,
		ProjectID: project.ID,
		TenantID:  project.TenantID,
	}, nil
}

// PollInsights polls an insights job until it reaches a terminal state
// ("success" or "error"), or the timeout is exceeded.
//
// Provide either Report (from GenerateInsights) or both ID and ProjectID.
func (r *Client) PollInsights(ctx context.Context, params PollInsightsParams, opts ...option.RequestOption) (*InsightsReport, error) {
	projectID, jobID, err := resolveInsightsIDs(params.Report, params.ProjectID, params.ID)
	if err != nil {
		return nil, fmt.Errorf("langsmith: PollInsights: %w", err)
	}

	rate := params.Rate
	if rate <= 0 {
		rate = 30 * time.Second
	}
	timeout := params.Timeout
	if timeout <= 0 {
		timeout = 30 * time.Minute
	}

	deadline := time.Now().Add(timeout)
	for {
		resp, err := r.Sessions.Insights.GetJob(ctx, projectID, jobID, opts...)
		if err != nil {
			return nil, fmt.Errorf("langsmith: poll insights: %w", err)
		}

		if params.Verbose {
			fmt.Printf("langsmith: polling insights %s at %s — status: %s\n",
				jobID, time.Now().UTC().Format(time.RFC3339), resp.Status)
		}

		if resp.Status == "success" || resp.Status == "error" {
			report := &InsightsReport{
				ID:     resp.ID,
				Name:   resp.Name,
				Status: resp.Status,
				Error:  resp.Error,
			}
			if params.Report != nil {
				report.ProjectID = params.Report.ProjectID
				report.TenantID = params.Report.TenantID
			} else {
				report.ProjectID = projectID
			}
			return report, nil
		}

		if time.Now().After(deadline) {
			return nil, fmt.Errorf("langsmith: timed out waiting for insights job %s after %v", jobID, timeout)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(rate):
		}
	}
}

// GetInsightsReport fetches the full results of a completed insights job,
// including clusters, the high-level summary report, and optionally all runs.
//
// Provide either Report (from GenerateInsights/PollInsights) or both ID and ProjectID.
// Set GetInsightsReportParams.IncludeRuns to true (the default) to also fetch runs.
func (r *Client) GetInsightsReport(ctx context.Context, params GetInsightsReportParams, opts ...option.RequestOption) (*InsightsReportResult, error) {
	projectID, jobID, err := resolveInsightsIDs(params.Report, params.ProjectID, params.ID)
	if err != nil {
		return nil, fmt.Errorf("langsmith: GetInsightsReport: %w", err)
	}

	resp, err := r.Sessions.Insights.GetJob(ctx, projectID, jobID, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: get insights job: %w", err)
	}

	result := &InsightsReportResult{
		ID:       resp.ID,
		Name:     resp.Name,
		Status:   resp.Status,
		ConfigID: resp.ConfigID,
		Metadata: resp.Metadata,
		Shape:    resp.Shape,
		Error:    resp.Error,
		Clusters: resp.Clusters,
	}
	if !resp.StartTime.IsZero() {
		t := resp.StartTime
		result.StartTime = &t
	}
	if !resp.EndTime.IsZero() {
		t := resp.EndTime
		result.EndTime = &t
	}
	if !resp.JSON.Report.IsNull() && !resp.JSON.Report.IsMissing() {
		report := resp.Report
		result.Report = &report
	}

	if params.IncludeRuns {
		runs, err := r.fetchAllInsightsRuns(ctx, projectID, jobID, "", opts...)
		if err != nil {
			return nil, fmt.Errorf("langsmith: fetch insights runs: %w", err)
		}
		result.Runs = runs
	}

	return result, nil
}

// fetchAllInsightsRuns fetches all runs for an insights job, paginating as needed.
// Pass a non-empty clusterID to filter runs to a specific cluster.
func (r *Client) fetchAllInsightsRuns(ctx context.Context, projectID, jobID, clusterID string, opts ...option.RequestOption) ([]map[string]interface{}, error) {
	const pageSize = 100
	var all []map[string]interface{}

	for offset := int64(0); ; offset += pageSize {
		q := SessionInsightGetRunsParams{
			Limit:  F(int64(pageSize)),
			Offset: F(offset),
		}
		if clusterID != "" {
			q.ClusterID = F(clusterID)
		}

		page, err := r.Sessions.Insights.GetRuns(ctx, projectID, jobID, q, opts...)
		if err != nil {
			return nil, err
		}

		all = append(all, page.Runs...)
		if int64(len(page.Runs)) < pageSize {
			break
		}
	}

	return all, nil
}

// ingestInsightsRuns creates a tracing project and ingests chat histories as runs.
func (r *Client) ingestInsightsRuns(ctx context.Context, chatHistories [][]map[string]interface{}, name string, opts ...option.RequestOption) (*TracerSessionWithoutVirtualFields, error) {
	if len(chatHistories) > maxInsightsChatHistories {
		chatHistories = chatHistories[:maxInsightsChatHistories]
	}

	projectName := name
	if projectName == "" {
		projectName = fmt.Sprintf("insights-%s", uuid.New().String()[:8])
	}

	project, err := r.Sessions.New(ctx, SessionNewParams{Name: F(projectName)}, opts...)
	if err != nil {
		return nil, fmt.Errorf("create insights project: %w", err)
	}

	now := time.Now().UTC()
	timePrefix := fmt.Sprintf("%s%06dZ", now.Format("20060102T150405"), now.Nanosecond()/1000)

	runs := make([]RunParam, 0, len(chatHistories))
	for _, history := range chatHistories {
		runID := uuid.New().String()
		inputs := map[string]interface{}{}
		if len(history) > 0 {
			inputs["messages"] = []map[string]interface{}{history[0]}
		}
		run := RunParam{
			ID:          F(runID),
			Name:        F("chat"),
			RunType:     F(RunRunTypeChain),
			Inputs:      F(inputs),
			Outputs:     F(map[string]interface{}{"messages": history}),
			SessionID:   F(project.ID),
			TraceID:     F(runID),
			DottedOrder: F(timePrefix + runID),
			StartTime:   F(now.Format(time.RFC3339Nano)),
			EndTime:     F(now.Format(time.RFC3339Nano)),
		}
		runs = append(runs, run)
	}

	if _, err := r.Runs.IngestBatch(ctx, RunIngestBatchParams{Post: F(runs)}, opts...); err != nil {
		return nil, fmt.Errorf("ingest runs: %w", err)
	}

	return project, nil
}

// ensureInsightsAPIKey verifies that an API key for the chosen model provider is
// available as a workspace secret, storing the provided key if not.
// Returns the resolved model name ("openai" or "anthropic").
func (r *Client) ensureInsightsAPIKey(ctx context.Context, openAIKey, anthropicKey, model string, opts ...option.RequestOption) (string, error) {
	hasOpenAI, hasAnthropic, secretErr := r.fetchInsightsSecrets(ctx, opts...)

	switch model {
	case "":
		switch {
		case hasOpenAI:
			return "openai", nil
		case hasAnthropic:
			return "anthropic", nil
		case openAIKey != "":
			if secretErr == nil {
				_ = r.storeWorkspaceSecret(ctx, "OPENAI_API_KEY", openAIKey, opts...)
			}
			return "openai", nil
		case anthropicKey != "":
			if secretErr == nil {
				_ = r.storeWorkspaceSecret(ctx, "ANTHROPIC_API_KEY", anthropicKey, opts...)
			}
			return "anthropic", nil
		default:
			return "", errors.New("no model API key available: provide OpenAIAPIKey or AnthropicAPIKey, or configure workspace secrets")
		}
	case "openai":
		if !hasOpenAI {
			if openAIKey == "" {
				return "", errors.New("model is \"openai\" but no OpenAI API key provided and none found in workspace secrets")
			}
			if secretErr == nil {
				_ = r.storeWorkspaceSecret(ctx, "OPENAI_API_KEY", openAIKey, opts...)
			}
		}
		return "openai", nil
	case "anthropic":
		if !hasAnthropic {
			if anthropicKey == "" {
				return "", errors.New("model is \"anthropic\" but no Anthropic API key provided and none found in workspace secrets")
			}
			if secretErr == nil {
				_ = r.storeWorkspaceSecret(ctx, "ANTHROPIC_API_KEY", anthropicKey, opts...)
			}
		}
		return "anthropic", nil
	default:
		return "", fmt.Errorf("unknown model %q: must be \"openai\" or \"anthropic\"", model)
	}
}

// fetchInsightsSecrets checks whether OPENAI_API_KEY and ANTHROPIC_API_KEY exist
// as workspace secrets. Errors are non-fatal — the caller decides how to proceed.
func (r *Client) fetchInsightsSecrets(ctx context.Context, opts ...option.RequestOption) (hasOpenAI, hasAnthropic bool, err error) {
	var resp []struct {
		SecretName string `json:"secret_name"`
	}
	q := insightsSecretsQueryParams{
		SecretNames: []string{"OPENAI_API_KEY", "ANTHROPIC_API_KEY"},
	}
	err = r.Get(ctx, "api/v1/workspaces/current/secrets", q, &resp, opts...)
	if err != nil {
		return false, false, err
	}
	for _, s := range resp {
		switch s.SecretName {
		case "OPENAI_API_KEY":
			hasOpenAI = true
		case "ANTHROPIC_API_KEY":
			hasAnthropic = true
		}
	}
	return hasOpenAI, hasAnthropic, nil
}

// storeWorkspaceSecret upserts a workspace secret.
func (r *Client) storeWorkspaceSecret(ctx context.Context, name, value string, opts ...option.RequestOption) error {
	body := map[string]string{
		"secret_name":  name,
		"secret_value": value,
	}
	return r.Post(ctx, "api/v1/workspaces/current/secrets", body, nil, opts...)
}

// resolveInsightsIDs extracts projectID and jobID from either a report or explicit IDs.
func resolveInsightsIDs(report *InsightsReport, projectID, jobID string) (string, string, error) {
	if report != nil {
		return report.ProjectID, report.ID, nil
	}
	if projectID == "" || jobID == "" {
		return "", "", errors.New("provide either Report or both ProjectID and ID")
	}
	return projectID, jobID, nil
}

// insightsSecretsQueryParams serialises secret_names as repeated query params.
type insightsSecretsQueryParams struct {
	SecretNames []string `query:"secret_names"`
}

func (r insightsSecretsQueryParams) URLQuery() (v url.Values) {
	return apiquery.MarshalWithSettings(r, apiquery.QuerySettings{
		ArrayFormat:  apiquery.ArrayQueryFormatRepeat,
		NestedFormat: apiquery.NestedQueryFormatBrackets,
	})
}
