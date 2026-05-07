// Package notifier provides webhook notification functionality for depwatch.
//
// It serializes dependency version change events into JSON payloads and
// delivers them to a configured webhook URL via HTTP POST requests.
//
// Example usage:
//
//	n := notifier.New("https://hooks.example.com/depwatch")
//	if err := n.Notify("github.com/some/dep", "v1.2.3", "v1.3.0"); err != nil {
//		log.Printf("failed to send notification: %v", err)
//	}
package notifier
