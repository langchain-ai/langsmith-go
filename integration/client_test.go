//go:build integration

package integration

import (
	"context"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/shared"
)

// dottedOrder returns a LangSmith dotted_order string for a root run.
// Format: YYYYMMDDTHHMMSSmmmmmmZ<run_id> (no dots — single part for root runs).
func dottedOrder(t time.Time, runID string) string {
	ts := t.UTC()
	return fmt.Sprintf("%s%06dZ%s", ts.Format("20060102T150405"), ts.Nanosecond()/1000, runID)
}

// childDottedOrder returns a dotted_order for a child run:
// {parentDottedOrder}.{childTimestamp}Z{childID}
func childDottedOrder(parentOrder string, t time.Time, childID string) string {
	ts := t.UTC()
	return fmt.Sprintf("%s.%s%06dZ%s", parentOrder, ts.Format("20060102T150405"), ts.Nanosecond()/1000, childID)
}

// --- Datasets CRUD ---

func TestDatasetCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	name := uniqueName("go-integ-dataset")

	// Create
	dataset, err := client.Datasets.New(ctx, langsmith.DatasetNewParams{
		Name:        langsmith.F(name),
		Description: langsmith.F("Integration test dataset"),
	})
	if err != nil {
		t.Fatalf("create dataset: %v", err)
	}
	if dataset.ID == "" {
		t.Fatal("expected non-empty dataset ID")
	}
	if dataset.Name != name {
		t.Errorf("name = %q, want %q", dataset.Name, name)
	}

	// Get
	got, err := client.Datasets.Get(ctx, dataset.ID)
	if err != nil {
		t.Fatalf("get dataset: %v", err)
	}
	if got.ID != dataset.ID {
		t.Errorf("get returned different ID: %s vs %s", got.ID, dataset.ID)
	}

	// List with name filter
	listed, err := client.Datasets.List(ctx, langsmith.DatasetListParams{
		Name: langsmith.F(name),
	})
	if err != nil {
		t.Fatalf("list datasets: %v", err)
	}
	if len(listed.Items) == 0 {
		t.Error("expected at least one dataset in list")
	}

	// Delete
	_, err = client.Datasets.Delete(ctx, dataset.ID)
	if err != nil {
		t.Fatalf("delete dataset: %v", err)
	}
}

// --- Examples CRUD ---

func TestExampleCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	dsName := uniqueName("go-integ-examples")

	dataset, err := client.Datasets.New(ctx, langsmith.DatasetNewParams{
		Name: langsmith.F(dsName),
	})
	if err != nil {
		t.Fatalf("create dataset: %v", err)
	}
	defer client.Datasets.Delete(ctx, dataset.ID)

	// Create individual example
	example, err := client.Examples.New(ctx, langsmith.ExampleNewParams{
		DatasetID: langsmith.F(dataset.ID),
		Inputs:    langsmith.F(map[string]interface{}{"question": "What is Go?"}),
		Outputs:   langsmith.F(map[string]interface{}{"answer": "A programming language."}),
	})
	if err != nil {
		t.Fatalf("create example: %v", err)
	}
	if example.ID == "" {
		t.Fatal("expected non-empty example ID")
	}

	// Get example
	got, err := client.Examples.Get(ctx, example.ID, langsmith.ExampleGetParams{})
	if err != nil {
		t.Fatalf("get example: %v", err)
	}
	if got.ID != example.ID {
		t.Errorf("get returned different ID")
	}

	// Update example
	_, err = client.Examples.Update(ctx, example.ID, langsmith.ExampleUpdateParams{
		Outputs: langsmith.F(map[string]interface{}{"answer": "A statically typed language."}),
	})
	if err != nil {
		t.Fatalf("update example: %v", err)
	}

	// Verify update persisted
	updated, err := client.Examples.Get(ctx, example.ID, langsmith.ExampleGetParams{})
	if err != nil {
		t.Fatalf("get updated example: %v", err)
	}
	if ans, ok := updated.Outputs["answer"]; !ok || ans != "A statically typed language." {
		t.Errorf("updated outputs = %v, want answer='A statically typed language.'", updated.Outputs)
	}

	// List examples
	listed, err := client.Examples.List(ctx, langsmith.ExampleListParams{
		Dataset: langsmith.F(dataset.ID),
	})
	if err != nil {
		t.Fatalf("list examples: %v", err)
	}
	if len(listed.Items) != 1 {
		t.Errorf("expected 1 example, got %d", len(listed.Items))
	}

	// Delete example
	_, err = client.Examples.Delete(ctx, example.ID)
	if err != nil {
		t.Fatalf("delete example: %v", err)
	}
}

