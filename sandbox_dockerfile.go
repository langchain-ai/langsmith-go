package langsmith

import (
	"archive/tar"
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/google/uuid"
	"github.com/moby/patternmatcher"
	"github.com/moby/patternmatcher/ignorefile"

	"github.com/langchain-ai/langsmith-go/option"
)

const (
	defaultSnapshotDockerfileTimeout = 10 * time.Minute
	snapshotBuildLogTailBytes        = 4096
)

// SnapshotDockerfileParams configures [SandboxSnapshotService.NewFromDockerfile].
type SnapshotDockerfileParams struct {
	// Name of the snapshot to publish. Required.
	Name string
	// Tag to publish the snapshot under within Name. Empty means "latest".
	Tag string
	// Dockerfile is the path to the Dockerfile, relative to Context unless
	// absolute. It may be outside Context. Empty means "Dockerfile".
	Dockerfile string
	// Context is the local build context directory. Empty means the current
	// working directory. Files excluded by the context's .dockerignore are not
	// uploaded, and a <Dockerfile>.dockerignore beside the Dockerfile takes
	// precedence over it, as in docker build.
	Context string
	// BuildArgs are passed to the build as Docker build args.
	BuildArgs map[string]string
	// Target is the Dockerfile stage to build. Empty builds the last stage.
	Target string
	// Description says what the snapshot's image can do. At most 1024 characters.
	Description string
	// FsCapacityBytes sizes both the builder sandbox's disk and the snapshot's
	// filesystem, and cannot be smaller than the builder's default disk. Zero
	// keeps the builder's default disk size for both.
	FsCapacityBytes int64
	// Vcpus sizes the temporary builder sandbox. BuildKit and its layer copies
	// contend for a single core by default, so an extra vCPU can cut a cold
	// build's wall time substantially. Zero uses the server default.
	Vcpus int64
	// MemBytes sizes the temporary builder sandbox's memory. Zero uses the
	// server default.
	MemBytes int64
	// Timeout bounds each step: starting the builder, uploading the context,
	// running the build, and waiting for the snapshot. Zero means 10 minutes.
	Timeout time.Duration
	// OnBuildLog, when set, receives build output as it streams.
	OnBuildLog func(string)
}

