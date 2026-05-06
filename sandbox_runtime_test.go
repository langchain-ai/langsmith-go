package langsmith_test

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
	"golang.org/x/net/websocket"
)

func TestSandboxBoxRunWithDataplaneURL(t *testing.T) {
	var gotPayload map[string]any
	var gotAPIKey string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost || r.URL.Path != "/execute" {
			t.Fatalf("unexpected request: %s %s", r.Method, r.URL.Path)
		}
		gotAPIKey = r.Header.Get("X-API-Key")
		if err := json.NewDecoder(r.Body).Decode(&gotPayload); err != nil {
			t.Fatalf("decode request body: %v", err)
		}
		w.Header().Set("Content-Type", "application/json")
		_ = json.NewEncoder(w).Encode(map[string]any{
			"stdout":    "hello\n",
			"stderr":    "",
			"exit_code": 0,
		})
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	res, err := client.Sandboxes.Boxes.RunWithDataplaneURL(context.Background(), srv.URL, langsmith.SandboxBoxRunParams{
		Command: langsmith.String("echo hello"),
		Timeout: langsmith.Int(12),
		Env: langsmith.F(map[string]string{
			"NAME": "value",
		}),
		CWD:   langsmith.String("/tmp"),
		Shell: langsmith.String("/bin/sh"),
	})
	if err != nil {
		t.Fatalf("RunWithDataplaneURL returned error: %v", err)
	}
	if res.Stdout != "hello\n" || res.Stderr != "" || res.ExitCode != 0 {
		t.Fatalf("unexpected result: %#v", res)
	}
	if !res.Success() {
		t.Fatal("expected result to be successful")
	}
	if gotAPIKey != "test-api-key" {
		t.Fatalf("expected API key header, got %q", gotAPIKey)
	}
	if gotPayload["command"] != "echo hello" {
		t.Fatalf("expected command payload, got %#v", gotPayload["command"])
	}
	if gotPayload["timeout"] != float64(12) {
		t.Fatalf("expected timeout payload, got %#v", gotPayload["timeout"])
	}
	if gotPayload["shell"] != "/bin/sh" {
		t.Fatalf("expected shell payload, got %#v", gotPayload["shell"])
	}
	if gotPayload["cwd"] != "/tmp" {
		t.Fatalf("expected cwd payload, got %#v", gotPayload["cwd"])
	}
	env, ok := gotPayload["env"].(map[string]any)
	if !ok || env["NAME"] != "value" {
		t.Fatalf("expected env payload, got %#v", gotPayload["env"])
	}
}

func TestSandboxBoxRunFetchesDataplaneURL(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/test-box":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"name":          "test-box",
				"status":        "ready",
				"dataplane_url": srvURL(r),
			})
		case r.Method == http.MethodPost && r.URL.Path == "/execute":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"stdout":    "ok",
				"stderr":    "",
				"exit_code": 0,
			})
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	res, err := client.Sandboxes.Boxes.Run(context.Background(), "test-box", langsmith.SandboxBoxRunParams{
		Command: langsmith.String("true"),
	})
	if err != nil {
		t.Fatalf("Run returned error: %v", err)
	}
	if res.Stdout != "ok" {
		t.Fatalf("expected stdout ok, got %q", res.Stdout)
	}
}

