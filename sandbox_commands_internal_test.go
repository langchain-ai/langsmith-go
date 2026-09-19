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

func TestNormalizeSandboxCommandStartParamsCloseStdin(t *testing.T) {
	tests := []struct {
		name string
		body SandboxCommandStartParams
		want bool
	}{
		{
			name: "defaults on for a non-PTY command",
			body: SandboxCommandStartParams{Command: String("cat")},
			want: true,
		},
		{
			name: "stays off under a PTY",
			body: SandboxCommandStartParams{Command: String("bash"), Pty: Bool(true)},
			want: false,
		},
		{
			name: "stays off under a PTY even when asked for",
			body: SandboxCommandStartParams{Command: String("bash"), Pty: Bool(true), CloseStdin: Bool(true)},
			want: false,
		},
		{
			name: "honours an explicit opt-out",
			body: SandboxCommandStartParams{Command: String("cat"), CloseStdin: Bool(false)},
			want: false,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			out, err := normalizeSandboxCommandStartParams(tc.body)
			if err != nil {
				t.Fatalf("normalizeSandboxCommandStartParams returned error: %v", err)
			}
			if got := out.CloseStdin.Value && out.CloseStdin.Present; got != tc.want {
				t.Fatalf("close_stdin: got %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNormalizeSandboxCommandStartParamsRejectsRunConfigWithDeprecatedFields(t *testing.T) {
	runConfig := F(SandboxRunConfig{User: String("app")})

	for name, body := range map[string]SandboxCommandStartParams{
		"env": {Command: String("echo ok"), Env: F(map[string]string{"A": "1"}), RunConfig: runConfig},
		"cwd": {Command: String("echo ok"), CWD: String("/tmp"), RunConfig: runConfig},
	} {
		t.Run(name, func(t *testing.T) {
			if _, err := normalizeSandboxCommandStartParams(body); err == nil {
				t.Fatal("expected an error combining RunConfig with the deprecated fields")
			}
		})
	}
}

func TestNormalizeSandboxCommandStartParamsForwardsRunConfig(t *testing.T) {
	out, err := normalizeSandboxCommandStartParams(SandboxCommandStartParams{
		Command:   String("echo ok"),
		RunConfig: F(SandboxRunConfig{User: String("app"), WorkDir: String("/workspace")}),
	})
	if err != nil {
		t.Fatalf("normalizeSandboxCommandStartParams returned error: %v", err)
	}
	if !out.RunConfig.Present {
		t.Fatal("expected run_config to be forwarded")
	}
	if got := out.RunConfig.Value.User.Value; got != "app" {
		t.Fatalf("unexpected run_config user: %q", got)
	}
}
