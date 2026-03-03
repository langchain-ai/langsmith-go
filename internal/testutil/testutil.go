package testutil

import (
	"net/http"
	"os"
	"strconv"
	"testing"
)

func CheckTestServer(t *testing.T, url string) bool {
	if _, err := http.Get(url); err != nil {
		const SKIP_MOCK_TESTS = "SKIP_MOCK_TESTS"
		if str, ok := os.LookupEnv(SKIP_MOCK_TESTS); ok {
			skip, parseErr := strconv.ParseBool(str)
			if parseErr != nil {
				t.Fatalf("strconv.ParseBool(os.LookupEnv(%s)) failed: %s", SKIP_MOCK_TESTS, parseErr)
			}
			if skip {
				t.Skip("The test will not run without a mock Prism server running against your OpenAPI spec")
				return false
			}
			t.Errorf("The test will not run without a mock Prism server running against your OpenAPI spec. You can set the environment variable %s to true to skip running any tests that require the mock server", SKIP_MOCK_TESTS)
			return false
		}
		t.Skipf("Mock server at %s not reachable (connection refused). Set SKIP_MOCK_TESTS=true to skip explicitly, or start the Prism mock server to run these tests.", url)
		return false
	}
	return true
}
