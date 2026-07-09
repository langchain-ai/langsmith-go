package langsmith

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"sync/atomic"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

func jwtWithSubject(t *testing.T, sub string) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":%q}`, sub)))
	return header + "." + payload + "."
}

func TestLoadProfileOptions_NoFile(t *testing.T) {
	t.Setenv("LANGSMITH_CONFIG_FILE", "/nonexistent/path/config.json")
	opts := loadProfileOptions()
	if len(opts) != 0 {
		t.Errorf("expected no options for missing file, got %d", len(opts))
	}
}

func TestLoadProfileOptions_ValidProfile(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "lsv2_pt_prodkey",
      "api_url": "https://prod.example.com",
      "workspace_id": "ws-prod"
    }
  }
}
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
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-key"
    },
    "staging": {
      "api_key": "staging-key"
    }
  }
}
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
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "default-key",
      "api_url": "https://default.example.com"
    }
  }
}
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
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "nonexistent",
  "profiles": {
    "prod": {
      "api_key": "key"
    }
  }
}
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

func TestLoadProfileOptions_InvalidJSON(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte("not valid json"), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)

	opts := loadProfileOptions()
	if len(opts) != 0 {
		t.Errorf("expected no options for invalid JSON, got %d", len(opts))
	}
}

func TestLoadProfileOptions_PartialFields(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "only-key"
    }
  }
}
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
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "profile-key"
    }
  }
}
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
	t.Setenv("LANGSMITH_CONFIG_FILE", "/nonexistent/path/config.json")
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
	accessToken := jwtWithSubject(t, "user-123")
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_url": "https://api.smith.langchain.com",
      "oauth": {
        "access_token": "` + accessToken + `"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if cfg.OAuthAccessToken != accessToken {
		t.Fatalf("expected profile access token to become OAuth access token, got %q", cfg.OAuthAccessToken)
	}
	if got := cfg.Request.Header.Get("authorization"); got != "Bearer "+accessToken {
		t.Fatalf("expected Authorization bearer header, got %q", got)
	}
	if got := cfg.Request.Header.Get("X-User-Id"); got != "user-123" {
		t.Fatalf("expected X-User-Id from access token subject, got %q", got)
	}
}

func TestLoadProfileOptions_APIKeyOverrideDropsOAuthUserID(t *testing.T) {
	clearAuthEnv(t)
	accessToken := jwtWithSubject(t, "profile-user")
	var gotAuth string
	var gotAPIKey string
	var gotUserID string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		gotAuth = r.Header.Get("Authorization")
		gotAPIKey = r.Header.Get("X-API-Key")
		gotUserID = r.Header.Get("X-User-Id")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
	}))
	defer ts.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_url": "` + ts.URL + `",
      "oauth": {
        "access_token": "` + accessToken + `"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := append(loadProfileOptions(), option.WithAPIKey("override-api-key"))
	cfg := applyOptions(t, opts)
	if cfg.OAuthAccessToken != "" {
		t.Fatalf("expected API key override to clear OAuth access token, got %q", cfg.OAuthAccessToken)
	}
	if got := cfg.Request.Header.Get("X-User-Id"); got != "" {
		t.Fatalf("expected API key override to clear initial X-User-Id, got %q", got)
	}

	preProfileOpts := append([]option.RequestOption{option.WithAPIKey("override-api-key")}, loadProfileOptions()...)
	cfg = applyOptions(t, preProfileOpts)
	if cfg.OAuthAccessToken != "" {
		t.Fatalf("expected pre-existing API key to suppress OAuth access token, got %q", cfg.OAuthAccessToken)
	}
	if got := cfg.Request.Header.Get("X-User-Id"); got != "" {
		t.Fatalf("expected pre-existing API key to suppress X-User-Id, got %q", got)
	}

	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, opts...); err != nil {
		t.Fatal(err)
	}
	if gotAPIKey != "override-api-key" {
		t.Fatalf("expected API key override, got %q", gotAPIKey)
	}
	if gotAuth != "" {
		t.Fatalf("expected OAuth Authorization to be removed, got %q", gotAuth)
	}
	if gotUserID != "" {
		t.Fatalf("expected OAuth X-User-Id to be removed, got %q", gotUserID)
	}
}

