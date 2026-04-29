package langsmith

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestLoadProfileOptions_NoFile(t *testing.T) {
	t.Setenv("LANGSMITH_CONFIG_FILE", "/nonexistent/path/config.toml")
	opts := loadProfileOptions()
	if len(opts) != 0 {
		t.Errorf("expected no options for missing file, got %d", len(opts))
	}
}

func TestLoadProfileOptions_ValidProfile(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `current_profile = "prod"

[prod]
api_key = "lsv2_pt_prodkey"
api_url = "https://prod.example.com"
workspace_id = "ws-prod"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	if len(opts) != 3 {
		t.Fatalf("expected 3 options (base_url, api_key, tenant_id), got %d", len(opts))
	}
}

func TestLoadProfileOptions_EnvProfileOverride(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `current_profile = "prod"

[prod]
api_key = "prod-key"

[staging]
api_key = "staging-key"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "staging")

	opts := loadProfileOptions()
	if len(opts) != 1 {
		t.Fatalf("expected 1 option (api_key only), got %d", len(opts))
	}
}

func TestLoadProfileOptions_FallbackToDefault(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_key = "default-key"
api_url = "https://default.example.com"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	if len(opts) != 2 {
		t.Fatalf("expected 2 options (base_url, api_key), got %d", len(opts))
	}
}

