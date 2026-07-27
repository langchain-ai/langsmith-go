package messagetranslators

// Terminal-reason mapping between Chat Completions and Anthropic Messages.
//
// Where the source's advisory reason disagrees with its own content, the
// content wins. A response that carries tool calls ended in a tool call, no
// matter what finish_reason says, and deriving the destination reason from the
// content is both unambiguous and lossless. This matters in practice: OpenAI
// itself reports finish_reason "tool_calls", but several OpenAI-compatible
// servers report "stop" alongside a populated tool_calls array, and rejecting
// those would make the translators unusable in front of them.
//
// Reasons that describe something other than how the turn ended are still
// rejected, because those are real features rather than a disagreement.

func chatCompletionsFinishReasonToAnthropic(finish string, hasTools bool, path string) (string, error) {
	switch finish {
	case "stop", "tool_calls":
		if hasTools {
			return "tool_use", nil
		}
		return "end_turn", nil
	case "length":
		return "max_tokens", nil
	case "content_filter":
		// Anthropic's refusal stop reason is the closest equivalent, and mapping
		// it keeps Azure OpenAI's filtered completions translatable.
		return "refusal", nil
	case "function_call":
		return "", at(path, ErrUnsupported)
	default:
		return "", at(path, ErrInvalidWireData)
	}
}

func anthropicStopReasonToChatCompletions(stop string, hasTools bool, path string) (string, error) {
	switch stop {
	case "end_turn", "stop_sequence", "tool_use":
		if hasTools {
			return "tool_calls", nil
		}
		return "stop", nil
	case "max_tokens":
		return "length", nil
	case "refusal":
		return "content_filter", nil
	case "pause_turn":
		return "", at(path, ErrUnsupported)
	default:
		return "", at(path, ErrInvalidWireData)
	}
}
