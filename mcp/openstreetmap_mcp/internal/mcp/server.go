// Package mcp implements the Model Context Protocol (MCP) server logic,
// exposing OpenStreetMap functionalities as tools.
package mcp

import (
	"context"
	"fmt"
	"strings"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
	"github.com/mlechner/mlc_toolretrieval/openstreetmap_mcp/internal/osm"
)

// Server wraps the MCP server and provides handlers for various OSM tools.
type Server struct {
	mcpServer *server.MCPServer
	osmClient *osm.Client
}

// NewServer creates a new Server instance with the given name, version, and OSM client.
func NewServer(name, version string, osmClient *osm.Client) *Server {
	s := server.NewMCPServer(name, version)
	return &Server{
		mcpServer: s,
		osmClient: osmClient,
	}
}

// GeocodeArgs defines the input arguments for the geocode_address tool.
type GeocodeArgs struct {
	Address string `json:"address" jsonschema:"description=The address or place name to geocode."`
}

// ReverseGeocodeArgs defines the input arguments for the reverse_geocode tool.
type ReverseGeocodeArgs struct {
	Latitude  float64 `json:"latitude" jsonschema:"description=The latitude coordinate."`
	Longitude float64 `json:"longitude" jsonschema:"description=The longitude coordinate."`
}

// FindNearbyArgs defines the input arguments for the find_nearby_places tool.
type FindNearbyArgs struct {
	Latitude   float64  `json:"latitude" jsonschema:"description=The latitude coordinate."`
	Longitude  float64  `json:"longitude" jsonschema:"description=The longitude coordinate."`
	Radius     float64  `json:"radius" jsonschema:"description=Search radius in meters."`
	Categories []string `json:"categories,omitempty" jsonschema:"description=Categories to search for (e.g. restaurant, cafe, school)."`
	Limit      int      `json:"limit,omitempty" jsonschema:"description=Maximum number of results to return.,default=10"`
}

// GetRouteArgs defines the input arguments for the get_route tool.
type GetRouteArgs struct {
	FromLat float64 `json:"from_lat"`
	FromLon float64 `json:"from_lon"`
	ToLat   float64 `json:"to_lat"`
	ToLon   float64 `json:"to_lon"`
	Mode    string  `json:"mode,omitempty" jsonschema:"description=The transportation mode.,enum=car,enum=bicycle,enum=foot,default=car"`
}

// SearchCategoryArgs defines the input arguments for the search_category tool.
type SearchCategoryArgs struct {
	Category string  `json:"category" jsonschema:"description=The category to search for (e.g., amenity=restaurant)."`
	MinLat   float64 `json:"min_lat"`
	MinLon   float64 `json:"min_lon"`
	MaxLat   float64 `json:"max_lat"`
	MaxLon   float64 `json:"max_lon"`
}

// FindPOIArgs defines the input arguments for specialized POI search tools like find_schools.
type FindPOIArgs struct {
	Latitude  float64 `json:"latitude"`
	Longitude float64 `json:"longitude"`
	Radius    float64 `json:"radius" jsonschema:"description=Search radius in meters."`
}

// RegisterTools registers all available OSM tools with the MCP server.
func (s *Server) RegisterTools() {
	// 1. Geocode Address
	s.mcpServer.AddTool(mcp.NewTool("geocode_address",
		mcp.WithDescription("Converts an address or place name to coordinates."),
		mcp.WithInputSchema[GeocodeArgs](),
	), s.handleGeocode)

	// 2. Reverse Geocode
	s.mcpServer.AddTool(mcp.NewTool("reverse_geocode",
		mcp.WithDescription("Converts coordinates to a human-readable address."),
		mcp.WithInputSchema[ReverseGeocodeArgs](),
	), s.handleReverseGeocode)

	// 3. Find Nearby Places
	s.mcpServer.AddTool(mcp.NewTool("find_nearby_places",
		mcp.WithDescription("Finds points of interest near a location."),
		mcp.WithInputSchema[FindNearbyArgs](),
	), s.handleFindNearby)

	// 4. Get Route
	s.mcpServer.AddTool(mcp.NewTool("get_route",
		mcp.WithDescription("Calculates a route between two points."),
		mcp.WithInputSchema[GetRouteArgs](),
	), s.handleGetRoute)

	// 5. Search Category in Bounding Box
	s.mcpServer.AddTool(mcp.NewTool("search_category",
		mcp.WithDescription("Searches for features of a specific category within a bounding box."),
		mcp.WithInputSchema[SearchCategoryArgs](),
	), s.handleSearchCategory)

	// 6. Find Schools
	s.mcpServer.AddTool(mcp.NewTool("find_schools",
		mcp.WithDescription("Finds schools near a location."),
		mcp.WithInputSchema[FindPOIArgs](),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleNearbyPOI(ctx, req, "amenity=school")
	})

	// 7. Find EV Charging Stations
	s.mcpServer.AddTool(mcp.NewTool("find_ev_charging_stations",
		mcp.WithDescription("Finds EV charging stations near a location."),
		mcp.WithInputSchema[FindPOIArgs](),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleNearbyPOI(ctx, req, "amenity=charging_station")
	})

	// 8. Find Parking
	s.mcpServer.AddTool(mcp.NewTool("find_parking",
		mcp.WithDescription("Finds parking facilities near a location."),
		mcp.WithInputSchema[FindPOIArgs](),
	), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		return s.handleNearbyPOI(ctx, req, "amenity=parking")
	})
}

