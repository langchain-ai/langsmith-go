package models

// WriteEndpoint identifies a LangSmith API endpoint to send traces to.
type WriteEndpoint struct {
	URL     string
	Key     string
	Project string
}
