package models

import "net/http"

// WriteEndpoint identifies a LangSmith API endpoint to send traces to.
type WriteEndpoint struct {
	URL              string
	Key              string // API key (sent as X-API-Key header).
	OAuthAccessToken string // OAuth access token (sent as Authorization header); takes precedence over Key.
	Project          string
}

// SetAuthHeader sets the appropriate authentication header on req.
// OAuth access token takes precedence over API key.
func (ep WriteEndpoint) SetAuthHeader(req *http.Request) {
	if ep.OAuthAccessToken != "" {
		req.Header.Set("Authorization", "Bearer "+ep.OAuthAccessToken)
	} else if ep.Key != "" {
		req.Header.Set("X-API-Key", ep.Key)
	}
}
