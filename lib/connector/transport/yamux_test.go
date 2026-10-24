package transport

import (
	"testing"
	"time"

	"github.com/hashicorp/yamux"
)

func TestYamuxConfigUsesBoundedCarrierDefaults(t *testing.T) {
	t.Parallel()

	config := YamuxConfig()
	if err := yamux.VerifyConfig(config); err != nil {
		t.Fatalf("VerifyConfig() error = %v", err)
	}
	if !config.EnableKeepAlive {
		t.Fatal("EnableKeepAlive = false, want true")
	}
	if config.KeepAliveInterval != 30*time.Second {
		t.Fatalf("KeepAliveInterval = %v, want %v", config.KeepAliveInterval, 30*time.Second)
	}
	if config.ConnectionWriteTimeout != 30*time.Second {
		t.Fatalf("ConnectionWriteTimeout = %v, want %v", config.ConnectionWriteTimeout, 30*time.Second)
	}
	if config.MaxStreamWindowSize != MaxWebSocketMessageSizeBytes {
		t.Fatalf("MaxStreamWindowSize = %d, want %d", config.MaxStreamWindowSize, MaxWebSocketMessageSizeBytes)
	}
	if config.AcceptBacklog != 128 {
		t.Fatalf("AcceptBacklog = %d, want 128", config.AcceptBacklog)
	}
	if config.StreamOpenTimeout != 15*time.Second {
		t.Fatalf("StreamOpenTimeout = %v, want %v", config.StreamOpenTimeout, 15*time.Second)
	}
	if config.StreamCloseTimeout != 5*time.Minute {
		t.Fatalf("StreamCloseTimeout = %v, want %v", config.StreamCloseTimeout, 5*time.Minute)
	}
	if config.Logger == nil || config.LogOutput != nil {
		t.Fatal("yamux logging must route through a Logger with LogOutput unset")
	}
}

func TestYamuxConfigReturnsIndependentConfigs(t *testing.T) {
	t.Parallel()

	first := YamuxConfig()
	first.AcceptBacklog = 1
	second := YamuxConfig()
	if second.AcceptBacklog != 128 {
		t.Fatalf("second AcceptBacklog = %d, want 128", second.AcceptBacklog)
	}
}
