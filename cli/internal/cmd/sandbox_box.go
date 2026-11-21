package cmd

import (
	"context"
	"fmt"
	"os"
	"slices"
	"strings"
	"time"

	"github.com/google/uuid"
	langsmith "github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/cli/internal/client"
	"github.com/langchain-ai/langsmith-go/cli/internal/cmdutil"
	"github.com/langchain-ai/langsmith-go/cli/internal/structured"
	"github.com/spf13/cobra"
)

const defaultBoxPollInterval = 2 * time.Second

func waitForBoxReady(ctx context.Context, c *client.Client, name string) (*langsmith.SandboxStatusResponse, error) {
	for {
		resp, err := c.SDK.Sandboxes.Boxes.GetStatus(ctx, name)
		if err != nil {
			return nil, err
		}
		switch resp.Status {
		case "ready":
			return resp, nil
		case "failed":
			return resp, fmt.Errorf("sandbox entered %s state", resp.Status)
		}
		timer := time.NewTimer(defaultBoxPollInterval)
		select {
		case <-ctx.Done():
			timer.Stop()
			return nil, ctx.Err()
		case <-timer.C:
		}
	}
}

type sandboxCreateInput struct {
	Snapshot      string
	SnapshotID    string
	VCPUs         int
	CPUMillicores int
	Memory        string
	RootFS        string
	ProxyConfig   string
	MountConfig   string

	DelegatePermissions []string

	IdleTTLSeconds         optionalInt64
	DeleteAfterStopSeconds optionalInt64
	NoPersistent           bool
	PreserveMemoryOnStop   bool
	// PreserveMemoryOnStopSet is true when the flag was given, so an explicit
	// false overrides a snapshot that carries memory preservation.
	PreserveMemoryOnStopSet bool
	RestoreMemory           string

	SandboxEnvVars []string
	Labels         []string
	TagValueIDs    []string
	ResourceTags   []string
	Volumes        []string
	RunConfig      sandboxRunConfigInput

	Wait            bool
	Console         bool
	Shell           string
	ForwardSSHAgent bool
	Env             []string
}

type sandboxServiceURLInput struct {
	Port             int
	ExpiresInSeconds int64
	Access           string
}

type sandboxDownloadURLInput struct {
	Path               string
	ExpiresInSeconds   int64
	ContentType        string
	ContentDisposition string
}

var sandboxBoxDetailRender = structured.PropertyList{
	Properties: []structured.Property{
		{Label: "Name", Template: "{{.Name}}"},
		{Label: "ID", Template: "{{.ID}}"},
		{Label: "Status", Template: "{{.Status}}"},
		{Label: "Size", Template: "{{.SizeClass}}"},
		{Label: "vCPU", Template: "{{formatCount .Vcpus}}"},
		{Label: "Memory", Template: "{{formatBytesOrDash .MemBytes}}"},
		{Label: "Rootfs", Template: "{{formatBytesOrDash .FsCapacityBytes}}"},
		{Label: "Snapshot", Template: "{{shortID .SnapshotID}}"},
		{Label: "Status message", Template: "{{.StatusMessage}}", OmitEmpty: true},
		{Label: "Idle TTL", Template: "{{formatCount .IdleTtlSeconds}}s"},
		{Label: "Delete after stop", Template: "{{if gt .DeleteAfterStopSeconds 0}}{{.DeleteAfterStopSeconds}}s{{end}}", OmitEmpty: true},
		{Label: "Persistent", Template: "{{.Persistent}}"},
		{Label: "Preserve memory", Template: "{{.PreserveMemoryOnStop}}"},
		{Label: "Labels", Template: "{{formatKV .Labels}}", OmitEmpty: true},
		{Label: "Run as", Template: "{{with .RunConfig.User}}{{.}}{{end}}{{with .RunConfig.WorkDir}} in {{.}}{{end}}", OmitEmpty: true},
		{
			Label:     "Volumes",
			Template:  "{{range $i, $v := .Volumes}}{{if $i}}\n{{end}}{{$v.VolumeName}} -> {{$v.MountPath}}{{if $v.ReadOnly}} (ro){{end}}{{end}}",
			OmitEmpty: true,
		},
		{
			Label:     "Delegated access",
			Template:  "{{.AccessDelegation.Mode}}{{with .AccessDelegation.Permissions}} ({{joinOrDash .}}){{end}}",
			OmitEmpty: true,
		},
		{Label: "Created", Template: "{{formatTime .CreatedAt}}"},
		{Label: "Stopped", Template: "{{with .StoppedAt}}{{formatTime .}}{{end}}", OmitEmpty: true},
	},
}

