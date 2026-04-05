package langsmith

import (
	"os"
	"path/filepath"
	"testing"

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
	t.Setenv("LANGSMITH_BEARER_TOKEN", "")
	t.Setenv("LANGSMITH_ORGANIZATION_ID", "")

	opts := DefaultClientOptions()
	// Should have at least: WithEnvironmentProduction + profile api_key
	if len(opts) < 2 {
		t.Errorf("expected at least 2 options (production env + profile key), got %d", len(opts))
	}
	// Verify it's usable (doesn't panic)
	_ = option.WithAPIKey("override")
}

func TestLoadProfileOptions_BearerToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
bearer_token = "eyJhbGciOiJSUzI1NiJ9.test"
api_url = "https://api.smith.langchain.com"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	if len(opts) != 2 {
		t.Fatalf("expected 2 options (base_url, bearer_token), got %d", len(opts))
	}
}

func TestLoadProfileOptions_BothAPIKeyAndBearerToken(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.toml")
	content := `[default]
api_key = "some-key"
bearer_token = "eyJhbGciOiJSUzI1NiJ9.test"
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	// Both are emitted — server decides which to honor (matches env var behavior)
	if len(opts) != 2 {
		t.Fatalf("expected 2 options (api_key + bearer_token), got %d", len(opts))
	}
}
