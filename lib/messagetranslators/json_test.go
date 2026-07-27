package messagetranslators

import (
	"encoding/json"
	"errors"
	"fmt"
	"testing"
)

func object(t *testing.T, b []byte) map[string]any {
	t.Helper()
	var v map[string]any
	if err := json.Unmarshal(b, &v); err != nil {
		t.Fatal(err)
	}
	return v
}

func requireErrorIs(t *testing.T, err, target error) {
	t.Helper()
	if !errors.Is(err, target) {
		t.Fatalf("error = %v, want errors.Is(_, %v)", err, target)
	}
}

func list(t *testing.T, v any) []any {
	t.Helper()
	a, ok := v.([]any)
	if !ok {
		t.Fatalf("not an array: %#v", v)
	}
	return a
}
func mapv(t *testing.T, v any) map[string]any {
	t.Helper()
	m, ok := v.(map[string]any)
	if !ok {
		t.Fatalf("not an object: %#v", v)
	}
	return m
}

func TestAnthropicRequestToResponsesTextImagesTools(t *testing.T) {
	in := []byte(`{
	 "model":"old","system":[{"type":"text","text":"be useful"}],"max_tokens":123,"temperature":0.2,"top_p":0.9,"stream":true,
	 "tools":[{"name":"weather","description":"lookup","input_schema":{"type":"object","properties":{"city":{"type":"string"}}}}],
	 "tool_choice":{"type":"tool","name":"weather"},
	 "messages":[
	  {"role":"user","content":[{"type":"text","text":"look"},{"type":"image","source":{"type":"url","url":"https://x/i.png"}},{"type":"image","source":{"type":"base64","media_type":"image/png","data":"aGk="}}]},
	  {"role":"assistant","content":[{"type":"text","text":"calling"},{"type":"tool_use","id":"call_1","name":"weather","input":{"city":"Paris"}}]},
	  {"role":"user","content":[{"type":"tool_result","tool_use_id":"call_1","content":"sunny"}]}
	 ]}`)
	b, err := AnthropicRequestToResponses(in, "new")
	if err != nil {
		t.Fatal(err)
	}
	r := object(t, b)
	if r["model"] != "new" || r["instructions"] != "be useful" || r["max_output_tokens"] != float64(123) {
		t.Fatalf("bad scalar mapping: %#v", r)
	}
	tools := list(t, r["tools"])
	if mapv(t, tools[0])["parameters"] == nil {
		t.Fatal("parameters missing")
	}
	if mapv(t, r["tool_choice"])["name"] != "weather" {
		t.Fatal("specific tool choice lost")
	}
	items := list(t, r["input"])
	if len(items) != 4 {
		t.Fatalf("want split ordered items, got %d: %#v", len(items), items)
	}
	content := list(t, mapv(t, items[0])["content"])
	if mapv(t, content[1])["image_url"] != "https://x/i.png" {
		t.Fatal("URL image lost")
	}
	if mapv(t, content[2])["image_url"] != "data:image/png;base64,aGk=" {
		t.Fatal("base64 image lost")
	}
	call := mapv(t, items[2])
	if call["type"] != "function_call" || call["arguments"] != `{"city":"Paris"}` {
		t.Fatalf("bad call: %#v", call)
	}
	out := mapv(t, items[3])
	if out["type"] != "function_call_output" || out["output"] != "sunny" {
		t.Fatalf("bad output: %#v", out)
	}
}

func TestResponsesRequestToAnthropic(t *testing.T) {
	in := []byte(`{"model":"r","instructions":"rules","max_output_tokens":44,"tool_choice":"required","tools":[{"type":"function","name":"f","parameters":{"type":"object"}}],"input":[{"type":"message","role":"user","content":[{"type":"input_text","text":"x"},{"type":"input_image","image_url":"data:image/jpeg;base64,aGk="}]},{"type":"function_call","call_id":"c","name":"f","arguments":"{\"n\":2}"},{"type":"function_call_output","call_id":"c","output":"ok"}]}`)
	b, err := ResponsesRequestToAnthropic(in, "a")
	if err != nil {
		t.Fatal(err)
	}
	a := object(t, b)
	if a["model"] != "a" || a["system"] != "rules" || a["max_tokens"] != float64(44) {
		t.Fatalf("bad mapping: %#v", a)
	}
	if mapv(t, a["tool_choice"])["type"] != "any" {
		t.Fatal("required not mapped")
	}
	ms := list(t, a["messages"])
	if len(ms) != 3 {
		t.Fatalf("got %d messages", len(ms))
	}
	image := mapv(t, list(t, mapv(t, ms[0])["content"])[1])
	source := mapv(t, image["source"])
	if source["type"] != "base64" || source["media_type"] != "image/jpeg" {
		t.Fatalf("bad image: %#v", image)
	}
	tool := mapv(t, list(t, mapv(t, ms[1])["content"])[0])
	if mapv(t, tool["input"])["n"] != float64(2) {
		t.Fatalf("bad arguments: %#v", tool)
	}
}

