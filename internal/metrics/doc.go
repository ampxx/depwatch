// Package metrics provides lightweight, goroutine-safe counters for
// tracking depwatch runtime statistics such as the total number of
// dependency checks performed, version changes detected, webhook
// notifications sent or failed, and persistent-store errors.
//
// All counters are backed by [sync/atomic] integers and are safe to
// update concurrently from multiple goroutines without additional
// locking.
//
// Use [New] to create a Metrics instance, increment the relevant
// counter fields directly, and call [Metrics.Snapshot] to obtain a
// point-in-time copy suitable for logging or JSON serialisation.
package metrics
