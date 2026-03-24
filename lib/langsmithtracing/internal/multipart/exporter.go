package multipart

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/textproto"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/google/uuid"
	"github.com/klauspost/compress/zstd"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/logger"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/models"
)

const (
	errorPreviewSize       = 4 << 10
	defaultBatchSizeLimit  = 20 << 20 // 20 MiB
	batchJSONOverheadBytes = 64       // conservative overhead for {"post":[],"patch":[]} framing + commas
)

// APIError is returned when the LangSmith API returns a non-2xx status code.
type APIError struct {
	StatusCode int
	Body       string
	RetryAfter time.Duration // parsed from Retry-After header; >0 for 429 responses
}

func (e *APIError) Error() string {
	return fmt.Sprintf("langsmith API error (status %d): %s", e.StatusCode, e.Body)
}

// Exporter sends batches of serialized operations to LangSmith.
// It prefers /runs/multipart but falls back to /runs/batch if the
// server returns 404.
type Exporter struct {
	client              *http.Client
	retry               RetryConfig
	logger              logger.Logger
	batchSizeLimitBytes int // max JSON payload per /runs/batch request; 0 uses defaultBatchSizeLimit
	compressionDisabled bool
	multipartDisabled   atomic.Bool
}

// NewExporter creates a new exporter.
func NewExporter(client *http.Client, retry RetryConfig, compressionDisabled bool, l logger.Logger) *Exporter {
	if client == nil {
		client = &http.Client{Timeout: 120 * time.Second}
	}
	if l == nil {
		l = logger.DefaultLogger{}
	}
	return &Exporter{
		client:              client,
		retry:               retry,
		logger:              l,
		batchSizeLimitBytes: defaultBatchSizeLimit,
		compressionDisabled: compressionDisabled,
	}
}

// Export sends a batch of operations to LangSmith. It tries the multipart
// endpoint first; on a 404 it falls back to the JSON batch endpoint and
// disables multipart for all subsequent calls.
func (e *Exporter) Export(ctx context.Context, endpoint models.WriteEndpoint, ops []*models.SerializedOp) error {
	if len(ops) == 0 {
		return nil
	}

	if e.multipartDisabled.Load() {
		return e.exportBatch(ctx, endpoint, ops)
	}

	err := e.exportMultipart(ctx, endpoint, ops)
	if err == nil {
		return nil
	}

	apiErr, ok := err.(*APIError)
	if !ok || apiErr.StatusCode != http.StatusNotFound {
		return err
	}

	e.logger.Warn("multipart endpoint returned 404; falling back to /runs/batch")
	e.multipartDisabled.Store(true)
	return e.exportBatch(ctx, endpoint, ops)
}

// exportMultipart sends ops to POST /runs/multipart as zstd-compressed
// multipart/form-data, with retries on transient failures.
func (e *Exporter) exportMultipart(ctx context.Context, endpoint models.WriteEndpoint, ops []*models.SerializedOp) error {
	var prePayloadBytes int64
	for _, op := range ops {
		prePayloadBytes += int64(op.SizeBytes())
	}

	attempts := max(e.retry.MaxAttempts, 1)
	var lastErr error
	var lastAPIErr *APIError

	for attempt := range attempts {
		if attempt > 0 {
			time.Sleep(e.retry.retryDelay(lastAPIErr, attempt-1))
		}

		boundary := uuid.New().String()
		pr, pw := io.Pipe()
		writeErrCh := make(chan error, 1)
		go func() {
			writeErrCh <- e.writeMultipartBody(pw, boundary, ops)
		}()

		err := e.doMultipartRequest(ctx, endpoint, pr, boundary, prePayloadBytes)
		if writeErr := <-writeErrCh; writeErr != nil {
			return writeErr
		}
		if err == nil {
			return nil
		}
		lastErr = err
		lastAPIErr = nil

		apiErr, ok := err.(*APIError)
		if !ok {
			if attempt < attempts-1 {
				e.logger.Warn("multipart request failed", "attempt", attempt+1, "max_attempts", attempts, "error", err)
			}
			continue
		}
		lastAPIErr = apiErr
		if !isRetryableStatus(apiErr.StatusCode) {
			return apiErr
		}
		if attempt < attempts-1 {
			e.logger.Warn("multipart request returned error status; retrying", "status", apiErr.StatusCode, "attempt", attempt+1, "max_attempts", attempts)
		}
	}
	return lastErr
}

