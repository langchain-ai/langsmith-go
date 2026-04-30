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
	src     io.ReadCloser
	buf     *bytes.Buffer
	onDone  func(io.Reader, error)
	onBytes func([]byte)
	once    sync.Once
}

// NewBufferedReader creates a BufferedReader that calls onDone with the
// buffered content when the source reaches EOF, is closed, or Read returns an error.
// The second argument to onDone is the error that ended the read, or nil for EOF/Close.
func NewBufferedReader(src io.ReadCloser, onDone func(io.Reader, error)) *BufferedReader {
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
		r.buf.Write(p[:n])
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

func (r *BufferedReader) trigger(readErr error) {
	r.once.Do(func() {
		if r.onDone != nil {
			r.onDone(r.buf, readErr)
		}
	})
}
