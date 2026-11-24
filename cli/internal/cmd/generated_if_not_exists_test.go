package cmd

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

// agentsConflictServer answers agent creates with 409 and serves one existing
// agent, recording each request as "METHOD path".
func agentsConflictServer(t *testing.T, requests *[]string) *httptest.Server {
	t.Helper()
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		*requests = append(*requests, r.Method+" "+r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v1/agents":
			w.Header().Set("retry-after-ms", "1")
			w.WriteHeader(http.StatusConflict)
			_, _ = w.Write([]byte(`{"title":"Conflict","status":409,"detail":"An Agent with this id already exists"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v1/agents/existing-agent":
			_, _ = w.Write([]byte(`{"id":"existing-agent","name":"Existing","source_type":"EXTERNAL","agent_type":"CUSTOM","created_at":"2026-10-01T00:00:00Z","environments":[]}`))
		default:
			w.WriteHeader(http.StatusNotFound)
			_, _ = w.Write([]byte(`{"title":"Not Found","status":404}`))
		}
	}))
	t.Cleanup(ts.Close)
	return ts
}

func TestGeneratedAgentsCreateIfNotExists(t *testing.T) {
	isolateGeneratedAuth(t)
	for _, tt := range []struct {
		name         string
		args         []string
		wantExit     int
		wantStdout   string
		wantStderr   string
		wantRetrieve string
	}{
		{
			name:         "returns the existing agent",
			args:         []string{"--id", "existing-agent", "--source-type", "EXTERNAL", "--production-only", "--if-not-exists", "--format", "json", "--name", "New"},
			wantStdout:   `"name": "Existing"`,
			wantRetrieve: "GET /api/v1/agents/existing-agent",
		},
		{
			name:         "takes --id=value",
			args:         []string{"--id=existing-agent", "--source-type=EXTERNAL", "--if-not-exists", "--transform", "id", "--format", "raw"},
			wantStdout:   "existing-agent",
			wantRetrieve: "GET /api/v1/agents/existing-agent",
		},
		{
			name:         "reports the conflict when the lookup finds nothing",
			args:         []string{"--id", "other-agent", "--source-type", "EXTERNAL", "--if-not-exists"},
			wantExit:     1,
			wantStderr:   "409 Conflict",
			wantRetrieve: "GET /api/v1/agents/other-agent",
		},
		{
			name:       "reports the conflict without the flag",
			args:       []string{"--id", "existing-agent", "--source-type", "EXTERNAL"},
			wantExit:   1,
			wantStderr: "409 Conflict",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			var requests []string
			ts := agentsConflictServer(t, &requests)
			args := append([]string{"--api-url", ts.URL, "--api-key", "k", "agents", "create"}, tt.args...)
			stdout, stderr, exitCode, err := runGeneratedCommandOutput(t, args...)
			if err != nil || exitCode != tt.wantExit {
				t.Fatalf("err = %v, exit code = %d, want %d\nstderr: %s", err, exitCode, tt.wantExit, stderr)
			}
			if !strings.Contains(stdout, tt.wantStdout) {
				t.Errorf("stdout = %q, want it to contain %q", stdout, tt.wantStdout)
			}
			if !strings.Contains(stderr, tt.wantStderr) {
				t.Errorf("stderr = %q, want it to contain %q", stderr, tt.wantStderr)
			}
			var retrieves []string
			for _, r := range requests {
				if strings.HasPrefix(r, "GET ") {
					retrieves = append(retrieves, r)
				}
			}
			if want := tt.wantRetrieve; (want == "" && len(retrieves) != 0) || (want != "" && (len(retrieves) != 1 || retrieves[0] != want)) {
				t.Errorf("lookups = %v, want %q", retrieves, want)
			}
		})
	}
}

func TestGeneratedIfNotExistsNeedsLiteralID(t *testing.T) {
	isolateGeneratedAuth(t)
	for _, args := range [][]string{
		{"--id", "@id.txt", "--source-type", "EXTERNAL", "--if-not-exists"},
		{"--source-type", "EXTERNAL", "--if-not-exists"},
	} {
		var requests []string
		ts := agentsConflictServer(t, &requests)
		_, _, _, err := runGeneratedCommandOutput(t, append([]string{"--api-url", ts.URL, "--api-key", "k", "agents", "create"}, args...)...)
		if err == nil || !strings.Contains(err.Error(), "--if-not-exists needs the id as a literal --id value") {
			t.Errorf("%v: err = %v, want it to ask for a literal --id", args, err)
		}
		if len(requests) != 0 {
			t.Errorf("%v: sent %v, want no requests", args, requests)
		}
	}
}

func TestGeneratedIfNotExistsRejectedOnOtherOperations(t *testing.T) {
	isolateGeneratedAuth(t)
	var requests []string
	ts := agentsConflictServer(t, &requests)
	stderr, exitCode, err := runGeneratedCommand(t, "", "--api-url", ts.URL, "--api-key", "k",
		"agents", "retrieve", "--agent-id", "existing-agent", "--if-not-exists")
	if err == nil && exitCode == 0 {
		t.Fatalf("agents retrieve accepted --if-not-exists; stderr: %s", stderr)
	}
	if len(requests) != 0 {
		t.Errorf("sent %v, want no requests", requests)
	}
}

// runGeneratedCommandOutput is runGeneratedCommand that also returns stdout.
func runGeneratedCommandOutput(t *testing.T, args ...string) (stdout, stderr string, exitCode int, err error) {
	t.Helper()
	prevExit := exitGenerated
	exitGenerated = func(code int) { exitCode = code }
	t.Cleanup(func() { exitGenerated = prevExit })

	root := NewRootCmd("dev", "dev")
	var errBuf bytes.Buffer
	root.SetArgs(args)
	root.SetErr(&errBuf)
	stdout = captureStdout(t, func() { err = root.Execute() })
	return stdout, errBuf.String(), exitCode, err
}