func sandboxCreateParams(name string, in *sandboxCreateInput) (langsmith.SandboxBoxNewParams, error) {
	params := langsmith.SandboxBoxNewParams{}
	if in.VCPUs != 0 {
		params.Vcpus = langsmith.F(int64(in.VCPUs))
	}
	if in.Memory != "" {
		memBytes, err := parseByteSize(in.Memory)
		if err != nil {
			return langsmith.SandboxBoxNewParams{}, fmt.Errorf("invalid --memory: %w", err)
		}
		params.MemBytes = langsmith.F(memBytes)
	}
	if name != "" {
		params.Name = langsmith.F(name)
	}
	if in.SnapshotID != "" {
		params.SnapshotID = langsmith.F(in.SnapshotID)
	}
	if in.Snapshot != "" {
		// The server takes snapshot_id (UUID) and snapshot_name (name or name:tag) as distinct fields.
		if _, err := uuid.Parse(in.Snapshot); err == nil {
			params.SnapshotID = langsmith.F(in.Snapshot)
		} else {
			params.SnapshotName = langsmith.F(in.Snapshot)
		}
	}
	if in.RootFS != "" {
		rootfsBytes, err := parseByteSize(in.RootFS)
		if err != nil {
			return langsmith.SandboxBoxNewParams{}, fmt.Errorf("invalid --rootfs-capacity: %w", err)
		}
		params.FsCapacityBytes = langsmith.F(rootfsBytes)
	}
	if in.ProxyConfig != "" {
		pc, err := loadJSONArg(in.ProxyConfig)
		if err != nil {
			return langsmith.SandboxBoxNewParams{}, fmt.Errorf("invalid --proxy-config: %w", err)
		}
		params.ProxyConfig = langsmith.Raw[langsmith.ProxyConfigParam](pc)
	}
	if len(in.DelegatePermissions) > 0 {
		params.AccessDelegation = langsmith.F(langsmith.SandboxBoxNewParamsAccessDelegation{
			Mode:        langsmith.F(langsmith.SandboxBoxNewParamsAccessDelegationModeExplicit),
			Permissions: langsmith.F(in.DelegatePermissions),
		})
	}
	if in.CPUMillicores != 0 {
		params.CPUMillicores = langsmith.F(int64(in.CPUMillicores))
	}
	if in.MountConfig != "" {
		mc, err := loadJSONArg(in.MountConfig)
		if err != nil {
			return langsmith.SandboxBoxNewParams{}, fmt.Errorf("invalid --mount-config: %w", err)
		}
		params.MountConfig = langsmith.Raw[langsmith.SandboxMountConfigParam](mc)
	}
	if in.IdleTTLSeconds.set {
		params.IdleTtlSeconds = langsmith.F(in.IdleTTLSeconds.value)
	}
	if in.DeleteAfterStopSeconds.set {
		params.DeleteAfterStopSeconds = langsmith.F(in.DeleteAfterStopSeconds.value)
	}
	if in.NoPersistent {
		if in.PreserveMemoryOnStopSet && in.PreserveMemoryOnStop {
			return langsmith.SandboxBoxNewParams{}, fmt.Errorf("--no-persistent cannot be combined with --preserve-memory-on-stop")
		}
		params.Persistent = langsmith.F(false)
	}
	if in.PreserveMemoryOnStopSet {
		params.PreserveMemoryOnStop = langsmith.F(in.PreserveMemoryOnStop)
	}
	switch in.RestoreMemory {
	case "", "if-present":
	case "always":
		params.RestoreMemory = langsmith.F(true)
	case "never":
		params.RestoreMemory = langsmith.F(false)
	default:
		return langsmith.SandboxBoxNewParams{}, fmt.Errorf("--restore-memory must be always, never or if-present (got %q)", in.RestoreMemory)
	}
	envVars, err := parseKeyValues("env-var", in.SandboxEnvVars)
	if err != nil {
		return langsmith.SandboxBoxNewParams{}, err
	}
	if len(envVars) > 0 {
		params.EnvVars = langsmith.F(envVars)
	}
	labels, err := parseKeyValues("label", in.Labels)
	if err != nil {
		return langsmith.SandboxBoxNewParams{}, err
	}
	if len(labels) > 0 {
		params.Labels = langsmith.F(labels)
	}
	if len(in.TagValueIDs) > 0 {
		params.TagValueIDs = langsmith.F(in.TagValueIDs)
	}
	volumes, err := parseVolumeSpecs(in.Volumes)
	if err != nil {
		return langsmith.SandboxBoxNewParams{}, err
	}
	if len(volumes) > 0 {
		attach := make([]langsmith.SandboxBoxNewParamsVolume, 0, len(volumes))
		for _, v := range volumes {
			p := langsmith.SandboxBoxNewParamsVolume{
				Volume:    langsmith.F(v.Volume),
				MountPath: langsmith.F(v.MountPath),
			}
			if v.ReadOnly {
				p.ReadOnly = langsmith.F(true)
			}
			attach = append(attach, p)
		}
		params.Volumes = langsmith.F(attach)
	}
	if !in.RunConfig.empty() {
		runEnv, err := in.RunConfig.envVars()
		if err != nil {
			return langsmith.SandboxBoxNewParams{}, err
		}
		rc := langsmith.SandboxBoxNewParamsRunConfig{}
		if in.RunConfig.User != "" {
			rc.User = langsmith.F(in.RunConfig.User)
		}
		if in.RunConfig.WorkDir != "" {
			rc.WorkDir = langsmith.F(in.RunConfig.WorkDir)
		}
		if len(runEnv) > 0 {
			rc.EnvVars = langsmith.F(runEnv)
		}
		params.RunConfig = langsmith.F(rc)
	}
	return params, nil
}

