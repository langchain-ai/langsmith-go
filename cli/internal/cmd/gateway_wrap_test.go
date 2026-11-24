package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestResolveGatewayURL(t *testing.T) {
	t.Setenv(envGatewayURL, "")
	cases := []struct {
		apiURL string
		want   string
	}{
		{"https://api.smith.langchain.com", "https://gateway.smith.langchain.com"},
		{"https://api.smith.langchain.com/api/v1", "https://gateway.smith.langchain.com"},
		{"https://eu.api.smith.langchain.com", "https://eu.gateway.smith.langchain.com"},
		{"https://beta.api.smith.langchain.com/", "https://beta.gateway.smith.langchain.com"},
		{"https://smith.langchain.com", "https://gateway.smith.langchain.com"},
		{"https://langsmith.example.com/api/v1", "https://langsmith.example.com/gateway"},
		{"https://example.com/langsmith/api/v1", "https://example.com/langsmith/gateway"},
		{"http://localhost:1980/api", "http://localhost:1980/gateway"},
	}
	for _, tc := range cases {
		got, err := resolveGatewayURL("", tc.apiURL)
		if err != nil || got != tc.want {
			t.Errorf("resolveGatewayURL(%q) = %q, %v; want %q", tc.apiURL, got, err, tc.want)
		}
	}
}

func TestResolveGatewayURLOverrides(t *testing.T) {
	t.Setenv(envGatewayURL, "http://localhost:8083/")
	got, err := resolveGatewayURL("", "https://api.smith.langchain.com")
	if err != nil || got != "http://localhost:8083" {
		t.Fatalf("env override: got %q, %v", got, err)
	}
	got, err = resolveGatewayURL("https://gw.example.com", "https://api.smith.langchain.com")
	if err != nil || got != "https://gw.example.com" {
		t.Fatalf("flag override: got %q, %v", got, err)
	}
	for _, bad := range []string{"ftp://gw.example.com", "https://user:pw@gw.example.com", "https://gw.example.com?x=1", "gw.example.com"} {
		if _, err := resolveGatewayURL(bad, ""); err == nil {
			t.Errorf("resolveGatewayURL(%q) accepted an invalid URL", bad)
		}
	}
}

const testWrapKey = "test-wrap-key"

func TestNewWrapKey(t *testing.T) {
	first, err := newWrapKey()
	if err != nil {
		t.Fatal(err)
	}
	second, err := newWrapKey()
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(first, "sk-") || len(first) < 40 || first == second {
		t.Fatalf("keys %q and %q are not distinct random sk- keys", first, second)
	}
}

func TestGatewayProxyAcceptsKeyInEitherHeader(t *testing.T) {
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, "{}")
	}))
	defer upstream.Close()
	auth := func(context.Context) (http.Header, error) {
		return http.Header{"Authorization": {"Bearer oauth-access-token"}}, nil
	}
	handler, err := newGatewayProxy(upstream.URL, testWrapKey, auth)
	if err != nil {
		t.Fatal(err)
	}
	for _, tc := range []struct{ header, value string }{
		{"Authorization", "Bearer " + testWrapKey},
		{"X-Api-Key", testWrapKey},
	} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/openai/v1/chat/completions", nil)
		req.Header.Set(tc.header, tc.value)
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusOK {
			t.Errorf("%s %q: status %d, want 200", tc.header, tc.value, rec.Code)
		}
	}
}

