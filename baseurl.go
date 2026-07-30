package langsmith

import (
	"net/url"
	"strings"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
	"github.com/langchain-ai/langsmith-go/option"
)

// apiPathPrefixes are the trailing path segments stripped from a configured base
// URL, longest first. Most generated request paths carry their own prefix
// ("api/v1/runs/query", "v2/runs/query"), so the base URL has to point at the
// deployment root. Users often configure LANGSMITH_ENDPOINT with the suffixed URL
// the Python and JS SDKs accept (https://api.smith.langchain.com/api/v1), which
// would otherwise double up.
var apiPathPrefixes = []string{"/api/v1", "/api"}

// normalizeBaseURL strips a trailing /api/v1 or /api path segment from a base URL.
//
// Only whole trailing path segments of the parsed URL are considered, so a host
// such as api.example.com is left alone. The scheme and host are never touched.
func normalizeBaseURL(base string) string {
	u, err := url.Parse(base)
	if err != nil {
		return base
	}
	if trimBaseURLPath(u) {
		return u.String()
	}
	return base
}

// trimBaseURLPath strips an API path prefix from u in place, reporting whether it
// changed anything.
func trimBaseURLPath(u *url.URL) bool {
	path := strings.TrimRight(u.Path, "/")
	for _, prefix := range apiPathPrefixes {
		if !strings.HasSuffix(path, prefix) {
			continue
		}
		trimmed := strings.TrimSuffix(path, prefix)
		// A base URL with a path keeps its trailing slash so that resolving a
		// relative request path against it does not drop the last segment.
		if trimmed != "" {
			trimmed += "/"
		}
		u.Path = trimmed
		// RawPath is only meaningful when it differs from Path; clearing it keeps
		// the two consistent after the rewrite.
		u.RawPath = ""
		return true
	}
	return false
}

// baseURLNormalizer strips a trailing /api/v1 or /api from whatever base URL the
// preceding options resolved to. It has to be applied last, after both
// [DefaultClientOptions] and any caller-supplied options.
//
// It is a named type rather than a [requestconfig.RequestOptionFunc] so that
// [withoutBaseURLNormalization] can identify and drop it.
type baseURLNormalizer struct{}

func (baseURLNormalizer) Apply(r *requestconfig.RequestConfig) error {
	if r.BaseURL == nil {
		return nil
	}
	// Copy rather than trim in place: option.WithBaseURL parses the URL once and
	// hands the same pointer to every client it is applied to.
	trimmed := *r.BaseURL
	if trimBaseURLPath(&trimmed) {
		r.BaseURL = &trimmed
	}
	return nil
}

func withNormalizedBaseURL() option.RequestOption { return baseURLNormalizer{} }

// withoutBaseURLNormalization returns opts with the base URL normalization
// dropped, yielding the base URL exactly as it was configured.
//
// Trace ingest needs this. The multipart exporter appends root-relative paths
// ("/runs/multipart", "/runs/batch") to the URL it is given, and the prefix those
// live behind is deployment specific — /api on self-hosted, nothing on SaaS. It
// cannot be reconstructed from the normalized root, so ingest keeps using the URL
// the caller configured.
func withoutBaseURLNormalization(opts []option.RequestOption) []option.RequestOption {
	kept := make([]option.RequestOption, 0, len(opts))
	for _, opt := range opts {
		if _, ok := opt.(baseURLNormalizer); ok {
			continue
		}
		kept = append(kept, opt)
	}
	return kept
}
