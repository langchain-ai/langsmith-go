package traceutil

import (
	"bytes"
	"encoding/json"
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
	// Search p directly and buffer only an unterminated trailing line, so a
	// long line is not rescanned once per Feed.
	for {
		if s.stopped {
			return
		}
		idx := bytes.IndexByte(p, '\n')
		if idx < 0 {
			s.buf.Write(p) // incomplete; the rest arrives in a later Feed
			return
		}
		line := bytes.TrimRight(p[:idx], "\r\n")
		if s.buf.Len() > 0 { // this line began in an earlier Feed
			s.buf.Write(line)
			line = s.buf.Bytes()
		}
		s.handle(line)
		s.buf.Reset()
		p = p[idx+1:]
	}
}

func (s *SSEScanner) handle(line []byte) {
	if !bytes.HasPrefix(line, []byte("data:")) {
		return
	}
	payload := bytes.TrimSpace(bytes.TrimPrefix(line, []byte("data:")))
	if len(payload) == 0 {
		return
	}
	if bytes.Equal(payload, []byte("[DONE]")) {
		if s.OnDone != nil {
			s.stopped = true
			s.OnDone()
		}
		return
	}
	var chunk map[string]any
	if err := json.Unmarshal(payload, &chunk); err != nil {
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
	s.handle(bytes.TrimRight(s.buf.Bytes(), "\r\n"))
	s.buf.Reset()
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
