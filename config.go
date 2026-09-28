package langsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	authpkg "github.com/langchain-ai/langsmith-go/internal/auth"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

const (
	oauthClientID       = "langsmith-cli"
	tokenRefreshLeeway  = time.Minute
	tokenRefreshTimeout = 10 * time.Second
)

// configProfile holds per-profile configuration from ~/.langsmith/config.json.
// A profile uses api_key (X-API-Key header) or OAuth access_token
// (Authorization: Bearer header) for authentication. access_token is written
// by `langsmith login` under the profile's oauth object.
type configProfile struct {
	APIKey      string      `json:"api_key,omitempty"`
	APIURL      string      `json:"api_url,omitempty"`
	WorkspaceID string      `json:"workspace_id,omitempty"`
	OAuth       configOAuth `json:"oauth,omitempty"`
}

type configOAuth struct {
	AccessToken  string `json:"access_token,omitempty"`
	RefreshToken string `json:"refresh_token,omitempty"`
	ExpiresAt    string `json:"expires_at,omitempty"`
	// Issuer is the authorization server that issued the tokens, recorded by
	// `langsmith auth login`. It can differ from APIURL on BYOC deployments
	// that delegate OAuth to another host, and refreshes must go there.
	Issuer string `json:"issuer,omitempty"`
}

type configFile struct {
	CurrentProfile string                   `json:"current_profile,omitempty"`
	Profiles       map[string]configProfile `json:"profiles,omitempty"`
}

type profileState struct {
	path        string
	cfg         configFile
	profileName string
}

type profileAuth struct {
	state    *profileState
	override bool
	mu       sync.Mutex
	// rejected is the last refresh token the server refused, so SDK retries do
	// not replay it. A new login writes a different token and clears the block.
	rejected string
}

type oauthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type oauthErrorResponse struct {
	Code             string `json:"error"`
	ErrorDescription string `json:"error_description"`
}

func (e oauthErrorResponse) Error() string {
	if e.ErrorDescription == "" {
		return e.Code
	}
	return e.Code + ": " + e.ErrorDescription
}

// loadProfileOptions reads ~/.langsmith/config.json and returns RequestOptions
// for the active profile. Returns nil if no config file exists or no matching
// profile is found.
//
// Profile selection priority:
//  1. LANGSMITH_PROFILE environment variable
//  2. current_profile key in config file
//  3. "default" profile
func loadProfileOptions() []option.RequestOption {
	opts, _ := loadProfileOptionsForProfile("", false)
	return opts
}

func loadProfileOptionsForProfile(profileName string, explicit bool) ([]option.RequestOption, error) {
	state, err := loadProfileState(profileName, explicit)
	if err != nil {
		return nil, err
	}
	if state == nil {
		return nil, nil
	}
	p, ok := state.cfg.Profiles[state.profileName]
	if !ok {
		return nil, nil
	}

	envAuthSet := os.Getenv("LANGSMITH_API_KEY") != ""
	envTenantSet := os.Getenv("LANGSMITH_TENANT_ID") != "" || os.Getenv("LANGSMITH_WORKSPACE_ID") != ""
	envEndpointSet := os.Getenv("LANGSMITH_ENDPOINT") != ""
	hasOAuth := p.OAuth.AccessToken != "" || p.OAuth.RefreshToken != ""

	var opts []option.RequestOption
	if !envEndpointSet {
		if p.APIURL != "" {
			opts = append(opts, option.WithBaseURL(p.APIURL))
		} else if explicit {
			opts = append(opts, resetBaseURL())
		}
	}
	if !envAuthSet {
		switch {
		case hasOAuth:
			opts = append(opts, withProfileAuth(&profileAuth{state: state, override: explicit}))
		case p.APIKey != "":
			opts = append(opts, option.WithAPIKey(p.APIKey))
		case explicit:
			opts = append(opts, resetAPIKey())
		}
	}
	if !envTenantSet {
		if p.WorkspaceID != "" {
			opts = append(opts, option.WithTenantID(p.WorkspaceID))
		} else if explicit {
			opts = append(opts, resetTenantID())
		}
	}
	return opts, nil
}

func loadProfileState(profileName string, explicit bool) (*profileState, error) {
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		if explicit {
			return nil, fmt.Errorf("reading LangSmith config: %w", err)
		}
		return nil, nil
	}

	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		if explicit {
			return nil, fmt.Errorf("parsing LangSmith config: %w", err)
		}
		return nil, nil
	}

	if profileName == "" {
		profileName = resolveProfileName(cfg)
	}
	if profileName == "" {
		return nil, nil
	}

	if _, ok := cfg.Profiles[profileName]; !ok {
		if explicit {
			return nil, fmt.Errorf("LangSmith profile not found: %s", profileName)
		}
		return nil, nil
	}
	return &profileState{path: path, cfg: cfg, profileName: profileName}, nil
}

