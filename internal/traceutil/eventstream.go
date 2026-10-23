package traceutil

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"hash/crc32"
	"io"
)

const (
	eventStreamPreludeLen = 12
	eventStreamMinLen     = eventStreamPreludeLen + 4
	eventStreamMaxLen     = eventStreamMinLen + 128<<10 + 16<<20
)

var errMalformedEventStream = errors.New("malformed event stream message")

// EventStreamScanner incrementally decodes the AWS event stream returned by
// Amazon Bedrock InvokeModelWithResponseStream. Each chunk event carries one
// provider-native JSON event, which is passed to onChunk. Exception messages
// are passed as {"type":"error","error":{"type":...,"message":...}}, matching
// the error event of an Anthropic Messages stream. Not safe for concurrent use.
type EventStreamScanner struct {
	onChunk func(map[string]any)
	// OnError, if set, is called once when the stream cannot be decoded;
	// feeding stops there.
	OnError func(error)
	stopped bool
	buf     []byte
}

func NewEventStreamScanner(onChunk func(map[string]any)) *EventStreamScanner {
	return &EventStreamScanner{onChunk: onChunk}
}

func (s *EventStreamScanner) Feed(p []byte) {
	if s.stopped {
		return
	}
	s.buf = append(s.buf, p...)
	for !s.stopped && len(s.buf) >= eventStreamPreludeLen {
		total := binary.BigEndian.Uint32(s.buf[0:4])
		headersLen := binary.BigEndian.Uint32(s.buf[4:8])
		if crc32.ChecksumIEEE(s.buf[0:8]) != binary.BigEndian.Uint32(s.buf[8:12]) ||
			total < eventStreamMinLen || total > eventStreamMaxLen || headersLen > total-eventStreamMinLen {
			s.fail(errMalformedEventStream)
			return
		}
		if uint32(len(s.buf)) < total {
			return
		}
		message := s.buf[:total]
		if crc32.ChecksumIEEE(message[:total-4]) != binary.BigEndian.Uint32(message[total-4:]) {
			s.fail(errMalformedEventStream)
			return
		}
		headersEnd := eventStreamPreludeLen + headersLen
		s.handle(message[eventStreamPreludeLen:headersEnd], message[headersEnd:total-4])
		if s.stopped {
			return
		}
		s.buf = s.buf[total:]
	}
}

// Finish reports a message left incomplete when the stream ended.
func (s *EventStreamScanner) Finish() {
	if !s.stopped && len(s.buf) > 0 {
		s.fail(io.ErrUnexpectedEOF)
	}
}

func (s *EventStreamScanner) fail(err error) {
	s.stopped = true
	s.buf = nil
	if s.OnError != nil {
		s.OnError(err)
	}
}

func (s *EventStreamScanner) handle(rawHeaders, payload []byte) {
	headers, err := parseEventStreamHeaders(rawHeaders)
	if err != nil {
		s.fail(err)
		return
	}
	switch headers[":message-type"] {
	case "event":
		if headers[":event-type"] != "chunk" {
			return
		}
		var part struct {
			Bytes []byte `json:"bytes"`
		}
		var chunk map[string]any
		if err := json.Unmarshal(payload, &part); err != nil {
			s.fail(fmt.Errorf("%w: %w", errMalformedEventStream, err))
			return
		}
		if err := json.Unmarshal(part.Bytes, &chunk); err != nil {
			s.fail(fmt.Errorf("%w: %w", errMalformedEventStream, err))
			return
		}
		if s.onChunk != nil {
			s.onChunk(chunk)
		}
	case "exception", "error":
		var detail struct {
			Message string `json:"message"`
		}
		_ = json.Unmarshal(payload, &detail)
		errorType := headers[":exception-type"]
		if errorType == "" {
			errorType = headers[":error-code"]
		}
		if detail.Message == "" {
			detail.Message = headers[":error-message"]
		}
		if s.onChunk != nil {
			s.onChunk(map[string]any{
				"type":  "error",
				"error": map[string]any{"type": errorType, "message": detail.Message},
			})
		}
	}
}

// parseEventStreamHeaders returns the string-valued headers of a message and
// skips headers of every other type.
func parseEventStreamHeaders(b []byte) (map[string]string, error) {
	headers := make(map[string]string)
	for len(b) > 0 {
		nameLen := int(b[0])
		if len(b) < nameLen+2 {
			return nil, errMalformedEventStream
		}
		name := string(b[1 : 1+nameLen])
		valueType := b[1+nameLen]
		b = b[nameLen+2:]
		var size int
		switch valueType {
		case 0, 1:
		case 2:
			size = 1
		case 3:
			size = 2
		case 4:
			size = 4
		case 5, 8:
			size = 8
		case 9:
			size = 16
		case 6, 7:
			if len(b) < 2 {
				return nil, errMalformedEventStream
			}
			size = int(binary.BigEndian.Uint16(b))
			b = b[2:]
		default:
			return nil, errMalformedEventStream
		}
		if len(b) < size {
			return nil, errMalformedEventStream
		}
		if valueType == 7 {
			headers[name] = string(b[:size])
		}
		b = b[size:]
	}
	return headers, nil
}

// ScanEventStream feeds everything read from br to scanner, and reports a
// truncated final message when the stream ends. Safe only when br is read by a
// single goroutine.
func ScanEventStream(br *BufferedReader, scanner *EventStreamScanner) {
	br.onBytes = scanner.Feed
	br.onFlush = scanner.Finish
}
