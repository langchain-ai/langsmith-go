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

// SandboxBrowserAttachmentTabService contains methods and other services that help
// with interacting with the langChain API.
//
// Beta: browsers are in beta, and these methods and types may change in
// incompatible ways before general availability.
//
// Note, unlike clients, this service does not read variables from the environment
// automatically. You should not instantiate this service directly, and instead use
// the [NewSandboxBrowserAttachmentTabService] method instead.
type SandboxBrowserAttachmentTabService struct {
	Options []option.RequestOption
}

// NewSandboxBrowserAttachmentTabService generates a new service that applies the
// given options to each request. These options are applied after the parent
// client's options (if there is one), and before any request-specific options.
func NewSandboxBrowserAttachmentTabService(opts ...option.RequestOption) (r *SandboxBrowserAttachmentTabService) {
	r = &SandboxBrowserAttachmentTabService{}
	r.Options = opts
	return
}

// Beta: open a tab in the window, at about:blank when no URL is given.
func (r *SandboxBrowserAttachmentTabService) New(ctx context.Context, name string, attachmentID string, body SandboxBrowserAttachmentTabNewParams, opts ...option.RequestOption) (res *BrowserTab, err error) {
	opts = slices.Concat(r.Options, opts)
	if name == "" {
		err = errors.New("missing required name parameter")
		return nil, err
	}
	if attachmentID == "" {
		err = errors.New("missing required attachment_id parameter")
		return nil, err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s/tabs", name, attachmentID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodPost, path, body, &res, opts...)
	return res, err
}

// Beta: close a tab in the window.
func (r *SandboxBrowserAttachmentTabService) Delete(ctx context.Context, name string, attachmentID string, tabID string, opts ...option.RequestOption) (err error) {
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
	if tabID == "" {
		err = errors.New("missing required tab_id parameter")
		return err
	}
	path := fmt.Sprintf("api/v2/sandboxes/browsers/%s/attachments/%s/tabs/%s", name, attachmentID, tabID)
	err = requestconfig.ExecuteNewRequest(ctx, http.MethodDelete, path, nil, nil, opts...)
	return err
}

type SandboxBrowserAttachmentTabNewParams struct {
	// URL to open; about:blank when empty.
	URL param.Field[string] `json:"url"`
}

func (r SandboxBrowserAttachmentTabNewParams) MarshalJSON() (data []byte, err error) {
	return apijson.MarshalRoot(r)
}
