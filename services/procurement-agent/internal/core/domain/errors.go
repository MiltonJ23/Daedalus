package domain

import "fmt"

// SearchNotFoundError is raised when a procurement search lookup fails.
type SearchNotFoundError struct {
	SearchID string
}

func (e *SearchNotFoundError) Error() string {
	return fmt.Sprintf("Procurement search not found: %s", e.SearchID)
}

// ResultNotFoundError is raised when an equipment result lookup fails.
type ResultNotFoundError struct {
	ResultID string
}

func (e *ResultNotFoundError) Error() string {
	return fmt.Sprintf("Equipment result not found: %s", e.ResultID)
}

// ValidationError carries field-level validation errors.
type ValidationError struct {
	Errors map[string]string
}

func (e *ValidationError) Error() string {
	return fmt.Sprintf("Validation failed: %v", e.Errors)
}

// InvalidDecisionError is raised when an unsupported decision value is supplied.
type InvalidDecisionError struct {
	Value string
}

func (e *InvalidDecisionError) Error() string {
	return fmt.Sprintf("Invalid decision %q: must be 'approved' or 'rejected'", e.Value)
}
