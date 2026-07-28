package main

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"reflect"
	"strings"
	"testing"
)

func TestParseSSE(t *testing.T) {
	input := ": comment\r\nevent: message_start\r\ndata: {\"type\":\r\ndata: \"start\"}\r\n\r\ndata: [DONE]\n\nevent: final\ndata: {\"ok\":true}" // no final newline
	events, err := parseSSE(strings.NewReader(input))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 3 {
		t.Fatalf("got %d events: %#v", len(events), events)
	}
	if events[0].Event != "message_start" {
		t.Fatalf("event name = %q", events[0].Event)
	}
	wantFirst := map[string]any{"type": "start"}
	if !reflect.DeepEqual(events[0].Data, wantFirst) {
		t.Fatalf("first data = %#v, want %#v", events[0].Data, wantFirst)
	}
	if events[1].Event != "" || events[1].Data != "[DONE]" {
		t.Fatalf("DONE event = %#v", events[1])
	}
	if got := events[2].Data.(map[string]any)["ok"]; got != true {
		t.Fatalf("final event data = %#v", events[2].Data)
	}
}

func TestParseSSEIgnoresEmptyAndUnknownFields(t *testing.T) {
	events, err := parseSSE(strings.NewReader("event: lonely\n\nid: 7\ndata\n\n"))
	if err != nil {
		t.Fatal(err)
	}
	if len(events) != 1 || events[0].Data != "" || events[0].Event != "" {
		t.Fatalf("events = %#v", events)
	}
}

func TestBuildRequests(t *testing.T) {
	o := options{anthropicModel: "anthropic-test", openAIModel: "openai-test"}
	chatCompletions := buildRequest(apiChatCompletions, "tool", "stream", o)
	streamOptions, ok := chatCompletions["stream_options"].(map[string]any)
	if !ok || streamOptions["include_usage"] != true {
		t.Fatalf("Chat Completions stream_options = %#v", chatCompletions["stream_options"])
	}
	if chatCompletions["temperature"] != 0 || chatCompletions["stream"] != true {
		t.Fatalf("Chat Completions controls = %#v", chatCompletions)
	}
	encoded, err := json.Marshal(chatCompletions)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), "São Paulo") || strings.Contains(string(encoded), "API_KEY") {
		t.Fatalf("unexpected encoded request: %s", encoded)
	}

	responses := buildRequest("responses", "tool", "completed", o)
	choice := responses["tool_choice"].(map[string]any)
	if choice["type"] != "function" || choice["name"] != "get_weather" {
		t.Fatalf("Responses tool choice = %#v", choice)
	}
	anthropic := buildRequest("messages", "tool", "completed", o)
	if anthropic["tool_choice"].(map[string]any)["type"] != "tool" {
		t.Fatalf("Anthropic tool choice = %#v", anthropic["tool_choice"])
	}
}

func TestChatCompletionsSelector(t *testing.T) {
	o, err := parseFlags([]string{"-api", apiChatCompletionsSelector})
	if err != nil {
		t.Fatal(err)
	}
	if o.api != apiChatCompletions {
		t.Fatalf("normalized API = %q", o.api)
	}
	if _, err := parseFlags([]string{"-api", "chat"}); err == nil {
		t.Fatal("ambiguous legacy chat selector was accepted")
	}
}

func TestFetchSendsButDoesNotCaptureSecret(t *testing.T) {
	const secret = "test-secret-never-snapshot"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("Authorization"); got != "Bearer "+secret {
			t.Errorf("Authorization = %q", got)
		}
		if r.URL.Path != "/v1/chat/completions" {
			t.Errorf("path = %q", r.URL.Path)
		}
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"response-id","choices":[]}`))
	}))
	defer server.Close()

	o := options{openAIBase: server.URL, openAIKey: secret}
	c := capture{provider: "openai", api: apiChatCompletions, mode: "completed", scenario: "text", request: map[string]any{"model": "test"}}
	s, err := fetch(context.Background(), server.Client(), o, c)
	if err != nil {
		t.Fatal(err)
	}
	if s.API != apiChatCompletions {
		t.Fatalf("snapshot API = %q", s.API)
	}
	if got := filename(c); got != "openai_chat_completions_completed_text.json" {
		t.Fatalf("filename = %q", got)
	}
	encoded, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(encoded), secret) || !strings.Contains(string(encoded), "response-id") {
		t.Fatalf("snapshot leaked secret or lost response: %s", encoded)
	}
}

func TestFetchCapturesAnthropicAPIVersion(t *testing.T) {
	const version = "2023-06-01"
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if got := r.Header.Get("anthropic-version"); got != version {
			t.Errorf("anthropic-version = %q", got)
		}
		if got := r.Header.Get("x-api-key"); got != "test-key" {
			t.Errorf("x-api-key = %q", got)
		}
		_, _ = w.Write([]byte(`{"id":"message-id"}`))
	}))
	defer server.Close()

	o := options{anthropicBase: server.URL, anthropicKey: "test-key", anthropicVer: version}
	c := capture{provider: "anthropic", api: apiMessages, mode: "completed", scenario: "text", request: map[string]any{"model": "test"}}
	s, err := fetch(context.Background(), server.Client(), o, c)
	if err != nil {
		t.Fatal(err)
	}
	if s.APIVersion != version {
		t.Fatalf("API version = %q", s.APIVersion)
	}
	encoded, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(encoded), `"api_version":"`+version+`"`) || strings.Contains(string(encoded), "captured_at") {
		t.Fatalf("unexpected snapshot metadata: %s", encoded)
	}
}

func TestWriteJSONAtomicOverwrite(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "fixture.json")
	if err := writeJSONAtomic(path, map[string]any{"value": 1}, false); err != nil {
		t.Fatal(err)
	}
	if err := writeJSONAtomic(path, map[string]any{"value": 2}, false); err == nil || !strings.Contains(err.Error(), "refusing to overwrite") {
		t.Fatalf("second write error = %v", err)
	}
	b, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(b), "2") {
		t.Fatalf("existing file changed: %s", b)
	}
	if err := writeJSONAtomic(path, map[string]any{"value": 2}, true); err != nil {
		t.Fatal(err)
	}
	b, _ = os.ReadFile(path)
	if !strings.Contains(string(b), "2") {
		t.Fatalf("overwrite did not replace file: %s", b)
	}
}

func TestEndpointURL(t *testing.T) {
	for _, tc := range []struct{ base, endpoint, want string }{
		{"https://example.com", "/v1/responses", "https://example.com/v1/responses"},
		{"https://example.com/v1", "/v1/responses", "https://example.com/v1/responses"},
		{"https://example.com/proxy", "/v1/messages", "https://example.com/proxy/v1/messages"},
	} {
		got, err := endpointURL(tc.base, tc.endpoint)
		if err != nil || got != tc.want {
			t.Errorf("endpointURL(%q) = %q, %v; want %q", tc.base, got, err, tc.want)
		}
	}
	if _, err := endpointURL("https://secret@example.com", "/v1/responses"); err == nil {
		t.Fatal("URL credentials were accepted")
	}
}
