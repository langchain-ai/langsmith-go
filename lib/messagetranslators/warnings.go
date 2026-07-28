package messagetranslators

import (
	"fmt"
	"sort"
	"sync"
)

// WarningCode identifies a non-fatal conversion condition.
type WarningCode string

const (
	// WarningUnknownField reports a field that is not part of the recognized
	// source schema. It is schema-drift telemetry, not a guarantee that the
	// field was semantically harmless.
	WarningUnknownField WarningCode = "unknown_field"
	// WarningLossyConversion reports recognized source metadata that is
	// intentionally dropped because it has no destination representation.
	WarningLossyConversion WarningCode = "lossy_conversion"
)

// Warning describes non-fatal schema drift or a lossy mapping. Path is the full
// JSON-style path to Field in the source wire object. Warning intentionally
// omits the source value so that callers can safely use these in log and metric
// labels.
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

// WarningCollector is a concurrency-safe WarningHandler target.
type WarningCollector struct {
	mu       sync.Mutex
	warnings []Warning
}

// HandleWarning records w. It can be passed to WithWarningHandler.
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

type wireInspector struct{ cfg config }

// fields builds an allowlist. Every call site is a package-level var: these are
// consulted once per inspected object, so building them per object would put an
// allocation on the hot path of every streamed event.
func fields(names ...string) map[string]struct{} {
	m := make(map[string]struct{}, len(names))
	for _, name := range names {
		m[name] = struct{}{}
	}
	return m
}

