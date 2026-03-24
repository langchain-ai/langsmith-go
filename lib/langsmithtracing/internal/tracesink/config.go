package tracesink

import (
	"runtime"
	"time"
)

// DrainConfig controls batching, drain behavior, and the worker pool for the
// trace sink. A fixed pool of MaxWorkers goroutines processes batches
// dispatched by a single dispatcher goroutine.
type DrainConfig struct {
	MaxBatchSize  int
	MaxBatchBytes int
	MaxQueueSize  int // buffered channel capacity; 0 uses default (10 000)
	DrainInterval time.Duration
	CloseTimeout  time.Duration // max time Close() will spend flushing; 0 = 60s default
	MaxWorkers    int           // fixed worker pool size; 0 uses default

	// Deprecated: no longer used. Workers are now a fixed pool.
	ScaleUpQueueTrigger int
	// Deprecated: no longer used. Workers are now a fixed pool.
	ScaleDownEmptyTrigger int
}

// DefaultDrainConfig returns production-grade defaults.
// MaxWorkers is capped to min(GOMAXPROCS, 32).
func DefaultDrainConfig() DrainConfig {
	maxWorkers := runtime.GOMAXPROCS(0)
	if maxWorkers > 32 {
		maxWorkers = 32
	}
	return DrainConfig{
		MaxBatchSize:  100,
		MaxBatchBytes: 20 * 1024 * 1024, // 20 MiB
		MaxQueueSize:  10_000,
		DrainInterval: 250 * time.Millisecond,
		CloseTimeout:  60 * time.Second,
		MaxWorkers:    maxWorkers,
	}
}