func TestLoadProfileOptions_RefreshesExpiredAccessToken(t *testing.T) {
	clearAuthEnv(t)
	tokenRequests := 0
	apiRequests := 0
	var apiAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
		case "/info":
			apiRequests++
			apiAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
			return
		default:
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
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_url": "` + ts.URL + `",
      "oauth": {
        "access_token": "old-access-token",
        "refresh_token": "old-refresh-token",
        "expires_at": "` + time.Now().Add(-time.Minute).UTC().Format(time.RFC3339) + `"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if cfg.OAuthAccessToken != "old-access-token" {
		t.Fatalf("expected initial access token before request, got %q", cfg.OAuthAccessToken)
	}
	if tokenRequests != 0 {
		t.Fatalf("expected no token request before API request, got %d", tokenRequests)
	}

	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, opts...); err != nil {
		t.Fatal(err)
	}
	if tokenRequests != 1 {
		t.Fatalf("expected one token request during API request, got %d", tokenRequests)
	}
	if apiRequests != 1 {
		t.Fatalf("expected one API request, got %d", apiRequests)
	}
	if apiAuth != "Bearer new-access-token" {
		t.Fatalf("expected refreshed bearer token on API request, got %q", apiAuth)
	}

	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(data), `"access_token": "new-access-token"`) {
		t.Fatalf("expected refreshed access token to be saved, got:\n%s", data)
	}
	if !strings.Contains(string(data), `"refresh_token": "new-refresh-token"`) {
		t.Fatalf("expected refreshed refresh token to be saved, got:\n%s", data)
	}
	if strings.Contains(string(data), `token_type`) || strings.Contains(string(data), `bearer_token`) {
		t.Fatalf("expected no token_type or bearer_token fields, got:\n%s", data)
	}
}

func TestLoadProfileOptions_RefreshesExpiredAccessTokenOnceAcrossAuthStates(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	var apiRequests atomic.Int32
	firstTokenStarted := make(chan struct{})
	releaseTokenRefresh := make(chan struct{})

	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			if tokenRequests.Add(1) == 1 {
				close(firstTokenStarted)
				<-releaseTokenRefresh
			}
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if got := r.FormValue("refresh_token"); got != "old-refresh-token" {
				t.Fatalf("expected old refresh token, got %q", got)
			}
			_ = json.NewEncoder(w).Encode(oauthTokenResponse{
				AccessToken:  "new-access-token",
				ExpiresIn:    300,
				RefreshToken: "new-refresh-token",
			})
		case "/info":
			apiRequests.Add(1)
			if got := r.Header.Get("Authorization"); got != "Bearer new-access-token" {
				t.Errorf("expected refreshed bearer token on API request, got %q", got)
			}
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_url": "` + ts.URL + `",
      "oauth": {
        "access_token": "old-access-token",
        "refresh_token": "old-refresh-token",
        "expires_at": "` + time.Now().Add(-time.Minute).UTC().Format(time.RFC3339) + `"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	optsA := loadProfileOptions()
	optsB := loadProfileOptions()

	var wg sync.WaitGroup
	wg.Add(2)
	request := func(opts []option.RequestOption) {
		defer wg.Done()
		var out map[string]string
		if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, opts...); err != nil {
			t.Error(err)
		}
	}
	go request(optsA)

	select {
	case <-firstTokenStarted:
	case <-time.After(time.Second):
		t.Fatal("first token request did not start")
	}

	go request(optsB)
	time.Sleep(25 * time.Millisecond)
	if got := tokenRequests.Load(); got != 1 {
		t.Fatalf("expected second auth state to wait on lock, token requests = %d", got)
	}

	close(releaseTokenRefresh)
	wg.Wait()

	if got := tokenRequests.Load(); got != 1 {
		t.Fatalf("expected one token request, got %d", got)
	}
	if got := apiRequests.Load(); got != 2 {
		t.Fatalf("expected two API requests, got %d", got)
	}
}

