// Command capturegoldens captures provider-native JSON and SSE fixtures.
package main

import (
	"bufio"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	schemaVersion              = 1
	defaultAnthropicModel      = "claude-haiku-4-5"
	defaultOpenAIModel         = "gpt-5.6-luna"
	defaultAnthropicVersion    = "2023-06-01"
	maxErrorBody               = 64 << 10
	apiAll                     = "all"
	apiMessages                = "messages"
	apiChatCompletions         = "chat_completions"
	apiChatCompletionsSelector = "chat-completions"
	apiResponses               = "responses"
)

type options struct {
	outputDir      string
	provider       string
	api            string
	scenario       string
	mode           string
	overwrite      bool
	dryRun         bool
	timeout        time.Duration
	anthropicModel string
	openAIModel    string
	anthropicBase  string
	openAIBase     string
	anthropicKey   string
	openAIKey      string
	anthropicVer   string
}

type capture struct {
	provider string
	api      string
	mode     string
	scenario string
	request  map[string]any
}

type snapshot struct {
	SchemaVersion int            `json:"schema_version"`
	Provider      string         `json:"provider"`
	API           string         `json:"api"`
	Mode          string         `json:"mode"`
	Scenario      string         `json:"scenario"`
	Request       map[string]any `json:"request"`
	Response      any            `json:"response,omitempty"`
	Events        *[]sseEvent    `json:"events,omitempty"`
}

type sseEvent struct {
	Event string `json:"event,omitempty"`
	Data  any    `json:"data"`
}

func main() {
	if err := run(os.Args[1:], os.Stdout); err != nil {
		fmt.Fprintln(os.Stderr, "capturegoldens:", err)
		os.Exit(1)
	}
}

func run(args []string, stdout io.Writer) error {
	o, err := parseFlags(args)
	if err != nil {
		return err
	}
	captures, err := plan(o)
	if err != nil {
		return err
	}
	if o.dryRun {
		enc := json.NewEncoder(stdout)
		enc.SetEscapeHTML(false)
		enc.SetIndent("", "  ")
		for _, c := range captures {
			if err := enc.Encode(map[string]any{
				"provider": c.provider, "api": c.api, "mode": c.mode,
				"scenario": c.scenario, "request": c.request,
			}); err != nil {
				return err
			}
		}
		fmt.Fprintf(stdout, "dry run: planned %d request(s); no network calls or files written\n", len(captures))
		return nil
	}
	if err := validateCredentials(o, captures); err != nil {
		return err
	}
	if err := os.MkdirAll(o.outputDir, 0o755); err != nil {
		return fmt.Errorf("create output directory: %w", err)
	}
	// Preflight every destination before spending money.
	if !o.overwrite {
		for _, c := range captures {
			name := filepath.Join(o.outputDir, filename(c))
			if _, err := os.Stat(name); err == nil {
				return fmt.Errorf("refusing to overwrite %s (use -overwrite)", name)
			} else if !errors.Is(err, os.ErrNotExist) {
				return fmt.Errorf("inspect %s: %w", name, err)
			}
		}
	}

	client := &http.Client{Timeout: o.timeout}
	for i, c := range captures {
		s, err := fetch(context.Background(), client, o, c)
		if err != nil {
			return fmt.Errorf("%s/%s/%s/%s: %w (captured %d of %d; existing files were left in place)", c.provider, c.api, c.mode, c.scenario, err, i, len(captures))
		}
		name := filepath.Join(o.outputDir, filename(c))
		if err := writeJSONAtomic(name, s, o.overwrite); err != nil {
			return fmt.Errorf("write %s: %w (captured %d of %d)", name, err, i, len(captures))
		}
		fmt.Fprintln(stdout, "wrote", name)
	}
	fmt.Fprintf(stdout, "captured %d fixture(s)\n", len(captures))
	return nil
}

