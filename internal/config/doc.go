// Package config handles loading and validating depwatch configuration from
// a YAML file.
//
// Configuration includes the path to the Go module file to watch, the webhook
// URL to notify on version changes, an optional poll interval (defaulting to
// 60 seconds), and an optional path for the persistent version store.
//
// Example usage:
//
//	cfg, err := config.Load("config.yaml")
//	if err != nil {
//		log.Fatal(err)
//	}
//	fmt.Println(cfg.WebhookURL)
package config
