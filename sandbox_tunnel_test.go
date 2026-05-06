package langsmith_test

import (
	"context"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/langchain-ai/langsmith-go"
	"github.com/langchain-ai/langsmith-go/option"
	"golang.org/x/net/websocket"
)

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
