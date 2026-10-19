// Package address implements the beta tracing addresses, re-exported by langsmith.
package address

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"
)

// ErrInvalidAddress is wrapped by invalid-value errors, which never echo the input.
var ErrInvalidAddress = errors.New("invalid address")

// ErrEnvAddress is wrapped by the error for env vars that cannot address a run.
var ErrEnvAddress = errors.New("the LANGSMITH_AGENT_* env vars cannot address a run")

const (
	EnvAgentID          = "LANGSMITH_AGENT_ID"
	EnvAgentEnvironment = "LANGSMITH_AGENT_ENVIRONMENT"
)

// The server's agent id rule: a DNS label.
var agentIDPattern = regexp.MustCompile(`^[a-z](?:[a-z0-9-]{0,61}[a-z0-9])?$`)

var uuidPattern = regexp.MustCompile(`^(?i:[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12})$`)

var environments = []string{"LOCAL", "DEVELOPMENT", "STAGING", "PRODUCTION"}

var errNoKind = fmt.Errorf("%w: the address is unset or has no kind", ErrInvalidAddress)

const (
	KindAgent      = "AGENT"
	KindExperiment = "EXPERIMENT"
	KindEvaluator  = "EVALUATOR"
)

// Address names the tracing project of a feature.
type Address interface {
	isAddress()
}

// AgentAddress is an agent environment. The zero value is unset. Build one with
// [NewAgentAddress]; a literal must set Kind to [KindAgent].
type AgentAddress struct {
	Kind        string `json:"kind,omitzero"`
	ID          string `json:"id,omitempty"`
	Environment string `json:"environment,omitempty"`
}

func NewAgentAddress(id, environment string) (AgentAddress, error) {
	return EnsureAgent(AgentAddress{Kind: KindAgent, ID: id, Environment: environment})
}

func (AgentAddress) isAddress() {}

func (a AgentAddress) IsZero() bool { return a == AgentAddress{} }

func (a *AgentAddress) UnmarshalJSON(data []byte) error {
	type plain AgentAddress
	var v plain
	if skip, err := decode(data, &v); skip || err != nil {
		return err
	}
	if v == (plain{}) { // the marshalled zero value
		*a = AgentAddress{}
		return nil
	}
	e, err := EnsureAgent(AgentAddress(v))
	if err == nil {
		*a = e
	}
	return err
}

// LRN returns the wire form, `lrn:agents/{id}/environments/{environment}`.
func LRN(a AgentAddress) string {
	return "lrn:agents/" + a.ID + "/environments/" + strings.ToLower(a.Environment)
}

// ExperimentAddress is an experiment's project. It cannot receive traces.
type ExperimentAddress struct {
	Kind string `json:"kind,omitzero"`
	ID   string `json:"id,omitempty"`
}

func NewExperimentAddress(id string) (ExperimentAddress, error) {
	a, err := Ensure(ExperimentAddress{Kind: KindExperiment, ID: id})
	if err != nil {
		return ExperimentAddress{}, err
	}
	return a.(ExperimentAddress), nil
}

func (ExperimentAddress) isAddress() {}

func (a *ExperimentAddress) UnmarshalJSON(data []byte) error {
	type plain ExperimentAddress
	var v plain
	if skip, err := decode(data, &v); skip || err != nil {
		return err
	}
	if v == (plain{}) { // the marshalled zero value
		*a = ExperimentAddress{}
		return nil
	}
	e, err := Ensure(ExperimentAddress(v))
	if err == nil {
		*a = e.(ExperimentAddress)
	}
	return err
}

// EvaluatorAddress is the workspace's evaluator project. It cannot receive traces.
type EvaluatorAddress struct {
	Kind string `json:"kind,omitzero"`
}

func NewEvaluatorAddress() EvaluatorAddress { return EvaluatorAddress{Kind: KindEvaluator} }

func (EvaluatorAddress) isAddress() {}

