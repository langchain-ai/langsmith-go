//go:build integration

package integration

import (
	"context"
	"strings"
	"testing"

	"github.com/langchain-ai/langsmith-go"
)

// repoOwner extracts the owner from a repo's FullName ("owner/handle").
// Falls back to the Owner field, or "-" if both are empty.
func repoOwner(fullName, owner string) string {
	if parts := strings.SplitN(fullName, "/", 2); len(parts) == 2 && parts[0] != "" {
		return parts[0]
	}
	if owner != "" {
		return owner
	}
	return "-"
}

// --- Repo (Prompt) CRUD ---

func TestRepoCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	handle := uniqueName("go-integ-prompt")

	// Create repo
	created, err := client.Repos.New(ctx, langsmith.RepoNewParams{
		RepoHandle:  langsmith.F(handle),
		IsPublic:    langsmith.F(false),
		Description: langsmith.F("Integration test prompt"),
		Tags:        langsmith.F([]string{"test", "go-sdk"}),
	})
	if err != nil {
		t.Fatalf("create repo: %v", err)
	}
	if created.Repo.RepoHandle != handle {
		t.Errorf("handle = %q, want %q", created.Repo.RepoHandle, handle)
	}
	owner := repoOwner(created.Repo.FullName, created.Repo.Owner)

	// Get repo
	got, err := client.Repos.Get(ctx, owner, handle)
	if err != nil {
		t.Fatalf("get repo: %v", err)
	}
	if got.Repo.RepoHandle != handle {
		t.Errorf("get returned different handle")
	}

	// Update repo
	_, err = client.Repos.Update(ctx, owner, handle, langsmith.RepoUpdateParams{
		Description: langsmith.F("Updated description"),
	})
	if err != nil {
		t.Fatalf("update repo: %v", err)
	}

	// Verify update persisted
	updated, err := client.Repos.Get(ctx, owner, handle)
	if err != nil {
		t.Fatalf("get updated repo: %v", err)
	}
	if updated.Repo.Description != "Updated description" {
		t.Errorf("description = %q, want 'Updated description'", updated.Repo.Description)
	}

	// List repos
	listed, err := client.Repos.List(ctx, langsmith.RepoListParams{
		Query: langsmith.F(handle),
	})
	if err != nil {
		t.Fatalf("list repos: %v", err)
	}
	if len(listed.Repos) == 0 {
		t.Error("expected at least one repo in list")
	}

	// Delete repo
	_, err = client.Repos.Delete(ctx, owner, handle)
	if err != nil {
		t.Fatalf("delete repo: %v", err)
	}
}

// --- Commit CRUD ---

func TestCommitCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	handle := uniqueName("go-integ-commit")

	// Create repo first
	created, err := client.Repos.New(ctx, langsmith.RepoNewParams{
		RepoHandle: langsmith.F(handle),
		IsPublic:   langsmith.F(false),
	})
	if err != nil {
		t.Fatalf("create repo: %v", err)
	}
	owner := repoOwner(created.Repo.FullName, created.Repo.Owner)
	defer client.Repos.Delete(ctx, owner, handle)

	// Create a commit with a simple prompt manifest
	manifest := map[string]interface{}{
		"lc":   1,
		"type": "constructor",
		"id":   []string{"langchain", "prompts", "prompt", "PromptTemplate"},
		"kwargs": map[string]interface{}{
			"input_variables": []string{"topic"},
			"template":        "Tell me a joke about {topic}",
		},
	}
	commit, err := client.Commits.New(ctx, owner, handle, langsmith.CommitNewParams{
		Manifest: langsmith.F[interface{}](manifest),
	})
	if err != nil {
		t.Fatalf("create commit: %v", err)
	}
	if commit == nil {
		t.Fatal("expected non-nil commit response")
	}

	// List commits
	listed, err := client.Commits.List(ctx, owner, handle, langsmith.CommitListParams{})
	if err != nil {
		t.Fatalf("list commits: %v", err)
	}
	if len(listed.Commits) == 0 {
		t.Error("expected at least one commit")
	}
}

// --- Settings ---

func TestGetSettings(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()

	settings, err := client.Settings.List(ctx)
	if err != nil {
		t.Fatalf("get settings: %v", err)
	}
	if settings.ID == "" {
		t.Error("expected non-empty tenant ID in settings")
	}
}

// --- Annotation Queue CRUD ---

func TestAnnotationQueueCRUD(t *testing.T) {
	client := newClient(t)
	ctx := context.Background()
	name := uniqueName("go-integ-queue")

	queue, err := client.AnnotationQueues.AnnotationQueues(ctx, langsmith.AnnotationQueueAnnotationQueuesParams{
		Name: langsmith.F(name),
	})
	if err != nil {
		t.Fatalf("create annotation queue: %v", err)
	}
	if queue.ID == "" {
		t.Fatal("expected non-empty queue ID")
	}

	// Get queue
	got, err := client.AnnotationQueues.Get(ctx, queue.ID)
	if err != nil {
		t.Fatalf("get annotation queue: %v", err)
	}
	if got.Name != name {
		t.Errorf("name = %q, want %q", got.Name, name)
	}

	// List queues
	listed, err := client.AnnotationQueues.GetAnnotationQueues(ctx, langsmith.AnnotationQueueGetAnnotationQueuesParams{})
	if err != nil {
		t.Fatalf("list annotation queues: %v", err)
	}
	if len(listed.Items) == 0 {
		t.Error("expected at least one annotation queue")
	}

	// Get size
	size, err := client.AnnotationQueues.GetSize(ctx, queue.ID, langsmith.AnnotationQueueGetSizeParams{})
	if err != nil {
		t.Fatalf("get queue size: %v", err)
	}
	if size.Size != 0 {
		t.Errorf("expected empty queue, got size %d", size.Size)
	}

	// Delete queue
	_, err = client.AnnotationQueues.Delete(ctx, queue.ID)
	if err != nil {
		t.Fatalf("delete annotation queue: %v", err)
	}
}
