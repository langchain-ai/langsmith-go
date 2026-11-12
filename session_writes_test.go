package langsmith_test

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSessionWriteContracts(t *testing.T) {
	for _, operation := range []string{"create", "update", "createV2", "updateV2"} {
		t.Run(operation, func(t *testing.T) {
			v2 := strings.HasSuffix(operation, "V2")
			create := strings.HasPrefix(operation, "create")
			path, tier, method := "/api/v1/sessions", "longlived", http.MethodPost
			if v2 {
				path, tier = "/api/v2/sessions", "LONG_LIVED"
			}
			if !create {
				path += "/project"
				method = http.MethodPatch
			}
			calls := 0
			client := langsmith.NewClient(
				option.WithBaseURL("https://example.com"),
				option.WithAPIKey("test-key"),
				option.WithHTTPClient(&http.Client{Transport: &closureTransport{
					fn: func(req *http.Request) (*http.Response, error) {
						calls++
						if req.Method != method || req.URL.Path != path {
							t.Errorf("request = %s %s, want %s %s", req.Method, req.URL.Path, method, path)
						}
						var body map[string]any
						if err := json.NewDecoder(req.Body).Decode(&body); err != nil {
							t.Fatal(err)
						}
						want := map[string]any{"name": "Example", "trace_tier": tier}
						if v2 && create {
							want["type"] = "PROJECT"
						}
						if !reflect.DeepEqual(body, want) {
							t.Errorf("body = %#v, want %#v", body, want)
						}
						return &http.Response{
							StatusCode: http.StatusOK,
							Header:     http.Header{"Content-Type": {"application/json"}},
							Body:       io.NopCloser(strings.NewReader(`{"id":"project","name":"Example","type":"PROJECT","trace_tier":"` + tier + `"}`)),
						}, nil
					},
				}}),
			)
			var gotTier string
			switch operation {
			case "create":
				result, err := client.Sessions.New(context.Background(), langsmith.SessionNewParams{
					Name: langsmith.F("Example"), TraceTier: langsmith.F(langsmith.SessionNewParamsTraceTierLonglived),
				})
				if err != nil {
					t.Fatal(err)
				}
				gotTier = string(result.TraceTier)
			case "update":
				result, err := client.Sessions.Update(context.Background(), "project", langsmith.SessionUpdateParams{
					Name: langsmith.F("Example"), TraceTier: langsmith.F(langsmith.SessionUpdateParamsTraceTierLonglived),
				})
				if err != nil {
					t.Fatal(err)
				}
				gotTier = string(result.TraceTier)
			case "createV2":
				result, err := client.Sessions.NewV2(context.Background(), langsmith.SessionNewV2Params{
					Type: langsmith.F(langsmith.SessionNewV2ParamsTypeProject), Name: langsmith.F("Example"), TraceTier: langsmith.F(langsmith.SessionNewV2ParamsTraceTierLongLived),
				})
				if err != nil {
					t.Fatal(err)
				}
				gotTier = string(result.TraceTier)
			case "updateV2":
				result, err := client.Sessions.UpdateV2(context.Background(), "project", langsmith.SessionUpdateV2Params{
					Name: langsmith.F("Example"), TraceTier: langsmith.F(langsmith.SessionUpdateV2ParamsTraceTierLongLived),
				})
				if err != nil {
					t.Fatal(err)
				}
				gotTier = string(result.TraceTier)
			}
			if calls != 1 || gotTier != tier {
				t.Errorf("calls = %d, response tier = %q; want 1, %q", calls, gotTier, tier)
			}
		})
	}
}