func TestAcquireOAuthRefreshDirLockBreaksExpiredTimestampLock(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	lockDir := filepath.Join(t.TempDir(), "config.json.oauth.lock.lock")
	if err := os.Mkdir(lockDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(lockDir, oauthRefreshLockTimestampFile),
		[]byte(now.Add(-oauthRefreshLockStaleAfter-time.Second).Format(time.RFC3339Nano)+"\n"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	lock, err := acquireOAuthRefreshDirLock(context.Background(), lockDir, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	defer lock.unlock()

	createdAt, ok := oauthRefreshDirLockCreatedAt(lockDir)
	if !ok {
		t.Fatal("expected refreshed lock timestamp")
	}
	if !createdAt.Equal(now) {
		t.Fatalf("expected refreshed lock timestamp %s, got %s", now, createdAt)
	}
}

func TestAcquireOAuthRefreshDirLockBreaksExpiredLockWithoutTimestamp(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	lockDir := filepath.Join(t.TempDir(), "config.json.oauth.lock.lock")
	if err := os.Mkdir(lockDir, 0700); err != nil {
		t.Fatal(err)
	}
	staleAt := now.Add(-oauthRefreshLockStaleAfter - time.Second)
	if err := os.Chtimes(lockDir, staleAt, staleAt); err != nil {
		t.Fatal(err)
	}

	lock, err := acquireOAuthRefreshDirLock(context.Background(), lockDir, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	defer lock.unlock()

	if _, err := os.Stat(filepath.Join(lockDir, oauthRefreshLockTimestampFile)); err != nil {
		t.Fatalf("expected timestamp file after taking stale lock: %v", err)
	}
}

func TestOAuthRefreshDirLockUnlockDoesNotRemoveNewOwner(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	lockDir := filepath.Join(t.TempDir(), "config.json.oauth.lock.lock")
	if err := os.Mkdir(lockDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := writeOAuthRefreshLockMetadata(
		lockDir,
		now.Add(-oauthRefreshLockStaleAfter-time.Second),
		"stale-owner",
	); err != nil {
		t.Fatal(err)
	}

	lock, err := acquireOAuthRefreshDirLock(context.Background(), lockDir, func() time.Time { return now })
	if err != nil {
		t.Fatal(err)
	}
	defer lock.unlock()

	staleLock := &oauthRefreshDirLock{path: lockDir, owner: "stale-owner"}
	if err := staleLock.unlock(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(lockDir); err != nil {
		t.Fatalf("expected stale owner unlock to leave new lock in place: %v", err)
	}
}

func TestAcquireOAuthRefreshDirLockWaitsForFreshLock(t *testing.T) {
	now := time.Unix(1700000000, 0).UTC()
	lockDir := filepath.Join(t.TempDir(), "config.json.oauth.lock.lock")
	if err := os.Mkdir(lockDir, 0700); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(
		filepath.Join(lockDir, oauthRefreshLockTimestampFile),
		[]byte(now.Add(-oauthRefreshLockStaleAfter+time.Second).Format(time.RFC3339Nano)+"\n"),
		0600,
	); err != nil {
		t.Fatal(err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Millisecond)
	defer cancel()
	_, err := acquireOAuthRefreshDirLock(ctx, lockDir, func() time.Time { return now })
	if !errors.Is(err, context.DeadlineExceeded) {
		t.Fatalf("expected deadline while fresh lock exists, got %v", err)
	}
	if _, err := os.Stat(lockDir); err != nil {
		t.Fatalf("expected fresh lock directory to remain: %v", err)
	}
}

func TestLoadProfileOptions_RefreshTokenOnlyBeforeProfileAPIKey(t *testing.T) {
	clearAuthEnv(t)
	tokenRequests := 0
	var apiAuth string
	var apiKey string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			tokenRequests++
			if err := r.ParseForm(); err != nil {
				t.Fatal(err)
			}
			if got := r.FormValue("refresh_token"); got != "profile-refresh-token" {
				t.Fatalf("expected profile refresh token, got %q", got)
			}
			_ = json.NewEncoder(w).Encode(oauthTokenResponse{
				AccessToken:  "new-access-token",
				ExpiresIn:    300,
				RefreshToken: "new-refresh-token",
			})
		case "/info":
			apiAuth = r.Header.Get("Authorization")
			apiKey = r.Header.Get("X-API-Key")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
		default:
			http.NotFound(w, r)
		}
	}))
	defer ts.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "profile-api-key",
      "api_url": "` + ts.URL + `",
      "oauth": {
        "refresh_token": "profile-refresh-token"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := loadProfileOptions()
	cfg := applyOptions(t, opts)
	if got := cfg.Request.Header.Get("X-API-Key"); got != "" {
		t.Fatalf("expected no initial X-API-Key before refresh, got %q", got)
	}
	if tokenRequests != 0 {
		t.Fatalf("expected no token request before API request, got %d", tokenRequests)
	}

	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, opts...); err != nil {
		t.Fatal(err)
	}
	if tokenRequests != 1 {
		t.Fatalf("expected one token request during API request, got %d", tokenRequests)
	}
	if apiKey != "" {
		t.Fatalf("expected refreshed OAuth to suppress profile API key, got %q", apiKey)
	}
	if apiAuth != "Bearer new-access-token" {
		t.Fatalf("expected refreshed bearer token on API request, got %q", apiAuth)
	}
}

func TestLoadProfileOptions_RefreshUsesProfileAPIURL(t *testing.T) {
	clearAuthEnv(t)
	tokenRequests := 0
	profileTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/oauth/token" {
			http.NotFound(w, r)
			return
		}
		tokenRequests++
		_ = json.NewEncoder(w).Encode(oauthTokenResponse{
			AccessToken:  "new-access-token",
			ExpiresIn:    300,
			RefreshToken: "new-refresh-token",
		})
	}))
	defer profileTS.Close()

	apiRequests := 0
	overrideTokenRequests := 0
	var apiAuth string
	apiTS := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/info":
			apiRequests++
			apiAuth = r.Header.Get("Authorization")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"ok": "true"})
		case "/oauth/token":
			overrideTokenRequests++
			http.Error(w, "refresh should use profile api_url", http.StatusTeapot)
		default:
			http.NotFound(w, r)
		}
	}))
	defer apiTS.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_url": "` + profileTS.URL + `",
      "oauth": {
        "access_token": "old-access-token",
        "refresh_token": "old-refresh-token",
        "expires_at": "` + time.Now().Add(-time.Minute).UTC().Format(time.RFC3339) + `"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	opts := append(loadProfileOptions(), option.WithBaseURL(apiTS.URL))
	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, opts...); err != nil {
		t.Fatal(err)
	}
	if tokenRequests != 1 {
		t.Fatalf("expected refresh against profile api_url, got %d token requests", tokenRequests)
	}
	if overrideTokenRequests != 0 {
		t.Fatalf("expected no refresh against request base URL, got %d", overrideTokenRequests)
	}
	if apiRequests != 1 {
		t.Fatalf("expected one API request, got %d", apiRequests)
	}
	if apiAuth != "Bearer new-access-token" {
		t.Fatalf("expected refreshed bearer token on API request, got %q", apiAuth)
	}
}