// reset* clear a value an earlier option set, so an explicitly selected profile
// that omits a dimension replaces current_profile's rather than inheriting it.
func resetBaseURL() option.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.BaseURL = nil
		return nil
	})
}

func resetAPIKey() option.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.APIKey = ""
		r.Request.Header.Del("X-API-Key")
		return nil
	})
}

func resetTenantID() option.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.TenantID = ""
		r.Request.Header.Del("X-Tenant-Id")
		return nil
	})
}

// WithProfile returns a request option that uses a named profile from the
// LangSmith config file. It is equivalent to selecting LANGSMITH_PROFILE for a
// single client without mutating process-wide environment variables.
func WithProfile(profileName string) option.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		opts, err := loadProfileOptionsForProfile(profileName, true)
		if err != nil {
			return err
		}
		return r.Apply(opts...)
	})
}

func withProfileAuth(auth *profileAuth) option.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		name, value, token := auth.currentAuthHeader()
		useProfileAuth := name != "" && (r.APIKey == "" || auth.override)
		if useProfileAuth && token != "" {
			r.OAuthAccessToken = token
			authpkg.SetUserIDHeaderFromAccessToken(r.Request.Header, token)
		}
		if useProfileAuth {
			if auth.override && !strings.EqualFold(name, "X-API-Key") {
				r.APIKey = ""
				r.Request.Header.Del("X-API-Key")
			}
			r.Request.Header.Set(name, value)
		}
		return r.Apply(option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
			if req.Header.Get("X-API-Key") != "" && !auth.override {
				req.Header.Del("Authorization")
				req.Header.Del("X-User-Id")
				return next(req)
			}
			name, value, token, err := auth.authHeader(req.Context())
			if err != nil {
				return nil, err
			}
			if name != "" {
				if auth.override && !strings.EqualFold(name, "X-API-Key") {
					req.Header.Del("X-API-Key")
				}
				if strings.EqualFold(name, "X-API-Key") {
					req.Header.Del("Authorization")
					req.Header.Del("X-User-Id")
				}
				req.Header.Set(name, value)
				if token != "" {
					authpkg.SetUserIDHeaderFromAccessToken(req.Header, token)
				}
			}
			return next(req)
		}))
	})
}

// ProfileAuthError reports that a profile's OAuth access token has expired and
// could not be refreshed. When the server rejected the refresh token, the user
// has to log in again.
type ProfileAuthError struct {
	Profile string
	Err     error
}

func (e *ProfileAuthError) Error() string {
	return fmt.Sprintf(
		"refreshing OAuth token for LangSmith profile %q: %v; run 'langsmith auth login --profile %s' to reauthenticate",
		e.Profile, e.Err, e.Profile,
	)
}

func (e *ProfileAuthError) Unwrap() error { return e.Err }

// ProfileAccessToken returns a current OAuth access token for a profile in the
// LangSmith config file, refreshing and saving it first when it is missing or
// about to expire. An empty profileName selects the active profile
// (LANGSMITH_PROFILE, then current_profile, then "default").
//
// Clients built with [WithProfile] or the default profile already refresh on
// their own; use this only when a bearer token is needed outside the client.
// Refreshes are serialized across processes, so concurrent callers share one
// rotation of the single-use refresh token.
func ProfileAccessToken(ctx context.Context, profileName string) (string, error) {
	state, err := loadProfileState(profileName, true)
	if err != nil {
		return "", err
	}
	if state == nil {
		return "", fmt.Errorf("no LangSmith profile selected")
	}
	p := state.cfg.Profiles[state.profileName]
	if p.OAuth.AccessToken == "" && p.OAuth.RefreshToken == "" {
		return "", fmt.Errorf("LangSmith profile %q has no OAuth credentials; run 'langsmith auth login --profile %s'", state.profileName, state.profileName)
	}
	p, err = (&profileAuth{state: state}).currentProfile(ctx)
	if err != nil {
		return "", err
	}
	return p.OAuth.AccessToken, nil
}

func (a *profileAuth) currentAuthHeader() (name string, value string, token string) {
	a.mu.Lock()
	defer a.mu.Unlock()
	p, ok := a.state.cfg.Profiles[a.state.profileName]
	if !ok {
		return "", "", ""
	}
	return currentAuthHeaderFromProfile(p)
}

func (a *profileAuth) authHeader(ctx context.Context) (name string, value string, token string, err error) {
	p, err := a.currentProfile(ctx)
	if err != nil {
		return "", "", "", err
	}
	name, value, token = authHeaderFromProfile(p)
	return name, value, token, nil
}

