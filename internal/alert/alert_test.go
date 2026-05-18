package alert_test

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/depwatch/internal/alert"
)

func TestNew_FieldsPopulated(t *testing.T) {
	before := time.Now().UTC()
	p := alert.New("github.com/foo/bar", "v1.0.0", "v1.1.0", alert.SeverityMinor)
	after := time.Now().UTC()

	if p.Module != "github.com/foo/bar" {
		t.Errorf("module: got %q", p.Module)
	}
	if p.OldVersion != "v1.0.0" {
		t.Errorf("old_version: got %q", p.OldVersion)
	}
	if p.NewVersion != "v1.1.0" {
		t.Errorf("new_version: got %q", p.NewVersion)
	}
	if p.Severity != alert.SeverityMinor {
		t.Errorf("severity: got %q", p.Severity)
	}
	if p.DetectedAt.Before(before) || p.DetectedAt.After(after) {
		t.Errorf("detected_at out of range: %v", p.DetectedAt)
	}
}

func TestEncode_ValidJSON(t *testing.T) {
	p := alert.New("github.com/foo/bar", "v1.0.0", "v2.0.0", alert.SeverityMajor)
	buf, err := p.Encode()
	if err != nil {
		t.Fatalf("Encode error: %v", err)
	}
	var out alert.Payload
	if err := json.NewDecoder(buf).Decode(&out); err != nil {
		t.Fatalf("decode error: %v", err)
	}
	if out.Module != p.Module {
		t.Errorf("round-trip module: got %q want %q", out.Module, p.Module)
	}
	if out.Severity != alert.SeverityMajor {
		t.Errorf("round-trip severity: got %q", out.Severity)
	}
}

func TestInferSeverity(t *testing.T) {
	cases := []struct {
		old, new string
		want     alert.Severity
	}{
		{"", "v1.0.0", alert.SeverityInfo},
		{"v1.0.0", "v1.0.1", alert.SeverityPatch},
		{"v1.0.0", "v1.1.0", alert.SeverityMinor},
		{"v1.2.3", "v2.0.0", alert.SeverityMajor},
		{"v2.1.0", "v2.1.5", alert.SeverityPatch},
	}
	for _, tc := range cases {
		got := alert.InferSeverity(tc.old, tc.new)
		if got != tc.want {
			t.Errorf("InferSeverity(%q, %q) = %q; want %q", tc.old, tc.new, got, tc.want)
		}
	}
}
