package transport

import (
	"log/slog"
	"time"

	"github.com/hashicorp/yamux"
)

// YamuxConfig returns a fresh configuration shared by both carrier peers.
func YamuxConfig() *yamux.Config {
	config := yamux.DefaultConfig()
	config.AcceptBacklog = 128
	config.EnableKeepAlive = true
	config.KeepAliveInterval = 30 * time.Second
	config.ConnectionWriteTimeout = 30 * time.Second
	config.MaxStreamWindowSize = MaxWebSocketMessageSizeBytes
	config.StreamOpenTimeout = 15 * time.Second
	// yamux Close is a half-close. Its timeout ultimately sends RST when the
	// peer never completes the other half, while leaving a substantial window
	// for legitimate response-after-EOF and asymmetric TCP flows.
	config.StreamCloseTimeout = 5 * time.Minute
	// Route yamux's internal logging into the service's structured logger at
	// debug level rather than its default raw os.Stderr writes. Its chatter on
	// transient, expected conditions (keepalive, resets, drains) then stays out
	// of the normal log stream but is recoverable with debug enabled, while
	// actionable failures still surface as returned errors. LogOutput must be
	// nil when Logger is set — yamux rejects both being set.
	config.Logger = slog.NewLogLogger(slog.Default().Handler(), slog.LevelDebug)
	config.LogOutput = nil
	return config
}
