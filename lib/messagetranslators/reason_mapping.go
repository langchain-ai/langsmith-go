package messagetranslators

func chatCompletionsFinishReasonToAnthropic(finish string, hasTools bool, path string) (string, error) {
	switch finish {
	case "stop":
		if hasTools {
			return "", at(path, ErrInvalidSequence)
		}
		return "end_turn", nil
	case "tool_calls":
		if !hasTools {
			return "", at(path, ErrInvalidSequence)
		}
		return "tool_use", nil
	case "length":
		return "max_tokens", nil
	case "content_filter", "function_call":
		return "", at(path, ErrUnsupported)
	default:
		return "", at(path, ErrInvalidWireData)
	}
}

func anthropicStopReasonToChatCompletions(stop string, hasTools bool, path string) (string, error) {
	switch stop {
	case "end_turn", "stop_sequence":
		if hasTools {
			return "", at(path, ErrInvalidSequence)
		}
		return "stop", nil
	case "tool_use":
		if !hasTools {
			return "", at(path, ErrInvalidSequence)
		}
		return "tool_calls", nil
	case "max_tokens":
		return "length", nil
	case "pause_turn", "refusal":
		return "", at(path, ErrUnsupported)
	default:
		return "", at(path, ErrInvalidWireData)
	}
}
