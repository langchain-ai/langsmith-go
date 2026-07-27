package messagetranslators

import (
	"fmt"
	"testing"
)

// Source parity: Python/JS request system, tool-choice, tools, and cache-control scenarios.
func TestV0RequestConfigurationMatrix(t *testing.T) {
	t.Run("system-string-and-block-concatenation", func(t *testing.T) {
		for _, tc := range []struct{ name, system, want string }{
			{"string", `"alpha"`, "alpha"},
			{"blocks-exact-no-separator", `[{"type":"text","text":"alpha"},{"type":"text","text":"β"},{"type":"text","text":""}]`, "alphaβ"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				got, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"system":`+tc.system+`,"messages":[{"role":"user","content":"x"}]}`), "r")
				if err != nil {
					t.Fatal(err)
				}
				if object(t, got)["instructions"] != tc.want {
					t.Fatalf("instructions = %#v", object(t, got)["instructions"])
				}
			})
		}
	})

	t.Run("system-invalid-and-cache", func(t *testing.T) {
		cases := []struct {
			name, value string
			target      error
		}{
			{"number", `3`, ErrInvalidWireData}, {"non-object-block", `["x"]`, ErrInvalidWireData},
			{"missing-text", `[{"type":"text"}]`, ErrInvalidWireData}, {"non-text", `[{"type":"image","text":"x"}]`, ErrUnsupported},
			{"cache-control", `[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}]`, nil},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"system":`+tc.value+`,"messages":[{"role":"user","content":"x"}]}`), "r")
				requireErrorIs(t, err, tc.target)
			})
		}
	})

	t.Run("tool-choice-both-directions", func(t *testing.T) {
		for _, tc := range []struct{ name, anth, responses string }{
			{"auto", `{"type":"auto"}`, `"auto"`}, {"any", `{"type":"any"}`, `"required"`},
			{"none-old-regression", `{"type":"none"}`, `"none"`}, {"named", `{"type":"tool","name":"f"}`, `{"type":"function","name":"f"}`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				ab, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tool_choice":`+tc.anth+`}`), "r")
				if err != nil {
					t.Fatal(err)
				}
				rb, err := ResponsesRequestToAnthropic([]byte(`{"model":"r","max_output_tokens":1,"input":"x","tool_choice":`+tc.responses+`}`), "a")
				if err != nil {
					t.Fatal(err)
				}
				if fmt.Sprint(object(t, ab)["tool_choice"]) == "<nil>" || fmt.Sprint(object(t, rb)["tool_choice"]) == "<nil>" {
					t.Fatal("choice dropped")
				}
			})
		}
	})

	t.Run("malformed-tool-choice", func(t *testing.T) {
		cases := []struct {
			name   string
			fn     func([]byte, string, ...Option) ([]byte, error)
			body   string
			target error
		}{
			{"anth-string", AnthropicRequestToResponses, `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tool_choice":"auto"}`, ErrInvalidWireData},
			{"anth-unknown", AnthropicRequestToResponses, `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tool_choice":{"type":"required"}}`, ErrUnsupported},
			{"anth-named-missing-name", AnthropicRequestToResponses, `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tool_choice":{"type":"tool"}}`, ErrInvalidWireData},
			{"responses-unknown-string", ResponsesRequestToAnthropic, `{"model":"r","max_output_tokens":1,"input":"x","tool_choice":"any"}`, ErrUnsupported},
			{"responses-bad-object", ResponsesRequestToAnthropic, `{"model":"r","max_output_tokens":1,"input":"x","tool_choice":{"type":"tool","name":"f"}}`, ErrUnsupported},
			{"responses-number", ResponsesRequestToAnthropic, `{"model":"r","max_output_tokens":1,"input":"x","tool_choice":1}`, ErrInvalidWireData},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) { _, err := tc.fn([]byte(tc.body), ""); requireErrorIs(t, err, tc.target) })
		}
	})

	t.Run("tools-defaults-and-validation", func(t *testing.T) {
		for _, tc := range []struct {
			name, body string
			fn         func([]byte, string, ...Option) ([]byte, error)
			target     error
		}{
			{"anth-default-schema", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tools":[{"name":"f"}]}`, AnthropicRequestToResponses, nil},
			{"responses-default-schema", `{"model":"r","max_output_tokens":1,"input":"x","tools":[{"type":"function","name":"f"}]}`, ResponsesRequestToAnthropic, nil},
			{"anth-description", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tools":[{"name":"f","description":3}]}`, AnthropicRequestToResponses, ErrInvalidWireData},
			{"anth-schema", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tools":[{"name":"f","input_schema":[]}]}`, AnthropicRequestToResponses, ErrInvalidWireData},
			{"anth-cache", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tools":[{"name":"f","cache_control":{"type":"ephemeral"}}]}`, AnthropicRequestToResponses, nil},
			{"anth-type", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tools":[{"name":"f","type":"custom"}]}`, AnthropicRequestToResponses, ErrUnsupported},
			{"anth-strict", `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"x"}],"tools":[{"name":"f","strict":false}]}`, AnthropicRequestToResponses, ErrUnsupported},
			{"responses-description", `{"model":"r","max_output_tokens":1,"input":"x","tools":[{"type":"function","name":"f","description":3}]}`, ResponsesRequestToAnthropic, ErrInvalidWireData},
			{"responses-schema", `{"model":"r","max_output_tokens":1,"input":"x","tools":[{"type":"function","name":"f","parameters":[]}]}`, ResponsesRequestToAnthropic, ErrInvalidWireData},
			{"responses-cache", `{"model":"r","max_output_tokens":1,"input":"x","tools":[{"type":"function","name":"f","cache_control":{"type":"ephemeral"}}]}`, ResponsesRequestToAnthropic, nil},
			{"responses-type", `{"model":"r","max_output_tokens":1,"input":"x","tools":[{"type":"web_search","name":"f"}]}`, ResponsesRequestToAnthropic, ErrUnsupported},
			{"responses-strict", `{"model":"r","max_output_tokens":1,"input":"x","tools":[{"type":"function","name":"f","strict":false}]}`, ResponsesRequestToAnthropic, ErrUnsupported},
		} {
			t.Run(tc.name, func(t *testing.T) {
				got, err := tc.fn([]byte(tc.body), "")
				if tc.target != nil {
					requireErrorIs(t, err, tc.target)
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				tools := list(t, object(t, got)["tools"])
				if len(tools) != 1 {
					t.Fatal("tool dropped")
				}
			})
		}
	})
}

// Source parity: Python/JS request history, tool pairing, Unicode, and cache metadata scenarios.
func TestV0RequestHistoryAndToolMatrix(t *testing.T) {
	t.Run("assistant-string-and-multiple-text", func(t *testing.T) {
		for _, content := range []string{`"hello"`, `[{"type":"text","text":"a"},{"type":"text","text":"b"}]`} {
			got, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":`+content+`}]}`), "r")
			if err != nil {
				t.Fatal(err)
			}
			items := list(t, object(t, got)["input"])
			assistant := mapv(t, items[1])
			for _, p := range list(t, assistant["content"]) {
				if mapv(t, p)["type"] != "output_text" {
					t.Fatalf("assistant part = %#v", p)
				}
			}
		}
	})

	t.Run("complete-resumed-history-two-calls-reordered-outputs", func(t *testing.T) {
		body := `{"model":"a","max_tokens":5,"messages":[` +
			`{"role":"user","content":"go"},` +
			`{"role":"assistant","content":[{"type":"text","text":"pre"},{"type":"tool_use","id":"c1","name":"f","input":{"雪":"☃"}},{"type":"tool_use","id":"c2","name":"g"}]},` +
			`{"role":"user","content":[{"type":"tool_result","tool_use_id":"c2","content":"two"},{"type":"tool_result","tool_use_id":"c1","content":[{"type":"text","text":"one"},{"type":"text","text":"!"}]}]},` +
			`{"role":"assistant","content":"resumed"}]}`
		got, err := AnthropicRequestToResponses([]byte(body), "r")
		if err != nil {
			t.Fatal(err)
		}
		items := list(t, object(t, got)["input"])
		if len(items) != 7 {
			t.Fatalf("items = %#v", items)
		}
		if mapv(t, items[2])["arguments"] != `{"雪":"☃"}` || mapv(t, items[5])["output"] != "one!" {
			t.Fatalf("history = %#v", items)
		}
		back, err := ResponsesRequestToAnthropic(got, "a")
		if err != nil {
			t.Fatal(err)
		}
		messages := list(t, object(t, back)["messages"])
		if len(messages) != 4 {
			t.Fatalf("back = %s", back)
		}
		assistantContent := list(t, mapv(t, messages[1])["content"])
		if mapv(t, mapv(t, assistantContent[1])["input"])["雪"] != "☃" {
			t.Fatalf("Responses-to-Anthropic Unicode arguments lost: %s", back)
		}
	})

	t.Run("tool-input-and-result-defaults", func(t *testing.T) {
		for _, tc := range []struct{ name, input, result, want string }{
			{"missing-and-missing", "", "", ""}, {"null-result", `,"input":{}`, `,"content":null`, ""},
			{"empty-array", `,"input":{}`, `,"content":[]`, ""}, {"false-is-error", `,"input":{}`, `,"content":"ok","is_error":false`, "ok"},
		} {
			t.Run(tc.name, func(t *testing.T) {
				body := `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"` + tc.input + `}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"c"` + tc.result + `}]}]}`
				got, err := AnthropicRequestToResponses([]byte(body), "r")
				if err != nil {
					t.Fatal(err)
				}
				items := list(t, object(t, got)["input"])
				if mapv(t, items[len(items)-1])["output"] != tc.want {
					t.Fatalf("output = %#v", items)
				}
			})
		}
	})

	t.Run("pairing-and-is-error-errors", func(t *testing.T) {
		base := `{"model":"a","max_tokens":1,"messages":%s}`
		cases := []struct {
			name, messages string
			target         error
		}{
			{"orphan", `[{"role":"user","content":[{"type":"tool_result","tool_use_id":"c"}]}]`, ErrInvalidSequence},
			{"unresolved", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"}]}]`, ErrInvalidSequence},
			{"duplicate", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"},{"type":"tool_use","id":"c","name":"g"}]}]`, ErrInvalidSequence},
			{"is-error-true", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"c","is_error":true}]}]`, ErrUnsupported},
			{"is-error-nonbool", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"c","is_error":"false"}]}]`, ErrInvalidWireData},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToResponses([]byte(fmt.Sprintf(base, tc.messages)), "r")
				requireErrorIs(t, err, tc.target)
			})
		}
	})

	t.Run("responses-pairing-errors", func(t *testing.T) {
		for _, tc := range []struct{ name, input string }{
			{"orphan", `[{"type":"message","role":"user","content":"u"},{"type":"function_call_output","call_id":"c","output":"x"}]`},
			{"unresolved", `[{"type":"message","role":"user","content":"u"},{"type":"function_call","call_id":"c","name":"f","arguments":"{}"}]`},
			{"duplicate", `[{"type":"message","role":"user","content":"u"},{"type":"function_call","call_id":"c","name":"f","arguments":"{}"},{"type":"function_call","call_id":"c","name":"g","arguments":"{}"}]`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := ResponsesRequestToAnthropic([]byte(`{"model":"r","max_output_tokens":1,"input":`+tc.input+`}`), "a")
				requireErrorIs(t, err, ErrInvalidSequence)
			})
		}
	})
}

