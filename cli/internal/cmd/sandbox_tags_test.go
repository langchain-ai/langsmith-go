package cmd

import (
	"encoding/json"
	"fmt"
	"net/http"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const workspaceTagsJSON = `[
	{"id":"key-env","key":"env","values":[{"id":"tv-prod","value":"prod"},{"id":"tv-dev","value":"dev"}]},
	{"id":"key-team","key":"team","values":[{"id":"tv-infra","value":"infra"}]}
]`

// sandboxTagServer serves the workspace's resource tags and records the
// request body or query sent to the sandbox endpoint.
func sandboxTagServer(t *testing.T, sandboxResponse string) (url string, body *map[string]any, query *map[string][]string, sandboxCalls *int) {
	t.Helper()
	body, query, sandboxCalls = &map[string]any{}, &map[string][]string{}, new(int)
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		if r.URL.Path == "/api/v1/workspaces/current/tags" {
			_, _ = w.Write([]byte(workspaceTagsJSON))
			return
		}
		*sandboxCalls++
		*query = r.URL.Query()
		if r.Body != nil && r.Method != http.MethodGet {
			require.NoError(t, json.NewDecoder(r.Body).Decode(body))
		}
		_, _ = w.Write([]byte(sandboxResponse))
	})
	return ts.URL, body, query, sandboxCalls
}

func TestSandboxCreateCmd_ResolvesResourceTagsByName(t *testing.T) {
	url, body, _, _ := sandboxTagServer(t, `{"id":"box-id","name":"my-vm","status":"provisioning"}`)

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "create", "my-vm",
		"--tag-value-id", "tv-other",
		"--resource-tag", "env=prod",
		"--resource-tag", "team=infra",
		"--resource-tag", "env=prod",
		"--format", "json")
	require.NoError(t, err)

	assert.Equal(t, []any{"tv-other", "tv-prod", "tv-infra"}, (*body)["tag_value_ids"])
}

func TestSandboxUpdateCmd_ResourceTagReplacesTags(t *testing.T) {
	url, body, _, _ := sandboxTagServer(t, `{"id":"box-id","name":"my-vm","status":"ready"}`)

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "update", "my-vm",
		"--resource-tag", "env=dev", "--format", "json")
	require.NoError(t, err)

	assert.Equal(t, map[string]any{"tag_value_ids": []any{"tv-dev"}}, *body)
}

func TestSandboxListCmd_FiltersByResourceTag(t *testing.T) {
	url, _, query, _ := sandboxTagServer(t, `{"items":[],"next_cursor":null}`)

	_, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "list",
		"--resource-tag", "team=infra", "--format", "json")
	require.NoError(t, err)

	assert.Equal(t, []string{"tv-infra"}, (*query)["tag_value_id"])
}

func TestSandboxCreateCmd_RejectsUnknownResourceTags(t *testing.T) {
	tests := []struct {
		spec string
		want string
	}{
		{"region=us", `no resource tag "region" in this workspace (keys: env, team)`},
		{"env=staging", `resource tag env has no value "staging" (values: dev, prod)`},
		{"env", "--resource-tag must be KEY=VALUE"},
	}
	for _, tt := range tests {
		t.Run(tt.spec, func(t *testing.T) {
			url, _, _, sandboxCalls := sandboxTagServer(t, `{}`)

			_, err := executeCommand(t, "--api-key", "test-key", "--api-url", url, "sandbox", "create", "my-vm",
				"--resource-tag", tt.spec)
			require.ErrorContains(t, err, tt.want)
			assert.Zero(t, *sandboxCalls, "no sandbox is created when a tag does not resolve")
		})
	}
}

func TestSandboxDeleteCmd_SeveralNamesBatchAndReportSkipped(t *testing.T) {
	var batches [][]string
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodPost, r.Method)
		require.Equal(t, "/api/v2/sandboxes/boxes/batch-delete", r.URL.Path)
		var req struct {
			Names []string `json:"names"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		batches = append(batches, req.Names)
		deleted := req.Names
		skipped := `[]`
		if req.Names[0] == "vm-0" {
			deleted, skipped = req.Names[1:], `[{"name":"vm-0","reason":"sandbox not found"}]`
		}
		raw, _ := json.Marshal(deleted)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"deleted":%s,"skipped":%s}`, raw, skipped)
	})

	args := []string{"--api-key", "test-key", "--api-url", ts.URL, "sandbox", "delete"}
	for i := range 51 {
		args = append(args, fmt.Sprintf("vm-%d", i))
	}
	out, err := executeCommand(t, args...)

	require.ErrorContains(t, err, "1 of 51 sandboxes not deleted")
	require.Len(t, batches, 2)
	assert.Len(t, batches[0], 50)
	assert.Equal(t, []string{"vm-50"}, batches[1])
	assert.Contains(t, out, "Deleted vm-1\n")
	assert.Contains(t, out, "Deleted vm-50\n")
	assert.Contains(t, out, "Skipped vm-0: sandbox not found\n")
}

func TestSandboxDeleteCmd_OneNameDeletesDirectly(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, http.MethodDelete, r.Method)
		require.Equal(t, "/api/v2/sandboxes/boxes/my-vm", r.URL.Path)
		w.WriteHeader(http.StatusNoContent)
	})

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "delete", "my-vm")
	require.NoError(t, err)
	assert.Contains(t, out, "Sandbox deleted.\n")
}

func TestSandboxUsageCmd(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v2/sandboxes/usage", r.URL.Path)
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"sandboxes_used":3,"sandboxes_limit":10,"sandbox_stopped_count":2,
			"cpu_used_millicores":3500,"cpu_limit_millicores":16000,
			"memory_used_bytes":8589934592,"memory_limit_bytes":68719476736,"snapshots_count":4,"registries_count":1}`))
	})

	out, err := executeCommand(t, "--api-key", "test-key", "--api-url", ts.URL, "sandbox", "usage")
	require.NoError(t, err)
	assert.Regexp(t, `Running:\s+3 / 10\n`, out)
	assert.Regexp(t, `Stopped:\s+2\n`, out)
	assert.Regexp(t, `CPU:\s+3500m / 16000m\n`, out)
}

func TestSandboxDeleteCmd_ReportsEarlierBatchesWhenALaterOneFails(t *testing.T) {
	ts := newTestServer(t, func(w http.ResponseWriter, r *http.Request) {
		var req struct {
			Names []string `json:"names"`
		}
		require.NoError(t, json.NewDecoder(r.Body).Decode(&req))
		if req.Names[0] != "vm-0" {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(`{"detail":"invalid"}`))
			return
		}
		raw, _ := json.Marshal(req.Names)
		w.Header().Set("Content-Type", "application/json")
		_, _ = fmt.Fprintf(w, `{"deleted":%s,"skipped":[]}`, raw)
	})

	args := []string{"--api-key", "test-key", "--api-url", ts.URL, "sandbox", "delete"}
	for i := range 52 {
		args = append(args, fmt.Sprintf("vm-%d", i))
	}
	out, err := executeCommand(t, args...)

	require.ErrorContains(t, err, "2 of 52 without a confirmed result")
	assert.Contains(t, out, "Deleted vm-0\n")
	assert.Contains(t, out, "Deleted vm-49\n")
	assert.NotContains(t, out, "vm-50")
}
