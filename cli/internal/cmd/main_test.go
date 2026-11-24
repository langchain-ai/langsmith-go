package cmd

import (
	"os"
	"path/filepath"
	"testing"

	lsconfig "github.com/langchain-ai/langsmith-go/cli/internal/config"
	"github.com/zalando/go-keyring"
)

// TestMain keeps tests off the developer's real keychain and config file.
// Tests that need saved profiles point LANGSMITH_CONFIG_FILE at their own.
func TestMain(m *testing.M) {
	keyring.MockInit()
	_ = os.Unsetenv(lsconfig.SecretStoreEnv)
	_ = os.Unsetenv("LANGSMITH_PROFILE")
	dir, err := os.MkdirTemp("", "langsmith-cmd-test")
	if err != nil {
		panic(err)
	}
	_ = os.Setenv(lsconfig.ConfigFileEnv, filepath.Join(dir, "config.json"))
	code := m.Run()
	_ = os.RemoveAll(dir)
	os.Exit(code)
}
