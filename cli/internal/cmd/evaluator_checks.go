package cmd

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

// secretPatterns match credential formats with distinctive prefixes, so
// ordinary strings in evaluator code don't match.
var secretPatterns = []struct {
	kind string
	re   *regexp.Regexp
}{
	{"LangSmith API key", regexp.MustCompile(`\blsv2_(?:pt|sk)_[A-Za-z0-9_]{20,}`)},
	{"Anthropic API key", regexp.MustCompile(`\bsk-ant-[A-Za-z0-9_-]{20,}`)},
	{"OpenAI API key", regexp.MustCompile(`\bsk-(?:proj-|svcacct-|admin-)?[A-Za-z0-9_-]{20,}`)},
	{"AWS access key ID", regexp.MustCompile(`\b(?:AKIA|ASIA)[0-9A-Z]{16}\b`)},
	{"Google API key", regexp.MustCompile(`\bAIza[0-9A-Za-z_-]{35}`)},
	{"GitHub token", regexp.MustCompile(`\b(?:gh[pousr]_[A-Za-z0-9]{30,}|github_pat_[A-Za-z0-9_]{40,})`)},
	{"Slack token", regexp.MustCompile(`\bxox[abprs]-[A-Za-z0-9-]{10,}`)},
	{"Stripe secret key", regexp.MustCompile(`\b[sr]k_live_[A-Za-z0-9]{20,}`)},
	{"private key", regexp.MustCompile(`-----BEGIN (?:[A-Z]+ )?PRIVATE KEY-----`)},
	// The user and password exclude {, }, $ and a % that doesn't start a
	// %XX escape, so URLs built from variables (f-strings, template literals,
	// % formatting) don't match but percent-encoded passwords do.
	{"password in a URL", regexp.MustCompile(`[a-z][a-z0-9+.-]*://(?:[^\s/:@'"{}$%]|%[0-9A-Fa-f]{2})+:(?:[^\s/@'"{}$%]|%[0-9A-Fa-f]{2})+@`)},
}

// checkForSecrets refuses an upload whose code or dependencies file contains
// something that looks like a credential. Uploaded code is stored on the
// evaluator, where anyone in the workspace can read it. The error names the
// line and the kind of credential, never the value. With ignore set, the
// findings go to stderr as a warning and the upload goes ahead.
func checkForSecrets(file, content string, ignore bool) error {
	var found []string
	for i, line := range strings.Split(content, "\n") {
		for _, p := range secretPatterns {
			if p.re.MatchString(line) {
				found = append(found, fmt.Sprintf("  %s:%d: looks like a %s", file, i+1, p.kind))
				break
			}
		}
	}
	if len(found) == 0 {
		return nil
	}
	if ignore {
		fmt.Fprintf(os.Stderr, "warning: uploading %s although it contains what looks like a hardcoded credential (--ignore-credential-check):\n%s\n", file, strings.Join(found, "\n"))
		return nil
	}
	return fmt.Errorf("not uploaded: %s contains what looks like a hardcoded credential, and uploaded code can be read by anyone in the workspace:\n%s\nRemove it. An advanced evaluator (--advanced) can use a workspace secret instead: create it in the workspace settings and pass --workspace-secret NAME. If these are not credentials, pass --ignore-credential-check", file, strings.Join(found, "\n"))
}

var (
	// import a, b.c as d / from x.y import z / from .x import z
	pyImportRe     = regexp.MustCompile(`^\s*import\s+(.+)$`)
	pyFromImportRe = regexp.MustCompile(`^\s*from\s+(\.*[\w.]*)\s+import\b`)
)

// errLocalImports is returned by checkLocalImports so callers and tests can
// tell it apart from other upload errors.
var errLocalImports = errors.New("local imports")

