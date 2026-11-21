package cmd

import (
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"sync/atomic"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func marshalBody(t *testing.T, params any) map[string]any {
	t.Helper()
	raw, err := json.Marshal(params)
	require.NoError(t, err)
	var body map[string]any
	require.NoError(t, json.Unmarshal(raw, &body))
	return body
}

func TestSandboxCreateCmd_SendsLifecycleMemoryAndRuntimeSettings(t *testing.T) {
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v2/sandboxes/boxes", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"provisioning"}`))
	})

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "create", "my-vm",
		"--restore-memory", "never",
		"--no-persistent",
		"--idle-ttl-seconds", "0",
		"--delete-after-stop-seconds", "3600",
		"--cpu-millicores", "500",
		"--env-var", "LOG_LEVEL=debug",
		"--label", "team=infra",
		"--tag-value-id", "tv-1",
		"--volume", "data:/data",
		"--volume", "models:/models:ro",
		"--user", "app",
		"--work-dir", "/srv/app",
		"--mount-config", `{"mounts":[]}`,
		"--format", "json",
	)
	require.NoError(t, err)

	assert.NotContains(t, body, "preserve_memory_on_stop")
	assert.Equal(t, false, body["restore_memory"])
	assert.Equal(t, false, body["persistent"])
	assert.Equal(t, float64(0), body["idle_ttl_seconds"])
	assert.Equal(t, float64(3600), body["delete_after_stop_seconds"])
	assert.Equal(t, float64(500), body["cpu_millicores"])
	assert.Equal(t, map[string]any{"LOG_LEVEL": "debug"}, body["env_vars"])
	assert.Equal(t, map[string]any{"team": "infra"}, body["labels"])
	assert.Equal(t, []any{"tv-1"}, body["tag_value_ids"])
	assert.Equal(t, []any{
		map[string]any{"volume": "data", "mount_path": "/data"},
		map[string]any{"volume": "models", "mount_path": "/models", "read_only": true},
	}, body["volumes"])
	assert.Equal(t, map[string]any{"user": "app", "work_dir": "/srv/app"}, body["run_config"])
	assert.Equal(t, map[string]any{"mounts": []any{}}, body["mount_config"])
}

func TestSandboxCreateCmd_SendsExplicitPreserveMemoryFalse(t *testing.T) {
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"provisioning"}`))
	})

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "create", "my-vm",
		"--snapshot", "warm", "--preserve-memory-on-stop=false", "--format", "json")
	require.NoError(t, err)

	assert.Equal(t, false, body["preserve_memory_on_stop"])
}

func TestSandboxCreateParams_OmitsUnsetOptionalSettings(t *testing.T) {
	params, err := sandboxCreateParams("my-vm", &sandboxCreateInput{})
	require.NoError(t, err)
	body := marshalBody(t, params)

	for _, field := range []string{
		"preserve_memory_on_stop", "restore_memory", "persistent", "idle_ttl_seconds",
		"delete_after_stop_seconds", "env_vars", "labels", "volumes", "run_config", "mount_config",
	} {
		assert.NotContains(t, body, field)
	}
}

func TestSandboxCreateParams_MemoryAndPersistence(t *testing.T) {
	params, err := sandboxCreateParams("", &sandboxCreateInput{PreserveMemoryOnStop: true, PreserveMemoryOnStopSet: true})
	require.NoError(t, err)
	body := marshalBody(t, params)
	assert.Equal(t, true, body["preserve_memory_on_stop"])
	assert.NotContains(t, body, "persistent")

	// An explicit false is sent, so it overrides a snapshot that carries memory preservation.
	params, err = sandboxCreateParams("", &sandboxCreateInput{PreserveMemoryOnStopSet: true})
	require.NoError(t, err)
	assert.Equal(t, false, marshalBody(t, params)["preserve_memory_on_stop"])

	_, err = sandboxCreateParams("", &sandboxCreateInput{NoPersistent: true, PreserveMemoryOnStop: true, PreserveMemoryOnStopSet: true})
	require.ErrorContains(t, err, "--no-persistent cannot be combined")
}

func TestSandboxCreateParams_RestoreMemory(t *testing.T) {
	tests := []struct {
		value string
		want  any
	}{
		{value: "", want: nil},
		{value: "if-present", want: nil},
		{value: "always", want: true},
		{value: "never", want: false},
	}
	for _, tc := range tests {
		t.Run(tc.value, func(t *testing.T) {
			params, err := sandboxCreateParams("", &sandboxCreateInput{RestoreMemory: tc.value})
			require.NoError(t, err)
			body := marshalBody(t, params)
			if tc.want == nil {
				assert.NotContains(t, body, "restore_memory")
				return
			}
			assert.Equal(t, tc.want, body["restore_memory"])
		})
	}

	_, err := sandboxCreateParams("", &sandboxCreateInput{RestoreMemory: "sometimes"})
	require.ErrorContains(t, err, "--restore-memory")
}

func TestParseVolumeSpecs_Rejects(t *testing.T) {
	for _, spec := range []string{"data", ":/data", "data:", "data:/data:rx", "a:/b:ro:x"} {
		_, err := parseVolumeSpecs([]string{spec})
		assert.Error(t, err, spec)
	}
}

func TestParseKeyValues_RequiresEquals(t *testing.T) {
	_, err := parseKeyValues("label", []string{"team"})
	require.ErrorContains(t, err, "--label must be KEY=VALUE")

	got, err := parseKeyValues("label", []string{"a=1", "b=x=y", "c="})
	require.NoError(t, err)
	assert.Equal(t, map[string]string{"a": "1", "b": "x=y", "c": ""}, got)
}

func TestSandboxCreateCmd_WaitPrintsReadySandbox(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/boxes":
			_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"provisioning"}`))
		case r.URL.Path == "/api/v2/sandboxes/boxes/my-vm/status":
			_, _ = w.Write([]byte(`{"status":"ready"}`))
		case r.URL.Path == "/api/v2/sandboxes/boxes/my-vm":
			_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"ready"}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "create", "my-vm", "--wait", "--format", "json")
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &got))
	assert.Equal(t, "ready", got["status"])
}

