package messagetranslators

import (
	"fmt"
	"sort"
	"sync"
)

// WarningCode identifies a non-fatal conversion condition.
type WarningCode string

const (
	// UnknownField reports a field that is not part of the recognized source schema.
	UnknownField WarningCode = "unknown_field"
	// LossyConversion is reserved for validated source data that is intentionally dropped.
	LossyConversion WarningCode = "lossy_conversion"

	// WarningUnknownField is an explicit-name alias for UnknownField.
	WarningUnknownField = UnknownField
	// WarningLossyConversion is an explicit-name alias for LossyConversion.
	WarningLossyConversion = LossyConversion
)

// Warning describes non-fatal schema drift or a lossy mapping. Path is the full
// JSON-style path to Field in the source wire object.
type Warning struct {
	Code    WarningCode
	Path    string
	Field   string
	Message string
}

// WarningHandler receives warnings synchronously, in source traversal order.
// A handler must not panic. Panics are deliberately not recovered: silently
// continuing after a callback panic can hide a broken gateway policy.
type WarningHandler func(Warning)

// ConversionOptions configures one conversion or stream converter. The zero
// value preserves the historical behavior of silently tolerating unknown fields.
type ConversionOptions struct {
	WarningHandler WarningHandler
}

func (o ConversionOptions) warn(w Warning) {
	if o.WarningHandler != nil {
		o.WarningHandler(w)
	}
}

// WarningCollector is a concurrency-safe WarningHandler target.
type WarningCollector struct {
	mu       sync.Mutex
	warnings []Warning
}

// HandleWarning records w. It can be passed as ConversionOptions.WarningHandler.
func (c *WarningCollector) HandleWarning(w Warning) {
	c.mu.Lock()
	c.warnings = append(c.warnings, w)
	c.mu.Unlock()
}

// Warnings returns a copy of all warnings collected so far.
func (c *WarningCollector) Warnings() []Warning {
	c.mu.Lock()
	defer c.mu.Unlock()
	return append([]Warning(nil), c.warnings...)
}

type wireInspector struct{ options ConversionOptions }

func fields(names ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(names))
	for _, name := range names {
		m[name] = struct{}{}
	}
	return m
}

