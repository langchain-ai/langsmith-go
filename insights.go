package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

const (
	defaultInsightsInstructions = "How are people using my agent? What are they asking about?"
	maxInsightsChatHistories    = 1000
)

// InsightsReport is returned after creating an insights job via [Client.GenerateInsights].
type InsightsReport struct {
	ID        string
	Name      string
	Status    string
	Error     string
	ProjectID string
	TenantID  string
	HostURL   string
}

// Link returns the LangSmith UI URL where this insights report can be viewed.
func (r *InsightsReport) Link() string {
	return fmt.Sprintf("%s/o/%s/projects/p/%s?tab=4&clusterJobId=%s",
		r.HostURL, r.TenantID, r.ProjectID, r.ID)
}

// InsightsReportResult contains the full results of a completed insights job.
type InsightsReportResult struct {
	ID        string
	Name      string
	Status    string
	StartTime *time.Time
	EndTime   *time.Time
	ConfigID  string
	Metadata  map[string]interface{}
	Shape     map[string]int64
	Error     string
	Clusters  []SessionInsightGetJobResponseCluster
	Report    *SessionInsightGetJobResponseReport
	Runs      []map[string]interface{}
}

type GenerateInsightsParams struct {
	ChatHistories   [][]map[string]interface{}
	Instructions    string
	Name            string
	Model           string
	OpenAIAPIKey    string
	AnthropicAPIKey string
}

type PollInsightsParams struct {
	Report    *InsightsReport
	ID        string
	ProjectID string
	Rate      time.Duration
	Timeout   time.Duration
	Verbose   bool
}

type GetInsightsReportParams struct {
	Report      *InsightsReport
	ID          string
	ProjectID   string
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

	instructions := params.Instructions
	if instructions == "" {
		instructions = defaultInsightsInstructions
	}

	jobParams := SessionInsightNewParams{
		CreateRunClusteringJobRequest: CreateRunClusteringJobRequestParam{
			Name:       F(params.Name),
			Model:      F(CreateRunClusteringJobRequestModel(model)),
			LastNHours: F(int64(1)),
			UserContext: F(map[string]string{
				"How are your agent traces structured?":          "The run.outputs.messages field contains a chat history between the user and the agent. This is all the context you need.",
				"What would you like to learn about your agent?": instructions,
			}),
		},
	}

	resp, err := r.Sessions.Insights.New(ctx, project.ID, jobParams, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: create insights job: %w", err)
	}

	report := &InsightsReport{
		ID:        resp.ID,
		Name:      resp.Name,
		Status:    resp.Status,
		Error:     resp.Error,
		ProjectID: project.ID,
		TenantID:  project.TenantID,
		HostURL:   r.hostURL(opts...),
	}

	fmt.Printf("The Insights Agent is running! This can take up to 30 minutes to complete."+
		" Once the report is completed, you'll be able to see results here: %s\n", report.Link())

	return report, nil
}

// PollInsights polls an insights job until it reaches a terminal state or the timeout is exceeded.
// Provide either Report or both ID and ProjectID, but not both.
func (r *Client) PollInsights(ctx context.Context, params PollInsightsParams, opts ...option.RequestOption) (*InsightsReport, error) {
	if (params.ID != "" || params.ProjectID != "") && params.Report != nil {
		return nil, errors.New("langsmith: PollInsights: specify exactly one of (ID and ProjectID) or Report")
	}
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

	maxTries := max(1, int(timeout/rate))
	for i := range maxTries {
		resp, err := r.Sessions.Insights.GetJob(ctx, projectID, jobID, opts...)
		if err != nil {
			return nil, fmt.Errorf("langsmith: poll insights: %w", err)
		}

		switch resp.Status {
		case "success":
			report := &InsightsReport{
				ID:      resp.ID,
				Name:    resp.Name,
				Status:  resp.Status,
				Error:   resp.Error,
				HostURL: r.hostURL(opts...),
			}
			if params.Report != nil {
				report.ProjectID = params.Report.ProjectID
				report.TenantID = params.Report.TenantID
			} else {
				report.ProjectID = projectID
			}
			fmt.Printf("Insights report completed! View the results at %s\n", report.Link())
			return report, nil
		case "error":
			return nil, fmt.Errorf("langsmith: failed to generate insights: %s", resp.Error)
		}

		if params.Verbose {
			fmt.Printf("Polling time: %v\n", time.Duration(i)*rate)
		}

		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(rate):
		}
	}

	return nil, fmt.Errorf("langsmith: timed out waiting for insights job %s after %v", jobID, timeout)
}

