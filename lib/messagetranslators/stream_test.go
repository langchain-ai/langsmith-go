package messagetranslators

import (
	"encoding/json"
	"errors"
	"testing"
)

func se(t *testing.T, typ, stringData string) SSEEvent {
	t.Helper()
	if !json.Valid([]byte(stringData)) {
		t.Fatalf("invalid fixture %s", stringData)
	}
	return SSEEvent{Event: typ, Data: []byte(stringData)}
}
func collect(t *testing.T, convert func(SSEEvent) ([]SSEEvent, error), inputs ...SSEEvent) []SSEEvent {
	t.Helper()
	var all []SSEEvent
	for _, in := range inputs {
		out, err := convert(in)
		if err != nil {
			t.Fatalf("convert %s: %v", in.Event, err)
		}
		all = append(all, out...)
	}
	return all
}
func names(es []SSEEvent) []string {
	r := make([]string, len(es))
	for i, e := range es {
		r[i] = e.Event
	}
	return r
}
func equalNames(t *testing.T, got []SSEEvent, want ...string) {
	t.Helper()
	n := names(got)
	if len(n) != len(want) {
		t.Fatalf("events %#v, want %#v", n, want)
	}
	for i := range n {
		if n[i] != want[i] {
			t.Fatalf("events %#v, want %#v", n, want)
		}
	}
}

func TestResponsesToAnthropicStreamTextToolAndUsage(t *testing.T) {
	c := NewResponsesToAnthropicStream("override")
	out := collect(t, c.Convert,
		se(t, "response.created", `{"type":"response.created","response":{"id":"resp_1","model":"r","status":"in_progress","usage":{"input_tokens":8,"output_tokens":0}}}`),
		se(t, "response.output_item.added", `{"type":"response.output_item.added","output_index":0,"item":{"type":"message","id":"m1"}}`),
		se(t, "response.content_part.added", `{"type":"response.content_part.added","output_index":0,"content_index":0,"part":{"type":"output_text","text":""}}`),
		se(t, "response.output_text.delta", `{"type":"response.output_text.delta","output_index":0,"content_index":0,"delta":"hel"}`),
		se(t, "response.output_text.delta", `{"type":"response.output_text.delta","output_index":0,"content_index":0,"delta":"lo"}`),
		se(t, "response.content_part.done", `{"type":"response.content_part.done","output_index":0,"content_index":0}`),
		se(t, "response.output_item.added", `{"type":"response.output_item.added","output_index":1,"item":{"type":"function_call","id":"fc","call_id":"call_1","name":"weather","arguments":""}}`),
		se(t, "response.function_call_arguments.delta", `{"type":"response.function_call_arguments.delta","output_index":1,"delta":"{\"city\":"}`),
		se(t, "response.function_call_arguments.delta", `{"type":"response.function_call_arguments.delta","output_index":1,"delta":"\"Paris\"}"}`),
		se(t, "response.output_item.done", `{"type":"response.output_item.done","output_index":1}`),
		se(t, "response.completed", `{"type":"response.completed","response":{"id":"resp_1","status":"completed","usage":{"input_tokens":8,"output_tokens":7}}}`),
	)
	equalNames(t, out, "message_start", "content_block_start", "content_block_delta", "content_block_delta", "content_block_stop", "content_block_start", "content_block_delta", "content_block_delta", "content_block_stop", "message_delta", "message_stop")
	start := object(t, out[0].Data)
	msg := mapv(t, start["message"])
	if msg["model"] != "override" || mapv(t, msg["usage"])["input_tokens"] != float64(8) {
		t.Fatalf("bad start: %#v", start)
	}
	toolStart := object(t, out[5].Data)
	if toolStart["index"] != float64(1) || mapv(t, toolStart["content_block"])["id"] != "call_1" {
		t.Fatalf("unstable tool block: %#v", toolStart)
	}
	delta := object(t, out[len(out)-2].Data)
	if mapv(t, delta["delta"])["stop_reason"] != "tool_use" || mapv(t, delta["usage"])["output_tokens"] != float64(7) {
		t.Fatalf("bad terminal delta: %#v", delta)
	}
	if err := c.Finish(); err != nil {
		t.Fatal(err)
	}
}

