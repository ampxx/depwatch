package scheduler_test

import (
	"context"
	"sync/atomic"
	"testing"
	"time"

	"github.com/depwatch/internal/scheduler"
)

func TestScheduler_RunsJobAtInterval(t *testing.T) {
	s := scheduler.New()
	var count atomic.Int64

	s.Register(scheduler.Job{
		Name:     "counter",
		Interval: 20 * time.Millisecond,
		Fn: func(ctx context.Context) error {
			count.Add(1)
			return nil
		},
	})

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	got := count.Load()
	if got < 2 {
		t.Errorf("expected at least 2 executions, got %d", got)
	}
}

func TestScheduler_StopsOnContextCancel(t *testing.T) {
	s := scheduler.New()
	var count atomic.Int64

	s.Register(scheduler.Job{
		Name:     "stopper",
		Interval: 10 * time.Millisecond,
		Fn: func(ctx context.Context) error {
			count.Add(1)
			return nil
		},
	})

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		time.Sleep(35 * time.Millisecond)
		cancel()
	}()
	s.Run(ctx)

	snap := count.Load()
	time.Sleep(30 * time.Millisecond)
	if count.Load() != snap {
		t.Error("job continued running after context was cancelled")
	}
}

func TestScheduler_MultipleJobs_RunConcurrently(t *testing.T) {
	s := scheduler.New()
	var a, b atomic.Int64

	s.Register(scheduler.Job{
		Name: "a", Interval: 15 * time.Millisecond,
		Fn: func(ctx context.Context) error { a.Add(1); return nil },
	})
	s.Register(scheduler.Job{
		Name: "b", Interval: 15 * time.Millisecond,
		Fn: func(ctx context.Context) error { b.Add(1); return nil },
	})

	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Millisecond)
	defer cancel()
	s.Run(ctx)

	if a.Load() < 1 || b.Load() < 1 {
		t.Errorf("expected both jobs to run, got a=%d b=%d", a.Load(), b.Load())
	}
}

func TestScheduler_NoJobs_ReturnsImmediately(t *testing.T) {
	s := scheduler.New()
	ctx, cancel := context.WithTimeout(context.Background(), 50*time.Millisecond)
	defer cancel()
	done := make(chan struct{})
	go func() {
		s.Run(ctx)
		close(done)
	}()
	select {
		case <-done:
		case <-time.After(100 * time.Millisecond):
			t.Error("Run did not return after context cancel with no jobs")
	}
}
