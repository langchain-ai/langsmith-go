package cmd

import (
	"context"
	"net/http"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// fakeBoxWaitClock replaces boxWaitClock with virtual time that advances only
// when the wait loop sleeps, and records each sleep.
func fakeBoxWaitClock(t *testing.T) *[]time.Duration {
	t.Helper()
	orig := boxWaitClock
	t.Cleanup(func() { boxWaitClock = orig })

	var mu sync.Mutex
	now := time.Date(2026, 1, 1, 0, 0, 0, 0, time.UTC)
	sleeps := &[]time.Duration{}
	boxWaitClock.now = func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return now
	}
	boxWaitClock.sleep = func(_ context.Context, d time.Duration) error {
		mu.Lock()
		defer mu.Unlock()
		now = now.Add(d)
		*sleeps = append(*sleeps, d)
		return nil
	}
	return sleeps
}

type fakeBoxServer struct {
	statuses     []string // returned in order; the last repeats
	statusCalls  int
	statusMsg    string
	createCalled bool
	startCalled  bool
}

func (f *fakeBoxServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		write := func(body string) {
			_, err := w.Write([]byte(body))
			require.NoError(t, err)
		}
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/boxes":
			f.createCalled = true
			write(`{"id":"box-id","name":"my-vm","status":"provisioning"}`)
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/boxes/my-vm/start":
			f.startCalled = true
			write(`{"id":"box-id","name":"my-vm","status":"starting"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sandboxes/boxes/my-vm/status":
			status := f.statuses[min(f.statusCalls, len(f.statuses)-1)]
			f.statusCalls++
			write(`{"status":"` + status + `","status_message":"` + f.statusMsg + `"}`)
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sandboxes/boxes/my-vm":
			write(`{"id":"box-id","name":"my-vm","status":"ready"}`)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	}
}

func runSandboxCmd(t *testing.T, srv *fakeBoxServer, args ...string) (string, error) {
	t.Helper()
	ts := newTestServer(t, srv.handler(t))
	return executeCommand(t, append([]string{"--api-key", "test-key", "--api-url", ts.URL, "--format", "json", "sandboxes"}, args...)...)
}

func TestSandboxCreate_WithoutWaitDoesNotPoll(t *testing.T) {
	fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"provisioning"}}

	out, err := runSandboxCmd(t, srv, "create", "my-vm")

	require.NoError(t, err)
	assert.True(t, srv.createCalled)
	assert.Zero(t, srv.statusCalls)
	assert.Contains(t, out, `"provisioning"`)
}

func TestSandboxCreate_WaitPollsWithBackoffUntilReady(t *testing.T) {
	sleeps := fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"provisioning", "provisioning", "provisioning", "provisioning", "provisioning", "provisioning", "ready"}}

	out, err := runSandboxCmd(t, srv, "create", "my-vm", "--wait")

	require.NoError(t, err)
	assert.Equal(t, 7, srv.statusCalls)
	assert.Equal(t, []time.Duration{
		500 * time.Millisecond, time.Second, 2 * time.Second, 4 * time.Second, 5 * time.Second, 5 * time.Second,
	}, *sleeps)
	assert.Contains(t, out, `"ready"`)
}

func TestSandboxCreate_WaitTimeoutExitsWithNextStep(t *testing.T) {
	sleeps := fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"provisioning"}}

	_, err := runSandboxCmd(t, srv, "create", "my-vm", "--wait", "--timeout", "10s")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out after 10s waiting for sandbox my-vm")
	assert.Contains(t, err.Error(), "last status: provisioning")
	assert.Contains(t, err.Error(), "langsmith sandboxes get my-vm --format json")
	var total time.Duration
	for _, d := range *sleeps {
		total += d
	}
	assert.Equal(t, 10*time.Second, total, "never sleeps past the deadline")
}

func TestSandboxCreate_WaitReportsFailure(t *testing.T) {
	fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"provisioning", "failed"}, statusMsg: "image pull failed"}

	_, err := runSandboxCmd(t, srv, "create", "my-vm", "--wait", "--workspace", "ws-test")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "sandbox my-vm entered failed state: image pull failed")
	assert.Contains(t, err.Error(), "langsmith sandboxes get my-vm")
	assert.Regexp(t, "langsmith sandboxes delete my-vm --api-url \\S+ --workspace ws-test`", err.Error())
	assert.NotContains(t, err.Error(), "test-key")
}

func TestSandboxStart_WaitFailureDoesNotSuggestDelete(t *testing.T) {
	fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"starting", "failed"}}

	_, err := runSandboxCmd(t, srv, "start", "my-vm")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "sandbox my-vm entered failed state")
	assert.NotContains(t, err.Error(), "sandboxes delete")
}

func TestSandboxCreate_TimeoutRequiresWait(t *testing.T) {
	srv := &fakeBoxServer{statuses: []string{"ready"}}

	_, err := runSandboxCmd(t, srv, "create", "my-vm", "--timeout", "1m")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--timeout only applies when waiting")
	assert.False(t, srv.createCalled)
}

func TestSandboxStart_WaitsByDefault(t *testing.T) {
	fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"starting", "ready"}}

	out, err := runSandboxCmd(t, srv, "start", "my-vm")

	require.NoError(t, err)
	assert.Equal(t, 2, srv.statusCalls)
	assert.Contains(t, out, "Sandbox my-vm started")
}

func TestSandboxStart_NoWaitReturnsImmediately(t *testing.T) {
	fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"starting"}}

	out, err := runSandboxCmd(t, srv, "start", "my-vm", "--wait=false")

	require.NoError(t, err)
	assert.True(t, srv.startCalled)
	assert.Zero(t, srv.statusCalls)
	assert.Contains(t, out, "Sandbox my-vm starting")
}

func TestSandboxStart_TimeoutExitsWithNextStep(t *testing.T) {
	fakeBoxWaitClock(t)
	srv := &fakeBoxServer{statuses: []string{"starting"}}

	_, err := runSandboxCmd(t, srv, "start", "my-vm", "--timeout", "30s")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out after 30s waiting for sandbox my-vm")
	assert.Contains(t, err.Error(), "langsmith sandboxes get my-vm --format json")
}

func TestSandboxStart_TimeoutBoundsASlowStatusRequest(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v2/sandboxes/boxes/my-vm/status" {
			select {
			case <-r.Context().Done():
				return
			case <-time.After(5 * time.Second):
			}
			_, _ = w.Write([]byte(`{"status":"ready"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"starting"}`))
	})

	start := time.Now()
	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "--format", "json",
		"sandboxes", "start", "my-vm", "--timeout", "200ms")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "timed out after 200ms waiting for sandbox my-vm")
	assert.Less(t, time.Since(start), 3*time.Second)
}

func TestSandboxStart_RejectsNonPositiveTimeout(t *testing.T) {
	srv := &fakeBoxServer{statuses: []string{"ready"}}

	_, err := runSandboxCmd(t, srv, "start", "my-vm", "--timeout", "0s")

	require.Error(t, err)
	assert.Contains(t, err.Error(), "--timeout must be greater than 0")
	assert.False(t, srv.startCalled)
}
