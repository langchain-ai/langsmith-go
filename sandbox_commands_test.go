package langsmith_test

import (
	"context"
	"encoding/base64"
	"encoding/json"
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

func TestSandboxCommandConsoleControls(t *testing.T) {
	executePayload := make(chan map[string]any, 1)
	resizePayload := make(chan map[string]any, 1)
	agentResponsePayload := make(chan map[string]any, 1)
	agentClosePayload := make(chan map[string]any, 1)

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
			t.Errorf("receive resize: %v", err)
			return
		}
		resizePayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"stdout","data":"shell ready\n","offset":0}`); err != nil {
			t.Errorf("send stdout: %v", err)
			return
		}
		agentRequest := base64.StdEncoding.EncodeToString([]byte("agent-request"))
		if err := websocket.Message.Send(ws, `{"type":"ssh_agent_data","channel_id":"ch-1","data":"`+agentRequest+`"}`); err != nil {
			t.Errorf("send agent data: %v", err)
			return
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive agent response: %v", err)
			return
		}
		agentResponsePayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"ssh_agent_close","channel_id":"ch-1"}`); err != nil {
			t.Errorf("send agent close: %v", err)
			return
		}
		if err := websocket.JSON.Receive(ws, &msg); err != nil {
			t.Errorf("receive close response: %v", err)
			return
		}
		agentClosePayload <- msg
		if err := websocket.Message.Send(ws, `{"type":"exit","exit_code":0}`); err != nil {
			t.Errorf("send exit: %v", err)
			return
		}
	}))
	defer srv.Close()

	agentData := make(chan []byte, 1)
	agentClose := make(chan string, 1)
	client := langsmith.NewClient(
		option.WithBaseURL("http://control-plane.test"),
		option.WithAPIKey("test-api-key"),
		option.WithMaxRetries(0),
	)
	handle, err := client.Sandboxes.Boxes.StartCommandWithDataplaneURLAndCallbacks(context.Background(), srv.URL, langsmith.SandboxCommandStartParams{
		Command:         langsmith.String("/bin/bash"),
		Pty:             langsmith.Bool(true),
		SSHAgentForward: langsmith.Bool(true),
	}, langsmith.SandboxCommandCallbacks{
		OnSSHAgentData: func(channelID string, data []byte) {
			if channelID != "ch-1" {
				t.Errorf("unexpected channel ID: %s", channelID)
			}
			agentData <- data
		},
		OnSSHAgentClose: func(channelID string) {
			agentClose <- channelID
		},
	})
	if err != nil {
		t.Fatalf("StartCommandWithDataplaneURLAndCallbacks returned error: %v", err)
	}

	execute := <-executePayload
	if execute["type"] != "execute" || execute["command"] != "/bin/bash" {
		t.Fatalf("unexpected execute payload: %#v", execute)
	}
	if execute["pty"] != true {
		t.Fatalf("unexpected pty: %#v", execute["pty"])
	}
	if execute["ssh_agent_forward"] != true {
		t.Fatalf("unexpected ssh_agent_forward: %#v", execute["ssh_agent_forward"])
	}

	if err := handle.Resize(120, 40); err != nil {
		t.Fatalf("Resize returned error: %v", err)
	}
	resize := <-resizePayload
	if resize["type"] != "resize" || resize["cols"] != float64(120) || resize["rows"] != float64(40) {
		t.Fatalf("unexpected resize payload: %#v", resize)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	chunk, ok, err := handle.Next(ctx)
	if err != nil {
		t.Fatalf("Next returned error: %v", err)
	}
	if !ok || chunk.Data != "shell ready\n" {
		t.Fatalf("unexpected output chunk: %#v", chunk)
	}

	if got := <-agentData; string(got) != "agent-request" {
		t.Fatalf("unexpected agent data: %q", string(got))
	}
	if err := handle.SendSSHAgentData("ch-1", []byte("agent-response")); err != nil {
		t.Fatalf("SendSSHAgentData returned error: %v", err)
	}
	agentResponse := <-agentResponsePayload
	if agentResponse["type"] != "ssh_agent_data" || agentResponse["channel_id"] != "ch-1" {
		t.Fatalf("unexpected agent response payload: %#v", agentResponse)
	}
	decoded, err := base64.StdEncoding.DecodeString(agentResponse["data"].(string))
	if err != nil {
		t.Fatalf("decode agent response: %v", err)
	}
	if string(decoded) != "agent-response" {
		t.Fatalf("unexpected agent response data: %q", string(decoded))
	}

	if got := <-agentClose; got != "ch-1" {
		t.Fatalf("unexpected agent close channel: %q", got)
	}
	if err := handle.CloseSSHAgentChannel("ch-1"); err != nil {
		t.Fatalf("CloseSSHAgentChannel returned error: %v", err)
	}
	closePayload := <-agentClosePayload
	if closePayload["type"] != "ssh_agent_close" || closePayload["channel_id"] != "ch-1" {
		t.Fatalf("unexpected agent close payload: %#v", closePayload)
	}

	result, err := handle.Result(ctx)
	if err != nil {
		t.Fatalf("Result returned error: %v", err)
	}
	if result.Stdout != "shell ready\n" || result.ExitCode != 0 {
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
