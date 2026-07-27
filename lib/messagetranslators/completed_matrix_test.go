package messagetranslators

import "testing"

func responsesCompleted(output, usage string) []byte {
	return []byte(`{"id":"resp","object":"response","model":"g","status":"completed","output":` + output + `,"usage":` + usage + `}`)
}
func anthropicCompleted(content, stop, usage string) []byte {
	return []byte(`{"id":"msg","type":"message","role":"assistant","model":"a","content":` + content + `,"stop_reason":"` + stop + `","stop_sequence":null,"usage":` + usage + `}`)
}

// Source parity: Python/JS completed text ordering, aggregation, tools, and envelope schema scenarios.
func TestV0CompletedPayloadAndOrderingMatrix(t *testing.T) {
	t.Run("simple-and-multiple-responses-text", func(t *testing.T) {
		cases := []struct {
			name, output string
			want         []string
		}{
			{"simple", `[{"type":"message","id":"m0","role":"assistant","status":"completed","content":[{"type":"output_text","text":"hello","annotations":[]}]}]`, []string{"hello"}},
			{"multiple-messages-parts", `[{"type":"message","id":"m0","role":"assistant","status":"completed","content":[{"type":"output_text","text":"a","annotations":[]},{"type":"output_text","text":"b"}]},{"type":"message","id":"m1","role":"assistant","status":"completed","content":[{"type":"output_text","text":"c"}]}]`, []string{"a", "b", "c"}},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				got, err := ResponsesResponseToAnthropic(responsesCompleted(tc.output, `{"input_tokens":0,"output_tokens":3}`), "a")
				if err != nil {
					t.Fatal(err)
				}
				cs := list(t, object(t, got)["content"])
				if len(cs) != len(tc.want) {
					t.Fatalf("content=%#v", cs)
				}
				for i, w := range tc.want {
					if mapv(t, cs[i])["text"] != w {
						t.Fatalf("content=%#v", cs)
					}
				}
			})
		}
	})

	t.Run("anthropic-contiguous-text-and-multiple-tools", func(t *testing.T) {
		content := `[{"type":"text","text":"a"},{"type":"text","text":"b"},{"type":"tool_use","id":"c1","name":"f","input":{"x":"☃"}},{"type":"tool_use","id":"c2","name":"g","input":{}},{"type":"text","text":"c"}]`
		got, err := AnthropicResponseToResponses(anthropicCompleted(content, "tool_use", `{"input_tokens":1,"output_tokens":2}`), "g")
		if err != nil {
			t.Fatal(err)
		}
		outs := list(t, object(t, got)["output"])
		if len(outs) != 4 {
			t.Fatalf("outputs=%#v", outs)
		}
		if len(list(t, mapv(t, outs[0])["content"])) != 2 || mapv(t, outs[1])["arguments"] != `{"x":"☃"}` || mapv(t, outs[2])["call_id"] != "c2" || mapv(t, list(t, mapv(t, outs[3])["content"])[0])["text"] != "c" {
			t.Fatalf("outputs=%#v", outs)
		}
	})

	t.Run("required-discriminators-and-fields", func(t *testing.T) {
		cases := []struct {
			name string
			fn   func([]byte, string) ([]byte, error)
			body []byte
		}{
			{"responses-object", ResponsesResponseToAnthropic, []byte(`{"id":"r","object":"list","model":"g","status":"completed","output":[],"usage":{"input_tokens":0,"output_tokens":0}}`)},
			{"responses-id", ResponsesResponseToAnthropic, []byte(`{"model":"g","status":"completed","output":[],"usage":{"input_tokens":0,"output_tokens":0}}`)},
			{"responses-output-type", ResponsesResponseToAnthropic, responsesCompleted(`[{"id":"m","role":"assistant","status":"completed","content":[]}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"responses-message-id", ResponsesResponseToAnthropic, responsesCompleted(`[{"type":"message","role":"assistant","status":"completed","content":[]}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"responses-message-role", ResponsesResponseToAnthropic, responsesCompleted(`[{"type":"message","id":"m","role":"user","status":"completed","content":[]}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"responses-message-status", ResponsesResponseToAnthropic, responsesCompleted(`[{"type":"message","id":"m","role":"assistant","content":[]}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"responses-function-id", ResponsesResponseToAnthropic, responsesCompleted(`[{"type":"function_call","call_id":"c","name":"f","arguments":"{}","status":"completed"}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"responses-function-status", ResponsesResponseToAnthropic, responsesCompleted(`[{"type":"function_call","id":"fc","call_id":"c","name":"f","arguments":"{}"}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"responses-content-type", ResponsesResponseToAnthropic, responsesCompleted(`[{"type":"message","id":"m","role":"assistant","status":"completed","content":[{"text":"x"}]}]`, `{"input_tokens":0,"output_tokens":0}`)},
			{"anthropic-type", AnthropicResponseToResponses, []byte(`{"id":"m","role":"assistant","model":"a","content":[],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":0}}`)},
			{"anthropic-role", AnthropicResponseToResponses, []byte(`{"id":"m","type":"message","role":"user","model":"a","content":[],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":0}}`)},
			{"anthropic-block-type", AnthropicResponseToResponses, anthropicCompleted(`[{"text":"x"}]`, `end_turn`, `{"input_tokens":0,"output_tokens":0}`)},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) { _, err := tc.fn(tc.body, ""); requireErrorIs(t, err, ErrInvalidWireData) })
		}
	})
}

// Source parity: Python/JS completed status/stop and usage validation scenarios.
func TestV0CompletedStatusAndUsageMatrix(t *testing.T) {
	t.Run("positive-status-stop", func(t *testing.T) {
		responses := []struct{ name, status, details, output, want string }{
			{"completed-end-turn", "completed", "", `[]`, "end_turn"},
			{"completed-tool-use", "completed", "", `[{"type":"function_call","id":"fc","call_id":"c","name":"f","arguments":"{}","status":"completed"}]`, "tool_use"},
			{"incomplete-max-tokens", "incomplete", `,"incomplete_details":{"reason":"max_output_tokens"}`, `[]`, "max_tokens"},
		}
		for _, tc := range responses {
			t.Run("responses-"+tc.name, func(t *testing.T) {
				body := []byte(`{"id":"r","model":"g","status":"` + tc.status + `"` + tc.details + `,"output":` + tc.output + `,"usage":{"input_tokens":0,"output_tokens":0}}`)
				got, err := ResponsesResponseToAnthropic(body, "a")
				if err != nil {
					t.Fatal(err)
				}
				if object(t, got)["stop_reason"] != tc.want {
					t.Fatalf("got=%s", got)
				}
			})
		}
		for _, stop := range []string{"end_turn", "stop_sequence", "max_tokens", "tool_use"} {
			t.Run("anthropic-"+stop, func(t *testing.T) {
				content := `[]`
				if stop == "tool_use" {
					content = `[{"type":"tool_use","id":"c","name":"f","input":{}}]`
				}
				got, err := AnthropicResponseToResponses(anthropicCompleted(content, stop, `{"input_tokens":0,"output_tokens":0}`), "g")
				if err != nil {
					t.Fatal(err)
				}
				status := object(t, got)["status"]
				if (stop == "max_tokens") != (status == "incomplete") {
					t.Fatalf("status=%v", status)
				}
			})
		}
	})

	t.Run("anthropic-cache-control-is-ignored", func(t *testing.T) {
		body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","cache_control":{"type":"ephemeral"},"content":[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}],"stop_reason":"end_turn","usage":{"input_tokens":0,"output_tokens":1}}`)
		if _, err := AnthropicResponseToResponses(body, "g"); err != nil {
			t.Fatal(err)
		}
	})

	t.Run("invalid-and-inconsistent-status-stop", func(t *testing.T) {
		for _, tc := range []struct {
			name, body string
			target     error
		}{
			{"responses-status", `{"id":"r","model":"g","status":"failed","output":[],"usage":{"input_tokens":0,"output_tokens":0}}`, ErrInvalidWireData},
			{"completed-with-incomplete-details", `{"id":"r","model":"g","status":"completed","incomplete_details":{"reason":"max_output_tokens"},"output":[],"usage":{"input_tokens":0,"output_tokens":0}}`, ErrInvalidWireData},
			{"incomplete-missing-details", `{"id":"r","model":"g","status":"incomplete","output":[],"usage":{"input_tokens":0,"output_tokens":0}}`, ErrInvalidWireData},
			{"incomplete-bad-details", `{"id":"r","model":"g","status":"incomplete","incomplete_details":"max_output_tokens","output":[],"usage":{"input_tokens":0,"output_tokens":0}}`, ErrInvalidWireData},
			{"incomplete-unknown-reason", `{"id":"r","model":"g","status":"incomplete","incomplete_details":{"reason":"content_filter"},"output":[],"usage":{"input_tokens":0,"output_tokens":0}}`, ErrUnsupported},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ResponsesResponseToAnthropic([]byte(tc.body), "a")
				requireErrorIs(t, err, tc.target)
			})
		}
		for _, tc := range []struct {
			name, content, stop string
			target              error
		}{
			{"tool-stop-no-tool", `[]`, `tool_use`, ErrInvalidSequence}, {"tool-with-end-turn", `[{"type":"tool_use","id":"c","name":"f","input":{}}]`, `end_turn`, ErrInvalidSequence},
			{"unsupported-stop", `[]`, `pause_turn`, ErrUnsupported}, {"unknown-stop", `[]`, `mystery`, ErrInvalidWireData},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicResponseToResponses(anthropicCompleted(tc.content, tc.stop, `{"input_tokens":0,"output_tokens":0}`), "g")
				requireErrorIs(t, err, tc.target)
			})
		}
	})

	t.Run("usage-positive", func(t *testing.T) {
		for _, usage := range []string{`{"input_tokens":0,"output_tokens":0}`, `{"input_tokens":10,"output_tokens":2,"total_tokens":12,"input_tokens_details":{"cached_tokens":3},"output_tokens_details":{"reasoning_tokens":7}}`} {
			got, err := ResponsesResponseToAnthropic(responsesCompleted(`[]`, usage), "a")
			if err != nil {
				t.Fatal(err)
			}
			if object(t, got)["usage"] == nil {
				t.Fatal("usage dropped")
			}
		}
		got, err := AnthropicResponseToResponses(anthropicCompleted(`[]`, `end_turn`, `{"input_tokens":2,"cache_read_input_tokens":3,"cache_creation_input_tokens":4,"output_tokens":5}`), "g")
		if err != nil {
			t.Fatal(err)
		}
		u := mapv(t, object(t, got)["usage"])
		if u["input_tokens"] != float64(9) || u["total_tokens"] != float64(14) {
			t.Fatalf("usage=%#v", u)
		}
	})

	t.Run("usage-invalid", func(t *testing.T) {
		for _, tc := range []struct{ name, usage string }{
			{"total-mismatch", `{"input_tokens":1,"output_tokens":2,"total_tokens":4}`}, {"cached-over-input", `{"input_tokens":1,"output_tokens":0,"input_tokens_details":{"cached_tokens":2}}`},
			{"negative", `{"input_tokens":-1,"output_tokens":0}`}, {"fractional", `{"input_tokens":1.5,"output_tokens":0}`}, {"string", `{"input_tokens":"1","output_tokens":0}`},
			{"null-input-details", `{"input_tokens":1,"output_tokens":0,"input_tokens_details":null}`}, {"string-input-details", `{"input_tokens":1,"output_tokens":0,"input_tokens_details":"x"}`},
			{"null-output-details", `{"input_tokens":1,"output_tokens":0,"output_tokens_details":null}`}, {"bad-reasoning", `{"input_tokens":1,"output_tokens":0,"output_tokens_details":{"reasoning_tokens":"1"}}`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ResponsesResponseToAnthropic(responsesCompleted(`[]`, tc.usage), "a")
				requireErrorIs(t, err, ErrInvalidWireData)
			})
		}
		for _, tc := range []struct{ name, usage string }{
			{"negative-read", `{"input_tokens":0,"output_tokens":0,"cache_read_input_tokens":-1}`}, {"fractional-create", `{"input_tokens":0,"output_tokens":0,"cache_creation_input_tokens":1.5}`}, {"null-read", `{"input_tokens":0,"output_tokens":0,"cache_read_input_tokens":null}`},
		} {
			t.Run("anthropic-"+tc.name, func(t *testing.T) {
				_, err := AnthropicResponseToResponses(anthropicCompleted(`[]`, `end_turn`, tc.usage), "g")
				requireErrorIs(t, err, ErrInvalidWireData)
			})
		}
	})
}