func TestSandboxUpdateCmd_SendsOnlyChangedFields(t *testing.T) {
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPatch, r.Method)
		require.Equal(t, "/api/v2/sandboxes/boxes/my-vm", r.URL.Path)
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"id":"box-id","name":"renamed","status":"ready"}`))
	})

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "update", "my-vm",
		"--name", "renamed", "--idle-ttl-seconds", "0", "--work-dir", "/srv", "--format", "json")
	require.NoError(t, err)

	assert.Equal(t, map[string]any{
		"name":             "renamed",
		"idle_ttl_seconds": float64(0),
		"run_config":       map[string]any{"work_dir": "/srv"},
	}, body)
}

func TestSandboxUpdateCmd_RequiresAChange(t *testing.T) {
	_, err := executeCommand(t, "--api-key", "test-key", "sandbox", "update", "my-vm")
	require.ErrorContains(t, err, "nothing to update")
}

func TestSandboxListCmd_FiltersAndLimit(t *testing.T) {
	var query map[string][]string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v2/sandboxes/boxes", r.URL.Path)
		query = r.URL.Query()
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"items":[{"name":"a"},{"name":"b"},{"name":"c"}],"next_cursor":null}`))
	})

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "list",
		"--status", "ready", "--label", "team=infra", "--label", "ephemeral", "--name-contains", "ci-", "--mine",
		"--limit", "2", "--format", "json")
	require.NoError(t, err)

	assert.Equal(t, []string{"ready"}, query["status"])
	assert.Equal(t, []string{"team=infra", "ephemeral"}, query["label"])
	assert.Equal(t, []string{"ci-"}, query["name_contains"])
	assert.Equal(t, []string{"me"}, query["created_by"])

	var got []map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &got))
	assert.Len(t, got, 2)
}

func TestSnapshotCaptureParams(t *testing.T) {
	params, err := snapshotCaptureParams("warm:v1", &snapshotCaptureInput{
		BoxName:             "my-vm",
		IncludeMemory:       true,
		Include:             []string{"resources", "lifecycle"},
		Labels:              []string{"stage=warm"},
		RunConfig:           sandboxRunConfigInput{User: "app", EnvVars: []string{"A=1"}},
		ImageConfig:         `{"exposed_ports":[8080]}`,
		ConfirmSecretValues: true,
	})
	require.NoError(t, err)
	body := marshalBody(t, params)

	assert.Equal(t, "warm", body["name"])
	assert.Equal(t, "v1", body["tag"])
	assert.Equal(t, true, body["include_memory"])
	assert.Equal(t, []any{"resources", "lifecycle"}, body["include"])
	assert.Equal(t, map[string]any{"stage": "warm"}, body["labels"])
	assert.Equal(t, map[string]any{"user": "app", "env_vars": map[string]any{"A": "1"}}, body["run_config"])
	assert.Equal(t, map[string]any{"exposed_ports": []any{float64(8080)}}, body["image_config"])
	assert.Equal(t, true, body["confirm_allow_snapshot_secret_values"])
}

