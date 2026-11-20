package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/cmdutil"
	"github.com/langchain-ai/langsmith-go/cli/internal/structured"
	"github.com/spf13/cobra"
)

var sandboxSnapshotCommand = structured.Parent{
	Use:   "snapshot",
	Short: "Manage sandbox snapshots (ext4 rootfs images)",
	Long: `Manage sandbox snapshots — ext4 rootfs images built from Docker images
or captured from running sandboxes. Snapshots are used to create new sandboxes.

A snapshot is named Docker-style: "name:tag" is a mutable pointer to one
immutable snapshot, and a bare name means "name:latest". Publishing over an
existing name:tag moves the pointer; the old snapshot stays addressable by id.

Examples:
  langsmith sandboxes snapshot list
  langsmith sandboxes snapshot build my-snap --docker-image ubuntu:24.04
  langsmith sandboxes snapshot capture my-snap:v2 --box my-vm
  langsmith sandboxes snapshot get my-snap:v2
  langsmith sandboxes snapshot delete <id>`,
	Children: []func() *cobra.Command{
		snapshotListCommand.Cobra,
		snapshotBuildCommand.Cobra,
		snapshotCaptureCommand.Cobra,
		snapshotGetCommand.Cobra,
		snapshotDeleteCommand.Cobra,
	},
}

var snapshotListCommand = structured.Command[struct{}]{
	Use:   "list",
	Short: "List all snapshots",
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		var snapshots []langsmith.SnapshotResponse
		pager := c.SDK.Sandboxes.Snapshots.ListAutoPaging(ctx, langsmith.SandboxSnapshotListParams{})
		for pager.Next() {
			snapshots = append(snapshots, pager.Current())
		}
		if err := pager.Err(); err != nil {
			return nil, fmt.Errorf("listing snapshots: %w", err)
		}
		return snapshots, nil
	},
	Render: structured.Table{
		Title: "Snapshots",
		Rows:  ".",
		Columns: []structured.Column{
			{Header: "ID", Template: "{{.ID}}"},
			{Header: "Name", Template: "{{.Name}}"},
			{Header: "Tags", Template: "{{joinOrDash .Tags}}"},
			{Header: "Image", Template: "{{dash .DockerImage}}"},
			{Header: "Status", Template: "{{.Status}}"},
			{Header: "Size", Template: "{{formatBytesOrDash .FsUsedBytes}}"},
			{Header: "Created", Template: "{{formatTime .CreatedAt}}"},
		},
	},
}

type snapshotBuildInput struct {
	DockerImage string
	Capacity    string
	RegistryID  string
	Description string

	Dockerfile string
	Context    string
	BuildArgs  []string
	Target     string
	VCPUs      int
	Memory     string
	Timeout    time.Duration
}

var snapshotBuildRender = structured.Template(`ID:          {{.ID}}
Name:        {{.Name}}
Tags:        {{joinOrDash .Tags}}
Image:       {{dash .DockerImage}}
Description: {{dash .Description}}
Status:      {{.Status}}
Size:        {{formatBytesOrDash .FsUsedBytes}}
Created:     {{formatTime .CreatedAt}}
`)

