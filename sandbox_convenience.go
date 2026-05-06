package langsmith

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strings"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

const (
	sandboxServiceAuthHeader           = "X-Langsmith-Sandbox-Service-Token"
	defaultSandboxServiceURLTTLSeconds = int64(600)
	maxSandboxServiceURLTTLSeconds     = int64(86400)
)

// SandboxResourceTimeoutError is returned when waiting for a sandbox resource
// exceeds the configured timeout.
type SandboxResourceTimeoutError struct {
	ResourceType string
	ResourceID   string
	LastStatus   string
	Timeout      time.Duration
}

func (e *SandboxResourceTimeoutError) Error() string {
	if e.LastStatus != "" {
		return fmt.Sprintf("langsmith: %s %q not ready after %s (last_status: %s)", e.ResourceType, e.ResourceID, e.Timeout, e.LastStatus)
	}
	return fmt.Sprintf("langsmith: %s %q not ready after %s", e.ResourceType, e.ResourceID, e.Timeout)
}

// SandboxResourceCreationError is returned when a sandbox resource reaches a
// failed provisioning state.
type SandboxResourceCreationError struct {
	ResourceType string
	ResourceID   string
	Message      string
}

func (e *SandboxResourceCreationError) Error() string {
	if e.Message != "" {
		return fmt.Sprintf("langsmith: %s %q failed: %s", e.ResourceType, e.ResourceID, e.Message)
	}
	return fmt.Sprintf("langsmith: %s %q failed", e.ResourceType, e.ResourceID)
}

// SandboxWaitParams configures polling for sandbox readiness.
type SandboxWaitParams struct {
	Timeout      time.Duration
	PollInterval time.Duration
}

// SnapshotWaitParams configures polling for snapshot readiness.
type SnapshotWaitParams struct {
	Timeout      time.Duration
	PollInterval time.Duration
}

// Sandbox is a convenience wrapper around generated sandbox box responses.
type Sandbox struct {
	ID              string
	Name            string
	DataplaneURL    string
	Status          string
	StatusMessage   string
	CreatedAt       string
	UpdatedAt       string
	TTLSeconds      int64
	IdleTTLSeconds  int64
	ExpiresAt       string
	SnapshotID      string
	Vcpus           int64
	MemBytes        int64
	FsCapacityBytes int64

	boxes *SandboxBoxService
}

// NewSandbox creates a sandbox and returns the convenience wrapper.
func (r *SandboxBoxService) NewSandbox(ctx context.Context, body SandboxBoxNewParams, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromNewResponse(res, r), nil
}

// GetSandbox retrieves a sandbox and returns the convenience wrapper.
func (r *SandboxBoxService) GetSandbox(ctx context.Context, name string, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromGetResponse(res, r), nil
}

// ListSandboxes lists sandboxes and returns convenience wrappers.
func (r *SandboxBoxService) ListSandboxes(ctx context.Context, query SandboxBoxListParams, opts ...option.RequestOption) ([]*Sandbox, error) {
	res, err := r.List(ctx, query, opts...)
	if err != nil {
		return nil, err
	}
	out := make([]*Sandbox, 0, len(res.Sandboxes))
	for i := range res.Sandboxes {
		out = append(out, sandboxFromListResponse(&res.Sandboxes[i], r))
	}
	return out, nil
}

// Wait polls the generated status endpoint until the sandbox is ready or failed.
func (r *SandboxBoxService) Wait(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*SandboxBoxGetResponse, error) {
	timeout := params.Timeout
	if timeout == 0 {
		timeout = 120 * time.Second
	}
	pollInterval := params.PollInterval
	if pollInterval == 0 {
		pollInterval = time.Second
	}
	deadline := time.Now().Add(timeout)
	lastStatus := ""

	for {
		status, err := r.GetStatus(ctx, name, opts...)
		if err != nil {
			return nil, err
		}
		lastStatus = status.Status
		switch status.Status {
		case "ready":
			return r.Get(ctx, name, opts...)
		case "failed":
			return nil, &SandboxResourceCreationError{
				ResourceType: "sandbox",
				ResourceID:   name,
				Message:      status.StatusMessage,
			}
		}

		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, &SandboxResourceTimeoutError{
				ResourceType: "sandbox",
				ResourceID:   name,
				LastStatus:   lastStatus,
				Timeout:      timeout,
			}
		}
		delay := minDuration(pollInterval, remaining)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
}