func TestExampleBulkCreate(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	dsName := uniqueName("go-integ-bulk")

	dataset, err := client.Datasets.New(ctx, langsmith.DatasetNewParams{
		Name: langsmith.F(dsName),
	})
	if err != nil {
		t.Fatalf("create dataset: %v", err)
	}
	defer client.Datasets.Delete(ctx, dataset.ID)

	results, err := client.Examples.Bulk.New(ctx, langsmith.ExampleBulkNewParams{
		Body: []langsmith.ExampleBulkNewParamsBody{
			{
				DatasetID: langsmith.F(dataset.ID),
				Inputs:    langsmith.F(map[string]interface{}{"q": "What is Go?"}),
				Outputs:   langsmith.F(map[string]interface{}{"a": "A language."}),
			},
			{
				DatasetID: langsmith.F(dataset.ID),
				Inputs:    langsmith.F(map[string]interface{}{"q": "What is Rust?"}),
				Outputs:   langsmith.F(map[string]interface{}{"a": "A language."}),
			},
		},
	})
	if err != nil {
		t.Fatalf("bulk create: %v", err)
	}
	if results == nil || len(*results) != 2 {
		t.Errorf("expected 2 bulk results, got %v", results)
	}
}

// --- Runs Batch Ingest & Query ---

func TestRunIngestAndQuery(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	sessionName := getSharedIntegrationProject()

	runID := uuid.New().String()
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)

	_, err := client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Post: langsmith.F([]langsmith.RunParam{
			{
				ID:          langsmith.F(runID),
				TraceID:     langsmith.F(runID),
				DottedOrder: langsmith.F(dottedOrder(now, runID)),
				Name:        langsmith.F("test-run"),
				RunType:     langsmith.F(langsmith.RunRunTypeChain),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"input": "hello"}),
				Outputs:     langsmith.F(map[string]interface{}{"output": "world"}),
			},
		}),
	})
	if err != nil {
		t.Fatalf("ingest batch: %v", err)
	}

	// Query with retries (runs may take a moment to be queryable), then assert round-trip like Python read_run
	var run *langsmith.RunQueryResponseRun
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		result, err := client.Runs.Query(ctx, langsmith.RunQueryParams{
			ID: langsmith.F([]string{runID}),
		})
		if err != nil {
			continue
		}
		if len(result.Runs) > 0 {
			run = &result.Runs[0]
			break
		}
	}
	if run == nil {
		t.Fatal("run not found after polling")
	}
	if run.ID != runID {
		t.Errorf("stored_run.id = %q, want %q", run.ID, runID)
	}
	if run.Name != "test-run" {
		t.Errorf("stored_run.name = %q, want 'test-run'", run.Name)
	}
	if run.StartTime.IsZero() {
		t.Error("stored_run.start_time should be set")
	}
	if run.Inputs == nil {
		t.Fatal("stored_run.inputs should be set")
	}
	if v, ok := run.Inputs["input"]; !ok || v != "hello" {
		t.Errorf("stored_run.inputs[input] = %v, want 'hello'", run.Inputs["input"])
	}
	if run.Outputs == nil {
		t.Fatal("stored_run.outputs should be set")
	}
	if v, ok := run.Outputs["output"]; !ok || v != "world" {
		t.Errorf("stored_run.outputs[output] = %v, want 'world'", run.Outputs["output"])
	}
	if run.RunType != langsmith.RunQueryResponseRunsRunTypeChain {
		t.Errorf("stored_run.run_type = %q, want 'chain'", run.RunType)
	}
}

