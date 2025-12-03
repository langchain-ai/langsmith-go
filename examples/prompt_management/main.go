package main

import (
	"context"
	"errors"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/langchain-ai/langsmith-go"
)

// Demonstrates how to manage prompts programmatically using the LangSmith Go SDK.
//
// This example shows:
//   - Creating prompt repositories (repos)
//   - Adding prompt content as commits (with variables)
//   - Listing prompts with filters
//   - Retrieving prompt content
//   - Pulling specific versions
//
// Prerequisites:
//   - LANGSMITH_API_KEY: Your LangSmith API key
//   - LANGSMITH_OWNER: Your LangSmith owner (defaults to "-" if not set)
//   - LANGCHAIN_BASE_URL or LANGSMITH_ENDPOINT: LangSmith API URL (https://api.smith.langchain.com)
//
// Running:
//
//	go run ./examples/prompt_management

const (
	defaultOwner      = "-"
	promptName        = "product-description-generator"
	fullNameSeparator = "/"
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	fmt.Println("=== LangSmith Prompt Management Example ===")

	client := langsmith.NewClient()
	ctx := context.Background()

	// 1. Check if prompt already exists
	fmt.Println("1. Checking for existing prompt using client.Repos.List()...")
	repo, promptExists, err := findPrompt(ctx, client, promptName)
	if err != nil {
		return fmt.Errorf("finding prompt: %w", err)
	}

	var owner string
	if promptExists {
		owner = extractOwner(repo)
		fmt.Printf("2. Prompt '%s' already exists\n", promptName)
		fmt.Printf("   View: https://smith.langchain.com/prompts/%s\n", repo.FullName)
		fmt.Println("   Skipping creation steps...")
	} else {
		fmt.Println("   ✓ No existing prompt found")

		// 2. Create prompt repository (metadata only)
		fmt.Println("2. Creating prompt repository using client.Repos.New()...")
		createResp, err := client.Repos.New(ctx, langsmith.RepoNewParams{
			RepoHandle:  langsmith.F(promptName),
			Description: langsmith.F("A product description generator that creates compelling descriptions for e-commerce products"),
			IsPublic:    langsmith.F(false),
		})
		if err != nil {
			return fmt.Errorf("creating repo: %w", err)
		}
		repo = &createResp.Repo
		owner = extractOwner(repo)
		fmt.Printf("   ✓ Created prompt repository: %s\n\n", promptName)
	}

	// 3. Add prompt content as a commit (only if prompt is new or has no commits)
	commitCreated, err := ensureCommit(ctx, client, promptName, owner, !promptExists)
	if err != nil {
		return fmt.Errorf("ensuring commit: %w", err)
	}

	// 4. List prompts
	if err := listPrompts(ctx, client); err != nil {
		return fmt.Errorf("listing prompts: %w", err)
	}

	// 5. Pull prompt content (retrieve latest commit)
	fmt.Println("5. Pulling prompt content using client.Get() with commit='latest'...")
	manifestResponse, err := retrieveLatestCommit(ctx, client, owner, promptName)
	if err != nil {
		return fmt.Errorf("retrieving commit: %w", err)
	}

	fmt.Println("   ✓ Retrieved latest commit manifest")
	promptContent := extractPromptContent(manifestResponse.Manifest)
	fmt.Printf("   ✓ Prompt content: \"%s\"\n\n", promptContent)

	// 6. Demonstrate using the prompt with a value
	demonstratePromptUsage(promptContent)

	// Summary
	printSummary(!promptExists, commitCreated)
	return nil
}

// extractOwner extracts the owner from a repo, handling nullable owner fields.
func extractOwner(repo *langsmith.RepoWithLookups) string {
	if repo.Owner != "" {
		return repo.Owner
	}
	// Parse owner from FullName (format: "owner/repo")
	if idx := strings.Index(repo.FullName, fullNameSeparator); idx >= 0 {
		return repo.FullName[:idx]
	}
	return defaultOwner
}

