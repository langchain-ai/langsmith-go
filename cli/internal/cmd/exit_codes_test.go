package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strings"
	"testing"

	langsmith "github.com/langchain-ai/langsmith-go"

	"github.com/langchain-ai/langsmith-go/cli/internal/exitcode"
)

// exitCalled stops a command at its exitProcess call.
type exitCalled struct{ code int }

// runExpectingExit runs a command that reports its failure through
// ExitErrorf and returns the exit code it asked for.
func runExpectingExit(t *testing.T, args ...string) (code int) {
	t.Helper()
	prev := exitProcess
	exitProcess = func(code int) { panic(exitCalled{code}) }
	t.Cleanup(func() { exitProcess = prev })
	defer func() {
		r := recover()
		exit, ok := r.(exitCalled)
		if !ok {
			panic(r)
		}
		code = exit.code
	}()
	_, err := executeCommand(t, args...)
	t.Fatalf("command did not exit; err = %v", err)
	return 0
}

func TestHandWrittenCommandsMapAPIErrorsToExitCodes(t *testing.T) {
	tests := map[int]int{
		http.StatusBadRequest:          exitcode.Input,
		http.StatusUnauthorized:        exitcode.Auth,
		http.StatusForbidden:           exitcode.Auth,
		http.StatusNotFound:            exitcode.NotFound,
		http.StatusConflict:            exitcode.Conflict,
		http.StatusUnprocessableEntity: exitcode.Input,
	}
	t.Setenv(exitcode.EnvVar, "detailed")
	for status, want := range tests {
		t.Run(http.StatusText(status), func(t *testing.T) {
			server := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(status)
				_, _ = w.Write([]byte(`{"detail":"rejected"}`))
			})
			if got := runExpectingExit(t, "--api-url", server.URL, "--api-key", "test-key", "projects", "list"); got != want {
				t.Errorf("exit code = %d, want %d", got, want)
			}
		})
	}
}

func TestExitErrorWithoutAnErrorIsInput(t *testing.T) {
	prev := exitProcess
	var got int
	exitProcess = func(code int) { got = code }
	t.Cleanup(func() { exitProcess = prev })
	ExitErrorf("--limit must be at most %d", 100)
	if got != exitcode.Input {
		t.Errorf("exit code = %d, want %d", got, exitcode.Input)
	}
}

func TestLegacyExitCodesWarnAboutTheDetailedCode(t *testing.T) {
	notFound := &langsmith.Error{StatusCode: http.StatusNotFound}

	t.Setenv(exitcode.EnvVar, "")
	var text bytes.Buffer
	if got := processExitCode(&text, notFound, false); got != exitcode.Input {
		t.Errorf("exit code = %d, want %d", got, exitcode.Input)
	}
	if !strings.Contains(text.String(), "from 1.0.0 exits 3 (not found)") || !strings.Contains(text.String(), "LANGSMITH_EXIT_CODES=detailed") {
		t.Errorf("warning = %q", text.String())
	}

	var js bytes.Buffer
	processExitCode(&js, notFound, true)
	var warning map[string]any
	if err := json.Unmarshal(js.Bytes(), &warning); err != nil {
		t.Fatalf("warning is not JSON: %q", js.String())
	}
	if warning["deprecated"] != true || warning["kind"] != "exit_code" || warning["replacement"] != "3" {
		t.Errorf("warning = %v", warning)
	}

	var none bytes.Buffer
	if got := processExitCode(&none, errors.New("--limit must be positive"), false); got != exitcode.Input || none.Len() != 0 {
		t.Errorf("bad input: exit code %d, warning %q; want %d and no warning", got, none.String(), exitcode.Input)
	}

	t.Setenv(exitcode.EnvVar, "detailed")
	var detailed bytes.Buffer
	if got := processExitCode(&detailed, notFound, false); got != exitcode.NotFound || detailed.Len() != 0 {
		t.Errorf("detailed: exit code %d, warning %q; want %d and no warning", got, detailed.String(), exitcode.NotFound)
	}
}

func TestGeneratedCommandErrorsWarnInJSON(t *testing.T) {
	isolateGeneratedAuth(t)
	writeGeneratedConfig(t, `{}`)
	t.Setenv(exitcode.EnvVar, "")
	_, _, err := runGeneratedCommand(t, "", "prompt-webhooks", "list", "--format", "json")
	if err == nil {
		t.Fatal("expected a not-authenticated error")
	}
	var w bytes.Buffer
	if got := processExitCode(&w, err, GetFormat() == "json" || generatedJSON); got != exitcode.Input {
		t.Errorf("exit code = %d, want %d", got, exitcode.Input)
	}
	var warning map[string]any
	if err := json.Unmarshal(w.Bytes(), &warning); err != nil || warning["replacement"] != "2" {
		t.Errorf("warning = %q, want a JSON warning naming exit code 2", w.String())
	}
}
