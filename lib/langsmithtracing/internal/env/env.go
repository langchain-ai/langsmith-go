package env

import (
	"log"
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
// Returns nil if unset, meaning all traces are kept.
func TracingSampleRate() *float64 {
	s := os.Getenv("LANGSMITH_TRACING_SAMPLING_RATE")
	if s == "" {
		s = os.Getenv("LANGCHAIN_TRACING_SAMPLING_RATE")
	}
	if s == "" {
		return nil
	}
	rate, err := strconv.ParseFloat(s, 64)
	if err != nil {
		log.Printf("[langsmith] invalid LANGSMITH_TRACING_SAMPLING_RATE %q: %v", s, err)
		return nil
	}
	if rate < 0 || rate > 1 {
		log.Printf("[langsmith] LANGSMITH_TRACING_SAMPLING_RATE must be between 0 and 1; got %f", rate)
		return nil
	}
	return &rate
}

var (
	envMetadataOnce sync.Once
	envMetadataMap  map[string]any

	runtimeEnvOnce sync.Once
	runtimeEnvMap  map[string]any
)

// RuntimeEnvironment returns SDK and platform info for injection into extra.runtime,
// matching the Python SDK's get_runtime_environment.
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
// suitable for injection into extra.metadata, matching the Python SDK's
// get_langchain_env_var_metadata and the langgraph-api's getDefaultMetadata.
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
