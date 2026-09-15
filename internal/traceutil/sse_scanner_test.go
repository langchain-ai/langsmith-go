package traceutil

import (
	"bytes"
	"errors"
	"fmt"
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

func TestSSEScanner_SkipsPreambleAndStopsAtDone(t *testing.T) {
	var got []map[string]any
	s := NewSSEScanner(func(c map[string]any) { got = append(got, c) })

	// Ignore metadata lines, then stop at [DONE] even without a callback.
	s.Feed([]byte("event: foo\n"))
	s.Feed([]byte("id: 1\n"))
	s.Feed([]byte("\n"))
	s.Feed([]byte("retry: 1000\n"))
	s.Feed([]byte("data: [DONE]\n"))
	s.Feed([]byte("data: {\"ok\":true}\n"))

	var want []map[string]any
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestOnFirstMatch(t *testing.T) {
	consumed, matches, fired := 0, 0, 0
	feed := OnFirstMatch(func(map[string]any) { consumed++ },
		func(c map[string]any) bool { matches++; return c["content"] != nil },
		func() { fired++ })
	feed(map[string]any{"role": "assistant"})
	feed(map[string]any{"content": "first"})
	feed(map[string]any{"content": "second"})
	if consumed != 3 || matches != 2 || fired != 1 {
		t.Fatalf("consumed=%d matches=%d fired=%d", consumed, matches, fired)
	}
}

func TestSSEScanner_StopsAtInvalidJSON(t *testing.T) {
	var got []map[string]any
	s := NewSSEScanner(func(c map[string]any) { got = append(got, c) })

	// Non-JSON data lines (e.g. an error body that happens to start with
	// "data:") must not panic or fire the callback.
	s.Feed([]byte("data: not json\n"))
	s.Feed([]byte("data: {\"valid\":true}\n"))

	var want []map[string]any
	if !reflect.DeepEqual(got, want) {
		t.Errorf("got %v, want %v", got, want)
	}
}

func TestScanSSE_Finalization(t *testing.T) {
	for _, closeEarly := range []bool{false, true} {
		for _, malformed := range []bool{false, true} {
			t.Run(fmt.Sprintf("close=%v/malformed=%v", closeEarly, malformed), func(t *testing.T) {
				body := "data: {\"last\":true}"
				if malformed {
					body = "data: {"
				}
				chunks, parseErrors, finalized := 0, 0, 0
				s := NewSSEScanner(func(map[string]any) { chunks++ })
				s.OnError = func(error) { parseErrors++ }
				br := NewBufferedReader(io.NopCloser(strings.NewReader(body)), func(*bytes.Buffer, error) {
					finalized++
					if malformed && parseErrors != 1 || !malformed && chunks != 1 {
						t.Fatal("onDone ran before final line was processed")
					}
				})
				ScanSSE(br, s)
				if closeEarly {
					if _, err := io.ReadFull(br, make([]byte, len(body))); err != nil {
						t.Fatal(err)
					}
				} else {
					if _, err := io.Copy(io.Discard, br); err != nil {
						t.Fatal(err)
					}
				}
				if err := br.Close(); err != nil {
					t.Fatal(err)
				}
				if finalized != 1 {
					t.Fatalf("finalized %d times", finalized)
				}
			})
		}
	}
}

func TestSSEParsingPathsAgree(t *testing.T) {
	for _, input := range []string{
		"event: message\r\ndata:{\"x\":1}\r\n\r\ndata: [DONE]\ndata: {\"ignored\":true}\n",
		"data: {\"x\":1}\ndata: malformed\ndata: {\"ignored\":true}\n",
		"data: {\"x\":1}",
		"data: {",
	} {
		want, wantErr := ParseSSEChunks(strings.NewReader(input))
		var got []map[string]any
		var gotErr error
		s := NewSSEScanner(func(chunk map[string]any) { got = append(got, chunk) })
		s.OnError = func(err error) { gotErr = err }
		for i := range len(input) {
			s.Feed([]byte(input[i : i+1]))
		}
		s.Finish()
		if !reflect.DeepEqual(got, want) || (gotErr == nil) != (wantErr == nil) {
			t.Fatalf("input %q: incremental=(%v, %v), whole-body=(%v, %v)", input, got, gotErr, want, wantErr)
		}
	}
}

func TestParseSSEChunksFunc_ReadError(t *testing.T) {
	readErr := errors.New("read failed")
	r := io.MultiReader(strings.NewReader("data: {\"x\":1}"), &errReader{err: readErr})
	chunks, err := ParseSSEChunks(r)
	if !errors.Is(err, readErr) || len(chunks) != 1 {
		t.Fatalf("got (%v, %v), want final chunk and read error", chunks, err)
	}
}
