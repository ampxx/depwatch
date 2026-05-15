// Package retry implements a simple exponential backoff retry helper for
// depwatch. It is used by the checker and notifier packages to make outbound
// HTTP calls more resilient against transient network errors.
//
// Usage:
//
//	cfg := retry.DefaultConfig() // 3 attempts, 500 ms base delay, 10 s cap
//	err := retry.Do(ctx, cfg, func() error {
//	    return doSomethingFallible()
//	})
//	if errors.Is(err, retry.ErrMaxAttempts) {
//	    // all attempts exhausted
//	}
package retry