func (a *EvaluatorAddress) UnmarshalJSON(data []byte) error {
	type plain EvaluatorAddress
	var v plain
	if skip, err := decode(data, &v); skip || err != nil {
		return err
	}
	if v == (plain{}) { // the marshalled zero value
		*a = EvaluatorAddress{}
		return nil
	}
	e, err := Ensure(EvaluatorAddress(v))
	if err == nil {
		*a = e.(EvaluatorAddress)
	}
	return err
}

// decode reports skip for JSON null, which leaves the value unset.
func decode(data []byte, v any) (skip bool, err error) {
	if string(data) == "null" {
		return true, nil
	}
	if json.Unmarshal(data, v) != nil {
		return false, fmt.Errorf("%w: not an address object", ErrInvalidAddress)
	}
	return false, nil
}

func Ensure(a Address) (Address, error) {
	switch v := a.(type) {
	case AgentAddress:
		return EnsureAgent(v)
	case *AgentAddress:
		if v != nil {
			return EnsureAgent(*v)
		}
	case ExperimentAddress:
		if v.Kind == "" {
			return nil, errNoKind
		}
		if v.Kind == KindExperiment && uuidPattern.MatchString(v.ID) {
			return ExperimentAddress{Kind: KindExperiment, ID: strings.ToLower(v.ID)}, nil
		}
		return nil, fmt.Errorf("%w: an experiment address needs kind %s and an id that is a UUID with hyphens", ErrInvalidAddress, KindExperiment)
	case *ExperimentAddress:
		if v != nil {
			return Ensure(*v)
		}
	case EvaluatorAddress:
		if v.Kind == "" {
			return nil, errNoKind
		}
		if v.Kind == KindEvaluator {
			return v, nil
		}
		return nil, fmt.Errorf("%w: an evaluator address needs kind %s", ErrInvalidAddress, KindEvaluator)
	case *EvaluatorAddress:
		if v != nil {
			return Ensure(*v)
		}
	}
	return nil, fmt.Errorf("%w: not an agent, experiment or evaluator address", ErrInvalidAddress)
}

// EnsureAgent rejects any kind but AGENT: only an agent can receive traces.
func EnsureAgent(a AgentAddress) (AgentAddress, error) {
	if a.Kind == "" {
		return AgentAddress{}, errNoKind
	}
	if a.Kind != KindAgent {
		return AgentAddress{}, fmt.Errorf("%w: only an %s address can receive traces", ErrInvalidAddress, KindAgent)
	}
	if !agentIDPattern.MatchString(a.ID) {
		return AgentAddress{}, fmt.Errorf("%w: agent id must be 1 to 63 lowercase ASCII letters, digits, or hyphens, start with a letter, and end with a letter or digit", ErrInvalidAddress)
	}
	env := strings.ToUpper(a.Environment)
	for _, e := range environments {
		if env == e {
			return AgentAddress{Kind: KindAgent, ID: a.ID, Environment: env}, nil
		}
	}
	return AgentAddress{}, fmt.Errorf("%w: environment must be one of %s", ErrInvalidAddress, strings.ToLower(strings.Join(environments, ", ")))
}

func OptionalAgent(a AgentAddress) (AgentAddress, error) {
	if a.IsZero() {
		return a, nil
	}
	return EnsureAgent(a)
}

// FromEnv returns the zero address if neither env var is set, and an error
// wrapping [ErrEnvAddress] if only one is set or a value is invalid.
func FromEnv() (AgentAddress, error) {
	id, env := os.Getenv(EnvAgentID), os.Getenv(EnvAgentEnvironment)
	switch {
	case id == "" && env == "":
		return AgentAddress{}, nil
	case id == "":
		return AgentAddress{}, fmt.Errorf("%w: %s is set but %s is not", ErrEnvAddress, EnvAgentEnvironment, EnvAgentID)
	case env == "":
		return AgentAddress{}, fmt.Errorf("%w: %s is set but %s is not", ErrEnvAddress, EnvAgentID, EnvAgentEnvironment)
	}
	a, err := NewAgentAddress(id, env)
	if err != nil {
		return AgentAddress{}, fmt.Errorf("%w: %s / %s: %w", ErrEnvAddress, EnvAgentID, EnvAgentEnvironment, err)
	}
	return a, nil
}
