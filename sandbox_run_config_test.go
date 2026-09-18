package langsmith_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
)

func TestSandboxRunSendsRunConfig(t *testing.T) {
	var body map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == "/execute" {
			require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
			w.Header().Set("Content-Type", "application/json")
			_, _ = w.Write([]byte(`{"stdout":"","stderr":"","exit_code":0}`))
			return
		}
		t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
	}))
	defer srv.Close()

	client := langsmith.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"))
	_, err := client.Sandboxes.Boxes.RunWithDataplaneURL(context.Background(), srv.URL, langsmith.SandboxBoxRunParams{
		Command: langsmith.F("echo hi"),
		RunConfig: langsmith.F(langsmith.SandboxRunConfig{
			User:    langsmith.F("app"),
			WorkDir: langsmith.F("/workspace"),
			EnvVars: langsmith.F(map[string]string{"A": "1"}),
		}),
	})
	require.NoError(t, err)

	assert.Equal(t, map[string]any{
		"user":     "app",
		"work_dir": "/workspace",
		"env_vars": map[string]any{"A": "1"},
	}, body["run_config"])
}

func TestSandboxRunRejectsRunConfigWithDeprecatedFields(t *testing.T) {
	client := langsmith.NewClient(option.WithBaseURL("http://example.invalid"), option.WithAPIKey("k"))

	for name, params := range map[string]langsmith.SandboxBoxRunParams{
		"env": {
			Command:   langsmith.F("echo hi"),
			Env:       langsmith.F(map[string]string{"A": "1"}),
			RunConfig: langsmith.F(langsmith.SandboxRunConfig{User: langsmith.F("app")}),
		},
		"cwd": {
			Command:   langsmith.F("echo hi"),
			CWD:       langsmith.F("/tmp"),
			RunConfig: langsmith.F(langsmith.SandboxRunConfig{User: langsmith.F("app")}),
		},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := client.Sandboxes.Boxes.RunWithDataplaneURL(context.Background(), "http://example.invalid", params)
			require.Error(t, err)
			assert.Contains(t, err.Error(), "cannot combine RunConfig")
		})
	}
}

