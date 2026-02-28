package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/mlcmcp/task-manager/internal/db"
	"github.com/mlcmcp/task-manager/internal/models"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

var Store *db.TaskStore

func HandleTaskCreate(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
	subject, _ := args["subject"].(string)
	description, _ := args["description"].(string)
	activeForm, _ := args["activeForm"].(string)
	metadata, _ := args["metadata"].(map[string]interface{})

	if subject == "" || description == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Error: subject and description are required"},
			},
			IsError: true,
		}, nil, nil
	}

	task := Store.Create(subject, description, activeForm, metadata)

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Task created with ID: %s", task.ID)},
		},
	}, nil, nil
}

func HandleTaskUpdate(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
	taskID, _ := args["taskId"].(string)
	if taskID == "" {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: taskId is required"}},
			IsError: true,
		}, nil, nil
	}

	task, err := Store.Update(taskID, func(t *models.Task) {
		if s, ok := args["status"].(string); ok {
			t.Status = models.TaskStatus(s)
		}
		if s, ok := args["subject"].(string); ok {
			t.Subject = s
		}
		if d, ok := args["description"].(string); ok {
			t.Description = d
		}
		if a, ok := args["activeForm"].(string); ok {
			t.ActiveForm = a
		}
		if o, ok := args["owner"].(string); ok {
			t.Owner = o
		}
		if m, ok := args["metadata"].(map[string]interface{}); ok {
			if t.Metadata == nil {
				t.Metadata = make(map[string]interface{})
			}
			for k, v := range m {
				if v == nil {
					delete(t.Metadata, k)
				} else {
					t.Metadata[k] = v
				}
			}
		}
		if b, ok := args["addBlocks"].([]interface{}); ok {
			for _, id := range b {
				if idStr, ok := id.(string); ok {
					t.Blocks = append(t.Blocks, idStr)
				}
			}
		}
		if b, ok := args["addBlockedBy"].([]interface{}); ok {
			for _, id := range b {
				if idStr, ok := id.(string); ok {
					t.BlockedBy = append(t.BlockedBy, idStr)
				}
			}
		}
	})

	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: fmt.Sprintf("Error: %v", err)}},
			IsError: true,
		}, nil, nil
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Task %s updated. Current status: %s", task.ID, task.Status)},
		},
	}, nil, nil
}

func HandleTaskGet(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
	taskID, _ := args["taskId"].(string)
	task, ok := Store.Get(taskID)
	if !ok {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: "Error: task not found"}},
			IsError: true,
		}, nil, nil
	}

	data, _ := json.MarshalIndent(task, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(data)},
		},
	}, nil, nil
}

func HandleTaskList(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
	tasks := Store.List()
	
	header := "### Active Task List\n"
	if Store.IsPlanMode() {
		header = "### Active Task List [PLAN MODE ACTIVE]\n"
	}

	if len(tasks) == 0 {
		return &mcp.CallToolResult{
			Content: []mcp.Content{&mcp.TextContent{Text: header + "No tasks found."}},
		}, nil, nil
	}

	var sb strings.Builder
	sb.WriteString(header)
	for _, t := range tasks {
		statusIcon := "⏳"
		if t.Status == models.StatusInProgress {
			statusIcon = "🚀"
		} else if t.Status == models.StatusCompleted {
			statusIcon = "✅"
		}
		sb.WriteString(fmt.Sprintf("- [%s] **ID: %s** - %s (%s)\n", statusIcon, t.ID, t.Subject, t.Status))
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: sb.String()},
		},
	}, nil, nil
}

func HandleEnterPlanMode(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
	Store.SetPlanMode(true)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Switched to PLAN MODE. I will now explore and design before implementing. Please provide tasks or context."},
		},
	}, nil, nil
}

func HandleExitPlanMode(ctx context.Context, request *mcp.CallToolRequest, args map[string]interface{}) (*mcp.CallToolResult, any, error) {
	Store.SetPlanMode(false)
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Exited PLAN MODE. Ready to implement approved changes."},
		},
	}, nil, nil
}