func (e *Exporter) doMultipartRequest(ctx context.Context, endpoint models.WriteEndpoint, pr *io.PipeReader, boundary string, prePayloadBytes int64) error {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint.URL+"/runs/multipart", pr)
	if err != nil {
		pr.CloseWithError(err)
		return fmt.Errorf("create request: %w", err)
	}
	req.Header.Set("Content-Type", "multipart/form-data; boundary="+boundary)
	endpoint.SetAuthHeader(req)
	if !e.compressionDisabled {
		req.Header.Set("Content-Encoding", "zstd")
		req.Header.Set("X-Pre-Compressed-Size", strconv.FormatInt(prePayloadBytes, 10))
	}

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("send multipart: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, errorPreviewSize))
	return &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
		RetryAfter: parseRetryAfter(resp),
	}
}

// batchChunk holds pre-serialized run JSONs grouped by kind, along with a
// running byte count used for splitting oversized payloads.
type batchChunk struct {
	Post  []json.RawMessage
	Patch []json.RawMessage
	bytes int
}

func newBatchChunk() *batchChunk {
	return &batchChunk{
		Post:  make([]json.RawMessage, 0),
		Patch: make([]json.RawMessage, 0),
		bytes: batchJSONOverheadBytes,
	}
}

func (c *batchChunk) add(kind models.OpKind, data json.RawMessage) {
	switch kind {
	case models.OpKindPost:
		c.Post = append(c.Post, data)
	case models.OpKindPatch:
		c.Patch = append(c.Patch, data)
	}
	c.bytes += len(data) + 1 // +1 for JSON comma separator
}

// exportBatch sends ops to POST /runs/batch as JSON: {"post": [...], "patch": [...]},
// splitting oversized payloads into multiple requests, retrying transient failures.
// Attachments are not supported on this endpoint and are silently dropped.
func (e *Exporter) exportBatch(ctx context.Context, endpoint models.WriteEndpoint, ops []*models.SerializedOp) error {
	limit := e.batchSizeLimitBytes
	if limit <= 0 {
		limit = defaultBatchSizeLimit
	}

	var chunks []*batchChunk
	cur := newBatchChunk()

	for _, op := range ops {
		run, err := e.opToRunJSON(op)
		if err != nil {
			return fmt.Errorf("build batch run JSON for %s: %w", op.ID, err)
		}
		runSize := len(run) + 1
		if cur.bytes+runSize > limit && (len(cur.Post) > 0 || len(cur.Patch) > 0) {
			chunks = append(chunks, cur)
			cur = newBatchChunk()
		}
		cur.add(op.Kind, run)
	}
	if len(cur.Post) > 0 || len(cur.Patch) > 0 {
		chunks = append(chunks, cur)
	}

	for i, chunk := range chunks {
		data, err := json.Marshal(struct {
			Post  []json.RawMessage `json:"post"`
			Patch []json.RawMessage `json:"patch"`
		}{Post: chunk.Post, Patch: chunk.Patch})
		if err != nil {
			return fmt.Errorf("marshal batch body (chunk %d/%d): %w", i+1, len(chunks), err)
		}
		if err := e.sendBatchWithRetry(ctx, endpoint, data, i+1, len(chunks)); err != nil {
			return err
		}
	}
	return nil
}

func (e *Exporter) sendBatchWithRetry(ctx context.Context, endpoint models.WriteEndpoint, data []byte, chunkNum, chunkTotal int) error {
	attempts := max(e.retry.MaxAttempts, 1)
	var lastErr error
	var lastAPIErr *APIError

	for attempt := range attempts {
		if attempt > 0 {
			time.Sleep(e.retry.retryDelay(lastAPIErr, attempt-1))
		}

		err := e.doBatchRequest(ctx, endpoint, data)
		if err == nil {
			return nil
		}
		lastErr = err
		lastAPIErr = nil

		apiErr, ok := err.(*APIError)
		if !ok {
			if attempt < attempts-1 {
				e.logger.Warn("batch request failed", "chunk", chunkNum, "chunks", chunkTotal, "attempt", attempt+1, "max_attempts", attempts, "error", err)
			}
			continue
		}
		lastAPIErr = apiErr
		if !isRetryableStatus(apiErr.StatusCode) {
			return apiErr
		}
		if attempt < attempts-1 {
			e.logger.Warn("batch request returned error status; retrying", "status", apiErr.StatusCode, "chunk", chunkNum, "chunks", chunkTotal, "attempt", attempt+1, "max_attempts", attempts)
		}
	}
	return lastErr
}

func (e *Exporter) doBatchRequest(ctx context.Context, endpoint models.WriteEndpoint, data []byte) error {
	req, err := http.NewRequestWithContext(ctx, "POST", endpoint.URL+"/runs/batch", bytes.NewReader(data))
	if err != nil {
		return fmt.Errorf("create batch request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")
	endpoint.SetAuthHeader(req)

	resp, err := e.client.Do(req)
	if err != nil {
		return fmt.Errorf("send batch: %w", err)
	}
	defer func() {
		io.Copy(io.Discard, resp.Body)
		resp.Body.Close()
	}()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return nil
	}

	bodyBytes, _ := io.ReadAll(io.LimitReader(resp.Body, errorPreviewSize))
	return &APIError{
		StatusCode: resp.StatusCode,
		Body:       string(bodyBytes),
		RetryAfter: parseRetryAfter(resp),
	}
}

