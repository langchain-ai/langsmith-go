package cmd

import (
	"bytes"
	"encoding/json"
	"net/http"
	"os"
	"path/filepath"
	"testing"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSplitSnapshotRef(t *testing.T) {
	tests := []struct {
		ref      string
		wantName string
		wantTag  string
	}{
		{ref: "my-snap", wantName: "my-snap", wantTag: ""},
		{ref: "my-snap:v2", wantName: "my-snap", wantTag: "v2"},
		{ref: "my-snap:2026081101", wantName: "my-snap", wantTag: "2026081101"},
		{ref: "my-snap:latest", wantName: "my-snap", wantTag: "latest"},
		// A trailing colon carries no tag, so the server applies its default.
		{ref: "my-snap:", wantName: "my-snap", wantTag: ""},
		// Split on the first colon only; the server rejects a colon inside a name.
		{ref: "my-snap:v2:extra", wantName: "my-snap", wantTag: "v2:extra"},
		{ref: "", wantName: "", wantTag: ""},
	}

	for _, tc := range tests {
		t.Run(tc.ref, func(t *testing.T) {
			name, tag := splitSnapshotRef(tc.ref)

			assert.Equal(t, tc.wantName, name)
			assert.Equal(t, tc.wantTag, tag)
		})
	}
}

func TestSnapshotRenderDescription(t *testing.T) {
	tests := []struct {
		name        string
		description string
		want        string
	}{
		{name: "empty is dangling", description: "", want: "Description: -"},
		{name: "set", description: "Python 3.12 with uv", want: "Description: Python 3.12 with uv"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			err := snapshotGetCommand.Render.RenderText(&out, langsmith.SnapshotResponse{
				ID:          "snap-1",
				Name:        "my-snap",
				Description: tc.description,
			})

			require.NoError(t, err)
			assert.Contains(t, out.String(), tc.want)
		})
	}
}

func TestJoinOrDashRendersSnapshotTags(t *testing.T) {
	tests := []struct {
		name string
		tags []string
		want string
	}{
		{name: "no tags is dangling", tags: nil, want: "Tags:        -"},
		{name: "one tag", tags: []string{"latest"}, want: "Tags:        latest"},
		{name: "many tags", tags: []string{"latest", "v2"}, want: "Tags:        latest, v2"},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var out bytes.Buffer

			err := snapshotGetCommand.Render.RenderText(&out, langsmith.SnapshotResponse{
				ID:   "snap-1",
				Name: "my-snap",
				Tags: tc.tags,
			})

			require.NoError(t, err)
			assert.Contains(t, out.String(), tc.want)
		})
	}
}

func TestParseBuildArgs(t *testing.T) {
	t.Setenv("FROM_ENV", "env-value")

	got, err := parseBuildArgs([]string{"VERSION=1", "EMPTY=", "WITH_EQUALS=a=b", "FROM_ENV", "UNSET_BUILD_ARG_FOR_TEST"})

	require.NoError(t, err)
	assert.Equal(t, map[string]string{
		"VERSION":     "1",
		"EMPTY":       "",
		"WITH_EQUALS": "a=b",
		"FROM_ENV":    "env-value",
	}, got)

	_, err = parseBuildArgs([]string{"=value"})
	assert.ErrorContains(t, err, "invalid --build-arg")
}

