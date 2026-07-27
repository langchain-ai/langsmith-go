package messagetranslators

import (
	"errors"
	"fmt"
	"strings"
	"testing"
)

func chatCompletionsChunk(id, model, choices, usage string) SSEEvent {
	if usage == "" {
		usage = "null"
	}
	return SSEEvent{Data: []byte(`{"id":"` + id + `","object":"chat.completion.chunk","created":10,"model":"` + model + `","choices":` + choices + `,"usage":` + usage + `}`)}
}
func chatCompletionsChoice(delta, finish string) string {
	return `[{"index":0,"delta":` + delta + `,"finish_reason":` + finish + `,"logprobs":null}]`
}

func TestChatCompletionsToAnthropicStreamTextLifecycleUsageAndDone(t *testing.T) {
	s := NewChatCompletionsToAnthropicStream("claude")
	out := collect(t, s.Convert,
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"role":"assistant","content":""}`, "null"), ""),
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"content":"hello ☃"}`, "null"), ""),
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, `"stop"`), ""),
		chatCompletionsChunk("c", "g", `[]`, `{"prompt_tokens":9,"completion_tokens":2,"total_tokens":11,"prompt_tokens_details":{"cached_tokens":3}}`),
		SSEEvent{Data: []byte("[DONE]")},
	)
	equalNames(t, out, "message_start", "content_block_start", "content_block_delta", "content_block_stop", "message_delta", "message_stop")
	start := mapv(t, object(t, out[0].Data)["message"])
	if start["role"] != "assistant" || start["model"] != "claude" || mapv(t, start["usage"])["input_tokens"] != float64(0) {
		t.Fatalf("start=%#v", start)
	}
	delta := object(t, out[2].Data)
	if mapv(t, delta["delta"])["text"] != "hello ☃" {
		t.Fatalf("delta=%#v", delta)
	}
	terminal := object(t, out[4].Data)
	if mapv(t, terminal["delta"])["stop_reason"] != "end_turn" || mapv(t, terminal["usage"])["output_tokens"] != float64(2) {
		t.Fatalf("terminal=%#v", terminal)
	}
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	if _, err := s.Convert(SSEEvent{Data: []byte("[DONE]")}); !errors.Is(err, ErrInvalidSequence) {
		t.Fatalf("after terminal=%v", err)
	}
}

func TestChatCompletionsToAnthropicStreamParallelToolsUnicodeAndFinishMappings(t *testing.T) {
	s := NewChatCompletionsToAnthropicStream("")
	out := collect(t, s.Convert,
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"role":"assistant"}`, "null"), ""),
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"content":"pre"}`, "null"), ""),
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":"{\"雪\":\""}},{"index":1,"id":"b","type":"function","function":{"name":"g","arguments":"{\"n\":"}}]}`, "null"), ""),
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":1,"function":{"arguments":"1}"}},{"index":0,"function":{"arguments":"☃\"}"}}]}`, "null"), ""),
		chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, `"tool_calls"`), ""),
		SSEEvent{Data: []byte("[DONE]")},
	)
	if names(out)[0] != "message_start" || names(out)[len(out)-1] != "message_stop" {
		t.Fatalf("events=%#v", names(out))
	}
	var starts, stops int
	for _, e := range out {
		if e.Event == "content_block_start" {
			starts++
		}
		if e.Event == "content_block_stop" {
			stops++
		}
	}
	if starts != 3 || stops != 3 {
		t.Fatalf("starts=%d stops=%d events=%#v", starts, stops, names(out))
	}
	terminal := object(t, out[len(out)-2].Data)
	if mapv(t, terminal["delta"])["stop_reason"] != "tool_use" {
		t.Fatalf("terminal=%#v", terminal)
	}

	for _, tc := range []struct{ finish, want string }{{`"stop"`, "end_turn"}, {`"length"`, "max_tokens"}} {
		x := NewChatCompletionsToAnthropicStream("a")
		got := collect(t, x.Convert, chatCompletionsChunk("x", "g", chatCompletionsChoice(`{"role":"assistant","content":"x"}`, "null"), ""), chatCompletionsChunk("x", "g", chatCompletionsChoice(`{}`, tc.finish), ""), SSEEvent{Data: []byte("[DONE]")})
		if mapv(t, object(t, got[len(got)-2].Data)["delta"])["stop_reason"] != tc.want {
			t.Fatalf("finish=%s", tc.finish)
		}
	}
}

