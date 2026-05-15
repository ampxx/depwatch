package semver_test

import (
	"testing"

	"github.com/yourorg/depwatch/internal/semver"
)

func TestParse_Valid(t *testing.T) {
	v, err := semver.Parse("v1.2.3")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 1 || v.Minor != 2 || v.Patch != 3 {
		t.Errorf("got %v, want v1.2.3", v)
	}
}

func TestParse_WithoutLeadingV(t *testing.T) {
	v, err := semver.Parse("0.9.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Major != 0 || v.Minor != 9 || v.Patch != 1 {
		t.Errorf("got %v", v)
	}
}

func TestParse_PreRelease(t *testing.T) {
	v, err := semver.Parse("v2.0.0-beta.1")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if v.Pre != "beta.1" {
		t.Errorf("expected pre=beta.1, got %q", v.Pre)
	}
}

func TestParse_Invalid(t *testing.T) {
	cases := []string{"v1.2", "v1.2.x", "", "latest"}
	for _, c := range cases {
		_, err := semver.Parse(c)
		if err == nil {
			t.Errorf("expected error for %q, got nil", c)
		}
	}
}

func TestIsNewer_PatchBump(t *testing.T) {
	newer, _ := semver.Parse("v1.2.4")
	older, _ := semver.Parse("v1.2.3")
	if !newer.IsNewer(older) {
		t.Error("v1.2.4 should be newer than v1.2.3")
	}
	if older.IsNewer(newer) {
		t.Error("v1.2.3 should not be newer than v1.2.4")
	}
}

func TestIsNewer_MajorBump(t *testing.T) {
	v2, _ := semver.Parse("v2.0.0")
	v1, _ := semver.Parse("v1.99.99")
	if !v2.IsNewer(v1) {
		t.Error("v2.0.0 should be newer than v1.99.99")
	}
}

func TestIsNewer_StableBeatsPreRelease(t *testing.T) {
	stable, _ := semver.Parse("v1.0.0")
	pre, _ := semver.Parse("v1.0.0-rc.1")
	if !stable.IsNewer(pre) {
		t.Error("v1.0.0 should be newer than v1.0.0-rc.1")
	}
	if pre.IsNewer(stable) {
		t.Error("pre-release should not be newer than stable")
	}
}

func TestIsNewer_SameVersion(t *testing.T) {
	a, _ := semver.Parse("v1.2.3")
	b, _ := semver.Parse("v1.2.3")
	if a.IsNewer(b) || b.IsNewer(a) {
		t.Error("equal versions should not report IsNewer")
	}
}

func TestString_RoundTrip(t *testing.T) {
	cases := []string{"v1.2.3", "v0.0.1-alpha"}
	for _, c := range cases {
		v, err := semver.Parse(c)
		if err != nil {
			t.Fatalf("parse %q: %v", c, err)
		}
		if got := v.String(); got != c {
			t.Errorf("String() = %q, want %q", got, c)
		}
	}
}