func TestLoadProfileOptions_OAuthAccessTokenOverridesProfileAPIKey(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "some-key",
      "oauth": {
        "access_token": "oauth-access-token"
      }
    }
  }
}
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
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "profile-key",
      "api_url": "https://profile.example.com",
      "workspace_id": "ws-profile",
      "oauth": {
        "access_token": "profile-access-token"
      }
    }
  }
}
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

func TestWithProfileOverridesDefaultProfile(t *testing.T) {
	clearAuthEnv(t)
	var gotAPIKey string
	var gotAuth string
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/info" {
			t.Fatalf("unexpected path: %s", r.URL.Path)
		}
		gotAPIKey = r.Header.Get("X-API-Key")
		gotAuth = r.Header.Get("Authorization")
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"version": "test"})
	}))
	defer ts.Close()

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "default",
  "profiles": {
    "default": {
      "api_key": "default-api-key"
    },
    "prod": {
      "oauth": {
        "access_token": "prod-access-token"
      }
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	client := NewClient(WithProfile("prod"), option.WithBaseURL(ts.URL), option.WithMaxRetries(0))
	if _, err := client.Info.List(context.Background()); err != nil {
		t.Fatal(err)
	}
	if gotAuth != "Bearer prod-access-token" {
		t.Fatalf("expected explicit profile bearer auth, got %q", gotAuth)
	}
	if gotAPIKey != "" {
		t.Fatalf("expected explicit OAuth profile to override default API key, got %q", gotAPIKey)
	}
}

// unsetTenantEnv removes the tenant/workspace env vars (absent, not empty:
// DefaultClientOptions treats a present-but-empty value as an explicit clear).
func unsetTenantEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"LANGSMITH_TENANT_ID", "LANGSMITH_WORKSPACE_ID"} {
		if orig, ok := os.LookupEnv(k); ok {
			k, orig := k, orig
			t.Cleanup(func() { _ = os.Setenv(k, orig) })
		}
		_ = os.Unsetenv(k)
	}
}

