package messagetranslators

import (
	"encoding/json"
	"fmt"
)

// fieldPolicy describes one source request field that has no destination
// representation.
//
// Many such fields are sent by SDKs even when the caller never set them,
// carrying the value the provider documents as the default. That value is
// semantically identical to omitting the field, so rejecting it would fail a
// request that asked for nothing unusual. accepts recognizes exactly that
// value; anything else is a real feature request that must not be silently
// dropped, and is rejected with ErrUnsupported.
type fieldPolicy struct {
	key     string
	accepts func(v any) bool // nil rejects any present, non-null value
	why     string           // optional detail appended to the error
}

// rejectUnsupportedFields applies policy groups in declaration order so that
// the field reported for a request violating several policies is deterministic.
//
// An explicit JSON null is treated as absent throughout: SDKs serialize unset
// optional fields that way, and null cannot express a feature request.
func rejectUnsupportedFields(o map[string]any, path string, groups ...[]fieldPolicy) error {
	for _, group := range groups {
		for _, policy := range group {
			v, present := o[policy.key]
			if !present || v == nil {
				continue
			}
			if policy.accepts != nil && policy.accepts(v) {
				continue
			}
			if policy.why != "" {
				return at(path+"."+policy.key, fmt.Errorf("%w: %s", ErrUnsupported, policy.why))
			}
			return at(path+"."+policy.key, ErrUnsupported)
		}
	}
	return nil
}

// rejectNonEmptyArray implements the v0 policy for metadata such as
// citations and annotations: an explicitly empty SDK-default array is harmless,
// but populated semantic metadata must never be silently discarded.
func rejectNonEmptyArray(o map[string]any, key, path string) error {
	v, exists := o[key]
	if !exists {
		return nil
	}
	a, ok := arr(v)
	if !ok {
		return at(path+"."+key, ErrInvalidWireData)
	}
	if len(a) != 0 {
		return at(path+"."+key, ErrUnsupported)
	}
	return nil
}

// rejectPresent rejects keys carrying a value. An explicit null is treated as
// absent, because SDKs serialize unset optional fields that way: OpenAI
// responses routinely include "function_call": null on messages that have none.
func rejectPresent(o map[string]any, path string, keys ...string) error {
	for _, key := range keys {
		if v, ok := o[key]; ok && v != nil {
			return at(path+"."+key, ErrUnsupported)
		}
	}
	return nil
}

func boolIs(want bool) func(any) bool {
	return func(v any) bool { b, ok := v.(bool); return ok && b == want }
}

func stringIn(want ...string) func(any) bool {
	return func(v any) bool {
		s, ok := str(v)
		if !ok {
			return false
		}
		for _, w := range want {
			if s == w {
				return true
			}
		}
		return false
	}
}

func numberIs(want int64) func(any) bool {
	return func(v any) bool {
		n, ok := v.(json.Number)
		if !ok {
			return false
		}
		i, err := n.Int64()
		return err == nil && i == want
	}
}

func emptyArray(v any) bool {
	a, ok := arr(v)
	return ok && len(a) == 0
}

func emptyObject(v any) bool {
	o, ok := obj(v)
	return ok && len(o) == 0
}

func arrayIs(want ...string) func(any) bool {
	return func(v any) bool {
		a, ok := arr(v)
		if !ok || len(a) != len(want) {
			return false
		}
		for n, x := range a {
			if s, ok := str(x); !ok || s != want[n] {
				return false
			}
		}
		return true
	}
}

// plainTextFormat matches the Chat Completions response_format default and,
// nested one level deeper, the Responses text default. Both say "return
// unstructured text", which is the only behavior these translators produce.
func plainTextFormat(v any) bool {
	o, ok := obj(v)
	return ok && len(o) == 1 && o["type"] == "text"
}

func responsesTextDefault(v any) bool {
	o, ok := obj(v)
	if !ok {
		return false
	}
	if len(o) == 0 {
		return true
	}
	if len(o) != 1 {
		return false
	}
	format, ok := obj(o["format"])
	return ok && len(format) == 0 || plainTextFormat(o["format"])
}

// Anthropic Messages request fields with no equivalent in either OpenAI format.
var anthropicRequestUnsupported = []fieldPolicy{
	{key: "thinking", why: "extended thinking has no destination equivalent"},
	{key: "top_k"},
	{key: "service_tier", accepts: stringIn("auto")},
	{key: "output_config"},
	{key: "metadata", accepts: emptyObject},
	{key: "mcp_servers", accepts: emptyArray},
	{key: "context_management"},
	{key: "container"},
}

// Responses has no stop-sequence control; Chat Completions does, so this group
// applies only to the Responses destination.
var anthropicToResponsesUnsupported = []fieldPolicy{
	{key: "stop_sequences", accepts: emptyArray, why: "Responses has no stop-sequence control"},
}

var responsesRequestUnsupported = []fieldPolicy{
	{key: "previous_response_id", why: "server-side conversation continuation has no destination equivalent"},
	{key: "conversation", why: "server-side conversation state has no destination equivalent"},
	{key: "background", accepts: boolIs(false), why: "background responses have no destination equivalent"},
	{key: "reasoning"},
	{key: "include", accepts: emptyArray},
	{key: "truncation", accepts: stringIn("disabled")},
	{key: "max_tool_calls"},
	{key: "parallel_tool_calls", accepts: boolIs(true)},
	{key: "text", accepts: responsesTextDefault, why: "structured output has no destination equivalent"},
	{key: "phase"},
	{key: "prompt_cache_key"},
	{key: "prompt_cache_retention"},
	{key: "store", accepts: boolIs(true), why: "response persistence is controlled by the destination provider"},
	{key: "top_logprobs", accepts: numberIs(0)},
	{key: "service_tier", accepts: stringIn("auto", "default")},
	{key: "metadata", accepts: emptyObject},
	{key: "user"},
	{key: "safety_identifier"},
}

var chatCompletionsRequestUnsupported = []fieldPolicy{
	{key: "functions", why: "legacy functions are not supported; use tools"},
	{key: "function_call", why: "legacy function_call is not supported; use tool_choice"},
	{key: "response_format", accepts: plainTextFormat, why: "structured output has no destination equivalent"},
	{key: "parallel_tool_calls", accepts: boolIs(true)},
	{key: "logprobs", accepts: boolIs(false)},
	{key: "top_logprobs", accepts: numberIs(0)},
	{key: "prediction"},
	{key: "modalities", accepts: arrayIs("text")},
	{key: "audio"},
	{key: "service_tier", accepts: stringIn("auto", "default")},
	{key: "store", accepts: boolIs(false), why: "response persistence is controlled by the destination provider"},
	{key: "metadata", accepts: emptyObject},
	{key: "seed"},
	{key: "user"},
	{key: "reasoning_effort"},
	{key: "web_search_options"},
	{key: "frequency_penalty", accepts: numberIs(0)},
	{key: "presence_penalty", accepts: numberIs(0)},
	{key: "logit_bias", accepts: emptyObject},
	{key: "verbosity", accepts: stringIn("medium")},
	{key: "safety_identifier"},
}
