package structured

import (
	"fmt"
	"sort"
	"strings"
	"text/template"
	"time"
)

func templateFuncs() template.FuncMap {
	return template.FuncMap{
		"dash":              dash,
		"formatBytes":       formatBytes,
		"formatBytesOrDash": formatBytesOrDash,
		"formatCount":       formatCount,
		"formatKV":          formatKV,
		"formatTime":        formatTime,
		"joinOrDash":        joinOrDash,
		"shortID":           shortID,
	}
}

func joinOrDash(items []string) string {
	if len(items) == 0 {
		return "-"
	}
	return strings.Join(items, ", ")
}

// formatKV renders a map as sorted "key=value" pairs, or "" when empty.
func formatKV(m map[string]string) string {
	if len(m) == 0 {
		return ""
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	pairs := make([]string, 0, len(keys))
	for _, k := range keys {
		pairs = append(pairs, k+"="+m[k])
	}
	return strings.Join(pairs, ", ")
}

func dash(s string) string {
	if s == "" {
		return "-"
	}
	return s
}

func formatBytes(b int64) string {
	const gb = 1024 * 1024 * 1024
	const mb = 1024 * 1024
	if b >= gb {
		return fmt.Sprintf("%.1f GB", float64(b)/float64(gb))
	}
	return fmt.Sprintf("%.0f MB", float64(b)/float64(mb))
}

func formatBytesOrDash(b int64) string {
	if b <= 0 {
		return "-"
	}
	return formatBytes(b)
}

func formatCount(n int64) string {
	if n <= 0 {
		return "-"
	}
	return fmt.Sprintf("%d", n)
}

func formatTime(ts string) string {
	t, err := time.Parse(time.RFC3339Nano, ts)
	if err != nil {
		return ts
	}
	return t.Local().Format("2006-01-02 15:04")
}

func shortID(id string) string {
	if id == "" {
		return "-"
	}
	if len(id) > 8 {
		return id[:8] + "..."
	}
	return id
}
