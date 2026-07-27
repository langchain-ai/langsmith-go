package messagetranslators

import (
	"strings"
	"testing"
)

func chatCompletionsCompleted(message, finish, usage string) []byte {
	return []byte(`{"id":"chatcmpl-source","object":"chat.completion","created":123,"model":"g","choices":[{"index":0,"message":` + message + `,"finish_reason":"` + finish + `","logprobs":null}],"usage":` + usage + `}`)
}

func TestChatCompletionsCompletedPayloadsIDsAndFinish(t *testing.T) {
	t.Run("text-tool-text-plus-tool-multiple-tools", func(t *testing.T) {
		cases := []struct {
			name, message, finish, stop string
			blocks                      int
		}{
			{"text", `{"role":"assistant","content":"hello"}`, "stop", "end_turn", 1},
			{"tool", `{"role":"assistant","content":null,"tool_calls":[{"id":"c1","type":"function","function":{"name":"f","arguments":"{}"}}]}`, "tool_calls", "tool_use", 1},
			{"text-tool", `{"role":"assistant","content":"pre","tool_calls":[{"id":"c1","type":"function","function":{"name":"f","arguments":"{\"雪\":\"☃\"}"}},{"id":"c2","type":"function","function":{"name":"g","arguments":"{}"}}]}`, "tool_calls", "tool_use", 3},
			{"length", `{"role":"assistant","content":"partial"}`, "length", "max_tokens", 1},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				b, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(tc.message, tc.finish, `{"prompt_tokens":3,"completion_tokens":2,"total_tokens":5}`), "claude")
				if err != nil {
					t.Fatal(err)
				}
				a := object(t, b)
				if a["model"] != "claude" || a["stop_reason"] != tc.stop || len(list(t, a["content"])) != tc.blocks {
					t.Fatalf("result=%s", b)
				}
				b2, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(tc.message, tc.finish, `{"prompt_tokens":3,"completion_tokens":2,"total_tokens":5}`), "claude")
				if err != nil || object(t, b2)["id"] != a["id"] {
					t.Fatalf("ID not deterministic: %s %s", b, b2)
				}
			})
		}
	})
	t.Run("anthropic-to-chat-completions-schema-and-model-override", func(t *testing.T) {
		body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"pre"},{"type":"tool_use","id":"c1","name":"f","input":{"雪":"☃"}},{"type":"tool_use","id":"c2","name":"g","input":{}}],"stop_reason":"tool_use","stop_sequence":null,"usage":{"input_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4,"output_tokens":5}}`)
		b, err := AnthropicResponseToChatCompletions(body, "g2")
		if err != nil {
			t.Fatal(err)
		}
		r := object(t, b)
		if r["object"] != "chat.completion" || r["model"] != "g2" || r["id"] == "" {
			t.Fatalf("result=%s", b)
		}
		if _, ok := r["created"].(float64); !ok {
			t.Fatalf("created type=%T", r["created"])
		}
		choices := list(t, r["choices"])
		c := mapv(t, choices[0])
		m := mapv(t, c["message"])
		if len(choices) != 1 || c["index"] != float64(0) || c["finish_reason"] != "tool_calls" || m["role"] != "assistant" || m["content"] != "pre" || len(list(t, m["tool_calls"])) != 2 {
			t.Fatalf("choice=%#v", c)
		}
		u := mapv(t, r["usage"])
		if u["prompt_tokens"] != float64(9) || u["completion_tokens"] != float64(5) || u["total_tokens"] != float64(14) || mapv(t, u["prompt_tokens_details"])["cached_tokens"] != float64(3) {
			t.Fatalf("usage=%#v", u)
		}
	})
}

