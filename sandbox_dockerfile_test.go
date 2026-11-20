package langsmith_test

import (
	"archive/tar"
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"sync"
	"testing"
	"time"

	"golang.org/x/net/websocket"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

// fakeDockerfileBuilder serves the control-plane and dataplane endpoints a
// Dockerfile snapshot build uses, recording what the client sent.
type fakeDockerfileBuilder struct {
	t             *testing.T
	buildExitCode int

	mu            sync.Mutex
	createBody    map[string]any
	uploadedPath  string
	uploadedTar   []byte
	buildCommand  string
	captureBody   map[string]any
	captureCalled bool
	deleted       []string
}

func (f *fakeDockerfileBuilder) serve() *httptest.Server {
	mux := http.NewServeMux()
	var srv *httptest.Server
	box := func(status string) map[string]any {
		return map[string]any{
			"id":            "box-id",
			"name":          "builder",
			"status":        status,
			"dataplane_url": srv.URL + "/dp",
		}
	}
	writeJSON := func(w http.ResponseWriter, v any) {
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(v)
	}

	mux.HandleFunc("POST /api/v2/sandboxes/boxes", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.mu.Lock()
		f.createBody = body
		f.mu.Unlock()
		writeJSON(w, box("provisioning"))
	})
	mux.HandleFunc("GET /api/v2/sandboxes/boxes/builder/status", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"status": "ready"})
	})
	mux.HandleFunc("GET /api/v2/sandboxes/boxes/builder", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, box("ready"))
	})
	mux.HandleFunc("DELETE /api/v2/sandboxes/boxes/{name}", func(w http.ResponseWriter, r *http.Request) {
		f.mu.Lock()
		f.deleted = append(f.deleted, r.PathValue("name"))
		f.mu.Unlock()
		w.WriteHeader(http.StatusNoContent)
	})
	mux.HandleFunc("POST /dp/upload", func(w http.ResponseWriter, r *http.Request) {
		file, _, err := r.FormFile("file")
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		data, _ := io.ReadAll(file)
		f.mu.Lock()
		f.uploadedPath = r.URL.Query().Get("path")
		f.uploadedTar = data
		f.mu.Unlock()
		writeJSON(w, map[string]any{})
	})
	mux.HandleFunc("POST /dp/execute", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"stdout": "", "stderr": "", "exit_code": 0})
	})
	mux.Handle("/dp/execute/ws", websocket.Handler(func(ws *websocket.Conn) {
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			f.t.Errorf("receive execute: %v", err)
			return
		}
		command, _ := msg["command"].(string)
		f.mu.Lock()
		f.buildCommand = command
		f.mu.Unlock()
		_ = websocket.Message.Send(ws, `{"type":"started","command_id":"cmd-1","pid":1}`)
		_ = websocket.Message.Send(ws, `{"type":"stderr","data":"#1 building\n","offset":0}`)
		exit, _ := json.Marshal(map[string]any{"type": "exit", "exit_code": f.buildExitCode})
		_ = websocket.Message.Send(ws, string(exit))
	}))
	mux.HandleFunc("POST /api/v2/sandboxes/boxes/builder/snapshot", func(w http.ResponseWriter, r *http.Request) {
		var body map[string]any
		_ = json.NewDecoder(r.Body).Decode(&body)
		f.mu.Lock()
		f.captureCalled = true
		f.captureBody = body
		f.mu.Unlock()
		writeJSON(w, map[string]any{"id": "snap-1", "name": body["name"], "status": "building"})
	})
	mux.HandleFunc("GET /api/v2/sandboxes/snapshots/snap-1", func(w http.ResponseWriter, r *http.Request) {
		writeJSON(w, map[string]any{"id": "snap-1", "name": "my-snap", "status": "ready"})
	})

	srv = httptest.NewServer(mux)
	f.t.Cleanup(srv.Close)
	return srv
}

func writeTestFiles(t *testing.T, root string, files map[string]string) {
	t.Helper()
	for name, content := range files {
		p := filepath.Join(root, filepath.FromSlash(name))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
			t.Fatal(err)
		}
	}
}

func tarFileNames(t *testing.T, data []byte) map[string]string {
	t.Helper()
	files := map[string]string{}
	tr := tar.NewReader(bytes.NewReader(data))
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			return files
		}
		if err != nil {
			t.Fatalf("read uploaded tar: %v", err)
		}
		switch hdr.Typeflag {
		case tar.TypeReg:
			content, _ := io.ReadAll(tr)
			files[hdr.Name] = string(content)
		case tar.TypeSymlink:
			files[hdr.Name] = "-> " + hdr.Linkname
		}
	}
}

func uploadedNames(t *testing.T, data []byte) string {
	t.Helper()
	files := tarFileNames(t, data)
	names := make([]string, 0, len(files))
	for name := range files {
		names = append(names, name)
	}
	sort.Strings(names)
	return strings.Join(names, ",")
}

func newDockerfileTestClient(srv *httptest.Server) *langsmith.Client {
	return langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
}

