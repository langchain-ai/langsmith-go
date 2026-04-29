package langsmith

import "github.com/langchain-ai/langsmith-go/lib/langsmithtracing"

// Re-exported tracing types so users can write langsmith.RunCreate{} etc.
// without importing the langsmithtracing package directly.

// TracingClient sends runs to LangSmith via the multipart ingestion endpoint.
type TracingClient = langsmithtracing.TracingClient

// RunCreate holds parameters for creating a new run (multipart post).
type RunCreate = langsmithtracing.RunCreate

// RunUpdate holds parameters for updating an existing run (multipart patch).
type RunUpdate = langsmithtracing.RunUpdate

// TracingAttachment is a binary file to upload alongside a run.
type TracingAttachment = langsmithtracing.Attachment

// TracingLogger is the interface used by [TracingClient] for diagnostic output.
type TracingLogger = langsmithtracing.Logger

// DrainConfig controls batching and auto-scaling behavior for the trace sink.
type DrainConfig = langsmithtracing.DrainConfig

// RunOp is a decoded run operation exposed to transform hooks.
type RunOp = langsmithtracing.RunOp

// RunTransformFunc is a pre-export transform hook.
type RunTransformFunc = langsmithtracing.RunTransformFunc

// TracingOption configures a TracingClient.
type TracingOption = langsmithtracing.Option

// DefaultDrainConfig returns production-grade defaults for the trace sink.
func DefaultDrainConfig() DrainConfig { return langsmithtracing.DefaultDrainConfig() }

// Tracing option constructors.
var (
	WithTracingAPIURL                     = langsmithtracing.WithAPIURL
	WithTracingAPIKey                     = langsmithtracing.WithAPIKey
	WithTracingProject                    = langsmithtracing.WithProject
	WithTracingDrain                      = langsmithtracing.WithDrainConfig
	WithSampleRate                        = langsmithtracing.WithSampleRate
	WithRunTransform                      = langsmithtracing.WithRunTransform
	WithMergeFilteredEnvIntoExtraMetadata = langsmithtracing.WithMergeFilteredEnvIntoExtraMetadata
	WithCompressionDisabled               = langsmithtracing.WithCompressionDisabled
	WithTracingLogger                     = langsmithtracing.WithLogger
)

// NewTracingClient creates a standalone TracingClient. Most users should use
// [Client.CreateRun] / [Client.UpdateRun] instead, which lazily initialize the
// underlying TracingClient on first use so REST-only clients pay no cost.
// It returns an error if tracing sampling env vars are set but invalid.
var NewTracingClient = langsmithtracing.NewTracingClient
