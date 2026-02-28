package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/mlcmcp/task-manager/internal/db"
	"github.com/mlcmcp/task-manager/internal/handlers"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: task-manager [options]\n")
		fmt.Fprintf(os.Stderr, "Options:\n")
		flag.PrintDefaults()
	}

	// 1. Setup flags
	dump := flag.Bool("dump", false, "Dump available tools as JSON and exit")
	dbPath := flag.String("db", "data/tasks.json", "Path to the tasks JSON file")
	flag.Parse()

	// 2. Initialize the Task Store
	store, err := db.NewTaskStore(*dbPath)
	if err != nil {
		log.Fatalf("Failed to initialize task store: %v", err)
	}
	handlers.Store = store

	// 3. Initialize the MCP server
	server := mcp.NewServer(&mcp.Implementation{
		Name:    "task-manager",
		Version: "1.0.0",
	}, nil)

	// 4. Define your tools
	tools := []mcp.Tool{
		{
			Name: "task__task_create__mlc",
			Description: `Use this tool to create a structured task list for your current coding session. This helps you track progress, organize complex tasks, and demonstrate thoroughness to the user.

## When to Use This Tool
- Complex multi-step tasks - When a task requires 3 or more distinct steps or actions.
- Non-trivial and complex tasks - Tasks that require careful planning or multiple operations.
- Plan mode - When using plan mode, create a task list to track the work.
- User explicitly requests todo list - When the user directly asks you to use the todo list.
- User provides multiple tasks - When users provide a list of things to be done (numbered or comma-separated).

## Task Fields
- subject: A brief, actionable title in imperative form (e.g., "Fix authentication bug in login flow")
- description: Detailed description of what needs to be done, including context and acceptance criteria.
- activeForm: Present continuous form shown in spinner when task is in_progress (e.g., "Fixing authentication bug").`,
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"subject": map[string]interface{}{
						"type":        "string",
						"description": "A brief, actionable title for the task (imperative form)",
					},
					"description": map[string]interface{}{
						"type":        "string",
						"description": "Detailed description of what needs to be done",
					},
					"activeForm": map[string]interface{}{
						"type":        "string",
						"description": "Present continuous form shown in status (e.g. 'Fixing bug...')",
					},
					"metadata": map[string]interface{}{
						"type": "object",
					},
				},
				"required": []string{"subject", "description"},
			},
		},
		{
			Name: "task__task_update__mlc",
			Description: `Use this tool to update a task in the task list.

## When to Use This Tool
- Mark tasks as in_progress: Immediately before you start working on a task.
- Mark tasks as completed: ONLY when you have FULLY accomplished the requirements.
- Manage dependencies: Use 'addBlocks' or 'addBlockedBy' to establish a workflow DAG.
- Update details: When requirements change or become clearer during implementation.

## Status Workflow
Status progresses: pending -> in_progress -> completed. Use 'deleted' to permanently remove a task created in error.`,
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"taskId": map[string]interface{}{"type": "string"},
					"status": map[string]interface{}{
						"type": "string",
						"enum": []string{"pending", "in_progress", "completed", "deleted"},
					},
					"subject":      map[string]interface{}{"type": "string"},
					"description":  map[string]interface{}{"type": "string"},
					"activeForm":   map[string]interface{}{"type": "string"},
					"addBlocks":    map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
					"addBlockedBy": map[string]interface{}{"type": "array", "items": map[string]interface{}{"type": "string"}},
				},
				"required": []string{"taskId"},
			},
		},
		{
			Name:        "task__task_get__mlc",
			Description: "Retrieve a task by its ID. Use this when you need the full description and context before starting work, or to understand task dependencies.",
			InputSchema: map[string]interface{}{
				"type": "object",
				"properties": map[string]interface{}{
					"taskId": map[string]interface{}{"type": "string"},
				},
				"required": []string{"taskId"},
			},
		},
		{
			Name:        "task__task_list__mlc",
			Description: "List all tasks in the task list. Use this to see available tasks, check overall progress, or find unblocked work after completing a task.",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name: "session__enter_plan_mode__mlc",
			Description: `Use this tool proactively when you're about to start a non-trivial implementation task. Getting user sign-off on your approach before writing code prevents wasted effort and ensures alignment.

## What Happens in Plan Mode
1. Thoroughly explore the codebase using search and read tools.
2. Design an implementation approach and create a task list.
3. Present your plan to the user for approval.
4. Exit plan mode with exit_plan_mode when ready to implement.`,
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
		{
			Name:        "session__exit_plan_mode__mlc",
			Description: "Exit Plan Mode and transition back to Implementation Mode. Call this once the user has approved your design and you are ready to start making changes.",
			InputSchema: map[string]interface{}{
				"type":       "object",
				"properties": map[string]interface{}{},
			},
		},
	}

	// 5. Register handlers
	mcp.AddTool(server, &tools[0], handlers.HandleTaskCreate)
	mcp.AddTool(server, &tools[1], handlers.HandleTaskUpdate)
	mcp.AddTool(server, &tools[2], handlers.HandleTaskGet)
	mcp.AddTool(server, &tools[3], handlers.HandleTaskList)
	mcp.AddTool(server, &tools[4], handlers.HandleEnterPlanMode)
	mcp.AddTool(server, &tools[5], handlers.HandleExitPlanMode)

	// 6. Handle discovery dump
	if *dump {
		json.NewEncoder(os.Stdout).Encode(tools)
		return
	}

	// 7. Run the server
	log.Printf("Task Manager MCP Server starting...")
	if err := server.Run(context.Background(), &mcp.StdioTransport{}); err != nil {
		log.Fatalf("Server failed: %v", err)
	}
}
