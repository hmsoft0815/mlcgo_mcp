package mcp

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mlechner/mlc_toolretrieval/openstreetmap_mcp/internal/osm"
)

func setupTestServer(t *testing.T, response string) (*Server, *httptest.Server) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, response)
	}))

	client := osm.NewClient("test-agent")
	client.NominatimBaseURL = ts.URL
	client.OverpassBaseURL = ts.URL
	client.OSRMBaseURL = ts.URL
	client.SetRateLimit(0)

	s := NewServer("test", "1.0.0", client)
	s.RegisterTools() // Ensure tools are registered
	return s, ts
}

func TestHandleGeocode(t *testing.T) {
	s, ts := setupTestServer(t, `[]`)
	defer ts.Close()

	args := GeocodeArgs{Address: "Berlin"}
	argsJSON, _ := json.Marshal(args)
	var argsMap map[string]any
	json.Unmarshal(argsJSON, &argsMap)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "geocode_address",
			Arguments: argsMap,
		},
	}

	res, err := s.handleGeocode(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGeocode failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("handleGeocode returned error: %v", res.Content)
	}
}

func TestHandleReverseGeocode(t *testing.T) {
	s, ts := setupTestServer(t, `{"place_id": 1}`)
	defer ts.Close()

	args := ReverseGeocodeArgs{Latitude: 52.5, Longitude: 13.4}
	argsJSON, _ := json.Marshal(args)
	var argsMap map[string]any
	json.Unmarshal(argsJSON, &argsMap)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "reverse_geocode",
			Arguments: argsMap,
		},
	}

	res, err := s.handleReverseGeocode(context.Background(), req)
	if err != nil {
		t.Fatalf("handleReverseGeocode failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("handleReverseGeocode returned error: %v", res.Content)
	}
}

func TestHandleFindNearby(t *testing.T) {
	s, ts := setupTestServer(t, `{"elements":[]}`)
	defer ts.Close()

	args := FindNearbyArgs{Latitude: 52.5, Longitude: 13.4, Radius: 1000}
	argsJSON, _ := json.Marshal(args)
	var argsMap map[string]any
	json.Unmarshal(argsJSON, &argsMap)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "find_nearby_places",
			Arguments: argsMap,
		},
	}

	res, err := s.handleFindNearby(context.Background(), req)
	if err != nil {
		t.Fatalf("handleFindNearby failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("handleFindNearby returned error: %v", res.Content)
	}
}

func TestHandleGetRoute(t *testing.T) {
	s, ts := setupTestServer(t, `{"routes":[]}`)
	defer ts.Close()

	args := GetRouteArgs{FromLat: 52.5, FromLon: 13.3, ToLat: 52.6, ToLon: 13.4}
	argsJSON, _ := json.Marshal(args)
	var argsMap map[string]any
	json.Unmarshal(argsJSON, &argsMap)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "get_route",
			Arguments: argsMap,
		},
	}

	res, err := s.handleGetRoute(context.Background(), req)
	if err != nil {
		t.Fatalf("handleGetRoute failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("handleGetRoute returned error: %v", res.Content)
	}
}

func TestHandleSearchCategory(t *testing.T) {
	s, ts := setupTestServer(t, `{"elements":[]}`)
	defer ts.Close()

	args := SearchCategoryArgs{Category: "amenity=restaurant", MinLat: 52.5, MinLon: 13.3, MaxLat: 52.6, MaxLon: 13.4}
	argsJSON, _ := json.Marshal(args)
	var argsMap map[string]any
	json.Unmarshal(argsJSON, &argsMap)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "search_category",
			Arguments: argsMap,
		},
	}

	res, err := s.handleSearchCategory(context.Background(), req)
	if err != nil {
		t.Fatalf("handleSearchCategory failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("handleSearchCategory returned error: %v", res.Content)
	}
}

func TestHandleNearbyPOI(t *testing.T) {
	s, ts := setupTestServer(t, `{"elements":[]}`)
	defer ts.Close()

	args := FindPOIArgs{Latitude: 52.5, Longitude: 13.4, Radius: 1000}
	argsJSON, _ := json.Marshal(args)
	var argsMap map[string]any
	json.Unmarshal(argsJSON, &argsMap)

	req := mcp.CallToolRequest{
		Params: mcp.CallToolParams{
			Name:      "find_schools",
			Arguments: argsMap,
		},
	}

	res, err := s.handleNearbyPOI(context.Background(), req, "amenity=school")
	if err != nil {
		t.Fatalf("handleNearbyPOI failed: %v", err)
	}
	if res.IsError {
		t.Fatalf("handleNearbyPOI returned error: %v", res.Content)
	}
}
