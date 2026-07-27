package messagetranslators

import (
	"errors"
	"fmt"
	"testing"
)

func chatCompletionsRequest(body string) []byte { return []byte(body) }

func TestChatCompletionsRequestSystemScalarsAndStops(t *testing.T) {
	t.Run("anthropic-system-string-and-blocks", func(t *testing.T) {
		for _, tc := range []struct{ name, system, want string }{
			{"string", `"rules"`, "rules"},
			{"blocks", `[{"type":"text","text":"a"},{"type":"text","text":"β"},{"type":"text","text":""}]`, "aβ"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				b, err := AnthropicRequestToChatCompletions(chatCompletionsRequest(`{"model":"a","max_tokens":7,"system":`+tc.system+`,"messages":[{"role":"user","content":"u"}]}`), "g")
				if err != nil {
					t.Fatal(err)
				}
				r := object(t, b)
				m := list(t, r["messages"])
				if r["model"] != "g" || mapv(t, m[0])["role"] != "system" || mapv(t, m[0])["content"] != tc.want {
					t.Fatalf("result=%s", b)
				}
			})
		}
	})
	t.Run("chat-completions-leading-system-developer-order", func(t *testing.T) {
		b, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_completion_tokens":9,"temperature":0.25,"top_p":0.75,"messages":[{"role":"system","content":"a"},{"role":"developer","content":[{"type":"text","text":"β"},{"type":"text","text":"c"}]},{"role":"user","content":"u"}]}`), "a")
		if err != nil {
			t.Fatal(err)
		}
		a := object(t, b)
		s := list(t, a["system"])
		if a["model"] != "a" || a["max_tokens"] != float64(9) || a["temperature"] != .25 || a["top_p"] != .75 || len(s) != 3 || mapv(t, s[0])["text"] != "a" || mapv(t, s[2])["text"] != "c" {
			t.Fatalf("result=%s", b)
		}
	})
	t.Run("nonleading-system", func(t *testing.T) {
		_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"system","content":"late"}]}`), "")
		requireErrorIs(t, err, ErrUnsupported)
	})
	t.Run("model-and-required-token-fields", func(t *testing.T) {
		cases := []struct {
			name, body, override string
			target               error
		}{
			{"override-allows-absent-model", `{"max_tokens":1,"messages":[{"role":"user","content":"u"}]}`, "a", nil},
			{"missing-model", `{"max_tokens":1,"messages":[{"role":"user","content":"u"}]}`, "", ErrInvalidWireData},
			{"empty-wire-model-even-override", `{"model":"","max_tokens":1,"messages":[{"role":"user","content":"u"}]}`, "a", ErrInvalidWireData},
			{"missing-limit", `{"model":"g","messages":[{"role":"user","content":"u"}]}`, "", ErrInvalidWireData},
			{"conflicting-limits", `{"model":"g","max_tokens":1,"max_completion_tokens":2,"messages":[{"role":"user","content":"u"}]}`, "", ErrInvalidWireData},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(tc.body), tc.override)
				if tc.target == nil {
					if err != nil {
						t.Fatal(err)
					}
				} else {
					requireErrorIs(t, err, tc.target)
				}
			})
		}
	})
	t.Run("stop-string-array-and-stream-options", func(t *testing.T) {
		for _, tc := range []struct {
			name, stop string
			count      int
		}{{"string", `"END"`, 1}, {"array", `["A","B"]`, 2}} {
			t.Run(tc.name, func(t *testing.T) {
				b, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"stream":true,"stream_options":{"include_usage":true},"stop":`+tc.stop+`,"messages":[{"role":"user","content":"u"}]}`), "")
				if err != nil {
					t.Fatal(err)
				}
				if len(list(t, object(t, b)["stop_sequences"])) != tc.count {
					t.Fatalf("result=%s", b)
				}
			})
		}
		bad := []string{
			`{"model":"g","max_tokens":1,"stream_options":{"include_usage":true},"messages":[{"role":"user","content":"u"}]}`,
			`{"model":"g","max_tokens":1,"stream":true,"stream_options":{},"messages":[{"role":"user","content":"u"}]}`,
			`{"model":"g","max_tokens":1,"stop":"","messages":[{"role":"user","content":"u"}]}`,
			`{"model":"g","max_tokens":1,"stop":["a",2],"messages":[{"role":"user","content":"u"}]}`,
		}
		for i, body := range bad {
			t.Run(fmt.Sprintf("bad-%d", i), func(t *testing.T) {
				_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(body), "")
				requireErrorIs(t, err, ErrInvalidWireData)
			})
		}
	})
}

