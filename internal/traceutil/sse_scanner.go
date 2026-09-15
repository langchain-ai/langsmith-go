package traceutil

import (
	"bytes"
	"encoding/json"
	"strings"
)

// SSEScanner incrementally parses Server-Sent Events. Not safe for concurrent use.
type SSEScanner struct {
	onChunk func(map[string]any)
	// OnDone, if set, is called at the "data: [DONE]" sentinel and feeding
	// stops there. OnError, if set, is called when a data line fails to
	// unmarshal and feeding stops there too. Registering either opts into
	// ParseSSEChunksFunc's semantics, which break and return on those lines;
	// leaving them nil keeps the default lenient scan, which skips and
	// continues (what OnFirstSSEMatch wants).
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
	s.buf.Write(p)
	for {
		if s.stopped {
			return
		}
		idx := bytes.IndexByte(s.buf.Bytes(), '\n')
		if idx < 0 {
			return
		}
		line := strings.TrimRight(string(s.buf.Next(idx+1)), "\r\n")
		s.handle(line)
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
		if s.OnDone != nil {
			s.stopped = true
			s.OnDone()
		}
		return
	}
	var chunk map[string]any
	if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
		if s.OnError != nil {
			s.stopped = true
			s.OnError(err)
		}
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
	if s.stopped || s.buf.Len() == 0 {
		return
	}
	line := strings.TrimRight(s.buf.String(), "\r\n")
	s.buf.Reset()
	s.handle(line)
}

// ScanSSE feeds everything read from br to scanner, and flushes it when the
// stream ends so a trailing partial line is not lost. Safe only when br is read
// by a single goroutine.
func ScanSSE(br *BufferedReader, scanner *SSEScanner) {
	br.onBytes = scanner.Feed
	br.onFlush = scanner.Finish
}

// OnFirstSSEMatch fires once on the first SSE chunk satisfying isMatch, then
// detaches from br. Safe only when br is read by a single goroutine.
func OnFirstSSEMatch(br *BufferedReader, isMatch func(map[string]any) bool, fire func()) {
	var fired bool
	scanner := NewSSEScanner(func(chunk map[string]any) {
		if fired || !isMatch(chunk) {
			return
		}
		fired = true
		fire()
		br.onBytes = nil
	})
	br.onBytes = func(b []byte) { scanner.Feed(b) }
}
