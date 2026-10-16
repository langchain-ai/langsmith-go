// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SandboxBrowserAttachmentService contains methods and other services that help
// with interacting with the langChain API.
//
// Beta: browsers are in beta, and these methods and types may change in
// incompatible ways before general availability.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxBrowserAttachmentService] method instead.
type SandboxBrowserAttachmentService struct {
	Options []option.RequestOption
	Tabs    *SandboxBrowserAttachmentTabService
}

// NewSandboxBrowserAttachmentService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewSandboxBrowserAttachmentService(opts ...option.RequestOption) (r *SandboxBrowserAttachmentService) {
	r = &SandboxBrowserAttachmentService{}
	r.Options = opts
	r.Tabs = NewSandboxBrowserAttachmentTabService(opts...)
	return
}

// Beta: give a sandbox its own window of the browser, or, with no sandbox, open a
// window for the person calling (or, with controller "agent", for a client to
// drive over api_cdp_url). The response's cdp_url is the Chrome DevTools endpoint
// the sandbox connects to, and api_cdp_url the same endpoint through the API; a
// client can open more tabs, which stay in its window. Windows share the browser's
// cookies and logins. Keeping each sandbox to its own window is best effort, not a
// security boundary: attach only sandboxes you would trust with every login in the
// browser. With share_localhost, the browser can also reach the servers the
// sandbox runs, at the response's localhost_alias. Attaching takes, for the
// browser and for the attached sandbox each, being its creator or a role with
// sandboxes:exec. With unrestricted, nothing is filtered: the sandbox gets the
// whole browser, including people's windows, every site's cookies, and commands
// that close it. It needs nothing beyond what any attach needs.
func (r *SandboxBrowserAttachmentService) New(ctx context.Context, name string, body SandboxBrowserAttachmentNewParams, opts ...option.RequestOption) (res *BrowserAttachment, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Beta: close the window and stop the sandbox from making new connections to the
// browser.
func (r *SandboxBrowserAttachmentService) Delete(ctx context.Context, name string, attachmentID string, opts ...option.RequestOption) (err error) {
	opts = slices.Concat(r.Options, opts)
	opts = append([]option.RequestOption{option.WithHeader("Accept", "*/*")}, opts...)
	if name == "" {
		err = errors.New("missing required name parameter")
		return err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

// Beta: return control of the window to the attached sandbox.
func (r *SandboxBrowserAttachmentService) HandBack(ctx context.Context, name string, attachmentID string, opts ...option.RequestOption) (res *BrowserAttachment, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s/handback", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// Beta: give a person control of the window. Until it is handed back, the attached
// sandbox's commands that read or act on the page fail with a "human has control"
// error, and it receives no page content.
func (r *SandboxBrowserAttachmentService) TakeOver(ctx context.Context, name string, attachmentID string, opts ...option.RequestOption) (res *BrowserAttachment, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s/takeover", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

type SandboxBrowserAttachmentNewParams struct {
	// Controller is who drives a window opened with no sandbox: "human" (default) for
	// a person in the live view, or "agent" for a client that drives it over
	// api_cdp_url. A sandbox's window always starts with "agent".
	Controller param.Field[SandboxBrowserAttachmentNewParamsController] `json:"controller"`
	// Name labels the window in the live view; defaults to the sandbox name.
	Name param.Field[string] `json:"name"`
	// Sandbox is the name or ID of the sandbox to attach. Leave it out to open a
	// window for the person calling, to use from the live view.
	Sandbox param.Field[string] `json:"sandbox"`
	// ShareLocalhost shares the sandbox's localhost with the browser, so pages can
	// open the servers it runs, such as dev servers, at the response's
	// localhost_alias. It reaches every port of the sandbox, including servers bound
	// only to its localhost. Off by default. Not available for a window opened with no
	// sandbox, or when the sandbox is itself a browser.
	ShareLocalhost param.Field[bool] `json:"share_localhost"`
	// Unrestricted gives the sandbox the whole browser over CDP, with nothing
	// filtered: every window and tab including people's, every site's cookies, and
	// commands that close or crash the browser. Taking over its window does not pause
	// it.
	Unrestricted param.Field[bool] `json:"unrestricted"`
}

func (r SandboxBrowserAttachmentNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Controller is who drives a window opened with no sandbox: "human" (default) for
// a person in the live view, or "agent" for a client that drives it over
// api_cdp_url. A sandbox's window always starts with "agent".
type SandboxBrowserAttachmentNewParamsController string

const (
	SandboxBrowserAttachmentNewParamsControllerAgent SandboxBrowserAttachmentNewParamsController = "agent"
	SandboxBrowserAttachmentNewParamsControllerHuman SandboxBrowserAttachmentNewParamsController = "human"
)

func (r SandboxBrowserAttachmentNewParamsController) IsKnown() bool {
	switch r {
	case SandboxBrowserAttachmentNewParamsControllerAgent, SandboxBrowserAttachmentNewParamsControllerHuman:
		return true
	}
	return false
}