func TestGatewayProxyReplacesCredentials(t *testing.T) {
	var got *http.Request
	var gotBody string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		got = r.Clone(context.Background())
		b, _ := io.ReadAll(r.Body)
		gotBody = string(b)
		w.Header().Set("Content-Type", "text/event-stream")
		fmt.Fprint(w, "event: message_stop\ndata: {}\n\n")
	}))
	defer upstream.Close()

	auth := func(context.Context) (http.Header, error) {
		return http.Header{"Authorization": {"Bearer oauth-access-token"}, "X-Tenant-Id": {"workspace-1"}}, nil
	}
	handler, err := newGatewayProxy(upstream.URL+"/gw", testWrapKey, auth)
	if err != nil {
		t.Fatal(err)
	}
	proxy := httptest.NewServer(handler)
	defer proxy.Close()

	req, _ := http.NewRequest(http.MethodPost, proxy.URL+"/anthropic/v1/messages?beta=true", strings.NewReader(`{"model":"claude"}`))
	req.Header.Set("Authorization", "Bearer "+testWrapKey)
	req.Header.Set("X-Api-Key", "should-not-forward")
	req.Header.Set("Cookie", "should-not-forward")
	req.Header.Set("OpenAI-Organization", "should-not-forward")
	req.Header.Set("Anthropic-Version", "2023-06-01")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatal(err)
	}
	body, _ := io.ReadAll(resp.Body)
	resp.Body.Close()

	if resp.StatusCode != http.StatusOK || !strings.Contains(string(body), "message_stop") {
		t.Fatalf("unexpected response %d: %s", resp.StatusCode, body)
	}
	if got.URL.Path != "/gw/anthropic/v1/messages" || got.URL.RawQuery != "beta=true" {
		t.Fatalf("forwarded to %s?%s", got.URL.Path, got.URL.RawQuery)
	}
	if got.Header.Get("Authorization") != "Bearer oauth-access-token" || got.Header.Get("X-Tenant-Id") != "workspace-1" {
		t.Fatalf("unexpected upstream auth headers: %v", got.Header)
	}
	if got.Header.Get("X-Api-Key") != "" || got.Header.Get("Cookie") != "" || got.Header.Get("OpenAI-Organization") != "" {
		t.Fatalf("caller credentials reached the gateway: %v", got.Header)
	}
	if got.Header.Get("Anthropic-Version") != "2023-06-01" || gotBody != `{"model":"claude"}` {
		t.Fatalf("request not forwarded intact: %v %q", got.Header, gotBody)
	}
}

func TestGatewayProxyRejectsRequestsWithoutSessionSecret(t *testing.T) {
	called := false
	upstream := httptest.NewServer(http.HandlerFunc(func(http.ResponseWriter, *http.Request) { called = true }))
	defer upstream.Close()
	authCalled := false
	auth := func(context.Context) (http.Header, error) {
		authCalled = true
		return http.Header{"Authorization": {"Bearer oauth-access-token"}}, nil
	}
	handler, err := newGatewayProxy(upstream.URL, testWrapKey, auth)
	if err != nil {
		t.Fatal(err)
	}
	for _, header := range []string{"", "Bearer wrong", testWrapKey, "Bearer " + testWrapKey + "x"} {
		rec := httptest.NewRecorder()
		req := httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
		if header != "" {
			req.Header.Set("Authorization", header)
		}
		handler.ServeHTTP(rec, req)
		if rec.Code != http.StatusUnauthorized {
			t.Errorf("Authorization %q: status %d, want 401", header, rec.Code)
		}
	}
	if called || authCalled {
		t.Fatalf("unauthenticated request reached the gateway (upstream=%v, token read=%v)", called, authCalled)
	}
}

func TestGatewayProxyReportsLoginFailure(t *testing.T) {
	auth := func(context.Context) (http.Header, error) {
		return nil, errors.New("refresh token revoked")
	}
	handler, err := newGatewayProxy("https://gateway.invalid", testWrapKey, auth)
	if err != nil {
		t.Fatal(err)
	}
	rec := httptest.NewRecorder()
	req := httptest.NewRequest(http.MethodPost, "/v1/messages", nil)
	req.Header.Set("Authorization", "Bearer "+testWrapKey)
	handler.ServeHTTP(rec, req)
	var body struct {
		Type  string `json:"type"`
		Error struct {
			Type    string `json:"type"`
			Message string `json:"message"`
		} `json:"error"`
	}
	if err := json.Unmarshal(rec.Body.Bytes(), &body); err != nil {
		t.Fatalf("not an Anthropic error body: %s", rec.Body)
	}
	if rec.Code != http.StatusUnauthorized || body.Error.Type != "authentication_error" || !strings.Contains(body.Error.Message, "refresh token revoked") {
		t.Fatalf("unexpected error %d: %s", rec.Code, rec.Body)
	}
}

