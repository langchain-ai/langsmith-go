package cmd

import (
	"encoding/json"
	"errors"
	"fmt"
	"maps"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
	"time"

	langsmith "github.com/langchain-ai/langsmith-go"
)

// testRule is a minimal JSON-serializable struct used by mock server handlers
// to produce responses that the SDK can decode into langsmith.Rule.
type testRule struct {
	ID             string         `json:"id"`
	DisplayName    string         `json:"display_name"`
	SamplingRate   float64        `json:"sampling_rate"`
	IsEnabled      bool           `json:"is_enabled"`
	DatasetID      string         `json:"dataset_id,omitempty"`
	ProjectID      string         `json:"project_id,omitempty"`
	CodeEvaluators []testCodeEval `json:"code_evaluators,omitempty"`
	Evaluators     []testLLMEval  `json:"evaluators,omitempty"`
}

type testCodeEval struct {
	Code     string `json:"code"`
	Language string `json:"language"`
}

type testLLMEval struct {
	Structured testLLMStructured `json:"structured"`
}

type testLLMStructured struct {
	HubRef          string            `json:"hub_ref,omitempty"`
	VariableMapping map[string]string `json:"variable_mapping,omitempty"`
}

// ==================== Pure function tests ====================

// ---------- extractPythonFunction ----------

func TestExtractPythonFunction_Simple(t *testing.T) {
	source := `import os

def my_func(run, example):
    score = 1
    return {"score": score}

def other_func():
    pass
`
	result := extractPythonFunction(source, "my_func")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "def my_func(run, example):") {
		t.Error("should contain function definition")
	}
	if !contains(result, `return {"score": score}`) {
		t.Error("should contain return statement")
	}
	if contains(result, "def other_func") {
		t.Error("should not contain other function")
	}
}

func TestExtractPythonFunction_LastFunction(t *testing.T) {
	source := `def first():
    pass

def last_func(x, y):
    return x + y
`
	result := extractPythonFunction(source, "last_func")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "def last_func(x, y):") {
		t.Error("should contain last_func definition")
	}
}

func TestExtractPythonFunction_NotFound(t *testing.T) {
	source := `def other(x):
    return x
`
	result := extractPythonFunction(source, "missing_func")
	if result != "" {
		t.Errorf("expected empty string, got %q", result)
	}
}

func TestExtractPythonFunction_WithIndentedBlocks(t *testing.T) {
	source := `def evaluate(run, example):
    if run.error:
        return {"score": 0}
    else:
        return {"score": 1}

class MyClass:
    pass
`
	result := extractPythonFunction(source, "evaluate")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if contains(result, "class MyClass") {
		t.Error("should not include class definition")
	}
}

func TestExtractPythonFunction_EmptySource(t *testing.T) {
	result := extractPythonFunction("", "anything")
	if result != "" {
		t.Errorf("expected empty for empty source, got %q", result)
	}
}

// ---------- extractJSFunction ----------

func TestExtractJSFunction_FunctionDeclaration(t *testing.T) {
	source := `function myEval(run, example) {
  if (run.error) {
    return { score: 0 };
  }
  return { score: 1 };
}

function other() {
  return 42;
}
`
	result := extractJSFunction(source, "myEval")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "function myEval(run, example)") {
		t.Error("should contain function declaration")
	}
	if contains(result, "function other") {
		t.Error("should not contain other function")
	}
}

func TestExtractJSFunction_ArrowFunction(t *testing.T) {
	source := `const myEval = (run, example) => {
  return { score: 1 };
}
`
	result := extractJSFunction(source, "myEval")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "const myEval") {
		t.Error("should contain const declaration")
	}
}

func TestExtractJSFunction_AsyncFunction(t *testing.T) {
	source := `export async function checkAccuracy(run, example) {
  const result = await check(run);
  return { score: result ? 1 : 0 };
}
`
	result := extractJSFunction(source, "checkAccuracy")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "async function checkAccuracy") {
		t.Error("should contain async function")
	}
}

func TestExtractJSFunction_ExportedConst(t *testing.T) {
	source := `export const myEval = (run, example) => {
  return { score: 1 };
}
`
	result := extractJSFunction(source, "myEval")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
}

func TestExtractJSFunction_NotFound(t *testing.T) {
	source := `function other() { return 1; }`
	result := extractJSFunction(source, "missing")
	if result != "" {
		t.Errorf("expected empty, got %q", result)
	}
}

func TestExtractJSFunction_NestedBraces(t *testing.T) {
	source := `function eval(run) {
  if (run.error) {
    return { key: { nested: true } };
  }
  return { score: 1 };
}
`
	result := extractJSFunction(source, "eval")
	if result == "" {
		t.Fatal("expected non-empty result")
	}
	if !contains(result, "return { score: 1 };") {
		t.Error("should contain the full function body")
	}
}

// ---------- detectLanguage ----------

func TestDetectLanguage_Python(t *testing.T) {
	lang, funcName := detectLanguage("eval.py")
	if lang != "python" {
		t.Errorf("expected python, got %q", lang)
	}
	if funcName != "perform_eval" {
		t.Errorf("expected perform_eval, got %q", funcName)
	}
}

func TestDetectLanguage_JavaScript(t *testing.T) {
	for _, ext := range []string{".js", ".ts", ".tsx", ".mjs"} {
		lang, funcName := detectLanguage("eval" + ext)
		if lang != "javascript" {
			t.Errorf("detectLanguage(eval%s): expected javascript, got %q", ext, lang)
		}
		if funcName != "performEval" {
			t.Errorf("detectLanguage(eval%s): expected performEval, got %q", ext, funcName)
		}
	}
}

func TestDetectLanguage_Unknown(t *testing.T) {
	lang, funcName := detectLanguage("eval.rb")
	if lang != "" || funcName != "" {
		t.Errorf("expected empty for .rb, got (%q, %q)", lang, funcName)
	}
}

func TestDetectLanguage_CaseInsensitive(t *testing.T) {
	lang, _ := detectLanguage("eval.PY")
	if lang != "python" {
		t.Errorf("expected python for .PY, got %q", lang)
	}
}

// ---------- renameJSFunction ----------

func TestRenameJSFunction_FunctionDecl(t *testing.T) {
	source := `function myEval(run, example) {
  return { score: 1 };
}
`
	result := renameJSFunction(source, "myEval")
	if !contains(result, "function performEval(") {
		t.Errorf("expected function renamed to performEval, got:\n%s", result)
	}
	if contains(result, "myEval") {
		t.Error("original name should be replaced")
	}
}

func TestRenameJSFunction_AsyncFunctionDecl(t *testing.T) {
	source := `async function myEval(run) {
  return { score: 1 };
}
`
	result := renameJSFunction(source, "myEval")
	if !contains(result, "async function performEval(") {
		t.Errorf("expected async function performEval, got:\n%s", result)
	}
}

func TestRenameJSFunction_ExportFunction(t *testing.T) {
	source := `export function myEval(run) {
  return { score: 1 };
}
`
	result := renameJSFunction(source, "myEval")
	if !contains(result, "function performEval(") {
		t.Errorf("expected export stripped, got:\n%s", result)
	}
	if contains(result, "export") {
		t.Error("export keyword should be removed")
	}
}

func TestRenameJSFunction_ArrowFunction(t *testing.T) {
	source := `const myEval = (run) => {
  return { score: 1 };
}
`
	result := renameJSFunction(source, "myEval")
	if !contains(result, "function performEval(run) {") {
		t.Errorf("expected arrow converted to function decl, got:\n%s", result)
	}
}

func TestRenameJSFunction_AsyncArrowFunction(t *testing.T) {
	source := `export const myEval = async (run, example) => {
  return { score: 1 };
}
`
	result := renameJSFunction(source, "myEval")
	if !contains(result, "async function performEval(run, example) {") {
		t.Errorf("expected async arrow converted, got:\n%s", result)
	}
}

// ---------- findRuleForReplace ----------

func TestFindRuleForReplace(t *testing.T) {
	code := []langsmith.RuleCodeEvaluator{{Code: "x"}}
	rules := []langsmith.Rule{
		{ID: "renamed", DisplayName: "test", EvaluatorID: "ev-1", EvaluatorName: "ready_for_task_grade", ProjectID: "proj-1"},
		{ID: "legacy", DisplayName: "accuracy", CodeEvaluators: code, DatasetID: "ds-1"},
		{ID: "other-target", DisplayName: "accuracy", CodeEvaluators: code, DatasetID: "ds-2"},
		{ID: "webhook", DisplayName: "notify", ProjectID: "proj-1"},
		{ID: "dup-a", DisplayName: "dup", EvaluatorID: "ev-2", EvaluatorName: "dup", ProjectID: "proj-2"},
		{ID: "dup-b", DisplayName: "dup-rule", EvaluatorID: "ev-3", EvaluatorName: "dup", ProjectID: "proj-2"},
	}
	tests := []struct {
		name, match, datasetID, projectID string
		wantID                            string
		wantErr                           bool
	}{
		{name: "evaluator name after a UI rename", match: "ready_for_task_grade", projectID: "proj-1", wantID: "renamed"},
		{name: "rule name", match: "test", projectID: "proj-1", wantID: "renamed"},
		{name: "legacy inline rule on its dataset", match: "accuracy", datasetID: "ds-1", wantID: "legacy"},
		{name: "name on another target", match: "accuracy", datasetID: "ds-other"},
		{name: "non-evaluator rule is never replaced", match: "notify", projectID: "proj-1"},
		{name: "no match", match: "missing", projectID: "proj-1"},
		{name: "ambiguous match", match: "dup", projectID: "proj-2", wantErr: true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := findRuleForReplace(rules, tt.match, tt.datasetID, tt.projectID)
			if tt.wantErr {
				if err == nil || !strings.Contains(err.Error(), "dup-a") || !strings.Contains(err.Error(), "dup-b") {
					t.Fatalf("expected ambiguity error listing both rules, got %v", err)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			gotID := ""
			if got != nil {
				gotID = got.ID
			}
			if gotID != tt.wantID {
				t.Errorf("expected %q, got %q", tt.wantID, gotID)
			}
		})
	}
}

func TestValidateEvaluatorTargetFlags(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name      string
		dataset   string
		project   string
		projectID string
		wantErr   string
	}{
		{name: "requires one target", wantErr: "must specify"},
		{name: "rejects dataset and project", dataset: "ds", project: "proj", wantErr: "only one of"},
		{name: "rejects dataset and project-id", dataset: "ds", projectID: "id", wantErr: "only one of"},
		{name: "rejects project and project-id", project: "proj", projectID: "id", wantErr: "only one of"},
		{name: "rejects all three", dataset: "ds", project: "proj", projectID: "id", wantErr: "only one of"},
		{name: "dataset only", dataset: "ds"},
		{name: "project only", project: "proj"},
		{name: "project-id only", projectID: "id"},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			err := validateEvaluatorTargetFlags(tt.dataset, tt.project, tt.projectID)
			if tt.wantErr == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), tt.wantErr) {
				t.Fatalf("expected %q error, got %v", tt.wantErr, err)
			}
		})
	}
}

