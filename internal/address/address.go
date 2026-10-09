// Package address implements the beta tracing addresses, re-exported by langsmith.
package address

import (
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

// APIAddress is an address as the query APIs take it. Environment is set only for AGENT.
type APIAddress struct {
	Kind        string
	ID          string
	Environment string
}

// Address names the tracing project of a feature.
type Address interface {
	ToAPIAddress() APIAddress
	isAddress()
}

// AgentAddress is an agent environment. The zero value is unset.
type AgentAddress struct {
	id          string
	environment string // upper case
}

func NewAgentAddress(id, environment string) (AgentAddress, error) {
	if !agentIDPattern.MatchString(id) {
		return AgentAddress{}, fmt.Errorf("%w: agent id must be 1 to 63 lowercase ASCII letters, digits, or hyphens, start with a letter, and end with a letter or digit", ErrInvalidAddress)
	}
	env := strings.ToUpper(environment)
	for _, e := range environments {
		if env == e {
			return AgentAddress{id: id, environment: env}, nil
		}
	}
	return AgentAddress{}, fmt.Errorf("%w: environment must be one of %s", ErrInvalidAddress, strings.ToLower(strings.Join(environments, ", ")))
}

func (AgentAddress) isAddress() {}

func (a AgentAddress) IsZero() bool { return a == AgentAddress{} }

func (a AgentAddress) ID() string { return a.id }

func (a AgentAddress) Environment() string { return a.environment }

func (a AgentAddress) ToAPIAddress() APIAddress {
	return APIAddress{Kind: "AGENT", ID: a.id, Environment: a.environment}
}

// LRN returns the wire form, `lrn:agents/{id}/environments/{environment}`.
func LRN(a AgentAddress) string {
	return "lrn:agents/" + a.id + "/environments/" + strings.ToLower(a.environment)
}

// ExperimentAddress is an experiment's project. It cannot receive traces.
type ExperimentAddress struct {
	id string
}

func NewExperimentAddress(id string) (ExperimentAddress, error) {
	if !uuidPattern.MatchString(id) {
		return ExperimentAddress{}, fmt.Errorf("%w: experiment id must be a UUID with hyphens", ErrInvalidAddress)
	}
	return ExperimentAddress{id: strings.ToLower(id)}, nil
}

func (ExperimentAddress) isAddress() {}

func (a ExperimentAddress) ID() string { return a.id }

func (a ExperimentAddress) ToAPIAddress() APIAddress {
	return APIAddress{Kind: "EXPERIMENT", ID: a.id}
}

// EvaluatorAddress is the workspace's evaluator project. It cannot receive traces.
type EvaluatorAddress struct{}

func NewEvaluatorAddress() EvaluatorAddress { return EvaluatorAddress{} }

func (EvaluatorAddress) isAddress() {}

func (EvaluatorAddress) ToAPIAddress() APIAddress { return APIAddress{Kind: "EVALUATOR"} }

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