var sandboxCreateCommand = structured.Command[*sandboxCreateInput]{
	Use:   "create [name]",
	Short: "Create a sandbox VM",
	Long: `Create a sandbox VM.

The --proxy-config flag accepts inline JSON or a file path prefixed with @.
The proxy config controls which outbound HTTP requests the sandbox proxy
allows and what headers to inject. Format:

  {
    "rules": [{
      "name": "openai",
      "match_hosts": ["api.openai.com"],
      "match_paths": [],
      "headers": [
        {"name": "Authorization", "type": "opaque", "value": "Bearer sk-..."},
        {"name": "X-Key", "type": "workspace_secret", "value": "Bearer {OPENAI_API_KEY}"}
      ],
      "env_vars": {"OPENAI_API_KEY": "proxy-injected"},
      "enabled": true
    }, {
      "name": "attribution",
      "type": "attribution",
      "headers": [
        {"name": "X-Cost-Center", "type": "plaintext", "value": "engineering"}
      ],
      "enabled": true
    }],
    "access_control": {
      "allow_list": ["*.openai.com", "*.anthropic.com"],
      "deny_list": []
    }
  }

Header types: "plaintext" (literal value), "opaque" (encrypted, hidden in API
responses), "workspace_secret" (resolved from workspace secrets via {KEY}).

An optional type: "attribution" rule discloses literal non-sensitive metadata to
all proxy-handled destinations, independently of credential rules. At most one
attribution rule is allowed, with up to 16 plaintext X- headers and no matching,
provider, relay, rewrite or environment fields. Credentials, reserved names and
templates are rejected, even when disabled. Enabled nonempty attribution cannot
be combined with a nonempty no_proxy list. Never include secrets or personal data.
This is not an egress-enforcement guarantee: direct-IP/nonstandard-port paths,
non-HTTP traffic and previously established opaque tunnels are not covered.
Omit the rule, set enabled to false, or send headers: [] to disable attribution.

A rule's "env_vars" are plaintext variables set for every command in the sandbox
while the rule is enabled, for tools that refuse to run without a credential
variable even though the proxy injects the real credential on the wire.

--delegate-permission lets code inside the sandbox call the LangSmith API as
you, with no API key of its own: LANGSMITH_ENDPOINT and a placeholder
LANGSMITH_API_KEY are set in the sandbox and the proxy swaps the credential in
on the wire. It grants only the permissions you list, and is repeatable.
Permissions are a ceiling re-checked on every request, not a snapshot, so
losing access yourself revokes it here too. Administering the workspace or
organization is never delegatable.

The grant belongs to the sandbox, not to you: anyone who can exec into it can
make calls under the grant, so list only what the sandbox needs.

Lifecycle:

  --idle-ttl-seconds         stop after this long without activity (default 600;
                             0 never stops it)
  --delete-after-stop-seconds delete this long after it stops (0 never deletes it)
  --preserve-memory-on-stop  save RAM on stop (explicit or idle), so the next
                             start resumes running processes instead of booting;
                             =false overrides a snapshot that turns it on
  --restore-memory           how to treat memory saved in the snapshot:
                             if-present (default) resumes from it when there is
                             some, always requires it, never cold-boots
  --no-persistent            discard every write when the sandbox stops; the
                             next start boots fresh from the snapshot. It cannot
                             preserve memory or have its filesystem captured.

--env-var sets a variable for every command in this sandbox. It is not carried
into snapshots captured from the sandbox unless the capture includes env_vars.
--user and --work-dir override the snapshot's user and working directory.

--volume (beta, needs volumes enabled for the workspace) attaches an existing
volume by name or ID at an absolute path, read-write unless suffixed with :ro.
At most 4. A volume has one read-write attachment or any number of read-only ones.

--mount-config accepts inline JSON or @file.json describing bucket, Git and
Context Hub mounts.

Examples:
  langsmith sandboxes create
  langsmith sandboxes create my-vm
  langsmith sandboxes create my-vm --snapshot <id-or-name>
  langsmith sandboxes create my-vm --snapshot-id <id>
  langsmith sandboxes create my-vm --snapshot-id <id> --rootfs-capacity 8gb
  langsmith sandboxes create my-vm --snapshot-id <id> --proxy-config @proxy.json
  langsmith sandboxes create my-vm --snapshot-id <id> --console
  langsmith sandboxes create my-vm --delegate-permission tracer_sessions:read
  langsmith sandboxes create my-vm --preserve-memory-on-stop --idle-ttl-seconds 1800
  langsmith sandboxes create scratch --no-persistent --delete-after-stop-seconds 60
  langsmith sandboxes create my-vm --env-var LOG_LEVEL=debug --label team=infra --wait
  langsmith sandboxes create my-vm --resource-tag env=prod
  langsmith sandboxes create my-vm --volume my-data:/data --volume models:/models:ro`,
	Args: cobra.MaximumNArgs(1),
	Input: func(cmd *cobra.Command) *sandboxCreateInput {
		in := &sandboxCreateInput{}
		cmd.Flags().StringVar(&in.Snapshot, "snapshot", in.Snapshot, "Snapshot ID or name to boot from")
		cmd.Flags().StringVar(&in.SnapshotID, "snapshot-id", in.SnapshotID, "Snapshot ID to boot from")
		cmd.Flags().IntVar(&in.VCPUs, "vcpus", in.VCPUs, "Number of vCPU cores")
		cmd.Flags().StringVar(&in.Memory, "memory", in.Memory, "Memory with unit (e.g. 4gb, 8gb); must be within 50% of 4gb per vCPU")
		cmd.Flags().StringVar(&in.RootFS, "rootfs-capacity", in.RootFS, "Root filesystem capacity with unit (e.g. 4gb, 8gb)")
		cmd.Flags().StringVar(&in.ProxyConfig, "proxy-config", in.ProxyConfig, "Proxy config as JSON or @file.json")
		cmd.Flags().StringArrayVar(&in.DelegatePermissions, "delegate-permission", nil, "Let code inside the sandbox call the LangSmith API as you, limited to this permission (repeatable)")
		cmd.Flags().IntVar(&in.CPUMillicores, "cpu-millicores", in.CPUMillicores, "CPU in millicores (e.g. 500 = 0.5 vCPU); takes precedence over --vcpus")
		cmd.Flags().StringVar(&in.MountConfig, "mount-config", in.MountConfig, "Bucket, Git and Context Hub mounts as JSON or @file.json")
		cmd.Flags().Var(&in.IdleTTLSeconds, "idle-ttl-seconds", "Stop the sandbox after this many idle seconds (0 never stops it; default 600)")
		cmd.Flags().Var(&in.DeleteAfterStopSeconds, "delete-after-stop-seconds", "Delete the sandbox this many seconds after it stops (0 never deletes it)")
		cmd.Flags().BoolVar(&in.NoPersistent, "no-persistent", in.NoPersistent, "Discard the filesystem's writes on every stop; each start boots fresh from the snapshot")
		cmd.Flags().BoolVar(&in.PreserveMemoryOnStop, "preserve-memory-on-stop", in.PreserveMemoryOnStop, "Save memory on stop so the next start resumes where it left off")
		cmd.Flags().StringVar(&in.RestoreMemory, "restore-memory", in.RestoreMemory, "Snapshot memory on boot: if-present (default), always or never")
		cmd.Flags().StringArrayVar(&in.SandboxEnvVars, "env-var", nil, "Environment variable for every command in the sandbox (KEY=VALUE, repeatable)")
		cmd.Flags().StringArrayVar(&in.Labels, "label", nil, "Label to set on the sandbox (KEY=VALUE, repeatable)")
		cmd.Flags().StringArrayVar(&in.TagValueIDs, "tag-value-id", nil, "Resource tag value ID to attach (repeatable)")
		addResourceTagFlag(cmd, &in.ResourceTags, "Resource tag to attach, KEY=VALUE by name (repeatable; the tag must already exist)")
		cmd.Flags().StringArrayVar(&in.Volumes, "volume", nil, "Attach a volume: VOLUME:/mount/path[:ro] (beta, repeatable, at most 4)")
		addSandboxRunConfigFlags(cmd, &in.RunConfig, "")
		cmd.Flags().BoolVar(&in.Wait, "wait", in.Wait, "Wait until the sandbox is ready before printing it")
		cmd.Flags().BoolVar(&in.Console, "console", in.Console, "Open an interactive console once the sandbox is ready")
		cmd.Flags().StringVar(&in.Shell, "shell", in.Shell, "Shell to use for --console (default: sandbox default, usually /bin/bash)")
		cmd.Flags().BoolVar(&in.ForwardSSHAgent, "forward-ssh-agent", in.ForwardSSHAgent, "Forward the local SSH agent (SSH_AUTH_SOCK) into the --console session")
		cmd.Flags().StringArrayVar(&in.Env, "env", nil, "Additional environment variable for the --console session (KEY or KEY=VALUE, repeatable)")
		cmd.Flags().String("jq", "", cmdutil.JQFlagUsage)
		return in
	},
	CustomOutput: true,
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxCreateInput, args []string) (any, error) {
		name := ""
		if len(args) > 0 {
			name = args[0]
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if in.Snapshot != "" && in.SnapshotID != "" {
			return nil, fmt.Errorf("use either --snapshot or --snapshot-id, not both")
		}

		in.PreserveMemoryOnStopSet = cmd.Flags().Changed("preserve-memory-on-stop")
		if in.TagValueIDs, err = withResourceTags(ctx, c, in.TagValueIDs, in.ResourceTags); err != nil {
			return nil, err
		}
		params, err := sandboxCreateParams(name, in)
		if err != nil {
			return nil, err
		}
		resp, err := c.SDK.Sandboxes.Boxes.New(ctx, params)
		if err != nil {
			return nil, fmt.Errorf("creating sandbox: %w", err)
		}

		if !in.Console {
			if !in.Wait {
				return nil, structured.Render(cmd, resp, sandboxBoxDetailRender)
			}
			if _, err := waitForBoxReady(ctx, c, resp.Name); err != nil {
				return nil, err
			}
			ready, err := c.SDK.Sandboxes.Boxes.Get(ctx, resp.Name)
			if err != nil {
				return nil, fmt.Errorf("getting sandbox: %w", err)
			}
			return nil, structured.Render(cmd, ready, sandboxBoxDetailRender)
		}

		if _, err := waitForBoxReady(ctx, c, resp.Name); err != nil {
			return nil, err
		}
		fmt.Fprintf(os.Stderr, "Sandbox %s ready, opening console...\n", resp.Name)
		return nil, runConsole(resp.Name, in.Shell, in.ForwardSSHAgent, in.Env)
	},
}

