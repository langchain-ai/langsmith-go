package address

import (
	"encoding/json"
	"errors"
	"strings"
	"testing"
)

func TestJSONRoundTrip(t *testing.T) {
	agent, _ := NewAgentAddress("support", "Staging")
	exp, _ := NewExperimentAddress("0190C3D4-0000-7000-8000-0000000000B1")
	for _, c := range []struct {
		name string
		v    any
		dst  any
		want string
	}{
		{"agent", agent, new(AgentAddress), `{"kind":"AGENT","id":"support","environment":"STAGING"}`},
		{"experiment", exp, new(ExperimentAddress), `{"kind":"EXPERIMENT","id":"0190c3d4-0000-7000-8000-0000000000b1"}`},
		{"evaluator", NewEvaluatorAddress(), new(EvaluatorAddress), `{"kind":"EVALUATOR"}`},
	} {
		got, err := json.Marshal(c.v)
		if err != nil || string(got) != c.want {
			t.Errorf("%s: marshal = %s, %v", c.name, got, err)
		}
		if err := json.Unmarshal([]byte(c.want), c.dst); err != nil {
			t.Errorf("%s: unmarshal: %v", c.name, err)
		}
	}
	var a AgentAddress
	if err := json.Unmarshal([]byte(`{"kind":"AGENT","id":"support","environment":"production"}`), &a); err != nil || a.Environment != "PRODUCTION" {
		t.Errorf("decode = %+v, %v", a, err)
	}
	for _, bad := range []string{`{"kind":"AGENT","id":"Bad_Id","environment":"LOCAL"}`, `{"kind":"EXPERIMENT","id":"0190c3d4-0000-7000-8000-0000000000b1"}`, `{"id":"support","environment":"LOCAL"}`, `[]`} {
		if err := json.Unmarshal([]byte(bad), &AgentAddress{}); !errors.Is(err, ErrInvalidAddress) {
			t.Errorf("%s: got %v", bad, err)
		}
	}
}

func TestValidation(t *testing.T) {
	agent := func(id, env string) error { _, err := NewAgentAddress(id, env); return err }
	experiment := func(id, _ string) error { _, err := NewExperimentAddress(id); return err }
	for _, c := range []struct {
		in, env string
		f       func(id, env string) error
		ok      bool
	}{
		{"support", "Staging", agent, true},
		{"a", "local", agent, true},
		{"Support", "local", agent, false},
		{"1abc", "local", agent, false},
		{"abc-", "local", agent, false},
		{strings.Repeat("x", 64), "local", agent, false},
		{"support", "zzenv", agent, false},
		{"0190C3D4-0000-7000-8000-0000000000B1", "", experiment, true},
		{"not-a-uuid-zz", "", experiment, false},
	} {
		err := c.f(c.in, c.env)
		if c.ok != (err == nil) || (err != nil && (!errors.Is(err, ErrInvalidAddress) || strings.Contains(err.Error(), c.in))) {
			t.Errorf("%q/%q: %v", c.in, c.env, err)
		}
	}
}

func TestEnsureRejectsLiteralValues(t *testing.T) {
	exp, _ := NewExperimentAddress("0190c3d4-0000-7000-8000-0000000000b1")
	for name, err := range map[string]error{
		"bad id":             ensureErr(AgentAddress{Kind: KindAgent, ID: "Bad_Id", Environment: "LOCAL"}),
		"no kind":            ensureErr(AgentAddress{ID: "support", Environment: "LOCAL"}),
		"non-agent in agent": ensureErr(AgentAddress{Kind: KindExperiment, ID: exp.ID}),
		"bad experiment":     ensureErr(ExperimentAddress{Kind: KindExperiment, ID: "x"}),
		"bad evaluator":      ensureErr(EvaluatorAddress{Kind: "AGENT"}),
	} {
		if !errors.Is(err, ErrInvalidAddress) {
			t.Errorf("%s: got %v", name, err)
		}
	}
	if _, err := EnsureAgent(AgentAddress{Kind: KindAgent, ID: "support", Environment: "Local"}); err != nil {
		t.Errorf("normalizable literal: %v", err)
	}
}

func ensureErr(a Address) error { _, err := Ensure(a); return err }

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

func TestZeroValueRoundTrip(t *testing.T) {
	type wrapper struct {
		Agent      AgentAddress      `json:"agent"`
		Experiment ExperimentAddress `json:"experiment"`
		Evaluator  EvaluatorAddress  `json:"evaluator"`
	}
	data, err := json.Marshal(wrapper{})
	if err != nil {
		t.Fatal(err)
	}
	var got wrapper
	if err := json.Unmarshal(data, &got); err != nil || got != (wrapper{}) {
		t.Errorf("round trip of %s: %+v, %v", data, got, err)
	}
	for _, dst := range []any{&AgentAddress{}, &ExperimentAddress{}, &EvaluatorAddress{}} {
		if err := json.Unmarshal([]byte(`{"kind":""}`), dst); err != nil {
			t.Errorf(`%T: {"kind":""}: %v`, dst, err)
		}
	}
	for _, dst := range []any{&AgentAddress{}, &ExperimentAddress{}} {
		if err := json.Unmarshal([]byte(`{"kind":"","id":"x"}`), dst); !errors.Is(err, ErrInvalidAddress) {
			t.Errorf(`%T: {"kind":"","id":"x"}: got %v`, dst, err)
		}
	}
}

func TestEmptyKindMessage(t *testing.T) {
	for name, err := range map[string]error{
		"agent":      ensureErr(AgentAddress{}),
		"experiment": ensureErr(ExperimentAddress{}),
		"evaluator":  ensureErr(EvaluatorAddress{}),
	} {
		if !errors.Is(err, ErrInvalidAddress) || !strings.Contains(err.Error(), "unset or has no kind") {
			t.Errorf("%s: %v", name, err)
		}
	}
	if err := ensureErr(AgentAddress{Kind: KindExperiment}); err == nil || !strings.Contains(err.Error(), "only an AGENT address can receive traces") {
		t.Errorf("non-agent kind: %v", err)
	}
}
