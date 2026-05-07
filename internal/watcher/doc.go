// Package watcher implements the core polling loop for depwatch.
//
// A Watcher periodically queries the Go module proxy for the latest
// version of each configured module. When a version change is detected
// it persists the new version via the store and dispatches an alert
// through the notifier.
//
// Typical usage:
//
//	w := watcher.New(cfg, checkerClient, notifierClient, storeClient)
//	w.Run(ctx) // blocks until ctx is cancelled
package watcher
