package traceutil

import "io"

const maxSSELineBytes = 16 << 20

// ParseSSEChunks reads an SSE stream and returns the parsed JSON objects
// from each "data: " line. It skips event/id/retry/empty lines and stops
// when it encounters a "data: [DONE]" sentinel.
func ParseSSEChunks(r io.Reader) ([]map[string]any, error) {
	var chunks []map[string]any
	err := ParseSSEChunksFunc(r, func(chunk map[string]any) {
		chunks = append(chunks, chunk)
	})
	return chunks, err
}

// ParseSSEChunksFunc parses the same stream as ParseSSEChunks but hands each
// object to onChunk instead of retaining it, so peak memory does not grow with
// the length of the stream. Parsing stops at the first malformed line, which is
// returned as an error; objects before it have already been passed to onChunk.
func ParseSSEChunksFunc(r io.Reader, onChunk func(map[string]any)) error {
	scanner := NewSSEScanner(onChunk)
	var parseErr error
	scanner.OnError = func(err error) { parseErr = err }
	buf := make([]byte, 32*1024)
	for {
		n, err := r.Read(buf)
		scanner.Feed(buf[:n])
		if scanner.stopped {
			return parseErr
		}
		if err != nil {
			scanner.Finish()
			if parseErr != nil {
				return parseErr
			}
			if err == io.EOF {
				return nil
			}
			return err
		}
	}
}
