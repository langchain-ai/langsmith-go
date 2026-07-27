package messagetranslators

import (
	"errors"
	"testing"
)

func TestChatCompletionsRequestConversions(t *testing.T) {
	anthropic := []byte(`{"model":"claude","system":"rules","max_tokens":20,"stream":true,"tools":[{"name":"f","input_schema":{"type":"object"}}],"tool_choice":{"type":"tool","name":"f"},"messages":[{"role":"user","content":[{"type":"text","text":"look"},{"type":"image","source":{"type":"url","url":"https://example/x.png"}}]},{"role":"assistant","content":[{"type":"text","text":"calling"},{"type":"tool_use","id":"c1","name":"f","input":{"x":"☃"}}]},{"role":"user","content":[{"type":"tool_result","tool_use_id":"c1","content":"ok"}]}]}`)
	b, err := AnthropicRequestToChatCompletions(anthropic, "gpt")
	if err != nil {
		t.Fatal(err)
	}
	r := object(t, b)
	if r["model"] != "gpt" || r["max_tokens"] != float64(20) {
		t.Fatalf("bad scalars: %#v", r)
	}
	if mapv(t, r["stream_options"])["include_usage"] != true {
		t.Fatal("stream usage was not requested")
	}
	messages := list(t, r["messages"])
	if len(messages) != 4 || mapv(t, messages[0])["role"] != "system" || mapv(t, messages[3])["role"] != "tool" {
		t.Fatalf("bad messages: %#v", messages)
	}
	call := mapv(t, list(t, mapv(t, messages[2])["tool_calls"])[0])
	if mapv(t, call["function"])["arguments"] != `{"x":"☃"}` {
		t.Fatalf("bad arguments: %#v", call)
	}

	back, err := ChatCompletionsRequestToAnthropic(b, "claude-new")
	if err != nil {
		t.Fatal(err)
	}
	a := object(t, back)
	if a["model"] != "claude-new" || a["max_tokens"] != float64(20) {
		t.Fatalf("bad reverse mapping: %#v", a)
	}
	if mapv(t, a["tool_choice"])["name"] != "f" {
		t.Fatal("named tool choice lost")
	}
}

func TestChatCompletionsCompletedResponseConversions(t *testing.T) {
	chatCompletions := []byte(`{"id":"chatcmpl-1","object":"chat.completion","created":1,"model":"gpt","choices":[{"index":0,"message":{"role":"assistant","content":"hello","tool_calls":[{"id":"c1","type":"function","function":{"name":"f","arguments":"{\"emoji\":\"☃\"}"}}]},"finish_reason":"tool_calls"}],"usage":{"prompt_tokens":10,"completion_tokens":3,"total_tokens":13,"prompt_tokens_details":{"cached_tokens":4}}}`)
	b, err := ChatCompletionsResponseToAnthropic(chatCompletions, "claude")
	if err != nil {
		t.Fatal(err)
	}
	a := object(t, b)
	if a["stop_reason"] != "tool_use" || mapv(t, a["usage"])["input_tokens"] != float64(6) || mapv(t, a["usage"])["cache_read_input_tokens"] != float64(4) {
		t.Fatalf("bad mapping: %#v", a)
	}
	back, err := AnthropicResponseToChatCompletions(b, "gpt-new")
	if err != nil {
		t.Fatal(err)
	}
	r := object(t, back)
	if mapv(t, list(t, r["choices"])[0])["finish_reason"] != "tool_calls" || mapv(t, r["usage"])["prompt_tokens"] != float64(10) {
		t.Fatalf("bad reverse mapping: %#v", r)
	}
}

