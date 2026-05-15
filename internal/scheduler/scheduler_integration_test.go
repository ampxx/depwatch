package scheduler_test

import (
	"context"
	"errors"
	"sync/atomic"
	"testing"
	"time"

	"github.com/depwatch/internal/scheduler"
)

// TestScheduler_ErrorDoesNotStopJob verifies that a job returning an error
// continues to be invoked on subsequent ticks.
func TestScheduler_ErrorDoesNotStopJob(t *testing.T) {
	s := scheduler.New()
	var count atomic.Int64

	s.Register(scheduler.Job{
		Name:     "flaky",
		Interval: 15 * time.Millisecond,
		Fn: func(ctx context.Context) error {
			count.Add(1)
			return errors.New("transient error")
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	if count.Load() < 2 {
		t.Errorf("job should keep running despite errors, got %d calls", count.Load())
	}
}

// TestScheduler_RegisterAfterRun_NotPickedUp ensures that jobs registered
// after Run is called are not started mid-flight (snapshot semantics).
func TestScheduler_RegisterAfterRun_NotPickedUp(t *testing.T) {
	s := scheduler.New()
	var lateCount atomic.Int64

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()

	go func() {
		time.Sleep(10 * time.Millisecond)
		s.Register(scheduler.Job{
			Name:     "late",
			Interval: 5 * time.Millisecond,
			Fn: func(ctx context.Context) error {
				lateCount.Add(1)
				return nil
			},
		})
	}()

	s.Run(ctx)

	if lateCount.Load() != 0 {
		t.Errorf("late-registered job should not run, but ran %d times", lateCount.Load())
	}
}