var sandboxServiceURLRender = structured.PropertyList{
	Properties: []structured.Property{
		{Label: "Browser URL", Template: "{{.BrowserURL}}"},
		{Label: "Service URL", Template: "{{.ServiceURL}}"},
		{Label: "LangSmith login", Template: "{{.Access}}", OmitEmpty: true},
		{Label: "Expires", Template: "{{formatTime .ExpiresAt}}", OmitEmpty: true},
		{Label: "Service Token", Template: "{{.Token}}", OmitEmpty: true},
	},
	Caption: `{{if .Token}}Example:
  curl -H "X-Langsmith-Sandbox-Service-Token: {{.Token}}" "{{.ServiceURL}}"{{else}}Open the URL in a browser; access is gated by your LangSmith login.{{end}}`,
}

var sandboxServiceURLCommand = structured.Command[*sandboxServiceURLInput]{
	Use:   "service-url <name> --port <port>",
	Short: "Generate an authenticated URL for a sandbox HTTP service",
	Long: `Generate an authenticated URL for an HTTP service running inside a sandbox.

By default this mints a short-lived service token: use the service URL with the
X-Langsmith-Sandbox-Service-Token header, or open the browser URL directly.

--access switches to LangSmith login instead, so the URL carries no token and
does not expire:

  restricted  anyone with sandboxes:read on this sandbox
  workspace   any member of the owning workspace

A login grant is durable, so token mode is refused while one is in place, and
--expires-in-seconds does not apply.

Examples:
  langsmith sandboxes service-url my-vm --port 8000
  langsmith sandboxes service-url my-vm --port 8000 --expires-in-seconds 3600
  langsmith sandboxes service-url my-vm --port 8000 --access workspace`,
	Args: cobra.ExactArgs(1),
	Input: func(cmd *cobra.Command) *sandboxServiceURLInput {
		in := &sandboxServiceURLInput{}
		cmd.Flags().IntVar(&in.Port, "port", in.Port, "Port inside the sandbox")
		cmd.Flags().Int64Var(&in.ExpiresInSeconds, "expires-in-seconds", in.ExpiresInSeconds, "URL TTL in seconds (token mode only)")
		cmd.Flags().StringVar(&in.Access, "access", in.Access, "Gate the URL behind LangSmith login instead of a token: restricted or workspace")
		_ = cmd.MarkFlagRequired("port")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxServiceURLInput, args []string) (any, error) {
		if in.Port < 1 || in.Port > 65535 {
			return nil, fmt.Errorf("--port must be between 1 and 65535 (got %d)", in.Port)
		}
		if cmd.Flags().Changed("expires-in-seconds") && in.ExpiresInSeconds < 1 {
			return nil, fmt.Errorf("--expires-in-seconds must be greater than 0")
		}
		access := langsmith.SandboxBoxGenerateServiceURLParamsAccess(in.Access)
		if cmd.Flags().Changed("access") &&
			access != langsmith.SandboxBoxGenerateServiceURLParamsAccessRestricted &&
			access != langsmith.SandboxBoxGenerateServiceURLParamsAccessWorkspace {
			return nil, fmt.Errorf(
				"--access must be restricted or workspace (got %q)", in.Access,
			)
		}
		// A login grant carries no token, so there is nothing for a TTL to expire.
		if cmd.Flags().Changed("expires-in-seconds") && cmd.Flags().Changed("access") {
			return nil, fmt.Errorf(
				"--expires-in-seconds does not apply to --access %s: "+
					"a LangSmith login URL carries no token and does not expire", in.Access,
			)
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		params := langsmith.SandboxBoxGenerateServiceURLParams{
			Port: langsmith.F(int64(in.Port)),
		}
		if cmd.Flags().Changed("expires-in-seconds") {
			params.ExpiresInSeconds = langsmith.F(in.ExpiresInSeconds)
		}
		if cmd.Flags().Changed("access") {
			params.Access = langsmith.F(access)
		}

		resp, err := c.SDK.Sandboxes.Boxes.GenerateServiceURL(ctx, args[0], params)
		if err != nil {
			return nil, fmt.Errorf("generating service URL: %w", err)
		}
		return resp, nil
	},
	Render: sandboxServiceURLRender,
}

var sandboxDownloadURLRender = structured.PropertyList{
	Properties: []structured.Property{
		{Label: "Download URL", Template: "{{.DownloadURL}}"},
		{Label: "Expires", Template: "{{if .ExpiresAt}}{{formatTime .ExpiresAt}}{{else}}never{{end}}"},
	},
	Caption: `Example:
  curl -LO "{{.DownloadURL}}"`,
}

var sandboxDownloadURLCommand = structured.Command[*sandboxDownloadURLInput]{
	Use:   "generate-download-url <name> --path <path>",
	Short: "Generate a link that downloads a sandbox file without credentials",
	Long: `Generate a link that downloads a single file from a sandbox.

The link carries its own token, so anyone with the URL can fetch that one file
with no LangSmith credential. It is pinned to the sandbox and the exact path,
and cannot be repointed at another file. Fetching wakes a stopped sandbox.

Do not modify the file after minting a link for it. The link is pinned to a
path, not to a snapshot of the contents, so a later write to that path may or
may not be reflected in what the link serves. Write a new file and mint a new
link when the contents change.

Without --expires-in-seconds the link never expires.

Examples:
  langsmith sandboxes generate-download-url my-vm --path /tmp/report.pdf
  langsmith sandboxes generate-download-url my-vm --path /tmp/report.pdf --expires-in-seconds 3600
  langsmith sandboxes generate-download-url my-vm --path /tmp/page.html --content-disposition inline`,
	Args: cobra.ExactArgs(1),
	Input: func(cmd *cobra.Command) *sandboxDownloadURLInput {
		in := &sandboxDownloadURLInput{}
		cmd.Flags().StringVar(&in.Path, "path", in.Path, "File path inside the sandbox")
		cmd.Flags().Int64Var(&in.ExpiresInSeconds, "expires-in-seconds", in.ExpiresInSeconds, "Link TTL in seconds (omit for a link that never expires)")
		cmd.Flags().StringVar(&in.ContentType, "content-type", in.ContentType, "Content-Type to serve the file as")
		cmd.Flags().StringVar(&in.ContentDisposition, "content-disposition", in.ContentDisposition, "Content-Disposition to serve the file with (attachment or inline)")
		_ = cmd.MarkFlagRequired("path")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxDownloadURLInput, args []string) (any, error) {
		if cmd.Flags().Changed("expires-in-seconds") && in.ExpiresInSeconds < 1 {
			return nil, fmt.Errorf("--expires-in-seconds must be greater than 0")
		}
		switch in.ContentDisposition {
		case "", "attachment", "inline":
		default:
			return nil, fmt.Errorf("--content-disposition must be attachment or inline (got %q)", in.ContentDisposition)
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		params := langsmith.SandboxBoxGenerateDownloadURLParams{
			Path: langsmith.F(in.Path),
		}
		if cmd.Flags().Changed("expires-in-seconds") {
			params.ExpiresInSeconds = langsmith.F(in.ExpiresInSeconds)
		}
		if in.ContentType != "" {
			params.ContentType = langsmith.F(in.ContentType)
		}
		if in.ContentDisposition != "" {
			params.ContentDisposition = langsmith.F(in.ContentDisposition)
		}

		resp, err := c.SDK.Sandboxes.Boxes.GenerateDownloadURL(ctx, args[0], params)
		if err != nil {
			return nil, fmt.Errorf("generating download URL: %w", err)
		}
		return resp, nil
	},
	Render: sandboxDownloadURLRender,
}

type sandboxListInput struct {
	sandboxListFilterInput
	TagValueIDs  []string
	ResourceTags []string
}

func sandboxListParams(in *sandboxListInput) langsmith.SandboxBoxListParams {
	params := langsmith.SandboxBoxListParams{}
	if in.NameContains != "" {
		params.NameContains = langsmith.F(in.NameContains)
	}
	if in.Status != "" {
		params.Status = langsmith.F(in.Status)
	}
	if in.Mine {
		params.CreatedBy = langsmith.F("me")
	}
	if len(in.Labels) > 0 {
		params.Label = langsmith.F(in.Labels)
	}
	if len(in.TagValueIDs) > 0 {
		params.TagValueID = langsmith.F(in.TagValueIDs)
	}
	return params
}

var sandboxListCommand = structured.Command[*sandboxListInput]{
	Use:   "list",
	Short: "List sandboxes",
	Long: `List sandboxes, optionally filtered.

Examples:
  langsmith sandboxes list
  langsmith sandboxes list --status ready
  langsmith sandboxes list --label team=infra --label ephemeral
  langsmith sandboxes list --resource-tag env=prod
  langsmith sandboxes list --name-contains ci- --limit 20`,
	Input: func(cmd *cobra.Command) *sandboxListInput {
		in := &sandboxListInput{}
		addSandboxListFilterFlags(cmd, &in.sandboxListFilterInput, "provisioning, ready, stopped, failed, deleting")
		cmd.Flags().StringArrayVar(&in.TagValueIDs, "tag-value-id", nil, "Only sandboxes with this resource tag value ID (repeatable)")
		addResourceTagFlag(cmd, &in.ResourceTags, "Only sandboxes with this resource tag, KEY=VALUE by name (repeatable)")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxListInput, args []string) (any, error) {
		if in.Limit < 0 {
			return nil, fmt.Errorf("--limit must not be negative")
		}
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if in.TagValueIDs, err = withResourceTags(ctx, c, in.TagValueIDs, in.ResourceTags); err != nil {
			return nil, err
		}
		var sandboxes []langsmith.SandboxResponse
		pager := c.SDK.Sandboxes.Boxes.ListAutoPaging(ctx, sandboxListParams(in))
		for pager.Next() {
			sandboxes = append(sandboxes, pager.Current())
			if in.Limit > 0 && len(sandboxes) >= in.Limit {
				break
			}
		}
		if err := pager.Err(); err != nil {
			return nil, fmt.Errorf("listing sandboxes: %w", err)
		}
		return sandboxes, nil
	},
	Render: structured.Table{
		Title: "Sandboxes",
		Rows:  ".",
		Columns: []structured.Column{
			{Header: "Name", Template: "{{.Name}}"},
			{Header: "Status", Template: "{{.Status}}"},
			{Header: "vCPU", Template: "{{formatCount .Vcpus}}"},
			{Header: "Mem", Template: "{{formatBytesOrDash .MemBytes}}"},
			{Header: "Rootfs", Template: "{{formatBytesOrDash .FsCapacityBytes}}"},
			{Header: "Snapshot", Template: "{{shortID .SnapshotID}}"},
			{Header: "Memory on stop", Template: "{{if .PreserveMemoryOnStop}}kept{{else}}-{{end}}"},
			{Header: "Labels", Template: "{{dash (formatKV .Labels)}}"},
			{Header: "Created", Template: "{{formatTime .CreatedAt}}"},
		},
	},
}

var sandboxGetCommand = structured.Command[struct{}]{
	Use:   "get <name>",
	Short: "Get a sandbox by name",
	Args:  cobra.ExactArgs(1),
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		resp, err := c.SDK.Sandboxes.Boxes.Get(ctx, args[0])
		if err != nil {
			return nil, fmt.Errorf("getting sandbox: %w", err)
		}

		return resp, nil
	},
	Render: sandboxBoxDetailRender,
}

type sandboxUpdateInput struct {
	Name                   string
	VCPUs                  int
	CPUMillicores          int
	Memory                 string
	RootFS                 string
	ProxyConfig            string
	IdleTTLSeconds         optionalInt64
	DeleteAfterStopSeconds optionalInt64
	TagValueIDs            []string
	ResourceTags           []string
	RunConfig              sandboxRunConfigInput
}

func sandboxUpdateParams(cmd *cobra.Command, in *sandboxUpdateInput) (langsmith.SandboxBoxUpdateParams, error) {
	params := langsmith.SandboxBoxUpdateParams{}
	changed := false
	if cmd.Flags().Changed("name") {
		params.Name = langsmith.F(in.Name)
		changed = true
	}
	if cmd.Flags().Changed("vcpus") {
		params.Vcpus = langsmith.F(int64(in.VCPUs))
		changed = true
	}
	if cmd.Flags().Changed("cpu-millicores") {
		params.CPUMillicores = langsmith.F(int64(in.CPUMillicores))
		changed = true
	}
	if cmd.Flags().Changed("memory") {
		memBytes, err := parseByteSize(in.Memory)
		if err != nil {
			return params, fmt.Errorf("invalid --memory: %w", err)
		}
		params.MemBytes = langsmith.F(memBytes)
		changed = true
	}
	if cmd.Flags().Changed("rootfs-capacity") {
		rootfsBytes, err := parseByteSize(in.RootFS)
		if err != nil {
			return params, fmt.Errorf("invalid --rootfs-capacity: %w", err)
		}
		params.FsCapacityBytes = langsmith.F(rootfsBytes)
		changed = true
	}
	if cmd.Flags().Changed("proxy-config") {
		pc, err := loadJSONArg(in.ProxyConfig)
		if err != nil {
			return params, fmt.Errorf("invalid --proxy-config: %w", err)
		}
		params.ProxyConfig = langsmith.Raw[langsmith.ProxyConfigParam](pc)
		changed = true
	}
	if in.IdleTTLSeconds.set {
		params.IdleTtlSeconds = langsmith.F(in.IdleTTLSeconds.value)
		changed = true
	}
	if in.DeleteAfterStopSeconds.set {
		params.DeleteAfterStopSeconds = langsmith.F(in.DeleteAfterStopSeconds.value)
		changed = true
	}
	if cmd.Flags().Changed("tag-value-id") || cmd.Flags().Changed("resource-tag") {
		params.TagValueIDs = langsmith.F(in.TagValueIDs)
		changed = true
	}
	if !in.RunConfig.empty() {
		rc := langsmith.SandboxBoxUpdateParamsRunConfig{}
		if in.RunConfig.User != "" {
			rc.User = langsmith.F(in.RunConfig.User)
		}
		if in.RunConfig.WorkDir != "" {
			rc.WorkDir = langsmith.F(in.RunConfig.WorkDir)
		}
		params.RunConfig = langsmith.F(rc)
		changed = true
	}
	if !changed {
		return params, fmt.Errorf("nothing to update (use --name, --vcpus, --cpu-millicores, --memory, --rootfs-capacity, --proxy-config, --idle-ttl-seconds, --delete-after-stop-seconds, --tag-value-id, --resource-tag, --user or --work-dir)")
	}
	return params, nil
}

var sandboxUpdateCommand = structured.Command[*sandboxUpdateInput]{
	Use:   "update <name>",
	Short: "Update sandbox settings",
	Long: `Update a sandbox's name, resources, lifecycle, run config or proxy configuration.

Resource changes (--vcpus, --cpu-millicores, --memory, --rootfs-capacity) take
effect on next start. Proxy config, idle TTL and delete-after-stop changes take
effect immediately. --user and --work-dir apply to commands started after the
update; a stopped sandbox picks them up on its next start.

--tag-value-id and --resource-tag replace the sandbox's resource tags with the
ones given.

The --proxy-config flag accepts inline JSON or @file.json. See "create --help"
for the proxy config JSON format.

Examples:
  langsmith sandboxes update my-vm --vcpus 4 --memory 16gb
  langsmith sandboxes update my-vm --name my-renamed-vm
  langsmith sandboxes update my-vm --idle-ttl-seconds 0
  langsmith sandboxes update my-vm --user app --work-dir /srv/app
  langsmith sandboxes update my-vm --resource-tag env=prod --resource-tag team=infra`,
	Args: cobra.ExactArgs(1),
	Input: func(cmd *cobra.Command) *sandboxUpdateInput {
		in := &sandboxUpdateInput{}
		cmd.Flags().StringVar(&in.Name, "name", in.Name, "Rename the sandbox")
		cmd.Flags().IntVar(&in.VCPUs, "vcpus", in.VCPUs, "Number of vCPU cores")
		cmd.Flags().IntVar(&in.CPUMillicores, "cpu-millicores", in.CPUMillicores, "CPU in millicores (e.g. 500 = 0.5 vCPU); takes precedence over --vcpus")
		cmd.Flags().StringVar(&in.Memory, "memory", in.Memory, "Memory with unit (e.g. 4gb, 8gb); must be within 50% of 4gb per vCPU")
		cmd.Flags().StringVar(&in.RootFS, "rootfs-capacity", in.RootFS, "Root filesystem capacity with unit (e.g. 4gb, 8gb)")
		cmd.Flags().StringVar(&in.ProxyConfig, "proxy-config", in.ProxyConfig, "Proxy config as JSON or @file.json")
		cmd.Flags().Var(&in.IdleTTLSeconds, "idle-ttl-seconds", "Stop the sandbox after this many idle seconds (0 never stops it)")
		cmd.Flags().Var(&in.DeleteAfterStopSeconds, "delete-after-stop-seconds", "Delete the sandbox this many seconds after it stops (0 never deletes it)")
		cmd.Flags().StringArrayVar(&in.TagValueIDs, "tag-value-id", nil, "Resource tag value ID; replaces the sandbox's tags (repeatable)")
		addResourceTagFlag(cmd, &in.ResourceTags, "Resource tag, KEY=VALUE by name; replaces the sandbox's tags (repeatable; the tag must already exist)")
		addSandboxRunConfigFlags(cmd, &in.RunConfig, "")
		return in
	},
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxUpdateInput, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}
		if in.TagValueIDs, err = withResourceTags(ctx, c, in.TagValueIDs, in.ResourceTags); err != nil {
			return nil, err
		}
		params, err := sandboxUpdateParams(cmd, in)
		if err != nil {
			return nil, err
		}

		resp, err := c.SDK.Sandboxes.Boxes.Update(ctx, args[0], params)
		if err != nil {
			return nil, fmt.Errorf("updating sandbox: %w", err)
		}

		return resp, nil
	},
	Render: sandboxBoxDetailRender,
}

