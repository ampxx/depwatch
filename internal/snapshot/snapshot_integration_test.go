package snapshot_test

import (
	"sync"
	"testing"

	"github.com/depwatch/internal/snapshot"
)

// TestSnapshot_ConcurrentWrites verifies that concurrent Set calls do not race.
func TestSnapshot_ConcurrentWrites(t *testing.T) {
	s := snapshot.New()
	var wg sync.WaitGroup
	modules := []string{
		"github.com/a/a",
		"github.com/b/b",
		"github.com/c/c",
		"github.com/d/d",
	}
	for _, mod := range modules {
		wg.Add(1)
		go func(m string) {
			defer wg.Done()
			s.Set(m, "v1.0.0")
		}(mod)
	}
	wg.Wait()

	for _, mod := range modules {
		v, ok := s.Get(mod)
		if !ok || v != "v1.0.0" {
			t.Errorf("module %s not set correctly", mod)
		}
	}
}

// TestSnapshot_FullCycle simulates a complete poll-cycle comparison workflow.
func TestSnapshot_FullCycle(t *testing.T) {
	prev := snapshot.New()
	prev.Set("github.com/unchanged/mod", "v1.0.0")
	prev.Set("github.com/updated/mod", "v1.0.0")

	curr := snapshot.New()
	curr.Set("github.com/unchanged/mod", "v1.0.0")
	curr.Set("github.com/updated/mod", "v1.1.0")
	curr.Set("github.com/new/mod", "v0.1.0")

	diffs := prev.Compare(curr)
	if len(diffs) != 2 {
		t.Fatalf("expected 2 diffs, got %d", len(diffs))
	}

	byModule := make(map[string]snapshot.Diff, len(diffs))
	for _, d := range diffs {
		byModule[d.Module] = d
	}

	if d, ok := byModule["github.com/updated/mod"]; !ok || d.NewVersion != "v1.1.0" {
		t.Error("expected updated/mod to appear in diffs")
	}
	if d, ok := byModule["github.com/new/mod"]; !ok || d.OldVersion != "" {
		t.Error("expected new/mod with empty old version")
	}
	if _, ok := byModule["github.com/unchanged/mod"]; ok {
		t.Error("unchanged/mod should not appear in diffs")
	}
}
