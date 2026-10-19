package langsmith_test

import (
	"errors"
	"net/url"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

func TestNewSessionResolveParams(t *testing.T) {
	agent, _ := langsmith.NewAgentAddress("support", "staging")
	exp, _ := langsmith.NewExperimentAddress("0190c3d4-0000-7000-8000-0000000000b1")
	for _, c := range []struct {
		name string
		addr langsmith.Address
		want url.Values
	}{
		{"agent", agent, url.Values{"kind": {"AGENT"}, "id": {"support"}, "environment": {"STAGING"}}},
		{"experiment", exp, url.Values{"kind": {"EXPERIMENT"}, "id": {"0190c3d4-0000-7000-8000-0000000000b1"}}},
		{"evaluator", langsmith.NewEvaluatorAddress(), url.Values{"kind": {"EVALUATOR"}}},
	} {
		p, err := langsmith.NewSessionResolveParams(c.addr)
		if err != nil || p.URLQuery().Encode() != c.want.Encode() {
			t.Errorf("%s: %v, %v", c.name, p.URLQuery(), err)
		}
	}
	if _, err := langsmith.NewSessionResolveParams(langsmith.AgentAddress{Kind: "AGENT", ID: "Bad_Id"}); !errors.Is(err, langsmith.ErrInvalidAddress) {
		t.Errorf("literal: got %v", err)
	}
}

func TestFeedbackWithAddress(t *testing.T) {
	agent, _ := langsmith.NewAgentAddress("support", "production")
	p, err := langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{Key: langsmith.F("k")}, agent)
	if err != nil || p.Address.Value != "lrn:agents/support/environments/production" {
		t.Errorf("address = %+v, %v", p.Address, err)
	}
	if _, err := langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{SessionID: langsmith.F("0190c3d4-0000-7000-8000-0000000000b1")}, agent); err == nil {
		t.Error("project beside address: want error")
	}
	if _, err := langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{}, langsmith.AgentAddress{Kind: "AGENT", ID: "Bad_Id"}); err == nil {
		t.Error("invalid address: want error")
	}
	_, err = langsmith.FeedbackWithAddress(langsmith.FeedbackCreateSchemaParam{}, langsmith.AgentAddress{})
	if !errors.Is(err, langsmith.ErrInvalidAddress) || !strings.Contains(err.Error(), "unset or has no kind") {
		t.Errorf("zero address: %v", err)
	}
	if _, err = langsmith.AddressField(langsmith.AgentAddress{}); !errors.Is(err, langsmith.ErrInvalidAddress) || !strings.Contains(err.Error(), "unset or has no kind") {
		t.Errorf("AddressField zero: %v", err)
	}
}
