// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package langsmith_test

import (
	"context"
	"errors"
	"os"
	"testing"

	"github.com/stainless-sdks/langsmith-api-go"
	"github.com/stainless-sdks/langsmith-api-go/internal/testutil"
	"github.com/stainless-sdks/langsmith-api-go/option"
)

func TestRepoNewWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Repos.New(context.TODO(), langsmith.RepoNewParams{
		IsPublic:    langsmith.F(true),
		RepoHandle:  langsmith.F("repo_handle"),
		Description: langsmith.F("description"),
		Readme:      langsmith.F("readme"),
		Tags:        langsmith.F([]string{"string"}),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRepoGet(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Repos.Get(
		context.TODO(),
		"owner",
		"repo",
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRepoUpdateWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Repos.Update(
		context.TODO(),
		"owner",
		"repo",
		langsmith.RepoUpdateParams{
			Description: langsmith.F("description"),
			IsArchived:  langsmith.F(true),
			IsPublic:    langsmith.F(true),
			Readme:      langsmith.F("readme"),
			Tags:        langsmith.F([]string{"string"}),
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

func TestRepoListWithOptionalParams(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Repos.List(context.TODO(), langsmith.RepoListParams{
		HasCommits:         langsmith.F(true),
		IsArchived:         langsmith.F(langsmith.RepoListParamsIsArchivedTrue),
		IsPublic:           langsmith.F(langsmith.RepoListParamsIsPublicTrue),
		Limit:              langsmith.F(int64(1)),
		Offset:             langsmith.F(int64(0)),
		Query:              langsmith.F("query"),
		SortDirection:      langsmith.F(langsmith.RepoListParamsSortDirectionAsc),
		SortField:          langsmith.F(langsmith.RepoListParamsSortFieldNumLikes),
		TagValueID:         langsmith.F([]string{"182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"}),
		Tags:               langsmith.F([]string{"string"}),
		TenantHandle:       langsmith.F("tenant_handle"),
		TenantID:           langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
		UpstreamRepoHandle: langsmith.F("upstream_repo_handle"),
		UpstreamRepoOwner:  langsmith.F("upstream_repo_owner"),
		WithLatestManifest: langsmith.F(true),
	})
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}

func TestRepoDelete(t *testing.T) {
	t.Skip("Prism tests are disabled")
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
		option.WithOrganizationID("My Organization ID"),
	)
	_, err := client.Repos.Delete(
		context.TODO(),
		"owner",
		"repo",
	)
	if err != nil {
		var apierr *langsmith.Error
		if errors.As(err, &apierr) {
			t.Log(string(apierr.DumpRequest(true)))
		}
		t.Fatalf("err should be nil: %s", err.Error())
	}
}