// checkLocalImports refuses code that imports modules from the evaluator's
// own folder. Only the file is uploaded, so those imports fail when the
// evaluator runs. Python relative imports, and imports of a module that is a
// .py file or a folder next to the evaluator, count as local; JavaScript
// imports of a path starting with ./, ../ or / do.
func checkLocalImports(file, code, language string) error {
	var modules []string
	seen := map[string]bool{}
	add := func(m string) {
		if !seen[m] {
			seen[m] = true
			modules = append(modules, m)
		}
	}
	switch language {
	case "python":
		for _, line := range pythonCodeLines(code) {
			for _, module := range pythonImportedModules(line) {
				if pythonModuleIsLocal(file, module) {
					add(module)
				}
			}
		}
	case "javascript":
		for _, spec := range jsImportSpecifiers(code) {
			if strings.HasPrefix(spec, "./") || strings.HasPrefix(spec, "../") || strings.HasPrefix(spec, "/") {
				add(spec)
			}
		}
	}
	if len(modules) == 0 {
		return nil
	}
	quoted := make([]string, len(modules))
	for i, m := range modules {
		quoted[i] = fmt.Sprintf("%q", m)
	}
	return fmt.Errorf("%w: not uploaded: %s imports %s from its own folder. Only %s is uploaded, so the evaluator would fail when it runs; move that code into the file, or publish it as a package and list it in --dependencies", errLocalImports, file, strings.Join(quoted, ", "), filepath.Base(file))
}

// pythonImportedModules returns the modules a line of Python imports.
func pythonImportedModules(line string) []string {
	if m := pyFromImportRe.FindStringSubmatch(line); m != nil {
		return []string{m[1]}
	}
	m := pyImportRe.FindStringSubmatch(line)
	if m == nil {
		return nil
	}
	var modules []string
	for _, part := range strings.Split(strings.SplitN(m[1], "#", 2)[0], ",") {
		if fields := strings.Fields(part); len(fields) > 0 {
			modules = append(modules, fields[0])
		}
	}
	return modules
}

// pythonCodeLines returns the lines of code with the contents of strings and
// comments blanked out, so import examples in prompts or docstrings are not
// read as imports. Line breaks are kept.
func pythonCodeLines(code string) []string {
	out := []byte(code)
	blank := func(from, to int) {
		for k := from; k < to && k < len(out); k++ {
			if out[k] != '\n' {
				out[k] = ' '
			}
		}
	}
	for i := 0; i < len(code); {
		switch c := code[i]; {
		case c == '#':
			end := strings.IndexByte(code[i:], '\n')
			if end < 0 {
				end = len(code) - i
			}
			blank(i, i+end)
			i += end
		case c == '"' || c == '\'':
			quote := code[i : i+1]
			if strings.HasPrefix(code[i:], strings.Repeat(quote, 3)) {
				quote = strings.Repeat(quote, 3)
			}
			j := i + len(quote)
			for j < len(code) && !strings.HasPrefix(code[j:], quote) {
				if code[j] == '\\' {
					j++
				} else if len(quote) == 1 && code[j] == '\n' {
					break
				}
				j++
			}
			end := min(j+len(quote), len(code))
			blank(i, end)
			i = end
		default:
			i++
		}
	}
	return strings.Split(string(out), "\n")
}

// jsImportSpecifiers returns the module specifiers a JavaScript file imports:
// import ... from "x", import "x", export ... from "x", import("x") and
// require("x"). Specifiers are read only from string literals that follow
// those keywords in code, so text inside comments, other strings and regex
// literals is ignored. Code inside template ${} expressions is scanned too.
func jsImportSpecifiers(code string) []string {
	s := &jsScanner{code: code}
	s.scan(false)
	var specifiers []string
	for _, m := range jsImportRe.FindAllStringSubmatch(s.masked.String(), -1) {
		var n int
		_, _ = fmt.Sscan(m[1], &n)
		specifiers = append(specifiers, s.literals[n])
	}
	return specifiers
}

