package langsmith

import (
	"errors"
	"reflect"
	"testing"
)

func mustAddress(t *testing.T, id, environment string) AgentAddress {
	t.Helper()
	a, err := NewAgentAddress(id, environment)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func setTracerEnv(t *testing.T, env map[string]string) {
	t.Helper()
	for _, name := range []string{"LANGSMITH_PROJECT", "LANGSMITH_AGENT_ID", "LANGSMITH_AGENT_ENVIRONMENT"} {
		t.Setenv(name, env[name])
	}
}

func TestOTelTracerHeaders(t *testing.T) {
	tests := []struct {
		name string
		env  map[string]string
		opts []OTelTracerOption
		want map[string]string
	}{
		{
			name: "defaults to the default project",
			want: map[string]string{"x-api-key": "key", "Langsmith-Project": "default"},
		},
		{
			name: "the project option names the project",
			opts: []OTelTracerOption{WithProjectName("mine")},
			want: map[string]string{"x-api-key": "key", "Langsmith-Project": "mine"},
		},
		{
			name: "the env project names the project",
			env:  map[string]string{"LANGSMITH_PROJECT": "from-env"},
			want: map[string]string{"x-api-key": "key", "Langsmith-Project": "from-env"},
		},
		{
			name: "an address replaces the project header",
			opts: []OTelTracerOption{WithAddress(mustAddress(t, "support-agent", "Staging"))},
			want: map[string]string{
				"x-api-key":                   "key",
				"Langsmith-Agent-Id":          "support-agent",
				"Langsmith-Agent-Environment": "staging",
			},
		},
		{
			name: "the env agent addresses the traces",
			env:  map[string]string{"LANGSMITH_AGENT_ID": "support-agent", "LANGSMITH_AGENT_ENVIRONMENT": "production"},
			want: map[string]string{
				"x-api-key":                   "key",
				"Langsmith-Agent-Id":          "support-agent",
				"Langsmith-Agent-Environment": "production",
			},
		},
		{
			name: "an explicit project beats the env agent",
			env:  map[string]string{"LANGSMITH_AGENT_ID": "support-agent", "LANGSMITH_AGENT_ENVIRONMENT": "production"},
			opts: []OTelTracerOption{WithProjectName("mine")},
			want: map[string]string{"x-api-key": "key", "Langsmith-Project": "mine"},
		},
		{
			name: "an explicit address beats the env project",
			env:  map[string]string{"LANGSMITH_PROJECT": "from-env"},
			opts: []OTelTracerOption{WithAddress(mustAddress(t, "support-agent", "local"))},
			want: map[string]string{
				"x-api-key":                   "key",
				"Langsmith-Agent-Id":          "support-agent",
				"Langsmith-Agent-Environment": "local",
			},
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setTracerEnv(t, test.env)
			cfg, err := resolveConfig(append([]OTelTracerOption{WithAPIKey("key")}, test.opts...))
			if err != nil {
				t.Fatal(err)
			}
			if got := cfg.headers(); !reflect.DeepEqual(got, test.want) {
				t.Fatalf("headers = %v, want %v", got, test.want)
			}
		})
	}
}

func TestOTelTracerRejectsInvalidAddresses(t *testing.T) {
	tests := map[string]AgentAddress{
		"no kind":       {ID: "support-agent", Environment: "LOCAL"},
		"not an agent":  {Kind: "EXPERIMENT", ID: "0190c3d4-0000-7000-8000-0000000000b1"},
		"invalid parts": {Kind: "AGENT", ID: "Bad_Id", Environment: "zz"},
	}
	for name, literal := range tests {
		t.Run(name, func(t *testing.T) {
			setTracerEnv(t, nil)
			_, err := resolveConfig([]OTelTracerOption{WithAPIKey("key"), WithAddress(literal)})
			if !errors.Is(err, ErrInvalidAddress) {
				t.Fatalf("error = %v, want one wrapping ErrInvalidAddress", err)
			}
		})
	}
}

func TestOTelTracerRefusesAmbiguousDestinations(t *testing.T) {
	tests := []struct {
		name    string
		env     map[string]string
		opts    []OTelTracerOption
		wantEnv bool
	}{
		{
			name: "a project and an address",
			opts: []OTelTracerOption{WithProjectName("mine"), WithAddress(mustAddress(t, "support-agent", "staging"))},
		},
		{
			name:    "an env that names a project and an agent",
			env:     map[string]string{"LANGSMITH_PROJECT": "p", "LANGSMITH_AGENT_ID": "support-agent", "LANGSMITH_AGENT_ENVIRONMENT": "staging"},
			wantEnv: true,
		},
		{
			name:    "an env that names only the agent id",
			env:     map[string]string{"LANGSMITH_AGENT_ID": "support-agent"},
			wantEnv: true,
		},
		{
			name:    "an env that names only the environment",
			env:     map[string]string{"LANGSMITH_AGENT_ENVIRONMENT": "staging"},
			wantEnv: true,
		},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			setTracerEnv(t, test.env)
			_, err := resolveConfig(append([]OTelTracerOption{WithAPIKey("key")}, test.opts...))
			if err == nil {
				t.Fatal("expected an error")
			}
			if test.wantEnv && !errors.Is(err, ErrEnvAddress) {
				t.Fatalf("error = %v, want one wrapping ErrEnvAddress", err)
			}
		})
	}
}