func TestSnapshotCaptureParams_Rejects(t *testing.T) {
	_, err := snapshotCaptureParams("s", &snapshotCaptureInput{BoxName: "b", Include: []string{"secrets"}})
	require.ErrorContains(t, err, "--include must be one of")

	_, err = snapshotCaptureParams("s", &snapshotCaptureInput{BoxName: "b", IncludeMemory: true, Checkpoint: "2026-03-29T00:09:28Z"})
	require.ErrorContains(t, err, "--include-memory")

	_, err = snapshotCaptureParams("s", &snapshotCaptureInput{})
	require.ErrorContains(t, err, "--from is required")
}

func TestSnapshotCaptureCmd_FromAndBoxAlias(t *testing.T) {
	for _, flag := range []string{"--from", "--box"} {
		t.Run(flag, func(t *testing.T) {
			var path string
			ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
				path = r.URL.Path
				w.Header().Set("Content-Type", "application/json")
				_, _ = w.Write([]byte(`{"id":"snap-1","name":"s","status":"building"}`))
			})

			_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "snapshot", "capture", "s",
				flag, "my-vm", "--format", "json")
			require.NoError(t, err)
			assert.Equal(t, "/api/v2/sandboxes/boxes/my-vm/snapshot", path)
		})
	}

	assert.True(t, snapshotCaptureCommand.Cobra().Flags().Lookup("box").Hidden)
}

func TestSnapshotBuildCmd_WaitPollsUntilReady(t *testing.T) {
	var polls atomic.Int32
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/snapshots":
			var body map[string]any
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			assert.Equal(t, map[string]any{"team": "infra"}, body["labels"])
			_, _ = w.Write([]byte(`{"id":"snap-1","name":"my-snap","status":"building"}`))
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sandboxes/snapshots/snap-1":
			polls.Add(1)
			_, _ = w.Write([]byte(`{"id":"snap-1","name":"my-snap","status":"ready","memory_snapshot_size_bytes":0}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "snapshot", "build", "my-snap",
		"--docker-image", "ubuntu:24.04", "--label", "team=infra", "--wait", "--format", "json")
	require.NoError(t, err)

	assert.Equal(t, int32(1), polls.Load())
	var got map[string]any
	require.NoError(t, json.Unmarshal([]byte(out), &got))
	assert.Equal(t, "ready", got["status"])
}

func TestSnapshotBuildCmd_WaitReportsFailure(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.Method == http.MethodPost {
			_, _ = w.Write([]byte(`{"id":"snap-1","status":"building"}`))
			return
		}
		_, _ = w.Write([]byte(`{"id":"snap-1","status":"failed","status_message":"image not found"}`))
	})

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "snapshot", "build", "my-snap",
		"--docker-image", "nope:latest", "--wait")
	require.ErrorContains(t, err, "image not found")
}

func TestSandboxExecCmd_SendsPerCommandRunConfig(t *testing.T) {
	t.Setenv("SANDBOX_EXEC_TEST_TOKEN", "from-local")
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sandboxes/boxes/my-vm":
			_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"ready","dataplane_url":"` + tsURL(t, r) + `"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/execute":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			_, _ = w.Write([]byte(`{"stdout":"","stderr":"","exit_code":0}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	captureStdout(t, func() {
		_, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "exec", "my-vm",
			"--user", "root", "--work-dir", "/srv", "--env", "SANDBOX_EXEC_TEST_TOKEN", "--env", "CI=1",
			"--timeout", "90", "--", "make", "test")
		require.NoError(t, err)
	})

	assert.Equal(t, map[string]any{
		"user":     "root",
		"work_dir": "/srv",
		"env_vars": map[string]any{"SANDBOX_EXEC_TEST_TOKEN": "from-local", "CI": "1"},
	}, body["run_config"])
	assert.Equal(t, float64(90), body["timeout_seconds"])
	assert.NotContains(t, body, "cwd")
	assert.NotContains(t, body, "env")
}

func TestParseSandboxFileRef(t *testing.T) {
	tests := []struct {
		arg     string
		sandbox string
		path    string
	}{
		{arg: "my-vm:/tmp/a.txt", sandbox: "my-vm", path: "/tmp/a.txt"},
		{arg: "./a.txt", path: "./a.txt"},
		{arg: "-", path: "-"},
		{arg: "notes:draft.txt", path: "notes:draft.txt"},
		{arg: "dir/my-vm:/x", path: "dir/my-vm:/x"},
		{arg: ":/tmp", path: ":/tmp"},
	}
	for _, tc := range tests {
		t.Run(tc.arg, func(t *testing.T) {
			got := parseSandboxFileRef(tc.arg)
			assert.Equal(t, tc.sandbox, got.Sandbox)
			assert.Equal(t, tc.path, got.Path)
		})
	}
}

func dataplaneFileServer(t *testing.T, files map[string][]byte) string {
	t.Helper()
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sandboxes/boxes/my-vm":
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"ready","dataplane_url":"` + tsURL(t, r) + `"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/upload":
			mr, err := r.MultipartReader()
			require.NoError(t, err)
			part, err := mr.NextPart()
			require.NoError(t, err)
			files[r.URL.Query().Get("path")] = readPart(t, part)
			w.WriteHeader(http.StatusOK)
		case r.Method == http.MethodGet && r.URL.Path == "/download":
			content, ok := files[r.URL.Query().Get("path")]
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/octet-stream")
			_, _ = w.Write(content)
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})
	return ts.URL
}

