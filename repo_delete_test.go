package langsmith_test

import (
	"context"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

// The repo delete endpoint returns 200 with a JSON body today and will return
// 204 No Content; Repos.Delete has to succeed on both.
func TestRepoDeleteAcceptsJSONAndNoContent(t *testing.T) {
	tests := []struct {
		name  string
		write func(w http.ResponseWriter)
	}{
		{
			name: "200 with JSON",
			write: func(w http.ResponseWriter) {
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusOK)
				_, _ = w.Write([]byte(`{}`))
			},
		},
		{
			name: "204 No Content",
			write: func(w http.ResponseWriter) {
				w.WriteHeader(http.StatusNoContent)
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if r.Method != http.MethodDelete || r.URL.Path != "/api/v1/repos/owner/repo" {
					http.Error(w, "unexpected request", http.StatusBadRequest)
					return
				}
				tt.write(w)
			}))
			defer srv.Close()

			client := langsmith.NewClient(
				option.WithBaseURL(srv.URL),
				option.WithAPIKey("test-api-key"),
				option.WithMaxRetries(0),
			)
			if _, err := client.Repos.Delete(context.Background(), "owner", "repo"); err != nil {
				t.Fatalf("Repos.Delete: %v", err)
			}
		})
	}
}

// Raw destinations get a 204 or empty 200 exactly as sent, not a substituted
// JSON body.
func TestNoContentKeepsRawResponses(t *testing.T) {
	for _, status := range []int{http.StatusOK, http.StatusNoContent} {
		t.Run(http.StatusText(status), func(t *testing.T) {
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(status)
			}))
			defer srv.Close()

			client := langsmith.NewClient(
				option.WithBaseURL(srv.URL),
				option.WithAPIKey("test-api-key"),
				option.WithMaxRetries(0),
			)

			got := []byte("unset")
			if err := client.Get(context.Background(), "api/v1/file", nil, &got); err != nil {
				t.Fatalf("Get into []byte: %v", err)
			}
			if len(got) != 0 {
				t.Fatalf("[]byte body = %q, want empty", got)
			}

			var res *http.Response
			if _, err := client.Repos.Delete(context.Background(), "owner", "repo", option.WithResponseBodyInto(&res)); err != nil {
				t.Fatalf("Repos.Delete into *http.Response: %v", err)
			}
			body, err := io.ReadAll(res.Body)
			_ = res.Body.Close()
			if err != nil {
				t.Fatalf("reading raw body: %v", err)
			}
			if len(body) != 0 {
				t.Fatalf("raw body = %q, want empty", body)
			}
		})
	}
}
