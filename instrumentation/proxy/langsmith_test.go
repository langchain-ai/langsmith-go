//go:build integration

package proxy_test

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"os/exec"
	"strings"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go/instrumentation/proxy"
	"github.com/langchain-ai/langsmith-go/internal/mockllm"
)

// TestLangSmith_TracesLand verifies that traces sent through the proxy
// actually appear in LangSmith by querying them back with the langsmith CLI.
//
// Requires:
//   - LANGSMITH_API_KEY env var
//   - langsmith CLI in PATH
func TestLangSmith_TracesLand(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set")
	}
	if _, err := exec.LookPath("langsmith"); err != nil {
		t.Skip("langsmith CLI not found in PATH")
	}

	project := fmt.Sprintf("__test_proxy_integ_%d_%d", time.Now().UnixNano(), rand.Intn(10000))
	t.Logf("project: %s", project)

	// Start Eliza as upstream
	upstream := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.DefaultHandler()))
	defer upstream.Close()

	// Start proxy with real LangSmith tracing
	p, err := proxy.New(proxy.Config{
		AnthropicUpstream: upstream.URL(),
		OpenAIUpstream:    upstream.URL(),
		LangSmithAPIKey:   apiKey,
		LangSmithProject:  project,
	})
	if err != nil {
		t.Fatalf("proxy.New: %v", err)
	}
	defer p.Shutdown(context.Background())

	proxyServer := httptest.NewServer(p)
	defer proxyServer.Close()

	// --- Send requests through proxy ---

	t.Log("sending OpenAI non-streaming request...")
	sendOpenAI(t, proxyServer.URL, `{"model":"gpt-4o","messages":[{"role":"user","content":"hello from proxy test"}]}`)

	t.Log("sending Anthropic non-streaming request...")
	sendAnthropic(t, proxyServer.URL, `{"model":"claude-sonnet-4-20250514","messages":[{"role":"user","content":"I am a proxy test"}],"max_tokens":100}`)

	t.Log("sending OpenAI streaming request...")
	sendOpenAI(t, proxyServer.URL, `{"model":"gpt-4o","messages":[{"role":"user","content":"streaming test"}],"stream":true}`)

	// Flush traces
	p.Shutdown(context.Background())

	// --- Poll LangSmith CLI for runs ---

	t.Log("polling LangSmith for traces...")
	runs := pollLangSmithCLI(t, project, 3, apiKey)

	if len(runs) < 3 {
		t.Fatalf("expected at least 3 runs, got %d", len(runs))
	}

	// Verify run shapes
	var foundOpenAI, foundAnthropic bool
	for _, run := range runs {
		name, _ := run["name"].(string)
		runType, _ := run["run_type"].(string)

		if runType != "llm" {
			t.Errorf("run %q: run_type = %q, want llm", name, runType)
		}

		if strings.Contains(name, "openai") {
			foundOpenAI = true
		}
		if strings.Contains(name, "anthropic") {
			foundAnthropic = true
		}
	}

	if !foundOpenAI {
		t.Error("expected at least one OpenAI run")
	}
	if !foundAnthropic {
		t.Error("expected at least one Anthropic run")
	}

	t.Logf("found %d runs in project %s", len(runs), project)
}