// ensureCommit ensures a commit exists for the prompt.
// Returns true if a commit was created, false if it already existed.
func ensureCommit(ctx context.Context, client *langsmith.Client, promptName, owner string, isNewPrompt bool) (bool, error) {
	if !isNewPrompt {
		hasCommit, err := hasLatestCommit(ctx, client, promptName, owner)
		if err != nil {
			return false, fmt.Errorf("checking for latest commit: %w", err)
		}
		if hasCommit {
			fmt.Println("3. Prompt already has content, skipping commit creation...")
			fmt.Println("   To update the prompt, delete it first or modify the content")
			return false, nil
		}
	}

	if err := createCommit(ctx, client, promptName, owner, ""); err != nil {
		return false, err
	}
	return true, nil
}

// listPrompts lists and filters prompts.
func listPrompts(ctx context.Context, client *langsmith.Client) error {
	fmt.Println("4. Listing prompts using client.Repos.List()...")

	allPrompts, err := client.Repos.List(ctx, langsmith.RepoListParams{
		IsPublic: langsmith.F(langsmith.RepoListParamsIsPublicFalse),
	})
	if err != nil {
		return fmt.Errorf("listing all prompts: %w", err)
	}
	fmt.Printf("   ✓ Found %d prompt(s) in your organization\n", len(allPrompts.Repos))

	productPrompts, err := client.Repos.List(ctx, langsmith.RepoListParams{
		Query:    langsmith.F("product"),
		IsPublic: langsmith.F(langsmith.RepoListParamsIsPublicFalse),
	})
	if err != nil {
		return fmt.Errorf("listing product prompts: %w", err)
	}
	fmt.Printf("   ✓ Found %d prompt(s) matching 'product' in your organization\n\n", len(productPrompts.Repos))
	return nil
}

// retrieveLatestCommit retrieves the latest commit manifest for a prompt.
func retrieveLatestCommit(ctx context.Context, client *langsmith.Client, owner, promptName string) (*langsmith.CommitManifestResponse, error) {
	var manifestResponse langsmith.CommitManifestResponse
	path := fmt.Sprintf("api/v1/commits/%s/%s/latest", owner, promptName)
	if err := client.Get(ctx, path, nil, &manifestResponse); err != nil {
		return nil, fmt.Errorf("retrieving latest commit: %w", err)
	}
	return &manifestResponse, nil
}

// demonstratePromptUsage demonstrates using a prompt template with a value.
func demonstratePromptUsage(promptContent string) {
	fmt.Println("6. Using the prompt with a value...")
	filledPrompt := strings.ReplaceAll(promptContent, "{product_name}", "Wireless Bluetooth Headphones")
	fmt.Printf("   Template: \"%s\"\n", promptContent)
	fmt.Printf("   Filled:   \"%s\"\n\n", filledPrompt)
}

// printSummary prints a summary of the operations performed.
func printSummary(repoCreated, commitCreated bool) {
	fmt.Println("=== Summary ===")
	if repoCreated {
		fmt.Println("✓ Created prompt repository")
	} else {
		fmt.Println("ℹ Prompt already existed (skipped creation)")
	}
	if commitCreated {
		fmt.Println("✓ Added prompt content as commit")
	} else {
		fmt.Println("ℹ Prompt content unchanged (skipped commit)")
	}
	fmt.Println("✓ Listed and filtered prompts")
	fmt.Println("✓ Retrieved prompt content")
	fmt.Println("✓ Demonstrated using prompt with a value")
	fmt.Println()
	fmt.Println("Learn more:")
	fmt.Println("https://docs.langchain.com/langsmith/manage-prompts-programmatically")
}

// ChatPromptBuilder helps build chat prompt manifests compatible with LangSmith.
type ChatPromptBuilder struct {
	messages       []Message
	inputVariables []string
}

// Message represents a single message in a chat prompt.
type Message struct {
	Type     MessageType
	Template string
}

// MessageType represents the type of a message.
type MessageType int

const (
	MessageTypeSystem MessageType = iota
	MessageTypeHuman
	MessageTypeAI
)