func TestWrapEnvPointsAtProxy(t *testing.T) {
	env := wrapEnv([]string{
		"PATH=/bin",
		"ANTHROPIC_API_KEY=sk-ant",
		"ANTHROPIC_BASE_URL=https://api.anthropic.com",
		"CLAUDE_CODE_USE_BEDROCK=1",
		"CLAUDE_CODE_OAUTH_TOKEN=x",
		"OPENAI_API_KEY=sk-openai",
		"OPENAI_ORG_ID=org-1",
	}, "http://127.0.0.1:1234", testWrapKey)
	want := []string{
		"PATH=/bin",
		"ANTHROPIC_BASE_URL=http://127.0.0.1:1234/anthropic",
		"ANTHROPIC_AUTH_TOKEN=" + testWrapKey,
		"OPENAI_BASE_URL=http://127.0.0.1:1234/openai/v1",
		"OPENAI_API_BASE=http://127.0.0.1:1234/openai/v1",
		"OPENAI_API_KEY=" + testWrapKey,
	}
	if strings.Join(env, "\n") != strings.Join(want, "\n") {
		t.Fatalf("wrapEnv = %q, want %q", env, want)
	}
}

func TestGatewayWrapRequiresCredentials(t *testing.T) {
	resetClaudeGlobals(t)
	isolateConfig(t)
	t.Setenv("LANGSMITH_API_KEY", "")
	_, err := executeCommand(t, "gateway", "wrap", "claude")
	if err == nil || !strings.Contains(err.Error(), "not authenticated") {
		t.Fatalf("expected the shared not-authenticated error, got %v", err)
	}
}

// TestGatewayWrapRunsClaudeCodeThroughGateway runs the test binary as Claude
// Code: TestFakeClaudeCode sends one request through the proxy it was given,
// which must carry the credentials every other command would use. The
// command is given without --, and its first argument looks like a flag.
func TestGatewayWrapRunsClaudeCodeThroughGateway(t *testing.T) {
	cases := []struct {
		name       string
		profile    string
		envAPIKey  string
		wantHeader string
		wantValue  string
	}{
		{
			name:       "oauth profile",
			profile:    `"oauth": {"access_token": "saved-access-token", "expires_at": "2999-01-01T00:00:00Z"}`,
			wantHeader: "Authorization",
			wantValue:  "Bearer saved-access-token",
		},
		{
			name:       "api key profile",
			profile:    `"api_key": "lsv2_pt_profile"`,
			wantHeader: "X-Api-Key",
			wantValue:  "lsv2_pt_profile",
		},
		{
			name:       "LANGSMITH_API_KEY over an oauth profile",
			profile:    `"oauth": {"access_token": "saved-access-token", "expires_at": "2999-01-01T00:00:00Z"}`,
			envAPIKey:  "lsv2_pt_env",
			wantHeader: "X-Api-Key",
			wantValue:  "lsv2_pt_env",
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetClaudeGlobals(t)
			var got http.Header
			var gotPath string
			upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				got, gotPath = r.Header.Clone(), r.URL.Path
				fmt.Fprint(w, `{"type":"message"}`)
			}))
			defer upstream.Close()

			writeClaudeTestConfig(t, tc.profile)
			t.Setenv("LANGSMITH_API_KEY", tc.envAPIKey)
			t.Setenv(envGatewayURL, upstream.URL)
			t.Setenv("ANTHROPIC_API_KEY", "sk-ant-should-not-be-used")
			t.Setenv("LANGSMITH_FAKE_CLAUDE", "1")

			exitCode := -1
			oldExit := exitProcess
			exitProcess = func(code int) { exitCode = code }
			t.Cleanup(func() { exitProcess = oldExit })

			_, err := executeCommand(t, "gateway", "wrap", os.Args[0], "-test.run=^TestFakeClaudeCode$")
			if err != nil {
				t.Fatalf("langsmith gateway wrap returned error: %v", err)
			}
			if exitCode != 7 {
				t.Fatalf("exit code = %d, want Claude Code's 7", exitCode)
			}
			if gotPath != "/anthropic/v1/messages" || got.Get(tc.wantHeader) != tc.wantValue {
				t.Fatalf("gateway got %s %q on %q", tc.wantHeader, got.Get(tc.wantHeader), gotPath)
			}
			if tc.wantHeader == "X-Api-Key" && got.Get("Authorization") != "" {
				t.Fatalf("session secret reached the gateway: %v", got)
			}
		})
	}
}

