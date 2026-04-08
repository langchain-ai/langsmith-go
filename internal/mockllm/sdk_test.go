package mockllm_test

import (
	"bytes"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"testing"

	"github.com/langchain-ai/langsmith-go/internal/mockllm"
)

func uvAvailable(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("uv"); err != nil {
		t.Skip("uv not found in PATH")
	}
}

func testdataDir(t *testing.T) string {
	t.Helper()
	_, thisFile, _, _ := runtime.Caller(0)
	return filepath.Join(filepath.Dir(thisFile), "testdata")
}

func runSDKTest(t *testing.T, script string) map[string]any {
	t.Helper()
	uvAvailable(t)

	srv := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.DefaultHandler()))
	defer srv.Close()

	scriptPath := filepath.Join(testdataDir(t), script)
	cmd := exec.Command("uv", "run", "--quiet", scriptPath)
	cmd.Env = append(os.Environ(),
		"ELIZA_BASE_URL="+srv.URL(),
		"OPENAI_API_KEY=fake",
		"ANTHROPIC_API_KEY=fake",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("uv run %s failed: %v\nstderr:\n%s\nstdout:\n%s", script, err, stderr.String(), stdout.String())
	}
	out := stdout.Bytes()

	var results map[string]any
	if err := json.Unmarshal(out, &results); err != nil {
		t.Fatalf("failed to parse JSON from %s: %v\nraw:\n%s", script, err, out)
	}
	return results
}

func TestSDK_OpenAI(t *testing.T) {
	r := runSDKTest(t, "test_openai_sdk.py")

	if v, _ := r["nonstreaming_content"].(string); v == "" {
		t.Error("expected non-empty nonstreaming content")
	}
	if v, _ := r["nonstreaming_model"].(string); v != "gpt-4o" {
		t.Errorf("model = %q, want gpt-4o", v)
	}
	if v, _ := r["nonstreaming_has_usage"].(bool); !v {
		t.Error("expected usage in response")
	}
	if v, _ := r["streaming_text"].(string); v == "" {
		t.Error("expected non-empty streaming text")
	}
	if v, _ := r["streaming_chunk_count"].(float64); v == 0 {
		t.Error("expected streaming chunks")
	}
	if v, _ := r["tool_finish_reason"].(string); v != "tool_calls" {
		t.Errorf("tool finish_reason = %q, want tool_calls", v)
	}
	if v, _ := r["tool_call_name"].(string); v == "" {
		t.Error("expected tool call name")
	}
	if v, _ := r["error_raised"].(bool); !v {
		t.Error("expected error for invalid API key")
	}
}

func TestSDK_Anthropic(t *testing.T) {
	r := runSDKTest(t, "test_anthropic_sdk.py")

	if v, _ := r["nonstreaming_content"].(string); v == "" {
		t.Error("expected non-empty nonstreaming content")
	}
	if v, _ := r["nonstreaming_model"].(string); v != "claude-sonnet-4-20250514" {
		t.Errorf("model = %q", v)
	}
	if v, _ := r["nonstreaming_has_usage"].(bool); !v {
		t.Error("expected usage")
	}
	if v, _ := r["nonstreaming_stop_reason"].(string); v != "end_turn" {
		t.Errorf("stop_reason = %q, want end_turn", v)
	}
	if v, _ := r["streaming_text"].(string); v == "" {
		t.Error("expected non-empty streaming text")
	}
	if v, _ := r["streaming_chunk_count"].(float64); v == 0 {
		t.Error("expected streaming chunks")
	}
	if v, _ := r["tool_stop_reason"].(string); v != "tool_use" {
		t.Errorf("tool stop_reason = %q, want tool_use", v)
	}
	if v, _ := r["tool_call_name"].(string); v == "" {
		t.Error("expected tool call name")
	}
	if v, _ := r["error_raised"].(bool); !v {
		t.Error("expected error for invalid API key")
	}
}

