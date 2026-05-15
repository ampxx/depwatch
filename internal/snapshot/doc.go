// Package snapshot provides a thread-safe, point-in-time record of module
// versions observed during a poll cycle.
//
// A Snapshot can be compared against a subsequent Snapshot to produce a list
// of Diff values, each describing a module whose version has changed (or that
// has appeared for the first time). Snapshots are immutable once cloned, so
// callers can safely retain a reference across goroutines without additional
// synchronisation.
package snapshot