func parseFlags(args []string) (options, error) {
	var o options
	fs := flag.NewFlagSet("capturegoldens", flag.ContinueOnError)
	fs.SetOutput(io.Discard)
	fs.StringVar(&o.outputDir, "output", "lib/messagetranslators/testdata/provider_goldens", "output directory")
	fs.StringVar(&o.provider, "provider", "all", "all, anthropic, or openai")
	fs.StringVar(&o.api, "api", apiAll, "all, messages, chat-completions, or responses")
	fs.StringVar(&o.scenario, "scenario", "all", "all, text, or tool")
	fs.StringVar(&o.mode, "mode", "all", "all, completed, or stream")
	fs.BoolVar(&o.overwrite, "overwrite", false, "replace existing fixtures")
	fs.BoolVar(&o.dryRun, "dry-run", false, "print requests without network calls or file writes")
	fs.DurationVar(&o.timeout, "timeout", 60*time.Second, "timeout for each HTTP request")
	fs.StringVar(&o.anthropicVer, "anthropic-version", envOr("ANTHROPIC_VERSION", defaultAnthropicVersion), "Anthropic API version")
	if err := fs.Parse(args); err != nil {
		return o, err
	}
	if fs.NArg() != 0 {
		return o, fmt.Errorf("unexpected arguments: %s", strings.Join(fs.Args(), " "))
	}
	if !oneOf(o.provider, "all", "anthropic", "openai") || !oneOf(o.api, apiAll, apiMessages, apiChatCompletionsSelector, apiResponses) ||
		!oneOf(o.scenario, "all", "text", "tool") || !oneOf(o.mode, "all", "completed", "stream") {
		return o, errors.New("invalid selection: see -provider, -api, -scenario, and -mode")
	}
	if o.api == apiChatCompletionsSelector {
		o.api = apiChatCompletions
	}
	if o.timeout <= 0 {
		return o, errors.New("-timeout must be positive")
	}
	o.anthropicModel = envOr("ANTHROPIC_MODEL", defaultAnthropicModel)
	o.openAIModel = envOr("OPENAI_MODEL", defaultOpenAIModel)
	o.anthropicBase = envOr("ANTHROPIC_BASE_URL", "https://api.anthropic.com")
	o.openAIBase = envOr("OPENAI_BASE_URL", "https://api.openai.com")
	o.anthropicKey = os.Getenv("ANTHROPIC_API_KEY")
	o.openAIKey = os.Getenv("OPENAI_API_KEY")
	return o, nil
}

func plan(o options) ([]capture, error) {
	apis := []struct{ provider, api string }{{"anthropic", apiMessages}, {"openai", apiChatCompletions}, {"openai", apiResponses}}
	var out []capture
	for _, a := range apis {
		if o.provider != "all" && o.provider != a.provider || o.api != "all" && o.api != a.api {
			continue
		}
		for _, scenario := range selected(o.scenario, "text", "tool") {
			for _, mode := range selected(o.mode, "completed", "stream") {
				out = append(out, capture{a.provider, a.api, mode, scenario, buildRequest(a.api, scenario, mode, o)})
			}
		}
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("-provider=%s and -api=%s select no API", o.provider, o.api)
	}
	return out, nil
}

func selected(got string, values ...string) []string {
	if got == "all" {
		return values
	}
	return []string{got}
}