// handleGeocode handles the geocode_address tool request.
func (s *Server) handleGeocode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args GeocodeArgs
	if err := req.BindArguments(&args); err != nil {
		return nil, err
	}

	results, err := s.osmClient.Geocode(ctx, args.Address)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(results)
}

// handleReverseGeocode handles the reverse_geocode tool request.
func (s *Server) handleReverseGeocode(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args ReverseGeocodeArgs
	if err := req.BindArguments(&args); err != nil {
		return nil, err
	}

	result, err := s.osmClient.ReverseGeocode(ctx, args.Latitude, args.Longitude)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(result)
}

// handleFindNearby handles the find_nearby_places tool request.
func (s *Server) handleFindNearby(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args FindNearbyArgs
	if err := req.BindArguments(&args); err != nil {
		return nil, err
	}
	if args.Limit == 0 {
		args.Limit = 10
	}

	var amenityFilter string
	if len(args.Categories) > 0 {
		amenityFilter = fmt.Sprintf(`["amenity"~"%s"]`, strings.Join(args.Categories, "|"))
	} else {
		amenityFilter = `["amenity"]`
	}

	query := fmt.Sprintf(`[out:json];node%s(around:%f,%f,%f);out %d;`,
		amenityFilter, args.Radius, args.Latitude, args.Longitude, args.Limit)

	results, err := s.osmClient.OverpassQuery(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(results.Elements)
}

// handleGetRoute handles the get_route tool request.
func (s *Server) handleGetRoute(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args GetRouteArgs
	if err := req.BindArguments(&args); err != nil {
		return nil, err
	}
	if args.Mode == "" {
		args.Mode = "car"
	}

	result, err := s.osmClient.GetRoute(ctx, args.FromLat, args.FromLon, args.ToLat, args.ToLon, args.Mode)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(result)
}

// handleSearchCategory handles the search_category tool request.
func (s *Server) handleSearchCategory(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args SearchCategoryArgs
	if err := req.BindArguments(&args); err != nil {
		return nil, err
	}

	// query format: node["amenity"="restaurant"](52.5,13.3,52.6,13.4);out;
	query := fmt.Sprintf(`[out:json];node[%s](%f,%f,%f,%f);out;`,
		args.Category, args.MinLat, args.MinLon, args.MaxLat, args.MaxLon)

	results, err := s.osmClient.OverpassQuery(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(results.Elements)
}

// handleNearbyPOI handles specialized POI search tools (e.g., find_schools).
func (s *Server) handleNearbyPOI(ctx context.Context, req mcp.CallToolRequest, category string) (*mcp.CallToolResult, error) {
	var args FindPOIArgs
	if err := req.BindArguments(&args); err != nil {
		return nil, err
	}

	query := fmt.Sprintf(`[out:json];node[%s](around:%f,%f,%f);out;`,
		category, args.Radius, args.Latitude, args.Longitude)

	results, err := s.osmClient.OverpassQuery(ctx, query)
	if err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	return mcp.NewToolResultJSON(results.Elements)
}

// ServeStdio starts the MCP server using the stdio transport.
func (s *Server) ServeStdio() error {
	return server.ServeStdio(s.mcpServer)
}

// ServeSSE starts the MCP server using the SSE transport on the specified address.
func (s *Server) ServeSSE(addr string) error {
	sseServer := server.NewSSEServer(s.mcpServer, server.WithStaticBasePath("/"))
	return sseServer.Start(addr)
}