func TestWithProfileClearsInheritedTenantWhenWorkspaceless(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "workspace_id": "prod-workspace-id"
    },
    "aws": {
      "api_key": "aws-api-key"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("aws")))
	if cfg.TenantID != "" {
		t.Fatalf("expected tenant cleared for workspace-less profile, got %q", cfg.TenantID)
	}
	if got := cfg.Request.Header.Get("X-Tenant-Id"); got != "" {
		t.Fatalf("expected no X-Tenant-Id header, got %q: current_profile tenant leaked", got)
	}
}

func TestWithProfileAppliesProfileWorkspace(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "workspace_id": "prod-workspace-id"
    },
    "staging": {
      "api_key": "staging-api-key",
      "workspace_id": "staging-workspace-id"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("staging")))
	if cfg.TenantID != "staging-workspace-id" {
		t.Fatalf("expected selected profile workspace tenant, got %q", cfg.TenantID)
	}
	if got := cfg.Request.Header.Get("X-Tenant-Id"); got != "staging-workspace-id" {
		t.Fatalf("expected X-Tenant-Id=staging-workspace-id, got %q", got)
	}
}

func TestWithProfileKeepsEnvTenantForWorkspacelessProfile(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)
	t.Setenv("LANGSMITH_TENANT_ID", "env-tenant")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "workspace_id": "prod-workspace-id"
    },
    "aws": {
      "api_key": "aws-api-key"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	// A tenant env var must win over the workspace-less profile's clear.
	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("aws")))
	if cfg.TenantID != "env-tenant" {
		t.Fatalf("expected env tenant to survive WithProfile, got %q", cfg.TenantID)
	}
	if got := cfg.Request.Header.Get("X-Tenant-Id"); got != "env-tenant" {
		t.Fatalf("expected X-Tenant-Id=env-tenant, got %q", got)
	}
}

func TestWithProfileClearsInheritedBaseURLWhenURLless(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)
	// Endpoint must be absent for the reset to apply (a present-but-empty value
	// makes DefaultClientOptions apply WithBaseURL("")). clearAuthEnv's cleanup
	// restores the original value after the test.
	_ = os.Unsetenv("LANGSMITH_ENDPOINT")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "api_url": "https://prod-host.example.com",
      "workspace_id": "prod-workspace-id"
    },
    "aws": {
      "api_key": "aws-api-key"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	// A url-less explicit profile must not inherit current_profile's api_url;
	// BaseURL resets to nil so the SDK falls back to the default endpoint.
	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("aws")))
	if cfg.BaseURL != nil {
		t.Fatalf("expected base URL cleared (nil => default) for url-less profile, got %q: current_profile URL leaked", cfg.BaseURL)
	}
}