// Source parity: Python/JS image modalities and known unsupported input/content boundaries.
func TestV0RequestMediaAndUnsupportedMatrix(t *testing.T) {
	t.Run("url-and-base64-roundtrip", func(t *testing.T) {
		body := []byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[{"type":"image","source":{"type":"url","url":"https://example/x.png"}},{"type":"image","source":{"type":"base64","media_type":"image/png","data":"4piD"}}]}]}`)
		got, err := AnthropicRequestToResponses(body, "r")
		if err != nil {
			t.Fatal(err)
		}
		back, err := ResponsesRequestToAnthropic(got, "a")
		if err != nil {
			t.Fatal(err)
		}
		content := list(t, mapv(t, list(t, object(t, back)["messages"])[0])["content"])
		if len(content) != 2 || mapv(t, mapv(t, content[0])["source"])["url"] != "https://example/x.png" || mapv(t, mapv(t, content[1])["source"])["data"] != "4piD" {
			t.Fatalf("roundtrip = %#v", content)
		}
	})

	t.Run("malformed-images-and-detail", func(t *testing.T) {
		anth := []struct {
			name, block string
			target      error
		}{
			{"empty-url", `{"type":"image","source":{"type":"url","url":""}}`, ErrInvalidWireData}, {"empty-media", `{"type":"image","source":{"type":"base64","media_type":"","data":"eA=="}}`, ErrInvalidWireData},
			{"empty-base64", `{"type":"image","source":{"type":"base64","media_type":"image/png","data":""}}`, ErrInvalidWireData}, {"bad-base64", `{"type":"image","source":{"type":"base64","media_type":"image/png","data":"!"}}`, ErrInvalidWireData},
		}
		for _, tc := range anth {
			t.Run("anth-"+tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[`+tc.block+`]}]}`), "r")
				requireErrorIs(t, err, tc.target)
			})
		}
		resp := []struct {
			name, url, extra string
			target           error
		}{
			{"empty", "", "", ErrInvalidWireData}, {"malformed-data", "data:image/png;base64", "", ErrUnsupported}, {"empty-media", "data:;base64,eA==", "", ErrInvalidWireData}, {"empty-data", "data:image/png;base64,", "", ErrInvalidWireData}, {"bad-data", "data:image/png;base64,!", "", ErrInvalidWireData}, {"detail", "https://x", `,"detail":"high"`, ErrUnsupported},
		}
		for _, tc := range resp {
			t.Run("responses-"+tc.name, func(t *testing.T) {
				_, err := ResponsesRequestToAnthropic([]byte(`{"model":"r","max_output_tokens":1,"input":[{"type":"message","role":"user","content":[{"type":"input_image","image_url":"`+tc.url+`"`+tc.extra+`}]}]}`), "a")
				requireErrorIs(t, err, tc.target)
			})
		}
	})

	t.Run("assistant-image", func(t *testing.T) {
		_, err := ResponsesRequestToAnthropic([]byte(`{"model":"r","max_output_tokens":1,"input":[{"type":"message","role":"user","content":"u"},{"type":"message","role":"assistant","content":[{"type":"input_image","image_url":"https://x"}]}]}`), "a")
		requireErrorIs(t, err, ErrInvalidWireData)
	})

	t.Run("unsupported-content-kinds", func(t *testing.T) {
		for _, typ := range []string{"document", "audio", "input_file", "file"} {
			t.Run("anth-"+typ, func(t *testing.T) {
				_, err := AnthropicRequestToResponses([]byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[{"type":"`+typ+`"}]}]}`), "r")
				requireErrorIs(t, err, ErrUnsupported)
			})
		}
		for _, typ := range []string{"input_audio", "input_file", "file_id", "document"} {
			t.Run("responses-"+typ, func(t *testing.T) {
				_, err := ResponsesRequestToAnthropic([]byte(`{"model":"r","max_output_tokens":1,"input":[{"type":"message","role":"user","content":[{"type":"`+typ+`"}]}]}`), "a")
				requireErrorIs(t, err, ErrUnsupported)
			})
		}
	})

	t.Run("missing-text-unknown-blocks-and-roles", func(t *testing.T) {
		cases := []struct {
			name   string
			fn     func([]byte, string, ...Option) ([]byte, error)
			body   string
			target error
		}{
			{"anth-missing-text", AnthropicRequestToResponses, `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[{"type":"text"}]}]}`, ErrInvalidWireData},
			{"responses-missing-text", ResponsesRequestToAnthropic, `{"model":"r","max_output_tokens":1,"input":[{"type":"message","role":"user","content":[{"type":"input_text"}]}]}`, ErrInvalidWireData},
			{"anth-unknown", AnthropicRequestToResponses, `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[{"type":"mystery"}]}]}`, ErrUnsupported},
			{"responses-unknown", ResponsesRequestToAnthropic, `{"model":"r","max_output_tokens":1,"input":[{"type":"mystery"}]}`, ErrUnsupported},
			{"anth-role", AnthropicRequestToResponses, `{"model":"a","max_tokens":1,"messages":[{"role":"system","content":"x"}]}`, ErrUnsupported},
			{"responses-role", ResponsesRequestToAnthropic, `{"model":"r","max_output_tokens":1,"input":[{"type":"message","role":"developer","content":"x"}]}`, ErrUnsupported},
		}
		for _, tc := range cases {
			t.Run(tc.name, func(t *testing.T) { _, err := tc.fn([]byte(tc.body), ""); requireErrorIs(t, err, tc.target) })
		}
	})

	t.Run("cache-controls", func(t *testing.T) {
		blocks := []struct{ name, content string }{
			{"normal-text", `[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}]`},
		}
		for _, tc := range blocks {
			t.Run(tc.name, func(t *testing.T) {
				body := `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":` + tc.content + `}]}`
				_, err := AnthropicRequestToResponses([]byte(body), "r")
				if err != nil {
					t.Fatal(err)
				}
			})
		}
		for _, tc := range []struct{ name, block string }{
			{"tool-result", `{"type":"tool_result","tool_use_id":"c","content":"x","cache_control":{"type":"ephemeral"}}`},
			{"nested-result-text", `{"type":"tool_result","tool_use_id":"c","content":[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}]}`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				body := `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f","cache_control":{"type":"ephemeral"}}]},{"role":"user","content":[` + tc.block + `]}]}`
				_, err := AnthropicRequestToResponses([]byte(body), "r")
				if err != nil {
					t.Fatal(err)
				}
			})
		}
	})

}
