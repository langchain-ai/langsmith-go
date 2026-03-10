package langsmithtracing

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"math/rand/v2"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/env"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/models"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/multipart"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/tracesink"
)

// Attachment is a binary file to upload alongside a run.
type Attachment = models.Attachment

// DrainConfig controls batching and auto-scaling behavior.
type DrainConfig = tracesink.DrainConfig

func DefaultDrainConfig() DrainConfig { return tracesink.DefaultDrainConfig() }

const (
	filteredTTL           = 5 * time.Minute
	filteredPruneInterval = 1 * time.Minute
)

// TracingClient sends runs to LangSmith via the multipart ingestion endpoint.
type TracingClient struct {
	sink    *tracesink.TraceSink
	project string

	sampleRate        *float64
	filteredMu        sync.Mutex
	filteredTraces    map[uuid.UUID]time.Time
	filteredLastPrune time.Time
}

// RunCreate holds parameters for creating a new run (multipart post).
// Setting Outputs, EndTime, or Error enables single-shot complete runs
// (create + end in one call) without a separate UpdateRun.
type RunCreate struct {
	ID          uuid.UUID
	TraceID     uuid.UUID
	ParentRunID *uuid.UUID
	Name        string
	RunType     string // "chain", "llm", "tool", etc.
	Inputs      map[string]any
	Outputs     map[string]any        // Set to create a complete run in one call.
	Extra       map[string]any
	Events      []map[string]any      // e.g. [{"name":"new_token","time":"...","kwargs":{...}}]
	Serialized  map[string]any        // Model manifest; only retained for "llm"/"prompt" run types.
	Tags        []string
	StartTime   time.Time
	EndTime     time.Time             // Set to create a complete run in one call.
	DottedOrder string
	Error       string                // Set to create a failed run in one call.
	Attachments map[string]Attachment // keyed by name; names must not contain '.'

	SessionName        string         // Per-run project name override (defaults to client project).
	SessionID          *uuid.UUID     // Per-run project ID override.
	ReferenceExampleID *uuid.UUID     // Links run to a dataset example (evaluations).
	InputAttachments   map[string]any // Attachment metadata for input fields.
	OutputAttachments  map[string]any // Attachment metadata for output fields.
}

// RunUpdate holds parameters for updating an existing run (multipart patch).
// All fields except ID, TraceID, and DottedOrder are optional.
type RunUpdate struct {
	ID          uuid.UUID
	TraceID     uuid.UUID
	Outputs     map[string]any
	Inputs      map[string]any        // override/replace inputs set at create time
	Extra       map[string]any        // override/replace extra set at create time
	Events      []map[string]any      // e.g. [{"name":"new_token","time":"..."}]
	Name        string                // rename the run
	RunType     string                // change run type ("chain", "llm", "tool", etc.)
	Tags        []string              // add or replace tags
	EndTime     time.Time
	StartTime   time.Time             // adjust start time
	DottedOrder string
	Error       string
	Attachments map[string]Attachment // keyed by name; names must not contain '.'

	SessionName        string         // Per-run project name override.
	SessionID          *uuid.UUID     // Per-run project ID override.
	ReferenceExampleID *uuid.UUID     // Links run to a dataset example (evaluations).
	InputAttachments   map[string]any // Attachment metadata for input fields.
	OutputAttachments  map[string]any // Attachment metadata for output fields.
}

// RunOp is a decoded run operation exposed to transform hooks.
// See [models.RunOp] for details.
type RunOp = models.RunOp

type RunTransformFunc = tracesink.RunTransformFunc

type Option func(*options)

type options struct {
	apiURL       string
	apiKey       string
	bearerToken  string
	project      string
	drainConfig  *tracesink.DrainConfig
	sampleRate   *float64
	runTransform RunTransformFunc
}

// WithAPIURL overrides the LangSmith API URL.
func WithAPIURL(url string) Option { return func(o *options) { o.apiURL = url } }