func TestWithProfileKeepsEnvEndpointForURLlessProfile(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)
	t.Setenv("LANGSMITH_ENDPOINT", "https://env-host.example.com")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "api_url": "https://prod-host.example.com",
      "workspace_id": "prod-workspace-id"
    },
    "aws": {
      "api_key": "aws-api-key"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	// LANGSMITH_ENDPOINT must win over the url-less profile's base-URL reset.
	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("aws")))
	if cfg.BaseURL == nil || cfg.BaseURL.String() != "https://env-host.example.com" {
		t.Fatalf("expected env endpoint to survive WithProfile, got %v", cfg.BaseURL)
	}
}

func TestWithProfileEnvTenantBeatsProfileWorkspace(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)
	t.Setenv("LANGSMITH_WORKSPACE_ID", "env-workspace-id")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "workspace_id": "prod-workspace-id"
    },
    "staging": {
      "api_key": "staging-api-key",
      "workspace_id": "staging-workspace-id"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	// A tenant env var must win over the selected profile's own workspace_id,
	// matching env-overrides-profile precedence.
	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("staging")))
	if cfg.TenantID != "env-workspace-id" {
		t.Fatalf("expected env workspace to win over profile workspace, got %q", cfg.TenantID)
	}
	if got := cfg.Request.Header.Get("X-Tenant-Id"); got != "env-workspace-id" {
		t.Fatalf("expected X-Tenant-Id=env-workspace-id, got %q", got)
	}
}

func TestWithProfileEnvEndpointBeatsProfileURL(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)
	t.Setenv("LANGSMITH_ENDPOINT", "https://env-host.example.com")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key",
      "api_url": "https://prod-host.example.com"
    },
    "staging": {
      "api_key": "staging-api-key",
      "api_url": "https://staging-host.example.com"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	// LANGSMITH_ENDPOINT must win over the selected profile's own api_url.
	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("staging")))
	if cfg.BaseURL == nil || cfg.BaseURL.String() != "https://env-host.example.com" {
		t.Fatalf("expected env endpoint to win over profile api_url, got %v", cfg.BaseURL)
	}
}

func TestWithProfileClearsInheritedAPIKeyForCredentiallessProfile(t *testing.T) {
	clearAuthEnv(t)
	unsetTenantEnv(t)
	// API key env must be absent (not empty): a present-but-empty value makes
	// DefaultClientOptions apply WithAPIKey(""), which would clear the ambient key
	// before WithProfile runs and mask the leak. clearAuthEnv's cleanup restores it.
	_ = os.Unsetenv("LANGSMITH_API_KEY")

	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "current_profile": "prod",
  "profiles": {
    "prod": {
      "api_key": "prod-api-key"
    },
    "bare": {
      "workspace_id": "bare-workspace-id"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")

	// An explicitly selected profile with no credentials must not inherit
	// current_profile's API key.
	cfg := applyOptions(t, append(DefaultClientOptions(), WithProfile("bare")))
	if cfg.APIKey != "" {
		t.Fatalf("expected inherited API key cleared for credential-less profile, got %q", cfg.APIKey)
	}
	if got := cfg.Request.Header.Get("X-API-Key"); got != "" {
		t.Fatalf("expected no X-API-Key header, got %q: current_profile key leaked", got)
	}
}

func TestWithProfileMissingProfileReturnsError(t *testing.T) {
	clearAuthEnv(t)
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	content := `{
  "profiles": {
    "default": {
      "api_key": "default-api-key"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(content), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)

	req, err := http.NewRequest(http.MethodGet, "https://example.com", nil)
	if err != nil {
		t.Fatal(err)
	}
	cfg := requestconfig.RequestConfig{Request: req, HTTPClient: http.DefaultClient}
	err = cfg.Apply(WithProfile("missing"))
	if err == nil {
		t.Fatal("expected missing profile error")
	}
	if !strings.Contains(err.Error(), "profile not found: missing") {
		t.Fatalf("unexpected error: %v", err)
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
	t.Setenv("LANGSMITH_ENDPOINT", "")
}