// sandboxBatchDeleteLimit is the most names one batch-delete request takes.
const sandboxBatchDeleteLimit = 50

type sandboxDeleteResult struct {
	Deleted []string                                         `json:"deleted"`
	Skipped []langsmith.SandboxBoxBatchDeleteResponseSkipped `json:"skipped"`
}

var sandboxDeleteCommand = structured.Command[struct{}]{
	Use:   "delete <name>...",
	Short: "Delete one or more sandboxes",
	Long: `Delete sandboxes by name.

With one name, a sandbox that does not exist is an error. With several, every
sandbox that can be deleted is, the rest are listed with the reason they were
skipped, and the command exits non-zero if any was skipped.

Examples:
  langsmith sandboxes delete my-vm
  langsmith sandboxes delete ci-1 ci-2 ci-3
  langsmith sandboxes list --name-contains ci- --format json --jq '.[].name' | xargs langsmith sandboxes delete`,
	Args: cobra.MinimumNArgs(1),
	Input: func(cmd *cobra.Command) struct{} {
		cmd.Flags().String("jq", "", cmdutil.JQFlagUsage)
		return struct{}{}
	},
	CustomOutput: true,
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if len(args) == 1 {
			if err := c.SDK.Sandboxes.Boxes.Delete(ctx, args[0]); err != nil {
				return nil, fmt.Errorf("deleting sandbox: %w", err)
			}
			return nil, structured.Render(cmd, sandboxMessage{Name: args[0], Message: "Sandbox deleted."}, structured.Template("{{.Message}}\n"))
		}

		result := sandboxDeleteResult{Deleted: []string{}, Skipped: []langsmith.SandboxBoxBatchDeleteResponseSkipped{}}
		var batchErr error
		answered := 0
		for batch := range slices.Chunk(args, sandboxBatchDeleteLimit) {
			resp, err := c.SDK.Sandboxes.Boxes.BatchDelete(ctx, langsmith.SandboxBoxBatchDeleteParams{Names: langsmith.F(batch)})
			if err != nil {
				batchErr = err
				break
			}
			answered += len(batch)
			result.Deleted = append(result.Deleted, resp.Deleted...)
			result.Skipped = append(result.Skipped, resp.Skipped...)
		}
		// Print what earlier batches did even when a later one fails, so the
		// sandboxes already deleted are on record.
		if err := structured.Render(cmd, result, structured.Template(`{{range .Deleted}}Deleted {{.}}
{{end}}{{range .Skipped}}Skipped {{.Name}}: {{.Reason}}
{{end}}`)); err != nil {
			return nil, err
		}
		if batchErr != nil {
			return nil, fmt.Errorf("deleting sandboxes (%d of %d without a confirmed result): %w", len(args)-answered, len(args), batchErr)
		}
		if len(result.Skipped) > 0 {
			return nil, fmt.Errorf("%d of %d sandboxes not deleted", len(result.Skipped), len(args))
		}
		return nil, nil
	},
}

