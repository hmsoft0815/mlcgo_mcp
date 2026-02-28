package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
)

// OverpassResponse represents the JSON response from the Overpass API interpreter.
// It contains metadata and a slice of map elements (nodes, ways, relations).
type OverpassResponse struct {
	Version   float64           `json:"version"`
	Generator string            `json:"generator"`
	Osm3s     map[string]any    `json:"osm3s"`
	Elements  []OverpassElement `json:"elements"`
}

// OverpassElement represents a single OpenStreetMap element (node, way, or relation).
// Depending on the element type, it may contain coordinates, tags, or a list of nodes.
type OverpassElement struct {
	Type     string            `json:"type"`
	ID       int64             `json:"id"`
	Lat      float64           `json:"lat,omitempty"`
	Lon      float64           `json:"lon,omitempty"`
	Center   *OverpassLocation `json:"center,omitempty"`
	Tags     map[string]string `json:"tags"`
	Nodes    []int64           `json:"nodes,omitempty"`
	Geometry []OverpassPoint   `json:"geometry,omitempty"`
}

// OverpassLocation represents a latitude and longitude coordinate.
type OverpassLocation struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// OverpassPoint represents a single point in a geometry or line.
type OverpassPoint struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// OverpassQuery executes a query against the Overpass API interpreter using a POST request.
// It is commonly used for searching POIs, bounding box searches, and data extraction.
// Note: This method honors the configured rate limit.
func (c *Client) OverpassQuery(ctx context.Context, query string) (*OverpassResponse, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}
	u := c.OverpassBaseURL
	data := url.Values{}
	data.Set("data", query)

	req, err := http.NewRequestWithContext(ctx, "POST", u, strings.NewReader(data.Encode()))
	if err != nil {
		return nil, fmt.Errorf("failed to create request: %w", err)
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("User-Agent", c.userAgent)

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf("unexpected status code: %d, body: %s", resp.StatusCode, string(body))
	}

	var result OverpassResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
