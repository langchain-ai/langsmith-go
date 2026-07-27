package messagetranslators

import (
	"fmt"
	"math"
)

// Usage reports token accounting observed on the source payload, expressed in
// Anthropic terms: InputTokens excludes cache reads, which are counted
// separately in CacheReadInputTokens.
//
// Usage is reported even when the destination wire format cannot carry it. That
// matters most for Responses-to-Anthropic and ChatCompletions-to-Anthropic
// streams, where input token counts arrive in a terminal event, long after the
// Anthropic message_start that would have had to carry them. Without this, a
// gateway that bills or meters from the translated stream would have to reparse
// the source stream to recover counts the converter already saw.
type Usage struct {
	InputTokens          int64
	OutputTokens         int64
	CacheReadInputTokens int64
}

// UsageHandler receives final token accounting once, when a converter observes
// a terminal event. It is not called for a truncated stream.
type UsageHandler func(Usage)

func token(v any, path string) (int64, error) { return integer(v, path, false) }

// openAIUsageToAnthropic accounts for cached_tokens being a subset of input_tokens.
func openAIUsageToAnthropic(u map[string]any, path string, cfg config) (map[string]any, error) {
	in, err := token(u["input_tokens"], path+".input_tokens")
	if err != nil {
		return nil, err
	}
	out, err := token(u["output_tokens"], path+".output_tokens")
	if err != nil {
		return nil, err
	}
	var cached int64
	if details, exists := u["input_tokens_details"]; exists {
		d, ok := obj(details)
		if !ok {
			return nil, at(path+".input_tokens_details", ErrInvalidWireData)
		}
		if v, exists := d["cached_tokens"]; exists {
			cached, err = token(v, path+".input_tokens_details.cached_tokens")
			if err != nil {
				return nil, err
			}
		}
	}
	if cached > in {
		return nil, at(path+".input_tokens_details.cached_tokens", fmt.Errorf("%w: cached_tokens exceeds input_tokens", ErrInvalidWireData))
	}
	// Anthropic has no reasoning-token usage category. Validate the wire value
	// before intentionally dropping it (Python/JS parity), rather than accepting
	// malformed details merely because the destination cannot represent them.
	if details, exists := u["output_tokens_details"]; exists {
		d, ok := obj(details)
		if !ok {
			return nil, at(path+".output_tokens_details", ErrInvalidWireData)
		}
		if v, exists := d["reasoning_tokens"]; exists {
			reasoning, err := token(v, path+".output_tokens_details.reasoning_tokens")
			if err != nil {
				return nil, err
			}
			if reasoning != 0 {
				cfg.lossy(path+".output_tokens_details.reasoning_tokens", "reasoning_tokens",
					"reasoning token count has no Anthropic usage category and was dropped")
			}
		}
	}
	if v, ok := u["total_tokens"]; ok {
		total, err := token(v, path+".total_tokens")
		if err != nil {
			return nil, err
		}
		if in > math.MaxInt64-out || total != in+out {
			return nil, at(path+".total_tokens", fmt.Errorf("%w: total_tokens must equal input_tokens + output_tokens", ErrInvalidWireData))
		}
	}
	return map[string]any{"input_tokens": in - cached, "output_tokens": out, "cache_read_input_tokens": cached}, nil
}

// anthropicUsageToOpenAI sums Anthropic's disjoint input token categories.
func anthropicUsageToOpenAI(u map[string]any, path string, cfg config) (map[string]any, error) {
	in, err := token(u["input_tokens"], path+".input_tokens")
	if err != nil {
		return nil, err
	}
	out, err := token(u["output_tokens"], path+".output_tokens")
	if err != nil {
		return nil, err
	}
	var read, create int64
	if v, ok := u["cache_read_input_tokens"]; ok {
		read, err = token(v, path+".cache_read_input_tokens")
		if err != nil {
			return nil, err
		}
	}
	if v, ok := u["cache_creation_input_tokens"]; ok {
		create, err = token(v, path+".cache_creation_input_tokens")
		if err != nil {
			return nil, err
		}
		if create != 0 {
			// OpenAI reports only cached reads, so cache writes can be summed into
			// the input total but their separate, differently-billed identity is
			// lost.
			cfg.lossy(path+".cache_creation_input_tokens", "cache_creation_input_tokens",
				"cache creation tokens have no OpenAI usage category and were folded into input_tokens")
		}
	}
	const maxInt64 = math.MaxInt64
	if read > maxInt64-in || create > maxInt64-in-read || out > maxInt64-in-read-create {
		return nil, at(path, fmt.Errorf("%w: token total overflows int64", ErrInvalidWireData))
	}
	totalIn := in + read + create
	return map[string]any{"input_tokens": totalIn, "output_tokens": out, "total_tokens": totalIn + out, "input_tokens_details": map[string]any{"cached_tokens": read}, "output_tokens_details": map[string]any{"reasoning_tokens": int64(0)}}, nil
}

