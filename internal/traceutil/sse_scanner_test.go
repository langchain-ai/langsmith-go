package traceutil

import (
	"io"
	"reflect"
	"strings"
	"testing"
)

func TestSSEScanner_ParsesAcrossWriteBoundaries(t *testing.T) {
	var got []map[string]any
	s := NewSSEScanner(func(c map[string]any) { got = append(got, c) })

	// Split a single data line across multiple Writes — the scanner must
	// hold partial bytes until a newline arrives.
	s.Feed([]byte("data: {\"a\""))
	s.Feed([]byte(":1}\n"))
	s.Feed([]byte("data: {\"b\":2}\ndata: {\"c\":3}\n"))

	want := []map[string]any{
		{"a": float64(1)},
		{"b": float64(2)},
		{"c": float64(3)},
	}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestSSEScanner_SkipsPreambleAndDone(t *testing.T) {
	var got []map[string]any
	s := NewSSEScanner(func(c map[string]any) { got = append(got, c) })

	// event:/id:/retry: lines, blank lines, and [DONE] should all be skipped.
	s.Feed([]byte("event: foo\n"))
	s.Feed([]byte("id: 1\n"))
	s.Feed([]byte("\n"))
	s.Feed([]byte("retry: 1000\n"))
	s.Feed([]byte("data: [DONE]\n"))
	s.Feed([]byte("data: {\"ok\":true}\n"))

	want := []map[string]any{{"ok": true}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestOnFirstSSEMatch_FiresOnceAndStopsParsing(t *testing.T) {
	src := &nopCloser{Reader: strings.NewReader(
		"data: {\"role\":\"assistant\"}\n" + // preamble — does not match
			"data: {\"content\":\"first\"}\n" + // matches → fires once
			"data: {\"content\":\"second\"}\n" + // would match again, but parsing must be stopped
			"data: {\"content\":\"third\"}\n",
	)}
	br := NewBufferedReader(src, nil)

	matches := 0
	fired := 0
	OnFirstSSEMatch(br,
		func(c map[string]any) bool {
			matches++
			_, ok := c["content"].(string)
			return ok
		},
		func() { fired++ },
	)

	if _, err := io.ReadAll(br); err != nil {
		t.Fatal(err)
	}

	if fired != 1 {
		t.Errorf("expected fire called once, got %d", fired)
	}
	// Only the preamble + the first matching chunk should have been parsed;
	// after the match, onBytes is detached and no further chunks are inspected.
	if matches != 2 {
		t.Errorf("expected scanner to stop after first match (2 isMatch calls), got %d", matches)
	}
}

func TestSSEScanner_SkipsInvalidJSON(t *testing.T) {
	var got []map[string]any
	s := NewSSEScanner(func(c map[string]any) { got = append(got, c) })

	// Non-JSON data lines (e.g. an error body that happens to start with
	// "data:") must not panic or fire the callback.
	s.Feed([]byte("data: not json\n"))
	s.Feed([]byte("data: {\"valid\":true}\n"))

	want := []map[string]any{{"valid": true}}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}
