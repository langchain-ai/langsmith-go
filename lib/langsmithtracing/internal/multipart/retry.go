package multipart

import (
	"context"
	"math"
	"math/rand/v2"
	"net/http"
	"strconv"
	"time"
)

// RetryConfig controls retry behavior for the multipart exporter.
type RetryConfig struct {
	MaxAttempts int
	BackoffBase time.Duration
	BackoffMax  time.Duration
}

// DefaultRetry returns a conservative retry configuration.
func DefaultRetry() RetryConfig {
	return RetryConfig{
		MaxAttempts: 3,
		BackoffBase: 500 * time.Millisecond,
		BackoffMax:  5 * time.Second,
	}
}

const defaultRetryAfter429 = 10 * time.Second

var retryableStatuses = map[int]bool{
	408: true, // Request Timeout
	425: true, // Too Early
	429: true, // Too Many Requests
	500: true, // Internal Server Error
	502: true, // Bad Gateway
	503: true, // Service Unavailable
	504: true, // Gateway Timeout
}

// isRetryableStatus returns true for HTTP status codes that warrant a retry.
func isRetryableStatus(code int) bool {
	return retryableStatuses[code]
}

// backoff returns the sleep duration for the given zero-based attempt,
// using exponential backoff with full jitter capped at BackoffMax.
func (rc RetryConfig) backoff(attempt int) time.Duration {
	base := float64(rc.BackoffBase) * math.Pow(2, float64(attempt))
	capped := math.Min(base, float64(rc.BackoffMax))
	return time.Duration(rand.Float64() * capped)
}

// retryDelay returns the duration to sleep before the next retry.
// For 429 responses, it uses the Retry-After header (defaulting to 10s).
// For other retryable errors, it uses exponential backoff with jitter.
func (rc RetryConfig) retryDelay(err *APIError, attempt int) time.Duration {
	if err != nil && err.RetryAfter > 0 {
		return err.RetryAfter
	}
	return rc.backoff(attempt)
}

// sleepWithContext sleeps for the given duration but returns early if
// the context is canceled. Returns ctx.Err() on cancellation, nil otherwise.
func sleepWithContext(ctx context.Context, d time.Duration) error {
	t := time.NewTimer(d)
	defer t.Stop()
	select {
	case <-t.C:
		return nil
	case <-ctx.Done():
		return ctx.Err()
	}
}

// parseRetryAfter extracts a Retry-After duration from an HTTP response.
// It supports the delta-seconds format (e.g. "10"). For 429 responses
// with no Retry-After header, it defaults to 10 seconds.
func parseRetryAfter(resp *http.Response) time.Duration {
	val := resp.Header.Get("Retry-After")
	if val != "" {
		if secs, err := strconv.Atoi(val); err == nil && secs > 0 {
			return time.Duration(secs) * time.Second
		}
	}
	if resp.StatusCode == 429 {
		return defaultRetryAfter429
	}
	return 0
}
