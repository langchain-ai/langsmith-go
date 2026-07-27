package messagetranslators

import (
	"errors"
	"testing"
)

func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("error = %v, want errors.Is(_, %v)", err, target)
	}
}

// These cases mirror the important Python/JS parity rule: SDK-default empty
// metadata is harmless, but semantic citation/annotation data is never dropped.
func TestParityCitationAndAnnotationPolicy(t *testing.T) {
	responsePrefix := `{"id":"r","model":"g","status":"completed","output":[{"type":"message","id":"m","role":"assistant","status":"completed","content":[{"type":"output_text","text":"x","annotations":`
	responseSuffix := `}]}],"usage":{"input_tokens":0,"output_tokens":1}}`
	if _, err := ResponsesResponseToAnthropic([]byte(responsePrefix+`[]`+responseSuffix), "c"); err != nil {
		t.Fatalf("empty annotations: %v", err)
	}
	_, err := ResponsesResponseToAnthropic([]byte(responsePrefix+`[{"type":"url_citation"}]`+responseSuffix), "c")
	requireErrorIs(t, err, ErrUnsupported)

	anthropic := func(citations string) []byte {
		return []byte(`{"id":"m","type":"message","role":"assistant","model":"c","content":[{"type":"text","text":"x","citations":` + citations + `}],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":1}}`)
	}
	if _, err := AnthropicResponseToResponses(anthropic(`[]`), "g"); err != nil {
		t.Fatalf("empty citations: %v", err)
	}
	_, err = AnthropicResponseToResponses(anthropic(`[{"type":"char_location"}]`), "g")
	requireErrorIs(t, err, ErrUnsupported)
}

func TestParityAnthropicTextToolTextAggregation(t *testing.T) {
	body := []byte(`{"id":"m","type":"message","role":"assistant","model":"c","content":[{"type":"text","text":"a"},{"type":"text","text":"b"},{"type":"tool_use","id":"call","name":"f","input":{"snowman":"☃"}},{"type":"text","text":"c"},{"type":"text","text":"d"}],"stop_reason":"tool_use","usage":{"input_tokens":1,"output_tokens":2}}`)
	got, err := AnthropicResponseToResponses(body, "g")
	if err != nil {
		t.Fatal(err)
	}
	outputs := list(t, object(t, got)["output"])
	if len(outputs) != 3 {
		t.Fatalf("outputs = %#v", outputs)
	}
	first := list(t, mapv(t, outputs[0])["content"])
	last := list(t, mapv(t, outputs[2])["content"])
	if len(first) != 2 || mapv(t, first[0])["text"] != "a" || mapv(t, first[1])["text"] != "b" || len(last) != 2 {
		t.Fatalf("contiguous text was not aggregated: %#v", outputs)
	}
	if mapv(t, outputs[1])["type"] != "function_call" || mapv(t, outputs[1])["arguments"] != `{"snowman":"☃"}` {
		t.Fatalf("tool ordering/unicode lost: %#v", outputs[1])
	}
}

func TestParityResponseRequiredFieldsAndReasoningUsageValidation(t *testing.T) {
	_, err := AnthropicResponseToResponses([]byte(`{"id":"m","model":"c","content":[],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":0}}`), "g")
	requireErrorIs(t, err, ErrInvalidWireData)
	_, err = ResponsesResponseToAnthropic([]byte(`{"id":"r","model":"g","status":"completed","usage":{"input_tokens":0,"output_tokens":0}}`), "c")
	requireErrorIs(t, err, ErrInvalidWireData)
	_, err = ResponsesResponseToAnthropic([]byte(`{"id":"r","model":"g","status":"completed","output":[],"usage":{"input_tokens":0,"output_tokens":0,"output_tokens_details":{"reasoning_tokens":"1"}}}`), "c")
	requireErrorIs(t, err, ErrInvalidWireData)
	if _, err = ResponsesResponseToAnthropic([]byte(`{"id":"r","model":"g","status":"completed","output":[],"usage":{"input_tokens":0,"output_tokens":0,"output_tokens_details":{"reasoning_tokens":7}}}`), "c"); err != nil {
		t.Fatalf("valid but intentionally unrepresentable reasoning count: %v", err)
	}
}