// chatCompletionsUsageToAnthropic converts Chat Completions usage. cached_tokens is a
// subset of prompt_tokens, as it is in the Responses API.
func chatCompletionsUsageToAnthropic(u map[string]any, path string, cfg config) (map[string]any, error) {
	in, err := token(u["prompt_tokens"], path+".prompt_tokens")
	if err != nil {
		return nil, err
	}
	out, err := token(u["completion_tokens"], path+".completion_tokens")
	if err != nil {
		return nil, err
	}
	var cached int64
	if v, ok := u["prompt_tokens_details"]; ok {
		d, ok := obj(v)
		if !ok {
			return nil, at(path+".prompt_tokens_details", ErrInvalidWireData)
		}
		if v, ok := d["cached_tokens"]; ok {
			cached, err = token(v, path+".prompt_tokens_details.cached_tokens")
			if err != nil {
				return nil, err
			}
		}
		if v, ok := d["audio_tokens"]; ok {
			n, e := token(v, path+".prompt_tokens_details.audio_tokens")
			if e != nil {
				return nil, e
			}
			if n != 0 {
				return nil, at(path+".prompt_tokens_details.audio_tokens", ErrUnsupported)
			}
		}
	}
	if cached > in {
		return nil, at(path+".prompt_tokens_details.cached_tokens", fmt.Errorf("%w: cached_tokens exceeds prompt_tokens", ErrInvalidWireData))
	}
	if v, ok := u["completion_tokens_details"]; ok {
		d, ok := obj(v)
		if !ok {
			return nil, at(path+".completion_tokens_details", ErrInvalidWireData)
		}
		// Keep parity with Responses usage: Anthropic has no reasoning-token
		// category, so validate reasoning_tokens and intentionally drop it. Audio
		// and speculative-prediction counts describe unsupported output modes and
		// may only be present as zero-valued SDK defaults.
		if v, ok := d["reasoning_tokens"]; ok {
			reasoning, e := token(v, path+".completion_tokens_details.reasoning_tokens")
			if e != nil {
				return nil, e
			}
			if reasoning != 0 {
				cfg.lossy(path+".completion_tokens_details.reasoning_tokens", "reasoning_tokens",
					"reasoning token count has no Anthropic usage category and was dropped")
			}
		}
		for _, key := range []string{"audio_tokens", "accepted_prediction_tokens", "rejected_prediction_tokens"} {
			if v, ok := d[key]; ok {
				n, e := token(v, path+".completion_tokens_details."+key)
				if e != nil {
					return nil, e
				}
				if n != 0 {
					return nil, at(path+".completion_tokens_details."+key, ErrUnsupported)
				}
			}
		}
	}
	if v, ok := u["total_tokens"]; ok {
		total, e := token(v, path+".total_tokens")
		if e != nil {
			return nil, e
		}
		if in > math.MaxInt64-out || total != in+out {
			return nil, at(path+".total_tokens", fmt.Errorf("%w: total_tokens must equal prompt_tokens + completion_tokens", ErrInvalidWireData))
		}
	}
	return map[string]any{"input_tokens": in - cached, "output_tokens": out, "cache_read_input_tokens": cached}, nil
}

func anthropicUsageToChatCompletions(u map[string]any, path string, cfg config) (map[string]any, error) {
	ru, err := anthropicUsageToOpenAI(u, path, cfg)
	if err != nil {
		return nil, err
	}
	return map[string]any{
		"prompt_tokens":     ru["input_tokens"],
		"completion_tokens": ru["output_tokens"],
		"total_tokens":      ru["total_tokens"],
		"prompt_tokens_details": map[string]any{
			"cached_tokens": ru["input_tokens_details"].(map[string]any)["cached_tokens"],
		},
	}, nil
}

// anthropicUsage and openAIUsage lift an already-validated usage map into the
// reporting struct, so the usage callback never re-derives token math.
func anthropicUsage(u map[string]any) Usage {
	return Usage{
		InputTokens:          u["input_tokens"].(int64),
		OutputTokens:         u["output_tokens"].(int64),
		CacheReadInputTokens: u["cache_read_input_tokens"].(int64),
	}
}

func openAIUsage(u map[string]any) Usage {
	cached := u["input_tokens_details"].(map[string]any)["cached_tokens"].(int64)
	return Usage{
		InputTokens:          u["input_tokens"].(int64) - cached,
		OutputTokens:         u["output_tokens"].(int64),
		CacheReadInputTokens: cached,
	}
}
