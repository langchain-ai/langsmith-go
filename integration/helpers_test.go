package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"sync"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/mockllm"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/sdk/trace/tracetest"
)

var (
	sharedMockServer     *mockllm.CombinedServer
	sharedMockServerOnce sync.Once
)

// getSharedMockServer returns a shared CombinedServer for all tests.
// Started once, lives for the test process lifetime.
func getSharedMockServer() *mockllm.CombinedServer {
	sharedMockServerOnce.Do(func() {
		sharedMockServer = mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.DefaultHandler()))
	})
	return sharedMockServer
}

// mockBaseURL returns the mock server URL if no real API key is set for the
// given provider. Returns ("", false) when real keys are available.
func mockBaseURL(provider string) (string, bool) {
	switch provider {
	case "anthropic":
		if os.Getenv("ANTHROPIC_API_KEY") != "" {
			return "", false
		}
	case "openai":
		if os.Getenv("OPENAI_API_KEY") != "" {
			return "", false
		}
	}
	return getSharedMockServer().URL(), true
}

var (
	sharedIntegrationProject     string
	sharedIntegrationProjectOnce sync.Once
)

// getSharedIntegrationProject returns one project name for the whole test run.
func getSharedIntegrationProject() string {
	sharedIntegrationProjectOnce.Do(func() {
		sharedIntegrationProject = uniqueName("__test_go_integ_tracing")
	})
	return sharedIntegrationProject
}

func requireEnv(t *testing.T, key string) string {
	t.Helper()
	val := os.Getenv(key)
	if val == "" {
		t.Skipf("%s not set, skipping integration test", key)
	}
	return val
}

func newClient(t *testing.T) *langsmith.Client {
	t.Helper()
	requireEnv(t, "LANGSMITH_API_KEY")
	return langsmith.NewClient()
}

func uniqueName(prefix string) string {
	return fmt.Sprintf("%s-%d-%d", prefix, time.Now().UnixNano(), rand.Intn(10000))
}

// tracedTP holds the resources for a test that sends traces to LangSmith.
type tracedTP struct {
	TP       *sdktrace.TracerProvider
	Exporter *tracetest.InMemoryExporter
	Tracer   *langsmith.Tracer
	Project  string
}

// newTracedTP creates a TracerProvider that sends traces to both an in-memory
// exporter (for local assertions) and the real LangSmith API.
// When prefix is empty, all tests share one project (getSharedIntegrationProject); trace tests use WithRunNameContext(ctx, name) with a hardcoded name per test type so runs are identifiable for polling.
// When prefix is non-empty, a unique project is used (prefix + timestamp + random).
// Shared project name uses go_integ_tracing__test prefix. If LANGSMITH_API_KEY is not set, it falls back to in-memory only.
func newTracedTP(t *testing.T, prefix string) *tracedTP {
	t.Helper()
	project := prefix
	if project == "" {
		project = getSharedIntegrationProject()
	} else {
		project = uniqueName(prefix)
	}
	exporter := tracetest.NewInMemoryExporter()
	tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))

	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		return &tracedTP{TP: tp, Exporter: exporter, Project: project}
	}

	tracerOpts := []langsmith.TracerOption{
		langsmith.WithAPIKey(apiKey),
		langsmith.WithProjectName(project),
	}
	if endpoint := os.Getenv("LANGSMITH_ENDPOINT"); endpoint != "" {
		host := strings.TrimPrefix(strings.TrimPrefix(endpoint, "https://"), "http://")
		tracerOpts = append(tracerOpts, langsmith.WithEndpoint(host))
	}
	ls, err := langsmith.New(tp, tracerOpts...)
	if err != nil {
		t.Fatalf("creating LangSmith tracer: %v", err)
	}

	return &tracedTP{TP: tp, Exporter: exporter, Tracer: ls, Project: project}
}

// Shutdown flushes and shuts down all resources.
func (tt *tracedTP) Shutdown(ctx context.Context) {
	tt.TP.ForceFlush(ctx)
	if tt.Tracer != nil {
		tt.Tracer.Shutdown(ctx)
	}
	tt.TP.Shutdown(ctx)
}

// Spans returns the locally captured spans.
func (tt *tracedTP) Spans() tracetest.SpanStubs {
	return tt.Exporter.GetSpans()
}

// SendsToLangSmith returns true if traces are being sent to LangSmith.
func (tt *tracedTP) SendsToLangSmith() bool {
	return tt.Tracer != nil
}

