package errors

import (
	"fmt"
	"net/http"
)

type ErrorCode string

const (
	CodeValidation      ErrorCode = "VALIDATION_ERROR"
	CodeNotFound        ErrorCode = "NOT_FOUND"
	CodeUnauthorized    ErrorCode = "UNAUTHORIZED"
	CodeForbidden       ErrorCode = "FORBIDDEN"
	CodeConflict        ErrorCode = "CONFLICT"
	CodeInternal        ErrorCode = "INTERNAL_SERVER_ERROR"
	CodePayment         ErrorCode = "PAYMENT_ERROR"
	CodeServiceUnavailable ErrorCode = "SERVICE_UNAVAILABLE"
	CodeBadGateway      ErrorCode = "BAD_GATEWAY"
)

type DaedalusError struct {
	Code    ErrorCode `json:"code"`
	Message string    `json:"message"`
	Details string    `json:"details,omitempty"`
	Status  int       `json:"status"`
}

func (e *DaedalusError) Error() string {
	return fmt.Sprintf("[%s] %s", e.Code, e.Message)
}

func (e *DaedalusError) HTTPStatus() int {
	return e.Status
}

func NewValidation(message string, details string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeValidation,
		Message: message,
		Details: details,
		Status:  http.StatusBadRequest,
	}
}

func NewNotFound(resource string, id string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeNotFound,
		Message: fmt.Sprintf("%s not found", resource),
		Details: fmt.Sprintf("ID: %s", id),
		Status:  http.StatusNotFound,
	}
}

func NewUnauthorized(reason string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeUnauthorized,
		Message: "Unauthorized access",
		Details: reason,
		Status:  http.StatusUnauthorized,
	}
}

func NewForbidden(reason string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeForbidden,
		Message: "Access forbidden",
		Details: reason,
		Status:  http.StatusForbidden,
	}
}

func NewConflict(message string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeConflict,
		Message: message,
		Status:  http.StatusConflict,
	}
}

func NewPaymentError(message string, details string) *DaedalusError {
	return &DaedalusError{
		Code:    CodePayment,
		Message: message,
		Details: details,
		Status:  http.StatusPaymentRequired,
	}
}

func NewInternal(message string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeInternal,
		Message: message,
		Status:  http.StatusInternalServerError,
	}
}

func NewServiceUnavailable(service string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeServiceUnavailable,
		Message: fmt.Sprintf("%s service temporarily unavailable", service),
		Status:  http.StatusServiceUnavailable,
	}
}

func NewBadGateway(details string) *DaedalusError {
	return &DaedalusError{
		Code:    CodeBadGateway,
		Message: "Bad gateway",
		Details: details,
		Status:  http.StatusBadGateway,
	}
}

func FromError(err error) *DaedalusError {
	if de, ok := err.(*DaedalusError); ok {
		return de
	}
	return NewInternal(err.Error())
}

func IsNotFound(err error) bool {
	if de, ok := err.(*DaedalusError); ok {
		return de.Code == CodeNotFound
	}
	return false
}

func IsUnauthorized(err error) bool {
	if de, ok := err.(*DaedalusError); ok {
		return de.Code == CodeUnauthorized
	}
	return false
}

func IsForbidden(err error) bool {
	if de, ok := err.(*DaedalusError); ok {
		return de.Code == CodeForbidden
	}
	return false
}

func IsValidation(err error) bool {
	if de, ok := err.(*DaedalusError); ok {
		return de.Code == CodeValidation
	}
	return false
}
