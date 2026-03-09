package traceutil

import (
	"bytes"
	"errors"
	"io"
	"strings"
	"testing"
)

type nopCloser struct {
	io.Reader
	closed bool
}

func (n *nopCloser) Close() error {
	n.closed = true
	return nil
}

type errReader struct {
	err error
}

func (e *errReader) Read([]byte) (int, error) {
	return 0, e.err
}

func TestBufferedReader_ReadToEOF(t *testing.T) {
	src := &nopCloser{Reader: strings.NewReader("hello world")}
	var captured string
	br := NewBufferedReader(src, func(r io.Reader, _ error) {
		data, _ := io.ReadAll(r)
		captured = string(data)
	})

	data, err := io.ReadAll(br)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "hello world" {
		t.Errorf("expected 'hello world', got %q", string(data))
	}
	if captured != "hello world" {
		t.Errorf("onDone captured %q, expected 'hello world'", captured)
	}
}

func TestBufferedReader_CloseBeforeEOF(t *testing.T) {
	src := &nopCloser{Reader: bytes.NewReader([]byte("abcdefghij"))}
	var captured string
	br := NewBufferedReader(src, func(r io.Reader, _ error) {
		data, _ := io.ReadAll(r)
		captured = string(data)
	})

	buf := make([]byte, 3)
	br.Read(buf)
	br.Close()

	if captured != "abc" {
		t.Errorf("onDone captured %q, expected 'abc'", captured)
	}
	if !src.closed {
		t.Error("expected underlying source to be closed")
	}
}

func TestBufferedReader_OnDoneCalledOnce(t *testing.T) {
	src := &nopCloser{Reader: strings.NewReader("data")}
	calls := 0
	br := NewBufferedReader(src, func(r io.Reader, _ error) {
		calls++
	})

	io.ReadAll(br)
	br.Close()

	if calls != 1 {
		t.Errorf("expected onDone called once, got %d", calls)
	}
}

func TestBufferedReader_NilOnDone(t *testing.T) {
	src := &nopCloser{Reader: strings.NewReader("data")}
	br := NewBufferedReader(src, nil)

	data, err := io.ReadAll(br)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if string(data) != "data" {
		t.Errorf("expected 'data', got %q", string(data))
	}
	br.Close()
}

func TestBufferedReader_ReadErrPropagation(t *testing.T) {
	readErr := errors.New("connection reset")
	src := &nopCloser{Reader: io.MultiReader(
		strings.NewReader("partial"),
		&errReader{err: readErr},
	)}
	var captured string
	var gotErr error
	br := NewBufferedReader(src, func(r io.Reader, err error) {
		data, _ := io.ReadAll(r)
		captured = string(data)
		gotErr = err
	})

	_, err := io.ReadAll(br)
	if !errors.Is(err, readErr) {
		t.Fatalf("expected read error %v, got %v", readErr, err)
	}
	if captured != "partial" {
		t.Errorf("onDone captured %q, expected 'partial'", captured)
	}
	if !errors.Is(gotErr, readErr) {
		t.Errorf("onDone error = %v, want %v", gotErr, readErr)
	}
}

func TestBufferedReader_PassesThroughData(t *testing.T) {
	content := strings.Repeat("x", 10000)
	src := &nopCloser{Reader: strings.NewReader(content)}
	br := NewBufferedReader(src, func(r io.Reader, _ error) {})

	data, err := io.ReadAll(br)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(data) != 10000 {
		t.Errorf("expected 10000 bytes, got %d", len(data))
	}
}