var snapshotBuildCommand = structured.Command[*snapshotBuildInput]{
	Use:   "build <name[:tag]>",
	Short: "Build a snapshot from a Docker image or a local Dockerfile",
	Long: `Build a snapshot from a Docker image (--docker-image) or a local Dockerfile
(--dockerfile). Without either flag, the Dockerfile in the --context directory
is used when there is one.

With --docker-image the tag is applied once the build finishes. Omit it and the
snapshot inherits the Docker image's tag (ubuntu:24.04 becomes 24.04), else
"latest". The command returns while the snapshot is still building.

A Dockerfile build follows docker build's context rules: --context defaults to
the current directory, --dockerfile is relative to the current directory and
defaults to <context>/Dockerfile, and files that .dockerignore excludes are not
uploaded (a <Dockerfile>.dockerignore takes precedence). The build runs in a
temporary sandbox with BuildKit and is exported as the snapshot, so no local
Docker is needed. Build output streams to stderr, and the command waits until
the snapshot is ready. An omitted tag means "latest". --capacity, when set,
sizes both the builder's disk and the snapshot; it defaults to the builder's
disk size and cannot be smaller. --vcpus and --memory size the builder, and an
extra vCPU speeds up cold builds.

	Examples:
	  langsmith sandboxes snapshot build my-snap
	  langsmith sandboxes snapshot build my-snap --docker-image ubuntu:24.04
	  langsmith sandboxes snapshot build my-snap:v2 --docker-image ubuntu:24.04
	  langsmith sandboxes snapshot build my-snap --docker-image ubuntu:24.04 --capacity 8gb
	  langsmith sandboxes snapshot build my-snap --docker-image ubuntu:24.04 --description "Python 3.12 with uv and ruff"
	  langsmith sandboxes snapshot build my-snap --dockerfile Dockerfile
	  langsmith sandboxes snapshot build my-snap:v2 --dockerfile docker/app.Dockerfile --context . --build-arg VERSION=2 --target runtime --vcpus 2`,
	Args: cobra.ExactArgs(1),
	Input: func(cmd *cobra.Command) *snapshotBuildInput {
		in := &snapshotBuildInput{
			Capacity: "4gb",
			Context:  ".",
			Timeout:  10 * time.Minute,
		}
		cmd.Flags().StringVar(&in.DockerImage, "docker-image", in.DockerImage, "Docker image to build from instead of a Dockerfile")
		cmd.Flags().StringVar(&in.Capacity, "capacity", in.Capacity, "Filesystem capacity with unit (e.g. 4gb, 8gb); with --dockerfile, defaults to the builder's disk size")
		cmd.Flags().StringVar(&in.RegistryID, "registry-id", in.RegistryID, "Registry ID for private images (--docker-image only)")
		cmd.Flags().StringVar(&in.Description, "description", in.Description, "What the snapshot's image can do, for handing to an agent (max 1024 chars)")
		cmd.Flags().StringVar(&in.Dockerfile, "dockerfile", in.Dockerfile, "Local Dockerfile to build from, relative to the current directory (default <context>/Dockerfile)")
		cmd.Flags().StringVar(&in.Context, "context", in.Context, "Build context directory for a Dockerfile build")
		cmd.Flags().StringArrayVar(&in.BuildArgs, "build-arg", nil, "Build arg for --dockerfile (KEY=VALUE, or KEY to take it from the environment; repeatable)")
		cmd.Flags().StringVar(&in.Target, "target", in.Target, "Dockerfile stage to build (--dockerfile only)")
		cmd.Flags().IntVar(&in.VCPUs, "vcpus", in.VCPUs, "vCPUs for the --dockerfile builder sandbox")
		cmd.Flags().StringVar(&in.Memory, "memory", in.Memory, "Memory with unit for the --dockerfile builder sandbox (e.g. 8gb)")
		cmd.Flags().DurationVar(&in.Timeout, "timeout", in.Timeout, "Time limit for each --dockerfile step: builder start, upload, build, and snapshot export")
		cmd.MarkFlagsMutuallyExclusive("docker-image", "dockerfile")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *snapshotBuildInput, args []string) (any, error) {
		name, tag := splitSnapshotRef(args[0])
		if in.DockerImage == "" && in.Dockerfile == "" {
			defaultDockerfile := filepath.Join(in.Context, "Dockerfile")
			if info, err := os.Stat(defaultDockerfile); err != nil || !info.Mode().IsRegular() {
				return nil, fmt.Errorf("no Dockerfile in %s: pass --dockerfile, --context, or --docker-image", in.Context)
			}
			in.Dockerfile = defaultDockerfile
		}

		capBytes, err := parseByteSize(in.Capacity)
		if err != nil {
			return nil, fmt.Errorf("invalid --capacity: %w", err)
		}

		if in.Dockerfile != "" {
			if in.RegistryID != "" {
				return nil, fmt.Errorf("--registry-id applies to --docker-image only")
			}
			if !cmd.Flags().Changed("capacity") {
				capBytes = 0
			}
			params, err := snapshotDockerfileParams(cmd, in, name, tag, capBytes)
			if err != nil {
				return nil, err
			}
			c, err := cmdutil.GetClient(cmd)
			if err != nil {
				return nil, err
			}
			resp, err := c.SDK.Sandboxes.Snapshots.NewFromDockerfile(ctx, params)
			if err != nil {
				return nil, fmt.Errorf("building snapshot from Dockerfile: %w", err)
			}
			return resp, nil
		}

		for _, flag := range []string{"context", "build-arg", "target", "vcpus", "memory", "timeout"} {
			if cmd.Flags().Changed(flag) {
				return nil, fmt.Errorf("--%s applies to --dockerfile only", flag)
			}
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		params := langsmith.SandboxSnapshotNewParams{
			Name:            langsmith.F(name),
			DockerImage:     langsmith.F(in.DockerImage),
			FsCapacityBytes: langsmith.F(capBytes),
		}
		if tag != "" {
			params.Tag = langsmith.F(tag)
		}
		if in.RegistryID != "" {
			params.RegistryID = langsmith.F(in.RegistryID)
		}
		if in.Description != "" {
			params.Description = langsmith.F(in.Description)
		}

		resp, err := c.SDK.Sandboxes.Snapshots.New(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("building snapshot: %w", err)
		}

		return resp, nil
	},
	Render: snapshotBuildRender,
}

// snapshotDockerfileParams turns snapshot build flags into the SDK's
// Dockerfile build parameters. Build output streams to stderr so stdout stays
// parseable.
func snapshotDockerfileParams(cmd *cobra.Command, in *snapshotBuildInput, name, tag string, capBytes int64) (langsmith.SnapshotDockerfileParams, error) {
	params := langsmith.SnapshotDockerfileParams{
		Name:            name,
		Tag:             tag,
		Context:         in.Context,
		Target:          in.Target,
		Description:     in.Description,
		FsCapacityBytes: capBytes,
		Timeout:         in.Timeout,
		OnBuildLog: func(s string) {
			fmt.Fprint(cmd.ErrOrStderr(), s)
		},
	}
	// Like docker build -f, a relative --dockerfile names a path from the
	// current directory, not from the context.
	dockerfile, err := filepath.Abs(in.Dockerfile)
	if err != nil {
		return params, fmt.Errorf("invalid --dockerfile: %w", err)
	}
	params.Dockerfile = dockerfile
	if in.Timeout <= 0 {
		return params, fmt.Errorf("--timeout must be positive")
	}
	if in.VCPUs < 0 {
		return params, fmt.Errorf("--vcpus must be positive")
	}
	params.Vcpus = int64(in.VCPUs)
	if in.Memory != "" {
		memBytes, err := parseByteSize(in.Memory)
		if err != nil {
			return params, fmt.Errorf("invalid --memory: %w", err)
		}
		params.MemBytes = memBytes
	}
	buildArgs, err := parseBuildArgs(in.BuildArgs)
	if err != nil {
		return params, err
	}
	params.BuildArgs = buildArgs
	return params, nil
}

// parseBuildArgs parses repeated --build-arg values. Like docker build, a bare
// KEY takes its value from the environment and is skipped when unset.
func parseBuildArgs(values []string) (map[string]string, error) {
	if len(values) == 0 {
		return nil, nil
	}
	args := make(map[string]string, len(values))
	for _, v := range values {
		key, value, hasValue := strings.Cut(v, "=")
		if key == "" {
			return nil, fmt.Errorf("invalid --build-arg %q: expected KEY=VALUE or KEY", v)
		}
		if !hasValue {
			var ok bool
			if value, ok = os.LookupEnv(key); !ok {
				continue
			}
		}
		args[key] = value
	}
	return args, nil
}

type snapshotCaptureInput struct {
	BoxName     string
	Checkpoint  string
	Description string
}

var snapshotCaptureCommand = structured.Command[*snapshotCaptureInput]{
	Use:   "capture <name[:tag]>",
	Short: "Capture a snapshot from a running sandbox",
	Long: `Capture a snapshot from a sandbox VM. If --checkpoint is specified, uses
that existing checkpoint (no VM interaction needed). Otherwise creates a
fresh checkpoint from the running VM's current state.

An omitted tag means "latest".

	Examples:
	  langsmith sandboxes snapshot capture my-snap --box my-vm
	  langsmith sandboxes snapshot capture my-snap:2026081101 --box my-vm
	  langsmith sandboxes snapshot capture my-snap --box my-vm --checkpoint 2026-03-29T00:09:28Z
	  langsmith sandboxes snapshot capture my-snap --box my-vm --description "repo cloned, deps installed"`,
	Args: cobra.ExactArgs(1),
	Input: func(cmd *cobra.Command) *snapshotCaptureInput {
		in := &snapshotCaptureInput{}
		cmd.Flags().StringVar(&in.BoxName, "box", in.BoxName, "Sandbox name to capture from (required)")
		cmd.Flags().StringVar(&in.Checkpoint, "checkpoint", in.Checkpoint, "Checkpoint timestamp to use (omit for fresh checkpoint)")
		cmd.Flags().StringVar(&in.Description, "description", in.Description, "What the snapshot's image can do, for handing to an agent (max 1024 chars)")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *snapshotCaptureInput, args []string) (any, error) {
		name, tag := splitSnapshotRef(args[0])
		if in.BoxName == "" {
			return nil, fmt.Errorf("--box is required")
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		params := langsmith.SandboxBoxNewSnapshotParams{
			Name: langsmith.F(name),
		}
		if tag != "" {
			params.Tag = langsmith.F(tag)
		}
		if in.Checkpoint != "" {
			params.Checkpoint = langsmith.F(in.Checkpoint)
		}
		if in.Description != "" {
			params.Description = langsmith.F(in.Description)
		}

		resp, err := c.SDK.Sandboxes.Boxes.NewSnapshot(ctx, in.BoxName, params)
		if err != nil {
			return nil, fmt.Errorf("capturing snapshot: %w", err)
		}

		return resp, nil
	},
	Render: structured.Template(`ID:          {{.ID}}
Name:        {{.Name}}
Tags:        {{joinOrDash .Tags}}
Image:       {{dash .DockerImage}}
Description: {{dash .Description}}
Status:      {{.Status}}
Size:        {{formatBytesOrDash .FsUsedBytes}}
Created:     {{formatTime .CreatedAt}}
`),
}

var snapshotGetCommand = structured.Command[struct{}]{
	Use:   "get <snapshot-id|name[:tag]>",
	Short: "Get a snapshot by ID or Docker-style reference",
	Args:  cobra.ExactArgs(1),
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		resp, err := c.SDK.Sandboxes.Snapshots.Get(ctx, args[0])
		if err != nil {
			return nil, fmt.Errorf("getting snapshot: %w", err)
		}

		return resp, nil
	},
	Render: structured.Template(`ID:          {{.ID}}
Name:        {{.Name}}
Tags:        {{joinOrDash .Tags}}
Image:       {{dash .DockerImage}}
Description: {{dash .Description}}
Status:      {{.Status}}
Size:        {{formatBytesOrDash .FsUsedBytes}}
Created:     {{formatTime .CreatedAt}}
`),
}

var snapshotDeleteCommand = structured.Command[struct{}]{
	Use:   "delete <snapshot-id|name[:tag]>",
	Short: "Delete a snapshot by ID or Docker-style reference",
	Args:  cobra.ExactArgs(1),
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if err := c.SDK.Sandboxes.Snapshots.Delete(ctx, args[0]); err != nil {
			return nil, fmt.Errorf("deleting snapshot: %w", err)
		}

		return sandboxMessage{Name: args[0], Message: "Snapshot deleted."}, nil
	},
	Render: structured.Template(`{{.Message}}
`),
}