func TestSandboxCommandHandleWebSocketLifecycle(t *testing.T) {
	executePayload := make(chan map[string]any, 1)
	inputPayload := make(chan map[string]any, 1)
	killPayload := make(chan map[string]any, 1)

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive execute: %v", err)
			return
		}
		executePayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"started","command_id":"cmd-123","pid":42}`); err != nil {
			t.Errorf("send started: %v", err)
			return
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive input: %v", err)
			return
		}
		inputPayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"stdout","data":"ready\n","offset":0}`); err != nil {
			t.Errorf("send stdout: %v", err)
			return
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive kill: %v", err)
			return
		}
		killPayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"exit","exit_code":137}`); err != nil {
			t.Errorf("send exit: %v", err)
			return
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	handle, err := client.Sandboxes.Boxes.StartCommandWithDataplaneURL(context.Background(), srv.URL, langsmith.SandboxCommandStartParams{
		Command:            langsmith.String("python3 -u -i"),
		TimeoutSeconds:     langsmith.Int(90),
		Shell:              langsmith.String("/bin/sh"),
		IdleTimeoutSeconds: langsmith.Int(-1),
		TTLSeconds:         langsmith.Int(120),
		Pty:                langsmith.Bool(true),
	})
	if err != nil {
		t.Fatalf("StartCommandWithDataplaneURL returned error: %v", err)
	}
	if handle.CommandID != "cmd-123" {
		t.Fatalf("expected command id, got %q", handle.CommandID)
	}
	if handle.PID != 42 {
		t.Fatalf("expected pid 42, got %d", handle.PID)
	}

	execute := <-executePayload
	if execute["type"] != "execute" || execute["command"] != "python3 -u -i" {
		t.Fatalf("unexpected execute payload: %#v", execute)
	}
	if execute["timeout_seconds"] != float64(90) {
		t.Fatalf("unexpected timeout_seconds: %#v", execute["timeout_seconds"])
	}
	if execute["shell"] != "/bin/sh" {
		t.Fatalf("unexpected shell: %#v", execute["shell"])
	}
	if execute["idle_timeout_seconds"] != float64(-1) {
		t.Fatalf("unexpected idle_timeout_seconds: %#v", execute["idle_timeout_seconds"])
	}
	if execute["ttl_seconds"] != float64(120) {
		t.Fatalf("unexpected ttl_seconds: %#v", execute["ttl_seconds"])
	}
	if execute["pty"] != true {
		t.Fatalf("unexpected pty: %#v", execute["pty"])
	}

	if err := handle.SendInput("print('hi')\n"); err != nil {
		t.Fatalf("SendInput returned error: %v", err)
	}
	input := <-inputPayload
	if input["type"] != "input" || input["data"] != "print('hi')\n" {
		t.Fatalf("unexpected input payload: %#v", input)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chunk, ok, err := handle.Next(ctx)
	if err != nil {
		t.Fatalf("Next returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected output chunk")
	}
	if chunk.Stream != "stdout" || chunk.Data != "ready\n" || chunk.Offset != 0 {
		t.Fatalf("unexpected chunk: %#v", chunk)
	}
	if handle.LastStdoutOffset() != int64(len("ready\n")) {
		t.Fatalf("unexpected stdout offset: %d", handle.LastStdoutOffset())
	}

	if err := handle.Kill(); err != nil {
		t.Fatalf("Kill returned error: %v", err)
	}
	kill := <-killPayload
	if kill["type"] != "kill" {
		t.Fatalf("unexpected kill payload: %#v", kill)
	}

	result, err := handle.Result(ctx)
	if err != nil {
		t.Fatalf("Result returned error: %v", err)
	}
	if result.Stdout != "ready\n" || result.ExitCode != 137 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestSandboxReconnectCommandWithDataplaneURL(t *testing.T) {
	reconnectPayload := make(chan map[string]any, 1)

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive reconnect: %v", err)
			return
		}
		reconnectPayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"stderr","data":"warn\n","offset":3}`); err != nil {
			t.Errorf("send stderr: %v", err)
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"exit","exit_code":0}`); err != nil {
			t.Errorf("send exit: %v", err)
			return
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	handle, err := client.Sandboxes.Boxes.ReconnectCommandWithDataplaneURL(context.Background(), srv.URL, "cmd-123", langsmith.SandboxCommandReconnectParams{
		StdoutOffset: langsmith.Int(7),
		StderrOffset: langsmith.Int(3),
	})
	if err != nil {
		t.Fatalf("ReconnectCommandWithDataplaneURL returned error: %v", err)
	}
	if handle.CommandID != "cmd-123" {
		t.Fatalf("expected command id, got %q", handle.CommandID)
	}

	reconnect := <-reconnectPayload
	if reconnect["type"] != "reconnect" || reconnect["command_id"] != "cmd-123" {
		t.Fatalf("unexpected reconnect payload: %#v", reconnect)
	}
	if reconnect["stdout_offset"] != float64(7) {
		t.Fatalf("unexpected stdout offset: %#v", reconnect["stdout_offset"])
	}
	if reconnect["stderr_offset"] != float64(3) {
		t.Fatalf("unexpected stderr offset: %#v", reconnect["stderr_offset"])
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chunk, ok, err := handle.Next(ctx)
	if err != nil {
		t.Fatalf("Next returned error: %v", err)
	}
	if !ok {
		t.Fatal("expected output chunk")
	}
	if chunk.Stream != "stderr" || chunk.Data != "warn\n" || chunk.Offset != 3 {
		t.Fatalf("unexpected chunk: %#v", chunk)
	}
	if handle.LastStderrOffset() != int64(8) {
		t.Fatalf("unexpected stderr offset: %d", handle.LastStderrOffset())
	}

	result, err := handle.Result(ctx)
	if err != nil {
		t.Fatalf("Result returned error: %v", err)
	}
	if result.Stderr != "warn\n" || result.ExitCode != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
}

func TestSandboxFileReadWriteWithDataplaneURL(t *testing.T) {
	var uploadedPath string
	var uploadedContent string
	var uploadAPIKey string
	var readAPIKey string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/upload":
			uploadAPIKey = r.Header.Get("X-API-Key")
			uploadedPath = r.URL.Query().Get("path")
			file, _, err := r.FormFile("file")
			if err != nil {
				t.Fatalf("read upload file: %v", err)
			}
			defer file.Close()
			data, err := io.ReadAll(file)
			if err != nil {
				t.Fatalf("read upload content: %v", err)
			}
			uploadedContent = string(data)
			w.WriteHeader(http.StatusNoContent)
		case r.Method == http.MethodGet && r.URL.Path == "/download":
			readAPIKey = r.Header.Get("X-API-Key")
			if r.URL.Query().Get("path") != "/app/test.txt" {
				t.Fatalf("unexpected read path: %q", r.URL.Query().Get("path"))
			}
			_, _ = w.Write([]byte("file contents"))
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	if err := client.Sandboxes.Boxes.WriteFileWithDataplaneURL(context.Background(), srv.URL, "/app/test.txt", []byte("hello world")); err != nil {
		t.Fatalf("WriteFileWithDataplaneURL returned error: %v", err)
	}
	if uploadedPath != "/app/test.txt" {
		t.Fatalf("unexpected upload path: %q", uploadedPath)
	}
	if uploadedContent != "hello world" {
		t.Fatalf("unexpected upload content: %q", uploadedContent)
	}
	if uploadAPIKey != "test-api-key" {
		t.Fatalf("expected upload API key header, got %q", uploadAPIKey)
	}

	content, err := client.Sandboxes.Boxes.ReadFileWithDataplaneURL(context.Background(), srv.URL, "/app/test.txt")
	if err != nil {
		t.Fatalf("ReadFileWithDataplaneURL returned error: %v", err)
	}
	if string(content) != "file contents" {
		t.Fatalf("unexpected read content: %q", string(content))
	}
	if readAPIKey != "test-api-key" {
		t.Fatalf("expected read API key header, got %q", readAPIKey)
	}
}

func TestSandboxWaitHelpers(t *testing.T) {
	var sandboxStatusCalls atomic.Int64
	var snapshotGetCalls atomic.Int64

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/test-box/status":
			if sandboxStatusCalls.Add(1) == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{"status": "provisioning"})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{"status": "ready"})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/boxes/test-box":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":            "box-id",
				"name":          "test-box",
				"status":        "ready",
				"dataplane_url": "http://dataplane.test",
			})
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/snapshots":
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "snap-1",
				"name":              "snap",
				"status":            "building",
				"fs_capacity_bytes": 1024,
			})
		case r.Method == http.MethodGet && r.URL.Path == "/v2/sandboxes/snapshots/snap-1":
			if snapshotGetCalls.Add(1) == 1 {
				_ = json.NewEncoder(w).Encode(map[string]any{
					"id":                "snap-1",
					"name":              "snap",
					"status":            "building",
					"fs_capacity_bytes": 1024,
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"id":                "snap-1",
				"name":              "snap",
				"status":            "ready",
				"fs_capacity_bytes": 1024,
			})
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)

	sandbox, err := client.Sandboxes.Boxes.WaitSandbox(context.Background(), "test-box", langsmith.SandboxWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("WaitSandbox returned error: %v", err)
	}
	if sandbox.Name != "test-box" || sandbox.DataplaneURL != "http://dataplane.test" {
		t.Fatalf("unexpected sandbox: %#v", sandbox)
	}
	if sandboxStatusCalls.Load() != 2 {
		t.Fatalf("expected 2 sandbox status calls, got %d", sandboxStatusCalls.Load())
	}

	snapshot, err := client.Sandboxes.Snapshots.NewAndWait(context.Background(), langsmith.SandboxSnapshotNewParams{
		Name:            langsmith.String("snap"),
		DockerImage:     langsmith.String("python:3.12-slim"),
		FsCapacityBytes: langsmith.Int(1024),
	}, langsmith.SnapshotWaitParams{
		Timeout:      time.Second,
		PollInterval: time.Millisecond,
	})
	if err != nil {
		t.Fatalf("NewAndWait returned error: %v", err)
	}
	if snapshot.ID != "snap-1" || snapshot.Status != "ready" {
		t.Fatalf("unexpected snapshot: %#v", snapshot)
	}
	if snapshotGetCalls.Load() != 2 {
		t.Fatalf("expected 2 snapshot get calls, got %d", snapshotGetCalls.Load())
	}
}

func TestSandboxServiceURLRefreshAndRequest(t *testing.T) {
	var serviceURLCalls atomic.Int64
	var firstPayload map[string]any
	var serviceAuthHeader string
	var serviceCustomHeader string

	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		switch {
		case r.Method == http.MethodPost && r.URL.Path == "/v2/sandboxes/boxes/test-box/service-url":
			call := serviceURLCalls.Add(1)
			var payload map[string]any
			if err := json.NewDecoder(r.Body).Decode(&payload); err != nil {
				t.Fatalf("decode service URL request: %v", err)
			}
			if call == 1 {
				firstPayload = payload
				_ = json.NewEncoder(w).Encode(map[string]any{
					"browser_url": srvURL(r) + "/browser",
					"service_url": srvURL(r) + "/service",
					"token":       "token-1",
					"expires_at":  time.Now().Add(time.Second).UTC().Format(time.RFC3339),
				})
				return
			}
			_ = json.NewEncoder(w).Encode(map[string]any{
				"browser_url": srvURL(r) + "/browser",
				"service_url": srvURL(r) + "/service",
				"token":       "token-2",
				"expires_at":  time.Now().Add(time.Hour).UTC().Format(time.RFC3339),
			})
		case r.Method == http.MethodGet && r.URL.Path == "/service/path":
			serviceAuthHeader = r.Header.Get("X-Langsmith-Sandbox-Service-Token")
			serviceCustomHeader = r.Header.Get("X-Test")
			_, _ = w.Write([]byte(`{"ok":true}`))
		default:
			http.Error(w, "unexpected "+r.Method+" "+r.URL.Path, http.StatusInternalServerError)
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL(srv.URL),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	service, err := client.Sandboxes.Boxes.Service(context.Background(), "test-box", langsmith.SandboxBoxGenerateServiceURLParams{
		Port: langsmith.Int(3000),
	})
	if err != nil {
		t.Fatalf("Service returned error: %v", err)
	}
	if firstPayload["port"] != float64(3000) {
		t.Fatalf("unexpected service port payload: %#v", firstPayload["port"])
	}
	if firstPayload["expires_in_seconds"] != float64(600) {
		t.Fatalf("expected default service URL TTL, got %#v", firstPayload["expires_in_seconds"])
	}

	token, err := service.Token(context.Background())
	if err != nil {
		t.Fatalf("Token returned error: %v", err)
	}
	if token != "token-2" {
		t.Fatalf("expected refreshed token, got %q", token)
	}
	if serviceURLCalls.Load() != 2 {
		t.Fatalf("expected 2 service URL calls, got %d", serviceURLCalls.Load())
	}

	resp, err := service.Get(context.Background(), "/path", http.Header{"X-Test": []string{"yes"}})
	if err != nil {
		t.Fatalf("service Get returned error: %v", err)
	}
	_ = resp.Body.Close()
	if serviceAuthHeader != "token-2" {
		t.Fatalf("expected service auth token header, got %q", serviceAuthHeader)
	}
	if serviceCustomHeader != "yes" {
		t.Fatalf("expected custom service header, got %q", serviceCustomHeader)
	}
}

func TestSandboxRunWithCallbacks(t *testing.T) {
	var stdout strings.Builder
	var stderr strings.Builder

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive execute: %v", err)
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"started","command_id":"cmd-123","pid":42}`); err != nil {
			t.Errorf("send started: %v", err)
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"stdout","data":"out","offset":0}`); err != nil {
			t.Errorf("send stdout: %v", err)
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"stderr","data":"err","offset":0}`); err != nil {
			t.Errorf("send stderr: %v", err)
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"exit","exit_code":0}`); err != nil {
			t.Errorf("send exit: %v", err)
			return
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	result, err := client.Sandboxes.Boxes.RunWithDataplaneURLAndCallbacks(context.Background(), srv.URL, langsmith.SandboxCommandStartParams{
		Command: langsmith.String("make test"),
	}, langsmith.SandboxCommandCallbacks{
		OnStdout: func(data string) { stdout.WriteString(data) },
		OnStderr: func(data string) { stderr.WriteString(data) },
	})
	if err != nil {
		t.Fatalf("RunWithDataplaneURLAndCallbacks returned error: %v", err)
	}
	if result.Stdout != "out" || result.Stderr != "err" || result.ExitCode != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
	if stdout.String() != "out" || stderr.String() != "err" {
		t.Fatalf("unexpected callbacks: stdout=%q stderr=%q", stdout.String(), stderr.String())
	}
}

func TestSandboxCommandHandleAutoReconnect(t *testing.T) {
	reconnectPayload := make(chan map[string]any, 1)
	var connections atomic.Int64

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		conn := connections.Add(1)
		var msg map[string]any
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive message: %v", err)
			return
		}
		if conn == 1 {
			if err := websocket.Message.Send(ws, `{"type":"started","command_id":"cmd-123","pid":42}`); err != nil {
				t.Errorf("send started: %v", err)
				return
			}
			if err := websocket.Message.Send(ws, `{"type":"stdout","data":"one","offset":0}`); err != nil {
				t.Errorf("send stdout: %v", err)
				return
			}
			_ = ws.Close()
			return
		}
		reconnectPayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"stdout","data":"two","offset":3}`); err != nil {
			t.Errorf("send reconnected stdout: %v", err)
			return
		}
		if err := websocket.Message.Send(ws, `{"type":"exit","exit_code":0}`); err != nil {
			t.Errorf("send exit: %v", err)
			return
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	handle, err := client.Sandboxes.Boxes.StartCommandWithDataplaneURL(context.Background(), srv.URL, langsmith.SandboxCommandStartParams{
		Command: langsmith.String("printf onetwo"),
	})
	if err != nil {
		t.Fatalf("StartCommandWithDataplaneURL returned error: %v", err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	result, err := handle.Result(ctx)
	if err != nil {
		t.Fatalf("Result returned error: %v", err)
	}
	if result.Stdout != "onetwo" || result.ExitCode != 0 {
		t.Fatalf("unexpected result: %#v", result)
	}
	reconnect := <-reconnectPayload
	if reconnect["type"] != "reconnect" || reconnect["command_id"] != "cmd-123" {
		t.Fatalf("unexpected reconnect payload: %#v", reconnect)
	}
	if reconnect["stdout_offset"] != float64(3) || reconnect["stderr_offset"] != float64(0) {
		t.Fatalf("unexpected reconnect offsets: %#v", reconnect)
	}
	if connections.Load() != 2 {
		t.Fatalf("expected 2 websocket connections, got %d", connections.Load())
	}
}

func TestSandboxTunnelWithDataplaneURL(t *testing.T) {
	serverDone := make(chan error, 1)

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		reader := &testWSFrameReader{ws: ws}

		msgType, flags, streamID, payload, err := reader.readFrame()
		if err != nil {
			serverDone <- err
			return
		}
		if msgType != 1 || flags != 1 || streamID != 1 || len(payload) != 0 {
			serverDone <- fmt.Errorf("unexpected open stream frame: type=%d flags=%d stream=%d payload=%v", msgType, flags, streamID, payload)
			return
		}

		msgType, flags, streamID, payload, err = reader.readFrame()
		if err != nil {
			serverDone <- err
			return
		}
		if msgType != 0 || flags != 0 || streamID != 1 || !equalBytes(payload, []byte{1, 31, 144}) {
			serverDone <- fmt.Errorf("unexpected tunnel connect frame: type=%d flags=%d stream=%d payload=%v", msgType, flags, streamID, payload)
			return
		}
		if err := sendTestYamuxFrame(ws, 0, 0, 1, []byte{0}); err != nil {
			serverDone <- err
			return
		}

		msgType, flags, streamID, payload, err = reader.readFrame()
		if err != nil {
			serverDone <- err
			return
		}
		if msgType != 0 || flags != 0 || streamID != 1 || string(payload) != "ping" {
			serverDone <- fmt.Errorf("unexpected tunnel data frame: type=%d flags=%d stream=%d payload=%q", msgType, flags, streamID, string(payload))
			return
		}
		if err := sendTestYamuxFrame(ws, 0, 0, 1, []byte("pong")); err != nil {
			serverDone <- err
			return
		}
		serverDone <- nil
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	tunnel, err := client.Sandboxes.Boxes.TunnelWithDataplaneURL(context.Background(), srv.URL, 8080, langsmith.SandboxTunnelParams{LocalPort: 0})
	if err != nil {
		t.Fatalf("TunnelWithDataplaneURL returned error: %v", err)
	}
	defer tunnel.Close()

	conn, err := net.DialTimeout("tcp", fmt.Sprintf("127.0.0.1:%d", tunnel.LocalPort), time.Second)
	if err != nil {
		t.Fatalf("dial local tunnel: %v", err)
	}
	defer conn.Close()
	_ = conn.SetDeadline(time.Now().Add(time.Second))
	if _, err := conn.Write([]byte("ping")); err != nil {
		t.Fatalf("write tunnel data: %v", err)
	}
	buf := make([]byte, 4)
	if _, err := io.ReadFull(conn, buf); err != nil {
		t.Fatalf("read tunnel data: %v", err)
	}
	if string(buf) != "pong" {
		t.Fatalf("unexpected tunnel response: %q", string(buf))
	}

	select {
	case err := <-serverDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for tunnel server")
	}
}

func TestSandboxOpenTunnelStreamWithDataplaneURL(t *testing.T) {
	serverDone := make(chan error, 1)
	remotePort := 9000

	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		reader := &testWSFrameReader{ws: ws}
		msgType, flags, streamID, payload, err := reader.readFrame()
		if err != nil {
			serverDone <- err
			return
		}
		if msgType != 1 || flags != 1 || streamID != 1 || len(payload) != 0 {
			serverDone <- fmt.Errorf("unexpected open stream frame: type=%d flags=%d stream=%d payload=%v", msgType, flags, streamID, payload)
			return
		}

		_, _, streamID, payload, err = reader.readFrame()
		if err != nil {
			serverDone <- err
			return
		}
		expectedHeader := []byte{byte(langsmith.SandboxTunnelProtocolVersion), byte(remotePort >> 8), byte(remotePort)}
		if !equalBytes(payload, expectedHeader) {
			serverDone <- fmt.Errorf("unexpected connect header: %v", payload)
			return
		}
		if err := sendTestYamuxFrame(ws, 0, 0, streamID, []byte{byte(langsmith.SandboxTunnelStatusOK)}); err != nil {
			serverDone <- err
			return
		}

		_, _, streamID, payload, err = reader.readFrame()
		if err != nil {
			serverDone <- err
			return
		}
		if string(payload) != "ping" {
			serverDone <- fmt.Errorf("unexpected stream payload: %q", string(payload))
			return
		}
		if err := sendTestYamuxFrame(ws, 0, 0, streamID, []byte("pong")); err != nil {
			serverDone <- err
			return
		}
		serverDone <- nil
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	stream, err := client.Sandboxes.Boxes.OpenTunnelStreamWithDataplaneURL(context.Background(), srv.URL, remotePort)
	if err != nil {
		t.Fatalf("OpenTunnelStreamWithDataplaneURL returned error: %v", err)
	}
	defer stream.Close()
	if stream.RemotePort != remotePort {
		t.Fatalf("expected remote port %d, got %d", remotePort, stream.RemotePort)
	}
	if _, err := stream.Write([]byte("ping")); err != nil {
		t.Fatalf("write stream: %v", err)
	}
	buf := make([]byte, 4)
	if _, err := io.ReadFull(stream, buf); err != nil {
		t.Fatalf("read stream: %v", err)
	}
	if string(buf) != "pong" {
		t.Fatalf("unexpected stream response: %q", string(buf))
	}

	select {
	case err := <-serverDone:
		if err != nil {
			t.Fatal(err)
		}
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for tunnel server")
	}
}

func TestSandboxOpenTunnelStreamStatusError(t *testing.T) {
	srv := httptest.NewServer(websocket.Handler(func(ws *websocket.Conn) {
		reader := &testWSFrameReader{ws: ws}
		_, _, streamID, _, err := reader.readFrame()
		if err != nil {
			t.Errorf("read open frame: %v", err)
			return
		}
		if _, _, _, _, err := reader.readFrame(); err != nil {
			t.Errorf("read connect frame: %v", err)
			return
		}
		if err := sendTestYamuxFrame(ws, 0, 0, streamID, []byte{byte(langsmith.SandboxTunnelStatusDialFailed)}); err != nil {
			t.Errorf("send status: %v", err)
			return
		}
	}))
	defer srv.Close()

	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	_, err := client.Sandboxes.Boxes.OpenTunnelStreamWithDataplaneURL(context.Background(), srv.URL, 5432)
	if err == nil {
		t.Fatal("expected status error")
	}
	var statusErr *langsmith.SandboxTunnelStatusError
	if !errors.As(err, &statusErr) {
		t.Fatalf("expected SandboxTunnelStatusError, got %T: %v", err, err)
	}
	if statusErr.Status != langsmith.SandboxTunnelStatusDialFailed || statusErr.RemotePort != 5432 {
		t.Fatalf("unexpected status error: %#v", statusErr)
	}
	if statusErr.Status.Reason() != "dial failed" {
		t.Fatalf("unexpected status reason: %q", statusErr.Status.Reason())
	}
}

type testWSFrameReader struct {
	ws  *websocket.Conn
	buf []byte
}

func (r *testWSFrameReader) read(n int) ([]byte, error) {
	for len(r.buf) < n {
		var msg []byte
		if err := websocket.Message.Receive(r.ws, &msg); err != nil {
			return nil, err
		}
		r.buf = append(r.buf, msg...)
	}
	out := append([]byte(nil), r.buf[:n]...)
	r.buf = r.buf[n:]
	return out, nil
}

func (r *testWSFrameReader) readFrame() (msgType byte, flags uint16, streamID uint32, payload []byte, err error) {
	header, err := r.read(12)
	if err != nil {
		return 0, 0, 0, nil, err
	}
	length := binary.BigEndian.Uint32(header[8:12])
	if length > 0 {
		payload, err = r.read(int(length))
		if err != nil {
			return 0, 0, 0, nil, err
		}
	}
	return header[1], binary.BigEndian.Uint16(header[2:4]), binary.BigEndian.Uint32(header[4:8]), payload, nil
}

func sendTestYamuxFrame(ws *websocket.Conn, msgType byte, flags uint16, streamID uint32, payload []byte) error {
	frame := make([]byte, 12+len(payload))
	frame[1] = msgType
	binary.BigEndian.PutUint16(frame[2:4], flags)
	binary.BigEndian.PutUint32(frame[4:8], streamID)
	binary.BigEndian.PutUint32(frame[8:12], uint32(len(payload)))
	copy(frame[12:], payload)
	return websocket.Message.Send(ws, frame)
}

func equalBytes(a []byte, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range a {
		if a[i] != b[i] {
			return false
		}
	}
	return true
}

func srvURL(r *http.Request) string {
	scheme := "http"
	if r.TLS != nil {
		scheme = "https"
	}
	return scheme + "://" + r.Host
}
