// Package scheduler provides a lightweight periodic job runner for depwatch.
//
// A Scheduler holds a set of named Jobs, each with its own interval and
// handler function. Calling Run blocks until the supplied context is
// cancelled, at which point all goroutines are stopped cleanly.
//
// Typical usage:
//
//	s := scheduler.New()
//	s.Register(scheduler.Job{
//		Name:     "version-check",
//		Interval: 5 * time.Minute,
//		Fn:       myCheckFn,
//	})
//	s.Run(ctx)
package scheduler
