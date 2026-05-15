// Package changelog provides a persistent, append-only log of dependency
// version change events detected by depwatch.
//
// Each time a monitored module transitions from one version to another, a
// timestamped Entry is appended to the Log and flushed to a JSON file on
// disk. The log can be reloaded across daemon restarts and queried by module
// path for audit or reporting purposes.
package changelog
