// Package backoff provides an exponential backoff Strategy used to compute
// wait durations between successive retry attempts.
//
// Usage:
//
//	s := backoff.New(backoff.DefaultConfig())
//	for attempt := 0; attempt < maxAttempts; attempt++ {
//		err := doSomething()
//		if err == nil {
//			break
//		}
//		time.Sleep(s.Next(attempt))
//	}
//
// The default configuration uses a 200 ms base, a 2× multiplier, a 30 s cap,
// and adds up to 30 % random jitter to spread load across concurrent callers.
package backoff
