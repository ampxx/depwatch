package watcher

import (
	"context"
	"time"

	"github.com/depwatch/internal/scheduler"
)

// ScheduledWatcher wraps a Runner inside a Scheduler Job so that
// version checks are driven by the central scheduler rather than
// by an ad-hoc ticker inside the runner itself.
type ScheduledWatcher struct {
	runner    *Runner
	scheduler *scheduler.Scheduler
	interval  time.Duration
}

// NewScheduledWatcher creates a ScheduledWatcher that will trigger the
// supplied Runner every interval using the provided Scheduler.
func NewScheduledWatcher(r *Runner, s *scheduler.Scheduler, interval time.Duration) *ScheduledWatcher {
	sw := &ScheduledWatcher{
		runner:    r,
		scheduler: s,
		interval:  interval,
	}
	s.Register(scheduler.Job{
		Name:     "depwatch-version-check",
		Interval: interval,
		Fn:       sw.tick,
	})
	return sw
}

// tick is the handler invoked by the Scheduler on each interval.
func (sw *ScheduledWatcher) tick(ctx context.Context) error {
	return sw.runner.RunOnce(ctx)
}
