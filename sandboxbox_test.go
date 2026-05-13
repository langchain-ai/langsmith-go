// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"encoding/json"
	"errors"
	"os"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSandboxBoxNewParamsOmitEmptySnapshotID(t *testing.T) {
	params := langsmith.SandboxBoxNewParams{
		Name: langsmith.F("my-vm"),
	}

	raw, err := json.Marshal(params)
	require.NoError(t, err)

	var body map[string]any
	require.NoError(t, json.Unmarshal(raw, &body))
	assert.Equal(t, "my-vm", body["name"])
	assert.NotContains(t, body, "snapshot_id")
}

func TestSandboxBoxNewWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.New(context.TODO(), langsmith.SandboxBoxNewParams{
		DeleteAfterStopSeconds: langsmith.F(int64(0)),
		FsCapacityBytes:        langsmith.F(int64(0)),
		IdleTtlSeconds:         langsmith.F(int64(0)),
		MemBytes:               langsmith.F(int64(0)),
		Name:                   langsmith.F("name"),
		ProxyConfig: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfig{
			AccessControl: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigAccessControl{
				AllowList: langsmith.F([]string{"string"}),
				DenyList:  langsmith.F([]string{"string"}),
			}),
			Callbacks: langsmith.F([]langsmith.SandboxBoxNewParamsProxyConfigCallback{{
				MatchHosts:  langsmith.F([]string{"string"}),
				TtlSeconds:  langsmith.F(int64(60)),
				URL:         langsmith.F("url"),
				FullRequest: langsmith.F(true),
				RequestHeaders: langsmith.F([]langsmith.SandboxBoxNewParamsProxyConfigCallbacksRequestHeader{{
					Name:  langsmith.F("name"),
					Type:  langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigCallbacksRequestHeadersTypePlaintext),
					IsSet: langsmith.F(true),
					Value: langsmith.F("value"),
				}}),
			}}),
			NoProxy: langsmith.F([]string{"string"}),
			Rules: langsmith.F([]langsmith.SandboxBoxNewParamsProxyConfigRule{{
				MatchHosts: langsmith.F([]string{"string"}),
				Name:       langsmith.F("name"),
				Enabled:    langsmith.F(true),
				Headers: langsmith.F([]langsmith.SandboxBoxNewParamsProxyConfigRulesHeader{{
					Name:  langsmith.F("name"),
					Type:  langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesHeadersTypePlaintext),
					IsSet: langsmith.F(true),
					Value: langsmith.F("value"),
				}}),
				MatchPaths: langsmith.F([]string{"string"}),
			}}),
		}),
		SnapshotID:   langsmith.F("snapshot_id"),
		SnapshotName: langsmith.F("snapshot_name"),
		Vcpus:        langsmith.F(int64(0)),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxGet(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.Get(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxUpdateWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.Update(
		context.TODO(),
		"name",
		langsmith.SandboxBoxUpdateParams{
			DeleteAfterStopSeconds: langsmith.F(int64(0)),
			FsCapacityBytes:        langsmith.F(int64(0)),
			IdleTtlSeconds:         langsmith.F(int64(0)),
			MemBytes:               langsmith.F(int64(0)),
			Name:                   langsmith.F("name"),
			ProxyConfig: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfig{
				AccessControl: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigAccessControl{
					AllowList: langsmith.F([]string{"string"}),
					DenyList:  langsmith.F([]string{"string"}),
				}),
				Callbacks: langsmith.F([]langsmith.SandboxBoxUpdateParamsProxyConfigCallback{{
					MatchHosts:  langsmith.F([]string{"string"}),
					TtlSeconds:  langsmith.F(int64(60)),
					URL:         langsmith.F("url"),
					FullRequest: langsmith.F(true),
					RequestHeaders: langsmith.F([]langsmith.SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeader{{
						Name:  langsmith.F("name"),
						Type:  langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigCallbacksRequestHeadersTypePlaintext),
						IsSet: langsmith.F(true),
						Value: langsmith.F("value"),
					}}),
				}}),
				NoProxy: langsmith.F([]string{"string"}),
				Rules: langsmith.F([]langsmith.SandboxBoxUpdateParamsProxyConfigRule{{
					MatchHosts: langsmith.F([]string{"string"}),
					Name:       langsmith.F("name"),
					Enabled:    langsmith.F(true),
					Headers: langsmith.F([]langsmith.SandboxBoxUpdateParamsProxyConfigRulesHeader{{
						Name:  langsmith.F("name"),
						Type:  langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesHeadersTypePlaintext),
						IsSet: langsmith.F(true),
						Value: langsmith.F("value"),
					}}),
					MatchPaths: langsmith.F([]string{"string"}),
				}}),
			}),
			Vcpus: langsmith.F(int64(0)),
		},
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxListWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.List(context.TODO(), langsmith.SandboxBoxListParams{
		Limit:         langsmith.F(int64(0)),
		NameContains:  langsmith.F("name_contains"),
		Offset:        langsmith.F(int64(0)),
		SortBy:        langsmith.F("sort_by"),
		SortDirection: langsmith.F("sort_direction"),
		Status:        langsmith.F("status"),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxDelete(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	err := client.Sandboxes.Boxes.Delete(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxNewSnapshotWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.NewSnapshot(
		context.TODO(),
		"name",
		langsmith.SandboxBoxNewSnapshotParams{
			Name:       langsmith.F("name"),
			Checkpoint: langsmith.F("checkpoint"),
		},
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxGenerateServiceURLWithOptionalParams(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.GenerateServiceURL(
		context.TODO(),
		"name",
		langsmith.SandboxBoxGenerateServiceURLParams{
			ExpiresInSeconds: langsmith.F(int64(0)),
			Port:             langsmith.F(int64(0)),
		},
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxGetStatus(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.GetStatus(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxStart(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	_, err := client.Sandboxes.Boxes.Start(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBoxStop(t *testing.T) {
	t.Skip("Mock server tests are disabled")
	baseURL := "http://localhost:4010"
	if envURL, ok := os.LookupEnv("TEST_API_BASE_URL"); ok {
		baseURL = envURL
	}
	if !testutil.CheckTestServer(t, baseURL) {
		return
	}
	client := langsmith.NewClient(
		option.WithBaseURL(baseURL),
		option.WithAPIKey("My API Key"),
		option.WithTenantID("My Tenant ID"),
	)
	err := client.Sandboxes.Boxes.Stop(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
