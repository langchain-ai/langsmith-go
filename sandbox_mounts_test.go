package langsmith

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestContextHubMountSerializesBackendShape(t *testing.T) {
	raw, err := json.Marshal(ContextHubMount(ContextHubMountParams{
		ID:              "memories",
		MountPath:       "/memories",
		Repo:            "-/my-agent",
		InitialPullOnly: true,
	}))
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, map[string]any{
		"id":         "memories",
		"type":       "contexthub",
		"mount_path": "/memories",
		"contexthub": map[string]any{
			"repo":              "-/my-agent",
			"initial_pull_only": true,
		},
	}, got)
}

func TestContextHubMountOmitsOptionalFields(t *testing.T) {
	raw, err := json.Marshal(ContextHubMount(ContextHubMountParams{
		ID:        "memories",
		MountPath: "/memories",
		Repo:      "-/my-agent",
	}))
	require.NoError(t, err)

	var got map[string]any
	require.NoError(t, json.Unmarshal(raw, &got))
	assert.Equal(t, map[string]any{
		"id":         "memories",
		"type":       "contexthub",
		"mount_path": "/memories",
		"contexthub": map[string]any{"repo": "-/my-agent"},
	}, got)
}

func TestContextHubMountUsableAsMountUnion(t *testing.T) {
	params := SandboxBoxNewParams{
		MountConfig: F(SandboxBoxNewParamsMountConfig{
			Mounts: F([]SandboxBoxNewParamsMountConfigMountUnion{
				ContextHubMount(ContextHubMountParams{
					ID:        "memories",
					MountPath: "/memories",
					Repo:      "-/my-agent",
				}),
			}),
		}),
	}

	raw, err := json.Marshal(params)
	require.NoError(t, err)
	assert.JSONEq(t, `{"mount_config":{"mounts":[{"id":"memories","type":"contexthub","mount_path":"/memories","contexthub":{"repo":"-/my-agent"}}]}}`, string(raw))
}