// Source parity: Python/JS malformed tools, citations/phase, and explicit unsupported output types.
func TestV0CompletedUnsupportedMatrix(t *testing.T) {
	t.Run("malformed-and-truncated-tool-arguments", func(t *testing.T) {
		for _, tc := range []struct {
			name, status, args string
			target             error
		}{{"completed-malformed", "completed", `"{"`, ErrInvalidWireData}, {"incomplete-truncated", "incomplete", `"{"`, ErrUnsupported}} {
			t.Run(tc.name, func(t *testing.T) {
				details := ""
				if tc.status == "incomplete" {
					details = `,"incomplete_details":{"reason":"max_output_tokens"}`
				}
				body := []byte(`{"id":"r","model":"g","status":"` + tc.status + `"` + details + `,"output":[{"type":"function_call","id":"fc","call_id":"c","name":"f","arguments":` + tc.args + `,"status":"` + tc.status + `"}],"usage":{"input_tokens":0,"output_tokens":0}}`)
				_, err := ResponsesResponseToAnthropic(body, "a")
				requireErrorIs(t, err, tc.target)
			})
		}
	})

	t.Run("annotations-citations-and-phase", func(t *testing.T) {
		for _, annotations := range []string{`[]`} {
			if _, err := ResponsesResponseToAnthropic(responsesCompleted(`[{"type":"message","id":"m","role":"assistant","status":"completed","content":[{"type":"output_text","text":"x","annotations":`+annotations+`}]}]`, `{"input_tokens":0,"output_tokens":0}`), "a"); err != nil {
				t.Fatal(err)
			}
		}
		for _, tc := range []struct {
			name, body string
			fn         func([]byte, string) ([]byte, error)
		}{
			{"nonempty-annotations", string(responsesCompleted(`[{"type":"message","id":"m","role":"assistant","status":"completed","content":[{"type":"output_text","text":"x","annotations":[{}]}]}]`, `{"input_tokens":0,"output_tokens":0}`)), ResponsesResponseToAnthropic},
			{"responses-phase", string(responsesCompleted(`[{"type":"message","id":"m","role":"assistant","status":"completed","content":[],"phase":"final"}]`, `{"input_tokens":0,"output_tokens":0}`)), ResponsesResponseToAnthropic},
			{"nonempty-citations", string(anthropicCompleted(`[{"type":"text","text":"x","citations":[{}]}]`, `end_turn`, `{"input_tokens":0,"output_tokens":0}`)), AnthropicResponseToResponses},
			{"anthropic-phase", string(anthropicCompleted(`[{"type":"text","text":"x","phase":"final"}]`, `end_turn`, `{"input_tokens":0,"output_tokens":0}`)), AnthropicResponseToResponses},
		} {
			t.Run(tc.name, func(t *testing.T) { _, err := tc.fn([]byte(tc.body), ""); requireErrorIs(t, err, ErrUnsupported) })
		}
	})

	t.Run("known-unsupported-responses-output", func(t *testing.T) {
		for _, typ := range []string{"reasoning", "thinking", "redacted_thinking", "refusal", "web_search_call", "file_search_call", "code_interpreter_call", "image_generation_call", "computer_call", "mcp_call", "server_tool_use", "server_tool_result", "audio", "file", "document"} {
			t.Run(typ, func(t *testing.T) {
				_, err := ResponsesResponseToAnthropic(responsesCompleted(`[{"type":"`+typ+`"}]`, `{"input_tokens":0,"output_tokens":0}`), "a")
				requireErrorIs(t, err, ErrUnsupported)
			})
		}
	})
	t.Run("known-unsupported-anthropic-content", func(t *testing.T) {
		for _, typ := range []string{"reasoning", "thinking", "redacted_thinking", "refusal", "server_tool_use", "web_search_tool_result", "audio", "file", "document"} {
			t.Run(typ, func(t *testing.T) {
				_, err := AnthropicResponseToResponses(anthropicCompleted(`[{"type":"`+typ+`"}]`, `end_turn`, `{"input_tokens":0,"output_tokens":0}`), "g")
				requireErrorIs(t, err, ErrUnsupported)
			})
		}
	})
	t.Run("supported-prefix-never-partial-success", func(t *testing.T) {
		body := responsesCompleted(`[{"type":"message","id":"m","role":"assistant","status":"completed","content":[{"type":"output_text","text":"safe"}]},{"type":"web_search_call"}]`, `{"input_tokens":0,"output_tokens":1}`)
		got, err := ResponsesResponseToAnthropic(body, "a")
		requireErrorIs(t, err, ErrUnsupported)
		if got != nil {
			t.Fatalf("partial payload returned: %s", got)
		}
	})

}