func TestChatCompletionsRequestToolChoiceAndDefinitions(t *testing.T) {
	t.Run("choices-both-directions", func(t *testing.T) {
		for _, tc := range []struct{ name, chatCompletions, anth string }{
			{"auto", `"auto"`, `{"type":"auto"}`}, {"required", `"required"`, `{"type":"any"}`}, {"none", `"none"`, `{"type":"none"}`},
			{"named", `{"type":"function","function":{"name":"f"}}`, `{"type":"tool","name":"f"}`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				cb, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"tool_choice":`+tc.chatCompletions+`,"messages":[{"role":"user","content":"u"}]}`), "a")
				if err != nil {
					t.Fatal(err)
				}
				if object(t, cb)["tool_choice"] == nil {
					t.Fatal("choice dropped")
				}
				ab, err := AnthropicRequestToChatCompletions(chatCompletionsRequest(`{"model":"a","max_tokens":1,"tool_choice":`+tc.anth+`,"messages":[{"role":"user","content":"u"}]}`), "g")
				if err != nil {
					t.Fatal(err)
				}
				if object(t, ab)["tool_choice"] == nil {
					t.Fatal("choice dropped")
				}
			})
		}
	})
	t.Run("malformed-choices", func(t *testing.T) {
		cases := []struct {
			name, body string
			fn         func([]byte, string) ([]byte, error)
			target     error
		}{
			{"chat-completions-unknown", `{"model":"g","max_tokens":1,"tool_choice":"any","messages":[{"role":"user","content":"u"}]}`, ChatCompletionsRequestToAnthropic, ErrUnsupported},
			{"chat-completions-number", `{"model":"g","max_tokens":1,"tool_choice":1,"messages":[{"role":"user","content":"u"}]}`, ChatCompletionsRequestToAnthropic, ErrInvalidWireData},
			{"chat-completions-named-missing", `{"model":"g","max_tokens":1,"tool_choice":{"type":"function","function":{}},"messages":[{"role":"user","content":"u"}]}`, ChatCompletionsRequestToAnthropic, ErrInvalidWireData},
			{"anth-string", `{"model":"a","max_tokens":1,"tool_choice":"auto","messages":[{"role":"user","content":"u"}]}`, AnthropicRequestToChatCompletions, ErrInvalidWireData},
			{"anth-unknown", `{"model":"a","max_tokens":1,"tool_choice":{"type":"required"},"messages":[{"role":"user","content":"u"}]}`, AnthropicRequestToChatCompletions, ErrUnsupported},
			{"anth-disable-parallel", `{"model":"a","max_tokens":1,"tool_choice":{"type":"auto","disable_parallel_tool_use":true},"messages":[{"role":"user","content":"u"}]}`, AnthropicRequestToChatCompletions, ErrUnsupported},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := tc.fn(chatCompletionsRequest(tc.body), "")
				requireErrorIs(t, err, tc.target)
			})
		}
	})
	t.Run("definitions-defaults-description-and-unicode-schema", func(t *testing.T) {
		b, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"tools":[{"type":"function","function":{"name":"f","description":"desc"}},{"type":"function","function":{"name":"雪","parameters":{"type":"object","properties":{"☃":{"type":"string"}}}}}],"messages":[{"role":"user","content":"u"}]}`), "a")
		if err != nil {
			t.Fatal(err)
		}
		tools := list(t, object(t, b)["tools"])
		if mapv(t, mapv(t, tools[0])["input_schema"])["type"] != "object" || mapv(t, tools[0])["description"] != "desc" || mapv(t, tools[1])["name"] != "雪" {
			t.Fatalf("tools=%#v", tools)
		}
		back, err := AnthropicRequestToChatCompletions(chatCompletionsRequest(`{"model":"a","max_tokens":1,"tools":[{"name":"f","description":"d"}],"messages":[{"role":"user","content":"u"}]}`), "g")
		if err != nil {
			t.Fatal(err)
		}
		fn := mapv(t, mapv(t, list(t, object(t, back)["tools"])[0])["function"])
		if fn["description"] != "d" || mapv(t, fn["parameters"])["type"] != "object" {
			t.Fatalf("tool=%#v", fn)
		}
	})
	t.Run("malformed-strict-and-hosted-tools", func(t *testing.T) {
		cases := []struct {
			name, tool string
			fn         func([]byte, string) ([]byte, error)
			target     error
		}{
			{"chat-completions-hosted", `{"type":"web_search"}`, ChatCompletionsRequestToAnthropic, ErrUnsupported},
			{"chat-completions-function-array", `{"type":"function","function":[]}`, ChatCompletionsRequestToAnthropic, ErrInvalidWireData},
			{"chat-completions-description", `{"type":"function","function":{"name":"f","description":2}}`, ChatCompletionsRequestToAnthropic, ErrInvalidWireData},
			{"chat-completions-schema", `{"type":"function","function":{"name":"f","parameters":[]}}`, ChatCompletionsRequestToAnthropic, ErrInvalidWireData},
			{"chat-completions-strict", `{"type":"function","function":{"name":"f","strict":true}}`, ChatCompletionsRequestToAnthropic, ErrUnsupported},
			{"anth-hosted", `{"type":"web_search_20250305","name":"search"}`, AnthropicRequestToChatCompletions, ErrUnsupported},
			{"anth-description", `{"name":"f","description":2}`, AnthropicRequestToChatCompletions, ErrInvalidWireData},
			{"anth-schema", `{"name":"f","input_schema":[]}`, AnthropicRequestToChatCompletions, ErrInvalidWireData},
			{"anth-strict", `{"name":"f","strict":true}`, AnthropicRequestToChatCompletions, ErrUnsupported},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				prefix := `{"model":"g","max_tokens":1,"tools":[` + tc.tool + `],"messages":[{"role":"user","content":"u"}]}`
				_, err := tc.fn(chatCompletionsRequest(prefix), "")
				requireErrorIs(t, err, tc.target)
			})
		}
	})
}