// WithAPIKey overrides the LangSmith API key.
func WithAPIKey(key string) Option { return func(o *options) { o.apiKey = key } }

// WithBearerToken sets a bearer token for authentication.
// When set, it takes precedence over the API key.
func WithBearerToken(token string) Option { return func(o *options) { o.bearerToken = token } }

// WithProject overrides the LangSmith project name.
func WithProject(name string) Option { return func(o *options) { o.project = name } }

// WithDrainConfig overrides the default drain/scaling configuration.
func WithDrainConfig(config DrainConfig) Option {
	return func(o *options) { o.drainConfig = &config }
}

// WithSampleRate sets the trace sampling rate (must be between 0 and 1).
// Overrides the LANGSMITH_TRACING_SAMPLING_RATE env var.
// Out-of-range values are clamped with a warning log.
func WithSampleRate(rate float64) Option {
	if rate < 0 || rate > 1 {
		log.Printf("[langsmith] WithSampleRate: rate %f out of range, clamping to [0, 1]", rate)
		rate = max(0, min(1, rate))
	}
	return func(o *options) { o.sampleRate = &rate }
}

// WithRunTransform sets a pre-export transform hook.
func WithRunTransform(fn RunTransformFunc) Option {
	return func(o *options) { o.runTransform = fn }
}

// NewTracingClient creates a TracingClient that sends runs via multipart ingestion.
// The context is propagated to HTTP requests during normal operation.
// Close always drains with a background context to guarantee delivery.
func NewTracingClient(ctx context.Context, opts ...Option) *TracingClient {
	cfg := options{
		apiURL:  env.APIURL(),
		apiKey:  env.APIKey(),
		project: env.Project(),
	}
	for _, o := range opts {
		o(&cfg)
	}

	drainCfg := tracesink.DefaultDrainConfig()
	if cfg.drainConfig != nil {
		drainCfg = *cfg.drainConfig
	}

	sampleRate := cfg.sampleRate
	if sampleRate == nil {
		sampleRate = env.TracingSampleRate()
	}

	endpoint := models.WriteEndpoint{
		URL:         cfg.apiURL,
		Key:         cfg.apiKey,
		BearerToken: cfg.bearerToken,
		Project:     cfg.project,
	}

	exp := multipart.NewExporter(nil, multipart.DefaultRetry())
	sink := tracesink.NewTraceSink(ctx, exp, drainCfg, endpoint, cfg.runTransform)

	return &TracingClient{
		sink:              sink,
		project:           cfg.project,
		sampleRate:        sampleRate,
		filteredTraces:    make(map[uuid.UUID]time.Time),
		filteredLastPrune: time.Now(),
	}
}

// CreateRun enqueues a run create (post) for multipart ingestion.
// If sampling is enabled, root runs are randomly kept/dropped and the
// decision is applied to all children in the same trace.
func (c *TracingClient) CreateRun(r *RunCreate) error {
	if !c.shouldSampleCreate(r.ID, r.TraceID) {
		return nil
	}
	if err := validateAttachmentNames(r.Attachments); err != nil {
		return err
	}

	sessionName := c.project
	if r.SessionName != "" {
		sessionName = r.SessionName
	}

	runInfo := map[string]any{
		"id":           r.ID.String(),
		"trace_id":     r.TraceID.String(),
		"name":         r.Name,
		"run_type":     r.RunType,
		"session_name": sessionName,
		"start_time":   r.StartTime.UTC().Format(time.RFC3339Nano),
		"dotted_order": r.DottedOrder,
	}
	if r.ParentRunID != nil {
		runInfo["parent_run_id"] = r.ParentRunID.String()
	}
	if len(r.Tags) > 0 {
		runInfo["tags"] = r.Tags
	}
	if r.SessionID != nil {
		runInfo["session_id"] = r.SessionID.String()
	}
	if r.ReferenceExampleID != nil {
		runInfo["reference_example_id"] = r.ReferenceExampleID.String()
	}
	if r.InputAttachments != nil {
		runInfo["input_attachments"] = r.InputAttachments
	}
	if r.OutputAttachments != nil {
		runInfo["output_attachments"] = r.OutputAttachments
	}
	if !r.EndTime.IsZero() {
		runInfo["end_time"] = r.EndTime.UTC().Format(time.RFC3339Nano)
	}
	if r.Error != "" {
		runInfo["status"] = "error"
	}

	runInfoBytes, err := json.Marshal(runInfo)
	if err != nil {
		return fmt.Errorf("marshal run info: %w", err)
	}

	extra := mergeRuntimeEnv(r.Extra)

	var serialized map[string]any
	if r.Serialized != nil && (r.RunType == "llm" || r.RunType == "prompt") {
		serialized = make(map[string]any, len(r.Serialized))
		for k, v := range r.Serialized {
			if k != "graph" {
				serialized[k] = v
			}
		}
	}

	op, err := buildOp(models.OpKindPost, r.ID, r.TraceID, runInfoBytes, r.Inputs, r.Outputs, extra, r.Events, r.Error, serialized, r.Attachments)
	if err != nil {
		return err
	}
	return c.sink.Submit(op)
}

