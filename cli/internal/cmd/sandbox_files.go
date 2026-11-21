package cmd

import (
	"context"
	"fmt"
	"io"
	"os"
	"path"
	"path/filepath"
	"strings"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/cmdutil"
	"github.com/langchain-ai/langsmith-go/cli/internal/structured"
	"github.com/spf13/cobra"
)

// sandboxFileRef is one side of a cp: a local path, or a path in a sandbox.
type sandboxFileRef struct {
	Sandbox string
	Path    string
}

func (r sandboxFileRef) remote() bool { return r.Sandbox != "" }

// parseSandboxFileRef reads NAME:/absolute/path as a sandbox path and anything
// else as a local one. Sandbox paths are always absolute, so a local path that
// contains a colon is not mistaken for one unless a slash follows the colon.
func parseSandboxFileRef(arg string) sandboxFileRef {
	if filepath.VolumeName(arg) != "" {
		return sandboxFileRef{Path: arg}
	}
	name, p, ok := strings.Cut(arg, ":")
	if !ok || name == "" || strings.ContainsAny(name, `/\`) || !strings.HasPrefix(p, "/") {
		return sandboxFileRef{Path: arg}
	}
	return sandboxFileRef{Sandbox: name, Path: p}
}

var sandboxCopyCommand = structured.Command[struct{}]{
	Use:   "cp <src> <dst>",
	Short: "Copy a file to or from a sandbox",
	Long: `Copy a single file between your machine and a sandbox.

Name the sandbox side NAME:/absolute/path. Use - as the local side to read from
stdin or write to stdout. A destination that is an existing local directory,
or a sandbox path ending in /, keeps the source's file name. Parent directories
inside the sandbox are created as needed. Directories are not copied; archive
them first, for example with tar. On Windows, a one-letter name followed by a
colon is a drive, so C:/report.pdf is a local path.

Examples:
  langsmith sandboxes cp ./config.yaml my-vm:/srv/app/config.yaml
  langsmith sandboxes cp ./build.tar.gz my-vm:/tmp/
  langsmith sandboxes cp my-vm:/tmp/report.pdf ./report.pdf
  langsmith sandboxes cp my-vm:/var/log/app.log - | tail
  tar czf - src | langsmith sandboxes cp - my-vm:/tmp/src.tar.gz`,
	Args: cobra.ExactArgs(2),
	Input: func(cmd *cobra.Command) struct{} {
		cmd.Flags().String("jq", "", cmdutil.JQFlagUsage)
		return struct{}{}
	},
	CustomOutput: true,
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		src, dst := parseSandboxFileRef(args[0]), parseSandboxFileRef(args[1])
		switch {
		case src.remote() && dst.remote():
			return nil, fmt.Errorf("copy between two sandboxes is not supported; copy through a local file")
		case !src.remote() && !dst.remote():
			return nil, fmt.Errorf("one of <src> and <dst> must be a sandbox path, NAME:/absolute/path")
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if dst.remote() {
			content, err := readLocalSource(src.Path)
			if err != nil {
				return nil, err
			}
			target := dst.Path
			if strings.HasSuffix(target, "/") {
				if src.Path == "-" {
					return nil, fmt.Errorf("name the destination file when copying from stdin")
				}
				target = path.Join(target, filepath.Base(src.Path))
			}
			if err := c.SDK.Sandboxes.Boxes.WriteFile(ctx, dst.Sandbox, target, content); err != nil {
				return nil, fmt.Errorf("writing %s:%s: %w", dst.Sandbox, target, err)
			}
			return nil, structured.Render(cmd, sandboxMessage{
				Name:    dst.Sandbox,
				Message: fmt.Sprintf("Copied %d bytes to %s:%s", len(content), dst.Sandbox, target),
			}, structured.Template("{{.Message}}\n"))
		}

		content, err := c.SDK.Sandboxes.Boxes.ReadFile(ctx, src.Sandbox, src.Path)
		if err != nil {
			return nil, fmt.Errorf("reading %s:%s: %w", src.Sandbox, src.Path, err)
		}
		if dst.Path == "-" {
			_, err := os.Stdout.Write(content)
			return nil, err
		}
		target := dst.Path
		if info, err := os.Stat(target); err == nil && info.IsDir() {
			target = filepath.Join(target, path.Base(src.Path))
		}
		if err := os.WriteFile(target, content, 0o644); err != nil {
			return nil, fmt.Errorf("writing %s: %w", target, err)
		}
		return nil, structured.Render(cmd, sandboxMessage{
			Name:    src.Sandbox,
			Message: fmt.Sprintf("Copied %d bytes to %s", len(content), target),
		}, structured.Template("{{.Message}}\n"))
	},
}

func readLocalSource(p string) ([]byte, error) {
	if p == "-" {
		return io.ReadAll(os.Stdin)
	}
	info, err := os.Stat(p)
	if err != nil {
		return nil, err
	}
	if info.IsDir() {
		return nil, fmt.Errorf("%s is a directory; cp copies single files", p)
	}
	return os.ReadFile(p)
}

type sandboxListFilesInput struct {
	Pattern string
	Limit   int
}

type sandboxListFilesResult struct {
	Matches   []langsmith.SandboxFileInfo `json:"matches"`
	Truncated bool                        `json:"truncated"`
}

var sandboxListFilesCommand = structured.Command[*sandboxListFilesInput]{
	Use:   "ls <name> [path]",
	Short: "List files in a sandbox",
	Long: `List the files under an absolute path in a sandbox (default /).

Without --pattern this lists the directory's own entries. --pattern matches
paths relative to the directory: "*.py" stays in it, "**/*.py" recurses.

Examples:
  langsmith sandboxes ls my-vm
  langsmith sandboxes ls my-vm /srv/app
  langsmith sandboxes ls my-vm /srv/app --pattern '**/*.py' --limit 200`,
	Args: cobra.RangeArgs(1, 2),
	Input: func(cmd *cobra.Command) *sandboxListFilesInput {
		in := &sandboxListFilesInput{Pattern: "*"}
		cmd.Flags().StringVar(&in.Pattern, "pattern", in.Pattern, "Glob pattern relative to the path; ** recurses")
		cmd.Flags().IntVar(&in.Limit, "limit", in.Limit, "Maximum number of entries (0 for the server default)")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxListFilesInput, args []string) (any, error) {
		dir := "/"
		if len(args) > 1 {
			dir = args[1]
		}
		if !strings.HasPrefix(dir, "/") {
			return nil, fmt.Errorf("path must be absolute (got %q)", dir)
		}
		if in.Limit < 0 {
			return nil, fmt.Errorf("--limit must not be negative")
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		params := langsmith.SandboxGlobParams{
			Pattern: langsmith.F(in.Pattern),
			Path:    langsmith.F(dir),
		}
		if in.Limit > 0 {
			params.Limit = langsmith.F(int64(in.Limit))
		}
		resp, err := c.SDK.Sandboxes.Boxes.Glob(ctx, args[0], params)
		if err != nil {
			return nil, fmt.Errorf("listing files: %w", err)
		}
		return sandboxListFilesResult{Matches: resp.Matches, Truncated: resp.Truncated}, nil
	},
	Render: structured.Template(`{{range .Matches}}{{if .IsDir}}d{{else}}-{{end}} {{printf "%12d" .SizeBytes}} {{formatTime .ModifiedAt}} {{.Path}}
{{end}}{{if .Truncated}}(truncated; narrow --pattern or raise --limit)
{{end}}`),
}
