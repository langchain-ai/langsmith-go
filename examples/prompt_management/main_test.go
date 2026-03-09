package main

import (
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

func TestChatPromptBuilderBuild(t *testing.T) {
	builder := &ChatPromptBuilder{}
	builder.
		SystemMessage("You are a helpful assistant.").
		UserMessage("Tell me about {topic}.")

	manifest, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	kwargs, ok := manifest["kwargs"].(map[string]interface{})
	if !ok {
		t.Fatal("expected kwargs in manifest")
	}
	messages, ok := kwargs["messages"].([]interface{})
	if !ok {
		t.Fatal("expected messages in kwargs")
	}
	if len(messages) != 2 {
		t.Errorf("expected 2 messages, got %d", len(messages))
	}

	inputVars, ok := kwargs["input_variables"].([]string)
	if !ok {
		t.Fatal("expected input_variables in kwargs")
	}
	if len(inputVars) != 1 || inputVars[0] != "topic" {
		t.Errorf("expected input_variables [topic], got %v", inputVars)
	}
}

func TestChatPromptBuilderEmpty(t *testing.T) {
	builder := &ChatPromptBuilder{}
	_, err := builder.Build()
	if err == nil {
		t.Error("expected error for empty builder")
	}
}

func TestChatPromptBuilderExplicitVars(t *testing.T) {
	builder := &ChatPromptBuilder{}
	builder.
		UserMessage("Hello {name}").
		InputVariables("name", "extra")

	manifest, err := builder.Build()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	kwargs := manifest["kwargs"].(map[string]interface{})
	vars := kwargs["input_variables"].([]string)
	if len(vars) != 2 {
		t.Errorf("expected 2 explicit vars, got %d", len(vars))
	}
}

func TestExtractVariablesFromTemplate(t *testing.T) {
	tests := []struct {
		template string
		expected []string
	}{
		{"Hello {name}, welcome to {place}", []string{"name", "place"}},
		{"No variables here", nil},
		{"{x} and {x} again", []string{"x"}},
		{"{a}{b}{c}", []string{"a", "b", "c"}},
	}
	for _, tt := range tests {
		got := extractVariablesFromTemplate(tt.template)
		if len(got) == 0 && len(tt.expected) == 0 {
			continue
		}
		if len(got) != len(tt.expected) {
			t.Errorf("template %q: expected %v, got %v", tt.template, tt.expected, got)
			continue
		}
		for i := range got {
			if got[i] != tt.expected[i] {
				t.Errorf("template %q: var %d: expected %q, got %q", tt.template, i, tt.expected[i], got[i])
			}
		}
	}
}

func TestGetMessageType(t *testing.T) {
	tests := []struct {
		name     string
		message  map[string]interface{}
		expected string
	}{
		{
			"system message",
			map[string]interface{}{"id": []interface{}{"langchain_core", "prompts", "chat", "SystemMessagePromptTemplate"}},
			"system",
		},
		{
			"AI message",
			map[string]interface{}{"id": []interface{}{"langchain_core", "prompts", "chat", "AIMessagePromptTemplate"}},
			"assistant",
		},
		{
			"human message",
			map[string]interface{}{"id": []interface{}{"langchain_core", "prompts", "chat", "HumanMessagePromptTemplate"}},
			"user",
		},
		{
			"no id",
			map[string]interface{}{},
			"user",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := getMessageType(tt.message); got != tt.expected {
				t.Errorf("expected %q, got %q", tt.expected, got)
			}
		})
	}
}

func TestExtractOwner(t *testing.T) {
	repo := &langsmith.RepoWithLookups{}
	repo.Owner = "myorg"
	repo.FullName = "myorg/my-prompt"
	if got := extractOwner(repo); got != "myorg" {
		t.Errorf("expected 'myorg', got %q", got)
	}

	repo2 := &langsmith.RepoWithLookups{}
	repo2.FullName = "otherorg/prompt"
	if got := extractOwner(repo2); got != "otherorg" {
		t.Errorf("expected 'otherorg', got %q", got)
	}

	repo3 := &langsmith.RepoWithLookups{}
	if got := extractOwner(repo3); got != defaultOwner {
		t.Errorf("expected %q, got %q", defaultOwner, got)
	}
}

func TestExtractPromptContent(t *testing.T) {
	manifest := map[string]interface{}{
		"kwargs": map[string]interface{}{
			"messages": []interface{}{
				map[string]interface{}{
					"id": []interface{}{"langchain_core", "prompts", "chat", "SystemMessagePromptTemplate"},
					"kwargs": map[string]interface{}{
						"prompt": map[string]interface{}{
							"kwargs": map[string]interface{}{
								"template": "You are helpful.",
							},
						},
					},
				},
			},
		},
	}
	content := extractPromptContent(manifest)
	if content == "" || content == "Unable to parse prompt content" {
		t.Errorf("expected prompt content, got: %q", content)
	}
}

func TestExtractPromptContentInvalid(t *testing.T) {
	content := extractPromptContent(map[string]interface{}{})
	if content != "Unable to parse prompt content" {
		t.Errorf("expected fallback message, got: %q", content)
	}
}