func TestResponseConversionsUsageStatusAndOrder(t *testing.T) {
	r := []byte(`{"id":"resp_1","model":"rm","status":"incomplete","incomplete_details":{"reason":"max_output_tokens"},"output":[{"type":"message","id":"m","role":"assistant","status":"incomplete","content":[{"type":"output_text","text":"hello"}]},{"type":"function_call","id":"fc","call_id":"c1","name":"f","arguments":"{\"x\":true}","status":"incomplete"}],"usage":{"input_tokens":10,"output_tokens":4,"input_tokens_details":{"cached_tokens":3},"output_tokens_details":{"reasoning_tokens":2}}}`)
	ab, err := ResponsesResponseToAnthropic(r, "am")
	if err != nil {
		t.Fatal(err)
	}
	a := object(t, ab)
	if a["stop_reason"] != "max_tokens" || a["model"] != "am" {
		t.Fatalf("bad terminal: %#v", a)
	}
	cs := list(t, a["content"])
	if mapv(t, cs[0])["type"] != "text" || mapv(t, cs[1])["type"] != "tool_use" {
		t.Fatal("content order lost")
	}
	if mapv(t, a["usage"])["cache_read_input_tokens"] != float64(3) {
		t.Fatal("cached usage lost")
	}
	back, err := AnthropicResponseToResponses(ab, "back")
	if err != nil {
		t.Fatal(err)
	}
	rr := object(t, back)
	if rr["status"] != "incomplete" || mapv(t, rr["incomplete_details"])["reason"] != "max_output_tokens" {
		t.Fatalf("bad status: %#v", rr)
	}
	outs := list(t, rr["output"])
	if len(outs) != 2 || mapv(t, outs[1])["call_id"] != "c1" {
		t.Fatal("outputs lost")
	}
}

