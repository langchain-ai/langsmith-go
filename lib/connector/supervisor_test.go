package connector

import (
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func TestCarrierSupervisorReservesDistinctRuntimes(t *testing.T) {
	runtimeA := uuid.Must(uuid.NewV7())
	runtimeB := uuid.Must(uuid.NewV7())
	supervisor := &carrierSupervisor{runtimeIDs: make(map[uuid.UUID]struct{})}

	require.True(t, supervisor.reserve(runtimeA))
	require.False(t, supervisor.reserve(runtimeA))
	require.True(t, supervisor.reserve(runtimeB))
	supervisor.release(runtimeA)
	require.True(t, supervisor.reserve(runtimeA))
}

func TestCarrierReconnectBackoff(t *testing.T) {
	t.Parallel()
	tests := []struct {
		name      string
		current   time.Duration
		uptime    time.Duration
		wantDelay time.Duration
		wantNext  time.Duration
	}{
		{name: "brief carrier", current: 500 * time.Millisecond, uptime: time.Second, wantDelay: 500 * time.Millisecond, wantNext: time.Second},
		{name: "maximum", current: 20 * time.Second, uptime: time.Second, wantDelay: 20 * time.Second, wantNext: 30 * time.Second},
		{name: "healthy carrier", current: 20 * time.Second, uptime: 61 * time.Second, wantDelay: 500 * time.Millisecond, wantNext: time.Second},
	}
	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			t.Parallel()
			delay, next := carrierReconnectBackoff(test.current, test.uptime)
			require.Equal(t, test.wantDelay, delay)
			require.Equal(t, test.wantNext, next)
		})
	}
}
