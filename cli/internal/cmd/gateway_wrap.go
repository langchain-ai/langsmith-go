package cmd

import (
	"context"
	"crypto/rand"
	"crypto/subtle"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"net/http/httputil"
	"net/url"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"
	"time"

	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/spf13/cobra"
)

// wrapEnvOverrides are removed from the wrapped command's environment so it
// cannot pick another credential, endpoint or provider over the proxy.
var wrapEnvOverrides = []string{
	"ANTHROPIC_API_KEY",
	"ANTHROPIC_AUTH_TOKEN",
	"ANTHROPIC_BASE_URL",
	"CLAUDE_CODE_OAUTH_TOKEN",
	"CLAUDE_CODE_USE_BEDROCK",
	"CLAUDE_CODE_USE_VERTEX",
	"CLAUDE_CODE_USE_FOUNDRY",
	"OPENAI_API_KEY",
	"OPENAI_BASE_URL",
	"OPENAI_API_BASE",
	"OPENAI_ORG_ID",
	"OPENAI_PROJECT_ID",
}

// wrapStrippedHeaders never reach the gateway: the proxy replaces the
// caller's credentials with the LangSmith ones.
var wrapStrippedHeaders = []string{
	"Authorization",
	"X-Api-Key",
	"X-Tenant-Id",
	"X-User-Id",
	"X-Service-Key",
	"Cookie",
	"Proxy-Authorization",
	"OpenAI-Organization",
	"OpenAI-Project",
}

var exitProcess = os.Exit

func newGatewayWrapCmd() *cobra.Command {
	var gatewayURL string
	cmd := &cobra.Command{
		Use:   "wrap [--] COMMAND [ARGS...]",
		Short: "Run a command with its Anthropic and OpenAI requests sent through the LLM Gateway",
		Long: `Run a command with its Anthropic and OpenAI API requests sent through the
LangSmith LLM Gateway, authenticated with your LangSmith credentials.

The command starts a proxy on 127.0.0.1 that attaches those credentials to each
request (refreshing an OAuth access token as it expires), and runs COMMAND
with the Anthropic and OpenAI SDK variables pointing at it:

  ANTHROPIC_BASE_URL                    <proxy>/anthropic
  OPENAI_BASE_URL, OPENAI_API_BASE      <proxy>/openai/v1
  ANTHROPIC_AUTH_TOKEN, OPENAI_API_KEY  an sk- key for this run only, which
                                        works only against this proxy

ANTHROPIC_API_KEY and other provider settings are removed so the command
cannot bypass the gateway. It exits with COMMAND's exit code.

Flags for langsmith go before COMMAND; everything from COMMAND on is passed
to it, so -- is only needed when COMMAND itself starts with a dash.

  langsmith gateway wrap claude
  langsmith gateway wrap claude -p "summarize this repo"
  langsmith --profile eu gateway wrap codex
  langsmith gateway wrap python agent.py`,
		Args: cobra.MinimumNArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			return runGatewayWrap(cmd, args, gatewayURL)
		},
	}
	// Flags after COMMAND belong to it.
	cmd.Flags().SetInterspersed(false)
	addGatewayURLFlag(cmd, &gatewayURL)
	return cmd
}

// runGatewayWrap runs argv behind the gateway proxy and exits with its exit
// code.
func runGatewayWrap(cmd *cobra.Command, argv []string, gatewayFlag string) error {
	ctx := cmd.Context()
	if ctx == nil {
		ctx = context.Background()
	}

	opts, err := resolveAuthenticatedOptions(currentGlobalFlags())
	if err != nil {
		return err
	}
	gatewayURL, err := resolveGatewayURL(gatewayFlag, opts.APIURL)
	if err != nil {
		return err
	}
	auth := sdkGatewayAuth(client.NewWithOptions(opts))
	// Fail before starting the command when the credentials don't resolve,
	// such as an OAuth login that can no longer be refreshed.
	headers, err := auth(ctx)
	if err != nil {
		return err
	}

	binPath, err := exec.LookPath(argv[0])
	if err != nil {
		return fmt.Errorf("%s not found: %w", argv[0], err)
	}

	key, err := newWrapKey()
	if err != nil {
		return err
	}
	handler, err := newGatewayProxy(gatewayURL, key, auth)
	if err != nil {
		return err
	}
	ln, err := net.Listen("tcp", "127.0.0.1:0")
	if err != nil {
		return fmt.Errorf("starting gateway proxy: %w", err)
	}
	srv := &http.Server{Handler: handler, ReadHeaderTimeout: 30 * time.Second}
	go func() { _ = srv.Serve(ln) }()
	defer func() {
		shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		_ = srv.Shutdown(shutdownCtx)
	}()

	fmt.Fprintf(cmd.ErrOrStderr(), "%s → LLM Gateway %s (%s)\n", filepath.Base(argv[0]), gatewayURL, describeGatewayAuth(headers, opts.ProfileName))

	child := exec.Command(binPath, argv[1:]...)
	child.Stdin = os.Stdin
	child.Stdout = cmd.OutOrStdout()
	child.Stderr = cmd.ErrOrStderr()
	child.Env = wrapEnv(os.Environ(), "http://"+ln.Addr().String(), key)
	if err := child.Start(); err != nil {
		return fmt.Errorf("starting %s: %w", argv[0], err)
	}

	// Signals sent to this process are forwarded so the command stops with it,
	// except a terminal Ctrl-C, which reaches the command through the
	// foreground process group already and would otherwise arrive twice.
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, os.Interrupt, syscall.SIGTERM)
	defer signal.Stop(sigs)
	done := make(chan error, 1)
	go func() { done <- child.Wait() }()
	for {
		select {
		case sig := <-sigs:
			if sig == os.Interrupt && terminalDeliversInterrupt() {
				continue
			}
			_ = child.Process.Signal(sig)
		case err := <-done:
			var exitErr *exec.ExitError
			if errors.As(err, &exitErr) {
				code := exitErr.ExitCode()
				if code < 0 {
					code = 1
				}
				// Release the listener before exiting, since os.Exit skips defers.
				shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
				_ = srv.Shutdown(shutdownCtx)
				cancel()
				exitProcess(code)
				return nil
			}
			return err
		}
	}
}

