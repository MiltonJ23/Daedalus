package domain_test

import (
	"testing"

	"github.com/MiltonJ23/Daedalus/AuthService/internal/domain"
	"github.com/go-playground/validator/v10"
)

func TestRegisterRequestValidation(t *testing.T) {
	validate := validator.New()

	tests := []struct {
		name    string
		req     domain.RegisterRequest
		wantErr bool
	}{
		{
			name:    "Valid request",
			req:     domain.RegisterRequest{Email: "entrepreneur@daedalus.cm", Password: "SecurePassword123!"},
			wantErr: false,
		},
		{
			name:    "Invalid email format",
			req:     domain.RegisterRequest{Email: "not-an-email", Password: "SecurePassword123!"},
			wantErr: true,
		},
		{
			name:    "Password too short",
			req:     domain.RegisterRequest{Email: "test@daedalus.cm", Password: "short"},
			wantErr: true,
		},
		{
			name:    "SQL Injection attempt in email",
			req:     domain.RegisterRequest{Email: "test@daedalus.cm' OR '1'='1", Password: "SecurePassword123!"},
			wantErr: true, // Rejeté par le validateur d'email
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := validate.Struct(tt.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("Struct validation error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