// WaitSandbox polls until ready and returns the convenience wrapper.
func (r *SandboxBoxService) WaitSandbox(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.Wait(ctx, name, params, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromGetResponse(res, r), nil
}

// StartAndWait starts a stopped sandbox and waits until it is ready.
func (r *SandboxBoxService) StartAndWait(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*SandboxBoxGetResponse, error) {
	if _, err := r.Start(ctx, name, opts...); err != nil {
		return nil, err
	}
	return r.Wait(ctx, name, params, opts...)
}

// StartSandbox starts a stopped sandbox and returns the convenience wrapper once
// it is ready.
func (r *SandboxBoxService) StartSandbox(ctx context.Context, name string, params SandboxWaitParams, opts ...option.RequestOption) (*Sandbox, error) {
	res, err := r.StartAndWait(ctx, name, params, opts...)
	if err != nil {
		return nil, err
	}
	return sandboxFromGetResponse(res, r), nil
}

// Wait polls until a snapshot reaches ready or failed status.
func (r *SandboxSnapshotService) Wait(ctx context.Context, snapshotID string, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	timeout := params.Timeout
	if timeout == 0 {
		timeout = 300 * time.Second
	}
	pollInterval := params.PollInterval
	if pollInterval == 0 {
		pollInterval = 2 * time.Second
	}
	deadline := time.Now().Add(timeout)
	lastStatus := ""

	for {
		snapshot, err := r.Get(ctx, snapshotID, opts...)
		if err != nil {
			return nil, err
		}
		lastStatus = snapshot.Status
		switch snapshot.Status {
		case "ready":
			return snapshot, nil
		case "failed":
			return nil, &SandboxResourceCreationError{
				ResourceType: "snapshot",
				ResourceID:   snapshotID,
				Message:      snapshot.StatusMessage,
			}
		}

		remaining := time.Until(deadline)
		if remaining <= 0 {
			return nil, &SandboxResourceTimeoutError{
				ResourceType: "snapshot",
				ResourceID:   snapshotID,
				LastStatus:   lastStatus,
				Timeout:      timeout,
			}
		}
		delay := minDuration(pollInterval, remaining)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(delay):
		}
	}
}

// NewAndWait creates a snapshot and waits until it is ready or failed.
func (r *SandboxSnapshotService) NewAndWait(ctx context.Context, body SandboxSnapshotNewParams, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	snapshot, err := r.New(ctx, body, opts...)
	if err != nil {
		return nil, err
	}
	return r.Wait(ctx, snapshot.ID, params, opts...)
}

// CaptureSnapshotAndWait captures a snapshot from a sandbox and waits until it
// is ready or failed.
func (r *SandboxBoxService) CaptureSnapshotAndWait(ctx context.Context, name string, body SandboxBoxNewSnapshotParams, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	snapshot, err := r.NewSnapshot(ctx, name, body, opts...)
	if err != nil {
		return nil, err
	}
	return NewSandboxSnapshotService(r.Options...).Wait(ctx, snapshot.ID, params, opts...)
}

// ReadFile reads a file from a named sandbox.
func (r *SandboxBoxService) ReadFile(ctx context.Context, name string, path string, opts ...option.RequestOption) ([]byte, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return r.ReadFileWithDataplaneURL(ctx, dataplaneURL, path, opts...)
}

// ReadFileWithDataplaneURL reads a file directly from a sandbox dataplane URL.
func (r *SandboxBoxService) ReadFileWithDataplaneURL(ctx context.Context, dataplaneURL string, path string, opts ...option.RequestOption) ([]byte, error) {
	opts = slices.Concat(r.Options, opts)
	requestURL, err := sandboxDataplaneURL(dataplaneURL, "download")
	if err != nil {
		return nil, err
	}
	u, err := url.Parse(requestURL)
	if err != nil {
		return nil, err
	}
	q := u.Query()
	q.Set("path", path)
	u.RawQuery = q.Encode()

	var out []byte
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, u.String(), nil, &out, opts...)
	return out, err
}

// WriteFile writes bytes to a file in a named sandbox.
func (r *SandboxBoxService) WriteFile(ctx context.Context, name string, path string, content []byte, opts ...option.RequestOption) error {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.Status, box.DataplaneURL)
	if err != nil {
		return err
	}
	return r.WriteFileWithDataplaneURL(ctx, dataplaneURL, path, content, opts...)
}

