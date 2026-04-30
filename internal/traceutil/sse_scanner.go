package traceutil

import (
	"bytes"
	"encoding/json"
	"strings"
)

// SSEScanner incrementally parses Server-Sent Events. Not safe for concurrent use.
type SSEScanner struct {
	onChunk func(map[string]any)
	buf     bytes.Buffer
}

func NewSSEScanner(onChunk func(map[string]any)) *SSEScanner {
	return &SSEScanner{onChunk: onChunk}
}

func (s *SSEScanner) Feed(p []byte) {
	s.buf.Write(p)
	for {
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
	if payload == "" || payload == "[DONE]" {
		return
	}
	var chunk map[string]any
	if err := json.Unmarshal([]byte(payload), &chunk); err != nil {
		return
	}
	if s.onChunk != nil {
		s.onChunk(chunk)
	}
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
