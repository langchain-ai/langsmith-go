//go:build integration

package integration

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

// oauthProfileName is the config profile the OAuth integration test uses. It
// must be an OAuth-logged-in profile (`langsmith auth login`), not API-key.
func oauthProfileName() string {
	if p := os.Getenv("LANGSMITH_OAUTH_PROFILE"); p != "" {
		return p
	}
	return "default"
}

// TestOAuthProfile_AuthenticatedRequest exercises the OAuth-bearer auth path
// against the real backend — the one mode `langsmith auth login` produces and
// the only one no other integration test covers (the rest use API keys).
//
// langsmith-go sends an X-User-Id header derived from the access-token JWT
// `sub`. The backend rejects the request with 403 when that id does not match
// the user it resolves from the same token, so a regression there is invisible
// to API-key auth and to the httptest-mocked unit tests.
func TestOAuthProfile_AuthenticatedRequest(t *testing.T) {
	profile := oauthProfileName()
	requireOAuthProfile(t, profile)

	// The profile is the only credential source: clear env so a stray
	// LANGSMITH_API_KEY/ENDPOINT can't mask the OAuth path under test.
	for _, k := range []string{"LANGSMITH_API_KEY", "LANGSMITH_ENDPOINT", "LANGSMITH_TENANT_ID", "LANGSMITH_WORKSPACE_ID"} {
		if v, ok := os.LookupEnv(k); ok {
			os.Unsetenv(k)
			t.Cleanup(func() { os.Setenv(k, v) })
		}
	}

	client := langsmith.NewClient(langsmith.WithProfile(profile))

	// Hit a smith-go /v2/sandboxes route — those run under the auth middleware
	// that compares X-User-Id against the user resolved from the token, so the
	// mismatch surfaces here as a 403 (the original report was this same call).
	if _, err := client.Sandboxes.Boxes.List(context.Background(), langsmith.SandboxBoxListParams{}); err != nil {
		var apiErr *langsmith.Error
		if errors.As(err, &apiErr) && apiErr.StatusCode == 403 {
			t.Fatalf("OAuth-profile request was forbidden (403) — likely the X-User-Id mismatch regression: %v", err)
		}
		t.Fatalf("OAuth-profile request failed: %v", err)
	}
}

// requireOAuthProfile skips unless profile exists and carries an OAuth access
// token. API-key profiles don't exercise the X-User-Id path, so they're skipped.
func requireOAuthProfile(t *testing.T, profile string) {
	t.Helper()
	path := os.Getenv("LANGSMITH_CONFIG_FILE")
	if path == "" {
		home, err := os.UserHomeDir()
		if err != nil {
			t.Skipf("cannot resolve home dir: %v", err)
		}
		path = filepath.Join(home, ".langsmith", "config.json")
	}
	data, err := os.ReadFile(path)
	if err != nil {
		t.Skipf("no langsmith config at %s; %s", path, loginHint(profile))
	}
	var cfg struct {
		Profiles map[string]struct {
			OAuth struct {
				AccessToken string `json:"access_token"`
			} `json:"oauth"`
		} `json:"profiles"`
	}
	if err := json.Unmarshal(data, &cfg); err != nil {
		t.Skipf("cannot parse %s: %v", path, err)
	}
	p, ok := cfg.Profiles[profile]
	if !ok {
		t.Skipf("profile %q not found in %s; %s", profile, path, loginHint(profile))
	}
	if p.OAuth.AccessToken == "" {
		t.Skipf("profile %q has no OAuth access token (API-key profiles don't exercise this path); %s", profile, loginHint(profile))
	}
}

// loginHint is the command to create the OAuth profile this test needs.
func loginHint(profile string) string {
	return fmt.Sprintf("run `langsmith auth login --profile %s`", profile)
}
