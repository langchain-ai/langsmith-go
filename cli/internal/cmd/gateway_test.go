package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

const gatewayModelsFixture = `{"object":"list","data":[
  {"id":"openai/gpt-5","object":"model","owned_by":"system"},
  {"id":"anthropic/claude-sonnet-5","object":"model"},
  {"id":"moonshotai/Kimi-K3","object":"model","owned_by":"langchain","supported_endpoints":["/v1/messages","/v1/responses"]},
  {"id":"anthropic/claude-haiku-5-5","object":"model"}
]}`

// fakeGateway serves the models fixture and records the auth headers of the
// last request.
func fakeGateway(t *testing.T, status int, body string) (*httptest.Server, *http.Header) {
	t.Helper()
	var got http.Header
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet || r.URL.Path != "/v1/models" {
			http.NotFound(w, r)
			return
		}
		got = r.Header.Clone()
		w.WriteHeader(status)
		fmt.Fprint(w, body)
	}))
	t.Cleanup(srv.Close)
	return srv, &got
}

func TestGatewayModelsSendsCredentials(t *testing.T) {
	cases := []struct {
		name       string
		profile    string
		envAPIKey  string
		wantHeader string
		wantValue  string
	}{
		{"oauth profile", `"oauth": {"access_token": "saved-access-token", "expires_at": "2999-01-01T00:00:00Z"}, "workspace_id": "00000000-0000-0000-0000-000000000123"`, "", "Authorization", "Bearer saved-access-token"},
		{"api key profile", `"api_key": "lsv2_pt_profile"`, "", "X-Api-Key", "lsv2_pt_profile"},
		{"LANGSMITH_API_KEY", `"oauth": {"access_token": "saved-access-token", "expires_at": "2999-01-01T00:00:00Z"}`, "lsv2_pt_env", "X-Api-Key", "lsv2_pt_env"},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			resetClaudeGlobals(t)
			srv, got := fakeGateway(t, http.StatusOK, gatewayModelsFixture)
			writeClaudeTestConfig(t, tc.profile)
			t.Setenv("LANGSMITH_API_KEY", tc.envAPIKey)

			if _, err := executeCommand(t, "gateway", "models", "--gateway-url", srv.URL, "--format", "json"); err != nil {
				t.Fatalf("gateway models returned error: %v", err)
			}
			if got.Get(tc.wantHeader) != tc.wantValue {
				t.Fatalf("gateway got %s %q; headers %v", tc.wantHeader, got.Get(tc.wantHeader), *got)
			}
			if strings.Contains(tc.profile, "workspace_id") && got.Get("X-Tenant-Id") != "00000000-0000-0000-0000-000000000123" {
				t.Fatalf("workspace not sent: %v", *got)
			}
		})
	}
}

func TestGatewayModelsListsAndFilters(t *testing.T) {
	resetClaudeGlobals(t)
	srv, _ := fakeGateway(t, http.StatusOK, gatewayModelsFixture)
	writeClaudeTestConfig(t, `"api_key": "lsv2_pt_profile"`)
	t.Setenv("LANGSMITH_API_KEY", "")
	t.Setenv(envGatewayURL, srv.URL)

	stdout, err := executeCommand(t, "gateway", "models", "--format", "json")
	if err != nil {
		t.Fatalf("gateway models returned error: %v", err)
	}
	var models []gatewayModel
	if err := json.Unmarshal([]byte(stdout), &models); err != nil {
		t.Fatalf("output is not a JSON array of models: %v\n%s", err, stdout)
	}
	var ids []string
	for _, m := range models {
		ids = append(ids, m.ID)
	}
	want := "anthropic/claude-haiku-5-5,anthropic/claude-sonnet-5,moonshotai/Kimi-K3,openai/gpt-5"
	if strings.Join(ids, ",") != want {
		t.Fatalf("models = %v, want sorted %s", ids, want)
	}
	if models[2].OwnedBy != "langchain" || len(models[2].SupportedEndpoints) != 2 {
		t.Fatalf("model fields not carried through: %+v", models[2])
	}

	stdout, err = executeCommand(t, "gateway", "models", "--provider", "Anthropic", "--jq", ".[].id")
	if err != nil {
		t.Fatalf("gateway models --provider returned error: %v", err)
	}
	if strings.TrimSpace(stdout) != "anthropic/claude-haiku-5-5\nanthropic/claude-sonnet-5" {
		t.Fatalf("--provider anthropic listed %q", stdout)
	}
}

func TestGatewayModelsReportsGatewayErrors(t *testing.T) {
	resetClaudeGlobals(t)
	srv, _ := fakeGateway(t, http.StatusForbidden, `{"error":"LLM Gateway is not enabled for this organization"}`)
	writeClaudeTestConfig(t, `"api_key": "lsv2_pt_profile"`)
	t.Setenv("LANGSMITH_API_KEY", "")

	_, err := executeCommand(t, "gateway", "models", "--gateway-url", srv.URL)
	if err == nil || !strings.Contains(err.Error(), "HTTP 403") || !strings.Contains(err.Error(), "not enabled for this organization") {
		t.Fatalf("expected the gateway's 403 and message, got %v", err)
	}
}

func TestGatewayModelsDoesNotFollowRedirects(t *testing.T) {
	resetClaudeGlobals(t)
	var leaked http.Header
	elsewhere := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		leaked = r.Header.Clone()
		fmt.Fprint(w, gatewayModelsFixture)
	}))
	defer elsewhere.Close()
	redirecting := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, elsewhere.URL+r.URL.Path, http.StatusFound)
	}))
	defer redirecting.Close()
	writeClaudeTestConfig(t, `"api_key": "lsv2_pt_profile", "workspace_id": "00000000-0000-0000-0000-000000000123"`)
	t.Setenv("LANGSMITH_API_KEY", "")

	_, err := executeCommand(t, "gateway", "models", "--gateway-url", redirecting.URL)
	if err == nil || !strings.Contains(err.Error(), "HTTP 302") {
		t.Fatalf("expected the redirect to be reported, got %v", err)
	}
	if leaked != nil {
		t.Fatalf("redirect target received the request with headers %v", leaked)
	}
}