// newWrapKey returns a random per-run key. It starts with sk-, like Anthropic
// and OpenAI keys, for clients that check the prefix before sending it.
func newWrapKey() (string, error) {
	b := make([]byte, 48)
	if _, err := rand.Read(b); err != nil {
		return "", fmt.Errorf("generating proxy key: %w", err)
	}
	return "sk-" + base64.RawURLEncoding.EncodeToString(b), nil
}

// carriesKey reports whether the request carries key, as a bearer token or
// X-Api-Key.
func carriesKey(r *http.Request, key string) bool {
	got, ok := strings.CutPrefix(r.Header.Get("Authorization"), "Bearer ")
	if !ok {
		got = r.Header.Get("X-Api-Key")
	}
	return subtle.ConstantTimeCompare([]byte(got), []byte(key)) == 1
}

type gatewayAuthContextKey struct{}

// newGatewayProxy forwards requests to the same path on the gateway, so
// /anthropic/... and /openai/... reach those provider routes. A request must
// carry the per-run key the wrapped command was given, as a bearer token or
// X-Api-Key; the proxy then replaces every caller credential with the
// LangSmith auth headers, so no other local process can borrow them through
// the listener.
func newGatewayProxy(gatewayURL, key string, auth gatewayAuth) (http.Handler, error) {
	target, err := url.Parse(strings.TrimRight(gatewayURL, "/"))
	if err != nil {
		return nil, err
	}
	proxy := &httputil.ReverseProxy{
		Rewrite: func(pr *httputil.ProxyRequest) {
			pr.SetURL(target)
			for _, h := range wrapStrippedHeaders {
				pr.Out.Header.Del(h)
			}
			resolved, _ := pr.In.Context().Value(gatewayAuthContextKey{}).(http.Header)
			for name, values := range resolved {
				pr.Out.Header[name] = values
			}
		},
		// Stream server-sent events to the client as they arrive.
		FlushInterval: -1,
		ErrorHandler: func(w http.ResponseWriter, r *http.Request, err error) {
			writeProxyError(w, http.StatusBadGateway, "api_error", "LangSmith gateway proxy: "+err.Error())
		},
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !carriesKey(r, key) {
			writeProxyError(w, http.StatusUnauthorized, "authentication_error",
				"request does not carry the credentials 'langsmith gateway wrap' gave this command (a setting such as Claude Code's apiKeyHelper may override them)")
			return
		}
		headers, err := auth(r.Context())
		if err != nil {
			writeProxyError(w, http.StatusUnauthorized, "authentication_error", "LangSmith credentials: "+err.Error())
			return
		}
		ctx := context.WithValue(r.Context(), gatewayAuthContextKey{}, headers)
		proxy.ServeHTTP(w, r.WithContext(ctx))
	}), nil
}

// writeProxyError writes an error in the Anthropic API's shape. OpenAI clients
// read the same error.message field, so both show it to the user.
func writeProxyError(w http.ResponseWriter, status int, errType, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(map[string]any{
		"type":  "error",
		"error": map[string]string{"type": errType, "message": message},
	})
}

// wrapEnv points the Anthropic and OpenAI SDK variables at the proxy and drops
// variables that would make the command use another credential or provider.
func wrapEnv(environ []string, proxyURL, key string) []string {
	out := make([]string, 0, len(environ)+5)
	for _, kv := range environ {
		name, _, _ := strings.Cut(kv, "=")
		drop := false
		for _, o := range wrapEnvOverrides {
			if strings.EqualFold(name, o) {
				drop = true
				break
			}
		}
		if !drop {
			out = append(out, kv)
		}
	}
	return append(out,
		"ANTHROPIC_BASE_URL="+proxyURL+"/anthropic",
		"ANTHROPIC_AUTH_TOKEN="+key,
		"OPENAI_BASE_URL="+proxyURL+"/openai/v1",
		"OPENAI_API_BASE="+proxyURL+"/openai/v1",
		"OPENAI_API_KEY="+key,
	)
}
