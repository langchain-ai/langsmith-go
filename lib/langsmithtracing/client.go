package langsmithtracing

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
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
// Use this to attach images, audio, PDFs, or any other binary content to traces.
type Attachment = models.Attachment

// DrainConfig controls batching and auto-scaling behavior.
// See [tracesink.DrainConfig] for field documentation.
type DrainConfig = tracesink.DrainConfig

// DefaultDrainConfig returns production-grade defaults for the trace sink.
func DefaultDrainConfig() DrainConfig { return tracesink.DefaultDrainConfig() }

// TracingClient sends runs to LangSmith via the multipart ingestion endpoint.
type TracingClient struct {
	sink          *tracesink.TraceSink
	writeEndpoint models.WriteEndpoint

	sampleRate     *float64
	filteredMu     sync.Mutex
	filteredTraces map[uuid.UUID]struct{} // traces that were sampled out
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
// All fields except ID, TraceID, EndTime, and DottedOrder are optional — only
// non-zero values are sent, matching the Python SDK's update_run behavior.
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

// RunTransformFunc is a pre-export transform hook matching the Python SDK's
// process_buffered_run_ops. It receives a batch of decoded run operations and
// returns (possibly modified) operations. The transform runs on every drain
// cycle, after batching but before coalescing and export.
type RunTransformFunc = tracesink.RunTransformFunc

// Option configures a TracingClient.
type Option func(*options)

type options struct {
	apiURL       string
	apiKey       string
	project      string
	drainConfig  *tracesink.DrainConfig
	sampleRate   *float64
	runTransform RunTransformFunc
}

// WithAPIURL overrides the LangSmith API URL.
func WithAPIURL(url string) Option { return func(o *options) { o.apiURL = url } }

// WithAPIKey overrides the LangSmith API key.
func WithAPIKey(key string) Option { return func(o *options) { o.apiKey = key } }

// WithProject overrides the LangSmith project name.
func WithProject(name string) Option { return func(o *options) { o.project = name } }

// WithDrainConfig overrides the default drain/scaling configuration.
func WithDrainConfig(config DrainConfig) Option {
	return func(o *options) { o.drainConfig = &config }
}

// WithSampleRate sets the trace sampling rate (0.0–1.0).
// 1.0 means send all traces, 0.0 means drop all, 0.5 means ~50%.
// Overrides the LANGSMITH_TRACING_SAMPLING_RATE env var.
func WithSampleRate(rate float64) Option {
	return func(o *options) { o.sampleRate = &rate }
}

// WithRunTransform sets a pre-export transform hook, matching the Python SDK's
// process_buffered_run_ops callback. The function receives each batch of run
// operations (decoded into maps) before coalescing and export, and can inspect
// or modify them. This is useful for enriching runs with external data,
// applying rate-limited APIs, or filtering sensitive fields in bulk.
func WithRunTransform(fn RunTransformFunc) Option {
	return func(o *options) { o.runTransform = fn }
}

// NewTracingClient creates a TracingClient that sends runs via multipart ingestion.
// The context is propagated to HTTP requests during normal operation; Close
// always drains with a background context to guarantee delivery.
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
		URL:     cfg.apiURL,
		Key:     cfg.apiKey,
		Project: cfg.project,
	}

	exp := multipart.NewExporter(nil, multipart.DefaultRetry())
	sink := tracesink.NewTraceSink(ctx, exp, drainCfg, endpoint, cfg.runTransform)

	return &TracingClient{
		sink:           sink,
		writeEndpoint:  endpoint,
		sampleRate:     sampleRate,
		filteredTraces: make(map[uuid.UUID]struct{}),
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

	sessionName := c.writeEndpoint.Project
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

	var inputsBytes, outputsBytes, extraBytes, eventsBytes, errorBytes, serializedBytes []byte
	if r.Inputs != nil {
		inputsBytes, err = json.Marshal(r.Inputs)
		if err != nil {
			return fmt.Errorf("marshal inputs: %w", err)
		}
	}
	if r.Outputs != nil {
		outputsBytes, err = json.Marshal(r.Outputs)
		if err != nil {
			return fmt.Errorf("marshal outputs: %w", err)
		}
	}
	if extra != nil {
		extraBytes, err = json.Marshal(extra)
		if err != nil {
			return fmt.Errorf("marshal extra: %w", err)
		}
	}
	if r.Events != nil {
		eventsBytes, err = json.Marshal(r.Events)
		if err != nil {
			return fmt.Errorf("marshal events: %w", err)
		}
	}
	if r.Error != "" {
		errorBytes, err = json.Marshal(r.Error)
		if err != nil {
			return fmt.Errorf("marshal error: %w", err)
		}
	}
	if r.Serialized != nil && (r.RunType == "llm" || r.RunType == "prompt") {
		filtered := make(map[string]any, len(r.Serialized))
		for k, v := range r.Serialized {
			if k != "graph" {
				filtered[k] = v
			}
		}
		serializedBytes, err = json.Marshal(filtered)
		if err != nil {
			return fmt.Errorf("marshal serialized: %w", err)
		}
	}

	op := &models.SerializedOp{
		Kind:        models.OpKindPost,
		ID:          r.ID,
		TraceID:     r.TraceID,
		RunInfo:     runInfoBytes,
		Inputs:      inputsBytes,
		Outputs:     outputsBytes,
		Events:      eventsBytes,
		Extra:       extraBytes,
		Error:       errorBytes,
		Serialized:  serializedBytes,
		Attachments: r.Attachments,
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
		"end_time":     r.EndTime.UTC().Format(time.RFC3339Nano),
		"dotted_order": r.DottedOrder,
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

	var inputsBytes, outputsBytes, eventsBytes, extraBytes, errorBytes []byte
	if r.Inputs != nil {
		inputsBytes, err = json.Marshal(r.Inputs)
		if err != nil {
			return fmt.Errorf("marshal inputs: %w", err)
		}
	}
	if r.Outputs != nil {
		outputsBytes, err = json.Marshal(r.Outputs)
		if err != nil {
			return fmt.Errorf("marshal outputs: %w", err)
		}
	}
	if r.Extra != nil {
		extraBytes, err = json.Marshal(r.Extra)
		if err != nil {
			return fmt.Errorf("marshal extra: %w", err)
		}
	}
	if r.Events != nil {
		eventsBytes, err = json.Marshal(r.Events)
		if err != nil {
			return fmt.Errorf("marshal events: %w", err)
		}
	}
	if r.Error != "" {
		errorBytes, err = json.Marshal(r.Error)
		if err != nil {
			return fmt.Errorf("marshal error: %w", err)
		}
	}

	op := &models.SerializedOp{
		Kind:        models.OpKindPatch,
		ID:          r.ID,
		TraceID:     r.TraceID,
		RunInfo:     runInfoBytes,
		Inputs:      inputsBytes,
		Outputs:     outputsBytes,
		Extra:       extraBytes,
		Events:      eventsBytes,
		Error:       errorBytes,
		Attachments: r.Attachments,
	}
	return c.sink.Submit(op)
}

// Close flushes pending operations and shuts down the client.
func (c *TracingClient) Close() {
	c.sink.Close()
}

// shouldSampleCreate decides whether a create (post) should be kept.
// For root runs (id == traceID), a random check is made against the sample rate.
// Child runs follow the decision of their root. If no sampling rate is set,
// all runs are kept.
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
		c.filteredTraces[traceID] = struct{}{}
		return false
	}
	return true
}