func TestChatCompletionsToAnthropicStreamEnvelopeAndSequenceMatrix(t *testing.T) {
	start := func(t *testing.T) *ChatCompletionsToAnthropicStream {
		t.Helper()
		s := NewChatCompletionsToAnthropicStream("a")
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"role":"assistant"}`, "null"), ""))
		return s
	}
	t.Run("required-envelope", func(t *testing.T) {
		base := string(chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"role":"assistant"}`, "null"), "").Data)
		for _, tc := range []struct{ name, body string }{
			{"object-missing", strings.Replace(base, `"object":"chat.completion.chunk",`, "", 1)},
			{"object-wrong", strings.Replace(base, `chat.completion.chunk`, `chat.completion`, 1)},
			{"id", strings.Replace(base, `"id":"c",`, "", 1)}, {"model", strings.Replace(base, `"model":"g",`, "", 1)},
			{"created", strings.Replace(base, `"created":10,`, "", 1)}, {"created-fraction", strings.Replace(base, `"created":10`, `"created":1.5`, 1)},
			{"choices", strings.Replace(base, chatCompletionsChoice(`{"role":"assistant"}`, "null"), `null`, 1)},
		} {
			t.Run(tc.name, func(t *testing.T) {
				s := NewChatCompletionsToAnthropicStream("a")
				convertError(t, s.Convert, SSEEvent{Data: []byte(tc.body)}, ErrInvalidWireData)
			})
		}
	})
	t.Run("identity-changes", func(t *testing.T) {
		for _, tc := range []struct {
			name, id, model string
			created         int
		}{{"id", "other", "g", 10}, {"model", "c", "other", 10}, {"created", "c", "g", 11}} {
			t.Run(tc.name, func(t *testing.T) {
				s := start(t)
				raw := string(chatCompletionsChunk(tc.id, tc.model, chatCompletionsChoice(`{}`, "null"), "").Data)
				raw = strings.Replace(raw, `"created":10`, fmt.Sprintf(`"created":%d`, tc.created), 1)
				convertError(t, s.Convert, SSEEvent{Data: []byte(raw)}, ErrInvalidSequence)
			})
		}
	})
	t.Run("choice-and-role-sequences", func(t *testing.T) {
		cases := []struct {
			name    string
			prepare bool
			choices string
			target  error
		}{
			{"two-choices", true, `[{"index":0,"delta":{},"finish_reason":null},{"index":0,"delta":{},"finish_reason":null}]`, ErrUnsupported},
			{"index-one", true, `[{"index":1,"delta":{},"finish_reason":null}]`, ErrUnsupported},
			{"choice-not-object", true, `[1]`, ErrInvalidWireData}, {"delta-not-object", true, `[{"index":0,"delta":null,"finish_reason":null}]`, ErrInvalidWireData},
			{"duplicate-role", true, chatCompletionsChoice(`{"role":"assistant"}`, "null"), ErrInvalidSequence},
			{"content-before-role", false, chatCompletionsChoice(`{"content":"x"}`, "null"), ErrInvalidSequence},
			{"wrong-role", false, chatCompletionsChoice(`{"role":"user"}`, "null"), ErrInvalidSequence},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				var s *ChatCompletionsToAnthropicStream
				if tc.prepare {
					s = start(t)
				} else {
					s = NewChatCompletionsToAnthropicStream("a")
				}
				convertError(t, s.Convert, chatCompletionsChunk("c", "g", tc.choices, ""), tc.target)
			})
		}
	})
	t.Run("finish-consistency", func(t *testing.T) {
		for _, tc := range []struct {
			name, finish string
			tools        bool
			target       error
		}{
			{"stop-with-tool", `"stop"`, true, ErrInvalidSequence}, {"tool-without-tool", `"tool_calls"`, false, ErrInvalidSequence},
			{"content-filter", `"content_filter"`, false, ErrUnsupported}, {"function-call", `"function_call"`, false, ErrUnsupported}, {"unknown", `"wat"`, false, ErrInvalidWireData},
		} {
			t.Run(tc.name, func(t *testing.T) {
				s := start(t)
				if tc.tools {
					collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":"{}"}}]}`, "null"), ""))
				}
				convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, tc.finish), ""), tc.target)
			})
		}
	})
	t.Run("done-and-finish-truncation", func(t *testing.T) {
		s := start(t)
		convertError(t, s.Convert, SSEEvent{Data: []byte("[DONE]")}, ErrInvalidSequence)
		if !errors.Is(s.Finish(), ErrTruncatedStream) {
			t.Fatalf("finish=%v", s.Finish())
		}
		x := start(t)
		collect(t, x.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, `"stop"`), ""))
		if !errors.Is(x.Finish(), ErrTruncatedStream) {
			t.Fatalf("finish=%v", x.Finish())
		}
	})
}

