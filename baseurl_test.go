package langsmith

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestNormalizeBaseURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name string
		in   string
		want string
	}{
		{"cloud api/v1 suffix", "https://api.smith.langchain.com/api/v1", "https://api.smith.langchain.com"},
		{"trailing slash", "https://api.smith.langchain.com/api/v1/", "https://api.smith.langchain.com"},
		{"bare api suffix", "https://api.smith.langchain.com/api", "https://api.smith.langchain.com"},
		{"bare api suffix trailing slash", "https://api.smith.langchain.com/api/", "https://api.smith.langchain.com"},
		{"self-hosted path prefix", "https://self-hosted.example.com/langsmith/api/v1", "https://self-hosted.example.com/langsmith/"},
		{"self-hosted path prefix bare api", "https://self-hosted.example.com/langsmith/api", "https://self-hosted.example.com/langsmith/"},
		{"already normalized", "https://api.smith.langchain.com", "https://api.smith.langchain.com"},
		{"localhost", "http://localhost:1984", "http://localhost:1984"},
		{"localhost with suffix", "http://localhost:1984/api/v1", "http://localhost:1984"},
		// A bare /v1 is not stripped, matching the Python and JS SDKs.
		{"bare v1 suffix", "https://api.smith.langchain.com/v1", "https://api.smith.langchain.com/v1"},
		// /api only counts as a whole trailing path segment.
		{"api mid-path", "https://api.smith.langchain.com/api/runs", "https://api.smith.langchain.com/api/runs"},
		{"api in host", "https://api.example.com", "https://api.example.com"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			if got := normalizeBaseURL(tt.in); got != tt.want {
				t.Errorf("normalizeBaseURL(%q) = %q, want %q", tt.in, got, tt.want)
			}
		})
	}
}

func TestNormalizeConfigURL(t *testing.T) {
	t.Parallel()

	tests := []struct {
		in   string
		want string
	}{
		{"https://api.smith.langchain.com/api/v1", "https://api.smith.langchain.com"},
		{"https://api.smith.langchain.com/api", "https://api.smith.langchain.com"},
		{"https://self-hosted.example.com/langsmith/api/v1/", "https://self-hosted.example.com/langsmith"},
		{"https://api.smith.langchain.com", "https://api.smith.langchain.com"},
	}

	for _, tt := range tests {
		if got := normalizeConfigURL(tt.in); got != tt.want {
			t.Errorf("normalizeConfigURL(%q) = %q, want %q", tt.in, got, tt.want)
		}
	}
}

func TestWithoutBaseURLNormalization(t *testing.T) {
	t.Parallel()

	opts := []option.RequestOption{
		option.WithBaseURL("https://self-hosted.example.com/api"),
		withNormalizedBaseURL(),
		option.WithAPIKey("test-api-key"),
	}

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatalf("http.NewRequest: %v", err)
	}
	cfg := requestconfig.RequestConfig{Request: req}
	if err := cfg.Apply(withoutBaseURLNormalization(opts)...); err != nil {
		t.Fatalf("Apply: %v", err)
	}

	if want := "https://self-hosted.example.com/api/"; cfg.BaseURL.String() != want {
		t.Errorf("base URL = %q, want %q", cfg.BaseURL, want)
	}
	if got, want := len(withoutBaseURLNormalization(opts)), len(opts)-1; got != want {
		t.Errorf("kept %d options, want %d", got, want)
	}
}

// TestClientTracingKeepsConfiguredBaseURL checks that trace ingest uses the base
// URL exactly as configured, without the normalization the generated request paths
// need. The prefix the ingest endpoint lives behind is deployment specific — /api
// on self-hosted, nothing on SaaS — so it cannot be reconstructed from the root.
func TestClientTracingKeepsConfiguredBaseURL(t *testing.T) {
	for _, suffix := range []string{"", "/api", "/api/v1"} {
		t.Run("suffix="+suffix, func(t *testing.T) {
			paths := make(chan string, 4)
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				select {
				case paths <- r.URL.Path:
				default:
				}
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
			}))
			defer server.Close()

			client := NewClient(
				option.WithBaseURL(server.URL+suffix),
				option.WithAPIKey("test-api-key"),
			)
			if err := client.CreateRun(&RunCreate{
				ID:        uuid.New(),
				Name:      "test-run",
				RunType:   "chain",
				StartTime: time.Now(),
			}); err != nil {
				t.Fatalf("CreateRun: %v", err)
			}
			client.Close() // drains the trace sink

			select {
			case got := <-paths:
				if want := suffix + "/runs/multipart"; got != want {
					t.Errorf("ingest path = %q, want %q", got, want)
				}
			case <-time.After(5 * time.Second):
				t.Fatal("no ingest request received")
			}
		})
	}
}

// TestClientNormalizesConfiguredBaseURL checks that a base URL configured with an
// /api or /api/v1 suffix does not double up with the prefix the generated request
// paths already carry.
func TestClientNormalizesConfiguredBaseURL(t *testing.T) {
	for _, suffix := range []string{"", "/api", "/api/v1"} {
		t.Run("suffix="+suffix, func(t *testing.T) {
			var gotPath string
			server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				gotPath = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				w.Write([]byte(`{}`))
			}))
			defer server.Close()

			client := NewClient(
				option.WithBaseURL(server.URL+suffix),
				option.WithAPIKey("test-api-key"),
			)
			_, err := client.Runs.Stats(context.Background(), RunStatsParams{})
			if err != nil {
				t.Fatalf("Runs.Stats: %v", err)
			}
			if want := "/api/v1/runs/stats"; gotPath != want {
				t.Errorf("request path = %q, want %q", gotPath, want)
			}
		})
	}
}
