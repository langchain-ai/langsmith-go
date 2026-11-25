package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"strings"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/output"
	"github.com/spf13/cobra"
)

func newHubPullCmd() *cobra.Command {
	var (
		dir       string
		commitRef string
		yes       bool
	)

	cmd := &cobra.Command{
		Use:   "pull [OWNER/]REPO[:COMMIT_OR_TAG] --dir PATH",
		Short: "Pull an agent or skill commit into a local directory",
		Args:  cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			owner, name, inlineRef, err := parseHubOwnerRepo(args[0])
			if err != nil {
				return err
			}
			ref := inlineRef
			if commitRef != "" {
				ref = commitRef
			}

			c := MustGetClient()
			ctx := context.Background()
			params := langsmith.RepoDirectoryListParams{}
			if ref != "" {
				params.Commit = langsmith.F(ref)
			}

			resp, err := c.SDK.Repos.Directories.List(ctx, owner, name, params)
			if err != nil {
				return fmt.Errorf("pulling %s/%s: %w", owner, name, err)
			}

			files, err := sdkFilesToHubFiles(resp.Files)
			if err != nil {
				return fmt.Errorf("decoding files for %s/%s: %w", owner, name, err)
			}
			written, linked, err := writeFilesToDirectory(cmd, dir, files, yes)
			if err != nil {
				return err
			}
			sort.Strings(written)
			sort.Strings(linked)

			absDir, _ := filepath.Abs(dir)
			out := map[string]any{
				"status":      "pulled",
				"owner":       owner,
				"repo":        name,
				"commit_id":   resp.CommitID,
				"commit_hash": resp.CommitHash,
				"dir":         absDir,
				"files":       written,
			}
			if len(linked) > 0 {
				out["linked_children"] = linked
			}
			if err := output.OutputJSON(out, ""); err != nil {
				ExitErrorf("%v", err)
			}
			return nil
		},
	}

	cmd.Flags().StringVar(&dir, "dir", "", "Destination directory (required; its contents are replaced)")
	cmd.Flags().StringVar(&commitRef, "commit", "", "Commit hash or tag name (overrides inline :ref)")
	addYesFlag(cmd, &yes, "Skip confirmation when replacing a non-empty directory")
	_ = cmd.MarkFlagRequired("dir")
	return cmd
}

func writeFilesToDirectory(cmd *cobra.Command, dest string, files map[string]hubFileEntry, yes bool) (written []string, linked []string, err error) {
	for path, entry := range files {
		if entry.Type != "file" {
			continue
		}
		if err := validateHubFilePath(path); err != nil {
			return nil, nil, err
		}
	}

	if err := confirmDirWipe(cmd, dest, yes); err != nil {
		return nil, nil, err
	}

	absDest, err := filepath.Abs(dest)
	if err != nil {
		return nil, nil, fmt.Errorf("resolving %s: %w", dest, err)
	}
	if err := os.RemoveAll(dest); err != nil {
		return nil, nil, fmt.Errorf("cleaning %s: %w", dest, err)
	}
	if err := os.MkdirAll(dest, 0o755); err != nil {
		return nil, nil, fmt.Errorf("creating %s: %w", dest, err)
	}
	for path, entry := range files {
		if entry.Type != "file" {
			linked = append(linked, fmt.Sprintf("%s (%s)", path, entry.Type))
			continue
		}
		out := filepath.Join(absDest, filepath.FromSlash(path))
		if rel, relErr := filepath.Rel(absDest, out); relErr != nil || rel == ".." || strings.HasPrefix(rel, ".."+string(filepath.Separator)) {
			return nil, nil, fmt.Errorf("path %q would escape %s", path, dest)
		}
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return nil, nil, fmt.Errorf("creating directory for %s: %w", path, err)
		}
		if err := os.WriteFile(out, []byte(entry.Content), 0o644); err != nil {
			return nil, nil, fmt.Errorf("writing %s: %w", path, err)
		}
		written = append(written, path)
	}
	return written, linked, nil
}

func validateHubFilePath(p string) error {
	if p == "" {
		return fmt.Errorf("hub returned empty path")
	}
	if filepath.IsAbs(p) || strings.HasPrefix(p, "/") {
		return fmt.Errorf("hub returned absolute path %q; refusing", p)
	}
	clean := filepath.Clean(filepath.FromSlash(p))
	if clean == ".." || strings.HasPrefix(clean, ".."+string(filepath.Separator)) {
		return fmt.Errorf("path %q would escape destination", p)
	}
	for _, seg := range strings.Split(filepath.ToSlash(clean), "/") {
		if seg == ".." {
			return fmt.Errorf("path %q contains traversal segment", p)
		}
	}
	return nil
}

// confirmDirWipe gates the destructive os.RemoveAll on dest. A missing or
// empty directory proceeds; any other confirms unless yes is set, including a
// directory pulled before, since it may hold local edits. A pulled directory
// (one with SKILL.md or AGENTS.md) used to be replaced with no prompt, so it
// counts as unprompted.
func confirmDirWipe(cmd *cobra.Command, dest string, yes bool) error {
	info, err := os.Stat(dest)
	if os.IsNotExist(err) {
		return nil
	}
	if err != nil {
		return fmt.Errorf("stat %s: %w", dest, err)
	}
	if !info.IsDir() {
		return fmt.Errorf("%s exists and is not a directory", dest)
	}
	entries, err := os.ReadDir(dest)
	if err != nil {
		return fmt.Errorf("reading %s: %w", dest, err)
	}
	if len(entries) == 0 || yes {
		return nil
	}
	return confirmDelete(cmd, deleteConfirmation{
		acceptsUnterminated: true,
		verb:                "replaces",
		target:              "the contents of the destination directory",
		identity:            fmt.Sprintf("Directory: %s (%d entries)", dest, len(entries)),
		unprompted:          hasHubMarker(dest),
	})
}

func hasHubMarker(dir string) bool {
	for _, marker := range []string{"SKILL.md", "AGENTS.md"} {
		if _, err := os.Stat(filepath.Join(dir, marker)); err == nil {
			return true
		}
	}
	return false
}