func TestResponsesToAnthropicIncompleteAndTruncation(t *testing.T) {
	c := NewResponsesToAnthropicStream("")
	collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"m"}}`), se(t, "response.incomplete", `{"response":{"id":"r","status":"incomplete","incomplete_details":{"reason":"max_output_tokens"},"usage":{"input_tokens":0,"output_tokens":9}}}`))
	if err := c.Finish(); err != nil {
		t.Fatal(err)
	}
	c2 := NewResponsesToAnthropicStream("")
	collect(t, c2.Convert, se(t, "response.created", `{"response":{"id":"r","model":"rm"}}`))
	if !errors.Is(c2.Finish(), ErrTruncatedStream) {
		t.Fatalf("got %v", c2.Finish())
	}
}

func TestAnthropicToResponsesStreamTextToolMaxTokens(t *testing.T) {
	c := NewAnthropicToResponsesStream("rm")
	out := collect(t, c.Convert,
		se(t, "message_start", `{"type":"message_start","message":{"id":"msg_a","model":"am","usage":{"input_tokens":11,"output_tokens":0}}}`),
		se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`),
		se(t, "content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"hi"}}`),
		se(t, "content_block_stop", `{"type":"content_block_stop","index":0}`),
		se(t, "content_block_start", `{"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"call_x","name":"f","input":{}}}`),
		se(t, "content_block_delta", `{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"{\"x\":"}}`),
		se(t, "content_block_delta", `{"type":"content_block_delta","index":1,"delta":{"type":"input_json_delta","partial_json":"1}"}}`),
		se(t, "content_block_stop", `{"type":"content_block_stop","index":1}`),
		se(t, "message_delta", `{"type":"message_delta","delta":{"stop_reason":"max_tokens","stop_sequence":null},"usage":{"output_tokens":6}}`),
		se(t, "message_stop", `{"type":"message_stop"}`),
	)
	equalNames(t, out, "response.created", "response.in_progress", "response.output_item.added", "response.content_part.added", "response.output_text.delta", "response.output_text.done", "response.content_part.done", "response.output_item.done", "response.output_item.added", "response.function_call_arguments.delta", "response.function_call_arguments.delta", "response.function_call_arguments.done", "response.output_item.done", "response.incomplete")
	terminal := object(t, out[len(out)-1].Data)
	r := mapv(t, terminal["response"])
	if r["status"] != "incomplete" || r["model"] != "rm" || mapv(t, r["usage"])["output_tokens"] != float64(6) {
		t.Fatalf("bad terminal: %#v", r)
	}
	outputs := list(t, r["output"])
	if len(outputs) != 2 || mapv(t, outputs[1])["call_id"] != "call_x" || mapv(t, outputs[1])["arguments"] != `{"x":1}` {
		t.Fatalf("bad outputs: %#v", outputs)
	}
	if err := c.Finish(); err != nil {
		t.Fatal(err)
	}
}

