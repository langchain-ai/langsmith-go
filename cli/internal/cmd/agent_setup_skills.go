package cmd

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/fs"
	"maps"
	"net/http"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"runtime"
	"slices"
	"strings"
	"sync"
	"time"
)

// skillMarkerFile is written into every skill directory setup installs. Its
// presence means setup owns the directory and may replace it on a later run;
// a directory without it is the user's and is left alone unless --force.
const skillMarkerFile = ".langsmith-setup.json"

const (
	skillsArchiveMaxBytes = 20 << 20
	skillsFileMaxBytes    = 4 << 20
)

// legacySkillSuccessors maps each skill langsmith-skills shipped to the
// langchain-skills skill that replaced it. Setup removes a legacy skill it
// installed only when the same run leaves the successor installed, since older
// langchain-skills releases lack some successors.
var legacySkillSuccessors = map[string]string{
	"langsmith-custom-apps":             "langchain-agent-experiments",
	"langsmith-dataset":                 "langchain-agent-experiments",
	"langsmith-evaluator":               "langchain-agent-experiments",
	"langsmith-experiment-diagnosis":    "langchain-agent-experiments",
	"langsmith-graph-evaluator":         "langchain-agent-experiments",
	"langsmith-online-eval-engineering": "langchain-agent-monitoring",
	"langsmith-trace":                   "langchain-agent-tracing",
}

// skillsSource is a repository that setup installs skills from, and the
// directory in its source tree that holds them.
type skillsSource struct {
	repo string
	dir  string
}

var (
	currentSkillsSource = skillsSource{repo: skillsRepo, dir: "skills/"}
	// --skills-ref pins written before the move name refs of langsmith-skills,
	// which kept its skills under config/skills/.
	legacySkillsSource = skillsSource{repo: "langchain-ai/langsmith-skills", dir: "config/skills/"}
)

var (
	skillNamePattern = regexp.MustCompile(`^[a-z0-9][a-z0-9._-]*$`)
	skillsRefPattern = regexp.MustCompile(`^[A-Za-z0-9][A-Za-z0-9._/-]*$`)
)

// skillFile is one file of a skill, keyed by its slash-separated path inside
// the skill directory.
type skillFile struct {
	data []byte
	mode fs.FileMode
}

// skillsPackage is the set of skills in one ref of a skills repository.
type skillsPackage struct {
	source skillsSource
	ref    string
	skills map[string]map[string]skillFile
}

func (p *skillsPackage) names() []string {
	var names []string
	for n := range p.skills {
		names = append(names, n)
	}
	slices.Sort(names)
	return names
}

// fetchSkillsArchive downloads the source tarball of repo at ref.
// This is a public GitHub download, not a LangSmith API call. Tests replace it.
var fetchSkillsArchive = func(ctx context.Context, repo, ref string) (io.ReadCloser, error) {
	u := "https://github.com/" + repo + "/archive/" + ref + ".tar.gz"
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, u, nil)
	if err != nil {
		cancel()
		return nil, err
	}
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		cancel()
		return nil, fmt.Errorf("downloading %s: %w", u, err)
	}
	if resp.StatusCode != http.StatusOK {
		_ = resp.Body.Close()
		cancel()
		return nil, fmt.Errorf("downloading %s: %s", u, resp.Status)
	}
	return cancelOnClose{ReadCloser: resp.Body, cancel: cancel}, nil
}

type cancelOnClose struct {
	io.ReadCloser
	cancel context.CancelFunc
}

func (c cancelOnClose) Close() error {
	err := c.ReadCloser.Close()
	c.cancel()
	return err
}

func validateSkillsRef(ref string) error {
	if !skillsRefPattern.MatchString(ref) || strings.Contains(ref, "..") {
		return fmt.Errorf("invalid --skills-ref %q", ref)
	}
	return nil
}

// cachedSkillsLoader downloads the skills once per run, on first use, so a run
// that only touches Claude Code (which installs the plugin) never downloads.
func cachedSkillsLoader(ref string) func(context.Context) (*skillsPackage, error) {
	var (
		once sync.Once
		pkg  *skillsPackage
		err  error
	)
	return func(ctx context.Context) (*skillsPackage, error) {
		once.Do(func() {
			pkg, err = loadSkills(ctx, currentSkillsSource, ref)
			if err != nil && ref != defaultSkillsRef {
				// A ref missing from langchain-skills, or one from before its
				// skills/ layout, is a pin made against langsmith-skills.
				if legacy, legacyErr := loadSkills(ctx, legacySkillsSource, ref); legacyErr == nil {
					pkg, err = legacy, nil
				}
			}
		})
		return pkg, err
	}
}

func loadSkills(ctx context.Context, src skillsSource, ref string) (*skillsPackage, error) {
	body, err := fetchSkillsArchive(ctx, src.repo, ref)
	if err != nil {
		return nil, err
	}
	defer body.Close()
	return parseSkillsArchive(body, src, ref)
}

