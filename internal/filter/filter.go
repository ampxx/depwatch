// Package filter provides dependency filtering logic for depwatch,
// allowing users to include or exclude specific modules from monitoring
// based on glob-style patterns.
package filter

import (
	"path"
	"strings"
)

// Filter decides which module paths should be watched.
type Filter struct {
	include []string
	exclude []string
}

// New creates a Filter from include and exclude pattern slices.
// If include is empty every module is considered included by default.
func New(include, exclude []string) *Filter {
	return &Filter{
		include: include,
		exclude: exclude,
	}
}

// Allow returns true when modulePath should be monitored.
// Exclusion patterns take precedence over inclusion patterns.
func (f *Filter) Allow(modulePath string) bool {
	for _, pat := range f.exclude {
		if matchPattern(pat, modulePath) {
			return false
		}
	}

	if len(f.include) == 0 {
		return true
	}

	for _, pat := range f.include {
		if matchPattern(pat, modulePath) {
			return true
		}
	}

	return false
}

// matchPattern matches a glob pattern against a module path.
// It normalises both strings to forward-slash paths before matching.
func matchPattern(pattern, modulePath string) bool {
	pattern = strings.TrimSpace(pattern)
	modulePath = strings.TrimSpace(modulePath)
	matched, err := path.Match(pattern, modulePath)
	if err != nil {
		// Treat malformed patterns as non-matching.
		return false
	}
	return matched
}
