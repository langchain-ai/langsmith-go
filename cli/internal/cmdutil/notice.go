package cmdutil

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	lsconfig "github.com/langchain-ai/langsmith-go/cli/internal/config"
)

// Notice kinds. Agents reading stderr under --format json match on these.
const (
	NoticeActiveProfile = "active_profile"
	NoticeTruncated     = "truncated"
	NoticeDefaultWindow = "default_time_window"
)

var (
	noticeOut   io.Writer = os.Stderr
	noticesSeen           = map[string]bool{}
)

// ResetNotices sends the next notices to w and lets each kind print again.
// Call it once per command run.
func ResetNotices(w io.Writer) {
	noticeOut = w
	noticesSeen = map[string]bool{}
}

// Notice writes one stderr line telling the reader something about the result
// that stdout doesn't show, such as which profile a command used or that a list
// was cut off. Each kind prints at most once per command run, so a command that
// fetches several pages warns once. With asJSON the line is a JSON object
// {"notice": kind, "message": ..., fields...}; otherwise it is "prefix: message".
// Stdout is never written.
func Notice(kind, prefix, message string, asJSON bool, fields map[string]any) {
	if noticesSeen[kind] {
		return
	}
	noticesSeen[kind] = true
	if !asJSON {
		fmt.Fprintf(noticeOut, "%s: %s\n", prefix, message)
		return
	}
	obj := make(map[string]any, len(fields)+2)
	for k, v := range fields {
		obj[k] = v
	}
	obj["notice"] = kind
	obj["message"] = message
	b, _ := json.Marshal(obj)
	fmt.Fprintln(noticeOut, string(b))
}

// NoteActiveProfile names the profile and workspace a command is about to use
// when more than one profile is saved, so a command run against the wrong one
// is noticed before its output is trusted. flagProfile is the --profile value.
func NoteActiveProfile(flagProfile string, opts client.Options, asJSON bool) {
	cfg, err := lsconfig.Load()
	if err != nil || len(cfg.Profiles) < 2 {
		return
	}
	name := cfg.ResolveProfileName(flagProfile, strings.TrimSpace(os.Getenv("LANGSMITH_PROFILE")))
	if _, ok := cfg.Profiles[name]; !ok {
		return
	}
	workspace := opts.WorkspaceID
	shown := workspace
	if shown == "" {
		shown = "the API key's default"
	}
	msg := fmt.Sprintf("using profile %q, workspace %s (%d profiles saved; pass --profile or --workspace to change)", name, shown, len(cfg.Profiles))
	if opts.ProfileName == "" {
		// The key came from --api-key or LANGSMITH_API_KEY; only the profile's
		// endpoint and workspace apply.
		msg = fmt.Sprintf("using the key from --api-key or LANGSMITH_API_KEY with profile %q's settings, workspace %s (%d profiles saved; pass --profile or --workspace to change)", name, shown, len(cfg.Profiles))
	}
	Notice(NoticeActiveProfile, "note", msg, asJSON, map[string]any{
		"profile":      name,
		"workspace_id": workspace,
	})
}
