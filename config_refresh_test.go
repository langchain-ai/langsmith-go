package langsmith

import (
	"context"
	"encoding/json"
	"errors"
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

// writeOAuthProfileConfig writes a config with one OAuth profile whose access
// token has already expired, plus fields the SDK does not model.
func writeOAuthProfileConfig(t *testing.T, apiURL, issuer string) string {
	t.Helper()
	oauth := map[string]any{
		"access_token":  "old-access-token",
		"refresh_token": "old-refresh-token",
		"expires_at":    time.Now().Add(-time.Minute).UTC().Format(time.RFC3339),
	}
	if issuer != "" {
		oauth["issuer"] = issuer
	}
	cfg := map[string]any{
		"current_profile":  "default",
		"future_top_level": "keep-me",
		"profiles": map[string]any{
			"default": map[string]any{
				"api_url":        apiURL,
				"future_profile": "keep-me-too",
				"oauth":          oauth,
			},
			"other": map[string]any{"api_key": "other-key"},
		},
	}
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		t.Fatal(err)
	}
	path := filepath.Join(t.TempDir(), "config.json")
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	t.Setenv("LANGSMITH_PROFILE", "")
	return path
}

func readConfigMap(t *testing.T, path string) map[string]any {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatal(err)
	}
	var out map[string]any
	if err := json.Unmarshal(data, &out); err != nil {
		t.Fatalf("config is not valid JSON: %v\n%s", err, data)
	}
	return out
}

func tokenServer(t *testing.T, tokenRequests *atomic.Int32, respond func(w http.ResponseWriter, r *http.Request)) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/oauth/token":
			tokenRequests.Add(1)
			if err := r.ParseForm(); err != nil {
				t.Error(err)
			}
			if got, want := r.FormValue("resource"), "http://"+r.Host; got != want {
				t.Errorf("expected resource %q, got %q", want, got)
			}
			respond(w, r)
		case "/info":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{"auth": r.Header.Get("Authorization")})
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(ts.Close)
	return ts
}

func rotateTo(accessToken, refreshToken string) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		_ = json.NewEncoder(w).Encode(oauthTokenResponse{
			AccessToken:  accessToken,
			ExpiresIn:    300,
			RefreshToken: refreshToken,
		})
	}
}

func TestRefreshPreservesFieldsTheSDKDoesNotModel(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	ts := tokenServer(t, &tokenRequests, rotateTo("new-access-token", "new-refresh-token"))
	path := writeOAuthProfileConfig(t, ts.URL, ts.URL)

	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, loadProfileOptions()...); err != nil {
		t.Fatal(err)
	}
	if out["auth"] != "Bearer new-access-token" {
		t.Fatalf("expected refreshed token on request, got %q", out["auth"])
	}

	cfg := readConfigMap(t, path)
	if cfg["future_top_level"] != "keep-me" || cfg["current_profile"] != "default" {
		t.Fatalf("top-level fields not preserved: %v", cfg)
	}
	profiles := cfg["profiles"].(map[string]any)
	if other := profiles["other"].(map[string]any); other["api_key"] != "other-key" {
		t.Fatalf("other profile not preserved: %v", other)
	}
	p := profiles["default"].(map[string]any)
	if p["future_profile"] != "keep-me-too" {
		t.Fatalf("profile field not preserved: %v", p)
	}
	oauth := p["oauth"].(map[string]any)
	if oauth["issuer"] != ts.URL {
		t.Fatalf("expected issuer %q to survive refresh, got %v", ts.URL, oauth["issuer"])
	}
	if oauth["access_token"] != "new-access-token" || oauth["refresh_token"] != "new-refresh-token" {
		t.Fatalf("expected rotated tokens saved, got %v", oauth)
	}
	info, err := os.Stat(path)
	if err != nil {
		t.Fatal(err)
	}
	if perm := info.Mode().Perm(); perm != 0600 {
		t.Fatalf("expected config mode 0600, got %#o", perm)
	}
}

func TestRefreshUsesProfileIssuer(t *testing.T) {
	clearAuthEnv(t)
	var issuerRequests, apiHostTokenRequests atomic.Int32
	issuer := tokenServer(t, &issuerRequests, rotateTo("issuer-access-token", "issuer-refresh-token"))
	api := tokenServer(t, &apiHostTokenRequests, rotateTo("wrong", "wrong"))
	writeOAuthProfileConfig(t, api.URL, issuer.URL)

	token, err := ProfileAccessToken(context.Background(), "")
	if err != nil {
		t.Fatal(err)
	}
	if token != "issuer-access-token" {
		t.Fatalf("expected token from issuer, got %q", token)
	}
	if issuerRequests.Load() != 1 || apiHostTokenRequests.Load() != 0 {
		t.Fatalf("expected refresh against issuer only; issuer=%d api=%d", issuerRequests.Load(), apiHostTokenRequests.Load())
	}
}

func TestProfileAccessTokenReturnsCurrentTokenWithoutRefreshing(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	ts := tokenServer(t, &tokenRequests, rotateTo("new-access-token", "new-refresh-token"))
	path := writeOAuthProfileConfig(t, ts.URL, "")
	cfg := readConfigMap(t, path)
	oauth := cfg["profiles"].(map[string]any)["default"].(map[string]any)["oauth"].(map[string]any)
	oauth["expires_at"] = time.Now().Add(time.Hour).UTC().Format(time.RFC3339)
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	token, err := ProfileAccessToken(context.Background(), "default")
	if err != nil {
		t.Fatal(err)
	}
	if token != "old-access-token" || tokenRequests.Load() != 0 {
		t.Fatalf("expected current token without refresh, got %q after %d refreshes", token, tokenRequests.Load())
	}
}