func TestChatCompletionsCompletedRequiredDiscriminatorsAndTypes(t *testing.T) {
	valid := string(chatCompletionsCompleted(`{"role":"assistant","content":"x"}`, "stop", `{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}`))
	cases := []struct {
		name, body string
		target     error
	}{
		{"object-missing", strings.Replace(valid, `"object":"chat.completion",`, "", 1), ErrInvalidWireData},
		{"object-wrong", strings.Replace(valid, `"chat.completion"`, `"chat.completion.chunk"`, 1), ErrInvalidWireData},
		{"id-missing", strings.Replace(valid, `"id":"chatcmpl-source",`, "", 1), ErrInvalidWireData},
		{"model-missing-even-override", strings.Replace(valid, `"model":"g",`, "", 1), ErrInvalidWireData},
		{"created-missing", strings.Replace(valid, `"created":123,`, "", 1), ErrInvalidWireData},
		{"created-fraction", strings.Replace(valid, `"created":123`, `"created":1.5`, 1), ErrInvalidWireData},
		{"choices-missing", strings.Replace(valid, `"choices":[{"index":0,"message":{"role":"assistant","content":"x"},"finish_reason":"stop","logprobs":null}],`, "", 1), ErrInvalidWireData},
		{"choices-two", strings.Replace(valid, `"choices":[{`, `"choices":[{"index":0,"message":{"role":"assistant","content":"y"},"finish_reason":"stop"},{`, 1), ErrUnsupported},
		{"index-missing", strings.Replace(valid, `"index":0,`, "", 1), ErrInvalidWireData},
		{"index-one", strings.Replace(valid, `"index":0`, `"index":1`, 1), ErrInvalidWireData},
		{"message-missing", strings.Replace(valid, `"message":{"role":"assistant","content":"x"}`, `"message":null`, 1), ErrInvalidWireData},
		{"role-missing", strings.Replace(valid, `"role":"assistant",`, "", 1), ErrInvalidWireData},
		{"role-user", strings.Replace(valid, `"role":"assistant"`, `"role":"user"`, 1), ErrInvalidWireData},
		{"content-missing", strings.Replace(valid, `,"content":"x"`, "", 1), ErrInvalidWireData},
		{"content-object", strings.Replace(valid, `"content":"x"`, `"content":{}`, 1), ErrInvalidWireData},
		{"finish-null", strings.Replace(valid, `"finish_reason":"stop"`, `"finish_reason":null`, 1), ErrInvalidWireData},
		{"usage-missing", strings.Replace(valid, `,"usage":{"prompt_tokens":1,"completion_tokens":1,"total_tokens":2}`, "", 1), ErrInvalidWireData},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ChatCompletionsResponseToAnthropic([]byte(tc.body), "override")
			requireErrorIs(t, err, tc.target)
		})
	}

	anth := `{"id":"m","type":"message","role":"assistant","model":"a","content":[],"stop_reason":"end_turn","stop_sequence":null,"usage":{"input_tokens":0,"output_tokens":0}}`
	for _, tc := range []struct{ name, body string }{
		{"type", strings.Replace(anth, `"type":"message",`, "", 1)}, {"role", strings.Replace(anth, `"role":"assistant"`, `"role":"user"`, 1)},
		{"id", strings.Replace(anth, `"id":"m",`, "", 1)}, {"model-even-override", strings.Replace(anth, `"model":"a",`, "", 1)},
		{"content", strings.Replace(anth, `"content":[]`, `"content":null`, 1)}, {"usage", strings.Replace(anth, `"usage":{"input_tokens":0,"output_tokens":0}`, `"usage":null`, 1)},
	} {
		t.Run("anthropic-"+tc.name, func(t *testing.T) {
			_, err := AnthropicResponseToChatCompletions([]byte(tc.body), "override")
			requireErrorIs(t, err, ErrInvalidWireData)
		})
	}
}

