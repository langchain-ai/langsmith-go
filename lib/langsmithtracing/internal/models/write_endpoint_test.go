package models

import (
	"encoding/base64"
	"fmt"
	"net/http"
	"testing"
)

func jwtWithSubject(t *testing.T, sub string) string {
	t.Helper()
	header := base64.RawURLEncoding.EncodeToString([]byte(`{"alg":"none"}`))
	payload := base64.RawURLEncoding.EncodeToString([]byte(fmt.Sprintf(`{"sub":%q}`, sub)))
	return header + "." + payload + "."
}

func TestWriteEndpointSetAuthHeaderAddsUserIDFromOAuthJWTSubject(t *testing.T) {
	accessToken := jwtWithSubject(t, "user-123")
	req, err := http.NewRequest(http.MethodPost, "https://example.com/runs/multipart", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	WriteEndpoint{OAuthAccessToken: accessToken}.SetAuthHeader(req)

	if got, want := req.Header.Get("Authorization"), "Bearer "+accessToken; got != want {
		t.Fatalf("Authorization = %q, want %q", got, want)
	}
	if got, want := req.Header.Get("X-User-Id"), "user-123"; got != want {
		t.Fatalf("X-User-Id = %q, want %q", got, want)
	}
}

func TestWriteEndpointSetExtraHeaders(t *testing.T) {
	req, err := http.NewRequest(http.MethodPost, "https://example.com/runs/multipart", nil)
	if err != nil {
		t.Fatalf("NewRequest: %v", err)
	}

	ep := WriteEndpoint{
		Key: "api-key",
		ExtraHeaders: map[string]string{
			"X-Tenant-Id": "tenant-1",
			"X-API-Key":   "override",
		},
	}
	ep.SetAuthHeader(req)
	ep.SetExtraHeaders(req)

	if got, want := req.Header.Get("X-Tenant-Id"), "tenant-1"; got != want {
		t.Errorf("X-Tenant-Id = %q, want %q", got, want)
	}
	if got, want := req.Header.Get("X-API-Key"), "override"; got != want {
		t.Errorf("X-API-Key = %q, want %q", got, want)
	}
}
