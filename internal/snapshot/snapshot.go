// Package snapshot records and compares the state of all watched
// dependencies at a point in time, enabling diff-based change detection
// across poll cycles.
package snapshot

import (
	"fmt"
	"sort"
	"strings"
	"sync"
)

// Entry holds the last-known version for a single module.
type Entry struct {
	Module  string
	Version string
}

// Diff describes a version change between two snapshots.
type Diff struct {
	Module  string
	OldVersion string
	NewVersion string
}

func (d Diff) String() string {
	return fmt.Sprintf("%s: %s → %s", d.Module, d.OldVersion, d.NewVersion)
}

// Snapshot is an immutable, point-in-time view of module versions.
type Snapshot struct {
	mu      sync.RWMutex
	entries map[string]string
}

// New returns an empty Snapshot.
func New() *Snapshot {
	return &Snapshot{entries: make(map[string]string)}
}

// Set records or updates the version for the given module.
func (s *Snapshot) Set(module, version string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries[module] = version
}

// Get returns the stored version and whether it exists.
func (s *Snapshot) Get(module string) (string, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	v, ok := s.entries[module]
	return v, ok
}

// Compare returns the list of Diffs between s (old) and next (new).
// Only modules present in next are considered.
func (s *Snapshot) Compare(next *Snapshot) []Diff {
	s.mu.RLock()
	defer s.mu.RUnlock()
	next.mu.RLock()
	defer next.mu.RUnlock()

	var diffs []Diff
	for mod, newVer := range next.entries {
		oldVer, exists := s.entries[mod]
		if !exists || oldVer != newVer {
			diffs = append(diffs, Diff{Module: mod, OldVersion: oldVer, NewVersion: newVer})
		}
	}
	sort.Slice(diffs, func(i, j int) bool {
		return strings.Compare(diffs[i].Module, diffs[j].Module) < 0
	})
	return diffs
}

// Clone returns a deep copy of the snapshot.
func (s *Snapshot) Clone() *Snapshot {
	s.mu.RLock()
	defer s.mu.RUnlock()
	c := New()
	for k, v := range s.entries {
		c.entries[k] = v
	}
	return c
}
