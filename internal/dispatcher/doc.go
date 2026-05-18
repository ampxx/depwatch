// Package dispatcher provides a fan-out mechanism that routes alert.Alert
// values to one or more Sender targets.
//
// Each Target can be independently configured with:
//   - A filter.Filter to restrict which modules trigger notifications.
//   - A throttle.Throttle to suppress repeated alerts within a time window.
//   - A MinSeverity threshold below which alerts are silently dropped.
//
// Usage:
//
//	d := dispatcher.New()
//	d.Add(&dispatcher.Target{
//	    Name:        "slack",
//	    Sender:      slackNotifier,
//	    MinSeverity: alert.SeverityMedium,
//	})
//	d.Dispatch(ctx, myAlert)
package dispatcher
