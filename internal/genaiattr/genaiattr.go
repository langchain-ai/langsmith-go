// Package genaiattr re-exports OpenTelemetry semconv v1.40.0 GenAI attribute
// keys/helpers alongside LangSmith-specific keys that are not (yet) part of
// the official semantic conventions.
//
// Standard gen_ai.* attributes come directly from the semconv package:
//
//	semconv.GenAIRequestModel(m)           → gen_ai.request.model = m
//	semconv.GenAIProviderNameGCPGemini     → gen_ai.provider.name = "gcp.gemini"
//
// LangSmith-specific and legacy keys are defined here because the semconv
// package does not cover them.
//
// See: https://opentelemetry.io/docs/specs/semconv/registry/attributes/gen-ai
package genaiattr

import "go.opentelemetry.io/otel/attribute"

// LangSmith-specific Gen AI attribute keys with no semconv equivalent.
//
// OTel semconv v1.37+ replaced the single-string prompt/completion model
// with structured per-message events (gen_ai.input.messages /
// gen_ai.output.messages). These keys are what the LangSmith OTLP converter
// reads; they were never part of an official semconv release.
//
// See: langchainplus/smith-go/otel/otel_converter.go (GenAIPrompt, GenAICompletion)
const (
	// PromptKey is a JSON-serialized string of input messages.
	// The converter uses this to populate run inputs.
	PromptKey = attribute.Key("gen_ai.prompt")

	// CompletionKey is a JSON-serialized string of output messages.
	// The converter uses this to populate run outputs.
	CompletionKey = attribute.Key("gen_ai.completion")

	// UsageTotalTokensKey is the total number of tokens used (input + output).
	UsageTotalTokensKey = attribute.Key("gen_ai.usage.total_tokens")

	// UsageReasoningTokensKey is the number of tokens used for reasoning/thinking.
	UsageReasoningTokensKey = attribute.Key("gen_ai.usage.details.reasoning_tokens")
)

// HTTP semantic convention attribute keys.
const (
	HTTPMethodKey = attribute.Key("http.method")
	HTTPURLKey    = attribute.Key("http.url")
)

// LangSmith-specific metadata attribute keys.
const (
	// StopReasonKey records the model's stop/finish reason in LangSmith metadata.
	StopReasonKey = attribute.Key("langsmith.metadata.stop_reason")
)

// Legacy Gen AI attribute keys that the LangSmith OTLP converter reads.
// These use underscore-separated names that predate the official semconv
// dotted format (gen_ai.usage.cache_read.input_tokens).
const (
	// CacheReadInputTokensKey is the legacy key the converter uses for cache-read tokens.
	CacheReadInputTokensKey = attribute.Key("gen_ai.usage.cache_read_input_tokens")

	// CacheCreationInputTokensKey is the legacy key the converter uses for cache-creation tokens.
	CacheCreationInputTokensKey = attribute.Key("gen_ai.usage.cache_creation_input_tokens")
)
