package domain

import "fmt"

// ProjectNotFoundError is raised when a project lookup fails.
type ProjectNotFoundError struct {
	ProjectID string
}

func (e *ProjectNotFoundError) Error() string {
	return fmt.Sprintf("Project not found: %s", e.ProjectID)
}

// ProjectValidationError carries field-level validation errors.
type ProjectValidationError struct {
	Errors map[string]string
}

func (e *ProjectValidationError) Error() string {
	return fmt.Sprintf("Validation failed: %v", e.Errors)
}

// ConfirmationRequiredError is raised when a destructive action needs explicit confirmation.
type ConfirmationRequiredError struct {
	Hint string
}

func (e *ConfirmationRequiredError) Error() string {
	return "Permanent deletion requires confirmation"
}

func NewConfirmationRequired() *ConfirmationRequiredError {
	return &ConfirmationRequiredError{
		Hint: "Add ?confirm=true to permanently delete this project",
	}
}