// SystemMessage adds a system message to the prompt.
func (b *ChatPromptBuilder) SystemMessage(template string) *ChatPromptBuilder {
	b.messages = append(b.messages, Message{Type: MessageTypeSystem, Template: template})
	return b
}

// UserMessage adds a user/human message to the prompt.
func (b *ChatPromptBuilder) UserMessage(template string) *ChatPromptBuilder {
	b.messages = append(b.messages, Message{Type: MessageTypeHuman, Template: template})
	return b
}

// AIMessage adds an AI/assistant message to the prompt.
func (b *ChatPromptBuilder) AIMessage(template string) *ChatPromptBuilder {
	b.messages = append(b.messages, Message{Type: MessageTypeAI, Template: template})
	return b
}

// InputVariables explicitly adds input variables (optional - variables are auto-detected from templates).
func (b *ChatPromptBuilder) InputVariables(variables ...string) *ChatPromptBuilder {
	b.inputVariables = append(b.inputVariables, variables...)
	return b
}

// Build creates the prompt manifest compatible with LangSmith.
func (b *ChatPromptBuilder) Build() (map[string]interface{}, error) {
	if len(b.messages) == 0 {
		return nil, errors.New("at least one message is required")
	}

	manifestMessages := make([]interface{}, 0, len(b.messages))
	for _, msg := range b.messages {
		manifestMessages = append(manifestMessages, msg.toManifest())
	}

	finalInputVariables := b.inputVariables
	if len(finalInputVariables) == 0 {
		// Auto-detect from all messages if not explicitly set
		seen := make(map[string]bool, len(b.messages))
		for _, msg := range b.messages {
			vars := extractVariablesFromTemplate(msg.Template)
			for _, v := range vars {
				if !seen[v] {
					finalInputVariables = append(finalInputVariables, v)
					seen[v] = true
				}
			}
		}
	}

	return map[string]interface{}{
		"lc":   1,
		"type": "constructor",
		"id":   []string{"langchain_core", "prompts", "chat", "ChatPromptTemplate"},
		"kwargs": map[string]interface{}{
			"input_variables": finalInputVariables,
			"messages":        manifestMessages,
		},
	}, nil
}

// extractVariablesFromTemplate extracts variable names from a template string.
// Variables are identified by the pattern {variableName}.
func extractVariablesFromTemplate(template string) []string {
	re := regexp.MustCompile(`\{([^}]+)\}`)
	matches := re.FindAllStringSubmatch(template, -1)
	variables := make([]string, 0, len(matches))
	seen := make(map[string]bool)
	for _, match := range matches {
		if len(match) > 1 {
			v := match[1]
			if !seen[v] {
				variables = append(variables, v)
				seen[v] = true
			}
		}
	}
	return variables
}

// toManifest converts a Message to its manifest representation.
func (m Message) toManifest() map[string]interface{} {
	var className string
	switch m.Type {
	case MessageTypeSystem:
		className = "SystemMessagePromptTemplate"
	case MessageTypeAI:
		className = "AIMessagePromptTemplate"
	case MessageTypeHuman:
		className = "HumanMessagePromptTemplate"
	}

	promptTemplate := map[string]interface{}{
		"lc":   1,
		"type": "constructor",
		"id":   []string{"langchain_core", "prompts", "prompt", "PromptTemplate"},
		"kwargs": map[string]interface{}{
			"template": m.Template,
		},
	}

	return map[string]interface{}{
		"lc":   1,
		"type": "constructor",
		"id":   []string{"langchain_core", "prompts", "chat", className},
		"kwargs": map[string]interface{}{
			"prompt": promptTemplate,
		},
	}
}

// getMessageType determines message type from the message object's id field.
func getMessageType(message map[string]interface{}) string {
	idObj, ok := message["id"]
	if !ok {
		return "user"
	}

	var className string
	switch v := idObj.(type) {
	case []interface{}:
		if len(v) > 0 {
			if str, ok := v[len(v)-1].(string); ok {
				className = str
			}
		}
	case []string:
		if len(v) > 0 {
			className = v[len(v)-1]
		}
	}

	switch {
	case strings.Contains(className, "System"):
		return "system"
	case strings.Contains(className, "AI") || strings.Contains(className, "Assistant"):
		return "assistant"
	default:
		return "user"
	}
}

