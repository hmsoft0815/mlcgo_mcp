package osm

import (
	"context"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

func TestClient_Geocode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/search" {
			t.Errorf("expected path /search, got %s", r.URL.Path)
		}
		query := r.URL.Query().Get("q")
		if query != "Berlin" {
			t.Errorf("expected query Berlin, got %s", query)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `[{"place_id":1001, "lat":"52.517", "lon":"13.388", "display_name":"Berlin, Germany"}]`)
	}))
	defer ts.Close()

	client := NewClient("test-agent")
	client.NominatimBaseURL = ts.URL
	client.SetRateLimit(0) // Disable rate limit for tests

	results, err := client.Geocode(context.Background(), "Berlin")
	if err != nil {
		t.Fatalf("Geocode failed: %v", err)
	}

	if len(results) != 1 || results[0].DisplayName != "Berlin, Germany" {
		t.Errorf("unexpected results: %+v", results)
	}
}

func TestClient_ReverseGeocode(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/reverse" {
			t.Errorf("expected path /reverse, got %s", r.URL.Path)
		}

		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"place_id":1001, "lat":"52.517", "lon":"13.388", "display_name":"Brandenburg Gate, Berlin"}`)
	}))
	defer ts.Close()

	client := NewClient("test-agent")
	client.NominatimBaseURL = ts.URL
	client.SetRateLimit(0)

	result, err := client.ReverseGeocode(context.Background(), 52.517, 13.388)
	if err != nil {
		t.Fatalf("ReverseGeocode failed: %v", err)
	}

	if result.DisplayName != "Brandenburg Gate, Berlin" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestClient_OverpassQuery(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"elements":[{"type":"node", "id":123, "lat":52.5, "lon":13.4, "tags":{"amenity":"school", "name":"Test School"}}]}`)
	}))
	defer ts.Close()

	client := NewClient("test-agent")
	client.OverpassBaseURL = ts.URL
	client.SetRateLimit(0)

	result, err := client.OverpassQuery(context.Background(), "dummy query")
	if err != nil {
		t.Fatalf("OverpassQuery failed: %v", err)
	}

	if len(result.Elements) != 1 || result.Elements[0].Tags["name"] != "Test School" {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestClient_GetRoute(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprint(w, `{"code":"Ok", "routes":[{"distance":1500, "duration":300}]}`)
	}))
	defer ts.Close()

	client := NewClient("test-agent")
	client.OSRMBaseURL = ts.URL
	client.SetRateLimit(0)

	result, err := client.GetRoute(context.Background(), 52.5, 13.3, 52.6, 13.4, "car")
	if err != nil {
		t.Fatalf("GetRoute failed: %v", err)
	}

	if len(result.Routes) != 1 || result.Routes[0].Distance != 1500 {
		t.Errorf("unexpected result: %+v", result)
	}
}

func TestClient_RateLimit(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprint(w, `[]`)
	}))
	defer ts.Close()

	client := NewClient("test-agent")
	client.NominatimBaseURL = ts.URL
	limit := 100 * time.Millisecond
	client.SetRateLimit(limit)

	start := time.Now()
	// First request
	client.Geocode(context.Background(), "test1")
	// Second request (should be delayed)
	client.Geocode(context.Background(), "test2")
	elapsed := time.Since(start)

	if elapsed < limit {
		t.Errorf("rate limit not enforced: took %v, expected at least %v", elapsed, limit)
	}
}
