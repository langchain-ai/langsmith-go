package cmdutil

import (
	"bytes"
	"encoding/json"
	"os"
	"path/filepath"
	"testing"

	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	lsconfig "github.com/langchain-ai/langsmith-go/cli/internal/config"
	"github.com/stretchr/testify/require"
)

func writeProfiles(t *testing.T, cfg string) {
	t.Helper()
	path := filepath.Join(t.TempDir(), "config.json")
	require.NoError(t, os.WriteFile(path, []byte(cfg), 0o600))
	t.Setenv(lsconfig.ConfigFileEnv, path)
	t.Setenv("LANGSMITH_PROFILE", "")
}

func captureNotices(t *testing.T) *bytes.Buffer {
	t.Helper()
	var buf bytes.Buffer
	ResetNotices(&buf)
	t.Cleanup(func() { ResetNotices(os.Stderr) })
	return &buf
}

func TestNoteActiveProfile(t *testing.T) {
	const two = `{"current_profile":"work","profiles":{"work":{"api_key":"k1","workspace_id":"ws-1"},"personal":{"api_key":"k2"}}}`

	t.Run("names the profile and workspace once per run", func(t *testing.T) {
		writeProfiles(t, two)
		buf := captureNotices(t)
		opts := client.Options{ProfileName: "work", WorkspaceID: "ws-1"}
		NoteActiveProfile("", opts, false)
		NoteActiveProfile("", opts, false)
		require.Equal(t, "note: using profile \"work\", workspace ws-1 (2 profiles saved; pass --profile or --workspace to change)\n", buf.String())
	})

	t.Run("json is one object per line", func(t *testing.T) {
		writeProfiles(t, two)
		buf := captureNotices(t)
		NoteActiveProfile("personal", client.Options{ProfileName: "personal"}, true)
		var got map[string]any
		require.NoError(t, json.Unmarshal(buf.Bytes(), &got))
		require.Equal(t, "active_profile", got["notice"])
		require.Equal(t, "personal", got["profile"])
		require.Equal(t, "", got["workspace_id"])
		require.Contains(t, got["message"], "the API key's default")
	})

	t.Run("silent with a single profile", func(t *testing.T) {
		writeProfiles(t, `{"current_profile":"work","profiles":{"work":{"api_key":"k1"}}}`)
		buf := captureNotices(t)
		NoteActiveProfile("", client.Options{ProfileName: "work"}, false)
		require.Empty(t, buf.String())
	})
}
