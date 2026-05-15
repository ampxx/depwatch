package filter_test

import (
	"testing"

	"github.com/yourorg/depwatch/internal/filter"
)

func TestAllow_NoPatterns_AllowsAll(t *testing.T) {
	f := filter.New(nil, nil)
	if !f.Allow("github.com/some/module") {
		t.Error("expected module to be allowed when no patterns are set")
	}
}

func TestAllow_IncludePattern_MatchesExact(t *testing.T) {
	f := filter.New([]string{"github.com/foo/bar"}, nil)
	if !f.Allow("github.com/foo/bar") {
		t.Error("expected exact include match to be allowed")
	}
	if f.Allow("github.com/foo/baz") {
		t.Error("expected non-matching module to be denied")
	}
}

func TestAllow_IncludeGlob_MatchesWildcard(t *testing.T) {
	f := filter.New([]string{"github.com/foo/*"}, nil)
	if !f.Allow("github.com/foo/bar") {
		t.Error("expected wildcard include to match")
	}
	if f.Allow("github.com/other/bar") {
		t.Error("expected non-matching module to be denied")
	}
}

func TestAllow_ExcludePattern_DeniesMatch(t *testing.T) {
	f := filter.New(nil, []string{"github.com/bad/module"})
	if f.Allow("github.com/bad/module") {
		t.Error("expected excluded module to be denied")
	}
	if !f.Allow("github.com/good/module") {
		t.Error("expected non-excluded module to be allowed")
	}
}

func TestAllow_ExcludeTakesPrecedenceOverInclude(t *testing.T) {
	f := filter.New(
		[]string{"github.com/foo/*"},
		[]string{"github.com/foo/bar"},
	)
	if f.Allow("github.com/foo/bar") {
		t.Error("expected exclude to take precedence over include")
	}
	if !f.Allow("github.com/foo/baz") {
		t.Error("expected non-excluded include match to be allowed")
	}
}

func TestAllow_MalformedPattern_TreatedAsNonMatch(t *testing.T) {
	// '[' is an invalid glob pattern in path.Match.
	f := filter.New([]string{"[invalid"}, nil)
	if f.Allow("github.com/anything") {
		t.Error("expected malformed include pattern to not match")
	}
}