func (i wireInspector) object(o map[string]any, path string, allowed map[string]struct{}) {
	if i.options.WarningHandler == nil || o == nil {
		return
	}
	unknown := make([]string, 0)
	for key := range o {
		if _, ok := allowed[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	sort.Strings(unknown)
	for _, key := range unknown {
		p := path + "." + key
		i.options.warn(Warning{Code: WarningUnknownField, Path: p, Field: key, Message: fmt.Sprintf("unknown source field %q at %s was ignored", key, p)})
	}
}

func eachObject(v any, fn func(int, map[string]any)) {
	if a, ok := arr(v); ok {
		for n, x := range a {
			if o, ok := obj(x); ok {
				fn(n, o)
			}
		}
	}
}

var anthropicEnvelopeFields = fields("model", "messages", "max_tokens", "system", "temperature", "top_p", "top_k", "stop_sequences", "stream", "tools", "tool_choice", "metadata", "thinking", "service_tier", "output_config", "cache_control", "mcp_servers", "context_management", "container", "inference_geo")
var anthropicResponseFields = fields("id", "type", "role", "model", "content", "stop_reason", "stop_sequence", "usage", "cache_control", "phase", "container")
var anthropicUsageFields = fields("input_tokens", "output_tokens", "cache_read_input_tokens", "cache_creation_input_tokens", "cache_creation", "server_tool_use", "service_tier", "inference_geo")

func inspectAnthropicObject(root map[string]any, response bool, options ConversionOptions, base string) {
	i := wireInspector{options}
	if response {
		i.object(root, base, anthropicResponseFields)
		inspectAnthropicContent(i, root["content"], base+".content")
		inspectAnthropicUsage(i, root["usage"], base+".usage")
		return
	}
	i.object(root, base, anthropicEnvelopeFields)
	if system, ok := arr(root["system"]); ok {
		for n, x := range system {
			if o, ok := obj(x); ok {
				i.object(o, fmt.Sprintf("%s.system[%d]", base, n), fields("type", "text", "cache_control", "citations"))
			}
		}
	}
	eachObject(root["messages"], func(n int, m map[string]any) {
		p := fmt.Sprintf("%s.messages[%d]", base, n)
		i.object(m, p, fields("role", "content", "name"))
		inspectAnthropicContent(i, m["content"], p+".content")
	})
	eachObject(root["tools"], func(n int, tool map[string]any) {
		i.object(tool, fmt.Sprintf("%s.tools[%d]", base, n), fields("name", "description", "input_schema", "type", "cache_control", "strict", "defer_loading", "allowed_callers"))
	})
	if choice, ok := obj(root["tool_choice"]); ok {
		i.object(choice, base+".tool_choice", fields("type", "name", "disable_parallel_tool_use"))
	}
}

func inspectAnthropicContent(i wireInspector, v any, path string) {
	eachObject(v, func(n int, block map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch block["type"] {
		case "text":
			i.object(block, p, fields("type", "text", "citations", "cache_control", "phase"))
		case "image":
			i.object(block, p, fields("type", "source", "cache_control", "phase"))
			if source, ok := obj(block["source"]); ok {
				i.object(source, p+".source", fields("type", "url", "media_type", "data"))
			}
		case "tool_use":
			i.object(block, p, fields("type", "id", "name", "input", "cache_control", "phase"))
		case "tool_result":
			i.object(block, p, fields("type", "tool_use_id", "content", "is_error", "cache_control", "phase"))
			eachObject(block["content"], func(j int, part map[string]any) {
				i.object(part, fmt.Sprintf("%s.content[%d]", p, j), fields("type", "text", "citations", "cache_control"))
			})
		default:
			// The discriminator is validated by the converter; do not describe the
			// rest of an object whose schema is not recognized as harmless drift.
		}
	})
}

func inspectAnthropicUsage(i wireInspector, v any, path string) {
	u, ok := obj(v)
	if !ok {
		return
	}
	i.object(u, path, anthropicUsageFields)
	if details, ok := obj(u["cache_creation"]); ok {
		i.object(details, path+".cache_creation", fields("ephemeral_5m_input_tokens", "ephemeral_1h_input_tokens"))
	}
	if details, ok := obj(u["server_tool_use"]); ok {
		i.object(details, path+".server_tool_use", fields("web_search_requests", "web_fetch_requests"))
	}
}

var responsesRequestFields = fields("model", "input", "instructions", "max_output_tokens", "temperature", "top_p", "stream", "stream_options", "tools", "tool_choice", "previous_response_id", "reasoning", "truncation", "max_tool_calls", "parallel_tool_calls", "text", "phase", "prompt_cache_key", "prompt_cache_retention", "metadata", "store", "include", "service_tier", "user", "safety_identifier", "top_logprobs", "background", "conversation")
var responsesResponseFields = fields("id", "object", "created_at", "completed_at", "status", "error", "incomplete_details", "instructions", "max_output_tokens", "max_tool_calls", "metadata", "model", "output", "parallel_tool_calls", "previous_response_id", "prompt_cache_key", "prompt_cache_retention", "reasoning", "safety_identifier", "service_tier", "store", "temperature", "text", "tool_choice", "tools", "top_logprobs", "top_p", "truncation", "usage", "user", "phase", "background", "conversation", "billing")
var responsesUsageFields = fields("input_tokens", "output_tokens", "total_tokens", "input_tokens_details", "output_tokens_details")

func inspectResponsesObject(root map[string]any, response bool, options ConversionOptions, base string) {
	i := wireInspector{options}
	if response {
		i.object(root, base, responsesResponseFields)
		inspectResponsesItems(i, root["output"], base+".output")
		inspectUsage(i, root["usage"], base+".usage", responsesUsageFields, fields("cached_tokens"), fields("reasoning_tokens"))
		if d, ok := obj(root["incomplete_details"]); ok {
			i.object(d, base+".incomplete_details", fields("reason"))
		}
		if e, ok := obj(root["error"]); ok {
			i.object(e, base+".error", fields("code", "message", "param", "type"))
		}
		return
	}
	i.object(root, base, responsesRequestFields)
	if streamOptions, ok := obj(root["stream_options"]); ok {
		i.object(streamOptions, base+".stream_options", fields("include_obfuscation"))
	}
	inspectResponsesItems(i, root["input"], base+".input")
	eachObject(root["tools"], func(n int, tool map[string]any) {
		i.object(tool, fmt.Sprintf("%s.tools[%d]", base, n), fields("type", "name", "description", "parameters", "strict", "cache_control"))
	})
	if choice, ok := obj(root["tool_choice"]); ok {
		i.object(choice, base+".tool_choice", fields("type", "name"))
	}
}

func inspectResponsesItems(i wireInspector, v any, path string) {
	eachObject(v, func(n int, item map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch item["type"] {
		case nil, "", "message":
			i.object(item, p, fields("id", "type", "role", "content", "status", "phase"))
			inspectResponsesParts(i, item["content"], p+".content")
		case "function_call":
			i.object(item, p, fields("id", "type", "call_id", "name", "arguments", "status", "phase"))
		case "function_call_output":
			i.object(item, p, fields("id", "type", "call_id", "output", "status", "phase"))
		}
	})
}

func inspectResponsesParts(i wireInspector, v any, path string) {
	eachObject(v, func(n int, part map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch part["type"] {
		case "input_text":
			i.object(part, p, fields("type", "text", "phase"))
		case "output_text":
			i.object(part, p, fields("type", "text", "annotations", "logprobs", "phase"))
		case "input_image":
			i.object(part, p, fields("type", "image_url", "detail", "file_id", "phase"))
		}
	})
}

var chatCompletionsRequestFields = fields("model", "messages", "max_tokens", "max_completion_tokens", "temperature", "top_p", "stop", "stream", "stream_options", "tools", "tool_choice", "n", "functions", "function_call", "response_format", "parallel_tool_calls", "logprobs", "top_logprobs", "prediction", "modalities", "audio", "service_tier", "store", "metadata", "seed", "user", "reasoning_effort", "web_search_options", "frequency_penalty", "presence_penalty", "logit_bias", "verbosity", "safety_identifier")
var chatCompletionsResponseFields = fields("id", "object", "created", "model", "choices", "usage", "service_tier", "system_fingerprint")
var chatCompletionsUsageFields = fields("prompt_tokens", "completion_tokens", "total_tokens", "prompt_tokens_details", "completion_tokens_details")

func inspectChatCompletionsObject(root map[string]any, response bool, options ConversionOptions, base string) {
	i := wireInspector{options}
	if response {
		i.object(root, base, chatCompletionsResponseFields)
		inspectChatCompletionsChoices(i, root["choices"], base+".choices", false)
		inspectUsage(i, root["usage"], base+".usage", chatCompletionsUsageFields, fields("cached_tokens", "audio_tokens"), fields("reasoning_tokens", "audio_tokens", "accepted_prediction_tokens", "rejected_prediction_tokens"))
		return
	}
	i.object(root, base, chatCompletionsRequestFields)
	if so, ok := obj(root["stream_options"]); ok {
		i.object(so, base+".stream_options", fields("include_usage", "include_obfuscation"))
	}
	eachObject(root["messages"], func(n int, m map[string]any) {
		p := fmt.Sprintf("%s.messages[%d]", base, n)
		i.object(m, p, fields("role", "content", "name", "tool_calls", "tool_call_id", "function_call", "audio", "refusal", "annotations"))
		inspectChatCompletionsParts(i, m["content"], p+".content")
		inspectChatCompletionsToolCalls(i, m["tool_calls"], p+".tool_calls", false)
	})
	eachObject(root["tools"], func(n int, tool map[string]any) {
		p := fmt.Sprintf("%s.tools[%d]", base, n)
		i.object(tool, p, fields("type", "function", "strict"))
		if f, ok := obj(tool["function"]); ok {
			i.object(f, p+".function", fields("name", "description", "parameters", "strict"))
		}
	})
	if tc, ok := obj(root["tool_choice"]); ok {
		i.object(tc, base+".tool_choice", fields("type", "function"))
		if f, ok := obj(tc["function"]); ok {
			i.object(f, base+".tool_choice.function", fields("name"))
		}
	}
}

func inspectChatCompletionsParts(i wireInspector, v any, path string) {
	eachObject(v, func(n int, part map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch part["type"] {
		case "text":
			i.object(part, p, fields("type", "text"))
		case "image_url":
			i.object(part, p, fields("type", "image_url"))
			if image, ok := obj(part["image_url"]); ok {
				i.object(image, p+".image_url", fields("url", "detail"))
			}
		case "input_audio":
			i.object(part, p, fields("type", "input_audio"))
		case "file":
			i.object(part, p, fields("type", "file"))
		case "refusal":
			i.object(part, p, fields("type", "refusal"))
		}
	})
}

func inspectChatCompletionsChoices(i wireInspector, v any, path string, stream bool) {
	eachObject(v, func(n int, choice map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		if stream {
			i.object(choice, p, fields("index", "delta", "finish_reason", "logprobs"))
			if d, ok := obj(choice["delta"]); ok {
				i.object(d, p+".delta", fields("role", "content", "tool_calls", "audio", "function_call", "refusal"))
				inspectChatCompletionsToolCalls(i, d["tool_calls"], p+".delta.tool_calls", true)
			}
		} else {
			i.object(choice, p, fields("index", "message", "finish_reason", "logprobs"))
			if m, ok := obj(choice["message"]); ok {
				i.object(m, p+".message", fields("role", "content", "tool_calls", "function_call", "audio", "refusal", "annotations"))
				inspectChatCompletionsParts(i, m["content"], p+".message.content")
				inspectChatCompletionsToolCalls(i, m["tool_calls"], p+".message.tool_calls", false)
			}
		}
	})
}

func inspectChatCompletionsToolCalls(i wireInspector, v any, path string, stream bool) {
	eachObject(v, func(n int, call map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		allowed := fields("id", "type", "function")
		if stream {
			allowed["index"] = struct{}{}
		}
		i.object(call, p, allowed)
		if f, ok := obj(call["function"]); ok {
			i.object(f, p+".function", fields("name", "arguments"))
		}
	})
}

func inspectUsage(i wireInspector, v any, path string, allowed, inputDetails, outputDetails map[string]struct{}) {
	u, ok := obj(v)
	if !ok {
		return
	}
	i.object(u, path, allowed)
	if inputDetails != nil {
		for _, key := range []string{"input_tokens_details", "prompt_tokens_details"} {
			if d, ok := obj(u[key]); ok {
				i.object(d, path+"."+key, inputDetails)
			}
		}
	}
	if outputDetails != nil {
		for _, key := range []string{"output_tokens_details", "completion_tokens_details"} {
			if d, ok := obj(u[key]); ok {
				i.object(d, path+"."+key, outputDetails)
			}
		}
	}
}
