package langsmith

import (
	"context"
	"testing"
	"testing/synctest"
	"time"

	"github.com/langchain-ai/langsmith-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// stubSandboxCommandAttempts swaps the attempt func for the duration of a test,
// recording each attempt's connect timeout and burning it on the fake clock.
func stubSandboxCommandAttempts(t *testing.T, err error) *[]time.Duration {
	t.Helper()
	timeouts := new([]time.Duration)
	original := startSandboxCommandAttemptFn
	t.Cleanup(func() { startSandboxCommandAttemptFn = original })
	startSandboxCommandAttemptFn = func(ctx context.Context, _ string, _ sandboxCommandStartRequest, connectTimeout time.Duration, _ ...option.RequestOption) (*SandboxCommandHandle, error) {
		*timeouts = append(*timeouts, connectTimeout)
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case <-time.After(connectTimeout):
		}
		return nil, err
	}
	return timeouts
}

func startCommandForTiming(ctx context.Context) error {
	svc := &SandboxBoxService{}
	_, err := svc.startCommandWithDataplaneURL(
		ctx,
		"https://dataplane.test/sb-1",
		SandboxCommandStartParams{Command: String("echo hi")},
		SandboxCommandCallbacks{},
	)
	return err
}

func TestSandboxCommandStartConnectBudget(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		start := time.Now()
		timeouts := stubSandboxCommandAttempts(t, &SandboxConnectTimeoutError{Message: "timed out"})

		err := startCommandForTiming(context.Background())

		var connErr *SandboxConnectTimeoutError
		require.ErrorAs(t, err, &connErr)

		assert.Equal(t, sandboxCommandConnectBudget, time.Since(start),
			"a blackholed handshake must cost exactly the budget, not attempts x timeout")

		// 30s dials at t=0, 30.5s, 61.5s, then a final dial clamped to the
		// 26.5s left; the 0.5+1+2s backoffs make up the rest of the budget.
		require.Equal(t, []time.Duration{
			sandboxCommandConnectTimeout,
			sandboxCommandConnectTimeout,
			sandboxCommandConnectTimeout,
			26500 * time.Millisecond,
		}, *timeouts)
		var total time.Duration
		for _, d := range *timeouts {
			assert.Positive(t, d, "a zero timeout would disable the dial deadline entirely")
			assert.LessOrEqual(t, d, sandboxCommandConnectTimeout)
			total += d
		}
		assert.Less(t, total, sandboxCommandConnectBudget,
			"backoff draws from the same budget, so dial time alone must be under it")
	})
}

func TestSandboxCommandStartStopsAtRetryLimitInsideBudget(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		// Attempts that fail instantly exhaust the retry count long before the
		// budget, so the retry limit is what ends the loop.
		original := startSandboxCommandAttemptFn
		t.Cleanup(func() { startSandboxCommandAttemptFn = original })
		var attempts int
		startSandboxCommandAttemptFn = func(_ context.Context, _ string, _ sandboxCommandStartRequest, _ time.Duration, _ ...option.RequestOption) (*SandboxCommandHandle, error) {
			attempts++
			return nil, &SandboxConnectTimeoutError{Message: "refused"}
		}

		start := time.Now()
		require.Error(t, startCommandForTiming(context.Background()))

		assert.Equal(t, sandboxCommandMaxAutoReconnects+1, attempts)
		// 500ms doubling to a 8s cap: 0.5+1+2+4+8.
		assert.Equal(t, 15500*time.Millisecond, time.Since(start))
	})
}

func TestSandboxCommandStartCancellationBeatsBackoff(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		stubSandboxCommandAttempts(t, &SandboxConnectTimeoutError{Message: "timed out"})

		ctx, cancel := context.WithTimeout(context.Background(), 45*time.Second)
		defer cancel()

		start := time.Now()
		err := startCommandForTiming(ctx)

		require.ErrorIs(t, err, context.DeadlineExceeded)
		assert.Equal(t, 45*time.Second, time.Since(start),
			"cancellation must win immediately, not at the end of the budget")
	})
}

func TestSandboxCommandStartDoesNotRetryPermanentError(t *testing.T) {
	synctest.Test(t, func(t *testing.T) {
		timeouts := stubSandboxCommandAttempts(t, &SandboxConnectionError{Message: "bad status"})

		start := time.Now()
		err := startCommandForTiming(context.Background())

		var connErr *SandboxConnectionError
		require.ErrorAs(t, err, &connErr)
		assert.Len(t, *timeouts, 1, "a rejected handshake must not be retried")
		assert.Equal(t, sandboxCommandConnectTimeout, time.Since(start))
	})
}
