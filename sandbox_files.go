package langsmith

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"mime/multipart"
	"net/http"
	"net/url"
	"slices"
	"strconv"
	"strings"

	"github.com/langchain-ai/langsmith-go/internal/apierror"
	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// ReadFile reads a file from a named sandbox.
func (r *SandboxBoxService) ReadFile(ctx context.Context, name string, path string, opts ...option.RequestOption) ([]byte, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.DataplaneURL)
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
	dataplaneURL, err := requireSandboxDataplaneURL(box.Name, box.DataplaneURL)
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
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.ReadFileWithDataplaneURL(ctx, dataplaneURL, path, opts...)
}

// WriteFile writes bytes to a file in this sandbox.
func (s *Sandbox) WriteFile(ctx context.Context, path string, content []byte, opts ...option.RequestOption) error {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.DataplaneURL)
	if err != nil {
		return err
	}
	return s.boxes.WriteFileWithDataplaneURL(ctx, dataplaneURL, path, content, opts...)
}

// SandboxFileInfo is one filesystem entry returned by a glob.
type SandboxFileInfo struct {
	Path       string `json:"path"`
	IsDir      bool   `json:"is_dir"`
	SizeBytes  int64  `json:"size_bytes"`
	ModifiedAt string `json:"modified_at"`
}

// SandboxGlobParams selects files and directories under a root.
type SandboxGlobParams struct {
	// Pattern matches each entry's path relative to Path. Supports ** for any
	// number of segments plus *, ? and [...] within one segment.
	Pattern param.Field[string] `json:"pattern" api:"required"`
	// Path is the absolute directory to search under.
	Path  param.Field[string] `json:"path" api:"required"`
	Limit param.Field[int64]  `json:"limit"`
}

