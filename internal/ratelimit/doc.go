// Package ratelimit implements a thread-safe token-bucket rate limiter used
// to throttle outbound webhook notifications in depwatch.
//
// A Limiter is created with a maximum burst size and a sustained refill rate
// (tokens per second). Each call to Allow consumes one token; when the bucket
// is empty Allow returns false and the caller should skip or defer the
// notification until tokens are available again.
//
// Example usage:
//
//	limiter := ratelimit.New(5, 0.5) // burst of 5, refill 1 token per 2 s
//	if limiter.Allow() {
//		notifier.Notify(ctx, module, oldVer, newVer)
//	}
package ratelimit