func TestStreamMalformedArgumentsErrorsAtFinalBlock(t *testing.T) {
	c := NewAnthropicToResponsesStream("")
	collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"am"}}`), se(t, "content_block_start", `{"index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{}}}`), se(t, "content_block_delta", `{"index":0,"delta":{"type":"input_json_delta","partial_json":"{"}}`))
	_, err := c.Convert(se(t, "content_block_stop", `{"index":0}`))
	if !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("got %v", err)
	}

	r := NewResponsesToAnthropicStream("")
	collect(t, r.Convert, se(t, "response.created", `{"response":{"id":"r","model":"rm"}}`), se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc_source","call_id":"c","name":"f"}}`), se(t, "response.function_call_arguments.delta", `{"output_index":0,"delta":"{"}`))
	_, err = r.Convert(se(t, "response.function_call_arguments.done", `{"output_index":0,"arguments":"{"}`))
	if !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("responses malformed final arguments: %v", err)
	}
}

func TestAnthropicResponsesEventSchemaAndUsage(t *testing.T) {
	c := NewAnthropicToResponsesStream("")
	out := collect(t, c.Convert,
		se(t, "message_start", `{"type":"message_start","message":{"id":"msg_source","model":"claude","usage":{"input_tokens":5,"cache_read_input_tokens":3,"cache_creation_input_tokens":2,"output_tokens":0}}}`),
		se(t, "message_delta", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":4}}`),
		se(t, "message_stop", `{"type":"message_stop"}`),
	)
	equalNames(t, out, "response.created", "response.in_progress", "response.completed")
	for i, e := range out {
		o := object(t, e.Data)
		if o["sequence_number"] != float64(i) {
			t.Fatalf("event %d sequence: %#v", i, o)
		}
	}
	created := mapv(t, object(t, out[0].Data)["response"])
	progress := mapv(t, object(t, out[1].Data)["response"])
	if created["usage"] != nil || progress["usage"] != nil {
		t.Fatalf("in-progress usage must be nil: %#v %#v", created, progress)
	}
	if created["id"] == "msg_source" {
		t.Fatalf("copied Anthropic message ID as response ID: %#v", created["id"])
	}
	terminal := mapv(t, object(t, out[2].Data)["response"])
	u := mapv(t, terminal["usage"])
	if u["input_tokens"] != float64(10) || u["output_tokens"] != float64(4) || u["total_tokens"] != float64(14) {
		t.Fatalf("terminal usage: %#v", u)
	}
}

func TestResponsesStreamTerminalOnlyInputUsageLimitation(t *testing.T) {
	c := NewResponsesToAnthropicStream("")
	out := collect(t, c.Convert,
		se(t, "response.created", `{"type":"response.created","response":{"id":"resp_x","model":"gpt","status":"in_progress"}}`),
		se(t, "response.completed", `{"type":"response.completed","response":{"id":"resp_x","status":"completed","usage":{"input_tokens":10,"output_tokens":2,"input_tokens_details":{"cached_tokens":3}}}}`),
	)
	equalNames(t, out, "message_start", "message_delta", "message_stop")
	startMessage := mapv(t, object(t, out[0].Data)["message"])
	if startMessage["id"] == "resp_x" {
		t.Fatalf("copied Responses ID as Anthropic message ID: %#v", startMessage["id"])
	}
	startUsage := mapv(t, startMessage["usage"])
	if startUsage["input_tokens"] != float64(0) {
		t.Fatalf("late terminal input must not be fabricated at start: %#v", startUsage)
	}
	deltaUsage := mapv(t, object(t, out[1].Data)["usage"])
	if len(deltaUsage) != 1 || deltaUsage["output_tokens"] != float64(2) {
		t.Fatalf("Anthropic terminal delta may contain output usage only: %#v", deltaUsage)
	}
}

func TestStreamSchemaErrorsFailedExtractionAndUnknownEvents(t *testing.T) {
	mismatch := NewResponsesToAnthropicStream("m")
	_, err := mismatch.Convert(se(t, "response.created", `{"type":"response.in_progress","response":{}}`))
	if !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("header mismatch: %v", err)
	}

	failed := NewResponsesToAnthropicStream("m")
	out := collect(t, failed.Convert, se(t, "response.failed", `{"type":"response.failed","response":{"id":"r","status":"failed","error":{"code":"server_error","message":"boom"}}}`))
	equalNames(t, out, "error")
	aerr := mapv(t, object(t, out[0].Data)["error"])
	if aerr["message"] != "boom" || aerr["type"] != "server_error" {
		t.Fatalf("failed error extraction: %#v", aerr)
	}

	anthropicErr := NewAnthropicToResponsesStream("")
	out = collect(t, anthropicErr.Convert, se(t, "error", `{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`))
	equalNames(t, out, "error")
	oe := object(t, out[0].Data)
	if oe["code"] != "overloaded_error" || oe["message"] != "busy" || oe["param"] != nil {
		t.Fatalf("invalid OpenAI error event: %#v", oe)
	}
	if _, nested := oe["error"]; nested {
		t.Fatalf("OpenAI error must be top-level: %#v", oe)
	}

	unknown := NewAnthropicToResponsesStream("")
	_, err = unknown.Convert(se(t, "mystery", `{"type":"mystery"}`))
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("unknown event: %v", err)
	}
}

