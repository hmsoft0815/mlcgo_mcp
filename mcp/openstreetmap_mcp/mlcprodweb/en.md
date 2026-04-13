# OpenStreetMap MCP Server

Empower your LLM with real-world geography. The OpenStreetMap MCP server provides a seamless bridge between modern AI models and the world's most comprehensive open map data.

## Why OpenStreetMap MCP?

Geography is the context for everything. Whether you are building an assistant that needs to find local schools, a logistics tool for routing, or a site analysis engine, this server provides the precise, structured data you need directly from OpenStreetMap.

## Key Capabilities

- **Find anything, anywhere**: Leverage the **Overpass API** to search for thousands of POI categories—from EV charging stations to local pharmacies—with high spatial precision.
- **Convert Adresses to Insights**: Use **Nominatim** to transform human-readable addresses into machine-readable coordinates (and vice versa) for accurate mapping.
- **Smart Navigation**: Integrate **OSRM-based routing** to calculate travel times and distances for cars, bicycles, or pedestrians.
- **Fair Usage by Design**: Built-in rate limiting ensures your integration stays respectful of OpenStreetMap's community-driven infrastructure.

## Getting Started

The server supports both **Stdio** for local desktop integration and **SSE** for remote services, facilitating deployment across a variety of environments.

## Quick Setup

### Claude Desktop
Add the following to your `claude_desktop_config.json`:

```json
{
  "mcpServers": {
    "openstreetmap": {
      "command": "openstreetmap_mcp"
    }
  }
}
```

### Gemini-CLI
Add to your `~/.gemini/settings.json`:

```json
{
  "mcpServers": {
    "openstreetmap": {
      "command": "openstreetmap_mcp"
    }
  }
}
```

### MCP-Tester
Add a new profile:

```bash
mcp-tester profile add osm -c "openstreetmap_mcp"
```
