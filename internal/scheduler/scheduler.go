package scheduler

import (
	"context"
	"sync"
	"time"
)

// Job represents a named, periodic task managed by the Scheduler.
type Job struct {
	Name     string
	Interval time.Duration
	Fn       func(ctx context.Context) error
}

// Scheduler runs a set of Jobs at their configured intervals.
type Scheduler struct {
	mu   sync.Mutex
	jobs []Job
}

// New returns an initialised, empty Scheduler.
func New() *Scheduler {
	return &Scheduler{}
}

// Register adds a Job to the Scheduler. It must be called before Run.
func (s *Scheduler) Register(j Job) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.jobs = append(s.jobs, j)
}

// Run starts all registered jobs and blocks until ctx is cancelled.
// Each job is executed in its own goroutine on its own ticker.
func (s *Scheduler) Run(ctx context.Context) {
	s.mu.Lock()
	jobs := make([]Job, len(s.jobs))
	copy(jobs, s.jobs)
	s.mu.Unlock()

	var wg sync.WaitGroup
	for _, j := range jobs {
		wg.Add(1)
		go func(job Job) {
			defer wg.Done()
			ticker := time.NewTicker(job.Interval)
			defer ticker.Stop()
			for {
				select {
				case <-ctx.Done():
					return
				case <-ticker.C:
					_ = job.Fn(ctx)
				}
			}
		}(j)
	}
	wg.Wait()
}