// WriteFileWithDataplaneURL writes bytes directly to a sandbox dataplane URL.
func (r *SandboxBoxService) WriteFileWithDataplaneURL(ctx context.Context, dataplaneURL string, path string, content []byte, opts ...option.RequestOption) error {
	opts = slices.Concat(r.Options, opts)
	requestURL, err := sandboxDataplaneURL(dataplaneURL, "upload")
	if err != nil {
		return err
	}
	u, err := url.Parse(requestURL)
	if err != nil {
		return err
	}
	q := u.Query()
	q.Set("path", path)
	u.RawQuery = q.Encode()

	var buf bytes.Buffer
	writer := multipart.NewWriter(&buf)
	part, err := writer.CreateFormFile("file", "file")
	if err != nil {
		return err
	}
	if _, err := part.Write(content); err != nil {
		return err
	}
	if err := writer.Close(); err != nil {
		return err
	}

	opts = slices.Concat([]option.RequestOption{
		option.WithRequestBody(writer.FormDataContentType(), &buf),
	}, opts)
	return requestconfig.ExecuteNewRequest(ctx, http.MethodPost, u.String(), nil, nil, opts...)
}

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

// Refresh fetches latest sandbox state and updates this object.
func (s *Sandbox) Refresh(ctx context.Context, opts ...option.RequestOption) error {
	box, err := s.boxes.Get(ctx, s.Name, opts...)
	if err != nil {
		return err
	}
	s.applyGetResponse(box)
	return nil
}