// parseSkillsArchive reads <src.dir><name>/... out of a GitHub source
// tarball (whose entries sit under one top-level directory). Only regular
// files with clean relative paths are kept, and sizes are capped.
func parseSkillsArchive(r io.Reader, src skillsSource, ref string) (*skillsPackage, error) {
	gz, err := gzip.NewReader(r)
	if err != nil {
		return nil, fmt.Errorf("reading skills archive: %w", err)
	}
	defer gz.Close()
	tr := tar.NewReader(io.LimitReader(gz, skillsArchiveMaxBytes+1))
	pkg := &skillsPackage{source: src, ref: ref, skills: map[string]map[string]skillFile{}}
	var total int64
	for {
		hdr, err := tr.Next()
		if errors.Is(err, io.EOF) {
			break
		}
		if err != nil {
			return nil, fmt.Errorf("reading skills archive: %w", err)
		}
		if hdr.Typeflag != tar.TypeReg {
			continue
		}
		parts := strings.SplitN(hdr.Name, "/", 2)
		if len(parts) != 2 {
			continue
		}
		rel, ok := strings.CutPrefix(parts[1], src.dir)
		if !ok {
			continue
		}
		name, file, ok := strings.Cut(rel, "/")
		if !ok || !skillNamePattern.MatchString(name) {
			continue
		}
		clean := path.Clean(file)
		if file == "" || clean != file || path.IsAbs(clean) || clean == ".." || strings.HasPrefix(clean, "../") || path.Base(clean) == skillMarkerFile {
			continue
		}
		if hdr.Size > skillsFileMaxBytes {
			return nil, fmt.Errorf("skills archive: %s is larger than %d bytes", hdr.Name, skillsFileMaxBytes)
		}
		total += hdr.Size
		if total > skillsArchiveMaxBytes {
			return nil, fmt.Errorf("skills archive is larger than %d bytes", skillsArchiveMaxBytes)
		}
		data, err := io.ReadAll(io.LimitReader(tr, hdr.Size))
		if err != nil {
			return nil, fmt.Errorf("reading skills archive: %w", err)
		}
		mode := fs.FileMode(0o644)
		if hdr.Mode&0o111 != 0 {
			mode = 0o755
		}
		if pkg.skills[name] == nil {
			pkg.skills[name] = map[string]skillFile{}
		}
		pkg.skills[name][clean] = skillFile{data: data, mode: mode}
	}
	for name, files := range pkg.skills {
		if _, ok := files["SKILL.md"]; !ok {
			delete(pkg.skills, name)
		}
	}
	if len(pkg.skills) == 0 {
		return nil, fmt.Errorf("no skills found in github.com/%s@%s", src.repo, ref)
	}
	return pkg, nil
}

// planSkillsDir plans copying every skill into root/<name>. Each skill is
// created when missing, left alone when identical, replaced when setup wrote
// it before, and skipped (unless --force) when the user owns a directory of
// that name.
func planSkillsDir(ctx context.Context, env *setupPlanEnv, agent, root string) setupStep {
	step := setupStep{Agent: agent, Component: "skills", Target: root}
	pkg, err := env.skills(ctx)
	if err != nil {
		step.Action, step.Error = stepFailed, err.Error()
		return step
	}
	env.loadedSkillsRepo = pkg.source.repo
	var create, update, conflict, remove []string
	for _, name := range pkg.names() {
		dir := filepath.Join(root, name)
		state, err := skillDirState(dir, pkg.skills[name])
		if err != nil {
			step.Action, step.Error = stepFailed, err.Error()
			return step
		}
		switch state {
		case "missing":
			create = append(create, name)
		case "owned":
			update = append(update, name)
		case "foreign":
			if env.force {
				update = append(update, name)
			} else {
				conflict = append(conflict, name)
			}
		}
	}

	for _, name := range slices.Sorted(maps.Keys(legacySkillSuccessors)) {
		successor := legacySkillSuccessors[name]
		if _, inPkg := pkg.skills[name]; inPkg || !pathExists(filepath.Join(root, name, skillMarkerFile)) {
			continue
		}
		if _, ok := pkg.skills[successor]; ok && !slices.Contains(conflict, successor) {
			remove = append(remove, name)
		}
	}

	var details []string
	if len(create) > 0 {
		details = append(details, "add "+strings.Join(create, ", "))
	}
	if len(update) > 0 {
		details = append(details, "update "+strings.Join(update, ", "))
	}
	if len(remove) > 0 {
		details = append(details, "remove "+strings.Join(remove, ", ")+" (replaced by the langchain-skills skills)")
	}
	if len(conflict) > 0 {
		details = append(details, "skip "+strings.Join(conflict, ", ")+" (directory not written by langsmith setup; pass --force to replace)")
	}
	step.Detail = strings.Join(details, "; ")
	switch {
	case len(update) > 0 || len(remove) > 0:
		step.Action = stepUpdate
	case len(create) > 0:
		step.Action = stepCreate
	case len(conflict) > 0:
		step.Action = stepSkip
	default:
		step.Action = stepUnchanged
		step.Detail = fmt.Sprintf("%d skills up to date", len(pkg.skills))
	}
	writes := append(append([]string{}, create...), update...)
	step.apply = func(context.Context) error {
		for _, name := range writes {
			if err := installSkillDir(root, name, pkg); err != nil {
				return err
			}
		}
		for _, name := range remove {
			if err := os.RemoveAll(filepath.Join(root, name)); err != nil {
				return fmt.Errorf("removing %s: %w", filepath.Join(root, name), err)
			}
		}
		return nil
	}
	return step
}