func TestLoadPromptMessagesFromJSON(t *testing.T) {
	t.Parallel()

	t.Run("pair format", func(t *testing.T) {
		t.Parallel()
		msgs, err := loadPromptMessagesFromJSON("prompt.json", []byte(
			`[["system","You are a judge."],["user","Q: {{input}}"]]`,
		))
		if err != nil {
			t.Fatal(err)
		}
		if len(msgs) != 2 || msgs[0][0] != "system" {
			t.Fatalf("unexpected messages: %#v", msgs)
		}
	})

	t.Run("role content objects", func(t *testing.T) {
		t.Parallel()
		msgs, err := loadPromptMessagesFromJSON("prompt.json", []byte(
			`[{"role":"system","content":"You are a judge."},{"role":"user","content":"Q: {{input}}"}]`,
		))
		if err != nil {
			t.Fatal(err)
		}
		if len(msgs) != 2 || msgs[1][1] != "Q: {{input}}" {
			t.Fatalf("unexpected messages: %#v", msgs)
		}
	})

	t.Run("invalid format", func(t *testing.T) {
		t.Parallel()
		_, err := loadPromptMessagesFromJSON("prompt.json", []byte(`{"messages":[]}`))
		if err == nil || !strings.Contains(err.Error(), "expected [[role,content]") {
			t.Fatalf("expected format hint, got %v", err)
		}
	})

	t.Run("rejects single-element pair", func(t *testing.T) {
		t.Parallel()
		_, err := loadPromptMessagesFromJSON("prompt.json", []byte(`[["system"]]`))
		if err == nil || !strings.Contains(err.Error(), "must be [role, content]") {
			t.Fatalf("expected length error, got %v", err)
		}
	})
}

func TestBuildLLMEvaluatorPayload_requiresModelConfig(t *testing.T) {
	t.Parallel()

	_, err := buildLLMEvaluatorPayload(
		"relevance", evaluatorTarget{projectID: "proj-1"},
		1.0, "", "", "prompt.json", "schema.json", "", nil,
	)
	if err == nil || !strings.Contains(err.Error(), "--model-config is required") {
		t.Fatalf("expected model-config error, got %v", err)
	}
}

// ==================== Command structure tests ====================

func TestEvaluatorCmd_Subcommands(t *testing.T) {
	cmd := newEvaluatorCmd()
	expected := map[string]bool{"get": false, "list": false, "upload": false, "create-llm": false, "delete": false, "rule": false}
	for _, sub := range cmd.Commands() {
		if _, ok := expected[sub.Name()]; ok {
			expected[sub.Name()] = true
		}
	}
	for name, found := range expected {
		if !found {
			t.Errorf("evaluator missing subcommand %q", name)
		}
	}
	for _, sub := range cmd.Commands() {
		if sub.Name() == "create" {
			t.Error("evaluator should not have 'create' subcommand")
		}
	}
}

func TestEvaluatorCmd_UseField(t *testing.T) {
	cmd := newEvaluatorCmd()
	if cmd.Use != "evaluators" {
		t.Errorf("expected Use=evaluators, got %q", cmd.Use)
	}
}

// ---------- Subcommand Use fields ----------

func TestEvaluatorListCmd_UseField(t *testing.T) {
	cmd := newEvaluatorListCmd()
	if cmd.Use != "list" {
		t.Errorf("expected Use=list, got %q", cmd.Use)
	}
}

func TestEvaluatorUploadCmd_UseField(t *testing.T) {
	cmd := newEvaluatorUploadCmd()
	if cmd.Use != "upload EVALUATOR_FILE" {
		t.Errorf("expected Use='upload EVALUATOR_FILE', got %q", cmd.Use)
	}
}

func TestEvaluatorDeleteCmd_UseField(t *testing.T) {
	cmd := newEvaluatorDeleteCmd()
	if cmd.Use != "delete EVALUATOR_ID" {
		t.Errorf("expected Use='delete EVALUATOR_ID', got %q", cmd.Use)
	}
}

// ---------- evaluator list flags ----------

func TestEvaluatorListCmd_Flags(t *testing.T) {
	cmd := newEvaluatorListCmd()
	f := cmd.Flags().Lookup("output")
	if f == nil {
		t.Fatal("--output flag not found")
	}
	if f.Shorthand != "o" {
		t.Errorf("expected shorthand 'o', got %q", f.Shorthand)
	}
	if f.DefValue != "" {
		t.Errorf("expected default empty, got %q", f.DefValue)
	}
}

// ---------- evaluator upload flags ----------

func TestEvaluatorUploadCmd_Flags(t *testing.T) {
	cmd := newEvaluatorUploadCmd()
	flags := map[string]string{
		"name":                    "",
		"function":                "",
		"dataset":                 "",
		"project":                 "",
		"sampling-rate":           "1",
		"trace-filter":            "",
		"replace":                 "false",
		"yes":                     "false",
		"advanced":                "false",
		"dependencies":            "",
		"input":                   "",
		"no-wait":                 "false",
		"ignore-credential-check": "false",
		"ignore-import-check":     "false",
	}
	for name, defVal := range flags {
		f := cmd.Flags().Lookup(name)
		if f == nil {
			t.Errorf("flag --%s not found", name)
			continue
		}
		if f.DefValue != defVal {
			t.Errorf("flag --%s: expected default %q, got %q", name, defVal, f.DefValue)
		}
	}
}

func TestEvaluatorUploadCmd_RequiredFlags(t *testing.T) {
	cmd := newEvaluatorUploadCmd()
	required := func(name string) bool {
		_, ok := cmd.Flags().Lookup(name).Annotations["cobra_annotation_bash_completion_one_required_flag"]
		return ok
	}
	if !required("name") {
		t.Error("flag --name not marked as required")
	}
	function := cmd.Flags().Lookup("function")
	if required("function") {
		t.Error("flag --function should be optional")
	}
	if _, ok := function.Annotations[deprecationAnnotation]; !ok || !function.Hidden {
		t.Error("flag --function should be deprecated and hidden")
	}
}

func TestEvaluatorUploadFunctionWarnsDeprecated(t *testing.T) {
	cmd := newEvaluatorUploadCmd()
	_ = cmd.Flags().Set("function", "grade")
	var buf strings.Builder
	warnDeprecated(&buf, cmd, nil, nil, true)
	var warning deprecationWarning
	if err := json.Unmarshal([]byte(buf.String()), &warning); err != nil {
		t.Fatalf("expected one JSON warning, got %q: %v", buf.String(), err)
	}
	if warning.Name != "--function" || warning.RemovedIn != "1.0.0" {
		t.Errorf("unexpected warning %+v", warning)
	}
}

func TestEvaluatorUploadFunctionRejectsWholeFileFlags(t *testing.T) {
	for _, flag := range [][2]string{
		{"advanced", "true"}, {"ignore-import-check", "true"}, {"dependencies", "requirements.txt"}, {"workspace-secret", "KEY"},
		{"input", "thread"}, {"no-wait", "true"},
	} {
		t.Run(flag[0], func(t *testing.T) {
			ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			})
			cleanup := setupTestEnv(t, ts.URL)
			defer cleanup()
			cmd := newEvaluatorUploadCmd()
			for k, v := range map[string]string{"name": "accuracy", "function": "grade", "project-id": testSessionID, flag[0]: flag[1]} {
				_ = cmd.Flags().Set(k, v)
			}
			err := runTestCommand(t, cmd, []string{"eval.py"})
			if err == nil || !strings.Contains(err.Error(), "--"+flag[0]+" cannot be combined with --function") {
				t.Fatalf("expected combination error, got %v", err)
			}
		})
	}
}

func TestEvaluatorUploadCmd_ExactArgs(t *testing.T) {
	cmd := newEvaluatorUploadCmd()
	if cmd.Args == nil {
		t.Fatal("expected Args validator")
	}
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error for 0 args")
	}
	if err := cmd.Args(cmd, []string{"eval.py"}); err != nil {
		t.Errorf("expected no error for 1 arg, got %v", err)
	}
	if err := cmd.Args(cmd, []string{"a", "b"}); err == nil {
		t.Error("expected error for 2 args")
	}
}

func TestEvaluatorUploadCmd_InvalidTargetReturnsError(t *testing.T) {
	cmd := newEvaluatorUploadCmd()
	_ = cmd.Flags().Set("dataset", "dataset")
	_ = cmd.Flags().Set("project", "project")

	runErr := runTestCommand(t, cmd, []string{"evaluator.py"})
	if runErr == nil || !contains(runErr.Error(), "only one of") {
		t.Fatalf("unexpected error: %v", runErr)
	}
}

// ---------- evaluator delete flags ----------

func TestEvaluatorDeleteCmd_Flags(t *testing.T) {
	cmd := newEvaluatorDeleteCmd()
	f := cmd.Flags().Lookup("yes")
	if f == nil {
		t.Fatal("--yes flag not found")
	}
	if f.DefValue != "false" {
		t.Errorf("expected default false, got %q", f.DefValue)
	}
}

func TestEvaluatorDeleteCmd_ExactArgs(t *testing.T) {
	cmd := newEvaluatorDeleteCmd()
	if err := cmd.Args(cmd, []string{}); err == nil {
		t.Error("expected error for 0 args")
	}
	if err := cmd.Args(cmd, []string{"my-eval"}); err != nil {
		t.Errorf("expected no error for 1 arg, got %v", err)
	}
}

// ==================== Execution tests ====================

