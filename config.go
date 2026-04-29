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

	"github.com/BurntSushi/toml"
	"github.com/langchain-ai/langsmith-go/option"
)

const (
	oauthClientID       = "langsmith-cli"
	tokenRefreshLeeway  = time.Minute
	tokenRefreshTimeout = 10 * time.Second
)

// configProfile holds per-profile configuration from ~/.langsmith/config.toml.
// A profile uses api_key (X-API-Key header), bearer_token, or access_token
// (Authorization: Bearer header) for authentication. access_token is written
// by `langsmith login` and takes precedence over bearer_token.
type configProfile struct {
	APIKey         string `toml:"api_key"`
	BearerToken    string `toml:"bearer_token"`
	APIURL         string `toml:"api_url"`
	WorkspaceID    string `toml:"workspace_id"`
	AccessToken    string `toml:"access_token"`
	RefreshToken   string `toml:"refresh_token"`
	TokenType      string `toml:"token_type"`
	TokenExpiresAt string `toml:"token_expires_at"`
}

type oauthTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
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

// loadProfileOptions reads ~/.langsmith/config.toml and returns RequestOptions
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

	var raw map[string]any
	if err := toml.Unmarshal(data, &raw); err != nil {
		return nil
	}

	profileName := resolveProfileName(raw)
	if profileName == "" {
		return nil
	}

	section, ok := raw[profileName].(map[string]any)
	if !ok {
		return nil
	}

	var p configProfile
	if v, ok := section["api_key"].(string); ok {
		p.APIKey = v
	}
	if v, ok := section["api_url"].(string); ok {
		p.APIURL = v
	}
	if v, ok := section["bearer_token"].(string); ok {
		p.BearerToken = v
	}
	if v, ok := section["access_token"].(string); ok {
		p.AccessToken = v
	}
	if v, ok := section["refresh_token"].(string); ok {
		p.RefreshToken = v
	}
	if v, ok := section["token_type"].(string); ok {
		p.TokenType = v
	}
	if v, ok := section["token_expires_at"].(string); ok {
		p.TokenExpiresAt = v
	}
	if v, ok := section["workspace_id"].(string); ok {
		p.WorkspaceID = v
	}
	refreshURL := p.APIURL
	if envURL := os.Getenv("LANGSMITH_ENDPOINT"); envURL != "" {
		refreshURL = envURL
	}
	envAuthSet := os.Getenv("LANGSMITH_API_KEY") != "" || os.Getenv("LANGSMITH_BEARER_TOKEN") != ""
	if shouldRefreshProfileToken(p) &&
		!envAuthSet {
		ctx, cancel := context.WithTimeout(context.Background(), tokenRefreshTimeout)
		defer cancel()
		if token, err := refreshOAuthToken(ctx, refreshURL, p.RefreshToken); err == nil {
			applyTokenResponse(&p, token, time.Now())
			section["access_token"] = p.AccessToken
			section["refresh_token"] = p.RefreshToken
			section["token_type"] = p.TokenType
			section["token_expires_at"] = p.TokenExpiresAt
			_ = saveRawConfig(path, raw)
		}
	}

	var opts []option.RequestOption
	if p.APIURL != "" {
		opts = append(opts, option.WithBaseURL(p.APIURL))
	}
	if !envAuthSet {
		if token := p.authBearerToken(); token != "" {
			opts = append(opts, option.WithBearerToken(token))
		} else if p.APIKey != "" {
			opts = append(opts, option.WithAPIKey(p.APIKey))
		}
	}
	if p.WorkspaceID != "" {
		opts = append(opts, option.WithTenantID(p.WorkspaceID))
	}
	return opts
}

// resolveProfileName determines which profile to use.
func resolveProfileName(raw map[string]any) string {
	// 1. LANGSMITH_PROFILE env var
	if name, ok := os.LookupEnv("LANGSMITH_PROFILE"); ok && name != "" {
		return name
	}
	// 2. current_profile from config
	if cp, ok := raw["current_profile"].(string); ok && cp != "" {
		return cp
	}
	// 3. "default" if it exists
	if _, ok := raw["default"]; ok {
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
	return filepath.Join(home, ".langsmith", "config.toml")
}

func (p configProfile) authBearerToken() string {
	if p.AccessToken != "" {
		return p.AccessToken
	}
	return p.BearerToken
}

func shouldRefreshProfileToken(p configProfile) bool {
	if p.RefreshToken == "" {
		return false
	}
	if p.AccessToken == "" {
		return true
	}
	expiresAt, err := time.Parse(time.RFC3339, p.TokenExpiresAt)
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
	p.AccessToken = token.AccessToken
	if token.RefreshToken != "" {
		p.RefreshToken = token.RefreshToken
	}
	p.TokenType = token.TokenType
	if p.TokenType == "" {
		p.TokenType = "Bearer"
	}
	if token.ExpiresIn > 0 {
		p.TokenExpiresAt = now.Add(time.Duration(token.ExpiresIn) * time.Second).UTC().Format(time.RFC3339)
	}
}

func saveRawConfig(path string, raw map[string]any) error {
	data, err := toml.Marshal(raw)
	if err != nil {
		return err
	}
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
