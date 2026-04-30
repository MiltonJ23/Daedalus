package domain

import "fmt"

type GoalNotFoundError struct{ GoalID string }

func (e *GoalNotFoundError) Error() string {
	return fmt.Sprintf("goal not found: %s", e.GoalID)
}

type TaskNotFoundError struct{ TaskID string }

func (e *TaskNotFoundError) Error() string {
	return fmt.Sprintf("task not found: %s", e.TaskID)
}

type ValidationError struct{ Errors map[string]string }

func (e *ValidationError) Error() string {
	return fmt.Sprintf("validation failed: %v", e.Errors)
}

type InvalidStatusError struct{ Value string }

func (e *InvalidStatusError) Error() string {
	return fmt.Sprintf("invalid status: %s", e.Value)
}