// NewFromDockerfile builds a snapshot from a local Dockerfile.
//
// It starts a temporary builder sandbox, uploads the build context, builds the
// image there with BuildKit, and exports the image as the snapshot. No local
// Docker daemon or registry is needed. The builder sandbox is deleted
// afterwards, on a best-effort basis.
func (r *SandboxSnapshotService) NewFromDockerfile(ctx context.Context, params SnapshotDockerfileParams, opts ...option.RequestOption) (*SnapshotResponse, error) {
	if params.Name == "" {
		return nil, errors.New("langsmith: snapshot name is required")
	}
	dockerfile := params.Dockerfile
	if dockerfile == "" {
		dockerfile = "Dockerfile"
	}
	contextDir := params.Context
	if contextDir == "" {
		contextDir = "."
	}
	inputs, err := resolveDockerBuildInputs(dockerfile, contextDir)
	if err != nil {
		return nil, err
	}
	inputTar, err := makeDockerBuildTar(inputs)
	if err != nil {
		return nil, err
	}
	timeout := params.Timeout
	if timeout == 0 {
		timeout = defaultSnapshotDockerfileTimeout
	}
	timeoutSeconds := int64((timeout + time.Second - 1) / time.Second)

	boxes := NewSandboxBoxService(r.Options...)
	builderParams := SandboxBoxNewParams{
		Name: F("snapshot-builder-" + shortRandomHex()),
	}
	if params.Vcpus > 0 {
		builderParams.Vcpus = F(params.Vcpus)
	}
	if params.MemBytes > 0 {
		builderParams.MemBytes = F(params.MemBytes)
	}
	if params.FsCapacityBytes > 0 {
		builderParams.FsCapacityBytes = F(params.FsCapacityBytes)
	}
	builder, err := boxes.NewSandbox(ctx, builderParams, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: create snapshot builder sandbox: %w", err)
	}
	defer func() {
		_ = builder.Delete(context.WithoutCancel(ctx), opts...)
	}()
	if builder.Status != "ready" {
		ready, err := boxes.Wait(ctx, builder.Name, SandboxWaitParams{Timeout: timeout}, opts...)
		if err != nil {
			return nil, fmt.Errorf("langsmith: wait for snapshot builder sandbox: %w", err)
		}
		builder.applyGetResponse(ready)
	}

	// Stage the build on the capacity-backed root filesystem, not /tmp. Inside
	// the sandbox /tmp is a RAM-backed tmpfs that FsCapacityBytes does not size,
	// and BuildKit's native snapshotter writes a full copy of every layer under
	// its root, so a /tmp build exhausts guest RAM.
	buildRoot := "/var/lib/langsmith-build/" + shortRandomHex()
	remoteInput := buildRoot + "/input"
	remoteTar := buildRoot + "/input.tar"
	imageRef := "langsmith-snapshot-build:" + strings.ReplaceAll(uuid.NewString(), "-", "")

	writeOpts := append([]option.RequestOption{option.WithRequestTimeout(timeout)}, opts...)
	if err := builder.WriteFile(ctx, remoteTar, inputTar, writeOpts...); err != nil {
		return nil, fmt.Errorf("langsmith: upload build context: %w", err)
	}
	extract, err := builder.Run(ctx, SandboxBoxRunParams{
		Command: F(fmt.Sprintf("rm -rf %[1]s && mkdir -p %[1]s && tar -xf %[2]s -C %[1]s",
			shellQuote(remoteInput), shellQuote(remoteTar))),
		Timeout: F(timeoutSeconds),
	}, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: unpack build context: %w", err)
	}
	if !extract.Success() {
		return nil, &SandboxResourceCreationError{
			ResourceType: "snapshot",
			ResourceID:   params.Name,
			Message:      "unpack build context: " + outputTail(extract.Stdout+extract.Stderr),
		}
	}

	var callbacks SandboxCommandCallbacks
	if params.OnBuildLog != nil {
		callbacks.OnStdout = params.OnBuildLog
		callbacks.OnStderr = params.OnBuildLog
	}
	build, err := builder.RunWithCallbacks(ctx, SandboxCommandStartParams{
		Command: F(dockerfileBuildCommand(dockerfileBuild{
			remoteContext:  path.Join(remoteInput, buildTarContextDir),
			dockerfileDir:  path.Join(remoteInput, buildTarDockerfileDir),
			dockerfileName: inputs.dockerfileName,
			imageRef:       imageRef,
			buildkitRoot:   buildRoot + "/buildkit-root",
			buildkitRun:    buildRoot + "/buildkit-run",
			buildArgs:      params.BuildArgs,
			target:         params.Target,
		})),
		TimeoutSeconds: F(timeoutSeconds),
	}, callbacks, opts...)
	if err != nil {
		return nil, fmt.Errorf("langsmith: run Dockerfile build: %w", err)
	}
	if !build.Success() {
		return nil, &SandboxResourceCreationError{
			ResourceType: "snapshot",
			ResourceID:   params.Name,
			Message:      fmt.Sprintf("Dockerfile build exited with code %d: %s", build.ExitCode, outputTail(build.Stdout+build.Stderr)),
		}
	}

	capture := SandboxBoxNewSnapshotParams{
		Name:        F(params.Name),
		DockerImage: F(imageRef),
	}
	if params.Tag != "" {
		capture.Tag = F(params.Tag)
	}
	if params.Description != "" {
		capture.Description = F(params.Description)
	}
	if params.FsCapacityBytes > 0 {
		capture.FsCapacityBytes = F(params.FsCapacityBytes)
	}
	return builder.CaptureSnapshotAndWait(ctx, capture, SnapshotWaitParams{Timeout: timeout}, opts...)
}

