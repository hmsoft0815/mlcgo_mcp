package handlers

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ExampleArgs defines the expected arguments for our example tool.
type ExampleArgs struct {
	Input string `json:"input" jsonschema:"description=Some input parameter"`
}

// HandleExample is the business logic for our example tool.
func HandleExample(ctx context.Context, req *mcp.CallToolRequest, args ExampleArgs) (*mcp.CallToolResult, any, error) {
	// 1. Process the request
	responseText := fmt.Sprintf("Skeleton received: %s", args.Input)

	// 2. Return a standard MCP content result using pointers to concrete types
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{
				Text: responseText,
			},
		},
	}, nil, nil
}
