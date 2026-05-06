package langsmith

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
)

const (
	sandboxServiceAuthHeader           = "X-Langsmith-Sandbox-Service-Token"
	defaultSandboxServiceURLTTLSeconds = int64(600)
	maxSandboxServiceURLTTLSeconds     = int64(86400)
)

// Service returns an auto-refreshing authenticated service URL helper.
func (r *SandboxBoxService) Service(ctx context.Context, name string, body SandboxBoxGenerateServiceURLParams, opts ...option.RequestOption) (*SandboxServiceURL, error) {
	body, err := normalizeSandboxServiceURLParams(body)
	if err != nil {
		return nil, err
	}
	res, err := r.GenerateServiceURL(ctx, name, body, opts...)
	if err != nil {
		return nil, err
	}
	return newSandboxServiceURL(res, func(ctx context.Context) (*SandboxServiceURL, error) {
		return r.Service(ctx, name, body, opts...)
	}), nil
}

// SandboxServiceURL is an authenticated URL for an HTTP service inside a sandbox.
type SandboxServiceURL struct {
	mu        sync.Mutex
	browser   string
	service   string
	token     string
	expiresAt string
	refresher func(context.Context) (*SandboxServiceURL, error)
}

func newSandboxServiceURL(res *SandboxBoxGenerateServiceURLResponse, refresher func(context.Context) (*SandboxServiceURL, error)) *SandboxServiceURL {
	return &SandboxServiceURL{
		browser:   res.BrowserURL,
		service:   res.ServiceURL,
		token:     res.Token,
		expiresAt: res.ExpiresAt,
		refresher: refresher,
	}
}

// BrowserURL returns a URL suitable for browser use, refreshing the token if it
// is near expiry.
func (s *SandboxServiceURL) BrowserURL(ctx context.Context) (string, error) {
	if err := s.refreshIfNeeded(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.browser, nil
}

// ServiceURL returns the base service URL, refreshing the token if it is near
// expiry.
func (s *SandboxServiceURL) ServiceURL(ctx context.Context) (string, error) {
	if err := s.refreshIfNeeded(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.service, nil
}

// Token returns the raw service token, refreshing it if it is near expiry.
func (s *SandboxServiceURL) Token(ctx context.Context) (string, error) {
	if err := s.refreshIfNeeded(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.token, nil
}

// ExpiresAt returns the token expiration timestamp, refreshing it if it is near
// expiry.
func (s *SandboxServiceURL) ExpiresAt(ctx context.Context) (string, error) {
	if err := s.refreshIfNeeded(ctx); err != nil {
		return "", err
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	return s.expiresAt, nil
}

// Request sends an HTTP request to the sandbox service and injects the service
// auth token.
func (s *SandboxServiceURL) Request(ctx context.Context, method string, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	if err := s.refreshIfNeeded(ctx); err != nil {
		return nil, err
	}
	s.mu.Lock()
	base := s.service
	token := s.token
	s.mu.Unlock()

	requestURL := strings.TrimRight(base, "/") + "/" + strings.TrimLeft(path, "/")
	req, err := http.NewRequestWithContext(ctx, method, requestURL, body)
	if err != nil {
		return nil, err
	}
	for key, values := range headers {
		for _, value := range values {
			req.Header.Add(key, value)
		}
	}
	req.Header.Set(sandboxServiceAuthHeader, token)
	return http.DefaultClient.Do(req)
}

// Get sends a GET request to the sandbox service.
func (s *SandboxServiceURL) Get(ctx context.Context, path string, headers http.Header) (*http.Response, error) {
	return s.Request(ctx, http.MethodGet, path, nil, headers)
}

// Post sends a POST request to the sandbox service.
func (s *SandboxServiceURL) Post(ctx context.Context, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	return s.Request(ctx, http.MethodPost, path, body, headers)
}

// Put sends a PUT request to the sandbox service.
func (s *SandboxServiceURL) Put(ctx context.Context, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	return s.Request(ctx, http.MethodPut, path, body, headers)
}

// Patch sends a PATCH request to the sandbox service.
func (s *SandboxServiceURL) Patch(ctx context.Context, path string, body io.Reader, headers http.Header) (*http.Response, error) {
	return s.Request(ctx, http.MethodPatch, path, body, headers)
}

// Delete sends a DELETE request to the sandbox service.
func (s *SandboxServiceURL) Delete(ctx context.Context, path string, headers http.Header) (*http.Response, error) {
	return s.Request(ctx, http.MethodDelete, path, nil, headers)
}

func (s *SandboxServiceURL) refreshIfNeeded(ctx context.Context) error {
	s.mu.Lock()
	if !s.shouldRefreshLocked() || s.refresher == nil {
		s.mu.Unlock()
		return nil
	}
	refresher := s.refresher
	s.mu.Unlock()

	fresh, err := refresher(ctx)
	if err != nil {
		return err
	}

	s.mu.Lock()
	s.browser = fresh.browser
	s.service = fresh.service
	s.token = fresh.token
	s.expiresAt = fresh.expiresAt
	s.refresher = fresh.refresher
	s.mu.Unlock()
	return nil
}

func (s *SandboxServiceURL) shouldRefreshLocked() bool {
	if s.expiresAt == "" {
		return false
	}
	expires, err := time.Parse(time.RFC3339, strings.Replace(s.expiresAt, "Z", "+00:00", 1))
	if err != nil {
		return false
	}
	return time.Until(expires) <= 30*time.Second
}

// Service returns an auto-refreshing service URL helper for this sandbox.
func (s *Sandbox) Service(ctx context.Context, body SandboxBoxGenerateServiceURLParams, opts ...option.RequestOption) (*SandboxServiceURL, error) {
	return s.boxes.Service(ctx, s.Name, body, opts...)
}

func normalizeSandboxServiceURLParams(body SandboxBoxGenerateServiceURLParams) (SandboxBoxGenerateServiceURLParams, error) {
	port := sandboxFieldValue(body.Port, int64(0))
	if port < 1 || port > 65535 {
		return SandboxBoxGenerateServiceURLParams{}, fmt.Errorf("port must be between 1 and 65535, got %d", port)
	}
	expiresInSeconds := sandboxFieldValue(body.ExpiresInSeconds, defaultSandboxServiceURLTTLSeconds)
	if expiresInSeconds < 1 || expiresInSeconds > maxSandboxServiceURLTTLSeconds {
		return SandboxBoxGenerateServiceURLParams{}, fmt.Errorf("expires_in_seconds must be between 1 and %d, got %d", maxSandboxServiceURLTTLSeconds, expiresInSeconds)
	}
	body.Port = F(port)
	body.ExpiresInSeconds = F(expiresInSeconds)
	return body, nil
}
