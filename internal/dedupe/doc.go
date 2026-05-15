// Package dedupe implements alert deduplication for depwatch.
//
// A Deduplicator tracks module@version pairs that have recently triggered
// notifications. When the same pair is seen again within the configured
// window the alert is suppressed, preventing webhook spam during polling
// cycles where the upstream version has not changed.
//
// Usage:
//
//	dd := dedupe.New(5 * time.Minute)
//	if !dd.IsDuplicate(module, version) {
//		// send alert
//	}
//
// Call Purge periodically (e.g. once per poll cycle) to reclaim memory.
package dedupe
