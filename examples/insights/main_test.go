package main

import (
	"encoding/json"
	"testing"
)

func TestBuildInsightsParams(t *testing.T) {
	data, err := json.Marshal(buildInsightsParams())
	if err != nil {
		t.Fatalf("marshaling params: %v", err)
	}

	var body map[string]any
	if err := json.Unmarshal(data, &body); err != nil {
		t.Fatalf("unmarshaling params: %v", err)
	}

	if got := body["name"]; got != reportName {
		t.Errorf("expected report name %q, got %v", reportName, got)
	}
	if _, ok := body["schedule_cron"]; ok {
		t.Errorf("unexpected recurring schedule in request: %s", data)
	}

	config, ok := body["config"].(map[string]any)
	if !ok {
		t.Fatalf("expected config object, got %T", body["config"])
	}
	if got := config["last_n_hours"]; got != float64(24) {
		t.Errorf("expected a 24-hour lookback, got %v", got)
	}
	if got := config["summary_prompt"]; got != "Summarize the user's intent and any failures." {
		t.Errorf("unexpected summary prompt: %v", got)
	}
	if _, ok := config["is_scheduled"]; ok {
		t.Errorf("unexpected scheduled flag in request: %s", data)
	}
}
