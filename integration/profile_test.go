//go:build integration

package integration

import (
	"context"
	"os"
	"path/filepath"
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

// newClientFromProfile creates a client that resolves credentials from a
// temporary config profile rather than environment variables.
func newClientFromProfile(t *testing.T, profileName, configContent string) *langsmith.Client {
	t.Helper()
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	if err := os.WriteFile(path, []byte(configContent), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	if profileName != "" {
		t.Setenv("LANGSMITH_PROFILE", profileName)
	} else {
		os.Unsetenv("LANGSMITH_PROFILE")
	}
	// Unset env vars so the profile is the only credential source.
	// os.Unsetenv is needed because t.Setenv("X", "") makes LookupEnv
	// return (true, ""), which the SDK interprets as an explicit empty override.
	envVars := []string{"LANGSMITH_API_KEY", "LANGSMITH_ENDPOINT", "LANGSMITH_TENANT_ID"}
	saved := make(map[string]string)
	for _, k := range envVars {
		if v, ok := os.LookupEnv(k); ok {
			saved[k] = v
		}
		os.Unsetenv(k)
	}
	t.Cleanup(func() {
		for k, v := range saved {
			os.Setenv(k, v)
		}
	})
	return langsmith.NewClient()
}

// TestProfileClient_ListSessions verifies that a client created entirely from
// a config profile (no env vars) can successfully authenticate and list sessions.
func TestProfileClient_ListSessions(t *testing.T) {
	// We need real credentials — read them from the env BEFORE clearing.
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set, skipping profile integration test")
	}
	endpoint := os.Getenv("LANGSMITH_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.smith.langchain.com"
	}

	config := `{
  "current_profile": "integ",
  "profiles": {
    "integ": {
      "api_key": "` + apiKey + `",
      "api_url": "` + endpoint + `"
    }
  }
}
`
	client := newClientFromProfile(t, "integ", config)

	resp, err := client.Sessions.List(context.Background(), langsmith.SessionListParams{
		Limit: langsmith.F(int64(1)),
	})
	if err != nil {
		t.Fatalf("list sessions via profile: %v", err)
	}
	// We just need to confirm the API accepted the request — at least 0 results is fine.
	_ = resp.Items
}

// TestProfileClient_DatasetCRUD verifies a full CRUD cycle using profile-based auth.
func TestProfileClient_DatasetCRUD(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set, skipping profile integration test")
	}
	endpoint := os.Getenv("LANGSMITH_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.smith.langchain.com"
	}

	config := `{
  "profiles": {
    "default": {
      "api_key": "` + apiKey + `",
      "api_url": "` + endpoint + `"
    }
  }
}
`
	// Uses "default" profile via fallback (no current_profile, no LANGSMITH_PROFILE).
	client := newClientFromProfile(t, "", config)

	ctx := context.Background()
	name := uniqueName("go-integ-profile-dataset")

	// Create
	dataset, err := client.Datasets.New(ctx, langsmith.DatasetNewParams{
		Name:        langsmith.F(name),
		Description: langsmith.F("Profile integration test"),
	})
	if err != nil {
		t.Fatalf("create dataset: %v", err)
	}
	defer client.Datasets.Delete(ctx, dataset.ID)

	if dataset.Name != name {
		t.Errorf("name = %q, want %q", dataset.Name, name)
	}

	// List
	listed, err := client.Datasets.List(ctx, langsmith.DatasetListParams{
		Name: langsmith.F(name),
	})
	if err != nil {
		t.Fatalf("list datasets: %v", err)
	}
	if len(listed.Items) == 0 {
		t.Error("expected at least one dataset")
	}
}

// TestProfileClient_EnvProfileSelection verifies that LANGSMITH_PROFILE env
// var correctly selects the named profile.
func TestProfileClient_EnvProfileSelection(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set, skipping profile integration test")
	}
	endpoint := os.Getenv("LANGSMITH_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.smith.langchain.com"
	}

	config := `{
  "current_profile": "wrong",
  "profiles": {
    "wrong": {
      "api_key": "invalid-key-should-not-be-used",
      "api_url": "` + endpoint + `"
    },
    "correct": {
      "api_key": "` + apiKey + `",
      "api_url": "` + endpoint + `"
    }
  }
}
`
	// LANGSMITH_PROFILE should override current_profile.
	client := newClientFromProfile(t, "correct", config)

	_, err := client.Sessions.List(context.Background(), langsmith.SessionListParams{
		Limit: langsmith.F(int64(1)),
	})
	if err != nil {
		t.Fatalf("expected success with 'correct' profile, got: %v", err)
	}
}

// TestProfileClient_EnvOverridesProfile verifies that environment variables
// take precedence over profile values. The profile has an invalid key, but
// the env var provides the real one.
func TestProfileClient_EnvOverridesProfile(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set, skipping profile integration test")
	}
	endpoint := os.Getenv("LANGSMITH_ENDPOINT")
	if endpoint == "" {
		endpoint = "https://api.smith.langchain.com"
	}

	// Profile has a bad key — env var should override it.
	dir := t.TempDir()
	path := filepath.Join(dir, "config.json")
	config := `{
  "profiles": {
    "default": {
      "api_key": "invalid-key-from-profile",
      "api_url": "` + endpoint + `"
    }
  }
}
`
	if err := os.WriteFile(path, []byte(config), 0600); err != nil {
		t.Fatal(err)
	}
	t.Setenv("LANGSMITH_CONFIG_FILE", path)
	// Set the real key via env — this should override the profile's bad key.
	t.Setenv("LANGSMITH_API_KEY", apiKey)
	t.Setenv("LANGSMITH_ENDPOINT", endpoint)
	os.Unsetenv("LANGSMITH_PROFILE")

	client := langsmith.NewClient()
	_, err := client.Sessions.List(context.Background(), langsmith.SessionListParams{
		Limit: langsmith.F(int64(1)),
	})
	if err != nil {
		t.Fatalf("env var should override profile's bad key, got: %v", err)
	}
}
