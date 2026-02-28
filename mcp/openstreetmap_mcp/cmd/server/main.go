package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/mlechner/mlc_toolretrieval/openstreetmap_mcp/internal/mcp"
	"github.com/mlechner/mlc_toolretrieval/openstreetmap_mcp/internal/osm"
)

func main() {
	transport := flag.String("transport", "stdio", "Transport to use (stdio or sse)")
	sseAddr := flag.String("sse-addr", ":8080", "Address to listen on for SSE")
	rateLimitSecs := flag.Int("osm-rate-limit", 5, "Minimum seconds between OpenStreetMap API requests")
	flag.Parse()

	osmClient := osm.NewClient("openstreetmap-mcp/1.0.0 (https://github.com/mlechner/mlc_toolretrieval/openstreetmap_mcp)")
	osmClient.SetRateLimit(time.Duration(*rateLimitSecs) * time.Second)
	mcpServer := mcp.NewServer("openstreetmap-mcp", "1.0.0", osmClient)
	mcpServer.RegisterTools()

	switch *transport {
	case "stdio":
		log.Println("Starting OpenStreetMap MCP server on stdio...")
		if err := mcpServer.ServeStdio(); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	case "sse":
		log.Printf("Starting OpenStreetMap MCP server on SSE (%s)...\n", *sseAddr)
		if err := mcpServer.ServeSSE(*sseAddr); err != nil {
			fmt.Fprintf(os.Stderr, "Error: %v\n", err)
			os.Exit(1)
		}
	default:
		fmt.Fprintf(os.Stderr, "Unknown transport: %s\n", *transport)
		os.Exit(1)
	}
}