func (i wireInspector) object(o map[string]any, path string, allowed map[string]struct{}) {
	if i.cfg.warnings == nil || o == nil {
		return
	}
	// Nil until proven otherwise: almost every real object is fully recognized,
	// and that case should not allocate.
	var unknown []string
	for key := range o {
		if _, ok := allowed[key]; !ok {
			unknown = append(unknown, key)
		}
	}
	if len(unknown) == 0 {
		return
	}
	sort.Strings(unknown)
	for _, key := range unknown {
		i.cfg.unknownField(path+"."+key, key)
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

func inspectAnthropicStopDetails(i wireInspector, o map[string]any, path string) {
	if details, present := o["stop_details"]; present && details != nil {
		i.cfg.lossy(path+".stop_details", "stop_details",
			"stop_details at "+path+".stop_details was dropped because it has no destination equivalent")
	}
}

func inspectAnthropicCaller(i wireInspector, block map[string]any, path string) {
	value, present := block["caller"]
	if !present || value == nil {
		return
	}
	caller, isObject := obj(value)
	if isObject {
		i.object(caller, path+".caller", anthropicCallerFields)
		if len(caller) == 1 && caller["type"] == "direct" {
			return
		}
	}
	i.cfg.lossy(path+".caller", "caller",
		"caller at "+path+".caller was dropped because its invocation metadata has no destination equivalent")
}

var (
	anthropicEnvelopeFields        = fields("model", "messages", "max_tokens", "system", "temperature", "top_p", "top_k", "stop_sequences", "stream", "tools", "tool_choice", "metadata", "thinking", "service_tier", "output_config", "cache_control", "mcp_servers", "context_management", "container", "inference_geo")
	anthropicResponseFields        = fields("id", "type", "role", "model", "content", "stop_reason", "stop_sequence", "stop_details", "usage", "cache_control", "phase", "container")
	anthropicUsageFields           = fields("input_tokens", "output_tokens", "cache_read_input_tokens", "cache_creation_input_tokens", "cache_creation", "server_tool_use", "service_tier", "inference_geo")
	anthropicSystemBlockFields     = fields("type", "text", "cache_control", "citations")
	anthropicMessageFields         = fields("role", "content", "name")
	anthropicToolFields            = fields("name", "description", "input_schema", "type", "cache_control", "strict", "defer_loading", "allowed_callers")
	anthropicToolChoiceFields      = fields("type", "name", "disable_parallel_tool_use")
	anthropicTextBlockFields       = fields("type", "text", "citations", "cache_control", "phase")
	anthropicImageBlockFields      = fields("type", "source", "cache_control", "phase")
	anthropicImageSourceFields     = fields("type", "url", "media_type", "data")
	anthropicToolUseBlockFields    = fields("type", "id", "name", "input", "caller", "cache_control", "phase")
	anthropicCallerFields          = fields("type", "tool_id")
	anthropicToolResultBlockFields = fields("type", "tool_use_id", "content", "is_error", "cache_control", "phase")
	anthropicToolResultPartFields  = fields("type", "text", "citations", "cache_control")
	anthropicCacheCreationFields   = fields("ephemeral_5m_input_tokens", "ephemeral_1h_input_tokens")
	anthropicServerToolUseFields   = fields("web_search_requests", "web_fetch_requests")
)

func inspectAnthropicObject(root map[string]any, response bool, cfg config, base string) {
	i := wireInspector{cfg}
	if response {
		i.object(root, base, anthropicResponseFields)
		inspectAnthropicStopDetails(i, root, base)
		inspectAnthropicContent(i, root["content"], base+".content")
		inspectAnthropicUsage(i, root["usage"], base+".usage")
		return
	}
	i.object(root, base, anthropicEnvelopeFields)
	if system, ok := arr(root["system"]); ok {
		for n, x := range system {
			if o, ok := obj(x); ok {
				i.object(o, fmt.Sprintf("%s.system[%d]", base, n), anthropicSystemBlockFields)
			}
		}
	}
	eachObject(root["messages"], func(n int, m map[string]any) {
		p := fmt.Sprintf("%s.messages[%d]", base, n)
		i.object(m, p, anthropicMessageFields)
		inspectAnthropicContent(i, m["content"], p+".content")
	})
	eachObject(root["tools"], func(n int, tool map[string]any) {
		i.object(tool, fmt.Sprintf("%s.tools[%d]", base, n), anthropicToolFields)
	})
	if choice, ok := obj(root["tool_choice"]); ok {
		i.object(choice, base+".tool_choice", anthropicToolChoiceFields)
	}
}

func inspectAnthropicContent(i wireInspector, v any, path string) {
	eachObject(v, func(n int, block map[string]any) {
		inspectAnthropicContentBlock(i, block, fmt.Sprintf("%s[%d]", path, n))
	})
}

func inspectAnthropicContentBlock(i wireInspector, block map[string]any, path string) {
	switch block["type"] {
	case "text":
		i.object(block, path, anthropicTextBlockFields)
	case "image":
		i.object(block, path, anthropicImageBlockFields)
		if source, ok := obj(block["source"]); ok {
			i.object(source, path+".source", anthropicImageSourceFields)
		}
	case "tool_use":
		i.object(block, path, anthropicToolUseBlockFields)
		inspectAnthropicCaller(i, block, path)
	case "tool_result":
		i.object(block, path, anthropicToolResultBlockFields)
		eachObject(block["content"], func(j int, part map[string]any) {
			i.object(part, fmt.Sprintf("%s.content[%d]", path, j), anthropicToolResultPartFields)
		})
	default:
		// The discriminator is validated by the converter; do not describe the
		// rest of an object whose schema is not recognized as harmless drift.
		return
	}
	if _, ok := block["cache_control"]; ok {
		i.cfg.lossy(path+".cache_control", "cache_control",
			"cache_control at "+path+" has no destination equivalent; cache behavior and billing may differ")
	}
}

func inspectAnthropicUsage(i wireInspector, v any, path string) {
	u, ok := obj(v)
	if !ok {
		return
	}
	i.object(u, path, anthropicUsageFields)
	if details, ok := obj(u["cache_creation"]); ok {
		i.object(details, path+".cache_creation", anthropicCacheCreationFields)
	}
	if details, ok := obj(u["server_tool_use"]); ok {
		i.object(details, path+".server_tool_use", anthropicServerToolUseFields)
	}
}

var (
	responsesRequestFields            = fields("model", "input", "instructions", "max_output_tokens", "temperature", "top_p", "stream", "stream_options", "tools", "tool_choice", "previous_response_id", "reasoning", "truncation", "max_tool_calls", "parallel_tool_calls", "text", "phase", "prompt_cache_key", "prompt_cache_retention", "metadata", "store", "include", "service_tier", "user", "safety_identifier", "top_logprobs", "background", "conversation")
	responsesResponseFields           = fields("id", "object", "created_at", "completed_at", "status", "error", "incomplete_details", "instructions", "max_output_tokens", "max_tool_calls", "metadata", "model", "output", "parallel_tool_calls", "previous_response_id", "prompt_cache_key", "prompt_cache_retention", "reasoning", "safety_identifier", "service_tier", "store", "temperature", "text", "tool_choice", "tools", "top_logprobs", "top_p", "truncation", "usage", "user", "phase", "background", "conversation", "billing")
	responsesUsageFields              = fields("input_tokens", "output_tokens", "total_tokens", "input_tokens_details", "output_tokens_details")
	responsesInputDetailFields        = fields("cached_tokens")
	responsesOutputDetailFields       = fields("reasoning_tokens")
	responsesIncompleteDetailFields   = fields("reason")
	responsesErrorFields              = fields("code", "message", "param", "type")
	responsesStreamOptionFields       = fields("include_obfuscation")
	responsesToolFields               = fields("type", "name", "description", "parameters", "strict", "cache_control")
	responsesToolChoiceFields         = fields("type", "name")
	responsesMessageItemFields        = fields("id", "type", "role", "content", "status", "phase")
	responsesFunctionCallItemFields   = fields("id", "type", "call_id", "name", "arguments", "status", "phase")
	responsesFunctionOutputItemFields = fields("id", "type", "call_id", "output", "status", "phase")
	responsesInputTextPartFields      = fields("type", "text", "phase")
	responsesOutputTextPartFields     = fields("type", "text", "annotations", "logprobs", "phase")
	responsesInputImagePartFields     = fields("type", "image_url", "detail", "file_id", "phase")
)

func inspectResponsesObject(root map[string]any, response bool, cfg config, base string) {
	i := wireInspector{cfg}
	if response {
		i.object(root, base, responsesResponseFields)
		inspectResponsesItems(i, root["output"], base+".output")
		inspectUsage(i, root["usage"], base+".usage", responsesUsageFields, responsesInputDetailFields, responsesOutputDetailFields)
		if d, ok := obj(root["incomplete_details"]); ok {
			i.object(d, base+".incomplete_details", responsesIncompleteDetailFields)
		}
		if e, ok := obj(root["error"]); ok {
			i.object(e, base+".error", responsesErrorFields)
		}
		return
	}
	i.object(root, base, responsesRequestFields)
	if streamOptions, ok := obj(root["stream_options"]); ok {
		i.object(streamOptions, base+".stream_options", responsesStreamOptionFields)
	}
	inspectResponsesItems(i, root["input"], base+".input")
	eachObject(root["tools"], func(n int, tool map[string]any) {
		i.object(tool, fmt.Sprintf("%s.tools[%d]", base, n), responsesToolFields)
	})
	if choice, ok := obj(root["tool_choice"]); ok {
		i.object(choice, base+".tool_choice", responsesToolChoiceFields)
	}
}

func inspectResponsesItems(i wireInspector, v any, path string) {
	eachObject(v, func(n int, item map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch item["type"] {
		case nil, "", "message":
			i.object(item, p, responsesMessageItemFields)
			inspectResponsesParts(i, item["content"], p+".content")
		case "function_call":
			i.object(item, p, responsesFunctionCallItemFields)
		case "function_call_output":
			i.object(item, p, responsesFunctionOutputItemFields)
		}
	})
}

func inspectResponsesParts(i wireInspector, v any, path string) {
	eachObject(v, func(n int, part map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch part["type"] {
		case "input_text":
			i.object(part, p, responsesInputTextPartFields)
		case "output_text":
			i.object(part, p, responsesOutputTextPartFields)
		case "input_image":
			i.object(part, p, responsesInputImagePartFields)
		}
	})
}

var (
	chatCompletionsRequestFields            = fields("model", "messages", "max_tokens", "max_completion_tokens", "temperature", "top_p", "stop", "stream", "stream_options", "tools", "tool_choice", "n", "functions", "function_call", "response_format", "parallel_tool_calls", "logprobs", "top_logprobs", "prediction", "modalities", "audio", "service_tier", "store", "metadata", "seed", "user", "reasoning_effort", "web_search_options", "frequency_penalty", "presence_penalty", "logit_bias", "verbosity", "safety_identifier")
	chatCompletionsResponseFields           = fields("id", "object", "created", "model", "choices", "usage", "service_tier", "system_fingerprint")
	chatCompletionsUsageFields              = fields("prompt_tokens", "completion_tokens", "total_tokens", "prompt_tokens_details", "completion_tokens_details")
	chatCompletionsInputDetailFields        = fields("cached_tokens", "audio_tokens")
	chatCompletionsOutputDetailFields       = fields("reasoning_tokens", "audio_tokens", "accepted_prediction_tokens", "rejected_prediction_tokens")
	chatCompletionsStreamOptionFields       = fields("include_usage", "include_obfuscation")
	chatCompletionsMessageFields            = fields("role", "content", "name", "tool_calls", "tool_call_id", "function_call", "audio", "refusal", "annotations")
	chatCompletionsToolFields               = fields("type", "function", "strict")
	chatCompletionsToolFunctionFields       = fields("name", "description", "parameters", "strict")
	chatCompletionsToolChoiceFields         = fields("type", "function")
	chatCompletionsToolChoiceFunctionFields = fields("name")
	chatCompletionsTextPartFields           = fields("type", "text")
	chatCompletionsImagePartFields          = fields("type", "image_url")
	chatCompletionsImageURLFields           = fields("url", "detail")
	chatCompletionsAudioPartFields          = fields("type", "input_audio")
	chatCompletionsFilePartFields           = fields("type", "file")
	chatCompletionsRefusalPartFields        = fields("type", "refusal")
	chatCompletionsStreamChoiceFields       = fields("index", "delta", "finish_reason", "logprobs")
	chatCompletionsDeltaFields              = fields("role", "content", "tool_calls", "audio", "function_call", "refusal")
	chatCompletionsChoiceFields             = fields("index", "message", "finish_reason", "logprobs")
	chatCompletionsChoiceMessageFields      = fields("role", "content", "tool_calls", "function_call", "audio", "refusal", "annotations")
	chatCompletionsToolCallFields           = fields("id", "type", "function")
	chatCompletionsStreamToolCallFields     = fields("id", "type", "function", "index")
	chatCompletionsToolCallFunctionFields   = fields("name", "arguments")
)

func inspectChatCompletionsObject(root map[string]any, response bool, cfg config, base string) {
	i := wireInspector{cfg}
	if response {
		i.object(root, base, chatCompletionsResponseFields)
		inspectChatCompletionsChoices(i, root["choices"], base+".choices", false)
		inspectUsage(i, root["usage"], base+".usage", chatCompletionsUsageFields, chatCompletionsInputDetailFields, chatCompletionsOutputDetailFields)
		return
	}
	i.object(root, base, chatCompletionsRequestFields)
	if so, ok := obj(root["stream_options"]); ok {
		i.object(so, base+".stream_options", chatCompletionsStreamOptionFields)
	}
	eachObject(root["messages"], func(n int, m map[string]any) {
		p := fmt.Sprintf("%s.messages[%d]", base, n)
		i.object(m, p, chatCompletionsMessageFields)
		inspectChatCompletionsParts(i, m["content"], p+".content")
		inspectChatCompletionsToolCalls(i, m["tool_calls"], p+".tool_calls", false)
	})
	eachObject(root["tools"], func(n int, tool map[string]any) {
		p := fmt.Sprintf("%s.tools[%d]", base, n)
		i.object(tool, p, chatCompletionsToolFields)
		if f, ok := obj(tool["function"]); ok {
			i.object(f, p+".function", chatCompletionsToolFunctionFields)
		}
	})
	if tc, ok := obj(root["tool_choice"]); ok {
		i.object(tc, base+".tool_choice", chatCompletionsToolChoiceFields)
		if f, ok := obj(tc["function"]); ok {
			i.object(f, base+".tool_choice.function", chatCompletionsToolChoiceFunctionFields)
		}
	}
}

func inspectChatCompletionsParts(i wireInspector, v any, path string) {
	eachObject(v, func(n int, part map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		switch part["type"] {
		case "text":
			i.object(part, p, chatCompletionsTextPartFields)
		case "image_url":
			i.object(part, p, chatCompletionsImagePartFields)
			if image, ok := obj(part["image_url"]); ok {
				i.object(image, p+".image_url", chatCompletionsImageURLFields)
			}
		case "input_audio":
			i.object(part, p, chatCompletionsAudioPartFields)
		case "file":
			i.object(part, p, chatCompletionsFilePartFields)
		case "refusal":
			i.object(part, p, chatCompletionsRefusalPartFields)
		}
	})
}

func inspectChatCompletionsChoices(i wireInspector, v any, path string, stream bool) {
	eachObject(v, func(n int, choice map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		if stream {
			i.object(choice, p, chatCompletionsStreamChoiceFields)
			if d, ok := obj(choice["delta"]); ok {
				i.object(d, p+".delta", chatCompletionsDeltaFields)
				inspectChatCompletionsToolCalls(i, d["tool_calls"], p+".delta.tool_calls", true)
			}
			return
		}
		i.object(choice, p, chatCompletionsChoiceFields)
		if m, ok := obj(choice["message"]); ok {
			i.object(m, p+".message", chatCompletionsChoiceMessageFields)
			inspectChatCompletionsParts(i, m["content"], p+".message.content")
			inspectChatCompletionsToolCalls(i, m["tool_calls"], p+".message.tool_calls", false)
		}
	})
}

func inspectChatCompletionsToolCalls(i wireInspector, v any, path string, stream bool) {
	allowed := chatCompletionsToolCallFields
	if stream {
		allowed = chatCompletionsStreamToolCallFields
	}
	eachObject(v, func(n int, call map[string]any) {
		p := fmt.Sprintf("%s[%d]", path, n)
		i.object(call, p, allowed)
		if f, ok := obj(call["function"]); ok {
			i.object(f, p+".function", chatCompletionsToolCallFunctionFields)
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
