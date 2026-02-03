package polymarket

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

const (
	// GammaAPIBaseURL is the base URL for the Polymarket Gamma API
	GammaAPIBaseURL = "https://gamma-api.polymarket.com"
)

// Client represents a client for interacting with the Polymarket API
type Client struct {
	BaseURL    string
	HTTPClient *http.Client
}

// NewClient creates a new Polymarket API client with default settings
func NewClient() *Client {
	return &Client{
		BaseURL:    GammaAPIBaseURL,
		HTTPClient: &http.Client{},
	}
}

// Tag represents a tag/category for an event
type Tag struct {
	ID    string `json:"id"`
	Label string `json:"label"`
	Slug  string `json:"slug"`
}

// Market represents a market within an event
type Market struct {
	ID            string `json:"id"`
	Question      string `json:"question"`
	ClobTokenIds  string `json:"clobTokenIds"`  // JSON string array
	Outcomes      string `json:"outcomes"`      // JSON string array
	OutcomePrices string `json:"outcomePrices"` // JSON string array
}

// Event represents a Polymarket event
type Event struct {
	ID      string   `json:"id"`
	Slug    string   `json:"slug"`
	Title   string   `json:"title"`
	Active  bool     `json:"active"`
	Closed  bool     `json:"closed"`
	Tags    []Tag    `json:"tags"`
	Markets []Market `json:"markets"`
}

// FetchActiveEventsOptions contains options for fetching active events
type FetchActiveEventsOptions struct {
	Limit int
}

// FetchActiveEvents fetches active events from the Polymarket API
func (c *Client) FetchActiveEvents(options *FetchActiveEventsOptions) ([]Event, error) {
	// Build URL with query parameters
	u, err := url.Parse(c.BaseURL + "/events")
	if err != nil {
		return nil, fmt.Errorf("failed to parse base URL: %w", err)
	}

	q := u.Query()
	q.Set("active", "true")
	q.Set("closed", "false")

	if options != nil && options.Limit > 0 {
		q.Set("limit", fmt.Sprintf("%d", options.Limit))
	}

	u.RawQuery = q.Encode()

	// Make HTTP request
	resp, err := c.HTTPClient.Get(u.String())
	if err != nil {
		return nil, fmt.Errorf("failed to make request: %w", err)
	}
	defer resp.Body.Close()

	// Check status code
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("API returned status %d: %s", resp.StatusCode, string(body))
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("failed to read response body: %w", err)
	}

	// Unmarshal JSON
	var events []Event
	if err := json.Unmarshal(body, &events); err != nil {
		return nil, fmt.Errorf("failed to unmarshal JSON: %w", err)
	}

	return events, nil
}
