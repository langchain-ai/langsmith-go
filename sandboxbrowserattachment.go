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
// sandboxes:exec. With unrestricted, nothing is filtered: the window's CDP clients
// (its sandbox, or anyone using its api_cdp_url) get the whole browser, including
// people's windows, every site's cookies, and commands that close it. It needs
// nothing beyond what any attach needs.
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

// Beta: size a window so its page is width by height CSS pixels. The window keeps
// that size until resized again, through this call or over CDP with
// Browser.setWindowBounds; the live view scales to fit it rather than resizing it.
// A browser that is not running cannot be resized.
func (r *SandboxBrowserAttachmentService) Update(ctx context.Context, name string, attachmentID string, body SandboxBrowserAttachmentUpdateParams, opts ...option.RequestOption) (res *BrowserAttachment, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPatch, path, body, &res, opts...)
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

// Beta: give the window to a person ("user") or back to the attached sandbox
// ("agent"). It takes effect at once and does not wait for anyone. While a person
// has control, the attached sandbox's commands that read or act on the page fail
// with a "human has control" error, and it receives no page content. The window's
// current controller is in the browser's attachments.
func (r *SandboxBrowserAttachmentService) SetController(ctx context.Context, name string, attachmentID string, body SandboxBrowserAttachmentSetControllerParams, opts ...option.RequestOption) (res *BrowserAttachment, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s/controller", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type SandboxBrowserAttachmentNewParams struct {
	// Controller is who drives a window opened with no sandbox: "user" (default) for a
	// person in the live view, or "agent" for a client that drives it over
	// api_cdp_url. A sandbox's window always starts with "agent".
	Controller param.Field[SandboxBrowserAttachmentNewParamsController] `json:"controller"`
	Height     param.Field[int64]                                       `json:"height"`
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
	// Unrestricted gives the window's CDP clients (its sandbox, or anyone using its
	// api_cdp_url) the whole browser, with nothing filtered: every window and tab
	// including people's, every site's cookies, and commands that close or crash the
	// browser. Taking over its window does not pause them.
	Unrestricted param.Field[bool] `json:"unrestricted"`
	// Width and Height size the window's page in CSS pixels; set both or neither.
	// Without them the window opens at 1280x800.
	Width param.Field[int64] `json:"width"`
}

func (r SandboxBrowserAttachmentNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Controller is who drives a window opened with no sandbox: "user" (default) for a
// person in the live view, or "agent" for a client that drives it over
// api_cdp_url. A sandbox's window always starts with "agent".
type SandboxBrowserAttachmentNewParamsController string

const (
	SandboxBrowserAttachmentNewParamsControllerAgent SandboxBrowserAttachmentNewParamsController = "agent"
	SandboxBrowserAttachmentNewParamsControllerUser  SandboxBrowserAttachmentNewParamsController = "user"
)

func (r SandboxBrowserAttachmentNewParamsController) IsKnown() bool {
	switch r {
	case SandboxBrowserAttachmentNewParamsControllerAgent, SandboxBrowserAttachmentNewParamsControllerUser:
		return true
	}
	return false
}

type SandboxBrowserAttachmentUpdateParams struct {
	Height param.Field[int64] `json:"height"`
	// Width and Height are the window's new page size in CSS pixels.
	Width param.Field[int64] `json:"width"`
}

func (r SandboxBrowserAttachmentUpdateParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

type SandboxBrowserAttachmentSetControllerParams struct {
	// Controller is "user" to give the window to a person, or "agent" to give it back
	// to the attached sandbox.
	Controller param.Field[SandboxBrowserAttachmentSetControllerParamsController] `json:"controller"`
}

func (r SandboxBrowserAttachmentSetControllerParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}

// Controller is "user" to give the window to a person, or "agent" to give it back
// to the attached sandbox.
type SandboxBrowserAttachmentSetControllerParamsController string

const (
	SandboxBrowserAttachmentSetControllerParamsControllerAgent SandboxBrowserAttachmentSetControllerParamsController = "agent"
	SandboxBrowserAttachmentSetControllerParamsControllerUser  SandboxBrowserAttachmentSetControllerParamsController = "user"
)

func (r SandboxBrowserAttachmentSetControllerParamsController) IsKnown() bool {
	switch r {
	case SandboxBrowserAttachmentSetControllerParamsControllerAgent, SandboxBrowserAttachmentSetControllerParamsControllerUser:
		return true
	}
	return false
}
