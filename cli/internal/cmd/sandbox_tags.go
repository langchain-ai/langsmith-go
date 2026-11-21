package cmd

import (
	"context"
	"fmt"
	"slices"
	"strings"

	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/spf13/cobra"
)

// addResourceTagFlag registers --resource-tag, the by-name counterpart of
// --tag-value-id.
func addResourceTagFlag(cmd *cobra.Command, specs *[]string, usage string) {
	cmd.Flags().StringArrayVar(specs, "resource-tag", nil, usage)
}

// withResourceTags returns ids followed by the tag value IDs that specs name,
// without duplicates. specs are KEY=VALUE pairs resolved against the
// workspace's resource tags; a tag that does not exist is an error, never
// created.
func withResourceTags(ctx context.Context, c *client.Client, ids, specs []string) ([]string, error) {
	if len(specs) == 0 {
		return ids, nil
	}
	tags, err := c.SDK.WorkspaceTags.List(ctx, langsmith.WorkspaceTagListParams{})
	if err != nil {
		return nil, fmt.Errorf("listing resource tags: %w", err)
	}
	var keys []langsmith.TagKeyWithValues
	if tags != nil {
		keys = *tags
	}

	out := slices.Clone(ids)
	for _, spec := range specs {
		key, value, ok := strings.Cut(spec, "=")
		if !ok || key == "" || value == "" {
			return nil, fmt.Errorf("--resource-tag must be KEY=VALUE (got %q)", spec)
		}
		id, err := resourceTagValueID(keys, key, value)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(out, id) {
			out = append(out, id)
		}
	}
	return out, nil
}

func resourceTagValueID(keys []langsmith.TagKeyWithValues, key, value string) (string, error) {
	names := make([]string, 0, len(keys))
	for _, k := range keys {
		names = append(names, k.Key)
		if k.Key != key {
			continue
		}
		values := make([]string, 0, len(k.Values))
		for _, v := range k.Values {
			if v.Value == value {
				return v.ID, nil
			}
			values = append(values, v.Value)
		}
		return "", fmt.Errorf("resource tag %s has no value %q (values: %s)", key, value, listOrNone(values))
	}
	return "", fmt.Errorf("no resource tag %q in this workspace (keys: %s)", key, listOrNone(names))
}

func listOrNone(items []string) string {
	if len(items) == 0 {
		return "none"
	}
	slices.Sort(items)
	return strings.Join(items, ", ")
}
