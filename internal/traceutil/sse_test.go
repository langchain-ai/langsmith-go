package traceutil

import (
	"strings"
	"testing"
)

func TestParseSSEChunks_Basic(t *testing.T) {
	input := "data: {\"id\":\"1\",\"text\":\"hello\"}\ndata: {\"id\":\"2\",\"text\":\"world\"}\ndata: [DONE]\n"
	chunks, err := ParseSSEChunks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
	if chunks[0]["id"] != "1" || chunks[0]["text"] != "hello" {
		t.Errorf("chunk 0: %v", chunks[0])
	}
	if chunks[1]["id"] != "2" || chunks[1]["text"] != "world" {
		t.Errorf("chunk 1: %v", chunks[1])
	}
}

func TestParseSSEChunks_SkipsNonDataLines(t *testing.T) {
	input := "event: message\ndata: {\"ok\":true}\nid: 123\n\n"
	chunks, err := ParseSSEChunks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk, got %d", len(chunks))
	}
}

func TestParseSSEChunks_StopsAtDone(t *testing.T) {
	input := "data: {\"a\":1}\ndata: [DONE]\ndata: {\"b\":2}\n"
	chunks, err := ParseSSEChunks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk (before [DONE]), got %d", len(chunks))
	}
}

func TestParseSSEChunks_MalformedJSON(t *testing.T) {
	input := "data: {\"ok\":true}\ndata: not-json\ndata: {\"ok\":false}\n"
	chunks, err := ParseSSEChunks(strings.NewReader(input))
	if err == nil {
		t.Fatal("expected error for malformed JSON")
	}
	if len(chunks) != 1 {
		t.Fatalf("expected 1 chunk (parsed before malformed line), got %d", len(chunks))
	}
}

func TestParseSSEChunks_Empty(t *testing.T) {
	chunks, err := ParseSSEChunks(strings.NewReader(""))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 0 {
		t.Fatalf("expected 0 chunks, got %d", len(chunks))
	}
}

func TestParseSSEChunks_NoDone(t *testing.T) {
	input := "data: {\"type\":\"message_start\"}\ndata: {\"type\":\"message_stop\"}\n"
	chunks, err := ParseSSEChunks(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(chunks) != 2 {
		t.Fatalf("expected 2 chunks, got %d", len(chunks))
	}
}