// Run executes a command and waits for completion.
func (s *Sandbox) Run(ctx context.Context, body SandboxBoxRunParams, opts ...option.RequestOption) (*SandboxExecutionResult, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.RunWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// StartCommand starts a streaming command in this sandbox.
func (s *Sandbox) StartCommand(ctx context.Context, body SandboxCommandStartParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.StartCommandWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// RunWithCallbacks starts a WebSocket command, invokes callbacks for output,
// and waits for completion.
func (s *Sandbox) RunWithCallbacks(ctx context.Context, body SandboxCommandStartParams, callbacks SandboxCommandCallbacks, opts ...option.RequestOption) (*SandboxExecutionResult, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.RunWithDataplaneURLAndCallbacks(ctx, dataplaneURL, body, callbacks, opts...)
}

// ReconnectCommand reconnects to a command stream.
func (s *Sandbox) ReconnectCommand(ctx context.Context, commandID string, body SandboxCommandReconnectParams, opts ...option.RequestOption) (*SandboxCommandHandle, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.ReconnectCommandWithDataplaneURL(ctx, dataplaneURL, commandID, body, opts...)
}

// ReadFile reads a file from this sandbox.
func (s *Sandbox) ReadFile(ctx context.Context, path string, opts ...option.RequestOption) ([]byte, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.ReadFileWithDataplaneURL(ctx, dataplaneURL, path, opts...)
}

// WriteFile writes bytes to a file in this sandbox.
func (s *Sandbox) WriteFile(ctx context.Context, path string, content []byte, opts ...option.RequestOption) error {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.Status, s.DataplaneURL)
	if err != nil {
		return err
	}
	return s.boxes.WriteFileWithDataplaneURL(ctx, dataplaneURL, path, content, opts...)
}

// Service returns an auto-refreshing service URL helper for this sandbox.
func (s *Sandbox) Service(ctx context.Context, body SandboxBoxGenerateServiceURLParams, opts ...option.RequestOption) (*SandboxServiceURL, error) {
	return s.boxes.Service(ctx, s.Name, body, opts...)
}

// Update updates this sandbox and refreshes this object's fields.
func (s *Sandbox) Update(ctx context.Context, body SandboxBoxUpdateParams, opts ...option.RequestOption) error {
	box, err := s.boxes.Update(ctx, s.Name, body, opts...)
	if err != nil {
		return err
	}
	s.applyUpdateResponse(box)
	return nil
}

// Start starts this sandbox and waits until it is ready.
func (s *Sandbox) Start(ctx context.Context, params SandboxWaitParams, opts ...option.RequestOption) error {
	box, err := s.boxes.StartAndWait(ctx, s.Name, params, opts...)
	if err != nil {
		return err
	}
	s.applyGetResponse(box)
	return nil
}

// Stop stops this sandbox.
func (s *Sandbox) Stop(ctx context.Context, opts ...option.RequestOption) error {
	if err := s.boxes.Stop(ctx, s.Name, opts...); err != nil {
		return err
	}
	s.Status = "stopped"
	s.DataplaneURL = ""
	return nil
}

// Delete deletes this sandbox.
func (s *Sandbox) Delete(ctx context.Context, opts ...option.RequestOption) error {
	return s.boxes.Delete(ctx, s.Name, opts...)
}

// CaptureSnapshot captures a snapshot from this sandbox.
func (s *Sandbox) CaptureSnapshot(ctx context.Context, body SandboxBoxNewSnapshotParams, opts ...option.RequestOption) (*SandboxBoxNewSnapshotResponse, error) {
	return s.boxes.NewSnapshot(ctx, s.Name, body, opts...)
}

// CaptureSnapshotAndWait captures a snapshot and waits until it is ready.
func (s *Sandbox) CaptureSnapshotAndWait(ctx context.Context, body SandboxBoxNewSnapshotParams, params SnapshotWaitParams, opts ...option.RequestOption) (*SandboxSnapshotGetResponse, error) {
	return s.boxes.CaptureSnapshotAndWait(ctx, s.Name, body, params, opts...)
}

func sandboxFromNewResponse(res *SandboxBoxNewResponse, boxes *SandboxBoxService) *Sandbox {
	if res == nil {
		return nil
	}
	return &Sandbox{
		ID:              res.ID,
		Name:            res.Name,
		DataplaneURL:    res.DataplaneURL,
		Status:          res.Status,
		StatusMessage:   res.StatusMessage,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		TTLSeconds:      res.TtlSeconds,
		IdleTTLSeconds:  res.IdleTtlSeconds,
		ExpiresAt:       res.ExpiresAt,
		SnapshotID:      res.SnapshotID,
		Vcpus:           res.Vcpus,
		MemBytes:        res.MemBytes,
		FsCapacityBytes: res.FsCapacityBytes,
		boxes:           boxes,
	}
}

func sandboxFromGetResponse(res *SandboxBoxGetResponse, boxes *SandboxBoxService) *Sandbox {
	if res == nil {
		return nil
	}
	sb := &Sandbox{boxes: boxes}
	sb.applyGetResponse(res)
	return sb
}

func sandboxFromListResponse(res *SandboxBoxListResponseSandbox, boxes *SandboxBoxService) *Sandbox {
	if res == nil {
		return nil
	}
	return &Sandbox{
		ID:              res.ID,
		Name:            res.Name,
		DataplaneURL:    res.DataplaneURL,
		Status:          res.Status,
		StatusMessage:   res.StatusMessage,
		CreatedAt:       res.CreatedAt,
		UpdatedAt:       res.UpdatedAt,
		TTLSeconds:      res.TtlSeconds,
		IdleTTLSeconds:  res.IdleTtlSeconds,
		ExpiresAt:       res.ExpiresAt,
		SnapshotID:      res.SnapshotID,
		Vcpus:           res.Vcpus,
		MemBytes:        res.MemBytes,
		FsCapacityBytes: res.FsCapacityBytes,
		boxes:           boxes,
	}
}

func (s *Sandbox) applyGetResponse(res *SandboxBoxGetResponse) {
	s.ID = res.ID
	s.Name = res.Name
	s.DataplaneURL = res.DataplaneURL
	s.Status = res.Status
	s.StatusMessage = res.StatusMessage
	s.CreatedAt = res.CreatedAt
	s.UpdatedAt = res.UpdatedAt
	s.TTLSeconds = res.TtlSeconds
	s.IdleTTLSeconds = res.IdleTtlSeconds
	s.ExpiresAt = res.ExpiresAt
	s.SnapshotID = res.SnapshotID
	s.Vcpus = res.Vcpus
	s.MemBytes = res.MemBytes
	s.FsCapacityBytes = res.FsCapacityBytes
}

func (s *Sandbox) applyUpdateResponse(res *SandboxBoxUpdateResponse) {
	s.ID = res.ID
	s.Name = res.Name
	s.DataplaneURL = res.DataplaneURL
	s.Status = res.Status
	s.StatusMessage = res.StatusMessage
	s.CreatedAt = res.CreatedAt
	s.UpdatedAt = res.UpdatedAt
	s.TTLSeconds = res.TtlSeconds
	s.IdleTTLSeconds = res.IdleTtlSeconds
	s.ExpiresAt = res.ExpiresAt
	s.SnapshotID = res.SnapshotID
	s.Vcpus = res.Vcpus
	s.MemBytes = res.MemBytes
	s.FsCapacityBytes = res.FsCapacityBytes
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

func minDuration(a time.Duration, b time.Duration) time.Duration {
	if a < b {
		return a
	}
	return b
}
