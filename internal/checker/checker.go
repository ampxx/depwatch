package checker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

// ModuleVersion holds the latest version info for a Go module.
type ModuleVersion struct {
	Version string    `json:"Version"`
	Time    time.Time `json:"Time"`
}

// Client fetches module version information from the Go module proxy.
type Client struct {
	HTTPClient *http.Client
	ProxyURL   string
}

// NewClient creates a new checker Client with sensible defaults.
func NewClient() *Client {
	return &Client{
		HTTPClient: &http.Client{Timeout: 10 * time.Second},
		ProxyURL:   "https://proxy.golang.org",
	}
}

// LatestVersion queries the Go module proxy for the latest version of the given module.
func (c *Client) LatestVersion(modulePath string) (*ModuleVersion, error) {
	url := fmt.Sprintf("%s/%s/@latest", c.ProxyURL, modulePath)

	resp, err := c.HTTPClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("fetching latest version for %s: %w", modulePath, err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("proxy returned status %d for module %s", resp.StatusCode, modulePath)
	}

	var mv ModuleVersion
	if err := json.NewDecoder(resp.Body).Decode(&mv); err != nil {
		return nil, fmt.Errorf("decoding response for %s: %w", modulePath, err)
	}

	return &mv, nil
}
