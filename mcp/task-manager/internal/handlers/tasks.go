package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/hmsoft0815/task-manager/internal/db"
	"github.com/hmsoft0815/task-manager/internal/models"
	"github.com/mark3labs/mcp-go/mcp"
)

var Store *db.TaskStore

type CreateTaskArgs struct {
	Subject     string                 `json:"subject" jsonschema:"description=A brief, actionable title for the task (imperative form)"`
	Description string                 `json:"description" jsonschema:"description=Detailed description of what needs to be done"`
	ActiveForm  string                 `json:"activeForm,omitempty" jsonschema:"description=Present continuous form shown in status (e.g. 'Fixing bug...')"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
}

type UpdateTaskArgs struct {
	TaskID       string   `json:"taskId"`
	Status       string   `json:"status,omitempty" jsonschema:"enum=pending,enum=in_progress,enum=completed,enum=deleted"`
	Subject      string   `json:"subject,omitempty"`
	Description  string   `json:"description,omitempty"`
	ActiveForm   string   `json:"activeForm,omitempty"`
	AddBlocks    []string `json:"addBlocks,omitempty"`
	AddBlockedBy []string `json:"addBlockedBy,omitempty"`
}

type GetTaskArgs struct {
	TaskID string `json:"taskId"`
}

func HandleTaskCreate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args CreateTaskArgs
	if err := request.BindArguments(&args); err != nil {
		return nil, err
	}

	if args.Subject == "" || args.Description == "" {
		return mcp.NewToolResultError("subject and description are required"), nil
	}

	task := Store.Create(args.Subject, args.Description, args.ActiveForm, args.Metadata)

	return mcp.NewToolResultText(fmt.Sprintf("Task created with ID: %s", task.ID)), nil
}

func HandleTaskUpdate(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args UpdateTaskArgs
	if err := request.BindArguments(&args); err != nil {
		return nil, err
	}

	if args.TaskID == "" {
		return mcp.NewToolResultError("taskId is required"), nil
	}

	task, err := Store.Update(args.TaskID, func(t *models.Task) {
		if args.Status != "" {
			t.Status = models.TaskStatus(args.Status)
		}
		if args.Subject != "" {
			t.Subject = args.Subject
		}
		if args.Description != "" {
			t.Description = args.Description
		}
		if args.ActiveForm != "" {
			t.ActiveForm = args.ActiveForm
		}
		if len(args.AddBlocks) > 0 {
			t.Blocks = append(t.Blocks, args.AddBlocks...)
		}
		if len(args.AddBlockedBy) > 0 {
			t.BlockedBy = append(t.BlockedBy, args.AddBlockedBy...)
		}
	})

	if err != nil {
		return mcp.NewToolResultErrorFromErr("Update failed", err), nil
	}

	return mcp.NewToolResultText(fmt.Sprintf("Task %s updated. Current status: %s", task.ID, task.Status)), nil
}

func HandleTaskGet(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	var args GetTaskArgs
	if err := request.BindArguments(&args); err != nil {
		return nil, err
	}

	task, ok := Store.Get(args.TaskID)
	if !ok {
		return mcp.NewToolResultError("task not found"), nil
	}

	data, _ := json.MarshalIndent(task, "", "  ")
	return mcp.NewToolResultText(string(data)), nil
}

func HandleTaskList(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	tasks := Store.List()

	header := "### Active Task List\n"
	if Store.IsPlanMode() {
		header = "### Active Task List [PLAN MODE ACTIVE]\n"
	}

	if len(tasks) == 0 {
		return mcp.NewToolResultText(header + "No tasks found."), nil
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

	return mcp.NewToolResultText(sb.String()), nil
}

func HandleEnterPlanMode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	Store.SetPlanMode(true)
	return mcp.NewToolResultText("Switched to PLAN MODE. I will now explore and design before implementing. Please provide tasks or context."), nil
}

func HandleExitPlanMode(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	Store.SetPlanMode(false)
	return mcp.NewToolResultText("Exited PLAN MODE. Ready to implement approved changes."), nil
}
