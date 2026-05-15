package snapshot_test

import (
	"testing"

	"github.com/depwatch/internal/snapshot"
)

func TestSet_AndGet(t *testing.T) {
	s := snapshot.New()
	s.Set("github.com/foo/bar", "v1.2.3")
	v, ok := s.Get("github.com/foo/bar")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if v != "v1.2.3" {
		t.Fatalf("expected v1.2.3, got %s", v)
	}
}

func TestGet_Missing(t *testing.T) {
	s := snapshot.New()
	_, ok := s.Get("github.com/missing/mod")
	if ok {
		t.Fatal("expected missing entry")
	}
}

func TestCompare_NoDiff(t *testing.T) {
	old := snapshot.New()
	old.Set("github.com/a/b", "v1.0.0")

	next := snapshot.New()
	next.Set("github.com/a/b", "v1.0.0")

	diffs := old.Compare(next)
	if len(diffs) != 0 {
		t.Fatalf("expected no diffs, got %d", len(diffs))
	}
}

func TestCompare_DetectsChange(t *testing.T) {
	old := snapshot.New()
	old.Set("github.com/a/b", "v1.0.0")

	next := snapshot.New()
	next.Set("github.com/a/b", "v1.1.0")

	diffs := old.Compare(next)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].OldVersion != "v1.0.0" || diffs[0].NewVersion != "v1.1.0" {
		t.Fatalf("unexpected diff: %s", diffs[0])
	}
}

func TestCompare_NewModuleCountsAsDiff(t *testing.T) {
	old := snapshot.New()
	next := snapshot.New()
	next.Set("github.com/new/mod", "v0.1.0")

	diffs := old.Compare(next)
	if len(diffs) != 1 {
		t.Fatalf("expected 1 diff, got %d", len(diffs))
	}
	if diffs[0].OldVersion != "" {
		t.Fatalf("expected empty old version for new module")
	}
}

func TestClone_IsIndependent(t *testing.T) {
	s := snapshot.New()
	s.Set("github.com/a/b", "v1.0.0")
	c := s.Clone()
	s.Set("github.com/a/b", "v2.0.0")

	v, _ := c.Get("github.com/a/b")
	if v != "v1.0.0" {
		t.Fatalf("clone should not reflect mutation, got %s", v)
	}
}

func TestDiff_String(t *testing.T) {
	d := snapshot.Diff{Module: "github.com/x/y", OldVersion: "v1.0.0", NewVersion: "v1.1.0"}
	expected := "github.com/x/y: v1.0.0 → v1.1.0"
	if d.String() != expected {
		t.Fatalf("expected %q, got %q", expected, d.String())
	}
}