var sandboxStartCommand = structured.Command[struct{}]{
	Use:   "start <name>",
	Short: "Start a stopped sandbox",
	Args:  cobra.ExactArgs(1),
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if _, err := c.SDK.Sandboxes.Boxes.Start(ctx, args[0]); err != nil {
			return nil, fmt.Errorf("starting sandbox: %w", err)
		}

		if _, err := waitForBoxReady(ctx, c, args[0]); err != nil {
			return nil, err
		}
		return sandboxMessage{Name: args[0], Message: fmt.Sprintf("Sandbox %s started", args[0])}, nil
	},
	Render: structured.Template(`{{.Message}}
`),
}

var sandboxStopCommand = structured.Command[struct{}]{
	Use:   "stop <name>",
	Short: "Stop a running sandbox (preserves data)",
	Args:  cobra.ExactArgs(1),
	Action: func(ctx context.Context, cmd *cobra.Command, in struct{}, args []string) (any, error) {
		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		if err := c.SDK.Sandboxes.Boxes.Stop(ctx, args[0]); err != nil {
			return nil, fmt.Errorf("stopping sandbox: %w", err)
		}

		return sandboxMessage{Name: args[0], Message: "Sandbox stopped."}, nil
	},
	Render: structured.Template(`{{.Message}}
`),
}

