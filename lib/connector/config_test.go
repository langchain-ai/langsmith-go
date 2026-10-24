package connector

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

const testConnectionID = "00000000-0000-7000-8000-000000000001"

func validSettings(t *testing.T) Settings {
	t.Helper()
	return Settings{
		Endpoint:     "https://api.example.com",
		ConnectionID: testConnectionID,
		APIKeyFile:   writeAPIKeyFile(t, "initial-placeholder\n"),
		Targets:      map[string]string{"database": "db.internal:5432"},
	}
}

func TestNewConfigAcceptsValidSettings(t *testing.T) {
	t.Parallel()

	config, err := NewConfig(validSettings(t))
	require.NoError(t, err)
	require.Equal(t, "wss://api.example.com/v1/connections/connect", config.Endpoint)
	require.Equal(t, uuid.MustParse(testConnectionID), config.ConnectionID)
	require.Equal(t, "db.internal:5432", mustLookupTarget(t, config.Targets, "database"))
}

func TestNewConfigDefaultsToSaaS(t *testing.T) {
	t.Parallel()

	settings := validSettings(t)
	settings.Endpoint = ""
	config, err := NewConfig(settings)
	require.NoError(t, err)
	require.Equal(t, "wss://api.smith.langchain.com/v1/connections/connect", config.Endpoint)
}

func TestNewConfigAcceptsInlineAPIKey(t *testing.T) {
	t.Parallel()

	settings := validSettings(t)
	settings.APIKeyFile = ""
	settings.APIKey = " inline-placeholder "
	config, err := NewConfig(settings)
	require.NoError(t, err)
	apiKey, err := config.ReadAPIKey()
	require.NoError(t, err)
	require.Equal(t, "inline-placeholder", apiKey)
}

func TestNewConfigRejectsInvalidSettings(t *testing.T) {
	t.Parallel()

	tests := map[string]func(*Settings){
		"no connection ID":      func(s *Settings) { s.ConnectionID = "" },
		"invalid connection ID": func(s *Settings) { s.ConnectionID = "not-a-uuid" },
		"nil connection ID":     func(s *Settings) { s.ConnectionID = uuid.Nil.String() },
		"no credential":         func(s *Settings) { s.APIKeyFile = "" },
		"both credentials":      func(s *Settings) { s.APIKey = "inline-placeholder" },
		"no targets":            func(s *Settings) { s.Targets = nil },
		"empty targets":         func(s *Settings) { s.Targets = map[string]string{} },
		"relay rules without relay mode": func(s *Settings) {
			s.RelayDeny = []string{"*.example.com"}
		},
	}
	for name, mutate := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := validSettings(t)
			mutate(&settings)
			_, err := NewConfig(settings)
			require.Error(t, err)
		})
	}
}

func TestNewConfigRejectsInvalidEndpoint(t *testing.T) {
	t.Parallel()

	endpoints := []string{
		// A WebSocket URL is the derived form, not something callers configure.
		"wss://api.example.com/v1/connections/connect",
		"ws://api.example.com",
		"https://",
		"ftp://api.example.com",
		"api.example.com",
		"https://user:password@api.example.com",
		"https://api.example.com?region=test",
		"https://api.example.com?",
		"https://api.example.com#fragment",
	}
	for _, endpoint := range endpoints {
		t.Run(endpoint, func(t *testing.T) {
			t.Parallel()

			settings := validSettings(t)
			settings.Endpoint = endpoint
			_, err := NewConfig(settings)
			require.Error(t, err)
		})
	}
}

func TestNewConfigRejectsUnreadableOrEmptyAPIKeyFile(t *testing.T) {
	t.Parallel()

	tests := map[string]string{
		"missing":    filepath.Join(t.TempDir(), "missing"),
		"empty":      writeAPIKeyFile(t, ""),
		"whitespace": writeAPIKeyFile(t, " \n\t"),
	}
	for name, apiKeyFile := range tests {
		t.Run(name, func(t *testing.T) {
			t.Parallel()

			settings := validSettings(t)
			settings.APIKeyFile = apiKeyFile
			_, err := NewConfig(settings)
			require.Error(t, err)
		})
	}
}

