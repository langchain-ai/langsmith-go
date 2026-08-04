package langsmith

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// selfHostedASServer serves authorization server metadata under /api, as
// self-hosted LangSmith does, and an HTML 200 at the bare-root .well-known path
// like the SPA catch-all.
func selfHostedASServer(t *testing.T) *httptest.Server {
	t.Helper()
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/.well-known/oauth-authorization-server":
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"issuer":                        srv.URL + "/api",
				"device_authorization_endpoint": srv.URL + "/api/oauth/device/code",
				"token_endpoint":                srv.URL + "/api/oauth/token",
			})
		case "/.well-known/oauth-authorization-server":
			w.Header().Set("Content-Type", "text/html")
			_, _ = w.Write([]byte("<!doctype html><html><body>app</body></html>"))
		default:
			http.NotFound(w, r)
		}
	}))
	t.Cleanup(srv.Close)
	return srv
}

func TestResolveTokenEndpoint_SelfHostedBareOrigin(t *testing.T) {
	srv := selfHostedASServer(t)

	got := resolveTokenEndpoint(context.Background(), srv.URL)
	if want := srv.URL + "/api/oauth/token"; got != want {
		t.Errorf("resolveTokenEndpoint(%q) = %q, want %q", srv.URL, got, want)
	}
}

func TestResolveTokenEndpoint_SelfHostedApiSuffix(t *testing.T) {
	srv := selfHostedASServer(t)

	got := resolveTokenEndpoint(context.Background(), srv.URL+"/api")
	if want := srv.URL + "/api/oauth/token"; got != want {
		t.Errorf("resolveTokenEndpoint(%q) = %q, want %q", srv.URL+"/api", got, want)
	}
}

func TestResolveTokenEndpoint_SaaSAtRoot(t *testing.T) {
	var srv *httptest.Server
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/.well-known/oauth-authorization-server" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]any{
				"issuer":                        srv.URL,
				"device_authorization_endpoint": srv.URL + "/oauth/device/code",
				"token_endpoint":                srv.URL + "/oauth/token",
			})
			return
		}
		http.NotFound(w, r)
	}))
	defer srv.Close()

	got := resolveTokenEndpoint(context.Background(), srv.URL)
	if want := srv.URL + "/oauth/token"; got != want {
		t.Errorf("resolveTokenEndpoint = %q, want %q", got, want)
	}
}

// Without a metadata document the endpoint keeps the configured mount point.
// The /api segment must survive here: normalizeBaseURL strips it for REST, but
// the authorization server on self-hosted lives under it.
func TestResolveTokenEndpoint_FallbackKeepsApiMount(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	got := resolveTokenEndpoint(context.Background(), srv.URL+"/api")
	if want := srv.URL + "/api/oauth/token"; got != want {
		t.Errorf("resolveTokenEndpoint = %q, want %q", got, want)
	}
}

func TestResolveTokenEndpoint_FallbackStripsApiV1(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		http.NotFound(w, r)
	}))
	defer srv.Close()

	got := resolveTokenEndpoint(context.Background(), srv.URL+"/api/v1")
	if want := srv.URL + "/oauth/token"; got != want {
		t.Errorf("resolveTokenEndpoint = %q, want %q", got, want)
	}
}

// Credentials are posted to this endpoint, so a document naming a different
// issuer or an off-origin endpoint must be ignored in favour of the fallback.
func TestResolveTokenEndpoint_RejectsUntrustedMetadata(t *testing.T) {
	for _, tc := range []struct {
		name string
		doc  func(base string) map[string]any
	}{
		{"issuer mismatch", func(base string) map[string]any {
			return map[string]any{
				"issuer":                        "https://evil.example.com",
				"device_authorization_endpoint": base + "/oauth/device/code",
				"token_endpoint":                base + "/oauth/token",
			}
		}},
		{"foreign endpoint", func(base string) map[string]any {
			return map[string]any{
				"issuer":                        base,
				"device_authorization_endpoint": base + "/oauth/device/code",
				"token_endpoint":                "https://evil.example.com/oauth/token",
			}
		}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			var srv *httptest.Server
			srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.URL.Path == "/.well-known/oauth-authorization-server" {
					w.Header().Set("Content-Type", "application/json")
					_ = json.NewEncoder(w).Encode(tc.doc(srv.URL))
					return
				}
				http.NotFound(w, r)
			}))
			defer srv.Close()

			got := resolveTokenEndpoint(context.Background(), srv.URL)
			if want := srv.URL + "/oauth/token"; got != want {
				t.Errorf("resolveTokenEndpoint = %q, want fallback %q", got, want)
			}
		})
	}
}