// skillDirState compares dir with the wanted files: "missing", "same",
// "owned" (setup's marker is present but contents differ), or "foreign".
func skillDirState(dir string, want map[string]skillFile) (string, error) {
	fi, err := os.Lstat(dir)
	if errors.Is(err, os.ErrNotExist) {
		return "missing", nil
	}
	if err != nil {
		return "", fmt.Errorf("stat %s: %w", dir, err)
	}
	if !fi.IsDir() {
		return "foreign", nil
	}
	type gotFile struct {
		data []byte
		exec bool
	}
	got := map[string]gotFile{}
	err = filepath.WalkDir(dir, func(p string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			return nil
		}
		rel, err := filepath.Rel(dir, p)
		if err != nil {
			return err
		}
		rel = filepath.ToSlash(rel)
		if rel == skillMarkerFile {
			return nil
		}
		if !d.Type().IsRegular() {
			got[rel] = gotFile{}
			return nil
		}
		info, err := d.Info()
		if err != nil {
			return err
		}
		data, err := os.ReadFile(p)
		if err != nil {
			return err
		}
		got[rel] = gotFile{data: data, exec: info.Mode().Perm()&0o100 != 0}
		return nil
	})
	if err != nil {
		return "", fmt.Errorf("reading %s: %w", dir, err)
	}
	// Windows has no exec bit, so only bytes are compared there.
	checkExec := runtime.GOOS != "windows"
	same := len(got) == len(want)
	if same {
		for rel, f := range want {
			g, ok := got[rel]
			if !ok || !bytes.Equal(g.data, f.data) || (checkExec && g.exec != (f.mode&0o100 != 0)) {
				same = false
				break
			}
		}
	}
	switch {
	case same:
		return "same", nil
	case pathExists(filepath.Join(dir, skillMarkerFile)):
		return "owned", nil
	default:
		return "foreign", nil
	}
}

// installSkillDir writes the skill into a temporary sibling directory, then
// swaps it into place so a failure never leaves a half-written skill.
func installSkillDir(root, name string, pkg *skillsPackage) error {
	if err := os.MkdirAll(root, 0o755); err != nil {
		return fmt.Errorf("creating %s: %w", root, err)
	}
	tmp, err := os.MkdirTemp(root, "."+name+".tmp-")
	if err != nil {
		return fmt.Errorf("creating temp dir in %s: %w", root, err)
	}
	defer func() { _ = os.RemoveAll(tmp) }()
	if err := os.Chmod(tmp, 0o755); err != nil {
		return err
	}
	for rel, f := range pkg.skills[name] {
		p := filepath.Join(tmp, filepath.FromSlash(rel))
		if err := os.MkdirAll(filepath.Dir(p), 0o755); err != nil {
			return fmt.Errorf("creating %s: %w", filepath.Dir(p), err)
		}
		if err := os.WriteFile(p, f.data, f.mode); err != nil {
			return fmt.Errorf("writing %s: %w", p, err)
		}
	}
	marker, _ := json.MarshalIndent(map[string]string{
		"source": "github.com/" + pkg.source.repo,
		"ref":    pkg.ref,
		"note":   "Installed by `langsmith setup`. Delete this file to keep your own edits from being replaced.",
	}, "", "  ")
	if err := os.WriteFile(filepath.Join(tmp, skillMarkerFile), append(marker, '\n'), 0o644); err != nil {
		return fmt.Errorf("writing marker: %w", err)
	}

	dest := filepath.Join(root, name)
	var old string
	if _, err := os.Lstat(dest); err == nil {
		old = filepath.Join(root, fmt.Sprintf(".%s.old-%d", name, time.Now().UnixNano()))
		if err := os.Rename(dest, old); err != nil {
			return fmt.Errorf("moving aside %s: %w", dest, err)
		}
	}
	if err := os.Rename(tmp, dest); err != nil {
		if old != "" {
			_ = os.Rename(old, dest)
		}
		return fmt.Errorf("installing %s: %w", dest, err)
	}
	if old != "" {
		_ = os.RemoveAll(old)
	}
	return nil
}