func TestChatCompletionsCompletedFinishConsistencyAndArguments(t *testing.T) {
	for _, tc := range []struct {
		name, message, finish string
		target                error
	}{
		{"legacy-function-call", `{"role":"assistant","content":"x"}`, "function_call", ErrUnsupported},
		{"unknown", `{"role":"assistant","content":"x"}`, "mystery", ErrInvalidWireData},
		{"malformed-arguments", `{"role":"assistant","content":null,"tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{"}}]}`, "tool_calls", ErrInvalidWireData},
		{"nonobject-arguments", `{"role":"assistant","content":null,"tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"[]"}}]}`, "tool_calls", ErrInvalidWireData},
		{"trailing-arguments", `{"role":"assistant","content":null,"tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{} {}"}}]}`, "tool_calls", ErrInvalidWireData},
		{"duplicate-call-id", `{"role":"assistant","content":null,"tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{}"}},{"id":"c","type":"function","function":{"name":"g","arguments":"{}"}}]}`, "tool_calls", ErrInvalidSequence},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(tc.message, tc.finish, `{"prompt_tokens":0,"completion_tokens":0,"total_tokens":0}`), "a")
			if out != nil {
				t.Fatalf("partial output=%s", out)
			}
			requireErrorIs(t, err, tc.target)
		})
	}

	t.Run("anthropic-duplicate-call-id", func(t *testing.T) {
		body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"tool_use","id":"c","name":"f","input":{}},{"type":"tool_use","id":"c","name":"g","input":{}}],"stop_reason":"tool_use","usage":{"input_tokens":0,"output_tokens":0}}`)
		out, err := AnthropicResponseToChatCompletions(body, "g")
		if out != nil {
			t.Fatalf("partial output=%s", out)
		}
		requireErrorIs(t, err, ErrInvalidSequence)
	})
	t.Run("anthropic-tool-then-text", func(t *testing.T) {
		body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"tool_use","id":"c","name":"f","input":{}},{"type":"text","text":"late"}],"stop_reason":"tool_use","usage":{"input_tokens":0,"output_tokens":0}}`)
		_, err := AnthropicResponseToChatCompletions(body, "g")
		requireErrorIs(t, err, ErrUnsupported)
	})
}