// opToRunJSON reconstructs a single run JSON object from a SerializedOp by
// merging RunInfo with the split-out fields. Attachments are dropped with a warning.
func (e *Exporter) opToRunJSON(op *models.SerializedOp) (json.RawMessage, error) {
	var m map[string]json.RawMessage
	if err := json.Unmarshal(op.RunInfo, &m); err != nil {
		return nil, fmt.Errorf("unmarshal run info: %w", err)
	}
	if op.Inputs != nil {
		m["inputs"] = op.Inputs
	}
	if op.Outputs != nil {
		m["outputs"] = op.Outputs
	}
	if op.Events != nil {
		m["events"] = op.Events
	}
	if op.Extra != nil {
		m["extra"] = op.Extra
	}
	if op.Error != nil {
		m["error"] = op.Error
	}
	if op.Serialized != nil {
		m["serialized"] = op.Serialized
	}
	if len(op.Attachments) > 0 {
		e.logger.Warn("attachments are not supported in batch mode; dropping", "count", len(op.Attachments), "run_id", op.ID)
	}
	return json.Marshal(m)
}

// writeMultipartBody streams multipart form data into pw.
// When compression is enabled (default), the entire body is zstd-compressed.
// When disabled via LANGSMITH_DISABLE_RUN_COMPRESSION the body is written uncompressed.
func (e *Exporter) writeMultipartBody(pw *io.PipeWriter, boundary string, ops []*models.SerializedOp) (err error) {
	defer func() {
		if writerErr := pw.CloseWithError(err); err == nil && writerErr != nil {
			err = writerErr
		}
	}()

	var dest io.Writer = pw
	if !e.compressionDisabled {
		zw, zErr := zstd.NewWriter(pw, zstd.WithEncoderLevel(zstd.SpeedDefault))
		if zErr != nil {
			return fmt.Errorf("create zstd writer: %w", zErr)
		}
		defer func() {
			if zclose := zw.Close(); err == nil && zclose != nil {
				err = zclose
			}
		}()
		dest = zw
	}

	mw := multipart.NewWriter(dest)
	if err = mw.SetBoundary(boundary); err != nil {
		return fmt.Errorf("set boundary: %w", err)
	}
	defer func() {
		if mclose := mw.Close(); err == nil && mclose != nil {
			err = mclose
		}
	}()

	writeJSONPart := func(name string, payload []byte) error {
		if len(payload) == 0 {
			return nil
		}
		h := textproto.MIMEHeader{}
		h.Set("Content-Disposition", `form-data; name="`+name+`"`)
		h.Set("Content-Type", "application/json")
		h.Set("Content-Length", strconv.Itoa(len(payload)))
		part, partErr := mw.CreatePart(h)
		if partErr != nil {
			return partErr
		}
		_, wErr := part.Write(payload)
		return wErr
	}

	writeBinaryPart := func(name, contentType string, data []byte) error {
		if len(data) == 0 {
			return nil
		}
		h := textproto.MIMEHeader{}
		h.Set("Content-Disposition", `form-data; name="`+name+`"`)
		h.Set("Content-Type", contentType)
		h.Set("Content-Length", strconv.Itoa(len(data)))
		part, partErr := mw.CreatePart(h)
		if partErr != nil {
			return partErr
		}
		_, wErr := part.Write(data)
		return wErr
	}

	for _, op := range ops {
		if len(op.RunInfo) == 0 {
			return fmt.Errorf("missing run info for run_id=%s", op.ID)
		}
		id := op.ID.String()
		prefix := string(op.Kind) + "." + id

		if err = writeJSONPart(prefix, op.RunInfo); err != nil {
			return err
		}
		if err = writeJSONPart(prefix+".inputs", op.Inputs); err != nil {
			return err
		}
		if err = writeJSONPart(prefix+".outputs", op.Outputs); err != nil {
			return err
		}
		if err = writeJSONPart(prefix+".events", op.Events); err != nil {
			return err
		}
		if err = writeJSONPart(prefix+".extra", op.Extra); err != nil {
			return err
		}
		if err = writeJSONPart(prefix+".error", op.Error); err != nil {
			return err
		}
		if err = writeJSONPart(prefix+".serialized", op.Serialized); err != nil {
			return err
		}

		for name, att := range op.Attachments {
			partName := "attachment." + id + "." + name
			ct := att.ContentType
			if ct == "" {
				ct = "application/octet-stream"
			}
			ct += "; length=" + strconv.Itoa(len(att.Data))
			if err = writeBinaryPart(partName, ct, att.Data); err != nil {
				return err
			}
		}
	}

	return nil
}
