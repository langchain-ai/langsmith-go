package main

import (
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/sashabaranov/go-openai"
)

func TestGetTestCases(t *testing.T) {
	cases := getTestCases()
	if len(cases) == 0 {
		t.Fatal("expected at least one test case")
	}
	for i, tc := range cases {
		if tc.question == "" {
			t.Errorf("test case %d has empty question", i)
		}
		if tc.expectedAnswer == "" {
			t.Errorf("test case %d has empty expected answer", i)
		}
	}
}

func TestGenerateExampleID_Deterministic(t *testing.T) {
	dsID := "550e8400-e29b-41d4-a716-446655440000"
	input := map[string]interface{}{"question": "test"}
	output := map[string]interface{}{"answer": "test"}

	id1, err := generateExampleID(dsID, input, output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	id2, err := generateExampleID(dsID, input, output)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if id1 != id2 {
		t.Errorf("expected deterministic IDs, got %s != %s", id1, id2)
	}
}

func TestGenerateExampleID_InvalidDatasetID(t *testing.T) {
	input := map[string]interface{}{"q": "a"}
	output := map[string]interface{}{"a": "1"}

	id, err := generateExampleID("not-a-uuid", input, output)
	if err != nil {
		t.Fatalf("expected fallback namespace, got error: %v", err)
	}
	if id == "" {
		t.Error("expected non-empty ID with fallback namespace")
	}
}

func TestExtractAnswer(t *testing.T) {
	resp := openai.ChatCompletionResponse{
		Choices: []openai.ChatCompletionChoice{
			{Message: openai.ChatCompletionMessage{Content: "Paris"}},
		},
	}
	if got := extractAnswer(resp); got != "Paris" {
		t.Errorf("expected 'Paris', got %q", got)
	}
}

func TestExtractAnswerEmpty(t *testing.T) {
	resp := openai.ChatCompletionResponse{}
	if got := extractAnswer(resp); got != "No response" {
		t.Errorf("expected 'No response', got %q", got)
	}
}

func TestBuildExperimentSpanAttributes(t *testing.T) {
	attrs := buildExperimentSpanAttributes("session-1", "example-1", "What is Go?")
	if len(attrs) == 0 {
		t.Fatal("expected attributes")
	}
}

func TestGetProjectName(t *testing.T) {
	t.Setenv("LANGSMITH_PROJECT", "")
	if name := getProjectName(); name != defaultProjectName {
		t.Errorf("expected %q, got %q", defaultProjectName, name)
	}

	t.Setenv("LANGSMITH_PROJECT", "custom")
	if name := getProjectName(); name != "custom" {
		t.Errorf("expected 'custom', got %q", name)
	}
}

func TestBuildDatasetURL(t *testing.T) {
	t.Setenv("LANGSMITH_ENDPOINT", "https://api.smith.langchain.com")
	ds := &langsmith.Dataset{}
	ds.TenantID = "tenant-1"
	ds.ID = "ds-1"
	url := buildDatasetURL(ds)
	if !strings.Contains(url, "tenant-1") || !strings.Contains(url, "ds-1") {
		t.Errorf("unexpected URL: %s", url)
	}
}
