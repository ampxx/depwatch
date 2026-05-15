// Package parser implements a lightweight go.mod parser for depwatch.
//
// It extracts require directives — both block-form and single-line — from a
// go.mod file and returns them as a slice of Dependency values that the
// watcher and checker packages can act upon.
//
// Inline comments (// ...) are stripped before processing, and indirect
// dependencies are included alongside direct ones so that the full
// dependency graph is monitored.
package parser
