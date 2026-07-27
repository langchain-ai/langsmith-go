package messagetranslators

import "testing"

func convertError(t *testing.T, convert func(SSEEvent) ([]SSEEvent, error), event SSEEvent, target error) {
	t.Helper()
	out, err := convert(event)
	if out != nil {
		t.Fatalf("error returned output: %#v", out)
	}
	requireErrorIs(t, err, target)
}

// Source parity: Python/JS streaming initial text, fragmented Unicode tools, ordering, and terminal equality.
func TestV0AnthropicToResponsesStreamPayloadMatrix(t *testing.T) {
	t.Run("initial-text-empty-and-nonempty", func(t *testing.T) {
		for _, tc := range []struct{ name, initial string }{{"empty", ""}, {"nonempty", "initial ☃"}} {
			t.Run(tc.name, func(t *testing.T) {
				initial := tc.initial
				c := NewAnthropicToResponsesStream("g")
				out := collect(t, c.Convert,
					se(t, "message_start", `{"message":{"id":"m","model":"a","usage":{"input_tokens":0,"output_tokens":0}}}`),
					se(t, "content_block_start", `{"index":0,"content_block":{"type":"text","text":"`+initial+`","cache_control":{"type":"ephemeral"}}}`),
					se(t, "content_block_stop", `{"index":0}`), se(t, "message_delta", `{"delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":0}}`), se(t, "message_stop", `{}`))
				part := mapv(t, object(t, out[3].Data)["part"])
				if part["text"] != initial {
					t.Fatalf("part=%#v", part)
				}
				terminal := mapv(t, object(t, out[len(out)-1].Data)["response"])
				doneItem := mapv(t, object(t, out[len(out)-2].Data)["item"])
				if mapv(t, list(t, mapv(t, list(t, terminal["output"])[0])["content"])[0])["text"] != initial || mapv(t, list(t, doneItem["content"])[0])["text"] != initial {
					t.Fatalf("terminal/done mismatch: %#v %#v", terminal, doneItem)
				}
			})
		}
	})

	t.Run("fragmented-unicode-tool-arguments", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		out := collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a","usage":{"input_tokens":0,"output_tokens":0}}}`), se(t, "content_block_start", `{"index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{}}}`), se(t, "content_block_delta", `{"index":0,"delta":{"type":"input_json_delta","partial_json":"{\"雪\":\""}}`), se(t, "content_block_delta", `{"index":0,"delta":{"type":"input_json_delta","partial_json":"☃\"}"}}`), se(t, "content_block_stop", `{"index":0}`), se(t, "message_delta", `{"delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":2}}`), se(t, "message_stop", `{}`))
		terminal := mapv(t, object(t, out[len(out)-1].Data)["response"])
		if mapv(t, list(t, terminal["output"])[0])["arguments"] != `{"雪":"☃"}` {
			t.Fatalf("terminal=%#v", terminal)
		}
	})

	t.Run("sequential-text-tool-text-stable-order", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		events := []SSEEvent{se(t, "message_start", `{"message":{"id":"m","model":"a","usage":{"input_tokens":1,"output_tokens":0}}}`)}
		events = append(events,
			se(t, "content_block_start", `{"index":7,"content_block":{"type":"text","text":"a"}}`), se(t, "content_block_stop", `{"index":7}`),
			se(t, "content_block_start", `{"index":2,"content_block":{"type":"tool_use","id":"c","name":"f","input":{}}}`), se(t, "content_block_delta", `{"index":2,"delta":{"type":"input_json_delta","partial_json":"{}"}}`), se(t, "content_block_stop", `{"index":2}`),
			se(t, "content_block_start", `{"index":9,"content_block":{"type":"text","text":"b"}}`), se(t, "content_block_stop", `{"index":9}`), se(t, "message_delta", `{"delta":{"stop_reason":"tool_use"},"usage":{"output_tokens":3}}`), se(t, "message_stop", `{}`))
		out := collect(t, c.Convert, events...)
		terminal := mapv(t, object(t, out[len(out)-1].Data)["response"])
		outputs := list(t, terminal["output"])
		if len(outputs) != 3 || mapv(t, outputs[0])["type"] != "message" || mapv(t, outputs[1])["type"] != "function_call" || mapv(t, outputs[2])["type"] != "message" {
			t.Fatalf("outputs=%#v", outputs)
		}
		for i, x := range outputs {
			if mapv(t, x)["status"] != "completed" {
				t.Fatalf("output %d=%#v", i, x)
			}
		}
	})

	t.Run("cache-read-create-terminal-and-malformed", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		out := collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a","usage":{"input_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4,"output_tokens":0}}}`), se(t, "message_delta", `{"delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":5}}`), se(t, "message_stop", `{}`))
		u := mapv(t, mapv(t, object(t, out[len(out)-1].Data)["response"])["usage"])
		if u["input_tokens"] != float64(9) || mapv(t, u["input_tokens_details"])["cached_tokens"] != float64(3) {
			t.Fatalf("usage=%#v", u)
		}
		for _, bad := range []string{`{"input_tokens":-1,"output_tokens":0}`, `{"input_tokens":1.5,"output_tokens":0}`, `{"input_tokens":0,"cache_creation_input_tokens":"1","output_tokens":0}`} {
			x := NewAnthropicToResponsesStream("g")
			convertError(t, x.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a","usage":`+bad+`}}`), ErrInvalidWireData)
		}
	})
}

