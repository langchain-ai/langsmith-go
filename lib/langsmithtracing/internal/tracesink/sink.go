package tracesink

import (
	"context"
	"fmt"
	"sync"
	"sync/atomic"
	"time"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/logger"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/models"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/multipart"
)

// RunTransformFunc is a pre-export transform hook. It receives a batch of decoded run
// operations and returns (possibly modified) operations. The transform runs
// on every drain cycle, after batching but before merging and export.
type RunTransformFunc func(ops []models.RunOp) []models.RunOp

type job struct {
	ctx   context.Context
	batch []*models.SerializedOp
}

// TraceSink asynchronously batches serialized operations and sends them
// via the multipart exporter. A single dispatcher goroutine reads from the
// queue channel, builds batches, and dispatches them to a fixed worker pool.
type TraceSink struct {
	exporter  *multipart.Exporter
	config    DrainConfig
	transform RunTransformFunc
	logger    logger.Logger
	endpoint  models.WriteEndpoint
	ctx       context.Context

	queue   chan *models.SerializedOp // producers: Submit; consumer: dispatcher
	jobs    chan job                  // producer: dispatcher; consumers: workers
	closed  atomic.Bool
	closeCh chan struct{} // signals dispatcher to drain and exit

	closeOnce sync.Once
	doneCh    chan struct{} // closed when dispatcher + all workers finish
}

// NewTraceSink creates and starts a trace sink. The provided context is
// propagated to HTTP requests during normal operation; Close always drains
// with a background context to guarantee delivery.
func NewTraceSink(ctx context.Context, exporter *multipart.Exporter, config DrainConfig, endpoint models.WriteEndpoint, transform RunTransformFunc, l logger.Logger) *TraceSink {
	if l == nil {
		l = logger.DefaultLogger{}
	}
	queueSize := config.MaxQueueSize
	if queueSize <= 0 {
		queueSize = 10_000
	}
	workers := config.MaxWorkers
	if workers <= 0 {
		workers = 1
	}

	s := &TraceSink{
		exporter:  exporter,
		config:    config,
		transform: transform,
		logger:    l,
		endpoint:  endpoint,
		ctx:       ctx,
		queue:     make(chan *models.SerializedOp, queueSize),
		jobs:      make(chan job, workers),
		closeCh:   make(chan struct{}),
		doneCh:    make(chan struct{}),
	}

	var wg sync.WaitGroup
	wg.Add(workers)
	for range workers {
		go func() {
			defer wg.Done()
			s.runWorker()
		}()
	}

	go func() {
		defer close(s.doneCh)
		cancel := s.runDispatcher()
		close(s.jobs)
		wg.Wait()
		if cancel != nil {
			cancel()
		}
	}()

	return s
}

// Submit adds a serialized operation to the queue.
//
// If the sink is closed or the queue is full, the operation is silently
// dropped.
func (s *TraceSink) Submit(op *models.SerializedOp) {
	if s.closed.Load() {
		s.logger.Error("Submit after close: dropping run", "run_id", op.ID)
		return
	}
	select {
	case s.queue <- op:
	default:
		s.logger.Error("Queue full: dropping run", "queue_size", len(s.queue), "max_queue_size", s.config.MaxQueueSize, "run_id", op.ID)
	}
}

// Close flushes remaining operations and shuts down the sink.
// The flush is bounded by CloseTimeout (default 60s); any items still
// queued after the deadline are dropped with a warning.
// Close is safe to call multiple times; only the first call has any effect.
func (s *TraceSink) Close() {
	s.closeOnce.Do(func() {
		s.closed.Store(true)
		close(s.closeCh)
		<-s.doneCh
	})
}

// runDispatcher is the single goroutine that reads from the queue,
// builds batches on a timer, and dispatches them to workers via the jobs channel.
// It returns a cancel function for the drain context (nil if no drain occurred).
// The caller must defer cancel until after all workers have finished.
func (s *TraceSink) runDispatcher() context.CancelFunc {
	ticker := time.NewTicker(s.config.DrainInterval)
	defer ticker.Stop()

	var pending *models.SerializedOp

	for {
		select {
		case <-ticker.C:
			pending = s.dispatchBatch(s.ctx, pending)

		case <-s.closeCh:
			ticker.Stop()
			return s.drainRemaining(pending)
		}
	}
}