// Top-level directories of the uploaded build input: the filtered context and,
// separately, the Dockerfile, which may live outside the context.
const (
	buildTarContextDir    = "context"
	buildTarDockerfileDir = "dockerfile"
)

// dockerBuildInputs are the local files a Dockerfile build uploads.
type dockerBuildInputs struct {
	// contextPath is the context directory with symlinks resolved, so the
	// archive walk starts at a real directory.
	contextPath string
	// dockerfilePath is the Dockerfile to read.
	dockerfilePath string
	// dockerfileName is the Dockerfile's base name as given, which BuildKit
	// uses to find <name>.dockerignore.
	dockerfileName string
	// ignorePath is the ignore file that filters the context, or empty.
	ignorePath string
	// dockerfileIgnore reports whether ignorePath is the Dockerfile-specific
	// <name>.dockerignore rather than the context's .dockerignore.
	dockerfileIgnore bool
}

// resolveDockerBuildInputs resolves the context directory, the Dockerfile, and
// the ignore file the way docker build does: <Dockerfile>.dockerignore beside
// the Dockerfile wins over <context>/.dockerignore.
func resolveDockerBuildInputs(dockerfile string, contextDir string) (dockerBuildInputs, error) {
	contextAbs, err := filepath.Abs(contextDir)
	if err != nil {
		return dockerBuildInputs{}, fmt.Errorf("langsmith: resolve build context: %w", err)
	}
	contextPath, err := filepath.EvalSymlinks(contextAbs)
	if err != nil {
		return dockerBuildInputs{}, fmt.Errorf("langsmith: build context: %w", err)
	}
	info, err := os.Stat(contextPath)
	if err != nil {
		return dockerBuildInputs{}, fmt.Errorf("langsmith: build context: %w", err)
	}
	if !info.IsDir() {
		return dockerBuildInputs{}, fmt.Errorf("langsmith: build context must be a directory: %s", contextPath)
	}
	dockerfilePath := dockerfile
	if !filepath.IsAbs(dockerfilePath) {
		dockerfilePath = filepath.Join(contextAbs, dockerfilePath)
	}
	info, err = os.Stat(dockerfilePath)
	if err != nil {
		return dockerBuildInputs{}, fmt.Errorf("langsmith: dockerfile: %w", err)
	}
	if !info.Mode().IsRegular() {
		return dockerBuildInputs{}, fmt.Errorf("langsmith: dockerfile must be a file: %s", dockerfilePath)
	}
	inputs := dockerBuildInputs{
		contextPath:    contextPath,
		dockerfilePath: dockerfilePath,
		dockerfileName: filepath.Base(dockerfilePath),
	}
	for i, candidate := range []string{dockerfilePath + ".dockerignore", filepath.Join(contextPath, ".dockerignore")} {
		if info, err := os.Stat(candidate); err == nil && info.Mode().IsRegular() {
			inputs.ignorePath = candidate
			inputs.dockerfileIgnore = i == 0
			break
		}
	}
	return inputs, nil
}

