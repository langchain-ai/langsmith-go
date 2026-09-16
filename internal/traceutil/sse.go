package traceutil

import (
	"bufio"
	"encoding/json"
	"io"
	"strings"
)

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
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), maxSSELineBytes)

	for scanner.Scan() {
		line := scanner.Text()

		if !strings.HasPrefix(line, "data: ") {
			continue
		}

		line = strings.TrimPrefix(line, "data: ")
		if line == "[DONE]" {
			break
		}

		var chunk map[string]any
		if err := json.Unmarshal([]byte(line), &chunk); err != nil {
			return err
		}

		onChunk(chunk)
	}

	return scanner.Err()
}
