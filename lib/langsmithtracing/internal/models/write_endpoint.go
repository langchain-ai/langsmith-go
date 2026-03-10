package models

import "net/http"

// WriteEndpoint identifies a LangSmith API endpoint to send traces to.
type WriteEndpoint struct {
	URL          string
	Key          string // API key (sent as X-API-Key header).
	BearerToken  string // Bearer token (sent as Authorization header); takes precedence over Key.
	Project      string
}

// SetAuthHeader sets the appropriate authentication header on req.
// Bearer token takes precedence over API key.
func (ep WriteEndpoint) SetAuthHeader(req *http.Request) {
	if ep.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+ep.BearerToken)
	} else if ep.Key != "" {
		req.Header.Set("X-API-Key", ep.Key)
	}
}
