package langsmith

// ContextHubMountParams describes a Context Hub repo mounted into a sandbox.
type ContextHubMountParams struct {
	// ID is a stable identifier for the mount.
	ID string
	// MountPath is the absolute guest path the repo tree is mirrored into.
	// Unlike bucket and git mounts it is not restricted to /mnt/mounts.
	MountPath string
	// Repo is the Context Hub repository, as "owner/repo" (e.g. "-/my-agent",
	// where "-" is the current workspace).
	Repo string
	// InitialPullOnly syncs the repo once at startup instead of polling for
	// updates for the sandbox's lifetime.
	InitialPullOnly bool
}

// ContextHubMount builds a read-only Context Hub mount spec for
// [SandboxBoxNewParams]. The repo's latest commit tree is mirrored into
// MountPath and the caller's credentials must have access to the repo. The sync
// is one-way: files written under MountPath in the sandbox are never pushed
// back to the repo, and the next sync overwrites them.
func ContextHubMount(params ContextHubMountParams) SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpec {
	contextHub := SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecContexthub{
		Repo: F(params.Repo),
	}
	if params.InitialPullOnly {
		contextHub.InitialPullOnly = F(true)
	}

	return SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpec{
		ID:         F(params.ID),
		Type:       F(SandboxBoxNewParamsMountConfigMountsSandboxapiContextHubRepoMountSpecTypeContexthub),
		MountPath:  F(params.MountPath),
		Contexthub: F(contextHub),
	}
}
