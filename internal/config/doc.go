// Package config handles loading and validating depwatch configuration
// from a YAML file.
//
// Configuration includes the path to the Go module file being monitored,
// the webhook URL to notify on version changes, and an optional poll
// interval that controls how frequently the proxy is queried.
//
// Example usage:
//
//	cfg, err := config.Load("config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(cfg.WebhookURL)
package config
