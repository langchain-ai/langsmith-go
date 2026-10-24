package connector

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestResponseHeaderScannerFindsTheFinalHeader(t *testing.T) {
	t.Parallel()

	tests := []struct {
		name   string
		chunks []string
		want   bool
	}{
		{
			name:   "CRLF header",
			chunks: []string{"HTTP/1.1 200 OK\r\nContent-Length: 2\r\n\r\nok"},
			want:   true,
		},
		{
			name:   "LF-only header",
			chunks: []string{"HTTP/1.1 200 OK\nContent-Length: 2\n\nok"},
			want:   true,
		},
		{
			name:   "header still arriving",
			chunks: []string{"HTTP/1.1 200 OK\r\nContent-Length:"},
		},
		{
			name:   "line of only carriage returns",
			chunks: []string{"HTTP/1.1 200 OK\r\n\r\r\n"},
		},
		{
			name:   "terminator split across reads",
			chunks: []string{"HTTP/1.1 204 No Content\r\n\r", "\n"},
			want:   true,
		},
		{
			name:   "status line split across reads",
			chunks: []string{"HTTP/1.1 1", "00 Continue\r\n\r\n"},
		},
		{
			name:   "only an informational response",
			chunks: []string{"HTTP/1.1 100 Continue\r\n\r\n"},
		},
		{
			name:   "informational then final response",
			chunks: []string{"HTTP/1.1 100 Continue\r\n\r\n", "HTTP/1.1 200 OK\r\n\r\n"},
			want:   true,
		},
		{
			name:   "switching protocols is final",
			chunks: []string{"HTTP/1.1 101 Switching Protocols\r\nUpgrade: websocket\r\n\r\n"},
			want:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()

			var scanner responseHeaderScanner
			got := false
			for _, chunk := range tt.chunks {
				got = scanner.scan([]byte(chunk))
				if got {
					break
				}
			}
			assert.Equal(t, tt.want, got)
		})
	}
}