func getSpanAttr(spans tracetest.SpanStubs, key string) (string, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes {
			if string(attr.Key) == key {
				return attr.Value.Emit(), true
			}
		}
	}
	return "", false
}

func getSpanAttrInt(spans tracetest.SpanStubs, key string) (int64, bool) {
	for _, s := range spans {
		for _, attr := range s.Attributes {
			if string(attr.Key) == key {
				return attr.Value.AsInt64(), true
			}
		}
	}
	return 0, false
}

// pollForRuns queries LangSmith for runs in a session, polling until at least
// minRuns appear (or minRuns matching runName if runName is set). Session is looked up by project name.
// When runName is non-empty, returned runs are filtered to those with run.Name == runName,
// so tests that share one project can find their run by unique name.
//
// Polling uses exponential backoff (500ms → 1s → 2s → 2s …) with a total
// budget of 10s for the session lookup and 10s for the run query, keeping the
// worst-case wall time to ~20s instead of the previous 60s.
func pollForRuns(t *testing.T, projectName string, minRuns int, runName string) []langsmith.RunSchema {
	t.Helper()
	client := langsmith.NewClient()
	ctx := context.Background()

	const pollBudget = 10 * time.Second

	var sessionID string
	deadline := time.Now().Add(pollBudget)
	wait := 500 * time.Millisecond
	for time.Now().Before(deadline) {
		time.Sleep(wait)
		sessions, err := client.Sessions.List(ctx, langsmith.SessionListParams{
			Name: langsmith.F(projectName),
		})
		if err != nil {
			t.Logf("poll sessions: %v", err)
		} else if len(sessions.Items) > 0 {
			sessionID = sessions.Items[0].ID
			break
		}
		wait = min(wait*2, 2*time.Second)
	}
	if sessionID == "" {
		t.Errorf("session %q not found after polling", projectName)
		return nil
	}

	selectFields := []langsmith.RunQueryParamsSelect{
		langsmith.RunQueryParamsSelectID,
		langsmith.RunQueryParamsSelectName,
		langsmith.RunQueryParamsSelectRunType,
		langsmith.RunQueryParamsSelectStartTime,
		langsmith.RunQueryParamsSelectEndTime,
		langsmith.RunQueryParamsSelectParentRunID,
		langsmith.RunQueryParamsSelectError,
		langsmith.RunQueryParamsSelectStatus,
		langsmith.RunQueryParamsSelectInputs,
		langsmith.RunQueryParamsSelectOutputs,
		langsmith.RunQueryParamsSelectPromptTokens,
		langsmith.RunQueryParamsSelectCompletionTokens,
		langsmith.RunQueryParamsSelectTotalTokens,
		langsmith.RunQueryParamsSelectMessages,
	}
	deadline = time.Now().Add(pollBudget)
	wait = 500 * time.Millisecond
	for time.Now().Before(deadline) {
		time.Sleep(wait)
		result, err := client.Runs.Query(ctx, langsmith.RunQueryParams{
			Session: langsmith.F([]string{sessionID}),
			Limit:   langsmith.F(int64(100)),
			Select:  langsmith.F(selectFields),
		})
		if err != nil {
			t.Logf("poll runs: %v", err)
		} else {
			runs := result.Runs
			if runName != "" {
				filtered := runs[:0]
				for _, r := range runs {
					if r.Name == runName {
						filtered = append(filtered, r)
					}
				}
				runs = filtered
			}
			if len(runs) >= minRuns {
				return runs
			}
		}
		wait = min(wait*2, 2*time.Second)
	}
	if runName != "" {
		t.Errorf("expected at least %d run(s) with name %q in session %q, not found after polling", minRuns, runName, projectName)
	} else {
		t.Errorf("expected at least %d runs in session %q, not found after polling", minRuns, projectName)
	}
	return nil
}

type LangSmithRunAssertions struct {
	WantName    string
	WantRunType langsmith.RunTypeEnum
	ExpectError bool   // if true, run.Error should be non-empty (e.g. exception raised → SDK patches run with error)
	WantStatus  string // if non-empty, assert run.Status equals this (e.g. "error" when call fails)
	WantInputs  bool   // assert run has non-empty inputs when backend populates
	WantOutputs bool   // assert run has non-empty outputs when backend populates

	// Content assertions (run contents match expected behaviors)
	WantUsage               bool   // run has PromptTokens and CompletionTokens > 0
	WantInputsContainSystem bool   // run inputs reflect system message (e.g. contain "system")
	WantInputsMessageCount  int    // run inputs have at least this many messages (0 = don't check)
	WantOutputsContainTools bool   // run outputs contain tool_calls / tool use (e.g. "get_weather")
	WantOutputsContainText  string // run outputs contain this substring (e.g. expected response text)
}