func readPart(t *testing.T, part *multipart.Part) []byte {
	t.Helper()
	b, err := io.ReadAll(part)
	require.NoError(t, err)
	return b
}

func TestSandboxCopyCmd_UploadIntoDirectoryKeepsName(t *testing.T) {
	files := map[string][]byte{}
	url := dataplaneFileServer(t, files)
	src := filepath.Join(t.TempDir(), "config.yaml")
	require.NoError(t, os.WriteFile(src, []byte("key: value\n"), 0o600))

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "cp", src, "my-vm:/srv/app/")
	require.NoError(t, err)

	assert.Equal(t, []byte("key: value\n"), files["/srv/app/config.yaml"])
	assert.Contains(t, out, "my-vm:/srv/app/config.yaml")
}

func TestSandboxCopyCmd_DownloadToFileAndStdout(t *testing.T) {
	files := map[string][]byte{"/tmp/report.txt": []byte("report body")}
	url := dataplaneFileServer(t, files)
	dir := t.TempDir()

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "cp", "my-vm:/tmp/report.txt", dir)
	require.NoError(t, err)
	got, err := os.ReadFile(filepath.Join(dir, "report.txt"))
	require.NoError(t, err)
	assert.Equal(t, "report body", string(got))

	var out string
	stdout := captureStdout(t, func() {
		out, err = executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "cp", "my-vm:/tmp/report.txt", "-")
	})
	require.NoError(t, err)
	assert.Empty(t, out)
	assert.Equal(t, "report body", stdout)
}

func TestSandboxCopyCmd_RequiresExactlyOneSandboxSide(t *testing.T) {
	_, err := executeCommand(t, "--api-key", "test-key", "sandbox", "cp", "a.txt", "b.txt")
	require.ErrorContains(t, err, "must be a sandbox path")

	_, err = executeCommand(t, "--api-key", "test-key", "sandbox", "cp", "vm-a:/x", "vm-b:/y")
	require.ErrorContains(t, err, "between two sandboxes")
}

func TestSandboxCopyCmd_RejectsLocalDirectory(t *testing.T) {
	url := dataplaneFileServer(t, map[string][]byte{})
	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "cp", t.TempDir(), "my-vm:/tmp/")
	require.ErrorContains(t, err, "is a directory")
}

func TestSandboxListFilesCmd(t *testing.T) {
	var body map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/api/v2/sandboxes/boxes/my-vm":
			_, _ = w.Write([]byte(`{"id":"box-id","name":"my-vm","status":"ready","dataplane_url":"` + tsURL(t, r) + `"}`))
		case r.Method == http.MethodPost && r.URL.Path == "/glob":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			_, _ = w.Write([]byte(`{"matches":[{"path":"/srv/app/main.py","is_dir":false,"size_bytes":42,"modified_at":"2026-11-19T10:00:00Z"},{"path":"/srv/app/lib","is_dir":true,"size_bytes":0,"modified_at":"2026-11-19T10:00:00Z"}],"truncated":true}`))
		default:
			t.Fatalf("unexpected request %s %s", r.Method, r.URL.Path)
		}
	})

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "ls", "my-vm", "/srv/app",
		"--pattern", "**/*", "--limit", "10", "--format", "pretty")
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"pattern": "**/*", "path": "/srv/app", "limit": float64(10)}, body)
	assert.Contains(t, out, "/srv/app/main.py")
	assert.Regexp(t, `(?m)^d\s+0 .* /srv/app/lib$`, out)
	assert.Contains(t, out, "truncated")
}

func TestSandboxListFilesCmd_RequiresAbsolutePath(t *testing.T) {
	_, err := executeCommand(t, "--api-key", "test-key", "sandbox", "ls", "my-vm", "srv")
	require.ErrorContains(t, err, "path must be absolute")
}