// UpdateRun enqueues a run update (patch) for multipart ingestion.
// If the run's trace was sampled out during CreateRun, the update is dropped.
func (c *TracingClient) UpdateRun(r *RunUpdate) error {
	if !c.shouldSampleUpdate(r.ID, r.TraceID) {
		return nil
	}
	if err := validateAttachmentNames(r.Attachments); err != nil {
		return err
	}

	runInfo := map[string]any{
		"id":           r.ID.String(),
		"trace_id":     r.TraceID.String(),
		"dotted_order": r.DottedOrder,
	}
	if !r.EndTime.IsZero() {
		runInfo["end_time"] = r.EndTime.UTC().Format(time.RFC3339Nano)
	}
	if r.Error != "" {
		runInfo["status"] = "error"
	}
	if r.Name != "" {
		runInfo["name"] = r.Name
	}
	if r.RunType != "" {
		runInfo["run_type"] = r.RunType
	}
	if !r.StartTime.IsZero() {
		runInfo["start_time"] = r.StartTime.UTC().Format(time.RFC3339Nano)
	}
	if len(r.Tags) > 0 {
		runInfo["tags"] = r.Tags
	}
	if r.SessionName != "" {
		runInfo["session_name"] = r.SessionName
	}
	if r.SessionID != nil {
		runInfo["session_id"] = r.SessionID.String()
	}
	if r.ReferenceExampleID != nil {
		runInfo["reference_example_id"] = r.ReferenceExampleID.String()
	}
	if r.InputAttachments != nil {
		runInfo["input_attachments"] = r.InputAttachments
	}
	if r.OutputAttachments != nil {
		runInfo["output_attachments"] = r.OutputAttachments
	}

	runInfoBytes, err := json.Marshal(runInfo)
	if err != nil {
		return fmt.Errorf("marshal run info: %w", err)
	}

	op, err := buildOp(models.OpKindPatch, r.ID, r.TraceID, runInfoBytes, r.Inputs, r.Outputs, r.Extra, r.Events, r.Error, nil, r.Attachments)
	if err != nil {
		return err
	}
	return c.sink.Submit(op)
}

