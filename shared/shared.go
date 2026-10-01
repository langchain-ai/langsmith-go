// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package shared

import (
	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
)

type AgentAddressParam struct {
	// `id` is the Agent's user-assigned id.
	ID param.Field[string] `json:"id" api:"required"`
	// `environment` is the Agent environment.
	Environment param.Field[AgentAddressEnvironment] `json:"environment" api:"required"`
	Kind        param.Field[AgentAddressKind]        `json:"kind" api:"required"`
}

func (r AgentAddressParam) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// `environment` is the Agent environment.
type AgentAddressEnvironment string

const (
	AgentAddressEnvironmentLocal       AgentAddressEnvironment = "LOCAL"
	AgentAddressEnvironmentDevelopment AgentAddressEnvironment = "DEVELOPMENT"
	AgentAddressEnvironmentStaging     AgentAddressEnvironment = "STAGING"
	AgentAddressEnvironmentProduction  AgentAddressEnvironment = "PRODUCTION"
)

func (r AgentAddressEnvironment) IsKnown() bool {
	switch r {
	case AgentAddressEnvironmentLocal, AgentAddressEnvironmentDevelopment, AgentAddressEnvironmentStaging, AgentAddressEnvironmentProduction:
		return true
	}
	return false
}

type AgentAddressKind string

const (
	AgentAddressKindAgent AgentAddressKind = "AGENT"
)

func (r AgentAddressKind) IsKnown() bool {
	switch r {
	case AgentAddressKindAgent:
		return true
	}
	return false
}
