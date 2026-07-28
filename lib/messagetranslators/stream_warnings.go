package messagetranslators

// Stream event allowlists. These are consulted once per SSE event, so they are
// package-level vars rather than per-event fields() calls.
var (
	anthropicEventMessageStartFields = fields("type", "message")
	anthropicEventBlockStartFields   = fields("type", "index", "content_block")
	anthropicEventBlockDeltaFields   = fields("type", "index", "delta")
	anthropicEventBlockStopFields    = fields("type", "index")
	anthropicEventMessageDeltaFields = fields("type", "delta", "usage")
	anthropicEventBareFields         = fields("type")
	anthropicEventErrorFields        = fields("type", "error")
	anthropicEventErrorBodyFields    = fields("type", "code", "message")
	anthropicEventStopDeltaFields    = fields("stop_reason", "stop_sequence", "stop_details")
	anthropicTextDeltaFields         = fields("type", "text", "citations")
	anthropicJSONDeltaFields         = fields("type", "partial_json")
	anthropicCitationDeltaFields     = fields("type", "citation")
)

func inspectDecodedAnthropicEvent(o map[string]any, typ string, cfg config) {
	i := wireInspector{cfg}
	switch typ {
	case "message_start":
		i.object(o, "event", anthropicEventMessageStartFields)
		if m, ok := obj(o["message"]); ok {
			inspectAnthropicObject(m, true, cfg, "event.message")
		}
	case "content_block_start":
		i.object(o, "event", anthropicEventBlockStartFields)
		if b, ok := obj(o["content_block"]); ok {
			inspectAnthropicContentBlock(i, b, "event.content_block")
		}
	case "content_block_delta":
		i.object(o, "event", anthropicEventBlockDeltaFields)
		if d, ok := obj(o["delta"]); ok {
			switch d["type"] {
			case "text_delta":
				i.object(d, "event.delta", anthropicTextDeltaFields)
			case "input_json_delta":
				i.object(d, "event.delta", anthropicJSONDeltaFields)
			case "citations_delta", "citation_delta":
				i.object(d, "event.delta", anthropicCitationDeltaFields)
			}
		}
	case "content_block_stop":
		i.object(o, "event", anthropicEventBlockStopFields)
	case "message_delta":
		i.object(o, "event", anthropicEventMessageDeltaFields)
		if d, ok := obj(o["delta"]); ok {
			i.object(d, "event.delta", anthropicEventStopDeltaFields)
			inspectAnthropicStopDetails(i, d, "event.delta")
		}
		inspectAnthropicUsage(i, o["usage"], "event.usage")
	case "message_stop", "ping":
		i.object(o, "event", anthropicEventBareFields)
	case "error":
		i.object(o, "event", anthropicEventErrorFields)
		if er, ok := obj(o["error"]); ok {
			i.object(er, "event.error", anthropicEventErrorBodyFields)
		}
	}
}

var (
	responsesEventResponseFields    = fields("type", "sequence_number", "response")
	responsesEventItemFields        = fields("type", "sequence_number", "output_index", "item")
	responsesEventPartFields        = fields("type", "sequence_number", "item_id", "output_index", "content_index", "part")
	responsesEventTextDeltaFields   = fields("type", "sequence_number", "item_id", "output_index", "content_index", "delta", "logprobs", "obfuscation")
	responsesEventTextDoneFields    = fields("type", "sequence_number", "item_id", "output_index", "content_index", "text", "annotations", "logprobs")
	responsesEventArgsDeltaFields   = fields("type", "sequence_number", "item_id", "output_index", "delta", "obfuscation")
	responsesEventArgsDoneFields    = fields("type", "sequence_number", "item_id", "output_index", "arguments")
	responsesEventUnsupportedFields = fields("type", "sequence_number", "item_id", "output_index", "content_index", "delta", "text", "annotation", "part")
	responsesEventErrorFields       = fields("type", "sequence_number", "error", "code", "message", "param")
	responsesEventErrorBodyFields   = fields("type", "code", "message", "param")
)

func inspectResponsesEvent(o map[string]any, typ string, cfg config) {
	i := wireInspector{cfg}
	switch typ {
	case "response.created", "response.in_progress", "response.completed", "response.incomplete", "response.failed":
		i.object(o, "event", responsesEventResponseFields)
		if r, ok := obj(o["response"]); ok {
			inspectResponsesObject(r, true, cfg, "event.response")
		}
	case "response.output_item.added", "response.output_item.done":
		i.object(o, "event", responsesEventItemFields)
		if item, ok := obj(o["item"]); ok {
			inspectResponseEventItem(i, item, "event.item")
		}
	case "response.content_part.added", "response.content_part.done":
		i.object(o, "event", responsesEventPartFields)
		if part, ok := obj(o["part"]); ok {
			inspectResponseEventPart(i, part, "event.part")
		}
	case "response.output_text.delta":
		i.object(o, "event", responsesEventTextDeltaFields)
	case "response.output_text.done":
		i.object(o, "event", responsesEventTextDoneFields)
	case "response.function_call_arguments.delta":
		i.object(o, "event", responsesEventArgsDeltaFields)
	case "response.function_call_arguments.done":
		i.object(o, "event", responsesEventArgsDoneFields)
	case "response.refusal.delta", "response.refusal.done", "response.output_text.annotation.added", "response.output_text.annotation.delta", "response.output_text.annotation.done", "response.reasoning_summary_part.added", "response.reasoning_summary_part.done", "response.reasoning_summary_text.delta", "response.reasoning_summary_text.done":
		// Recognized but unsupported event families still get envelope drift telemetry.
		i.object(o, "event", responsesEventUnsupportedFields)
	case "error":
		i.object(o, "event", responsesEventErrorFields)
		if er, ok := obj(o["error"]); ok {
			i.object(er, "event.error", responsesEventErrorBodyFields)
		}
	}
}

func inspectResponseEventItem(i wireInspector, item map[string]any, path string) {
	switch item["type"] {
	case "message":
		i.object(item, path, responsesMessageItemFields)
		inspectResponsesParts(i, item["content"], path+".content")
	case "function_call":
		i.object(item, path, responsesFunctionCallItemFields)
	}
}

func inspectResponseEventPart(i wireInspector, part map[string]any, path string) {
	if part["type"] == "output_text" {
		i.object(part, path, responsesOutputTextPartFields)
	}
}

var (
	chatCompletionsEventErrorFields     = fields("error")
	chatCompletionsEventErrorBodyFields = fields("message", "type", "code", "param")
)

func inspectChatCompletionsEvent(o map[string]any, cfg config) {
	i := wireInspector{cfg}
	if er, ok := obj(o["error"]); ok {
		i.object(o, "event", chatCompletionsEventErrorFields)
		i.object(er, "event.error", chatCompletionsEventErrorBodyFields)
		return
	}
	if o["object"] != "chat.completion.chunk" {
		return
	}
	i.object(o, "event", chatCompletionsResponseFields)
	inspectChatCompletionsChoices(i, o["choices"], "event.choices", true)
	inspectUsage(i, o["usage"], "event.usage", chatCompletionsUsageFields, chatCompletionsInputDetailFields, chatCompletionsOutputDetailFields)
}
