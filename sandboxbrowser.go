// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"time"

	"github.com/langchain-ai/langsmith-go/internal/apijson"
	"github.com/langchain-ai/langsmith-go/internal/param"
	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// SandboxBrowserService contains methods and other services that help with
// interacting with the langChain API.
//
// Beta: browsers are in beta, and these methods and types may change in
// incompatible ways before general availability.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxBrowserService] method instead.
type SandboxBrowserService struct {
	Options     []option.RequestOption
	Attachments *SandboxBrowserAttachmentService
}

// NewSandboxBrowserService generates a new service that applies the given options
// to each request. These options are applied after the parent client's options (if
// there is one), and before any request-specific options.
func NewSandboxBrowserService(opts ...option.RequestOption) (r *SandboxBrowserService) {
	r = &SandboxBrowserService{}
	r.Options = opts
	r.Attachments = NewSandboxBrowserAttachmentService(opts...)
	return
}

// Beta: create a browser sandbox, which runs one persistent Chrome that other
// sandboxes attach to. Logins made in it persist across restarts. It keeps its
// memory when it stops and is never deleted automatically.
func (r *SandboxBrowserService) New(ctx context.Context, body SandboxBrowserNewParams, opts ...option.RequestOption) (res *BrowserResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	path := "api/v2/sandboxes/browsers"
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Beta: get a browser and the windows attached to it.
func (r *SandboxBrowserService) Get(ctx context.Context, name string, opts ...option.RequestOption) (res *BrowserResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodGet, path, nil, &res, opts...)
	return res, err
}

// Beta: load the unpacked extensions placed in the browser's
// /var/lib/browserd/extensions (one per directory) that its Chrome does not have
// yet. Chrome loads them while running, so windows, tabs and connections are left
// as they are. Extensions already loaded stay loaded, even if removed from the
// directory, until Chrome next launches.
func (r *SandboxBrowserService) LoadExtensions(ctx context.Context, name string, opts ...option.RequestOption) (res *SandboxBrowserLoadExtensionsResponse, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/extensions/load", name)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, nil, &res, opts...)
	return res, err
}

// BrowserAttachment is one window of a browser (beta).
type BrowserAttachment struct {
	ID     string `json:"id"`
	ActURL string `json:"act_url"`
	// APICDPURL is the same Chrome DevTools endpoint through the LangSmith API, for
	// clients outside any sandbox. It takes any API credential: headers as for any API
	// call (X-Api-Key, or Authorization: Bearer with an OAuth access token, plus
	// X-Tenant-Id where needed), or, for clients that take only a URL, api_key, token
	// and tenant_id query parameters.
	APICdpURL string `json:"api_cdp_url"`
	// CDPURL is the Chrome DevTools endpoint the attached sandbox connects to, for
	// example with Playwright's connect_over_cdp. It works only from inside that
	// sandbox.
	CdpURL string `json:"cdp_url"`
	// Clients are the sandboxes currently connected to this window, with "api" for
	// each client connected through api_cdp_url. An unrestricted window's connections
	// are not listed.
	Clients []string `json:"clients"`
	// Controller is "agent" while the attached sandbox drives the window and "user"
	// while a person has taken it over.
	Controller    BrowserAttachmentController `json:"controller"`
	CreatedAt     string                      `json:"created_at"`
	Height        int64                       `json:"height"`
	LinksExpireAt time.Time                   `json:"links_expire_at" format:"date-time"`
	// LocalhostAlias is the sandbox's localhost as the browser sees it: open
	// http://<localhost_alias>:<port> in a window to reach a server listening on that
	// port in the sandbox, whether bound to localhost or to all interfaces (any port
	// but 80 and 443). While this is the only sandbox sharing its localhost,
	// http://localhost:<port> works too. Pages on other sites can navigate to a
	// plain-HTTP server there but cannot fetch from it. Set only for windows attached
	// with share_localhost.
	LocalhostAlias string `json:"localhost_alias"`
	Name           string `json:"name"`
	// Reason is why the agent asked a person to take over, when it did.
	Reason string `json:"reason"`
	// Sandbox is the attached sandbox; empty for a window a person opened.
	Sandbox string       `json:"sandbox"`
	Tabs    []BrowserTab `json:"tabs"`
	// Unrestricted windows give their sandbox the whole browser over CDP, with nothing
	// filtered: every window and tab, every site's cookies, and commands that close or
	// crash the browser. Taking over such a window does not pause its sandbox, and the
	// LSBrowserControl domain is not available.
	Unrestricted bool `json:"unrestricted"`
	// ViewURL and ActURL open the window on the sandbox service domain with no
	// LangSmith login, for an agent to send to a person: ViewURL to watch it, ActURL
	// also to use it while its controller is "user". Anyone holding a link can open
	// the window until LinksExpireAt. Each response mints new ones; they are returned
	// only to those who may use the browser.
	ViewURL string `json:"view_url"`
	// Width and Height are the window's page size in CSS pixels, when it was given
	// one; a window without one keeps Chrome's default of 1280x800.
	Width int64                 `json:"width"`
	JSON  browserAttachmentJSON `json:"-"`
}

