package langsmith

import (
	"io"
	"net/http"
	"strings"

	"github.com/langchain-ai/langsmith-go/internal/requestconfig"
)

// withNoContentAsEmptyJSON makes a 204 No Content read as an empty JSON object
// when the response is decoded into a typed destination. Methods generated for
// an endpoint that returns 200 with a body decode the response; without this, an
// endpoint that moves to 204 would make them fail on the empty body.
//
// Raw destinations ([]byte, string, *http.Response) get the response exactly as
// sent. The destination is read when the request runs, since a caller's
// option.WithResponseBodyInto is applied after the client defaults.
func withNoContentAsEmptyJSON() requestconfig.RequestOption {
	return requestconfig.RequestOptionFunc(func(cfg *requestconfig.RequestConfig) error {
		cfg.Middlewares = append(cfg.Middlewares, func(req *http.Request, next func(*http.Request) (*http.Response, error)) (*http.Response, error) {
			res, err := next(req)
			if err != nil || res == nil || res.StatusCode != http.StatusNoContent || !decodesJSON(cfg.ResponseBodyInto) {
				return res, err
			}
			_ = res.Body.Close()
			res.Body = io.NopCloser(strings.NewReader("{}"))
			res.ContentLength = 2
			res.Header = res.Header.Clone()
			if res.Header == nil {
				res.Header = http.Header{}
			}
			res.Header.Set("Content-Type", "application/json")
			return res, nil
		})
		return nil
	})
}

// decodesJSON reports whether the request config decodes the response body as
// JSON into dst, rather than handing it over raw or not reading it.
func decodesJSON(dst any) bool {
	switch dst.(type) {
	case nil, *[]byte, *string, **string, **http.Response:
		return false
	}
	return true
}
