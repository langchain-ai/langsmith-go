package langsmith

import (
	"errors"

	"github.com/langchain-ai/langsmith-go/internal/address"
	"github.com/langchain-ai/langsmith-go/internal/param"
)

// Address (beta) names the tracing project of a feature. Only an [AgentAddress] can receive traces.
type Address = address.Address

// APIAddress (beta) is an address as the query APIs take it.
type APIAddress = address.APIAddress

// AgentAddress (beta) is an agent environment. The zero value is unset.
type AgentAddress = address.AgentAddress

// ExperimentAddress (beta) is an experiment's project.
type ExperimentAddress = address.ExperimentAddress

// EvaluatorAddress (beta) is the workspace's evaluator project.
type EvaluatorAddress = address.EvaluatorAddress

// ErrInvalidAddress (beta) is wrapped by invalid-value errors, which never echo the input.
var ErrInvalidAddress = address.ErrInvalidAddress

// ErrEnvAddress (beta) is wrapped by the error for env vars that cannot address a run.
var ErrEnvAddress = address.ErrEnvAddress

// NewAgentAddress (beta) validates an agent id (a DNS label) and an environment
// (local, development, staging or production, in any case).
func NewAgentAddress(id, environment string) (AgentAddress, error) {
	return address.NewAgentAddress(id, environment)
}

// NewExperimentAddress (beta) validates a hyphenated experiment UUID.
func NewExperimentAddress(id string) (ExperimentAddress, error) {
	return address.NewExperimentAddress(id)
}

// NewEvaluatorAddress (beta) returns the evaluator address.
func NewEvaluatorAddress() EvaluatorAddress { return address.NewEvaluatorAddress() }

// AgentAddressFromEnv (beta) reads LANGSMITH_AGENT_ID and LANGSMITH_AGENT_ENVIRONMENT:
// both or neither, else an error wrapping [ErrEnvAddress].
func AgentAddressFromEnv() (AgentAddress, error) { return address.FromEnv() }

// NewSessionResolveParams (beta) returns the [SessionService.Resolve] params for an address.
func NewSessionResolveParams(a Address) SessionResolveParams {
	api := a.ToAPIAddress()
	p := SessionResolveParams{Kind: F(SessionResolveParamsKind(api.Kind))}
	if api.ID != "" {
		p.ID = F(api.ID)
	}
	if api.Environment != "" {
		p.Environment = F(SessionResolveParamsEnvironment(api.Environment))
	}
	return p
}

// AddressField (beta) returns the value for a generated `Address` field.
func AddressField(a AgentAddress) param.Field[string] { return F(address.LRN(a)) }

// FeedbackWithAddress (beta) addresses the feedback; it errors if SessionID is set.
func FeedbackWithAddress(params FeedbackCreateSchemaParam, a AgentAddress) (FeedbackCreateSchemaParam, error) {
	if a.IsZero() {
		return params, errors.New("langsmith: the feedback address is unset")
	}
	if params.SessionID.Present {
		return params, errors.New("langsmith: feedback is addressed by project or by address, not both")
	}
	params.Address = AddressField(a)
	return params, nil
}
