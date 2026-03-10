package tracesink

import (
	"runtime"
	"time"
)

// DrainConfig controls batching, drain behavior, and auto-scaling for the
// trace sink. When the pending queue exceeds ScaleUpQueueTrigger, additional
// worker goroutines are spawned (up to MaxWorkers) to export batches
// concurrently. Workers exit after ScaleDownEmptyTrigger consecutive empty
// drain cycles, matching the Python SDK's sub-thread scaling model.
type DrainConfig struct {
	MaxBatchSize  int
	MaxBatchBytes int
	MaxQueueSize  int // pending ops beyond this limit are dropped with a warning (0 = unlimited)
	DrainInterval time.Duration

	ScaleUpQueueTrigger   int // queue depth that triggers a new worker goroutine
	MaxWorkers            int // max concurrent export goroutines (excluding the main loop)
	ScaleDownEmptyTrigger int // consecutive empty drains before a worker exits
}

// DefaultDrainConfig returns production-grade defaults.
// MaxWorkers is capped to min(GOMAXPROCS, 32), matching the Python SDK ceiling.
func DefaultDrainConfig() DrainConfig {
	maxWorkers := runtime.GOMAXPROCS(0)
	if maxWorkers > 32 {
		maxWorkers = 32
	}
	return DrainConfig{
		MaxBatchSize:          100,
		MaxBatchBytes:         20 * 1024 * 1024, // 20 MiB
		MaxQueueSize:          10_000,
		DrainInterval:         250 * time.Millisecond,
		ScaleUpQueueTrigger:   200,
		MaxWorkers:            maxWorkers,
		ScaleDownEmptyTrigger: 4,
	}
}