// TestRunCreateAndUpdate mirrors Python: create run, update run (patch), then read_run and assert.
func TestRunCreateAndUpdate(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	sessionName := getSharedIntegrationProject()

	runID := uuid.New().String()
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)

	// Create run
	_, err := client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Post: langsmith.F([]langsmith.RunParam{
			{
				ID:          langsmith.F(runID),
				TraceID:     langsmith.F(runID),
				DottedOrder: langsmith.F(dottedOrder(now, runID)),
				Name:        langsmith.F("create-update-run"),
				RunType:     langsmith.F(langsmith.RunRunTypeChain),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"x": 1}),
				Outputs:     langsmith.F(map[string]interface{}{"y": 2}),
			},
		}),
	})
	if err != nil {
		t.Fatalf("ingest run: %v", err)
	}

	// Update run (patch) with new outputs (trace_id, dotted_order, end_time required by batch patch API)
	dotted := dottedOrder(now, runID)
	_, err = client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Patch: langsmith.F([]langsmith.RunParam{
			{
				ID:          langsmith.F(runID),
				TraceID:     langsmith.F(runID),
				DottedOrder: langsmith.F(dotted),
				EndTime:     langsmith.F(nowStr),
				SessionName: langsmith.F(sessionName),
				Outputs:     langsmith.F(map[string]interface{}{"y": 99, "z": "updated"}),
			},
		}),
	})
	if err != nil {
		t.Fatalf("patch run: %v", err)
	}

	// Poll until run is visible and patch has been applied (eventual consistency)
	var stored *langsmith.RunQueryResponseRun
	for i := 0; i < 15; i++ {
		time.Sleep(2 * time.Second)
		result, err := client.Runs.Query(ctx, langsmith.RunQueryParams{
			ID: langsmith.F([]string{runID}),
		})
		if err != nil {
			continue
		}
		if len(result.Runs) == 0 {
			continue
		}
		r := &result.Runs[0]
		stored = r
		// Patch applied when outputs contain updated values
		if r.Outputs != nil {
			vy := r.Outputs["y"]
			vz := r.Outputs["z"]
			if vy != nil && vz != nil {
				yOk := false
				switch n := vy.(type) {
				case float64:
					yOk = n == 99
				case int:
					yOk = n == 99
				}
				if !yOk && vy == "99" {
					yOk = true
				}
				if yOk && vz == "updated" {
					break
				}
			}
		}
	}
	if stored == nil {
		t.Fatal("run not found after polling")
	}
	if stored.ID != runID {
		t.Errorf("stored_run.id = %q, want %q", stored.ID, runID)
	}
	if stored.Name != "create-update-run" {
		t.Errorf("stored_run.name = %q, want 'create-update-run'", stored.Name)
	}
	if stored.StartTime.IsZero() {
		t.Error("stored_run.start_time should be set")
	}
	if stored.Outputs == nil {
		t.Fatal("stored_run.outputs should be set")
	}
	if v, ok := stored.Outputs["y"]; !ok {
		t.Error("stored_run.outputs should have 'y'")
	} else if n, ok := v.(float64); ok {
		if n != 99 {
			t.Errorf("stored_run.outputs[y] = %v, want 99", v)
		}
	} else if v != "99" {
		t.Errorf("stored_run.outputs[y] = %v, want 99", v)
	}
	if v, ok := stored.Outputs["z"]; !ok || v != "updated" {
		t.Errorf("stored_run.outputs[z] = %v, want 'updated'", stored.Outputs["z"])
	}
}

