package cmd

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"slices"
	"strconv"
	"strings"
	"sync"
	"testing"
)

const destructiveTestProjectID = "519bb9dd-079b-4488-8610-e330951ea3e4"

// destructiveCase is how the tests below reach the confirmation of a
// hand-written command: the arguments after `langsmith` (keep says whether the
// local data it would replace must survive), the skip flag the refusal must
// name ("" when the command has none), whether the command used to run with
// no prompt, whether its earlier prompt took a "y" with no newline, whether it
// changes only local files, and whether it is generated.
type destructiveCase struct {
	args                func(t *testing.T, keep bool) []string
	skipFlag            string
	unprompted          bool
	acceptsUnterminated bool
	local               bool
	generated           bool
}

// nonEmptyDir returns a directory holding one file and, when keep is set,
// checks at cleanup that the file is still there.
func nonEmptyDir(t *testing.T, dir string, keep bool) string {
	t.Helper()
	seedFile(t, dir, "keep.txt", "local edits")
	if keep {
		t.Cleanup(func() {
			if _, err := os.Stat(filepath.Join(dir, "keep.txt")); err != nil {
				t.Errorf("expected %s to be left alone: %v", dir, err)
			}
		})
	}
	return dir
}

func fixedArgs(args ...string) func(*testing.T, bool) []string {
	return func(*testing.T, bool) []string { return args }
}

// handWrittenDestructiveCases covers every hand-written command marked with
// confirmsAnnotation. Generated deletes need no entry: they confirm before
// the generated tree parses any flag.
var handWrittenDestructiveCases = map[string]destructiveCase{
	"langsmith apps delete": {args: fixedArgs("apps", "delete", "my-app"), acceptsUnterminated: true},
	"langsmith apps pull": {
		args: func(t *testing.T, keep bool) []string {
			dir := t.TempDir()
			t.Chdir(dir)
			nonEmptyDir(t, filepath.Join(dir, "my-app"), keep)
			return []string{"apps", "pull", "my-app"}
		},
		skipFlag:            "--force",
		acceptsUnterminated: true,
		local:               true,
	},
	"langsmith datasets delete": {args: fixedArgs("datasets", "delete", "eval-set"), acceptsUnterminated: true},
	"langsmith evaluators create-llm": {
		args: func(t *testing.T, _ bool) []string {
			modelConfig := filepath.Join(t.TempDir(), "model.json")
			seedFile(t, filepath.Dir(modelConfig), "model.json", `{"type":"chat"}`)
			return []string{"evaluators", "create-llm", "--name", "accuracy", "--dataset", "eval-set", "--replace",
				"--hub-ref", "my-org/new:latest", "--model-config", modelConfig}
		},
		acceptsUnterminated: true,
	},
	"langsmith evaluators delete": {args: fixedArgs("evaluators", "delete", testEvaluatorID)},
	"langsmith evaluators rule delete": {
		args: fixedArgs("evaluators", "rule", "delete", "accuracy", "--dataset", "eval-set"),
	},
	"langsmith evaluators upload": {
		args: func(t *testing.T, _ bool) []string {
			dir := t.TempDir()
			seedFile(t, dir, "eval.py", "def check(run, example):\n    return {}\n")
			return []string{"evaluators", "upload", filepath.Join(dir, "eval.py"), "--name", "accuracy",
				"--function", "check", "--dataset", "eval-set", "--replace"}
		},
		acceptsUnterminated: true,
	},
	"langsmith examples delete": {args: fixedArgs("examples", "delete", replaceRuleID), acceptsUnterminated: true},
	"langsmith hub delete":      {args: fixedArgs("hub", "delete", "acme/my-skill"), acceptsUnterminated: true},
	"langsmith hub pull": {
		args: func(t *testing.T, keep bool) []string {
			return []string{"hub", "pull", "acme/my-skill", "--dir", nonEmptyDir(t, t.TempDir(), keep)}
		},
		acceptsUnterminated: true,
		local:               true,
	},
	"langsmith insights config delete": {
		args: fixedArgs("insights", "config", "delete", replaceRuleID, "--project-id", destructiveTestProjectID),
	},
	"langsmith profiles delete": {
		args: func(t *testing.T, _ bool) []string {
			config := filepath.Join(t.TempDir(), "config.json")
			seedFile(t, filepath.Dir(config), "config.json", `{"profiles": {"work": {"api_key": "work-key"}}}`)
			if err := os.Chmod(config, 0o600); err != nil {
				t.Fatal(err)
			}
			t.Setenv("LANGSMITH_CONFIG_FILE", config)
			return []string{"profiles", "delete", "work"}
		},
		unprompted: true,
		local:      true,
	},
	"langsmith projects delete": {
		args:     fixedArgs("projects", "delete", "--project-id", destructiveTestProjectID),
		skipFlag: "",
	},
	"langsmith prompt delete":             {args: fixedArgs("prompt", "delete", "acme/my-prompt"), acceptsUnterminated: true},
	"langsmith sandboxes delete":          {args: fixedArgs("sandboxes", "delete", "my-box"), unprompted: true},
	"langsmith sandboxes snapshot delete": {args: fixedArgs("sandboxes", "snapshot", "delete", "my-snap"), unprompted: true},
}

