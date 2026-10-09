package langsmithtracing

import (
	"errors"
	"fmt"
	"time"

	"github.com/google/uuid"

	"github.com/langchain-ai/langsmith-go/internal/address"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/env"
	ilog "github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/logger"
)

// AgentAddress (beta) is an agent environment, the only address that can receive traces.
type AgentAddress = address.AgentAddress

// NewAgentAddress (beta) validates an agent id and environment.
func NewAgentAddress(id, environment string) (AgentAddress, error) {
	return address.NewAgentAddress(id, environment)
}

var errBothDestinations = errors.New("langsmith: a run is addressed by project or by address, not both")

// destination is a project or an address, never both.
type destination struct {
	project string
	address AgentAddress
}

// addressedTraceTTL bounds how long an addressed trace is remembered for its children.
const addressedTraceTTL = time.Hour

type addressedTrace struct {
	address AgentAddress
	at      time.Time
}

// resolveClientDestination returns the client default: options, then env. A bad
// env is envErr, not a failure, so runs naming their own destination still trace.
func resolveClientDestination(project string, addr AgentAddress, l ilog.Logger) (dest destination, envErr, err error) {
	switch {
	case project != "" && !addr.IsZero():
		return destination{}, nil, errBothDestinations
	case project != "":
		return destination{project: project}, nil, nil
	case !addr.IsZero():
		return destination{address: addr}, nil, nil
	}
	envProject := env.ProjectIfSet()
	envAddr, envErr := address.FromEnv()
	if envErr == nil && envProject != "" && !envAddr.IsZero() {
		envErr = fmt.Errorf("%w: a project and an address are both set in the environment", address.ErrEnvAddress)
	}
	if envErr != nil {
		// A fixed message: the error is built from the env values, and each
		// untraced run returns it to its caller.
		l.Warn("runs that name no destination are not traced: LANGSMITH_PROJECT and LANGSMITH_AGENT_ID / LANGSMITH_AGENT_ENVIRONMENT do not name a single valid destination")
		return destination{}, envErr, nil
	}
	if !envAddr.IsZero() {
		return destination{address: envAddr}, nil, nil
	}
	if envProject == "" {
		envProject = "default"
	}
	return destination{project: envProject}, nil, nil
}

// resolveDestination: the run's own project or address, then its addressed
// trace's, then the client's.
func (c *TracingClient) resolveDestination(id, traceID uuid.UUID, parent *uuid.UUID, sessionName string, sessionID *uuid.UUID, addr AgentAddress) (destination, error) {
	named := sessionName != "" || sessionID != nil
	switch {
	case named && !addr.IsZero():
		return destination{}, errBothDestinations
	case named:
		if sessionName == "" {
			sessionName = c.dest.project
		}
		return destination{project: sessionName}, nil
	case !addr.IsZero():
		c.rememberTrace(id, traceID, addr)
		return destination{address: addr}, nil
	}
	if parent != nil {
		if a, ok := c.addressedTrace(traceID); ok {
			return destination{address: a}, nil
		}
	}
	if c.envErr != nil {
		return destination{}, c.envErr
	}
	if !c.dest.address.IsZero() {
		c.rememberTrace(id, traceID, c.dest.address)
	}
	return c.dest, nil
}

func (c *TracingClient) applyDestination(runInfo map[string]any, d destination) {
	if d.address.IsZero() {
		if d.project != "" {
			runInfo["session_name"] = d.project
		}
		return
	}
	c.betaOnce.Do(func() {
		c.logger.Warn("addressing runs to an agent is in beta and enabled per workspace; a workspace without it rejects the run, so the trace is lost rather than falling back to a project")
	})
	runInfo["address"] = address.LRN(d.address)
}

// rememberTrace records a root run's address for its children.
func (c *TracingClient) rememberTrace(id, traceID uuid.UUID, a AgentAddress) {
	if id != traceID {
		return
	}
	c.addressedMu.Lock()
	defer c.addressedMu.Unlock()
	now := time.Now()
	if now.Sub(c.addressedLastPrune) >= filteredPruneInterval {
		c.addressedLastPrune = now
		for k, v := range c.addressedTraces {
			if now.Sub(v.at) > addressedTraceTTL {
				delete(c.addressedTraces, k)
			}
		}
	}
	c.addressedTraces[traceID] = addressedTrace{address: a, at: now}
}

func (c *TracingClient) addressedTrace(traceID uuid.UUID) (AgentAddress, bool) {
	c.addressedMu.Lock()
	defer c.addressedMu.Unlock()
	t, ok := c.addressedTraces[traceID]
	return t.address, ok
}