// jsImportRe matches an import keyword that is not a property, such as
// Buffer.from, followed by a masked string literal. A spread (...) before the
// keyword still counts.
var jsImportRe = regexp.MustCompile(`(?:^|[^.\w$]|\.\.\.)(?:from|import|require)\s*\(?\s*"#(\d+)"`)

// jsScanner copies code to masked with comments dropped and each string,
// template or regex literal replaced by a "#N" placeholder (literals[N]).
type jsScanner struct {
	code     string
	i        int
	masked   strings.Builder
	literals []string
	// last is the last code character written, used to tell a regex literal
	// from division.
	last byte
}

func (s *jsScanner) emit(b byte) {
	s.masked.WriteByte(b)
	if b != ' ' && b != '\n' && b != '\t' && b != '\r' {
		s.last = b
	}
}

func (s *jsScanner) literal(text string) {
	s.literals = append(s.literals, text)
	fmt.Fprintf(&s.masked, ` "#%d" `, len(s.literals)-1)
	s.last = '"'
}

// regexAllowed reports whether a / here starts a regex literal rather than
// division: at the start, after an operator or opening bracket, or after a
// keyword such as return.
func (s *jsScanner) regexAllowed() bool {
	if s.last == 0 || strings.IndexByte("(,=:[!&|?{};+-*%<>~^", s.last) >= 0 {
		return true
	}
	out := strings.TrimRight(s.masked.String(), " \t\r\n")
	for _, kw := range []string{"return", "typeof", "case", "do", "else", "in", "of", "void", "yield", "await"} {
		if strings.HasSuffix(out, kw) {
			before := strings.TrimSuffix(out, kw)
			if before == "" || !isJSIdentByte(before[len(before)-1]) {
				return true
			}
		}
	}
	return false
}

func isJSIdentByte(b byte) bool {
	return b == '_' || b == '$' || b == '.' || b >= '0' && b <= '9' || b >= 'a' && b <= 'z' || b >= 'A' && b <= 'Z'
}

// scan copies code until the end, or, with inExpr, until the } that closes a
// template ${} expression.
func (s *jsScanner) scan(inExpr bool) {
	depth := 0
	for s.i < len(s.code) {
		c := s.code[s.i]
		rest := s.code[s.i:]
		switch {
		case strings.HasPrefix(rest, "//"):
			end := strings.IndexByte(rest, '\n')
			if end < 0 {
				end = len(rest)
			}
			s.i += end
		case strings.HasPrefix(rest, "/*"):
			end := strings.Index(rest[2:], "*/")
			if end < 0 {
				s.i = len(s.code)
			} else {
				s.i += end + 4
			}
			s.emit(' ')
		case c == '"' || c == '\'':
			j := s.i + 1
			for j < len(s.code) && s.code[j] != c && s.code[j] != '\n' {
				if s.code[j] == '\\' {
					j++
				}
				j++
			}
			s.literal(s.code[s.i+1 : min(j, len(s.code))])
			s.i = j + 1
		case c == '`':
			s.template()
		case c == '/' && s.regexAllowed():
			s.regex()
		case inExpr && c == '{':
			depth++
			s.emit(c)
			s.i++
		case inExpr && c == '}':
			if depth == 0 {
				s.i++
				return
			}
			depth--
			s.emit(c)
			s.i++
		default:
			s.emit(c)
			s.i++
		}
	}
}

// template records a template literal's text as one literal and scans each
// ${} expression as code after it.
func (s *jsScanner) template() {
	start := s.i + 1
	s.i = start
	var text strings.Builder
	var exprs []int
	for s.i < len(s.code) && s.code[s.i] != '`' {
		switch {
		case s.code[s.i] == '\\':
			text.WriteString(s.code[s.i:min(s.i+2, len(s.code))])
			s.i += 2
		case strings.HasPrefix(s.code[s.i:], "${"):
			text.WriteString("${}")
			s.i += 2
			exprs = append(exprs, s.i)
			// Skip the expression here; it is scanned as code below.
			sub := &jsScanner{code: s.code, i: s.i}
			sub.scan(true)
			s.i = sub.i
		default:
			text.WriteByte(s.code[s.i])
			s.i++
		}
	}
	s.i++
	s.literal(text.String())
	for _, at := range exprs {
		sub := &jsScanner{code: s.code, i: at, literals: s.literals}
		sub.scan(true)
		s.literals = sub.literals
		// The ; keeps separate expressions from reading as one statement.
		s.masked.WriteString(" ; " + sub.masked.String() + " ; ")
	}
}