// Source parity: Python/JS Responses stream multi-item lifecycle, closure-on-error, and terminal policies.
func TestV0ResponsesToAnthropicStreamPayloadMatrix(t *testing.T) {
	t.Run("fragmented-unicode-tool-arguments", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		out := collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`), se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc","call_id":"c","name":"f"}}`), se(t, "response.function_call_arguments.delta", `{"output_index":0,"delta":"{\"雪\":\""}`), se(t, "response.function_call_arguments.delta", `{"output_index":0,"delta":"☃\"}"}`), se(t, "response.function_call_arguments.done", `{"output_index":0,"arguments":"{\"雪\":\"☃\"}"}`), se(t, "response.output_item.done", `{"output_index":0,"item":{"type":"function_call","id":"fc","arguments":"{\"雪\":\"☃\"}"}}`), se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":{"input_tokens":0,"output_tokens":2}}}`))
		if mapv(t, object(t, out[2].Data)["delta"])["partial_json"] != "{\"雪\":\"" {
			t.Fatalf("events=%#v", out)
		}
	})

	t.Run("multiple-messages-content-parts-stable-indexes", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		out := collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`),
			se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m0","role":"assistant"}}`), se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":"a"}}`), se(t, "response.content_part.done", `{"output_index":0,"content_index":0}`), se(t, "response.content_part.added", `{"output_index":0,"content_index":1,"part":{"type":"output_text","text":"b"}}`), se(t, "response.content_part.done", `{"output_index":0,"content_index":1}`), se(t, "response.output_item.done", `{"output_index":0}`),
			se(t, "response.output_item.added", `{"output_index":1,"item":{"type":"message","id":"m1","role":"assistant"}}`), se(t, "response.content_part.added", `{"output_index":1,"content_index":0,"part":{"type":"output_text","text":"c"}}`), se(t, "response.content_part.done", `{"output_index":1,"content_index":0}`), se(t, "response.output_item.done", `{"output_index":1}`), se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":{"input_tokens":0,"output_tokens":3}}}`))
		var indexes []float64
		for _, e := range out {
			if e.Event == "content_block_start" {
				indexes = append(indexes, object(t, e.Data)["index"].(float64))
			}
		}
		if len(indexes) != 3 || indexes[0] != 0 || indexes[1] != 1 || indexes[2] != 2 {
			t.Fatalf("indexes=%#v", indexes)
		}
	})

	t.Run("cache-terminal-preservation-and-malformed", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		out := collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g","usage":{"input_tokens":9,"output_tokens":0,"input_tokens_details":{"cached_tokens":3}}}}`), se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":{"input_tokens":9,"output_tokens":2,"input_tokens_details":{"cached_tokens":3}}}}`))
		startUsage := mapv(t, mapv(t, object(t, out[0].Data)["message"])["usage"])
		if startUsage["input_tokens"] != float64(6) || startUsage["cache_read_input_tokens"] != float64(3) {
			t.Fatalf("usage=%#v", startUsage)
		}
		for _, usage := range []string{`{"input_tokens":-1,"output_tokens":0}`, `{"input_tokens":1,"output_tokens":0,"input_tokens_details":null}`, `{"input_tokens":1,"output_tokens":0,"total_tokens":2}`} {
			x := NewResponsesToAnthropicStream("a")
			collect(t, x.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`))
			convertError(t, x.Convert, se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":`+usage+`}}`), ErrInvalidWireData)
		}
	})

	t.Run("error-closes-open-text-and-tool", func(t *testing.T) {
		for _, kind := range []string{"text", "tool"} {
			t.Run(kind, func(t *testing.T) {
				c := NewResponsesToAnthropicStream("a")
				inputs := []SSEEvent{se(t, "response.created", `{"response":{"id":"r","model":"g"}}`)}
				if kind == "text" {
					inputs = append(inputs, se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m"}}`), se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":""}}`))
				} else {
					inputs = append(inputs, se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc","call_id":"c","name":"f"}}`))
				}
				inputs = append(inputs, se(t, "response.failed", `{"response":{"error":{"code":"boom","message":"failed"}}}`))
				out := collect(t, c.Convert, inputs...)
				if out[len(out)-2].Event != "content_block_stop" || out[len(out)-1].Event != "error" {
					t.Fatalf("events=%#v", names(out))
				}
			})
		}
	})
}

