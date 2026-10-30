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
		AccessDelegation: langsmith.F(langsmith.SandboxBoxNewParamsAccessDelegation{
			Mode:        langsmith.F(langsmith.SandboxBoxNewParamsAccessDelegationModeInherit),
			Permissions: langsmith.F([]string{"string"}),
		}),
		CPUMillicores:          langsmith.F(int64(0)),
		DeleteAfterStopSeconds: langsmith.F(int64(0)),
		EnvVars: langsmith.F(map[string]string{
			"foo": "string",
		}),
		FsCapacityBytes: langsmith.F(int64(0)),
		IdleTtlSeconds:  langsmith.F(int64(0)),
		Labels: langsmith.F(map[string]string{
			"foo": "string",
		}),
		MemBytes: langsmith.F(int64(0)),
		MountConfig: langsmith.F(langsmith.SandboxMountConfigParam{
			Auth: langsmith.F(langsmith.SandboxMountConfigAuthParam{
				Aws: langsmith.F[langsmith.SandboxMountConfigAuthAwsUnionParam](langsmith.SandboxMountConfigAuthAwsSandboxesSandboxAwsMountRoleAuthConfigParam{
					RoleArn:      langsmith.F("x"),
					SubjectScope: langsmith.F("tag:2fdcD1Dc-bbDb-2BBD-0Afa-1A3C33cFaADc"),
				}),
				Azure: langsmith.F[langsmith.SandboxMountConfigAuthAzureUnionParam](langsmith.SandboxMountConfigAuthAzureSandboxesSandboxAzureMountWorkloadIdentityAuthConfigParam{
					ClientID:     langsmith.F("E1CB97d8-Ddf3-aAae-0Afa-1A3C33cFaADc"),
					TenantID:     langsmith.F("E1CB97d8-Ddf3-aAae-0Afa-1A3C33cFaADc"),
					SubjectScope: langsmith.F("tag:2fdcD1Dc-bbDb-2BBD-0Afa-1A3C33cFaADc"),
				}),
				Gcp: langsmith.F[langsmith.SandboxMountConfigAuthGcpUnionParam](langsmith.SandboxMountConfigAuthGcpSandboxesSandboxGcpMountWorkloadIdentityAuthConfigParam{
					WorkloadIdentityProvider: langsmith.F("x"),
					ServiceAccountEmail:      langsmith.F("service_account_email"),
					SubjectScope:             langsmith.F("tag:2fdcD1Dc-bbDb-2BBD-0Afa-1A3C33cFaADc"),
				}),
			}),
			Mounts: langsmith.F([]langsmith.SandboxMountConfigMountsUnionParam{langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecParam{
				ID:        langsmith.F("id"),
				MountPath: langsmith.F("mount_path"),
				S3: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecS3Param{
					Bucket:      langsmith.F("bucket"),
					Region:      langsmith.F("region"),
					EndpointURL: langsmith.F("endpoint_url"),
					PathStyle:   langsmith.F(true),
					Prefix:      langsmith.F("prefix"),
				}),
				Type: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecTypeS3),
				AzureBlob: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecAzureBlobParam{
					Container:      langsmith.F("container"),
					StorageAccount: langsmith.F("26f1klzn5713"),
					Prefix:         langsmith.F("prefix"),
				}),
				Cache: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecCacheParam{
					MaxSizeBytes:     langsmith.F(int64(0)),
					WritebackSeconds: langsmith.F(int64(0)),
				}),
				Contexthub: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecContexthubParam{
					Repo:            langsmith.F("repo"),
					InitialPullOnly: langsmith.F(true),
				}),
				Encryption: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecEncryptionParam{
					KmsKey: langsmith.F("kms_key"),
				}),
				Gcs: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecGcsParam{
					Bucket: langsmith.F("bucket"),
					Prefix: langsmith.F("prefix"),
				}),
				Git: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitParam{
					RemoteURL: langsmith.F("remote_url"),
					Ref: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefParam{
						Name: langsmith.F("name"),
						Type: langsmith.F(langsmith.SandboxMountConfigMountsSandboxapiS3BucketMountSpecGitRefTypeBranch),
					}),
					RefreshIntervalSeconds: langsmith.F(int64(1)),
				}),
				ReadOnly: langsmith.F(true),
			}}),
		}),
		Name:                 langsmith.F("name"),
		Persistent:           langsmith.F(true),
		PreserveMemoryOnStop: langsmith.F(true),
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
					RoleArn: langsmith.F("x"),
					SessionTags: langsmith.F(map[string]string{
						"foo": "string",
					}),
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
				Rewrite: langsmith.F(langsmith.ProxyConfigRulesRewriteParam{
					Port:        langsmith.F(int64(1)),
					Credentials: langsmith.F(langsmith.ProxyConfigRulesRewriteCredentialsStrip),
					Host:        langsmith.F("host"),
					Sandbox:     langsmith.F("sandbox"),
				}),
				Type: langsmith.F("type"),
			}}),
		}),
		RestoreMemory: langsmith.F(true),
		RunConfig: langsmith.F(langsmith.SandboxBoxNewParamsRunConfig{
			EnvVars: langsmith.F(map[string]string{
				"foo": "string",
			}),
			User:    langsmith.F("user"),
			WorkDir: langsmith.F("work_dir"),
		}),
		Snapshot:     langsmith.F("snapshot"),
		SnapshotID:   langsmith.F("snapshot_id"),
		SnapshotName: langsmith.F("snapshot_name"),
		TagValueIDs:  langsmith.F([]string{"string"}),
		Vcpus:        langsmith.F(int64(0)),
		Volumes: langsmith.F([]langsmith.SandboxBoxNewParamsVolume{{
			MountPath: langsmith.F("mount_path"),
			Volume:    langsmith.F("volume"),
			ReadOnly:  langsmith.F(true),
		}}),
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
			CPUMillicores:          langsmith.F(int64(0)),
			DeleteAfterStopSeconds: langsmith.F(int64(0)),
			FsCapacityBytes:        langsmith.F(int64(0)),
			IdleTtlSeconds:         langsmith.F(int64(0)),
			MemBytes:               langsmith.F(int64(0)),
			Name:                   langsmith.F("name"),
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
						RoleArn: langsmith.F("x"),
						SessionTags: langsmith.F(map[string]string{
							"foo": "string",
						}),
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
					Rewrite: langsmith.F(langsmith.ProxyConfigRulesRewriteParam{
						Port:        langsmith.F(int64(1)),
						Credentials: langsmith.F(langsmith.ProxyConfigRulesRewriteCredentialsStrip),
						Host:        langsmith.F("host"),
						Sandbox:     langsmith.F("sandbox"),
					}),
					Type: langsmith.F("type"),
				}}),
			}),
			RunConfig: langsmith.F(langsmith.SandboxBoxUpdateParamsRunConfig{
				EnvVars: langsmith.F(map[string]string{
					"foo": "string",
				}),
				User:    langsmith.F("user"),
				WorkDir: langsmith.F("work_dir"),
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
		Cursor:        langsmith.F("cursor"),
		Label:         langsmith.F([]string{"string"}),
		Limit:         langsmith.F(int64(0)),
		NameContains:  langsmith.F("name_contains"),
		Offset:        langsmith.F(int64(0)),
		PageSize:      langsmith.F(int64(0)),
		SortBy:        langsmith.F("sort_by"),
		SortDirection: langsmith.F("sort_direction"),
		SortOrder:     langsmith.F("sort_order"),
		Status:        langsmith.F("status"),
		TagValueID:    langsmith.F([]string{"string"}),
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
			Name:                             langsmith.F("name"),
			Checkpoint:                       langsmith.F("checkpoint"),
			ConfirmAllowSnapshotSecretValues: langsmith.F(true),
			Description:                      langsmith.F("description"),
			DockerImage:                      langsmith.F("docker_image"),
			FsCapacityBytes:                  langsmith.F(int64(0)),
			ImageConfig: langsmith.F(langsmith.SandboxBoxNewSnapshotParamsImageConfig{
				Cmd:          langsmith.F([]string{"string"}),
				Entrypoint:   langsmith.F([]string{"string"}),
				ExposedPorts: langsmith.F([]int64{int64(0)}),
			}),
			Include:       langsmith.F([]langsmith.SandboxBoxNewSnapshotParamsInclude{langsmith.SandboxBoxNewSnapshotParamsIncludeAll}),
			IncludeMemory: langsmith.F(true),
			Labels: langsmith.F(map[string]string{
				"foo": "string",
			}),
			RunConfig: langsmith.F(langsmith.SandboxBoxNewSnapshotParamsRunConfig{
				EnvVars: langsmith.F(map[string]string{
					"foo": "string",
				}),
				User:    langsmith.F("user"),
				WorkDir: langsmith.F("work_dir"),
			}),
			Tag: langsmith.F("tag"),
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

func TestSandboxBoxDeleteServiceURLWithOptionalParams(t *testing.T) {
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
	err := client.Sandboxes.Boxes.DeleteServiceURL(
		context.TODO(),
		"name",
		langsmith.SandboxBoxDeleteServiceURLParams{
			Port: langsmith.F(int64(0)),
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

func TestSandboxBoxGenerateDownloadURLWithOptionalParams(t *testing.T) {
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
	_, err := client.Sandboxes.Boxes.GenerateDownloadURL(
		context.TODO(),
		"name",
		langsmith.SandboxBoxGenerateDownloadURLParams{
			Path:               langsmith.F("path"),
			ContentDisposition: langsmith.F("content_disposition"),
			ContentType:        langsmith.F("content_type"),
			CspSandbox:         langsmith.F(true),
			CspSandboxFlags:    langsmith.F([]langsmith.SandboxBoxGenerateDownloadURLParamsCspSandboxFlag{langsmith.SandboxBoxGenerateDownloadURLParamsCspSandboxFlagAllowDownloads}),
			CspSourceBundles:   langsmith.F([]langsmith.SandboxBoxGenerateDownloadURLParamsCspSourceBundle{langsmith.SandboxBoxGenerateDownloadURLParamsCspSourceBundleCdnjs}),
			ExpiresInSeconds:   langsmith.F(int64(0)),
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
			Access:           langsmith.F(langsmith.SandboxBoxGenerateServiceURLParamsAccessRestricted),
			EmbedOrigins:     langsmith.F([]string{"string"}),
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

func TestSandboxBoxListServiceURLsWithOptionalParams(t *testing.T) {
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
	_, err := client.Sandboxes.Boxes.ListServiceURLs(
		context.TODO(),
		"name",
		langsmith.SandboxBoxListServiceURLsParams{
			Cursor:   langsmith.F("cursor"),
			PageSize: langsmith.F(int64(0)),
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