func TestChatCompletionsToAnthropicStreamToolDeltaUsageErrorsAndRollback(t *testing.T) {
	start := func(t *testing.T) *ChatCompletionsToAnthropicStream {
		s := NewChatCompletionsToAnthropicStream("a")
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"role":"assistant"}`, "null"), ""))
		return s
	}
	t.Run("tool-delta-errors", func(t *testing.T) {
		cases := []struct {
			name, delta string
			target      error
		}{
			{"empty", `{"tool_calls":[]}`, ErrInvalidWireData}, {"missing-index", `{"tool_calls":[{"id":"a","type":"function","function":{"name":"f","arguments":"{}"}}]}`, ErrInvalidWireData},
			{"gap-index", `{"tool_calls":[{"index":1,"id":"a","type":"function","function":{"name":"f","arguments":"{}"}}]}`, ErrInvalidSequence},
			{"missing-id", `{"tool_calls":[{"index":0,"type":"function","function":{"name":"f","arguments":"{}"}}]}`, ErrInvalidWireData},
			{"wrong-type", `{"tool_calls":[{"index":0,"id":"a","type":"custom","function":{"name":"f","arguments":"{}"}}]}`, ErrInvalidWireData},
			{"missing-name", `{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"arguments":"{}"}}]}`, ErrInvalidWireData},
			{"arguments-number", `{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":2}}]}`, ErrInvalidWireData},
			{"duplicate-index-event", `{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":""}},{"index":0,"function":{"arguments":"{}"}}]}`, ErrInvalidSequence},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				s := start(t)
				convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(tc.delta, "null"), ""), tc.target)
			})
		}
	})
	t.Run("identity-changes-and-duplicate-ids", func(t *testing.T) {
		s := start(t)
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":"{}"}}]}`, "null"), ""))
		convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":1,"id":"a","type":"function","function":{"name":"g","arguments":"{}"}}]}`, "null"), ""), ErrInvalidSequence)
		convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"id":"other","function":{"arguments":""}}]}`, "null"), ""), ErrInvalidSequence)
	})
	t.Run("malformed-arguments-at-finish-and-rollback", func(t *testing.T) {
		s := start(t)
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":"{"}}]}`, "null"), ""))
		convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, `"tool_calls"`), ""), ErrInvalidWireData)
		// The failed finish must not poison state; complete JSON and finish again.
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"function":{"arguments":"}"}}]}`, "null"), ""), chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, `"tool_calls"`), ""), SSEEvent{Data: []byte("[DONE]")})
		if err := s.Finish(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("usage-sequence-duplicate-and-mismatch", func(t *testing.T) {
		u := `{"prompt_tokens":2,"completion_tokens":1,"total_tokens":3}`
		s := start(t)
		convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, "null"), u), ErrInvalidSequence)
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{}`, `"stop"`), ""), chatCompletionsChunk("c", "g", `[]`, u))
		convertError(t, s.Convert, chatCompletionsChunk("c", "g", `[]`, u), ErrInvalidSequence)
		x := NewChatCompletionsToAnthropicStream("a")
		// A first and finishing chunk may carry usage, but a changed usage-only chunk is duplicate/invalid.
		collect(t, x.Convert, chatCompletionsChunk("x", "g", chatCompletionsChoice(`{"role":"assistant","content":"x"}`, `"stop"`), u))
		convertError(t, x.Convert, chatCompletionsChunk("x", "g", `[]`, `{"prompt_tokens":3,"completion_tokens":1,"total_tokens":4}`), ErrInvalidSequence)
	})
	t.Run("unsupported-deltas-and-error-closes-blocks", func(t *testing.T) {
		for _, field := range []string{"audio", "function_call", "refusal"} {
			t.Run(field, func(t *testing.T) {
				s := start(t)
				convertError(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"`+field+`":{}}`, "null"), ""), ErrUnsupported)
			})
		}
		s := start(t)
		collect(t, s.Convert, chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"content":"open"}`, "null"), ""), chatCompletionsChunk("c", "g", chatCompletionsChoice(`{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":"{"}}]}`, "null"), ""))
		out := collect(t, s.Convert, SSEEvent{Event: "error", Data: []byte(`{"error":{"message":"boom","code":"server_error"}}`)})
		if len(out) != 3 || out[0].Event != "content_block_stop" || out[1].Event != "content_block_stop" || out[2].Event != "error" {
			t.Fatalf("events=%#v", names(out))
		}
		if err := s.Finish(); err != nil {
			t.Fatal(err)
		}
	})
}

