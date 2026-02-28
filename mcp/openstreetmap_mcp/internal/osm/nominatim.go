package osm

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
)

// GeocodeResponse represents the structure returned by the Nominatim search API.
// It includes the place ID, OSM identity, coordinates, and display name.
type GeocodeResponse struct {
	PlaceID     int64    `json:"place_id"`
	Licence     string   `json:"licence"`
	OSMType     string   `json:"osm_type"`
	OSMID       int64    `json:"osm_id"`
	BoundingBox []string `json:"boundingbox"`
	Lat         string   `json:"lat"`
	Lon         string   `json:"lon"`
	DisplayName string   `json:"display_name"`
	Class       string   `json:"class"`
	Type        string   `json:"type"`
	Importance  float64  `json:"importance"`
}

// Geocode converts an address or place name into geographic coordinates using the Nominatim search API.
// It returns a slice of matching GeocodeResponse results.
// Note: This method honors the configured rate limit.
func (c *Client) Geocode(ctx context.Context, address string) ([]GeocodeResponse, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}
	u, _ := url.Parse(c.NominatimBaseURL + "/search")
	q := u.Query()
	q.Set("q", address)
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
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

	var results []GeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return results, nil
}

// ReverseGeocodeResponse represents the structure returned by the Nominatim reverse geocoding API.
// It contains detailed address information for a given coordinate.
type ReverseGeocodeResponse struct {
	PlaceID     int64             `json:"place_id"`
	Licence     string            `json:"licence"`
	OSMType     string            `json:"osm_type"`
	OSMID       int64             `json:"osm_id"`
	Lat         string            `json:"lat"`
	Lon         string            `json:"lon"`
	DisplayName string            `json:"display_name"`
	Address     map[string]string `json:"address"`
	BoundingBox []string          `json:"boundingbox"`
}

// ReverseGeocode converts geographic coordinates into a human-readable address using the Nominatim reverse geocoding API.
// Note: This method honors the configured rate limit.
func (c *Client) ReverseGeocode(ctx context.Context, lat, lon float64) (*ReverseGeocodeResponse, error) {
	if err := c.waitRateLimit(ctx); err != nil {
		return nil, err
	}
	u, _ := url.Parse(c.NominatimBaseURL + "/reverse")
	q := u.Query()
	q.Set("lat", fmt.Sprintf("%f", lat))
	q.Set("lon", fmt.Sprintf("%f", lon))
	q.Set("format", "json")
	u.RawQuery = q.Encode()

	req, err := http.NewRequestWithContext(ctx, "GET", u.String(), nil)
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

	var result ReverseGeocodeResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	return &result, nil
}
