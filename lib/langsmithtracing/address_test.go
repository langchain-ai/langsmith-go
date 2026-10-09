package langsmithtracing_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/google/uuid"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing"
)

type capture struct {
	mu  sync.Mutex
	ops []langsmithtracing.RunOp
}

func (c *capture) hook(ops []langsmithtracing.RunOp) []langsmithtracing.RunOp {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.ops = append(c.ops, ops...)
	return ops
}

func (c *capture) byName(name string) map[string]any {
	c.mu.Lock()
	defer c.mu.Unlock()
	for _, op := range c.ops {
		if op.Data["name"] == name {
			return op.Data
		}
	}
	return nil
}

func newAddressClient(t *testing.T, opts ...langsmithtracing.Option) (*langsmithtracing.TracingClient, *capture) {
	t.Helper()
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusAccepted)
	}))
	t.Cleanup(srv.Close)
	rec := &capture{}
	opts = append([]langsmithtracing.Option{
		langsmithtracing.WithAPIURL(srv.URL),
		langsmithtracing.WithAPIKey("test-key"),
		langsmithtracing.WithRunTransform(rec.hook),
	}, opts...)
	c, err := langsmithtracing.NewTracingClient(context.Background(), opts...)
	if err != nil {
		t.Fatal(err)
	}
	return c, rec
}

func clearAddressEnv(t *testing.T) {
	t.Helper()
	for _, k := range []string{"LANGSMITH_AGENT_ID", "LANGSMITH_AGENT_ENVIRONMENT", "LANGSMITH_PROJECT", "LANGCHAIN_PROJECT"} {
		t.Setenv(k, "")
	}
}

func mustAgent(t *testing.T, id, env string) langsmithtracing.AgentAddress {
	t.Helper()
	a, err := langsmithtracing.NewAgentAddress(id, env)
	if err != nil {
		t.Fatal(err)
	}
	return a
}

func rootRun(name string, mod func(*langsmithtracing.RunCreate)) *langsmithtracing.RunCreate {
	id := uuid.New()
	r := &langsmithtracing.RunCreate{ID: id, TraceID: id, Name: name, RunType: "chain", StartTime: time.Now(), DottedOrder: "x"}
	if mod != nil {
		mod(r)
	}
	return r
}

func TestRunAddress(t *testing.T) {
	clearAddressEnv(t)
	c, rec := newAddressClient(t)
	if err := c.CreateRun(rootRun("addressed", func(r *langsmithtracing.RunCreate) {
		r.Address = mustAgent(t, "support", "Production")
	})); err != nil {
		t.Fatal(err)
	}
	c.Close()
	d := rec.byName("addressed")
	if d["address"] != "lrn:agents/support/environments/production" {
		t.Errorf("address = %v", d["address"])
	}
	if _, ok := d["session_name"]; ok {
		t.Errorf("an addressed run must not carry session_name: %v", d["session_name"])
	}
}

func TestRunWithProjectAndAddressIsRejected(t *testing.T) {
	clearAddressEnv(t)
	c, _ := newAddressClient(t)
	defer c.Close()
	err := c.CreateRun(rootRun("both", func(r *langsmithtracing.RunCreate) {
		r.SessionName = "my-secret-project"
		r.Address = mustAgent(t, "support", "local")
	}))
	if err == nil {
		t.Fatal("want error")
	}
	if strings.Contains(err.Error(), "my-secret-project") || strings.Contains(err.Error(), "support") {
		t.Errorf("error echoes input: %v", err)
	}
}

func TestChildFollowsAddressedRoot(t *testing.T) {
	clearAddressEnv(t)
	c, rec := newAddressClient(t, langsmithtracing.WithProject("default-project"))
	r := rootRun("root", func(r *langsmithtracing.RunCreate) { r.Address = mustAgent(t, "support", "production") })
	_ = c.CreateRun(r)
	if err := c.CreateRun(&langsmithtracing.RunCreate{
		ID: uuid.New(), TraceID: r.TraceID, ParentRunID: &r.ID, Name: "child", RunType: "llm",
		StartTime: time.Now(), DottedOrder: "x.y",
	}); err != nil {
		t.Fatal(err)
	}
	c.Close()
	if got := rec.byName("child"); got["address"] != "lrn:agents/support/environments/production" || got["session_name"] != nil {
		t.Errorf("child = %v", got)
	}
}