func TestChatCompletionsRequestUnsupportedAndSequence(t *testing.T) {
	cases := []struct {
		name, body string
		sentinel   error
	}{
		{"legacy", `{"model":"g","max_tokens":1,"functions":[],"messages":[{"role":"user","content":"x"}]}`, ErrUnsupported},
		{"structured", `{"model":"g","max_tokens":1,"response_format":{"type":"json_object"},"messages":[{"role":"user","content":"x"}]}`, ErrUnsupported},
		{"orphan", `{"model":"g","max_tokens":1,"messages":[{"role":"user","content":"x"},{"role":"tool","tool_call_id":"missing","content":"x"}]}`, ErrInvalidSequence},
		{"parallel-control", `{"model":"g","max_tokens":1,"parallel_tool_calls":false,"messages":[{"role":"user","content":"x"}]}`, ErrUnsupported},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			_, err := ChatCompletionsRequestToAnthropic([]byte(tc.body), "")
			if !errors.Is(err, tc.sentinel) {
				t.Fatalf("got %v", err)
			}
		})
	}
}

func TestChatCompletionsToAnthropicStreamFragmentedTools(t *testing.T) {
	s := NewChatCompletionsToAnthropicStream("claude")
	chunks := []string{
		`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{"role":"assistant"},"finish_reason":null}]}`,
		`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{"content":"hi"},"finish_reason":null}]}`,
		`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{"tool_calls":[{"index":0,"id":"a","type":"function","function":{"name":"f","arguments":"{\"x\":"}},{"index":1,"id":"b","type":"function","function":{"name":"g","arguments":"{\"y\":"}}]},"finish_reason":null}]}`,
		`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{"tool_calls":[{"index":1,"function":{"arguments":"\"☃\"}"}},{"index":0,"function":{"arguments":"1}"}}]},"finish_reason":null}]}`,
		`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[{"index":0,"delta":{},"finish_reason":"tool_calls"}]}`,
		`{"id":"c","object":"chat.completion.chunk","created":1,"model":"g","choices":[],"usage":{"prompt_tokens":9,"completion_tokens":4,"total_tokens":13,"prompt_tokens_details":{"cached_tokens":2}}}`,
	}
	var out []SSEEvent
	for _, raw := range chunks {
		got, err := s.Convert(SSEEvent{Data: []byte(raw)})
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, got...)
	}
	got, err := s.Convert(SSEEvent{Data: []byte("[DONE]")})
	if err != nil {
		t.Fatal(err)
	}
	out = append(out, got...)
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	if len(out) < 10 || out[0].Event != "message_start" || out[len(out)-1].Event != "message_stop" {
		t.Fatalf("bad output lifecycle: %#v", out)
	}
}

func TestAnthropicToChatCompletionsStream(t *testing.T) {
	s := NewAnthropicToChatCompletionsStream("gpt")
	events := []SSEEvent{
		{Event: "message_start", Data: []byte(`{"type":"message_start","message":{"id":"m","type":"message","role":"assistant","model":"claude","content":[],"stop_reason":null,"stop_sequence":null,"usage":{"input_tokens":5,"cache_read_input_tokens":2,"output_tokens":0}}}`)},
		{Event: "content_block_start", Data: []byte(`{"type":"content_block_start","index":0,"content_block":{"type":"tool_use","id":"c","name":"f","input":{}}}`)},
		{Event: "content_block_delta", Data: []byte(`{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"{\"x\":"}}`)},
		{Event: "content_block_delta", Data: []byte(`{"type":"content_block_delta","index":0,"delta":{"type":"input_json_delta","partial_json":"1}"}}`)},
		{Event: "content_block_stop", Data: []byte(`{"type":"content_block_stop","index":0}`)},
		{Event: "message_delta", Data: []byte(`{"type":"message_delta","delta":{"stop_reason":"tool_use","stop_sequence":null},"usage":{"output_tokens":3}}`)},
		{Event: "message_stop", Data: []byte(`{"type":"message_stop"}`)},
	}
	var out []SSEEvent
	for _, e := range events {
		got, err := s.Convert(e)
		if err != nil {
			t.Fatal(err)
		}
		out = append(out, got...)
	}
	if err := s.Finish(); err != nil {
		t.Fatal(err)
	}
	if string(out[len(out)-1].Data) != "[DONE]" {
		t.Fatalf("missing done: %#v", out)
	}
	u := object(t, out[len(out)-2].Data)
	if mapv(t, u["usage"])["prompt_tokens"] != float64(7) {
		t.Fatalf("bad usage: %#v", u)
	}
}
