package cmd

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"sort"
	"strings"
	"time"

	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/structured"
	"github.com/spf13/cobra"
)

const envGatewayURL = "LANGSMITH_GATEWAY_URL"

// gatewayRequestTimeout bounds one gateway API call. Listing models asks every
// configured provider, so it allows for a slow upstream.
const gatewayRequestTimeout = 60 * time.Second

// gatewayErrorBodyLimit caps how much of an error response is shown.
const gatewayErrorBodyLimit = 2 << 10

// gatewayHTTPClient does not follow redirects. Go drops Authorization on a
// redirect to another host but forwards X-API-Key and X-Tenant-Id, and the
// gateway API has no reason to redirect, so a redirect is reported instead.
var gatewayHTTPClient = &http.Client{
	CheckRedirect: func(*http.Request, []*http.Request) error { return http.ErrUseLastResponse },
}

var gatewayCommand = structured.Parent{
	Use:   "gateway",
	Short: "Use the LangSmith LLM Gateway (experimental)",
	Long: `Use the LangSmith LLM Gateway with the same LangSmith credentials as every
other command: an OAuth login from 'langsmith auth login', or an API key.

The gateway URL is derived from the API URL (api.smith.langchain.com uses
gateway.smith.langchain.com; self-hosted uses <origin>/gateway). Override it
with --gateway-url or LANGSMITH_GATEWAY_URL.

To run a command, such as Claude Code, with its Anthropic and OpenAI requests
sent through the gateway, use 'langsmith gateway wrap claude'.`,
	Children: []func() *cobra.Command{
		gatewayModelsCommand.Cobra,
		newGatewayWrapCmd,
	},
}

type gatewayModelsInput struct {
	Provider   string
	GatewayURL string
}

// gatewayModel is one entry of the gateway's OpenAI-shaped /v1/models list.
// OwnedBy is the upstream listing's owner (an OpenAI organization, for
// example), not necessarily the provider the gateway routes the model to.
type gatewayModel struct {
	ID                 string   `json:"id"`
	OwnedBy            string   `json:"owned_by"`
	SupportedEndpoints []string `json:"supported_endpoints,omitempty"`
}

var gatewayModelsCommand = structured.Command[*gatewayModelsInput]{
	Use:   "models",
	Short: "List the models the LLM Gateway serves your workspace",
	Long: `List the models the LLM Gateway serves your workspace, across every configured
provider and LangChain-hosted models, after your workspace's model access
policies. Use a model's ID as the model name in requests to the gateway's
unified endpoints (/v1/chat/completions, /v1/messages, /v1/responses): a
<provider>/<model> ID routes to that provider, and an ID without a provider
prefix is a LangChain-hosted model.

  langsmith gateway models
  langsmith gateway models --provider anthropic
  langsmith gateway models --jq '.[].id'`,
	Args: cobra.NoArgs,
	Input: func(cmd *cobra.Command) *gatewayModelsInput {
		in := &gatewayModelsInput{}
		addGatewayURLFlag(cmd, &in.GatewayURL)
		cmd.Flags().StringVar(&in.Provider, "provider", "", "Only list models routed to this provider, whose IDs start with <provider>/ (for example anthropic, openai, baseten)")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *gatewayModelsInput, args []string) (any, error) {
		if ctx == nil {
			ctx = context.Background()
		}
		var list struct {
			Data []gatewayModel `json:"data"`
		}
		if err := gatewayGet(ctx, in.GatewayURL, "/v1/models", &list); err != nil {
			return nil, err
		}
		prefix := strings.ToLower(strings.TrimSuffix(in.Provider, "/")) + "/"
		models := make([]gatewayModel, 0, len(list.Data))
		for _, m := range list.Data {
			if in.Provider == "" || strings.HasPrefix(strings.ToLower(m.ID), prefix) {
				models = append(models, m)
			}
		}
		sort.SliceStable(models, func(i, j int) bool { return models[i].ID < models[j].ID })
		return models, nil
	},
	Render: structured.Table{
		Rows: ".",
		Columns: []structured.Column{
			{Header: "Model", Template: "{{.ID}}"},
			{Header: "Owned by", Template: "{{dash .OwnedBy}}"},
			{Header: "Endpoints", Template: "{{joinOrDash .SupportedEndpoints}}"},
		},
	},
}

// addGatewayURLFlag registers --gateway-url, which overrides the gateway URL
// derived from the API URL.
func addGatewayURLFlag(cmd *cobra.Command, target *string) {
	cmd.Flags().StringVar(target, "gateway-url", "", "LLM Gateway URL (default: derived from the API URL) [env: "+envGatewayURL+"]")
}

