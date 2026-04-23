package errors

import (
	"net/http"
	"testing"
)

func TestNewValidation(t *testing.T) {
	err := NewValidation("invalid input", "field: email")
	if err.Code != CodeValidation {
		t.Errorf("code = %v, want %v", err.Code, CodeValidation)
	}
	if err.Status != http.StatusBadRequest {
		t.Errorf("status = %d, want %d", err.Status, http.StatusBadRequest)
	}
}

func TestNewNotFound(t *testing.T) {
	err := NewNotFound("project", "123")
	if err.Code != CodeNotFound {
		t.Errorf("code = %v, want %v", err.Code, CodeNotFound)
	}
	if err.Status != http.StatusNotFound {
		t.Errorf("status = %d, want %d", err.Status, http.StatusNotFound)
	}
}

func TestNewUnauthorized(t *testing.T) {
	err := NewUnauthorized("token expired")
	if err.Code != CodeUnauthorized {
		t.Errorf("code = %v, want %v", err.Code, CodeUnauthorized)
	}
	if err.Status != http.StatusUnauthorized {
		t.Errorf("status = %d, want %d", err.Status, http.StatusUnauthorized)
	}
}

func TestNewForbidden(t *testing.T) {
	err := NewForbidden("insufficient permissions")
	if err.Code != CodeForbidden {
		t.Errorf("code = %v, want %v", err.Code, CodeForbidden)
	}
	if err.Status != http.StatusForbidden {
		t.Errorf("status = %d, want %d", err.Status, http.StatusForbidden)
	}
}

func TestNewPaymentError(t *testing.T) {
	err := NewPaymentError("payment failed", "card declined")
	if err.Code != CodePayment {
		t.Errorf("code = %v, want %v", err.Code, CodePayment)
	}
	if err.Status != http.StatusPaymentRequired {
		t.Errorf("status = %d, want %d", err.Status, http.StatusPaymentRequired)
	}
}

func TestDaedalusErrorError(t *testing.T) {
	err := NewNotFound("project", "123")
	errStr := err.Error()
	if errStr == "" {
		t.Error("Error() returned empty string")
	}
}

func TestIsNotFound(t *testing.T) {
	err := NewNotFound("project", "123")
	if !IsNotFound(err) {
		t.Error("IsNotFound() = false, want true")
	}

	err2 := NewValidation("invalid", "")
	if IsNotFound(err2) {
		t.Error("IsNotFound() = true, want false")
	}
}

func TestIsUnauthorized(t *testing.T) {
	err := NewUnauthorized("token expired")
	if !IsUnauthorized(err) {
		t.Error("IsUnauthorized() = false, want true")
	}
}

func TestIsForbidden(t *testing.T) {
	err := NewForbidden("insufficient permissions")
	if !IsForbidden(err) {
		t.Error("IsForbidden() = false, want true")
	}
}

func TestIsValidation(t *testing.T) {
	err := NewValidation("invalid input", "")
	if !IsValidation(err) {
		t.Error("IsValidation() = false, want true")
	}
}

func TestFromError(t *testing.T) {
	original := NewNotFound("project", "123")
	converted := FromError(original)
	if converted.Code != CodeNotFound {
		t.Errorf("code = %v, want %v", converted.Code, CodeNotFound)
	}
}

func TestHTTPStatus(t *testing.T) {
	tests := []struct {
		name   string
		err    *DaedalusError
		status int
	}{
		{"validation", NewValidation("", ""), http.StatusBadRequest},
		{"not found", NewNotFound("", ""), http.StatusNotFound},
		{"unauthorized", NewUnauthorized(""), http.StatusUnauthorized},
		{"forbidden", NewForbidden(""), http.StatusForbidden},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.err.HTTPStatus() != tt.status {
				t.Errorf("HTTPStatus() = %d, want %d", tt.err.HTTPStatus(), tt.status)
			}
		})
	}
}
