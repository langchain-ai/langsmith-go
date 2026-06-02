// Package runsv2 provides a hand-written client for the LangSmith v2 runs
// query endpoint (POST /v2/runs/query), which is backed by SmithDB.
//
// This sits alongside the generated SDK because the v2 endpoint is not yet
// part of the OpenAPI spec that drives Stainless code generation. The package
// follows the same convention as lib/langsmithtracing: a self-contained
// subpackage outside the regenerated surface.
package runsv2

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

// SortOrder is the start_time sort direction returned in a query result.
type SortOrder string

const (
	SortOrderAsc  SortOrder = "ASC"
	SortOrderDesc SortOrder = "DESC"
)

// SelectField identifies a run property that can be requested in QueryRequest.Selects.
// Values mirror smith-go/runs/v2/query/types.go RunSelectField.
type SelectField string

const (
	SelectID                     SelectField = "ID"
	SelectName                   SelectField = "NAME"
	SelectRunType                SelectField = "RUN_TYPE"
	SelectStatus                 SelectField = "STATUS"
	SelectStartTime              SelectField = "START_TIME"
	SelectEndTime                SelectField = "END_TIME"
	SelectLatencySeconds         SelectField = "LATENCY_SECONDS"
	SelectFirstTokenTime         SelectField = "FIRST_TOKEN_TIME"
	SelectError                  SelectField = "ERROR"
	SelectErrorPreview           SelectField = "ERROR_PREVIEW"
	SelectExtra                  SelectField = "EXTRA"
	SelectMetadata               SelectField = "METADATA"
	SelectEvents                 SelectField = "EVENTS"
	SelectInputs                 SelectField = "INPUTS"
	SelectInputsPreview          SelectField = "INPUTS_PREVIEW"
	SelectOutputs                SelectField = "OUTPUTS"
	SelectOutputsPreview         SelectField = "OUTPUTS_PREVIEW"
	SelectManifest               SelectField = "MANIFEST"
	SelectParentRunIDs           SelectField = "PARENT_RUN_IDS"
	SelectProjectID              SelectField = "PROJECT_ID"
	SelectTraceID                SelectField = "TRACE_ID"
	SelectThreadID               SelectField = "THREAD_ID"
	SelectDottedOrder            SelectField = "DOTTED_ORDER"
	SelectTags                   SelectField = "TAGS"
	SelectReferenceExampleID     SelectField = "REFERENCE_EXAMPLE_ID"
	SelectReferenceDatasetID     SelectField = "REFERENCE_DATASET_ID"
	SelectAppPath                SelectField = "APP_PATH"
	SelectIsRoot                 SelectField = "IS_ROOT"
	SelectTotalTokens            SelectField = "TOTAL_TOKENS"
	SelectPromptTokens           SelectField = "PROMPT_TOKENS"
	SelectCompletionTokens       SelectField = "COMPLETION_TOKENS"
	SelectPromptTokenDetails     SelectField = "PROMPT_TOKEN_DETAILS"
	SelectCompletionTokenDetails SelectField = "COMPLETION_TOKEN_DETAILS"
	SelectTotalCost              SelectField = "TOTAL_COST"
	SelectPromptCost             SelectField = "PROMPT_COST"
	SelectCompletionCost         SelectField = "COMPLETION_COST"
	SelectPromptCostDetails      SelectField = "PROMPT_COST_DETAILS"
	SelectCompletionCostDetails  SelectField = "COMPLETION_COST_DETAILS"
	SelectAttachments            SelectField = "ATTACHMENTS"
	SelectIsInDataset            SelectField = "IS_IN_DATASET"
	SelectShareURL               SelectField = "SHARE_URL"
	SelectFeedbackStats          SelectField = "FEEDBACK_STATS"
	SelectPriceModelID           SelectField = "PRICE_MODEL_ID"
)

// QueryRequest is the JSON body for POST /v2/runs/query.
// Field names and semantics mirror smith-go/runs/v2/query/types.go QueryRunsRequestBody.
type QueryRequest struct {
	IDs               []string      `json:"ids,omitempty"`
	TraceID           *string       `json:"trace_id,omitempty"`
	RunType           *string       `json:"run_type,omitempty"`
	ProjectIDs        []string      `json:"project_ids,omitempty"`
	ReferenceExamples []string      `json:"reference_examples,omitempty"`
	MinStartTime      *string       `json:"min_start_time,omitempty"`
	MaxStartTime      *string       `json:"max_start_time,omitempty"`
	HasError          *bool         `json:"has_error,omitempty"`
	AIQuery           *string       `json:"ai_query,omitempty"`
	Filter            *string       `json:"filter,omitempty"`
	TraceFilter       *string       `json:"trace_filter,omitempty"`
	TreeFilter        *string       `json:"tree_filter,omitempty"`
	IsRoot            *bool         `json:"is_root,omitempty"`
	Cursor            *string       `json:"cursor,omitempty"`
	PageSize          *uint64       `json:"page_size,omitempty"`
	Selects           []SelectField `json:"selects,omitempty"`
	SortOrder         *SortOrder    `json:"sort_order,omitempty"`
}

// QueryResponse is the response body for POST /v2/runs/query.
type QueryResponse struct {
	Items      []Run   `json:"items"`
	NextCursor *string `json:"next_cursor,omitempty"`
	HasMore    bool    `json:"has_more"`
}