func TestEvaluatorUploadReplacePatchesExistingCodeEvaluator(t *testing.T) {
	evaluatorFile := t.TempDir() + "/eval.py"
	if err := os.WriteFile(
		evaluatorFile,
		[]byte("def check_accuracy(run, example):\n    return {\"score\": 1}\n"),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	var sawDelete bool
	var patchBody map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/sessions" && r.Method == "GET":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"id": "project-1", "name": "my-project"},
			})
		case r.URL.Path == "/api/v2/runs/rules" && r.Method == "GET":
			_ = writeRulePage(w, []testRule{
				{
					ID:           "existing-rule",
					DisplayName:  "accuracy",
					SamplingRate: 0.25,
					IsEnabled:    true,
					ProjectID:    "project-1",
					CodeEvaluators: []testCodeEval{
						{Code: "def perform_eval(run, example):\n    return {\"score\": 0}", Language: "python"},
					},
				},
			})
		case r.URL.Path == "/api/v1/runs/rules/existing-rule" && r.Method == "PATCH":
			if err := json.NewDecoder(r.Body).Decode(&patchBody); err != nil {
				t.Fatalf("decoding patch body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":           "existing-rule",
				"display_name": "accuracy",
				"session_id":   "project-1",
			})
		case r.URL.Path == "/api/v1/runs/rules/existing-rule" && r.Method == "DELETE":
			sawDelete = true
			http.Error(w, "delete should not be called", http.StatusInternalServerError)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})

	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		cmd := newEvaluatorUploadCmd()
		_ = cmd.Flags().Set("name", "accuracy")
		_ = cmd.Flags().Set("function", "check_accuracy")
		_ = cmd.Flags().Set("project", "my-project")
		_ = cmd.Flags().Set("sampling-rate", "0.5")
		_ = cmd.Flags().Set("replace", "true")
		_ = cmd.Flags().Set("yes", "true")
		_ = runTestCommand(t, cmd, []string{evaluatorFile})
	})

	if sawDelete {
		t.Fatal("upload --replace should patch the existing evaluator, not delete it")
	}
	if patchBody == nil {
		t.Fatal("expected PATCH body")
	}
	if patchBody["display_name"] != "accuracy" {
		t.Errorf("expected display_name=accuracy, got %v", patchBody["display_name"])
	}
	if patchBody["session_id"] != "project-1" {
		t.Errorf("expected session_id=project-1, got %v", patchBody["session_id"])
	}
	if patchBody["sampling_rate"] != 0.5 {
		t.Errorf("expected sampling_rate=0.5, got %v", patchBody["sampling_rate"])
	}
	evaluators, ok := patchBody["code_evaluators"].([]any)
	if !ok || len(evaluators) != 1 {
		t.Fatalf("expected one code evaluator, got %#v", patchBody["code_evaluators"])
	}
	codeEvaluator, ok := evaluators[0].(map[string]any)
	if !ok {
		t.Fatalf("expected code evaluator object, got %#v", evaluators[0])
	}
	if codeEvaluator["language"] != "python" {
		t.Errorf("expected language=python, got %v", codeEvaluator["language"])
	}
	code, _ := codeEvaluator["code"].(string)
	if !strings.Contains(code, "def perform_eval(") {
		t.Errorf("expected uploaded code to be renamed to perform_eval, got:\n%s", code)
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if result["rule_id"] != "existing-rule" {
		t.Errorf("expected output rule_id=existing-rule, got %v", result["rule_id"])
	}
}

func TestEvaluatorCreateLLMReplacePatchesExistingEvaluator(t *testing.T) {
	modelConfigPath := t.TempDir() + "/model.json"
	if err := os.WriteFile(
		modelConfigPath,
		[]byte(`{"type":"chat","config":{"model":"test-model"}}`),
		0o600,
	); err != nil {
		t.Fatal(err)
	}

	var sawDelete bool
	var patchBody map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v1/sessions" && r.Method == "GET":
			_ = json.NewEncoder(w).Encode([]map[string]any{
				{"id": "project-1", "name": "my-project"},
			})
		case r.URL.Path == "/api/v2/runs/rules" && r.Method == "GET":
			_ = writeRulePage(w, []testRule{
				{
					ID:           "existing-rule",
					DisplayName:  "relevance",
					SamplingRate: 0.25,
					IsEnabled:    true,
					ProjectID:    "project-1",
					Evaluators: []testLLMEval{
						{Structured: testLLMStructured{HubRef: "my-org/old:latest"}},
					},
				},
			})
		case r.URL.Path == "/api/v1/runs/rules/existing-rule" && r.Method == "PATCH":
			if err := json.NewDecoder(r.Body).Decode(&patchBody); err != nil {
				t.Fatalf("decoding patch body: %v", err)
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":           "existing-rule",
				"display_name": "relevance",
				"session_id":   "project-1",
			})
		case r.URL.Path == "/api/v1/runs/rules/existing-rule" && r.Method == "DELETE":
			sawDelete = true
			http.Error(w, "delete should not be called", http.StatusInternalServerError)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})

	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		cmd := newEvaluatorCreateLLMCmd()
		_ = cmd.Flags().Set("name", "relevance")
		_ = cmd.Flags().Set("project", "my-project")
		_ = cmd.Flags().Set("hub-ref", "my-org/relevance:latest")
		_ = cmd.Flags().Set("model-config", modelConfigPath)
		_ = cmd.Flags().Set("sampling-rate", "0.5")
		_ = cmd.Flags().Set("replace", "true")
		_ = cmd.Flags().Set("yes", "true")
		_ = runTestCommand(t, cmd, nil)
	})

	if sawDelete {
		t.Fatal("create-llm --replace should patch the existing evaluator, not delete it")
	}
	if patchBody == nil {
		t.Fatal("expected PATCH body")
	}
	if patchBody["display_name"] != "relevance" || patchBody["evaluators"] == nil {
		t.Fatalf("expected LLM evaluator replacement payload, got %#v", patchBody)
	}

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if result["rule_id"] != "existing-rule" {
		t.Errorf("expected output rule_id=existing-rule, got %v", result["rule_id"])
	}
}

// ==================== evaluators and rules ====================
const (
	testEvaluatorID  = "c065b9d1-a6cb-4f2c-be93-46ed1ac6d8d0"
	testEvaluatorID2 = "5b1d7f0e-2f4a-4c55-9a0e-7d6c1e2b3a41"
	testRuleID       = "9f5c80d7-0983-438b-b7e9-4064fc5510ab"
	testRuleID2      = "0e8f2a6c-1b3d-4e5f-8a9b-0c1d2e3f4a5b"
	testSessionID    = "3d2c1b0a-9f8e-4d7c-b6a5-4f3e2d1c0b9a"
)

// serveEvaluatorPage answers GET /platform/evaluators with evaluators on the
// first page and an empty page after it, which ends auto-paging.
func serveEvaluatorPage(w http.ResponseWriter, r *http.Request, evaluators []map[string]any) {
	w.Header().Set("Content-Type", "application/json")
	if off := r.URL.Query().Get("offset"); off != "" && off != "0" {
		evaluators = []map[string]any{}
	}
	_ = json.NewEncoder(w).Encode(map[string]any{"evaluators": evaluators, "total": len(evaluators)})
}

func readyForTaskGrade() map[string]any {
	return map[string]any{
		"id":            testEvaluatorID,
		"name":          "ready_for_task_grade",
		"type":          "llm",
		"feedback_keys": []string{"ready"},
		"llm_evaluator": map[string]any{"prompt_repo_handle": "ready-for-task"},
		"run_rules": []map[string]any{
			{"id": testRuleID, "session_id": testSessionID, "session_name": "vanta-agent"},
		},
	}
}

func TestEvaluatorList_ShowsEvaluatorNameAndAttachedRules(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/platform/evaluators" && r.Method == http.MethodGet {
			serveEvaluatorPage(w, r, []map[string]any{readyForTaskGrade()})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		if err := runTestCommand(t, newEvaluatorListCmd(), nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 evaluator, got %d", len(result))
	}
	if result[0]["name"] != "ready_for_task_grade" || result[0]["id"] != testEvaluatorID {
		t.Errorf("expected evaluator identity, got %v", result[0])
	}
	rules, _ := result[0]["rules"].([]any)
	if len(rules) != 1 {
		t.Fatalf("expected 1 attached rule, got %v", result[0]["rules"])
	}
	rule := rules[0].(map[string]any)
	if rule["rule_id"] != testRuleID || rule["project"] != "vanta-agent" {
		t.Errorf("expected rule to point at vanta-agent, got %v", rule)
	}
}

