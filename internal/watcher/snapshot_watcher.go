package watcher

import (
	"context"
	"log"
	"time"

	"github.com/depwatch/internal/snapshot"
)

// SnapshotWatcher polls all configured modules and emits Diffs by comparing
// successive Snapshots. It integrates with the existing Runner infrastructure
// without modifying it.
type SnapshotWatcher struct {
	client   VersionChecker
	modules  []string
	interval time.Duration
	prev     *snapshot.Snapshot
}

// VersionChecker is satisfied by checker.Client and its decorators.
type VersionChecker interface {
	LatestVersion(ctx context.Context, module string) (string, error)
}

// NewSnapshotWatcher constructs a SnapshotWatcher for the given modules.
func NewSnapshotWatcher(client VersionChecker, modules []string, interval time.Duration) *SnapshotWatcher {
	return &SnapshotWatcher{
		client:   client,
		modules:  modules,
		interval: interval,
		prev:     snapshot.New(),
	}
}

// Run polls on every interval tick until ctx is cancelled.
// Detected diffs are passed to onChange.
func (sw *SnapshotWatcher) Run(ctx context.Context, onChange func([]snapshot.Diff)) {
	ticker := time.NewTicker(sw.interval)
	defer ticker.Stop()

	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			diffs := sw.poll(ctx)
			if len(diffs) > 0 {
				onChange(diffs)
			}
		}
	}
}

func (sw *SnapshotWatcher) poll(ctx context.Context) []snapshot.Diff {
	curr := snapshot.New()
	for _, mod := range sw.modules {
		v, err := sw.client.LatestVersion(ctx, mod)
		if err != nil {
			log.Printf("snapshot_watcher: error fetching %s: %v", mod, err)
			continue
		}
		curr.Set(mod, v)
	}
	diffs := sw.prev.Compare(curr)
	sw.prev = curr.Clone()
	return diffs
}