// shouldSampleUpdate decides whether an update (patch) should be kept.
// Updates for traces that were sampled out are dropped. Root-run patches
// clean up the filtered set to prevent unbounded growth.
func (c *TracingClient) shouldSampleUpdate(id, traceID uuid.UUID) bool {
	if c.sampleRate == nil {
		return true
	}

	c.filteredMu.Lock()
	defer c.filteredMu.Unlock()

	if _, filtered := c.filteredTraces[traceID]; !filtered {
		return true
	}
	if id == traceID {
		delete(c.filteredTraces, traceID)
	}
	return false
}

// mergeRuntimeEnv injects SDK/platform info into extra.runtime and filtered
// LANGCHAIN_*/LANGSMITH_* env vars into extra.metadata, matching the Python
// SDK's _insert_runtime_env. User-provided keys take precedence in both maps.
func mergeRuntimeEnv(extra map[string]any) map[string]any {
	if extra == nil {
		extra = make(map[string]any)
	}

	// extra.runtime: SDK info as base, user keys override.
	sdkRuntime := env.RuntimeEnvironment()
	userRuntime, _ := extra["runtime"].(map[string]any)
	merged := make(map[string]any, len(sdkRuntime)+len(userRuntime))
	for k, v := range sdkRuntime {
		merged[k] = v
	}
	for k, v := range userRuntime {
		merged[k] = v
	}
	extra["runtime"] = merged

	// extra.metadata: env vars as base, user keys override.
	envMeta := env.LangChainEnvMetadata()
	if len(envMeta) > 0 {
		metadata, _ := extra["metadata"].(map[string]any)
		if metadata == nil {
			metadata = make(map[string]any, len(envMeta))
		}
		for k, v := range envMeta {
			if _, exists := metadata[k]; !exists {
				metadata[k] = v
			}
		}
		extra["metadata"] = metadata
	}

	return extra
}

func validateAttachmentNames(attachments map[string]Attachment) error {
	for name := range attachments {
		if strings.Contains(name, ".") {
			return fmt.Errorf("langsmith: attachment name %q must not contain '.'", name)
		}
	}
	return nil
}