func TestSnapshotBuildRejectsInvalidFlagCombinations(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		t.Errorf("unexpected request %s %s", r.Method, r.URL.Path)
		http.Error(w, "unexpected", http.StatusInternalServerError)
	})

	tests := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{name: "no source and no Dockerfile", args: []string{"--context", t.TempDir()}, wantErr: "no Dockerfile in"},
		{name: "both sources", args: []string{"--docker-image", "ubuntu:24.04", "--dockerfile", "Dockerfile"}, wantErr: "none of the others can be"},
		{name: "build arg with image", args: []string{"--docker-image", "ubuntu:24.04", "--build-arg", "A=1"}, wantErr: "--build-arg applies to --dockerfile only"},
		{name: "vcpus with image", args: []string{"--docker-image", "ubuntu:24.04", "--vcpus", "2"}, wantErr: "--vcpus applies to --dockerfile only"},
		{name: "registry with dockerfile", args: []string{"--dockerfile", "Dockerfile", "--registry-id", "reg-1"}, wantErr: "--registry-id applies to --docker-image only"},
		{name: "bad memory", args: []string{"--dockerfile", "Dockerfile", "--memory", "lots"}, wantErr: "invalid --memory"},
		{name: "bad build arg", args: []string{"--dockerfile", "Dockerfile", "--build-arg", "=1"}, wantErr: "invalid --build-arg"},
	}
	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			args := append([]string{"sandboxes", "snapshot", "build", "my-snap", "--api-key", "test-api-key", "--api-url", ts.URL}, tc.args...)

			_, err := executeCommand(t, args...)

			assert.ErrorContains(t, err, tc.wantErr)
		})
	}
}

func TestSnapshotBuildFromDockerfileSizesBuilder(t *testing.T) {
	contextDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(contextDir, "Dockerfile"), []byte("FROM scratch\n"), 0o644))

	var createBody map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/boxes" {
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		}
		http.Error(w, `{"detail":"quota exceeded"}`, http.StatusForbidden)
	})

	_, err := executeCommand(t, "sandboxes", "snapshot", "build", "my-snap:v2",
		"--api-key", "test-api-key", "--api-url", ts.URL,
		"--context", contextDir,
		"--capacity", "8gb", "--vcpus", "2", "--memory", "8gb")

	assert.ErrorContains(t, err, "building snapshot from Dockerfile")
	assert.Equal(t, float64(2), createBody["vcpus"])
	assert.Equal(t, float64(8<<30), createBody["mem_bytes"])
	assert.Equal(t, float64(8<<30), createBody["fs_capacity_bytes"])
}

func TestSnapshotBuildFromDockerfileLeavesCapacityToBuilderByDefault(t *testing.T) {
	contextDir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(contextDir, "Dockerfile"), []byte("FROM scratch\n"), 0o644))

	var createBody map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/boxes" {
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		}
		http.Error(w, `{"detail":"quota exceeded"}`, http.StatusForbidden)
	})

	_, err := executeCommand(t, "sandboxes", "snapshot", "build", "my-snap",
		"--api-key", "test-api-key", "--api-url", ts.URL,
		"--context", contextDir)

	require.Error(t, err)
	require.NotNil(t, createBody)
	assert.NotContains(t, createBody, "fs_capacity_bytes")
}

func TestSnapshotBuildResolvesDockerfileFromWorkingDirectory(t *testing.T) {
	root := t.TempDir()
	require.NoError(t, os.MkdirAll(filepath.Join(root, "build"), 0o755))
	require.NoError(t, os.MkdirAll(filepath.Join(root, "app"), 0o755))
	require.NoError(t, os.WriteFile(filepath.Join(root, "build", "app.Dockerfile"), []byte("FROM scratch\n"), 0o644))
	t.Chdir(root)

	var createBody map[string]any
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		if r.Method == http.MethodPost && r.URL.Path == "/api/v2/sandboxes/boxes" {
			_ = json.NewDecoder(r.Body).Decode(&createBody)
		}
		http.Error(w, `{"detail":"quota exceeded"}`, http.StatusForbidden)
	})

	_, err := executeCommand(t, "sandboxes", "snapshot", "build", "my-snap",
		"--api-key", "test-api-key", "--api-url", ts.URL,
		"--dockerfile", "build/app.Dockerfile", "--context", "app")

	assert.ErrorContains(t, err, "quota exceeded")
	assert.NotNil(t, createBody, "the Dockerfile outside the context was not found from the working directory")
}
