package mockllm

import (
	"io"
	"net/http"
	"strings"
	"testing"
	"testing/synctest"
	"time"
)

func TestJitteredDelay_Range(t *testing.T) {
	synctest.Run(func() {
		delay := JitteredDelay(10*time.Millisecond, 50*time.Millisecond)
		for i := 0; i < 100; i++ {
			d := delay()
			if d < 10*time.Millisecond || d >= 50*time.Millisecond {
				t.Errorf("delay %v out of range [10ms, 50ms)", d)
			}
		}
	})
}

func TestDefaultStreamDelay_NonZero(t *testing.T) {
	synctest.Run(func() {
		delay := DefaultStreamDelay()
		if delay == nil {
			t.Fatal("DefaultStreamDelay should not return nil")
		}
		d := delay()
		if d < 30*time.Millisecond || d >= 120*time.Millisecond {
			t.Errorf("default delay %v out of expected range [30ms, 120ms)", d)
		}
	})
}

func TestNoDelay_Nil(t *testing.T) {
	if NoDelay() != nil {
		t.Error("NoDelay should return nil")
	}
}

func TestElizaHandler_DefaultJitter(t *testing.T) {
	synctest.Run(func() {
		h := ElizaHandler()
		resp := h(Request{
			Model:    "test",
			Messages: []Message{{Role: "user", Content: "hello"}},
		})
		if resp.StreamDelay == nil {
			t.Error("Eliza should have StreamDelay set by default")
		}
	})
}

func TestElizaHandler_RespondImmediately(t *testing.T) {
	synctest.Run(func() {
		h := ElizaHandler()

		// First: jitter enabled by default
		resp1 := h(Request{
			Model:    "test",
			Messages: []Message{{Role: "user", Content: "hello"}},
		})
		if resp1.StreamDelay == nil {
			t.Error("expected jitter enabled by default")
		}

		// Disable jitter
		resp2 := h(Request{
			Model:    "test",
			Messages: []Message{{Role: "user", Content: "Please respond immediately"}},
		})
		if resp2.StreamDelay != nil {
			t.Error("respond immediately response should not have delay")
		}

		// Subsequent calls: jitter stays disabled
		resp3 := h(Request{
			Model:    "test",
			Messages: []Message{{Role: "user", Content: "hello again"}},
		})
		if resp3.StreamDelay != nil {
			t.Error("jitter should stay disabled after 'respond immediately'")
		}
	})
}

// TestStreamDelay_UsedInStreaming verifies that StreamDelay is actually called
// during SSE streaming. This test uses real I/O (not synctest) because
// httptest.Server requires real goroutines.
func TestStreamDelay_UsedInStreaming(t *testing.T) {
	var delays int
	countingDelay := func() time.Duration {
		delays++
		return 0 // don't actually sleep in tests
	}

	h := func(req Request) Response {
		return Response{
			Content:      "one two three four five",
			InputTokens:  5,
			OutputTokens: 5,
			StopReason:   "end_turn",
			StreamDelay:  countingDelay,
		}
	}

	srv := NewCombinedServer(WithHandler(h))
	defer srv.Close()

	resp, err := http.Post(srv.URL()+"/v1/chat/completions",
		"application/json",
		strings.NewReader(`{"model":"test","messages":[{"role":"user","content":"hi"}],"stream":true}`))
	if err != nil {
		t.Fatalf("request: %v", err)
	}
	io.ReadAll(resp.Body)
	resp.Body.Close()

	if delays == 0 {
		t.Error("expected StreamDelay to be called at least once during streaming")
	}
	t.Logf("StreamDelay called %d times", delays)
}
