// Package osm provides a client for interacting with various OpenStreetMap APIs.
// This includes Nominatim for geocoding, Overpass for querying map features,
// and OSRM for routing services.
package osm

import (
	"context"
	"net/http"
	"sync"
	"time"
)

// Client handles communication with the OpenStreetMap APIs.
// It encapsulates an HTTP client and provides methods for interacting
// with Geocoding, Overpass, and Routing services.
type Client struct {
	// httpClient is the underlying HTTP client used for making requests.
	httpClient *http.Client
	// userAgent is the string used in the User-Agent header for all requests,
	// which is required by OpenStreetMap's usage policies.
	userAgent string
	// rateLimit is the minimum duration between consecutive requests to OSM services.
	// Defaults to 5 seconds to comply with fair usage policies.
	rateLimit time.Duration
	// lastRequestMu protects lastRequest time.
	lastRequestMu sync.Mutex
	// lastRequest tracks the time of the last issued request.
	lastRequest time.Time

	// API Base URLs (configurable for testing)
	NominatimBaseURL string
	OverpassBaseURL  string
	OSRMBaseURL      string
}

// NewClient creates a new OpenStreetMap client with the given User-Agent string.
// A descriptive User-Agent is recommended according to OSM guidelines.
// It defaults to a 5-second rate limit between requests.
func NewClient(userAgent string) *Client {
	return &Client{
		httpClient:       &http.Client{},
		userAgent:        userAgent,
		rateLimit:        5 * time.Second,
		NominatimBaseURL: "https://nominatim.openstreetmap.org",
		OverpassBaseURL:  "https://overpass-api.de/api/interpreter",
		OSRMBaseURL:      "https://router.project-osrm.org",
	}
}

// SetRateLimit updates the minimum duration between consecutive requests.
func (c *Client) SetRateLimit(d time.Duration) {
	c.rateLimit = d
}

// waitRateLimit blocks until the rate limit duration has passed since the last request.
func (c *Client) waitRateLimit(ctx context.Context) error {
	c.lastRequestMu.Lock()
	defer c.lastRequestMu.Unlock()

	now := time.Now()
	elapsed := now.Sub(c.lastRequest)
	if elapsed < c.rateLimit {
		waitTime := c.rateLimit - elapsed
		select {
		case <-time.After(waitTime):
		case <-ctx.Done():
			return ctx.Err()
		}
	}
	c.lastRequest = time.Now()
	return nil
}