// TestRunBatchIngestRoundTrip mirrors Python test_batch_ingest_runs: batch ingest 3 runs
// (2 root + 1 child) with distinct inputs/outputs/error, then list_runs and assert round-trip.
func TestRunBatchIngestRoundTrip(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	sessionName := getSharedIntegrationProject()

	run1ID := uuid.New().String()
	run2ID := uuid.New().String()
	run3ID := uuid.New().String()
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)
	do1 := dottedOrder(now, run1ID)
	do2 := dottedOrder(now.Add(time.Millisecond), run2ID)
	do3 := childDottedOrder(do1, now.Add(2*time.Millisecond), run3ID)

	_, err := client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Post: langsmith.F([]langsmith.RunParam{
			{
				ID:          langsmith.F(run1ID),
				TraceID:     langsmith.F(run1ID),
				DottedOrder: langsmith.F(do1),
				Name:        langsmith.F("run1"),
				RunType:     langsmith.F(langsmith.RunRunTypeChain),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"input1": 1, "input2": 2}),
				Outputs:     langsmith.F(map[string]interface{}{"output1": 3, "output2": 4}),
			},
			{
				ID:          langsmith.F(run2ID),
				TraceID:     langsmith.F(run2ID),
				DottedOrder: langsmith.F(do2),
				Name:        langsmith.F("run2"),
				RunType:     langsmith.F(langsmith.RunRunTypeChain),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"input1": 5, "input2": 6}),
				Outputs:     langsmith.F(map[string]interface{}{"output1": 4, "output2": 5}),
			},
			{
				ID:          langsmith.F(run3ID),
				TraceID:     langsmith.F(run1ID),
				DottedOrder: langsmith.F(do3),
				ParentRunID: langsmith.F(run1ID),
				Name:        langsmith.F("run3"),
				RunType:     langsmith.F(langsmith.RunRunTypeLlm),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"input1": 1, "input2": 2}),
				Error:       langsmith.F("error"),
			},
		}),
	})
	if err != nil {
		t.Fatalf("batch ingest: %v", err)
	}

	// Query by run IDs (shared project has many runs; we only need our three)
	var runs []langsmith.RunQueryResponseRun
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		result, err := client.Runs.Query(ctx, langsmith.RunQueryParams{
			ID: langsmith.F([]string{run1ID, run2ID, run3ID}),
		})
		if err != nil {
			continue
		}
		if len(result.Runs) >= 3 {
			runs = result.Runs
			break
		}
	}
	if len(runs) != 3 {
		t.Fatalf("expected exactly 3 runs, got %d", len(runs))
	}

	byID := make(map[string]langsmith.RunQueryResponseRun)
	for _, r := range runs {
		byID[r.ID] = r
	}

	// run1: inputs {input1:1, input2:2}, outputs {output1:3, output2:4}
	if r, ok := byID[run1ID]; !ok {
		t.Errorf("run1 %q not found", run1ID)
	} else {
		assertMap(t, "run1.inputs", r.Inputs, map[string]interface{}{"input1": float64(1), "input2": float64(2)})
		assertMap(t, "run1.outputs", r.Outputs, map[string]interface{}{"output1": float64(3), "output2": float64(4)})
	}
	// run2: inputs {input1:5, input2:6}, outputs {output1:4, output2:5}
	if r, ok := byID[run2ID]; !ok {
		t.Errorf("run2 %q not found", run2ID)
	} else {
		assertMap(t, "run2.inputs", r.Inputs, map[string]interface{}{"input1": float64(5), "input2": float64(6)})
		assertMap(t, "run2.outputs", r.Outputs, map[string]interface{}{"output1": float64(4), "output2": float64(5)})
	}
	// run3: inputs {input1:1, input2:2}, error "error"
	if r, ok := byID[run3ID]; !ok {
		t.Errorf("run3 %q not found", run3ID)
	} else {
		assertMap(t, "run3.inputs", r.Inputs, map[string]interface{}{"input1": float64(1), "input2": float64(2)})
		if r.Error != "error" {
			t.Errorf("run3.error = %q, want 'error'", r.Error)
		}
	}
}

func assertMap(t *testing.T, label string, got map[string]interface{}, want map[string]interface{}) {
	t.Helper()
	if got == nil {
		t.Errorf("%s = nil, want %v", label, want)
		return
	}
	for k, vWant := range want {
		vGot, ok := got[k]
		if !ok {
			t.Errorf("%s[%q] missing, want %v", label, k, vWant)
			continue
		}
		// JSON unmarshals numbers as float64
		switch w := vWant.(type) {
		case float64:
			if g, ok := vGot.(float64); !ok || g != w {
				t.Errorf("%s[%q] = %v, want %v", label, k, vGot, vWant)
			}
		default:
			if vGot != vWant {
				t.Errorf("%s[%q] = %v, want %v", label, k, vGot, vWant)
			}
		}
	}
}