// Run is one entry in QueryResponse.Items. Fields are pointers so unset values
// (those not requested via Selects, or null in the response) can be distinguished
// from zero values. JSON-typed fields (inputs, outputs, extra, metadata, manifest)
// are returned as raw bytes so callers can decode them on demand.
type Run struct {
	ID                 *string         `json:"id,omitempty"`
	TraceID            *string         `json:"trace_id,omitempty"`
	Name               *string         `json:"name,omitempty"`
	RunType            *string         `json:"run_type,omitempty"`
	StartTime          *string         `json:"start_time,omitempty"`
	EndTime            *string         `json:"end_time,omitempty"`
	Status             *string         `json:"status,omitempty"`
	ParentRunIDs       []string        `json:"parent_run_ids,omitempty"`
	ProjectID          *string         `json:"project_id,omitempty"`
	ReferenceExampleID *string         `json:"reference_example_id,omitempty"`
	ReferenceDatasetID *string         `json:"reference_dataset_id,omitempty"`
	DottedOrder        *string         `json:"dotted_order,omitempty"`
	Tags               []string        `json:"tags,omitempty"`
	ThreadID           *string         `json:"thread_id,omitempty"`
	AppPath            *string         `json:"app_path,omitempty"`
	IsRoot             *bool           `json:"is_root,omitempty"`
	TotalTokens        *int64          `json:"total_tokens,omitempty"`
	PromptTokens       *int64          `json:"prompt_tokens,omitempty"`
	CompletionTokens   *int64          `json:"completion_tokens,omitempty"`
	TotalCost          *float64        `json:"total_cost,omitempty"`
	PromptCost         *float64        `json:"prompt_cost,omitempty"`
	CompletionCost     *float64        `json:"completion_cost,omitempty"`
	FirstTokenTime     *string         `json:"first_token_time,omitempty"`
	LatencySeconds     *float64        `json:"latency_seconds,omitempty"`
	Inputs             json.RawMessage `json:"inputs,omitempty"`
	Outputs            json.RawMessage `json:"outputs,omitempty"`
	InputsPreview      *string         `json:"inputs_preview,omitempty"`
	OutputsPreview     *string         `json:"outputs_preview,omitempty"`
	Error              *string         `json:"error,omitempty"`
	ErrorPreview       *string         `json:"error_preview,omitempty"`
	Extra              json.RawMessage `json:"extra,omitempty"`
	Metadata           json.RawMessage `json:"metadata,omitempty"`
	Manifest           json.RawMessage `json:"manifest,omitempty"`
	Attachments        json.RawMessage `json:"attachments,omitempty"`
	IsInDataset        *bool           `json:"is_in_dataset,omitempty"`
	ShareURL           *string         `json:"share_url,omitempty"`
	FeedbackStats      json.RawMessage `json:"feedback_stats,omitempty"`
	PriceModelID       *string         `json:"price_model_id,omitempty"`
}

// HTTPError is returned by Client.Query when the server responds with a non-2xx
// status. Callers can type-assert on this to drive fallback behaviour (e.g. retry
// against the v1 endpoint when v2 returns 4xx).
type HTTPError struct {
	StatusCode int
	Body       string
}

func (e *HTTPError) Error() string {
	if e.Body == "" {
		return fmt.Sprintf("runsv2: HTTP %d", e.StatusCode)
	}
	return fmt.Sprintf("runsv2: HTTP %d: %s", e.StatusCode, e.Body)
}

// Client calls the v2 runs query endpoint.
type Client struct {
	baseURL    string
	apiKey     string
	httpClient *http.Client
}

// Option configures a Client.
type Option func(*Client)

// WithHTTPClient overrides the default http.Client.
func WithHTTPClient(c *http.Client) Option {
	return func(cl *Client) {
		if c != nil {
			cl.httpClient = c
		}
	}
}

// NewClient builds a Client for the given LangSmith endpoint and API key.
// baseURL may include a trailing "/api/v1" suffix (common for self-hosted
// LANGSMITH_ENDPOINT values); the suffix is stripped so the v2 path is
// appended cleanly. apiKey is sent as X-API-Key on every request and is
// never logged by this package.
func NewClient(baseURL, apiKey string, opts ...Option) *Client {
	c := &Client{
		baseURL:    normalizeBaseURL(baseURL),
		apiKey:     apiKey,
		httpClient: &http.Client{Timeout: 60 * time.Second},
	}
	for _, opt := range opts {
		opt(c)
	}
	return c
}

func normalizeBaseURL(raw string) string {
	u := strings.TrimRight(raw, "/")
	u = strings.TrimSuffix(u, "/api/v1")
	return strings.TrimRight(u, "/")
}

// Query posts the request to /v2/runs/query and decodes the response.
// A non-2xx status is returned as *HTTPError.
func (c *Client) Query(ctx context.Context, body QueryRequest) (*QueryResponse, error) {
	payload, err := json.Marshal(body)
	if err != nil {
		return nil, fmt.Errorf("runsv2: marshal request: %w", err)
	}
	endpoint, err := url.JoinPath(c.baseURL, "v2", "runs", "query")
	if err != nil {
		return nil, fmt.Errorf("runsv2: build url: %w", err)
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(payload))
	if err != nil {
		return nil, fmt.Errorf("runsv2: new request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	if c.apiKey != "" {
		req.Header.Set("X-API-Key", c.apiKey)
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("runsv2: send request: %w", err)
	}
	defer resp.Body.Close()

	respBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("runsv2: read response: %w", err)
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil, &HTTPError{StatusCode: resp.StatusCode, Body: string(respBytes)}
	}

	var out QueryResponse
	if err := json.Unmarshal(respBytes, &out); err != nil {
		return nil, fmt.Errorf("runsv2: decode response: %w", err)
	}
	return &out, nil
}

// Ptr is a small helper for constructing pointer fields on QueryRequest.
func Ptr[T any](v T) *T { return &v }