// findPrompt finds a prompt by exact repo handle match.
// Returns the repo if found, a boolean indicating if it was found, and any error.
func findPrompt(ctx context.Context, client *langsmith.Client, promptName string) (*langsmith.RepoWithLookups, bool, error) {
	repos, err := client.Repos.List(ctx, langsmith.RepoListParams{
		Query:    langsmith.F(promptName),
		IsPublic: langsmith.F(langsmith.RepoListParamsIsPublicFalse),
	})
	if err != nil {
		return nil, false, fmt.Errorf("listing repos: %w", err)
	}

	for _, repo := range repos.Repos {
		if repo.RepoHandle == promptName {
			return &repo, true, nil
		}
	}
	return nil, false, nil
}

// hasLatestCommit checks if a prompt repository has a latest commit.
func hasLatestCommit(ctx context.Context, client *langsmith.Client, promptName, owner string) (bool, error) {
	var manifestResponse langsmith.CommitManifestResponse
	path := fmt.Sprintf("api/v1/commits/%s/%s/latest", owner, promptName)
	err := client.Get(ctx, path, nil, &manifestResponse)
	if err != nil {
		var apiErr *langsmith.Error
		if errors.As(err, &apiErr) && apiErr.StatusCode == 404 {
			return false, nil
		}
		return false, fmt.Errorf("checking for latest commit: %w", err)
	}
	return true, nil
}

// createCommit creates a commit with prompt content.
func createCommit(ctx context.Context, client *langsmith.Client, promptName, owner, parentCommit string) error {
	fmt.Println("3. Adding prompt content using client.Post()...")

	builder := &ChatPromptBuilder{}
	builder.
		SystemMessage("You are a professional copywriter specializing in e-commerce product descriptions.").
		UserMessage("Write a compelling product description for {product_name}. Highlight key features and benefits.").
		InputVariables("product_name")

	manifest, err := builder.Build()
	if err != nil {
		return fmt.Errorf("building manifest: %w", err)
	}

	params := map[string]interface{}{
		"manifest": manifest,
	}
	if parentCommit != "" {
		params["parent_commit"] = parentCommit
	}

	path := fmt.Sprintf("api/v1/commits/%s/%s", owner, promptName)
	var createResp interface{}
	if err := client.Post(ctx, path, params, &createResp); err != nil {
		return fmt.Errorf("posting commit: %w", err)
	}

	fmt.Println("   ✓ Added prompt content as commit")
	return nil
}

// extractPromptContent extracts prompt content from a manifest for display (shows all messages).
func extractPromptContent(manifest map[string]interface{}) string {
	kwargs, ok := manifest["kwargs"].(map[string]interface{})
	if !ok {
		return "Unable to parse prompt content"
	}

	messagesObj, ok := kwargs["messages"]
	if !ok {
		return "Unable to parse prompt content"
	}

	messages, ok := messagesObj.([]interface{})
	if !ok {
		return "Unable to parse prompt content"
	}

	parts := make([]string, 0, len(messages))
	for i, messageObj := range messages {
		message, ok := messageObj.(map[string]interface{})
		if !ok {
			continue
		}

		msgKwargs, ok := message["kwargs"].(map[string]interface{})
		if !ok {
			continue
		}

		prompt, ok := msgKwargs["prompt"].(map[string]interface{})
		if !ok {
			continue
		}

		promptKwargs, ok := prompt["kwargs"].(map[string]interface{})
		if !ok {
			continue
		}

		template, ok := promptKwargs["template"].(string)
		if !ok {
			continue
		}

		messageType := getMessageType(message)
		prefix := ""
		if i > 0 {
			prefix = " | "
		}
		parts = append(parts, fmt.Sprintf("%s%s: \"%s\"", prefix, messageType, template))
	}

	if len(parts) == 0 {
		return "Unable to parse prompt content"
	}

	return strings.Join(parts, "")
}
