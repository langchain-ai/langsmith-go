// File generated from our OpenAPI spec by Stainless. See CONTRIBUTING.md for details.

package cmd

import (
	"context"
	"fmt"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/generated/internal/apiquery"
	"github.com/langchain-ai/langsmith-go/cli/internal/generated/internal/requestflag"
	"github.com/langchain-ai/langsmith-go/option"
	"github.com/tidwall/gjson"
	"github.com/urfave/cli/v3"
)

var sandboxesVolumesCreate = cli.Command{
	Name:    "create",
	Usage:   "**Beta:** This API is in active development and may change without notice.\nRequires sandbox volumes to be enabled for your workspace. Create an empty\nvolume that can be attached to sandboxes at a mount path. The size is rounded up\nto a whole GiB.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Name of the volume, unique within the workspace.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[int64]{
			Name:     "size-bytes",
			Usage:    "Capacity in bytes, rounded up to a whole GiB.",
			Required: true,
			BodyPath: "size_bytes",
		},
		&requestflag.Flag[string]{
			Name:     "description",
			Usage:    "Free-text description, at most 1024 characters.",
			BodyPath: "description",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "labels",
			Usage:    "Labels to filter and identify the volume by. Keys are lowercase letters, digits, dots and hyphens.",
			BodyPath: "labels",
		},
	},
	Action:          handleSandboxesVolumesCreate,
	HideHelpCommand: true,
}

var sandboxesVolumesRetrieve = cli.Command{
	Name:    "retrieve",
	Usage:   "**Beta:** This API is in active development and may change without notice.\nRequires sandbox volumes to be enabled for your workspace. Get a sandbox volume\nby ID or name, with its current attachments.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "volume",
			Required:  true,
			PathParam: "volume",
		},
	},
	Action:          handleSandboxesVolumesRetrieve,
	HideHelpCommand: true,
}

var sandboxesVolumesUpdate = cli.Command{
	Name:    "update",
	Usage:   "**Beta:** This API is in active development and may change without notice.\nRequires sandbox volumes to be enabled for your workspace. Update a volume's\ndescription or labels. Each field sent replaces the stored value; omitted fields\nare unchanged. Null, an empty description, or an empty labels object clears it.\nOnly the volume's creator, or a role with sandboxes:exec, may update it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "volume",
			Required:  true,
			PathParam: "volume",
		},
		&requestflag.Flag[*string]{
			Name:     "description",
			Usage:    "Description replaces the description; null or an empty string clears it.",
			BodyPath: "description",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "labels",
			Usage:    "Labels replaces every label; null or an empty object clears them.",
			BodyPath: "labels",
		},
	},
	Action:          handleSandboxesVolumesUpdate,
	HideHelpCommand: true,
}

var sandboxesVolumesList = cli.Command{
	Name:    "list",
	Usage:   "**Beta:** This API is in active development and may change without notice.\nRequires sandbox volumes to be enabled for your workspace. List the workspace's\nsandbox volumes with their current attachments.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "cursor",
			Usage:     "Opaque pagination cursor from a prior response's next_cursor",
			QueryPath: "cursor",
		},
		&requestflag.Flag[string]{
			Name:      "name-contains",
			Usage:     "Filter by name substring",
			QueryPath: "name_contains",
		},
		&requestflag.Flag[int64]{
			Name:      "page-size",
			Usage:     "Number of results per page",
			Default:   20,
			QueryPath: "page_size",
		},
		&requestflag.Flag[int64]{
			Name:  "max-items",
			Usage: "The maximum number of items to return (use -1 for unlimited).",
		},
	},
	Action:          handleSandboxesVolumesList,
	HideHelpCommand: true,
}

var sandboxesVolumesDelete = cli.Command{
	Name:    "delete",
	Usage:   "**Beta:** This API is in active development and may change without notice.\nRequires sandbox volumes to be enabled for your workspace. Delete a volume and\nits data. A volume attached to any sandbox cannot be deleted. Storage is\nreclaimed asynchronously. Only the volume's creator, or a role with\nsandboxes:exec, may delete it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "volume",
			Required:  true,
			PathParam: "volume",
		},
	},
	Action:          handleSandboxesVolumesDelete,
	HideHelpCommand: true,
}

