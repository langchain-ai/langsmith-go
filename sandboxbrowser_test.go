// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/internal/testutil"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxBrowserNewWithOptionalParams(t *testing.T) {
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
	_, err := client.Sandboxes.Browsers.New(context.TODO(), langsmith.SandboxBrowserNewParams{
		Hidpi:          langsmith.F(true),
		IdleTtlSeconds: langsmith.F(int64(0)),
		MemBytes:       langsmith.F(int64(0)),
		Name:           langsmith.F("name"),
		ProxyConfig: langsmith.F(langsmith.ProxyConfigParam{
			AccessControl: langsmith.F(langsmith.ProxyConfigAccessControlParam{
				AllowList: langsmith.F([]string{"string"}),
				DenyList:  langsmith.F([]string{"string"}),
			}),
			Callbacks: langsmith.F([]langsmith.ProxyConfigCallbackParam{{
				MatchHosts:  langsmith.F([]string{"string"}),
				TtlSeconds:  langsmith.F(int64(60)),
				URL:         langsmith.F("url"),
				FullRequest: langsmith.F(true),
				RequestHeaders: langsmith.F([]langsmith.ProxyConfigCallbacksRequestHeaderParam{{
					Name:  langsmith.F("name"),
					Type:  langsmith.F(langsmith.ProxyConfigCallbacksRequestHeadersTypePlaintext),
					IsSet: langsmith.F(true),
					Value: langsmith.F("value"),
				}}),
			}}),
			Description: langsmith.F("description"),
			NoProxy:     langsmith.F([]string{"string"}),
			Rules: langsmith.F([]langsmith.ProxyConfigRuleParam{{
				Name: langsmith.F("name"),
				Aws: langsmith.F[langsmith.ProxyConfigRulesAwsUnionParam](langsmith.ProxyConfigRulesAwsSandboxesProxyAwsRoleConfigParam{
					RoleArn:      langsmith.F("x"),
					SubjectScope: langsmith.F("tag:2fdcD1Dc-bbDb-2BBD-0Afa-1A3C33cFaADc"),
				}),
				Azure: langsmith.F[langsmith.ProxyConfigRulesAzureUnionParam](langsmith.ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigParam{
					ClientID:     langsmith.F("E1CB97d8-Ddf3-aAae-0Afa-1A3C33cFaADc"),
					TenantID:     langsmith.F("E1CB97d8-Ddf3-aAae-0Afa-1A3C33cFaADc"),
					Resources:    langsmith.F([]langsmith.ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResource{langsmith.ProxyConfigRulesAzureSandboxesProxyAzureWorkloadIdentityConfigResourceHTTPSStorageAzureCom}),
					SubjectScope: langsmith.F("tag:2fdcD1Dc-bbDb-2BBD-0Afa-1A3C33cFaADc"),
				}),
				Description: langsmith.F("description"),
				Enabled:     langsmith.F(true),
				EnvVars: langsmith.F(map[string]string{
					"foo": "string",
				}),
				Gcp: langsmith.F[langsmith.ProxyConfigRulesGcpUnionParam](langsmith.ProxyConfigRulesGcpSandboxesProxyGcpWorkloadIdentityConfigParam{
					WorkloadIdentityProvider: langsmith.F("x"),
					Scopes:                   langsmith.F([]string{"string"}),
					ServiceAccountEmail:      langsmith.F("service_account_email"),
					SubjectScope:             langsmith.F("tag:2fdcD1Dc-bbDb-2BBD-0Afa-1A3C33cFaADc"),
				}),
				Headers: langsmith.F([]langsmith.ProxyConfigRulesHeaderParam{{
					Name:  langsmith.F("name"),
					Type:  langsmith.F(langsmith.ProxyConfigRulesHeadersTypePlaintext),
					IsSet: langsmith.F(true),
					Value: langsmith.F("value"),
				}}),
				MatchHeaders: langsmith.F([]string{"string"}),
				MatchHosts:   langsmith.F([]string{"string"}),
				MatchPaths:   langsmith.F([]string{"string"}),
				RelayID:      langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
				Type:         langsmith.F("type"),
			}}),
		}),
		Vcpus: langsmith.F(int64(0)),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBrowserGet(t *testing.T) {
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
	_, err := client.Sandboxes.Browsers.Get(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestSandboxBrowserLoadExtensions(t *testing.T) {
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
	_, err := client.Sandboxes.Browsers.LoadExtensions(context.TODO(), "name")
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