// TestLangSmith_RunContents verifies that run inputs, outputs, and token
// usage are populated correctly in LangSmith.
func TestLangSmith_RunContents(t *testing.T) {
	apiKey := os.Getenv("LANGSMITH_API_KEY")
	if apiKey == "" {
		t.Skip("LANGSMITH_API_KEY not set")
	}
	if _, err := exec.LookPath("langsmith"); err != nil {
		t.Skip("langsmith CLI not found in PATH")
	}

	project := fmt.Sprintf("__test_proxy_contents_%d_%d", time.Now().UnixNano(), rand.Intn(10000))

	upstream := mockllm.NewCombinedServer(mockllm.WithHandler(mockllm.DefaultHandler()))
	defer upstream.Close()

	p, err := proxy.New(proxy.Config{
		AnthropicUpstream: upstream.URL(),
		OpenAIUpstream:    upstream.URL(),
		LangSmithAPIKey:   apiKey,
		LangSmithProject:  project,
	})
	if err != nil {
		t.Fatalf("proxy.New: %v", err)
	}
	defer p.Shutdown(context.Background())

	proxyServer := httptest.NewServer(p)
	defer proxyServer.Close()

	// Send a simple request
	sendOpenAI(t, proxyServer.URL, `{"model":"gpt-4o","messages":[{"role":"user","content":"What is Go?"}]}`)

	p.Shutdown(context.Background())

	runs := pollLangSmithCLI(t, project, 1, apiKey)
	if len(runs) == 0 {
		t.Fatal("no runs found")
	}

	run := runs[0]

	// Check inputs
	if inputs, ok := run["inputs"].(map[string]any); ok {
		inputsJSON, _ := json.Marshal(inputs)
		if !strings.Contains(string(inputsJSON), "What is Go?") {
			t.Errorf("run inputs should contain prompt text, got: %s", inputsJSON)
		}
	} else {
		t.Error("run should have inputs")
	}

	// Check outputs
	if outputs, ok := run["outputs"].(map[string]any); ok {
		outputsJSON, _ := json.Marshal(outputs)
		if len(outputsJSON) == 0 || string(outputsJSON) == "{}" {
			t.Error("run outputs should be non-empty")
		}
	} else {
		t.Error("run should have outputs")
	}

	// Check token usage (nested under token_usage)
	if tokenUsage, ok := run["token_usage"].(map[string]any); ok {
		if pt, _ := tokenUsage["prompt_tokens"].(float64); pt == 0 {
			t.Errorf("token_usage.prompt_tokens = %v, want > 0", tokenUsage["prompt_tokens"])
		}
		if ct, _ := tokenUsage["completion_tokens"].(float64); ct == 0 {
			t.Errorf("token_usage.completion_tokens = %v, want > 0", tokenUsage["completion_tokens"])
		}
	} else {
		t.Errorf("expected token_usage object, got %v", run["token_usage"])
	}
}

// --- Helpers ---

func sendOpenAI(t *testing.T, proxyURL, body string) {
	t.Helper()
	resp, err := http.Post(proxyURL+"/v1/chat/completions", "application/json", strings.NewReader(body))
	if err != nil {
		t.Fatalf("openai request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("openai status: %d", resp.StatusCode)
	}
}

func sendAnthropic(t *testing.T, proxyURL, body string) {
	t.Helper()
	req, _ := http.NewRequest("POST", proxyURL+"/v1/messages", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("x-api-key", "fake")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		t.Fatalf("anthropic request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()
	if resp.StatusCode != 200 {
		t.Fatalf("anthropic status: %d", resp.StatusCode)
	}
}

// pollLangSmithCLI uses the langsmith CLI to query runs, polling until
// at least minRuns appear or the timeout is reached.
func pollLangSmithCLI(t *testing.T, project string, minRuns int, apiKey string) []map[string]any {
	t.Helper()

	deadline := time.Now().Add(30 * time.Second)
	wait := 1 * time.Second

	for time.Now().Before(deadline) {
		time.Sleep(wait)

		cmd := exec.Command("langsmith", "run", "list",
			"--project", project,
			"--limit", "20",
			"--full",
			"--format", "json",
		)
		cmd.Env = append(os.Environ(), "LANGSMITH_API_KEY="+apiKey)

		var stderr strings.Builder
		cmd.Stderr = &stderr
		out, err := cmd.Output()
		if err != nil {
			t.Logf("langsmith run list: %v (stderr: %s)", err, stderr.String())
			wait = min(wait*2, 4*time.Second)
			continue
		}

		var runs []map[string]any
		if err := json.Unmarshal(out, &runs); err != nil {
			t.Logf("parse runs: %v (raw: %s)", err, string(out[:min(len(out), 200)]))
			wait = min(wait*2, 4*time.Second)
			continue
		}

		if len(runs) >= minRuns {
			return runs
		}

		t.Logf("found %d/%d runs, polling...", len(runs), minRuns)
		wait = min(wait*2, 4*time.Second)
	}

	t.Errorf("expected at least %d runs in project %q, timed out", minRuns, project)
	return nil
}