func TestStreamIndexDeltaStopAndTerminalValidation(t *testing.T) {
	a := NewAnthropicToResponsesStream("gpt")
	collect(t, a.Convert, se(t, "message_start", `{"type":"message_start","message":{"id":"m","model":"claude","usage":{"input_tokens":0,"output_tokens":0}}}`))
	_, err := a.Convert(se(t, "content_block_start", `{"type":"content_block_start","index":-1,"content_block":{"type":"text","text":""}}`))
	if !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("negative index: %v", err)
	}

	a = NewAnthropicToResponsesStream("gpt")
	collect(t, a.Convert,
		se(t, "message_start", `{"type":"message_start","message":{"id":"m","model":"claude","usage":{"input_tokens":0,"output_tokens":0}}}`),
		se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`),
	)
	_, err = a.Convert(se(t, "content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","text":"x"}}`))
	if !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("wrong delta discriminator: %v", err)
	}
	collect(t, a.Convert, se(t, "content_block_stop", `{"type":"content_block_stop","index":0}`))
	_, err = a.Convert(se(t, "content_block_stop", `{"type":"content_block_stop","index":0}`))
	if !errors.Is(err, ErrInvalidSequence) {
		t.Fatalf("duplicate stop: %v", err)
	}

	noDelta := NewAnthropicToResponsesStream("gpt")
	collect(t, noDelta.Convert, se(t, "message_start", `{"type":"message_start","message":{"id":"m","model":"claude","usage":{"input_tokens":0,"output_tokens":0}}}`))
	_, err = noDelta.Convert(se(t, "message_stop", `{"type":"message_stop"}`))
	if !errors.Is(err, ErrInvalidSequence) {
		t.Fatalf("missing message_delta: %v", err)
	}

	badStop := NewAnthropicToResponsesStream("gpt")
	collect(t, badStop.Convert, se(t, "message_start", `{"type":"message_start","message":{"id":"m","model":"claude","usage":{"input_tokens":0,"output_tokens":0}}}`))
	_, err = badStop.Convert(se(t, "message_delta", `{"type":"message_delta","delta":{"stop_reason":"pause_turn"},"usage":{"output_tokens":0}}`))
	if !errors.Is(err, ErrUnsupported) {
		t.Fatalf("pause_turn: %v", err)
	}
}

func TestStreamErrorsAndInvalidSequence(t *testing.T) {
	r := NewResponsesToAnthropicStream("")
	out := collect(t, r.Convert, se(t, "error", `{"type":"error","error":{"message":"boom"}}`))
	equalNames(t, out, "error")
	if err := r.Finish(); err != nil {
		t.Fatal(err)
	}
	a := NewAnthropicToResponsesStream("")
	out = collect(t, a.Convert, se(t, "error", `{"type":"error","error":{"message":"boom"}}`))
	equalNames(t, out, "error")
	if err := a.Finish(); err != nil {
		t.Fatal(err)
	}
	r2 := NewResponsesToAnthropicStream("")
	_, err := r2.Convert(se(t, "response.output_text.delta", `{"output_index":0,"content_index":0,"delta":"x"}`))
	if !errors.Is(err, ErrInvalidSequence) {
		t.Fatalf("got %v", err)
	}
	a2 := NewAnthropicToResponsesStream("")
	collect(t, a2.Convert, se(t, "message_start", `{"message":{"id":"m","model":"am"}}`))
	if !errors.Is(a2.Finish(), ErrTruncatedStream) {
		t.Fatalf("got %v", a2.Finish())
	}
}