type sandboxExecInput struct {
	User           string
	WorkDir        string
	Env            []string
	TimeoutSeconds int64
	Shell          string
	Stream         bool
}

func sandboxExecRunConfig(in *sandboxExecInput) (langsmith.SandboxRunConfig, bool, error) {
	env, err := parseEnvSpecs("env", in.Env)
	if err != nil {
		return langsmith.SandboxRunConfig{}, false, err
	}
	rc := langsmith.SandboxRunConfig{}
	set := false
	if in.User != "" {
		rc.User = langsmith.F(in.User)
		set = true
	}
	if in.WorkDir != "" {
		rc.WorkDir = langsmith.F(in.WorkDir)
		set = true
	}
	if len(env) > 0 {
		rc.EnvVars = langsmith.F(env)
		set = true
	}
	return rc, set, nil
}

var sandboxExecCommand = structured.Command[*sandboxExecInput]{
	Use:   "exec <name> -- <command>",
	Short: "Execute a command inside a sandbox",
	Long: `Execute a one-off command inside a running sandbox and print its output.

Output is printed when the command finishes. --stream prints it as it is
produced instead, for long-running commands. The exit code is the command's.

--user, --work-dir and --env apply to this command only, over the sandbox's
own run config. A bare --env KEY takes the value from your local environment.

Examples:
  langsmith sandboxes exec my-vm -- uname -a
  langsmith sandboxes exec my-vm -- ls -la /
  langsmith sandboxes exec my-vm --work-dir /srv/app -- make test
  langsmith sandboxes exec my-vm --user root -- apt-get update
  langsmith sandboxes exec my-vm --env GITHUB_TOKEN --env CI=1 -- ./build.sh
  langsmith sandboxes exec my-vm --stream --timeout 1800 -- npm run build`,
	Args: cobra.MinimumNArgs(1),
	Input: func(cmd *cobra.Command) *sandboxExecInput {
		in := &sandboxExecInput{}
		cmd.Flags().StringVar(&in.User, "user", in.User, "Run as this user: name, uid, name:group or uid:gid")
		cmd.Flags().StringVar(&in.WorkDir, "work-dir", in.WorkDir, "Working directory for the command")
		cmd.Flags().StringArrayVar(&in.Env, "env", nil, "Environment variable for the command (KEY or KEY=VALUE, repeatable)")
		cmd.Flags().Int64Var(&in.TimeoutSeconds, "timeout", in.TimeoutSeconds, "Kill the command after this many seconds")
		cmd.Flags().StringVar(&in.Shell, "shell", in.Shell, "Shell that runs the command (default: sandbox default)")
		cmd.Flags().BoolVar(&in.Stream, "stream", in.Stream, "Print output as it is produced")
		return in
	},
	CustomOutput: true,
	Action: func(ctx context.Context, cmd *cobra.Command, in *sandboxExecInput, args []string) (any, error) {
		name := args[0]

		cmdArgs := cmd.ArgsLenAtDash()
		if cmdArgs < 0 || cmdArgs >= len(args) {
			return nil, fmt.Errorf("usage: langsmith sandboxes exec <name> -- <command>")
		}
		command := args[cmdArgs:]
		if len(command) == 0 {
			return nil, fmt.Errorf("no command specified")
		}
		if in.TimeoutSeconds < 0 {
			return nil, fmt.Errorf("--timeout must not be negative")
		}
		runConfig, hasRunConfig, err := sandboxExecRunConfig(in)
		if err != nil {
			return nil, err
		}

		c, err := cmdutil.GetClient(cmd)
		if err != nil {
			return nil, err
		}

		var result *langsmith.SandboxExecutionResult
		if in.Stream {
			params := langsmith.SandboxCommandStartParams{
				Command: langsmith.F(sandboxShellCommand(command)),
			}
			if in.TimeoutSeconds > 0 {
				params.TimeoutSeconds = langsmith.F(in.TimeoutSeconds)
			}
			if in.Shell != "" {
				params.Shell = langsmith.F(in.Shell)
			}
			if hasRunConfig {
				params.RunConfig = langsmith.F(runConfig)
			}
			result, err = c.SDK.Sandboxes.Boxes.RunWithCallbacks(ctx, name, params, langsmith.SandboxCommandCallbacks{
				OnStdout: func(s string) { fmt.Fprint(os.Stdout, s) },
				OnStderr: func(s string) { fmt.Fprint(os.Stderr, s) },
			})
			if err != nil {
				return nil, fmt.Errorf("execute: %w", err)
			}
		} else {
			params := langsmith.SandboxBoxRunParams{
				Command: langsmith.F(sandboxShellCommand(command)),
			}
			if in.TimeoutSeconds > 0 {
				params.Timeout = langsmith.F(in.TimeoutSeconds)
			}
			if in.Shell != "" {
				params.Shell = langsmith.F(in.Shell)
			}
			if hasRunConfig {
				params.RunConfig = langsmith.F(runConfig)
			}
			result, err = c.SDK.Sandboxes.Boxes.Run(ctx, name, params)
			if err != nil {
				return nil, fmt.Errorf("execute: %w", err)
			}
			if result.Stdout != "" {
				fmt.Print(result.Stdout)
			}
			if result.Stderr != "" {
				fmt.Fprint(os.Stderr, result.Stderr)
			}
		}

		if result.ExitCode != 0 {
			os.Exit(int(result.ExitCode))
		}
		return nil, nil
	},
}

func newSandboxExecCmd() *cobra.Command {
	return sandboxExecCommand.Cobra()
}

func sandboxShellCommand(args []string) string {
	quoted := make([]string, 0, len(args))
	for _, arg := range args {
		quoted = append(quoted, shellQuote(arg))
	}
	return strings.Join(quoted, " ")
}
