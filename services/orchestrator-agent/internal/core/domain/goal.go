package domain

import (
	"time"

	"github.com/google/uuid"
)

// Goal lifecycle states.
const (
	GoalStatusPending   = "pending"
	GoalStatusRunning   = "running"
	GoalStatusCompleted = "completed"
	GoalStatusFailed    = "failed"
)

// SubTask lifecycle states.
const (
	TaskStatusPending    = "pending"
	TaskStatusDispatched = "dispatched"
	TaskStatusInProgress = "in_progress"
	TaskStatusCompleted  = "completed"
	TaskStatusFailed     = "failed"
)

// Canonical sub-task types — one Redis Stream per type.
const (
	TaskTypeProcurementSearch = "procurement_search"
	TaskTypeLayoutGeneration  = "layout_generation"
	TaskTypeThreeDRender      = "three_d_render"
	TaskTypeCostAnalysis      = "cost_analysis"
)

// Goal — a high-level user objective decomposed into a sub-task graph.
type Goal struct {
	ID          string    `json:"id"`
	UserID      string    `json:"user_id"`
	ProjectID   string    `json:"project_id"`
	Description string    `json:"description"`
	Status      string    `json:"status"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

func NewGoal(userID, projectID, description string) Goal {
	now := time.Now().UTC()
	return Goal{
		ID:          uuid.New().String(),
		UserID:      userID,
		ProjectID:   projectID,
		Description: description,
		Status:      GoalStatusPending,
		CreatedAt:   now,
		UpdatedAt:   now,
	}
}

// SubTask — atomic unit of work delegated to a specialised agent.
type SubTask struct {
	ID        string                 `json:"id"`
	GoalID    string                 `json:"goal_id"`
	Type      string                 `json:"type"`
	Status    string                 `json:"status"`
	Payload   map[string]interface{} `json:"payload"`
	DependsOn []string               `json:"depends_on"`
	StreamID  string                 `json:"stream_id"`
	Error     string                 `json:"error,omitempty"`
	CreatedAt time.Time              `json:"created_at"`
	UpdatedAt time.Time              `json:"updated_at"`
}

func NewSubTask(goalID, taskType string, payload map[string]interface{}, deps []string) SubTask {
	now := time.Now().UTC()
	if payload == nil {
		payload = map[string]interface{}{}
	}
	if deps == nil {
		deps = []string{}
	}
	return SubTask{
		ID:        uuid.New().String(),
		GoalID:    goalID,
		Type:      taskType,
		Status:    TaskStatusPending,
		Payload:   payload,
		DependsOn: deps,
		CreatedAt: now,
		UpdatedAt: now,
	}
}

// ApplyStatus mutates the sub-task with a validated lifecycle transition.
func (t *SubTask) ApplyStatus(status, errMsg string) error {
	switch status {
	case TaskStatusPending, TaskStatusDispatched, TaskStatusInProgress, TaskStatusCompleted, TaskStatusFailed:
		t.Status = status
		t.Error = errMsg
		t.UpdatedAt = time.Now().UTC()
		return nil
	default:
		return &InvalidStatusError{Value: status}
	}
}

// TaskGraph bundles a goal with all its sub-tasks for API responses + logging.
type TaskGraph struct {
	Goal  Goal      `json:"goal"`
	Tasks []SubTask `json:"tasks"`
}
