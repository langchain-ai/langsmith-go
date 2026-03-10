package tracesink

import (
	"context"
	"log"
	"sync"
	"time"

	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/models"
	"github.com/langchain-ai/langsmith-go/lib/langsmithtracing/internal/multipart"
)

// RunTransformFunc is a pre-export transform hook matching the Python SDK's
// process_buffered_run_ops callback. It receives a batch of decoded run
// operations and returns (possibly modified) operations. The transform runs
// on every drain cycle, after batching but before coalescing and export.
type RunTransformFunc func(ops []models.RunOp) []models.RunOp

// TraceSink asynchronously batches serialized operations and sends them
// via the multipart exporter. Under load it spawns additional worker
// goroutines to export batches concurrently.
type TraceSink struct {
	exporter  *multipart.Exporter
	config    DrainConfig
	transform RunTransformFunc
	endpoint  models.WriteEndpoint
	ctx       context.Context

	mu     sync.Mutex
	queue  []*models.SerializedOp
	notify chan struct{} // signaled on Submit to wake runLoop

	done   chan struct{}
	closed bool

	workerMu    sync.Mutex
	workerCount int
	workerWg    sync.WaitGroup
}

// NewTraceSink creates and starts a trace sink. The provided context is
// propagated to HTTP requests during normal operation; Close always drains
// with a background context to guarantee delivery.
func NewTraceSink(ctx context.Context, exporter *multipart.Exporter, config DrainConfig, endpoint models.WriteEndpoint, transform RunTransformFunc) *TraceSink {
	s := &TraceSink{
		exporter:  exporter,
		config:    config,
		transform: transform,
		endpoint:  endpoint,
		ctx:       ctx,
		notify:    make(chan struct{}, 1),
		done:      make(chan struct{}),
	}
	go s.runLoop()
	return s
}

// Submit adds a serialized operation to the queue.
// If the queue is at capacity (MaxQueueSize), the operation is dropped with
// a warning log, matching the JS SDK's backpressure behavior.
//
// Submit does not wake the drain loop immediately — the ticker provides
// a natural batching window so that create+update pairs for the same run
// are coalesced before export.
func (s *TraceSink) Submit(op *models.SerializedOp) error {
	s.mu.Lock()
	defer s.mu.Unlock()
	if s.closed {
		return nil
	}
	if s.config.MaxQueueSize > 0 && len(s.queue) >= s.config.MaxQueueSize {
		log.Printf("[langsmith] queue full (%d items); dropping run %s", len(s.queue), op.ID)
		return nil
	}
	s.queue = append(s.queue, op)
	return nil
}

// Close flushes remaining operations and shuts down the sink.
// It spawns workers for a concurrent drain when many items remain.
func (s *TraceSink) Close() {
	s.mu.Lock()
	if s.closed {
		s.mu.Unlock()
		return
	}
	s.closed = true
	s.mu.Unlock()

	// Wake runLoop so it sees closed immediately.
	select {
	case s.notify <- struct{}{}:
	default:
	}
	<-s.done

	s.mu.Lock()
	queueLen := len(s.queue)
	s.mu.Unlock()
	if queueLen > s.config.MaxBatchSize {
		s.spawnWorkers(s.config.MaxWorkers, context.Background())
	}

	for s.drainOnce(context.Background()) {
	}
	s.workerWg.Wait()
}

func (s *TraceSink) runLoop() {
	defer close(s.done)
	ticker := time.NewTicker(s.config.DrainInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
		case <-s.notify:
		}
		s.mu.Lock()
		closed := s.closed
		queueLen := len(s.queue)
		s.mu.Unlock()
		if closed {
			return
		}
		if queueLen > s.config.ScaleUpQueueTrigger {
			s.spawnWorkers(1, s.ctx)
		}
		s.drainOnce(s.ctx)
	}
}

// spawnWorkers starts up to n new worker goroutines, respecting MaxWorkers.
func (s *TraceSink) spawnWorkers(n int, ctx context.Context) {
	s.workerMu.Lock()
	defer s.workerMu.Unlock()
	for i := 0; i < n && s.workerCount < s.config.MaxWorkers; i++ {
		s.workerCount++
		s.workerWg.Add(1)
		go s.runWorker(ctx)
	}
}

// takeBatch removes up to one export-sized batch from the queue.
func (s *TraceSink) takeBatch() []*models.SerializedOp {
	s.mu.Lock()
	defer s.mu.Unlock()
	if len(s.queue) == 0 {
		return nil
	}

	end := len(s.queue)
	if end > s.config.MaxBatchSize {
		end = s.config.MaxBatchSize
	}
	batchBytes := 0
	for i := 0; i < end; i++ {
		sz := s.queue[i].Size()
		if i > 0 && batchBytes+sz > s.config.MaxBatchBytes {
			end = i
			break
		}
		batchBytes += sz
	}

	batch := make([]*models.SerializedOp, end)
	copy(batch, s.queue[:end])
	s.queue = s.queue[end:]
	return batch
}

// drainOnce takes one batch and exports it. Returns true if work was found.
func (s *TraceSink) drainOnce(ctx context.Context) bool {
	batch := s.takeBatch()
	if len(batch) == 0 {
		return false
	}

	if s.transform != nil {
		var err error
		batch, err = s.applyTransform(batch)
		if err != nil {
			log.Printf("[langsmith] transform error: %v", err)
			return true
		}
		if len(batch) == 0 {
			return true
		}
	}

	coalesced, err := models.Coalesce(batch)
	if err != nil {
		log.Printf("[langsmith] coalesce error: %v", err)
		return true
	}
	if err := s.exporter.Export(ctx, s.endpoint, coalesced); err != nil {
		log.Printf("[langsmith] export error: %v", err)
	}
	return true
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

	ops = s.transform(ops)

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

// runWorker drains batches until the queue is consistently empty.
// During shutdown (closed && empty queue) it exits immediately.
func (s *TraceSink) runWorker(ctx context.Context) {
	defer s.workerWg.Done()
	defer func() {
		s.workerMu.Lock()
		s.workerCount--
		s.workerMu.Unlock()
	}()

	emptyRuns := 0
	for emptyRuns <= s.config.ScaleDownEmptyTrigger {
		if s.drainOnce(ctx) {
			emptyRuns = 0
			continue
		}
		s.mu.Lock()
		closed := s.closed
		empty := len(s.queue) == 0
		s.mu.Unlock()
		if closed && empty {
			return
		}
		emptyRuns++
		time.Sleep(s.config.DrainInterval)
	}
}
