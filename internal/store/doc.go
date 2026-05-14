// Package store provides a simple persistent key-value store for tracking
// the last known versions of Go module dependencies.
//
// The store serializes version data to a JSON file on disk, allowing depwatch
// to detect changes across restarts. It is safe for concurrent use via
// internal locking.
//
// Example usage:
//
//	s, err := store.New("/var/lib/depwatch/versions.json")
//	if err != nil {
//		log.Fatal(err)
//	}
//	s.Set("github.com/some/module", "v1.2.3")
//	version, ok := s.Get("github.com/some/module")
package store