func TestRunIngestWithChildRuns(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	sessionName := getSharedIntegrationProject()

	parentID := uuid.New().String()
	childID := uuid.New().String()
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)
	parentDO := dottedOrder(now, parentID)
	childDO := childDottedOrder(parentDO, now.Add(time.Millisecond), childID)

	_, err := client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Post: langsmith.F([]langsmith.RunParam{
			{
				ID:          langsmith.F(parentID),
				TraceID:     langsmith.F(parentID),
				DottedOrder: langsmith.F(parentDO),
				Name:        langsmith.F("parent-chain"),
				RunType:     langsmith.F(langsmith.RunRunTypeChain),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"input": "hello"}),
				Outputs:     langsmith.F(map[string]interface{}{"output": "world"}),
				Tags:        langsmith.F([]string{"foo", "bar"}),
			},
			{
				ID:          langsmith.F(childID),
				TraceID:     langsmith.F(parentID),
				DottedOrder: langsmith.F(childDO),
				ParentRunID: langsmith.F(parentID),
				Name:        langsmith.F("child-llm"),
				RunType:     langsmith.F(langsmith.RunRunTypeLlm),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"prompt": "hello"}),
				Outputs:     langsmith.F(map[string]interface{}{"text": "world"}),
			},
		}),
	})
	if err != nil {
		t.Fatalf("ingest batch with child: %v", err)
	}

	var parentRun, childRun *langsmith.RunQueryResponseRun
	for i := 0; i < 10; i++ {
		time.Sleep(2 * time.Second)
		result, err := client.Runs.Query(ctx, langsmith.RunQueryParams{
			ID: langsmith.F([]string{parentID, childID}),
		})
		if err != nil {
			continue
		}
		for i := range result.Runs {
			r := &result.Runs[i]
			if r.ID == parentID {
				parentRun = r
			}
			if r.ID == childID {
				childRun = r
			}
		}
		if parentRun != nil && childRun != nil {
			break
		}
	}
	if parentRun == nil {
		t.Fatal("parent run not found after polling")
	}
	if childRun == nil {
		t.Fatal("child run not found after polling")
	}

	// Nested run structure (like Python test_nested_runs)
	if parentRun.ParentRunID != "" {
		t.Errorf("parent parent_run_id = %q, want empty (root)", parentRun.ParentRunID)
	}
	if childRun.ParentRunID != parentID {
		t.Errorf("child parent_run_id = %q, want %q", childRun.ParentRunID, parentID)
	}
	if parentRun.RunType != langsmith.RunQueryResponseRunsRunTypeChain {
		t.Errorf("parent run_type = %q, want 'chain'", parentRun.RunType)
	}
	if childRun.RunType != langsmith.RunQueryResponseRunsRunTypeLlm {
		t.Errorf("child run_type = %q, want 'llm'", childRun.RunType)
	}
	if parentRun.Name != "parent-chain" {
		t.Errorf("parent name = %q, want 'parent-chain'", parentRun.Name)
	}
	if childRun.Name != "child-llm" {
		t.Errorf("child name = %q, want 'child-llm'", childRun.Name)
	}
	// Tags on chain run (order may vary)
	tagSet := make(map[string]bool)
	for _, tag := range parentRun.Tags {
		tagSet[tag] = true
	}
	if !tagSet["foo"] || !tagSet["bar"] {
		t.Errorf("parent tags = %v, want to include [foo, bar]", parentRun.Tags)
	}
	// Inputs on child
	if childRun.Inputs == nil {
		t.Fatal("child inputs should be set")
	}
	if v, ok := childRun.Inputs["prompt"]; !ok || v != "hello" {
		t.Errorf("child inputs[prompt] = %v, want 'hello'", childRun.Inputs["prompt"])
	}
	if v, ok := childRun.Outputs["text"]; !ok || v != "world" {
		t.Errorf("child outputs[text] = %v, want 'world'", childRun.Outputs["text"])
	}
}

