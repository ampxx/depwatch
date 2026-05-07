package notifier

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// Payload represents the webhook notification body.
type Payload struct {
	Module  string `json:"module"`
	OldVersion string `json:"old_version"`
	NewVersion string `json:"new_version"`
	Timestamp  string `json:"timestamp"`
}

// Notifier sends webhook notifications.
type Notifier struct {
	webhookURL string
	client     *http.Client
}

// New creates a new Notifier with the given webhook URL.
func New(webhookURL string) *Notifier {
	return &Notifier{
		webhookURL: webhookURL,
		client: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// Notify sends a webhook notification about a dependency version change.
func (n *Notifier) Notify(module, oldVersion, newVersion string) error {
	payload := Payload{
		Module:     module,
		OldVersion: oldVersion,
		NewVersion: newVersion,
		Timestamp:  time.Now().UTC().Format(time.RFC3339),
	}

	body, err := json.Marshal(payload)
	if err != nil {
		return fmt.Errorf("notifier: failed to marshal payload: %w", err)
	}

	resp, err := n.client.Post(n.webhookURL, "application/json", bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("notifier: failed to send webhook: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		return fmt.Errorf("notifier: webhook returned non-2xx status: %d", resp.StatusCode)
	}

	return nil
}