// loadDockerignore parses an ignore file with Docker's own pattern rules.
func loadDockerignore(ignorePath string) (*patternmatcher.PatternMatcher, error) {
	if ignorePath == "" {
		return nil, nil
	}
	f, err := os.Open(ignorePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	patterns, err := ignorefile.ReadAll(f)
	if err != nil {
		return nil, fmt.Errorf("read %s: %w", ignorePath, err)
	}
	if len(patterns) == 0 {
		return nil, nil
	}
	return patternmatcher.New(patterns)
}

// makeDockerBuildTar archives the build input: the context under context/,
// without the files .dockerignore excludes, and the Dockerfile under
// dockerfile/ with its Dockerfile-specific ignore file, if any, so BuildKit
// applies the same patterns. Symlinks are archived as links, not followed;
// other special files are skipped.
func makeDockerBuildTar(inputs dockerBuildInputs) ([]byte, error) {
	var buf bytes.Buffer
	tw := tar.NewWriter(&buf)
	if err := writeDockerContext(tw, inputs); err != nil {
		return nil, fmt.Errorf("langsmith: archive build context: %w", err)
	}
	if err := writeDockerfileInput(tw, inputs); err != nil {
		return nil, fmt.Errorf("langsmith: archive dockerfile: %w", err)
	}
	if err := tw.Close(); err != nil {
		return nil, fmt.Errorf("langsmith: archive build context: %w", err)
	}
	return buf.Bytes(), nil
}

func writeDockerContext(tw *tar.Writer, inputs dockerBuildInputs) error {
	pm, err := loadDockerignore(inputs.ignorePath)
	if err != nil {
		return err
	}
	contextPath := inputs.contextPath
	return filepath.WalkDir(contextPath, func(p string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		rel, err := filepath.Rel(contextPath, p)
		if err != nil {
			return err
		}
		action := archiveEntry
		if rel != "." && pm != nil {
			if action, err = dockerignoreAction(pm, rel, d.IsDir()); err != nil {
				return err
			}
		}
		switch action {
		case omitTree:
			return filepath.SkipDir
		case omitEntry:
			return nil
		}
		name := buildTarContextDir
		if rel != "." {
			name = path.Join(buildTarContextDir, filepath.ToSlash(rel))
		}
		return writeTarEntry(tw, p, name)
	})
}

type contextEntryAction int

const (
	archiveEntry contextEntryAction = iota
	// omitEntry leaves out the entry but still walks a directory's children.
	omitEntry
	// omitTree leaves out a directory and everything below it.
	omitTree
)

// dockerignoreAction applies Docker's context-walk rules to rel. An excluded
// directory is still walked when an exception pattern (!) could re-include
// something inside it.
func dockerignoreAction(pm *patternmatcher.PatternMatcher, rel string, isDir bool) (contextEntryAction, error) {
	excluded, err := pm.MatchesOrParentMatches(rel)
	if err != nil || !excluded {
		return archiveEntry, err
	}
	if !isDir {
		return omitEntry, nil
	}
	if !pm.Exclusions() {
		return omitTree, nil
	}
	dirPrefix := rel + string(filepath.Separator)
	for _, pattern := range pm.Patterns() {
		if pattern.Exclusion() && strings.HasPrefix(pattern.String()+string(filepath.Separator), dirPrefix) {
			return omitEntry, nil
		}
	}
	return omitTree, nil
}

func writeDockerfileInput(tw *tar.Writer, inputs dockerBuildInputs) error {
	if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeDir, Name: buildTarDockerfileDir + "/", Mode: 0o755}); err != nil {
		return err
	}
	if err := writeTarFile(tw, inputs.dockerfilePath, path.Join(buildTarDockerfileDir, inputs.dockerfileName)); err != nil {
		return err
	}
	if inputs.dockerfileIgnore {
		return writeTarFile(tw, inputs.ignorePath, path.Join(buildTarDockerfileDir, inputs.dockerfileName+".dockerignore"))
	}
	return nil
}

// writeTarFile archives the regular file at p, following symlinks, as name.
func writeTarFile(tw *tar.Writer, p, name string) error {
	data, err := os.ReadFile(p)
	if err != nil {
		return err
	}
	if err := tw.WriteHeader(&tar.Header{Typeflag: tar.TypeReg, Name: name, Mode: 0o644, Size: int64(len(data))}); err != nil {
		return err
	}
	_, err = tw.Write(data)
	return err
}

// writeTarEntry archives p as name without following a symlink, skipping
// special files.
func writeTarEntry(tw *tar.Writer, p, name string) error {
	info, err := os.Lstat(p)
	if err != nil {
		return err
	}
	link := ""
	switch {
	case info.Mode()&fs.ModeSymlink != 0:
		if link, err = os.Readlink(p); err != nil {
			return err
		}
	case info.IsDir(), info.Mode().IsRegular():
	default:
		return nil
	}
	hdr, err := tar.FileInfoHeader(info, link)
	if err != nil {
		return err
	}
	hdr.Name = name
	if info.IsDir() {
		hdr.Name += "/"
	}
	hdr.Uname, hdr.Gname = "", ""
	if err := tw.WriteHeader(hdr); err != nil {
		return err
	}
	if !info.Mode().IsRegular() {
		return nil
	}
	f, err := os.Open(p)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = io.Copy(tw, f)
	return err
}

