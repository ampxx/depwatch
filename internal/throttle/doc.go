// Package throttle implements a per-module notification cooldown mechanism
// for depwatch.
//
// When a dependency version change is detected, the watcher may fire
// repeatedly during a polling window. Throttle ensures that webhook
// notifications for the same module are suppressed until the configured
// cooldown duration has elapsed since the last successful send.
//
// Usage:
//
//	th := throttle.New(10 * time.Minute)
//	if th.Allow(moduleName) {
//		// send notification
//	}
package throttle
