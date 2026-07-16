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
	scanner := bufio.NewScanner(r)
	scanner.Buffer(make([]byte, 0, 64*1024), maxSSELineBytes)
	var chunks []map[string]any

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
			return chunks, err
		}

		chunks = append(chunks, chunk)
	}

	if err := scanner.Err(); err != nil {
		return chunks, err
	}

	return chunks, nil
}
