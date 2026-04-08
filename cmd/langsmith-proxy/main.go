// Command langsmith-proxy starts a reverse proxy that adds LangSmith tracing
// to Anthropic and OpenAI API calls.
//
// Usage:
//
//	langsmith-proxy [options]              # standalone proxy server
//	langsmith-proxy [options] CMD [ARGS]  # wrap a command with tracing
//
// In wrap mode, the proxy starts on a random port and execs CMD with
// ANTHROPIC_BASE_URL and OPENAI_BASE_URL pointed at the proxy. The proxy
// shuts down when CMD exits.
//
// Environment (PROXY_ prefix takes precedence, falls back to unprefixed):
//
//	PROXY_LANGSMITH_API_KEY   — LangSmith API key (required for real tracing)
//	PROXY_LANGSMITH_PROJECT   — project name, supports Go templates (default: "default")
//	PROXY_LANGSMITH_ENDPOINT  — LangSmith endpoint (default: api.smith.langchain.com)
//
// Then point any SDK at the proxy:
//
//	export ANTHROPIC_BASE_URL=http://localhost:14355
//	export OPENAI_BASE_URL=http://localhost:14355/v1
package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"os/user"
	"path/filepath"
	"text/template"

	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"

	"github.com/langchain-ai/langsmith-go/instrumentation/proxy"
)

