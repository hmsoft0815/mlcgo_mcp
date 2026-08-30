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

// version is stamped from the VERSION file at build time via
// `-ldflags -X main.version`. "dev" is what an unstamped build reports, which
// is honest — until 2026-08-30 the number was written out twice as a literal
// here and could not move without an edit.
var version = "dev"

const serverName = "openstreetmap-mcp"

func main() {
	transport := flag.String("transport", "stdio", "Transport to use (stdio or sse)")
	sseAddr := flag.String("sse-addr", ":8080", "Address to listen on for SSE")
	rateLimitSecs := flag.Int("osm-rate-limit", 5, "Minimum seconds between OpenStreetMap API requests")
	showVersion := flag.Bool("version", false, "Show version and exit")
	flag.Parse()

	if *showVersion {
		fmt.Printf("%s version %s\n", serverName, version)
		return
	}

	// The User-Agent identifies us to OpenStreetMap's servers; it carried the
	// same hardcoded number and would have kept claiming 1.0.0 forever.
	osmClient := osm.NewClient(fmt.Sprintf(
		"%s/%s (https://github.com/mlechner/mlc_toolretrieval/openstreetmap_mcp)", serverName, version))
	osmClient.SetRateLimit(time.Duration(*rateLimitSecs) * time.Second)
	mcpServer := mcp.NewServer(serverName, version, osmClient)
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
