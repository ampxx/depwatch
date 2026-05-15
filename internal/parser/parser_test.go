package parser_test

import (
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/depwatch/internal/parser"
)

const sampleGoMod = `module github.com/example/app

go 1.21

require (
	github.com/foo/bar v1.2.3
	github.com/baz/qux v0.9.0 // indirect
)

require github.com/single/dep v2.0.0
`

func writeTempGoMod(t *testing.T, content string) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "go.mod")
	if err := os.WriteFile(p, []byte(content), 0o644); err != nil {
		t.Fatalf("write temp go.mod: %v", err)
	}
	return p
}

func TestParse_BlockRequire(t *testing.T) {
	deps, err := parser.Parse(strings.NewReader(sampleGoMod))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %d", len(deps))
	}
	assertDep(t, deps[0], "github.com/foo/bar", "v1.2.3")
	assertDep(t, deps[1], "github.com/baz/qux", "v0.9.0")
	assertDep(t, deps[2], "github.com/single/dep", "v2.0.0")
}

func TestParse_EmptyFile(t *testing.T) {
	deps, err := parser.Parse(strings.NewReader("module github.com/x/y\n\ngo 1.21\n"))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 0 {
		t.Fatalf("expected 0 deps, got %d", len(deps))
	}
}

func TestParseFile_Success(t *testing.T) {
	p := writeTempGoMod(t, sampleGoMod)
	deps, err := parser.ParseFile(p)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 3 {
		t.Fatalf("expected 3 deps, got %d", len(deps))
	}
}

func TestParseFile_Missing(t *testing.T) {
	_, err := parser.ParseFile("/nonexistent/go.mod")
	if err == nil {
		t.Fatal("expected error for missing file, got nil")
	}
}

func TestParse_InlineCommentStripped(t *testing.T) {
	input := "require github.com/a/b v1.0.0 // some comment\n"
	deps, err := parser.Parse(strings.NewReader(input))
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(deps) != 1 {
		t.Fatalf("expected 1 dep, got %d", len(deps))
	}
	assertDep(t, deps[0], "github.com/a/b", "v1.0.0")
}

func assertDep(t *testing.T, d parser.Dependency, wantPath, wantVersion string) {
	t.Helper()
	if d.Path != wantPath {
		t.Errorf("path: got %q, want %q", d.Path, wantPath)
	}
	if d.Version != wantVersion {
		t.Errorf("version: got %q, want %q", d.Version, wantVersion)
	}
}
