package address

import (
	"errors"
	"strings"
	"testing"
)

func TestValidation(t *testing.T) {
	long := strings.Repeat("x", 64)
	for name, f := range map[string]func() (APIAddress, error){
		"agent": func() (APIAddress, error) {
			a, err := NewAgentAddress("support", "Staging")
			return a.ToAPIAddress(), err
		},
		"experiment": func() (APIAddress, error) {
			a, err := NewExperimentAddress("0190C3D4-0000-7000-8000-0000000000B1")
			return a.ToAPIAddress(), err
		},
	} {
		got, err := f()
		if err != nil || (got.Kind == "AGENT") != (name == "agent") || (name == "agent" && got.Environment != "STAGING") {
			t.Errorf("%s: %+v, %v", name, got, err)
		}
	}
	if got := NewEvaluatorAddress().ToAPIAddress(); got.Kind != "EVALUATOR" {
		t.Errorf("evaluator: %+v", got)
	}
	bad := map[string]func(string) error{
		"Support":       func(s string) error { _, err := NewAgentAddress(s, "local"); return err },
		"1abc":          func(s string) error { _, err := NewAgentAddress(s, "local"); return err },
		"abc-":          func(s string) error { _, err := NewAgentAddress(s, "local"); return err },
		long:            func(s string) error { _, err := NewAgentAddress(s, "local"); return err },
		"zzenv":         func(s string) error { _, err := NewAgentAddress("support", s); return err },
		"not-a-uuid-zz": func(s string) error { _, err := NewExperimentAddress(s); return err },
	}
	for in, f := range bad {
		err := f(in)
		if !errors.Is(err, ErrInvalidAddress) || strings.Contains(err.Error(), in) {
			t.Errorf("%q: %v", in, err)
		}
	}
}

func TestFromEnv(t *testing.T) {
	for _, c := range []struct {
		id, env string
		ok      bool
	}{{"", "", true}, {"support", "local", true}, {"support", "", false}, {"", "local", false}, {"Bad_Id", "local", false}} {
		t.Setenv(EnvAgentID, c.id)
		t.Setenv(EnvAgentEnvironment, c.env)
		_, err := FromEnv()
		if c.ok != (err == nil) || (err != nil && !errors.Is(err, ErrEnvAddress)) {
			t.Errorf("%q/%q: %v", c.id, c.env, err)
		}
	}
}
