package main

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/sashabaranov/go-openai"
)

func TestCreateWeatherToolDefinition(t *testing.T) {
	def := createWeatherToolDefinition()
	if def.Name != "get_weather" {
		t.Errorf("expected name 'get_weather', got %q", def.Name)
	}
	if def.Description == "" {
		t.Error("expected non-empty description")
	}
	if def.Parameters == nil {
		t.Error("expected non-nil parameters")
	}
}

func TestExecuteToolWeather(t *testing.T) {
	result := executeTool("get_weather", `{"location":"Paris"}`)
	var parsed map[string]interface{}
	if err := json.Unmarshal([]byte(result), &parsed); err != nil {
		t.Fatalf("result is not valid JSON: %v", err)
	}
	if loc, ok := parsed["location"].(string); !ok || loc != "Paris" {
		t.Errorf("expected location 'Paris', got %v", parsed["location"])
	}
	if _, ok := parsed["temperature"]; !ok {
		t.Error("expected temperature in result")
	}
}

func TestExecuteToolUnknown(t *testing.T) {
	result := executeTool("nonexistent", `{}`)
	if !strings.Contains(result, "Unknown tool") {
		t.Errorf("expected 'Unknown tool' error, got: %s", result)
	}
}

func TestExecuteToolInvalidArgs(t *testing.T) {
	result := executeTool("get_weather", "not json")
	if !strings.Contains(result, "error") {
		t.Errorf("expected error for invalid JSON args, got: %s", result)
	}
}

func TestBuildPromptText(t *testing.T) {
	messages := []openai.ChatCompletionMessage{
		{Role: "user", Content: "Hello"},
		{Role: "assistant", Content: "Hi there"},
	}
	text := buildPromptText(messages)
	if !strings.Contains(text, "user: Hello") {
		t.Errorf("expected user message in prompt, got: %s", text)
	}
	if !strings.Contains(text, "assistant: Hi there") {
		t.Errorf("expected assistant message in prompt, got: %s", text)
	}
}

func TestBuildPromptTextWithToolCalls(t *testing.T) {
	messages := []openai.ChatCompletionMessage{
		{
			Role: "assistant",
			ToolCalls: []openai.ToolCall{
				{
					Type:     openai.ToolTypeFunction,
					Function: openai.FunctionCall{Name: "get_weather", Arguments: `{"location":"SF"}`},
				},
			},
		},
	}
	text := buildPromptText(messages)
	if !strings.Contains(text, "get_weather") {
		t.Errorf("expected tool call in prompt, got: %s", text)
	}
}

func TestLoadConfigMissingKeys(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "")
	t.Setenv("LANGSMITH_API_KEY", "")
	_, err := loadConfig()
	if err == nil {
		t.Error("expected error when OPENAI_API_KEY is missing")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
	t.Setenv("OPENAI_API_KEY", "sk-test")
	t.Setenv("LANGSMITH_API_KEY", "ls-test")
	t.Setenv("LANGSMITH_PROJECT", "")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.projectName != defaultProjectName {
		t.Errorf("expected default project name %q, got %q", defaultProjectName, cfg.projectName)
	}
}
