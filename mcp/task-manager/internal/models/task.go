package models

import "time"

type TaskStatus string

const (
	StatusPending    TaskStatus = "pending"
	StatusInProgress TaskStatus = "in_progress"
	StatusCompleted  TaskStatus = "completed"
	StatusDeleted    TaskStatus = "deleted"
)

type Task struct {
	ID          string                 `json:"id"`
	Subject     string                 `json:"subject"`
	Description string                 `json:"description"`
	ActiveForm  string                 `json:"activeForm"`
	Status      TaskStatus             `json:"status"`
	Owner       string                 `json:"owner,omitempty"`
	Blocks      []string               `json:"blocks,omitempty"`
	BlockedBy   []string               `json:"blockedBy,omitempty"`
	Metadata    map[string]interface{} `json:"metadata,omitempty"`
	CreatedAt   time.Time              `json:"createdAt"`
	UpdatedAt   time.Time              `json:"updatedAt"`
}