func TestAnthropicToChatCompletionsStreamTextToolsSchemaAndDone(t *testing.T) {
	t.Run("nonempty-initial-text-deltas-finish-usage-done", func(t *testing.T) {
		s := NewAnthropicToChatCompletionsStream("g")
		out := collect(t, s.Convert,
			se(t, "message_start", `{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"a","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":2,"cache_read_input_tokens":3,"output_tokens":0}}}`),
			se(t, "content_block_start", `{"type":"content_block_start","index":7,"content_block":{"type":"text","text":"initial ","cache_control":{"type":"ephemeral"}}}`),
			se(t, "content_block_delta", `{"type":"content_block_delta","index":7,"delta":{"type":"text_delta","text":"☃"}}`),
			se(t, "content_block_stop", `{"type":"content_block_stop","index":7}`),
			se(t, "message_delta", `{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":2}}`),
			se(t, "message_stop", `{"type":"message_stop"}`),
		)
		if len(out) != 6 || string(out[len(out)-1].Data) != "[DONE]" {
			t.Fatalf("out=%#v", out)
		}
		var id string
		for i, e := range out[:len(out)-1] {
			o := object(t, e.Data)
			if i == 0 {
				id = o["id"].(string)
			}
			if o["id"] != id || o["object"] != "chat.completion.chunk" || o["model"] != "g" || o["created"] == nil {
				t.Fatalf("chunk=%#v", o)
			}
		}
		if mapv(t, mapv(t, list(t, object(t, out[1].Data)["choices"])[0])["delta"])["content"] != "initial " {
			t.Fatalf("initial=%s", out[1].Data)
		}
		u := mapv(t, object(t, out[len(out)-2].Data)["usage"])
		if u["prompt_tokens"] != float64(5) || u["completion_tokens"] != float64(2) || u["total_tokens"] != float64(7) {
			t.Fatalf("usage=%#v", u)
		}
		if err := s.Finish(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("multiple-tools-fragmented-unicode-stable-index", func(t *testing.T) {
		s := NewAnthropicToChatCompletionsStream("g")
		out := collect(t, s.Convert,
			se(t, "message_start", `{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"a","content":[],"usage":{"input_tokens":0,"output_tokens":0}}}`),
			se(t, "content_block_start", `{"type":"content_block_start","index":5,"content_block":{"type":"tool_use","id":"a","name":"f","input":{}}}`),
			se(t, "content_block_start", `{"type":"content_block_start","index":2,"content_block":{"type":"tool_use","id":"b","name":"g","input":{}}}`),
			se(t, "content_block_delta", `{"type":"content_block_delta","index":5,"delta":{"type":"input_json_delta","partial_json":"{\"雪\":\""}}`),
			se(t, "content_block_delta", `{"type":"content_block_delta","index":2,"delta":{"type":"input_json_delta","partial_json":"{}"}}`),
			se(t, "content_block_delta", `{"type":"content_block_delta","index":5,"delta":{"type":"input_json_delta","partial_json":"☃\"}"}}`),
			se(t, "content_block_stop", `{"type":"content_block_stop","index":2}`), se(t, "content_block_stop", `{"type":"content_block_stop","index":5}`),
			se(t, "message_delta", `{"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":3}}`), se(t, "message_stop", `{"type":"message_stop"}`))
		indexes := []float64{}
		for _, e := range out {
			if string(e.Data) == "[DONE]" {
				continue
			}
			o := object(t, e.Data)
			cs := list(t, o["choices"])
			if len(cs) == 0 {
				continue
			}
			d := mapv(t, mapv(t, cs[0])["delta"])
			if v := d["tool_calls"]; v != nil {
				indexes = append(indexes, mapv(t, list(t, v)[0])["index"].(float64))
			}
		}
		if len(indexes) != 5 || indexes[0] != 0 || indexes[1] != 1 || indexes[2] != 0 || indexes[3] != 1 || indexes[4] != 0 {
			t.Fatalf("indexes=%#v", indexes)
		}
	})
}

func TestAnthropicToChatCompletionsStreamSequenceMalformedUnsupportedAndRollback(t *testing.T) {
	start := func(t *testing.T) *AnthropicToChatCompletionsStream {
		s := NewAnthropicToChatCompletionsStream("g")
		collect(t, s.Convert, se(t, "message_start", `{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"a","content":[],"usage":{"input_tokens":0,"output_tokens":0}}}`))
		return s
	}
	t.Run("start-envelope", func(t *testing.T) {
		for _, body := range []string{
			`{"type":"message_start","message":{"id":"m","role":"assistant","model":"a","content":[]}}`,
			`{"type":"message_start","message":{"id":"m","type":"message","role":"user","model":"a","content":[]}}`,
			`{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","content":[]}}`,
			`{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"a","content":[{}]}}`,
		} {
			s := NewAnthropicToChatCompletionsStream("g")
			convertError(t, s.Convert, SSEEvent{Event: "message_start", Data: []byte(body)}, ErrInvalidWireData)
		}
	})
	t.Run("duplicate-unknown-stopped-blocks-and-duplicate-id", func(t *testing.T) {
		s := start(t)
		text := se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`)
		collect(t, s.Convert, text)
		convertError(t, s.Convert, text, ErrInvalidSequence)
		convertError(t, s.Convert, se(t, "content_block_delta", `{"type":"content_block_delta","index":9,"delta":{"type":"text_delta","text":"x"}}`), ErrInvalidSequence)
		collect(t, s.Convert, se(t, "content_block_stop", `{"type":"content_block_stop","index":0}`))
		convertError(t, s.Convert, se(t, "content_block_stop", `{"type":"content_block_stop","index":0}`), ErrInvalidSequence)
		x := start(t)
		collect(t, x.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{}}}`), se(t, "content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{}"}}`), se(t, "content_block_stop", `{"type":"content_block_stop","index":0}`))
		convertError(t, x.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":1,"content_block":{"type":"tool_use","id":"c","name":"g","input":{}}}`), ErrInvalidSequence)
	})
	t.Run("text-after-tool-and-wrong-deltas", func(t *testing.T) {
		s := start(t)
		collect(t, s.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{}}}`))
		convertError(t, s.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":1,"content_block":{"type":"text","text":"late"}}`), ErrUnsupported)
		convertError(t, s.Convert, se(t, "content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"text_delta","text":"x"}}`), ErrInvalidWireData)
		x := start(t)
		collect(t, x.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`))
		convertError(t, x.Convert, se(t, "content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{}"}}`), ErrInvalidWireData)
		convertError(t, x.Convert, se(t, "content_block_delta", `{"type":"content_block_delta","index":0,"delta":{"type":"citations_delta"}}`), ErrUnsupported)
	})
	t.Run("failed-start-rolls-back", func(t *testing.T) {
		s := start(t)
		convertError(t, s.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{"bad":1}}}`), ErrInvalidWireData)
		out := collect(t, s.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":"ok"}}`))
		if len(out) != 1 {
			t.Fatalf("out=%#v", out)
		}
	})
	t.Run("terminal-and-finish-errors", func(t *testing.T) {
		for _, tc := range []struct {
			name, delta string
			target      error
		}{
			{"open-block", `{"type":"message_delta","delta":{"stop_reason":"end_turn"},"usage":{"output_tokens":0}}`, ErrInvalidSequence},
			{"missing-stop", `{"type":"message_delta","delta":{},"usage":{"output_tokens":0}}`, ErrInvalidWireData},
			{"bad-stop-sequence", `{"type":"message_delta","delta":{"stop_reason":"stop_sequence","stop_sequence":null},"usage":{"output_tokens":0}}`, ErrInvalidWireData},
			{"unexpected-stop-sequence", `{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":"x"},"usage":{"output_tokens":0}}`, ErrInvalidWireData},
			{"unknown-stop", `{"type":"message_delta","delta":{"stop_reason":"wat"},"usage":{"output_tokens":0}}`, ErrInvalidWireData},
		} {
			t.Run(tc.name, func(t *testing.T) {
				s := start(t)
				if tc.name == "open-block" {
					collect(t, s.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"text","text":""}}`))
				}
				convertError(t, s.Convert, SSEEvent{Event: "message_delta", Data: []byte(tc.delta)}, tc.target)
			})
		}
		s := start(t)
		if !errors.Is(s.Finish(), ErrTruncatedStream) {
			t.Fatal("expected truncated")
		}
		convertError(t, s.Convert, se(t, "message_stop", `{"type":"message_stop"}`), ErrInvalidSequence)
	})
	t.Run("terminal-event-after-terminal-ping-error", func(t *testing.T) {
		s := start(t)
		if out, err := s.Convert(se(t, "ping", `{"type":"ping"}`)); err != nil || len(out) != 0 {
			t.Fatalf("ping out=%#v err=%v", out, err)
		}
		collect(t, s.Convert, se(t, "message_delta", `{"type":"message_delta","delta":{"stop_reason":"end_turn","stop_sequence":null},"usage":{"output_tokens":0}}`), se(t, "message_stop", `{"type":"message_stop"}`))
		convertError(t, s.Convert, se(t, "ping", `{"type":"ping"}`), ErrInvalidSequence)
		x := start(t)
		out := collect(t, x.Convert, se(t, "error", `{"type":"error","error":{"type":"overloaded_error","message":"busy"}}`))
		if len(out) != 1 || object(t, out[0].Data)["error"] == nil {
			t.Fatalf("error=%#v", out)
		}
		if err := x.Finish(); err != nil {
			t.Fatal(err)
		}
	})
	t.Run("explicit-unsupported-blocks", func(t *testing.T) {
		for _, typ := range []string{"thinking", "redacted_thinking", "server_tool_use", "web_search_tool_result", "audio", "document"} {
			t.Run(typ, func(t *testing.T) {
				s := start(t)
				convertError(t, s.Convert, se(t, "content_block_start", `{"type":"content_block_start","index":0,"content_block":{"type":"`+typ+`"}}`), ErrUnsupported)
			})
		}
	})
	_ = fmt.Sprint
}