var sandboxesVolumesClone = cli.Command{
	Name:    "clone",
	Usage:   "**Beta:** This API is in active development and may change without notice.\nRequires sandbox volumes to be enabled for your workspace. Create a new volume\nholding a copy-on-write copy of an existing one. A volume attached read-write to\na sandbox that is not stopped cannot be cloned. Only the volume's creator, or a\nrole with sandboxes:exec, may clone it.",
	Suggest: true,
	Flags: []cli.Flag{
		&requestflag.Flag[string]{
			Name:      "volume",
			Required:  true,
			PathParam: "volume",
		},
		&requestflag.Flag[string]{
			Name:     "name",
			Usage:    "Name of the new volume, unique within the workspace.",
			Required: true,
			BodyPath: "name",
		},
		&requestflag.Flag[string]{
			Name:     "description",
			Usage:    "Free-text description of the new volume, at most 1024 characters.",
			BodyPath: "description",
		},
		&requestflag.Flag[map[string]any]{
			Name:     "labels",
			Usage:    "Labels for the new volume. Keys are lowercase letters, digits, dots and hyphens.",
			BodyPath: "labels",
		},
	},
	Action:          handleSandboxesVolumesClone,
	HideHelpCommand: true,
}

func handleSandboxesVolumesCreate(ctx context.Context, cmd *cli.Command) error {
	client := langsmith.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatRepeat,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := langsmith.SandboxVolumeNewParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Sandboxes.Volumes.New(ctx, params, options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "sandboxes:volumes create",
		Transform:      transform,
	})
}

func handleSandboxesVolumesRetrieve(ctx context.Context, cmd *cli.Command) error {
	client := langsmith.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("volume") && len(unusedArgs) > 0 {
		cmd.Set("volume", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatRepeat,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Sandboxes.Volumes.Get(ctx, cmd.Value("volume").(string), options...)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "sandboxes:volumes retrieve",
		Transform:      transform,
	})
}

func handleSandboxesVolumesUpdate(ctx context.Context, cmd *cli.Command) error {
	client := langsmith.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("volume") && len(unusedArgs) > 0 {
		cmd.Set("volume", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatRepeat,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := langsmith.SandboxVolumeUpdateParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Sandboxes.Volumes.Update(
		ctx,
		cmd.Value("volume").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "sandboxes:volumes update",
		Transform:      transform,
	})
}

func handleSandboxesVolumesList(ctx context.Context, cmd *cli.Command) error {
	client := langsmith.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()

	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatRepeat,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	params := langsmith.SandboxVolumeListParams{}

	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	if format == "raw" {
		var res []byte
		options = append(options, option.WithResponseBodyInto(&res))
		_, err = client.Sandboxes.Volumes.List(ctx, params, options...)
		if err != nil {
			return err
		}
		obj := gjson.ParseBytes(res)
		return ShowJSON(obj, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "sandboxes:volumes list",
			Transform:      transform,
		})
	} else {
		iter := client.Sandboxes.Volumes.ListAutoPaging(ctx, params, options...)
		maxItems := int64(-1)
		if cmd.IsSet("max-items") {
			maxItems = cmd.Value("max-items").(int64)
		}
		return ShowJSONIterator(iter, maxItems, ShowJSONOpts{
			ExplicitFormat: explicitFormat,
			Format:         format,
			RawOutput:      cmd.Root().Bool("raw-output"),
			Title:          "sandboxes:volumes list",
			Transform:      transform,
		})
	}
}

func handleSandboxesVolumesDelete(ctx context.Context, cmd *cli.Command) error {
	client := langsmith.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("volume") && len(unusedArgs) > 0 {
		cmd.Set("volume", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatRepeat,
		EmptyBody,
		false,
	)
	if err != nil {
		return err
	}

	return client.Sandboxes.Volumes.Delete(ctx, cmd.Value("volume").(string), options...)
}

func handleSandboxesVolumesClone(ctx context.Context, cmd *cli.Command) error {
	client := langsmith.NewClient(getDefaultRequestOptions(cmd)...)
	unusedArgs := cmd.Args().Slice()
	if !cmd.IsSet("volume") && len(unusedArgs) > 0 {
		cmd.Set("volume", unusedArgs[0])
		unusedArgs = unusedArgs[1:]
	}
	if len(unusedArgs) > 0 {
		return fmt.Errorf("Unexpected extra arguments: %v", unusedArgs)
	}

	options, err := flagOptions(
		cmd,
		apiquery.NestedQueryFormatBrackets,
		apiquery.ArrayQueryFormatRepeat,
		ApplicationJSON,
		false,
	)
	if err != nil {
		return err
	}

	params := langsmith.SandboxVolumeCloneParams{}

	var res []byte
	options = append(options, option.WithResponseBodyInto(&res))
	_, err = client.Sandboxes.Volumes.Clone(
		ctx,
		cmd.Value("volume").(string),
		params,
		options...,
	)
	if err != nil {
		return err
	}

	obj := gjson.ParseBytes(res)
	format := cmd.Root().String("format")
	explicitFormat := cmd.Root().IsSet("format")
	transform := cmd.Root().String("transform")
	return ShowJSON(obj, ShowJSONOpts{
		ExplicitFormat: explicitFormat,
		Format:         format,
		RawOutput:      cmd.Root().Bool("raw-output"),
		Title:          "sandboxes:volumes clone",
		Transform:      transform,
	})
}
