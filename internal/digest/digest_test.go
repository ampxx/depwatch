package digest_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/yourorg/depwatch/internal/digest"
)

func writeTempFile(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "go.sum")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("writeTempFile: %v", err)
	}
	return p
}

func TestChanged_FirstCallReturnsFalse(t *testing.T) {
	p := writeTempFile(t, "module github.com/foo/bar v1.0.0")
	d := digest.New()

	changed, err := d.Changed(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected false on first call, got true")
	}
}

func TestChanged_SameContentReturnsFalse(t *testing.T) {
	p := writeTempFile(t, "module github.com/foo/bar v1.0.0")
	d := digest.New()

	if _, err := d.Changed(p); err != nil {
		t.Fatalf("seed call failed: %v", err)
	}
	changed, err := d.Changed(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if changed {
		t.Error("expected false for identical content, got true")
	}
}

func TestChanged_ModifiedContentReturnsTrue(t *testing.T) {
	p := writeTempFile(t, "module github.com/foo/bar v1.0.0")
	d := digest.New()

	if _, err := d.Changed(p); err != nil {
		t.Fatalf("seed call failed: %v", err)
	}

	if err := os.WriteFile(p, []byte("module github.com/foo/bar v1.1.0"), 0o644); err != nil {
		t.Fatalf("overwrite failed: %v", err)
	}

	changed, err := d.Changed(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if !changed {
		t.Error("expected true after content change, got false")
	}
}

func TestReset_TreatsPathAsUnseen(t *testing.T) {
	p := writeTempFile(t, "module github.com/foo/bar v1.0.0")
	d := digest.New()

	if _, err := d.Changed(p); err != nil {
		t.Fatalf("seed call failed: %v", err)
	}
	d.Reset(p)

	changed, err := d.Changed(p)
	if err != nil {
		t.Fatalf("unexpected error after reset: %v", err)
	}
	if changed {
		t.Error("expected false after Reset (treated as first call), got true")
	}
}

func TestHashFile_MissingFile(t *testing.T) {
	_, err := digest.HashFile("/nonexistent/path/go.sum")
	if err == nil {
		t.Error("expected error for missing file, got nil")
	}
}
