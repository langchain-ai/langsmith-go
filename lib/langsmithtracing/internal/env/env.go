package env

import (
	"fmt"
	"os"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/langchain-ai/langsmith-go/internal"
)

const defaultAPIURL = "https://api.smith.langchain.com"

// APIURL returns the LangSmith API base URL from LANGSMITH_ENDPOINT or LANGCHAIN_ENDPOINT.
func APIURL() string {
	if v := os.Getenv("LANGSMITH_ENDPOINT"); v != "" {
		return v
	}
	if v := os.Getenv("LANGCHAIN_ENDPOINT"); v != "" {
		return v
	}
	return defaultAPIURL
}

// APIKey returns the API key from LANGSMITH_API_KEY or LANGCHAIN_API_KEY.
func APIKey() string {
	if v := os.Getenv("LANGSMITH_API_KEY"); v != "" {
		return v
	}
	return os.Getenv("LANGCHAIN_API_KEY")
}

// Project returns the project name from LANGSMITH_PROJECT or LANGCHAIN_PROJECT.
func Project() string {
	if v := os.Getenv("LANGSMITH_PROJECT"); v != "" {
		return v
	}
	if v := os.Getenv("LANGCHAIN_PROJECT"); v != "" {
		return v
	}
	return "default"
}

// TracingSampleRate returns the sampling rate from
// LANGSMITH_TRACING_SAMPLING_RATE or LANGCHAIN_TRACING_SAMPLING_RATE.
// Returns (nil, nil) if unset, meaning all traces are kept.
// Returns an error if the variable is set but not a valid float in [0, 1].
func TracingSampleRate() (*float64, error) {
	envName := "LANGSMITH_TRACING_SAMPLING_RATE"
	s := os.Getenv(envName)
	if s == "" {
		envName = "LANGCHAIN_TRACING_SAMPLING_RATE"
		s = os.Getenv(envName)
	}
	if s == "" {
		return nil, nil
	}
	rate, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return nil, fmt.Errorf("langsmith: invalid %s %q: %w", envName, s, err)
	}
	if rate < 0 || rate > 1 {
		return nil, fmt.Errorf("langsmith: %s must be between 0 and 1; got %f", envName, rate)
	}
	return &rate, nil
}

// CompressionDisabled returns true if LANGSMITH_DISABLE_RUN_COMPRESSION or
// LANGCHAIN_DISABLE_RUN_COMPRESSION is set to a truthy value ("true", "1", "yes").
func CompressionDisabled() bool {
	for _, key := range []string{
		"LANGSMITH_DISABLE_RUN_COMPRESSION",
		"LANGCHAIN_DISABLE_RUN_COMPRESSION",
	} {
		v := strings.ToLower(strings.TrimSpace(os.Getenv(key)))
		if v == "true" || v == "1" || v == "yes" {
			return true
		}
	}
	return false
}

var (
	envMetadataOnce sync.Once
	envMetadataMap  map[string]any

	runtimeEnvOnce sync.Once
	runtimeEnvMap  map[string]any
)

// RuntimeEnvironment returns SDK and platform info for injection into extra.runtime.
// The same map may be returned on each call; do not mutate it—copy first if you need to edit.
func RuntimeEnvironment() map[string]any {
	runtimeEnvOnce.Do(func() {
		runtimeEnvMap = map[string]any{
			"library":         "langsmith",
			"sdk":             "langsmith-go",
			"sdk_version":     internal.PackageVersion,
			"runtime":         "go",
			"runtime_version": runtime.Version(),
			"platform":        runtime.GOOS + "/" + runtime.GOARCH,
		}
	})
	return runtimeEnvMap
}

// LangChainEnvMetadata returns filtered LANGCHAIN_*/LANGSMITH_* env vars
// suitable for merging into extra.metadata. The tracing client only uses this when
// WithMergeFilteredEnvIntoExtraMetadata(true) is set.
// The same map may be returned on each call; do not mutate it—copy first if you need to edit.
func LangChainEnvMetadata() map[string]any {
	envMetadataOnce.Do(func() {
		excluded := map[string]bool{
			"LANGSMITH_API_KEY":    true,
			"LANGCHAIN_API_KEY":    true,
			"LANGCHAIN_PROJECT":    true,
			"LANGSMITH_PROJECT":    true,
			"LANGCHAIN_SESSION":    true,
			"LANGSMITH_TRACING":    true,
			"LANGCHAIN_TRACING":    true,
			"LANGCHAIN_TRACING_V2": true,
		}
		envMetadataMap = make(map[string]any)
		var revisionID string
		for _, e := range os.Environ() {
			k, v, ok := strings.Cut(e, "=")
			if !ok {
				continue
			}
			if !strings.HasPrefix(k, "LANGCHAIN_") && !strings.HasPrefix(k, "LANGSMITH_") {
				continue
			}
			if excluded[k] {
				continue
			}
			lower := strings.ToLower(k)
			if strings.Contains(lower, "key") || strings.Contains(lower, "secret") ||
				strings.Contains(lower, "token") || strings.Contains(lower, "endpoint") {
				continue
			}
			if k == "LANGCHAIN_REVISION_ID" {
				revisionID = v
				continue
			}
			envMetadataMap[k] = v
		}
		if revisionID != "" {
			envMetadataMap["revision_id"] = revisionID
		}
	})
	return envMetadataMap
}
