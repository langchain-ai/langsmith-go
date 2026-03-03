package main

import (
	"os"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

func TestGetIndividualExamples(t *testing.T) {
	examples := getIndividualExamples()
	if len(examples) == 0 {
		t.Fatal("expected at least one individual example")
	}
	for i, ex := range examples {
		if ex.question == "" {
			t.Errorf("example %d has empty question", i)
		}
		if ex.answer == "" {
			t.Errorf("example %d has empty answer", i)
		}
	}
}

func TestGetBulkExamples(t *testing.T) {
	examples := getBulkExamples()
	if len(examples) == 0 {
		t.Fatal("expected at least one bulk example")
	}
}

func TestBuildExampleParams(t *testing.T) {
	params := buildExampleParams("test-dataset-id", exampleData{
		question: "What is Go?",
		answer:   "A programming language.",
	})
	if params.DatasetID.Value == "" {
		t.Error("expected non-empty DatasetID")
	}
}

func TestBuildBulkExampleBody(t *testing.T) {
	body := buildBulkExampleBody("test-dataset-id", exampleData{
		question: "What is Go?",
		answer:   "A programming language.",
	})
	if body.DatasetID.Value == "" {
		t.Error("expected non-empty DatasetID")
	}
}

func TestBuildDatasetURL(t *testing.T) {
	os.Setenv("LANGSMITH_ENDPOINT", "https://api.smith.langchain.com")
	defer os.Unsetenv("LANGSMITH_ENDPOINT")

	ds := &langsmith.Dataset{}
	ds.TenantID = "tenant-1"
	ds.ID = "ds-1"
	url := buildDatasetURL(ds)
	if !strings.Contains(url, "tenant-1") || !strings.Contains(url, "ds-1") {
		t.Errorf("unexpected URL: %s", url)
	}
}

func TestBuildDatasetURLDefaultEndpoint(t *testing.T) {
	os.Unsetenv("LANGSMITH_ENDPOINT")

	ds := &langsmith.Dataset{}
	ds.TenantID = "t"
	ds.ID = "d"
	url := buildDatasetURL(ds)
	if !strings.HasPrefix(url, "https://smith.langchain.com") {
		t.Errorf("expected default endpoint in URL, got: %s", url)
	}
}
