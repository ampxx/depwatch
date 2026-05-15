package changelog_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/depwatch/internal/changelog"
)

func tempLogPath(t *testing.T) string {
	t.Helper()
	return filepath.Join(t.TempDir(), "changelog.json")
}

func TestRecord_AppendsEntry(t *testing.T) {
	log, err := changelog.New(tempLogPath(t))
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	if err := log.Record("github.com/foo/bar", "v1.0.0", "v1.1.0"); err != nil {
		t.Fatalf("Record: %v", err)
	}
	entries := log.Entries()
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	e := entries[0]
	if e.Module != "github.com/foo/bar" || e.OldVer != "v1.0.0" || e.NewVer != "v1.1.0" {
		t.Errorf("unexpected entry: %+v", e)
	}
	if e.DetectedAt.IsZero() {
		t.Error("DetectedAt should not be zero")
	}
}

func TestEntriesFor_FiltersCorrectly(t *testing.T) {
	log, _ := changelog.New(tempLogPath(t))
	_ = log.Record("github.com/foo/bar", "v1.0.0", "v1.1.0")
	_ = log.Record("github.com/baz/qux", "v2.0.0", "v2.1.0")
	_ = log.Record("github.com/foo/bar", "v1.1.0", "v1.2.0")

	got := log.EntriesFor("github.com/foo/bar")
	if len(got) != 2 {
		t.Fatalf("expected 2 entries for foo/bar, got %d", len(got))
	}
}

func TestLog_PersistsAcrossReload(t *testing.T) {
	path := tempLogPath(t)
	log1, _ := changelog.New(path)
	_ = log1.Record("github.com/foo/bar", "v1.0.0", "v1.1.0")

	log2, err := changelog.New(path)
	if err != nil {
		t.Fatalf("reload: %v", err)
	}
	if len(log2.Entries()) != 1 {
		t.Errorf("expected 1 persisted entry, got %d", len(log2.Entries()))
	}
}

func TestNew_InvalidFile(t *testing.T) {
	path := filepath.Join(t.TempDir(), "bad.json")
	if err := os.WriteFile(path, []byte("not json{"), 0o644); err != nil {
		t.Fatal(err)
	}
	_, err := changelog.New(path)
	if err == nil {
		t.Error("expected error for invalid JSON, got nil")
	}
}

func TestEntries_ReturnsCopy(t *testing.T) {
	log, _ := changelog.New(tempLogPath(t))
	_ = log.Record("github.com/a/b", "v1.0.0", "v1.0.1")
	copy1 := log.Entries()
	copy1[0].Module = "mutated"
	copy2 := log.Entries()
	if copy2[0].Module == "mutated" {
		t.Error("Entries should return an independent copy")
	}
}
