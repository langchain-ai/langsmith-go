package cmd

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/spf13/cobra"
)

// optionalInt64 is an int flag that records whether it was set, so 0 can be
// sent as a value (for example idle_ttl_seconds=0 disables the idle stop).
type optionalInt64 struct {
	value int64
	set   bool
}

func (o *optionalInt64) String() string {
	if !o.set {
		return ""
	}
	return strconv.FormatInt(o.value, 10)
}

func (o *optionalInt64) Set(s string) error {
	v, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return fmt.Errorf("must be an integer")
	}
	o.value, o.set = v, true
	return nil
}

func (o *optionalInt64) Type() string { return "int" }

// parseKeyValues parses repeated KEY=VALUE flag values into a map.
func parseKeyValues(flag string, specs []string) (map[string]string, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(specs))
	for _, spec := range specs {
		key, value, ok := strings.Cut(spec, "=")
		if !ok || key == "" {
			return nil, fmt.Errorf("--%s must be KEY=VALUE (got %q)", flag, spec)
		}
		out[key] = value
	}
	return out, nil
}

// parseEnvSpecs parses repeated KEY or KEY=VALUE flag values. A bare KEY takes
// its value from the local environment.
func parseEnvSpecs(flag string, specs []string) (map[string]string, error) {
	if len(specs) == 0 {
		return nil, nil
	}
	out := make(map[string]string, len(specs))
	for _, spec := range specs {
		key, value, hasValue := strings.Cut(spec, "=")
		if key == "" {
			return nil, fmt.Errorf("--%s must be KEY or KEY=VALUE", flag)
		}
		if !hasValue {
			local, ok := os.LookupEnv(key)
			if !ok {
				return nil, fmt.Errorf("--%s %s is not set in the local environment", flag, key)
			}
			value = local
		}
		out[key] = value
	}
	return out, nil
}

// sandboxVolumeSpec is one --volume VOLUME:/mount/path[:ro|:rw] value.
type sandboxVolumeSpec struct {
	Volume    string
	MountPath string
	ReadOnly  bool
}

func parseVolumeSpecs(specs []string) ([]sandboxVolumeSpec, error) {
	out := make([]sandboxVolumeSpec, 0, len(specs))
	for _, spec := range specs {
		parts := strings.Split(spec, ":")
		if len(parts) < 2 || len(parts) > 3 || parts[0] == "" || parts[1] == "" {
			return nil, fmt.Errorf("--volume must be VOLUME:/mount/path[:ro] (got %q)", spec)
		}
		v := sandboxVolumeSpec{Volume: parts[0], MountPath: parts[1]}
		if len(parts) == 3 {
			switch parts[2] {
			case "ro":
				v.ReadOnly = true
			case "rw":
			default:
				return nil, fmt.Errorf("--volume mode must be ro or rw (got %q)", parts[2])
			}
		}
		out = append(out, v)
	}
	return out, nil
}

// sandboxRunConfigInput is the user, working directory and environment
// commands in a sandbox run with.
type sandboxRunConfigInput struct {
	User    string
	WorkDir string
	EnvVars []string
}

// addSandboxRunConfigFlags registers --user and --work-dir, and --env-var when
// envUsage is not empty.
func addSandboxRunConfigFlags(cmd *cobra.Command, in *sandboxRunConfigInput, envUsage string) {
	cmd.Flags().StringVar(&in.User, "user", in.User, "User commands run as: name, uid, name:group or uid:gid")
	cmd.Flags().StringVar(&in.WorkDir, "work-dir", in.WorkDir, "Absolute working directory commands start in")
	if envUsage != "" {
		cmd.Flags().StringArrayVar(&in.EnvVars, "env-var", nil, envUsage)
	}
}

func (in *sandboxRunConfigInput) empty() bool {
	return in.User == "" && in.WorkDir == "" && len(in.EnvVars) == 0
}

func (in *sandboxRunConfigInput) envVars() (map[string]string, error) {
	return parseKeyValues("env-var", in.EnvVars)
}

// sandboxListFilterInput is the filters shared by sandbox and snapshot lists.
type sandboxListFilterInput struct {
	NameContains string
	Status       string
	Mine         bool
	Labels       []string
	Limit        int
}

func addSandboxListFilterFlags(cmd *cobra.Command, in *sandboxListFilterInput, statuses string) {
	cmd.Flags().StringVar(&in.NameContains, "name-contains", in.NameContains, "Only names containing this text")
	cmd.Flags().StringVar(&in.Status, "status", in.Status, "Only this status: "+statuses)
	cmd.Flags().BoolVar(&in.Mine, "mine", in.Mine, "Only those you created (needs a user identity, not a service key)")
	cmd.Flags().StringArrayVar(&in.Labels, "label", nil, "Only those with this label: KEY to require the key, KEY=VALUE to match the value (repeatable; all must match)")
	cmd.Flags().IntVar(&in.Limit, "limit", in.Limit, "Maximum number of results (0 for all)")
}
