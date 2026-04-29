package langsmith

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"time"

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
}

type configFile struct {
	CurrentProfile string                   `json:"current_profile,omitempty"`
	Profiles       map[string]configProfile `json:"profiles,omitempty"`
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
	path := configPath()
	data, err := os.ReadFile(path)
	if err != nil {
		return nil
	}

	var cfg configFile
	if err := json.Unmarshal(data, &cfg); err != nil {
		return nil
	}

	profileName := resolveProfileName(cfg)
	if profileName == "" {
		return nil
	}

	p, ok := cfg.Profiles[profileName]
	if !ok {
		return nil
	}

	refreshURL := p.APIURL
	if envURL := os.Getenv("LANGSMITH_ENDPOINT"); envURL != "" {
		refreshURL = envURL
	}
	envAuthSet := os.Getenv("LANGSMITH_API_KEY") != ""
	if shouldRefreshProfileToken(p) &&
		!envAuthSet {
		ctx, cancel := context.WithTimeout(context.Background(), tokenRefreshTimeout)
		defer cancel()
		if token, err := refreshOAuthToken(ctx, refreshURL, p.OAuth.RefreshToken); err == nil {
			applyTokenResponse(&p, token, time.Now())
			cfg.Profiles[profileName] = p
			_ = saveConfig(path, cfg)
		}
	}

	var opts []option.RequestOption
	if p.APIURL != "" {
		opts = append(opts, option.WithBaseURL(p.APIURL))
	}
	if !envAuthSet {
		if token := p.OAuth.AccessToken; token != "" {
			opts = append(opts, withOAuthAccessToken(token))
		} else if p.APIKey != "" {
			opts = append(opts, option.WithAPIKey(p.APIKey))
		}
	}
	if p.WorkspaceID != "" {
		opts = append(opts, option.WithTenantID(p.WorkspaceID))
	}
	return opts
}

func withOAuthAccessToken(token string) option.RequestOption {
	return requestconfig.RequestOptionFunc(func(r *requestconfig.RequestConfig) error {
		r.OAuthAccessToken = token
		return r.Apply(option.WithHeader("authorization", "Bearer "+token))
	})
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
		strings.TrimRight(normalizeConfigURL(apiURL), "/")+"/oauth/token",
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

func saveConfig(path string, cfg configFile) error {
	data, err := json.MarshalIndent(cfg, "", "  ")
	if err != nil {
		return err
	}
	data = append(data, '\n')
	if err := os.MkdirAll(filepath.Dir(path), 0700); err != nil {
		return err
	}
	if err := os.WriteFile(path, data, 0600); err != nil {
		return err
	}
	return os.Chmod(path, 0600)
}

func normalizeConfigURL(apiURL string) string {
	u := strings.TrimRight(apiURL, "/")
	return strings.TrimSuffix(u, "/api/v1")
}
