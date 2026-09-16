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

func TestRepoDirectoryListWithOptionalParams(t *testing.T) {
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
	_, err := client.Repos.Directories.List(
		context.TODO(),
		"owner",
		"repo",
		langsmith.RepoDirectoryListParams{
			Commit: langsmith.F("commit"),
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

func TestRepoDirectoryDeleteWithOptionalParams(t *testing.T) {
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
	err := client.Repos.Directories.Delete(
		context.TODO(),
		"owner",
		"repo",
		langsmith.RepoDirectoryDeleteParams{
			RepoType: langsmith.F(langsmith.RepoDirectoryDeleteParamsRepoTypeAgent),
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

func TestRepoDirectoryCommitWithOptionalParams(t *testing.T) {
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
	_, err := client.Repos.Directories.Commit(
		context.TODO(),
		"owner",
		"repo",
		langsmith.RepoDirectoryCommitParams{
			Files: langsmith.F(map[string]langsmith.RepoDirectoryCommitParamsFilesUnion{
				"agents/pinned": langsmith.RepoDirectoryCommitParamsFilesDirectoryAgentEntryInput{
					RepoHandle: langsmith.F("review-agent"),
					Type:       langsmith.F(langsmith.RepoDirectoryCommitParamsFilesDirectoryAgentEntryInputTypeAgent),
					CommitID:   langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
					Selector: langsmith.F[langsmith.DirectorySelectorUnionParam](langsmith.DirectorySelectorDirectoryCommitSelectorParam{
						CommitID: langsmith.F("0198f3ab-7c2d-7def-8a91-23456789abcd"),
						Type:     langsmith.F(langsmith.DirectorySelectorDirectoryCommitSelectorTypeCommit),
					}),
				},
				"skills/current": langsmith.RepoDirectoryCommitParamsFilesDirectorySkillEntryInput{
					RepoHandle: langsmith.F("shared-skill"),
					Type:       langsmith.F(langsmith.RepoDirectoryCommitParamsFilesDirectorySkillEntryInputTypeSkill),
					CommitID:   langsmith.F("182bd5e5-6e1a-4fe4-a799-aa6d9a6ab26e"),
					Selector: langsmith.F[langsmith.DirectorySelectorUnionParam](langsmith.DirectorySelectorDirectoryLatestSelectorParam{
						Type: langsmith.F(langsmith.DirectorySelectorDirectoryLatestSelectorTypeLatest),
					}),
				},
			}),
			ParentCommit: langsmith.F("parent_commit"),
			SkipWebhooks: langsmith.F(true),
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