// destructiveTestServer answers the reads commands make before they confirm
// and records every other request. Unless allowWrites, those fail.
func destructiveTestServer(t *testing.T, allowWrites bool) (string, *sync.Mutex, *[]string) {
	t.Helper()
	var mu sync.Mutex
	var writes []string
	rules := &fullReplaceRuleServer{rule: codeRule()}
	ruleHandler := rules.handler(t)
	apps := appsSourceServer(t, []customApp{{ID: testAppID, Name: "my-app"}}, map[string][]byte{testAppID: tarGz(t, map[string]string{"src/App.tsx": "fresh"})})

	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodGet {
			mu.Lock()
			writes = append(writes, r.Method+" "+r.URL.Path)
			mu.Unlock()
			if !allowWrites {
				http.Error(w, "unexpected write", http.StatusInternalServerError)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte("{}"))
			return
		}
		w.Header().Set("Content-Type", "application/json")
		switch {
		case strings.HasPrefix(r.URL.Path, "/v1/platform/custom-apps"):
			apps.ServeHTTP(w, r)
		case r.URL.Path == "/api/v1/platform/evaluators/"+testEvaluatorID:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": testEvaluatorID, "name": "accuracy"})
		case r.URL.Path == "/api/v2/sessions/"+destructiveTestProjectID:
			_ = json.NewEncoder(w).Encode(map[string]any{"id": destructiveTestProjectID, "tenant_id": "t", "name": "production"})
		case strings.HasPrefix(r.URL.Path, "/api/v1/platform/hub/repos/"):
			_ = json.NewEncoder(w).Encode(map[string]any{
				"commit_hash": "h1",
				"files":       map[string]any{"SKILL.md": map[string]any{"type": "file", "content": "x"}},
			})
		default:
			ruleHandler(w, r)
		}
	})
	t.Cleanup(setupTestEnv(t, ts.URL))
	return ts.URL, &mu, &writes
}

// destructiveCases returns how to reach the confirmation of every command
// that deletes or overwrites data, keyed by command path.
func destructiveCases(t *testing.T) map[string]destructiveCase {
	t.Helper()
	cases := map[string]destructiveCase{}
	for _, c := range leafCommands(NewRootCmd("test", "test")) {
		if !confirmsDelete(c) {
			continue
		}
		if c.op != nil {
			cases[c.path] = destructiveCase{
				args:      fixedArgs(strings.Fields(strings.TrimPrefix(c.path, "langsmith "))...),
				skipFlag:  "--yes",
				generated: true,
			}
			continue
		}
		tc, ok := handWrittenDestructiveCases[c.path]
		if !ok {
			t.Fatalf("%s confirms but has no entry in handWrittenDestructiveCases", c.path)
		}
		if tc.skipFlag == "" && c.path != "langsmith projects delete" {
			tc.skipFlag = "--yes"
		}
		cases[c.path] = tc
	}
	return cases
}

// runDestructive runs args against a test server with stdin that is not a
// terminal and returns the error, stderr and the writes the server received.
func runDestructive(t *testing.T, tc destructiveCase, stdin string, allowWrites bool) (error, string, []string) {
	t.Helper()
	url, mu, writes := destructiveTestServer(t, allowWrites)
	root := NewRootCmd("test", "test")
	var stderr bytes.Buffer
	root.SetIn(strings.NewReader(stdin))
	root.SetErr(&stderr)
	root.SetOut(io.Discard)
	// Global flags go last: generated commands read them from their own arguments.
	root.SetArgs(append(tc.args(t, !allowWrites), "--api-key", "test-api-key", "--api-url", url))
	var err error
	captureStdout(t, func() { err = root.Execute() })
	mu.Lock()
	defer mu.Unlock()
	return err, stderr.String(), slices.Clone(*writes)
}

