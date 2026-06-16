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
		EnvVars: langsmith.F(map[string]string{
			"foo": "string",
		}),
		FsCapacityBytes: langsmith.F(int64(0)),
		IdleTtlSeconds:  langsmith.F(int64(0)),
		MemBytes:        langsmith.F(int64(0)),
		Mounts: langsmith.F([]langsmith.SandboxBoxNewParamsMountUnion{langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpec{
			ID:        langsmith.F("id"),
			MountPath: langsmith.F("mount_path"),
			S3: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecS3{
				Bucket:      langsmith.F("bucket"),
				EndpointURL: langsmith.F("endpoint_url"),
				Region:      langsmith.F("region"),
				PathStyle:   langsmith.F(true),
				Prefix:      langsmith.F("prefix"),
			}),
			Type: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecTypeS3),
			Cache: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecCache{
				MaxSizeBytes:     langsmith.F(int64(0)),
				WritebackSeconds: langsmith.F(int64(0)),
			}),
			Gcs: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGcs{
				Bucket: langsmith.F("bucket"),
				Prefix: langsmith.F("prefix"),
			}),
			Git: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGit{
				RemoteURL: langsmith.F("remote_url"),
				Ref: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRef{
					Name: langsmith.F("name"),
					Type: langsmith.F(langsmith.SandboxBoxNewParamsMountsSandboxapiS3BucketMountSpecGitRefTypeBranch),
				}),
				RefreshIntervalSeconds: langsmith.F(int64(1)),
			}),
			ReadOnly: langsmith.F(true),
		}}),
		Name: langsmith.F("name"),
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
				Name: langsmith.F("name"),
				Aws: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesAws{
					AccessKeyID: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyID{
						Type:  langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesAwsAccessKeyIDTypePlaintext),
						IsSet: langsmith.F(true),
						Value: langsmith.F("value"),
					}),
					SecretAccessKey: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKey{
						Type:  langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesAwsSecretAccessKeyTypePlaintext),
						IsSet: langsmith.F(true),
						Value: langsmith.F("value"),
					}),
				}),
				Enabled: langsmith.F(true),
				Gcp: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesGcp{
					Scopes: langsmith.F([]string{"string"}),
					ServiceAccountJson: langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJson{
						Type:  langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesGcpServiceAccountJsonTypePlaintext),
						IsSet: langsmith.F(true),
						Value: langsmith.F("value"),
					}),
				}),
				Headers: langsmith.F([]langsmith.SandboxBoxNewParamsProxyConfigRulesHeader{{
					Name:  langsmith.F("name"),
					Type:  langsmith.F(langsmith.SandboxBoxNewParamsProxyConfigRulesHeadersTypePlaintext),
					IsSet: langsmith.F(true),
					Value: langsmith.F("value"),
				}}),
				MatchHosts: langsmith.F([]string{"string"}),
				MatchPaths: langsmith.F([]string{"string"}),
				Type:       langsmith.F("type"),
			}}),
		}),
		RestoreMemory: langsmith.F(true),
		SnapshotID:    langsmith.F("snapshot_id"),
		SnapshotName:  langsmith.F("snapshot_name"),
		TagValueIDs:   langsmith.F([]string{"string"}),
		Vcpus:         langsmith.F(int64(0)),
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
					Name: langsmith.F("name"),
					Aws: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesAws{
						AccessKeyID: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyID{
							Type:  langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesAwsAccessKeyIDTypePlaintext),
							IsSet: langsmith.F(true),
							Value: langsmith.F("value"),
						}),
						SecretAccessKey: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKey{
							Type:  langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesAwsSecretAccessKeyTypePlaintext),
							IsSet: langsmith.F(true),
							Value: langsmith.F("value"),
						}),
					}),
					Enabled: langsmith.F(true),
					Gcp: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesGcp{
						Scopes: langsmith.F([]string{"string"}),
						ServiceAccountJson: langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJson{
							Type:  langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesGcpServiceAccountJsonTypePlaintext),
							IsSet: langsmith.F(true),
							Value: langsmith.F("value"),
						}),
					}),
					Headers: langsmith.F([]langsmith.SandboxBoxUpdateParamsProxyConfigRulesHeader{{
						Name:  langsmith.F("name"),
						Type:  langsmith.F(langsmith.SandboxBoxUpdateParamsProxyConfigRulesHeadersTypePlaintext),
						IsSet: langsmith.F(true),
						Value: langsmith.F("value"),
					}}),
					MatchHosts: langsmith.F([]string{"string"}),
					MatchPaths: langsmith.F([]string{"string"}),
					Type:       langsmith.F("type"),
				}}),
			}),
			TagValueIDs: langsmith.F([]string{"string"}),
			Vcpus:       langsmith.F(int64(0)),
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
		CreatedBy:     langsmith.F("created_by"),
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
			Name:            langsmith.F("name"),
			Checkpoint:      langsmith.F("checkpoint"),
			DockerImage:     langsmith.F("docker_image"),
			FsCapacityBytes: langsmith.F(int64(0)),
			IncludeMemory:   langsmith.F(true),
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