func TestChatCompletionsCompletedUsageMatrix(t *testing.T) {
	t.Run("zero-cache-and-reasoning-intentionally-dropped", func(t *testing.T) {
		for _, usage := range []string{
			`{"prompt_tokens":0,"completion_tokens":0,"total_tokens":0}`,
			`{"prompt_tokens":10,"completion_tokens":3,"total_tokens":13,"prompt_tokens_details":{"cached_tokens":4,"audio_tokens":0},"completion_tokens_details":{"reasoning_tokens":7,"audio_tokens":0,"accepted_prediction_tokens":0,"rejected_prediction_tokens":0}}`,
		} {
			b, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(`{"role":"assistant","content":"x"}`, "stop", usage), "a")
			if err != nil {
				t.Fatal(err)
			}
			u := mapv(t, object(t, b)["usage"])
			if usage[17] == '1' && (u["input_tokens"] != float64(6) || u["cache_read_input_tokens"] != float64(4)) {
				t.Fatalf("usage=%#v", u)
			}
		}
	})
	invalid := []struct {
		name, usage string
		target      error
	}{
		{"total-mismatch", `{"prompt_tokens":1,"completion_tokens":2,"total_tokens":4}`, ErrInvalidWireData},
		{"cached-over-prompt", `{"prompt_tokens":1,"completion_tokens":0,"prompt_tokens_details":{"cached_tokens":2}}`, ErrInvalidWireData},
		{"negative-prompt", `{"prompt_tokens":-1,"completion_tokens":0}`, ErrInvalidWireData},
		{"negative-completion", `{"prompt_tokens":0,"completion_tokens":-1}`, ErrInvalidWireData},
		{"fraction", `{"prompt_tokens":1.5,"completion_tokens":0}`, ErrInvalidWireData},
		{"null-prompt-details", `{"prompt_tokens":1,"completion_tokens":0,"prompt_tokens_details":null}`, ErrInvalidWireData},
		{"null-completion-details", `{"prompt_tokens":1,"completion_tokens":0,"completion_tokens_details":null}`, ErrInvalidWireData},
		{"bad-reasoning", `{"prompt_tokens":1,"completion_tokens":0,"completion_tokens_details":{"reasoning_tokens":1.5}}`, ErrInvalidWireData},
		{"nonzero-audio-input", `{"prompt_tokens":1,"completion_tokens":0,"prompt_tokens_details":{"audio_tokens":1}}`, ErrUnsupported},
		{"nonzero-audio-output", `{"prompt_tokens":1,"completion_tokens":0,"completion_tokens_details":{"audio_tokens":1}}`, ErrUnsupported},
		{"nonzero-prediction", `{"prompt_tokens":1,"completion_tokens":0,"completion_tokens_details":{"accepted_prediction_tokens":1}}`, ErrUnsupported},
		{"overflow-sum", `{"prompt_tokens":9223372036854775807,"completion_tokens":1,"total_tokens":0}`, ErrInvalidWireData},
	}
	for _, tc := range invalid {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(`{"role":"assistant","content":"x"}`, "stop", tc.usage), "a")
			requireErrorIs(t, err, tc.target)
		})
	}
	t.Run("anthropic-overflow", func(t *testing.T) {
		body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[],"stop_reason":"end_turn","usage":{"input_tokens":9223372036854775807,"cache_read_input_tokens":1,"output_tokens":0}}`)
		_, err := AnthropicResponseToChatCompletions(body, "g")
		requireErrorIs(t, err, ErrInvalidWireData)
	})
}

func TestChatCompletionsCompletedUnsupportedMetadataAndNoPartialOutput(t *testing.T) {
	baseUsage := `{"prompt_tokens":0,"completion_tokens":0}`
	for _, tc := range []struct {
		name, message string
		target        error
	}{
		{"function-call", `{"role":"assistant","content":"x","function_call":{}}`, ErrUnsupported},
		{"refusal", `{"role":"assistant","content":"x","refusal":"no"}`, ErrUnsupported},
		{"audio", `{"role":"assistant","content":"x","audio":{}}`, ErrUnsupported},
		{"annotations-nonempty", `{"role":"assistant","content":"x","annotations":[{}]}`, ErrUnsupported},
		{"annotations-null", `{"role":"assistant","content":"x","annotations":null}`, ErrInvalidWireData},
		{"multipart-refusal", `{"role":"assistant","content":[{"type":"text","text":"safe"},{"type":"refusal","refusal":"no"}]}`, ErrUnsupported},
		{"multipart-audio", `{"role":"assistant","content":[{"type":"text","text":"safe"},{"type":"input_audio"}]}`, ErrUnsupported},
	} {
		t.Run(tc.name, func(t *testing.T) {
			out, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(tc.message, "stop", baseUsage), "a")
			if out != nil {
				t.Fatalf("partial=%s", out)
			}
			requireErrorIs(t, err, tc.target)
		})
	}
	if _, err := ChatCompletionsResponseToAnthropic(chatCompletionsCompleted(`{"role":"assistant","content":"x","annotations":[],"refusal":null,"audio":null}`, "stop", baseUsage), "a"); err != nil {
		t.Fatalf("empty defaults: %v", err)
	}

	for _, typ := range []string{"thinking", "redacted_thinking", "refusal", "audio", "server_tool_use", "web_search_tool_result", "document"} {
		t.Run("anthropic-"+typ, func(t *testing.T) {
			body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":[{"type":"text","text":"safe"},{"type":"` + typ + `"}],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":0}}`)
			out, err := AnthropicResponseToChatCompletions(body, "g")
			if out != nil {
				t.Fatalf("partial=%s", out)
			}
			requireErrorIs(t, err, ErrUnsupported)
		})
	}
}