func TestLoadProfileOptions_NoMatchingProfile(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `current_profile = "nonexistent"

[prod]
api_key = "key"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	if len(opts) != 0 {
		t.Errorf("expected no options for missing profile, got %d", len(opts))
	}
}

func TestLoadProfileOptions_InvalidTOML(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	if err := os.WriteFile(path, []byte("not valid [[ toml"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)

	opts := loadProfileOptions()
	if len(opts) != 0 {
		t.Errorf("expected no options for invalid TOML, got %d", len(opts))
	}
}

func TestLoadProfileOptions_PartialFields(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_key = "only-key"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	if len(opts) != 1 {
		t.Fatalf("expected 1 option (api_key only), got %d", len(opts))
	}
}

func TestDefaultClientOptions_IncludesProfileOptions(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_key = "profile-key"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")
	// Clear env vars so profile is the only source
	t.Setenv("LANGSMITH_API_KEY", "")
	t.Setenv("LANGSMITH_ENDPOINT", "")
	t.Setenv("LANGSMITH_TENANT_ID", "")

	opts := DefaultClientOptions()
	// Should have at least: WithEnvironmentProduction + profile api_key
	if len(opts) < 2 {
		t.Errorf("expected at least 2 options (production env + profile key), got %d", len(opts))
	}
	// Verify it's usable (doesn't panic)
	_ = option.WithAPIKey("override")
}

func TestDefaultClientOptions_WorkspaceIDEnvAlias(t *testing.T) {
	t.Setenv("LANGSMITH_CONFIG_FILE", "/nonexistent/path/config.toml")
	t.Setenv("LANGSMITH_API_KEY", "")
	t.Setenv("LANGSMITH_ENDPOINT", "")
	t.Setenv("LANGSMITH_TENANT_ID", "tenant-env")
	t.Setenv("LANGSMITH_WORKSPACE_ID", "workspace-env")

	cfg := applyOptions(t, DefaultClientOptions())
	if cfg.TenantID != "workspace-env" {
		t.Fatalf("expected LANGSMITH_WORKSPACE_ID to override tenant ID, got %q", cfg.TenantID)
	}
	if got := cfg.Request.Header.Get("X-Tenant-Id"); got != "workspace-env" {
		t.Fatalf("expected X-Tenant-Id header from workspace env, got %q", got)
	}
}

func TestLoadProfileOptions_OAuthAccessToken(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_url = "https://api.smith.langchain.com"

[default.oauth]
access_token = "test-access-token"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if cfg.OAuthAccessToken != "test-access-token" {
		t.Fatalf("expected profile access token to become OAuth access token, got %q", cfg.OAuthAccessToken)
	}
	if got := cfg.Request.Header.Get("authorization"); got != "Bearer test-access-token" {
		t.Fatalf("expected Authorization bearer header, got %q", got)
	}
}

func TestLoadProfileOptions_RefreshesExpiredAccessToken(t *testing.T) {
	tokenRequests := 0
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			http.NotFound(w, r)
			return
		}
		tokenRequests++
		if err := r.ParseForm(); err != nil {
			t.Fatal(err)
		}
		if got := r.FormValue("grant_type"); got != "refresh_token" {
			t.Fatalf("expected refresh_token grant, got %q", got)
		}
		if got := r.FormValue("client_id"); got != oauthClientID {
			t.Fatalf("expected client_id %q, got %q", oauthClientID, got)
		}
		if got := r.FormValue("refresh_token"); got != "old-refresh-token" {
			t.Fatalf("expected old refresh token, got %q", got)
		}
		_ = json.NewEncoder(w).Encode(oauthTokenResponse{
			AccessToken:  "new-access-token",
			ExpiresIn:    300,
			RefreshToken: "new-refresh-token",
		})
	}))
	defer ts.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_url = "` + ts.URL + `"

[default.oauth]
access_token = "old-access-token"
refresh_token = "old-refresh-token"
expires_at = "` + time.Now().Add(-time.Minute).UTC().Format(time.RFC3339) + `"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")
	t.Setenv("LANGSMITH_API_KEY", "")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if cfg.OAuthAccessToken != "new-access-token" {
		t.Fatalf("expected refreshed access token, got %q", cfg.OAuthAccessToken)
	}
	if tokenRequests != 1 {
		t.Fatalf("expected one token request, got %d", tokenRequests)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `access_token = "new-access-token"`) {
		t.Fatalf("expected refreshed access token to be saved, got:\n%s", data)
	}
	if !strings.Contains(string(data), `refresh_token = "new-refresh-token"`) {
		t.Fatalf("expected refreshed refresh token to be saved, got:\n%s", data)
	}
	if strings.Contains(string(data), `token_type`) || strings.Contains(string(data), `bearer_token`) {
		t.Fatalf("expected no token_type or bearer_token fields, got:\n%s", data)
	}
}

func TestLoadProfileOptions_OAuthAccessTokenOverridesProfileAPIKey(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_key = "some-key"

[default.oauth]
access_token = "oauth-access-token"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if cfg.APIKey != "" {
		t.Fatalf("expected profile OAuth token to take precedence over profile API key")
	}
	if cfg.OAuthAccessToken != "oauth-access-token" {
		t.Fatalf("expected profile OAuth token, got %q", cfg.OAuthAccessToken)
	}
	if got := cfg.Request.Header.Get("X-API-Key"); got != "" {
		t.Fatalf("expected no X-API-Key header, got %q", got)
	}
}

func TestLoadProfileOptions_EnvAuthSuppressesProfileAuth(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_key = "profile-key"
api_url = "https://profile.example.com"
workspace_id = "ws-profile"

[default.oauth]
access_token = "profile-access-token"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")
	t.Setenv("LANGSMITH_API_KEY", "env-api-key")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if cfg.APIKey != "" || cfg.OAuthAccessToken != "" {
		t.Fatalf("expected profile auth to be suppressed when env auth is set")
	}
	if got := cfg.Request.Header.Get("authorization"); got != "" {
		t.Fatalf("expected no profile Authorization header, got %q", got)
	}
	if cfg.TenantID != "ws-profile" {
		t.Fatalf("expected profile workspace to remain available, got %q", cfg.TenantID)
	}
}

func applyOptions(t *testing.T, opts []option.RequestOption) requestconfig.RequestConfig {
	t.Helper()
	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestconfig.RequestConfig{Request: req, HTTPClient: http.DefaultClient}
	if err := cfg.Apply(opts...); err != nil {
		t.Fatal(err)
	}
	return cfg
}

func clearAuthEnv(t *testing.T) {
	t.Helper()
	t.Setenv("LANGSMITH_API_KEY", "")
}