func (r SandboxGlobParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// SandboxGlobResult holds the entries matching a glob. Truncated reports that
// the server hit its result cap or deadline, so the answer is partial.
type SandboxGlobResult struct {
	Matches   []SandboxFileInfo `json:"matches"`
	Truncated bool              `json:"truncated"`
}

// SandboxGrepParams searches file contents for a literal string.
type SandboxGrepParams struct {
	// Pattern is literal text, not a regular expression.
	Pattern param.Field[string] `json:"pattern" api:"required"`
	// Path is the absolute directory to search under.
	Path param.Field[string] `json:"path" api:"required"`
	// Glob restricts which files are searched. A bare *.py matches by basename
	// at any depth; a pattern containing / or ** matches the path relative to
	// Path.
	Glob  param.Field[string] `json:"glob"`
	Limit param.Field[int64]  `json:"limit"`
}

func (r SandboxGrepParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// SandboxGrepMatch is one matching line. Line is 1-based.
type SandboxGrepMatch struct {
	Path string `json:"path"`
	Line int64  `json:"line"`
	Text string `json:"text"`
}

// SandboxGrepResult holds the lines matching a search. Truncated reports that
// the server hit its result cap or deadline.
type SandboxGrepResult struct {
	Matches   []SandboxGrepMatch `json:"matches"`
	Truncated bool               `json:"truncated"`
}

// Glob finds files and directories matching a pattern in a named sandbox.
func (r *SandboxBoxService) Glob(ctx context.Context, name string, body SandboxGlobParams, opts ...option.RequestOption) (*SandboxGlobResult, error) {
	dataplaneURL, err := r.resolveDataplaneURL(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	return r.GlobWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// GlobWithDataplaneURL runs a glob directly against a sandbox dataplane URL.
func (r *SandboxBoxService) GlobWithDataplaneURL(ctx context.Context, dataplaneURL string, body SandboxGlobParams, opts ...option.RequestOption) (res *SandboxGlobResult, err error) {
	opts = slices.Concat(r.Options, opts)
	path, err := sandboxDataplaneURL(dataplaneURL, "glob")
	if err != nil {
		return nil, err
	}
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Ls lists a directory's immediate entries in a named sandbox, without
// recursing: the non-recursive glob case.
func (r *SandboxBoxService) Ls(ctx context.Context, name string, path string, opts ...option.RequestOption) (*SandboxGlobResult, error) {
	return r.Glob(ctx, name, SandboxGlobParams{Pattern: F("*"), Path: F(path)}, opts...)
}

// Grep searches file contents in a named sandbox for a literal string.
func (r *SandboxBoxService) Grep(ctx context.Context, name string, body SandboxGrepParams, opts ...option.RequestOption) (*SandboxGrepResult, error) {
	dataplaneURL, err := r.resolveDataplaneURL(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	return r.GrepWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// GrepWithDataplaneURL runs a grep directly against a sandbox dataplane URL.
func (r *SandboxBoxService) GrepWithDataplaneURL(ctx context.Context, dataplaneURL string, body SandboxGrepParams, opts ...option.RequestOption) (res *SandboxGrepResult, err error) {
	opts = slices.Concat(r.Options, opts)
	path, err := sandboxDataplaneURL(dataplaneURL, "grep")
	if err != nil {
		return nil, err
	}
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// SandboxFileStat is what a HEAD on a sandbox file reports, without
// transferring it. ETag is opaque: compare it, never parse it.
type SandboxFileStat struct {
	SizeBytes    int64
	ETag         string
	LastModified string
	ContentType  string
}

// SandboxReadRangeParams selects the bytes a ranged read returns. Provide
// Start (with optional End), or SuffixBytes.
type SandboxReadRangeParams struct {
	Start param.Field[int64]
	// End is inclusive, and clamped to EOF rather than rejected.
	End param.Field[int64]
	// SuffixBytes returns the last N bytes. Cannot be combined with Start/End.
	SuffixBytes param.Field[int64]
	// IfRange is an ETag from an earlier chunk. The range is honored only
	// while the file still matches it.
	IfRange string
	// IfNoneMatch is an ETag the caller already holds. An unchanged file
	// answers with Unchanged set and no bytes.
	IfNoneMatch string
}

// SandboxFileChunk is the bytes a ranged read returned, and where they sit in
// the file.
type SandboxFileChunk struct {
	Content    []byte
	ETag       string
	TotalBytes int64
	Start      int64
	// Partial reports a 206. A false here after a ranged request means the
	// file changed and the server sent it whole -- restart from zero rather
	// than appending.
	Partial bool
	// Unchanged reports that the caller's IfNoneMatch still matches, so no
	// bytes were returned.
	Unchanged    bool
	LastModified string
}

// End is the offset just past the last byte returned.
func (c SandboxFileChunk) End() int64 {
	return c.Start + int64(len(c.Content))
}

// StatFile reports a file's size and validators in a named sandbox without
// transferring it.
func (r *SandboxBoxService) StatFile(ctx context.Context, name string, path string, opts ...option.RequestOption) (*SandboxFileStat, error) {
	dataplaneURL, err := r.resolveDataplaneURL(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	return r.StatFileWithDataplaneURL(ctx, dataplaneURL, path, opts...)
}

// StatFileWithDataplaneURL stats a file directly against a sandbox dataplane URL.
func (r *SandboxBoxService) StatFileWithDataplaneURL(ctx context.Context, dataplaneURL string, path string, opts ...option.RequestOption) (*SandboxFileStat, error) {
	opts = slices.Concat(r.Options, opts)
	u, err := sandboxDownloadURL(dataplaneURL, path)
	if err != nil {
		return nil, err
	}

	var res *http.Response
	opts = append(opts, option.WithResponseInto(&res))
	if err := requestconfig.ExecuteNewRequest(ctx, http.MethodHead, u, nil, nil, opts...); err != nil {
		return nil, sandboxFileRequestError(err, path)
	}
	return &SandboxFileStat{
		SizeBytes:    res.ContentLength,
		ETag:         res.Header.Get("ETag"),
		LastModified: res.Header.Get("Last-Modified"),
		ContentType:  res.Header.Get("Content-Type"),
	}, nil
}

// ReadFileRange reads part of a file in a named sandbox, for chunked reads and
// resumed downloads.
func (r *SandboxBoxService) ReadFileRange(ctx context.Context, name string, path string, body SandboxReadRangeParams, opts ...option.RequestOption) (*SandboxFileChunk, error) {
	dataplaneURL, err := r.resolveDataplaneURL(ctx, name, opts...)
	if err != nil {
		return nil, err
	}
	return r.ReadFileRangeWithDataplaneURL(ctx, dataplaneURL, path, body, opts...)
}

// ReadFileRangeWithDataplaneURL reads a byte range directly from a sandbox
// dataplane URL.
func (r *SandboxBoxService) ReadFileRangeWithDataplaneURL(ctx context.Context, dataplaneURL string, path string, body SandboxReadRangeParams, opts ...option.RequestOption) (*SandboxFileChunk, error) {
	opts = slices.Concat(r.Options, opts)
	rangeHeader, err := sandboxRangeHeader(body)
	if err != nil {
		return nil, err
	}
	u, err := sandboxDownloadURL(dataplaneURL, path)
	if err != nil {
		return nil, err
	}

	opts = append(opts, option.WithHeader("Range", rangeHeader))
	if body.IfRange != "" {
		opts = append(opts, option.WithHeader("If-Range", body.IfRange))
	}
	if body.IfNoneMatch != "" {
		opts = append(opts, option.WithHeader("If-None-Match", body.IfNoneMatch))
	}

	var res *http.Response
	var content []byte
	opts = append(opts, option.WithResponseInto(&res))
	if err := requestconfig.ExecuteNewRequest(ctx, http.MethodGet, u, nil, &content, opts...); err != nil {
		return nil, sandboxFileRequestError(err, path)
	}
	return sandboxFileChunkFromResponse(res, content), nil
}

func sandboxDownloadURL(dataplaneURL string, path string) (string, error) {
	requestURL, err := sandboxDataplaneURL(dataplaneURL, "download")
	if err != nil {
		return "", err
	}
	u, err := url.Parse(requestURL)
	if err != nil {
		return "", err
	}
	q := u.Query()
	q.Set("path", path)
	u.RawQuery = q.Encode()
	return u.String(), nil
}

func (r *SandboxBoxService) resolveDataplaneURL(ctx context.Context, name string, opts ...option.RequestOption) (string, error) {
	box, err := r.Get(ctx, name, opts...)
	if err != nil {
		return "", err
	}
	return requireSandboxDataplaneURL(box.Name, box.DataplaneURL)
}

// sandboxRangeHeader renders a byte range as an RFC 9110 Range header value.
func sandboxRangeHeader(body SandboxReadRangeParams) (string, error) {
	if body.SuffixBytes.Present {
		if body.Start.Present || body.End.Present {
			return "", errors.New("langsmith: cannot combine SuffixBytes with Start/End")
		}
		if body.SuffixBytes.Value <= 0 {
			return "", errors.New("langsmith: SuffixBytes must be positive")
		}
		return fmt.Sprintf("bytes=-%d", body.SuffixBytes.Value), nil
	}
	if !body.Start.Present {
		return "", errors.New("langsmith: provide Start (with optional End), or SuffixBytes")
	}
	if body.Start.Value < 0 {
		return "", errors.New("langsmith: Start must not be negative")
	}
	if !body.End.Present {
		return fmt.Sprintf("bytes=%d-", body.Start.Value), nil
	}
	if body.End.Value < body.Start.Value {
		return "", errors.New("langsmith: End must not precede Start")
	}
	return fmt.Sprintf("bytes=%d-%d", body.Start.Value, body.End.Value), nil
}

func sandboxFileChunkFromResponse(res *http.Response, content []byte) *SandboxFileChunk {
	chunk := &SandboxFileChunk{
		ETag:         res.Header.Get("ETag"),
		LastModified: res.Header.Get("Last-Modified"),
	}
	if res.StatusCode == http.StatusNotModified {
		chunk.Unchanged = true
		return chunk
	}
	chunk.Content = content
	if res.StatusCode == http.StatusPartialContent {
		chunk.Partial = true
		chunk.Start, chunk.TotalBytes = parseSandboxContentRange(res.Header.Get("Content-Range"))
		return chunk
	}
	// A stale If-Range answers 200 with the whole file; the caller has to
	// restart rather than append, so report it from byte zero.
	chunk.TotalBytes = int64(len(content))
	return chunk
}

// parseSandboxContentRange reads the first byte offset and total size out of a
// Content-Range header.
func parseSandboxContentRange(value string) (start int64, total int64) {
	spec, ok := strings.CutPrefix(value, "bytes ")
	if !ok {
		return 0, 0
	}
	rangePart, totalPart, _ := strings.Cut(strings.TrimSpace(spec), "/")
	first, _, _ := strings.Cut(rangePart, "-")
	if parsed, err := strconv.ParseInt(strings.TrimSpace(first), 10, 64); err == nil {
		start = parsed
	}
	if parsed, err := strconv.ParseInt(strings.TrimSpace(totalPart), 10, 64); err == nil {
		total = parsed
	}
	return start, total
}

// sandboxFileRequestError maps a 416 onto an operation error: its body is
// text/plain, not the endpoint's usual JSON error shape.
func sandboxFileRequestError(err error, path string) error {
	var aerr *apierror.Error
	if errors.As(err, &aerr) && aerr.StatusCode == http.StatusRequestedRangeNotSatisfiable {
		return &SandboxOperationError{
			Operation: "read",
			Message:   fmt.Sprintf("requested range for %q starts past the end of the file", path),
		}
	}
	return err
}

// Glob finds files and directories matching a pattern in this sandbox.
func (s *Sandbox) Glob(ctx context.Context, body SandboxGlobParams, opts ...option.RequestOption) (*SandboxGlobResult, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.GlobWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// Ls lists a directory's immediate entries in this sandbox, without recursing.
func (s *Sandbox) Ls(ctx context.Context, path string, opts ...option.RequestOption) (*SandboxGlobResult, error) {
	return s.Glob(ctx, SandboxGlobParams{Pattern: F("*"), Path: F(path)}, opts...)
}

// Grep searches file contents in this sandbox for a literal string.
func (s *Sandbox) Grep(ctx context.Context, body SandboxGrepParams, opts ...option.RequestOption) (*SandboxGrepResult, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.GrepWithDataplaneURL(ctx, dataplaneURL, body, opts...)
}

// StatFile reports a file's size and validators in this sandbox without
// transferring it.
func (s *Sandbox) StatFile(ctx context.Context, path string, opts ...option.RequestOption) (*SandboxFileStat, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.StatFileWithDataplaneURL(ctx, dataplaneURL, path, opts...)
}

// ReadFileRange reads part of a file in this sandbox.
func (s *Sandbox) ReadFileRange(ctx context.Context, path string, body SandboxReadRangeParams, opts ...option.RequestOption) (*SandboxFileChunk, error) {
	dataplaneURL, err := requireSandboxDataplaneURL(s.Name, s.DataplaneURL)
	if err != nil {
		return nil, err
	}
	return s.boxes.ReadFileRangeWithDataplaneURL(ctx, dataplaneURL, path, body, opts...)
}