// splitSnapshotRef splits a Docker-style name[:tag] on the first colon. The
// server rejects a colon inside a name, so an empty tag leaves the default to it.
func splitSnapshotRef(ref string) (name, tag string) {
	name, tag, _ = strings.Cut(ref, ":")
	return name, tag
}

func parseByteSize(s string) (int64, error) {
	s = strings.TrimSpace(s)
	if s == "" {
		return 0, fmt.Errorf("empty size")
	}

	i := 0
	for i < len(s) && (s[i] >= '0' && s[i] <= '9' || s[i] == '.') {
		i++
	}
	if i == 0 {
		return 0, fmt.Errorf("invalid size %q: no number", s)
	}
	if i == len(s) {
		return 0, fmt.Errorf("invalid size %q: unit required (e.g. 512mb, 2gb)", s)
	}

	num, err := strconv.ParseFloat(s[:i], 64)
	if err != nil {
		return 0, fmt.Errorf("invalid size %q: %w", s, err)
	}

	unit := strings.ToLower(strings.TrimSpace(s[i:]))
	var multiplier float64
	switch unit {
	case "b":
		multiplier = 1
	case "k", "kb", "kib":
		multiplier = 1024
	case "m", "mb", "mib":
		multiplier = 1024 * 1024
	case "g", "gb", "gib":
		multiplier = 1024 * 1024 * 1024
	case "t", "tb", "tib":
		multiplier = 1024 * 1024 * 1024 * 1024
	default:
		return 0, fmt.Errorf("invalid size %q: unknown unit %q", s, unit)
	}

	bytes := int64(num * multiplier)

	const minBytes = 1024 * 1024
	const maxBytes = 1024 * 1024 * 1024 * 1024 * 1024
	if bytes < minBytes {
		return 0, fmt.Errorf("invalid size %q: must be at least 1mb", s)
	}
	if bytes > maxBytes {
		return 0, fmt.Errorf("invalid size %q: must be less than 1pb", s)
	}

	return bytes, nil
}

func loadJSONArg(s string) (json.RawMessage, error) {
	var data []byte
	if strings.HasPrefix(s, "@") {
		var err error
		data, err = os.ReadFile(s[1:])
		if err != nil {
			return nil, fmt.Errorf("reading file %s: %w", s[1:], err)
		}
	} else {
		data = []byte(s)
	}
	var raw json.RawMessage
	if err := json.Unmarshal(data, &raw); err != nil {
		return nil, fmt.Errorf("invalid JSON: %w", err)
	}
	return raw, nil
}
