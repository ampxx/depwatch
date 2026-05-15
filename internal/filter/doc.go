// Package filter provides include/exclude pattern matching for Go module
// paths monitored by depwatch.
//
// Patterns follow the same glob syntax as [path.Match]. Exclusion patterns
// always take precedence over inclusion patterns, mirroring the behaviour
// of common CI ignore files.
//
// Usage:
//
//	f := filter.New(cfg.IncludePatterns, cfg.ExcludePatterns)
//	if f.Allow(modulePath) {
//	    // proceed with version check
//	}
package filter
