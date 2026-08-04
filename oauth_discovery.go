package langsmith

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	wellKnownOAuthPath      = "/.well-known/oauth-authorization-server"
	oauthDiscoveryMaxBody   = 1 << 20
	oauthDiscoveryPerTryTTL = 5 * time.Second
)

type oauthServerMetadata struct {
	Issuer                      string `json:"issuer"`
	DeviceAuthorizationEndpoint string `json:"device_authorization_endpoint"`
	TokenEndpoint               string `json:"token_endpoint"`
}

// oauthConfigURL reduces a configured API URL to the mount point the
// authorization server is served from. Unlike normalizeConfigURL it keeps a
// trailing "/api": REST lives at <origin>/api/v1 so its base URL is the origin,
// but on self-hosted the authorization server lives under <origin>/api.
func oauthConfigURL(apiURL string) string {
	u := strings.TrimRight(apiURL, "/")
	return strings.TrimSuffix(u, "/api/v1")
}

// resolveTokenEndpoint returns the token endpoint for apiURL, preferring the
// RFC 8414 metadata document the deployment serves and falling back to
// <mount>/oauth/token when none is available. It never fails: callers treat an
// unreachable endpoint as a failed refresh.
func resolveTokenEndpoint(ctx context.Context, apiURL string) string {
	for _, base := range oauthDiscoveryCandidates(apiURL) {
		if meta := fetchOAuthMetadata(ctx, base); meta != nil {
			return meta.TokenEndpoint
		}
	}
	return strings.TrimRight(oauthConfigURL(apiURL), "/") + "/oauth/token"
}

// oauthDiscoveryCandidates lists metadata base URLs to probe, most specific
// first: the configured mount point, then the self-hosted and SaaS locations.
func oauthDiscoveryCandidates(apiURL string) []string {
	given := strings.TrimRight(oauthConfigURL(apiURL), "/")
	origin := strings.TrimRight(normalizeConfigURL(apiURL), "/")

	out := make([]string, 0, 3)
	seen := make(map[string]bool, 3)
	for _, c := range []string{given, origin + "/api", origin} {
		if c == "" || c == "/api" || seen[c] {
			continue
		}
		seen[c] = true
		out = append(out, c)
	}
	return out
}

func fetchOAuthMetadata(ctx context.Context, base string) *oauthServerMetadata {
	reqCtx, cancel := context.WithTimeout(ctx, oauthDiscoveryPerTryTTL)
	defer cancel()

	req, err := http.NewRequestWithContext(reqCtx, http.MethodGet, base+wellKnownOAuthPath, nil)
	if err != nil {
		return nil
	}
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return nil
	}

	body, err := io.ReadAll(io.LimitReader(resp.Body, oauthDiscoveryMaxBody))
	if err != nil {
		return nil
	}

	var doc oauthServerMetadata
	// The self-hosted SPA catch-all answers unknown paths with an HTML 200.
	if err := json.Unmarshal(body, &doc); err != nil {
		return nil
	}
	if doc.TokenEndpoint == "" || doc.DeviceAuthorizationEndpoint == "" {
		return nil
	}
	if err := validateOAuthMetadata(&doc, base); err != nil {
		return nil
	}
	return &doc
}

// validateOAuthMetadata enforces that the document describes the deployment we
// probed: RFC 8414 requires the issuer to match the URL the well-known path was
// built from, and every endpoint must share the issuer's origin. Refresh tokens
// are posted to these URLs, so an unvalidated document could redirect
// credentials to another host.
func validateOAuthMetadata(doc *oauthServerMetadata, base string) error {
	if strings.TrimRight(doc.Issuer, "/") != strings.TrimRight(base, "/") {
		return fmt.Errorf("issuer %q does not match %q", doc.Issuer, base)
	}
	issuer, err := url.Parse(doc.Issuer)
	if err != nil {
		return err
	}
	for _, ep := range []string{doc.DeviceAuthorizationEndpoint, doc.TokenEndpoint} {
		u, err := url.Parse(ep)
		if err != nil {
			return err
		}
		if u.Scheme != issuer.Scheme || u.Host != issuer.Host {
			return fmt.Errorf("endpoint %q is not on issuer origin", ep)
		}
	}
	return nil
}