// currentProfile returns the profile with a usable access token, refreshing it
// first when it is missing or about to expire.
func (a *profileAuth) currentProfile(ctx context.Context) (configProfile, error) {
	a.mu.Lock()
	defer a.mu.Unlock()
	p, ok := a.state.cfg.Profiles[a.state.profileName]
	if !ok {
		return configProfile{}, fmt.Errorf("LangSmith profile not found: %s", a.state.profileName)
	}
	if !shouldRefreshProfileToken(p) {
		return p, nil
	}

	cfg, fresh, err := refreshProfileLocked(ctx, a.state.path, a.state.profileName, a.rejected)
	if cfg.Profiles != nil {
		a.state.cfg = cfg
	}
	if err == nil {
		return fresh, nil
	}
	var oauthErr oauthErrorResponse
	if errors.As(err, &oauthErr) && oauthErr.Code == "invalid_grant" {
		a.rejected = fresh.OAuth.RefreshToken
	}
	if fresh.OAuth.RefreshToken != "" {
		p = fresh
	}
	// A token inside the refresh leeway still works, so a failed refresh only
	// matters once it has actually expired.
	if accessTokenUsable(p, time.Now()) {
		return p, nil
	}
	return p, &ProfileAuthError{Profile: a.state.profileName, Err: err}
}

var errRefreshTokenRejected = oauthErrorResponse{Code: "invalid_grant", ErrorDescription: "refresh token was already rejected"}

// refreshProfileLocked refreshes profileName's OAuth token while holding a lock
// file next to the config. Refresh tokens are single use: the server rotates
// them, and replaying a rotated one is treated as theft and revokes every token
// for the identity. Re-reading the config after taking the lock lets a caller
// that waited reuse the rotation another process just saved.
//
// The returned profile is the one read under the lock, even on error.
func refreshProfileLocked(ctx context.Context, path, profileName, rejected string) (configFile, configProfile, error) {
	ctx, cancel := context.WithTimeout(ctx, tokenRefreshTimeout)
	defer cancel()

	lock, err := acquireOAuthRefreshLock(ctx, path+".oauth.lock")
	if err != nil {
		return configFile{}, configProfile{}, err
	}
	defer lock.Unlock()

	cfg, p, err := loadProfileConfigFromPath(path, profileName)
	if err != nil {
		return configFile{}, configProfile{}, err
	}
	if !shouldRefreshProfileToken(p) {
		return cfg, p, nil
	}
	if rejected != "" && p.OAuth.RefreshToken == rejected {
		return cfg, p, errRefreshTokenRejected
	}

	tokenURL := p.APIURL
	if p.OAuth.Issuer != "" {
		tokenURL = p.OAuth.Issuer
	}
	token, err := refreshOAuthToken(ctx, tokenURL, p.OAuth.RefreshToken)
	if err != nil {
		return cfg, p, err
	}
	applyTokenResponse(&p, token, time.Now())
	if err := saveProfileOAuth(path, profileName, p.OAuth); err != nil {
		return cfg, p, fmt.Errorf("saving refreshed OAuth token: %w", err)
	}
	cfg.Profiles[profileName] = p
	return cfg, p, nil
}

func loadProfileConfigFromPath(path, profileName string) (configFile, configProfile, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return configFile{}, configProfile{}, err
	}

	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return configFile{}, configProfile{}, err
	}
	p, ok := cfg.Profiles[profileName]
	if !ok {
		return configFile{}, configProfile{}, fmt.Errorf("LangSmith profile not found: %s", profileName)
	}
	return cfg, p, nil
}

func currentAuthHeaderFromProfile(p configProfile) (name string, value string, token string) {
	if p.OAuth.AccessToken != "" {
		return "Authorization", "Bearer " + p.OAuth.AccessToken, p.OAuth.AccessToken
	}
	if p.OAuth.RefreshToken != "" {
		return "", "", ""
	}
	return authHeaderFromProfile(p)
}

func authHeaderFromProfile(p configProfile) (name string, value string, token string) {
	if p.OAuth.AccessToken != "" {
		return "Authorization", "Bearer " + p.OAuth.AccessToken, p.OAuth.AccessToken
	}
	if p.APIKey != "" {
		return "X-API-Key", p.APIKey, ""
	}
	return "", "", ""
}

// resolveProfileName determines which profile to use.
func resolveProfileName(cfg configFile) string {
	// 1. LANGSMITH_PROFILE env var
	if name, ok := os.LookupEnv("LANGSMITH_PROFILE"); ok && name != "" {
		return name
	}
	// 2. current_profile from config
	if cfg.CurrentProfile != "" {
		return cfg.CurrentProfile
	}
	// 3. "default" if it exists
	if _, ok := cfg.Profiles["default"]; ok {
		return "default"
	}
	return ""
}

