// Package semver provides utilities for parsing and comparing semantic
// version strings as used by the Go module proxy.
package semver

import (
	"fmt"
	"strconv"
	"strings"
)

// Version represents a parsed semantic version.
type Version struct {
	Major int
	Minor int
	Patch int
	Pre   string // pre-release label, e.g. "beta.1"
}

// Parse parses a semver string such as "v1.2.3" or "v2.0.0-beta.1".
// The leading "v" is optional.
func Parse(s string) (Version, error) {
	s = strings.TrimPrefix(s, "v")

	pre := ""
	if idx := strings.IndexByte(s, '-'); idx >= 0 {
		pre = s[idx+1:]
		s = s[:idx]
	}

	parts := strings.SplitN(s, ".", 3)
	if len(parts) != 3 {
		return Version{}, fmt.Errorf("semver: invalid version %q", s)
	}

	parse := func(p string) (int, error) {
		n, err := strconv.Atoi(p)
		if err != nil {
			return 0, fmt.Errorf("semver: non-numeric segment %q", p)
		}
		return n, nil
	}

	major, err := parse(parts[0])
	if err != nil {
		return Version{}, err
	}
	minor, err := parse(parts[1])
	if err != nil {
		return Version{}, err
	}
	patch, err := parse(parts[2])
	if err != nil {
		return Version{}, err
	}

	return Version{Major: major, Minor: minor, Patch: patch, Pre: pre}, nil
}

// IsNewer reports whether v is strictly newer than other.
// Pre-release versions are considered older than the equivalent stable release.
func (v Version) IsNewer(other Version) bool {
	switch {
	case v.Major != other.Major:
		return v.Major > other.Major
	case v.Minor != other.Minor:
		return v.Minor > other.Minor
	case v.Patch != other.Patch:
		return v.Patch > other.Patch
	}
	// Same numeric version: stable (empty pre) beats pre-release.
	if v.Pre == "" && other.Pre != "" {
		return true
	}
	return false
}

// String returns the canonical "vMAJOR.MINOR.PATCH" representation.
func (v Version) String() string {
	s := fmt.Sprintf("v%d.%d.%d", v.Major, v.Minor, v.Patch)
	if v.Pre != "" {
		s += "-" + v.Pre
	}
	return s
}
