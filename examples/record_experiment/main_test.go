package main

import (
	"strings"
	"testing"
	"time"
)

func TestDefineExperimentResults(t *testing.T) {
	results := defineExperimentResults()
	if len(results) == 0 {
		t.Fatal("expected at least one experiment result")
	}
	for i, r := range results {
		if r.Input == nil {
			t.Errorf("result %d has nil input", i)
		}
		if r.ReferenceOutput == nil {
			t.Errorf("result %d has nil reference output", i)
		}
		if r.ActualOutput == nil {
			t.Errorf("result %d has nil actual output", i)
		}
		if r.StartTime.IsZero() {
			t.Errorf("result %d has zero start time", i)
		}
		if !r.EndTime.After(r.StartTime) {
			t.Errorf("result %d end time should be after start time", i)
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

func TestGenerateExampleID_DifferentInputsDifferentIDs(t *testing.T) {
	dsID := "550e8400-e29b-41d4-a716-446655440000"

	id1, _ := generateExampleID(dsID, map[string]interface{}{"q": "a"}, map[string]interface{}{"a": "1"})
	id2, _ := generateExampleID(dsID, map[string]interface{}{"q": "b"}, map[string]interface{}{"a": "2"})
	if id1 == id2 {
		t.Error("different inputs should produce different IDs")
	}
}

func TestFormatDottedOrder(t *testing.T) {
	ts := time.Date(2026, 1, 15, 10, 30, 45, 123456000, time.UTC)
	runID := "abc-123"
	result := formatDottedOrder(ts, runID)

	if !strings.HasPrefix(result, "20260115T103045") {
		t.Errorf("expected timestamp prefix, got: %s", result)
	}
	if !strings.Contains(result, "123456Z") {
		t.Errorf("expected microseconds, got: %s", result)
	}
	if !strings.HasSuffix(result, runID) {
		t.Errorf("expected run ID suffix, got: %s", result)
	}
}

func TestBuildDatasetURL(t *testing.T) {
	url := buildDatasetURL("tenant-1", "ds-1")
	if url != "https://smith.langchain.com/o/tenant-1/datasets/ds-1" {
		t.Errorf("unexpected URL: %s", url)
	}
}

func TestBuildSessionURL(t *testing.T) {
	url := buildSessionURL("tenant-1", "session-1")
	if url != "https://smith.langchain.com/o/tenant-1/projects/p/session-1" {
		t.Errorf("unexpected URL: %s", url)
	}
}

func TestGenerateUUID(t *testing.T) {
	id, err := generateUUID()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(id) != 36 {
		t.Errorf("expected UUID length 36, got %d", len(id))
	}
}
