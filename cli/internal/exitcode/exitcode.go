// Package exitcode maps CLI errors to the detailed process exit codes
// documented in the README and the agent reference, so scripts and agents can
// tell a retryable failure from one they must fix first.
package exitcode

import (
	"context"
	"errors"
	"net"
	"net/http"
	"net/url"
	"os"
	"strings"

	"github.com/langchain-ai/langsmith-go/cli/internal/client"
)

// EnvVar selects the exit codes. Until DefaultIn, every failure exits Input
// unless it is set to "detailed".
const EnvVar = "LANGSMITH_EXIT_CODES"

// DefaultIn is the major version that makes the detailed codes the default.
const DefaultIn = "1.0.0"

// Detailed reports whether EnvVar opts in to the detailed codes.
func Detailed() bool {
	return strings.EqualFold(strings.TrimSpace(os.Getenv(EnvVar)), "detailed")
}

// Name returns a short description of a detailed exit code.
func Name(code int) string {
	switch code {
	case OK:
		return "success"
	case Auth:
		return "authentication"
	case NotFound:
		return "not found"
	case Conflict:
		return "conflict or limit"
	case Timeout:
		return "--wait timeout"
	case Unavailable:
		return "server error or unreachable"
	case RateLimited:
		return "rate limited"
	default:
		return "bad input"
	}
}

// Detailed exit codes. Every command uses them; `sandboxes exec` is the
// exception and exits with the remote command's own code.
const (
	// OK means the command succeeded.
	OK = 0
	// Input means the command line or request was invalid, or the command
	// failed for a reason not covered below. Fix the input before retrying.
	Input = 1
	// Auth means credentials are missing or were rejected (HTTP 401 or 403).
	Auth = 2
	// NotFound means the resource does not exist (HTTP 404 or 410).
	NotFound = 3
	// Conflict means the request is valid but the resource's state or the
	// plan refuses it (HTTP 402 or 409). Don't retry it unchanged.
	Conflict = 4
	// Timeout means a --wait gave up before the operation finished. The
	// operation may still complete.
	Timeout = 5
	// Unavailable means the server failed (HTTP 5xx), timed out the request
	// (HTTP 408), or could not be reached.
	// Retrying later may succeed.
	Unavailable = 6
	// RateLimited means the API rejected the request as too frequent (HTTP
	// 429). Retry with backoff.
	RateLimited = 7
)

// FromHTTPStatus returns the exit code for an API error response status.
func FromHTTPStatus(status int) int {
	switch {
	case status < 400:
		return OK
	case status == http.StatusUnauthorized, status == http.StatusForbidden:
		return Auth
	case status == http.StatusNotFound, status == http.StatusGone:
		return NotFound
	case status == http.StatusPaymentRequired, status == http.StatusConflict:
		return Conflict
	case status == http.StatusTooManyRequests:
		return RateLimited
	case status == http.StatusRequestTimeout, status >= 500:
		return Unavailable
	default:
		return Input
	}
}

// codedError gives an error an explicit exit code.
type codedError struct {
	err  error
	code int
}

func (e *codedError) Error() string { return e.err.Error() }

func (e *codedError) Unwrap() error { return e.err }

// WithCode returns err with an explicit exit code, which For returns in
// preference to anything it would infer from the error.
func WithCode(err error, code int) error {
	if err == nil {
		return nil
	}
	return &codedError{err: err, code: code}
}

// For returns the exit code for err: OK for nil, a code set by WithCode, the
// code for the API response status behind err, NotFound for a lookup that
// matched nothing, Unavailable for a request that failed in transit or timed
// out, and Input otherwise, including a canceled request and an invalid URL.
func For(err error) int {
	if err == nil {
		return OK
	}
	var coded *codedError
	if errors.As(err, &coded) {
		return coded.code
	}
	if status := client.HTTPStatus(err); status >= 400 {
		return FromHTTPStatus(status)
	}
	if errors.Is(err, client.ErrNotFound) {
		return NotFound
	}
	if errors.Is(err, context.Canceled) {
		return Input
	}
	var urlErr *url.Error
	if errors.As(err, &urlErr) && urlErr.Op == "parse" {
		return Input
	}
	var netErr net.Error
	if urlErr != nil || errors.As(err, &netErr) {
		return Unavailable
	}
	return Input
}
