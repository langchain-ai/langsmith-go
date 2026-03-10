package main

import (
	"testing"
)

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