func TestSandboxRunAcceptsDeprecatedFieldsAlone(t *testing.T) {
	var body map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.NoError(t, json.NewDecoder(r.Body).Decode(&body))
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"stdout":"","stderr":"","exit_code":0}`))
	}))
	defer srv.Close()

	client := langsmith.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"))
	_, err := client.Sandboxes.Boxes.RunWithDataplaneURL(context.Background(), srv.URL, langsmith.SandboxBoxRunParams{
		Command: langsmith.F("echo hi"),
		CWD:     langsmith.F("/tmp"),
		Env:     langsmith.F(map[string]string{"A": "1"}),
	})
	require.NoError(t, err)

	assert.Equal(t, "/tmp", body["cwd"])
	assert.Nil(t, body["run_config"])
}

func TestSandboxGlobAndGrep(t *testing.T) {
	var globBody, grepBody map[string]any

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch r.URL.Path {
		case "/glob":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&globBody))
			_, _ = w.Write([]byte(`{"matches":[{"path":"/w/app.py","is_dir":false,"size_bytes":12,"modified_at":"2026-07-14T00:00:00Z"}],"truncated":true}`))
		case "/grep":
			require.NoError(t, json.NewDecoder(r.Body).Decode(&grepBody))
			_, _ = w.Write([]byte(`{"matches":[{"path":"/w/app.py","line":12,"text":"# TODO"}],"truncated":false}`))
		default:
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"))
	ctx := context.Background()

	globbed, err := client.Sandboxes.Boxes.GlobWithDataplaneURL(ctx, srv.URL, langsmith.SandboxGlobParams{
		Pattern: langsmith.F("**/*.py"),
		Path:    langsmith.F("/w"),
		Limit:   langsmith.F(int64(200)),
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"pattern": "**/*.py", "path": "/w", "limit": float64(200)}, globBody)
	assert.True(t, globbed.Truncated)
	require.Len(t, globbed.Matches, 1)
	assert.Equal(t, int64(12), globbed.Matches[0].SizeBytes)

	grepped, err := client.Sandboxes.Boxes.GrepWithDataplaneURL(ctx, srv.URL, langsmith.SandboxGrepParams{
		Pattern: langsmith.F("TODO"),
		Path:    langsmith.F("/w"),
		Glob:    langsmith.F("*.py"),
	})
	require.NoError(t, err)
	assert.Equal(t, map[string]any{"pattern": "TODO", "path": "/w", "glob": "*.py"}, grepBody)
	require.Len(t, grepped.Matches, 1)
	assert.Equal(t, int64(12), grepped.Matches[0].Line)
}

func TestSandboxReadFileRange(t *testing.T) {
	tests := []struct {
		name    string
		params  langsmith.SandboxReadRangeParams
		handler func(w http.ResponseWriter, r *http.Request)
		verify  func(t *testing.T, chunk *langsmith.SandboxFileChunk, sentRange string)
	}{
		{
			name:   "partial content reports its offsets",
			params: langsmith.SandboxReadRangeParams{Start: langsmith.F(int64(10)), End: langsmith.F(int64(13))},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Range", "bytes 10-13/100")
				w.Header().Set("ETag", `"abc"`)
				w.WriteHeader(http.StatusPartialContent)
				_, _ = w.Write([]byte("0123"))
			},
			verify: func(t *testing.T, chunk *langsmith.SandboxFileChunk, sentRange string) {
				assert.Equal(t, "bytes=10-13", sentRange)
				assert.True(t, chunk.Partial)
				assert.Equal(t, int64(10), chunk.Start)
				assert.Equal(t, int64(14), chunk.End())
				assert.Equal(t, int64(100), chunk.TotalBytes)
				assert.Equal(t, `"abc"`, chunk.ETag)
				assert.Equal(t, []byte("0123"), chunk.Content)
			},
		},
		{
			name:   "stale if-range reports a whole file from zero",
			params: langsmith.SandboxReadRangeParams{Start: langsmith.F(int64(10)), IfRange: `"old"`},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, `"old"`, r.Header.Get("If-Range"))
				_, _ = w.Write([]byte("whole"))
			},
			verify: func(t *testing.T, chunk *langsmith.SandboxFileChunk, sentRange string) {
				assert.Equal(t, "bytes=10-", sentRange)
				assert.False(t, chunk.Partial)
				assert.Equal(t, int64(0), chunk.Start)
				assert.Equal(t, int64(5), chunk.TotalBytes)
			},
		},
		{
			name:   "unchanged file returns no bytes",
			params: langsmith.SandboxReadRangeParams{Start: langsmith.F(int64(0)), IfNoneMatch: `"same"`},
			handler: func(w http.ResponseWriter, r *http.Request) {
				assert.Equal(t, `"same"`, r.Header.Get("If-None-Match"))
				w.Header().Set("ETag", `"same"`)
				w.WriteHeader(http.StatusNotModified)
			},
			verify: func(t *testing.T, chunk *langsmith.SandboxFileChunk, sentRange string) {
				assert.True(t, chunk.Unchanged)
				assert.Empty(t, chunk.Content)
			},
		},
		{
			name:   "suffix range",
			params: langsmith.SandboxReadRangeParams{SuffixBytes: langsmith.F(int64(4))},
			handler: func(w http.ResponseWriter, r *http.Request) {
				w.Header().Set("Content-Range", "bytes 96-99/100")
				w.WriteHeader(http.StatusPartialContent)
				_, _ = w.Write([]byte("tail"))
			},
			verify: func(t *testing.T, chunk *langsmith.SandboxFileChunk, sentRange string) {
				assert.Equal(t, "bytes=-4", sentRange)
				assert.Equal(t, int64(96), chunk.Start)
			},
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			var sentRange string
			srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				sentRange = r.Header.Get("Range")
				tc.handler(w, r)
			}))
			defer srv.Close()

			client := langsmith.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"))
			chunk, err := client.Sandboxes.Boxes.ReadFileRangeWithDataplaneURL(context.Background(), srv.URL, "/big.bin", tc.params)
			require.NoError(t, err)
			tc.verify(t, chunk, sentRange)
		})
	}
}

func TestSandboxReadFileRangeRejectsInvalidRanges(t *testing.T) {
	client := langsmith.NewClient(option.WithBaseURL("http://example.invalid"), option.WithAPIKey("k"))

	for name, params := range map[string]langsmith.SandboxReadRangeParams{
		"no range":          {},
		"suffix with start": {Start: langsmith.F(int64(1)), SuffixBytes: langsmith.F(int64(2))},
		"end before start":  {Start: langsmith.F(int64(5)), End: langsmith.F(int64(1))},
		"negative start":    {Start: langsmith.F(int64(-1))},
		"zero suffix":       {SuffixBytes: langsmith.F(int64(0))},
	} {
		t.Run(name, func(t *testing.T) {
			_, err := client.Sandboxes.Boxes.ReadFileRangeWithDataplaneURL(context.Background(), "http://example.invalid", "/big.bin", params)
			require.Error(t, err)
		})
	}
}

func TestSandboxReadFileRangePastEOF(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Range", "bytes */100")
		w.WriteHeader(http.StatusRequestedRangeNotSatisfiable)
		_, _ = w.Write([]byte("invalid range"))
	}))
	defer srv.Close()

	client := langsmith.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"))
	_, err := client.Sandboxes.Boxes.ReadFileRangeWithDataplaneURL(context.Background(), srv.URL, "/big.bin", langsmith.SandboxReadRangeParams{
		Start: langsmith.F(int64(500)),
	})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "past the end of the file")
}

func TestSandboxStatFile(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, http.MethodHead, r.Method)
		assert.Equal(t, "/big.bin", r.URL.Query().Get("path"))
		w.Header().Set("ETag", `"abc"`)
		w.Header().Set("Content-Type", "application/octet-stream")
		w.Header().Set("Content-Length", "100")
		w.WriteHeader(http.StatusOK)
	}))
	defer srv.Close()

	client := langsmith.NewClient(option.WithBaseURL(srv.URL), option.WithAPIKey("k"))
	stat, err := client.Sandboxes.Boxes.StatFileWithDataplaneURL(context.Background(), srv.URL, "/big.bin")
	require.NoError(t, err)
	assert.Equal(t, int64(100), stat.SizeBytes)
	assert.Equal(t, `"abc"`, stat.ETag)
	assert.Equal(t, "application/octet-stream", stat.ContentType)
}
