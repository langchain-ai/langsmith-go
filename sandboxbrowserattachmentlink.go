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

// SandboxBrowserAttachmentLinkService contains methods and other services that
// help with interacting with the langChain API.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxBrowserAttachmentLinkService] method instead.
type SandboxBrowserAttachmentLinkService struct {
	Options []option.RequestOption
}

// NewSandboxBrowserAttachmentLinkService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewSandboxBrowserAttachmentLinkService(opts ...option.RequestOption) (r *SandboxBrowserAttachmentLinkService) {
	r = &SandboxBrowserAttachmentLinkService{}
	r.Options = opts
	return
}

// Beta: mint a view link and an act link for one window, as every attachment
// response does, with a chosen lifetime and the origins allowed to embed the
// window's pages in an iframe. Without embed_origins no page may frame them.
// view_embed_url and act_embed_url carry the same tokens in the URL fragment, for
// an iframe, so no server sees them. An embedded page sends its parent window
// messages when control changes and when it connects or disconnects (see the
// browser guide). Either link opens the window to anyone who holds it until it
// expires; detaching the window ends it early. A stopped browser starts when links
// are minted and when a link is opened, and resumes with its windows. Minting
// takes being the browser's creator or a role with sandboxes:exec.
func (r *SandboxBrowserAttachmentLinkService) New(ctx context.Context, name string, attachmentID string, body SandboxBrowserAttachmentLinkNewParams, opts ...option.RequestOption) (res *BrowserWindowLinks, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s/links", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

type BrowserWindowLinks struct {
	ActEmbedURL string `json:"act_embed_url"`
	// ActURL also lets the person take control of the window, use it, and hand it
	// back.
	ActURL    string    `json:"act_url"`
	ExpiresAt time.Time `json:"expires_at" format:"date-time"`
	// ViewEmbedURL and ActEmbedURL open the same window in an iframe with the token in
	// the URL fragment, which no server, log or Referer header sees. Viewer options
	// follow the token in the fragment, as in {act_embed_url}&fit=follow&zoom=1.
	ViewEmbedURL string `json:"view_embed_url"`
	// ViewURL shows the window and takes no input.
	ViewURL string                 `json:"view_url"`
	JSON    browserWindowLinksJSON `json:"-"`
}

// browserWindowLinksJSON contains the JSON metadata for the struct
// [BrowserWindowLinks]
type browserWindowLinksJSON struct {
	ActEmbedURL  apijson.Field
	ActURL       apijson.Field
	ExpiresAt    apijson.Field
	ViewEmbedURL apijson.Field
	ViewURL      apijson.Field
	raw          string
	ExtraFields  map[string]apijson.Field
}

func (r *BrowserWindowLinks) UnmarshalJSON(data []byte) (err error) {
	return apijson.UnmarshalRoot(data, r)
}

func (r browserWindowLinksJSON) RawJSON() string {
	return r.raw
}

type SandboxBrowserAttachmentLinkNewParams struct {
	// EmbedOrigins are the origins whose pages may show the window's pages in an
	// iframe, such as https://app.example.com: https, or http on localhost, with no
	// path. Without them no page may frame the window.
	EmbedOrigins param.Field[[]string] `json:"embed_origins"`
	// TTLSeconds is how long the links work, from 60 seconds to 24 hours. Defaults to
	// 24 hours.
	TtlSeconds param.Field[int64] `json:"ttl_seconds"`
}

func (r SandboxBrowserAttachmentLinkNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