// configPath returns the path to the config file.
func configPath() string {
	if v := os.Getenv("LANGSMITH_CONFIG_FILE"); v != "" {
		return v
	}
	home, err := os.UserHomeDir()
	if err != nil {
		return ""
	}
	return filepath.Join(home, ".langsmith", "config.json")
}

func shouldRefreshProfileToken(p configProfile) bool {
	if p.OAuth.RefreshToken == "" {
		return false
	}
	if p.OAuth.AccessToken == "" {
		return true
	}
	expiresAt, err := time.Parse(time.RFC3339, p.OAuth.ExpiresAt)
	if err != nil {
		return false
	}
	return !expiresAt.After(time.Now().Add(tokenRefreshLeeway))
}

// accessTokenUsable reports whether p has an access token that has not expired.
// A token without a parseable expiry is assumed usable.
func accessTokenUsable(p configProfile, now time.Time) bool {
	if p.OAuth.AccessToken == "" {
		return false
	}
	expiresAt, err := time.Parse(time.RFC3339, p.OAuth.ExpiresAt)
	return err != nil || expiresAt.After(now)
}

func refreshOAuthToken(ctx context.Context, apiURL, refreshToken string) (*oauthTokenResponse, error) {
	if apiURL == "" {
		apiURL = "https://api.smith.langchain.com"
	}
	values := url.Values{
		"grant_type":    {"refresh_token"},
		"client_id":     {oauthClientID},
		"refresh_token": {refreshToken},
	}
	req, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		resolveTokenEndpoint(ctx, apiURL),
		bytes.NewBufferString(values.Encode()),
	)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		var oauthErr oauthErrorResponse
		if err := json.Unmarshal(body, &oauthErr); err == nil && oauthErr.Code != "" {
			return nil, oauthErr
		}
		return nil, fmt.Errorf("HTTP %d: %s", resp.StatusCode, string(body))
	}
	var token oauthTokenResponse
	if err := json.Unmarshal(body, &token); err != nil {
		return nil, err
	}
	if token.AccessToken == "" {
		return nil, fmt.Errorf("oauth refresh response missing access_token")
	}
	return &token, nil
}

func applyTokenResponse(p *configProfile, token *oauthTokenResponse, now time.Time) {
	p.OAuth.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		p.OAuth.RefreshToken = token.RefreshToken
	}
	if token.ExpiresIn > 0 {
		p.OAuth.ExpiresAt = now.Add(time.Duration(token.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
	}
}

// saveProfileOAuth writes a profile's rotated tokens back to the config file.
// It patches the raw JSON rather than re-encoding configFile so fields the SDK
// does not model survive, and replaces the file atomically so a concurrent
// reader never sees it truncated.
func saveProfileOAuth(path, profileName string, oauth configOAuth) error {
	data, err := os.ReadFile(path)
	if err != nil {
		return err
	}
	var root map[string]json.RawMessage
	if err := json.Unmarshal(data, &root); err != nil {
		return err
	}
	var profiles map[string]map[string]json.RawMessage
	if err := json.Unmarshal(root["profiles"], &profiles); err != nil {
		return err
	}
	profile, ok := profiles[profileName]
	if !ok {
		return fmt.Errorf("LangSmith profile not found: %s", profileName)
	}
	fields := map[string]json.RawMessage{}
	if raw, ok := profile["oauth"]; ok {
		if err := json.Unmarshal(raw, &fields); err != nil {
			return err
		}
	}
	for key, value := range map[string]string{
		"access_token":  oauth.AccessToken,
		"refresh_token": oauth.RefreshToken,
		"expires_at":    oauth.ExpiresAt,
	} {
		if value == "" {
			delete(fields, key)
			continue
		}
		encoded, err := json.Marshal(value)
		if err != nil {
			return err
		}
		fields[key] = encoded
	}

	if profile["oauth"], err = json.Marshal(fields); err != nil {
		return err
	}
	if root["profiles"], err = json.Marshal(profiles); err != nil {
		return err
	}
	out, err := json.MarshalIndent(root, "", "  ")
	if err != nil {
		return err
	}
	return writeFileAtomic(path, append(out, '\n'), 0600)
}

func writeFileAtomic(path string, data []byte, perm os.FileMode) error {
	f, err := os.CreateTemp(filepath.Dir(path), ".langsmith-config-*")
	if err != nil {
		return err
	}
	tmp := f.Name()
	defer os.Remove(tmp)
	if _, err := f.Write(data); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Chmod(perm); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Sync(); err != nil {
		_ = f.Close()
		return err
	}
	if err := f.Close(); err != nil {
		return err
	}
	return os.Rename(tmp, path)
}

func normalizeConfigURL(apiURL string) string {
	return strings.TrimRight(normalizeBaseURL(strings.TrimRight(apiURL, "/")), "/")
}