func TestRequiredRequestValidation(t *testing.T) {
	validAnthropic := `{"model":"claude","max_tokens":1,"messages":[{"role":"user","content":"hi"}]}`
	validResponses := `{"model":"gpt","max_output_tokens":1,"input":"hi"}`
	if _, err := AnthropicRequestToResponses([]byte(validAnthropic), ""); err != nil {
		t.Fatal(err)
	}
	if _, err := ResponsesRequestToAnthropic([]byte(validResponses), ""); err != nil {
		t.Fatal(err)
	}

	cases := []struct {
		name string
		fn   func([]byte, string) ([]byte, error)
		body string
	}{
		{"anthropic-model", AnthropicRequestToResponses, `{"max_tokens":1,"messages":[{"role":"user","content":"x"}]}`},
		{"anthropic-max-zero", AnthropicRequestToResponses, `{"model":"m","max_tokens":0,"messages":[{"role":"user","content":"x"}]}`},
		{"anthropic-max-fraction", AnthropicRequestToResponses, `{"model":"m","max_tokens":1.5,"messages":[{"role":"user","content":"x"}]}`},
		{"anthropic-empty-messages", AnthropicRequestToResponses, `{"model":"m","max_tokens":1,"messages":[]}`},
		{"responses-model", ResponsesRequestToAnthropic, `{"max_output_tokens":1,"input":"x"}`},
		{"responses-max-string", ResponsesRequestToAnthropic, `{"model":"m","max_output_tokens":"1","input":"x"}`},
		{"responses-empty-input", ResponsesRequestToAnthropic, `{"model":"m","max_output_tokens":1,"input":[]}`},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			if _, err := tc.fn([]byte(tc.body), ""); !errors.Is(err, ErrInvalidWireData) {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestUsageArithmeticAndValidation(t *testing.T) {
	body := []byte(`{"id":"msg_source","type":"message","role":"assistant","model":"claude","content":[{"type":"text","text":"ok"}],"stop_reason":"end_turn","usage":{"input_tokens":5,"cache_read_input_tokens":3,"cache_creation_input_tokens":2,"output_tokens":4}}`)
	converted, err := AnthropicResponseToResponses(body, "gpt")
	if err != nil {
		t.Fatal(err)
	}
	r := object(t, converted)
	u := mapv(t, r["usage"])
	if u["input_tokens"] != float64(10) || u["output_tokens"] != float64(4) || u["total_tokens"] != float64(14) {
		t.Fatalf("bad usage: %#v", u)
	}
	if mapv(t, u["input_tokens_details"])["cached_tokens"] != float64(3) {
		t.Fatalf("bad cached usage: %#v", u)
	}

	openAI := []byte(`{"id":"r","model":"gpt","status":"completed","output":[],"usage":{"input_tokens":10,"output_tokens":4,"input_tokens_details":{"cached_tokens":3}}}`)
	converted, err = ResponsesResponseToAnthropic(openAI, "claude")
	if err != nil {
		t.Fatal(err)
	}
	au := mapv(t, object(t, converted)["usage"])
	if au["input_tokens"] != float64(7) || au["cache_read_input_tokens"] != float64(3) {
		t.Fatalf("bad subset subtraction: %#v", au)
	}

	bad := []byte(`{"status":"completed","output":[],"usage":{"input_tokens":2,"output_tokens":0,"input_tokens_details":{"cached_tokens":3}}}`)
	if _, err := ResponsesResponseToAnthropic(bad, "claude"); !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("got %v", err)
	}
}

func TestTerminalStatusStopSemanticsAndIDs(t *testing.T) {
	for _, status := range []string{"", "in_progress", "mystery", "failed"} {
		body := []byte(`{"status":"` + status + `","output":[]}`)
		if _, err := ResponsesResponseToAnthropic(body, "m"); err == nil {
			t.Fatalf("accepted status %q", status)
		}
	}
	base := `{"id":"msg_same","type":"message","role":"assistant","model":"m","content":[],"usage":{"input_tokens":0,"output_tokens":0},"stop_reason":"%s"}`
	for _, stop := range []string{"pause_turn", "refusal", "mystery", ""} {
		if _, err := AnthropicResponseToResponses([]byte(fmt.Sprintf(base, stop)), "r"); err == nil {
			t.Fatalf("accepted stop %q", stop)
		}
	}
	b1, err := AnthropicResponseToResponses([]byte(fmt.Sprintf(base, "end_turn")), "r")
	if err != nil {
		t.Fatal(err)
	}
	b2, err := AnthropicResponseToResponses([]byte(fmt.Sprintf(base, "end_turn")), "r")
	if err != nil {
		t.Fatal(err)
	}
	r1, r2 := object(t, b1), object(t, b2)
	if r1["id"] == "msg_same" || r1["id"] != r2["id"] {
		t.Fatalf("IDs must be destination-specific and deterministic: %#v %#v", r1["id"], r2["id"])
	}
}

func TestSemanticRejectionsAndRoleOrder(t *testing.T) {
	anthropicBase := `{"model":"m","max_tokens":1,"messages":[{"role":"user","content":"x"}],%s}`
	for _, field := range []string{`"output_config":{"format":{"type":"json_schema"}}`, `"tools":[{"name":"f","strict":true}]`, `"tools":[{"type":"web_search_20250305","name":"search"}]`} {
		if _, err := AnthropicRequestToResponses([]byte(fmt.Sprintf(anthropicBase, field)), ""); !errors.Is(err, ErrUnsupported) {
			t.Fatalf("field %s: %v", field, err)
		}
	}
	badRole := `{"model":"m","max_tokens":1,"messages":[{"role":"user","content":[{"type":"tool_use","id":"c","name":"f","input":{}}]}]}`
	if _, err := AnthropicRequestToResponses([]byte(badRole), ""); !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("got %v", err)
	}
	isError := `{"model":"m","max_tokens":1,"messages":[{"role":"user","content":[{"type":"tool_result","tool_use_id":"c","content":"bad","is_error":true}]}]}`
	if _, err := AnthropicRequestToResponses([]byte(isError), ""); !errors.Is(err, ErrUnsupported) {
		t.Fatalf("got %v", err)
	}
	for _, field := range []string{`"parallel_tool_calls":false`, `"text":{"format":{"type":"json_schema"}}`, `"tools":[{"type":"function","name":"f","strict":true}]`} {
		body := []byte(`{"model":"m","max_output_tokens":1,"input":"x",` + field + `}`)
		if _, err := ResponsesRequestToAnthropic(body, ""); !errors.Is(err, ErrUnsupported) {
			t.Fatalf("field %s: %v", field, err)
		}
	}
	orphan := `{"model":"m","max_output_tokens":1,"input":[{"type":"function_call_output","call_id":"c","output":"x"}]}`
	if _, err := ResponsesRequestToAnthropic([]byte(orphan), ""); !errors.Is(err, ErrInvalidSequence) {
		t.Fatalf("got %v", err)
	}
}

func TestMalformedAndUnsupported(t *testing.T) {
	cases := []struct {
		name   string
		fn     func([]byte, string) ([]byte, error)
		body   string
		target error
	}{
		{"top-level", AnthropicRequestToResponses, `[]`, ErrInvalidWireData},
		{"trailing", AnthropicRequestToResponses, `{} xyz`, ErrInvalidWireData},
		{"bad-arguments", ResponsesRequestToAnthropic, `{"input":[{"type":"function_call","call_id":"c","name":"f","arguments":"{"}]}`, ErrInvalidWireData},
		{"arguments-not-object", ResponsesResponseToAnthropic, `{"output":[{"type":"function_call","call_id":"c","name":"f","arguments":"[]"}]}`, ErrInvalidWireData},
		{"thinking", AnthropicRequestToResponses, `{"thinking":{"type":"enabled"},"messages":[]}`, ErrUnsupported},
		{"nontext-tool-output", AnthropicRequestToResponses, `{"model":"m","max_tokens":10,"messages":[{"role":"user","content":[{"type":"tool_result","tool_use_id":"c","content":[{"type":"image"}]}]}]}`, ErrUnsupported},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := tc.fn([]byte(tc.body), "")
			if !errors.Is(err, tc.target) {
				t.Fatalf("error %v, want %v", err, tc.target)
			}
			var ce *ConversionError
			if err != nil && !errors.As(err, &ce) && tc.name != "top-level" && tc.name != "trailing" {
				t.Fatalf("wanted path-aware error, got %T", err)
			}
		})
	}
}

