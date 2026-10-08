package env

import (
	"testing"
)

func TestTracingSampleRate(t *testing.T) {
	t.Run("unset", func(t *testing.T) {
		t.Setenv("LANGSMITH_TRACING_SAMPLING_RATE", "")
		t.Setenv("LANGCHAIN_TRACING_SAMPLING_RATE", "")
		r, err := TracingSampleRate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r != nil {
			t.Fatalf("expected nil rate when unset, got %v", *r)
		}
	})

	t.Run("valid", func(t *testing.T) {
		t.Setenv("LANGSMITH_TRACING_SAMPLING_RATE", "0.5")
		t.Setenv("LANGCHAIN_TRACING_SAMPLING_RATE", "")
		r, err := TracingSampleRate()
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if r == nil || *r != 0.5 {
			t.Fatalf("want 0.5, got %v", r)
		}
	})

	t.Run("invalid parse", func(t *testing.T) {
		t.Setenv("LANGSMITH_TRACING_SAMPLING_RATE", "not-a-number")
		t.Setenv("LANGCHAIN_TRACING_SAMPLING_RATE", "")
		_, err := TracingSampleRate()
		if err == nil {
			t.Fatal("expected error for invalid float")
		}
	})

	t.Run("out of range", func(t *testing.T) {
		t.Setenv("LANGSMITH_TRACING_SAMPLING_RATE", "2")
		t.Setenv("LANGCHAIN_TRACING_SAMPLING_RATE", "")
		_, err := TracingSampleRate()
		if err == nil {
			t.Fatal("expected error for rate > 1")
		}
	})
}

func TestRevisionID(t *testing.T) {
	tests := []struct {
		name      string
		langsmith string
		langchain string
		want      string
	}{
		{"unset", "", "", ""},
		{"langsmith", "new", "", "new"},
		{"legacy langchain", "", "legacy", "legacy"},
		{"langsmith wins", "new", "legacy", "new"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Setenv("LANGSMITH_REVISION_ID", tt.langsmith)
			t.Setenv("LANGCHAIN_REVISION_ID", tt.langchain)
			if got := RevisionID(); got != tt.want {
				t.Fatalf("RevisionID() = %q, want %q", got, tt.want)
			}
		})
	}
}