func TestReadAPIKeyReadsFileOnEveryCall(t *testing.T) {
	t.Parallel()

	settings := validSettings(t)
	config, err := NewConfig(settings)
	require.NoError(t, err)

	require.NoError(t, os.WriteFile(settings.APIKeyFile, []byte("rotated-placeholder\n"), 0o600))
	apiKey, err := config.ReadAPIKey()
	require.NoError(t, err)
	require.Equal(t, "rotated-placeholder", apiKey)
}

// TestCarrierEndpointDerivesTunnelURLFromAPIURL: callers configure the same
// LangSmith API URL they use everywhere else; the WebSocket scheme and carrier
// path are internal details supplied here. A base path is preserved so a
// self-hosted install mounted under a prefix keeps working.
func TestCarrierEndpointDerivesTunnelURLFromAPIURL(t *testing.T) {
	t.Parallel()

	for _, test := range []struct{ endpoint, want string }{
		{"https://api.smith.langchain.com", "wss://api.smith.langchain.com/v1/connections/connect"},
		{"https://api.smith.langchain.com/", "wss://api.smith.langchain.com/v1/connections/connect"},
		{"https://langsmith.corp.example.com/api", "wss://langsmith.corp.example.com/api/v1/connections/connect"},
		{"https://langsmith.corp.example.com/api/", "wss://langsmith.corp.example.com/api/v1/connections/connect"},
		{"https://langsmith.corp.example.com:8443", "wss://langsmith.corp.example.com:8443/v1/connections/connect"},
		// Plain HTTP is for local development only, and maps to ws.
		{"http://localhost:1984", "ws://localhost:1984/v1/connections/connect"},
		{"http://127.0.0.1:1984", "ws://127.0.0.1:1984/v1/connections/connect"},
		{"http://[::1]:1984", "ws://[::1]:1984/v1/connections/connect"},
	} {
		t.Run(test.endpoint, func(t *testing.T) {
			t.Parallel()

			got, err := carrierEndpoint(test.endpoint)
			require.NoError(t, err)
			require.Equal(t, test.want, got)
		})
	}
}

func TestCarrierEndpointRejectsPlainHTTPForNonLoopbackHosts(t *testing.T) {
	t.Parallel()

	for _, endpoint := range []string{
		"http://langsmith.internal:1984",
		"http://10.0.0.10:1984",
		"http://169.254.169.254:1984",
	} {
		t.Run(endpoint, func(t *testing.T) {
			t.Parallel()

			_, err := carrierEndpoint(endpoint)
			require.Error(t, err)
		})
	}
}

func TestValidateCarrierURLRejectsPlaintextNonLoopbackHost(t *testing.T) {
	t.Parallel()

	err := validateCarrierURL("ws://langsmith.internal/v1/connections/connect")
	require.Error(t, err)
}

func TestCarrierEndpointRejectsSurroundingWhitespace(t *testing.T) {
	t.Parallel()

	for _, endpoint := range []string{" https://api.example.com", "https://api.example.com\n"} {
		_, err := carrierEndpoint(endpoint)
		require.Error(t, err)
	}
}

func TestNewClientAcceptsAnInlineAPIKey(t *testing.T) {
	t.Parallel()

	config, err := NewConfig(Settings{
		APIKey:       "lsv2_placeholder",
		ConnectionID: testConnectionID,
		Targets:      map[string]string{"database": "db.internal:5432"},
	})
	require.NoError(t, err)

	client, err := NewClient(config)
	require.NoError(t, err)
	require.NotNil(t, client)
}

func mustLookupTarget(t *testing.T, targets TargetMap, name string) string {
	t.Helper()

	address, ok := targets.Lookup(name)
	require.True(t, ok)
	return address
}

func writeAPIKeyFile(t *testing.T, contents string) string {
	t.Helper()

	path := filepath.Join(t.TempDir(), "api-key")
	require.NoError(t, os.WriteFile(path, []byte(contents), 0o600))
	return path
}
