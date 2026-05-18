// Package alert provides structured alert payloads sent to webhook endpoints
// when a dependency version change is detected.
package alert

import (
	"bytes"
	"encoding/json"
	"fmt"
	"time"
)

// Severity indicates the urgency level of an alert.
type Severity string

const (
	SeverityInfo  Severity = "info"
	SeverityWarn  Severity = "warn"
	SeverityPatch Severity = "patch"
	SeverityMinor Severity = "minor"
	SeverityMajor Severity = "major"
)

// Payload is the JSON body posted to a webhook when a version change is detected.
type Payload struct {
	Module     string    `json:"module"`
	OldVersion string    `json:"old_version"`
	NewVersion string    `json:"new_version"`
	Severity   Severity  `json:"severity"`
	DetectedAt time.Time `json:"detected_at"`
}

// New constructs a Payload, deriving severity from the version bump type.
func New(module, oldVersion, newVersion string, severity Severity) Payload {
	return Payload{
		Module:     module,
		OldVersion: oldVersion,
		NewVersion: newVersion,
		Severity:   severity,
		DetectedAt: time.Now().UTC(),
	}
}

// Encode serialises the payload to a compact JSON byte buffer.
func (p Payload) Encode() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	if err := json.NewEncoder(&buf).Encode(p); err != nil {
		return nil, fmt.Errorf("alert: encode payload: %w", err)
	}
	return &buf, nil
}

// InferSeverity returns a Severity based on which semver component changed.
// It expects canonical semver strings prefixed with 'v' (e.g. "v1.2.3").
func InferSeverity(oldVersion, newVersion string) Severity {
	if oldVersion == "" {
		return SeverityInfo
	}
	var oldMajor, oldMinor int
	var newMajor, newMinor int
	fmt.Sscanf(oldVersion, "v%d.%d.", &oldMajor, &oldMinor)
	fmt.Sscanf(newVersion, "v%d.%d.", &newMajor, &newMinor)
	switch {
	case newMajor > oldMajor:
		return SeverityMajor
	case newMinor > oldMinor:
		return SeverityMinor
	default:
		return SeverityPatch
	}
}