func TestAssistantHistoryUsesOutputTextAndRejectsImages(t *testing.T) {
	body := []byte(`{"model":"claude","max_tokens":10,"messages":[{"role":"user","content":"hi"},{"role":"assistant","content":[{"type":"text","text":"hello"}]}]}`)
	converted, err := AnthropicRequestToResponses(body, "gpt")
	if err != nil {
		t.Fatal(err)
	}
	items := list(t, object(t, converted)["input"])
	assistantContent := list(t, mapv(t, items[1])["content"])
	if mapv(t, assistantContent[0])["type"] != "output_text" {
		t.Fatalf("assistant history must use output_text: %#v", assistantContent)
	}

	bad := []byte(`{"model":"claude","max_tokens":10,"messages":[{"role":"assistant","content":[{"type":"image","source":{"type":"url","url":"https://example.com/a.png"}}]}]}`)
	if _, err := AnthropicRequestToResponses(bad, "gpt"); !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("assistant image: %v", err)
	}
}

func TestResponsesResponseRequiresUsage(t *testing.T) {
	body := []byte(`{"id":"resp_1","model":"gpt","status":"completed","output":[]}`)
	if _, err := ResponsesResponseToAnthropic(body, "claude"); !errors.Is(err, ErrInvalidWireData) {
		t.Fatalf("missing usage: %v", err)
	}
}

func TestSynthesizedResponsesEnvelopeHasCoreSchemaFields(t *testing.T) {
	body := []byte(`{"id":"msg_1","type":"message","role":"assistant","model":"claude","content":[{"type":"text","text":"ok"}],"stop_reason":"end_turn","usage":{"input_tokens":1,"output_tokens":1}}`)
	converted, err := AnthropicResponseToResponses(body, "gpt")
	if err != nil {
		t.Fatal(err)
	}
	r := object(t, converted)
	if r["created_at"] == nil || r["completed_at"] == nil || r["error"] != nil || r["incomplete_details"] != nil {
		t.Fatalf("incomplete Responses envelope: %#v", r)
	}
	part := mapv(t, list(t, mapv(t, list(t, r["output"])[0])["content"])[0])
	if _, ok := part["logprobs"].([]any); !ok {
		t.Fatalf("output_text.logprobs missing: %#v", part)
	}
	u := mapv(t, r["usage"])
	if mapv(t, u["output_tokens_details"])["reasoning_tokens"] != float64(0) {
		t.Fatalf("output token details missing: %#v", u)
	}
}
