package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/mlcmcp/memory-server/internal/handlers"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: memory-server [options]\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	dump := flag.Bool("dump", false, "Dump tool definitions as JSON and exit")
	flag.Parse()

	home, _ := os.UserHomeDir()
	dbDir := filepath.Join(home, ".local", "share", "mcp-proxy")
	os.MkdirAll(dbDir, 0755)
	dbPath := filepath.Join(dbDir, "memory.db")

	handler, err := handlers.NewMemoryHandler(dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize SQLite: %v", err)
	}

	server := mcp.NewServer(&mcp.Implementation{
		Name:    "memory-server",
		Version: "1.1.2",
	}, nil)

	// FLAT TOOLS - Maximum compatibility
	tools := []mcp.Tool{
		{
			Name:        "memory__memorize__mlc",
			Description: "Store a new fact or observation about an entity",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"entity":      map[string]interface{}{"type": "string", "description": "The name of the thing (e.g. 'Oly' or 'Project')"},
					"category":    map[string]interface{}{"type": "string", "description": "Category (e.g. 'Person', 'Setting')"},
					"observation": map[string]interface{}{"type": "string", "description": "The actual fact to remember"},
				},
				"required":             []string{"entity", "observation"},
			},
		},
		{
			Name:        "memory__search_nodes__mlc",
			Description: "Search for stored facts in the knowledge graph",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"query": map[string]interface{}{"type": "string", "description": "The search term"},
				},
				"required":             []string{"query"},
			},
		},
		{
			Name:        "memory__read_graph__mlc",
			Description: "Read all stored memories",
			InputSchema: map[string]interface{}{
				"type":                 "object",
				"properties":           map[string]interface{}{},
			},
		},
	}

	textResult := func(text string) *mcp.CallToolResult {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: text},
			},
		}
	}

	// Simplified Handlers
	mcp.AddTool(server, &tools[0], func(ctx context.Context, req *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
		entity, _ := args["entity"].(string)
		entType, _ := args["category"].(string)
		obs, _ := args["observation"].(string)

		hArgs := map[string]interface{}{
			"entities": []interface{}{
				map[string]interface{}{
					"name": entity,
					"entityType": entType,
					"observations": []interface{}{obs},
				},
			},
		}
		res, err := handler.CreateEntities(hArgs)
		if err != nil { return nil, nil, err }
		return textResult(res.(string)), nil, nil
	})

	mcp.AddTool(server, &tools[1], func(ctx context.Context, req *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
		res, err := handler.SearchNodes(args)
		if err != nil { return nil, nil, err }
		return textResult(res.(string)), nil, nil
	})

	mcp.AddTool(server, &tools[2], func(ctx context.Context, req *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
		res, err := handler.ReadGraph(args)
		if err != nil { return nil, nil, err }
		return textResult(res.(string)), nil, nil
	})

	if *dump {
		json.NewEncoder(os.Stdout).Encode(tools)
		return
	}

	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}