// gatewayGet sends an authenticated GET to the gateway and decodes its JSON
// response into out.
func gatewayGet(ctx context.Context, gatewayFlag, path string, out any) error {
	opts, err := resolveAuthenticatedOptions(currentGlobalFlags())
	if err != nil {
		return err
	}
	gatewayURL, err := resolveGatewayURL(gatewayFlag, opts.APIURL)
	if err != nil {
		return err
	}
	headers, err := sdkGatewayAuth(client.NewWithOptions(opts))(ctx)
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(ctx, gatewayRequestTimeout)
	defer cancel()
	req, err := http.NewRequestWithContext(ctx, http.MethodGet, gatewayURL+path, nil)
	if err != nil {
		return err
	}
	for name, values := range headers {
		req.Header[name] = values
	}
	req.Header.Set("Accept", "application/json")

	resp, err := gatewayHTTPClient.Do(req)
	if err != nil {
		return fmt.Errorf("calling the LLM Gateway at %s: %w", gatewayURL, err)
	}
	defer resp.Body.Close()
	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		body, _ := io.ReadAll(io.LimitReader(resp.Body, gatewayErrorBodyLimit))
		return fmt.Errorf("LLM Gateway returned HTTP %d for GET %s: %s", resp.StatusCode, path, strings.TrimSpace(string(body)))
	}
	if err := json.NewDecoder(resp.Body).Decode(out); err != nil {
		return fmt.Errorf("decoding LLM Gateway response from GET %s: %w", path, err)
	}
	return nil
}

// gatewayAuth returns the LangSmith auth headers to send to the gateway.
type gatewayAuth func(ctx context.Context) (http.Header, error)

// sdkGatewayAuth returns the headers the SDK would authenticate a LangSmith
// request with, following the same flag, environment and profile precedence
// as every other command. The SDK refreshes an expiring OAuth token under a
// lock shared with other processes.
func sdkGatewayAuth(c *client.Client) gatewayAuth {
	return func(ctx context.Context) (http.Header, error) {
		headers, err := c.SDK.AuthHeaders(ctx)
		if err != nil {
			return nil, err
		}
		if headers.Get("Authorization") == "" && headers.Get("X-API-Key") == "" {
			return nil, errors.New("no LangSmith credentials; run 'langsmith auth login' or set LANGSMITH_API_KEY")
		}
		return headers, nil
	}
}

// describeGatewayAuth names the credential kind, without revealing it.
func describeGatewayAuth(headers http.Header, profile string) string {
	kind := "LangSmith API key"
	if strings.HasPrefix(headers.Get("Authorization"), "Bearer ") {
		kind = "LangSmith OAuth"
	}
	if profile != "" {
		return fmt.Sprintf("%s, profile %q", kind, profile)
	}
	return kind
}

// resolveGatewayURL returns the gateway origin (with any base path), without a
// trailing slash. It follows the frontend's gatewayBaseUrl: SaaS hosts map
// <prefix>.api.smith.langchain.com to <prefix>.gateway.smith.langchain.com,
// and other deployments serve the gateway at <origin><basePath>/gateway.
func resolveGatewayURL(flagValue, apiURL string) (string, error) {
	raw := strings.TrimSpace(flagValue)
	if raw == "" {
		raw = strings.TrimSpace(os.Getenv(envGatewayURL))
	}
	if raw == "" {
		root := strings.TrimSuffix(strings.TrimRight(client.NormalizeURL(apiURL), "/"), "/api")
		u, err := url.Parse(root)
		if err != nil || u.Host == "" {
			return "", fmt.Errorf("cannot derive the gateway URL from API URL %q; pass --gateway-url", apiURL)
		}
		host := u.Hostname()
		const appSuffix = "smith.langchain.com"
		apiSuffix := "api." + appSuffix
		switch {
		case host == apiSuffix || strings.HasSuffix(host, "."+apiSuffix):
			raw = u.Scheme + "://" + strings.TrimSuffix(host, apiSuffix) + "gateway." + appSuffix
		case host == appSuffix || strings.HasSuffix(host, "."+appSuffix):
			raw = u.Scheme + "://" + strings.TrimSuffix(host, appSuffix) + "gateway." + appSuffix
		default:
			raw = root + "/gateway"
		}
	}
	u, err := url.Parse(strings.TrimRight(raw, "/"))
	if err != nil || (u.Scheme != "https" && u.Scheme != "http") || u.Host == "" || u.User != nil || u.RawQuery != "" || u.Fragment != "" {
		return "", fmt.Errorf("invalid gateway URL %q: must be an http(s) URL with no credentials, query, or fragment", raw)
	}
	return u.String(), nil
}
