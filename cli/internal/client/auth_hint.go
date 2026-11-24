package client

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"

	"github.com/langchain-ai/langsmith-go/option"
)

// authHint prints one hint on the first 401 or 403, naming the command that
// fixes it, so an agent fixes its credentials instead of retrying them on every
// command. A key from another region (an EU key sent to the US endpoint) is the
// most common cause, and the API answers it with 403 as well as 401. An
// org-scoped key sent without a workspace also gets a 403.
type authHint struct {
	once sync.Once
	w    io.Writer
}

const regionAdvice = "EU keys need https://eu.api.smith.langchain.com (set LANGSMITH_ENDPOINT, --api-url, or the selected profile's API URL)"

type expectForbiddenKey struct{}

// ExpectForbidden marks requests made with ctx as ones whose 403 the caller
// handles, so a 403 on them doesn't print the credentials hint.
func ExpectForbidden(ctx context.Context) context.Context {
	return context.WithValue(ctx, expectForbiddenKey{}, true)
}

func (h *authHint) check(req *http.Request, res *http.Response) {
	if res.StatusCode == http.StatusForbidden && req.Context().Value(expectForbiddenKey{}) != nil {
		return
	}
	var msg string
	switch res.StatusCode {
	case http.StatusUnauthorized:
		msg = fmt.Sprintf("the API rejected these credentials (HTTP 401) at %s. Check the key hasn't been revoked and belongs to this region: %s. Or run `langsmith auth login`.", req.URL.Host, regionAdvice)
	case http.StatusForbidden:
		if requiresWorkspace(res) {
			msg = "this API key can access more than one workspace, so the request must name one (HTTP 403). Run `langsmith workspaces list` to find its ID, then pass `--workspace <id>` or save it with `langsmith workspaces set-default <id>`."
			break
		}
		msg = fmt.Sprintf("the API refused this request (HTTP 403) at %s. Either the key lacks permission for this workspace or action (run `langsmith auth info` to see the profile and workspace in use, and `langsmith workspaces list` for your role in each), or it belongs to another region: %s.", req.URL.Host, regionAdvice)
	default:
		return
	}
	h.once.Do(func() { fmt.Fprintf(h.w, "hint: %s\n", msg) })
}

// orgScopedKeyRequiresWorkspace is the error code the API returns when an
// org-scoped key sends no workspace.
const orgScopedKeyRequiresWorkspace = "org_scoped_key_requires_workspace"

// requiresWorkspace reports whether a 403 says the request named no workspace.
// It reads the start of the body and puts it back for the caller.
func requiresWorkspace(res *http.Response) bool {
	if res.Body == nil {
		return false
	}
	head, _ := io.ReadAll(io.LimitReader(res.Body, 4096))
	res.Body = struct {
		io.Reader
		io.Closer
	}{io.MultiReader(bytes.NewReader(head), res.Body), res.Body}
	return bytes.Contains(head, []byte(orgScopedKeyRequiresWorkspace))
}

func (h *authHint) option() option.RequestOption {
	return option.WithMiddleware(func(req *http.Request, next option.MiddlewareNext) (*http.Response, error) {
		res, err := next(req)
		if err == nil {
			h.check(req, res)
		}
		return res, err
	})
}

type authHintTransport struct {
	hint *authHint
	base http.RoundTripper
}

func (t authHintTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	res, err := t.base.RoundTrip(req)
	if err == nil {
		t.hint.check(req, res)
	}
	return res, err
}

var processAuthHint = &authHint{w: os.Stderr}

// WrapDefaultTransportWithAuthHint adds the auth hint to requests sent through
// http.DefaultTransport, which generated commands use, until restore is called.
func WrapDefaultTransportWithAuthHint() (restore func()) {
	prev := http.DefaultTransport
	http.DefaultTransport = authHintTransport{hint: processAuthHint, base: prev}
	return func() { http.DefaultTransport = prev }
}
