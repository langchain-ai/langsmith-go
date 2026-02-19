package traceutil

import (
	"bytes"
	"io"
	"sync"
)

// BufferedReader saves data read from the source and triggers onDone
// when fully read or closed. This allows response bodies to stream
// through to the caller while still capturing the full content for
// span tagging.
type BufferedReader struct {
	src    io.ReadCloser
	buf    *bytes.Buffer
	onDone func(io.Reader)
	once   sync.Once
}

// NewBufferedReader creates a BufferedReader that calls onDone with the
// buffered content when the source reaches EOF or is closed.
func NewBufferedReader(src io.ReadCloser, onDone func(io.Reader)) *BufferedReader {
	return &BufferedReader{
		src:    src,
		buf:    &bytes.Buffer{},
		onDone: onDone,
	}
}

func (r *BufferedReader) Read(p []byte) (int, error) {
	n, err := r.src.Read(p)
	if n > 0 {
		r.buf.Write(p[:n])
	}
	if err == io.EOF {
		r.trigger()
	}
	return n, err
}

func (r *BufferedReader) Close() error {
	r.trigger()
	return r.src.Close()
}

func (r *BufferedReader) trigger() {
	r.once.Do(func() {
		if r.onDone != nil {
			r.onDone(r.buf)
		}
	})
}