func TestSnapshotNewFromDockerfile(t *testing.T) {
	contextDir := t.TempDir()
	writeTestFiles(t, contextDir, map[string]string{
		"docker/app.Dockerfile": "FROM ubuntu:24.04\n",
		"src/main.py":           "print('hi')\n",
		".git/HEAD":             "ref: refs/heads/main\n",
	})
	if err := os.Symlink("src/main.py", filepath.Join(contextDir, "main.py")); err != nil {
		t.Fatal(err)
	}

	fake := &fakeDockerfileBuilder{t: t}
	client := newDockerfileTestClient(fake.serve())

	var logs strings.Builder
	snapshot, err := client.Sandboxes.Snapshots.NewFromDockerfile(context.Background(), langsmith.SnapshotDockerfileParams{
		Name:            "my-snap",
		Tag:             "v2",
		Dockerfile:      "docker/app.Dockerfile",
		Context:         contextDir,
		BuildArgs:       map[string]string{"VERSION": "1", "A": "it's"},
		Target:          "runtime",
		Description:     "Python app",
		FsCapacityBytes: 8 << 30,
		Vcpus:           2,
		Timeout:         5 * time.Second,
		OnBuildLog:      func(s string) { logs.WriteString(s) },
	})
	if err != nil {
		t.Fatalf("NewFromDockerfile returned error: %v", err)
	}
	if snapshot.ID != "snap-1" || snapshot.Status != "ready" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()

	if fake.createBody["vcpus"] != float64(2) || fake.createBody["fs_capacity_bytes"] != float64(8<<30) {
		t.Errorf("builder sandbox not sized from params: %#v", fake.createBody)
	}
	if !strings.HasPrefix(fake.uploadedPath, "/var/lib/langsmith-build/") {
		t.Errorf("context uploaded to %q, want under /var/lib/langsmith-build/", fake.uploadedPath)
	}
	files := tarFileNames(t, fake.uploadedTar)
	want := "context/.git/HEAD,context/docker/app.Dockerfile,context/main.py,context/src/main.py,dockerfile/app.Dockerfile"
	if got := uploadedNames(t, fake.uploadedTar); got != want {
		t.Errorf("uploaded files = %s, want %s", got, want)
	}
	if files["context/main.py"] != "-> src/main.py" {
		t.Errorf("symlink archived as %q, want a link to src/main.py", files["context/main.py"])
	}

	for _, want := range []string{
		"'filename=app.Dockerfile'",
		"'target=runtime'",
		`'build-arg:A=it'"'"'s' '--opt' 'build-arg:VERSION=1'`,
		"| docker load",
	} {
		if !strings.Contains(fake.buildCommand, want) {
			t.Errorf("build command missing %q:\n%s", want, fake.buildCommand)
		}
	}
	if logs.String() != "#1 building\n" {
		t.Errorf("OnBuildLog received %q", logs.String())
	}

	imageRef, _ := fake.captureBody["docker_image"].(string)
	if imageRef == "" || !strings.Contains(fake.buildCommand, "'--output' 'type=docker,name="+imageRef+"'") {
		t.Errorf("captured image %q is not the image the build produced", imageRef)
	}
	if fake.captureBody["name"] != "my-snap" || fake.captureBody["tag"] != "v2" ||
		fake.captureBody["description"] != "Python app" || fake.captureBody["fs_capacity_bytes"] != float64(8<<30) {
		t.Errorf("unexpected capture body: %#v", fake.captureBody)
	}
	if len(fake.deleted) != 1 || fake.deleted[0] != "builder" {
		t.Errorf("builder sandbox deletions = %v, want [builder]", fake.deleted)
	}
}