// regex skips a regex literal, including character classes and flags.
func (s *jsScanner) regex() {
	j := s.i + 1
	inClass := false
	for j < len(s.code) && s.code[j] != '\n' {
		c := s.code[j]
		if c == '\\' {
			j += 2
			continue
		}
		if c == '[' {
			inClass = true
		} else if c == ']' {
			inClass = false
		} else if c == '/' && !inClass {
			break
		}
		j++
	}
	j++
	for j < len(s.code) && isJSIdentByte(s.code[j]) && s.code[j] != '.' {
		j++
	}
	s.i = min(j, len(s.code))
	s.emit(' ')
	s.last = ')'
}

// pythonModuleIsLocal reports whether module comes from file's folder: a
// relative import, or a .py file or a folder of Python files next to file.
// The evaluator file itself doesn't count, since it is uploaded under another
// name.
func pythonModuleIsLocal(file, module string) bool {
	if strings.HasPrefix(module, ".") {
		return true
	}
	top := strings.SplitN(module, ".", 2)[0]
	if top == "" {
		return false
	}
	dir := filepath.Dir(file)
	if sibling := filepath.Join(dir, top+".py"); filepath.Clean(sibling) != filepath.Clean(file) {
		if _, err := os.Stat(sibling); err == nil {
			return true
		}
	}
	pyFiles, _ := filepath.Glob(filepath.Join(dir, top, "*.py"))
	return len(pyFiles) > 0
}

var (
	pyTopFuncRe = regexp.MustCompile(`(?m)^(?:async\s+)?def\s+([A-Za-z]\w*)\s*\(`)
	jsTopFuncRe = regexp.MustCompile(`(?m)^(?:export\s+)?(?:async\s+)?function\s*\*?\s*([A-Za-z$][\w$]*)\s*\(|^(?:export\s+)?(?:const|let|var)\s+([A-Za-z$][\w$]*)\s*=\s*(?:async\s*)?(?:function\b|\([^)]*\)\s*=>|[A-Za-z$][\w$]*\s*=>)`)
)

// withEntryPointHint adds a rename hint to the server's error for code that
// lacks the entry point, naming the top-level functions the file defines.
// Thread and trajectory input (input set) is only passed to evaluate(ctx).
func withEntryPointHint(err error, code, language, input string) error {
	if err == nil || !strings.Contains(err.Error(), "must define a function named") {
		return err
	}
	var names []string
	if language == "python" {
		for _, m := range pyTopFuncRe.FindAllStringSubmatch(code, -1) {
			names = append(names, m[1])
		}
	} else {
		for _, m := range jsTopFuncRe.FindAllStringSubmatch(code, -1) {
			names = append(names, m[1]+m[2])
		}
	}
	entry, name := "perform_eval(run, example)", "perform_eval"
	if language != "python" {
		entry, name = "function performEval(run, example)", "performEval"
	}
	hint := fmt.Sprintf("hint: the file must define %s, or evaluate(ctx) for an advanced evaluator (--advanced)", entry)
	if input != "" {
		name = "evaluate"
		hint = fmt.Sprintf("hint: --input %s needs the file to define evaluate(ctx)", input)
	}
	if len(names) > 0 {
		hint += fmt.Sprintf("; it defines %s, so rename the one that scores to %s", strings.Join(names, ", "), name)
	}
	return fmt.Errorf("%w\n%s", err, hint)
}