func TestSDK_LangChain(t *testing.T) {
	r := runSDKTest(t, "test_langchain.py")

	if v, _ := r["openai_invoke"].(string); v == "" {
		t.Error("expected non-empty openai invoke content")
	}
	if v, _ := r["openai_stream_chunks"].(float64); v == 0 {
		t.Error("expected openai stream chunks")
	}
	if v, _ := r["openai_stream_text"].(string); v == "" {
		t.Error("expected non-empty openai stream text")
	}
	if v, _ := r["anthropic_invoke"].(string); v == "" {
		t.Error("expected non-empty anthropic invoke content")
	}
	if v, _ := r["anthropic_stream_chunks"].(float64); v == 0 {
		t.Error("expected anthropic stream chunks")
	}
	if v, _ := r["anthropic_stream_text"].(string); v == "" {
		t.Error("expected non-empty anthropic stream text")
	}
}

func TestSDK_LangChain_Advanced(t *testing.T) {
	r := runSDKTest(t, "test_langchain_advanced.py")

	if v, _ := r["openai_multiturn_nonempty"].(bool); !v {
		t.Error("expected non-empty openai multi-turn response")
	}
	if v, _ := r["anthropic_system_nonempty"].(bool); !v {
		t.Error("expected non-empty anthropic system message response")
	}
	if v, _ := r["openai_tool_calls"].(float64); v == 0 {
		t.Error("expected openai tool calls")
	}
	if v, _ := r["openai_tool_name"].(string); v == "" {
		t.Error("expected openai tool call name")
	}
	if v, _ := r["anthropic_tool_calls"].(float64); v == 0 {
		t.Error("expected anthropic tool calls")
	}
	if v, _ := r["anthropic_tool_name"].(string); v == "" {
		t.Error("expected anthropic tool call name")
	}
	if v, _ := r["openai_system_stream_chunks"].(float64); v == 0 {
		t.Error("expected openai system stream chunks")
	}
	if v, _ := r["anthropic_multiturn_stream_chunks"].(float64); v == 0 {
		t.Error("expected anthropic multi-turn stream chunks")
	}
}

func npxAvailable(t *testing.T) {
	t.Helper()
	if _, err := exec.LookPath("npx"); err != nil {
		t.Skip("npx not found in PATH")
	}
}

func TestSDK_AISDK(t *testing.T) {
	npxAvailable(t)
	uvAvailable(t) // reuse server setup

	srv := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.DefaultHandler()))
	defer srv.Close()

	// Install deps and run
	testdata := testdataDir(t)
	installCmd := exec.Command("pnpm", "install", "--silent")
	installCmd.Dir = testdata
	if out, err := installCmd.CombinedOutput(); err != nil {
		t.Fatalf("pnpm install failed: %v\n%s", err, out)
	}

	cmd := exec.Command("node", filepath.Join(testdata, "test_ai_sdk.mjs"))
	cmd.Dir = testdata
	cmd.Env = append(os.Environ(),
		"ELIZA_BASE_URL="+srv.URL(),
		"OPENAI_API_KEY=fake",
	)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr
	if err := cmd.Run(); err != nil {
		t.Fatalf("node test_ai_sdk.mjs failed: %v\nstderr:\n%s\nstdout:\n%s", err, stderr.String(), stdout.String())
	}

	var r map[string]any
	if err := json.Unmarshal(stdout.Bytes(), &r); err != nil {
		t.Fatalf("failed to parse JSON: %v\nraw:\n%s", err, stdout.String())
	}

	if v, _ := r["generate_text"].(string); v == "" {
		t.Error("expected non-empty generateText result")
	}
	if v, _ := r["stream_has_text"].(bool); !v {
		t.Error("expected non-empty streamText result")
	}
	if v, _ := r["tool_calls_count"].(float64); v == 0 {
		t.Error("expected tool calls from AI SDK")
	}
}
