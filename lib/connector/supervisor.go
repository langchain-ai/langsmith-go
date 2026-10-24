package connector

import (
	"context"
	"io"
	"math/rand/v2"
	"net"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hashicorp/yamux"
)

const (
	desiredCarrierCount     = 2
	initialReconnectBackoff = 500 * time.Millisecond
	maxReconnectBackoff     = 30 * time.Second
	healthyCarrierMinUptime = 60 * time.Second
)

type carrierSession struct {
	sessionID uuid.UUID
	runtimeID uuid.UUID
	yamux     *yamux.Session
	control   net.Conn
	transport io.Closer
}

func (c *carrierSession) close() {
	_ = c.control.Close()
	_ = c.yamux.Close()
	_ = c.transport.Close()
}

type carrierSupervisor struct {
	client     *Client
	mu         sync.Mutex
	runtimeIDs map[uuid.UUID]struct{}
}

// Run maintains two registered carriers on distinct edge runtimes until the
// context is cancelled. Each slot reconnects independently, so one live carrier
// remains usable while the other is being replaced.
func (c *Client) Run(ctx context.Context) error {
	if err := ctx.Err(); err != nil {
		return err
	}
	supervisor := &carrierSupervisor{
		client:     c,
		runtimeIDs: make(map[uuid.UUID]struct{}, desiredCarrierCount),
	}
	var slots sync.WaitGroup
	slots.Add(desiredCarrierCount)
	for range desiredCarrierCount {
		go func() {
			defer slots.Done()
			supervisor.runSlot(ctx)
		}()
	}
	slots.Wait()
	return ctx.Err()
}

func (s *carrierSupervisor) runSlot(ctx context.Context) {
	logger := s.client.logger
	backoff := initialReconnectBackoff
	for ctx.Err() == nil {
		carrier, err := s.client.runCarrier(ctx)
		if err != nil {
			if ctx.Err() != nil {
				return
			}
			logger.WarnContext(ctx, "connector could not register with LangSmith; retrying",
				"error", err, "backoff", backoff.String())
			if !waitForCarrierRetry(ctx, backoff) {
				return
			}
			backoff = min(2*backoff, maxReconnectBackoff)
			continue
		}
		if !s.reserve(carrier.runtimeID) {
			logger.DebugContext(ctx, "connector carrier reached an edge runtime it already uses; reconnecting",
				"session_id", carrier.sessionID, "edge_runtime_id", carrier.runtimeID)
			carrier.close()
			if !waitForCarrierRetry(ctx, backoff) {
				return
			}
			backoff = min(2*backoff, maxReconnectBackoff)
			continue
		}
		logger.InfoContext(ctx, "connector registered with LangSmith",
			"session_id", carrier.sessionID, "edge_runtime_id", carrier.runtimeID)

		started := time.Now()
		serveErr := s.client.serveCarrier(ctx, carrier)
		s.release(carrier.runtimeID)
		if ctx.Err() != nil {
			return
		}
		uptime := time.Since(started)
		delay, nextBackoff := carrierReconnectBackoff(backoff, uptime)
		logger.WarnContext(ctx, "connector carrier closed; reconnecting",
			"session_id", carrier.sessionID,
			"error", serveErr,
			"uptime", uptime.Round(time.Second).String(),
			"backoff", delay.String())
		if !waitForCarrierRetry(ctx, delay) {
			return
		}
		backoff = nextBackoff
	}
}

func carrierReconnectBackoff(current, uptime time.Duration) (time.Duration, time.Duration) {
	if uptime >= healthyCarrierMinUptime {
		current = initialReconnectBackoff
	}
	return current, min(2*current, maxReconnectBackoff)
}

func (s *carrierSupervisor) reserve(runtimeID uuid.UUID) bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, exists := s.runtimeIDs[runtimeID]; exists {
		return false
	}
	s.runtimeIDs[runtimeID] = struct{}{}
	return true
}

func (s *carrierSupervisor) release(runtimeID uuid.UUID) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.runtimeIDs, runtimeID)
}

func waitForCarrierRetry(ctx context.Context, backoff time.Duration) bool {
	delay := backoff
	if jitter := backoff / 2; jitter > 0 {
		delay += time.Duration(rand.Int64N(int64(jitter)))
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return false
	case <-timer.C:
		return true
	}
}