// dispatchBatch collects one batch from the queue and sends it to workers.
// Returns any leftover op that didn't fit in the batch (for the next cycle).
func (s *TraceSink) dispatchBatch(ctx context.Context, pending *models.SerializedOp) *models.SerializedOp {
	batch, leftover := s.collectBatch(pending)
	if len(batch) > 0 {
		select {
		case s.jobs <- job{ctx: ctx, batch: batch}:
		case <-ctx.Done():
			s.logger.Warn("context canceled; dropping batch", "batch_size", len(batch))
		}
	}
	return leftover
}

// collectBatch non-blockingly drains available ops from the queue into a batch,
// respecting MaxBatchSize and MaxBatchBytes. If an op would exceed the byte
// limit, it's returned as leftover for the next batch.
func (s *TraceSink) collectBatch(pending *models.SerializedOp) (batch []*models.SerializedOp, leftover *models.SerializedOp) {
	var batchBytes int

	if pending != nil {
		batch = append(batch, pending)
		batchBytes += pending.SizeBytes()
	}

	for len(batch) < s.config.MaxBatchSize {
		select {
		case op := <-s.queue:
			sz := op.SizeBytes()
			if len(batch) > 0 && s.config.MaxBatchBytes > 0 && batchBytes+sz > s.config.MaxBatchBytes {
				return batch, op
			}
			batch = append(batch, op)
			batchBytes += sz
		default:
			return batch, nil
		}
	}
	return batch, nil
}

// drainRemaining flushes all items left in the queue during shutdown.
// It returns a cancel function for the drain context; the caller must
// not cancel it until all workers have finished processing.
func (s *TraceSink) drainRemaining(pending *models.SerializedOp) context.CancelFunc {
	timeout := s.config.CloseTimeout
	if timeout <= 0 {
		timeout = 60 * time.Second
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)

	for {
		batch, leftover := s.collectBatch(pending)
		pending = leftover
		if len(batch) == 0 {
			return cancel
		}
		select {
		case s.jobs <- job{ctx: ctx, batch: batch}:
		case <-ctx.Done():
			remaining := len(s.queue)
			if pending != nil {
				remaining++
			}
			s.logger.Warn("close timed out; dropping pending items", "timeout", timeout, "remaining", remaining)
			return cancel
		}
	}
}

// runWorker processes batches from the jobs channel until it's closed.
func (s *TraceSink) runWorker() {
	for j := range s.jobs {
		s.processBatch(j.ctx, j.batch)
	}
}

// processBatch applies the transform hook, merges patches into posts, and exports.
func (s *TraceSink) processBatch(ctx context.Context, batch []*models.SerializedOp) {
	if s.transform != nil {
		var err error
		batch, err = s.applyTransform(batch)
		if err != nil {
			s.logger.Error("transform error", "error", err)
			return
		}
		if len(batch) == 0 {
			return
		}
	}

	merged, err := models.MergePatchToPost(batch)
	if err != nil {
		s.logger.Error("merge patch to post error", "error", err)
		return
	}
	if err := s.exporter.Export(ctx, s.endpoint, merged); err != nil {
		s.logger.Error("export error", "error", err)
	}
}

func (s *TraceSink) applyTransform(batch []*models.SerializedOp) ([]*models.SerializedOp, error) {
	ops := make([]models.RunOp, len(batch))
	for i, sop := range batch {
		decoded, err := models.DeserializeOp(sop)
		if err != nil {
			return nil, err
		}
		ops[i] = decoded
	}

	var transformErr error
	func() {
		defer func() {
			if r := recover(); r != nil {
				transformErr = fmt.Errorf("transform panicked: %v", r)
			}
		}()
		ops = s.transform(ops)
	}()
	if transformErr != nil {
		return nil, transformErr
	}

	result := make([]*models.SerializedOp, len(ops))
	for i, op := range ops {
		sop, err := models.SerializeOp(op)
		if err != nil {
			return nil, err
		}
		result[i] = sop
	}
	return result, nil
}
