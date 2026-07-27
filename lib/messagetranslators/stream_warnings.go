package messagetranslators

func eventObjectForWarnings(e SSEEvent) (map[string]any, string, bool) {
	if string(e.Data) == "[DONE]" {
		return nil, "", false
	}
	o, err := decodeObject(e.Data)
	if err != nil {
		return nil, "", false
	}
	typ := e.Event
	if typ == "" {
		typ, _ = str(o["type"])
	}
	return o, typ, true
}

func inspectAnthropicEvent(e SSEEvent, options ConversionOptions) {
	o, typ, ok := eventObjectForWarnings(e)
	if !ok {
		return
	}
	i := wireInspector{options}
	switch typ {
	case "message_start":
		i.object(o, "event", fields("type", "message"))
		if m, ok := obj(o["message"]); ok {
			inspectAnthropicObject(m, true, options, "event.message")
		}
	case "content_block_start":
		i.object(o, "event", fields("type", "index", "content_block"))
		if b, ok := obj(o["content_block"]); ok {
			switch b["type"] {
			case "text":
				i.object(b, "event.content_block", fields("type", "text", "citations", "cache_control", "phase"))
			case "tool_use":
				i.object(b, "event.content_block", fields("type", "id", "name", "input", "cache_control", "phase"))
			}
		}
	case "content_block_delta":
		i.object(o, "event", fields("type", "index", "delta"))
		if d, ok := obj(o["delta"]); ok {
			switch d["type"] {
			case "text_delta":
				i.object(d, "event.delta", fields("type", "text", "citations"))
			case "input_json_delta":
				i.object(d, "event.delta", fields("type", "partial_json"))
			case "citations_delta", "citation_delta":
				i.object(d, "event.delta", fields("type", "citation"))
			}
		}
	case "content_block_stop":
		i.object(o, "event", fields("type", "index"))
	case "message_delta":
		i.object(o, "event", fields("type", "delta", "usage"))
		if d, ok := obj(o["delta"]); ok {
			i.object(d, "event.delta", fields("stop_reason", "stop_sequence"))
		}
		inspectAnthropicUsage(i, o["usage"], "event.usage")
	case "message_stop", "ping":
		i.object(o, "event", fields("type"))
	case "error":
		i.object(o, "event", fields("type", "error"))
		if er, ok := obj(o["error"]); ok {
			i.object(er, "event.error", fields("type", "code", "message"))
		}
	}
}

func inspectResponsesEvent(e SSEEvent, options ConversionOptions) {
	o, typ, ok := eventObjectForWarnings(e)
	if !ok {
		return
	}
	i := wireInspector{options}
	common := []string{"type", "sequence_number"}
	allow := func(names ...string) map[string]struct{} { return fields(append(common, names...)...) }
	switch typ {
	case "response.created", "response.in_progress", "response.completed", "response.incomplete", "response.failed":
		i.object(o, "event", allow("response"))
		if r, ok := obj(o["response"]); ok {
			inspectResponsesObject(r, true, options, "event.response")
		}
	case "response.output_item.added", "response.output_item.done":
		i.object(o, "event", allow("output_index", "item"))
		if item, ok := obj(o["item"]); ok {
			inspectResponseEventItem(i, item, "event.item")
		}
	case "response.content_part.added", "response.content_part.done":
		i.object(o, "event", allow("item_id", "output_index", "content_index", "part"))
		if part, ok := obj(o["part"]); ok {
			inspectResponseEventPart(i, part, "event.part")
		}
	case "response.output_text.delta":
		i.object(o, "event", allow("item_id", "output_index", "content_index", "delta", "logprobs", "obfuscation"))
	case "response.output_text.done":
		i.object(o, "event", allow("item_id", "output_index", "content_index", "text", "annotations", "logprobs"))
	case "response.function_call_arguments.delta":
		i.object(o, "event", allow("item_id", "output_index", "delta", "obfuscation"))
	case "response.function_call_arguments.done":
		i.object(o, "event", allow("item_id", "output_index", "arguments"))
	case "response.refusal.delta", "response.refusal.done", "response.output_text.annotation.added", "response.output_text.annotation.delta", "response.output_text.annotation.done", "response.reasoning_summary_part.added", "response.reasoning_summary_part.done", "response.reasoning_summary_text.delta", "response.reasoning_summary_text.done":
		// Recognized but unsupported event families still get envelope drift telemetry.
		i.object(o, "event", allow("item_id", "output_index", "content_index", "delta", "text", "annotation", "part"))
	case "error":
		i.object(o, "event", allow("error", "code", "message", "param"))
		if er, ok := obj(o["error"]); ok {
			i.object(er, "event.error", fields("type", "code", "message", "param"))
		}
	}
}

func inspectResponseEventItem(i wireInspector, item map[string]any, path string) {
	switch item["type"] {
	case "message":
		i.object(item, path, fields("id", "type", "status", "role", "content", "phase"))
		inspectResponsesParts(i, item["content"], path+".content")
	case "function_call":
		i.object(item, path, fields("id", "type", "status", "call_id", "name", "arguments", "phase"))
	}
}

func inspectResponseEventPart(i wireInspector, part map[string]any, path string) {
	if part["type"] == "output_text" {
		i.object(part, path, fields("type", "text", "annotations", "logprobs", "phase"))
	}
}

func inspectChatCompletionsEvent(e SSEEvent, options ConversionOptions) {
	if string(e.Data) == "[DONE]" {
		return
	}
	o, err := decodeObject(e.Data)
	if err != nil {
		return
	}
	i := wireInspector{options}
	if er, ok := obj(o["error"]); ok {
		i.object(o, "event", fields("error"))
		i.object(er, "event.error", fields("message", "type", "code", "param"))
		return
	}
	if o["object"] != "chat.completion.chunk" {
		return
	}
	i.object(o, "event", chatCompletionsResponseFields)
	inspectChatCompletionsChoices(i, o["choices"], "event.choices", true)
	inspectUsage(i, o["usage"], "event.usage", chatCompletionsUsageFields, fields("cached_tokens", "audio_tokens"), fields("reasoning_tokens", "audio_tokens", "accepted_prediction_tokens", "rejected_prediction_tokens"))
}
