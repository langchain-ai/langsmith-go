package langsmith

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

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
