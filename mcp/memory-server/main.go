package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/hmsoft0815/memory-server/internal/handlers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type MemorizeArgs struct {
	Entity      string `json:"entity" jsonschema:"description=The name of the thing (e.g. 'Oly' or 'Project')"`
	Category    string `json:"category,omitempty" jsonschema:"description=Category (e.g. 'Person', 'Setting')"`
	Observation string `json:"observation" jsonschema:"description=The actual fact to remember"`
}

type SearchNodesArgs struct {
	Query string `json:"query" jsonschema:"description=The search term"`
}

type ReadGraphArgs struct{}

func main() {
	// Ensure logging is on stderr for MCP protocol safety
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	v := flag.Bool("version", false, "Print version and exit")
	transportType := flag.String("transport", "stdio", "Transport type (stdio or sse)")
	addr := flag.String("addr", ":3000", "SSE address (if transport=sse)")
	flag.Parse()

	if *v {
		fmt.Println("memory-server 1.2.0")
		return
	}

	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".local", "share", "mcp-proxy")
	os.MkdirAll(dbDir, 0755)
	dbPath := filepath.Join(dbDir, "memory.db")

	handler, err := handlers.NewMemoryHandler(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite: %v", err)
	}

	// Initialize the MCP server
	s := server.NewMCPServer("memory-server", "1.2.0")

	// Register tools
	registerTools(s, handler)

	// Run the server
	runServer(s, *transportType, *addr)
}

func registerTools(s *server.MCPServer, handler *handlers.MemoryHandler) {
	s.AddTool(mcp.NewTool("memory__memorize__mlc",
		mcp.WithDescription("Store a new fact or observation about an entity"),
		mcp.WithInputSchema[MemorizeArgs](),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args MemorizeArgs
		if err := request.BindArguments(&args); err != nil {
			return nil, err
		}

		hArgs := map[string]interface{}{
			"entities": []interface{}{
				map[string]interface{}{
					"name":         args.Entity,
					"entityType":   args.Category,
					"observations": []interface{}{args.Observation},
				},
			},
		}
		res, err := handler.CreateEntities(hArgs)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(res.(string)), nil
	})

	s.AddTool(mcp.NewTool("memory__search_nodes__mlc",
		mcp.WithDescription("Search for stored facts in the knowledge graph"),
		mcp.WithInputSchema[SearchNodesArgs](),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		var args SearchNodesArgs
		if err := request.BindArguments(&args); err != nil {
			return nil, err
		}

		hArgs := map[string]interface{}{
			"query": args.Query,
		}
		res, err := handler.SearchNodes(hArgs)
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(res.(string)), nil
	})

	s.AddTool(mcp.NewTool("memory__read_graph__mlc",
		mcp.WithDescription("Read the entire knowledge graph"),
	), func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		res, err := handler.ReadGraph(map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		return mcp.NewToolResultText(res.(string)), nil
	})

	s.AddPrompt(mcp.NewPrompt("memory__knowledge_extraction__mlc",
		mcp.WithPromptDescription("Deep extraction of related entities and observations for a specific topic."),
		mcp.WithArgument("topic", mcp.ArgumentDescription("The topic to research in the knowledge graph"), mcp.RequiredArgument()),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		topic := request.Params.Arguments["topic"]
		return &mcp.GetPromptResult{
			Description: "Knowledge extraction for " + topic,
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Text: "You are a knowledge graph expert. Search the memory for all entities and observations related to '" + topic + "'. Construct a narrative summary that connects these pieces of information, highlighting key relationships and historical context stored in the graph.",
					},
				},
			},
		}, nil
	})
}

func runServer(s *server.MCPServer, transportType, addr string) {
	switch transportType {
	case "stdio":
		log.Println("Memory MCP Server starting on stdio...")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	case "sse":
		log.Printf("Memory MCP Server starting on SSE at %s...", addr)
		sseServer := server.NewSSEServer(s)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("SSE Server failed: %v", err)
		}
	default:
		log.Fatalf("Unknown transport: %s", transportType)
	}
}