func main() {
	addr := flag.String("addr", "", "listen address (default \"localhost:14355\", or \"localhost:0\" in wrap mode)")
	anthropicUpstream := flag.String("anthropic-upstream", "", "Anthropic upstream URL (default: https://api.anthropic.com)")
	openaiUpstream := flag.String("openai-upstream", "", "OpenAI upstream URL (default: https://api.openai.com)")
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Usage:
  langsmith-proxy [options]              Standalone proxy server
  langsmith-proxy [options] CMD [ARGS]   Wrap a command with tracing

Options:
  --addr ADDRESS              Listen address (default "localhost:14355", or ":0" in wrap mode)
  --anthropic-upstream URL    Anthropic upstream URL (default: https://api.anthropic.com)
  --openai-upstream URL       OpenAI upstream URL (default: https://api.openai.com)

Environment (PROXY_ prefix takes precedence, falls back to unprefixed):
  PROXY_LANGSMITH_API_KEY     LangSmith API key (required for real tracing)
  PROXY_LANGSMITH_PROJECT     LangSmith project name, supports Go templates
                              Default: "proxy-{{.User}}-{{.Cmd}}" in wrap mode, "default" otherwise
                              Template fields: {{.User}}, {{.Dir}}, {{.Cmd}}
  PROXY_LANGSMITH_ENDPOINT    LangSmith endpoint (default: api.smith.langchain.com)

Examples:
  langsmith-proxy claude
  langsmith-proxy python my_agent.py
  langsmith-proxy codex exec "hello"
  langsmith-proxy                          # standalone on :14355
`)
	}
	flag.Parse()

	wrapArgs := flag.Args()
	wrapMode := len(wrapArgs) > 0

	// Default addr: fixed port for standalone, random port for wrap mode.
	if *addr == "" {
		if wrapMode {
			*addr = "localhost:0"
		} else {
			*addr = "localhost:14355"
		}
	}

	apiKey := proxyEnv("LANGSMITH_API_KEY")
	cmdName := ""
	if wrapMode {
		cmdName = filepath.Base(wrapArgs[0])
	}
	projectTmpl := proxyEnv("LANGSMITH_PROJECT")
	if projectTmpl == "" && wrapMode {
		projectTmpl = "proxy-{{.User}}-{{.Cmd}}"
	}
	project := renderProject(projectTmpl, cmdName)
	endpoint := proxyEnv("LANGSMITH_ENDPOINT")

	cfg := proxy.Config{
		AnthropicUpstream: *anthropicUpstream,
		OpenAIUpstream:    *openaiUpstream,
		LangSmithAPIKey:   apiKey,
		LangSmithProject:  project,
		LangSmithEndpoint: endpoint,
	}

	if apiKey == "" {
		if !wrapMode {
			log.Println("WARNING: LANGSMITH_API_KEY not set — traces will be printed to stdout instead of sent to LangSmith")
		}
		exporter, err := stdouttrace.New(stdouttrace.WithPrettyPrint())
		if err != nil {
			log.Fatalf("Failed to create stdout exporter: %v", err)
		}
		tp := sdktrace.NewTracerProvider(sdktrace.WithSyncer(exporter))
		cfg.TracerProvider = tp
	}

	p, err := proxy.New(cfg)
	if err != nil {
		log.Fatalf("Failed to create proxy: %v", err)
	}

	// Start listener (may use :0 for random port).
	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("Failed to listen on %s: %v", *addr, err)
	}
	actualAddr := ln.Addr().String()

	aUp := *anthropicUpstream
	if aUp == "" {
		aUp = "https://api.anthropic.com"
	}
	oUp := *openaiUpstream
	if oUp == "" {
		oUp = "https://api.openai.com"
	}

	if !wrapMode {
		log.Printf("LangSmith tracing proxy listening on %s", actualAddr)
		if project != "" {
			log.Printf("  Project:   %s", project)
		}
		log.Printf("  Anthropic: /v1/messages          -> %s", aUp)
		log.Printf("  OpenAI:    /v1/chat/completions  -> %s", oUp)
		log.Printf("  OpenAI:    /v1/responses         -> %s", oUp)
		log.Printf("  Passthru:  /v1/models            -> %s", oUp)
	}

	server := &http.Server{Handler: p}
	go server.Serve(ln)

	if wrapMode {
		os.Exit(runWrapped(wrapArgs, actualAddr, p, server))
	}

	// Standalone mode: wait for signal.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()
	log.Println("Shutting down...")
	p.Shutdown(context.Background())
	server.Shutdown(context.Background())
}

// runWrapped execs the command with ANTHROPIC_BASE_URL and OPENAI_BASE_URL
// set to the proxy, waits for it to finish, then shuts down the proxy.
func runWrapped(args []string, proxyAddr string, p *proxy.Proxy, server *http.Server) int {
	baseURL := "http://" + proxyAddr

	cmd := exec.Command(args[0], args[1:]...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Env = append(os.Environ(),
		"ANTHROPIC_BASE_URL="+baseURL,
		"OPENAI_BASE_URL="+baseURL+"/v1",
	)

	err := cmd.Run()

	// Shut down proxy and flush traces.
	p.Shutdown(context.Background())
	server.Shutdown(context.Background())

	if err != nil {
		if exitErr, ok := err.(*exec.ExitError); ok {
			return exitErr.ExitCode()
		}
		log.Printf("Command failed: %v", err)
		return 1
	}
	return 0
}

// proxyEnv reads PROXY_<name> first, falls back to <name>.
func proxyEnv(name string) string {
	if v := os.Getenv("PROXY_" + name); v != "" {
		return v
	}
	return os.Getenv(name)
}

// projectContext is the template context for the project name.
type projectContext struct {
	User string // current OS username
	Dir  string // basename of the working directory
	Cmd  string // wrapped command name (empty in standalone mode)
}

// renderProject processes the project name as a Go text/template.
func renderProject(raw, cmdName string) string {
	if raw == "" {
		return ""
	}

	t, err := template.New("project").Parse(raw)
	if err != nil {
		log.Printf("WARNING: invalid project template %q: %v", raw, err)
		return raw
	}

	ctx := projectContext{Cmd: cmdName}
	if u, err := user.Current(); err == nil {
		ctx.User = u.Username
	}
	if wd, err := os.Getwd(); err == nil {
		ctx.Dir = filepath.Base(wd)
	}

	var buf bytes.Buffer
	if err := t.Execute(&buf, ctx); err != nil {
		log.Printf("WARNING: project template execution failed: %v", err)
		return raw
	}
	return buf.String()
}
