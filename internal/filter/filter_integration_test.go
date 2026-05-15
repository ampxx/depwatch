package filter_test

import (
	"testing"

	"github.com/yourorg/depwatch/internal/filter"
)

// TestFilter_TypicalWorkflow exercises a realistic configuration where an
// organisation wants to watch all internal modules but ignore test stubs.
func TestFilter_TypicalWorkflow(t *testing.T) {
	f := filter.New(
		[]string{
			"github.com/myorg/*",
			"golang.org/x/*",
		},
		[]string{
			"github.com/myorg/testutil",
		},
	)

	cases := []struct {
		module  string
		wantOK  bool
	}{
		{"github.com/myorg/api", true},
		{"github.com/myorg/worker", true},
		{"github.com/myorg/testutil", false}, // excluded
		{"golang.org/x/sync", true},
		{"golang.org/x/net", true},
		{"github.com/thirdparty/lib", false}, // not included
		{"github.com/myorg/deep/nested", false}, // single-star doesn't cross slashes
	}

	for _, tc := range cases {
		t.Run(tc.module, func(t *testing.T) {
			got := f.Allow(tc.module)
			if got != tc.wantOK {
				t.Errorf("Allow(%q) = %v, want %v", tc.module, got, tc.wantOK)
			}
		})
	}
}

// TestFilter_EmptyExclude_IncludeStillApplied ensures that an empty exclude
// list does not inadvertently block all modules.
func TestFilter_EmptyExclude_IncludeStillApplied(t *testing.T) {
	f := filter.New([]string{"github.com/safe/*"}, []string{})
	if !f.Allow("github.com/safe/pkg") {
		t.Error("expected included module to be allowed with empty exclude list")
	}
	if f.Allow("github.com/unsafe/pkg") {
		t.Error("expected non-included module to be denied")
	}
}