func buildOp(
	kind models.OpKind, id, traceID uuid.UUID, runInfo []byte,
	inputs, outputs, extra map[string]any, events []map[string]any,
	runErr string, serialized map[string]any, attachments map[string]Attachment,
) (*models.SerializedOp, error) {
	op := &models.SerializedOp{
		Kind:        kind,
		ID:          id,
		TraceID:     traceID,
		RunInfo:     runInfo,
		Attachments: attachments,
	}
	var err error
	if inputs != nil {
		if op.Inputs, err = json.Marshal(inputs); err != nil {
			return nil, fmt.Errorf("marshal inputs: %w", err)
		}
	}
	if outputs != nil {
		if op.Outputs, err = json.Marshal(outputs); err != nil {
			return nil, fmt.Errorf("marshal outputs: %w", err)
		}
	}
	if extra != nil {
		if op.Extra, err = json.Marshal(extra); err != nil {
			return nil, fmt.Errorf("marshal extra: %w", err)
		}
	}
	if events != nil {
		if op.Events, err = json.Marshal(events); err != nil {
			return nil, fmt.Errorf("marshal events: %w", err)
		}
	}
	if runErr != "" {
		if op.Error, err = json.Marshal(runErr); err != nil {
			return nil, fmt.Errorf("marshal error: %w", err)
		}
	}
	if serialized != nil {
		if op.Serialized, err = json.Marshal(serialized); err != nil {
			return nil, fmt.Errorf("marshal serialized: %w", err)
		}
	}
	return op, nil
}

// Close flushes pending operations and shuts down the client.
func (c *TracingClient) Close() {
	c.sink.Close()
}

// shouldSampleCreate decides whether a create (post) should be kept.
// For root runs (id == traceID), a random check is made against the sample rate.
// Child runs follow the decision of their root. 
func (c *TracingClient) shouldSampleCreate(id, traceID uuid.UUID) bool {
	if c.sampleRate == nil {
		return true
	}

	c.filteredMu.Lock()
	defer c.filteredMu.Unlock()

	if _, filtered := c.filteredTraces[traceID]; filtered {
		return false
	}

	if id == traceID {
		if rand.Float64() < *c.sampleRate {
			return true
		}
		now := time.Now()
		c.filteredTraces[traceID] = now
		c.pruneFilteredLocked(now)
		return false
	}
	return true
}

// pruneFilteredLocked removes entries older than filteredTTL.
func (c *TracingClient) pruneFilteredLocked(now time.Time) {
	if now.Sub(c.filteredLastPrune) < filteredPruneInterval {
		return
	}
	c.filteredLastPrune = now
	cutoff := now.Add(-filteredTTL)
	for id, ts := range c.filteredTraces {
		if ts.Before(cutoff) {
			delete(c.filteredTraces, id)
		}
	}
}

// shouldSampleUpdate decides whether an update (patch) should be kept.
// Updates for traces that were sampled out are dropped.
// Cleanup of the filteredTraces map is handled by TTL-based pruning in
// pruneFilteredLocked, not here, to avoid prematurely removing entries
// that child-run updates still need to consult.
func (c *TracingClient) shouldSampleUpdate(id, traceID uuid.UUID) bool {
	if c.sampleRate == nil {
		return true
	}

	c.filteredMu.Lock()
	defer c.filteredMu.Unlock()

	_, filtered := c.filteredTraces[traceID]
	return !filtered
}

func mergeRuntimeEnv(extra map[string]any) map[string]any {
	result := make(map[string]any, len(extra)+2)
	for k, v := range extra {
		result[k] = v
	}

	baseRuntime := env.RuntimeEnvironment()
	userRuntime, _ := result["runtime"].(map[string]any)
	runtime := make(map[string]any, len(baseRuntime)+len(userRuntime))
	for k, v := range baseRuntime {
		runtime[k] = v
	}
	for k, v := range userRuntime {
		runtime[k] = v
	}
	result["runtime"] = runtime

	envMeta := env.LangChainEnvMetadata()
	if len(envMeta) > 0 {
		oldMeta, _ := result["metadata"].(map[string]any)
		metadata := make(map[string]any, len(envMeta)+len(oldMeta))
		for k, v := range envMeta {
			metadata[k] = v
		}
		for k, v := range oldMeta {
			metadata[k] = v
		}
		result["metadata"] = metadata
	}

	return result
}

func validateAttachmentNames(attachments map[string]Attachment) error {
	for name := range attachments {
		if strings.Contains(name, ".") {
			return fmt.Errorf("langsmith: attachment name %q must not contain '.'", name)
		}
	}
	return nil
}
