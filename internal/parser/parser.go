// Package parser provides utilities for parsing go.mod files and
// extracting module dependency information.
package parser

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"strings"
)

// Dependency represents a single module dependency extracted from a go.mod file.
type Dependency struct {
	Path    string
	Version string
}

// ParseFile reads the go.mod file at the given path and returns all
// require directives as a slice of Dependency.
func ParseFile(path string) ([]Dependency, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("parser: open %q: %w", path, err)
	}
	defer f.Close()
	return Parse(f)
}

// Parse reads go.mod content from r and returns all require directives
// as a slice of Dependency.
func Parse(r io.Reader) ([]Dependency, error) {
	var deps []Dependency
	inBlock := false

	scanner := bufio.NewScanner(r)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		if line == "require (" {
			inBlock = true
			continue
		}
		if inBlock && line == ")" {
			inBlock = false
			continue
		}

		var raw string
		if inBlock {
			raw = line
		} else if strings.HasPrefix(line, "require ") {
			raw = strings.TrimPrefix(line, "require ")
		} else {
			continue
		}

		// Strip inline comments.
		if idx := strings.Index(raw, "//"); idx != -1 {
			raw = strings.TrimSpace(raw[:idx])
		}
		if raw == "" {
			continue
		}

		parts := strings.Fields(raw)
		if len(parts) < 2 {
			continue
		}
		deps = append(deps, Dependency{Path: parts[0], Version: parts[1]})
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parser: scan: %w", err)
	}
	return deps, nil
}