func TestSnapshotNewFromDockerfileBuildFailure(t *testing.T) {
	contextDir := t.TempDir()
	writeTestFiles(t, contextDir, map[string]string{"Dockerfile": "FROM scratch\n"})

	fake := &fakeDockerfileBuilder{t: t, buildExitCode: 1}
	client := newDockerfileTestClient(fake.serve())

	_, err := client.Sandboxes.Snapshots.NewFromDockerfile(context.Background(), langsmith.SnapshotDockerfileParams{
		Name:    "my-snap",
		Context: contextDir,
		Timeout: 5 * time.Second,
	})
	var creationErr *langsmith.SandboxResourceCreationError
	if !errors.As(err, &creationErr) {
		t.Fatalf("expected SandboxResourceCreationError, got %v", err)
	}
	if !strings.Contains(creationErr.Message, "exited with code 1") || !strings.Contains(creationErr.Message, "#1 building") {
		t.Errorf("error does not carry the exit code and build output: %q", creationErr.Message)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	if fake.captureCalled {
		t.Error("snapshot captured after a failed build")
	}
	if len(fake.deleted) != 1 {
		t.Errorf("builder sandbox deletions = %v, want one", fake.deleted)
	}
}

func TestSnapshotNewFromDockerfileAcceptsDockerfileOutsideContext(t *testing.T) {
	root := t.TempDir()
	writeTestFiles(t, root, map[string]string{
		"build/Dockerfile": "FROM scratch\n",
		"ctx/readme.txt":   "hi\n",
	})

	fake := &fakeDockerfileBuilder{t: t}
	client := newDockerfileTestClient(fake.serve())

	_, err := client.Sandboxes.Snapshots.NewFromDockerfile(context.Background(), langsmith.SnapshotDockerfileParams{
		Name:       "my-snap",
		Dockerfile: filepath.Join(root, "build", "Dockerfile"),
		Context:    filepath.Join(root, "ctx"),
		Timeout:    5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewFromDockerfile returned error: %v", err)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	if got, want := uploadedNames(t, fake.uploadedTar), "context/readme.txt,dockerfile/Dockerfile"; got != want {
		t.Errorf("uploaded files = %s, want %s", got, want)
	}
}

func TestSnapshotNewFromDockerfileAppliesDockerignore(t *testing.T) {
	contextDir := t.TempDir()
	writeTestFiles(t, contextDir, map[string]string{
		"Dockerfile":          "FROM scratch\n",
		".dockerignore":       "# local only\n*.log\n!keep.log\nsecrets/\n!secrets/public.pem\nnode_modules\n/build\n",
		"app.log":             "noise\n",
		"keep.log":            "kept\n",
		"secrets/key":         "private\n",
		"secrets/public.pem":  "public\n",
		"node_modules/x/a.js": "dep\n",
		"build/out":           "artifact\n",
		"src/build/kept":      "source\n",
	})

	fake := &fakeDockerfileBuilder{t: t}
	client := newDockerfileTestClient(fake.serve())

	_, err := client.Sandboxes.Snapshots.NewFromDockerfile(context.Background(), langsmith.SnapshotDockerfileParams{
		Name:    "my-snap",
		Context: contextDir,
		Timeout: 5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewFromDockerfile returned error: %v", err)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	want := "context/.dockerignore,context/Dockerfile,context/keep.log,context/secrets/public.pem,context/src/build/kept,dockerfile/Dockerfile"
	if got := uploadedNames(t, fake.uploadedTar); got != want {
		t.Errorf("uploaded files = %s, want %s", got, want)
	}
}

func TestSnapshotNewFromDockerfilePrefersDockerfileSpecificIgnore(t *testing.T) {
	contextDir := t.TempDir()
	writeTestFiles(t, contextDir, map[string]string{
		".dockerignore":                      "a.txt\n",
		"docker/app.Dockerfile":              "FROM scratch\n",
		"docker/app.Dockerfile.dockerignore": "b.txt\n",
		"a.txt":                              "a\n",
		"b.txt":                              "b\n",
	})

	fake := &fakeDockerfileBuilder{t: t}
	client := newDockerfileTestClient(fake.serve())

	_, err := client.Sandboxes.Snapshots.NewFromDockerfile(context.Background(), langsmith.SnapshotDockerfileParams{
		Name:       "my-snap",
		Context:    contextDir,
		Dockerfile: "docker/app.Dockerfile",
		Timeout:    5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewFromDockerfile returned error: %v", err)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	files := tarFileNames(t, fake.uploadedTar)
	if _, ok := files["context/b.txt"]; ok {
		t.Error("context/b.txt uploaded although app.Dockerfile.dockerignore excludes it")
	}
	if _, ok := files["context/a.txt"]; !ok {
		t.Error("context/a.txt left out although only the context's .dockerignore excludes it")
	}
	if files["dockerfile/app.Dockerfile.dockerignore"] != "b.txt\n" {
		t.Errorf("Dockerfile-specific ignore file not uploaded beside the Dockerfile: %v", files)
	}
}

func TestSnapshotNewFromDockerfileFollowsSymlinkedContext(t *testing.T) {
	root := t.TempDir()
	writeTestFiles(t, root, map[string]string{
		"app/Dockerfile":  "FROM scratch\n",
		"app/src/main.py": "print('hi')\n",
	})
	linkedContext := filepath.Join(root, "current")
	if err := os.Symlink(filepath.Join(root, "app"), linkedContext); err != nil {
		t.Fatal(err)
	}

	fake := &fakeDockerfileBuilder{t: t}
	client := newDockerfileTestClient(fake.serve())

	_, err := client.Sandboxes.Snapshots.NewFromDockerfile(context.Background(), langsmith.SnapshotDockerfileParams{
		Name:       "my-snap",
		Context:    linkedContext,
		Dockerfile: filepath.Join(linkedContext, "Dockerfile"),
		Timeout:    5 * time.Second,
	})
	if err != nil {
		t.Fatalf("NewFromDockerfile returned error: %v", err)
	}

	fake.mu.Lock()
	defer fake.mu.Unlock()
	files := tarFileNames(t, fake.uploadedTar)
	if files["context/Dockerfile"] != "FROM scratch\n" || files["context/src/main.py"] == "" {
		t.Errorf("uploaded context from a symlinked directory = %v, want its Dockerfile and src/main.py", files)
	}
	if !strings.Contains(fake.buildCommand, "'filename=Dockerfile'") {
		t.Errorf("build command does not name the Dockerfile:\n%s", fake.buildCommand)
	}
}
