package langsmith

import (
	"errors"
	"testing"
)

func TestNormalizeSandboxRunParamsDefaults(t *testing.T) {
	params, timeout, err := normalizeSandboxRunParams(SandboxBoxRunParams{Command: String("echo ok")})
	if err != nil {
		t.Fatalf("normalizeSandboxRunParams returned error: %v", err)
	}
	if timeout != defaultSandboxCommandTimeoutSeconds {
		t.Fatalf("unexpected timeout: %d", timeout)
	}
	if params.Timeout.Value != defaultSandboxCommandTimeoutSeconds {
		t.Fatalf("expected default timeout field, got %d", params.Timeout.Value)
	}
	if params.Shell.Value != defaultSandboxCommandShell {
		t.Fatalf("expected default shell, got %q", params.Shell.Value)
	}
}

func TestNormalizeSandboxCommandStartParamsDefaultsAndMissingCommand(t *testing.T) {
	payload, err := normalizeSandboxCommandStartParams(SandboxCommandStartParams{Command: String("echo ok")})
	if err != nil {
		t.Fatalf("normalizeSandboxCommandStartParams returned error: %v", err)
	}
	if payload.Type.Value != "execute" {
		t.Fatalf("unexpected payload type: %q", payload.Type.Value)
	}
	if payload.TimeoutSeconds.Value != defaultSandboxCommandTimeoutSeconds {
		t.Fatalf("unexpected timeout: %d", payload.TimeoutSeconds.Value)
	}
	if payload.Shell.Value != defaultSandboxCommandShell {
		t.Fatalf("unexpected shell: %q", payload.Shell.Value)
	}
	if payload.IdleTimeoutSeconds.Value != defaultSandboxCommandIdleTimeout {
		t.Fatalf("unexpected idle timeout: %d", payload.IdleTimeoutSeconds.Value)
	}
	if payload.TTLSeconds.Value != defaultSandboxCommandTTLSeconds {
		t.Fatalf("unexpected ttl: %d", payload.TTLSeconds.Value)
	}

	if _, err := normalizeSandboxCommandStartParams(SandboxCommandStartParams{}); err == nil {
		t.Fatal("expected missing command error")
	}
	if _, _, err := normalizeSandboxRunParams(SandboxBoxRunParams{}); err == nil {
		t.Fatal("expected missing run command error")
	}
}

func TestSandboxErrorFromWSMessage(t *testing.T) {
	timeoutErr := sandboxErrorFromWSMessage(sandboxWSMessage{
		Type:      "error",
		ErrorType: "CommandTimeout",
		Error:     "deadline exceeded",
	}, "cmd-1")
	var commandTimeout *SandboxCommandTimeoutError
	if !errors.As(timeoutErr, &commandTimeout) {
		t.Fatalf("expected SandboxCommandTimeoutError, got %T: %v", timeoutErr, timeoutErr)
	}

	notFoundErr := sandboxErrorFromWSMessage(sandboxWSMessage{
		Type:      "error",
		ErrorType: "CommandNotFound",
	}, "cmd-1")
	if notFoundErr.Error() != "command not found: cmd-1 [CommandNotFound]" {
		t.Fatalf("unexpected command not found error: %v", notFoundErr)
	}
}