func TestParityIncompleteTruncatedArgumentsPolicy(t *testing.T) {
	base := func(status, args string) []byte {
		details := ""
		if status == "incomplete" {
			details = `,"incomplete_details":{"reason":"max_output_tokens"}`
		}
		return []byte(`{"id":"r","model":"g","status":"` + status + `"` + details + `,"output":[{"type":"function_call","id":"fc","call_id":"c","name":"f","arguments":` + args + `,"status":"` + status + `"}],"usage":{"input_tokens":0,"output_tokens":1}}`)
	}
	_, err := ResponsesResponseToAnthropic(base("incomplete", `"{"`), "c")
	requireErrorIs(t, err, ErrUnsupported)
	_, err = ResponsesResponseToAnthropic(base("completed", `"{"`), "c")
	requireErrorIs(t, err, ErrInvalidWireData)
}

func TestParityAnthropicInitialStreamTextAndCacheRead(t *testing.T) {
	c := NewAnthropicToResponsesStream("g")
	out := collect(t, c.Convert,
		se(t, "message_start", `{"message":{"id":"m","model":"c","usage":{"input_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4,"output_tokens":0}}}`),
		se(t, "content_block_start", `{"index":0,"content_block":{"type":"text","text":"initial"}}`),
		se(t, "content_block_delta", `{"index":0,"delta":{"type":"text_delta","text":"+delta"}}`),
		se(t, "content_block_stop", `{"index":0}`),
		se(t, "message_delta", `{"delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":2}}`),
		se(t, "message_stop", `{}`),
	)
	// No synthetic output_text.delta is emitted for the start payload itself.
	equalNames(t, out, "response.created", "response.in_progress", "response.output_item.added", "response.content_part.added", "response.output_text.delta", "response.output_text.done", "response.content_part.done", "response.output_item.done", "response.completed")
	added := mapv(t, object(t, out[3].Data)["part"])
	if added["text"] != "initial" {
		t.Fatalf("initial text lost: %#v", added)
	}
	terminal := mapv(t, object(t, out[len(out)-1].Data)["response"])
	u := mapv(t, terminal["usage"])
	if u["input_tokens"] != float64(9) || mapv(t, u["input_tokens_details"])["cached_tokens"] != float64(3) {
		t.Fatalf("cache categories lost: %#v", u)
	}
}

func TestParityResponsesStreamFSMAndFailureClosure(t *testing.T) {
	t.Run("content-on-tool", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("c")
		collect(t, c.Convert,
			se(t, "response.created", `{"response":{"id":"r","model":"g"}}`),
			se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc","call_id":"call","name":"f"}}`),
		)
		_, err := c.Convert(se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":""}}`))
		requireErrorIs(t, err, ErrInvalidSequence)
	})

	t.Run("item-done-with-open-content", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("c")
		collect(t, c.Convert,
			se(t, "response.created", `{"response":{"id":"r","model":"g"}}`),
			se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m"}}`),
			se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":""}}`),
		)
		_, err := c.Convert(se(t, "response.output_item.done", `{"output_index":0}`))
		requireErrorIs(t, err, ErrInvalidSequence)
	})

	t.Run("failure-closes-open-block", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("c")
		out := collect(t, c.Convert,
			se(t, "response.created", `{"response":{"id":"r","model":"g"}}`),
			se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m"}}`),
			se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":""}}`),
			se(t, "response.failed", `{"response":{"error":{"code":"server_error","message":"boom"}}}`),
		)
		equalNames(t, out, "message_start", "content_block_start", "content_block_stop", "error")
	})
}

func TestParityStreamUnsupportedMetadataAndTruncatedCall(t *testing.T) {
	c := NewResponsesToAnthropicStream("c")
	collect(t, c.Convert,
		se(t, "response.created", `{"response":{"id":"r","model":"g"}}`),
		se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m"}}`),
	)
	_, err := c.Convert(se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":"","annotations":[{"type":"url_citation"}]}}`))
	requireErrorIs(t, err, ErrUnsupported)

	truncated := NewResponsesToAnthropicStream("c")
	collect(t, truncated.Convert,
		se(t, "response.created", `{"response":{"id":"r","model":"g"}}`),
		se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc","call_id":"call","name":"f"}}`),
		se(t, "response.function_call_arguments.delta", `{"output_index":0,"delta":"{"}`),
	)
	_, err = truncated.Convert(se(t, "response.incomplete", `{"response":{"id":"r","status":"incomplete","usage":{"input_tokens":0,"output_tokens":1}}}`))
	requireErrorIs(t, err, ErrUnsupported)
}
