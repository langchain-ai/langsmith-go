package auth

import (
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strings"
)

// SetUserIDHeaderFromAccessToken sets X-User-Id from the access token's JWT
// subject when present.
func SetUserIDHeaderFromAccessToken(header http.Header, accessToken string) {
	if userID := UserIDFromAccessToken(accessToken); userID != "" {
		header.Set("X-User-Id", userID)
	}
}

// UserIDFromAccessToken extracts the subject from a JWT access token.
func UserIDFromAccessToken(accessToken string) string {
	parts := strings.Split(accessToken, ".")
	if len(parts) < 2 {
		return ""
	}

	payload, err := base64.RawURLEncoding.DecodeString(parts[1])
	if err != nil {
		payload, err = base64.URLEncoding.DecodeString(parts[1])
		if err != nil {
			return ""
		}
	}

	var claims struct {
		Subject string `json:"sub"`
	}
	if err := json.Unmarshal(payload, &claims); err != nil {
		return ""
	}
	return claims.Subject
}