func TestChatCompletionsRequestContentImagesAndHistory(t *testing.T) {
	t.Run("string-multipart-url-base64-roundtrip", func(t *testing.T) {
		body := `{"model":"g","max_tokens":2,"messages":[{"role":"user","content":[{"type":"text","text":"look"},{"type":"image_url","image_url":{"url":"https://example/x.png"}},{"type":"image_url","image_url":{"url":"data:image/png;base64,4piD"}}]}]}`
		b, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(body), "a")
		if err != nil {
			t.Fatal(err)
		}
		content := list(t, mapv(t, list(t, object(t, b)["messages"])[0])["content"])
		if len(content) != 3 || mapv(t, mapv(t, content[1])["source"])["type"] != "url" || mapv(t, mapv(t, content[2])["source"])["data"] != "4piD" {
			t.Fatalf("content=%#v", content)
		}
		back, err := AnthropicRequestToChatCompletions(b, "g2")
		if err != nil {
			t.Fatal(err)
		}
		if object(t, back)["model"] != "g2" {
			t.Fatalf("back=%s", back)
		}
	})
	t.Run("malformed-images", func(t *testing.T) {
		parts := []struct {
			name, part string
			target     error
		}{
			{"missing-image-url", `{"type":"image_url"}`, ErrInvalidWireData}, {"empty-url", `{"type":"image_url","image_url":{"url":""}}`, ErrInvalidWireData},
			{"detail", `{"type":"image_url","image_url":{"url":"https://x","detail":"high"}}`, ErrUnsupported}, {"bad-data-source", `{"type":"image_url","image_url":{"url":"data:image/png,abc"}}`, ErrUnsupported},
			{"bad-base64", `{"type":"image_url","image_url":{"url":"data:image/png;base64,!"}}`, ErrInvalidWireData},
		}
		for _, tc := range parts {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"messages":[{"role":"user","content":[`+tc.part+`]}]}`), "a")
				requireErrorIs(t, err, tc.target)
			})
		}
		_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"image_url","image_url":{"url":"https://x"}}]}]}`), "a")
		requireErrorIs(t, err, ErrUnsupported)
	})
	t.Run("two-calls-reordered-results-empty-resumed-unicode", func(t *testing.T) {
		body := `{"model":"g","max_tokens":5,"messages":[{"role":"user","content":"go"},{"role":"assistant","content":"pre","tool_calls":[{"id":"c1","type":"function","function":{"name":"f","arguments":"{\"雪\":\"☃\"}"}},{"id":"c2","type":"function","function":{"name":"g","arguments":"{}"}}]},{"role":"tool","tool_call_id":"c2","content":""},{"role":"tool","tool_call_id":"c1","content":[{"type":"text","text":"one"},{"type":"text","text":"!"}]},{"role":"user","content":"resumed"}]}`
		b, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(body), "a")
		if err != nil {
			t.Fatal(err)
		}
		ms := list(t, object(t, b)["messages"])
		if len(ms) != 3 {
			t.Fatalf("messages=%#v", ms)
		}
		assistant := list(t, mapv(t, ms[1])["content"])
		if mapv(t, mapv(t, assistant[1])["input"])["雪"] != "☃" {
			t.Fatalf("messages=%#v", ms)
		}
		results := list(t, mapv(t, ms[2])["content"])
		if len(results) != 3 || mapv(t, results[0])["tool_use_id"] != "c2" || mapv(t, results[0])["content"] != "" || mapv(t, results[1])["content"] != "one!" || mapv(t, results[2])["text"] != "resumed" {
			t.Fatalf("results=%#v", results)
		}
		back, err := AnthropicRequestToChatCompletions(b, "g2")
		if err != nil {
			t.Fatal(err)
		}
		bm := list(t, object(t, back)["messages"])
		if len(bm) != 5 || mapv(t, bm[2])["tool_call_id"] != "c2" || mapv(t, bm[4])["content"] == nil {
			t.Fatalf("back=%s", back)
		}
	})
	t.Run("call-pairing-and-arguments-errors", func(t *testing.T) {
		cases := []struct {
			name, messages string
			target         error
		}{
			{"duplicate", `[{"role":"user","content":"u"},{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{}"}},{"id":"c","type":"function","function":{"name":"g","arguments":"{}"}}]}]`, ErrInvalidSequence},
			{"orphan", `[{"role":"user","content":"u"},{"role":"tool","tool_call_id":"c","content":"x"}]`, ErrInvalidSequence},
			{"unresolved", `[{"role":"user","content":"u"},{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{}"}}]}]`, ErrInvalidSequence},
			{"missing-arguments", `[{"role":"user","content":"u"},{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"f"}}]}]`, ErrInvalidWireData},
			{"non-object", `[{"role":"user","content":"u"},{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"[]"}}]}]`, ErrInvalidWireData},
			{"malformed", `[{"role":"user","content":"u"},{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{"}}]}]`, ErrInvalidWireData},
			{"trailing", `[{"role":"user","content":"u"},{"role":"assistant","tool_calls":[{"id":"c","type":"function","function":{"name":"f","arguments":"{} {}"}}]}]`, ErrInvalidWireData},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"messages":`+tc.messages+`}`), "a")
				requireErrorIs(t, err, tc.target)
			})
		}
	})
	t.Run("anthropic-ordering-loss", func(t *testing.T) {
		base := `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":%s}]}`
		ok := fmt.Sprintf(base, `[{"type":"text","text":"pre"},{"type":"tool_use","id":"c","name":"f","input":{}}]`)
		// Complete the call so the history itself is valid.
		ok = ok[:len(ok)-2] + `,{"role":"user","content":[{"type":"tool_result","tool_use_id":"c","content":"ok"}]}]}`
		if _, err := AnthropicRequestToChatCompletions(chatCompletionsRequest(ok), "g"); err != nil {
			t.Fatal(err)
		}
		bad := fmt.Sprintf(base, `[{"type":"tool_use","id":"c","name":"f","input":{}},{"type":"text","text":"post"}]`)
		_, err := AnthropicRequestToChatCompletions(chatCompletionsRequest(bad), "g")
		requireErrorIs(t, err, ErrUnsupported)
	})
}

