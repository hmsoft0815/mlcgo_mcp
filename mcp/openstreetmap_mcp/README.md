# OpenStreetMap MCP Server

Copyright (c) 2026 Michael Lechner. All rights reserved. (Go)

A Model Context Protocol (MCP) server for OpenStreetMap data, providing tools for geocoding, routing, and spatial analysis.

## Features
- **Golang Implementation:** Clean and efficient Go code.
- **Transports:** Supports both `stdio` and `SSE` (Server-Sent Events).
- **Structured Replies:** Uses `mcp-go` to provide structured JSON results.
- **Clean Schema:** Automatically generated JSON schemas from Go structs.
- **API Integrations:**
  - **[Nominatim](https://nominatim.org/):** Geocoding and Reverse Geocoding.
  - **[OSRM](http://project-osrm.org/):** Routing (car, bicycle, foot).
  - **[Overpass API](https://wiki.openstreetmap.org/wiki/Overpass_API):** Spatial queries (nearby POIs, schools, EV stations, parking).

## OpenStreetMap Fair Usage
This server interacts with public [OpenStreetMap](https://www.openstreetmap.org/) services. It is critical to adhere to their usage policies:
- **Rate Limiting:** Do not make more than one request per second. This server enforces a **default 5-second delay** between requests (configurable).
- **User-Agent:** Always provide a descriptive User-Agent (pre-configured in this server).
- **Usage Policies:** Please review the official [Nominatim Usage Policy](https://operations.osmfoundation.org/policies/nominatim/) and [Overpass API Usage Policy](https://operations.osmfoundation.org/policies/overpass/).

## Tools
- `geocode_address`: Converts an address to coordinates using Nominatim.
- `reverse_geocode`: Converts coordinates to an address using Nominatim.
- `find_nearby_places`: Finds POIs near a location using Overpass API.
- `get_route`: Calculates routing between two points using OSRM.
- `search_category`: Searches for features in a bounding box using Overpass API.
- `find_schools`: Locates schools near coordinates using Overpass API.
- `find_ev_charging_stations`: Locates EV chargers using Overpass API.
- `find_parking`: Locates parking facilities using Overpass API.

## Build
```bash
make build
```

## Usage

### Stdio Transport (Default)
```bash
./osm-mcp --transport stdio --osm-rate-limit 5
```

### SSE Transport
```bash
./osm-mcp --transport sse --sse-addr :8080 --osm-rate-limit 5
```

### Configuration Flags
- `--transport`: `stdio` or `sse` (default: `stdio`)
- `--sse-addr`: Address for SSE server (default: `:8080`)
- `--osm-rate-limit`: Minimum seconds to wait between OSM API calls (default: `5`)

## License
MIT
