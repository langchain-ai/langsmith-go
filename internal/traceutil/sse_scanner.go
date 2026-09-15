package traceutil

import (
	"bufio"
	"bytes"
	"encoding/json"
	"strings"
)

// SSEScanner incrementally parses Server-Sent Events. Not safe for concurrent use.
type SSEScanner struct {
	onChunk func(map[string]any)
	// Callbacks report terminal events; parsing stops even when they are nil.
	OnDone  func()
	OnError func(error)
	stopped bool
	buf     bytes.Buffer
}

func NewSSEScanner(onChunk func(map[string]any)) *SSEScanner {
	return &SSEScanner{onChunk: onChunk}
}

func (s *SSEScanner) Feed(p []byte) {
	if s.stopped {
		return
	}
	for len(p) > 0 && !s.stopped {
		idx := bytes.IndexByte(p, '\n')
		n := len(p)
		if idx >= 0 {
			n = idx + 1
		}
		remaining := maxSSELineBytes - s.buf.Len()
		if n > remaining || (n == remaining && idx < 0) {
			s.fail(bufio.ErrTooLong)
			return
		}
		s.buf.Write(p[:n])
		p = p[n:]
		if idx < 0 {
			return
		}
		line := strings.TrimRight(s.buf.String(), "\r\n")
		s.buf.Reset()
		s.handle(line)
	}
}

func (s *SSEScanner) stop() {
	s.stopped = true
	s.buf = bytes.Buffer{}
}

func (s *SSEScanner) fail(err error) {
	s.stop()
	if s.OnError != nil {
		s.OnError(err)
	}
}

func (s *SSEScanner) handle(line string) {
	if !strings.HasPrefix(line, "data:") {
		return
	}
	payload := strings.TrimSpace(strings.TrimPrefix(line, "data:"))
	if payload == "" {
		return
	}
	if payload == "[DONE]" {
		s.stop()
		if s.OnDone != nil {
			s.OnDone()
		}
		return
	}
	var chunk map[string]any
	if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
		s.fail(err)
		return
	}
	if s.onChunk != nil {
		s.onChunk(chunk)
	}
}

// Finish handles any bytes left after the final newline. A bufio-based parse
// surfaces those as a last (usually malformed) line; without this they would be
// silently dropped, so a truncated stream would look clean.
func (s *SSEScanner) Finish() {
	if s.stopped {
		return
	}
	line := strings.TrimRight(s.buf.String(), "\r\n")
	s.stop()
	s.handle(line)
}

// ScanSSE feeds everything read from br to scanner, and flushes it when the
// stream ends so a trailing partial line is not lost. Safe only when br is read
// by a single goroutine.
func ScanSSE(br *BufferedReader, scanner *SSEScanner) {
	br.onBytes = scanner.Feed
	br.onFlush = scanner.Finish
}

// OnFirstMatch wraps a chunk consumer and fires once on the first matching
// chunk. Later chunks still reach consume, but no longer call isMatch.
func OnFirstMatch(consume func(map[string]any), isMatch func(map[string]any) bool, fire func()) func(map[string]any) {
	var fired bool
	return func(chunk map[string]any) {
		consume(chunk)
		if !fired && isMatch(chunk) {
			fired = true
			fire()
		}
	}
}