func TestFakeClaudeCode(t *testing.T) {
	if os.Getenv("LANGSMITH_FAKE_CLAUDE") != "1" {
		t.Skip("run by TestGatewayWrapRunsClaudeCodeThroughGateway")
	}
	if os.Getenv("ANTHROPIC_API_KEY") != "" {
		os.Exit(2)
	}
	req, _ := http.NewRequest(http.MethodPost, os.Getenv("ANTHROPIC_BASE_URL")+"/v1/messages", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("ANTHROPIC_AUTH_TOKEN"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(3)
	}
	os.Exit(7)
}

func resetClaudeGlobals(t *testing.T) {
	t.Helper()
	oldKey, oldURL, oldProfile, oldWorkspace := flagAPIKey, flagAPIURL, flagProfile, flagWorkspaceID
	t.Cleanup(func() {
		flagAPIKey, flagAPIURL, flagProfile, flagWorkspaceID = oldKey, oldURL, oldProfile, oldWorkspace
	})
	flagAPIKey, flagAPIURL, flagProfile, flagWorkspaceID = "", "", "", ""
	t.Setenv("LANGSMITH_ENDPOINT", "")
	t.Setenv("LANGSMITH_PROFILE", "")
	t.Setenv("LANGSMITH_WORKSPACE_ID", "")
	t.Setenv("LANGSMITH_TENANT_ID", "")
}

func writeClaudeTestConfig(t *testing.T, profileAuth string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	config := `{"current_profile": "dev", "profiles": {"dev": {` + profileAuth + `}}}`
	if err := os.WriteFile(path, []byte(config), 0o600); err != nil {
		t.Fatal(err)
	}
}

// TestGatewayWrapRunsOpenAIClientThroughGateway runs the test binary as an
// OpenAI client: TestFakeOpenAIClient posts to OPENAI_BASE_URL.
func TestGatewayWrapRunsOpenAIClientThroughGateway(t *testing.T) {
	resetClaudeGlobals(t)
	var gotAuth, gotPath string
	upstream := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth, gotPath = r.Header.Get("Authorization"), r.URL.Path
		fmt.Fprint(w, `{"object":"chat.completion"}`)
	}))
	defer upstream.Close()

	writeClaudeTestConfig(t, `"oauth": {"access_token": "saved-access-token", "expires_at": "2999-01-01T00:00:00Z"}`)
	t.Setenv("LANGSMITH_API_KEY", "")
	t.Setenv(envGatewayURL, upstream.URL)
	t.Setenv("OPENAI_API_KEY", "sk-openai-should-not-be-used")
	t.Setenv("LANGSMITH_FAKE_OPENAI", "1")

	exitCode := -1
	oldExit := exitProcess
	exitProcess = func(code int) { exitCode = code }
	t.Cleanup(func() { exitProcess = oldExit })

	_, err := executeCommand(t, "gateway", "wrap", "--", os.Args[0], "-test.run=^TestFakeOpenAIClient$")
	if err != nil {
		t.Fatalf("langsmith gateway wrap returned error: %v", err)
	}
	if exitCode != 7 {
		t.Fatalf("exit code = %d, want the command's 7", exitCode)
	}
	if gotAuth != "Bearer saved-access-token" || gotPath != "/openai/v1/chat/completions" {
		t.Fatalf("gateway got Authorization %q on %q", gotAuth, gotPath)
	}
}

func TestFakeOpenAIClient(t *testing.T) {
	if os.Getenv("LANGSMITH_FAKE_OPENAI") != "1" {
		t.Skip("run by TestGatewayWrapRunsOpenAIClientThroughGateway")
	}
	req, _ := http.NewRequest(http.MethodPost, os.Getenv("OPENAI_BASE_URL")+"/chat/completions", strings.NewReader("{}"))
	req.Header.Set("Authorization", "Bearer "+os.Getenv("OPENAI_API_KEY"))
	resp, err := http.DefaultClient.Do(req)
	if err != nil || resp.StatusCode != http.StatusOK {
		os.Exit(3)
	}
	os.Exit(7)
}