// Source parity: Python/JS stream FSM rejection matrix and failed-start rollback.
func TestV0ResponsesStreamSequenceMatrix(t *testing.T) {
	startMessage := func(t *testing.T) *ResponsesToAnthropicStream {
		c := NewResponsesToAnthropicStream("a")
		collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`), se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m"}}`))
		return c
	}
	t.Run("output-after-item-done", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			event SSEEvent
		}{
			{"content-added", se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":"x"}}`)},
			{"text-delta", se(t, "response.output_text.delta", `{"output_index":0,"content_index":0,"delta":"x"}`)},
			{"text-done", se(t, "response.output_text.done", `{"output_index":0,"content_index":0,"text":"x"}`)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				c := startMessage(t)
				collect(t, c.Convert, se(t, "response.output_item.done", `{"output_index":0}`))
				convertError(t, c.Convert, tc.event, ErrInvalidSequence)
			})
		}
	})
	t.Run("wrong-item-kind", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`), se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc","call_id":"c","name":"f"}}`))
		convertError(t, c.Convert, se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":"x"}}`), ErrInvalidSequence)
	})
	t.Run("item-id-mismatch", func(t *testing.T) {
		c := startMessage(t)
		convertError(t, c.Convert, se(t, "response.content_part.added", `{"item_id":"other","output_index":0,"content_index":0,"part":{"type":"output_text","text":"x"}}`), ErrInvalidSequence)
	})
	t.Run("duplicate-start", func(t *testing.T) {
		c := startMessage(t)
		convertError(t, c.Convert, se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"message","id":"m"}}`), ErrInvalidSequence)
	})
	t.Run("duplicate-done", func(t *testing.T) {
		c := startMessage(t)
		collect(t, c.Convert, se(t, "response.output_item.done", `{"output_index":0}`))
		convertError(t, c.Convert, se(t, "response.output_item.done", `{"output_index":0}`), ErrInvalidSequence)
	})
	t.Run("done-kind-mismatch", func(t *testing.T) {
		c := startMessage(t)
		convertError(t, c.Convert, se(t, "response.output_item.done", `{"output_index":0,"item":{"type":"function_call","id":"m"}}`), ErrInvalidSequence)
	})
	t.Run("terminal-with-open-block", func(t *testing.T) {
		c := startMessage(t)
		collect(t, c.Convert, se(t, "response.content_part.added", `{"output_index":0,"content_index":0,"part":{"type":"output_text","text":""}}`))
		convertError(t, c.Convert, se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":{"input_tokens":0,"output_tokens":0}}}`), ErrInvalidSequence)
	})

	t.Run("terminal-envelope-matrix", func(t *testing.T) {
		for _, tc := range []struct {
			name, terminal string
			target         error
		}{
			{"id-mismatch", `{"response":{"id":"other","status":"completed","usage":{"input_tokens":0,"output_tokens":0}}}`, ErrInvalidSequence},
			{"id-missing", `{"response":{"status":"completed","usage":{"input_tokens":0,"output_tokens":0}}}`, ErrInvalidWireData},
			{"status-mismatch", `{"response":{"id":"r","status":"incomplete","usage":{"input_tokens":0,"output_tokens":0}}}`, ErrInvalidWireData},
			{"usage-missing", `{"response":{"id":"r","status":"completed"}}`, ErrInvalidWireData},
			{"incomplete-details-missing", `{"response":{"id":"r","status":"incomplete","usage":{"input_tokens":0,"output_tokens":0}}}`, ErrInvalidWireData},
			{"incomplete-details-unsupported", `{"response":{"id":"r","status":"incomplete","incomplete_details":{"reason":"content_filter"},"usage":{"input_tokens":0,"output_tokens":0}}}`, ErrUnsupported},
		} {
			t.Run(tc.name, func(t *testing.T) {
				c := NewResponsesToAnthropicStream("a")
				collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`))
				typ := "response.completed"
				if tc.name == "incomplete-details-missing" || tc.name == "incomplete-details-unsupported" {
					typ = "response.incomplete"
				}
				convertError(t, c.Convert, se(t, typ, tc.terminal), tc.target)
			})
		}
	})
	t.Run("terminal-usage-mismatch", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g","usage":{"input_tokens":4,"output_tokens":0,"input_tokens_details":{"cached_tokens":1}}}}`))
		convertError(t, c.Convert, se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":{"input_tokens":5,"output_tokens":1,"input_tokens_details":{"cached_tokens":1}}}}`), ErrInvalidSequence)
	})
	t.Run("event-after-terminal", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`), se(t, "response.completed", `{"response":{"id":"r","status":"completed","usage":{"input_tokens":0,"output_tokens":0}}}`))
		convertError(t, c.Convert, se(t, "response.in_progress", `{}`), ErrInvalidSequence)
	})
}

// Source parity: Python/JS Anthropic stream sequencing, discriminator, JSON framing, and unsupported lifecycle events.
func TestV0AnthropicStreamSequenceAndFramingMatrix(t *testing.T) {
	t.Run("after-and-duplicate-message-delta", func(t *testing.T) {
		for _, tc := range []struct {
			name  string
			after SSEEvent
		}{
			{"event-after-delta", se(t, "content_block_start", `{"index":0,"content_block":{"type":"text","text":""}}`)},
			{"duplicate-delta", se(t, "message_delta", `{"delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":0}}`)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				c := NewAnthropicToResponsesStream("g")
				collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a"}}`), se(t, "message_delta", `{"delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":0}}`))
				convertError(t, c.Convert, tc.after, ErrInvalidSequence)
			})
		}
	})
	t.Run("wrong-delta-discriminator", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a"}}`), se(t, "content_block_start", `{"index":0,"content_block":{"type":"text","text":""}}`))
		convertError(t, c.Convert, se(t, "content_block_delta", `{"index":0,"delta":{"type":"input_json_delta","partial_json":"{}"}}`), ErrInvalidWireData)
	})
	t.Run("failed-block-start-does-not-poison-state", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a"}}`))
		convertError(t, c.Convert, se(t, "content_block_start", `{"index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{"already":1}}}`), ErrInvalidWireData)
		out := collect(t, c.Convert, se(t, "content_block_start", `{"index":0,"content_block":{"type":"text","text":"ok"}}`))
		if len(out) != 2 {
			t.Fatalf("events=%#v", names(out))
		}
	})
	t.Run("event-after-terminal", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a"}}`), se(t, "message_delta", `{"delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":0}}`), se(t, "message_stop", `{}`))
		convertError(t, c.Convert, se(t, "ping", `{}`), ErrInvalidSequence)
	})

	t.Run("header-data-and-json-framing", func(t *testing.T) {
		cases := []struct {
			name   string
			e      SSEEvent
			target error
		}{
			{"header-only", SSEEvent{Event: "ping", Data: []byte(`{}`)}, nil}, {"data-only", SSEEvent{Data: []byte(`{"type":"ping"}`)}, nil},
			{"mismatch", SSEEvent{Event: "ping", Data: []byte(`{"type":"message_start"}`)}, ErrInvalidWireData}, {"missing", SSEEvent{Data: []byte(`{}`)}, ErrInvalidWireData},
			{"malformed", SSEEvent{Event: "ping", Data: []byte(`{`)}, ErrInvalidWireData}, {"trailing", SSEEvent{Event: "ping", Data: []byte(`{} {}`)}, ErrInvalidWireData},
			{"nonobject", SSEEvent{Event: "ping", Data: []byte(`[]`)}, ErrInvalidWireData}, {"done", SSEEvent{Data: []byte(`[DONE]`)}, ErrUnsupported},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				c := NewAnthropicToResponsesStream("g")
				out, err := c.Convert(tc.e)
				if tc.target == nil {
					if err != nil || len(out) != 0 {
						t.Fatalf("out=%#v err=%v", out, err)
					}
					return
				}
				requireErrorIs(t, err, tc.target)
			})
		}
	})
}