// GetInsightsReport fetches the full results of a completed insights job.
// Provide either Report or both ID and ProjectID, but not both.
func (r *Client) GetInsightsReport(ctx context.Context, params GetInsightsReportParams, opts ...option.RequestOption) (*InsightsReportResult, error) {
	if params.Report != nil && (params.ID != "" || params.ProjectID != "") {
		return nil, errors.New("langsmith: GetInsightsReport: specify exactly one of (ID and ProjectID) or Report")
	}
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

func (r *Client) ingestInsightsRuns(ctx context.Context, chatHistories [][]map[string]interface{}, name string, opts ...option.RequestOption) (*TracerSessionWithoutVirtualFields, error) {
	if len(chatHistories) > maxInsightsChatHistories {
		fmt.Printf("langsmith: warning: can only generate insights over %d items; truncating to first %d\n",
			maxInsightsChatHistories, maxInsightsChatHistories)
		chatHistories = chatHistories[:maxInsightsChatHistories]
	}

	projectName := name
	if projectName == "" {
		projectName = "insights " + time.Now().Format("2006-01-02 15:04:05")
	}

	project, err := r.Sessions.New(ctx, SessionNewParams{Name: F(projectName)}, opts...)
	if err != nil {
		return nil, fmt.Errorf("create insights project: %w", err)
	}

	now := time.Now().UTC()
	startTime := now.Add(-time.Second)
	timePrefix := fmt.Sprintf("%s%06dZ", now.Format("20060102T150405"), now.Nanosecond()/1000)

	runs := make([]RunParam, 0, len(chatHistories))
	for _, history := range chatHistories {
		runID := uuid.New().String()
		inputs := map[string]interface{}{}
		if len(history) > 0 {
			inputs["messages"] = history[:1]
		}
		run := RunParam{
			ID:          F(runID),
			Name:        F("trace"),
			RunType:     F(RunRunTypeChain),
			Inputs:      F(inputs),
			Outputs:     F(map[string]interface{}{"messages": history}),
			SessionID:   F(project.ID),
			TraceID:     F(runID),
			DottedOrder: F(timePrefix + runID),
			StartTime:   F(startTime.Format(time.RFC3339Nano)),
			EndTime:     F(now.Format(time.RFC3339Nano)),
		}
		runs = append(runs, run)
	}

	if _, err := r.Runs.IngestBatch(ctx, RunIngestBatchParams{Post: F(runs)}, opts...); err != nil {
		return nil, fmt.Errorf("ingest runs: %w", err)
	}

	return project, nil
}

func (r *Client) ensureInsightsAPIKey(ctx context.Context, openAIKey, anthropicKey, model string, opts ...option.RequestOption) (string, error) {
	workspaceKeys, _ := r.fetchWorkspaceSecretKeys(ctx, opts...)

	hasOpenAI := workspaceKeys["OPENAI_API_KEY"]
	hasAnthropic := workspaceKeys["ANTHROPIC_API_KEY"]

	switch model {
	case "":
		switch {
		case hasOpenAI:
			return "openai", nil
		case hasAnthropic:
			return "anthropic", nil
		case openAIKey != "":
			_ = r.storeWorkspaceSecret(ctx, "OPENAI_API_KEY", openAIKey, opts...)
			return "openai", nil
		case anthropicKey != "":
			_ = r.storeWorkspaceSecret(ctx, "ANTHROPIC_API_KEY", anthropicKey, opts...)
			return "anthropic", nil
		default:
			return "", errors.New("must specify OpenAIAPIKey or AnthropicAPIKey")
		}
	case "openai":
		if !hasOpenAI {
			if openAIKey == "" {
				return "", errors.New("model is \"openai\" but no OpenAI API key provided and none found in workspace secrets")
			}
			_ = r.storeWorkspaceSecret(ctx, "OPENAI_API_KEY", openAIKey, opts...)
		}
		return "openai", nil
	case "anthropic":
		if !hasAnthropic {
			if anthropicKey == "" {
				return "", errors.New("model is \"anthropic\" but no Anthropic API key provided and none found in workspace secrets")
			}
			_ = r.storeWorkspaceSecret(ctx, "ANTHROPIC_API_KEY", anthropicKey, opts...)
		}
		return "anthropic", nil
	default:
		return "", fmt.Errorf("unknown model %q: must be \"openai\" or \"anthropic\"", model)
	}
}

func (r *Client) fetchWorkspaceSecretKeys(ctx context.Context, opts ...option.RequestOption) (map[string]bool, error) {
	var resp []struct {
		Key string `json:"key"`
	}
	err := r.Get(ctx, "api/v1/workspaces/current/secrets", nil, &resp, opts...)
	if err != nil {
		return nil, err
	}
	keys := make(map[string]bool, len(resp))
	for _, s := range resp {
		keys[s.Key] = true
	}
	return keys, nil
}

func (r *Client) storeWorkspaceSecret(ctx context.Context, key, value string, opts ...option.RequestOption) error {
	body := []map[string]string{{"key": key, "value": value}}
	return r.Post(ctx, "api/v1/workspaces/current/secrets", body, nil, opts...)
}

func (r *Client) hostURL(opts ...option.RequestOption) string {
	cfg, err := requestconfig.NewRequestConfig(context.Background(), http.MethodGet, "", nil, nil,
		append(r.Options, opts...)...)
	if err != nil || cfg.BaseURL == nil {
		return "https://api.smith.langchain.com"
	}
	return strings.TrimRight(cfg.BaseURL.String(), "/")
}

func resolveInsightsIDs(report *InsightsReport, projectID, jobID string) (string, string, error) {
	if report != nil {
		return report.ProjectID, report.ID, nil
	}
	if projectID == "" || jobID == "" {
		return "", "", errors.New("provide either Report or both ProjectID and ID")
	}
	return projectID, jobID, nil
}
