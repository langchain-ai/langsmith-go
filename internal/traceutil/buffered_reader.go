package traceutil

import (
	"bytes"
	"io"
	"sync"
)

// BufferedReader saves data read from the source and triggers onDone
// when the source reaches EOF, is closed, or Read returns an error (e.g. context.Canceled).
// onDone receives the buffered content and the error that ended the read (nil for EOF/Close).
// This allows response bodies to stream through while capturing content and real errors for span tagging.
type BufferedReader struct {
	src      io.ReadCloser
	buf      *bytes.Buffer
	bufLimit int
	onDone   func(*bytes.Buffer, error)
	onBytes  func([]byte)
	onFlush  func()
	once     sync.Once
}

// NewBufferedReader creates a BufferedReader that calls onDone with the
// buffered content when the source reaches EOF, is closed, or Read returns an error.
// The second argument to onDone is the error that ended the read, or nil for EOF/Close.
func NewBufferedReader(src io.ReadCloser, onDone func(*bytes.Buffer, error)) *BufferedReader {
	return &BufferedReader{
		src:    src,
		buf:    &bytes.Buffer{},
		onDone: onDone,
	}
}

func (r *BufferedReader) Read(p []byte) (int, error) {
	n, err := r.src.Read(p)
	if n > 0 {
		if r.onBytes != nil {
			r.onBytes(p[:n])
		}
		if r.bufLimit == 0 {
			r.buf.Write(p[:n])
		} else if remaining := r.bufLimit - r.buf.Len(); remaining > 0 {
			r.buf.Write(p[:min(n, remaining)])
		}
	}
	if err != nil {
		r.trigger(err)
	}
	return n, err
}

func (r *BufferedReader) Close() error {
	r.trigger(nil) // closed by consumer, no read error
	return r.src.Close()
}

// LimitBuffer caps how much of the source is retained for onDone.
func (r *BufferedReader) LimitBuffer(n int) {
	r.bufLimit = n
}

func (r *BufferedReader) trigger(readErr error) {
	r.once.Do(func() {
		if r.onFlush != nil {
			r.onFlush()
		}
		if r.onDone != nil {
			r.onDone(r.buf, readErr)
		}
	})
}

// TruncateString converts at most the first limit bytes of data to a string,
// appending "..." if data was longer.
func TruncateString(data []byte, limit int) string {
	if len(data) <= limit {
		return string(data)
	}
	return string(data[:limit]) + "..."
}
