package langsmith

import (
	"os"
	"path/filepath"

	"github.com/BurntSushi/toml"
	"github.com/langchain-ai/langsmith-go/option"
)

// configProfile holds per-profile configuration from ~/.langsmith/config.toml.
type configProfile struct {
	APIKey      string `toml:"api_key"`
	APIURL      string `toml:"api_url"`
	WorkspaceID string `toml:"workspace_id"`
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
	if v, ok := section["workspace_id"].(string); ok {
		p.WorkspaceID = v
	}

	var opts []option.RequestOption
	if p.APIURL != "" {
		opts = append(opts, option.WithBaseURL(p.APIURL))
	}
	if p.APIKey != "" {
		opts = append(opts, option.WithAPIKey(p.APIKey))
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