// Source parity: explicit unsupported Responses/Anthropic stream event and content-block boundaries.
func TestV0StreamUnsupportedMatrix(t *testing.T) {
	t.Run("responses-events", func(t *testing.T) {
		for _, typ := range []string{"response.reasoning_summary_part.added", "response.reasoning_summary_part.done", "response.reasoning_summary_text.delta", "response.reasoning_summary_text.done", "response.refusal.delta", "response.refusal.done", "response.output_text.annotation.added", "response.output_text.annotation.delta", "response.output_text.annotation.done"} {
			t.Run(typ, func(t *testing.T) {
				c := NewResponsesToAnthropicStream("a")
				_, err := c.Convert(SSEEvent{Event: typ, Data: []byte(`{}`)})
				requireErrorIs(t, err, ErrUnsupported)
			})
		}
	})
	t.Run("anthropic-blocks", func(t *testing.T) {
		for _, typ := range []string{"thinking", "redacted_thinking", "signature", "citation", "server_tool_use", "web_search_tool_result", "tool_search_tool_result"} {
			t.Run(typ, func(t *testing.T) {
				c := NewAnthropicToResponsesStream("g")
				collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a"}}`))
				convertError(t, c.Convert, se(t, "content_block_start", `{"index":0,"content_block":{"type":"`+typ+`"}}`), ErrUnsupported)
			})
		}
	})
	t.Run("anthropic-citation-delta", func(t *testing.T) {
		c := NewAnthropicToResponsesStream("g")
		collect(t, c.Convert, se(t, "message_start", `{"message":{"id":"m","model":"a"}}`), se(t, "content_block_start", `{"index":0,"content_block":{"type":"text","text":""}}`))
		convertError(t, c.Convert, se(t, "content_block_delta", `{"index":0,"delta":{"type":"citation_delta"}}`), ErrUnsupported)
	})
	t.Run("incomplete-truncated-tool-policy", func(t *testing.T) {
		c := NewResponsesToAnthropicStream("a")
		collect(t, c.Convert, se(t, "response.created", `{"response":{"id":"r","model":"g"}}`), se(t, "response.output_item.added", `{"output_index":0,"item":{"type":"function_call","id":"fc","call_id":"c","name":"f"}}`), se(t, "response.function_call_arguments.delta", `{"output_index":0,"delta":"{"}`))
		convertError(t, c.Convert, se(t, "response.incomplete", `{"response":{"id":"r","status":"incomplete","usage":{"input_tokens":0,"output_tokens":1}}}`), ErrUnsupported)
	})

}