func buildRequest(api, scenario, mode string, o options) map[string]any {
	stream := mode == "stream"
	prompt := "Reply with exactly GOLDEN_OK and nothing else."
	if scenario == "tool" {
		prompt = "Call get_weather once. Use the location exactly \"São Paulo, Brazil\"."
	}
	parameters := map[string]any{
		"type":                 "object",
		"properties":           map[string]any{"location": map[string]any{"type": "string"}},
		"required":             []string{"location"},
		"additionalProperties": false,
	}
	var r map[string]any
	switch api {
	case apiMessages:
		r = map[string]any{
			"model": o.anthropicModel, "max_tokens": 64, "temperature": 0,
			"messages": []any{map[string]any{"role": "user", "content": prompt}},
			"stream":   stream,
		}
		if scenario == "tool" {
			r["tools"] = []any{map[string]any{"name": "get_weather", "description": "Get the weather for a location.", "input_schema": parameters}}
			r["tool_choice"] = map[string]any{"type": "tool", "name": "get_weather"}
		}
	case apiChatCompletions:
		r = map[string]any{
			"model": o.openAIModel, "max_completion_tokens": 64, "temperature": 0,
			"messages": []any{map[string]any{"role": "user", "content": prompt}},
			"stream":   stream,
		}
		if stream {
			r["stream_options"] = map[string]any{"include_usage": true}
		}
		if scenario == "tool" {
			r["tools"] = []any{map[string]any{"type": "function", "function": map[string]any{"name": "get_weather", "description": "Get the weather for a location.", "parameters": parameters}}}
			r["tool_choice"] = map[string]any{"type": "function", "function": map[string]any{"name": "get_weather"}}
		}
	case apiResponses:
		r = map[string]any{
			"model": o.openAIModel, "max_output_tokens": 64, "temperature": 0,
			"input": prompt, "stream": stream,
		}
		if scenario == "tool" {
			r["tools"] = []any{map[string]any{"type": "function", "name": "get_weather", "description": "Get the weather for a location.", "parameters": parameters}}
			r["tool_choice"] = map[string]any{"type": "function", "name": "get_weather"}
		}
	}
	return r
}

func validateCredentials(o options, captures []capture) error {
	for _, c := range captures {
		if c.provider == "anthropic" && o.anthropicKey == "" {
			return errors.New("ANTHROPIC_API_KEY is required for the selected captures")
		}
		if c.provider == "openai" && o.openAIKey == "" {
			return errors.New("OPENAI_API_KEY is required for the selected captures")
		}
	}
	return nil
}

func fetch(ctx context.Context, client *http.Client, o options, c capture) (snapshot, error) {
	body, err := json.Marshal(c.request)
	if err != nil {
		return snapshot{}, err
	}
	base, endpoint, key := o.openAIBase, "/v1/"+map[string]string{apiChatCompletions: "chat/completions", apiResponses: "responses"}[c.api], o.openAIKey
	if c.provider == "anthropic" {
		base, endpoint, key = o.anthropicBase, "/v1/messages", o.anthropicKey
	}
	u, err := endpointURL(base, endpoint)
	if err != nil {
		return snapshot{}, err
	}
	req, err := http.NewRequestWithContext(ctx, http.MethodPost, u, bytes.NewReader(body))
	if err != nil {
		return snapshot{}, err
	}
	req.Header.Set("Content-Type", "application/json")
	if c.provider == "anthropic" {
		req.Header.Set("x-api-key", key)
		req.Header.Set("anthropic-version", o.anthropicVer)
	} else {
		req.Header.Set("Authorization", "Bearer "+key)
	}
	resp, err := client.Do(req)
	if err != nil {
		return snapshot{}, err
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		b, readErr := io.ReadAll(io.LimitReader(resp.Body, maxErrorBody+1))
		if readErr != nil {
			return snapshot{}, fmt.Errorf("HTTP %s (reading error body: %v)", resp.Status, readErr)
		}
		truncated := ""
		if len(b) > maxErrorBody {
			b, truncated = b[:maxErrorBody], " (truncated)"
		}
		return snapshot{}, fmt.Errorf("HTTP %s: %s%s", resp.Status, strings.TrimSpace(string(b)), truncated)
	}
	s := snapshot{SchemaVersion: schemaVersion, Provider: c.provider, API: c.api, Mode: c.mode, Scenario: c.scenario, Request: c.request}
	if c.mode == "stream" {
		var events []sseEvent
		events, err = parseSSE(resp.Body)
		if events == nil {
			events = []sseEvent{}
		}
		s.Events = &events
	} else {
		s.Response, err = decodeJSON(resp.Body)
	}
	return s, err
}

func decodeJSON(r io.Reader) (any, error) {
	dec := json.NewDecoder(r)
	dec.UseNumber()
	var v any
	if err := dec.Decode(&v); err != nil {
		return nil, fmt.Errorf("decode JSON response: %w", err)
	}
	var extra any
	if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
		if err == nil {
			return nil, errors.New("decode JSON response: trailing JSON value")
		}
		return nil, fmt.Errorf("decode JSON response: trailing data: %w", err)
	}
	return v, nil
}

