package langsmith

import (
	"context"
	"testing"
	"time"
)

func TestNormalizeSandboxServiceURLParams(t *testing.T) {
	params, err := normalizeSandboxServiceURLParams(SandboxBoxGenerateServiceURLParams{
		Port: Int(8080),
	})
	if err != nil {
		t.Fatalf("normalizeSandboxServiceURLParams returned error: %v", err)
	}
	if params.Port.Value != 8080 {
		t.Fatalf("unexpected port: %d", params.Port.Value)
	}
	if params.ExpiresInSeconds.Value != defaultSandboxServiceURLTTLSeconds {
		t.Fatalf("expected default ttl, got %d", params.ExpiresInSeconds.Value)
	}

	if _, err := normalizeSandboxServiceURLParams(SandboxBoxGenerateServiceURLParams{Port: Int(0)}); err == nil {
		t.Fatal("expected invalid port error")
	}
	if _, err := normalizeSandboxServiceURLParams(SandboxBoxGenerateServiceURLParams{
		Port:             Int(8080),
		ExpiresInSeconds: Int(maxSandboxServiceURLTTLSeconds + 1),
	}); err == nil {
		t.Fatal("expected invalid ttl error")
	}
}

func TestSandboxServiceURLAccessorsWithoutRefresh(t *testing.T) {
	service := newSandboxServiceURL(&SandboxBoxGenerateServiceURLResponse{
		BrowserURL: "https://browser.example",
		ServiceURL: "https://service.example",
		Token:      "token-a",
		ExpiresAt:  time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
	}, nil)

	browserURL, err := service.BrowserURL(context.Background())
	if err != nil {
		t.Fatalf("BrowserURL returned error: %v", err)
	}
	if browserURL != "https://browser.example" {
		t.Fatalf("unexpected browser URL: %q", browserURL)
	}
	serviceURL, err := service.ServiceURL(context.Background())
	if err != nil {
		t.Fatalf("ServiceURL returned error: %v", err)
	}
	if serviceURL != "https://service.example" {
		t.Fatalf("unexpected service URL: %q", serviceURL)
	}
	token, err := service.Token(context.Background())
	if err != nil {
		t.Fatalf("Token returned error: %v", err)
	}
	if token != "token-a" {
		t.Fatalf("unexpected token: %q", token)
	}
	expiresAt, err := service.ExpiresAt(context.Background())
	if err != nil {
		t.Fatalf("ExpiresAt returned error: %v", err)
	}
	if expiresAt == "" {
		t.Fatal("expected expiration timestamp")
	}
}
