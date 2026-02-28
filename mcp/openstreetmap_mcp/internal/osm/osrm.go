package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

// OSRMResponse is the top-level response structure from the OSRM route API.
// It contains the status code, a list of calculated routes, and waypoints.
type OSRMResponse struct {
	Code      string      `json:"code"`
	Routes    []OSRMRoute `json:"routes"`
	Waypoints []any       `json:"waypoints"`
}

// OSRMRoute represents a calculated route between points, including distance and duration.
type OSRMRoute struct {
	Geometry   any             `json:"geometry"`
	Legs       []OSRMRouteLeg  `json:"legs"`
	WeightName string          `json:"weight_name"`
	Weight     float64         `json:"weight"`
	Duration   float64         `json:"duration"`
	Distance   float64         `json:"distance"`
}

// OSRMRouteLeg represents a segment of a route between two waypoints, containing multiple steps.
type OSRMRouteLeg struct {
	Steps    []any   `json:"steps"`
	Weight   float64 `json:"weight"`
	Duration float64 `json:"duration"`
	Distance float64 `json:"distance"`
	Summary  string  `json:"summary"`
}

// GetRoute calculates the best path between two coordinates using the OSRM routing engine.
// Supported profiles include "car", "bicycle", and "foot".
// Note: This method honors the configured rate limit.
func (c *Client) GetRoute(ctx context.Context, fromLat, fromLon, toLat, toLon float64, mode string) (*OSRMResponse, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}
	// Mode defaults to car. OSRM public demo uses different endpoints for these.
	profile := "car"
	switch mode {
	case "bike", "bicycle":
		profile = "bicycle"
	case "foot", "walking":
		profile = "foot"
	}

	u := fmt.Sprintf("%s/route/v1/%s/%f,%f;%f,%f?overview=full&geometries=geojson", c.OSRMBaseURL, profile, fromLon, fromLat, toLat, toLon)

	req, err := http.NewRequestWithContext(ctx, "GET", u, nil)
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var result OSRMResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