// parseSSE implements the SSE line and event framing rules. Unknown fields and
// comments are ignored; data lines are joined with a newline.
func parseSSE(r io.Reader) ([]sseEvent, error) {
	br := bufio.NewReader(r)
	var events []sseEvent
	var eventName string
	var data []string
	dispatch := func() {
		if len(data) == 0 {
			eventName = ""
			return
		}
		raw := strings.Join(data, "\n")
		var value any
		dec := json.NewDecoder(strings.NewReader(raw))
		dec.UseNumber()
		if err := dec.Decode(&value); err != nil {
			value = raw
		} else {
			var extra any
			if err := dec.Decode(&extra); !errors.Is(err, io.EOF) {
				value = raw
			}
		}
		events = append(events, sseEvent{Event: eventName, Data: value})
		eventName, data = "", nil
	}
	for {
		line, err := br.ReadString('\n')
		if len(line) != 0 {
			line = strings.TrimSuffix(line, "\n")
			line = strings.TrimSuffix(line, "\r")
			if line == "" {
				dispatch()
			} else if line[0] != ':' {
				field, value, found := strings.Cut(line, ":")
				if !found {
					value = ""
				} else if strings.HasPrefix(value, " ") {
					value = value[1:]
				}
				switch field {
				case "event":
					eventName = value
				case "data":
					data = append(data, value)
				}
			}
		}
		if err != nil {
			if errors.Is(err, io.EOF) {
				dispatch()
				return events, nil
			}
			return nil, fmt.Errorf("read SSE stream: %w", err)
		}
	}
}

func endpointURL(base, endpoint string) (string, error) {
	u, err := url.Parse(strings.TrimRight(base, "/"))
	if err != nil || (u.Scheme != "http" && u.Scheme != "https") || u.Host == "" || u.User != nil {
		return "", fmt.Errorf("invalid base URL %q (must be http(s), with no credentials)", base)
	}
	if strings.HasSuffix(u.Path, "/v1") && strings.HasPrefix(endpoint, "/v1/") {
		u.Path += strings.TrimPrefix(endpoint, "/v1")
	} else {
		u.Path += endpoint
	}
	return u.String(), nil
}

func filename(c capture) string {
	return strings.Join([]string{c.provider, c.api, c.mode, c.scenario}, "_") + ".json"
}

func writeJSONAtomic(path string, value any, overwrite bool) error {
	var b bytes.Buffer
	enc := json.NewEncoder(&b)
	enc.SetEscapeHTML(false)
	enc.SetIndent("", "  ")
	if err := enc.Encode(value); err != nil {
		return err
	}
	tmp, err := os.CreateTemp(filepath.Dir(path), ".capturegoldens-*")
	if err != nil {
		return err
	}
	tmpName := tmp.Name()
	defer os.Remove(tmpName)
	ok := false
	defer func() {
		if !ok {
			tmp.Close()
		}
	}()
	if err := tmp.Chmod(0o644); err != nil {
		return err
	}
	if _, err := tmp.Write(b.Bytes()); err != nil {
		return err
	}
	if err := tmp.Sync(); err != nil {
		return err
	}
	if err := tmp.Close(); err != nil {
		return err
	}
	ok = true
	if overwrite {
		return os.Rename(tmpName, path)
	}
	// A hard link publishes the complete temp file atomically and fails if the
	// destination exists, avoiding a check-then-rename overwrite race.
	if err := os.Link(tmpName, path); err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("refusing to overwrite %s (use -overwrite)", path)
		}
		return err
	}
	return os.Remove(tmpName)
}

func envOr(name, fallback string) string {
	if value := os.Getenv(name); value != "" {
		return value
	}
	return fallback
}

func oneOf(got string, values ...string) bool {
	for _, value := range values {
		if got == value {
			return true
		}
	}
	return false
}