// assertLangSmithRunFields asserts on id, name, run_type, start_time, end_time,
// parent_run_id, inputs, outputs, error, and tags where relevant (like Python
// integration tests). revision_id is not in the Go run schema.
func assertLangSmithRunFields(t *testing.T, r *langsmith.RunSchema, a LangSmithRunAssertions) {
	t.Helper()
	if r == nil {
		t.Fatal("run is nil")
	}
	if r.ID == "" {
		t.Error("run.id should be set")
	}
	if r.Name != a.WantName {
		t.Errorf("run.name = %q, want %q", r.Name, a.WantName)
	}
	if r.RunType != a.WantRunType {
		t.Errorf("run.run_type = %q, want %q", r.RunType, a.WantRunType)
	}
	if r.StartTime.IsZero() {
		t.Error("run.start_time should be set")
	}
	if r.EndTime.IsZero() {
		t.Error("run.end_time should be set")
	}
	if r.ParentRunID != "" {
		t.Errorf("run.parent_run_id = %q, want empty (root run)", r.ParentRunID)
	}
	if a.WantStatus != "" {
		if r.Status != a.WantStatus {
			t.Errorf("run.status = %q, want %q", r.Status, a.WantStatus)
		}
	}
	if a.ExpectError {
		if r.Error == "" {
			t.Error("run.error should be set for error case")
		}
	} else {
		if r.Error != "" {
			t.Errorf("run.error = %q, want empty for success", r.Error)
		}
	}
	// Inputs/outputs: when backend populates from OTLP trace, assert non-empty.
	if a.WantInputs {
		if r.Inputs == nil || len(r.Inputs) == 0 {
			t.Error("run.inputs should be non-empty (backend should populate from trace)")
		}
	}
	if a.WantOutputs {
		if r.Outputs == nil || len(r.Outputs) == 0 {
			t.Error("run.outputs should be non-empty (backend should populate from trace)")
		}
	}

	// Content assertions: run contents match expected behaviors
	if a.WantUsage {
		if r.PromptTokens <= 0 {
			t.Errorf("run.prompt_tokens = %d, want > 0 (usage in run)", r.PromptTokens)
		}
		if r.CompletionTokens <= 0 {
			t.Errorf("run.completion_tokens = %d, want > 0 (usage in run)", r.CompletionTokens)
		}
	}
	if a.WantInputsContainSystem {
		inputsStr := mapToJSONString(r.Inputs)
		if inputsStr == "" || !strings.Contains(inputsStr, "system") {
			t.Error("run.inputs should contain system message content (expected 'system' in inputs)")
		}
	}
	if a.WantInputsMessageCount > 0 {
		n := runInputsMessageCount(r)
		if n < a.WantInputsMessageCount {
			t.Errorf("run.inputs should have at least %d messages, got %d", a.WantInputsMessageCount, n)
		}
	}
	if a.WantOutputsContainTools {
		outputsStr := mapToJSONString(r.Outputs)
		if outputsStr == "" {
			t.Error("run.outputs should be non-empty for tool calling")
		} else if !strings.Contains(outputsStr, "tool") && !strings.Contains(outputsStr, "get_weather") {
			t.Error("run.outputs should contain tool_calls or tool use (e.g. get_weather)")
		}
	}
	if a.WantOutputsContainText != "" {
		outputsStr := mapToJSONString(r.Outputs)
		if outputsStr == "" || !strings.Contains(outputsStr, a.WantOutputsContainText) {
			t.Errorf("run.outputs should contain expected response text %q", a.WantOutputsContainText)
		}
	}

	// Tags: assert only if we ever set them from tracer (currently we don't)
	_ = r.Tags
}

// mapToJSONString returns a JSON string of m for substring assertions; returns "" if m is nil or marshal fails.
func mapToJSONString(m map[string]interface{}) string {
	if m == nil {
		return ""
	}
	b, err := json.Marshal(m)
	if err != nil {
		return ""
	}
	return string(b)
}

// runInputsMessageCount returns the number of messages in run inputs (Inputs["messages"] or run.Messages).
func runInputsMessageCount(r *langsmith.RunSchema) int {
	if r.Inputs != nil {
		if msgs, ok := r.Inputs["messages"].([]interface{}); ok {
			return len(msgs)
		}
	}
	if len(r.Messages) > 0 {
		return len(r.Messages)
	}
	return 0
}