func TestChatCompletionsRequestExplicitUnsupportedMatrix(t *testing.T) {
	for _, field := range []string{"functions", "function_call", "response_format", "parallel_tool_calls", "logprobs", "top_logprobs", "prediction", "modalities", "audio", "service_tier", "store", "metadata", "seed", "user", "reasoning_effort", "web_search_options", "frequency_penalty", "presence_penalty", "logit_bias", "verbosity", "safety_identifier"} {
		t.Run("top-"+field, func(t *testing.T) {
			_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"`+field+`":{},"messages":[{"role":"user","content":"u"}]}`), "a")
			requireErrorIs(t, err, ErrUnsupported)
		})
	}
	for _, tc := range []struct{ name, extra string }{
		{"message-name", `,"name":"x"`}, {"message-function-call", `,"function_call":{}`}, {"message-audio", `,"audio":{}`}, {"message-refusal", `,"refusal":"no"`}, {"message-annotations", `,"annotations":[{}]`},
	} {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"messages":[{"role":"user","content":"u"`+tc.extra+`}]}`), "a")
			requireErrorIs(t, err, ErrUnsupported)
		})
	}
	for _, typ := range []string{"input_audio", "file", "refusal", "document"} {
		t.Run("content-"+typ, func(t *testing.T) {
			_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"messages":[{"role":"user","content":[{"type":"`+typ+`"}]}]}`), "a")
			requireErrorIs(t, err, ErrUnsupported)
		})
	}
	for _, n := range []string{"0", "2"} {
		t.Run("n-"+n, func(t *testing.T) {
			_, err := ChatCompletionsRequestToAnthropic(chatCompletionsRequest(`{"model":"g","max_tokens":1,"n":`+n+`,"messages":[{"role":"user","content":"u"}]}`), "a")
			if n == "0" {
				if !errors.Is(err, ErrInvalidWireData) {
					t.Fatalf("err=%v", err)
				}
			} else {
				requireErrorIs(t, err, ErrUnsupported)
			}
		})
	}
}
