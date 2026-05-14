package watcher

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/yourorg/depwatch/internal/config"
	"github.com/yourorg/depwatch/internal/metrics"
)

// Runner orchestrates multiple Watcher instances, one per configured
// dependency, and manages their lifecycle under a shared context.
type Runner struct {
	cfg      *config.Config
	metrics  *metrics.Metrics
	watchers []*Watcher
}

// NewRunner creates a Runner from the provided configuration and metrics
// collector. It initialises one Watcher per dependency defined in cfg.
func NewRunner(cfg *config.Config, m *metrics.Metrics) (*Runner, error) {
	if cfg == nil {
		return nil, fmt.Errorf("watcher: config must not be nil")
	}
	if m == nil {
		return nil, fmt.Errorf("watcher: metrics must not be nil")
	}

	watchers := make([]*Watcher, 0, len(cfg.Dependencies))
	for _, dep := range cfg.Dependencies {
		w, err := New(dep, cfg, m)
		if err != nil {
			return nil, fmt.Errorf("watcher: failed to create watcher for %s: %w", dep.Module, err)
		}
		watchers = append(watchers, w)
	}

	return &Runner{
		cfg:      cfg,
		metrics:  m,
		watchers: watchers,
	}, nil
}

// Run starts all watchers concurrently and blocks until ctx is cancelled.
// Each watcher polls its dependency at the interval defined in the config.
// Run returns only after all watcher goroutines have exited.
func (r *Runner) Run(ctx context.Context) {
	if len(r.watchers) == 0 {
		log.Println("runner: no dependencies configured, nothing to watch")
		<-ctx.Done()
		return
	}

	done := make(chan struct{}, len(r.watchers))

	for _, w := range r.watchers {
		w := w // capture loop variable
		go func() {
			defer func() { done <- struct{}{} }()
			w.Watch(ctx)
		}()
	}

	log.Printf("runner: watching %d dependenc(ies) with poll interval %s",
		len(r.watchers), r.cfg.PollInterval.String())

	// Wait for all watchers to finish after context cancellation.
	for range r.watchers {
		select {
		case <-done:
		case <-time.After(10 * time.Second):
			log.Println("runner: timed out waiting for watchers to stop")
			return
		}
	}

	log.Println("runner: all watchers stopped")
}

// Len returns the number of watchers managed by this Runner.
func (r *Runner) Len() int {
	return len(r.watchers)
}
