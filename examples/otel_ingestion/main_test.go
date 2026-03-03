package main

import (
	"testing"
)

func TestLoadConfigMissingKey(t *testing.T) {
	t.Setenv("LANGSMITH_API_KEY", "")
	_, err := loadConfig()
	if err == nil {
		t.Error("expected error when LANGSMITH_API_KEY is missing")
	}
}

func TestLoadConfigDefaults(t *testing.T) {
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

func TestLoadConfigCustomProject(t *testing.T) {
	t.Setenv("LANGSMITH_API_KEY", "ls-test")
	t.Setenv("LANGSMITH_PROJECT", "my-project")

	cfg, err := loadConfig()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cfg.projectName != "my-project" {
		t.Errorf("expected 'my-project', got %q", cfg.projectName)
	}
}