// --- Feedback CRUD ---

func TestFeedbackCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	sessionName := getSharedIntegrationProject()

	runID := uuid.New().String()
	now := time.Now().UTC()
	nowStr := now.Format(time.RFC3339Nano)

	_, err := client.Runs.IngestBatch(ctx, langsmith.RunIngestBatchParams{
		Post: langsmith.F([]langsmith.RunParam{
			{
				ID:          langsmith.F(runID),
				TraceID:     langsmith.F(runID),
				DottedOrder: langsmith.F(dottedOrder(now, runID)),
				Name:        langsmith.F("feedback-test-run"),
				RunType:     langsmith.F(langsmith.RunRunTypeChain),
				SessionName: langsmith.F(sessionName),
				StartTime:   langsmith.F(nowStr),
				EndTime:     langsmith.F(nowStr),
				Inputs:      langsmith.F(map[string]interface{}{"input": "test"}),
			},
		}),
	})
	if err != nil {
		t.Fatalf("ingest run for feedback: %v", err)
	}

	// Wait for run to be available
	time.Sleep(3 * time.Second)

	// Create feedback
	fb, err := client.Feedback.New(ctx, langsmith.FeedbackNewParams{
		FeedbackCreateSchema: langsmith.FeedbackCreateSchemaParam{
			Key:   langsmith.F("correctness"),
			RunID: langsmith.F(runID),
			Score: langsmith.F[langsmith.FeedbackCreateSchemaScoreUnionParam](shared.UnionFloat(0.9)),
		},
	})
	if err != nil {
		t.Fatalf("create feedback: %v", err)
	}
	if fb.ID == "" {
		t.Fatal("expected non-empty feedback ID")
	}

	// Get feedback and verify both key and score
	got, err := client.Feedback.Get(ctx, fb.ID, langsmith.FeedbackGetParams{})
	if err != nil {
		t.Fatalf("get feedback: %v", err)
	}
	if got.Key != "correctness" {
		t.Errorf("feedback key = %q, want 'correctness'", got.Key)
	}
	if score, ok := got.Score.(shared.UnionFloat); !ok {
		t.Errorf("feedback score type = %T, want shared.UnionFloat", got.Score)
	} else if float64(score) != 0.9 {
		t.Errorf("feedback score = %v, want 0.9", score)
	}

	// Delete feedback
	_, err = client.Feedback.Delete(ctx, fb.ID)
	if err != nil {
		t.Fatalf("delete feedback: %v", err)
	}
}

// --- Sessions CRUD ---

func TestSessionCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	name := uniqueName("go-integ-session")

	session, err := client.Sessions.New(ctx, langsmith.SessionNewParams{
		Name:        langsmith.F(name),
		Description: langsmith.F("Integration test session"),
	})
	if err != nil {
		t.Fatalf("create session: %v", err)
	}
	if session.ID == "" {
		t.Fatal("expected non-empty session ID")
	}

	// Get session
	got, err := client.Sessions.Get(ctx, session.ID, langsmith.SessionGetParams{})
	if err != nil {
		t.Fatalf("get session: %v", err)
	}
	if got.Name != name {
		t.Errorf("name = %q, want %q", got.Name, name)
	}

	// List sessions
	listed, err := client.Sessions.List(ctx, langsmith.SessionListParams{
		Name: langsmith.F(name),
	})
	if err != nil {
		t.Fatalf("list sessions: %v", err)
	}
	if len(listed.Items) == 0 {
		t.Error("expected at least one session in list")
	}

	// Delete session — the API returns an empty body, which the SDK cannot
	// deserialize into interface{}. Tolerate that specific error.
	_, err = client.Sessions.Delete(ctx, session.ID)
	if err != nil && !strings.Contains(err.Error(), "expected destination type") {
		t.Fatalf("delete session: %v", err)
	}
}
