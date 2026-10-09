package langsmith_test

import (
	"net/url"
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

func TestNewSessionResolveParams(t *testing.T) {
	agent, err := langsmith.NewAgentAddress("support", "staging")
	if err != nil {
		t.Fatal(err)
	}
	exp, err := langsmith.NewExperimentAddress("0190c3d4-0000-7000-8000-0000000000b1")
	if err != nil {
		t.Fatal(err)
	}
	for _, c := range []struct {
		name string
		addr langsmith.Address
		want url.Values
	}{
		{"agent", agent, url.Values{"kind": {"AGENT"}, "id": {"support"}, "environment": {"STAGING"}}},
		{"experiment", exp, url.Values{"kind": {"EXPERIMENT"}, "id": {"0190c3d4-0000-7000-8000-0000000000b1"}}},
		{"evaluator", langsmith.NewEvaluatorAddress(), url.Values{"kind": {"EVALUATOR"}}},
	} {
		t.Run(c.name, func(t *testing.T) {
			got := langsmith.NewSessionResolveParams(c.addr).URLQuery()
			if got.Encode() != c.want.Encode() {
				t.Errorf("query = %v, want %v", got, c.want)
			}
		})
	}
}

func TestFeedbackWithAddress(t *testing.T) {
	agent, _ := langsmith.NewAgentAddress("support", "production")
	p, err := langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{Key: langsmith.F("k")}, agent)
	if err != nil {
		t.Fatal(err)
	}
	if !p.Address.Present || p.Address.Value != "lrn:agents/support/environments/production" {
		t.Errorf("address = %+v", p.Address)
	}
	if _, err := langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{SessionID: langsmith.F("0190c3d4-0000-7000-8000-0000000000b1")}, agent); err == nil {
		t.Error("project beside address: want error")
	}
	if _, err := langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{}, langsmith.AgentAddress{}); err == nil {
		t.Error("unset address: want error")
	}
}