func TestProfileAccessTokenErrors(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	ts := tokenServer(t, &tokenRequests, rotateTo("unused", "unused"))
	writeOAuthProfileConfig(t, ts.URL, "")

	if _, err := ProfileAccessToken(context.Background(), "missing"); err == nil {
		t.Fatal("expected error for missing profile")
	}
	if _, err := ProfileAccessToken(context.Background(), "other"); err == nil || !strings.Contains(err.Error(), "OAuth") {
		t.Fatalf("expected no-OAuth error for API key profile, got %v", err)
	}
}

// Many processes that start after the access token expires must share one
// rotation. Replaying the rotated refresh token makes the server revoke every
// token for the identity, which forces the user to log in again.
func TestProfileAccessTokenConcurrentCallersRotateOnce(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	var mu sync.Mutex
	valid := "old-refresh-token"
	ts := tokenServer(t, &tokenRequests, func(w http.ResponseWriter, r *http.Request) {
		mu.Lock()
		defer mu.Unlock()
		if r.FormValue("refresh_token") != valid {
			w.WriteHeader(http.StatusBadRequest)
			_ = json.NewEncoder(w).Encode(oauthErrorResponse{Code: "invalid_grant", ErrorDescription: "refresh token has been revoked"})
			return
		}
		valid = "new-refresh-token"
		time.Sleep(20 * time.Millisecond)
		rotateTo("new-access-token", valid)(w, r)
	})
	writeOAuthProfileConfig(t, ts.URL, "")

	const callers = 8
	var wg sync.WaitGroup
	errs := make(chan error, callers)
	for range callers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			token, err := ProfileAccessToken(context.Background(), "")
			if err == nil && token != "new-access-token" {
				err = errors.New("unexpected token " + token)
			}
			errs <- err
		}()
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		if err != nil {
			t.Fatal(err)
		}
	}
	if got := tokenRequests.Load(); got != 1 {
		t.Fatalf("expected exactly one refresh, got %d", got)
	}
}

func TestExpiredTokenWithRejectedRefreshReturnsReauthError(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	ts := tokenServer(t, &tokenRequests, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusBadRequest)
		_ = json.NewEncoder(w).Encode(oauthErrorResponse{Code: "invalid_grant", ErrorDescription: "refresh token has been revoked"})
	})
	writeOAuthProfileConfig(t, ts.URL, "")

	var out map[string]string
	err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out,
		append(loadProfileOptions(), option.WithMaxRetries(2))...)
	var authErr *ProfileAuthError
	if !errors.As(err, &authErr) {
		t.Fatalf("expected ProfileAuthError, got %v", err)
	}
	if authErr.Profile != "default" || !strings.Contains(err.Error(), "langsmith auth login --profile default") {
		t.Fatalf("expected reauth hint for profile, got %q", err.Error())
	}
	// Retries must not replay a refresh token the server already rejected.
	if got := tokenRequests.Load(); got != 1 {
		t.Fatalf("expected one refresh attempt across retries, got %d", got)
	}
	if _, err := ProfileAccessToken(context.Background(), ""); !errors.As(err, &authErr) {
		t.Fatalf("expected ProfileAuthError from ProfileAccessToken, got %v", err)
	}
}

func TestUnexpiredTokenIsUsedWhenRefreshFails(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	ts := tokenServer(t, &tokenRequests, func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusServiceUnavailable)
	})
	path := writeOAuthProfileConfig(t, ts.URL, "")
	cfg := readConfigMap(t, path)
	oauth := cfg["profiles"].(map[string]any)["default"].(map[string]any)["oauth"].(map[string]any)
	// Inside the refresh leeway, but not yet expired.
	oauth["expires_at"] = time.Now().Add(30 * time.Second).UTC().Format(time.RFC3339)
	data, _ := json.Marshal(cfg)
	if err := os.WriteFile(path, data, 0600); err != nil {
		t.Fatal(err)
	}

	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, loadProfileOptions()...); err != nil {
		t.Fatal(err)
	}
	if out["auth"] != "Bearer old-access-token" {
		t.Fatalf("expected still-valid token to be used, got %q", out["auth"])
	}
}

func TestExplicitBearerBeatsDefaultProfile(t *testing.T) {
	clearAuthEnv(t)
	var tokenRequests atomic.Int32
	ts := tokenServer(t, &tokenRequests, rotateTo("profile-access-token", "profile-refresh-token"))
	writeOAuthProfileConfig(t, ts.URL, "")

	explicit := jwtWithSubject(t, "explicit-user")
	var gotUserID string
	opts := append(loadProfileOptions(),
		option.WithHeader("Authorization", "Bearer "+explicit),
		option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			gotUserID = req.Header.Get("X-User-Id")
			return next(req)
		}),
	)
	var out map[string]string
	if err := requestconfig.ExecuteNewRequest(context.Background(), http.MethodGet, "/info", nil, &out, opts...); err != nil {
		t.Fatal(err)
	}
	if out["auth"] != "Bearer "+explicit {
		t.Fatalf("expected explicit bearer, got %q", out["auth"])
	}
	if gotUserID != "explicit-user" {
		t.Fatalf("expected X-User-Id from explicit bearer, got %q", gotUserID)
	}
	if got := tokenRequests.Load(); got != 0 {
		t.Fatalf("expected no refresh of the unused default profile, got %d", got)
	}
}
