package client

import (
	"bytes"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestAuthHintOption(t *testing.T) {
	// The API sends the code as a Problem Details extension member; older
	// servers send it in a plain JSON body.
	const workspaceProblem = `{"type":"about:blank","title":"Forbidden","status":403,"detail":"This API key is org-scoped and requires workspace specification.","instance":"/api/v1/sessions","error":"org_scoped_key_requires_workspace","message":"This API key is org-scoped and requires workspace specification."}`
	const workspaceBody = `{"error":"org_scoped_key_requires_workspace","message":"This API key is org-scoped and requires workspace specification."}`
	tests := []struct {
		name            string
		status          int
		contentType     string
		body            string
		expectForbidden bool
		want            string
	}{
		{name: "401", status: http.StatusUnauthorized, want: "rejected these credentials (HTTP 401)"},
		{name: "403", status: http.StatusForbidden, want: "run `langsmith auth info`"},
		{name: "403 without workspace", status: http.StatusForbidden, contentType: "application/problem+json", body: workspaceProblem, want: "Run `langsmith workspaces list`"},
		{name: "403 without workspace from an older server", status: http.StatusForbidden, body: workspaceBody, want: "Run `langsmith workspaces list`"},
		{name: "403 expected", status: http.StatusForbidden, expectForbidden: true},
		{name: "401 expected", status: http.StatusUnauthorized, expectForbidden: true, want: "rejected these credentials (HTTP 401)"},
		{name: "404", status: http.StatusNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			body := tt.body
			if body == "" {
				body = `{"detail":"Invalid token"}`
			}
			contentType := tt.contentType
			if contentType == "" {
				contentType = "application/json"
			}
			ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Type", contentType)
				w.WriteHeader(tt.status)
				_, _ = w.Write([]byte(body))
			}))
			defer ts.Close()
			var stderr bytes.Buffer
			hint := &authHint{w: &stderr}
			sdk := langsmith.NewClient(option.WithBaseURL(ts.URL), option.WithAPIKey("k"), option.WithMaxRetries(0), hint.option())

			ctx := context.Background()
			if tt.expectForbidden {
				ctx = ExpectForbidden(ctx)
			}
			for range 2 {
				_, err := sdk.Sessions.List(ctx, langsmith.SessionListParams{})
				// The hint reads the body, so the error must still carry it.
				if err == nil || !strings.Contains(err.Error(), body) {
					t.Fatalf("err = %v, want it to include the response body", err)
				}
			}

			wantCount := 0
			if tt.want != "" {
				wantCount = 1
			}
			if got := strings.Count(stderr.String(), "hint:"); got != wantCount || !strings.Contains(stderr.String(), tt.want) {
				t.Fatalf("stderr = %q, want %d hint containing %q", stderr.String(), wantCount, tt.want)
			}
		})
	}
}

func TestAuthHintTransport(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusUnauthorized)
	}))
	defer ts.Close()
	var stderr bytes.Buffer
	c := &http.Client{Transport: authHintTransport{hint: &authHint{w: &stderr}, base: http.DefaultTransport}}

	res, err := c.Get(ts.URL)
	if err != nil {
		t.Fatal(err)
	}
	res.Body.Close()

	if !strings.Contains(stderr.String(), "HTTP 401") {
		t.Fatalf("stderr = %q, want the 401 hint", stderr.String())
	}
}
