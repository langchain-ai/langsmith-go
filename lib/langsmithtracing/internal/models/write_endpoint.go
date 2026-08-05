package models

import (
	"net/http"

	"github.com/langchain-ai/langsmith-go/internal/auth"
)

// WriteEndpoint identifies a LangSmith API endpoint to send traces to.
type WriteEndpoint struct {
	URL              string
	Key              string // API key (sent as X-API-Key header).
	OAuthAccessToken string // OAuth access token (sent as Authorization header); takes precedence over Key.
	Project          string

	// ExtraHeaders are sent on every request to this endpoint. They are applied
	// after the authentication headers, so an entry may replace one of them.
	ExtraHeaders map[string]string
}

// SetAuthHeader sets the appropriate authentication header on req.
// OAuth access token takes precedence over API key.
func (ep WriteEndpoint) SetAuthHeader(req *http.Request) {
	if ep.OAuthAccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+ep.OAuthAccessToken)
		auth.SetUserIDHeaderFromAccessToken(req.Header, ep.OAuthAccessToken)
	} else if ep.Key != "" {
		req.Header.Set("X-API-Key", ep.Key)
	}
}

// SetExtraHeaders applies ExtraHeaders to req, replacing any header req
// already carries under the same name.
func (ep WriteEndpoint) SetExtraHeaders(req *http.Request) {
	for name, value := range ep.ExtraHeaders {
		req.Header.Set(name, value)
	}
}
