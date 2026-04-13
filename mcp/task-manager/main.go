package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/hmsoft0815/task-manager/internal/db"
	"github.com/hmsoft0815/task-manager/internal/handlers"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Ensure logging is on stderr for MCP protocol safety
	log.SetOutput(os.Stderr)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	v := flag.Bool("version", false, "Print version and exit")
	transportType := flag.String("transport", "stdio", "Transport type (stdio or sse)")
	addr := flag.String("addr", ":3001", "SSE address (if transport=sse)")
	dbPath := flag.String("db", "data/tasks.json", "Path to the tasks JSON file")
	flag.Parse()

	if *v {
		fmt.Println("task-manager 1.2.0")
		return
	}

	// 1. Initialize the Task Store
	store, err := db.NewTaskStore(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize task store: %v", err)
	}
	handlers.Store = store

	// 2. Initialize the MCP server
	s := server.NewMCPServer(
		"task-manager",
		"1.2.0",
	)

	// 3. Register tools
	registerTools(s)

	// 4. Run the server
	runServer(s, *transportType, *addr)
}

func registerTools(s *server.MCPServer) {
	s.AddTool(mcp.NewTool("task__task_create__mlc",
		mcp.WithDescription("Create a new task"),
		mcp.WithInputSchema[handlers.CreateTaskArgs](),
	), handlers.HandleTaskCreate)

	s.AddTool(mcp.NewTool("task__task_update__mlc",
		mcp.WithDescription("Update an existing task"),
		mcp.WithInputSchema[handlers.UpdateTaskArgs](),
	), handlers.HandleTaskUpdate)

	s.AddTool(mcp.NewTool("task__task_get__mlc",
		mcp.WithDescription("Get task details by ID"),
		mcp.WithInputSchema[handlers.GetTaskArgs](),
	), handlers.HandleTaskGet)

	s.AddTool(mcp.NewTool("task__task_list__mlc",
		mcp.WithDescription("List all tasks"),
	), handlers.HandleTaskList)

	s.AddTool(mcp.NewTool("task__enter_plan_mode__mlc",
		mcp.WithDescription("Enter research/planning mode"),
	), handlers.HandleEnterPlanMode)

	s.AddTool(mcp.NewTool("task__exit_plan_mode__mlc",
		mcp.WithDescription("Exit research/planning mode"),
	), handlers.HandleExitPlanMode)

	s.AddPrompt(mcp.NewPrompt("task__session_bootstrap__mlc",
		mcp.WithPromptDescription("Bootstraps a new coding session by analyzing pending tasks and establishing priorities."),
	), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
		return &mcp.GetPromptResult{
			Description: "Lead developer prioritized work plan bootstrap",
			Messages: []mcp.PromptMessage{
				{
					Role: mcp.RoleUser,
					Content: mcp.TextContent{
						Text: "You are a lead developer. Review the entire task list. Identify all 'pending' and 'in_progress' tasks. Analyze their dependencies and suggest a prioritized work plan for this session, explaining why each task should be tackled in the proposed order.",
					},
				},
			},
		}, nil
	})
}

func runServer(s *server.MCPServer, transportType, addr string) {
	switch transportType {
	case "stdio":
		log.Println("Task Manager MCP Server starting on stdio...")
		if err := server.ServeStdio(s); err != nil {
			log.Fatalf("Server failed: %v", err)
		}
	case "sse":
		log.Printf("Task Manager MCP Server starting on SSE at %s...", addr)
		sseServer := server.NewSSEServer(s)
		if err := sseServer.Start(addr); err != nil {
			log.Fatalf("SSE Server failed: %v", err)
		}
	default:
		log.Fatalf("Unknown transport: %s", transportType)
	}
}