func notTerminal(t *testing.T) {
	t.Helper()
	old := inputIsTerminal
	inputIsTerminal = func(io.Reader) bool { return false }
	t.Cleanup(func() { inputIsTerminal = old })
}

// TestDestructiveCommandsRequireYes runs every command that prompted before it
// confirmed, with stdin that is not a terminal and gives no answer, or a "y"
// with no newline where the earlier prompt needed a complete line: each must
// refuse, name the flag that skips the prompt, and change nothing.
func TestDestructiveCommandsRequireYes(t *testing.T) {
	notTerminal(t)
	for path, tc := range destructiveCases(t) {
		if tc.unprompted && !requireYesWithoutTerminal {
			continue
		}
		stdins := []string{""}
		if !tc.acceptsUnterminated || requireYesWithoutTerminal {
			stdins = append(stdins, "y")
		}
		for _, stdin := range stdins {
			t.Run(path+" stdin="+strconv.Quote(stdin), func(t *testing.T) {
				err, stderr, writes := runDestructive(t, tc, stdin, false)
				if err == nil || !strings.Contains(err.Error(), "stdin is not a terminal") {
					t.Fatalf("expected a refusal because stdin is not a terminal, got %v\nstderr: %s", err, stderr)
				}
				if tc.skipFlag != "" && !strings.Contains(err.Error(), "rerun the command with "+tc.skipFlag) {
					t.Errorf("refusal should name %s: %v", tc.skipFlag, err)
				}
				if !strings.Contains(stderr, "WARNING: This permanently") {
					t.Errorf("expected the warning on stderr, got %q", stderr)
				}
				var reported ReportedError
				if errors.As(err, &reported) && !strings.Contains(stderr, err.Error()) {
					t.Errorf("a reported error must be on stderr: %v\nstderr: %s", err, stderr)
				}
				if len(writes) > 0 {
					t.Errorf("expected no writes, got %v", writes)
				}
			})
		}
	}
}

// TestDestructiveCommandsWarnWithoutTerminal covers what scripts relied on
// before confirmation required a terminal: a command that ran with no prompt
// still runs with no answer, and a piped "y" still confirms the others (with
// no newline only where the earlier prompt took one). Both
// complete and print a deprecation warning that names the skip flag.
// Generated deletes confirm through the same confirmDelete, but the generated
// tree exits the process on the arguments this test gives them.
func TestDestructiveCommandsWarnWithoutTerminal(t *testing.T) {
	if requireYesWithoutTerminal {
		t.Skip("confirmation requires a terminal")
	}
	notTerminal(t)
	for path, tc := range destructiveCases(t) {
		if tc.generated {
			continue
		}
		t.Run(path, func(t *testing.T) {
			stdin := "y\n"
			if tc.acceptsUnterminated {
				stdin = "y"
			}
			if tc.unprompted {
				stdin = ""
			}
			err, stderr, writes := runDestructive(t, tc, stdin, true)
			if err != nil {
				t.Fatalf("expected the command to complete, got %v\nstderr: %s", err, stderr)
			}
			want := "deprecated; from " + strictConfirmationVersion + " it will refuse"
			if !strings.Contains(stderr, want) {
				t.Errorf("expected a deprecation warning containing %q, got %q", want, stderr)
			}
			if tc.skipFlag != "" && !strings.Contains(stderr, "Pass "+tc.skipFlag+" now") {
				t.Errorf("warning should name %s: %q", tc.skipFlag, stderr)
			}
			if !tc.local && len(writes) == 0 {
				t.Errorf("expected the change to reach the API")
			}
		})
	}
}

// TestDeletingCommandsConfirm keeps new hand-written commands that delete,
// remove, or replace from skipping confirmation.
func TestDeletingCommandsConfirm(t *testing.T) {
	// These lose no data: `projects issues runs add` links the run again.
	reversible := map[string]bool{"langsmith projects issues runs remove": true}
	for _, c := range leafCommands(NewRootCmd("test", "test")) {
		if c.op != nil || reversible[c.path] {
			continue
		}
		words := strings.Split(c.cobra.Name(), "-")
		destructive := slices.Contains(words, "delete") || slices.Contains(words, "remove") ||
			c.cobra.LocalFlags().Lookup("replace") != nil
		if destructive && !confirmsDelete(c) {
			t.Errorf("%s deletes or replaces data but does not confirm; use confirmDelete and addYesFlag", c.path)
		}
	}
}
