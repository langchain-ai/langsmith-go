package connector

import "bytes"

// statusCodeEnd is the length of a status line up to and including its code.
const statusCodeEnd = len("HTTP/1.1 100")

// responseHeaderScanner finds the end of an HTTP response's final header in
// the raw bytes an origin sends, reading them the way the edge's HTTP client
// does: a header block ends at an empty line, a line may end in LF or CRLF,
// and an informational 1xx response other than 101 Switching Protocols is
// followed by another header block.
type responseHeaderScanner struct {
	// statusLine holds the current block's first bytes, up to the status code.
	statusLine []byte
	lineLength int
	lastByte   byte
	pastStatus bool
}

// scan consumes chunk and reports whether it completed a final response header.
func (s *responseHeaderScanner) scan(chunk []byte) bool {
	for _, b := range chunk {
		if b != '\n' {
			if !s.pastStatus && len(s.statusLine) < statusCodeEnd {
				s.statusLine = append(s.statusLine, b)
			}
			s.lineLength++
			s.lastByte = b
			continue
		}
		// Go's header reader drops only the one CR before LF, so a line of
		// more than that is not the empty line that ends a block.
		empty := s.lineLength == 0 || (s.lineLength == 1 && s.lastByte == '\r')
		s.lineLength = 0
		if !empty {
			s.pastStatus = true
			continue
		}
		if s.informational() {
			*s = responseHeaderScanner{}
			continue
		}
		return true
	}
	return false
}

func (s *responseHeaderScanner) informational() bool {
	if len(s.statusLine) < statusCodeEnd || !bytes.HasPrefix(s.statusLine, []byte("HTTP/")) {
		return false
	}
	code := string(s.statusLine[statusCodeEnd-3:])
	return code[0] == '1' && code != "101"
}