// browserAttachmentJSON contains the JSON metadata for the struct
// [BrowserAttachment]
type browserAttachmentJSON struct {
	ID             apijson.Field
	ActURL         apijson.Field
	APICdpURL      apijson.Field
	CdpURL         apijson.Field
	Clients        apijson.Field
	Controller     apijson.Field
	CreatedAt      apijson.Field
	Height         apijson.Field
	LinksExpireAt  apijson.Field
	LocalhostAlias apijson.Field
	Name           apijson.Field
	Reason         apijson.Field
	Sandbox        apijson.Field
	Tabs           apijson.Field
	Unrestricted   apijson.Field
	ViewURL        apijson.Field
	Width          apijson.Field
	raw            string
	ExtraFields    map[string]apijson.Field
}

func (r *BrowserAttachment) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r browserAttachmentJSON) RawJSON() string {
	return r.raw
}

// Controller is "agent" while the attached sandbox drives the window and "user"
// while a person has taken it over.
type BrowserAttachmentController string

const (
	BrowserAttachmentControllerAgent BrowserAttachmentController = "agent"
	BrowserAttachmentControllerUser  BrowserAttachmentController = "user"
)

func (r BrowserAttachmentController) IsKnown() bool {
	switch r {
	case BrowserAttachmentControllerAgent, BrowserAttachmentControllerUser:
		return true
	}
	return false
}

// BrowserResponse describes a browser and the windows attached to it (beta).
type BrowserResponse struct {
	ID string `json:"id"`
	// Attachments are listed only while the browser is running.
	Attachments []BrowserAttachment `json:"attachments"`
	// HiDPI reports whether pages render at twice the resolution. It is reported only
	// while the browser is running.
	Hidpi  bool                `json:"hidpi"`
	Name   string              `json:"name"`
	Status string              `json:"status"`
	JSON   browserResponseJSON `json:"-"`
}

// browserResponseJSON contains the JSON metadata for the struct [BrowserResponse]
type browserResponseJSON struct {
	ID          apijson.Field
	Attachments apijson.Field
	Hidpi       apijson.Field
	Name        apijson.Field
	Status      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BrowserResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r browserResponseJSON) RawJSON() string {
	return r.raw
}

// BrowserTab is a tab in a browser window (beta).
type BrowserTab struct {
	ID    string         `json:"id"`
	Title string         `json:"title"`
	URL   string         `json:"url"`
	JSON  browserTabJSON `json:"-"`
}

// browserTabJSON contains the JSON metadata for the struct [BrowserTab]
type browserTabJSON struct {
	ID          apijson.Field
	Title       apijson.Field
	URL         apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *BrowserTab) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r browserTabJSON) RawJSON() string {
	return r.raw
}

type SandboxBrowserLoadExtensionsResponse struct {
	// Loaded are the ids of the extensions this call loaded; empty when Chrome already
	// had every extension in the directory.
	Loaded []string                                 `json:"loaded"`
	JSON   sandboxBrowserLoadExtensionsResponseJSON `json:"-"`
}

// sandboxBrowserLoadExtensionsResponseJSON contains the JSON metadata for the
// struct [SandboxBrowserLoadExtensionsResponse]
type sandboxBrowserLoadExtensionsResponseJSON struct {
	Loaded      apijson.Field
	raw         string
	ExtraFields map[string]apijson.Field
}

func (r *SandboxBrowserLoadExtensionsResponse) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r sandboxBrowserLoadExtensionsResponseJSON) RawJSON() string {
	return r.raw
}

type SandboxBrowserNewParams struct {
	// HiDPI renders pages at twice the resolution, so the live view is sharp on HiDPI
	// screens, at four times the pixels for the browser to render and stream. Off by
	// default, and fixed once the browser is created.
	Hidpi param.Field[bool] `json:"hidpi"`
	// IdleTTLSeconds stops the browser after this long without use. A stopped browser
	// resumes with its memory on the next connection. Defaults to the sandbox default;
	// 0 disables the idle stop.
	IdleTtlSeconds param.Field[int64] `json:"idle_ttl_seconds"`
	MemBytes       param.Field[int64] `json:"mem_bytes"`
	// Name of the browser sandbox. Generated when omitted, as "browser-" followed by
	// the browser's ID.
	Name param.Field[string] `json:"name"`
	// VCPUs and MemBytes size the browser sandbox, as for any sandbox. With neither
	// set, a browser gets 2 vCPUs and 8 GiB.
	Vcpus param.Field[int64] `json:"vcpus"`
}

func (r SandboxBrowserNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