func TestEvaluatorGet_ByNameRequiresExactMatch(t *testing.T) {
	var fetchedID string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v1/platform/evaluators" && r.Method == http.MethodGet:
			// The API's name filter is a substring match.
			serveEvaluatorPage(w, r, []map[string]any{
				{"id": testEvaluatorID2, "name": "ready_for_task_grade_v2", "type": "llm"},
				readyForTaskGrade(),
			})
		case strings.HasPrefix(r.URL.Path, "/api/v1/platform/evaluators/") && r.Method == http.MethodGet:
			fetchedID = strings.TrimPrefix(r.URL.Path, "/api/v1/platform/evaluators/")
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(readyForTaskGrade())
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	out := captureStdout(t, func() {
		if err := runTestCommand(t, newEvaluatorGetCmd(), []string{"ready_for_task_grade"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if fetchedID != testEvaluatorID {
		t.Errorf("expected to fetch %s, fetched %q", testEvaluatorID, fetchedID)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if result["name"] != "ready_for_task_grade" {
		t.Errorf("expected name=ready_for_task_grade, got %v", result["name"])
	}
	llm, _ := result["llm_evaluator"].(map[string]any)
	if llm["prompt_repo_handle"] != "ready-for-task" {
		t.Errorf("expected llm_evaluator detail, got %v", result["llm_evaluator"])
	}
}

func TestEvaluatorGet_RuleNamePointsToItsEvaluator(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/platform/evaluators":
			serveEvaluatorPage(w, r, nil)
		case "/api/v2/runs/rules":
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{{
				"id": testRuleID, "display_name": "test", "project_id": testSessionID,
				"evaluator_id": testEvaluatorID, "evaluator_name": "ready_for_task_grade",
			}})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	err := runTestCommand(t, newEvaluatorGetCmd(), []string{"test"})
	if err == nil || !strings.Contains(err.Error(), "ready_for_task_grade") || !strings.Contains(err.Error(), "rule name") {
		t.Fatalf("expected error naming the rule's evaluator, got %v", err)
	}
}

func TestEvaluatorGet_DuplicateNamesAsksForID(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v1/platform/evaluators" {
			serveEvaluatorPage(w, r, []map[string]any{
				{"id": testEvaluatorID, "name": "accuracy", "type": "code"},
				{"id": testEvaluatorID2, "name": "accuracy", "type": "llm"},
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	err := runTestCommand(t, newEvaluatorGetCmd(), []string{"accuracy"})
	if err == nil || !strings.Contains(err.Error(), testEvaluatorID) || !strings.Contains(err.Error(), testEvaluatorID2) {
		t.Fatalf("expected error listing both IDs, got %v", err)
	}
}

func TestEvaluatorDelete_RejectsNameWithoutCallingAPI(t *testing.T) {
	var calls int
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		http.Error(w, "unexpected", http.StatusInternalServerError)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	cmd := newEvaluatorDeleteCmd()
	_ = cmd.Flags().Set("yes", "true")
	err := runTestCommand(t, cmd, []string{"accuracy"})
	if err == nil || !strings.Contains(err.Error(), "evaluator rule delete") {
		t.Fatalf("expected error pointing to rule delete, got %v", err)
	}
	if calls != 0 {
		t.Errorf("expected no API calls, got %d", calls)
	}
}

func TestEvaluatorDelete_WithRulesRequiresDeleteRules(t *testing.T) {
	var sawDelete bool
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			sawDelete = true
		}
		if r.URL.Path == "/api/v1/platform/evaluators/"+testEvaluatorID && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(readyForTaskGrade())
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	cmd := newEvaluatorDeleteCmd()
	_ = cmd.Flags().Set("yes", "true")
	err := runTestCommand(t, cmd, []string{testEvaluatorID})
	if err == nil || !strings.Contains(err.Error(), "--delete-rules") || !strings.Contains(err.Error(), "vanta-agent") {
		t.Fatalf("expected --delete-rules error naming the project, got %v", err)
	}
	if sawDelete {
		t.Error("expected no DELETE request")
	}
}

func TestEvaluatorDelete_DeletesEvaluatorAndRules(t *testing.T) {
	var deleteQuery string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v1/platform/evaluators/"+testEvaluatorID {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		switch r.Method {
		case http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(readyForTaskGrade())
		case http.MethodDelete:
			deleteQuery = r.URL.RawQuery
			w.WriteHeader(http.StatusNoContent)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	out := captureStdout(t, func() {
		cmd := newEvaluatorDeleteCmd()
		_ = cmd.Flags().Set("delete-rules", "true")
		_ = cmd.Flags().Set("yes", "true")
		if err := runTestCommand(t, cmd, []string{testEvaluatorID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if deleteQuery != "delete_run_rules=true" {
		t.Errorf("expected delete_run_rules=true, got %q", deleteQuery)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if result["status"] != "deleted" || result["rules_deleted"] != float64(1) {
		t.Errorf("unexpected output: %v", result)
	}
}

func TestEvaluatorRuleList_ShowsRuleAndEvaluatorNames(t *testing.T) {
	var sessionFilter string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet {
			sessionFilter = r.URL.Query().Get("project_id")
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{{
				"id": testRuleID, "display_name": "test", "sampling_rate": 1.0, "is_enabled": true,
				"project_id": testSessionID, "project_name": "vanta-agent",
				"evaluator_id": testEvaluatorID, "evaluator_name": "ready_for_task_grade",
				"evaluator_version": 3,
			}})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		cmd := newEvaluatorRuleListCmd()
		_ = cmd.Flags().Set("project-id", testSessionID)
		if err := runTestCommand(t, cmd, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if sessionFilter != testSessionID {
		t.Errorf("expected project_id filter %s, got %q", testSessionID, sessionFilter)
	}
	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if len(result) != 1 {
		t.Fatalf("expected 1 rule, got %d", len(result))
	}
	got := result[0]
	if got["rule_id"] != testRuleID || got["rule_name"] != "test" {
		t.Errorf("expected rule identity, got %v", got)
	}
	if got["evaluator_id"] != testEvaluatorID || got["evaluator_name"] != "ready_for_task_grade" {
		t.Errorf("expected evaluator identity, got %v", got)
	}
	if _, ok := got["evaluator_version"]; ok {
		t.Error("evaluator_version is a rule format version and should not be shown")
	}
}

func TestEvaluatorGet_DuplicateNamesAlsoPointsToRuleEvaluator(t *testing.T) {
	const ruleEvaluatorID = "7e590565-6032-4117-8087-659a58ac172d"
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/api/v1/platform/evaluators":
			serveEvaluatorPage(w, r, []map[string]any{
				{"id": testEvaluatorID, "name": "test", "type": "code"},
				{"id": testEvaluatorID2, "name": "test", "type": "code"},
			})
		case "/api/v2/runs/rules":
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{
				{"id": testRuleID, "display_name": "test", "evaluator_id": testEvaluatorID, "evaluator_name": "test"},
				{"id": testRuleID2, "display_name": "test", "evaluator_id": ruleEvaluatorID, "evaluator_name": "ready_for_task_grade"},
			})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	err := runTestCommand(t, newEvaluatorGetCmd(), []string{"test"})
	if err == nil {
		t.Fatal("expected an ambiguity error")
	}
	msg := err.Error()
	for _, want := range []string{testEvaluatorID, testEvaluatorID2, "ready_for_task_grade", ruleEvaluatorID} {
		if !strings.Contains(msg, want) {
			t.Errorf("expected error to mention %q, got %v", want, msg)
		}
	}
	if strings.Count(msg, testEvaluatorID) != 1 {
		t.Errorf("expected evaluators named %q not to be repeated as rule hints, got %v", "test", msg)
	}
}

func TestEvaluatorRuleList_SendsFilters(t *testing.T) {
	var query url.Values
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet {
			query = r.URL.Query()
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{
				{"id": testRuleID, "display_name": "test", "project_id": testSessionID, "evaluator_id": testEvaluatorID, "evaluator_name": "ready_for_task_grade"},
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		cmd := newEvaluatorRuleListCmd()
		_ = cmd.Flags().Set("project-id", testSessionID)
		_ = cmd.Flags().Set("evaluator-id", testEvaluatorID)
		_ = cmd.Flags().Set("name-contains", "te")
		_ = cmd.Flags().Set("type", "project")
		if err := runTestCommand(t, cmd, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	for key, want := range map[string]string{
		"project_id":    testSessionID,
		"evaluator_id":  testEvaluatorID,
		"name_contains": "te",
		"type":          "PROJECT",
	} {
		if got := query.Get(key); got != want {
			t.Errorf("expected %s=%q, got %q", key, want, got)
		}
	}
	var result []map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if len(result) != 1 || result[0]["rule_id"] != testRuleID || result[0]["session_id"] != testSessionID {
		t.Fatalf("expected rule %s on project %s, got %v", testRuleID, testSessionID, result)
	}
}

func TestEvaluatorRuleList_RejectsUnknownType(t *testing.T) {
	cmd := newEvaluatorRuleListCmd()
	_ = cmd.Flags().Set("type", "thread")
	err := runTestCommand(t, cmd, nil)
	if err == nil || !strings.Contains(err.Error(), "--type") {
		t.Fatalf("expected a --type error, got %v", err)
	}
}

func TestEvaluatorRuleList_PagesWithLimitAndCursor(t *testing.T) {
	ids := []string{testRuleID, testRuleID2, testEvaluatorID, testEvaluatorID2}
	var cursors, pageSizes []string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/api/v2/runs/rules" {
			http.Error(w, "not found", http.StatusNotFound)
			return
		}
		cursor, pageSize := r.URL.Query().Get("cursor"), r.URL.Query().Get("page_size")
		cursors, pageSizes = append(cursors, cursor), append(pageSizes, pageSize)
		start := 0
		if cursor != "" {
			start, _ = strconv.Atoi(strings.TrimPrefix(cursor, "c"))
		}
		size, _ := strconv.Atoi(pageSize)
		end := min(start+size, len(ids))
		var page []map[string]any
		for _, id := range ids[start:end] {
			page = append(page, map[string]any{"id": id, "display_name": id, "evaluator_id": testEvaluatorID})
		}
		var next any
		if end < len(ids) {
			next = fmt.Sprintf("c%d", end)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{"items": page, "next_cursor": next})
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		cmd := newEvaluatorRuleListCmd()
		_ = cmd.Flags().Set("limit", "3")
		_ = cmd.Flags().Set("cursor", "c1")
		if err := runTestCommand(t, cmd, nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	if !slices.Equal(cursors, []string{"c1"}) || !slices.Equal(pageSizes, []string{"3"}) {
		t.Errorf("expected one request at cursor c1 for 3 rules, got cursors %v sizes %v", cursors, pageSizes)
	}
	var got struct {
		Rules   []map[string]any `json:"rules"`
		Cursors map[string]any   `json:"cursors"`
	}
	if err := json.Unmarshal([]byte(out), &got); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if len(got.Rules) != 3 || got.Rules[0]["rule_id"] != testRuleID2 || got.Cursors["next"] != nil {
		t.Errorf("expected the last 3 rules and no next cursor, got %v", got)
	}

	cursors, pageSizes = nil, nil
	out = captureStdout(t, func() {
		if err := runTestCommand(t, newEvaluatorRuleListCmd(), nil); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	var all []map[string]any
	if err := json.Unmarshal([]byte(out), &all); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if len(all) != len(ids) || !slices.Equal(pageSizes, []string{"100"}) {
		t.Errorf("expected every rule from one 100-rule page, got %d rules with page sizes %v", len(all), pageSizes)
	}
}

// writeRulePage writes a one-page GET /api/v2/runs/rules response.
func writeRulePage(w http.ResponseWriter, items any) error {
	return json.NewEncoder(w).Encode(map[string]any{"items": items, "next_cursor": nil})
}

func TestEvaluatorRuleGet_ShowsCodeEvaluator(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/runs/rules" && r.URL.Query().Get("id") == testRuleID {
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{{
				"id": testRuleID, "display_name": "accuracy", "dataset_id": "ds-1",
				"evaluator_id": testEvaluatorID, "evaluator_name": "accuracy_v2",
				"code_evaluators": []map[string]any{{"code": "def perform_eval(run, example):\n  return {}", "language": "PYTHON"}},
			}})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	out := captureStdout(t, func() {
		if err := runTestCommand(t, newEvaluatorRuleGetCmd(), []string{testRuleID}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if result["type"] != "code" || result["language"] != "python" || result["code"] == "" {
		t.Errorf("expected code evaluator detail, got %v", result)
	}
	if result["evaluator_name"] != "accuracy_v2" {
		t.Errorf("expected evaluator_name=accuracy_v2, got %v", result["evaluator_name"])
	}
}

func TestEvaluatorRuleDelete_NameNeedsTarget(t *testing.T) {
	var calls int
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		calls++
		http.Error(w, "unexpected", http.StatusInternalServerError)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	cmd := newEvaluatorRuleDeleteCmd()
	_ = cmd.Flags().Set("yes", "true")
	err := runTestCommand(t, cmd, []string{"accuracy"})
	if err == nil || !strings.Contains(err.Error(), "--project") {
		t.Fatalf("expected error asking for a target, got %v", err)
	}
	if calls != 0 {
		t.Errorf("expected no API calls, got %d", calls)
	}
}

func TestEvaluatorRuleDelete_RefusesAmbiguousName(t *testing.T) {
	var sawDelete bool
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			sawDelete = true
		}
		if r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{
				{"id": testRuleID, "display_name": "accuracy", "project_id": testSessionID, "evaluator_id": testEvaluatorID},
				{"id": testRuleID2, "display_name": "accuracy", "project_id": testSessionID, "evaluator_id": testEvaluatorID2},
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	cmd := newEvaluatorRuleDeleteCmd()
	_ = cmd.Flags().Set("project-id", testSessionID)
	_ = cmd.Flags().Set("yes", "true")
	err := runTestCommand(t, cmd, []string{"accuracy"})
	if err == nil || !strings.Contains(err.Error(), testRuleID) || !strings.Contains(err.Error(), testRuleID2) {
		t.Fatalf("expected error listing both rule IDs, got %v", err)
	}
	if sawDelete {
		t.Error("expected no DELETE request")
	}
}

func TestEvaluatorRuleDelete_DeletesSingleMatch(t *testing.T) {
	var deletedPath string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet:
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{
				{"id": testRuleID, "display_name": "accuracy", "project_id": testSessionID, "evaluator_id": testEvaluatorID},
				{"id": testRuleID2, "display_name": "toxicity", "project_id": testSessionID},
			})
		case r.Method == http.MethodDelete:
			deletedPath = r.URL.Path
			w.WriteHeader(http.StatusNoContent)
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	out := captureStdout(t, func() {
		cmd := newEvaluatorRuleDeleteCmd()
		_ = cmd.Flags().Set("project-id", testSessionID)
		_ = cmd.Flags().Set("yes", "true")
		if err := runTestCommand(t, cmd, []string{"accuracy"}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if deletedPath != "/api/v1/runs/rules/"+testRuleID {
		t.Errorf("expected DELETE of rule %s, got %q", testRuleID, deletedPath)
	}
	var result map[string]any
	if err := json.Unmarshal([]byte(out), &result); err != nil {
		t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
	}
	if result["rule_id"] != testRuleID || result["evaluator_id"] != testEvaluatorID {
		t.Errorf("unexpected output: %v", result)
	}
}

func TestEvaluatorRuleList_HidesNonEvaluatorRulesUnlessAll(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/api/v2/runs/rules" {
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{
				{"id": testRuleID, "display_name": "test", "evaluator_id": testEvaluatorID, "evaluator_name": "ready_for_task_grade"},
				{"id": testRuleID2, "display_name": "legacy", "code_evaluators": []map[string]any{{"code": "x", "language": "python"}}},
				{"id": "1a2b3c4d-0000-4000-8000-000000000001", "display_name": "notify", "webhooks": []map[string]any{{"url": "https://example.com"}}},
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	listNames := func(all bool) []string {
		out := captureStdout(t, func() {
			cmd := newEvaluatorRuleListCmd()
			if all {
				_ = cmd.Flags().Set("all", "true")
			}
			if err := runTestCommand(t, cmd, nil); err != nil {
				t.Errorf("unexpected error: %v", err)
			}
		})
		var result []map[string]any
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("failed to parse output JSON: %v\noutput: %s", err, out)
		}
		var names []string
		for _, r := range result {
			names = append(names, r["rule_name"].(string))
		}
		return names
	}

	if got := strings.Join(listNames(false), ","); got != "test,legacy" {
		t.Errorf("expected evaluator and legacy inline rules only, got %s", got)
	}
	if got := strings.Join(listNames(true), ","); got != "test,legacy,notify" {
		t.Errorf("expected --all to include the webhook rule, got %s", got)
	}
}

func TestEvaluatorRuleDelete_RefusesNonEvaluatorRule(t *testing.T) {
	var sawDelete bool
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodDelete {
			sawDelete = true
		}
		if r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet {
			w.Header().Set("Content-Type", "application/json")
			_ = writeRulePage(w, []map[string]any{
				{"id": testRuleID, "display_name": "notify", "project_id": testSessionID, "webhooks": []map[string]any{{"url": "https://example.com"}}},
			})
			return
		}
		http.Error(w, "not found", http.StatusNotFound)
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	cmd := newEvaluatorRuleDeleteCmd()
	_ = cmd.Flags().Set("yes", "true")
	err := runTestCommand(t, cmd, []string{testRuleID})
	if err == nil || !strings.Contains(err.Error(), "not an evaluator rule") {
		t.Fatalf("expected refusal for a webhook rule, got %v", err)
	}
	if sawDelete {
		t.Error("expected no DELETE request")
	}
}

func TestEvaluatorUploadReplaceMatchesRenamedEvaluator(t *testing.T) {
	evaluatorFile := t.TempDir() + "/eval.py"
	if err := os.WriteFile(evaluatorFile, []byte("def grade(run, example):\n    return {\"score\": 1}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var patchedPath string
	var patchBody map[string]any
	var sawPost bool
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet:
			_ = writeRulePage(w, []map[string]any{{
				"id": testRuleID, "display_name": "test", "project_id": testSessionID,
				"evaluator_id": testEvaluatorID, "evaluator_name": "ready_for_task_grade",
			}})
		case r.URL.Path == "/api/v1/runs/rules" && r.Method == http.MethodPost:
			sawPost = true
			http.Error(w, "should not create", http.StatusInternalServerError)
		case r.URL.Path == "/api/v1/platform/evaluators/"+testEvaluatorID && r.Method == http.MethodGet:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": testEvaluatorID, "type": "code", "code_evaluator": map[string]any{"code": "def perform_eval(run, example): ..."}})
		case r.Method == http.MethodPatch:
			patchedPath = r.URL.Path
			_ = json.NewDecoder(r.Body).Decode(&patchBody)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": testRuleID, "evaluator_id": testEvaluatorID})
		default:
			http.Error(w, "not found", http.StatusNotFound)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	captureStdout(t, func() {
		cmd := newEvaluatorUploadCmd()
		_ = cmd.Flags().Set("name", "ready_for_task_grade")
		_ = cmd.Flags().Set("function", "grade")
		_ = cmd.Flags().Set("project-id", testSessionID)
		_ = cmd.Flags().Set("replace", "true")
		_ = cmd.Flags().Set("yes", "true")
		if err := runTestCommand(t, cmd, []string{evaluatorFile}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})

	if sawPost {
		t.Fatal("expected --replace to update the existing evaluator, not create a second one")
	}
	if patchedPath != "/api/v1/runs/rules/"+testRuleID {
		t.Fatalf("expected PATCH of rule %s, got %q", testRuleID, patchedPath)
	}
	if patchBody["display_name"] != "test" {
		t.Errorf("expected the rule to keep its name, got %v", patchBody["display_name"])
	}
}

func TestEvaluatorUploadWithoutReplaceRefusesRenamedEvaluator(t *testing.T) {
	evaluatorFile := t.TempDir() + "/eval.py"
	if err := os.WriteFile(evaluatorFile, []byte("def grade(run, example):\n    return {\"score\": 1}\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	var sawWrite bool
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			sawWrite = true
		}
		w.Header().Set("Content-Type", "application/json")
		_ = writeRulePage(w, []map[string]any{{
			"id": testRuleID, "display_name": "test", "project_id": testSessionID,
			"evaluator_id": testEvaluatorID, "evaluator_name": "ready_for_task_grade",
		}})
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()

	cmd := newEvaluatorUploadCmd()
	_ = cmd.Flags().Set("name", "ready_for_task_grade")
	_ = cmd.Flags().Set("function", "grade")
	_ = cmd.Flags().Set("project-id", testSessionID)
	err := runTestCommand(t, cmd, []string{evaluatorFile})
	if err == nil || !strings.Contains(err.Error(), "already exists") {
		t.Fatalf("expected already-exists error, got %v", err)
	}
	if sawWrite {
		t.Error("expected no create or update request")
	}
}

// ---------- whole-file upload ----------

// wholeFileServer fakes the evaluators API and the rules endpoint for uploads
// without --function. An advanced evaluator's build reports buildPolls BUILDING
// responses before it settles on finalStatus.
type wholeFileServer struct {
	mu            sync.Mutex
	evaluator     map[string]any // stored evaluator, nil until created
	createBody    map[string]any
	updateBody    map[string]any
	ruleBody      map[string]any
	rulePatch     map[string]any
	rules         []map[string]any
	ruleCreateErr bool
	buildPolls    int
	finalStatus   string
	buildError    string
	gets          int
	deletes       int
	updates       int
	noRebuild     bool // updates keep the stored build status, as when nothing changed
}

func (s *wholeFileServer) codeEvaluator() map[string]any {
	ce, _ := s.evaluator["code_evaluator"].(map[string]any)
	return ce
}

func (s *wholeFileServer) store(body map[string]any) {
	ce, _ := body["code_evaluator"].(map[string]any)
	stored := map[string]any{}
	if s.evaluator != nil {
		stored = s.codeEvaluator()
	}
	for k, v := range ce {
		stored[k] = v
	}
	if stored["advanced_features_enabled"] == true && !s.noRebuild {
		stored["evaluator_build_status"] = "ENQUEUED"
	}
	s.evaluator = map[string]any{"id": testEvaluatorID, "name": "accuracy", "type": "code", "code_evaluator": stored}
}

func (s *wholeFileServer) handler(t *testing.T) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		s.mu.Lock()
		defer s.mu.Unlock()
		w.Header().Set("Content-Type", "application/json")
		evaluatorPath := "/api/v1/platform/evaluators/" + testEvaluatorID
		decode := func() map[string]any {
			var body map[string]any
			if err := json.NewDecoder(r.Body).Decode(&body); err != nil {
				t.Errorf("decoding %s %s: %v", r.Method, r.URL.Path, err)
			}
			return body
		}
		switch {
		case r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet:
			_ = writeRulePage(w, s.rules)
		case r.URL.Path == "/api/v1/platform/evaluators" && r.Method == http.MethodPost:
			s.createBody = decode()
			s.store(s.createBody)
			_ = json.NewEncoder(w).Encode(map[string]any{"evaluator": s.evaluator})
		case r.URL.Path == evaluatorPath && r.Method == http.MethodPatch:
			s.updates++
			if s.evaluator != nil && buildInProgress(langsmith.OnlineCodeEvaluatorEvaluatorBuildStatus(fmt.Sprint(s.codeEvaluator()["evaluator_build_status"]))) {
				w.WriteHeader(http.StatusConflict)
				_ = json.NewEncoder(w).Encode(map[string]any{"detail": "Cannot update evaluator: a build is in progress; retry when status is READY or FAILED"})
				return
			}
			s.updateBody = decode()
			s.store(s.updateBody)
			_ = json.NewEncoder(w).Encode(map[string]any{"evaluator": s.evaluator})
		case r.URL.Path == evaluatorPath && r.Method == http.MethodGet:
			s.gets++
			ce := s.codeEvaluator()
			if buildInProgress(langsmith.OnlineCodeEvaluatorEvaluatorBuildStatus(fmt.Sprint(ce["evaluator_build_status"]))) {
				if s.buildPolls > 0 {
					s.buildPolls--
					ce["evaluator_build_status"] = "BUILDING"
				} else {
					ce["evaluator_build_status"] = s.finalStatus
					ce["evaluator_build_error"] = s.buildError
				}
			}
			_ = json.NewEncoder(w).Encode(s.evaluator)
		case r.URL.Path == evaluatorPath && r.Method == http.MethodDelete:
			s.deletes++
			w.WriteHeader(http.StatusNoContent)
		case r.URL.Path == "/api/v1/runs/rules" && r.Method == http.MethodPost:
			s.ruleBody = decode()
			if s.ruleCreateErr {
				w.WriteHeader(http.StatusUnprocessableEntity)
				_ = json.NewEncoder(w).Encode(map[string]any{"detail": "code_evaluator_input=trajectory requires group_by=trajectory_id"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"id": testRuleID, "evaluator_id": testEvaluatorID})
		case r.URL.Path == "/api/v1/runs/rules/"+testRuleID && r.Method == http.MethodPatch:
			s.rulePatch = decode()
			_ = json.NewEncoder(w).Encode(map[string]any{"id": testRuleID})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
			http.Error(w, "not found", http.StatusNotFound)
		}
	}
}

const wholeFileSource = "import json\n\nTHRESHOLD = 0.5\n\n\ndef _score(text):\n    return len(text) > 3\n\n\ndef perform_eval(run, example):\n    return {\"score\": _score(json.dumps(run))}\n"

// runWholeFileUpload uploads wholeFileSource without --function and returns
// the parsed output and the command's error.
func runWholeFileUpload(t *testing.T, srv *wholeFileServer, flags map[string]string) (map[string]any, error) {
	t.Helper()
	return runWholeFileUploadSource(t, srv, wholeFileSource, flags)
}

// runWholeFileUploadSource is runWholeFileUpload with the given file contents.
func runWholeFileUploadSource(t *testing.T, srv *wholeFileServer, source string, flags map[string]string) (map[string]any, error) {
	t.Helper()
	prevInterval := evaluatorBuildPollInterval
	evaluatorBuildPollInterval = time.Millisecond
	t.Cleanup(func() { evaluatorBuildPollInterval = prevInterval })

	dir := t.TempDir()
	evaluatorFile := dir + "/eval.py"
	if err := os.WriteFile(evaluatorFile, []byte(source), 0o600); err != nil {
		t.Fatal(err)
	}
	ts := newTestServer(t, srv.handler(t))
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	cmd := newEvaluatorUploadCmd()
	for k, v := range map[string]string{"name": "accuracy", "project-id": testSessionID} {
		_ = cmd.Flags().Set(k, v)
	}
	for k, v := range flags {
		if k == "dependencies" {
			path := dir + "/requirements.txt"
			if err := os.WriteFile(path, []byte(v), 0o600); err != nil {
				t.Fatal(err)
			}
			v = path
		}
		items := []string{v}
		if k == "workspace-secret" {
			items = strings.Split(v, ",")
		}
		for _, item := range items {
			if err := cmd.Flags().Set(k, item); err != nil {
				t.Fatalf("setting --%s: %v", k, err)
			}
		}
	}
	var runErr error
	out := captureStdout(t, func() { runErr = runTestCommand(t, cmd, []string{evaluatorFile}) })
	var result map[string]any
	if runErr == nil {
		if err := json.Unmarshal([]byte(out), &result); err != nil {
			t.Fatalf("parsing output %q: %v", out, err)
		}
	}
	return result, runErr
}

func TestEvaluatorUploadWholeFileSendsFileUnchanged(t *testing.T) {
	srv := &wholeFileServer{}
	result, err := runWholeFileUpload(t, srv, map[string]string{"sampling-rate": "0.5", "trace-filter": "eq(is_root, true)"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv.createBody["name"] != "accuracy" || srv.createBody["type"] != "code" {
		t.Errorf("evaluator create body = %#v", srv.createBody)
	}
	ce, _ := srv.createBody["code_evaluator"].(map[string]any)
	want := map[string]any{"code": wholeFileSource, "language": "python"}
	if !reflect.DeepEqual(ce, want) {
		t.Errorf("code_evaluator = %#v, want exactly %#v (no advanced_features_enabled without --advanced)", ce, want)
	}
	wantRule := map[string]any{
		"display_name": "accuracy", "sampling_rate": 0.5, "is_enabled": true, "include_extended_stats": false,
		"evaluator_id": testEvaluatorID, "evaluator_version": float64(3), "session_id": testSessionID,
		"trace_filter": "eq(is_root, true)",
	}
	if !reflect.DeepEqual(srv.ruleBody, wantRule) {
		t.Errorf("rule body = %#v, want %#v", srv.ruleBody, wantRule)
	}
	if srv.gets != 0 {
		t.Errorf("a non-advanced upload should not poll the build, got %d GETs", srv.gets)
	}
	if result["rule_id"] != testRuleID || result["evaluator_id"] != testEvaluatorID || result["advanced"] != false {
		t.Errorf("unexpected output %#v", result)
	}
}

func TestEvaluatorUploadWholeFileAdvancedWaitsForBuild(t *testing.T) {
	srv := &wholeFileServer{buildPolls: 2, finalStatus: "READY"}
	result, err := runWholeFileUpload(t, srv, map[string]string{
		"advanced": "true", "dependencies": "requests==2.32.3\n", "workspace-secret": "OPENAI_KEY,SERP_KEY", "input": "trajectory",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	ce, _ := srv.createBody["code_evaluator"].(map[string]any)
	want := map[string]any{
		"code": wholeFileSource, "language": "python", "advanced_features_enabled": true,
		"dependencies": "requests==2.32.3\n", "workspace_secrets_keys": []any{"OPENAI_KEY", "SERP_KEY"},
		"code_evaluator_input": "trajectory",
	}
	if !reflect.DeepEqual(ce, want) {
		t.Errorf("code_evaluator = %#v, want %#v", ce, want)
	}
	if srv.ruleBody["group_by"] != "trajectory_id" || srv.ruleBody["evaluator_id"] != testEvaluatorID {
		t.Errorf("rule body = %#v, want group_by trajectory_id and the evaluator ID", srv.ruleBody)
	}
	if srv.gets != 3 {
		t.Errorf("expected 3 build polls, got %d", srv.gets)
	}
	if result["build_status"] != "READY" || result["advanced"] != true {
		t.Errorf("unexpected output %#v", result)
	}
}

func TestEvaluatorUploadWholeFileAdvancedBuildFailed(t *testing.T) {
	srv := &wholeFileServer{finalStatus: "FAILED", buildError: "pip install failed: No matching distribution found for nosuchpkg"}
	_, err := runWholeFileUpload(t, srv, map[string]string{"advanced": "true", "dependencies": "nosuchpkg\n"})
	if err == nil || !strings.Contains(err.Error(), "build failed: pip install failed: No matching distribution") {
		t.Fatalf("expected the build error, got %v", err)
	}
	if srv.deletes != 0 {
		t.Error("a failed build should keep the evaluator so it can be fixed with --replace")
	}
}

func TestEvaluatorUploadWholeFileNoWaitSkipsBuild(t *testing.T) {
	srv := &wholeFileServer{finalStatus: "FAILED"}
	result, err := runWholeFileUpload(t, srv, map[string]string{"advanced": "true", "no-wait": "true"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv.gets != 0 {
		t.Errorf("--no-wait should not poll the build, got %d GETs", srv.gets)
	}
	if result["build_status"] != "ENQUEUED" {
		t.Errorf("build_status = %v, want ENQUEUED", result["build_status"])
	}
}

func TestEvaluatorUploadWholeFileDeletesEvaluatorWhenRuleFails(t *testing.T) {
	srv := &wholeFileServer{ruleCreateErr: true}
	_, err := runWholeFileUpload(t, srv, nil)
	if err == nil || !strings.Contains(err.Error(), "requires group_by=trajectory_id") {
		t.Fatalf("expected the server's rule error, got %v", err)
	}
	if srv.deletes != 1 {
		t.Errorf("expected the new evaluator to be deleted once, got %d deletes", srv.deletes)
	}
	if !strings.Contains(err.Error(), "evaluator "+testEvaluatorID+" was deleted") {
		t.Errorf("expected the error to say the evaluator was deleted, got %v", err)
	}
}

func TestEvaluatorUploadWholeFileReplace(t *testing.T) {
	existingRule := map[string]any{
		"id": testRuleID, "display_name": "accuracy", "project_id": testSessionID,
		"evaluator_id": testEvaluatorID, "evaluator_name": "accuracy", "sampling_rate": 0.2,
		"filter": "eq(is_root, true)", "is_enabled": true,
	}
	t.Run("keeps mode", func(t *testing.T) {
		srv := &wholeFileServer{rules: []map[string]any{existingRule}}
		srv.store(map[string]any{"code_evaluator": map[string]any{"code": "old", "language": "python"}})
		if _, err := runWholeFileUpload(t, srv, map[string]string{"replace": "true", "yes": "true"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		ce, _ := srv.updateBody["code_evaluator"].(map[string]any)
		if want := (map[string]any{"code": wholeFileSource, "language": "python"}); !reflect.DeepEqual(ce, want) {
			t.Errorf("update code_evaluator = %#v, want %#v", ce, want)
		}
		if srv.createBody != nil || srv.ruleBody != nil {
			t.Error("--replace should update, not create")
		}
		if srv.rulePatch["sampling_rate"] != 0.2 || srv.rulePatch["filter"] != "eq(is_root, true)" || srv.rulePatch["session_id"] != testSessionID {
			t.Errorf("rule PATCH should keep the rule's settings, got %#v", srv.rulePatch)
		}
		for _, key := range []string{"code_evaluators", "evaluator_id"} {
			if _, ok := srv.rulePatch[key]; ok {
				t.Errorf("rule PATCH should not send %s", key)
			}
		}
	})
	t.Run("advanced waits for the running build", func(t *testing.T) {
		srv := &wholeFileServer{rules: []map[string]any{existingRule}, buildPolls: 1, finalStatus: "READY"}
		srv.store(map[string]any{"code_evaluator": map[string]any{"code": "old", "language": "python", "advanced_features_enabled": true}})
		result, err := runWholeFileUpload(t, srv, map[string]string{"replace": "true", "yes": "true"})
		if err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if srv.updates < 2 || srv.updateBody == nil {
			t.Errorf("expected the update to be retried after the build, got %d attempts", srv.updates)
		}
		if _, ok := srv.updateBody["code_evaluator"].(map[string]any)["advanced_features_enabled"]; ok {
			t.Error("replace without --advanced should keep the stored mode")
		}
		if result["advanced"] != true || result["build_status"] != "READY" {
			t.Errorf("unexpected output %#v", result)
		}
	})
	t.Run("unchanged advanced evaluator reports its failed build", func(t *testing.T) {
		srv := &wholeFileServer{rules: []map[string]any{existingRule}, noRebuild: true}
		srv.store(map[string]any{"code_evaluator": map[string]any{
			"code": wholeFileSource, "language": "python", "advanced_features_enabled": true,
			"evaluator_build_status": "FAILED", "evaluator_build_error": "pip install failed",
		}})
		_, err := runWholeFileUpload(t, srv, map[string]string{"replace": "true", "yes": "true"})
		if err == nil || !strings.Contains(err.Error(), "build failed: pip install failed") {
			t.Fatalf("expected the stored build error, got %v", err)
		}
	})
	t.Run("refusals leave the evaluator unchanged", func(t *testing.T) {
		withAlerts := maps.Clone(existingRule)
		withAlerts["alerts"] = []any{map[string]any{"routing_key": "k", "type": "pagerduty"}}
		for name, tc := range map[string]struct {
			rule  map[string]any
			flags map[string]string
			want  string
		}{
			"alerts":          {withAlerts, nil, "alert"},
			"to trajectory":   {existingRule, map[string]string{"input": "trajectory"}, "trajectory input"},
			"from trajectory": {map[string]any{"id": testRuleID, "display_name": "accuracy", "project_id": testSessionID, "evaluator_id": testEvaluatorID, "evaluator_name": "accuracy", "group_by": "trajectory_id"}, map[string]string{"input": "thread"}, "trajectory input"},
		} {
			t.Run(name, func(t *testing.T) {
				srv := &wholeFileServer{rules: []map[string]any{tc.rule}}
				srv.store(map[string]any{"code_evaluator": map[string]any{"code": "old", "language": "python"}})
				flags := map[string]string{"replace": "true", "yes": "true"}
				maps.Copy(flags, tc.flags)
				_, err := runWholeFileUpload(t, srv, flags)
				if err == nil || !strings.Contains(err.Error(), tc.want) {
					t.Fatalf("expected a refusal mentioning %q, got %v", tc.want, err)
				}
				if srv.updates != 0 || srv.rulePatch != nil {
					t.Errorf("expected no writes, got %d evaluator updates and rule PATCH %#v", srv.updates, srv.rulePatch)
				}
			})
		}
	})
	t.Run("--advanced switches mode", func(t *testing.T) {
		srv := &wholeFileServer{rules: []map[string]any{existingRule}, finalStatus: "READY"}
		srv.store(map[string]any{"code_evaluator": map[string]any{"code": "old", "language": "python"}})
		if _, err := runWholeFileUpload(t, srv, map[string]string{"replace": "true", "yes": "true", "advanced": "true"}); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
		if srv.updateBody["code_evaluator"].(map[string]any)["advanced_features_enabled"] != true {
			t.Errorf("update body = %#v, want advanced_features_enabled", srv.updateBody)
		}
	})
}

func TestEvaluatorUploadFunctionReplaceRefusesAdvancedEvaluator(t *testing.T) {
	srv := &wholeFileServer{rules: []map[string]any{{
		"id": testRuleID, "display_name": "accuracy", "project_id": testSessionID,
		"evaluator_id": testEvaluatorID, "evaluator_name": "accuracy",
	}}}
	srv.store(map[string]any{"code_evaluator": map[string]any{"code": "x", "language": "python", "advanced_features_enabled": true}})
	_, err := runWholeFileUpload(t, srv, map[string]string{"function": "perform_eval", "replace": "true", "yes": "true"})
	if err == nil || !strings.Contains(err.Error(), "advanced evaluator") {
		t.Fatalf("expected refusal, got %v", err)
	}
	if srv.rulePatch != nil || srv.updateBody != nil {
		t.Error("expected no update")
	}
}

// The deprecated --function path sends the extracted, renamed function inline
// to the rules endpoint, exactly as before whole-file uploads.
func TestEvaluatorUploadFunctionSendsRuleWithInlineCode(t *testing.T) {
	evaluatorFile := t.TempDir() + "/eval.py"
	if err := os.WriteFile(evaluatorFile, []byte(wholeFileSource+"\n\ndef grade(run, example):\n    return {\"score\": 1}\n"), 0o600); err != nil {
		t.Fatal(err)
	}
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.URL.Path == "/api/v2/runs/rules" && r.Method == http.MethodGet:
			_ = writeRulePage(w, []map[string]any{})
		case r.URL.Path == "/api/v1/runs/rules" && r.Method == http.MethodPost:
			_ = json.NewDecoder(r.Body).Decode(&body)
			_ = json.NewEncoder(w).Encode(map[string]any{"id": testRuleID, "evaluator_id": testEvaluatorID})
		default:
			t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})
	cleanup := setupTestEnv(t, ts.URL)
	defer cleanup()
	flagOutputFormat = "json"

	out := captureStdout(t, func() {
		cmd := newEvaluatorUploadCmd()
		for k, v := range map[string]string{"name": "accuracy", "function": "grade", "project-id": testSessionID} {
			_ = cmd.Flags().Set(k, v)
		}
		if err := runTestCommand(t, cmd, []string{evaluatorFile}); err != nil {
			t.Errorf("unexpected error: %v", err)
		}
	})
	want := map[string]any{
		"display_name": "accuracy", "sampling_rate": 1.0, "is_enabled": true, "include_extended_stats": false,
		"session_id": testSessionID,
		"code_evaluators": []any{map[string]any{
			"code": "def perform_eval(run, example):\n    return {\"score\": 1}\n", "language": "python",
		}},
	}
	if !reflect.DeepEqual(body, want) {
		t.Errorf("rule body = %#v, want %#v", body, want)
	}
	var result map[string]any
	_ = json.Unmarshal([]byte(out), &result)
	wantOut := map[string]any{"status": "uploaded", "rule_id": testRuleID, "evaluator_id": testEvaluatorID, "name": "accuracy", "target": "project"}
	if !reflect.DeepEqual(result, wantOut) {
		t.Errorf("output = %#v, want %#v", result, wantOut)
	}
}

// ==================== helper ====================

func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

func TestEvaluatorUploadRefusesHardcodedSecrets(t *testing.T) {
	key := fakeKey("sk-proj-", 32)
	source := "import json\n\nOPENAI_KEY = \"" + key + "\"\n\ndef perform_eval(run, example):\n    return {}\n"
	for name, tc := range map[string]struct {
		source string
		flags  map[string]string
	}{
		"in the code":              {source, nil},
		"in the dependencies file": {wholeFileSource, map[string]string{"advanced": "true", "dependencies": "--index-url https://user:hunter2@pypi.example.com/simple\nrequests\n"}},
		"with --function": {
			"def perform_eval(run, example):\n    client = OpenAI(api_key=\"" + key + "\")\n    return {}\n",
			map[string]string{"function": "perform_eval"},
		},
	} {
		t.Run(name, func(t *testing.T) {
			srv := &wholeFileServer{}
			var err error
			out := captureStdout(t, func() { _, err = runWholeFileUploadSource(t, srv, tc.source, tc.flags) })
			if err == nil || !strings.Contains(err.Error(), "hardcoded credential") {
				t.Fatalf("expected a credential refusal, got %v", err)
			}
			if strings.Contains(err.Error()+out, key) || strings.Contains(err.Error(), "hunter2") {
				t.Errorf("the error should not repeat the credential: %v", err)
			}
			if srv.createBody != nil || srv.ruleBody != nil {
				t.Error("expected nothing to be uploaded")
			}
		})
	}
}

func TestEvaluatorUploadIgnoreCredentialCheck(t *testing.T) {
	source := "def perform_eval(run, example):\n    key = \"" + fakeKey("sk-proj-", 32) + "\"\n    return {}\n"
	for name, flags := range map[string]map[string]string{
		"whole file": {"ignore-credential-check": "true"},
		"--function": {"ignore-credential-check": "true", "function": "perform_eval"},
	} {
		t.Run(name, func(t *testing.T) {
			srv := &wholeFileServer{}
			if _, err := runWholeFileUploadSource(t, srv, source, flags); err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if srv.createBody == nil && srv.ruleBody == nil {
				t.Error("expected the upload to go ahead")
			}
		})
	}
}

func TestEvaluatorUploadLocalImports(t *testing.T) {
	const source = "from .helpers import score\n\ndef evaluate(ctx):\n    return score(ctx)\n"
	t.Run("advanced upload is refused", func(t *testing.T) {
		srv := &wholeFileServer{}
		_, err := runWholeFileUploadSource(t, srv, source, map[string]string{"advanced": "true"})
		if !errors.Is(err, errLocalImports) || !strings.Contains(err.Error(), `".helpers"`) {
			t.Fatalf("expected a local import refusal naming .helpers, got %v", err)
		}
		if srv.createBody != nil {
			t.Error("expected nothing to be uploaded")
		}
	})
	t.Run("replace of an advanced evaluator is refused", func(t *testing.T) {
		srv := &wholeFileServer{rules: []map[string]any{{
			"id": testRuleID, "display_name": "accuracy", "project_id": testSessionID,
			"evaluator_id": testEvaluatorID, "evaluator_name": "accuracy",
		}}}
		srv.store(map[string]any{"code_evaluator": map[string]any{"code": "old", "language": "python", "advanced_features_enabled": true}})
		_, err := runWholeFileUploadSource(t, srv, source, map[string]string{"replace": "true", "yes": "true"})
		if !errors.Is(err, errLocalImports) {
			t.Fatalf("expected a local import refusal, got %v", err)
		}
		if srv.updates != 0 {
			t.Error("expected no evaluator update")
		}
	})
	t.Run("legacy upload is left to the server", func(t *testing.T) {
		srv := &wholeFileServer{}
		if _, err := runWholeFileUploadSource(t, srv, source, nil); err != nil {
			t.Fatalf("unexpected error: %v", err)
		}
	})
}

func TestCheckLocalImports(t *testing.T) {
	dir := t.TempDir()
	for _, p := range []string{"helpers.py", "pkg/__init__.py", "alone/json.py", "data/notes.txt"} {
		if err := os.MkdirAll(filepath.Dir(filepath.Join(dir, p)), 0o700); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(filepath.Join(dir, p), nil, 0o600); err != nil {
			t.Fatal(err)
		}
	}
	for name, tc := range map[string]struct {
		file, code string
		want       []string // modules named in the error; none means accepted
	}{
		"python packages":           {"eval.py", "import json, re\nfrom rapidfuzz import fuzz\nimport numpy as np\n", nil},
		"python relative":           {"eval.py", "from . import utils\nfrom .lib.x import y\n", []string{`"."`, `".lib.x"`}},
		"python sibling module":     {"eval.py", "import json, helpers as h\nfrom pkg.sub import thing\n", []string{`"helpers"`, `"pkg.sub"`}},
		"python indented import":    {"eval.py", "def evaluate(ctx):\n    import helpers\n    return {}\n", []string{`"helpers"`}},
		"js packages":               {"eval.js", "import fuzz from 'fuzzball';\nconst _ = require(\"lodash\");\n", nil},
		"js relative":               {"eval.ts", "import { score } from './score.js';\nexport * from \"../shared\";\nconst cfg = require('./cfg.json');\nconst m = await import(`./lazy`);\n", []string{`"./score.js"`, `"../shared"`, `"./cfg.json"`, `"./lazy"`}},
		"js side-effect import":     {"eval.mjs", "import './setup.js';\n", []string{`"./setup.js"`}},
		"js commented imports":      {"eval.js", "// import { score } from './score.js';\n/* const x = require('./x');\n*/\nconst u = 'https://example.com';\n", nil},
		"python docstring example":  {"eval.py", "\"\"\"Usage:\n\nfrom .helpers import score\nimport helpers\n\"\"\"\nimport json\n", nil},
		"python after docstring":    {"eval.py", "'''Doc.'''\nfrom .helpers import score\n", []string{`".helpers"`}},
		"python file named json":    {"alone/json.py", "import json\n", nil},
		"python non-code folder":    {"eval.py", "import data\n", nil},
		"js Buffer.from":            {"eval.js", "const b = Buffer.from('./expected');\nconst a = Array.from(\"./x\");\n", nil},
		"js regex with quote":       {"eval.js", "const q = /['\"]/g;\nimport { score } from './score.js';\n", []string{`"./score.js"`}},
		"js regex with comment":     {"eval.js", "const ok = /[/*]/.test(s);\nconst h = require('./helpers');\n", []string{`"./helpers"`}},
		"js division not regex":     {"eval.js", "const r = a / b / c;\nconst h = require('./helpers');\n", []string{`"./helpers"`}},
		"js import in ${}":          {"eval.js", "const msg = `score: ${require('./score').value}`;\n", []string{`"./score"`}},
		"js nested template":        {"eval.js", "const s = `a ${`b ${import('./deep')}`} c`;\nconst x = require('lodash');\n", []string{`"./deep"`}},
		"js spread require":         {"eval.js", "const cfg = { ...require('./defaults'), x: 1 };\n", []string{`"./defaults"`}},
		"js separate expressions":   {"eval.js", "const s = `${from} ${'./not-a-module'}`;\n", nil},
		"js /* in a string":         {"eval.js", "const glob = 'src/*.js';\nimport { score } from './score.js';\nconst end = '*/';\n", []string{`"./score.js"`}},
		"js import text in string":  {"eval.js", "const hint = \"import x from './x'\";\nconst r = require(`lodash`);\n", nil},
		"js import in template":     {"eval.js", "const cfg = await import(`./cfg.js`);\n", []string{`"./cfg.js"`}},
		"python prompt with quotes": {"eval.py", "PROMPT = 'Wrap code in \"\"\" fences'\nfrom .helpers import score\n", []string{`".helpers"`}},
		"python escaped quote":      {"eval.py", "\"\"\"Say \\\"\"\" hi.\"\"\"\nfrom .helpers import score\n", []string{`".helpers"`}},
		"python import in string":   {"eval.py", "EXAMPLE = \"from .helpers import score\"\nimport json\n", nil},
		"python import in comment":  {"eval.py", "# from .helpers import score\nimport json  # import helpers\n", nil},
	} {
		t.Run(name, func(t *testing.T) {
			language, _ := detectLanguage(tc.file)
			err := checkLocalImports(filepath.Join(dir, tc.file), tc.code, language)
			if len(tc.want) == 0 {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil {
				t.Fatalf("expected an error naming %v", tc.want)
			}
			for _, m := range tc.want {
				if !strings.Contains(err.Error(), m) {
					t.Errorf("error %q does not name %s", err, m)
				}
			}
		})
	}
}

// fakeKey builds a credential-shaped test value at run time, so the source
// holds no string that secret scanners read as a real key.
func fakeKey(prefix string, n int) string {
	const chars = "ABCDEFGHJKLMNPQRSTUVWXYZ23456789"
	b := []byte(prefix)
	for i := range n {
		b = append(b, chars[(i*7+3)%len(chars)])
	}
	return string(b)
}

func TestCheckForSecrets(t *testing.T) {
	for name, tc := range map[string]struct {
		line string
		kind string // empty means accepted
	}{
		"langsmith":        {`KEY = "` + fakeKey("lsv2_pt_", 34) + `"`, "LangSmith API key"},
		"anthropic":        {`key = "` + fakeKey("sk-ant-api03-", 26) + `"`, "Anthropic API key"},
		"openai":           {`openai.api_key = '` + fakeKey("sk-", 30) + `'`, "OpenAI API key"},
		"aws":              {`AWS = "` + fakeKey("AKIA", 16) + `"`, "AWS access key ID"},
		"google":           {`g = "` + fakeKey("AIza", 35) + `"`, "Google API key"},
		"github":           {`t = "` + fakeKey("ghp_", 36) + `"`, "GitHub token"},
		"slack":            {`s = "` + fakeKey("xoxb-", 24) + `"`, "Slack token"},
		"stripe":           {`s = "` + fakeKey("sk_"+"live_", 24) + `"`, "Stripe secret key"},
		"private key":      {"-----BEGIN RSA " + "PRIVATE KEY-----", "private key"},
		"url password":     {`DB = "postgres://admin:s3cret@db.internal:5432/app"`, "password in a URL"},
		"env lookup":       {`key = os.environ["OPENAI_API_KEY"]`, ""},
		"placeholder":      {`# set OPENAI_API_KEY to sk-...`, ""},
		"url without pw":   {`URL = "https://user@example.com/path"`, ""},
		"f-string url":     {`url = f"postgresql://{user}:{password}@{host}/db"`, ""},
		"template url":     {"const u = `https://${user}:${token}@api.example.com`;", ""},
		"%-format url":     {`u = "https://%s:%s@host" % (a, b)`, ""},
		"encoded password": {`DB = "postgres://admin:p%40ss%21word@db.internal/app"`, "password in a URL"},
		"ordinary string":  {`label = "task-classification-for-long-thread-names"`, ""},
	} {
		t.Run(name, func(t *testing.T) {
			err := checkForSecrets("eval.py", "import os\n"+tc.line+"\n", false)
			if tc.kind == "" {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				return
			}
			if err == nil || !strings.Contains(err.Error(), "eval.py:2: looks like a "+tc.kind) {
				t.Fatalf("expected eval.py:2 flagged as a %s, got %v", tc.kind, err)
			}
		})
	}
}

func TestEvaluatorUploadIgnoreImportCheck(t *testing.T) {
	const source = "from .helpers import score\n\ndef evaluate(ctx):\n    return score(ctx)\n"
	srv := &wholeFileServer{finalStatus: "READY"}
	if _, err := runWholeFileUploadSource(t, srv, source, map[string]string{"advanced": "true", "ignore-import-check": "true"}); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if srv.createBody == nil {
		t.Error("expected the upload to go ahead")
	}
}

func TestWithEntryPointHint(t *testing.T) {
	serverErr := errors.New(`400 Bad Request: code must define a function named "perform_eval"`)
	for name, tc := range map[string]struct {
		code, language string
		input          string
		want           string // empty means no hint
	}{
		"thread input": {"def check_thread(run):\n    return {}\n", "python", "thread", "needs the file to define evaluate(ctx); it defines check_thread, so rename the one that scores to evaluate"},
		"python":       {"import json\n\ndef _helper():\n    pass\n\nasync def check_accuracy(run, example):\n    return {}\n", "python", "", "it defines check_accuracy, so rename the one that scores to perform_eval"},
		"javascript":   {"export function grade(run) {}\nconst score = (run) => 1;\nconst x = 1;\n", "javascript", "", "it defines grade, score"},
		"no functions": {"X = 1\n", "python", "", "the file must define perform_eval(run, example)"},
		"other errors": {"def f(): pass\n", "python", "", ""},
	} {
		t.Run(name, func(t *testing.T) {
			in := serverErr
			if tc.want == "" {
				in = errors.New("500 Internal Server Error")
			}
			got := withEntryPointHint(in, tc.code, tc.language, tc.input)
			if !errors.Is(got, in) {
				t.Errorf("the server's error should be kept, got %v", got)
			}
			if tc.want == "" {
				if got.Error() != in.Error() {
					t.Errorf("expected no hint, got %v", got)
				}
				return
			}
			if !strings.Contains(got.Error(), tc.want) {
				t.Errorf("hint %q does not contain %q", got, tc.want)
			}
		})
	}
}
