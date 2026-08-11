package main

import (
	"context"
	"fmt"
	"os"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates mounting a Context Hub repo into a sandbox.
//
// The mount is read-only: the repo's latest commit tree is mirrored into the
// mount path and kept in sync for the sandbox's lifetime, but guest writes are
// never pushed back to the repo. Unlike bucket and git mounts, a Context Hub
// mount can target any guest path outside the system directories.
//
// Prerequisites:
//   - LANGSMITH_API_KEY: an API key with access to the Context Hub repo
//
// Running:
//
//	go run ./examples/sandbox_context_hub_mount -- <owner>/<repo>
const mountPath = "/memories"

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	repo := "-/my-agent"
	if len(os.Args) > 1 {
		repo = os.Args[1]
	}

	mount := langsmith.ContextHubMount(langsmith.ContextHubMountParams{
		ID:        "memories",
		MountPath: mountPath,
		Repo:      repo,
	})

	client := langsmith.NewClient()
	ctx := context.Background()

	sandbox, err := client.Sandboxes.Boxes.NewSandbox(ctx, langsmith.SandboxBoxNewParams{
		Name: langsmith.String("context-hub-mount-sandbox"),
		MountConfig: langsmith.F(langsmith.SandboxBoxNewParamsMountConfig{
			Mounts: langsmith.F([]langsmith.SandboxBoxNewParamsMountConfigMountUnion{mount}),
		}),
	})
	if err != nil {
		return err
	}
	defer func() {
		if err := sandbox.Delete(ctx); err != nil {
			fmt.Fprintf(os.Stderr, "Failed to delete sandbox: %v\n", err)
		}
	}()

	result, err := sandbox.Run(ctx, langsmith.SandboxBoxRunParams{
		Command: langsmith.String("ls -R " + mountPath),
	})
	if err != nil {
		return err
	}
	fmt.Printf("Context Hub repo %q mounted at %s:\n%s", repo, mountPath, result.Stdout)
	return nil
}
