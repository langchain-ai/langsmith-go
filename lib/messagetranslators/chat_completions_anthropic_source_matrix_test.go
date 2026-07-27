package messagetranslators

import (
	"fmt"
	"testing"
)

func TestAnthropicToChatCompletionsRequestValidationMatrix(t *testing.T) {
	t.Run("scalars-stop-stream-and-model", func(t *testing.T) {
		b, err := AnthropicRequestToChatCompletions([]byte(`{"model":"a","max_tokens":8,"temperature":0.2,"top_p":0.8,"stop_sequences":["A","B"],"stream":true,"messages":[{"role":"user","content":"u"}]}`), "g")
		if err != nil {
			t.Fatal(err)
		}
		r := object(t, b)
		if r["model"] != "g" || r["max_tokens"] != float64(8) || r["temperature"] != .2 || r["top_p"] != .8 || len(list(t, r["stop"])) != 2 || mapv(t, r["stream_options"])["include_usage"] != true {
			t.Fatalf("result=%s", b)
		}
	})
	t.Run("required-and-range-errors", func(t *testing.T) {
		for _, tc := range []struct{ name, body string }{
			{"model", `{"max_tokens":1,"messages":[{"role":"user","content":"u"}]}`},
			{"max", `{"model":"a","messages":[{"role":"user","content":"u"}]}`},
			{"messages-missing", `{"model":"a","max_tokens":1}`},
			{"messages-empty", `{"model":"a","max_tokens":1,"messages":[]}`},
			{"temperature", `{"model":"a","max_tokens":1,"temperature":2,"messages":[{"role":"user","content":"u"}]}`},
			{"top-p", `{"model":"a","max_tokens":1,"top_p":-1,"messages":[{"role":"user","content":"u"}]}`},
			{"stream-type", `{"model":"a","max_tokens":1,"stream":"true","messages":[{"role":"user","content":"u"}]}`},
			{"stop-type", `{"model":"a","max_tokens":1,"stop_sequences":"x","messages":[{"role":"user","content":"u"}]}`},
			{"stop-empty", `{"model":"a","max_tokens":1,"stop_sequences":[""],"messages":[{"role":"user","content":"u"}]}`},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToChatCompletions([]byte(tc.body), "")
				requireErrorIs(t, err, ErrInvalidWireData)
			})
		}
	})
	t.Run("system-errors", func(t *testing.T) {
		for _, tc := range []struct {
			name, system string
			target       error
		}{
			{"number", `1`, ErrInvalidWireData}, {"nonobject", `["x"]`, ErrInvalidWireData},
			{"missing-text", `[{"type":"text"}]`, ErrInvalidWireData}, {"wrong-type", `[{"type":"image","text":"x"}]`, ErrUnsupported},
			{"cache", `[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}]`, nil},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToChatCompletions([]byte(`{"model":"a","max_tokens":1,"system":`+tc.system+`,"messages":[{"role":"user","content":"u"}]}`), "g")
				requireErrorIs(t, err, tc.target)
			})
		}
	})
	t.Run("images-malformed-sources", func(t *testing.T) {
		for _, tc := range []struct {
			name, image string
			target      error
		}{
			{"source-missing", `{"type":"image"}`, ErrInvalidWireData},
			{"source-type", `{"type":"image","source":{"type":"file","file_id":"x"}}`, ErrUnsupported},
			{"url-empty", `{"type":"image","source":{"type":"url","url":""}}`, ErrInvalidWireData},
			{"base64-media", `{"type":"image","source":{"type":"base64","media_type":"","data":"eA=="}}`, ErrInvalidWireData},
			{"base64-data", `{"type":"image","source":{"type":"base64","media_type":"image/png","data":""}}`, ErrInvalidWireData},
			{"base64-invalid", `{"type":"image","source":{"type":"base64","media_type":"image/png","data":"!"}}`, ErrInvalidWireData},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToChatCompletions([]byte(`{"model":"a","max_tokens":1,"messages":[{"role":"user","content":[`+tc.image+`]}]}`), "g")
				requireErrorIs(t, err, tc.target)
			})
		}
	})
	t.Run("tool-result-content-and-sequence", func(t *testing.T) {
		prefix := `{"model":"a","max_tokens":1,"messages":[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f","input":{}}]},{"role":"user","content":[%s]}]}`
		for _, tc := range []struct {
			name, result string
			target       error
		}{
			{"missing-content", `{"type":"tool_result","tool_use_id":"c"}`, nil},
			{"null-content", `{"type":"tool_result","tool_use_id":"c","content":null}`, nil},
			{"empty-array", `{"type":"tool_result","tool_use_id":"c","content":[]}`, nil},
			{"text-array", `{"type":"tool_result","tool_use_id":"c","content":[{"type":"text","text":"a"},{"type":"text","text":"β"}]}`, nil},
			{"content-number", `{"type":"tool_result","tool_use_id":"c","content":1}`, ErrInvalidWireData},
			{"content-image", `{"type":"tool_result","tool_use_id":"c","content":[{"type":"image"}]}`, ErrUnsupported},
			{"nested-cache", `{"type":"tool_result","tool_use_id":"c","content":[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}]}`, nil},
			{"is-error-type", `{"type":"tool_result","tool_use_id":"c","is_error":"false"}`, ErrInvalidWireData},
			{"is-error", `{"type":"tool_result","tool_use_id":"c","is_error":true}`, ErrUnsupported},
		} {
			t.Run(tc.name, func(t *testing.T) {
				b, err := AnthropicRequestToChatCompletions([]byte(fmt.Sprintf(prefix, tc.result)), "g")
				if tc.target != nil {
					requireErrorIs(t, err, tc.target)
					return
				}
				if err != nil {
					t.Fatal(err)
				}
				ms := list(t, object(t, b)["messages"])
				if mapv(t, ms[len(ms)-1])["role"] != "tool" {
					t.Fatalf("result=%s", b)
				}
			})
		}
	})
	t.Run("pairing-and-message-errors", func(t *testing.T) {
		for _, tc := range []struct {
			name, messages string
			target         error
		}{
			{"orphan", `[{"role":"user","content":[{"type":"tool_result","tool_use_id":"c"}]}]`, ErrInvalidSequence},
			{"duplicate-call", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"},{"type":"tool_use","id":"c","name":"g"}]}]`, ErrInvalidSequence},
			{"unresolved", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"}]}]`, ErrInvalidSequence},
			{"assistant-before-result", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"}]},{"role":"assistant","content":"late"}]`, ErrInvalidSequence},
			{"user-text-before-result", `[{"role":"user","content":"u"},{"role":"assistant","content":[{"type":"tool_use","id":"c","name":"f"}]},{"role":"user","content":[{"type":"text","text":"late"}]}]`, ErrUnsupported},
			{"bad-role", `[{"role":"system","content":"u"}]`, ErrUnsupported},
			{"name", `[{"role":"user","name":"x","content":"u"}]`, ErrUnsupported},
			{"empty-content", `[{"role":"user","content":[]}]`, ErrInvalidWireData},
		} {
			t.Run(tc.name, func(t *testing.T) {
				_, err := AnthropicRequestToChatCompletions([]byte(`{"model":"a","max_tokens":1,"messages":`+tc.messages+`}`), "g")
				requireErrorIs(t, err, tc.target)
			})
		}
	})
	t.Run("explicit-top-level-unsupported", func(t *testing.T) {
		for _, field := range []string{"thinking", "top_k", "service_tier", "output_config", "mcp_servers", "context_management"} {
			t.Run(field, func(t *testing.T) {
				_, err := AnthropicRequestToChatCompletions([]byte(`{"model":"a","max_tokens":1,"`+field+`":{},"messages":[{"role":"user","content":"u"}]}`), "g")
				requireErrorIs(t, err, ErrUnsupported)
			})
		}
	})
	t.Run("top-level-cache-control-is-ignored", func(t *testing.T) {
		_, err := AnthropicRequestToChatCompletions([]byte(`{"model":"a","max_tokens":1,"cache_control":{"type":"ephemeral"},"system":[{"type":"text","text":"s","cache_control":{"type":"ephemeral"}}],"tools":[{"name":"f","cache_control":{"type":"ephemeral"}}],"messages":[{"role":"user","content":"u"}]}`), "g")
		if err != nil {
			t.Fatal(err)
		}
	})
}

func TestAnthropicToChatCompletionsCompletedValidationMatrix(t *testing.T) {
	usage := `{"input_tokens":0,"output_tokens":0}`
	for _, tc := range []struct {
		name, content, stop string
		target              error
	}{
		{"text-citations-empty", `[{"type":"text","text":"x","citations":[]}]`, "end_turn", nil},
		{"text-citations", `[{"type":"text","text":"x","citations":[{}]}]`, "end_turn", ErrUnsupported},
		{"text-missing", `[{"type":"text"}]`, "end_turn", ErrInvalidWireData},
		{"tool-input-null", `[{"type":"tool_use","id":"c","name":"f","input":null}]`, "tool_use", nil},
		{"tool-input-array", `[{"type":"tool_use","id":"c","name":"f","input":[]}]`, "tool_use", ErrInvalidWireData},
		{"tool-missing-id", `[{"type":"tool_use","name":"f","input":{}}]`, "tool_use", ErrInvalidWireData},
		{"stop-sequence", `[{"type":"text","text":"x"}]`, "stop_sequence", nil},
		{"pause", `[]`, "pause_turn", ErrUnsupported},
		{"unknown-stop", `[]`, "wat", ErrInvalidWireData},
		{"cache-control", `[{"type":"text","text":"x","cache_control":{"type":"ephemeral"}}]`, "end_turn", nil},
		{"phase", `[{"type":"text","text":"x","phase":"final"}]`, "end_turn", ErrUnsupported},
	} {
		t.Run(tc.name, func(t *testing.T) {
			body := []byte(`{"id":"m","type":"message","role":"assistant","model":"a","content":` + tc.content + `,"stop_reason":"` + tc.stop + `","usage":` + usage + `}`)
			out, err := AnthropicResponseToChatCompletions(body, "g")
			if tc.target == nil {
				if err != nil {
					t.Fatal(err)
				}
				if out == nil {
					t.Fatal("missing output")
				}
				return
			}
			if out != nil {
				t.Fatalf("partial output=%s", out)
			}
			requireErrorIs(t, err, tc.target)
		})
	}
}