type dockerfileBuild struct {
	remoteContext  string
	dockerfileDir  string
	dockerfileName string
	imageRef       string
	buildkitRoot   string
	buildkitRun    string
	buildArgs      map[string]string
	target         string
}

// dockerfileBuildCommand returns a bash script that starts buildkitd, builds
// the Dockerfile, and loads the result into the sandbox's Docker daemon as
// imageRef.
func dockerfileBuildCommand(b dockerfileBuild) string {
	socketAddr := "unix://" + path.Join(b.buildkitRun, "buildkitd.sock")
	logPath := path.Join(b.buildkitRun, "buildkitd.log")

	buildctl := []string{
		"buildctl", "--addr", socketAddr, "build",
		"--progress=plain",
		"--frontend", "dockerfile.v0",
		"--local", "context=" + b.remoteContext,
		"--local", "dockerfile=" + b.dockerfileDir,
		"--opt", "filename=" + b.dockerfileName,
		"--output", "type=docker,name=" + b.imageRef,
	}
	if b.target != "" {
		buildctl = append(buildctl, "--opt", "target="+b.target)
	}
	keys := make([]string, 0, len(b.buildArgs))
	for k := range b.buildArgs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		buildctl = append(buildctl, "--opt", "build-arg:"+k+"="+b.buildArgs[k])
	}
	quoted := make([]string, len(buildctl))
	for i, arg := range buildctl {
		quoted[i] = shellQuote(arg)
	}

	lines := []string{
		"set -euo pipefail",
		fmt.Sprintf("mkdir -p %s %s", shellQuote(b.buildkitRoot), shellQuote(b.buildkitRun)),
		fmt.Sprintf("buildkitd --addr %s --root %s --oci-worker=true --containerd-worker=false "+
			"--oci-worker-snapshotter=native --oci-worker-binary buildkit-runc > %s 2>&1 &",
			shellQuote(socketAddr), shellQuote(b.buildkitRoot), shellQuote(logPath)),
		"buildkitd_pid=$!",
		`cleanup() { kill "$buildkitd_pid" >/dev/null 2>&1 || true; }`,
		"trap cleanup EXIT",
		"for i in $(seq 1 300); do",
		fmt.Sprintf("  if buildctl --addr %s debug workers >/dev/null 2>&1; then break; fi", shellQuote(socketAddr)),
		fmt.Sprintf(`  if ! kill -0 "$buildkitd_pid" >/dev/null 2>&1; then cat %s; exit 1; fi`, shellQuote(logPath)),
		fmt.Sprintf(`  if [ "$i" = 300 ]; then cat %s; exit 1; fi`, shellQuote(logPath)),
		"  sleep 0.1",
		"done",
		"for i in $(seq 1 300); do",
		"  if docker info >/dev/null 2>&1; then break; fi",
		`  if [ "$i" = 300 ]; then docker info; exit 1; fi`,
		"  sleep 0.1",
		"done",
		strings.Join(quoted, " ") + " | docker load",
		fmt.Sprintf("rm -rf %s || true", shellQuote(b.buildkitRoot)),
		"",
	}
	return strings.Join(lines, "\n")
}

func shellQuote(s string) string {
	return "'" + strings.ReplaceAll(s, "'", `'"'"'`) + "'"
}

func shortRandomHex() string {
	return strings.ReplaceAll(uuid.NewString(), "-", "")[:12]
}

func outputTail(s string) string {
	s = strings.TrimSpace(s)
	if len(s) > snapshotBuildLogTailBytes {
		s = "..." + s[len(s)-snapshotBuildLogTailBytes:]
	}
	return s
}
